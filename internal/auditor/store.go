package auditor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// JSONGraphStore implements GraphStore for graphify-out/graph.json metadata updates.
type JSONGraphStore struct {
	mu sync.Mutex
}

// NewJSONGraphStore initializes a new JSONGraphStore.
func NewJSONGraphStore() *JSONGraphStore {
	return &JSONGraphStore{}
}

// SaveGraphData saves the entire GraphData structure to targetPath atomically.
func (s *JSONGraphStore) SaveGraphData(ctx context.Context, targetPath string, graph *GraphData) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	cleanPath := filepath.Clean(targetPath)
	outBytes, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal graph data: %w", err)
	}

	tmpFile := cleanPath + ".tmp"
	if err := os.WriteFile(tmpFile, outBytes, 0644); err != nil {
		return fmt.Errorf("failed to write tmp graph file: %w", err)
	}

	if err := os.Rename(tmpFile, cleanPath); err != nil {
		_ = os.Remove(tmpFile)
		return fmt.Errorf("atomic rename failed for %s: %w", cleanPath, err)
	}

	return nil
}

// SaveAuditedEdges updates graphify-out/graph.json with provenance and confidence annotations.
func (s *JSONGraphStore) SaveAuditedEdges(ctx context.Context, targetPath string, edges []AuditedEdge) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	cleanPath := filepath.Clean(targetPath)
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If file does not exist, build an initial audit JSON report artifact
			initialStructure := map[string]interface{}{
				"audited_edges": edges,
				"total_audited": len(edges),
			}
			outBytes, marshalErr := json.MarshalIndent(initialStructure, "", "  ")
			if marshalErr != nil {
				return fmt.Errorf("failed to marshal initial audit structure: %w", marshalErr)
			}
			return os.WriteFile(cleanPath, outBytes, 0644)
		}
		return fmt.Errorf("failed to read graph store target %s: %w", cleanPath, err)
	}

	var graphData GraphData
	if err := json.Unmarshal(data, &graphData); err != nil {
		return fmt.Errorf("failed to unmarshal graph JSON at %s: %w", cleanPath, err)
	}

	// Index audited edges for quick lookup
	auditMap := make(map[string]AuditedEdge)
	for _, edge := range edges {
		auditMap[edge.ID] = edge
	}

	// Update matching link edges
	updatedLinks := make([]map[string]interface{}, 0, len(graphData.Links))
	for _, link := range graphData.Links {
		source, _ := link["source"].(string)
		target, _ := link["target"].(string)
		linkID := fmt.Sprintf("%s->%s", source, target)

		if audited, exists := auditMap[linkID]; exists {
			link["provenance"] = audited.ProvenanceStatus
			link["confidence"] = audited.Confidence
			if audited.ASTNodePattern != "" {
				link["ast_pattern"] = audited.ASTNodePattern
			}
			// Omit pruned phantom edges from output graph links if confidence is 0.0
			if audited.Confidence > 0.0 {
				updatedLinks = append(updatedLinks, link)
			}
		} else {
			updatedLinks = append(updatedLinks, link)
		}
	}

	graphData.Links = updatedLinks

	outBytes, err := json.MarshalIndent(graphData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated graph JSON: %w", err)
	}

	tmpFile := cleanPath + ".tmp"
	if err := os.WriteFile(tmpFile, outBytes, 0644); err != nil {
		return fmt.Errorf("failed to write tmp graph file: %w", err)
	}

	if err := os.Rename(tmpFile, cleanPath); err != nil {
		_ = os.Remove(tmpFile)
		return fmt.Errorf("atomic rename failed for %s: %w", cleanPath, err)
	}

	return nil
}

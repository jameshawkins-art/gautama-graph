package auditor

import (
	"fmt"
	"sort"
)

// TarjanSCCDetector implements CycleDetector using Tarjan's Strongly Connected Components algorithm.
type TarjanSCCDetector struct{}

// NewTarjanSCCDetector constructs a new TarjanSCCDetector instance.
func NewTarjanSCCDetector() *TarjanSCCDetector {
	return &TarjanSCCDetector{}
}

// FindCycles identifies circular reference chains and strongly connected components across the document link graph.
func (d *TarjanSCCDetector) FindCycles(graph *DocGraph) *CycleReport {
	report := &CycleReport{
		Cycles: make([]CircularCycle, 0),
	}

	if graph == nil || len(graph.Nodes) == 0 {
		return report
	}

	// 1. Build adjacency list representation
	adj := make(map[string][]string)
	allNodes := make(map[string]bool)

	for id := range graph.Nodes {
		allNodes[id] = true
		adj[id] = make([]string, 0)
	}

	for _, edge := range graph.Edges {
		allNodes[edge.SourceID] = true
		allNodes[edge.TargetID] = true
		adj[edge.SourceID] = append(adj[edge.SourceID], edge.TargetID)

		// Check for direct self-loop
		if edge.SourceID == edge.TargetID {
			cycleID := fmt.Sprintf("cycle-self-%s", edge.SourceID)
			report.Cycles = append(report.Cycles, CircularCycle{
				CycleID:  cycleID,
				Length:   1,
				DocChain: []string{edge.SourceID, edge.SourceID},
			})
		}
	}

	// 2. Tarjan's Algorithm State
	index := make(map[string]int)
	lowlink := make(map[string]int)
	onStack := make(map[string]bool)
	var stack []string
	currentIndex := 0

	for node := range allNodes {
		index[node] = -1
	}

	// Helper function for recursive traversal
	var strongConnect func(v string)
	strongConnect = func(v string) {
		index[v] = currentIndex
		lowlink[v] = currentIndex
		currentIndex++
		stack = append(stack, v)
		onStack[v] = true

		for _, w := range adj[v] {
			if index[w] == -1 {
				strongConnect(w)
				if lowlink[w] < lowlink[v] {
					lowlink[v] = lowlink[w]
				}
			} else if onStack[w] {
				if index[w] < lowlink[v] {
					lowlink[v] = index[w]
				}
			}
		}

		// Root of an SCC
		if lowlink[v] == index[v] {
			var scc []string
			for {
				topIdx := len(stack) - 1
				w := stack[topIdx]
				stack = stack[:topIdx]
				onStack[w] = false
				scc = append(scc, w)
				if w == v {
					break
				}
			}

			// If SCC has > 1 nodes, it is a circular dependency loop
			if len(scc) > 1 {
				// Reverse for topological cycle representation and close the loop
				for i, j := 0, len(scc)-1; i < j; i, j = i+1, j-1 {
					scc[i], scc[j] = scc[j], scc[i]
				}
				docChain := append([]string(nil), scc...)
				docChain = append(docChain, scc[0]) // close cycle

				cycleID := fmt.Sprintf("cycle-%d-%s", len(report.Cycles)+1, scc[0])
				report.Cycles = append(report.Cycles, CircularCycle{
					CycleID:  cycleID,
					Length:   len(scc),
					DocChain: docChain,
				})
			}
		}
	}

	// Sort node keys for deterministic graph traversal order
	sortedNodes := make([]string, 0, len(allNodes))
	for n := range allNodes {
		sortedNodes = append(sortedNodes, n)
	}
	sort.Strings(sortedNodes)

	for _, node := range sortedNodes {
		if index[node] == -1 {
			strongConnect(node)
		}
	}

	report.TotalCycles = len(report.Cycles)
	return report
}

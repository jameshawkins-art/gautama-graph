package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/jameshawkins-art/gautama-graph/internal/runner"
)

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSliceFlag) Set(val string) error {
	*s = append(*s, val)
	return nil
}

// handleQueryCommand parses and executes the `query` subcommand.
func handleQueryCommand(args []string) error {
	fs := flag.NewFlagSet("query", flag.ContinueOnError)
	dfsFlag := fs.Bool("dfs", false, "Use depth-first traversal instead of breadth-first")
	budgetFlag := fs.Int("budget", 2000, "Cap output at N tokens")
	workspaceFlag := fs.String("workspace", "", "Path to workspace root (defaults to CWD)")
	graphFlag := fs.String("graph", "", "Path to graph.json (defaults to <workspace>/graphify-out/graph.json)")
	jsonFlag := fs.Bool("json", false, "Emit machine-readable JSON output")

	var contextFlags stringSliceFlag
	fs.Var(&contextFlags, "context", "Filter by edge context (repeatable)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 || strings.TrimSpace(fs.Args()[0]) == "" {
		return fmt.Errorf("question argument is required. Usage: gautama-graph query \"<question>\" [flags]")
	}

	question := fs.Args()[0]
	opts := runner.QueryOptions{
		WorkspaceRoot: *workspaceFlag,
		GraphPath:     *graphFlag,
		DFS:           *dfsFlag,
		BudgetTokens:  *budgetFlag,
		JSONOutput:    *jsonFlag,
		Context:       contextFlags,
	}

	svc := runner.NewDefaultQueryService(nil, nil)
	res, err := svc.Query(context.Background(), question, opts)
	if err != nil {
		return err
	}

	fmt.Print(res.Output)
	if !strings.HasSuffix(res.Output, "\n") {
		fmt.Println()
	}
	return nil
}

// handlePathCommand parses and executes the `path` subcommand.
func handlePathCommand(args []string) error {
	fs := flag.NewFlagSet("path", flag.ContinueOnError)
	workspaceFlag := fs.String("workspace", "", "Path to workspace root (defaults to CWD)")
	graphFlag := fs.String("graph", "", "Path to graph.json (defaults to <workspace>/graphify-out/graph.json)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 2 || strings.TrimSpace(fs.Args()[0]) == "" || strings.TrimSpace(fs.Args()[1]) == "" {
		return fmt.Errorf("both start and end nodes are required. Usage: gautama-graph path \"<nodeA>\" \"<nodeB>\" [flags]")
	}

	nodeA := fs.Args()[0]
	nodeB := fs.Args()[1]

	opts := runner.PathOptions{
		WorkspaceRoot: *workspaceFlag,
		GraphPath:     *graphFlag,
	}

	svc := runner.NewDefaultQueryService(nil, nil)
	res, err := svc.Path(context.Background(), nodeA, nodeB, opts)
	if err != nil {
		return err
	}

	fmt.Print(res.Output)
	if !strings.HasSuffix(res.Output, "\n") {
		fmt.Println()
	}
	return nil
}

// handleExplainCommand parses and executes the `explain` subcommand.
func handleExplainCommand(args []string) error {
	fs := flag.NewFlagSet("explain", flag.ContinueOnError)
	workspaceFlag := fs.String("workspace", "", "Path to workspace root (defaults to CWD)")
	graphFlag := fs.String("graph", "", "Path to graph.json (defaults to <workspace>/graphify-out/graph.json)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 || strings.TrimSpace(fs.Args()[0]) == "" {
		return fmt.Errorf("concept identifier is required. Usage: gautama-graph explain \"<concept>\" [flags]")
	}

	concept := fs.Args()[0]
	opts := runner.ExplainOptions{
		WorkspaceRoot: *workspaceFlag,
		GraphPath:     *graphFlag,
	}

	svc := runner.NewDefaultQueryService(nil, nil)
	res, err := svc.Explain(context.Background(), concept, opts)
	if err != nil {
		return err
	}

	fmt.Print(res.Output)
	if !strings.HasSuffix(res.Output, "\n") {
		fmt.Println()
	}
	return nil
}

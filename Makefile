.PHONY: test audit-docs audit-ast audit graphify-update audit-remediate graphify-query graphify-path graphify-explain

PATH := $(HOME)/go/bin:$(HOME)/.local/bin:$(PATH)

test:
	GOWORK=off go test -v -race ./...

# Knowledge Graph Query Targets (Zero-Host-Dependency)
graphify-query:
	@if [ -z "$(Q)" ]; then \
		echo "❌ Error: Q is required. Usage: make graphify-query Q=\"<question>\" [DFS=1] [BUDGET=2000]"; \
		exit 1; \
	fi
	GOWORK=off go run ./cmd/gautama-graph query "$(Q)" \
		$(if $(filter 1 true,$(DFS)),--dfs,) \
		$(if $(BUDGET),--budget $(BUDGET),)

graphify-path:
	@if [ -z "$(A)" ] || [ -z "$(B)" ]; then \
		echo "❌ Error: A and B are required. Usage: make graphify-path A=\"<symbol1>\" B=\"<symbol2>\""; \
		exit 1; \
	fi
	GOWORK=off go run ./cmd/gautama-graph path "$(A)" "$(B)"

graphify-explain:
	@if [ -z "$(C)" ]; then \
		echo "❌ Error: C is required. Usage: make graphify-explain C=\"<concept>\""; \
		exit 1; \
	fi
	GOWORK=off go run ./cmd/gautama-graph explain "$(C)"

# Pipeline & Audit Targets
audit-docs:
	GOWORK=off go run ./cmd/graphify-doc-audit -strict

audit-ast:
	GOWORK=off go run ./cmd/graphify-ast-audit

graphify-update:
	GOWORK=off go run ./cmd/gautama-graph

audit-remediate:
	GOWORK=off go run ./cmd/graphify-doc-audit -fix

audit: graphify-update audit-docs audit-ast

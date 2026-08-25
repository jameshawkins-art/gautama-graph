.PHONY: audit-docs audit-ast audit graphify-update audit-remediate

PATH := $(HOME)/go/bin:$(HOME)/.local/bin:$(PATH)

audit-docs:
	GOWORK=off go run github.com/jameshawkins-art/gautama-graph/cmd/graphify-doc-audit -strict

audit-ast:
	GOWORK=off go run github.com/jameshawkins-art/gautama-graph/cmd/graphify-ast-audit

graphify-update:
	GOWORK=off go run github.com/jameshawkins-art/gautama-graph/cmd/gautama-graph

audit-remediate:
	GOWORK=off go run github.com/jameshawkins-art/gautama-graph/cmd/graphify-doc-audit -fix

audit: graphify-update audit-docs audit-ast

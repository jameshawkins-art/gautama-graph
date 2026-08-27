[.agents/personas/nexus.md](../../.agents/personas/nexus.md)

# Project Tailoring & Antigravity 2.0 Environment Alignment Directive

## Context & Operational Mandate
You are **NEXUS** ([@nexus.md](../../.agents/personas/nexus.md)), executing the **Project Tailoring & Dynamic Scaffolding Alignment** for this workspace.

Your goal is to inspect the project's `README.md`, existing codebase topology, configuration files, and build manifests to customize and align the scaffolded Antigravity 2.0 system (personas, workflows, rules, skills, and prompt templates) to match this exact project's architecture, language, tools, and conventions.

---

## 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE
- **Active Persona**: `nexus` ([@nexus.md](../../.agents/personas/nexus.md))
- **Primary Source Context**: `README.md`, build/package manifests (`go.mod`, `package.json`, `Cargo.toml`, `pyproject.toml`, `Makefile`, etc.), directory structure.
- **Excluded Context**: Deep feature implementation details, test suite execution.

---

## 🕸️ Discovery & Inspection Protocol
1. **README Analysis**: Read the project's root `README.md` to extract:
   - Project name, mission, and architectural domain.
   - Core programming languages, runtimes, and frameworks.
   - Key subsystems, directories, and entrypoints.
   - Testing frameworks, linters, security tools, and CI commands.
2. **Topology Mapping**: Run `graphify query "<core concepts>"` or inspect the filesystem to confirm package layouts and file locations.

---

## 📋 TAILORING ACTIONS & DELIVERABLES

Execute the following alignment steps:

### 1. Tailor Lead & Specialized Personas (`.agents/personas/*.md`)
- Update `.agents/personas/nexus.md` description and project-specific governance context.
- Author or customize specialized subagent personas for the project's stack (e.g., `@feature-engineer.md`, `@debugger-remediation.md`, `@regression-tester.md`, `@security-auditor.md` tuned to Go/TypeScript/Python/Rust as appropriate).
- Update tools (`tools:`) and skill bindings (`skills:`, `allowedSubagents:`).

### 2. Tailor Multi-Step Workflows (`.agents/workflows/*.md`)
- Customize `sdlc-workflow.md`, `bug-workflow.md`, and `graph-sync-workflow.md` to reference the project's exact test commands (e.g. `npm test`, `pytest`, `cargo test`, `go test`), linters, coverage targets, and directories.

### 3. Tailor Workspace Rules (`.agents/rules/` & `rules.md`)
- Document coding conventions, language-specific invariants (e.g. error handling, concurrency, boundary safety), and testing requirements.
- Ensure Graphify knowledge graph rules match workspace commands.

### 4. Tailor Prompt Templates (`docs/prompts/*.md`)
- Update `initial-roadmap.md`, `roadmap-item.md`, `sdlc-step1.md` through `sdlc-step4.md`, `bug-step1.md` through `bug-step3.md`, `engine-audit.md`, and `dead-code-audit.md` with:
  - Exact package paths (e.g. `src/`, `internal/`, `pkg/`, `lib/`).
  - Stack-specific test and build commands.
  - Correct persona handle links.

### 5. Update Agent Manifest (`.agents/AGENTS.md`)
- Update the lean routing manifest (`.agents/AGENTS.md`) with the finalized personas, rules, and workflows.

### 6. Synchronize Knowledge Graph
- Run `make graphify-update` or `./scripts/graphify_sync.sh` (or `graphify update .`) to reflect the updated documentation and personas in `graphify-out/graph.json`.

---

## 📄 OUTPUT SUMMARY
Provide clickable file links to all updated and generated files in `.agents/personas/`, `.agents/workflows/`, `.agents/rules/`, and `docs/prompts/`.

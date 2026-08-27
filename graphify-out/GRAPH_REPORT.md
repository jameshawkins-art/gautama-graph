# Graph Report - gautama-graph  (2026-08-26)

## Corpus Check
- 118 files · ~82,756 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1028 nodes · 1431 edges · 93 communities (82 shown, 11 thin omitted)
- Extraction: 95% EXTRACTED · 5% INFERRED · 0% AMBIGUOUS · INFERRED: 67 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `78ffc0a7`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- docs/prompts/bug-step1.md
- time.Time
- doc_auditor.go
- What You Must Do When Invoked
- docs/prompts/sdlc-step3.md
- Graphify Knowledge Graph Auditing Subsystem
- auditor/types.go
- Global Workspace Execution & Gatekeeping Rules (rules.md)
- Section 1: Core Behavioral Rules & Guardrails
- Pre-SDLC Roadmap Item Formulation Directive
- Master Engineering Roadmap Formulation Directive (Pre-SDLC)
- Requirements Specification: Deep AST Multi-Package Import & Interface Implementation Resolution
- Requirements Specification: Encapsulated Graphify Binary Manager & Single-Entrypoint Orchestrator
- Feature Roadmap Item 002: Deep AST Multi-Package Import & Interface Implementation Resolution (🟢 COMPLETED V1.2.0)
- Technical Architecture Blueprint: Encapsulated Graphify Binary Manager & Single-Entrypoint Orchestrator
- graphify reference: extra exports and benchmark
- Architecture Blueprint: Streaming AST IPC Bridge & Persistent Subprocess Daemon Pool
- Step Execution Protocol
- Lead Engine Architect & Security Code Audit Directive
- Step Execution Protocol
- Core Capabilities & Specializations
- graphify reference: query, path, explain
- Feature Roadmap Item 001: Encapsulated Graphify Binary Manager & Single-Entrypoint Orchestrator (🟢 COMPLETED V1.1.0)
- Gautama Graph: Antigravity 2.0 Multi-Agent Team Manifest
- Core Directives & Rules
- graphify reference: add a URL and watch a folder
- graphify reference: commit hook and native CLAUDE.md integration
- graphify reference: incremental update and cluster-only
- AGENTS.md
- audit_python_file
- Section 1: Core Behavioral Rules & Guardrails
- graphify reference: GitHub clone and cross-repo merge
- graphify reference: transcribe video and audio
- extraction-spec.md
- 📋 TAILORING ACTIONS & DELIVERABLES
- Detailed Item Specifications
- scripts/graphify_sync.sh
- github.com/jameshawkins-art/gautama-graph
- workerSession
- Core Directives & Rules
- Core Directives & Rules
- runner_test.go
- context.Context
- docs/prompts/sdlc-step2.md
- agent-execute.sh
- gautama-studio-execute.sh
- Requirements Specification: Streaming AST IPC Bridge & Persistent Subprocess Daemon Pool
- docs/prompts/sdlc-step1.md
- testing.T
- Technical Architecture Blueprint: Deep AST Multi-Package Import & Interface Implementation Resolution
- Core Capabilities & Specializations
- docs/prompts/sdlc-step4.md
- Feature Roadmap Item 003: Streaming AST IPC Bridge & Persistent Subprocess Daemon Pool (🟢 COMPLETED V1.3.0)
- audit_python_file
- Requirements Specification: Markdown Doc Link Auto-Remediation & Circular Cycle Detector
- Architecture Blueprint: Markdown Doc Link Auto-Remediation & Circular Cycle Detector
- 📋 TAILORING ACTIONS & DELIVERABLES
- Audit Findings Across the 5 Technical Pillars
- docs/prompts/bug-step3.md
- Feature Roadmap Item 004: Markdown Doc Link Auto-Remediation & Circular Cycle Detector (🟢 COMPLETED V1.4.0)
- templates/prompts/sdlc-step1.md
- Requirements Specification: Antigravity Environment Scaffolder & Knowledge Setup CLI
- Architecture Blueprint: Antigravity Environment Scaffolder & Knowledge Setup CLI
- Feature Roadmap Item 005: Antigravity Environment Scaffolder & Knowledge Setup CLI
- templates/rules/graphify.md
- Workflow: graphify
- templates/prompts/bug-step1.md
- templates/scripts/graphify_sync.sh
- agents_snippet.md
- templates/prompts/bug-step2.md
- templates/prompts/bug-step3.md
- Unused Symbol, Dead Code & Test-Only Function Audit Directive
- templates/prompts/sdlc-step2.md
- templates/prompts/sdlc-step3.md
- Unused Symbol, Dead Code & Test-Only Function Audit Directive
- Architecture, Code Quality & Security Audit Directive
- templates/prompts/sdlc-step4.md
- Master Engineering Roadmap Formulation Directive (Pre-SDLC)
- Pre-SDLC Roadmap Item Formulation Directive
- Step Execution Protocol
- Step Execution Protocol
- Core Directives & Rules
- Section 1: Core Behavioral Rules & Guardrails
- Graphify Knowledge Graph & Topology Rules
- docs/prompts/bug-step2.md
- 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES
- Knowledge Graph Synchronization & Integrity Workflow Guide
- JSONGraphStore
- NewDefaultPythonASTBridge
- Knowledge Graph Synchronization & Integrity Workflow Guide
- 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES
- 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES

## God Nodes (most connected - your core abstractions)
1. `Engine` - 21 edges
2. `NewDefaultEngine()` - 17 edges
3. `workerSession` - 17 edges
4. `NewEngine()` - 12 edges
5. `NewDefaultScaffolderService()` - 12 edges
6. `What You Must Do When Invoked` - 12 edges
7. `NewDefaultIPCWorkerPool()` - 10 edges
8. `PackageSymbolTable` - 10 edges
9. `AuditedEdge` - 10 edges
10. `/graphify` - 10 edges

## Surprising Connections (you probably didn't know these)
- `handleAntigravityCommand()` --calls--> `NewDefaultScaffolderService()`  [EXTRACTED]
  cmd/gautama-graph/antigravity.go → internal/scaffold/scaffolder.go
- `main()` --calls--> `NewDefaultEngine()`  [EXTRACTED]
  cmd/graphify-ast-audit/main.go → internal/auditor/engine.go
- `main()` --calls--> `NewDocGraphAuditor()`  [EXTRACTED]
  cmd/graphify-doc-audit/main.go → internal/auditor/doc_auditor.go
- `main()` --calls--> `NewStandardOrchestrator()`  [EXTRACTED]
  cmd/gautama-graph/main.go → internal/runner/orchestrator.go
- `main()` --calls--> `NewDefaultDocRemediatorService()`  [EXTRACTED]
  cmd/graphify-doc-audit/main.go → internal/auditor/doc_remediator.go

## Import Cycles
- None detected.

## Communities (93 total, 11 thin omitted)

### Community 0 - "docs/prompts/bug-step1.md"
Cohesion: 0.10
Nodes (19): 1.1 Root-Cause Isolation & Minimal Surgical Patches, 1.2 IPC Bridge & Subprocess Crash Defense (`python/ast_auditor_bridge.py`), 1.3 Concurrency & Mutex Safety, 1.4 AST Walkers & Link Graph Recursion Guardrails, 2.1 Root Cause Analysis & Remediation Proposal (`remediation_proposal.md`), 2.2 Patch Feasibility Proposal (`patch_feasibility_proposal.md`), 2.3 JSON Remediation Meta-Artifact (`remediation_meta.json`), Debugger & Remediation Skill (`gautama-mechanic`) (+11 more)

### Community 1 - "time.Time"
Cohesion: 0.17
Nodes (11): embed.FS, io/fs.FileMode, time.Time, ValidatePathBoundary(), DefaultScaffolderService, ScaffoldAction, ScaffoldActionType, ScaffolderService (+3 more)

### Community 2 - "doc_auditor.go"
Cohesion: 0.15
Nodes (17): BrokenLinkResult, DefaultDocGraphParser, DefaultDocGraphStore, DocGraphAuditor, DocGraphParser, DocGraphStore, DocNodeResult, regexp.Regexp (+9 more)

### Community 3 - "What You Must Do When Invoked"
Cohesion: 0.08
Nodes (24): For /graphify add and --watch, For /graphify query, For the commit hook and native CLAUDE.md integration, For --update and --cluster-only, /graphify, Honesty Rules, Interpreter guard for subcommands, Part A - Structural extraction for code files (+16 more)

### Community 4 - "docs/prompts/sdlc-step3.md"
Cohesion: 0.20
Nodes (9): 1. Hardened Go Engine Implementation & Production Call-Site Wiring (`@feature-engineer.md`, `@nexus.md`), 2. Table-Driven Test Suites & Regression Verification (`@regression-tester.md`), 3. SQA Certification & Deliverable Handoff (`@regression-tester.md`, `@security-auditor.md`, `@nexus.md`), CONTEXT & OBJECTIVE, 🕸️ MANDATORY GRAPHIFY DISCOVERY & POST-IMPLEMENTATION SYNC, 📄 OUTPUT FILE REQUIREMENT, 🛑 PHASE 3 & 4 EXECUTION CONSTRAINTS, 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES (+1 more)

### Community 5 - "Graphify Knowledge Graph Auditing Subsystem"
Cohesion: 0.11
Nodes (18): 1. AST Code Relationship Auditor (`cmd/graphify-ast-audit`), 2. Markdown Documentation Graph Auditor (`cmd/graphify-doc-audit`), 🏛️ Architecture & Components, CLI Flags, CLI Flags, 🛠️ CLI Utilities & Usage, Example 1: Auditing Code Relationships in Go, Example 2: Auditing Documentation Graph Integrity (+10 more)

### Community 6 - "auditor/types.go"
Cohesion: 0.05
Nodes (50): ASTParser, CircularCycle, Config, CrossPackageEvaluator, CycleReport, DefaultCrossPackageEvaluator, DefaultInterfaceResolver, DefaultPackageSymbolIndexer (+42 more)

### Community 7 - "Global Workspace Execution & Gatekeeping Rules (rules.md)"
Cohesion: 0.29
Nodes (7): 1. System Architecture & Technical Stack Boundaries, 2. Module Gatekeeper & Lifecycle Protocol, 3. Filesystem Safety & Atomic Persistence, 4. Concurrency, Stream Hygiene & Code Quality, 5. Test-Driven Development (TDD) & Production Call-Site Invariant, Global Workspace Execution & Gatekeeping Rules (rules.md), Language & Runtime Boundaries

### Community 8 - "Section 1: Core Behavioral Rules & Guardrails"
Cohesion: 0.16
Nodes (11): 1.1 Public API Export & Go Naming Conventions, 1.2 Deterministic Path & Traversal Protection, 1.3 Atomic File Persistence via Temporary Buffers, 1.4 Context Propagation & Subprocess Discipline, 1.5 Versioning & Git Tag Consistency, 1.6 Production Call-Site Invariant & TDD Alignment, 2.1 Feature Implementation Plan (`feature_implementation_plan.md`), 2.2 Patch Feasibility Proposal (`patch_feasibility_proposal.md`) (+3 more)

### Community 9 - "Pre-SDLC Roadmap Item Formulation Directive"
Cohesion: 0.29
Nodes (7): 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE, Context & Operational Mandate, 🕸️ Mandatory Graphify Knowledge Graph Discovery (Token Optimization), 📄 OUTPUT FILE REQUIREMENT, Pre-SDLC Roadmap Item Formulation Directive, 📋 REQUIRED DELIVERABLES & OUTPUT FORMAT, 🛑 STRICT PRE-SDLC EXECUTION CONSTRAINTS

### Community 10 - "Master Engineering Roadmap Formulation Directive (Pre-SDLC)"
Cohesion: 0.29
Nodes (7): 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE, Context & Operational Mandate, 🕸️ Mandatory Graphify Knowledge Graph Discovery (Token Optimization), Master Engineering Roadmap Formulation Directive (Pre-SDLC), 📄 OUTPUT FILE REQUIREMENT, 📋 REQUIRED DELIVERABLES & OUTPUT FORMAT, 🛑 STRICT PRE-SDLC EXECUTION CONSTRAINTS

### Community 11 - "Requirements Specification: Deep AST Multi-Package Import & Interface Implementation Resolution"
Cohesion: 0.14
Nodes (14): 1. Executive Summary & Strategic Scope, 2. Go Interface Specifications & Domain Model Contracts, 3.1 Zero-Trust Path Confinement (`@security-auditor.md`), 3.2 Two-Phase Atomic Persistence Protocol, 3. Filesystem Confinement, Zero-Trust Safety & Two-Phase Persistence Plan, 4.1 STRIDE Threat Model, 4.2 Standard Library & Memory Safety Guardrails, 4. Cyber Security Threat Modeling & Subprocess Safety (+6 more)

### Community 12 - "Requirements Specification: Encapsulated Graphify Binary Manager & Single-Entrypoint Orchestrator"
Cohesion: 0.13
Nodes (15): 1.1 Context & Background, 1.2 Strategic Goal & Functional Scope, 1. Executive Summary & Problem Scope, 2.1 Domain Data Models & Enums, 2.2 Go Interface Contracts (ISP Compliant), 2. Go Interface & Data Model Specifications, 3.1 Local Binary Caching Strategy, 3.2 Atomic Persistence Protocol (+7 more)

### Community 13 - "Feature Roadmap Item 002: Deep AST Multi-Package Import & Interface Implementation Resolution (🟢 COMPLETED V1.2.0)"
Cohesion: 0.29
Nodes (7): 1. Executive Summary & Strategic Objective, 2. Subsystem / Engine Component Matrix, 3. Phased Master Task Matrix, 4. Definition of Done (DoD) & Acceptance Criteria, Feature Roadmap Item 002: Deep AST Multi-Package Import & Interface Implementation Resolution (🟢 COMPLETED V1.2.0), Problem Statement, Strategic Solution & Target Architecture

### Community 14 - "Technical Architecture Blueprint: Encapsulated Graphify Binary Manager & Single-Entrypoint Orchestrator"
Cohesion: 0.17
Nodes (12): 1. Subsystem Architecture & System Topology, 2. Granular Go Interface Contracts (`internal/runner/types.go`), 3. End-to-End Execution Sequence Flow, 4.1 Subprocess Command Invariant, 4.2 Deadlock-Free Stream Capture, 4. Subprocess Lifecycle & Stream Safety Design, 5.1 TLS & SHA-256 Checksum Verification, 5.2 Zero-Trust Filesystem Confinement (+4 more)

### Community 15 - "graphify reference: extra exports and benchmark"
Cohesion: 0.22
Nodes (8): graphify reference: extra exports and benchmark, Step 6b - Wiki (only if --wiki flag), Step 7 - Neo4j export (only if --neo4j or --neo4j-push flag), Step 7a - FalkorDB export (only if --falkordb or --falkordb-push flag), Step 7b - SVG export (only if --svg flag), Step 7c - GraphML export (only if --graphml flag), Step 7d - MCP server (only if --mcp flag), Step 8 - Token reduction benchmark (only if total_words > 5000)

### Community 16 - "Architecture Blueprint: Streaming AST IPC Bridge & Persistent Subprocess Daemon Pool"
Cohesion: 0.11
Nodes (18): 1. System Architecture & High-Level Topology, 2.1 Domain Models & Framing Types, 2.2 Go Subsystem Interfaces, 2. Go Domain Interface Architecture & Contracts, 3.1 Protocol Specification, 3.2 Python Daemon Specification (`python/ast_daemon.py`), 3. Streaming NDJSON Wire Protocol & Python Worker Architecture, 4.1 Process Supervisor Details (`internal/auditor/ipc_bridge.go`) (+10 more)

### Community 17 - "Step Execution Protocol"
Cohesion: 0.25
Nodes (7): Gautama Graph SDLC Workflow Guide, Step 1: Feature Inception (`sdlc-step1`), Step 2: Implementation (`sdlc-step2`), Step 3: Regression & Fuzzing Gate (`sdlc-step3`), Step 4: Security Audit & Release Gate (`sdlc-step4`), Step-by-Step Context Isolation Matrix, Step Execution Protocol

### Community 18 - "Lead Engine Architect & Security Code Audit Directive"
Cohesion: 0.25
Nodes (8): 1. Audit Scope & Target Codepaths, 2. Audit Execution & Deliverables, 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE, Context & Operational Mandate, Lead Engine Architect & Security Code Audit Directive, 🕸️ Mandatory Graphify Knowledge Graph Scoping (Token Optimization), 📄 OUTPUT FILE REQUIREMENT, Target Codepaths & Subsystems:

### Community 19 - "Step Execution Protocol"
Cohesion: 0.29
Nodes (6): Gautama Graph Bug Remediation Workflow Guide, Step 1: Defect Capture & Isolation (`bug-step1`), Step 2: Surgical Remediation (`bug-step2`), Step 3: Regression & Security Gate (`bug-step3`), Step-by-Step Context Isolation Matrix, Step Execution Protocol

### Community 20 - "Core Capabilities & Specializations"
Cohesion: 0.25
Nodes (8): 1. Antigravity 2.0 Dynamic Scaffolding & Prompt Ops Governance, 2. Context Hygiene & Anti-Bloat Governance, 3. Multi-Agent Orchestration & SDLC Gatekeeping, 4. Gautama Graph Engine Architecture & Integrity Governance, 5. Graphify Knowledge Graph & Token Optimization Mandate, Core Capabilities & Specializations, Coupling Constraints & Operational Rules, NEXUS PROTOCOL DIRECTIVE

### Community 21 - "graphify reference: query, path, explain"
Cohesion: 0.33
Nodes (5): For /graphify explain, For /graphify path, graphify reference: query, path, explain, Step 0 — Constrained query expansion (REQUIRED before traversal), Step 1 — Traversal

### Community 22 - "Feature Roadmap Item 001: Encapsulated Graphify Binary Manager & Single-Entrypoint Orchestrator (🟢 COMPLETED V1.1.0)"
Cohesion: 0.29
Nodes (7): 1. Executive Summary & Strategic Objective, 2. Subsystem / Engine Component Matrix, 3. Phased Master Task Matrix, 4. Definition of Done (DoD) & Acceptance Criteria, Feature Roadmap Item 001: Encapsulated Graphify Binary Manager & Single-Entrypoint Orchestrator (🟢 COMPLETED V1.1.0), Problem Statement, Strategic Solution & Target Architecture

### Community 23 - "Gautama Graph: Antigravity 2.0 Multi-Agent Team Manifest"
Cohesion: 0.40
Nodes (4): Gautama Graph: Antigravity 2.0 Multi-Agent Team Manifest, Lifecycle Handoff Schema, Routing Manifest, Team Layout

### Community 24 - "Core Directives & Rules"
Cohesion: 0.29
Nodes (7): 1. Zero-Trust Path Traversal Defense, 2. Unsafe Memory & CGo Prohibition, 3. Subprocess & Command Injection Defense, 4. Dependency Scanning & API SemVer Enforcement, Core Directives & Rules, Deliverables & Meta-Artifacts, Security & Compliance Auditor Persona Specification (`gautama-gatekeeper`)

### Community 25 - "graphify reference: add a URL and watch a folder"
Cohesion: 0.50
Nodes (3): For /graphify add, For --watch, graphify reference: add a URL and watch a folder

### Community 26 - "graphify reference: commit hook and native CLAUDE.md integration"
Cohesion: 0.50
Nodes (3): For git commit hook, For native CLAUDE.md integration, graphify reference: commit hook and native CLAUDE.md integration

### Community 27 - "graphify reference: incremental update and cluster-only"
Cohesion: 0.50
Nodes (3): For --cluster-only, For --update (incremental re-extraction), graphify reference: incremental update and cluster-only

### Community 28 - "AGENTS.md"
Cohesion: 0.17
Nodes (8): Agent System Manifest, Execution Rules (`.agents/rules/`), Workflows (`.agents/workflows/`), 1. 🔴 RED STAGE: Failing Test First, 2. 🟢 GREEN STAGE: Minimal Production Code & Production Call-Site Wiring, 3. 🔵 REFACTOR & INTEGRATION AUDIT STAGE: DRY & AST Verification, Test-Driven Development (TDD) Red-Green-Refactor & Production Integration Invariant, Workflow: graphify

### Community 29 - "audit_python_file"
Cohesion: 0.67
Nodes (3): audit_python_file(), main(), Any

### Community 30 - "Section 1: Core Behavioral Rules & Guardrails"
Cohesion: 0.40
Nodes (5): 1.1 Path Traversal Defense & Zero-Trust Filesystem Confinement, 1.2 Unsafe Memory & CGo Ban, 1.3 Subprocess Isolation & Command Injection Defense, 1.4 Dependency Vulnerability Auditing & Public API Tracking, Section 1: Core Behavioral Rules & Guardrails

### Community 34 - "📋 TAILORING ACTIONS & DELIVERABLES"
Cohesion: 0.15
Nodes (12): 1. Tailor Lead & Specialized Personas (`.agents/personas/*.md`), 2. Tailor Multi-Step Workflows (`.agents/workflows/*.md`), 3. Tailor Workspace Rules (`.agents/rules/` & `rules.md`), 4. Tailor Prompt Templates (`docs/prompts/*.md`), 5. Update Agent Manifest (`.agents/AGENTS.md`), 6. Synchronize Knowledge Graph, 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE, Context & Operational Mandate (+4 more)

### Community 35 - "Detailed Item Specifications"
Cohesion: 0.18
Nodes (11): Architectural Subsystems Overview, Detailed Item Specifications, Executive Summary & Strategic Mission, Gautama Graph Master Product & Architecture Roadmap, Item 001: Encapsulated Graphify Binary Manager & Single-Entrypoint Orchestrator, Item 002: Deep AST Multi-Package Import & Interface Implementation Resolution, Item 003: Streaming AST IPC Bridge & Persistent Subprocess Daemon Pool, Item 004: Markdown Doc Link Auto-Remediation & Circular Cycle Detector (+3 more)

### Community 38 - "workerSession"
Cohesion: 0.18
Nodes (10): DefaultIPCWorkerPool, PoolStats, workerSession, WorkerState, WorkerStats, bufio.Scanner, io.WriteCloser, os/exec.Cmd (+2 more)

### Community 39 - "Core Directives & Rules"
Cohesion: 0.29
Nodes (7): 1. Public API Design & Export Standards, 2. Filesystem Safety & Atomic Persistence, 3. Context & Subprocess Hygiene, 4. TDD & Production Call-Site Invariant, Core Directives & Rules, Deliverables & Meta-Artifacts, Feature Engineer Persona Specification (`gautama-builder`)

### Community 40 - "Core Directives & Rules"
Cohesion: 0.29
Nodes (7): 1. Test Isolation & Hermetic Environments, 2. Table-Driven Tests & Edge Cases, 3. Concurrency, Race Detection & Coverage Gates, 4. TDD Verification & Production Call-Site Invariant, Core Directives & Rules, Deliverables & Meta-Artifacts, Regression & Test Automation Persona Specification (`gautama-guard`)

### Community 41 - "runner_test.go"
Cohesion: 0.06
Nodes (41): handleAntigravityCommand(), main(), net/http.Client, time.Duration, DocGraphAuditorService, ASTGraphAuditorService, NewDefaultReleaseDownloader(), NewDefaultBinaryManager() (+33 more)

### Community 42 - "context.Context"
Cohesion: 0.33
Nodes (4): AuditedEdge, CandidateEdge, IPCResponse, context.Context

### Community 43 - "docs/prompts/sdlc-step2.md"
Cohesion: 0.25
Nodes (6): Feature Engineer Skill (`gautama-builder`), CONTEXT & OBJECTIVE, 🕸️ MANDATORY GRAPHIFY KNOWLEDGE GRAPH MAPPING (TOKEN OPTIMIZATION), 📄 OUTPUT FILE REQUIREMENT, 🛑 PHASE 2 EXECUTION CONSTRAINTS, 🔒 STEP 2 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE

### Community 46 - "Requirements Specification: Streaming AST IPC Bridge & Persistent Subprocess Daemon Pool"
Cohesion: 0.11
Nodes (19): 1.1 Context & Problem Statement, 1.2 Target Vision, 1. Executive Summary & Problem Scope, 2.1 IPC Protocol & Data Contracts, 2.2 IPC Worker & Pool Interfaces, 2. Go Interface Contracts & Domain Models, 3.1 Framing Specification, 3.2 Python Daemon Specification (`python/ast_daemon.py`) (+11 more)

### Community 47 - "docs/prompts/sdlc-step1.md"
Cohesion: 0.18
Nodes (10): 1. Requirements Scope & Go Interface Specifications (`@nexus.md`, `@feature-engineer.md`), 2. Filesystem Confinement & Two-Phase Persistence Plan (`@feature-engineer.md`, `@security-auditor.md`), 3. Cyber Security Threat Modeling & Subprocess Safety (`@security-auditor.md`), 4. Definition of Done (DoD) & Acceptance Criteria (`@nexus.md`, `@feature-engineer.md`), CONTEXT & OBJECTIVE, 🕸️ MANDATORY GRAPHIFY KNOWLEDGE GRAPH DISCOVERY (TOKEN OPTIMIZATION), 📄 OUTPUT FILE REQUIREMENT, 🛑 PHASE 1 EXECUTION CONSTRAINTS (+2 more)

### Community 48 - "testing.T"
Cohesion: 0.07
Nodes (46): CycleDetector, DefaultASTParser, DefaultDocRemediatorService, DocRemediationPlan, HeadingAnchorTable, TarjanSCCDetector, TestGautamaGraphCLI_AntigravityInstall_DryRun(), TestGautamaGraphCLI_AntigravityInstall_Execution() (+38 more)

### Community 49 - "Technical Architecture Blueprint: Deep AST Multi-Package Import & Interface Implementation Resolution"
Cohesion: 0.12
Nodes (17): 1. Executive Architecture Summary & System Topology, 2.1 Domain Data Structures & Contracts, 2.2 Concrete Subsystem Components, 2. Detailed Go Interface & Struct Architecture, 3.1 Workspace Package Discovery & Compilation (`indexer.go`), 3.2 Cross-Package Selector Call Evaluation (`cross_package.go`), 3.3 Implicit Interface Satisfaction Resolution (`resolver.go`), 3. Algorithmic Multi-Package Indexing & Type-Checking Engine (+9 more)

### Community 50 - "Core Capabilities & Specializations"
Cohesion: 0.15
Nodes (11): 1. Antigravity 2.0 Dynamic Scaffolding & Prompt Ops Governance, 2. Context Hygiene & Anti-Bloat Governance, 3. Multi-Agent Orchestration & SDLC Gatekeeping, 4. Knowledge Graph Architecture & Graphify Governance, Core Capabilities & Specializations, Coupling Constraints & Operational Rules, NEXUS PROTOCOL DIRECTIVE, 1. 🔴 RED STAGE: Failing Test First (+3 more)

### Community 51 - "docs/prompts/sdlc-step4.md"
Cohesion: 0.22
Nodes (8): 1. Security Verification & Dependency Audit (`@security-auditor.md`, `@nexus.md`), 2. Compilation & Master Knowledge Graph Sync (`@nexus.md`, `@regression-tester.md`), 3. Release Matrix Update & Production Synchronization (`@nexus.md`), CONTEXT & OBJECTIVE, 📄 OUTPUT FILE REQUIREMENT, 🛑 PHASE 5 & 6 EXECUTION AUTHORITY, 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES, 🔒 STEP 4 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE

### Community 53 - "Feature Roadmap Item 003: Streaming AST IPC Bridge & Persistent Subprocess Daemon Pool (🟢 COMPLETED V1.3.0)"
Cohesion: 0.29
Nodes (7): 1. Executive Summary & Strategic Objective, 2. Subsystem / Engine Component Matrix, 3. Phased Master Task Matrix, 4. Definition of Done (DoD), Feature Roadmap Item 003: Streaming AST IPC Bridge & Persistent Subprocess Daemon Pool (🟢 COMPLETED V1.3.0), Problem Statement, Strategic Solution & Target Architecture

### Community 54 - "audit_python_file"
Cohesion: 0.67
Nodes (3): audit_python_file(), main(), Any

### Community 55 - "Requirements Specification: Markdown Doc Link Auto-Remediation & Circular Cycle Detector"
Cohesion: 0.10
Nodes (21): 1.1 Context & Problem Statement, 1.2 Target Vision, 1. Executive Summary & Problem Scope, 2.1 Domain Models & Data Structures, 2.2 Subsystem Interfaces, 2. Go Interface Contracts & Domain Models, 3.1 Canonical Relative Path Calculation, 3.2 Fuzzy Basename & Slug Matching (+13 more)

### Community 56 - "Architecture Blueprint: Markdown Doc Link Auto-Remediation & Circular Cycle Detector"
Cohesion: 0.11
Nodes (18): 1. System Architecture & High-Level Topology, 2.1 Domain Data Structures, 2.2 Subsystem Interfaces, 2. Go Interface Architecture & Domain Contracts, 3.1 Canonical Relative Path Normalization, 3.2 Fuzzy Basename & Tree Distance Matching, 3.3 GitHub-Flavored Markdown (GFM) Anchor Extraction, 3.4 Tarjan's Strongly Connected Components (SCC) Cycle Detector (+10 more)

### Community 57 - "📋 TAILORING ACTIONS & DELIVERABLES"
Cohesion: 0.15
Nodes (12): 1. Tailor Lead & Specialized Personas (`.agents/personas/*.md`), 2. Tailor Multi-Step Workflows (`.agents/workflows/*.md`), 3. Tailor Workspace Rules (`.agents/rules/` & `rules.md`), 4. Tailor Prompt Templates (`docs/prompts/*.md`), 5. Update Agent Manifest (`.agents/AGENTS.md`), 6. Synchronize Knowledge Graph, 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE, Context & Operational Mandate (+4 more)

### Community 58 - "Audit Findings Across the 5 Technical Pillars"
Cohesion: 0.20
Nodes (10): 1. Pillar 1: Go Interface Abstractions & Export Standards, 2. Pillar 2: Zero-Trust Path Confinement & Traversal Defense, 3. Pillar 3: Atomic Persistence & Concurrency Hygiene, 4. Pillar 4: Subprocess Lifecycle & Stream Safety, 5. Pillar 5: Test Suite Integrity & Vulnerability Assessment, Audit Findings Across the 5 Technical Pillars, Executive Summary, Lead Engine Architect & Security Code Audit Report (+2 more)

### Community 59 - "docs/prompts/bug-step3.md"
Cohesion: 0.33
Nodes (5): CONTEXT & OBJECTIVE, 🕸️ MANDATORY GRAPHIFY DISCOVERY & POST-FIX GRAPH SYNC, 📄 OUTPUT FILE REQUIREMENT, 🛑 PHASE B3 & B4 EXECUTION CONSTRAINTS, 🔒 STEP B3 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE

### Community 60 - "Feature Roadmap Item 004: Markdown Doc Link Auto-Remediation & Circular Cycle Detector (🟢 COMPLETED V1.4.0)"
Cohesion: 0.29
Nodes (7): 1. Executive Summary & Strategic Objective, 2. Subsystem / Engine Component Matrix, 3. Phased Master Task Matrix, 4. Definition of Done (DoD), Feature Roadmap Item 004: Markdown Doc Link Auto-Remediation & Circular Cycle Detector (🟢 COMPLETED V1.4.0), Problem Statement, Strategic Solution & Target Architecture

### Community 61 - "templates/prompts/sdlc-step1.md"
Cohesion: 0.18
Nodes (10): 1. Requirements Scope & Interface Specifications, 2. Filesystem Confinement & Atomic Persistence Plan, 3. Cyber Security Threat Modeling & Dependency Safety, 4. Definition of Done (DoD) & Acceptance Criteria, CONTEXT & OBJECTIVE, 🕸️ MANDATORY GRAPHIFY KNOWLEDGE GRAPH DISCOVERY (TOKEN OPTIMIZATION), 📄 OUTPUT FILE REQUIREMENT, 🛑 PHASE 1 EXECUTION CONSTRAINTS (+2 more)

### Community 62 - "Requirements Specification: Antigravity Environment Scaffolder & Knowledge Setup CLI"
Cohesion: 0.10
Nodes (20): 1.1 Context & Problem Statement, 1.2 Target Vision, 1. Executive Summary & Problem Scope, 2.1 Domain Models & Data Structures, 2.2 Subsystem Interfaces, 2. Go Interface Contracts & Domain Models, 3.1 Embedded Template Rendering, 3.2 Non-Destructive File Merging (+12 more)

### Community 63 - "Architecture Blueprint: Antigravity Environment Scaffolder & Knowledge Setup CLI"
Cohesion: 0.12
Nodes (17): 1. System Architecture & High-Level Topology, 2.1 Domain Data Structures, 2.2 Subsystem Interfaces, 2. Go Interface Architecture & Domain Contracts, 3.1 Embedded Template Filesystem Layout, 3.2 Plan Formulation & Non-Destructive Scaffolding, 3.3 Atomic Two-Phase File Commit, 3. Core Algorithms & Implementation Details (+9 more)

### Community 64 - "Feature Roadmap Item 005: Antigravity Environment Scaffolder & Knowledge Setup CLI"
Cohesion: 0.29
Nodes (7): 1. Executive Summary & Strategic Objective, 2. Subsystem / Engine Component Matrix, 3. Phased Master Task Matrix, 4. Definition of Done (DoD), Feature Roadmap Item 005: Antigravity Environment Scaffolder & Knowledge Setup CLI, Problem Statement, Strategic Solution & Target Architecture

### Community 67 - "templates/prompts/bug-step1.md"
Cohesion: 0.20
Nodes (9): 1. Defect Classification & Impact Analysis, 2. Deterministic Minimal Reproduction Case, 3. Bug Specification Document (`docs/bugs/bug-<description>-<id>.md`), CONTEXT & OBJECTIVE, 🕸️ MANDATORY GRAPHIFY KNOWLEDGE GRAPH DISCOVERY (TOKEN OPTIMIZATION), 📄 OUTPUT FILE REQUIREMENT, 🛑 PHASE B1 EXECUTION CONSTRAINTS, 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES (+1 more)

### Community 70 - "templates/prompts/bug-step2.md"
Cohesion: 0.20
Nodes (9): 1. Root Cause Analysis (RCA), 2. Technical Remediation Blueprint & Surgical Diff Preview, 3. Security & Side-Effect Assessment, CONTEXT & OBJECTIVE, 🕸️ MANDATORY GRAPHIFY KNOWLEDGE GRAPH MAPPING (TOKEN OPTIMIZATION), 📄 OUTPUT FILE REQUIREMENT, 🛑 PHASE B2 EXECUTION CONSTRAINTS, 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES (+1 more)

### Community 71 - "templates/prompts/bug-step3.md"
Cohesion: 0.20
Nodes (9): 1. Targeted Code Remediation & Refactor, 2. Regression & SQA Verification, 3. Security Check & Defect Closure, CONTEXT & OBJECTIVE, 🕸️ MANDATORY GRAPHIFY DISCOVERY & POST-FIX GRAPH SYNC, 📄 OUTPUT FILE REQUIREMENT, 🛑 PHASE B3 & B4 EXECUTION CONSTRAINTS, 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES (+1 more)

### Community 72 - "Unused Symbol, Dead Code & Test-Only Function Audit Directive"
Cohesion: 0.20
Nodes (9): 1. Static AST Symbol & Caller Graph Extraction, 2. Defect Classification Taxonomy, 📋 Audit Scope & Methodology, 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE, Context & Operational Mandate, 🕸️ Mandatory Graphify Knowledge Graph Scoping, 📄 OUTPUT FILE REQUIREMENT, Unused Symbol, Dead Code & Test-Only Function Audit Directive (+1 more)

### Community 73 - "templates/prompts/sdlc-step2.md"
Cohesion: 0.20
Nodes (9): 1. Engine & Interface Architecture Blueprint, 2. Subprocess & API Contracts, 3. Cyber Security Architecture & Hardening Strategy, CONTEXT & OBJECTIVE, 🕸️ MANDATORY GRAPHIFY KNOWLEDGE GRAPH MAPPING (TOKEN OPTIMIZATION), 📄 OUTPUT FILE REQUIREMENT, 🛑 PHASE 2 EXECUTION CONSTRAINTS, 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES (+1 more)

### Community 74 - "templates/prompts/sdlc-step3.md"
Cohesion: 0.20
Nodes (9): 1. Hardened Implementation & Production Call-Site Wiring, 2. Table-Driven Test Suites & Regression Verification, 3. SQA Certification & Deliverable Handoff, CONTEXT & OBJECTIVE, 🕸️ MANDATORY GRAPHIFY DISCOVERY & POST-IMPLEMENTATION SYNC, 📄 OUTPUT FILE REQUIREMENT, 🛑 PHASE 3 & 4 EXECUTION CONSTRAINTS, 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES (+1 more)

### Community 75 - "Unused Symbol, Dead Code & Test-Only Function Audit Directive"
Cohesion: 0.22
Nodes (9): 1. Static AST Symbol & Caller Graph Extraction, 2. Defect Classification Taxonomy, 📋 Audit Scope & Methodology, 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE, Context & Operational Mandate, 🕸️ Mandatory Graphify Knowledge Graph Scoping, 📄 OUTPUT FILE REQUIREMENT, Unused Symbol, Dead Code & Test-Only Function Audit Directive (+1 more)

### Community 76 - "Architecture, Code Quality & Security Audit Directive"
Cohesion: 0.22
Nodes (8): 1. Audit Scope & Target Codepaths, 2. Audit Execution & Deliverables, Architecture, Code Quality & Security Audit Directive, 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE, Context & Operational Mandate, 🕸️ Mandatory Graphify Knowledge Graph Scoping (Token Optimization), 📄 OUTPUT FILE REQUIREMENT, Target Codepaths & Subsystems:

### Community 77 - "templates/prompts/sdlc-step4.md"
Cohesion: 0.22
Nodes (8): 1. Security Verification & Dependency Audit, 2. Compilation & Master Knowledge Graph Sync, 3. Release Matrix Update & Production Synchronization (`@nexus.md`), CONTEXT & OBJECTIVE, 📄 OUTPUT FILE REQUIREMENT, 🛑 PHASE 5 & 6 EXECUTION AUTHORITY, 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES, 🔒 STEP 4 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE

### Community 78 - "Master Engineering Roadmap Formulation Directive (Pre-SDLC)"
Cohesion: 0.25
Nodes (7): 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE, Context & Operational Mandate, 🕸️ Mandatory Graphify Knowledge Graph Discovery (Token Optimization), Master Engineering Roadmap Formulation Directive (Pre-SDLC), 📄 OUTPUT FILE REQUIREMENT, 📋 REQUIRED DELIVERABLES & OUTPUT FORMAT, 🛑 STRICT PRE-SDLC EXECUTION CONSTRAINTS

### Community 79 - "Pre-SDLC Roadmap Item Formulation Directive"
Cohesion: 0.25
Nodes (7): 🔒 CONTEXT ISOLATION & PROGRESSIVE DISCLOSURE, Context & Operational Mandate, 🕸️ Mandatory Graphify Knowledge Graph Discovery (Token Optimization), 📄 OUTPUT FILE REQUIREMENT, Pre-SDLC Roadmap Item Formulation Directive, 📋 REQUIRED DELIVERABLES & OUTPUT FORMAT, 🛑 STRICT PRE-SDLC EXECUTION CONSTRAINTS

### Community 80 - "Step Execution Protocol"
Cohesion: 0.25
Nodes (7): SDLC Feature Engineering Workflow Guide, Step 1: Feature Inception (`sdlc-step1`), Step 2: Implementation (`sdlc-step2`), Step 3: Regression & Fuzzing Gate (`sdlc-step3`), Step 4: Security Audit & Release Gate (`sdlc-step4`), Step-by-Step Context Isolation Matrix, Step Execution Protocol

### Community 81 - "Step Execution Protocol"
Cohesion: 0.29
Nodes (6): Bug Remediation Workflow Guide, Step 1: Defect Capture & Isolation (`bug-step1`), Step 2: Surgical Remediation (`bug-step2`), Step 3: Regression & Security Gate (`bug-step3`), Step-by-Step Context Isolation Matrix, Step Execution Protocol

### Community 82 - "Core Directives & Rules"
Cohesion: 0.33
Nodes (6): 1. Minimal Surgical Patches & Repro-First, 2. IPC Bridge & Subprocess Resilience, 3. Concurrency & Mutex Discipline, Core Directives & Rules, Debugger & Remediation Persona Specification (`gautama-mechanic`), Deliverables & Meta-Artifacts

### Community 83 - "Section 1: Core Behavioral Rules & Guardrails"
Cohesion: 0.17
Nodes (11): 1.1 Test Isolation & Hermetic Filesystem Usage, 1.2 Table-Driven Test Patterns & Boundary Testing, 1.3 Race Detection & Concurrency Verification, 1.4 Coverage Gate & Quality Thresholds, 1.5 TDD Red-Green-Refactor Protocol & Production Call-Site Auditing, 2.1 Unit Execution Report (`unit_execution_report.md`), 2.2 Regression Matrix Checklist (`regression_matrix.md`), 2.3 Test Verification Meta-Artifact (`test_verification_meta.json`) (+3 more)

### Community 84 - "Graphify Knowledge Graph & Topology Rules"
Cohesion: 0.40
Nodes (5): 1. Graph Discovery First (Token Optimization), 2. Post-Modification Update & AST Pruning, 3. Documentation Link Standards, 4. Full Pipeline Synchronization, Graphify Knowledge Graph & Topology Rules

### Community 85 - "docs/prompts/bug-step2.md"
Cohesion: 0.17
Nodes (10): 2.1 Security Audit Report (`security_audit_report.md`), 2.2 Security Compliance Checklist (`compliance_checklist.md`), 2.3 JSON Security Verification Meta-Artifact (`security_verification_meta.json`), Section 2: Expected Antigravity Artifact Deliverables, Security & Compliance Auditor Skill (`gautama-gatekeeper`), CONTEXT & OBJECTIVE, 🕸️ MANDATORY GRAPHIFY KNOWLEDGE GRAPH MAPPING (TOKEN OPTIMIZATION), 📄 OUTPUT FILE REQUIREMENT (+2 more)

### Community 86 - "📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES"
Cohesion: 0.50
Nodes (4): 1. Root Cause Analysis (RCA) (`@debugger-remediation.md`, `@nexus.md`), 2. Technical Remediation Blueprint & Surgical Diff Preview (`@debugger-remediation.md`), 3. Security & Side-Effect Assessment (`@security-auditor.md`), 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES

### Community 87 - "Knowledge Graph Synchronization & Integrity Workflow Guide"
Cohesion: 0.50
Nodes (3): Knowledge Graph Synchronization & Integrity Workflow Guide, Master Script Invocation, Pipeline Execution Stages

### Community 88 - "JSONGraphStore"
Cohesion: 0.40
Nodes (4): GraphData, JSONGraphStore, sync.Mutex, NewJSONGraphStore()

### Community 89 - "NewDefaultPythonASTBridge"
Cohesion: 0.50
Nodes (4): DefaultPythonASTBridge, pyBridgePayload, pyBridgeResponse, NewDefaultPythonASTBridge()

### Community 90 - "Knowledge Graph Synchronization & Integrity Workflow Guide"
Cohesion: 0.50
Nodes (3): Knowledge Graph Synchronization & Integrity Workflow Guide, Master Script Invocation, Pipeline Execution Stages

### Community 91 - "📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES"
Cohesion: 0.50
Nodes (4): 1. Targeted Code Remediation & Refactor (`@debugger-remediation.md`), 2. Regression & SQA Verification (`@regression-tester.md`), 3. Security Check & Defect Closure (`@security-auditor.md`, `@nexus.md`), 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES

### Community 92 - "📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES"
Cohesion: 0.50
Nodes (4): 1. Go Engine & Interface Architecture Blueprint (`@feature-engineer.md`, `@nexus.md`), 2. Python IPC Bridge & Subprocess Contract (`@feature-engineer.md`, `@security-auditor.md`), 3. Cyber Security Architecture & Hardening Strategy (`@security-auditor.md`), 📋 REQUIRED DELIVERABLES & PERSONA RESPONSIBILITIES

## Knowledge Gaps
- **514 isolated node(s):** `github.com/jameshawkins-art/gautama-graph`, `IPCSession`, `DocRemediatorService`, `OrchestratorService`, `graphify_sync.sh script` (+509 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **11 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Architecture Blueprint: Antigravity Environment Scaffolder & Knowledge Setup CLI` connect `Architecture Blueprint: Antigravity Environment Scaffolder & Knowledge Setup CLI` to `.agents/personas/nexus.md`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Why does `Requirements Specification: Streaming AST IPC Bridge & Persistent Subprocess Daemon Pool` connect `Requirements Specification: Streaming AST IPC Bridge & Persistent Subprocess Daemon Pool` to `.agents/personas/nexus.md`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **Why does `Architecture Blueprint: Markdown Doc Link Auto-Remediation & Circular Cycle Detector` connect `Architecture Blueprint: Markdown Doc Link Auto-Remediation & Circular Cycle Detector` to `.agents/personas/nexus.md`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Are the 11 inferred relationships involving `NewDefaultEngine()` (e.g. with `NewDefaultCrossPackageEvaluator()` and `NewDefaultSelectorEvaluator()`) actually correct?**
  _`NewDefaultEngine()` has 11 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/jameshawkins-art/gautama-graph`, `IPCSession`, `DocRemediatorService` to the rest of the system?**
  _514 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `docs/prompts/bug-step1.md` be split into smaller, more focused modules?**
  _Cohesion score 0.09523809523809523 - nodes in this community are weakly interconnected._
- **Should `What You Must Do When Invoked` be split into smaller, more focused modules?**
  _Cohesion score 0.08 - nodes in this community are weakly interconnected._
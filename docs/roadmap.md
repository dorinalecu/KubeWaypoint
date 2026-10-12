# Development Roadmap

This roadmap describes direction, not authorization to implement all milestones.
Implement only the current approved increment. Each phase must provide a useful,
verifiable outcome before adding the next source of complexity.

The project serves engineers working with AKS using public information. Internal
Microsoft workflows, customer environments, upgrade readiness, and non-public
material from other projects are outside its scope.

## Phase 0: Engineering Foundation

Define product boundaries, evidence semantics, runtime responsibilities, an
architecture decision, development guidance, and the incremental Go plan.

**Exit:** The [architecture](architecture.md), [ADR](adr/0001-go-first-deterministic-stateless-mvp.md),
[development guide](development.md), and [Phase 1 plan](phase-1-plan.md) agree.
There is no application or future-phase scaffolding.

## Phase 1: Go MVP

Build one local, stateless, deterministic API with versioned JSON, a small AKS
taxonomy, bounded validation, explicit failures, safe logs, correlation IDs,
timeouts, shutdown, and synthetic tests.

**Learning focus:** Go types, errors, composition, HTTP, context, and tests.
**Exit:** All Phase 1 acceptance criteria pass without cloud or AI dependencies.

## Phase 2: Python Analysis Service

Move analysis to Python while Go retains public API behavior and orchestration.
Define a versioned internal contract and validate Python output. Keep analysis
deterministic initially; test unavailable, slow, and incompatible responses.
Avoid silently retaining two rule engines.

**Learning focus:** Python services, contract testing, distributed failure modes.
**Exit:** Both services work independently and end-to-end; failures remain visible.
Public behavior remains compatible or is explicitly versioned.

## Phase 3: Containerization

Package the working services with reproducible builds, non-root execution,
external configuration, health checks, and a documented local startup flow.
Keep secrets out of images. Do not introduce Kubernetes solely to run locally.

**Learning focus:** Build/runtime separation, networking, process lifecycle.
**Exit:** A clean checkout can run both services using the documented container
flow, including correct shutdown.

## Phase 4: Optional LLM Integration

Introduce a replaceable Python provider boundary with explicit enablement,
validated structured output, deadlines, cost bounds, and no-AI operation.
Treat input as untrusted data and model output as unverified analysis.
Evaluate against the deterministic baseline using synthetic scenarios.

**Learning focus:** Probabilistic dependencies, output validation, evaluation.
**Exit:** AI failures and limitations are visible, data transfer is intentional,
and engineer review remains mandatory. Select the runtime provider separately
from the coding assistant.

## Phase 5: Public AKS Knowledge and RAG

Curate licensed public sources with provenance and freshness information.
Retrieve relevant material and attach references that actually support the
recommendations. Choose retrieval tooling from measured needs; a vector database
is not a prerequisite.

**Learning focus:** Retrieval relevance, grounding, source quality.
**Exit:** Recommendations can be traced to public evidence, and unsupported
claims and retrieval limitations are evaluated.

## Phase 6: Observability and Evaluation

Expand the tests and logging established in Phase 1 with latency/error metrics,
useful tracing, representative synthetic scenarios, regression evaluation, and
dependency failure testing. Measure both software reliability and analysis
quality without logging raw case content.

**Learning focus:** Operational visibility, measurable quality, regression control.
**Exit:** Changes can be compared against a documented baseline, including
ambiguity, unsupported claims, and failures. Classification confidence is not
represented as a root-cause probability.

## Phase 7: Azure Infrastructure

Choose hosting in a separate ADR, weighing AKS learning goals against cost and
operational complexity. Add repeatable infrastructure, identity, secret handling,
network controls, access protection, budgets, and cleanup instructions.

**Learning focus:** Infrastructure as Code, identity, networking, operations.
**Exit:** An explicitly authorized environment can be deployed, validated, and
removed reproducibly. An unauthenticated local MVP is not exposed publicly.

## Phase 8: CI/CD and GitOps

Automate established checks, build and version images, and introduce controlled
deployment, verification, promotion, and rollback. Use GitOps if the chosen
hosting architecture justifies it. Tests exist before this milestone; pipelines
automate them rather than introducing correctness for the first time.

**Learning focus:** Auditable delivery, release management, recovery.
**Exit:** Delivery and rollback are documented, repeatable, and verified.

## Phase 9: Controlled Agentic Case Intake

Explore an iterative workflow that asks questions, incorporates reported
answers, retrieves public evidence, and revises an investigation plan.
Define allowed steps, execution budgets, stopping conditions, and human review.
Any durable workflow state requires a new persistence decision.

**Learning focus:** Workflow orchestration, state, bounded execution, auditability.
**Exit:** Steps and evidence transitions are traceable without autonomous cluster
access, diagnostic execution, or customer communication.

## Portfolio checkpoints

Phases 0-1 deliver an explainable Go application. Phases 2-3 deliver a
reproducible Go/Python system. Phases 4-6 add optional AI and measurable analysis
quality. Phases 7-8 add cloud operation and delivery. Phase 9 adds a controlled
multi-step workflow.

Later milestones may change based on demonstrated needs. Record significant
decisions instead of building speculative extensibility now.

# ADR 0001: Start with a Go-First, Deterministic, Stateless MVP

- Status: Accepted for Phase 1
- Date: 2026-10-12
- Scope: Initial runtime architecture and analysis approach

## Context

KubeWaypoint is a public portfolio project combining practical AKS investigation
knowledge with software engineering learning. No application exists yet.
The first useful outcome is consistent intake analysis, not autonomous
troubleshooting or a production support platform.

Go and Python are both intended learning goals. Introducing two runtimes, network
failures, AI output uncertainty, and deployment infrastructure simultaneously
would make errors harder to isolate and weaken the educational progression.

The product must distinguish reported information from verified evidence, expose
its limitations, and require engineer review. Development must use original
content, synthetic cases, and public references, with no internal systems or
customer data.

## Decision

Build Phase 1 as one local Go process using the standard library where practical.

- Expose a versioned analysis API and health/readiness endpoints.
- Keep domain analysis separate from HTTP handlers, but in the same process.
- Use transparent deterministic rules for node registration and DNS.
- Return explicit unknown and insufficient-information outcomes.
- Preserve reported evidence semantics and require engineer review.
- Keep the runtime stateless, without a database or case history.
- Bound requests and define explicit errors, timeouts, logs, and shutdown.
- Establish synthetic unit, HTTP, and behavioral tests before adding services.
- Bind to loopback by default; do not expose an unauthenticated MVP publicly.

Do not create Python, provider, retrieval, container, cloud, or CI scaffolding.
In Phase 2, move analysis to Python behind a versioned internal contract while
Go retains the public API and orchestration. Introduce LLMs only in a later,
explicitly enabled milestone.

## Alternatives considered

| Alternative | Why not start here? |
| --- | --- |
| Go API plus Python service immediately | Adds distribution concerns before analysis behavior and contracts are established |
| Python-only API | Valid technically, but does not meet the intended Go service-learning milestone |
| LLM-first analysis | Makes baseline behavior, unsupported claims, costs, and data transfer harder to control |
| Database-backed application | Adds storage and data-retention responsibilities without a current user need |
| Kubernetes-first delivery | Adds operational work before there is a useful application to deploy |

## Consequences

### Benefits

- Local execution requires no cloud account or external AI calls.
- Rules and failure behavior can be tested and explained.
- Transport, domain modeling, and Go lifecycle concerns can be learned together.
- The later Python implementation has an observable behavioral baseline.
- There is no application-owned store of case content to operate.

### Costs and limitations

- Deterministic extraction and summaries will be conservative and limited.
- The initial taxonomy will not cover arbitrary AKS problems.
- Rules cannot establish root cause or independently verify reports.
- Moving analysis to Python involves migration work, not a free replacement.
- Statelessness does not provide history, durable review, or secure memory erasure.
- Local-only operation does not solve authentication or production deployment.

## Revisit triggers

Revisit this decision when Phase 2 starts, when a demonstrated workflow needs
durable state, or before any remote deployment. Record the change in a new ADR;
do not silently reinterpret this decision to authorize later infrastructure.

## Verification

Phase 1 is complete only when a local Go API processes synthetic cases, follows
the [documented contract](../architecture.md), exposes visible failures, and
passes the [implementation-plan acceptance criteria](../phase-1-plan.md).

# Phase 1: Incremental Implementation Plan

## Objective

Deliver a locally runnable Go API that turns a synthetic AKS problem description
into a deterministic, engineer-reviewed investigation starting point.

Follow the [architecture and planned contract](architecture.md) and
[ADR 0001](adr/0001-go-first-deterministic-stateless-mvp.md). This document is a
plan, not a claim that implementation or commands already exist.

All increments are sequential: each uses the preceding behavior and tests.
Explain the increment before edits. Delegate code and tests to Sonnet 5.5;
the orchestrator owns decisions, acceptance criteria, validation, and review.
Do not commit or push without a relevant explicit request.

Current milestone: increment 1. The agreed
[analysis contract](analysis-contract.md) fixes nested shapes, synthetic
expectations, and the 0 / 0.5 / 0.8 support rubric. Product input/output is English
only. Go 1.27.2 is the minimum version and initial validation toolchain.

## Increment 1: Domain contract and synthetic expectations

Before implementing rules, settle the nested response shapes and exact scoring
with the engineer using a small table of synthetic scenarios.

- Include positive node-registration and DNS descriptions.
- Include empty/whitespace input, a vague report, an unsupported actionable
  symptom, and mixed node-registration/DNS evidence.
- Include negated and tentative statements and an isolated generic keyword.
- Specify expected category, score, evidence, gaps, questions, and checks.
- Give rules stable IDs and document their match conditions and support tiers.
- Resolve any active multi-category evidence to unknown, without arbitrary
  category ordering or comparison of tentative versus explicit tiers.
- Confirm that zero-confidence unresolved outcomes do not claim root cause.

Then use the selected Go version, create the root module
`github.com/dorinalecu/KubeWaypoint`, and add domain types and the first focused
tests. Create only packages used by this increment.

**Acceptance:** Contracts and scenario expectations are explicit, score bounds
and ambiguity behavior are documented, the module builds, and initial tests
pass. No HTTP, Python, container, or provider scaffolding is added.

## Increment 2: In-process deterministic analyzer

Implement analysis in `internal/analysis`.

- Classify with the agreed rule table and explicit unresolved outcomes. The
  rubric uses the highest active tier within one category, not summed weights.
- Produce a conservative summary and supported reported-fact extraction.
- Associate classification evidence with rule IDs and input excerpts.
- Identify gaps and generate focused, deduplicated questions and initial checks.
- Respect cancellation and separate unsupported input from invalid input.
- Keep ordering stable and avoid mutating caller-owned input.
- Do not manufacture verified facts or conclusions.

**Acceptance:** Table-driven tests cover all increment-1 scenarios, negation,
ambiguity, confidence boundaries, empty collections, and deterministic output.
Tests verify that tentative statements are not strengthened into assertions.

## Increment 3: Versioned HTTP API

Add `internal/httpapi` and the minimal runnable `cmd/kubewaypoint` entry point.

- Wire `POST /api/v1/analyze`, `GET /healthz`, and `GET /readyz`.
- Enforce the 32 KiB body limit and all field limits in the architecture.
- Reject malformed JSON, invalid UTF-8, duplicate keys, unknown fields, wrong
  types, explicit null fields, and trailing JSON values.
- Accept JSON media-type parameters; reject unsupported content types.
- Emit stable success and error shapes, non-null collection arrays, and the
  review flag. Keep domain analysis independent of HTTP.
- Implement correlation ID generation/validation and safe error responses.
- Handle missing routes and unsupported methods consistently.
- Listen on loopback by default; expose no remote access implicitly.

**Acceptance:** `httptest` verifies statuses, headers, exact response shapes,
method handling, content types, JSON edge cases, and all input boundaries.
Test each numeric limit at the maximum and above it; test non-ASCII input using
UTF-8 byte limits rather than character counts. Errors must not echo input.
Document a real startup command and a synthetic request/response when runnable.

## Increment 4: Operational behavior and safe observability

Complete startup configuration, structured logs, cancellation, timeouts, and
shutdown using standard-library mechanisms where practical.

Use these initial local defaults, subject to measured revision:

| Budget | Default |
| --- | --- |
| Header read | 5 seconds |
| Complete request read | 10 seconds |
| Response write | 10 seconds |
| Idle connection | 60 seconds |
| Analysis after decoding | 2 seconds |
| Graceful shutdown | 10 seconds |

These are separate budgets, not an end-to-end latency promise. An analysis
deadline can produce a safe 503 only while the connection remains writable.
Socket-level timeouts or client disconnects may prevent a JSON response; record
a safe failure instead of claiming every timeout returns an envelope.

- Validate startup settings and fail visibly on invalid configuration.
- Use an allowlisted logger; never include raw input or error text carrying it.
- Set readiness only after initialization and clear it before shutdown.
- Drain in-flight requests within the shutdown budget.
- Propagate request cancellation into the analyzer.
- Surface deadline, initialization, write, and shutdown failures explicitly.

**Acceptance:** Tests verify log fields and absence of synthetic content markers,
correlation isolation across concurrent requests, cancellation, deadline
handling, and readiness transitions. A real-process check verifies startup,
loopback responsiveness, shutdown, and invalid-startup exit behavior.
Add race testing where supported and report any unavailable prerequisite.

## Increment 5: End-to-end verification and usage documentation

Exercise the runnable API against the scenario set and complete the local
documentation.

- Show actual startup, configuration, request, response, and stop commands.
- Explain every response field, score limitation, and unresolved outcome.
- Verify same-input determinism excluding generated identifiers.
- Verify that no case files, persistence, external calls, or diagnostic
  execution were introduced.
- Document taxonomy coverage and known extraction limitations.
- Run formatting verification, vet, and the complete Go test suite.

**Acceptance:** A fresh local checkout with the documented Go toolchain can run
and validate the service without credentials. All earlier criteria remain met.
The engineer can explain the request flow, rule behavior, errors, and lifecycle.

## Phase 1 completion checklist

- [x] Supported Go version and exact build/test commands are documented.
- [ ] Analysis and probes work on loopback.
- [ ] All request limits and error/status contracts have focused tests.
- [ ] Response shape, review flag, versions, evidence, and collections are tested.
- [ ] Taxonomy, scoring, ambiguity, and deterministic ordering are documented.
- [ ] Reported statements remain separate from verified conclusions.
- [ ] Logs and errors do not expose content, metadata, or arbitrary identifiers.
- [ ] Deadlines, cancellation, readiness, and shutdown are verified.
- [ ] Formatting, vet, and tests pass; race-test limitations are explicit.
- [ ] No future-phase scaffolding, real case data, or external AI calls exist.

Completion authorizes a review of Phase 2, not its automatic implementation.

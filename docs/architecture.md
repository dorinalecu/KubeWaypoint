# Architecture and MVP Behavior

## Purpose and boundaries

KubeWaypoint helps engineers turn an incomplete AKS problem description into an
investigation starting point. It is a public, independent educational project
with practical acceptance criteria, not an internal support tool.

The initial audience includes platform engineers, SREs, Kubernetes
administrators, and support engineers. AKS is the technical specialization, not
an employment or access requirement. Upgrade readiness belongs to the separate
KubeUpgradeCompass project.

Use original content, synthetic scenarios, and licensed public sources only.
The application does not connect to clusters, execute diagnostics, send customer
messages, or transfer content to external AI providers in Phase 1.

All product inputs, outputs, interfaces, and documentation are English.
Multilingual analysis is outside the project scope. Valid UTF-8 is an encoding
requirement, not a promise of language detection or multilingual understanding.

## Phase 1: one Go process

```text
Client
  |
  v
Go HTTP API
  |-- bounded decoding and request validation
  |-- correlation ID and request deadline
  v
In-process deterministic analyzer
  |-- conservative summary and reported fact extraction
  |-- technical classification and supporting rule IDs
  |-- missing information and clarifying questions
  |-- recommended initial validation steps
  v
Versioned JSON response requiring engineer review
```

HTTP handlers own transport concerns, not classification rules. The analyzer
accepts domain input and returns domain output without depending on HTTP.
Transport validation and domain constraints are separate responsibilities:
invalid input must not become an apparently successful analysis.

Keep the initial module at the repository root. The intended organization is
`cmd/kubewaypoint` for startup, `internal/httpapi` for HTTP behavior, and
`internal/analysis` for domain types, rules, and analysis. Create these packages
only when they contain working code. Prefer concrete types and direct function
calls; introduce a small interface only when an actual consumer needs a
replaceable boundary.

Rules are centrally maintained, versioned with the implementation, and reviewed
alongside tests. No plugin system, configurable rule language, or database is
needed.

## Evidence semantics

| Concept | Meaning in this project |
| --- | --- |
| Reported observation or fact | A statement extracted from the input, not independently verified |
| Hypothesis | A possible explanation requiring evidence; not a confirmed cause |
| Verified conclusion | A conclusion supported by reviewed validation evidence |
| Validation step | A proposed way to obtain or check evidence; not a completed action |

Phase 1 cannot independently verify a conclusion and must not generate one.
It recommends checks rather than asserting causal explanations. Reported facts
retain a supporting input excerpt where extraction is possible. Generated text
must not strengthen tentative language, invent an environment detail, or turn
sequence into causation.

The summary is conservative and extractive, not an LLM-style rewrite.
Unrecognized details may remain in the summary without being promoted to
structured facts.

## Classification

The initial outcomes are `node_registration`, `dns`,
`insufficient_information`, and `unknown`.

- Use `insufficient_information` when no actionable symptom can be identified.
- Use `unknown` for an actionable symptom outside the supported categories or
  conflicting category evidence that cannot be resolved conservatively.
- Do not force a supported category from an isolated generic keyword.
- Keep match evidence and tie handling explicit and testable.
- Deduplicate questions and checks; keep their ordering stable.

Confidence is a heuristic support tier, not a calibrated probability or
root-cause probability. Use 0.8 for an explicit supported symptom, 0.5 for a
tentative supported symptom, and 0 for unresolved classification. Repetition does
not increase confidence. Active evidence in both supported categories resolves
to `unknown` with 0, even if the evidence has different strengths. Negated
symptom assertions do not support a classification.

The [analysis contract](analysis-contract.md) defines stable rule IDs, nested
shapes, scenario expectations, and the limitations of this rubric. The rule
engine is not part of increment 1.

## Planned public contract

`POST /api/v1/analyze` accepts one JSON object. Field names use `snake_case`.

| Request field | Constraint |
| --- | --- |
| `description` | Required string; nonempty after trimming; at most 16 KiB of UTF-8 text |
| `case_id` | Optional nonblank opaque string, at most 128 UTF-8 bytes; synthetic identifiers only during development |
| `metadata` | Optional string-to-string object, at most 16 entries; keys 1-64 UTF-8 bytes, values at most 256 UTF-8 bytes |

The complete request body is limited to 32 KiB, including JSON encoding.
The body limit and individual field limits both apply. Invalid UTF-8, unknown
fields, wrong types, duplicate JSON object keys, trailing JSON values, and
explicit `null` values for these fields are rejected. Optional fields may be
omitted. Metadata is advisory and does not override description-based rules in
the MVP.

The successful response has these required fields:

| Field | Purpose |
| --- | --- |
| `schema_version` | Public response contract version, initially `1.0` |
| `rules_version` | Version of the deterministic rules used |
| `analysis_id` | Server-generated identifier for this analysis; not a persisted record |
| `correlation_id` | Request correlation identifier |
| `requires_engineer_review` | Always `true` |
| `summary` | Conservative restatement of the reported issue |
| `category` | One of the four defined outcomes |
| `confidence` | Supported-category heuristic score |
| `classification_evidence` | Matched rule identifiers and supporting input excerpts |
| `reported_facts` | Extracted statements with supporting input excerpts |
| `missing_information` | Explicit information gaps |
| `clarifying_questions` | Focused questions associated with those gaps |
| `initial_checks` | Proposed checks with a purpose and expected evidence |

Collection fields are arrays, including when empty. The
[analysis contract](analysis-contract.md) defines their nested field names.
Do not add hypotheses or verified conclusions implicitly during implementation. A valid
unknown or insufficient-information analysis returns HTTP 200; it is not a
transport error. Case content and metadata are not echoed wholesale.

Same input and rule version must yield the same analysis content and ordering.
Generated analysis and correlation identifiers are excluded from this guarantee.
Document incompatible contract changes rather than silently changing v1.

## HTTP, lifecycle, and errors

Require `application/json` for analysis requests, allowing media-type parameters.
Use a JSON error envelope containing `schema_version`, `correlation_id`, and
`error` with stable `code` and safe `message` fields.

| Status | Meaning |
| --- | --- |
| 400 | Invalid JSON, field validation, or correlation header |
| 404 | Unknown route |
| 405 | Unsupported method; include the appropriate `Allow` header |
| 413 | Request body exceeds its bound |
| 415 | Unsupported media type |
| 503 | Service not ready or analysis deadline exceeded while a response can still be written |
| 500 | Unexpected internal failure; no internal details in the response |

`X-Correlation-ID` is returned on responses. A supplied value must contain 1-64
ASCII characters from letters, digits, `.`, `_`, and `-`. Generate a value when
absent; on an invalid value, reject the request and use a new safe value for the
error. Never trust or log an unvalidated header.

`GET /healthz` reports process health. `GET /readyz` returns 200 only after
initialization and while accepting analysis; return 503 during shutdown.
Probe bodies contain status only, not case content or configuration.

Default listening is loopback-only at `127.0.0.1:8080`. Phase 1 has no
authentication and is not suitable for public exposure. Startup configuration
must be validated, and failures must be visible with a nonzero exit status.
Define explicit header-read, request-read, response-write, idle, analysis, and
shutdown deadlines in increment 4. Propagate cancellation into analysis.
Graceful shutdown stops admission before draining in-flight requests.

## Data and observability

No database, request history, or case files are created. Case data is used in
memory only; statelessness is not a guarantee of secure memory erasure.

Use an allowlist for structured logs: correlation ID, known route name, method,
status, duration, versions, and safe error code. Do not log descriptions,
summaries, extracted excerpts, metadata, case IDs, arbitrary URLs, headers, or
raw validation errors. Inputs must not be placed in panic diagnostics.

## Phase 2 and beyond

```text
Client -> Go public API -> versioned internal HTTP/JSON -> Python analyzer
                                                        |
                                                        +-- deterministic rules
                                                        +-- optional LLM later
```

Go owns public validation, orchestration, deadlines, safe observability, and
validation of the Python response. Python owns analysis. The deterministic
analysis responsibility moves to Python; do not maintain duplicate rule engines
without a separate justified decision.

Contract tests must protect public behavior across the move. Dependency failures
remain explicit; a future fallback must be deliberate and identified in output.
LLM providers and retrieval are later choices, not current dependencies.

See [ADR 0001](adr/0001-go-first-deterministic-stateless-mvp.md),
the [implementation plan](phase-1-plan.md), and the [roadmap](roadmap.md).

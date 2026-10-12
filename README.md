# KubeWaypoint

An engineer-reviewed AKS case-intake assistant that turns incomplete problem
descriptions into structured investigation starting points.

**Status:** Phase 0 - Engineering Foundation. This repository currently contains
product documentation only. There is no runnable application yet.

## Why KubeWaypoint?

Support engineers often receive incomplete descriptions of complex issues.
Before troubleshooting can begin, they need to normalize the problem, identify
missing evidence, ask focused questions, and choose useful initial checks.

KubeWaypoint aims to reduce repetitive intake work and make that process more
consistent. It does not replace the engineer's technical judgment, customer
relationship, or responsibility for the investigation.

The project also demonstrates a progression from AKS and cloud troubleshooting
expertise to practical software engineering: Go services, explicit contracts,
testing, observability, and eventually Python-based AI capabilities and
cloud-native delivery.

This is an independent portfolio project, not an official Microsoft product.
Its implementation and documentation will be original, without importing
internal procedures, support-pack content, or customer data from other projects.

## Product behavior

Given an unstructured AKS case description, KubeWaypoint will help the engineer:

- Restate the reported issue concisely.
- Identify a likely technical category.
- Extract facts reported in the input without treating them as independently
  verified evidence.
- Identify missing information.
- Formulate clarifying questions.
- Recommend initial diagnostic checks.

The analysis must distinguish reported facts, observations, hypotheses, missing
information, and validation steps. It must not present a possible cause as a
confirmed root cause without supporting evidence.

All output is advisory and requires engineer review. Recommended checks are not
executed automatically, and customer communication is not sent automatically.

## First MVP: a deterministic Go service

The first implementation milestone will provide:

- `POST /api/v1/analyze` for a case description, an optional case identifier,
  and optional metadata that must not contain secrets.
- A structured response containing an analysis identifier, normalized summary,
  primary category, confidence score, extracted facts, missing information,
  clarifying questions, recommended initial checks, correlation ID, and schema
  version.
- `GET /healthz` for process health.
- `GET /readyz` for the ability to process requests.
- Bounded and validated requests, consistent error objects, and appropriate HTTP
  status codes.
- A small, centrally defined AKS taxonomy, starting with node registration and
  DNS, plus explicit insufficient-information and unknown outcomes.
- Transparent, testable classification rules.
- Structured logs, correlation IDs, intentional server timeouts, and graceful
  shutdown.
- Unit and HTTP handler tests using synthetic case descriptions.

The confidence score will describe support for a classification, not the
probability of a root cause. Its exact scoring rules and limitations will be
defined and tested during implementation.

The MVP will run locally without cloud credentials or external AI calls. It will
be stateless: analyses will not be stored in a database.

### Explicitly outside the first MVP

- A Python service or real LLM integration.
- Retrieval-augmented generation, embeddings, or a vector database.
- A web interface, authentication, or case history.
- Dockerfiles, Kubernetes manifests, Terraform, or Azure infrastructure.
- GitHub Actions workflows, GitOps, or a service mesh.
- Autonomous diagnostic execution or customer communication.

## Project boundary

KubeWaypoint focuses on case intake and initial troubleshooting plans. Assessing
cluster upgrade readiness is a separate product concern, not a capability
promised by this project's roadmap.

[KubeUpgradeCompass](https://github.com/dorinalecu/KubeUpgradeCompass) is a
separate, planned portfolio project for evidence-based AKS upgrade-readiness
assessment. It has no implementation yet and is deferred while KubeWaypoint is
developed. Neither project requires the other to run.

## Architecture direction

The initial request flow stays inside one process:

```text
Client
  -> Go HTTP API
  -> request validation and correlation ID
  -> deterministic analyzer
  -> structured, versioned JSON response
```

Business logic will be separate from HTTP handlers. Keeping the first analyzer
in-process lets us establish domain behavior and tests before introducing
distributed-system failure modes.

A later milestone will move analysis behind an explicit, versioned Go-to-Python
contract. Python will eventually support deterministic and replaceable LLM
providers. Curated AKS knowledge and retrieval come after optional LLM
integration.

Dependency failures must remain visible. Any future fallback must be intentional,
documented, and identified in the response rather than disguised as a successful
primary analysis.

## Development roadmap

| Phase | Outcome |
| --- | --- |
| 0. Engineering Foundation | Product scope, architecture decisions, development guidance, and a precise Go implementation plan |
| 1. Go MVP | A locally runnable, deterministic case-analysis API with tests |
| 2. Python Analysis Service | Two independently testable services with contract and integration tests |
| 3. Containerization | Reproducible, non-root local container execution |
| 4. LLM Integration | Optional AI analysis with replaceable providers and validated output |
| 5. AKS Knowledge and RAG | Curated retrieval with source references and grounding |
| 6. Observability and Evaluation | Measurable behavior, representative scenarios, and failure testing |
| 7. Azure Infrastructure | Repeatable infrastructure with security, cost, and cleanup guidance |
| 8. CI/CD and GitOps | Auditable builds, deployment, promotion, and rollback |
| 9. Agentic Case Intake | Controlled workflow steps with human review and an audit trail |

Only the current milestone will be implemented. Future phases do not justify
unused scaffolding or dependencies.

## Local development

Phase 0 requires Git and a text editor. Phase 1 will additionally require a
supported Go toolchain; the minimum Go version will be selected and recorded when
the module is created.

There are no application run or test commands yet. Once the Go module exists,
the baseline checks, run from its directory, will be:

```text
gofmt -l .
go vet ./...
go test ./...
```

Formatting verification should report no files, and vet and tests must pass.
Exact application startup instructions will be added with the runnable service.
No Azure subscription, Python installation, Docker installation, or AI provider
key is required for the first MVP.

## Engineering and data-safety principles

- Prefer the standard library and add dependencies only for a demonstrated need.
- Keep domain logic testable, contracts explicit, and configuration outside code.
- Propagate request context and handle errors rather than hiding them.
- Validate behavior with deterministic tests, including malformed input and
  ambiguous classifications.
- Never commit customer case data, credentials, tokens, or other secrets.
- Use synthetic examples without real customer or resource identifiers.
- Do not log complete case descriptions or optional metadata by default.
- Keep external AI data transfer disabled unless explicitly configured in a
  future milestone.
- Document limitations and require engineer review before external use.
- Use small, focused changes and conventional commit messages such as `docs:`,
  `feat:`, `fix:`, and `test:`.

## Next foundation step

Document the current and future architecture and record the decision to start
with a Go-first, deterministic, stateless MVP. Then complete the local
development guide, roadmap details, and the step-by-step Phase 1 plan before
writing application code.

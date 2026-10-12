# KubeWaypoint

An engineer-reviewed AKS case-intake assistant that turns incomplete problem
descriptions into structured investigation starting points.

**Status:** Phase 1 - Domain foundation. The first increment introduces Go domain
contracts, input validation, and tests. There is no analyzer or runnable HTTP
service yet.

## Why KubeWaypoint?

Platform engineers, SREs, Kubernetes administrators, and support engineers often
receive incomplete descriptions of complex issues. Before troubleshooting can
begin, they need to normalize the problem, identify missing evidence, ask focused
questions, and choose useful initial checks.

KubeWaypoint aims to reduce repetitive intake work and make that process more
consistent. It does not replace the engineer's technical judgment or
responsibility for the investigation. Here, "case intake" means structuring a
reported technical problem, not integrating with an internal support system.

The project also demonstrates a progression from AKS and cloud troubleshooting
expertise to practical software engineering: Go services, explicit contracts,
testing, observability, and eventually Python-based AI capabilities and
cloud-native delivery.

This is an independent public portfolio project, not an official Microsoft
product or a tool restricted to Microsoft employees. It requires no internal
systems or privileged access. Its implementation and documentation will be
original, using synthetic examples and licensed public references rather than
internal procedures, support-pack content, or customer data from other projects.

The product uses English for inputs, outputs, documentation, and interfaces.
Multilingual analysis is outside the project scope.

## Engineering foundation

- [Architecture and MVP behavior](docs/architecture.md)
- [ADR 0001: Go-first, deterministic, stateless MVP](docs/adr/0001-go-first-deterministic-stateless-mvp.md)
- [Development guide](docs/development.md)
- [Detailed roadmap](docs/roadmap.md)
- [Incremental Phase 1 implementation plan](docs/phase-1-plan.md)
- [Analysis contract, scoring, and synthetic scenarios](docs/analysis-contract.md)

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

The confidence score describes support for a classification, not the probability
of a root cause. The agreed rubric is 0 for unresolved classification, 0.5 for a
tentative reported symptom, and 0.8 for an explicit reported symptom. Repetition
does not increase the score, and active evidence in both supported categories
resolves to unknown. See the [contract](docs/analysis-contract.md) for limitations.

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

The Go module is at the repository root and requires Go 1.27.2 or newer.
The initial validation toolchain is Go 1.27.2 on Windows amd64.

Run the following from the repository root:

```text
gofmt -l .
go build ./...
go vet ./...
go test ./...
```

Formatting verification should report no files, and vet and tests must pass.
Build verifies the domain package; it does not produce an application executable.
There is no application startup command yet. See the
[development guide](docs/development.md) for using the locally installed Go
toolchain without changing persistent PATH settings. No Azure subscription,
Python installation, Docker installation, or AI provider key is required.

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

## Next step

Review the domain foundation and synthetic scenario expectations, then implement
the deterministic analyzer in increment 2 of the
[Phase 1 plan](docs/phase-1-plan.md) when requested. Do not implement later phases
or install their dependencies in advance.

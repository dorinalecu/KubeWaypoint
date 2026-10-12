# Development Guide

## Current state

Phase 0 is documented. Phase 1 increment 1 adds the root Go module and domain
types, validation, and tests in `internal/analysis`. There is no analyzer,
runnable service, or Python environment yet.

Work in the existing project folder without creating a worktree unless requested.
Do not install dependencies or introduce later-phase tooling preemptively.

## Prerequisites by milestone

| Milestone | Required locally |
| --- | --- |
| Phase 0 | Git and a text editor |
| Phase 1 | A supported stable Go toolchain |
| Phase 2 | A supported Python runtime and the tools selected for that milestone |
| Phase 3 | A compatible container runtime |
| Later cloud or AI phases | Explicitly authorized accounts, configuration, and credentials |

The minimum Go version is 1.27.2, verified against the
[official release manifest](https://go.dev/dl/?mode=json).
The initially tested toolchain is Go 1.27.2 on Windows amd64.
The module is `github.com/dorinalecu/KubeWaypoint` at the repository root.
Record additional tested environments as they are verified, not by assumption.

Python, containers, Azure credentials, and AI provider keys are not Phase 1
prerequisites. An unavailable future toolchain does not block documentation.

## Working agreements

- Explain the increment's scope, behavior, and trade-offs before editing.
- Keep code, documentation, identifiers, commit messages, and all product
  input/output in English. Multilingual analysis is outside scope.
- Prefer standard libraries, direct composition, and explicit error handling.
- Create packages and directories only when they have a current purpose.
- Keep implementation and directly related documentation synchronized.
- Never revert unrelated work or commit or push without a relevant request.
- Use conventional commit subjects such as `docs:`, `feat:`, `fix:`, and `test:`.
- Authorized commits use a verified GitHub noreply author email and the
  `Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>` trailer.

For the engineer's current development workflow, architecture, task boundaries,
acceptance criteria, and review belong to the orchestrator; all coding, including
tests and fixes, is delegated to Claude Sonnet 5.5. This is a contributor workflow
preference, not a runtime model dependency or a requirement for other users.

## Baseline Go validation

Run these commands from the module root:

```text
gofmt -l .
go build ./...
go vet ./...
go test ./...
```

`gofmt -l .` must report no files; its exit status alone does not establish
formatting compliance. Build, vet, and tests must pass. There is no executable
entry point yet; build validates packages only. Start with the smallest relevant
test selection while iterating, then run the baseline checks for an increment.
Use existing tooling rather than introducing a parallel test framework.

Add race testing when concurrency is introduced, where the toolchain supports
it. On Windows, the Go race detector may require a compatible C toolchain;
report a missing prerequisite rather than claiming that race testing passed.

Tests must not need cloud credentials, network access to external services,
customer environments, or an AI provider. Use unit tests for domain behavior,
`httptest` for handlers, and real-process lifecycle checks when startup and
shutdown exist.

### Personal Windows toolchain

With the engineer's permission, Go 1.27.2 was installed from the official Windows
amd64 archive into `%LOCALAPPDATA%\Programs\Go\1.27.2\go`. Its SHA-256 was verified
against the official release manifest. No machine or user PATH setting was
changed. This location is a local convenience, not a repository dependency.

In a PowerShell terminal, the following enables it for that terminal only:

```powershell
$env:Path = "$env:LOCALAPPDATA\Programs\Go\1.27.2\go\bin;$env:Path"
$env:GOTOOLCHAIN = "local"
$env:GOPROXY = "off"
go version
```

Then run the baseline commands from the repository root. `GOTOOLCHAIN=local`
prevents automatic toolchain downloads; `GOPROXY=off` is suitable for the current
standard-library-only module. Revisit dependency configuration only when an
approved dependency is actually added.

## Synthetic scenarios and data safety

Use short, original scenarios covering DNS, node registration, insufficient
information, unsupported symptoms, and conflicting category evidence.
The agreed [contract scenarios](analysis-contract.md) describe future analyzer
expectations; increment-1 contract fixtures do not prove classification behavior.
Examples must not be copied or adapted from real cases or internal material.
Use reserved example domains and obviously synthetic identifiers when needed.

Review both committed fixtures and runtime logs for accidental disclosure.
Do not log raw descriptions, metadata, summaries, or evidence excerpts. Do not
use a real customer's input merely because the application runs locally.
External AI transfer stays disabled unless separately authorized and configured
in a later milestone.

## Local execution policy

Phase 1 listens on loopback by default. It has no authentication and must not be
made publicly reachable. Document the actual startup command, supported
configuration, example request, response, and shutdown procedure when implemented.
Do not describe planned commands as available commands.

## Definition of a completed increment

The intended behavior works, relevant positive and negative tests pass, failures
are visible, limitations are documented, and no unused scaffolding or unrelated
changes were added. A model's assertion of success is not verification.

Follow the [Phase 1 plan](phase-1-plan.md) and
[architecture](architecture.md) for concrete acceptance criteria.

# Analysis Contract and Synthetic Expectations

## Status

This document fixes Phase 1 domain shapes and future analyzer expectations.
Increment 1 implements domain types and input validation only. It does not
implement the analyzer, public response envelope, JSON request decoding, HTTP
server, or classification rules.

All product input/output is English. These original synthetic scenarios are not
adapted from customer cases or internal procedures.

## Input and validation boundary

The domain input contains `description`, optional `case_id`, and optional
`metadata`. Domain validation enforces valid UTF-8, nonblank description and
provided case ID, and the field limits in the [architecture](architecture.md).
Byte limits apply before trimming. Validation does not mutate input.

An omitted case ID is represented by a nil pointer; a provided empty or
whitespace-only ID is invalid. Omitted or empty metadata is valid; keys must be
nonempty and values may be empty. Metadata does not select or override category.

Domain validation errors contain only fixed field names and error codes, never
input values or metadata keys. The future HTTP layer will map them into a safe
error envelope. JSON null, duplicate/unknown keys, media types, and the 32 KiB
body limit belong to that transport layer and are not implemented in increment 1.

## Domain result

`Result` contains these JSON fields:

| Field | Type |
| --- | --- |
| `rules_version` | String, initially `1.0` |
| `requires_engineer_review` | Boolean; generated results must use `true` |
| `summary` | String |
| `category` | String enum: `node_registration`, `dns`, `insufficient_information`, `unknown` |
| `confidence` | Number: 0, 0.5, or 0.8 |
| `classification_evidence` | Array of rule matches |
| `reported_facts` | Array of supported reported statements |
| `missing_information` | Array of information gaps |
| `clarifying_questions` | Array of questions |
| `initial_checks` | Array of proposed checks |

Nested objects have exactly these fields:

| Object | String fields |
| --- | --- |
| Classification evidence | `rule_id`, `source_excerpt` |
| Reported fact | `statement`, `source_excerpt` |
| Information gap | `id`, `description` |
| Clarifying question | `id`, `information_id`, `question` |
| Initial check | `id`, `description`, `purpose`, `expected_evidence` |

`information_id` references an information-gap ID from the same result. IDs for
questions and checks are stable template identifiers, not customer identifiers.
An excerpt is verbatim reported content; it is not independent verification.
A check describes an action for engineer review, not an executed command.

Generated results must provide every collection as an array, even when empty.
Increment-1 serialization tests use explicitly initialized empty slices.
A zero-value Go result is not a valid generated analysis; an analyzer will be
responsible for producing complete results in increment 2.

The future HTTP response adds `schema_version` (`1.0`), `analysis_id`, and
`correlation_id` at the top level. These are not domain-result fields.
No timestamps, hypotheses, verified conclusions, case-history links, or implicit
fallbacks are added.

## Confidence rubric

| Active support | Category outcome | Confidence |
| --- | --- | --- |
| Explicit recognized symptom in one supported category | That category | 0.8 |
| Tentative recognized symptom in one supported category | That category | 0.5 |
| Both supported categories have active evidence | `unknown` | 0 |
| Actionable symptom outside supported categories | `unknown` | 0 |
| No actionable symptom | `insufficient_information` | 0 |

Within a single category, use the highest active tier. Repeated wording or
multiple rules do not accumulate confidence. Cross-category conflict is evaluated
before choosing a tier, so explicit DNS plus tentative node-registration evidence
still yields unknown.

The numeric values are fixed heuristic labels, not empirically calibrated
probabilities. In particular, 0.8 does not mean an 80% chance that DNS is broken
or that DNS is the root cause. Unresolved outcomes use 0 for support of a
supported category, not certainty that no problem exists.

## Initial rule catalog for increment 2

These are narrow English symptom rules, not general natural-language inference.
Match case-insensitively with token boundaries; normalized matching must retain
the original excerpt. Do not classify from a category word alone.

| Rule ID | Category | Positive symptom forms | Tier |
| --- | --- | --- | --- |
| `NR001` | `node_registration` | Nodes fail to register or cannot register with the cluster | Explicit |
| `NR002` | `node_registration` | Nodes may be failing to register, may fail to register, or might fail to register | Tentative |
| `DNS001` | `dns` | Pods cannot resolve DNS names, or DNS resolution fails | Explicit |
| `DNS002` | `dns` | DNS resolution may be failing, may fail, or might fail | Tentative |

Tentative markers such as `may`, `might`, or `possibly` in the symptom clause
downgrade an otherwise explicit recognized symptom to the tentative rule.
Negated symptom assertions do not activate a rule. A recognized failure form
such as `cannot resolve` is an active symptom, not a negation of the problem.
Positive success statements are not failure evidence.

Increment 2 must test negation scope and document its limitations. It must not
pretend to understand arbitrary English phrasing. For conservative matching,
punctuation-delimited clauses separate unrelated assertions; quoted, historical,
or linguistically complex descriptions may remain unresolved.

Deduplicate matches by rule ID, retaining the first supporting excerpt.
Order evidence by rule ID and template collections by stable ID. An actionable
unsupported symptom must not be reduced to a generic category keyword match.

## Synthetic scenario matrix

The following are expected future analyzer outcomes, not passing analyzer tests
in increment 1. Summaries retain reported wording and facts retain source
excerpts without claims of independent verification.

| ID | Description | Expected category | Score | Evidence | Follow-up profile |
| --- | --- | --- | --- | --- | --- |
| `S01` | New nodes fail to register with the cluster. | `node_registration` | 0.8 | `NR001` | Node |
| `S02` | Pods cannot resolve DNS names. | `dns` | 0.8 | `DNS001` | DNS |
| `S03` | New nodes may be failing to register with the cluster. | `node_registration` | 0.5 | `NR002` | Node |
| `S04` | Pod DNS resolution may be failing. | `dns` | 0.5 | `DNS002` | DNS |
| `S05` | The cluster is broken. | `insufficient_information` | 0 | None | General |
| `S06` | DNS. | `insufficient_information` | 0 | None | General |
| `S07` | Pods repeatedly restart with OOMKilled. | `unknown` | 0 | None | General |
| `S08` | Pods can resolve DNS names. New nodes register successfully. | `insufficient_information` | 0 | None | General |
| `S09` | New nodes do not fail to register with the cluster. | `insufficient_information` | 0 | None | General |
| `S10` | DNS resolution does not fail. | `insufficient_information` | 0 | None | General |
| `S11` | New nodes fail to register with the cluster. Pods cannot resolve DNS names. | `unknown` | 0 | `DNS001`, `NR001` | Conflict |
| `S12` | New nodes may be failing to register with the cluster. Pods cannot resolve DNS names. | `unknown` | 0 | `DNS001`, `NR002` | Conflict |
| `S13` | Pods cannot resolve DNS names. Pods cannot resolve DNS names. | `dns` | 0.8 | `DNS001` once | DNS |
| `S14` | DNS resolution fails. New nodes do not fail to register with the cluster. | `dns` | 0.8 | `DNS001` | DNS |

Empty and whitespace-only descriptions fail input validation; they are not
successful insufficient-information results. Add independent boundary and
encoding cases in validation tests, not category rules.

For these minimal descriptions, none reports start time, precise scope, or a
specific diagnostic error. The follow-up profiles below therefore identify
those gaps. Increment 2 must also test richer inputs where an already reported
detail removes the corresponding gap and question.

## Initial follow-up profiles

Each row defines a gap description and its associated question. Gap ID is the
question's `information_id`; question ID is `q_` followed by the gap ID.
Do not ask for a detail that has already been extracted from the report.

| Profile | Gap ID | Missing information | Question |
| --- | --- | --- | --- |
| Node | `node_error` | Exact node-registration error is not reported. | Is there a node-registration error message, and if so, what is its exact text? |
| Node | `node_scope` | Affected node or node-pool scope is not reported. | Are all new nodes affected, or only a particular node pool? |
| Node | `node_start` | Symptom start time is not reported. | When was the node-registration issue first reported? |
| DNS | `dns_error` | Exact DNS error is not reported. | Is there a DNS error message, and if so, what is its exact text? |
| DNS | `dns_scope` | Affected workload scope is not reported. | Which workloads and namespaces are affected? |
| DNS | `dns_start` | Symptom start time is not reported. | When did the DNS symptom start? |
| General | `general_detail` | An actionable symptom or sufficient diagnostic detail is missing. | What operation fails, and what exact error is observed? |
| General | `general_scope` | Affected scope is not reported. | Which workloads or cluster components are affected? |
| Conflict | `conflict_order` | Relationship between the reported symptoms is unclear. | When did each symptom start, and which was observed first? |
| Conflict | `conflict_scope` | Scope of each reported symptom is unclear. | Which nodes and workloads are affected by each symptom? |

Proposed initial checks are advisory:

| Profile | Check ID | Description | Purpose | Expected evidence |
| --- | --- | --- | --- | --- |
| Node | `check_node_error` | Review available node-registration error output. | Establish the reported failure stage without assuming a cause. | Exact error text and the registration stage where it occurs. |
| DNS | `check_dns_error` | Review available DNS lookup error output from an affected workload. | Distinguish the reported failure mode without assuming a DNS component is at fault. | Queried name, error text, and affected workload context. |
| General | `check_symptom_detail` | Review existing output for the reported failing operation. | Establish a concrete investigation starting point. | Operation, error text, and affected component. |
| Conflict | `check_symptom_timeline` | Compare available timestamps and scope for both symptoms. | Test whether the symptoms should be investigated together. | Separate timelines and affected scopes; no inferred causal relationship. |

Profiles are deliberately small. They are investigation prompts, not internal
support procedures or automatic diagnostic execution.

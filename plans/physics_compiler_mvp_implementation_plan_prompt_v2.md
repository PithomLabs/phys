# Physics Compiler MVP — Implementation Plan Prompt

## Role

You are the **implementation architect**.

Do **not** implement the system. Produce the implementation plan that a coding agent will execute.

The companion file `physics_compiler_mvp_specs_v2.md` is the **normative implementation specification**. It is authoritative. Do not redesign it, reinterpret it, or introduce alternative architectures.

The purpose of this prompt is to turn that fixed specification into a sequenced engineering plan.

---

## Mission

Plan the smallest coherent Go implementation that demonstrates this thesis:

> An AI agent can read a machine-readable physics corpus, construct typed physical primitives, perform symbolic derivations under explicit assumptions, have MRC reject structurally invalid operations, record committed reasoning with provenance, formulate a provisional hypothesis, and hand a reviewable research artifact to humans — without the system adjudicating physical truth.

The AI is the theorist.

The Physics Compiler is the formal symbolic pen-and-paper substrate.

Human researchers and empirical reality remain the final scientific authority.

---

## Six product requirements

The plan MUST demonstrate all six:

1. **Fundamental primitives + MRC:** the library represents physical equations using typed physical primitives and enforces MRC through Go typing, constructor authority, and library operation contracts.
2. **Frameworks as evaluation targets:** packages such as special relativity are scoped frameworks with explicit assumptions, limits, anomalies, and derivation targets; they are not the final ontology of nature.
3. **Machine-readable physics corpus:** packages plus validated semantic manifests provide the AI-readable corpus.
4. **Collective limitations are explicit:** assumptions, provenance, regimes, conventions, anomalies, candidate status, and human handoff are represented rather than hidden.
5. **No truth adjudication:** the library checks formal/structural validity only. It never computes physical truth, truth scores, or theory rankings.
6. **Hypothesis formation:** the AI may introduce provisional concepts and formulate falsifiable research candidates, including possible theory bridges, but promotion and empirical adjudication remain human-controlled.

---

## Non-negotiable architecture

The plan MUST preserve these decisions:

- Go is the source language.
- No `.phys` language.
- No custom lexer/parser/frontend.
- No general-purpose mathematics ontology.
- No full CAS.
- No theorem prover.
- Symbolic reasoning only.
- No numerical simulation/runtime.
- No automatic empirical validation.
- No truth scores or theory rankings.
- No automatic promotion of hypotheses.
- No EBP 2.1 integration.
- MRC has one semantic source of truth.
- `physvet` is deferred from MVP.
- Plan Z's `Identify`/provenance/ledger concept is retained.
- DeepSeek's domain-package corpus structure is retained.
- Qwen's "AI as theorist / library as pen and paper" framing is retained, without epistemic overclaims.
- Candidate concepts have open representational freedom but closed trusted authority.
- MRC itself is fallible and versioned; MVP permits revision by versioning, not by runtime bypass.

---

## Scope of the MVP

Only these populated implementation areas are permitted:

```text
core/
ops/
mechanics/
relativity/
hypothesis/
docs/paper-translation/
```

Do not add stub packages for:

```text
electromagnetism
qm
qft
statmech
```

Special relativity is the only relativity content in MVP.

---

## Required implementation-plan output

Return one implementation plan with these exact sections:

### 1. Executive architecture

Explain the final MVP architecture in enough detail to orient the implementer, but do not invent alternatives.

### 2. Exact repository tree

Show every source/test/data/document file the plan expects to create or modify.

Tests MUST use Go's normal adjacent `*_test.go` convention.

### 3. Dependency-ordered implementation sequence

Give an ordered sequence in which each step has:

- purpose
- inputs/dependencies
- files affected
- implementation result
- verification gate

Do not write code.

### 4. Core model implementation

Describe the exact implementation of the types defined by the specs, including the relationship among:

```text
domain nominal types
core.Object
core.Expr
dimensions
assumptions
conventions
provenance
```

Do not redesign their relationships.

### 5. MRC implementation map

Create a table:

| Rule | Exact check | Enforcement point | Failure | Test |

Cover every MVP MRC rule.

Distinguish clearly between:

- Go compile/package boundaries
- constructor authority
- operation-time MRC
- deferred `physvet`

### 6. Symbolic engine

Describe only the specified MVP expression nodes and operations.

Explain canonicalization, equality, hashing, and exact rational encoding.

Do not expand the symbolic engine beyond the specification.

### 7. Physics corpus

Explain how `mechanics/manifest.json` and `relativity/manifest.json` are embedded, decoded into typed structures, validated, and cross-checked against executable constructors/relations.

### 8. Assumptions, conventions, provenance, and candidate containment

Explain the exact propagation/containment laws already fixed by the specs.

### 9. Derivation ledger

Describe session authority, step recording, `Identify`, hash chaining, replay validation, commit/seal boundaries, and tamper detection.

### 10. Hypothesis and ResearchCandidate

Describe candidate creation, provenance contamination, falsifiability fields, anomaly references, sealing, and the human boundary.

### 11. Canonical vertical slices

Describe the exact implementation/test sequence for:

- classical mechanics
- `F = ma`
- `E = mc²`
- explicit `Identify`
- candidate contamination
- anomaly/falsifiability
- ledger integrity

Do not replace the specified derivations with new ones.

### 12. Coverage matrix

This is mandatory.

Map **every normative MUST in the specifications** and every MRC rule ID to at least one concrete test.

No requirement may remain untested merely because it is "documentation."

### 13. Acceptance gate

List the exact pass conditions for MVP completion.

### 14. Explicit deferrals

Repeat all v0.5+ deferrals that must not leak into the implementation.

### 15. Scope accounting

Provide a concise estimate of:

- source files
- test files
- manifest/document files
- major exported types
- major operations

Target **no more than 40 repository files total** for MVP, excluding `go.mod`, `go.sum`, generated artifacts, and VCS metadata. If the plan exceeds 40, justify every excess file against a normative MVP requirement.

The estimate exists to detect bloat.

---

## Planning rules

The plan agent MUST NOT:

- propose multiple architectures
- reopen settled decisions
- substitute a different symbolic representation
- invent a custom parser
- turn `physvet` into MVP scope
- add deferred domain packages
- create "future-proof" abstractions without a required MVP use
- silently omit a normative requirement
- use vague language such as "the implementer can decide"

If the specification still contains a genuine contradiction or under-specified point, create an `Open Spec Items` section at the end of the plan and identify it precisely instead of improvising. Do not invent a silent resolution.

The plan MUST also include a coverage matrix mapping every normative MUST in the specifications and every `MRC-001` through `MRC-008` rule to its implementation location and at least one test.

The plan MUST also include a coverage matrix mapping every normative MUST in the specifications and every `MRC-001` through `MRC-008` rule to its implementation location and at least one test.

---

## Final constraint

The plan is successful only if a separate coding agent can execute it **without making architectural decisions**.

The implementation should be small enough to understand as one vertical slice, yet complete enough to demonstrate:

```text
AI reads corpus
      ↓
typed physical construction
      ↓
assumptions + conventions
      ↓
symbolic derivation
      ↓
MRC enforcement
      ↓
explicit physical identification
      ↓
canonical/replayable derivation
      ↓
provisional hypothesis
      ↓
falsifiability + anomaly context
      ↓
sealed ResearchCandidate
      ↓
human review
```

Anything not required for this loop belongs outside the MVP.

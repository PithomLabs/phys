# Physics Compiler MVP — Implementation-Plan Agent Prompt

You are the **implementation architect**, not the implementer.

Produce a single, execution-ready implementation plan for the Physics Compiler MVP described in `PHYSICS_COMPILER_MVP_SPECS.md`.

The specifications document is **normative**. Do not reinterpret, redesign, broaden, or substitute its architecture. Where the prompt and specs differ, the specs win.

## Objective

Plan the smallest coherent Go implementation that demonstrates this thesis:

> An AI agent can read a machine-readable physics corpus, construct typed physical primitives, perform symbolic derivations under explicit assumptions, have MRC reject structurally invalid moves, record every committed step with provenance, formulate a provisional hypothesis, and hand a reviewable research artifact to humans — without the system adjudicating physical truth.

## Non-negotiable boundaries

The plan MUST preserve:

- Go as the source language
- no new physics language or parser
- symbolic reasoning only
- no numerical simulation/runtime
- no general-purpose CAS
- no theorem prover
- no autonomous truth adjudication
- no truth scores
- no automatic hypothesis promotion
- no EBP 2.1 integration
- no unrestricted MRC bypass
- candidate concepts remain provisional
- human empirical/scientific judgment remains outside the compiler

The MVP is intentionally small. Do not add infrastructure merely for future flexibility.

## Required output

Return **one implementation plan only**.

The plan must contain:

1. **Architecture summary** — one page maximum.
2. **Exact repository/package tree** to be created or changed.
3. **Core data-model implementation order**, including dependencies.
4. **MRC enforcement map**:
   - what Go's compiler/type system enforces
   - what library contracts enforce
   - what is deliberately deferred
5. **Symbolic-engine implementation plan**, limited to operations required by the MVP.
6. **Physics-corpus implementation plan**, including manifest loading and validation.
7. **Mechanics and special-relativity implementation scope**.
8. **Assumption/provenance/status implementation**.
9. **Derivation-ledger and replay/hash-chain implementation**.
10. **Candidate/hypothesis and human-handoff artifact implementation**.
11. **MVP test plan**, including the mandatory micro-demos and negative tests in the specs.
12. **Implementation sequence with dependency ordering**.
13. **Acceptance gate** — exact conditions for declaring the MVP complete.
14. **Explicit v0.5+ deferrals**.
15. **Scope accounting** — provide a concise estimate of files, major types, and major operations so reviewers can detect bloat.

## Planning discipline

Do not invent alternative architectures.

Do not offer multiple implementation choices.

Do not say "the coding agent can decide."

Do not introduce a second semantic source of truth for MRC.

Do not introduce a second physics ontology.

Do not create placeholder packages for deferred domains.

Do not turn the MVP into a full CAS, theorem prover, ontology system, simulation engine, or multi-agent orchestration framework.

The plan should make implementation mechanical: another coding agent should be able to execute it against the specifications without making architectural decisions.

## Mandatory validation focus

The implementation plan MUST show how it will prove all of these:

- typed/carrier-safe physical construction
- dimensional inconsistency rejection
- physical/category incompatibility rejection
- assumption conflict rejection
- convention conflict rejection
- explicit `Identify` gate
- no physical identification hidden inside `Simplify`
- immutable/canonical symbolic representation
- exact rational round-trip
- deterministic equality/hash
- ledger replay validation
- ledger tamper detection
- candidate contamination remains provisional
- candidate promotion is blocked
- at least one falsifiability condition is present on a sealed research candidate
- corpus manifests agree with executable package behavior
- one classical relation and the canonical relativistic mass-energy relation derive/reproduce deterministically
- corpus/framework status is carried as data and never computed as physical truth

## Final instruction

Treat `PHYSICS_COMPILER_MVP_SPECS.md` as the contract.

Your job is to transform that contract into a **sequenced engineering plan with no architectural ambiguity and no unnecessary scope**.

Do not write implementation code.

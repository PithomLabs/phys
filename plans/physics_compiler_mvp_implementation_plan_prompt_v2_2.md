# Physics Compiler MVP — Implementation Plan Prompt v2.2

## 0. Handoff-critical instruction

**VERY IMPORTANT:** this repository's Go module path MUST be exactly:

```text
github.com/PithomLabs/phys
```

The implementation plan MUST use that module path everywhere. A different module path is a hard failure. All imports within the repository MUST resolve under `github.com/PithomLabs/phys/...`.

The companion file:

```text
physics_compiler_mvp_specs_v2_2.md
```

is the **normative implementation specification**. It is authoritative. Do not redesign it, reinterpret it, or introduce alternative architectures.

Your role is **implementation architect, not implementer**. You may include exact type/function signatures and small structural sketches, but do not write implementation bodies.

The target state is that a separate coding agent can implement the MVP without making architectural decisions.

This v2.2 pair is the final pre-handoff contract after the multi-agent review cycle. Do not reopen settled architecture. The required output is an executable implementation plan, not a new architecture proposal.

---

## Mission

Plan the smallest coherent Go implementation that demonstrates this thesis:

> An AI agent can read a machine-readable physics corpus, construct typed physical primitives, perform symbolic derivations under explicit assumptions, have MRC reject structurally invalid operations, record committed reasoning with provenance, formulate a provisional hypothesis, and hand a reviewable research artifact to humans — without the system adjudicating physical truth.

The AI is the theorist.

The Physics Compiler is the formal symbolic pen-and-paper substrate.

Human researchers and empirical reality remain the final scientific authority.

The MVP proves the architecture, not the scientific completeness of physics.

---

## Six product requirements

The implementation plan MUST cover all six:

1. **Fundamental primitives + MRC:** typed physical primitives, constructor authority, and operation-time MRC enforcement.
2. **Frameworks as evaluation targets:** mechanics and special relativity are scoped frameworks with explicit assumptions, limits, provenance, anomalies, and derivation targets; neither is the final ontology of nature.
3. **Machine-readable corpus:** package code plus validated semantic manifests are the AI-readable physics corpus.
4. **Collective limitations are explicit:** dimensions, physical kinds, assumptions, conventions, provenance, regimes, anomalies, candidate containment, and human handoff are explicit artifacts.
5. **No truth adjudication:** no truth score, probability of truth, theory ranking, automatic corpus-status inference, or empirical adjudication.
6. **Hypothesis formation:** AI-callable candidate concepts and provisional ResearchCandidates exist, with falsifiability structure and human-only scientific promotion/adjudication.

---

## Non-negotiable architecture

Preserve every decision below exactly as specified:

- Go is the implementation language.
- Module path is `github.com/PithomLabs/phys`.
- Go standard library only for MVP.
- No `.phys` language.
- No lexer/parser/frontend.
- No custom compiler language.
- No general mathematical ontology.
- No full CAS.
- No theorem prover.
- No numerical execution or simulation.
- No empirical-data ingestion.
- No automatic empirical validation.
- No truth scores or theory rankings.
- No automatic hypothesis promotion.
- No EBP 2.1 integration.
- MRC has one semantic source of truth in the normative specs and core operation contracts.
- `physvet` is deferred from MVP.
- `core.Object` is the common immutable carrier, but its actual authoritative representation lives behind the module's `internal/kernel` boundary; public `core` exposes the stable API.
- Domain packages expose thin nominal wrappers around `core.Object` and exact `CoreObject() core.Object` accessors.
- Pure symbolic operations live in `ops`; ledger/session orchestration lives in `session`; there is no `core -> ops` dependency.
- `Identify` is a session method, not a package-level `ops.Identify`, and has no ambient/global session state.
- Candidate concepts are open in representation but closed in trusted authority.
- `APPROXIMATED` is reserved for post-MVP work; MVP operations do not produce it.
- MRC is versioned and fallible; MVP has no runtime bypass or exception override API.
- External `ResearchCandidate` JSON loads only into `UnverifiedResearchCandidate`; only the `session` package may yield a trusted sealed `ResearchCandidate`: `Session.Seal()` is the canonical minting action, while `UnverifiedResearchCandidate.Validate()` is an external-loading entrypoint that MUST delegate to the same private session sealing/validation path; unverified parsing exposes no `core.Object` values.

Do not re-open any of these decisions in the implementation plan.

---

## Scope of MVP

Only these populated implementation areas are permitted:

```text
core/
ops/
session/
mechanics/
relativity/
hypothesis/
docs/
internal/kernel/
```

Do not add stub packages for:

```text
electromagnetism
qm
qft
statmech
review
physvet
```

Special relativity is the only relativity content in MVP.

---

# Required implementation-plan output

Return one Markdown implementation plan with these exact sections, in this exact order.

## 1. Executive architecture

Explain the final dependency graph and authority boundaries exactly as specified:

```text
internal/kernel  ← actual immutable object/Expr authority
       ↑
core             ← public API/type aliases
       ↑
ops              ← pure transformations
       ↑
session         ← ledger/replay/seal authority

mechanics / relativity / hypothesis → core + internal/kernel
session MUST NOT import mechanics, relativity, or hypothesis
```

Make the absence of `core -> ops` explicit.

## 2. Exact repository tree

Show every file the implementation will create or modify using the normative tree in the specs.

Do not add convenience files. Tests MUST use adjacent `*_test.go` files.

## 3. Dependency-ordered implementation sequence

Give an ordered sequence. Every step MUST contain:

- purpose
- dependencies
- exact files affected
- implementation result
- verification gate

The sequence MUST begin with the immutable core/kernel representation and canonicalization, then metadata/provenance, then pure operations, then corpus packages, then session/ledger/replay, then hypothesis/research artifacts, then end-to-end acceptance.

## 4. Core model implementation

Describe the exact implementation of:

```text
internal/kernel.Object
core.Object
core.Expr
Kind
Dimension
AssumptionSet
ConventionSet
Provenance
CorpusStatus
```

Explain how the internal authority boundary prevents external construction of trusted `core.Object` values while still allowing domain and operation packages inside the module to mint validated immutable objects.

## 5. MRC implementation map

Include this table:

| Rule | Exact semantic check | Enforcement location | Failure | Test |
|---|---|---|---|---|

Cover `MRC-001` through `MRC-008` exactly.

Distinguish:

- Go package/internal visibility
- trusted constructor/mint authority
- operation-time MRC
- session authority
- candidate containment
- deferred `physvet`

## 6. Symbolic engine

Describe only the exact MVP node set:

```text
Symbol
Rational
Add
Mul
Neg
Pow
Sqrt
Call
Relation
BranchSet
```

Describe canonicalization, exact `math/big.Rat`, deterministic child sorting, canonical equality, SHA-256 hashing, canonical JSON, the constructor-vs-`Simplify` normalization boundary, the exact rational `Pow` rules needed by the golden traces, the bounded differentiation rules, the substitution-based limit rules, the fixed `lorentz_factor` body, the pattern-only `Solve`, and the narrow `SelectBranch` behavior.

Do not propose a general CAS, root finder, theorem prover, parser, or symbolic function registry.

## 7. Physics corpus

Describe:

- `mechanics/manifest.json`
- `relativity/manifest.json`
- `go:embed`
- strict typed manifest structs
- `constructor` mapping
- canonical expression decoding
- constructor cross-check
- provenance/dimension/kind/assumption comparison
- anomaly/limitation records

The manifest `statement` field is explanatory text only. It MUST NOT be parsed as mathematics.

## 8. Assumptions, conventions, provenance, and candidate containment

State the exact merge and conflict rules, deterministic limited assumption entailment used only for sign/nonzero preconditions, provenance propagation law, corpus-status separation, and hypothesis contamination law.

Do not invent a provenance ordering that is not in the specs.

## 9. Derivation ledger

Describe:

- `session.Session`
- exact nine session actions
- draft buffer
- commit/seal state machine
- pure `ops` execution
- session-only `Identify`
- canonical `InputCanonicals`
- canonical `OutputCanonical`
- operation parameters
- genesis hash
- `StepEnvelope`
- SHA-256 chain
- replay
- tamper detection
- integrity vs authenticity limitation
- non-candidate derivations stop at `Concluded`; `Seal` fails with `ProvenanceError` unless valid `HYPOTHESIS` metadata is present
- external ledger/candidate parsing is unverified and cannot manufacture trusted public objects

Explicitly explain how replay avoids an import cycle.

## 10. Hypothesis and ResearchCandidate

Describe:

- hypothesis candidate concept construction
- required explicit `Kind` and `Dimension`
- fixed `HYPOTHESIS` provenance
- ResearchCandidate trusted construction is owned exclusively by the `session` package; `Session.Seal()` is the canonical minting action, and `UnverifiedResearchCandidate.Validate()` MUST delegate to the same private session sealing/validation path rather than implement an independent trusted-construction path.
- candidate contamination
- Prediction
- FalsificationCondition
- RecoveryClaim
- AnomalyReference
- Review / Challenge artifacts
- external reviewer boundary
- validation means artifact integrity, not physical truth

## 11. Canonical vertical slices

Describe the exact implementation/test sequence for:

- typed mechanics
- `F = ma`
- dimension rejection
- physical-kind rejection
- assumption conflict
- convention conflict
- nontrivial differentiation of `1/2 m v²`
- special-relativistic energy-momentum relation
- `E = mc²` derivation by `Substitute → Simplify → Solve → SelectBranch`
- `lim(v→0) LorentzFactor(v) = 1`
- explicit `Identify` firewall
- hypothesis contamination
- falsifiability + anomaly
- manifest validation
- exact-rational round-trip
- ledger tamper/replay
- determinism
- sealed candidate handoff

Do not replace the normative vertical slices with alternatives.

## 12. Coverage matrix

This section is mandatory.

Map **every normative `MUST`/`MUST NOT` clause in the specs** and every MRC rule ID to:

```text
Requirement ID
Implementation location
Primary test
Secondary verification where applicable
```

Do not use vague rows such as "covered by integration tests".

The matrix is a first-class planning artifact, not an afterthought.

## 13. Acceptance gate

List the exact pass conditions, including:

```text
go test ./...
go vet ./...
```

and all normative A–S acceptance tests.

A plan is not accepted if it invents architecture, adds prohibited scope, or leaves a normative requirement without a test mapping.

## 14. Explicit deferrals

Repeat all v0.5+ deferrals exactly from the specs, especially:

- `physvet`
- general series expansion
- full 1905 Einstein derivation
- broader symbolic analysis
- integration
- richer tensor/index machinery
- GR/QM/QFT/stat-mech/EM
- numerical execution
- empirical adapters
- reviewer orchestration
- MRC exception workflow
- automatic corpus governance
- automatic hypothesis promotion
- EBP 2.1

## 15. Scope accounting

State:

- exact file count: 39 including `go.mod`
- exact split: 24 Go source files, 10 Go test files, 2 JSON manifests, 3 config/documentation files
- major exported types
- major operations
- expected dependency count

The normative tree is intentionally below the 40-file budget. Do not exceed it.

---

## Planning rules

The implementation plan MUST NOT:

- propose multiple architectures
- reopen settled decisions
- substitute another object abstraction
- move ledger/session authority back into `core`
- introduce a public generic `NewObject`
- introduce a public `ops.Identify`
- introduce global/ambient session state
- add a `.phys` parser
- turn `physvet` into MVP scope
- add deferred domain packages
- build a general-purpose CAS
- build a general root solver
- build a theorem prover
- build numerical execution
- use mathematical source strings as executable input
- use `map[string]any` for canonical artifacts
- use floats for exact symbolic/dimensional values
- silently omit a normative requirement
- use phrases such as "the coding agent can decide"

The plan MAY include type signatures, struct field lists, dependency diagrams, exact algorithm steps, and test pseudocode. It MUST NOT include implementation bodies.

If the specification appears to contain a genuine contradiction, do not silently resolve it. Put the precise issue under an unnumbered `Open Spec Items` note after section 15. For this v2.2 specification the expected value is:

```text
Open Spec Items: NONE
```

A planner that reports architectural uncertainty here has failed the handoff goal.

---

## Anti-overengineering guardrails

The following are binding:

### Solver

`Solve` is pattern matching only. MVP accepts exactly the energy-quadratic pattern:

```text
Relation(eq, Pow(Symbol, 2), Expr)
```

with one target symbol, and returns:

```text
BranchSet(+Sqrt(Expr), Neg(Sqrt(Expr)))
```

Any other form returns `UnsupportedOperationError`.

### Calculus

`Differentiate` is bounded to the exact recursive rules specified for constants, symbols, Add, Mul, Neg, and integer-power `Pow`. `Call`, `Sqrt`, `Relation`, and `BranchSet` differentiation are unsupported in MVP.

`Limit` is direct substitution plus simplification. The only `Call` case is the fixed `lorentz_factor` definition.

### Function definition

The `lorentz_factor` definition is fixed to:

```text
1 / Sqrt(1 - Pow(v/c, 2))
```

There is no general runtime function-definition registry.

### Thin wrappers

Domain wrappers contain one `core.Object` and expose `CoreObject() core.Object`. All generic operations consume `core.Object`. Do not duplicate symbolic logic in domain packages.

### Manifest validation

Manifests use canonical expression JSON, never equation strings for executable semantics. Each item carries a `constructor` identifier, and tests use a static function map defined in the test package.

### Ledger replay

Every committed step retains the canonical inputs and canonical output needed for replay. Hashes verify integrity; retained canonical forms make replay possible. Do not claim that SHA-256 alone provides authenticity.

### External artifact parsing

`ParseResearchCandidateJSON` returns only `UnverifiedResearchCandidate`. It MUST NOT return `ResearchCandidate` and MUST NOT expose `core.Object` values from externally supplied JSON. Validation must be explicit before a sealed candidate is obtained.

---

## Final constraint

The plan succeeds only when an independent coding agent can implement the MVP with **no architectural design decisions left open**.

`session` MUST NOT import `hypothesis`, `mechanics`, or `relativity`; those packages supply objects to the session through `core.Object`.

The intended proof loop is:

```text
AI reads corpus
      ↓
typed physical construction
      ↓
assumptions + conventions
      ↓
pure symbolic operations
      ↓
MRC enforcement
      ↓
session ledger
      ↓
explicit identification when needed
      ↓
replayable derivation
      ↓
provisional hypothesis
      ↓
falsifiability + anomaly context
      ↓
sealed ResearchCandidate
      ↓
human review
```

The system stops there.

It does not decide what nature believes.

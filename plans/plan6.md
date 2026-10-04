# Physics Compiler MVP — Normative Implementation Specification

**Status:** Normative MVP specification  
**Audience:** Coding/implementation agent  
**Purpose:** Remove architectural ambiguity so implementation is mechanical rather than interpretive.

---

## 0. Authority and interpretation

This document is the **authoritative implementation contract** for the MVP.

The implementation plan prompt is subordinate to this document.

The implementation agent MUST NOT redesign the architecture, introduce alternatives, or choose among competing approaches unless a requirement here is literally impossible in the existing repository.

The six product requirements are:

1. The library MUST provide fundamental physical primitives capable of representing physics equations, with MRC enforced through the Go type system and library contracts.
2. Physics frameworks such as relativity MUST be represented as explicitly scoped corpus/evaluation targets with assumptions and limits; common equations such as `F = ma` and `E = mc²` MUST be reusable corpus artifacts without making any framework the final ontology of nature.
3. Packages and machine-readable metadata MUST form the AI-readable physics corpus.
4. The architecture MUST expose assumptions, provenance, limitations, anomalies/failures, and reviewable derivations.
5. The library MUST NOT adjudicate physical truth.
6. The system MUST enable provisional AI-generated hypotheses and research candidates while stopping before scientific/empirical adjudication.

---

# 1. MVP thesis

The MVP is a **formal symbolic pen-and-paper substrate** for AI agents.

The AI performs the conceptual physics reasoning.

The system provides:

- typed physical primitives
- symbolic expressions
- symbolic operations
- MRC constraints
- assumption propagation
- provenance
- canonicalization
- derivation ledger
- machine-readable corpus
- provisional hypothesis objects
- human-review artifacts

The system does not:

- simulate nature
- fit measurements
- perform numerical optimization
- determine physical truth
- rank theories by truth
- automatically promote hypotheses
- integrate EBP 2.1

---

# 2. Scope of the MVP

## 2.1 Publicly populated packages

The MVP MUST populate only:

```text
core/
ops/
mechanics/
relativity/
hypothesis/
cmd/physvet/       # NOT implemented in MVP
```

There MUST NOT be stub directories for:

```text
electromagnetism/
qm/
qft/
statmech/
```

Those are deferred.

## 2.2 Static analyzer decision

A custom `physvet` analyzer is **deferred from the MVP**.

Do not implement `cmd/physvet/`.

For MVP, "compiler-level MRC" means:

- nominal Go types for distinct physical categories
- constructor-controlled object creation
- unexported authoritative fields
- package visibility boundaries
- interfaces with unexported methods where needed to prevent external impostors

MRC checks that cannot be expressed by Go's static type system are enforced by the library operation contracts.

A future `physvet` MUST consume the same library contracts and metadata rather than becoming a second semantic source of truth.

## 2.3 MVP canonical physics content

### Mechanics

Populate only enough to represent and test:

- `Mass`
- `Time`
- `Position`
- `Velocity`
- `Acceleration`
- `Force`
- `Momentum`
- `Energy`
- `KineticEnergy`
- `NewtonSecondLaw`
- `MomentumRelation`
- `KineticEnergyRelation`

### Special relativity

Populate only:

- `Spacetime`
- `MinkowskiMetric`
- `RestMass`
- `Energy`
- `ThreeMomentum`
- `FourMomentum`
- `SpeedOfLight`
- `LorentzFactor`
- `EnergyMomentumRelation`
- `RestFrame`
- `MassEnergyRelation`

Do not implement general relativity in MVP.

The `relativity` package therefore means **special-relativity corpus content only** for this release.

---

# 3. Repository layout

Create this exact structure:

```text
core/
    object.go
    dimension.go
    expr.go
    assumption.go
    convention.go
    provenance.go
    corpus.go
    ledger.go
    research.go
    review.go
    errors.go
    canonical.go

ops/
    arithmetic.go
    simplify.go
    substitute.go
    differentiate.go
    limit.go
    identify.go
    compare.go
    solve.go

mechanics/
    primitives.go
    relations.go
    manifest.json
    manifest_test.go
    relations_test.go

relativity/
    primitives.go
    relations.go
    manifest.json
    manifest_test.go
    derivation_test.go
    anomaly_test.go

hypothesis/
    candidate.go
    candidate_test.go

docs/
    paper-translation/
        common-notation.md
        framework-mapping.md
        ambiguity-resolution.md
```

No other top-level package is required for the MVP.

The implementation MUST preserve the repository's existing Go module path and existing Go toolchain declaration. Do not invent a new module path.

---

# 4. Core object model

## 4.1 Public physical types

The following MUST be distinct nominal Go types. Do NOT alias all of them to one exported `Object` type:

```text
Mass
Time
Position
Velocity
Acceleration
Force
Momentum
Energy
Spacetime
MinkowskiMetric
RestMass
ThreeMomentum
FourMomentum
```

Each domain type MUST contain immutable authoritative state originating from `core`.

Strings may be used only as **symbol names**.

Strings MUST NEVER be used as carriers for relational physical structure.

Valid pattern:

```go
mechanics.Mass("m")
```

because `"m"` is only the symbolic name.

Required relational constructors MUST use typed carriers:

```text
MinkowskiMetric(Spacetime)
FourMomentum(Spacetime)
Metric(Spacetime)
State(System)    # future, not MVP
```

Do not accept arbitrary string identifiers where a physical carrier is required.

## 4.2 Constructor authority

Authoritative fields MUST be unexported.

AI/caller code MUST NOT be able to instantiate a trusted physical object by filling a struct literal.

Domain constructors are the only trusted construction path.

## 4.3 Immutability

Trusted physical objects are immutable after construction.

Operations return new values.

No operation mutates an existing physics object, expression, assumption set, provenance record, or ledger step.

---

# 5. Dimensions

## 5.1 Representation

Use seven SI base dimensions:

```text
M   mass
L   length
T   time
I   electric current
Θ   thermodynamic temperature
N   amount of substance
J   luminous intensity
```

Represent exponents exactly using `math/big.Rat`.

Do NOT use floating-point values for dimensional exponents.

## 5.2 Required operations

`Dimension` MUST support:

- equality
- multiplication/addition of exponents
- division/subtraction of exponents
- integer/rational power
- deterministic canonical serialization

## 5.3 MVP dimension behavior

Required physical dimensions:

```text
Mass                  M
Time                  T
Position              L
Velocity              L T^-1
Acceleration          L T^-2
Force                 M L T^-2
Momentum              M L T^-1
Energy                M L^2 T^-2
```

`SpeedOfLight` has:

```text
L T^-1
```

## 5.4 Dimension MRC

Addition/equality of incompatible dimensions MUST fail.

Multiplication and division MUST combine dimensions.

Power MUST transform dimensions using exact rational exponents.

---

# 6. Symbolic expression engine

## 6.1 Closed representation

Implement one exported immutable `core.Expr` handle backed by a **closed internal expression tree**.

Do not expose an extensible public expression interface.

The internal node set MUST be limited to:

```text
Symbol
Rational
Add
Mul
Neg
Pow
Sqrt
Relation
Derivative
Limit
Function
```

No generic arbitrary callback nodes.

No reflection-based symbolic representation.

No dependency on a general external CAS.

## 6.2 Exact scalar representation

Use:

```text
math/big.Rat
```

for exact rational coefficients.

Physical constants such as `c` are symbolic typed quantities, not measured floating-point values.

## 6.3 Canonicalization

Canonicalization MUST be deterministic.

At minimum:

- flatten nested addition
- flatten nested multiplication
- sort commutative operands deterministically
- combine exact rational coefficients
- normalize negative signs
- normalize integer/rational powers
- normalize relation operand ordering where semantically appropriate
- normalize dimensions
- preserve non-commutative order where required

Never use expression strings as the equality mechanism.

## 6.4 Equality and hashing

Expose deterministic public operations:

```text
Eq(a, b)
Hash(expr)
```

`Hash` MUST be based on canonical structure, not display text.

Use SHA-256 for canonical hashes.

The canonical serialization MUST be stable across runs.

---

# 7. Assumption system

## 7.1 First-class assumptions

Every derived symbolic result MAY carry an immutable `AssumptionSet`.

Assumptions MUST NOT be represented only as comments.

Minimum assumption kinds:

```text
Domain
Regime
Constraint
Convention
Approximation
MathPrecondition
PhysicalAssumption
```

## 7.2 Assumption structure

Each assumption contains:

```text
Kind
Key
Value
```

`Key` identifies the logical subject.

Examples:

```text
(Regime, velocity_relation, v << c)
(Constraint, momentum, p = 0)
(Convention, metric_signature, -+++)
(MathPrecondition, variable, x >= 0)
(Approximation, series_order, 2)
```

## 7.3 Merge policy

When an operation combines inputs:

```text
result assumptions =
    union(input assumptions)
    + required operation assumptions
    + newly introduced assumptions
```

Duplicate identical assumptions collapse to one.

Two different values for the same `(Kind, Key)` MUST produce an `AssumptionConflictError`.

The MVP MUST NOT attempt theorem-proving or logical subsumption.

Do not infer that one inequality or regime is more specific than another.

---

# 8. Conventions

Conventions are first-class immutable metadata.

Examples:

```text
metric signature
index position convention
coordinate convention
sign convention
```

If an operation receives incompatible convention values for the same convention key, it MUST return a `ConventionConflictError`.

Convention metadata MUST participate in provenance and canonical identity.

---

# 9. Provenance and status

## 9.1 Provenance status

The library owns and mints these statuses:

```text
DEFINED
POSTULATED
DERIVED
IDENTIFIED
APPROXIMATED
HYPOTHESIS
```

## 9.2 Corpus status

Human-curated corpus status is separate:

```text
ESTABLISHED
CONTESTED
SUPERSEDED
FALSIFIED
```

The library MAY carry, render, and serialize corpus status.

The library MUST NOT infer or compute corpus status.

There MUST be no field named `Truth`, `TruthScore`, `ConfidenceInReality`, or equivalent.

## 9.3 Provenance fields

A provenance record MUST retain:

```text
Status
Source
Framework
Parent hashes / premise references
Assumption hash
Convention hash
MRC version
```

---

# 10. MRC

## 10.1 MRC rule identifiers

All MRC failures and committed derivation steps MUST identify the ruleset version.

Use:

```text
MRC version: "mrc-v0.4"
```

Minimum stable rule IDs:

```text
MRC-001  constructor/carrier integrity
MRC-002  dimensional compatibility
MRC-003  physical/category compatibility
MRC-004  assumption compatibility
MRC-005  convention compatibility
MRC-006  explicit physical identification
MRC-007  provenance/session authority
MRC-008  candidate containment
```

These IDs are stable identifiers, not an exception mechanism.

## 10.2 Fallible MRC

MRC itself is not treated as metaphysical truth.

The MVP MUST preserve enough versioning/provenance that a future human-validated MRC revision can explicitly supersede a prior rule.

Do NOT implement a runtime exception/override registry in MVP.

Do NOT provide:

```text
DisableMRC()
IgnoreRule(...)
ForceValid(...)
```

or equivalent bypasses.

Future exceptions are expected to enter through a versioned ruleset revision, with their own human evidence/review process, outside the MVP.

## 10.3 Failure behavior

MRC violations MUST return structured Go errors.

The library MUST NOT panic for normal MRC violations.

Required categories include:

```text
DimensionMismatchError
CategoryMismatchError
AssumptionConflictError
ConventionConflictError
IdentifyError
ProvenanceError
CandidateContainmentError
```

All must satisfy the standard `error` interface and be inspectable through `errors.As`.

---

# 11. Operation contracts

All operations MUST:

1. validate applicable MRC rules
2. validate dimensions where relevant
3. validate assumptions/conventions where relevant
4. construct a new immutable result
5. attach provenance
6. be recordable by a derivation session

No operation may silently bypass the session when producing a committed derivation.

## 11.1 Required MVP operations

Implement only:

```text
Add
Subtract
Multiply
Divide
Pow
Sqrt
Simplify
Substitute
Differentiate
Limit
Compare
Solve
Identify
```

Do NOT implement `Integrate` in MVP.

## 11.2 Simplify vs Identify

`Simplify` is purely mechanical.

It MUST NOT:

- infer physical identity
- rename a quantity into another quantity
- assert equivalence because expressions "look similar"
- promote a result to `IDENTIFIED`

`Identify` is the only MVP operation that records a physical identification.

Required semantic contract:

```text
Identify(a, b, justification)
```

MUST:

- accept two explicit physical operands
- require a non-empty recorded justification
- check dimensions/structure before creating the identification
- create an `IDENTIFIED` provenance record
- create a first-class identification artifact
- append the identification to the derivation ledger
- never be called implicitly by `Simplify`

The justification is the AI's stated physical reason and is reviewable evidence, not proof of truth.

---

# 12. Derivation session and ledger

## 12.1 Session authority

A `core.Session` is the only authority that can mint a committed `Derivation`.

A `Derivation` MUST NOT have a public ordinary constructor.

Domain objects and expressions may exist independently, but trusted derivation provenance is session-owned.

## 12.2 Minimum session actions

The session MUST support:

```text
Postulate
Declare
Define
Step
Identify
Conclude
Draft
Commit
Seal
```

`Draft` operations do not contribute to the committed ledger until explicitly committed.

## 12.3 Ledger step contents

Every committed step MUST contain:

```text
StepID
Index
Label
Operation
Input hashes
Output hash
Assumption hash
Convention hash
Provenance status
MRC version
Previous step hash
Current step hash
```

## 12.4 Hash chain

Compute:

```text
currentHash =
    SHA256(canonical(step_without_currentHash) + previousHash)
```

Do not include timestamps or nondeterministic process data in canonical step hashing.

## 12.5 Replay validation

A committed derivation MUST be replayable.

`Validate()` MUST:

1. verify the hash chain
2. resolve each prior canonical output
3. re-execute each mechanical step
4. reproduce each recorded canonical output
5. verify assumption and convention closure
6. verify every `Identify` still has explicit justification
7. verify final conclusion matches the recorded output

Any mismatch MUST invalidate the derivation.

This is the primary anti-tampering mechanism for the MVP.

---

# 13. Canonical physics-corpus manifests

## 13.1 Canonical source

Each populated domain package MUST contain:

```text
manifest.json
```

The manifest is the canonical machine-readable semantic source for corpus metadata.

Go doc comments MAY mirror the manifest but MUST NOT be the canonical metadata source.

## 13.2 Manifest minimum schema

Each manifest MUST contain:

```json
{
  "schema_version": "0.4",
  "framework_id": "...",
  "framework_name": "...",
  "corpus_status": "...",
  "assumptions": [],
  "domain": [],
  "limits": [],
  "items": []
}
```

Each item MUST support:

```text
id
kind
name
statement
provenance_status
source
assumptions
derivable_from
reduces_to
known_limits
anomalies
falsification_conditions
```

Not every field is required for every item, but the schema MUST be explicit.

## 13.3 Runtime access

Each package MUST expose:

```text
Manifest()
```

The JSON should be loaded using Go's embedded-file mechanism.

Do not parse Go source code or doc comments at runtime.

## 13.4 Manifest validation

Manifest tests MUST cross-check mechanically verifiable claims.

Examples:

- manifest dimensions vs actual constructor dimensions
- manifest relation identity vs actual canonical expression
- declared reduction target vs executable reduction test
- referenced item IDs actually exist

A manifest inconsistency MUST fail tests.

---

# 14. Mechanics corpus

The mechanics manifest MUST contain at least:

```text
NewtonSecondLaw
MomentumRelation
KineticEnergyRelation
```

with canonical expressions equivalent to:

```text
F = ma
p = mv
K = 1/2 mv²
```

These are corpus artifacts.

Their presence MUST NOT be interpreted by the compiler as metaphysical declarations of final truth.

---

# 15. Special-relativity corpus

The relativity manifest MUST explicitly identify the framework as:

```text
special_relativity
```

It MUST carry framework assumptions and limits.

Minimum framework-level assumptions:

```text
Minkowski spacetime
Lorentz symmetry
special-relativistic regime
no gravitational dynamics in this MVP framework
```

Do not claim that this metadata constitutes proof that the framework is universally true.

Minimum canonical relations:

```text
LorentzFactor
FourMomentum
EnergyMomentumRelation
MassEnergyRelation
```

with canonical mass-energy relation equivalent to:

```text
E = m c²
```

under the stated rest-frame condition.

---

# 16. Framework evaluation anatomy

A framework is considered an MVP evaluation target only if it has:

1. at least one deterministic derivation test
2. at least one executable reduction/limit test
3. at least one known limitation/anomaly record
4. at least one falsifiability-related field where applicable

A text field saying "reduces to classical mechanics" is NOT sufficient.

The reduction MUST be executable.

---

# 17. Canonical MVP derivations

## 17.1 Classical derivation

The first canonical test MUST construct and reproduce:

```text
F = ma
```

through typed mechanics primitives.

This proves basic physical construction and category/dimension enforcement.

## 17.2 Relativistic derivation

The primary relativistic integration test MUST derive:

```text
E = m c²
```

from the energy-momentum relation:

```text
E² = (pc)² + (mc²)²
```

under:

```text
p = 0
m >= 0
c > 0
```

The derivation MUST:

1. construct typed mass, momentum, energy, and speed-of-light quantities
2. construct the energy-momentum relation
3. apply the rest-frame substitution
4. simplify
5. take the positive-energy root under an explicit precondition
6. produce the canonical result `E = mc²`
7. record every step

The MVP does NOT implement Einstein's full 1905 series-expansion derivation.

The 1905 historical derivation, tagged series truncation, and broader series algebra are deferred to v0.5.

This is intentional scope control.

---

# 18. Identify micro-test

Because the physical-identification boundary is core to the architecture, the MVP MUST include one explicit `Identify` test separate from the mechanical mass-energy derivation.

The test MUST prove:

- `Simplify` cannot create an `IDENTIFIED` step
- `Identify(a,b,justification)` creates one
- the justification is recorded
- the identification is visible in the review artifact
- dimensions/structural compatibility are checked
- removing the identification leaves the derivation incomplete when that identification is required

The test need not assert that the physical identification is true.

It asserts only that the system treats it as an explicit, reviewable physical claim.

---

# 19. Hypothesis layer

## 19.1 Candidate objects

A candidate physical concept MUST be representable without changing trusted core object definitions.

Use a distinct candidate provenance/status:

```text
HYPOTHESIS
```

Candidate objects are provisional.

Candidate-derived results MUST remain candidate/provisional.

They MUST NOT mint:

```text
ESTABLISHED
DERIVED trusted corpus
```

or equivalent trusted corpus authority.

## 19.2 Closed authority

The MVP MUST implement:

> open representation, closed authority

The AI may introduce a candidate concept.

Only trusted constructors and sessions can mint trusted corpus provenance.

There is no AI-callable promotion operation.

## 19.3 Candidate acceptance constraints

A sealed `ResearchCandidate` MUST contain:

```text
hypothesis
premises
assumptions
derivation
framework dependencies
predictions
at least one falsification condition
recovery/reduction claims where applicable
anomaly references where applicable
```

A candidate missing a required falsification condition MUST NOT be sealed as a complete research candidate.

---

# 20. Popperian/falsifiability structure

The MVP does not compute truth.

It does support explicit falsifiability structure.

A `FalsificationCondition` MUST identify:

```text
target claim
condition/observation that would contradict it
scope/regime
```

A `Prediction` MUST identify:

```text
observable/quantity
predicted relation
assumptions
framework/candidate source
```

The library MUST not attach a probability of being true.

---

# 21. Anomaly/failure registry

Each populated framework MUST contain at least one structured limitation/anomaly/failure record.

For MVP this is metadata, not a numerical empirical database.

Minimum fields:

```text
id
framework
description
status
related_assumptions
related_items
research_relevance
```

The candidate layer MUST be able to cite an anomaly record.

An anomaly record does not cause the compiler to conclude that the underlying framework is false.

---

# 22. ResearchCandidate

## 22.1 Authority

Only a `Session` that has validated and sealed a derivation can create a trusted `ResearchCandidate`.

Do not provide a public constructor that creates a trusted candidate from arbitrary fields.

## 22.2 Minimum logical schema

```text
ID
Hypothesis
Premises
Assumptions
Derivation
FrameworkDependencies
Predictions
FalsificationConditions
RecoveryClaims
AnomalyReferences
ReviewHistory
MRCVersion
LedgerHash
```

## 22.3 Serialization

The artifact MUST be deterministic JSON.

No timestamps or random identifiers may affect canonical artifact hashing.

External JSON MUST be treated as untrusted input until:

- schema validation
- canonicalization
- ledger hash validation
- replay validation

have passed.

---

# 23. Review artifact

The MVP defines review data structures but does NOT implement AI orchestration.

Minimum review structures:

```text
Challenge
Review
```

`Challenge` fields:

```text
StepID
Category
Severity
Description
```

Minimum challenge categories:

```text
CategoryError
DimensionError
AssumptionConflict
ConventionConflict
UnsupportedIdentification
InvalidReduction
ProvenanceProblem
CandidateOverreach
```

`Review` contains:

```text
DerivationID
Challenges
ReviewerNotes
```

The reviewer AI is external to the core library.

---

# 24. Paper-translation documentation

The MVP MUST include minimal documentation for the AI translation path:

```text
docs/paper-translation/common-notation.md
docs/paper-translation/framework-mapping.md
docs/paper-translation/ambiguity-resolution.md
```

These documents MUST explain:

- common physics notation used by the MVP
- mapping to `mechanics` or `relativity`
- ambiguous symbols and how context resolves them
- the distinction between mathematical notation and physical object identity

This documentation is part of the AI corpus, but is not the canonical machine-readable semantic source; `manifest.json` remains canonical metadata.

---

# 25. Required negative tests

The MVP test suite MUST include failing cases for:

1. incompatible dimensions
2. incompatible physical/category operation
3. assumption conflict
4. convention conflict
5. direct construction of trusted object using forbidden fields
6. `Simplify` attempting to create physical identity
7. missing `Identify` justification
8. candidate result attempting trusted promotion
9. corrupted ledger hash
10. modified intermediate ledger output
11. nondeterministic canonical serialization
12. malformed/inconsistent corpus manifest

Each failure MUST produce a structured, inspectable error.

---

# 26. Determinism tests

Repeated runs of the same derivation MUST produce identical:

```text
canonical expression
object hash
step hashes
derivation hash
manifest-derived relation hash
ResearchCandidate canonical JSON
```

No timestamps, map iteration order, pointer addresses, or random IDs may leak into canonical artifacts.

---

# 27. MVP corpus/evaluation tests

The minimum pass suite is:

### Test A — Mechanics

Reproduce `F = ma`.

### Test B — Dimensions

Reject an invalid physical addition/equality.

### Test C — Assumptions

Reject contradictory assumptions.

### Test D — Conventions

Reject incompatible conventions.

### Test E — Relativity

Derive `E = mc²` from the energy-momentum relation under the rest-frame condition.

### Test F — Identify

Demonstrate explicit physical-identification gating.

### Test G — Candidate containment

Create a novel candidate concept, derive with it, verify candidate contamination, and verify promotion is blocked.

### Test H — Falsifiability

Seal a candidate only when a falsification condition is present.

### Test I — Anomaly inspection

Load one anomaly/failure record and attach it to a candidate.

### Test J — Ledger integrity

Corrupt a recorded step and verify replay/hash validation fails.

### Test K — Manifest validation

Change a manifest dimension/relation declaration and verify the package test fails.

---

# 28. Human boundary

The MVP MUST stop at:

```text
ResearchCandidate
        ↓
human / empirical evaluation
```

The library does not:

- decide whether the candidate is true
- decide whether a corpus framework is ultimately true
- declare a theory confirmed
- decide which theory nature chose
- automatically update corpus status
- automatically promote a candidate
- perform empirical adjudication

---

# 29. MRC exception/revision boundary

The MVP handles the possibility that MRC itself is wrong through:

```text
MRC versioning
stable rule IDs
provenance recording
deterministic derivation artifacts
```

The MVP does NOT implement human exception approval.

A future exception MUST conceptually be handled by:

```text
existing MRC rule
        ↓
human-validated exception proposal
        ↓
evidence and falsification record
        ↓
new versioned/scoped rule
        ↓
regression suite
```

This future workflow is outside MVP.

No code path in MVP may simply disable a rule.

---

# 30. Explicit non-goals

Do NOT implement:

- a new `.phys` language
- parser/scanner/frontend
- general mathematical ontology
- public `math.Scalar`, `math.Vector`, `math.Tensor` hierarchy
- general CAS
- numerical physics
- numerical integration
- numerical optimization
- Monte Carlo
- GPU execution
- theorem prover
- general tensor algebra framework
- full differential geometry engine
- general relativity
- quantum mechanics
- QFT
- statistical mechanics
- automated empirical-data ingestion
- truth scores
- theory rankings
- autonomous scientific adjudication
- automatic corpus promotion
- MRC exception registry
- multi-agent orchestration runtime
- EBP 2.1 integration

---

# 31. v0.5+ deferrals

Explicitly defer:

1. `physvet` static-analysis tooling
2. Einstein's complete 1905 series-expansion derivation
3. tagged general series algebra beyond the MVP
4. broader mathematical expression families
5. integration
6. richer tensor/index machinery
7. general relativity
8. quantum mechanics
9. QFT
10. statistical mechanics
11. empirical-data adapters
12. human MRC exception workflow
13. automatic corpus governance tooling
14. reviewer-AI orchestration

---

# 32. Definition of done

The MVP is complete only when:

- the package tree matches the specification
- all trusted object construction is carrier/type controlled
- dimensions are exact
- assumptions are first-class and conflicting assumptions fail
- conventions are first-class and conflicting conventions fail
- symbolic expressions are canonical and deterministic
- equality/hash are structural
- the derivation ledger is session-owned
- the ledger is hash-chained
- ledger replay detects tampering
- `Simplify` never creates physical identity
- `Identify` is explicit and provenance-recorded
- mechanics corpus reproduces `F = ma`
- relativity corpus derives `E = mc²`
- manifests are embedded, queryable, and test-validated
- at least one anomaly/limitation is represented
- candidate objects remain provisional
- candidate promotion is impossible through the MVP API
- sealed research candidates contain assumptions, derivation, predictions, and falsification conditions
- no artifact can report physical truth
- all required negative and determinism tests pass
- no deferred package or deferred subsystem has been added as scaffolding

The finished MVP should demonstrate one narrow but complete loop:

```text
AI reads corpus
    ↓
constructs typed physics primitives
    ↓
states assumptions
    ↓
performs symbolic derivation
    ↓
MRC rejects invalid moves
    ↓
explicit physical insight is separately recorded
    ↓
derivation becomes canonical + replayable
    ↓
AI may formulate a provisional candidate
    ↓
candidate carries predictions/falsification conditions
    ↓
sealed research artifact
    ↓
human judgment
```

That is the complete MVP. Anything not required to make this loop work belongs outside the MVP.



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

# Physics Compiler MVP — Normative Implementation Specification v2

**Status:** NORMATIVE  
**Purpose:** Define the smallest implementation that can be executed without architectural invention.  
**Source language:** Go using the ordinary Go toolchain.  
**MVP mode:** Symbolic only.  
**Truth boundary:** The system does not adjudicate physical truth.

---

# 0. Authority, scope, and terminology

This document is the authoritative contract for the MVP.

The implementation plan is subordinate to this specification.

The coding agent MUST implement the architecture described here exactly unless the existing repository makes a literal requirement impossible. It MUST NOT introduce an alternative architecture merely because another design is more familiar.

The term **physics corpus** means the machine-readable collection of domain packages, manifests, and auxiliary translation documentation.

The term **framework** means a scoped body of physical assumptions, definitions, relations, and derivations such as special relativity.

The term **MRC** means the formal checks that constrain physically meaningful symbolic construction and transformation.

MRC is itself versioned and fallible. The MVP records its version but does not implement an exception/override mechanism.

---

# 1. Product requirements

The MVP MUST satisfy all six requirements:

1. Fundamental physical primitives with MRC enforced by Go typing, constructor authority, and operation contracts.
2. Frameworks such as special relativity represented as explicitly scoped evaluation targets with assumptions, limits, anomalies, and derivation cases.
3. Packages and semantic manifests form a machine-readable physics corpus.
4. Assumptions, provenance, conventions, regimes, anomalies, and limitations are explicit.
5. The system never adjudicates physical truth or emits truth scores/rankings.
6. The AI can formulate provisional candidate concepts and research candidates, with falsifiability structure and human handoff.

---

# 2. Explicit non-goals

The MVP MUST NOT implement:

- `.phys` or any other physics language
- lexer/parser/frontend
- custom compiler
- general mathematical ontology
- full CAS
- theorem prover
- numerical execution
- numerical integration
- numerical optimization
- simulation
- Monte Carlo
- GPU execution
- empirical-data ingestion
- automatic theory ranking
- truth scoring
- physical-truth adjudication
- automatic corpus-status inference
- automatic hypothesis promotion
- MRC bypass/override registry
- MRC exception approval workflow
- `physvet`
- electromagnetism package
- quantum-mechanics package
- QFT package
- statistical-mechanics package
- general relativity
- broad differential geometry
- general tensor algebra
- multi-agent orchestration runtime
- EBP 2.1 integration

---

# 3. Exact MVP repository layout

Create or modify only:

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
    errors.go
    canonical.go
    object_test.go
    dimension_test.go
    expr_test.go
    assumption_test.go
    convention_test.go
    provenance_test.go
    ledger_test.go
    research_test.go

ops/
    arithmetic.go
    simplify.go
    substitute.go
    differentiate.go
    limit.go
    compare.go
    solve.go
    identify.go
    branch.go
    operations_test.go
    negative_test.go

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
    reduction_test.go
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

Go test files MUST be adjacent to the code they test.

All construction-authority negative tests MUST use external test packages where package visibility could otherwise make the test invalid; for example:

```text
package mechanics_test
```

Do not create empty placeholder directories for deferred domains.

---

# 4. Core physical-object abstraction

## 4.1 Two-layer representation

The MVP uses two levels:

1. distinct nominal domain types for API/type clarity;
2. one validated `core.Object` carrier used by generic symbolic operations.

A domain value is a thin immutable wrapper around a `core.Object`.

Required pattern:

```go
type Mass struct {
    object core.Object
}

func (m Mass) CoreObject() core.Object {
    return m.object
}
```

Every populated domain type MUST expose the exact accessor:

```text
CoreObject() core.Object
```

Generic `ops` functions MUST accept `core.Object` values. Domain callers obtain those values only through `CoreObject()`.

The MVP MUST NOT expose a generic `NewObject(kind, ...)` constructor.

Fixed constructors in the populated domain packages are the trusted construction API.

## 4.2 `core.Object`

`core.Object` is a concrete immutable value, not an open user-implementable interface.

Required logical fields:

```text
valid
name
kind
dimension
expr
assumptions
conventions
provenance
```

All authoritative fields MUST be unexported.

The zero value is invalid.

All operations MUST reject an invalid zero `core.Object`.

The exact read-only accessors are:

```go
Valid() bool
Name() string
Kind() Kind
Dimension() Dimension
Expr() Expr
Assumptions() AssumptionSet
Conventions() ConventionSet
Provenance() Provenance
```

These accessors return immutable values/views. No mutator is permitted.

## 4.3 Trusted construction

Only the fixed constructors defined by the populated core/domain APIs may create valid trusted objects.

No public generic object factory is permitted.

Each constructor fixes its own:

```text
Kind
Dimension
Provenance
Assumptions
Conventions
Expression
```

and callers cannot mutate them.

## 4.4 Strings

Strings are permitted only for symbolic names, identifiers, and metadata.

Strings MUST NEVER carry physical relationships or symbolic expressions.

The symbolic engine MUST NOT parse or evaluate mathematical source strings.

# 5. Physical kind

Dimension is not sufficient for physical category checking.

Define a finite MVP `Kind` enum containing at least:

```text
Mass
Time
Position
Velocity
Acceleration
Force
Momentum
Energy
RestMass
ThreeMomentum
FourMomentum
SpeedOfLight
Spacetime
MinkowskiMetric
Expression
Relation
BranchSet
```

The exact enum names MUST be stable.

## 5.1 Why `Kind` exists

Two quantities may have identical dimensions while not being additively compatible.

Example:

```text
Energy
Torque
```

are dimensionally identical but are not automatically the same physical kind.

MRC-003 therefore checks `Kind` independently of dimensional compatibility.

---

# 6. Dimensions

## 6.1 Base dimensions

Use exactly seven SI base dimensions:

```text
M   mass
L   length
T   time
I   electric current
Θ   thermodynamic temperature
N   amount of substance
J   luminous intensity
```

Each exponent MUST use `math/big.Rat`.

Never use floating-point dimensional exponents.

## 6.2 Required dimension values

```text
Mass            = M
Time            = T
Position        = L
Velocity        = L T^-1
Acceleration     = L T^-2
Force            = M L T^-2
Momentum        = M L T^-1
Energy          = M L^2 T^-2
SpeedOfLight    = L T^-1
```

## 6.3 Dimension operations

Implement:

```text
Equal
Multiply
Divide
Pow
Canonical
```

Dimension arithmetic must be exact and deterministic.

---

# 7. Symbolic expression representation

## 7.1 Closed expression tree

`core.Expr` is an immutable symbolic handle backed by a closed set of internal node types.

The implementation MUST use an internal closed node mechanism inside `core`.

The node set is exactly:

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
```

Do not add:

```text
Derivative
Limit
Function
Arbitrary
Eval
Callback
RawString
```

as expression nodes.

`Differentiate` and `Limit` are operations that transform expressions; they do not require dedicated unevaluated derivative/limit nodes in MVP.

`Call` represents a symbolic function application such as:

```text
LorentzFactor(v)
```

## 7.2 Symbol

A `Symbol` node contains only a symbol name.

The symbol name is data, not executable code.

## 7.3 Rational

Use `math/big.Rat`.

Canonical rational form:

```text
num/den
```

where:

- denominator is positive
- numerator/denominator are reduced
- zero is always `0/1`

## 7.4 Add and Mul

`Add` and `Mul` are n-ary internal nodes.

Canonicalization MUST:

- flatten nested same-operation nodes
- combine exact rational coefficients where applicable
- remove additive zero / multiplicative identity where safe
- sort commutative children deterministically
- preserve order only for operations marked non-commutative; MVP contains no non-commutative multiplication node

## 7.5 Neg

Normalize repeated negation:

```text
Neg(Neg(x)) → x
```

## 7.6 Pow

Represent exact powers.

Integer and rational exponents are allowed.

Do not silently simplify an expression in a way that requires unstated sign/domain assumptions.

## 7.7 Sqrt

`Sqrt(x)` is symbolic.

The simplifier MUST NOT transform:

```text
sqrt(x²)
```

into `x`

without a supplied assumption permitting it.

For the MVP positive-energy mass-energy derivation, the relevant sign assumption is carried explicitly.

## 7.8 Call

`Call` has:

```text
function_id
arguments
```

MVP function IDs include:

```text
lorentz_factor
```

Additional arbitrary runtime functions are forbidden.

---

# 8. Canonical symbolic representation

Canonicalization is deterministic.

For `Add` and `Mul`, children are sorted by:

1. node kind ordinal
2. canonical child hash
3. canonical byte sequence as final tie-breaker

Do not use Go map iteration order.

Do not use pointer addresses.

Do not use timestamps.

Do not use random values.

## 8.1 Equality

Expose:

```text
Eq(a, b)
```

which compares canonical symbolic structure plus the relevant physical metadata.

Equality MUST NOT be implemented by comparing display strings.

## 8.2 Hash

Expose:

```text
Hash(x)
```

using SHA-256 over canonical serialization.

The hash MUST be deterministic across processes and runs.

---

# 9. Canonical JSON encoding

The MVP uses deterministic JSON for external artifacts and manifests.

The canonical expression JSON schema is:

```json
{"kind":"symbol","name":"E"}
```

```json
{"kind":"rational","value":"1/2"}
```

```json
{"kind":"add","terms":[ ... ]}
```

```json
{"kind":"mul","factors":[ ... ]}
```

```json
{"kind":"neg","expr": ...}
```

```json
{"kind":"pow","base": ...,"exp": ...}
```

```json
{"kind":"sqrt","expr": ...}
```

```json
{"kind":"call","function":"lorentz_factor","args":[ ... ]}
```

```json
{"kind":"relation","operator":"eq","lhs": ...,"rhs": ...}
```

Objects must use explicit structs, not `map[string]any`, for canonical encoding/decoding.

Canonical JSON object field order MUST be fixed by the encoder.

Rational values MUST be encoded as canonical strings such as `"1/2"`; JSON floating-point values MUST NOT represent exact symbolic rationals.

Round-trip:

```text
Expr → canonical JSON → Expr
```

MUST reproduce identical canonical form and hash.

---

# 10. Assumption system

## 10.1 Assumption kinds

The MVP supports:

```text
Domain
Regime
Constraint
Convention
Approximation
MathPrecondition
PhysicalAssumption
```

## 10.2 Structure

Each assumption contains:

```text
Kind
Key
CanonicalValue
```

`CanonicalValue` MUST have deterministic serialization.

Example:

```text
Kind: Constraint
Key: momentum
CanonicalValue: p=0
```

## 10.3 Merge

For combined inputs:

```text
result =
    union(input assumptions)
    + required operation assumptions
    + newly introduced assumptions
```

Identical assumptions deduplicate.

Two different canonical values with the same `(Kind, Key)` MUST produce:

```text
AssumptionConflictError
```

The MVP does not implement logical subsumption.

Do not infer that:

```text
v << c
```

implies any other regime automatically.

---

# 11. Conventions

Conventions are immutable key/value metadata.

MVP examples:

```text
metric.signature = -+++
```

If the same convention key appears with different values, return:

```text
ConventionConflictError
```

Conventions participate in canonicalization and provenance.

---

# 12. Provenance and corpus status

## 12.1 Provenance status

The library mints:

```text
DEFINED
POSTULATED
DERIVED
IDENTIFIED
APPROXIMATED
HYPOTHESIS
```

## 12.2 Provenance contamination law

When an operation combines provenance-bearing inputs:

- if any input has `HYPOTHESIS`, the output has `HYPOTHESIS`
- otherwise a mechanical result may be `DERIVED`
- approximation-producing operations may produce `APPROXIMATED`
- `Identify` produces `IDENTIFIED`
- no operation may upgrade `HYPOTHESIS` to a trusted provenance state

The rule applies transitively through nested expressions and committed derivation steps.

Candidate provenance cannot be laundered through:

```text
Simplify
Solve
Limit
Compare
SelectBranch
Identify
```

## 12.3 Corpus status

Human-curated corpus status is separate:

```text
ESTABLISHED
CONTESTED
SUPERSEDED
FALSIFIED
```

The library may carry, render, and serialize corpus status.

The library MUST NOT infer or compute corpus status.

There MUST be no API returning:

```text
PHYSICALLY_TRUE
TRUTH_SCORE
PROBABILITY_OF_TRUTH
BEST_THEORY
```

## 12.4 Provenance records

Every committed result retains:

```text
Status
Source
Framework
Parent hashes / premise references
Assumption hash
Convention hash
MRC version
```

# 13. MRC rules

The MVP MRC version is:

```text
mrc-v0.4
```

Each rule has a stable identifier.

## MRC-001 — Constructor/carrier integrity

Valid trusted physical objects can only originate from the fixed core/domain constructors.

Authoritative fields are unexported and there is no generic public object factory.

The zero `core.Object` is invalid.

Enforcement is primarily Go package visibility plus constructor design.

## MRC-002 — Dimensional compatibility

- `Add`/`Subtract`: dimensions must match
- `Compare(..., eq)`: dimensions must match
- `Multiply`: dimensions multiply exactly
- `Divide`: dimensions divide exactly
- `Pow`: dimensions transform by exact rational exponent

Failure:

```text
DimensionMismatchError
```

## MRC-003 — Physical-kind compatibility

Dimension equality does not imply physical compatibility.

- `Add`/`Subtract` require identical `Kind`, except `Expression + Expression`
- `Compare(..., eq)` may compare a named physical quantity with a derived `Expression` when dimensions match
- comparison of two named quantities requires compatible `Kind`
- `Multiply`/`Divide` may produce `Expression` and therefore do not require equal operand kinds
- `Kind` is constructor-assigned and immutable

MRC-003 is distinct from MRC-002.

Failure:

```text
CategoryMismatchError
```

## MRC-004 — Assumption compatibility

Operations merge assumptions according to §10. Conflicting `(Kind, Key)` values fail.

Failure:

```text
AssumptionConflictError
```

## MRC-005 — Convention compatibility

Conflicting convention values for the same key fail.

Failure:

```text
ConventionConflictError
```

## MRC-006 — Explicit physical identification

`Identify` is the only operation that creates an `IDENTIFIED` provenance event.

It requires two explicit operands and a non-empty justification.

`Simplify` MUST NOT invoke it.

Failure:

```text
IdentifyError
```

## MRC-007 — Session/provenance authority

A trusted committed derivation can be minted only by a `Session`.

`Derivation` has no public arbitrary-field constructor.

Committed steps are created only through session operations.

Failure:

```text
ProvenanceError
```

## MRC-008 — Candidate containment

Any output depending on a `HYPOTHESIS` input becomes `HYPOTHESIS`.

No operation may raise candidate provenance to a trusted status.

No promotion API exists in MVP.

Failure:

```text
CandidateContainmentError
```

# 14. Operation contracts

All operations return new immutable values and never mutate inputs.

Applicable MRC checks execute before a result is committed to a session ledger.

All successful results propagate assumptions, conventions, and provenance.

## 14.1 Add

```text
Add(a, b core.Object) (core.Object, error)
```

Checks `MRC-002`, `MRC-003`, `MRC-004`, `MRC-005`, `MRC-008`.

## 14.2 Subtract

```text
Subtract(a, b core.Object) (core.Object, error)
```

Same semantic contract as `Add`.

## 14.3 Multiply

```text
Multiply(a, b core.Object) (core.Object, error)
```

Result kind is `Expression`.

Checks `MRC-004`, `MRC-005`, `MRC-008`.

## 14.4 Divide

```text
Divide(a, b core.Object) (core.Object, error)
```

Result kind is `Expression`.

If `b` is symbolic, add:

```text
MathPrecondition
Key: denominator
Value: denominator != 0
```

Checks `MRC-004`, `MRC-005`, `MRC-008`.

## 14.5 Pow

```text
Pow(base core.Object, exponent *big.Rat) (core.Object, error)
```

Dimensions transform exactly.

Required sign/domain conditions are added as assumptions rather than inferred.

## 14.6 Simplify

```text
Simplify(x core.Object) (core.Object, error)
```

Performs only mechanical symbolic normalization.

It MUST NOT make physical identifications or change physical meaning.

## 14.7 Substitute

```text
Substitute(target, variable, replacement core.Object) (core.Object, error)
```

Replacement dimensions must match the substituted quantity.

Assumptions and provenance propagate.

## 14.8 Differentiate

```text
Differentiate(target, wrt core.Object) (core.Object, error)
```

Symbolic differentiation only.

The result uses the existing expression nodes; there is no `Derivative` node.

Unsupported function derivatives return `UnsupportedOperationError`.

## 14.9 Limit

```text
Limit(target, variable, value core.Object) (core.Object, error)
```

MVP scope:

- direct symbolic substitution for supported analytic forms
- the `LorentzFactor(v)` case
- no L'Hôpital rule
- no series-limit engine

Unsupported or singular forms return `UnsupportedOperationError`.

Required limit:

```text
lim(v→0) LorentzFactor(v) = 1
```

This is a self-consistency/reduction check, not a truth judgment.

## 14.10 Compare

```text
Compare(a, b core.Object, op RelationOperator) (core.Object, error)
```

Relation operators:

```text
eq
neq
lt
lte
gt
gte
```

Result kind is `Relation`.

For `eq`, dimensions must match.

## 14.11 Solve

```text
Solve(relation core.Object, target core.Object) (core.Object, error)
```

MVP scope:

- equality relations only
- one target symbol
- polynomial-style isolation only for forms required by `E = mc²`
- returns `BranchSet`
- no general root solver

Example:

```text
E² = X
```

may yield:

```text
E = +sqrt(X)
E = -sqrt(X)
```

## 14.12 SelectBranch

```text
SelectBranch(branches core.Object, constraint core.Object) (core.Object, error)
```

It:

- selects only a branch compatible with the constraint
- adds the constraint to assumptions
- preserves provenance
- returns a new immutable object

For mass-energy:

```text
E >= 0
```

It never invokes `Identify`.

## 14.13 Identify

```text
Identify(a, b core.Object, justification string) (core.Object, error)
```

It:

- accepts two explicit operands
- requires non-empty justification
- checks applicable structural/dimensional constraints
- creates a relation with `IDENTIFIED` provenance
- records the identification in the active session

The justification is reviewable AI reasoning, not proof of physical truth.

# 15. Session and derivation ledger

## 15.1 Session actions

The MVP session supports exactly nine actions:

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

`Draft` and `Commit` are session state transitions, but both are exposed as explicit session operations.

## 15.2 Step vs Identify

`Identify` is a specialized committed step.

The ledger stores it as:

```text
Operation = "Identify"
StepKind = "Identification"
```

Ordinary `Step` stores:

```text
Operation = operation name
StepKind = "Transformation"
```

## 15.3 Commit

`Commit` moves a draft step into the immutable committed ledger.

After commit, the step cannot be mutated.

## 15.4 Seal

`Seal` closes the session's derivation/research artifact.

After sealing:

- committed steps cannot change
- the final ledger hash is fixed
- a `ResearchCandidate` may be minted if required fields are present

## 15.5 Step structure

Each committed step MUST contain:

```text
StepID
Index
Label
StepKind
Operation
InputHashes
OutputHash
AssumptionHash
ConventionHash
ProvenanceStatus
MRCVersion
PreviousStepHash
CurrentStepHash
```

## 15.6 Hash chain

Do not concatenate variable-length fields directly.

Construct a canonical structured `StepEnvelope` containing:

```text
previousHash
canonicalStep
```

and SHA-256 the canonical serialization of that envelope.

This defines:

```text
currentHash =
SHA256(CanonicalJSON(StepEnvelope{
    PreviousHash: previousHash,
    Step: canonicalStep
}))
```

## 15.7 Replay validation

`Validate()` MUST:

1. verify the first hash
2. verify every subsequent hash chain link
3. reconstruct inputs from canonical hashes
4. re-execute every mechanical step
5. reproduce every recorded output hash
6. verify assumption/convention closure
7. verify `Identify` justifications
8. verify provenance monotonicity
9. verify final conclusion

Any mismatch invalidates the derivation.

---

# 16. Candidate objects

## 16.1 Purpose

Candidate objects allow an AI agent to introduce genuinely new physical concepts.

They are provisional.

## 16.2 Construction

The hypothesis package exposes an AI-callable constructor for a provisional candidate concept.

A candidate receives:

```text
Provenance.Status = HYPOTHESIS
```

## 16.3 Containment

Any result depending on a candidate is also `HYPOTHESIS`.

This is provenance contamination, not a separate trust score.

No operation may upgrade a candidate-derived result to trusted provenance.

## 16.4 Trusted corpus boundary

There is no MVP promotion operation.

A candidate cannot become a trusted corpus item through normal AI execution.

Human curation outside the MVP is the only promotion path.

# 17. Falsifiability structures

The candidate artifact supports:

```text
Prediction
FalsificationCondition
RecoveryClaim
AnomalyReference
```

## 17.1 Prediction

Minimum fields:

```text
id
observable
relation
assumptions
```

## 17.2 FalsificationCondition

Minimum fields:

```text
id
targetClaim
contradictingCondition
regime
```

The library records these structures.

It does not determine whether the condition occurred.

---

# 18. Anomaly registry

Each populated framework MUST contain at least one anomaly/limitation record.

Minimum fields:

```text
id
framework
description
status
relatedAssumptions
relatedItems
researchRelevance
```

MVP examples may be historical or conceptual.

The anomaly itself does not cause automatic corpus-status change.

A candidate may reference the anomaly.

---

# 19. Corpus manifests

## 19.1 Canonical source

Each populated package contains:

```text
manifest.json
```

It is the canonical machine-readable semantic source.

Go doc comments MAY mirror metadata but are not canonical.

## 19.2 Manifest fields

Top-level:

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

Each item supports:

```text
id
kind
name
statement
canonical_expr
provenance_status
source
assumptions
derivable_from
reduces_to
known_limits
anomalies
falsification_conditions
```

`statement` is human-readable.

`canonical_expr` is machine-checkable and uses the canonical JSON expression schema from §9.

## 19.3 Loading

Use `go:embed` and typed Go manifest structs.

Each populated package exposes:

```text
Manifest()
```

Raw `map[string]any` is forbidden for the public manifest API.

## 19.4 Validation

Manifest validation is a pure function over input bytes or typed data.

Production uses embedded bytes.

Tests may pass modified bytes directly to validate corruption/inconsistency.

For executable items:

1. decode `canonical_expr`
2. call the specified domain constructor/relation
3. obtain the object's `Expr()`
4. compare canonical structural equality
5. compare canonical SHA-256 hash
6. compare mechanically verifiable metadata

No mathematical source-string parser is permitted or required.

# 20. Mechanics corpus

Mechanics is the initial classical corpus in MVP and is not required to declare a predecessor reduction target.

The mechanics package contains only:

```text
Mass
Time
Position
Velocity
Acceleration
Force
Momentum
Energy
KineticEnergy
NewtonSecondLaw
MomentumRelation
KineticEnergyRelation
```

Minimum relations:

```text
F = ma
p = mv
K = 1/2 mv²
```

The corpus status is human-curated data.

The compiler does not infer that these relations are metaphysically final.

---

# 21. Special-relativity corpus

The MVP `relativity` package means special relativity only.

It contains:

```text
Spacetime
MinkowskiMetric
RestMass
Energy
ThreeMomentum
FourMomentum
SpeedOfLight
LorentzFactor
EnergyMomentumRelation
RestFrame condition
MassEnergyRelation
```

Framework assumptions MUST include at least:

```text
Minkowski spacetime
Lorentz symmetry
special-relativistic regime
no gravitational dynamics in this package
```

The manifest MUST identify known limitations and an anomaly/limitation record.

---

# 22. Rest frame

`RestFrame` is NOT a physical object.

It is represented as an assumption/constraint:

```text
Kind: Constraint
Key: momentum
Value: p=0
```

It participates in assumption propagation and provenance.

---

# 23. Canonical MVP derivations

## 23.1 Classical relation

Construct and reproduce:

```text
F = ma
```

through typed mechanics primitives.

## 23.2 Differentiation micro-test

Execute:

```text
Differentiate(Velocity(v), Time(t))
```

and verify the expected symbolic derivative structure.

This exercises the generic operation layer.

## 23.3 Mass-energy relation

Construct:

```text
E² = (pc)² + (mc²)²
```

Then perform, in order:

1. `Substitute(p, 0)`
2. `Simplify`
3. `Solve` for `E`
4. `SelectBranch` using `E >= 0`

The final canonical result MUST be:

```text
E = mc²
```

The result MUST arise from the actual operations.

Hardcoding `E = mc²` as the test output is forbidden.

This is a substrate/architecture test, not a claim that the compiler independently established physical truth.

## 23.4 Relativistic self-reduction

Execute:

```text
Limit(LorentzFactor(v), v, 0)
```

and reproduce:

```text
1
```

This is the MVP's executable self-consistency/reduction test.

The meaningful low-velocity recovery of Newtonian kinetic energy requires series expansion and is deferred to v0.5.

# 24. Identify micro-test

A dedicated test MUST demonstrate:

```text
Simplify(a)
```

cannot produce an `IDENTIFIED` provenance event.

Then:

```text
Identify(a, b, justification)
```

must:

- produce `IDENTIFIED`
- add the justification
- append an identification ledger entry
- remain visible in the final research/review artifact

The test validates the firewall, not physical truth.

---

# 25. Candidate contamination test

A candidate concept MUST be created through the hypothesis API.

The test MUST:

1. create a candidate object
2. derive a result using it
3. verify the result remains `HYPOTHESIS`
4. verify the result cannot be converted to trusted provenance
5. verify no promotion API exists
6. verify a sealed trusted corpus artifact cannot be produced from the candidate

This is the primary candidate-containment test.

---

# 26. Falsifiability/anomaly micro-test

The MVP MUST:

1. load the framework anomaly record
2. create a candidate referencing that anomaly
3. attach at least one prediction
4. attach at least one falsification condition
5. seal the candidate as a provisional `ResearchCandidate`
6. verify that no scientific-truth status is produced

---

# 27. ResearchCandidate

## 27.1 Authority

Only a validated `Session` may seal and mint a `ResearchCandidate`.

Sealing does not grant corpus authority.

The hypothesis package may create provisional candidate concepts marked `HYPOTHESIS`; only the session may seal the final research artifact.

There is no public promotion operation.

## 27.2 Required fields

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

## 27.3 Validation

Expose:

```text
Validate()
```

which performs:

```text
schema validation
→ canonicalization
→ ledger hash validation
→ derivation replay
→ provenance/candidate validation
```

Validation establishes artifact integrity, not physical truth.

# 28. Review artifact

The MVP defines review structures but does not implement reviewer-AI orchestration.

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

`Challenge` contains:

```text
StepID
Category
Severity
Description
```

`Review` contains:

```text
DerivationID
Challenges
ReviewerNotes
```

The reviewer consumes the sealed artifact externally.

---

# 29. Paper-translation documentation

The following are **auxiliary AI guidance**, not canonical manifest metadata:

```text
docs/paper-translation/common-notation.md
docs/paper-translation/framework-mapping.md
docs/paper-translation/ambiguity-resolution.md
```

They explain:

- common notation
- physics-object mapping
- framework identification
- ambiguity resolution
- distinction between notation and physical identity

The manifests remain the canonical machine-readable semantic source.

---

# 30. MRC fallibility and revision

MRC versioning is required.

Each derivation records:

```text
MRCVersion = "mrc-v0.4"
```

Each MRC rule has a stable ID.

The MVP MUST NOT provide an MRC bypass.

The conceptual future workflow is:

```text
current rule
    ↓
human-validated exception proposal
    ↓
evidence + falsification/review
    ↓
new scoped/versioned rule
    ↓
regression suite
```

MVP implements only the version/provenance foundation.

MRC is a model of current methodological constraints, not a claim of metaphysical infallibility.

---

# 31. Mandatory error taxonomy

Implement typed errors:

```text
DimensionMismatchError
CategoryMismatchError
AssumptionConflictError
ConventionConflictError
IdentifyError
ProvenanceError
CandidateContainmentError
InvalidObjectError
UnsupportedOperationError
ManifestValidationError
LedgerValidationError
```

All satisfy Go's standard `error` interface.

They MUST be inspectable with `errors.As`.

Normal MRC violations MUST NOT panic.

---

# 32. Determinism

The same input derivation executed twice MUST produce identical:

```text
canonical expression
expression hash
object canonical representation
assumption hash
convention hash
step hashes
derivation hash
manifest relation hash
ResearchCandidate canonical JSON
```

No:

```text
timestamps
random UUIDs
pointer addresses
map iteration order
environment-dependent ordering
```

may affect canonical artifacts.

---

# 33. Mandatory negative tests

The MVP MUST contain tests that fail for:

1. dimension mismatch
2. physical-kind mismatch with equal dimensions
3. assumption conflict
4. convention conflict
5. invalid zero `core.Object`
6. attempted trusted-object fabrication
7. `Simplify` creating physical identity
8. missing `Identify` justification
9. candidate contamination bypass
10. nonexistent promotion path
11. corrupted ledger
12. modified intermediate result
13. inconsistent manifest
14. exact-rational round-trip failure
15. nondeterministic canonicalization

Construction-authority negatives MUST be executed from external test packages where relevant.

---

# 34. Coverage matrix requirement

The implementation repository MUST contain, in the implementation plan, a matrix mapping:

```text
spec requirement
→ implementation file
→ test
```

Every normative MUST in this document MUST map to a concrete test or to an explicit, testable API invariant.

Every MRC rule ID:

```text
MRC-001 ... MRC-008
```

MUST map to at least one test.

A specification item without a test mapping is not considered complete.

---

# 35. Required MVP acceptance tests

The MVP is complete only when all of these pass:

### A — Typed mechanics

Construct valid mechanics objects through domain constructors.

### B — Classical relation

Reproduce `F = ma`.

### C — Dimension rejection

Reject incompatible dimensions.

### D — Category rejection

Reject distinct physical kinds even when dimensions match.

### E — Assumption conflict

Reject contradictory assumptions.

### F — Convention conflict

Reject contradictory conventions.

### G — Differentiation

Execute the specified symbolic derivative path.

### H — Relativistic relation

Construct the energy-momentum relation.

### I — Mass-energy derivation

Derive `E = mc²` through `Substitute → Simplify → Solve → SelectBranch`.

### J — Limit

Recover `LorentzFactor(v→0)=1`.

### K — Identify firewall

Verify `Simplify` cannot create `IDENTIFIED`.

### L — Candidate containment

Verify candidate provenance is contagious and cannot be promoted.

### M — Falsifiability

Verify a research candidate contains a prediction and falsification condition.

### N — Anomaly

Verify a candidate can reference a framework anomaly.

### O — Manifest

Verify manifests load and match executable constructors.

### P — Rational round-trip

Verify exact rational canonical JSON round-trip.

### Q — Ledger tamper

Verify hash-chain/replay validation detects altered steps.

### R — Determinism

Verify repeated derivations produce identical canonical artifacts.

### S — Candidate handoff

Verify a session alone can seal a `ResearchCandidate`.

---

# 36. Static `physvet` boundary

No `physvet` implementation is part of MVP.

The future analyzer MUST consume the same core semantics and contracts.

It MUST NOT independently redefine:

- dimensions
- physical kinds
- assumptions
- conventions
- provenance
- operation contracts

This is a v0.5+ architecture constraint.

---

# 37. v0.5+ deferrals

The following remain explicitly deferred:

1. `physvet`
2. Einstein's full 1905 series-expansion derivation
3. general series expansion/truncation engine
4. broader symbolic analysis
5. symbolic integration
6. richer tensor/index machinery
7. general relativity
8. quantum mechanics
9. QFT
10. statistical mechanics
11. electromagnetism
12. empirical-data adapters
13. MRC human exception workflow
14. reviewer-AI orchestration
15. automatic corpus governance
16. automatic hypothesis promotion
17. EBP 2.1 integration
18. numerical execution
19. general theory-management platform

---

# 37.1 MVP size guardrail

Target no more than **40 repository files total**, excluding `go.mod`, `go.sum`, generated artifacts, and VCS metadata.

Any plan or implementation exceeding the cap MUST justify every excess file against a normative MVP requirement.

Convenience-only files, deferred-domain scaffolding, compatibility shims, and placeholder abstractions are prohibited.

# 38. Final architectural boundary


The MVP is successful when this loop works deterministically:

```text
AI reads corpus
      ↓
typed physical construction
      ↓
explicit assumptions/conventions
      ↓
symbolic operations
      ↓
MRC enforcement
      ↓
explicit physical identification
      ↓
canonical derivation ledger
      ↓
candidate hypothesis
      ↓
prediction + falsification condition
      ↓
anomaly context
      ↓
sealed ResearchCandidate
      ↓
human review
```

The system stops there.

It does not tell the human what nature believes.

---

# 38.1 Coverage requirement

The implementation plan MUST map every normative MUST in this specification and every MRC rule ID to a concrete implementation location and at least one test.

No normative requirement is complete without a verification path.

# 39. Design principle


The governing principle for the MVP is:

> **Open representation, closed authority; explicit assumptions, explicit insight; deterministic mechanics, fallible MRC; machine-readable physics, human scientific judgment.**




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

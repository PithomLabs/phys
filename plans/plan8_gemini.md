# Implementation Plan: Physics Compiler MVP (v2.1)

## 1. Executive architecture

The module path is strictly pinned to `[github.com/PithomLabs/phys](https://github.com/PithomLabs/phys)`. The architecture follows a strict, unidirectional layered dependency graph that prevents import cycles and isolates trusted authority boundaries:

```text
internal/kernel  ← actual immutable object/Expr authority & validation
       ↑
core             ← public facade & type aliases (Object, Expr, Kind, Dimension, etc.)
       ↑
ops              ← pure, stateless symbolic transformations
       ↑
session          ← derivation ledger, step recording, replay, & candidate sealing
       ↑
mechanics / relativity / hypothesis → core + internal/kernel

```

```text
[Dependency Enforcement Rules]
1. core MUST NOT import ops or session.
2. ops MUST NOT import session.
3. session MAY import core, ops, and internal/kernel.
4. session MUST NOT import mechanics, relativity, or hypothesis.
5. mechanics, relativity, and hypothesis import core and internal/kernel; they interact with session solely by passing core.Object handles.

```

### Authority Boundaries

* **`internal/kernel` (Mint Authority):** Owns the concrete, unexported fields of `kernel.Object` and `kernel.Expr`. The unexported function `kernel.MintObject(ObjectSpec)` is the sole production entry point for instantiating valid physical objects. Go's `internal/` package rule prevents external modules from importing `internal/kernel`, enforcing constructor authority across public API boundaries.
* **`core` (Public Facade):** Exposes read-only type aliases (`core.Object = kernel.Object`) and value constructors for metadata (`Dimension`, `AssumptionSet`, `ConventionSet`, `Provenance`). It contains zero generic object factories or mutators.
* **`ops` (Stateless Symbolic Engine):** Pure, deterministic function library (`Add`, `Substitute`, `Differentiate`, etc.). `ops` retains zero ambient or global session state and cannot read or write to any ledger.
* **`session` (Derivation & Replay Authority):** Holds derivation draft state, ledger commit logic, hash-chain calculation, step replay validation, and `ResearchCandidate` sealing.
* **Domain Packages (`mechanics`, `relativity`, `hypothesis`):** Thin nominal wrappers exposing `CoreObject() core.Object` and fixed, domain-specific physical constructors.

---

## 2. Exact repository tree

The implementation comprises exactly **39 files** (including `go.mod`), adhering strictly to the <40 file guardrail:

```text
go.mod
README.md

internal/
    kernel/
        types.go
        mint.go

core/
    object.go
    expr.go
    dimension.go
    assumption.go
    convention.go
    provenance.go
    corpus.go
    canonical.go
    errors.go
    object_test.go
    expr_test.go

ops/
    arithmetic.go
    simplify.go
    transform.go
    relation.go
    dispatch.go
    operations_test.go
    negative_test.go

session/
    session.go
    ledger.go
    research_candidate.go
    session_test.go

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

hypothesis/
    candidate.go
    candidate_test.go

docs/
    paper-translation.md

```

---

## 3. Dependency-ordered implementation sequence

```text
Step 1: Kernel Foundation (internal/kernel)
├── Purpose: Establish unexported Object/Expr backing structs and invariant validation.
├── Dependencies: Go standard library (math/big, crypto/sha256).
├── Files: internal/kernel/types.go, internal/kernel/mint.go
├── Result: kernel.Object, kernel.Expr, ObjectSpec, and MintObject.
└── Gate: internal/kernel package compiles; unit tests verify unexported field safety.

Step 2: Core Model & Metadata Facade (core)
├── Purpose: Expose public type aliases, metadata types, canonical JSON, and typed errors.
├── Dependencies: internal/kernel
├── Files: core/errors.go, core/dimension.go, core/assumption.go, core/convention.go, 
│          core/provenance.go, core/corpus.go, core/canonical.go, core/expr.go, 
│          core/object.go, core/object_test.go, core/expr_test.go
├── Result: Immutable core.Object alias, dimension calculus, structural equality, SHA-256 hashing.
└── Gate: `go test ./core/...` passes; exact-rational and expression canonical JSON round-trips pass.

Step 3: Pure Symbolic Engine (ops)
├── Purpose: Implement stateless symbolic operations, simplification, calculus, solver, and dispatch.
├── Dependencies: core, internal/kernel
├── Files: ops/arithmetic.go, ops/simplify.go, ops/transform.go, ops/relation.go, 
│          ops/dispatch.go, ops/operations_test.go, ops/negative_test.go
├── Result: Pure transformations (Add, Substitute, Differentiate, Limit, Solve, SelectBranch, Apply).
└── Gate: `go test ./ops/...` passes; all negative MRC error assertions (MRC-001–MRC-005, MRC-008) pass.

Step 4: Physics Corpus Packages (mechanics, relativity)
├── Purpose: Provide domain-specific primitive wrappers, fixed relations, and embedded semantic manifests.
├── Dependencies: core, internal/kernel
├── Files: mechanics/primitives.go, mechanics/relations.go, mechanics/manifest.json, mechanics/manifest_test.go, mechanics/relations_test.go,
│          relativity/primitives.go, relativity/relations.go, relativity/manifest.json, relativity/manifest_test.go, relativity/derivation_test.go
├── Result: Mechanics and relativity domain primitives, relations, and validated embed manifests.
└── Gate: Manifest cross-check tests verify constructor expression/dimension/kind agreement.

Step 5: Derivation Ledger & Replay Engine (session)
├── Purpose: State-machine derivation orchestration, step logging, SHA-256 hash chaining, replay validation, and Identify ownership.
├── Dependencies: core, ops, internal/kernel
├── Files: session/ledger.go, session/session.go, session/research_candidate.go, session/session_test.go
├── Result: Session lifecycle, step ledger, replay validator, and ResearchCandidate sealing.
└── Gate: `go test ./session/...` passes; tamper detection, replay verification, and hash-chain checks pass.

Step 6: Provisional Hypothesis & Falsifiability (hypothesis)
├── Purpose: AI-callable candidate concept construction with enforced HYPOTHESIS provenance and containment.
├── Dependencies: core, internal/kernel
├── Files: hypothesis/candidate.go, hypothesis/candidate_test.go
├── Result: NewCandidateConcept, prediction/falsification structures, containment enforcement.
└── Gate: `go test ./hypothesis/...` passes; export surface inspection verifies no promotion API exists.

Step 7: Canonical Vertical Slices & Documentation (docs)
├── Purpose: End-to-end verification of E = mc² derivation, calculus micro-test, limit micro-test, and paper translation artifact.
├── Dependencies: All repository packages.
├── Files: docs/paper-translation.md, relativity/derivation_test.go, mechanics/relations_test.go
├── Result: End-to-end paper translation verification and acceptance gate execution.
└── Gate: `go test ./...` and `go vet ./...` pass clean across all packages.

```

---

## 4. Core model implementation

### Authority Isolation

`internal/kernel/types.go` defines the unexported `kernel.Object` struct:

```go
type Object struct {
    valid        bool
    name         string
    kind         Kind
    dimension    Dimension
    expr         Expr
    assumptions  AssumptionSet
    conventions  ConventionSet
    provenance   Provenance
    corpusStatus CorpusStatus
}

```

`core/object.go` exposes `type Object = kernel.Object`. Because all fields are unexported and `core` exposes no `NewObject` factory, external callers cannot instantiate `core.Object{}` with `valid: true`. The zero value is `valid: false`.

`internal/kernel/mint.go` provides the sole production minting entry point:

```go
type ObjectSpec struct {
    Name          string
    Kind          Kind
    Dimension     Dimension
    Expr          Expr
    Assumptions   AssumptionSet
    Conventions   ConventionSet
    Provenance    Provenance
    CorpusStatus  CorpusStatus
}

func MintObject(spec ObjectSpec) (Object, error)

```

`MintObject` performs invariant checking: validates non-empty names where required, verifies expression handles, checks dimension validity, normalizes metadata sets, and returns an immutable `kernel.Object`.

### Data Types & Accessors

* **`core.Object` Read-Only Accessors:** `Valid() bool`, `Name() string`, `Kind() Kind`, `Dimension() Dimension`, `Expr() Expr`, `Assumptions() AssumptionSet`, `Conventions() ConventionSet`, `Provenance() Provenance`, `CorpusStatus() CorpusStatus`.
* **Kind (18 Stable Ordinals):** `Mass`, `Time`, `Position`, `Velocity`, `Acceleration`, `Force`, `Momentum`, `Energy`, `KineticEnergy`, `RestMass`, `ThreeMomentum`, `FourMomentum`, `SpeedOfLight`, `Spacetime`, `MinkowskiMetric`, `Expression`, `Relation`, `BranchSet`.
* **Dimension (7 SI Base Exponents):** Represented internally by `math/big.Rat` exponents for $[M, L, T, I, \Theta, N, J]$. Public fixed constructors: `Dimensionless()`, `DimensionMass()`, `DimensionLength()`, `DimensionTime()`, `DimensionVelocity()`, `DimensionAcceleration()`, `DimensionForce()`, `DimensionMomentum()`, `DimensionEnergy()`. Arithmetic methods: `Equal`, `Multiply`, `Divide`, `Pow`, `CanonicalJSON`, `Hash`.
* **AssumptionSet & Assumption:** Closed union value (`TextValue(string)` or `ExprValue(Expr)`). Kinds: `Domain`, `Regime`, `Constraint`, `Convention`, `Approximation`, `MathPrecondition`, `PhysicalAssumption`. `Merge` performs deterministic set union with conflict checking (returning `AssumptionConflictError`). Bounded entailment methods: `EntailsNonNegative(expr)`, `EntailsPositive(expr)`, `EntailsNonZero(expr)`.
* **ConventionSet & Convention:** Key/value string metadata. `Merge` detects conflicting values for identical keys (returning `ConventionConflictError`).
* **Provenance:** Status enum (`DEFINED`, `POSTULATED`, `DERIVED`, `IDENTIFIED`, `APPROXIMATED`, `HYPOTHESIS`). Read-only metadata: `Status`, `Source`, `Framework`, `ParentHashes`, `AssumptionHash`, `ConventionHash`, `MRCVersion` (`mrc-v0.4`), `Justification`. Constructed via `NewProvenance(...)`.
* **CorpusStatus:** Enum (`NONE`, `ESTABLISHED`, `CONTESTED`, `SUPERSEDED`, `FALSIFIED`). Human-curated repository metadata; never programmatically inferred or computed.

---

## 5. MRC implementation map

| Rule | Exact Semantic Check | Enforcement Location | Failure | Test |
| --- | --- | --- | --- | --- |
| **MRC-001** | Object constructed strictly via approved domain/session mint paths; zero-value `!Valid()` rejected. | `internal/kernel/mint.go`, `core/object.go` | `InvalidObjectError` / `ProvenanceError` | `core/object_test.go`, `ops/negative_test.go` |
| **MRC-002** | `Add`/`Subtract`/`Compare` require exact `Dimension.Equal()`; `Multiply`/`Divide`/`Pow` compute derived dimensions; `Substitute`/`Differentiate` validate dimensional calculus. | `ops/arithmetic.go`, `ops/transform.go`, `ops/relation.go` | `DimensionMismatchError` | `ops/negative_test.go` |
| **MRC-003** | `Add`/`Subtract` require same named Kind or `Expression`+`Expression`; `Compare` allows same Kind or `Expression` match; Inequalities require scalar ordered kinds; `Multiply`/`Divide` yield `Expression`. | `ops/arithmetic.go`, `ops/relation.go` | `CategoryMismatchError` | `ops/negative_test.go` |
| **MRC-004** | Merging assumptions across operations enforces set union; duplicate keys with non-identical values trigger conflict error. | `core/assumption.go`, `ops/*` | `AssumptionConflictError` | `ops/negative_test.go` |
| **MRC-005** | Merging conventions across operations enforces set union; duplicate keys with non-identical values trigger conflict error. | `core/convention.go`, `ops/*` | `ConventionConflictError` | `ops/negative_test.go` |
| **MRC-006** | `IDENTIFIED` status created strictly via `Session.Identify` with non-empty justification, valid operands, and matching dimensions/kinds. | `session/session.go` | `IdentifyError` | `ops/negative_test.go`, `session/session_test.go` |
| **MRC-007** | Committed derivation ledgers and sealed `ResearchCandidate` artifacts created strictly via `session.Session`. | `session/session.go`, `session/ledger.go` | `ProvenanceError` | `session/session_test.go` |
| **MRC-008** | If any operand has `HYPOTHESIS` provenance, downstream results MUST inherit `HYPOTHESIS`. No API exists to promote hypotheses to trusted status. | `ops/*`, `session/session.go`, `hypothesis/candidate.go` | `CandidateContainmentError` | `hypothesis/candidate_test.go` |

---

## 6. Symbolic engine

### Closed Expression Node Set (10 Nodes)

Node ordinals: `0: Symbol`, `1: Rational`, `2: Add`, `3: Mul`, `4: Neg`, `5: Pow`, `6: Sqrt`, `7: Call`, `8: Relation`, `9: BranchSet`.
Relation operators: `RelationEq` ("eq"), `RelationNeq` ("neq"), `RelationLt` ("lt"), `RelationLte` ("lte"), `RelationGt` ("gt"), `RelationGte` ("gte").

```go
// Public Expr inspection API
func NewSymbol(name string) (Expr, error)
func NewRational(value *big.Rat) Expr
func NewAdd(terms ...Expr) Expr
func NewMul(factors ...Expr) Expr
func NewNeg(expr Expr) Expr
func NewPow(base Expr, exponent *big.Rat) Expr
func NewSqrt(expr Expr) Expr
func NewCall(functionID string, args ...Expr) (Expr, error)
func NewRelation(op RelationOperator, lhs, rhs Expr) Expr
func NewBranchSet(target Expr, branches ...Expr) Expr

```

### Simplification & Canonical Normalization Rules

1. **Rational Representation:** `math/big.Rat` reduced fraction (`num/den`), positive denominator, zero normalized to `0/1`.
2. **Child Sorting (3-Key Order):** Sorted by (i) `ExprKind` ordinal, (ii) lowercase hex SHA-256 child hash byte comparison, (iii) canonical JSON byte sequence comparison using `bytes.Compare`.
3. **N-ary Flattening & Identities:** Flatten nested `Add`/`Mul`; remove identities ($x + 0 \to x$, $x \cdot 1 \to x$, $x^1 \to x$, $x^0 \to 1$ when $x \neq 0$ safe); combine rational constants.
4. **Sign Normalization:** $\text{Mul}(-1, x) \to \text{Neg}(x)$, $\text{Neg}(\text{Rational}(q)) \to \text{Rational}(-q)$, $\text{Neg}(\text{Neg}(x)) \to x$.
5. **Repeated Powers & Sqrt:** Combine $x \cdot x \to x^2$; $(x^a)^b \to x^{a \cdot b}$ for non-negative integer rationals. Apply $\sqrt{x^2} \to x$ **only** when `AssumptionSet.EntailsNonNegative(x)` returns true.
6. **Bounded Sign Entailment:** Structural checking without a theorem prover:
* $\text{Rational}(q)$ with $q \ge 0$ is non-negative.
* $\text{Symbol}(x)$ with explicit $x \ge 0$ assumption is non-negative.
* $\text{Mul}(f_1, \dots, f_n)$ where all factors are non-negative is non-negative.
* $\text{Pow}(x, n)$ with even integer $n \ge 0$ is non-negative.



### Bounded Calculus, Solver, & Dispatch

* **`Differentiate(target, wrt)`:** Bounded recursive differentiation for `Rational` ($\to 0$), `Symbol` ($\to 1$ if match, else $0$), `Add` (sum of derivatives), `Mul` (n-ary product rule), `Neg` ($\text{Neg}(d/dx)$), and `Pow` with non-negative integer exponent ($n \cdot b^{n-1} \cdot db/dx$). Unsupported nodes (`Sqrt`, `Call`, `Relation`, `BranchSet`) return `UnsupportedOperationError`.
* **`Limit(target, variable, value)`:** Expands fixed function bodies (`lorentz_factor` $\to 1 / \sqrt{1 - (v/c)^2}$), performs direct symbol substitution, and applies simplification.
* **`Solve(relation, target)`:** Pattern matches exact quadratic form $\text{Relation}(\text{eq}, \text{Pow}(\text{Symbol}(\text{target}), 2), \text{Expr})$. Returns $\text{BranchSet}(\text{target}, [\sqrt{\text{rhs}}, -\sqrt{\text{rhs}}])$. All other forms return `UnsupportedOperationError`.
* **`SelectBranch(branches, constraint)`:** Accepts $\text{BranchSet}(\text{target}, [\text{pos}, \text{Neg}(\text{pos})])$ and constraint $\text{Relation}(\text{gte}, \text{target}, 0)$. Filters non-negative branch and returns simplified `Expression`.
* **`Apply(id, inputs, params)`:** Closed dispatch function over operation IDs: `add`, `subtract`, `multiply`, `divide`, `pow`, `simplify`, `substitute`, `differentiate`, `limit`, `compare`, `solve`, `select_branch`. (`identify` is excluded from `ops.Apply`).

---

## 7. Physics corpus

### Semantic Manifest Architecture

`mechanics/manifest.json` and `relativity/manifest.json` are embedded into their respective packages via `go:embed`. Manifest loading deserializes into strict, typed structs (`Manifest`, `ManifestItem`, `ManifestDomain`, `ManifestLimit`, `ManifestAnomaly`, `ManifestReduction`). Dynamic `map[string]any` is prohibited. The `statement` field is documentation text and is never parsed as executable mathematics.

```go
// Pure bytes validator in core/corpus.go
func ValidateManifestBytes(data []byte) (Manifest, error)

```

### Constructor Cross-Check Mechanism

Package unit tests (`manifest_test.go`) maintain a static function map resolving constructor identifiers:

```go
var constructorMap = map[string]func() core.Object{
    "mechanics.NewtonSecondLaw": mechanics.NewtonSecondLaw,
    "relativity.EnergyMomentumRelation": relativity.EnergyMomentumRelation,
    // ...
}

```

The cross-check test iterates over every item in `manifest.json`:

1. Resolves `item.Constructor` against `constructorMap`.
2. Invokes constructor function to obtain `core.Object`.
3. Verifies `object.Expr()` canonical JSON and hash match `item.CanonicalExpr`.
4. Verifies `object.Dimension()` equals `item.Dimension`.
5. Verifies `object.Kind()` equals `item.Kind`.
6. Verifies `object.Provenance().Status` equals `item.ProvenanceStatus`.

---

## 8. Assumptions, conventions, provenance, and candidate containment

### Merge and Conflict Rules

* **Assumptions:** Union of input sets. Key collision with differing canonical values returns `AssumptionConflictError`. Denominator preconditions auto-generate keys formatted as `"denominator/" + HashExpr(b)`.
* **Conventions:** Union of input sets. Key collision with differing canonical values returns `ConventionConflictError`.

### Deterministic Provenance Propagation Law

$$\text{Status}_{\text{out}} = \begin{cases} \text{HYPOTHESIS} & \text{if } \exists \, i : \text{Status}_{\text{in}, i} = \text{HYPOTHESIS} \\ \text{IDENTIFIED} & \text{if action is } \texttt{Session.Identify} \\ \text{DERIVED} & \text{for pure } \texttt{ops} \text{ transformations} \end{cases}$$

For assertion steps in `Session`:

* `Postulate`: Requires input `POSTULATED`.
* `Define`: Requires input `DEFINED`.
* `Declare`: Preserves input status.

### Hypothesis Contamination Law

`HYPOTHESIS` provenance is transitive and sticky. If an input carries `HYPOTHESIS`, all downstream operations produce `HYPOTHESIS`. Any attempt to validate a hypothesis-dependent artifact under a trusted corpus status fails with `CandidateContainmentError`.

---

## 9. Derivation ledger

### Lifecycle State Machine

$$\texttt{New} \xrightarrow{\text{Draft()}} \texttt{Drafting} \xrightarrow{\text{Commit()}} \texttt{Committed} \xrightarrow{\text{Conclude()}} \texttt{Concluded} \xrightarrow{\text{Seal()}} \texttt{Sealed}$$

Mutations (`Postulate`, `Declare`, `Define`, `Step`, `Identify`) are permitted only in `Drafting`. After `Seal()`, all mutation attempts return `ProvenanceError`.

### Step Ledger & Hash Chain Structure

Each committed step generates a deterministic `StepID` (`step-000001`, `step-000002`, ...). Genesis `PreviousStepHash` is 64 lowercase zeros (`0000000000000000000000000000000000000000000000000000000000000000`).

```go
type StepEnvelope struct {
    PreviousHash string   `json:"previous_hash"`
    Step         StepBody `json:"step"`
}

// CurrentStepHash = SHA256(CanonicalJSON(StepEnvelope))

```

Every step retains `InputCanonicals[]`, `ParamsCanonical`, and `OutputCanonical` JSON bytes alongside `InputHashes[]` and `OutputHash`.

### Replay & Anti-Tampering Engine

`Session.Validate()` and `Ledger.Validate()` execute step-by-step verification:

1. Re-verify genesis hash and `StepID`/index sequence.
2. Decode `InputCanonicals[]` and verify `HashObject(decoded) == InputHash`.
3. Decode `OutputCanonical` and verify `HashObject(decoded) == OutputHash`.
4. Re-calculate hash chain across all `StepEnvelope` structures.
5. Re-execute transformation steps via `ops.Apply(OperationID, inputs, params)` and identification steps via the session path.
6. Verify replayed output hash matches stored `OutputHash`. Divergence returns `LedgerValidationError`.

### Import Cycle Prevention

`session` imports `core` and `ops`. `core` and `ops` do not import `session`. `Session.Identify` is owned exclusively by `session/session.go`; no package-level `ops.Identify` exists.

---

## 10. Hypothesis and ResearchCandidate

### Candidate Concept Construction

`hypothesis.NewCandidateConcept` builds provisional objects without exposing a general factory:

```go
func NewCandidateConcept(
    id string,
    kind core.Kind,
    dimension core.Dimension,
    expr core.Expr,
    assumptions core.AssumptionSet,
    conventions core.ConventionSet,
) (core.Object, error)

```

`NewCandidateConcept` forces `Provenance.Status = HYPOTHESIS` and `CorpusStatus = NONE`.

### Sealed ResearchCandidate Artifact

Minted exclusively via `Session.Seal()`. Contains immutable getters:
`ID()`, `Hypothesis()`, `Premises()`, `Assumptions()`, `Derivation()`, `FrameworkDependencies()`, `Predictions()`, `FalsificationConditions()`, `RecoveryClaims()`, `AnomalyReferences()`, `ReviewHistory()`, `MRCVersion()`, `LedgerHash()`.

External loader: `ParseResearchCandidateJSON(data []byte) (ResearchCandidate, error)`. `Validate()` verifies hash integrity, ledger replay, and candidate containment without adjudicating physical truth.

---

## 11. Canonical vertical slices

### Vertical Slice Test Matrix

```text
Slice 1: Typed Mechanics & F = ma
├── Operations: NewForce(), NewMass(), NewAcceleration(), NewtonSecondLaw()
└── Assertions: F = m*a canonical expr match; Kind/Dimension equality; Manifest agreement.

Slice 2: Negative Rejections
├── Operations: Add(Mass, Time); Add(Force, Torque); Merge conflicting assumptions/conventions.
└── Assertions: Mismatched dimensions -> DimensionMismatchError; Mismatched kinds -> CategoryMismatchError; Conflicting keys -> AssumptionConflictError / ConventionConflictError.

Slice 3: Kinetic Energy Calculus
├── Operations: Construct K = 1/2*m*v^2; Differentiate(K, v); Compare(result, Momentum, eq).
└── Assertions: Derivative yields canonical m*v; Compare with named Momentum succeeds under MRC-003.

Slice 4: E = mc² Derivation
├── Sequence: 1. EnergyMomentumRelation() [E^2 = (p*c)^2 + (m*c^2)^2]
│             2. Substitute(relation, p, ZeroThreeMomentum()) -> E^2 = (0*c)^2 + (m*c^2)^2
│             3. Simplify() -> E^2 = (m*c^2)^2
│             4. Solve(eq, E) -> BranchSet(E, [Sqrt((m*c^2)^2), Neg(Sqrt((m*c^2)^2))])
│             5. Compare(E, ZeroEnergy(), gte) -> E >= 0
│             6. SelectBranch(branches, E >= 0) -> Sqrt((m*c^2)^2)
│             7. Simplify() using m >= 0 assumption -> m*c^2
└── Assertions: Final expression is exactly m*c^2; derivation is purely mechanical; MassEnergyRelation() constructor is NOT called.

Slice 5: Relativistic Limit
├── Operations: Limit(LorentzFactor(), Velocity(), ZeroVelocity())
└── Assertions: Expands fixed body 1/Sqrt(1-(v/c)^2); substitutes v=0; simplifies denominator -> canonical rational 1.

Slice 6: Identify Firewall
├── Operations: Simplify(x); Session.Identify(a, b, "justification").
└── Assertions: Simplify output is DERIVED; Identify output is IDENTIFIED + recorded step in ledger.

Slice 7: Candidate Contamination & Containment
├── Operations: NewCandidateConcept(); derive downstream result; attempt candidate sealing.
└── Assertions: Downstream status is HYPOTHESIS; fails validation if claimed as ESTABLISHED; export inspection confirms zero promotion functions.

Slice 8: Manifest Validation & Canonical JSON Round-Trip
├── Operations: ParseManifest(relativity); Expr/Object -> CanonicalJSON -> Parse -> Hash check.
└── Assertions: Manifest items match executable constructors; JSON round-trip yields byte-identical bytes and hashes.

Slice 9: Ledger Tamper & Replay Detection
├── Operations: Build derivation; mutate OutputCanonical in step; execute Session.Validate().
└── Assertions: Hash mismatch / replay divergence triggers LedgerValidationError.

```

---

## 12. Coverage matrix

| Requirement / Rule ID | Implementation Location | Primary Test | Secondary Verification |
| --- | --- | --- | --- |
| **REQ-000-01** | `go.mod` | `go.mod` inspection | `go test ./...` |
| **REQ-000-02** | `go.mod` | `go.mod` dependency check | Build verification |
| **REQ-000-03** | `docs/paper-translation.md` | Architecture verification | Documentation audit |
| **REQ-000-04** | `docs/paper-translation.md` | Threat model audit | Code comment review |
| **REQ-001-01** | `internal/kernel/mint.go`, `ops/*` | `ops/operations_test.go` | Vertical Slice 1–4 |
| **REQ-001-02** | `mechanics/*`, `relativity/*` | `mechanics/relations_test.go` | `relativity/derivation_test.go` |
| **REQ-001-03** | `mechanics/manifest.json`, `relativity/manifest.json` | `mechanics/manifest_test.go` | `relativity/manifest_test.go` |
| **REQ-001-04** | `core/dimension.go`, `core/assumption.go` | `core/object_test.go` | Vertical Slice 2 |
| **REQ-001-05** | Entire repository | Export surface inspection | No truth score types exist |
| **REQ-001-06** | `hypothesis/candidate.go`, `session/research_candidate.go` | `hypothesis/candidate_test.go` | `session/session_test.go` |
| **REQ-002-01..26** | Non-goal compliance | File tree & export surface inspection | Zero prohibited features |
| **REQ-003-01..03** | Repository layout | File tree count check (39 files) | External package tests |
| **REQ-004-01..03** | `go.mod`, package imports | `go vet ./...` / import graph test | Dependency cycle check |
| **REQ-005-01..12** | `internal/kernel/types.go`, `core/object.go` | `core/object_test.go` | `ops/negative_test.go` |
| **REQ-006-01..03** | `core/object.go`, `ops/arithmetic.go` | `ops/negative_test.go` | Vertical Slice 2 |
| **REQ-007-01..03** | `core/dimension.go` | `core/object_test.go` | `ops/operations_test.go` |
| **REQ-008-01..06** | `core/expr.go`, `ops/simplify.go` | `core/expr_test.go` | `ops/operations_test.go` |
| **REQ-009-01..06** | `core/canonical.go`, `ops/simplify.go` | `core/expr_test.go` | `ops/operations_test.go` |
| **REQ-010-01..07** | `core/canonical.go` | `core/expr_test.go` | Vertical Slice 8 |
| **REQ-011-01..02** | `core/assumption.go` | `core/object_test.go` | `ops/negative_test.go` |
| **REQ-012-01..02** | `core/convention.go` | `core/object_test.go` | `ops/negative_test.go` |
| **REQ-013-01..05** | `core/provenance.go`, `core/corpus.go` | `core/object_test.go` | `hypothesis/candidate_test.go` |
| **REQ-015-01** | `ops/*` | `ops/operations_test.go` | Parallel test execution |
| **REQ-016-01** | `session/session.go`, `session/ledger.go` | `session/session_test.go` | Vertical Slice 9 |
| **REQ-024-01..02** | `hypothesis/candidate.go` | `hypothesis/candidate_test.go` | Export surface inspection |
| **REQ-032-01..22** | `ops/negative_test.go`, `session/session_test.go` | Mandatory negative test suite | Acceptance test suite |
| **MRC-001..008** | Core / Ops / Session enforcement | MRC Test Suite in Section 5 | Vertical Slices 1–9 |

---

## 13. Acceptance gate

The implementation is complete only when all of the following conditions pass:

1. **Compilation & Standard Tooling Clean:**
```bash
go test ./...
go vet ./...

```


2. **File Budget Guardrail:** Total repository file count MUST NOT exceed 39 files.
3. **Zero Third-Party Dependencies:** `go.mod` references only standard library packages under `go 1.24`.
4. **Acceptance Tests A–S All Pass:**
* **A (Typed Mechanics):** Construct mechanics primitives via fixed constructors.
* **B (Classical Relation):** Validate $F = m \cdot a$ manifest/constructor agreement.
* **C (Dimension Rejection):** Dimension mismatch returns `DimensionMismatchError`.
* **D (Category Rejection):** Kind mismatch with equal dimensions returns `CategoryMismatchError`.
* **E (Assumption Conflict):** Conflicting assumption values return `AssumptionConflictError`.
* **F (Convention Conflict):** Conflicting convention values return `ConventionConflictError`.
* **G (Differentiation):** $d/dv(1/2 m v^2) \to m \cdot v$; `Compare` with `Momentum` succeeds.
* **H (Relativistic Relation):** Energy-momentum relation constructor validation.
* **I (E = mc² Derivation):** Full vertical substitution, solve, branch selection derivation.
* **J (Relativistic Limit):** $\lim_{v \to 0} \gamma(v) = 1$ via fixed body expansion.
* **K (Identify Firewall):** `Simplify` yields `DERIVED`; `Identify` yields `IDENTIFIED`.
* **L (Candidate Contamination):** `HYPOTHESIS` input propagates transitively.
* **M (Falsifiability & Anomaly):** Sealed candidate retains prediction and anomaly reference.
* **N (Manifest Validation):** Manifest JSON deserialization and constructor cross-checking.
* **O (Exact-Rational Round-Trip):** Lossless JSON byte round-trip for exact rationals.
* **P (Ledger Tamper/Replay):** Mutated canonical outputs trigger `LedgerValidationError`.
* **Q (Determinism):** Dual derivation runs produce byte-identical canonical JSON and hashes.
* **R (Sealed Candidate Handoff):** `ResearchCandidate` validates internally without truth score.
* **S (No Promotion API):** AST inspection verifies absence of exported hypothesis promotion methods.



---

## 14. Explicit deferrals

The following items are explicitly deferred to version v0.5+ and MUST NOT be included in the MVP codebase:

* `physvet` CLI tool or linter package.
* General series expansion, Taylor series, or asymptotic truncation engines.
* Full 1905 Einstein mass-energy paper derivation.
* General computer algebra system (CAS), integration, or differential equations solver.
* General differential geometry, general relativity, tensor index calculus.
* Electromagnetism (`em`), Quantum Mechanics (`qm`), Quantum Field Theory (`qft`), Statistical Mechanics (`statmech`).
* Numerical execution, simulation, floating-point evaluation, GPU execution.
* Empirical data ingestion, data adapters, or experimental datasets.
* Multi-agent reviewer orchestration runtime or automated review generation.
* Runtime MRC exception approval workflow or rule overrides.
* Automated hypothesis promotion or corpus status inference.
* Evidence-Based Profiling (EBP 2.1) integration.

---

## 15. Scope accounting

* **Exact Total File Count:** 39 files
* **Source / Test / Manifest / Doc Split:**
* Go Source Files: 21
* Go Test Files: 12
* Manifest JSON Files: 2
* Config & Documentation Files: 4 (`go.mod`, `README.md`, `docs/paper-translation.md`, and this plan document)


* **Major Exported Types:** `Object`, `Expr`, `Kind`, `Dimension`, `Assumption`, `AssumptionSet`, `Convention`, `ConventionSet`, `Provenance`, `CorpusStatus`, `Manifest`, `Session`, `Ledger`, `Step`, `ResearchCandidate`, `Prediction`, `FalsificationCondition`, `RecoveryClaim`, `AnomalyReference`, `Review`, `Challenge`.
* **Major Operations:** `Add`, `Subtract`, `Multiply`, `Divide`, `Pow`, `Simplify`, `Substitute`, `Differentiate`, `Limit`, `Compare`, `Solve`, `SelectBranch`, `Apply`, `Identify`.
* **Expected External Dependencies:** 0 (Standard library `go 1.24` only).

Open Spec Items: NONE
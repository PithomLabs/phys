# Physics Compiler MVP — Normative Implementation Specification v2.2

## 0. Authority, module identity, and governing boundary

This document is the **normative implementation specification** for the Physics Compiler MVP.

The implementation-planning prompt is subordinate to this document. Where wording differs, this document wins.

### 0.1 Exact Go module identity

The Go module MUST be:

```text
github.com/PithomLabs/phys
```

This is a hard invariant.

All repository-local imports MUST use the `github.com/PithomLabs/phys/...` path.

The MVP MUST use the Go standard library only.

The baseline Go language version is:

```text
go 1.24
```

A later installed Go toolchain MAY build the module, but `go.mod` MUST retain the module path above and MUST NOT introduce third-party dependencies for MVP.

**REQ-000-01**: `go.mod` MUST declare `module github.com/PithomLabs/phys`.

**REQ-000-02**: MVP MUST NOT require third-party Go modules.

### 0.2 Product boundary

The MVP is a deterministic symbolic physics substrate intended to be usable by an AI agent as machine-readable physics corpus plus formal pen-and-paper machinery.

The AI is the theorist/formulator.

The library is the formal symbolic substrate.

Human researchers and empirical reality remain the final scientific authority.

The library checks formal structure, dimensions, physical-category compatibility, assumptions, conventions, provenance containment, and artifact integrity. It does not adjudicate physical truth.

**REQ-000-03**: The implementation MUST preserve this human scientific authority boundary.

### 0.3 Threat model for constructor authority

The MVP constructor-authority guarantee is a **public API/package-boundary guarantee**, not a malicious-contributor security boundary.

It protects external callers, including AI-generated code, from directly manufacturing trusted physical objects through the public API.

Code already inside this repository can import its own `internal/kernel` package and therefore has minting authority. Repository contributors are trusted maintainers for this MVP.

MVP MUST NOT claim that an adversarial contributor with write access to the source tree cannot bypass the internal authority boundary.

**REQ-000-04**: Documentation MUST describe constructor authority as an external/public API boundary, not as cryptographic or hostile-source protection.

---

# 1. Product requirements

The implementation MUST demonstrate all of the following.

**REQ-001-01 — Fundamental primitives + MRC**

The library MUST represent physical equations using typed physical primitives and MUST enforce MRC through Go package boundaries, constructor authority, immutable metadata, and operation-time semantic checks.

**REQ-001-02 — Frameworks as evaluation targets**

`mechanics` and `relativity` MUST be represented as explicit frameworks with assumptions, conventions, limitations/anomalies, provenance, and executable derivation targets.

They MUST NOT be represented as the final ontology of nature.

**REQ-001-03 — Machine-readable corpus**

Each populated framework package MUST provide a validated machine-readable semantic manifest.

The manifests are the canonical machine-readable corpus metadata for the populated packages.

**REQ-001-04 — Explicit limitations**

The implementation MUST make dimensions, physical kinds, assumptions, conventions, provenance, corpus status, framework scope, anomalies/limitations, and candidate containment explicit.

**REQ-001-05 — No truth adjudication**

The library MUST NOT compute physical truth, truth scores, probabilities of truth, theory rankings, or automatic corpus-status judgments.

**REQ-001-06 — Hypothesis formation**

The AI-facing hypothesis package MUST allow provisional candidate concepts. The session layer MUST allow a sealed `ResearchCandidate` artifact containing a derivation, predictions, falsification conditions, recovery claims, anomaly references, and review history.

Promotion to trusted corpus authority remains outside MVP and human-controlled.

---

# 2. Explicit non-goals

The MVP MUST NOT implement any of the following.

**REQ-002-01** `.phys` or another physics-specific language.

**REQ-002-02** A lexer, parser, or frontend for physics source text.

**REQ-002-03** A custom compiler language.

**REQ-002-04** A general mathematical ontology.

**REQ-002-05** A full computer algebra system.

**REQ-002-06** A theorem prover or general symbolic proof assistant.

**REQ-002-07** Numerical execution.

**REQ-002-08** Numerical integration, optimization, simulation, Monte Carlo, GPU execution, or numerical ODE solving.

**REQ-002-09** Empirical-data ingestion or automated empirical validation.

**REQ-002-10** Automatic theory ranking, theory selection, physical-truth scoring, or probability-of-truth scoring.

**REQ-002-11** Automatic corpus-status inference.

**REQ-002-12** Automatic hypothesis promotion.

**REQ-002-13** An MRC runtime bypass or override API.

**REQ-002-14** An MRC exception-approval workflow.

**REQ-002-15** `physvet` implementation.

**REQ-002-16** Electromagnetism package.

**REQ-002-17** Quantum mechanics package.

**REQ-002-18** Quantum field theory package.

**REQ-002-19** Statistical mechanics package.

**REQ-002-20** General relativity.

**REQ-002-21** General differential geometry or general tensor/index algebra.

**REQ-002-22** General symbolic integration.

**REQ-002-23** General symbolic series expansion/truncation.

**REQ-002-24** General multi-agent orchestration runtime.

**REQ-002-25** EBP 2.1 integration.

**REQ-002-26** A general theory-management platform.

---

# 3. Exact MVP repository layout

The implementation MUST create or modify only the following files.

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

The baseline contains **39 repository files including `go.mod`** and **38 files when `go.mod` is excluded**, which is below the 40-file MVP guardrail.

No additional source, test, data, helper, compatibility, or future-scaffolding files are permitted.

**REQ-003-01**: The implementation MUST use only the files listed above unless a normative requirement in this document is impossible to satisfy without one additional file; such an exception would itself be a specification defect and is not expected in v2.2.

**REQ-003-02**: Tests MUST use adjacent Go `*_test.go` files.

**REQ-003-03**: Constructor-authority negative tests MUST use external package tests when package visibility is the subject of the assertion.

---

# 4. Package dependency graph and authority boundary

The exact dependency direction is:

```text
internal/kernel
      ↓
core
      ↓
ops
      ↓
session
      ↑
mechanics / relativity / hypothesis
```

More precisely:

```text
internal/kernel  ← concrete immutable Object/Expr/Kind/Dimension authority
       ↑
core             ← public aliases + package-level façade functions
       ↑
ops              ← pure symbolic transformations
       ↑
session         ← derivation/ledger/replay/seal authority

mechanics ----→ core + internal/kernel
relativity ---→ core + internal/kernel
hypothesis ---→ core + internal/kernel
```

The concrete backing definitions for `Object`, `Expr`, `Kind`, and `Dimension` live in `internal/kernel`.
The corresponding files in `core/` expose aliases and package-level façade constructors/helpers; they do not declare methods on aliased types.
`ops` and `session` may import `internal/kernel` directly for validated minting and replay decoding. This does not create a dependency cycle.

`core` MUST NOT import `ops`.

`core` MUST NOT import `session`.

`ops` MUST NOT import `session`.

`session` MAY import `core`, `ops`, and `internal/kernel`.

The internal package is used specifically to solve the Go visibility problem created by immutable `core.Object` fields plus multiple cooperating packages.

**REQ-004-01**: No dependency cycle may exist.

**REQ-004-02**: `core` MUST remain independent of `ops` and `session`.

**REQ-004-03**: Public external callers MUST NOT be able to import `internal/kernel` because of Go's `internal/` package rule.

---

# 5. Core physical-object abstraction

## 5.1 Two-layer representation

The public model has two layers:

1. nominal domain types such as `Mass`, `Velocity`, `Energy`, `FourMomentum`, and `KineticEnergy`;
2. one immutable generic `core.Object` carrier used by `ops`.

A domain wrapper contains exactly one `core.Object` value.

Required pattern:

```go
type Mass struct {
    object core.Object
}

func (m Mass) CoreObject() core.Object {
    return m.object
}
```

**REQ-005-01**: Every populated domain type MUST expose the exact accessor:

```go
CoreObject() core.Object
```

**REQ-005-02**: Generic symbolic operations MUST accept `core.Object`, not a collection of unrelated domain interfaces.

**REQ-005-03**: Domain wrappers MUST not duplicate symbolic expression logic.

**REQ-005-04**: The MVP MUST NOT expose a generic public `NewObject(kind, ...)` factory.

**REQ-005-05**: No domain package may expose a public constructor that permits callers to choose `Kind`, `Dimension`, `Provenance`, or corpus status arbitrarily for a trusted object.

## 5.2.1 Exact internal object mint API

The internal kernel MUST expose exactly one production minting entry point for `Object` values:

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

`MintObject` MUST:

1. reject an invalid expression handle;
2. reject empty names where the constructor contract requires a name;
3. reject invalid dimensions;
4. validate provenance/hash metadata;
5. validate assumption/convention canonicality;
6. normalize all immutable metadata into canonical internal form;
7. return a fully valid immutable `Object` or a typed error.

`Object`, `Expr`, `Kind`, `Dimension`, and the metadata value types used by `ObjectSpec` are concrete kernel types. `core` aliases them for the public API.

`MintObject` exists only under `internal/kernel` and MUST NOT be re-exported from `core`.

No second object-minting helper with a different semantic contract may exist in MVP.

## 5.2 Internal implementation of `core.Object`

The actual concrete implementation of `core.Object` resides in:

```text
internal/kernel/types.go
```

`core/object.go` exposes the public alias:

```go
type Object = kernel.Object
```

`core` similarly aliases the foundational value types required by the public API:

```text
Kind
Dimension
Expr
Assumption
AssumptionSet
Convention
ConventionSet
Provenance
CorpusStatus
```

The authoritative fields of `kernel.Object` are unexported.

Logical fields are:

```text
valid
name
kind
dimension
expr
assumptions
conventions
provenance
corpusStatus
```

The zero value is invalid.

**REQ-005-06**: `kernel.Object` fields MUST remain unexported.

**REQ-005-07**: `core.Object{}` MUST be invalid.

**REQ-005-08**: No mutator may exist for authoritative object metadata.

## 5.3 Exact read-only accessors

The exact required object accessors are:

```go
Valid() bool
Name() string
Kind() Kind
Dimension() Dimension
Expr() Expr
Assumptions() AssumptionSet
Conventions() ConventionSet
Provenance() Provenance
CorpusStatus() CorpusStatus
```

The returned metadata MUST be immutable from the caller's perspective.

**REQ-005-09**: All listed accessors MUST exist with the exact names and return categories above.

## 5.4 Trusted construction

The only production code permitted to mint a valid `core.Object` is:

- fixed public constructors in populated domain packages;
- fixed candidate construction in `hypothesis` with forced `HYPOTHESIS` provenance;
- pure operation results in `ops` built from valid operation inputs;
- session-specific derived artifacts such as `Identify`, `Postulate`, `Define`, and `Declare`, subject to the exact session rules in §16;
- internal kernel test fixtures used only by tests.

All such construction routes call the internal kernel mint mechanism.

The internal mint function is intentionally located below the module's public import boundary:

```text
internal/kernel.MintObject(...)
```

The public `core` package MUST NOT forward this function.

**REQ-005-10**: The public package set MUST contain no generic trusted-object constructor.

**REQ-005-11**: All production-created valid objects MUST pass the kernel's invariant validation before becoming visible to callers.

## 5.5 Invalid-object rule

All public operations MUST reject `!Valid()` inputs with:

```text
InvalidObjectError
```

No operation may panic for an invalid input in ordinary use.

**REQ-005-12**: Invalid zero objects MUST be rejected before operation-specific work.

---

# 6. Physical kind

`Dimension` and `Kind` are intentionally independent.

The MVP `Kind` enumeration is exactly:

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

The exact names and stable ordinals MUST NOT change within `mrc-v0.4`.

**REQ-006-01**: The implementation MUST include every Kind listed above.

**REQ-006-02**: `Kind` MUST be constructor-assigned and immutable.

## 6.1 Why kind exists

Two quantities may have equal dimensions yet not be additively compatible.

For example, energy and torque share dimensions but do not automatically become the same physical category.

MRC-003 therefore checks `Kind` separately from `Dimension`.

## 6.2 Kind compatibility table

For `Add` and `Subtract`:

```text
same named Kind + same Dimension → allowed
Expression + Expression          → allowed
named Kind + different named Kind → rejected
named Kind + Expression           → rejected
```

For `Compare`:

```text
named Kind + same named Kind       → allowed when dimensions match
named Kind + Expression             → allowed when dimensions match
Expression + named Kind             → allowed when dimensions match
Expression + Expression             → allowed when dimensions match
```

For inequality operators (`lt`, `lte`, `gt`, `gte`), both sides MUST also be ordered scalar quantities under the bounded MVP ordering rule:

```text
Mass
RestMass
Time
Energy
KineticEnergy
SpeedOfLight
Expression
```

Velocity is intentionally not used by the MVP inequality path.

For `Multiply` and `Divide`, operand kinds do not need to match; the result kind is `Expression`.

For `Pow`, the result kind is `Expression`.

For `Solve`, the input relation is `Relation`, the target is any valid physical quantity whose expression is a single symbol, and the result kind is `BranchSet`.

For `SelectBranch`, the input kind MUST be `BranchSet` and the constraint kind MUST be `Relation`.

**REQ-006-03**: The compatibility table above is normative and MUST be implemented exactly.

---

# 7. Dimensions

## 7.0 Exact public Dimension API

The public `core` package MUST expose fixed constructors for all dimensions used by MVP domains:

```go
func Dimensionless() Dimension
func DimensionMass() Dimension
func DimensionLength() Dimension
func DimensionTime() Dimension
func DimensionVelocity() Dimension
func DimensionAcceleration() Dimension
func DimensionForce() Dimension
func DimensionMomentum() Dimension
func DimensionEnergy() Dimension
```

`Dimension` MUST expose:

```go
func (d Dimension) Valid() bool
func (d Dimension) Equal(other Dimension) bool
func (d Dimension) Multiply(other Dimension) Dimension
func (d Dimension) Divide(other Dimension) Dimension
func (d Dimension) Pow(exponent *big.Rat) Dimension
func (d Dimension) CanonicalJSON() ([]byte, error)
func (d Dimension) Hash() [32]byte
```

No generic public dimension parser is required.

## 7.1 Base dimensions

Use exactly seven SI base dimensions:

```text
M  mass
L  length
T  time
I  electric current
Θ  thermodynamic temperature
N  amount of substance
J  luminous intensity
```

Each exponent MUST be represented internally using `math/big.Rat`.

No floating-point value is permitted for dimension arithmetic or canonical dimension encoding.

**REQ-007-01**: Dimension arithmetic MUST use exact rational exponents backed by `math/big.Rat`.

## 7.2 Required dimensions

```text
Mass             = M
Time             = T
Position         = L
Velocity         = L T^-1
Acceleration      = L T^-2
Force            = M L T^-2
Momentum         = M L T^-1
Energy           = M L^2 T^-2
RestMass         = M
ThreeMomentum    = M L T^-1
FourMomentum     = M L T^-1
SpeedOfLight     = L T^-1
Spacetime        = L
MinkowskiMetric  = dimensionless
Expression       = operation-derived
Relation         = dimension of its left/right physical quantity when defined
BranchSet        = dimension of its target quantity
```

For dimensionless values, all seven exponents are zero.

`Relation` is not used as an operand for arithmetic dimension composition.

## 7.3 Dimension operations

The dimension implementation MUST provide deterministic operations equivalent to:

```text
Equal
Multiply
Divide
Pow
Canonical
```

`Pow` multiplies every base exponent by the exact rational exponent.

All arithmetic MUST return canonical reduced rationals with positive denominators.

**REQ-007-02**: Dimension equality MUST be structural/exact.

**REQ-007-03**: Dimension multiplication/division/power MUST be deterministic.

---

# 8. Symbolic expression representation

## 8.0 Exact public expression construction and inspection API

The public `core` package MUST expose the following constructors/helpers:

```go
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

The following read-only inspection methods MUST exist on `Expr`:

```go
Valid() bool
Kind() ExprKind
SymbolName() string
RationalValue() *big.Rat
Children() []Expr
Base() Expr
Exponent() *big.Rat
FunctionID() string
Arguments() []Expr
RelationOperator() RelationOperator
Left() Expr
Right() Expr
BranchTarget() Expr
Branches() []Expr
```

A method MAY return an empty/invalid value when the method does not apply to the node kind; operation code MUST check `Expr.Kind()` before using kind-specific accessors.

The zero `Expr` value is invalid.

`RationalValue()` MUST return a defensive copy.

`Exponent()` MUST return a defensive copy.

Slice-returning accessors MUST return copies so callers cannot mutate the internal expression tree.

## 8.1 Closed expression tree

`core.Expr` is an immutable symbolic handle backed by a closed internal representation.

The exact MVP node set is:

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

No other node type is permitted in MVP.

The following are explicitly forbidden as expression nodes:

```text
Derivative
Limit
Function
Arbitrary
Eval
Callback
RawString
```

`Differentiate` and `Limit` are operations transforming supported expressions.

`Call` is symbolic application only.

`BranchSet` is a bounded symbolic container for the two-branch quadratic solution used by the MVP; it is not a general set-algebra engine.

**REQ-008-01**: The symbolic node set MUST be exactly the list above.

**REQ-008-02**: No runtime callback or executable function may be stored inside an expression.

## 8.2 Node ordinals

The canonical node-kind ordinal ordering is:

```text
0 Symbol
1 Rational
2 Add
3 Mul
4 Neg
5 Pow
6 Sqrt
7 Call
8 Relation
9 BranchSet
```

These ordinals are part of canonicalization and MUST remain stable for `mrc-v0.4` artifacts.

## 8.2.1 Exact expression/operator enum surface

`ExprKind` MUST expose these stable values:

```text
ExprSymbol
ExprRational
ExprAdd
ExprMul
ExprNeg
ExprPow
ExprSqrt
ExprCall
ExprRelation
ExprBranchSet
```

`RelationOperator` MUST expose exactly:

```text
RelationEq
RelationNeq
RelationLt
RelationLte
RelationGt
RelationGte
```

The canonical JSON strings are:

```text
eq
neq
lt
lte
gt
gte
```

No alias operator names are accepted by canonical JSON.

## 8.3 Symbol

A Symbol contains only a non-empty identifier.

Symbol names are data, never code.

Examples:

```text
m
v
a
F
p
E
c
t
```

No unit or dimension is encoded into the symbol name.

**REQ-008-03**: Symbol content MUST never be interpreted as executable source.

## 8.4 Rational

Use `math/big.Rat` internally.

Canonical form:

```text
num/den
```

with:

- reduced numerator/denominator;
- positive denominator;
- zero normalized to `0/1`.

Exact symbolic coefficients MUST NOT use `float64`.

**REQ-008-04**: Rational canonical serialization MUST be deterministic and exact.

## 8.5 Add and Mul

`Add` and `Mul` are n-ary nodes.

They are commutative in MVP.

Canonicalization MUST:

1. flatten nested same-operation nodes;
2. remove safe identities;
3. combine exact rational coefficients;
4. combine directly repeated compatible powers where explicitly covered by §9;
5. normalize negative scalar factors according to the normal-form rule in §9.5;
6. sort children deterministically.

There is no noncommutative multiplication branch in MVP.

**REQ-008-05**: No Go map iteration order may affect expression ordering.

## 8.6 Neg

`Neg` is a first-class node.

The canonicalizer MUST apply:

```text
Neg(Neg(x)) → x
```

`Neg(x)` is the preferred canonical representation of a negative symbolic result.

## 8.7 Pow

`Pow(base, exponent)` stores an exact rational exponent.

The MVP MAY construct any rational exponent needed by the specified operations, but canonical rewrite rules remain bounded.

`Pow(x, 1)` simplifies to `x`.

`Pow(x, 0)` simplifies to `1` only when `x` is a valid nonzero-safe base under the operation's assumptions; the MVP derivations do not depend on ambiguous `0^0` behavior.

**REQ-008-06**: `Pow` MUST never silently use floating-point exponents.

## 8.8 Sqrt

`Sqrt(x)` is a dedicated symbolic node.

Canonical simplification rules are defined in §9.6.

A sign-dependent rewrite MUST NOT occur without the required assumption entailment.

## 8.9 Call

`Call` has:

```text
function_id
arguments[]
```

The only MVP function ID is:

```text
lorentz_factor
```

No arbitrary function registry exists.

The public relativity constructor produces the expression:

```text
LorentzFactor(v)
```

as a Call node.

---

# 9. Canonicalization and mechanical rewrite rules

Canonicalization is a mechanical representation rule. It is not a physical identification rule.

## 9.1 Canonical equality

`core.EqualExpr(a,b)` MUST compare canonical symbolic structure.

It MUST NOT compare display strings.

`core.EqualObject(a,b)` MUST compare canonical object representation, including authoritative physical metadata.

No equality operation may infer physical truth.

**REQ-009-01**: Expression equality MUST be structural.

**REQ-009-02**: Object equality MUST include all authoritative metadata fields that participate in canonical object serialization.

## 9.2 Canonical hashing

Use SHA-256 over canonical bytes.

Required public hash operations include:

```text
HashExpr
HashObject
HashAssumptionSet
HashConventionSet
```

The exact returned Go representation MAY be `[32]byte` with a hexadecimal helper; canonical JSON MUST use the lowercase hexadecimal form where a string hash field is required.

**REQ-009-03**: Hashes MUST be deterministic across runs, processes, and machines for identical canonical input.

## 9.3 Child sorting

For children of `Add` and `Mul`, sort by:

1. node-kind ordinal;
2. lowercase hexadecimal canonical SHA-256 child hash, compared lexicographically byte-for-byte;
3. canonical child byte sequence, compared lexicographically byte-for-byte.

The implementation MUST use `bytes.Compare`-equivalent byte-level ordering.

No pointer addresses, map iteration, timestamps, or locale-sensitive ordering are permitted.

**REQ-009-04**: Canonical child ordering MUST follow the three-key order above.

## 9.4 Rational combination

For a `Mul` node:

```text
Rational(q1) · Rational(q2) → Rational(q1*q2)
```

For an `Add` node, all direct rational terms are combined into one rational term when the exact sum is nonzero.

Examples:

```text
1/2 · 2 → 1
1/2 + 1/2 → 1
```

**REQ-009-05**: Rational combination MUST be exact.

## 9.5 Canonical scalar/sign normal form

The preferred forms are:

```text
Mul(-1, x) → Neg(x)
Neg(Rational(q)) → Rational(-q)
Neg(Neg(x)) → x
```

A negative rational coefficient on a symbolic product is represented as a `Neg` around a positive-coefficient `Mul`.

`Mul` canonical numeric coefficients are therefore non-negative unless the entire expression is a rational scalar.

## 9.6 Identity, zero, and exact rational-power rules

Constructor-time normalization is limited to structural canonicalization only: flattening same-kind `Add`/`Mul`, exact rational combination, sign normal form, and deterministic child ordering. The algebraic rewrites in this section and §9.7–§9.8 are `Simplify`-time rules.

The simplifier MUST apply these exact mechanical rules:

```text
Add(x, 0) → x
Add()     → 0
Mul(x, 1) → x
Mul()     → 1
Mul(x, 0) → 0  when x is a finite supported symbolic factor under current assumptions
Pow(x,1)  → x
Neg(Neg(x)) → x
Sqrt(1) → 1
Sqrt(0) → 0
```

For exact rational powers with an integer exponent:

```text
Pow(Rational(q), n) → Rational(q^n) for integer n >= 0
Pow(Rational(q), n) → Rational(q^n) for integer n < 0 when q != 0
Pow(Rational(0), n) → UnsupportedOperationError for integer n < 0
Pow(Rational(1), n) → Rational(1) for any exact rational exponent n
Pow(Rational(0), e) → Rational(0) for positive exact rational exponent e
```

`0^0` remains unsupported rather than being silently assigned a value.

For MVP, a **finite supported symbolic factor** is any expression composed only of `Symbol`, `Rational`, `Add`, `Mul`, `Neg`, non-negative-integer `Pow`, and sign-safe `Sqrt`, with every `Call` expanded only when its fixed MVP definition is finite under current assumptions. Primitive domain symbols are treated as finite physical quantities; the simplifier does not model IEEE floating-point infinities or NaNs.

For a symbolic denominator, `0` times an expression is accepted only when the expression is finite under the current assumptions; the required MVP derivations do not create an indeterminate `0/0` term.

## 9.7 Repeated powers

The simplifier MUST support the exact combinations required by the MVP derivations:

```text
x · x → Pow(x,2)
Pow(x,2) · Pow(x,2) → Pow(x,4)
```

and, for a fixed nonnegative integer exponent:

```text
Pow(Pow(x,a),b) → Pow(x,a*b)
```

only when `a` and `b` are non-negative integer rationals.

The implementation MUST NOT add general branch-sensitive power algebra.

## 9.8 Sqrt rewrite rules and bounded sign entailment

The simplifier MUST implement:

```text
Sqrt(Rational(q)) → exact rational root
```

when `q` is a non-negative perfect-square rational.

For symbolic terms:

```text
Sqrt(Pow(x,2)) → x
```

is permitted only when the current `AssumptionSet` entails `x >= 0`.

The MVP's bounded sign entailment is purely structural and consists of:

```text
Rational(q), q >= 0                         → nonnegative
Symbol(x), explicit x >= 0 assumption       → nonnegative
Symbol(x), explicit x > 0 assumption        → positive
Mul(f1,...,fn), every factor nonnegative    → nonnegative
Pow(x,n), even nonnegative integer n        → nonnegative
```

A product containing a positive rational coefficient and a nonnegative symbolic factor is nonnegative.

Positive entails nonzero.

No other logical inference is permitted.

This is not a theorem prover.

**REQ-009-06**: Sign-sensitive simplification MUST use only this bounded structural entailment.

## 9.9 Relation normalization

For a Relation node, recursively canonicalize the two sides and preserve the exact relation operator.

The relation operator set is:

```text
eq
neq
lt
lte
gt
gte
```

The simplifier MUST NOT replace `Relation(eq, a, b)` with an `IDENTIFIED` provenance state.

---

# 10. Canonical JSON schema

Canonical JSON rules apply to **all canonical artifacts**, not only expressions.

This includes:

```text
Expr
Object
Dimension
AssumptionSet
ConventionSet
Provenance
Manifest
OperationParams
Step
StepEnvelope
Ledger
ResearchCandidate
Prediction
FalsificationCondition
RecoveryClaim
AnomalyReference
Review
Challenge
```

## 10.1 Canonical JSON requirements

**REQ-010-01**: Canonical artifacts MUST use explicit typed structs.

**REQ-010-02**: Canonical artifacts MUST NOT use `map[string]any` or equivalent dynamic maps.

**REQ-010-03**: Canonical field order MUST be fixed by the struct/encoder declaration order and remain stable.

**REQ-010-04**: JSON numeric values MUST NOT encode exact symbolic rationals or dimension exponents.

**REQ-010-05**: Exact rationals MUST serialize as canonical strings such as `"1/2"`.

**REQ-010-06**: Canonical artifacts MUST NOT contain timestamps, random IDs, pointer addresses, environment-dependent fields, or map-order-dependent values.

## 10.2 Expression JSON

Canonical expression forms are exactly:

```json
{"kind":"symbol","name":"E"}
```

```json
{"kind":"rational","value":"1/2"}
```

```json
{"kind":"add","terms":[...]}
```

```json
{"kind":"mul","factors":[...]}
```

```json
{"kind":"neg","expr":{...}}
```

```json
{"kind":"pow","base":{...},"exp":"2/1"}
```

```json
{"kind":"sqrt","expr":{...}}
```

```json
{"kind":"call","function":"lorentz_factor","args":[...]}
```

```json
{"kind":"relation","operator":"eq","lhs":{...},"rhs":{...}}
```

```json
{"kind":"branch_set","target":{...},"branches":[...]}
```

No other expression JSON kind is allowed.

## 10.3 Dimension JSON

Dimension JSON field order MUST be:

```json
{
  "m":"1/1",
  "l":"0/1",
  "t":"0/1",
  "i":"0/1",
  "theta":"0/1",
  "n":"0/1",
  "j":"0/1"
}
```

## 10.4 Object JSON

Canonical object JSON field order MUST be:

```text
schema_version
valid
name
kind
dimension
expr
assumptions
conventions
provenance
corpus_status
```

`valid` is always `true` for a serialized valid object.

The canonical object representation is retained by the session ledger for replay.

## 10.5 Round-trip requirement

The following MUST reproduce byte-identical canonical form and identical hash:

```text
Expr
→ canonical JSON
→ Expr
```

and, through the internal replay decoder:

```text
Object
→ canonical JSON
→ Object
```

**REQ-010-07**: Canonical round-trip tests MUST exist for expressions and representative object artifacts.

---

# 11. Assumption system

## 11.0 Exact assumption metadata API

`AssumptionSet` is immutable and MUST expose:

```go
func NewTextAssumption(kind AssumptionKind, key, value string) (Assumption, error)
func NewExprAssumption(kind AssumptionKind, key string, value Expr) (Assumption, error)
func NewAssumptionSet(values ...Assumption) AssumptionSet

func (s AssumptionSet) Values() []Assumption
func (s AssumptionSet) Merge(other AssumptionSet) (AssumptionSet, error)
func (s AssumptionSet) Equal(other AssumptionSet) bool
func (s AssumptionSet) CanonicalJSON() ([]byte, error)
func (s AssumptionSet) Hash() [32]byte
func (s AssumptionSet) EntailsNonNegative(expr Expr) bool
func (s AssumptionSet) EntailsPositive(expr Expr) bool
func (s AssumptionSet) EntailsNonZero(expr Expr) bool
```

`Assumption` MUST expose read-only accessors for `Kind`, `Key`, and typed Value mode/content.

`AssumptionSet.Values()` MUST return a copy.

## 11.1 Assumption kinds

The exact MVP assumption kind values are:

```text
AssumptionDomain
AssumptionRegime
AssumptionConstraint
AssumptionConvention
AssumptionApproximation
AssumptionMathPrecondition
AssumptionPhysicalAssumption
```

Canonical JSON strings are:

```text
domain
regime
constraint
convention
approximation
math_precondition
physical_assumption
```

The exact MVP assumption kinds are:

```text
Domain
Regime
Constraint
Convention
Approximation
MathPrecondition
PhysicalAssumption
```

`Approximation` is reserved in MVP; no MVP operation produces `APPROXIMATED` provenance.

## 11.2 Structured assumption value

An assumption MUST NOT encode a physical relationship in a free-form string.

Each assumption has:

```text
Kind
Key
Value
```

where `Value` is the following closed union:

```text
TextValue(string)
ExprValue(core.Expr)
```

Canonical JSON is typed, for example:

```json
{
  "kind":"constraint",
  "key":"momentum_zero",
  "value":{
    "mode":"expr",
    "expr":{
      "kind":"relation",
      "operator":"eq",
      "lhs":{"kind":"symbol","name":"p"},
      "rhs":{"kind":"rational","value":"0/1"}
    }
  }
}
```

TextValue is permitted for non-mathematical metadata such as regime identifiers.

**REQ-011-01**: Mathematical assumptions MUST use structured expression values, never equation strings.

## 11.3 Construction

Public helpers MAY construct assumptions because assumptions are metadata, not trusted physical objects.

The core API MUST provide deterministic constructors equivalent to:

```text
NewTextAssumption(kind, key, text)
NewExprAssumption(kind, key, expr)
NewAssumptionSet(...)
```

Keys MUST be non-empty.

## 11.4 Merge

When combining objects:

```text
result assumptions =
    union(all input assumptions)
    + required operation assumptions
    + explicitly introduced assumptions
```

Exact duplicate assumptions are deduplicated.

Two assumptions with the same `(Kind, Key)` and different canonical values MUST yield:

```text
AssumptionConflictError
```

No subsumption or theorem proving is performed.

**REQ-011-02**: Assumption merge MUST be deterministic set union plus exact conflict detection.

## 11.5 Denominator preconditions

A symbolic division by `b` adds:

```text
Kind = MathPrecondition
Key = "denominator/" + HashExpr(b.Expr())
Value = ExprValue(Compare(b, 0, neq))
```

The key therefore differs for different denominators.

This prevents unrelated divisions from silently deduplicating or falsely conflicting.

## 11.6 Sign assumptions used in MVP

`RestMass()` MUST carry:

```text
Constraint
Key: rest_mass_nonnegative
Value: m >= 0
```

`SpeedOfLight()` MUST carry:

```text
Constraint
Key: speed_of_light_positive
Value: c > 0
```

These are explicit assumptions carried by the corresponding objects.

The mass-energy simplification uses them through the bounded entailment rules in §9.8.

---

# 12. Conventions

## 12.0 Exact convention metadata API

`ConventionSet` MUST expose:

```go
func NewConvention(key, value string) (Convention, error)
func NewConventionSet(values ...Convention) ConventionSet

func (s ConventionSet) Values() []Convention
func (s ConventionSet) Merge(other ConventionSet) (ConventionSet, error)
func (s ConventionSet) Equal(other ConventionSet) bool
func (s ConventionSet) CanonicalJSON() ([]byte, error)
func (s ConventionSet) Hash() [32]byte
```

`Values()` MUST return a copy.


Conventions are separate immutable metadata from assumptions, even though `AssumptionKind.Convention` remains a reserved category for future corpus work.

A Convention is:

```text
Key
Value
```

Both are non-empty strings.

MVP relativity convention:

```text
metric.signature = -+++
```

A convention conflict means same key, different canonical value.

The result is:

```text
ConventionConflictError
```

Convention sets participate in canonical serialization and object hashing.

**REQ-012-01**: Convention conflicts MUST be rejected deterministically.

**REQ-012-02**: No physical equation may be encoded solely as a convention string.

---

# 13. Provenance and corpus status

## 13.0 Exact provenance API

`Provenance` MUST expose read-only accessors:

```go
Status() ProvenanceStatus
Source() string
Framework() string
ParentHashes() []string
AssumptionHash() [32]byte
ConventionHash() [32]byte
MRCVersion() string
Justification() string
```

`ParentHashes()` MUST return a copy.

## 13.0.1 Exact provenance construction helpers

The public `core` package MUST expose read-only status values and one deterministic constructor for provenance metadata:

```go
type ProvenanceStatus string

func NewProvenance(
    status ProvenanceStatus,
    source string,
    framework string,
    parentHashes []string,
    assumptionHash [32]byte,
    conventionHash [32]byte,
    mrcVersion string,
    justification string,
) (Provenance, error)
```

`NewProvenance` validates status-specific invariants but does not mint a `core.Object`.

Object minting remains exclusively governed by §5.4 and `internal/kernel.MintObject`.

## 13.1 Provenance statuses

The exact provenance status enum is:

```text
DEFINED
POSTULATED
DERIVED
IDENTIFIED
APPROXIMATED
HYPOTHESIS
```

`APPROXIMATED` is reserved in MVP and is exempt from the MVP status-coverage requirement because no MVP operation produces it.

It becomes active when series expansion/truncation is implemented in v0.5+.

**REQ-013-01**: `APPROXIMATED` MUST exist in the enum but MUST NOT be minted by an MVP operation.

## 13.2 Deterministic provenance propagation

The MVP uses the following exact propagation law.

For every pure `ops` transformation:

```text
if any input status == HYPOTHESIS:
    output status = HYPOTHESIS
else:
    output status = DERIVED
```

This applies to:

```text
Add
Subtract
Multiply
Divide
Pow
Simplify
Substitute
Differentiate
Limit
Compare
Solve
SelectBranch
```

For `Session.Identify`:

```text
if either input status == HYPOTHESIS:
    output status = HYPOTHESIS
else:
    output status = IDENTIFIED
```

For session assertion actions:

```text
Postulate → records a valid object whose provenance status is POSTULATED
Define    → records a valid object whose provenance status is DEFINED
Declare   → preserves the input's existing provenance status
```

No MVP production domain constructor is required to emit `POSTULATED`; the `Postulate` action is exercised with an `internal/kernel` test fixture whose provenance is explicitly `POSTULATED`. `Define` is exercised by the fixed domain constructors.

No MVP operation may create `IDENTIFIED` except `Session.Identify`.

No MVP operation may raise `HYPOTHESIS` to any trusted status.

**REQ-013-02**: Provenance propagation MUST be deterministic and MUST follow this exact status law.

## 13.3 Provenance record

Every object provenance record contains:

```text
Status
Source
Framework
ParentHashes[]
AssumptionHash
ConventionHash
MRCVersion
Justification
```

`Justification` is non-empty only for `IDENTIFIED` or a recorded identification attempt involving `HYPOTHESIS`.

`ParentHashes` preserve exact input object hashes in operation input order.

`MRCVersion` is always:

```text
mrc-v0.4
```

for MVP artifacts.

**REQ-013-03**: Provenance parent references MUST be based on canonical object hashes.

## 13.4 Corpus status

The exact corpus-status values are:

```text
NONE
ESTABLISHED
CONTESTED
SUPERSEDED
FALSIFIED
```

Canonical JSON strings are lowercase snake-case equivalents.

Corpus status is a separate human-curated axis:

```text
NONE
ESTABLISHED
CONTESTED
SUPERSEDED
FALSIFIED
```

`NONE` is used for derived and candidate artifacts unless a specific corpus constructor specifies otherwise.

The populated manifest files are the human-curated source of corpus status.

The implementing agent MUST NOT infer or revise corpus status.

For MVP baseline manifests:

```text
mechanics framework corpus status = ESTABLISHED
special_relativity framework corpus status = ESTABLISHED
```

These values are repository metadata, not a runtime judgment by the library.

**REQ-013-04**: The library MUST never compute corpus status.

**REQ-013-05**: No API may expose any of the following as a generated physical judgment:

```text
PHYSICALLY_TRUE
TRUTH_SCORE
PROBABILITY_OF_TRUTH
BEST_THEORY
```

## 13.5 Source and framework inheritance

For derived results:

- if all input `Source` values are identical and non-empty, preserve that source;
- otherwise use empty source;
- if all input `Framework` values are identical and non-empty, preserve that framework;
- otherwise use empty framework.

No semantic inference occurs.

---

# 14. MRC version and rules

The MVP MRC version is:

```text
mrc-v0.4
```

Each rule identifier is stable.

The exact MVP rule set is:

```text
MRC-001 Constructor/carrier integrity
MRC-002 Dimensional compatibility
MRC-003 Physical-kind compatibility
MRC-004 Assumption compatibility
MRC-005 Convention compatibility
MRC-006 Explicit physical identification
MRC-007 Session/provenance authority
MRC-008 Candidate containment
```

## 14.1 MRC-001 — Constructor/carrier integrity

A valid trusted object MUST originate through an allowed construction path from §5.4.

`core.Object` authoritative fields are unexported behind `internal/kernel`.

There is no public generic object factory.

The zero object is invalid.

Primary enforcement:

```text
Go package visibility
internal/ package boundary
kernel invariant validation
fixed domain constructors
pure operation outputs
session authority
```

Failure:

```text
InvalidObjectError
```

or, where the attempted construction path itself is being rejected by a public API:

```text
ProvenanceError
```

## 14.2 MRC-002 — Dimensional compatibility

`Add` and `Subtract` require equal dimensions.

Every `Compare` operator requires equal dimensions.

`Multiply` multiplies dimensions.

`Divide` divides dimensions.

`Pow` raises dimensions by the exact rational exponent.

`Substitute` requires the replacement dimension to equal the variable dimension.

`Differentiate` uses dimension calculus:

```text
d(target)/d(wrt) dimension = target dimension / wrt dimension
```

The resulting object has that derived dimension.

Failure:

```text
DimensionMismatchError
```

## 14.3 MRC-003 — Physical-kind compatibility

Implement the exact §6.2 compatibility table.

Failure:

```text
CategoryMismatchError
```

## 14.4 MRC-004 — Assumption compatibility

Operations merge assumptions according to §11.

Conflicting keyed assumptions fail.

Failure:

```text
AssumptionConflictError
```

## 14.5 MRC-005 — Convention compatibility

Operations merge conventions and reject conflicting keys.

Failure:

```text
ConventionConflictError
```

## 14.6 MRC-006 — Explicit physical identification

Only `Session.Identify` creates an `IDENTIFIED` provenance result.

It requires:

- two valid explicit operands;
- compatible dimensions;
- compatible kinds according to Compare-equivalence rules;
- non-empty justification.

The identification is logged in the session draft buffer.

A justification is reviewable reasoning, not a proof of physical truth.

If either operand is `HYPOTHESIS`, output remains `HYPOTHESIS` under MRC-008.

Failure:

```text
IdentifyError
```

## 14.7 MRC-007 — Session/provenance authority

Only `session.Session` may create a committed derivation ledger and seal a `ResearchCandidate`.

No exported arbitrary-field derivation constructor exists.

All committed steps MUST pass through Session.

Failure:

```text
ProvenanceError
```

## 14.8 MRC-008 — Candidate containment

If any operation result depends directly on an object with `HYPOTHESIS` provenance, the result MUST be `HYPOTHESIS`.

Containment is transitive through all downstream operations and session steps.

An attempted artifact state in which a hypothesis-dependent result is presented with a trusted provenance status MUST fail validation with:

```text
CandidateContainmentError
```

There is no promotion API.

---

# 15. Operation API and common operation semantics

All pure operations:

- never mutate inputs;
- validate object validity first;
- validate applicable MRC rules;
- merge assumptions and conventions;
- create a fresh immutable result;
- propagate provenance deterministically;
- preserve deterministic canonical form.

They do not write to any ledger and do not read ambient/session state.

**REQ-015-01**: Every pure operation MUST be free of global/ambient session state.

## 15.1 Add

Exact signature:

```go
func Add(a, b core.Object) (core.Object, error)
```

Checks:

```text
MRC-001
MRC-002
MRC-003
MRC-004
MRC-005
MRC-008
```

Result:

```text
same Kind as operands
same Dimension
DERIVED unless HYPOTHESIS contamination applies
```

`Expression + Expression` returns `Expression`.

## 15.2 Subtract

Exact signature:

```go
func Subtract(a, b core.Object) (core.Object, error)
```

Same compatibility contract as Add.

## 15.3 Multiply

Exact signature:

```go
func Multiply(a, b core.Object) (core.Object, error)
```

Result kind:

```text
Expression
```

Result dimension:

```text
a.dimension * b.dimension
```

## 15.4 Divide

Exact signature:

```go
func Divide(a, b core.Object) (core.Object, error)
```

Result kind:

```text
Expression
```

Result dimension:

```text
a.dimension / b.dimension
```

Symbolic denominator assumptions are added exactly as §11.5 specifies.

## 15.5 Pow

Exact signature:

```go
func Pow(base core.Object, exponent *big.Rat) (core.Object, error)
```

The implementation MUST clone/read the input exponent and MUST NOT retain a mutable caller-owned pointer.

Result kind:

```text
Expression
```

Result dimension:

```text
base.dimension ^ exponent
```

The exact exponent is serialized canonically.

## 15.6 Simplify

Exact signature:

```go
func Simplify(x core.Object) (core.Object, error)
```

Simplify is mechanical only.

It MUST NOT:

```text
create IDENTIFIED provenance
create a physical identification relation
change corpus status
invoke Session.Identify
consult session state
```

Its output status is `DERIVED` unless hypothesis contamination applies.

## 15.7 Substitute

Exact signature:

```go
func Substitute(target, variable, replacement core.Object) (core.Object, error)
```

Requirements:

- `variable.Expr()` MUST be a single `Symbol` node;
- replacement dimension MUST equal variable dimension;
- variable and replacement kind MUST be identical;
- all occurrences of the exact variable symbol MUST be replaced in the target expression;
- target kind remains unchanged;
- assumptions/conventions/provenance are merged according to general rules.

## 15.8 Differentiate

Exact signature:

```go
func Differentiate(target, wrt core.Object) (core.Object, error)
```

`wrt.Expr()` MUST be a single `Symbol` node.

The bounded MVP derivative engine supports only:

```text
Rational        → 0
Symbol          → 1 if same symbol, otherwise 0
Add             → derivative of each term
Mul             → n-ary product rule
Neg             → Neg(derivative)
Pow             → integer-constant power rule
```

For `Pow(base,n)` where `n` is a non-negative integer rational:

```text
d(base^n)/dx = n * base^(n-1) * d(base)/dx
```

The result MUST be simplified after differentiation.

Unsupported:

```text
Sqrt
Call
Relation
BranchSet
non-integer symbolic exponent
```

These return `UnsupportedOperationError`.

Every successful `Differentiate` call returns a `core.Object` with `Kind=Expression` and the exact derivative dimension computed from the target dimension divided by the `wrt` dimension.

The MVP does not create a `Derivative` node.

## 15.9 Limit

Exact signature:

```go
func Limit(target, variable, value core.Object) (core.Object, error)
```

The MVP performs direct substitution followed by simplification.

For `Call(function="lorentz_factor")`, it MUST expand the fixed definition in §15.9.1 first.

No L'Hôpital algorithm exists.

No series expansion exists.

No general asymptotic engine exists.

Unsupported or singular results return `UnsupportedOperationError`.

### 15.9.1 Fixed Lorentz-factor function body

The only MVP function ID is:

```text
lorentz_factor
```

Its fixed internal definition is:

```text
1 / Sqrt(1 - Pow(v/c, 2))
```

The body is an ordinary existing expression tree built from the allowed nodes.

There is no general runtime function-definition registry.

The relativity object carries the `c > 0` assumption from `SpeedOfLight()`.

For the required limit:

```text
Limit(LorentzFactor(), Velocity(), ZeroVelocity())
```

or its exact object-equivalent variable/value forms MUST:

1. expand the body;
2. replace `v` with `0`;
3. simplify exact zeros, powers, square roots, and the denominator;
4. produce canonical rational `1/1`.

The implementation MUST NOT simply return `1` because the function ID is `lorentz_factor`; it must execute the fixed body path.

Every successful `Limit` call returns a `core.Object` with `Kind=Expression` and the exact resulting limit dimension.

## 15.10 Compare

Exact signature:

```go
func Compare(a, b core.Object, op core.RelationOperator) (core.Object, error)
```

All operators require equal dimensions.

`eq` and `neq` require kinds compatible under §6.2.

Inequalities require the ordered-kind rule in §6.2.

Result:

```text
Kind = Relation
Dimension = dimension of the compared quantities
```

The result provenance is `DERIVED` unless hypothesis contamination applies.

`Compare` does not decide whether an equality is physically true. It constructs a relation artifact.

## 15.11 Solve

Exact signature:

```go
func Solve(relation core.Object, target core.Object) (core.Object, error)
```

The MVP solver is pattern matching only.

Accepted exact input shape:

```text
Relation(eq,
    Pow(Symbol(target), 2),
    Expr
)
```

The target expression MUST be a single Symbol matching the squared Symbol in the relation.

Result:

```text
BranchSet(
    target,
    [
        Sqrt(rhs),
        Neg(Sqrt(rhs))
    ]
)
```

The result kind is `BranchSet` and its dimension is the target dimension.

Any other equation form returns `UnsupportedOperationError`.

No linear solver, polynomial root finder, matrix solver, or general equation isolation is permitted.

## 15.12 SelectBranch

Exact signature:

```go
func SelectBranch(branches, constraint core.Object) (core.Object, error)
```

The MVP accepts exactly:

```text
BranchSet(target, [positiveBranch, Neg(positiveBranch)])
```

and a constraint of the form:

```text
Relation(gte, target, zero)
```

where `zero` has the same kind and dimension as target.

The selected branch is the first non-negative branch. The selected constraint is also recorded as a structured assumption with deterministic key:

```text
selected_branch/<HashExpr(constraint.Expr())>
```

The operation MUST:

1. validate the branch set shape;
2. validate constraint compatibility;
3. merge the constraint into assumptions;
4. simplify the selected branch under the current assumptions;
5. preserve provenance, subject to hypothesis contamination;
6. return the selected solution as `Expression` with the target dimension.

The MVP MUST support the mass-energy positive-energy branch and MUST simplify:

```text
Sqrt(Pow(m*c^2,2)) → m*c^2
```

because `m >= 0` is carried by `RestMass` and `c^2` is structurally nonnegative.

No arbitrary branch selection rules are permitted.

## 15.13 Pure-operation dispatch

`OperationParams` has this exact field order and canonical JSON shape:

```json
{
  "kind":"pow",
  "exponent":"2/1",
  "operator":"eq",
  "justification":""
}
```

Unused fields are encoded with their empty/default values. The only permitted `OperationParams.Kind` values are:

```text
empty
pow
compare
identify
```

`pow` requires `Exponent`; `compare` requires `Operator`; `identify` requires `Justification`; `empty` is used for all other operations and assertion actions. `identify` is stored in a session step but is never accepted by `ops.Apply`.

`ops` MUST provide one fixed dispatch mechanism for session replay.

Exact conceptual API:

```go
type OperationID string

type OperationParams struct {
    Kind          string
    Exponent      string
    Operator      core.RelationOperator
    Justification string
}

func Apply(id OperationID, inputs []core.Object, params OperationParams) (core.Object, error)
```

The exact `OperationID` values are:

```text
add
subtract
multiply
divide
pow
simplify
substitute
differentiate
limit
compare
solve
select_branch
```

`identify` is deliberately excluded from `ops.Apply`; identification is session-owned.

The dispatch switch MUST be closed over exactly those IDs.

No plugin registry or runtime operation registration exists.

---

# 16. Session and derivation ledger

The ledger/session layer exists in `session/` specifically to prevent an import cycle between the core representation and pure operations.

`core` has no ledger implementation.

`session` imports both `core` and `ops` and is therefore the replay/orchestration boundary.

## 16.1 Exact nine session actions

The session supports exactly:

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

`Draft` and `Commit` are explicit state transitions/actions.

## 16.2 Session lifecycle

The MVP has this exact state machine:

```text
New
 ↓
Drafting
 ↓
Committed
 ↓
Concluded
 ↓
Sealed
```

Rules:

- `Draft` may be called exactly once and transitions `New → Drafting`.
- `Postulate`, `Declare`, `Define`, `Step`, and `Identify` are legal only in `Drafting`.
- `Commit` is legal only in `Drafting` and transitions `Drafting → Committed`. An empty draft MUST fail with `LedgerValidationError` and remain `Drafting`.
- `Conclude` is legal only in `Committed` and transitions to `Concluded`.
- `Seal` is legal only in `Concluded` when candidate metadata contains a valid `HYPOTHESIS`; otherwise it MUST fail with `ProvenanceError` and leave the session `Concluded`.
- A pure derivation without a candidate ends at `Concluded` and is validated with `Session.Validate()`; it is not sealed as a `ResearchCandidate`.
- `Validate` is legal in every session state, never mutates state, and in `New` validates only the empty-session invariants.
- After `Sealed`, all mutation actions fail with `ProvenanceError`.

**REQ-016-01**: The session state machine MUST match the lifecycle above.

## 16.2.1 Exact candidate-draft metadata

The `Draft` action carries all candidate metadata required for final sealing so the session needs no additional hidden mutation actions.

Exact conceptual type:

```go
type DraftMetadata struct {
    Hypothesis                 core.Object
    Premises                   []core.Object
    Assumptions                core.AssumptionSet
    FrameworkDependencies      []FrameworkDependency
    Predictions                []Prediction
    FalsificationConditions    []FalsificationCondition
    RecoveryClaims             []RecoveryClaim
    AnomalyReferences          []AnomalyReference
    ReviewHistory              []Review
}
```

For a pure mechanical derivation with no research candidate, `Hypothesis` MAY be the invalid zero `core.Object`.

For an MVP sealed `ResearchCandidate`, `Hypothesis` MUST be valid and `HYPOTHESIS`; `Seal` MUST reject an invalid or non-`HYPOTHESIS` hypothesis with `ProvenanceError`.

All slices are copied at Draft time.

The session MUST NOT expose separate candidate-metadata mutation methods in MVP.

## 16.3 Exact session method signatures

The session MUST expose methods equivalent to:

```go
func New() *Session

func (s *Session) Draft(derivationID, label string, metadata DraftMetadata) error

func (s *Session) Postulate(label string, object core.Object) error

func (s *Session) Declare(label string, object core.Object) error

func (s *Session) Define(label string, object core.Object) error

func (s *Session) Step(label string, operation ops.OperationID, inputs []core.Object, params ops.OperationParams) (core.Object, error)

func (s *Session) Identify(a, b core.Object, justification string) (core.Object, error)

func (s *Session) Conclude(object core.Object) error

func (s *Session) Commit() error

func (s *Session) Seal() (ResearchCandidate, error)

func (s *Session) Validate() error

func (s *Session) CanonicalJSON() ([]byte, error)
```

Exact pointer/value receiver details MAY follow ordinary Go conventions, but the semantic signatures above are fixed.

## 16.4 Draft

`Draft` records:

```text
DerivationID
Label
MRCVersion
GenesisHash
```

It clears no existing state because a session starts empty.

`DerivationID` MUST be non-empty and caller-supplied. It is therefore deterministic and not auto-generated.

## 16.5 Assertion actions

`Postulate(label,obj)` requires:

```text
obj.Valid() == true
obj.Provenance().Status == POSTULATED
```

MVP does not require a public domain constructor that emits `POSTULATED`; the action is exercised with an `internal/kernel` test fixture.

`Define(label,obj)` requires:

```text
obj.Valid() == true
obj.Provenance().Status == DEFINED
```

`Declare(label,obj)` requires only valid object and preserves existing provenance.

These actions record assertion steps without invoking `ops.Apply`.

During replay, assertion steps validate that the retained canonical object is unchanged.

## 16.6 Step execution

`Session.Step` MUST:

1. validate session state;
2. validate every input object;
3. reject `operation == identify`;
4. call the fixed `ops.Apply` dispatch;
5. create a canonical step record containing exact inputs, params, and output;
6. append the step to the draft buffer;
7. return the operation result.

The session does not permit callers to supply an arbitrary output object for a mechanical operation.

The operation result comes from `ops.Apply` and is then recorded.

## 16.7 Identify ownership

There is **no** package-level:

```text
ops.Identify
```

`Session.Identify` is the only identification API.

It MUST:

1. validate session state;
2. validate both objects;
3. validate dimensions and kinds under MRC-006;
4. require non-empty justification after trimming whitespace;
5. construct a Relation(eq, a.Expr(), b.Expr());
6. set provenance to `IDENTIFIED` unless hypothesis contamination forces `HYPOTHESIS`;
7. preserve the merged assumptions/conventions;
8. record the justification in provenance;
9. append an identification step to the draft buffer;
10. return the immutable relation object.

It MUST NOT consult or modify any global active-session state.

## 16.8 Draft buffer

`Step` and `Identify` append fully specified step material to the session's draft buffer.

Each draft entry contains enough material to later form a committed hash-chain record.

No draft entry is externally mutable through the public API.

## 16.9 Commit

`Commit` freezes the draft buffer into the committed ledger in recording order.

Commit MUST:

1. assign deterministic indices beginning at one;
2. assign `StepID` as:

```text
step-000001
step-000002
...
```

3. set the genesis/previous hash chain;
4. compute each current step hash;
5. preserve canonical inputs/outputs/params;
6. make the committed ledger immutable through the public API.

After commit, the draft buffer MUST be empty.

## 16.10 Conclude

`Conclude(object)` is legal only after `Commit`.

The supplied object MUST:

- be valid;
- have an object hash equal to the last committed step's output hash;
- become the session conclusion.

The conclusion is not a second derivation step.

## 16.11 Seal

`Seal` is legal only after `Conclude`.

Before sealing, the session MUST run full validation.

On success:

- the ledger becomes immutable;
- `DerivationHash` is fixed as the final committed current step hash;
- the session becomes sealed;
- a `ResearchCandidate` MAY be minted only if the candidate-required fields in §18 are present.

`Seal` returns a `ResearchCandidate` for the MVP's candidate handoff test.

## 16.12 Step structure

Each committed `Step` MUST contain these semantic fields in this exact order:

```text
StepID
Index
Label
StepKind
Operation
InputHashes
InputCanonicals
ParamsCanonical
OutputHash
OutputCanonical
AssumptionHash
ConventionHash
ProvenanceStatus
MRCVersion
PreviousStepHash
CurrentStepHash
```

`StepKind` values:

```text
Assertion
Transformation
Identification
```

Operations for assertion steps are:

```text
postulate
declare
define
```

Transformation operations are the exact `ops.OperationID` values in §15.13.

Identification operation is:

```text
identify
```

## 16.13 Input retention

Hashes are not invertible and therefore are not sufficient for replay.

Each committed step MUST retain:

```text
InputCanonicals[]
OutputCanonical
ParamsCanonical
```

alongside the corresponding hashes.

For every input:

```text
HashObject(decoded InputCanonical) == InputHash
```

MUST hold.

For output:

```text
HashObject(decoded OutputCanonical) == OutputHash
```

MUST hold.

This is the replay substrate.

## 16.14 Operation parameters

`ParamsCanonical` MUST encode operation parameters needed for deterministic replay.

Examples:

For `Pow`:

```text
Exponent = "2/1"
```

For `Compare`:

```text
Operator = "eq"
```

For `Identify`:

```text
Justification = non-empty string
```

For operations with no parameters, `OperationParams` uses its fixed empty/default values.

No parameter may be omitted when needed to replay the operation.

## 16.15 Genesis hash

The exact genesis previous-step hash is 64 lowercase zero characters:

```text
0000000000000000000000000000000000000000000000000000000000000000
```

The first committed step's `PreviousStepHash` MUST equal this value.

## 16.16 StepID

`StepID` is deterministic and derived solely from `Index`.

No random or UUID-based step IDs exist in MVP.

## 16.17 StepEnvelope hash

The current step hash is:

```text
SHA256(
  CanonicalJSON(
    StepEnvelope{
      PreviousHash: previousStepHash,
      Step:         canonicalStepBody
    }
  )
)
```

`canonicalStepBody` contains all Step fields except `CurrentStepHash` itself.

This avoids recursive self-reference.

The hash-chain is an integrity mechanism.

It is **not** an authenticity signature.

Anyone who can rewrite the entire ledger can recompute the hashes.

Authenticity/signing is outside MVP.

## 16.18 Derivation hash

For a non-empty committed ledger:

```text
DerivationHash = CurrentStepHash of the last committed step
```

The sealed candidate carries this hash.

## 16.19 Replay validation

`Session.Validate` MUST perform, in order:

1. validate session structure/state;
2. validate genesis previous hash;
3. validate every StepID/index pair;
4. decode every retained input canonical object;
5. verify every InputHash;
6. decode every retained output canonical object;
7. verify every OutputHash;
8. verify assumption and convention hashes;
9. verify the current hash chain by recomputing every StepEnvelope hash;
10. verify operation parameters against the operation ID;
11. replay every assertion step as a retained-object consistency check;
12. replay every transformation through `ops.Apply`;
13. replay every identification through the session-owned identification construction path;
14. compare replayed output hash to retained OutputHash;
15. validate provenance propagation and MRC version;
16. validate conclusion hash;
17. validate candidate-containment invariants.

A mismatch MUST return:

```text
LedgerValidationError
```

or the more specific typed error wrapped by it.

## 16.20 Tamper behavior

Changing only `OutputCanonical` MUST cause detection through hash mismatch.

Changing `OutputCanonical` and recomputing only `OutputHash` MUST still be detected by current-step hash mismatch.

Changing the retained output and recomputing all hashes MUST still be detected by replay divergence unless the changed data happens to be exactly the same valid operation output.

This is the primary anti-tampering test.

---

## 16.21 Exact Ledger value surface

The public ledger value type is immutable after commit.

It MUST expose:

```go
type Ledger struct { /* unexported storage */ }

func (l Ledger) Steps() []Step
func (l Ledger) DerivationHash() string
func (l Ledger) Validate() error
func (l Ledger) CanonicalJSON() ([]byte, error)
func ParseLedgerJSON(data []byte) (Ledger, error)
```

`Step` itself MUST have unexported storage with read-only accessors for every field listed in §16.12.

`ParseLedgerJSON` loads an artifact for validation; it does not create session authority or mint trusted objects outside the loaded record.

# 17. Corpus manifests

Each populated domain package has exactly one canonical machine-readable manifest:

```text
mechanics/manifest.json
relativity/manifest.json
```

The manifest bytes are loaded with `go:embed`.

## 17.0 Exact public manifest API

`core/corpus.go` MUST define typed structures equivalent to:

```go
type Manifest struct {
    SchemaVersion  string
    FrameworkID    string
    FrameworkName  string
    CorpusStatus   CorpusStatus
    Assumptions    []Assumption
    Domain         []ManifestDomain
    Limits         []ManifestLimit
    Anomalies      []ManifestAnomaly
    Items          []ManifestItem
}

type ManifestDomain struct {
    ID          string
    Name        string
    Description string
}

type ManifestLimit struct {
    ID          string
    Description string
}

type ManifestAnomaly struct {
    ID                    string
    Framework             string
    Description           string
    Status                string
    RelatedAssumptions    []string
    RelatedItems          []string
    ResearchRelevance     string
}

type ManifestReduction struct {
    ID        string
    Condition string
    Result    string
}

type ManifestItem struct {
    ID                      string
    Constructor             string
    Kind                    string
    Name                    string
    Statement               string
    CanonicalExpr           Expr
    Dimension               Dimension
    ProvenanceStatus        ProvenanceStatus
    Source                  string
    Assumptions             []Assumption
    DerivableFrom           []string
    ReducesTo               []ManifestReduction
    KnownLimits             []string
    Anomalies               []string
    FalsificationConditions []string
}
```

The package MUST expose:

```go
func ParseManifest(data []byte) (Manifest, error)
func ValidateManifestBytes(data []byte) (Manifest, error)
func CanonicalManifestJSON(m Manifest) ([]byte, error)
func (m Manifest) Hash() [32]byte
```

`ParseManifest` and `ValidateManifestBytes` perform only schema/canonical-form validation. Constructor cross-checking occurs in the domain manifest tests.

## 17.1 Manifest top-level schema

The exact typed structure is:

```text
Manifest
    schema_version
    framework_id
    framework_name
    corpus_status
    assumptions[]
    domain[]
    limits[]
    anomalies[]
    items[]
```

No dynamic maps are used in the canonical in-memory schema.

## 17.2 Domain entries

A domain entry contains:

```text
id
name
description
```

## 17.3 Limit entries

A limit entry contains:

```text
id
description
```

## 17.4 Anomaly entries

An anomaly entry contains:

```text
id
framework
description
status
related_assumptions[]
related_items[]
research_relevance
```

Status is descriptive metadata, not automatically generated scientific truth.

## 17.5 Manifest item schema

Each item contains, in order:

```text
id
constructor
kind
name
statement
canonical_expr
dimension
provenance_status
source
assumptions
derivable_from
reduces_to
known_limits
anomalies
falsification_conditions
```

`constructor` is a non-empty string naming the Go constructor identifier.

Example:

```json
{
  "id":"NewtonSecondLaw",
  "constructor":"mechanics.NewtonSecondLaw",
  "kind":"Relation",
  "name":"Newton's second law",
  "statement":"Force equals mass times acceleration.",
  "canonical_expr":{
    "kind":"relation",
    "operator":"eq",
    "lhs":{"kind":"symbol","name":"F"},
    "rhs":{"kind":"mul","factors":[
      {"kind":"symbol","name":"m"},
      {"kind":"symbol","name":"a"}
    ]}
  },
  "dimension":{
    "m":"1/1",
    "l":"1/1",
    "t":"-2/1",
    "i":"0/1",
    "theta":"0/1",
    "n":"0/1",
    "j":"0/1"
  },
  "provenance_status":"DEFINED",
  "source":"Newton, Principia",
  "assumptions":[],
  "derivable_from":[],
  "reduces_to":[],
  "known_limits":[],
  "anomalies":[],
  "falsification_conditions":[]
}
```

The exact corpus data belongs in the manifest files; the schema is normative here.

## 17.6 Statement field

`statement` is human-readable documentation.

It MUST NOT be parsed as mathematics.

There is no source-equation parser.

The machine-checkable mathematical representation is `canonical_expr`.

## 17.7 Canonical manifest expression validation

Manifest loading MUST:

1. decode JSON into typed Go structs;
2. decode `canonical_expr` into the closed `core.Expr` node set;
3. validate required field presence and enum values;
4. validate deterministic canonical structure.

A pure validator equivalent to:

```text
ValidateManifestBytes(data []byte) (Manifest, error)
```

MUST exist in `core/corpus.go` or the equivalent exact file from the repository tree.

The validator performs no constructor calls and no external I/O beyond reading its input bytes.

## 17.8 Constructor cross-check

The manifest test MUST contain a static test-only mapping:

```go
map[string]func() core.Object
```

The mapping MUST resolve every manifest `constructor` identifier.

No reflection-based constructor registry is used.

The test MUST:

1. parse manifest bytes;
2. resolve `constructor` using the static map;
3. call the constructor;
4. compare manifest `canonical_expr` with `object.Expr()` structurally and by hash;
5. compare manifest `dimension` with `object.Dimension()`;
6. compare manifest `kind` with `object.Kind()`;
7. compare manifest `provenance_status` with `object.Provenance().Status`;
8. compare manifest assumptions with object assumptions;
9. compare source/framework metadata where specified.

This cross-check verifies that the machine-readable corpus agrees with executable constructors.

## 17.9 Corpus status ownership

The manifest files MUST contain explicit corpus status values.

The implementation MUST load those values rather than derive them.

Tests MUST verify that loading a manifest preserves the declared status.

No implementation code may rewrite it.

---

# 18. Mechanics corpus

Mechanics is the initial classical framework corpus.

## 18.0 Exact mechanics constructor surface

The mechanics package MUST expose these exact constructor signatures:

```go
func NewMass() Mass
func NewTime() Time
func NewPosition() Position
func NewVelocity() Velocity
func NewAcceleration() Acceleration
func NewForce() Force
func NewMomentum() Momentum
func NewEnergy() Energy
func NewKineticEnergy(m Mass, v Velocity) KineticEnergy

func NewtonSecondLaw() core.Object
func MomentumRelation() core.Object
func KineticEnergyRelation() core.Object
```

All zero-argument constructors return the fixed canonical symbol/artifact described by the manifest.

## 18.1 Populated primitives

The package MUST provide thin nominal wrappers for:

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
```

All wrappers expose:

```text
CoreObject() core.Object
```

The canonical symbols are:

```text
m
 t
x
v
a
F
p
E
K
```

The actual constructor naming MUST be:

```text
NewMass
NewTime
NewPosition
NewVelocity
NewAcceleration
NewForce
NewMomentum
NewEnergy
NewKineticEnergy
```

`NewKineticEnergy()` is not a zero-argument corpus manifest constructor because its expression requires Mass and Velocity values; the manifest contains the relation constructor instead.

## 18.2 Mechanics relations

The package MUST provide zero-argument fixed relation constructors:

```text
NewtonSecondLaw()
MomentumRelation()
KineticEnergyRelation()
```

Canonical relation expressions:

```text
F = m*a
p = m*v
K = 1/2*m*v^2
```

`KineticEnergyRelation()` is the corpus artifact.

A parameterized `NewKineticEnergy(m,v)` constructor MUST produce the expression:

```text
1/2 * m.Expr() * Pow(v.Expr(),2)
```

with `Kind = KineticEnergy` and `Dimension = Energy`.

## 18.3 Mechanics framework metadata

Framework ID:

```text
classical_mechanics
```

Framework name:

```text
Classical Mechanics
```

Corpus status:

```text
ESTABLISHED
```

The manifest MUST record the classical/nonrelativistic scope and a framework limitation/anomaly record describing the regime boundary.

The limitation is a scope record, not an automated claim that classical mechanics is "false".

---

# 19. Special-relativity corpus

Only special relativity is included.

## 19.0 Exact relativity constructor surface

The relativity package MUST expose these exact constructor signatures:

```go
func NewSpacetime() Spacetime
func NewMinkowskiMetric() MinkowskiMetric
func NewRestMass() RestMass
func NewEnergy() Energy
func NewThreeMomentum() ThreeMomentum
func NewFourMomentum() FourMomentum
func NewSpeedOfLight() SpeedOfLight
func NewVelocity() Velocity

func LorentzFactor() core.Object
func EnergyMomentumRelation() core.Object
func MassEnergyRelation() core.Object
func ZeroThreeMomentum() core.Object
func ZeroEnergy() core.Object
func ZeroVelocity() core.Object
func RestFrameAssumption() core.Assumption
```

The constructor names are part of the manifest `constructor` field mapping.

## 19.1 Populated primitives

The package MUST provide thin nominal wrappers for:

```text
Spacetime
MinkowskiMetric
RestMass
Energy
ThreeMomentum
FourMomentum
SpeedOfLight
```

It MUST also provide fixed relation/function constructors:

```text
LorentzFactor()
EnergyMomentumRelation()
MassEnergyRelation()
```

The fixed symbols are:

```text
s   spacetime coordinate symbol
E   energy symbol
m   rest mass symbol
p   three-momentum symbol
P   four-momentum symbol
c   speed of light symbol
v   velocity symbol
```

## 19.2 RestFrame

`RestFrame` is not an object.

It is a structured assumption:

```text
Kind: Constraint
Key: rest_frame
Value: ExprValue(
    Relation(eq, Symbol("p"), Rational(0))
)
```

The package MAY expose a fixed helper:

```text
RestFrameAssumption() core.Assumption
```

This is the only role of RestFrame in MVP.

It is not a convention, framework, physical object, or independent kind.

## 19.3 Fixed special-relativity assumptions

The relativity manifest MUST identify these framework assumptions:

```text
Minkowski spacetime
Lorentz symmetry
No gravitational dynamics in package
Special-relativistic regime
```

The Minkowski metric canonical symbol is:

```text
η
```

The manifest also records convention:

```text
metric.signature = -+++
```

This convention attaches to `NewMinkowskiMetric()` and is preserved by objects derived from that metric through pure operations.

## 19.4 Lorentz factor

`LorentzFactor()` returns a valid `core.Object` with:

```text
Kind = Expression
Dimension = dimensionless
Expr = Call("lorentz_factor", [Symbol("v")])
```

The object carries the `c > 0` assumption.

Its fixed body is the expression in §15.9.1.

## 19.5 Energy-momentum relation

`EnergyMomentumRelation()` MUST return:

```text
E^2 = (p*c)^2 + (m*c^2)^2
```

with:

```text
Kind = Relation
Dimension = Energy^2
```

The object MUST carry the explicit `m >= 0` and `c > 0` assumptions through its constituent domain objects.

## 19.6 Mass-energy relation

`MassEnergyRelation()` is a reusable corpus artifact:

```text
E = m*c^2
```

It is retained in the corpus because established physics equations are valid reusable corpus artifacts even when the framework is not treated as an ultimate ontology.

Its manifest provenance MUST be `DERIVED` and its metadata MUST state:

```text
derivable_from = [EnergyMomentumRelation, RestFrame]
```

The constructor itself is a fixed corpus constructor; the vertical derivation test MUST NOT call it to obtain the final result.

## 19.7 Zero constructors

The package MUST expose:

```text
ZeroThreeMomentum()
ZeroEnergy()
ZeroVelocity()
```

Exact values:

```text
ZeroThreeMomentum:
    Kind = ThreeMomentum
    Dimension = Momentum
    Expr = Rational(0)
    Provenance = DEFINED
    Assumptions includes RestFrameAssumption()

ZeroEnergy:
    Kind = Energy
    Dimension = Energy
    Expr = Rational(0)
    Provenance = DEFINED

ZeroVelocity:
    Kind = Velocity
    Dimension = Velocity
    Expr = Rational(0)
    Provenance = DEFINED
```

`ZeroThreeMomentum()` carries the explicit rest-frame constraint so `Substitute` cannot silently drop the condition `p = 0`.

These fixed constructors are required to make the canonical derivations implementable without a generic object factory.

## 19.8 Symbol isolation

The mechanics and relativity packages may both use canonical symbol names such as `v`, but MVP derivations treat symbols package-locally and do not mix domain packages in one derivation. No implicit cross-package symbol namespace or automatic renaming mechanism exists.

## 19.9 Relativity limitation/anomaly

The manifest MUST contain at least one structured scope limitation/anomaly record, minimum example:

```text
id: no_gravity
framework: special_relativity
description: Package covers flat-spacetime special relativity and does not implement gravitational dynamics.
status: scope_limit
research_relevance: framework boundary for candidate theory bridges.
```

No runtime inference from this record is required.

---

# 20. Canonical E = mc² vertical derivation

This is the primary architecture test.

It is not a claim that the library has independently discovered or empirically established mass-energy equivalence.

It verifies that the formal substrate can mechanically carry the required derivation while keeping assumptions and branch choices explicit.

## 20.1 Required sequence

The test MUST execute actual operations in this order:

```text
construct EnergyMomentumRelation()
↓
construct ZeroThreeMomentum()  # includes RestFrameAssumption()
↓
Substitute(relation, ThreeMomentum(), ZeroThreeMomentum())
↓
Simplify(...)
↓
Solve(..., Energy())
↓
construct Compare(Energy(), ZeroEnergy(), gte)
↓
SelectBranch(..., energy_nonnegative)
```

No hardcoded final result is allowed.

The test MUST NOT call `MassEnergyRelation()` to obtain the answer.

## 20.2 Golden trace

The required golden conceptual trace is:

### Step 1 — input relation

```text
E^2 = (p*c)^2 + (m*c^2)^2
```

### Step 2 — substitute rest-frame zero momentum

`ZeroThreeMomentum()` carries `RestFrameAssumption()`; `Substitute` therefore merges the explicit condition `p = 0` into the derivation metadata.

```text
E^2 = (0*c)^2 + (m*c^2)^2
```

### Step 3 — mechanical simplify

```text
E^2 = (m*c^2)^2
```

### Step 4 — solve exact accepted quadratic form

```text
BranchSet(
  E,
  [
    Sqrt((m*c^2)^2),
    Neg(Sqrt((m*c^2)^2))
  ]
)
```

### Step 5 — explicit positive-energy constraint

```text
E >= 0
```

constructed by:

```text
Compare(Energy(), ZeroEnergy(), gte)
```

### Step 6 — select positive branch and simplify

Using the explicit `m >= 0` assumption and the structural fact that `c^2` is nonnegative:

```text
Sqrt((m*c^2)^2) → m*c^2
```

### Final selected solution

```text
m*c^2
```

The derivation MAY then use `Compare(Energy(), selectedSolution, eq)` to construct an explicit relation artifact if a final equation object is required for the candidate handoff.

That comparison remains `DERIVED`, not `IDENTIFIED`.

## 20.3 Why the sign step is valid in the MVP

The simplifier is not allowed to assume `sqrt(x^2)=x` universally.

For the specific expression:

```text
x = m*c^2
```

MVP's bounded sign checker establishes:

```text
m >= 0
c^2 >= 0
m*c^2 >= 0
```

without theorem proving beyond the explicit structural rules in §9.8.

---

# 21. Calculus micro-test

The derivative test MUST be nontrivial and MUST exercise the Expression-vs-named-kind comparison path.

Construct:

```text
m := mechanics.NewMass()
v := mechanics.NewVelocity()
K := mechanics.NewKineticEnergy(m, v)
```

Then:

```text
Differentiate(K.CoreObject(), v.CoreObject())
```

Expected canonical result:

```text
m*v
```

The result kind is `Expression`.

Then compare with named `Momentum`:

```text
Compare(result, mechanics.NewMomentum().CoreObject(), eq)
```

This MUST be accepted under MRC-003 because one side is `Expression`, the other is named `Momentum`, and dimensions match.

The derivative implementation therefore exercises:

```text
product rule
power rule
exact rational simplification
MRC-003 Expression-vs-named comparison
```

The old trivial test `d(v)/dt` is not used.

---

# 22. Relativistic limit micro-test

Execute:

```text
Limit(LorentzFactor(), Velocity(), ZeroVelocity())
```

Expected canonical result:

```text
1
```

The result MUST come from expansion of the fixed Lorentz-factor body and direct substitution followed by simplification.

A test that merely matches the function ID and returns `1` without traversing the fixed body is invalid.

This test demonstrates framework self-consistency/reduction support, not empirical validation.

The meaningful Newtonian kinetic-energy low-velocity series reduction remains deferred until series expansion exists in v0.5+.

---

# 23. Identify firewall micro-test

The test MUST establish both sides of the firewall.

First:

```text
Simplify(x)
```

MUST NOT produce `IDENTIFIED` provenance and MUST NOT append a ledger identification event.

Second:

```text
Session.Identify(a,b,justification)
```

MUST:

- require explicit operands;
- produce `IDENTIFIED` for non-hypothesis inputs;
- record the justification;
- record a ledger `Identification` step;
- remain visible after sealing.

If either operand is `HYPOTHESIS`, the result remains `HYPOTHESIS` and the identification attempt is preserved for review.

---

# 24. Hypothesis candidate concept

The `hypothesis` package is AI-callable provisional construction only.

## 24.1 Candidate concept definition

Two different concepts are deliberately distinguished:

1. **Candidate concept** — a `core.Object` with `Provenance.Status == HYPOTHESIS`.
2. **ResearchCandidate** — a sealed derivation artifact created only by `session.Session.Seal`.

## 24.2 Candidate concept constructor

The exact constructor is conceptually:

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

Requirements:

- `id` non-empty;
- `kind` explicit and immutable after construction;
- `dimension` explicit and immutable after construction;
- `expr` valid;
- assumptions/conventions deterministic;
- provenance forced to `HYPOTHESIS`;
- corpus status forced to `NONE`;
- no caller-supplied provenance status;
- no promotion API.

**REQ-024-01**: Candidate concept construction MUST force `HYPOTHESIS` provenance.

**REQ-024-02**: Candidate concept construction MUST require explicit `Kind` and `Dimension`.

## 24.3 Candidate containment test

The test MUST:

1. create a candidate concept;
2. derive with it;
3. verify downstream result is `HYPOTHESIS`;
4. verify `Simplify`, `Solve`, `Limit`, `Compare`, `SelectBranch`, and `Session.Identify` do not upgrade it;
5. verify no promotion API exists;
6. verify a hypothesis-dependent ResearchCandidate cannot validate as trusted corpus authority.

The last assertion is a validation/integrity assertion, not a scientific-truth judgment.

## 24.4 No promotion API

The MVP source tree MUST NOT expose exported functions/methods named or semantically equivalent to:

```text
Promote
Trust
ApproveHypothesis
PromoteToEstablished
SetCorpusStatus
```

The candidate test SHOULD use a standard-library AST/export-surface inspection to verify the absence of such exported declarations rather than attempting to compile intentionally invalid Go source inside a test.

---

# 25. Falsifiability structures

The candidate system MUST represent research-test structure without deciding whether the candidate is true or false.

## 25.1 Prediction

Fields:

```text
ID
Observable
Relation
Assumptions
```

`Relation` is a canonical `core.Expr` relation node.

## 25.2 FalsificationCondition

Fields:

```text
ID
TargetClaim
ContradictingCondition
Regime
```

`ContradictingCondition` is a canonical expression, not free-form executable text.

## 25.3 RecoveryClaim

Fields:

```text
ID
Description
FromFramework
Condition
```

`Condition` is structured canonical expression/assumption data.

## 25.4 AnomalyReference

Fields:

```text
ID
Framework
Description
```

Anomaly references point to manifest/anomaly records but do not perform automatic inference.

## 25.5 Required falsifiability test

At least one provisional candidate MUST be constructed with:

- one prediction;
- one falsification condition;
- one anomaly reference.

The sealed artifact MUST preserve them.

No truth status is generated by this test.

---

# 26. ResearchCandidate artifact

`ResearchCandidate` is a sealed, reviewable research artifact.

## 26.0 Exact ResearchCandidate support types

The `session` package MUST define the following canonical data types:

```go
type FrameworkDependency struct {
    FrameworkID       string
    ManifestHash      string
    AssumptionHashes  []string
}

type Prediction struct {
    ID          string
    Observable  string
    Relation    core.Expr
    Assumptions core.AssumptionSet
}

type FalsificationCondition struct {
    ID                  string
    TargetClaim         string
    ContradictingCondition core.Expr
    Regime              string
}

type RecoveryClaim struct {
    ID            string
    Description   string
    FromFramework string
    Condition     core.Expr
}

type AnomalyReference struct {
    ID          string
    Framework   string
    Description string
}

type Challenge struct {
    StepID      string
    Category    string
    Severity    string
    Description string
}

type Review struct {
    DerivationID string
    Challenges   []Challenge
    ReviewerNotes string
}
```

All these structures MUST have deterministic canonical encoders and MUST NOT use dynamic maps.

## 26.1 Authority

Only a successfully validated `session.Session` may mint a `ResearchCandidate` in MVP.

A decoded external artifact is not considered minted by the decoder; it is merely loaded for validation.

No API promotes a candidate into corpus authority.

## 26.2 Required fields

Fields, in canonical order:

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

## 26.3 Hypothesis field

`Hypothesis` MUST be a candidate concept object with `HYPOTHESIS` provenance.

The candidate artifact MUST NOT replace it with a trusted status.

## 26.4 Premises

`Premises` are the canonical object inputs that materially seed the derivation.

They are retained as immutable data.

## 26.5 Framework dependencies

Each framework dependency contains:

```text
FrameworkID
ManifestHash
AssumptionHashes[]
```

This makes framework scope explicit and deterministic.

## 26.6 Derivation

`Derivation` contains the sealed canonical session ledger.

Its final hash MUST equal `LedgerHash`.

## 26.7 Review history

`ReviewHistory` contains zero or more `Review` artifacts.

The MVP does not orchestrate reviewer agents.

## 26.8 Validation pipeline

`ResearchCandidate.Validate()` MUST perform:

```text
schema validation
→ canonicalization
→ ledger hash validation
→ derivation replay
→ provenance/candidate validation
→ framework-reference integrity validation
→ falsifiability/anomaly field validation
```

Within `session`, framework-reference integrity means only deterministic structural validity: non-empty framework IDs, valid manifest-hash format, non-empty assumption-hash strings where required, and well-formed anomaly IDs. The `session` package MUST NOT import `mechanics` or `relativity`. Tests that need to verify that a referenced anomaly actually exists in a manifest perform that cross-check externally.

Successful validation means:

```text
artifact is internally consistent
```

It does not mean:

```text
hypothesis is true
```

## 26.9 Candidate containment at artifact level

A `ResearchCandidate` created by `Session.Seal()` is a sealed value. Its public accessors may expose validated `core.Object` values because those values have already passed session validation.

External JSON MUST NOT be decoded directly into a trusted `ResearchCandidate` value. Public deserialization therefore uses an explicitly unverified wrapper.

## 26.10 ResearchCandidate public surface

`ResearchCandidate` MUST have unexported storage and MUST NOT expose an arbitrary public constructor.

Trusted sealed `ResearchCandidate` construction is owned exclusively by the `session` package. `Session.Seal()` is the canonical production minting action. `UnverifiedResearchCandidate.Validate()` is an external-loading entrypoint and MUST delegate to the same unexported session sealing/validation constructor; it MUST NOT contain an independent trusted-construction path or bypass session validation.

It MUST expose:

```go
func (c ResearchCandidate) ID() string
func (c ResearchCandidate) Hypothesis() core.Object
func (c ResearchCandidate) Premises() []core.Object
func (c ResearchCandidate) Assumptions() core.AssumptionSet
func (c ResearchCandidate) Derivation() Ledger
func (c ResearchCandidate) FrameworkDependencies() []FrameworkDependency
func (c ResearchCandidate) Predictions() []Prediction
func (c ResearchCandidate) FalsificationConditions() []FalsificationCondition
func (c ResearchCandidate) RecoveryClaims() []RecoveryClaim
func (c ResearchCandidate) AnomalyReferences() []AnomalyReference
func (c ResearchCandidate) ReviewHistory() []Review
func (c ResearchCandidate) MRCVersion() string
func (c ResearchCandidate) LedgerHash() string
func (c ResearchCandidate) Validate() error
func (c ResearchCandidate) CanonicalJSON() ([]byte, error)
```

Collection accessors MUST return copies.

The external loader MUST have the exact unverified surface:

```go
type UnverifiedResearchCandidate struct { /* unexported storage */ }

func ParseResearchCandidateJSON(data []byte) (UnverifiedResearchCandidate, error)
func (u UnverifiedResearchCandidate) CanonicalJSON() ([]byte, error)
func (u UnverifiedResearchCandidate) Validate() (ResearchCandidate, error)
```

`ParseResearchCandidateJSON` MUST NOT return `ResearchCandidate` and MUST NOT expose `core.Object` accessors. It only validates the outer canonical artifact schema and stores canonical bytes for later validation.

`UnverifiedResearchCandidate.Validate()` is the only public path from externally supplied candidate JSON to a trusted `ResearchCandidate`, and it MUST execute the full validation pipeline in §26.8 before returning success.

A candidate artifact MUST fail validation with `CandidateContainmentError` if:

- `Hypothesis` is not `HYPOTHESIS`;
- a hypothesis-dependent derivation output is marked trusted;
- a promotion-like state is claimed internally.

Candidate artifact validation MUST NOT treat self-recomputed hashes as authenticity proof. The validator establishes internal consistency only.
# 27. Review artifact

The review artifact is data only. The `Challenge` and `Review` structures are defined in `core/corpus.go`; no `review/` package exists in MVP.

## 27.1 Challenge

Fields:

```text
StepID
Category
Severity
Description
```

## 27.2 Review

Fields:

```text
DerivationID
Challenges
ReviewerNotes
```

## 27.3 Categories

Exact MVP categories:

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

## 27.4 Orchestration boundary

The library defines the artifact schemas but does not implement:

- reviewer selection;
- agent-to-agent orchestration;
- review queues;
- automatic challenge generation;
- automatic candidate correction.

The intended external loop is:

```text
ResearchCandidate JSON
      ↓
external reviewer
      ↓
Review JSON
      ↓
human/formulator
```

---

# 28. Manifest data requirements

## 28.1 Mechanics manifest required items

The manifest MUST contain machine-readable entries corresponding to at least:

```text
Mass
Time
Position
Velocity
Acceleration
Force
Momentum
Energy
NewtonSecondLaw
MomentumRelation
KineticEnergyRelation
```

`KineticEnergy` as a parameterized constructor need not appear as a zero-argument manifest item because the canonical corpus item is represented by `KineticEnergyRelation`.

## 28.2 Relativity manifest required items

The manifest MUST contain machine-readable entries corresponding to at least:

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
MassEnergyRelation
```

`RestFrame` is represented as an assumption, not a manifest physical object.

## 28.3 Limitation/anomaly requirements

Each populated framework manifest MUST include at least one structured limitation/anomaly record.

The candidate test MUST demonstrate referencing at least one such record.

---

# 29. MRC fallibility and revision boundary

The library records:

```text
MRCVersion = mrc-v0.4
```

Each rule ID is stable.

There is no runtime MRC override.

There is no runtime exception registry.

Future MRC revision conceptually follows:

```text
current rule
    ↓
human-validated exception proposal
    ↓
evidence + falsification/review
    ↓
new scoped/versioned MRC rule
    ↓
regression suite
```

The MVP implements only:

- version recording;
- stable rule identifiers;
- deterministic rule enforcement.

The MRC is a fallible methodological model, not a metaphysical declaration that the library's rules can never be revised.

---

# 30. Mandatory error taxonomy

The exact typed errors are:

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

All MUST implement Go's `error` interface.

All MUST be inspectable with `errors.As`.

Normal MRC violations MUST NOT panic.

The errors SHOULD contain deterministic diagnostic fields such as operation ID, offending kinds/dimensions, and stable rule ID; they MUST NOT contain nondeterministic values.

---

# 31. Determinism requirements

Two identical derivation runs MUST produce byte-identical canonical artifacts.

The following MUST be deterministic:

```text
expression canonical JSON
expression hash
object canonical JSON
object hash
assumption canonical JSON/hash
convention canonical JSON/hash
manifest canonical representation
step canonical representation
step hash chain
ledger canonical JSON
ResearchCandidate canonical JSON
ResearchCandidate LedgerHash
```

No artifact may depend on:

```text
timestamps
random UUIDs
pointer addresses
map iteration order
environment-dependent path ordering
locale-sensitive ordering
process IDs
hostnames
```

## 31.1 Determinism test

Run the same canonical derivation twice in the same test process and compare all canonical outputs.

The test SHOULD additionally run the derivation through a fresh Session instance to verify that no hidden mutable global state exists.

---

# 32. Mandatory negative tests

The implementation MUST contain tests for all of the following:

**REQ-032-01** dimension mismatch;

**REQ-032-02** physical-kind mismatch with equal dimensions;

**REQ-032-03** assumption conflict;

**REQ-032-04** convention conflict;

**REQ-032-05** invalid zero `core.Object`;

**REQ-032-06** absence of public generic object factory;

**REQ-032-07** `Simplify` cannot create `IDENTIFIED` provenance;

**REQ-032-08** missing Identify justification;

**REQ-032-09** candidate contamination bypass;

**REQ-032-10** absence of promotion API;

**REQ-032-10a** externally supplied `ResearchCandidate` JSON loads only as `UnverifiedResearchCandidate`; no public parse path returns trusted `ResearchCandidate` or `core.Object` values;

**REQ-032-11** corrupted ledger hash chain;

**REQ-032-12** modified intermediate canonical output;

**REQ-032-13** modified output with recomputed local hash but stale step hash;

**REQ-032-14** modified output with recomputed full chain but replay divergence;

**REQ-032-15** inconsistent manifest;

**REQ-032-16** exact-rational canonical JSON round-trip;

**REQ-032-17** nondeterministic canonicalization attempt;

**REQ-032-18** unsupported general solver form;

**REQ-032-19** unsupported derivative node/function form;

**REQ-032-20** unsupported general limit form;

**REQ-032-21** invalid session state transition;

**REQ-032-22** post-seal mutation attempt.

---

# 33. Explicit acceptance tests A–S

The MVP is complete only when every test A–S passes.

## A — Typed mechanics

Construct valid mechanics domain values only through fixed constructors and expose their `core.Object` values.

## B — Classical relation

Construct and validate:

```text
F = m*a
```

and verify manifest/constructor agreement.

## C — Dimension rejection

Attempt an invalid dimensional Add/Compare and receive `DimensionMismatchError`.

## D — Category rejection

Attempt to combine distinct named kinds with equal dimensions and receive `CategoryMismatchError`.

## E — Assumption conflict

Construct two valid assumption-bearing objects with the same assumption key but conflicting structured values and verify `AssumptionConflictError`.

## F — Convention conflict

Create conflicting convention sets and verify `ConventionConflictError`.

## G — Differentiation

Differentiate `1/2*m*v^2` with respect to `v` and obtain canonical `m*v`.

Then compare to named `Momentum` using `Compare(..., eq)` and verify success.

## H — Relativistic relation

Construct `EnergyMomentumRelation()` and verify canonical expression, dimensions, assumptions, and manifest match.

## I — Mass-energy derivation

Execute the full §20 sequence using actual operations and obtain final selected solution:

```text
m*c^2
```

No hardcoding of the output is permitted.

## J — Relativistic limit

Execute the fixed-body Lorentz-factor limit and obtain exact `1`.

## K — Identify firewall

Verify `Simplify` cannot mint `IDENTIFIED`, while `Session.Identify` can create and record an identification with explicit justification.

## L — Candidate containment

Verify hypothesis provenance remains contagious through mechanical operations and cannot be promoted.

## M — Falsifiability

Verify sealed candidate contains Prediction and FalsificationCondition artifacts.

## N — Anomaly

Verify the candidate references a manifest/framework anomaly or scope limitation.

## O — Manifest

Verify both manifests load through `go:embed`, decode into strict typed structs, validate, and cross-check every item against the static constructor mapping.

## P — Exact rational round-trip

Verify exact rational and expression canonical JSON round-trips preserve canonical bytes and hash. Also verify exact `Pow(Rational(0),2) → 0` and `Pow(Rational(1),-1) → 1` simplification required by the canonical derivations.

## Q — Ledger tamper/replay

Verify all tampering forms in §32 are detected by hash validation and/or replay.

## R — Determinism

Verify repeated equivalent derivations produce byte-identical canonical artifacts and hashes.

## S — Candidate handoff

Verify `Session.Seal()` after full validation produces a sealed `ResearchCandidate` whose `LedgerHash` matches the sealed derivation.

---

# 34. Framework evaluation anatomy

Each populated framework MUST be evaluated through three complementary mechanisms.

## 34.1 Deterministic derivation test

Mechanics:

```text
F = m*a
```

Relativity:

```text
E = m*c^2
```

through the exact paths above.

## 34.2 Executable self-consistency/reduction

Special relativity:

```text
lim(v→0) LorentzFactor(v) = 1
```

This is the MVP executable reduction.

Newtonian kinetic-energy recovery through an explicit series expansion is deferred.

## 34.3 Limitation/anomaly record

Each framework manifest MUST contain at least one limitation/anomaly record.

These records are machine-readable context for AI/human research.

The library does not infer what should be done about an anomaly.

---

# 35. Paper translation documentation

The single document:

```text
docs/paper-translation.md
```

is auxiliary guidance.

It MUST contain sections equivalent to:

```text
Common notation
Framework mapping
Ambiguity resolution
```

It is not canonical machine-readable physics metadata.

The manifests are canonical machine-readable semantic metadata.

The translation document may explain how conventional paper notation maps to the library's typed representation, but it MUST NOT be required for runtime execution.

---

# 36. `physvet` boundary

No `physvet` implementation is part of MVP.

A future `physvet` implementation MUST consume the same public/core semantic contracts and metadata.

It MUST NOT independently redefine:

```text
dimensions
physical kinds
assumptions
conventions
provenance
operation contracts
corpus status
```

This is a v0.5+ architectural constraint.

---

# 37. Explicit v0.5+ deferrals

The following are outside MVP:

1. `physvet` implementation.
2. General symbolic series expansion.
3. General truncation-order propagation.
4. Full 1905 Einstein low-velocity series derivation.
5. Symbolic integration.
6. Broader symbolic calculus.
7. General equation solving.
8. General power/branch algebra.
9. Richer tensor/index machinery.
10. General relativity.
11. Quantum mechanics.
12. Quantum field theory.
13. Statistical mechanics.
14. Electromagnetism.
15. Numerical execution.
16. Empirical-data adapters.
17. Reviewer-AI orchestration.
18. Automatic corpus governance.
19. Automatic hypothesis promotion.
20. MRC exception/override workflow.
21. EBP 2.1 integration.
22. General theory-management platform.

No deferred feature may leak into MVP through placeholder packages or speculative abstractions.

---

# 38. Scope guardrail

The repository file baseline is the exact tree in §3.

The total file count is:

```text
39 including go.mod
38 excluding go.mod
```

The MVP file budget is:

```text
<= 40 repository files total
```

The count includes `go.mod` and excludes only `go.sum`, generated artifacts, and VCS metadata.

Exact §3 composition: 24 Go source files, 10 Go test files, 2 JSON manifest files, and 3 configuration/documentation files (`go.mod`, `README.md`, `docs/paper-translation.md`) = 39 files total.

Because the normative tree already fits within the cap, the implementation plan MUST NOT add files outside §3.

Convenience-only files, placeholder abstractions, compatibility shims, future scaffolding, generated registries, or empty packages are prohibited.

---

# 39. Normative requirement and test coverage matrix

The implementation plan MUST map every normative `MUST`/`MUST NOT` clause in this specification to an implementation location and at least one test. Where a normative clause has no explicit `REQ-*` label, the implementation plan MUST assign a deterministic coverage ID such as `REQ-§15.8-MUST-01`; no normative clause may be omitted from the matrix.

The minimum mandatory matrix is reproduced below so the coding/plan agent has an explicit target.

| Requirement | Implementation location | Primary test |
|---|---|---|
| REQ-000-01 | go.mod | package/build gate |
| REQ-000-02 | go.mod/repository | `go test ./...` |
| REQ-000-03 | core/session architecture | candidate handoff |
| REQ-000-04 | README + core/internal boundary | constructor authority test |
| REQ-001-01 | core/ops | A,C,D |
| REQ-001-02 | mechanics/relativity manifests | H,I,J,O |
| REQ-001-03 | core/corpus + manifests | O |
| REQ-001-04 | core metadata + session | A–S integration |
| REQ-001-05 | banned API/source-surface tests | L,S + API-surface check |
| REQ-001-06 | hypothesis/session | L,M,N,S |
| REQ-002-01..26 | repository tree | source-tree/package audit |
| REQ-003-01 | repository tree | scope audit |
| REQ-003-02 | *_test.go placement | repository audit |
| REQ-003-03 | external tests | MRC-001 test |
| REQ-004-01 | package graph | `go test ./...` |
| REQ-004-02 | core imports | source audit |
| REQ-004-03 | internal/kernel | external-package authority test |
| REQ-005-01..12 | core/object + domains + kernel | A,C,D,L |
| REQ-006-01..03 | core/object + ops | A,D,G |
| REQ-007-01..03 | core/dimension | C,P,R |
| REQ-008-01..06 | core/expr + canonical | P,R,G,I,J |
| REQ-009-01..06 | core/canonical + ops/simplify | G,I,J,P,R |
| REQ-010-01..07 | core/canonical/corpus/session | O,P,Q,R,S |
| REQ-011-01..02 | core/assumption | E,I,J |
| REQ-012-01..02 | core/convention | F |
| REQ-013-01..05 | core/provenance + session | K,L,S |
| MRC-001 | internal/kernel + core + domains | C/A constructor authority |
| MRC-002 | ops | C |
| MRC-003 | ops | D,G |
| MRC-004 | ops/core assumptions | E |
| MRC-005 | ops/core conventions | F |
| MRC-006 | session | K |
| MRC-007 | session | Q,S |
| MRC-008 | ops/session/hypothesis | L,S |
| REQ-015-01 | ops | R + source audit |
| REQ-016-01 | session/session.go | session state tests |
| REQ-024-01..02 | hypothesis | L |
| REQ-032-01..22 | *_test.go | negative test suite |

The implementation plan MUST expand grouped rows such as `REQ-002-01..26` into individual requirement rows before finalizing the plan.

`APPROXIMATED` is exempt from required status-production coverage in MVP under §13.1.

---

# 40. Final architectural invariants

The MVP is successful only when the following deterministic loop works:

```text
AI reads manifest/corpus
        ↓
typed domain construction
        ↓
immutable core.Object
        ↓
explicit assumptions + conventions
        ↓
pure symbolic operations
        ↓
MRC checks
        ↓
session draft/commit
        ↓
explicit Identify when required
        ↓
canonical replayable ledger
        ↓
provisional hypothesis
        ↓
prediction + falsification condition
        ↓
anomaly/scope context
        ↓
sealed ResearchCandidate
        ↓
human review
```

The system stops there.

It does not decide what nature believes.

The governing design principle is:

> **Open representation, closed authority; explicit assumptions, explicit insight; deterministic mechanics, fallible MRC; machine-readable physics, human scientific judgment.**

---

# 41. Final handoff condition — no architectural blockers

This specification intentionally resolves the architectural issues identified in the review cycle.

The resolved choices are normative:

1. **Go module identity:** `github.com/PithomLabs/phys`.
2. **Constructor authority:** actual object implementation in `internal/kernel`; public `core.Object` aliases the immutable kernel type; no generic public factory.
3. **Ledger placement:** `session` owns ledger/replay/seal so `core` does not import `ops`.
4. **Replay substrate:** committed steps retain canonical inputs, canonical outputs, and canonical operation parameters; hashes are verification/integrity fields.
5. **Identify ownership:** only `Session.Identify`; no ambient session state; no `ops.Identify`.
6. **Structured assumptions:** mathematical assumptions use canonical expression values rather than relationship strings.
7. **Branch representation:** `BranchSet` is a bounded expression node.
8. **Pow result:** `Kind=Expression`.
9. **Lorentz factor:** one fixed internal expression body, no general function registry.
10. **Sqrt sign handling:** exact bounded entailment rules using explicit `m>=0` and `c>0` assumptions.
11. **Zero values:** fixed dimensioned zero constructors exist where canonical derivations require them.
12. **Differentiation test:** nontrivial derivative of `1/2*m*v^2` with comparison to named Momentum.
13. **Solver:** exact pattern matching only; no general solver.
14. **Manifest mapping:** every item declares `constructor`; tests resolve through a static function map.
15. **File budget:** exact repository tree is below the 40-file guardrail.
16. **Provenance:** exact contamination/propagation law is fixed.
17. **Candidate terminology:** candidate concept and ResearchCandidate are separate artifacts.
18. **MRC fallibility:** versioned and non-overridable in MVP.
19. **External candidate loading:** `ParseResearchCandidateJSON` returns only `UnverifiedResearchCandidate`; trusted `ResearchCandidate` is obtainable only after validation, and unverified artifacts expose no trusted `core.Object` accessors.
20. **Expression normalization boundary:** constructors perform only structural normalization; semantic rewrite rules, including exact rational power identities, belong to `Simplify`.
21. **Session seal gate:** non-hypothesis derivations may conclude and validate but cannot seal a `ResearchCandidate`; sealing requires valid `HYPOTHESIS` metadata.

**Handoff invariant:** a coding agent MUST NOT need to choose among alternative architectures to implement this MVP. Any implementation that requires reopening one of these decisions is non-conforming.


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

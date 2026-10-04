# DESIGNING A PHYSICS COMPILER IN GO
## From Program A MVP to a Full Physics Language

**Grounding:** The file proposes a two-phase evolution. Phase A (Program A MVP) discovers the language by using it. Phase B (Physics Compiler) formalizes the language from observed patterns. This design describes Phase B, using Phase A as the semantic backend.

---

## Part I: The Architecture

### The Pipeline

```
Physics Source (.phys)
        │
        ▼
    ┌─────────┐
    │  Lexer  │
    └────┬────┘
         │ tokens
         ▼
    ┌─────────┐
    │ Parser  │
    └────┬────┘
         │ AST
         ▼
    ┌──────────┐
    │ Type     │
    │ Checker  │
    └────┬─────┘
         │ typed AST
         ▼
    ┌──────────┐
    │ Dimension│
    │ Checker  │
    └────┬─────┘
         │ dimensioned AST
         ▼
    ┌──────────┐
    │ Physics  │
    │ Lowering │
    └────┬─────┘
         │ Physics IR
         ▼
    ┌──────────┐
    │ Program A│  (semantic backend)
    │ Semantic │
    │ Passes   │
    └────┬─────┘
         │
    ┌────┴────┬──────────┬──────────┐
    ▼         ▼          ▼          ▼
Program B   Lean      CAS      Numerical
(executor) (proof)   (symbolic) (simulation)
```

### The Division of Labor

| Layer | Job | Fails At |
|---|---|---|
| **Lexer** | Text → tokens | Illegal characters, malformed strings |
| **Parser** | Tokens → AST | Syntax errors |
| **Type Checker** | AST → typed AST | Type mismatches (Scalar + Vector) |
| **Dimension Checker** | Typed AST → dimensioned AST | Dimensional mismatches (Length + Time) |
| **Lowering** | Dimensioned AST → Physics IR | Structural gaps (undefined references) |
| **Program A** | Physics IR → semantic verdict | Semantic errors (circularity, category errors, mode violations) |
| **Backends** | Physics IR → results | Computational errors |

**The compiler does not compute.** It lowers to IR. Program B computes.

---

## Part II: The Physics Language

### Design Principles

1. **Physics semantics, not just math syntax.** A `Group` is not a `ContinuousSymmetryOfAction`. The language distinguishes.
2. **Types carry physical role.** `Metric` ≠ `Tensor[2,0]`. `State` ≠ `Vector`.
3. **Modes are first-class.** `derive`, `represent`, `identify`, `interpret`, `construct`, `emerge`, `recover` are keywords.
4. **Regimes are first-class.** `regime microscopic`, `regime IR`, `regime continuum` are declarations.
5. **Dimensional analysis is built in.** Length, Time, Mass, Action, Charge are primitive dimensions.
6. **Assumptions are explicit.** `assume` is a keyword. Imports are declared.
7. **Silences are declared.** `silent` is a keyword. The Born rule, gauge group, particle content are declared as silences.

### The Grammar (Surface Level)

```
program     := declaration*
declaration := stage | field | action | symmetry | derive
             | assume | silent | regime | observable
             | construct | emerge | recover | identify
             | interpret | represent

stage       := "stage" IDENT "{" stage_body "}"
stage_body  := (dimension | metric | topology | causality)*

field       := "field" IDENT ":" TYPE

action      := "action" IDENT "{" action_body "}"

symmetry    := "symmetry" IDENT ":" "ContinuousSymmetry(" IDENT ")"

derive      := "derive" IDENT "{" derive_body "}"

assume      := "assume" IDENT ":" TYPE

silent      := "silent" IDENT

regime      := "regime" IDENT "{" regime_body "}"

observable  := "observable" IDENT ":" TYPE
```

### Example: Problem of Time

```physics
// Phase 1: Causal substrate
stage CausalSubstrate {
    primitive: Event
    structure: PartialOrder(≺)
    property: LocallyFinite
}

// Phase 2: Discrete volume
assume SprinklingDensity : Scalar
field Cardinality : (Event, Event) → Nat
relation CardinalityApproximatesVolume {
    forall x, y in CausalSubstrate:
        Cardinality(x, y) ≈ SprinklingDensity * Volume(x, y)
}

// Phase 3: Reconstruction
derive ConformalMetric {
    from CausalSubstrate.PartialOrder
    using MalamentTheorem
    yields g_hat : ConformalClass(Metric)
}

derive ConformalFactor {
    from Cardinality
    using FixedPointCondition
    yields Ω : Scalar
}

derive PhysicalMetric {
    from ConformalMetric, ConformalFactor
    yields g : Metric
    where g = Ω² * g_hat
}

// Phase 4: Proper time
derive ProperTime {
    from PhysicalMetric
    yields τ : (Path) → Scalar
    where τ(γ) = ∫_γ √(-g(dx, dx))
}

// Phase 5: Physical time
identify PhysicalTime {
    as ProperTime
    along ClockWorldline
    requires ExplicitBridge
}

// Phase 6: Dynamics
action EinsteinHilbert {
    term: (1/16πG) * ∫ R √(-g) d⁴x
    regime: IR
}

// Phase 7: Measure
silent BornRule
interpret Measure {
    as |Ψ|²
    where Ψ : WaveFunctional
}
```

### What the Language Enforces

| Construct | What It Enforces |
|---|---|
| `stage` | The Stage is declared explicitly. It has primitive, structure, properties. |
| `assume` | Assumptions are declared. They are not hidden. |
| `silent` | Silences are declared. The compiler flags attempts to derive them. |
| `derive` | Derivations have explicit `from`, `using`, `yields`. Dependencies are explicit. |
| `identify` | Identifications require `requires ExplicitBridge`. They cannot be silent. |
| `interpret` | Interpretations require `as`. They cannot be silent. |
| `regime` | Claims have regimes. Silent regime movement is flagged. |
| `construct` | Constructions are distinct from derivations. |
| `emerge` | Emergence requires scaling/limit. |
| `recover` | Recovery is distinct from derivation. |

---

## Part III: The Compiler Pipeline

### Stage 1 — Lexer

**Input:** Source text.
**Output:** Token stream.

```go
// internal/lexer/token.go
package lexer

type TokenKind int

const (
    EOF TokenKind = iota
    IDENT
    KEYWORD
    NUMBER
    OPERATOR
    DELIMITER
    STRING
    COMMENT
)

type Token struct {
    Kind    TokenKind
    Lexeme  string
    Line    int
    Column  int
}

type Lexer struct {
    input    string
    pos      int
    line     int
    column   int
    tokens   []Token
}

func (l *Lexer) Tokenize() ([]Token, error) {
    // Scan input, emit tokens.
    // Keywords: stage, field, action, symmetry, derive,
    //           assume, silent, regime, observable,
    //           construct, emerge, recover, identify,
    //           interpret, represent
    // ...
}
```

**Errors produced:** Illegal characters, unterminated strings.

### Stage 2 — Parser

**Input:** Token stream.
**Output:** AST.

```go
// internal/ast/ast.go
package ast

type Node interface {
    Pos() Position
}

type Program struct {
    Declarations []Declaration
}

type Declaration interface {
    Node
    declNode()
}

type StageDecl struct {
    Name       string
    Primitives []Primitive
    Structure  Structure
    Properties []Property
    Pos        Position
}

type DeriveDecl struct {
    Name       string
    From       []Reference
    Using      []Reference
    Yields     []Binding
    Where      []Equation
    Regime     string
    Pos        Position
}

type IdentifyDecl struct {
    Name       string
    As         Reference
    Along      Reference
    Requires   []Requirement
    Pos        Position
}

type SilentDecl struct {
    Name string
    Pos  Position
}

type AssumeDecl struct {
    Name string
    Type TypeExpr
    Pos  Position
}

// ... etc.
```

**Errors produced:** Syntax errors, unexpected tokens.

### Stage 3 — Type Checker

**Input:** AST.
**Output:** Typed AST.

```go
// internal/types/types.go
package types

type Type interface {
    String() string
    Underlying() Type
    Equals(Type) bool
}

type ScalarType struct{}
type VectorType struct{ Dim int }
type TensorType struct{ Up, Down int }
type OperatorType struct{ Domain, Codomain Type }
type StateType struct{ Space Type }
type ObservableType struct{ Space Type }
type FieldType struct{ Domain, Codomain Type }
type MetricType struct{}
type PartialOrderType struct{ Carrier Type }
type ManifoldType struct{ Dim int }
type TopologyType struct{}
type ActionType struct{}
type ConstraintType struct{}
type GeneratorType struct{}
type InvariantType struct{}
type ParameterType struct{}

// Physical role wrappers
type PhysicalTime struct{}
type CoordinateTime struct{}
type ProperTime struct{}
type RGScale struct{}
type CausalOrder struct{}
type ThermodynamicArrow struct{}
```

**Type checking rules:**

```go
// internal/sema/typecheck.go
package sema

func (c *Checker) checkDerive(d *ast.DeriveDecl) error {
    // 1. Resolve `from` references.
    // 2. Resolve `using` references.
    // 3. Resolve `yields` bindings.
    // 4. Check that yielded types are compatible with the derivation.
    // 5. Check that no yielded type appears in `from` (circularity).
    // ...
}

func (c *Checker) checkIdentify(i *ast.IdentifyDecl) error {
    // 1. Resolve `as` reference.
    // 2. Check that `requires` includes ExplicitBridge.
    // 3. Check that the identified types are compatible.
    // ...
}
```

**Errors produced:** Type mismatches, undefined symbols, incompatible types.

### Stage 4 — Dimension Checker

**Input:** Typed AST.
**Output:** Dimensioned AST.

```go
// internal/dimensions/dimensions.go
package dimensions

type Dimension struct {
    Length      int
    Time        int
    Mass        int
    Charge      int
    Temperature int
    Action      int
    // ...
}

var (
    Dimensionless = Dimension{}
    Length        = Dimension{Length: 1}
    Time          = Dimension{Time: 1}
    Velocity      = Dimension{Length: 1, Time: -1}
    Energy        = Dimension{Mass: 1, Length: 2, Time: -2}
    Action        = Dimension{Mass: 1, Length: 2, Time: -1}
)

func (d Dimension) Mul(other Dimension) Dimension {
    return Dimension{
        Length:      d.Length + other.Length,
        Time:        d.Time + other.Time,
        Mass:        d.Mass + other.Mass,
        // ...
    }
}

func (d Dimension) Equals(other Dimension) bool {
    return d == other
}
```

**Dimension checking rules:**

```go
// internal/sema/dimcheck.go
package sema

func (c *Checker) checkEquation(eq *ast.Equation) error {
    lhs := c.dimensionOf(eq.LHS)
    rhs := c.dimensionOf(eq.RHS)
    if !lhs.Equals(rhs) {
        return errorf("dimensional mismatch: %s vs %s", lhs, rhs)
    }
    return nil
}
```

**Errors produced:** Dimensional mismatches, incompatible units.

### Stage 5 — Physics Lowering

**Input:** Dimensioned AST.
**Output:** Physics IR.

```go
// internal/ir/ir.go
package ir

type Hypothesis struct {
    ID          string
    Goal        string
    Concepts    []Concept
    Claims      []Claim
    Assumptions []Assumption
    Silences    []Silence
    Regimes     []Regime
    Obligations []Obligation
}

type Concept struct {
    Name       string
    Component  Component  // STAGE, DANCE, DANCER, JOINT
    MathType   MathType
    PhysicalRole string
}

type Claim struct {
    ID         string
    Mode       Mode       // DERIVE, IDENTIFY, REPRESENT, ...
    Relation   Relation   // DERIVES, IDENTIFIES, REPRESENTS, ...
    Subject    Reference
    Object     Reference
    Premises   []Reference
    Regime     string
    Obligations []Obligation
}

type Mode int

const (
    ModePostulate Mode = iota
    ModeRepresent
    ModeDerive
    ModeIdentify
    ModeInterpret
    ModeConstruct
    ModeEmerge
    ModeRecover
    ModeReconstruct
    ModeInfer
    ModeConstrain
    ModeCouple
)

type Component int

const (
    ComponentStage Component = iota
    ComponentDance
    ComponentDancer
    ComponentJoint
)
```

**Lowering rules:**

```go
// internal/lower/lower.go
package lower

func LowerProgram(p *ast.Program) (*ir.Hypothesis, error) {
    h := &ir.Hypothesis{
        ID:   generateID(p),
        Goal: extractGoal(p),
    }
    for _, decl := range p.Declarations {
        switch d := decl.(type) {
        case *ast.StageDecl:
            h.Concepts = append(h.Concepts, lowerStage(d)...)
        case *ast.DeriveDecl:
            h.Claims = append(h.Claims, lowerDerive(d)...)
        case *ast.IdentifyDecl:
            h.Claims = append(h.Claims, lowerIdentify(d)...)
        case *ast.SilentDecl:
            h.Silences = append(h.Silences, lowerSilent(d)...)
        case *ast.AssumeDecl:
            h.Assumptions = append(h.Assumptions, lowerAssume(d)...)
        // ...
        }
    }
    return h, nil
}
```

### Stage 6 — Program A Semantic Passes

**Input:** Physics IR.
**Output:** Semantic verdict.

```go
// internal/sema/program_a.go
package sema

func RunProgramA(h *ir.Hypothesis) *Verdict {
    v := &Verdict{HypothesisID: h.ID}
    
    // Pass 1: Component assignment
    v.ComponentErrors = checkComponents(h)
    
    // Pass 2: Mode consistency
    v.ModeErrors = checkModes(h)
    
    // Pass 3: Circularity
    v.CircularityErrors = checkCircularity(h)
    
    // Pass 4: Masquerade detection
    v.MasqueradeErrors = checkMasquerades(h)
    
    // Pass 5: Regime consistency
    v.RegimeErrors = checkRegimes(h)
    
    // Pass 6: Silence violations
    v.SilenceErrors = checkSilences(h)
    
    // Pass 7: Native mathematics jurisdiction
    v.JurisdictionErrors = checkJurisdiction(h)
    
    // Pass 8: Obligation generation
    v.Obligations = generateObligations(h)
    
    // Overall status
    v.Status = computeStatus(v)
    
    return v
}
```

**The semantic passes:**

| Pass | What It Checks | Example Error |
|---|---|---|
| **Component** | Every concept has a component. | `MissingComponent: Ω` |
| **Mode** | Declared modes match dependencies. | `ModeMismatch: DERIVE used for CONSTRUCT` |
| **Circularity** | Derivation graph is acyclic. | `Circular: A derives B, B derives A` |
| **Masquerade** | Modes are not secretly stronger. | `Masquerade: REPRESENT declared as IDENTIFY` |
| **Regime** | Claims do not silently change regime. | `RegimeEscape: UV result claimed in IR` |
| **Silence** | Silences are not violated. | `SilenceViolation: attempts to derive BornRule` |
| **Jurisdiction** | Math is used natively. | `Jurisdiction: Topology used to generate dynamics` |
| **Obligation** | Unresolved items are flagged. | `Obligation: ExplicitBridge required for IDENTIFY` |

### Stage 7 — Backend Dispatch

**Input:** Semantic verdict + Physics IR.
**Output:** Backend-specific artifacts.

```go
// internal/backend/backend.go
package backend

type Backend interface {
    Name() string
    CanHandle(h *ir.Hypothesis) bool
    Emit(h *ir.Hypothesis) ([]byte, error)
}

type ProgramBBackend struct{}
func (b *ProgramBBackend) Emit(h *ir.Hypothesis) ([]byte, error) {
    // Emit execution plan for Program B.
}

type LeanBackend struct{}
func (b *LeanBackend) Emit(h *ir.Hypothesis) ([]byte, error) {
    // Emit Lean proof obligations.
}

type CASBackend struct{}
func (b *CASBackend) Emit(h *ir.Hypothesis) ([]byte, error) {
    // Emit symbolic computation.
}
```

---

## Part IV: The Go Package Structure

```
physicsc/
├── cmd/
│   └── physicsc/
│       └── main.go
├── internal/
│   ├── lexer/
│   │   ├── token.go
│   │   └── lexer.go
│   ├── parser/
│   │   ├── parser.go
│   │   └── precedence.go
│   ├── ast/
│   │   ├── ast.go
│   │   └── print.go
│   ├── types/
│   │   ├── types.go
│   │   ├── physical.go
│   │   └── unify.go
│   ├── dimensions/
│   │   ├── dimensions.go
│   │   └── arithmetic.go
│   ├── ir/
│   │   ├── ir.go
│   │   └── validate.go
│   ├── sema/
│   │   ├── typecheck.go
│   │   ├── dimcheck.go
│   │   ├── program_a.go
│   │   ├── circularity.go
│   │   ├── masquerade.go
│   │   ├── regime.go
│   │   ├── silence.go
│   │   └── jurisdiction.go
│   ├── lower/
│   │   └── lower.go
│   ├── backend/
│   │   ├── backend.go
│   │   ├── program_b.go
│   │   ├── lean.go
│   │   └── cas.go
│   └── diag/
│       ├── diagnostic.go
│       └── repair.go
├── language/
│   └── stdlib/
│       ├── bm.phys
│       ├── cdt.phys
│       ├── cst.phys
│       └── as.phys
├── examples/
│   ├── problem_of_time.phys
│   ├── backreaction.phys
│   └── emergence_spacetime.phys
└── go.mod
```

---

## Part V: The Standard Library

The standard library provides reference semantic packs for the base frameworks.

```physics
// language/stdlib/bm.phys
module BohmianMechanics {
    concept Configuration : Dancer
    concept WaveFunction : Dancer
    concept Guidance : Dance
    concept QuantumEquilibrium : Joint
    
    postulate ActualConfiguration {
        Q ∈ Configuration
    }
    
    postulate GuidanceLaw {
        Q̇ = j / ρ
    }
    
    assume Equivariance {
        ρ₀ = |Ψ₀|² → ρ_t = |Ψ_t|²
    }
    
    silent BornRule
}
```

```physics
// language/stdlib/cst.phys
module CausalSetTheory {
    stage CausalSet {
        primitive: Event
        structure: PartialOrder(≺)
        property: LocallyFinite
    }
    
    field Cardinality : (Event, Event) → Nat
    
    derive DiscreteVolume {
        from Cardinality
        yields V : (Interval) → Scalar
        where V(I) ≈ Cardinality(I) / ϱ
    }
    
    silent SpacetimeDimension
    silent LorentzianSignature
}
```

```physics
// language/stdlib/cdt.phys
module CausalDynamicalTriangulations {
    stage Triangulation {
        primitive: Simplex
        structure: CausalGluing
        property: Foliated
    }
    
    action Regge {
        term: (1/8πG) * Σ A_h δ_h
        regime: microscopic
    }
    
    assume Ensemble {
        Z = Σ_T exp(i S_Regge[T])
    }
    
    silent ContinuumLimit
}
```

```physics
// language/stdlib/as.phys
module AsymptoticSafety {
    field EffectiveAction : (Scale) → Action
    
    derive RGFlow {
        from EffectiveAction
        yields β : (Coupling) → Coupling
        where β(g) = k dg/dk
    }
    
    assume UVFixedPoint {
        β(g*) = 0
    }
    
    silent ExistenceOfFixedPoint
}
```

---

## Part VI: A Worked Example

### Source: Problem of Time

```physics
module ProblemOfTime {
    import BohmianMechanics
    import CausalSetTheory
    import CausalDynamicalTriangulations
    import AsymptoticSafety
    
    // H1: Causal order is primary.
    derive CausalOrder {
        from CausalSetTheory.CausalSet
        yields ≺ : PartialOrder
    }
    
    // H2: Cardinality is discrete volume.
    derive DiscreteVolume {
        from CausalSetTheory.Cardinality
        yields V : (Interval) → Scalar
        where V(I) ≈ Cardinality(I) / ϱ
        regime: microscopic
    }
    
    // H3: Causal order reconstructs the conformal metric.
    derive ConformalMetric {
        from ≺
        using MalamentTheorem
        yields g_hat : ConformalClass(Metric)
        regime: continuum
    }
    
    // H4: Cardinality reconstructs the conformal factor.
    derive ConformalFactor {
        from V
        using FixedPointCondition
        yields Ω : Scalar
        regime: continuum
    }
    
    // H5: Proper time is derived from the metric.
    derive ProperMetric {
        from g_hat, Ω
        yields g : Metric
        where g = Ω² * g_hat
    }
    
    derive ProperTime {
        from g
        yields τ : (Path) → Scalar
        where τ(γ) = ∫_γ √(-g(dx, dx))
    }
    
    // H6: Physical time is proper time along causal paths.
    identify PhysicalTime {
        as ProperTime
        along ClockWorldline
        requires ExplicitBridge
    }
    
    // H7: The Dance is consistent with the reconstructed time.
    action EinsteinHilbert {
        term: (1/16πG) * ∫ R √(-g) d⁴x
        regime: IR
    }
    
    derive Dynamics {
        from EinsteinHilbert
        yields EOM : (Field) → Equation
        where δS/δφ = 0
    }
    
    // H8: The measure is consistent with the reconstructed time.
    silent BornRule
    
    interpret Measure {
        as |Ψ|²
        where Ψ : WaveFunctional
        regime: IR
    }
}
```

### Compilation Output

```
$ physicsc problem_of_time.phys

[LEXER]      OK — 247 tokens
[PARSER]     OK — 12 declarations
[TYPECHECK]  OK — all types resolved
[DIMCHECK]   OK — all dimensions consistent
[LOWERING]   OK — 12 concepts, 8 claims, 3 assumptions, 2 silences
[PROGRAM A]  ...
  [Component]    OK — all concepts assigned
  [Mode]         OK — modes consistent
  [Circularity]  WARNING — H3 depends on H1; H4 depends on H2; H5 depends on H3 and H4
  [Masquerade]   OK — no masquerades detected
  [Regime]       OK — regimes declared
  [Silence]      OK — BornRule declared silent
  [Jurisdiction] OK — math used natively
  [Obligation]   2 obligations:
                 - ExplicitBridge for PhysicalTime.identify
                 - FixedPointCondition for ConformalFactor.derive
[VERDICT]    WELL_FORMED UNDER SPEC v0.1
             Status: WELL_FORMED
             Obligations: 2 open
             Output: problem_of_time.ir.json
```

### IR Output (Partial)

```json
{
  "id": "problem_of_time",
  "goal": "Investigate whether causal order can provide physical time",
  "concepts": [
    {"name": "CausalOrder", "component": "STAGE", "mathType": "PartialOrder"},
    {"name": "PhysicalTime", "component": "STAGE", "mathType": "ProperTime"},
    {"name": "ProperMetric", "component": "STAGE", "mathType": "Metric"}
  ],
  "claims": [
    {
      "id": "H3",
      "mode": "DERIVE",
      "relation": "DERIVES",
      "subject": "CausalOrder",
      "object": "ConformalMetric",
      "regime": "continuum",
      "obligations": []
    },
    {
      "id": "H6",
      "mode": "IDENTIFY",
      "relation": "IDENTIFIES",
      "subject": "ProperTime",
      "object": "PhysicalTime",
      "requires": ["ExplicitBridge"],
      "obligations": ["ExplicitBridge"]
    }
  ],
  "silences": ["BornRule"],
  "assumptions": ["SprinklingDensity", "Equivariance"]
}
```

---

## Part VII: Design Principles and Anti-Patterns

### Principles

1. **The compiler does not compute.** It lowers to IR. Program B computes. This is the A/B boundary.
2. **Every construct has semantics.** `derive` is not `construct`. `identify` is not `represent`. The language enforces the distinction.
3. **Types carry physical role.** `PhysicalTime` ≠ `RGScale` ≠ `CoordinateTime`. The type system enforces the distinction.
4. **Dimensions are checked.** `Length + Time` fails at compile time.
5. **Regimes are first-class.** Silent regime movement is a compile error.
6. **Silences are declared.** `silent BornRule` is a declaration. The compiler flags attempts to violate it.
7. **Assumptions are explicit.** `assume` is a keyword. No hidden premises.
8. **Obligations are explicit.** The compiler does not silently accept. It flags what is unresolved.

### Anti-Patterns

| Anti-Pattern | Why It's Wrong |
|---|---|
| **LaTeX with a parser** | The language must encode physics semantics, not math syntax. |
| **Everything is a tensor** | Physical roles matter. `State` ≠ `Vector`. |
| **No distinction between modes** | `derive` and `construct` have different obligations. |
| **Silent imports** | Imports must be declared. Hidden imports are the original sin. |
| **No regime tracking** | Regimes are where category errors hide. |
| **No silences** | Silences are the constitution. The compiler must respect them. |
| **No obligations** | The compiler does not decide. It flags what requires external adjudication. |

---

## Part VIII: The One-Sentence Summary

> **A physics compiler in Go is a multi-pass system: lexer → parser → type checker → dimension checker → lowering → Program A semantic passes → backend dispatch. The language expresses physics semantics, not just math syntax: `derive` ≠ `construct`, `identify` requires `ExplicitBridge`, `silent` declares silences, `assume` declares imports, `regime` declares scope. Types carry physical role: `PhysicalTime` ≠ `RGScale` ≠ `CoordinateTime`. Dimensions are checked at compile time. The compiler does not compute. It lowers to a Physics IR that Program A checks for semantic coherence and that Program B executes. The compiler is the frontend. Program A is the semantic backend. Program B is the execution backend. The Go SDK becomes one frontend among others. The language is discovered from actual hypotheses, not designed in the abstract. The result is a physics compiler that is honest about its imports, explicit about its silences, and testable in its claims. The constitution's prime directive is "Let reality choose the mathematics." The compiler follows this. It does not choose the mathematics. It chooses the tests. Reality will choose the theory.**

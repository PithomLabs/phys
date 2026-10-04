# The Physics Compiler: A Clean Go Design

## The Central Design Principle

There are two completely different things that have been conflated in every previous iteration:

**Human conceptual scaffolding** — "Dancer, Dance, Stage," "Jurisdiction Matrix," "Generator vs. Invariant," "Reality-First Constitution." These are *mental models* that help physicists think clearly. They belong in documentation, not in compiler code.

**Compiler-enforceable structure** — Types, dimensions, domains, codomains, operation compatibility, dependency acyclicity. These are *formal constraints* that a machine can check deterministically.

The compiler does not know what a "Dancer" is. It knows that a `PartialOrder` cannot be the domain of a `Differentiate` operation. That single type-level fact catches the category error that the "Jurisdiction Matrix" was trying to express philosophically.

The compiler does not litigate, adjudicate, or interpret. It parses, type-checks, and emits. Humans and other tools interpret the output.

---

## 1. Architecture: Following the Go Compiler

The Go compiler itself follows a strict pipeline: `scanner → parser → AST → type checker → IR → backend`. We mirror this exactly.

```
program-a/
  cmd/
    program-a/
      main.go                  # Entry point
  internal/
    token/
      token.go                 # Token definitions
    scanner/
      scanner.go               # Lexer
    ast/
      ast.go                   # Abstract syntax tree nodes
    parser/
      parser.go                # Recursive descent parser
    types/
      types.go                 # Math and physics type system
      universe.go              # Built-in types and operations
    typecheck/
      typecheck.go             # Type checking pass
    semantics/
      semantics.go             # Semantic consistency pass
    ir/
      ir.go                    # Intermediate representation
    lower/
      lower.go                 # AST → IR lowering
    emit/
      emit.go                  # JSON/YAML artifact output
  spec/
    types.yaml                 # Declarative type definitions
  examples/
    problem-of-time.phys
    backreaction.phys
```

No database. No distributed system. No executor. No EBP. No adversarial review. No context graph. No Solvent integration. Just a compiler.

---

## 2. The Type System: Where the Real Work Lives

This is the heart of the compiler. Every mathematical object and physical concept gets a precise type. The type system itself prevents category errors — no philosophical rules needed.

### 2.1 Math Types

```go
package types

// MathKind is the fundamental category of a mathematical object.
type MathKind int

const (
    Scalar MathKind = iota
    Vector
    Covector
    Tensor
    Operator
    Functional
    Manifold
    Metric
    Connection
    Curvature
    PartialOrder
    Graph
    Group
    Algebra
    Measure
    Topology
    Spectrum
    Set
)

// Dimension represents physical dimensions for dimensional analysis.
type Dimension struct {
    Length    int // power of L
    Mass      int // power of M
    Time      int // power of T
    Charge    int // power of Q
    Temperature int // power of Θ
}

// Regime marks the scale domain of an object.
type Regime int

const (
    RegimeUnspecified Regime = iota
    RegimeUV
    RegimeIR
    RegimeContinuum
    RegimeDiscrete
    RegimeMesoscopic
)

// MathType is the full type of a mathematical object.
type MathType struct {
    Kind      MathKind
    Domain    *MathType // for maps/operators: input type
    Codomain  *MathType // for maps/operators: output type
    Dim       Dimension // physical dimensions
    Regime    Regime
    Name      string    // human-readable label
}
```

### 2.2 Physics Roles

This is where we encode the "Dancer/Dance/Stage" distinction — but as a flat, machine-checkable enum, not a philosophical framework.

```go
// PhysRole classifies what job a concept plays in the hypothesis.
type PhysRole int

const (
    RoleUnspecified PhysRole = iota
    RoleState         // a physical configuration or state
    RoleObservable    // a measurable quantity
    RoleDynamics      // an equation of motion, flow, or evolution law
    RoleConstraint    // a restriction on allowed states or dynamics
    RoleGeometry      // spacetime structure, causal structure, metric
    RoleInvariant     // a quantity preserved under some transformation
    RoleParameter     // a free or fixed numerical value
    RoleSymmetry      // a transformation group or algebra
    RoleBridge        // an explicit translation between two descriptions
)
```

### 2.3 Operations

Operations are typed. Each operation declares what `MathKind` it accepts and produces. This is where category errors are caught — not by a philosophical matrix, but by the type checker.

```go
// OpKind is a mathematical or physical operation.
type OpKind int

const (
    // Math operations
    OpAdd OpKind = iota
    OpMultiply
    OpCompose
    OpDifferentiate
    OpIntegrate
    OpTrace
    OpTensorProduct
    OpDirectSum
    OpQuotient
    OpLimit
    OpPullback

    // Physics operations (semantic modes)
    OpPostulate
    OpRepresent
    OpDerive
    OpIdentify
    OpInterpret
    OpConstruct
    OpRecover
    OpConstrain
)

// OpSignature declares what types an operation accepts and produces.
type OpSignature struct {
    Op       OpKind
    Inputs   []MathKind // accepted input kinds
    Output   MathKind   // produced output kind
    MinArgs  int
    MaxArgs  int        // -1 = variadic
}
```

### 2.4 The Operation Table (The Real "Jurisdiction Matrix")

This table replaces the philosophical Jurisdiction Matrix. It says: *these operations are defined for these types*. If you try to `Differentiate` a `PartialOrder`, the type checker rejects it — not because a rule says "Arithmetic is silent on Dynamics," but because `Differentiate` simply does not accept `PartialOrder` as input.

```go
var OperationTable = []OpSignature{
    // Calculus operations require continuous types
    {OpDifferentiate, []MathKind{Scalar, Vector, Tensor, Functional, Manifold}, Scalar, 1, 2},
    {OpIntegrate, []MathKind{Scalar, Vector, Tensor, Functional}, Scalar, 1, 2},
    {OpLimit, []MathKind{Scalar, Vector, Tensor, Functional, Spectrum}, Scalar, 1, 1},

    // Algebraic operations
    {OpCompose, []MathKind{Operator, Functional}, Operator, 2, 2},
    {OpTrace, []MathKind{Operator, Tensor}, Scalar, 1, 1},
    {OpTensorProduct, []MathKind{Vector, Tensor, Algebra}, Tensor, 2, -1},

    // Geometric operations
    {OpPullback, []MathKind{Metric, Curvature, Tensor}, Tensor, 2, 2},

    // Physics semantic operations
    {OpDerive, nil, Scalar, 1, -1},       // checked by semantic pass
    {OpRepresent, nil, Scalar, 2, 2},     // source → target
    {OpIdentify, nil, Scalar, 2, 2},      // stronger than represent
    {OpInterpret, nil, Scalar, 2, 2},     // math → physics
    {OpConstrain, nil, Scalar, 2, 2},     // invariant → dynamics
}
```

The key insight: `OpDifferentiate` does not list `PartialOrder`, `Group`, `Set`, `Topology`, or `Spectrum` as valid inputs. The type checker enforces this mechanically. No philosophy required.

---

## 3. The AST: What the Parser Produces

```go
package ast

import "program-a/internal/token"

// A File is the top-level node.
type File struct {
    Decls []Decl
}

// Decl is a top-level declaration.
type Decl interface {
    declNode()
    Pos() token.Pos
}

// ObjectDecl declares a mathematical or physical object.
// Example: object causal_order : PartialOrder { role geometry; regime discrete }
type ObjectDecl struct {
    Name   string
    Type   string // references a MathKind by name
    Role   string // references a PhysRole by name
    Regime string
    Props  map[string]string
}

// OperationDecl declares an operation on objects.
// Example: derive physical_time from causal_order using differentiate
type OperationDecl struct {
    Mode    string   // "derive", "represent", "identify", etc.
    Target  string   // name of the result
    Sources []string // names of inputs
    Method  string   // operation used
}

// BridgeDecl declares an explicit translation between descriptions.
// Example: bridge causal_to_metric : causal_order → metric_tensor via reconstruction
type BridgeDecl struct {
    Name   string
    From   string
    To     string
    Method string
}

// AssumptionDecl declares an assumption.
// Example: assume local_finiteness of causal_order
type AssumptionDecl struct {
    Target string
    Prop   string
}
```

---

## 4. The Compiler Pipeline

### 4.1 Scanner and Parser

Standard recursive descent. The scanner tokenizes the `.phys` source. The parser produces the AST. This is routine compiler engineering — no physics-specific logic here.

### 4.2 Type Checker

This is where math-level category errors are caught.

```go
package typecheck

func Check(file *ast.File, universe *types.Universe) []Error {
    var errs []Error

    // Pass 1: Resolve all object declarations to MathTypes
    env := make(map[string]*types.MathType)
    for _, decl := range file.Decls {
        if obj, ok := decl.(*ast.ObjectDecl); ok {
            mt, err := universe.Resolve(obj.Type, obj.Role, obj.Regime)
            if err != nil {
                errs = append(errs, Error{Pos: obj.Pos(), Msg: err.Error()})
                continue
            }
            env[obj.Name] = mt
        }
    }

    // Pass 2: Check all operations for type compatibility
    for _, decl := range file.Decls {
        if op, ok := decl.(*ast.OperationDecl); ok {
            sig, ok := universe.LookupOp(op.Method)
            if !ok {
                errs = append(errs, Error{Pos: op.Pos(), Msg: "unknown operation: " + op.Method})
                continue
            }
            for _, src := range op.Sources {
                srcType, ok := env[src]
                if !ok {
                    errs = append(errs, Error{Pos: op.Pos(), Msg: "undefined object: " + src})
                    continue
                }
                if !sig.Accepts(srcType.Kind) {
                    errs = append(errs, Error{
                        Pos: op.Pos(),
                        Msg: fmt.Sprintf("type error: operation %s does not accept %s (kind %s)",
                            op.Method, src, srcType.Kind),
                    })
                }
            }
        }
    }

    // Pass 3: Dimensional analysis
    // ... check that Add/Multiply operands have compatible dimensions ...

    return errs
}
```

**What this catches:**
- `differentiate(causal_order)` → type error: `PartialOrder` not in `Differentiate`'s input list
- `integrate(spin_group)` → type error: `Group` not in `Integrate`'s input list
- `add(length, time)` → dimension error: `L ≠ T`
- `compose(partial_order, metric)` → type error: `PartialOrder` not in `Compose`'s input list

No philosophical rules. No "Jurisdiction Matrix." Just types.

### 4.3 Semantic Checker

This is where physics-level consistency is checked. It operates on the typed AST.

```go
package semantics

func Check(file *ast.File, env map[string]*types.MathType) []Error {
    var errs []Error

    // Pass 1: Check for circular derivation dependencies
    graph := buildDerivationGraph(file)
    if cycles := findCycles(graph); len(cycles) > 0 {
        for _, c := range cycles {
            errs = append(errs, Error{Msg: "circular derivation: " + formatCycle(c)})
        }
    }

    // Pass 2: Check that IDENTIFY operations have explicit justification
    for _, decl := range file.Decls {
        if op, ok := decl.(*ast.OperationDecl); ok && op.Mode == "identify" {
            if !hasExplicitBridge(file, op.Sources[0], op.Target) {
                errs = append(errs, Error{
                    Pos: op.Pos(),
                    Msg: fmt.Sprintf("IDENTIFY requires explicit bridge: %s → %s", op.Sources[0], op.Target),
                })
            }
        }
    }

    // Pass 3: Check that DERIVE chains don't secretly import their conclusion
    for _, decl := range file.Decls {
        if op, ok := decl.(*ast.OperationDecl); ok && op.Mode == "derive" {
            if importsConclusion(file, op, env) {
                errs = append(errs, Error{
                    Pos: op.Pos(),
                    Msg: "import-as-derivation: target property already present in assumptions",
                })
            }
        }
    }

    return errs
}
```

**What this catches:**
- `derive time from causal_order` where `causal_order` was defined with an assumption that already contains temporal structure → import-as-derivation
- `identify rg_scale as physical_time` without an explicit bridge → missing bridge
- `derive A from B; derive B from C; derive C from A` → circular dependency

### 4.4 IR Lowering and Emission

The typed, semantically checked AST is lowered to a flat IR and emitted as JSON.

```go
package ir

type Hypothesis struct {
    Objects    []Object    `json:"objects"`
    Operations []Operation `json:"operations"`
    Bridges    []Bridge    `json:"bridges"`
    Assumptions []string   `json:"assumptions"`
    Obligations []string   `json:"obligations"`
}

type Object struct {
    Name     string `json:"name"`
    MathKind string `json:"math_kind"`
    PhysRole string `json:"phys_role"`
    Regime   string `json:"regime"`
    Dim      string `json:"dimension"`
}

type Operation struct {
    Mode    string   `json:"mode"`
    Target  string   `json:"target"`
    Sources []string `json:"sources"`
    Method  string   `json:"method"`
}
```

This JSON is the output. Program B reads it and decides what to compute. A human reads it and decides what it means. The compiler does neither.

---

## 5. What the Compiler Does NOT Do

This list is as important as what it does do:

| Excluded | Why |
|---|---|
| EBP / evidence tracking | Not the compiler's job. That's an external system. |
| Adversarial review | Not the compiler's job. That's a separate AI agent. |
| Null separation tests | Not the compiler's job. That's Program B. |
| "Dancer/Dance/Stage" as compiler types | Human mental model. The compiler uses `PhysRole` enum. |
| "Jurisdiction Matrix" as a hardcoded map | Replaced by operation signatures in the type system. |
| "Generator vs. Invariant" as a rule | Replaced by type constraints: invariants have kind `Spectrum`/`Topology`, which are not valid inputs to dynamics operations. |
| Computing consequences | Program B's job. |
| Interpreting results | Human's job. |
| Judging truth | Nobody's job at compile time. |
| Maintaining research context | Solvent's job. |
| Tracking hypothesis versions | Git's job. |

---

## 6. The Language: What a `.phys` File Looks Like

```physics
// problem-of-time.phys

object causal_order : PartialOrder {
    role geometry
    regime discrete
}

object local_finiteness : Constraint {
    role constraint
    target causal_order
}

object spacetime_metric : Metric {
    role geometry
    regime continuum
}

object rg_scale : Scalar {
    role parameter
    regime uv
}

object relational_clock : Scalar {
    role observable
    regime ir
    dim T
}

assume local_finiteness of causal_order

bridge causal_to_metric : causal_order → spacetime_metric
    via reconstruction
    obligation continuum_limit
    obligation lorentzian_signature

derive relational_clock from spacetime_metric
    using integration

// This would FAIL type checking:
// derive physical_time from causal_order using differentiate
// Error: Differentiate does not accept PartialOrder
```

---

## 7. The Complete Pipeline

```text
.phys source file
      │
      ▼
  ┌─────────┐
  │ Scanner │  tokenizes
  └────┬────┘
       ▼
  ┌─────────┐
  │ Parser  │  produces AST
  └────┬────┘
       ▼
  ┌─────────────┐
  │ Type Checker│  checks math types, dimensions, operation signatures
  └────┬────────┘
       ▼
  ┌─────────────────┐
  │ Semantic Checker│  checks circularity, import-as-derive, missing bridges
  └────┬────────────┘
       ▼
  ┌─────────┐
  │ Lowerer │  AST → IR
  └────┬────┘
       ▼
  ┌─────────┐
  │ Emitter │  IR → JSON artifact
  └────┬────┘
       ▼
  hypothesis.json  →  Program B / Human analysis
```

---

## 8. Summary of Design Principles

1. **Conceptual ≠ Compiler.** "Dancer/Dance/Stage" is documentation. The compiler uses `PhysRole` and `MathKind` enums. The "Jurisdiction Matrix" is documentation. The compiler uses operation signatures.

2. **Types prevent category errors.** You cannot differentiate a partial order because `Differentiate` does not accept `PartialOrder`. No philosophical rule needed.

3. **The compiler does not interpret.** It parses, type-checks, and emits. Humans and other tools interpret the output.

4. **Standard Go compiler architecture.** Scanner → Parser → AST → Type Checker → Semantic Checker → IR → Emitter. No exotic patterns.

5. **Minimal.** No database, no distributed system, no executor, no EBP, no adversarial review, no context graph. Just a compiler.

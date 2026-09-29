# THE PHYSICS COMPILER — CORE DESIGN PRINCIPLES
## Separation of Human Mental Models from Compiler Concerns

**Status:** Foundational design clarification. The compiler is a tool, not a philosopher.

---

## Part I: The Separation Principle

The design has been conflating two distinct things:

$$
\boxed{
\text{Human mental models} \neq \text{Compiler concerns}
}
$$

| Human Mental Models | Compiler Concerns |
|---|---|
| The Dancer/Dance/Stage metaphor | The grammar of the language |
| The "where math shines" narrative | The type system |
| The jurisdiction matrix as philosophical prose | The semantic rules |
| The "reality-first" methodology | The representation of math structures |
| The constitution as moral framework | The representation of physics structures |
| The nine maps as survey documents | The error taxonomy |
| The "silences" as philosophical commitments | The standard library |
| The promotion states as research methodology | The analysis tools |

**The compiler does not have a philosophy.** It has a grammar, a type system, a semantic checker, a standard library, and analysis tools. That is all.

**The compiler does not adjudicate physics.** It represents what the human means and gives the human tools to analyze the representation. It does not decide whether the physics is "true."

**The compiler does not interpret the program.** It parses, type-checks, semantically analyzes, and reports. The interpretation is the human's job.

This is the Go compiler analogy:

| Go Compiler | Physics Compiler |
|---|---|
| Checks syntax | Checks syntax |
| Checks types | Checks types (math + physics) |
| Checks scoping | Checks scoping |
| Checks interface satisfaction | Checks jurisdiction compliance |
| Provides `go vet` | Provides `phys vet` |
| Provides `go test` | Provides `phys test` |
| Provides standard library | Provides standard library |
| Does NOT decide if your algorithm is efficient | Does NOT decide if your physics is true |
| Does NOT decide if your business logic is correct | Does NOT decide if your theory is correct |
| Does NOT interpret your program | Does NOT interpret your program |

---

## Part II: The Compiler's Actual Job

The compiler has exactly five jobs:

1. **Represent human intent.** The human writes a hypothesis in the physics language. The compiler represents it as a typed AST, then as a physics IR.
2. **Enforce syntactic guidelines.** The grammar defines what is a well-formed program.
3. **Enforce semantic guidelines.** The type system and semantic checker define what is a coherent program.
4. **Provide analysis tools.** The compiler provides tools for the human to analyze the program: dependency graph, obligation report, jurisdiction report, regime report, etc.
5. **Output the representation.** The compiler outputs the representation in a form that Program B can execute, or that the human can inspect.

The compiler does **not** have these jobs:

- Adjudicate physics
- Litigate papers
- Decide truth
- Interpret results
- Select theories
- Promote claims
- Maintain evidence
- Track research methodology

---

## Part III: The Grammar of Mathematics

The compiler must represent mathematical structures as a language. The grammar of mathematics is:

### 3.1 Math Primitives

```
math Scalar
math Vector
math Covector
math Tensor
math Operator
math Function
math Functional
math Set
math Relation
math PartialOrder
math Graph
math Manifold
math Metric
math Connection
math Curvature
math Topology
math Measure
math ProbabilityMeasure
math Group
math Algebra
math Category
math Spectrum
math Invariant
math Parameter
math Limit
math Variation
math Expectation
math Correlation
math Distribution
```

### 3.2 Math Operations

```
op differentiate : Function → Function
op integrate : Function → Scalar
op vary : Functional → Function
op substitute : (Expr, Symbol, Expr) → Expr
op simplify : Expr → Expr
op factor : Expr → Expr
op linearize : Expr → Expr
op quantize : Classical → Quantum
op coarse_grain : Micro → Macro
op discretize : Continuous → Discrete
op continuum_limit : Discrete → Continuous
op take_expectation : (State, Observable) → Scalar
op project : (State, Subspace) → State
op commute : (Operator, Operator) → Operator
op compose : (Function, Function) → Function
op compare : (Expr, Expr) → Relation
```

### 3.3 Math Relations

```
rel equals : (Expr, Expr)
rel approximates : (Expr, Expr, Error)
rel is_limit_of : (Expr, Expr)
rel is_dual_to : (Structure, Structure)
rel is_representation_of : (Structure, Structure)
rel is_isomorphic_to : (Structure, Structure)
rel is_substructure_of : (Structure, Structure)
rel is_extension_of : (Structure, Structure)
```

### 3.4 Math Types

Every math primitive has a type. Every operation has a domain and codomain. The compiler checks type compatibility.

```
type Scalar : MathType
type Vector : MathType
type Tensor : MathType
type Operator : MathType
...
```

---

## Part IV: The Grammar of Physics

The compiler must represent physics structures as a language. The grammar of physics is:

### 4.1 Physics Primitives

```
physics Particle
physics Field
physics State
physics Observable
physics Configuration
physics Action
physics Lagrangian
physics Hamiltonian
physics Constraint
physics Generator
physics Flow
physics Spacetime
physics Metric
physics Connection
physics Curvature
physics Symmetry
physics Conservation
physics Charge
physics Current
physics Energy
physics Momentum
physics Entropy
physics Temperature
physics OrderParameter
physics Phase
physics Defect
physics Soliton
physics Instanton
physics Anyon
physics Spectrum
physics Scattering
physics Decay
physics Coupling
physics Mass
physics Spin
physics Gauge
physics Vacuum
physics GroundState
physics ExcitedState
physics Measurement
physics Record
physics System
physics Subsystem
physics Environment
physics Observer
physics Regime
physics Limit
physics CouplingConstant
physics InitialCondition
physics BoundaryCondition
```

### 4.2 Physics Operations

```
op evolve : State → State
op measure : (State, Observable) → Record
op prepare : (System, Procedure) → State
op scatter : (State, State) → Amplitude
op decay : State → State
op interact : (Field, Field) → Interaction
op break_symmetry : (Symmetry, Vacuum) → Vacuum
op gauge : (Field, GaugeGroup) → Field
op quantize : ClassicalField → QuantumField
op renormalize : (Theory, Scale) → Theory
op flow : (Coupling, Scale) → Coupling
op condense : (Field, Vacuum) → Vacuum
op excite : (GroundState, Quantum) → ExcitedState
op annihilate : (State, State) → State
op create : (Vacuum, Quantum) → State
op propagate : (Field, Spacetime) → Field
op couple : (Field, Field, Coupling) → Interaction
op backreact : (Matter, Geometry) → Geometry
op reconstruct : (Data, Theory) → Structure
op emerge : (Micro, Limit) → Macro
```

### 4.3 Physics Relations

```
rel conserves : (Current, Symmetry)
rel breaks : (Symmetry, Vacuum)
rel couples_to : (Field, Field)
rel is_gauge_of : (Force, Group)
rel is_representation_of : (Particle, Group)
rel is_limit_of : (Theory, Theory)
rel is_dual_to : (Theory, Theory)
rel is_regime_of : (Physics, Regime)
rel is_consistent_with : (Theory, Observation)
rel is_derived_from : (Result, Premise)
rel is_assumed_in : (Assumption, Theory)
```

### 4.4 Physics Types

Every physics primitive has a type. Every operation has a domain and codomain. The compiler checks type compatibility.

```
type Particle : PhysicsType
type Field : PhysicsType
type State : PhysicsType
type Observable : PhysicsType
...
```

---

## Part V: The Mapping Between Grammars

The compiler must map between the math grammar and the physics grammar. This is where the jurisdiction kernel lives.

### 5.1 The Mapping Table

```
mapping MathToPhysics {
    Scalar          → PhysicalQuantity
    Vector          → PhysicalVector
    Tensor          → PhysicalTensor
    Operator        → Observable
    Function        → Field
    Functional      → Action
    Manifold        → Spacetime
    Metric          → SpacetimeMetric
    Connection      → GaugeField
    Curvature       → FieldStrength
    Symmetry        → PhysicalSymmetry
    Spectrum        → PhysicalSpectrum
    Invariant       → PhysicalInvariant
    ...
}
```

### 5.2 The Jurisdiction Check

For every operation in the physics grammar, the compiler checks:

1. **The math operation is well-typed.** (Math type check)
2. **The physics operation is well-typed.** (Physics type check)
3. **The math operation is native to the physics operation's component.** (Jurisdiction check)
4. **The math operation is native to the physics operation's regime.** (Jurisdiction check)
5. **The math operation does not violate a silence.** (Silence check)

### 5.3 The Error Report

If any check fails, the compiler reports:

```
JURISDICTION-ERROR
  Operation: generate_dynamics
  MathBranch: Topology
  Job: Generate
  Component: Dance
  Regime: IR
  Reason: Topology is silent on the Generate job.
  NativeMath: Calculus
  Suggestion: Use a calculus operation, or reclassify the job.
```

---

## Part VI: Syntactic Guidelines

The grammar enforces:

1. **Every claim has a mode.** (derive, represent, identify, interpret, construct, emerge, recover, etc.)
2. **Every claim has a regime.** (microscopic, IR, continuum, UV, classical, quantum, etc.)
3. **Every claim has a component.** (Stage, Dance, Dancer, Invariant, Joint)
4. **Every claim has a math branch.** (Calculus, Geometry, Algebra, Topology, Arithmetic, Spectra, Symmetry, Duality, Emergence, Entropy, Complexity, DegreesOfFreedom, Langlands)
5. **Every claim has a source.** (definition, theorem, assumption, conjecture, etc.)
6. **Every import is declared.** (No hidden imports)
7. **Every silence is declared.** (No hidden silences)
8. **Every assumption is explicit.** (No hidden assumptions)

If any of these are missing, the compiler reports a syntax error.

---

## Part VII: Semantic Guidelines

The semantic checker enforces:

1. **Math type compatibility.** (No scalar + vector)
2. **Physics type compatibility.** (No particle + field)
3. **Jurisdiction compatibility.** (No topology-as-dynamics)
4. **Regime compatibility.** (No IR claim in UV regime)
5. **Mode compatibility.** (No derive-as-represent)
6. **Component compatibility.** (No Stage-as-Dancer)
7. **Silence compatibility.** (No derive-Born-rule)
8. **Import compatibility.** (No undeclared import)
9. **Representation compatibility.** (No incompatible representations)
10. **Direction compatibility.** (No inverted direction of explanation)

If any of these fail, the compiler reports a semantic error.

---

## Part VIII: The Standard Library

The compiler provides a standard library of math and physics constructs. The standard library is:

1. **Math constructs.** Types, operations, relations for each math branch.
2. **Physics constructs.** Types, operations, relations for each physics domain.
3. **Jurisdiction constructs.** The mapping between math and physics.
4. **Regime constructs.** The regime taxonomy and transitions.
5. **Constraint constructs.** The constraint taxonomy.
6. **Symmetry constructs.** The symmetry taxonomy.
7. **Limit constructs.** The limit taxonomy.
8. **Emergence constructs.** The emergence taxonomy.
9. **Measurement constructs.** The measurement taxonomy.
10. **Analysis constructs.** Tools for analyzing the program.

The standard library is versioned. It is extensible. It is auditable.

---

## Part IX: What the Compiler Does NOT Do

The compiler does **not**:

1. **Adjudicate physics.** It does not decide whether a theory is true.
2. **Litigate papers.** It does not decide whether a paper is correct.
3. **Decide truth.** It does not have a truth criterion.
4. **Interpret results.** It does not interpret the output of the program.
5. **Select theories.** It does not choose between competing theories.
6. **Promote claims.** It does not promote claims to a higher status.
7. **Maintain evidence.** It does not track evidence for or against a claim.
8. **Track research methodology.** It does not enforce a research methodology.
9. **Enforce philosophy.** It does not enforce a philosophy of science.
10. **Interpret intent.** It does not interpret the human's intent. It represents the human's intent.

The compiler's job is representation and analysis. It represents what the human means. It gives the human tools to analyze the representation. It does not interpret.

---

## Part X: The Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   HUMAN INTENT                                                  │
│                                                                 │
│   The human writes a hypothesis in the physics language.        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   PHYSICS LANGUAGE (Source)                                     │
│                                                                 │
│   The grammar of mathematics and the grammar of physics.        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   LEXER + PARSER                                                │
│                                                                 │
│   Tokenizes and parses the source into an AST.                  │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   TYPE CHECKER                                                  │
│                                                                 │
│   Checks math types and physics types.                          │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   JURISDICTION CHECKER                                          │
│                                                                 │
│   Checks math-physics mapping, regime, component, silence.      │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   SEMANTIC CHECKER                                              │
│                                                                 │
│   Checks mode, import, representation, direction.               │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   PHYSICS IR                                                    │
│                                                                 │
│   The typed, checked representation of the hypothesis.          │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   ANALYSIS TOOLS                                                │
│                                                                 │
│   Dependency graph, obligation report, jurisdiction report,    │
│   regime report, silence report, etc.                          │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│   OUTPUT                                                        │
│                                                                 │
│   The representation, the analysis, the errors.                 │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Part XI: The Go Standard Library Analogy

The Go compiler provides:

- `fmt` for formatting
- `strings` for string manipulation
- `math` for mathematical functions
- `sort` for sorting
- `io` for input/output
- `net` for networking
- `os` for operating system interface
- `sync` for synchronization
- `time` for time

The physics compiler provides:

- `math` for mathematical structures
- `physics` for physical structures
- `jurisdiction` for math-physics mapping
- `regime` for regime taxonomy
- `analysis` for analysis tools
- `report` for reporting
- `graph` for dependency graphs
- `obligation` for obligation tracking
- `error` for error reporting

The standard library is the set of pre-built constructs that the human can use without redefining them.

---

## Part XII: The Core Design Principles

1. **Represent, don't interpret.** The compiler represents what the human means. It does not interpret.
2. **Check, don't adjudicate.** The compiler checks the program against the grammar, the type system, and the semantic rules. It does not adjudicate truth.
3. **Separate mental models from compiler concerns.** The Dancer/Dance/Stage metaphor is a human mental model. The compiler does not use it. It uses types, operations, and relations.
4. **Enforce syntax and semantics.** The compiler enforces the grammar and the semantic rules. It does not enforce a philosophy.
5. **Provide tools, not verdicts.** The compiler provides analysis tools. The human uses them to draw verdicts.
6. **Follow Go's design.** The compiler is a tool. It has a grammar, a type system, a semantic checker, a standard library, and analysis tools. It does not have a philosophy.
7. **Be explicit.** Every import is declared. Every silence is declared. Every assumption is explicit. Every regime is explicit. Every component is explicit.
8. **Be extensible.** The standard library is extensible. The jurisdiction kernel is extensible. The grammar is extensible.
9. **Be auditable.** The compiler's rules are explicit and auditable. The compiler can be checked against its own rules.
10. **Be minimal.** The compiler does not do more than it needs to. It represents, checks, analyzes, and reports.

---

## The One-Paragraph Summary

> **The physics compiler is a tool for representing human intent and analyzing it, not for adjudicating physics. It has five jobs: represent human intent, enforce syntactic guidelines, enforce semantic guidelines, provide analysis tools, and output the representation. It does not adjudicate physics, litigate papers, decide truth, interpret results, select theories, promote claims, maintain evidence, track research methodology, enforce philosophy, or interpret intent. The design separates human mental models (the Dancer/Dance/Stage metaphor, the jurisdiction matrix as philosophical prose, the nine maps as survey documents) from compiler concerns (the grammar, the type system, the semantic rules, the representation of math and physics structures, the error taxonomy, the standard library, the analysis tools). The compiler has a grammar for mathematics (primitives, operations, relations, types) and a grammar for physics (primitives, operations, relations, types), and a mapping between them (the jurisdiction kernel). The compiler enforces syntactic guidelines (every claim has a mode, regime, component, math branch, source; every import is declared; every silence is declared; every assumption is explicit) and semantic guidelines (math type compatibility, physics type compatibility, jurisdiction compatibility, regime compatibility, mode compatibility, component compatibility, silence compatibility, import compatibility, representation compatibility, direction compatibility). The compiler provides a standard library of math constructs, physics constructs, jurisdiction constructs, regime constructs, constraint constructs, symmetry constructs, limit constructs, emergence constructs, measurement constructs, and analysis constructs. The compiler follows Go's design: it is a tool with a grammar, a type system, a semantic checker, a standard library, and analysis tools. It is explicit, extensible, auditable, and minimal. The compiler's job is to represent what the human means and give the human tools to analyze the representation. It is not the compiler's job to interpret the program.**

# Physics Compiler — Definitive Specification v0.1

**Status:** FROZEN REFERENCE  
**Purpose:** Normative design reference for the Physics Compiler  
**Scope:** Symbolic theoretical physics only  
**Implementation language:** Go  
**Source language:** `.phys`

---

## 0. Executive Definition

The **Physics Compiler** is a programming language and compiler for expressing **theoretical physics symbolically**, in a form analogous to writing mathematics and physics on paper, but with formal compiler guarantees.

Its primary job is to represent, check, transform, and expose the structure of symbolic mathematics and theoretical physics.

The compiler works with:

- symbols
- exact symbolic constants
- expressions
- mathematical structures
- physical quantities and structures
- equations and relations
- tensors and indices
- dimensions
- operators
- functions and mappings
- derivatives, integrals, limits, and variations
- assumptions and definitions
- compositions and transformations
- typed relationships between mathematics and physics

The compiler does **not** evaluate physical values, run simulations, execute numerical models, or adjudicate whether a physical theory is true.

The governing idea is:

> **Write theoretical physics as symbols; let the compiler enforce the formal grammar and structural relationships of those symbols.**

---

# 1. Non-Negotiable Scope

## 1.1 Symbolic Only

The compiler is concerned with **symbolic mathematical manipulation**.

Examples of valid activity:

```text
x² + 2x + 1
→ (x + 1)²
```

```text
∂μ Aν - ∂ν Aμ
```

```text
δS = 0
```

```text
Gμν + Λgμν = 8πG Tμν
```

```text
[p, q] = -iℏ
```

These are symbolic structures.

Exact coefficient normalization, simplification, collection, substitution, and other algebraic transformations are permitted when they are part of symbolic semantics.

For example:

```text
2 + 2
```

may normalize to:

```text
4
```

because this is symbolic/algebraic normalization, not execution of a physical quantity.

## 1.2 No Physical-Value Execution

The compiler does not:

- assign experimental values to physical variables
- execute equations against numerical inputs
- simulate trajectories
- solve numerical differential equations
- perform numerical PDE solving
- execute Monte Carlo calculations
- run numerical optimization
- perform parameter fitting
- generate CPU/GPU numerical kernels
- execute scientific workloads
- provide a numerical runtime
- behave as a simulation language

A downstream mathematical or computational system may later consume symbolic output, but such systems are external to the Physics Compiler.

## 1.3 No Scientific Adjudication

The compiler does not decide:

- whether a theory is true
- whether a theory is physically realized
- whether a hypothesis is experimentally supported
- whether a theory is accepted by the physics community
- whether one theory is preferable to another
- whether a paper is correct
- whether a proposed mechanism is empirically adequate

The compiler can determine formal properties of the program.

It may therefore produce:

```text
type mismatch
undefined symbol
invalid tensor contraction
dimension mismatch
unsatisfied interface
invalid operator application
missing required domain
invalid symbolic transformation
```

It must not produce:

```text
physically true
physically false
scientifically correct
scientifically invalid
accepted physics
rejected physics
```

---

# 2. Design Philosophy

The compiler follows the engineering philosophy of Go:

- small and explicit compiler core
- deterministic behavior
- clear grammar
- strong static typing
- structural interfaces
- explicit composition
- no implicit conversions where they would hide meaning
- excellent diagnostics
- reusable libraries
- canonical formatting
- source/spec/conformance discipline
- strong tooling around the compiler

The Physics Compiler is **not** a compiler implementation of the broader Reality-First research methodology.

Human conceptual frameworks remain useful to humans, but they are not automatically compiler ontology.

---

# 3. Human Mental Models vs Compiler Concerns

The following remain conceptual/research tools rather than compiler primitives:

- Reality-First
- No Sacred Cows
- Dancer / Dance / Stage
- broader research constitutions
- paper-review methodology
- evidence promotion
- adversarial scientific review
- research obligations
- external context graphs
- EBP/Solvent workflow
- theory ranking
- scientific verdicts

They may appear in documentation, examples, or external tooling.

The compiler itself works with:

- grammar
- symbols
- types
- interfaces
- operations
- dimensions
- index structures
- domains and codomains
- composition
- provenance markers where formally represented
- symbolic transformations
- diagnostics
- IR

The compiler should not encode a philosophical ontology merely because it is useful for human reasoning.

---

# 4. Governing Architectural Principle

The compiler must distinguish four things:

```text
mathematical structure
        ≠
symbolic representation
        ≠
physical structure
        ≠
physical interpretation
```

These may be related, but the relationships must be represented explicitly.

The central category-error rule is:

> **A mathematical object must not silently acquire a physical meaning that its formal type, interface, domain/codomain, dimensions, index structure, or explicit bridge does not support.**

This is the compiler-native form of the **Math Reality Check (MRC)**.

---

# 5. Math Reality Check (MRC)

## 5.1 Purpose

MRC is a **first-class compiler semantic layer** whose purpose is to prevent category errors between mathematical structures and physical variables/structures.

MRC is specifically concerned with questions such as:

- Is this mathematical object being used as a mathematical object, or as a physical object?
- Does the physical quantity have an explicit mathematical carrier?
- Is a mathematical structure being silently identified with a physical quantity?
- Does a cross-domain transformation have an explicit construction?
- Does the mathematical operation support the intended physical use?
- Are domain/codomain requirements satisfied?
- Are dimensions compatible?
- Are tensor/index constraints preserved?
- Does the declared physical type actually provide the capability required by the operation?

## 5.2 MRC Is Not a Physics-Truth Engine

MRC may determine:

```text
FORMALLY COMPATIBLE
```

or:

```text
CATEGORY ERROR
```

only where the conclusion follows from the language's formal contracts.

It does not determine:

```text
NATURE CONFIRMS THIS
```

or:

```text
THIS THEORY IS TRUE
```

## 5.3 No Universal Jurisdiction Matrix

MRC must **not** be implemented as a universal hardcoded:

```text
MathType × PhysicalRole × Regime
```

matrix.

The compiler should not contain a giant table saying which mathematical structures Nature allows to generate which physical structures.

Instead, compatibility should arise through:

- types
- structural interfaces
- operator signatures
- constructors
- domain/codomain
- explicit mappings
- composition
- dimensions
- tensor/index rules
- declared mathematical structures
- declared physical structures

The compiler enforces the formal contracts that the language defines.

## 5.4 Cross-Domain Use Must Be Explicit

A mathematical object and a physical object are different categories unless the program explicitly constructs a relationship.

For example, this should not silently pass merely because the names look similar:

```phys
N : Integer
t : Time

t = N
```

The compiler requires a declared construction or mapping whose type contract explains how an integer-valued mathematical object supplies a `Time`.

A formally declared construction may instead look conceptually like:

```phys
growth_step : Step
clock : Temporal[Step]
t : Time = clock.duration(growth_step)
```

The compiler does not decide whether the construction models Nature correctly. It checks that the construction is explicitly represented and type-compatible.

---

# 6. MRC Enforcement Levels

MRC has two fundamentally different classes of findings.

## 6.1 Hard Compiler Errors

These are mechanically decidable contradictions in the language.

Examples:

```text
Scalar + Vector
```

when the addition signature requires equal scalar structure.

```text
Function[A,B](C)
```

where `C` is not compatible with `A`.

```text
∇μ Vμ
```

when the required differentiable/manifold structure is absent.

```text
Tμν Vμ
```

when indices cannot be legally contracted.

```text
length + time
```

when addition requires dimension equality.

```text
mathematical PartialOrder
```

used directly where an operation requires a `Temporal` capability, without a declared construction satisfying that capability.

These are compile-time failures.

## 6.2 Heuristic / Suspicion Findings

Some category errors cannot be established mechanically.

Those belong to `phys vet`.

Examples:

- unusual cross-domain representation
- unusual regime crossing
- use of a mathematically legitimate invariant in an unexpected physical context
- reliance on an unproven symbolic derivation
- suspicious semantic naming
- a formally legal construction whose physical motivation is unclear

`phys vet` may report:

```text
warning: unusual mathematical-to-physical mapping
```

but it must not claim:

```text
error: Nature cannot work this way
```

---

# 7. The Fundamental Type Separation

The compiler should represent at least these conceptual layers:

```text
MathType
PhysicsType
SymbolicExpression
Bridge / Mapping
```

A mathematical type describes mathematical structure.

A physics type describes a physical quantity or physical structure.

A bridge explicitly relates the two.

The compiler must not collapse everything into one universal type tag.

For example:

```text
Scalar
Vector
Tensor
Manifold
Group
Operator
```

are mathematical structures.

Meanwhile:

```text
Time
Energy
Field
State
Observable
Spacetime
Particle
System
```

are physical structures.

A physical quantity normally has a mathematical carrier and physical semantics.

For example, a symbolic physical quantity may be modeled conceptually as:

```text
Quantity[
    value-type,
    dimension,
    physical-role
]
```

without reducing its physical identity to its mathematical carrier.

Thus:

```text
Time
```

is not merely:

```text
Scalar
```

and:

```text
Energy
```

is not merely:

```text
Dimension[M L² T⁻²]
```

The latter describes part of its mathematical structure, not the entirety of its physical meaning.

---

# 8. Mathematical Grammar

The compiler must provide a real expression language.

## 8.1 Core Expression Nodes

At minimum:

```text
Symbol
Literal
Negation
Sum
Product
Quotient
Power
FunctionApplication
Composition
Equality
Approximation
Inequality
Tuple
SetExpression
IndexExpression
TensorExpression
OperatorApplication
Derivative
Integral
Limit
Variation
Substitution
Contraction
Trace
Determinant
Transpose
Adjoint
Commutator
Anticommutator
```

The exact AST naming may evolve, but the semantic capabilities are required by the language.

## 8.2 Expressions Are Symbolic

A source such as:

```phys
H = p^2 / (2*m) + V(x)
```

creates a symbolic expression tree.

The compiler does not evaluate `H`.

It understands:

- `H` is a symbol
- `p`, `m`, `x` are referenced symbols
- `/` creates a quotient
- `p^2` creates a power
- `V(x)` is a function application
- the expression is symbolically structured
- dimensions and types may be checked

---

# 9. Mathematical Structures

The standard library should support, as appropriate:

## Algebra

- Scalar
- Vector
- Covector
- Matrix
- Tensor
- Function
- Functional
- Operator
- Group
- Algebra
- Representation
- Invariant
- Spectrum

## Geometry

- Manifold
- TangentSpace
- CotangentSpace
- Metric
- Connection
- Curvature
- DifferentialForm
- Bundle

## Topology

- TopologicalSpace
- Homotopy-related structures
- WindingNumber
- CharacteristicClass
- Defect-like structures

## Analysis

- Measure
- Distribution
- Derivative
- Integral
- Limit
- Variation
- Expectation
- Correlation

Not every item becomes a keyword or compiler built-in.

The compiler provides generic mechanisms; the standard library supplies reusable mathematical structures.

---

# 10. Physics Grammar

The standard library should support theoretical-physics constructs such as:

- System
- Subsystem
- State
- Field
- Particle
- Observable
- Measurement
- Configuration
- Action
- Lagrangian
- Hamiltonian
- Constraint
- Symmetry
- ConservationLaw
- Charge
- Current
- Energy
- Momentum
- AngularMomentum
- Spacetime
- Coupling
- InitialCondition
- BoundaryCondition
- Vacuum
- GroundState
- ExcitedState
- Spectrum
- Scattering
- Decay
- Phase
- OrderParameter
- Defect
- Soliton
- GaugeStructure

These should generally be library-level constructs rather than a closed compiler ontology.

---

# 11. Domains and Codomains

Functions, operators, mappings, and transformations must expose explicit domains and codomains where relevant.

Conceptually:

```text
Function[A,B]
Operator[A,B]
Map[A,B]
```

The compiler uses these signatures for mechanical compatibility checking.

Example:

```text
f : Vector → Scalar
```

allows:

```text
f(v)
```

when `v : Vector`.

It rejects:

```text
f(Metric)
```

when `Metric` is not a valid input type.

This is one of the primary mechanisms for preventing category errors.

---

# 12. Structural Interfaces / Capabilities

The compiler adopts Go-style structural interfaces.

A capability says what an object formally supports.

Conceptual examples:

```go
type Differentiable interface {
    Derivative(...)
}

type Measurable interface {
    Measure(...)
}

type Temporal interface {
    Duration(...) Interval
}

type Observable interface {
    Spectrum(...) Set
}

type Dynamical interface {
    Flow(...) Transform
}
```

These are **formal language contracts**, not metaphysical claims.

A type satisfies an interface because its defined structure satisfies the required operations.

Interfaces are preferable to large hardcoded semantic tables.

---

# 13. Composition

Composition is a first-class language principle.

The compiler should support structures that acquire additional mathematical capabilities through composition.

Conceptual examples:

```text
VectorSpace + InnerProduct
    → InnerProductSpace
```

```text
InnerProductSpace + Completeness
    → HilbertSpace
```

```text
Manifold + Metric
    → MetricStructure
```

```text
Metric + Connection
    → CurvatureStructure
```

```text
Field + Spacetime
    → FieldOnSpacetime
```

Composition should normally be represented through types, generic structures, constructors, and interfaces rather than a central semantic matrix.

---

# 14. Generics / Parametric Structure

The language should support generic structures where they make symbolic mathematics cleaner.

Potential examples:

```text
Tensor[Rank, Space]
Field[T]
State[S]
Operator[A,B]
Function[A,B]
```

Generics must remain simple enough to preserve Go-like readability.

They are intended for structural abstraction, not for creating a meta-language for arbitrary theorem proving.

---

# 15. Tensor and Index System

Tensor/index syntax is core to theoretical physics.

The language must represent:

- tensor rank
- covariant indices
- contravariant indices
- index labels
- free indices
- dummy/contracted indices
- tensor products
- contractions
- index renaming
- raising indices
- lowering indices
- symmetry properties
- associated spaces/manifolds

The compiler must distinguish:

```text
Tμν
```

from:

```text
Tμ
```

and must structurally understand contraction.

For:

```text
∇μ Vμ
```

the compiler should be able to reason about:

- derivative operator
- index
- variance
- contraction
- tensor rank
- underlying space/manifold

Index notation must not be stored merely as opaque text.

---

# 16. Dimensions

Dimensions are compile-time symbolic structure.

A dimension is represented as an exponent vector over declared base dimensions.

Conceptually:

```text
Length       L
Mass         M
Time         T
Charge       Q
Temperature  Θ
```

Examples:

```text
velocity = L T⁻¹
energy   = M L² T⁻²
```

The compiler should reject:

```text
length + time
```

when the operation requires dimension-compatible operands.

Dimension checking is distinct from numerical unit conversion.

The compiler does not need a numerical unit-conversion engine.

---

# 17. Symbolic Calculus

The compiler language must represent:

- derivative
- partial derivative
- covariant derivative
- integral
- limit
- variation
- functional derivative where supported

Examples:

```text
∂μ Aν
```

```text
∫ L d⁴x
```

```text
lim[x→0] f(x)
```

```text
δS
```

These are symbolic operators.

Any transformation rules must preserve symbolic structure and operate on the expression representation.

---

# 18. Operator Algebra

The language must support operators that are not assumed to commute.

It should represent constructs such as:

```text
[A,B]
```

```text
AB
```

```text
A ∘ B
```

```text
A†
```

```text
Tr(A)
```

The type system must distinguish ordinary scalar algebra from operator algebra.

The language must not assume:

```text
AB = BA
```

unless an explicit algebraic rule establishes commutativity.

---

# 19. Symbolic Transformation Engine

The compiler should support a modular symbolic transformation subsystem.

Core candidates:

- substitution
- simplification
- normalization
- collection
- factorization
- expansion
- differentiation
- algebraic rewriting
- index contraction
- index renaming
- raising/lowering indices
- operator composition
- variation
- canonicalization

Transformations should operate on typed symbolic representations.

They should not require numerical execution.

The transformation system should be rule-driven and extensible.

It should not become a general theorem prover.

---

# 20. Equations and Relations

The language should distinguish:

```text
definition
equation
identity
constraint
approximation
assumption
relation
```

For example:

```phys
define F := ma
```

is structurally different from:

```phys
postulate F = ma
```

and from:

```phys
derive E = mc^2
```

The exact syntax is to be finalized by the language specification, but these semantic distinctions must exist.

---

# 21. Physical Variables

A physical variable must carry explicit structural information.

A bare mathematical scalar should not silently become a physical quantity.

For example:

```text
x : Scalar
```

does not automatically mean:

```text
x : Position
```

A position requires a declared physical or geometric structure.

Similarly:

```text
n : Integer
```

does not automatically mean:

```text
n : Time
```

or:

```text
n : QuantumNumber
```

or:

```text
n : ParticleCount
```

A physical interpretation requires an explicit physical type or construction.

This rule is central to MRC.

---

# 22. Mathematical-to-Physical Bridges

Cross-domain relationships must be represented explicitly.

Possible constructs include:

```text
representation
identification
interpretation
construction
mapping
embedding
measurement
observable
```

These relationships must be typed.

For example:

```text
Representation[MathematicalStructure, PhysicalStructure]
```

is different from:

```text
Identity[MathematicalStructure, PhysicalStructure]
```

and both are different from merely using one as notation for the other.

The compiler checks the formal contract of the bridge.

It does not determine whether the bridge is physically true.

---

# 23. Representation ≠ Identity

The language must preserve the distinction:

```text
Mathematical representation
```

versus:

```text
Physical identification
```

For example:

```text
HilbertSpace represents StateSpace
```

is not equivalent to:

```text
HilbertSpace is StateSpace
```

Likewise:

```text
MetricTensor represents Geometry
```

is not automatically:

```text
MetricTensor is physical spacetime itself
```

The compiler requires these stronger relationships to be explicitly declared.

MRC must prevent accidental collapse of these categories.

---

# 24. Discrete ≠ Arithmetic

A central MRC rule from the research material is preserved as a conceptual safeguard, but expressed through formal language structure.

An integer-valued quantity is not automatically an arithmetic structure.

Examples:

- winding number
- Chern number
- mode number
- representation label
- discrete spectrum index

may arise naturally from topology, symmetry, spectral structure, or other mathematical constructions.

The compiler should therefore not infer:

```text
Integer → ArithmeticPhysics
```

merely from the existence of an integer.

Likewise:

```text
discrete → arithmetic
```

must not be an implicit semantic rule.

Arithmetic-specific structures must be explicitly represented by their mathematical definitions and interfaces.

---

# 25. Continuous ≠ Discrete

The compiler should preserve mathematical regime distinctions where they are formally required.

Examples:

```text
Continuous
Discrete
Finite
Infinite
Compact
Noncompact
Differentiable
Non-differentiable
```

These are mathematical properties where formally meaningful.

A transformation such as:

```text
discretize(ContinuousStructure)
```

must be represented explicitly.

Likewise:

```text
continuum_limit(DiscreteStructure)
```

must be a declared symbolic construction.

The compiler does not decide whether a particular physical theory's continuum limit succeeds empirically.

---

# 26. Provenance

Provenance may be represented structurally.

Conceptual forms:

```text
Assumed[T]
Derived[T]
Proved[T]
```

These should be treated as formal type information, not a scientific truth-ranking system.

For example:

```text
Derived[EnergyExpression]
```

means the compiler recognizes the object as produced by a declared derivational construct.

It does not mean:

```text
Nature has proven the energy expression.
```

`Proved[T]` must refer to an actual recognized formal proof/certificate mechanism if implemented.

The compiler must not allow an arbitrary annotation to counterfeit a proof certificate.

---

# 27. Dependencies and Derivations

The compiler should support explicit derivational dependencies.

The Go architecture provides the model:

- source declarations establish names
- imports establish dependencies
- declaration-before-use avoids hidden circularity
- derivational modules should have deterministic dependency structure

However:

> **The Physics Compiler must not literally equate the physical derivation graph with the Go package import graph.**

Go's module/package architecture is an engineering inspiration.

Physical couplings and self-consistency relationships can be represented as symbolic values or structures without turning them into package import cycles.

---

# 28. Coupling vs Derivation

The language must distinguish:

```text
A is derived from B
```

from:

```text
A is coupled to B
```

and from:

```text
A represents B
```

and:

```text
A constrains B
```

A system may contain mutual physical coupling without requiring derivational recursion.

This is important for structures such as matter/geometry or interacting fields.

---

# 29. Regimes

Regime information may exist as explicit annotations where useful:

```text
continuum
discrete
classical
quantum
UV
IR
equilibrium
nonequilibrium
```

But regimes are not a universal truth classifier.

Compile-time rules should enforce only formally declared compatibility constraints.

Unusual regime usage may be reported by `phys vet`.

---

# 30. Compiler Pipeline

The canonical pipeline is:

```text
Source
  ↓
Scanner
  ↓
Parser
  ↓
AST
  ↓
Name Resolution
  ↓
Type System
  ↓
Interface Satisfaction
  ↓
Expression Typing
  ↓
Tensor / Index Checking
  ↓
Dimension Checking
  ↓
MRC Formal Checks
  ↓
Composition
  ↓
Symbolic Semantic Analysis
  ↓
Typed Physics IR
  ↓
Formatting / Transformation / Inspection / Documentation
```

There is no numerical execution stage.

---

# 31. Scanner

The scanner is responsible only for lexical structure.

It recognizes:

- identifiers
- keywords
- literals
- operators
- delimiters
- index syntax
- unicode mathematical symbols where supported
- comments
- annotations

Physics meaning belongs later in the pipeline.

---

# 32. Parser

The parser converts valid source into AST.

The parser should not embed physical truth rules.

It should recognize the formal language.

A recursive-descent parser is a reasonable initial implementation strategy.

---

# 33. AST

The AST must be designed around the language rather than around Program A's old artifact schema.

It should support:

- declarations
- symbols
- types
- expressions
- indexed expressions
- tensors
- equations
- definitions
- assumptions
- symbolic operations
- derivations
- mappings
- bridges
- interfaces
- modules/packages

The AST should preserve source locations for diagnostics.

---

# 34. Name Resolution

Name resolution must be deterministic.

It should handle:

- local declarations
- package/module imports
- qualified identifiers
- aliases
- generic parameters
- scope
- shadowing rules
- declaration-before-use where adopted by the language

Unknown symbols are compiler errors.

---

# 35. Type Checker

The type checker is central.

It must check:

- type identity
- assignability
- generic instantiation
- function arguments
- operator arguments
- function domains/codomains
- interface satisfaction
- tensor structure
- index compatibility
- dimension compatibility
- physical/mathematical category boundaries

The type system is the primary enforcement mechanism for category correctness.

---

# 36. MRC Compiler Pass

MRC operates over the typed representation.

It should verify:

1. mathematical and physical categories are explicitly represented;
2. cross-domain relationships have explicit bridges;
3. an operation's required capabilities are actually implemented;
4. a mathematical carrier is not silently substituted for a physical quantity;
5. a physical quantity is not reduced to an unrelated mathematical structure;
6. dimensions remain consistent;
7. tensor/index structure remains valid;
8. domain/codomain contracts are respected;
9. representations are not silently treated as identities;
10. symbolic constructions do not acquire unsupported physical roles by naming alone.

MRC should use the strongest available formal contract.

When a problem is mechanically decidable, it is a compiler error.

When it is not mechanically decidable, it belongs in `phys vet` or human interpretation.

---

# 37. Physics IR

The Physics IR is a **typed symbolic IR**.

It is not:

- an execution IR
- a numerical IR
- an LLVM-like machine IR
- a simulation graph
- a runtime object model

It should represent canonical symbolic structures such as:

```text
Symbol
Constant
Sum
Product
Power
FunctionApplication
Derivative
Integral
Limit
Variation
Tensor
Index
Contraction
OperatorApplication
Equation
Constraint
Mapping
Bridge
Declaration
Type
Capability
```

The IR should preserve sufficient structure for symbolic transformation and inspection.

It must not simply be a serialized copy of the AST.

---

# 38. IR Invariants

The IR should provide deterministic invariants such as:

- resolved names
- normalized type references
- explicit operator signatures
- explicit tensor/index relationships
- normalized dimensions
- explicit physical/mathematical bridges
- canonical symbolic structure
- provenance markers where applicable

Equivalent source forms may lower to equivalent canonical IR.

---

# 39. Standard Library Philosophy

The compiler core should remain small.

The standard library supplies mathematical and physics vocabulary.

Potential package families:

```text
std/algebra
std/linear
std/calculus
std/tensor
std/geometry
std/topology
std/measure
std/operators
std/analysis
std/relativity
std/qm
std/statmech
std/field
std/dynamics
```

These are candidates, not mandatory package names.

The standard library is documentation and reusable formal structure.

It is not a hidden authority on physical truth.

---

# 40. Math Branch Metadata

The earlier Reality-First material identified broad mathematical branches and their natural areas of use, including:

- Algebra
- Arithmetic
- Calculus
- Chaos
- Complexity
- Degrees of Freedom
- Duality
- Emergence
- Entropy
- Geometry
- Invariants
- Langlands
- Spectra
- Symmetry
- Topology

That material is useful as a design/reference layer.

It must not become a universal compiler jurisdiction table.

Where useful, branch information may appear in:

- package documentation
- type documentation
- `phys doc`
- `phys vet` heuristics
- educational examples

Formal compiler behavior must come from actual type/interface/operator contracts.

---

# 41. Toolchain

The toolchain should follow Go's model.

## `physc`

Primary compiler.

Responsibilities:

- parse
- resolve
- type-check
- MRC-check
- build typed symbolic IR
- emit canonical artifacts
- report diagnostics

## `phys fmt`

Canonical source formatter.

It should provide deterministic formatting for:

- declarations
- expressions
- tensor/index notation
- equations
- symbolic operators

## `phys vet`

Heuristic static analysis.

Possible findings:

- suspicious category crossings
- unusual mathematical-to-physical mapping
- unsupported regime use
- unproven derived structures
- suspicious unused assumptions
- potentially accidental identity claims
- symbolic patterns likely to deserve human inspection

`phys vet` is advisory.

## `phys doc`

Generate documentation from declarations and library metadata.

## `phys types`

Expose reusable type-checking functionality for editors and other tools.

## `phys test`

Support tests of:

- parser behavior
- type rules
- interface satisfaction
- symbolic transformations
- MRC behavior
- tensor/index rules
- dimensional rules
- standard-library contracts

---

# 42. Diagnostics

Diagnostics should be precise and Go-like.

A diagnostic should include, where possible:

- source position
- error/warning code
- concise message
- relevant types
- relevant symbols
- violated contract
- useful context
- optional suggested correction

Example:

```text
./time.ph:17:12: E0214:
cannot use causal_order as Temporal

  causal_order has type PartialOrder
  Temporal requires Duration() Interval

  no declared construction from PartialOrder to Temporal
```

The message should explain the formal mismatch rather than make philosophical claims.

---

# 43. Formatting and Canonicalization

Source and IR should be deterministic.

Canonicalization should normalize, where formally justified:

- equivalent syntax forms
- expression ordering where commutativity has been established
- index naming where safe
- type references
- dimensions
- declaration ordering when permitted

The compiler must not apply algebraic identities that are not valid for the relevant algebraic structure.

For example, noncommuting operators must not be reordered as though they were scalar multiplication.

---

# 44. Testing Strategy

The specification must be paired with a conformance corpus.

Test classes should include:

## Lexer tests

- identifiers
- unicode symbols
- operators
- indexes
- literals

## Parser tests

- declarations
- expressions
- equations
- derivatives
- integrals
- tensor notation
- mappings

## Type tests

- valid operations
- invalid operations
- generic types
- interface satisfaction
- domain/codomain failures

## MRC tests

- math-only valid use
- physical-only valid use
- explicit math→physics bridge
- missing bridge
- illegal identity
- scalar wrongly used as physical variable
- integer/discrete/arithmetic confusion
- regime mismatch
- mathematically valid but physically role-incompatible construction

## Tensor tests

- valid contraction
- invalid contraction
- free-index mismatch
- variance errors
- raising/lowering

## Dimension tests

- compatible addition
- invalid addition
- correct multiplication
- correct division
- incomplete dimension metadata

## Symbolic transformation tests

- substitution
- simplification
- derivative
- variation
- index renaming
- contraction
- operator algebra
- canonicalization

## Determinism tests

Repeated compilation of identical source must yield identical semantic output.

---

# 45. Canonical Stress Tests

The initial implementation must be capable of representing at least the following classes of expressions.

## Classical mechanics

```phys
L = 1/2 * m * v^2 - V(x)
```

## Field theory

```phys
F[μν] = ∂[μ] A[ν] - ∂[ν] A[μ]
```

## Variational principle

```phys
δS = 0
```

## Quantum operator algebra

```phys
[p, q] = -i * ℏ
```

## Relativity

```phys
G[μν] + Λ * g[μν] = 8π * G * T[μν]
```

## Statistical mechanics

```phys
Z = Tr(exp(-β * H))
```

The compiler must treat these as symbolic theoretical constructs.

---

# 46. Canonical MRC Stress Tests

The compiler should also exercise category boundaries.

## 46.1 Mathematical integer ≠ physical time

```phys
n : Integer
t : Time

t = n
```

This requires an explicit and valid construction.

## 46.2 Partial order ≠ temporal duration

```phys
C : PartialOrder
t : Time

t = C
```

This is invalid unless the program provides an explicit construction satisfying the required temporal contract.

## 46.3 Topological invariant ≠ dynamical generator

A topological invariant may constrain a system if an explicit mathematical construction supports that relationship.

It must not automatically become a dynamical generator merely because both participate in the same theory.

## 46.4 Representation ≠ identity

A declaration of representation must not silently satisfy an identity relationship.

---

# 47. What the Compiler Guarantees

For a fixed language specification and implementation, the compiler guarantees only what it formally checks.

The compiler can guarantee:

- syntactic well-formedness
- name resolution
- static type correctness
- interface satisfaction
- domain/codomain compatibility
- dimension consistency
- tensor/index consistency
- formally defined symbolic transformation correctness
- explicit math/physics category boundaries
- deterministic compilation/canonicalization

It cannot guarantee:

- truth of physics
- empirical validity
- experimental confirmation
- correctness of an unformalized physical interpretation
- completeness of the mathematical theory
- that a formally expressible model corresponds to Nature

---

# 48. Explicitly Rejected Architectures

The following are rejected for v0:

## 48.1 Numerical Code Generation

No:

- LLVM backend
- C backend
- CUDA
- OpenCL
- GPU kernels
- numerical runtime

## 48.2 Simulation Engine

No internal:

- ODE solver
- PDE solver
- Monte Carlo engine
- lattice engine
- numerical integrator
- optimizer

## 48.3 Scientific Adjudication Engine

No:

- truth scores
- scientific rankings
- theory rankings
- evidence debt
- paper litigation
- adversarial review
- empirical promotion states

## 48.4 Giant Physics Ontology

Do not encode every physical noun as a compiler keyword.

## 48.5 Giant Math Enumeration

Do not make the compiler dependent on a closed list of every possible mathematical structure.

## 48.6 Universal Jurisdiction Matrix

Do not hardcode a matrix deciding all permitted mathematical-to-physical relationships.

## 48.7 Dancer/Dance/Stage Types

Do not use those concepts as compiler ontology.

## 48.8 General Theorem Prover

No theorem proving inside compiler core.

## 48.9 General-Purpose CAS

The symbolic subsystem should be language-focused and modular, not an attempt to reproduce an unrestricted computer algebra system.

---

# 49. Program A Relationship

Program A is an earlier prototype for structured theoretical-physics semantic checking.

The Physics Compiler is not:

```text
Program A + parser
```

It is a new language architecture.

Program A may supply:

- ideas
- test cases
- vocabulary
- experience
- useful semantic constraints
- historical implementation knowledge

but the Physics Compiler must have its own:

- grammar
- expression language
- AST
- type system
- interfaces
- composition model
- tensor/index system
- symbolic engine
- Physics IR
- standard library
- toolchain

Program A's limitations must not become architectural constraints on the new compiler.

---

# 50. External Systems Boundary

The Physics Compiler may export symbolic representations for consumption by external systems.

Possible future consumers include:

- computer algebra systems
- theorem provers
- numerical engines
- simulation tools
- external mathematical software

These systems are external.

The Physics Compiler does not absorb their responsibilities.

The clean boundary is:

```text
                  Physics Compiler
                         │
                         │ typed symbolic IR
                         ↓
              external symbolic/computational
                       systems
```

---

# 51. The Compiler's Fundamental Contract

The defining contract is:

> **The compiler represents human mathematical/physical intent in a formal symbolic language and enforces the formal relationships expressed by that language.**

The compiler does not interpret Nature.

It does not decide what a symbol "really is" in the physical world.

It ensures that the program does not silently violate the formal distinctions the language has declared.

---

# 52. The MRC Principle in One Sentence

The entire Math Reality Check should reduce to this compiler principle:

> **No mathematical structure may silently masquerade as a physical structure, and no physical variable may silently collapse into an unrelated mathematical structure; any such relationship must be explicit and formally typed.**

This is the core category-error protection.

---

# 53. Definitive Architectural Summary

The Physics Compiler is:

```text
               SYMBOLIC THEORETICAL PHYSICS
                          │
                          ↓
                    .phys language
                          │
                          ↓
                     Scanner
                          │
                          ↓
                      Parser
                          │
                          ↓
                        AST
                          │
                          ↓
                   Name Resolution
                          │
                          ↓
                     Type System
                          │
             ┌────────────┼────────────┐
             ↓            ↓            ↓
         Interfaces    Dimensions   Tensor/Index
             │            │            │
             └────────────┼────────────┘
                          ↓
                         MRC
                          │
                          ↓
                     Composition
                          │
                          ↓
              Symbolic Semantic Analysis
                          │
                          ↓
                  Typed Physics IR
                          │
             ┌────────────┼─────────────┐
             ↓            ↓             ↓
         phys fmt      phys doc      phys vet
             │
             ↓
     symbolic transformation/
       inspection tooling
```

No physical-value execution exists in this architecture.

No scientific adjudication exists in this architecture.

No numerical runtime exists in this architecture.

The heart of the system is:

```text
Symbolic Expression
+
Mathematical Type System
+
Physics Type System
+
Interfaces
+
Composition
+
Tensor/Index Semantics
+
Dimensions
+
Math Reality Check
+
Typed Symbolic IR
```

That is the definitive v0.1 foundation.

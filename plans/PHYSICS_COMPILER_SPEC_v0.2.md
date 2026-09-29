# Physics Compiler — Definitive Specification v0.2

**Status:** FROZEN REFERENCE  
**Supersedes:** `PHYSICS_COMPILER_SPEC_v0.1.md`  
**Scope:** Symbolic theoretical physics  
**Implementation language:** Go  
**Source language:** `.phys`

---

## 0. Executive Definition

The **Physics Compiler** is a programming language and compiler for enabling AI agents to represent, manipulate, reason about, and derive **theoretical physics symbolically**.

The compiler is not primarily a notation parser.

It is a **semantic representation language** for physics.

A human physics paper may write:

```text
G_{μν} + Λ g_{μν} = 8πG T_{μν}
```

but the Physics Compiler should represent the same content as a structured symbolic construction such as:

```text
Equality(
    Add(
        EinsteinTensor(spacetime, metric),
        Scale(cosmologicalConstant, MetricTensor(spacetime, metric))
    ),
    Scale(
        Multiply(Eight, NewtonConstant),
        StressEnergyTensor(matter)
    )
)
```

The compiler therefore works from **named variables, typed objects, constructors, operations, and relationships**, not from human mathematical typography.

The important distinction is:

> **Mathematical notation is a representation used by physicists; the Physics Compiler is a formal semantic representation used by AI agents and tools.**

The compiler can emit conventional mathematical notation for humans, but conventional mathematical notation is not the fundamental source representation.

---

# 1. Core Mission

The Physics Compiler exists to give AI agents a formal language in which they can:

1. read or receive physics concepts;
2. represent mathematical and physical objects as typed variables;
3. compose those objects;
4. apply allowed mathematical and physical operations;
5. symbolically transform expressions;
6. formulate hypotheses;
7. reason through explicit symbolic transformations;
8. perform symbolic calculations;
9. derive new equations from existing symbolic structures;
10. expose the complete derivation and representation for inspection by humans or other AI agents;
11. support a separate adversarial AI that can inspect the same formal representation.

The compiler is therefore a **symbolic reasoning substrate for theoretical physics**.

It is not merely a syntax checker.

---

# 2. The Central Representation Model

The fundamental object in the system is not a string of mathematical notation.

It is a typed symbolic structure.

For example:

```text
Equality
 ├── Add
 │    ├── EinsteinTensor
 │    └── Scale(Λ, MetricTensor)
 └── Scale(8G, StressEnergyTensor)
```

is the canonical semantic model.

Likewise, a source declaration might conceptually be:

```phys
var metric Metric[spacetime]
var cosmologicalConstant Scalar
var stressEnergy StressEnergyTensor[matter]

einsteinTensor := EinsteinTensor(spacetime, metric)

geometry := Scale(cosmologicalConstant, metric)

rhs := Scale(
    Multiply(Eight, newtonConstant),
    stressEnergy
)

equation := Equality(
    Add(einsteinTensor, geometry),
    rhs,
)
```

The exact surface syntax is to be finalized by the language specification, but the **semantic shape** above is normative.

---

# 3. Source Language Philosophy

The source language is intentionally **Go-like** in its programming model.

AI agents should primarily work with:

- identifiers
- variable declarations
- assignments
- functions
- constructors
- method calls
- explicit operations
- explicit types
- explicit composition

rather than LaTeX-like notation.

### Preferred

```phys
metric := Metric(spacetime)
field := ScalarField(spacetime)
energy := Energy(field)

gradient := Gradient(energy)
```

### Not the primary source form

```text
g_{μν}
∇_μ φ
G_{μν} + Λg_{μν} = 8πGT_{μν}
```

The latter may be emitted by formatting/documentation tools, but it should not be the core language representation.

This is deliberate:

> **The language should expose semantic objects and allowed operations directly so AI agents can reason over structure rather than infer structure from notation.**

---

# 4. Variables

Variables are first-class.

They should behave broadly like variables in Go.

Conceptually:

```phys
var metric Metric
var field ScalarField
var action Action
var state State
```

and:

```phys
metric := Metric(spacetime)
field := ScalarField(spacetime)
action := Action(field, metric)
```

Variables have:

- names
- static types
- scope
- source position
- optional provenance
- symbolic value/reference

A variable does not contain a runtime numerical value.

It refers to a symbolic object.

---

# 5. Assignment

Assignment creates symbolic bindings.

For example:

```phys
metric := Metric(spacetime)
curvature := Curvature(metric)
action := EinsteinHilbertAction(metric)
```

The compiler records the symbolic construction.

Assignment never implies numerical evaluation.

---

# 6. Constants and Exact Symbolic Values

The language supports exact symbolic constants and literals where mathematically appropriate.

Examples:

```phys
pi
hbar
c
e
G
sqrt(2)
```

These are symbolic values.

The compiler may simplify exact symbolic arithmetic, but it does not substitute measured numerical values.

For example:

```text
Multiply(2, 3)
```

may normalize to:

```text
6
```

while:

```text
Multiply(NewtonConstant, Mass)
```

remains symbolic.

---

# 7. Mathematical Objects and Physical Objects

The compiler must maintain a formal distinction between:

### Mathematical objects

Examples:

```text
Scalar
Vector
Covector
Tensor
Function
Functional
Operator
Set
Relation
Group
Algebra
Manifold
Metric
Connection
Curvature
Topology
Measure
Distribution
Spectrum
```

### Physical objects

Examples:

```text
System
Subsystem
State
Field
Particle
Observable
Measurement
Spacetime
Energy
Momentum
Charge
Current
Action
Lagrangian
Hamiltonian
Constraint
Symmetry
ConservationLaw
Coupling
InitialCondition
BoundaryCondition
Vacuum
GroundState
ExcitedState
```

These categories must not collapse into one generic "physics object" type.

However, mathematical and physical structures may be **composed explicitly**.

---

# 8. Mathematical Objects Are Not Automatically Physical Objects

A mathematical object does not become a physical quantity merely because its type can be represented numerically.

Examples:

```text
Integer
Scalar
Vector
Group
Topology
Spectrum
```

must not silently become:

```text
Time
Energy
Particle
Spacetime
Dynamics
```

without an explicit formal construction.

For example:

```phys
n := IntegerSymbol("n")
time := Time(n)
```

must only compile when `Time(...)` is a declared construction whose argument contract accepts the supplied object.

The compiler is enforcing the language's formal contract.

It is not making a claim about Nature.

---

# 9. Math Reality Check (MRC)

MRC is a **compiler-native structural rule**.

Its purpose is:

> **Prevent category errors between mathematical objects and physical objects.**

MRC is not a scientific truth detector.

It is the mechanism that prevents accidental semantic collapse.

Examples of protected distinctions:

```text
mathematical scalar
    ≠
physical time
```

```text
integer
    ≠
physical particle number
```

```text
partial order
    ≠
physical duration
```

```text
topological invariant
    ≠
dynamical generator
```

```text
representation
    ≠
identity
```

```text
mathematical spectrum
    ≠
physical observable
```

unless an explicit typed construction defines the relationship.

---

# 10. MRC Is Based on Formal Contracts

MRC must not be implemented as:

```text
MathType × PhysicalRole × Regime
```

or another universal hardcoded jurisdiction matrix.

Instead, MRC relies on:

- static types
- structural interfaces
- operator signatures
- function domains/codomains
- constructors
- explicit mappings
- explicit representations
- explicit identifications
- explicit interpretations
- dimensions
- tensor/index constraints
- composition
- declared transformations

The compiler should ask:

> **Does the program actually construct the relationship it claims?**

not:

> **Does the compiler's philosophical table say Nature permits it?**

---

# 11. Explicit Math–Physics Relationships

The language needs explicit constructs for relationships such as:

```text
represents
identifies
interprets
maps
embeds
constructs
constrains
couples
```

For example:

```phys
representation := Represents(metric, spacetimeGeometry)
```

is different from:

```phys
identity := Identifies(metric, physicalSpacetime)
```

and from:

```phys
interpretation := Interprets(solution, physicalObservable)
```

The compiler checks the formal contracts of these operations.

---

# 12. Representation Is Not Identity

This distinction is mandatory.

A mathematical structure used to represent a physical object is not automatically identical to that physical object.

For example:

```text
MetricTensor
```

may represent geometric structure in a physical theory without the language silently inferring:

```text
MetricTensor == physical spacetime
```

Likewise:

```text
HilbertSpace
```

may represent the state structure of a theory without automatically becoming:

```text
physical state itself
```

The stronger relationship must be explicitly represented.

---

# 13. Physics Variables as Typed Semantic Objects

A physics variable should carry more than a raw mathematical carrier.

For example, an energy variable may conceptually contain:

```text
physical type:
    Energy

mathematical carrier:
    Scalar

dimension:
    M L² T⁻²

domain/context:
    System or State

operations:
    addition
    scaling
    expectation
    conservation-related transformations
```

This is why:

```text
Energy
```

is not merely:

```text
Scalar
```

and why:

```text
Time
```

is not merely:

```text
Scalar
```

The physical type supplies semantic structure and allowed operations.

---

# 14. Operations Are the Primary Reasoning Mechanism

The language should emphasize **operations over notation**.

Examples:

```text
differentiate
integrate
vary
substitute
simplify
factor
expand
collect
limit
compose
contract
raiseIndex
lowerIndex
transpose
adjoint
commute
anticommute
trace
determinant
project
expect
normalize
transform
map
represent
identify
interpret
derive
constrain
couple
```

The operation set should be compositional and typed.

The AI agent reasons by applying these operations to symbolic objects.

---

# 15. Operations Must Be Generic Across Physics Domains

Operations such as:

```text
differentiate
integrate
substitute
compose
simplify
vary
project
expect
contract
```

must not be restricted to one branch of physics.

Their applicability is determined by the types and interfaces of their arguments.

For example, `differentiate` may operate on:

- scalar functions
- vector-valued functions
- tensor fields
- actions/functionals
- operator-valued expressions
- field expressions

when their declared structures satisfy the required capability.

Similarly, `integrate` may apply to different mathematical structures when a valid domain and measure are supplied.

This keeps the language centered on **general symbolic operations** rather than creating a separate bespoke operation vocabulary for every physics subfield.

---

# 16. Interfaces / Capabilities

The compiler follows Go-style structural interfaces.

Capabilities describe what an object can formally support.

Conceptual examples:

```go
type Differentiable interface {
    Derivative(variable Variable) Expression
}

type Integrable interface {
    Integrate(domain Domain) Expression
}

type Observable interface {
    Spectrum() Spectrum
}

type Operator interface {
    Domain() Type
    Codomain() Type
}

type Tensor interface {
    Rank() int
}
```

The exact interface syntax belongs to the language specification.

The principle is normative:

> **A type may participate in an operation because it satisfies a formal capability contract.**

This is the principal mechanism for broad operation reuse across physics domains.

---

# 17. No Hidden Semantic Meaning From Names

Names alone must not determine semantics.

For example:

```text
time
energy
state
field
space
```

do not become valid physics objects merely because the identifier has a familiar name.

Likewise:

```text
x
n
G
H
S
T
```

do not receive physics meaning from conventional mathematical naming.

Meaning comes from:

- declared types
- constructors
- operations
- interfaces
- relationships

This is critical for AI robustness.

---

# 18. Mathematical Notation Ingestion

The compiler itself does not need to parse arbitrary mathematical typography.

A physics paper may contain:

```text
G_{μν}
```

The AI agent or an upstream ingestion system translates that notation into the Physics Compiler's semantic vocabulary, for example:

```text
EinsteinTensor(spacetime, metric)
```

The compiler then checks and manipulates the semantic representation.

This creates a clean separation:

```text
Physics paper
    ↓
AI paper-understanding / translation
    ↓
Physics Compiler language
    ↓
typed symbolic representation
    ↓
symbolic reasoning/transformation
```

The AI translation layer may use conventional notation heavily.

The compiler does not need to.

---

# 19. Human Mathematical Notation Is an Output Format

The compiler should be capable of rendering a semantic expression into human-oriented representations.

For example:

```text
Equality(
    Add(
        EinsteinTensor(spacetime, metric),
        Scale(cosmologicalConstant, MetricTensor(spacetime, metric))
    ),
    Scale(
        Multiply(Eight, NewtonConstant),
        StressEnergyTensor(matter)
    )
)
```

may be rendered as conventional mathematical notation for:

- papers
- documentation
- debugging
- human review

The notation layer is therefore **presentation**, not core semantics.

---

# 20. Symbolic Expression Model

The internal expression system must represent:

```text
Symbol
Constant
Constructor
Call
Sum
Product
Quotient
Power
FunctionApplication
Composition
Equality
Approximation
Inequality
Relation
Derivative
Integral
Limit
Variation
Substitution
Contraction
TensorExpression
Index
OperatorApplication
Trace
Determinant
Transpose
Adjoint
Commutator
Anticommutator
Expectation
Projection
Transform
```

Not every construct must appear as source syntax.

The compiler's internal representation must be able to represent these structures.

---

# 21. Symbolic Calculations

The compiler must support **symbolic calculation**.

This includes transformations such as:

```text
differentiate(expression, variable)
integrate(expression, domain)
substitute(expression, variable, replacement)
simplify(expression)
expand(expression)
factor(expression)
collect(expression, symbol)
vary(functional, variable)
takeLimit(expression, point)
contract(tensor, indices)
raiseIndex(tensor, metric)
lowerIndex(tensor, metric)
commutator(A, B)
expect(state, observable)
compose(f, g)
```

These operations can produce new symbolic expressions.

For example:

```text
simplify(
    Add(
        Power(x, 2),
        Scale(2, x),
        One
    )
)
```

may produce the canonical symbolic structure equivalent to:

```text
Power(Add(x, One), 2)
```

The compiler is performing symbolic algebra, not numerical evaluation.

---

# 22. Derivation

The language must support explicit symbolic derivations.

A derivation is a sequence of mechanically represented transformations.

Conceptually:

```phys
start := lagrangian(...)
step1 := vary(start, field)
step2 := simplify(step1)
equation := Equality(step2, Zero)
```

or a structured derivation construct.

The compiler can track:

```text
input expression
→ transformation
→ result
```

and retain the derivation chain in IR.

This gives AI agents a machine-readable reasoning trace.

---

# 23. Derivation Does Not Mean Scientific Truth

A successful symbolic derivation means:

> the declared transformation rules produced the result.

It does not mean:

> Nature confirms the result.

The distinction must remain explicit.

---

# 24. Hypothesis Representation

The language must be able to express hypotheses as symbolic programs.

A hypothesis may contain:

```text
postulates
definitions
objects
relations
equations
constraints
derivations
transformations
regimes
representations
identifications
interpretations
```

For example:

```phys
postulate spacetime := Manifold(...)
postulate metric := Metric(spacetime)
postulate field := ScalarField(spacetime)

lagrangian := Lagrangian(field, metric)

equation := EulerLagrange(
    lagrangian,
    field
)
```

The result is a symbolic hypothesis representation.

The compiler checks its formal consistency.

---

# 25. AI Reasoning Model

The intended AI workflow is:

```text
physics paper
      ↓
extract concepts
      ↓
translate notation → semantic objects
      ↓
declare variables
      ↓
compose objects
      ↓
apply symbolic operations
      ↓
produce equations
      ↓
transform equations
      ↓
derive further equations
      ↓
formulate hypothesis
      ↓
export typed symbolic representation
```

A second AI agent should be able to consume that same representation and inspect the symbolic reasoning.

The compiler itself does not need to contain the adversarial reviewer.

It needs to provide a sufficiently rich and transparent representation for that reviewer.

---

# 26. AI-Friendly Design

The language should optimize for:

- explicitness
- discoverable operations
- typed objects
- composability
- predictable syntax
- inspectable transformations
- precise diagnostics
- canonical representation
- low ambiguity

Avoid requiring an AI agent to infer complex semantics from typography.

Prefer:

```text
Differentiate(field, coordinate)
```

over forcing the agent to construct:

```text
∂φ/∂x
```

Prefer:

```text
Scale(cosmologicalConstant, MetricTensor(...))
```

over:

```text
Λgμν
```

Prefer:

```text
Equality(lhs, rhs)
```

over making an AI parser infer all physical meaning from a handwritten equation string.

---

# 27. Tensor and Index Semantics

Indices must remain explicit in the internal representation even though conventional index notation is not the core source syntax.

The compiler must represent:

- index identity
- free vs bound indices
- covariant vs contravariant variance
- tensor rank
- contraction
- tensor product
- raising/lowering
- index renaming
- symmetry properties
- associated vector spaces/manifolds

For example:

```text
EinsteinTensor(spacetime, metric)
```

may have an internal tensor structure corresponding to a rank-two covariant tensor.

An operation such as:

```text
Contract(A, index1, B, index2)
```

must verify the index compatibility mechanically.

The compiler must not merely store `μ` or `ν` as decorative characters.

---

# 28. Dimensions

Dimensions remain compile-time symbolic structure.

Examples:

```text
Velocity = Length / Time
Energy = Mass * Length² / Time²
```

Dimension checking must work through expression composition.

For example:

```text
Add(Length, Time)
```

should fail when the operation requires equal dimensions.

Dimensions are a formal property.

They do not constitute physical-truth adjudication.

---

# 29. Regimes

Regimes may be explicitly represented where relevant:

```text
classical
quantum
relativistic
nonrelativistic
continuum
discrete
equilibrium
nonequilibrium
UV
IR
```

The compiler enforces only formally declared compatibility rules.

A regime is not a claim that a theory is correct in that regime.

---

# 30. Composition

Composition is one of the central design principles.

The compiler should support:

### Mathematical composition

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

### Physical composition

```text
Field + Spacetime
    → FieldOnSpacetime
```

```text
Matter + Geometry
    → CoupledSystem
```

### Operator composition

```text
Function[A,B] + Function[B,C]
    → Function[A,C]
```

Composition should happen through formal types, constructors, and interfaces.

---

# 31. Generic Structures

Generics should support reusable symbolic structures such as:

```text
Tensor[Rank, Space]
Field[T]
State[S]
Operator[A,B]
Function[A,B]
Observable[S,A]
```

Generics should make mathematical structure explicit without creating a new meta-language.

---

# 32. Algebraic Structure

The compiler must preserve the algebraic properties of objects.

Examples:

- scalar multiplication
- vector addition
- noncommuting multiplication
- linear operators
- associative composition where declared
- commutative operations where declared
- Lie brackets
- anticommutators
- adjoints

The compiler must never reorder an expression merely because it resembles ordinary scalar algebra.

For example:

```text
A * B
```

must remain distinct from:

```text
B * A
```

unless the relevant types/rules establish commutativity.

---

# 33. Quantum Mechanics and Quantum Field Theory

The compiler must be designed so that the same symbolic operation system can support quantum structures.

Important constructs include:

```text
HilbertSpace
StateVector
DensityOperator
Observable
Operator
Hamiltonian
PositionOperator
MomentumOperator
Commutator
Expectation
Projection
Unitary
EvolutionOperator
CreationOperator
AnnihilationOperator
FieldOperator
Vacuum
CorrelationFunction
Propagator
```

These should be composable through general language mechanisms.

Operations such as:

```text
Differentiate
Integrate
Compose
Commute
Expect
Project
Transform
Trace
```

must be reusable wherever the involved types satisfy the necessary contracts.

---

# 34. Cross-Domain Operations in Physics

The language should allow operations to be shared across fields.

For example, the same conceptual:

```text
Differentiate
```

operation may be used for:

- classical mechanics
- field theory
- general relativity
- quantum mechanics
- quantum field theory
- statistical mechanics

with different mathematical structures supplying the appropriate derivative semantics.

Likewise:

```text
Integrate
```

may apply to:

- ordinary functions
- fields
- measures
- path-integral expressions
- expectation constructions

when the relevant type and domain contracts are satisfied.

The architecture must therefore avoid a separate hardcoded "calculus engine for each physics branch."

---

# 35. Standard Library

The standard library supplies reusable mathematical and physical structures.

Candidate packages include:

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
std/qft
std/statmech
std/fields
std/dynamics
```

The compiler core remains small.

The library supplies domain vocabulary and operation implementations/contracts.

---

# 36. MRC and the Standard Library

MRC does not hardcode a universal theory of physics.

Instead, the standard library exposes precise definitions.

For example:

```text
std/qm.Operator
std/qm.Observable
std/qm.HilbertSpace
```

can define the operations and relationships that those objects support.

Likewise:

```text
std/geometry.Manifold
std/geometry.Metric
std/geometry.Connection
```

provide formal structures.

The compiler then checks composition and interface satisfaction.

---

# 37. Arithmetic and Discreteness

The system must preserve the distinction:

```text
integer-valued
```

from:

```text
arithmetic structure
```

and:

```text
discrete
```

from:

```text
arithmetic origin
```

A quantity may be integer-valued because of:

- topology
- symmetry
- counting
- representation labels
- spectral structure

without becoming an arithmetic object.

The compiler must require an explicit mathematical construction before assigning stronger mathematical structure.

---

# 38. Physics Paper Translation

A future AI paper-reading system should be able to convert:

```text
paper notation
```

into:

```text
Physics Compiler semantic objects
```

For example:

```text
paper symbol:
Gμν
```

becomes something such as:

```text
EinsteinTensor(spacetime, metric)
```

and:

```text
paper symbol:
Tμν
```

becomes:

```text
StressEnergyTensor(matter)
```

and:

```text
paper equation:
Gμν + Λgμν = 8πGTμν
```

becomes the typed symbolic structure:

```text
Equality(
    Add(
        EinsteinTensor(spacetime, metric),
        Scale(
            cosmologicalConstant,
            MetricTensor(spacetime, metric)
        )
    ),
    Scale(
        Multiply(Eight, NewtonConstant),
        StressEnergyTensor(matter)
    )
)
```

The compiler then operates on that structure.

---

# 39. Symbolic IR

The Physics IR is a typed symbolic intermediate representation.

It is not an execution IR.

It represents:

```text
symbols
types
constructors
calls
expressions
operators
equations
transformations
tensor structure
dimensions
physical structures
mathematical structures
bridges
provenance
```

The IR must be richer than a serialized AST.

It should contain resolved and normalized semantic information.

---

# 40. IR Example

A conceptual IR for the Einstein equation may look like:

```text
Equality
 ├── Add
 │    ├── Apply
 │    │    └── EinsteinTensor
 │    │         ├── Spacetime
 │    │         └── Metric
 │    └── Scale
 │         ├── CosmologicalConstant
 │         └── Apply
 │              └── MetricTensor
 │                   ├── Spacetime
 │                   └── Metric
 └── Scale
      ├── Multiply
      │    ├── ExactConstant(8)
      │    └── NewtonConstant
      └── Apply
           └── StressEnergyTensor
                └── Matter
```

This structure is what downstream AI tools should inspect and manipulate.

---

# 41. Transformation Trace

Every nontrivial symbolic transformation should be representable as:

```text
input
operation
arguments
output
```

For example:

```text
Derivative(
    expression = Lagrangian(field, metric),
    variable   = field
)
```

produces:

```text
EulerLagrangeExpression(...)
```

The compiler should be able to retain this relationship for inspection.

This allows an AI reviewer to follow:

```text
equation A
    ↓ differentiate
equation B
    ↓ substitute
equation C
    ↓ simplify
equation D
```

---

# 42. Symbolic Derivation and New Insights

The architecture must support chained symbolic operations such that an AI agent can produce new equations.

The compiler does not judge whether the new equation is a correct description of Nature.

It verifies the formal symbolic transformations for which rules exist.

For example:

```text
postulate L
derive dL_dφ := Differentiate(L, φ)
derive eom   := EulerLagrange(L, φ)
derive canon := Simplify(eom)
```

The resulting symbolic chain is machine-readable.

This is the basis for AI-generated theoretical reasoning.

---

# 43. External Adversarial AI

A separate AI system may consume:

- source code
- AST
- typed IR
- transformation trace
- diagnostics
- rendered equations
- dependency information

and perform adversarial analysis.

That analysis is intentionally external.

The Physics Compiler should expose enough structured information that the adversarial agent does not need to reconstruct the meaning from raw text.

---

# 44. Compiler Tools

The toolchain should include:

## `physc`

Compile and type-check `.phys` programs.

Produce:

- diagnostics
- typed symbolic IR
- canonical artifacts
- transformation information

## `phys fmt`

Canonical source formatting.

## `phys vet`

Advisory static analysis for patterns not strictly mechanically decidable.

Potential findings:

- unusual math→physics bridges
- suspicious representation/identity usage
- cross-regime constructions
- unused assumptions
- unsupported or weakly supported transformations
- suspicious symbolic derivation patterns

`phys vet` is not a physics-truth engine.

## `phys doc`

Generate human-readable documentation and conventional mathematical renderings.

## `phys types`

Expose reusable type-system and analysis APIs.

## `phys test`

Test symbolic rules, language semantics, transformations, standard-library contracts, and conformance cases.

---

# 45. Diagnostics

Diagnostics should resemble Go diagnostics in clarity.

Example:

```text
./time.phys:21:10: E0214:
cannot use causalOrder as Temporal

  got:      PartialOrder
  required: Temporal

  reason:
  no declared construction provides the Temporal capability
  for PartialOrder

  MRC:
  mathematical structure cannot be used as a physical
  temporal quantity without an explicit compatible construction
```

The diagnostic should identify the formal reason.

It should not claim:

```text
Nature forbids this.
```

---

# 46. Canonical Formatting

The compiler must have a canonical internal representation.

Equivalent representations should normalize deterministically where the algebra permits.

Examples:

```text
Add(a, b)
```

and:

```text
Add(b, a)
```

may canonicalize identically only when the involved operation is actually commutative.

Noncommutative expressions must preserve order.

Index renaming may be canonicalized where mathematically safe.

---

# 47. Determinism

For identical source and identical library/specification versions:

- parsing is deterministic;
- name resolution is deterministic;
- type checking is deterministic;
- MRC is deterministic;
- symbolic transformations governed by deterministic rules are deterministic;
- canonical IR is deterministic;
- formatting is deterministic.

No timestamps or nondeterministic metadata should affect semantic identity.

---

# 48. No Numerical Runtime

The compiler has no runtime numerical evaluator.

It does not:

- execute variables
- sample values
- simulate systems
- fit parameters
- solve numerical PDEs
- run Monte Carlo
- invoke GPU kernels
- compile to CPU numerical instructions

A symbolic quantity remains symbolic.

---

# 49. No Premature External Backends

Do not make any numerical backend part of the core architecture.

Future systems may consume the symbolic IR:

```text
Physics Compiler
       ↓
Typed Symbolic IR
       ↓
external CAS / theorem prover / numerical system
```

Those systems are consumers, not compiler components.

---

# 50. No General Theorem Prover

The compiler may mechanically apply declared transformation rules.

It does not need to prove arbitrary mathematical theorems.

Formal proof systems may later consume compiler representations externally.

---

# 51. No Research Methodology

The compiler must not embed:

- literature review workflow
- evidence grading
- scientific promotion
- adversarial review methodology
- paper litigation
- research debt
- external context graphs
- theory rankings
- hypothesis truth scores

The compiler supplies formal symbolic representation and analysis.

---

# 52. Program A Relationship

Program A is an earlier theoretical-physics semantic prototype.

The Physics Compiler is not:

```text
Program A + parser
```

The Physics Compiler requires its own:

- language
- source model
- semantic object model
- expression system
- type system
- interfaces
- MRC rules
- symbolic transformation engine
- tensor/index system
- dimensions
- symbolic IR
- standard library
- AI tooling

Program A may provide:

- prior experience
- vocabulary
- candidate tests
- implementation ideas

but Program A's architecture is not the governing architecture.

---

# 53. Definitive Compiler Boundary

The definitive separation is:

```text
┌───────────────────────────────────────────────┐
│ AI PAPER UNDERSTANDING                        │
│ Reads conventional mathematics and physics   │
│ and translates it into semantic objects.      │
└───────────────────────┬───────────────────────┘
                        ↓
┌───────────────────────────────────────────────┐
│ PHYSICS COMPILER LANGUAGE                      │
│ Go-like variables, types, constructors,       │
│ operations, compositions, equations,           │
│ transformations, hypotheses.                  │
└───────────────────────┬───────────────────────┘
                        ↓
┌───────────────────────────────────────────────┐
│ SYMBOLIC COMPILER CORE                         │
│ Scanner → Parser → Resolver → Type Checker   │
│ → Interfaces → MRC → Symbolic Analysis → IR │
└───────────────────────┬───────────────────────┘
                        ↓
┌───────────────────────────────────────────────┐
│ AI REASONING / ADVERSARIAL ANALYSIS            │
│ Consumes source, IR, traces, diagnostics.      │
│ Remains external to compiler semantics.       │
└───────────────────────────────────────────────┘
```

---

# 54. What the Compiler Guarantees

The compiler guarantees the formal properties defined by the language:

- syntactic correctness
- name resolution
- static type correctness
- interface satisfaction
- valid operation signatures
- mathematical structure consistency
- physical structure consistency
- mathematical/physical category separation
- tensor/index consistency
- dimension consistency
- valid declared symbolic transformations
- deterministic symbolic representation

It does not guarantee:

- physical truth
- experimental confirmation
- empirical adequacy
- correctness of an unformalized physical interpretation
- acceptance by the scientific community
- truth of an AI-generated hypothesis

---

# 55. The Defining Principle

The entire architecture can be reduced to this:

> **AI agents should reason about theoretical physics by constructing typed semantic objects and operating on them through explicit symbolic operations.**

Mathematical notation is one way humans write those ideas.

The Physics Compiler's job is to provide the deeper, machine-readable structure underneath that notation.

---

# 56. Final Architectural Test

The v0 language must be able to represent, inspect, transform, and derive symbolic theoretical-physics structures without requiring conventional mathematical notation in the source.

The test is not:

```text
Can the compiler parse Gμν + Λgμν = 8πGTμν?
```

The test is:

```text
Can an AI construct:

metric
einsteinTensor
cosmologicalConstant
stressEnergy
rhs
lhs
equation

using typed variables and explicit operations,
then transform and derive from those objects?
```

The answer must be yes.

The resulting system should make it natural for an AI agent to move from:

```text
physics paper notation
```

to:

```text
semantic physics representation
```

to:

```text
symbolic operations
```

to:

```text
derived equations
```

to:

```text
new theoretical structures
```

while preserving enough formal structure for another AI agent to inspect the reasoning.

---

# 57. V0 Architectural Lock

The following are locked for implementation:

1. **Symbolic theoretical physics only.**
2. **No physical-value execution.**
3. **Source language is Go-like and semantic, not LaTeX-like.**
4. **Variables and assignment are first-class.**
5. **Mathematical and physical objects are distinct typed structures.**
6. **MRC is compiler-native and exists specifically to prevent math↔physics category errors.**
7. **MRC is implemented through formal contracts, not a universal jurisdiction matrix.**
8. **Explicit operations are the primary reasoning mechanism.**
9. **Operations are generic and compositional across physics domains.**
10. **Interfaces/capabilities are central.**
11. **Tensor/index semantics are explicit internally.**
12. **Dimensions are symbolic compile-time structure.**
13. **Symbolic differentiation, integration, variation, substitution, simplification, operator algebra, and related transformations are core capabilities.**
14. **Derivations are explicit symbolic transformation chains.**
15. **The Physics IR is a typed symbolic IR, not an execution IR.**
16. **Human mathematical notation is primarily an input to the AI translation layer and an output of documentation/rendering tools.**
17. **A separate adversarial AI may consume compiler representations but is not part of compiler semantics.**
18. **No numerical runtime, simulation engine, theorem prover, scientific adjudicator, or research-methodology engine is part of v0.**
19. **Program A does not define the Physics Compiler architecture.**
20. **The compiler's purpose is to make theoretical physics machine-representable and machine-operable at the symbolic level.**

---

# 58. One-Sentence Definition

> **The Physics Compiler is a Go-like symbolic programming language for theoretical physics in which AI agents declare mathematical and physical objects as typed variables, compose them, apply formally allowed symbolic operations, derive new equations, and expose the resulting reasoning as a canonical machine-readable representation.**

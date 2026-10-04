# Physics Compiler — Go Library Specification v0.3

**Status:** FROZEN REFERENCE  
**Supersedes:** `PHYSICS_COMPILER_SPEC_v0.2.md`  
**Implementation:** Go library  
**Language:** Go  
**Domain:** Theoretical physics  
**Representation:** Symbolic physics only

---

## 0. The Fundamental Simplification

The previous design was unnecessarily trying to invent a new programming language.

We already have an excellent language:

> **Go.**

The Physics Compiler should therefore **not** create:

- a new `.phys` language
- a scanner
- a parser
- a custom syntax
- a custom type-language
- a custom module system
- a custom AST frontend
- a custom compiler frontend

Instead:

> **The Physics Compiler is a Go library for representing and symbolically manipulating theoretical-physics objects.**

Go already provides:

- variables
- assignment
- functions
- methods
- structs
- interfaces
- generics
- packages
- imports
- dependency management
- static typing
- compiler diagnostics
- tests
- documentation

Use those mechanisms directly.

The project should maximize the power of Go rather than reproduce Go.

---

# 1. Mission

The library exists so that AI agents can express theoretical physics as Go programs.

The intended workflow is:

```text
physics paper
      ↓
AI understands paper notation
      ↓
AI maps notation to Physics Compiler objects
      ↓
AI declares Go variables
      ↓
AI composes physics objects
      ↓
AI applies physics operations
      ↓
symbolic transformation
      ↓
derived physics objects / equations
      ↓
new symbolic insight
      ↓
another AI can inspect the program and derivation
```

The library is therefore a **symbolic reasoning substrate for AI agents working with theoretical physics**.

---

# 2. The Most Important Boundary

The public library should contain **physics objects only**.

There is no public taxonomy such as:

```text
Scalar
Vector
Covector
Tensor
Function
Functional
Set
Group
Algebra
Manifold
Topology
Measure
```

These are mathematical concepts.

They may be referenced in documentation to help an AI translate a paper into the appropriate physics object, but they are **not the public object model of the library**.

The public API instead deals with objects such as:

```text
Mass
Energy
Momentum
Time
Position
Velocity
Field
State
Observable
Particle
System
Spacetime
Metric
Action
Lagrangian
Hamiltonian
Constraint
Symmetry
ConservationLaw
Charge
Current
Temperature
Entropy
PartitionFunction
QuantumState
QuantumObservable
QuantumField
Vacuum
Interaction
Propagator
ScatteringProcess
```

The library represents physics.

Mathematics is the internal machinery and the documentation bridge.

---

# 3. Go Is the Source Language

AI agents write ordinary Go.

Example:

```go
metric := relativity.Metric(spacetime)
matter := relativity.StressEnergy(medium)

einstein := relativity.EinsteinTensor(metric)
cosmological := relativity.CosmologicalConstantTerm(metric, lambda)
equation := relativity.EinsteinEquation(einstein, cosmological, matter)
```

Variables work exactly as Go variables work.

Assignment works exactly as Go assignment works.

Functions work exactly as Go functions work.

There is no `.phys` language.

There is no second compiler frontend.

The ordinary Go compiler performs language-level checking.

---

# 4. What Makes This a "Physics Compiler"

The term describes the role of the library and its tooling, not a new programming language.

The Go compiler handles:

```text
Go syntax
Go types
Go interfaces
Go generics
Go packages
Go dependencies
Go compilation
```

The Physics library handles:

```text
physics objects
physics relationships
physics operations
symbolic transformations
physics-specific type contracts
symbolic derivations
category boundaries
physics-domain composition
```

This separation is intentional.

---

# 5. Symbolic, Not Numerical

The library manipulates symbolic physics.

A variable such as:

```go
mass := mechanics.Mass("m")
```

is symbolic.

A variable such as:

```go
velocity := mechanics.Velocity("v")
```

is symbolic.

An operation such as:

```go
energy := mechanics.KineticEnergy(mass, velocity)
```

produces a symbolic physical object.

It does not evaluate the object's physical value.

The library does not:

- substitute measured values
- run numerical simulations
- solve numerical PDEs
- perform numerical optimization
- execute scientific workloads
- run Monte Carlo
- generate GPU kernels
- provide a numerical physics runtime

The word **calculate** in this project means:

> **perform a formal symbolic transformation and derive another symbolic physics object.**

---

# 6. Exact Symbolic Constants

Exact constants are allowed as ingredients of symbolic physics.

Examples:

```go
pi
hbar
c
G
kB
e
```

Exact symbolic arithmetic may be performed when useful.

For example, symbolic construction may contain the exact factor `8`.

Measured physical values are outside the library's purpose.

---

# 7. No Public Mathematical Object Layer

Do not create public APIs such as:

```go
math.Scalar
math.Vector
math.Tensor
math.Function
math.Group
math.Manifold
math.Measure
```

The project is not a general mathematics library.

A physics object may internally contain whatever symbolic structure is necessary to perform physics operations, but that implementation detail should remain behind the physics API.

For example:

```go
field := quantum.Field(...)
```

is preferable to requiring an AI to first construct:

```go
math.Function[...]
math.Manifold[...]
math.Tensor[...]
```

just to express a field.

---

# 8. Public Physics Object Model

The public API should be organized around physics domains.

## Classical mechanics

Potential objects:

```text
Mass
Position
Velocity
Acceleration
Force
Momentum
AngularMomentum
Energy
KineticEnergy
PotentialEnergy
Lagrangian
Action
EquationOfMotion
Trajectory
Constraint
```

## Classical fields

```text
Field
ScalarField
VectorField
GaugeField
FieldConfiguration
FieldEquation
Current
Charge
EnergyDensity
LagrangianDensity
```

## Relativity

```text
Spacetime
Metric
ProperTime
Worldline
Geodesic
Connection
Curvature
EinsteinTensor
StressEnergy
EinsteinEquation
CosmologicalConstant
```

## Quantum mechanics

```text
QuantumSystem
QuantumState
Observable
Hamiltonian
Evolution
Measurement
Expectation
Transition
Amplitude
DensityState
Projection
QuantumChannel
```

## Quantum field theory

```text
QuantumField
FieldState
Vacuum
Interaction
InteractionHamiltonian
Propagator
Correlation
GreenFunction
ScatteringProcess
SMatrix
CreationProcess
AnnihilationProcess
GaugeField
FermionField
BosonField
```

## Statistical mechanics

```text
Ensemble
Microstate
Macrostate
Temperature
Entropy
FreeEnergy
PartitionFunction
EquationOfState
CorrelationFunction
Response
```

These are examples, not a requirement that every term become a separate public type.

Prefer meaningful physics types with composable construction over a giant taxonomy.

---

# 9. Physics Objects Must Be Composable

Composition is fundamental.

Examples:

```go
spacetime := relativity.Spacetime(...)
metric := relativity.Metric(spacetime)
field := qft.Field(spacetime)
system := qft.System(field, metric)
```

and:

```go
state := qm.State(system)
hamiltonian := qm.Hamiltonian(system)
evolution := qm.Evolution(hamiltonian, state)
```

and:

```go
matter := relativity.StressEnergy(medium)
geometry := relativity.Metric(spacetime)
equation := relativity.EinsteinEquation(geometry, matter)
```

The library should make valid compositions easy and invalid compositions difficult.

---

# 10. Operations Are the Main API

The central API is the set of symbolic physics operations.

Examples:

```text
Differentiate
Integrate
Vary
Substitute
Simplify
Expand
Factor
Collect
Limit
Compose
Contract
RaiseIndex
LowerIndex
Transpose
Adjoint
Commute
Anticommute
Trace
Project
Expect
Transform
Normalize
FourierTransform
LaplaceTransform
Derive
Constrain
Couple
Evolve
Measure
```

These are operations on **physics objects**.

The AI should be able to use the same operation across many physics domains when the object supports it.

---

# 11. Operation Contracts

Operations should be expressed through Go interfaces wherever appropriate.

Conceptually:

```go
type Differentiable interface {
    Differentiate(by PhysicsVariable) PhysicsExpression
}

type Integrable interface {
    Integrate(over PhysicsDomain) PhysicsExpression
}

type Varyable interface {
    Vary(with PhysicsVariable) PhysicsExpression
}
```

The actual API should avoid unnecessarily forcing every object to implement every operation.

Prefer small composable interfaces.

For example:

```go
type Differentiator interface {
    Differentiate(Variable) Expression
}
```

is preferable to one enormous interface containing dozens of methods.

---

# 12. Symbolic Physics Expression

Although the public API contains physics objects rather than mathematical objects, an internal symbolic representation is necessary.

This internal representation exists to support:

- composition
- transformation
- simplification
- derivation
- canonicalization
- inspection

It is an implementation mechanism.

It is **not a public mathematical ontology**.

A useful internal abstraction may be:

```go
type Expression struct {
    ...
}
```

where every expression is constructed from physics-domain values and operations.

The public contract should remain physics-oriented.

---

# 13. Example: Einstein Equation

Do not require this:

```text
G_{μν} + Λg_{μν} = 8πGT_{μν}
```

as source syntax.

Instead the AI should construct:

```go
spacetime := relativity.Spacetime("M")
metric := relativity.Metric(spacetime)
matter := relativity.StressEnergy("T")
lambda := relativity.CosmologicalConstant("Λ")
gravity := relativity.EinsteinTensor(metric)

equation := relativity.EinsteinEquation(
    gravity,
    lambda,
    matter,
)
```

The library's internal symbolic representation can correspond to:

```text
EinsteinEquation
 ├── EinsteinTensor(metric)
 ├── CosmologicalConstantTerm(lambda, metric)
 └── StressEnergyTerm(matter)
```

The library may render that to conventional notation for humans.

The AI does not need to reason from the notation.

---

# 14. Example: Quantum Mechanics

The AI should be able to write:

```go
system := qm.System("particle")
state := qm.State(system)
hamiltonian := qm.Hamiltonian(system)

evolution := qm.Evolution(state, hamiltonian)
energy := qm.EnergyObservable(hamiltonian)

expectation := qm.Expect(state, energy)
```

The objects remain symbolic.

The library can transform them and derive new symbolic objects.

---

# 15. Example: Symbolic Differentiation

The same operation should work across physics domains.

Classical:

```go
acceleration := physics.Differentiate(velocity, time)
```

Field theory:

```go
gradient := physics.Differentiate(field, coordinate)
```

Quantum:

```go
response := physics.Differentiate(expectation, parameter)
```

Relativity:

```go
curvatureDerivative := physics.Differentiate(curvature, coordinate)
```

The operation's applicability is determined by the object's declared capabilities.

No branch-specific implementation should be required merely because the operation is mathematically the same.

---

# 16. Example: Symbolic Integration

Likewise:

```go
action := physics.Integrate(lagrangianDensity, spacetimeDomain)
```

```go
partition := physics.Integrate(
    statisticalWeight,
    stateSpace,
)
```

```go
amplitude := physics.Integrate(
    fieldExpression,
    pathDomain,
)
```

The operation remains generic.

The physics object supplies the domain-specific contract.

---

# 17. Example: Variation

Variation is a first-class physics operation:

```go
fieldEquation := physics.Vary(action, field)
```

The result should be a symbolic physics object appropriate to the transformation.

For a Lagrangian field theory, the resulting object may be an equation of motion.

---

# 18. Example: Derivation

AI agents should be able to construct chains such as:

```go
action := theory.Action(...)
variation := physics.Vary(action, field)
equation := physics.Simplify(variation)
```

The system should retain the symbolic relationship:

```text
Action
  ↓ Vary
Variation
  ↓ Simplify
EquationOfMotion
```

This gives an AI reviewer a machine-readable derivation.

---

# 19. Derivation Is Symbolic Reasoning

The library should make it natural to build chains such as:

```text
postulate
    ↓
define
    ↓
compose
    ↓
differentiate / integrate / vary
    ↓
simplify
    ↓
derive
    ↓
new equation
```

The library does not determine whether the final equation is physically true.

It determines whether the requested symbolic operations are formally supported.

---

# 20. Physics Operations Across Branches

The operation vocabulary should be deliberately cross-domain.

A single operation should be reusable whenever its contract applies.

For example:

```text
Differentiate
Integrate
Vary
Substitute
Simplify
Compose
Transform
Project
Expect
Trace
Contract
Commute
```

should not be duplicated as:

```text
QM_Differentiate
QFT_Differentiate
GR_Differentiate
StatMech_Differentiate
```

unless the semantics genuinely differ.

The object supplies the domain-specific behavior.

---

# 21. Interfaces as the Physics Grammar

Go interfaces are the primary mechanism for expressing shared capabilities.

Possible small interfaces:

```go
type Differentiable interface {
    Differentiate(Variable) Expression
}

type Integrable interface {
    Integrate(Domain) Expression
}

type Evolvable interface {
    Evolve(Parameter) Expression
}

type Observable interface {
    Expect(State) Expression
}

type Transformable interface {
    Transform(Transformation) Expression
}
```

These names are illustrative.

The final API must keep interfaces small and composable.

---

# 22. No Giant Type Matrix

Do not implement:

```text
PhysicsType × MathType × Regime × Operation
```

as a central table.

The applicability of an operation should come from:

- Go type signatures
- interfaces
- constructors
- domain requirements
- explicit physics relationships

This keeps the library extensible.

---

# 23. Math Reality Check in v0.3

MRC is retained, but dramatically simplified.

The library does **not** need a separate mathematical ontology to enforce MRC.

Instead:

> **MRC becomes a design rule for the AI-facing documentation and a consequence of the physics API's type boundaries.**

The public API only accepts physics objects.

The documentation teaches an AI agent how paper mathematics maps into those objects.

---

# 24. MRC Documentation Layer

The repository should contain documentation such as:

```text
docs/ai/
    physics-objects.md
    notation-to-physics.md
    mrc.md
```

The documentation answers:

```text
Paper notation
      ↓
What physics object does this denote?
      ↓
What operation does it participate in?
      ↓
What physics object results?
```

For example:

```text
paper: Gμν
meaning: EinsteinTensor
library: relativity.EinsteinTensor

paper: Tμν
meaning: StressEnergy
library: relativity.StressEnergy

paper: H
meaning: Hamiltonian
library: qm.Hamiltonian

paper: ψ
meaning: QuantumState / StateVector depending on context
library: qm.State
```

The exact mapping is contextual.

The AI must not infer physics identity from notation alone.

---

# 25. MRC Category-Error Rule

The fundamental rule remains:

> **A mathematical symbol in a paper is not itself a physics object. The AI must translate the symbol into the appropriate physical object before operating on it.**

For example:

```text
n
```

does not automatically mean:

```text
Time
```

or:

```text
ParticleNumber
```

or:

```text
QuantumNumber
```

The AI must determine the physical role from context and select the appropriate library object.

Likewise:

```text
G
```

may mean very different physical objects depending on context.

The documentation exists to guide that translation.

---

# 26. The Library Prevents Category Errors After Translation

Once translated into physics objects, Go's type system and the library's operation contracts take over.

Example:

```go
order := somePhysicsObject
time := physics.Time(order)
```

should fail unless a valid physics constructor explicitly accepts that source object.

Likewise, an operation expecting a Hamiltonian should not accept an unrelated physical object.

The protection is:

```text
AI translation
      ↓
Physics type
      ↓
Go type checking
      ↓
physics operation contract
```

---

# 27. Representation vs Identity

The library should distinguish relationships such as:

```text
Represents
Identifies
Interprets
Constrains
Couples
```

but only where these relationships are actually useful to the physics object model.

Do not build a large epistemic framework.

The purpose is to preserve important distinctions in AI representations.

---

# 28. Physical Objects May Contain Mathematical Structure Internally

This is permitted and expected.

For example, a `Metric` necessarily has mathematical structure internally.

A `Hamiltonian` necessarily has symbolic algebraic structure internally.

A `Field` necessarily has domain and transformation structure internally.

The rule is:

> **Do not expose a general mathematical ontology as the primary public API.**

The implementation can use whatever internal structures are necessary.

---

# 29. Internal Symbolic Engine

The internal engine should be as small and reusable as practical.

It may contain:

```text
expression nodes
symbol references
operation nodes
transformation rules
canonicalization
rewrite rules
derivation links
```

But these are implementation details.

Do not make the user or AI explicitly construct low-level mathematical AST nodes unless there is a compelling reason.

Prefer:

```go
result := physics.Differentiate(field, coordinate)
```

over:

```go
expr := ast.NewDerivative(...)
```

---

# 30. Symbolic Transformation Rules

Transformations should be represented as explicit rules.

Examples:

```text
Differentiate
Integrate
Substitute
Simplify
Expand
Factor
Collect
Vary
Contract
Commute
Trace
Project
Transform
```

The transformation subsystem should be extensible by physics package.

For example:

```text
std/mechanics
std/relativity
std/qm
std/qft
std/statmech
```

may contribute domain-specific transformations.

---

# 31. No Universal CAS Requirement

The project does not need to become a general-purpose computer algebra system.

The objective is:

> **enough symbolic manipulation to reason about and derive theoretical physics.**

Implement operations according to physics usefulness and compositionality.

Do not make arbitrary mathematical theorem proving a prerequisite.

---

# 32. Canonical Symbolic Representation

Every physics object produced by the library should have a deterministic symbolic representation.

This supports:

- comparison
- transformation
- caching
- derivation inspection
- serialization
- AI review
- documentation
- reproducible tests

Two equivalent constructions should normalize identically where the relevant transformation rules justify equivalence.

Do not reorder noncommuting operations.

---

# 33. Derivation Trace

The library should retain enough information to expose:

```text
input object
operation
arguments
result
```

For example:

```text
Lagrangian
    ↓ Differentiate(field)
EulerLagrangeTerm
    ↓ Simplify()
EquationOfMotion
```

This is essential for AI reasoning and adversarial review.

---

# 34. Adversarial AI Boundary

The adversarial reviewer is a separate AI system.

It can consume:

- Go source
- physics objects
- symbolic representation
- operation traces
- derived equations
- documentation
- diagnostics

The Physics library does not perform the adversarial review itself.

Its job is to make the reasoning representation sufficiently explicit for another agent to review.

---

# 35. AI Paper Translation Boundary

Paper understanding belongs to the AI layer.

The workflow is:

```text
paper notation
    ↓
AI semantic interpretation
    ↓
Physics Compiler object
```

The library does not need to understand arbitrary LaTeX.

This drastically reduces compiler complexity while preserving the semantic goal.

---

# 36. Example End-to-End Workflow

A paper describes:

```text
H
ψ
O
⟨ψ|O|ψ⟩
```

The AI translates these into:

```go
system := qm.System("system")
state := qm.State(system)
hamiltonian := qm.Hamiltonian(system)
observable := qm.Observable(system)
```

Then:

```go
expectation := qm.Expect(state, observable)
```

and:

```go
evolution := qm.Evolve(state, hamiltonian)
```

and perhaps:

```go
response := physics.Differentiate(expectation, parameter)
```

The AI now reasons using typed physics objects rather than manually manipulating paper notation.

---

# 37. Standard Library Organization

A practical initial repository may look like:

```text
physics/
    core/
    mechanics/
    fields/
    relativity/
    qm/
    qft/
    statmech/
    operations/
    transform/
    inspect/
    docs/
    examples/
    tests/
```

The exact package decomposition should be kept small initially.

Avoid premature fragmentation.

---

# 38. Core Package

`core` should contain only shared infrastructure such as:

- physics object identity
- symbolic references
- operation representation
- derivation trace
- common errors
- common composition mechanisms

It should not contain every physics type.

---

# 39. Domain Packages

Each domain package owns its physics vocabulary.

For example:

```go
qm.Hamiltonian
qm.State
qm.Observable
qm.Measurement
```

while:

```go
relativity.Spacetime
relativity.Metric
relativity.EinsteinTensor
relativity.StressEnergy
```

This follows Go's package philosophy.

Domain packages should expose coherent APIs rather than one enormous `physics` package.

---

# 40. Operations Package

Common operations should be reusable:

```go
operations.Differentiate(...)
operations.Integrate(...)
operations.Vary(...)
operations.Substitute(...)
operations.Simplify(...)
operations.Compose(...)
operations.Transform(...)
```

Domain packages may provide specialized operations when needed.

The operation system should prefer generic interfaces over domain-specific duplication.

---

# 41. Inspection

Provide explicit inspection APIs so AI agents can inspect what they have constructed.

Examples:

```go
inspect.Type(object)
inspect.Structure(object)
inspect.Operations(object)
inspect.Derivation(object)
inspect.Symbolic(object)
inspect.Render(object)
```

The exact API is to be designed during implementation, but the capability is required.

---

# 42. Rendering

The library may render a physics object into conventional notation.

For example:

```go
inspect.Render(equation)
```

may produce human-readable notation such as:

```text
Gμν + Λgμν = 8πGTμν
```

The rendered notation is presentation.

The semantic object remains the source of truth for the AI/toolchain.

---

# 43. Serialization

Physics objects and derivation traces should have deterministic serialization.

The serialization must preserve:

- object identity
- types
- relationships
- symbolic operations
- derivation order
- canonical structure

JSON is a reasonable interoperability format.

The serialization format is an output representation, not a second programming language.

---

# 44. Testing

Tests should be organized around physics operations and objects.

Examples:

```text
construct Hamiltonian
differentiate Hamiltonian
derive equation of motion
simplify equation
integrate Lagrangian density
construct expectation
compute symbolic commutator
transform field
construct Einstein equation
derive curvature relation
```

Tests should verify symbolic structure, not numerical physics.

---

# 45. Cross-Domain Conformance Tests

The same operation should be tested across domains.

For example:

```text
Differentiate:
    mechanics
    relativity
    field theory
    quantum mechanics
    QFT
    statistical mechanics
```

and:

```text
Integrate:
    mechanics
    field theory
    QFT
    statistical mechanics
```

and:

```text
Simplify:
    all supported symbolic physics structures
```

This ensures that operations are genuinely reusable.

---

# 46. Invalid Operation Tests

Tests must also establish useful failure modes.

Examples:

```text
Differentiate(object without differentiable physical structure)
Integrate(object without an integration domain)
Expect(object that is not an observable)
Evolve(object without evolution law)
Commute(object that is not operator-like
)
```

Errors should identify the violated physics contract.

---

# 47. What Is Explicitly Removed From v0.3

The following are removed from the architecture:

```text
custom .phys language
custom scanner
custom parser
custom AST frontend
custom grammar
custom physics type language
custom module system
public math-object ontology
public math branch taxonomy
universal jurisdiction matrix
Dancer/Dance/Stage compiler types
research methodology
evidence management
scientific adjudication
truth scoring
paper litigation
numerical runtime
simulation engine
GPU/CPU code generation
```

They are not needed to accomplish the core objective.

---

# 48. What Remains From the Earlier MRC Work

The useful insight remains:

> **Do not let mathematical notation masquerade as a physical object.**

But the implementation is now radically simpler:

```text
paper notation
      ↓
AI translation/documentation
      ↓
typed physics object
      ↓
Go type system
      ↓
physics operation contracts
```

The old MRC jurisdiction machinery is therefore **retired**.

MRC survives as a documentation principle and as formal API discipline.

---

# 49. Program A Relationship

Program A is now clearly separate.

Program A was a declarative semantic checker.

The new system is a Go physics library designed for symbolic reasoning.

Program A can provide:

- historical lessons
- test cases
- useful domain distinctions

but its artifact model, verdict model, SemanticSpec, candidate state system, and MRC engine should not be imported into this library.

---

# 50. Definitive Architecture

The entire system now becomes:

```text
                 PHYSICS PAPERS
                       │
                       ↓
              AI PAPER UNDERSTANDING
                       │
             notation → physics meaning
                       │
                       ↓
                  ORDINARY GO
                       │
         ┌─────────────┼─────────────┐
         ↓             ↓             ↓
      variables     functions     interfaces
         │             │             │
         └─────────────┼─────────────┘
                       ↓
               PHYSICS LIBRARY
                       │
          ┌────────────┼─────────────┐
          ↓            ↓             ↓
      physics       operations    composition
      objects
          │            │             │
          └────────────┼─────────────┘
                       ↓
             SYMBOLIC TRANSFORMATION
                       ↓
                  DERIVATION
                       ↓
               NEW EQUATIONS
                       ↓
             INSPECTION / TRACE
                       ↓
              ADVERSARIAL AI
```

There is no second programming language.

There is no custom compiler frontend.

There is no public mathematics ontology.

There is no numerical execution engine.

---

# 51. The Core Design Principle

> **Use Go to represent physics objects, and use typed symbolic operations to reason over them.**

Mathematical notation is how physics papers communicate.

The AI translates that notation into physics objects.

The Physics library provides the operations.

The Go compiler provides the programming-language semantics.

The symbolic engine provides the transformation machinery.

The adversarial AI provides independent review.

---

# 52. One-Sentence Definition

> **The Physics Compiler is a Go library that gives AI agents a typed, compositional, symbolically manipulable vocabulary of theoretical-physics objects and operations, enabling them to translate physics literature into machine-readable physics, reason through symbolic transformations, and derive new equations without introducing a new programming language or a public mathematical ontology.**

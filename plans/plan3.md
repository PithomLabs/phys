## prompt

just to clarify, this physics compiler is purely for theoretical physics only (symbols only), there is no actual computing of values, that is what I design it for (symbol manipulation as if on pen and paper)


## answer


Yes. That clarification materially sharpens the architecture.

The **Physics Compiler is a symbolic theoretical-physics language**, not a numerical physics programming language.

Its purpose is essentially:

> **Write physics the way a theoretical physicist writes mathematics on paper, but in a formal language that can parse, type-check, compose, transform, and inspect the symbolic structure.**

There is **no evaluation of physical values**.

### What the compiler does

```text
.phys source
    ↓
Scanner
    ↓
Parser
    ↓
AST
    ↓
Name Resolution
    ↓
Symbol / Type System
    ↓
Dimensional / Tensor / Index Checking
    ↓
Interface & Composition Checking
    ↓
Symbolic Semantic Analysis
    ↓
Symbolic Physics IR
    ↓
Transformation / Inspection / Pretty Printing
```

The central object is therefore **symbolic structure**, not executable values.

For example:

```phys
L = 1/2 m v^2
```

The compiler should understand that:

* `L` is a symbol
* `m` is a symbol with a declared mathematical/physical type
* `v` is a symbol
* multiplication and exponentiation are symbolic operations
* the expression has a well-formed structure
* dimensional relations can be checked symbolically

But it should **not** calculate:

```text
m = 2 kg
v = 3 m/s
L = 9 J
```

That is outside the language.

---

## The compiler should behave more like mathematical notation than Mathematica

There is an important distinction.

### In scope

Symbolic transformations such as:

```text
d/dx (x² + 3x)
```

→

```text
2x + 3
```

or

```text
A B C
```

with declared operator types and algebraic relationships.

Likewise:

```text
∇_μ V^μ
```

should be represented structurally, allowing the compiler to understand:

* index position
* tensor rank
* contraction
* covariance/contravariance
* manifold/context
* derivative operator
* resulting type

And:

```text
H = p² / 2m + V(x)
```

is a symbolic Hamiltonian expression.

### Out of scope

The compiler does not:

* substitute numerical values
* run simulations
* solve numerical differential equations
* evaluate integrals numerically
* run Monte Carlo
* perform lattice calculations
* execute PDE solvers
* perform GPU computation
* generate numerical kernels
* determine empirical predictions by computation

Even symbolic computation should be understood as **formal manipulation of expressions**, not execution of a physical model.

---

# This also changes how I would define "Physics IR"

The Physics IR should **not** be an execution IR.

It should be a **typed symbolic representation**.

For example:

```text
Expression
 ├── Symbol
 ├── NumberLiteral
 ├── Sum
 ├── Product
 ├── Power
 ├── FunctionApplication
 ├── Derivative
 ├── Integral
 ├── Limit
 ├── TensorExpression
 ├── Contraction
 ├── OperatorApplication
 ├── Composition
 └── Relation
```

with physics-aware annotations attached to those mathematical structures.

So this:

```phys
EinsteinEq:
    G[μν] + Λ g[μν] = 8πG T[μν]
```

becomes something closer to:

```text
Equality
 ├── Add
 │    ├── EinsteinTensor[μν]
 │    └── Scale(Λ, MetricTensor[μν])
 └── Scale(8πG, StressEnergyTensor[μν])
```

The compiler understands the **symbolic grammar** of the equation.

It does not "run Einstein's equations."

---

# The mathematical expression system is therefore central

This is the biggest consequence of your clarification.

For Program A, we intentionally avoided a general mathematical expression AST.

For the **Physics Compiler, that restriction is no longer appropriate**.

A real Physics Compiler needs a serious symbolic expression language because the language itself must be capable of expressing things like:

```text
S = ∫ d⁴x √(-g) L
```

```text
δS = 0
```

```text
[p, q] = -iℏ
```

```text
F_{μν} = ∂_μ A_ν - ∂_ν A_μ
```

```text
R^ρ_{ σμν}
```

```text
Z = Tr(e^{-βH})
```

```text
ψ(x,t) = ...
```

Those are **expressions**, not declarations of executable computations.

---

# Symbolic algebra should be treated as language semantics

This is where the Gemini/DeepSeek material becomes much more important.

The compiler should know formal operations such as:

```text
+
-
*
/
^
=
≈
≠
∘
∂
∇
∫
Σ
Π
lim
δ
[ , ]
Tr
exp
log
det
```

and structural mathematical operations such as:

```text
differentiate
integrate
substitute
simplify
factor
expand
collect
contract
raise index
lower index
transpose
adjoint
commute
compose
vary
take limit
```

But these are **symbolic transformations**.

For example:

```phys
expr := x^2 + 2*x + 1

simplify(expr)
```

can produce:

```text
(x + 1)^2
```

without ever assigning a numerical value to `x`.

---

# This also clarifies the role of physics

The compiler should distinguish:

### Mathematical structure

```text
Scalar
Vector
Covector
Tensor
Function
Operator
Manifold
Metric
Connection
Group
Algebra
Measure
Distribution
Spectrum
...
```

from:

### Physical structure

```text
Field
State
Observable
Particle
System
Spacetime
Action
Lagrangian
Hamiltonian
Constraint
Symmetry
Charge
Current
Energy
Momentum
Coupling
InitialCondition
BoundaryCondition
...
```

But those should generally be **library-level typed constructs**, not a gigantic compiler-level enumeration.

For example:

```phys
field ScalarField φ on M
metric g on M
connection ∇ compatible_with g
```

The compiler checks whether the symbolic constructions are formally valid.

It does not decide:

> "Is φ really a physically valid field?"

That remains outside the compiler.

---

# Interfaces become even more important

The Go-style capability idea fits this beautifully.

For example, conceptually:

```go
type Differentiable interface {
    Derivative(...) Expression
}

type Tensor interface {
    Rank() int
    Variance() Variance
}

type Operator[A, B] interface {
    Domain() A
    Codomain() B
}
```

These are **symbolic contracts**.

They let the compiler reason compositionally.

For example:

```text
Tensor
    +
Index
    +
Contraction
    ↓
TensorExpression
```

or:

```text
Function[A,B]
    +
Function[B,C]
    ↓
Function[A,C]
```

This is much closer to the way mathematics actually works.

---

# "No actual computing of values" should be a hard architectural constraint

I would therefore add this explicitly to the locked specification:

> **The Physics Compiler is a symbolic compiler for theoretical physics. It manipulates symbols, expressions, mathematical structures, and formal relationships. It does not evaluate physical quantities or execute numerical computation.**

And consequently:

### Core compiler

```text
language
grammar
AST
symbols
types
interfaces
generics
composition
dimensions
tensor/index semantics
symbolic expressions
symbolic transformations
physics structures
IR
formatting
inspection
documentation
```

### Explicitly outside the compiler

```text
numerical evaluation
simulation
measurement processing
data fitting
parameter estimation
numerical PDE solving
Monte Carlo
GPU execution
CPU kernels
scientific workloads
experimental analysis
```

That also means the previous **"Program B / computational backend"** concept should not be treated as a planned part of the Physics Compiler itself. It can exist externally as a future consumer of exported symbolic representations, but **the Physics Compiler has no need to become an execution compiler**.

The clean mental model is:

> **Go compiler + computer algebra notation + theoretical-physics type system.**

Not:

> physics simulation language.

That distinction should be reflected throughout `PHYSICS_COMPILER_PLAN.md`.


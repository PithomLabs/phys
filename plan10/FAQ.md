## question

Does the codebase has a mechanism to prevent category errors in physical and mathematics primitives when manipulating variables during derivation

## answer

Yes — **the codebase has a real mechanism for preventing a substantial class of physical and mathematical category errors during derivation**, and it is one of the central design features of the MVP.

The key point is that it uses **two independent type-like checks**:

### 1. Dimension checking — “Do the units make sense?”

Every trusted `core.Object` carries an immutable `Dimension`. The system has exact rational dimension arithmetic, rather than floating-point unit metadata. specs_v2_3

For operations:

- `Add` / `Subtract` require equal dimensions.
- `Compare` requires equal dimensions.
- `Multiply` composes dimensions.
- `Divide` divides dimensions.
- `Pow` raises the dimension to an exact rational exponent.
- `Substitute` requires the replacement to have the **same dimension as the variable**.
- `Differentiate` computes dimension as `target / wrt`. 

So something like

```text
Energy + Velocity
```

cannot silently proceed. The dimensional check rejects it with `DimensionMismatchError`. The validation order is deterministic: dimensional failure occurs before the physical-kind check. plan10_v2_3

### 2. Physical-kind checking — “Even with the same units, are these the same category?”

This is the more interesting part.

`Dimension` and `Kind` are deliberately **not the same thing**. The specification explicitly exists to prevent cases where two quantities have the same dimensions but are not automatically interchangeable. For example, energy and another energy-dimensional quantity can still belong to different semantic categories. specs_v2_3

The MVP has a closed `Kind` set including:

```text
Mass
RestMass
Energy
KineticEnergy
Momentum
ThreeMomentum
FourMomentum
Velocity
...
Expression
Relation
BranchSet
```

and MRC-003 applies a strict compatibility matrix. specs_v2_3

For example:

```text
Mass + RestMass
```

has the same dimension, but different named kinds, so it is rejected with:

```text
CategoryMismatchError
```

That is specifically tested as a same-dimension category-error case, so the test isolates MRC-003 rather than accidentally triggering dimensional failure. plan10_v2_3

### 3. Variables are not just raw strings

This is especially important for your question about **manipulating variables during derivation**.

A substitution variable must itself be a valid physical object whose expression is a **single `Symbol`**, and the replacement must have:

```text
same Dimension
same Kind
```

as the variable. physics_compiler_mvp_specs_v2_1

So the system does not merely say:

```text
replace "m" with X
```

It effectively constrains the substitution to something like:

```text
m : Mass[M]
X : Mass[M]
```

rather than allowing:

```text
m : Mass[M]
X : Length[L]
```

That is a meaningful safeguard against an AI accidentally changing the semantic type of a variable while performing algebra.

### 4. Derivation operations preserve the type information

The symbolic expression tree itself is deliberately more general than the physical object wrapper, but every resulting `Object` gets newly computed metadata.

For example:

```text
Multiply(m, a)
```

produces:

```text
Kind       = Expression
Dimension  = M × (L/T²)
```

rather than pretending the result is still `Mass` or `Acceleration`. specs_v2_3

Likewise:

```text
Differentiate(K, v)
```

gets a dimension derived mechanically as:

```text
[K] / [v]
```

and the resulting kind is `Expression`. specs_v2_3

That is exactly the behavior you want from a symbolic substrate: once you perform general algebra, the system does not invent a stronger physical category than it actually knows.

### 5. Assumptions and conventions are another layer of protection

There is also a third category of semantic consistency beyond dimensions and kinds.

Operations merge assumptions and conventions, and conflicting values for the same keyed assumption or convention cause explicit errors rather than silently combining incompatible premises. plan10_v2_3

For example, a derivation cannot silently end up carrying both incompatible versions of a premise under the same `(Kind, Key)`.

### But there is an important limitation

I would **not** describe the current MVP as a complete prevention mechanism for *all* category errors.

It prevents errors that the current ontology knows how to express:

```text
dimension mismatch
physical-kind mismatch
invalid substitution typing
invalid relation operands
incompatible assumptions
incompatible conventions
invalid branch-selection structures
```

But it does **not yet have a full semantic type system for physics**.

For example, it does not currently know that:

```text
position x
velocity v
four-vector Xμ
metric tensor gμν
scalar field φ
vector field Vμ
operator H
wavefunction ψ
```

belong to different mathematical structures with rules for index variance, tensor rank, manifold domain, coordinate charts, function spaces, operator domains, etc.

That's a much deeper level of category checking than the current `Kind + Dimension` model.

And this distinction matters for your longer-term TOE/compiler idea:

> **The MVP already has the seed of a physics type system, but it is not yet a full mathematical-physics type system.**

The current architecture gives you a particularly good foundation for extending that later without confusing it with ordinary algebra.

In fact, the separation:

```text
Kind
Dimension
Expr
Assumptions
Conventions
Provenance
```

is quite useful because a future **kernel growth gate** could ask whether a proposed new invariant is genuinely foundational and theory-neutral, rather than dumping every new physical category into the kernel.

So, at the current frozen level: **yes, category-error prevention is built in and tested; but it is currently strongest at the physical-quantity level, not yet at the full tensor/mathematical-structure level.**
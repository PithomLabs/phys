Below is the consolidated **coding-agent prompt**. It is intentionally written as a **plan-writing / implementation-planning task**, not an implementation task. The goal is to make the agent understand the architectural intent behind Plan 11 → Plan 11.1, the frozen kernel philosophy, and the GR stress-test purpose before proposing any code changes.

---

# Coding Agent Prompt — Write Revised Plan 11.1 Implementation Plan for GR Stress Test

## Role

You are a senior Go architect designing the implementation plan for the next phase of the `phys` project.

Your task is **not to implement General Relativity yet**.

Your task is to produce a cohesive implementation plan for **Plan 11.1 Revised: General Relativity as a Kernel Stress Test**, incorporating all architectural decisions, adversarial findings, and governance constraints described below.

The output must be an implementation plan that another coding agent can execute without reopening architectural decisions.

---

# 1. Project Philosophy (Non-Negotiable)

The inspiration is Go itself.

The system has:

```
small robust core
        +
rich composable ecosystem
```

The kernel must remain:

- math-agnostic
- physics-theory-agnostic
- minimal
- composable
- extensible
- invariant-preserving

The kernel is **not** a physics textbook.

The kernel does **not** contain:

- General Relativity objects
- Quantum Mechanics objects
- QFT objects
- tensors
- metrics
- manifolds
- Hilbert spaces
- Hamiltonians
- Lagrangians
- field equations
- domain-specific mathematical structures

Those belong to userland packages.

The kernel only provides primitives rich enough that userland theories can be built on top.

The analogy is:

```
Go language/runtime
        ↓
standard library
        ↓
third-party ecosystem
```

The equivalent architecture is:

```
phys kernel
        ↓
shared mathematical libraries (future Level 2)
        ↓
theory packages
        ↓
GR / QM / QFT / future theories
```

---

# 2. Current Architecture Rule

The primary criteria:

## Kernel must not answer:

"What mathematical objects belong to physics?"

Instead:

## Kernel must answer:

"What primitive capabilities are universally useful across theories?"

Examples that may eventually become kernel/shared primitives:

- exact rational arithmetic
- immutable symbolic expressions
- dimensional algebra
- provenance tracking
- assumption tracking
- canonical replay
- deterministic transformations

Examples that remain userland:

- Schwarzschild metric
- Christoffel symbols
- tensors
- curvature tensors
- Einstein tensor
- Hilbert space
- wavefunction
- operators
- field tensors

---

# 3. Purpose of GR Implementation

GR is not being added to the kernel.

GR is a **stress test**.

The question is:

> Can the current kernel primitives support a serious theory package without being modified?

The experiment must discover:

1. What the kernel can represent naturally.
2. What requires userland mathematics.
3. Where the abstraction boundary is correct.
4. Whether any kernel limitation represents a true foundational gap.

Do not start from:

"GR needs tensors, therefore add tensors."

That violates the experiment.

---

# 4. Repository / Module Boundary

The frozen `phys` tree remains untouched.

GR must live outside the kernel repository.

Expected architecture:

```
github.com/PithomLabs/phys

    frozen kernel
    core
    ops
    session
    existing mechanics
    existing relativity


github.com/PithomLabs/phys-gr

    Level-1 GR implementation
    depends on phys
    never imported by phys
```

No reverse dependency.

No kernel contamination.

No GR-specific exceptions inside `phys`.

---

# 5. Revised GR Workload Anchor

The primary workload:

```
coordinates/chart
        ↓
metric
        ↓
inverse metric
        ↓
Christoffel symbols
        ↓
covariant derivative
        ↓
Riemann curvature
        ↓
Ricci tensor/scalar
        ↓
Einstein tensor
        ↓
vacuum equation
        ↓
Schwarzschild exterior
        ↓
weak-field/Newtonian limit
```

This remains the chosen stress path.

However:

Do not claim Schwarzschild is transcendental-free.

Standard spherical coordinates include:

```
sin²(theta)
```

Therefore:

- trigonometric functions belong to GR userland math
- they do not imply kernel expansion

---

# 6. Critical Pin: Differentiation Boundary

This is the most important clarification.

The frozen kernel contains:

```
Differentiate()
```

but Pin B intentionally limits negative integer power differentiation.

Example:

```
Differentiate(Pow(r,-1))
```

must fail closed:

```
UnsupportedOperationError
```

This behavior is intentional.

The GR workload must not silently expand the kernel.

Required execution sequence:

## GR-3a

First test:

```
phys.Differentiate(metric expression)
```

Expected:

```
failure
+
logged as SPEC-INTENDED-BOUND
```

## GR-3b

If GR needs continuation:

implement:

```
phys-gr/symbolic/differentiate.go
```

as a theory-local differentiator.

Rules:

- must use public `core.Expr`
- must not modify `phys`
- must not pretend kernel supports GR differentiation
- must preserve kernel invariants

Classification:

```
kernel limitation:
SPEC-INTENDED-BOUND

GR package capability:
PACKAGE-SOLVABLE
```

Both results must be recorded.

---

# 7. Expression Boundary

Respect the frozen expression model.

The kernel exposes:

- immutable expressions
- constructors
- read-only inspection

The kernel does not expose arbitrary new expression nodes.

Do not introduce:

```
Tensor
Derivative
Metric
Field
Manifold
```

as kernel expressions.

Those are GR package structures.

---

# 8. Object / Type System Rule

Do not weaken the existing authority model.

The kernel has:

- immutable Object
- controlled construction
- provenance
- corpus status
- MRC validation

No generic:

```
NewObject(kind,...)
```

factory.

GR cannot manufacture trusted kernel objects by choosing:

```
Kind
Dimension
Provenance
CorpusStatus
```

arbitrarily.

If GR needs richer semantic structures:

keep them in GR userland.

---

# 9. Required Plan Adjustments

The revised Plan 11.1 document must include:

---

## A. Pin A correction

The Identify regression fixture must be:

Wrong:

```
Identify(Energy, E/c²)
```

because dimensions differ.

Correct:

```
Identify(RestMass, E/c²)
```

Assertions:

```
Kind == Relation
Provenance == IDENTIFIED
CorpusStatus == NONE
Dimension == Mass
```

Must not become:

```
RestMass
```

---

## B. Explicit C1-C8 Concern Ledger

Define all concerns in the plan.

Example:

```
C1 Type decay / semantic loss

C2 Derived-expression identification

C3 Assumption consistency

C4 Symbol identity

C5 General mathematical functions

C6 Kind ontology extensibility

C7 Identify regression

C8 Differentiation regression
```

Each concern must have:

- description
- classification
- routing decision

---

## C. Silent Wrongness Probe

Do not rely only on GR workload.

Create explicit adversarial tests.

Examples:

- symbol ambiguity
- incorrect field differentiation
- invalid contraction
- dimension-valid semantic errors

The dangerous outcome is:

```
accepted
+
wrong
```

not:

```
unsupported
```

---

## D. Level-2 Governance

Add:

### L1 → L2 promotion requires:

```
named second theory consumer
+
worked example
+
shared semantics
+
no assumption leakage
```

### L2 → L3 requires:

```
proof that generic invariant cannot be preserved above kernel
```

Level 2 cannot become a hidden kernel extension mechanism.

---

## E. Growth Evidence Record

Do not reference an undefined "20-field record."

Enumerate it.

Must include:

- capability
- minimal counterexample
- failure category
- Level-1 workaround
- Level-2 workaround
- why workaround fails if applicable
- genericity
- foundationality
- second consumer
- worked example
- MRC impact
- replay impact
- migration impact
- kernel change required
- human review evidence

---

# 10. Failure Classification

Every obstacle must map exactly to one:

```
SPEC-INTENDED-BOUND

PACKAGE-SOLVABLE

REPRESENTABLE-BUT-UNFAITHFUL

UNREPRESENTABLE

SILENTLY-WRONG
```

No wishlist items.

No:

"it would be nice if the kernel had tensors."

Only evidence.

---

# 11. Expected Plan Structure

Produce a document containing:

```
1. Executive summary

2. Architectural principles

3. Frozen baseline assumptions

4. Module boundaries

5. GR workload

6. Representation strategy

7. Kernel interaction points

8. Level-1 GR package design

9. GR-local mathematics required

10. Failure taxonomy

11. C1-C8 concern ledger

12. Silent-wrongness testing strategy

13. Level-1 → Level-2 promotion rules

14. Level-2 → Level-3 Growth Gate

15. Evidence record format

16. Test strategy

17. Mutation testing

18. Expected outcomes

19. Explicit non-goals

20. Definition of done
```

---

# 12. Final Decision Principle

The plan must preserve this rule:

> Do not make the kernel smarter because GR is complicated. Make userland richer first. Only promote something into the kernel when multiple theories independently prove that the primitive itself is foundational, theory-neutral, and impossible to safely implement above the kernel boundary.

The success condition is not:

"GR works."

The success condition is:

"We learned exactly where the kernel abstraction boundary holds, where userland must grow, and whether any genuine foundational primitive is missing."

---

Produce only the implementation plan. Do not write code. Do not modify architecture. Do not propose kernel growth unless the evidence criteria above are satisfied.
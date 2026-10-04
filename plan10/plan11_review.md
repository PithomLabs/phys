Yes. After consolidating the four reviews, I would **lock the following as the pre-GR architecture and test contract**.

The consensus is strong: the six concerns do **not** justify kernel expansion. The useful result is that they expose the boundary between the frozen trusted substrate and the richer mathematical/semantic layers that will eventually sit above it. DeepSeek, Gemini and Qwen independently converge on the three-level model, while Z’s adversarial pass identifies a few places that must be pinned before the GR stress test begins. plan11_deepseek plan11_qwen

# Consolidated pre-GR decision

## 1. Architecture is now explicitly three-level

```text
LEVEL 1 — Theory-local userland
    general_relativity/
    quantum/
    qft/
    ...
        ↓ demonstrated reuse

LEVEL 2 — Shared mathematical / semantic libraries
    tensor/
    symbols/
    semantic/
    constraints/
    ...
        ↓ proven need for trusted invariant

LEVEL 3 — Trusted kernel
    internal/kernel
    core
    ops
    session
```

Level 1 and Level 2 cannot modify Level 3. A Level-1 theory does not jump directly into the kernel. A Level-2 promotion is human-curated and requires actual demonstrated reuse, not merely "two theories could plausibly use this." Kernel promotion requires an unsuccessful higher-level attempt to preserve the required invariant. plan11_deepseek

This is the cleanest expression of the Go analogy: the kernel is analogous to the language/runtime primitives; richer libraries accumulate above it. Tensors, Hilbert spaces, symbols, differential forms, etc. are not automatically kernel concepts. plan11_gemini

## 2. The six concerns are retained, but their routing is settled

| Concern | Consolidated disposition |
|---|---|
| **Type decay** | Real limitation; **Level 2**, not kernel |
| **Substitution catch-22** | Real limitation; **Level 2 semantic ascription**, do not weaken `Substitute` |
| **Assumption semantic conflict** | Real limitation; **Level 2 constraint/entailment layer** |
| **Raw string symbol identity** | Real limitation; **Level 2 symbol/entity layer** |
| **Transcendentals / symbolic exponents** | Not current MVP capability; future extension |
| **Closed Kind ontology** | Correctly closed; do not turn it into a subtype hierarchy |

That triage is directly supported by the reviews. plan11_deepseek plan11_qwen

### Important refinement on type decay

The observed failure is real:

```text
m × c²  → Expression, M L² T⁻²
r × F   → Expression, M L² T⁻²

Expression + Expression → permitted
```

So the kernel can no longer distinguish the physical meaning of the two derived quantities. But that does **not** mean the kernel needs more physics ontology. The missing capability is semantic ascription above the kernel. plan11_deepseek

This is actually consistent with the core design principle:

> **Algebra does not manufacture physical meaning.**

---

# 3. Three corrections must be incorporated before GR

This is the part I would treat as **pre-GR conformance cleanup**, not kernel redesign.

## A. Pin `Session.Identify`'s result kind

Z identified a legitimate ambiguity: the frozen material specifies that `Identify` constructs

```text
Relation(eq, a.Expr(), b.Expr())
```

and records an `IDENTIFIED` provenance result, but does not explicitly pin the object's `Kind`. plan10_v2_3

That matters because `Compare` already permits named-vs-Expression comparisons.

I would **not adopt Z's suggested `Kind = Expression`**.

The more mechanically coherent rule is:

> **`Session.Identify` always returns `KindRelation`, because it creates an identification relation; it never retypes either operand and never converts an `Expression` into a nominal physical kind.**

That preserves the already specified construction semantics and keeps the semantic-ascription capability out of the trusted session path. The returned expression is literally a `Relation(eq, ...)`, and ordinary `Compare` results are likewise `KindRelation`. plan10_v2_3 specs_v2_3

Before GR, this should be mechanically pinned by implementation + test.

**Required test:**

```text
Identify(Energy, Expression(E), "...")
    => KindRelation
    => Expr = Relation(eq, ...)
    => Provenance = IDENTIFIED
```

and, critically:

```text
Identify(Energy, E/c², ...)
```

must **not** result in `RestMass`.

That closes the only plausible covert-ascription route.

## B. Explicitly close the negative-integer differentiation ambiguity

Z is also correct that the current wording leaves a small ambiguity around:

```text
Pow(x, -1)
```

The supported differentiation rule is explicitly for **non-negative integer** exponents, while the unsupported wording mentions non-integer symbolic exponents. specs_v2_3

Before GR, pin this as:

> **`Differentiate` supports only non-negative integer constant exponents. Negative integer exponents are currently unsupported and return `UnsupportedOperationError`.**

Do **not** add rational-function differentiation to the MVP merely to make GR easier.

That keeps the MVP bounded and makes a GR encounter with `r⁻¹` an honest "spec-intended bound" rather than an ambiguous implementation failure.

## C. Narrow the existing documentation claim about physical-kind checking

The reviews correctly point out that the present protection is conditional.

The accurate statement is:

> **Named-kind compatibility protects additive operations when named kinds are retained. `Expression + Expression` with equal dimensions is explicitly permitted by MRC-003; preserving or restoring higher-level semantic category after algebra is outside the MVP.**

This should replace any wording suggesting that the current kernel mechanically catches every energy-vs-torque-style mismatch. plan11_z

---

# 4. Keep these things OUT of the kernel

This should be an explicit GR rule.

The GR implementation must not cause us to add:

```text
Tensor
Manifold
MetricTensor
Connection
Curvature
EinsteinTensor
CoordinateChart
HilbertSpace
WaveFunction
GaugeField
Field
Spinor
...
```

as kernel kinds merely because GR needs them.

Likewise, we should not add:

```text
semantic subtyping
entity identity
SMT solving
general logical entailment
general calculus
tensor algebra
```

to `internal/kernel` merely because the GR implementation is inconvenient without them.

The reviewers consistently endorse that boundary. plan11_gemini plan11_qwen

---

# 5. GR is now explicitly a workload, not a redesign exercise

The first GR implementation should therefore be treated as:

> **A hostile workload against the existing substrate.**

Not:

> "Let's build the GR features the kernel appears to lack."

The workflow should be:

```text
GR requirement
      ↓
Can general_relativity express it itself?
      |
      ├─ YES → Level 1
      |
      └─ NO
           ↓
Can a generic shared library express it?
           |
           ├─ YES → Level 2 candidate
           |
           └─ NO
                ↓
Can a provisional Level-2 abstraction be built
without changing kernel semantics?
                |
                ├─ YES → remain above kernel
                |
                └─ NO
                     ↓
                investigate Level 3
```

And even then, **Level 3 is not automatically "grow kernel."**

The candidate must pass the Growth Gate.

---

# 6. Every GR failure gets classified

This is one of the most valuable additions from the adversarial review.

For every obstacle encountered, the GR implementation log should classify it as exactly one of:

```text
SPEC-INTENDED-BOUND
    MVP deliberately does not support it.

PACKAGE-SOLVABLE
    general_relativity can handle it locally.

REPRESENTABLE-BUT-UNFAITHFUL
    current substrate can encode it, but loses important semantics.

UNREPRESENTABLE
    current substrate cannot faithfully encode it.

SILENTLY-WRONG
    current substrate accepts it and produces a result whose
    semantic error is not mechanically exposed.
```

Only the last three constitute serious evidence requiring deeper architectural investigation.

And even `UNREPRESENTABLE` is **not** evidence of kernel growth until a higher-level attempt has failed.

That is the major methodological improvement from the consolidated reviews.

---

# 7. Add the Level-2 governance pins now, but build nothing yet

The reviews identified four governance requirements for the new shared-library tier. plan11_z

### Level 1 → Level 2

Requires:

```text
named second theory
+
worked second-theory example
+
demonstrated common abstraction
```

"GR and QM might both need tensors" is not enough.

### Level 2 home

Before actually creating shared libraries, decide their repository/module boundary. They should not silently violate the frozen MVP 42-file closed-world tree.

### Level-2 non-leakage

A shared library must preserve the same provenance, assumption, canonicalization, and integrity discipline. Moving something to Level 2 cannot become a loophole around the kernel.

### Level-2 no-bypass rule

A shared library cannot silently depend on a kernel capability that does not exist.

If it genuinely requires a new trusted invariant:

```text
Level 2 attempt
      ↓
document failure
      ↓
Growth Gate
      ↓
possible Level 3 candidate
```

---

# 8. The concern ledger should now be frozen

This prevents future agents from repeatedly reopening settled arguments.

```text
C1  Expression type decay
    → Level 2 semantic/ascription investigation

C2  Derived-expression → nominal-category ascription
    → Level 2 semantic/ascription investigation

C3  Semantic assumption contradiction
    → Level 2 constraint/entailment investigation

C4  Symbol/entity identity
    → Level 2 namespaced symbol investigation

C5  Transcendentals
    → future extension, not MVP defect

C6  Symbolic exponents
    → future extension, not MVP defect

C7  Negative-integer differentiation
    → explicitly unsupported MVP case

C8  Identify result kind
    → KindRelation, no retyping
```

The purpose is not merely documentation. It prevents an AI coding agent from seeing a familiar problem during GR and immediately proposing "add it to the kernel."

---

# 9. What we should NOT do before GR

I would **not** build any of these now:

```text
semantic/
symbols/
constraints/
tensor/
manifold/
calculus/
hilbert/
```

The point of the stress test is precisely to discover which abstractions are actually required.

The only pre-GR work should be the three mechanical clarifications above plus the Growth Gate instrumentation/governance.

The reviews explicitly converge on this: the next step is to run GR and let it generate evidence rather than pre-solving the anticipated failures. plan11_deepseek

---

# Final consolidated state

I would lock this as the starting position:

```text
KERNEL
  robust
  trusted
  math-agnostic
  theory-agnostic
  intentionally closed
  no growth established

MVP OPERATIONS
  bounded
  exact
  deterministic
  fail-closed where unsupported
  not a CAS
  not a theorem prover

LEVEL 1
  theories implement their own mathematics as needed

LEVEL 2
  emerges only from demonstrated cross-theory reuse
  human-curated
  remains above kernel

LEVEL 3
  only generic trusted primitives/invariants
  requires full Growth Gate evidence

GR
  first serious workload
  no preselected kernel additions
  every failure classified
```

The most important conclusion is:

> **We are now testing whether the kernel is expressive enough, not trying to make the kernel expressive enough for GR.**

That is the correct stress test for the Go-inspired architecture. The consolidated reviews support this direction, with the `Identify` result-kind ambiguity and negative-integer differentiation wording resolved before the experiment starts. plan11_z

**I would consider this the pre-GR lock.** The next step should be the actual `general_relativity` workload plan, starting from the simplest faithful GR representation and carrying the failure-level ledger throughout.
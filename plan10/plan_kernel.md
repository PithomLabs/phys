## prompt

now just what we did with PithomLabs/solvent I want you to identify the scenarios, gates, constraints, requirements needed for kernel growth gate, the relativity package is only about special relativity, one target scenario for the kernel growth gate design is about general relativity, we only add primitives to the kernel for established facts like general relativity, any theorem like equivariance theorem, Valentini H-Theorem should be outside the kernel as part of the user's libraries of packages, got it? also, refresh me about the metadata design for each package. I looked at relativity package and it does not have a brief README.md for AI agents so that they need not grok the code and other metadata that we briefly discussed earlier


## answer

Yes, got it. And the Solvent precedent is useful here: the **Growth Gate should be a narrow admission mechanism, not permission to continually reopen the trusted core**. In Solvent, the rule was essentially “new kernel behavior only when a genuinely new durable invariant/transition cannot be expressed outside the kernel; otherwise use adapters/domain packs.” solvent-agentjacking-v0-vision-…

For Physics Compiler, I would translate that into an epistemic/formal rule.

## 1. The key distinction

We should separate three things:

```text
GENERIC KERNEL
    = universal formal machinery

ESTABLISHED THEORY PACKAGE
    = trusted, bounded physical knowledge

USER LIBRARY / HYPOTHESIS PACKAGE
    = theorems, models, interpretations, proposed structures
```

So:

```text
internal/kernel/
    must remain theory-neutral

relativity/
    Special Relativity only

general_relativity/
    General Relativity only

user libraries
    equivariance theorem
    Valentini H-Theorem
    derivations
    specialized results
    competing formulations
    research artifacts
```

The important correction to the phrase “add primitives to the kernel for established facts like general relativity” is: **GR primitives should become part of an established `general_relativity` corpus package, not automatically part of `internal/kernel`.** The kernel grows only when the *generic formal substrate itself* is missing something fundamental. The existing architecture already establishes that theory packages own their objects, relations, assumptions, conventions, scope, limitations, anomalies, derivations and corpus metadata. plan7(5)

That distinction will prevent the kernel from eventually becoming a giant physics ontology.

---

# 2. Kernel Growth Gate — scenarios

I would define these scenarios for the gate.

### Scenario A — New established physics primitive

Example:

> Add a General Relativity primitive such as spacetime metric / curvature-related quantity.

Gate asks:

```text
Is this a generic mathematical/semantic primitive
needed by the formal substrate?

        YES → potentially kernel

        NO → established theory package
```

For GR specifically, the default outcome should be:

```text
general_relativity/
    primitive
    relation
    assumptions
    metadata
```

not `internal/kernel`.

### Scenario B — New established mathematical capability

Example:

> A future established theory requires a generic tensor/index object that the current formal substrate literally cannot represent.

This can be a **kernel-growth candidate**.

But it must demonstrate:

```text
existing representation insufficient
AND
capability is theory-neutral
AND
multiple future frameworks can use it
AND
cannot safely live in a package
```

### Scenario C — Established theorem

Example:

```text
equivariance theorem
Valentini H-Theorem
```

Outcome:

```text
NOT KERNEL

→ user/theory library
```

Even if the theorem is mathematically established, it is still a **derived proposition**, not foundational substrate machinery.

### Scenario D — Derived relation within a framework

Example:

```text
E = mc²
```

The existing `MassEnergyRelation` may remain a corpus artifact in `relativity`, while the actual derivation is performed by the generic operation/session machinery. The specification explicitly distinguishes the reusable corpus relation from the derivation that must not simply retrieve it. physics_compiler_mvp_specs_v2

### Scenario E — New hypothesis

Outcome:

```text
hypothesis/
    → provisional object
    → HYPOTHESIS
    → never kernel
```

The specification already establishes that candidate-derived results remain HYPOTHESIS and cannot automatically enter trusted corpus authority. specs_v2_3

---

# 3. Kernel Growth Gate — gates

I would make the admission process roughly:

```text
G0  Scope / request classification
G1  Established-status verification
G2  Primitive-vs-theorem classification
G3  Genericity test
G4  Package-isolation test
G5  Existing-capability sufficiency test
G6  Assumption/scope/falsifiability test
G7  Canonical representation test
G8  Cross-framework contamination test
G9  Regression/mutation test
G10 Human freeze decision
```

### The crucial gate

**G3 — Genericity**

Ask:

> Does this belong to the universal formal language of the library, or does it belong to one physical framework?

If it contains theory-specific assumptions, laws, postulates, interpretations or empirical scope, it belongs outside the generic kernel.

---

# 4. Hard constraints

A proposed kernel addition must satisfy all of these:

```text
1. established or formally foundational for the substrate
2. theory-neutral
3. reusable across multiple frameworks
4. cannot reasonably live in an existing package
5. does not encode a specific theory's assumptions
6. does not embed a theorem/derivation as primitive machinery
7. does not create truth adjudication
8. does not create automatic promotion
9. does not weaken provenance containment
10. does not introduce hidden ambient assumptions
11. has deterministic canonical representation
12. has adversarial regression coverage
```

And the default should be:

> **When uncertain, keep it outside the kernel.**

That mirrors the Solvent growth philosophy: extensions should normally live in the surrounding layers, and kernel expansion should be rare and deliberate. solvent-agentjacking-v0-vision-…

---

# 5. General Relativity as the target stress test

GR is actually an excellent Growth Gate test because it is much larger than the current SR package.

A request such as:

> “Add General Relativity.”

must **not** translate into:

```text
put GR into internal/kernel
```

Instead the gate should decompose it:

```text
General Relativity
      │
      ├── established primitives
      ├── framework assumptions
      ├── conventions
      ├── relations
      ├── limits
      ├── anomalies
      ├── derivations
      └── falsifiability metadata
              ↓
      general_relativity/
```

Only if GR exposes a genuinely missing **generic formal capability** should the gate create a kernel-growth proposal for that capability.

This gives us a very powerful invariant:

> **Adding a theory must not require turning the generic kernel into that theory.**

---

# 6. Package metadata — the design we already have

You remembered correctly that we discussed packages as the AI's physics corpus rather than just ordinary Go APIs.

The authoritative machine-readable layer is currently `manifest.json`. The specification defines package-level metadata as:

```text
schema_version
framework_id
framework_name
corpus_status
assumptions[]
domain[]
limits[]
anomalies[]
items[]
```

and each item carries:

```text
id
constructor
kind
name
statement
canonical_expr
dimension
provenance_status
source
assumptions
derivable_from
reduces_to
known_limits
anomalies
falsification_conditions
```

specs_v2_3

The manifest is the **canonical semantic source**; Go comments can mirror it but are not authoritative. It is typed, validated, embedded, and constructor-cross-checked. physics_compiler_mvp_specs_v2

So I would think of package metadata as three layers:

```text
README.md
    ↓
human/AI orientation

manifest.json
    ↓
canonical machine-readable semantics

Go package
    ↓
executable constructors + operations
```

---

# 7. The missing `relativity/README.md`

You're right: this is currently a hole.

The manifest tells an agent **what exists**, but a short README should tell an AI agent **what this package is before it reads the code**.

I would make every future established framework package have a very small AI-facing README with approximately:

```text
# Special Relativity

Framework ID: special_relativity
Status: ESTABLISHED

Scope:
    Special relativity only.

Not included:
    General relativity
    gravitational dynamics
    quantum mechanics

Core assumptions:
    Minkowski spacetime
    Lorentz symmetry
    special-relativistic regime
    no gravitational dynamics

Convention:
    metric signature -+++

Core primitives:
    Spacetime
    MinkowskiMetric
    RestMass
    Energy
    ThreeMomentum
    FourMomentum
    SpeedOfLight
    Velocity

Core relations:
    LorentzFactor
    EnergyMomentumRelation
    MassEnergyRelation

Canonical derivation target:
    E = mc²

Important limitation:
    This package does not implement gravitational dynamics.

Machine-readable metadata:
    manifest.json
```

The current specification already requires the relativity package to identify itself as `special_relativity`, carry its framework assumptions and limits, and contain the relevant corpus metadata. physics_compiler_mvp_specs_v2

### One important freeze implication

I would **not modify the frozen 39-file MVP merely to add this README now**.

Instead, make:

```text
Package README requirement
    =
post-freeze Growth Gate requirement
```

Then the Growth Gate can add `relativity/README.md`, `mechanics/README.md`, and eventually `general_relativity/README.md` as **documentation-layer package metadata**, without confusing that with kernel growth.

That gives us a clean rule:

```text
kernel growth ≠ corpus growth ≠ package documentation growth
```

And I think that distinction is exactly what we need before we start adding General Relativity.
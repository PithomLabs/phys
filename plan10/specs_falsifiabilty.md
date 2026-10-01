## prompt

one more constraint, given that we only encode established physics facts, for this special relativity exercise it must be construed as external package in a separate folder outside the kernel since we can safely say that special relativity although proven already must be treated with the same Karl Popper criterion of falsifiability per EBP 2.1, that hypothesis or theorem or anything with its set of assumptions or postulates should be treated as separate package with OS folder as our unit of separation to separate other hypotheses, what do you think? if this makes sense, include additional prompt to enforce this


## answer

Yes. I think this is a **good architectural constraint**, with one important refinement: treat it as an **epistemic/package isolation rule**, not as a claim that established special relativity is merely an unverified hypothesis.

Your EBP 2.1 source frames falsifiability as an eval requirement and emphasizes that the protocol is an “honest-accounting” framework rather than a truth engine. :chatgpt-content-reference{index="0"} That fits very well with separating a physical framework from the trusted semantic kernel.

The key distinction should be:

```text
internal/kernel
    = formal machinery
    = no physical theory
    = no theory-specific assumptions/postulates/laws

mechanics/
    = classical-mechanics corpus

relativity/
    = special-relativity corpus
    = its own assumptions, conventions, postulates/framework scope

hypothesis/
    = provisional/untrusted candidate space
    = never trusted merely because it is structurally valid
```

So the kernel does **not** say "special relativity is true." It provides machinery capable of representing and transforming propositions under explicit assumptions.

This also gives you a strong future scaling rule:

```text
ONE PHYSICAL FRAMEWORK / THEORY
        ↓
ONE TOP-LEVEL OS FOLDER / PACKAGE
        ↓
OWN ASSUMPTIONS
OWN POSTULATES
OWN CORPUS
OWN DERIVATIONS
OWN LIMITATIONS
OWN ANOMALIES
```

Then a future theory can be evaluated **against** another theory without their axioms silently becoming one undifferentiated trusted corpus.

That is particularly useful for your eventual TOE work: a proposed theory should not gain special-relativity premises merely because both happen to live in the same compiler.

### One correction to the terminology

I would avoid saying:

> "special relativity must be treated as a hypothesis."

Instead:

> **Special relativity is an established physical framework whose propositions remain explicitly assumption- and falsifiability-bounded.**

That preserves the distinction between **epistemic status** and **architectural isolation**.

The current v2.3 design already has the correct physical separation mechanically: `relativity/` is outside `internal/kernel`, and the specification explicitly says the kernel should not contain the physical corpus. The specification also explicitly forbids automatic truth adjudication and automatic corpus-status inference. :chatgpt-content-reference{index="1"}

## Add this to the implementation-agent prompt

I would append the following as **FREEZE GATE H — THEORY / EPISTEMIC PACKAGE ISOLATION**:

---

# 18. FREEZE GATE H — THEORY / EPISTEMIC PACKAGE ISOLATION

Enforce a strict separation between the generic formal kernel and physical theory/framework corpora.

The architectural principle is:

```text
internal/kernel
    = theory-neutral semantic machinery

top-level physics package
    = one bounded physical framework/corpus
```

For the current MVP:

```text
internal/kernel/
    = no special-relativity corpus
    = no mechanics corpus
    = no physical laws
    = no theory-specific assumptions/postulates

mechanics/
    = classical-mechanics corpus

relativity/
    = special-relativity corpus

hypothesis/
    = provisional/untrusted candidate space
```

### A. Special relativity remains outside the kernel

Treat `relativity/` as an external theory/framework package relative to the kernel.

Its:

- physical kinds
- relations
- formulas
- assumptions
- conventions
- framework scope
- anomalies/limitations
- derivations
- corpus metadata

must remain outside `internal/kernel`.

Do not move special-relativity knowledge into the kernel merely because the kernel needs it for the E=mc² demonstration.

The kernel may contain only the explicitly specification-mandated generic identifiers already permitted, including:

```text
KindMinkowski
LorentzFactorFunctionID = "lorentz_factor"
```

Those identifiers are semantic hooks mandated by the frozen specification; they do not authorize embedding the corresponding physical theory into the kernel.

### B. OS-folder/package is the unit of physical-theory separation

Treat the top-level package/folder as the architectural isolation boundary for a physical framework.

Future physical frameworks must not be merged into a common physics package merely for convenience.

Conceptually:

```text
mechanics/
relativity/
future_theory/
future_hypothesis/
```

Each framework owns its own:

```text
assumptions
postulates
corpus objects
relations
scope
limitations
anomalies
derivations
```

A framework must not silently inherit another framework's assumptions.

Cross-framework reasoning, when eventually introduced, must occur through explicit objects/relations and explicit assumptions rather than package-global ambient state.

### C. Do not interpret established status as unrestricted truth

`ESTABLISHED` means trusted corpus status within the current curated corpus.

It does not mean:

```text
universally true
unfalsifiable
outside empirical scrutiny
axiom of the entire compiler
```

Physical claims remain associated with their explicit:

```text
assumptions
framework
scope
limitations
falsifiability
```

This is consistent with the EBP principle that the system is an honest-accounting/evaluation substrate rather than a truth engine.

Do not implement EBP 2.1 inside the kernel as part of this freeze.

### D. Hypothesis isolation

The `hypothesis/` package must remain separate from trusted physical-theory packages.

A provisional hypothesis:

```text
must not become a member of relativity/
must not become a member of mechanics/
must not enter internal/kernel/
must not acquire ESTABLISHED status automatically
```

Hypothesis-derived expressions remain subject to the existing HYPOTHESIS-contamination rules.

A valid mathematical derivation does not promote a hypothesis into trusted corpus authority.

### E. No ambient assumption leakage

Add regression tests ensuring that operations cannot obtain physical assumptions merely from package co-location.

For example:

```text
relativity assumptions
    must be explicitly present in relativity-derived objects

mechanics assumptions
    must be explicitly present in mechanics-derived objects

hypothesis assumptions
    must remain hypothesis-associated
```

No package may rely on a hidden global assumption registry to make an expression valid.

### F. Package-boundary regression test

Add a structural/import-boundary test ensuring:

```text
internal/kernel
    does not import
        mechanics
        relativity
        hypothesis
```

and preserve the already-required dependency direction.

The domain packages may depend on the generic formal substrate, but the generic substrate must not depend on a particular physical theory.

### G. Freeze constraint

Do not create a new directory structure during this hardening pass.

The existing top-level separation is sufficient for MVP:

```text
mechanics/
relativity/
hypothesis/
internal/kernel/
```

This gate establishes the architectural rule for future corpus growth without expanding the current MVP.

---

### Why I would add this now

It strengthens the exact thing you're trying to prove later.

The AI should eventually be able to reason like:

```text
[Special Relativity corpus]
        +
[another theory corpus]
        +
[explicit bridge assumptions]
        ↓
      AI reasoning
        ↓
   candidate derivation
        ↓
   EBP-style evaluation
```

rather than:

```text
everything physics-related
        ↓
one giant trusted knowledge pool
        ↓
AI mixes assumptions invisibly
```

That distinction will become **very important once you start doing TOE/hypothesis work**. The E=mc² demonstration is actually an excellent reason to enforce it now: `EnergyMomentumRelation` and `RestFrameAssumption` belong to the relativity package; the kernel merely supplies the machinery that lets the agent manipulate them. :chatgpt-content-reference{index="2"}

I would therefore add **Gate H** to the implementation prompt, while keeping the existing 39-file tree and not adding any new theory packages during this freeze.
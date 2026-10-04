I reviewed **Plan 8 — Freeze Delta** against the entire chain of decisions and the latest adversarial review. It is now the strongest consolidated version so far. It correctly incorporates G1–G11, preserves Plan 6/7, keeps the 39-file tree, and explicitly defers Einstein-1905 corpus work. :chatgpt-content-reference{index="0"}

I would approve it **with three small wording/precision amendments before implementation**.

### 1. H2 is correct, but make the allowlist structural rather than textual

Plan 8 correctly fixed the earlier H2 mistake by allowing the specification-mandated enums and the two distinctive identifiers. :chatgpt-content-reference{index="1"}

However, the implementation agent should not interpret this as:

> “Any occurrence of these strings anywhere in `internal/kernel` is automatically allowed.”

It should verify the **specific mandated declarations/usages**:

```text
Kind enum
ExprKind enum
AssumptionKind enum
ProvenanceStatus enum
CorpusStatus enum
KindMinkowski
LorentzFactorFunctionID
closed-world validation required by those definitions
```

Everything outside those mandated structures remains forbidden.

Likewise, the Lorentz exception should be function-scoped exactly as Plan 8 says, because `"v"` and `"c"` are legitimate inside the fixed Lorentz body but must not become generic symbol-dispatch logic. :chatgpt-content-reference{index="2"}

### 2. H3 needs both key and value preservation

Plan 8 correctly incorporates the review's stronger assumption-preservation rule:

> inherited assumptions are preserved **verbatim**, not merely by key. :chatgpt-content-reference{index="3"}

That is exactly right.

The implementation invariant should effectively be:

```text
output assumptions
    =
exact union(input assumptions)
+
explicit operation-generated assumptions
```

with identical `(Kind, Key, Value)` preserved and conflicting values rejected.

That aligns with the specification's deterministic merge law. :chatgpt-content-reference{index="4"}

### 3. H6 should explicitly test all four `OperationParams.Kind` values

Plan 8 correctly superseded the old `pow`-only rule with symmetric per-kind validation. :chatgpt-content-reference{index="5"}

The four schema kinds are:

```text
empty
pow
compare
identify
```

The implementation should explicitly verify:

```text
empty
    → Exponent=""
      Operator=""
      Justification=""

pow
    → Exponent required
      Operator=""
      Justification=""

compare
    → Operator required
      Exponent=""
      Justification=""

identify
    → Justification required
      Exponent=""
      Operator=""
```

The specification confirms this exact four-kind schema. :chatgpt-content-reference{index="6"}

And ordinary operation fixtures should use the canonical empty form; the old `"operator":"eq"` `pow` representation belongs only in the negative test.

---

## Everything else in Plan 8 is net-valid

The following decisions are now properly locked:

**Theory isolation:** special relativity remains `ESTABLISHED`, but its assumptions, conventions, scope and empirical status remain bounded inside `relativity/`; they do not become ambient kernel axioms. :chatgpt-content-reference{index="7"}

**Reverse constructor allowlist:** explicit and closed, rather than inferred from status. :chatgpt-content-reference{index="8"}

**Package leakage:** constructor-output checks, not exact assumption-set pinning, so legitimate mixed-framework derivations remain possible. :chatgpt-content-reference{index="9"}

**Corpus status:** schema remains capable of `CONTESTED`, etc.; the freeze test pins the current baseline to `ESTABLISHED` rather than making alternative statuses illegal. :chatgpt-content-reference{index="10"}

**Source provenance:** both corpora are aligned to the Plan 10 pins. :chatgpt-content-reference{index="11"}

**README epistemic wording:** correctly distinguishes the formal E=mc² derivation from Einstein's 1905 reconstruction and from empirical proof. :chatgpt-content-reference{index="12"}

**Mutation methodology:** disposable copy, restore per mutant, `-count=1`, and claim only actual rerun results. :chatgpt-content-reference{index="13"}

**Non-blockers:** crypto authenticity, hostile `unsafe`/`reflect`/`linkname`, general fuzzing, EBP runtime, and Einstein-1905 corpus remain correctly outside this freeze. :chatgpt-content-reference{index="14"} :chatgpt-content-reference{index="15"}

## My final verdict

**Plan 8 is the correct implementation plan for the final kernel freeze.**

Make the three precision amendments above, then give it to the coding agent. Do **not** add another architectural phase.

The resulting freeze boundary is now very clear:

```text
internal/kernel
    = generic formal machinery

mechanics/
    = established classical-mechanics framework

relativity/
    = established special-relativity framework,
      explicitly bounded by its assumptions/scope

hypothesis/
    = provisional candidates

AI agent
    = composes reasoning/derivations

session/
    = records and audits the derivation

E=mc²
    = primary demonstration that the agent can compose
      a non-hardcoded derivation from encoded SR premises
```

And the 1905 Einstein reconstruction stays exactly where it belongs: **after the kernel is frozen**.
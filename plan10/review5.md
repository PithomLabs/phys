I consolidated all four reviews against the completed software audit and the completed physics-corpus audit. The overall picture is very strong, but I would make **one final, narrow hardening pass before declaring the kernel frozen**.

The important distinction is that the current implementation has already passed the substantive software and corpus audits: 18/18 mutation kills, full replay, candidate containment, relativity/mechanics encoding, exact corpus cross-checks, 39-file architecture, and clean-copy verification. :chatgpt-content-reference{index="0"} :chatgpt-content-reference{index="1"}

The remaining concerns are mostly about **closing trust-boundary gaps mechanically**, rather than redesigning anything.

## Consolidated findings

### P1 — Make malformed AI/API inputs impossible to turn into panics

Gemini identified three concrete robustness gaps: `Limit` receiving a non-symbol variable, `SelectBranch` relying too much on structural/index assumptions, and `ops.Apply` accessing positional inputs without an explicit arity guard. :chatgpt-content-reference{index="2"}

These are worth retaining even though they were not shown to be current failures.

The frozen rule should be:

```text
malformed operation input
        ↓
structured error
        ↓
never panic
```

In particular:

- every operation has an explicit required arity;
- `Apply` checks arity before indexing;
- `Substitute`, `Differentiate`, and `Limit` require the relevant variable/wrt object to be a single Symbol;
- `SelectBranch` validates the complete BranchSet shape before selecting.

This is particularly important for an AI-facing library because malformed generated calls are an expected operating condition.

### P1 — Test the provenance law across **all 12 pure operations**

DeepSeek's finding is net-valid: the plan's Test L listed only seven operations even though v2.3 defines the HYPOTHESIS propagation law for all twelve. :chatgpt-content-reference{index="3"}

The current implementation audit did not establish that this complete 12-operation matrix is mechanically tested.

Freeze criterion:

```text
Add
Subtract
Multiply
Divide
Pow
Simplify
Substitute
Differentiate
Limit
Compare
Solve
SelectBranch

+ Session.Identify
```

For every operation, verify:

```text
any HYPOTHESIS input
→ HYPOTHESIS output
```

and:

```text
no HYPOTHESIS input
→ DERIVED output
```

with Identify separately producing `IDENTIFIED` only when permitted.

This is a core epistemic invariant, not just a test nicety. The specification explicitly defines this propagation law. :chatgpt-content-reference{index="4"}

### P1 — Lock the **corpus-status content**, not merely its enum validity

This is, in my view, the most important new finding from Z's self-review.

The corpus audit says corpus status is human-curated and loaded from manifests, never computed. :chatgpt-content-reference{index="5"}

But Z correctly observes that simply validating:

```text
ESTABLISHED
CONTESTED
SUPERSEDED
FALSIFIED
NONE
```

does not prevent a valid-value corruption such as:

```text
ESTABLISHED → CONTESTED
```

from passing schema validation. :chatgpt-content-reference{index="6"}

For a system whose purpose is to distinguish **established corpus knowledge from research/hypothesis material**, this deserves a hard guard.

The manifest verification should explicitly assert the baseline corpus status:

```text
mechanics              = ESTABLISHED
special_relativity     = ESTABLISHED
```

and verify that candidate/derived artifacts are not silently promoted.

### P1 — Close the reverse corpus-constructor loophole

Z also identifies a real asymmetry:

```text
manifest → constructor
```

is checked, but:

```text
constructor → manifest
```

may not be.

An extra exported constructor could theoretically create an `ESTABLISHED` physics object that has no manifest entry and therefore evade manifest-to-constructor checking. :chatgpt-content-reference{index="7"}

For the user's intended model:

> **all trusted established physics belongs to the explicitly curated corpus**

that reverse direction is important.

The final mechanical test should establish a bijection/allowlist:

```text
manifest constructors
+
explicitly permitted non-manifest constructors

= all trusted physics constructors
```

Anything else producing an `ESTABLISHED` corpus object should fail the audit.

### P1 — Mechanically enforce the kernel's physics boundary

This is the one point where the four agents' discussions converge with an important correction.

Qwen/DeepSeek initially framed the requirement as “no physics in kernel,” but the corpus audit established that the spec itself mandates limited physics identifiers inside the kernel:

```text
KindMinkowski
lorentz_factor
```

and their associated closed-world validation. :chatgpt-content-reference{index="8"}

So the correct freeze rule is:

> **`internal/kernel` may contain only explicitly enumerated spec-mandated physics identifiers; no unmandated physics law, symbol semantics, framework assumption, or domain-specific rule may be added.**

Likewise in `ops`:

```text
generic symbolic machinery
      +
explicitly mandated Lorentz-function body
```

but no:

```text
if symbol == "c" ...
if symbol == "m" ...
```

The audit found the current implementation clean in this respect. :chatgpt-content-reference{index="9"}

For long-term kernel hardening, I agree with DeepSeek that this boundary should eventually have a mechanical source audit, but **use an allowlist, not a blanket “no physics symbol” scan**. :chatgpt-content-reference{index="10"}

### P2 — Validate falsifiability expressions as real `core.Expr`

Gemini's candidate finding is worth incorporating.

The candidate model explicitly uses structured `core.Expr` for Prediction, FalsificationCondition, and RecoveryClaim, and candidate validation includes falsifiability/anomaly validation. 

Therefore a trusted sealed candidate should never contain an invalid expression handle.

Add:

```text
Prediction.Relation.Valid()
FalsificationCondition.ContradictingCondition.Valid()
RecoveryClaim.Condition.Valid()
```

before trusted candidate construction.

Do **not** require these expressions to be scientifically true or related to the hypothesis; that would violate the research-artifact boundary. This is purely structural validity.

### P2 — Run the manifest corruption battery as a real mutation test

DeepSeek initially called this a major gap; Z later correctly noted that the existing machinery detects many of the corruptions. :chatgpt-content-reference{index="12"}

The remaining lesson is stronger: **analytical reasoning about what a test should catch is not equivalent to actually mutating the corpus and observing detection.**

The final hardening suite should actually mutate, in a disposable copy:

```text
F=ma
p=mv
K
E² relation
m≥0
c>0
Minkowski convention
RestFrame
provenance
corpus_status
dimension
extra trusted constructor/relation
```

and observe the intended detector.

That also directly addresses Z's R2/R3 observations. :chatgpt-content-reference{index="13"}

### P2 — Freeze the ambiguous serialization/plan pins explicitly

Z identifies two genuine documentation ambiguities:

- the `OperationParams` example uses `operator:"eq"` for a `pow` example even though unused parameters are supposed to be empty;
- several deterministic corpus values are pinned in the plan but not clearly separated from normative spec requirements. :chatgpt-content-reference{index="14"}

The `OperationParams` issue is especially worth recording because the v2.3 text itself currently contains the contradictory example. :chatgpt-content-reference{index="15"}

I would freeze the operational rule as:

```json
{
  "kind": "pow",
  "exponent": "2/1",
  "operator": "",
  "justification": ""
}
```

and treat that as an implementation pin/erratum rather than modifying the frozen specification.

### P3 — Rational serialization deserves one explicit invariant test

Qwen's `big.Rat` warning is technically useful, but I would **not change the implementation merely because of it**. The current corpus audit already verified deterministic canonical expressions and hashes. :chatgpt-content-reference{index="16"}

The useful hardening is simply to ensure a test explicitly asserts:

```text
1/2 → "1/2"
-3/4 → "-3/4"
2 → "2/1"
0 → "0/1"
```

and rejects JSON numeric tokens.

No need for a new abstraction unless the current source inspection shows `encoding/json` is directly marshaling `*big.Rat`.

---

## Findings I would **not** carry forward

### “Identify justification should be semantically meaningful”

Reject this as an MVP requirement.

The spec requires a **non-empty trimmed justification**, not a human-quality assessment. Gemini's proposal to have consumers generate a Challenge for a meaningless justification crosses into review/orchestration territory. :chatgpt-content-reference{index="17"}

The kernel should not decide whether `"a"` is a good scientific justification.

### “go:embed must only exist in test files”

Not a normative blocker.

The v2.3 spec requires the manifest bytes to use `go:embed`, but does not require the directive specifically to live in `manifest_test.go`. :chatgpt-content-reference{index="18"}

The current implementation's test-only placement is a reasonable architecture choice, but don't turn this into a semantic kernel invariant.

### “No physics whatsoever in internal/kernel”

Too strong.

The current audit correctly found `lorentz_factor` and `KindMinkowski` there, and those are explicitly mandated by the closed-world design. :chatgpt-content-reference{index="19"}

The right boundary is **no unmandated physics**.

### “Cryptographic authenticity”

Already an accepted MVP limitation, not a missing feature. The specification explicitly describes the ledger as an integrity mechanism, not an authenticity signature. :chatgpt-content-reference{index="20"}

The documentation should keep saying so.

---

# The consolidated freeze checklist

I would require these **seven gates** before declaring the kernel truly frozen:

| Gate | Freeze requirement |
|---|---|
| **1. No-panic boundary** | malformed operation inputs produce typed errors, never slice/index panics |
| **2. Provenance law** | all 12 pure ops + Identify explicitly tested for HYPOTHESIS/DERIVED/IDENTIFIED propagation |
| **3. Corpus authority** | corpus status is content-locked; no valid-value drift; reverse constructor→manifest check |
| **4. Kernel boundary** | only explicitly allowed physics identifiers exist in kernel; no hidden physics laws |
| **5. Trusted candidate integrity** | all embedded candidate expressions are valid before trusted construction |
| **6. Corpus mutation proof** | actual corruption battery detects semantic corpus mutations |
| **7. Serialization pins** | exact OperationParams/rational/corpus pins are explicitly frozen and internally consistent |

Everything else is either already strongly established by the independent audits or belongs to the next layer.

## Final architecture I would freeze

```text
                 ┌───────────────────────────┐
                 │       internal/kernel     │
                 │                           │
                 │ generic symbolic substrate│
                 │ dimensions                │
                 │ assumptions               │
                 │ provenance                │
                 │ canonicalization          │
                 │ mint authority            │
                 │ + only mandated closed-   │
                 │   world physics IDs       │
                 └─────────────┬─────────────┘
                               │
                            core
                               │
                    ┌──────────┴──────────┐
                    │                     │
                 generic ops        physics corpus
                                      │
                              ┌───────┴───────┐
                              │               │
                         mechanics        relativity
```

Then:

```text
established physics
        ↓
mechanics / relativity corpus

new hypothesis
        ↓
external research artifact
        ↓
session / hypothesis
        ↓
ResearchCandidate
```

That separation is exactly what the corpus audit ultimately demonstrated: all required mechanics and special-relativity knowledge is explicitly encoded in the designated corpus layers, while the kernel contains no Newtonian law, GR, QM, EM, thermodynamics, or other hidden physics. :chatgpt-content-reference{index="21"} :chatgpt-content-reference{index="22"}

**My consolidated verdict: the MVP is functionally and physically well-founded enough to freeze, but I would not make the freeze irreversible until the seven hardening gates above are mechanically demonstrated. None requires re-architecting the system.**
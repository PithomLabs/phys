I consolidated the attached **Gemini + Z adversarial review** with our prior freeze verdict and the decisions we have since locked. The review is useful, but several of its strongest claims are either over-broad or already superseded by decisions we made afterward. The result is a tighter freeze plan, not a redesign.

## Consolidated verdict

### Net-valid and should be incorporated

| Finding | Verdict | Freeze action |
|---|---|---|
| **G1 — H2 kernel allowlist too literal** | **VALID** | Correct H2 to distinguish *spec-mandated kernel enums/machinery* from theory-specific physics content. |
| **G2 — scanner must ignore comments** | **VALID** | AST/source checks inspect actual syntax, not comments. |
| **G3 — Lorentz `v`/`c` literals need function-scoped allowance** | **VALID** | Permit them only in the named Lorentz-body construction path. |
| **G4 — assumption subset must preserve values, not merely keys** | **VALID** | Test inherited assumptions are preserved verbatim, including canonical values. |
| **G5 — reverse allowlist should be status-independent** | **VALID, with refinement** | All exported **object-producing trusted constructors** must be manifest-backed or explicitly allowlisted. Metadata helpers are not object constructors. |
| **G6 — H4 must be scoped to constructor outputs** | **VALID** | Mixed-framework derivations may legitimately contain both frameworks' assumptions; package-leak tests target package-owned constructor outputs. |
| **G8 — establish precedence** | **VALID** | Add explicit `specs_v2_3.md` > Plan 7 precedence. |
| **G9 — Plan 7 document outside module tree** | **VALID** | Planning document does not count toward the 39-file implementation tree. |
| **G10 — document epistemic rule** | **VALID** | Put the H0 wording in README, including the status of the E=mc² demo. |
| **G11 — resolve inherited source/status pins** | **VALID** | Verify and explicitly freeze the baseline content pins rather than relying only on structural validation. |

The review's summary is correct that these are the meaningful changes needed before implementation. :chatgpt-content-reference{index="0"} :chatgpt-content-reference{index="1"}

### Valid but needs narrower interpretation

**H3 assumption-forgery concern:** the underlying concern is worth hardening, but the review overstates the current architecture. The specification already says assumption merging is deterministic set union plus exact conflict detection, and operation-generated assumptions such as denominator preconditions are explicitly defined. :chatgpt-content-reference{index="2"}

The correct freeze invariant is therefore:

```text
output assumptions =
    input assumptions, preserved verbatim
    +
    explicitly enumerated operation-generated assumptions
```

with:

```text
same key + different value → AssumptionConflictError
```

So we should **enumerate the per-operation generated-key set and test that no other assumptions can spontaneously appear**, rather than introduce a general assumption-trust mechanism.

### Gemini's runtime panic concern

**Valid at the boundary, but not as a demand for a general adversarial parser/fuzzer.**

We already have the right architecture: invalid objects must be rejected before operation-specific work, and public operations must not panic on invalid input. :chatgpt-content-reference{index="3"}

The freeze should therefore require:

```text
malformed structure
    → deterministic error
    → no panic
```

with the targeted malformed-shape battery already in Gate A. We do **not** need to turn the MVP into a generalized malicious-AST sandbox.

### Cryptographic authenticity

**Not a freeze blocker.**

The review is technically correct that byte hashes do not provide cryptographic authenticity against an attacker who is able to rewrite the complete artifact and its hashes. But that is an **authenticity/threat-model problem**, not a reason to redesign the MVP's replay system.

The existing specification deliberately separates replay integrity from truth adjudication and keeps cryptographic authenticity outside this MVP. We retain that boundary.

### AST-evasion via `unsafe`, `reflect`, `go:linkname`

**Reject as a freeze requirement.**

The source-boundary test exists to prove repository/package architecture, not to construct a hostile-process sandbox. We do not need to defend `internal/kernel` against a deliberately malicious Go compiler-level modification using `unsafe` or `go:linkname`.

That would turn a conformance test into a security sandbox, which is outside this MVP.

### G7 — `pow`-only versus symmetric unused-field validation

**Superseded by our locked decision.**

We explicitly chose:

> **Symmetric strict per-kind unused-field rejection.**

So the final rule is not merely:

```text
pow.operator must be empty
```

It is:

```text
for every OperationParams kind:
    meaningful fields = explicitly defined for that kind
    every other field = canonical empty/default
```

This is now the stronger frozen rule and should replace the review's "pow-only" concern.

---

# Important correction to H2

This is the most important technical change to Plan 7.

The review correctly caught that:

> “only `KindMinkowski` and `LorentzFactorFunctionID`”

is too literal and would fail against the actual specification because `internal/kernel` necessarily contains the complete **spec-mandated `Kind` enumeration and related foundational enums**. The specification explicitly requires the 18 physical/object kinds and stable ordinals. :chatgpt-content-reference{index="4"}

So H2 should **not** mean "only these two physics-looking strings may occur in kernel source."

It should mean:

```text
ALLOWED KERNEL PHYSICAL MACHINERY

1. spec-mandated closed Kind enum
2. spec-mandated AssumptionKind enum
3. spec-mandated ProvenanceStatus enum
4. spec-mandated CorpusStatus enum
5. KindMinkowski
6. LorentzFactorFunctionID / "lorentz_factor"
7. closed-world validation required by those specifications
```

And explicitly forbidden:

```text
Newton's law
E=mc²
mass-energy relation
momentum relation
rest-frame physics
Minkowski assumptions
Lorentz symmetry assumptions
special-relativity assumptions
mechanics assumptions
theory-specific formulas
theory-specific constants
symbol-name dispatch masquerading as generic algebra
```

That is much more faithful to the actual architecture.

The same applies to `ops`: the Lorentz implementation is a **small mandated closed-world physics hook**, and the test must permit the entire implementation cluster needed to execute that hook rather than trying to whitelist one textual function body. The review's G1/G3 point is therefore correct. :chatgpt-content-reference{index="5"}

---

# Corpus pin issue: source and corpus status

The source drift finding should be resolved before freeze, but **not by silently changing either side**.

Plan 10 explicitly pinned:

```text
NewtonSecondLaw → "Newton, Principia"
other mechanics items → "Classical Mechanics corpus"

EnergyMomentumRelation,
MassEnergyRelation → "Einstein, 1905"

other relativity items → "Special Relativity corpus"
```

:chatgpt-content-reference{index="6"}

The earlier audit found the implementation currently uses `"Newton, Principia"` for all mechanics items, which is a Plan-vs-implementation pin deviation but not a semantic physics defect. :chatgpt-content-reference{index="7"}

Therefore the freeze action should be:

> **Resolve the source-pin discrepancy explicitly before freeze. Do not merely test manifest-vs-constructor equality.**

Given that Plan 10 is the existing data pin, I would align the implementation to Plan 10 rather than silently redefine the plan.

Likewise for corpus status:

```text
generic schema validation
    → accepts valid status enum values

baseline freeze test
    → requires current mechanics + relativity manifest corpus_status == ESTABLISHED
```

Do **not** make `"CONTESTED"` an invalid schema value; that would conflict with the purpose of the corpus-status enum. The freeze test should pin the current baseline to `ESTABLISHED`, not prohibit future statuses universally.

The specification explicitly says manifest corpus status is **loaded, not derived**. :chatgpt-content-reference{index="8"}

---

# Final consolidated Plan 7 delta

I would now consider these the **authoritative additions/corrections** before implementation:

### H0 — Epistemic wording
Keep:

```text
special relativity = ESTABLISHED
```

and:

> An established physical framework remains explicitly bounded by its assumptions, conventions, scope and empirical status; it does not become an ambient axiom of the generic kernel.

Already correct. :chatgpt-content-reference{index="9"}

### H1 — Import boundary
Keep, but explicitly state it proves **architecture**, not resistance to malicious compiler-level techniques.

### H2 — Kernel boundary
Replace the two-item textual allowlist with a **spec-mandated machinery allowlist**:

```text
closed foundational enums +
KindMinkowski +
LorentzFactorFunctionID +
required closed-world validation
```

and reject actual theory content beyond those mandated mechanisms.

### H3 — Assumption subset
Strengthen from key-only to:

```text
preserve inherited assumptions exactly
+
allow only explicitly enumerated generated assumptions
+
reject unrelated generated assumptions
```

The spec's merge semantics are the foundation. :chatgpt-content-reference{index="10"}

### H4 — Package leakage
Scope checks to **domain constructor outputs**, not arbitrary mixed operations.

Never exact-pin an entire assumption set.

Use forbidden cross-framework keys.

### H5 — E=mc² ownership
Keep exactly as written. The current final result's assumption provenance is already consistent with the intended trace. :chatgpt-content-reference{index="11"}

### New H6 — OperationParams symmetric strictness

Promote our latest locked decision into Plan 7:

```text
each operation kind explicitly defines meaningful params
all unused fields must be empty/default
non-empty unused fields are rejected
```

Update all ordinary success fixtures to canonical form.

Add negative tests for intentionally populated unused fields.

### New H7 — Baseline corpus-content pins

Explicitly verify:

```text
current framework corpus_status = ESTABLISHED
source pins match Plan 10
constructor metadata matches manifest
```

This is a freeze-content test, not a generic schema restriction.

### Documentation

README should explain:

```text
kernel = formal machinery
theory package = bounded established framework
hypothesis = provisional candidate space
```

and explicitly state:

```text
E=mc² MVP demonstration
= formal derivation from encoded SR premises
≠ historical reconstruction of Einstein 1905
≠ empirical discovery/truth adjudication
```

---

## Bottom line

After consolidating the review with our earlier decisions, I would **not reject Plan 7**. I would amend it in a few precise places before implementation:

**G1, G2, G3, G4, G5, G6, G8, G9, G10 and G11 are net-valid.**

**Cryptographic authenticity and hostile `unsafe`/`go:linkname` evasion are not freeze blockers.**

**The generalized runtime-panic concern is handled by the existing no-panic boundary, not by building a general adversarial AST engine.**

**The `pow` issue is superseded by our later locked decision for symmetric per-kind unused-field rejection.**

Most importantly, the revised H2 must distinguish **mandatory semantic/type machinery in the kernel** from **actual physical theory content**. That preserves both requirements: the kernel remains theory-neutral in substance, while still conforming to the specification's required physical-kind and metadata enums. :chatgpt-content-reference{index="12"}

With those changes, Plan 7 is in the right shape for the implementation pass.
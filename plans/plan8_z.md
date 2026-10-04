# Dispatch-Gate Review: v2.1 Prompt + Specs

**Verdict up front: NO architectural blockers remain. GO for handoff after eight small patches (~45 minutes of editing, all single-paragraph or single-line).** This is the first pair in the corpus where I traced the canonical derivation end-to-end through the specified rewrite rules and found it *almost* closes — the residual findings are derivation-level, not architectural, and all have paste-ready fixes below.

---

## 1. What was verified

**All eleven patches from the prior round are applied and correct:**

| Prior patch | Status in v2.1 |
|---|---|
| Canonical outputs retained for replay | ✅ prompt §9 (`InputCanonicals`/`OutputCanonical`), guardrail "Ledger replay" |
| Zero-quantity constructors | ⚠️ partial — `ZeroVelocity()` appears (§15.9.1); see Patch 4 |
| Identify as session method, no ambient state | ✅ §14.6, prompt non-negotiables |
| Fixed `lorentz_factor` body, executed not shortcut | ✅ §15.9.1, including "MUST NOT simply return 1" |
| `Pow` result kind = Expression | ✅ §15.5 |
| Nontrivial differentiation test | ✅ prompt §11 (d(½mv²)/dv) |
| Manifest validator as pure function over bytes | ✅ §19.4 |
| Rational-as-string canonical JSON + round-trip test | ✅ §10.1–10.5 |
| `APPROXIMATED` reserved with coverage exemption | ✅ §11.1, §13.1, REQ-013-01 |
| StepEnvelope struct hashing (no concatenation) | ✅ prompt §9 |
| External test packages for authority negatives | ✅ REQ-003-03 |

**The new `internal/kernel` architecture is correct and necessary.** v2 had a latent contradiction: `core.Object` with unexported fields means only package `core` can mint, but domain packages and `ops` must also mint — with no public factory allowed. v2.1 resolves this with the standard Go pattern (implementation in `internal/kernel`, type alias in `core`, mint authority below the public boundary). I verified: the dependency graph is acyclic; `core.Object = kernel.Object` aliases preserve unexported-field protection for external callers; REQ-004-02/03 hold; `session` can replay without importing domain packages via the §15.13 dispatch mechanism — which answers the prompt's "how does replay avoid an import cycle" question by construction.

**§0.3 (threat-model honesty) is exemplary** — the corpus's own discipline finally applied to itself: constructor authority is a public-API boundary, not hostile-source protection.

**File arithmetic checks out:** §3 layout sums to exactly 39 including `go.mod`, matching the stated count and the 40-file guardrail with one file of headroom.

---

## 2. The derivation trace (evidence of coherence)

I executed the mass-energy chain mentally against the normative rules:

```text
E² = (pc)² + (mc²)²                     Compare: dims equal ✓, Expression operands ✓
Substitute(p → ZeroMomentum)             dims/kind match ✓ (§15.7)
  Mul(0, c) → 0                          §9.6 Mul(x,0) ✓
  Pow(0, 2) → 0                          ⚠️ RULE MISSING (Patch 1)
Simplify                                 Add(x,0) → x ✓
Solve                                    pattern matches §15.11 exactly ✓
  BranchSet(E, [Sqrt(rhs), Neg(Sqrt(rhs))])
SelectBranch(E ≥ 0)                      Energy ∈ ordered kinds ✓ (§6.2)
  Sqrt(Pow(m·c², 2)) → m·c²              §9.8 entailment: m≥0 explicit ✓, Pow(c,2) even ✓ → base nonnegative ✓
Conclude: Relation(eq, E, m·c²)          canonical ✓
```

And the differentiation test: d(½·m·v²)/dv via the n-ary product rule → ½·m·2v → rational combination (§9.4) → `m·v`, dimension Energy/Velocity = M L T⁻¹ ✓, comparable to Momentum under the named+Expression clause ✓.

**The architecture closes. Two mechanical rules are missing** — the only place the trace breaks.

---

## 3. Patches (paste-ready)

**Patch 1 — BLOCKING for acceptance test J. §9.6, add two rules:**

> ```text
> Pow(Rational(1), e) → 1        for any exact rational e
> Pow(Rational(0), e) → 0        when e is a positive rational
> ```
> Both are exact, assumption-free ring identities. Without them, the `Limit(LorentzFactor, v, 0)` path reaches `Pow(1, -1)` after `Sqrt(1) → 1` and cannot canonicalize to `1/1` through the fixed body path that §15.9.1 mandates.

**Patch 2 — BLOCKING for acceptance test D as envisioned. §6.1, add:**

> Test D's equal-dimension, distinct-kind rejection MUST use kinds that exist in the MVP enum — e.g. `Add(KineticEnergy, Energy)` or `Add(Mass, RestMass)`. The §6.1 torque example is illustrative prose; no Torque Kind or constructor exists in MVP, so a torque-based negative test is unconstructible.

**Patch 3 — determinism. §16 (session), add one line:**

> The genesis hash (the `previousHash` of the first committed step) MUST be the fixed constant `SHA256(CanonicalJSON(StepEnvelope{PreviousHash: "", Step: {}}))` — a fixed literal, never generated per-session. Otherwise §32's cross-run derivation-hash determinism fails.

**Patch 4 — enumerate the zero constants. §20/§21 (corpus content), add:**

> The relativity package MUST expose fixed constant constructors `ZeroVelocity()`, `ZeroMomentum()`, `ZeroEnergy()` — Kind matching their named counterpart, expr `Rational(0)`, provenance `DEFINED`, required by §15.9.1, §23.3, and §15.12 respectively.

**Patch 5 — §13.2 wording ambiguity:**

> "Postulate → POSTULATED input required / Define → DEFINED input required" reads as a constraint on inputs. Intended meaning: the `Postulate` action mints `POSTULATED` status; `Define` mints `DEFINED`; `Declare` preserves the input's existing status. Reword accordingly.

**Patch 6 — dependency-diagram footnote (§4):**

> The §4 diagram shows primary authority flow. `ops` and `session` additionally import `internal/kernel` (minting per §5.4; replay decoding per §10.5). This is not a cycle.

**Patch 7 — fabrication-negative test form (REQ-003-03):**

> "Attempted trusted-object fabrication" is a compile-time property for external callers (unexported fields). The executable test form is: (a) zero-`core.Object` rejected with `InvalidObjectError` at every public operation; (b) an API-surface assertion that no exported constructor in `core`/domain packages accepts `Kind`/`Dimension`/`Provenance` parameters. Both run from external test packages.

**Patch 8 — placement:**

> `Challenge`/`Review` structs (§28) live in `core/corpus.go` — they are artifact data types, and no `review` package may be created (§3).

---

## 4. Pre-handoff checklist (the excerpt is truncated at §15.13)

I can fully verify everything up to §15.13 plus the entire prompt. §16 onward (session, ledger, hypothesis, manifests, canonical derivations, negative tests, acceptance list) is cut off in what you sent me. Before handing off, diff v2 → v2.1 for exactly five things — all existed in v2 and the prompt cross-references them, so this is a two-minute check, not a re-review:

1. **§16 session:** nine actions present; the §15.13 dispatch table covers every `OperationID` that can appear in a committed step (including the identification step kind).
2. **Step struct:** carries `InputCanonicals []` and `OutputCanonical` alongside hashes — hashes verify, canonicals replay.
3. **§23 derivations:** unchanged from v2 (F=ma, KE differentiation, the Substitute→Simplify→Solve→SelectBranch chain, Limit) with the "hardcoding E = mc² is forbidden" clause intact.
4. **§33 negatives + §35 tests A–S + §34/§38.1 coverage-matrix requirement:** present and unrenumbered relative to the prompt's references.
5. **Manifest item schema:** includes the `constructor` identifier field and `canonical_expr` (prompt §7 depends on both).

---

## 5. Why this is the handoff point

Three structural facts say stop reviewing and start building:

1. **The finding classes have changed.** Every round of this corpus, findings were *missing mechanisms* (provenance, conventions, canonical identity) or *contradictions* (§13 vs §26, §16 vs §31). This round's findings are *missing ring identities and test-pair enumerations* — the species of defect a coding agent fixes in minutes, not one that forces replanning.
2. **The last architectural question is answered.** The kernel-mint boundary was the final genuinely open design point in the entire corpus, and v2.1 answers it with the correct Go idiom, an explicit dependency contract, and an honest threat model.
3. **The planner now has zero authorized freedom.** Every anti-overengineering guardrail (solver pattern, calculus bounds, fixed function body, thin wrappers, dispatch replay) is normative, the Open Spec Items protocol is active, and with Patches 1–8 applied, "Open Spec Items: NONE" is a *truthful* expectation rather than pressure to hide uncertainty.

Apply the eight patches, run the five-point diff checklist, and hand both documents to the coding agent with the prompt's own final sentence as the acceptance frame: the plan succeeds when the loop runs deterministically and the system stops where judgment begins.
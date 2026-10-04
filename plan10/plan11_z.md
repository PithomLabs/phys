# Adversarial Review: The Six-Concern Assessment

## Verdict

**The triage is mostly sound and the verdict ("robust for MVP purpose, not a universal type system, no kernel growth evidenced") is one I endorse** — it's consistent with everything frozen, and the three-level promotion ladder is a genuine improvement over my own prior gate design (I split ops/kernel/spec but missed the shared-library tier). However: **one claim's dismissal is overbroad and re-opens a live spec ambiguity, one rebuttal is aspirational rather than mechanical because a load-bearing API detail is unpinned, and the assessment twice cites `plan9.2_claude` in normative-supporting roles** — the document plan10_v2_3 explicitly demotes to "historical/background only — never normative." Note also I cannot see the original six-claim review, so I audit the assessment's characterizations against the frozen documents.

---

## Six-claim scorecard (verified against frozen docs)

| # | Assessment verdict | My verification |
|---|---|---|
| 1 Type decay | Real limitation | ✅ Mechanically confirmed: `m·c²` and `r·F` both decay to `Expression`/`M·L²T⁻²`; §6.2 permits Expression+Expression Add. But see A3 — and note **no golden trace in any visible document requires object-level Add of two decayed Expressions**, which strengthens "accepted design trade, not defect": tightening it would break nothing in the corpus, yet it's spec-pinned, so it stays. |
| 2 Substitution catch-22 | Real; don't loosen Substitute; ascription lives above | ✅ Architecture-correct — **but the rebuttal is currently aspirational, not mechanical. See A1.** |
| 3 Assumption coexistence | Real; not a logic engine | ✅ Verified: `v << c` (TextValue regime) and `v = 0.99c` (expressible as `Relation(lt, v, c)` or eq-form) carry different keys → no conflict; entailment is closed-list, no contradiction detection. Consistent with §11.4 "no subsumption" and §9.8 "not a theorem prover." |
| 4 Symbol collision | Real; userland fix | ✅ Verified: §8.3 string identity; Differentiate/Substitute are syntactic. Note the frameworks *deliberately* share symbol strings (`m`, `p`, `E` appear in both packages' canonical symbol sets), so the hazard is specifically **cross-framework composition** — already procedurally guarded by AGENTS.md §9/§12. |
| 5 Transcendentals/symbolic exponents | "Does not apply to MVP" | ⚠️ **Half right.** Transcendental functions: structurally inexpressible (single `lorentz_factor` ID, closed-world Call validation — audit FIND-002). Symbolic exponents: inexpressible (`NewPow` takes `*big.Rat`). But the dismissal blanket-covers a live ambiguity it shouldn't. See A2. |
| 6 Closed ontology | Keep closed; no OO hierarchy | ✅ Correct, and the observation that `Kind` growth requires spec evolution (ordinals pinned for `mrc-v0.4`) is right. |

---

## 🔴 A1 — `Session.Identify`'s result Kind is unpinned, and that makes the claim-2 rebuttal mechanically unverified

The assessment asserts Identify "does not turn arbitrary algebra into a trusted named type." Check every visible source: spec §14.6 (operand requirements only), plan §9 (status, justification, step record — no result Kind), AGENTS.md §5, slice K. **None pins the Kind of the returned object.**

Why this matters: Compare permits *named + Expression* operands with equal dimensions. So `Identify(Energy-named, E/c²-Expression, justification)` is a legal call — and whatever Kind rule the implementation improvised (operand A's kind? some merge?) is an **unpinned open choice**. If it isn't pinned to `Expression`, then Identify is a **covert semantic-ascription path through the session authority** — the most-trusted tier in the system — gated by nothing but a non-empty justification string. That is precisely the capability the assessment says is missing and should live above the kernel; it may already exist accidentally, inside the kernel's trust boundary.

Under the §41 handoff invariant ("no remaining architectural choices"), an unpinned result Kind is itself freeze-blocking. **Pin it now**: recommended `Kind = Expression` always (Identify records a justified relation between objects; it never re-types) — then the claim-2 rebuttal becomes true mechanically, and the ascription capability stays where the assessment wants it.

---

## 🟠 A2 — Claim 5's dismissal is overbroad: the §15.8 negative-integer ambiguity is live

"No current silent failure of the sort described" is not established. The bounded power rule covers **non-negative integer** `n`; the unsupported list says "non-integer **symbolic** exponent" — which literally excludes negative *integers* from the ban (they're integers) while the rule excludes them from coverage. So `Differentiate` of `Pow(r, -1)` — the first step of any rational-function or Schwarzschild-adjacent derivation — falls in a spec gap: fail-closed error, missing rule, or worse, depending on implementer reading. I flagged this in the GR review as requiring a pin *before* Pass 1; this assessment's claim-5 verdict ("the hazard class is empty, record as future work") would wave it through. **Correct record:** transcendentals and symbolic exponents → future extension work (agreed); negative-integer power differentiation → **current spec ambiguity, pin before GR Pass 1 or the experiment's failure log opens polluted.**

---

## 🟠 A3 — AGENTS.md overstates the mechanical protection that claim 1 shows decays

AGENTS.md §8 checklist: *"Physical kinds compatible (§6.2 table; energy ≢ torque-style mismatches)."* Claim 1 demonstrates exactly the case where this promise fails: **both operands decayed to Expression with equal dimensions passes** — the energy-vs-torque pair the spec's own §6.1 exists to catch is catchable only when at least one side retains a named kind. The checklist item implies mechanical detection that degrades to nothing on the most common derivation outputs. One-line fix: *"(named-kind protection applies when at least one operand is a named kind; equal-dimension Expression+Expression mixing is permitted by design — semantic category of decayed results is the agent's responsibility)."* The assessment's framing ("solve at a higher layer") is right for the fix, but the current doc *overclaims* the present protection.

---

## 🟡 A4 — The three-level ladder is the right addition; it needs four governance pins

The Level 1 (theory-local) → Level 2 (shared library) → Level 3 (kernel) ladder fills a real gap — my prior three-track gate had no home for "used by two theories, not yet substrate-worthy." To survive contact:

1. **Promotion evidence format**: reuse the Pass-4 BEFORE/AFTER counterexample for 2→3; for 1→2 require a *named second consumer with worked example* — the assessment itself replaced "plausibly require" with named consumers at kernel level; hold Level 2 to the same standard.
2. **Home rule**: a Level-2 shared library cannot fit the frozen 39-file tree. Where do shared libraries live — post-MVP same-module packages (spec amendment), separate modules? Pin before the first library exists.
3. **Non-leakage extends upward**: my H3/H4 gates (assumption-subset invariant, forbidden-key checks) currently bind MVP packages. A shared tensor package consumed by GR and QFT is exactly where cross-framework assumption bleed re-enters. The H4 pattern must be a condition of Level-2 promotion, not just kernel admission.
4. **No bypass**: Level 2 must not become a side door to Level 3 — a shared library cannot *require* kernel semantics it doesn't have; if it needs them, that's a Gate submission, not a library API.

---

## 🟡 A5 — Citation hygiene regression

`plan9.2_claude` is cited twice in support of normative claims (assumption-boundedness, symbol package-locality). plan10_v2_3's header: plan9 is "historical/background only — never normative." This repeats the earlier `physics_compiler_mvp_specs_v2_1` stale citation — in a project whose core discipline is provenance, supporting arguments with demoted documents is self-undermining. Both claims happen to verify against v2.3 directly (§11.4; §8.3 + the per-package symbol pins), so the fix is citation swap, not substance. The trailing "AGENTS" reference is also truncated/dangling.

---

## Verified-correct (no action)

The claim-1 example computes correctly under the frozen ops (`Multiply`/`Pow` metadata laws, §6.2 Expression+Expression). The integrity-core inventory is accurate. The Go-ecosystem analogy is apt and, importantly, the discipline it protects — *"algebra does not magically manufacture physical meaning"* — is mechanically real today (Simplify's prohibitions, §15.6) **except** at the one unpinned point (A1). The final four-line classification is a fair summary.

---

## Bottom line

Endorse the assessment's conclusions and adopt the ladder — it improves the Growth Gate I previously reviewed. Before this becomes guidance for agents: pin `Session.Identify`'s result Kind (🔴 — the claim-2 rebuttal is only true once this exists), narrow claim 5's dismissal to carve out the §15.8 negative-integer ambiguity (it's current, not future), add the Expression-decay caveat to AGENTS.md §8, attach the four ladder governance pins, and purge the `plan9.2` citations. With those, "no evidence for kernel growth established" holds — and A1's fix is what makes it hold *mechanically* rather than by intention.
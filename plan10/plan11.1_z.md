# Adversarial Review: `plan11.1.md` (GR Implementation-Plan Phase)

## Verdict

**Conditional GO.** The PASS 0/PASS 1 split is sound governance, and the plan correctly absorbs A1–A5 — including one resolution that is *better* than my own recommendation (conceded below). But there is **one internal contradiction between Pin B and the workload anchor that will either stall the experiment or tempt the implementing agent into exactly the freeze violation Pin B forbids**, plus a dimensional error baked into Pin A's own test fixture — ironic in a project about dimension checking — and three dropped governance items from the A4 pin set.

---

## 1. Prior-findings disposition (verified)

| Prior | Disposition in plan11.1 | My verification |
|---|---|---|
| A1 (Identify Kind unpinned) | **Pin A: `KindRelation`, test-only** | 🟢 **Conceded — and improved.** plan11_review's `KindRelation` is strictly safer than my `Expression` suggestion: a Relation-kind object is maximally inert (rejected by Add's table, unusable as named quantity, Solve-inert unless it matches the single quadratic pattern). The no-covert-ascription property holds *mechanically*, not by convention. My recommendation is superseded. |
| A2 (§15.8 negative-integer ambiguity) | **Pin B: fail-closed, test-only** | 🟢 Resolved in the correct direction; see 🔴 P1 for the consequence nobody wrote down. |
| A3 (AGENTS.md §8 overclaim) | **Pin C: one-line caveat** | 🟢 Matches. |
| A4 (ladder governance) | D3 (separate modules) + §6 named-second-consumer | ⚠️ **Partial** — see 🟠 P4. |
| A5 (citation hygiene) | §8 reconciliation | ⚠️ Acceptable for a reconciliation record, but the GR plan itself must cite only frozen sources. |
| F2 (failure taxonomy pre-registration) | §5, exactly five classes | 🟢 Adopted verbatim, including silent-wrongness as primary thread. |
| F4 (Outcome-A likelihood) | §6 "Default outcome: NO KERNEL GROWTH"; C6 closed | 🟢 Expectation set correctly. |
| F6 (post-freeze, out-of-tree, independent veto) | §3 §17, D3 | 🟢 Mostly — see 🟠 P3 for the Level-1 gap. |

The static-spherical anchor is well chosen: it avoids transcendentals entirely (Schwarzschild is rational), isolates index structure inside the thread (D2 — no premature generic system), and the diagonal metric keeps inversion within `Pow(g_ii,−1)` reach. GR-0…GR-8's "where the substrate bends, breaks, or silently misrepresents" is the right logging posture.

---

## 🔴 P1 — Pin B and the workload anchor are mutually unsatisfiable as written; the resolution must be declared now

Pin B locks: `Differentiate(Pow(x, negative-integer)) → UnsupportedOperationError`. The anchor demands: `coordinates → metric → inverse → Christoffel → …` over **static spherical spacetime, whose every metric component is a rational function of `r`** — `g_tt = −(1−2M/r)`, `g_rr = (1−2M/r)⁻¹` — and the pinned division representation is `a/b ≡ Mul(a, Pow(b,−1))`. Therefore **every differentiation in GR-3 through GR-8 traverses negative-integer powers and hits the just-pinned rejection.** Kernel `Differentiate` is unusable for the entire primary thread past GR-2.

The plan's own §4 says the thread "forces the hard encounters honestly (`r⁻¹`/`r⁻²` negative powers, `∂g/∂x` differentiation demands)" — so the authors know the encounter happens. But the resolution is unstated, and there is exactly one: **the GR/physmath module implements its own bounded rational-function differentiator** (products/powers/sums of symbols and rationals, minted through the `core` façade constructors), leaving kernel `Differentiate` frozen. Note this does not violate REQ-005-03 — that rule binds MVP *domain wrappers*, not a separate post-freeze module — but it must be said, because the failure mode of leaving it implicit is precisely the one Pin B's "Do NOT expand differentiation to make GR easier" warns against: an agent at GR-3, facing a pinned rejection *and* an anchor that demands Christoffel symbols, "helpfully" extends `Differentiate` — violating both the pin and the frozen tree.

**Fix (one paragraph in plan11.1 §4 or PASS 1 §7):** declare package-level differentiation as the pre-registered workaround; classify the encounter `SPEC-INTENDED-BOUND` + `PACKAGE-SOLVABLE` up front; state explicitly that this classification **cannot constitute Level-3 evidence** (the invariant is preserved by the workaround — the gate question is whether the *package* differentiator can preserve the kernel's guarantees, and it can, since it mints through the same façade).

---

## 🟠 P2 — Pin A's negative-test fixture contains a dimensional error and tests nothing

"`Identify(Energy, E/c², …)` must not yield `RestMass`" — compute: `E/c²` has dimension `ML²T⁻² / L²T⁻² = M` (**a mass** — `E = mc² ⟺ m = E/c²`), while the Energy operand has `ML²T⁻²`. The call fails MRC-002 and returns an error: the test passes **vacuously**, exercising dimension rejection, not kind behavior. As a bonus problem, if a fixture-builder "fixes" it by declaring the `E/c²` object with Energy dimension, the fixture encodes a physics error.

The sharp version uses a **dimension-compatible** pair so the kind path is actually exercised: `Identify(RestMass, E/c²)` (both `M`; named+Expression legal per Compare-kind rules) → assert result `Kind == KindRelation`, **not** `RestMass` and not the named operand's kind. (`Identify(Energy, mc²)` works equally.) Also add to the positive pin: result `CorpusStatus == NONE` and result dimension == shared operand dimension — both are plan-pinned behaviors currently unasserted.

---

## 🟠 P3 — The GR package's home is unpinned (D3 covers Level-2 only)

D3 pins Level-2 libraries as separate modules. The **Level-1 `general_relativity/` package's location is stated nowhere** — and REQ-002-20 plus the frozen 42-file tree forbid it from living in `phys`. One line fixes it: *GR is a separate module/workspace depending one-way on frozen `phys`, same rule as Level-2; the phys tree is untouched by all GR work, not only promotions.* PASS 1 §6 then elaborates.

---

## 🟠 P4 — Two A4 governance pins were dropped from the ladder

1. **Non-leakage extension:** the H3/H4 pattern (assumption-subset invariant, forbidden-key checks) must be a **condition of L1→L2 promotion**, not merely kernel admission. A shared tensor library consumed by GR and a future QFT module is exactly where cross-framework assumption bleed re-enters. §6's L1→L2 criteria (named second theory + worked example + shared semantics) omit it.
2. **No-bypass rule:** a Level-2 library must not require kernel semantics it doesn't have; if it does, that's a Gate submission, not a library API. §6's L2→L3 burden implies it; state it.

---

## 🟡 P5–P9 — Smaller pins

- **P5:** §6 routes **C7/C8**, but this document defines only six concerns. The routing table isn't self-contained; PASS 1 §11 must restate all eight with one-line definitions and source attribution (plan11_review), or a reader cannot execute §6.
- **P6:** **Schedule the silent-wrongness probe explicitly.** The primary thread won't surface the benchmark hazard naturally: with explicit coordinate components, `M`, `G`, `c` are *genuinely constant* symbols, so `Differentiate`'s silent-zero-on-non-coordinate-symbol behavior is *correct* throughout GR-3..GR-8. The silently-wrong class needs a deliberate probe (differentiate an expression where a symbol denotes a field of the variable; observe silent `0` with `DERIVED` provenance; log with the REQ-015-01 tension — the externality blocker is that a field registry is ambient state). Add it as a scheduled test, not an hoped-for byproduct.
- **P7:** The "20-field Growth Evidence Record" is cited before definition — it must be enumerated in PASS 1 §13, not inherited as a count.
- **P8:** "42-file re-freeze AUTHORIZED (adv_review11)" — instrument unverifiable in this packet, and `specs_v2_3.md` byte-frozen + 42 files = a standing §3 deviation that must be recorded under REQ-003-01's specification-defect clause. Cite where the authorization and exception are documented. (Corroborating detail: 42 = 39 + AGENTS.md + two package READMEs — consistent with my B1 resolution path.)
- **P9:** Pin B resolves the §15.8 ambiguity *de facto* while the spec text remains silent. Record in PASS 1 §2 (Frozen Baseline invariants): negative-integer-power differentiation is fail-closed **by pinned test**, spec text unpinned — so a future spec reader doesn't "discover" the gap and expand it.

---

## Bottom line

The plan's structure — pins before planning, hard-gate halt semantics, taxonomy-logged workload, default NO-GROWTH, separate modules — is the correct operationalization of everything this review chain has established, and Pin A's `KindRelation` resolution is genuinely superior to what I proposed. Before PASS 0/1 execute: **write down the package-level differentiator (P1 — the contradiction is real and the failure mode is a freeze violation), fix Pin A's fixture to a dimension-compatible pair (P2), pin the GR module home (P3), restore the two dropped A4 pins (P4), and fold in P5–P9.** With those, the plan is ready — and P1's declaration is what keeps "do NOT expand differentiation to make GR easier" from colliding with "reach Christoffel symbols" in the same agent's context window.
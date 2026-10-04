# Adversarial Review: `GR_IMPLEMENTATION_PLAN_20261002.md` (v2) — Line-by-Line

## Verdict

**Architecture: holds. Blockers to the frozen `phys` kernel: none. The stated target — "no ambiguity for coding agent about architectural decisions" — is not yet met.** The v2 revision correctly lands every structural lock from the consolidation, but **three of my four 🔴 pins were acknowledged in prose yet never actually pinned** (derivative representation, reduction protocol, tensor arithmetic ops), and the bridge itself introduced a new determinism-critical unpinned choice (bridge object `id`). The changelog's claim "implicit implementation choices pinned" is therefore overstated. What remains is a finite pin list — nothing requires redesign.

---

## 1. Prior-pin resolution audit

| Pin | v2 disposition | Verdict |
|---|---|---|
| PIN-1 derivative representation | §6.4 acknowledges `A′/B′ (+ higher derivatives as curvature requires)` but the §6.2 node list (`…Sin, Cos, UnknownFunction`) still contains **no derivative node, no derivative-tagged form, no max-order pin** | 🔴 **Unresolved** — see §2 |
| PIN-2 bridge granularity | §6.2 says conversion requires the "entire expression" be kernel-representable; §7's D12 example bridges a "kernel-compatible **term**" | 🟠 **Wording tension** — see §3 |
| PIN-3 reduction protocol | §5: "bounded differential/algebraic reduction (problem-specific bounded local solver; NO general ODE solver)" — a label, not a method | 🔴 **Unresolved** — see §2 |
| PIN-4 tensor ops | §8 list is still exactly `Contract / RaiseIndex / LowerIndex / ApplyMetric / ApplyInverseMetric` — no `ScalarMul`, no `TensorAdd/Sub`; symmetry-inheritance rule absent | 🔴 **Unresolved** — see §2 |
| PIN-5 symbol strings | Coordinate *order* pinned (`x⁰, r, θ, φ`); exact Symbol strings, time-coordinate symbol (`t` vs `x⁰`), mass symbol (`M` appears in §5/§7 as display notation) all unpinned | 🟠 **Partial** |
| PIN-6 GRStep chaining | §7 GRStep has `CurrentHash` but **no `PreviousHash`** — unchained, and the omission is undeclared (v1 said "replay/hash trace," consolidation recommended chained envelope) | 🟡 Open — choose chain or record deliberate unchained rationale |
| PIN-7 ZeroTest boundary | §9.3 "zero detection" — no rule for UnknownFunction-bearing expressions (undecidable class) | 🟡 Open |
| PIN-8 node-set closure + Sqrt | Node list given; closure/amendment rule unstated; Sqrt unaddressed | 🟡 Open — see √\|g\| trap, §4 |
| PIN-9 bridge conventions | §8: "attach compatible `ConventionSet` metadata for frozen conflict detection" | 🟢 Resolved |
| PIN-10 bridge justification | §6 justification present; ESCALATE-TO-SPEC pre-designation not added (generic §18 covers it) | 🟢 Adequate |

**Bookkeeping:** supersession handled (changelog + v1 preserved; v2 declared execution contract — the v1 stale probe text can't mislead an agent executing v2) 🟢. D-namespace collision not fixed (D1–D18 unnamespaced vs plan11.1's D1–D3), but v2 is self-contained and plan11.1 is preservation-only — low residual risk 🟡. **C-1 still unfixed**: `plan10_v2_3.md` §9's `OperationParams` example still shows `operator:"eq"` on a `pow` artifact that plan7's strict rejection requires the implementation to refuse — still one line, still on the pre-PASS-0 list. ⚪ still unverifiable in packet: `session/session.go:278` grounding; `adv_review11.md` (now load-bearing for §3.3 — attach it).

---

## 2. The three acknowledged-but-unpinned decisions (guaranteed encounters)

**🔴 U1 — Derivative-of-UnknownFunction representation (bites at GR-3b, certainly at GR-5).** Riemann of the ansatz contains **A″(r) and B″(r)** — this is arithmetic fact, not possibility. The plan requires `dA/dr = A'(r)` "+ higher derivatives as curvature requires" but never says what A′ *is* as a node: `Deriv{fn, var, order}`? A distinct `UnknownFunction("A′", r)`? An order field on `UnknownFunction`? Related rule needed: are `A` and `A′` independent atoms for Normalize/like-term collection (they must be, or nothing cancels), and is `A″` reachable (max order 2 for this workload)? **Pin:** one representation + the independence rule + `maxDerivativeOrder = 2`.

**🔴 U2 — Vacuum-reduction protocol (bites at GR-7→GR-8, the climax of D14).** The plan bans the general escape (no ODE solver, §5/§19) but supplies no specific method — the agent is squeezed between prohibition and invention. The honest derivation has a standard first-integral structure: combine `R_tt/A + R_rr/B = 0` → `A′/A + B′/B = 0` → `A·B = const` → boundary conditions (`A,B → 1` at flat infinity) → `B = 1/A` → integrate `R_rr = 0` (or power-law ansatz + coefficient comparison via Normalize/ZeroTest) → `A = 1 − 2M/r` with `M` as the pinned integration constant. **Pin this exact protocol** (combination trick, boundary conditions, ansatz form, coefficient-comparison mechanic). Otherwise the agent improvises the one step D14 exists to exercise.

**🔴 U3 — Tensor arithmetic (bites at GR-6/GR-7).** `G_μν = R_μν − ½R·g_μν` needs: Contract ✓, Contract ✓, then **scalar-times-tensor and tensor-minus-tensor** — neither in the §8 ops list. Add `ScalarMul` and `Add/Sub`, plus one symmetry-inheritance sentence (G_μν symmetric because both summands are; how `Symmetries` is populated by construction).

**🔴 U4 (new this round) — Bridge object contract.** `NewCandidateConcept(id, …)` **requires a non-empty `id`** (plan §4: hypothesis `id` requires names), and the id enters canonical JSON → `HashObject` → GRStep determinism. Unpinned, an agent will improvise `"bridge1", "bridge2", …` — nondeterministic across runs, silently breaking D11. Pin the full bridge contract: **id = deterministic function of content** (e.g., `"gr/" + hex(HashExpr(bridgedExpr))`), plus the standard attached assumption set (which assumptions ride every bridge? recommend a minimal pinned GR set: `Constraint/r_positive = (r > 0)` + regime entry — also what makes kernel `Simplify` legally do `Sqrt(Pow(r,2))→r`-class work during D12 reuse) alongside the already-pinned ConventionSet.

---

## 3. 🟠 Two-wording ambiguities

- **Bridge granularity (PIN-2).** Reconcile §6.2 ("entire expression") with §7 ("kernel-compatible term → bridge"): one sentence — *"`ToCoreExpr` is total over its argument: it succeeds iff that argument's tree contains only kernel-mirror nodes; callers may invoke it on any subexpression, e.g. individual terms of a sum; mixed sums are decomposed before bridging."* Without it, D12's reuse rule is unenforceable-by-reading.
- **Symbol strings (PIN-5).** GRStep and bridge hashes depend on exact strings. Pin: `t, r, theta, phi, M` (ASCII per the accepted `eta` precedent — note §8 currently writes `θ, φ`); resolve time-coordinate as `t` (matching the §5 ansatz `dt²` and the SR package) with `x⁰` as display-only; pin `M` as the mass symbol in both geometric and SI subtracks.

---

## 4. 🟡 Remaining pins (one editing session)

1. **GRStep chaining** — add `PreviousHash` (envelope chain, consistent with the project's own tamper-evidence design) or record unchained as deliberate with rationale.
2. **ZeroTest boundary** — decidable for rational+`Sin`/`Cos` concrete expressions; UnknownFunction-bearing → "cannot decide; keep symbolic" except structural identity; pin it (vacuum equations are exactly the undecidable-looking class).
3. **Node-set closure** — declare the §6.2 set closed; amendments require plan revision; decide **Sqrt: excluded** — and add the **√|g| trap line**: determinant-form Christoffel identities (common in textbooks) are forbidden; the standard component formula is the pinned route (otherwise the agent invents a local Sqrt node mid-GR-3).
4. **§7 vs §12 conflict** — §7 says kernel `Differentiate` is "GR-3a probe only," but §12 Probe A legitimately calls it in the independent suite. Add: "§7 governs workload usage; §12's adversarial suite may invoke any frozen op."
5. **Geometric↔SI mass mapping** — the tensor track yields `M` as a *length* (geometrized); the correspondence subtrack needs SI mass. Pin the conversion (`M_geo = G·M_SI/c²`) at the GR-8 boundary or the dimension audit (`GM/r² → acceleration`) can be assembled inconsistently.
6. **§9.1 wording** — "negative-integer/rational terms": clarify this means *rational-function terms* (not fractional-exponent support), keeping the differentiator out of CAS territory.
7. Optional (recommended, not required): use `Compare` + `Session.Identify` to *record* the GR-8 Newtonian correspondence (HYPOTHESIS operands → HYPOTHESIS output — correct epistemics, and it exercises session authority, which the workload otherwise never touches); and schedule the free boundary probe `Limit(g_rr, r, 0) → UnsupportedOperationError` (singular-substitution refusal) alongside GR-3a.
8. Cosmetic: duplicated "Pin A:" in §16; `go.work` shown inside `phys-gr/` (workspaces conventionally live at the parent connecting both modules — pin its home); "GR-3a must fail loudly" → "the probe test fails (a derivative was returned)."

---

## 5. What v2 gets right (verified, no action)

The bridge resolution is correctly grounded (§5.4 path 2, REQ-024-02, §11.3) and genuinely strengthens the test — full MRC/assumption/contamination surface now exercised, with GR outputs epistemically honest as `HYPOTHESIS`. §3.4's "compiler-checked guarantee" is technically accurate (cross-module `internal/` enforcement). The `replace`-into-internals retirement is correct. Silent-wrongness rescoping (§10/§12) is coherent: the frozen kernel's deliberate acceptances are pre-classified ledger items (C1–C8); `SILENTLY-WRONG` now correctly targets Level-1 machinery bugs; Probe A's ACCEPTED-CORRECT reclassification retires the former benchmark counterexample properly. The §14 wording lock, D17 rename, independent-classification rule, evidence fields 18–20, mutation protocol, PASS0 file-membership check (verified: all three targets ∈ 42-set), and the three-false-readings paragraph all land as prescribed. Appendix B's rejection list remains justified item-by-item — including my withdrawn placeholder proposal.

---

## Bottom line

**No architectural blockers — the frozen kernel, the module boundary, the bridge, and the gate are all correctly locked. But the plan's own success condition ("no remaining architectural choice for the implementation agent") fails on four items the locked pass structure guarantees the agent will face**: U1 (A″ representation at GR-5), U2 (the reduction method at GR-7/8), U3 (tensor add/scale at GR-7), U4 (bridge id/assumptions at every contact). All four are one-paragraph pins with concrete recommended text above; add the two 🟠 wording fixes and sweep the 🟡 list in the same pass, fix C-1 in `plan10_v2_3.md` §9 while you're editing documents, and attach `adv_review11.md` to the packet. With that single revision, the target condition is genuinely met — and the changelog's "implicit implementation choices pinned" claim becomes true rather than aspirational.
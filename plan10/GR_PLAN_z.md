# Adversarial Review: `GR_IMPLEMENTATION_PLAN.md` (Plan 11.1 Revised) + generating prompt

## Verdict

**Conditional GO.** All nine prior findings are resolved — two of them exactly as prescribed, one (the Pin A fixture) even mutation-validated. The prompt↔plan conformance is exact on the required 20-section structure. But the revision introduced **one genuine architectural contradiction the plan must resolve before PASS 1 executes** (the plan lists Object-level kernel ops as the GR interaction surface while simultaneously forbidding GR from ever owning a `core.Object`), plus an unpinned trig-composition representation and a wrong operation mapping for the Newtonian limit. None of these threatens the frozen `phys` tree — all are phys-gr-side decisions the plan was supposed to pin and didn't.

---

## 1. Prior-findings disposition (verified)

| Prior | Disposition | Verification |
|---|---|---|
| P1 (Pin B × anchor contradiction) | GR-3a/GR-3b ordered sequence, §7 + §9.1 | 🟢 Resolved — probe first, local differentiator only after, "must not be represented as kernel capability," and PACKAGE-SOLVABLE is excluded from serious evidence by §10. |
| P2 (Pin A vacuous fixture) | `Identify(RestMass, E/c²)`, dim M | 🟢 Resolved — dimension-compatible, named+Expression legal per Compare-kind rules, so the kind path is genuinely exercised. **§17 even mutation-tests the fixture back to the broken form** — that's the right discipline. |
| P3 (GR module home) | §4 sibling module, one-way dep | 🟢 Resolved (see N4 for build mechanics). |
| P4 (dropped A4 pins) | §13 non-leakage + no-bypass rule | 🟢 Both restored verbatim. |
| P5 (C7/C8 undefined) | §11 self-contained ledger | 🟢 Faithful — classifications and routing match the settled record. |
| P6 (silent-wrongness scheduling) | §12 six probes, field-differentiation explicit | 🟢 Present with expected outcome pre-registered (see N5 note on its containment path). |
| P7 (20-field record) | §15 enumerates all 20 | 🟢 Includes artifact-compat, non-goal collision, silent-wrongness analysis, named second consumer, independent review, human approval. |
| P8 (42-file citation) | §3 cites `adv_review11.md` | 🟢 Cited; instrument remains ⚪ unverifiable in this packet. |
| P9 (Pin B de facto vs spec text) | Pin B distinguishes frozen text from pin | 🟢 Resolved. |

The §14 wording change ("Default outcome: NO KERNEL GROWTH" → "no kernel modification presumed; all outcomes admissible, determined from recorded evidence") is an **improvement**: the old default could suppress genuine evidence ("userland coped ⇒ done"), while the anti-growth asymmetry now lives where it belongs — in the evidential burden (fields 4–7, 14–18), the named-second-consumer requirement, and "burden of proof is on growth," which is retained.

---

## 🔴 N1 — The plan is self-contradictory about what GR components *are*, and the interaction surface is unexecutable as written

- §6.2/§6.4: GR scalar components are **`core.Expr`** inside GR-local structs; GR objects are "**not `core.Object` instances** — the frozen module has no public generic trusted-object factory."
- §7's kernel interaction table: `Add/Subtract/…/Simplify/Differentiate` used for "rational component algebra … **MRC enforced**."
- But **every kernel op signature takes `core.Object`** (§15.1–15.12) — MRC-002/003 operate on Objects' Dimensions and Kinds, which bare `Expr`s don't have.
- And phys-gr is a **separate module**: `internal/kernel.MintObject` is unimportable (REQ-004-03). The public Object-producing surface is: fixed domain constructors (wrong physics), ops (Object→Object), session actions (derive from existing Objects), and **`hypothesis.NewCandidateConcept(id, kind, dimension, expr, assumptions, conventions)`** — the *only* public path minting an Object from raw (Kind, Dimension, Expr, assumptions).

So either (a) GR stays Expr-only — then **no kernel op can be called on any GR component, the §7 table is decorative, and even the GR-3a probe `phys.Differentiate(Pow(r,−1))` cannot be built** (its operands must be Objects) — or (b) GR wraps components via the hypothesis constructor, which §6.4 currently denies.

**Recommendation: adopt (b) and pin it.** Wrap components as `Kind=Expression` Objects with exact Dimensions (public `Dimension` Multiply/Divide/Pow make curvature `L⁻²` etc. constructible), GR regime assumptions as real `AssumptionSet`s, minted via `NewCandidateConcept`. This is strictly better for the experiment: the full MRC surface becomes exercisable (the actual thing under test), kernel assumption machinery becomes usable by GR — e.g., attach `Constraint/r_nonnegative = (r ≥ 0)` and kernel `Simplify` will legally perform `Sqrt(Pow(r,2)) → r` via bounded entailment — and every GR-derived quantity carries **`HYPOTHESIS` provenance by the contamination law, which is epistemically exactly right** for userland physics. Log the friction itself as boundary evidence: *the only public arbitrary mint in the frozen substrate is the hypothesis path* — that is by design (trusted objects are curated), and the discomfort external theory packages feel at it is a genuine Growth-Gate data point, not a defect to route around.

---

## 🟠 N2 — The trig/rational composition boundary is unpinned (`r²sin²θ`)

§9.2 correctly assigns trig to phys-gr and rational algebra to `core.Expr` — but the Schwarzschild angular components are **both at once**. `g_φφ = r²sin²θ` cannot be a pure rational `Expr`, and a GR-local struct can't be handed to kernel `Simplify`. The plan must pick a representation before GR-1 golden states exist. Recommended: **placeholder symbol** (`Symbol("sin_theta")` in the `Expr`) + GR-local semantic/differentiation table keyed on the placeholder, with the unfaithfulness **logged as `REPRESENTABLE-BUT-UNFAITHFUL`** (the symbol is not the function). This is not a workaround to hide — it is precisely C4/C5 evidence, and the local explicit-argument table is clean of ambient state.

## 🟠 N3 — §7's `Limit` row mis-maps the Newtonian limit

"Limit | Newtonian-limit step" is wrong: kernel `Limit` is substitution+simplify plus the fixed `lorentz_factor` expansion — it cannot express `r→∞` (no infinity representation) and series expansion doesn't exist. GR-8's weak-field correspondence is **`Compare` to build the relation + `Session.Identify` with justification for the φ = −GM/r correspondence** + GR-local asymptotic bookkeeping. Pin this now, or the implementing agent burns a pass discovering `Limit` can't do it. Note the elegant consequence: since the GR-side operand is hypothesis-minted, the correspondence lands as `HYPOTHESIS` by contamination — the correct epistemic status for a correspondence identification. Also pre-register the assumption-form trap: entailment reads only `x≥0`/`x>0`-shaped assumptions, so attach the simple form (`r ≥ 0`) alongside the physical regime (`r > 2M`) or `Sqrt(Pow(r,2))` won't simplify.

## 🟠 N4 — Module/build mechanics and frozen-`phys` integrity are unstated

A sibling module requiring `github.com/PithomLabs/phys` needs `go.work` or a `replace` to resolve locally — and "No `replace` into kernel internals" is technically meaningless as written (replace maps module roots, not internals). Pin: *Go workspace (or replace to the **pristine, hash-verified** local phys checkout only; no modified forks, no vendoring)*; GR work in a separate clone/worktree; frozen-`phys` git hash recorded in each evidence record. §17's mutations also inherit the Gate H disposable-copy/`-count=1`/restore discipline — state it.

---

## 🟡 N5 — Small pins (fold into one pass)

1. **§7 scope note:** "Differentiate — GR-3a probe only" governs the workload; §12 probe 2 legitimately uses kernel `Differentiate` in the separate adversarial suite. One clause.
2. **GR-3b should reuse kernel `Simplify`** for cleanup (its output terms — zero-product, `Add(x,0)` — are kernel-simplifiable); local rules only for the negative-power/trig gaps. One line in §9.1.
3. **Golden-state pins:** GR adopts signature `−+++` (consistent with the relativity package convention) and a units policy (geometric `G=c=1` with dimension bookkeeping via public Dimension arithmetic) before GR-1 — otherwise golden components are agent-choice.
4. Pin A has a duplicated "Pin A:" typo; §17 "fail loudly" should read "the probe test fails (derivative returned)."

---

## What holds (hand-checked)

The GR-3a classification (`SPEC-INTENDED-BOUND`) is correct per the taxonomy; the workload still avoids transcendental *kernel* encounters (sin²θ now honestly claimed, contra the earlier draft); `Differentiate` of `t^(4/3)` remains out of thread scope; the C1–C8 ledger is consistent with the settled six-concern record plus the two pins; §18's three false-readings paragraph is the best governance writing in the document chain; non-goals §19 correctly retains the E=mc² firewall and the REQ-002-21 tripwire (field 12 of the evidence record).

## ⚪ Unverifiable in this packet

`session/session.go:278` KindRelation grounding (Pin A's "zero production diff expected" claim), the existence/content of `plan10/adv_review11.md`, and the 42-file membership list.

---

## Bottom line

**Fix N1 before PASS 1 executes** — it is the one finding that makes the execution contract internally unexecutable, and its resolution (hypothesis-wrapped, dimension-carrying, assumption-carrying Objects) *strengthens* the experiment by putting the full MRC/assumption/contamination machinery in GR's reach rather than leaving the stress test at the untyped `Expr` level. Pin N2's placeholder strategy and N3's Identify-based Newtonian limit in the same pass, add N4's build/integrity mechanics, and sweep N5. With those, this is an execution-ready contract — and note what N1's resolution gives you for free: the first `phys-gr` evidence records will document that an external theory package *can* build fully MRC-guarded, assumption-carrying, contamination-honest physics on the frozen public API without touching a single frozen byte. That is the strongest possible `NO-GROWTH` evidence the Growth Gate could ask for.
# Adversarial Review: `plan10_v2_3.md` vs. `specs_v2_3.md` (MVP scope)

## Verdict

**No architectural blockers. Conditional GO** — the architecture is unchanged and sound; all five prior findings were resolved. What remains is one internal contradiction in the replay-decoding wording, one self-inconsistent pinned fixture, and a cluster of under-pinned corpus data that will cause agent-guessing (determinism risk). All fixable in a single editing pass.

---

## 1. Prior findings — resolution audit

| Prior ID | Status | Notes |
|---|---|---|
| F1 relativity item list | 🟢 Resolved | 10 items pinned; `Velocity` + 3 zeros explicitly non-manifest (§4, §7, Step 11 all agree now). |
| F2 `Pow(x,0)` | 🟢 Resolved | §8.7 added to Simplify set; `0^0` still errors. See N5 for residual ambiguity. |
| F3 kernel constructors | 🟢 Resolved | Step 1 owns all validation once; `core` façades delegate. |
| F4 set ordering | 🟢 Resolved | Canonicalize→sort-by-bytes→dedup; feeds hash chain. |
| F5 slice D fixture | 🟢 Resolved | `Mass` vs `RestMass` (both dim `M`) — correct this time. |
| 🟡-1..8 | 🟢 Resolved | Params schema, §8.5/§9.6 boundary, Relation/BranchSet simplify, SelectBranch 6-step, defensive copies, kernel helper ownership, `schema_version:"1"` provenance-of-pin, decoder invariant validation. |
| v2.3 Reconciliation A/B/C | 🟢 Absorbed | Challenge/Review sole site in `core/corpus.go`; `≠0` vs `>0` entailment split correctly distinguishes admissibility from sign; positional table matches README's §15.13.1 claims. |

I re-traced the golden derivations against the revised rules (limit body retains `Mul(1,…)` until Simplify; `Mul(0, c⁻¹)→0` via `c>0 ⊨ c≠0`; `Sqrt(Pow(m·c²,2))` unwraps via `m≥0` + even exponent) — all still hold.

---

## 2. New findings this round

### 🔴 N1 — Replay-decoder wording contradicts spec §4 and plan §9

Plan §4/Step 3: *"no exported `DecodeObjectJSON`/`ParseObject` **anywhere**"*. But:

- Spec §4 explicitly permits `session` to import `internal/kernel` "for validated minting **and replay decoding**."
- Plan §9 requires `session.Validate` to decode `InputCanonicals`/`OutputCanonical` — which only `session` does.

If the decoder is unexported inside `kernel`, session replay **cannot compile**. The evident intent (per the 🟡-8 note) is "no decoder exposed through the *public* API (`core`)" — which is correct and spec-conformant, since `internal/kernel` is already module-internal. Fix the wording: *"Object decoding is exported from `internal/kernel` only (module-internal visibility); no public package exposes any object decoder; `session` accesses it via its permitted direct kernel import (spec §4)."* As written, a literal-minded agent either breaks replay or invents a workaround.

### 🔴 N2 — The pinned `OperationParams` canonical-JSON example violates its own binding rule

Plan §9 gives: `{"kind":"pow","exponent":"2/1","operator":"eq","justification":""}` — a `pow` params carrying `operator:"eq"`. That contradicts both the binding table (`pow` ⇒ only `Exponent` used) and "unused fields encode their empty/default values." This example is exactly the kind of fixture that becomes a test vector and then a conformance bug. Correct form: `{"kind":"pow","exponent":"2/1","operator":"","justification":""}`.

### 🟠 N3 — Corpus data pins still incomplete (determinism gap)

The plan's selling point is "no open choices," but three corpus details remain agent-choice:

1. **RestFrameAssumption exact form is not in the plan.** It's pinned in `prompt2.md` (Phase 6: `Constraint`, key `rest_frame`, value `Relation(eq, Symbol("p"), Rational(0))`), but prompt2 is an *audit* prompt, not the coding agent's authority. Import the pin into plan §7 verbatim, including the warning that nothing may treat rest-frame as ambient state.
2. **Per-item assumption matrix is under-pinned.** `rest_mass_nonnegative`/`speed_of_light_positive` are pinned for the two named carriers — but what does `LorentzFactor()` carry (its body divides by `c`)? What about `Energy`, `ThreeMomentum`, `FourMomentum`, `Spacetime`, `MinkowskiMetric`, and all mechanics items (presumably none)? The manifest↔constructor cross-check forces *agreement* but not *content*; two runs can deterministically produce different corpora. Add an explicit item→assumptions table (mostly "∅").
3. **Zero constructors' Kind/Dimension unpinned.** `ZeroThreeMomentum` is pinned via prompt2 (Kind `ThreeMomentum`, dim Momentum, `DEFINED`, carries rest-frame). `ZeroEnergy`/`ZeroVelocity` get only names. Pin Kind `Energy`/`Velocity`, dims Energy / `L T⁻¹`, expr `Rational(0)`, `DEFINED`, no assumptions.

### 🟡 N4 — The `go/parser` ban is stated absolutely but the plan's own audits need it

README: "`go/parser`, `go/scanner`, `go/token` banned." But the export-surface scans, import-independence scans, and `TestNoForbiddenExportedAPI` require *parsing Go source* — which needs `go/parser` + `go/ast` in **test files**. Plan §7 scopes the ban to non-test files correctly; README doesn't. An over-literal agent will either refuse the scans or hand-roll a scanner. Scope the README statement to production/non-test code.

### 🟡 N5 — "Nonzero-safe base" for `Pow(x,0)→1` is still undefined

Plan copies spec §8.7's phrase verbatim, which defers the ambiguity. The deterministic pin consistent with "`0^0` errors": rewrite applies iff base is not structurally `Rational(0)` and not entailed-zero; otherwise `UnsupportedOperationError`. No golden trace needs it, but pin it before someone "resolves" it in code.

### 🟡 N6 — Unused-params strictness in `Session.Step` unpinned

Binding table lists "params fields used." Must `Step` *reject* a `subtract` call carrying `Justification:"x"`, or ignore it? Recommend strict reject (smaller tamper surface, deterministic ledgers). One sentence fixes it.

### 🟡 N7 — `TestApplyPositionalInputs` "reorder rejection" is only meaningful for order-sensitive IDs

`Apply("add", [b,a])` is indistinguishable from a valid call — commutative ops cannot reject swaps. Scope the reorder assertions to `subtract, divide, substitute, differentiate, limit, compare, solve, select_branch` + params-binding failures (`pow` without `Exponent`). Swapped-input *tamper* on commutative ops is still caught by replay divergence — but that's the ledger's job, not `Apply`'s.

### 🟡 N8 — `Challenge`/`Review` canonical encoders missing from Step 9 deliverables

Spec §10.1 lists Review/Challenge as canonical artifacts (REQ-010-01..03 apply). Plan Step 9 names the manifest functions but not their canonical JSON/equality — yet `ReviewHistory` in the candidate needs them. Add to `core/corpus.go` deliverables.

### 🟡 N9 (README) — `go build ./...` cannot verify session import-independence

It verifies cycle-freedom only; `session` importing `mechanics` would still build. The real verifier is the AST import scan (`TestSessionImportIndependence`). Fix the README attribution.

### 📌 N10 (accepted limitation, no action) — Simplify engine is build-gated only until Step 8

The most intricate code (entailment, finiteness, §9.6–9.9) gets no test for three steps. Within MVP constraints this is tolerable since Step 8's suite is comprehensive; optionally pull `ops/operations_test.go` creation into Step 5 with incremental cases — the file budget doesn't change.

---

## 3. MVP-scope limitations correctly accepted (verified, keep documented)

Closed §9.8 entailment list (no derived inference); single-pattern `Solve`; `lorentz_factor` as the only `Call`; no L'Hôpital/series; integrity-not-authenticity chain (README states it); corpus status human-curated, never computed; `APPROXIMATED`/`Approximation` reserved; no object parser in the public API (round-trip via in-module tests only — spec §10.5 frames this as internal, so conformant); manifest minimum-set pinning of spec's "at least." All consistent with §2 non-goals and correctly surfaced in README.

## ⚪ Unverifiable in this packet (request full docs at hand-off)

Both documents truncate again — spec at §15.9, plan mid-§9. Unchecked: §15.11–15.13 exact wording (plan claims to quote §15.13/§15.13.1), §16.12/§16.19/§16.20, §17/§19/§26–§28/§32/§41, the plan's A–S slice table and final gate, and `TestManifestRejectsInconsistent`'s survival into v2.3. These cover precisely the session/candidate/manifest regions where residual risk concentrates.

---

## Bottom line

**Conditional GO.** Apply N1–N3 (the only ones that can produce a wrong artifact rather than an imprecise one), fold N4–N9 into the same pass, and attach complete documents to the hand-off packet. Nothing requires re-architecting; the v2.3 revision demonstrably closed every prior finding, including the three spec reconciliations.


## Audit Review

# Does the Review Hold? — Re-assessment Against `adv_review5.md`

## Short answer

**Yes — the review stands, and the audit strengthens it.** The two documents review different artifacts: mine targeted the *plan text*; `adv_review5.md` targets the *implemented corpus*. The audit independently corroborates the architecture, the golden derivations, and the corpus encoding — and it **empirically confirms the exact risk class my N3 finding predicted** (FIND-001 is a plan-pin drift). Meanwhile, my doc-level findings (N1, N2, N4–N9) are untouched, because the corpus audit's scope excludes plan text, session mechanics, and params contracts entirely.

---

## 1. Prior-finding disposition against audit evidence

| Prior finding | Status after audit | Evidence |
|---|---|---|
| N1 replay-decoder wording contradiction | **Still open** | Audit never touches session replay/decoding — out of its scope. Plan text still self-contradicts spec §4. |
| N2 `OperationParams` example violates own binding | **Still open** | Params canonical JSON not examined by corpus audit. |
| N3 corpus pins incomplete | **Materialized + partially resolved** | FIND-001 is precisely the predicted drift mode. Zero-ctor/RestFrame forms were implemented sensibly (matching the prompt2 pins I said to import) but nothing *forced* that outcome. Retro-pin needed. |
| N4 `go/parser` ban scope | Behavior correct, doc text still overbroad | Audit §5 confirms non-test files are clean — i.e., the implementation adopted the correct interpretation; README still states it absolutely. |
| N5 `Pow(x,0)` nonzero-safe | Still open, unexercised | No golden trace needs it; remains unpinned. |
| N6 params strictness, N7 reorder scope, N8 Review/Challenge encoders, N9 README build claim | **Still open** | All out of corpus-audit scope. Note: audit §9's six-way cross-check list confirms N8 — no Review/Challenge encoder verification exists anywhere yet. |
| N10 Simplify test gating | Moot | Implementation complete, tests pass. |
| ⚪ unseen-section list | **Partially de-risked** | §17.8 cross-check, §19.3/§19.5/§19.6 pins, §28 item counts now have executable evidence. §15.13, §16, §26 remain unaudited (see R6). |

**Mutual corroboration:** the audit's three derivation traces (§6) match my independent hand-traces step-for-step — including the `Mul(0, c⁻¹)→0` admissibility via `c>0 ⊨ c≠0` and the `Sqrt(Pow(m·c²,2))` unwrap. Two independent methods converging on the same traces is meaningful evidence the simplify/entailment design is executable as pinned.

---

## 2. New findings from adversarially reviewing the audit itself

### 🟠 R1 — FIND-001 needs a disposition, and its sibling is unchecked
The audit correctly notes the spec is silent on non-`NewtonSecondLaw` source strings — but the *plan* pinned them (`"Classical Mechanics corpus"`), and the implementation deviated. Since spec > plan, either resolution is defensible — but the §41 hand-off invariant ("no remaining choices") means plan and code cannot stay divergent. Pick one authority and document it. **Critically: the audit's §3 relativity table drops the `Source` column that §2 has.** Given mechanics drifted, relativity sources (`"Einstein, 1905"` for EMR/MER, `"Special Relativity corpus"` otherwise) are unverified. Verify before closing FIND-001.

### 🟠 R2 — Corpus-status drift is not actually detected
The cross-check list (§9: expr, dimension, kind, provenance_status, assumptions, source) **omits `corpus_status`**. The corruption row "Change corpus status → MintObject rejects invalid corpus status" only proves invalid-*enum* rejection — flipping `established`→`contested` (a valid value) in the manifest would pass every cited detector. This is exactly the human-curated axis the spec says must not be revised; the baseline pin needs a mechanical detector (item-level comparison, or a framework-level `corpus_status == "established"` assertion). My first-round slice-O inventory had "corpus status preserved" — it vanished somewhere between plan revisions; restore it.

### 🟠 R3 — The "add undocumented trusted relation" row overclaims
Detection as described only covers the **manifest→constructor** direction (extra manifest item → unresolved or mismatched map entry). The reverse — a new exported domain constructor minting `ESTABLISHED` objects with **no manifest entry** — is trusted physics invisible to every cited check. Also a mechanism error in the audit: `DisallowUnknownFields` blocks unknown *fields*, not extra *items*. Fix is cheap: pin the constructor allowlist in the export-surface audit (expected manifest constructors + the four pinned non-manifest constructors), making the manifest↔constructor relation a checked bijection. This is test content in existing files — no REQ-003-01 impact.

### 🟡 R4 — `KineticEnergy` DERIVED/NONE was never pinned anywhere
The implementation's choice is sensible (spec §13.4's `NONE` default for non-corpus artifacts), but no document pinned it — same gap class as N3. Retro-pin it, including Kind, so the next parameterized constructor doesn't re-litigate it.

### 🟡 R5 — Audit hygiene nits
- FIND-002 says `expandCall` is in `ops/transform.go`; §7 says `simplify.go:392-398`. One is wrong (likely a typo — harmless, but fix for the record).
- The corruption battery appears **analytically derived, not executed**. Prompt2 Phase 22 required introducing corruptions in a disposable copy and observing detection. The analytical table is a reasonable proxy for most rows, but R2/R3 show precisely why analytical derivation fails — run the battery.
- The Lorentz-limit non-hardcoding argument ("result carries `speed_of_light_positive`, a hardcoded `1` would lack this") is weaker than it looks — a fabricated result object could carry assumptions too. The body-expansion source citations are the stronger evidence; keep both, weight them correctly.

### 🟡 R6 — Coverage boundary: corpus PASS ≠ conformance PASS
The audit's own scope leaves unaudited: §15.13.1 positional enforcement in `Session.Step`/replay, §16 ledger chain/17-step validation/tamper triple, §26 candidate validation and containment-artifact rejection, REQ-032-08..22, slices K/Q/R/S, and the determinism double-run. My prior hand-off requirement stands unchanged: **this corpus audit complements but does not replace the software-conformance audit, and complete (untruncated) spec + plan documents are still needed** — both remain truncated in the packet (spec at §15.9, plan mid-§9).

---

## 3. Updated verdict

| Layer | Verdict |
|---|---|
| Architecture | Unchanged — sound; audit adds kernel-separation evidence (FIND-002's classification is correct: `lorentz_factor` ID + `KindMinkowski` are spec-mandated closed-world content, not leakage) |
| Corpus encoding | **Independently PASS** per audit; credible within scope |
| Plan doc | Conditional GO **unchanged** — N1–N3 text fixes still owed |
| Implementation | One plan-pin divergence (FIND-001) requiring disposition; R1–R3 as cheap hardening |

**Bottom line:** the review holds in full. The audit's PASS is genuine for what it examined, and it converted several of my ⚪ items into verified ones. But it also validated the review's central thesis in the most concrete way possible: **cross-checks enforce consistency, not content — any plan pin without a mechanical detector will drift** (FIND-001), and the two drift directions it missed (corpus status R2, extra constructors R3) are the same lesson applied to its own test suite. Close R1–R3, disposition FIND-001, apply the still-open plan fixes, and run the conformance audit over the session/candidate region before hand-off.
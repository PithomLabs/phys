# Consolidated Review: v3 (line-by-line re-audit) + `GR_IMPLEMENTATION_PLAN_20261004.md` (v4 hardening plan)

## Verdict

**Architecture: unchanged and blocker-free. But the v4 hardening plan, as written, would introduce two new contradictions with the frozen spec and its own pinned script — both guaranteed agent stalls — and its stale-reference sweep misses three dangling refs I can name by line.** The single most serious finding: **the F2 `Subst` pin forbids the one operation the pinned 14-step script cannot run without** (function→composite substitution), and the v4 author's own F4 needs a substitution form F2 also forbids. Second: **the GR-8 `Identify` expectation (`IDENTIFIED`) contradicts the frozen contamination law** — the operands are forced-HYPOTHESIS bridges, so `IDENTIFIED` is unreachable through the public API, and the agent hits a wall with no legal workaround. Both are one-paragraph fixes inside the v4 work items; neither touches architecture.

---

## 1. v3 findings → v4 disposition

| v3 finding | v4 item | Verdict |
|---|---|---|
| N-A Pin B target file doesn't exist | H1 → `ops/negative_test.go` | ✅ **Verified**: file is in the frozen 39-set; membership check will pass. |
| N-B `Subst`/`Reduce` undefined | F2 defines `Subst`; `Reduce` still undefined | ❌ **F2 is wrong — see V4-1.** `Reduce` remains unpinned (define or delete from D31/§7). |
| N-C bridge assumption filter/keys | E4 pins the filter; **keys still unpinned** | 🟠 Half-resolved — keys are hash-bearing (bridge ID), still cross-agent nondeterminism. |
| Nit: `§36 footer` / `20.12` dangling | B1 | ✅ |
| Nit: duplicated "Pin A:" / "fail loudly" | **No work item** — §25 lock claims fixed, v3 text still has both, A1 copies verbatim | 🟡 persists unless patched |
| Nit: Rational zero-form/overflow | not addressed | 🟡 open |
| Nit: Add/Mul arity unwrap post-normalize | not addressed | 🟡 open |
| Nit: GRStep `OperationID` vocabulary | not addressed (D2 pins schema, not the ID namespace) | 🟡 open |
| Nit: `symbolic/weakfield.go MAY` | D4 | ✅ sole home pinned |
| Nit: √\|g\| exclusion | implicit via pinned Γ formula + STOP | 🟡 optional one-liner |
| Golden-sign label (`R_thth` only after B=1/A) | G1 pins the **general** form `1−1/B+rB'/(2B²)−rA'/(2AB)` | ✅ **Independently re-derived and confirmed** — and it's an improvement: certification can now run in the (A,B) system |

---

## 2. 🔴 New findings in the v4 plan itself

### V4-1 — F2's `Subst` pin forbids the operation the pinned script requires (contradicts D28 + F4)

Trace the script's substitution needs mechanically:

- **Step 8** (flatness ⇒ `k1=1`): the only mechanical route is substituting the *asserted boundary limits* `A→1`, `B→1` into `AB=k1` — **UnknownFunction→Rational substitution**. (F4's "`u=1/r, u→0`" cannot work here: at step 8, A and B are still unknown — there is no closed form to rewrite in u. The u-mechanism only applies to closed-form expressions.)
- **Steps 10–11**: `R_θθ = 0` with `B=1/A` ⇒ `(rA)'=1` — requires **B→1/A** (function→composite) plus **B′→−A′/A²** (its derivative atom, via chain rule).
- **Step 12**: certifying candidate `A→1+k2/r` inside `(rA)'−1` — **function→composite**.
- **Step 14**: vacuum check with both closed forms — **function→composite, twice**.
- **F4/F3 themselves**: `u→0` is **Symbol→Rational** (not Symbol→Symbol); `mu=eps·r` is Symbol→composite.

F2 ("allow Symbol→Symbol; composite → `UnrepresentableSubstitutionError`") blocks every one of these, and its own example (`A(r)→A(x)`) describes argument-relabeling — a different operation from what the script needs. **Corrected contract to pin instead:**

```text
Subst(e, target, replacement):
  target ∈ { Symbol, UnknownFunction atom (Name, Args, Order=all-zero) }
  replacement = arbitrary GRExpr
  composite target (Add/Mul/…) → UnrepresentableSubstitutionError
  UnknownFunction f → g:  all derivative atoms f^(k) (same Name+Args)
      are simultaneously rewritten to Diff^k(g, arg) via the local
      differentiator (max order 2 enforced; NonZeroEvidenceError
      where a negative-power base lacks evidence)
  Symbol substitution never rewrites inside UnknownFunction.Args
      (function identity is pinned to declared argument symbols)
  result: Normalize
```

Without this, D28's certification primitive does not exist and the agent must either improvise it (banned) or STOP at GR-7 (guaranteed).

### V4-2 — GR-8 `Identify` expectation contradicts the frozen contamination law

v3 §23 (copied verbatim into v4 per A1) says the correspondence produces "`KindRelation`/**`IDENTIFIED`**/`NONE`" and "expected `IDENTIFIED`/`NONE`." But the operands are, by v3's own §23, "**HYPOTHESIS** kernel-representable scalars" — and frozen §13.2 is explicit: *Identify with a HYPOTHESIS operand ⇒ HYPOTHESIS output*. There is **no public path to non-HYPOTHESIS operands** here (the bridge is the only raw mint, forced HYPOTHESIS; `k2` and `−2mu` are not domain constructors). So the pinned expectation is unreachable: the agent's test fails, and no conforming fix exists. **Patch: expected result = `Kind=Relation`, provenance `HYPOTHESIS`, corpus `NONE`** — which is also the epistemically correct outcome (a mass correspondence identified from hypothesis-grade GR scalars is hypothesis-grade). H3 must state the provenance explicitly, and the §23 "expected IDENTIFIED" text must be patched in v4.

### 🟠 V4-3 — F3 silently supersedes frozen D34; no precedence clause exists

F3 redefines `Truncate` from D34's general "retains O(eps⁰..¹)" degree-truncation to "closed-form Schwarzschild only + `mu=eps·r` + single rule `Pow(1+u,−1)→1−u`." The rule checks out (u=−2ε ⇒ 1/(1−2ε)→1+2ε ✓), and the restriction is *good* — but v3 is byte-frozen and v4 patches create conflicts with a document labeled "normative execution contract." v4 must add one explicit clause: **"v4 supersedes v3 wherever they differ; v3 is preserved for lineage only"** — and an itemized supersession note for D34 (keeping its no-`APPROXIMATED` and honesty-record semantics).

### 🟠 V4-4 — C1's "distribution only where explicitly required" is decision language

The ZeroTest certifications (step 5: `R_tt/A + R_rr/B = (AB)'/(rAB²)`) require common-denominator assembly, which requires distributing products over sums. Pin the permitted context in one sentence: *distribution is permitted solely within common-denominator normalization and is bounded by the workload's expression depth; all other expansion is rejected.* Also pin Normalize as **iterate-to-fixpoint with a pass cap** (exceeding the cap ⇒ `UNDECIDED`, never loop).

### 🟠 V4-5 — B1's stale-ref sweep is incomplete; A2's grep list misses the actual offenders

By line: **§5** ("Full execution tables are in §27"), **§8** ("bridgeability census (§27)"), and **§20 item 2** ("tables (§27)") all dangle — the tables live in **§24**; v3 has no §27. None is in B1's list. Likewise A2's banned-phrase grep misses the decision language actually present: "preferably also `a_r`" (§23), "(if any)" (§23 Identify), "where representable" (§21), "where included" (D34). Extend both lists; make `a_r` a definite include (it is one `Diff` call), delete "(if any)" (H3 settles the question), and rewrite the §21 parenthetical per E4.

---

## 3. 🟡 Remaining pins (one pass each)

1. **D40's "8 artifact kinds" vs Appendix D's 8 schemas disagree**: v3 §26 lists `Index` (alongside `IndexSlot`); v4 D.2–D.8 drops `Index` and adds `BridgePayload`. Reconcile the canonical list (recommend: `GRExpr, IndexSlot, Tensor, Metric, Connection, Curvature, GRStep, BridgePayload`).
2. **`adversarial/` is a test-only package** — a directory containing only `_test.go` fails `go build ./...` ("no non-test Go files"). phys-gr has no file freeze, so: distribute SW tests adjacent to their targets and keep `adversarial/` evidence as .md, or add one non-test file. Pin the choice and pin phys-gr's gate as `go test ./... && go vet ./...`.
3. **GRStep `OperationID` namespace** (hash-bearing): pin e.g. `gr:Diff`, `gr:Normalize`, `kernel:<op>`.
4. **Rational edge pins**: zero form `0/1`; `num`/`den` overflow → named error, never silent wrap.
5. **Post-normalize arity invariant**: rational collection can drop an `Add`/`Mul` below the ≥2 floor — pin the unwrap rule (1 term → the term; 0 → `Rational 0`).
6. **G1 notation**: bare `R` for the θθ component collides with Ricci scalar `R` (§23) — use `R_θθ` throughout.
7. **v4 item count**: header says "28 hardening items"; the work items enumerate 30 (A2+B2+C6+D4+E4+F4+G2+H3+I3). Reconcile for audit hygiene.
8. **SI correspondence scope**: H3 pins only `k2 ↔ −2mu` (geometric). State that the SI relation `k2 = −2GM/c²` is dimension-audit-only (GR-8d), no `Identify`.
9. **Pin A "restore broken fixture" (§16)**: classify it as a §17 disposable-copy mutant (in-place restoration of the authoritative tree would violate the freeze protocol it sits next to).
10. **D35 mechanics**: pin that the `Compare` relation is recorded via `Session.Step` (`Kind:"compare"`, `Operator:"eq"`) before `Identify`, lifecycle `Draft → Identify → Conclude`, no `Seal` (DraftMetadata without a Hypothesis field is legal precisely because Seal is never called).
11. **D=4**: state it once explicitly (it is currently implicit in `4^Rank` and `[4][4][4]`).

---

## 4. Verified-correct in v4 (no action)

H1's target file exists in the frozen tree ✓. G1's general formula independently re-derived ✓ (and correctly specialized: B=1/A ⇒ `1−A−rA'` ✓). F3's single rule verified for `u=−2ε` ✓ and the `g_tt` correspondence closes exactly (`−(1−2ε) = −(1+2Φ/c²)` with `Φ/c²=−μ/r=−ε` ✓). E2/E3 are consistent with v3 §21 once read as caller-driven decomposition (ToCoreExpr all-or-nothing per call) ✓. Payload field order pinned, `schema_version` exclusion justified ✓. C6 sidecar-only hashing ✓. D3's `Connection{[4][4][4], ChartID}` with Γ-index order documented ✓. B2's footer set ✓. H2's GR-3a sequence (bridge → `ops.Differentiate` → replay, never `Session.Step`) matches D33 ✓ — note its dimension closure works: geometrized-Dimensionless probe keeps the Differentiate dimension rule trivially satisfied.

---

## Bottom line

The consolidation target is reachable, but **not by executing v4 as written**: V4-1 and V4-2 are contradictions the executing agent cannot resolve without either violating a pin or stalling to human adjudication — precisely the outcomes the plan's own STOP machinery then triggers. Amend the v4 work items before execution: rewrite F2 to the corrected `Subst` contract (with derivative-atom auto-rewrite), patch the §23/H3 expected provenance to `HYPOTHESIS`, add the v4-supersedes-v3 clause, pin the bridge assumption keys and the distribution context, extend B1/A2's sweeps to the named offenders, and sweep the eleven 🟡 items. All are edits to the work-item list — none reopens D1–D40. With that single revision pass, `UNPINNED_IMPLEMENTATION_CHOICES = 0` becomes a claim the document has earned, and the coding agent inherits a contract in which every remaining discovery routes to evidence, STOP, or the Growth Gate — never to a decision.
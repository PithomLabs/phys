# Consolidated Review — v3 (`GR_IMPLEMENTATION_PLAN_20261003.md`) + v4 Hardening Plan

Read: v3 line by line against `specs_v2_3.md`; v4 work items line by line against every v3 finding; both against the four-prior-review trail. Target: no architectural blockers, no residual architectural choice for the coding agent.

---

## 1. Coverage check — every prior v3 finding maps to a v4 work item

I re-verified each finding from my v3 review against v4's work-item list. All are covered.

| v3 finding | v4 work item | Status |
|---|---|---|
| B1 §20.12 references non-existent §36 | v4-B1, B2 | ✓ |
| B2 GRExpr canonical JSON unpinned | v4-D1, D2 | ✓ |
| B3 GRStep + other artifact JSON unpinned | v4-D2, D3 | ✓ |
| B4 Bridge payload canonical form unpinned | v4-E1 | ✓ |
| B5 "exact normalization proof" undefined | v4-F1 | ✓ |
| B6 §25 hashing = §26/Appendix C stale | v4-C4 | ✓ |
| H1 R_thth "(ansatz)" tag misleading | v4-G1 | ✓ |
| H2 Session.Identify operands unspecified | v4-H3 | ✓ |
| H3 Test files absent from §8 | v4-D4 | ✓ |
| M1 Fields 19/20 formats unpinned | v4-C6 | ✓ |
| M2 Canonicalization location unnamed | v4-C5 | ✓ |
| M3 Non-bridgeable assumption rule | v4-E4 | ✓ |
| M4 Vacuum relation representation | v4-G2 | ✓ |
| O1 Footer claim false | v4-B2, I2 | ✓ |
| O2 Timestamp exclusion not cross-noted | v4-C6 | ✓ |
| O3 §8 tree incomplete | v4-D4 | ✓ |

Confirmed reading of v4's fix mechanics:

- **G1's formula is correct.** `R_θθ = 1 − 1/B + rB'/(2B²) − rA'/(2AB)`; under `B = 1/A` it reduces exactly to `1 − A − rA'`. Verified algebraically.
- **H1's Pin B target change is correct.** `ops/differentiate_test.go` is not in the frozen `specs_v2_3.md §3` tree; `ops/negative_test.go` is. v4 correctly redirects.
- **E4 correctly closes M3.** `A(r) != 0` cannot reach the bridge because `A(r)` is `UnknownFunction`. v4 pins this as GR-replay-only — no silent drop, no silent conversion.
- **F3's `Pow(1+u,−1) → 1−u` rule is sufficient.** Weak-field `B = 1/A` at the pinned order requires exactly this single rule; `A = 1 − 2ε` needs none.

---

## 2. Residual items not addressed by v4

None of the following are architectural blockers, but under the plan's own "zero unpinned choices" standard they need pinning in v4's execution, or explicit acknowledgment.

### R1. Hash concatenation semantics for `PreviousStepHash`

v3 §26 and v4-C2 give:

```
CurrentStepHash = SHA-256(PreviousStepHash || CanonicalJSON(step-without-CurrentStepHash))
```

`PreviousStepHash` is defined as a 64-character hex string (genesis = 64 zeros). The `||` is undefined:

- **(a)** concatenate the *string* `PreviousStepHash` (ASCII hex) with the UTF-8 bytes of the canonical JSON, then SHA-256; or
- **(b)** decode `PreviousStepHash` from hex to 32 raw bytes, prepend those bytes, then SHA-256.

Two implementations that agree on everything else will produce different step hashes. Pin one in v4 (recommend **(a)** — the string form is what appears in `Step.PreviousStepHash` and what appears on replay).

### R2. Test package placement rule not pinned

v4-D4 says the §8 tree must list "adjacent `_test.go`". It does not say whether each test is `package symbolic` (internal) or `package symbolic_test` (external). Go permits both, and the frozen `phys` plan (`plan10_v2_3.md` §2) makes this a first-class distinction. Pin the rule (recommend: internal by default; external only where the test asserts public-surface behavior).

### R3. `GR_WORKLOAD_IMPLEMENTED=0` dropped from footer

v3's footer has five values; v4-B2 lists four, omitting `GR_WORKLOAD_IMPLEMENTED=0`. If intentional (status flag, not invariant), state so. If accidental, restore.

### R4. "11-point checklist" in v4-I2 is not enumerated

v4-I2 says the self-audit runs against an "11-point checklist" and then lists roughly seven comma-separated items with ambiguous grouping. Name the eleven points explicitly, or change "11-point" to the actual count.

### R5. `kind` case in GRExpr JSON not pinned explicitly

v4-D1 uses `symbol{name}`, `rational{num,den}`, … suggesting lowercase `kind` values matching the frozen `specs_v2_3.md §10.2` convention (`"kind":"symbol"`). This is strongly implied but not stated. One sentence closes it.

### R6. `Index` vs `IndexSlot` in the canonical-JSON coverage list

v3 §26 lists eight types: `GRExpr/Index/IndexSlot/Tensor/Metric/Connection/Curvature/GRStep`. v4-D2 lists eight: `GRExpr` (D.1) + `IndexSlot/Tensor/Metric/Connection/Curvature/GRStep/BridgePayload` (D.2–D.8). `Index` is silently replaced by `BridgePayload`. If `Index` was a typo for `IndexSlot`, say so. If distinct (e.g., a string label), v4 is dropping it.

### R7. `Chart` and `ConventionSet` not in v4-D2's canonical-JSON list

Both appear in v3 §8 as GR-local structures. Neither appears in v4-D2. If they do not receive standalone canonical JSON (because they are embedded in `GRStep` inputs or are always the pinned chart), state this. Otherwise their canonical forms are unpinned.

### R8. `Subst` scope

v4-F2 pins `Subst` to `Symbol → Symbol` only. That is a bounded, defensible choice — but the plan should say *what uses it*. If nothing in GR-0…GR-8 exercises `Subst`, either remove it from the D31 name list or record why it exists.

### R9. "candidate-and-certify" mechanics not pinned

v3 §23 steps 7 and 12 say "candidate-and-certify" without defining the procedure. v4 does not touch this. Under the plan's own target ("zero unpinned choices"), pin the procedure: e.g. "the certifier substitutes the candidate expression into the residual, runs `Normalize`, then `ZeroTest`; success requires `ZERO`; `UNDECIDED` fails closed."

### R10. v4-B1 conflates two fix categories

v4-B1 lists "current `20.12`, `25` locks line, `Appendix C` checklist" as targets for the §36/§20.12 reference fix. But §25's stale reference is the `hashing = §26/Appendix C` one, which is separately handled by v4-C4; and v3's Appendix C contains no §36 reference. The v4-B1 list is therefore either over-broad (mixing in C4's target) or mis-scoped. Tighten.

### R11. Chart canonical identity not pinned across components

The plan says the canonical chart is `spherical-static`. But nothing pins that the metric, connection, curvature, and every tensor all carry `ChartID = "spherical-static"` throughout the workload. If any component may carry a different `ChartID`, the cross-chart validation rules (v3 §8, v4-D3) are exercised. If not, state "single-chart workload; cross-chart behavior is only exercised by the adversarial suite."

None of R1–R11 reopens an architectural decision. Each is a one-line pin in v4's execution. R1 and R2 are the highest-leverage: R1 affects replay determinism (a hard invariant) and R2 affects the tree's own completeness claim.

---

## 3. Verdict

**v4 as a work plan is complete against every prior v3 finding. The resulting v4 document, if executed as v4's work items specify, will not introduce any new architectural blockers, and will not leave any of the previously-identified architectural choices open to the coding agent.**

The eleven residual items above are all one-line disambiguations, not re-decisions. They are the class of "unpinned implementation choice" the plan exists to eliminate; folding them into v4's work items is required to reach the plan's own stated bar:

> `ARCHITECTURAL_BLOCKERS = 0`
> `UNPINNED_IMPLEMENTATION_CHOICES = 0`

Recommend:

1. **Execute v4's work items A–I** as written — they cover every prior finding.
2. **Add R1–R11** to v4's item list before execution (they are additive; none conflicts with A–I).
3. **After v4 executes**, re-run the self-audit mechanically against the same checklist that produced the v3 findings, and only then set `UNPINNED_IMPLEMENTATION_CHOICES = 0`. v3's footer claimed zero while six blockers existed; the audit procedure itself is what needs the verification, not just the counter.

With R1–R11 folded in, v4 is the first version that genuinely meets the target:

> **The coding agent discovers mathematics, not architecture. It may hit implementation bugs, performance limits, or genuine kernel-boundary evidence — but it does not choose what belongs in the kernel, how GR interoperates with `phys`, what the symbolic representation is, how replay hashes are computed, what conventions apply, or what constitutes promotion.**
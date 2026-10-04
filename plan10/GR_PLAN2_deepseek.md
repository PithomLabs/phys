# Adversarial Review — `GR_IMPLEMENTATION_PLAN_20261002.md` (v2)

Read alongside: `specs_v2_3.md`, `plan10_v2_3.md`, the v1 plan, and the six-review consolidation. Target unchanged: no architectural blockers, no residual architectural choice for the coding agent.

---

## 0. Bottom line

v2 is a genuine improvement. It resolves the v1 §6/§7 bridge contradiction, corrects the silent-wrongness semantics, retires the two imprecise phrasings, renames the terminal state, and moves the "explicit in phys-gr" pins closer to the plan surface.

**But v2 does not close the gap.** Of the fourteen items the consolidated review asked for beyond D1–D18, v2 fully resolved **six**, partially resolved **three**, and left **five** substantively open — plus it introduced **four new findings**. Two of the open items are genuine architectural blockers under the plan's own target (curvature sign conventions; bridge `Dimension` value), and two more are high-severity (GR-local replay contract; `UnknownFunction` shape). The plan is not yet at the "no architectural choice remains" bar it claims in its header.

Do not authorize PASS 0 yet.

---

## 1. What v2 correctly resolved (do not re-litigate)

Directly checked against the prior consolidated review:

| Prior finding | v2 resolution | Verdict |
|---|---|---|
| C3-1 (§6.2 contradiction with D5–D7) | §6 fully rewritten as "Locked bridge" | **Resolved** |
| C3-2 (`Limit` for Newtonian limit) | §7 row rewritten: "NOT used for GR-8" | **Resolved** |
| C3-3 ("does not bypass `phys.Differentiate`") | §7: old phrasing explicitly retired | **Resolved** |
| C3-4 (§13/§18 terminal-state contradiction) | §18 renamed to `LEVEL-2-CANDIDATE-PENDING`; §13 consistent | **Resolved** |
| C3-5 (42-file vs 39-file) | §3.3: "39 + 3 = 42"; cites `plan10/adv_review11.md` | **Resolved** |
| C3-6 (silent-wrongness Probe 2) | §12: Probe A/B/C corrected | **Resolved** |
| C3-7 ("no replace into kernel internals") | §4: retired as imprecise | **Resolved** |
| C3-8 (§15 fields 18/20 formats) | §15: enumerated | **Resolved** |
| D1–D18 | Confirmed against `specs_v2_3.md` | **Adopt as-is** |

These are closed. The remainder of this review does not revisit them.

---

## 2. Prior findings still open — v2 did not close them

### BL-1. Curvature sign conventions are still deferred, not pinned

v2 §8:

> "GR curvature convention **explicitly defined in `phys-gr`** (kernel acquires no GR curvature semantics)."

This is the same deferral as v1. It says *where* the convention will be recorded (in `phys-gr`), not *what* it is. The coding agent cannot derive Schwarzschild from `A(r)`, `B(r)` without a sign convention, because:

- The Christoffel formula's sign determines all subsequent curvature signs.
- The Riemann definition (`R^ρ_σμν = ∂_μ Γ − ∂_ν Γ + ΓΓ − ΓΓ` vs its negative) determines Ricci.
- The Ricci contraction (`R_μν = R^ρ_μρν` vs `R^ρ_ρμν`) determines the Einstein tensor.
- The Einstein tensor definition (`G_μν = R_μν − ½ g_μν R`) determines the vacuum equation's form.
- The Schwarzschild derivation's intermediate `A(r)`, `B(r)` equations depend on the choice.

The plan must pin the exact formulas. My prior C2-2 recommended the MTW/Wald set (compatible with the existing `−+++` signature). v2 does not adopt any of it. **This is a blocker under the plan's own target.**

### BL-2. Bridge `Dimension` value is still undefined

v2 §6 diagram lists:

```
core.Object
  Kind = Expression
  explicit Dimension
  explicit assumptions
```

But v2 never says what "explicit Dimension" is. The plan's own §8 (Units) has two subtracks with different dimensional requirements:

- Geometrized tensor track (`G = c = 1`): every scalar is formally dimensionless, but `r` still has length-like character.
- SI Newtonian subtrack (`GM/(c²r)`, `GM/r²`): scalars carry SI dimensions.

The plan does not state which rule applies to the bridge. Two mutually incompatible readings are open:

- **(a)** Every geometrized bridge uses `core.Dimensionless()`. SI audit re-bridges with SI dimensions on a separate path.
- **(b)** Bridge `Dimension` reflects the physical dimension of the subexpression regardless of the geometrized-unit convention.

The coding agent must choose. My prior C2-3 recommendation was (a). v2 does not adopt it. **This is a blocker under the plan's own target.**

### HI-1. GR-local replay contract is still incomplete

v2 §7 gives:

```
GRStep: StepID / OperationID / InputCanonical / ParamsCanonical
        / OutputCanonical / CurrentHash   (SHA-256 over canonical JSON)
```

Prior C1-2 asked for seven specific pins. v2 provides:

| Required | Status in v2 |
|---|---|
| Genesis hash | Not stated |
| `StepID` format (deterministic from index) | Not stated |
| Hash algorithm | SHA-256 (present) |
| Canonical JSON encoding rules | Not stated |
| Field order of `GRStep` | Implied by listing; not declared as fixed |
| Bridge hash as a separate `GRStep` or a field | Not stated |
| `Decode(Canonical(x)) == x` verification | Not stated |

Six of seven still open. The coding agent cannot implement `GRStep` without choosing. **High severity** — this is the L1 replay foundation the plan claims (D11) to have pinned.

### HI-2. `UnknownFunction` shape is still undefined

v2 §6:

> `UnknownFunction("A", r)`, `UnknownFunction("B", r)` with `dA/dr = A'(r)`, `dB/dr = B'(r)` (+ higher derivatives as curvature requires).

Prior C2-4 asked for:

| Required | Status in v2 |
|---|---|
| `Args` type — `[]Symbol` or `[]GRExpr`? | Not stated (examples show a bare symbol) |
| Arity fixed at construction | Not stated |
| Representation of `A'(r)` — `UnknownFunction("A'", r)`? A dedicated `Derivative` node? | Not stated |
| Derivative of `A(r)` w.r.t. a variable other than `r` | Not stated |

The plan names the node but not the node's shape. The coding agent must choose. **High severity** — this is the unknown-function facility `Schwarzschild` derivation depends on.

### HI-3. Differentiation base-case set not pinned

Prior C2-6 asked for the complete base-case table:

| Case | v2 status |
|---|---|
| `Symbol(x)` vs `Symbol(x)` → `1` | Not stated |
| `Symbol(x)` vs `Symbol(y)`, `x≠y` → `0` | Not stated |
| `UnknownFunction(f,[x])` vs `Symbol(x)` → `f'(x)` | Stated by example |
| `UnknownFunction(f,[x])` vs `Symbol(y)`, `y∉args` → `0` | Not stated |
| `Sin(x)` vs `Symbol(x)` → `Cos(x)` | Stated |
| `Sin(x)` vs `Symbol(y)` → `0` | Not stated |

Two cases resolved by example; four still open. The last (`Sin(x)` vs non-`x`) matters for the angular-metric sector. **High severity.**

### HI-4. Chart representation not pinned

Prior C2-1 asked for a concrete chart shape. v2 §8 lists `coordinates/ chart, symbols, assumptions/regime, conventions` but does not define a struct. The coding agent must choose whether `ChartID` is an enum, a string, or a struct with fields. **Medium-high severity** — `ChartID` appears in §8's tensor engine (D13) and every component is chart-scoped.

### MD-1. Independent reviewer identity still unnamed

Prior C1-3 asked who performs the independent classification. v2 §15 field 19 requires "agent classification + independent classification + disagreement resolution" but does not name the independent reviewer. **Medium severity** — the coding agent cannot fill field 19 without knowing whether the reviewer is a human, a second AI agent, or a checklist.

### MD-2. Bridge failure error not named

Prior C2-9 asked for a named GR-local error when `ToCoreExpr()` cannot project (trig, unknown function present). v2 §6 says:

> "No fake conversion through placeholder symbols."

but does not name the error. The coding agent must choose. **Medium severity** — the error appears in every projection attempt.

### MD-3. Schwarzschild endpoint classification not pinned

Prior C2-5 asked whether the endpoint is `Session.Identify` or GR-local substitution. v2 §5 says "PRIMARY derives... SECONDARY verifies..." but does not say which mechanism produces the final equation artifact. **Medium severity** — affects session evidence logging vs GR-local replay logging.

### MD-4. Bridge decision procedure still imprecise

Prior C1-1 asked for eager/lazy/hybrid. v2 §6 says "bridged into a **temporary** generic `core.Object` **ONLY when** a frozen kernel operation is deliberately being exercised" — this reads as **lazy**, but "deliberately being exercised" is not a decision procedure. When a scalar is constructed with `r` and `t` but never fed to a kernel op, does it become a `core.Object` or stay GR-local? **Medium severity.**

### MD-5. `symbolic/` file list still incomplete

v2 §8 adds files:

```
expr.go, differentiate.go, normalize.go, bridge.go, weakfield.go
```

Prior C2-8 (from my review) recommended also:

```
canonical.go, replay.go, construct.go (or fold into expr.go), zero.go
(if not folded into normalize.go), expr_test.go, differentiate_test.go
```

Appendix C declares GR artifact hashing (`GRExpr`, `Index`, `Tensor`, `Metric`, `Connection`, `Curvature`, `GRStep`) but §8 provides no file for the canonical JSON / SHA-256 implementation. Similarly §7 requires a `GRStep` replay engine but §8 provides no `replay.go`. Either fold into existing files (state so) or add them. **Medium severity** — the tree must match the plan's own capabilities.

### MD-6. PASS0 executor and evidence-record classification still unspecified

Prior C1-4 and C1-5 both stand. v2 §16 adds a post-PASS0 baseline record but does not specify:

- whether PASS 0 and the GR workload run under the same agent invocation,
- whether the Growth Evidence Record is a canonical `phys` artifact (it contains a UTC timestamp per §15 field 20, which `specs_v2_3.md` §31 forbids in canonical artifacts — so either the record is declared non-canonical, or the timestamp must be removed).

**Medium severity** for both.

### MD-7. Three authorized documentation files not named by path

v2 §3.3 cites `plan10/adv_review11.md` for the 42-file freeze, but does not inline the 3 documentation filenames. If `adv_review11.md` contains them, cite that fact; if not, list them. Otherwise PASS 0's pre-edit membership check relies on the coding agent opening a document the plan has not guaranteed exists in an accessible form. **Medium severity** — technically a PASS 0 pre-requisite failure mode.

---

## 3. New findings introduced by v2

### N-1. §6 uses `GRExpr / GRScalar` but `GRScalar` is undefined

v2 §6 diagram:

```
GRExpr / GRScalar
```

But §6, §8, and Appendix C reference only `GRExpr`. `GRScalar` is undefined and used once. Either define it (perhaps as an alias or a distinct scalar-only subtype) or remove it. **Low-medium severity** — but this is exactly the class of ambiguity that the plan exists to eliminate.

### N-2. §8 "GR curvature convention explicitly defined in phys-gr" reads as a deferral

Even if BL-1 is fixed by inlining the formulas, the current sentence's form invites the coding agent to believe the pin lives elsewhere. Reword to "GR curvature convention pinned in this document (and mirrored into `phys-gr/docs/conventions.md` as implementation material)" or similar. The plan cannot cite itself for a pin it does not state. **Sub-finding of BL-1.**

### N-3. §7's GR-3a probe does not say whether it runs through `Session.Step`

v2 §7:

> "GR-3a — Kernel boundary probe: kernel-representable `Pow(r,-1)` via HYPOTHESIS `core.Object` → `phys.Differentiate(...)` → `UnsupportedOperationError`"

Does the probe call `ops.Differentiate` directly, or `Session.Step(label, "differentiate", ...)`? This matters because:

- A direct call produces no session ledger entry; the failure is recorded only in GR-local evidence.
- A `Session.Step` call produces a failed session step, which the session's own §16.6 processing must handle (rejected or recorded as failed?).

`specs_v2_3.md` §16.6 describes `Session.Step` calling `ops.Apply` and appending to the draft buffer. It does not describe a failure path (does a failed `ops.Apply` append a step or abort?). This is out of the GR plan's scope to fix, but the GR plan must say which interface the GR-3a probe uses. **Medium severity.**

### N-4. §15 field 14 refers to "GR artifacts §20" but §20 is Definition of Done

v2 §15 field 14:

> "artifact compatibility (canonical form / hash — GR artifacts **§20** + `core.Object` hashes for bridges)"

But §20 in v2 is "Definition of Done." GR artifact hashing is in Appendix C. The cross-reference is wrong. **Low severity — a citation fix.**

---

## 4. Items fully resolved by v2 (for completeness)

To make sure nothing was missed, here is the complete list of what v2 closed and what remains, itemized:

**Resolved (8):** C3-1, C3-2, C3-3, C3-4, C3-5 (partial — see MD-7), C3-6, C3-7, C3-8.

**Still open (13):** BL-1, BL-2, HI-1, HI-2, HI-3, HI-4, MD-1, MD-2, MD-3, MD-4, MD-5, MD-6, MD-7.

**New (4):** N-1, N-2, N-3, N-4.

Total remaining issues: **17**.

---

## 5. Final consolidated decision list

The prior review's D1–D18 are confirmed and adopted as-is. The following pin-list must be added to v3 to reach the target. Each item is a single decision that must be stated in the plan text.

```
D19  Bridge decision procedure. Recommend: lazy/temporary — a
     GRExpr becomes a HYPOTHESIS core.Object only at the call site of
     a kernel op, and is discarded after reintegration. State this
     explicitly and say whether any scalar ever persists as a
     core.Object beyond a single op-call.

D20  Bridge Dimension. Recommend: core.Dimensionless() for every
     geometrized-unit bridge; SI Dimension only for the explicitly-SI
     Newtonian subtrack. State the rule; do not leave "explicit
     Dimension" undefined.

D21  Curvature sign conventions. Pin at plan level in MTW/Wald form
     compatible with −+++: Christoffel, Riemann, Ricci, scalar,
     Einstein, exactly. Reference §8.

D22  GR-local replay contract. Pin genesis hash (recommend 64 zeros,
     matching kernel), StepID format (`gr-step-000001`), canonical
     JSON encoding rules, `GRStep` field order, whether bridges are
     separate GRSteps or fields, and the `Decode(Canonical(x)) == x`
     round-trip requirement.

D23  UnknownFunction shape. Pin: `{Name string, Args []Symbol}` with
     fixed arity; `d f(args)/dx` for `x ∈ args` = `UnknownFunction(f',
     args)`; for `x ∉ args` = `Rational(0)`; higher derivatives as
     `f''`, `f'''`, etc.

D24  Differentiation base-case set. Pin the complete table:
     Symbol vs Symbol same → 1; Symbol vs Symbol different → 0;
     UnknownFunction vs arg → derivative; UnknownFunction vs
     non-arg → 0; Sin(x) vs x → Cos(x); Sin(x) vs y ≠ x → 0;
     Cos(x) vs x → −Sin(x); Cos(x) vs y ≠ x → 0.

D25  Chart representation. Pin: `Chart{ID string, Coordinates
     []Symbol, Order []int}`, immutable, canonical-JSON-hashable.

D26  Independent reviewer identity. Recommend: human maintainer for
     UNREPRESENTABLE and SILENTLY-WRONG; second AI agent for initial
     second-opinion classification of REPRESENTABLE-BUT-UNFAITHFUL.
     State the rule.

D27  Bridge failure error. Pin: `UnrepresentableKernelProjectionError`
     (GR-local), returned by `ToCoreExpr()` when trig or unknown
     functions are present.

D28  Schwarzschild endpoint. Pin: GR-local substitution produces the
     derived form; optional SI-bridged `Compare` produces the
     verification form. Never `Session.Identify` on the derivation
     path.

D29  symbolic/ file list. Add `canonical.go` (GRExpr/GRStep canonical
     JSON + SHA-256) and `replay.go` (GRLedger, Validate). Fold
     `construct.go` into `expr.go` if preferred; state the fold.
     Add `expr_test.go` and `differentiate_test.go`.

D30  PASS0 executor. Pin: PASS0 and GR workload may be the same agent
     if the split is recorded in the evidence log; otherwise separate
     invocations.

D31  Growth Evidence Record classification. Pin: it is a Level-1
     planning artifact, not a canonical phys artifact; `specs_v2_3.md`
     §31 determinism rules do not apply; the UTC timestamp field is
     therefore permissible.

D32  Three documentation files. Either inline their paths in §3.3, or
     state that `plan10/adv_review11.md` lists them and PASS0 must
     read that document to perform the membership check.

D33  GR-3a interface. Pin: does the probe call `ops.Differentiate`
     directly, or `Session.Step(label, "differentiate", ...)`? State
     the choice and, if session-mediated, the expected failure
     behavior at the session boundary (step not appended; error
     returned; session state unchanged).

D34  §6 terminology. Remove `GRScalar` or define it as an alias of
     `GRExpr`.

D35  §15 field 14 cross-reference. Change "GR artifacts §20" to
     "GR artifacts Appendix C."
```

Additionally, the §15 field 20 UTC timestamp rule needs a cross-reference to D31 so it does not appear to violate `specs_v2_3.md` §31.

---

## 6. Verdict

v2 is on the right trajectory. It is not finished.

- **Two blockers** (BL-1 curvature signs, BL-2 bridge Dimension) directly contradict the plan's own target statement, because they require the coding agent to make an architectural choice about mathematical content.
- **Two high-severity items** (HI-1 replay contract, HI-2 UnknownFunction shape) will produce ad-hoc implementations that cannot be reconciled later without rework.
- **Nine medium items** are recoverable in a single editorial pass.
- **Four new items** were introduced by v2 itself and are trivially fixable.

Applying D19–D35 produces a plan that meets the stated target:

> **A coding agent should have no remaining architectural choice while implementing GR. It may discover mathematical difficulty, implementation bugs, performance limits, or genuine kernel-boundary evidence — but it must not have to decide what belongs in the kernel, how GR interoperates with `phys`, what the symbolic representation is, how replay works, what sign conventions apply, or what constitutes promotion.**

Do not authorize PASS 0 on v2. Authorize on v3 with D19–D35 folded in.
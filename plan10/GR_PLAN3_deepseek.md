# Adversarial Review — `GR_IMPLEMENTATION_PLAN_20261003.md` (v3)

Read against `specs_v2_3.md`, `plan10_v2_3.md`, the v1/v2 plans, and the six-review consolidation. Target unchanged: no architectural blockers, no residual architectural choice for the coding agent.

---

## 0. Bottom line

v3 is a genuinely close successor. All eighteen prior findings I raised against v2 are addressed in form. Direct spot-checks against the frozen spec confirm the pins are consistent: MTW curvature formulas, `Dimensionless` geometrized bridge, 11-field `GRStep`, `UnknownFunction{Name, Args, DerivativeOrder}` with order sum ≤ 2, complete `Diff` table, `Chart{ID, Coordinates, Order}`, `UnrepresentableKernelProjectionError`, lazy bridge, GR-3a `ops.Differentiate`-direct, four terminal states, three named documentation files, executable evidence fields 18–20.

**But v3 does not meet its own header claim.** The document opens with `No architectural decision remains for the implementation agent` and closes with `ARCHITECTURAL_BLOCKERS = 0 / UNPINNED_IMPLEMENTATION_CHOICES = 0`. Six findings below contradict both statements. Four are blockers. Two are high-severity. Three are medium.

Do not authorize PASS 0 on v3. Do so on v4 with items B1–B6 and H1–H3 folded in.

---

## 1. What v3 correctly resolved (calibration, do not re-litigate)

Checked line-by-line against my prior v2 findings and the six-review consolidation:

| Prior finding | v3 location | Verdict |
|---|---|---|
| BL-1 (curvature sign conventions deferred) | §23 — Christoffel/Riemann/Ricci/scalar/Einstein pinned, MTW, `−+++`, `Λ=0` | **Resolved** |
| BL-2 (bridge Dimension undefined) | §22 — geometrized → `Dimensionless`; SI audit separate path | **Resolved** |
| HI-1 (GR-local replay contract) | §26 — 11 fields, 64-zero genesis, `gr-step-NNNNNN`, hash construction, round-trip law | **Resolved** |
| HI-2 (`UnknownFunction` shape) | §6 — `{Name, Args, DerivativeOrder}`, arity pin, sum ≤ 2, atom semantics | **Resolved** |
| HI-3 (differentiation base cases) | §7 — full `Diff` table with product/chain rules and dependency zeros | **Resolved** |
| HI-4 (chart representation) | §8 — `Chart{ID, Coordinates, Order}`, canonical ID `spherical-static` | **Resolved** |
| MD-1 (independent reviewer identity) | §10 — "second-agent review … agent classification + independent classification + disagreement record + human resolution" | **Resolved** |
| MD-2 (bridge failure error) | §6 / §21 — `UnrepresentableKernelProjectionError` | **Resolved** |
| MD-3 (Schwarzschild endpoint) | §23 — `Session.Identify` solely for the integration-constant correspondence | **Resolved** |
| MD-4 (bridge decision procedure) | §21 — lazy, at deliberate kernel-op call sites only | **Resolved** |
| MD-5 (`symbolic/` file list) | §8 — `expr.go`, `differentiate.go`, `normalize.go`, `bridge.go`, `weakfield.go` (fold pinned), `replay/replay.go` | **Resolved** (with H3 caveat) |
| MD-6 (PASS0 executor and evidence-record classification) | §16 — "PASS0 is a separate task; the GR agent never performs PASS0 changes"; §26 — canonical artifacts exclude review metadata and UTC timestamps | **Resolved** |
| MD-7 (three doc files not named) | §3.3 — `AGENTS.md`, `mechanics/README.md`, `relativity/README.md` | **Resolved** |
| N-1 (`GRScalar` undefined) | Removed from diagrams and text | **Resolved** |
| N-2 (§8 curvature-deferral wording) | §23 pins the formulas in-document | **Resolved** |
| N-3 (GR-3a interface) | §7/§33 — direct `ops.Differentiate`, never `Session.Step` | **Resolved** |
| N-4 (§15 field 14 cross-reference) | Corrected (no longer cites §20) | **Resolved** |

These are closed. The remainder of this review does not revisit them.

---

## 2. Blockers — the plan claims zero, the plan has these

### B1. §20 item 12 references a non-existent section

> §20.12: "Machine checklist values verified (§36 footer)."

The document ends at §26 + Appendices A, B, C. There is no §36. `specs_v2_3.md` §36 covers the `physvet` boundary — inapplicable. The reference is stale from an earlier draft. This is precisely the class of defect the plan's own audit claim is supposed to catch.

**Fix:** either inline the machine-checklist values in §20 (preferred — they are short: `PASS0_REQUIRED=1`, `ARCHITECTURAL_BLOCKERS=0`, etc.), or reference the actual closing block by name.

### B2. Canonical JSON schema for `GRExpr` is not pinned

§26 claims "fixed field order" but never states the field order. Compare with `specs_v2_3.md` §10.2, which pins the core `Expr` JSON exhaustively:

```
{"kind":"symbol","name":"E"}
{"kind":"rational","value":"1/2"}
{"kind":"add","terms":[...]}
{"kind":"mul","factors":[...]}
{"kind":"neg","expr":{...}}
{"kind":"pow","base":{...},"exp":"2/1"}
{"kind":"sqrt","expr":{...}}
{"kind":"call","function":"lorentz_factor","args":[...]}
{"kind":"relation","operator":"eq","lhs":{...},"rhs":{...}}
{"kind":"branch_set","target":{...},"branches":[...]}
```

GRExpr has 9 nodes (`Symbol`, `Rational`, `Add`, `Mul`, `Neg`, `Pow`, `Sin`, `Cos`, `UnknownFunction`). The plan gives only type signatures (`Symbol{name}`, `Rational{num, den}`, …) and no `kind` discriminator rule, no field order per node, no `UnknownFunction` key naming. The coding agent must choose:

- Is there a `"kind"` field? (kernel convention says yes; the plan is silent.)
- Does `UnknownFunction` serialize as `{"kind":"unknown_function","name":"A","args":[...],"orders":[...]}` or `{"kind":"unknown_function","name":"A","args":[...],"derivative_order":[...]}`?
- Are `Add.terms` and `Mul.factors` JSON arrays?
- What is the JSON for `Pow{base, exp int}` — `"exp":2` or `"exp":"2/1"`? (kernel uses `*big.Rat`; GR-local uses `int`.)

Without this pin, canonical bytes are undefined. Canonical bytes feed SHA-256. Hashes feed the `GRStep` chain. The replay round-trip law (§26: `Decode(Canonical(x)) == x`) cannot be tested against an unspecified schema. This is an architectural choice that the plan claims not to leave.

**Fix:** pin the JSON schema for all 9 `GRExpr` node kinds, in kernel §10.2 style, with explicit `kind` discriminators and field order.

### B3. Canonical JSON schema for `GRStep` and the other GR artifacts is not pinned

§26 lists `GRStep{StepID, Index, OperationID, InputCanonicals, InputHashes, ParamsCanonical, OutputCanonical, OutputHash, KernelBridgeHash, PreviousStepHash, CurrentStepHash}` and says "fixed field order," but does not say what the JSON looks like. Same problem as B2, one level up. Also unaddressed:

- Is `Index` an integer or a string?
- Are `InputCanonicals` and `InputHashes` arrays, and do they serialize as strings (hash) or nested objects (canonical)?
- `KernelBridgeHash` "empty when no bridge" — empty string, or omitted?
- Field order in JSON.

Same for `Index`, `IndexSlot`, `Tensor`, `Metric`, `Connection`, `Curvature` (§26 says all eight are covered; only the field lists in §8 hint at the shapes).

**Fix:** pin the canonical JSON for each of the eight GR artifact types.

### B4. Bridge payload canonical form is not pinned

§21:

> "Bridge ID: `gr/bridge/<hex SHA-256 of canonical bridge payload>`; payload = canonical `core.Expr` + `Kind` + `Dimension` + `Assumptions` + `Conventions`."

What is the canonical payload bytes? Concatenation of the five canonical serializations with what separator? Null bytes? Length prefixes? A wrapper JSON object? The SHA-256 is not well-defined until this is pinned. Two implementations that agree on all five canonical pieces but choose different separators produce different bridge IDs for the same bridge — the exact class of determinism failure the plan exists to prevent.

**Fix:** pin the payload form: e.g., a JSON object `{"expr":..., "kind":..., "dimension":..., "assumptions":..., "conventions":...}` with fixed field order, canonical JSON, UTF-8, SHA-256 over those bytes.

### B5. §13's "exact normalization proof" is not defined

§13:

> "`ZeroTest ∈ {ZERO, NONZERO, UNDECIDED}`; `ZERO` solely on exact normalization proof; `NONZERO` solely on direct structural proof; otherwise `UNDECIDED`."

"Exact normalization proof" is not defined. Candidates:

- `Normalize(expr) == Rational(0)`;
- `Normalize(expr) == Rational(0)` and the normalization trace is finite and free of `UNDECIDED`-bearing identities;
- Some weaker condition.

Without this, the coding agent decides what `ZERO` means. Every downstream vacuum-equation check (`ZeroTest(ZERO)` in §23) depends on it. `SILENTLY-WRONG` classification (SW-11 zero-controls) depends on it.

**Fix:** pin the definition. Recommended: `ZeroTest(e) = ZERO` iff `Normalize(e)` reduces, using only the pinned §13 rule list, to `Rational(0)`.

### B6. §25's "hashing = §26/Appendix C" is a stale cross-reference

§25:

> "Locks: ... valid cross-references (hashing = §26/Appendix C)."

In v3, Appendix C is `Self-Audit Record`, not hashing. In v2, Appendix C was `GR Artifact Hashing`. This cross-reference was not updated when Appendix C was repurposed.

**Fix:** `§26` alone, or a new appendix.

---

## 3. High-severity — the plan permits an architectural choice

### H1. §23's `R_thth = 1 − A − rA' (ansatz)` is misleading; the formula requires `B = 1/A`

The 14-step script (§23) computes Ricci components at step 4 with the general ansatz (independent `A`, `B`). The formula `R_θθ = 1 − A − rA'` is only valid **after** step 9 has established `B = 1/A`. The `(ansatz)` tag on the golden-sign line can be read as "valid for the general ansatz" — which it is not.

Under MTW with general `A`, `B` (from standard references, e.g. Carroll):

```
R_θθ = 1 − 1/B + r B'/(2 B²) − r A'/(2 A B)
```

Substituting `B = 1/A`:

```
R_θθ = 1 − A − r A'
```

The formula is correct **only** in the post-step-9 regime. §23 step 10 places the certification correctly (after step 9), but the header line "Golden signs: `R_thth = 1 − A − rA'` (ansatz)" should not claim generality.

Same issue for §23 step 5: `R_tt/A + R_rr/B = (AB)'/(r A B²)` is valid for general `A`, `B` under this ansatz (I verified), so the tag there is fine. Only the `R_thth` line needs correction.

**Fix:** change `(ansatz)` to `(after B = 1/A)` and, if desired, list the general-ansatz form alongside.

### H2. §23's `Session.Identify` use is not fully specified

§23:

> "`Session.Identify` use (if any): solely for the integration-constant correspondence (`k2 = −2GM/c²` SI / `k2 = −2mu` geometric) in the GR-8 recording segment, producing a `KindRelation`/`IDENTIFIED`/`NONE` relation artifact; exact call site recorded in the GR trace."

`Session.Identify(a, b, justification)` requires two `core.Object` operands with equal dimensions and compatible kinds. The plan does not say which two operands, or what the justification is.

Candidate readings:

- Identify between a derived `core.Object` expressing `k2` (in geometric units) and one expressing `−2mu`, both dimensionless.
- Identify between a derived geometric-units expression and its SI unit-expression, dimensionally reconciled by explicit dimension assignment on the bridges.
- Identify between two forms of the constant (`k2` derived and `k2` expected from the Newtonian correspondence).

These are mutually incompatible identifications, with different evidence traces. The coding agent must decide.

**Fix:** state the two operands explicitly and the justification template.

### H3. Test files not listed in §8; test placement not specified

§8 gives a directory tree but includes no `_test.go` files. §16 says "GR tests: golden states per pass, GR-3a/3b pair, …" but §8 does not name the test files, does not pin test package placement (`package physgr` vs `package physgr_test`), and does not say whether tests live adjacent to source (Go idiom) or in a separate tree.

The v1 plan listed `symbolic/expr_test.go` and `symbolic/differentiate_test.go`. v3 removed them from the tree without stating an alternative.

**Fix:** either list test files in §8 (adjacent, one per package, named `<source>_test.go`), or state the rule explicitly ("every `phys-gr` package has an adjacent `<pkg>_test.go`; external where authority/visibility is under test").

---

## 4. Medium — disambiguation

### M1. §15 fields 19–20 formats not pinned

Field 19: "`IndependentReview`: reviewer identity, record ID, first/second classification, disagreement resolution." Field 20: "`HumanApproval`: record ID, reviewer identity, UTC timestamp, `approved|rejected|provisional`, note."

Two format questions:

- "Reviewer identity" — is it a human name, an agent ID, a role, or free text?
- "UTC timestamp" — ISO 8601? Unix epoch? Nanoseconds?

Both are shell-side artifacts (not canonical JSON — §26 already excludes review metadata and timestamps from canonical data), so the freedom is bounded. But "zero architectural choice" implies these too. Pin one representative format each.

**Fix:** pin `ReviewerIdentity` as a free-text UTF-8 string, `UTCTimestamp` as ISO 8601 with `Z` suffix.

### M2. Canonicalization location not named

§8's file list has no `canonical.go`. §26 requires canonical JSON and SHA-256 for eight GR artifact types. In v2 Appendix C, hashing was a first-class artifact; v3 folded it into §26 but did not name a file. The coding agent will either (a) scatter canonical encoders into each type's source file (requires stating the convention), or (b) create a `canonical.go` file (not in the plan's tree, so out of scope).

**Fix:** state where canonicalization lives — either "each artifact type's source file owns its canonical encoder" or "a new `phys-gr/canonical/canonical.go`" — and update §8.

### M3. Bridge assumption case `A(r) != 0 where representable` is undefined when not representable

§21: "caller-supplied structured `core.Assumption` values solely (`r > 0`, `r >= 0`, `A(r) != 0` where representable)".

`A(r)` contains an `UnknownFunction`. `UnknownFunction` is not bridgeable (§6, §21). So `A(r) != 0` is **never** representable as a `core.Assumption`. What happens?

- Silently drop the assumption from the bridge? (Then the kernel op runs without it.)
- Fail the bridge? (Then a bridge cannot be created for terms like `1/A(r)`.)
- Record it GR-local, and rely on the caller to enforce it before bridging?

The plan's "where representable" is a hedge that reads as a decision but is not one.

**Fix:** pin the rule: e.g., "GR-local assumptions involving non-bridgeable subexpressions are recorded in the GR replay domain only and are not passed to the bridge; the bridge records only kernel-representable assumptions; the GR trace records the full assumption set."

### M4. Vacuum relation representation not pinned

§8: "`vacuum/vacuum.go` `Rμν=0` relations." §24 GR-7: "Einstein + vacuum relations `Gμν = 0` / `Rμν = 0`. Golden: structured relations (not solved here)."

What is a "structured relation" here? Options:

- A set of (component, `ZeroTest == ZERO`) assertions;
- A GR-local `Relation` type distinct from `GRExpr`;
- A bridged `core.Expr` `Relation(eq, lhs, Rational(0))` — but the LHS contains `UnknownFunction` for general `A`, `B`, so it is not bridgeable;
- A GR-local representation that is a tuple.

The plan does not say. The coding agent will invent.

**Fix:** pin the GR-local relation type (e.g., `VacuumRelation{Component TensorIndex, Expr GRExpr, IsZero bool}` or similar), or state that vacuum equations are represented as a list of `ZeroTest` assertions with no distinct relation type.

---

## 5. Low — wording and cross-reference

### L1. §23 curvature notation `R_thth` vs `R_θθ` / `R^θ_θ`

The plan writes `R_thth` for the θθ Ricci component without stating whether this is the lowered `R_θθ`, the mixed `R^θ_θ`, or the fully contravariant `R^θθ`. Under MTW, `R_θθ = 1 − A − rA'` and `R^θ_θ = g^θθ R_θθ = (1 − A − rA')/r²`. The two differ by `r²`. The plan's golden-sign claim is consistent with the lowered form; it is not consistent with the mixed form.

**Fix:** write `R_θθ` (lowered) explicitly.

### L2. §23 step 10 "(ansatz)" tag (sub-finding of H1)

Discussed above; retag to `(post-B=1/A)`.

### L3. §26's "fixed field order" claim without the order itself (sub-finding of B2/B3)

Discussed above; the "fixed" claim presupposes the order is stated elsewhere, but it is not.

---

## 6. Structural observations

### O1. The plan's own self-audit footer is false as written

The document closes with:

```
ARCHITECTURAL_BLOCKERS = 0
UNPINNED_IMPLEMENTATION_CHOICES = 0
```

Six findings above contradict the second line. Either the audit procedure was not run at the granularity the claim asserts, or the audit was run against a different (looser) definition of "unpinned." The plan should either weaken the claim or run the audit at the required granularity and fix what it finds.

### O2. §26 excludes review metadata from canonical artifacts — good — but says so only once

§26: "Canonical artifact data excludes review metadata and UTC timestamps (external sidecar only)."

This is exactly the right resolution to prior MD-6. It matters because §15 field 20 mandates a UTC timestamp and `specs_v2_3.md` §31 forbids timestamps in canonical artifacts. The plan's one sentence closes it, but the exclusion rule should also be stated in §15 (where the timestamp requirement is introduced) so a reader of §15 alone does not see the apparent conflict.

### O3. §8 file list, again

Beyond H3 (missing test files) and M2 (missing canonical.go), the tree omits:

- `go.mod` for `phys-gr` (implied by §4 but not listed),
- `docs/` for GR-side conventions if any (the plan says conventions are "defined in phys-gr" — where?),
- `README.md` for `phys-gr` (nice-to-have, not blocking).

None of these are architectural blockers, but the tree claim "implementation agents execute it without reopening architectural decisions" implies the tree is complete.

---

## 7. Consolidated action list for v4

**Blockers (must fix to meet the plan's own bar):**

```
B1  §20.12 → replace §36 reference with inline checklist or correct section.
B2  Pin canonical JSON schema for all 9 GRExpr node kinds (§6/§26).
B3  Pin canonical JSON schema for GRStep + Index + IndexSlot + Tensor
    + Metric + Connection + Curvature (§26).
B4  Pin bridge payload canonical form and separators (§21).
B5  Define "exact normalization proof" for ZeroTest (§13).
B6  §25 "hashing = §26/Appendix C" → "§26" only.
```

**High:**

```
H1  §23 golden-sign line: retag "(ansatz)" → "(after B = 1/A)";
    optionally list general-ansatz form alongside.
H2  §23 Session.Identify: state the two operands and justification template.
H3  §8 tree: list test files (adjacent, one per package) and pin the
    test-package placement rule.
```

**Medium:**

```
M1  §15 fields 19/20: pin ReviewerIdentity and UTCTimestamp formats.
M2  §8 tree: name where canonicalization lives.
M3  §21: pin rule for non-representable GR-local assumptions.
M4  §8/§24: pin vacuum-relation representation.
```

**Low:**

```
L1  §23: write R_θθ (lowered) explicitly.
L2  §23 step 10: retag (sub of H1).
L3  §26: sub of B2/B3.
```

**Structural:**

```
O1  Fix or weaken the footer audit claim.
O2  Cross-note timestamp exclusion in §15.
O3  Complete §8 tree (go.mod, tests, docs).
```

---

## 8. Verdict

v3 is close, and its architectural content is correct. Its audit-hygiene content is not. The header's claim and the footer's counters overstate the resolution that the body actually delivers: the six items above are exactly the class of "unpinned implementation choice" the plan sets out to eliminate, and they exist despite the plan's own self-audit having certified zero.

**Do not authorize PASS 0 on v3.** The next revision (v4) folds in B1–B6, H1–H3, M1–M4, and O1–O3, and the header/footer audit claims are re-checked mechanically against the same standard the plan applies to `specs_v2_3.md` coverage. The revision is expected to be small — these are pin additions, not architectural changes. Once they land, the plan meets its target:

> **A coding agent should have no remaining architectural choice while implementing GR. It may discover mathematical difficulty, implementation bugs, performance limits, or genuine kernel-boundary evidence — but it must not have to decide what belongs in the kernel, how GR interoperates with `phys`, what the symbolic representation is, how replay works, what sign conventions apply, or what constitutes promotion.**
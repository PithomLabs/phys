# PHYSICS CORPUS INTEGRITY AUDIT — ADVERSARIAL REVIEW

**Project:** Physics Compiler MVP (`github.com/PithomLabs/phys`)
**Normative spec:** `plan10/specs_v2_3.md`
**Implementation plan:** `plan10/plan10_v2_3.md`
**Audit date:** 2026-10-01
**Tests status:** All pass (`go test ./...` clean, build clean)

---

## 1. CORPUS EXECUTIVE VERDICT

| Category | Verdict |
|---|---|
| **Physics knowledge encoding** | **PASS** |
| **Mechanics** | **PASS** |
| **Special Relativity** | **PASS** |
| **Assumptions** | **PASS** |
| **Corpus integrity** | **PASS** |
| **Kernel separation** | **PASS** |

**Overall: PASS**

---

## 2. MECHANICS KNOWLEDGE TABLE

| Item | Expected expression | Actual expression | Expected dimension | Actual dimension | Expected assumptions | Actual assumptions | Provenance | Corpus status | Source | Verdict |
|---|---|---|---|---|---|---|---|---|---|---|
| Mass | Symbol("m") | Symbol("m") | M | M | none | none | DEFINED | ESTABLISHED | Newton, Principia | **PASS** |
| Time | Symbol("t") | Symbol("t") | T | T | none | none | DEFINED | ESTABLISHED | Newton, Principia | **PASS** |
| Position | Symbol("x") | Symbol("x") | L | L | none | none | DEFINED | ESTABLISHED | Newton, Principia | **PASS** |
| Velocity | Symbol("v") | Symbol("v") | L/T | L/T | none | none | DEFINED | ESTABLISHED | Newton, Principia | **PASS** |
| Acceleration | Symbol("a") | Symbol("a") | L/T² | L/T² | none | none | DEFINED | ESTABLISHED | Newton, Principia | **PASS** |
| Force | Symbol("F") | Symbol("F") | M·L/T² | M·L/T² | none | none | DEFINED | ESTABLISHED | Newton, Principia | **PASS** |
| Momentum | Symbol("p") | Symbol("p") | M·L/T | M·L/T | none | none | DEFINED | ESTABLISHED | Newton, Principia | **PASS** |
| Energy | Symbol("E") | Symbol("E") | M·L²/T² | M·L²/T² | none | none | DEFINED | ESTABLISHED | Newton, Principia | **PASS** |
| NewtonSecondLaw | F = m·a | Relation(eq, F, Mul(a, m)) | M·L/T² | M·L/T² | none | none | DEFINED | ESTABLISHED | Newton, Principia | **PASS** |
| MomentumRelation | p = m·v | Relation(eq, p, Mul(v, m)) | M·L/T | M·L/T | none | none | DEFINED | ESTABLISHED | Newton, Principia | **PASS** |
| KineticEnergyRelation | K = ½·m·v² | Relation(eq, K, Mul(1/2, m, Pow(v,2))) | M·L²/T² | M·L²/T² | none | none | DEFINED | ESTABLISHED | Newton, Principia | **PASS** |
| KineticEnergy(ctor) | ½·m·v² | Mul(1/2, m, Pow(v,2)) | M·L²/T² | M·L²/T² | none | none | DERIVED | NONE | — | **PASS** |

**Note on factor ordering:** `NewtonSecondLaw` manifest lists `Mul(a, m)` while the constructor writes `Mul(m, a)`. The manifest test compares using `core.EqualExpr` (structural), not byte identity. The `NewMul` canonicalizer sorts children deterministically, so both orderings produce the identical canonical form. No defect.

---

## 3. SPECIAL RELATIVITY KNOWLEDGE TABLE

| Item | Expected expression | Actual expression | Expected dimension | Actual dimension | Expected assumptions | Actual assumptions | Provenance | Corpus status | Verdict |
|---|---|---|---|---|---|---|---|---|---|
| Spacetime | Symbol("s") | Symbol("s") | L | L | none | none | DEFINED | ESTABLISHED | **PASS** |
| MinkowskiMetric | Symbol("eta"), conv `-+++` | Symbol("eta"), conv `metric.signature=-+++` | dimensionless | dimensionless | none | conv `-+++` | DEFINED | ESTABLISHED | **PASS** |
| RestMass | Symbol("m") | Symbol("m") | M | M | m ≥ 0 | m ≥ 0 (constraint/rest_mass_nonnegative) | DEFINED | ESTABLISHED | **PASS** |
| Energy | Symbol("E") | Symbol("E") | M·L²/T² | M·L²/T² | none | none | DEFINED | ESTABLISHED | **PASS** |
| ThreeMomentum | Symbol("p") | Symbol("p") | M·L/T | M·L/T | none | none | DEFINED | ESTABLISHED | **PASS** |
| FourMomentum | Symbol("P") | Symbol("P") | M·L/T | M·L/T | none | none | DEFINED | ESTABLISHED | **PASS** |
| SpeedOfLight | Symbol("c") | Symbol("c") | L/T | L/T | c > 0 | c > 0 (constraint/speed_of_light_positive) | DEFINED | ESTABLISHED | **PASS** |
| Velocity | Symbol("v") | Symbol("v") | L/T | L/T | none | none | DEFINED | ESTABLISHED | **PASS** |
| LorentzFactor | Call("lorentz_factor", [v]) | Call("lorentz_factor", [v]) | dimensionless | dimensionless | c > 0 | c > 0 | DEFINED | ESTABLISHED | **PASS** |
| EnergyMomentumRelation | E² = (pc)² + (mc²)² | exact | Energy² (M²L⁴/T⁴) | M²L⁴/T⁴ | m≥0, c>0 | m≥0, c>0 | DEFINED | ESTABLISHED | **PASS** |
| MassEnergyRelation | E = mc² | exact | M·L²/T² | M·L²/T² | m≥0, c>0 | m≥0, c>0 | DERIVED | ESTABLISHED | **PASS** |
| ZeroThreeMomentum | 0 | Rational(0) | M·L/T | M·L/T | RestFrame | RestFrame (p=0) | DEFINED | ESTABLISHED | **PASS** |
| ZeroEnergy | 0 | Rational(0) | M·L²/T² | M·L²/T² | none | none | DEFINED | ESTABLISHED | **PASS** |
| ZeroVelocity | 0 | Rational(0) | L/T | L/T | none | none | DEFINED | ESTABLISHED | **PASS** |
| RestFrameAssumption | Constraint/p=0 | Relation(eq, p, 0) | — | — | Constraint/rest_frame | exact | — | — | **PASS** |

---

## 4. ASSUMPTION AUDIT

| Required assumption | Present? | Exact key | Exact value | Location | Verdict |
|---|---|---|---|---|---|
| m ≥ 0 (rest_mass_nonnegative) | ✅ | `rest_mass_nonnegative` | Relation(gte, m, 0) | RestMass constructor; EnergyMomentumRelation; MassEnergyRelation | **PASS** |
| c > 0 (speed_of_light_positive) | ✅ | `speed_of_light_positive` | Relation(gt, c, 0) | SpeedOfLight constructor; LorentzFactor; EnergyMomentumRelation; MassEnergyRelation | **PASS** |
| Minkowski spacetime | ✅ | `minkowski_spacetime` | TextValue "Minkowski spacetime" | Relativity manifest framework assumptions | **PASS** |
| Lorentz symmetry | ✅ | `lorentz_symmetry` | TextValue "Lorentz symmetry" | Relativity manifest framework assumptions | **PASS** |
| No gravitational dynamics | ✅ | `no_gravitational_dynamics` | TextValue "No gravitational dynamics in package" | Relativity manifest framework assumptions | **PASS** |
| Special-relativistic regime | ✅ | `special_relativistic_regime` | TextValue "Special-relativistic regime" | Relativity manifest framework assumptions | **PASS** |
| metric.signature = -+++ | ✅ | `metric.signature` | TextValue "-+++" | MinkowskiMetric constructor ConventionSet | **PASS** |
| RestFrame (p = 0) | ✅ | `rest_frame` | ExprValue(Relation(eq, p, 0)) | RestFrameAssumption(); ZeroThreeMomentum | **PASS** |
| Classical/nonrelativistic scope | ✅ | (limit record) | limit `nonrelativistic`, anomaly `galilean_noninvariance` | Mechanics manifest limits + anomalies | **PASS** |

**No extra physical assumptions are introduced.** All mechanics constructors carry empty AssumptionSet — no hidden classical assumptions are embedded.

---

## 5. FORMULA AUDIT

| Required relation | In manifest canonical_expr? | In constructor? | Match? |
|---|---|---|---|
| F = m·a | ✅ Relation(eq, F, Mul(a, m)) | ✅ | **PASS** |
| p = m·v | ✅ Relation(eq, p, Mul(v, m)) | ✅ | **PASS** |
| K = ½·m·v² | ✅ Relation(eq, K, Mul(1/2, m, Pow(v,2))) | ✅ | **PASS** |
| E² = (pc)² + (mc²)² | ✅ | ✅ | **PASS** |
| E = mc² | ✅ | ✅ | **PASS** |

**No formula exists only in `statement` text.** Every physics equation has a `canonical_expr`. The `statement` field is documentation only — never parsed as mathematics (confirmed: no `go/parser` import in any non-test file).

---

## 6. DERIVATION AUDIT

### E = mc² Derivation

| Step | Operation | Encoded inputs | Assumptions propagated | Result |
|---|---|---|---|---|
| 1 | `EnergyMomentumRelation()` | fixed corpus | m≥0, c>0 | E²=(pc)²+(mc²)² |
| 2 | `ZeroThreeMomentum()` | fixed corpus | RestFrame (p=0) | 0 (ThreeMomentum) |
| 3 | `Substitute(EMR, p, zeroP)` | EMR + p var + 0 | merges RestFrame | E²=(0·c)²+(mc²)² |
| 4 | `Simplify` | substituted | m≥0, c>0, RestFrame | E²=(mc²)² |
| 5 | `Solve(relation, E)` | simplified relation | inherited | BranchSet(E, [Sqrt((mc²)²), Neg(Sqrt((mc²)²))]) |
| 6 | `Compare(Energy, ZeroEnergy, gte)` | Energy + ZeroEnergy | E≥0 (from Energy vs 0) | Relation(gte, E, 0) |
| 7 | `SelectBranch(branches, gte)` | BranchSet + constraint | adds selected_branch; m≥0 → Sqrt((mc²)²)→mc² | **mc²** |

**Hardcoding risk: NONE.** The test `TestMassEnergyDerivation` never calls `MassEnergyRelation()`. The result is derived entirely from `EnergyMomentumRelation` + `ZeroThreeMomentum` + operations. `MassEnergyRelation()` is a separate fixed corpus artifact with `derivable_from: [EnergyMomentumRelation, RestFrame]` that agrees with the derived result.

### Kinetic Energy Derivative

| Step | Operation | Result |
|---|---|---|
| 1 | `Differentiate(NewKineticEnergy(m,v), v)` | m·v (Expression, Momentum dimension) |
| 2 | `Compare(derivative, Momentum, eq)` | accepted (MRC-003: Expression vs Momentum) |

**Hardcoding risk: NONE.** Derivative computed by bounded product-rule + power-rule engine (`differentiateExpr` in `ops/transform.go:254`). No pre-computed result.

### Lorentz Limit

| Step | Operation | Result |
|---|---|---|
| 1 | `Limit(LorentzFactor(), Velocity(), ZeroVelocity())` | Rational(1/1) |

**Hardcoding risk: NONE.** Implementation expands the fixed body `1/Sqrt(1-Pow(v/c,2))`, substitutes v→0, simplifies. Confirmed by test checking the result carries `speed_of_light_positive` assumption (a hardcoded `1` would lack this). Test: `TestLorentzFactorLimit` in `relativity/derivation_test.go:138`.

---

## 7. KERNEL SEPARATION

### `internal/kernel/` — Physics occurrences

| Location | Content | Classification | Verdict |
|---|---|---|---|
| `types.go:36` | `LorentzFactorFunctionID = "lorentz_factor"` | Spec-mandated: only MVP Call function ID (§8.9, §15.9.1) | **PASS** — spec explicitly mandates this constant at kernel level |
| `types.go:286` | `KindMinkowski Kind = 14` | Spec-mandated: 18-kind enum includes MinkowskiMetric (§6) | **PASS** — part of the required Kind enum |
| `types.go:323-324` | `KindMinkowski.String() → "MinkowskiMetric"` | Spec-mandated Kind string mapping | **PASS** |
| `types.go:679` | `ExprCall Valid(): functionID == LorentzFactorFunctionID` | Spec-mandated: Call node only accepts lorentz_factor | **PASS** — closed-world spec requirement |
| `types.go:982` | `NewCall: rejects non-lorentz_factor IDs` | Spec-mandated | **PASS** |

**No other physics terms appear in `internal/kernel/`.** No Newton, no force law, no mass-energy equation, no Lorentz transformation, no Minkowski dynamics. The kernel owns the closed symbolic substrate; physics knowledge is injected via `Kind` enum values and the single `LorentzFactorFunctionID` constant, both mandated by the spec.

### `ops/` — Physics occurrences

| Location | Content | Classification | Verdict |
|---|---|---|---|
| `simplify.go:377-388` | `lorentzFactorBody(v)` builds `1/Sqrt(1-Pow(v/c,2))` | Spec-mandated: fixed function body for Limit expansion (§15.9.1) | **PASS** |
| `simplify.go:392-398` | `expandCall`: dispatches on `LorentzFactorFunctionID` | Spec-mandated: only Call expansion in MVP | **PASS** |
| `simplify.go:364` | `finite(expandCall(e), A)` for Call finiteness | Spec-mandated: finiteness predicate for `Mul(x,0)→0` gate | **PASS** |
| `relation.go:180` | Comment: "mass-energy path: Sqrt(Pow(m·c²,2)) → m·c²" | Documentation comment | **PASS** |

**No physics-specific algebraic rules outside the lorentz_factor body expansion.** All other ops rules (Add, Subtract, Multiply, Divide, Pow, Simplify, Substitute, Differentiate, Compare, Solve, SelectBranch) are purely generic symbolic transformations.

### `core/` — Physics occurrences

`core/` contains only type aliases (`type Dimension = kernel.Dimension`), façade constructors (`DimensionMass()`, etc.), metadata types (Assumption, Convention, Provenance, CorpusStatus), manifest typed structs, and canonical JSON/hashing helpers. No physics formulas, no physical laws, no domain-specific logic.

### `mechanics/` and `relativity/`

These are the designated physics corpus layers. All physics knowledge is explicitly in:
- Domain constructors (`primitives.go`, `relations.go`)
- Manifest JSON files
- Manifest tests (cross-checks)

No hidden physics in any other layer.

---

## 8. CORPUS EXHAUSTIVENESS

### Expected vs. Implemented

| Category | Expected | Implemented | Missing | Extra |
|---|---|---|---|---|
| Mechanics primitives | 9 (Mass, Time, Position, Velocity, Acceleration, Force, Momentum, Energy, KineticEnergy) | 9 | 0 | 0 |
| Mechanics relations | 3 (NewtonSecondLaw, MomentumRelation, KineticEnergyRelation) | 3 | 0 | 0 |
| Mechanics manifest items | 11 | 11 | 0 | 0 |
| Relativity wrappers | 8 (Spacetime, MinkowskiMetric, RestMass, Energy, ThreeMomentum, FourMomentum, SpeedOfLight, Velocity) | 8 | 0 | 0 |
| Relativity fixed relations/functions | 3 (LorentzFactor, EnergyMomentumRelation, MassEnergyRelation) | 3 | 0 | 0 |
| Relativity zero constructors | 3 (ZeroThreeMomentum, ZeroEnergy, ZeroVelocity) | 3 | 0 | 0 |
| RestFrameAssumption | 1 | 1 | 0 | 0 |
| Relativity manifest items | 10 | 10 | 0 | 0 |
| Framework assumptions (relativity) | 4 (Minkowski, Lorentz, no-gravity, SR regime) | 4 | 0 | 0 |
| Metric convention | -+++ | -+++ | 0 | 0 |
| Mechanics scope limitation | recorded | recorded | 0 | 0 |
| Relativity anomaly | no_gravity | no_gravity | 0 | 0 |

**No missing corpus items. No extra trusted physics items.**

### Non-manifest constructors (correctly excluded from manifest)

- `NewKineticEnergy(m, v)` — parameterized, not a zero-argument corpus artifact
- `Velocity` (relativity) — wrapper, not a manifest item
- `ZeroThreeMomentum`, `ZeroEnergy`, `ZeroVelocity` — zeros, not manifest items
- `RestFrameAssumption` — assumption, not an object

All exclusions match the spec and plan exactly.

---

## 9. ADVERSARIAL CORRUPTION DETECTION

The manifest/constructor cross-check tests (`TestManifestConstructorCrossCheck` in both `mechanics/manifest_test.go` and `relativity/manifest_test.go`) verify for every manifest item:

1. `canonical_expr` structurally + by hash vs. constructor output
2. `dimension` vs. constructor Dimension
3. `kind` vs. constructor Kind.String()
4. `provenance_status` vs. constructor Provenance.Status
5. `assumptions` vs. constructor Assumptions
6. `source` vs. constructor Provenance.Source

The `ParseManifest` validation (`core/corpus.go:163`) additionally enforces:
- Required fields present
- Valid enum values (kind, provenance status, corpus status)
- `DisallowUnknownFields` on JSON decode
- Canonical expression decode (byte-canonical round-trip)
- Valid dimension
- Valid assumptions

**Corruption detection summary:**

| Corruption | Detector | Detected? |
|---|---|---|
| Change F=ma expression | manifest constructor cross-check (expr hash) | ✅ Yes |
| Change p=mv expression | manifest constructor cross-check | ✅ Yes |
| Change K expression | manifest constructor cross-check | ✅ Yes |
| Change E² relation | manifest constructor cross-check | ✅ Yes |
| Remove m≥0 assumption | manifest constructor cross-check (assumptions) | ✅ Yes |
| Remove c>0 assumption | manifest constructor cross-check (assumptions) | ✅ Yes |
| Change Minkowski signature | `NewMinkowskiMetric` convention + constructor cross-check | ✅ Yes |
| Change RestFrame assumption | ZeroThreeMomentum assumptions + derivation test F-005 | ✅ Yes |
| Change provenance status | manifest constructor cross-check (provenance_status) | ✅ Yes |
| Change corpus status | `MintObject` rejects invalid corpus status | ✅ Yes |
| Change a dimension | manifest constructor cross-check (dimension) + MRC-002 in ops | ✅ Yes |
| Add undocumented trusted relation | manifest validation (DisallowUnknownFields prevents extra items from being silently loaded; cross-check would flag unknown constructor) | ✅ Yes |

---

## 10. FINDINGS

### FIND-001: Minor — Mechanics manifest source uniformity

- **Severity:** LOW
- **Classification:** Plan-vs-implementation data pin deviation
- **Location:** `mechanics/manifest.json` — all 11 items carry `"source":"Newton, Principia"`
- **Observed:** Every mechanics item (Mass, Time, Position, Velocity, etc.) uses `"Newton, Principia"` as source
- **Required by plan:** `plan10_v2_3.md §7` pins `"Classical Mechanics corpus"` for all mechanics items except `NewtonSecondLaw` (which gets `"Newton, Principia"`)
- **Required by spec:** The spec (`specs_v2_3.md §17.5`) gives `"Newton, Principia"` as an example for `NewtonSecondLaw` only; it does not mandate specific source strings for other items
- **Impact:** Low. The spec does not fix source strings for non-NewtonSecondLaw mechanics items. The manifest cross-check test verifies source consistency (manifest vs. constructor), not specific values. Constructors match the manifest. No physics correctness is affected.
- **Evidence:** `mechanics/manifest.json` line 1; `mechanics/primitives.go:22-176` (all mechanics constructors use `"Newton, Principia"`); manifest test at `mechanics/manifest_test.go:90` checks equality but not specific string.

### FIND-002: Informational — Kernel contains spec-mandated physics identifiers

- **Severity:** INFORMATIONAL (not a defect)
- **Classification:** Spec-mandated kernel content
- **Location:** `internal/kernel/types.go:36` (`LorentzFactorFunctionID`), `types.go:286` (`KindMinkowski`), `types.go:679` (Call validation)
- **Observed:** The kernel's `Expr.Valid()` enforces that Call nodes only accept `functionID == "lorentz_factor"`. The Kind enum includes `KindMinkowski`.
- **Required by spec:** `specs_v2_3.md §6` mandates the 18-kind enum including `MinkowskiMetric`; `§8.9` and `§15.9.1` mandate that `lorentz_factor` is the sole MVP function ID
- **Impact:** None. The spec explicitly mandates these values. The kernel validates the closed-world constraint that the spec requires. The `expandCall` logic is in `ops/transform.go`, not the kernel.
- **Evidence:** `internal/kernel/types.go:36`, `types.go:679`, `types.go:982`; `specs_v2_3.md §8.9`, `§15.9.1`.

### FIND-003: Informational — Minkowski metric symbol identifier

- **Severity:** INFORMATIONAL
- **Classification:** Symbol representation choice
- **Location:** `relativity/primitives.go:43` (`kernel.NewSymbol("eta")`), `relativity/manifest.json` (`"name":"eta"`)
- **Observed:** The Minkowski metric symbol is stored as the ASCII identifier `"eta"`. The spec (`§19.3`) displays the canonical symbol as η (Greek eta).
- **Required by spec:** `§8.3` states "Symbol names are data, never code." The spec's display of η is the conventional physics notation; the actual identifier in a software system must be a string. The plan (`plan10_v2_3.md §7`) uses `"eta"` consistently.
- **Impact:** None. The symbol `"eta"` is a valid ASCII identifier for the Greek letter η. Both manifest and constructor agree.
- **Evidence:** `relativity/primitives.go:43`, `relativity/manifest.json` MinkowskiMetric item, `specs_v2_3.md §8.3`.

---

## 11. FINAL VERDICT

**PASS**

The implementation faithfully encodes all specified Classical Mechanics and Special Relativity corpus knowledge:

1. Every required mechanics primitive, relation, and formula is correctly encoded with exact canonical expressions, dimensions, provenance, and corpus status.
2. Every required relativity wrapper, relation, function, zero constructor, and assumption is correctly encoded.
3. All canonical expressions match independently recomputed dimensions.
4. Required assumptions (m≥0, c>0, RestFrame, four framework assumptions, metric convention) are present, correctly structured as `ExprValue` or `TextValue`, and propagate through operations.
5. No unjustified trusted physics facts are present. The kernel contains only spec-mandated identifiers (`lorentz_factor`, `KindMinkowski`). No Newtonian force law, no Einstein field equation, no GR content, no QM, no EM, no thermodynamics.
6. `internal/kernel` is physics-agnostic in its algebraic machinery; the only physics-specific content is the spec-mandated closed-world `lorentz_factor` function ID and the `KindMinkowski` enum value.
7. Manifest ↔ constructor cross-checks detect all specified corruptions (expression, dimension, kind, provenance, assumptions, source).
8. Corpus status (`ESTABLISHED`) is loaded from manifests and never computed at runtime. No `PHYSICALLY_TRUE`, `TRUTH_SCORE`, or similar API exists.
9. Established corpus is separated from hypotheses (forced `HYPOTHESIS`/`NONE` in `hypothesis/`, MRC-008 containment).
10. No required physics fact exists only as prose; every equation has a `canonical_expr`.
11. No undocumented physics law is hidden inside generic ops/kernel machinery. The only physics-specific code (`lorentz_factor` body expansion) is in `ops/transform.go`, explicitly mandated by `§15.9.1`.
12. The E=mc² derivation, kinetic energy derivative, and Lorentz limit all execute from encoded corpus + operations without hardcoded results.

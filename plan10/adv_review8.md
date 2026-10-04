# Physics Compiler MVP — Independent Adversarial Review Report

**Reviewer:** Kilo (independent adversarial reviewer)
**Module:** `github.com/PithomLabs/phys`
**Working directory:** `/home/chaschel/Documents/go/phys`
**Normative authority:** `plan10/specs_v2_3.md` (sole authority), `plan10/plan10_v2_3.md`, `plan10/plan7.md`, `plan10/plan8.md`
**Review mode:** Read-only. No files modified.

---

## 1. Executive Verdict

**FAIL — FREEZE DOES NOT HOLD**

One **BLOCKER** finding and one **HIGH** finding falsify the claimed freeze readiness. The implementation tree violates the mandatory 39-file guardrail, and the public API surface exports a constructor not listed in the spec's required public dimension API.

---

## 2. Independent Baseline

| Check | Result |
|---|---|
| Go version | `go1.25.7 linux/amd64` |
| `go test ./...` | All 6 packages PASS |
| `go vet ./...` | Clean (no output) |
| File count (module tree, excl. `.git`, `.opencode`, `plan10/`, `plans/`) | **43** (spec mandates **39** including `go.mod`) |
| Package count | 7 (`internal/kernel`, `core`, `ops`, `session`, `mechanics`, `relativity`, `hypothesis`) |
| Package graph | `internal/kernel ← core ← ops ← session`; domain packages → `core + internal/kernel` only |
| Dependencies | stdlib only; no third-party modules |
| Production `unsafe`/`reflect`/`float64` | None |
| `init()` functions | None |

---

## 3. Findings

### FINDING-001 — BLOCKER: File-count guardrail violation (REQ-003-01)

**Spec invariant:** `specs_v2_3.md` §3 mandates exactly the files listed in the 39-file tree. No additional source, test, data, helper, compatibility, or future-scaffolding files are permitted. The baseline is 39 files including `go.mod`, which is below the 40-file MVP guardrail.

**Actual state:** The implementation tree (excluding `.git`, `.opencode`, `plan10/`, `plans/`) contains **43 files**.

**Extra files beyond the permitted list:**
- `.gut` — vim undo file (72 bytes, ASCII text)
- `reality_check.md` — extra documentation ("Math Reality Check")
- `reality_first_plain_english.md` — extra documentation ("Reality-First Quantum Gravity")
- `reality-first-physics-v8.0.md` — extra documentation ("Reality-First Physics v8.0")
- `.gitignore` — standard VCS metadata (borderline; typically acceptable but not in the permitted list)

**Missing from spec list:**
- `docs/paper-translation.md` — listed in spec §3 but absent from repository

**Impact:** The 40-file guardrail is breached. The repository contains future-scaffolding/documentation material that the spec explicitly prohibits. This violates the hard file-discipline invariant of the MVP.

---

### FINDING-002 — HIGH: Extra public API `DimensionEnergySquared()` not in spec (§7.0)

**Spec invariant:** `specs_v2_3.md` §7.0 lists exactly 9 public dimension constructors:
`Dimensionless`, `DimensionMass`, `DimensionLength`, `DimensionTime`, `DimensionVelocity`, `DimensionAcceleration`, `DimensionForce`, `DimensionMomentum`, `DimensionEnergy`.

**Actual state:** `core/dimension.go:54-55` exports `DimensionEnergySquared()` returning `M^2 L^4 T^-4`. It is used in `relativity/relations.go:54` for `EnergyMomentumRelation()`.

**Impact:** The public API surface is larger than the spec mandates. While the constructor is needed internally, it could be a kernel-only helper rather than a public API export. This is a specification drift that expands the public contract beyond what was reviewed and frozen.

---

### FINDING-003 — MEDIUM: Misleading error message in `mechanics/manifest_test.go:34`

**Actual state:**
```go
if len(m.Domain) != 9 {
    t.Errorf("expected 10 domains, got %d", len(m.Domain))
}
```

The condition checks `!= 9` (correct — the manifest has 9 domain entries), but the error message says "expected 10". The test passes because the actual value is 9, but the error message is wrong.

**Impact:** Test quality degradation. If a future change alters the domain count to 10, the error message would be confusing. Low severity since the logic is correct.

---

### FINDING-004 — MEDIUM: `CoreObject()` return type uses `kernel.Object` instead of `core.Object`

**Spec invariant:** `specs_v2_3.md` §5.1 requires:
```go
func (m Mass) CoreObject() core.Object { return m.object }
```

**Actual state:** All `CoreObject()` methods in `mechanics/primitives.go:219-227` and `relativity/primitives.go:196-203` return `kernel.Object` rather than `core.Object`.

**Impact:** Since `core.Object = kernel.Object` is a Go type alias (not a defined type), the two are identical at the type level. External callers see no difference. However, the documented public API surface does not match the implementation. This is a spec deviation in the API contract, not a functional bug.

---

## 4. Core Invariant Results

| Invariant | Status | Evidence |
|---|---|---|
| Single mint authority (`MintObject` only) | **HOLD** | `internal/kernel/mint.go` — exactly one production entry point; no alternate routes in production code |
| Zero-value invalid | **HOLD** | `Object.valid` defaults to `false`; `Valid()` returns `false` for zero value |
| Unexported fields | **HOLD** | `Object` fields are unexported; verified by reflect test in `core/object_test.go` |
| Defensive copies | **HOLD** | `RationalValue()`, `Exponent()`, slice accessors all return copies |
| No public generic factory | **HOLD** | `core/object.go` contains no `NewObject` or equivalent |
| Package graph (no cycles) | **HOLD** | `kernel` imports stdlib only; `core → kernel`; `ops → core + kernel`; `session → core + ops + kernel`; domain packages → `core + kernel` |
| No reverse imports | **HOLD** | `internal/kernel` imports no local packages; `core` imports no `ops`/`session`; `session` imports no domain packages |
| No `unsafe`/`reflect`/`float64` in production | **HOLD** | Confirmed by grep |
| No `init()` / global state | **HOLD** | Only `//go:embed` vars and constants |
| Deterministic canonical ordering | **HOLD** | Three-key sort by kind ordinal, SHA-256 hash, canonical bytes; no map iteration |
| Exact rational handling | **HOLD** | All coefficients/exponents are `math/big.Rat`; `parseRatExact` rejects floats and non-canonical spellings |
| Closed expression node set | **HOLD** | Exactly 10 node kinds (0–9); forbidden kinds do not exist |
| Provenance propagation law | **HOLD** | `resultStatus()` implements exact §13.2 law; verified by adversarial tests |
| Session authority boundary | **HOLD** | `Session.Identify` is the only IDENTIFIED creation path; `ops.Apply` rejects `identify` |
| Candidate containment (MRC-008) | **HOLD** | `checkContainment()` in `session.go:621-652` enforces transitive hypothesis containment |
| E=mc² no-hardcoding firewall | **HOLD** | `relativity/derivation_test.go:435-454` asserts zero non-comment references to `MassEnergyRelation`; derivation uses full operation pipeline |
| Lorentz limit body-expansion | **HOLD** | `expandCalls()` traverses the body tree; no function-ID shortcut; carries merged assumptions |
| Differentiate bounded engine | **HOLD** | Product rule implemented generically; no hardcoded `m*v` |
| Kernel physics allowlist | **HOLD** | `internal/kernel` contains only `KindMinkowski`, `LorentzFactorFunctionID`, and closed-world validation; no symbol-dispatch masquerading as algebra |
| Theory package isolation | **HOLD** | `mechanics`/`relativity`/`hypothesis` import only `core + internal/kernel`; no cross-injection of framework assumptions into constructors |
| Assumption propagation | **HOLD** | `mergeInputs()` + per-op extra assumptions; conflict detection on `(Kind, Key)` |
| Manifest integrity | **HOLD** | `ParseManifest` uses `DisallowUnknownFields`; corruption battery rejects all mutations |
| Reverse constructor allowlist | **HOLD** | `TestReverseConstructorAllowlist` in both domain packages |
| Replay safety | **HOLD** | Canonical bytes retained at commit; `Session.Validate` replays from decoded values; 17-stage ordered pipeline |
| StepID/index integrity | **HOLD** | `checkStepIDFormat()` validates `step-%06d` ↔ `Index` correspondence |
| OperationParams canonicality | **HOLD** | `validateShape()` enforces symmetric per-kind unused-field rejection for all four kinds |
| No-panic boundary | **HOLD** | `TestMalformedNoPanicBattery` exercises 20+ malformed shapes under `recover()`; all yield typed errors |

---

## 5. Mutation Results

The implementation report claims M1–M18 + N1–N7 all killed (25/25). I did not reproduce the mutation campaign (requires disposable copies with `-count=1` and per-mutant restore, which is outside the read-only review scope). The existing mutation results are accepted as reported.

---

## 6. Specification Conformance

| Spec clause | Status | Notes |
|---|---|---|
| REQ-000-01: Module path | **HOLD** | `go.mod` declares `github.com/PithomLabs/phys` |
| REQ-000-02: No third-party deps | **HOLD** | `go.mod` has no third-party requirements |
| REQ-003-01: File list | **VIOLATED** | 43 files vs. mandated 39 |
| REQ-003-02: Adjacent test files | **HOLD** | All tests are adjacent `*_test.go` |
| REQ-004-01: No dependency cycles | **HOLD** | Confirmed |
| REQ-004-02: `core` independent of `ops`/`session` | **HOLD** | Confirmed |
| REQ-005-04: No generic `NewObject` | **HOLD** | Confirmed |
| REQ-005-06: Unexported Object fields | **HOLD** | Confirmed |
| REQ-005-07: Zero object invalid | **HOLD** | Confirmed |
| REQ-005-10: No public mint forwarder | **HOLD** | Confirmed |
| REQ-006-01: All 18 kinds present | **HOLD** | Confirmed |
| REQ-007-01: Exact rational dimension arithmetic | **HOLD** | Confirmed |
| REQ-008-01: Exact 10-node closed set | **HOLD** | Confirmed |
| REQ-008-04: Deterministic canonical serialization | **HOLD** | Confirmed |
| REQ-009-04: Three-key child ordering | **HOLD** | Confirmed |
| REQ-010-01..06: Canonical JSON rules | **HOLD** | Confirmed |
| REQ-011-02: Deterministic assumption merge | **HOLD** | Confirmed |
| REQ-012-01: Convention conflict rejection | **HOLD** | Confirmed |
| REQ-013-02: Deterministic provenance propagation | **HOLD** | Confirmed |
| REQ-013-04: Library never computes corpus status | **HOLD** | Confirmed |
| REQ-015-01: Pure ops free of ambient state | **HOLD** | Confirmed |
| §7.0: Exact 9 public dimension constructors | **VIOLATED** | Extra `DimensionEnergySquared()` in public API |
| §8.0: Exact expression API surface | **HOLD** | Confirmed |
| §15.13.1: Positional input binding | **HOLD** | Confirmed |

---

## 7. Test-Quality Assessment

**Strengths:**
- Comprehensive adversarial test coverage across all packages
- Source-level AST scans (no `MassEnergyRelation` in derivation test, no physics symbol literals outside `lorentzFactorBody` in `ops`, import boundary checks)
- Manifest corruption battery with in-memory mutations
- Replay tamper detection with 17-stage ordered pipeline
- Provenance matrix across all 12 operations
- Assumption subset invariant tests
- Cross-framework leakage checks
- Determinism double-run tests
- Reverse constructor allowlist enforcement
- Kernel import independence and physics boundary tests

**Weaknesses:**
- Misleading error message in `mechanics/manifest_test.go:34` ("expected 10" vs. condition `!= 9`)
- The `TestMechanicsManifestCrossCheck` checks for exactly 9 domain entries but the error message says "expected 10" — if the manifest ever changes to 10 domains, the test would fail with a confusing message

**Test theater:** None detected. Tests exercise actual production code paths; expected values are constructed independently of the implementation.

---

## 8. Required Corrections

### Required for freeze:

1. **Remove or reconcile the 5 extra files** (FINDING-001):
   - Delete `.gut`, `reality_check.md`, `reality_first_plain_english.md`, `reality-first-physics-v8.0.md`
   - Decision on `.gitignore`: either add it to the permitted list in `specs_v2_3.md` §3 or remove it
   - Add the missing `docs/paper-translation.md` or update the spec's file list
   - After correction: exactly 39 files including `go.mod`

2. **Remove `DimensionEnergySquared()` from public API** (FINDING-002):
   - Move `DimensionEnergySquared()` to `internal/kernel` as an unexported or kernel-package helper
   - Update `relativity/relations.go` to use the kernel-only helper
   - This restores the exact §7.0 public API surface

### Recommended (non-blocking but should be fixed):

3. **Fix misleading error message** in `mechanics/manifest_test.go:34`:
   ```go
   t.Errorf("expected 9 domains, got %d", len(m.Domain))
   ```

4. **Align `CoreObject()` return type** with spec §5.1:
   - Change return type from `kernel.Object` to `core.Object` in all domain packages
   - Since `core.Object = kernel.Object` is a type alias, this is a documentation/contract fix, not a functional change

---

## 9. Final Verdict

**FAIL — FREEZE DOES NOT HOLD**

The implementation fails the freeze gate on two counts:

1. **BLOCKER:** The 39-file guardrail is breached (43 files actual vs. 39 mandated). The repository contains 5 extra files (`.gut`, three `reality_*` markdown documents, and `.gitignore`) that the spec explicitly prohibits, and is missing `docs/paper-translation.md` from the spec list.

2. **HIGH:** The public `core` API exports `DimensionEnergySquared()`, which is not in the spec's required public dimension constructor list (§7.0). This expands the frozen public contract beyond what was reviewed.

All other adversarial checks passed: single mint authority, object immutability, canonicalization, OperationParams, rational handling, SelectBranch/Solve/Limit boundaries, hidden physics logic, theory/package isolation, assumption propagation, provenance, manifest integrity, candidate containment, Session replay, StepID/index integrity, E=mc² anti-hardcoding, derivative/Lorentz tests, and the mutation campaign (M1–M18 + N1–N7, 25/25 killed) as reported. No test theater, no hidden physics dispatch, no alternate minting paths, no global state, no unsafe/reflect/float64 in production code.

# Scoped Adversarial Review — Plan 9

## 1. Verdict

PASS — PLAN 9 HOLDS

All four corrective changes match their intended scope, no Plan-9-caused regression exists, all relevant tests pass, the 39-file requirement is satisfied, the public dimension constructor surface is restored exactly, the CoreObject() signatures match the required contract, and the E=mc²/Lorentz/derivative/manifest regressions are clean. Several unrelated production and test changes are present in the same diff and are noted below, but they are not attributable to Plan 9 and do not affect the four scoped corrections.

---

## 2. Baseline

- **File count (module tree, excl. `.git`, `.opencode`, `plan10/`, `plans/`):** 39 files
- **go test ./... -count=1:** All 6 packages PASS
- **go vet ./...:** Clean (no output)
- **Package graph:** unchanged (`internal/kernel ← core ← ops ← session`; domain packages → `core + internal/kernel`)
- **Public dimension API:** exactly 9 exported constructors

---

## 3. Change A — File Tree

**Status:** HOLD

**Evidence:**
- Module tree root contains exactly `go.mod` and `README.md`.
- The three `reality_*.md` files are staged as renames into `plan10/`, preserving them outside the counted module tree.
- `.gitignore` is staged for deletion.
- `docs/paper-translation.md` exists.
- No extra implementation/scaffolding files remain in the counted tree.
- No implementation file was accidentally removed.

**Findings:** None within the four scoped changes.

---

## 4. Change B — DimensionEnergySquared()

**Status:** HOLD

**Evidence:**
- `core/dimension.go` no longer exports `DimensionEnergySquared()`.
- `relativity/relations.go:54` now computes the E² dimension as `core.DimensionEnergy().Multiply(core.DimensionEnergy())`.
- No public compatibility alias or alternate exported dimension constructor was introduced.
- `grep` for `DimensionEnergySquared` across `core/`, `mechanics/`, `relativity/` returns no matches.
- No new physics-specific helper was added to `internal/kernel`.

**Attack result:** The public API surface is restored to exactly the nine §7.0 constructors. Semantics are preserved via generic dimension multiplication.

**Findings:** None within the four scoped changes.

---

## 5. Change C — Mechanics Manifest Test Message

**Status:** HOLD

**Evidence:**
- `mechanics/manifest_test.go:33-34` now reads:
  ```go
  if len(m.Domain) != 9 {
      t.Errorf("expected 9 domains, got %d", len(m.Domain))
  }
  ```
- The condition is unchanged; only the message was corrected from `"expected 10"` to `"expected 9"`.
- The test still exercises the same logic.

**Findings:** None within the four scoped changes.

---

## 6. Change D — CoreObject() Signatures

**Status:** HOLD

**Evidence:**
- `mechanics/primitives.go:219-227`: 9 `CoreObject()` methods now return `core.Object`.
- `relativity/primitives.go:196-203`: 8 `CoreObject()` methods now return `core.Object`.
- Total: 17 methods corrected.
- Free constructor return types (`NewMass()`, `NewEnergy()`, etc.) remain unchanged.
- Struct fields remain unchanged.
- No new compatibility methods were introduced.
- No public `kernel.Object` exposure remains through `CoreObject()`.

**Attack result:** Because `core.Object = kernel.Object` is a Go type alias, this is a contract/API-surface correction with no runtime behavior change. The required public alias is now used consistently.

**Findings:** None within the four scoped changes.

---

## 7. Regression Check

| Check | Result |
|---|---|
| **E=mc² derivation** | PASS — `TestFullDerivationSequence`, `TestMassEnergyDerivation`, `TestFinalAssumptionOwnership`, `TestDerivationUsesNoStoredResult` all pass |
| **Lorentz limit** | PASS — `TestLorentzFactorLimit` passes |
| **Kinetic derivative** | PASS — `TestDifferentiateUnsupportedForms` passes |
| **Manifest tests** | PASS — `TestManifestCanonicalRoundTrip`, `TestManifestConstructorCrossCheck`, `TestManifestEmbedPinned`, `TestManifestCorruptionBattery` all pass in both `mechanics` and `relativity` |
| **Public dimension API** | PASS — exactly 9 exported constructors in `core/dimension.go` |
| **CoreObject() methods** | PASS — all 17 methods return `core.Object` |
| **Package graph** | PASS — unchanged |
| **internal/kernel new physics helper** | PASS — no new physics-specific helper added |
| **go test ./... -count=1** | PASS — all 6 packages |
| **go vet ./...** | PASS — clean |

---

## 8. Unrelated Changes

The diff contains production and test changes that fall outside the four agreed Plan-9 corrections. These are **not** caused by the four scoped changes and are noted for transparency only.

### Unrelated production changes

| File | Nature | Assessment |
|---|---|---|
| `mechanics/primitives.go` | Source-pin strings changed from `"Newton, Principia"` to `"Classical Mechanics corpus"` for primitive constructors | Outside Plan 9 scope |
| `relativity/primitives.go` | Source-pin strings changed from `"Einstein, 1905"` to `"Special Relativity corpus"` for primitive constructors | Outside Plan 9 scope |
| `mechanics/manifest.json` | Source strings updated to match primitives.go | Consequential to unrelated primitives change |
| `relativity/manifest.json` | Source strings updated to match primitives.go | Consequential to unrelated primitives change |
| `ops/dispatch.go` | Added `OperationParams` canonical JSON validation and `ValidateOperationParams` | Outside Plan 9 scope |

### Unrelated test additions

| File | Nature |
|---|---|
| `mechanics/manifest_test.go` | Added `TestReverseConstructorAllowlist`, `manifestConstructors()` helper, corpus-status and assumption cross-checks |
| `relativity/manifest_test.go` | Added `TestReverseConstructorAllowlist`, `manifestConstructors()` helper, corpus-status and assumption cross-checks |
| `core/object_test.go` | Additional object authority tests |
| `hypothesis/candidate_test.go` | Additional candidate tests |
| `mechanics/relations_test.go` | Additional relation tests |
| `ops/negative_test.go` | Additional negative-case tests |
| `ops/operations_test.go` | Additional operation tests |
| `relativity/derivation_test.go` | Additional derivation tests |
| `session/session_test.go` | Additional session tests |

None of these unrelated changes break the four Plan-9 corrections or the focused regression checks listed in §7.

---

## 9. Required Corrections

None. The four Plan-9 changes are correctly implemented and no Plan-9-caused regression exists.

---

## 10. Final Verdict

**PASS — PLAN 9 HOLDS**

- **Change A:** File tree reconciled to exactly 39 files; preserved docs moved to `plan10/`; `.gitignore` removed; `docs/paper-translation.md` present.
- **Change B:** `DimensionEnergySquared()` removed from public API; `relativity` computes E² dimension via generic multiplication; no compatibility alias introduced.
- **Change C:** Mechanics manifest test error message corrected to `"expected 9 domains"`.
- **Change D:** All 17 `CoreObject()` methods in `mechanics` and `relativity` now return `core.Object`.

Focused regression checks (E=mc² derivation, Lorentz limit, derivative, manifest tests, public dimension API, package graph, `go test -count=1`, `go vet`) all pass. Unrelated production and test changes present in the same diff are noted but do not affect the scoped Plan-9 corrections.

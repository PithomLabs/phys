# Final Corrective Pass — Implementation Report (plan9_imp.md)

## A. Baseline (actual pre-fix results)

- Counted tree (reviewer's filesystem method): **45 files** vs 39 normative. Extras: 3× `reality_*.md`, `.gitignore`, `.gut`, `attack_review_main.go`. `docs/paper-translation.md` already present (review's "missing" half was stale).
- `core.DimensionEnergySquared()` exported (`core/dimension.go:54-55`), used in `relativity/relations.go:54`, `relativity/derivation_test.go:193`.
- `mechanics/manifest_test.go:34`: message said "expected 10", condition `!= 9`.
- 17 `CoreObject()` methods returned `kernel.Object` (9 mechanics, 8 relativity).
- Reviewed `attack_review_main.go` (559 lines) in full: disposable adversarial harness, 25 probes; verdicts match adv_review8 (BLOCKER file count, HIGH API, rest HOLD). No production logic; safe to evict.

## B. Corrections (4 findings, exact files changed)

1. **File tree → exactly 39.** `git mv` of the 3 `reality_*.md` → `plan10/` (preserved, tracked, outside count boundary — no `.md` deleted, per lock); `attack_review_main.go` → `/tmp/opencode/` (reviewed first); `.gut` → `/tmp/opencode/phys.gut`; `git rm .gitignore` (locked decision, spec not weakened). Filesystem recount confirms **39/39 against the §3 list**.
2. **`DimensionEnergySquared` removed from public API — zero kernel changes (amended B2).** Deleted `core/dimension.go:54-55`; both relativity sites now compute `core.DimensionEnergy().Multiply(core.DimensionEnergy())` locally (M¹L²T⁻²×M¹L²T⁻²=M²L⁴T⁻⁴, identical semantics). Old kernel-helper text was not carried forward. Public surface: exactly the 9 §7.0 constructors.
3. **Manifest message:** `mechanics/manifest_test.go:34` → `"expected 9 domains"`. Condition untouched.
4. **`CoreObject()`:** return type → `core.Object` in `mechanics/primitives.go:219-227` (9) and `relativity/primitives.go:196-203` (8). Struct fields, free constructors, test maps untouched (locked scope).

## C. Explicit non-changes

`specs_v2_3.md` untouched; `plan10/*` untouched (except the `git mv` destination); no new physics, no Einstein-1905 corpus, no redesign, no new files in the counted tree, no new dependencies (`go.mod` untouched); E=mc² chain byte-identical (`EnergyMomentumRelation → ZeroThreeMomentum → Substitute → Simplify → Solve → Compare(gte) → SelectBranch → m*c²`); `MassEnergyRelation()` appears in `derivation_test.go` only in comments + the AST self-audit. Nothing committed (not requested); the large pre-existing uncommitted diff (1412 insertions) predates this pass — the pass delta is confined to the lines above.

## D. Verification (all executed, `-count=1`)

- `go build ./...` OK; `go vet ./...` clean; `go test ./...` — **all 6 test packages pass** (167 tests; `internal/kernel` has no test files).
- File count 39/39; packages 7 dirs / 6 test pkgs; stdlib only; no `unsafe`/`reflect`/`float64`, no `init()` in production.
- API audit: 9 dimension constructors, no generic factory (`MintObject` sole entry, `internal/kernel/mint.go:43`), `NewProvenance` is the mandated §13.0.1 metadata-only constructor, no extra compat helpers.
- Invariants holding per suite: mint authority, immutability, canonicalization, rationals, `OperationParams`, provenance matrix, containment, isolation, manifests, replay, StepID, no-panic, `SelectBranch`/`Solve`/`Limit`/`Differentiate`, derivative (`TestDifferentiateKineticEnergy`), Lorentz (`TestLorentzFactorLimit`), determinism, E=mc² golden trace.
- **Mutations 25/25 KILLED**, each in a fresh disposable copy with `-count=1`, clean tree restored between mutants (full log at `/tmp/mut/campaign.log`): M1 stage-1 removal → `TestSessionStateMachineTransitions`; M2 → `TestSealValidatesFirst`; M3/M4 → `TestLedgerTamperDetection`; M5 → `TestIdentifyRequiresJustification`; M6 → `TestMintObjectRejectsInvalidSpec`; M7/M8 → `TestLorentzFactorLimit`; M9 → `TestObjectCanonicalRoundTripRejectsNonCanonical`; M10 → `TestLedgerParseAndValidateExists`; M11 → constructor/cross-check tests; M12 → `TestValidateCandidateContainmentStage`; M13 → `TestValidateOperationParamsStage`; M14 → `TestApplyPositionalInputs`; M15 second mint path → `TestNoGenericFactory` (former survivor now killed); M16 → derivation tests; M17 → manifest/cross-check tests; M18 → Limit tests; N1 pow-shape → `TestParamsCanonicalRoundTrip`; N2 → `TestHypothesisContamination`; N3 → `TestExprConstructorSurface`; N4 (redesigned to prediction-relation gate after first probe survived) → `TestInvalidCandidateExpressions`; N5/N6/N7 → respective contract tests.
- Two non-blocking observations (out of scope, not changed): `TestRepositoryTreeExact`'s allowlist is broader than §3 (membership-only, passes); `NewCandidateConcept` empty-id rejection has no dedicated test.

## E. Final status

**FREEZE-READY**

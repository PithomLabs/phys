# Adversarial Code Review — Physics Compiler MVP
**Date:** 2026-09-30  
**Auditor:** Independent adversarial reviewer (read-only)  
**Normative authority:** `specs_v2_3.md`  
**Implementation plan:** `plan10_v2_3.md`  
**Remediation history:** `plan3(20260930-112009).md`

---

## 1. EXECUTIVE VERDICT

**Implementation coverage:** PARTIAL  
**Implementation correctness:** FAIL  
**Architecture:** PARTIAL  
**Test adequacy:** FAIL  

**Overall:** FAIL

The implementation does **not** faithfully implement `specs_v2_3.md`. Multiple critical and high-severity defects remain in the session/ledger/candidate subsystem. The test suite is structurally inadequate: controlled mutation testing proves that at least seven independent semantic mutations in `Session.Validate`, `Session.Seal`, `Session.Identify`, and `UnverifiedResearchCandidate.Validate` all **survive** undetected. The core symbolic engine (`core`, `ops`, `mechanics`, `relativity` primitives) is substantially correct, but the session/ledger/candidate layer — the highest-priority semantic area — contains significant unimplemented and incorrectly implemented behavior.

---

## 2. BASELINE

**Build:** PASS (`go build ./...` clean)  
**Vet:** PASS (`go vet ./...` clean)  
**Test:** PASS (all 7 packages `ok`)  
**Go version:** go1.25.7  

**Package inventory (7 local packages):**
- `github.com/PithomLabs/phys/core`
- `github.com/PithomLabs/phys/hypothesis`
- `github.com/PithomLabs/phys/internal/kernel`
- `github.com/PithomLabs/phys/mechanics`
- `github.com/PithomLabs/phys/ops`
- `github.com/PithomLabs/phys/relativity`
- `github.com/PithomLabs/phys/session`

**File inventory (MVP tree, excluding `.opencode/`, `plans/`, `plan10/`):**
- 39 files including `go.mod`, `README.md`, `docs/paper-translation.md`
- 24 Go source files, 10 Go test files, 2 JSON manifests, 3 config/doc files
- Exact plan-pin file count matches

**Dependency inventory:**
- Zero third-party dependencies
- 7 local packages + Go stdlib only
- Module path: `github.com/PithomLabs/phys`

**Non-MVP artifacts (separated from MVP tree, not removed):**
- `.opencode/` directory (Node.js tooling artifacts)
- `plan10/` directory (plan documents)
- `plans/` directory (historical plans)
- `reality_check.md`, `reality-first-physics-v8.0.md`, `reality_first_plain_english.md`

---

## 3. PLAN 3 REMEDIATION VERIFICATION

### P0

| # | STATUS | CODE EVIDENCE | TEST EVIDENCE |
|---|--------|---------------|---------------|
| **P0-1** | **FAIL** — 17-stage replay only partially implemented | `session/session.go:264-331` implements stages 1,2,3,5,7,8,9,12,13. Missing stages 4 (decode retained inputs — inputs are already objects, not bytes, so no explicit decode), 10 (operation params vs op ID), 11 (assertion replay), 14 (replayed output hash comparison), 15 (provenance/MRC), 16 (conclusion hash), 17 (candidate containment). `replayTransformation` compares Expr/Dimension but NOT hash. `replayIdentification` compares Expr but NOT hash. | Mutation M1 (skip InputHash), M2 (skip OutputHash), M3 (skip transformation replay), M4 (skip identification replay), M5 (skip CurrentStepHash), M8 (tamper step hash) all **SURVIVED**. No test in the suite exercises `Session.Validate` for any of these stages. |
| **P0-2** | **PASS** — Internal kernel object canonical encode/decode exists | `internal/kernel/mint.go:170-230` (`LoadObjectJSON`) decodes with full invariant validation, enforces `Encode(Decode(b)) == b`. No exported `UnmarshalJSON`/`MarshalJSON` on `core.Object`. | Round-trip verified by `LoadObjectJSON` byte-equality check. |
| **P0-3** | **PASS** — Public Ledger/Step surface exported | `session/ledger.go:82-104` exports `Ledger`, `Steps()`, `DerivationHash()`, `Validate()`, `CanonicalJSON()`, `ParseLedgerJSON`. `Step` has read-only accessors for all §16.12 fields. | Surface exists; not independently tested by AST-level surface audit. |
| **P0-4** | **PARTIAL** — `Seal() (ResearchCandidate, error)` exists but drops metadata | `session/session.go:246-261` returns `(ResearchCandidate, error)`. However, `mintResearchCandidate()` (line 431-455) initializes `frameworkDependencies`, `predictions`, `falsificationConditions`, `recoveryClaims`, `anomalyReferences`, `reviewHistory` as empty slices instead of copying from `s.metadata`. | `TestSealRequiresConcludedAndHypothesis` (line 146-158) is a **stub** — calls `s.Draft()` then `_ = s`. No seal is ever invoked. |
| **P0-5** | **FAIL** — `UnverifiedResearchCandidate.Validate` is non-functional | `session/research_candidate.go:241-251` calls `json.Unmarshal(uc.bytes, &c)` into `ResearchCandidate` which has all **unexported** fields. Go's `encoding/json` silently ignores unexported fields, so `c` is always zero-valued. `c.Validate()` then fails on `c.id == ""`. The function **never** produces a `ResearchCandidate`. It does not delegate to any shared seal/validation constructor. | Mutation making `Validate` always succeed **SURVIVED**. No test verifies that external candidate JSON can be loaded successfully. |
| **P0-6** | **PASS** — `go:embed` manifest loading implemented | `mechanics/primitives.go:12-13` and `relativity/primitives.go:12-13` both use `//go:embed manifest.json`. | Tests access `ManifestJSON` directly. |

### P1

| # | STATUS | CODE EVIDENCE | TEST EVIDENCE |
|---|--------|---------------|---------------|
| **P1-1** | **PASS** — Exact Session signatures | All 12 required methods exist with exact signatures (§16.3). | No AST-level signature audit test exists. |
| **P1-2** | **PARTIAL** — DraftMetadata has wrong ReviewHistory type | `session/session.go:33-45` types `ReviewHistory []core.Object` instead of `[]core.Review` as specified in §16.2.1. Reviews are silently dropped. | No test verifies ReviewHistory preservation. |
| **P1-3** | **PASS** — ResearchCandidate accessors exact | `session/research_candidate.go:86-126` exposes 13 data accessors + `Validate` + `CanonicalJSON`. Collection accessors return copies. | No AST-level accessor audit test exists. |
| **P1-4** | **PARTIAL** — Outer-schema validation exists but is trivial | `ParseResearchCandidateJSON` (line 253-261) unmarshals into `map[string]interface{}` for outer-schema validation. However, `UnverifiedResearchCandidate.Validate` is completely broken (see P0-5). | No test verifies outer-schema rejection of malformed JSON. |
| **P1-5** | **PARTIAL** — Validation pipeline stages 1-3,5,7 implemented; stages 4,6,8-17 missing | `ResearchCandidate.Validate` (line 132-190) checks schema, ledger hash, derivation replay, hypothesis status, framework refs, falsifiability. Does NOT scan derivation outputs for HYPOTHESIS contamination. | No test verifies derivation-output provenance scan. |
| **P1-6** | **PASS** — Session-owned identification replay helper exists | `session/session.go:364-400` (`replayIdentification`) reconstructs Relation(eq,a,b), minted via `kernel.MintObject`, does NOT use `ops.Apply`. | Helper exists but is never invoked by meaningful tests. |

### P2

| # | STATUS | CODE EVIDENCE | TEST EVIDENCE |
|---|--------|---------------|---------------|
| **P2-1** | **FAIL** — Most A–S acceptance tests missing or stubbed | `TestTypedMechanicsConstructors` (A): not present. `TestNewtonSecondLawManifestMatch` (B): not present. `TestDimensionMismatch` (C): present in `ops/negative_test.go`. `TestKindMismatchEqualDimensions` (D): present. `TestAssumptionConflict` (E): present. `TestConventionConflict` (F): present. `TestDifferentiateKineticEnergy` (G): not present. `TestEnergyMomentumRelation` (H): not present. `TestMassEnergyDerivation` (I): not present. `TestLorentzFactorLimit` (J): not present. `TestSimplifyNeverIdentifies` / `TestSessionIdentifyRecords` (K): stubs. `TestHypothesisContamination` (L): stub. `TestCandidateArtifactContainmentRejected` (L): stub. `TestSealedCandidatePreservesFalsifiability` (M): stub. `TestCandidateReferencesManifestAnomaly` (N): stub. `TestMechanicsManifestCrossCheck` (O): partial. `TestRelativityManifestCrossCheck` (O): partial. `TestExactRationalRoundTrip` (P): present. `TestLedgerTamperDetection` (Q): present but tests wrong target. `TestDerivationDeterminism` (R): stub. `TestSealResearchCandidate` (S): stub. | 8 of 12 session tests are stubs or test wrong target. |
| **P2-2** | **FAIL** — Zero-assertion/tautological tests present | `TestIdentifyRequiresJustification` (line 30-42): drafts session then `_ = s` — no assertion. `TestSealRequiresConcludedAndHypothesis` (line 146-158): stub. `TestCandidateContainmentBasic` (line 160-164): `_ = true`. `TestDerivationDeterminismBasic` (line 108-120): stub. | All pass trivially. |
| **P2-3** | **PARTIAL** — Justification trimmed but untested | `Session.Identify` (line 179) uses `strings.TrimSpace(justification) == ""`. Mutation removing `TrimSpace` **SURVIVED**. | `TestSessionIdentifyRequiresJustificationNonEmpty` (line 55-71) tests empty justification but NOT whitespace-padded justification. |
| **P2-4** | **PARTIAL** — Typed state errors returned but untested | `ProvenanceError` returned for state violations. `LedgerValidationError` returned for validation failures. | No test verifies exact error types for state transitions. |
| **P2-5** | **FAIL** — Mutation probes not killed | See mutation results section. 7 of 10 session-related mutation probes **SURVIVED**. | See §17. |

---

## 4. CRITICAL / HIGH FINDINGS

### F-001: Session.Validate Missing 7 of 17 Stages
- **Severity:** CRITICAL
- **Classification:** IMPLEMENTED INCORRECTLY
- **Spec anchor:** §16.19
- **Plan anchor:** §12.7 (REQ-§16.19-MUST-01)
- **File/function:** `session/session.go:264-331`
- **Observed:** Stages 10 (operation params vs op ID), 11 (assertion replay), 14 (replayed output hash comparison), 15 (provenance/MRC validation), 16 (conclusion hash), 17 (candidate containment) are completely absent. Stages 4, 6 (decode retained canonicals) are implicit since canonicals are stored as `core.Object` values, not bytes.
- **Required:** All 17 stages in exact order with `LedgerValidationError` on any mismatch.
- **Current test:** No test exercises `Session.Validate` for any missing stage. Existing tests (`TestValidateMutation_EmptyLedger`, `TestValidateMutation_StepHashTamper`, `TestSessionValidateMutation`) test `Ledger.Validate`, not `Session.Validate`, and test a *separate* parsed ledger rather than the session's own ledger.
- **Mutation result:** M1 (InputHash), M2 (OutputHash), M3 (transformation replay), M4 (identification replay), M5 (CurrentStepHash), M8 (step hash tamper) all **SURVIVED**.

### F-002: UnverifiedResearchCandidate.Validate Non-Functional
- **Severity:** CRITICAL
- **Classification:** IMPLEMENTED INCORRECTLY
- **Spec anchor:** §26.10, REQ-032-10a
- **Plan anchor:** §12.10
- **File/function:** `session/research_candidate.go:241-251`
- **Observed:** `json.Unmarshal(uc.bytes, &c)` into `ResearchCandidate` (all unexported fields) silently ignores all JSON fields. `c` is always zero-valued. `c.Validate()` always fails on `c.id == ""`. The function **never** returns a `ResearchCandidate`.
- **Required:** Execute full validation pipeline (§26.8) before returning success. Must delegate to shared seal/validation constructor.
- **Current test:** `TestCandidateContainmentBasic` (line 160-164) is `_ = true`. No test verifies external candidate loading.
- **Mutation result:** Mutation making `Validate` always succeed **SURVIVED**.

### F-003: Seal Drops Candidate Metadata
- **Severity:** CRITICAL
- **Classification:** IMPLEMENTED INCORRECTLY
- **Spec anchor:** §16.11, §26.2
- **Plan anchor:** §12.10
- **File/function:** `session/session.go:431-455` (`mintResearchCandidate`)
- **Observed:** `frameworkDependencies`, `predictions`, `falsificationConditions`, `recoveryClaims`, `anomalyReferences`, `reviewHistory` initialized as empty slices. `s.metadata` fields are never copied.
- **Required:** Copy all metadata slices from `s.metadata` to candidate.
- **Current test:** `TestSealRequiresConcludedAndHypothesis` is a stub. No test verifies metadata preservation.
- **Mutation result:** Removing Seal's call to `s.Validate()` **SURVIVED**.

### F-004: Identify Missing Kind Compatibility (MRC-003)
- **Severity:** HIGH
- **Classification:** IMPLEMENTED INCORRECTLY
- **Spec anchor:** §16.7.3, §6.2
- **Plan anchor:** §12.6
- **File/function:** `session/session.go:175-212` (`Session.Identify`)
- **Observed:** Checks only `a.Dimension().Equal(b.Dimension())`. Does NOT check kind compatibility. Two objects of different named kinds with equal dimensions (e.g., `Mass` vs `RestMass`) are allowed to identify.
- **Required:** Enforce Compare-kind rules (same named kind, or named+Expression, or Expression+Expression — not different named kinds).
- **Current test:** `TestSessionIdentifyRecords` not present.
- **Mutation result:** Removing HYPOTHESIS contamination from Identify **SURVIVED** (M12).

### F-005: ZeroThreeMomentum Missing RestFrameAssumption
- **Severity:** HIGH
- **Classification:** IMPLEMENTED INCORRECTLY
- **Spec anchor:** §19.7, §20.1
- **Plan anchor:** §11 (Step 11)
- **File/function:** `relativity/relations.go:91-106`
- **Observed:** `Assumptions: kernel.NewAssumptionSet()` — empty. Missing `RestFrameAssumption()`.
- **Required:** `Assumptions: kernel.NewAssumptionSet(RestFrameAssumption())`
- **Current test:** `derivation_test.go` logs assumptions but never asserts presence of `RestFrameAssumption`.

### F-006: Relativity Assumption Keys Mismatch Manifest
- **Severity:** HIGH
- **Classification:** IMPLEMENTED INCORRECTLY
- **Spec anchor:** §11.6, §19.4, §19.5
- **Plan anchor:** §8 (assumption keys)
- **File/function:** `relativity/relations.go:15,46,48,72,74`
- **Observed:** 
  - `LorentzFactor()`: key `"c_positive"` (should be `"speed_of_light_positive"`)
  - `EnergyMomentumRelation()`: keys `"m_nonnegative"`, `"c_positive"` (should be `"rest_mass_nonnegative"`, `"speed_of_light_positive"`)
  - `MassEnergyRelation()`: keys `"m_nonnegative"`, `"c_positive"` (should be `"rest_mass_nonnegative"`, `"speed_of_light_positive"`)
- **Required:** Exact keys matching manifest's `assumptions` arrays.
- **Current test:** `manifest_test.go` never asserts assumption keys.

### F-007: Relativity Manifest Missing 4 Framework Assumptions
- **Severity:** HIGH
- **Classification:** IMPLEMENTED INCORRECTLY
- **Spec anchor:** §19.3
- **Plan anchor:** §8
- **File/function:** `relativity/manifest.json` (top-level `assumptions` array)
- **Observed:** Only 2 items: `rest_mass_nonnegative`, `speed_of_light_positive`. Missing: Minkowski spacetime, Lorentz symmetry, No gravitational dynamics in package, Special-relativistic regime.
- **Required:** 4 framework-level assumptions as specified in §19.3.

### F-008: ResearchCandidate.Validate Skips Derivation Output Provenance
- **Severity:** HIGH
- **Classification:** IMPLEMENTED INCORRECTLY
- **Spec anchor:** MRC-008, §14.8, §26.8
- **Plan anchor:** §26.9
- **File/function:** `session/research_candidate.go:132-190`
- **Observed:** Checks `Hypothesis` field status but does NOT iterate derivation ledger steps to verify no output with HYPOTHESIS provenance is presented with a trusted status.
- **Required:** Scan derivation outputs; fail with `CandidateContainmentError` if any hypothesis-dependent step output carries trusted status.
- **Current test:** `TestCandidateContainmentBasic` is a no-op.

### F-009: CandidateContainmentError Missing Unwrap()
- **Severity:** HIGH
- **Classification:** IMPLEMENTED INCORRECTLY
- **Spec anchor:** §30
- **Plan anchor:** §30
- **File/function:** `internal/kernel/types.go:118-125`
- **Observed:** `CandidateContainmentError` implements `Error()` but not `Unwrap()` or `Is()`. `errors.As` cannot unwrap to this type.
- **Required:** Implement `Unwrap() error` (return nil) or `Is(error) bool` per §30.
- **Current test:** `core/expr_test.go` line ~1452 performs `errors.As` for all 11 error types including `CandidateContainmentError` — this call **fails at runtime**.

---

## 5. MEDIUM / LOW FINDINGS

### F-010: DraftMetadata.ReviewHistory Typed []core.Object
- **Severity:** MEDIUM
- **Classification:** IMPLEMENTED INCORRECTLY
- **Spec anchor:** §16.2.1
- **File/function:** `session/session.go:44`
- **Observed:** `ReviewHistory []core.Object` instead of `[]core.Review`.

### F-011: Session.CanonicalJSON Missing schema_version
- **Severity:** MEDIUM
- **Classification:** PLAN-LEVEL CONFORMANCE GAP
- **Spec anchor:** §10.4
- **Plan anchor:** §4
- **File/function:** `session/session.go:419-425`
- **Observed:** `sessionDTO` has no `SchemaVersion` field.

### F-012: Ledger.CanonicalJSON Missing schema_version
- **Severity:** MEDIUM
- **Classification:** PLAN-LEVEL CONFORMANCE GAP
- **Spec anchor:** §10.4
- **Plan anchor:** §4
- **File/function:** `session/ledger.go:141-145`
- **Observed:** DTO has no `SchemaVersion` field.

### F-013: Prediction.Relation Typed core.Object Instead of core.Expr
- **Severity:** MEDIUM
- **Classification:** IMPLEMENTED INCORRECTLY
- **Spec anchor:** §25.1, §26.0
- **Plan anchor:** §26.0
- **File/function:** `session/research_candidate.go:23,30,38`
- **Observed:** `Prediction.Relation`, `FalsificationCondition.ContradictingCondition`, `RecoveryClaim.Condition` all typed `core.Object` instead of `core.Expr`.

### F-014: MintObjectMust Exists in Production Code
- **Severity:** MEDIUM
- **Classification:** PLAN-LEVEL CONFORMANCE GAP
- **Spec anchor:** §5.2.1
- **Plan anchor:** P3-5
- **File/function:** `internal/kernel/mint.go:91-100`
- **Observed:** `MintObjectMust` is a second minting helper in production code. Comment says "FOR TEST USE ONLY" but it is not in a `_test.go` file.

### F-015: session_test.go Imports mechanics
- **Severity:** MEDIUM
- **Classification:** PLAN-LEVEL CONFORMANCE GAP
- **Spec anchor:** §4
- **Plan anchor:** §2, P3-3
- **File/function:** `session/session_test.go:9`
- **Observed:** Imports `"github.com/PithomLabs/phys/mechanics"` and uses it for fixtures throughout. Plan §2 states `session_test.go` should be external package but does not list mechanics import.

---

## 6. TEST QUALITY FINDINGS

### Zero-Assertion Tests
- `TestIdentifyRequiresJustification` (`session/session_test.go:30-42`): drafts session, then `_ = s`. No assertion.
- `TestSealRequiresConcludedAndHypothesis` (`session/session_test.go:146-158`): drafts session, then `_ = s`. No seal ever invoked.
- `TestCandidateContainmentBasic` (`session/session_test.go:160-164`): `_ = true`.
- `TestDerivationDeterminismBasic` (`session/session_test.go:108-120`): drafts session, then `_ = s`.

### Tautological Tests
- `TestSessionSurfaceNineActions` (`session/session_test.go:126-140`): drafts session, checks `s != nil`.
- `TestSessionStateMachineTransitions` (`session/session_test.go:88-102`): drafts session, checks `s != nil`.

### Weakly Assertioned Tests
- `TestValidateMutation_StepHashTamper` (`session/session_test.go:183-245`): creates a *separate* `Ledger2` from tampered JSON and calls `ledger2.Validate()`. Does NOT corrupt the session's own ledger and call `s.Validate()`. Tests the wrong target.
- `TestSessionValidateMutation` (`session/session_test.go:247-301`): same structural flaw — modifies JSON, re-parses into `ledger2`, calls `ledger2.Validate()`.
- `TestLedgerParseAndValidateExists` (`session/session_test.go:81-86`): only checks that parse returns error for empty steps.

### Tests That Can Survive Behavior Removal
- All session mutation probes (M1-M5, M8, M11, M12) **SURVIVED** — removing validation stages, tamper detection, justification trimming, and HYPOTHESIS contamination all pass tests.

---

## 7. MUTATION RESULTS

| Mutant | Expected | Actual | Killed/Survived |
|--------|----------|--------|-----------------|
| M1: Skip InputHash verification | FAIL | ok | **SURVIVED** |
| M2: Skip OutputHash verification | FAIL | ok | **SURVIVED** |
| M3: Skip transformation replay | FAIL | ok | **SURVIVED** |
| M4: Skip identification replay | FAIL | ok | **SURVIVED** |
| M5: Skip CurrentStepHash recomputation | FAIL | ok | **SURVIVED** |
| M8: Tamper step hash undetectable | FAIL | ok | **SURVIVED** |
| M9: Pow(x,0) always returns 1 | FAIL | FAIL | **KILLED** |
| M10: MRC-003 broken | FAIL | FAIL | **KILLED** |
| M11: No justification trimming | FAIL | ok | **SURVIVED** |
| M12: No HYPOTHESIS contamination in Identify | FAIL | ok | **SURVIVED** |

**7 of 10 probes SURVIVED.** The 3 killed (M9, M10) are in `ops`, not in session/ledger/candidate.

---

## 8. COVERAGE AUDIT

### REQ-032 / MUST-MUST NOT Coverage
- REQ-032-01..07, 18..20: Partially covered in `ops/negative_test.go`.
- REQ-032-08 (justification trimming): **Partially covered** — empty justification tested, whitespace-padded justification NOT tested.
- REQ-032-09 (candidate containment): **Not covered** — `TestCandidateArtifactContainmentRejected` is a stub.
- REQ-032-10a (UnverifiedResearchCandidate): **Not covered** — `UnverifiedResearchCandidate.Validate` is non-functional.
- REQ-032-11..14 (ledger tamper detection): **Incorrectly tested** — tests verify `Ledger.Validate`, not `Session.Validate`.
- REQ-032-21/22 (typed state errors): **Not independently tested** — no test verifies exact error types for state violations.

### MRC Coverage
- MRC-001: PASS (constructor authority verified)
- MRC-002: PASS (dimension checks in ops)
- MRC-003: **PARTIAL** — ops.Add/Subtract/Compare enforce it; `Session.Identify` does NOT.
- MRC-004: PASS (assumption conflict detection)
- MRC-005: PASS (convention conflict detection)
- MRC-006: **PARTIAL** — `Session.Identify` exists but lacks kind check; justification trimming untested; HYPOTHESIS contamination untested.
- MRC-007: **PARTIAL** — Session authority exists but `Session.Validate` incomplete; `UnverifiedResearchCandidate.Validate` broken.
- MRC-008: **PARTIAL** — ops contamination law correct; `ResearchCandidate.Validate` does NOT scan derivation outputs.

### A–S Acceptance Coverage
| Test | Status |
|-------|--------|
| A TestTypedMechanicsConstructors | **NOT PRESENT** |
| B TestNewtonSecondLawManifestMatch | **NOT PRESENT** |
| C TestDimensionMismatch | PRESENT |
| D TestKindMismatchEqualDimensions | PRESENT |
| E TestAssumptionConflict | PRESENT |
| F TestConventionConflict | PRESENT |
| G TestDifferentiateKineticEnergy | **NOT PRESENT** |
| H TestEnergyMomentumRelation | **NOT PRESENT** |
| I TestMassEnergyDerivation | **NOT PRESENT** |
| J TestLorentzFactorLimit | **NOT PRESENT** |
| K TestSimplifyNeverIdentifies / TestSessionIdentifyRecords | **STUBS** |
| L TestHypothesisContamination / TestCandidateArtifactContainmentRejected | **STUBS** |
| M TestSealedCandidatePreservesFalsifiability | **STUB** |
| N TestCandidateReferencesManifestAnomaly | **STUB** |
| O TestMechanics/RelativityManifestCrossCheck | PARTIAL |
| P TestExactRationalRoundTrip | PRESENT |
| Q TestLedgerTamperDetection | PRESENT (wrong target) |
| R TestDerivationDeterminism | **STUB** |
| S TestSealResearchCandidate | **STUB** |

---

## 9. ARCHITECTURE / IMPORT / EXPORT AUDIT

**Package graph:** PASS
- `core` → `internal/kernel` only
- `ops` → `core`, `internal/kernel`
- `session` → `core`, `ops`, `internal/kernel`
- `mechanics` → `core`, `internal/kernel`
- `relativity` → `core`, `internal/kernel`
- `hypothesis` → `core`, `internal/kernel`
- No cycles

**Public API boundary:** PASS
- `MintObject` exists only in `internal/kernel`
- `core` exposes no mint forwarder
- `Object` fields unexported

**Mint authority:** PARTIAL
- `MintObjectMust` exists in production code (`internal/kernel/mint.go:91-100`) — plan says only `MintObject` is specified.

**Forbidden APIs:** PASS (in production code)
- No `Promote`, `Trust`, `ApproveHypothesis`, `SetCorpusStatus`, `OverrideMRC`, `BypassMRC`
- No `float32`/`float64` in non-test sources
- No `math/rand`, `go/parser`, `go/scanner`, `go/token`, `net/http`, `os/exec`
- No `ops.Identify`
- No review package

**map[string]any:** Present only in `session/research_candidate.go:255` for outer-schema validation — acceptable per §26.10.

---

## 10. PHYSICS-TRACE AUDIT

### F=ma
`mechanics.NewtonSecondLaw()` constructs `Relation(eq, Symbol("F"), Mul(Symbol("m"), Symbol("a")))` with Kind=Relation, Dimension=Force. Manifest matches. **PASS.**

### d/dv (1/2 mv²)
`mechanics.NewKineticEnergy(m, v)` constructs `Mul(Rational(1/2), m.Expr(), Pow(v.Expr(), 2))` with Kind=KineticEnergy, Dimension=Energy. `Differentiate` on this expression with `wrt=v` applies product rule and power rule: `1/2 * m * 2 * v^(2-1) * d(v)/dv` → `m * v`. The result is simplified. **PASS** (ops/transform.go:254-334).

### E=mc² derivation
The required sequence is:
1. `EnergyMomentumRelation()` — present
2. `ZeroThreeMomentum()` — **MISSING RestFrameAssumption** (F-005)
3. `Substitute` — works
4. `Simplify` — works
5. `Solve(..., Energy())` — works
6. `Compare(Energy(), ZeroEnergy(), gte)` — works
7. `SelectBranch` — works (mass-energy path: `Sqrt(Pow(Mul(m,Pow(c,2)),2)) → m*c²` because `m≥0` and `c²` nonnegative)

The derivation CAN mechanically execute, but step 2 loses the explicit rest-frame constraint metadata.

### lim(v→0) γ
`Limit(LorentzFactor(), Velocity(), ZeroVelocity())` expands fixed body `1/Sqrt(1-Pow(v/c,2))`, substitutes `v=0`, simplifies → `Rational(1)`. **PASS** — body traversal is real, not an ID shortcut.

### Hypothesis contamination
`ops` contamination law correct: any HYPOTHESIS input → HYPOTHESIS output for all 12 ops. **PASS.**
`Session.Identify` contamination: implemented but **UNTESTED** (M12 survived).
`ResearchCandidate` derivation-output scan: **MISSING** (F-008).

### Candidate sealing
`Session.Seal` checks Hypothesis status and calls `Validate()`, but `Validate()` is incomplete and `mintResearchCandidate` drops all metadata slices. **FAIL.**

---

## 11. PLAN-LEVEL GAPS

### Specification Violations
1. §16.19 — 17-stage Validate incomplete (F-001)
2. §26.10 — UnverifiedResearchCandidate.Validate non-functional (F-002)
3. §16.11 — Seal drops metadata (F-003)
4. §16.7.3 — Identify missing kind check (F-004)
5. §19.7 — ZeroThreeMomentum missing RestFrameAssumption (F-005)
6. §11.6 / §19.4/5 — Relativity assumption keys mismatch (F-006)
7. §19.3 — Relativity manifest incomplete (F-007)
8. §26.8 — ResearchCandidate.Validate skips derivation-output provenance scan (F-008)
9. §30 — CandidateContainmentError missing Unwrap (F-009)
10. §16.2.1 — DraftMetadata.ReviewHistory wrong type (F-010)
11. §25.1, §26.0 — Prediction.Relation typed Object instead of Expr (F-013)

### Plan 10 Conformance Gaps
1. **P3-3** — `session_test.go` imports `mechanics` (plan §2 says external session tests should not import domain packages)
2. **P3-4** — Relativity assumption keys wrong
3. **P3-5** — `MintObjectMust` exists in production code
4. **P3-6** — `sessionDTO` and Ledger DTO missing `schema_version: "1"`
5. **Plan §11** — A–S acceptance tests: 8 missing, 4 stubs, 2 test wrong target

### Documentation-Only Gaps
- README is accurate and does not make false claims.
- `docs/paper-translation.md` has exactly the three required sections: Common notation, Framework mapping, Ambiguity resolution. **PASS.**

---

## 12. REMEDIATION

Ordered by severity:

1. **[CRITICAL]** Implement missing `Session.Validate` stages: operation params validation, assertion replay, replayed output hash comparison, provenance/MRC validation, conclusion hash, candidate containment derivation-output scan.
2. **[CRITICAL]** Fix `UnverifiedResearchCandidate.Validate` to use a typed intermediate DTO and execute the full validation pipeline (§26.8) before returning success.
3. **[CRITICAL]** Fix `mintResearchCandidate` to copy ALL metadata slices from `s.metadata` (FrameworkDependencies, Predictions, FalsificationConditions, RecoveryClaims, AnomalyReferences, ReviewHistory).
4. **[HIGH]** Add kind compatibility check to `Session.Identify` (MRC-003 / §6.2 Compare-kind rules).
5. **[HIGH]** Add `RestFrameAssumption()` to `ZeroThreeMomentum()`.
6. **[HIGH]** Fix relativity assumption keys to match manifest: `"rest_mass_nonnegative"`, `"speed_of_light_positive"`.
7. **[HIGH]** Add 4 missing framework assumptions to `relativity/manifest.json`.
8. **[HIGH]** Add `Unwrap() error` to `CandidateContainmentError`.
9. **[HIGH]** Add derivation-output provenance scan to `ResearchCandidate.Validate`.
10. **[MEDIUM]** Fix `DraftMetadata.ReviewHistory` type to `[]core.Review`.
11. **[MEDIUM]** Add `schema_version: "1"` to `sessionDTO` and Ledger canonical JSON DTOs.
12. **[MEDIUM]** Change `Prediction.Relation`, `FalsificationCondition.ContradictingCondition`, `RecoveryClaim.Condition` from `core.Object` to `core.Expr`.
13. **[MEDIUM]** Move `MintObjectMust` to test-only or remove from production.
14. **[MEDIUM]** Replace stub tests with real assertions (TestSealResearchCandidate, TestCandidateContainmentBasic, TestDerivationDeterminism, etc.).
15. **[MEDIUM]** Fix mutation-susceptible tests: TestValidateMutation_StepHashTamper and TestSessionValidateMutation should corrupt the session's own ledger and call `s.Validate()`, not parse a separate ledger.

---

## 13. FINAL VERDICT

**FAIL**

The current implementation cannot survive an independent adversarial attempt to falsify its conformance and correctness claims. The session/ledger/candidate layer — the highest-priority semantic area per the specification — contains multiple critical defects that are invisible to the existing test suite. Seven independent mutation probes in core validation and sealing logic all survived undetected. The test suite has four zero-assertion stubs and tests the wrong target for tamper detection. While the core symbolic engine (`core`, `ops`) and the domain packages (`mechanics`, `relativity` primitives) are substantially correct, the session/ledger/candidate subsystem is incomplete and contains significant unimplemented specification requirements.

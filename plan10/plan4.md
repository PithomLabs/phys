# Implementation Plan — Physics Compiler MVP Remediation (Post-Adversarial Review) — Plan 4

**Normative:** `plan10/specs_v2_3.md` | **Plan:** `plan10/plan10_v2_3.md` | **Findings:** `plan10/adv_review.md` + `plan10/review.md`
**Constraints locked by user:** breaking `DraftMetadata` type fix; layered validation; typed DTO + shared ctor; remove `MintObjectMust` from prod; test-first on F-009.

## 1. Actual current state (from source, not history)

| # | Question | Answer in current tree |
|---|---|---|
| 1 | Canonical session step | `session/session.go:stepEntry` (unexported draft buffer: label/kind/op/inputs/params/output as `core.Object` values) |
| 2 | Committed ledger step | `session/ledger.go:stepEnvelope` (16 section 16.12 fields; `InputCanonicals []core.Object`, `OutputCanonical core.Object`, `ParamsCanonical []byte`) |
| 3 | Where canonicals stored | In `stepEnvelope` at commit time (`newLedgerInternal`); retained as live `core.Object`, not bytes |
| 4 | How decoded | **Not decoded** — replay reuses live `core.Object` values directly; `internal/kernel.LoadObjectJSON` exists (`mint.go:170`) but session replay never calls it |
| 5 | `Session.Validate` | `session.go:264-331`: partial — structure/genesis/StepID/chain, InputHash/OutputHash/AssumptionHash/ConventionHash, transformation + identification replay (Expr/Dimension compare only). **Missing:** params-vs-opID, assertion replay, replayed-output-hash compare, provenance/MRC, conclusion hash, containment |
| 6 | `Ledger.Validate` | `ledger.go:107`: hash-only (index/StepID/CurrentStepHash/chain). No replay, no hash-vs-canonical checks |
| 7 | `Seal` to candidate | `session.go:246-261` + `mintResearchCandidate:431-455`: checks Concluded + HYPOTHESIS + `Validate()`, then builds `ResearchCandidate` with **empty slices** for all 6 metadata collections |
| 8 | External candidate load | `research_candidate.go:253-261` outer-schema only to `UnverifiedResearchCandidate{bytes}`; `Validate:241-251` does `json.Unmarshal(uc.bytes, &c)` into unexported-field struct = **always zero-valued, never succeeds** |
| 9 | Candidate validation | `ResearchCandidate.Validate:132-190`: schema/ledger-hash/`Ledger.Validate`/top-level Hypothesis/framework/falsifiability/review-category. **Missing:** derivation-output provenance scan (MRC-008) |
| 10 | Metadata copied/dropped | `DraftMetadata` (`session.go:33-45`) wrongly typed (`FrameworkDependencies/Predictions/etc []core.Object`, `ReviewHistory []core.Object`); `mintResearchCandidate` drops all of them anyway |
| 11 | Provenance/MRC checked | Only in `ops` + ad-hoc `Identify` contamination; no session replay provenance check, no containment scan |
| 12 | Identify replay reconstruction | `replayIdentification:364-400` duplicates (not shares) `Identify:175-212` logic; both miss kind-compat, merged assumptions/conventions, justification-in-provenance verification |

**Corroborating defects:** `relativity/relations.go` uses `c_positive`/`m_nonnegative` keys, `ZeroThreeMomentum` has empty assumptions; `relativity/manifest.json` top-level assumptions only 2 items; `Prediction.Relation` etc typed `core.Object` not `core.Expr`; `MintObjectMust` in prod `mint.go:91`; `session_test.go` is stub-heavy + imports `mechanics`.

## 2. Root causes (fix these, not symptoms)

1. **Two replay models:** `Session.Validate` vs `Ledger.Validate` diverged; neither is authoritative.
2. **No canonical-bytes replay substrate:** objects retained live, so decode + invariant-check + hash-verify stages are vacuous.
3. **Identify split-brain:** construct path vs replay path duplicated, both incomplete (no MRC-003 kind rule, no merged metadata, no contamination proof).
4. **Two candidate construction paths (actually 1.5 broken paths):** `Seal` vs `Unverified.Validate` share nothing; latter is non-functional due to unexported-field unmarshal.
5. **Wrong semantic model in types:** `DraftMetadata` + `Prediction/FalsificationCondition/RecoveryClaim` use `core.Object` where spec mandates `core.Expr`; `ReviewHistory` wrong element type — so metadata *cannot* survive even if copied.
6. **Tests prove existence, not semantics:** stubs, `_ = s`, wrong-target tamper tests → 7/10 mutants survive.

## 3. Implementation plan (within exact 39-file tree, no new packages/files)

### Stage 0 — Baseline + scaffolding (no semantics change)
- In disposable copy only: `go version; go build ./...; go vet ./...; go test ./...; go list ./...`. Record baseline.

### Stage 1 — Type-model correction (breaking, user-approved)
**Files:** `session/session.go`, `session/research_candidate.go`
- `DraftMetadata` to exact v2.3: `Hypothesis core.Object; Premises []core.Object; Assumptions core.AssumptionSet; FrameworkDependencies []FrameworkDependency; Predictions []Prediction; FalsificationConditions []FalsificationCondition; RecoveryClaims []RecoveryClaim; AnomalyReferences []AnomalyReference; ReviewHistory []core.Review`. Keep `DerivationID/Label` handling; copy all slices at `Draft`.
- `Prediction.Relation -> core.Expr`; `FalsificationCondition.ContradictingCondition -> core.Expr`; `RecoveryClaim.Condition -> core.Expr`.
- Update `mintResearchCandidate`, DTOs, canonical JSON, all callers/tests consistently. No compat shim.

### Stage 2 — Relativity metadata fix
**Files:** `relativity/relations.go`, `relativity/manifest.json`
- Keys to `rest_mass_nonnegative`, `speed_of_light_positive` in `LorentzFactor`, `EnergyMomentumRelation`, `MassEnergyRelation`.
- `ZeroThreeMomentum`: `Kind=ThreeMomentum, Dimension=Momentum, Expr=Rational(0), Provenance=DEFINED, Assumptions={RestFrameAssumption()}`.
- Manifest top-level `assumptions[4]` with pinned kinds/keys/values (Minkowski spacetime, Lorentz symmetry, No gravitational dynamics, Special-relativistic regime). Keep `metric.signature=-+++`, 10 items, `Velocity`/zeros not items.
- Verify `NewRestMass`/`NewSpeedOfLight` in `primitives.go` already correct — do not touch.

### Stage 3 — Kernel/test-hygiene (small) + C3
**Files:** `internal/kernel/mint.go` (no new files)
- C3 — eliminate `MintObjectMust` from production: refactor all production callers (domain constructors, ops, session, hypothesis) to explicit `obj, err := kernel.MintObject(spec)` handling. Do NOT create `internal/kernel/*_test.go` (not in the 39-file tree). A Must-style helper may persist ONLY inside an already-existing `*_test.go` file in the tree and test-only; otherwise delete it. AST-verify single prod `func MintObject`.
- F-009: **do not add `Unwrap`**; add direct `errors.As(err,&CandidateContainmentError{})` probe test first; only change wrapping if probe fails.

### Stage 4 — Session core: Identify + private reconstruction helper
**File:** `session/session.go`
- `Identify`: keep state/TrimSpace-justification checks; add validity of both operands; dimension equality; **Compare-kind rules** (same named kind, named+Expression, Expression+Expression allowed; different named kinds to `CategoryMismatchError`; Mass vs RestMass is the negative fixture); build `Relation(eq,a.Expr(),b.Expr())`; contamination (either HYPOTHESIS to HYPOTHESIS else IDENTIFIED); **merged assumptions/conventions**; justification in provenance; correct Source/Framework inheritance; `Kind:"identify"` params.
- Extract `reconstructIdentification(a,b core.Object, justification string) (core.Object,error)`: pure, no session mutation. Both `Identify` and `replayIdentification` call it; replay then compares full canonical bytes/hash, not just Expr.

### Stage 5 — Layered validation (user-approved: structural Ledger, authoritative Session, shared leaf primitives) + C1/C2
**Files:** `session/session.go`, `session/ledger.go` (no new files; private helpers inside these)
- C1 — canonical-byte replay substrate (mandatory): Draft: live `core.Object` values are acceptable in the unexported draft buffer. Commit: canonicalize each input/output object ONCE via `CanonicalObjectJSON`, retain the resulting IMMUTABLE CANONICAL BYTES in Step as the authoritative replay substrate, and compute InputHash/OutputHash FROM those retained bytes. `ParamsCanonical` (already `[]byte`) is retained likewise. Validate (both Ledger structural checks and Session 17-stage replay): decode retained bytes via `internal/kernel.LoadObjectJSON` (invariant check + `Encode(Decode(b))==b`), verify hashes against decoded values, and replay ONLY from decoded retained values. Never regenerate canonical bytes from live objects during validation. Public `Step.InputCanonicals()/OutputCanonical()` expose decoded value copies, while committed `Step` internally retains canonical bytes.
- C2 — leaf-only helpers + exact ordering: shared helpers MUST be atomic leaf primitives only, e.g. `checkStepIDFormat, decodeRetainedInput(i), verifyInputHash(i), decodeRetainedOutput, verifyOutputHash, verifyMetadataHashes, verifyChainLink, decodeParams, checkParamsBinding`. No shared helper may internally perform more than one section 16.19 stage. `Session.Validate` alone orchestrates stages 1..17 in exact order (section 16.19-MUST-01). `Ledger.Validate` calls only the structural subset and MUST NOT invoke assertion/transformation/identification replay, provenance, conclusion, or containment stages.
- `Ledger.Validate` (structural/integrity only): structure, StepID/index, `HashObject(decoded Input)==InputHash`, output hash, assumption/convention hashes, `ParamsCanonical` validity, genesis/chain, `CurrentStepHash` recompute. Return `LedgerValidationError`/wrapped typed errors. Must **not** call `Session.Validate`.
- `Session.Validate` — exact section 16.19 17-stage order: (1) structure/state (2) genesis (3) StepID/index (4) decode retained inputs via internal kernel codec (`CanonicalObjectJSON` to `LoadObjectJSON` round-trip check) (5) InputHash (6) decode output (7) OutputHash (8) Assumption/ConventionHash (9) chain recompute (10) params-vs-opID via `ops.ValidateOperationParams` (11) assertion replay (12) transformation replay via `ops.Apply` at section 15.13.1 indices (13) identification replay via private helper (14) replayed-output hash compare (15) provenance/MRC (16) conclusion hash (17) containment. Any mismatch to `LedgerValidationError`.
- `replayTransformation`: compare full `EqualObject`/hash + dimension/kind, not Expr-only.
- Fix `CanonicalJSON` DTOs (`sessionDTO`, ledger DTO) to include `schema_version:"1"` (DTO-only, never kernel state).

### Stage 6 — Seal + external loader: one trusted path
**Files:** `session/session.go`, `session/research_candidate.go`
- Private `sealOrValidateCandidate(ledger Ledger, meta DraftMetadata, conclusionHash string) (ResearchCandidate,error)` used by **both** `Seal` and `UnverifiedResearchCandidate.Validate`.
- `Seal`: require Concluded; require valid HYPOTHESIS hypothesis (`ProvenanceError` otherwise, stay Concluded); run full `Validate()`; run containment; **copy all metadata slices**; `DerivationHash==LedgerHash`; transition to Sealed only on success; mark ledger sealed.
- `ParseResearchCandidateJSON` to `UnverifiedResearchCandidate` only (outer-schema check). Add typed intermediate DTO (exported-field struct for JSON only); `Unverified.Validate`: decode DTO to canonical-bytes check to ledger-hash check to reconstruct ledger via `ParseLedgerJSON`+structural validate to reconstitute `core.Object`s via kernel decoder to delegate to shared ctor to return trusted candidate. No `core.Object` accessors on unverified type; no exported-field change to `ResearchCandidate`; no second trust path.
- C4 — replay-derived containment (mandatory algorithm): for each committed step, using retained canonical inputs + `ParamsCanonical`: assertion → retained output provenance is authoritative; check status-class rules (Postulate→POSTULATED, Define→DEFINED, Declare→preserved) and byte-equality with retained output. Transformation → `ops.Apply` replay from decoded inputs/params gives expected object; expected provenance = HYPOTHESIS if ANY input HYPOTHESIS else DERIVED (§13.2); compare expected status + full canonical hash vs retained. Identification → session-owned `reconstructIdentification` gives expected object; compare full canonical hash + status (IDENTIFIED unless HYPOTHESIS input). Containment violation (expected HYPOTHESIS but retained trusted status, or retained `corpus_status != NONE` on hypothesis-derived output, or top-level Hypothesis not HYPOTHESIS) → `CandidateContainmentError`. Containment is thus a replay consequence shared by `Session.Validate` stage 17 and `ResearchCandidate.Validate` §26.8 via the same private helper. Keep existing framework/falsifiability/review-category checks.

### Stage 7 — Tests: replace stubs with semantic A–S + `TestApplyPositionalInputs`
**Files (only):** `session/session_test.go`, `hypothesis/candidate_test.go`, `ops/negative_test.go`/`operations_test.go`, `mechanics/relations_test.go`+`manifest_test.go`, `relativity/derivation_test.go`+`manifest_test.go`, `core/expr_test.go`
- Implement normative names exactly: A,B,G,H,I,J,K(`TestSimplifyNeverIdentifies`+`TestSessionIdentifyRecords`),L(`TestHypothesisContamination`+`TestCandidateArtifactContainmentRejected`),M,N,O (`TestMechanicsManifestCrossCheck`, `TestRelativityManifestCrossCheck`),P,Q,R,S + `TestApplyPositionalInputs`.
- Q must corrupt **session's own committed ledger** then `s.Validate()` for all section 16.20 forms (output-only, output+hash, full-chain-recompute to replay divergence, broken chain). Delete wrong-target `Ledger2.Validate` pattern.
- S: full Draft to derive to Commit to Conclude to Seal with real hypothesis/metadata; assert metadata preservation, `LedgerHash==DerivationHash`, canonical round-trip via `ParseResearchCandidateJSON(...).Validate()`, post-seal mutation to `ProvenanceError`.
- K: Simplify-never-IDENTIFIED; Identify records justification+step; empty/whitespace justification fails; Mass-vs-RestMass kind rejection; HYPOTHESIS contamination.
- R: same derivation twice in-process + fresh Session to byte-identical canonical bytes/hashes.
- Keep `session_test.go` as `package session_test`; P3-3 mechanics-import is lower priority.

### Stage 8 — Mutation + determinism + audit gates
- Disposable-copy mutation probes (min 21 listed in brief section 12): each must be killed by Q/K/L/S/R tests. Record MUTANT/EXPECTED/ACTUAL/KILLED table.
- Canonical audit: Expr/Dimension/AssumptionSet/ConventionSet/Provenance/Object/Manifest/Params/Step/Envelope/Ledger/Candidate/Prediction/FalsificationCondition/RecoveryClaim/AnomalyReference/Review/Challenge — typed structs, fixed order, sorted sets, `"num/den"`, no maps/timestamps/RNG/pointers, lowercase hashes.
- Source audit: 39 files, stdlib only, `go list` graph (`core->kernel`, `ops->core+kernel`, `session->core+ops+kernel`, domains->`core+kernel`), no cycles, no `ops.Identify`/promotion/bypass/`float`/parser/net APIs, no mutable package state.

## 4. Acceptance checklist (final gate)

`go build ./...` + `go vet ./...` + `go test ./...` clean **plus**: 7 packages, 39-file tree, zero third-party deps, A–S names all present and semantic, `TestApplyPositionalInputs`, REQ-032-01..22+10a, MRC-001..008, determinism double-run, all mutants killed, clean-worktree rerun.

## 5. Risks / non-goals

- No new packages/files, no CAS/parser/numerics/orchestration/physvet/promotion, no RFC 8785 substitution, no hardcoded E=mc2/Lorentz outputs, no test-weakening.
- `specs_v2_3.md` / `plan10_v2_3.md` untouched.
- Remaining plan-level items (P3-3 test import, DTO `schema_version` value `"1"`) documented as conformance notes, not architecture drivers.

## 6. Delivery report skeleton (for implementation agent)

1. Root causes 2. Changes (file+function+purpose) 3. Test changes (which test proves which behavior) 4. Mutation table 5. Spec coverage (REQ/MRC/A–S) 6. Architecture confirmation 7. Remaining gaps (spec vs test-quality vs plan-conformance vs unverified) 8. Verdict PASS/PASS WITH FINDINGS/FAIL.

*Status: IMPLEMENTED. All stages executed 2026-10-01; `go build ./...`, `go vet ./...`, `go test ./... -count=1` green (7 packages); 23/23 mutation probes KILLED in disposable copy; 39-file tree and package graph unchanged; `specs_v2_3.md` / `plan10_v2_3.md` untouched.*

## 7. Implementation record (build agent, 2026-10-01)

- `session/session.go` (rewritten): exact v2.3 `DraftMetadata`; `reconstructIdentification` (Compare-kind rules, merged metadata, TrimSpace, contamination); 17-stage `Session.Validate` over leaf helpers; `checkContainment` (C4 replay consequence); `Seal` via `buildTrustedCandidate` + `candidate.Validate`; conclusion record; `schema_version:"1"` session DTO.
- `session/ledger.go` (rewritten): C1 byte substrate (`InputCanonicals`/`OutputCanonical` as retained canonical `json.RawMessage`, hashes computed from retained bytes); `Step` accessors decode value copies; strict typed `ParseLedgerJSON`; structural `Ledger.Validate`; leaf primitives (C2); `schema_version:"1"` ledger DTO.
- `session/research_candidate.go` (rewritten): `Prediction.Relation` / `FalsificationCondition.ContradictingCondition` / `RecoveryClaim.Condition` → `core.Expr` with canonical wire codecs; single `buildTrustedCandidate`; full §26.8 `Validate` (incl. replay-derived containment); typed `candidateWire` DTO; `UnverifiedResearchCandidate.Validate` decodes wire → shared ctor → validates; strict outer-schema parse.
- `relativity/relations.go`: exact assumption keys; `ZeroThreeMomentum` carries `RestFrameAssumption()`.
- `relativity/manifest.json`: +4 framework assumptions (6 total); anomaly `no_gravity`/`scope_limit` per plan pin.
- `internal/kernel/mint.go`: `MintObjectMust` eliminated from production (test callers use local helper / explicit handling).
- Tests: `session/session_test.go` rewritten (unsafe fault injection on the session's own ledger, per-stage message pins, Q/R/S/K real); `hypothesis/candidate_test.go` rewritten (L/M/N real + AST forbidden-API scan + `errors.As` containment probe); `TestTypedMechanicsConstructors` + `TestDifferentiateKineticEnergyOps` (A/G); `TestEnergyMomentumRelation` + `TestMassEnergyDerivation` + real `TestLorentzFactorLimit` (H/I/J); `TestManifestEmbedPinned` (both manifests).
- Mutation matrix: 23/23 KILLED (M1-M8 stage isolation incl. message pins, M10-M12 Identify, M13 Seal-validate via params-tamper Seal test, M14 metadata, M15 unverified, M16 containment, M17/M18 relativity, M19 embed, M20/M21 hardcodes via merged-assumption asserts, M22 assertion replay, M23 conclusion).
- F-009: verified `errors.As(err, &CandidateContainmentError{})` directly (no `Unwrap` added).
- Known plan-level notes: `session_test.go` still imports `mechanics` (P3-3, test-only, unchanged); pre-existing `gofmt` drift in untouched files left as-is.


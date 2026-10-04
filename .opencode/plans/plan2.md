# Physics Compiler MVP — Adversarial Audit Remediation Plan (Plan 2)

**Generated from Phase 22 audit report** | Spec: `plan10/specs_v2_3.md` | Plan: `plan10/plan10_v2_3.md`  
**Repo**: `/home/chaschel/Documents/go/phys` | **Sandbox**: `/tmp/opencode/phys-audit`

---

## Executive Summary

The MVP implements the happy path for kernel minting, expression canonicalization, 12 ops, and basic session drafting — but **critical spec mandates are absent or incorrect**: ledger replay (§16.19), ResearchCandidate handoff (§26/§16.3), manifest `go:embed` (§33-O), Object JSON codec, typed errors in session, and 33% of normative §11 tests missing. Test quality is critically weak (8 zero-assertion tests covering entire L/M/N/Q/S acceptance surface).

**Baseline `go test ./...` FAILS** in delivered tree (`core.TestRepositoryTreeExact` vs `.opencode/plans/plan.md`).

---

## Classification Summary

| Classification | Count | Key Examples |
|---|---|---|
| **IMPLEMENTED+TESTED** | 11 (REQ-032) + kernel/expr/ops happy paths | REQ-032-01..07,15..20; MintObject; canonical sort; 12 ops |
| **IMPLEMENTED but INSUFFICIENTLY TESTED** | 7 | `session` import independence; error taxonomy 4 codes; hypothesis forced status; kernel dead exports; `MintObjectMust` prod presence |
| **PARTIALLY IMPLEMENTED** | 9 | Session drafting (8/12 sigs); ResearchCandidate types (surface wrong); Hypothesis containment; §24.3 test; `NewCandidateConcept` sig; session state machine; relativity assumption keys; manifest load (no embed); canonical round-trip |
| **NOT IMPLEMENTED** | 6 | Session Validate (14/17 stages); Session Seal→Candidate; Ledger/Step exported surface; Ledger tamper detection (4 forms); `go:embed` manifests; Object JSON codec; §11 tests S/Q/R/P/K |
| **IMPLEMENTED INCORRECTLY** | 8 | `Seal()` signature; `Conclude` hash vs Object; `Identify` no trim; `CandidateContainmentError` never thrown; `ParseResearchCandidateJSON` no validation; `DraftMetadata` missing fields; `ResearchCandidate` surface; README claims |
| **IMPLEMENTED BUT ARCHITECTURALLY NON-CONFORMING** | 4 | `MintObjectMust` in prod; `Session.Validate` returns `ledger` not `error`; test package assignments; unexported `Ledger`/`Step` |
| **UNVERIFIED** | 1 | `PH-` labels (none exist) |

---

## Remediation Priority List (Minimal to Spec Conformance)

### P0 — Critical Blockers (Must Fix First)

| # | Action | Spec Anchor | Est. Effort |
|---|---|---|---|
| **P0-1** | Implement `Session.Validate` full 17-stage replay (recompute `CurrentStepHash`, verify `InputHash`/`OutputHash`, decode `InputCanonicals`/`OutputCanonical`, call `ops.Apply` for transformations/identifications, return `LedgerValidationError`) | §16.19 L2754-2781, REQ-032-11..14 | High |
| **P0-2** | Add `MarshalJSON`/`UnmarshalJSON` to `core.Object` (canonical RFC 8785) — enables ledger round-trip | §10 L851, §16.12 L2720, §26.2 | Medium |
| **P0-3** | Make `Session.Seal()` return `(ResearchCandidate, error)`; mint candidate from sealed ledger + `DraftMetadata` candidate fields | §16.3 L2427-2455, §26.6 L3896, §33-S L4356 | High |
| **P0-4** | Fix `UnverifiedResearchCandidate.Validate()` — unmarshal into exported DTO, populate unexported fields, run shared validation ctor with `session.Validate` | §26.10 L3972, L3982, REQ-032-10a L4236 | High |
| **P0-5** | Add `go:embed` for `mechanics/manifest.json` and `relativity/manifest.json`; load via embed | §33-O L2825, L4340 | Low |

### P1 — Signature/Structure Alignment

| # | Action | Spec Anchor | Est. Effort |
|---|---|---|---|
| **P1-1** | Align `Session` 12 signatures to §16.3 table: add `New()`, add `label` params, `Conclude(core.Object)`, `Validate() error`, `Seal() (ResearchCandidate,error)` | §16.3 L2427-2455 | Medium |
| **P1-2** | Extend `DraftMetadata` with 8 missing candidate fields; copy slices at `Draft` time | §16.2.1 L2404-2423 | Medium |
| **P1-3** | Fix `ResearchCandidate` surface: `ID()`, `Premises() []core.Object`, `Derivation() Ledger` (exported), value returns | §26.10 L3951-3965 | Medium |
| **P1-4** | Add `ParseResearchCandidateJSON` outer-schema validation; store `canonicalBytes` | §26.10 L3980, REQ-032-10a | Low |
| **P1-5** | Implement §26.8 stage 2 (canonicalization) and stage 4 (derivation replay via `session.Validate`) | §26.8 L3910-3920 | Medium |

### P2 — Test Coverage & Quality

| # | Action | Spec Anchor | Est. Effort |
|---|---|---|---|
| **P2-1** | Implement 9 missing §11 tests (S, Q, R, P-round-trip, K-session, I, B, A, H) with real assertions | Plan §11 L424-473 | Medium |
| **P2-2** | Replace 8 zero-assertion tests with real assertions; fix 10 tautological tests | §33 L4326-4356 | Low |
| **P2-3** | Fix REQ-032-08: `TrimSpace` on justification; construct `IdentifyError` | §16.7.4 L2527, REQ-032-08 L4226 | Low |
| **P2-4** | Fix REQ-032-21/22: return typed `ProvenanceError`/`LedgerValidationError` on state violations | §16.1 L2390/2393, REQ-032-21/22 | Low |

### P3 — Housekeeping & Polish

| # | Action | Spec Anchor | Est. Effort |
|---|---|---|---|
| **P3-1** | Remove 73 extra files or update Plan §38 budget; fix `TestRepositoryTreeExact` allow-list | Plan §38 L4491-4504 | Low |
| **P3-2** | Fix `README.md` 5 false claims | — | Trivial |
| **P3-3** | Align test packages to Plan §9 pins (external `mechanics_test`, `relativity_test`, `hypothesis_test` with cross-imports) | Plan §9 L77-89 | Low |
| **P3-4** | Fix assumption keys in relativity: `rest_mass_nonnegative`, `speed_of_light_positive` | §11.6 L1434/1442, Plan §8 L345 | Low |
| **P3-5** | Adjudicate `MintObjectMust`: move to test-only or remove from production | §5.2.1 | Low |

---

## Mutation Probe Baseline (Must Preserve)

| Probe | Current Result | Must Remain |
|---|---|---|
| Kernel mint dim check | KILLED by `TestMintObjectRejectsInvalidSpec` | ✓ |
| Expr canonical sort | KILLED by `TestCanonicalOrderingStable` + 5 cascading | ✓ |
| Session ledger step hash | **SURVIVES** (no test exercises Validate) | Must be KILLED after P0-1 |
| Session Validate empty ledger | **SURVIVES** | Must be KILLED after P0-1 |

---

## File/Tree Integrity Gates

- Pre/post SHA-256 digest (excl `.git`, `.opencode`) must remain IDENTICAL for all 39 plan-tree files
- `git status` must show no source modifications (only `.opencode/plans/plan2.md` added)
- Plan §3 tree (39 files) must remain 39/39 present
- `go list -deps` must show 7 local + 97 stdlib = 0 third-party modules
- `go build ./...`, `go vet ./...` must remain clean
- `go test ./...` must PASS (currently fails due to `.opencode/plans/plan.md`)

---

## Implementation Order

1. **P0-5** (`go:embed` manifests) — Low effort, unblocks acceptance O
2. **P0-2** (Object JSON codec) — Enables ledger canonical bytes
3. **P0-1** (Session.Validate 17-stage replay) — Core fix, kills mutation probes
4. **P0-3** (Seal returns ResearchCandidate) — Requires P0-1, P0-2, P1-2
5. **P0-4** (UnverifiedResearchCandidate.Validate fix) — Requires P0-3, P1-3
6. **P1-1** (Session signatures) — Mechanical but pervasive
7. **P1-2** (DraftMetadata fields) — Required for P0-3 candidate minting
8. **P1-3** (ResearchCandidate surface) — Required for P0-3/P0-4
9. **P1-4, P1-5** (Parse validation + pipeline stages)
10. **P2-1, P2-2, P2-3, P2-4** (Test coverage)
11. **P3-1..5** (Housekeeping)

---

## Notes

- All work must be done in the repo at `/home/chaschel/Documents/go/phys`
- Run `go build ./...`, `go vet ./...`, `go test ./...` in sandbox `/tmp/opencode/phys-audit` after each logical group
- Verify mutation probes still killed (or newly killed for session validation)
- Do NOT modify files outside the 39 plan-tree files unless updating Plan §38 budget
- The `.opencode/plans/plan.md` (pre-audit artifact) will cause `TestRepositoryTreeExact` to fail — either exclude `.opencode` in the test or remove it (but keep `plan2.md` here)
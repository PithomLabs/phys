# Physics Compiler MVP — Corrected Remediation Plan (Plan 3)

**Generated from Phase 22 audit + user architectural review**  
Spec: `plan10/specs_v2_3.md` | Plan: `plan10/plan10_v2_3.md`  
Repo: `/home/chaschel/Documents/go/phys` | Sandbox: `/tmp/opencode/phys-audit`

---

## Core Principle

The audit findings are correct; the remediation wording in Plan 2 had architectural errors. This plan corrects them.

> **Canonical artifacts, replay, candidate validation and manifest semantics are executable integrity mechanisms — not documentation decoration.** (v2.3 philosophy)

---

## P0 — Actual Semantic Blockers (Must Fix First)

| # | Action | Spec Anchor | Notes |
|---|---|---|---|
| **P0-1** | Complete `Session.Validate` 17-stage replay: state/genesis/StepID/index **→** recompute `CurrentStepHash` for every step **→** verify `InputHash` matches `HashObject(InputCanonical)` **→** verify `OutputHash` matches `HashObject(OutputCanonical)` **→** verify `AssumptionHash`/`ConventionHash` **→** **transformation steps**: replay via `ops.Apply` and compare result **→** **identification steps**: replay via **session-owned identification construction path** (not `ops.Apply`) **→** provenance/MRC validation **→** candidate containment validation **→** return `LedgerValidationError` on any mismatch | §16.19 L2754-2781, REQ-032-11..14 L4238-4244 | **Critical**: identification replay uses private session path, not `ops.Apply` (Identify firewall). Depends on P1-6. |
| **P0-2** | Fix internal kernel object canonical encode/decode: add/fix `EncodeObject`/`DecodeObject` (or equivalent) in `internal/kernel` with **invariant validation before replay admission** and **exact `Encode(Decode(b)) == b`** round-trip. **Do NOT add public `UnmarshalJSON`/`MarshalJSON` to `core.Object`** — that would expose public deserialization and undermine constructor authority. | §10 L851, §16.12 L2720, §26.2, REQ-010-07 | Ledger `InputCanonicals`/`OutputCanonical` must round-trip through internal codec. RFC 8785 is not the spec — v2.3 defines its own typed canonical JSON contract. |
| **P0-3** | Fix public `Ledger`/`Step` surface to exact v2.3 API: export `Ledger`, `Step`, `ParseLedgerJSON(data []byte) (Ledger, error)`, `Ledger.Steps() []Step`, `Ledger.DerivationHash() string`, `Ledger.Validate() error`, `Ledger.CanonicalJSON() ([]byte, error)`, `Step` with read-only accessors for every §16.12 field. | §16.21 L2796-2814, REQ-§16.21-MUST-01/02, Plan L719-720 | Currently `ledger`/`stepEnvelope` are unexported; `ParseLedgerJSON` returns unexported type. |
| **P0-4** | Fix `Session.Seal()` → `(ResearchCandidate, error)`: mint a `ResearchCandidate` from the sealed ledger + full `DraftMetadata` candidate fields. | §16.3 L2427-2455, §26.6 L3896, §33-S L4356 | Requires P0-1, P0-2, P1-2 (DraftMetadata). |
| **P0-5** | Fix `UnverifiedResearchCandidate.Validate()`: delegate to the **same trusted validation/sealing path** used by `Session.Seal()` — no second trusted-construction mechanism. Implementation choice (DTO, builder, shared private ctor) is free; the semantic contract is mandatory. | §26.10 L3972, L3982, REQ-032-10a L4236, Plan §13 L170, L415 | **Not** "unmarshal into exported DTO" — that's an implementation detail. The requirement: single validation path. |
| **P0-6** | Implement required `go:embed` manifest loading: `//go:embed manifest.json` in `mechanics/primitives.go` and `relativity/primitives.go`; load via embed, not `os.ReadFile`. | §33-O L2825, L4340, Acceptance O | Low effort, unblocks acceptance O. |

---

## P1 — Contract Alignment (Exact Spec Signatures)

| # | Action | Spec Anchor |
|---|---|---|
| **P1-1** | Align `Session` 12 signatures to §16.3 table: add `New() *Session`, add `label` param to `Postulate`/`Declare`/`Define`/`Step`, `Conclude(core.Object)`, `Validate() error`, `Seal() (ResearchCandidate, error)`. | §16.3 L2427-2455 |
| **P1-2** | Extend `DraftMetadata` with 8 missing candidate fields; copy slices at `Draft` time. | §16.2.1 L2404-2423 |
| **P1-3** | Fix `ResearchCandidate` accessors: `ID()`, `Premises() []core.Object`, `Derivation() Ledger` (exported), value returns (not pointer). | §26.10 L3951-3965, Plan L416 |
| **P1-4** | `ParseResearchCandidateJSON`: add outer-schema validation; store `canonicalBytes`. | §26.10 L3980, REQ-032-10a |
| **P1-5** | Implement §26.8 candidate validation pipeline stage 2 (canonicalization) and stage 4 (derivation replay via `session.Validate`). | §26.8 L3910-3920 |
| **P1-6** | Add session-owned identification replay helper (private) used by `Validate` for identification steps — distinct from `ops.Apply`. | §16.19 identification replay clause |

---

## P2 — Proof Quality (Mutation-Killing Tests)

| # | Action | Spec Anchor |
|---|---|---|
| **P2-1** | Implement 9 missing §11 acceptance tests with **substantive assertions**: `TestSealResearchCandidate` (S), `TestLedgerTamperDetection` (Q), `TestDerivationDeterminism` (R), `TestExactRationalRoundTrip` (P), `TestSessionIdentifyRecords` (K), `TestMassEnergyDerivation` (I), `TestNewtonSecondLawConstruct` (B), `TestTypedMechanicsConstructors` (A), `TestEnergyMomentumRelation` (H). `TestObjectCanonicalRoundTrip` required as additional object-level coverage. | Plan §11 L424-473 |
| **P2-2** | Replace 8 zero-assertion tests with real assertions; fix 10 tautological tests (see audit Phase B). | §33 L4326-4356 |
| **P2-3** | Fix REQ-032-08: `TrimSpace` on justification; construct `IdentifyError`. | §16.7.4 L2527, REQ-032-08 L4226 |
| **P2-4** | Fix REQ-032-21/22: return typed `ProvenanceError`/`LedgerValidationError` on state violations. | §16.1 L2390/2393, REQ-032-21/22 |
| **P2-5** | Ensure mutation probes are **killed**: session ledger step-hash mutation, empty-ledger mutation must fail after P0-1. | Audit mutation results |

---

## P3 — Plan-Level Conformance & Cleanup

| # | Action | Spec Anchor | Classification |
|---|---|---|---|
| **P3-1** | **Repository baseline reconciliation**: distinguish (a) MVP implementation tree (39 files per §3), (b) pre-existing project/audit/planning artifacts (`.opencode/`, `plans/`, `plan10/`, reality docs, `.gitignore`, `.gut`), (c) genuine prohibited implementation files. Remove only (c). **`TestRepositoryTreeExact` MUST continue to assert the exact 39-file MVP implementation tree** — do not weaken its allow-list. Pre-existing audit/planning artifacts outside that tree must be classified separately and MUST NOT be admitted as implementation-tree files. | Plan §38 L4491-4504 | Cleanup — **do not blindly delete 73 files; do not weaken the scope gate** |
| **P3-2** | Fix `README.md` 5 false claims (audit #18). | — | Trivial |
| **P3-3** | Align test packages to Plan §9 pins where required: external `mechanics_test` (uses `ops`), external `relativity_test`, external `hypothesis_test` (uses `ops`, `session`, `relativity`). | Plan §9 L77-89 | Plan-level test-package conformance |
| **P3-4** | Fix assumption keys in relativity primitives: `rest_mass_nonnegative`, `speed_of_light_positive` (not `m_nonnegative`/`c_positive`). | §11.6 L1434/1442, Plan §8 L345, §11 H L460 |
| **P3-5** | Adjudicate `MintObjectMust`: move to test-only or remove from production (per §5.2.1 only `MintObject` is specified). | §5.2.1 | Adjudication |
| **P3-6** | Manifest canonical-byte pin: ensure `json.Marshal(manifest)` matches `manifest.json` bytes exactly (mechanics currently indented). | Plan §12 matrix | **Plan 10 conformance pin** |

---

## Mutation Probe Baseline (Must Preserve & Extend)

| Probe | Current | After P0-1 Must Be |
|---|---|---|
| Kernel mint dim check | **KILLED** by `TestMintObjectRejectsInvalidSpec` | **KILLED** |
| Expr canonical sort | **KILLED** by `TestCanonicalOrderingStable` + 5 cascading | **KILLED** |
| Session ledger step hash | **SURVIVES** (no test exercises Validate) | **KILLED** by new `TestLedgerTamperDetection` / replay tests |
| Session Validate empty ledger | **SURVIVES** | **KILLED** |

---

## File/Tree Integrity Gates

- Pre/post SHA-256 digest (excl `.git`, `.opencode`) identical for all 39 plan-tree files
- `git status` shows no source modifications (only `.opencode/plans/plan3.md` added)
- Plan §3 tree: 39/39 present
- `go list -deps`: 7 local + 97 stdlib = 0 third-party
- `go build ./...`, `go vet ./...` clean
- `go test ./...` **PASS** (fix `TestRepositoryTreeExact` baseline per P3-1)

---

## Implementation Order (Corrected)

1. **P0-6** (`go:embed` manifests) — Low, unblocks acceptance O
2. **P0-2** (Internal kernel object canonical encode/decode with invariant validation) — Enables ledger round-trip
3. **P0-3** (Export `Ledger`/`Step` surface) — Required for P0-1 replay decode
4. **P1-6** (Session-owned identification replay helper) — Required for P0-1 identification replay
5. **P0-1** (Session.Validate 17-stage replay with correct identification path) — Core fix, kills mutation probes
6. **P1-2** (DraftMetadata candidate fields) — Required for P0-4
7. **P1-1** (Session signatures) — Mechanical but pervasive
8. **P0-4** (Seal returns ResearchCandidate) — Requires P0-1, P0-2, P1-2
9. **P1-3** (ResearchCandidate accessors) — Required for P0-4/P0-5
10. **P0-5** (UnverifiedResearchCandidate.Validate single validation path) — Requires P0-4, P1-3
11. **P1-4, P1-5** (Candidate validation pipeline)
12. **P2-1..P2-5** (Test coverage + mutation kills)
13. **P3-1..P3-6** (Baseline reconciliation, README, test packages, assumption keys, MintObjectMust, manifest bytes)

---

## Notes

- All work in `/home/chaschel/Documents/go/phys`
- Run `go build ./...`, `go vet ./...`, `go test ./...` in sandbox `/tmp/opencode/phys-audit` after each logical group
- Verify mutation probes: old kills preserved, new session validation probes killed
- Do NOT add public `UnmarshalJSON`/`MarshalJSON` to `core.Object`
- Do NOT use `ops.Apply` for identification replay in `Validate`
- `TestRepositoryTreeExact` must remain strict on the 39-file MVP tree (P3-1)
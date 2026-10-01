# Physics Compiler MVP — `github.com/PithomLabs/phys`

Module: `github.com/PithomLabs/phys`  
Go: `1.24` (stdlib only — no third-party dependencies)

## Scope (MVP — specs_v2_3.md / Plan 10 v2.3)

This is a deterministic symbolic physics compiler. It defines:
- Immutable `core.Object` carriers (§5) with fixed provenance (`POSTULATED`, `DEFINED`, `DERIVED`, `IDENTIFIED`, `HYPOTHESIS`, `APPROXIMATED` (reserved; no MVP operation produces this status — REQ-013-01)) and corpus status (`NONE`, `ESTABLISHED`, `CONTESTED`, `SUPERSEDED`, `FALSIFIED`).
- A closed expression set (`core.Expr`) with 10 node kinds (§8.1/§8.2) and deterministic canonical JSON / SHA-256 (§9).
- Twelve pure operations (`ops.Apply` — §15) evaluated at exact positional indices (§15.13.1): `add`, `subtract`, `multiply`, `divide`, `pow`, `simplify`, `substitute`, `differentiate`, `limit`, `compare`, `solve`, `select_branch`. No plugin registry.
- A session authority (`session.Session` — §16) that is the only mechanism to commit a ledger, produce an `IDENTIFIED` status (`Session.Identify` — MRC-006), or seal a `ResearchCandidate` (MRC-007 / §26).
- Candidate artifacts (`ResearchCandidate`) with falsifiability structures (`Prediction`, `FalsificationCondition`, `RecoveryClaim`, `AnomalyReference`) and review/challenge consumption (`core.Challenge` / `core.Review` — §27). No promotion APIs exist.

## Non-goals (explicitly deferred — REQ-002-01 / specs_v2_3.md §2 / §36)

No `.phys` language, no lexer/parser/frontend (`go/parser`, `go/scanner`, `go/token` banned), no custom compiler language, no truth-score/promotion APIs (`TRUST`, `PROMOTE`, `APPROVE`), no cryptographic signing/authenticity mechanism, no adversarial-source protection claim, no deferred `physvet` (v0.5 contract recorded in specs only).

## Authority boundaries

**constructor authority (MRC-001 / REQ-000-03 / REQ-000-04):** `internal/kernel.Object` fields are unexported; `MintObject` validates the 7-point contract (§5.2.1) before visibility. `core.Object` is an alias; `core` exposes no mint forwarder. This is a **public/package-boundary** mechanism (`go` `internal/` visibility), not cryptographic or hostile-source protection. The README explicitly states: `"not cryptographic"`.

**Operation authority (MRC-002..005 / MRC-008):** `ops` functions are pure; input `HYPOTHESIS` propagates to `HYPOTHESIS` output (contamination law — §13.2). `Simplify` never yields `IDENTIFIED`; `Identify` is session-only.

**Session authority (MRC-006 / MRC-007):** Only `session.Session` commits a ledger, records `Identification`, or produces `IDENTIFIED`. `ParseLedgerJSON` returns an unverified `Ledger` for validation only (§16.20); `ParseResearchCandidateJSON` returns only `UnverifiedResearchCandidate` (§26.10) — no `core.Object` accessors.

**Candidate containment (MRC-008 / §26.8 / §26.9):** Any hypothesis-dependent derivation step presented with a trusted provenance status, or with `CorpusStatus` ≠ `NONE`, or a `Hypothesis` field not `HYPOTHESIS`, fails validation (`CandidateContainmentError`). No promotion API exists (`Promote`, `Trust`, `ApproveHypothesis`, `BypassMRC`, `OverrideMRC` — all absent).

## Integrity vs authenticity (§16.17 / §26.10)

The SHA-256 step chain (`PreviousStepHash` → `CurrentStepHash`) detects mutation (§16.19–§16.20). It does **not** prove authenticity: anyone rewriting the full ledger can recompute hashes. Signing/authenticity is outside MVP (specs_v2_3.md §16.17 / §26.10). `physvet` deferred — REQ-002-15 (§36). `ResearchCandidate.Validate()` asserts structural consistency, not physical truth.

## Dependency rules (§4 / REQ-004-01 / REQ-004-02)

Architecture (locked):
```
internal/kernel  ← Object/Expr/Kind/Dimension/metadata authority
       ↑
core             ← public API aliases + corpus/manifests (§26 / §27)
       ↑
ops              ← pure transformations (12 ops; no session/ledger)
       ↑
session          ← ledger/replay/seal authority (imports core, ops, kernel)

mechanics / relativity  → core + kernel (never session)
hypothesis             → core + kernel (never mechanics/relativity at build time)
```
`core` imports only `internal/kernel` + stdlib. `ops` never imports `session`. `session` never imports `mechanics`/`relativity`/`hypothesis`. Import independence verified by `go build ./...` (cycle-free).

## MRC references

- `mrc-v0.4` — pinned for MVP artifacts (§5.4 / §10.4 / §26.2 / §16.12).
- MRC-001: constructor integrity (§5.2.1).
- MRC-002..005: operation-time constraints (§15 / §6.2 / §8 / §9).
- MRC-006: identification firewall (`Identify` only via `Session.Identify` — §14.6 / §16.7).
- MRC-007: session authority (ledger commit / seal — §16 / §26).
- MRC-008: candidate containment (§26.8 / §26.9 / §26.10).

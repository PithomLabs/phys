# AUDIT PLAN — Physics Compiler MVP vs `specs_v2_3.md` (v2.3)

Adversarial, independent, read-only audit of the completed Physics Compiler MVP.

Authoritative sources: `plan10/specs_v2_3.md` (sole normative authority), `plan10/plan10_v2_3.md` (implementation plan; its §11 test names and §12 matrix are normative only where the plan explicitly says so). Plan claims of "PASS" are not evidence. `go test ./...` green is necessary but not sufficient.

> Location note: written to `.opencode/plans/plan.md` because repository edits are permission-locked to this path, and a root-level `plan.md` would make `core.TestRepositoryTreeExact` fail, contaminating the audit baseline. The audit report must record this file as a post-hoc audit artifact.

---

## 0. Ground rules (apply to every phase)

- **Original repo is never touched by the audit itself.** All `go build/vet/test/list` runs happen in a sandbox copy at `/tmp/opencode/phys-audit` (rsync without `.git`). Final gate: `git status --porcelain` on the original must be empty (aside from this plan file, created at user request before the audit).
- **Evidence protocol:** a requirement is only PASS with (a) spec clause + line, (b) implementation `file:line`/function, (c) test name + test-body line proving the behavior, (d) for high-risk items, an adversarial probe result. Plan "PASS" claims, function names, and `go test` green are never evidence on their own. Conflicting evidence → report the conflict. Unverifiable → `UNVERIFIED`.
- **Classification per requirement:** `IMPLEMENTED+TESTED` / `IMPLEMENTED but INSUFFICIENTLY TESTED` / `PARTIALLY IMPLEMENTED` / `NOT IMPLEMENTED` / `IMPLEMENTED INCORRECTLY` / `IMPLEMENTED BUT ARCHITECTURALLY NON-CONFORMING`.
- **Delegation:** parallel `explore` subagents for mechanical scans (REQ extraction, import graph, AST/export scans, test-body quality review), each with the evidence protocol; **I personally verify every CRITICAL/HIGH finding and every acceptance A–S verdict.** Subagents may not write files.

## 1. Phase 1 — Baseline (sandbox)

```text
go version
go build ./...      # in /tmp/opencode/phys-audit
go vet ./...
go test ./... -count=1 -v > test_full.log   (record pass/fail/skip counts per package)
go list ./...
go list -deps ./...
```

Inventory: `find` (excluding `.git`) vs plan10 §2's 39-file tree; classify every file as spec-listed / pre-existing planning doc (`plans/`, `plan10/`, `reality*.md`, `.gitignore`, `.gut`) / **extra**. Report both counts (total repo files vs 39-file MVP set) and adjudicate against spec §3 ("create or modify only") + §38. Count: Go source / test / JSON / config-doc; 7 package dirs (`internal/kernel`, `core`, `ops`, `session`, `mechanics`, `relativity`, `hypothesis`) vs total repo dir count. Verify module path, `go 1.24`, zero third-party deps.

## 2. Phase 2 — Mechanical requirement extraction

Scripted (read-only) extraction from `specs_v2_3.md`:

1. All `REQ-*` labels → unique set (verify plan's claimed **115**).
2. All lines containing `MUST`/`MUST NOT` → group by section (verify claimed **126** MUST-bearing sections); each gets a deterministic `REQ-§N-MUST-nn` ID where unlabeled.
3. MRC-001..008, acceptance A–S, §15.13.1 table, REQ-032-01..22+10a.
4. Diff against plan §12 rows **and** against the actual code/test inventory; produce the master coverage table (requirement → impl location → test → verdict). Plan §12 claims are inputs, not conclusions.

## 3. Phases 3–9 — Correctness deep dives (targeted reads)

| Phase | Primary files | Critical checks |
|---|---|---|
| 3 Kernel/mint | `internal/kernel/mint.go` (262), `types.go` (2470), `core/object.go` | §5.2.1 7-point contract each; zero-value invalid; defensive copies; single mint path; **`MintObjectMust` at mint.go:93 — is this a "second object-minting helper" (§5.2.1 last line)?**; decode path runs same invariants |
| 4 Kinds/dims/MRC | `kernel/types.go`, `ops/arithmetic.go`, `ops/relation.go` | 18 kinds + ordinals; §6.2 matrix cell-by-cell; MRC evaluation order 001→002→003→004→005→008; Mass vs RestMass fixture; exact typed errors (errors.As, not "some error") |
| 5 Expr/canonicalization | `kernel/types.go`, `ops/simplify.go` (441) | 10 nodes/ordinals, no callbacks, constructor-vs-Simplify boundary, child 3-key sort, Pow(x,1)/Pow(x,0)-gated/0^0/0^-1/negative+repeated powers, Sqrt rules + sign entailment, `base≠0` vs `base>0`, relation/BranchSet simplify never IDENTIFIED, defensive copies, canonical bytes |
| 6 Operations/positional | `ops/*.go`, `session/session.go:145`, replay in `Validate` | exact 12 signatures; §15.13.1 positions enforced **identically** in `ops.Apply`, `Session.Step`, replay; `OperationParams` field order/canonical JSON/Kind binding; inspect `TestApplyPositionalInputs` body (prelim: substantive — all 12 + reorder + identify-exclusion, ops/operations_test.go:191–376); no `ops.Identify`, no package-level mutable state |
| 7 Calculus | `ops/transform.go` (403), `ops/relation.go` | Differentiate bounded rule set only; Limit executes fixed body (source trace — no function-ID shortcut); Solve exact pattern; SelectBranch 6-step incl. `selected_branch/<hash>` |
| 8 Mechanics | `mechanics/primitives.go`, `relations.go`, `manifest.json` | per-constructor expr/dim/provenance/assumptions inspected, not just existence; Test D fixture Mass vs RestMass everywhere |
| 9 Relativity | `relativity/*.go`, both manifests | 10-item manifest exactly, Velocity wrapper, RestFrame assumption; §20 chain traced in `TestFullDerivationSequence` (prelim: executes chain, but only *logs* golden states — assert quality); §21, §22 traces; no hardcoded E=mc² |

## 4. Phases 10–12 — Canonical JSON, session/replay, tamper

- Canonical encoders for all artifact classes (Expr, Dimension, AssumptionSet, ConventionSet, Object, Manifest, OperationParams, Step, StepEnvelope, Ledger, ResearchCandidate): field order, rational `"1/2"`, set ordering, lowercase hex, `schema_version:"1"` only on DTO, no `map[string]any`, no map iteration.
- **Walk the 17-stage §16.19 pipeline in `Session.Validate` (session.go:244+)** — verify each stage exists, in order, that replay uses retained `InputCanonicals`/`OutputCanonical`/`ParamsCanonical` (never infers from hashes), and Identify replay uses a pure session-owned helper.
- Preliminary hot spot: `session/ledger.go` is only **30 lines**; Step/StepEnvelope/chain logic appears to live in `session.go` — verify §16.12 16-field order, genesis, StepID, envelope formula, and §16.21 exported `Ledger` surface (actual `ParseLedgerJSON` returns **unexported** `ledger`).
- Determinism: re-run the §20 derivation twice + fresh session in a test probe script (in sandbox) and byte-compare all canonical artifacts.
- **Tamper probes on `/tmp` copy:** 4 §16.20 modes; additionally mutate the *test* expectations / weaken replay in the copy to prove the tests would fail (Phase 12 + 20 evidence).

## 5. Phases 13–17 — Hypothesis, manifests, errors, source/AST, import graph

- `hypothesis/candidate.go` (50 lines): forced HYPOTHESIS/NONE, no caller provenance, no promotion; ResearchCandidate: 14 accessors + Validate + CanonicalJSON (prelim: accessors exist, research_candidate.go:85–125), private seal path shared with `UnverifiedResearchCandidate.Validate`, §26.8 stage order, containment → `CandidateContainmentError`.
- Manifests: manual item-by-item comparison of both JSON files vs constructors (11/10 sets prelim verified by constructor IDs); strict decode/unknown-field rejection; canonical bytes (**prelim: `CanonicalManifestJSON` = `json.Marshal` compact, but `mechanics/manifest.json` is indented → file-bytes ≠ canonical bytes; neither test checks file-bytes == canonical**); **no `go:embed` anywhere (grep verified) despite spec §33-O requiring it** — tests use `os.ReadFile`.
- 11 typed errors + errors.As + deterministic diagnostics + no-panic.
- Repo-wide AST scans (test-assisted `TestForbiddenSourceSurface`, plus independent scans): floats, `math.*`, parser/scanner/token, rand/http/exec, `map[string]any`, exported `NewObject`, promotion/truth/bypass names, `ops.Identify`, Lorentz-ID shortcut, package-level mutables, `review/` package, deferred packages.
- Import graph from `go list -deps` + per-file import grep (prelim: edges look conformant — kernel→∅, core→kernel, ops→core+kernel, session→core+ops+kernel, domains→core+kernel; confirm `Challenge`/`Review` only in `core/corpus.go`, no cycles).

## 6. Phases 18–20 — Acceptance, negative suite, test quality

- For each A–S + `TestApplyPositionalInputs`: existence under **exact plan name/package**, then body inspection.
- **Preliminary test-inventory verdict (verify in execution):** 115 distinct test funcs exist, but **~50 plan-normative names are absent**, including acceptance-critical `TestTypedMechanicsConstructors` (A), `TestEnergyMomentumRelation` (H), `TestMassEnergyDerivation` (I — replaced by `TestFullDerivationSequence`), `TestSessionIdentifyRecords` (K), `TestLedgerTamperDetection` (Q), `TestDerivationDeterminism` (R — renamed `...Basic`), `TestSealResearchCandidate` (S), `TestNewtonSecondLawConstruct` (B). Plan §11 declares names normative → coverage finding regardless of behavior; then judge behavior separately (e.g. is J actually executed? prelim: `relativity.TestLorentzFactorLimit` is a `t.Log` stub; check `ops.TestLorentzFactorLimit:915`).
- Negative suite REQ-032-01..22+10a mapped to tests; **preliminary hot spot: `session/session_test.go` (129 lines) and `hypothesis/candidate_test.go` (44 lines) contain assertion-free stubs** (e.g. `TestCandidateContainmentBasic` body is `_ = true`; Q-section only parses `{"steps":[]}` and logs) → REQ-032-09/10a/11–14/21/22, K/Q/R/S likely `IMPLEMENTED but INSUFFICIENTLY TESTED` or `NOT IMPLEMENTED` pending code checks.
- Quality audit: non-nil-only asserts, no dimension/provenance checks, `t.Logf`-only tests, hardcoded-answer tests, internal-vs-external package mismatches (prelim: plan pins `mechanics_test`/`relativity_test`, actual both internal `mechanics`/`relativity`), mutation probes (break positional binding / replay retention in the sandbox → does any test fail?).

## 7. Phase 21 — Manual traces (code-walk, not test-based)

10 traces as specified (Mass+RestMass rejection; KE differentiation; NewtonSecondLaw+manifest; E=mc² §20; Lorentz limit; Step→Commit→replay; 3 tamper modes; contamination; containment; Seal→candidate→Parse→Validate), each reported `INPUT → VALIDATION → TRANSFORMATION → METADATA → OUTPUT → HASH/REPLAY EFFECT` with file:line. Specific check: **`Seal() error` (session.go:228) vs spec §16.3 `Seal() (ResearchCandidate, error)` — if no accessor returns the sealed candidate, trace 10 cannot complete as specified.**

## 8. Phase 22 — Report (terminal only)

The 11-section report exactly as demanded: EXECUTIVE VERDICT (coverage/correctness/architecture/overall) → baseline → coverage counts → correctness findings (severity/REQ/file/impl-vs-spec/why/test-catches?) → test-quality findings (separate) → architecture/API → determinism → session/replay → physics derivations → **exact remediation list (ordered; each tagged implementation-change / test-change / docs-only / no-action)** → final status. No report file written; remediations are "smallest conforming correction", no redesign.

## 9. Preliminary heat map (recon hypotheses — to confirm, not conclusions)

**Likely FAIL/PARTIAL areas (priority order):**

1. Acceptance test-name inventory (~50 missing normative names; A/H/I/K/Q/R/S affected).
2. Session & hypothesis test substance (stub bodies; Q/R/S/K/L/M/N likely unproven).
3. `Seal`/`Validate`/`Draft`/`Postulate`… signature deviations from spec §16.3 (spec says "semantic signatures … are fixed").
4. Missing `go:embed` (spec §33-O) + manifest file bytes not canonical (plan §7).
5. `kernel.MintObjectMust` as candidate second mint helper (§5.2.1).
6. Exported `Ledger`/`Step` surface: unexported `ledger` return types vs §16.21.
7. `TestRepositoryTreeExact` tolerates extra pre-existing files and never asserts the count 39.
8. Manifest `domain` counts (mechanics test asserts 9 with message "expected 10"; plan pins `domain[1]`) vs spec §17.2/§18.3.

**Preliminary clean areas (to re-verify):** module identity, stdlib-only, import-graph direction, no floats/forbidden imports in non-test sources, 10-item relativity/11-item mechanics constructor sets, README integrity-vs-authenticity + constructor-authority statements, docs §35 sections, `TestApplyPositionalInputs` substance, §20 chain actually executed in `TestFullDerivationSequence`.

## 10. Execution order & integrity gate

`1 baseline → 2 extraction → (3–9 parallel subagent sweeps, I verify) → (10–12 own deep read + probes) → (13–17 mixed) → (18–20 own) → 21 traces → 22 report`; probes only in `/tmp/opencode/phys-audit*`; final `git status --porcelain` check on the original (only `.opencode/plans/plan.md` may appear, as a pre-audit artifact).

---

Scope decisions:

- Pre-existing `plans/`, `plan10/`, `reality*.md` are **reported as a counted discrepancy with both interpretations** (total repo files vs implementation-created files), adjudicated on spec §3/§38 wording rather than silently excluded.
- All `go` commands run in the sandbox copy, so the original is provably untouched by audit execution.
- This plan file itself must be disclosed in the Phase 1 inventory as an audit artifact (permission-locked location; not part of the 39-file MVP tree).

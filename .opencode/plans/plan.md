# AUDIT PLAN — Physics Compiler MVP vs `specs_v2_3.md` (v2.3) — REVISED (rev 2)

Adversarial, independent, read-only audit of the completed Physics Compiler MVP.

Authoritative sources: `plan10/specs_v2_3.md` (sole normative authority), `plan10/plan10_v2_3.md` (implementation plan; its §11 test names and §12 matrix are normative only where the plan explicitly says so). Plan claims of "PASS" are not evidence. `go test ./...` green is necessary but not sufficient.

> Location note: written to `.opencode/plans/plan.md` because repository edits are permission-locked to this path, and a root-level `plan.md` would make `core.TestRepositoryTreeExact` fail. The audit report must disclose this file as a pre-audit artifact.

---

## 0. Ground rules (apply to every phase)

- **Original repo is never touched by the audit itself.** All `go build/vet/test/list` runs happen in a sandbox copy at `/tmp/opencode/phys-audit` (rsync without `.git`, without `.opencode`). Integrity gate: SHA-256 digest of the full file tree (excluding `.git`) captured before and after execution must be identical; `git status --porcelain` output before vs after must be identical (both legitimately show the untracked implementation — the initial commit contains only `plan9.md`).
- **Evidence protocol:** a requirement is only PASS with (a) spec clause + line, (b) implementation `file:line`/function, (c) test name + test-body line proving the behavior, (d) for high-risk items, an adversarial probe result. Plan "PASS" claims, function names, and `go test` green are never evidence on their own. Conflicting evidence → report the conflict. Unverifiable from available evidence → `UNVERIFIED`.
- **Classification per requirement (formal, closed list):** `IMPLEMENTED+TESTED` / `IMPLEMENTED but INSUFFICIENTLY TESTED` / `PARTIALLY IMPLEMENTED` / `NOT IMPLEMENTED` / `IMPLEMENTED INCORRECTLY` / `IMPLEMENTED BUT ARCHITECTURALLY NON-CONFORMING` / `UNVERIFIED`.
- **Every finding reports four statuses:** `SPEC STATUS` / `PLAN STATUS` / `IMPLEMENTATION STATUS` / `TEST STATUS`. A Plan 10 pin is never reported as a spec violation. Template:
  - `go:embed` → Spec: FAIL / Plan: FAIL / Impl: FAIL / Test: INSUFFICIENT
  - `manifest file byte-identical to canonical bytes` → Spec: no explicit requirement found / Plan: FAIL / Impl: FAIL / Test: FAIL
- **Staleness rule:** every "preliminary" recon observation is a *hypothesis* only and must be re-established against the current repository snapshot at audit time. No earlier estimate (e.g. "~50 missing test names") may be carried in as a premise; the A–S inventory is recomputed fresh, with §11 names as the expected set.
- **Delegation:** parallel `explore` subagents for mechanical scans (REQ extraction, import graph, AST/export scans, test-body quality review), each with the evidence protocol; **I personally verify every CRITICAL/HIGH finding and every acceptance A–S verdict.** Subagents may not write files.

## 1. Phase 1 — Baseline (sandbox)

```text
go version
go build ./...      # in /tmp/opencode/phys-audit
go vet ./...
go test ./... -count=1 -v > test_full.log   (record pass/fail/skip counts per package)
go list ./...
go list -deps ./...   # must contain the seven local packages plus standard-library
                      # packages and zero third-party modules (NOT "stdlib only"
                      # literally — local packages are expected)
```

Inventory: `find` (excluding `.git`, `.opencode`) vs plan10 §2's 39-file tree; classify every file as spec-listed / pre-existing planning doc (`plans/`, `plan10/`, `reality*.md`, `.gitignore`, `.gut`) / **extra**. Report both counts (total repo files vs 39-file MVP set) and adjudicate against spec §3 ("create or modify only") + §38. Count: Go source / test / JSON / config-doc; 7 package dirs (`internal/kernel`, `core`, `ops`, `session`, `mechanics`, `relativity`, `hypothesis`) vs total repo dir count. Verify module path, `go 1.24`, zero third-party deps. Capture pre-execution tree digest + git status.

## 2. Phase 2 — Mechanical requirement extraction

Scripted (read-only) extraction from `specs_v2_3.md`:

1. All `REQ-*` labels → unique set (verify plan's claimed **115**).
2. All lines containing `MUST`/`MUST NOT` → group by section (verify claimed **126** MUST-bearing sections); each gets a deterministic `REQ-§N-MUST-nn` ID where unlabeled.
3. MRC-001..008, acceptance A–S, §15.13.1 table, REQ-032-01..22+10a.
4. Diff against plan §12 rows **and** against the actual code/test inventory; produce the master coverage table (requirement → impl location → test → verdict). Plan §12 claims are inputs, not conclusions.

## 3. Phases 3–9 — Correctness deep dives (targeted reads)

| Phase | Primary files | Critical checks |
|---|---|---|
| 3 Kernel/mint | `internal/kernel/mint.go`, `types.go`, `core/object.go` | §5.2.1 7-point contract each; zero-value invalid; defensive copies; single mint path; **`MintObjectMust` adjudicated against the five criteria: (a) mints independently? (b) different semantic contract? (c) wraps the same `MintObject` path? (d) test-only? (e) mere `(Object,error)`→panic adapter?** — then classified, not prejudged; decode path runs same invariants |
| 4 Kinds/dims/MRC | `kernel/types.go`, `ops/arithmetic.go`, `ops/relation.go` | 18 kinds + ordinals; §6.2 matrix cell-by-cell; MRC evaluation order 001→002→003→004→005→008; Mass vs RestMass fixture; exact typed errors (errors.As, not "some error") |
| 5 Expr/canonicalization | `kernel/types.go`, `ops/simplify.go` | 10 nodes/ordinals, no callbacks, constructor-vs-Simplify boundary, child 3-key sort, Pow(x,1)/Pow(x,0)-gated/0^0/0^-1/negative+repeated powers, Sqrt rules + sign entailment, `base≠0` vs `base>0`, relation/BranchSet simplify never IDENTIFIED, defensive copies, canonical bytes |
| 6 Operations/positional | `ops/*.go`, `session/session.go`, replay in `Validate` | exact 12 signatures; §15.13.1 positions enforced **identically** in `ops.Apply`, `Session.Step`, replay; `OperationParams` field order/canonical JSON/Kind binding; inspect `TestApplyPositionalInputs` body; no `ops.Identify`, no package-level mutable state |
| 7 Calculus | `ops/transform.go`, `ops/relation.go` | Differentiate bounded rule set only; Limit executes fixed body (source trace — no function-ID shortcut); Solve exact pattern; SelectBranch 6-step incl. `selected_branch/<hash>` |
| 8 Mechanics | `mechanics/primitives.go`, `relations.go`, `manifest.json` | per-constructor expr/dim/provenance/assumptions inspected, not just existence; Test D fixture Mass vs RestMass everywhere |
| 9 Relativity | `relativity/*.go`, both manifests | 10-item manifest exactly, Velocity wrapper, RestFrame assumption; §20 chain trace; §21, §22 traces; no hardcoded E=mc² |

## 4. Phases 10–12 — Canonical JSON, session/replay, tamper

- Canonical encoders for all artifact classes (Expr, Dimension, AssumptionSet, ConventionSet, Object, Manifest, OperationParams, Step, StepEnvelope, Ledger, ResearchCandidate): field order, rational `"1/2"`, set ordering, lowercase hex, `schema_version:"1"` only on DTO, no `map[string]any`, no map iteration.
- **HIGH-PRIORITY — `Ledger` public surface (§16.21):** verify exact public `Ledger` type and `func ParseLedgerJSON(data []byte) (Ledger, error)` plus exported `Ledger` methods against §16.21. If the actual return type is an unexported `ledger`, classify `IMPLEMENTED INCORRECTLY` unless a genuine exported alias/type equivalence satisfies the exact public surface. Same treatment for `Step` accessors (§16.21-02).
- **Walk the 17-stage §16.19 pipeline in `Session.Validate`** — verify each stage exists, in order, that replay uses retained `InputCanonicals`/`OutputCanonical`/`ParamsCanonical` (never infers from hashes), and Identify replay uses a pure session-owned helper.
- Verify §16.12 16-field Step order, genesis, StepID, envelope formula, state machine vs §16.3 signatures.
- Determinism: re-run the §20 derivation twice + fresh session in the sandbox and byte-compare all canonical artifacts.
- **Tamper probes on `/tmp` copy:** 4 §16.20 modes; plus mutation protocol (below).

## 5. Phases 13–17 — Hypothesis, manifests, errors, source/AST, import graph

- `hypothesis/candidate.go`: forced HYPOTHESIS/NONE, no caller provenance, no promotion; ResearchCandidate: 14 accessors + Validate + CanonicalJSON, private seal path shared with `UnverifiedResearchCandidate.Validate`, §26.8 stage order, containment → `CandidateContainmentError`.
- **`go:embed` — direct normative check (not test-quality):** spec §33-O ("Verify both manifests load through `go:embed`") and the spec's "manifest bytes are loaded with `go:embed`". If grep shows no `go:embed`, classify as a **direct normative implementation gap** (Spec: FAIL), independent of test quality.
- **Manifest canonical bytes — two dimensions reported separately:**
  - *Spec v2.3:* are `ParseManifest`/`CanonicalManifestJSON` deterministic and correct?
  - *Plan 10:* are the committed `manifest.json` bytes exactly `CanonicalManifestJSON(ParseManifest(file))` (plan §7 pin)?
- Manual item-by-item comparison of both JSON files vs constructors; strict decode/unknown-field rejection; corpus status; anomaly records.
- 11 typed errors + errors.As + deterministic diagnostics + no-panic.
- Repo-wide AST scans: floats, `math.*`, parser/scanner/token, rand/http/exec, `map[string]any`, exported `NewObject`, promotion/truth/bypass names, `ops.Identify`, Lorentz-ID shortcut, package-level mutables, `review/` package, deferred packages.
- Import graph from `go list -deps` + per-file import grep (kernel→∅, core→kernel, ops→core+kernel, session→core+ops+kernel, domains→core+kernel; `Challenge`/`Review` only in `core/corpus.go`; no cycles).

## 6. Phases 18–20 — Acceptance, negative suite, test quality

- **Recompute the current A–S inventory from the repository at audit time.** §11 names are the authoritative expected set. No earlier "missing names" estimate is a premise.
- For each A–S + `TestApplyPositionalInputs`: existence under exact plan name/package, then **body inspection** (name alone is not coverage).
- Negative suite REQ-032-01..22+10a mapped to tests with required error modes; look for tests that pass for the wrong reason (assertion-free stubs, `t.Log`-only, non-nil-only).
- Quality audit: no dimension/provenance checks, hardcoded-answer tests, internal-vs-external package mismatches vs plan §2 table, tests that would survive a forbidden API being introduced.

### Mutation protocol (explicit; sandbox only)

```text
baseline passes in /tmp sandbox
→ introduce exactly one controlled mutant
→ expected test fails (record which test, which assertion)
→ restore mutant
→ baseline passes again
```

Any mutant that does **not** kill a test → test-quality finding. Planned mutants: weaken replay to ignore retained inputs; break §15.13.1 positional binding; remove `go:embed`-equivalent loading path... (adjust to actual implementation).

## 7. Phase 21 — Manual traces (code-walk, not test-based)

10 traces (Mass+RestMass rejection; KE differentiation; NewtonSecondLaw+manifest; E=mc² §20; Lorentz limit; Step→Commit→replay; 3 tamper modes; contamination; containment; Seal→candidate→Parse→Validate), each reported `INPUT → VALIDATION → TRANSFORMATION → METADATA → OUTPUT → HASH/REPLAY EFFECT` with file:line.

## 8. Phase 22 — Report (terminal only)

11 sections: EXECUTIVE VERDICT → baseline → coverage counts → correctness findings (with four-status matrix) → test-quality findings (separate) → architecture/API → determinism → session/replay → physics derivations → exact remediation list (ordered; each tagged implementation-change / test-change / docs-only / no-action) → final audit status. No report file written.

## 9. Prioritized probe list (hypotheses only — every item re-established fresh)

1. Acceptance test-name inventory (recompute; §11 expected set).
2. Session/hypothesis test substance (assertion-free stubs suspected).
3. §16.3 signature conformance (`Seal`, `Validate`, `Draft`, `Postulate/Declare/Define/Step` labels, `Conclude`).
4. `go:embed` presence (spec §33-O direct gap if absent) + manifest byte dimension split (spec vs plan).
5. `MintObjectMust` five-criteria adjudication (§5.2.1).
6. `Ledger`/`Step` exported surface vs §16.21 (**high-priority correctness check**).
7. `TestRepositoryTreeExact` scope (skips `plan10/`, `plans/`; never asserts count 39).
8. Manifest `domain` counts vs spec §17.2/§18.3 + plan pins.

**Preliminary clean areas (re-verify):** module identity, stdlib-only, import-graph direction, no floats/forbidden imports in non-test sources, manifest item sets, README integrity-vs-authenticity + constructor-authority statements, docs §35 sections.

## 10. Execution order & integrity gate

`1 baseline → 2 extraction → (3–9 parallel subagent sweeps, I verify) → (10–12 own deep read + probes) → (13–17 mixed) → (18–20 own) → 21 traces → 22 report`.

Gate: pre/post SHA-256 tree digest identical (excluding `.git`); pre/post `git status --porcelain` identical; probes only under `/tmp/opencode/`; `.opencode/plans/plan.md` disclosed as pre-audit artifact.

---

Scope decisions:

- Pre-existing `plans/`, `plan10/`, `reality*.md` reported as a counted discrepancy with both interpretations (total repo files vs implementation-created files), adjudicated on spec §3/§38 wording.
- All `go` commands run in the sandbox copy; original provably untouched by audit execution.
- Spec correctness and plan conformance are never conflated; every finding carries SPEC/PLAN/IMPLEMENTATION/TEST statuses.

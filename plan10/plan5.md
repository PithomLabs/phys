# Final Test-Proof Remediation Plan — Plan 5 (M12 / M13 / M18)

**Disposition:** tests-only remediation expected; production changes only if reproduction uncovers a real defect. No architecture, file-tree, package, or spec/plan-doc changes.

**Authoritative sources:** `specs_v2_3.md` (sole normative authority), `plan10_v2_3.md` (normative A–S names), `plan4.md` (prior remediation, frozen), `adv_review2.md` (mutation findings only, not authority).

**Locked decisions:** M18 test lives in `ops/operations_test.go`; campaign roster is exactly the audit's M1–M18 translated to current file/function locations (stale refs such as `validate.go` / `load.go` must be remapped, never copied blindly).

## Phase 0 — Baseline + code-truth check (read-only first, then run)

1. `go build ./...`, `go vet ./...`, `go test ./...` — must be green before changes.
2. Confirm in current source: `Session.Validate` stage 10 (`session/session.go` → `checkParamsBinding`) and stage 17 (→ `checkContainment`) exist; `Limit` (`ops/transform.go`) does expand → substitute → simplify; `checkContainment` recomputes the HYPOTHESIS law + corpus-NONE rule with `CandidateContainmentError`.
3. Record that the audit's stale file refs (`session/validate.go`, `kernel/load.go`, `mechanics/mechanics.go`) do not exist; every mutant below is mapped to its real location with the translation documented.

## Phase 1 — Reproduce the three survivors first

In a disposable copy (`/tmp`, removed `.git/plan10/plans/.opencode`), apply each mutation singly, run focused tests with `-count=1`, restore:

- **M12:** neuter stage 17 only (e.g., `checkContainment` → `return nil`, or skip the stage-17 call). Observe whether the existing `containment_corpus_status` subtest kills it.
- **M13:** neuter stage 10 only (`checkParamsBinding` call or body). Observe whether the existing `params_binding` subtest kills it.
- **M18:** `Limit` → hardcoded `Rational(1)` result. Observe whether `TestLorentzFactorLimit` (+ merged-assumption asserts) kills it.
- Record the honest outcome per mutant plus the **first validation stage reached before failure** (via the error message/stage tag). A new test is accepted only if its failure is attributable to the exact removed invariant — evidenced by the message pin pointing at that stage — not to a later pipeline stage. Two branches:
  - **(a) Already killed:** the audit ran stale code; still add the three named tests in Phase 2 as explicit, audit-traceable proof.
  - **(b) Genuinely survives:** root-cause why (e.g., tamper caught by an earlier stage, message not pinned, assumption not propagated), then design the new test to isolate exactly that stage.

## Phase 2 — Add the three named tests (minimum, no renames)

- **`TestValidateCandidateContainmentStage`** (`session/session_test.go`): real HYPOTHESIS object → `Declare` → `Commit` → tamper retained step via the established unsafe fault-injection helpers so the output stays structurally canonical with recomputed `InputHashes/OutputHash/CurrentStepHash`, keeping the HYPOTHESIS dependency as HYPOTHESIS (provenance propagation stays correct, so stage 15 has nothing to reject) while the retained output carries `CorpusStatus != NONE` (upgrade both input+output bytes of the assertion step, since transformation replay subsumes containment — replay equality enforces the law there). Call `s.Validate()`; assert `errors.As(err, &CandidateContainmentError)` **and** a containment-indicating message (stages 1–16 must pass by construction). Do NOT use the trusted-provenance mutation as primary (stage 15 owns provenance/MRC). Also assert `Ledger.Validate` alone does not produce the containment verdict (proves stage-17 specificity). Reuse the `mirrorStep`/`ledgerMirror` fidelity check.
- **`TestValidateOperationParamsStage`** (`session/session_test.go`): commit a valid `add` step on `[a, b]`, then replace `ParamsCanonical` with a canonically valid `Kind:"pow"` object (`{"kind":"pow","exponent":"2/1","operator":"","justification":""}` — shape-valid per `ParseOperationParams`, binding-invalid per §15.13 which requires `empty` for `add`), recompute step/chain hashes so stages 1–9 pass. Prove `ParseOperationParams` succeeds in-test. Assert `Session.Validate()` rejects **with the stage-10 binding message** (`"operation params:"`), and assert `Ledger.Validate()` **passes** the same ledger. Rationale: `ops.Apply` re-validates params at replay (`ops/dispatch.go`), so a removed stage 10 surfaces as a stage-12 replay failure with a different message — the message pin is what isolates stage 10, not merely the rejection. Do NOT weaken the independent `ops.Apply` check to make this test work. Removing stage 10 must make the test fail on the message assertion.
- **`TestLimitDirectSubstitutionNonUnity`** (`ops/operations_test.go`, using existing `sym`/`rat`/`exprObj`/`defined` helpers): `Limit(Add(x,2), x, 3)` with dimensionless operands → exact `{"kind":"rational","value":"5/1"}`, `Kind==Expression`, dimension dimensionless. Under a `Limit→1` mutant this fails; `TestLorentzFactorLimit` (exact `1/1` + assumption asserts) stays untouched. No asymptotics, no `>1` criterion.
- Quality gates: deterministic values only; M12 message must implicate containment (not structural); M13 params must be shape-valid (prove by `ParseOperationParams` success in-test); M18 expected value ≠ 1.

## Phase 3 — 18-mutant campaign (audit IDs translated, exactly 18)

Map each to its current location before running; killed = the intended semantic test fails for the intended reason (build breakage does not count).

| ID | Translation (verify during mapping) | Expected killer |
|---|---|---|
| M1 | Remove Validate stage 1 (structure/state) | structural/empty-ledger tests |
| M2 | Seal skips validation | `TestSealValidatesFirst` |
| M3 | Corrupt retained canonical bytes | Q output/hash subtests |
| M4 | Swap committed step order | chain-linkage subtests |
| M5 | Bypass Identify path (use `ops.Apply` for Identify / drop session helper) | K tests |
| M6 | Weaken `MintObject` invariant | `TestMintObjectRejectsInvalidSpec` |
| M7 | Remove/bypass `Limit` body traversal | J + new non-unity Limit test |
| M8 | Break `LorentzFactor` constructor | H/J/manifest cross-check |
| M9 | Skip `LoadObjectJSON` invariant/canonical-equality check | round-trip + tamper tests |
| M10 | Accept malformed `stepEnvelope` | strict-parse tests |
| M11 | Corrupt mass dimension | D (Mass vs RestMass) + dimension tests |
| M12 | Remove stage 17 | **new** containment-stage test |
| M13 | Remove stage 10 | **new** params-stage test |
| M14 | Corrupt transform op behavior (e.g., Add) | C/positional/arithmetic tests |
| M15 | Weaken mint authority (second mint path) | `TestNoGenericFactory` |
| M16 | Break mass-energy relation | I (+ manifest cross-check) |
| M17 | Corrupt mechanics force unit | B/manifest cross-check |
| M18 | Hardcode `Limit→1` | **new** non-unity Limit test |

Run all 18 in the disposable copy (backup/restore per mutant, `-count=1`), then restore pristine and re-run the full suite. Target: **18/18 KILLED**.

## Phase 4 — Regression + clean-worktree verification

- Re-run `go build/vet/test ./...` in place; spot-check the Task 8 lists (Validate's 13 aspects, candidate's 4, relativity's 6, mint's 4) via the named suites (Q/R/S/K/L/M/N/O/I/J + `TestNoGenericFactory`).
- Fresh copy: verify the exact 39-file tree, 7 packages, `go list -m all` stdlib-only, no deferred packages, no new exports, `specs_v2_3.md` + `plan10_v2_3.md` unmodified (do not touch any doc until 18/18 is measured).
- Deliver the 9-section report (baseline, root cause, changes, test proof, 18-row mutation table, regression, architecture, gaps, final PASS / PASS WITH FINDINGS / FAIL verdict).

## Explicit non-goals

No production refactor unless Phase 1 finds a real defect (then the smallest spec-conformant fix + explanation); no `Unwrap` on `CandidateContainmentError`; no test weakening/removal; no A–S renames; no hook-only-for-mutants helpers (reuse the existing unsafe fault-injection); no doc-claim updates before 18/18 is measured.

*Status: IMPLEMENTED 2026-10-01. 18/18 meaningful mutants KILLED (disposable copy, per-mutant backup/restore, `-count=1`); `go build ./...`, `go vet ./...`, `go test ./... -count=1` green in place and in clean copy; 39-file tree, 7 packages, stdlib-only preserved; `specs_v2_3.md` / `plan10_v2_3.md` untouched.*

## Implementation record

- Phase 1 reproduction: M12 killed already by `containment_corpus_status`; M13 killed already by `params_binding` message pin (first-failure-stage without stage 10 = stage 12 replay); M18 killed at relativity level by merged-assumption asserts but survived at ops level (ops J checks only `1/1` bytes) — confirming the audit against current code for M18 only.
- Added: `TestValidateCandidateContainmentStage`, `TestValidateOperationParamsStage` (`session/session_test.go`); `TestLimitDirectSubstitutionNonUnity` (`ops/operations_test.go`, `Limit(Add(x,2),x,3)→5/1` + 3-parent-hash proof of substitution path).
- Added minimal subtests inside normative tests (no new top-level names, deviation documented): `step_id_format` in Q (M10), `New().Validate()` assertion in `TestSessionStateMachineTransitions` (M1).
- M8 mapping fix: killer is `TestManifestConstructorCrossCheck` (per-item dimension check), not the summary `TestRelativityManifestCrossCheck`.
- No production semantic changes. `ops.Apply` param validation deliberately untouched (M13 isolation relies on distinguishing stage 10 from stage 12).


## Implementation record (Plan 6: M9/M15, tests only)

- Reproduced both survivors in disposable copy: M9 (full suite green without re-encode check), M15 (`TestNoGenericFactory` green with `MintObjectAlt`). No production defect found.
- `TestObjectCanonicalRoundTripRejectsNonCanonical` (`core/object_test.go`): canonical round-trips byte-identical + hash-equal; top-level-only whitespace variant (nested bytes intact, so sub-component checks pass) rejected with `"non-canonical object JSON"`. First `MarshalIndent` variant discarded: nested bytes change, so sub-checks reject before the top-level check (wrong attribution).
- `TestNoGenericFactory` strengthened: primary exact-set invariant over kernel `Object`-returning functions = `{MintObject, LoadObjectJSON}` (any third name fails); secondary single-literal-site + decoder-shape checks; existing name counts/bans kept.
- Full translated M1–M18 campaign: **18/18 KILLED** (per-mutant restore, `-count=1`, kills attributed to intended tests).
- No production changes. `specs_v2_3.md` / `plan10_v2_3.md` untouched.

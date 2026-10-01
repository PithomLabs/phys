# Adversarial Review 11 — Re-review of Documentation Amendment v2.4 Remediation

**Reviewer:** Independent adversarial reviewer (senior Go architect)
**Date:** 2026-10-01
**Status:** DOCUMENTATION PASS — 42-FILE RE-FREEZE AUTHORIZED
**Scope:** F-001, F-002, F-003 remediation verification; narrow documentation-vs-source consistency review; residual AGENTS.md wording check; 42-file tree enforcement; frozen-spec non-regression.

---

## A. Re-review Verdict

**F-001 — RESOLVED.** `TestRepositoryTreeExact` now admits exactly 42 paths. The five formerly admitted extra files (`.gitignore`, `.gut`, `reality-first-physics-v8.0.md`, `reality_check.md`, `reality_first_plain_english.md`) are absent. Independent intrusion testing in a disposable copy proves the test fails when each of those five paths is added to the filesystem. There is no wildcard, no `*.md` exception, no `<=42` loose rule, and no plan/audit/tooling allowance.

**F-002 — RESOLVED.** `AGENTS.md` §2 now correctly states that provenance status and corpus status are two independent axes, using `MassEnergyRelation` (provenance `DERIVED`, corpus status `ESTABLISHED`) as the canonical example. The complete provenance set (`DEFINED`, `POSTULATED`, `DERIVED`, `IDENTIFIED`, `APPROXIMATED`, `HYPOTHESIS`) and the complete corpus-status set (`NONE`, `ESTABLISHED`, `CONTESTED`, `SUPERSEDED`, `FALSIFIED`) are both present. Pure-op propagation, `HYPOTHESIS` contamination, `Session.Identify` as the sole `IDENTIFIED` creator, and corpus status as human-curated metadata are all correctly stated.

**F-003 — RESOLVED.** `mechanics/README.md` now states that anomaly `galilean_noninvariance` is on `NewtonSecondLaw` and `KineticEnergyRelation`, matching the manifest's `related_items` exactly.

**F-004 — NOT A FINDING.** The phrase in `AGENTS.md` §5:

> "construction fixed constructors mint typed objects (the only way in)"

is **not** materially misleading. The table row labeled "construction" describes one specific mathematical act among several (transformation, relation artifact, bounded solving, branch selection, identification, hypothesis formation). The parenthetical colloquially means "the only way to create the initial typed domain objects" — not "the only valid trusted-object minting path in the whole system." The full document (§6 hypothesis formation, §5 identification, §5.4 trusted construction in the spec) correctly enumerates all five legitimate minting paths: fixed domain constructors, `hypothesis.NewCandidateConcept`, pure `ops` results, session-derived artifacts, and internal test fixtures. An AI agent reading the complete document would not be misled into believing that hypothesis objects or session-identified objects are invalid.

---

## B. Findings

**No findings.** F-001, F-002, and F-003 are all resolved. F-004 is not a finding. The narrow documentation-vs-source scan revealed no materially false claims.

---

## C. Verified Invariants

| Invariant | Status |
|---|---|
| Exact 42-file tree | **CONFIRMED** — 42 non-plan files present; 5 formerly admitted extras absent |
| Closed-world tree test | **CONFIRMED** — `TestRepositoryTreeExact` has exactly 42 allowlisted paths; intrusion test fails for all 5 formerly admitted paths |
| v2.3 spec byte-for-byte frozen | **CONFIRMED** — `md5sum b31d0cd9d457b59c2aabd183d0f9562b` matches git HEAD |
| Exactly three new repository documentation files | **CONFIRMED** — `AGENTS.md`, `mechanics/README.md`, `relativity/README.md` |
| No manifest changes | **CONFIRMED** — `mechanics/manifest.json` and `relativity/manifest.json` unchanged |
| No production-code changes | **CONFIRMED** — `git diff` shows only `core/object_test.go` modified |
| Package boundaries unchanged | **CONFIRMED** — import graph unchanged |
| Documentation matches canonical metadata | **CONFIRMED** — READMEs match manifests |
| E=mc² anti-shortcut rule preserved | **CONFIRMED** — `relativity/README.md` and `derivation_test.go` consistent |
| Special Relativity boundary preserved | **CONFIRMED** — `relativity/README.md` explicitly excludes GR, gravitational dynamics, curved spacetime, black-hole physics, quantum mechanics, quantum gravity |
| Provenance/corpus two-axis model preserved | **CONFIRMED** — `AGENTS.md` §2 correctly describes independent axes |

---

## D. Independent Test Results

```
go build ./...          # PASS — no output, no errors
go vet ./...            # PASS — no output, no errors
go test ./...           # PASS
  core                  ok
  hypothesis            ok
  mechanics             ok
  ops                   ok
  relativity            ok
  session               ok
```

Targeted checks:
- `TestRepositoryTreeExact` — PASS (clean tree); **FAIL** (as required) when any of the 5 formerly admitted intrusion files is added to a disposable copy
- `TestFullDerivationSequence` — PASS
- `TestMassEnergyDerivation` — PASS
- `TestLorentzFactorLimit` — PASS
- `TestDerivationUsesNoStoredResult` — PASS
- `TestNoRelativityLeakage` — PASS
- `TestManifestConstructorCrossCheck` (mechanics + relativity) — PASS
- `TestMechanicsManifestCrossCheck` — PASS
- `TestRelativityManifestCrossCheck` — PASS
- `TestReverseConstructorAllowlist` (mechanics + relativity) — PASS
- `TestDifferentiateKineticEnergyOps` — PASS
- `TestFinalAssumptionOwnership` — PASS

Intrusion test (disposable copy):
```
cd /tmp/phys-intrusion-test && go test ./core -v -run TestRepositoryTreeExact
--- FAIL: TestRepositoryTreeExact (0.00s)
    object_test.go:811: unexpected file in repository tree: .gitignore
    object_test.go:811: unexpected file in repository tree: .gut
    object_test.go:811: unexpected file in repository tree: reality-first-physics-v8.0.md
    object_test.go:811: unexpected file in repository tree: reality_check.md
    object_test.go:811: unexpected file in repository tree: reality_first_plain_english.md
```

---

## E. Re-freeze Decision

**DOCUMENTATION PASS — 42-FILE RE-FREEZE AUTHORIZED**

# Adversarial Review 10 — Documentation Amendment v2.4

**Reviewer:** Independent adversarial reviewer (senior Go architect)
**Date:** 2026-10-01
**Status:** DOCUMENTATION FAIL — 42-FILE RE-FREEZE NOT AUTHORIZED
**Scope:** AGENTS.md, mechanics/README.md, relativity/README.md, TestRepositoryTreeExact, directly affected tests, actual manifests, actual exported API/package boundaries, relevant production implementation, relevant derivation tests.

---

## 1. Executive Verdict

The documentation amendment is **not** accurate enough to accept in its current state, and the repository is **not** safe to re-freeze at exactly 42 files. Three independent defects block re-freeze authorization:

- A **BLOCKER** in `core/object_test.go` (TestRepositoryTreeExact): the allowed-map contains 5 paths that are not part of the 42-file invariant. The test admits files it must reject.
- A **HIGH** epistemic-model defect in `AGENTS.md` §2: the phrase "Every object you touches belongs to exactly one layer" is semantically inaccurate because provenance status and corpus status are independent dimensions. `MassEnergyRelation` is a concrete counterexample (provenance `DERIVED`, corpus status `ESTABLISHED`).
- A **MEDIUM** documentation-consistency defect in `mechanics/README.md`: the anomaly scope ("on all three relations") does not match the manifest (`related_items` lists only two relations).

Build, vet, and the full test suite pass independently. The actual repository tree already contains exactly 42 non-plan files. The defects are in the documentation and the tree-test allowlist, not in the production implementation.

---

## 2. Findings Table

| FINDING-ID | Severity | File | Exact Section / Line | Observed Claim | Authoritative Evidence | Why It Is Wrong or Misleading | Minimal Correction |
|---|---|---|---|---|---|---|---|
| F-001 | **BLOCKER** | `core/object_test.go` (TestRepositoryTreeExact) | Lines 750–799 | Allowed map admits `.gitignore`, `.gut`, `reality-first-physics-v8.0.md`, `reality_check.md`, `reality_first_plain_english.md` | `specs_v2_3.md` §3 lists exactly 39 implementation files; user requirement states the 42-file invariant must be exact; actual `find` shows the repository contains exactly the 42 expected paths and none of the 5 extras | The test does not enforce the exact 42-file allowlist. If any of the 5 extra files were added to the repository, the test would not reject them, violating the freeze invariant. | Remove the 5 extra entries from the `allowed` map so it contains exactly the 42 paths (39 spec files + `AGENTS.md` + `mechanics/README.md` + `relativity/README.md`). |
| F-002 | **HIGH** | `AGENTS.md` | Line 24 | "Every object you touches belongs to exactly one layer. Never confuse them." | `specs_v2_3.md` §13.1 lists provenance statuses (`DEFINED`, `POSTULATED`, `DERIVED`, `IDENTIFIED`, `APPROXIMATED`, `HYPOTHESIS`); §13.4 lists corpus statuses (`NONE`, `ESTABLISHED`, `CONTESTED`, `SUPERSEDED`, `FALSIFIED`); `relativity/relations.go` lines 75–86 mint `MassEnergyRelation` with `StatusDerived` and `CorpusEstablished` | Provenance and corpus status are independent axes. `MassEnergyRelation` simultaneously carries provenance `DERIVED` and corpus status `ESTABLISHED`, belonging to two layers at once. The wording implies a false mutual exclusion. | Rephrase to make the two dimensions explicit, e.g.: "Objects carry two independent status axes: provenance (DEFINED / POSTULATED / DERIVED / IDENTIFIED / HYPOTHESIS) and corpus status (NONE / ESTABLISHED / …). Do not confuse them." |
| F-003 | **MEDIUM** | `mechanics/README.md` | Lines 52–53 | Anomaly `galilean_noninvariance` is "on all three relations" | `mechanics/manifest.json` anomaly entry: `"related_items":["NewtonSecondLaw","KineticEnergyRelation"]` — only two relations, not three | The README scope ("all three") does not match the manifest scope (two items). | Correct to "on two of the three relations (`NewtonSecondLaw`, `KineticEnergyRelation`)" or match the manifest exactly. |

---

## 3. Required Corrections

**F-001 (BLOCKER) — `core/object_test.go`**

The `allowed` map in `TestRepositoryTreeExact` must be reduced from 47 entries to exactly 42 entries. The 5 entries to remove are:

- `.gitignore`
- `.gut`
- `reality-first-physics-v8.0.md`
- `reality_check.md`
- `reality_first_plain_english.md`

The remaining 42 entries must be exactly the 39 files listed in `specs_v2_3.md` §3 plus:

- `AGENTS.md`
- `mechanics/README.md`
- `relativity/README.md`

No `*.md` exception, no `<=42` loose rule, and no admission of plan/audit/tooling artifacts.

**F-002 (HIGH) — `AGENTS.md`**

Replace the sentence "Every object you touches belongs to exactly one layer. Never confuse them." with wording that reflects the independent dimensions of provenance and corpus status, using `MassEnergyRelation` as the canonical example if helpful.

**F-003 (MEDIUM) — `mechanics/README.md`**

Change "on all three relations" to match the manifest's `related_items` list.

---

## 4. Verified Invariants

- The **actual repository tree** contains exactly 42 non-plan files: the 39 frozen implementation files from `specs_v2_3.md` §3 plus `AGENTS.md`, `mechanics/README.md`, and `relativity/README.md`. No extra files exist outside `plan10/` and `plans/`.
- `plan10/specs_v2_3.md` is byte-for-byte frozen and unchanged.
- The **v2.4 amendment document** is not present in the repository tree.
- `mechanics/manifest.json` reports `framework_id: classical_mechanics`, `corpus_status: established`, 9 domains, 11 items, anomaly `galilean_noninvariance`.
- `relativity/manifest.json` reports `framework_id: special_relativity`, `corpus_status: established`, 6 assumptions (2 constraints + 4 framework), 8 domains, 10 items, anomaly `no_gravity`.
- `MassEnergyRelation` in `relativity/relations.go` is minted with `StatusDerived`, `CorpusEstablished`, and `derivable_from: ["EnergyMomentumRelation", "RestFrame"]`.
- Package boundaries are enforced: `mechanics` and `relativity` import only `core` + `internal/kernel`; `session` imports `core`, `ops`, `internal/kernel` but not any domain package; `ops` imports only `core`; `core` imports only `internal/kernel`.
- All 12 supported operations (`Add`, `Subtract`, `Multiply`, `Divide`, `Pow`, `Simplify`, `Substitute`, `Differentiate`, `Limit`, `Compare`, `Solve`, `SelectBranch`) are closed over in `ops/dispatch.go`. `identify` is deliberately excluded.
- `Limit` expands `lorentz_factor` via fixed-body traversal (`expandCalls` in `ops/transform.go`), never by ID shortcut.
- `Differentiate` supports only the bounded engine specified in `specs_v2_3.md` §15.8.
- `Solve` accepts only the exact quadratic pattern `Relation(eq, Pow(Symbol(target),2), rhs)`.
- `Session.Identify` is the sole creation path for `IDENTIFIED`; no `ops.Identify` exists.
- No promotion API exists anywhere in the exported surface.
- The E=mc² derivation test (`relativity/derivation_test.go`) never references `MassEnergyRelation` as a premise (enforced by AST scan in `TestDerivationUsesNoStoredResult`).
- The Lorentz-limit test (`TestLorentzFactorLimit`) asserts body expansion by checking that the merged `speed_of_light_positive` assumption is present on the result.
- The mechanics differentiation test (`TestDifferentiateKineticEnergyOps`) asserts the bounded engine and MRC-003 Expression-vs-named comparison path.

---

## 5. Test Results

Independent execution on the unmodified repository:

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

Relevant targeted tests:

- `TestRepositoryTreeExact` — PASS (but see F-001: the allowlist is too permissive even though no offending files currently exist)
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

All MRC rules (MRC-001 through MRC-008) are exercised and green in the existing test suite.

---

## 6. Re-Freeze Authorization

The independent reviewer finds that the documentation amendment contains a **BLOCKER** tree-test defect and a **HIGH** epistemic-model wording defect, plus a **MEDIUM** mechanics-anomaly consistency defect. These must be remediated before re-freeze.

**DOCUMENTATION FAIL — 42-FILE RE-FREEZE NOT AUTHORIZED**

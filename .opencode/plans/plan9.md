# Final Corrective Pass — Plan (staged for plan10/plan9.md, no implementation yet)

NOTE: writes to `plan10/plan9.md` are blocked by the workspace permission rules (edit denied outside narrow paths).
This is a staging copy at `.opencode/plans/plan9.md`. Move it to `plan10/plan9.md` with:
`mv .opencode/plans/plan9.md plan10/plan9.md`
after adjusting permissions, or instruct me and I will do the move. Content below is the plan body.

---

Status: planning only. No implementation changes made. Awaiting explicit instruction before any build-phase execution.

User decisions recorded:
1. No `.md` deletions — the three `reality_*.md` docs stay in the repo. Nothing is deleted just for the sake of deletion.
2. `CoreObject()` methods only — free constructors keep returning `kernel.Object`.

## A. Baseline (from read-only evidence)

- Filesystem counted tree (excluding `.git/`, `.opencode/`, `plan10/`, `plans/`, i.e. the reviewer's method): **45 files**. Normative: **39**.
- The 6 extras on disk: `reality_check.md`, `reality_first_plain_english.md`, `reality-first-physics-v8.0.md`, `.gitignore`, `.gut` (72B gut metadata, git-ignored), `attack_review_main.go` (untracked helper, only partially reviewed so far — header + grep snippets only).
- Correction to the review: `docs/paper-translation.md` **now exists** on disk and is tracked — the "missing file" half of FINDING-001 is already resolved. Only the extras remain.
- `core/dimension.go:54-55` still exports `DimensionEnergySquared()`; used in `relativity/relations.go:54` and `relativity/derivation_test.go:193`.
- `mechanics/manifest_test.go:34` still says `"expected 10 domains"` against a `!= 9` check.
- `CoreObject()` returns `kernel.Object` in `mechanics/primitives.go:219-227` (9 methods) and `relativity/primitives.go:196-203` (8 methods). Both files already import `core`, and `core.Object = kernel.Object` is an alias, so the change is signature-only.
- E=mc² firewall intact: `derivation_test.go` references `MassEnergyRelation` only in comments and the AST self-audit; derivation chain unchanged.

## B. Planned corrections (4 findings)

### B1. File tree — relocate, do not delete (per user constraint)

- Full-read `attack_review_main.go` first (user instruction: review those `.go` files first), then **move** (never `rm`) out of the counted tree:
  - 3x `reality_*.md` → `plan10/` via `git mv` (content preserved, still tracked, still in repo; `plan10/` is already excluded from the implementation-tree count).
  - `attack_review_main.go` → outside the repo (e.g. `/tmp/opencode/`) after review.
  - `.gut` → outside the repo (regenerable tool metadata, untracked+ignored).
- **`.gitignore` = REMOVE (decision locked) — see §E: it is tracked and not in the §3 list. Arithmetic: 45 − 3 (md moved) − 1 (attack moved) − 1 (.gut moved) = **40 with `.gitignore` kept, 39 with it removed**. There is no other path to 39 — the rest of the tree already matches §3 exactly. `specs_v2_3.md` is not touched either way.

### B2. Remove public `DimensionEnergySquared()` — 3 files

- NO `internal/kernel` change (amended B2): compute the E² dimension locally in `relativity` via the existing generic machinery `core.DimensionEnergy().Multiply(core.DimensionEnergy())` (M¹L²T⁻² × M¹L²T⁻² = M²L⁴T⁻⁴). No physics-specific naming enters the theory-neutral kernel.
- `relativity/relations.go:54`: `core.DimensionEnergySquared()` → `core.DimensionEnergy().Multiply(core.DimensionEnergy())`.
- `core/dimension.go:54-55`: delete the public constructor + its doc comment.
- `relativity/derivation_test.go:193`: same local `Multiply` expression for the dimension assertion.
- Verify afterwards: exactly the 9 §7.0 constructors remain public; `NewDimension`/`ParseDimensionJSON` untouched (reviewer already accepted them).

### B3. Manifest message — 1 line

- `mechanics/manifest_test.go:34`: `"expected 10 domains"` → `"expected 9 domains"`. Condition untouched.

### B4. `CoreObject()` signatures — 2 files, 17 methods

- Return type `kernel.Object` → `core.Object` in `mechanics/primitives.go:219-227` and `relativity/primitives.go:196-203`. Struct fields (`Object kernel.Object`), free constructors, and test maps (`func() kernel.Object`) untouched per user option 1.

## C. Explicit non-changes

`specs_v2_3.md`, `plan10/*`, E=mc² chain (`EnergyMomentumRelation → ZeroThreeMomentum → Substitute → Simplify → Solve → Compare(gte) → SelectBranch → m*c²`), `MassEnergyRelation()` absence from executable test references, no new physics/corpus/dependencies/files; all currently-holding invariants left alone.

## D. Verification (build phase, after approval)

1. `go test ./... -count=1`, `go vet ./...`.
2. File-tree recount by the reviewer's filesystem method; `git status` review (note: tree currently has pre-existing modifications in 17 tracked files — inspect before/after).
3. Fresh public-API audit: exported dimension constructors (= 9), no generic factories/mint forwarders/provenance mint paths/compat helpers (`go doc`, `go/ast` scans as the existing tests do).
4. Freeze invariants: mint authority, immutability, canonicalization, rationals, `OperationParams`, provenance, containment, isolation, manifests, replay, StepID, no-panic, `SelectBranch`/`Solve`/`Limit`/`Differentiate`, E=mc², Lorentz limit.
5. Mutations M1–M18 + N1–N7 in disposable copies with `-count=1`, restoring clean tree between mutants; prior 25/25 claimed only if actually rerun and still killed.
6. Report in the A–E format with exactly one of `FREEZE-READY` / `NOT FREEZE-READY`.

## E. Honest status forecast — one decision deferred

With the no-deletion constraint honored via relocation and `.gitignore` removed, the tree lands at **exactly 39 files**. Keeping it means the report must conclude **NOT FREEZE-READY on REQ-003-01 alone** (all other blockers fixed, everything else verified). Removing that single tracked VCS file is the only spec-compliant route to `FREEZE-READY`.

Execution order once instructed: review `attack_review_main.go` → relocations → B2 → B3 → B4 → tests/vet → audits → mutations → final report.

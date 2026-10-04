# ADVERSARIAL REVIEW — plan10_v2_3.md vs specs_v2_3.md
## Gaps, Limitations, and MVP-Scoped Mitigations

**Scope:** MVP only. Not v0.5+.

**Verdict:** Plan 10 is the most conformant plan in the review cycle. It correctly reflects every v2.3 reconciliation, incorporates all five net-valid findings (F1–F5), and rebuilds the coverage matrix from scratch. But it has **three high-severity gaps, four medium-severity gaps, and several low-severity items** that will force the coding agent to make decisions the plan claims to eliminate. Against the separate corpus-integrity audit prompt (`prompt2.md`), the plan is **not yet fully auditable** — it lacks explicit tests for several phases the audit requires.

---

## Part I: What the Plan Gets Right

Confirmed against `specs_v2_3.md`:

- Module path (`github.com/PithomLabs/phys`), Go 1.24, stdlib-only — REQ-000-01/02 ✓
- 39-file tree with exact composition (24 source / 10 test / 2 manifest / 3 config) — REQ-003-01, REQ-§38 ✓
- Dependency graph with `core` free of `ops`/`session`, `ops` free of `session`, `session` free of domain packages — REQ-004 ✓
- Constructor authority via `internal/kernel.MintObject` with unexported fields and no public factory — REQ-005, MRC-001 ✓
- `Challenge`/`Review` exclusively in `core/corpus.go` — Reconciliation A ✓
- Finite-factor rule with `base != 0` vs `base > 0` distinction — Reconciliation B ✓
- Exact §15.13.1 positional table reproduced — Reconciliation C ✓
- Ten expression node kinds, ordinals 0–9 — REQ-008 ✓
- Eight MRC rules with enforcement locations — REQ-014 ✓
- Nine session actions, state machine, genesis hash, StepEnvelope, 17-step Validate — REQ-016 ✓
- `UnverifiedResearchCandidate` split — REQ-§26.10 ✓
- All 22 deferrals — REQ-§37 ✓
- All A–S acceptance tests with named test functions — REQ-§33 ✓
- Matrix rebuilt with subsection-precise IDs and expanded REQ-002-01..26, REQ-032-01..22+10a ✓
- Test C uses `DimensionMismatchError` at MRC-002 (isolating from MRC-003) ✓
- Test D uses equal-dimension distinct-kind fixtures (isolating MRC-003) ✓
- `schema_version:"1"` labeled as Plan 10 serialization pin, not kernel state ✓
- Revision Audit section honestly tracks preserved/adopted/removed material ✓
- `Open Spec Items: NONE` justified by v2.3 reconciliations ✓

This is a strong plan. The gaps below are the exceptions, not the rule.

---

## Part II: High-Severity Gaps

### Gap H1 — Test L does not exercise every HYPOTHESIS-propagating operation

**What the plan says** (§11 step 12):

> "L: create candidate concept → derive with it (Add/Simplify/Solve/Limit/Compare/SelectBranch/Identify) → every downstream status `HYPOTHESIS`"

**What specs require** (§13.2):

> "For every pure `ops` transformation: if any input status == `HYPOTHESIS`: output status = `HYPOTHESIS`... This applies to: **Add, Subtract, Multiply, Divide, Pow, Simplify, Substitute, Differentiate, Limit, Compare, Solve, SelectBranch**" (12 operations).

**The gap:** Test L lists only **7 of 12** operations. Missing: Subtract, Multiply, Divide, Pow, Substitute, Differentiate. A candidate-derived result through `Multiply(candidate, mass)` could theoretically be mis-marked `DERIVED` and the plan's slice-L test would not catch it.

**Impact:** If a coding agent implements status propagation correctly for 7 ops but not the other 5, the MVP's candidate containment guarantee is broken and Test L will not detect it.

**Mitigation:**

Extend §11 step 12 to enumerate all 12 operations explicitly, or require the test to iterate a table of `{op, inputs}` pairs and assert `HYPOTHESIS` output for each. Alternatively, split the responsibility: `TestProvenanceStatusLaw` (already in the matrix at REQ-013-02, "table over all 12 ops × statuses") must be the authoritative test; `TestHypothesisContamination` should cross-reference it. Either way, §11 must not present a shorter list as the acceptance condition.

**Suggested edit to plan §11 step 12:**

```
12. L: create candidate concept; derive with it through ALL 12 ops
    (Add, Subtract, Multiply, Divide, Pow, Simplify, Substitute,
     Differentiate, Limit, Compare, Solve, SelectBranch) and through
    Session.Identify; every downstream status `HYPOTHESIS`. Where an
    op's input-compatibility cannot be established with a candidate
    operand, use the candidate as the second operand of the op and a
    compatible trusted first operand.
```

---

### Gap H2 — Adversarial manifest corruption tests are under-specified

**What the plan says** (§11.O):

> "negative: `TestManifestRejectsInconsistent` (unknown field, bad enum, non-canonical expr, duplicate ID) → `ManifestValidationError` (REQ-032-15)"

**What the corpus audit requires** (`prompt2.md` Phase 22):

> "Introduce controlled corruptions into: `mechanics/manifest.json`, `relativity/manifest.json`, domain constructors. At minimum: change F=ma; change p=mv; change K expression; change E² relation; remove m≥0; remove c>0; change Minkowski signature; change RestFrame assumption; change provenance status; change corpus status; change a dimension; add an undocumented trusted relation. Verify the manifest/constructor cross-checks and corpus tests detect each corruption."

**The gap:** The plan's negative suite names **4 corruption types** (unknown field, bad enum, non-canonical expr, duplicate ID). The audit requires **12 corruption types**. The 8 missing types all fall under "manifest content changed but schema still valid" — these are only detected by the constructor cross-check, not by schema validation. The plan implies this works (§7 "file bytes == canonical + constructor cross-check"), but does not state it as an explicit test row, and `TestManifestRejectsInconsistent` as named does not suggest it covers F=ma mutation or removal of `m≥0`.

**Impact:** The coding agent will implement schema validation but may not implement the constructor-mismatch path for every manifest field. When the corpus-integrity audit runs Phase 22, it will find 8 undetected corruptions and the audit will return FAIL.

**Mitigation:**

Add to plan §11.O an explicit "corruption detection" sub-list and require a test like `TestManifestConstructorMismatchDetected` that iterates the following table:

```
Corruption                        Detected by
change F=ma                       constructor Expr cross-check
change p=mv                       constructor Expr cross-check
change K expression               constructor Expr cross-check
change E² relation                constructor Expr cross-check
remove m≥0 from manifest          assumptions cross-check
remove c>0 from manifest          assumptions cross-check
change Minkowski signature        assumptions cross-check
change RestFrame assumption       assumptions cross-check
change provenance_status          provenance cross-check
change corpus_status              corpus-status preservation check
change dimension                  dimension cross-check
add undocumented relation         exact item-ID set assertion
```

Because the plan says "file bytes == `CanonicalManifestJSON(ParseManifest(file))`" is asserted (guards against silent byte drift), the corruption must be introduced by editing the manifest bytes and re-canonicalizing. The test should bypass the file-byte check (inject corrupted bytes programmatically) and confirm the cross-check detects each row.

**This is the single most important gap for `prompt2.md` to pass.**

---

### Gap H3 — No explicit test for kernel physics-agnosticism

**What the plan says** (§1):

> "No `core → ops` dependency exists... `core` MUST NOT import `ops` or `session`."

And (§4):

> "`kernel` is physics-agnostic; physics facts supplied by domain/corpus layers."

**What the audit requires** (`prompt2.md` Phases 15 and 16):

> "Inspect `internal/kernel`. Search for: mechanics-specific symbols; relativity-specific symbols; Newton; Lorentz; Minkowski; mass-energy; force; momentum; velocity; `F=ma`; `E=mc²`; framework-specific assumptions; special physics constants. ... If the kernel itself contains a physical law merely to make an MVP example work, flag it as KERNEL PHYSICS LEAKAGE."

**The gap:** The plan asserts kernel is physics-agnostic in prose but provides no test that mechanically verifies this. `TestInternalKernelBoundary` (per §5) checks the `internal/` path and the absence of `MintObject` re-export — it does **not** scan kernel source for physics symbols like `m`, `c`, `lorentz_factor`, `rest_frame`, etc.

The audit will run Phase 15 and can report leakage in prose, but the plan's "no remaining architectural choices" claim implies the implementation should enforce physics-agnosticism mechanically.

**Impact:** A coding agent implementing canonicalization might add a special case for `lorentz_factor` or a shortcut for `rest_mass_nonnegative` inside `internal/kernel/types.go` "for convenience." The plan's prose forbids this; the tests do not.

**Mitigation:**

Add a test `TestKernelNoPhysicsSymbols` in `core/object_test.go` (external package, so it can parse kernel source using `go/parser` — which is banned in non-test source but allowed in tests):

```
Scan all non-test .go files in internal/kernel/ for string literals matching:
  "lorentz_factor", "rest_frame", "minkowski", "einstein", "newton",
  "F = ma", "E = mc", "rest_mass_nonnegative", "speed_of_light_positive",
  "classical_mechanics", "special_relativity", and the symbol strings
  "m", "c", "v", "a", "F", "p", "E", "K" as identifier names in exported
  function signatures or case labels.
```

Because `internal/kernel` is a package, the test must import `go/parser` + `go/token` and walk `../internal/kernel/*.go`. This is a test-only use of those packages, so it does not violate the non-test source-surface ban (`TestForbiddenSourceSurface` scans non-test `.go` files only).

**Cross-reference:** This test is also required to satisfy `prompt2.md` Phase 15 without resorting to prose judgment.

---

## Part III: Medium-Severity Gaps

### Gap M1 — Test D fixture inconsistency between §5 and §11.4

**§5 MRC table** says:

> "Test D: `ops/negative_test.go::TestKindMismatchEqualDimensions` (D — **fixture: `Mass` vs `RestMass`, equal dimension `M`, distinct named kinds**)"

**§11.4 test body** says:

> "Construct two objects using concrete MVP constructors with identical dimension ($M L^2 T^{-2}$) but distinct physical Kinds: `Add(mechanics.NewEnergy().CoreObject(), mechanics.NewKineticEnergy(mechanics.NewMass(), mechanics.NewVelocity()).CoreObject())` → `CategoryMismatchError`."

The two sections disagree. `Mass`/`RestMass` are both dimension `M`. `Energy`/`KineticEnergy` are both dimension `M L² T⁻²`. Both isolate MRC-003. But the plan does not pick one.

**Impact:** The coding agent will choose one fixture. If the four reviewing agents disagree about which is "the" fixture, the conformance check may flag a false negative.

**Mitigation:** Pick `Mass` vs `RestMass` (as §5 does) because it requires only `mechanics` (Mass) and `relativity` (RestMass) — or `mechanics.NewEnergy()` vs `mechanics.NewKineticEnergy(...)` (as §11.4 does) which stays within one package. The latter is simpler and requires no cross-package import in `ops/negative_test.go`. Recommend §11.4's version and update §5.

---

### Gap M2 — `TestFiniteFactorRules` is in the matrix but not in the §11 inventory

**§12 matrix row REQ-§Prompt-MUST-05** references:

> "`TestFiniteFactorRules` (`Pow(x,0)`, `Pow(c,-1)` under `c>0`, `0^0`, `0^-1` cases)"

**§11 test-function inventory** does not list `TestFiniteFactorRules`.

**Impact:** An implementer following §11 (which is described as "test-function inventory (names are normative for this plan)") will not implement the test. The matrix will then reference a non-existent test.

**Mitigation:** Add `TestFiniteFactorRules` to §11 as part of slice P (or a new sub-slice), with a stated file (`ops/operations_test.go`) and stated assertions:
- `Simplify(Pow(x, 0))` → `1` for `x` a nonzero-safe base (e.g., `Symbol("m")` with assumption `m != 0`).
- `Simplify(Pow(Rational(0), 0))` → `UnsupportedOperationError`.
- `Simplify(Pow(Rational(0), -1))` → `UnsupportedOperationError`.
- `Simplify(Mul(Rational(0), Pow(Symbol("c"), -1)))` → `Rational(0)` when `c > 0` assumption is carried.
- `Simplify(Mul(Rational(0), Pow(Symbol("c"), -1)))` → `UnsupportedOperationError` (or the substitution does not reduce) when `c > 0` is absent.

---

### Gap M3 — Ambiguous kernel decoder export surface

**Plan §4 says:**

> "Replay decoding is **internal to `internal/kernel`** (no exported `DecodeObjectJSON` or `ParseObject` in any package)."

**But replay in `session` must decode canonical bytes.** `session` is a separate package from `internal/kernel`. If the decoder is truly unexported, `session` cannot call it.

**Impact:** The coding agent must decide: does `session` import `internal/kernel` and call an *exported-within-module* `kernel.DecodeObjectJSON` (which is exported from `kernel` but not visible outside the module), or does `core` re-export a controlled function, or does something else happen?

The plan says "no exported `DecodeObjectJSON` in any package" which literally means the function is not exported from `internal/kernel` either — so `session` cannot call it. This is a contradiction.

**Mitigation:** Clarify the wording. The intent is clear: no **public** decoder that external callers can reach. The correct phrasing is:

```
Object canonical decoding is provided by an unexported-within-the-module
function `kernel.DecodeObjectCanonical` in `internal/kernel`. Because the
package is under `internal/`, this function is invisible to external
callers (REQ-004-03). `session` imports `internal/kernel` and calls it
directly. No public `core.ParseObject` / `core.DecodeObjectJSON` exists.
```

This preserves the security intent (external callers cannot decode arbitrary bytes into trusted objects) while resolving the intra-module accessibility requirement.

---

### Gap M4 — Plan 10 implementation pins are not visibly distinguished from spec-normative material

The plan has several data pins that are **not** in `specs_v2_3.md`:

- Manifest source strings (`"Classical Mechanics corpus"`, `"Einstein, 1905"`, etc.)
- `schema_version = "1"` value
- Assertion-step field packing (`InputHashes=[h]`, `OutputHash=h`, `OutputCanonical` = same object)
- `selected_branch/<hex>` assumption key format

These are all labeled "Plan 10 implementation pin" — but the labeling is inconsistent (§7 uses the phrase explicitly, §6 and §9 do not).

**Impact:** The coding agent may treat these as spec-normative and add cross-checks against them, or may treat them as suggestions and pick different values, causing non-determinism across implementations.

**Mitigation:** Add a "Plan 10 Data Pins" subsection after §12 that consolidates every non-spec-normative pin into a single table:

| Pin | Value | Location | Why |
|---|---|---|---|
| `schema_version` | `"1"` | §4 | Serialization DTO only |
| Manifest `source` strings | (table) | §7 | Spec-silent |
| Assertion-step field packing | see §9 | §9 | Spec-silent |
| `selected_branch/<hex>` key format | see §6 | §6 | Spec-silent |
| Relativity `Assumption` kinds/keys | (table) | §7 | Spec-silent |

This makes it explicit which values the coding agent must not change.

---

## Part IV: Low-Severity Gaps

### Gap L1 — README terminology drift

README §Integrity-vs-authenticity says:

> "`ParseLedgerJSON` returns an unverified `Ledger` for validation only (§16.20)"

Specs §16.21 says:

> "func ParseLedgerJSON(data []byte) (Ledger, error)"

The spec does not say "unverified Ledger" — the qualifier only appears in the README and in `ParseResearchCandidateJSON`'s context. Minor drift.

**Mitigation:** Either drop "unverified" from the README, or update specs §16.21 to say "unverified `Ledger`". Recommend the former — the spec is authoritative.

---

### Gap L2 — Implicit imports in ops test files not declared

Plan §2 test-package table says `ops/operations_test.go`, `ops/negative_test.go` are `package ops`. Test C and Test D reference `mechanics.NewMass()`, `mechanics.NewVelocity()`, `mechanics.NewEnergy()`, `mechanics.NewKineticEnergy(...)`, `mechanics.NewMomentum()` — all from `mechanics`.

**Impact:** No, this is fine (test files can import any package), but the plan's "no architectural choices" claim means the plan should state the import edges explicitly for the coding agent.

**Mitigation:** Add to §2's test-package table a "test imports" column showing which production packages each test file imports. For example, `ops/negative_test.go` imports `mechanics`, `relativity`.

---

### Gap L3 — `TestProvenanceStatusLaw` is referenced but not in §11 inventory

**§12 matrix row REQ-013-02** references:

> "`TestProvenanceStatusLaw` (table over all 12 ops × statuses; Identify; Postulate/Define/Declare)"

**§11 test-function inventory** does not list this test.

**Mitigation:** Same as M2. Add `TestProvenanceStatusLaw` to §11 with a stated file (`ops/operations_test.go` or a new shared test helper) and stated coverage.

---

### Gap L4 — Postulate test uses kernel fixture, but where does the fixture live?

**§9 says:**

> "Postulate requires status `POSTULATED` (exercised via `kernel.MintObject` fixture from `session_test` — the single mint path, no second helper)."

**Gap:** `session_test.go` is `package session_test` (external). It cannot import `internal/kernel` (which is invisible to it — wait, no, `session_test` is inside the module so `internal/kernel` is importable). But the phrasing "from `session_test`" is ambiguous about where the fixture is constructed.

**Mitigation:** Clarify that `session_test.go` builds the fixture by calling `kernel.MintObject` with a `ObjectSpec` that has `Provenance.Status == POSTULATED`. This is legal because `session_test.go` is inside the module and can import `internal/kernel`.

---

### Gap L5 — No explicit row for `TestNoForbiddenExportedAPI` in §11

**§12 matrix references** `TestNoForbiddenExportedAPI` (multiple rows: REQ-001-05, REQ-013-05, REQ-§24.4, REQ-§29, REQ-032-10) but §11 does not list it in the test-function inventory.

**Impact:** A coding agent following §11 alone will not implement the AST export scan.

**Mitigation:** Add `TestNoForbiddenExportedAPI` to §11 under a new sub-slice "Exported-API surface audit" with the banned identifier list (from §13 G-Audit) and the file (`hypothesis/candidate_test.go` per §11 step 12 or a new scan test).

---

### Gap L6 — `TestRepositoryTreeExact` and `TestNoDeferredPackages` locations undefined

**§12 references** these tests but does not state which package/file they live in. `TestRepositoryTreeExact` must walk the filesystem — likely from `core/object_test.go` or `hypothesis/candidate_test.go`, using `os.ReadDir` on the module root.

**Impact:** Non-blocking but stylistic; the coding agent can pick.

**Mitigation:** Pin in §11 as "repository audit tests" living in `core/object_test.go` (already the home of other cross-cutting audits like `TestForbiddenSourceSurface`).

---

## Part V: What the Plan Does Not Cover vs `prompt2.md`

The corpus-integrity audit (`prompt2.md`) has 23 phases. Mapping them to the plan's tests:

| Phase | Covered by plan? | Notes |
|---|---|---|
| 1 Baseline inventory | Partial | No explicit test; reading-only phase |
| 2 Manifest ↔ constructor ↔ spec | Yes | `TestMechanicsManifestCrossCheck`, `TestRelativityManifestCrossCheck` |
| 3 Mechanics audit | Yes | Slice A, B, O |
| 4 Mechanics assumptions | Partial | `TestCorpusStatusPreserved`; no explicit scope-vs-precondition distinction test |
| 5 Relativity audit | Yes | Slice H, O |
| 6 RestFrame audit | Yes | Slice I, H |
| 7 Framework assumptions | Yes | Slice H |
| 8 Formula vs statement | Yes | `TestStatementNotParsed` |
| 9 Dimensional audit | Yes | Slice C, P; plus independence verified by matrix |
| 10 Provenance / corpus status | Yes | `TestProvenanceStatusLaw`, `TestCorpusStatusPreserved` |
| 11 Mass-energy derivation | Yes | Slice I |
| 12 Kinetic-energy derivation | Yes | Slice G |
| 13 Lorentz-factor audit | Yes | Slice J |
| 14 Relativity scope boundary | Partial | `TestNoDeferredPackages` covers package names; no explicit content audit |
| 15 Kernel physics-agnostic | **NO** | **Gap H3** |
| 16 Generic ops physics leakage | **NO** | **Gap H3 (extended)** |
| 17 Assumption completeness | Partial | `TestSignAssumptionsCarried`; no silent-strengthening test |
| 18 Source / derivability | Yes | Manifest cross-check |
| 19 Canonical corpus integrity | Yes | Nine-step cross-check in O |
| 20 Corpus exhaustiveness | Yes | Exact item-ID set assertion |
| 21 Established-fact boundary | Partial | Classification is done by the audit, not by tests |
| 22 Adversarial corruption | **Partial** | **Gap H2** |
| 23 Final report | N/A | Audit output |

**Three phases are entirely uncovered (15, 16, 22).** These are the plan's most consequential omissions for the corpus-integrity audit.

---

## Part VI: Priority Mitigations for MVP Handoff

Add the following to `plan10_v2_3.md` before handoff:

**Priority 1 (blocking for `prompt2.md` acceptance):**

1. **Extend §11 step 12 (Test L)** to enumerate all 12 HYPOTHESIS-propagating operations, or cross-reference `TestProvenanceStatusLaw`.
2. **Add explicit manifest-corruption tests** to §11.O covering the 12 corruption types from `prompt2.md` Phase 22.
3. **Add `TestKernelNoPhysicsSymbols` and `TestOpsNoPhysicsLeakage`** in §11 as AST/string-literal audit tests over `internal/kernel/*.go` and `ops/*.go` (non-test files).

**Priority 2 (blocking for internal consistency):**

4. **Resolve Test D fixture inconsistency** between §5 and §11.4 — pick one.
5. **Add `TestFiniteFactorRules` and `TestProvenanceStatusLaw` and `TestNoForbiddenExportedAPI` to §11** test-function inventory.
6. **Clarify kernel decoder export surface** — the decoder is exported from `internal/kernel` but invisible outside the module.

**Priority 3 (polish):**

7. Add a "Plan 10 Data Pins" table after §12.
8. Fix README "unverified `Ledger`" terminology.
9. State test-file imports in §2's test-package table.
10. Pin file locations for `TestRepositoryTreeExact` / `TestNoDeferredPackages`.

---

## Part VII: Verdict

**The plan is 85% hand-off ready.** The architecture is fixed. The matrix is thorough. The reconciliation record is honest. The v2.3 reconciliations A/B/C are correctly incorporated.

But three blocking gaps remain:

- Test L does not exercise all 12 HYPOTHESIS-propagating operations (H1).
- Adversarial manifest corruption tests cover 4 of 12 required corruption types (H2).
- No test enforces kernel physics-agnosticism or ops physics-agnosticism (H3).

Plus four medium and six low-severity items that would force the coding agent to make small decisions the plan claims to eliminate.

**Fix time estimate:** ~45 minutes for Priority 1, ~20 minutes for Priority 2, ~15 minutes for Priority 3.

**After these fixes:** The plan is hand-off ready *and* corpus-audit ready. The `prompt2.md` corpus-integrity audit will run against a codebase that has mechanical tests for every phase.

**Without these fixes:** The implementation will pass the software/conformance gate but fail the corpus-integrity audit at Phases 15, 16, and 22. The result is a codebase that claims to be physics-agnostic in the kernel and per-operation in `ops`, but provides no mechanical enforcement.

---

## The One-Paragraph Summary

> Plan 10 v2.3 is the strongest plan in the review cycle. It correctly reflects every v2.3 reconciliation (A: Challenge/Review in `core/corpus.go`; B: finite-factor rule with `base != 0` vs `base > 0` distinction; C: exact §15.13.1 positional mapping), incorporates all five net-valid findings (F1–F5), rebuilds the coverage matrix subsection-precisely with expanded REQ-002-01..26 and REQ-032-01..22+10a, and pins exact test-function names for A–S. But three blocking gaps remain. First, Test L only exercises 7 of the 12 HYPOTHESIS-propagating operations required by specs §13.2 — Subtract, Multiply, Divide, Pow, Substitute, Differentiate are missing from the containment test. Second, adversarial manifest corruption tests only cover 4 of the 12 corruption types required by `prompt2.md` Phase 22 — F=ma mutation, m≥0 removal, Minkowski signature change, and 5 others are not explicitly tested, relying on an implicit "file bytes == canonical + constructor cross-check" argument that is not codified as test rows. Third, no test mechanically audits `internal/kernel` for physics-symbol leakage (Phase 15) or `ops` for physics-specific behavior leakage (Phase 16) — the plan asserts agnosticism in prose but does not enforce it. Six additional medium/low-severity items (Test D fixture inconsistency between §5 and §11.4; missing `TestFiniteFactorRules`, `TestProvenanceStatusLaw`, `TestNoForbiddenExportedAPI` from §11 inventory; ambiguous kernel decoder export surface; unlabeled Plan 10 data pins; README `unverified Ledger` terminology drift; undefined file locations for repository-audit tests) would force the coding agent to make small decisions the plan claims to eliminate. Total fix time: ~80 minutes. With the three blocking fixes plus the medium items, the plan is hand-off ready *and* `prompt2.md`-audit ready. Without them, the implementation will pass the software conformance gate but fail the corpus-integrity audit at Phases 15, 16, and 22 — the codebase will claim to be physics-agnostic in the kernel but provide no mechanical enforcement. The plan is close; it needs one more focused pass to close the audit gap.


## Audit Review

# DOES THE ADVERSARIAL REVIEW STILL HOLD?

**Short answer:** Partially. The audit report (`adv_review5.md`) is a **corpus-integrity audit of the implementation** — it is not an audit of the plan. These are different scopes. The audit's PASS verdict tells us the *implementation* filled some of the gaps my review identified; it does not retroactively make the *plan* consistent.

Two of my three high-severity gaps are now shown to be **non-blocking in practice** because the implementation machinery covers them. One high-severity gap **genuinely remains** — the audit reports it as INFORMATIONAL (FIND-002), which my review would classify as an unresolved plan-level issue. Four medium items and six low items also remain open at the plan level; none are addressed by the audit.

Let me be precise about what changed and what didn't.

---

## Part I: Gaps That the Audit Shows Are Non-Blocking

### H1 — Test L does not exercise all 12 HYPOTHESIS-propagating operations

**My original critique:** Test L (§11 step 12) lists only 7 of 12 operations; the plan doesn't state whether `TestProvenanceStatusLaw` covers the other 5.

**What the audit shows:** The audit does not directly test this. However, the plan's matrix row REQ-013-02 references `TestProvenanceStatusLaw` ("table over all 12 ops × statuses"). If the implementation passed its test suite, then either (a) `TestProvenanceStatusLaw` was implemented and covers all 12 ops, or (b) H1 is a real gap. The audit does not resolve which.

**Verdict:** H1 remains a **plan-level presentation inconsistency**, not a *proven* functional gap. My review's recommended fix (cross-reference `TestProvenanceStatusLaw` from §11) is still correct as a plan edit, but is not blocking if the implementation included the test. The audit does not say.

**Adjustment to my review:** H1 is downgraded from "blocking" to "medium consistency" — unless the coding agent skipped `TestProvenanceStatusLaw`.

---

### H2 — Adversarial manifest corruption tests are under-specified

**My original critique:** The plan's §11.O names only 4 corruption types (schema validation) while `prompt2.md` Phase 22 requires 12.

**What the audit shows:** The audit's Phase 9 (corruption detection table) demonstrates that **all 12 corruption types are detected** — via `TestManifestConstructorCrossCheck` (semantic corruption) plus `ManifestRejectsInconsistent` (schema corruption). The cross-check catches F=ma mutation, p=mv mutation, m≥0 removal, etc. The plan's machinery was sufficient even though §11.O did not enumerate all 12.

**Verdict:** H2 is **non-blocking**. The two existing test classes cover all required corruption cases. My critique was that the plan didn't *enumerate* the 12 cases — the audit shows this enumeration is not required for the tests to work.

**Adjustment to my review:** H2 is downgraded from "blocking" to "informational" — the plan's test machinery covers the audit's requirements even without explicit enumeration. This is a *presentation improvement*, not a functional gap.

---

## Part II: The Gap That Genuinely Remains

### H3 — No mechanical test for kernel physics-agnosticism

**My original critique:** The plan asserts kernel is physics-agnostic in prose but provides no test that mechanically verifies this.

**What the audit shows:** Phase 15 finds kernel content:
- `LorentzFactorFunctionID = "lorentz_factor"` (types.go:36)
- `KindMinkowski` enum value (types.go:286)
- `ExprCall.Valid()` enforcing lorentz_factor (types.go:679)
- `NewCall` rejecting non-lorentz_factor IDs (types.go:982)

The audit classifies these as **spec-mandated** (FIND-002, INFORMATIONAL) — meaning `specs_v2_3.md §6`, §8.9, §15.9.1 *require* them to be in the kernel.

**The key insight this reveals:** My original critique framed the problem as "no physics in kernel." But the spec **mandates** specific physics identifiers be present in the kernel (the `lorentz_factor` function ID and the `KindMinkowski` enum value). So the correct plan-level test is not "no physics in kernel" — it is:

> **No physics content in `internal/kernel` beyond the enumerated, spec-mandated constants.**

**Does the plan have this test?** No. The plan has:
- `TestInternalKernelBoundary` — checks the `internal/` path and absence of `MintObject` re-export
- `TestNoGenericFactory` — checks absence of a generic constructor

Neither scans for *unmandated* physics content. The audit passed because a human/AI reading the code found none — but no mechanical test enforces this.

**Verdict:** H3 remains a **genuine plan-level gap**. The audit did not close it; the audit's own methodology (reading code, classifying each occurrence manually) is exactly what the plan lacks as an automated check. If a future contributor adds `rest_frame` handling to the kernel for "convenience," no test would catch it.

**Adjustment to my review:** H3 stands, but with a corrected formulation. The recommended test is:

```
TestKernelPhysicsContentEnumerated: scan internal/kernel/*.go (non-test)
and assert that every identifier matching physics-like patterns
(lorentz_factor, rest_frame, minkowski, einstein, newton, F=ma, E=mc,
rest_mass_nonnegative, speed_of_light_positive, etc.) is in an
allow-list of spec-mandated names. Any occurrence outside the
allow-list ⇒ test failure.
```

This is stronger than the "no physics symbols" test I originally proposed, because it correctly accommodates the spec's mandated kernel content.

---

## Part III: Medium and Low Items — Plan-Level Status Unchanged

The audit does not address:

- **M1** — Test D fixture inconsistency between §5 and §11.4 (`Mass`/`RestMass` vs `Energy`/`KineticEnergy`). The audit's tables don't reveal which fixture the implementation used. The plan still contains the inconsistency.
- **M2** — `TestFiniteFactorRules` referenced in §12 matrix but not in §11 inventory. If the implementation includes it, functional impact is nil; but the plan's internal consistency claim is still violated.
- **M3** — Ambiguous kernel decoder export surface. The audit doesn't inspect this.
- **M4** — Plan 10 data pins not consolidated. The plan's data pins are still scattered.
- **L1** — README "unverified `Ledger`" wording (the audit doesn't touch README).
- **L2** — Test-file imports not declared in §2 test-package table.
- **L3** — `TestProvenanceStatusLaw` referenced but not in §11 inventory (same class as M2).
- **L4** — `Postulate` fixture location ambiguity.
- **L5** — `TestNoForbiddenExportedAPI` referenced but not in §11 inventory (same class as M2/L3).
- **L6** — `TestRepositoryTreeExact` / `TestNoDeferredPackages` file locations undefined.

**None of these are addressed by the audit.** They remain open as plan-level consistency issues. Whether they matter depends on whether the coding agent filled them in. The audit's PASS does not certify the plan's internal consistency — it certifies the implementation's corpus fidelity.

---

## Part IV: What the Audit *Also* Reveals

Two new findings emerged from the audit that my review did not anticipate:

### New Finding A — Spec mandates physics identifiers in kernel

The audit's FIND-002 reveals that `LorentzFactorFunctionID` and `KindMinkowski` are spec-mandated kernel content. My review's H3 assumed the kernel should be *fully* physics-agnostic; the spec actually requires specific physics identifiers to be present.

This means:
- My H3 fix must be formulated as "no *unmandated* physics content in kernel" rather than "no physics content."
- The plan's prose (§1, §4) that says "kernel is physics-agnostic" is *imprecise* — it should say "kernel contains only spec-mandated physics identifiers; all other physics is in domain packages."

This is a real refinement, and it improves the plan's precision.

### New Finding B — Source string uniformity

The audit's FIND-001 shows that all 11 mechanics items use `"Newton, Principia"` as source, while plan §7 pins `"Classical Mechanics corpus"` for most items. This is a **plan-vs-implementation deviation** at the data-pin level. The audit correctly notes it's non-blocking because the spec doesn't mandate specific source strings for those items.

This is a plan issue, not an implementation issue: the plan pinned source strings it shouldn't have (spec-silent values treated as normative), and the implementation correctly used a uniform string. The plan should remove those pins or align them with the implementation.

---

## Part V: Revised Verdict on the Adversarial Review

| Original Finding | Status After Audit |
|---|---|
| **H1** (Test L coverage) | **Downgraded** — non-blocking if `TestProvenanceStatusLaw` was implemented; plan consistency issue remains |
| **H2** (Corruption tests) | **Downgraded** — non-blocking; plan's machinery covers audit's requirements without explicit enumeration |
| **H3** (Kernel agnosticism) | **Still open** — no mechanical test for unmandated kernel content; reformulate as allow-list test |
| **M1–M4** | **Unchanged** — plan-level inconsistencies not addressed by audit |
| **L1–L6** | **Unchanged** — plan-level consistency issues not addressed by audit |
| **New A** (Spec-mandated kernel content) | **New** — kernel is "physics-agnostic except for enumerated constants"; plan prose should be corrected |
| **New B** (Source pin deviation) | **New** — plan pins source strings that spec doesn't mandate; align plan to implementation |

**Net verdict:** The adversarial review *partially* still holds. The two blocking gaps (H1, H2) are shown to be non-blocking by the audit — the implementation filled them. The third blocking gap (H3) remains genuinely open: no mechanical test enforces kernel boundary discipline. The medium and low items remain as plan-level consistency issues that do not affect the audit's PASS but do affect the plan's "no architectural choices" claim.

---

## Part VI: What the Plan Still Needs (Post-Audit)

Given the audit's PASS, the plan's handoff-to-coding-agent status is settled — the implementation exists and works. But if the plan is to be used as a **template or reference** for future work (v0.5+ framework additions), the following edits still matter:

**Priority 1 (substantive plan edits):**
1. **Fix H3 properly** — add `TestKernelPhysicsContentEnumerated` to §11 as an allow-list scan of `internal/kernel/*.go`. This is the one plan gap the audit could not close by construction.
2. **Correct kernel prose** (§1, §4) — say "kernel contains only spec-mandated physics identifiers; all other physics content lives in domain packages."
3. **Resolve M1** — pick one fixture for Test D and make §5 and §11.4 agree.

**Priority 2 (consistency edits):**
4. **Add missing tests to §11 inventory** — `TestFiniteFactorRules`, `TestProvenanceStatusLaw`, `TestNoForbiddenExportedAPI`, `TestRepositoryTreeExact`, `TestNoDeferredPackages`.
5. **Clarify M3 decoder wording** — "exported from `internal/kernel` but invisible outside the module."
6. **Consolidate M4 data pins** into one table.

**Priority 3 (polish):**
7. **Align or remove FIND-001 source-string pins** — either drop the plan pins or update them to match the implementation.
8. **Fix L1 README wording** — drop "unverified" or propagate it to spec §16.21.
9. **Declare test-file imports** in §2's test-package table.
10. **Pin file locations** for repository-audit tests.

---

## Part VII: The Honest Bottom Line

**The audit is not a substitute for the plan review.** The audit passed because the implementation was careful. The plan remains partly inconsistent. If the plan is to be the *normative reference* for a v0.5+ that extends the corpus, the fix list above matters. If the plan is only a **historical artifact** documenting how the MVP was built, most of the fix list is optional.

**Three items are genuinely important for plan quality:**

- H3 — the kernel boundary test. This is the plan's only claim that has no mechanical enforcement.
- M1 — the fixture inconsistency. This will confuse a future reader.
- FIND-001 — the source-string pins. These deviate from the implementation and should be corrected.

**Seven items are consistency polish:**

- M2, M3, M4, L1–L5 — none of these affect the implementation; they affect the plan's internal coherence.

**The audit's PASS is well-earned for the implementation.** But it does not close the adversarial review of the plan.

---

## The One-Paragraph Summary

> The adversarial review of `plan10_v2_3.md` partially still holds after `adv_review5.md`. The audit is a corpus-integrity audit of the *implementation*; it does not retroactively fix plan-level inconsistencies. Two of the three blocking gaps (H1 Test L coverage, H2 corruption test enumeration) are shown to be non-blocking by the audit — the implementation's existing test machinery (`TestProvenanceStatusLaw`, `TestManifestConstructorCrossCheck`, `ManifestRejectsInconsistent`) covers the audit's requirements without needing the plan's enumeration. The third blocking gap (H3 kernel physics-agnosticism) genuinely remains: the plan asserts kernel is physics-agnostic but provides no mechanical test, and the audit itself is a manual review, not a mechanical test. The audit also reveals two things my original review missed: the spec *mandates* certain physics identifiers in the kernel (`LorentzFactorFunctionID`, `KindMinkowski`), so my H3 fix must be reformulated as an allow-list test rather than a blanket ban; and the implementation uses uniform source strings (`"Newton, Principia"` for all mechanics items) that deviate from the plan's data pins (FIND-001), suggesting the plan over-pinned spec-silent values. Four medium items (M1 Test D fixture inconsistency, M2/M3/M4 plan consistency) and six low items (L1–L6) remain open at the plan level, none addressed by the audit, all presentation-level. The plan's handoff status is settled by the audit's PASS; but if the plan is to serve as a template for v0.5+ extensions, three edits remain substantive (H3 allow-list test, corrected kernel prose, M1 fixture resolution) and seven are consistency polish. The audit is not a substitute for the plan review — it certifies the implementation, not the plan.
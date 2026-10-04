I'm partway through the v2.1 spec (about 3,000 of 4,016 lines read) and checking whether the fixes closed the earlier blockers; reading the remainder now.**Verdict: one trust-boundary hole and about eight small spec patches away from a clean handoff.** v2.1 closes every blocker from my last review. It is not yet at "Open Spec Items: NONE", which is what your prompt promises the planner.

## Earlier blockers: closed

- **B1 (minting authority):** `internal/kernel.MintObject` plus the explicit threat model (§0.3) settles it.
- **B2 (import cycle):** replay lives in `session`, which imports `ops`, so there is no `core → ops` edge.
- **B3 (replay):** steps now retain `InputCanonicals`, `OutputCanonical`, `ParamsCanonical`, a genesis hash and a `StepEnvelope`.
- **B4 (string assumptions):** assumptions are structured (`ExprValue`), and denominator keys are hash-qualified.
- **B5 (`BranchSet`):** it is now a real expression node.
- **B6 (`E = mc²` trace):** §20 gives a golden trace, and the sign entailment is defined and works as written.
- **B7 (file count):** the tree is 39 files including `go.mod`, 38 without it.
- **Differentiation and limit tests:** they are now well-posed. `d(½mv²)/dv = mv` is checked against `Momentum`, and the `lorentz_factor` body is executed rather than looked up.

## Fix before handoff

**1. Deserialization is a public object factory (architectural).**
- `ParseResearchCandidateJSON` is public, and `Premises()` and `Hypothesis()` return `core.Object`. Crafted JSON with `"status":"DEFINED"` and arbitrary kind and dimension therefore yields a valid trusted-looking object to any external caller.
- Ledger validation doesn't stop this: assertion steps are only checked for self-consistency (§16.19 step 11), and §16.17 concedes that anyone can recompute the hashes.
- This breaks REQ-005-10 and the §0.3 guarantee.
- Smallest fix: candidate accessors return canonical JSON (or hashes) for premises and the hypothesis, not `core.Object`.

**2. The mandated golden traces can't run with the rules as written.**
- **Missing power rules:**
  - §9.6 and §9.7 have no `Pow(0,n)→0` and no `Pow(Rational,int)` evaluation.
  - Both `E = mc²` step 3 (`(0·c)²`) and the limit (`Pow(1,-1)`) need them.
- **Undefined predicate:** "valid finite symbolic factor" (§9.6). Suggested definition: no `Pow` with a negative exponent whose base isn't entailed nonzero.
- **Unclear normalization boundary:**
  - Constructor-time versus `Simplify`-time normalization isn't specified.
  - Golden step 2 shows `(0*c)^2` unsimplified after `Substitute`, so zero-annihilation must be `Simplify`-only.
  - Suggested rule: constructors do only flatten, rational-combine, sign normal form and sort. Everything in §9.6–9.8 is `Simplify`-only.

**3. The rest-frame assumption is never attached to the derivation.**
- §19.2 and the manifest's `derivable_from=[…, RestFrame]` say it should be, but the §20.1 sequence never merges `RestFrameAssumption()`.
- The derived `E = mc²` therefore silently drops its `p=0` condition, which REQ-001-04 forbids.
- Fix: have `ZeroThreeMomentum()` carry it.
- To avoid circularity, build the `Relation(eq, p, 0)` expression directly rather than via `Compare(p, ZeroThreeMomentum(), eq)`.

**4. Session semantics.**
- **`Seal` semantics for non-candidate sessions:**
  - `Seal` requires a `HYPOTHESIS` object in `DraftMetadata`, so a pure `E = mc²` session can never reach `Sealed`.
  - Either say so (such sessions end at `Concluded` and `Ledger.Validate()`) or make the hypothesis optional.
  - Acceptance test K ("visible after sealing") then needs a candidate.
- **Stale reference:** §16.11 cites "§18" for the candidate fields; it should be §26.
- **Framework-reference validation:** `session` can't import manifests, so "framework-reference integrity" (§26.8) must mean well-formedness only. State that; the anomaly cross-check belongs in the test.
- **`POSTULATED` producer:** no production constructor yields `POSTULATED`, so `Session.Postulate` is unreachable outside kernel fixtures. Either make one relativity constructor `POSTULATED` or state that it is fixture-tested only.
- **`SelectBranch` assumption key:** the key is unspecified (`energy_nonnegative` in §20.1 is a label). Suggested: `"selected_branch/"+HashExpr(constraint)`, mirroring §11.5.

**5. File hygiene.** This matters because the corrupted sections are the API-surface ones.
- **Literal `\n` escapes:** 39 lines contain them, inside code blocks in §§5.2.1, 7.0, 8.0, 8.2.1, 11.0, 12.0, 13.0, 15.13, 16.2.1, 16.21, 17.0, 18.0, 19.0, 26.0 and 26.10.
- **Duplicated lists:** the assumption-kind list in §11.1 and the corpus-status list in §13.4 appear twice.
- **Orphaned heading:** §26.9 has no body; §26.10 sits directly under it, and the text that belongs to §26.9 follows §26.10's code block.
- **Invalid manifest example:** the JSON example in §17.5 fails to parse; I ran it through `json.load`. It is missing a closing brace after `rhs`.

## Small items the planner can settle without risk

- **Type mismatch:** `ParentHashes()` returns `[]string`, while `NewProvenance` takes `[][32]byte`.
- **Shared provenance rule:** name a single place for the §13.2 and §13.5 propagation rule; `core/provenance.go` fits.
- **Kernel/core split:** one sentence would settle it. Methods on aliased types and all encoders and decoders live in `kernel`; `core` holds constructors and package-level functions.
- **Testing REQ-004-03:** an in-module `_test` package can import `internal/kernel`, so this can't be tested from inside the module. Use a path or AST assertion instead.
- **Underspecified details:**
  - `Commit` on an empty draft
  - the states in which `Validate()` is legal
  - the `OperationParams.Kind` values
- **Prompt sequencing:** §3 puts corpus packages before `session`, but the relativity derivation tests use `session`. Split those tests across steps.
- **Plan size:** the spec has 232 MUST lines, 177 of them unlabeled, so the §12 matrix will run to roughly 290 rows. Allow the matrix as a separate second deliverable, or the planner may truncate it.

Once items 1–5 are patched, "Open Spec Items: NONE" is a defensible expectation. I can draft the replacement text for each patch, including the repaired API sections, if you'd like.



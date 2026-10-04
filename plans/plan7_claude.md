I have the prompt in context, but the specs file isn't, so I'll read it from disk before reviewing.**Verdict: not ready to hand off.** The framing is strong: non-goals, error taxonomy, determinism rules, acceptance tests A–S, and honest MRC fallibility. But the spec has several structural defects, and the prompt's "no architectural decisions" promise depends on them being fixed first. A planner following this prompt would either invent architecture or produce an Open Spec Items section with dozens of entries.

## Blockers

**B1. MRC-001 can't be enforced with Go visibility as written.** `core.Object` has unexported fields, but `ops`, `mechanics`, `relativity` and `hypothesis` are different packages and must all create objects. That requires an exported minting function in `core`, which is exactly the "generic public object factory" §4.1 and §4.3 forbid. The spec has no threat model either. If agent code lives outside the module, Go's `internal/` gives real closure, but the layout forbids adding it. If the agent can add files to the module, no package trick helps. The spec must choose a threat model and a kernel location.

**B2. `Validate()` in `core/ledger.go` must re-execute `ops` steps, but `ops` imports `core`.** That is an import cycle. Fixes such as an injected executor, an `init()` registry, or moving replay elsewhere each have authority consequences. An injected executor lets a caller pass a permissive one. The spec must pick one.

**B3. Replay is impossible with the specified step record.** §15.5 stores only hashes, but §15.7 says to "reconstruct inputs from canonical hashes", and a hash can't be inverted. The step record is also missing:
- operation parameters (Pow exponent, Compare operator, SelectBranch constraint)
- the Identify justification (§24 and §15.7 step 7 need it)
- a genesis `previousHash`
- the meaning of `StepID`

Also undefined:
- a canonical JSON schema for `Object`, since §9 covers only `Expr`
- whether `Hash`/`Eq` operate on `Expr` or `Object`
- what "relevant physical metadata" in §8.1 means

**B4. String-valued assumptions contradict §4.4 and §7.** §4.4 says strings never carry relationships, yet assumptions use `p=0`, `E >= 0` and `denominator != 0`. Simplify must consult sign assumptions to decide `sqrt(x²)→x` (§7.7), and it can't do that without parsing those strings. Assumption values need a structured form.

The `Divide` rule is also broken. It uses the fixed key `denominator`, so either:
- the value is the literal string, and a second denominator's precondition silently dedupes away, or
- the value includes the expression, and two divisions raise a spurious `AssumptionConflictError`.

**B5. `BranchSet` has no representation.** `Object` has one `expr`, and the closed node set has no set or branch node. `Solve` can't return its result without violating the layout.

**B6. The canonical `E = mc²` derivation isn't derivable from the stated rules.**
- Simplify rules for zero-annihilation (`0·c`), `Pow(Mul, int)` distribution and like-base combination (`c·c²`) aren't specified.
- `sqrt(m²c⁴)→mc²` needs `m ≥ 0` and `c > 0`. Nothing introduces them, and the required sequence ends at `SelectBranch` with no final Simplify. `E ≥ 0` doesn't justify that step.
- `Substitute(p, 0)` and `E >= 0` need dimensioned zero objects, but there is no constructor for them. The same gap applies to the `½` in `K = ½mv²`.
- `Neg(x)` versus `Mul(-1, x)` and `Sqrt(x)` versus `Pow(x, ½)` are duplicate representations with no chosen normal form, so `Eq` is unreliable.

The spec needs a normative rewrite table plus a golden trace with intermediate canonical forms.

**B7. The spec violates its own size cap.** §3 lists 47 files (25 Go source, 17 tests, 2 JSON, 3 docs), against the §37.1 cap of 40. §3 also says "create or modify only" these files, so the plan can't trim or add any. Either raise the cap or cut 7 files, for example by merging `canonical.go` into `expr.go`, `convention` into `assumption`, and the relativity test files.

## High-severity gaps

- **Ill-posed tests.**
  - `Differentiate(Velocity(v), Time(t))` on independent symbols is 0, and there is no way to declare `v(t)`. Use something like `d(½at²)/dt` instead.
  - `Limit(lorentz_factor(v), v, 0)` uses an opaque `Call` with no `c` argument, so the "reduction check" is a hardcoded lookup checked against itself. §23.3 forbids this pattern for `E = mc²` but not here.
- **Negative tests that can't be run as written.**
  - "Fabrication" of an unexported struct literal is a compile error, not a test failure.
  - "Verify no promotion API exists" needs a method, such as a `go/parser` export-surface allowlist.
  - Assumption and convention conflict tests are hard to construct when constructors fix those fields. Which API lets a test create two conflicting objects?
- **Session semantics are undefined.**
  - `Postulate`, `Declare`, `Define`, `Conclude` and `Draft` have no defined meaning or signatures.
  - `ops.Identify(a, b, j)` takes no session but "records in the active session", which implies global state.
  - It's unclear how `Step` dispatches an operation by name.
- **Provenance rules are incomplete.**
  - There is no combination rule for `DEFINED`/`POSTULATED`/`IDENTIFIED` inputs.
  - "Monotonicity" (§15.7) has no defined order.
  - `APPROXIMATED` is unreachable with the MVP ops.
  - Whether Simplify yields `DERIVED` or preserves status is unclear ("may be DERIVED").
  - The `HYPOTHESIS` versus `IDENTIFIED` precedence in `Identify` is unstated.
  - MRC-008 names a failure, but contamination is silent tainting. When does `CandidateContainmentError` fire?
- **Undefined types and terms.**
  - "Trusted corpus artifact" (§25.6) doesn't exist as an API, so that assertion is vacuous.
  - `RecoveryClaim` and `AnomalyReference` have no fields.
  - Candidate concepts have no `Kind`; the enum lacks one.
  - "Compatible Kind" (MRC-003) is undefined, as is the dimension rule for `lt`/`gte`/`neq`.
  - `KineticEnergy` is in the mechanics list but not the `Kind` enum.
- **Manifest gaps.**
  - Items have no field naming the constructor to cross-check (§19.4).
  - The top-level schema has no anomalies array, though §18 requires records.
  - `domain` and `limits` element types are undefined.
  - `reduces_to` is undefined while the reduction test is deferred.
  - Who fills `corpus_status`? If the implementing agent does, that's curation by the AI.
- **Missing project basics.** The spec mentions "the existing repository" but never names it. No Go version, module path, or stdlib-only rule is given, and MRC is never expanded. The six-step tamper claim also needs a caveat. An unkeyed hash chain can be recomputed end to end, so it gives integrity, not authenticity, unless `LedgerHash` is anchored externally.
- **Thesis versus acceptance.** All 19 acceptance tests are scripted derivations. None exercises "an AI agent", so the MVP validates the substrate, not the thesis. That's fine, but it should be stated.

## Defects in the prompt itself

- "Plan Z", "DeepSeek" and "Qwen" (lines 60–62) appear nowhere in the spec, so the planner can't honor them. Cite spec sections instead.
- "MRC has one semantic source of truth" is never resolved. Which artifact is it?
- The coverage-matrix paragraph is duplicated (lines 237 and 239).
- "Do not write code" conflicts with "exact implementation of the types". Say that signatures and type sketches are allowed but function bodies are not.
- Section 15 demands justification for files over 40, but the spec forces 47, so the prompt should say how to handle that. §34 also says the repository must contain the matrix "in the implementation plan", which is muddled.
- There are 65 MUST/MUST NOT lines but no requirement IDs, so "every normative MUST" is not mechanically checkable. Add IDs like `REQ-§7.7-1`.
- Requiring tests for documentation contradicts §29, where the docs are auxiliary. Also, no package is designated to host such a test.
- "Exact sections" 1–15 plus "Open Spec Items" needs a fixed numbering (16), and there is no output size or format guidance.

## Suggested path

1. Patch the spec to v2.1 covering B1–B7. Resolve B1 and B2 together: choose the threat model, then the kernel location (a root `internal/` package is the cleanest Go-native option, at +1 file).
2. Add the normative rewrite table and golden `E = mc²` trace.
3. Fix the prompt items above.
4. Re-run this review. Only genuinely minor gaps should remain for Open Spec Items.

I can draft the v2.1 spec patch for B1–B7 and the high-severity items as concrete replacement text, with the cheapest coherent choice for each, if you'd like.
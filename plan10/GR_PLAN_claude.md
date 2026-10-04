I've read the GR plan and checked it against the kernel spec text, and I'm writing up the review now. The spec I had on hand is v2.1, not `specs_v2_3`. Anything marked (v2.1) below should be re-checked against v2_3 before you act on it.

## Verdict

As a governance charter it's strong. As an execution plan a coding agent could follow without making architectural decisions, it isn't ready. The main reason is that the plan never decides how GR scalars reach `ops` at all.

## Blockers

**B1. `ops` needs `core.Object`, but §6.4 says GR structures are not Objects.**
- §7 lists `Add/Multiply/Simplify/Substitute` as the GR contact surface, and those only accept Objects.
- Outside `phys`, the only ways to get one are:
  - the fixed constructors, whose symbols are fixed (`m,t,x,v,a,F,p,E,K` plus `c`);
  - `hypothesis.NewCandidateConcept` (v2.1), which takes any Expr and any Kind but forces `HYPOTHESIS` provenance.
- So every GR scalar that touches `ops` is `HYPOTHESIS`, and so is every `Identify` result.
- The plan never states this. It also means GR can never become trusted corpus from outside `phys`.
- Decide it explicitly and pre-register it as an expected finding. GR-3a also needs an `r` Object, which the probe as written doesn't have.

**B2. Trig representation contradicts itself.**
- §7 says GR-3b works over public `core.Expr` only. §9.2 says trig is GR-local.
- `NewCall` has one function ID (`lorentz_factor`) and no registry (v2.1). So `sin θ` can't be a node.
- That leaves two options. One is an opaque `Symbol("sinθ")`, which is exactly the C4 trap: the kernel differentiates it to 0 silently. The other is a parallel expression type, which is a CAS.
- The plan assumes trig is unavoidable, but that depends on the coordinate choice. With `x = cosθ`, the angular metric becomes `r²/(1−x²)` and `r²(1−x²)`, which is fully rational. Pre-register this as a GR-0 representation experiment, or the plan commits to a parallel CAS up front.

**B3. A required GR-local capability is missing from §9: rational-function normalization and zero-testing.**
- The v2.1 Simplify rule set has no expansion, like-term collection, or common-denominator cancellation.
- Showing `R_μν = 0` for Schwarzschild needs all three.
- This is the dominant cost of GR-5 through GR-8. A zero-test false positive is also the most dangerous silent-wrongness vector, and probe 3 doesn't cover it.

**B4. GR-local operations can't be replayed by `session`.**
- `ops.Apply` is closed over 12 IDs, so a local differentiate or simplify can only enter the ledger as a `Declare` assertion.
- Replay then checks consistency only. The hardest math ends up as unverified assertions, while §3 claims replay integrity is preserved.
- The plan also dropped the Canonicalization/Replay/Artifact impact section that plan11.1 §3 required.
- This is also the most plausible real Growth Gate evidence class: verifiable externally defined operations.

**B5. The experiment can't produce a growth verdict as designed.**
- Userland can reimplement almost anything, so nearly everything lands in `PACKAGE-SOLVABLE`.
- The gate requires "unlivable above the kernel". The only things that truly can't be done above it are authority-bearing: trusted mint, replay of new operations, and canonical hashing of new node types.
- Pre-register which observations would count as a `KERNEL-GROWTH CANDIDATE`, or `NO-GROWTH` is close to guaranteed.
- Two related problems:
  - `LEVEL-2-GROWTH` needs a named second consumer, and a GR-only run has none (SR or mechanics might qualify). Add a `LEVEL-2-CANDIDATE-PENDING` outcome.
  - The agent that hits an obstacle also classifies it. Require independent classification for anything other than `PACKAGE-SOLVABLE`.

## Significant

- **Classification contradicts itself.** §10 and §20.5 say each obstacle gets exactly one category, but §7 records GR-3a and GR-3b separately. Split "kernel encounter" from "userland resolution". `SPEC-INTENDED-BOUND` for negative exponents also overclaims if the spec text only defines non-negative exponents, which is the plan's own P9 distinction. Check whether v2_3's `Divide` emits `Pow(b,-1)`. If it does, the kernel produces expressions its own `Differentiate` rejects, and that should be recorded.
- **The kernel isn't "no metrics".**
  - The v2.1 Kind enum has `MinkowskiMetric`, `Spacetime` and `FourMomentum`.
  - The convention `metric.signature = -+++` already exists.
  - The plan misses two natural integration points: the `M→0` flat-limit check, and `ConventionSet` as a guard against the classic GR sign and signature errors.
- **GR-8 has no operational definition.**
  - Kernel `Limit` is substitution only.
  - Series expansion is deferred (REQ-002-23), and `APPROXIMATED` is reserved and unreachable.
  - `v << c` can't be written as a `Relation`.
  - §7's "`Limit` | Newtonian-limit step" is the wrong tool. Pre-register the expected classification and a feasible worked example.
  - Likewise `Solve` can't derive Schwarzschild from the vacuum equations; it can only verify them. Say so. §7's "single-pattern t²" description of `Solve` is also wrong.
- **Units are undecided.**
  - GR components carry no `Dimension`, so MRC-002 is never exercised.
  - Coordinate-basis components have index-dependent dimensions, which one scalar `Dimension` can't represent.
  - Choose between geometrized units and SI with explicit G and c. Add this as concern C9.
- **It isn't yet an execution plan.**
  - GR-0 through GR-8 are one-liners.
  - The ten per-pass attributes the plan promises are absent.
  - There are no worked examples or golden values.
  - Human veto points are missing.
  - There's no stopping rule.
  - Resource bounds are called "observational" but nothing is measured, and Riemann's 256 components will swell expressions.
- **Baseline pinning is missing.**
  - No version or commit pin for `phys`. A local `replace` or `go.work` follows the working tree, so PASS0 edits and mutation runs bleed in.
  - "No replace into kernel internals" is confused, since `replace` targets the module root.
  - PASS0 changes content in the frozen files, but the cited authorization (adv_review11) covers the earlier content. Add a new freeze record with hashes, and have GR evidence record the `phys` hash it ran against.
- **Mutation testing is mis-specified.**
  - Four of the listed mutants (the Pin B guard, `Substitute` kind equality, the `gte` constraint, the assumption-merge conflict rule) modify the kernel. They must run in a disposable worktree, with the frozen tree hash-verified afterwards.
  - The Pin A mutation changes the fixture, not production code. The load-bearing mutant is "`Identify` returns the operand's kind".

## Smaller items

- Probes 4 and 5 are the same test.
- Missing probes:
  - mixing charts;
  - treating Γ as a tensor;
  - misapplied symmetries;
  - sign and signature conventions;
  - unknown-function dependencies, which can't be represented, so the plan should restrict itself to explicit metrics.
- Probe 2 should be framed carefully. The kernel treats a symbol as constant in t exactly as specified, so the missing piece is a dependency declaration.
- `general_relativity` and `phys-gr` are used inconsistently, "Pin A: Pin A:" is duplicated, and the E=mc² firewall bullet has no GR motivation.
- Evidence-record fields 19–20 (independent review, human approval) have no defined process.
- How GR tensors are serialized and hashed for artifacts is unspecified.

## What's solid

- The separate-module boundary is real: the `internal/` closure holds for `github.com/PithomLabs/phys-gr`.
- The Pin A fixture correction is right.
- The self-contained C1–C8 ledger and the no-bypass rule are good.
- Treating checkpoints and `NO-GROWTH` as valid outcomes is sound.

## Order of fixes before handoff

1. Decide the Object bridge (B1).
2. Decide scalar and coordinate representation, with the `x = cosθ` experiment (B2).
3. Add the GR-local normalizer and the replay-boundary treatment (B3, B4).
4. Pre-register the falsifiers and the independent-classification rule (B5).
5. Fill the per-pass tables, including the Newtonian-limit worked example.
6. Pin the baseline and define the mutation harness.

I can draft the amendment text for these if you'd like.
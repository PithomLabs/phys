I'm re-checking the GR-plan claims against `specs_v2_3` and `plan10_v2_3`. The v2.3 spec confirms most of them so far. I'm looking at ops dispatch, Declare/Identify, and the file count next.I checked every kernel-dependent claim in my review against `specs_v2_3.md` and `plan10_v2_3.md` instead of the v2.1 text. Every blocker still holds, and the check changed three findings. The Newtonian limit is a spec-explicit bound, the negative-power classification is shakier than the plan assumes, and the plan's "42-file tree" conflicts with the spec's 39. The revised review follows.

## Verdict

As a governance charter it's strong. As an execution plan a coding agent could follow without making architectural decisions, it isn't ready. The plan never decides how GR scalars reach `ops` at all.

## Blockers

**B1. `ops` needs `core.Object`, but §6.4 says GR structures are not Objects.**
- §7 lists `Add/Multiply/Simplify/Substitute` as the GR contact surface, and those only accept Objects.
- Outside `phys`, the only ways to get one are the fixed constructors (fixed symbols) and `hypothesis.NewCandidateConcept` (§24.2).
  - That constructor takes any Expr and any Kind, but forces `HYPOTHESIS` provenance and `NONE` corpus status.
- Per §16.7 and MRC-006/008, `Identify` on a hypothesis operand stays `HYPOTHESIS`. So every GR scalar that touches `ops`, and every GR identification, is `HYPOTHESIS`.
- The plan never says this. It also means GR can never become trusted corpus from outside `phys`.
- Decide it explicitly and pre-register it as an expected finding. GR-3a also needs an `r` Object, which the probe as written lacks.

**B2. Trig representation contradicts itself.**
- §7 says GR-3b works over public `core.Expr` only. §9.2 says trig is GR-local.
- Per §8.9, `lorentz_factor` is the only function ID and no function registry exists.
- `NewCall` returns an error, but the spec doesn't state that it rejects other IDs. Make "does `NewCall("sin", …)` fail?" an explicit GR-0 probe instead of an assumption.
- If it fails, the options are an opaque `Symbol("sinθ")`, which is the C4 trap (the kernel differentiates it to 0 silently), or a parallel expression type, which is a CAS.
- Trig isn't forced by GR. With `x = cosθ` the angular metric is `r²/(1−x²)` and `r²(1−x²)`, which is fully rational and representable with negative integer `Pow`. Pre-register this as a GR-0 experiment, or the plan commits to a parallel CAS up front.

**B3. A required GR-local capability is missing from §9: rational-function normalization and zero-testing.**
- §9.4–§9.8 contain no expansion, no like-term collection, no common-denominator cancellation, and no `Pow` distribution.
- `Pow(Pow(x,a),b)` combines only for non-negative integers.
- Showing `R_μν = 0` for Schwarzschild needs all of these. This is the dominant cost of GR-5 through GR-8.
- A zero-test false positive is also the most dangerous silent-wrongness vector, and probe 3 doesn't cover it.

**B4. GR-local operations can't be replayed by `session`.**
- Per §15.13/§15.13.1, `ops.Apply` is closed over 12 IDs and "no plugin registry or runtime operation registration exists."
- A local differentiate or simplify can therefore enter the ledger only through `Declare`, and §16.5 and §16.19 step 11 replay assertions only as retained-object consistency checks.
- The hardest math becomes unverified assertions, while §3 claims replay integrity is preserved.
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

- **Classification contradicts itself.**
  - §10 and §20.5 say each obstacle gets exactly one category, but §7 records GR-3a and GR-3b separately. Split "kernel encounter" from "userland resolution".
  - `SPEC-INTENDED-BOUND` for negative exponents also overclaims. §15.8 defines the power rule only for non-negative integer `n` and lists `Sqrt`, `Call`, `Relation`, `BranchSet` and non-integer symbolic exponents as unsupported. Negative integer exponents appear in neither list.
  - So the bound is spec-silent and pin-defined, which is exactly the plan's own P9 distinction. Either add a `SPEC-SILENT / PIN-DEFINED` category or ratify the pin before classifying.
  - There's an asymmetry worth recording: §9.6 explicitly admits negative-integer `Pow` in expression trees, and the Lorentz body uses `v/c`. The kernel builds expressions its own `Differentiate` rejects. The spec doesn't pin how `Divide` represents quotients, so confirm in the implementation.
- **The kernel isn't "no metrics".**
  - The §6 Kind enum includes `Spacetime`, `MinkowskiMetric` and `FourMomentum`.
  - The convention `metric.signature = -+++` exists (§12).
  - The plan misses two natural integration points: the `M→0` flat-limit check, and `ConventionSet` as a guard against the classic GR sign and signature errors.
- **GR-8 is a spec-explicit bound, not a surprise.**
  - §15.9 says `Limit` is substitution only, with no series expansion and no asymptotic engine.
  - REQ-002-23 defers series expansion, `APPROXIMATED` is reserved until v0.5+, and the spec already defers the Newtonian kinetic-energy reduction for this reason.
  - `v << c` also can't be written as a `Relation`.
  - §7's "`Limit` | Newtonian-limit step" is the wrong tool. Pre-register `SPEC-INTENDED-BOUND` for the approximation step and a feasible worked example.
  - `Solve` accepts only `Relation(eq, Pow(Symbol(target),2), Expr)` (§15.11). The plan's "single-pattern t²" is wrong, and Schwarzschild can be verified but not derived.
- **Units are undecided.**
  - GR components carry no `Dimension`, so MRC-002 is never exercised.
  - Coordinate-basis components have index-dependent dimensions, which one scalar `Dimension` can't represent.
  - Choose between geometrized units and SI with explicit G and c. Add this as concern C9.
- **It isn't yet an execution plan.**
  - GR-0 through GR-8 are one-liners. The ten per-pass attributes the plan promises are absent.
  - There are no worked examples or golden values, no human veto points, and no stopping rule.
  - Resource bounds are called "observational" but nothing is measured, and Riemann's 256 components will swell expressions.
- **Baseline pinning and the file count don't reconcile.**
  - Both `specs_v2_3` (§3, §38) and `plan10_v2_3` pin exactly 39 files, with `TestRepositoryTreeExact (39 files)`.
  - The GR plan cites a 42-file tree authorized by adv_review11, which isn't in the uploads, so I can't verify it.
  - The plan must say which test now enforces 42 and how that squares with a byte-frozen `specs_v2_3`.
  - It also needs a version or commit pin for `phys`. A local `replace` follows the working tree, so PASS0 edits and mutation runs bleed in. "No replace into kernel internals" is confused, since `replace` targets the module root.
  - PASS0 changes frozen-file content, so it needs a new freeze record with hashes, and GR evidence should record the `phys` hash it ran against.
- **Mutation testing is mis-specified.**
  - Four mutants (the Pin B guard, `Substitute` kind equality, the `gte` constraint, the assumption-merge conflict rule) modify the kernel. They must run in a disposable worktree, with the frozen tree hash-verified afterwards.
  - The Pin A mutation changes the fixture, not production code. The load-bearing mutant is "`Identify` returns the operand's kind".

## Smaller items

- Probes 4 and 5 are the same test.
- Missing probes:
  - mixing charts;
  - treating Γ as a tensor;
  - misapplied symmetries;
  - sign and signature conventions;
  - unknown-function dependencies, which can't be represented, so the plan should restrict itself to explicit metrics.
- Probe 2 should be framed carefully. The kernel treats a symbol as constant in t exactly as §15.8 specifies, so the missing piece is a dependency declaration.
- `general_relativity` and `phys-gr` are used inconsistently, "Pin A: Pin A:" is duplicated, and the E=mc² firewall bullet has no GR motivation.
- Evidence-record fields 19–20 (independent review, human approval) have no defined process.
- How GR tensors are serialized and hashed for artifacts is unspecified.

## What's solid

- The separate-module boundary is real: Go's `internal/` rule keeps `github.com/PithomLabs/phys-gr` out of `internal/kernel`.
- The Pin A fixture correction is right. §16.7 constructs `Relation(eq, a.Expr(), b.Expr())`, so `KindRelation` is spec-grounded.
- The self-contained C1–C8 ledger and the no-bypass rule are good.
- Treating checkpoints and `NO-GROWTH` as valid outcomes is sound.

## Order of fixes before handoff

1. Decide the Object bridge (B1).
2. Decide scalar and coordinate representation, with the `x = cosθ` experiment (B2).
3. Add the GR-local normalizer and the replay-boundary treatment (B3, B4).
4. Pre-register the falsifiers and the independent-classification rule (B5).
5. Fill the per-pass tables, including a Newtonian-limit worked example that is pre-classified `SPEC-INTENDED-BOUND`.
6. Reconcile the 39/42 count, pin the baseline, and define the mutation harness.

I can draft the amendment text for these if you'd like.
# GR Implementation Plan — Consolidated Review and Locked Decision Register

Status: consolidation of all reviews of `GR_IMPLEMENTATION_PLAN.md`. Document-only. It edits no spec, kernel, manifest, or tree.

Purpose: leave **no architectural decision** to the coding agent. Every point where the plan was silent, contradictory, or ambiguous is pinned below as a numbered decision (`DL-nn`). The amended plan is the original plan plus §§5–10 of this document.

---

## 0. Sources and how they were read

| Tag | Source | Notes |
|---|---|---|
| **R-A** | My three review passes of `GR_IMPLEMENTATION_PLAN.md`; the third re-verified every kernel-dependent claim against `specs_v2_3.md` / `plan10_v2_3.md` | Findings tagged B1–B5 and S-items |
| **R-B** | `GR_PLAN_review.md`, a consolidation of six reviews (it cites claude, claude2, deepseek, qwen, z, gemini) | The six individual reviews were **not** uploaded. Attributions below come only from R-B's own citations. R-B's last line ends mid-sentence ("the target condition is finally achievable:"); nothing after it is assumed |
| **Plan** | `GR_IMPLEMENTATION_PLAN.md` (+ `plan11.md`, `plan11_1.md`, `plan11_1_overview.md`, `plan11_1_review.md`) | Reviewed target |
| **Spec** | `specs_v2_3.md` (read directly) | Authority for every kernel claim |
| **Not verifiable** | `adv_review11.md` (42-file re-freeze authorization) | Cited by Plan and R-B, not uploaded; handled as a PASS0 verification action (DL-02) |

---

## 1. Verdict

The architecture is sound and both review streams agree on it. The Plan was a governance charter, not an execution contract.

After consolidation there are **no remaining architectural blockers**, provided the decisions in §5 are adopted. Four items are not decisions but human actions or a known spec-text ambiguity (§11). None changes the architecture.

R-B resolves most of what R-A found and adds three things R-A had not decided (derivation-first Schwarzschild, concrete tensor structure, units). R-B also leaves **nineteen gaps** of its own (§4). Those gaps, not R-B's lock list, are what would still force the coding agent to improvise.

---

## 2. Reconciliation matrix

Disposition: **ADOPT** (as stated), **ADOPT+** (adopted with a required refinement), **REJECT**, **ADD** (new in this consolidation).

| # | Finding | Source | Disposition | Resolved in |
|---|---|---|---|---|
| 1 | §6 ("GR objects are not `core.Object`") contradicts §7 (uses `ops`, which needs `core.Object`) | R-A B1, R-B §2 | ADOPT+ (bridge consequences were missing) | DL-05, DL-06, DL-07 |
| 2 | Bridged objects are all `HYPOTHESIS`/`NONE`; GR can never be trusted corpus from outside `phys` | R-A B1 | ADD (R-B uses the bridge but never states the consequence) | DL-06, E1 |
| 3 | Trig/`UnknownFunction` cannot live in `core.Expr`; one coherent local algebra, no placeholder symbols | R-A B2, R-B §3–5 | ADOPT+ (R-A's `x = cosθ` idea becomes a measured census) | DL-10, DL-19 (GR-1/3 census) |
| 4 | Rational normalization and zero-testing are mandatory Level-1 machinery | R-A B3, R-B §8 | ADOPT+ (soundness, three-valued result, oracle, negative controls were missing) | DL-13, §8 |
| 5 | GR-local operations cannot enter `Session.Step` (closed 12 ops) | R-A B4, R-B §7 | ADOPT+ (R-B's `GRStep` omits `PreviousHash`, genesis, closed op set, link to kernel) | DL-17 |
| 6 | Experiment cannot open the gate as designed (unfalsifiable) | R-A B5 | ADD (R-B has no pre-registered falsifiers) | DL-22, DL-23, §9 |
| 7 | `LEVEL-2-GROWTH` unreachable in a GR-only run | R-A B5, R-B §17 | ADOPT+ (same flaw exists in the kernel outcome, see #8) | DL-22 |
| 8 | `KERNEL-GROWTH CANDIDATE` needs a proven Level-2 attempt, but the Plan forbids pre-building Level 2 | — | ADD (neither review caught it) | DL-22 |
| 9 | Independent classification for anything other than `PACKAGE-SOLVABLE` | R-A B5, R-B §18 | ADOPT | DL-24 |
| 10 | "Exactly one category per obstacle" contradicted by recording GR-3a and GR-3b | R-A | ADOPT (R-B keeps the contradiction) | DL-21 |
| 11 | `SPEC-INTENDED-BOUND` for negative exponents: spec is silent in the unsupported list | R-A | ADOPT+ (resolved by citation, see DL-21) | DL-21, E2 |
| 12 | Kernel is not "no metrics": Kind enum has `Spacetime`, `MinkowskiMetric`, `FourMomentum`; convention `metric.signature = -+++` exists | R-A | ADOPT (R-B uses the convention but not the wording fix) | DL-01, DL-09 |
| 13 | `phys.Limit` is the wrong tool for the Newtonian limit | R-A, R-B §12 | ADOPT+ (R-B has no bounded definition) | DL-16 |
| 14 | `Solve` accepts only `Relation(eq, Pow(Symbol,2), Expr)`; Schwarzschild cannot be solved by it | R-A, Spec §15.11 | ADOPT | DL-15 |
| 15 | Units undecided; tensor-component dimensions are index-dependent | R-A, R-B §13 | ADOPT+ (R-B's "dimension" assignment for bridged objects is undefined) | DL-07 |
| 16 | Per-pass attributes promised but absent; no worked examples or golden values | R-A, Plan §5 | ADD | §6 |
| 17 | Human veto points / stop rule missing | R-A | ADOPT+ | DL-26, DL-27 |
| 18 | Resource bounds "observational" but nothing is measured | R-A, R-B (A11) | ADOPT | DL-18, DL-26 |
| 19 | Baseline pinning, `replace` wording, PASS0 re-freeze record | R-A, R-B §21–22 | ADOPT+ | DL-02, DL-03 |
| 20 | Mutation harness must run on a disposable worktree; Pin A mutant must mutate production | R-A, R-B §21 | ADOPT+ | DL-04 |
| 21 | 39 vs 42 files | R-A | RESOLVED by R-B §21 (39 implementation + 3 authorized documentation files) | DL-02 |
| 22 | Probes 4 and 5 duplicate; missing probes (charts, Γ-as-tensor, symmetry, sign conventions, unknown functions) | R-A, R-B §14–16 | ADOPT | §8 |
| 23 | Field-function probe mislabels kernel behavior | R-A, R-B §16 | ADOPT | SW-1, SW-2 |
| 24 | Derivation-first Schwarzschild from unknown `A(r)`, `B(r)` | R-B §11 (qwen) | ADOPT+ (machine-certified script, not machine discovery; integration constant needs an identification step) | DL-15 |
| 25 | Concrete tensor/index structure before GR-2 | R-B §14 (qwen) | ADOPT+ | DL-14 |
| 26 | Kernel Contact Ledger instead of any ratio; no "Substrate Efficacy Ratio" | R-B §10 | ADOPT | DL-18, DL-25 |
| 27 | Reuse kernel `Simplify`/arithmetic where representable | R-B §9 | ADOPT+ (expected reuse is small; Spec §9.6–9.8 has no expansion or like-term collection) | DL-12 |
| 28 | Evidence-record fields 18 and 20 defined | R-B §19 | ADOPT | DL-24 |
| 29 | GR artifact hashing and canonical JSON | R-B §20 | ADOPT+ | DL-10, DL-17 |
| 30 | Remove "does not bypass `phys.Differentiate`" | R-B §6 | ADOPT | §10 |
| 31 | Remove "No replace into kernel internals" | R-A, R-B §22 | ADOPT | DL-03 |
| 32 | Wording defects (duplicate "Pin A:", `general_relativity` vs `phys-gr`, `t²` solver description, unmotivated E=mc² firewall bullet) | R-A | ADOPT | §10 |
| 33 | Substrate Efficacy Ratio; kernel tensor/trig/function nodes; placeholder `sin_theta`; weaken `Substitute`; extend `Differentiate`/`Limit` for GR; automatic Level-2 promotion | R-B (rejections) | ADOPT the rejections | §5 (DL-25) and §10 |

---

## 3. Conflicts between reviews and how they are resolved

**C-1. Trig: placeholder symbol vs local node vs rational chart.**
R-A listed two options (opaque `Symbol("sinθ")` or a parallel expression type) and offered a third: the rational chart `x = cosθ`, which makes `sin²θ = 1 − x²` and `dθ² = dx²/(1−x²)`. R-B rejects placeholders and wants real `Sin`/`Cos` nodes, and says the trig requirement "should remain in the workload rather than be evaded."

Resolution: R-B's `Sin`/`Cos` nodes in a local algebra are the **primary representation**. Placeholders stay forbidden. The rational chart is **not** an evasion. It is a legitimate second chart used for one purpose: a **bridgeability census** (DL-19, GR-1 and GR-3) that measures how much of the pipeline can reach the kernel in each chart. R-B's own rule (D12, "use the kernel whenever representable") cannot be evaluated honestly without this measurement. Without it, nearly every angular-sector component is unbridgeable and D12 reuse would be silently ~0.

**C-2. Verify-only vs derivation-first Schwarzschild.**
R-A recommended verify-only because unknown functions are not representable in `core.Expr`. R-B recommends derivation-first with local `UnknownFunction`.

Resolution: R-B wins, because the local algebra now represents unknown functions. But "derive" is only honest as a **certified derivation script** (DL-15): every step's candidate is supplied by the script and machine-certified by differentiation plus zero-test. The machine does not "discover" the solution. Explicit-metric verification is both the early golden oracle (GR-5/6) and the independent final check (GR-8).

**C-3. `phys.Limit` for the Newtonian limit (Plan §7) vs local reduction.**
Both reviews agree it is wrong. R-B says "local weak-field reduction" but does not define it. Defined in DL-16.

**C-4. Where kernel sessions appear.**
Neither review decides. Decided in DL-28: only the GR-8 identification segment uses `Session`, and only up to `Conclude`.

**C-5. "Default outcome: NO KERNEL GROWTH" vs "no modification presumed; burden on growth."**
Both appear in Plan §14. Resolved by one wording (DL-22): no kernel modification is presumed, all outcomes admissible, and the burden of proof is on the proposer of any promotion. These are compatible.

**C-6. Probe 2 ("bare symbol differentiates to 0 → SILENTLY-WRONG").**
R-A softened it; R-B reverses it. Resolution: SW-1 is `ACCEPTED-CORRECT` under frozen semantics (Spec §15.8, "Symbol → 1 if same symbol, otherwise 0"). `SILENTLY-WRONG` is reserved for results wrong under a *legitimate* API use (DL-21).

---

## 4. Defects in R-B's lock list (what R-B still leaves ambiguous)

R-B's D1–D18 are accepted as a base. Each gap below is closed by a decision in §5.

| Gap | Problem | Closed by |
|---|---|---|
| G-1 | Bridge consequence unstated: every bridged result and every `Identify` result is `HYPOTHESIS` (Spec §24.2, §16.7, MRC-006/008) | DL-06 |
| G-2 | With real `Sin`/`Cos`, nearly the entire angular sector is unbridgeable; D12 reuse is unmeasured | DL-19 census |
| G-3 | "Explicit Dimension" for bridged objects is undefined, and tensor components have index-dependent dimension | DL-07 |
| G-4 | `KERNEL-GROWTH-CANDIDATE` needs a Level-2 attempt that the Plan forbids pre-building | DL-22 |
| G-5 | No pre-registered falsifiers; the gate may be unopenable | DL-23 |
| G-6 | One obstacle receives two classifications (GR-3a, GR-3b) | DL-21 |
| G-7 | `GRStep` lacks previous-hash chain, genesis constant, closed operation set, and a link to kernel hashes | DL-17 |
| G-8 | "Bounded differential/algebraic reduction" is not specified; integration constant `k2` appears without an identification step | DL-15 |
| G-9 | `ZeroTest` soundness not required; a false-ZERO is the most dangerous failure; no negative controls | DL-13, §8 |
| G-10 | `UnknownFunction` node undefined (derivative orders, atoms) | DL-10, DL-12 |
| G-11 | `ToCore`/`FromCore` node set open; half-integer exponents and `Sqrt` round-trip unspecified | DL-05 |
| G-12 | Symbol namespace unspecified; GR symbols can silently collide with kernel fixed symbols (`m t x v a F p E K c`) | DL-08 |
| G-13 | "Curvature convention explicitly defined" is not defined | DL-09 |
| G-14 | Mutation list mutates fixtures instead of production code in the Pin A case | DL-04 |
| G-15 | No stop rule; capability inventory open-ended | DL-11, DL-26 |
| G-16 | Kernel Contact Ledger has no fields | DL-18 |
| G-17 | `go.work` for local use, but no check that the checkout equals the recorded baseline | DL-03 |
| G-18 | Per-pass tables, golden values, negative controls absent | §6 |
| G-19 | Local operation names could be confused with kernel operations (`Differentiate`, `Simplify`, `Limit`, `Solve`) | DL-12 |

---

## 5. Locked decision register (normative)

Format: each decision is final for this plan. Changing one requires human amendment (DL-26).

### A. Boundary, baseline, tooling

**DL-01 Module boundary.** `github.com/PithomLabs/phys` is the frozen Level-3 substrate. All GR work lives in the sibling module `github.com/PithomLabs/phys-gr`, which depends one-way on `phys` through the public `core`, `ops`, `session`, `hypothesis`, `mechanics`, and `relativity` packages only. `mechanics` and `relativity` inside `phys` are grandfathered. Go's `internal/` rule keeps `phys-gr` out of `internal/kernel`.
*Wording correction:* the kernel does not contain GR, tensor, or metric **semantics**. It does contain nominal placeholder Kinds (`Spacetime`, `MinkowskiMetric`, `FourMomentum`, Spec §6) with no components. Every sentence saying "no metrics" in the Plan is replaced by that statement.

**DL-02 Baseline.** Spec §3/§38 pin 39 implementation files. The operational freeze is 42 files (39 + 3 authorized documentation files), cited from `adv_review11` and not verified here.
PASS0 step 0, before any edit: verify that each PASS0 target file (`session/session_test.go`, the `ops` test file(s), `AGENTS.md`) is a member of the 42-file set, and record the inventory.
PASS0 closes with a **baseline record** (stored in `phys-gr/evidence/baseline/`) containing: `phys` git commit SHA; working-tree-clean flag; SHA-256 of `specs_v2_3.md`; SHA-256 of both manifests; SHA-256 of each modified PASS0 file; the exact 42-file inventory; Go toolchain version; results of `go build ./...`, `go vet ./...`, `go test ./... -count=1`. Every GR evidence record cites the baseline record ID.

**DL-03 Module mechanics.**
- Local development: an uncommitted `go.work` joins sibling `phys` and `phys-gr`.
- CI/production: `phys-gr/go.mod` requires `github.com/PithomLabs/phys` at the exact baseline version/pseudo-version, with `go.sum` committed.
- Forbidden: `replace` to a modified fork, vendoring `phys`, any `internal/kernel` import.
- Before every GR test run, `tools/verify_baseline.sh` compares the `phys` checkout commit and cleanliness with the baseline record. A mismatch halts the run and is never worked around.
- The phrase "No replace into kernel internals" is deleted.

**DL-04 Mutation protocol.**
- Never mutate the authoritative `phys` checkout. Use a disposable clone or worktree, apply the mutant, run tests, discard, then re-verify the authoritative checkout (`git status --porcelain` empty and the baseline hash check passes).
- The mutant list, all of which are *production* mutants:
  1. `session.Identify` returns the operand's Kind instead of `KindRelation` (kills Pin A).
  2. Remove the negative-exponent guard in `differentiateExpr` (kills Pin B and the GR-3a probe).
  3. Weaken `Substitute` kind equality (kills SW-5).
  4. Strip the `gte`/shape check in `SelectBranch`.
  5. Drop the assumption-merge conflict rule (kills SW-6).
  6. `symbolic.ZeroTest` mutated to always return `ZERO` (kills the negative controls and oracle).
  7. `tensor.Contract` mutated to skip the variance check (kills SW-7).
- Flipping a *fixture* (for example `Identify(Energy, E/c²)`) is not a mutation. It may be kept only as a vacuity check on the old fixture.

### B. Bridge and representation

**DL-05 Bridge.** Package `phys-gr/bridge`.
- `ToCore(GRExpr) (core.Expr, error)` is defined only on the subset `{Symbol, Rational, Add, Mul, Neg, Pow(integer exponent)}`. Any `Sin`, `Cos`, `UnknownFunction`, or non-integer exponent returns `ErrNotKernelRepresentable`. There is no placeholder conversion.
- `FromCore(core.Expr) (GRExpr, error)` is closed over the same subset. `Sqrt`, `Call`, `Relation`, and `BranchSet` return `ErrNotLocalRepresentable`.
- The kernel carrier is created only through `hypothesis.NewCandidateConcept(id, KindExpression, dimension, expr, assumptions, conventions)`.
- `GRExpr` has no half-integer exponents and no `Sqrt`.

**DL-06 Bridge consequences (pre-registered, not discoveries).** Every bridged object, and every `ops` result or `Session.Identify` result derived from one, carries `HYPOTHESIS` provenance and `NONE` corpus status (Spec §24.2, §16.7, MRC-006/008). Therefore:
1. GR can never produce trusted (`DEFINED`/`DERIVED`-from-corpus) objects from outside `phys`. This is by design; GR trusted corpus would require a GR package inside `phys`, which DL-01 forbids. Registry entry **E1**.
2. Mathematical facts (for example a normalization identity) are tagged `HYPOTHESIS`, because the kernel has no "mathematical, not physical" provenance. Registry entry **E12**, concern **C10**: observation only.

**DL-07 Kind, Dimension, ids.**
- Bridged objects use `KindExpression`, except the SI audit (GR-8d), which deliberately uses named kinds (`Mass`, `Position`) to exercise MRC-003.
- Geometrized phase (GR-0 … GR-8c): units `G = c = 1`; coordinates `x0, r` and parameter `M` are treated as pure numbers; bridged `Dimension` = `Dimensionless`. MRC-002 is therefore **vacuous** in this phase. Recorded as concern **C9** (a tensor component's dimension depends on its coordinate basis and cannot be a single scalar `Dimension`). Registry entry **E7**, classification `REPRESENTABLE-BUT-UNFAITHFUL`. Component dimensional metadata stays Level-1.
- SI phase (GR-8d): real dimensions through `core.Dimension` (Multiply/Divide/Pow). `G` has dimension `L³ M⁻¹ T⁻²`, built from the public Dimension API.
- Object `id` format: `gr/<chart>/<name>`.

**DL-08 Symbol namespace.**
- Kernel fixed symbols (reserved, never used for GR quantities): `m t x v a F p E K c`.
- GR symbols: `x0 r theta phi M G eps lambda k1 k2 A B` plus `UnknownFunction` names.
- Exception: in the SI audit (GR-8d), the symbol `c` is used **deliberately** and is the same entity as `relativity.NewSpeedOfLight()`.
- A uniqueness check in `bridge` rejects bridging two GR quantities under one symbol.
- SW-3 deliberately violates this at the kernel layer to measure C4.

**DL-09 Conventions (pinned).**
- Signature `-+++`, kernel key `metric.signature` (same key and value as the relativity package).
- Riemann (MTW): `R^ρ_{σμν} = ∂_μ Γ^ρ_{νσ} − ∂_ν Γ^ρ_{μσ} + Γ^ρ_{μλ}Γ^λ_{νσ} − Γ^ρ_{νλ}Γ^λ_{μσ}`; Ricci `R_{μν} = R^ρ_{μρν}`; scalar `R = g^{μν}R_{μν}`; Einstein `G_{μν} = R_{μν} − ½ R g_{μν}`.
- Coordinate order `(x0, r, theta, phi)`, 0-based indices.
- Kernel `ConventionSet` keys attached to every bridged object: `metric.signature=-+++`, `riemann.convention=MTW`, `coordinate.order=x0,r,theta,phi`.
- Convention probes use `metric.signature=+---` and `riemann.convention=Weinberg`.

### C. Local algebra and tensor engine

**DL-10 Single bounded local algebra** (`phys-gr/symbolic`).
- Node set: `Symbol`, `Rational`, `Add`, `Mul`, `Neg`, `Pow` (integer exponent only, negative allowed), `Sin(arg)`, `Cos(arg)`, `UnknownFunction{Name, Args[], DerivOrders[]}`.
- Not included: `Sqrt`, `Exp`, `Log`, half-integer exponents, any function not listed.
- Exact `math/big.Rat`; no floats; no `map[string]any` in canonical artifacts.
- Canonical JSON and SHA-256 hashing with child ordering mirroring Spec §9.3: node-kind ordinal, then canonical child hash, then bytes. Rationals encoded `"num/den"`.
- Node-kind ordinals are fixed in the order listed above.
- The artifact format is Level-1 and does not become a `phys` artifact.

**DL-11 Closed capability inventory.** `phys-gr` implements exactly: `Normalize`, `ZeroTest`, `Diff`, `Subst`, `bridge` (ToCore/FromCore), the single trig relation `Cos² → 1 − Sin²` (even powers only), `Truncate` (DL-16), reduction patterns RP-INT and RP-PROD (DL-15), and the tensor engine (DL-14). Any capability outside this list requires a human-approved amendment. It is never added by the implementing agent. This is the scope-creep guard against building a CAS.

**DL-12 Local operation naming and kernel reuse.**
- Local operations are named `Diff`, `Normalize`, `Subst`, `Truncate`, `Reduce`. They never reuse the kernel names `Differentiate`, `Simplify`, `Substitute`, `Limit`, `Solve`, and never wrap, shadow, or re-export them. `Diff` is "an independent Level-1 operation, not an extension or replacement of the kernel operation."
- Reuse rule: when a subexpression is `ToCore`-representable *and* the intended operation is one the kernel supports, the kernel operation is exercised and logged in the Kernel Contact Ledger. Otherwise the local operation runs.
- Expected reuse is small: Spec §9.6–9.8 contain no expansion, like-term collection, or common-denominator cancellation. That is evidence (E4), not a defect.

**DL-13 `Normalize` and `ZeroTest`.**
- `Normalize`: expand over a common denominator, collect like terms, apply `Cos² → 1 − Sin²`, produce a canonical form over atoms. Atoms are `Symbol`, `Sin(arg)`, `Cos(arg)`, `UnknownFunction` with each derivative order a distinct atom.
- `ZeroTest` returns one of `ZERO`, `NONZERO`, `UNDECIDED`, and never guesses. Rule: `ZERO` iff the expanded numerator over the common denominator normalizes to `0` (denominators assumed nonzero, recorded as assumptions).
- Soundness cross-check (required): an exact-rational point oracle in **test** code evaluates the expression at several rational points. `Sin`/`Cos` use Pythagorean points (for example `3/5, 4/5`); `UnknownFunction` is instantiated with concrete polynomial test functions using an independently written exact derivative. Any `ZERO` verdict that evaluates nonzero is a `SILENTLY-WRONG` finding against `phys-gr` itself.
- Required negative controls: §8 SW-10.
- Cancellation beyond exact division by denominator bases is out of scope.

**DL-14 Tensor engine** (`phys-gr/tensor`, `phys-gr/metric`, `phys-gr/connection`).
- `Tensor`: `Rank`, `IndexSlots[]`, `Components` (dense, fixed row-major order over indices `0..3`), `Symmetries`, `ChartID`.
- `IndexSlot`: slot number, variance (`covariant`|`contravariant`), index domain (`spacetime`).
- Declared symmetries are **verified** by `ZeroTest` at construction and stored with a certificate hash.
- Operations: `Contract(pairs)`, `RaiseIndex`, `LowerIndex`, `ApplyMetric`, `ApplyInverseMetric`. Contraction requires one covariant and one contravariant slot, the same chart, and the same domain.
- `Metric`: symmetric `Tensor` (0,2) plus inverse, verified on construction (`g·g⁻¹ = δ` via `ZeroTest`).
- `Connection` is a **distinct type** (Christoffel symbols), not a `Tensor`, and no tensor operation accepts it.
- No implicit Einstein summation. No coordinate transformations are implemented.
- Component keys are fixed tuples, never strings.

**DL-15 Schwarzschild derivation: certified script.**
Ansatz: `ds² = −A(r) dt² + B(r) dr² + r² dΩ²`, with `A`, `B` `UnknownFunction`s. Assumptions recorded: `A ≠ 0`, `B ≠ 0`, `r ≠ 0`, static, spherically symmetric.
Steps (each *certified*: the script supplies a candidate and the machine verifies it by `Diff` plus `ZeroTest`; the machine does not search):
1. Compute `R_tt`, `R_rr`, `R_θθ` through the pipeline and check them against the standard textbook forms in §6 (GR-6).
2. Certify `R_tt/A + R_rr/B = (AB)'/(r A B²)`; conclude `(AB)' = 0` ⇒ `AB = k1` (**RP-PROD**).
3. Boundary condition `A → 1, B → 1` as `r → ∞` is an **imposed assumption**, not computed (there is no limit engine) ⇒ `k1 = 1`. Registry entry E11.
4. Substitute `B = 1/A`: `R_θθ ∝ 1 − A − rA'` ⇒ `(rA)' = 1`.
5. **RP-INT** (candidate-and-certify): candidate `rA = r + k2`, certified by `Diff`. Hence `A = 1 + k2/r`, `B = 1/A`.
6. Certify all three vacuum equations `ZERO` for general `k2`.
7. Identify `k2` with `−2GM/c²` in the SI subtrack (DL-16, DL-28).
Explicit-metric verification is separate and independent (GR-8b).

**DL-16 Weak-field correspondence (bounded; replaces `phys.Limit`).**
- `Truncate(expr, λ, n)` extracts the coefficients of `λ⁰..λⁿ` from an expression that is a **polynomial in λ**, where `λ` scales `M` (or `k2`). A non-polynomial argument returns `ErrNotInClass`. No series expansion is implemented.
- Required results: exact flat check (`M → 0` or `k2 → 0` gives flat spherical metric); `g_00` correspondence; `a_r = −c² Γ^r_00` leading term `−GM/r²` with `k2 = −2GM/c²`.
- Required checks run through `Session.Identify` (DL-28) and kernel dimension audits:
  `GM/(c²r)` is dimensionless; `GM/r²` has acceleration dimension; `k2` and `GM/c²` both have dimension `L`.
- `APPROXIMATED` provenance is **not** activated. The approximation semantics stay Level-1. Registry entry E5.

**DL-17 GR trace (Level-1 replay).** Package `phys-gr/trace`. Mirrors the kernel ledger design (Spec §16) at Level 1.
- `GRStep` fields in order: `StepID`, `Index`, `OperationID`, `InputCanonicals`, `InputHashes`, `ParamsCanonical`, `OutputCanonical`, `OutputHash`, `KernelBridgeHash` (optional; the `core.Object` hash when the step used the kernel), `PreviousStepHash`, `CurrentStepHash`.
- Genesis previous hash: SHA-256 of the ASCII string `phys-gr-genesis-v1`. `CurrentStepHash` = SHA-256 over canonical JSON of all fields except itself.
- Closed `OperationID` set: `normalize`, `zero_test`, `diff`, `subst`, `bridge_to`, `bridge_from`, `truncate`, `reduce_rp_prod`, `reduce_rp_int`, `contract`, `raise_index`, `lower_index`, `assemble_christoffel`, `assemble_riemann`, `assemble_ricci`, `assemble_einstein`.
- `trace.Validate()` recomputes hashes and replays every step through the corresponding local operation. Replay detects any change that is not an exact valid operation output.
- The trace holds **no authority**: integrity only, not authenticity. It cannot be loaded into a `phys` `Session`.
- The two replay domains link only by hash equality of bridged object canonical forms.

**DL-18 Kernel Contact Ledger.** Descriptive only; no ratio, no threshold. One entry per contact with `phys`:
`pass`, `kernel_operation`, `inputs (object hashes)`, `outcome` (`succeeded` | `intentionally-unsupported` | `error: <typed error>`), `local_operation_invoked_instead`, `boundary_note`, `phys baseline ID`.
Observational metrics (never correctness thresholds): expression node counts before/after, `Normalize` term counts, wall-clock per pass.

### D. Experiment design

**DL-19 Pass structure.** GR-0 … GR-8 per §6, with the bridgeability census at GR-1 and GR-3.
The census records, for each of two charts — polar `(x0, r, θ, φ)` and rational `(x0, r, x = cosθ, φ)` — how many metric components (GR-1) and Christoffel components (GR-3) are `ToCore`-representable. It measures; it does not create a second pipeline.

**DL-20 Silent-wrongness suite.** Fixed list SW-1 … SW-13 in §8. `SILENTLY-WRONG` is reserved for a representation that is legitimate under the defined API and nevertheless yields a mathematically or semantically wrong *accepted* result. User misuse is not evidence against the kernel.

**DL-21 Failure taxonomy and recording.** Five categories unchanged. Each **obstacle** gets exactly one category for the *kernel encounter*. The *userland resolution* is a separate field of the same record. GR-3a and GR-3b are one record: encounter `SPEC-INTENDED-BOUND`, resolution `PACKAGE-SOLVABLE`.
Negative-exponent differentiation is classified `SPEC-INTENDED-BOUND` by citation: Spec §15.8 states the engine "supports only" the listed rules and defines the power rule for non-negative integer `n`, so negative `n` is excluded by construction. The explicit unsupported list is non-exhaustive, and PASS0 Pin B makes the exclusion executable.

**DL-22 Terminal outcomes (four, renamed).**
- `NO-GROWTH`: boundary correct; all needs met in userland.
- `LEVEL-2-CANDIDATE-PENDING`: an abstraction appears reusable; promotion needs a named second established consumer, a worked example, and human curation. It is never promoted inside this plan.
- `KERNEL-GROWTH-CANDIDATE-PENDING`: a full evidence record is opened. Fields 5 and 7 (the Level-2 attempt and its failure reason) are `NOT ATTEMPTED`, because Level 2 is not pre-built here. The record needs a follow-on provisional Level-2 attempt and human approval before it can advance. Status: `PROVISIONAL-INCOMPLETE`.
- `ESCALATE-TO-SPEC`: frozen-spec ambiguity discovered; human adjudication before any code change.
Wording lock: "No kernel modification is presumed. All Growth Gate outcomes remain admissible and must be determined from recorded evidence. The burden of proof lies with the proposer of any promotion."

**DL-23 Pre-registered growth-eligible classes.** Only these evidence classes can ever lead to `KERNEL-GROWTH-CANDIDATE-PENDING`:
- **GE-1** Replay-verifiability of externally defined operations (E6): the kernel is the only party that can attest replay.
- **GE-2** Authority to mint non-`HYPOTHESIS` objects for externally curated theory packages (E1).
- **GE-3** Cross-layer canonical identity (hashing) for new node kinds (E3, E6).
- **GE-4** Provenance/assumption propagation across the layer boundary, including mathematical-vs-physical conflation (E1, E12).
Not eligible: convenience, performance, missing trig/function/tensor support, unit handling, missing simplifier power (E4, E5, E10, E11). Each eligible record must answer the "unlivable above the kernel" question for a **theory-neutral** form of the capability.

**DL-24 Evidence record and review.**
- The 20 fields, enumerated in Plan §15, stand. Field 18 is `MinimalKernelChange` = affected package/file, affected exported identifiers, change kind (add API | alter invariant | add operation | add expression node | alter canonical encoding | alter replay semantics | alter MRC), no source code. Field 20 is `HumanApproval` = record ID, reviewer identifier, UTC timestamp, decision (`approved`|`rejected`|`provisional`), note. No cryptographic signing.
- Classification: `PACKAGE-SOLVABLE` may be classified by the implementing agent. Everything else needs an **independent classification** by a reviewer who did not author the pass (a numbered adversarial-review artifact), with disagreement resolution recorded. A record without human approval is `PROVISIONAL` and cannot trigger a kernel modification.

**DL-25 Rejected proposals (permanent).** Substrate Efficacy Ratio or any efficacy threshold; kernel tensor, trig, or unknown-function nodes; placeholder `sin_theta`-style symbols; weakening `Substitute`; extending `Differentiate` or `Limit` for GR; automatic Level-2 promotion; treating bare-symbol differentiation as a kernel bug.

**DL-26 Stop rule and resource policy.**
- A pass is `BLOCKED` when it cannot complete within the DL-11 capability inventory. Adding a capability is never the agent's decision: halt for human adjudication.
- A human sets a wall-clock limit `T_pass` in the run header. Exceeding it records `BLOCKED(resource)`, which is **observation-only** and is never growth evidence.
- Intermediate checkpoints are valid outcomes. Completed capability, remaining gap, and classification are recorded.

**DL-27 Human veto points.**
| Gate | After | Human decides |
|---|---|---|
| V0 | PASS0 | baseline record accepted; GR may start |
| V1 | amended plan | plan accepted (PASS1 complete) |
| V2 | GR-0 | API audit, conventions, symbol table, `T_pass`, independent reviewer named |
| V3 | GR-3 | boundary result (3a/3b record) accepted; census reviewed |
| V4 | GR-8 | derivation, verification, correspondence accepted |
| V5 | before any Growth Gate record advances | independent classification and approval |

**DL-28 Kernel sessions.** `session.Session` is used **only** for the GR-8d identification segment: `Draft → Declare (bridged objects) → Identify → Commit → Conclude`. `Seal` is not exercised and no `ResearchCandidate` is minted. Whether `Validate()` is legal before `Seal` is ambiguous in the frozen spec (§11, S-1). The segment calls `CanonicalJSON()`, and calls `Validate()` only if legal; otherwise it records the observed behavior and does not work around it.

**DL-29 PASS0 pins (final).**
- **Pin A** (`session/session_test.go`): operands `RestMass` and `E/c²`, where `E/c² = ops.Divide(relativity Energy, ops.Pow(SpeedOfLight, 2/1))` (dimension `M`). Run `Identify` inside a `Drafting` session. Assert `Kind == KindRelation`, `Expr == Relation(eq, …)`, `Provenance IDENTIFIED`, `CorpusStatus NONE`, `Dimension == Mass`, `Kind != RestMass`, operands unchanged.
- **Pin B** (`ops` tests): `Differentiate(Pow(x, negative integer))` → `UnsupportedOperationError`. Frozen-spec text (non-negative rule) is distinguished from the pin.
- **Pin C** (`AGENTS.md` §8 checklist, one line): named-kind compatibility protects additive operations only while a named kind is retained; equal-dimension `Expression + Expression` mixing is permitted by design; semantic preservation is higher-layer work. MRC-003 unchanged.
- **Observation O-1** (not a pin): record the frozen behavior of `Validate()` at `Committed`/`Concluded` and of `Seal` without candidate fields. Do not pin an unspecified behavior.
- Any non-conformance halts for human adjudication.

**DL-30 Wording corrections** — see §10.

---

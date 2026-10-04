# Plan 11.1 Revised — GR Implementation Plan (GR as Kernel Stress Test)

Version: v2 — 2026-10-02 (audit revision).
Status: execution contract. Planning document only.
No code, corpus, spec, or tree changes made by this document.

Changelog / audit trail:
- v1 (2026-10-01): initial 20-section execution contract
  (`plan10/GR_IMPLEMENTATION_PLAN.md`).
- v2 (2026-10-02): six-review reconciliation. Locks the GR ↔
  `core.Object` bridge (HYPOTHESIS wrapper), single bounded GR-local
  algebra (real `Sin`/`Cos`/`UnknownFunction`, no placeholders),
  GR-3a/GR-3b independent-op wording, GR-local replay trace,
  derivation-first Schwarzschild, GR-local Newtonian reduction (not
  `phys.Limit`), units/conventions pins, concrete tensor engine,
  corrected silent-wrongness semantics, `LEVEL-2-CANDIDATE-PENDING`
  terminal state, independent-review rule, precise evidence fields,
  GR artifact hashing, baseline/mutation protocol, `go.work`
  mechanics. No architecture reopened; implicit implementation
  choices pinned. Previous version preserved as
  `GR_IMPLEMENTATION_PLAN_20261001.md`.
- Snapshot of this version preserved as
  `GR_IMPLEMENTATION_PLAN_20261002.md`.

Lineage: `plan11.md` → `plan11_review.md` → `plan11.1.md` →
`plan11.1_review.md` → v1 → this v2 revision.
`plan11.1.md` is preserved. Frozen inputs: `plan10/specs_v2_3.md`
(byte-for-byte frozen, SHA-256 recorded at run time), both
`manifest.json` files (frozen), `AGENTS.md` workflow contract.

```
phys
 └── frozen kernel
     ├── no GR
     ├── no tensors
     ├── no metrics
     ├── no transcendental math
     └── no theory semantics

phys-gr
 └── Level-1 General Relativity
     ├── metric
     ├── tensors
     ├── indices
     ├── curvature
     ├── GR-local differentiation
     └── Schwarzschild / Newtonian limit

Growth Gate
 └── decides if anything deserves promotion
```

Consolidated locks (D1–D18, see Appendix A):
D1 kernel stays math-/theory-agnostic. D2 GR is Level 1 in sibling
module `github.com/PithomLabs/phys-gr`. D3 existing
`mechanics/`/`relativity/` inside `phys` are grandfathered (predate
this architecture); all NEW Level-1 packages follow the external-module
rule. D4 GR physical structures stay Level-1, not kernel Kinds. D5
kernel-compatible GR scalar subexpressions may be wrapped as
HYPOTHESIS `core.Object`s (`KindExpression`, explicit Dimension). D6
single bounded GR-local algebra. D7 algebra includes `Sin`, `Cos`,
`UnknownFunction`. D8 no opaque `sin_theta`/placeholder symbols. D9
GR-3a exercises frozen `phys.Differentiate` before GR-3b. D10
GR-local ops never modify/shadow/re-export kernel ops. D11 GR-local
math has its own deterministic replay/hash trace. D12 kernel
`Simplify`/arithmetic reused wherever genuinely representable. D13
tensor/index algebra is concrete Level-1 infrastructure. D14
Schwarzschild derived from unknown `A(r)`,`B(r)`, then independently
verified. D15 Newtonian correspondence is GR-local weak-field
reduction, not `phys.Limit`. D16 no substrate-efficacy ratio. D17
Level-2 promotion never automatic (`LEVEL-2-CANDIDATE-PENDING` only).
D18 kernel growth needs full evidence + independent review + human
approval.

---

## 1. Executive Summary

General Relativity (GR) is a hostile workload against the frozen `phys`
substrate, not a kernel feature. The question under test is:

> Can current kernel primitives support a serious theory package
> without being modified?

The experiment records three things separately for every encounter:
what `phys` can do, what `phys` deliberately refuses to do, and what
richer userland mathematics must provide. Success is boundary knowledge,
not "GR works." `NO-GROWTH` is an acceptable successful outcome.

Structure is gated:

```text
PASS0 — Conformance Pins (prerequisite gate, must be green)
        ↓
PASS1 — this GR Implementation Plan (assumes PASS0 baseline)
        ↓
GR Workload Execution (GR-0…GR-8 + adversarial probes)
        ↓
Growth Gate evaluation
```

PASS0 proves the kernel baseline. GR proves the boundary. The Growth
Gate decides whether anything deserves promotion.

v2 clarification: the plan previously left the GR ↔ `core.Object`
boundary implicit (§6 said GR structures are not `core.Object` while §7
expected GR to use object-based `ops`, which is unexecutable as
written since every `ops` function operates on `core.Object`). This
revision locks the bridge: Level-1 `GRExpr` structures plus a
temporary HYPOTHESIS `core.Object` wrapper used only when a frozen
kernel operation is deliberately exercised. No architectural choice
remains for the implementation agent (see §§6–7).

---

## 2. Architectural Principles

1. **Go-inspired layering.** Small robust core + rich composable
   ecosystem. The kernel is analogous to language/runtime primitives;
   theories are analogous to libraries.
2. **Kernel answers "what primitive capabilities are universally useful
   across theories?"** — exact rational arithmetic, immutable symbolic
   expressions, dimensional algebra, provenance/assumption tracking,
   canonical replay, deterministic transformations. The kernel must not
   answer "what mathematical objects belong to physics?"
3. **Kernel is math-agnostic and theory-agnostic**: minimal,
   composable, extensible, invariant-preserving. It is not a physics
   textbook.
4. **Algebra does not manufacture physical meaning.** Derived
   `Expression` results carry no nominal physical category by virtue
   of algebra alone.
5. **Do not make the kernel smarter because GR is complicated. Make
   userland richer first.** Promote into the kernel only when multiple
   theories independently prove the primitive is foundational,
   theory-neutral, and impossible to safely implement above the kernel
   boundary.
6. **Three levels, strict direction:**
   `LEVEL 1 (theory-local) → LEVEL 2 (shared math, human-curated)
   → LEVEL 3 (trusted kernel)`. No level modifies the level below it
   except through the Growth Gate, and Level 3 changes touch `phys`
   only via full gate evidence plus human approval.
7. **No ratio-based growth.** The Growth Gate concerns invariants, not
   quantity of code or share of operations performed by the kernel
   (see §18).

---

## 3. Frozen Baseline Assumptions

1. `specs_v2_3.md` is byte-for-byte frozen. This plan does not edit it.
2. Both theory `manifest.json` files (mechanics, relativity) are frozen.
3. Operational repository freeze: **39 frozen implementation files +
   3 authorized documentation files = exact 42 files.** Authorization
   record is `plan10/adv_review11.md` (`DOCUMENTATION PASS — 42-FILE
   RE-FREEZE AUTHORIZED`), confirmed by independent re-review (exact
   42-file membership + closed-world tree test). PASS0 adds content to
   existing authorized files only; exact 42-file membership is
   unchanged. PASS0 must verify before editing that its three target
   files are already members of that 42-file set.
4. `internal/kernel` owns trusted construction and invariants.
   External callers cannot import `internal/` — Go enforces this; the
   separate-module boundary (D2) is therefore a compiler-checked
   guarantee, not convention. All GR contact with `phys` is through
   the public `core` / `ops` / `session` facade.
5. `core.Expr` is immutable, closed to the MVP node set; the only MVP
   `Call` function id is `lorentz_factor`. `Pow` exponent is `*big.Rat`.
6. Trusted objects originate only through fixed domain constructors,
   hypothesis construction, pure operations, or session artifacts.
   There is no generic public object factory. GR cannot mint trusted
   kernel objects with arbitrary Kind/Dimension/Provenance/CorpusStatus.
   The sole GR bridge is the frozen hypothesis constructor, which
   accepts an arbitrary `Expr` plus explicit `Kind`/`Dimension` while
   forcibly assigning `HYPOTHESIS`/`NONE` (see §6).
7. Operations are the 12 supported IDs only; `identify` never via
   `ops.Apply`. `Simplify` cannot create `IDENTIFIED`, change corpus
   status, invoke `Identify`, or consult session state.
8. PASS0 is a mandatory prerequisite gate (see §16). The GR workload
   begins only after PASS0 pins are implemented and verified green with
   `go build ./...`, `go vet ./...`, `go test ./... -count=1`, and a
   post-PASS0 baseline record is produced (commit SHA,
   working-tree-clean status, SHA-256 of `specs_v2_3.md` + both
   manifests + each modified PASS0 file, exact 42-file inventory,
   build/vet/test results). Every GR evidence record cites the exact
   `phys` commit/baseline it ran against.

Non-negotiable invariants preserved by all GR work: immutability, mint
authority, no generic factory, exact rationals, canonicalization,
ordinal stability, MRC, provenance law + hypothesis contamination,
containment, session/ledger/replay/hash integrity.

---

## 4. Module Boundaries

```text
github.com/PithomLabs/phys      frozen Level-3 substrate
    internal/kernel, core, ops, session,
    mechanics (classical only, grandfathered),
    relativity (SR only, grandfathered),
    hypothesis (candidate space)

github.com/PithomLabs/phys-gr   Level-1 GR implementation
    depends one-way on phys (public API only)
    never imported by phys
```

Rules:

- All NEW theory work (including all of GR) occurs outside the `phys`
  repository tree as a sibling local module (e.g. `~/go/phys-gr` with
  its own `go.mod` requiring `github.com/PithomLabs/phys`).
  `mechanics/` and `relativity/` inside `phys` are grandfathered
  because they predate this architecture; they do not set precedent.
- Local development uses `go.work` connecting sibling `phys` +
  `phys-gr`. Production/CI pins an exact `phys` module
  version/commit. Forbidden: `replace` to a modified `phys` fork,
  vendored `phys`, direct `internal/kernel` import. (The phrase "no
  replace into kernel internals" is retired as imprecise — `replace`
  addresses modules, not packages.) Before any GR run, the `phys`
  checkout must equal the recorded clean baseline.
- Level-1 theory packages follow the same one-way dependency
  discipline as Level-2 libraries. A theory package failing due to
  missing kernel capability is Growth Gate evidence, not permission to
  bypass the boundary.
- Future Level-2 libraries are likewise separate Go modules (e.g.
  `physmath`) depending one-way on frozen `phys`; `phys` never depends
  back. GR implementation, GR-local tests, and GR-local evidence all
  live in `phys-gr` only. Shared Level-2 candidates live in a separate
  future module. Kernel changes are impossible without the Growth Gate.
  PASS0 is the sole pre-GR exception (prerequisite gate against
  already-existing `phys` test/documentation files).

---

## 5. GR Workload

Canonical static, spherically symmetric thread (locked). Linearized
gravity is reserved for later, not the entry.

```text
PASS0   Baseline Regression Verification (gate)
PASS1   GR Plan Verification (this document acknowledged)
GR-0    Representation + conventions + units + public API audit
GR-1    Static spherical metric (general ansatz with A(r), B(r))
GR-2    Inverse metric + concrete tensor/index engine
GR-3    Christoffel (3a kernel probe + 3b GR-local continuation)
GR-4    Covariant derivative
GR-5    Riemann curvature
GR-6    Ricci tensor + scalar curvature
GR-7    Einstein tensor + vacuum equations
GR-8    Schwarzschild derivation + independent verification
        + flat M→0 check + weak-field/Newtonian correspondence
Adversarial suite (independent, see §12)
Growth Gate
```

Each pass records: objective, files/packages, dependencies,
capabilities, worked examples, tests, expected failures, evidence,
level implicated, exit criteria — and where the substrate bends,
breaks, or silently misrepresents.

Derivation-first rule (D14): GR-8 PRIMARY derives Schwarzschild from
the general static spherical ansatz

```text
ds² = -A(r)dt² + B(r)dr² + r²dΩ²
```

with `A(r)`, `B(r)` unknown functions through Christoffel → Riemann
→ Ricci → vacuum equations → bounded differential/algebraic reduction
(problem-specific bounded local solver; NO general ODE solver).
SECONDARY target independently verifies the resulting closed-form
Schwarzschild metric. Verification alone is insufficient (it would not
exercise unknown-function differentiation or the reduction path).

GR-8 completion is the default target (full pipeline + `mechanics ↔
general_relativity` bridge without kernel promotion). Intermediate
stopping points are valid evidence checkpoints; record completed
capability + remaining gap + classification.

---

## 6. Representation Strategy

**Locked bridge (resolves the §6/§7 contradiction): GR mathematical /
physical structures remain Level-1 objects. A GR scalar expression is
bridged into a temporary generic `core.Object` ONLY when a frozen
kernel operation is deliberately being exercised.**

```text
              phys-gr Level 1
                    │
             GRExpr / GRScalar
                    │
       ┌────────────┴────────────┐
       │                         │
local-only expression      kernel-representable
       │                         │
       │                  HYPOTHESIS core.Object
       │                  Kind = Expression
       │                  explicit Dimension
       │                  explicit assumptions
       │                         │
       │                         ▼
       │                     phys / ops
       │
       └──── GR-local mathematics
```

Bridge for a kernel-compatible scalar:

```text
GRExpr
  ↓ ToCoreExpr()           (fails if Sin/Cos/UnknownFunction present)
core.Expr
  ↓ hypothesis.NewCandidateConcept(...)   (frozen constructor)
core.Object
  Kind = Expression, Provenance = HYPOTHESIS, CorpusStatus = NONE
```

This is explicitly permitted by the frozen hypothesis constructor and
is a BETTER stress test than bare `core.Expr`: it exercises
dimensional checks, MRC, assumptions, provenance contamination,
`Simplify`, `Compare`, etc., while preserving the trust boundary.

Critical restrictions:

1. A GR-local expression containing a Level-1-only construct (`Sin`,
   `Cos`, `UnknownFunction`) CANNOT be converted to `core.Expr`.
   No fake conversion through placeholder symbols.
2. One coherent bounded symbolic algebra — the agent must NOT create
   `core.Expr` + ad-hoc trig fragments. `phys-gr/symbolic/` defines:
   `Expr, Symbol, Rational, Add, Mul, Neg, Pow, Sin, Cos,
   UnknownFunction`, with structural conversion `GRExpr → core.Expr`
   only when the entire expression is kernel-representable. This is
   NOT a general CAS; it is bounded and sufficient for the workload.
   It owns unknown functions, trig, rational normalization, bounded
   differentiation/substitution, bounded first-order weak-field
   reduction, zero testing.
3. Real `Sin`/`Cos` nodes (D7–D8): `gφφ = r²·Sin(θ)²` with local
   rules `d Sin(x)/dx = Cos(x)`, `d Cos(x)/dx = -Sin(x)`. The
   reviewers' placeholder workaround (`sin_theta → Symbol` + secret
   derivative table) is REJECTED — it recreates the semantic
   ambiguity under measurement. Standard spherical Schwarzschild
   requires the angular structure; trig stays in the workload.
4. Unknown functions are genuine nodes (D7): `UnknownFunction("A", r)`,
   `UnknownFunction("B", r)` with `dA/dr = A'(r)`, `dB/dr = B'(r)`
   (+ higher derivatives as curvature requires). NEVER represent
   `A(r)` as `core.Symbol("A")` — the kernel correctly treats that as
   `r`-independent.
5. No new kernel expression nodes. Symbol identity stays syntactic at
   the kernel layer; namespacing/scoping are GR-local until Level-2
   evidence exists.

---

## 7. Kernel Interaction Points

Allowed contact surface (public facade only):

| Operation | Use | Constraint |
|---|---|---|
| Constructors (`NewPow`, `NewAdd`, …) | build kernel-representable scalar components | public API only |
| Read-only inspection | traverse `core.Expr` | never `internal/kernel` |
| `Add/Subtract/Multiply/Divide/Pow/Simplify/Substitute` | rational component algebra on HYPOTHESIS wrappers | pure, session-free; MRC enforced |
| `Differentiate` | GR-3a probe only | expect fail-closed on negative-integer powers |
| `Compare` | build relation artifacts | proves nothing by itself |
| `Solve` | single-pattern `t²` only | `BranchSet ±Sqrt` |
| `SelectBranch` | explicit `gte` constraint | recorded as structured assumption |
| `Session.Identify` | explicit identification | `KindRelation` only, never retyping |
| `Limit` | NOT used for GR-8 (see below) | Lorentz-factor body only; no asymptotic engine |

Forbidden: new `Kind`s, generic factories, kernel expression nodes,
`internal/` imports, assumption/provenance bypass, silent kernel
extension, expanding `ops` with GR-local operation IDs.

Differentiation boundary (D9–D10):

```text
GR-3a — Kernel boundary probe:
    kernel-representable Pow(r,-1) via HYPOTHESIS core.Object
    → phys.Differentiate(...)
    → UnsupportedOperationError
    → classified SPEC-INTENDED-BOUND. Record kernel boundary.

GR-3b — Level-1 continuation (only after GR-3a recorded):
    phys-gr/symbolic differentiator runs. This is an INDEPENDENT
    Level-1 operation, not an extension, replacement, shadow, or
    re-export of phys.Differentiate. Classified PACKAGE-SOLVABLE.
    Must not be represented as kernel differentiation capability.
```

(The old "does not bypass phys.Differentiate" phrasing is retired as
contradictory — GR-3b necessarily takes over the mathematical task
after GR-3a fails; independence is the correct relation.)

Reuse rule (D12): use the kernel whenever the relevant subexpression
is genuinely representable in `core.Expr`; use GR-local rules only
for absent capabilities. Example: GR-local differentiator emits a
kernel-compatible term → bridge to HYPOTHESIS wrapper →
`phys.Simplify` → reintegrate into `GRExpr`. This measures actual
kernel reuse instead of forcing everything through it or abandoning it.

Session rule: `phys.Session` records kernel-level operations and
trusted derivation events only (12-op vocabulary stays closed — a
GR-local differentiator can never enter `Session.Step` as a new
operation ID, and `ops` is NOT expanded to admit it). `phys-gr`
keeps its own deterministic replayable derivation trace:

```text
GRStep: StepID / OperationID / InputCanonical / ParamsCanonical
        / OutputCanonical / CurrentHash   (SHA-256 over canonical JSON)
```

Kernel-compatible scalar bridges additionally record the
corresponding `core.Object` hash. Two replay domains: Level-3 `phys`
Session replay; Level-1 `phys-gr` mathematical replay. If external
operations' exclusion from the trusted session ever proves to be a
generic theory-independent invariant problem, that is Growth Gate
evidence; otherwise it remains a deliberate boundary.

Newtonian-limit rule (D15): GR-8 does NOT use `phys.Limit` (frozen
`Limit` is direct-substitution + fixed Lorentz body, not an
infinity/asymptotic engine). GR-8 uses GR-local weak-field reduction
+ exact correspondence extraction + kernel dimension audit where
possible. Two checkpoints: exact flat limit (`M → 0` → flat spherical
metric); weak-field correspondence with `ε = GM/(c²r)`, bounded
first-order reduction recovering `Φ = -GM/r` (preferably also
`a_r = -GM/r²`). Approximation/regime semantics stay Level 1; the
frozen `APPROXIMATED` provenance capability is NOT activated.

---

## 8. Level-1 GR Package Design

Sibling module `github.com/PithomLabs/phys-gr` (one-way dep on `phys`):

```text
phys-gr/
    go.mod / go.work (local) 
    coordinates/    chart, symbols, assumptions/regime, conventions
    metric/         gμν + inverse, contraction gμν·gνρ=δμρ
    index/          variance, raising/lowering, contraction rules
    connection/     Christoffel assembly (uses symbolic/)
    curvature/      Riemann, Ricci, scalar, Einstein
    vacuum/         Rμν=0 structured relations
    solution/       Schwarzschild derivation + verification
    limit/          weak-field / Newtonian correspondence
    symbolic/
        expr.go            GRExpr algebra (Add/Mul/Neg/Pow/Sin/Cos/UnknownFunction)
        differentiate.go   GR-local bounded differentiator
        normalize.go       Normalize + ZeroTest (common-denominator,
                           rational-factor/like-term collection,
                           negative-power normalization, Sin²+Cos²→1)
        bridge.go          ToCoreExpr + HYPOTHESIS wrapping (refuses
                           Sin/Cos/UnknownFunction subexpressions)
        weakfield.go       bounded first-order ε reduction
    evidence/
        failure_log.md     per-obstacle classified records
        growth_records/    Growth Evidence Records
        kernel_contact.md  descriptive Kernel Contact Ledger (no ratio)
    adversarial/
        silent_wrongness_test.go
```

Dependency rule: `phys-gr → phys` only. `phys` never imports
`phys-gr`. GR-local math never imports `internal/kernel`.

Tensor engine (D13, concrete before GR-2 — "emergent" alone is too
vague to implement):

```text
Tensor: Rank / IndexSlots[] / Components / Symmetries / ChartID
IndexSlot: slot number / variance (covariant|contravariant)
           / index domain
Component keys: fixed tuples, not arbitrary strings.
Ops: Contract / RaiseIndex / LowerIndex / ApplyMetric /
     ApplyInverseMetric — enforcing same chart, compatible domains,
     valid contraction pairing, variance correctness, declared
     symmetries.
```

Explicit: **Christoffel symbols are connection coefficients, not
tensors** (adversarial probe).

Conventions (pinned GR contract): coordinate order `(x⁰, r, θ, φ)`;
metric signature `-+++` (matches existing SR package baseline);
GR curvature convention explicitly defined in `phys-gr` (kernel
acquires no GR curvature semantics). Convention-mismatch probes:
wrong signature, wrong Riemann sign, mixed convention sets, mixed
charts. Where scalar GR objects bridge into `core.Object`, attach
compatible `ConventionSet` metadata for frozen conflict detection.

Units (pinned): tensor calculation in geometrized units `G = c = 1`;
GR-8 physical-correspondence subtrack restores SI for scalar
Newtonian quantities. SI audit explicitly tests `GM/(c²r) →
dimensionless` and `GM/r² → acceleration` via HYPOTHESIS carriers
where kernel-representable. Coordinate-component dimensions stay
GR-local metadata (component dimensionality is
index/coordinate-basis dependent — never forced into kernel scalar
semantics).

---

## 9. GR-Local Mathematics Required

Owned by `phys-gr`, never kernel:

1. **Bounded differentiator** (negative-integer/rational terms,
   `Sin`/`Cos` rules, `UnknownFunction` `A'(r)`/`B'(r)` + higher
   derivatives as needed); independent Level-1 op (§7).
2. **Trigonometric nodes + identities** (`r²dθ² + r²Sin²θdφ²`;
   `Sin²+Cos² → 1` minimum). MVP stays closed.
3. **Normalize + ZeroTest** (mandatory Level-1 machinery —
   dominant practical cost of curvature stages): common-denominator,
   rational-factor/like-term collection, constant-factor and
   negative-power normalization, zero detection. Bounded; NOT a
   general CAS. Purpose: make `Rμν = 0` testable without pretending
   `phys.Simplify` is a GR CAS.
4. **Index/tensor engine** (§8), inverse-metric/contraction checks,
   curvature assembly, assumption/regime bookkeeping.

None of this implies kernel expansion.

---

## 10. Failure Taxonomy

Every obstacle maps to exactly one category, logged with requirement,
minimal example, current vs expected result, category rationale,
Level-1/Level-2 workarounds attempted, and Level-3 implication.
No wishlist entries.

```text
SPEC-INTENDED-BOUND
    MVP deliberately does not support it.
PACKAGE-SOLVABLE
    general_relativity handles it locally (implementation agent
    may classify this alone).
REPRESENTABLE-BUT-UNFAITHFUL
    Substrate encodes it but loses important semantics.
    (independent reviewer required — see §15.)
UNREPRESENTABLE
    Substrate cannot faithfully encode it.
    (independent reviewer required.)
SILENTLY-WRONG
    A representation legitimate under the defined Level-1 API
    nevertheless produces a mathematically wrong accepted result.
    (independent reviewer required.)
```

Classification authority: `PACKAGE-SOLVABLE` may be classified by the
implementation agent; anything else requires independent review
(agent classification + independent classification + disagreement
resolution) before Level-3 consideration (anti-confirmation-bias rule).

---

## 11. C1-C8 Concern Ledger

Self-contained; frozen routing.

| ID | Concern | Description | Classification | Routing |
|---|---|---|---|---|
| C1 | Derived-kind/type decay | `m·c²` and `r×F` both → `Expression M L²T⁻²`; `Add` permitted, category lost | Real limitation | Level-2 semantic/ascription investigation; never kernel, never weaken `Substitute` |
| C2 | Derived-expression ascription | No path from derived `Expression` (e.g. `E/c²`) to nominal kind via algebra; `Identify` does not retype | Real limitation | Level-2 semantic ascription; `Substitute` stays strict |
| C3 | Assumption semantic inconsistency | Keyed-merge `(Kind,Key)` conflict only; `v<<c` vs `v=0.99c` can coexist unless structurally keyed | Real limitation | Level-2 constraint/entailment layer |
| C4 | Symbol/entity identity collision | `Symbol("m")` syntactic; electron vs proton mass ambiguous; `Differentiate`/`Substitute` use string equality | Real limitation | Level-2 namespaced symbol/entity layer |
| C5 | Transcendental/general-function capability | MVP has no `sin/exp/ln`, no symbolic exponents (`Pow` takes `*big.Rat`; single `Call` id) | Not MVP defect | Future extension; GR trig stays Level-1 |
| C6 | Closed nominal Kind ontology | New first-class `Kind` requires curation/spec evolution; ordinals part of `mrc-v0.4` | Closed by design | Do not build subtype hierarchy in kernel |
| C7 | Identify result-kind regression pin | `Identify` must return `KindRelation`, never retype operand | Conformance pin | PASS0 Pin A (see §16) |
| C8 | Negative-integer differentiation pin | `Differentiate` supports non-negative integer exponents only; negative fail closed | Conformance pin | PASS0 Pin B (see §16) |

C1–C4 → architectural/userland concerns. C5 → future extension.
C6 → closed by design. C7–C8 → conformance pins, not growth findings.

---

## 12. Silent-Wrongness Testing Strategy

Mandatory independent adversarial suite alongside (not inside) the GR
workload. The dangerous outcome is `accepted + wrong`, not
`unsupported`. `SILENTLY-WRONG` is reserved for representations
legitimate under the defined Level-1 API that nevertheless produce a
mathematically wrong accepted result — user misuse is not kernel
evidence.

Corrected field-function semantics:

```text
Probe A: core Symbol("f") differentiated wrt r → 0
         → ACCEPTED-CORRECT (bare symbol really is syntactically
         constant under frozen semantics — NOT a kernel bug).
Probe B: GR UnknownFunction("f", r) → f'(r)
         → expected Level-1 behavior.
Probe C: malformed/ambiguous GR-local dependency representation
         → must be rejected/diagnosed locally.
```

Remaining probes (each asserts expected vs observed):

1. Symbol ambiguity (electron vs proton `Symbol("m")` mixing).
2. Invalid contraction/raising, mixed charts, Γ-treated-as-tensor,
   symmetry errors, sign-convention errors in GR-local structures —
   must not pass silently.
3. Dimension-valid semantic errors (energy-vs-torque mixing; same
   dimension wrong category substitution stays rejected or flagged
   unfaithful).
4. Type-decay addition (`E_expr + τ_expr`).
5. Inconsistent-assumption coexistence.
6. Convention mismatches (wrong signature / Riemann sign / mixed sets).

---

## 13. Level-1 → Level-2 Promotion Rules

A Level-1 abstraction becomes a Level-2 shared-library candidate only
with all of:

```text
shared abstraction
+ named second established consumer (not hypothetical)
+ worked second-theory example
+ demonstrated common (shared) semantics
+ assumption / provenance non-leakage
+ no forbidden framework-specific state leaking across consumers
```

"GR and QM might both need tensors" is insufficient. Level-2 home is a
separate Go module; promotion never touches the 42-file tree and never
happens silently inside this plan (see §18: GR yields at most
`LEVEL-2-CANDIDATE-PENDING`). No-bypass rule applies.

---

## 14. Level-2 → Level-3 Growth Gate

Level-3 (kernel primitive) requires the full Growth Evidence Record
(§15) plus proof of all of: generic, theory-neutral, non-semantic,
foundational, reusable across frameworks, invariant-bearing, and
unlivable above the kernel (provisional Level-2 attempt built without
changing kernel semantics, documented failure to preserve the
invariant). Burden of proof is on growth.

Wording lock: "No kernel modification is presumed. All Growth Gate
outcomes remain admissible and must be determined from recorded
evidence."

---

## 15. Evidence Record Format

Every `UNREPRESENTABLE` / `REPRESENTABLE-BUT-UNFAITHFUL` /
`SILENTLY-WRONG` obstacle that survives L1 and provisional-L2
attempts gets one record with these enumerated fields:

1. capability; 2. minimal counterexample; 3. failure category (§10);
   4. L1 attempt; 5. L2 attempt; 6. L1 failure reason; 7. L2 failure
   reason; 8. genericity; 9. foundationality; 10. second consumer
   (named); 11. second-consumer worked example; 12. non-goal collision
   check; 13. silent-wrongness analysis; 14. artifact compatibility
   (canonical form / hash — GR artifacts §20 + `core.Object` hashes
   for bridges); 15. MRC impact; 16. replay impact (both replay
   domains); 17. migration impact;
   18. `MinimalKernelChange`: affected package/file; affected exported
   identifier(s) if any; change kind (add API / alter invariant / add
   operation / add expression node / alter canonical encoding / alter
   replay semantics / alter MRC); NO source code;
   19. independent review reference (agent classification +
   independent classification + disagreement resolution — mandatory
   for non-`PACKAGE-SOLVABLE`);
   20. `HumanApproval`: record ID; reviewer name/identifier; UTC
   timestamp; decision (`approved` / `rejected` / `provisional`);
   decision note. No cryptographic signing required by the MVP. A
   record without human approval remains PROVISIONAL and cannot
   trigger kernel modification.

---

## 16. Test Strategy

**PASS0 — Conformance Pins (prerequisite gate, executes first).**
Any pin failure blocks; never worked around; halt for human
adjudication.

- **Pin A** (`session/session_test.go`): Pin A: `Identify(RestMass,
  E/c²)` → `KindRelation`, dimension `M`, provenance `IDENTIFIED`,
  corpus status `NONE`; must not produce `RestMass`. (Corrected
  fixture: `Energy` vs `E/c²` would fail dimension-compatibility
  first and test nothing.) Also assert `Expr` is `Relation(eq, …)`.
  Grounding: `session/session.go` already mints `KindRelation` —
  test pin only, zero production diff expected.
- **Pin B** (`ops` tests): `Differentiate(Pow(x, negative-integer))`
  → `UnsupportedOperationError`. Do not expand differentiation.
  Distinguish frozen-spec text (non-negative integer rule) from the
  conformance pin (explicit fail-closed). Test pin only.
- **Pin C** (`AGENTS.md` §8 checklist one-line caveat): named-kind
  compatibility protects additive operations only while a named kind
  is retained; equal-dimension `Expression + Expression` mixing is
  permitted by design; semantic preservation is higher-layer work.
  Do not change MRC-003. No file-count change.
- Pre-edit check: all three PASS0 target files already members of the
  exact 42-file set. Gate: `go build ./...`, `go vet ./...`,
  `go test ./... -count=1` green. Emit post-PASS0 baseline record
  (§3.8) for all GR evidence to cite.

**PASS1 / GR workload tests:** per-pass worked examples with golden
states (E=mc² chain intact, never a GR premise), GR-3a/GR-3b pair,
trig-ownership test, contraction tests, vacuum-equation and
Schwarzschild golden components (derivation + verification), flat
`M→0` and Newtonian-correspondence checks, Kernel Contact Ledger
entries (descriptive — which kernel op invoked / succeeded /
intentionally unsupported / which local op ran / where the boundary
lies; NO ratio, NO threshold).

---

## 17. Mutation Testing

Negative/mutation probes to confirm fail-closed behavior is real.
Protocol: NEVER mutate the authoritative `phys` checkout. Use a
disposable clone/worktree, run mutant + tests, discard mutant,
re-verify authoritative checkout clean. (No corpus-manifest
re-hashing merely because tests changed.)

- Flip Pin A fixture back to `Identify(Energy, E/c²)` → must fail
  dimension check (proves corrected fixture non-vacuous).
- Remove Pin B negative-exponent guard → GR-3a probe must fail
  loudly (proves pin load-bearing).
- Weaken `Substitute` kind-equality → type-decay suite must catch
  masquerading nominal quantities.
- Strip `gte` branch constraint from `SelectBranch` → branch tests
  must reject.
- Drop assumption-merge conflict rule → contamination tests must
  reject.
- Resource/AST bounds are observational instrumentation only, never
  mathematical correctness thresholds.

---

## 18. Expected Outcomes

Four admissible terminal states (evidence-determined). NOTE the
Level-2 rename (D17): GR-only execution cannot establish the named
second consumer, so automatic `LEVEL-2-GROWTH` is impossible:

```text
NO-GROWTH
    Kernel boundary correct; all GR needs met in userland.
    Acceptable success.

LEVEL-2-CANDIDATE-PENDING
    GR exposed an apparently reusable abstraction, but promotion
    is blocked until a named second established theory provides a
    worked example and the human curator approves a separate shared
    module. If mechanics/relativity is genuinely demonstrated as
    second consumer during this work, record the evidence — but do
    NOT silently promote inside this plan.

KERNEL-GROWTH-CANDIDATE
    Full evidence record + proven L2 failure to preserve a
    generic invariant + independent review + human decision.

ESCALATE-TO-SPEC
    Frozen-spec ambiguity discovered; human adjudication
    required before any code change.
```

No kernel modification is presumed. Prevent three false readings:
"kernel passed GR" (false — userland did work), "kernel failed,
therefore grow it" (premature), "Level 1 worked around it, therefore
kernel had no limitation" (also false).

---

## 19. Explicit Non-Goals

- No GR objects, tensors, metrics, manifolds, connections, curvature,
  Hilbert spaces, wavefunctions, operators, fields, spinors as kernel
  kinds. No kernel trig/unknown-function nodes.
- No semantic subtyping, entity identity, SMT solving, general
  entailment, general calculus, or tensor algebra in
  `internal/kernel`.
- No `semantic/`, `symbols/`, `constraints/`, `tensor/`, `manifold/`,
  `calculus/`, `hilbert/` pre-building before GR generates evidence.
- No transcendentals, symbolic exponents, rational-function
  differentiation, general ODE solving, general CAS, or asymptotic
  `Limit` expansion to make GR easier.
- No placeholder `sin_theta`/function-name symbols with secret
  derivative tables.
- No weakening of `Substitute`, `Identify` authority, MRC, provenance
  law, or closed-world tree. No `ops` expansion for GR-local ops.
- No substrate-efficacy ratio or growth threshold.
- No stored corpus result (notably `MassEnergyRelation`) used as a
  derivation premise (firewall violation).
- No empirical-validity claims from formal consistency; no corpus
  status assigned by computation; no `APPROXIMATED` activation by
  stealth.

---

## 20. Definition of Done

1. PASS0 pins implemented, green, baseline record emitted (§§3, 16).
2. `phys-gr` sibling module executes GR-0…GR-8 (or principled
   checkpoint with gap classification) without modifying `phys`.
3. GR-3a kernel failure + GR-3b independent local continuation both
   recorded with correct classifications.
4. Trig/unknown-function ownership demonstrated (real GR-local `Sin`,
   `Cos`, `UnknownFunction`; kernel `Call` unchanged; no
   placeholders).
5. C1–C8 ledger entries evidenced; failure log classifies every
   obstacle exactly once (§10) with Kernel Contact Ledger entries.
6. Silent-wrongness suite executed with corrected semantics (§12);
   every `SILENTLY-WRONG` has an evidence record or a Level-1/Level-2
   containment argument.
7. Schwarzschild derived from `A(r)`,`B(r)` + independently verified;
   flat `M→0` + weak-field correspondence (`Φ = -GM/r`) demonstrated
   via GR-local reduction.
8. Tensor engine concrete with ChartID/variance/symmetry enforcement;
   conventions pinned; GR artifact hashes recorded (§8 + Appendix C).
9. Promotion decisions follow §§13–14 with enumerated records;
   independent review present for all non-`PACKAGE-SOLVABLE`; Level 2
   never bypasses the kernel, Level 3 never presumed.
10. Final report states: what `phys` can do, what it deliberately
    refuses, what userland provides — i.e. **boundary learned.
    NO-GROWTH is an acceptable successful outcome.**

Final decision principle: do not make the kernel smarter because GR
is complicated. Make userland richer first. Only promote something
into the kernel when multiple theories independently prove that the
primitive itself is foundational, theory-neutral, and impossible to
safely implement above the kernel boundary.

---

## Appendix A — Consolidated Decisions D1–D18

D1 kernel math-/theory-agnostic. D2 GR Level-1 sibling module. D3
grandfathered `mechanics/`/`relativity/`. D4 GR structures not kernel
Kinds. D5 HYPOTHESIS bridge for kernel-compatible scalars. D6 single
bounded GR-local algebra. D7 `Sin`/`Cos`/`UnknownFunction` genuine
nodes. D8 no placeholders. D9 GR-3a before GR-3b. D10 GR-local ops
independent of kernel ops. D11 GR-local replay/hash trace. D12 kernel
reuse where representable. D13 concrete tensor engine. D14
derivation-first Schwarzschild. D15 GR-local Newtonian reduction. D16
no efficacy ratio. D17 `LEVEL-2-CANDIDATE-PENDING` only. D18 growth
needs evidence + review + approval.

## Appendix B — Explicitly Rejected

❌ substrate efficacy ratio/threshold ❌ kernel tensor types ❌ kernel
trig nodes ❌ kernel unknown-function nodes ❌ placeholder
`sin_theta` ❌ weakening `Substitute` ❌ extending
`phys.Differentiate` ❌ expanding `phys.Limit` ❌ automatic Level-2
promotion ❌ bare-symbol differentiation as kernel silent-wrongness bug
❌ `ops` expansion for GR-local ops ❌ general ODE solver / general CAS.

## Appendix C — GR Artifact Hashing

`phys-gr` defines canonical JSON + SHA-256 for: `GRExpr`, `Index`,
`Tensor`, `Metric`, `Connection`, `Curvature`, GR derivation steps
(`GRStep`). Level-1 format; never part of `phys` canonical artifacts.
Audit chain: GR local artifact → canonical JSON → SHA-256 → evidence
record. Kernel-compatible scalar bridges additionally retain their
`core.Object` hashes.

---

*Document-only artifact. Modifies no specs, kernel, manifests, corpus,
or frozen tree. Implementation agents execute it without reopening
architectural decisions; any spec ambiguity escalates to human review.*

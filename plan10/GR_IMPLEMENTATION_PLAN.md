# Plan 11.1 Revised — GR Implementation Plan (GR as Kernel Stress Test)

Status: execution contract. Planning document only.
No code, corpus, spec, or tree changes made by this document.

Lineage: `plan11.md` → `plan11_review.md` → `plan11.1.md` →
`plan11.1_review.md` → this document.
This document supersedes nothing historical; `plan11.1.md` is preserved.
Frozen inputs: `plan10/specs_v2_3.md` (byte-for-byte frozen),
both `manifest.json` files (frozen), `AGENTS.md` workflow contract.

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

---

## 3. Frozen Baseline Assumptions

1. `specs_v2_3.md` is byte-for-byte frozen. This plan does not edit it.
2. Both theory `manifest.json` files (mechanics, relativity) are frozen.
3. The 42-file closed-world tree is frozen. Authorization record is
   `plan10/adv_review11.md` (`DOCUMENTATION PASS — 42-FILE RE-FREEZE
   AUTHORIZED`). PASS0 adds content to existing authorized files only;
   exact 42-file membership is unchanged.
4. `internal/kernel` owns trusted construction and invariants.
   External callers cannot import `internal/`; all GR contact with
   `phys` is through the public `core` / `ops` / `session` facade.
5. `core.Expr` is immutable, closed to the MVP node set; the only MVP
   `Call` function id is `lorentz_factor`. `Pow` exponent is `*big.Rat`.
6. `core.Object` creation is restricted; there is no generic
   `NewObject(kind, …)` factory. GR cannot mint trusted kernel objects
   with arbitrary Kind/Dimension/Provenance/CorpusStatus.
7. Operations are the 12 supported IDs only; `identify` never via
   `ops.Apply`. `Simplify` cannot create `IDENTIFIED`, change corpus
   status, invoke `Identify`, or consult session state.
8. PASS0 is a mandatory prerequisite gate (see §16). The GR workload
   begins only after PASS0 pins are implemented and verified green with
   `go build ./...`, `go vet ./...`, `go test ./... -count=1`.

Non-negotiable invariants preserved by all GR work: immutability, mint
authority, no generic factory, exact rationals, canonicalization,
ordinal stability, MRC, provenance law + hypothesis contamination,
containment, session/ledger/replay/hash integrity.

---

## 4. Module Boundaries

```text
github.com/PithomLabs/phys      frozen Level-3 substrate
    internal/kernel, core, ops, session,
    mechanics (classical only), relativity (SR only),
    hypothesis (candidate space)

github.com/PithomLabs/phys-gr   Level-1 GR implementation
    depends one-way on phys (public API only)
    never imported by phys
```

Rules:

- All GR work occurs outside the `phys` repository tree (sibling local
  module, e.g. `~/go/phys-gr` with its own `go.mod` requiring
  `github.com/PithomLabs/phys`). No `replace` into kernel internals.
- Level-1 theory packages follow the same one-way dependency
  discipline as Level-2 libraries: they live outside frozen `phys`,
  depend only on the public `phys` API, and never modify the kernel
  tree. A theory package failing due to missing kernel capability is
  Growth Gate evidence, not permission to bypass the boundary.
- Future Level-2 libraries are likewise separate Go modules (e.g.
  `physmath`) depending one-way on frozen `phys`; `phys` never depends
  back. Only a proven Level-3 primitive touches `phys`, via the full
  Growth Gate.

---

## 5. GR Workload

Canonical static, spherically symmetric thread (locked). Index algebra
is developed INSIDE this thread (GR-2), never as a standalone generic
system first. Linearized gravity is reserved for later, not the entry.

```text
GR-0  Representation baseline (coordinates, charts, assumptions)
GR-1  Metric structure gμν
GR-2  Inverse metric + index algebra emergence
GR-3  Christoffel symbols (differentiation boundary encounter)
GR-4  Covariant derivative
GR-5  Riemann curvature
GR-6  Ricci tensor / scalar
GR-7  Einstein tensor + vacuum equation
GR-8  Schwarzschild exterior + weak-field / Newtonian limit
```

Each pass records: objective, files/packages, dependencies,
capabilities, worked examples, tests, expected failures, evidence,
level implicated, exit criteria — and where the substrate bends,
breaks, or silently misrepresents.

GR-8 completion is the default target because it exercises the full
pipeline from representation through physical interpretation and
bridges `mechanics ↔ general_relativity` without kernel promotion.
Intermediate stopping points (e.g. GR-0…GR-3 boundary discovery) are
valid evidence checkpoints; do not force artificial completion to claim
success. Record completed capability + remaining gap + classification.

This workload forces the hard encounters honestly: `r⁻¹`/`r⁻²`
negative powers, `gμν`/`Γ`/`R` index identity, `∂g/∂x`
differentiation demands, `gμν·gνρ=δμρ` contraction, `Rμν=0`
structured relations, Newtonian-limit assumption/regime semantics,
and `sin²θ` trigonometric demand (see §9).

---

## 6. Representation Strategy

1. **GR-domain mathematical structures (`sin²θ`, trigonometric
   functions, tensors, metric, connection, curvature) remain Level-1
   `phys-gr` structures.** They are not `core.Object` kinds unless a
   future Growth Gate independently proves a theory-neutral kernel
   primitive.
2. Scalar rational algebra uses `core.Expr` (constructors +
   read-only inspection). Everything with index/variance/chart/frame
   semantics is a GR-local Go struct holding `core.Expr` components
   plus explicit index metadata.
3. No new kernel expression nodes (`Tensor`, `Derivative`, `Metric`,
   `Field`, `Manifold`) are introduced. The kernel expression model
   stays closed.
4. GR-specific objects (tensors, metrics, indices, connections,
   curvature tensors) are initially Level-1 Go structures, not
   `core.Object` instances — the frozen module has no public generic
   trusted-object factory.
5. Symbol identity stays syntactic at the kernel layer (`Symbol` is a
   string id). Namespacing, scoping, entity refs, coordinate refs are
   GR-local concerns until Level-2 evidence exists.

---

## 7. Kernel Interaction Points

Allowed contact surface only:

| Operation | Use | Constraint |
|---|---|---|
| Constructors (`NewPow`, `NewAdd`, …) | build scalar components | public API only |
| Read-only inspection | traverse `core.Expr` | never `internal/kernel` |
| `Add/Subtract/Multiply/Divide/Pow/Simplify/Substitute` | rational component algebra | pure, session-free; MRC enforced |
| `Differentiate` | GR-3a probe only | expect fail-closed on negative-integer powers |
| `Compare` | build relation artifacts | proves nothing by itself |
| `Solve` | single-pattern `t²` only | `BranchSet ±Sqrt` |
| `SelectBranch` | explicit `gte` constraint | recorded as structured assumption |
| `Session.Identify` | explicit identification | `KindRelation` only, never retyping |
| `Limit` | Newtonian-limit step | variable/value explicit; `Call` traversal fixed-body `lorentz_factor` only |

Forbidden: new `Kind`s, generic factories, kernel expression nodes,
`internal/` imports, assumption/provenance bypass, silent kernel
extension via Level-1/Level-2 code.

Critical differentiation boundary (most important pin):

```text
GR-3a — Kernel boundary probe:
    phys.Differentiate(Pow(r,-1))
      → UnsupportedOperationError
      → classified SPEC-INTENDED-BOUND

GR-3b — Level-1 continuation (only after GR-3a recorded):
    phys-gr/symbolic/differentiate.go operating over public
    core.Expr only. Classified PACKAGE-SOLVABLE. Must not be
    represented as kernel differentiation capability.
```

GR-3b must not modify `phys.Differentiate`, add negative-power
support to the kernel, create a hidden kernel extension, or claim
kernel-level GR differentiation support.

---

## 8. Level-1 GR Package Design

Sibling module `github.com/PithomLabs/phys-gr` (one-way dep on `phys`):

```text
phys-gr/
    go.mod
    coordinates/    chart, symbols, assumptions/regime
    metric/         gμν + inverse, contraction gμν·gνρ=δμρ
    index/          variance, raising/lowering, contraction rules
    connection/     Christoffel assembly (uses symbolic/)
    curvature/      Riemann, Ricci, scalar, Einstein
    vacuum/         Rμν=0 structured relations
    solution/       Schwarzschild exterior
    limit/          weak-field / Newtonian correspondence
    symbolic/
        differentiate.go   GR-local bounded differentiator
        trig.go            sin/cos userland math
    evidence/
        failure_log.md     per-obstacle classified records
        growth_records/    Growth Evidence Records
    adversarial/
        silent_wrongness_test.go
```

Dependency rule: `phys-gr → phys` only. `phys` never imports
`phys-gr`. GR-local math never imports `internal/kernel`.

---

## 9. GR-Local Mathematics Required

Owned by `phys-gr`, never kernel:

1. **Bounded differentiator** (`symbolic/differentiate.go`): handles
   negative-integer powers and rational metric terms over public
   `core.Expr`; preserves kernel invariants; does not bypass
   `phys.Differentiate`.
2. **Trigonometric functions**: standard Schwarzschild angular sector
   `r²dθ² + r²sin²θdφ²` requires `sin(θ)`/`cos(θ)`. GR Level 1 owns
   this representation; the MVP stays closed (`Call =
   lorentz_factor` only). Scalar rational algebra → `core.Expr`;
   trig → GR-local math; tensor/index semantics → GR-local math.
3. **Index/tensor algebra**: variance, contraction, raising/lowering,
   symmetry bookkeeping — emergent in GR-2 inside the workload.
4. **Inverse metric / contraction checks**, curvature assembly,
   assumption/regime bookkeeping for the Newtonian limit.

None of this implies kernel expansion.

---

## 10. Failure Taxonomy

Every obstacle maps to exactly one category, logged with requirement,
minimal example, current vs expected result, category rationale,
Level-1/Level-2 workarounds attempted, and Level-3 implication.
No wishlist entries ("it would be nice if the kernel had tensors" is
not evidence).

```text
SPEC-INTENDED-BOUND
    MVP deliberately does not support it.
    (e.g. Differentiate(Pow(r,-1)) fail-closed.)

PACKAGE-SOLVABLE
    general_relativity handles it locally.
    (e.g. GR-local differentiator continues Christoffel.)

REPRESENTABLE-BUT-UNFAITHFUL
    Substrate encodes it but loses important semantics.
    (e.g. Expression addition erases energy-vs-torque category.)

UNREPRESENTABLE
    Substrate cannot faithfully encode it.
    (Not kernel-growth evidence until L1+L2 attempts fail.)

SILENTLY-WRONG
    Substrate accepts it and produces a result whose semantic
    error is not mechanically exposed. Primary adversarial thread.
```

Only the last three constitute serious evidence requiring deeper
investigation, and even `UNREPRESENTABLE` is not growth evidence
until higher-level attempts have failed.

---

## 11. C1-C8 Concern Ledger

Self-contained; frozen routing. Prevents reopening settled arguments.

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
workload. GR success alone cannot establish kernel semantic safety.
The dangerous outcome is `accepted + wrong`, not `unsupported`.

Probes (each asserts expected vs observed, classifies
`Rejected-correctly / Accepted-correctly / SILENTLY-WRONG`):

1. Symbol ambiguity: same `Symbol("m")` for electron vs proton mass
   in one derivation; check for undetected mixing.
2. Field differentiation: field `f(t)` represented as bare
   `Symbol("f")`; `Differentiate(f-expr, t)` yields `0` (accepted,
   `DERIVED`) where nonzero expected → `SILENTLY-WRONG` candidate;
   investigate L1/L2 scoping fix.
3. Invalid contraction / raising: wrong variance or index pairing in
   GR-local structures must not pass silently.
4. Dimension-valid semantic errors: energy-vs-torque `Expression +
   Expression` mixing; invalid substitution with same dimension but
   wrong category must stay rejected or be flagged unfaithful.
5. Type-decay addition: `E_expr + τ_expr` both `M L²T⁻²`.
6. Inconsistent-assumption coexistence: contradictory regime
   assumptions coexisting without structural conflict.

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
separate Go module; promotion never touches the 42-file tree.
Level 2 cannot become a hidden kernel extension: if it needs the
kernel to do X to preserve a generic trusted invariant, that is
Growth Gate evidence, not permission to sneak X into a Level-2 API
(no-bypass rule).

---

## 14. Level-2 → Level-3 Growth Gate

Level-3 (kernel primitive) requires the full Growth Evidence Record
(§15) plus proof of all of: generic, theory-neutral, non-semantic,
foundational, reusable across frameworks, invariant-bearing, and
unlivable above the kernel (provisional Level-2 attempt built without
changing kernel semantics, documented failure to preserve the
invariant). Burden of proof is on growth. Default outcome: no kernel
modification presumed; all gate outcomes remain admissible and are
determined from recorded evidence.

Wording lock: "No kernel modification is presumed. All Growth Gate
outcomes remain admissible and must be determined from recorded
evidence." (replaces bare "Default outcome: NO KERNEL GROWTH.")

---

## 15. Evidence Record Format

Every `UNREPRESENTABLE` / `REPRESENTABLE-BUT-UNFAITHFUL` /
`SILENTLY-WRONG` obstacle that survives L1 and provisional-L2
attempts gets one record with these enumerated fields (no undefined
"20-field record" references):

1. capability; 2. minimal counterexample; 3. failure category (one of
   §10); 4. L1 attempt; 5. L2 attempt; 6. L1 failure reason; 7. L2
   failure reason; 8. genericity; 9. foundationality; 10. second
   consumer (named); 11. second-consumer worked example; 12. non-goal
   collision check; 13. silent-wrongness analysis; 14. artifact
   compatibility (canonical form / hash); 15. MRC impact; 16. replay
   impact; 17. migration impact; 18. minimal kernel change (exact
   diff-shape, no implementation); 19. independent review reference;
   20. human approval (signature/decision).

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
- Gate: `go build ./...`, `go vet ./...`,
  `go test ./... -count=1` green; baseline re-verified and frozen.

**PASS1 / GR workload tests:** per-pass worked examples with golden
states (notably the E=mc² derivation chain stays intact and is never
used as a GR premise), kernel-vs-local differentiation pair
(GR-3a/GR-3b), trig-ownership test, contraction tests, vacuum-equation
and Schwarzschild golden components, Newtonian-limit correspondence.

---

## 17. Mutation Testing

Negative/mutation probes to confirm fail-closed behavior is real:

- Flip Pin A fixture back to `Identify(Energy, E/c²)` → must fail
  dimension check (proves corrected fixture is non-vacuous).
- Remove Pin B negative-exponent guard → GR-3a probe must fail
  loudly (proves pin is load-bearing).
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

Four admissible terminal states (evidence-determined):

```text
NO-GROWTH
    Kernel boundary correct; all GR needs met in userland.
    Acceptable success.

LEVEL-2-GROWTH
    Demonstrated cross-theory abstraction (named 2nd consumer +
    worked example); new shared module above kernel.

KERNEL-GROWTH CANDIDATE
    Full evidence record + proven L2 failure to preserve a
    generic invariant; minimal kernel-change proposal for
    human decision.

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
  kinds.
- No semantic subtyping, entity identity, SMT solving, general
  entailment, general calculus, or tensor algebra in
  `internal/kernel`.
- No `semantic/`, `symbols/`, `constraints/`, `tensor/`, `manifold/`,
  `calculus/`, `hilbert/` pre-building before GR generates evidence.
- No transcendentals, symbolic exponents, or rational-function
  differentiation added to the MVP to make GR easier.
- No weakening of `Substitute`, `Identify` authority, MRC, provenance
  law, or closed-world tree.
- No stored corpus result (notably `MassEnergyRelation`) used as a
  derivation premise (firewall violation).
- No empirical-validity claims from formal consistency; no corpus
  status assigned by computation.

---

## 20. Definition of Done

1. PASS0 pins implemented, green, baseline frozen (§16).
2. `phys-gr` sibling module executes GR-0…GR-8 (or principled
   checkpoint with gap classification) without modifying `phys`.
3. GR-3a kernel failure + GR-3b local continuation both recorded with
   correct classifications.
4. Trig ownership demonstrated (GR-local `sin`; kernel `Call`
   unchanged).
5. C1–C8 ledger entries evidenced; failure log classifies every
   obstacle exactly once (§10).
6. Silent-wrongness suite executed; every `SILENTLY-WRONG` has an
   evidence record or a Level-1/Level-2 containment argument.
7. Promotion decisions follow §§13–14 with enumerated records; Level 2
   never bypasses the kernel, Level 3 never presumed.
8. Final report states: what `phys` can do, what it deliberately
   refuses, what userland provides — i.e. **boundary learned.
   NO-GROWTH is an acceptable successful outcome.**

Final decision principle: do not make the kernel smarter because GR
is complicated. Make userland richer first. Only promote something
into the kernel when multiple theories independently prove that the
primitive itself is foundational, theory-neutral, and impossible to
safely implement above the kernel boundary.

---

*Document-only artifact. Modifies no specs, kernel, manifests, corpus,
or frozen tree. Implementation agents execute it without reopening
architectural decisions; any spec ambiguity escalates to human review.*

# Plan 11.1 Revised — GR Implementation Plan v3 (Final Pre-Handoff Execution Contract)

Version: v3 — 20261003.
Status: normative execution contract. Planning document only.
No code, corpus, spec, tree, test, manifest, or module changes made by this document.

INPUT: `plan10/GR_IMPLEMENTATION_PLAN_20261002.md` (v2).
OUTPUT: this file `plan10/GR_IMPLEMENTATION_PLAN_20261003.md` (v3).
Append-only versioning: v2 artifacts are untouched. No existing plan file was revised.

Changelog / audit trail:
- v1 (20261001): initial 20-section execution contract.
- v2 (20261002): six-review reconciliation. Locked GR ↔ `core.Object`
  HYPOTHESIS bridge concept, single bounded GR-local algebra, real
  `Sin`/`Cos`/`UnknownFunction`, GR-3a/GR-3b independence,
  GR-local replay, derivation-first Schwarzschild, GR-local Newtonian
  reduction, units/conventions pins, `LEVEL-2-CANDIDATE-PENDING`.
- v3 (20261003): final pre-handoff pass. Keeps D1–D18 unchanged, adds
  D19–D40, and integrates every pin into the normative sections:
  exact `GRExpr` algebra, `UnknownFunction` representation (max total
  order 2), complete `Diff` table, exact bridge contract
  (`ToCoreExpr`/`FromCoreExpr`/`ToCoreObject`, lazy, content-derived
  ID, dimension/assumption policy), `mu` vs `M` separation, structured
  assumption translation, chart pin, dense row-major tensor model,
  MTW curvature formulas with golden signs, machine-certified
  Schwarzschild script (derive + independent verify), three-valued
  `ZeroTest` with oracle, GR-local naming, full replay protocol,
  GR-3a direct-`ops` path, Newtonian `Truncate` procedure, GR-8
  session recording segment, pending-state outcomes, provisional-L2
  escalation path, independent-review workflow, executable Fields
  18–20, canonical JSON/SHA-256, module file plan, PASS0 prerequisite,
  `go.work` mechanics, GR-0…GR-8 execution tables, bridgeability
  census, expanded silent-wrongness suite, stop conditions, wording
  corrections, spec-precedence rule. Zero architectural choice remains
  for the implementation agent.

Lineage: `plan11.md` → `plan11_review.md` → `plan11.1.md` →
`plan11.1_review.md` → v1 → v2 → this v3.
Frozen authority: `plan10/specs_v2_3.md` is authoritative over every
historical planning document. Where any older planning text conflicts
with `specs_v2_3.md`, the frozen specification governs and the
historical text is non-normative. This document never modifies the
frozen spec.

```
phys
 └── frozen kernel (Level 3)
     ├── no GR metric semantics
     ├── no GR tensor semantics
     ├── no transcendental math
     └── no theory semantics
     (Existing nominal SR-related Kinds such as Spacetime,
      MinkowskiMetric, FourMomentum remain part of the frozen MVP.)

phys-gr (Level 1, sibling module github.com/PithomLabs/phys-gr)
 ├── GRExpr algebra: Diff / Normalize / ZeroTest / Subst / Truncate / Reduce
 └── GR structures: Chart / Tensor / Metric / Connection / Curvature
 └── exact kernel subset → temporary HYPOTHESIS Object → frozen phys ops
 └── reintegrate GRExpr

Growth Gate: decides promotion from recorded evidence only.
```

Decision register (normative; D1–D18 retained, D19–D40 added in v3):
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
D19 `GRExpr` node set closed to 9 nodes, immutable, canonical order,
integer exponents (negative permitted locally), max derivative order 2.
D20 `UnknownFunction` struct `{Name, Args, DerivativeOrder}` with
`len` equality, arity pin, order-sum bound, atom semantics. D21 full
`Diff` table pinned including product/chain rules and dependency
zeros. D22 bridge triple
(`ToCoreExpr`/`FromCoreExpr`/`ToCoreObject`) with lazy temporary
objects, named bridge error, content-derived ID over pinned payload.
D23 dimension policy (geometrized `Dimensionless` bridge; separate SI
audit path) with `mu` vs `M` separation (`mu = G*M/c^2`, `k2 = -2*mu`).
D24 structured assumption translation, caller-owned, no hidden
injection, regime-vs-precondition split. D25 chart pin
(`spherical-static`, `[t,r,theta,phi]`, ASCII names). D26 dense
row-major tensor model (`Rank 0..4`, `4^Rank` components,
lexicographic flattening, `SymmetryRule{Permutation,Sign}`),
eight required ops, validation + symmetry inheritance, Connection is
not a Tensor. D27 MTW curvature formulas + golden signs, `-+++`,
`Lambda = 0` vacuum specialization. D28 14-step machine-certified
Schwarzschild script (candidate-and-certify, no general ODE solver).
D29 PRIMARY derivation vs SECONDARY verification artifacts with
independent success criteria. D30 `Normalize` rule set + three-valued
`ZeroTest` + test-only oracle + negative controls. D31 GR-local naming
(`Diff/Normalize/ZeroTest/Subst/Truncate/Reduce`), kernel names never
shadowed. D32 replay protocol (11 fields, 64-zero genesis,
`gr-step-NNNNNN`, hash construction, round-trip law). D33 GR-3a uses
`ops.Differentiate` directly, never `Session.Step`. D34 Newtonian
`Truncate(expr,eps,1)` retaining `O(eps^0..1)`, `Phi = -GM/r`
(`a_r = -GM/r^2` where included), no `APPROXIMATED`. D35 GR-8 session
use limited to final correspondence segment ending at `Conclude`.
D36 terminal states
(`NO-GROWTH/LEVEL-2-CANDIDATE-PENDING/KERNEL-GROWTH-CANDIDATE-PENDING/ESCALATE-TO-SPEC`),
no automatic promotion. D37 provisional-L2 escalation path with
pre-registered falsifiers. D38 second-agent review workflow for
non-`PACKAGE-SOLVABLE`. D39 executable Fields 18–20. D40 canonical
JSON/SHA-256 for all 8 GR artifact kinds with timestamp exclusion.

Explicitly rejected (permanent): kernel tensor types; kernel
`Sin`/`Cos`; kernel `UnknownFunction`; extending `phys.Differentiate`;
expanding `phys.Limit`; general ODE solver; general CAS; placeholder
`sin_theta`; arbitrary rational-exponent local calculus; Substrate
Efficacy Ratio or any reuse threshold; automatic Level-2 promotion;
bare-symbol differentiation treated as kernel defect; GR inside `phys`;
`ops` expansion for GR-local ops; symbolic integration; asymptotic
engine; arbitrary function registry; general theorem proving.

> No architectural decision remains for the implementation agent.
> Remaining discoveries belong to implementation evidence, failure
> classification, or human Growth Gate adjudication.

---

## 1. Executive Summary

GR is a hostile workload against frozen `phys`, not a kernel feature:

> Can current kernel primitives support a serious theory package without being modified?

Record separately: what `phys` can do, what it deliberately refuses,
what userland provides. `NO-GROWTH` is acceptable success.

Gating:

```text
PASS0 (separate conformance task, prerequisite gate)
  → PASS1 (this plan acknowledged)
  → GR-0…GR-8 + adversarial suite + bridgeability census
  → Growth Gate
```

v3 = same architecture as v2, every implementation choice pinned.

---

## 2. Architectural Principles

1. Go-inspired layering: small core + rich ecosystem.
2. Kernel answers primitive capabilities useful across theories
   (exact rationals, immutable expressions, dimensional algebra,
   provenance/assumption tracking, canonical replay, deterministic
   transforms); never "which mathematical objects belong to physics".
3. Kernel math-agnostic, theory-agnostic, minimal, composable,
   invariant-preserving.
4. Algebra does not manufacture physical meaning.
5. Userland first; promotion only on multi-theory foundational proof.
6. Level direction: Level 1 (theory-local) is consumed by Level 2
   (shared, human-curated), which may inform Level 3 (trusted kernel)
   solely through the Growth Gate. Lower levels are never modified by
   upper levels outside the gate.
7. Invariants, not operation counts, drive growth (D16/D37).

---

## 3. Frozen Baseline Assumptions

1. `specs_v2_3.md` byte-for-byte frozen and authoritative.
2. Both theory manifests frozen.
3. Operational freeze: 39 implementation files + 3 documentation files
   (`AGENTS.md`, `mechanics/README.md`, `relativity/README.md`) = 42.
   PASS0 verifies membership before editing; 42-count unchanged.
4. `internal/kernel` unreachable from `phys-gr` (Go compiler boundary).
   Contact via public `core`/`ops`/`session`/`hypothesis` facade only.
5. `core.Expr` closed MVP node set; sole `Call` id `lorentz_factor`;
   `Pow` exponent `*big.Rat`.
6. Trusted objects originate solely via fixed domain constructors,
   hypothesis construction, pure operations, session artifacts. No
   generic factory. GR bridge uses the frozen hypothesis constructor
   (arbitrary `Expr` + explicit `Kind`/`Dimension` → forced
   `HYPOTHESIS`/`NONE`).
7. 12 operation IDs only; `identify` never via `ops.Apply`.
8. PASS0 prerequisite gate (§16); post-PASS0 baseline record required
   (commit SHA, clean-tree status, SHA-256 of spec + manifests +
   modified files, 42-inventory, build/vet/test). GR evidence cites it.

---

## 4. Module Boundaries

`github.com/PithomLabs/phys` (frozen) vs
`github.com/PithomLabs/phys-gr` (Level 1, one-way dep, public API).
`mechanics/` + `relativity/` grandfathered (D3). New Level-1 work
always external. Local `go.work` at the shared parent connecting
pristine siblings; CI pins exact `phys` commit. Local `replace`
permitted solely to the pristine hash-verified checkout, never to a
modified fork. No vendoring. No `internal/kernel` import. Workspace
location example: parent directory containing `phys/` and `phys-gr/`
side by side; exact path recorded in the GR run record.

---

## 5. GR Workload (GR-0…GR-8 locked thread)

Static spherical thread; linearized gravity deferred.

```text
GR-0 representation + conventions + units + API audit
GR-1 general static spherical metric with A(r), B(r)
GR-2 inverse metric + concrete tensor/index engine
GR-3 Christoffel (3a kernel probe + 3b GR-local Diff)
GR-4 covariant derivative
GR-5 Riemann
GR-6 Ricci + scalar
GR-7 Einstein + vacuum equations
GR-8 Schwarzschild derivation + verification + flat check + weak-field
```

Full execution tables are in §27. Default target is the complete
thread; checkpoints with gap classification are valid evidence.

---

## 6. Representation Strategy (GRExpr algebra — D19/D20)

Node set closed to exactly 9 nodes. Additions are STOP conditions:

```text
Symbol{name}  Rational{p,q reduced, q>0}  Add[terms]  Mul[factors]
Neg[x]  Pow{base, n int}  Sin[x]  Cos[x]  UnknownFunction{...}
```

1. Immutable representation. All constructors return new values.
2. Exact fields: `Symbol{name ASCII}`; `Rational{num int, den int>0,
   gcd=1}`; `Add{terms []GRExpr ≥2}`; `Mul{factors []GRExpr ≥2}`;
   `Neg{x GRExpr}`; `Pow{base GRExpr, exp int}`; `Sin{x}`;
   `Cos{x}`; `UnknownFunction{Name, Args []Symbol,
   DerivativeOrder []uint8}`.
3. Canonical child ordering: `Add`/`Mul` terms sorted by
   `CanonicalBytes`; `Neg` single child; `Pow` base-then-exp.
4. Equality: structural equality over canonical form.
5. Serialization: canonical JSON with fixed field order, exact
   rational `num/den` strings, deterministic arrays; no maps in the
   authoritative path.
6. Exponent domain: integers only; negative integers permitted in
   `phys-gr`. No symbolic exponents, no rational-exponent calculus.
   `Pow` with negative exponent requires nonzero-base evidence at
   use sites per operation contract.
7. Unsupported named errors: `UnsupportedExponentError`,
   `UnsupportedNodeError`, `DerivativeOrderExceededError`,
   `UnrepresentableKernelProjectionError`, `NonZeroEvidenceError`.
8. Maximum supported derivative order: total order 2 (§7).

UnknownFunction pin (D20): `Name` (ASCII, no punctuation encoding
derivative state); `Args` ordered symbols (arity fixed per function;
`A`,`B` arity 1 with `Args=[r]`); `DerivativeOrder` vector with
`len == len(Args)`, entries ≥0, `sum ≤ 2`. `A(r)=[0]`,
`A'(r)=[1]`, `A''(r)=[2]`. Distinct orders are distinct algebraic
atoms for normalization yet structurally related (same `Name`+`Args`).
Differentiating wrt a non-argument returns zero (§7). No generic
`Derivative` node.

---

## 7. Kernel Interaction Points + Diff table (D21/D31/D33)

Contact surface: constructors, read-only inspection, arithmetic +
`Simplify`/`Subst`-equivalent kernel ops on HYPOTHESIS wrappers,
`Differentiate` (GR-3a probe only), `Compare`, single-pattern `Solve`,
constrained `SelectBranch`, `Session.Identify` (`KindRelation` only).
`Limit` never performs GR-8. No new Kinds/nodes/factories/`ops` IDs.

GR-local names (D31): `Diff`, `Normalize`, `ZeroTest`, `Subst`,
`Truncate`, `Reduce`. Kernel names `Differentiate`, `Simplify`,
`Limit`, `Solve`, `Substitute` are never used for GR-local functions.

Complete `Diff(e, x)` contract (independent Level-1 operation,
distinct from frozen kernel `Differentiate`):

```text
Diff(Rational, x) = 0
Diff(Symbol(x), x) = 1
Diff(Symbol(y), x) = 0 for y != x
Diff(Add(u..), x) = Add(Diff(u,x)..)
Diff(Mul(u..), x) = n-ary product rule sum
Diff(Neg(u), x) = Neg(Diff(u,x))
Diff(Pow(u,n), x) = Mul(n, Pow(u,n-1), Diff(u,x)), n integer
Diff(Sin(u), x) = Mul(Cos(u), Diff(u,x))
Diff(Cos(u), x) = Mul(Neg(Sin(u)), Diff(u,x))
Diff(UnknownFunction(f,args,orders), x) =
    orders with x-slot incremented if x ∈ args and sum+1 ≤ 2,
    else 0 when x ∉ args,
    else DerivativeOrderExceededError
```

Chain rule, product rule, integer-only powers, dependency zeros,
max-order enforcement, named errors. GR-3a path (D33): bridge →
`ops.Differentiate` directly (pure public operation), never
`Session.Step` (failed-step session semantics are out of scope for a
boundary probe). GR-3b logged solely in the GR replay domain.

Reuse rule: kernel ops whenever the subexpression is in the exact
kernel-mirror subset; GR-local rules solely for absent capabilities.

---

## 8. Level-1 GR Package Design (file plan — D26/D32)

```text
phys-gr/
  coordinates/chart.go        Chart (D25) + ConventionSet
  tensor/tensor.go            Tensor/IndexSlot/SymmetryRule (D26)
  tensor/ops.go              TensorAdd/Subtract/ScalarMul/Contract/
                             RaiseIndex/LowerIndex/ApplyMetric/ApplyInverseMetric
  metric/metric.go            Metric + inverse + contraction check
  connection/connection.go    Connection coefficients (NOT Tensor)
  curvature/curvature.go      Riemann/Ricci/scalar/Einstein (D27)
  vacuum/vacuum.go            Rμν=0 relations
  solution/schwarzschild.go   14-step script + verification (D28/D29)
  limit/weakfield.go          Truncate + Phi/a_r extraction (D34)
  symbolic/expr.go            GRExpr 9-node algebra + canonicalization (D19/D20)
  symbolic/differentiate.go   Diff table (D21)
  symbolic/normalize.go       Normalize + ZeroTest + oracle interface (D30)
  symbolic/bridge.go          ToCoreExpr/FromCoreExpr/ToCoreObject (D22–D24)
  symbolic/weakfield.go       MAY fold into limit/weakfield.go; plan records
                             canonical home as limit/weakfield.go with
                             symbolic/weakfield.go as thin forwarder if retained
  replay/replay.go            GRStep chain + hashing (D32)
  evidence/                   failure_log.md + growth_records/ + kernel_contact.md
  adversarial/                SW-1…SW-13 suite (§12)
```

Merging permitted solely as stated above. No shared Level-2 packages.

Tensor model (D26): `Tensor{Rank 0..4, Slots []IndexSlot, Components
[]GRExpr dense row-major length 4^Rank lexicographic, Symmetries
[]SymmetryRule{Permutation, Sign ±1}, ChartID}`; `IndexSlot{slot,
variance covariant|contravariant, domain=spacetime}`; component keys
are positional tuples via `offset(i)=(((i0·D+i1)·D+i2)…)`; absent
components are explicit zeros (dense storage makes omission
impossible). Validation: chart compatibility, variance, pairing,
domain, declared symmetries (verified before reliance); symmetry
inheritance pinned per op (Add/Sub preserve intersection; ScalarMul
preserves; Contract inherits non-contracted symmetries with sign
adjustment; Raise/Lower preserve). Connection is not a Tensor:
generic tensor ops reject connection coefficients by type.

Chart (D25): `Chart{ID, Coordinates []Symbol, Order []int}`;
canonical `ID="spherical-static"`, coords `[t,r,theta,phi]`,
`Order=[0,1,2,3]`; ASCII names
`t,r,theta,phi,A,B,M,mu,G,c,eps,k1,k2`; `x^0` documentation only.
Rational chart `x=cos(theta)` is a secondary bridgeability census
(§27), never the primary representation.

---

## 9. GR-Local Mathematics Required

`Diff` (§7); `Normalize` + three-valued `ZeroTest` (§13);
`Sin`/`Cos` + workload identities; `UnknownFunction` calculus;
`Truncate`; tensor engine; inverse/contraction checks. Bounded to this
workload (no general CAS/ODE/integration/asymptotics/registry).

---

## 10. Failure Taxonomy

`SPEC-INTENDED-BOUND` / `PACKAGE-SOLVABLE` (agent-classifiable) /
`REPRESENTABLE-BUT-UNFAITHFUL` / `UNREPRESENTABLE` / `SILENTLY-WRONG`
(second-agent review mandatory: agent classification + independent
classification + disagreement record + human resolution).

---

## 11. C1–C8 Concern Ledger

Per v2 (unchanged routing): C1–C4 Level-2 investigations; C5 future;
C6 closed; C7–C8 PASS0 pins. Full table retained from v2 §11.

---

## 12. Silent-Wrongness Testing Strategy (SW-1…SW-13)

Legitimacy definition: a representation is legitimate when it uses the
pinned Level-1 API exactly as specified (correct node types, chart,
variance, conventions). Misuse of the API is user error, never
`SILENTLY-WRONG`. Corrected semantics: bare `Symbol("f")` wrt `r` →
0 is `ACCEPTED-CORRECT`.

Suite: SW-1 symbol/entity ambiguity; SW-2 `UnknownFunction`
dependency; SW-3 chart mixing; SW-4 Γ-as-tensor (must reject);
SW-5 invalid contraction; SW-6 invalid symmetry; SW-7
dimension-valid semantic mismatch; SW-8 type-decay addition; SW-9
assumption inconsistency; SW-10 convention/signature/Riemann-sign
mismatch; SW-11 `ZeroTest` false-zero (negative controls); SW-12
bridge round-trip/canonical mismatch; SW-13 tensor
canonicalization/replay inconsistency.

---

## 13. Level-1 → Level-2 Promotion Rules

Named second established consumer + worked example + shared semantics
+ non-leakage + no-bypass. Insufficient: hypothetical reuse. Home:
separate module. Never automatic.

Normalize pin (D30): common-denominator, rational-factor/like-term
collection, negative-power normalization, exact rational
normalization, minimum trig set (`Cos² → 1−Sin²`,
`Sin²+Cos² → 1`). `ZeroTest ∈ {ZERO, NONZERO, UNDECIDED}`; `ZERO`
solely on exact normalization proof; `NONZERO` solely on direct
structural proof; otherwise `UNDECIDED` (never guess).
`UnknownFunction` terms compare by `(Name,Args,Order)` atoms.
Negative controls prove known non-zeros are never `ZERO`. Test oracle
is test-only: it falsifies false `ZERO`, never promotes samples to
proof.

---

## 14. Level-2 → Level-3 Growth Gate

Full record + generic/theory-neutral/non-semantic/foundational/
reusable/invariant-bearing/unlivable-above proof. Wording lock
retained. Pre-registered falsifiers for genuine evidence: trusted
invariant, canonical integrity, deterministic replay, MRC
enforcement, or immutable construction cannot be preserved above the
kernel. "Userland required more code" is never evidence.

---

## 15. Evidence Record Format (Fields 18–20 executable)

Fields 1–17 per v2. Field 18 `MinimalKernelChange`: package/file,
exported identifiers, change kind (add API / alter invariant / add
operation / add expression node / alter canonical encoding / alter
replay semantics / alter MRC), invariant/MRC/canonical/replay impact,
no source code. Field 19 `IndependentReview`: reviewer identity,
record ID, first/second classification, disagreement resolution.
Field 20 `HumanApproval`: record ID, reviewer identity, UTC timestamp,
`approved|rejected|provisional`, note. No signing required. No human
approval ⇒ `PROVISIONAL` ⇒ no kernel modification.

---

## 16. Test Strategy (PASS0 prerequisite)

PASS0 is a separate task; the GR agent never performs PASS0 changes.
It verifies: 42-membership, post-PASS0 baseline record, clean
checkout, recorded commit + SHA-256 values, build/vet/test green.
Absent/mismatched ⇒ `STOP`, `ESCALATE-TO-SPEC`/human review, no GR.

Pins: Pin A `Identify(RestMass, E/c²)` → `KindRelation`/`M`/
`IDENTIFIED`/`NONE`, operands unchanged. Pin B
`Differentiate(Pow(x,-1))` → `UnsupportedOperationError`, target
`ops/differentiate_test.go` (membership-checked; halt rather than
create). Pin C `AGENTS.md` caveat. Test-sanity (not production
mutation): restore broken `Identify(Energy, E/c²)` fixture.

GR tests: golden states per pass, GR-3a/3b pair, trig ownership,
contraction, vacuum + Schwarzschild goldens, flat + Newtonian checks,
Kernel Contact Ledger rows
(`pass/kernel_operation/inputs/outcome/local_operation/boundary_note/baseline_id`
+ observational `node count/term count/wall-clock`), descriptive only.

---

## 17. Mutation Testing

Disposable clone/worktree only; destroy mutant; re-verify
authoritative checkout. Production mutants: `Identify` returns
operand Kind; remove negative-power guard; weaken `Substitute`
equality; remove `SelectBranch` constraint; remove assumption
conflict; `ZeroTest` always `ZERO`; `Contract` skips variance check.
AST bounds are observational instrumentation.

---

## 18. Expected Outcomes

`NO-GROWTH` / `LEVEL-2-CANDIDATE-PENDING` (reusable appearance
recorded; promotion blocked pending named second consumer + human
curation; no library created here) /
`KERNEL-GROWTH-CANDIDATE-PENDING` (record evidence → human-authorized
provisional L2 attempt → gate advancement solely on L2 failure) /
`ESCALATE-TO-SPEC`. No ratio, no threshold, no automatic promotion.

---

## 19. Explicit Non-Goals

Per rejection list above (§15 v2 list retained + general
CAS/ODE/integration/asymptotics/registry/prover). SR corpus: GR never
modifies `relativity/`; may cite SR as external reference; SR stays
Special Relativity only.

---

## 20. Definition of Done

1. PASS0 record verified; baseline matches checkout (§16) else STOP.
2. `phys-gr` executes GR-0…GR-8 tables (§27) modifying only `phys-gr`.
3. GR-3a direct-`ops` failure + GR-3b independent `Diff` recorded.
4. Real `Sin`/`Cos`/`UnknownFunction` demonstrated; kernel `Call`
   unchanged; no placeholders.
5. Every obstacle classified once (§10); ledger + contact rows complete.
6. SW-1…SW-13 executed with corrected semantics; `SILENTLY-WRONG`
   records or containment arguments present.
7. Schwarzschild derived (14 steps) + independently verified; flat
   `M→0`; weak-field `Phi` (plus `a_r` where included) via `Truncate`.
8. Tensor/chart/convention/unit pins honored; artifact hashes recorded.
9. Fields 18–20 complete; independent review for non-`PACKAGE-SOLVABLE`;
   human approval for any candidate.
10. Bridgeability census recorded (descriptive, no score).
11. Final report: what `phys` does / deliberately refuses / userland
    provides. `NO-GROWTH` acceptable.
12. Machine checklist values verified (§36 footer).

---

## 21. Bridge Contract (D22–D24)

```text
ToCoreExpr(GRExpr) (core.Expr, error)
FromCoreExpr(core.Expr) (GRExpr, error)
ToCoreObject(GRExpr, Dimension, AssumptionSet, ConventionSet)
    (core.Object, string bridgeID, error)
```

Mirror subset: `Symbol/Rational/Add/Mul/Neg/Pow(int)` solely.
`Sin`/`Cos`/`UnknownFunction` ⇒
`UnrepresentableKernelProjectionError`; no placeholder encoding.
Lazy: temporary `core.Object` created solely at deliberate kernel-op
call sites, discarded after reintegration into `GRExpr`; never a
persistent GR-domain object. Mixed expressions project representable
subexpressions/terms individually; GR-local whole remains
authoritative. `FromCoreExpr` accepts the same subset; `Sqrt`/`Call`/
`Relation`/`BranchSet` never import. Round-trip law:
`FromCoreExpr(ToCoreExpr(t)) == t` for mirror-subset terms.

Bridge ID: `gr/bridge/<hex SHA-256 of canonical bridge payload>`;
payload = canonical `core.Expr` + `Kind` + `Dimension` +
`Assumptions` + `Conventions`. Deterministic; no UUID/timestamp/
random/counter. Assumptions: caller-supplied structured
`core.Assumption` values solely (`r > 0`, `r >= 0`, `A(r) != 0` where
representable); no raw strings; no hidden injection. Regime
assumptions (physical, GR-local) are distinct from mathematical
preconditions (kernel-op requirements); record both separately.

---

## 22. Bridge Dimension / Units (D23)

Geometrized tensor track (GR-0…GR-8c): `G = c = 1`, bridge
`core.Dimensionless()`. Dimensional-meaning loss recorded as GR
limitation evidence, never kernel defect. SI audit (GR-8d): explicit
SI dimensions on a separate bridge path. Tensor-component
dimensionality stays GR-local metadata (basis-dependent). Symbols:
geometric length `mu`; SI mass `M`; `mu = G*M/c^2`; `k2 = -2*mu`;
`A = 1 − 2*mu/r` geometric; `mu/r` dimensionless pre-bridge. SI audit
verifies `G*M/c² → L`, `GM/(c²r) → 1`, `GM/r² → L/T²`. `M` never
silently changes meaning.

---

## 23. GR-8 Procedures (D27–D29/D34/D35)

Curvature (MTW, `-+++`, `[t,r,theta,phi]`, `Lambda = 0`):

```text
Gamma^rho_mu_nu = ½ g^rhoσ (∂μ gνσ + ∂ν gμσ − ∂σ gμν)
R^rho_σμν = ∂μ Γ^rho_νσ − ∂ν Γ^rho_μσ
            + Γ^rho_μλ Γ^λ_νσ − Γ^rho_νλ Γ^λ_μσ
Rμν = R^rho_μρν ; R = g^μν Rμν ; Gμν = Rμν − ½R gμν
Vacuum: Gμν = 0 ⇔ Rμν = 0
```

Golden signs: `R_thth = 1 − A − rA'` (ansatz); Schwarzschild
`A = 1 − 2mu/r` gives vanishing Ricci; `g_tt ≈ −(1+2Φ/c²)`.

14-step script: (1) ansatz metric; (2) inverse; (3) Γ; (4) required
Ricci (`R_tt,R_rr,R_thth`); (5) certify
`R_tt/A + R_rr/B = (AB)'/(rAB²)`; (6) vacuum ⇒ `(AB)' = 0`;
(7) `AB = k1` via candidate-and-certify; (8) asymptotic flatness
(`A,B → 1` as `r → ∞`) ⇒ `k1 = 1`; (9) `B = 1/A`; (10) certify
`R_thth = 1 − A − rA'`; (11) vacuum ⇒ `(rA)' = 1`; (12)
candidate-and-certify `rA = r + k2`; (13) `A = 1+k2/r`, `B = 1/A`;
(14) verify all vacuum components. Machine certifies supplied
candidates; it does not discover Schwarzschild autonomously. No
general ODE solver. `Session.Identify` use (if any): solely for the
integration-constant correspondence (`k2 = −2GM/c²` SI /
`k2 = −2mu` geometric) in the GR-8 recording segment, producing a
`KindRelation`/`IDENTIFIED`/`NONE` relation artifact; exact call site
recorded in the GR trace.

Derivation (PRIMARY) vs verification (SECONDARY): PRIMARY runs the
14 steps from `A(r),B(r)`; success = certified chain + `Rμν = 0`
via `ZeroTest(ZERO)`; failure = any step `UNDECIDED`/nonzero with
classified evidence. SECONDARY starts from closed-form
`A = 1−2mu/r, B = 1/A` constructed independently (never reusing
PRIMARY intermediates) and verifies all vacuum components;
success/failure likewise. SECONDARY never rescues PRIMARY.

Newtonian: `eps = G*M/(c²r)`; `Truncate(expr, eps, 1)` retains
`O(eps^0),O(eps^1)`, discards higher; `g_tt = −(1+2Φ/c²)+O(eps²)`;
`Phi = −GM/r`; `a_r = −GM/r²` where included. No `APPROXIMATED`.

GR-8 session segment: sole `phys.Session` use; objects compared are
HYPOTHESIS kernel-representable scalars (weak-field `g_tt`
expansion vs `−(1+2Φ/c²)` form); relation via `Compare`; constant
correspondence via `Identify` as above; expected
`IDENTIFIED`/`NONE`; GR trace holds all mathematics, session ledger
holds solely the correspondence relations.

---

## 24. GR-0…GR-8 Execution Tables (D25–D29/D33/D34)

Each pass carries: Objective / Inputs / Outputs / Data structures /
Kernel contact / Level-1 operations / Worked example / Golden result /
Negative tests / Expected failures / Evidence artifacts /
Classification / Exit criteria / Stop condition / Human review point.

- GR-0: audit chart/conventions/units/public API. Golden: canonical
  chart + ASCII names + `mu`/`M` split recorded; `ToCoreExpr` domain
  table acknowledged. Stop: any missing pin.
- GR-1: ansatz metric `diag(−A,B,r²,r²Sin²θ)` as `Tensor` rank 2 +
  `UnknownFunction` A/B. Kernel: none (unrepresentable whole;
  projectable rational coefficients where useful). Golden: component
  list + chart binding.
- GR-2: inverse + engine ops (`Contract/Raise/Lower/ApplyMetric/
  ApplyInverseMetric`) with validation. Golden: `g·g⁻¹ = δ`.
  Negatives: chart mismatch, variance error, bad pairing rejected.
- GR-3: 3a direct `ops.Differentiate(Pow(r,−1))` ⇒
  `UnsupportedOperationError` (`SPEC-INTENDED-BOUND`); 3b
  `Diff`-based Γ golden set. Census: spherical vs `x=cosθ`
  representability (descriptive).
- GR-4: covariant derivative via `Diff` + Γ + `Normalize`. Golden:
  metric-compatibility spot check.
- GR-5: Riemann via D27 formula (needs `A''/B''`, order 2). Golden:
  symmetries + selected components.
- GR-6: Ricci contraction + scalar. Golden: `R_tt,R_rr,R_thth`
  forms; script steps 4–5 certified.
- GR-7: Einstein + vacuum relations `Gμν = 0` / `Rμν = 0`. Golden:
  structured relations (not solved here).
- GR-8: 14-step derivation + independent verification + `M→0` flat
  check + `Truncate` correspondence + session recording segment.
  Golden: `A = 1−2mu/r`, `Φ = −GM/r`. Stop: any general-solver need.

`T_pass` wall-clock budgets are human-set execution observations;
overrun never constitutes growth evidence.

---

## 25. Stop Conditions + Spec Precedence + Wording Locks

STOP (no improvisation): baseline invalid; frozen API contradicts
plan; required Level-1 representation unspecified; new kernel
operation/Kind/expression-node need appears; result unclassifiable;
replay/hash unpreservable; human adjudication required; capability
outside §§6–9 inventory (`BLOCKED`, never invent abstraction).

Precedence: `specs_v2_3.md` authoritative; historical planning
examples non-normative (includes the old `OperationParams` example
inconsistency); freeze references preserved but PASS0 verification
required rather than silent authority.

Locks: single "Pin A:"; always `phys-gr`; "no GR metric/tensor
semantics" (never bare "no metrics"); kernel `Differentiate` vs
GR-local `Diff`; kernel `Limit` vs GR-local weak-field reduction;
retired "does not bypass" phrasing; corrected `Solve` (single
`Relation(eq,t²,rhs)` → `±Sqrt` branch set); valid cross-references
(hashing = §26/Appendix C); §2 direction unambiguous.

---

## 26. Canonicalization / Hashing (D32/D40)

Canonical JSON: fixed field order, deterministic arrays, exact
rational `num/den` strings, deterministic map serialization (maps
excluded from authoritative paths), UTF-8, no whitespace variance.
Numbers: integers plain; rationals as joined strings. Strings: ASCII
symbol names. Covers `GRExpr/Index/IndexSlot/Tensor/Metric/
Connection/Curvature/GRStep`. Canonical artifact data excludes
review metadata and UTC timestamps (external sidecar only).

`GRStep{StepID `gr-step-NNNNNN`, Index, OperationID, InputCanonicals,
InputHashes, ParamsCanonical, OutputCanonical, OutputHash,
KernelBridgeHash (empty when no bridge), PreviousStepHash,
CurrentStepHash}`; genesis `PreviousStepHash` = 64 zeros;
`CurrentStepHash = SHA-256(PreviousStepHash ||
CanonicalJSON(step-without-CurrentStepHash))`. Round-trip law:
`Decode(Canonical(x)) == x` before acceptance. Kernel `Object` hashes
recorded alongside (not merged) for bridged scalars.

---

## Appendix A — Decision Register D1–D40

Per header list. Uniqueness verified. New decisions D19–D40 carry
exact sections above; no decision contradicts another.

## Appendix B — Rejected (permanent)

Kernel tensors/trig/unknown-functions; extending `Differentiate`;
expanding `Limit`; general ODE/CAS/integration/asymptotics/registry/
prover; placeholders; rational-exponent calculus; efficacy ratio or
threshold; automatic L2; bare-symbol blame; GR inside `phys`; `ops`
expansion.

## Appendix C — Self-Audit Record

Re-read v2 + review set; checklist built per review finding
(resolved / rejected / human-action / spec ambiguity — all items in
§§21–28 of the review addressed above; architecture-redesign items
rejected per Appendix B). Searched for `TBD`, `to be decided`,
`define later`, `implementation-defined`, `coding agent may choose`,
`future decision`, `as needed` (zero occurrences in normative force);
illustrative language eliminated. Cross-references verified.
Decision IDs unique. Definition of Done (§20) checks every pin
(PASS0, bridge, Diff, trig/unknown-functions, classification,
SW suite, derivation + verification, flat + weak-field, tensor,
conventions, Fields 18–20, census).

---

ARCHITECTURAL_BLOCKERS = 0
UNPINNED_IMPLEMENTATION_CHOICES = 0
KERNEL_CHANGES_AUTHORIZED = 0
GR_WORKLOAD_IMPLEMENTED = 0
PASS0_REQUIRED = 1

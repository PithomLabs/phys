# Plan 11.1 Revised — GR Implementation Plan v4-F (Full Corrected Execution Contract)

Version: v4-F — 20261005.
Status: normative execution contract. Planning document only.
No code, corpus, spec, tree, test, manifest, or module changes made by this document.
This file is the SOLE normative execution contract for the coding agent.

Source: `plan10/GR_IMPLEMENTATION_PLAN_20261003.md` (v3, frozen lineage).
Hardening lineage: `plan10/GR_IMPLEMENTATION_PLAN_20261004.md` (v4-H plan, immutable).
Output: this file `plan10/GR_IMPLEMENTATION_PLAN_20261005.md` (v4-F).
Append-only versioning: v2, v3, v4-H artifacts are untouched.

Authority (normative precedence):
1. `plan10/specs_v2_3.md` (frozen) — highest authority.
2. This v4-F contract — execution hardening; where v4-F explicitly amends an
   implementation rule, v4-F governs execution.
3. v3 — immutable historical lineage only.
Use this file as the execution-building instruction. Attached review files
(Z / DeepSeek / Gemini / Qwen) are evidence/guidance only. This v4-F contract
is the sole normative execution contract.
Frozen `specs_v2_3.md` remains authoritative over both v3 and v4-F.

Changelog: v4-F applies the GR_PLAN4_plan.md build plan: V4-1 Subst correction,
V4-2 H3 HYPOTHESIS correction, V4-3 precedence, V4-4 bounded distribution+cap,
V4-5 stale-ref/language sweep, Reduce removal, assumption keys, IndexSlot
reconciliation, chart propagation, OperationID namespace, rational/AST guards,
adversarial layout, flat-space check, DeepSeek R1-R11, selective Gemini
procedures, Qwen frozen-API pins, 15-rule Implementation Discipline,
normative Appendix D, 26.1-26.3 pipeline. No architecture change, no kernel
change. D orders of magnitude: base D1-D40 retained as amended below.

Authoritative counts (single source of truth for audit):
- Base decisions D1-D40: 40 items.
- v4-F amendments V4-1 through V4-13: 13 items. Total decisions: 53.
- R1-R11 residuals are folded into V4 items/sections per Appendix A table;
  they are not a second count.
- Closing checklist: 12 items C1-C12 in Section 20. The number "12" below is
  the only checklist count; any other number elsewhere is non-normative.

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
  ├── GRExpr algebra: Diff / Normalize / ZeroTest / Subst / Truncate
  └── GR structures: Chart / Tensor / Metric / Connection / Curvature
  └── exact kernel subset → temporary HYPOTHESIS Object → frozen phys ops
  └── reintegrate GRExpr

Growth Gate: decides promotion from recorded evidence only.
```

Decision register: D1-D18 per v3 unchanged. D19-D40 per v3 with amendments:
D31 no longer contains Reduce (removed); D34 narrowed per V4-F Truncate rule;
Section 23/H3 Identify expectation corrected to HYPOTHESIS per frozen
Section 13.2 contamination law. V4-1..V4-13 defined in Appendix A and pinned
in normative sections below.

## 1. Executive Summary

GR is a hostile workload against frozen `phys`, not a kernel feature. Record
separately what `phys` can do, refuses, and what userland provides.
`NO-GROWTH` is acceptable success.
Gating: PASS0 (prerequisite) → PASS1 → GR-0…GR-8 + adversarial suite +
bridgeability census → Growth Gate. v4-F = v3 architecture with every
implementation choice pinned plus V4 corrections.

## 2. Architectural Principles

Per v3 (small core + ecosystem; kernel primitive capabilities only; minimal
composable invariant-preserving kernel; algebra manufactures no meaning;
userland first; Level direction 1→2→3 through Growth Gate only; invariants
not counts drive growth).

## 3. Frozen Baseline Assumptions

Per v3 items 1-8 (frozen spec/manifests; 42-file operational freeze verified
by PASS0; internal/kernel unreachable; core.Expr closed node set; trusted
objects via fixed constructors/hypothesis/ops/session; 12 operation IDs;
PASS0 baseline record). 42-count unchanged by GR.

## 4. Module Boundaries

`github.com/PithomLabs/phys` frozen vs `github.com/PithomLabs/phys-gr`
Level-1 one-way dep on public API only. mechanics/relativity grandfathered.
Local go.work at shared parent; CI pins phys commit; replace only to
pristine hash-verified checkout; no vendoring; no internal/kernel import.

## 5. GR Workload (GR-0…GR-8 locked thread)

Static spherical thread; linearized gravity deferred. Full execution tables
are in Section 24 (corrected; any Section 27 pointer is stale and void).

```
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

## 6. Representation Strategy (GRExpr algebra)

Node set closed to exactly 9 nodes; additions are STOP conditions:
Symbol / Rational / Add / Mul / Neg / Pow / Sin / Cos / UnknownFunction.
Immutable; exact fields per Appendix D.1; canonical child ordering;
structural equality over canonical form; integer exponents only (negative
permitted locally with nonzero-base evidence); max derivative order 2;
named errors: UnsupportedExponentError, UnsupportedNodeError,
DerivativeOrderExceededError, UnrepresentableKernelProjectionError,
NonZeroEvidenceError, UnrepresentableSubstitutionError.
Rational guards: zero is always 0/1; no silent fixed-width wrap (named error,
no int64 pin). Post-normalize arity: 0 terms → Rational(0); 1 term → unwrap;
>=2 → Add/Mul node. kind strings lowercase per frozen Section 10.2.
UnknownFunction pin: Name ASCII; Args ordered symbols, arity fixed (A,B arity
1, Args=[r]); DerivativeOrder len == len(Args), entries >= 0, sum <= 2.
A(r)=[0], A'(r)=[1], A''(r)=[2]; distinct orders are distinct atoms for
normalization, structurally related by Name+Args. Diff wrt non-argument → 0.
No generic Derivative node. A, A', A'' dependency semantics per Section 7.

## 7. Kernel Interaction Points + Diff/Subst/Truncate (V4-1 corrected)

Contact surface: constructors, read-only inspection, arithmetic +
Simplify/Subst-equivalent kernel ops on HYPOTHESIS wrappers, Differentiate
(GR-3a probe only), Compare, single-pattern Solve, constrained SelectBranch,
Session.Identify (KindRelation only). Limit never performs GR-8. No new
Kinds/nodes/factories/ops IDs. Reduce does not exist.
GR-local names: Diff, Normalize, ZeroTest, Subst, Truncate. Kernel names
Differentiate, Simplify, Limit, Solve, Substitute never used for GR-local
functions. Sqrt/BranchSet/Relation/Call results never import into GRExpr;
kernel Compare/Solve outputs remain core.Object artifacts.

Diff(e,x) table: Rational→0; Symbol(x)→1; Symbol(y≠x)→0; Add→Add(Diff);
Mul→n-ary product rule; Neg→Neg(Diff); Pow(u,n)→Mul(n,Pow(u,n-1),Diff(u,x))
integer n; Sin→Mul(Cos(u),Diff); Cos→Mul(Neg(Sin(u)),Diff);
UnknownFunction(f,args,orders): x in args and sum+1≤2 → incremented atom;
x not in args → 0; else DerivativeOrderExceededError. Dependency checklist:
Diff(A(r),t)=0; Diff(A(r),r)=A'(r); Diff(A'(r),r)=A''(r);
Diff(A''(r),r)=DerivativeOrderExceededError. GR-3a: bridge →
ops.Differentiate directly, never Session.Step.

Subst (V4-1 corrected contract):
Subst(e, target, replacement): target ∈ {Symbol, zero-order UnknownFunction
atom}; replacement = arbitrary GRExpr; composite target →
UnrepresentableSubstitutionError; UnknownFunction f→F rewrites f,f',f''
consistently via local Diff with order cap 2; Symbol subst never rewrites
UnknownFunction.Args; result = Normalize(result); negative-power evidence
enforced. Workload uses: Symbol→Symbol renaming; candidate certification
(A→1, B→1, B→1/A, A→1+k2/r); mu→eps*r weak-field; u→0 boundary.
UnknownFunction arguments themselves only Symbol→Symbol. Not a general engine.
Step 8 boundary evaluation is implemented through this Subst contract
(bounded boundary-condition substitution); no "only mechanical route" claim.

Truncate: closed-form Schwarzschild weak-field only; mu=eps*r substitution,
bounded algebraic reduction, retain O(eps^0..1); single narrow rule
Pow(1+u,-1)→1-u ONLY when u explicitly first-order in eps
(e.g. (1-2eps)^-1→1+2eps); no recursion to arbitrary powers; no
UnknownFunction expansion; no Taylor/engine/solver. No APPROXIMATED.

## 8. Level-1 GR Package Design (file plan)

phys-gr/ : go.mod; coordinates/chart.go; tensor/tensor.go; tensor/ops.go
(TensorAdd/Subtract/ScalarMul/Contract/RaiseIndex/LowerIndex/ApplyMetric/
ApplyInverseMetric); metric/metric.go; connection/connection.go (NOT Tensor);
curvature/curvature.go; vacuum/vacuum.go (VacuumAssertion collection);
solution/schwarzschild.go; limit/weakfield.go (SOLE canonical home; no
symbolic/weakfield.go); symbolic/expr.go; symbolic/differentiate.go;
symbolic/normalize.go; symbolic/bridge.go; replay/replay.go; evidence/
(failure_log.md, kernel_contact.md, growth_records/); adversarial evidence
as .md with tests adjacent to targets (see Section 16). Nothing else.
Tensor model per v3 D26 (Rank 0..4, 4^Rank dense row-major lexicographic,
SymmetryRule{Permutation,Sign}, validation + inheritance, Connection
rejected by type). D=4 stated explicitly. Chart canonical
ID=spherical-static, coords [t,r,theta,phi], Order=[0,1,2,3]; every
GR-0…GR-8 Metric/Tensor/Connection/Curvature artifact carries
ChartID=spherical-static; cross-chart only adversarial. Connection:
Components[rho][mu][nu] = Γ^rho_mu_nu; single global index ordering, no
second local ordering. Tensor ops: validate structure, ChartID, rank/slots,
variance legality, symmetries → allocate fresh output → compute → validate
output; never mutate inputs, never partial-write before validation.

## 9. GR-Local Mathematics Required

Diff (Section 7); Normalize + three-valued ZeroTest (Section 13);
Sin/Cos + workload identities; UnknownFunction calculus; Subst; Truncate;
tensor engine; inverse/contraction checks. Bounded to this workload only.

## 10. Failure Taxonomy

SPEC-INTENDED-BOUND / PACKAGE-SOLVABLE / REPRESENTABLE-BUT-UNFAITHFUL /
UNREPRESENTABLE / SILENTLY-WRONG (second-agent review mandatory with
classification + independent classification + disagreement + human
resolution).

## 11. C1–C8 Concern Ledger

Per v3 routing (C1-C4 Level-2; C5 future; C6 closed; C7-C8 PASS0 pins).

## 12. Silent-Wrongness Suite (SW-1…SW-13)

Per v3 with corrected semantics (bare Symbol→0 is ACCEPTED-CORRECT).
SW-4 Γ-as-tensor must reject; SW-11 negative controls; SW-12 round-trip;
SW-13 replay inconsistency.

## 13. Normalize + ZeroTest (one-way, bounded)

Normalize pipeline (fixed order, one-way, terminating):
1 numeric simplification; 2 flatten Add/Mul; 3 canonical child ordering;
4 rational/like-term collection; 5 pinned trig identities only
(Cos^2→1-Sin^2, Sin^2+Cos^2→1); 6 zero reduction. No reverse trig; no
unrestricted distributive expansion; no rule reintroduces an earlier-stage
form. Distribution permitted SOLELY for bounded common-denominator
normalization required by GR workload (e.g. R_tt/A+R_rr/B residual); bounded
by workload expression depth. Iterate to fixpoint with finite pass cap; cap
exceeded → UNDECIDED, never loop. UnknownFunction atoms compare by
(Name,Args,Order). ZeroTest ∈ {ZERO,NONZERO,UNDECIDED}: ZERO only when
Normalize→Rational(0); NONZERO only on direct structural proof; else
UNDECIDED. No numerics/sampling/heuristics. Test oracle is test-only
falsifier of false ZERO. Negative controls required.

## 14. Level-2 → Level-3 Growth Gate

Per v3 (full record + generic/theory-neutral/non-semantic/foundational/
reusable/invariant-bearing/unlivable-above proof; pre-registered falsifiers).

## 15. Evidence Record Format (Fields 18-20 executable)

Fields 1-17 per v3. Field 18 MinimalKernelChange (no source). Field 19
IndependentReview. Field 20 HumanApproval (UTC ISO-8601 Z). Timestamps and
reviewer identity (UTF-8 string) are external sidecar metadata ONLY; they
MUST NOT enter canonical JSON or any canonical hash. No approval →
PROVISIONAL → no kernel modification.

## 16. Test Strategy (PASS0 prerequisite)

PASS0 separate; GR agent verifies only; mismatch → STOP, ESCALATE-TO-SPEC.
Pin A Identify(RestMass,E/c²)→KindRelation/M/IDENTIFIED/NONE. Pin B
Differentiate(Pow(x,-1))→UnsupportedOperationError, target
ops/negative_test.go ONLY (ops/differentiate_test.go does not exist; never
create). Pin C AGENTS.md caveat. Fixture restoration of broken
Identify(Energy,E/c²) is a Section 17 disposable-copy mutant, never in-place
production mutation. GR tests: golden states per pass; GR-3a/3b pair (3b
golden = 9 independent non-zero Γ symmetry classes with Γ^rho_mu_nu =
Γ^rho_nu_mu in full dense array); trig ownership; componentwise
g·g^-1=δ golden (all 16 residual components via ZeroTest, Kronecker pattern);
curvature symmetry validation checks (Γ symmetric; Riemann antisymmetric;
Ricci symmetric); vacuum + Schwarzschild goldens; flat mu→0 check recovering
spherical Minkowski; Newtonian Phi AND a_r (a_r definite include, one Diff
call) via Truncate; contact ledger rows. Test packages: internal
`package <pkg>` default; `<pkg>_test` only for public-API isolation. Level-1
gate: `go test ./...` and `go vet ./...`. Adversarial tests adjacent to
targets; adversarial/ holds .md evidence (a test-only directory with no
non-test Go file fails go build).

## 17. Mutation Testing

Per v3 (disposable clone only) plus fixture mutant above.

## 18. Expected Outcomes

NO-GROWTH / LEVEL-2-CANDIDATE-PENDING / KERNEL-GROWTH-CANDIDATE-PENDING /
ESCALATE-TO-SPEC. No ratio/threshold/automatic promotion.

## 19. Explicit Non-Goals

Per v3 rejection list +: no Reduce invention; no int64 Rational pin; no
generic sorted-keys (fixed field order governs); no HYPOTHESIS-as-Kind
BridgePayload; no trig.go reopening; no general CAS/ODE/integration/
asymptotics/registry/prover/solver/series/Limit engine. SR corpus untouched.

## 20. Definition of Done (closing checklist C1-C12, authoritative)

C1 PASS0 record verified else STOP (Section 16). C2 phys-gr executes
GR-0…GR-8 tables (Section 24) modifying only phys-gr. C3 GR-3a direct-ops
failure + GR-3b independent Diff recorded. C4 real Sin/Cos/UnknownFunction,
kernel Call unchanged, no placeholders. C5 obstacles classified once
(Section 10), ledger + contact rows complete. C6 SW-1…SW-13 with corrected
semantics. C7 Schwarzschild derived (14 steps) + independently verified;
flat mu→0; weak-field Phi AND a_r via Truncate. C8 tensor/chart/convention/
unit pins honored; artifact hashes recorded. C9 Fields 18-20 complete;
independent review for non-PACKAGE-SOLVABLE; human approval for candidates.
C10 bridgeability census descriptive, no score. C11 final report
(phys-can / refuses / userland-provides); NO-GROWTH acceptable. C12 machine
closing audit block verified (footer of this file). Any Section 36 pointer
is stale and void; any Section 27 pointer for execution tables is stale and
void (tables are Section 24); hashing pointer is Section 26 only.

## 21. Bridge Contract (strict)

ToCoreExpr(GRExpr)→(core.Expr,error); FromCoreExpr(core.Expr)→(GRExpr,error);
ToCoreObject(GRExpr,Dimension,AssumptionSet,ConventionSet)→
(core.Object,bridgeID,error). Mirror subset Symbol/Rational/Add/Mul/Neg/
Pow(int) only. ToCoreExpr is STRICT: any Sin/Cos/UnknownFunction node →
UnrepresentableKernelProjectionError for that call (no silent partial
projection); caller-driven extraction of representable subexpressions only;
GR-local whole authoritative. FromCoreExpr same subset; Sqrt/Call/Relation/
BranchSet never import into GRExpr; kernel Sqrt/BranchSet results stay
kernel-side. Round-trip: FromCoreExpr(ToCoreExpr(t))==t for mirror terms.
BridgeID = gr/bridge/<lowercase-hex SHA-256 UTF-8 CanonicalJSON payload>.
Canonical payload field order: 1 expr, 2 kind, 3 dimension, 4 assumptions,
5 conventions. Payload is NOT a canonical Object; frozen Section 10.4 Object
order MUST NOT reorder it. No schema_version in payload. No UUID/timestamp/
counter/random. BridgePayload.assumptions = exact canonical representation
of retained structured AssumptionSet in frozen core canonical ordering; no
GR string serialization, no map iteration, no regenerated keys. Only
kernel-representable assumptions enter bridge; GR-local non-bridgeable
assumptions (e.g. A(r)!=0) remain GR-replay-only; never silently dropped or
converted. Assumptions caller-supplied structured core.Assumption only;
regime (physical, GR-local) vs mathematical preconditions recorded
separately. Lazy temporary HYPOTHESIS objects at deliberate call sites only.

## 22. Bridge Dimension / Units

Geometrized track GR-0…GR-8c: G=c=1, bridge core.Dimensionless();
meaning-loss recorded as GR limitation. SI audit GR-8d separate path. mu
(geometric length) vs M (SI mass) never conflated: mu=G*M/c^2, k2=-2*mu,
A=1-2*mu/r geometric; mu/r dimensionless pre-bridge. SI audit verifies
G*M/c^2→L, GM/(c^2 r)→1, GM/r^2→L/T^2.

## 23. GR-8 Procedures

MTW -+++, [t,r,theta,phi], Lambda=0:
Gamma^rho_mu_nu = 1/2 g^rhoσ(∂μ gνσ+∂ν gμσ-∂σ gμν);
R^rho_σμν = ∂μ Γ^rho_νσ-∂ν Γ^rho_μσ+Γ^rho_μλ Γ^λ_νσ-Γ^rho_νλ Γ^λ_μσ;
Rμν=R^rho_μρν; R=g^μν Rμν; Gμν=Rμν-1/2 R gμν. Vacuum Gμν=0 ⇔ Rμν=0.
R_θθ (lowered; bare R collides with Ricci scalar — never use bare R for the
θθ component): general ansatz R_θθ = 1-1/B+rB'/(2B^2)-rA'/(2AB); ONLY after
B=1/A: R_θθ = 1-A-rA'. 14-step script: (1) ansatz; (2) inverse; (3) Γ;
(4) R_tt,R_rr,R_θθ; (5) certify R_tt/A+R_rr/B=(AB)'/(rAB^2); (6) vacuum →
(AB)'=0; (7) candidate-and-certify AB=k1; (8) bounded boundary-condition
evaluation via Subst (assert A→1,B→1 into AB=k1; structural u=1/r,u=0 form
for closed-form checks; never r=infinity, never Limit) → k1=1; (9) B=1/A;
(10) certify R_θθ=1-A-rA'; (11) vacuum → (rA)'=1; (12) candidate-and-certify
rA=r+k2; (13) A=1+k2/r,B=1/A; (14) verify all vacuum. Candidate-and-certify
(mechanical): construct residual from authoritative state → Subst candidate
→ Normalize → ZeroTest; ZERO certifies; UNDECIDED fails closed; NONZERO
rejects; never certify by pattern match. PRIMARY owns its objects+trace;
SECONDARY builds closed-form A=1-2mu/r,B=1/A from scratch with separate
objects+trace, reusing no PRIMARY intermediates; SECONDARY never rescues
PRIMARY. Vacuum represented as GR-local VacuumAssertion{ComponentIndex,
Expression GRExpr, ZeroResult} collection; kernel Relation not used for
general A(r),B(r) equations. Newtonian: eps=G*M/(c^2 r);
Truncate(expr,eps,1); g_tt=-(1+2Φ/c^2)+O(eps^2); Phi=-GM/r; a_r=-GM/r^2.
Flat check: mu→0 reduces Schwarzschild to spherical Minkowski. GR-8 session
segment (sole phys.Session use, frozen API only; GR math stays in GR replay):
Compare weak-field g_tt expansion vs -(1+2Φ/c^2) form recorded via
Session.Step (Kind compare, Operator eq), lifecycle Draft→Identify→
Conclude, no Seal; constant correspondence via Session.Identify ONLY for
k2↔-2mu geometric (operands: k2 expression; -2mu expression), justification
required, result Kind=Relation, Provenance=HYPOTHESIS, CorpusStatus=NONE
(per frozen Section 13.2 contamination law; IDENTIFIED unreachable for
HYPOTHESIS operands). SI k2=-2GM/c^2 is dimensional-audit only, no Identify.
Exact call site in GR trace.

## 24. GR-0…GR-8 Execution Tables

Each pass: Objective/Inputs/Outputs/Data structures/Kernel contact/Level-1
ops/Worked example/Golden result/Negative tests/Expected failures/Evidence/
Classification/Exit criteria/Stop/Human review. GR-0 chart/conventions/units/
API audit (golden: canonical chart + ASCII + mu/M split + ToCoreExpr domain).
GR-1 ansatz diag(-A,B,r^2,r^2Sin^2θ) rank-2 Tensor + A/B UnknownFunction
(golden: component list + chart binding). GR-2 inverse + engine ops golden
g·g^-1=δ componentwise (Section 16). GR-3 3a UnsupportedOperationError
(SPEC-INTENDED-BOUND), 3b Diff Γ golden (9 symmetry classes); census
spherical vs x=cosθ descriptive. GR-4 covariant derivative golden
metric-compatibility. GR-5 Riemann (needs A''/B'' order 2) golden symmetries
+ components. GR-6 Ricci golden R_tt,R_rr,R_θθ + steps 4-5 certified. GR-7
Einstein + vacuum relations (structured, not solved). GR-8 14-step +
verification + mu→0 flat + Truncate correspondence + session segment; golden
A=1-2mu/r, Phi=-GM/r, a_r=-GM/r^2. T_pass overrun never growth evidence.

## 25. Stop Conditions + Spec Precedence + Wording Locks

STOP: baseline invalid; frozen API contradicts plan; required Level-1
representation unspecified; new kernel op/Kind/node need; unclassifiable
result; replay/hash unpreservable; human adjudication required; capability
outside Sections 6-9 inventory (BLOCKED, never invent). Precedence:
specs_v2_3.md > this v4-F > v3 lineage. Normative-language sweep (mandatory
validation): zero occurrences in normative force of TBD, to be decided,
define later, implementation-defined, coding agent may choose, future
decision, as needed, MAY, preferably, (if any), where representable, where
included — except quoted audit text. Locks: single Pin A:; always phys-gr;
kernel Differentiate vs GR Diff; kernel Limit vs GR weak-field; retired
does-not-bypass phrasing; Solve single Relation(eq,t^2,rhs)→±Sqrt; cross-refs:
execution tables = Section 24, hashing = Section 26, DoD checklist = Section
20; Section 2 direction unambiguous.

## 26. Canonicalization / Hashing

26.1 Pipeline (Section 13 order, one-way, capped). 26.2 Hash construction:
GRExpr→CanonicalJSON→UTF-8→SHA-256; GRStep CurrentStepHash =
SHA-256(string-concat PreviousStepHash ASCII-hex || UTF-8
CanonicalJSON(step-without-CurrentStepHash)) — string-concat form pinned
(R1); genesis PreviousStepHash = 64 zeros; BridgeID per Section 21.
Canonical JSON: fixed field order (NOT generic sorted-keys), deterministic
arrays, exact rational num/den strings, no maps in authoritative paths,
UTF-8, no whitespace variance, kind discriminator lowercase, ASCII names.
Each artifact type owns its encoder; no canonical/ package. Chart and
ConventionSet have NO standalone canonical artifact schema; when present in
a GR artifact they use the fixed embedded representation specified by that
artifact's Appendix D schema and are never independently hashed as GR kinds.
Artifact data excludes review metadata and UTC timestamps (sidecar only).
26.3 Round-trip + replay verification: Decode(Canonical(x))==x with
byte-identical re-encode before acceptance; per GRStep verify InputHash[i]==
SHA256(InputCanonical[i]), OutputHash==SHA256(OutputCanonical),
PreviousStepHash==prior CurrentStepHash, recomputed CurrentStepHash match,
then operation-specific replay. Kernel Object hashes recorded alongside.

## 27. Implementation Discipline — Non-Architectural Guidelines

Procedures, not architecture: (1) validate inputs before computation, never
mutate; (2) fresh outputs; (3) single global index ordering ([rho][mu][nu]
Connection, lexicographic row-major Tensor); (4) 9 Christoffel symmetry
classes with lower-index symmetry in dense array; (5) componentwise inverse
vs Kronecker delta; (6) Γ/Riemann/Ricci symmetry validation checks;
(7) candidate-and-certify loop only; (8) PRIMARY/SECONDARY isolation;
(9) recompute all hashes in replay; (10) decode/re-encode before accept;
(11) ChartID spherical-static everywhere GR-0…GR-8; (12) Subst bounded;
(13) Truncate bounded; (14) never add CAS/series/limit/solver/registry/hidden
rewrites; (15) STOP and report uncovered detail, never invent architecture.

## Appendix A — Decision Registers

D1-D40 per header as amended (D31 sans Reduce; D34 narrowed; H3 corrected).
V4-1 Subst correction (Sec 7). V4-2 HYPOTHESIS Identify (Secs 16,23).
V4-3 precedence (header, Sec 25). V4-4 distribution+cap (Sec 13).
V4-5 stale-ref/language sweep (Secs 5,20,24,25 + validation).
V4-6 Reduce removal (Secs 7-9,19). V4-7 assumption keys (Sec 21).
V4-8 IndexSlot inventory + embedded Chart/ConventionSet (Secs 8,26,D).
V4-9 chart propagation (Secs 8,16,27). V4-10 OperationID namespace (Secs
8,26,D). V4-11 rational/AST guards (Secs 6,13). V4-12 adversarial layout +
gates (Secs 8,16). V4-13 flat-space check (Secs 16,23). R1-R11 disposition:
R1→Sec 26.2; R2→Sec 16; R3→footer note; R4→Sec 20 C1-C12; R5→Secs 6,26,D;
R6→V4-8; R7→V4-8; R8→Sec 7 uses; R9→Secs 13,23; R10→Secs 20,25;
R11→V4-9. IDs unique; no decision contradicts another.

## Appendix B — Rejected (permanent)

Per v3 list plus: Reduce invention; int64 Rational pin; generic sorted-keys;
HYPOTHESIS-as-Kind payload; trig.go reopening; general expansion/CAS;
"only mechanical route" philosophy wording; IDENTIFIED for HYPOTHESIS
operands; r=infinity / Limit engine; silent partial projection.

## Appendix C — Self-Audit Record

Audited v3 line-by-line + V4-H plan + Z/DeepSeek/Gemini/Qwen reviews;
normative-language grep zero; cross-refs verified (tables Sec 24, hashing
Sec 26, checklist Sec 20); stale Sec 36/27 pointers removed; vague phrases
removed; counts verified (53 decisions, 12 checklist items); validation
per Section 20 C12 performed before setting footer flags.

## Appendix D — Canonical JSON Schemas (normative)

Fixed field order; no whitespace; UTF-8; deterministic arrays; no maps in
authoritative paths; kind discriminator lowercase.
D.1 GRExpr:
{"kind":"symbol","name":"r"}
{"kind":"rational","num":"1","den":"2"} (strings; zero "0"/"1")
{"kind":"add","terms":[]}
{"kind":"mul","factors":[]}
{"kind":"neg","expr":{}}
{"kind":"pow","base":{},"exp":"-1"} (integer string)
{"kind":"sin","expr":{}}
{"kind":"cos","expr":{}}
{"kind":"unknown_function","name":"A","args":["r"],"derivative_order":[0]}
Args symbols only; no expression args; no Derivative node.
D.2 IndexSlot: {"slot":0,"variance":"covariant|contravariant","domain":"spacetime"}
ordered by slot. D.3 Tensor: {"rank":2,"slots":[],"components":[],
"symmetries":[{"permutation":[],"sign":1}],"chart_id":"spherical-static"}
rank 0..4; components dense row-major lexicographic length 4^Rank with
explicit zeros; symmetries SymmetryRule list. D.4 Metric:
{"chart_id":"","coords":[],"g":[],"ginv":[]} with g/ginv Tensor payloads.
D.5 Connection (NOT Tensor):
{"chart_id":"spherical-static","components":[]} as [rho][mu][nu] nested
arrays of GRExpr canonical JSON. D.6 Curvature:
{"chart_id":"","kind":"riemann|ricci|scalar|einstein","components":[]|{}}
D.7 GRStep field order: step_id, index, operation_id, input_canonicals,
input_hashes, params_canonical, output_canonical, output_hash,
kernel_bridge_hash ("" when none), previous_step_hash, current_step_hash.
operation_id namespace (exact, closed): gr:diff, gr:normalize, gr:zero_test,
gr:subst, gr:truncate, gr:contract, gr:raise_index, gr:lower_index,
gr:apply_metric, gr:apply_inverse_metric, gr:tensor_add, gr:tensor_subtract,
gr:tensor_scalar_mul, gr:curvature, gr:vacuum_check, gr:bridge_project,
kernel:<frozen-ops-ID-lowercase> (e.g. kernel:differentiate, kernel:compare,
kernel:identify-ledger-only). No other OperationID values permitted.
input_canonicals, input_hashes, params_canonical, output_canonical,
output_hash, kernel_bridge_hash ("" when none), previous_step_hash,
current_step_hash. D.8 BridgePayload field order: expr, kind, dimension,
assumptions, conventions. D.9 Embedded Chart/ConventionSet: no standalone
schema; embedded fixed forms per owning artifact schema; never hashed alone.

---

PASS0_REQUIRED = 1
ARCHITECTURAL_BLOCKERS = 0
UNPINNED_IMPLEMENTATION_CHOICES = 0
KERNEL_CHANGES_AUTHORIZED = 0
GR_WORKLOAD_IMPLEMENTED = 0 (status flag only, not an invariant; retained
for v3 lineage continuity)

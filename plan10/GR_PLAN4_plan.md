# GR_PLAN4 Build Plan — v4-F full corrected contract (PLAN ONLY, no implementation)

Status: planning document only. Do NOT treat this file as the execution contract.
Implementation target (on next instruction): `plan10/GR_IMPLEMENTATION_PLAN_20261005.md` (v4-F).

## 0. Versioning / authority

- Preserve untouched: `GR_IMPLEMENTATION_PLAN_20261002.md` (v2),
  `GR_IMPLEMENTATION_PLAN_20261003.md` (v3 frozen lineage),
  `GR_IMPLEMENTATION_PLAN_20261004.md` (v4-H, 129-line hardening plan).
- On next instruction only: create `GR_IMPLEMENTATION_PLAN_20261005.md` (v4-F)
  as the single normative execution contract for the coding agent.
- v4-F header block:
  Source = 20261003; Hardening lineage = 20261004; Output = 20261005;
  Status = normative full v4 execution contract;
  Authority = specs_v2_3.md > v4-F > v3 lineage.
- Add V4-3 precedence clause: v3 immutable lineage; v4-F governs where it
  amends a rule; frozen specs_v2_3.md governs both.

## 1. Base build

- Copy v3 Sections 1-26 + Appendices A-C verbatim into v4-F draft, then patch
  v4-F only. Carry D1-D40 forward; add v4 delta IDs (V4-1..., R1-R11).
Use this file (`GR_PLAN4_plan.md`) as the execution-building instruction.
  The attached review files (Z / DeepSeek / Gemini / Qwen) are evidence/guidance only.
  The final 20261005 contract itself is the sole normative execution contract.

## 2. Critical contradictions (Z, must-fix)

- V4-1 Subst (refined): target = Symbol | zero-order UnknownFunction atom;
  replacement = arbitrary GRExpr; composite target -> error;
  f -> F rewrites f,f',f'' via local Diff, total order <= 2;
  Symbol subst never rewrites UnknownFunction.Args; result Normalize;
  negative-power evidence still enforced. Covers A->1, B->1, B->1/A,
  A->1+k2/r, mu->eps*r, u->0. Not a general CAS engine. Do NOT adopt Z's
  "only mechanical route" philosophy wording.
- V4-2 H3 provenance: GR-8 Identify(k2 <-> -2mu) -> Kind=Relation,
  Provenance=HYPOTHESIS, Corpus=NONE, justification required. Patch H3 item
  and Section 23 expected-result text. SI k2=-2GM/c^2 stays audit-only.

## 3. Z V4-3..V4-5 + items 6-13

- Bounded distribution only for common-denominator normalization; Normalize =
  fixed passes to fixpoint with finite cap; cap exceeded -> UNDECIDED.
- Stale-ref sweep: Sections 5, 8, 20 `-> Section 27` refs point to Section 24;
  extend banned-phrase grep (preferably also a_r, (if any), where
  representable, where included); a_r definite include.
- Remove Reduce from Section 7 / D31 / OperationID / GRStep / tests / tree /
  checklist. GR-local ops = Diff, Normalize, ZeroTest, Subst, Truncate.
- Bridge assumptions = exact canonical retained AssumptionSet in frozen core
  ordering; no string serialization, no map order, no regenerated keys.
- Canonical inventory = GRExpr, IndexSlot, Tensor, Metric, Connection,
  Curvature, GRStep, BridgePayload; drop bare Index unless defined.
- Chart propagation: all GR-0..GR-8 Metric/Tensor/Connection/Curvature carry
  ChartID = spherical-static; cross-chart only adversarial.
- OperationID namespace: gr:diff, gr:normalize, gr:zero_test, gr:subst,
  gr:truncate, gr:contract, gr:raise_index, gr:lower_index,
  gr:apply_metric, gr:apply_inverse_metric, ... + kernel:<frozen-op-id>.
- Rational/AST guards: zero always 0/1; post-normalize 0 terms -> Rational(0),
  1 term -> unwrap, >=2 -> node; no silent fixed-width wrap (no int64 pin).
- adversarial/: adjacent tests (internal package default, _test pkg only for
  public-API isolation); evidence as .md; gate `go test ./... && go vet ./...`.
- Flat-space check mu -> 0 -> spherical Minkowski as validation procedure.

## 4. DeepSeek R1-R11

- R1 hash || = string-concat (ASCII hex + UTF-8 canon) -> SHA-256.
- R2 test packages per above. R3 state GR_WORKLOAD_IMPLEMENTED=0 disposition.
- R4 enumerate checklist points explicitly. R5 kind lowercase per frozen 10.2.
- R6/R7 reconcile Index vs IndexSlot; Chart and ConventionSet have no standalone
  canonical artifact schema. When present in a GR artifact, they are serialized
  using their fixed embedded representation as specified by that artifact's
  Appendix D schema. They are never independently hashed as GR artifact kinds.
- R8 document Subst workload uses. R9 mechanical candidate-and-certify
  (residual -> Subst -> Normalize -> ZeroTest; ZERO certifies, UNDECIDED
  fails closed). R10 tighten B1 vs C4 scope. R11 single-chart pin.

## 5. Gemini selective + Qwen + discipline

- Adopt: PRIMARY/SECONDARY isolation; 9 Christoffel symmetry classes;
  8-step validate-before-compute + fresh outputs; single [rho][mu][nu]
  ordering; componentwise g*g^-1 = delta golden; curvature symmetry checks;
  full replay re-verification; one-way normalization; narrow (1-2eps)^-1 ->
  1+2eps at O(eps1); structural u=1/r,u=0 flatness; Diff dependency table.
- Reject: Rational int64 pin; generic sorted-keys; BridgePayload example with
  kind HYPOTHESIS / string assumptions; extra trig.go reopening tree.
- Qwen: GR-3a frozen-API sequence; GR-8 existing phys.Session only
  (Compare via Session.Step then Identify, Draft -> Identify -> Conclude).
- Append 15-rule Implementation Discipline section (non-architectural).

## 6. Validation

- Normative-language sweep, stale-ref grep, schema/hash/bridge/ZeroTest/
  Identify determinism, Subst/Truncate bounds, no Limit engine, no kernel
  change, 26.1-26.3 + normative Appendix D + audit-only Appendix C checks.
- Set PASS0_REQUIRED=1, ARCHITECTURAL_BLOCKERS=0,
  UNPINNED_IMPLEMENTATION_CHOICES=0, KERNEL_CHANGES_AUTHORIZED=0 only after
  audit passes.

Next: await explicit instruction to implement GR_IMPLEMENTATION_PLAN_20261005.md.

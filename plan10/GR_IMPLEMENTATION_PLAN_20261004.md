# GR v4 Hardening Plan (agreed, not yet executed)

Source: `plan10/GR_IMPLEMENTATION_PLAN_20261003.md` (v3, frozen)
Output: `plan10/GR_IMPLEMENTATION_PLAN_20261004.md` (v4, to be created)
Rule: append-only. Leave v3 byte-for-byte unchanged. Do not modify
v2, v3, `specs_v2_3.md`, manifests, or `phys`. No code changes in this step.

Status: planning document only. Awaiting explicit instruction before
implementing v4.

## Locked decisions

1. Output file: create `GR_IMPLEMENTATION_PLAN_20261004.md` (v4),
   header `v4 - 20261004`, changelog + lineage `v2 -> v3 -> v4`.
2. Structure split:
   - `26.1-26.3` = normative pipeline/hash/round-trip behavioral rules.
   - New normative `Appendix D` = exact canonical JSON schemas.
   - `Appendix C` = self-audit record only.
   - `Appendix D` is normative, not illustrative.
3. Bridge payload order: `expr, kind, dimension, assumptions, conventions`.
   Distinct from frozen `specs_v2_3.md 10.4` Object order. No `schema_version`.
4. Scope: 28 hardening items (16 original + 12 net-valid review findings).
   No new architecture. No CAS/ODE/Limit engine. No `Sqrt` in `GRExpr`.
   No kernel changes. No new Kinds/AST nodes.

## Work items (v4 agent contract)

### A. Setup + global audit
- A1. Copy v3 -> v4 verbatim, then patch v4 only.
- A2. Grep v4 for TBD / to be decided / define later / implementation-defined /
      coding agent may choose / future decision / as needed / MAY in normative
      force -> must be 0 after pass (except quoted audit text).

### B. Definition of Done (20 + footer)
- B1. Replace all stale `36 footer` / `20.12` refs (current `20.12`,
      `25` locks line, `Appendix C` checklist) with explicit closing audit
      block reference.
- B2. Footer must contain:
      `PASS0_REQUIRED=1, ARCHITECTURAL_BLOCKERS=0,`
      `UNPINNED_IMPLEMENTATION_CHOICES=0, KERNEL_CHANGES_AUTHORIZED=0`.

### C. Canonicalization pipeline (26.1-26.3)
- C1. `26.1`: deterministic normalization order only:
      1.numeric 2.flatten Add/Mul 3.canonical child ordering
      4.rational collection 5.pinned trig identities only 6.zero reduction.
      Distribution only where explicitly required by bounded GR workload.
      Reject general/full polynomial-expansion CAS.
- C2. `26.2`: hash construction:
      `GRExpr -> CanonicalJSON -> UTF-8 -> SHA-256`;
      `GRStep: SHA-256(PreviousStepHash || Canon(step-without-current))`;
      bridge per D.
- C3. `26.3`: round-trip acceptance `Decode(Canonical(x)) == x`.
- C4. Fix xref `hashing = 26/Appendix C` -> `hashing = 26` (25 + anywhere).
- C5. Pin canonicalization ownership: each artifact type owns its encoder;
      no new `canonical/` package.
- C6. Repeat in 15: Field-20 UTC timestamp + reviewer identity
      (`UTF-8 string`, `ISO-8601 UTC Z`) are sidecar-only, never hashed.

### D. Normative schemas (new Appendix D)
- D1. `D.1 GRExpr` 9 nodes, exact field order + JSON:
      `symbol{name}`, `rational{num,den}`, `add{terms}`, `mul{factors}`,
      `neg{expr}`, `pow{base,exp}`, `sin{expr}`, `cos{expr}`,
      `unknown_function{name,args,derivative_order}`.
      UTF-8, no whitespace variance, deterministic arrays, no maps in
      authoritative paths, `kind` discriminator, Args symbols-only,
      no generic Derivative node.
- D2. `D.2-D.8`: `IndexSlot, Tensor, Metric, Connection, Curvature,`
      `GRStep, BridgePayload` - each with field order, array/int/string/
      empty-value encoding.
- D3. `Connection{Components [4][4][4]GRExpr; ChartID}`,
      `Components[rho][mu][nu] = Gamma^rho_mu_nu`, not a Tensor.
- D4. `8` tree: list `go.mod + sources + adjacent _test.go + evidence`,
      nothing else. Sole weak-field home `limit/weakfield.go`;
      delete `symbolic/weakfield.go MAY` wording.

### E. Bridge contract (21)
- E1. Canonical `BridgePayload{expr,kind,dimension,assumptions,conventions}`,
      `BridgeID = gr/bridge/<hex SHA-256 UTF8 canon payload>`.
      No UUID/timestamp/counter/random.
- E2. Strict `ToCoreExpr`: any non-mirror node ->
      `UnrepresentableKernelProjectionError`. Caller-driven subexpression
      extraction only; no silent partial projection.
- E3. `Relation/BranchSet/Sqrt/Call` never import into `GRExpr`.
      Kernel `Compare`/`Solve` results stay `core.Object` artifacts.
- E4. Non-bridgeable assumptions (e.g. `A(r)!=0`) stay GR-replay-only;
      only kernel-representable assumptions enter bridge; never drop/convert
      silently.

### F. ZeroTest / Subst / Truncate / Step 8
- F1. `ZeroTest=ZERO` only if pinned `Normalize -> Rational(0)`;
      `NONZERO` only on structural proof; else `UNDECIDED`. No numerics.
- F2. `Subst`: allow `Symbol->Symbol` (e.g. `A(r)->A(x)`); composite ->
      `UnrepresentableSubstitutionError`.
- F3. `Truncate`: closed-form Schwarzschild only, `mu=eps*r` substitution,
      bounded reduction; single bounded rule `Pow(1+u,-1)->1-u` at O(eps1).
      No Taylor/UnknownFunction/arbitrary solver.
- F4. Step 8: structural `u=1/r, u->0` boundary evaluation only. No `Limit`.

### G. Schwarzschild wording (23-24)
- G1. Use `R_thth` lowered form only (document as R_theta-theta). General
      `R = 1-1/B+rB'/(2B^2)-rA'/(2AB)`; simplified `1-A-rA'` only after
      `B=1/A`.
- G2. Vacuum as GR-local `VacuumAssertion{ComponentIndex, GRExpr,
      ZeroResult}` collection; kernel `Relation` not used for general
      `A(r),B(r)` equations.

### H. Tests / GR-3a / GR-8 session (16, 23)
- H1. Pin B -> `ops/negative_test.go`,
      `Differentiate(Pow(x,-1)) -> UnsupportedOperationError`.
      Remove `ops/differentiate_test.go` refs.
- H2. GR-3a sequence:
      `GRExpr -> ToCoreObject(GRExpr,Dimension,AssumptionSet,ConventionSet)`
      `-> ops.Differentiate -> UnsupportedOperationError -> replay`,
      never `Session.Step`.
- H3. GR-8 `Identify`: operands `k2 <-> -2mu`, purpose constant
      correspondence only, result `KindRelation/NONE`. Session uses frozen
      `phys` API only; math stays in GR replay domain.

### I. Invariants + validation
- I1. Preserve `phys-gr -> phys` only; forbid `phys -> phys-gr`,
      kernel changes, new Kinds/AST nodes, GR semantics in kernel.
- I2. Self-audit v4 against 11-point checklist (no stale refs, GRExpr/artifact
      JSON deterministic, bridge/ZeroTest/Identify/test-locations
      deterministic, Subst/Truncate bounded, no Limit engine, no kernel
      changes, zero architectural choices remain).
- I3. Report files/sections/validation only. No further architecture.

## Next step
Await explicit user instruction to execute v4 implementation.

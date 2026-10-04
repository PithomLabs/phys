```python
import os
print("Current dir:", os.getcwd())
for root, dirs, files in os.walk("."):
    for f in files:
        if "GR" in f or "plan" in f or "plan10" in root or f.endswith(".md"):
            print(os.path.join(root, f))


```

```text
Current dir: /home/bard
./GR_IMPLEMENTATION_PLAN_20261004.md
./GR_IMPLEMENTATION_PLAN_20261003.md
./GR_IMPLEMENTATION_PLAN_20261002.md
./GR_PLAN_review.md
./plan11.1_overview.md
./GR_IMPLEMENTATION_PLAN.md
./plan11.1.md
./AGENTS.md
./specs_falsifiabilty.md
./plan7.md
./adv_review5.md
./prompt2.md
./plan10_v2_3.md
./README.md
./specs_v2_3.md
./plan9.2.md
./plan9.1_claude.md
./physics_compiler_mvp_specs_v2_2.md
./physics_compiler_mvp_implementation_plan_prompt_v2_2.md


```

```python
import os

files_to_check = [
    'GR_PLAN_review.md',
    'adv_review5.md',
    'GR_IMPLEMENTATION_PLAN_20261004.md',
    'GR_IMPLEMENTATION_PLAN_20261003.md'
]

for fname in files_to_check:
    if os.path.exists(fname):
        print(f"=== {fname} ===")
        with open(fname, 'r') as f:
            content = f.read()
            print(f"Length: {len(content)} chars")
            print(content[:1000])
            print("...\n")


```

```text
=== GR_PLAN_review.md ===
Length: 26390 chars
I read all six reviews end-to-end and reconciled them against the current `GR_IMPLEMENTATION_PLAN.md` and the frozen `specs_v2_3` API details. The reviewers are converging on the same conclusion: the architecture is sound, but the plan still leaves a few **actual implementation choices** implicit. Those must be pinned now so the coding agent cannot invent architecture while implementing.

The most important unresolved point is the GR ↔ `core.Object` boundary. The reviews independently identified that §6 says GR structures are not `core.Object`, while §7 simultaneously expects GR to use object-based `ops`; because every `ops` function operates on `core.Object`, the current document is literally unexecutable at that point. GR_PLAN_claude GR_PLAN_z

# Consolidated architectural lock

## 1. GR remains Level 1 and remains outside `phys`

This survives every review.

```text
[github.com/PithomLabs/phys](https://github.com/PithomLabs/phys)
    frozen Level 3 substrate

[github.com/PithomLabs/phys-gr](https://github.com/PithomLabs/phys-gr)
    Level 1 General Relativity

...

=== adv_review5.md ===
Length: 20605 chars
# PHYSICS CORPUS INTEGRITY AUDIT — ADVERSARIAL REVIEW

**Project:** Physics Compiler MVP (`[github.com/PithomLabs/phys](https://github.com/PithomLabs/phys)`)
**Normative spec:** `plan10/specs_v2_3.md`
**Implementation plan:** `plan10/plan10_v2_3.md`
**Audit date:** 2026-10-01
**Tests status:** All pass (`go test ./...` clean, build clean)

---

## 1. CORPUS EXECUTIVE VERDICT

| Category | Verdict |
|---|---|
| **Physics knowledge encoding** | **PASS** |
| **Mechanics** | **PASS** |
| **Special Relativity** | **PASS** |
| **Assumptions** | **PASS** |
| **Corpus integrity** | **PASS** |
| **Kernel separation** | **PASS** |

**Overall: PASS**

---

## 2. MECHANICS KNOWLEDGE TABLE

| Item | Expected expression | Actual expression | Expected dimension | Actual dimension | Expected assumptions | Actual assumptions | Provenance | Corpus status | Source | Verdict |
|---|---|---|---|---|---|---|---|---|---|---|
| Mass | Symbol("m") | Symbol("m") | M | M | none | none | DEFINED | ESTABLISHED | Newton, Principia | **PASS** |
| Time |
...

=== GR_IMPLEMENTATION_PLAN_20261004.md ===
Length: 6519 chars
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
4. Scope: 28 hardening
...

=== GR_IMPLEMENTATION_PLAN_20261003.md ===
Length: 32733 chars
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
  D19–D40, and integrates every pin into th
...



```

```python
import os

# Let's see all text in GR_PLAN_review.md and GR_IMPLEMENTATION_PLAN_20261004.md and any other review files
with open("GR_PLAN_review.md") as f:
    plan_review = f.read()

with open("GR_IMPLEMENTATION_PLAN_20261004.md") as f:
    plan_20261004 = f.read()

print("GR_PLAN_review sections:")
for line in plan_review.split("\n"):
    if line.startswith("#"):
        print(line)

print("\nGR_IMPLEMENTATION_PLAN_20261004.md content:")
print(plan_20261004)


```

```text
GR_PLAN_review sections:
# Consolidated architectural lock
## 1. GR remains Level 1 and remains outside `phys`
# 2. The correct GR/core.Object bridge
### Critical restriction
# 3. `phys-gr` gets one coherent bounded symbolic algebra
# 4. Use real `Sin`/`Cos`, not `sin_theta` placeholder symbols
# 5. Unknown functions are also genuine GR-local nodes
# 6. GR-3a / GR-3b is now unambiguous
# 7. GR-local operations do NOT become `phys.Session` operations
# 8. Rational normalization and zero testing are mandatory Level-1 machinery
# 9. Kernel `Simplify` should still be reused where applicable
# 10. Do not create a "Substrate Efficacy Ratio"
# 11. Schwarzschild should be derivation-first
# 12. Newtonian limit is NOT `phys.Limit`
### Exact flat-limit checkpoint
### Weak-field/Newtonian correspondence
# 13. Units: pin them, but keep tensor dimensions out of kernel semantics
# 14. Tensor structure must be concrete before GR-2
# 15. Conventions are part of the GR contract
# 16. Silent-wrongness suite: correct its semantics
# 17. Level-2 outcome must be changed
# 18. Growth evidence must be independently classified
# 19. Define the 20-field evidence artifacts precisely
# 20. GR artifact hashing
# 21. Baseline and mutation discipline
# 22. Local module mechanics
# 23. The frozen 42-file `phys` tree stays untouched by GR
# 24. Revised GR pass structure
# Final consolidated decisions
## Overall disposition of the six reviews

GR_IMPLEMENTATION_PLAN_20261004.md content:
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



```

```python
import os

# Check all review files
reviews = [f for f in os.listdir('.') if 'review' in f.lower() or 'plan' in f.lower()]
print("Matching files:", reviews)

for f in sorted(reviews):
    print(f"File: {f}, Size: {os.path.getsize(f)} bytes")


```

```text
Matching files: ['GR_IMPLEMENTATION_PLAN_20261004.md', 'GR_IMPLEMENTATION_PLAN_20261003.md', 'GR_IMPLEMENTATION_PLAN_20261002.md', 'GR_PLAN_review.md', 'plan11.1_overview.md', 'GR_IMPLEMENTATION_PLAN.md', 'plan11.1.md', 'plan7.md', 'adv_review5.md', 'plan10_v2_3.md', 'plan9.2.md', 'plan9.1_claude.md', 'physics_compiler_mvp_implementation_plan_prompt_v2_2.md']
File: GR_IMPLEMENTATION_PLAN.md, Size: 24564 bytes
File: GR_IMPLEMENTATION_PLAN_20261002.md, Size: 39482 bytes
File: GR_IMPLEMENTATION_PLAN_20261003.md, Size: 33149 bytes
File: GR_IMPLEMENTATION_PLAN_20261004.md, Size: 6519 bytes
File: GR_PLAN_review.md, Size: 26684 bytes
File: adv_review5.md, Size: 20892 bytes
File: physics_compiler_mvp_implementation_plan_prompt_v2_2.md, Size: 16035 bytes
File: plan10_v2_3.md, Size: 148323 bytes
File: plan11.1.md, Size: 6972 bytes
File: plan11.1_overview.md, Size: 9595 bytes
File: plan7.md, Size: 6406 bytes
File: plan9.1_claude.md, Size: 100557 bytes
File: plan9.2.md, Size: 120867 bytes


```

```python
with open("GR_PLAN_review.md") as f:
    text = f.read()

lines = text.split("\n")
print("GR_PLAN_review.md summary of headers and key sections:")
for i, line in enumerate(lines):
    if line.startswith("#"):
        print(f"Line {i}: {line}")


```

```text
GR_PLAN_review.md summary of headers and key sections:
Line 4: # Consolidated architectural lock
Line 6: ## 1. GR remains Level 1 and remains outside `phys`
Line 25: # 2. The correct GR/core.Object bridge
Line 84: ### Critical restriction
Line 94: # 3. `phys-gr` gets one coherent bounded symbolic algebra
Line 147: # 4. Use real `Sin`/`Cos`, not `sin_theta` placeholder symbols
Line 191: # 5. Unknown functions are also genuine GR-local nodes
Line 234: # 6. GR-3a / GR-3b is now unambiguous
Line 271: # 7. GR-local operations do NOT become `phys.Session` operations
Line 321: # 8. Rational normalization and zero testing are mandatory Level-1 machinery
Line 365: # 9. Kernel `Simplify` should still be reused where applicable
Line 391: # 10. Do not create a "Substrate Efficacy Ratio"
Line 437: # 11. Schwarzschild should be derivation-first
Line 494: # 12. Newtonian limit is NOT `phys.Limit`
Line 518: ### Exact flat-limit checkpoint
Line 526: ### Weak-field/Newtonian correspondence
Line 554: # 13. Units: pin them, but keep tensor dimensions out of kernel semantics
Line 577: # 14. Tensor structure must be concrete before GR-2
Line 632: # 15. Conventions are part of the GR contract
Line 662: # 16. Silent-wrongness suite: correct its semantics
Line 700: # 17. Level-2 outcome must be changed
Line 734: # 18. Growth evidence must be independently classified
Line 772: # 19. Define the 20-field evidence artifacts precisely
Line 813: # 20. GR artifact hashing
Line 849: # 21. Baseline and mutation discipline
Line 900: # 22. Local module mechanics
Line 929: # 23. The frozen 42-file `phys` tree stays untouched by GR
Line 956: # 24. Revised GR pass structure
Line 1011: # Final consolidated decisions
Line 1057: ## Overall disposition of the six reviews


```

```python
with open("adv_review5.md") as f:
    text = f.read()

lines = text.split("\n")
print("adv_review5.md summary of headers:")
for i, line in enumerate(lines):
    if line.startswith("#"):
        print(f"Line {i}: {line}")


```

```text
adv_review5.md summary of headers:
Line 0: # PHYSICS CORPUS INTEGRITY AUDIT — ADVERSARIAL REVIEW
Line 10: ## 1. CORPUS EXECUTIVE VERDICT
Line 25: ## 2. MECHANICS KNOWLEDGE TABLE
Line 46: ## 3. SPECIAL RELATIVITY KNOWLEDGE TABLE
Line 68: ## 4. ASSUMPTION AUDIT
Line 86: ## 5. FORMULA AUDIT
Line 100: ## 6. DERIVATION AUDIT
Line 102: ### E = mc² Derivation
Line 116: ### Kinetic Energy Derivative
Line 125: ### Lorentz Limit
Line 135: ## 7. KERNEL SEPARATION
Line 137: ### `internal/kernel/` — Physics occurrences
Line 149: ### `ops/` — Physics occurrences
Line 160: ### `core/` — Physics occurrences
Line 164: ### `mechanics/` and `relativity/`
Line 175: ## 8. CORPUS EXHAUSTIVENESS
Line 177: ### Expected vs. Implemented
Line 196: ### Non-manifest constructors (correctly excluded from manifest)
Line 207: ## 9. ADVERSARIAL CORRUPTION DETECTION
Line 245: ## 10. FINDINGS
Line 247: ### FIND-001: Minor — Mechanics manifest source uniformity
Line 258: ### FIND-002: Informational — Kernel contains spec-mandated physics identifiers
Line 268: ### FIND-003: Informational — Minkowski metric symbol identifier
Line 280: ## 11. FINAL VERDICT


```

The consolidated execution specification and architectural contract for the general relativity theory module (`phys-gr`) synthesizes all design decisions from v3 (`GR_IMPLEMENTATION_PLAN_20261003.md`) and hardening items from v4 (`GR_IMPLEMENTATION_PLAN_20261004.md`). Every architectural choice, mathematical rule, error condition, JSON schema, and execution boundary is explicitly pinned, leaving zero unpinned implementation decisions for the coding agent.

---

## 1. Executive Architectural Verdict & Module Boundaries

The kernel (`github.com/PithomLabs/phys`) remains strictly frozen, theory-agnostic, and math-agnostic. All general relativity calculations take place within the standalone Level-1 sibling module `github.com/PithomLabs/phys-gr`.

```
phys (Level 3 - Frozen Kernel)
 └── No GR metric, tensor, or transcendental semantics[cite: 1]
 └── Internal packages unreachable via Go compiler boundaries[cite: 1]
 └── Public surface only: core, ops, session, hypothesis[cite: 1]

phys-gr (Level 1 - Sibling Module)
 ├── GRExpr Algebra (9 closed immutable nodes)[cite: 1, 2]
 ├── Local Tensor Engine (Rank 0..4 dense row-major array)[cite: 1]
 ├── GR-Local Calculus (Diff, Normalize, ZeroTest, Subst, Truncate)[cite: 1]
 ├── HYPOTHESIS Bridge (ToCoreExpr / FromCoreExpr / ToCoreObject)[cite: 1]
 └── GR Replay & Artifact Serialization Engine[cite: 1, 2]

```

* **Module Dependency Rule**: Dependency flows strictly one-way from `phys-gr` to `phys` (`phys-gr -> phys`). Imports from `phys` into `phys-gr` are forbidden, and no code inside `phys` may reference `phys-gr`.


* **Operational Freeze**: The `phys` repository contains exactly 42 inventory files (39 implementation files + 3 documentation files). PASS0 verifies this inventory before any execution.


* **Workspace Setup**: Local development uses `go.work` in a common parent directory containing `phys/` and `phys-gr/` side-by-side. No source forks or vendoring are allowed.



---

## 2. Complete Decision Register (D1–D40 Summary)

| ID | Decision Summary | Normative Rule & Scope |
| --- | --- | --- |
| **D1** | Kernel Math-/Theory-Agnostic | Kernel provides primitive capability only; never dictates physical meaning.

 |
| **D2** | Level-1 Sibling Module | `phys-gr` created as an independent Go module.

 |
| **D3** | Grandfathered Packages | Existing `mechanics/` and `relativity/` remain untouched.

 |
| **D4** | GR Structures External | Metrics, connections, and tensors are not kernel `Kind`s.

 |
| **D5** | HYPOTHESIS Bridge | Projections use the frozen hypothesis constructor (`HYPOTHESIS`/`NONE`).

 |
| **D6** | Bounded GR-Local Algebra | `GRExpr` bounded strictly to the closed node set required for static spherical GR.

 |
| **D7** | Genuine Transcendental Nodes | `Sin`, `Cos`, and `UnknownFunction` exist natively in `GRExpr`.

 |
| **D8** | No Placeholder Variables | String placeholders (e.g., `sin_theta`) are explicitly forbidden.

 |
| **D9** | GR-3a Before GR-3b | Probe kernel `ops.Differentiate` bound (GR-3a) before executing local GR-3b.

 |
| **D10** | Independent Local Ops | Local calculus (`Diff`, `Normalize`) operates independently of kernel ops.

 |
| **D11** | Local Replay Chain | Replay trace recorded in `phys-gr` with SHA-256 state hashes.

 |
| **D12** | Kernel Mirror Reuse | Mirror-subset expressions project to kernel ops when exact representation exists.

 |
| **D13** | Concrete Tensor Engine | Dense row-major storage for rank 0..4 tensors with symmetry rules.

 |
| **D14** | Derivation-First Schwarzschild | 14-step candidate-and-certify derivation precedes verification.

 |
| **D15** | Newtonian Reduction | Local `Truncate` extracts potential $\Phi = -GM/r$ and $a_r = -GM/r^2$.

 |
| **D16** | No Efficacy Ratio | Code reuse percentage or efficiency ratios never trigger kernel growth.

 |
| **D17** | Candidate Pending State | Level-2 proposals stay in `LEVEL-2-CANDIDATE-PENDING`.

 |
| **D18** | Growth Gate Adjudication | Kernel growth requires evidence, independent review, and human approval.

 |
| **D19** | Closed Node Set | Exactly 9 immutable `GRExpr` node types.

 |
| **D20** | `UnknownFunction` Representation | `{Name, Args, DerivativeOrder}` with total order sum $\le 2$.

 |
| **D21** | Complete `Diff` Table | Explicit rule table covering product, chain, power, and unknown derivatives.

 |
| **D22** | Bridge Triple Contract | `ToCoreExpr`, `FromCoreExpr`, and `ToCoreObject` with content-derived IDs.

 |
| **D23** | Dimension & Unit Policy | Geometrized dimensionless bridge ($\mu = GM/c^2$, $k_2 = -2\mu$) with separate SI audit.

 |
| **D24** | Assumption Translation | Structured caller-supplied assumptions; regime vs. precondition split.

 |
| **D25** | Static Spherical Chart Pin | Canonical chart `spherical-static` with coordinates $[t, r, \theta, \phi]$.

 |
| **D26** | Tensor Model | Dense array ($4^{\text{Rank}}$ components), 8 required ops; Connection is not a Tensor.

 |
| **D27** | MTW Curvature Sign Conventions | Sign convention $(- + + +)$, MTW curvature formulas, $\Lambda = 0$.

 |
| **D28** | Schwarzschild Script | 14-step machine-certified candidate-and-certify script.

 |
| **D29** | Dual Artifact Verification | Primary derivation and secondary verification maintain independent criteria.

 |
| **D30** | Three-Valued `ZeroTest` | Returns `ZERO`, `NONZERO`, or `UNDECIDED`.

 |
| **D31** | Disjoint Naming | GR local names (`Diff`, `Normalize`, etc.) never shadow kernel names.

 |
| **D32** | Replay Hash Protocol | `gr-step-NNNNNN` format, 64-zero genesis, payload SHA-256 chaining.

 |
| **D33** | Direct-`ops` Probe Path | GR-3a calls `ops.Differentiate` directly, bypassing `Session.Step`.

 |
| **D34** | Weak-Field `Truncate` | Bounded substitution $\mu = \epsilon r$ retaining $O(\epsilon^0.. \epsilon^1)$ terms.

 |
| **D35** | Session Recording Scope | `phys.Session` restricted strictly to final scalar correspondence checks.

 |
| **D36** | Terminal States | `NO-GROWTH`, `LEVEL-2-CANDIDATE-PENDING`, `KERNEL-GROWTH-CANDIDATE-PENDING`, `ESCALATE-TO-SPEC`.

 |
| **D37** | Provisional Escalation Path | Escalation requires pre-registered falsifiers.

 |
| **D38** | Independent Second-Agent Review | Non-`PACKAGE-SOLVABLE` classifications require dual-agent consensus.

 |
| **D39** | Executable Evidence Fields | Fields 18–20 strictly structured for machine validation.

 |
| **D40** | Canonical JSON Serialization | Fixed key order, deterministic encoding, timestamps excluded from hash payload.

 |

---

## 3. Expression Algebra (`GRExpr`) & Node Specifications

The `GRExpr` type is an immutable Go interface. The node set is strictly closed to 9 types; any additional node type constitutes a stop condition.

```
GRExpr = Symbol | Rational | Add | Mul | Neg | Pow | Sin | Cos | UnknownFunction

```

### Node Data Structures

1. **`Symbol`**: `{ Name string }` (ASCII identifier, e.g., `"r"`, `"t"`, `"theta"`, `"phi"`, `"A"`, `"B"`, `"M"`, `"mu"`, `"G"`, `"c"`, `"eps"`, `"k1"`, `"k2"`).


2. **`Rational`**: `{ Num int64, Den int64 }` where $\text{Den} > 0$ and $\gcd(\vert{}\text{Num}\vert{}, \text{Den}) = 1$.


3. **`Add`**: `{ Terms []GRExpr }` ($\ge 2$ terms, canonical order).


4. **`Mul`**: `{ Factors []GRExpr }` ($\ge 2$ factors, canonical order).


5. **`Neg`**: `{ X GRExpr }`.


6. **`Pow`**: `{ Base GRExpr, Exp int64 }` (Exponents are strictly integers; negative exponents permitted locally).


7. **`Sin`**: `{ X GRExpr }`.


8. **`Cos`**: `{ X GRExpr }`.


9. **`UnknownFunction`**: `{ Name string, Args []Symbol, DerivativeOrder []uint8 }`.


* `Name`: ASCII string (e.g., `"A"`, `"B"`).


* `Args`: Slice of symbols defining function arity (e.g., `[Symbol("r")]`).


* `DerivativeOrder`: Vector matching `Args` length. Entries $\ge 0$, total sum $\le 2$.


* Representational mapping: $A(r) \implies [0]$, $A'(r) \implies [1]$, $A''(r) \implies [2]$.





### Canonical Order and Simplification Pipeline

Normalizing a `GRExpr` uses a 6-stage deterministic pipeline:

1. **Numeric Simplification**: Evaluate explicit rational arithmetic ($1/2 + 1/3 \to 5/6$, $a^0 \to 1$, $a^1 \to a$).


2. **Flattening**: Flatten nested `Add` and `Mul` nodes (`Add(Add(a, b), c) -> Add(a, b, c)`).


3. **Canonical Child Ordering**: Sort children of `Add` and `Mul` lexicographically by their canonical JSON byte representations.


4. **Rational Collection**: Combine like-term algebraic factors and scalar coefficients.


5. **Trigonometric Reduction**: Apply bounded trig identities ($\cos^2(x) \to 1 - \sin^2(x)$, $\sin^2(x) + \cos^2(x) \to 1$).


6. **Zero Reduction**: Eliminate zero terms in `Add` ($x + 0 \to x$) and absorb zero factors in `Mul` ($x \cdot 0 \to 0$).



---

## 4. Differentiation, Normalization & Three-Valued `ZeroTest`

### Complete `Diff(e, x)` Differentiation Table

`Diff` operates purely within Level 1 and is distinct from kernel `phys.Differentiate`.

$$\begin{aligned} \text{Diff}(\text{Rational}(p, q), x) &\to 0 \\ \text{Diff}(\text{Symbol}(y), x) &\to \begin{cases} 1 & \text{if } y == x \\ 0 & \text{if } y \neq x \end{cases} \\ \text{Diff}(\text{Add}(u_1, \dots, u_n), x) &\to \text{Add}(\text{Diff}(u_1, x), \dots, \text{Diff}(u_n, x)) \\ \text{Diff}(\text{Mul}(u_1, \dots, u_n), x) &\to \sum_{i=1}^n \text{Mul}\left( \text{Diff}(u_i, x), \prod_{j \neq i} u_j \right) \\ \text{Diff}(\text{Neg}(u), x) &\to \text{Neg}(\text{Diff}(u, x)) \\ \text{Diff}(\text{Pow}(u, n), x) &\to \text{Mul}(n, \text{Pow}(u, n-1), \text{Diff}(u, x)) \\ \text{Diff}(\text{Sin}(u), x) &\to \text{Mul}(\text{Cos}(u), \text{Diff}(u, x)) \\ \text{Diff}(\text{Cos}(u), x) &\to \text{Mul}(\text{Neg}(\text{Sin}(u)), \text{Diff}(u, x)) \\ \text{Diff}(\text{UnknownFunction}(f, \mathbf{args}, \mathbf{orders}), x) &\to \begin{cases}  \text{UnknownFunction}(f, \mathbf{args}, \mathbf{orders} + \mathbf{e}_k) & \text{if } x = \mathbf{args}[k] \text{ and } \sum \text{orders} + 1 \le 2 \\ 0 & \text{if } x \notin \mathbf{args} \\ \text{Error}(\text{DerivativeOrderExceededError}) & \text{if } \sum \text{orders} + 1 > 2 \end{cases} \end{aligned}$$

### Three-Valued `ZeroTest` Logic

`ZeroTest(e)` evaluates expression equivalence to zero without floating-point or numerical sampling:

```
               ┌────────────────────────────────────────┐
               │         Normalize(expression)          │
               └───────────────────┬────────────────────┘
                                   │
         ┌─────────────────────────┼────────────────────────┐
         ▼                         ▼                        ▼
┌─────────────────┐       ┌─────────────────┐      ┌─────────────────┘
│  Rational(0/1)  │       │ Direct Structural│      │  Unresolved /   │
│   Exact Form    │       │ Non-Zero Proof  │      │ Bounded Identity│
└────────┬────────┘       └────────┬────────┘      └────────┬────────┘
         │                         │                        │
         ▼                         ▼                        ▼
       ZERO                     NONZERO                 UNDECIDED
```[cite: 1, 2]

* **`ZERO`**: Returned if and only if `Normalize(e)` reduces structurally to `Rational(0, 1)`[cite: 1, 2].
* **`NONZERO`**: Returned if and only if `Normalize(e)` yields a non-zero constant rational, or a single non-zero symbol/function atom with non-zero coefficient[cite: 1, 2].
* **`UNDECIDED`**: Returned when symbolic cancellation cannot be proven using the 6-stage normalization rules[cite: 1, 2].

---

## 5. Bridge Protocol & Serialization Schemas

### HYPOTHESIS Bridge Contracts

The bridge converts expressions between `GRExpr` and `core.Expr` for kernel testing[cite: 1].

```go
ToCoreExpr(expr GRExpr) (core.Expr, error)
FromCoreExpr(coreExpr core.Expr) (GRExpr, error)
ToCoreObject(expr GRExpr, dim core.Dimension, assumptions core.AssumptionSet, conv ConventionSet) (core.Object, string, error)
```[cite: 1]

* **Mirror Subset**: Projections via `ToCoreExpr` are restricted to `Symbol`, `Rational`, `Add`, `Mul`, `Neg`, and `Pow(int64)`[cite: 1]. Attempting to project `Sin`, `Cos`, or `UnknownFunction` returns `UnrepresentableKernelProjectionError`[cite: 1, 2].
* **Bridge ID Generation**: Deterministic SHA-256 hash derived over the canonical bridge payload[cite: 1, 2]:

$$\text{BridgeID} = \text{"gr/bridge/"} + \text{Hex}(\text{SHA-256}(\text{CanonicalJSON}(\text{BridgePayload})))$$
[cite: 1, 2]

### Bridge Payload Canonical Structure

The canonical payload for generating a Bridge ID requires exact field ordering[cite: 2]:

```json
{
  "expr": { "kind": "rational", "num": 1, "den": 1 },
  "kind": "HYPOTHESIS",
  "dimension": "Dimensionless",
  "assumptions": ["r > 0"],
  "conventions": { "metric_signature": "-+++" }
}
```[cite: 2]

### Normative JSON Schemas (Appendix D)

All JSON artifacts are serialized without whitespace, with sorted keys, and without trailing commas[cite: 1, 2].

#### D.1 `GRExpr` JSON Representations

```json
{"kind": "symbol", "name": "r"}
{"kind": "rational", "num": 1, "den": 2}
{"kind": "add", "terms": [...]}
{"kind": "mul", "factors": [...]}
{"kind": "neg", "expr": {...}}
{"kind": "pow", "base": {...}, "exp": -1}
{"kind": "sin", "expr": {...}}
{"kind": "cos", "expr": {...}}
{"kind": "unknown_function", "name": "A", "args": ["r"], "derivative_order": [1]}
```[cite: 2]

#### D.2 `Tensor` JSON Schema

```json
{
  "rank": 2,
  "slots": [
    {"slot": 0, "variance": "covariant", "domain": "spacetime"},
    {"slot": 1, "variance": "covariant", "domain": "spacetime"}
  ],
  "components": [...],
  "symmetries": [
    {"permutation": [1, 0], "sign": 1}
  ],
  "chart_id": "spherical-static"
}
```[cite: 2]

#### D.3 `Connection` JSON Schema (Not a Tensor)

```json
{
  "chart_id": "spherical-static",
  "components": [
    [
      [{...}, {...}, {...}, {...}],
      [{...}, {...}, {...}, {...}],
      [{...}, {...}, {...}, {...}],
      [{...}, {...}, {...}, {...}]
    ]
  ]
}
```[cite: 2]

---

## 6. Tensor Engine, Chart & Curvature Physics

### Chart Definition (D25)

The primary coordinate chart is pinned as `spherical-static`[cite: 1]:
* **Coordinates**: $x^0 = t$, $x^1 = r$, $x^2 = \theta$, $x^3 = \phi$[cite: 1].
* **Symbols**: `"t"`, `"r"`, `"theta"`, `"phi"`[cite: 1].

### Metric Ansatz (GR-1)

The general static spherical metric tensor $g_{\mu\nu}$ is diagonal[cite: 1]:

$$g_{\mu\nu} = \begin{pmatrix} -A(r) & 0 & 0 & 0 \\ 0 & B(r) & 0 & 0 \\ 0 & 0 & r^2 & 0 \\ 0 & 0 & 0 & r^2 \sin^2\theta \end{pmatrix}$$
[cite: 1]

Inverse metric $g^{\mu\nu}$ components[cite: 1]:

$$g^{\mu\nu} = \text{diag}\left(-\frac{1}{A(r)}, \frac{1}{B(r)}, \frac{1}{r^2}, \frac{1}{r^2 \sin^2\theta}\right)$$
[cite: 1]

### MTW Curvature Formulas & Golden Signs (D27)

* **Metric Signature**: $(-, +, +, +)$[cite: 1, 2].
* **Christoffel Symbols**:

$$\Gamma^\rho_{\mu\nu} = \frac{1}{2} g^{\rho\sigma} \left( \frac{\partial g_{\nu\sigma}}{\partial x^\mu} + \frac{\partial g_{\mu\sigma}}{\partial x^\nu} - \frac{\partial g_{\mu\nu}}{\partial x^\sigma} \right)$$
[cite: 1]

* **Riemann Curvature Tensor**:

$$R^\rho_{\ \sigma\mu\nu} = \frac{\partial \Gamma^\rho_{\nu\sigma}}{\partial x^\mu} - \frac{\partial \Gamma^\rho_{\mu\sigma}}{\partial x^\nu} + \Gamma^\rho_{\mu\lambda} \Gamma^\lambda_{\nu\sigma} - \Gamma^\rho_{\nu\lambda} \Gamma^\lambda_{\mu\sigma}$$
[cite: 1]

* **Ricci Tensor**: $R_{\mu\nu} = R^\rho_{\ \mu\rho\nu}$[cite: 1].
* **Ricci Scalar**: $R = g^{\mu\nu} R_{\mu\nu}$[cite: 1].
* **Einstein Tensor**: $G_{\mu\nu} = R_{\mu\nu} - \frac{1}{2} R g_{\mu\nu}$[cite: 1].
* **Vacuum Equations**: $G_{\mu\nu} = 0 \iff R_{\mu\nu} = 0$ (for $\Lambda = 0$)[cite: 1].

#### Lowered $R_{\theta\theta}$ Formulation

The full expression for $R_{\theta\theta}$ prior to metric substitution is[cite: 2]:

$$R_{\theta\theta} = 1 - \frac{1}{B(r)} + \frac{r B'(r)}{2 B(r)^2} - \frac{r A'(r)}{2 A(r) B(r)}$$
[cite: 2]

When $B(r) = \frac{1}{A(r)}$, this simplifies to the golden form[cite: 1, 2]:

$$R_{\theta\theta} = 1 - A(r) - r A'(r) = 1 - (r A(r))'$$
[cite: 1, 2]

---

## 7. 14-Step Schwarzschild Derivation, Verification & Weak-Field Limits

The Schwarzschild solution is produced via a 14-step candidate-and-certify script[cite: 1].


```

┌─────────────────────────────────────────────────────────────────────────┐
│              14-Step Machine-Certified Derivation Script                │
├─────────────────────────────────────────────────────────────────────────┤
│  1. Define Metric Ansatz: g_μν = diag(-A(r), B(r), r², r² sin²θ)        │
│  2. Calculate Inverse Metric: g^μν = diag(-1/A, 1/B, 1/r², 1/(r² sin²θ))│
│  3. Compute Christoffel Symbols: Γ^ρ_μν                                 │
│  4. Compute Non-Zero Ricci Components: R_tt, R_rr, R_θθ                 │
│  5. Certify Linear Combination: R_tt/A + R_rr/B = (AB)' / (r A B²)      │
│  6. Impose Vacuum Condition: R_tt/A + R_rr/B = 0 => (AB)' = 0            │
│  7. Certify Intermediate Solution: A(r) B(r) = k₁                       │
│  8. Apply Boundary Condition (r -> ∞): A -> 1, B -> 1 => k₁ = 1         │
│  9. Substitute Reciprocal Constraint: B(r) = 1 / A(r)                   │
│ 10. Certify Angular Component: R_θθ = 1 - A(r) - r A'(r)               │
│ 11. Impose Angular Vacuum Condition: (r A(r))' = 1                      │
│ 12. Certify Radial Solution: r A(r) = r + k₂ => A(r) = 1 + k₂/r         │
│ 13. Construct Metric Solutions: A(r) = 1 + k₂/r, B(r) = (1 + k₂/r)⁻¹    │
│ 14. Verify All Vacuum Components: ZeroTest(R_μν) == ZERO                │
└─────────────────────────────────────────────────────────────────────────┘

```[cite: 1]

### Primary Derivation vs. Secondary Verification (D29)

* **Primary Derivation**: Executes Steps 1–14 starting from unknown $A(r)$ and $B(r)$[cite: 1]. Success requires every certified step to evaluate to `ZeroTest(ZERO)`[cite: 1].
* **Secondary Verification**: Constructs $A(r) = 1 - \frac{2\mu}{r}$ and $B(r) = \left(1 - \frac{2\mu}{r}\right)^{-1}$ independently, without reusing primary intermediates, and verifies that $R_{\mu\nu} = 0$ for all components[cite: 1].

### Weak-Field Reduction & Newtonian Correspondence (D34)

Weak-field reduction uses `Truncate(expr, eps, 1)` where $\mu = \epsilon \cdot r$[cite: 1, 2]:

$$g_{tt} = -\left(1 - \frac{2\mu}{r}\right) = -\left(1 - \frac{2 \epsilon r}{r}\right) = -(1 - 2\epsilon) = -\left(1 + \frac{2\Phi}{c^2}\right) + O(\epsilon^2)$$
[cite: 1, 2]

This extracts the Newtonian gravitational potential $\Phi$ and radial acceleration $a_r$[cite: 1]:

$$\Phi = -\frac{G M}{r}, \quad a_r = -\frac{\partial \Phi}{\partial r} = -\frac{G M}{r^2}$$
[cite: 1]

---

## 8. Replay Protocol, Artifact Hashing & Execution Tables (GR-0…GR-8)

### Replay Step Structure (`GRStep`)

Every step in `phys-gr` appends a deterministic `GRStep` record to the execution trace[cite: 1, 2]:

```go
type GRStep struct {
    StepID            string    // "gr-step-000001"
    Index             uint64    // Sequential index starting at 1
    OperationID       string    // e.g., "gr.diff", "gr.normalize"
    InputCanonicals   []string  // Canonical JSON strings of inputs
    InputHashes       []string  // SHA-256 hashes of input canonicals
    ParamsCanonical   string    // Canonical JSON of parameters
    OutputCanonical   string    // Canonical JSON of output
    OutputHash        string    // SHA-256 of OutputCanonical
    KernelBridgeHash  string    // SHA-256 of BridgePayload (or empty)
    PreviousStepHash  string    // SHA-256 of preceding GRStep
    CurrentStepHash   string    // SHA-256 of current GRStep
}
```[cite: 1]

The genesis step uses a `PreviousStepHash` of 64 zero characters (`0000000000000000000000000000000000000000000000000000000000000000`)[cite: 1].

### Step Hash Calculation Formula

$$\text{CurrentStepHash} = \text{SHA-256}(\text{PreviousStepHash} \parallel \text{CanonicalJSON}(\text{GRStep}_{\text{without CurrentStepHash}}))$$
[cite: 1, 2]

### Execution Thread Table (GR-0…GR-8)

| Pass ID | Stage Objective | Key Input / Operation | Golden Outcome |
| :--- | :--- | :--- | :--- |
| **GR-0** | API & Chart Audit | Chart `spherical-static`[cite: 1] | Coordinates $[t,r,\theta,\phi]$ validated[cite: 1]. |
| **GR-1** | Metric Ansatz | Unknown $A(r), B(r)$[cite: 1] | $g_{\mu\nu} = \text{diag}(-A, B, r^2, r^2 \sin^2\theta)$ constructed[cite: 1]. |
| **GR-2** | Inverse & Contraction | $g_{\mu\nu}$ tensor[cite: 1] | $g^{\mu\alpha} g_{\alpha\nu} = \delta^\mu_\nu$ verified[cite: 1]. |
| **GR-3a** | Kernel Differentiation Probe | `ToCoreObject(1/r)`[cite: 1, 2] | Probes kernel; returns `UnsupportedOperationError`[cite: 1, 2]. |
| **GR-3b** | Local Christoffel Computation | `Diff` over metric components[cite: 1] | 9 non-zero $\Gamma^\rho_{\mu\nu}$ components computed[cite: 1]. |
| **GR-4** | Covariant Derivative | $\nabla_\alpha g_{\mu\nu}$[cite: 1] | Metric compatibility $\nabla_\alpha g_{\mu\nu} = 0$ verified[cite: 1]. |
| **GR-5** | Riemann Tensor | $\partial \Gamma + \Gamma \Gamma$[cite: 1] | Riemann components $R^\rho_{\ \sigma\mu\nu}$ generated[cite: 1]. |
| **GR-6** | Ricci Tensor & Scalar | Contraction over Riemann[cite: 1] | $R_{tt}, R_{rr}, R_{\theta\theta}, R_{\phi\phi}$ computed[cite: 1]. |
| **GR-7** | Einstein Tensor & Vacuum | $G_{\mu\nu} = R_{\mu\nu} - \frac{1}{2} R g_{\mu\nu}$[cite: 1] | Formulates $R_{\mu\nu} = 0$ system[cite: 1]. |
| **GR-8a** | 14-Step Derivation | Candidate-and-certify script[cite: 1] | Certified derivation of $A(r) = 1 - 2\mu/r$[cite: 1]. |
| **GR-8b** | Independent Verification | Closed-form metric[cite: 1] | Direct verification that $R_{\mu\nu} = 0$[cite: 1]. |
| **GR-8c** | Flat-Spacetime Limit | $M \to 0$ ($\mu \to 0$)[cite: 1] | Recovers Minkowski metric $\eta_{\mu\nu}$[cite: 1]. |
| **GR-8d** | Weak-Field & Session Link | `Truncate` at $O(\epsilon^1)$[cite: 1, 2] | Recovers $\Phi = -GM/r$; session records `IDENTIFIED`[cite: 1, 2]. |

---

## 9. Silent-Wrongness Test Suite (SW-1…SW-13) & Failure Taxonomy

### Obstacle Failure Taxonomy

Every execution obstacle must be classified into one of five mutually exclusive categories[cite: 1]:

1. `SPEC-INTENDED-BOUND`: Expected operation boundary defined in frozen spec[cite: 1].
2. `PACKAGE-SOLVABLE`: Issue resolvable completely inside `phys-gr`[cite: 1].
3. `REPRESENTABLE-BUT-UNFAITHFUL`: Kernel accepts expression but alters mathematical semantics[cite: 1].
4. `UNREPRESENTABLE`: Mathematical structure cannot be projected into kernel[cite: 1].
5. `SILENTLY-WRONG`: Computation completes without error but produces incorrect physics[cite: 1].

### Silent-Wrongness Adversarial Test Suite

| Test ID | Targeted Failure Mode | Test Input / Condition | Required Guard Mechanism |
| :--- | :--- | :--- | :--- |
| **SW-1** | Symbol/Entity Ambiguity | Bare `Symbol("f")` differentiated wrt `r`[cite: 1] | `Diff` treats non-argument symbol as constant ($0$)[cite: 1]. |
| **SW-2** | Unknown Function Dependency | `UnknownFunction("A", [r])` diff wrt $t$[cite: 1] | Evaluates to $0$ since $t \notin \text{Args}$[cite: 1]. |
| **SW-3** | Chart Coordinate Mixing | Combining components from different charts[cite: 1] | Rejects tensor ops with mismatched `ChartID`[cite: 1]. |
| **SW-4** | Connection Treated as Tensor | Passing $\Gamma^\rho_{\mu\nu}$ to tensor contraction[cite: 1] | Type-system separation blocks non-tensor inputs[cite: 1, 2]. |
| **SW-5** | Invalid Index Contraction | Contracting two contravariant indices[cite: 1] | Asserts variance pair is covariant/contravariant[cite: 1]. |
| **SW-6** | Unverified Tensor Symmetry | Assuming symmetry without declared rule[cite: 1] | Enforces explicit `SymmetryRule` checks[cite: 1]. |
| **SW-7** | Semantic Dimension Mismatch | Adding geometric length $\mu$ to SI mass $M$[cite: 1] | Explicit separation of $\mu$ (length) and $M$ (mass)[cite: 1]. |
| **SW-8** | Type-Decay Addition | Adding rank-2 tensor to scalar[cite: 1] | Strict type checks prevent mixed rank additions[cite: 1]. |
| **SW-9** | Assumption Inconsistency | Combining $r > 0$ with $r = 0$ assumptions[cite: 1] | Caller-owned `AssumptionSet` validation[cite: 1]. |
| **SW-10** | Metric Signature Mismatch | Mixing $(+ - - -)$ and $(-, +, +, +)$[cite: 1] | Pin metric signature globally to $(-, +, +, +)$[cite: 1, 2]. |
| **SW-11** | False Zero Confirmation | Complex expression wrongly evaluating to $0$[cite: 1] | `ZeroTest` returns `UNDECIDED` if not proven[cite: 1, 2]. |
| **SW-12** | Round-Trip Canonical Shift | `FromCoreExpr(ToCoreExpr(e)) != e`[cite: 1] | Strict round-trip equality assertion[cite: 1, 2]. |
| **SW-13** | Replay Hash Mismatch | Altering step inputs without updating step hash[cite: 1] | Cryptographic hash chain validation[cite: 1, 2]. |

---

## 10. Final Execution Checklist & Zero-Ambiguity Hardening Verification

Before handoff to the coding agent, the repository state must verify all closing audit conditions[cite: 1, 2]:


```

[✓] PASS0_REQUIRED = 1                     (42/42 repository inventory files verified)
[✓] ARCHITECTURAL_BLOCKERS = 0             (Zero unresolved architectural choices)
[✓] UNPINNED_IMPLEMENTATION_CHOICES = 0    (All math, API, and schemas strictly pinned)
[✓] KERNEL_CHANGES_AUTHORIZED = 0          (Kernel phys repository remains byte-for-byte frozen)

```

### Module File Tree (`phys-gr/`)

The implementation for `phys-gr` is contained within the following file structure[cite: 1, 2]:


```

phys-gr/
├── go.mod
├── coordinates/
│   ├── chart.go             # Chart structure and coordinate definitions
│   └── chart_test.go
├── tensor/
│   ├── tensor.go            # Rank 0..4 dense row-major tensor model
│   ├── ops.go               # Contract, raise/lower, apply metric ops
│   └── tensor_test.go
├── metric/
│   ├── metric.go            # Metric tensor and inverse calculation
│   └── metric_test.go
├── connection/
│   ├── connection.go        # Christoffel symbols (non-tensor structure)
│   └── connection_test.go
├── curvature/
│   ├── curvature.go         # Riemann, Ricci, and Einstein calculations
│   └── curvature_test.go
├── vacuum/
│   ├── vacuum.go            # Vacuum field equation assertions
│   └── vacuum_test.go
├── solution/
│   ├── schwarzschild.go     # 14-step Schwarzschild derivation and verification
│   └── schwarzschild_test.go
├── limit/
│   ├── weakfield.go         # Truncate and Newtonian limit extraction
│   └── weakfield_test.go
├── symbolic/
│   ├── expr.go              # GRExpr 9-node algebra and canonical ordering
│   ├── differentiate.go     # Complete Diff function and lookup table
│   ├── normalize.go         # 6-stage normalization and ZeroTest
│   ├── bridge.go            # ToCoreExpr, FromCoreExpr, and ToCoreObject
│   └── symbolic_test.go
├── replay/
│   ├── replay.go            # GRStep hashing and replay chain ledger
│   └── replay_test.go
├── adversarial/
│   └── silent_wrongness_test.go  # Executable SW-1...SW-13 test suite
└── evidence/
├── failure_log.md       # Categorized failure taxonomy log
├── kernel_contact.md    # Kernel interaction records
└── growth_records/      # Fields 18-20 candidate records

```[cite: 1, 2]

```
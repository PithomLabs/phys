```python
with open('GR_IMPLEMENTATION_PLAN_20261003.md', 'r') as f:
    plan_v3 = f.read()

print("Plan v3 length:", len(plan_v3))
print("Plan v3 line count:", len(plan_v3.splitlines()))


```

```text
Plan v3 length: 32733
Plan v3 line count: 731


```

```python
import glob

files = glob.glob('**/*.go', recursive=True)
print("Go files:", files)

for f in sorted(files):
    with open(f, 'r') as fp:
        lines = fp.readlines()
        print(f"{f}: {len(lines)} lines")


```

```text
Go files: ['relations.go', 'primitives.go', 'manifest_test.go', 'derivation_test.go']
derivation_test.go: 455 lines
manifest_test.go: 319 lines
primitives.go: 203 lines
relations.go: 151 lines


```

```python
import glob, os

for root, dirs, files in os.walk('.'):
    # skip .git or __pycache__ if any
    if '.git' in root or '__pycache__' in root:
        continue
    for f in files:
        p = os.path.join(root, f)
        print(p, os.path.getsize(p))


```

```text
./GR_IMPLEMENTATION_PLAN_20261003.md 33149
./GR_IMPLEMENTATION_PLAN_20261002.md 39482
./GR_PLAN_review.md 26684
./plan11.1_overview.md 9595
./GR_IMPLEMENTATION_PLAN.md 24564
./plan11.1.md 6972
./relations.go 6162
./primitives.go 6843
./manifest_test.go 10324
./manifest.json 8893
./derivation_test.go 14973
./AGENTS.md 11686
./specs_falsifiabilty.md 8740
./plan7.md 6406
./adv_review5.md 20892
./prompt2.md 25681
./plan10_v2_3.md 148323
./README.md 6028
./specs_v2_3.md 112006
./plan9.2.md 120867
./plan9.1_claude.md 100557
./physics_compiler_mvp_specs_v2_2.md 109762
./physics_compiler_mvp_implementation_plan_prompt_v2_2.md 16035


```

```python
with open('GR_IMPLEMENTATION_PLAN_20261003.md', 'r') as f:
    text = f.read()

import re

sections = re.split(r'\n(?=## )', text)
for i, s in enumerate(sections):
    lines = s.strip().splitlines()
    header = lines[0] if lines else ""
    print(f"Section {i}: {header} ({len(lines)} lines)")


```

```text
Section 0: # Plan 11.1 Revised — GR Implementation Plan v3 (Final Pre-Handoff Execution Contract) (124 lines)
Section 1: ## 1. Executive Summary (21 lines)
Section 2: ## 2. Architectural Principles (18 lines)
Section 3: ## 3. Frozen Baseline Assumptions (22 lines)
Section 4: ## 4. Module Boundaries (13 lines)
Section 5: ## 5. GR Workload (GR-0…GR-8 locked thread) (20 lines)
Section 6: ## 6. Representation Strategy (GRExpr algebra — D19/D20) (40 lines)
Section 7: ## 7. Kernel Interaction Points + Diff table (D21/D31/D33) (41 lines)
Section 8: ## 8. Level-1 GR Package Design (file plan — D26/D32) (49 lines)
Section 9: ## 9. GR-Local Mathematics Required (8 lines)
Section 10: ## 10. Failure Taxonomy (8 lines)
Section 11: ## 11. C1–C8 Concern Ledger (6 lines)
Section 12: ## 12. Silent-Wrongness Testing Strategy (SW-1…SW-13) (18 lines)
Section 13: ## 13. Level-1 → Level-2 Promotion Rules (18 lines)
Section 14: ## 14. Level-2 → Level-3 Growth Gate (10 lines)
Section 15: ## 15. Evidence Record Format (Fields 18–20 executable) (13 lines)
Section 16: ## 16. Test Strategy (PASS0 prerequisite) (21 lines)
Section 17: ## 17. Mutation Testing (10 lines)
Section 18: ## 18. Expected Outcomes (10 lines)
Section 19: ## 19. Explicit Non-Goals (8 lines)
Section 20: ## 20. Definition of Done (21 lines)
Section 21: ## 21. Bridge Contract (D22–D24) (30 lines)
Section 22: ## 22. Bridge Dimension / Units (D23) (13 lines)
Section 23: ## 23. GR-8 Procedures (D27–D29/D34/D35) (50 lines)
Section 24: ## 24. GR-0…GR-8 Execution Tables (D25–D29/D33/D34) (37 lines)
Section 25: ## 25. Stop Conditions + Spec Precedence + Wording Locks (21 lines)
Section 26: ## 26. Canonicalization / Hashing (D32/D40) (20 lines)
Section 27: ## Appendix A — Decision Register D1–D40 (4 lines)
Section 28: ## Appendix B — Rejected (permanent) (7 lines)
Section 29: ## Appendix C — Self-Audit Record (21 lines)


```

```python
# Let's inspect Sections 6, 7, 8, 21, 22, 23, 24, 26 in detail
def print_sec(idx):
    print(f"=== SECTION {idx} ===")
    print(sections[idx])
    print("\n")

for i in [6, 7, 8, 21, 22, 23, 24, 26]:
    print_sec(i)


```

```text
=== SECTION 6 ===
## 6. Representation Strategy (GRExpr algebra — D19/D20)

Node set closed to exactly 9 nodes. Additions are STOP conditions:

```text
Symbol{name}  Rational{p,q reduced, q>0}  Add[terms]  Mul[factors]
Neg[x]  Pow{base, n int}  Sin[x]  Cos[x]  UnknownFunction{...}

```

1. Immutable representation. All constructors return new values.
2. Exact fields: `Symbol{name ASCII}`; `Rational{num int, den int>0, gcd=1}`; `Add{terms []GRExpr ≥2}`; `Mul{factors []GRExpr ≥2}`;
`Neg{x GRExpr}`; `Pow{base GRExpr, exp int}`; `Sin{x}`;
`Cos{x}`; `UnknownFunction{Name, Args []Symbol, DerivativeOrder []uint8}`.
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

=== SECTION 7 ===

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

=== SECTION 8 ===

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

Tensor model (D26): `Tensor{Rank 0..4, Slots []IndexSlot, Components []GRExpr dense row-major length 4^Rank lexicographic, Symmetries []SymmetryRule{Permutation, Sign ±1}, ChartID}`; `IndexSlot{slot, variance covariant|contravariant, domain=spacetime}`; component keys
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

=== SECTION 21 ===

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

=== SECTION 22 ===

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

=== SECTION 23 ===

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

=== SECTION 24 ===

## 24. GR-0…GR-8 Execution Tables (D25–D29/D33/D34)

Each pass carries: Objective / Inputs / Outputs / Data structures /
Kernel contact / Level-1 operations / Worked example / Golden result /
Negative tests / Expected failures / Evidence artifacts /
Classification / Exit criteria / Stop condition / Human review point.

* GR-0: audit chart/conventions/units/public API. Golden: canonical
chart + ASCII names + `mu`/`M` split recorded; `ToCoreExpr` domain
table acknowledged. Stop: any missing pin.
* GR-1: ansatz metric `diag(−A,B,r²,r²Sin²θ)` as `Tensor` rank 2 +
`UnknownFunction` A/B. Kernel: none (unrepresentable whole;
projectable rational coefficients where useful). Golden: component
list + chart binding.
* GR-2: inverse + engine ops (`Contract/Raise/Lower/ApplyMetric/ ApplyInverseMetric`) with validation. Golden: `g·g⁻¹ = δ`.
Negatives: chart mismatch, variance error, bad pairing rejected.
* GR-3: 3a direct `ops.Differentiate(Pow(r,−1))` ⇒
`UnsupportedOperationError` (`SPEC-INTENDED-BOUND`); 3b
`Diff`-based Γ golden set. Census: spherical vs `x=cosθ`
representability (descriptive).
* GR-4: covariant derivative via `Diff` + Γ + `Normalize`. Golden:
metric-compatibility spot check.
* GR-5: Riemann via D27 formula (needs `A''/B''`, order 2). Golden:
symmetries + selected components.
* GR-6: Ricci contraction + scalar. Golden: `R_tt,R_rr,R_thth`
forms; script steps 4–5 certified.
* GR-7: Einstein + vacuum relations `Gμν = 0` / `Rμν = 0`. Golden:
structured relations (not solved here).
* GR-8: 14-step derivation + independent verification + `M→0` flat
check + `Truncate` correspondence + session recording segment.
Golden: `A = 1−2mu/r`, `Φ = −GM/r`. Stop: any general-solver need.

`T_pass` wall-clock budgets are human-set execution observations;
overrun never constitutes growth evidence.

---

=== SECTION 26 ===

## 26. Canonicalization / Hashing (D32/D40)

Canonical JSON: fixed field order, deterministic arrays, exact
rational `num/den` strings, deterministic map serialization (maps
excluded from authoritative paths), UTF-8, no whitespace variance.
Numbers: integers plain; rationals as joined strings. Strings: ASCII
symbol names. Covers `GRExpr/Index/IndexSlot/Tensor/Metric/ Connection/Curvature/GRStep`. Canonical artifact data excludes
review metadata and UTC timestamps (external sidecar only).

`GRStep{StepID `gr-step-NNNNNN`, Index, OperationID, InputCanonicals, InputHashes, ParamsCanonical, OutputCanonical, OutputHash, KernelBridgeHash (empty when no bridge), PreviousStepHash, CurrentStepHash}`; genesis `PreviousStepHash` = 64 zeros;
`CurrentStepHash = SHA-256(PreviousStepHash || CanonicalJSON(step-without-CurrentStepHash))`. Round-trip law:
`Decode(Canonical(x)) == x` before acceptance. Kernel `Object` hashes
recorded alongside (not merged) for bridged scalars.

---

```



```python?code_reference&code_event_index=6
import sympy as sp

t, r, theta, phi = sp.symbols('t r theta phi', real=True)
A = sp.Function('A')(r)
B = sp.Function('B')(r)

# Metric g_uv
# x^0=t, x^1=r, x^2=theta, x^3=phi
# signature - + + +
g = sp.Matrix([
    [-A, 0, 0, 0],
    [0, B, 0, 0],
    [0, 0, r**2, 0],
    [0, 0, 0, r**2 * sp.sin(theta)**2]
])

g_inv = g.inv()

coords = [t, r, theta, phi]
dim = 4

# Christoffel symbols Gamma^rho_mu_nu
Gamma = sp.MutableDenseNDimArray.zeros(dim, dim, dim)
for rho in range(dim):
    for mu in range(dim):
        for nu in range(dim):
            val = 0
            for sigma in range(dim):
                val += sp.Rational(1, 2) * g_inv[rho, sigma] * (
                    sp.diff(g[nu, sigma], coords[mu]) +
                    sp.diff(g[mu, sigma], coords[nu]) -
                    sp.diff(g[mu, nu], coords[sigma])
                )
            Gamma[rho, mu, nu] = sp.simplify(val)

# Ricci tensor R_mu_nu = R^rho_mu_rho_nu
# MTW sign convention:
# R^rho_sigma_mu_nu = d_mu Gamma^rho_nu_sigma - d_nu Gamma^rho_mu_sigma + Gamma^rho_mu_lambda Gamma^lambda_nu_sigma - Gamma^rho_nu_lambda Gamma^lambda_mu_sigma
# R_mu_nu = R^rho_mu_rho_nu

R = sp.Matrix.zeros(dim, dim)
for mu in range(dim):
    for nu in range(dim):
        val = 0
        for rho in range(dim):
            # d_rho Gamma^rho_nu_mu - d_nu Gamma^rho_rho_mu + sum_lambda (Gamma^rho_rho_lambda Gamma^lambda_nu_mu - Gamma^rho_nu_lambda Gamma^lambda_rho_mu)
            d1 = sp.diff(Gamma[rho, nu, mu], coords[rho])
            d2 = sp.diff(Gamma[rho, rho, mu], coords[nu])
            t1 = sum(Gamma[rho, rho, lam] * Gamma[lam, nu, mu] for lam in range(dim))
            t2 = sum(Gamma[rho, nu, lam] * Gamma[lam, rho, mu] for lam in range(dim))
            val += d1 - d2 + t1 - t2
        R[mu, nu] = sp.simplify(val)

print("R_00 (R_tt):", R[0,0])
print("R_11 (R_rr):", R[1,1])
print("R_22 (R_thth):", R[2,2])
print("R_33 (R_phph):", R[3,3])

# Check Step 5: R_tt / A + R_rr / B
lhs = sp.simplify(R[0,0] / A + R[1,1] / B)
print("\nR_tt / A + R_rr / B =", lhs)

# Compare with (A*B)' / (r * A * B^2)
rhs = sp.simplify(sp.diff(A*B, r) / (r * A * B**2))
print("(A*B)' / (r * A * B^2) =", rhs)
print("Difference:", sp.simplify(lhs - rhs))


```

```text
R_00 (R_tt): Derivative(A(r), (r, 2))/(2*B(r)) - Derivative(A(r), r)*Derivative(B(r), r)/(4*B(r)**2) - Derivative(A(r), r)**2/(4*A(r)*B(r)) + Derivative(A(r), r)/(r*B(r))
R_11 (R_rr): -Derivative(A(r), (r, 2))/(2*A(r)) + Derivative(A(r), r)*Derivative(B(r), r)/(4*A(r)*B(r)) + Derivative(A(r), r)**2/(4*A(r)**2) + Derivative(B(r), r)/(r*B(r))
R_22 (R_thth): r*Derivative(B(r), r)/(2*B(r)**2) - r*Derivative(A(r), r)/(2*A(r)*B(r)) + 1 - 1/B(r)
R_33 (R_phph): (r*A(r)*Derivative(B(r), r)/2 - r*B(r)*Derivative(A(r), r)/2 + A(r)*B(r)**2 - A(r)*B(r))*sin(theta)**2/(A(r)*B(r)**2)

R_tt / A + R_rr / B = Derivative(B(r), r)/(r*B(r)**2) + Derivative(A(r), r)/(r*A(r)*B(r))
(A*B)' / (r * A * B^2) = Derivative(B(r), r)/(r*B(r)**2) + Derivative(A(r), r)/(r*A(r)*B(r))
Difference: 0


```

The revised plan **`GR_IMPLEMENTATION_PLAN_20261003.md` (v3)** is exceptionally thorough and succeeds in pinning down the macro-architecture, module boundaries, and execution gates. However, a strict line-by-line adversarial review reveals **8 critical ambiguities and structural edge cases** that will cause a coding agent to stall, diverge, or make non-deterministic design decisions during execution.

---

## 1. Critical Architectural Ambiguities & Edge Cases

### A. The `Sqrt` / Fractional Exponent Paradox (D19 vs D28 / §21 vs §25)

* **The Issue:** D19 restricts `GRExpr` exponents strictly to `int` (`Pow{base, exp int}`), and §21 specifies that `FromCoreExpr` rejects `Sqrt` with `UnrepresentableKernelProjectionError`. However, §25 mentions returning `±Sqrt` from kernel `Solve`, and metric derivations frequently encounter square-root relationships.
* **Agent Failure Mode:** If `Solve` returns a `core.Expr` containing `Sqrt`, `FromCoreExpr` will throw `UnrepresentableKernelProjectionError`, breaking the candidate-and-certify workflow in Step 7/13 of GR-8.
* **Required Pin:** Clarify that all metric components in `phys-gr` must be formulated strictly in polynomial or rational form (e.g., solving for $A$ and $B$, or $A^2$ where appropriate) so `Sqrt` is never imported into `GRExpr`. Explicitly define how `FromCoreExpr` handles kernel branch sets containing `Sqrt` (e.g., unwrapping $x^2 = c \implies x = \text{candidate}$ without importing `Sqrt`).

### B. `UnknownFunction` Argument Immutability vs `Subst` (D20 / §6 / §7)

* **The Issue:** D20 pins `UnknownFunction` as `{Name string, Args []Symbol, DerivativeOrder []uint8}`. Notice that `Args` is strictly `[]Symbol`, **not** `[]GRExpr`.
* **Agent Failure Mode:** If `Subst(expr, var, replacement)` is called where `var` is a coordinate in `Args` and `replacement` is an expression (e.g., $r \to r' + k$ or $r \to e^\rho$), the resulting `UnknownFunction` cannot be represented in the 9-node algebra.
* **Required Pin:** Explicitly pin the domain of `Subst` for `UnknownFunction`: `Subst` on `UnknownFunction` is restricted strictly to **Symbol-to-Symbol renaming** (e.g., $r \to r'$). Any attempt to substitute a composite `GRExpr` into an `UnknownFunction` argument must immediately return an `UnrepresentableSubstitutionError`.

### C. `Truncate` Series Expansion Boundary (D34 / §23 / §9)

* **The Issue:** D34 defines `Truncate(expr, eps, max_order)` as retaining $O(\epsilon^0 \dots \epsilon^1)$. However, §9 explicitly bans an "asymptotics engine" or general series expansion.
* **Agent Failure Mode:** For Schwarzschild $g_{tt} = -(1 - 2\mu/r)$, $g_{tt}$ is polynomial in $\mu$. But for $g_{rr} = B(r) = (1 - 2\mu/r)^{-1} = \text{Pow}(\text{Add}(1, \text{Neg}(\dots)), -1)$, $\epsilon = \mu/r$ is inside a negative power. Standard polynomial filtering cannot truncate $(1 - \epsilon)^{-1}$ into $1 + \epsilon + O(\epsilon^2)$ without series expansion.
* **Required Pin:** Specify whether `Truncate` performs exact geometric series expansion for degree $-1$ binomials of the form $(1 + u)^{-1} \to 1 - u$, or if `Truncate` operates strictly on expanded polynomial forms in `GRExpr`.

### D. Asymptotic Limit Verification in Step 8 without `Limit` (§23 / §7 / §9)

* **The Issue:** §23 Step 8 requires certifying "asymptotic flatness ($A, B \to 1$ as $r \to \infty$) $\Rightarrow k_1 = 1$". However, §7 and §9 state that neither the kernel `Limit` nor a Level-1 limit engine is available.
* **Agent Failure Mode:** The coding agent will attempt to write an asymptotic evaluator or stall when trying to execute Step 8 programmatically.
* **Required Pin:** Define Step 8 strictly as a **structural substitution evaluation** at $1/r = 0$ (or setting $r^{-1} \to 0$ via `Subst`), rather than a continuous analytical limit calculation.

---

## 2. Technical & Data Structure Ambiguities

### E. Replay Hash Concatenation Byte Format (D32 / §26)

* **The Issue:** D32 defines `CurrentStepHash = SHA-256(PreviousStepHash || CanonicalJSON(...))`.
* **Agent Failure Mode:** The agent may concatenate `PreviousStepHash` as a 64-character ASCII hex string, while another implementation concatenates it as 32 raw binary bytes, leading to hash mismatch across platforms/runs.
* **Required Pin:** Pin the input payload byte sequence:

$$\text{Payload} = \text{ASCII\_BYTES}(\text{PreviousStepHash}_{\text{hex64}}) \mathbin{\Vert} \text{UTF8\_BYTES}(\text{CanonicalJSON})$$



### F. `Normalize` Expansion and Confluence Rules (D30 / §13)

* **The Issue:** D30 specifies trig rules ($\cos^2\theta \to 1 - \sin^2\theta$, $\sin^2\theta + \cos^2\theta \to 1$), rational normalization, and like-term collection, but does not specify term expansion rules.
* **Agent Failure Mode:** Without mandatory distributive expansion ($\text{Mul}(\text{Add}(\dots), \text{Add}(\dots)) \to \text{Add}(\dots)$), like terms inside nested products will fail to combine, causing `ZeroTest` to return `UNDECIDED` on valid zeros. Conversely, unconstrained expansion of trig identities can cause infinite rewrite loops.
* **Required Pin:** Specify the exact rewrite order in `Normalize`:
1. Distribute products over sums (full polynomial expansion).
2. Simplify numerical `Rational` coefficients.
3. Apply trig reduction ($\cos^2 x \to 1 - \sin^2 x$).
4. Collect like terms over canonical sorted ordering.



### G. Struct Definition & Indexing Order for `Connection` (D26 / §8)

* **The Issue:** D26 notes that `Connection` is NOT a `Tensor` and generic tensor operations reject it. However, the exact Go struct layout and index ordering for $\Gamma^\rho_{\mu\nu}$ are omitted.
* **Agent Failure Mode:** The agent may implement $\Gamma$ as `[4][4][4]GRExpr` with ambiguous index mapping (e.g., whether slot 0 is the upper index $\rho$ or lower index $\mu$).
* **Required Pin:** Define `Connection` explicitly:
```go
type Connection struct {
    // Components[rho][mu][nu] corresponds to \Gamma^\rho_{\mu\nu}
    Components [4][4][4]GRExpr
    ChartID    string
}

```



### H. Subexpression Extraction Interface for the Bridge (§21 / D22)

* **The Issue:** §21 states that "mixed expressions project representable subexpressions/terms individually", but D22 only defines `ToCoreExpr(GRExpr)`.
* **Agent Failure Mode:** The agent will not know whether `ToCoreExpr` should automatically strip/fail on non-representable nodes or if a helper function is needed.
* **Required Pin:** Explicitly state that `ToCoreExpr` returns `UnrepresentableKernelProjectionError` if **any** node in the tree is outside the mirror subset. Add an explicit helper function `ExtractRepresentableSubexpr(e GRExpr) []GRExpr` for caller-driven term extraction.

---

## 3. Summary of Concrete Edits for v4

| Issue | Target Section | Required Plan Amendment |
| --- | --- | --- |
| **A. `Sqrt` Handling** | §21, §25, D19 | Pin that `FromCoreExpr` rejects `Sqrt`, and metric solvers must operate on squared components ($A, B$). |
| **B. `Subst` Scope** | §7, D20 | Restrict `Subst` on `UnknownFunction` to single `Symbol` replacements only. |
| **C. `Truncate` Binomials** | §23, D34 | Specify explicit rule $(1 + u)^{-1} \to 1 - u$ for negative power truncation up to $O(\epsilon^1)$. |
| **D. Step 8 Limit** | §23 | Replace "analytical limit $r \to \infty$" with "algebraic evaluation under $r^{-1} \to 0$". |
| **E. Replay Bytes** | §26, D32 | Explicitly specify ASCII Hex string concatenation for `PreviousStepHash`. |
| **F. `Normalize` Order** | §13, D30 | Define a strict 4-step canonical transformation pipeline to guarantee confluence. |
| **G. `Connection` Layout** | §8, D26 | Pin `[4][4][4]GRExpr` where index 0 is $\rho$ (upper), index 1 is $\mu$, index 2 is $\nu$. |
| **H. Bridge Projection** | §21 | Clarify `ToCoreExpr` strictness and add `ExtractRepresentableSubexpr` contract. |
I read all five new reviews in full and reconciled them against the current v2 plan and the frozen `specs_v2_3` material. The reviews are broadly convergent: **the architecture itself is sound; the remaining problem is execution-level under-specification**. The strongest review identifies four guaranteed implementation choices still left open, while the other reviews supply the concrete resolutions. GR_PLAN2_deepseek GR_PLAN2_z

The target should therefore be **v3 = same architecture, but every implementation choice pinned**.

# Consolidated verdict

I would classify the findings this way:

| Area | Net result |
|---|---|
| Kernel architecture | **GREEN** |
| Level 1 / 2 / 3 boundaries | **GREEN** |
| GR external module | **GREEN** |
| HYPOTHESIS bridge concept | **GREEN** |
| Growth Gate philosophy | **GREEN** |
| GR workload choice | **GREEN** |
| GR↔`core.Expr` representation | **RED — must pin exact contract** |
| Unknown-function representation | **RED — must pin** |
| Tensor representation/operations | **RED — must pin** |
| Schwarzschild reduction | **RED — must pin exact script** |
| Curvature conventions | **RED — must pin exact formulas** |
| Replay/hash protocol | **RED — must pin exact format** |
| Baseline/PASS0 mechanics | **YELLOW — pin exact execution protocol** |
| Silent-wrongness testing | **GREEN after minor additions** |
| Level-2/Growth outcome model | **GREEN after pending-state correction** |

The important point is that **none of these require putting GR, tensors, trig, or unknown functions into the kernel**. The current v2's core architectural decisions remain correct. GR_IMPLEMENTATION_PLAN_20261002

---

# 1. LOCK the GR mathematical representation

This is the biggest remaining issue.

The v2 says GR structures use `GRExpr`, but it still leaves enough detail unspecified that two agents could implement incompatible symbolic systems. GR_IMPLEMENTATION_PLAN_20261002

The final contract should be:

```text
phys-gr/symbolic/GRExpr

Node set exactly:

Symbol
Rational
Add
Mul
Neg
Pow(integer exponent)
Sin
Cos
UnknownFunction
```

No `Sqrt`, no arbitrary functions, no symbolic exponents, no `Derivative` node.

The **Gemini suggestion that the local power rule should accept arbitrary rational exponents should be rejected**. It is broader than the agreed bounded workload. Keep integer exponents, including negative integers.

The Z review correctly identifies that the derivative representation itself must be pinned, not just mentioned. GR_PLAN2_z

---

# 2. LOCK `UnknownFunction`

This should be one of the most precise pins.

I recommend:

```go
type UnknownFunction struct {
    Name          string
    Args          []Symbol
    DerivativeOrder []uint8
}
```

with:

```text
len(DerivativeOrder) == len(Args)

A(r):
    Name = "A"
    Args = ["r"]
    DerivativeOrder = [0]

A'(r):
    Name = "A"
    Args = ["r"]
    DerivativeOrder = [1]

A''(r):
    Name = "A"
    Args = ["r"]
    DerivativeOrder = [2]
```

Maximum derivative order for this workload:

```text
2
```

Higher order is unsupported.

This is superior to encoding `"A'"` as a function name because `A`, `A'`, and `A''` remain structurally related while still being distinct algebraic atoms for normalization.

Differentiation must increment the appropriate derivative-order slot.

The reviews correctly identify this as a guaranteed GR-5 encounter rather than a speculative issue. GR_PLAN2_z

---

# 3. LOCK the complete local differentiation table

The implementation agent must not invent this.

```text
d Rational / dx                         = 0

d Symbol(x) / dx                        = 1
d Symbol(y) / dx, y != x                = 0

d Add(u,v) / dx                         = u' + v'
d Mul(u,v,...) / dx                      = n-ary product rule
d Neg(u) / dx                            = -u'

d Pow(u,n) / dx                          = n*u^(n-1)*u'
                                         n ∈ integers

d Sin(u) / dx                            = Cos(u)*u'
d Cos(u) / dx                            = -Sin(u)*u'

d UnknownFunction(f,args,orders) / dx    =
    increment derivative order of x
    if x is an argument

d UnknownFunction(...)/dx                = 0
    when x is not an argument
```

And:

```text
maximum derivative order = 2
```

No symbolic exponent calculus.

This closes the HI-3 / U1 class of ambiguity. GR_PLAN2_deepseek

---

# 4. LOCK the GR ↔ kernel bridge completely

The bridge is conceptually correct, but the agent still needs the exact procedure.

Use:

```text
GRExpr
  ↓
ToCoreExpr()
  ↓
core.Expr
  ↓
hypothesis.NewCandidateConcept()
  ↓
temporary core.Object
  ↓
kernel operation
  ↓
FromCoreExpr()
  ↓
GRExpr
```

The bridge is **lazy**:

> A GR expression becomes a `core.Object` only at a deliberate kernel-operation call site. The object is temporary and does not persist as a GR domain object after the operation.

The bridge must never modify the underlying GR expression.

The reviewers correctly identified this as a remaining deterministic choice. GR_PLAN2_deepseek

## Exact projection boundary

`ToCoreExpr()` accepts exactly:

```text
Symbol
Rational
Add
Mul
Neg
Pow(integer exponent)
```

and recursively traverses any subexpression.

Therefore this is legal:

```text
ToCoreExpr(term)
```

even when the larger GR expression contains unrepresentable siblings.

A mixed sum may be decomposed into representable subterms for kernel reuse, but **the GR-local complete expression remains authoritative**.

Anything containing:

```text
Sin
Cos
UnknownFunction
```

fails with a named error:

```text
UnrepresentableKernelProjectionError
```

No placeholder-symbol conversion.

This reconciles the bridge-granularity criticism. GR_PLAN2_deepseek

---

# 5. LOCK bridge `Dimension`, assumptions and ID

### Dimension

Use:

```text
GR-0 … GR-8c:
    geometrized units G = c = 1
    bridge Dimension = Dimensionless

GR-8d:
    explicit SI dimensions
    for the dedicated dimensional audit
```

The loss of physical dimensional meaning in the geometrized bridge is explicitly recorded as a GR limitation/evidence item; it is not a kernel defect.

### Assumptions

Do **not** automatically inject `r > 0`, `A != 0`, etc. into every bridge.

Instead:

```text
ToCoreObject(
    expr,
    dimension,
    explicit structured core.AssumptionSet,
    explicit ConventionSet
)
```

The GR-local caller owns the assumptions.

When an assumption is representable in the kernel vocabulary, it is translated into a structured `core.Assumption`, never a raw string.

This incorporates the Qwen finding while avoiding hidden automatic semantic inference. GR_PLAN2_z

### ID

Do not use:

```text
bridge1
bridge2
timestamp
UUID
```

Use a deterministic content-derived ID:

```text
gr/bridge/<hex SHA-256 of canonical bridge payload>
```

where the payload includes at least:

```text
canonical core.Expr
Kind
Dimension
Assumptions
Conventions
```

This is better than the older `gr/<chart>/<name>` proposal because the object ID then corresponds deterministically to its content. Qwen and Z both correctly identified this as a missing determinism pin. GR_PLAN2_z GR_PLAN2_deepseek

---

# 6. LOCK `FromCoreExpr`

This was missed in v2.

Define:

```go
func ToCoreExpr(e GRExpr) (core.Expr, error)
func FromCoreExpr(e core.Expr) (GRExpr, error)
func ToCoreObject(...) (core.Object, error)
```

`FromCoreExpr` accepts exactly the same kernel-compatible subset.

Therefore:

```text
Sqrt
Call
Relation
BranchSet
```

are not imported into `GRExpr`.

That closes the round-trip ambiguity identified by Gemini. GR_PLAN2_qwen

---

# 7. LOCK the tensor representation

The v2 tensor concept is correct but still needs an executable data model. GR_IMPLEMENTATION_PLAN_20261002

I would choose a **dense row-major representation**, not a Go map:

```go
type Tensor struct {
    Rank       int
    Slots      []IndexSlot
    Components []GRExpr
    Symmetries []SymmetryRule
    ChartID    string
}
```

Rules:

```text
0 <= Rank <= 4

number of components = 4^Rank

Components ordered lexicographically / row-major

Index domain = spacetime

Variance = Covariant | Contravariant

ChartID = immutable canonical string
```

This avoids map-iteration nondeterminism.

The `SymmetryRule` should itself have a fixed canonical structure:

```text
Permutation
Sign (+1 or -1)
```

Declared symmetries are verified before being relied upon.

Required tensor operations:

```text
TensorAdd
TensorSubtract
ScalarMul
Contract
RaiseIndex
LowerIndex
ApplyMetric
ApplyInverseMetric
```

This directly closes the missing U3 tensor arithmetic identified by Z. GR_PLAN2_z

And:

> **Connection is not a Tensor.**

The coding agent must not let generic tensor operations accept Christoffel symbols as though they were tensors.

---

# 8. LOCK the chart

Define:

```go
type Chart struct {
    ID          string
    Coordinates []Symbol
    Order       []int
}
```

For this workload:

```text
ID = "spherical-static"
Coordinates = [t, r, theta, phi]
Order = [0,1,2,3]
```

Use exact ASCII symbol names:

```text
t
r
theta
phi
A
B
M
G
eps
k1
k2
```

rather than allowing the implementation agent to choose Unicode versus ASCII representations.

The rational chart:

```text
(x0, r, x=cos(theta), phi)
```

is a **secondary bridgeability census only**, not a coordinate-transformation engine.

This is how the earlier trig-versus-rational-chart disagreement should be resolved: polar is the canonical GR representation; the rational representation measures how much more of the workload can reach the kernel. GR_PLAN2_claude

---

# 9. LOCK curvature conventions mathematically

This is a genuine blocker until pinned.

Use the exact MTW convention already proposed:

\[
\Gamma^\rho_{\mu\nu}
=
\frac12 g^{\rho\sigma}
(\partial_\mu g_{\nu\sigma}
+\partial_\nu g_{\mu\sigma}
-\partial_\sigma g_{\mu\nu})
\]

\[
R^\rho{}_{\sigma\mu\nu}
=
\partial_\mu\Gamma^\rho_{\nu\sigma}
-\partial_\nu\Gamma^\rho_{\mu\sigma}
+\Gamma^\rho_{\mu\lambda}\Gamma^\lambda_{\nu\sigma}
-\Gamma^\rho_{\nu\lambda}\Gamma^\lambda_{\mu\sigma}
\]

\[
R_{\mu\nu}=R^\rho{}_{\mu\rho\nu}
\]

\[
R=g^{\mu\nu}R_{\mu\nu}
\]

\[
G_{\mu\nu}
=
R_{\mu\nu}-\frac12Rg_{\mu\nu}
\]

and workload conventions:

```text
signature = -+++
Λ = 0
coordinate order = t,r,theta,phi
```

The review is correct that merely saying "defined in phys-gr" leaves the coding agent a mathematical choice. GR_PLAN2_deepseek

---

# 10. LOCK the GR-local normalization and ZeroTest semantics

`Normalize` must be bounded to:

```text
common denominator
rational factor collection
like-term collection
negative integer power normalization
Cos² → 1 − Sin²
exact rational simplification
```

No general CAS.

`ZeroTest` must be:

```text
ZERO
NONZERO
UNDECIDED
```

with:

```text
ZERO:
    exact normalization proves zero

NONZERO:
    only when the bounded algebra has a direct structural proof

UNDECIDED:
    otherwise
```

Never guess.

The independent point oracle is test-only. It can falsify a false `ZERO`; it cannot turn a sampled result into a proof.

This resolves the B3/G9 issue. GR_PLAN2_claude

---

# 11. LOCK the Schwarzschild derivation script

This should not remain "bounded solver."

The review is right: that still gives an implementation agent too much freedom. GR_PLAN2_z

Use this exact script:

```text
1. Start with:

   ds² = -A(r)dt² + B(r)dr² + r²dΩ²

2. Compute Γ.

3. Compute R_tt, R_rr, R_theta theta.

4. Certify:

   R_tt/A + R_rr/B = (AB)' / (r A B²)

5. Vacuum condition implies:

   (AB)' = 0

6. Therefore:

   AB = k1

7. Impose asymptotic flatness:

   A → 1
   B → 1
   as r → ∞

   Therefore:

   k1 = 1

8. Substitute:

   B = 1/A

9. Certify:

   R_theta theta = 1 - A - rA'

10. Vacuum gives:

    (rA)' = 1

11. Candidate-and-certify:

    rA = r + k2

12. Therefore:

    A = 1 + k2/r
    B = 1/A

13. Independently verify all vacuum components.

14. In SI correspondence:

    k2 = -2GM/c²
```

The machine **certifies candidate steps**; it does not discover the solution. That distinction must be explicit. This reconciles the verify-only/derive-first conflict. GR_PLAN2_claude

---

# 12. LOCK geometric/SI mass semantics

I would make one small improvement over the reviewers' symbol recommendation.

Do not let the same symbol implicitly change physical dimension.

Use:

```text
mu
    geometric Schwarzschild length parameter

M
    SI mass

mu = G*M/c²
k2 = -2*mu
```

Then:

```text
A = 1 - 2*mu/r
```

in the geometric GR representation.

This makes:

```text
mu/r
```

dimensionless even before entering the geometrized "Dimensionless" bridge.

The SI audit can then explicitly verify:

```text
G*M/c² → L
G*M/(c²*r) → 1
G*M/r² → L/T²
```

This is cleaner than making `M` silently mean a mass in one stage and a length parameter in another. The underlying reviewers correctly identified the dimensional ambiguity; the cleanest resolution is to separate the symbols. GR_PLAN2_deepseek

---

# 13. LOCK GR-8 Newtonian reduction

Do **not** use `phys.Limit`.

The current v2 is already correct about that. GR_IMPLEMENTATION_PLAN_20261002

Pin the actual procedure:

```text
ε = GM/(c²r)

retain ε^0 and ε^1 only

g_tt = -(1 + 2Φ/c²) + O(ε²)

Φ = -GM/r

a_r = -GM/r²
```

Use a bounded local operation:

```text
Truncate(expr, eps, 1)
```

not a general series engine.

Then the correspondence may be recorded through:

```text
Compare(...)
Session.Identify(...)
```

with HYPOTHESIS operands/output where appropriate.

This gives session machinery one deliberate GR use instead of leaving the entire session layer unused.

---

# 14. LOCK GR-local replay

The current v2 says it has a replay trace but still leaves critical details open. GR_PLAN2_deepseek

Mirror the kernel as closely as possible:

```text
GRStep fields:

StepID
Index
OperationID
InputCanonicals
InputHashes
ParamsCanonical
OutputCanonical
OutputHash
KernelBridgeHash
PreviousStepHash
CurrentStepHash
```

Use:

```text
Genesis PreviousStepHash =
0000000000000000000000000000000000000000000000000000000000000000
```

and:

```text
StepID:
gr-step-000001
gr-step-000002
...
```

`CurrentStepHash` is:

```text
SHA-256(
    PreviousStepHash ||
    CanonicalJSON(step-without-CurrentStepHash)
)
```

Canonical JSON must have fixed field order and deterministic array ordering.

Require:

```text
Decode(Canonical(x)) == x
```

before accepting the artifact.

This is preferable to the alternative custom genesis string because it mirrors the existing kernel integrity model.

---

# 15. GR-3a must call `ops.Differentiate` directly

This is a subtle but important pin.

Do:

```text
phys-gr
   ↓
hypothesis bridge
   ↓
ops.Differentiate(...)
```

Do **not** use `Session.Step` for the intentional failed probe.

Why?

Because the experiment is measuring the pure kernel operation boundary, and the frozen spec does not clearly define the session behavior of a failed step. The review correctly identified this as a remaining ambiguity. GR_PLAN2_deepseek

The GR trace records the probe outcome separately.

---

# 16. Local operations must be clearly distinct from kernel operations

Use:

```text
Diff
Normalize
ZeroTest
Subst
Truncate
Reduce
```

Never:

```text
Differentiate
Simplify
Limit
Solve
Substitute
```

for GR-local functions.

This prevents agents from creating a pseudo-shadow API that could later be mistaken for an extension of `phys`.

The current v2 already points in this direction. GR_IMPLEMENTATION_PLAN_20261002

---

# 17. The Level-2 outcome stays `PENDING`

The earlier contradiction is correctly resolved.

GR alone cannot establish a second established consumer.

Therefore:

```text
NO-GROWTH

LEVEL-2-CANDIDATE-PENDING

KERNEL-GROWTH-CANDIDATE-PENDING

ESCALATE-TO-SPEC
```

remain the terminal states.

There is **no automatic Level-2 promotion** during this plan.

The reviews correctly converged on this. GR_IMPLEMENTATION_PLAN_20261002

---

# 18. But `KERNEL-GROWTH-CANDIDATE` needs a provisional Level-2 path

This is an important subtlety from the consolidated review.

The plan cannot simultaneously say:

```text
never pre-build Level 2
```

and:

```text
kernel candidate requires Level-2 failure
```

The solution is:

```text
GR execution
    ↓
potential generic trusted-invariant failure
    ↓
KERNEL-GROWTH-CANDIDATE-PENDING
    ↓
human-authorized follow-on Level-2 attempt
    ↓
failure/success
    ↓
only then Growth Gate advancement
```

So the current GR run does **not** build a speculative Level-2 library in advance.

The Level-2 attempt is created only if evidence triggers it.

That resolves G4 without weakening the no-prebuilding rule. GR_PLAN2_claude

---

# 19. PASS0 is a separate prerequisite task

Given our earlier decision, make this explicit:

> **PASS0 is executed by a separate conformance task before the GR implementation agent begins.**

The GR implementation agent does not decide whether PASS0 is green.

It verifies the baseline record exists and matches the recorded `phys` checkout.

PASS0 includes:

```text
Pin A:
Identify(RestMass, E/c²)
    → KindRelation
    → M
    → IDENTIFIED
    → NONE
    → operands unchanged

Pin B:
Differentiate(Pow(x,-1))
    → UnsupportedOperationError

Pin C:
AGENTS.md caveat
```

The exact Pin B test target should be pinned to the existing file:

```text
ops/differentiate_test.go
```

provided it is confirmed as a member of the 42-file set; otherwise halt rather than create a file.

The PASS0 target membership check is mandatory.

---

# 20. The 42-file baseline needs one final explicit inventory pin

The operational freeze is:

```text
39 implementation files
+
AGENTS.md
mechanics/README.md
relativity/README.md
=
42
```

The earlier re-freeze review confirms the three documentation files and exact closed-world test. adv_review11(2)

The GR plan should simply list these three paths explicitly instead of requiring the coding agent to discover them from `adv_review11`.

---

# 21. Mutation testing: separate test-sanity from code mutations

The current v2 has already improved this, but the final plan should distinguish:

### Test-sanity check

```text
Restore broken Pin A fixture
Identify(Energy, E/c²)
```

This is **not** production mutation testing.

### Production mutations

Use disposable worktree only:

```text
Identify returns operand Kind
remove negative-power guard
weaken Substitute kind equality
remove SelectBranch constraint
remove assumption conflict
ZeroTest always ZERO
Contract skips variance check
```

Then destroy the mutant tree and re-verify authoritative `phys`.

This matches the review consensus. GR_PLAN2_claude

---

# 22. Silent-wrongness suite should be expanded slightly

Keep the corrected semantics:

```text
bare Symbol("f") / t → 0
    ACCEPTED-CORRECT
```

Do not call that kernel wrong.

Then add:

```text
SW-1 Symbol/entity ambiguity
SW-2 UnknownFunction dependency representation
SW-3 chart mixing
SW-4 treating Γ as tensor
SW-5 invalid index contraction
SW-6 invalid symmetry assumption
SW-7 dimension-valid semantic mismatch
SW-8 type-decay addition
SW-9 assumption inconsistency
SW-10 convention mismatch
SW-11 ZeroTest false-zero
SW-12 bridge canonicalization inconsistency
SW-13 tensor canonicalization/replay inconsistency
```

The exact numbering can be whatever v3 chooses, but **the plan must define the complete set**, rather than vaguely saying "etc.".

---

# 23. Correct the "no metrics" wording

This is a genuine factual wording issue.

The kernel already has nominal kinds such as:

```text
Spacetime
MinkowskiMetric
FourMomentum
```

The correct statement is:

> **The kernel contains no General Relativity metric semantics or tensor semantics. Existing nominal SR-related Kinds remain part of the frozen MVP.**

That is more precise than:

```text
kernel:
    no metrics
```

The review correctly identified this. GR_PLAN2_claude

---

# 24. Lock the role of the existing SR corpus

GR should **not** modify `relativity`.

It may use SR as an external established reference where relevant.

The existing `relativity` package remains:

```text
Special Relativity
```

only.

No GR semantics get smuggled into it.

This preserves the already-established theory boundary.

---

# 25. Define the human review process

The evidence record now has fields 19–20, but we need a deterministic workflow.

Use:

```text
Implementation agent
      ↓
classifies obstacle
      ↓
independent reviewer
      ↓
records second classification
      ↓
disagreement resolved
      ↓
human curator decides promotion
```

For the independent reviewer:

> A second agent/reviewer who did not author the relevant implementation pass.

The human curator remains the final authority for:

```text
LEVEL-2 promotion
KERNEL-GROWTH-CANDIDATE
ESCALATE-TO-SPEC
```

This satisfies the "independent classification" requirement without requiring a second human for every ordinary obstacle. GR_PLAN2_claude

---

# 26. Add the stop rule

The capability inventory must be closed.

If the agent reaches something outside it:

```text
BLOCKED
```

not:

```text
invent new abstraction
```

Likewise each pass gets a human-selected:

```text
T_pass
```

Wall-clock overrun is an execution observation, never a kernel-growth argument.

This was explicitly identified as missing. GR_PLAN2_claude

---

# 27. Keep the Kernel Contact Ledger, but do not turn it into a score

This remains correct:

```text
pass
kernel_operation
inputs
outcome
local_operation_invoked_instead
boundary_note
phys_baseline_id
```

Plus observational measurements:

```text
node count
term count
wall-clock time
```

No:

```text
kernel usage percentage
minimum efficacy ratio
automatic growth threshold
```

The rejection of Gemini's "Substrate Efficacy Ratio" should remain permanent. GR_PLAN2_claude GR_PLAN_qwen

---

# 28. One documentation precedence rule

There is one old planning-document inconsistency noted by Z around the `OperationParams` example. The normative authority remains:

```text
specs_v2_3.md
```

not an inconsistent example in an older planning document.

Therefore add:

> **Where any historical planning document conflicts with `specs_v2_3.md`, the frozen specification is authoritative. Historical planning examples are non-normative.**

That prevents an implementation agent from choosing the wrong parameter mapping without requiring us to alter the frozen specification.

---

# What I would explicitly reject from the reviews

These recommendations should **not** enter v3:

```text
❌ kernel tensor types
❌ kernel Sin/Cos
❌ kernel UnknownFunction
❌ extending phys.Differentiate for GR
❌ expanding phys.Limit for GR
❌ general ODE solver
❌ general CAS
❌ placeholder sin_theta
❌ arbitrary rational-exponent local calculus
❌ Substrate Efficacy Ratio
❌ efficacy threshold
❌ automatic Level-2 promotion
❌ treating bare-symbol differentiation as a kernel defect
❌ putting GR inside phys
```

The current v2 already correctly rejects most of these. GR_IMPLEMENTATION_PLAN_20261002

---

# The actual v3 target

After all of the above, the execution architecture becomes:

```text
                FROZEN phys
                    │
             public API only
                    │
                    ▼
             github.com/PithomLabs/phys-gr
                    │
          ┌─────────┴─────────┐
          │                   │
     GRExpr algebra       GR structures
          │                   │
          ├── Diff            ├── Chart
          ├── Normalize       ├── Tensor
          ├── ZeroTest        ├── Metric
          ├── Subst           ├── Connection
          └── Truncate        └── Curvature
          │
          ▼
     exact kernel subset
          │
          ▼
  temporary HYPOTHESIS Object
          │
          ▼
      frozen phys ops
          │
          ▼
      reintegrate GRExpr
```

with:

```text
GR-3a
    frozen kernel probe

GR-3b
    independent Level-1 Diff

GR-7/8
    certified bounded Schwarzschild reduction

GR-8
    explicit SI/Newtonian correspondence

all along:
    GR replay trace
    Kernel Contact Ledger
    adversarial suite
    Growth Evidence
```

The five reviews therefore do **not** justify another architectural redesign. They establish a finite final pin set. The most important remaining pins are exactly the ones the strongest line-by-line review identified: derivative representation, reduction protocol, tensor arithmetic, bridge ID/assumptions, plus the associated deterministic replay, chart, convention, dimension, and baseline mechanics. GR_PLAN2_z GR_PLAN2_deepseek

**My recommendation: do one final v3 editorial pass that incorporates these pins, then run the adversarial reviews once more. Do not start PASS0 or GR implementation until that review says the coding agent has zero remaining architectural choices.**
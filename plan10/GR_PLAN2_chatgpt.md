I read the revised v2 carefully against the six reviews and the earlier locked decisions. **It is substantially better, but I would not authorize the coding-agent handoff yet.** The Growth Gate architecture is now sound; the remaining problems are mostly places where the document still leaves an implementation decision implicit.

The two most important unresolved issues are real architectural/execution blockers.

## What is now correctly locked

The revision successfully incorporates the major conclusions:

- GR is external Level 1, while `mechanics` and `relativity` are grandfathered. GR_IMPLEMENTATION_PLAN_20261002
- The `core.Object` bridge is explicitly HYPOTHESIS-only. GR_IMPLEMENTATION_PLAN_20261002
- GR-3a must exercise frozen `phys.Differentiate` before GR-3b continues locally. GR_IMPLEMENTATION_PLAN_20261002
- Trig and unknown functions are genuinely Level-1 constructs rather than fake placeholder symbols. GR_IMPLEMENTATION_PLAN_20261002
- Schwarzschild is derivation-first, not merely verification. GR_IMPLEMENTATION_PLAN_20261002
- Newtonian reduction is explicitly removed from `phys.Limit`. GR_IMPLEMENTATION_PLAN_20261002
- Level-2 promotion is now `LEVEL-2-CANDIDATE-PENDING`, so GR alone cannot fabricate a second consumer. GR_IMPLEMENTATION_PLAN_20261002
- Silent-wrongness is correctly separated from user misuse. GR_IMPLEMENTATION_PLAN_20261002
- The substrate-efficacy ratio has correctly been rejected. GR_IMPLEMENTATION_PLAN_20261002

So the architecture itself is in good shape.

# Remaining blockers

## 1. The GR expression representation is still internally inconsistent

This is the most important remaining issue.

Section 6 says:

> "Everything with index/variance/chart/frame semantics is a GR-local Go struct holding `core.Expr` components"

and the kernel bridge is built around that model. GR_IMPLEMENTATION_PLAN_20261002

But the same section says that `Sin`, `Cos`, and `UnknownFunction` are Level-1 nodes that **cannot** be converted to `core.Expr`. GR_IMPLEMENTATION_PLAN_20261002

Therefore this cannot work as written:

```text
Metric component
    ↓
core.Expr
    ↓
gφφ = r² Sin(θ)²
```

because `Sin(θ)` is not a `core.Expr`.

The plan needs one definitive rule:

> **All GR-local mathematical structures — scalar expressions, metric components, tensor components, connection coefficients, curvature components — use `phys-gr`'s `GRExpr` type. `core.Expr` is only a bridge representation for the subset of `GRExpr` that is exactly representable by the frozen kernel.**

So the model becomes:

```text
GRExpr
   │
   ├── local-only
   │     Sin
   │     Cos
   │     UnknownFunction
   │
   └── kernel-representable subset
              ↓
         ToCoreExpr()
              ↓
         core.Object
```

This should also add the inverse operation:

```text
FromCoreExpr(core.Expr) -> GRExpr
```

because the plan explicitly says kernel `Simplify` results are to be "reintegrate[d] into `GRExpr`," but no reverse conversion is defined. GR_IMPLEMENTATION_PLAN_20261002

**I would not proceed until this is pinned.**

---

## 2. The Level-1 algebra is named, but not actually specified enough

The plan says `phys-gr/symbolic/expr.go` contains:

```text
Expr, Symbol, Rational, Add, Mul, Neg, Pow,
Sin, Cos, UnknownFunction
```

and gives some responsibilities, but it does not define the actual representation contract. GR_IMPLEMENTATION_PLAN_20261002

For a coding agent, at minimum the plan needs to pin:

```text
node fields
node tags / ordinals
canonical child ordering
equality semantics
immutable representation
canonical serialization
supported differentiation rules
supported normalization rules
supported substitution rules
unsupported cases
```

For example, the differentiator must explicitly define:

```text
d(Symbol(x))/dx
d(Symbol(x))/dy
d(Rational)/dx
d(Add)/dx
d(Mul)/dx
d(Neg)/dx
d(Pow(base,n))/dx
d(Sin(x))/dx
d(Cos(x))/dx
d(UnknownFunction(f,x))/dx
```

including which exponents and function forms are supported.

Likewise `Normalize`/`ZeroTest` needs an exact rule set rather than just a responsibility list. The current plan identifies common-denominator and like-term work, but not the complete bounded rewrite contract. GR_IMPLEMENTATION_PLAN_20261002

This doesn't require building a general CAS. It requires making the **bounded GR-local algebra deterministic enough that two coding agents would implement the same thing.**

---

## 3. The tensor engine is still underspecified

The plan improved this substantially, but it still stops short of an executable data contract.

It currently says:

```text
Tensor:
    Rank / IndexSlots[] / Components / Symmetries / ChartID

IndexSlot:
    slot / variance / index domain

Component keys:
    fixed tuples, not arbitrary strings
```

with operations such as `Contract`, `RaiseIndex`, `LowerIndex`, etc. GR_IMPLEMENTATION_PLAN_20261002

But the coding agent still has to choose:

- exact Go representation of an index slot;
- exact component-key representation;
- coordinate index domain;
- how absent components are represented;
- how zero components are represented;
- exact symmetry representation;
- whether symmetry is metadata or mechanically enforced;
- ordering of components;
- how contraction enumerates indices;
- how chart identity is represented.

Pin these.

For example, explicitly choose:

```text
Coordinate indices = {0,1,2,3}
ChartID = canonical string/identifier
Variance = Covariant | Contravariant
ComponentKey = canonical immutable four-index tuple representation
Tensor components = map[ComponentKey]GRExpr
```

or whatever exact representation you prefer.

The point is not that this particular representation is mandatory; the point is that the coding agent must not choose it.

---

## 4. The curvature convention is still not actually pinned

The plan says:

> "GR curvature convention explicitly defined in `phys-gr`"

but that leaves the convention to the coding agent. GR_IMPLEMENTATION_PLAN_20261002

That is too important to leave open.

Pin the exact definitions, e.g. the chosen sign/order for:

\[
R^\rho{}_{\sigma\mu\nu}
\]

then:

\[
R_{\sigma\nu}
\]

then:

\[
R=g^{\mu\nu}R_{\mu\nu}
\]

then:

\[
G_{\mu\nu}=R_{\mu\nu}-\frac12 Rg_{\mu\nu}.
\]

Also pin whether the vacuum equation is represented as:

```text
Gμν = 0
```

or:

```text
Rμν = 0
```

or both, and which one is the canonical workload relation.

The coding agent should not get to choose a curvature convention because that choice changes every downstream golden value.

---

## 5. Schwarzschild derivation still lacks the operational derivation contract

The plan correctly says:

> start with `A(r), B(r)` and derive through Christoffel → Riemann → Ricci → vacuum equations → bounded reduction. GR_IMPLEMENTATION_PLAN_20261002

But it still doesn't tell the agent exactly **which reduction** is intended.

You need to pin:

```text
coordinate gauge:
    areal radius r

ansatz:
    ds² = -A(r)dt² + B(r)dr² + r²dΩ²

domain:
    exterior region

boundary conditions:
    asymptotic flatness
    A(r) → 1
    B(r) → 1

time normalization:
    fixed at infinity

integration constant:
    exact name and relation to physical mass
```

I would also explicitly distinguish the geometric mass parameter from SI mass.

For example:

```text
μ = Schwarzschild length parameter
μ = 2GM/c²
```

rather than letting the same `M` silently change dimensional meaning between geometric and SI phases.

Then specify the bounded solution sequence rather than telling the agent to "solve the equations." There is no general ODE solver by design. GR_IMPLEMENTATION_PLAN_20261002

---

## 6. Newtonian reduction is still not concrete enough

The plan correctly rejects `phys.Limit` and says GR-8 uses bounded first-order weak-field reduction with:

\[
\epsilon=\frac{GM}{c^2r},
\qquad
\Phi=-\frac{GM}{r}.
\]

GR_IMPLEMENTATION_PLAN_20261002

But the actual reduction procedure is still unspecified.

Pin:

```text
weak-field parameter
    ε = GM/(c²r)

retained order
    O(ε)

discarded order
    ε² and higher

metric weak-field extraction
    exact formula

potential identification
    exact formula

Newtonian acceleration derivation
    exact formula, if included
```

Also specify exactly how the geometric-unit solution is restored to SI. Otherwise `G`, `c`, and the mass parameter can become ambiguous at the boundary between GR-8 stages.

---

## 7. The per-pass execution detail is still missing

This is the clearest "it's still a planning document" issue.

The plan says:

> "Each pass records: objective, files/packages, dependencies, capabilities, worked examples, tests, expected failures, evidence, level implicated, exit criteria"

but then GR-0 through GR-8 are still essentially one-line labels. GR_IMPLEMENTATION_PLAN_20261002

That promised information needs to actually be present.

For each pass, provide a table:

```text
GR-1
Objective:
Inputs:
Outputs:
Data structures:
Kernel operations exercised:
Local operations exercised:
Worked example:
Golden result:
Negative tests:
Expected failures:
Evidence artifacts:
Growth classification:
Exit criteria:
```

Do that for all GR-0 through GR-8.

Otherwise the implementation agent will still make design decisions during execution.

---

## 8. GR replay/hash protocol needs one more level of precision

The plan establishes a separate GR replay trace:

```text
GRStep:
    StepID / OperationID / InputCanonical /
    ParamsCanonical / OutputCanonical / CurrentHash
```

and separate SHA-256 hashing. GR_IMPLEMENTATION_PLAN_20261002

Appendix C says GRExpr, Tensor, Metric, Connection, Curvature, and GRStep receive canonical JSON + SHA-256. GR_IMPLEMENTATION_PLAN_20261002

But the actual canonical protocol isn't pinned.

Before implementation, specify:

```text
GR schema/version
node/type tags
field ordering
array ordering
number encoding
map ordering
hash input
genesis hash
chain formula
replay procedure
```

Otherwise two implementations could produce mathematically equivalent but byte-different GR artifacts.

---

# Governance/editorial items that still need cleanup

These are smaller, but they matter given your stated goal of **zero ambiguity**.

### Invalid cross-reference

Section 3 refers to:

> "post-PASS0 baseline record (§3.8)"

but Section 3 has no §3.8. GR_IMPLEMENTATION_PLAN_20261002

Create the actual subsection or correct the reference.

### Ambiguous level-direction statement

Section 2 still says:

> "No level modifies the level below it except through the Growth Gate"

which is ambiguous because Level 1/2/3 direction is itself hierarchical. GR_IMPLEMENTATION_PLAN_20261002

Use:

> **No level modifies another level except through the explicitly defined promotion/Growth Gate process.**

### Old "false reading" wording survived

Section 18 still says:

> `"kernel passed GR" (false — userland did work)`

while the plan simultaneously defines `NO-GROWTH` as a successful outcome. GR_IMPLEMENTATION_PLAN_20261002

Use the earlier corrected distinction:

```text
"Kernel passed GR" = acceptable only in the sense that the
kernel boundary was sufficient; it does NOT mean the kernel
performed the GR mathematics.
```

### Allowed-operation table needs an exercise status

The table lists `Compare`, `Solve`, `SelectBranch`, `Session.Identify`, etc., but the workload doesn't clearly say when each is actually exercised. GR_IMPLEMENTATION_PLAN_20261002

For each entry say:

```text
EXERCISED in GR-x
or
AVAILABLE BUT NOT EXERCISED
```

Otherwise the coding agent may invent unnecessary GR uses for them.

### PASS1 is still vague

`PASS1 — GR Plan Verification (this document acknowledged)` is not really a defined verification gate. GR_IMPLEMENTATION_PLAN_20261002

Either define exactly what PASS1 verifies or remove it and make GR-0 the first execution pass after PASS0.

### Human-review process needs an actual procedure

The evidence format now defines fields 19–20, which is good. GR_IMPLEMENTATION_PLAN_20261002

But it does not yet specify the actual artifact/process for:

```text
independent classification
disagreement resolution
human approval
```

Define where these records live and what constitutes a completed review. Otherwise field 19/20 remains a form rather than a gate.

---

# One thing I would NOT change

I would **not** adopt any of the suggestions that push us toward:

```text
kernel tensors
kernel trig
kernel UnknownFunction
general ODE solver
general CAS
substrate-efficacy ratio
automatic Level-2 promotion
```

The current v2 correctly rejects those directions. GR_IMPLEMENTATION_PLAN_20261002

Likewise, the `HYPOTHESIS core.Object` bridge is conceptually good. It gives external GR code access to the existing MRC/assumption/provenance machinery without pretending GR objects are native kernel physics types. GR_IMPLEMENTATION_PLAN_20261002

# Final disposition

I would classify the revised plan as:

**Architecture: GREEN**

**Growth-Gate governance: GREEN**

**Kernel boundary: GREEN**

**Execution contract: YELLOW — not yet handoff-ready**

The remaining work is principally to turn the declared Level-1 abstractions into **exact contracts** rather than new architectural ideas.

The most important final lock is:

```text
phys
    ↓
public core.Expr subset only

phys-gr
    ↓
one coherent GRExpr algebra
    ↓
Tensor / Metric / Connection / Curvature all use GRExpr
    ↓
ToCoreExpr / FromCoreExpr only for exact kernel subset
    ↓
temporary HYPOTHESIS core.Object
    ↓
frozen phys operations
```

Then pin the exact tensor representation, curvature convention, Schwarzschild reduction, Newtonian reduction, GR canonicalization/replay, and per-pass golden workloads.

Once those are added, **I would expect the coding-agent handoff to satisfy your actual target: no architectural blockers, no unresolved ownership questions, and no place where the implementation agent has to invent the architecture.**
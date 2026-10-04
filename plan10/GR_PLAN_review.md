I read all six reviews end-to-end and reconciled them against the current `GR_IMPLEMENTATION_PLAN.md` and the frozen `specs_v2_3` API details. The reviewers are converging on the same conclusion: the architecture is sound, but the plan still leaves a few **actual implementation choices** implicit. Those must be pinned now so the coding agent cannot invent architecture while implementing.

The most important unresolved point is the GR ↔ `core.Object` boundary. The reviews independently identified that §6 says GR structures are not `core.Object`, while §7 simultaneously expects GR to use object-based `ops`; because every `ops` function operates on `core.Object`, the current document is literally unexecutable at that point. GR_PLAN_claude GR_PLAN_z

# Consolidated architectural lock

## 1. GR remains Level 1 and remains outside `phys`

This survives every review.

```text
github.com/PithomLabs/phys
    frozen Level 3 substrate

github.com/PithomLabs/phys-gr
    Level 1 General Relativity
    one-way dependency on phys
```

No GR code enters the 42-file `phys` tree. New Level-1 packages follow this rule; the existing `mechanics/` and `relativity/` packages are grandfathered because they predate this architecture.

The separate-module boundary is not merely organizational: it ensures Go's `internal/` boundary prevents `phys-gr` from reaching `internal/kernel`. GR_PLAN_deepseek

---

# 2. The correct GR/core.Object bridge

This is the most important new lock.

The reviews are right that saying:

> "GR objects are not `core.Object`"

and simultaneously saying:

> "GR uses `ops.Add`, `ops.Multiply`, `ops.Simplify`, ..."

cannot both literally be true.

The resolution is **not** to make GR physical objects into kernel objects.

Instead:

> **GR mathematical/physical structures remain Level-1 objects. A GR scalar expression is bridged into a temporary generic `core.Object` only when a frozen kernel operation is deliberately being exercised.**

The exact model should be:

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

For a kernel-compatible scalar:

```text
GRExpr
  ↓ ToCoreExpr()
core.Expr
  ↓ hypothesis.NewCandidateConcept(...)
core.Object
  Kind = Expression
  Provenance = HYPOTHESIS
  CorpusStatus = NONE
```

This is explicitly permitted by the frozen hypothesis constructor, which accepts an arbitrary `Expr`, explicit `Kind` and `Dimension`, while forcibly assigning `HYPOTHESIS` and `NONE`. specs_v2_3

This is actually a **better stress test** than leaving GR at bare `core.Expr`: it lets the GR package exercise dimensional checks, MRC, assumptions, provenance contamination, `Simplify`, `Compare`, etc., while preserving the trust boundary.

### Critical restriction

A GR-local expression that contains a Level-1-only construct such as `Sin`, `Cos`, or `UnknownFunction` **cannot** be converted to `core.Expr`.

That is deliberate.

There is no fake conversion through placeholder symbols.

---

# 3. `phys-gr` gets one coherent bounded symbolic algebra

This resolves the trig/unknown-function contradiction.

The coding agent must **not** create:

```text
core.Expr + ad-hoc trig fragments
```

as a half-and-half representation.

Instead, define one Level-1 GR-local algebra, for example:

```text
phys-gr/symbolic/

    Expr
    Symbol
    Rational
    Add
    Mul
    Neg
    Pow
    Sin
    Cos
    UnknownFunction
```

with structural conversion:

```text
GRExpr → core.Expr
```

only when the entire expression is expressible in the frozen kernel language.

This is **not a general CAS**. It is a bounded symbolic layer sufficient for the GR workload.

The local algebra owns:

- unknown functions such as `A(r)` and `B(r)`;
- `sin(θ)` and `cos(θ)`;
- rational-function normalization;
- bounded differentiation;
- substitution;
- bounded first-order weak-field reduction;
- zero testing.

The kernel continues to own its existing closed AST. The frozen core expression API already provides the read-only constructors and inspectors needed for this external layer. specs_v2_3

---

# 4. Use real `Sin`/`Cos`, not `sin_theta` placeholder symbols

This is where I would **reject the reviewers' suggested placeholder-symbol workaround**.

The reviews correctly identify the danger:

```text
sin_theta → Symbol
Differentiate(..., theta)
    → 0
```

But the solution should not be a secret derivative lookup table over arbitrary string symbols.

That recreates exactly the semantic ambiguity we are trying to measure.

Instead:

```text
Sin(theta)
Cos(theta)
```

are genuine **GR-local nodes**.

The local differentiator implements:

```text
d Sin(x) / dx = Cos(x)
d Cos(x) / dx = -Sin(x)
```

This resolves the angular metric honestly:

```text
gφφ = r² Sin(θ)²
```

without pretending the kernel understands trigonometric functions.

The standard Schwarzschild spherical-coordinate representation does require this angular structure, so the trig requirement should remain in the workload rather than be evaded.

---

# 5. Unknown functions are also genuine GR-local nodes

The Qwen review raised an important point that the other reviewers did not sufficiently develop.

For a derivation rather than mere verification, use:

```text
A(r)
B(r)
```

as unknown functions.

Therefore:

```text
UnknownFunction("A", r)
UnknownFunction("B", r)
```

must be representable in `phys-gr`.

The local differentiator must implement:

```text
d A(r) / dr = A'(r)
d B(r) / dr = B'(r)
```

and higher derivatives as required by the curvature calculation.

Do **not** represent `A(r)` as:

```text
core.Symbol("A")
```

because the kernel correctly treats that as a symbol independent of `r`.

This is important: the dedicated silent-wrongness test should first establish that behavior as **ACCEPTED-CORRECT**, not incorrectly blame the kernel. A bare symbol really is syntactically constant according to the frozen semantics. The GR-local `UnknownFunction` abstraction is what supplies the missing dependency semantics. GR_PLAN_deepseek

---

# 6. GR-3a / GR-3b is now unambiguous

The final rule should be:

```text
GR-3a
    Create kernel-representable:
        Pow(r, -1)
    through a HYPOTHESIS core.Object.

    Call:
        phys.Differentiate(...)

    Expected:
        UnsupportedOperationError

    Record kernel boundary.

GR-3b
    Run phys-gr/symbolic differentiator.

    This is the Level-1 continuation.

    It does not modify, replace, shadow, or re-export
    phys.Differentiate.
```

The phrase **"does not bypass `phys.Differentiate`" should be removed**.

It is contradictory because GR-3b necessarily takes over the mathematical task after GR-3a fails. The correct wording is:

> **GR-3b is an independent Level-1 operation, not an extension or replacement of the kernel operation.**

This directly resolves the reviews' most concrete wording criticism. GR_PLAN_deepseek

---

# 7. GR-local operations do NOT become `phys.Session` operations

This is an important architectural decision.

The kernel session's transformation vocabulary is deliberately closed at 12 operations. The reviews correctly note that a GR-local differentiator cannot simply enter `Session.Step` as a new operation ID. GR_PLAN_claude2

Do **not** solve this by expanding `ops`.

Instead:

```text
phys.Session
    records kernel-level operations and trusted derivation events

phys-gr derivation trace
    records Level-1 mathematical operations
```

The GR trace must itself be deterministic and replayable:

```text
GRStep:
    StepID
    OperationID
    InputCanonical
    ParamsCanonical
    OutputCanonical
    CurrentHash
```

with SHA-256 over canonical JSON.

Kernel-compatible scalar bridges additionally record the corresponding `core.Object` hash.

This gives us two separate replay domains:

```text
Level 3:
    phys Session replay

Level 1:
    phys-gr mathematical replay
```

If the fact that external mathematical operations cannot enter the trusted session eventually proves to be a **generic, theory-independent invariant problem**, that becomes genuine Growth Gate evidence.

Otherwise, it remains a deliberate boundary.

---

# 8. Rational normalization and zero testing are mandatory Level-1 machinery

This was correctly identified by both Claude reviews.

The local GR algebra must provide bounded:

```text
Normalize
ZeroTest
```

sufficient for the Schwarzschild workload.

At minimum:

```text
common-denominator normalization
rational-factor collection
like-term collection
constant factor normalization
negative-power normalization
zero detection
```

and the necessary bounded trig identities:

```text
Sin(x)^2 + Cos(x)^2 → 1
```

The implementation must **not** become a general symbolic algebra system.

The purpose is to make:

```text
Rμν = 0
```

testable without pretending `phys.Simplify` is a GR CAS.

This was identified as the dominant practical cost of the curvature stages. GR_PLAN_claude

---

# 9. Kernel `Simplify` should still be reused where applicable

The reviewers correctly point out that moving all algebra to GR-local code would waste a major part of the stress test.

So the rule should be:

> **Use the kernel whenever the relevant subexpression is genuinely representable in `core.Expr`; use GR-local rules only for capabilities absent from the kernel.**

For example:

```text
GR-local differentiator
        ↓
produces kernel-compatible term
        ↓
bridge → HYPOTHESIS core.Object
        ↓
phys.Simplify
        ↓
reintegrate result into GRExpr
```

Thus we measure actual reuse of the kernel rather than either forcing everything through it or abandoning it.

---

# 10. Do not create a "Substrate Efficacy Ratio"

I **reject** Gemini's proposed ratio:

\[
\frac{\text{kernel operations}}
{\text{kernel operations}+\text{userland operations}}
\]

with a threshold determining kernel growth. GR_PLAN_qwen

That would be a serious architectural mistake.

Why?

Because the architecture explicitly says:

```text
GR tensor algebra → userland
GR differentiation → userland
GR trigonometry → userland
GR curvature → userland
```

Counting those operations as "kernel bypasses" would artificially penalize the architecture we intentionally designed.

Instead record a **Kernel Contact Ledger** descriptively:

```text
kernel operation invoked
kernel operation succeeded
kernel operation intentionally unsupported
local operation invoked
kernel/userland boundary
```

No ratio.

No threshold.

No "below X% ⇒ grow kernel."

The Growth Gate concerns **invariants**, not quantity of code or percentage of operations performed by the kernel.

---

# 11. Schwarzschild should be derivation-first

The Qwen review raises a valuable distinction.

I recommend locking:

```text
Primary:
    derive Schwarzschild from the general static spherical ansatz

Secondary:
    verify the resulting Schwarzschild solution independently
```

The primary derivation starts with:

```text
ds² =
    -A(r) dt²
    +B(r) dr²
    +r² dΩ²
```

with:

```text
A(r), B(r)
```

unknown.

Then:

```text
A,B
 ↓
Christoffel
 ↓
Riemann
 ↓
Ricci
 ↓
vacuum equations
 ↓
bounded differential/algebraic reduction
 ↓
Schwarzschild form
```

The local solver must be **problem-specific and bounded**. Do not build a general ODE solver.

The closed-form Schwarzschild metric then becomes a separate verification target.

This is stronger than merely plugging the known answer into `Rμν` and discovering that the algebra engine can simplify it. The Qwen review correctly notes that verification alone doesn't exercise unknown-function differentiation or the actual reduction path. GR_PLAN_qwen

---

# 12. Newtonian limit is NOT `phys.Limit`

This needs to be removed from the current operation mapping.

The frozen `Limit` operation is specifically bounded around direct substitution and the fixed Lorentz-factor body; it is not an infinity/asymptotic engine. specs_v2_3

So GR-8 should use:

```text
GR-local weak-field reduction
+
exact correspondence extraction
+
kernel dimension audit where possible
```

not:

```text
phys.Limit(...)
```

The plan should include two distinct checkpoints:

### Exact flat-limit checkpoint

```text
M → 0
```

giving the flat spherical metric.

### Weak-field/Newtonian correspondence

Introduce:

```text
ε = GM/(c²r)
```

and perform a bounded first-order reduction in `phys-gr`.

Recover:

```text
Φ = -GM/r
```

and, preferably, the corresponding radial acceleration:

```text
a_r = -GM/r²
```

The approximation/regime semantics remain Level 1. Do not pretend the frozen `APPROXIMATED` provenance capability has been activated.

This directly addresses the reviews' concern about the wrong `Limit` mapping. GR_PLAN_claude2

---

# 13. Units: pin them, but keep tensor dimensions out of kernel semantics

I recommend:

> **GR tensor calculation uses geometrized units `G = c = 1`; the GR-8 physical-correspondence subtrack restores SI dimensions for the scalar Newtonian quantities.**

This keeps the tensor algebra tractable without pretending the kernel's scalar `Dimension` is a tensor-unit system.

The SI audit should explicitly test:

```text
GM/(c²r) → dimensionless
GM/r²     → acceleration
```

using HYPOTHESIS `core.Object` carriers where the expression is kernel-representable.

Coordinate-component dimensions remain GR-local metadata because a tensor component's dimensional behavior depends on its index/coordinate basis.

This turns the reviewers' "units undecided" observation into a useful two-layer test rather than forcing tensor dimensions into the kernel. GR_PLAN_claude2

---

# 14. Tensor structure must be concrete before GR-2

Qwen is right that "emergent tensor algebra" is too vague for a coding agent.

Pin the Level-1 representation:

```text
Tensor:
    Rank
    IndexSlots[]
    Components
    Symmetries
    ChartID
```

Each `IndexSlot` has:

```text
slot number
variance: covariant | contravariant
index domain
```

Component keys are fixed tuples, not arbitrary strings.

Required operations:

```text
Contract
RaiseIndex
LowerIndex
ApplyMetric
ApplyInverseMetric
```

The GR package must enforce:

```text
same chart
compatible index domains
valid contraction pairing
variance correctness
declared symmetries
```

Also explicitly:

> **Christoffel symbols are connection coefficients, not tensors.**

This becomes an adversarial probe.

The reviewers specifically identified missing probes around mixed charts, treating Γ as a tensor, symmetry errors, and sign conventions. GR_PLAN_claude

---

# 15. Conventions are part of the GR contract

Pin:

```text
coordinate order:
    (x⁰, r, θ, φ)

metric signature:
    -+++

GR curvature convention:
    explicitly defined in phys-gr
```

The existing SR package already uses the `metric.signature = -+++` convention, so GR should use the same signature baseline. The kernel itself does not acquire GR curvature semantics. GR_PLAN_claude2

Include convention mismatch probes:

```text
wrong signature
wrong Riemann sign convention
mixed convention sets
mixed charts
```

Where scalar GR objects are bridged into `core.Object`, attach compatible `ConventionSet` metadata so the frozen kernel can detect conflicts.

---

# 16. Silent-wrongness suite: correct its semantics

The field-function probe should not claim:

```text
core.Symbol("f")
→ derivative 0
→ SILENTLY-WRONG
```

That is not a kernel bug.

Instead:

```text
Probe A:
    core Symbol("f") differentiated wrt r
    → 0
    → ACCEPTED-CORRECT

Probe B:
    GR UnknownFunction("f", r)
    → f'(r)
    → expected Level-1 behavior

Probe C:
    malformed or ambiguous GR-local dependency representation
    → should be rejected/diagnosed locally
```

Then reserve `SILENTLY-WRONG` for:

> a representation that is legitimate under the defined Level-1 API and nevertheless produces a mathematically wrong accepted result.

This keeps the taxonomy scientifically useful rather than turning user misuse into evidence against the kernel. GR_PLAN_deepseek

---

# 17. Level-2 outcome must be changed

This is a genuine contradiction identified by DeepSeek.

The plan currently permits:

```text
LEVEL-2-GROWTH
```

but simultaneously requires a named second established consumer, which a GR-only execution cannot necessarily establish. GR_PLAN_deepseek

Change the terminal outcomes to:

```text
NO-GROWTH

LEVEL-2-CANDIDATE-PENDING

KERNEL-GROWTH-CANDIDATE

ESCALATE-TO-SPEC
```

`LEVEL-2-CANDIDATE-PENDING` means:

> GR has exposed an abstraction that appears reusable, but promotion cannot occur until a named second established theory provides a worked example and the human curator approves the shared module.

If an existing established theory (`mechanics` or `relativity`) is genuinely demonstrated as the second consumer during this work, record that evidence, but **do not silently promote the library inside this plan**. Promotion becomes a separate human-curated action.

This makes Level 1 → Level 2 and Level 2 → Level 3 genuinely sequential.

---

# 18. Growth evidence must be independently classified

The coding agent cannot both discover and adjudicate a potential kernel-growth finding.

Therefore:

```text
PACKAGE-SOLVABLE
    may be classified by implementation agent

anything else:
    independent reviewer required
```

For:

```text
REPRESENTABLE-BUT-UNFAITHFUL
UNREPRESENTABLE
SILENTLY-WRONG
```

the evidence record must contain:

```text
agent classification
+
independent classification
+
disagreement resolution
```

before it becomes a candidate for Level-3 consideration.

That addresses the confirmation-bias concern without artificially presupposing kernel growth.

---

# 19. Define the 20-field evidence artifacts precisely

The previous "exact diff-shape" wording is too vague.

Field 18 should mean:

```text
MinimalKernelChange:
    affected package/file
    affected exported identifier(s), if any
    change kind:
        add API
        alter invariant
        add operation
        add expression node
        alter canonical encoding
        alter replay semantics
        alter MRC
    no source code
```

Field 20:

```text
HumanApproval:
    record ID
    reviewer name/identifier
    UTC timestamp
    decision:
        approved
        rejected
        provisional
    decision note
```

No cryptographic signing is required by the MVP.

A record without human approval remains **PROVISIONAL** and cannot trigger a kernel modification.

---

# 20. GR artifact hashing

This was correctly flagged as missing.

`phys-gr` must define canonical JSON and SHA-256 hashing for:

```text
GRExpr
Index
Tensor
Metric
Connection
Curvature
GR derivation steps
```

The format is Level 1.

It does **not** become part of `phys` canonical artifacts.

A GR derivation can therefore be audited:

```text
GR local artifact
       ↓
canonical JSON
       ↓
SHA-256
       ↓
evidence record
```

while kernel-compatible scalar bridges retain their ordinary `core.Object` hashes.

---

# 21. Baseline and mutation discipline

The reviewers' baseline concern is valid, but the "manifest re-hashing" concern needs to be narrowed.

The actual repository has an operational 42-file re-freeze authorized after the three documentation additions; the independent re-review explicitly confirmed exact 42-file membership and the closed-world tree test. adv_review11(2) adv_review11(2)

So the plan should say:

```text
v2.3 specification:
    39 frozen implementation files

operational repository freeze:
    exact 42 files
    39 implementation + 3 authorized documentation files
```

PASS0 must verify before editing that its three target files are already members of that 42-file set.

Then produce a **new post-PASS0 baseline record** containing:

```text
phys git commit SHA
phys working-tree-clean status
SHA-256 of specs_v2_3.md
SHA-256 of both manifests
SHA-256 of each modified PASS0 file
exact 42-file inventory
go build result
go vet result
go test result
```

The GR test evidence records the exact `phys` commit/baseline it used.

For mutations:

```text
never mutate authoritative phys checkout

use disposable clone/worktree
run mutant
run tests
discard mutant
re-verify authoritative checkout
```

This is much cleaner than trying to re-manifest the corpus manifests merely because tests changed.

---

# 22. Local module mechanics

The phrase "No replace into kernel internals" should be removed because `replace` addresses modules, not individual internal packages.

Pin:

```text
local development:
    go.work connecting sibling phys and phys-gr

production/CI:
    exact phys module version/commit

forbidden:
    replace to a modified phys fork
    vendored phys
    direct internal/kernel import
```

Before a GR run:

```text
phys checkout == recorded clean baseline
```

This directly addresses the build/immutability concerns. GR_PLAN_deepseek

---

# 23. The frozen 42-file `phys` tree stays untouched by GR

This now means:

```text
GR implementation
    → phys-gr only

GR-local tests
    → phys-gr only

GR-local evidence
    → phys-gr only

shared Level-2 candidates
    → separate future Go module

kernel changes
    → impossible without Growth Gate
```

PASS0 is the one pre-GR exception because it is a separate prerequisite gate against the already-existing `phys` test/documentation files.

The actual frozen architecture already states that trusted objects originate through fixed domain constructors, hypothesis construction, pure operations, or session artifacts, and there is no generic public object factory. specs_v2_3

---

# 24. Revised GR pass structure

The final execution contract should therefore look like:

```text
PASS0
Baseline Regression Verification
    ↓
PASS1
GR Plan Verification
    ↓
GR-0
representation + conventions + units + public API audit
    ↓
GR-1
static spherical metric
    ↓
GR-2
inverse metric + concrete tensor/index engine
    ↓
GR-3
Christoffel
    ├── 3a kernel negative-power probe
    └── 3b GR-local differentiation
    ↓
GR-4
covariant derivative
    ↓
GR-5
Riemann
    ↓
GR-6
Ricci + scalar curvature
    ↓
GR-7
Einstein tensor + vacuum equations
    ↓
GR-8
Schwarzschild derivation
    +
independent Schwarzschild verification
    +
flat M→0 check
    +
weak-field/Newtonian correspondence
    ↓
Independent adversarial suite
    ↓
Growth Gate
```

This incorporates the Qwen requirement that Level-1 mathematical machinery be concrete rather than merely "emergent." GR_PLAN_qwen

---

# Final consolidated decisions

After removing duplicates, false positives, and recommendations that conflict with the architecture, I would lock these **18 decisions**:

```text
D1   phys kernel remains math-/theory-agnostic.

D2   GR is Level 1 and lives in sibling module github.com/PithomLabs/phys-gr.

D3   Existing mechanics/relativity inside phys are grandfathered.

D4   GR physical structures remain Level-1 structures, not kernel Kinds.

D5   Kernel-compatible GR scalar subexpressions may be wrapped as
     HYPOTHESIS core.Objects with KindExpression and explicit Dimension.

D6   Non-kernel GR expressions use a single bounded GR-local algebra.

D7   GR-local algebra includes Sin, Cos, UnknownFunction.

D8   No opaque sin_theta/function-name placeholder representation.

D9   GR-3a must exercise frozen phys.Differentiate before GR-3b local continuation.

D10  GR-local operations never modify, shadow, or re-export kernel operations.

D11  GR-local mathematics has its own deterministic replay/hash trace.

D12  Kernel Simplify/arithmetic are reused whenever the subexpression is
     genuinely kernel-representable.

D13  Tensor/index algebra is concrete Level-1 infrastructure, not test-file maps.

D14  Schwarzschild is derived from unknown A(r), B(r), then independently verified.

D15  Newtonian correspondence is GR-local bounded weak-field reduction,
     not phys.Limit.

D16  No substrate-efficacy ratio or threshold.

D17  Level-2 promotion is never automatic; GR may produce only a
     LEVEL-2-CANDIDATE-PENDING requiring a named second consumer and human curation.

D18  Kernel growth requires full evidence, independent review, and human approval.
```

## Overall disposition of the six reviews

The reviews are therefore **not asking for another architectural redesign**.

They have identified a handful of genuine missing decisions:

```text
GR ↔ core.Object bridge          → LOCKED
GR symbolic representation       → LOCKED
unknown-function representation  → LOCKED
tensor representation            → LOCKED
GR-local replay                  → LOCKED
Newtonian-limit mechanism        → LOCKED
units/conventions                → LOCKED
Level-2 terminal state           → LOCKED
baseline/mutation protocol       → LOCKED
review/approval protocol         → LOCKED
```

And several suggestions should explicitly **not** be adopted:

```text
❌ substrate efficacy ratio
❌ kernel tensor types
❌ kernel trig nodes
❌ kernel unknown-function nodes
❌ placeholder sin_theta representation
❌ weakening Substitute
❌ extending phys.Differentiate for GR
❌ expanding phys.Limit for GR
❌ automatic Level-2 promotion
❌ treating bare-symbol differentiation as a kernel silent-wrongness bug
```

That is the important consolidation. The reviews' strongest shared observation is that the current plan is **architecturally right but operationally under-specified**; the changes above turn those implicit assumptions into explicit implementation contracts. GR_PLAN_deepseek GR_PLAN_z

With these decisions, the target condition is finally achievable:

> **A coding agent should have no remaining architectural choice to make while implementing GR. It may discover mathematical difficulty, implementation bugs, performance limits, or genuine kernel-boundary evidence—but it must not have to decide what belongs in the kernel, how GR interoperates with `phys`, what the symbolic representation is, how replay works, or what constitutes promotion.**
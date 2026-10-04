## prompt

now we go to the kernel growth gate, as we tackle general relativity as our growth gate target, that would inform us of that gate criteria and just to clarify do we need to tackle the issue per below? But it does **not yet have a full semantic type system for physics**.

For example, it does not currently know that:
```
position x
velocity v
four-vector Xμ
metric tensor gμν
scalar field φ
vector field Vμ
operator H
wavefunction ψ
```

belong to different mathematical structures with rules for index variance, tensor rank, manifold domain, coordinate charts, function spaces, operator domains, etc.

That's a much deeper level of category checking than the current `Kind + Dimension` model.

And this distinction matters for your longer-term TOE/compiler idea:

> **The MVP already has the seed of a physics type system, but it is not yet a full mathematical-physics type system.**

The current architecture gives you a particularly good foundation for extending that later without confusing it with ordinary algebra.

In fact, the separation:
```
Kind
Dimension
Expr
Assumptions
Conventions
Provenance
```

is quite useful because a future **kernel growth gate** could ask whether a proposed new invariant is genuinely foundational and theory-neutral, rather than dumping every new physical category into the kernel.

So, at the current frozen level: **yes, category-error prevention is built in and tested; but it is currently strongest at the physical-quantity level, not yet at the full tensor/mathematical-structure level.**


## answer

Yes. **General Relativity should be the first serious workload for the Kernel Growth Gate, but the absence of a full mathematical-physics type system should *not* itself trigger kernel growth.**

The right sequence is:

> **Do not build a tensor/manifold/type system first. Make GR try to use the frozen kernel exactly as it stands, record what fails, and let those failures supply the evidence for the Growth Gate.**

The distinction is important.

### What the frozen kernel already guarantees

The current substrate already separates at least two kinds of semantic checking:

```text
Dimension
Physical Kind
```

with operation-time MRC checks. `Add`/`Subtract`/`Compare` require dimensional compatibility, MRC-003 separately checks physical-kind compatibility, and `Substitute` requires both dimension and kind agreement between variable and replacement. plan10_v2_3 physics_compiler_mvp_specs_v2_1

The symbolic layer is intentionally very small: exactly ten expression node types, with bounded calculus and no general CAS/prover. specs_v2_3

That is enough for the current mechanics/SR corpus. It is **not enough to express full tensor/manifold semantics faithfully**.

### What we should *not* do

Do not now say:

> “GR needs tensors, therefore add Tensor, Index, Manifold, Metric, Connection, Curvature, CoordinateChart, ScalarField, VectorField, … to `internal/kernel`.”

That would completely invert the Growth Gate.

It would be **kernel design driven by the needs of one theory**.

Instead, GR becomes the experiment:

```text
Frozen Kernel
      ↓
General Relativity package
      ↓
Attempt meaningful GR derivations
      ↓
Record semantic failures
      ↓
Can existing substrate express them faithfully?
      ↓
        NO
        ↓
Can GR package / user library solve it without weakening guarantees?
      ↓
   YES → keep outside kernel
      |
      NO
      ↓
Is missing capability generic, foundational,
theory-neutral, reusable, and minimal?
      ↓
   YES → candidate kernel growth
```

That is the actual Growth Gate.

## The key distinction: representable vs faithfully typable

This is where your quoted issue becomes useful.

There is a huge difference between:

```text
“Can I encode gμν somehow?”
```

and:

```text
“Can the system know that gμν is a rank-2 tensor,
with the required index structure, and reject
an operation that is mathematically meaningless?”
```

The first is easy: encode something into an `Expr`.

The second is the actual question.

For example, a future GR derivation may need to distinguish:

```text
scalar
vector
covector
rank-(1,1) tensor
metric tensor
connection coefficients
Riemann tensor
Ricci tensor
scalar curvature
```

and ultimately enforce things such as index compatibility, contraction legality, tensor rank, and transformation behavior.

The current `Kind + Dimension + Expr` model does not express those structures today. That is a **known capability boundary**, not yet a demonstrated reason for kernel growth.

## GR should therefore be a "growth-gate workload," not a "kernel-expansion project"

I would define the GR phase in four passes.

### Pass 1 — Try to stay entirely outside the kernel

Create a prospective:

```text
general_relativity/
```

package.

Do not modify `internal/kernel`.

Try to represent a carefully selected set of established GR concepts and derivations using the existing substrate plus package-local structures.

The package should cover enough genuine GR to expose the boundary, not merely create placeholder types.

For example:

```text
manifold
metric
connection
curvature
Einstein tensor
stress-energy tensor
Einstein field equation
```

Then attempt actual derivations rather than just constructing data structures.

### Pass 2 — Identify the first *semantic* failure

For each failure ask:

> Is this merely missing domain vocabulary, or is the frozen formal substrate incapable of enforcing the required invariant?

For example:

```text
GR package can define “Metric”
```

is not a kernel problem.

But:

```text
GR package cannot make tensor contraction
type-safe because core.Object has no generic mathematical
structure capable of representing or checking the required
index constraints
```

could become a kernel candidate.

That distinction is crucial.

### Pass 3 — Apply the Growth Gate to every proposed primitive

Every proposed kernel addition must pass all of these:

**Foundational**  
It represents a mathematical/formal concept rather than a GR law.

**Theory-neutral**  
It must make sense outside GR.

**Reusable**  
At least multiple future frameworks or mathematical domains should plausibly require it.

**Semantically necessary**  
Without it, the formal system cannot enforce an invariant that matters.

**Not replaceable externally**  
A theory package, user library, adapter, or higher layer cannot provide the same guarantee without duplicating or bypassing trusted substrate semantics.

**Minimal**  
Add the smallest primitive/invariant that closes the demonstrated gap.

**No theorem leakage**  
The primitive must not encode Einstein equations, equivalence principle conclusions, curvature theorems, conservation theorems, or particular GR results.

This preserves the distinction you already established: **GR belongs in `general_relativity/`; generic mathematical substrate needed by GR and other theories may qualify for kernel growth.**

### Pass 4 — Require a concrete counterexample

This is perhaps the most important gate criterion.

A proposed kernel change should come with:

```text
BEFORE
What invalid or ambiguous state can currently be represented?

FAILURE
What actual GR derivation demonstrates the problem?

OUTSIDE-KERNEL ATTEMPT
Why can't the GR package/user library enforce the invariant?

CANDIDATE PRIMITIVE
What is the smallest generic capability required?

CROSS-THEORY TEST
Where would the same capability be useful outside GR?

AFTER
What invalid state becomes impossible or formally rejectable?
```

That is the physics equivalent of the Solvent “new durable fact / atomic transition” test. The Solvent precedent explicitly says kernel changes should occur only when the new invariant cannot safely be expressed outside the kernel. solvent-agentjacking-v0-vision-…

## So what about the "full semantic type system"?

**Do not make that a prerequisite.**

Instead, make it an **output of the GR experiment**.

There are three possible outcomes.

### Outcome A — GR fits

You discover that the required mathematics can be safely represented through domain/user-level abstractions while retaining the kernel's guarantees.

Then:

```text
Kernel growth = 0
```

and that is a successful Growth Gate result.

### Outcome B — one or two generic primitives are missing

For example, GR demonstrates a genuinely reusable formal capability that cannot be expressed outside the kernel without weakening correctness.

Then add **only those primitives**.

Not “the GR type system.”

Something more like:

```text
missing generic invariant
        ↓
smallest generic kernel primitive
        ↓
GR package uses it
        ↓
future theories can use it too
```

### Outcome C — the current mathematical model is fundamentally too weak

Suppose GR demonstrates that meaningful physics derivation requires a general mathematical-structure system that cannot coexist with the current object model.

Then that is a much larger finding.

At that point the Growth Gate should **stop rather than automatically grow the kernel**.

Because the conclusion may be:

> the frozen MVP's scalar physical-object substrate has reached its intended boundary and a new mathematical layer should be designed above it.

That could be a better architectural answer than turning `internal/kernel` into a massive tensor-calculus kernel.

## This also changes how I would define "category error"

Your earlier observation is correct, but we should sharpen it.

Today the system protects against:

```text
Dimension category errors
Physical-kind category errors
Assumption conflicts
Convention conflicts
Provenance/epistemic category errors
Operation-shape errors
```

The tests explicitly exercise dimension rejection and equal-dimension/different-kind rejection. plan10_v2_3

GR will ask for a deeper layer:

```text
mathematical-structure category errors
```

such as:

```text
scalar ↔ vector
vector ↔ covector
rank-1 ↔ rank-2 tensor
free index ↔ contracted index
tensor on one manifold ↔ tensor on another manifold
coordinate expression ↔ invariant object
```

But **we should not decide in advance that these all belong in the kernel.**

That is exactly what the Growth Gate is supposed to determine.

### The principle I would lock now

> **General Relativity is the stress test for kernel sufficiency, not the justification for kernel expansion.**

And:

> **The absence of a full mathematical-physics type system is a Growth Gate question, not a pre-Growth Gate defect.**

The current `Kind + Dimension + Expr + Assumptions + Conventions + Provenance` architecture gives us enough machinery to run that experiment without prematurely committing the kernel to a tensor-calculus ontology. AGENTS

The next thing I would design is therefore the **formal Kernel Growth Gate itself**: a GR-driven decision matrix that an implementation/review agent must execute for every proposed new primitive, with explicit PASS/FAIL evidence requirements and a hard “no kernel growth” outcome.


## Z

# Adversarial Review: GR-as-Growth-Gate-Workload Proposal

## Verdict

**The principle is correct and should be adopted** — "GR is the stress test for kernel sufficiency, not the justification for kernel expansion" is the right lock, and it's exactly what AGENTS.md §13 anticipated. But the proposal has **one structural blind spot (kernel vs. ops vs. spec growth), one methodology flaw that would pollute the entire experiment (unclassified failure logs), and it underestimates how much GR will land in Outcome A**. Most importantly: the strongest genuine kernel-growth evidence already exists and should be pre-registered now — the `Differentiate` constant-symbol semantics, where the external fix is *blocked by the project's own purity rule*. That's the benchmark counterexample the matrix should be built around.

---

## 1. What's right (verified against the frozen documents)

- The representable-vs-faithfully-typable distinction is sound and matches the spec's closed-world posture (§8.1, §10.2, REQ-008-01).
- The seven Pass-3 criteria are consistent with AGENTS.md §13's growth discipline, and Pass 4's BEFORE/AFTER counterexample format is the correct operationalization — it's the project's own falsifiability ethos applied to kernel evolution.
- Outcome C ("stop rather than grow; design a layer above") is the most valuable part of the proposal. Most growth-gate designs lack a legitimate no-go exit; this one has it.
- Not touching the type system now is correct per REQ-002-21 and REQ-002-04 (both still normative non-goals).

---

## 2. 🔴 F1 — The gate gates the wrong axis: kernel ≠ ops ≠ spec

The proposal's entire decision machinery targets `internal/kernel`. But map where GR friction will actually land:

| Growth axis | Examples GR will demand | Existing gate coverage |
|---|---|---|
| **ops** (new rules in existing ops) | negative-integer power rule for `Differentiate`; more `Solve` patterns; more `Limit` bodies | ❌ none — silently permitted or silently forbidden |
| **kernel** (new node types, metadata, enums) | field-aware symbols; second `Call` id | ✅ Pass 3/4 |
| **spec** (canonical JSON §10.2 is closed; ordinals pinned) | *any* new Expr node kind | ❌ unaddressed — and this is the killer |

The third row is the structural omission: a new expression node type is not a code change. §10.2 says "No other expression JSON kind is allowed"; §8.2 ordinals are "stable for mrc-v0.4 artifacts." Kernel growth is **artifact-breaking**: every existing hash, ledger, and sealed candidate depends on the closed canonical form. The decision matrix must include (a) an **append-only ordinal policy** (new nodes get ordinals 10+, never renumber — preserves old artifacts' sort order), (b) an **mrc-version bump + replay/migration policy** as a mandatory cost line in any GROW verdict, and (c) recognition that GROW on a node type is a *spec amendment*, which per REQ-000-03 routes to human authority, not to the implementing agent.

Without the three-track split, the experiment resolves in one of two bad ways: ops extensions get smuggled in without any gate (most GR friction lives there), or everything gets the full kernel ceremony and the team drowns. **Fix: three gates with different thresholds — ops rules (generic + tested + fail-closed ⇒ lightweight), kernel primitives (full counterexample ceremony + spec amendment), non-goal changes (human decision only).**

Related ambiguity to pin *before* the experiment: spec §15.8's unsupported list says "non-integer symbolic exponent" while the rule is stated only for non-negative integers — **is `Differentiate` of `Pow(x,-1)` an error or a missing rule?** The plan text (§6) implies unsupported; the spec's unsupported list literally excludes negative integers from the ban. This is load-bearing for GR (every rational-function metric hits it in step one) and must be pinned now, not discovered mid-experiment.

---

## 3. 🟠 F2 — Pass 2's triage will drown: pre-register the failure taxonomy

Pass 1's failure log will be **dominated by spec-intended boundedness rejections** — rejections the spec *designed*, which are not evidence of kernel insufficiency. If Pass 2 triages every `UnsupportedOperationError` as a candidate finding, the experiment either triggers premature growth or exhausts itself on noise. I can predict the first session's log:

| # | Predicted first failure | Class |
|---|---|---|
| 1 | `d/dr` of `Pow(r,-1)` — i.e., **any rational function** (Schwarzschild `g_rr`) | Spec-intended bound (ops) + spec ambiguity (§2) |
| 2 | `d/dt` of `t^(4/3)` — FRW matter era | Spec-intended bound |
| 3 | `Differentiate` of `Sqrt`/`Call` | Spec-intended bound |
| 4 | Second `Call` id rejected (`sin`, `exp`, Γ) at kernel validation | Representation gap (kernel-adjacent) |
| 5 | `Solve` accepts one quadratic pattern; no linear solves (inverting `g_μν`) | Spec-intended bound — **package workaround exists** (Cramer's rule over explicit rational entries is pure `Add/Mul/Pow`, within kernel capability) |
| 6 | No general `Limit` (Newtonian limit of `g_00`) | Spec-intended bound |
| 7 | **`Differentiate` treating a field symbol as constant → silently wrong result** | **Silent-wrongness hazard — see F3** |
| 8 | Index structure expressible only as opaque symbol names | Representable-but-unfaithful — **but see F4** |
| 9 | Conventions are inert metadata (`metric.signature` never consulted by entailment) | Candidate evidence class |

**Fix:** require the experiment plan to pre-register exactly this taxonomy with five classes — *spec-intended bound / representable-but-unfaithful / unrepresentable / silently-wrong / package-solvable* — and rule that **only the last three classes constitute gate evidence**. Note class 5's workaround: this is why feasibility of Pass 1 must be pinned to *concrete-metric* derivations (explicit coordinate functions, Cramer inversions) and explicitly exclude abstract-index manipulations, which are infeasible in Pass 1 no matter what.

---

## 4. 🔴 F3 — The benchmark counterexample already exists; register it now

The single most likely genuine kernel-growth candidate — and the perfect seed for the Pass-4 format:

- **BEFORE:** `Differentiate(f, x)` where `f` contains `Symbol("g_xx")` returns `0` **silently** — a wrong derivative certified `DERIVED` with valid provenance. Every other failure mode in the system fails closed; this one fails *wrong*.
- **FAILURE:** any abstract-metric GR derivation (Christoffel symbols).
- **OUTSIDE-KERNEL ATTEMPT:** package discipline ("only differentiate coordinate-explicit expressions") works for curated workflows — but the *public operation* remains silently wrong for every other caller. The natural package-level fix — a registry of which symbols are fields — is **banned by REQ-015-01** (no ambient state). This is the sharpest result available: **the purity requirement that makes the kernel trustworthy is precisely what blocks the external fix.** That is textbook "not replaceable externally."
- **CANDIDATE:** minimal field-semantics primitive (or fail-closed rejection of non-coordinate symbols under differentiation — note even that needs a "which symbols are coordinates" carrier, so it is not obviously smaller).
- **CROSS-THEORY:** continuum mechanics, EM, QM wavefunctions — concrete second consumers.
- **AFTER:** silent wrongness eliminated; fail-closed preserved.

This example also motivates the **missing eighth gate criterion: *no silent failure mode*** — a new primitive must convert currently-silent wrong states into errors, or prove it introduces none. Fail-closed preservation is the deepest invariant the MVP owns; the gate must protect it explicitly, because it's the invariant most easily sacrificed for expressiveness.

---

## 5. 🟠 F4 — The proposal underestimates Outcome A: index safety is externally enforceable

The proposal treats "cannot make contraction type-safe" as a prima facie kernel candidate. Attack: a `general_relativity/` package can enforce index structure **in Go, at package level** — a typed `Tensor` struct with index kinds that validates contractions before emitting per-component scalar `core.Object`s. Within the project's trust model (REQ-000-03/04: in-repo packages are trusted; `ops`-level MRC enforcement is the *same* trust tier), package-level index checking is exactly as strong as `ops`-level dimension/kind checking is today. The kernel never sees index structure; each emitted object is a well-formed scalar expression carrying full provenance and canonicality. The corpus layer is likewise reusable without growth (typed manifests are theory-agnostic — GR needs no kernel change to get machine-readable corpus).

Consequence: the "not replaceable externally" criterion must demand a **worked bypass analysis** — show the invariant cannot be enforced at package level *within the trust model and without violating a spec rule* (as F3's REQ-015-01 tension does). Index algebra will largely fail that test — meaning **Outcome A plausibly covers more of GR than the proposal implies, and Outcome C is more likely than Outcome B.** Set that expectation now so a thin evidence yield isn't read as experiment failure. (One category error from the proposal's list is worth flagging as likely-permanently-out-of-scope: manifold identity of tensors — that's structural mathematical identity, deep water; classify as package convention.)

---

## 6. 🟡 F5 — Matrix hardening (beyond the eight criteria)

1. **Merge Pass 3 and Pass 4** into one instrument — Pass 4's BEFORE/AFTER *is* the operationalization of "semantically necessary"; as separate passes it invites criteria-passing without evidence.
2. **Replace "plausibly require it"** (reusable criterion) with **a named second consumer plus a worked example**. "Plausibly" does no work; the corpus auditor would reject it as prose-load-bearing fact — the exact thing this project bans.
3. **Non-goal collision check as an explicit gate question:** does the candidate constitute or enable REQ-002-21 (general tensor/index algebra) or REQ-002-04? The likeliest failure of the whole experiment is month-three, when the "one small primitive" is named `Index`. Pre-commit: the gate may *recommend* non-goal revision but **cannot authorize it** — that routes to human/spec authority (REQ-000-03), same as Outcome C.
4. **Mutation-tested rejection paths** (extend Gate H methodology): a primitive isn't accepted until its rejection behavior survives the mutant battery in a disposable copy.
5. **Artifact-compat cost line** (from F1): mrc-version bump, append-only ordinals, replay/migration story — mandatory in every GROW verdict.

---

## 7. 🟡 F6 — Sequencing and governance gaps

- **Post-freeze, out-of-tree, explicit.** REQ-002-20 bans GR in the MVP; plan7's freeze discipline is "no new files." The GR experiment must be sequenced *after* freeze acceptance, in a **separate workspace/branch**, with the frozen MVP untouched. Otherwise this repeats the B1 file-guardrail breach pattern at larger scale. One sentence in the experiment charter fixes it.
- **Independent adversarial reviewer with veto.** The proposal mentions a review agent; make the separation structural — proposer cannot self-certify, reviewer runs the taxonomy and counterexample audits, reports in the established A–E/§17 format with verdict `GROW / NO-GROW / ESCALATE-TO-SPEC`.
- **The Solvent citation is dangling.** "solvent-agentjacking-v0-vision-…" is truncated and unverifiable in this packet. In a project whose discipline is provenance, the Growth Gate's founding precedent cannot be an unverifiable reference. Attach the document or mark the precedent as informal.
- **Outcome C needs an artifact clause:** the frozen MVP remains *valid* for scalar-quantity physics (mechanics, SR) — a new mathematical layer is additive, never a rewrite of what's frozen. Pre-commit that Outcome C is an acceptable, even likely, result (per F4), so the experiment isn't steered toward GROW to avoid "wasted" effort.

---

## Bottom line

Adopt the locked principle verbatim — it is the correct growth epistemology and consistent with everything this project has built. But before Pass 1 runs: (1) split the gate into ops/kernel/spec tracks with the artifact-compatibility policy, since most GR friction is ops-level and any node change is a spec amendment; (2) pre-register the failure taxonomy above so spec-intended rejections don't masquerade as evidence; (3) register the `Differentiate` field-symbol counterexample as the benchmark — it is currently the only known failure that is simultaneously *silently wrong*, *externally unfixable without violating REQ-015-01*, and *cross-theory*; (4) add the no-silent-failure criterion and the REQ-002-21 tripwire; (5) sequence the experiment post-freeze, out-of-tree, with an independent veto-holding reviewer. Run honestly, this experiment's most probable output is Outcome A for index structure plus a narrow, well-evidenced candidate list — and a likely Outcome C recommendation to design the mathematical layer *above* the frozen kernel rather than inside it. That would be the Growth Gate working exactly as intended.


## Gemini

While the dual `Dimension + Kind` architecture prevents basic operational category errors, several deeper structural and semantic edge cases in the MVP design limit its ability to catch subtle derivation errors:

### 1. "Type Decay" Neuters Category Checking in Multi-Step Derivations

When you perform operations on typed primitives (e.g., `Mass`, `Velocity`), the resulting object's `Kind` automatically degrades to `KindExpression` (or `KindRelation`).

* **The Problem:** Once an object becomes `KindExpression`, subsequent operations compare `Dimension` only, because `KindExpression` is compatible with any matching dimension.
* **Concrete Failure Mode:** Energy ($E$) and Torque ($\tau$) share identical dimensions ($M L^2 T^{-2}$). If $E$ is derived via $m c^2$ (`KindExpression`) and $\tau$ is derived via $r \times F$ (`KindExpression`), the system will happily allow `Add(E_expr, Torque_expr)`. The `CategoryMismatchError` protection is completely lost after the first operation step.

### 2. The Substitution Catch-22

The specification requires that a replacement object in `Substitute` must match both the **`Dimension`** and **`Kind`** of the target variable.

* **The Trap:** If variable $m$ has `Kind = RestMass`, substituting a derived algebraic quantity like $E / c^2$ creates a dilemma:
1. $E / c^2$ has `Kind = Expression`. If strict `Kind` matching is enforced, `Substitute` rejects $E / c^2$ because `KindExpression ≠ KindRestMass`. To fix this, you would be forced to run `Session.Identify` to manually re-type $E / c^2$ as `RestMass`.
2. If `Substitute` relaxes `Kind` matching to accept any `KindExpression` whose dimensions match $m$, it creates a category leak: any generic expression with dimension $M$ can now be substituted into $m$ without verifying if it represents a valid mass quantity.



### 3. Syntactic Key Collision vs. Semantic Assumption Conflicts

Assumption and convention conflict resolution relies on exact string key collisions under the same `(Kind, Key)` pair.

* **The Problem:** The engine checks string identity, not mathematical or logical satisfiability (SMT/SAT).
* **Concrete Failure Mode:** If one premise asserts `v << c` under key `"non_relativistic"` and another premise asserts `v = 0.99c` under key `"relativistic_speed"`, the assumption sets merge without throwing a conflict error. Incompatible physical limits can silently contaminate the same derivation unless human authors standardize every key string globally across all packages.

### 4. Global String-Based Symbol Collision Risk

`Symbol` nodes inside expression trees are identified purely by raw string names (`"m"`, `"c"`, `"p"`).

* **The Problem:** Symbols lack lexical scoping, namespaces, or instance identity.
* **Concrete Failure Mode:** In derivations involving multiple objects of the same physical type (e.g., mass of an electron $m_e$ vs. mass of a proton $m_p$, or target velocity $v_1$ vs. source velocity $v_2$), string substitution or simplification targeting `"m"` or `"v"` will globally match and mutate all identical string symbols in the expression tree, regardless of which physical entity they were intended to represent.

### 5. Transcendental Functions & Non-Rational Exponent Hazards

`Dimension` arithmetic uses exact rational power vectors ($M^a L^b T^c \dots$).

* **The Problem:** Physical functions like $\sin(\omega t)$, $\exp(-E / k T)$, or $\ln(V_2 / V_1)$ mathematically require strictly **dimensionless** arguments ($M^0 L^0 T^0$). Furthermore, raising an expression to a symbolic exponent (e.g., $x^n$ where $n$ is a variable symbol rather than a concrete rational number) breaks dimensional power evaluation unless $x$ is explicitly dimensionless.
* **Missing Safeguard:** Without explicit checks enforcing that transcendental function arguments and variable exponents evaluate to dimensionless quantities, invalid dimensional powers can propagate undetected through standard `Call` or `Pow` nodes.

### 6. Closed Ontology Matrix vs. Lack of Subtyping

The MVP relies on a fixed, closed set of `Kind` enums and an explicit matrix (`MRC-003`).

* **The Problem:** Physical categories exist in a hierarchy (e.g., `RestMass` $\subset$ `InvariantMass` $\subset$ `Mass` $\subset$ `ScalarQuantity`).
* **Limitation:** Because the engine uses discrete enum matching rather than an extensible subtyping lattice, adding any new domain or concept requires modifying the core kernel rather than defining hierarchical relationships in domain packages.
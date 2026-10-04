IMPORTANT: Reference the files mentioned here under plan10 folder but follow the procedures in this file, not in those reference files.

## Z

Yes. **Z adds genuinely net-valid insights, and two of them are more serious than the earlier “micro-pin” level.** After checking Z against v4 and the frozen `specs_v2_3.md`, I would treat Z as an additional hardening pass rather than merely commentary.

The key point is that Z found **two internal contradictions in v4 itself**:

1. **F2 `Subst` is too restrictive to execute the locked Schwarzschild script.**
2. **H3's expected `IDENTIFIED` provenance is impossible for HYPOTHESIS bridge operands under the frozen provenance law.**

Those are real findings, not stylistic preferences. GR_PLAN4_z

## 1. V4-1 — `Subst` must be broadened: **definitely adopt**

The current v4 says:

```text
Subst: Symbol -> Symbol
composite -> UnrepresentableSubstitutionError
```

But the actual GR-8 workload necessarily needs substitutions such as:

```text
A(r) → 1
B(r) → 1
B(r) → 1/A(r)
A(r) → 1 + k2/r
B(r) → (1 + k2/r)^-1
mu → eps*r
u → 0
```

and, crucially, derivative-bearing forms must follow the substitution consistently:

```text
A(r)       → F(r)
A'(r)      → F'(r)
A''(r)     → F''(r)
```

Z traces these conflicts directly against Steps 8, 10–14 and F3/F4. GR_PLAN4_z

So **the principle of Z's correction is valid**:

```text
Subst target:
  Symbol
  OR zero-order UnknownFunction atom

replacement:
  arbitrary GRExpr

composite expression as target:
  error

if substituting UnknownFunction f(r) with F(r):
  f(r)      → F(r)
  f'(r)     → Diff(F,r)
  f''(r)    → Diff(Diff(F,r),r)

max derivative order = 2
negative-power requirements still enforced
result normalized
```

One refinement: I would **not blindly copy Z's “UnknownFunction f → g” wording** because the replacement need not be another function `g`; it can be an arbitrary GRExpr such as `1+k2/r`. The normative rule should be **function-atom substitution with derivative closure**, not function-to-function substitution.

That is a much better coding-agent contract.

---

## 2. V4-2 — H3 provenance: **definitely adopt**

This one is unequivocal.

The frozen specification says:

```text
if either input status == HYPOTHESIS:
    Session.Identify output = HYPOTHESIS
```

and separately:

```text
if no input is HYPOTHESIS:
    Session.Identify output = IDENTIFIED
```

specs_v2_3

The GR bridge explicitly creates temporary **HYPOTHESIS** objects. Therefore the GR-8 correspondence cannot legitimately expect:

```text
KindRelation
IDENTIFIED
NONE
```

It must be:

```text
KindRelation
HYPOTHESIS
NONE
```

The identification *attempt* can still be recorded and reviewed; it simply cannot upgrade hypothesis-grade material to trusted `IDENTIFIED`. Z is exactly right on this. GR_PLAN4_z

So H3 should become:

```text
GR-8 Identify:
    operands = k2 ↔ -2mu
    purpose = integration-constant correspondence
    result Kind = Relation
    result Provenance = HYPOTHESIS
    CorpusStatus = NONE
    justification = required
```

The SI correspondence `k2 = -2GM/c²` should remain a dimensional-audit statement rather than a separate Identify event.

---

# 3. V4-3 — Explicit v3/v4 precedence: **adopt**

Z correctly notices that v3 is frozen, while v4 deliberately changes some v3 implementation wording.

Add:

```text
v3 remains immutable lineage.
Where v4 explicitly amends an implementation rule, v4 governs execution.
Frozen specs_v2_3.md remains authoritative over both.
```

This is not a new architecture. It prevents the coding agent from seeing two conflicting normative sentences and having to decide which one wins. GR_PLAN4_z

The hierarchy should be:

```text
specs_v2_3.md
    ↓ highest authority
v4
    ↓ execution hardening
v3
    ↓ immutable historical lineage
```

---

# 4. V4-4 — bounded distribution + normalization pass cap: **adopt**

This is a good implementation procedure.

The GR residual

\[
R_{tt}/A + R_{rr}/B
\]

requires enough algebraic expansion to expose the common denominator. Simply saying “distribution where needed” leaves too much discretion.

Pin:

```text
Distribution is permitted only for the bounded common-denominator
normalization required by the GR workload.

No general distributive expansion engine.

Normalize executes the fixed passes to a stable result, but with a
finite pass cap. If the cap is exceeded, return UNDECIDED rather
than looping.
```

Z's motivation is sound: deterministic termination is part of the failure-closed design. GR_PLAN4_z

---

# 5. V4-5 — stale-reference and vague-language sweep: **adopt**

This is excellent audit hygiene.

Z identifies actual dangling references not covered by the present A2/B1 sweep:

```text
§5 → §27
§8 → §27
§20 → §27
```

while the execution tables are actually in §24. GR_PLAN4_z

It also catches phrases that still leave discretion:

```text
preferably also a_r
(if any)
where representable
where included
```

These should be made deterministic where the workload already requires the behavior.

I would add an explicit **“normative-language sweep”** to the v4 validation procedure.

---

# 6. `Reduce` remains undefined: **adopt**

This is a straightforward catch.

v3 lists:

```text
Diff
Normalize
ZeroTest
Subst
Truncate
Reduce
```

but `Reduce` has no concrete contract.

There are only two clean options:

```text
A. Define Reduce precisely
B. Remove Reduce from the public GR-local operation vocabulary
```

For this project I would choose **B unless GR-8 actually requires a distinct `Reduce` operation**. Do not let the coding agent invent what `Reduce` means.

Z correctly identifies this as an implementation-choice gap. GR_PLAN4_z

---

# 7. Bridge assumption canonical keys: **adopt**

This is subtle but important.

Because assumptions participate in the BridgePayload hash, “we retain only bridgeable assumptions” is insufficient. Two agents could retain the same assumptions but serialize them differently.

Pin:

```text
BridgePayload.assumptions =
    the exact canonical representation of the retained structured
    AssumptionSet, using the frozen core canonical ordering.

No GR-specific assumption-string serialization.
No map iteration order.
No regenerated keys.
```

This reuses the frozen representation rather than inventing another assumption format.

---

# 8. `Index` vs `IndexSlot`: **adopt**

Z correctly catches a schema inventory mismatch.

Current v3 canonicalization inventory includes:

```text
Index
IndexSlot
Tensor
...
```

while the proposed Appendix D list replaces `Index` with `BridgePayload`. GR_PLAN4_z

This must be settled explicitly before coding.

I would retain only the types that actually exist in the v4 design:

```text
GRExpr
IndexSlot
Tensor
Metric
Connection
Curvature
GRStep
BridgePayload
```

and remove `Index` **unless v4 actually defines a distinct `Index` type elsewhere**.

---

# 9. Chart propagation: **adopt**

Pin:

```text
Every GR-0…GR-8 Metric, Tensor, Connection, and Curvature artifact:
ChartID = "spherical-static"
```

Cross-chart rejection is adversarial testing, not part of the primary workload.

This is a useful consistency rule because the chart identity otherwise exists mostly as metadata while tensor operations depend on it. GR_PLAN4_z

---

# 10. `OperationID` namespace: **adopt**

Because `OperationID` is hash-bearing, this should not remain free-form.

Use one exact namespace scheme. I would use:

```text
gr:diff
gr:normalize
gr:zero_test
gr:subst
gr:truncate
gr:contract
gr:raise_index
gr:lower_index
gr:apply_metric
gr:apply_inverse_metric
...
```

and:

```text
kernel:<frozen-operation-id>
```

for kernel contacts.

The specific capitalization in Z's example is less important than having **one exact namespace**.

---

# 11. Rational and AST invariant guards: **adopt**

Two useful implementation-level rules:

### Canonical zero

```text
Rational zero is always 0/1.
```

### Add/Mul arity

After normalization:

```text
0 terms → Rational(0)
1 term  → unwrap to that term
≥2      → Add/Mul node
```

This prevents canonicalization from producing an AST that violates the very node invariants the plan says are mandatory.

Z identifies both explicitly. GR_PLAN4_z

I would also add:

```text
No fixed-width rational arithmetic may silently wrap.
```

But I would **not** hard-code `int64` merely to solve this; the representation choice should remain consistent with the existing v4 contract.

---

# 12. `adversarial/` package structure: **adopt**

This is another very practical coding detail.

A directory containing only `_test.go` files does not constitute an ordinary Go package in the way the build tree may expect.

So decide now whether:

```text
adversarial/
    silent_wrongness_test.go
```

is replaced by tests adjacent to the packages they test, or the directory gets a non-test Go file.

I favor **adjacent tests**, because the adversarial suite can then exercise package boundaries deliberately using `_test` packages, while evidence remains documentation.

Also explicitly pin:

```text
go test ./...
go vet ./...
```

as the Level-1 gate.

---

# 13. Flat-space check: **adopt**

Z's residual list doesn't emphasize this, but its workload structure implies an excellent cheap check:

```text
mu → 0
```

must reduce the Schwarzschild metric to the spherical Minkowski metric.

This is already present in the broader GR workload history, so it is a good **implementation validation procedure**, not a new architecture. Gemini also includes this in the GR-8 thread. GR_PLAN4_gemini

---

# One thing I would **not** adopt from Z literally

Z says:

> “the only mechanical route” for Step 8 is substitution.

I would not make that wording normative.

The **problem it identifies is real**: F4's `u=1/r → 0` mechanism cannot operate while `A(r)` and `B(r)` are still unknown. Therefore the plan needs a bounded boundary-condition substitution mechanism.

But the architecture should say:

```text
Step 8 requires an explicit bounded boundary-condition evaluation.
This is implemented through the broadened GR-local Subst contract.
```

rather than claiming that substitution is mathematically the only conceivable mechanism.

That keeps the implementation contract precise without making an unnecessary philosophical claim.

---

# Net result

Z adds **more than procedural polish**. It changes our status from:

```text
v4 = mostly complete
```

to:

```text
v4 still has 2 genuine internal contradictions
+ several deterministic implementation pins
```

The most important additions are:

| Z finding | Disposition |
|---|---|
| F2 `Subst` contradiction | **ADOPT — critical** |
| H3 `IDENTIFIED` vs HYPOTHESIS | **ADOPT — critical** |
| v3/v4 precedence | **ADOPT** |
| bounded distribution + normalization cap | **ADOPT** |
| stale-reference sweep | **ADOPT** |
| undefined `Reduce` | **ADOPT** |
| canonical assumption keys | **ADOPT** |
| Index/IndexSlot reconciliation | **ADOPT** |
| chart propagation | **ADOPT** |
| OperationID namespace | **ADOPT** |
| rational/AST invariants | **ADOPT** |
| adversarial test-package layout | **ADOPT** |
| flat-space check | **ADOPT** |
| “only mechanical route” wording | **REFINE, don't adopt literally** |

Z's central conclusion is therefore correct: **do not execute the present v4 work plan yet.** The two contradictions must be corrected first, otherwise the coding agent is forced either to violate a normative pin or stop. GR_PLAN4_z

And importantly, **none of these require expanding the kernel or changing the overall GR architecture**. They are exactly the kind of implementation-contract corrections the hardening pass is supposed to catch. GR_PLAN4_z

## DeepSeek, Gemini, Qwen


Yes. After reading the **DeepSeek, Gemini, and Qwen v4 reviews as a set**, I would add a layer of **implementation guidelines/procedures** beyond the 28 hardening decisions. They do not require architectural changes; they make the coding agent less likely to produce a technically compliant but sloppy implementation.

The important point is to **use Gemini selectively**. Some of its details are excellent; a few are actually inconsistent with the already-frozen contract and should not be copied.

## What I would adopt

### 1. Make the candidate-and-certify procedure completely mechanical

This is one of the strongest net-valid additions from DeepSeek.

For every Schwarzschild candidate:

```text
candidate expression
    ↓
construct residual from authoritative current state
    ↓
Substitute candidate into residual
    ↓
Normalize
    ↓
ZeroTest
    ↓
ZERO       → certify candidate
UNDECIDED → fail closed
NONZERO   → reject candidate
```

Do not certify merely because the expression matches the expected pattern. This directly closes DeepSeek's R9 gap. GR_PLAN4_deepseek

I would make this an explicit coding-agent rule.

---

### 2. Make the primary/secondary Schwarzschild runs independently auditable

Gemini's derivation table is useful here, but I would strengthen it:

```text
PRIMARY:
    owns its own GRExpr/Tensor/Connection/Curvature objects
    owns its own replay trace

SECONDARY:
    constructs closed-form A,B from scratch
    owns separate objects
    owns separate replay trace
    MUST NOT reuse PRIMARY intermediate AST/tensor artifacts
```

The secondary is verification, not a recovery path for a failed primary. Gemini's detailed workload makes this distinction concrete. GR_PLAN4_gemini

---

### 3. Correct the “9 non-zero Christoffel components” wording

Gemini's table says **9 non-zero Christoffel components**, but for the dense `[rho][mu][nu]` representation there are **9 independent non-zero patterns under the lower-index symmetry** and more than 9 populated ordered array entries.

So the implementation guideline should say:

```text
GR-3b golden:
9 independent non-zero Γ^rho_mu_nu symmetry classes,
with Γ^rho_mu_nu = Γ^rho_nu_mu.
```

This is a useful correction rather than merely copying Gemini's wording. The underlying workload and metric are otherwise correctly specified. GR_PLAN4_gemini

---

### 4. Add explicit tensor validation-before-computation discipline

The reviews establish the tensor model, but the coding procedure can be clearer:

```text
Every tensor operation:

1. validate input tensor structure
2. validate ChartID
3. validate rank/slot compatibility
4. validate variance/contraction legality
5. validate declared symmetries
6. allocate fresh output
7. compute components
8. validate output invariants
```

No operation should partially modify an output before discovering an invalid input.

This fits the existing dense/symmetry design. GR_PLAN3_gemini

---

### 5. Pin a single component-index convention everywhere

Gemini's `Connection` layout is worth adopting:

```text
Connection.Components[rho][mu][nu]
    = Γ^rho_mu_nu
```

and the same positional convention should be used consistently in Christoffel, Riemann, Ricci contractions and canonical serialization. GR_PLAN4_gemini

I would state additionally:

```text
No package may introduce a second local index ordering.
```

---

### 6. Make metric/inverse verification componentwise and deterministic

For GR-2:

```text
g^{μα} g_{αν}
    → Contract
    → δ^μ_ν
    → ZeroTest each resulting residual component
```

The expected zero/nonzero pattern of the Kronecker delta should be part of the golden test, rather than checking only one component.

Gemini's explicit inverse metric and contraction workload makes this a natural implementation test. GR_PLAN4_gemini

---

### 7. Add curvature invariants as cheap internal checks

This is a good implementation-quality addition, without creating new mathematical machinery:

```text
R^rho_sigma_mu_nu = -R^rho_sigma_nu_mu
R_mu_nu = R_nu_mu
```

and for the Christoffel construction:

```text
Γ^rho_mu_nu = Γ^rho_nu_mu
```

These should be **validation checks**, not additional algebraic capabilities.

---

### 8. Make replay verification stronger than hash comparison alone

The existing v4 replay contract says canonical round-trip and hash chaining, but the coding procedure should explicitly verify:

```text
for each GRStep:

InputHash[i] == SHA256(InputCanonical[i])
OutputHash    == SHA256(OutputCanonical)

PreviousStepHash == prior CurrentStepHash

CurrentStepHash == recomputed hash

Decode(Canonical(x))
    → Encode(...)
    → byte-identical canonical bytes
```

Then perform operation-specific replay.

This turns the replay engine into an actual integrity verifier rather than a structure that merely stores hashes. The reviews correctly emphasize deterministic replay as a core boundary property. GR_PLAN4_deepseek

---

### 9. Make chart binding explicit across the complete workload

DeepSeek R11 is worth adopting:

```text
Every Metric, Tensor, Connection, and Curvature artifact
created during GR-0…GR-8 MUST have:

ChartID = "spherical-static"
```

Cross-chart behavior is tested only by the adversarial suite.

This prevents the canonical chart from being merely documentation while individual objects quietly use another chart. GR_PLAN4_deepseek

---

### 10. Pin GR-local `Subst` usage, not merely its syntax

DeepSeek R8 asks a good question: why have `Subst` if no workload uses it? GR_PLAN4_deepseek

I would retain `Subst`, because it is genuinely useful for the workload, but document concrete uses:

```text
Allowed workload uses:
- Symbol → Symbol coordinate renaming
- bounded algebraic substitutions used by candidate certification
- mu → eps*r in weak-field reduction

UnknownFunction arguments:
- only Symbol → Symbol
```

Do not let this become a general substitution engine.

---

### 11. Preserve explicit `UnknownFunction` dependency semantics

Gemini's differentiation table is useful as an implementation checklist:

```text
Diff(A(r), t) = 0
Diff(A(r), r) = A'(r)
Diff(A'(r), r) = A''(r)
Diff(A''(r), r) = DerivativeOrderExceededError
```

The important bit is that **`A`, `A'`, and `A''` are structurally related but distinct atoms** for normalization. Gemini's representation makes this explicit. GR_PLAN4_gemini

---

### 12. Make normalization one-way

Gemini's six-stage pipeline is useful:

```text
numeric
→ flatten
→ canonical ordering
→ rational collection
→ bounded trig reduction
→ zero reduction
```

but the agent should be instructed:

```text
No reverse trig identities.
No unrestricted distributive expansion.
No rewrite rule may reintroduce an earlier-stage form.
Normalization terminates after the fixed bounded pipeline.
```

This retains the useful part of Gemini while avoiding accidental CAS behavior. GR_PLAN4_gemini

---

### 13. Keep the bounded inverse-binomial rule extremely narrow

The current v4 already has:

```text
Pow(1+u,-1) → 1-u
```

Gemini is useful in clarifying how it should actually be applied:

```text
Apply only when u is explicitly first-order in eps.

Example:
(1 - 2eps)^-1
→ 1 + 2eps + O(eps²)

Do not apply recursively to arbitrary negative powers.
Do not expand UnknownFunction expressions.
```

That gives the coding agent a usable procedure without creating a general series engine. GR_PLAN3_gemini

---

### 14. Make the asymptotic-flatness check structural, not “limit-like”

Use:

```text
u = 1/r

rewrite candidate in u

evaluate boundary condition under:
u = 0
```

not:

```text
r = infinity
```

and not a new `Limit` abstraction.

Gemini's proposal is valid here and aligns with the existing v4 boundary. GR_PLAN3_gemini

---

### 15. Add explicit test-package convention

DeepSeek R2 is valid.

Pin:

```text
Default:
package <package>

Use:
package <package>_test

only when intentionally testing public API isolation.
```

That removes another implementation choice without changing architecture. GR_PLAN4_deepseek

---

## Gemini details I would **not** adopt

This is important because Gemini's review contains some technically useful but **incorrect or conflicting** specifics.

### Do not adopt its `Rational{Num int64, Den int64}` pin

The current contract does not choose `int64`, and adding that would unnecessarily narrow the representation.

### Do not adopt “sorted keys” for canonical JSON

The contract is **fixed field order**, not generic alphabetic JSON-key sorting. These are different notions.

### Do not copy its BridgePayload example literally

Gemini gives:

```json
"kind": "HYPOTHESIS"
```

and string assumptions such as:

```json
"assumptions": ["r > 0"]
```

That conflicts with the established model.

The bridge uses a **Kind** such as `Expression`; `HYPOTHESIS` is provenance. Assumptions are structured typed values, not relationship strings. The frozen spec explicitly requires typed structured assumption values. specs_v2_3

### Do not add `trig.go` merely because Gemini names one

The current v4 contract has already decided that `Sin`/`Cos` live in the bounded symbolic layer. File placement should remain whatever the final v4 tree says, rather than allowing Gemini's file name to silently reopen the tree.

---

# What DeepSeek and Qwen contribute beyond Gemini

DeepSeek's strongest contribution is **audit completeness** rather than new mathematics. Its R1–R11 residual list catches things like hash-byte semantics, test package placement, canonical-schema coverage, chart propagation, and the undefined candidate-certify procedure. GR_PLAN4_deepseek

Qwen's strongest contribution is **concrete frozen-API procedure**, especially the GR-3a path and the requirement that GR-8 use the existing `phys.Session` rather than inventing a parallel session abstraction. GR_PLAN3_qwen

---

# Recommended addition to the coding-agent prompt

I would add one final section to the handoff prompt:

```text
## Implementation Discipline — Non-Architectural Guidelines

These are implementation procedures, not architecture decisions.

1. Validate inputs before computation and never mutate inputs.
2. Every tensor operation produces a fresh output.
3. Use one global GR index ordering:
   [rho][mu][nu] for Connection,
   lexicographic row-major ordering for Tensor.
4. Treat 9 Christoffel entries as 9 independent non-zero symmetry
   classes; preserve lower-index symmetry in the full dense array.
5. Verify metric inverse componentwise against the Kronecker delta.
6. Verify Γ symmetry, Riemann antisymmetry, and Ricci symmetry where
   applicable.
7. Candidate-and-certify always means:
   construct residual → substitute candidate → Normalize → ZeroTest.
   Only ZERO certifies. UNDECIDED fails closed.
8. Primary and secondary Schwarzschild runs must not share intermediate
   derivation artifacts.
9. Recompute every input/output/step hash during replay validation.
10. Decode/re-encode canonical artifacts before accepting replay.
11. All GR-0…GR-8 structures use ChartID = "spherical-static".
12. Keep Subst strictly bounded to its declared symbol/substitution cases.
13. Keep Truncate strictly bounded to the closed-form weak-field workload.
14. Never add a general CAS, series engine, limit engine, solver,
    function registry, or hidden rewrite system.
15. Where an implementation detail is not covered by this contract,
    STOP and report the ambiguity; do not invent architecture.
```

That is the layer I think was missing: **the architecture is now locked; these rules constrain implementation technique without pretending implementation technique itself is architecture.**

The v4 plan already establishes the overarching separation — frozen `phys`, Level-1 `phys-gr`, closed GRExpr, bounded calculus, deterministic replay, and no kernel changes. GR_IMPLEMENTATION_PLAN_20261004 The added guidelines simply make that contract harder for an implementation agent to accidentally undermine.



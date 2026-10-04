I reviewed the three adversarial reviews together (DeepSeek, Gemini, Qwen) against the current **GR_IMPLEMENTATION_PLAN_20261003.md v3**. The target standard is:

> **No architectural blockers, no ambiguity for coding agent, no accidental kernel expansion.**

Overall conclusion:

**v3 architecture is correct. The remaining issues are not architectural redesigns; they are execution-contract gaps.** The reviewers agree that the core decisions are sound: `phys-gr` isolation, GR-local algebra, candidate-and-certify Schwarzschild derivation, HYPOTHESIS bridge, three-valued `ZeroTest`, and no kernel growth. DeepSeek explicitly confirms the prior major issues are resolved (curvature conventions, bridge dimensions, replay, UnknownFunction, Diff table, chart representation, bridge failure behavior, etc.). GR_PLAN3_deepseek

The required v4 changes are therefore **pin completion**, not architecture changes.

---

# Consolidated Review Outcome

## ✅ Keep unchanged (already correct)

These should NOT be reopened:

### 1. Module boundary

Keep:

```
phys
  frozen kernel

phys-gr
  Level-1 sibling module
  one-way dependency
```

No GR semantics enter `phys`.

---

### 2. GRExpr model

Keep:

```
Symbol
Rational
Add
Mul
Neg
Pow(integer)
Sin
Cos
UnknownFunction
```

No:

- generic CAS
- symbolic integration
- rational exponent calculus
- kernel tensor types
- kernel trig nodes

---

### 3. GR-3 split

Keep:

```
GR-3a:
prove kernel boundary

GR-3b:
continue with GR-local Diff
```

No modification of `phys.Differentiate`.

---

### 4. Tensor design

Keep:

- dense row-major
- rank 0–4
- lexicographic ordering
- explicit zeros
- Connection not a Tensor

---

# Required v4 patches

## BLOCKER-LEVEL FIXES

These are the items where a coding agent could legitimately make different implementations.

---

# B1 — Fix stale section references

### Problem

Definition of Done references:

```
§36 footer
```

but v3 ends at:

```
§26 + Appendix A/B/C
```

DeepSeek correctly identifies this as a self-audit failure. GR_PLAN3_deepseek

### Fix

Replace:

```
Machine checklist values verified (§36 footer)
```

with:

```
Machine checklist values verified in closing audit block:
ARCHITECTURAL_BLOCKERS=0
UNPINNED_IMPLEMENTATION_CHOICES=0
KERNEL_CHANGES_AUTHORIZED=0
```

---

# B2 — Pin canonical JSON schema

This is the biggest remaining determinism issue.

Current:

> "canonical JSON with fixed field order"

is insufficient.

A coding agent still decides:

- does every node have `"kind"`?
- what is field ordering?
- how is UnknownFunction encoded?
- how are integers encoded?

DeepSeek correctly identifies that hashes depend on canonical bytes, so unspecified JSON means unspecified replay identity. GR_PLAN3_deepseek

## Add:

Explicit schemas.

Example:

```json
{
  "kind":"symbol",
  "name":"r"
}
```

```json
{
  "kind":"pow",
  "base": {...},
  "exp":"-1"
}
```

```json
{
  "kind":"unknown_function",
  "name":"A",
  "args":["r"],
  "derivative_order":[1]
}
```

Rules:

- `kind` always first field
- fixed field ordering
- arrays only
- no maps in canonical path

---

# B3 — Pin artifact JSON schemas

Same problem exists above GRExpr.

Need schemas for:

```
GRStep
Index
IndexSlot
Tensor
Metric
Connection
Curvature
```

Currently:

```
fixed field order
```

but order is not defined.

DeepSeek flags this as separate from GRExpr because replay artifacts also depend on it. GR_PLAN3_deepseek

Add appendix:

```
Canonical Artifact Schemas
```

with exact JSON examples.

---

# B4 — Pin bridge hash payload

Current:

```
SHA256(
 canonical core.Expr +
 Kind +
 Dimension +
 Assumptions +
 Conventions
)
```

Ambiguous.

Need:

```json
{
 "expr": ...,
 "kind": ...,
 "dimension": ...,
 "assumptions": ...,
 "conventions": ...
}
```

Hash:

```
SHA256(UTF8(canonical-json))
```

DeepSeek correctly notes concatenation ambiguity otherwise. GR_PLAN3_deepseek

---

# B5 — Define ZeroTest proof rule

Current:

```
ZERO only on exact normalization proof
```

Too vague.

Pin:

```
ZeroTest(expr)=ZERO iff:

Normalize(expr)
using only the pinned rewrite rules
returns Rational(0)
```

No:

- heuristics
- numerical sampling
- assumptions-based guessing

DeepSeek correctly identifies vacuum verification depends on this. GR_PLAN3_deepseek

---

# B6 — Fix hashing reference

Change:

```
hashing = §26/Appendix C
```

to:

```
hashing = §26
```

Appendix C is self audit now.

DeepSeek identifies this as stale migration text. GR_PLAN3_deepseek

---

# HIGH PRIORITY IMPLEMENTATION CLARIFICATIONS

---

# H1 — Correct Rθθ wording

Current:

```
R_thth = 1-A-rA' (ansatz)
```

is misleading.

It is only true after:

```
B = 1/A
```

For general metric:

```
diag(-A,B,r²,r²sin²θ)
```

the expression is:

```
Rθθ =
1 - 1/B
+ rB'/(2B²)
- rA'/(2AB)
```

After substitution:

```
B=1/A

Rθθ=1-A-rA'
```

Fix wording:

```
Rθθ = 1-A-rA'
(after B=1/A)
```

DeepSeek validates this distinction. GR_PLAN3_deepseek

---

# H2 — Specify Session.Identify operands

Current:

```
Identify(k2 correspondence)
```

Too vague.

Pin exact operation:

Example:

```
Operand A:
geometric k2 expression

Operand B:
-2mu expression

Both:
Dimensionless

Justification:
Schwarzschild integration constant correspondence
```

The coding agent should not invent another identification.

---

# H3 — Add test layout

Pin:

```
Every package:

package_name/
    source.go
    source_test.go
```

Examples:

```
symbolic/
    expr.go
    expr_test.go

tensor/
    tensor.go
    tensor_test.go
```

Qwen also caught a concrete test-path problem: Pin B references a file that does not exist in the frozen tree. GR_PLAN3_qwen

Change:

```
ops/differentiate_test.go
```

to the actual existing test location.

---

# Gemini-specific findings — mostly valid, incorporate selectively

Gemini provided deeper implementation-level details. Several are valid and should be added.

---

## G1 — UnknownFunction Substitution rule ✅ VALID

Current:

```
UnknownFunction Args []Symbol
```

means:

```
A(r)
```

not:

```
A(r+k)
```

Therefore:

Pin:

```
Subst inside UnknownFunction:
only Symbol → Symbol replacement

Composite expression replacement:
UnsupportedSubstitutionError
```

This avoids silently expanding the algebra.

---

## G2 — Truncate boundary ✅ VALID

Important.

The plan forbids a general asymptotic engine.

Therefore pin:

```
Truncate is not Taylor expansion.

Supported only:
closed-form Schwarzschild weak-field expressions.
```

No:

- arbitrary series
- UnknownFunction expansion
- general asymptotics

Implement only:

```
eps = GM/(c²r)
```

weak-field extraction.

---

## G3 — r→∞ handling ✅ VALID

Do not create a Limit engine.

Replace:

```
limit r→∞
```

with:

```
structural asymptotic evaluation:

1/r → 0
```

This matches the bounded philosophy.

---

## G4 — Replay hash byte format ✅ VALID

Pin:

```
PreviousStepHash:
ASCII hex string

CanonicalJSON:
UTF-8 bytes

Payload:
ASCII(previousHashHex)
+
UTF8(canonicalJSON)
```

Otherwise two implementations can hash differently.

---

## G5 — Normalize ordering ⚠️ VALID BUT KEEP BOUNDED

Need a deterministic order.

Add:

```
Normalize pipeline:

1. numeric simplification
2. flatten Add/Mul
3. canonical sorting
4. rational collection
5. pinned trig identities
6. zero reduction
```

Avoid unlimited rewrite loops.

---

## G6 — Connection struct layout ✅ VALID

Pin:

```go
type Connection struct {
    Components [4][4][4]GRExpr
    ChartID string
}
```

meaning:

```
Components[rho][mu][nu]
=
Γ^rho_mu_nu
```

---

## G7 — Bridge extraction helper ⚠️ PARTIALLY VALID

The concept is useful.

But do NOT make `ToCoreExpr` silently extract.

Keep:

```
ToCoreExpr:
strict

any unsupported node:
error
```

Optional helper:

```
ExtractRepresentableTerms()
```

for caller-controlled decomposition.

---

# Findings NOT to add

## Reject: any architecture expansion

Do NOT add:

- kernel changes
- GR semantics in kernel
- generic CAS
- general ODE solver
- general symbolic limit engine
- general theorem prover

These contradict the whole purpose.

---

# Final v4 checklist

Before coding handoff:

```
[ ] Canonical JSON schemas pinned
[ ] GRStep schema pinned
[ ] Bridge payload schema pinned
[ ] ZeroTest exact definition pinned
[ ] Hash byte encoding pinned
[ ] Rθθ wording fixed
[ ] Session Identify operands fixed
[ ] Test layout fixed
[ ] Pin B target corrected
[ ] UnknownFunction substitution restricted
[ ] Truncate scope restricted
[ ] r→∞ replaced by structural evaluation
[ ] Connection indexing fixed
[ ] Normalize pipeline fixed
```

After these changes:

**Architectural blockers: 0**

**Implementation choices remaining: 0**

The coding agent should only face:
- mathematical difficulty,
- bugs,
- performance constraints,
- evidence classification,

not design decisions.
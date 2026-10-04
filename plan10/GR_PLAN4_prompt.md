Use the following prompt for the implementation agent. The previous review should be provided as **guidance context only**; this prompt is the **execution contract**. It intentionally does not reopen architecture. The agent must only apply the specified deterministic patches.

---

# Implementation Prompt — GR Implementation Plan v4 Hardening Patches

## Role

You are implementing the final hardening pass for `GR_IMPLEMENTATION_PLAN_20261003.md`.

The architectural review has already concluded the architecture is correct. Your task is **not** to redesign anything.

Your objective:

> Remove all remaining implementation ambiguity so that a coding agent implementing `phys-gr` has zero architectural decisions left to make.

Do not add new theory scope.
Do not expand kernel scope.
Do not modify `phys`.
Do not introduce new abstractions beyond the pinned fixes below.

---

# Scope

Modify only:

```
GR_IMPLEMENTATION_PLAN_20261003.md
```

(or the designated v4 copy if the workflow requires a new version).

Do not modify:

```
specs_v2_3.md
phys kernel
manifests
existing implementation files
```

This is a documentation contract hardening pass.

---

# Required Changes

## 1. Fix Definition of Done stale references

Replace stale:

```
§36 footer
```

references.

The document currently ends at:

```
§26
Appendix A
Appendix B
Appendix C
```

Replace with an explicit closing audit block reference.

The checklist must explicitly contain:

```
PASS0_REQUIRED=1
ARCHITECTURAL_BLOCKERS=0
UNPINNED_IMPLEMENTATION_CHOICES=0
KERNEL_CHANGES_AUTHORIZED=0
```

---

# 2. Pin GRExpr canonical JSON schema

Add a dedicated subsection under canonicalization.

The following rules are mandatory:

## General rules

Canonical JSON:

- UTF-8
- no whitespace variance
- deterministic arrays
- no maps in authoritative paths
- field ordering fixed
- every node has `"kind"` discriminator

---

## Required node schemas

Pin exact field ordering.

Example:

### Symbol

```json
{
  "kind":"symbol",
  "name":"r"
}
```

### Rational

```json
{
  "kind":"rational",
  "num":"1",
  "den":"2"
}
```

### Add

```json
{
  "kind":"add",
  "terms":[]
}
```

### Mul

```json
{
  "kind":"mul",
  "factors":[]
}
```

### Neg

```json
{
  "kind":"neg",
  "expr":{}
}
```

### Pow

```json
{
  "kind":"pow",
  "base":{},
  "exp":"-1"
}
```

### Sin

```json
{
  "kind":"sin",
  "expr":{}
}
```

### Cos

```json
{
  "kind":"cos",
  "expr":{}
}
```

### UnknownFunction

```json
{
  "kind":"unknown_function",
  "name":"A",
  "args":["r"],
  "derivative_order":[0]
}
```

Rules:

- `UnknownFunction.Args` remain symbols only.
- No expression arguments.
- No generic derivative node.

---

# 3. Pin canonical schemas for GR artifacts

Add canonical schema definitions for:

```
IndexSlot
Tensor
Metric
Connection
Curvature
GRStep
```

Each must specify:

- field order
- array representation
- integer encoding
- string encoding

The agent must not choose serialization format.

---

# 4. Pin bridge hash payload format

Replace ambiguous:

```
canonical core.Expr + Kind + Dimension + Assumptions + Conventions
```

with:

Canonical payload:

```json
{
 "expr":{},
 "kind":{},
 "dimension":{},
 "assumptions":{},
 "conventions":{}
}
```

Rules:

```
BridgeID =
gr/bridge/<SHA256(hex UTF8 canonical JSON payload)>
```

No:

- UUID
- timestamps
- counters
- random values

---

# 5. Define ZeroTest exactly

Replace:

```
exact normalization proof
```

with:

```
ZeroTest(expr)=ZERO only when:

Normalize(expr)
using only pinned normalization rules
returns Rational(0).

Otherwise:
NONZERO only on structural proof.
Otherwise:
UNDECIDED.
```

Forbidden:

- numerical approximation
- sampling
- heuristic simplification

---

# 6. Fix hashing cross references

Replace:

```
hashing = §26/Appendix C
```

with:

```
hashing = §26
```

Appendix C remains only self-audit metadata.

---

# 7. Correct Schwarzschild Ricci wording

Replace:

```
R_thth = 1 − A − rA' (ansatz)
```

with:

```
R_θθ = 1 − A − rA'
(after B = 1/A)
```

Clarify that the general ansatz expression is:

```
R_θθ =
1 - 1/B
+ rB'/(2B²)
- rA'/(2AB)
```

The simplified form is only after:

```
B = 1/A
```

---

# 8. Pin Session.Identify call

The GR-8 Identify operation must be deterministic.

Document:

Operands:

```
Operand A:
geometric integration constant expression k2

Operand B:
-2mu geometric expression
```

Purpose:

```
constant correspondence only
```

Result:

```
KindRelation
CorpusStatus NONE
```

No other Identify usage is permitted.

---

# 9. Pin test file layout

The file tree must explicitly include tests.

Rule:

Every package has adjacent tests:

Example:

```
symbolic/
    expr.go
    expr_test.go

tensor/
    tensor.go
    tensor_test.go
```

External package tests are only allowed when testing public API behavior.

---

# 10. Correct Pin B test target

Do not reference nonexistent:

```
ops/differentiate_test.go
```

Use the existing frozen test location.

The contract should state:

```
ops/negative_test.go

Differentiate(Pow(x,-1))
returns UnsupportedOperationError
```

No differentiation expansion.

---

# 11. Restrict UnknownFunction substitution

Add:

```
Substitution rules:

Allowed:
Symbol → Symbol

Example:
A(r) → A(x)

Forbidden:
Symbol → composite GRExpr

Example:
A(r+k)
A(exp(x))
```

Composite replacement returns:

```
UnrepresentableSubstitutionError
```

---

# 12. Restrict Truncate scope

Document:

`Truncate` is NOT a general series engine.

Allowed only:

```
closed-form Schwarzschild weak-field expressions
```

Workflow:

```
mu = eps*r substitution
bounded algebraic reduction
retain required orders
```

Forbidden:

- Taylor engine
- expansion of UnknownFunction
- arbitrary asymptotic solver

---

# 13. Define asymptotic flatness evaluation

Do not implement a Limit engine.

Replace:

```
r → infinity
```

with:

Structural evaluation:

```
1/r → 0
```

Only for the Schwarzschild candidate certification.

---

# 14. Pin normalization ordering

Add deterministic normalization pipeline:

```
1. numeric simplification
2. flatten Add/Mul
3. canonical child ordering
4. rational collection
5. pinned trig identities only
6. zero reduction
```

No unrestricted rewrite system.

---

# 15. Pin Connection layout

Clarify:

```
Connection components:

Components[rho][mu][nu]

represents:

Γ^rho_mu_nu
```

Connection is not Tensor.

---

# 16. Preserve architecture invariants

The final document must continue enforcing:

Allowed:

```
phys-gr -> phys
```

Forbidden:

```
phys -> phys-gr
kernel changes
new Kinds
new kernel AST nodes
GR semantics in kernel
```

---

# Validation Checklist Before Completion

Before reporting completion verify:

```
[ ] No stale section references
[ ] GRExpr JSON fully deterministic
[ ] GR artifact JSON fully deterministic
[ ] Bridge hash deterministic
[ ] ZeroTest deterministic
[ ] Session Identify deterministic
[ ] Test locations deterministic
[ ] UnknownFunction substitution bounded
[ ] Truncate bounded
[ ] No Limit engine introduced
[ ] No kernel changes required
[ ] No architectural decisions remain
```

Final response should summarize only:

1. files changed  
2. sections changed  
3. validation performed  

Do not propose further architecture changes.
# Adversarial Review of the GR Implementation Plan (Plan 11.1 Revised)

## Executive Verdict
The architectural philosophy of the plan is flawless: GR is a stress test, the kernel remains frozen, and Level 1 (theory-local) mathematics absorbs the complexity. However, the plan contains **three critical mathematical gaps** in the Level 1 machinery that will cause the coding agent to fail when attempting the actual Schwarzschild workload. The plan assumes Level 1 can just "handle" GR math, but it does not specify *how* Level 1 will represent and differentiate the specific mathematical objects that GR requires (unknown functions, trigonometry, and tensors). 

Additionally, there are necessary refinements to the PASS0 gates and the semantic meaning of Pin A.

Below is the adversarial deconstruction of the gaps and the exact amendments required before handoff.

---

## ❌ CRITICAL GAP 1: The "Unknown Function" Differentiation Trap
**The Problem:** To derive the Schwarzschild metric, one must start with a general spherically symmetric metric containing unknown functions, e.g., $g_{tt} = -A(r)$ and $g_{rr} = B(r)$. Computing the Christoffel symbols requires differentiating $A(r)$ and $B(r)$ with respect to $r$.
**The Kernel Failure:** 
- If $A(r)$ is represented as a bare `core.Symbol("A")`, the kernel's `Differentiate(A, r)` returns `0`. This is **silently wrong**.
- If $A(r)$ is represented as a `core.Call("A", [r])`, the kernel's `Differentiate` returns `UnsupportedOperationError` (since `Call` is restricted to `lorentz_factor`).
**The Fix:** The plan must explicitly mandate that Level 1 (`phys-gr/symbolic/`) defines a representation for **unknown functions of a coordinate** (e.g., an `UnknownFunction` AST node) and implements the chain rule for them in the GR-local differentiator. The differentiator must return a new `UnknownFunction` representing the derivative (e.g., $A'(r)$).

## ❌ CRITICAL GAP 2: Trigonometric Derivatives in the Angular Sector
**The Problem:** The Schwarzschild metric in spherical coordinates contains the angular sector $r^2 d\theta^2 + r^2 \sin^2\theta d\phi^2$. Computing the Ricci tensor requires differentiating $\sin(\theta)$ with respect to $\theta$.
**The Kernel Failure:** The frozen kernel has no trigonometric nodes.
**The Fix:** The plan explicitly mentions `trig.go` in passing, but it must be elevated to a **mandatory Level 1 component**. `phys-gr/symbolic/trig.go` must define `Sin` and `Cos` AST nodes and their derivative rules ($d(\sin x)/dx = \cos x$, $d(\cos x)/dx = -\sin x$). Without this, the angular components of the curvature tensors cannot be computed.

## ❌ CRITICAL GAP 3: Tensor Data Structures are "Emergent" but Undefined
**The Problem:** The plan states: *"Index/tensor algebra: variance, contraction, raising/lowering... emergent in GR-2 inside the workload."*
**The Risk:** If tensor algebra is just "emergent in the workload," the coding agent will write ad-hoc, untestable spaghetti code inside the test files. 
**The Fix:** The plan must mandate a concrete Level 1 data structure for tensors before the workload begins. For example:
```go
// phys-gr/tensor/tensor.go
type Variance bool // Covariant or Contravariant
type Tensor struct {
    Rank       int
    Variance   []Variance
    Components map[string]core.Expr // e.g., "tt" -> expr, "tr" -> expr
}
func (t Tensor) Contract(other Tensor, idx1, idx2 int) Tensor { ... }
```
Level 1 must own the tensor engine; the workload only supplies the components.

---

## ⚠️ REFINEMENT 1: Pin A Semantics and `Identify`
**The Issue:** The review correctly fixes the Pin A fixture to `Identify(RestMass, E/c²)`. However, the assertion `result.Kind() != KindRestMass` is trivially true and slightly misses the point.
**The Reality:** `Identify` *always* returns a new object of `KindRelation`. It does not retype or mutate its operands. 
**The Fix:** Pin A must explicitly assert **operand immutability** and **result properties**:
1. `result.Kind() == KindRelation`
2. `result.Expr().RelationOperator() == RelationEq`
3. `result.Provenance().Status() == IDENTIFIED`
4. `result.CorpusStatus() == NONE`
5. `result.Dimension() == M`
6. The original `RestMass` object remains `KindRestMass` (proving `Identify` did not mutate the input).

## ⚠️ REFINEMENT 2: PASS0 is a Regression Gate, Not a Development Phase
**The Issue:** The plan describes PASS0 as implementing Pins A, B, and C. But the review notes these are "test pin only, zero production diff expected."
**The Fix:** Rename PASS0 to **"PASS0 — Baseline Regression Verification"**. Its sole purpose is to prove that the *already implemented* frozen kernel behaves exactly as the spec demands. If Pin B fails (e.g., the kernel accidentally allows negative-integer differentiation), the baseline is corrupted, and the GR stress test is invalid. PASS0 adds no new features; it only adds regression tests to `core/object_test.go` and `ops/negative_test.go`.

## ⚠️ REFINEMENT 3: Derivation vs. Verification of Schwarzschild
**The Issue:** The workload says "Schwarzschild exterior". Does the plan *verify* the known Schwarzschild metric, or does it *derive* it from general spherical symmetry?
**The Fix:** To truly stress the Level 1 machinery, the workload must be a **derivation**. 
- **Verification** (plugging in $1 - 2M/r$ and checking $R_{\mu\nu} = 0$) only tests rational algebra.
- **Derivation** (starting with $A(r)$ and $B(r)$, computing $R_{\mu\nu} = 0$, and solving for $A$ and $B$) forces the "unknown function" differentiation trap and the algebraic solver to actually work. The plan must explicitly mandate the derivation path.

---

## 🛠️ REQUIRED AMENDMENTS TO THE PLAN

Inject the following section into the GR Implementation Plan before handoff:

```markdown
### Addendum: Level 1 Mathematical Machinery (Mandatory Pre-Requisites)
Before executing the GR workload (GR-0 to GR-8), the Level 1 module 
(`phys-gr`) MUST implement the following mathematical machinery. 
These are not kernel extensions; they are Level 1 userland structures 
operating over public `core.Expr`.

1. Unknown Function Representation (`phys-gr/symbolic/func.go`):
   - Define an AST node for unknown functions of a coordinate (e.g., A(r)).
   - The GR-local differentiator MUST implement the chain rule for these 
     nodes, returning a new node representing the derivative (e.g., A'(r)).
   - This prevents the "silently-wrong" trap where bare symbols differentiate to 0.

2. Trigonometric Functions (`phys-gr/symbolic/trig.go`):
   - Define AST nodes for Sin and Cos.
   - The GR-local differentiator MUST implement their derivatives 
     (d(sin x) = cos x, d(cos x) = -sin x).
   - Required for the angular sector of the Schwarzschild metric.

3. Tensor Engine (`phys-gr/tensor/tensor.go`):
   - Define a concrete data structure for tensors (Rank, Variance, Components map).
   - Implement contraction, raising/lowering (using the metric), and 
     symmetry bookkeeping.
   - The GR workload will use this engine; it will not use ad-hoc maps in tests.

4. Workload Mandate:
   - The Schwarzschild workload (GR-8) MUST be a DERIVATION from general 
     spherical symmetry (starting with unknown functions A(r) and B(r)), 
     not a mere verification of the known metric. This forces the Level 1 
     machinery to solve the actual differential equations.
```

### Final Handoff Status
With the addition of the **Level 1 Mathematical Machinery** addendum and the clarifications to **Pin A** and **PASS0**, the plan is fully bulletproof. It correctly isolates the kernel, defines the exact boundaries of Level 1, and anticipates the specific mathematical bottlenecks that will arise when computing curvature tensors. 

**You are cleared to hand off to the coding agent.**
# Adversarial Review of `GR_IMPLEMENTATION_PLAN_20261003.md` (v3)

## Executive Verdict: CLEARED FOR HANDOFF (with 6 surgical patches)

This is a masterclass in execution contracting. The plan successfully translates the dense, adversarial philosophy of the Reality-First Constitution and the v2.3 specs into a deterministic, Go-idiomatic engineering blueprint for a hostile workload. The separation of the 14-step "candidate-and-certify" Schwarzschild script from a general ODE solver is a brilliant architectural constraint that perfectly respects the kernel's bounded symbolic engine.

However, a line-by-line audit against the frozen 42-file tree and the v2.3 API surface reveals **6 micro-ambiguities** that will cause a coding agent to hesitate, invent unauthorized files, or misinterpret the bridge mechanics. 

Below is the exact audit and the required patches to make the document 100% bulletproof.

---

### ❌ Patch 1: The Phantom File in PASS0 Pin B
**The Error:** Section 16 (Test Strategy) states: *"Pin B... target `ops/differentiate_test.go` (membership-checked; halt rather than create)."*
**The Violation:** `ops/differentiate_test.go` **does not exist** in the frozen 42-file tree. The tree only contains `ops/operations_test.go` and `ops/negative_test.go`. If the agent follows this, it will either create a forbidden file or halt in confusion.
**The Fix:** Change the target to the correct existing file.
> **Inject into §16:** "Pin B (`ops/negative_test.go`): `Differentiate(Pow(x, -1))` → `UnsupportedOperationError`. Do not expand differentiation. Target the existing `TestDifferentiateUnsupportedForms` or add a specific `TestGR3aNegativePowerRejection` within `ops/negative_test.go`."

### ❌ Patch 2: The `Relation` / `BranchSet` Bridge Gap
**The Error:** Section 21 (Bridge Contract) defines the mirror subset as: *"Symbol/Rational/Add/Mul/Neg/Pow(int) solely."* But Section 23 (GR-8) requires using `ops.Compare` and `ops.SelectBranch` in the session segment. These operations require `Relation` and `BranchSet` nodes.
**The Violation:** The agent will realize `Relation` is not in the `GRExpr` AST mirror subset and won't know how to construct it for the kernel.
**The Fix:** Clarify that complex structural nodes are built via `core` constructors, not the `GRExpr` AST.
> **Inject into §21:** "Mirror subset for scalar algebra: `Symbol/Rational/Add/Mul/Neg/Pow(int)` solely. Complex structural nodes (`Relation`, `BranchSet`) are **not** part of the `GRExpr` AST. They are constructed directly via `core.NewRelation` and `core.NewBranchSet` using bridged scalar `core.Expr` operands."

### ❌ Patch 3: The `weakfield.go` Schrodinger's File
**The Error:** Section 8 (File Plan) states: *"symbolic/weakfield.go MAY fold into limit/weakfield.go; plan records canonical home as limit/weakfield.go with symbolic/weakfield.go as thin forwarder if retained."*
**The Violation:** Coding agents hate "MAY". It will likely create both files to be safe, violating the clean architecture, or guess wrong.
**The Fix:** Mandate a single canonical location.
> **Replace in §8:** "Canonical home for weak-field reduction is `limit/weakfield.go`. Do not create `symbolic/weakfield.go`."

### ❌ Patch 4: The GR-3a Bridge Execution Sequence
**The Error:** Section 7 states: *"GR-3a uses `ops.Differentiate` directly, never `Session.Step`."*
**The Violation:** The agent knows `ops.Differentiate` takes `core.Object`, but the GR workload operates on `GRExpr`. The agent needs the exact mechanical sequence to bridge, execute, and reintegrate.
**The Fix:** Provide the exact 4-step code flow for GR-3a.
> **Inject into §7:** "GR-3a execution sequence: 
> 1. `coreExpr, _ := bridge.ToCoreExpr(grExpr)` 
> 2. `coreObj, _ := bridge.ToCoreObject(coreExpr, dim, assumptions, conventions)` 
> 3. `resultObj, err := ops.Differentiate(coreObj, wrtObj)` (Expect `UnsupportedOperationError`) 
> 4. Record the failure as `SPEC-INTENDED-BOUND` in the GR replay trace."

### ❌ Patch 5: The `Truncate` Scope Ambiguity
**The Error:** Section 23 mentions `Truncate(expr, eps, 1)` for the Newtonian limit.
**The Violation:** The agent might try to implement a general symbolic Taylor expansion engine for `UnknownFunction`s, which is explicitly forbidden (No general series expansion).
**The Fix:** Constrain `Truncate` to the explicit closed-form metric.
> **Inject into §23:** "Truncation is performed **only** on the explicit closed-form Schwarzschild metric obtained from GR-8 PRIMARY/SECONDARY (where `A` and `B` are explicit rational functions of `r` and `mu`), not on the generic `UnknownFunction` ansatz. It is implemented as a bounded algebraic substitution (e.g., `mu = eps * r`) followed by `Simplify` to drop `O(eps^2)` terms. No general series engine is built."

### ❌ Patch 6: The Session API Invocation
**The Error:** Section 23 mentions the "GR-8 session segment" but doesn't specify how the agent instantiates the session.
**The Violation:** The agent might invent a custom GR session wrapper.
**The Fix:** Explicitly bind it to the frozen `phys` session API.
> **Inject into §23:** "The GR-8 session segment uses the standard frozen `phys` session API exactly as specified in `specs_v2_3.md` §16: `s := session.New()`, `s.Draft(...)`, `s.Step(...)` for `Compare`, `s.Identify(...)` for the constant correspondence, `s.Commit()`, `s.Conclude()`. The `DraftMetadata` will contain the GR-8 hypothesis."

---

## Final Architectural Validation

With those 6 patches applied, the plan achieves **zero architectural blockers**. 

### Why this plan is a masterpiece of constraint:
1. **The "Candidate-and-Certify" Pattern:** By forcing the agent to propose an algebraic identity (e.g., `(AB)' = 0`) and use `ZeroTest` to certify it, you completely bypass the need for a symbolic ODE solver while still rigorously deriving the Schwarzschild metric. This is a profound alignment with the kernel's bounded symbolic engine.
2. **The HYPOTHESIS Bridge:** Using `hypothesis.NewCandidateConcept` to temporarily wrap GR scalars into `core.Object` is the exact correct mechanism. It allows the GR package to leverage the kernel's dimensional analysis, MRC, and `Simplify` without polluting the trusted corpus or inventing new kernel Kinds.
3. **The 3-Valued `ZeroTest`:** Forcing `UNDECIDED` instead of guessing `ZERO` is the ultimate defense against the "Silent-Wrongness" adversarial suite. It ensures the system fails loud rather than failing wrong.
4. **Strict Module Isolation:** The `phys-gr` sibling module pattern perfectly enforces the Growth Gate. The kernel cannot accidentally acquire GR semantics because the Go compiler physically prevents the import.

### Handoff Status: CLEARED
Apply the 6 patches above. Once applied, hand the document to the coding agent. It has exactly what it needs to execute the most rigorous stress test the Physics Compiler MVP has ever faced, without making a single unauthorized architectural decision.
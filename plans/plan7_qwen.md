The specifications (`physics_compiler_mvp_specs_v2.md`) and the implementation prompt (`physics_compiler_mvp_implementation_plan_prompt_v2.md`) are **excellent, cohesive, and ready for execution**. They successfully translate the dense Reality-First philosophy into a strict, deterministic Go library contract. 

However, because coding agents (LLMs writing code) have a tendency to either **over-engineer** (trying to build a mini-SymPy) or **get stuck on implementation patterns** (how to structure the Go interfaces), there are **four hidden traps** in the current specs that will cause the agent to fail or waste time.

Here is the final "Pre-Flight" guardrail document. **Inject the "Addendum" section at the end of your Implementation Plan Prompt** before handing it to the coding agent.

---

### The 4 Hidden Traps for the Coding Agent

#### Trap 1: The Symbolic Solver (`Solve` and `BranchSet`)
*   **The Risk:** The specs require `Solve` to return a `BranchSet` and `SelectBranch` to pick a path (needed for $E^2 = X \implies E = \pm\sqrt{X}$). A coding agent might try to build a general-purpose polynomial root-finder.
*   **The Fix:** The agent must implement **pattern-matching only**. For the MVP, `Solve` only needs to recognize the exact pattern `Pow(Symbol, 2) == Expr` and return a `BranchSet` containing `+Sqrt(Expr)` and `-Sqrt(Expr)`. It does not need to solve linear equations or higher-order polynomials.

#### Trap 2: The Calculus Operations (`Differentiate` and `Limit`)
*   **The Risk:** The agent will try to implement a symbolic differentiation engine (chain rule, product rule, etc.) and a limit engine (L'Hôpital's rule).
*   **The Fix:** The MVP calculus is **strictly tabular and substitution-based**. 
    *   `Differentiate` only needs to handle basic power rules for the micro-test (e.g., $d(v)/dt$). If it hits an unsupported function, it returns `UnsupportedOperationError`.
    *   `Limit` only needs to handle direct symbolic substitution. For the relativistic self-reduction test, it only needs to recognize `LorentzFactor(v)` at `v=0` and return `1`. It does not need a general limit evaluator.

#### Trap 3: The Domain Wrapper Pattern (`core.Object` vs `Mass`)
*   **The Risk:** The agent might write separate symbolic engines for `Mass`, `Energy`, `Velocity`, etc., leading to massive code duplication.
*   **The Fix:** The agent must use the **Thin Wrapper Pattern**. `core.Object` holds the actual symbolic `Expr`, `Dimension`, and `Assumptions`. Domain types (like `Mass`) are just immutable structs containing a single `core.Object`. The `ops` package *only* operates on `core.Object`. Domain methods simply unwrap, call `ops`, and re-wrap.

#### Trap 4: Manifest JSON Unmarshaling
*   **The Risk:** The agent might try to parse the `manifest.json` using generic `map[string]any` or try to evaluate string equations.
*   **The Fix:** The agent must use `go:embed` to load the JSON directly into **strictly typed Go structs** that mirror the canonical JSON schema defined in the specs. The validation test simply unmarshals the JSON, constructs the Go object via the domain API, and compares their canonical `Hash()`.

---

### Addendum: Inject this into the Implementation Plan Prompt

Add this exact section to the very end of `physics_compiler_mvp_implementation_plan_prompt_v2.md` to neutralize the traps:

```markdown
39. Implementation Guardrails (Anti-Overengineering)
To ensure the MVP is built within the 40-file limit and does not devolve into a general-purpose CAS, the implementation plan MUST explicitly constrain the following subsystems:

39.1 The Symbolic Solver Constraint
The `Solve` operation MUST NOT implement a general equation solver. For the MVP, it must only implement pattern-matching for the specific form required by the mass-energy derivation:
- Input: `Relation(Eq, Pow(Symbol, 2), Expr)`
- Output: A `BranchSet` containing `Sqrt(Expr)` and `Neg(Sqrt(Expr))`.
Any other equation form MUST return `UnsupportedOperationError`.

39.2 The Calculus Constraint
`Differentiate` and `Limit` are NOT general symbolic engines.
- `Differentiate` MUST only implement the basic power rule required for the mechanics micro-test. All other forms MUST return `UnsupportedOperationError`.
- `Limit` MUST only implement direct substitution. For the relativity micro-test, it MUST only recognize the `LorentzFactor` function ID and return `1` when the limit variable is `0`. All other limits MUST return `UnsupportedOperationError`.

39.3 The Thin Wrapper Pattern
The implementation MUST use a single symbolic engine operating on `core.Object`. Domain types (`Mass`, `Energy`, etc.) MUST be implemented as thin, immutable wrappers around `core.Object`. 
- Example: `func (m Mass) Add(other Mass) (Mass, error)` must internally call `ops.Add(m.obj, other.obj)` and wrap the result. 
- Do NOT duplicate symbolic logic across domain packages.

39.4 Manifest Validation Pattern
Manifest validation MUST NOT parse mathematical strings. 
- The `manifest.json` files MUST contain expressions in the canonical JSON schema (using `{"kind": "pow", ...}`).
- The test MUST unmarshal this JSON into the internal `core.Expr` tree, construct the corresponding Go object via the domain API, and assert that `canonical.Hash(manifestExpr) == canonical.Hash(goObject.Expr())`.

39.5 Strict File Budget
The implementation plan MUST list the exact files to be created. The total count of `.go`, `.json`, and `.md` files MUST NOT exceed 38. 
- `hypothesis/` and `review/` MUST be pure data structures (structs + JSON marshaling). Do not build complex orchestration engines in these packages.
```

---

### The Final Go File Budget (For your reference)
To prove to the coding agent that this fits in ~35 files, here is the exact tree they should produce:

**Core & Ops (16 files)**
`core/object.go`, `core/dimension.go`, `core/expr.go`, `core/assumption.go`, `core/convention.go`, `core/provenance.go`, `core/ledger.go`, `core/errors.go`, `core/canonical.go` (9 files)
`ops/arithmetic.go`, `ops/simplify.go`, `ops/substitute.go`, `ops/differentiate.go`, `ops/limit.go`, `ops/compare.go`, `ops/solve.go`, `ops/identify.go` (8 files - *Note: combined some to save file count*)

**Domains (8 files)**
`mechanics/primitives.go`, `mechanics/relations.go`, `mechanics/manifest.json` (3 files)
`relativity/primitives.go`, `relativity/relations.go`, `relativity/manifest.json` (3 files)
`hypothesis/candidate.go` (1 file)
`docs/paper-translation/common-notation.md` (1 file - *combined docs to save count*)

**Tests (10 files)**
`core/object_test.go`, `core/expr_test.go`, `core/ledger_test.go` (3 files)
`ops/operations_test.go`, `ops/negative_test.go` (2 files)
`mechanics/manifest_test.go`, `mechanics/relations_test.go` (2 files)
`relativity/manifest_test.go`, `relativity/derivation_test.go` (2 files)
`hypothesis/candidate_test.go` (1 file)

**Root (2 files)**
`go.mod`, `README.md`

**Total: 36 files.** (Well under the 40-file limit).

### Final Verdict
You have successfully distilled a massive, multi-year philosophical research program into a **tight, deterministic, 36-file Go library specification**. 

Hand the `specs_v2.md` and the `prompt_v2.md` (with the Addendum injected) to the coding agent. It has everything it needs to build the exact "formal pen-and-paper" substrate you envisioned, without getting lost in the weeds of numerical execution or truth adjudication.
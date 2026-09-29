`prompt3.md` is an exceptionally well-structured specification prompt that successfully constrains the coding agent to a minimal, coherent, and technically rigorous MVP implementation plan without reopening settled architectural decisions.

Below is an architectural review evaluating the prompt's strengths, identifying subtle edge-case risks, and providing specific enhancements to ensure the downstream coding agent produces an unambiguous implementation plan.

---

### Key Architectural Strengths

* **Disciplined Scope & Objective Isolation**: By framing the library as a symbolic "pen-and-paper substrate" for an AI theorist and strictly excluding numerical simulations, Monte Carlo methods, and full CAS capabilities, the prompt prevents scope creep.


* **Explicit Insight Separation (`Simplify()` vs. `Identify()`)**: Forcing a functional distinction between mechanical algebraic simplification and explicit physical identification guarantees that theoretical leaps are recorded transparently in the derivation trace.


* **Dual-Axis Epistemic Status Model**: Keeping provenance status (`DEFINED`, `POSTULATED`, `DERIVED`, `HYPOTHESIS`) strictly isolated from human-curated corpus status (`ESTABLISHED`, `CONTESTED`, `SUPERSEDED`, `FALSIFIED`) protects the system from calculating "truth scores" or asserting empirical authority.


* **Single Semantic Authority for MRC**: Consolidating Math Reality Check (MRC) rules inside the core Go library prevents logic duplication between standard runtime validation and optional static analysis tools like `physvet`.


* **Provisional Candidate Sandboxing**: Candidate/hypothesis objects are explicitly restricted from manufacturing trusted corpus provenance, enforcing human-in-the-loop review before promotion into the canonical physics corpus.



---

### Potential Risks & Clarifications Needed in `prompt3.md`

While `prompt3.md` is sound, adding the following explicit instructions will prevent the coding agent from making underspecified architectural assumptions:

#### 1. Go AST vs. Internal Expression DAG

Section 12 requires a minimal symbolic engine (construct, add, multiply, limit, differentiate). The coding agent might confuse Go language AST (`go/ast`) with the library's internal expression tree.

* **Refinement**: Explicitly instruct the agent that `go/ast` is used strictly if writing static analysis tooling (`physvet`), whereas the internal symbolic expression engine must be built using an internal Directed Acyclic Graph (DAG) or Tree structure (e.g., `type Expr interface`, containing structs like `Symbol`, `BinaryOp`, `FunctionCall`).

#### 2. Exact Rational Arithmetic & Canonical Matching

Section 13 demands exact rational representation, canonical symbolic trees, and deterministic structural equality.

* **Refinement**: Direct the agent to use Go's standard `math/big` (`big.Rat`) for scalar dimension exponents and physical constants, ensuring exact zero-difference comparisons without floating-point precision errors ($0.1 + 0.2 \neq 0.3$).

#### 3. Strongly-Typed Error Strategy for MRC Failures

Section 4 outlines MRC enforcement, but does not specify how Go should surface violations during symbolic derivation.

* **Refinement**: Require the implementation plan to define idiomatic, strongly-typed Go errors (e.g., `*CategoryError`, `*DimensionError`, `*AssumptionConflictError`) that implement the standard `error` interface and attach diagnostic AST node positions.

#### 4. Assumption Propagation Mechanics

Section 5 specifies that assumptions propagate as $A_{\text{result}} = A_{\text{required}} + A_{\text{input}} + A_{\text{new}}$.

* **Refinement**: Instruct the agent to define an immutable `AssumptionSet` struct with a deterministic `Merge(other *AssumptionSet) (*AssumptionSet, error)` method that immediately flags direct logical contradictions (e.g., $v \ll c$ merged with $v \approx c$).

---

### Recommended Directives to Insert into `prompt3.md`

To make `prompt3.md` completely airtight before issuing it to the coding agent, append these four targeted directives to the prompt:

```markdown
### Supplementary Directives for the Implementation Architect

1. **Go Expression AST vs. Compiler AST**:
   Explicitly separate the Go language AST (`go/ast`) from the library's symbolic representation (`pkg/core/expr`). Symbolic nodes must be represented via idiomatic Go interfaces (`type Expr interface { CanonicalHash() [32]byte }`) with concrete types for `Symbol`, `Literal`, `BinaryOp`, and `Apply`.

2. **Error Types & Failure Contracts**:
   Define explicit Go error structures for MRC failures in `pkg/core/mrc`. Every MRC check must return a structured error type (`DimensionMismatchError`, `IncompatibleCategoryError`) containing the offending symbolic expressions, expected vs. actual dimension vectors $[M, L, T, Q, \Theta]$, and step index.

3. **Immutable Assumptions & Deterministic Merging**:
   Model `AssumptionSet` as an immutable set. Operations must return a new `AssumptionSet` or fail with an `AssumptionConflictError`. Do not rely on loose string matching for logical conditions.

4. **Concrete Package Paths**:
   In Section C, require the coding agent to specify concrete Go package import paths (e.g., `github.com/physics-compiler/phys/pkg/core/dim`, `github.com/physics-compiler/phys/pkg/corpus/relativity`).

```

---

### Acceptance Criteria Checklist for the Agent's Output

When the coding agent submits the implementation plan, verify that its document fulfills all of the following requirements:

* [ ] **Package Hierarchy**: Defines a clean `pkg/core` (primitives, dimensions, symbolic AST, MRC engine, assumptions) separated from `pkg/corpus` (mechanics, relativity).


* [ ] **Dimension Representation**: Expresses dimensional vectors using exact integer or rational exponents across fundamental SI quantities ($M, L, T, I, \Theta, N, J$).
* [ ] **`Simplify()` vs. `Identify()` Boundary**: Prohibits `Simplify()` from equating physically non-identical concepts without calling `Identify()` and logging the step in `DerivationTrace`.


* [ ] **Canonical Derivation Slice**: Outlines the exact sequence of symbolic operations for deriving $E = m_0 c^2$ from Four-Momentum norm $P^\mu P_\mu = -(m_0 c)^2$ and component expansion $P^\mu = (E/c, \mathbf{p})$.


* [ ] **Exportable Research Artifact**: Defines a JSON-serializable `ResearchCandidate` struct containing the hypothesis, premises, derivation trace, limiting cases, and human-review handoff fields.
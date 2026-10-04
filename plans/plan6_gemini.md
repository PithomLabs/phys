The specification (`PHYSICS_COMPILER_MVP_SPECS.md`) and implementation prompt (`PHYSICS_COMPILER_MVP_IMPLEMENTATION_PLAN_PROMPT.md`) form a tight, normative contract for building a deterministic, symbolic "pen-and-paper substrate" in Go.

The documents successfully constrain the system scope—strictly excluding numerical simulations, full Computer Algebra System (CAS) engines, theorem provers, and automated physical truth adjudication.

---

## Architectural Alignment & Core Strengths

1. **Strict Authority Hierarchy**: The prompt explicitly declares the specs document as normative, ensuring that any downstream coding agent cannot alter the core architecture, substitute languages, or broaden scope.
2. **Carrier-Safe Physical Construction**: Enforces distinct Go nominal types for physical quantities (e.g., `Mass`, `Velocity`, `FourMomentum`) with unexported state, preventing AI agents from instantiating unvalidated physical objects via open struct literals.
3. **Exact Rational Dimensionality**: Base SI dimensions ($M, L, T, I, \Theta, N, J$) use `math/big.Rat` exponents, preventing floating-point rounding errors during dimensional checks ($L^1 T^{-1} \cdot T^1 = L^1$).
4. **Epistemic & Provenance Isolation**: Keeps provenance status (`DEFINED`, `POSTULATED`, `DERIVED`, `IDENTIFIED`, `APPROXIMATED`, `HYPOTHESIS`) separate from human-curated corpus status (`ESTABLISHED`, `CONTESTED`, `SUPERSEDED`, `FALSIFIED`), eliminating "truth scores" or computed empirical authority.
5. **Functional Separation of `Simplify` vs. `Identify**`: Mechanical algebraic reduction (`Simplify`) is strictly barred from equating physically non-identical concepts. Physical leaps require an explicit `Identify(a, b, justification)` call logged in the session ledger with reviewable human/AI justification.
6. **Hash-Chained Derivation Ledger & Replay Validation**: Derivation sessions produce SHA-256 hash-chained steps. Tampering with intermediate expressions or steps invalidates `Validate()`, guaranteeing cryptographic auditability.

---

## Edge Cases & Architectural Risks

While the specifications are comprehensive, the following technical edge cases require explicit directives in the implementation prompt to prevent the downstream coding agent from making underspecified design decisions.

### 1. Expression AST Representation vs. Go Compiler AST

* **Risk**: The prompt mentions Go AST in the context of MRC, which could lead the coding agent to confuse Go standard library AST (`go/ast`) with the physics library's internal expression engine (`core.Expr`).
* **Resolution**: Clarify in the prompt that `go/ast` is strictly for static analysis tools (deferred to v0.5+), whereas `core.Expr` must be implemented via an internal immutable Directed Acyclic Graph (DAG) or composite tree interface with closed node types (`Symbol`, `Rational`, `Add`, `Mul`, `Neg`, `Pow`, `Sqrt`, `Relation`, `Derivative`, `Limit`, `Function`).

### 2. Canonicalization Determinism & Map Iteration Non-Determinism

* **Risk**: Canonical sorting during expression simplification (`Simplify`) or SHA-256 hashing (`Hash`) can become non-deterministic if Go `map` iteration is used without explicit sorting.
* **Resolution**: Require canonicalization algorithms to sort expression terms using deterministic byte-level comparisons of node types and canonical hashes (e.g., lexicographical sorting of commutative operands in `Add` and `Mul`).

### 3. Assumption Conflict Mechanics

* **Risk**: Section 7.3 requires `AssumptionConflictError` when two different values exist for the same `(Kind, Key)`. Without explicit key formatting rules, string mismatches like `velocity_relation` vs `vel_rel` could bypass merge validation.
* **Resolution**: Mandate a standardized `Key` naming standard (e.g., lowercase snake_case namespaces like `regime.velocity`, `convention.metric_signature`) and treat `AssumptionSet` as an immutable set with deterministic union logic.

### 4. Candidate Containment & Shielding against Corpus Leakage

* **Risk**: An AI agent might attempt to pass a provisional candidate object (`HYPOTHESIS` status) into a standard domain relation to mint a trusted `DERIVED` result.
* **Resolution**: Enforce that any expression or object built with a `HYPOTHESIS` input automatically poisons the downstream result, tagging it strictly as `HYPOTHESIS`. Prohibit session routines from sealing a non-hypothesis candidate into a trusted corpus manifest.

---

## Recommended Directives to Append to `prompt.md`

Add the following targeted directives to `physics_compiler_mvp_implementation_plan_prompt.md` to ensure the resulting implementation plan is fully unambiguous:

```markdown
### Targeted Directives for the Implementation Plan

1. **Closed Expression Engine (`core.Expr`)**:
   - Represent symbolic expressions via an internal `type exprNode interface` with unexported methods to enforce closed implementation within `core`.
   - Use `math/big.Rat` for exact scalar representation inside `Rational` AST nodes and dimension vector exponents.
   - Prohibit reflection-based or generic `interface{}` evaluation nodes.

2. **Canonical Hashing & Determinism Contract**:
   - Define canonical serialization as a deterministic byte sequence representing the normalized expression tree.
   - Commutative nodes (`Add`, `Mul`) MUST sort child nodes lexicographically by their canonical SHA-256 hash before serialization.
   - Ensure `Hash()` never incorporates pointer addresses, system timestamps, map iteration output, or random UUIDs.

3. **Strongly-Typed Error Taxonomy**:
   - Detail concrete Go error structs implementing the standard `error` interface:
     - `DimensionMismatchError`
     - `CategoryMismatchError`
     - `AssumptionConflictError`
     - `ConventionConflictError`
     - `IdentifyError`
     - `ProvenanceError`
     - `CandidateContainmentError`
   - All errors must expose inspectable fields (e.g., expected vs. actual dimensions) accessible via `errors.As`.

4. **Concrete Package Export Boundaries**:
   - Ensure `mechanics` and `relativity` expose domain-specific constructors that wrap core carriers.
   - Unexported internal fields must prevent callers outside `core` from bypassing constructor-level Math Reality Checks (MRC).

```

---

## Package Mapping & Core Implementation Structure

The implementation plan produced by the agent should map directly to the package layout below:

| Package | Primary Responsibilities | Key Types & Structs |
| --- | --- | --- |
| `core/` | Base physical carriers, exact dimensions, symbolic AST, assumptions, conventions, provenance, hash-chained derivation ledger, and candidate objects. | `Dimension`, `Expr`, `AssumptionSet`, `ConventionSet`, `Provenance`, `Session`, `Derivation`, `Step`, `ResearchCandidate`, `FalsificationCondition` |
| `ops/` | Symbolic transformation engines and Math Reality Check (MRC) rule enforcement. | `Add()`, `Subtract()`, `Multiply()`, `Divide()`, `Pow()`, `Sqrt()`, `Simplify()`, `Substitute()`, `Differentiate()`, `Limit()`, `Compare()`, `Solve()`, `Identify()` |
| `mechanics/` | Classical mechanics primitives, $F=ma$, $p=mv$, $K=\frac{1}{2}mv^2$, and embedded JSON manifest validation. | `Mass`, `Time`, `Position`, `Velocity`, `Acceleration`, `Force`, `Momentum`, `Energy`, `NewtonSecondLaw()` |
| `relativity/` | Special relativity primitives, spacetime metric, four-momentum, $E^2 = (pc)^2 + (m c^2)^2$, and $E=mc^2$ rest-frame derivation. | `Spacetime`, `MinkowskiMetric`, `RestMass`, `ThreeMomentum`, `FourMomentum`, `SpeedOfLight`, `LorentzFactor`, `EnergyMomentumRelation()` |
| `hypothesis/` | Provisional candidate object sandboxing and research candidate formulation. | `CandidateConcept`, `ProvisionalRelation`, `ResearchCandidate` |
| `docs/paper-translation/` | AI corpus documentation mapping mathematical physics notation to code carriers. | `common-notation.md`, `framework-mapping.md`, `ambiguity-resolution.md` |

---

## Verification & Acceptance Checklist

When evaluating the output generated by the downstream implementation plan agent, verify that the plan satisfies all mandatory criteria:

* [ ] **Exact Go Module Layout**: Strictly adheres to the 6 allowed directories (`core`, `ops`, `mechanics`, `relativity`, `hypothesis`, `docs`) without creating stub directories for deferred domains (`qm`, `qft`, `statmech`).
* [ ] **Carrier-Safe Primitives**: Confirms distinct nominal types for all physical quantities with unexported internal fields.
* [ ] **Exact Dimensions**: Uses `math/big.Rat` across all 7 SI base quantities ($M, L, T, I, \Theta, N, J$).
* [ ] **No Integrated CAS/Theorem Prover**: Confirms symbolic operations are limited strictly to the 13 specified methods. Excludes `Integrate()`.
* [ ] **Explicit `Identify()` Boundary**: Confirms `Simplify()` is strictly algebraic and cannot generate `IDENTIFIED` provenance or rename physical concepts.
* [ ] **Hash-Chained Ledger Replay**: Specifies cryptographic SHA-256 step chaining and step-by-step re-execution in `Validate()`.
* [ ] **Manifest Embedding**: Exposes `Manifest()` in `mechanics` and `relativity` via `embed.FS` parsing `manifest.json`.
* [ ] **Mandatory Test Suite**: Outlines test cases for all specified micro-demos (**Test A** through **Test K**), including explicit negative tests for dimension mismatch, category mismatch, assumption conflict, convention conflict, and corrupted ledger replay.
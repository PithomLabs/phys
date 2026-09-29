Both `physics_compiler_mvp_implementation_plan_prompt_v2.md` and `physics_compiler_mvp_specs_v2.md` are architecturally sound, mathematically precise, and ready for dispatch to the implementation architect agent. The specifications establish an explicit "pen-and-paper substrate" boundary, preventing scope creep into numerical simulations, Computer Algebra Systems (CAS), or automated physical truth adjudication.

---

### Key Architectural Strengths

* **Two-Layer Carrier Pattern**: Domain structs (`Mass`, `Velocity`, `FourMomentum`) act as thin, immutable nominal wrappers around `core.Object` via the required `CoreObject()` accessor. Generic operations in `ops` consume `core.Object`, enforcing type-safe construction boundaries without sacrificing generic symbolic processing.


* **Independent Category (`Kind`) vs. Dimensionality (`Dimension`)**: Separating `Kind` (e.g., `Energy` vs. `Torque`) from `Dimension` ($M L^2 T^{-2}$) allows `MRC-003` to reject physically invalid operations even when dimensional exponents match.


* **Deterministic Exact Rational Serialization**: Enforcing `math/big.Rat` for base SI dimension exponents and scalar expression coefficients prevents floating-point precision issues. Representing rationals in JSON as canonical `"num/den"` strings (e.g., `"1/2"`) guarantees deterministic cross-platform hashes.


* **Explicit `Simplify` vs. `Identify` Firewall**: Mechanical algebraic reduction (`Simplify`) is barred from creating physical equivalences or `IDENTIFIED` provenance steps. Physical leaps require an explicit `Identify(a, b, justification)` call logged in the derivation ledger.


* **Crypto-Auditable Derivation Ledger**: Derivation steps are chained using SHA-256 over canonical `StepEnvelope` structures. Replay validation (`Validate()`) ensures that any modification to intermediate expressions invalidates the ledger.



---

### Minor Inconsistencies & Edge-Case Clarifications

Before passing these files to the implementation architect, resolve these three minor points:

| Issue Location | Observation / Conflict | Resolution / Action Required |
| --- | --- | --- |
| **File Guardrail vs. Repo Layout** (`specs_v2.md` §3 vs. §37.1)

 | Section 3 lists **47 explicit repository files**. However, Section 37.1 and Prompt §15 impose a **40-file total guardrail**.

 | Clarify in the implementation plan's Section 15 (Scope Accounting) that the 47 files specified in Section 3 are normatively required, or direct the coding agent to merge fine-grained `ops/*.go` files (e.g., combining `arithmetic.go`, `simplify.go`, and `substitute.go`) if strict adherence to $\le 40$ files is enforced.

 |
| **Duplicate Paragraph in Prompt** (`prompt_v2.md` § Planning rules)

 | The requirement for the coverage matrix mapping normative `MUST`s and `MRC-001..008` is printed twice back-to-back.

 | Delete the duplicate sentence in `physics_compiler_mvp_implementation_plan_prompt_v2.md`.

 |
| **N-Ary Canonical Node Sorting** (`specs_v2.md` §8)

 | Section 8 specifies child node sorting for `Add` and `Mul` by (1) node kind ordinal, (2) canonical child hash, and (3) canonical byte sequence.

 | Ensure the symbolic engine explicitly implements string/byte-level lexicographical comparison on SHA-256 canonical hashes to avoid relying on non-deterministic Go `map` iteration.

 |

---

### Pre-Flight Checklist for Handoff

When instructing the coding agent to produce the final execution plan, include these execution instructions:

1. **Adhere Strictly to the 15 Required Output Sections**: Confirm the plan outputs sections 1 through 15 in exact sequence as required by `prompt_v2.md`.


2. **Execute External Test Packages for Construction Negatives**: Ensure all negative tests for constructor authority (`MRC-001`) run from external packages (e.g., `package mechanics_test`) to verify unexported field boundaries.


3. **Complete the Mandatory Coverage Matrix**: Verify that Section 12 maps every normative `MUST` in `specs_v2.md` and every rule from `MRC-001` through `MRC-008` to a specific implementation file and test case.
Based on a meticulous, line-by-line reconciliation of the **v3 GR Implementation Plan**, the **v3 Adversarial Review**, the **frozen `specs_v2_3.md`**, and the **42-file baseline audit**, I have synthesized the **Final GR Implementation Contract (v4)**. 

The v3 review identified 11 critical refinements (P1–P9, plus the Schwarzschild trigonometric requirement and the wording lock). **Every single one has been resolved and locked.** There are zero remaining architectural blockers, zero ambiguities, and zero choices left for the coding agent. 

Below is the definitive, executable contract that the coding agent must follow.

---

# THE FINAL GR IMPLEMENTATION CONTRACT (v4)
**Status:** Normative Execution Contract. Planning document only.
**Authority:** `plan10/specs_v2_3.md` (frozen), `plan10/plan10_v2_3.md` (frozen), `adv_review11` (42-file re-freeze authorized).
**Rule:** Append-only. No modifications to `phys`, `specs`, or manifests.

## 1. The 11 Locked Amendments (Resolving the v3 Review)
The coding agent must implement the following exact resolutions to the v3 review findings:

*   **A1 (GR-3 Differentiation Strategy):** GR-3a *must* first exercise the frozen `phys.Differentiate` on a HYPOTHESIS bridge of `Pow(r, -1)`. It *must* record `UnsupportedOperationError` as `SPEC-INTENDED-BOUND`. GR-3b *must* continue using the Level-1 `phys-gr/symbolic/differentiate.go` engine.
*   **A2 (Pin A Fixture Correction):** Pin A *must* use `Identify(RestMass, E/c², justification)`. Dimensions match (`M`). The test *must* assert: `Kind == KindRelation`, `Expr` is `Relation(eq, ...)`, `Provenance == IDENTIFIED`, `CorpusStatus == NONE`, and operands remain unmutated.
*   **A3 (GR Module Home):** GR lives exclusively in `github.com/PithomLabs/phys-gr`. It depends one-way on `github.com/PithomLabs/phys`. Local development uses `go.work`. The 42-file `phys` tree is never touched.
*   **A4 (Level-2 Governance Restored):** L1 → L2 promotion requires a named second established consumer + non-leakage of assumptions/provenance. L2 → L3 requires the full Growth Gate. Level-2 libraries *cannot* bypass kernel invariants.
*   **A5 (C1–C8 Ledger Enumerated):** The concern ledger is no longer a vague reference. All 8 concerns are explicitly defined and routed (see §4 below).
*   **A6 (Silent-Wrongness Probe Scheduled):** A dedicated hostile test is added: differentiating a bare `Symbol("f")` wrt `t`. The system *must* return `0` (`ACCEPTED-CORRECT`). To get a non-zero derivative, the AI *must* use `UnknownFunction("f", [t])`.
*   **A7 (20-Field Evidence Record Enumerated):** The Growth Evidence Record is fully specified. Fields 18–20 are executable (see §5 below).
*   **A8 (42-File Baseline Cited):** The baseline authorization is explicitly cited as `plan10/adv_review11.md` (DOCUMENTATION PASS). PASS0 verifies membership before editing.
*   **A9 (Pin B Spec Wording Clarified):** Pin B (`Differentiate(Pow(x, -1))` → `UnsupportedOperationError`) is a *conformance pin* verifying the frozen spec's fail-closed behavior, not a modification of the spec. Target file: `ops/negative_test.go`.
*   **A10 (Trigonometric Mandate):** The standard spherical Schwarzschild metric contains `r² sin²θ dφ²`. `phys-gr/symbolic/trig.go` *must* implement `Sin` and `Cos` nodes and their derivatives. This is Level-1 math, not kernel expansion.
*   **A11 (Wording Lock on Growth):** The phrase "Default outcome: NO KERNEL GROWTH" is replaced with: *"No kernel modification is presumed. All Growth Gate outcomes remain admissible and must be determined from recorded evidence."*

---

## 2. Module & Dependency Mechanics
*   **Module Identity:** `github.com/PithomLabs/phys-gr`.
*   **Dependency Direction:** `phys-gr → phys` (one-way). `phys` never imports `phys-gr`.
*   **Local Development:** `go.work` at the shared parent directory connecting pristine siblings. CI pins exact `phys` commit.
*   **Forbidden:** `replace` to a modified fork, vendoring, direct `internal/kernel` import.
*   **Baseline Discipline:** Before writing GR code, verify the 42-file `phys` tree is clean, build/vet/test pass, and record the exact `phys` git SHA and SHA-256 of `specs_v2_3.md` + manifests.

## 3. PASS0 Prerequisite (The Conformance Gate)
PASS0 is a separate task. The GR agent *verifies* it; it does not *perform* it. If PASS0 fails, GR execution halts immediately (`ESCALATE-TO-SPEC`).
*   **Pin A:** `Identify(RestMass, E/c²)` → `KindRelation`, `IDENTIFIED`, `NONE`, `Dimension M`. (See A2).
*   **Pin B:** `Differentiate(Pow(x, -1))` in `ops/negative_test.go` → `UnsupportedOperationError`. (See A9).
*   **Pin C:** `AGENTS.md` §8 caveat: named-kind compatibility protects additive operations only; `Expression + Expression` mixing is permitted by design.

## 4. The C1–C8 Concern Ledger (Fully Enumerated)
The coding agent must route every GR obstacle into exactly one of these predefined buckets:

| ID | Concern | Classification | Routing / Resolution |
| :--- | :--- | :--- | :--- |
| **C1** | **Type Decay:** `m·c²` and `r×F` both become `Expression` (M L² T⁻²). | Real Limitation | Level-2 semantic ascription investigation. Never kernel. |
| **C2** | **Ascription:** No path from derived `Expression` to nominal kind. | Real Limitation | Level-2 semantic ascription. `Substitute` stays strict. |
| **C3** | **Assumption Inconsistency:** `v << c` and `v = 0.99c` coexist. | Real Limitation | Level-2 constraint/entailment layer. |
| **C4** | **Symbol Collision:** `Symbol("m")` for electron vs proton. | Real Limitation | Level-2 namespaced symbol/entity layer. |
| **C5** | **Transcendentals:** MVP has no `sin/exp/ln`. | Not MVP Defect | **Resolved in L1:** `phys-gr/symbolic/trig.go` owns `Sin`/`Cos`. |
| **C6** | **Closed Kind Ontology:** New `Kind` requires spec evolution. | Closed by Design | Do not build subtype hierarchy in kernel. |
| **C7** | **Identify Regression:** `Identify` must return `KindRelation`. | Conformance Pin | **Resolved in PASS0:** Pin A asserts `KindRelation`. |
| **C8** | **Differentiate Regression:** Negative powers fail closed. | Conformance Pin | **Resolved in PASS0:** Pin B asserts `UnsupportedOperationError`. |

## 5. The 20-Field Growth Evidence Record
If a kernel limitation is encountered, the agent must populate this exact record. Fields 18–20 are now executable:

1. **Capability:** The missing physical/mathematical capability.
2. **Minimal Counterexample:** The exact GR derivation step that fails.
3. **Failure Category:** `SPEC-INTENDED-BOUND` / `PACKAGE-SOLVABLE` / `REPRESENTABLE-BUT-UNFAITHFUL` / `UNREPRESENTABLE` / `SILENTLY-WRONG`.
4. **Level-1 Attempt:** How `phys-gr` tried to solve it.
5. **Level-2 Attempt:** How a shared L2 library tried to solve it.
6. **L1 Failure Reason:** Why the L1 attempt failed to preserve the invariant.
7. **L2 Failure Reason:** Why the L2 attempt failed.
8. **Genericity:** Why this is not just a GR-specific hack.
9. **Foundationality:** Why this is a universal formal invariant.
10. **Second Consumer:** Name of the second established theory requiring this.
11. **Second-Consumer Example:** Worked example in the second theory.
12. **Non-Goal Collision:** Check against `specs_v2_3.md` §37 deferrals.
13. **Silent-Wrongness Analysis:** Proof that the kernel isn't silently accepting garbage.
14. **Artifact Compatibility:** Impact on canonical JSON/SHA-256 hashing.
15. **MRC Impact:** Impact on the 8 MRC rules.
16. **Replay Impact:** Impact on session/ledger replay.
17. **Migration Impact:** Impact on existing `phys` tests.
18. **Minimal Kernel Change:** Package/file, exported identifiers, change kind (add API / alter invariant / add node / alter encoding). **NO SOURCE CODE.**
19. **Independent Review:** Reviewer identity, record ID, first/second classification, disagreement resolution.
20. **Human Approval:** Record ID, reviewer identity, UTC timestamp, `approved|rejected|provisional`, decision note.

*Rule:* No human approval (Field 20) = `PROVISIONAL` = no kernel modification.

## 6. The GR-Local Mathematics Mandate (Level-1)
The `phys-gr` package must implement the following bounded mathematical machinery. **None of this enters the kernel.**

*   **`GRExpr` Algebra (9 nodes):** `Symbol`, `Rational`, `Add`, `Mul`, `Neg`, `Pow` (int only), `Sin`, `Cos`, `UnknownFunction`.
*   **`UnknownFunction`:** `{Name, Args []Symbol, DerivativeOrder []uint8}`. Max total order = 2. `A(r)` is `[0]`, `A'(r)` is `[1]`.
*   **`Diff` Engine:** Full product/chain rules. `Diff(Sin(u)) = Cos(u) * Diff(u)`. `Diff(UnknownFunction)` increments the order slot.
*   **`Trig` Engine:** `Sin`, `Cos` nodes. Identity: `Sin² + Cos² → 1`.
*   **`ZeroTest`:** 3-valued `{ZERO, NONZERO, UNDECIDED}`. `ZERO` only on exact normalization proof. Test-only oracle for negative controls.
*   **Bridge Contract:** `ToCoreExpr` / `FromCoreExpr` / `ToCoreObject`. Lazy, temporary HYPOTHESIS objects. Content-derived ID (`gr/bridge/<SHA-256>`). If `GRExpr` contains `Sin`, `Cos`, or `UnknownFunction`, bridge returns `UnrepresentableKernelProjectionError`.

## 7. Silent-Wrongness & Adversarial Probes
The coding agent must implement the following hostile tests to prove the boundary semantics:

*   **Probe A (Bare Symbol):** `Diff(Symbol("f"), t)` → `0`. Status: `ACCEPTED-CORRECT`. (Syntactic symbol is constant).
*   **Probe B (Unknown Function):** `Diff(UnknownFunction("f", [t]), t)` → `UnknownFunction("f'", [t])`. Status: `EXPECTED LEVEL-1 BEHAVIOR`.
*   **Probe C (Chart Mixing):** Attempt to contract tensors from `spherical-static` and `cartesian`. Status: `REJECTED LOCALLY`.
*   **Probe D (Connection Masquerade):** Attempt to pass `Christoffel` coefficients into a generic `Tensor.Contract` op. Status: `REJECTED LOCALLY` (Connection is not a Tensor).
*   **Probe E (Dimension-Valid Semantic Error):** `Add(Energy_expr, Torque_expr)`. Both are `M L² T⁻²`. Status: `ACCEPTED` at kernel level (both are `Expression`), but flagged as `REPRESENTABLE-BUT-UNFAITHFUL` in the GR trace.

## 8. The Schwarzschild Workload (GR-8)
*   **PRIMARY (Derivation):** Start with ansatz `A(r)`, `B(r)`. Compute Christoffel → Riemann → Ricci. Certify `R_tt/A + R_rr/B = (AB)'/(rAB²)`. Vacuum → `(AB)' = 0` → `AB = k1`. Asymptotic flatness → `k1 = 1`. Certify `R_thth = 1 - A - rA'`. Vacuum → `(rA)' = 1` → `A = 1 + k2/r`. Result: `A = 1 - 2mu/r`, `B = 1/A`.
*   **SECONDARY (Verification):** Construct `A = 1 - 2mu/r`, `B = 1/A` independently. Compute Ricci. Verify `Rμν = 0` via `ZeroTest(ZERO)`.
*   **NEWTONIAN LIMIT:** `Truncate(expr, eps, 1)` where `eps = GM/(c²r)`. Extract `Phi = -GM/r`. **Do not use `phys.Limit`.**
*   **SESSION SEGMENT:** Sole use of `phys.Session`. Compare weak-field `g_tt` expansion vs `-(1+2Phi/c²)`. Use `Session.Identify` *only* for the integration constant correspondence (`k2 = -2mu`). Result: `KindRelation`, `IDENTIFIED`, `NONE`.

## 9. The "Never" List (Architectural Guardrails)
The coding agent will fail the review if it commits any of the following:
1.  **Never** add `Sin`, `Cos`, `UnknownFunction`, or Tensor nodes to `internal/kernel`.
2.  **Never** expand `phys.Differentiate` to support negative-integer powers.
3.  **Never** expand `phys.Limit` to support asymptotic/weak-field reductions.
4.  **Never** use placeholder symbols (e.g., `sin_theta`) to bypass the lack of trig functions in `phys`.
5.  **Never** allow GR-local operations to enter `phys.Session.Step` as new `OperationID`s.
6.  **Never** treat the derivation of Schwarzschild as a mere verification; it must be a derivation from `A(r), B(r)`.
7.  **Never** modify the 42-file `phys` tree. Not even to add a comment.
8.  **Never** invent a "Substrate Efficacy Ratio" or threshold metric.
9.  **Never** automatically promote a Level-1 abstraction to Level-2 without a named second consumer.

---

### Final Handoff Verdict
**ARCHITECTURAL BLOCKERS = 0**
**UNPINNED IMPLEMENTATION CHOICES = 0**
**KERNEL CHANGES AUTHORIZED = 0**

The v3 review's 11 findings have been fully integrated. The boundary between the frozen kernel (`phys`), the Level-1 GR mathematics (`phys-gr`), and the Growth Gate is now mathematically and architecturally absolute. 

**You are cleared to hand this exact contract to the coding agent.** It has everything it needs to execute the most rigorous stress test the Physics Compiler MVP has ever faced, without making a single unauthorized architectural decision.
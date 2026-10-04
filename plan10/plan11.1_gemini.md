An adversarial review of **Plan 11.1 (`plan11.1.md`)** identifies key architectural friction points, hidden contradictions, and execution risks embedded within the proposed GR implementation strategy.

---

## 1. Core Architectural Friction & Deadlocks

### Critical Blocker 1: The Pin B Differentiation Paradox vs. GR Anchor Deadlock

* **The Conflict:** PASS 0 Pin B explicitly locks and tests that `Differentiate(Pow(x, negative-integer))` fail-closes with `UnsupportedOperationError`. However, Section 4 anchors the entire GR workload on static spherically symmetric spacetime (Schwarzschild metric), where $g_{00} = -(1 - 2GM/r)$ and $g_{11} = (1 - 2GM/r)^{-1}$.
* **The Failure Mode:** Step GR-3 (Christoffel symbols $\Gamma^\lambda_{\mu\nu} = \frac{1}{2} g^{\lambda\sigma} (\partial_\mu g_{\nu\sigma} + \partial_\nu g_{\mu\sigma} - \partial_\sigma g_{\mu\nu})$) requires differentiating terms containing $r^{-1}$ and $r^{-2}$. Because Pin B locks `Differentiate` to reject negative exponents in the kernel, **GR-3 will unconditionally fail at the kernel boundary**.
* **The Trap:** Either Level 2 (`physmath`) must build an entire parallel differentiation engine outside `phys` (rendering `phys.Differentiate` dead code for field theory), or the workload will stall at GR-3 and be classified as `SPEC-INTENDED-BOUND`. Pre-locking Pin B while anchoring on Schwarzschild guarantees that `phys` cannot execute its own workload without a complete Level-2 bypass.

### Critical Blocker 2: Level-2 Semantic Ascription vs. MRC-003 Interoperability

* **The Conflict:** PASS 0 Pin A enforces that `Session.Identify` produces `KindRelation` and explicitly prohibits re-labeling derived expressions (e.g., $E/c^2$) as named physical kinds (e.g., `RestMass`). Section 6 routes semantic ascription entirely to Level 2.
* **The Failure Mode:** If Level 2 maintains semantic types via external wrappers around `phys.Object`, those ascribed objects cannot be passed back into `phys` kernel functions that enforce MRC-003 named-kind checking.
* **The Trap:** If Level 2 derives $m_{\text{eff}} = E/c^2$, wraps it as `RestMass` in `physmath`, and then attempts to pass it to a `phys` kernel function expecting `KindRestMass`, `phys` will inspect the underlying `phys.Object` (`KindExpression`) and reject the operation under MRC-003. Level 2 semantic ascription is isolated from the kernel's type checker, creating an unusable abstraction wall.

### Critical Blocker 3: Export Surface Lockout under the 42-File Freeze

* **The Conflict:** Rule D3 isolates Level 2 into a separate Go module (`physmath`) depending one-way on frozen `phys`, with the 42-file tree locked.
* **The Failure Mode:** Level 2 can only interact with `phys` through its currently exported APIs. If Level 2 index algebra or tensor builders require inspecting internal AST node structures, unexported fields, or private helper functions (e.g., custom AST pattern matching for contraction), Level 2 will hit a hard wall.
* **The Trap:** Because the 42-file tree cannot be touched, missing export methods cannot be added to `phys` without invoking the Level-3 Growth Gate process. The plan assumes the current public API of `phys` is sufficient for external module expansion without having verified the export surface requirements of Level 2.

### Critical Blocker 4: AST Explosion and Equivalence Failure in Tensor Contractions

* **The Conflict:** Section 4 embeds index algebra ($g_{\mu\nu} g^{\nu\rho} = \delta_\mu^\rho$) and tensor derivatives ($R_{\mu\nu} = \partial_\alpha \Gamma^\alpha_{\mu\nu} - \dots$) into Level 2.
* **The Failure Mode:** Schwarzschild curvature expressions involve dense sums of products of rational terms ($r^{-1}$, $r^{-2}$) and trigonometric functions ($\sin\theta$, $\cos\theta$).
* **The Trap:** `phys.Simplify` provides bounded, scalar-level algebraic simplification. It does not perform rational function canonicalization, common denominator expansion, or trigonometric reduction. Generating $R_{\mu\nu}$ via explicit matrix operations over `phys.Object` will yield massive, unsimplified ASTs. When verifying the vacuum field equation $R_{\mu\nu} = 0$ (GR-5), `phys` will fail to reduce the unsimplified AST to `0`, resulting in a false `REPRESENTABLE-BUT-UNFAITHFUL` or `UNREPRESENTABLE` classification.

---

## 2. Process & Governance Vulnerabilities

### Vulnerability 1: Confirmation Bias in the "NO KERNEL GROWTH" Default

Section 6 pre-locks the default outcome of the GR stress test to `NO KERNEL GROWTH`.

* **Risk:** When Level 2 inevitably encounters substrate limitations (such as Pin B differentiation failure or AST explosion), the evaluation framework is incentivized to classify every failure as `SPEC-INTENDED-BOUND` or `PACKAGE-SOLVABLE` via Level-2 workarounds.
* **Consequence:** If Level 2 is forced to duplicate differentiation, symbol scoping, and simplification to bypass kernel constraints, `phys` degrades into a passive data wrapper rather than an active mathematical physics substrate.

### Vulnerability 2: Triviality of Pin A Verification

Pin A requires testing that `Identify(Energy, E/c², ...)` does not yield `RestMass`.

* **Risk:** In `session/session.go:278`, `Identify` already hardcodes the return type to `KindRelation`. Testing that it does not return `RestMass` verifies a condition that was never implemented.
* **Consequence:** Pin A creates a false sense of security regarding type safety without addressing the underlying requirement: how userland safely binds derived expressions to physical categories for downstream kernel operations.

### Vulnerability 3: Ambiguity in PASS 0 Mutation Boundaries

Section 0 claims `specs_v2_3.md` and the 42-file tree are frozen. Yet Section 2 (PASS 0) mandates adding new tests in `session/session_test.go`, `ops` test files, and editing `AGENTS.md`.

* **Risk:** Modifying test files and documentation inside `phys` technically mutates the 42-file repository tree and checksums prior to PASS 1.
* **Consequence:** The baseline definition is ambiguous. PASS 0 must explicitly declare whether test file edits are exempt from the 42-file freeze manifest or if a re-manifesting step is required before PASS 1.

---

## 3. Required Plan Adjustments

1. **Clarify Pin B vs. GR-3 Strategy:** Explicitly state in Section 4 whether GR-3 will bypass `phys.Differentiate` via a `physmath` symbolic differentiation pass, or if GR-3 is expected to fail at the kernel boundary and trigger a formal failure log entry.
2. **Define Level-2 $\to$ Kernel Type Bridge:** Define how a Level-2 ascribed object (`RestMass` wrapper around `KindExpression`) can be used in kernel operations without violating MRC-003 or requiring kernel modification.
3. **Audit Export Surface before PASS 1:** Run a static analysis pass on `phys` public exports to verify that AST constructors, node inspectors, and expression iteration methods are sufficient for external module consumption.
4. **Establish AST Size & Simplification Limits:** Include explicit AST depth/node bounds and simplification expectations for GR-4 through GR-6 in the workload anchor to distinguish between kernel simplification failure and structural invalidity.
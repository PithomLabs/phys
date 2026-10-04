### 1. Critical Architectural Friction & Structural Paradoxes

#### Paradox 1: Trigonometric Hybrid Modeling ($g_{33} = r^2 \sin^2\theta$) & AST Disconnect

* **The Conflict:** Section 6 mandates that tensor components in $g_{\mu\nu}$ are stored as `core.Expr` nodes. Section 9 states that trigonometric functions ($\sin\theta, \cos\theta$) are owned by `phys-gr` userland (`symbolic/trig.go`) because the MVP expression node set is closed (`Call = lorentz_factor` only).


* **The Failure Mode:** In the Schwarzschild metric, $g_{33} = r^2 \sin^2\theta$. Because `core.Expr` cannot represent $\sin\theta$, $g_{33}$ cannot exist as a pure `core.Expr`.


* If `phys-gr` models $\sin\theta$ as an uninterpreted `core.Symbol("sin_theta")`, taking $\partial_\theta g_{33}$ via `phys-gr/symbolic/differentiate.go` or `phys.Differentiate` with respect to `core.Symbol("theta")` will return `0` due to syntactic variable mismatch. This directly triggers the **`SILENTLY-WRONG`** failure probe defined in Section 12.


* If `phys-gr` models $g_{33}$ as a hybrid wrapper (combining `core.Expr` for $r^2$ with custom trig AST nodes), `phys` kernel functions (`Add`, `Multiply`, `Simplify`, `Substitute`) can no longer process $g_{33}$ or perform tensor contractions on $g_{\mu\nu}$.




* **Impact:** The kernel is cut off from processing the angular components of Schwarzschild spacetime, forcing `phys-gr` to duplicate rational component algebra for non-`core.Expr` structures.

#### Paradox 2: The `phys.Limit` Trailing-Body Lockout in GR-8

* **The Conflict:** Section 7 lists `phys.Limit` as an approved kernel interaction point for the Newtonian limit in GR-8, but explicitly restricts it: *"`Call` traversal fixed-body `lorentz_factor` only"*.


* **The Failure Mode:** Weak-field and Newtonian limit calculations (GR-8) require evaluating asymptotic limits of metric components as $r \to \infty$ or $GM/(c^2 r) \to 0$. None of these expressions contain Special Relativity's `lorentz_factor` ($\gamma = 1/\sqrt{1 - v^2/c^2}$).


* **Impact:** Passing Schwarzschild metric components into `phys.Limit` will fail or return unchanged expressions because `phys.Limit` only traverses `lorentz_factor` AST nodes. `phys-gr` will be forced to bypass `phys.Limit` entirely and build a standalone limit evaluator in `phys-gr/limit/`.



#### Paradox 3: `manifest.json` Hash Integrity vs. PASS 0 Test File Mutation

* **The Conflict:** Section 0 and Section 3 declare that `specs_v2_3.md`, both `manifest.json` files, and the 42-file tree are byte-for-byte frozen. Section 2 and Section 16 instruct PASS 0 to modify `session/session_test.go`, `ops/differentiate_test.go` (or ops test files), and `AGENTS.md`.


* **The Failure Mode:** If `manifest.json` tracks SHA-256 cryptographic content hashes of repository files rather than plain file paths, editing `session_test.go` alters its SHA-256 digest, causing `go test ./...` baseline verification to fail immediately.
* **Impact:** PASS 0 cannot execute without either breaking `manifest.json` integrity or requiring an explicit manifest re-hash authorization step that contradicts the "byte-for-byte frozen" declaration.



#### Paradox 4: Go Module `replace` Directive Ambiguity in Sibling Repositories

* **The Conflict:** Section 4 defines `phys-gr` as a sibling Go module (`[github.com/PithomLabs/phys-gr](https://github.com/PithomLabs/phys-gr)`) depending on `[github.com/PithomLabs/phys](https://github.com/PithomLabs/phys)` and strictly states: *"No `replace` into kernel internals"*.


* **The Failure Mode:** During local development and CI runs, if `[github.com/PithomLabs/phys](https://github.com/PithomLabs/phys)` is located on local disk and not published to a remote Git tag, `go build` and `go test` inside `phys-gr` will fail to resolve the dependency without a `replace` directive in `phys-gr/go.mod` (e.g., `replace [github.com/PithomLabs/phys](https://github.com/PithomLabs/phys) => ../phys`).
* **Impact:** The wording "No `replace` into kernel internals" is ambiguous. Execution agents may interpret it as a total ban on `replace` directives in `go.mod`, stalling local compilation.



---

### 2. Governance & Evaluation Blindspots

#### Blindspot 1: Triviality of Pin A Verification

Section 16 mandates asserting that `Identify(RestMass, E/c²)` returns `KindRelation` and does not produce `RestMass`. As acknowledged in Section 3, `session/session.go:278` already hardcodes `KindRelation` as its return type. Testing this verifies static code already present in the codebase. It does not solve or test how higher layers can associate derived scalar expressions ($E/c^2$) back to physical categories (`RestMass`) for downstream type-checked kernel operations.

#### Blindspot 2: The "Userland Bypass Fallacy" & Missing Substrate Efficacy Metric

Section 18 defines `NO-GROWTH` as an acceptable outcome if all GR needs are met in userland (`phys-gr`).

* **The Risk:** If `phys-gr` is forced to implement its own differentiator (`symbolic/differentiate.go`), custom trig handler (`symbolic/trig.go`), custom limit evaluator (`limit/`), custom tensor algebra (`index/`), and custom expression simplifier because `phys` rejects or fails every operation, `phys-gr` will have built an entire parallel computer algebra system on top of a dormant kernel.


* **Consequence:** Declaring `NO-GROWTH` under total userland bypass misrepresents substrate inadequacy as kernel success. The plan lacks a metric to distinguish between *effective kernel support* and *complete userland bypass*.

---

### 3. Required Plan Adjustments

1. **Resolve Angular Metric Representation in Section 6 & Section 9:**
Specify how $g_{33} = r^2 \sin^2\theta$ is represented in `phys-gr`. If `sin_theta` is modeled as an opaque `core.Symbol`, explicitly require `phys-gr/symbolic/differentiate.go` to maintain a symbolic derivative lookup table ($\partial_\theta \text{Symbol}("sin\_theta") \to \text{Symbol}("cos\_theta")$) to avoid `0`-derivative silent wrongness.


2. **Clarify `phys.Limit` Scope in Section 7:**
Explicitly state that `phys.Limit` cannot process GR metric limits due to its `lorentz_factor` restriction, and record this step as a `SPEC-INTENDED-BOUND` probe during GR-8, routing Newtonian limit evaluation to `phys-gr/limit/`.


3. **Clarify `manifest.json` Protocol in Section 0 & Section 16:**
Specify whether `manifest.json` checks file membership or cryptographic SHA-256 hashes. If hashes are checked, explicitly authorize a PASS 0 re-manifesting step for `session_test.go` and `ops` test files prior to freezing PASS 1.


4. **Disambiguate `replace` Directives in Section 4:**
Clarify Section 4: *"Local module path replacement (`replace [github.com/PithomLabs/phys](https://github.com/PithomLabs/phys) => ../phys`) in `phys-gr/go.mod` is permitted for local development; referencing unexported `phys/internal/...` packages is strictly forbidden."*

5. **Establish a Substrate Efficacy Metric in Section 18 & Section 20:**
Add a quantitative metric to the Definition of Done:

$$\text{Substrate Efficacy Ratio} = \frac{\text{Kernel-Evaluated AST Operations}}{\text{Kernel-Evaluated AST Operations} + \text{Userland-Bypassed Operations}}$$



If the ratio falls below an agreed threshold, the outcome must be classified as `UNREPRESENTABLE` / `KERNEL-GROWTH CANDIDATE` rather than a blanket `NO-GROWTH`.
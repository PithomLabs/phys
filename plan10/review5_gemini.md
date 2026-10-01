An adversarial review of the `plan10_v2_3.md` implementation plan reveals several internal logical gaps, unchecked edge cases, and structural limitations. Because the scope strictly excludes hostile-source protection (relying on Go `internal/` package boundaries) and cryptographic authenticity, this review focuses on API abuse, state corruption, and mathematical bypasses that could occur during derivation replay or artifact ingestion.

### 1. Symbolic & Mathematical Execution Gaps

* **`Limit` Variable Type Vulnerability:** The `Limit` operation requires three inputs: `target`, `variable`, and `value`. The plan specifies "substitution of value for variable's single symbol".


* *Adversarial Gap:* If a malicious or malformed step passes a non-symbol expression (e.g., `Add(x, y)` or `Rational(2)`) as the `variable`, the substitution logic could panic or silently fail to replace anything, resulting in an invalid limit evaluation.
* *MVP Resolution:* Add a strict precondition check in `ops.Limit`: `if variable.Expr().Kind() != kernel.ExprSymbol { return UnsupportedOperationError }`.


* **`SelectBranch` Structural Blindness:** The exact six-step contract for `SelectBranch` requires verifying that branches are "structurally `[b, Neg(b)]`" and then selecting "the first (structurally nonnegative `Sqrt`) branch".


* *Adversarial Gap:* An adversary could manually construct a `BranchSet` where the array order is inverted (`[Neg(b), b]`). If the implementation blindly selects index `0`, it will select the negative branch, violating the physical entailment rules.
* *MVP Resolution:* The operation must dynamically inspect the root node of the selected branch to ensure it is not `Neg(x)` rather than relying on slice index `0`.


* **Positional Array Bounds Panics:** The `ops.Apply` dispatch uses exact positional indices up to `inputs[2]` (for `substitute`). The plan notes that "unused input positions are omitted".


* *Adversarial Gap:* If a corrupted ledger JSON provides fewer inputs than an operation requires (e.g., calling `substitute` with only 1 input), direct access to `inputs[2]` will cause a Go runtime panic (out of bounds), crashing the validation loop rather than returning a structured error.
* *MVP Resolution:* `Session.Step` and `ops.Apply` must enforce a strict `len(inputs) == required_arity(operation_id)` check before evaluating parameters.



### 2. State Machine & Session Vulnerabilities

* **Ledger Hash Collision via Step Index Manipulation:** The ledger tamper detection relies on `CurrentStepHash = SHA256(CanonicalJSON(StepEnvelope{PreviousHash, Step: body}))`. The plan specifies that `StepID` is "derived solely from index".


* *Adversarial Gap:* An adversary could forge a ledger where `Index` jumps (e.g., from 1 to 5) but the `StepID` is manually hardcoded as `step-000005` to match. If the replay engine blindly accepts the parsed `StepID` instead of recomputing it iteratively from the replay loop count, the hash chain will successfully validate a skipped-step ledger.
* *MVP Resolution:* During `Session.Validate()`'s 17-step pipeline, the engine must strictly recompute `StepID` based on the validation loop's internal counter, ignoring the JSON's `StepID` until the equality assertion.




* **Identify Justification Bypasses:** `Session.Identify` requires a "non-empty trimmed justification".


* *Adversarial Gap:* Because there is no semantic review, an adversary can supply arbitrary whitespace or single-character strings (e.g., `"a"`), satisfying the non-empty requirement and polluting the provenance graph with meaningless justifications.


* *MVP Resolution:* This is an accepted limitation of the MVP. Since external reviewers consume the `Review` artifact, this is addressed by ensuring the UI/consumer-side throws a `Challenge` with `Category: ProvenanceProblem` when a justification is mathematically meaningless.





### 3. Artifact & Provenance Limitations

* **Disconnected Falsifiability Expressions:** `Prediction`, `FalsificationCondition`, and `RecoveryClaim` carry `core.Expr` nodes. The validation pipeline checks that "IDs non-empty when entries present".


* *Adversarial Gap:* A malicious user could construct a `ResearchCandidate` where the falsifiability expressions are syntactically invalid (e.g., broken trees, invalid `big.Rat` pointers) or completely unrelated to the derived hypothesis. The plan explicitly states these structures record test structure "without any truth/false decision". If `Seal` does not validate the internal `Expr` objects, it mints a trusted candidate containing corrupt mathematical nodes.


* *MVP Resolution:* Add a step in `UnverifiedResearchCandidate.Validate` (and `Seal`) that calls `.Valid()` on every embedded `core.Expr` within the falsifiability structures before returning a trusted `ResearchCandidate`.


* **The Authenticity Illusion:** The plan explicitly defers cryptographic signing, stating that "anyone rewriting the full ledger can recompute hashes".


* *Adversarial Gap:* While known to the spec, the practical danger is that an external consumer checking `.Validate() == nil` might mistakenly assume the `ResearchCandidate` is genuinely authored by the claimant, not just structurally valid.


* *MVP Resolution:* Ensure that `UnverifiedResearchCandidate.Validate()` is strictly documented in its GoDoc to assert *structural and mathematical consistency only*, and requires transport-layer authenticity (e.g., verified Git commits) for identity verification.


## Audit Review

Yes, the previous adversarial review still holds completely.

The attached "adv_review5.md" document is strictly a **Physics Corpus Integrity Audit**, which evaluates the correctness of the encoded domain knowledge and kernel separation rather than the runtime security or state machine robustness of the underlying Go implementation.

Here is why the operational and state machine vulnerabilities remain unaddressed despite the overall "PASS" verdict in the attached audit:

* **Divergent Audit Scopes:** "adv_review5.md" focuses exclusively on verifying that mechanics and special relativity concepts (e.g., Mass, Time, Minkowski spacetime) are properly initialized with the correct dimensions and assumptions. It verifies that the `internal/kernel/` remains physics-agnostic, but it does not evaluate the input validation limits, array bounds, or structural edge cases of the `ops` package.


* **Derivations Audited for Hardcoding, Not Abuse:** The audit confirms that derivations like E=mc², the kinetic energy derivative, and the Lorentz limit execute strictly from the encoded corpus and operations without any hardcoded results. However, it does not test what happens when a user feeds malformed inputs into these operations—such as passing non-symbol variables to the `Limit` operation or inverted arrays to `SelectBranch`—which leaves the previously identified runtime panics and mathematical bypasses open.


* **Corruption Detection is Corpus-Bound, Not Ledger-Bound:** The adversarial corruption detection detailed in "adv_review5.md" successfully verifies that changes to canonical expressions (like F=ma), dimensional parameters, or assumptions (like m≥0 and c>0) are caught by manifest cross-checks and hash comparisons. It does not test or mitigate maliciously constructed session ledgers where step indices are manipulated to spoof hash collisions, nor does it confirm whether expressions nested inside a `ResearchCandidate`'s falsifiability conditions are mathematically valid.



In summary, "adv_review5.md" successfully proves the MVP correctly and cleanly models physics without embedding domain secrets in the engine. It does not invalidate or resolve the previous findings regarding API fragility, slice bounds-checking panics, and validation bypasses in the session replay loop.
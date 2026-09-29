# REVIEW OF THE MVP PROMPT AND SPECS
## First Principles, Internal Consistency, Architectural Completeness

**Status:** Both documents are a significant improvement over previous iterations. They are the best specification artifacts in this project so far. But they still do not achieve the claim of "no architectural decisions for the coding agent." This review identifies where.

---

## Part I: Overall Verdict

**The prompt is good.** It is clear about the normative status of the specs, disciplined about non-negotiable boundaries, and specific about the required deliverables.

**The specs are mostly good.** They are precise where precision is possible, explicit about deferrals, and disciplined about scope.

**But together, they leave at least ten architectural decisions to the coding agent.** The specs claim to remove architectural ambiguity. They remove most of it. They do not remove all of it. And the prompt's own discipline rule ("do not say 'the coding agent can decide'") is violated by the specs themselves.

The prompt says:

> "another coding agent should be able to execute it against the specifications without making architectural decisions."

This is the test. It does not pass yet.

---

## Part II: The Prompt — Strengths and Weaknesses

### Strengths

| Element | Why It Works |
|---|---|
| **Normative hierarchy** ("specs win") | Removes the ambiguity from previous prompts. |
| **Non-negotiable boundaries** | Prevents the coding agent from reopening settled questions. |
| **15 required outputs** | Forces structure. The coding agent cannot produce a vague plan. |
| **Planning discipline** | Explicitly forbids alternatives, "agent can decide," second sources of truth. |
| **Mandatory validation focus** | Lists the 15 things the plan must prove. Clear acceptance target. |
| **Scope accounting** | Requires the plan to report files, types, operations. Enables bloat detection. |

### Weaknesses

**Weakness 1 — The prompt inherits the specs' gaps.**

The prompt says "treat the specs as the contract." If the specs are incomplete, the prompt is incomplete. The prompt does not add any constraint that would force the coding agent to resolve the specs' gaps.

**Weakness 2 — The prompt does not require the plan to list the architectural decisions the coding agent had to make.**

A good plan would say: "Here are the places where the specs left a choice, and here is what I chose." The prompt does not require this. The coding agent will make choices and hide them.

**Weakness 3 — The prompt's scope accounting requirement is unspecified.**

"Provide a concise estimate of files, major types, and major operations so reviewers can detect bloat" — but what counts as bloat? Is 20 files too many? 50? 100? The plan needs a reference point.

**Weakness 4 — The prompt does not reference the six product requirements.**

The specs list them (§0). The prompt does not. The coding agent reads the specs, which is fine, but the prompt should echo them for emphasis.

**Weakness 5 — "One page maximum" for the architecture summary is arbitrary.**

The architecture of this MVP cannot be summarized in one page. The constraint is aesthetic, not functional. The coding agent will either violate it or omit critical information.

---

## Part III: The Specs — Section by Section

### §0–1 (Authority, Thesis) — Good

Clear. No issues.

### §2 (Scope)

**§2.1 vs §3 inconsistency.** §2.1 lists `cmd/physvet/` as a package to populate but marks it "NOT implemented in MVP." §3 does not list it in the required tree. The coding agent will be confused: is the directory created (empty) or not created?

**Fix:** Remove `cmd/physvet/` from §2.1's list. Defer it entirely.

**§2.2 static analyzer decision.** Clear and correct. `physvet` is deferred.

**§2.3 canonical physics content.** Good. Concrete list.

**§2.3 says "The relativity package therefore means special-relativity corpus content only for this release."** Good. But this should also apply to any framework package in v0.5+.

### §3 (Repository layout) — Good, with one issue

The tree is explicit. But `docs/paper-translation/` is required. This is unusual for a code MVP. The spec should say why the docs are part of the deliverable (they are part of the AI corpus; see §24).

**Missing:** The tree does not show where tests live. Are they adjacent to the code (`*_test.go` next to source) or in a separate `testdata/` directory? Go convention is adjacent, but the spec should say.

### §4 (Core object model) — Critical gap

**§4.1 says:** "Do NOT alias all of them to one exported `Object` type."

**The gap:** If the domain types are distinct nominal types, how do the operations accept them? The spec lists operations (§11.1) that must accept physics objects generically:

```
Add
Subtract
Multiply
Divide
Simplify
Substitute
Differentiate
```

But these cannot accept "any physics object" if the physics objects are distinct types with no common interface.

**The spec does not specify the common abstraction.** Is it:
- An interface `core.Object` that all domain types implement?
- A generic type parameter?
- A struct-embedding pattern?
- A separate `Expr` type that wraps physics objects?

This is the single most important missing piece. Without it, the coding agent must invent the architecture.

**Fix:** Specify the common abstraction. For example:

```go
// core.Object is the interface satisfied by all physics objects.
type Object interface {
    Name() string
    Dimension() Dimension
    Expr() Expr
    Provenance() Provenance
}
```

Then domain types (e.g., `mechanics.Mass`) are concrete types that implement `core.Object`.

**§4.2 constructor authority.** Good. Unexported fields, domain constructors only.

**§4.3 immutability.** Good.

### §5 (Dimensions) — Good

Precise. `math/big.Rat` for exponents. 7 SI base dimensions. Operations listed.

**Minor:** The spec does not say how dimensions are serialized. §6.4 says hashing is canonical, but dimension canonicalization is not specified.

### §6 (Symbolic expression engine) — Good, with ambiguities

**§6.1 closed representation.** Good. Node set is explicit.

**Ambiguity 1:** `Derivative` and `Limit` are listed as nodes. Are these unevaluated (representing `∂f/∂x` as a symbol) or evaluated (the result of `Differentiate`)? If both, how does the API distinguish them?

**Ambiguity 2:** `Function` is listed but not defined. Is it a symbolic function symbol (e.g., `f(x)`)? An application? A higher-order function?

**Ambiguity 3:** The spec says the node set is "limited to" the listed nodes. But operations like `Compare` and `Solve` (§11.1) will need to produce new expression forms. Are comparison results expressions? Are solutions expressions? The spec does not say.

**§6.2 exact scalar representation.** Good.

**§6.3 canonicalization.** Good list. But "normalize relation operand ordering where semantically appropriate" is vague. Semantically appropriate according to what rule? The spec does not say.

**§6.4 equality and hashing.** Good.

### §7 (Assumption system) — Good

Precise. 7 assumption kinds. Merge policy is explicit: union, dedup, conflict error.

**Minor:** The spec says "Two different values for the same `(Kind, Key)` MUST produce an `AssumptionConflictError`." This is correct for the MVP but the spec should note that this is a conservative policy (v0.5 may allow specialization).

**§7.3 "Do not infer that one inequality or regime is more specific than another."** Good. Explicitly notes the MVP limitation.

### §8 (Conventions) — Brief but fine

### §9 (Provenance and status) — Good

Two axes, distinct. No `Truth` field. Good.

### §10 (MRC) — Critical gap

**§10.1 MRC rule IDs.** 8 IDs listed. But the spec does not say what each rule actually checks.

- **MRC-001** constructor/carrier integrity — vague
- **MRC-002** dimensional compatibility — clear
- **MRC-003** physical/category compatibility — vague. What is "category" here? The earlier designs had "Component" (Stage/Dance/Dancer) but this spec drops it. Without it, what's the check?
- **MRC-004** assumption compatibility — clear (§7)
- **MRC-005** convention compatibility — clear (§8)
- **MRC-006** explicit physical identification — clear (§11.2)
- **MRC-007** provenance/session authority — vague. How is session authority checked?
- **MRC-008** candidate containment — vague. How is containment enforced?

**The gap:** MRC-001, MRC-003, MRC-007, MRC-008 are underspecified. The coding agent will have to invent their semantics. This contradicts the prompt's "no architectural decisions" requirement.

**Fix:** For each rule, specify the check. For example:

```
MRC-001: Every physics object MUST be constructed via a domain constructor.
         Direct struct literal construction is forbidden because fields are unexported.
         Enforcement: Go's package visibility.
         
MRC-003: An operation MUST only accept operands whose domain roles are compatible.
         The compatible-role set is declared per operation.
         Example: Add requires operands with the same Dimension and the same Role.
         Example: Differentiate requires the second operand to be a Time or Coordinate.
         Enforcement: runtime check inside each operation.
```

Without this, the coding agent is making up the MRC rules.

### §11 (Operation contracts) — Good, with gaps

**§11.1 lists 13 operations.** Good.

**§11.2 Simplify vs Identify.** Clear.

**The gap:** What does each operation do? The spec says:

> "All operations MUST: validate applicable MRC rules, validate dimensions, validate assumptions/conventions, construct a new immutable result, attach provenance, be recordable by a derivation session."

But the spec does not specify:
- Which rules each operation checks
- What each operation returns
- What each operation's signature is

For example, does `Solve` take an equation and a variable and return a solution? Or does it take an expression and return a root set? The spec does not say.

**Fix:** For each operation, specify:
- Signature
- Rules checked
- Return type
- Behavior on failure

### §12 (Derivation session and ledger) — Good, with errors

**§12.2 count error:** Says "8 session actions" but lists 9: `Postulate, Declare, Define, Step, Identify, Conclude, Draft, Commit, Seal`.

**§12.2 ambiguity:** What's the difference between `Step` and `Identify`? If `Identify` is a special step, why isn't it a step kind? If `Identify` is a separate session action, does the ledger record it as a step?

**§12.3 step contents.** Good list. But what is "Operation"? If `Identify` is a step, `Operation` would be `"Identify"`. If `Identify` is separate, `Operation` doesn't apply.

**§12.4 hash chain.** Good formula.

**§12.5 replay validation.** Good.

### §13 (Manifests) — Good, with one gap

**§13.4 "manifest relation identity vs actual canonical expression."** How? The manifest has a `statement` field (a string). The Go code has a constructor that produces a canonical expression. The test must compare them. But comparing a string to a canonical expression requires parsing the string or rendering the expression. Which? The spec does not say.

**Fix:** Either:
- The manifest's `statement` field is a canonical expression string, and the test parses it.
- The manifest's `statement` field is only documentation, and the test compares the constructor's output to a hardcoded canonical expression in the test.

The first is more work. The second is simpler for MVP but weaker.

### §14 (Mechanics corpus) — Good

### §15 (Special-relativity corpus) — Good

### §16 (Framework evaluation anatomy) — Good

The requirement that "a text field saying 'reduces to classical mechanics' is NOT sufficient" is good. It forces executable reductions.

### §17 (Canonical MVP derivations) — Good

**§17.2 says "apply the rest-frame substitution."** What is `RestFrame`? §2.3 lists it as an object. Is it an assumption? A convention? A special object? The spec does not say.

**Fix:** Specify what `RestFrame` is. My recommendation: it is a constraint assumption `(Constraint, momentum, p = 0)`.

### §18 (Identify micro-test) — Good

### §19 (Hypothesis layer) — Critical gap

**§19.1 candidate objects.** Good.

**§19.2 "open representation, closed authority."** Good slogan. But the mechanism is not specified. How does the API prevent a candidate object from being treated as trusted?

**The gap:** The spec says "There is no AI-callable promotion operation." This is a negative statement. What is the positive mechanism?

Options:
- A separate type: `CandidateObject` vs `TrustedObject`
- A provenance field that is immutable after construction
- A session-scoped flag that marks the object as candidate

**Fix:** Specify the mechanism. For example:

```go
type Provenance struct {
    Status Status  // DEFINED, POSTULATED, ..., HYPOTHESIS
    ...
}

// A candidate object has Status == HYPOTHESIS.
// Operations that require trusted input MUST reject HYPOTHESIS objects.
// The session tracks which steps produce candidate-derived results.
```

### §20 (Popperian structure) — Good

### §21 (Anomaly registry) — Good

### §22 (ResearchCandidate) — Good, with one gap

**§22.3 "external JSON untrusted until validation."** Who validates? What is the validation flow? The spec says "schema validation, canonicalization, ledger hash validation, replay validation" but does not specify the API for each.

**Fix:** Specify the validation pipeline. For example:

```go
func (c *ResearchCandidate) Validate() error
```

which runs all four checks in sequence.

### §23 (Review artifact) — Brief but fine

### §24 (Paper-translation docs) — Confusing

**§24 says "part of the AI corpus but not canonical."** If they're part of the corpus, they're canonical. If they're not canonical, they're not part of the corpus. The distinction is unclear.

**Fix:** Say they are auxiliary documentation, not machine-readable canonical metadata. The manifests are canonical. The paper-translation docs are guidance.

### §25 (Negative tests) — Good

### §26 (Determinism tests) — Good

### §27 (MVP corpus/evaluation tests) — Good

### §28–32 — Good

---

## Part IV: The Critical Gaps, Ranked

| Rank | Gap | Severity |
|---|---|---|
| **1** | The common abstraction for physics objects is unspecified (§4.1). Without it, the operations cannot be typed. | **Blocking** |
| **2** | The MRC rule semantics are underspecified (§10). MRC-001, MRC-003, MRC-007, MRC-008 have no defined checks. | **Blocking** |
| **3** | The operation signatures are unspecified (§11). The coding agent must invent them. | **Blocking** |
| **4** | The candidate containment mechanism is unspecified (§19.2). The slogan is clear; the mechanism is not. | **Severe** |
| **5** | The `Derivative`, `Limit`, and `Function` node semantics are ambiguous (§6.1). Unevaluated vs. evaluated. | **Severe** |
| **6** | The manifest-constructor cross-check mechanism is unspecified (§13.4). How is the test implemented? | **Severe** |
| **7** | The `RestFrame` object type is unspecified (§17.2). Assumption? Convention? Object? | **Severe** |
| **8** | The `Step` vs. `Identify` ledger relationship is ambiguous (§12.2, §12.3). | **Moderate** |
| **9** | The validation pipeline API for external JSON is unspecified (§22.3). | **Moderate** |
| **10** | The paper-translation docs' status is contradictory (§24). Part of corpus but not canonical. | **Minor** |

---

## Part V: The Internal Inconsistencies

| Location | Issue |
|---|---|
| **§2.1 vs §3** | `cmd/physvet/` is listed in §2.1 but not in §3's tree. |
| **§12.2** | Says "8 session actions" but lists 9. |
| **§2.1** | Says "The MVP MUST populate only: ... cmd/physvet/ # NOT implemented in MVP." Contradiction: populated or not? |
| **§6.1 vs §11.1** | §6.1's node set includes `Derivative` and `Limit`; §11.1 includes `Differentiate` and `Limit` as operations. Are these the same thing? |
| **§13.1 vs §24** | §13.1 says manifest.json is "canonical machine-readable semantic source"; §24 says paper-translation docs are "part of the AI corpus but not the canonical machine-readable semantic source." If both are in the corpus but only one is canonical, the term "corpus" is overloaded. |
| **§12.2** | `Draft` and `Commit` are listed as session actions. Are they actions or state transitions? The spec doesn't distinguish. |
| **§3 tree** | Lists `docs/paper-translation/` but does not list `testdata/` or any test directory. Where do test fixtures live? |

---

## Part VI: The Prompt's Own Discipline Rule Is Violated by the Specs

The prompt says:

> "Do not say 'the coding agent can decide.'"

The specs leave at least ten decisions to the coding agent:

1. The common physics object abstraction
2. The MRC-001, MRC-003, MRC-007, MRC-008 semantics
3. The operation signatures
4. The candidate containment mechanism
5. The `Derivative`, `Limit`, `Function` node semantics
6. The manifest-constructor cross-check
7. The `RestFrame` type
8. The `Step` vs `Identify` ledger relationship
9. The external JSON validation API
10. The paper-translation docs' status

This is the central problem. The prompt's discipline rule cannot be satisfied unless the specs themselves are complete.

---

## Part VII: What the Specs Get Right

It would be unfair to end on criticism alone. The specs are the best artifacts in this project. They get the following right:

1. **Normative hierarchy.** Specs win over prompt.
2. **Scope discipline.** §2.3, §30, §31 are explicit about what is in and out.
3. **Two-axis status model.** Provenance vs. corpus status is a real contribution.
4. **Explicit MRC versioning.** No runtime override. Future revision through versioning.
5. **Ledger hash chain.** Concrete and implementable.
6. **Determinism requirements.** SHA-256, no timestamps, canonical serialization.
7. **Replay validation.** The primary anti-tampering mechanism.
8. **Manifest tests.** Cross-checking is specified as a requirement.
9. **Identify gate.** Simplify cannot perform physical identity. Explicit.
10. **Candidate containment.** The principle is clear even if the mechanism is not.
11. **Negative test list.** §25 is comprehensive.
12. **Definition of done.** §32 is specific.
13. **The MVP thesis statement.** §1 is crisp.
14. **Explicit non-goals.** §30 is thorough.
15. **v0.5 deferrals.** §31 is honest about what's postponed.

---

## Part VIII: What to Fix Before Submitting

The specs need seven specific fixes before the prompt will produce a plan with no architectural ambiguity.

### Fix 1 — Specify the common physics object abstraction

Add to §4:

```go
type Object interface {
    Name() string
    Dimension() Dimension
    Expr() Expr
    Provenance() Provenance
}
```

Domain types implement this interface. Operations accept `Object`.

### Fix 2 — Specify MRC-001, MRC-003, MRC-007, MRC-008

Add to §10:

```
MRC-001: All physics objects are constructed via domain constructors.
         Direct struct literal construction is impossible because fields are unexported.
         Enforcement: Go's package visibility.

MRC-003: An operation MAY only accept operands whose roles are in its declared compatible set.
         Each operation declares its compatible roles.
         Example: Add requires operands with the same Dimension and the same Role.
         Enforcement: runtime check.

MRC-007: All committed steps are recorded by a Session.
         No operation may produce a trusted derivation without a Session.
         Enforcement: Derivation has no public constructor; only Session.Commit produces one.

MRC-008: Candidate objects carry Status == HYPOTHESIS.
         Operations requiring trusted input MUST reject candidate objects.
         Enforcement: provenance check in each operation.
```

### Fix 3 — Specify the operation signatures

Add to §11:

For each operation, give:
- Signature
- Rules checked
- Return type

For example:

```go
func Add(a, b core.Object) (core.Object, error)
// Rules: MRC-002 (dimension), MRC-003 (role)
// Returns: a new Object with combined provenance

func Identify(a, b core.Object, justification string) (core.Object, error)
// Rules: MRC-006
// Returns: a new Object with Status == IDENTIFIED
// Side effect: appends an identification to the session ledger
```

### Fix 4 — Specify the candidate containment mechanism

Add to §19.2:

```go
// Every Object has a Provenance with a Status field.
// Trusted objects have Status in {DEFINED, POSTULATED, DERIVED, IDENTIFIED, APPROXIMATED}.
// Candidate objects have Status == HYPOTHESIS.
//
// Operations that require trusted input check the input's Status.
// If the input is HYPOTHESIS, the operation fails with CandidateContainmentError.
//
// ResearchCandidate may only be sealed if all its premises are trusted.
// A ResearchCandidate containing candidate premises is marked as provisional,
// and cannot be sealed as a complete research artifact.
```

### Fix 5 — Resolve the expression node ambiguity

Add to §6.1:

```
Derivative, Limit, and Function are UNEVALUATED nodes.
They represent symbolic derivatives, limits, and function applications.
The operations Differentiate and Limit EVALUATE these nodes and return their results.
```

### Fix 6 — Specify the manifest-constructor cross-check

Add to §13.4:

```
For each manifest item with a "statement" field, the test:
  1. Parses the statement string into a canonical expression using a simple parser.
  2. Constructs the corresponding expression via the Go constructor.
  3. Compares the two canonical expressions for structural equality.

If parsing is not feasible for MVP, the manifest item may include a "canonical_expr"
field containing a pre-serialized expression, and the test compares directly.
```

### Fix 7 — Specify the `RestFrame` type

Add to §17.2:

```
RestFrame is a constraint assumption: (Constraint, momentum, p = 0).
It is not an object.
It is not a convention.
It is an assumption applied to the derivation.
```

### Minor fixes

- **§2.1**: Remove `cmd/physvet/` from the list of populated packages.
- **§12.2**: Fix the count (8 → 9, or remove one action).
- **§12.2**: Specify the distinction between `Step` and `Identify`.
- **§24**: Clarify that paper-translation docs are auxiliary documentation, not part of the canonical machine-readable corpus.
- **§3**: Add a note about where test files live (adjacent to source, Go convention).

---

## Part IX: The Final Assessment

**The prompt is 90% there. The specs are 80% there. Together they are 75% there.**

The remaining 25% is the ten critical gaps and seven internal inconsistencies. These are fixable. None require rethinking the architecture. They require filling in the details the specs have deferred.

The prompt's discipline rule — "no architectural decisions for the coding agent" — is achievable, but only if the specs are completed first.

**The recommended sequence:**

1. Fix the seven critical gaps in the specs.
2. Fix the seven internal inconsistencies in the specs.
3. Update the prompt to echo the six product requirements.
4. Remove the "one page maximum" constraint from the prompt.
5. Add a required deliverable to the prompt: "List all architectural decisions the specs left unspecified, and state what was chosen."
6. Submit.

With these fixes, the prompt and specs will reliably produce a consistent, high-quality implementation plan.

---

## The One-Paragraph Summary

> The prompt and specs are the best artifacts in this project. The prompt is disciplined, the specs are precise where they can be, and the scope is explicit. But the combined documents still leave ten architectural decisions to the coding agent: the common physics object abstraction, the MRC-001/003/007/008 semantics, the operation signatures, the candidate containment mechanism, the `Derivative`/`Limit`/`Function` node semantics, the manifest-constructor cross-check, the `RestFrame` type, the `Step` vs. `Identify` ledger relationship, the external JSON validation API, and the paper-translation docs' status. Seven internal inconsistencies compound the problem: `cmd/physvet/` is in §2.1 but not §3; §12.2 says "8 actions" but lists 9; the term "corpus" is overloaded; and several counts and lists do not agree. The prompt's own discipline rule — "do not say 'the coding agent can decide'" — is violated by the specs themselves. The seven critical fixes are: specify the common object interface, specify the MRC rule semantics, specify the operation signatures, specify the candidate containment mechanism, resolve the expression node ambiguity, specify the manifest-constructor cross-check, and specify the `RestFrame` type. The five minor fixes are: remove `cmd/physvet/` from §2.1, fix the §12.2 count, specify `Step` vs. `Identify`, clarify the paper-translation docs' status, and note where tests live. With these twelve fixes, the prompt and specs will achieve their stated goal. Without them, the coding agent will make architectural decisions the prompt forbids, and the four-agent review will identify inconsistencies that could have been prevented. The right move is to spend one more iteration on the specs, fix the twelve issues, and then submit. The architecture is sound. The specification is not yet complete.
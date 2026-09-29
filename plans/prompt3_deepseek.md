# REVIEW OF PROMPT3.MD AS A CODING AGENT BRIEF
## Is This a Good Prompt for Producing an Implementation Plan?

**Status:** Prompt review. Not a plan review. The question is whether this prompt will reliably produce a consistent, high-quality implementation plan when given to a competent coding agent.

---

## Part I: Overall Assessment

The prompt is a **well-structured brief**. It specifies the mission, the constraints, the non-goals, and the deliverable format. It is disciplined about scope and explicit about deferrals.

But it is **not yet a prompt that will reliably produce a good plan**. Three reasons:

1. **It references external documents that are not included.** Plan Z, DeepSeek's package architecture, Gemini's `physvet`, Qwen's framing — all are cited but not provided. The coding agent must reconstruct them from context. If the agent has not read them, it will invent its own interpretations, which may not match the synthesis the user has in mind.

2. **It leaves seven critical design decisions entirely to the coding agent without suggested defaults.** For each, the agent can go in multiple directions, and the prompt does not say which direction is preferred. This will produce plans that are internally consistent but mutually incompatible across agents.

3. **It has four internal tensions** that will force the agent to make judgment calls the prompt does not authorize. The agent will resolve them one way; a reviewer will resolve them another.

The prompt is a good brief for an **expert** who has read the four plans and understands the synthesis. It is not a good brief for a coding agent that needs to produce a consistent plan without reinventing the architecture.

---

## Part II: Strengths

| Strength | Why It Matters |
|---|---|
| **§2 Non-negotiable decisions** | Prevents the agent from reopening settled questions. Good. |
| **§17 Vertical slice requirement** | Forces the agent to commit to a concrete end-to-end test. Good. |
| **§19 Explicit non-goals** | Prevents scope creep. Good. |
| **§20 Deliverable format (A–N)** | Forces structure. Good. |
| **§21 Discipline rule** | Gives the agent a scope criterion. Good. |
| **The provenance/corpus status distinction (§6)** | Separates two axes that are easy to conflate. Good. |
| **The candidate/corpus distinction (§9)** | Preserves the human-review boundary. Good. |

The prompt's instincts are correct. The execution has gaps.

---

## Part III: The Seven Critical Gaps

For each gap, I identify: what is missing, why the coding agent will get it wrong without guidance, and what the prompt should add.

### Gap 1 — The machine-readable corpus format

**What the prompt says (§6):**

> Do not make comments the sole canonical metadata source. Prefer a machine-readable representation that can be validated by tests.

**What is missing:** What is the machine-readable representation? JSON sidecars? Go struct tags? A generated manifest? A separate corpus file?

**Why the agent will get it wrong:** The prompt forbids comments but does not specify an alternative. A coding agent will default to the simplest option: either keep comments (violating the prompt) or invent a JSON format ad hoc. Both are wrong.

**What the prompt should add:** A concrete default. For example:

> Unless the agent has a strong reason otherwise, use a `corpus.json` file per package, containing the machine-readable metadata for each exported symbol. The Go doc comments may mirror the JSON, but the JSON is canonical.

### Gap 2 — The `Identify` gate at the code level

**What the prompt says (§2B):**

> `Identify()` is the explicit gate for a physical identification/insight and must be recorded in provenance.

**What is missing:** What does `Identify` take as input? What does it return? Does it need a justification string? Does it modify the objects or create a new one?

**Why the agent will get it wrong:** The prompt defines `Identify` conceptually but not operationally. The coding agent will guess. If the guess differs from the reviewer's expectation, the plan fails.

**What the prompt should add:** A signature. For example:

```go
// Identify asserts that two physics objects represent the same physical quantity.
// It records the identification in the ledger with the supplied justification.
// It is the only operation that can introduce a physical insight.
func Identify(a, b Object, justification string) (Object, error)
```

The exact signature is less important than having one.

### Gap 3 — Assumption propagation

**What the prompt says (§5):**

> result assumptions = required operation assumptions + input assumptions + newly introduced conditions

**What is missing:** Is `+` a union, a filter, or a compatibility check? What happens when two assumptions conflict? Does the operation fail, or does it record the conflict?

**Why the agent will get it wrong:** The formula suggests a union, but assumption conflicts are common in physics (e.g., "small angle approximation" vs. "arbitrary angle"). If the union is naive, assumptions accumulate without check. If the check is too strict, valid derivations fail.

**What the prompt should add:** A policy. For example:

> Assumptions are propagated by union. When two assumptions have the same `Kind` and different `Value`, the operation fails unless one is a specialization of the other. The failure must be surfaced to the caller.

### Gap 4 — Derivation ledger enforcement

**What the prompt says (§14):**

> The AI must not be able to hand-author a trusted derivation artifact that bypasses the ledger.

**What is missing:** How is this enforced? Via API design? Via session state? Via cryptographic signing? Via a separate build step?

**Why the agent will get it wrong:** The prompt states a requirement but not a mechanism. A coding agent will default to "the API doesn't expose a way to construct a derivation directly," which is fine but under-specified.

**What the prompt should add:** A mechanism. For example:

> Derivations are produced only by the `Ledger` type. The `Ledger` is initialized by a session, records each operation as it happens, and can produce a `Derivation` artifact only at the end. There is no public constructor for `Derivation`.

### Gap 5 — The `ResearchCandidate` schema

**What the prompt says (§16):**

> The natural terminal artifact is a `ResearchCandidate` or equivalent. It should contain only the information necessary for a human researcher to inspect the proposal, such as: hypothesis, premises, assumptions, derivation, provenance, ...

**What is missing:** The schema. What fields? What types? What format?

**Why the agent will get it wrong:** The list is descriptive, not prescriptive. Different agents will produce different schemas.

**What the prompt should add:** A minimal schema. For example:

```go
type ResearchCandidate struct {
    Hypothesis        Object
    Premises          []Object
    Assumptions       []Assumption
    Derivation        Derivation
    Provenance        Provenance
    FrameworkDeps     []string
    Limits            []Limit
    Predictions       []Prediction
    Falsification     []Condition
    ReviewHistory     []Review
}
```

### Gap 6 — MRC failure behavior

**What the prompt says (§20E):**

> For each MRC responsibility, state: mechanism, enforcement location, failure behavior.

**What is missing:** What is failure behavior? Panic? Return error? Log and continue? Which failures are hard (block execution) and which are soft (warn and continue)?

**Why the agent will get it wrong:** The prompt asks for "failure behavior" without defining what it means. Different agents will produce different interpretations.

**What the prompt should add:** A default. For example:

> MRC failures are returned as Go errors. The caller decides whether to abort the derivation. The library does not panic. Soft warnings are surfaced via a `Diagnostics` field on the result.

### Gap 7 — The MVP package scope

**What the prompt says (§3C):**

> Use a small shared core plus domain packages such as: core, mechanics, electromagnetism, relativity, qm, qft, statmech.

**What is missing:** Which packages are in the MVP? The phrase "such as" suggests the list is illustrative, but §21 says to defer anything not in the vertical slice. If the vertical slice is E=mc², then `relativity` is required, `mechanics` is useful, and `qm`, `qft`, `statmech` are not.

**Why the agent will get it wrong:** The agent might implement all listed packages (over-engineering) or only `core` (under-engineering).

**What the prompt should add:** An explicit MVP scope. For example:

> The MVP includes `core`, `operations`, `mechanics`, `relativity`, and `hypothesis`. All other domain packages are deferred to v0.5.

---

## Part IV: The Four Internal Tensions

### Tension 1 — "Non-negotiable" vs. "decide explicitly"

**§2 says:**

> Treat the following as decisions, not questions to reopen.

**§2D says:**

> For MVP, decide explicitly whether physvet is needed now.

The second sentence reopens a decision the first sentence says not to reopen. The agent will be confused: is `physvet` decided or not?

**Fix:** Either decide `physvet` in the prompt ("physvet is deferred to v0.5"), or explicitly mark it as the one exception ("all decisions are settled except physvet, which the agent must decide").

### Tension 2 — "Smallest coherent MVP" vs. the required concepts list

**The prompt repeatedly says "smallest coherent MVP."**

**§3 lists nine required concepts:** Physics Object, Dimension, Assumption, Convention, Symbolic Expression, Provenance, Derivation, Knowledge/Corpus Status, Candidate/Hypothesis Object.

**§4 lists six MRC requirements.**

**§5 lists seven assumption kinds.**

**§6 lists two status axes.**

That is not small. The agent will struggle to reconcile "smallest" with the requirement list.

**Fix:** Acknowledge the tension. Say: "The required concepts are the minimum for the vertical slice. Do not add concepts beyond this list. Each concept should be implemented minimally; the plan should specify what the minimum is."

### Tension 3 — "No math ontology" vs. "mathematical machinery may exist internally"

**§3 says:**

> Do not build a general-purpose mathematics ontology.

**§3 also says:**

> Mathematical machinery may exist internally.

The boundary is unclear. What counts as "internal machinery"? A symbolic expression type? A tensor index system? A dimension algebra?

**Fix:** Give examples. For example: "The internal symbolic expression type is required. The tensor index system is deferred. The dimension algebra is required."

### Tension 4 — The E=mc² vertical slice requires more than the prompt acknowledges

**§17 says:**

> Use a canonical derivation such as mass-energy equivalence as the principal integration test, while also demonstrating at least one classical relation.

**What is unacknowledged:** E=mc² is derived from the energy-momentum relation E² = (pc)² + (mc²)². This requires:
- Four-momentum (special relativity)
- The Minkowski metric
- Lorentz invariance
- The rest-frame condition (p=0)

That is not a small vertical slice. The agent might underestimate the scope.

**Fix:** Acknowledge the scope. Say: "The E=mc² vertical slice requires the four-momentum, the Minkowski metric, and the rest-frame condition. The MVP must include these. If the scope is too large, use a simpler vertical slice such as F=ma and defer E=mc² to v0.5."

---

## Part V: What's Missing

### Missing 1 — The six criteria are referenced but not listed

**§21 says:** "the six project criteria"

**What the agent needs:** A list of the six criteria. Without them, the discipline rule cannot be applied.

**Fix:** Inline the six criteria. For example:

> The six project criteria are: (1) fundamental primitives with MRC, (2) framework packages as eval targets with stated assumptions, (3) library as machine-readable corpus, (4) collective gap identification, (5) no truth adjudication, (6) hypothesis formulation support.

### Missing 2 — The referenced plans are not included

**§2 references Plan Z, DeepSeek, Gemini, Qwen.** The prompt assumes the agent has read them. If the agent has not, it will invent its own interpretations.

**Fix:** Either include the plans as appendices, or summarize the key points of each plan in the prompt itself. For example:

> Plan Z's key contributions: the Einstein test, the `Identify` gate, provenance, draft/commit. DeepSeek's key contribution: the package architecture. Gemini's key contribution: `physvet`. Qwen's key contribution: the framing. These contributions are incorporated in the decisions above.

### Missing 3 — The testing strategy

**§20M asks for "acceptance criteria."** But the prompt does not say how tests are organized. Table-driven tests? Golden files? Property-based tests?

**Fix:** Suggest a default. For example:

> Tests should be table-driven for operations and golden-file-based for derivations. The canonical derivations should have golden files that any change must reproduce.

### Missing 4 — The multi-agent workflow

**§15 mentions the adversarial reviewer.** But the prompt does not specify how the formulator and reviewer communicate.

**Fix:** Specify the artifact boundary. For example:

> The formulator produces a `ResearchCandidate` artifact (JSON). The reviewer consumes it and produces a `Review` artifact (JSON). The formulator can consume the `Review` and revise. The library does not orchestrate this; it only defines the artifact schemas.

### Missing 5 — The bootstrap order is suggested but not required

**§18 says:** "Choose the smallest credible implementation order. A likely sequence is: core → canonicalization → assumptions → operations → mechanics → canonical derivation → relativity → hypothesis → review."

The word "likely" suggests the sequence is a suggestion, not a requirement. The agent might choose differently.

**Fix:** Make the sequence a requirement unless the agent has a strong reason to deviate. For example:

> The implementation order is: (1) core representation, (2) canonicalization and exactness, (3) assumptions and provenance, (4) operations, (5) mechanics, (6) first canonical derivation (F=ma), (7) relativity, (8) second canonical derivation (E=mc²), (9) hypothesis model, (10) review artifact. Deviate only if a dependency requires it, and state the reason.

---

## Part VI: Recommendations

The prompt is good. It needs eight specific fixes before it will reliably produce a good plan:

1. **Inline the six criteria.** The agent cannot apply the discipline rule without them.
2. **Define the machine-readable corpus format.** Default to a `corpus.json` per package.
3. **Define the `Identify` signature.** At minimum, specify what it takes and returns.
4. **Define the assumption propagation policy.** Union with conflict detection is a reasonable default.
5. **Define the derivation ledger enforcement.** The `Ledger` is the only public constructor.
6. **Define the `ResearchCandidate` schema.** Even a minimal one is better than none.
7. **Define MRC failure behavior.** Errors, not panics; soft warnings via diagnostics.
8. **Resolve the four tensions.** Decide `physvet`; acknowledge the scope of the requirement list; clarify the math ontology boundary; acknowledge the scope of the E=mc² vertical slice.

With these fixes, the prompt becomes a reliable brief. Without them, it depends on the coding agent's competence and alignment with the reviewer's expectations.

---

## Part VII: Verdict

**The prompt is a good brief for an expert. It is not yet a good brief for a coding agent.**

The prompt's structure is correct. Its scope discipline is correct. Its non-goals are correct. Its deliverable format is correct. But it leaves too many critical decisions to the agent's judgment, and it references documents the agent may not have.

The user should either:

**Option A:** Add the eight fixes and resubmit the prompt.

**Option B:** Attach the four plans (Z, DeepSeek, Gemini, Qwen) and the six criteria to the prompt as appendices.

**Option C:** Accept that the prompt will produce a plan whose quality depends on the agent's competence, and use the four-agent review to catch inconsistencies.

Option A is the safest. Option B is the most thorough. Option C is the fastest but riskiest.

---

## The One-Paragraph Summary

> The prompt is a well-structured brief with correct instincts. It specifies the mission, the constraints, the non-goals, and the deliverable format. But it has seven critical gaps: the machine-readable corpus format is unspecified, the `Identify` gate is not defined at the code level, assumption propagation is sketched but not specified, derivation ledger enforcement is stated but not mechanized, the `ResearchCandidate` schema is descriptive not prescriptive, MRC failure behavior is undefined, and the MVP package scope is ambiguous. It has four internal tensions: "non-negotiable" vs. "decide explicitly," "smallest MVP" vs. the required concepts list, "no math ontology" vs. "internal math machinery," and the E=mc² vertical slice requiring more scope than the prompt acknowledges. It has five missing elements: the six criteria are referenced but not listed, the four plans are referenced but not included, the testing strategy is unspecified, the multi-agent workflow is unspecified, and the bootstrap order is suggested rather than required. The prompt should be fixed with eight specific changes before it is given to a coding agent: inline the six criteria, define the corpus format, define `Identify`, define assumption propagation, define ledger enforcement, define the `ResearchCandidate` schema, define MRC failure behavior, and resolve the four tensions. With these fixes, the prompt will reliably produce a consistent, high-quality plan. Without them, the prompt depends on the coding agent's competence and alignment with the reviewer's expectations, which is too risky for a multi-agent review. The prompt is 80% there. The remaining 20% is what determines whether the plan is good.
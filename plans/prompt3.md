You are the senior implementation architect for the Physics Compiler MVP.

Your task is **NOT to write code**. Produce a concrete, implementation-ready plan for the smallest coherent MVP based strictly on the consolidated architecture below.

The plan will be independently reviewed by four AI agents. Therefore, be precise, internally consistent, conservative in scope, and explicit about what is deliberately deferred.

## 1. Mission

Build the minimum viable Go library that gives an AI agent a **formal pen-and-paper substrate for symbolic physics reasoning**.

The AI is the theorist.

The library/compiler provides:

* fundamental physical primitives
* symbolic transformation
* compiler/library-level Math Reality Check (MRC)
* assumptions and constraints
* provenance and derivation tracking
* a machine-readable physics corpus
* explicit handling of physical identification/insight
* a controlled mechanism for proposing novel physics concepts

The system does **not** determine whether a hypothesis is true in nature.

The system's purpose is to enable AI agents to formulate and inspect hypotheses about nature, including possible theory-bridging ideas, while leaving empirical/scientific adjudication to humans.

## 2. Core architectural decisions already made

Treat the following as decisions, not questions to reopen:

### A. Go is the source language

Do not invent a new physics programming language.

Use ordinary Go, Go packages, Go types, interfaces, compiler tooling, and tests.

### B. Plan Z is the conceptual spine

Retain:

* explicit `Postulate`
* `Declare`
* `Define`
* `Step`
* `Identify`
* `Conclude`
* derivation provenance
* draft/commit distinction where useful
* canonical derivation records

The critical distinction is:

**mechanical symbolic manipulation != physical insight**

`Simplify()` must not silently perform a physical identification.

`Identify()` is the explicit gate for a physical identification/insight and must be recorded in provenance.

### C. DeepSeek's package architecture is the corpus structure

Use a small shared core plus domain packages such as:

```text
core
mechanics
electromagnetism
relativity
qm
qft
statmech
```

Do not attempt to populate all of theoretical physics for the MVP.

### D. MRC has one semantic source of truth

The library owns the MRC semantics.

A future/static-analysis tool such as `physvet` may consume those semantics, but must not become a second independent physics rule engine.

For MVP, decide explicitly whether `physvet` is needed now. If not required for the minimal vertical slice, defer it rather than duplicating MRC logic.

### E. MRC itself is fallible

The current MRC rules are not metaphysical truths.

The architecture must permit future human-validated exceptions and revisions without silently weakening existing rules.

An exception must be:

* explicit
* scoped
* versioned
* provenance-bearing
* human-approved
* subject to documented falsification/evidence requirements

Do not design an unrestricted "override" mechanism.

### F. Physical truth is outside compiler authority

The system may report things such as:

```text
FORMALLY_VALID
CATEGORY_ERROR
DIMENSION_ERROR
ASSUMPTION_CONFLICT
UNRESOLVED
```

It must never produce:

```text
PHYSICALLY_TRUE
THEORY_CONFIRMED
TRUTH_SCORE
```

Corpus status is curated input/data, not something the compiler computes.

## 3. Required core concepts

The MVP must determine the minimum viable representation for:

```text
Physics Object
Dimension
Assumption
Convention
Symbolic Expression
Provenance
Derivation
Knowledge/Corpus Status
Candidate/Hypothesis Object
```

Use the smallest model that can support the required end-to-end examples.

Do not build a general-purpose mathematics ontology.

Do not expose a giant taxonomy of mathematical structures unless the MVP demonstrably requires one.

Mathematical machinery may exist internally.

## 4. MRC requirements

MRC must operate on the following broad classes of constraints:

1. physical/category compatibility
2. dimensional consistency
3. operation applicability
4. assumption/constraint compatibility
5. explicit physical identification
6. provenance/derivation integrity

The plan must state **which checks are enforced by Go's type system, which are enforced by runtime/library contracts, and which are deferred to static analysis tooling**.

Do not create a giant centralized jurisdiction matrix.

Prefer typed constructors, interfaces, operation contracts, explicit assumptions, and canonical structures.

## 5. Assumption system

Assumptions are first-class data, not comments.

The plan must specify the minimum model and propagation behavior for:

```text
domain
regime
constraint
convention
approximation/truncation
mathematical precondition
physical assumption
```

Operations must propagate assumptions and surface contradictions.

Example principle:

```text
result assumptions =
    required operation assumptions
    + input assumptions
    + newly introduced conditions
```

Conflicts must not disappear silently.

Do not over-engineer a theorem prover.

## 6. Corpus model

Physics packages are part of the AI's machine-readable physics corpus.

The MVP must define the minimum machine-readable metadata for a physics artifact, including where applicable:

```text
definition
framework
assumptions
domain/regime
provenance/source
derivable-from
known limits/reductions
known failure/anomaly
corpus status
```

Distinguish:

### Provenance status

For example:

```text
DEFINED
POSTULATED
DERIVED
IDENTIFIED
APPROXIMATED
HYPOTHESIS
```

### Human-curated corpus status

For example:

```text
ESTABLISHED
CONTESTED
SUPERSEDED
FALSIFIED
```

The library may carry and render corpus status but must not infer it.

Do not make comments the sole canonical metadata source.

Prefer a machine-readable representation that can be validated by tests.

## 7. Frameworks such as relativity

A framework such as relativity is **not part of the metaphysical core**.

Treat it as a scoped physics framework/corpus package.

The MVP must allow a framework package to serve as an AI evaluation target through:

* known derivations
* explicit assumptions
* domain/regime boundaries
* known limits/reductions
* known anomalies/failures
* distinguishing predictions where practical

Do not build an enormous evaluation framework.

Demonstrate the concept with a very small canonical suite.

## 8. Common equations

The MVP must demonstrate that ordinary foundational physics relations can live in the corpus independently of treating any single modern theory as "final."

At minimum, plan for a small set of canonical examples such as:

```text
F = ma
momentum = mv
kinetic energy
mass-energy relation E = mc²
```

Use these as corpus/derivation examples.

Do not attempt to encode every textbook equation.

## 9. Hypothesis mechanism

The AI must be able to introduce a genuinely new candidate physical concept.

However:

**open representation does not imply open authority.**

Candidate objects must be distinguishable from trusted corpus objects.

Candidate-derived results must remain marked as candidate/provisional and must not be able to manufacture trusted corpus provenance.

Promotion into the trusted corpus must be a human-governed operation outside ordinary AI execution.

Keep this mechanism minimal.

Do not build a full theory-management platform.

## 10. Popper/falsifiability support

The MVP must support the representation of:

```text
candidate hypothesis
assumptions
predicted consequence
distinguishing prediction
potential falsification condition
recovery/reduction to an existing framework
```

Do not claim that the compiler has falsified or confirmed a theory.

The compiler records and checks the formal structure of these claims.

Humans/empirical reality determine their scientific status.

## 11. Anomalies and failed theories

The corpus should not contain only successful theories.

Plan a minimal representation for known anomalies/failure cases because these provide important AI research targets.

However, keep the MVP small: implement only enough structure to demonstrate that an AI can inspect an anomaly, identify the relevant framework/assumptions, and formulate a candidate response.

## 12. Symbolic engine

Define the minimum internal symbolic representation required for the MVP.

At minimum, determine support for the operations necessary for the canonical vertical slice:

```text
construct
add/subtract
multiply/divide
power
substitute
differentiate
integrate where genuinely required
simplify
limit/reduction
identify
```

Do not build a complete CAS.

Do not implement numerical simulation, numerical optimization, Monte Carlo, GPU execution, or general scientific computing.

The goal is symbolic derivation.

## 13. Canonicalization and exactness

Retain the existing mechanical safeguards:

* immutable/controlled physics objects
* constructor-controlled authoritative metadata
* exact rational representation where required
* canonical symbolic representation
* deterministic structural equality
* deterministic serialization
* canonical hashing
* explicit convention metadata
* convention conflict detection

Do not replace these with string comparison or informal normalization.

## 14. Derivation ledger

The derivation ledger is a core MVP artifact.

Every committed derivation must preserve enough information to reconstruct:

```text
premises
assumptions
operations
intermediate results
identifications
approximations
conclusion
provenance
```

The AI must not be able to hand-author a trusted derivation artifact that bypasses the ledger.

The session/ledger should be the authority for trusted derivation provenance.

## 15. Adversarial review boundary

The MVP should produce a machine-readable derivation/research artifact that a separate reviewer AI can inspect.

Define the minimum artifact and review structures.

The reviewer should be able to identify issues such as:

```text
category error
dimension error
hidden assumption
assumption conflict
unsupported identification
invalid reduction
overclaim
circular reasoning
provenance problem
```

Do not build an autonomous truth adjudicator.

## 16. Human handoff

The natural terminal artifact is a `ResearchCandidate` or equivalent.

It should contain only the information necessary for a human researcher to inspect the proposal, such as:

```text
hypothesis
premises
assumptions
derivation
provenance
known framework dependencies
limits/reductions
distinguishing predictions
potential falsification conditions
review/challenge history
```

The library stops there.

It does not decide whether the candidate is scientifically true.

## 17. MVP vertical slice

The implementation plan MUST identify one end-to-end vertical slice and make that the primary acceptance target.

Prefer a minimal sequence such as:

```text
AI reads corpus
   ↓
constructs physics objects
   ↓
states assumptions
   ↓
performs symbolic derivation
   ↓
MRC checks each step
   ↓
explicit Identify where physical insight is required
   ↓
derivation ledger is produced
   ↓
canonical result is generated
   ↓
candidate/research artifact can be reviewed
```

Use a canonical derivation such as mass-energy equivalence as the principal integration test, while also demonstrating at least one classical relation.

The plan should explain exactly why the selected vertical slice exercises the architecture rather than merely proving that arithmetic works.

## 18. Bootstrap order

Choose the smallest credible implementation order.

A likely sequence is:

```text
core representation
→ canonicalization/exactness
→ assumptions/provenance
→ operations
→ mechanics
→ canonical derivation
→ relativity
→ hypothesis/research artifact
→ review artifact
```

Adjust only when a dependency requires it.

Do not implement all domain packages before proving the core vertical slice.

## 19. Explicit non-goals for MVP

The implementation plan must explicitly exclude:

* new programming language
* new parser/frontend
* general mathematics ontology
* full CAS
* numerical physics engine
* simulation framework
* automated empirical validation
* autonomous truth adjudication
* truth scores
* comprehensive physics corpus
* full theorem prover
* unrestricted exception bypass
* large multi-agent orchestration system
* integration with EBP 2.1 itself
* automatic promotion of hypotheses into trusted physics

## 20. Deliverable format

Return an **implementation plan only**, not code.

Structure it as:

### A. MVP objective

One paragraph.

### B. Architectural boundary

What belongs in the MVP and what explicitly does not.

### C. Package/module layout

Show the proposed Go package tree.

### D. Core data model

List the minimal core types and their responsibilities.

### E. MRC enforcement model

For each MRC responsibility, state:

```text
mechanism
enforcement location
failure behavior
```

### F. Symbolic operation set

Only operations required for the MVP vertical slice.

### G. Physics corpus structure

How mechanics and relativity are represented and how machine-readable semantics are exposed.

### H. Assumptions/provenance/status model

Keep the two status axes separate.

### I. Candidate/hypothesis model

Explain how novel concepts are represented without bypassing trusted authority.

### J. Derivation ledger and research artifact

Define the minimum contents and ownership.

### K. Canonical MVP derivations/tests

Specify the exact small set of derivations that prove the architecture works.

### L. Implementation sequence

Ordered engineering steps with dependencies.

### M. Acceptance criteria

Concrete tests that must pass before calling the MVP complete.

### N. Explicit deferrals

List important capabilities intentionally postponed to v0.5+.

## 21. Discipline rule

Whenever a proposed component does not directly support the MVP vertical slice or one of the six project criteria, **defer it**.

Do not add infrastructure "for future flexibility" unless the MVP genuinely requires it.

The goal is not to build the final Physics Compiler.

The goal is to build the **smallest system that convincingly demonstrates the core thesis**:

> An AI can use a machine-readable physics corpus and a typed symbolic substrate to perform non-numerical physics derivations, while MRC constrains formal/category errors, physical insight remains explicit and reviewable, hypotheses remain provisional, and scientific truth remains outside the compiler's authority.

Produce the plan with enough technical specificity that another coding agent could implement it without having to reinvent the architecture.

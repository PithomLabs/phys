**Plan 4Z is the best baseline among the four options**, though it should be implemented as a synthesized hybrid combining Plan 4Z's conceptual spine, DeepSeek's package structure, Gemini's `physvet` static analysis, and an explicit Assumptions & Knowledge-Status layer. Plan 4Z wins because it uses the **Einstein derivation as a design test**, distinguishing mechanical symbolic transformations from creative physical insights while making those insights explicit through reviewable ledgers.

---

### Evaluation Across the 6 Criteria

| Criterion | Plan 4 DeepSeek | Plan 4 Gemini | Plan 4 Qwen | **Plan 4Z (Best Baseline)**<br> |
| --- | --- | --- | --- | --- |
| **1. Primitives & Compiler-Level MRC** | Moderate typed construction; delegates math engine internally.

 | **Strongest mechanism**: Introduces `physvet` static analysis tool before AST symbolic execution.

 | Conceptually strong, but lacks static compiler tooling details.

 | **Strongest overall architecture**: Defines 5 enforcement levels (typed construction, dimensions, contracts, `Identify`, provenance).

 |
| **2. Theories as Derivation Eval Targets (Popperian Falsifiability)** | Teaches theories as standard library packages, but lacks explicit regime boundaries.

 | Focuses heavily on formal proofs, risking treating valid proofs as absolute physical truth.

 | Weak implementation detail; uses epistemically overstrong claims.

 | **Strongest fit**: Explicitly scopes frameworks (`Postulate`, `Declare`, `Step`) with domain limitations.

 |
| **3. Packages as Machine-Readable AI Corpus** | **Strongest package structure**: Clean domain separation (`mechanics`, `relativity`, `qm`, `qft`, `statmech`).

 | Moderate; focuses on verification workflows rather than semantic corpus features.

 | Weak; lacks clear Go package and API layouts.

 | Moderate package layout, but strong metadata tagging for AI reading.

 |
| **5. Non-Adjudication of Physical Truth** | **Strong**: Focuses on symbolic equivalence rather than empirical verification.

 | Some overreach: "Physical sanity" terminology risks blurring validity with empirical reality.

 | Overreach: Claims "absolute epistemic safety" against false theories.

 | **Strongest**: Explicitly maintains that formal mathematical validity $\neq$ physical truth.

 |
| **6. AI Hypothesis Formulation & Research Program** | Good serialization for external agents, but lacks explicit hypothesis primitives.

 | Good verification steps, but restricts agents to existing library objects.

 | Strong conceptually (Program A = pen & paper), but lacks concrete primitives.

 | **Strongest foundation**: `Identify` gate, draft/commit ledgers, and provenance tracing.

 |

---

### Detailed Analysis by Criterion

#### 1. Fundamental Primitives & Compiler-Level Math Reality Check (MRC)

* **Plan 4Z** establishes a 5-tier enforcement model: typed construction, dimensional checking, operation contracts, explicit `Identify` gates, and provenance closure. Instead of viewing MRC as a generic Go interface check, Z recognizes that physical operations (e.g., $E + T$) require explicit physical bridges (e.g., $k_B \cdot T$) rather than raw type compatibility.


* **Gemini** contributes a crucial static analysis tool (`physvet`), intercepting Go code *before* symbolic execution to statically verify dimensional vector flow $[M, L, T, Q, \Theta]$, category bridges, and index variance.



#### 2. Relativity as an Evaluation Target (Popperian Falsifiability)

* Per Karl Popper, scientific theories like Special Relativity are not metaphysical axioms, but constrained, falsifiable frameworks.


* **Plan 4Z** treats `relativity` as an explicitly scoped physical framework with defined postulates and domain boundaries rather than an unassailable ontology. This allows an AI agent to use `relativity` as a target derivation trace to verify if specific conclusions hold under stated assumptions.



#### 3. Library Packages as the AI's Physics Corpus

* For AI agents to "know" and "understand" physics, Go packages must serve as machine-readable knowledge artifacts rather than standard API function signatures.


* **DeepSeek** delivers the best Go package layout (`core`, `mechanics`, `relativity`, `qm`, `qft`, `statmech`). When enriched with Z's semantic metadata (definitions, prerequisites, domain regimes, known limits), the AI agent can read and reason across the corpus as a structured knowledge network.



#### 5. Non-Adjudication of Physical Truth

* The library must evaluate **formal structural validity**, never **empirical physical truth**.


* **Plan 4Z** excels here by explicitly isolating symbolic mechanics from physical interpretation. It prohibits the compiler from auto-evaluating nature's validity, rejecting **Qwen's** "absolute epistemic safety" claims and **Gemini's** "physical sanity" overreach.



#### 6. AI Hypothesis Formulation & Research Program Facilitation

* To discover novel bridges (such as Theory of Everything couplings under Evidence-Based Profiling frameworks), the system must allow AI agents to propose new physical constructs, perform formal transformations, and generate derivation traces.


* **Plan 4Z** separates mechanical simplification (`Simplify()`) from conceptual breakthroughs (`Identify()`), ensuring that any theoretical leap proposed by an AI agent is flagged in an immutable derivation ledger for human review.



---

### 4. Collective Gaps & Limitations Across All Four Plans

Despite their individual strengths, all four plans collectively miss five critical architectural requirements:

1. **Lack of a First-Class Assumption & Regime System**: None of the plans adequately model mathematical and physical constraints (e.g., $x \ge 0$, $[A, B] = 0$, $v \ll c$, flat spacetime). Equations cannot exist as bare formulas; they must carry explicit premises, domains of validity, and approximation orders.


2. **Missing Knowledge-Status Metadata**: The plans fail to distinguish between different classes of physical statements. The library must tag artifacts with explicit epistemic statuses: `Definition`, `Postulate`, `Law`, `Identity`, `DerivedRelation`, `Approximation`, `ModelAssumption`, or `Hypothesis`.


3. **Closed Conceptual Vocabularies**: Existing plans restrict AI agents to choosing pre-defined library objects (`Mass`, `Energy`, `Metric`). To discover new physics or unification bridges, the library must feature a `hypothesis.NewObject()` mechanism allowing agents to declare candidate primitives without polluting the standard corpus.


4. **Lack of Machine-Readable Package Semantics**: Packages currently act like standard code APIs. To function as a true corpus, package items must expose structured metadata detailing physical definitions, prerequisites, counterexamples, and limit behaviors.


5. **Undefined Human-Handoff Artifact**: The plans do not define where the AI agent's execution terminates. The engine requires a structured `ResearchCandidate` export (containing premises, derivation trace, limiting cases, and predicted observables) that halts execution for human evaluation, explicitly preventing autonomous "truth scoring".



---

### Final Architecture Recommendation (v0.4 Blueprint)

Construct the **v0.4 Physics Compiler** by combining the best elements of the evaluated plans:

$$\text{v0.4 Architecture} = \underbrace{\text{Plan 4Z Spine}}_{\text{Einstein Test, Identify Gate, Ledger}} + \underbrace{\text{DeepSeek Packages}}_{\text{Domain Layout \& Core Engine}} + \underbrace{\text{Gemini } \mathtt{physvet}}_{\text{Static MRC Pass}} + \underbrace{\text{First-Class Assumptions \& Knowledge Status}}_{\text{Collective Gap Fixes}}$$

```text
                  HUMAN RESEARCHER (Empirical Evaluation)
                                    ▲
                                    │ (Reviewable ResearchCandidate)
                             AI THEORIST AGENT
                                    │
                       Reads Machine-Readable Corpus
                                    │
                           Emits Ordinary Go Code
                                    │
                                    ▼
                       ┌─────────────────────────┐
                       │   physvet Static Pass   │ (MRC Pass: Dimensions, Index Variance)
                       └────────────┬────────────┘
                                    │
                                    ▼
                       ┌─────────────────────────┐
                       │   Core Physics Library  │ (Primitives, Bridges, Operators)
                       └────────────┬────────────┘
                                    │
                                    ▼
                       ┌─────────────────────────┐
                       │   Physics Domain Corpus │ (Mechanics, Relativity, QM, QFT)
                       └────────────┬────────────┘
                                    │
                                    ▼
                       ┌─────────────────────────┐
                       │    Derivation Engine    │ (Explicit Identify() & Step Ledgers)
                       └─────────────────────────┘

```

By enforcing MRC statically, scoping theories falsiably, encoding physics semantically, and exporting structured derivation traces, the library enables AI agents to formulate meaningful theoretical hypotheses while keeping empirical adjudication strictly in human hands.
Based on the six criteria provided, **Plan Z is the strongest conceptual foundation**, but it cannot stand alone. It must be synthesized with **DeepSeek’s package architecture** and **Gemini’s static analysis mechanism** to create the definitive v0.4 Physics Compiler. 

Here is the evaluation of the plans against your criteria, the identification of their collective blind spots, and the blueprint for the unified architecture.

---

### 1. The Verdict: Why Plan Z is the Spine
**Plan Z** is the only plan that correctly identifies the core epistemological boundary of theoretical physics: **the distinction between mechanical symbolic manipulation and physical insight.** 

By introducing the `Identify` gate (separating it from `Simplify`), Plan Z ensures that the compiler never silently conflates a mathematical equivalence with a physical identification (e.g., $E=mc^2$ is not just an algebraic simplification; it is a physical identification of mass and energy). Furthermore, Plan Z strictly adheres to **Criterion 5**, refusing to act as a "truth machine," whereas Plan Qwen dangerously overreaches by claiming "absolute epistemic safety."

However, Plan Z lacks a robust package structure for the AI to read (Criterion 3), and Plan DeepSeek lacks the rigorous MRC enforcement mechanisms. Therefore, the final architecture uses **Z for the epistemological engine**, **DeepSeek for the physics corpus**, and **Gemini for the MRC enforcement layer**.

---

### 2. Evaluation Against the Criteria

#### Criterion 1: Fundamental Primitives + MRC at the Compiler Level
*   **The Solution:** We adopt **Plan Z’s** typed physical primitives (`Mass`, `Metric`, `State`) combined with **Plan Gemini’s `physvet`** static analysis layer. 
*   **How it works:** The AI writes standard Go code using physics primitives. Before the code compiles, `physvet` intercepts it to check dimensional consistency, tensor variance, non-commutativity, and category boundaries (e.g., preventing a `PartialOrder` from being passed where a `Temporal` duration is required). The MRC is enforced mechanically via Go interfaces, not by a hardcoded philosophical matrix.

#### Criterion 2: Relativity as a Popperian Eval Target (Not Absolute Truth)
*   **The Solution:** The library must treat `relativity` not as the ontology of nature, but as a **Formal Framework** with explicit boundaries.
*   **How it works:** A package like `relativity` is encoded with machine-readable metadata:
    ```go
    // Framework Metadata
    Status: EstablishedRelation
    Assumptions: [LorentzInvariance, EquivalencePrinciple, ContinuousManifold]
    Domain: Macroscopic, WeakToStrongGravity
    Limits: FailsAtSingularity, IncompatibleWithQuantumSuperposition
    ```
    The AI agent uses this package as an **evaluation target**. It can derive consequences *within* the framework, or it can propose a `CandidateTheory` that modifies these assumptions. The compiler checks if the candidate theory correctly `Recover()`s the Relativity framework in the appropriate IR (Infrared) limit, satisfying Popperian falsifiability.

#### Criterion 3: Packages as the AI Physics Corpus
*   **The Solution:** We adopt **Plan DeepSeek’s** domain-specific package structure (`mechanics`, `relativity`, `qm`, `qft`, `statmech`).
*   **How it works:** These packages are not just collections of Go functions; they are the AI’s "textbooks." When the AI needs to understand how a field propagates, it reads the `qft.Propagator` documentation and code structure. The packages encode the *current consensus of physical knowledge*, allowing the AI to "know" physics by interacting with the corpus, rather than relying on the LLM's latent space.

#### Criterion 5: No Truth Adjudication
*   **The Solution:** Strict adherence to **Plan Z’s** boundary. 
*   **How it works:** The compiler outputs `FORMALLY_VALID`, `CATEGORY_ERROR`, or `UNRESOLVED`. It **never** outputs `PHYSICALLY_TRUE`. The library guarantees that the AI's derivation is logically sound and mathematically typed; it explicitly refuses to judge whether nature actually behaves that way. That is the job of Program B (execution) and human empirical review.

#### Criterion 6: Ultimate Goal (AI Hypothesis Formulation for EBP 2.1)
*   **The Solution:** The library must support an **Open Hypothesis Vocabulary**.
*   **How it works:** The AI is not restricted to combining existing objects (like `Mass` and `Energy`). To find TOE bridges, the AI must be able to invent *new* physical concepts. The library provides a `CandidateObject` sandbox where the AI can declare a new construct (e.g., `PreGeometrySpin`), assign it provisional dimensions and transformation rules, and attempt to derive its coupling to the `relativity` corpus. If it survives the semantic checks, it becomes a `ResearchCandidate` artifact for human review.

---

### 3. The Collective Gaps (What All Four Plans Missed)

While the four plans provide excellent foundations, they collectively miss four critical architectural requirements necessary for a true AI research engine:

#### Gap 1: The Assumption and Constraint Engine
None of the plans adequately model the *assumptions* required for symbolic transformations. 
*   *Example:* $\sqrt{x^2} = x$ is only valid if $x \ge 0$. $AB = BA$ is only valid if the objects commute. 
*   *Fix:* Every symbolic operation must carry an `Assumption` tag. If the AI derives a result, the compiler must output the exact conditions under which that derivation holds. If the AI ignores an assumption, the compiler flags a `CONDITIONAL_VIOLATION`.

#### Gap 2: Machine-Readable Package Semantics
DeepSeek’s packages are great, but they are just Go code. To serve as an AI corpus, they need semantic tags.
*   *Fix:* Every exported function/object in the corpus must have a structured comment block that the AI can parse:
    ```go
    // @Status: EstablishedRelation
    // @DerivableFrom: [LorentzTransform, FourMomentum]
    // @ReducesTo: classical.KineticEnergy (when v << c)
    // @FalsificationCondition: Violates unitarity at Planck scale
    ```

#### Gap 3: The Falsifiability / Discriminator Structure
The plans allow the AI to derive equations, but they don't natively support expressing *competing predictions*.
*   *Fix:* The library needs a `Discriminator` primitive. The AI should be able to write:
    ```go
    candidate := theory.Propose(...)
    predictionA := candidate.Predict(observable)
    predictionB := relativity.Predict(observable)
    physics.AssertNotEqual(predictionA, predictionB) // Creates a testable discriminator
    ```

#### Gap 4: The "Recovery" and "Limit" Operations
To satisfy the Popperian requirement (Criterion 2), the library must natively support taking limits.
*   *Fix:* Introduce `physics.Recover(candidate, target, regime)` and `physics.Limit(expression, parameter, value)`. The AI must be able to formally prove that its new TOE bridge reduces to the `relativity` corpus in the macroscopic, low-energy limit.

---

### 4. The Unified v0.4 Architecture

To realize the Reality-First research program (EBP 2.1), the final system looks like this:

```text
                 HUMAN RESEARCH INTENT (EBP 2.1 Goals)
                                   │
                                   ▼
                             AI THEORIST
                     (Reads Physics Corpus Packages)
                                   │
                                   ▼
                          ORDINARY GO SOURCE
                     (Uses Z's Postulate/Identify)
                                   │
                                   ▼
                     ┌────────────────────────┐
                     │       PHYSVET          │ <-- Gemini's MRC enforcement
                     │ (Dimensions, Tensors,  │
                     │  Category Boundaries)  │
                     └───────────┬────────────┘
                                 │
                                 ▼
                     ┌────────────────────────┐
                     │    PHYSICS CORE (Z)    │
                     │ (Symbolic Engine,      │
                     │  Assumption Tracker,   │
                     │  Provenance Ledger)    │
                     └───────────┬────────────┘
                                 │
                                 ▼
                     ┌────────────────────────┐
                     │   PHYSICS CORPUS       │ <-- DeepSeek's Packages
                     │ (mechanics, relativity,│    (With Semantic Metadata)
                     │  qm, qft, statmech)    │
                     └───────────┬────────────┘
                                 │
                                 ▼
                      SYMBOLIC DERIVATION
                      (Including Limits/Recovery)
                                 │
                                 ▼
                      DERIVATION LEDGER
                                 │
                     ┌───────────┴────────────┐
                     ▼                        ▼
              KNOWN RELATION             AI HYPOTHESIS
              (Corpus derivation)        (Candidate Sandbox)
                     │                        │
                     └───────────┬────────────┘
                                 ▼
                     RESEARCH CANDIDATE ARTIFACT
                     ├── Formal Validity: PASS
                     ├── Assumptions: [List]
                     ├── Provenance: [List]
                     ├── Falsification Conditions: [List]
                     └── Recovery Limits: [List]
                                 │
                                 ▼
                        HUMAN / EMPIRICAL REALITY
```

### Summary
By adopting **Plan Z's epistemological rigor**, **DeepSeek's corpus architecture**, and **Gemini's static analysis**, while patching the **Assumption Engine** and **Candidate Sandbox** gaps, you create a system that perfectly fulfills your criteria. It gives the AI a rigorous, machine-readable "pen and paper" to formulate TOE bridges, strictly enforces the Math Reality Check, treats established physics as Popperian frameworks rather than absolute truths, and definitively stops at the boundary of formal validity, leaving the adjudication of physical truth to reality and humans.
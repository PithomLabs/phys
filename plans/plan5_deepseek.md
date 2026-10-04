# WHICH PLAN IS BEST — AN INDEPENDENT ASSESSMENT
## Against the Six Criteria, With Where I Disagree With the Attached Review

**Status:** Independent evaluation. The attached review argues for Plan 4Z. I largely agree but disagree on three points and identify five gaps the review itself missed.

---

## Part I: The Verdict

**Plan 4Z is the best foundation.** The review is correct on this.

But the review overstates the case in three places and understates it in two. And the synthesis it proposes is too ambitious for where the user is (the user has repeatedly pushed for simplicity).

My verdict is:

> **Use Z's conceptual spine. Use DeepSeek's package architecture. Use Gemini's `physvet` as a narrow gate, not a layer. Reject Qwen's epistemic overreach. Add an assumptions/status model. But build the minimum, not the maximum.**

---

## Part II: Assessment Against the Six Criteria

### Criterion 1 — Fundamental primitives + MRC enforcement

| Plan | Assessment | Rating |
|---|---|---|
| **DeepSeek** | Provides the package structure (`core`, `mechanics`, `relativity`, `qm`, ...) but does not specify the primitives. MRC is implied by types but not designed. | **Medium** |
| **Gemini** | Strongest technical mechanism — `physvet` static analyzer with three passes (category/bridge, dimensions, tensor/noncommutativity). But the mechanism is a checker, not a primitive system. | **Medium-High** |
| **Qwen** | Frames the problem well (AI is theorist, compiler is pen-and-paper) but does not deliver the primitives. | **Low-Medium** |
| **Z** | Strongest conceptual model — five enforcement levels (typed construction, dimensional checking, operation contracts, explicit `Identify`, provenance/ledger closure). But the primitives are still abstract. | **High (conceptual)** |

**The gap none of them fill:** None of the four plans says what a physics object *is*, at the code level. They all describe *what the library should do* but not *what the library's core type is*.

**My addition:** The core type must be decided first. My recommendation:

```go
type Object struct {
    Name       string
    Role       Role
    Component  Component   // set by constructor only
    Dimension  Dimension
    Structure  Expr        // internal symbolic structure
    Assumptions []Assumption
    Provenance Provenance
}
```

Where `Component` and `Assumptions` are set by the constructor, not the AI. This is the primitive. MRC enforces through it.

### Criterion 2 — Library packages as eval targets with stated assumptions

| Plan | Assessment | Rating |
|---|---|---|
| **DeepSeek** | Has package structure but does not clearly distinguish "framework" from "truth." | **Medium** |
| **Gemini** | Tends to treat relativity as the ontology. Some overreach. | **Low-Medium** |
| **Qwen** | Explicitly says relativity should be a scoped framework, not the ontology. But does not architect it. | **Medium (conceptual)** |
| **Z** | Strongest — treats relativity as one explicitly scoped framework among many, with assumptions and provenance. | **High** |

**The gap none of them fully fill:** None of the four describes *how* a package declares itself as a framework. There is no "framework manifest."

**My addition:** Every package should ship with a `framework.json` (or equivalent) that declares:

```json
{
  "name": "special_relativity",
  "status": "established_relation",
  "assumptions": ["Minkowski_spacetime", "Lorentz_symmetry", "no_gravity"],
  "domain": "flat_spacetime, v << c excluded",
  "reduces_to": ["classical_mechanics"],
  "source": "Einstein 1905",
  "known_limits": ["does_not_include_gravity", "does_not_include_quantum"]
}
```

The AI reads this manifest as part of its corpus. This makes relativity a *refinement*, not a *truth*.

### Criterion 3 — Library packages as physics corpus

| Plan | Assessment | Rating |
|---|---|---|
| **DeepSeek** | Strongest package architecture. Domain packages (`mechanics`, `relativity`, `qm`, `qft`, `statmech`) organized like Go's standard library. | **High** |
| **Gemini** | Moderate. Focuses on static analysis, not corpus. | **Medium** |
| **Qwen** | Weak. Conceptual only. | **Low** |
| **Z** | Moderate. Has the epistemic model but not the package architecture. | **Medium** |

**The gap none of them fill:** The packages are code. The AI reads code. But code alone is not a corpus. The AI needs machine-readable *semantics* — definitions, assumptions, derivations, limits — not just function signatures.

**My addition:** Every exported type and function must carry a doc comment in a structured format that the AI can parse. Example:

```go
// EnergyMomentumRelation returns E² = (pc)² + (mc²)².
//
// Status: established_relation
// Framework: special_relativity
// Assumes: Minkowski_spacetime, Lorentz_symmetry
// Derives_from: four_momentum_norm
// Reduces_to: classical_kinetic_energy (v << c)
// Source: Einstein 1905, extended by Dirac 1928
func EnergyMomentumRelation(m Mass, p Momentum, c SpeedOfLight) Object
```

The AI parses these comments. They are the corpus.

### Criterion 4 — Collective gaps

The attached review identifies:
- Assumptions/regime system
- Knowledge-status/provenance system
- Open hypothesis vocabulary
- Machine-readable package semantics
- Human-handoff artifact

These are correct. I add five more the review missed:

**Gap A — The paper translation protocol.** None of the plans specifies *how* the AI translates a paper's notation to library objects. This is the AI's most common task. Without a protocol, the AI invents one ad hoc, leading to inconsistency.

**Gap B — The adversarial review protocol.** The v0.3 spec mentions an "adversarial AI boundary" but no plan specifies what the reviewer consumes or produces. The challenge/response schema is undefined.

**Gap C — The equality and canonicalization spec.** When are two physics objects equal? How are expressions canonicalized? None of the plans specifies this. Without it, the derivation trace is unreadable.

**Gap D — The symbolic expression type.** All four plans talk about "operations on physics objects" but none defines the internal `Expr` type. This is the engine. It cannot be hand-waved.

**Gap E — The multi-agent workflow.** The user mentions two AIs (formulator and adversarial reviewer). None of the plans specifies how they interact. What does the formulator write? What does the reviewer read? How does the reviewer's challenge feed back?

### Criterion 5 — No truth adjudication

| Plan | Assessment | Rating |
|---|---|---|
| **DeepSeek** | Explicitly rejects truth adjudication. | **High** |
| **Gemini** | Has "physical sanity" language that risks overreach. | **Medium** |
| **Qwen** | Claims "absolute epistemic safety" — a clear overreach. | **Low** |
| **Z** | Strongest. Distinguishes formal validity from physical truth explicitly. | **High** |

The review is correct to reject Qwen's language and to narrow Gemini's. Z wins here.

### Criterion 6 — AI hypothesis formulation

| Plan | Assessment | Rating |
|---|---|---|
| **DeepSeek** | Has the package architecture but no hypothesis layer. | **Medium** |
| **Gemini** | Focuses on validation, not formulation. | **Low-Medium** |
| **Qwen** | Strongest conceptual framing (AI = theorist). | **Medium-High (conceptual)** |
| **Z** | Strongest architectural foundation (`Declare`, `Hypothesis`, draft/commit). | **High** |

**The gap none of them fill:** None of the four plans says how a *new* physics concept is introduced. They assume the AI composes existing library objects. But a ToE search may require genuinely new primitives.

**My addition:** A hypothesis layer with its own type space:

```go
type CandidateObject struct {
    Name        string
    Role        Role          // may be a new role
    Component   Component
    Dimension   Dimension
    Assumptions []Assumption
    Status      "candidate"    // not "established"
}
```

Candidate objects are not corpus objects. They are sandbox objects. They can be used in derivations but cannot be promoted to corpus without human review.

---

## Part III: Where I Disagree With the Attached Review

### Disagreement 1 — The synthesis is too ambitious

The review proposes:

> Z + DeepSeek + Gemini + new assumptions layer + open hypothesis layer

This is four plans plus two additions. That is not a synthesis. It is a superset.

The user has pushed for simplicity repeatedly. The right move is:

> **Z's conceptual spine (Identify gate, provenance, draft/commit) + DeepSeek's package architecture + one assumptions/status model. Nothing else.**

`physvet` is not needed for v0.4. It is a refinement for v0.5 or later. The library can enforce MRC through interfaces and constructors without a separate static analyzer.

### Disagreement 2 — The review understates Z's weaknesses

Z's `Identify` gate is a genuine insight. But the review does not note that Z is still primarily a *derivation* architecture, not a *knowledge* architecture. Z tells you how to derive, not how to store physics. DeepSeek's package architecture is what stores physics. The review acknowledges this but treats it as a minor point. It is not minor — it is the corpus the AI reads.

### Disagreement 3 — The review is silent on the paper translation protocol

The AI's primary task is reading papers. The review discusses hypothesis formulation and adversarial review but not the input side. This is a significant omission.

---

## Part IV: What I Would Actually Build

### The Minimum Viable Architecture

```text
┌─────────────────────────────────────────────────────────┐
│                    AI FORMULATOR                        │
│                                                         │
│  reads papers → translates to Go code                  │
│                                                         │
└──────────────────────────┬──────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                 PHYSICS LIBRARY                         │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │  core                                           │   │
│  │  - Object { Role, Component, Dimension,         │   │
│  │             Structure, Assumptions, Provenance }│   │
│  │  - Assumption                                   │   │
│  │  - Provenance                                   │   │
│  │  - Derivation                                   │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │  operations                                     │   │
│  │  - Differentiate, Integrate, Vary, Substitute,  │   │
│  │    Simplify, Compose, Identify, ...             │   │
│  │  - each checks Component, Dimension, Role       │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │  mechanics / relativity / qm / qft / statmech  │   │
│  │  - domain constructors                          │   │
│  │  - domain operations                            │   │
│  │  - framework manifests (JSON in doc comments)   │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │  hypothesis                                     │   │
│  │  - CandidateObject                              │   │
│  │  - sandbox for novel concepts                   │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
└──────────────────────────┬──────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│              DERIVATION LEDGER                          │
│                                                         │
│  every operation recorded: input, op, args, output     │
│                                                         │
└──────────────────────────┬──────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│           ADVERSARIAL REVIEWER                          │
│                                                         │
│  reads ledger → challenges → produces review            │
│                                                         │
└──────────────────────────┬──────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                  HUMAN RESEARCHER                       │
│                                                         │
│  final judgment                                         │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### The Core Type

```go
type Object struct {
    Name        string
    Role        Role
    Component   Component      // set by constructor
    Dimension   Dimension
    Structure   Expr           // internal symbolic structure
    Assumptions []Assumption   // what must hold
    Provenance  Provenance     // where it came from
}
```

### The Assumption Type

```go
type Assumption struct {
    Kind   AssumptionKind  // Domain, Regime, Convention, Approximation, ...
    Name   string
    Value  string
}
```

### The Provenance Type

```go
type Provenance struct {
    Status   Status        // Definition, Postulate, Established, Derived, Approximation, Hypothesis
    Source   string        // paper, experiment, derivation
    Framework string       // which framework (mechanics, relativity, ...)
}
```

### The Operations

Each operation is a function that:

1. Checks `Component` compatibility
2. Checks `Dimension` compatibility
3. Checks `Role` compatibility
4. Checks `Assumption` compatibility
5. Produces a new `Object`
6. Records the step in the ledger

```go
func Differentiate(f Object, wrt Object) (Object, error) {
    if err := checkComponent(f, Dance, Stage); err != nil {
        return Object{}, err
    }
    if err := checkDimension(f, wrt); err != nil {
        return Object{}, err
    }
    // ...
    result := Object{...}
    ledger.Record("Differentiate", f, wrt, result)
    return result, nil
}
```

### The Hypothesis Layer

```go
type CandidateObject struct {
    Object                       // embeds Object
    Status  Status = Hypothesis  // marked as candidate
    Review  Review               // pending human review
}

func DeclareCandidate(name string, role Role, dim Dimension) CandidateObject
```

Candidate objects are usable in derivations but cannot be promoted to the corpus without human action.

---

## Part V: The Five Gaps the Review Itself Missed

### Gap 1 — The paper translation protocol

The AI's workflow is:

```
paper → notation → library objects → derivation
```

The middle step is undefined. The library should ship with a translation guide:

```
docs/paper-translation/
    common-notation.md        // g_μν, T_μν, ψ, H, ...
    framework-mapping.md      // "this paper is in the framework of X"
    ambiguity-resolution.md   // "G might mean Newton's constant or the Einstein tensor"
```

The AI reads these. They are part of the corpus.

### Gap 2 — The adversarial review protocol

The reviewer needs a schema:

```go
type Challenge struct {
    Target    StepID
    Issue     string
    Category  ChallengeCategory  // CategoryError, HiddenAssumption, Circular, Overreach, ...
    Severity  Severity
}

type Review struct {
    Derivations []Derivation
    Challenges  []Challenge
    Verdict     string  // "sound", "unsound", "needs_revision"
}
```

The formulator responds. The cycle continues until the reviewer is satisfied or the human intervenes.

### Gap 3 — The equality and canonicalization spec

Two `Object`s are equal if:
- Same `Role`
- Same `Component`
- Same `Dimension`
- Structurally equivalent `Structure` (after canonicalization)
- Same `Assumptions` (up to ordering)

Canonicalization rules:
- Collect like terms
- Sort commutative factors
- Normalize dimensions
- Normalize indices

Without this, the ledger is unreadable.

### Gap 4 — The symbolic expression type

The `Expr` type is internal but must be specified:

```go
type Expr interface {
    Kind() ExprKind
    String() string
}

type Symbol struct { Name string }
type Literal struct { Value string; Type string }
type Sum struct { Terms []Expr }
type Product struct { Factors []Expr }
type Power struct { Base, Exp Expr }
type Derivative struct { F Expr; Wrt string }
type Integral struct { F Expr; Over string }
// ...
```

This is the engine. It must be designed before the operations can be implemented.

### Gap 5 — The multi-agent workflow

The formulator and reviewer are separate AIs. The workflow:

1. Formulator writes Go code and derives consequences.
2. Formulator outputs a `ResearchCandidate` artifact.
3. Reviewer reads the artifact and produces a `Review`.
4. Formulator reads the `Review` and revises.
5. Cycle until resolved.
6. Human receives the final artifact plus the review history.

The library must define the artifact schema.

---

## Part VI: The Correct Answer

**Plan 4Z is the best foundation.** But it is not sufficient alone.

The right architecture is:

| From | Take |
|---|---|
| **Z** | The Einstein Test, the `Identify` gate, provenance, draft/commit |
| **DeepSeek** | The package architecture (core + domain packages) |
| **Gemini** | Deferred to v0.5 — `physvet` is not needed for v0.4 |
| **Qwen** | The framing (AI = theorist), but not the epistemic claims |
| **New** | Assumptions model, knowledge-status model, hypothesis layer, paper translation protocol, adversarial review protocol |

**The minimum viable v0.4 is:**

1. `core` package: `Object`, `Assumption`, `Provenance`, `Derivation`, `Expr`
2. `operations` package: `Differentiate`, `Integrate`, `Vary`, `Substitute`, `Simplify`, `Compose`, `Identify`
3. Domain packages: `mechanics`, `relativity`, `qm`, `qft`, `statmech`
4. `hypothesis` package: `CandidateObject`, `DeclareCandidate`
5. `ledger` package: `Record`, `Inspect`
6. `review` package: `Challenge`, `Review`, `Respond`
7. Documentation: framework manifests in doc comments; paper translation guide

That is what to build. Not more.

---

## The One-Paragraph Summary

> **Plan 4Z is the best foundation, but the attached review's synthesis is too ambitious for where the project is.** Z wins on criteria 1, 2, 5, and 6 because it uses the Einstein derivation as the design test, distinguishes mechanical steps from insight, and handles provenance explicitly. DeepSeek wins on criterion 3 because it has the best package architecture. Gemini's `physvet` is valuable but premature — it should be deferred to v0.5, not incorporated into v0.4. Qwen's framing (AI = theorist, compiler = pen and paper) is correct but its claims of "absolute epistemic safety" must be rejected. The collective gaps the four plans missed are: the paper translation protocol, the adversarial review protocol, the equality and canonicalization spec, the symbolic expression type, and the multi-agent workflow. The attached review missed these five. The right architecture is: Z's conceptual spine + DeepSeek's package architecture + one assumptions/status model + the five missing pieces. The right minimum viable v0.4 is: `core` (Object, Assumption, Provenance, Derivation, Expr), `operations` (Differentiate, Integrate, Vary, Substitute, Simplify, Compose, Identify), domain packages (mechanics, relativity, qm, qft, statmech), `hypothesis` (CandidateObject, DeclareCandidate), `ledger` (Record, Inspect), `review` (Challenge, Review, Respond), and documentation (framework manifests in doc comments, paper translation guide). The library's goal is to give AI agents a formally disciplined representation of physics, let them derive new equations and formulate new hypotheses, and make every leap visible to human researchers. It does not adjudicate truth. It does not discover the theory of everything. It does not replace the physicist. It is the pen and paper that the AI uses to reason. That is the correct scope. Build that.**
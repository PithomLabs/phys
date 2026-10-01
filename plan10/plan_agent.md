## prompt

before we proceed with kernel growth gates, give prompt for the implementation agent to read the codebase in order to write AGENTS.md for AI agents, the goal is to document the process of reasoning around physics and mathematics, how to derive hypotheses given assumptions and other metadata per manifest.json with our reference implementations under mechanics and relativity package folders, now is also the time to include guidelines to write README.md for AI agents in each package folder, the goal is to create the ultimate checklist for AI agents, focus only with AGENTS.md and README.md guidelines and checklist for AI agents, we already have README.md at the root folder of the codebase


## answer

Yes. This is the right next documentation step, and I would keep it **strictly documentation-only**: the agent should read the entire codebase and manifests first, then produce the AI-facing reasoning contract without changing kernel/domain behavior.

The design should separate:

```text
AGENTS.md
    = global instructions for how an AI agent reasons with the library

package/README.md
    = concise orientation to one specific physics framework/package
```

The `manifest.json` remains the canonical machine-readable semantic source; the README is an orientation layer, not a replacement for the manifest. The manifest already carries framework identity, status, assumptions, domains, limits, anomalies, and per-item metadata such as provenance, source, derivability, reductions, known limits, and falsification conditions. physics_compiler_mvp_specs_v2 specs_v2_3

I would give the implementation agent this prompt:

---

# IMPLEMENTATION PROMPT — AI REASONING DOCUMENTATION LAYER

You are the documentation implementation agent for the Physics Compiler.

Your task is to **study the entire existing codebase and physics corpus first**, then write the AI-agent documentation layer consisting only of:

```text
AGENTS.md
mechanics/README.md
relativity/README.md
```

The root repository already has `README.md`.

This task is **documentation only**.

Do not modify production code, tests, manifests, kernel semantics, operation semantics, public APIs, or the frozen physics corpus.

Do not perform kernel growth.

Do not add General Relativity.

Do not add Einstein's 1905 derivation.

Do not add new theorem implementations.

Do not add new physics primitives.

The purpose is to document **how an AI agent should reason with the existing system and physics packages**.

---

## 1. READ BEFORE WRITING

Read the complete repository, not just selected files.

At minimum inspect:

```text
internal/kernel/
core/
ops/
session/
mechanics/
relativity/
hypothesis/
README.md
plan10/specs_v2_3.md
plan10/plan10_v2_3.md
plan10/plan7.md
plan10/plan8.md
plan10/plan9.md
```

Also inspect:

```text
mechanics/manifest.json
relativity/manifest.json
```

Read the actual Go implementations corresponding to the manifest constructors and the derivation tests.

Do not infer documentation from filenames alone.

The documentation must reflect the **actual current implementation**, while the normative semantic rules come from `specs_v2_3.md`.

---

# 2. DOCUMENTATION AUTHORITY

Use this precedence:

```text
specs_v2_3.md
    ↓
actual implementation
    ↓
manifest.json
    ↓
tests / derivation examples
    ↓
AGENTS.md
    ↓
package README.md
```

Do not invent behavior.

Do not document capabilities that are not actually implemented.

Where the implementation is deliberately bounded, explain the boundary.

The manifest is the canonical machine-readable semantic source for populated framework packages. Go documentation may summarize it, but must not replace it. physics_compiler_mvp_specs_v2

---

# 3. AGENTS.md PURPOSE

`AGENTS.md` is the **global operating manual for AI agents using this repository**.

Its primary purpose is:

> Teach an AI agent how to reason about physics and mathematics using the typed primitives, assumptions, metadata, operations, manifests, sessions, and hypothesis machinery without confusing formal derivation with physical truth.

The document must make the agent's reasoning workflow explicit.

It should be a practical checklist, not an essay.

---

# 4. AGENTS.md — REQUIRED CONTENT

## A. What this system is

Explain:

```text
AI agent
    = theorist / formulator / reasoner

Physics Compiler
    = formal symbolic substrate + corpus

Human researchers / empirical reality
    = scientific authority
```

The library checks formal structure, dimensions, categories, assumptions, conventions, provenance, containment, and artifact integrity. It does not adjudicate physical truth. specs_v2_3

---

## B. The three epistemic layers

Teach the agent to distinguish:

```text
ESTABLISHED FRAMEWORK
    e.g. relativity/

DERIVED ARTIFACT
    result obtained from explicit operations

HYPOTHESIS
    provisional candidate material
```

Make clear:

```text
ESTABLISHED
≠ universally true
≠ final ontology
```

and:

```text
formally derivable
≠ empirically established
```

Hypothesis-derived results remain HYPOTHESIS and cannot be promoted automatically. specs_v2_3

---

# 5. AGENTS.md — REASONING WORKFLOW

Document a canonical AI reasoning workflow.

Use a checklist like:

```text
[ ] 1. Identify the physical question.
[ ] 2. Identify the framework/package involved.
[ ] 3. Read that package's README.
[ ] 4. Read the relevant manifest entries.
[ ] 5. Determine framework assumptions and conventions.
[ ] 6. Determine regime/domain/scope.
[ ] 7. Identify the available typed primitives.
[ ] 8. Identify the relevant relations.
[ ] 9. Identify known limits/limitations/anomalies.
[ ] 10. Determine whether the desired statement already exists as corpus data.
[ ] 11. If deriving it, do NOT use the stored target relation as the derivation source.
[ ] 12. Construct the derivation from explicit premises.
[ ] 13. Use only operations supported by the MVP.
[ ] 14. Preserve assumptions throughout.
[ ] 15. Check dimensions and physical kind at each step.
[ ] 16. Check provenance at each step.
[ ] 17. Check for HYPOTHESIS contamination.
[ ] 18. Record every meaningful derivation step through Session when required.
[ ] 19. Make branch choices and preconditions explicit.
[ ] 20. Distinguish formal result from physical interpretation.
[ ] 21. Record unresolved assumptions/limitations.
[ ] 22. Only then formulate a candidate hypothesis if needed.
```

This must be one of the central sections of `AGENTS.md`.

---

# 6. AGENTS.md — REASONING FROM ASSUMPTIONS

Explain that equations are not to be treated as context-free strings.

An AI agent must ask:

```text
What assumptions make this transformation valid?
What regime am I in?
What convention is being used?
What dimensions and kinds are involved?
What branch conditions apply?
```

The existing assumption system explicitly models deterministic union, exact conflict detection, structured expression-valued assumptions, and operation-generated preconditions. specs_v2_3

Teach the agent:

```text
Do not silently introduce an assumption.

Do not silently discard an assumption.

Do not treat a theorem-like transformation as valid merely because
the expressions look algebraically similar.

Do not replace structured assumptions with prose.
```

---

# 7. AGENTS.md — MATHEMATICAL REASONING DISCIPLINE

Document the distinction between:

```text
construction
transformation
identification
hypothesis formation
```

The agent should understand:

```text
Simplify()
    = mechanical transformation

Compare()
    = constructs a relation artifact

Solve()
    = bounded supported pattern

SelectBranch()
    = explicit branch selection under constraint

Identify()
    = explicit conceptual identification owned by Session
```

Do not describe the MVP as a general CAS or theorem prover. Its solver and calculus capabilities are deliberately bounded. For example, `Solve` accepts one exact quadratic pattern and does not provide general equation solving. specs_v2_3

---

# 8. AGENTS.md — DERIVATION CHECKLIST

Require the AI to check every proposed derivation for:

```text
[ ] Premises explicitly identified
[ ] Framework identified
[ ] Assumptions explicitly identified
[ ] Conventions identified
[ ] Physical regime identified
[ ] Input objects valid
[ ] Dimensions compatible
[ ] Physical kinds compatible
[ ] Operation supported by MVP
[ ] Operation parameters canonical
[ ] No hidden assumption introduced
[ ] No assumption silently discarded
[ ] Provenance correct
[ ] Hypothesis contamination preserved
[ ] Branch selection explicit
[ ] Preconditions explicit
[ ] No corpus result used as an unexplained derivation shortcut
[ ] Final expression canonical
[ ] Session/ledger used when a durable derivation is required
[ ] Scientific interpretation clearly separated from formal result
```

---

# 9. AGENTS.md — HOW TO DERIVE HYPOTHESES

This is especially important.

Document a disciplined hypothesis workflow:

```text
OBSERVATION / QUESTION
        ↓
FRAMEWORK
        ↓
KNOWN ESTABLISHED PREMISES
        ↓
ASSUMPTIONS
        ↓
FORMAL TRANSFORMATION
        ↓
UNRESOLVED GAP / ANOMALY
        ↓
CANDIDATE HYPOTHESIS
        ↓
DERIVATION
        ↓
PREDICTIONS
        ↓
FALSIFICATION CONDITIONS
        ↓
RECOVERY / LIMIT CLAIMS
        ↓
RESEARCH CANDIDATE
        ↓
HUMAN EVALUATION
```

The AI must understand that it may formulate a candidate, but cannot promote it into trusted corpus authority.

The candidate model already supports:

```text
Prediction
FalsificationCondition
RecoveryClaim
AnomalyReference
```

and the library records those structures rather than deciding whether they occurred in nature. specs_v2_3

---

# 10. AGENTS.md — MANIFEST-FIRST WORKFLOW

Teach the AI:

> Before reasoning about a package, inspect its manifest and README.

For each manifest item, explain how the AI should use:

```text
id
constructor
kind
name
statement
canonical_expr
dimension
provenance_status
source
assumptions
derivable_from
reduces_to
known_limits
anomalies
falsification_conditions
```

These are already the package's semantic vocabulary. specs_v2_3

Explicitly teach:

```text
statement
    = orientation/documentation

canonical_expr
    = machine-checkable mathematical object
```

The agent must never extract mathematical meaning from the prose `statement` field when `canonical_expr` exists. specs_v2_3

---

# 11. AGENTS.md — THEORY PACKAGE BOUNDARIES

Teach the AI the package architecture:

```text
internal/kernel
    = generic formal machinery

mechanics/
    = Classical Mechanics

relativity/
    = Special Relativity only

hypothesis/
    = provisional candidate space
```

A framework's assumptions must come explicitly from that framework.

Do not merge framework assumptions implicitly.

Do not treat one framework as the universal ontology of nature.

The established architecture explicitly keeps theory packages outside the generic kernel. plan7(5)

---

# 12. AGENTS.md — CURRENT RELATIVITY BOUNDARY

State clearly:

```text
relativity/
    = SPECIAL RELATIVITY ONLY
```

Therefore the agent must not assume that this package contains:

```text
General Relativity
gravitational dynamics
curved-spacetime dynamics
quantum gravity
```

The current relativity framework explicitly carries assumptions including Minkowski spacetime, Lorentz symmetry, special-relativistic regime, and no gravitational dynamics in the package. physics_compiler_mvp_specs_v2

---

# 13. AGENTS.md — IMPORTANT E=mc² RULE

Document the MVP's primary demonstration.

The AI must distinguish:

```text
stored corpus artifact:
    MassEnergyRelation

derivation:
    EnergyMomentumRelation
    → ZeroThreeMomentum
    → Substitute
    → Simplify
    → Solve
    → Compare
    → SelectBranch
    → m*c²
```

The stored `MassEnergyRelation` must not be used as the derivation source when demonstrating reasoning capability.

State explicitly:

```text
This is a formal derivation from encoded special-relativistic premises.

It is not:
- Einstein's full 1905 historical reconstruction
- empirical discovery
- physical-truth adjudication
```

---

# 14. AGENTS.md — WHEN TO FORM A HYPOTHESIS

The agent should form a hypothesis only when the requested reasoning goes beyond the established corpus/framework.

Checklist:

```text
[ ] Existing corpus does not already answer the question.
[ ] The desired statement requires a new assumption, structure, or bridge.
[ ] The candidate is explicitly marked HYPOTHESIS.
[ ] The candidate's assumptions are explicit.
[ ] The candidate's predictions are explicit where applicable.
[ ] At least one falsification condition is defined where applicable.
[ ] Recovery/limit claims are explicit where applicable.
[ ] Anomalies motivating the candidate are recorded where relevant.
[ ] No trusted corpus status is assigned automatically.
```

---

# 15. AGENTS.md — DO NOT DO

Include a compact anti-pattern checklist:

```text
[ ] Do not invent physics not present in the corpus.
[ ] Do not invent assumptions.
[ ] Do not silently change framework boundaries.
[ ] Do not treat prose statements as machine equations.
[ ] Do not bypass dimensional/category checks.
[ ] Do not bypass provenance.
[ ] Do not promote hypotheses.
[ ] Do not call a stored result a derivation.
[ ] Do not use unsupported symbolic operations.
[ ] Do not assume a theory is universally true.
[ ] Do not infer empirical validation from formal consistency.
[ ] Do not modify the kernel merely to make one theory easier to express.
```

---

# 16. PACKAGE README GUIDELINES

Create/refresh:

```text
mechanics/README.md
relativity/README.md
```

Each package README should be intentionally short.

The README is an **AI orientation document**, not a duplicate of `manifest.json`.

Target:

```text
~1 page
```

unless the actual package complexity makes slightly more necessary.

---

# 17. PACKAGE README REQUIRED TEMPLATE

Every established framework package README must use the following structure:

```text
# <Framework Name>

## Framework Identity
Framework ID:
Corpus Status:

## Scope
What this package represents.

## Not in Scope
What it explicitly does not represent.

## Assumptions
Core framework assumptions.

## Conventions
Important conventions.

## Core Primitives
The important typed objects.

## Core Relations
The important equations/relations.

## Derivation Targets
What the package's executable derivations demonstrate.

## Limits and Known Boundaries
Known limits, exclusions, anomalies.

## AI Reasoning Guidance
How an AI agent should use this package.

## Canonical Metadata
manifest.json
```

Do not invent fields outside the actual package semantics unless clearly labeled as orientation.

---

# 18. RELATIVITY README

For `relativity/README.md`, explicitly say:

```text
Framework ID: special_relativity
Corpus Status: ESTABLISHED

Scope:
Special relativity only.

Not in scope:
General relativity
gravitational dynamics
quantum mechanics
quantum gravity

Core assumptions:
Minkowski spacetime
Lorentz symmetry
special-relativistic regime
no gravitational dynamics

Convention:
metric signature -+++

Important primitives:
Spacetime
MinkowskiMetric
RestMass
Energy
ThreeMomentum
FourMomentum
SpeedOfLight
Velocity

Important relations:
LorentzFactor
EnergyMomentumRelation
MassEnergyRelation

Primary derivation target:
E = mc²

Important limitation:
No gravitational dynamics.
```

These values must be verified against the actual manifest/code rather than guessed. The specification already defines this package as special relativity and requires these framework boundaries. physics_compiler_mvp_specs_v2

Also explain:

```text
MassEnergyRelation
    = reusable corpus artifact

derivation_test.go
    = executable demonstration of deriving the result
```

Do not present the README as claiming Einstein's historical 1905 derivation.

---

# 19. MECHANICS README

For `mechanics/README.md`, identify:

```text
Framework ID
Corpus status
Classical/nonrelativistic scope
Core primitives
Core relations
Known limits
Anomalies
Typical derivation targets
AI usage guidance
manifest.json
```

At minimum cover:

```text
Mass
Time
Position
Velocity
Acceleration
Force
Momentum
Energy
KineticEnergy

F = ma
p = mv
K = 1/2 mv²
```

Use the actual current manifest and implementation for all metadata.

---

# 20. README RULES FOR ALL FUTURE PACKAGES

The documentation must establish a reusable template for future packages.

An AI agent encountering a new framework should be able to answer within minutes:

```text
What is this framework?
What does it assume?
What regime does it describe?
What does it exclude?
What primitives does it provide?
What relations does it provide?
What can I derive?
What are its known limits?
What anomalies/limitations are recorded?
Where is the canonical machine-readable metadata?
```

Do not require the AI to understand the Go implementation merely to discover these answers.

---

# 21. IMPORTANT SEPARATION: THEOREMS VS PRIMITIVES

Document this rule explicitly in `AGENTS.md`:

```text
ESTABLISHED PHYSICAL PRIMITIVE
    may belong to an established theory package

GENERIC FORMAL PRIMITIVE
    may belong to internal/kernel if it is truly theory-neutral

THEOREM / DERIVED RELATION / SPECIALIZED RESULT
    normally belongs to the relevant user/theory package

HYPOTHESIS
    belongs to hypothesis/research space
```

Examples for future reasoning:

```text
General Relativity primitive
    → general_relativity package

Equivariance theorem
    → user/theory library

Valentini H-Theorem
    → user/theory library

Novel TOE bridge
    → hypothesis/research package
```

Do not put theorem content into `internal/kernel` merely because the theorem is mathematically established.

The kernel remains the smallest generic formal substrate.

---

# 22. KERNEL-GROWTH BOUNDARY

This task is documentation-only, but `AGENTS.md` should prepare agents for the future Growth Gate.

Document:

> Before requesting a kernel addition, first determine whether the capability belongs in an existing framework package, user library, adapter, or higher-level reasoning layer.

Only consider kernel growth when the missing capability is genuinely:

```text
generic
foundational
theory-neutral
reusable across multiple frameworks
not safely representable outside kernel
```

Do not implement this gate now.

Only document the reasoning discipline.

---

# 23. OUTPUT REQUIREMENTS

Produce exactly:

```text
AGENTS.md
mechanics/README.md
relativity/README.md
```

No other files.

Do not modify root `README.md`.

Do not modify:

```text
internal/kernel/**
core/**
ops/**
session/**
mechanics/*.go
relativity/*.go
hypothesis/**
manifest.json
specifications
plans
```

unless correcting a broken documentation reference is absolutely necessary; prefer not to touch them.

---

# 24. QUALITY CHECKLIST

Before completing the documentation, verify:

```text
[ ] AGENTS.md describes actual current behavior.
[ ] AGENTS.md teaches an explicit physics-reasoning workflow.
[ ] AGENTS.md teaches assumption-first reasoning.
[ ] AGENTS.md teaches provenance reasoning.
[ ] AGENTS.md teaches hypothesis containment.
[ ] AGENTS.md teaches manifest-first package discovery.
[ ] AGENTS.md distinguishes formal validity from physical truth.
[ ] AGENTS.md distinguishes established framework from final ontology.
[ ] AGENTS.md distinguishes theorem from primitive.
[ ] AGENTS.md explains E=mc² demonstration correctly.
[ ] AGENTS.md explains current relativity boundary.
[ ] mechanics README is package-specific.
[ ] relativity README is package-specific.
[ ] both READMEs identify scope and exclusions.
[ ] both READMEs identify assumptions/conventions.
[ ] both READMEs identify important primitives/relations.
[ ] both READMEs identify limits/anomalies.
[ ] both READMEs point AI agents to manifest.json.
[ ] README content does not duplicate the entire manifest.
[ ] no unsupported capability is documented.
[ ] no production code changed.
[ ] no physics corpus changed.
[ ] no kernel semantics changed.
```

---

# 25. FINAL REPORT

Return:

```text
## Documentation Implementation Report

### Files created/updated
- AGENTS.md
- mechanics/README.md
- relativity/README.md

### Code changes
NONE

### Corpus changes
NONE

### Kernel changes
NONE

### Verification
- documentation cross-checked against implementation
- manifests cross-checked
- E=mc² description cross-checked
- relativity scope cross-checked
- mechanics scope cross-checked

### Final status
DOCUMENTATION-READY
```

The goal is not merely to explain the repository to humans.

The goal is:

> **An AI agent should be able to enter this repository, read `AGENTS.md`, read the relevant package `README.md`, inspect the manifest, and know how to reason about the physics without first reverse-engineering the Go implementation.**

Do not turn `AGENTS.md` into generic software-engineering instructions. Its central subject is **disciplined machine reasoning over physics and mathematics using this formal substrate**.

---

This also preserves the architectural principle we've established: the package is the AI's bounded physics textbook, while the generic kernel remains the formal substrate. The manifest is the canonical semantic layer, and the README is the fast orientation layer. physics_compiler_mvp_specs_v2

And importantly, this documentation pass does **not** accidentally turn General Relativity, equivariance results, Valentini's H-Theorem, or any future theorem into kernel content. Those remain package/library material unless a later, explicit Growth Gate establishes that something genuinely belongs in the generic substrate.
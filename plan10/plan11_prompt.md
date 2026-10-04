You are the senior implementation planner for the Physics Compiler repository.

Your task is NOT to implement General Relativity yet.

Your task is to write a rigorous implementation plan for the first General Relativity (GR) workload, which will serve as the kernel stress test under the consolidated three-level Growth Gate defined below.

The plan must be grounded in the repository as it exists today. Inspect the actual repository, frozen specifications, current implementation, tests, AGENTS.md, mechanics/README.md, relativity/README.md, and the four attached review documents before producing the plan.

Authoritative baseline:
- `specs_v2_3.md` remains frozen.
- The existing MVP is frozen at the current Level-3 architecture.
- `internal/kernel` is the trusted substrate.
- `core`, `ops`, and `session` form the trusted Level-3 machinery.
- `mechanics` and `relativity` are reference implementations/corpora, not kernel-owned physics ontology.
- General Relativity must be implemented as an external/theory-local package, e.g. `general_relativity/`, outside `internal/kernel`.
- Do NOT redesign the kernel in anticipation of GR.
- Do NOT implement GR in this task.
- Do NOT create a new shared mathematical library merely because GR may need one.
- The purpose of this phase is to produce the implementation plan and the stress-test methodology.

Architectural principle to preserve:

    The kernel is like Go itself:
    small in what it knows,
    rich in generic primitives,
    robust in invariant enforcement,
    composable and extensible,
    but agnostic about particular mathematics and physical theories.

The long-term architecture has three levels:

LEVEL 1 — Theory-local userland
    general_relativity/
    quantum/
    qft/
    etc.

LEVEL 2 — Shared mathematical/semantic libraries
    only created after demonstrated reuse by at least two established theories
    and human curation.

LEVEL 3 — Trusted kernel
    only generic trusted primitives/invariants that higher levels cannot
    safely express or enforce.

General Relativity MUST NOT jump directly from Level 1 to Level 3.

==================================================
1. CORE PURPOSE OF THE GR WORKLOAD
==================================================

Treat GR as a hostile workload against the existing substrate.

The question is NOT:

    "What does GR need that we should add to the kernel?"

The question is:

    "How far can a faithful GR implementation be built above the
     existing substrate, and where do actual capability boundaries occur?"

The GR implementation should therefore attempt to express real GR structures and derivations using the existing architecture before proposing any new abstraction.

The plan must explicitly preserve this principle:

    GR is a workload for evaluating kernel sufficiency,
    not a justification for kernel expansion.

The default outcome is NO KERNEL GROWTH.

==================================================
2. FAILURE CLASSIFICATION
==================================================

Every material obstacle discovered during GR implementation must be classified into exactly one of these categories:

1. SPEC-INTENDED-BOUND
   The MVP intentionally does not support the capability.

2. PACKAGE-SOLVABLE
   `general_relativity/` can express or solve the problem entirely in
   theory-local userland.

3. REPRESENTABLE-BUT-UNFAITHFUL
   The current substrate can encode something, but doing so loses
   semantically important information.

4. UNREPRESENTABLE
   The current substrate cannot faithfully encode the required structure.

5. SILENTLY-WRONG
   The current substrate accepts an operation/representation but can
   produce a semantically incorrect result without exposing the problem.

The plan must require the GR workload log to record, for every significant
failure:

- requirement/capability being attempted;
- minimal example;
- current representation;
- observed result;
- expected semantic result;
- failure category;
- why the category was chosen;
- whether a Level-1 workaround exists;
- whether a Level-2 workaround exists;
- whether a Level-3 kernel candidate is actually implicated.

Do not allow the failure log to become a wishlist of desirable mathematical
features.

==================================================
3. THREE-LEVEL PROMOTION RULE
==================================================

The plan must use this exact routing model:

LEVEL 1 → LEVEL 2
    Requires demonstrated reuse by at least two established theories,
    with a named second consumer and a worked example.

LEVEL 2 → LEVEL 3
    Requires a demonstrated failure of the Level-2 abstraction to preserve
    or enforce a genuinely generic invariant safely without kernel support.

"Tensors are useful to GR and QFT" is not sufficient evidence.

"Two theories use this mathematical object" is only a candidate signal for
Level 2, not automatic promotion.

"GR cannot be implemented conveniently without this" is not evidence for
Level 3.

The burden of proof for kernel growth is:

    generic
    theory-neutral
    mathematically/theoretically non-semantic
    foundational
    reusable
    trusted/invariant-bearing
    cannot reasonably live in Level 1 or Level 2
    cannot be safely expressed above the kernel.

==================================================
4. PRE-GR CONFORMANCE PINS
==================================================

The plan must include these as mandatory preconditions before the main GR
workload begins.

A. Session.Identify result kind

The current frozen material describes Identify as constructing:

    Relation(eq, a.Expr(), b.Expr())

but does not sufficiently pin the result Kind.

The implementation plan must explicitly pin the intended behavior:

    Session.Identify always returns KindRelation.

It must NEVER:
- retype an Expression into Energy, Mass, RestMass, etc.;
- act as a semantic-ascription mechanism;
- create a nominal physical type merely because the caller supplied a
  justification string.

Include a regression test such as:

    Identify(named physical object, expression object, justification)
        -> KindRelation
        -> relation expression
        -> IDENTIFIED provenance

The plan must explicitly verify that no Identify path becomes a covert
semantic-ascription mechanism.

Do not alter `Substitute` to weaken its existing dimension + identical-kind
requirement.

B. Negative-integer power differentiation

The MVP differentiation rules support non-negative integer constant powers.
The current wording creates ambiguity around negative integer powers.

Before GR Pass 1, pin this behavior:

    Differentiate(Pow(x, negative-integer))
        -> UnsupportedOperationError

Do NOT expand MVP differentiation merely to make GR easier.

This ensures that a negative-power encounter in GR is a deliberate
SPEC-INTENDED-BOUND rather than an ambiguous implementation defect.

C. Physical-kind documentation wording

The plan must identify any documentation wording that implies the current
kernel catches every energy-vs-torque category mismatch.

The accurate statement is:

    Named-kind compatibility protects additive operations when named kinds
    are retained. Expression + Expression with equal dimensions is permitted
    by design. Higher-level preservation/ascription of semantic category is
    outside the MVP.

Do not change MRC-003 merely to hide this limitation.

D. Frozen-spec discipline

Do not modify `specs_v2_3.md` in this planning task.

Where a clarification is needed because the frozen spec is underspecified,
record it explicitly as a pre-GR conformance pin / post-freeze clarification
for later human approval and implementation.

Do not silently reinterpret the frozen specification.

==================================================
5. CURRENT SIX-CONCERN LEDGER
==================================================

The plan must carry forward the following already-adjudicated concerns and
must NOT reopen them as automatic kernel-growth proposals:

C1 — Type decay
    Derived physical objects become KindExpression.
    Equal-dimension Expression + Expression is permitted.
    Example:
        m*c² -> Expression, M L² T^-2
        r*F -> Expression, M L² T^-2
        Add(...) -> accepted
    Classification:
        real limitation;
        investigate as a Level-2 semantic/ascription capability;
        not automatically kernel growth.

C2 — Substitution catch-22
    Derived expression E/c² cannot substitute into RestMass because
    Substitute requires identical Kind.
    Classification:
        real limitation;
        do NOT weaken Substitute;
        investigate semantic ascription above kernel.

C3 — Assumption semantic conflict
    Assumptions use deterministic exact-key merging.
    Distinct keys may encode logically inconsistent conditions.
    Classification:
        real limitation;
        constraint/entailment logic belongs above kernel.

C4 — Symbol identity collision
    Symbol identity is currently raw string identity.
    Classification:
        real limitation;
        investigate namespaced/entity-aware symbol layer above kernel.

C5 — Transcendental functions / symbolic exponents
    Arbitrary transcendental calls are not part of the MVP.
    Pow uses exact rational exponents, not symbolic exponent objects.
    Classification:
        not a current MVP defect;
        future extension commitment only.

C6 — Closed Kind ontology
    Kind is deliberately closed and ordinal-stable in mrc-v0.4.
    Classification:
        keep closed;
        do NOT create a physics subtype hierarchy in the kernel.

The GR plan must treat these as the existing concern ledger.

Do not repeatedly relitigate them.

==================================================
6. WHAT GR MUST TEST
==================================================

Design a staged GR workload that begins with minimal structures and
progressively exercises more demanding mathematics.

The exact GR workload should be derived from the actual implementation
requirements, not invented merely to justify a preferred architecture.

At minimum, investigate the ability to represent and compose:

- spacetime/manifold concepts;
- coordinates/charts;
- metric structure;
- inverse metric;
- tensor components and index structure;
- contraction;
- raising/lowering indices;
- connection/Christoffel symbols;
- covariant differentiation;
- curvature;
- Ricci tensor/scalar;
- Einstein tensor;
- stress-energy tensor;
- Einstein field equation;
- simple exact GR derivations or identities sufficient to exercise the
  substrate.

Do not assume that all of these belong in a common library.

Initially place them in `general_relativity/` as Level-1 workload structures.

Where some mathematics naturally belongs in an internal GR-local helper
package, say so explicitly.

Do not prematurely create:
    tensor/
    manifold/
    semantic/
    symbols/
    constraints/
    calculus/

as shared packages.

Instead, use the GR workload to discover whether any of these become genuine
Level-2 candidates.

==================================================
7. REPRESENTATION STRATEGY
==================================================

The plan must explain how GR can first be attempted using the current
`core.Object` + existing expression machinery.

Do not assume that an existing `Kind` must be created for every GR object.

Prefer higher-level Go structures/wrappers where possible.

For example, the plan should examine whether structures resembling:

    Tensor
    Index
    Component
    Metric
    CoordinateChart
    Connection

can initially live as domain-local data structures that contain or refer to
`core.Object` values, while preserving the kernel's invariants.

The plan must identify exactly where scalar `core.Object` is sufficient and
where richer GR-local structures become necessary.

Do not introduce a tensor type into the kernel merely because tensor algebra
is difficult.

==================================================
8. MULTI-THEORY REUSE TEST
==================================================

Whenever the GR implementation appears to produce a broadly reusable
mathematical abstraction, mark it as a potential Level-2 candidate.

However, the plan must not promote it immediately.

Require:

    named second established theory
    +
    worked second-theory example
    +
    demonstrated shared semantics

before Level-1 → Level-2 promotion is considered.

The second theory should be concrete, such as Quantum Mechanics or QFT,
not "future theories might use this."

If no second consumer is demonstrated during this phase, keep the abstraction
GR-local.

==================================================
9. LEVEL-2 CANDIDATE TEST
==================================================

For every possible Level-2 candidate, document:

- what the abstraction is;
- why GR needs it;
- exactly where it first appears;
- whether GR-local implementation works;
- whether a second theory actually uses it;
- worked second-theory example if available;
- what generic invariant it preserves;
- whether it can remain outside the kernel;
- whether it depends on theory semantics;
- whether it requires assumptions or provenance semantics;
- whether it can be implemented without modifying the frozen kernel.

A Level-2 candidate should be promoted only by human curation after evidence.

==================================================
10. LEVEL-3 KERNEL CANDIDATE TEST
==================================================

If and ONLY IF the GR workload exposes something that appears to require
kernel support, create a Growth Evidence Record.

That record must include:

1. Capability missing
2. Minimal concrete counterexample
3. Failure category
4. Level-1 implementation attempt
5. Level-2 implementation attempt
6. Why Level 1 failed
7. Why Level 2 failed
8. Genericity argument
9. Foundationality argument
10. Named second consumer
11. Worked second-consumer example
12. Non-goal/collision analysis
13. Silent-wrongness analysis
14. Artifact/canonicalization compatibility impact
15. MRC impact
16. Replay/hash impact
17. Migration/backward-compatibility impact
18. Minimum possible kernel change
19. Independent reviewer requirement
20. Human approval requirement

No kernel change should be proposed merely because a mathematically richer
library would be convenient.

==================================================
11. SILENT-WRONGNESS IS THE MOST IMPORTANT GR TEST
==================================================

The GR implementation must actively search for cases where the current
substrate ACCEPTS a derivation but represents it incorrectly.

Pay particular attention to:

- symbol identity ambiguity;
- field/component symbol differentiation;
- index contraction mistakes;
- raising/lowering mistakes;
- dimensionally valid but semantically invalid operations;
- loss of semantic type after algebraic derivation;
- assumptions that coexist but are logically inconsistent;
- substitutions that preserve dimensions but lose semantic identity.

Do not only test "does the system reject bad input?"

Also test:

    "Can it accept a bad derivation and silently make it look valid?"

Any such case must receive especially strong review.

==================================================
12. KERNEL INVARIANTS THAT MUST REMAIN UNTOUCHED
==================================================

The plan must explicitly state that GR must preserve:

- immutable trusted objects;
- internal kernel mint authority;
- no public generic trusted-object factory;
- exact rational dimensions;
- deterministic canonicalization;
- stable expression node ordinals;
- stable Kind ordinals for mrc-v0.4;
- MRC enforcement;
- assumption/convention integrity;
- provenance propagation;
- hypothesis contamination;
- session authority;
- derivation ledger integrity;
- deterministic replay;
- artifact hash compatibility.

Any proposed GR feature that bypasses these rules is architecturally invalid.

==================================================
13. REPOSITORY / FREEZE DISCIPLINE
==================================================

The existing MVP tree is frozen.

Do not silently add files to the frozen implementation tree during this
planning task.

The plan must distinguish:

    existing frozen files
    proposed GR files
    possible future Level-2 files
    possible kernel changes (if ever justified)

The plan must state where the eventual GR package would live without
prematurely changing the frozen MVP.

Do not edit the existing mechanics or Special Relativity corpus merely to
make GR fit.

Special Relativity remains Special Relativity.

General Relativity is separate.

==================================================
14. IMPLEMENTATION PLAN FORMAT
==================================================

Produce one implementation-plan document, for example:

    GR_IMPLEMENTATION_PLAN.md

Do NOT write implementation code.

Do NOT modify the kernel.

Do NOT modify `specs_v2_3.md`.

Do NOT create speculative shared libraries.

The plan must contain:

1. Executive Summary
2. Frozen Baseline and Non-Negotiable Invariants
3. Three-Level Architecture
4. GR Workload Objectives
5. Pre-GR Conformance Pins
6. GR Package/Repository Structure
7. Representation Strategy
8. Staged GR Implementation Phases
9. Canonical GR Workloads / Worked Examples
10. Failure Taxonomy and Evidence Logging
11. Existing Six-Concern Ledger
12. Level-1 → Level-2 Promotion Criteria
13. Level-2 → Level-3 Growth Gate
14. Silent-Wrongness Test Strategy
15. Mutation/Negative Testing Strategy
16. Canonicalization / Replay / Artifact Impact Analysis
17. Human Review / Independent Veto Points
18. Definition of Done for the GR Stress Test
19. Explicit Non-Goals
20. Expected Outcomes:
       NO-GROWTH
       LEVEL-2-GROWTH
       KERNEL-GROWTH CANDIDATE
       ESCALATE-TO-SPEC

For every implementation phase include:

- objective;
- files/packages expected to change;
- dependencies;
- exact capabilities exercised;
- representative worked examples;
- tests;
- expected failure modes;
- evidence to capture;
- promotion level implicated;
- exit criteria.

==================================================
15. DEFINITION OF DONE
==================================================

The implementation plan is complete only if another senior coding agent could
execute it without having to invent architectural decisions.

The plan must make clear:

- what gets built first;
- what deliberately does not get built;
- how GR is represented;
- how failures are recorded;
- how each failure is classified;
- how Level-1 versus Level-2 versus Level-3 decisions are made;
- what constitutes actual evidence for kernel insufficiency;
- what constitutes merely a missing userland mathematical abstraction;
- what would stop the experiment;
- who makes the final promotion decision.

Do not conclude in advance that any kernel primitive will be added.

The expected default is:

    GR stresses the substrate,
    evidence is collected,
    abstractions remain in userland unless demonstrated reusable,
    and the kernel remains unchanged unless the Growth Gate proves otherwise.

==================================================
16. REVIEW MATERIAL
==================================================

Use the current repository plus the attached review documents as the
architectural input.

In particular, reconcile the four review documents and preserve only the
net-valid conclusions.

Important consensus:
- the MVP is robust for its intended scope;
- the six concerns do not establish kernel growth;
- the three-level promotion ladder is the correct evolution model;
- GR is the appropriate next workload;
- mathematical and physical objects should initially remain above the kernel;
- shared abstractions emerge only through demonstrated reuse and human
  curation;
- kernel growth is exceptional and evidence-driven.

Also inspect the actual repository before finalizing the plan. Do not rely
solely on historical plan documents when the repository contains newer
implementation state.

Finally, clearly separate:
    facts verified from the repository,
    requirements from the frozen specification,
    decisions already locked by the architecture,
    proposed implementation steps,
    and open questions that genuinely require human judgment.
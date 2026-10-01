I reviewed **Plan 7 — Final Kernel Freeze + Theory/Epistemic Package Isolation** carefully against the consolidated verdict and the decisions we locked afterward. It is now substantially coherent: it preserves Plan 6 unchanged, adds only Gate H, locks the reverse constructor allowlist and `pow` encoding, and explicitly makes assumption non-leakage a two-level proof. :chatgpt-content-reference{index="0"}

The important epistemic distinction is also correctly captured: **special relativity remains `ESTABLISHED`; package isolation does not demote it to `HYPOTHESIS`**. Its assumptions, conventions, scope and empirical status remain explicit rather than becoming ambient axioms of the kernel. :chatgpt-content-reference{index="1"}

I would make only one small clarification before handing it to the implementation agent: the import-boundary wording `kernel←core←ops←session` should be interpreted as the dependency direction **kernel ← core ← ops ← session**, with domain packages depending only on the permitted substrate, exactly as the rest of Plan 7 states. No architectural change is needed.

The assumption strategy is also correctly locked to **both levels**: an operation-level subset invariant, plus cross-package forbidden-key tests rather than brittle exact assumption sets. :chatgpt-content-reference{index="2"}

## Final implementation prompt

```text
IMPLEMENTATION TASK
===================

Implement Plan 7 — Final Kernel Freeze + Theory/Epistemic Package Isolation.

You are performing the final hardening pass before freezing the Physics Compiler MVP kernel.

This is NOT a redesign.
This is NOT a feature phase.
This is NOT a physics-corpus expansion.

Treat the attached Plan 7 as authoritative for this pass.

==================================================
1. FROZEN SCOPE
==================================================

Preserve unchanged unless a concrete gate below strictly requires a production fix:

- specs_v2_3.md
- Plan 10 architecture
- existing 39-file implementation tree
- package boundaries
- dependency direction
- internal/kernel authority
- core trusted facade/type aliases
- ops purity
- session replay/ledger architecture
- existing mechanics corpus
- existing special-relativity corpus
- existing E=mc² derivation path
- stdlib-only requirement

Do NOT add:

- new physics domains
- General Relativity
- Quantum Mechanics
- Electromagnetism
- Einstein-1905 historical derivation corpus
- parser
- CAS
- numerical engine
- theorem prover
- general solver
- automatic hypothesis promotion
- truth scoring/adjudication
- EBP 2.1 runtime integration
- AI planner/orchestrator

Einstein's original 1905 mass-energy derivation is explicitly deferred until AFTER kernel freeze.

Do not solve future research problems now.

==================================================
2. LOCKED DECISIONS
==================================================

Reverse non-manifest constructor allowlist is CLOSED and EXPLICIT:

- NewKineticEnergy
- Velocity
- ZeroThreeMomentum
- ZeroEnergy
- ZeroVelocity
- hypothesis concept

Do NOT infer permission from provenance, corpus status, naming, or any other property.

Every other ESTABLISHED-yielding trusted constructor must map to an approved manifest item.

OperationParams `pow` is STRICT:

Canonical operational form:

{"kind":"pow","exponent":"2/1","operator":"","justification":""}

A non-empty unused `operator` is rejected.

Therefore:

{"kind":"pow","exponent":"2/1","operator":"eq","justification":""}

MUST be rejected.

Do not modify specs_v2_3.md.

==================================================
3. PRIMARY MVP ACCEPTANCE TARGET
==================================================

Preserve and protect the E=mc² architecture test.

Required actual operation sequence:

EnergyMomentumRelation
    ↓
ZeroThreeMomentum
    ↓
Substitute
    ↓
Simplify
    ↓
Solve(Energy)
    ↓
Compare(Energy, ZeroEnergy, gte)
    ↓
SelectBranch
    ↓
m*c²

The derivation must use actual operations.

The test MUST NOT call MassEnergyRelation() to obtain the final result.

Harden the prohibition with a mechanical source-level assertion in:

relativity/derivation_test.go

Do not replace this with a weaker final-value assertion.

Also preserve:

Differentiate(KineticEnergy, v) → m*v

and:

Limit(LorentzFactor, v, 0) → 1

Neither may depend on hardcoded final-result shortcuts.

==================================================
4. BASELINE FIRST
==================================================

Before modifications:

1. Inspect the complete repository.
2. Verify the existing 39-file tree.
3. Run:

go test ./...
go vet ./...

4. Confirm the current derivation and integrity tests.
5. Establish the baseline before changing anything.

Do not assume a prior audit result remains true without executing the checks.

==================================================
5. GATE A — NO-PANIC BOUNDARY
==================================================

Implement/complete negative tests in existing files only:

ops/negative_test.go
session/session_test.go

Cover at minimum:

- Apply wrong arity
- unknown operation
- malformed/invalid object
- non-symbol variable for Substitute
- non-symbol variable for Differentiate
- non-symbol variable for Limit
- malformed Solve input
- non-relation Solve input
- unsupported Solve shape
- malformed BranchSet
- wrong branch count
- wrong positive/negative relationship
- invalid branch relation operator
- wrong branch target
- wrong zero value
- malformed Session.Step

Required behavior:

- typed error
- never panic
- no partial session mutation

Failed Session.Step must leave the draft and ledger unchanged.

Limit must require the contractually valid single-symbol variable.

SelectBranch must inspect and validate the complete structural contract before selecting a branch.

Do NOT build a general expression validator.

Only make production changes if the negative battery exposes a real production failure.

==================================================
6. GATE B — COMPLETE PROVENANCE MATRIX
==================================================

Test all 12 pure operations:

- Add
- Multiply
- Negate
- Pow
- Sqrt
- Substitute
- Simplify
- Differentiate
- Limit
- Solve
- Compare
- SelectBranch

Plus Session.Identify separately.

Required law:

clean/non-hypothesis input
    → DERIVED

any HYPOTHESIS input
    → HYPOTHESIS

Pure operations must never accidentally manufacture:

- IDENTIFIED
- POSTULATED
- DEFINED

Session.Identify remains the sole identification authority.

Required:

clean Identify → IDENTIFIED

hypothesis-contaminated Identify → HYPOTHESIS

Use the central provenance logic already present.

Do not add semantic prose-quality checks to Identify justification beyond existing non-empty/trimmed behavior.

==================================================
7. GATE C — CORPUS LOCK + REVERSE ALLOWLIST
==================================================

Strengthen manifest constructor cross-checking.

Per applicable manifest item, verify consistency for:

- canonical expression
- canonical expression hash
- kind
- dimension
- provenance
- corpus status
- assumptions
- source
- other already-required metadata

Both current frameworks remain ESTABLISHED.

Derived and candidate artifacts must never silently become ESTABLISHED.

Reverse mapping test:

For each domain package:

every ESTABLISHED-yielding exported constructor
    ↔
approved manifest item

every non-manifest constructor
    ∈
the pinned allowlist

Do NOT create inference rules.

==================================================
8. GATE D — THEORY-NEUTRAL KERNEL
==================================================

internal/kernel is generic formal semantic/mint machinery.

It must NOT contain physical theory content.

The only explicitly permitted physics-specific identifiers are:

- KindMinkowski
- LorentzFactorFunctionID / "lorentz_factor"

plus only their closed-world validation.

Do not embed special-relativity formulas in internal/kernel.

Do not embed mechanics laws in internal/kernel.

Do not add hidden physical constants, symbols, assumptions or framework logic.

The ops package may contain only the fixed Lorentz-factor body already required by the specification.

Reject symbol-name shortcuts such as:

symbol == "c"
symbol == "m"

when they would turn generic symbolic machinery into hidden physics-specific behavior.

==================================================
9. GATE H — THEORY / EPISTEMIC PACKAGE ISOLATION
==================================================

This is the critical additional freeze rule.

Treat the top-level package/folder as the isolation boundary for a physical framework.

Current architecture:

internal/kernel/
    = theory-neutral formal machinery

mechanics/
    = classical-mechanics corpus

relativity/
    = special-relativity corpus

hypothesis/
    = provisional/untrusted candidate space

Special relativity is NOT demoted to HYPOTHESIS.

Its corpus status remains ESTABLISHED.

The rule is:

An established physical framework remains explicitly bounded by its assumptions,
conventions, scope and empirical status; it does not become an ambient axiom of
the generic kernel.

Theory packages own their:

- physical objects
- relations
- assumptions
- postulates/framework premises
- conventions
- scope
- limitations
- anomalies
- derivations
- corpus metadata

Do not move these into internal/kernel.

Do not create a global theory-assumption registry.

Do not allow a result to acquire assumptions merely because the code executing it belongs to a particular package.

==================================================
10. GATE H1 — IMPORT BOUNDARY
==================================================

Use the existing AST-walk infrastructure in:

core/object_test.go

Assert:

internal/kernel non-test source imports ZERO local packages:

- mechanics
- relativity
- hypothesis
- core
- ops
- session

Preserve dependency direction:

kernel ← core ← ops ← session

with domain packages depending only on the permitted substrate.

No reverse dependencies.

Do not add new production dependencies merely to perform this test.

==================================================
11. GATE H2 — KERNEL PHYSICS SOURCE ALLOWLIST
==================================================

Extend the AST/source inspection test.

Non-test internal/kernel source may contain only:

- KindMinkowski
- LorentzFactorFunctionID
- "lorentz_factor"

and their explicitly required closed-world validation.

Reject additional theory-specific:

- symbols
- laws
- constants
- assumptions
- framework names

Do not reject the two explicitly permitted identifiers.

==================================================
12. GATE H3 — OPS ASSUMPTION SUBSET INVARIANT
==================================================

Prove assumption non-leakage at the pure-operation level.

For every pure operation:

output assumption keys ⊆
    input assumption keys
    UNION
    operation-explicitly-required/generated keys

Examples of legitimate operation-generated keys include:

- denominator preconditions
- selected_branch/*
- other keys explicitly required by the operation contract

An operation MUST NOT invent an unrelated framework assumption.

Enumerate each operation's explicitly required/generated keys in the implementation.

Do not use exact full assumption-set equality.

The invariant must permit legitimate assumption accumulation.

==================================================
13. GATE H4 — PACKAGE FORBIDDEN-KEY TESTS
==================================================

Prove cross-framework non-leakage at package level.

Do NOT pin exact complete assumption sets.

Instead assert forbidden intersections.

Mechanics results:

mechanics assumptions
    ∩
relativity-only keys
    = empty

Relativity results:

relativity assumptions
    ∩
mechanics-only keys
    = empty

Relativity-only keys currently include:

- rest_frame
- rest_mass_nonnegative
- speed_of_light_positive
- minkowski_spacetime
- lorentz_symmetry
- no_gravitational_dynamics
- special_relativistic_regime

Mechanics-only keys must be enumerated from the mechanics manifest.

Hypothesis assumptions remain hypothesis-associated.

No ESTABLISHED material may silently leak into candidate material.

Do both directions.

==================================================
14. GATE H5 — E=mc² ASSUMPTION OWNERSHIP
==================================================

In relativity/derivation_test.go verify that every assumption attached to the final m*c² result has a traceable source.

At minimum:

- rest_mass_nonnegative
- speed_of_light_positive
- selected_branch/*
- rest-frame-derived assumptions

must trace to:

- an explicit input object
- or a legitimate operation step

Do not accept assumptions merely because the final result is known.

Keep the no-MassEnergyRelation source assertion.

Keep Differentiate and Lorentz-limit derivations unhardcoded.

==================================================
15. GATE E — CANDIDATE EXPRESSION VALIDITY
==================================================

Verify trusted candidate validation rejects invalid expressions in:

Prediction.Relation

FalsificationCondition.ContradictingCondition

RecoveryClaim.Condition

Structural validity only.

Do not infer empirical truth.

Do not implement EBP runtime logic.

==================================================
16. GATE F — MANIFEST CORRUPTION BATTERY
==================================================

Use in-memory corruption only.

Do not modify production manifest files on disk.

Independently corrupt:

- canonical expression
- canonical expression hash
- dimension
- kind
- provenance
- corpus status
- assumptions
- source
- unknown constructor
- extra manifest item

Every corruption must be rejected.

==================================================
17. GATE G — SERIALIZATION
==================================================

`pow`:

- non-empty operator → reject
- non-empty irrelevant justification → reject according to frozen validation semantics
- canonical unused fields → empty/default

Rational canonicalization must cover:

1/2
-3/4
2 → 2/1
0 → 0/1

Reject JSON-number forms where exact string representation is required.

Do not rely on direct *big.Rat JSON marshaling.

Verify deterministic canonical bytes on repeated runs.

==================================================
18. REPLAY SAFETY
==================================================

Do not alter:

- InputCanonicals
- OutputCanonical
- ParamsCanonical
- stored hashes
- byte-based replay
- Session.Validate stage ordering

Canonical retained bytes remain authoritative for replay.

==================================================
19. MUTATION REGRESSION
==================================================

Run the complete existing M1–M18 mutation suite in a disposable copy.

Use:

-count=1

Restore the original implementation before each mutation.

Add cheap mutants for:

- accepting non-empty pow.operator
- allowing an extra ESTABLISHED constructor
- allowing an unapproved kernel physics identifier
- accepting invalid candidate Expr
- accepting malformed Limit shape
- accepting malformed SelectBranch shape
- allowing ambient cross-package assumption injection

Protect against regression of:

- hardcoded Limit result
- hardcoded E=mc²
- containment bypass
- operation compatibility bypass
- non-canonical object acceptance
- second mint path
- provenance promotion
- physics leakage

Report ONLY mutations actually executed.

==================================================
20. NON-BLOCKERS
==================================================

Do NOT expand scope for:

- justification prose quality
- absolute zero-physics kernel rule beyond the explicit allowlist
- go:embed file location
- cryptographic authenticity
- AI planner
- EBP runtime integration
- Einstein-1905 corpus

==================================================
21. FILE DISCIPLINE
==================================================

No new implementation files.

No renames.

No cross-package moves.

Use existing test files for Gate H:

- core/object_test.go → import boundary + kernel physics allowlist
- ops tests → assumption subset
- domain/hypothesis tests → forbidden-key checks
- relativity/derivation_test.go → E=mc² assumption ownership

Use:

- ops/negative_test.go
- session/session_test.go

for Gate A.

Use existing manifest tests for Gate F.

Production modifications are allowed ONLY when an executable test demonstrates a real violation.

Prefer the smallest production fix plus focused regression test.

==================================================
22. FINAL VERIFICATION
==================================================

Run:

go test ./...
go vet ./...

Then verify:

- 39-file tree
- package graph
- import boundaries
- kernel source allowlist
- theory/package isolation
- assumption non-leakage
- reverse constructor allowlist
- manifest integrity
- candidate containment
- candidate expression validity
- complete provenance matrix
- no-panic battery
- serialization pins
- byte replay
- E=mc² trace
- kinetic derivative
- Lorentz limit
- mutation campaign
- clean working tree
- no unintended dependencies

Do not declare PASS based only on code inspection when an executable test can establish the result.

==================================================
23. FINAL REPORT — EXACT FORMAT
==================================================

A. BASELINE

Report actual pre-change:

- go test ./...
- go vet ./...
- tree/package status

B. CHANGES IMPLEMENTED

For each Gate A–H:

- PASS / NOT NEEDED
- files changed
- production changes, if any
- tests added/changed
- concrete evidence

C. EXPLICIT NON-CHANGES

Confirm:

- specs_v2_3 unchanged
- Plan 10 preserved
- 39-file architecture preserved
- package graph preserved
- no new physics domains
- no Einstein-1905 corpus
- no parser
- no CAS
- no numerical engine
- no theorem prover
- no solver expansion
- no automatic promotion
- no truth adjudication
- no EBP runtime integration

D. VERIFICATION

Report actual results for:

- go test ./...
- go vet ./...
- E=mc² derivation
- no-panic battery
- provenance matrix
- manifest corruption battery
- reverse allowlist
- kernel physics boundary
- theory isolation
- assumption non-leakage
- candidate validation
- serialization
- replay
- mutations

E. FINAL STATUS

Use exactly:

FREEZE-READY

or:

NOT FREEZE-READY

If NOT FREEZE-READY, list only concrete blockers supported by failed tests or executable evidence.

Never declare FREEZE-READY based solely on inspection.

==================================================
24. FINAL ARCHITECTURAL PRINCIPLE
==================================================

The frozen MVP must establish:

AI agent
    ↓
formal reasoning substrate
    ↓
explicit theory objects
    ↓
bounded theory package
    ↓
mechanical derivation
    ↓
auditable result

The kernel provides reasoning machinery.

The theory package provides bounded physical knowledge.

The hypothesis package provides provisional candidate knowledge.

No physical framework becomes an ambient axiom of the kernel.

The primary proof-of-concept remains the AI-composable derivation:

EnergyMomentumRelation
→ ZeroThreeMomentum
→ Substitute
→ Simplify
→ Solve
→ Compare
→ SelectBranch
→ m*c²

Do NOT implement the future Einstein-1905 reconstruction now.

That is explicitly deferred until after the kernel is frozen.
```

This version is what I would hand to the implementation agent. Plan 7's own integration section already confirms the intended coexistence of Gates A–G with the new isolation rules, including the explicit no-new-files discipline and mutation additions. :chatgpt-content-reference{index="3"}

The conceptual architecture is now clean: **the kernel is not the physics; the theory package is the physics; the hypothesis package is provisional; the AI is the reasoner; and the ledger is the audit trail.**
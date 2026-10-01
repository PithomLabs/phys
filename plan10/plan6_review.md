I reviewed the attached **Plan 6 — Final Kernel Freeze / Hardening Pass** against the consolidated verdict and the additional theory/package-isolation constraint. The plan is already tightly aligned with the earlier freeze decisions: narrow hardening only, no Einstein-1905 reconstruction, explicit reverse allowlist, strict `pow` rejection, the full E=mc² operation chain, gates A–G, mutation regression, and the 39-file/no-redesign discipline. :chatgpt-content-reference{index="0"} :chatgpt-content-reference{index="1"}

The one substantive addition I would make is **Gate H: Theory/Epistemic Package Isolation**. That does not conflict with Plan 6; it makes its kernel boundary more rigorous. In particular, `relativity/` should remain a separate top-level theory package whose assumptions, postulates, conventions, corpus, scope and limitations are owned by that package, while `internal/kernel` remains theory-neutral except for the two explicitly permitted specification-mandated identifiers.

One important wording refinement: **do not classify special relativity itself as “a hypothesis”** merely because it is isolated this way. Its current corpus status can remain `ESTABLISHED`; the architectural rule is that an established framework is still **bounded by its explicit assumptions and falsifiability conditions** and must not become ambient truth inside the kernel.

Here is the consolidated implementation prompt I would now give the coding agent:

---

# IMPLEMENTATION PROMPT — PLAN 6 FINAL KERNEL FREEZE + THEORY ISOLATION

You are the implementation agent for the Go Physics Compiler MVP.

Execute **Plan 6 — Final Kernel Freeze / Hardening Pass**, incorporating the additional **Theory/Epistemic Package Isolation** requirement below.

This is the **last narrow hardening pass before MVP kernel freeze**.

## 1. Scope

This task is:

> **narrow hardening + mechanical proof only**

It is **not** a redesign and **not** a feature phase.

The following are explicitly out of scope:

- new physics domains
- General Relativity
- Quantum Mechanics
- Electromagnetism
- Einstein 1905 historical corpus reconstruction
- parser
- CAS
- numerical engine
- theorem prover
- general solver
- automatic hypothesis promotion
- truth scoring/adjudication
- EBP 2.1 runtime integration
- AI planner/orchestrator

Einstein's original 1905 mass-energy derivation remains **deferred until after kernel freeze**.

The attached Plan 6 is the base plan. Preserve its scope and decisions. :chatgpt-content-reference{index="2"}

---

# 2. Frozen architecture

Unless a concrete freeze gate below strictly requires a change, treat these as immutable:

- `specs_v2_3.md`
- Plan 10 architecture
- 39-file tree
- package boundaries
- dependency direction
- `internal/kernel` authority model
- `core` trusted type facade
- `ops` purity
- `session` replay/ledger model
- existing mechanics corpus
- existing special-relativity corpus
- existing E=mc² derivation path
- stdlib-only requirement

Do not add files.

Do not rename files.

Do not move packages.

Do not refactor unrelated code.

Do not broaden the corpus.

---

# 3. Locked decisions

These decisions are already made and must be implemented exactly.

### Reverse constructor allowlist

Use a **pinned explicit allowlist**:

```text
NewKineticEnergy
Velocity
ZeroThreeMomentum
ZeroEnergy
ZeroVelocity
hypothesis concept
```

Do **not** infer permission from provenance, corpus status, naming, or any other rule.

Every other `ESTABLISHED`-yielding trusted constructor must correspond to an approved manifest item.

### `pow` OperationParams

Strictly reject `pow` parameters containing a non-empty unused `operator`.

Canonical operational form:

```json
{"kind":"pow","exponent":"2/1","operator":"","justification":""}
```

The inconsistent textual `"operator":"eq"` example in v2.3 must not override this implementation invariant.

Do not modify the specification text.

---

# 4. PRIMARY ACCEPTANCE TARGET — E=mc²

Do not weaken the existing primary architecture test.

It must execute:

```text
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
```

using actual symbolic operations.

It must **never call `MassEnergyRelation()` to obtain the answer**.

Harden the no-call rule with a mechanical source assertion in:

```text
relativity/derivation_test.go
```

The existing kinetic-energy derivative test must remain:

```text
Differentiate(KineticEnergy, v) → m*v
```

The existing Lorentz-limit test must remain:

```text
Limit(LorentzFactor, v, 0) → 1
```

Neither may rely on hardcoded output behavior.

Plan 6 explicitly establishes this acceptance target. :chatgpt-content-reference{index="3"}

---

# 5. GATE A — NO-PANIC OPERATION BOUNDARY

Add/complete negative-test coverage in existing files only.

Test:

- wrong-arity `ops.Apply`
- unknown operation
- invalid object
- non-symbol variable for `Substitute`
- non-symbol variable for `Differentiate`
- non-symbol variable for `Limit`
- malformed/non-relation `Solve`
- malformed/non-quadratic `Solve`
- malformed `BranchSet`
- wrong branch count
- incorrect negation relationship
- incorrect relation operator
- incorrect target
- incorrect zero value
- invalid `Session.Step`

For every malformed case:

```text
error, not panic
```

For failed `Session.Step`:

```text
draft remains unchanged
no partial ledger append
```

`SelectBranch` must validate its complete v2.3 structural contract rather than blindly selecting branch 0.

`Limit` must require the appropriate single-symbol variable representation.

Do not build a general expression validator.

Production changes are justified only if the negative battery reveals a genuine panic or boundary violation.

Plan 6 defines this gate and its intended battery. :chatgpt-content-reference{index="4"}

---

# 6. GATE B — 12-OP PROVENANCE MATRIX

Test all 12 pure operations:

```text
Add
Multiply
Negate
Pow
Sqrt
Substitute
Simplify
Differentiate
Limit
Solve
Compare
SelectBranch
```

plus `Session.Identify` separately.

Required law:

```text
clean trusted/non-hypothesis inputs
    → DERIVED

any HYPOTHESIS input
    → HYPOTHESIS
```

Pure operations must never produce:

```text
IDENTIFIED
POSTULATED
DEFINED
```

as an accidental result provenance.

`Session.Identify` remains the sole identification authority.

Test:

```text
Session.Identify(clean, clean) → IDENTIFIED
Session.Identify(hypothesis, ...) → HYPOTHESIS
```

Do not introduce justification prose-quality semantics beyond the existing non-empty/trimmed rule.

Use the existing central provenance implementation rather than creating parallel logic.

Plan 6 defines this matrix. :chatgpt-content-reference{index="5"}

---

# 7. GATE C — CORPUS CONTENT LOCK + REVERSE ALLOWLIST

Strengthen existing manifest cross-checking.

For every applicable trusted item verify:

- canonical expression
- canonical expression hash
- dimension
- kind
- provenance
- corpus status
- assumptions
- source
- relevant existing metadata

Both current frameworks remain `ESTABLISHED`.

Derived/candidate artifacts must not silently become `ESTABLISHED`.

### Reverse allowlist

For each domain package:

```text
every ESTABLISHED-yielding exported constructor
    ↔
manifest item

every non-manifest constructor
    ∈
pinned explicit allowlist
```

Do not infer constructor permission automatically.

Plan 6 explicitly locks this model. :chatgpt-content-reference{index="6"}

---

# 8. GATE D — THEORY-NEUTRAL KERNEL + THEORY PACKAGE ISOLATION

This gate incorporates the additional architectural requirement.

## 8.1 Kernel boundary

`internal/kernel` is the generic formal semantic/mint substrate.

It must contain **no physical theory corpus**.

The only physics-specific identifiers permitted there are the already specification-mandated closed-world identifiers:

```text
KindMinkowski
LorentzFactorFunctionID = "lorentz_factor"
```

and only the minimal validation necessary to support those mandated identifiers.

Do not add special-relativity equations or physics assumptions to the kernel.

Do not add mechanics laws to the kernel.

Do not add future TOE concepts to the kernel.

## 8.2 Theory package is the unit of separation

Treat each top-level physical framework package/folder as an explicit epistemic/architectural boundary.

Current structure:

```text
mechanics/
relativity/
hypothesis/
internal/kernel/
```

Conceptually:

```text
internal/kernel
    = formal machinery

mechanics
    = classical mechanics framework/corpus

relativity
    = special relativity framework/corpus

hypothesis
    = provisional/untrusted candidate space
```

The `relativity/` package owns its:

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

Do not move any of these into `internal/kernel`.

## 8.3 Important epistemic distinction

Do **not** reinterpret package isolation as saying special relativity has `HYPOTHESIS` corpus status.

The baseline relativity corpus may remain:

```text
CorpusStatus = ESTABLISHED
```

The architectural rule is instead:

> An established physical framework remains explicitly bounded by its assumptions, conventions, scope and empirical status; it does not become an ambient axiom of the generic kernel.

The compiler remains a formal/integrity substrate, not a truth adjudicator.

## 8.4 No ambient assumption leakage

A result must not acquire relativity assumptions merely because code happens to execute inside or alongside the `relativity` package.

Likewise mechanics assumptions must not appear in unrelated results.

Framework assumptions must enter derivations through explicit trusted objects/relations/constraints.

No hidden global physical-assumption registry.

## 8.5 Import-boundary test

Add structural regression coverage using existing test infrastructure so that:

```text
internal/kernel
    MUST NOT import
        mechanics
        relativity
        hypothesis
```

The direction must remain:

```text
kernel
   ↑
core
   ↑
ops
   ↑
domain/session layers as permitted by architecture
```

Do not introduce reverse dependencies.

## 8.6 Kernel physics source allowlist

Extend the existing kernel source inspection test so that non-test files under:

```text
internal/kernel/
```

may contain only the explicitly permitted physics-specific identifiers.

Reject future accidental insertion of theory-specific symbols or laws.

Preserve the two specification-mandated identifiers.

Plan 6 already requires the kernel physics allowlist; this gate adds the explicit theory-package/epistemic interpretation and assumption-isolation tests. :chatgpt-content-reference{index="7"}

---

# 9. GATE E — TRUSTED CANDIDATE EXPRESSION VALIDITY

For trusted candidate validation, embedded expressions must be structurally valid.

Test malformed:

```text
Prediction.Relation
FalsificationCondition.ContradictingCondition
RecoveryClaim.Condition
```

and verify candidate validation/sealing fails.

This is structural validity only.

Do not interpret this as empirical truth validation.

Do not add EBP implementation.

Plan 6 already identifies this exact gate. :chatgpt-content-reference{index="8"}

---

# 10. GATE F — MANIFEST CORRUPTION BATTERY

In existing manifest tests only, exercise independent in-memory corruption of:

```text
canonical expression
canonical expression hash
dimension
kind
provenance status
corpus status
assumptions
source
unknown constructor
extra manifest item
```

Each corruption must be rejected by existing validation/cross-check machinery.

Never mutate production corpus files or disk state.

Plan 6 defines this battery. :chatgpt-content-reference{index="9"}

---

# 11. GATE G — SERIALIZATION AND CANONICALIZATION PINS

Preserve exact deterministic serialization.

### `pow`

Reject:

```json
{"kind":"pow","exponent":"2/1","operator":"eq","justification":""}
```

Accept:

```json
{"kind":"pow","exponent":"2/1","operator":"","justification":""}
```

Add regression tests.

### Rational serialization

Test exact canonical representations for:

```text
1/2
-3/4
2 → 2/1
0 → 0/1
```

Reject JSON-number representations where strings are required.

Do not rely on direct `*big.Rat` JSON marshaling.

### Determinism

Verify repeated canonicalization produces identical bytes.

Do not introduce map-order dependence.

Plan 6 specifies these pins. :chatgpt-content-reference{index="10"}

---

# 12. REPLAY SAFETY

Do not alter:

```text
InputCanonicals
OutputCanonical
ParamsCanonical
```

or their hashes.

Preserve byte-based replay.

Preserve the existing `Session.Validate` stage ordering and semantics.

Do not optimize away retained canonical bytes.

Plan 6 explicitly says this remains untouched. :chatgpt-content-reference{index="11"}

---

# 13. MUTATION REGRESSION

Rerun the full existing M1–M18 mutation suite in a disposable copy.

Use:

```bash
-count=1
```

and restore the original implementation between mutants.

Also add inexpensive mutations corresponding to the new gates:

```text
accept non-empty pow.operator
allow extra ESTABLISHED constructor
allow unapproved kernel physics identifier
accept invalid candidate expression
accept malformed Limit shape
accept malformed SelectBranch shape
allow ambient cross-package assumption
```

Protect especially against regressions for:

- hardcoded Limit result
- hardcoded E=mc² result
- candidate containment bypass
- operation-compatibility bypass
- non-canonical object acceptance
- second mint path
- provenance promotion
- kernel physics leakage

Report only actually executed mutation results.

Plan 6 requires actual rerun results, not inspection claims. :chatgpt-content-reference{index="12"}

---

# 14. NON-BLOCKERS — DO NOT EXPAND SCOPE

Do not turn these into implementation work:

- justification prose quality
- absolute “zero physics references” rule beyond the explicit kernel allowlist
- location of `go:embed`
- cryptographic authenticity
- AI planning/orchestration
- EBP 2.1 runtime integration
- Einstein's 1905 derivation

Plan 6 explicitly identifies these as non-blockers. :chatgpt-content-reference{index="13"}

---

# 15. FILE DISCIPLINE

No new files.

No renames.

No cross-package moves unless an actual architecture violation requires one.

Prefer:

```text
small production fix
+
focused regression test
```

over broad refactoring.

Any change outside a named freeze gate is out of scope.

:chatgpt-content-reference{index="14"}

---

# 16. FINAL VERIFICATION

Run:

```bash
go test ./...
go vet ./...
```

Then verify:

```text
39-file tree
package graph
import boundaries
kernel physics allowlist
theory/package isolation
manifest integrity
reverse constructor allowlist
candidate containment
candidate expression validity
12-op provenance matrix
no-panic battery
serialization pins
byte replay
E=mc² trace
mutation results
clean working tree
no unintended dependencies
```

Do not declare success from inspection when executable verification is available.

---

# 17. FINAL REPORT

Return exactly these sections.

## A. Baseline

Actual results before changes.

## B. Changes implemented

For each gate A–H:

```text
PASS / NOT NEEDED
files changed
implementation change
tests added/updated
```

## C. Explicit non-changes

Confirm:

```text
spec v2.3 unchanged
Plan 10 architecture preserved
39-file tree preserved
package graph preserved
no new physics domains
no Einstein-1905 corpus
no parser
no CAS
no numerical engine
no theorem prover
no solver expansion
no automatic promotion
no truth adjudication
no EBP runtime integration
```

## D. Verification

Report actual results of:

```text
go test ./...
go vet ./...
E=mc² derivation
no-panic battery
provenance matrix
manifest corruption battery
reverse allowlist
kernel physics boundary
theory/package isolation
candidate expression validity
serialization tests
replay validation
mutation campaign
```

## E. Final status

Use exactly one:

```text
FREEZE-READY
```

or:

```text
NOT FREEZE-READY
```

If `NOT FREEZE-READY`, list only concrete blockers supported by failed tests or executable evidence.

---

# FINAL DESIGN PRINCIPLE

The frozen MVP must establish a clean separation:

```text
             AI AGENT
                 │
                 │ composes derivation
                 ▼
        ┌─────────────────┐
        │ formal substrate│
        │  internal/core  │
        │      + ops      │
        └────────┬────────┘
                 │
       explicit theory objects
                 │
        ┌────────┴────────┐
        ▼                 ▼
   mechanics/         relativity/
   corpus             corpus
                         │
                 explicit assumptions
                 explicit scope
                 explicit provenance
                 explicit limitations
```

The kernel supplies **reasoning machinery**, not physical truth.

A theory package supplies **bounded physical knowledge**.

A hypothesis package supplies **provisional candidate knowledge**.

The future TOE layer will eventually connect these packages through explicit bridges and evaluations.

**Do not implement that future layer now.**

The immediate objective is to leave the kernel sufficiently deterministic, isolated, auditable and mutation-resistant that the subsequent AI-reasoning experiment—especially the blind, operation-composed derivation of \(E=mc^2\)—has a trustworthy substrate.
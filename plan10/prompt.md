You are revising the attached **Plan 10 implementation plan** for the Physics Compiler MVP.

## Governing source of truth

Read the complete attached files before editing:

1. `specs_v2_3.md` — **normative and authoritative**
2. `plan10.md` — the current Plan 10 to revise
3. The adversarial review supplied below — use only findings that remain valid after checking them against the complete v2.3 specification.

The final revised plan must be governed by **`specs_v2_3.md`**, not v2.2, `plan9.md`, or any historical specification.

### Primary objective

Produce a revised **Plan 10 v2.3** that:

* preserves the existing Plan 10 architecture wherever it is compatible with v2.3;
* preserves useful execution detail, test specificity, and dependency-ordered implementation sequencing;
* incorporates all **net-valid** findings below;
* removes stale, contradictory, or v2.2-only material;
* does **not** introduce architectural refactoring merely to make Plan 10 look like Plan 9.2;
* leaves the coding agent with one deterministic architecture and no architectural choices to invent;
* remains within the exact 39-file repository tree and ≤40-file MVP budget;
* does not modify the normative specification.

This is a **precision/conformance revision, not an architectural redesign**.

Do not add packages, files, abstractions, registries, frameworks, generic factories, CAS functionality, numerical execution, or future scaffolding unless the v2.3 specification explicitly requires them.

---

# 1. First, establish the invariant architecture

The following architecture is already sound and should be PRESERVED unless you find a direct contradiction in v2.3:

```text
internal/kernel  ← concrete immutable authority
       ↑
core             ← public aliases + façade
       ↑
ops              ← pure symbolic transformations
       ↑
session          ← ledger/replay/seal authority

mechanics   → core + internal/kernel
relativity  → core + internal/kernel
hypothesis  → core + internal/kernel
```

Preserve these architectural properties:

* `core` MUST NOT import `ops` or `session`.
* `ops` MUST NOT import `session`.
* `session` MUST NOT import `mechanics`, `relativity`, or `hypothesis`.
* `session` may use `core`, `ops`, and `internal/kernel`.
* `internal/kernel` owns the concrete immutable representation.
* `core.Object` and `core.Expr` remain aliases over kernel types.
* `kernel.MintObject` remains the sole production object mint.
* `ops` remains pure/stateless.
* `Session.Identify` remains the only identification API.
* `session` remains the owner of ledger/replay/sealing.
* `ResearchCandidate` remains session-owned.
* hypothesis contamination remains enforced.
* no generic public object factory exists.
* no truth-adjudication API exists.
* no MRC bypass/override exists.

Do NOT introduce `core.SealedDerivation`, a new `review/` package, a domain registry, or any alternative layering.

---

# 2. Replace all v2.2 authority language with v2.3

The current Plan 10 begins by treating v2.2 / `plan9.md` as normative.

That MUST be corrected.

The revised plan must explicitly state:

```text
Normative source:
specs_v2_3.md
```

and must describe v2.3 as the sole authoritative specification.

All references to v2.2 or `plan9.md` must either:

* be removed; or
* be explicitly marked historical/background only.

Do not leave statements such as:

```text
Normative source: plan9.md = specs v2.2
```

anywhere in the revised plan.

---

# 3. Reconcile all v2.3-specific changes

The revised Plan 10 MUST explicitly incorporate the three reconciliations that v2.3 introduced:

### A. Challenge / Review ownership

v2.3 explicitly establishes:

```text
core/corpus.go
    ↓
Challenge
Review

session
    ↓
consumes core.Challenge / core.Review
```

Therefore:

* `Challenge` and `Review` are defined exclusively in `core/corpus.go`.
* `session` MUST NOT redefine them.
* `session` MUST NOT introduce aliases that imply a second canonical declaration.
* no `review/` package exists.

Do not retain the current Plan 10 wording:

```text
Challenge/Review re-declared as aliases of core types
```

Replace it with the exact v2.3 model: `session` consumes `core.Challenge` and `core.Review`.

### B. Finite symbolic factors / negative integer powers

Use the v2.3 §9.6 rule exactly.

The implementation plan must distinguish:

```text
base != 0
    → non-singularity/admissibility

base > 0
    → sign-sensitive rewrite
```

Do not introduce a Lorentz-specific simplification hack.

The general finite-factor rule must allow the required `c^-1` case under `c > 0`, while retaining:

```text
0^0 → UnsupportedOperationError
```

and rejecting unsafe negative powers such as `0^-1`.

### C. Normative `ops.Apply` positional mapping

The plan MUST include the exact v2.3 §15.13.1 mapping for all 12 operations and explicitly state that:

* `ops.Apply`;
* `Session.Step`;
* replay in `Session.Validate`

must use exactly the same positional mapping.

Inputs MUST NOT be reordered.

`identify` is excluded from `ops.Apply`.

Also include the exact `OperationParams` schema from v2.3 §15.13:

```go
type OperationParams struct {
    Kind          string
    Exponent      string
    Operator      core.RelationOperator
    Justification string
}
```

with canonical JSON:

```json
{
  "kind":"pow",
  "exponent":"2/1",
  "operator":"eq",
  "justification":""
}
```

and permitted kinds:

```text
empty
pow
compare
identify
```

Pin the required parameter behavior:

* `pow` → `Exponent`
* `compare` → `Operator`
* `identify` → `Justification`
* all other operations/assertions → `empty`
* `identify` never enters `ops.Apply`

Do not leave `OperationParams` opaque.

---

# 4. Apply the net-valid adversarial findings

Treat the following five findings as valid and incorporate them.

## F1 — Fix relativity manifest consistency

v2.3 §28.2 requires exactly this 10-item relativity manifest baseline:

```text
Spacetime
MinkowskiMetric
RestMass
Energy
ThreeMomentum
FourMomentum
SpeedOfLight
LorentzFactor
EnergyMomentumRelation
MassEnergyRelation
```

`Velocity` is a required relativity wrapper/constructor, but it is NOT one of the ten required manifest items.

Therefore make the plan internally consistent:

* 8 nominal wrappers in the ordinary wrapper surface:
  `Spacetime`, `MinkowskiMetric`, `RestMass`, `Energy`,
  `ThreeMomentum`, `FourMomentum`, `SpeedOfLight`, `Velocity`
* fixed relation/function constructors:
  `LorentzFactor`, `EnergyMomentumRelation`, `MassEnergyRelation`
* zero constructors:
  `ZeroThreeMomentum`, `ZeroEnergy`, `ZeroVelocity`
* exactly 10 required manifest items as specified above
* no ambiguity over whether `Velocity` is a manifest item

The manifest constructor mapping must reflect the exact v2.3 §28.2 set.

## F2 — Add `Pow(x,0)` to Simplify

v2.3 §8.7 explicitly requires:

```text
Pow(x,0) → 1
```

only when `x` is a valid nonzero-safe base under the operation's assumptions.

Add this to the Plan 10 simplification rules.

Retain:

```text
0^0 → UnsupportedOperationError
```

Do not turn `Pow(x,0)` into an unconditional rewrite.

Keep all §9.6 exact rational-power rules intact.

## F3 — Make kernel-side metadata construction explicit

The plan currently implies metadata constructors without specifying how immutable kernel values can actually be created.

Make the implementation sequence explicit.

Because kernel metadata types have unexported authoritative fields, Step 1 MUST establish the kernel-side construction primitives required by the public core façade for:

* `Dimension`
* `Assumption`
* `AssumptionSet`
* `Convention`
* `ConventionSet`
* `Provenance`

Then Step 2's public `core` constructors delegate to those kernel-side constructors/helpers.

State explicitly that validation/invariant logic for each metadata type exists in ONE authoritative place in `internal/kernel`; `core` does not duplicate divergent validation rules.

This is an implementation detail inside the already-fixed architecture, not a new package.

## F4 — Pin deterministic set ordering

v2.3 requires deterministic canonical artifacts and hashes, but does not leave set ordering implementation-dependent.

The revised plan MUST pin one deterministic representation/order for:

```text
AssumptionSet
ConventionSet
```

Use a clear deterministic canonical rule; preferred implementation:

```text
canonicalize each element
→ obtain canonical element bytes
→ sort lexicographically by canonical bytes
→ store in deterministic slice order
```

Exact duplicates are deduplicated.

For assumptions, conflict detection remains by:

```text
(Kind, Key)
```

with differing canonical values producing `AssumptionConflictError`.

For conventions, conflict detection remains:

```text
same Key + different Value
    → ConventionConflictError
```

Do not use Go map iteration to determine canonical order.

Explicitly state that the deterministic ordering participates in:

```text
canonical JSON
→ set hash
→ object hash
→ ledger hash chain
→ candidate determinism
```

## F5 — Fix Test D

The current Test D example mentioning KineticEnergy vs Momentum is dimensionally invalid.

Replace it with a genuinely equal-dimension / different-kind pair.

Preferred examples:

```text
Mass vs RestMass
```

or:

```text
Velocity vs SpeedOfLight
```

or:

```text
Momentum vs ThreeMomentum
```

or:

```text
Position vs Spacetime
```

Use one exact pair consistently throughout Plan 10.

The test must isolate MRC-003 and not accidentally trigger MRC-002.

---

# 5. Incorporate these additional precision improvements

These are also net-valid and should be included.

### Simplify constructor boundary

Explicitly document:

```text
Constructor-time:
    flatten Add/Mul
    combine exact rationals
    sign normal form
    deterministic ordering

Simplify-time:
    Add/Mul identities
    zero/one rules
    rational Pow rewrites
    repeated powers
    Sqrt rewrites
    relation recursive simplification
```

Resolve the v2.3 §8.5 / §9.6 boundary explicitly rather than leaving it implicit.

### Simplify on Relation and BranchSet

Because the E=mc² path performs `Substitute` followed by `Simplify` on a relation, explicitly state:

* `Relation` is recursively simplified on both sides while preserving the relation operator;
* `BranchSet` is recursively simplified within its bounded structure as required;
* `Simplify` never creates `IDENTIFIED`;
* no relation simplification is interpreted as physical identification.

Do not expand this into general equation reasoning.

### SelectBranch

Use the exact v2.3 §15.12 contract:

```text
BranchSet(target,[positiveBranch,Neg(positiveBranch)])
+
Relation(gte,target,zero)
```

where `zero` has matching kind and dimension.

The operation must:

1. validate branch-set shape;
2. validate constraint compatibility;
3. merge the constraint into assumptions;
4. simplify the selected branch under current assumptions;
5. preserve provenance subject to contamination;
6. return `Expression` with target dimension.

Use:

```text
selected_branch/<HashExpr(constraint.Expr())>
```

as the deterministic assumption key.

Do not invent a broad branch-selection engine.

If v2.3 does not explicitly specify a separate sign-entailment-failure error, do not invent one merely to fill the gap. State the narrow behavior required by the specified mass-energy path and use existing `UnsupportedOperationError` only where the broader operation contract already requires it.

### Defensive copies

Explicitly mention defensive copies for:

* `NewRational` input `*big.Rat`;
* `NewPow` input exponent;
* `Pow` caller-supplied exponent;
* `Expr.RationalValue()`;
* `Expr.Exponent()`;
* slice-returning accessors.

No caller-owned mutable pointer may become part of immutable kernel state.

### Kernel helper ownership

Explicitly state where these foundational mechanisms are implemented:

```text
EntailsNonNegative
EntailsPositive
EntailsNonZero
AssumptionSet.Merge
ConventionSet.Merge
metadata equality
metadata canonicalization
metadata hashing
```

Prefer `internal/kernel/types.go` within the fixed file budget.

Do not add helper packages/files.

### Object decoding boundary

Do NOT introduce a public `DecodeObjectJSON` or `ParseObject` API.

Object canonical decoding should remain internal to `internal/kernel` and used by session replay/internal validation.

The decoder must validate the same authoritative invariants needed for a valid object before the object enters replay.

A canonical byte sequence alone is not sufficient if its decoded object violates kernel invariants.

### `schema_version`

Do NOT treat `"1"` as a v2.3 normative requirement unless the attached v2.3 spec explicitly establishes that value.

The spec requires `schema_version` in the canonical object/artifact structures, but the implementation plan must distinguish:

```text
authoritative kernel storage
```

from

```text
canonical serialization DTO/artifact representation
```

Do not add serialization-only `schema_version` state to `kernel.Object` merely to satisfy JSON ordering.

If retaining a `"1"` value as an implementation pin, label it explicitly as a **Plan 10 implementation pin**, not as a statement that v2.3 itself mandates `"1"`.

### Canonical codec visibility

Do not write that `core` calls unexported Go identifiers inside `internal/kernel`.

Use a valid Go visibility contract such as:

```text
internal/kernel
    exports controlled canonical helper functions to in-module packages

core
    exposes the public canonical façade
```

while keeping object decoding internal-only where appropriate.

---

# 6. Remove stale v2.2 material

Perform a systematic scan of Plan 10 for v2.2 residue.

Particularly remove/correct:

* `plan9.md` as normative source;
* v2.2-only assumption-kind lists;
* any old Challenge/Review placement;
* any old `OperationParams` ambiguity;
* any old manifest item counts;
* any stale `REQ-*` references whose semantics changed in v2.3;
* any text that contradicts the v2.3 specification.

Do not blindly preserve old Plan 10 text merely because it sounds plausible.

For every contract, v2.3 wins.

For example, the current Plan 10 coverage material mentioning assumption kinds such as:

```text
fixed_constant
state_relation
framework_relation
```

must be reconciled with the actual v2.3 assumption-kind set:

```text
Domain
Regime
Constraint
Convention
Approximation
MathPrecondition
PhysicalAssumption
```

and corresponding canonical JSON values:

```text
domain
regime
constraint
convention
approximation
math_precondition
physical_assumption
```

Do not carry forward v2.2 semantics that conflict with v2.3.

---

# 7. Rebuild the coverage matrix against v2.3

This is important.

Plan 10 already has a strong coverage matrix; preserve its level of detail.

However, rebuild/check it against the COMPLETE `specs_v2_3.md`.

The revised matrix MUST:

* cover every explicit `REQ-*`;
* expand `REQ-002-01..26` individually;
* cover every normative unnumbered `MUST` / `MUST NOT`;
* assign deterministic coverage IDs where needed;
* cover every MRC-001..008;
* cover A–S;
* provide implementation location;
* provide at least one named test;
* include the v2.3 §15.13.1 positional mapping;
* include §§18–41 where normative clauses exist;
* contain no references to obsolete v2.2-only requirements.

Do not merely copy the old Plan 10 matrix.

Use the actual attached v2.3 spec as the authoritative checklist.

At the end, verify mechanically that no normative v2.3 clause is left unmapped.

---

# 8. Preserve Plan 10's execution strengths

Keep Plan 10's useful structure:

```text
purpose
dependencies
files
result
gate
```

Keep its dependency-ordered implementation sequence.

Keep the explicit test-package plan.

Keep the A–S acceptance inventory.

Keep the final repository/file-budget audit.

Keep source/AST/import audits where they materially enforce the specification.

Keep the explicit determinism run:

```text
same derivation twice
+
fresh Session instance
```

Keep the anti-hardcoding checks for the Lorentz-factor body and the E=mc² derivation.

The objective is to make Plan 10 **cleaner and more exact**, not shorter at the expense of implementation determinism.

---

# 9. Important separation from Plan 9.2

Do NOT attempt to merge Plan 10 into Plan 9.2.

Plan 10 is intentionally an independent implementation plan.

It should remain independently recognizable while converging on the same v2.3 normative architecture.

Do not copy Plan 9.2 wholesale.

Use Plan 9.2 only as a secondary sanity reference where helpful; the attached `specs_v2_3.md` remains the authority.

---

# 10. Final consistency checks before writing the revised plan

Before considering the revision complete, perform these checks:

### Architecture

Confirm:

```text
kernel → core → ops → session
domains → core + kernel
no cycles
```

### Files

Confirm exactly:

```text
39 files total including go.mod
```

and no file outside the v2.3 §3 tree.

### Operations

Confirm exactly 12 pure `ops` operations:

```text
Add
Subtract
Multiply
Divide
Pow
Simplify
Substitute
Differentiate
Limit
Compare
Solve
SelectBranch
```

and `identify` is excluded from `ops.Apply`.

### Session

Confirm exactly nine session actions:

```text
Draft
Postulate
Declare
Step
Identify
Commit
Conclude
Seal
```

with the exact v2.3 lifecycle and read-only validation methods.

### Candidate

Confirm:

```text
ParseResearchCandidateJSON
    → UnverifiedResearchCandidate only

UnverifiedResearchCandidate.Validate()
    → trusted ResearchCandidate only after full validation
```

and no direct trusted external JSON construction.

### Relativity

Confirm the exact ten required manifest items and separate `Velocity` wrapper.

### Simplifier

Confirm:

```text
Pow(x,0)
```

is present with the nonzero-safe gate and:

```text
0^0
```

remains unsupported.

### Canonicalization

Confirm deterministic ordering for:

```text
Expr
AssumptionSet
ConventionSet
Object
Manifest
OperationParams
Step
StepEnvelope
Ledger
ResearchCandidate
```

with no map-order dependence.

### Canonical Object

Confirm:

```text
schema_version
valid
name
kind
dimension
expr
assumptions
conventions
provenance
corpus_status
```

in exact canonical order, while keeping serialization-only metadata separate from authoritative kernel storage.

### Final wording

The revised plan MUST NOT claim:

```text
Open Spec Items: NONE
```

until you have actually verified the complete v2.3 specification and coverage matrix.

Once complete, the plan SHOULD state that the architecture is fixed and the coding agent has no remaining architectural choices.

---

# 11. Deliverable

Produce the revised document as:

```text
plan10_v2_3.md
```

Do not modify `specs_v2_3.md`.

Do not modify Plan 9.2.

Do not implement code.

Do not create any repository source files.

The only deliverable in this task is the revised implementation plan.

At the end of the revised plan, include a concise:

```text
## Revision Audit

### Preserved Architecture
...

### Net-valid Review Findings Adopted
...

### Stale v2.2 Material Removed
...

### v2.3 Coverage Verification
...

### Open Spec Items
NONE
```

Use `NONE` only after checking the complete attached v2.3 specification.

The resulting Plan 10 must be independently handoff-ready: a coding agent reading only `specs_v2_3.md` + `plan10_v2_3.md` should not need to choose an architecture, invent missing APIs, infer positional operation mappings, guess manifest membership, or resolve contradictory implementation instructions.

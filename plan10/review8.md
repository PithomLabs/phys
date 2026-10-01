The adversarial review found **two real freeze blockers**, and I agree with both. The rest of the core architecture survived the attack.

### What the review established

The strongest part is that it independently confirmed the important invariants:

- single mint authority
- immutable objects
- package/dependency boundaries
- canonicalization and exact rationals
- provenance and hypothesis containment
- assumption propagation
- theory-package isolation
- manifest integrity
- replay integrity
- `OperationParams` strictness
- malformed-input boundary
- E=mc² anti-hardcoding
- Lorentz-limit and derivative paths

The reviewer also found no test theater and no hidden physics dispatch. adv_review8

So this is **not a fundamental architectural failure**. It is a final freeze-discipline failure.

## Blocker 1 — repository/file-set violation

The repository has **43 files where the frozen specification requires 39**.

The concrete extras identified are:

```text
.gut
reality_check.md
reality_first_plain_english.md
reality-first-physics-v8.0.md
.gitignore
```

while the specified:

```text
docs/paper-translation.md
```

is missing. adv_review8

This is a legitimate blocker because the specification explicitly freezes the implementation tree.

The implementation agent should therefore **not simply delete everything until it reaches 39**. It must reconcile each item against the normative 39-file list:

```text
delete accidental .gut
delete the three future-scaffolding/reality documents
decide .gitignore according to the exact spec/file-list rule
restore/create docs/paper-translation.md
verify exactly 39
```

The reviewer is correct that this must be resolved before freeze. adv_review8

## Blocker 2 — `DimensionEnergySquared()` public API leakage

This is also valid.

The specification's public dimension API has exactly nine dimension constructors; the implementation exposes:

```go
DimensionEnergySquared()
```

publicly in `core`, and `relativity` uses it. adv_review8

That violates the frozen public surface.

The proper fix is exactly the one suggested:

```text
core.DimensionEnergySquared()
        ↓ remove public export

internal/kernel
    private/kernel-only helper
        ↓
relativity/relations.go
```

The reviewer correctly observes that the helper is useful internally but does not belong in the public contract. adv_review8

### Two additional cleanup items

These aren't freeze blockers, but should be fixed in the same tiny pass:

1. Correct:

```go
expected 10 domains
```

to:

```go
expected 9 domains
```

because the actual condition is `len(m.Domain) != 9`. adv_review8

2. Change domain wrapper signatures from:

```go
CoreObject() kernel.Object
```

to the specification-facing:

```go
CoreObject() core.Object
```

Since `core.Object` is an alias, behavior doesn't change; this is API/spec conformity. adv_review8

## One thing I would *not* adopt from the review

The reviewer accepted the claimed mutation results rather than reproducing them because of the read-only review constraint. adv_review8

So we should **not regard 25/25 as independently verified** yet. It remains an implementation-agent claim until rerun in a disposable copy.

That distinction matters for the final freeze.

---

# Prompt for the implementation agent

Use this as the **final corrective pass after the adversarial review**:

```text id="ku0u1o"
FINAL CORRECTIVE PASS — POST-ADVERSARIAL REVIEW

The independent adversarial review reported:

FAIL — FREEZE DOES NOT HOLD

Two concrete blockers were found:

1. frozen 39-file implementation tree is violated
2. public DimensionEnergySquared() exceeds the frozen public API

The reviewer independently found the core semantic/security invariants otherwise holding.

Your task is to correct ONLY these concrete findings plus the two minor API/test-quality issues below, then rerun the complete freeze verification.

==================================================
1. NORMATIVE AUTHORITY
==================================================

specs_v2_3.md is authoritative.

Plan 10, Plan 7, and Plan 8 are additive freeze requirements.

Do not redesign anything.

Do not add features.

Do not add physics.

Do not add the Einstein-1905 corpus.

Do not modify the E=mc² reasoning architecture.

==================================================
2. BLOCKER A — RESTORE EXACT 39-FILE TREE
==================================================

The actual module tree was independently found to contain 43 files instead of the exact frozen 39-file implementation tree.

Reconcile the repository against the exact file list in specs_v2_3.md §3.

Known extras reported:

- .gut
- reality_check.md
- reality_first_plain_english.md
- reality-first-physics-v8.0.md
- .gitignore

Known missing required file:

- docs/paper-translation.md

Do not reach 39 by arbitrary deletion.

Perform an exact file-by-file comparison against the normative 39-file list.

Required outcome:

- every normative file exists
- no extra implementation/documentation/helper/scaffolding file exists inside the counted tree
- exactly 39 files
- no accidental generated artifacts
- no accidental editor/VCS artifact counted contrary to the frozen rule

For .gitignore specifically, resolve it strictly according to the normative spec.
Do not silently change the file-count definition.

Do NOT modify specs_v2_3 merely to accommodate the current repository.

==================================================
3. BLOCKER B — REMOVE PUBLIC DimensionEnergySquared()
==================================================

The public core API currently exports:

    DimensionEnergySquared()

This is not part of the frozen §7.0 public dimension constructor surface.

Remove it from the public core API.

Do NOT remove the underlying mathematical capability if it is required internally.

Refactor the smallest possible way:

- move the helper to internal/kernel as an unexported or kernel-only helper
- update relativity/relations.go to use that helper
- preserve exact dimension semantics
- preserve all existing tests
- do not create a new public compatibility wrapper
- do not create another public dimension constructor

Afterward verify that the public core API exposes exactly the nine specified dimension constructors:

- Dimensionless
- DimensionMass
- DimensionLength
- DimensionTime
- DimensionVelocity
- DimensionAcceleration
- DimensionForce
- DimensionMomentum
- DimensionEnergy

No extra public dimension constructor is permitted.

==================================================
4. MINOR FIX — MECHANICS MANIFEST TEST MESSAGE
==================================================

Correct the misleading test message in:

mechanics/manifest_test.go

Current logic checks:

    len(m.Domain) != 9

Therefore the error must say:

    expected 9 domains, got %d

Do not change the actual test condition.

==================================================
5. MINOR FIX — CoreObject() API TYPE
==================================================

Align domain wrapper signatures with the normative public API.

Use:

    func (m Mass) CoreObject() core.Object

rather than exposing kernel.Object directly in the method signature.

Apply to all affected mechanics/relativity domain wrappers.

Because core.Object is a type alias, preserve identical runtime semantics.

Do not expose internal kernel types unnecessarily through public domain signatures.

==================================================
6. DO NOT TOUCH SUCCESSFUL INVARIANTS
==================================================

Do not alter the already-passing behavior for:

- single mint authority
- object immutability
- canonicalization
- exact rationals
- OperationParams
- provenance
- assumption propagation
- candidate containment
- package isolation
- manifest integrity
- replay
- StepID/index
- no-panic boundary
- SelectBranch
- Solve
- Limit
- Differentiate
- E=mc² derivation
- Lorentz limit

Do not “fix” things that are already correct.

==================================================
7. E=mc² MUST REMAIN IDENTICAL
==================================================

Verify the exact chain remains:

EnergyMomentumRelation
→ ZeroThreeMomentum
→ Substitute
→ Simplify
→ Solve(Energy)
→ Compare(Energy, ZeroEnergy, gte)
→ SelectBranch
→ m*c²

MassEnergyRelation() must remain absent from executable/non-comment
references in relativity/derivation_test.go.

Do not use this corrective pass to change the derivation.

==================================================
8. RE-RUN ALL FREEZE TESTS
==================================================

Run:

go test ./...
go vet ./...

Then independently verify:

- exact 39-file tree
- package count
- package graph
- public API surface
- kernel mint authority
- kernel physics boundary
- theory/package isolation
- assumption non-leakage
- reverse constructor allowlist
- provenance matrix
- manifest corruption battery
- candidate validation
- replay integrity
- OperationParams canonicality
- malformed-input battery
- E=mc² derivation
- derivative
- Lorentz limit

==================================================
9. MUTATION CAMPAIGN
==================================================

Independently rerun:

- M1–M18
- N1–N7

Use disposable copies.

Use:

    -count=1

Restore the clean implementation between mutants.

Do not report the previous claimed 25/25 result as independently verified unless
you actually rerun the campaign.

Also ensure the previously established mutants remain killed after the API/file
corrections.

==================================================
10. PUBLIC API AUDIT
==================================================

Perform a fresh exported-symbol audit.

Specifically search for:

- generic object factories
- alternate mint paths
- extra public dimension constructors
- public kernel types leaking unnecessarily
- alternate provenance constructors that can mint trusted objects
- undocumented public compatibility helpers

The public API must remain exactly within the frozen contract.

==================================================
11. FINAL REPORT
==================================================

Use:

A. Baseline

Actual pre-fix results.

B. Corrections

For each of the four findings:

- file-tree correction
- DimensionEnergySquared correction
- manifest error-message correction
- CoreObject() signature correction

State exact files changed.

C. Explicit non-changes

Confirm:

- specs_v2_3 unchanged
- Plan 10 unchanged
- no new physics
- no Einstein-1905 corpus
- no architecture redesign
- no new files
- 39-file tree restored
- no new dependencies
- E=mc² path unchanged

D. Verification

Report actual:

- go test ./...
- go vet ./...
- exact file count
- public API audit
- mutation results
- E=mc² result
- all freeze invariants

E. Final status

Use exactly one:

FREEZE-READY

or

NOT FREEZE-READY

Do not declare FREEZE-READY unless every concrete blocker is actually
corrected and the corresponding executable verification passes.
```

### My current status assessment

I would call the system **“architecturally freeze-ready but repository-contract not yet freeze-ready.”**

That distinction is important. The adversarial review did not uncover a conceptual failure in the kernel, provenance model, theory isolation, or E=mc² reasoning path. It found two **hard contractual violations at the repository/API boundary**. adv_review8

Once those are corrected and the mutation campaign is independently rerun, we should do **one final adversarial verification**, rather than immediately declaring the kernel frozen.
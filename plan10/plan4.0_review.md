I reviewed **Plan 4 as the remediation plan**. It is substantially better and addresses the real root causes rather than just patching the adversarial findings. The key architecture choices are now correct: breaking typed candidate metadata, layered validation, one shared candidate trust path, session-owned identification replay, and no public object decoder. :chatgpt-content-reference{index="0"}

I would, however, make **four corrections before implementation**.

## 1. Most important: make the canonical-byte replay substrate explicit

Plan 4 correctly identifies the root problem:

> “No canonical-bytes replay substrate: objects retained live, so decode + invariant-check + hash-verify stages are vacuous.” :chatgpt-content-reference{index="1"}

But Stage 5 still leaves an ambiguity. The current committed structure is described as:

```text
InputCanonicals []core.Object
OutputCanonical core.Object
```

and then Stage 5 says to decode canonical objects through `LoadObjectJSON`. :chatgpt-content-reference{index="2"} :chatgpt-content-reference{index="3"}

That is not enough. **The committed ledger must retain the canonical representations as the authoritative replay substrate**, not merely retain live `core.Object` values and regenerate canonical bytes later.

The implementation plan should explicitly say:

```text
Draft:
    live core.Object values are acceptable

Commit:
    canonicalize each input/output object once
    retain immutable canonical bytes in Step
    compute hashes from those retained canonical bytes

Validate:
    decode retained bytes
    validate kernel invariants
    verify hashes
    replay from decoded retained values
```

Otherwise an implementation could still silently regenerate canonical JSON from a live object and technically “use the codec” without proving that the persisted replay material is what was actually validated.

This is directly grounded in §16.13: `InputCanonicals[]`, `OutputCanonical`, and `ParamsCanonical` are explicitly the replay substrate. :chatgpt-content-reference{index="4"}

**I would add this as a mandatory sentence to Stage 5.**

---

## 2. Keep the 17-stage ordering visible; shared helpers must be leaf primitives

The layered-validation choice is correct:

```text
Ledger.Validate()
    structural/integrity

Session.Validate()
    full 17-stage replay
```

That matches v2.3, which assigns the full ordered pipeline to `Session.Validate`, while `Ledger` has its own public `Validate()` surface. 

But Plan 4 says:

> “Shared private primitives (hash-verify, chain-verify, StepID/index, params-decode): implement once, call from both.” :chatgpt-content-reference{index="6"}

That is potentially dangerous.

A helper such as:

```text
validateStep()
```

that internally checks five things could accidentally destroy the required stage ordering.

Instead specify:

> **Shared helpers MUST be atomic/leaf-level primitives. `Session.Validate` is responsible for orchestrating them in the exact §16.19 order. No shared helper may silently perform later validation stages.**

So:

```text
Session.Validate
  stage 3 → checkStepID(index)
  stage 4 → decodeInput()
  stage 5 → verifyInputHash()
  ...
```

rather than:

```text
validateStep()
  → does stages 3,4,5,8,9...
```

This will make the mutation tests and `TestValidatePipelineOrder` much more trustworthy.

---

## 3. `MintObjectMust` should probably be eliminated, not moved to a new test file

Plan 4 says:

> “Remove `MintObjectMust` from prod; move helper to `_test.go` or replace callers...” :chatgpt-content-reference{index="7"}

The **remove-from-production** part is correct.

But the exact 39-file tree has only:

```text
internal/kernel/types.go
internal/kernel/mint.go
```

There is no `internal/kernel/*_test.go`. :chatgpt-content-reference{index="8"}

Adding one would violate the exact file tree unless the helper is placed in an already-existing test file.

So the cleanest instruction is:

> **Prefer replacing production callers with explicit `(Object, error)` handling. Do not create a new test file merely to preserve `MintObjectMust`. A test-only helper may live in an existing test file if genuinely useful.**

That keeps the 39-file invariant airtight.

---

## 4. Make candidate-containment recomputation algorithmically explicit

Stage 6 says:

> “`ResearchCandidate.Validate`: add derivation-output scan — recompute per-step provenance...” :chatgpt-content-reference{index="9"}

This is directionally right, but still too vague for an implementation agent.

The important question is:

> **How exactly is expected provenance recomputed?**

The stored step has inputs, operation, output, provenance status, etc. For a transformation, the expected provenance comes from the operation's provenance propagation law. For Identify, the expected status depends on the input statuses. For assertions, the retained object itself is the source of truth.

The plan should explicitly require:

```text
assertion:
    retained object's provenance is authoritative for consistency

transformation:
    replay operation from retained canonical inputs/params
    derive expected provenance under ops provenance law
    compare expected vs retained output provenance

identification:
    reconstruct through session-owned helper
    compare full canonical output/hash/provenance
```

Then MRC-008 becomes a **replay consequence**, not a separate ad-hoc scan.

That would eliminate the risk of implementing a superficial:

```text
if candidate.Hypothesis.Status == HYPOTHESIS ...
```

check while missing contaminated intermediate outputs.

---

# One thing Plan 4 gets exactly right

The `UnverifiedResearchCandidate` solution is correct.

The plan explicitly keeps:

```text
ResearchCandidate
    unexported fields

UnverifiedResearchCandidate
    unexported fields

external JSON
    ↓
typed intermediate representation
    ↓
canonical/validation checks
    ↓
shared trusted construction path
```

rather than exporting candidate fields simply to satisfy `encoding/json`. :chatgpt-content-reference{index="10"}

That is the correct authority model and matches v2.3's explicit requirement that `UnverifiedResearchCandidate.Validate()` delegate to the same private trusted path as `Session.Seal()`. :chatgpt-content-reference{index="11"}

---

# One minor typo

Stage 7 says:

> `Ox2`

for the two O tests. :chatgpt-content-reference{index="12"}

Change that to:

```text
O:
TestMechanicsManifestCrossCheck
TestRelativityManifestCrossCheck
```

Purely editorial.

---

# My verdict on Plan 4

**Plan 4 is now the correct remediation plan.** I would not redesign it.

Make these final constraints explicit:

```text
1. committed Step retains canonical BYTES, not live Objects, as replay authority
2. Session.Validate owns the exact 17-stage ordering
3. shared validation helpers are leaf primitives only
4. do not create a new file for MintObjectMust; preferably eliminate it
5. define the exact provenance recomputation used by candidate containment
```

The rest of the plan is well aligned with the specification, including the exact session surface, 17-stage validation, typed candidate model, relativity assumptions, and A–S rehabilitation. :chatgpt-content-reference{index="13"}

**After those five clarifications, I would give Plan 4 to the implementation agent.**

Yes. Use the following as the **final implementation-agent prompt for the kernel freeze**. It deliberately excludes any Einstein-historical corpus expansion or new physics additions.

The primary MVP target remains the machine-composed modern special-relativistic derivation of \(E=mc^2\), using the required operation chain and explicitly **not** using `MassEnergyRelation()` as the derivation source. :chatgpt-content-reference{index="0"}

---

# IMPLEMENTATION PROMPT — FINAL KERNEL FREEZE / HARDENING PASS

You are the implementation agent for the Go Physics Compiler MVP.

Your task is to perform the **final narrow hardening and freeze pass** on the existing implementation.

This is **not a redesign** and **not a new feature phase**.

The objective is to close the remaining robustness and integrity gaps identified by the consolidated adversarial review, while preserving the already-approved architecture, specification, corpus, package graph, and 39-file implementation shape.

## 1. NON-NEGOTIABLE FREEZE CONSTRAINTS

Treat the following as frozen unless a change is strictly necessary to implement one of the explicit hardening gates below:

- `specs_v2_3.md`
- Plan 10 architecture
- 39-file implementation tree
- package boundaries
- dependency direction
- `internal/kernel` authority model
- `core` facade/type aliases
- `ops` pure symbolic transformation model
- `session` ledger/replay model
- existing mechanics corpus
- existing special-relativity corpus
- existing E=mc² derivation path
- no parser
- no CAS
- no numerical engine
- no general theorem prover
- no general solver
- no automatic hypothesis promotion
- no truth adjudication
- stdlib-only implementation

Do **not** broaden the scope.

Do **not** add new physics domains.

Do **not** add General Relativity, Quantum Mechanics, Electromagnetism, thermodynamics, cosmology, or any other new corpus.

Do **not** add a corpus reconstruction of Einstein's 1905 mass-energy argument.

That historical derivation is explicitly deferred until **after kernel freeze**.

The purpose of this task is to freeze the formal substrate that will later support broader physics corpora.

---

# 2. PRIMARY MVP ACCEPTANCE TARGET

Preserve and protect the existing primary architecture test:

```text
EnergyMomentumRelation()
    ↓
ZeroThreeMomentum()
    ↓
Substitute(...)
    ↓
Simplify(...)
    ↓
Solve(..., Energy())
    ↓
Compare(Energy(), ZeroEnergy(), gte)
    ↓
SelectBranch(...)
    ↓
m*c²
```

The derivation must continue to use actual symbolic operations.

It must **not** call:

```go
MassEnergyRelation()
```

to obtain the answer.

The terminal expression must arise from the operation pipeline, with assumptions and branch selection preserved.

This is a demonstration of **AI-composed formal reasoning over encoded physics knowledge**, not a claim that the software independently discovered or empirically established the equation. :chatgpt-content-reference{index="1"}

Do not weaken this test while hardening the system.

---

# 3. FIRST: REPRODUCE CURRENT BASELINE

Before modifying code:

1. Inspect the entire existing repository.
2. Verify the current implementation state.
3. Run:

```bash
go test ./...
go vet ./...
```

4. Confirm the existing derivation tests and integrity tests.
5. Do not assume prior audit claims are still true; verify them.

Establish a baseline before making changes.

---

# 4. FREEZE GATE A — NO-PANIC MALFORMED-INPUT BOUNDARY

Harden operation dispatch and operation implementations so malformed but structurally decoded inputs return errors rather than panic.

### Required checks

Before positional indexing:

- verify exact/minimum arity as appropriate
- verify required node positions exist
- verify node kinds before type assertions
- verify operands before accessing child indexes
- verify all required fields before dereferencing

Pay particular attention to:

```text
ops.Apply
Session.Step
Differentiate
Limit
Solve
SelectBranch
```

### Limit

Require the limit variable to be exactly one valid `Symbol` where the contract requires a symbol.

Reject malformed/non-symbol/multi-node variable representations.

Do not silently reinterpret an arbitrary expression as the limit variable.

### SelectBranch

Do **not** blindly select branch index 0.

Validate the exact structural contract defined by v2.3:

```text
BranchSet(
    target,
    [
        positiveBranch,
        Neg(positiveBranch)
    ]
)

constraint:

Relation(gte, target, zero)
```

Validate:

- target structure
- branch count
- positive branch shape
- negative branch shape
- negation relationship
- constraint operator
- constraint target
- constraint zero value
- compatibility between target and constraint

Only then select the mathematically valid branch under the MVP's bounded rules.

### Tests

Add negative tests for malformed expressions/operations covering the above cases.

Tests must assert:

- an error is returned
- no panic occurs
- no partially committed session state is produced

Do not turn this into a general expression validator.

---

# 5. FREEZE GATE B — PROVENANCE LAW

Make provenance behavior mechanically explicit and test it comprehensively.

For every pure symbolic operation:

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

enforce:

```text
normal trusted/non-hypothesis inputs
    → DERIVED

any HYPOTHESIS input
    → HYPOTHESIS
```

Pure operations must never silently manufacture:

```text
IDENTIFIED
POSTULATED
DEFINED
```

as a result provenance.

Preserve the distinction:

- `Identify` is session authority
- pure ops are derivation machinery
- hypothesis contamination is sticky
- no promotion mechanism exists

### Identify

Keep `Session.Identify` as the sole identification authority.

Test:

```text
pure op → never IDENTIFIED
Session.Identify → IDENTIFIED when permitted
HYPOTHESIS involved → remains HYPOTHESIS
```

Do not add semantic-quality requirements to justification text beyond the existing specification. Do not turn justification validation into an editorial/review system.

---

# 6. FREEZE GATE C — CORPUS AUTHORITY CONTENT LOCK

The existing corpus must remain a **closed trusted baseline**.

Mechanically enforce that trusted corpus metadata cannot drift silently.

For the established baseline, verify and lock:

```text
CorpusStatus = ESTABLISHED
```

for the appropriate trusted corpus items.

Constructor behavior must agree with the manifest.

Add/strengthen content validation so trusted constructors cannot silently diverge from manifest truth for relevant fields, including:

- canonical expression
- canonical expression hash
- kind
- dimension
- provenance
- corpus status
- assumptions
- source
- other already-required manifest metadata

### Reverse constructor allowlist

Add a reverse allowlist test:

```text
trusted constructor
    ↔
approved manifest item
```

Every trusted named physics constructor must correspond to an explicitly approved corpus item.

No extra trusted physics artifact should exist merely because a constructor happens to exist.

Do not redesign the constructor system.

---

# 7. FREEZE GATE D — KERNEL PHYSICS BOUNDARY

Preserve the strict separation:

```text
internal/kernel = generic semantic/mint machinery
domain packages = trusted physics corpus
```

Audit `internal/kernel` for physics leakage.

The only physics-specific identifiers permitted there are the identifiers explicitly mandated by the frozen specification, namely:

```text
KindMinkowski
LorentzFactorFunctionID = "lorentz_factor"
```

Do not add hidden laws, equations, Newtonian constants, relativity formulas, GR concepts, QM concepts, etc. to the kernel.

Add a mechanical regression test or static inspection test that rejects introduction of additional unapproved physics-specific identifiers into `internal/kernel`.

Do not make the test so broad that it rejects the two explicitly permitted specification-mandated identifiers.

---

# 8. FREEZE GATE E — TRUSTED CANDIDATE EMBEDDED-EXPRESSION VALIDITY

Harden candidate validation.

Any embedded structured expression in trusted candidate fields must be checked as a valid canonical expression before the candidate is accepted as internally valid.

Pay particular attention to:

```text
Prediction.Relation
FalsificationCondition.ContradictingCondition
RecoveryClaim.Condition
```

Validate their embedded `core.Expr` values structurally.

Reject invalid/malformed/non-canonical expressions.

Do not confuse this with scientific truth validation.

The candidate-validation success condition remains:

```text
internal consistency + schema integrity + replay integrity
```

not empirical truth.

Preserve the existing distinction between trusted `ResearchCandidate` data and external/unverified research candidates.

---

# 9. FREEZE GATE F — MANIFEST CORRUPTION BATTERY

Execute a deliberate corruption test battery against the manifest/constructor integrity layer.

Each corruption must be independently rejected.

At minimum exercise corruption of:

```text
canonical expression
canonical expression hash
dimension
kind
provenance status
corpus status
assumptions
source
```

where each field is applicable.

The point is to demonstrate that a trusted constructor cannot continue accepting a manifest item after its authoritative content has been altered inconsistently.

Keep these as negative tests; do not introduce runtime mutation into production corpus state.

---

# 10. FREEZE G — SERIALIZATION / CANONICALIZATION / OPERATION PINS

Freeze the exact operational serialization rules already established.

### Rational values

Preserve exact rational representation.

No floating-point serialization.

Canonical rational representation must remain exact and deterministic.

### OperationParams

There is a known inconsistency in the v2.3 textual example around unused `OperationParams` fields for `pow`.

For implementation purposes, freeze the operational behavior already adopted by the code:

```text
unused fields = empty/default values
```

Therefore a `pow` parameter object must not invent an unrelated operator such as `"eq"` merely because of the example inconsistency.

Add a regression test for the chosen operational encoding.

Do **not** rewrite the frozen v2.3 specification as part of this task.

### Determinism

Preserve deterministic ordering and canonical byte generation.

No map-iteration-dependent canonical output.

No non-deterministic serialization.

---

# 11. REPLAY / LEDGER SAFETY

Do not disturb the existing replay authority.

The ledger must continue to retain the canonical byte representations necessary for deterministic validation:

```text
InputCanonicals
OutputCanonical
ParamsCanonical
```

with their hashes.

Replay must remain byte-based rather than depending on regeneration from potentially changed live values.

Preserve the existing Session.Validate ordering and stage semantics.

Do not reorder validation stages merely to simplify implementation.

---

# 12. MUTATION-RESISTANCE REGRESSION

The previous audit established meaningful mutation coverage.

After changes:

- rerun the existing mutation campaign if available
- verify all previously killed meaningful mutants remain killed
- add new mutants corresponding to the newly hardened gates where practical

Especially protect against regressions that would reintroduce:

```text
hardcoded Limit result
hardcoded E=mc² result
bypassed candidate-containment validation
bypassed operation compatibility validation
non-canonical trusted object acceptance
alternate generic mint paths
```

Do not claim a mutation result that was not actually rerun.

---

# 13. DO NOT OVER-HARDEN

The following are **not blockers** and must not become new scope:

### Do not require sophisticated Identify justification semantics

A non-empty/trimmed justification requirement is sufficient for the MVP.

Do not build a prose-quality analyzer.

### Do not enforce “zero physics references whatsoever” in kernel

The specification explicitly permits its small set of mandated physics identifiers.

### Do not move `go:embed` merely because one reviewer preferred another file location

The requirement is functional: manifests must be embedded and loaded correctly.

### Do not attempt cryptographic authenticity

Integrity, canonicalization, hashing, and replay are in scope.

Authenticity/trust of external artifacts is a later concern.

### Do not add an AI planner

The kernel provides the formal substrate.

The agent that chooses the derivation sequence is outside this hardening pass.

---

# 14. PRESERVE THE CURRENT E=mc² ARCHITECTURE TEST

After all hardening changes, the complete E=mc² derivation must still succeed.

Required conceptual trace:

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

The final result must arise from symbolic transformation and explicit branch selection.

The test must continue to demonstrate that `MassEnergyRelation()` is **not** used to obtain the answer. :chatgpt-content-reference{index="2"}

Also preserve:

```text
Differentiate(KineticEnergy, v) → m*v
```

and:

```text
Limit(LorentzFactor, v, 0) → 1
```

without hardcoded outputs. :chatgpt-content-reference{index="3"}

---

# 15. FILE / ARCHITECTURE DISCIPLINE

Do not casually add files.

Do not rename packages.

Do not move functionality between packages unless required to correct a concrete violation of the frozen architecture.

Prefer:

```text
small production fix
+
focused regression test
```

over architectural refactoring.

If a proposed change cannot be justified as one of the seven freeze gates above, do not implement it.

---

# 16. REQUIRED FINAL VERIFICATION

At the end, run at minimum:

```bash
go test ./...
go vet ./...
```

Also perform:

```text
API/import boundary verification
39-file tree verification
package graph verification
kernel physics-boundary verification
manifest integrity verification
candidate containment verification
provenance verification
malformed-input/no-panic verification
E=mc² derivation verification
mutation regression verification
```

Verify no unintended generated files or dependency changes.

The repository must remain clean and reproducible.

---

# 17. FINAL REPORT FORMAT

Return a concise implementation report containing exactly:

### A. Baseline
What passed before changes.

### B. Changes implemented
For each of the seven freeze gates, state:

```text
PASS / NOT NEEDED
files changed
what was hardened
tests added/updated
```

### C. Explicit non-changes
Confirm:

```text
spec v2.3 unchanged
39-file architecture preserved
no new physics corpus
no Einstein historical derivation corpus
no parser
no CAS
no solver expansion
no automatic hypothesis promotion
no kernel redesign
```

### D. Verification
Report actual results of:

```text
go test ./...
go vet ./...
mutation tests
E=mc² derivation
manifest corruption battery
no-panic tests
provenance tests
kernel boundary test
candidate-expression validation
```

### E. Final status

Use only one of:

```text
FREEZE-READY
```

or

```text
NOT FREEZE-READY
```

If `NOT FREEZE-READY`, enumerate only concrete remaining blockers with file/test evidence.

Do not declare `FREEZE-READY` based on intent or code inspection alone when an executable verification can be performed.

---

## FINAL PRINCIPLE

This is the last **hardening pass before freezing the MVP kernel**.

Be conservative.

**Do not redesign what already works.  
Do not add capabilities merely because they are interesting.  
Do not expand the physics corpus.  
Do not solve future research problems now.**

The goal is to leave behind a small, deterministic, auditable formal substrate in which an AI agent can later be given physics knowledge and asked to **compose a derivation**, with the \(E=mc^2\) derivation serving as the primary proof-of-concept.
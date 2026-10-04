Use this as the **post-freeze adversarial code-review prompt**. The key is that the reviewer must treat the reported `FREEZE-READY` result as a **claim to attack, not evidence to trust**. The review should be independent, mutation-oriented, and explicitly distinguish specification compliance from robustness/security concerns.

# ADVERSARIAL CODE REVIEW PROMPT — POST PLAN 7/8 FREEZE

You are the **independent adversarial reviewer** of the Physics Compiler MVP after the implementation agent reported:

```text
FREEZE-READY
```

Do not assume that verdict is correct.

Your job is to attempt to **falsify the implementation's claimed freeze readiness**.

This is a review only.

## 1. Review posture

Be adversarial, skeptical, and evidence-driven.

Do NOT:

- trust the implementation agent's report
- accept passing tests as sufficient proof
- redesign the architecture
- propose unrelated features
- add new physics domains
- implement Einstein's 1905 derivation
- implement EBP 2.1
- require cryptographic authenticity
- demand protection against malicious compiler/linker-level attacks such as `unsafe`/`go:linkname`
- broaden the MVP into a general CAS/theorem prover/fuzzer

The question is:

> **Does the actual implementation satisfy the frozen specification and Plan 7/8 hardening requirements, including under adversarial modification?**

Use executable evidence wherever possible.

---

# 2. Authoritative materials

Treat the following as the governing references:

1. `specs_v2_3.md` — authoritative specification
2. Plan 10 architecture
3. Plan 7 / Plan 8 freeze deltas
4. the actual repository source
5. tests actually present in the repository

Resolve conflicts in this order:

```text
specs_v2_3.md
    >
Plan 10
    >
Plan 7/8 additive freeze gates
    >
implementation/report
```

Do not silently reinterpret the specification to accommodate the implementation.

---

# 3. Start by independently establishing the baseline

Inspect the actual repository.

Record:

```text
Go version
module path
file count
package count
package graph
source/test file split
dependencies
git status
```

Run:

```bash
go test ./...
go vet ./...
```

Do not rely on the implementation report's statement that these passed.

Confirm or falsify every reported baseline claim independently.

---

# 4. Verify the exact architectural boundary

Independently establish:

```text
internal/kernel
core
ops
session
mechanics
relativity
hypothesis
```

Verify:

```text
internal/kernel
    = trusted mint/invariant machinery

core
    = public trusted facade

ops
    = pure symbolic transformations

session
    = derivation/replay/ledger authority

mechanics
    = bounded classical mechanics corpus

relativity
    = bounded special-relativity corpus

hypothesis
    = provisional candidate space
```

Attack the dependency direction.

Look specifically for:

- reverse imports
- hidden registries
- reflection-based constructor discovery
- alternate minting paths
- package-level global state
- operation logic that consults session state
- theory knowledge accidentally embedded in generic machinery

Verify the 39-file production architecture claim independently.

Do not count planning documents outside the module tree as implementation files.

---

# 5. Attack the single-mint authority

The specification requires the internal kernel mint boundary to remain controlled.

Attempt to find **any alternate production route** that creates a trusted valid object without going through the intended kernel mint authority.

Search for:

```text
NewObject
Object{
factory functions
struct literals
alternate Mint*
clone/copy constructors
deserialization constructors
reflection
unsafe construction
```

Distinguish legitimate:

```text
internal/kernel test fixtures
```

from production escape paths.

Verify the public API does not expose a generic trusted-object constructor.

Mutation target:

```text
add a second valid-object construction path
```

The existing mutation report claims the second mint-path class was killed. Reproduce that claim independently.

---

# 6. Attack object immutability

Verify:

- authoritative `core.Object` fields are unexported
- zero-value object is invalid
- accessors cannot mutate state
- returned slices are defensive copies
- assumption/convention/provenance metadata cannot be mutated through aliases
- `*big.Rat` inputs do not remain caller-owned mutable state
- no hidden mutable references survive object construction

Try malicious caller-side mutation wherever the public API exposes slices, pointers, or nested values.

Add temporary negative tests in a disposable review copy if needed.

---

# 7. Attack OperationParams canonicality

Independently inspect the complete `OperationParams` schema.

The frozen rule is **symmetric per-kind unused-field rejection**.

Verify all four kinds:

```text
empty
pow
compare
identify
```

Expected semantics:

```text
empty
    Exponent=""
    Operator=""
    Justification=""

pow
    Exponent required
    Operator=""
    Justification=""

compare
    Operator required
    Exponent=""
    Justification=""

identify
    Justification required
    Exponent=""
    Operator=""
```

Attempt mutations:

```text
pow.operator = "eq"
pow.justification = "x"

compare.exponent = "2/1"
compare.justification = "x"

identify.exponent = "2/1"
identify.operator = eq

empty.exponent != ""
empty.operator != ""
empty.justification != ""
```

Verify canonical fixtures use the empty/default representation.

Check that normalization does not silently accept multiple semantically equivalent byte encodings.

---

# 8. Attack canonicalization

Review the complete expression canonicalization rules.

Attempt to find paths producing multiple canonical encodings for equivalent structures.

Test:

- Add flattening
- Mul flattening
- exact rational reduction
- sign normalization
- deterministic child sorting
- duplicate handling
- zero/one simplifications where specified
- canonical JSON field ordering
- deterministic assumption ordering
- deterministic convention ordering
- deterministic manifest encoding

Look for:

```text
map iteration
pointer addresses
Go runtime-dependent ordering
unstable sorting
locale-sensitive behavior
float serialization
non-canonical JSON numbers
```

Perform repeated canonicalization runs and compare exact bytes and hashes.

---

# 9. Attack exact rational handling

Verify all rational serialization is exact.

Test:

```text
1/2
-3/4
2 → 2/1
0 → 0/1
```

Attempt:

```text
1.5
0.5
1e-3
JSON numeric literals
```

where exact string representation is required.

Look specifically for accidental use of:

```text
float64
json.Number
direct *big.Rat marshaling
```

that could introduce non-canonical representations.

---

# 10. Attack operation arity and malformed AST boundaries

Do not merely replay the existing malformed fixtures.

Construct **new malformed structures** not represented by the current tests.

Attack:

```text
Apply
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
Session.Step
Session.Identify
```

Vary:

```text
zero inputs
one input
too many inputs
nil-ish/zero objects
wrong kinds
wrong expressions
nested malformed expressions
unexpected relation operators
unexpected branch structure
invalid rational exponents
invalid symbol structures
```

Every malformed case should produce a deterministic error rather than a panic.

For Session operations verify:

```text
failed operation
    →
no partial ledger mutation
no state transition
no hidden draft mutation
```

Do not demand a general-purpose fuzzing engine.

The question is whether the **bounded MVP contracts are actually closed against malformed inputs**.

---

# 11. Attack SelectBranch specifically

This is a high-value target.

The implementation must accept only:

```text
BranchSet(
    target,
    [
        positiveBranch,
        Neg(positiveBranch)
    ]
)
```

plus:

```text
Relation(gte, target, zero)
```

where `zero` has matching kind and dimension.

Try:

```text
three branches
two unrelated branches
Neg of wrong expression
reversed branches
wrong target
wrong operator
wrong zero
wrong kind
wrong dimension
constraint on another symbol
nested unrelated BranchSet
```

Verify that no malformed form can cause branch 0 to be selected merely because it happens to exist.

Verify `selected_branch/<hash>` is deterministic.

---

# 12. Attack Solve

The MVP solver must remain extremely narrow.

Accepted shape:

```text
Relation(eq, Pow(Symbol(target), 2), Expr)
```

Attempt to make it accept:

```text
linear equations
multiple variables
reversed equations
nested target expressions
non-symbol target
wrong powers
negative powers
cubic equations
multiple relations
```

These must remain unsupported.

The reviewer should specifically test whether implementation details accidentally broaden the solver beyond the specification.

---

# 13. Attack Limit

The Lorentz path is allowed only through its bounded fixed body.

Verify:

```text
Limit(LorentzFactor(), Velocity(), ZeroVelocity())
→ 1
```

is obtained by:

```text
expand fixed body
substitute
simplify
```

and not by:

```text
if functionID == "lorentz_factor" { return 1 }
```

Also verify a non-Lorentz limit still performs ordinary direct substitution + simplification rather than incorrectly returning `1`.

Attack:

```text
non-symbol variable
wrong variable
wrong value
unsupported function
unknown function ID
malformed Call
```

---

# 14. Attack hidden physics logic

Search `internal/kernel` and `ops` for:

```text
symbol == "m"
symbol == "c"
symbol == "p"
symbol == "E"
physics-specific switch branches
formula-specific constants
law-specific shortcuts
```

The only permitted special-relativity machinery is the explicitly specification-mandated closed-world mechanism.

For the Lorentz body:

```text
"v"
"c"
```

may appear only inside the specifically permitted fixed body-construction path.

Try to determine whether generic algebra has accidentally become:

```text
special-relativity algebra
```

through symbol-name dispatch.

Also ensure the kernel does not contain the actual equations of:

```text
Newton's law
p = mv
E = mc²
energy-momentum relation
```

as hidden executable shortcuts.

---

# 15. Attack theory/package isolation

The architectural principle is:

```text
kernel
    = formal machinery

mechanics
    = classical mechanics framework

relativity
    = special relativity framework

hypothesis
    = provisional candidates
```

Verify:

```text
internal/kernel
    imports no theory packages

mechanics
    cannot inject relativity assumptions into its constructor outputs

relativity
    cannot inject mechanics assumptions into its constructor outputs

hypothesis
    cannot silently become established corpus material
```

Do not exact-pin complete assumption sets.

Test forbidden cross-framework keys.

Critically distinguish:

```text
constructor output isolation
```

from:

```text
legitimate mixed-framework derivation
```

A mixed derivation may legitimately accumulate assumptions from multiple explicit input theories.

---

# 16. Attack assumption propagation

This is a major red-team target.

The required law is:

```text
output assumptions
    =
exact inherited input assumptions
    +
explicitly enumerated operation-generated assumptions
```

Verify inherited assumptions preserve:

```text
Kind
Key
Value
```

not merely the key name.

Attempt to mutate:

```text
same key
different value
```

and verify:

```text
AssumptionConflictError
```

Attempt to inject an assumption not present in:

```text
inputs
or
the operation's explicit generated-key set
```

and verify rejection.

Verify legitimate generated assumptions such as:

```text
denominator/*
selected_branch/*
```

continue working.

---

# 17. Attack provenance

Independently verify the complete provenance law.

For every pure operation:

```text
clean input
    → DERIVED

any HYPOTHESIS input
    → HYPOTHESIS
```

Pure operations must never create:

```text
IDENTIFIED
POSTULATED
DEFINED
```

unless explicitly required by a session-specific assertion path outside pure ops.

Then attack:

```text
Session.Identify
```

Verify:

```text
clean operands → IDENTIFIED
hypothesis contamination → HYPOTHESIS
```

Try to create an `IDENTIFIED` artifact through:

```text
Simplify
Compare
Solve
SelectBranch
Multiply
Substitute
```

and ensure that is impossible.

---

# 18. Attack corpus authority

Verify the pinned constructor/manifest relationship.

The reverse rule must be:

```text
every exported object-producing constructor
    ∈
manifest constructor IDs
    ∪
pinned non-manifest allowlist
```

Pinned allowlist:

```text
NewKineticEnergy
Velocity
ZeroThreeMomentum
ZeroEnergy
ZeroVelocity
hypothesis concept
```

Do not infer permission from status.

Attempt to add a new constructor that returns:

```text
DEFINED
DERIVED
ESTABLISHED
```

and verify the reverse allowlist catches it.

Check that parameterized `NewKineticEnergy` is correctly handled as non-manifest.

---

# 19. Attack manifest integrity

Independently perform in-memory corruption.

Mutate one field at a time:

```text
canonical expression
canonical expression hash
dimension
kind
provenance
corpus status
assumptions
source
constructor
item set
```

Verify rejection.

Specifically test:

```text
ESTABLISHED → CONTESTED
```

must fail the **freeze-baseline** check but remain schema-legal.

Do not incorrectly make `CONTESTED` an invalid corpus-status enum.

Verify source pins:

```text
NewtonSecondLaw
    → Newton, Principia

other mechanics items
    → Classical Mechanics corpus

EnergyMomentumRelation
MassEnergyRelation
    → Einstein, 1905

other relativity items
    → Special Relativity corpus
```

Check both manifest and constructor metadata.

---

# 20. Attack candidate containment

Try to create a candidate where:

```text
HYPOTHESIS-dependent result
```

is presented as:

```text
ESTABLISHED
DEFINED
DERIVED
IDENTIFIED
```

Candidate validation must reject that.

Also attack embedded:

```text
Prediction.Relation
FalsificationCondition.ContradictingCondition
RecoveryClaim.Condition
```

with:

```text
invalid Expr
non-canonical Expr
malformed node
```

and verify validation/sealing fails.

Remember:

```text
candidate may reference established premises
```

That is allowed.

What must not happen is a **hypothesis-derived result becoming trusted status**.

---

# 21. Attack Session replay

Verify retained canonical bytes are actually the authority for replay:

```text
InputCanonicals
OutputCanonical
ParamsCanonical
```

and corresponding hashes.

Attempt:

```text
mutated live values
reordered inputs
altered params
changed output object
changed step index
changed StepID
changed previous hash
changed current hash
changed assumption hash
changed convention hash
```

and verify validation catches the alteration.

Verify `Session.Validate` follows the required stage order.

Do not allow a successful replay merely because regenerated live objects happen to look equivalent.

---

# 22. Attack StepID/index integrity

Verify:

```text
StepID
    ↔
Index
```

is recomputed/validated rather than trusted from serialized input.

Attempt:

```text
step-000002 with Index=1
step-000001 with Index=2
duplicate StepID
skipped index
reordered steps
```

and verify rejection.

---

# 23. Attack E=mc² specifically

This is the **primary MVP experiment**.

The reviewer must verify both:

### A. Actual derivation

The trace remains:

```text
EnergyMomentumRelation
→ ZeroThreeMomentum
→ Substitute
→ Simplify
→ Solve(Energy)
→ Compare(Energy, ZeroEnergy, gte)
→ SelectBranch
→ m*c²
```

### B. Anti-hardcoding firewall

`relativity/derivation_test.go` must contain zero executable/non-comment references to:

```text
MassEnergyRelation
```

Attempt to introduce a future call and verify the source-level guard fails.

Also verify the result's assumptions are traceable to:

```text
RestMass
SpeedOfLight
ZeroThreeMomentum / RestFrame
SelectBranch
```

and are not attached merely because the expected answer is known.

Do not evaluate this as “Einstein's original 1905 derivation.”

It is the modern special-relativistic derivation from encoded SR premises.

---

# 24. Attack derivative and Lorentz tests

For:

```text
Differentiate(KineticEnergy, v)
```

verify the result is actually produced by the bounded derivative engine:

```text
→ m*v
```

not a special-case kinetic-energy shortcut.

For:

```text
Limit(LorentzFactor, v, 0)
```

verify the body-expansion path is actually executed.

Mutation targets:

```text
hardcode m*v
hardcode 1
```

must be killed.

---

# 25. Mutation campaign

The implementation report claims:

```text
M1–M18: 18/18 killed
N1–N7: 7/7 killed
```

Do NOT accept this without independent verification.

Use disposable copies.

For each mutation:

```text
apply exactly one mutation
run relevant tests
record result
restore clean state
```

Use:

```bash
-count=1
```

At minimum independently reproduce mutations targeting:

```text
hardcoded E=mc²
hardcoded Limit
second mint path
candidate containment bypass
non-canonical acceptance
operation compatibility bypass
pow operator acceptance
extra ESTABLISHED constructor
kernel physics leakage
invalid candidate Expr
malformed Limit acceptance
malformed SelectBranch acceptance
ambient assumption injection
```

Report exact results.

---

# 26. Search for test theater

This is a critical adversarial pass.

Look for tests that are:

- testing helper code rather than production code
- asserting constants copied from the implementation
- constructing expected values through the same flawed pathway
- using hardcoded expected outputs without testing the intermediate derivation
- disabled/skipped
- build-tagged away
- conditional on environment inappropriately
- testing only the positive path
- never exercising the supposedly critical guard
- source-scanning the test itself rather than the actual production behavior
- mutating a copy incorrectly so the mutation never actually reaches the tested code

For every important security/integrity property ask:

> **Could the implementation be wrong while this test still passes?**

If yes, report the gap.

---

# 27. Search for specification drift

Compare actual implementation against every relevant frozen invariant.

Especially check:

```text
OperationParams
canonical expression node set
Kind enum
AssumptionKind enum
ProvenanceStatus
CorpusStatus
MRC-001..008
Session lifecycle
Session.Validate stage order
manifest schema
constructor cross-check
E=mc² chain
kernel physics allowlist
theory package isolation
```

Do not assume that because an audit report says “PASS,” the implementation still conforms.

---

# 28. Classify findings

Every finding must be classified:

```text
BLOCKER
HIGH
MEDIUM
LOW
INFORMATIONAL
NOT-A-FINDING
```

Use:

### BLOCKER

A direct frozen-spec violation, a trusted-boundary escape, incorrect derivation semantics, or a reproducible integrity failure that invalidates the claimed MVP freeze.

### HIGH

A meaningful flaw that could permit bypass of a core invariant but does not necessarily invalidate all MVP behavior.

### MEDIUM

A real robustness/conformance issue with bounded impact.

### LOW

Minor drift, documentation, test hygiene, or maintainability issue.

### INFORMATIONAL

Interesting observation without a required fix.

### NOT-A-FINDING

Reviewer's concern contradicted by the specification or current executable behavior.

Do not inflate severity.

---

# 29. No redesign rule

Do not recommend:

- new architecture
- new abstractions
- new packages
- new DSL
- new physics domains
- general theorem proving
- general CAS
- general fuzzing infrastructure
- cryptographic signatures
- hostile-process sandboxing
- EBP integration

unless the finding demonstrates that the **existing frozen architecture is internally inconsistent**.

The objective is:

```text
hardening the frozen kernel
```

not:

```text
designing v0.5
```

---

# 30. Final verdict

Conclude with exactly one:

```text
PASS — FREEZE HOLDS
```

or:

```text
FAIL — FREEZE DOES NOT HOLD
```

Do not use “mostly pass,” “conditional pass,” or similar wording.

If `FAIL`, identify the minimum concrete blockers that must be corrected before freeze.

If `PASS`, explicitly state why the major attack surfaces were independently defeated.

---

# 31. Required final report

Use exactly this structure:

```text
# Adversarial Review — Physics Compiler MVP

## 1. Executive Verdict

PASS — FREEZE HOLDS
or
FAIL — FREEZE DOES NOT HOLD

## 2. Independent Baseline

go version:
go test ./...:
go vet ./...:
tree/package verification:

## 3. Findings

### FINDING-001
Severity:
Category:
Location:
Attack:
Expected:
Observed:
Evidence:
Impact:
Disposition:

...

## 4. Core Invariant Results

Kernel mint authority:
Object immutability:
Canonicalization:
OperationParams:
No-panic boundary:
Provenance:
Assumption non-leakage:
Corpus authority:
Manifest integrity:
Candidate containment:
Replay integrity:
Package/theory isolation:
E=mc² anti-hardcoding:
Derivative:
Lorentz limit:

## 5. Mutation Results

M1–M18:
N1–N7:
Additional mutations:

## 6. Specification Conformance

Any deviations from:
- specs_v2_3
- Plan 10
- Plan 7
- Plan 8

## 7. Test-Quality Assessment

Identify any tests that could pass despite a broken implementation.

## 8. Required Corrections

Only concrete freeze blockers.

## 9. Final Verdict

PASS — FREEZE HOLDS
or
FAIL — FREEZE DOES NOT HOLD
```

## Final instruction to reviewer

Your strongest result is **not** “I found no bugs.”

Your strongest result is:

> **I attempted concrete ways to violate the frozen invariants, including independent mutations and malformed inputs, and can show exactly which attempts were rejected by which executable controls.**

Treat the implementation agent's `FREEZE-READY` report as the **first thing to attack**.
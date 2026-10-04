Yes — **they still need adversarial review**, and I would not declare these documentation files final yet.

I reviewed the uploaded `AGENTS.md`, `mechanics/README.md`, and `relativity/README.md` against the frozen v2.3 design material. Overall they are **substantially accurate**, but I found a couple of places where an adversarial reviewer should intervene before the 42-file re-freeze.

### 1. `AGENTS.md`: one real semantic problem

This passage is too strong:

> “Every object you touch belongs to exactly one layer.”

followed by:

> `ESTABLISHED FRAMEWORK` / `DERIVED ARTIFACT` / `HYPOTHESIS`

The repository actually treats **provenance and corpus status as separate dimensions**. For example, `MassEnergyRelation` is an `ESTABLISHED` corpus item while its provenance is `DERIVED`. So an object can simultaneously be a corpus item and a derived artifact. The document itself later acknowledges that corpus status is separately human-curated, which exposes the inconsistency. AGENTS

I would change that section to something like:

```text
Every object has a provenance status.
Corpus status is a separate, human-curated property and must never be
confused with provenance.
```

Then describe the three epistemic categories as **reasoning categories**, not mutually exclusive object states.

### 2. `mechanics/README.md`: anomaly identifier needs verification

The mechanics README says:

> `galilean_noninvariance`

as the anomaly identifier. README

The frozen Plan 10 metadata I checked uses the pinned mechanics anomaly identifier **`nonrelativistic_regime`**. plan10_v2_3

That is exactly the sort of tiny metadata discrepancy an AI agent could propagate into later reasoning. The README should use the **actual manifest identifier**, not a semantically reasonable invented synonym.

Everything else in the mechanics README is broadly aligned: classical/nonrelativistic scope, constructor surface, `KineticEnergy` treatment, derivation target, and package boundary. README README

### 3. `relativity/README.md`: looks strong

I don't see a comparable substantive problem there. Its special-relativity-only boundary, six assumptions, `-+++` convention, constructor distinction, E=mc² derivation chain, Lorentz limit, and explicit GR exclusion line up with the frozen design. README README

### 4. `AGENTS.md` is otherwise doing the right job

The manifest-first workflow, assumption discipline, dimension/kind checking, provenance containment, prohibition on using stored results as derivation premises, special-relativity boundary, and theorem-vs-primitive/kernel-growth distinction are all appropriate for the intended AI-agent role. AGENTS AGENTS AGENTS AGENTS

### Verdict

**Do the adversarial review.** Not because the documents are broadly wrong, but because they are becoming part of the **AI reasoning control surface**. A one-word manifest mismatch or an epistemic-model ambiguity is more important here than it would be in an ordinary README.

I would make the two corrections above first, then have the independent reviewer perform a **targeted doc-vs-code-vs-manifest adversarial pass**, followed by the already agreed full suite and 42-file re-freeze gate.


## prompt

You are the INDEPENDENT ADVERSARIAL REVIEWER for the Physics Compiler MVP documentation amendment.

Your review is independent of the implementation agent. Treat all implementation-agent claims, self-audit results, summaries, and reports as UNTRUSTED until independently verified.

GOAL
Determine whether the documentation amendment is accurate enough to accept and whether the repository is safe to re-freeze at exactly 42 files.

You are reviewing ONLY this documentation amendment and its directly affected conformance surface.

AUTHORIZED NEW REPOSITORY FILES — EXACTLY THREE
1. AGENTS.md
2. mechanics/README.md
3. relativity/README.md

IMPORTANT REPOSITORY RULE
The v2.4 amendment document itself is an external planning/audit artifact and MUST NOT be added to the repository tree.

FROZEN SPEC RULE
plan10/specs_v2_3.md remains byte-for-byte frozen.
Do not modify it.
Do not reinterpret it.
Do not silently reconcile it with general knowledge.
Use the frozen v2.3 specification as the normative baseline.

REFERENCE PRECEDENCE
When checking correctness, use this order:

1. specs_v2_3.md
2. actual production implementation
3. actual manifest.json files
4. existing tests
5. AGENTS.md / package README documentation

The documentation must describe the repository as it actually exists.
Documentation does not define semantics.

SCOPE
Review:
- AGENTS.md
- mechanics/README.md
- relativity/README.md
- TestRepositoryTreeExact or equivalent repository-tree allowlist
- any directly affected tests
- actual manifests
- actual exported API/package boundaries
- relevant production implementation
- relevant derivation tests

Do NOT:
- redesign the architecture
- add README files to other packages
- modify hypothesis/, core/, ops/, session/, or internal/kernel merely to improve documentation
- change physics semantics
- add new corpus items
- add new assumptions
- add a general parser/CAS/prover
- reopen the kernel freeze
- create a generic documentation allowlist
- introduce a loose “<=42 files” rule

The intended final repository invariant is:

39 frozen implementation files
+ AGENTS.md
+ mechanics/README.md
+ relativity/README.md
= EXACTLY 42 repository files

TestRepositoryTreeExact MUST admit exactly those 42 paths and reject every other repository file.

==================================================
A. AGENTS.md ADVERSARIAL REVIEW
==================================================

Verify every substantive claim against the actual code/spec.

Pay particular attention to:

1. EPISTEMIC MODEL

AGENTS.md currently describes three epistemic layers:

    ESTABLISHED FRAMEWORK
    DERIVED ARTIFACT
    HYPOTHESIS

It also says:

    “Every object you touch belongs to exactly one layer.”

Determine whether that sentence is semantically accurate.

Specifically check whether provenance status and corpus status are independent dimensions.

Verify the special case of MassEnergyRelation:
- corpus item
- corpus status ESTABLISHED
- provenance status DERIVED

If provenance and corpus status are independent, the “exactly one layer” wording is potentially misleading and must be reported as a finding.

Do NOT fix it by inventing a new ontology.
Recommend the smallest wording correction that accurately reflects the implementation/spec.

2. PROVENANCE

Verify:
- DEFINED
- POSTULATED
- DERIVED
- IDENTIFIED
- APPROXIMATED
- HYPOTHESIS
- pure-op propagation
- HYPOTHESIS contamination
- Identify-only creation of IDENTIFIED
- no promotion API

Check that AGENTS.md does not imply that formal derivability creates empirical truth or corpus authority.

3. CORPUS STATUS

Verify that AGENTS.md correctly describes corpus status as human-curated metadata and does not imply that operations compute or infer it.

4. REASONING WORKFLOW

Verify the checklist against the actual supported system:
- framework identification
- README-first
- manifest-first
- assumptions
- regime
- typed primitives
- relation selection
- dimension checks
- physical-kind checks
- provenance
- hypothesis contamination
- session recording
- branch/precondition handling
- formal-vs-physical distinction

Ensure no workflow step claims capabilities the MVP does not possess.

5. MATHEMATICAL CAPABILITY CLAIMS

Verify the exact 12 operations:

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

Verify:
- bounded Solve pattern only
- bounded Differentiate engine only
- direct-substitution Limit behavior
- Lorentz-factor fixed-body handling
- no general solver
- no integration
- no series
- no numerics
- no general theorem prover
- no CAS semantics

6. ASSUMPTION LANGUAGE

Verify all concrete assumptions named in AGENTS.md:
- RestMass m >= 0
- SpeedOfLight c > 0
- denominator nonzero preconditions
- framework assumptions
- convention handling
- bounded sign/nonzero entailment

Do not accept a claim of general theorem proving or unrestricted logical entailment.

7. MANIFEST CLAIMS

Verify every listed manifest field against the actual manifest schema:

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

Verify that AGENTS.md correctly says:
- README is orientation
- manifest is canonical machine-readable metadata
- statement is prose, not a machine equation
- canonical_expr is the machine-checkable expression
- constructor cross-check exists
- corpus status is preserved rather than derived

8. PACKAGE BOUNDARIES

Verify these statements against imports and actual code:

internal/kernel = generic trusted formal machinery
core = public facade
ops = pure symbolic transformations
session = derivation/replay/seal authority
mechanics = Classical Mechanics ONLY
relativity = Special Relativity ONLY
hypothesis = provisional candidate space

Verify:
- session does not import domain packages
- domains do not import each other
- kernel contains no accidental theory-specific physics
- relativity does not contain GR
- theorem material is not being presented as kernel primitives

9. KERNEL-GROWTH GUIDANCE

Verify that AGENTS.md says, in substance:
- established physical primitives belong in theory packages
- generic formal primitives may belong in kernel only when foundational/theory-neutral/reusable and not reasonably expressible outside kernel
- theorems/derived results belong in theory/user libraries
- hypothesis belongs in hypothesis/research space
- this document does NOT itself authorize kernel growth

Ensure AGENTS.md does not accidentally authorize adding General Relativity, equivariance theorems, Valentini H-Theorem, or other theory-specific theorems into internal/kernel.

10. E=mc² WORKFLOW

Verify the exact derivation chain against the implementation/test:

EnergyMomentumRelation
→ ZeroThreeMomentum
→ Substitute
→ Simplify
→ Solve(Energy)
→ Compare(Energy, ZeroEnergy, gte)
→ SelectBranch
→ m*c²

Verify that AGENTS.md:
- does not use MassEnergyRelation as the derivation premise
- does not hardcode the result
- does not describe this as Einstein’s historical 1905 reconstruction
- does not describe it as empirical discovery
- does not equate formal derivation with truth

==================================================
B. mechanics/README.md ADVERSARIAL REVIEW
==================================================

Verify every claim against:
- mechanics/manifest.json
- mechanics/primitives.go
- mechanics/relations.go
- mechanics tests
- frozen spec

Check:

1. Framework identity:
   classical_mechanics

2. Scope:
   point-particle Newtonian mechanics
   nonrelativistic regime

3. Not-in-scope:
   relativistic
   quantum
   gravitational
   continuum physics

4. Primitive constructor names and whether each is:
   - manifest item
   - non-manifest helper
   - parameterized object

5. Relations:
   F = m*a
   p = m*v
   K = 1/2*m*v^2

6. KineticEnergy handling:
   verify the README does not incorrectly imply that parameterized
   NewKineticEnergy is itself a manifest item.

7. Derivation target:
   Differentiate(KineticEnergy, Velocity) -> m*v

8. Assumptions/conventions:
   Verify the README accurately reflects the manifest and constructors.

9. Limits and anomaly IDENTIFIERS:
   This is a HIGH-PRIORITY CHECK.

Do NOT trust the README's anomaly identifier.
Read the ACTUAL mechanics/manifest.json and determine the exact canonical anomaly ID.

The current README says:

    galilean_noninvariance

The frozen Plan 10 metadata previously pinned:

    nonrelativistic_regime

Determine which identifier is actually present in the current repository manifest.

If the README differs from the manifest, report it as a documentation consistency defect.

Do not silently change the manifest.

10. Verify the claimed “9 domains” count against the actual manifest.

11. Verify the README does not claim that the classical framework is universally valid or that its limitation is an empirical truth computed by the library.

==================================================
C. relativity/README.md ADVERSARIAL REVIEW
==================================================

Verify every claim against:
- relativity/manifest.json
- relativity/primitives.go
- relativity/relations.go
- derivation_test.go
- frozen v2.3 spec

Check:

1. Framework ID:
   special_relativity

2. Explicit boundary:
   Special Relativity ONLY

3. Explicit exclusions:
   General Relativity
   gravitational dynamics
   curved spacetime
   black-hole physics
   quantum mechanics
   quantum gravity

4. Assumptions:
   rest_mass_nonnegative
   speed_of_light_positive
   minkowski_spacetime
   lorentz_symmetry
   no_gravitational_dynamics
   special_relativistic_regime

Verify exact keys/types/ownership against the manifest.

5. Convention:
   metric.signature = -+++

6. Constructor surface:
   Spacetime
   MinkowskiMetric
   RestMass
   Energy
   ThreeMomentum
   FourMomentum
   SpeedOfLight
   Velocity

Verify which are manifest items and which are helper/non-manifest constructors.

7. Relations:
   LorentzFactor
   EnergyMomentumRelation
   MassEnergyRelation

8. MassEnergyRelation:
   verify:
   - corpus item exists
   - provenance DERIVED
   - corpus status ESTABLISHED
   - derivable_from metadata
   - README does not confuse these properties

9. E=mc² derivation chain:
   compare README exactly with derivation_test.go

10. Lorentz limit:
   verify README correctly says fixed-body Call expansion rather than a name-based shortcut

11. E² dimension:
   verify the README accurately describes this as local generic dimension multiplication, not special kernel machinery.

12. Manifest counts:
   verify “10 items, 8 domains” against the actual manifest.

==================================================
D. DOCUMENTATION MUST NOT INVENT SEMANTICS
==================================================

Look aggressively for statements that are reasonable-sounding but unsupported.

Examples of prohibited documentation claims include:
- the library proves physical truth
- corpus status is inferred
- a theorem is a primitive
- a prose statement is machine-checked as an equation
- arbitrary symbolic solving exists
- arbitrary calculus exists
- arbitrary function definitions exist
- arbitrary logical entailment exists
- hypotheses can become established automatically
- package boundaries are merely stylistic
- kernel can safely absorb theory-specific concepts because they are “important”

Report every such claim found.

==================================================
E. DOCUMENTATION VS ACTUAL CODE
==================================================

For every material claim in the three documents, independently verify against code.

Use actual source where necessary.

Pay special attention to:
- exported API names
- constructor names
- package imports
- Kind names
- Dimension behavior
- operation names
- assumption keys
- provenance statuses
- corpus statuses
- manifest fields
- exact manifest item counts
- anomaly/limitation identifiers
- E=mc² derivation sequence
- Lorentz-limit behavior
- candidate/hypothesis handling

Do not merely compare prose with the implementation agent’s report.

==================================================
F. TREE / FREEZE REVIEW
==================================================

Independently inspect the repository tree.

Expected:

39 frozen implementation files
+ AGENTS.md
+ mechanics/README.md
+ relativity/README.md
= 42 exact repository files

The external v2.4 amendment document must NOT count toward the 42.

Verify TestRepositoryTreeExact:
- accepts exactly the 39 original files
- accepts exactly the 3 new documentation paths
- rejects every other file
- does not become a loose <=42 check
- does not add a generic "*.md" exception
- does not accidentally admit plan/audit/tooling artifacts

Verify that plan10/specs_v2_3.md remains unchanged.

==================================================
G. REGRESSION VERIFICATION
==================================================

Run independently:

go build ./...
go vet ./...
go test ./...

Also run the relevant documentation/tree tests.

Check for regressions in:
- MRC-001 through MRC-008
- canonicalization
- dimension checking
- physical-kind checking
- assumption conflicts
- convention conflicts
- provenance
- HYPOTHESIS containment
- manifest cross-checks
- session replay
- candidate validation
- E=mc² derivation
- Lorentz limit
- differentiation
- package isolation
- repository-tree exactness

Do not accept “the implementation agent reported green” as evidence.

==================================================
H. REVIEW METHOD
==================================================

For each finding, provide:

FINDING-ID
Severity: BLOCKER / HIGH / MEDIUM / LOW / INFORMATIONAL
File
Exact section or line
Observed claim
Authoritative evidence
Why it is wrong or potentially misleading
Minimal correction

Distinguish carefully between:
- factual defect
- documentation ambiguity
- wording improvement
- unsupported claim
- harmless stylistic preference

Do not inflate stylistic issues into blockers.

==================================================
I. REQUIRED SPECIAL CHECKS
==================================================

SPECIAL CHECK 1 — EPISTEMIC LAYER WORDING

Determine whether:

“Every object you touch belongs to exactly one layer”

is inconsistent with the independent dimensions of:
- provenance status
- corpus status

Use MassEnergyRelation as a concrete test case.

SPECIAL CHECK 2 — MECHANICS ANOMALY ID

Determine the exact anomaly ID in mechanics/manifest.json.

Compare it with the README.

Do not rely on prior reports.

SPECIAL CHECK 3 — RELATIVITY ASSUMPTION COUNT

The README says “Read the 6 framework assumptions”.

Verify that exactly 6 assumptions are actually intended by the current manifest/constructor model:
- 2 constraints
- 4 framework assumptions

SPECIAL CHECK 4 — DOCUMENTATION AUTHORITY

Verify that every README explicitly remains orientation-only and does not supersede manifest or code.

SPECIAL CHECK 5 — KERNEL GROWTH

Verify that the documentation does NOT imply:
- General Relativity belongs in relativity/
- GR belongs in kernel
- theorem results belong in kernel
- equivariance theorem belongs in kernel
- Valentini H-Theorem belongs in kernel

The intended rule is:
- established GR primitives -> future separate general_relativity package
- generic formal substrate primitives -> kernel only through future growth gate
- theorems/derived results -> theory/user libraries
- hypotheses -> hypothesis/research layer

==================================================
J. FINAL VERDICT
==================================================

At the end, produce exactly:

1. Executive verdict
2. Findings table
3. Required corrections
4. Verified invariants
5. Test results
6. Re-freeze authorization

The final line MUST be exactly one of:

DOCUMENTATION PASS — 42-FILE RE-FREEZE AUTHORIZED

or

DOCUMENTATION FAIL — 42-FILE RE-FREEZE NOT AUTHORIZED

IMPORTANT:
The independent reviewer alone controls the re-freeze authorization.

Do not modify files during this review unless explicitly instructed in a later remediation task.

If findings exist, do not “fix while reviewing.”
Report them precisely so a separate remediation step can address them.

The reviewer’s conclusion must be based on independently verified evidence, not the implementation agent’s self-audit.
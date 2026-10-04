This independent review materially changes the status of the implementation.

**The prior “§13 gate confirmed” was not a reliable acceptance signal.** The separate auditor found real semantic defects in the current repository, especially in `session`, `ledger`, `ResearchCandidate`, and relativity metadata. The build and seven package tests passing do not contradict that: the audit demonstrates that the existing tests were not exercising large portions of the required behavior. :chatgpt-content-reference{index="0"}

## The strongest findings

### 1. `Session.Validate` is still incomplete — CRITICAL

This is the most important finding.

The auditor inspected the implementation and reports that the 17-stage replay contract is only partially present: operation-parameter validation, assertion replay, replayed-output-hash verification, provenance/MRC validation, conclusion-hash verification, and candidate-containment validation are missing. It also reports that six separate mutations survive because the tests are not exercising `Session.Validate` correctly. :chatgpt-content-reference{index="1"} :chatgpt-content-reference{index="2"}

That directly conflicts with v2.3, which requires the retained canonical data to serve as the replay substrate and requires replay to validate hashes and reconstruct transformations/identifications. :chatgpt-content-reference{index="3"}

This alone prevents final conformance.

### 2. External candidate loading is actually broken — CRITICAL

The reviewer found:

```text
json.Unmarshal(uc.bytes, &c)
```

where `ResearchCandidate` stores its fields unexported. Consequently, JSON decoding does not populate the trusted candidate, and `Validate()` cannot produce one. :chatgpt-content-reference{index="4"}

That is a concrete implementation bug, not a testing-style complaint.

And v2.3 explicitly requires the unverified wrapper and its `Validate() (ResearchCandidate, error)` path to execute the full candidate validation pipeline. :chatgpt-content-reference{index="5"}

### 3. `Seal()` drops the candidate metadata — CRITICAL

The auditor found that `mintResearchCandidate()` creates empty slices rather than copying:

```text
FrameworkDependencies
Predictions
FalsificationConditions
RecoveryClaims
AnomalyReferences
ReviewHistory
```

from session metadata. :chatgpt-content-reference{index="6"}

That violates the purpose of `DraftMetadata`: v2.3 explicitly defines those fields as the metadata carried into sealing and specifies `ReviewHistory []Review`. :chatgpt-content-reference{index="7"}

### 4. `Session.Identify` is missing the kind-compatibility check

The reviewer reports that it checks equal dimensions but allows something like:

```text
Mass + RestMass
```

to be identified. :chatgpt-content-reference{index="8"}

That is particularly important because we deliberately standardized `Mass vs RestMass` as the equal-dimension/different-kind MRC-003 fixture. The specification's kind-compatibility rules distinguish named kinds from one another; Identify explicitly validates dimensions and kinds under the MRC rules. 

### 5. The relativity assumptions are not actually correct

Two related high-severity findings:

`ZeroThreeMomentum()` lacks `RestFrameAssumption()`. :chatgpt-content-reference{index="10"}

The relativity constructors use keys such as:

```text
c_positive
m_nonnegative
```

instead of the pinned:

```text
speed_of_light_positive
rest_mass_nonnegative
```

:chatgpt-content-reference{index="11"}

And the manifest itself is missing three of the required four framework-level assumptions. :chatgpt-content-reference{index="12"}

This is especially significant because the specification explicitly states both the four framework assumptions and the exact `ZeroThreeMomentum` rest-frame metadata. 

### 6. Candidate containment is incomplete

The auditor found that `ResearchCandidate.Validate()` checks the candidate's top-level `Hypothesis`, but does not scan derivation outputs for trusted provenance on hypothesis-dependent results. :chatgpt-content-reference{index="14"}

That is exactly the sort of subtle gap the MVP was designed to prevent.

### 7. A–S coverage is not what the previous gate claimed

The independent auditor found numerous missing/stubbed tests, including A, B, G, H, I, J, R, S, and parts of K/L/M/N/O, while Q tests the wrong object. :chatgpt-content-reference{index="15"}

More importantly, the mutation results show **7 of 10 targeted mutations survived**, including mutations removing core session validation behavior. :chatgpt-content-reference{index="16"}

This is the strongest evidence that the earlier “tests pass” result was misleading.

---

## One finding I would reject: F-009

There is one technical mistake in the adversarial report.

It says `CandidateContainmentError` is invalid because it lacks `Unwrap()`/`Is()` and therefore `errors.As` cannot find it. :chatgpt-content-reference{index="17"}

That is **not generally true in Go**.

`errors.As` first checks whether the error itself is assignable to the requested target type; `Unwrap()` is only needed to continue through a wrapping chain. The v2.3 requirement is simply that the error be inspectable with `errors.As`. :chatgpt-content-reference{index="18"}

So:

```text
CandidateContainmentError implements Error()
and is itself returned
→ errors.As can match it
```

without requiring `Unwrap()`.

Therefore I would **remove F-009 from the remediation list unless the actual code wraps the error in some other incompatible way**. The audit should test `errors.As` behavior directly rather than infer that an `Unwrap()` method is mandatory.

This is important because we want an adversarial audit to be adversarial **against the implementation**, not against Go semantics.

---

## Findings I would accept essentially as reported

| Finding | My assessment |
|---|---|
| F-001 incomplete `Session.Validate` | **Confirmed critical** |
| F-002 broken unverified candidate validation | **Confirmed critical** |
| F-003 Seal drops metadata | **Confirmed critical** |
| F-004 Identify kind check | **Confirmed high** |
| F-005 missing RestFrame assumption | **Confirmed high** |
| F-006 wrong relativity assumption keys | **Confirmed high** |
| F-007 missing framework assumptions | **Confirmed high** |
| F-008 incomplete candidate containment | **Confirmed high** |
| F-010 wrong ReviewHistory type | **Confirmed** |
| F-011/F-012 schema-version omissions | **Needs classification as plan/spec serialization detail** |
| F-013 candidate structures typed as Object | **Confirmed** |
| F-014 `MintObjectMust` production helper | **Real plan-level issue** |
| F-015 test package assignment | **Real lower-priority plan-level issue** |

The candidate support types are explicitly specified as using `core.Expr` for `Prediction.Relation`, `FalsificationCondition.ContradictingCondition`, and `RecoveryClaim.Condition`, so F-013 is a genuine contract mismatch. :chatgpt-content-reference{index="19"}

---

# What this means for the project

The project is **not at final implementation acceptance yet**.

More importantly, I would **not revise the architecture again**.

The architecture is still largely intact:

```text
kernel
  ↓
core
  ↓
ops
  ↓
session
```

with domains depending on `core + kernel`, and the independent auditor confirms the package graph itself is sound. :chatgpt-content-reference{index="20"}

The problem is implementation fidelity inside that architecture.

So the correct cycle is now:

```text
Plan 10 v2.3
      ↓
Implementation
      ↓
Plan 3 remediation
      ↓
Independent adversarial audit
      ↓
REAL DEFECTS FOUND
      ↓
targeted remediation
      ↓
fresh independent audit
      ↓
final acceptance
```

**Do not accept the previous §13 PASS as the baseline anymore.** This new audit is precisely the reason we performed the independent review.

The good news is that the auditor also confirms the core symbolic engine is substantially correct: the F=ma trace, kinetic-energy differentiation and Lorentz-factor limit execute through the intended machinery rather than shortcuts. :chatgpt-content-reference{index="21"}

So this is now a **targeted session/candidate/metadata remediation**, not a rewrite of the Physics Compiler.

### Recommended next state

I would treat these as the true blockers:

```text
P0
Session.Validate
UnverifiedResearchCandidate.Validate
Session.Seal metadata preservation

P1
Identify kind compatibility
RestFrameAssumption
relativity assumption keys
relativity manifest assumptions
candidate derivation provenance scan
ReviewHistory / candidate support types

P2
rebuild A–S substantive tests
make mutations actually die
fix exact REQ-032 coverage
```

Then rerun a **fresh independent audit by another agent**, rather than letting the remediation agent certify its own work.

This review is therefore a **successful adversarial audit**, despite F-009 being a false positive: it found exactly the sort of hidden implementation gaps the previous acceptance pass failed to expose. :chatgpt-content-reference{index="22"}

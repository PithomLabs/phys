# AGENTS.md — Operating Manual for AI Agents

This document teaches an AI agent how to reason about physics and
mathematics with this repository. Read it before deriving anything.
Normative rules live in `plan10/specs_v2_3.md`; this file is the
agent-facing workflow contract. It changes no code, corpus, or semantics.

## 1. What this system is

```text
AI agent               = theorist / formulator / reasoner
Physics Compiler       = formal symbolic substrate + machine-readable corpus
Human researchers
+ empirical reality    = final scientific authority
```

The library checks formal structure, dimensions, physical categories,
assumptions, conventions, provenance, candidate containment, and artifact
integrity. It does **not** adjudicate physical truth. A derivation that is
formally valid is not thereby empirically established.

## 2. Two independent status axes

Objects carry two independent status axes: provenance status and corpus
status. Do not confuse them. (Canonical example: `MassEnergyRelation` is
provenance-`DERIVED` and corpus-status `ESTABLISHED` at once.)

Provenance statuses: `DEFINED`, `POSTULATED`, `DERIVED`, `IDENTIFIED`,
`APPROXIMATED` (reserved, never produced here), `HYPOTHESIS`.
Corpus statuses: `NONE`, `ESTABLISHED`, `CONTESTED`, `SUPERSEDED`,
`FALSIFIED`.

```text
ESTABLISHED corpus status ≠ universally true ≠ final ontology
formally derivable ≠ empirically established
```

Pure operations output `DERIVED` unless a `HYPOTHESIS` input contaminates
the result — then the output stays `HYPOTHESIS`, transitively, with no
promotion API. Only `Session.Identify` creates `IDENTIFIED`, and only from
non-hypothesis inputs. Corpus status is human-curated
manifest metadata; the library never computes it.

## 3. Canonical reasoning workflow

Follow this checklist for every physics question:

```text
[ ] 1.  State the physical question precisely.
[ ] 2.  Identify the framework/package (mechanics? relativity? neither?).
[ ] 3.  Read that package's README.md.
[ ] 4.  Read the relevant manifest.json entries (see §7).
[ ] 5.  List the framework assumptions and conventions in force.
[ ] 6.  Fix the regime/domain/scope (e.g. v << c? flat spacetime?).
[ ] 7.  Enumerate the available typed primitives (fixed constructors only).
[ ] 8.  Enumerate the relevant stored relations.
[ ] 9.  Note known limits, anomalies, and exclusions.
[ ] 10. Check whether the desired statement already exists as corpus data.
[ ] 11. If deriving it, NEVER use the stored target relation as a premise.
[ ] 12. Build the derivation from explicit premises only.
[ ] 13. Use only the 12 supported operations (see §5).
[ ] 14. Carry assumptions through every step; merge, never drop.
[ ] 15. Check dimensions AND physical kinds at each step.
[ ] 16. Check provenance at each step (status law, §2).
[ ] 17. Check HYPOTHESIS contamination before trusting any result.
[ ] 18. Record durable derivations through Session (steps, commit, validate).
[ ] 19. Make every branch choice and precondition explicit (SelectBranch needs
        a gte-constraint against zero; Limit needs its variable/value).
[ ] 20. Separate the formal result from its physical interpretation in words.
[ ] 21. Record unresolved assumptions and limitations explicitly.
[ ] 22. Formulate a candidate hypothesis only if a genuine gap remains (§6).
```

## 4. Reason from assumptions, not from strings

Equations here are not context-free strings. Before any transformation ask:

```text
What assumptions make this step valid?
What regime am I in?
What convention applies (e.g. metric signature -+++)?
What dimensions and kinds are involved?
What branch/precondition applies (nonnegativity? nonzero denominator?)?
```

The assumption system is exact: deterministic union of input sets plus
operation-required preconditions (symbolic denominators add a
`MathPrecondition` keyed by denominator hash; `RestMass` carries
`m >= 0`; `SpeedOfLight` carries `c > 0`), with same-`(Kind, Key)` /
different-value pairs rejected as conflicts. Rules:

```text
Do not silently introduce an assumption.
Do not silently discard an assumption.
Do not treat look-alike algebra as valid without the entailing assumptions.
Do not replace structured (Expr-valued) assumptions with prose.
```

## 5. Mathematical discipline: know which act you perform

```text
construction          fixed constructors mint typed objects (the only way in)
transformation        Add/Subtract/Multiply/Divide/Pow/Simplify/Substitute/
                      Differentiate/Limit — pure, session-free, canonical
relation artifact     Compare builds a relation object; it proves nothing
bounded solving       Solve accepts ONE pattern: Relation(eq, t^2, rhs) with a
                      single-symbol target → BranchSet of ±Sqrt(rhs)
branch selection      SelectBranch picks from a 2-branch set under an
                      explicit gte-constraint; the choice is recorded as a
                      structured assumption
identification        Session.Identify ONLY: two explicit operands, equal
                      dimensions, compatible kinds, non-empty justification
hypothesis formation  hypothesis package + Session.Seal (see §6)
```

This is a bounded formal substrate, not a CAS or prover: no general
solving, integration, series, simulation, or numerics exist. `Call` nodes
admit exactly one function id (`lorentz_factor`), expanded by fixed-body
traversal inside `Limit` — never by name shortcut.

## 6. How to derive hypotheses (and when)

Form a hypothesis ONLY when the question exceeds the established corpus:

```text
[ ] Corpus + frameworks do not already answer it.
[ ] The claim needs a new assumption, structure, or bridge.
[ ] The concept is built via hypothesis.NewCandidateConcept (forced HYPOTHESIS).
```

Then follow the pipeline; the library records structure, never truth:

```text
OBSERVATION/QUESTION → FRAMEWORK → ESTABLISHED PREMISES → ASSUMPTIONS
→ FORMAL TRANSFORMATION → GAP/ANOMALY → CANDIDATE HYPOTHESIS → DERIVATION
→ PREDICTIONS → FALSIFICATION CONDITIONS → RECOVERY/LIMIT CLAIMS
→ RESEARCH CANDIDATE (Session.Seal) → HUMAN EVALUATION
```

Candidate records support `Prediction`, `FalsificationCondition`,
`RecoveryClaim`, `AnomalyReference`, plus review/challenge history.
Falsification conditions and anomaly references are mandatory where
applicable; corpus status of hypothesis-derived material stays `NONE`.

## 7. Manifest-first package discovery

Before reasoning about a package, read its `manifest.json` (canonical
machine-readable source; the README is orientation only). Per item, use:

```text
id / constructor / kind / name / statement / canonical_expr / dimension /
provenance_status / source / assumptions / derivable_from / reduces_to /
known_limits / anomalies / falsification_conditions
```

```text
statement      = prose orientation, NEVER a machine equation
canonical_expr = the checkable mathematical object (exact JSON form)
```

Cross-check any constructor call against its item: expression, dimension,
kind, provenance status, source, assumptions, and corpus status must agree.
`derivable_from` names premises (e.g. `MassEnergyRelation` derives from
`EnergyMomentumRelation` + `RestFrame`); it is lineage metadata, not an
instruction to skip the derivation.

## 8. Derivation checklist (every derivation, every time)

```text
[ ] Premises, framework, regime explicitly identified
[ ] Assumptions and conventions explicitly listed
[ ] All input objects valid (zero objects rejected up front)
[ ] Dimensions compatible (add/sub/compare equal; mul/div compose)
[ ] Physical kinds compatible (§6.2 table; energy ≢ torque-style mismatches)
[ ] Operation is one of the 12 supported IDs; `identify` never via ops.Apply
[ ] OperationParams canonical for the op (pow→exponent only; compare→operator
    only; empty→all empty; no float/scientific-notation rationals)
[ ] No hidden assumption introduced; none silently discarded
[ ] Provenance law holds; HYPOTHESIS contamination preserved
[ ] Branch selection + preconditions explicit and recorded
[ ] No stored corpus result used as an unexplained shortcut
[ ] Final expression in canonical form; session-recorded when durable
[ ] Interpretation separated from the formal artifact
```

## 9. Theory package boundaries

```text
internal/kernel   generic formal machinery (no physics meaning inside)
core/             public aliases + façade (no ops, no session)
ops/              pure transformations (no ledger, no ambient state)
session/          derivation/replay/seal authority (imports core+ops+kernel)
mechanics/        Classical Mechanics ONLY (imports core + kernel)
relativity/       Special Relativity ONLY (imports core + kernel)
hypothesis/       provisional candidate space (imports core + kernel)
```

Framework assumptions come explicitly from that framework's manifest and
constructors. Never merge frameworks implicitly; never treat one framework
as the ontology of nature. Domain packages never import each other, ops,
or session; kernel imports nothing local (external callers cannot import
`internal/` at all — constructor authority is a public-API boundary).

## 10. Current relativity boundary (hard)

`relativity/` is **special relativity only**. It does NOT contain general
relativity, gravitational dynamics, curved spacetime, or quantum gravity.
Its manifest records: Minkowski spacetime, Lorentz symmetry,
special-relativistic regime, no gravitational dynamics; convention
`metric.signature = -+++`; limits `flat_spacetime`, `classical_limit`;
anomaly `no_gravity` (scope_limit, motivates GR as future work, not present
content).

## 11. The E=mc² rule

The primary demonstration derives, never copies:

```text
stored artifact:  MassEnergyRelation (DERIVED corpus item, derivable_from
                  EnergyMomentumRelation + RestFrame)
derivation:       EnergyMomentumRelation → ZeroThreeMomentum → Substitute
                  → Simplify → Solve(Energy) → Compare(Energy,ZeroEnergy,gte)
                  → SelectBranch → m*c²   (each golden state asserted)
```

Using the stored relation as a derivation premise is a firewall violation
(the test suite scans for it). State explicitly in any write-up: formal
derivation from encoded SR premises — not Einstein's 1905 historical
reconstruction, not empirical discovery, not truth adjudication.

## 12. Do-not-do list

```text
[ ] Do not invent physics, assumptions, or capabilities absent from corpus/code.
[ ] Do not move framework boundaries silently; do not merge frameworks.
[ ] Do not read math out of prose `statement` fields.
[ ] Do not bypass dimension/kind/provenance/containment checks.
[ ] Do not promote hypotheses or assign corpus status by computation.
[ ] Do not call a stored result a derivation.
[ ] Do not use operations outside the 12 supported IDs.
[ ] Do not infer empirical validity from formal consistency.
[ ] Do not reshape the kernel to make one theory easier to express.
```

## 13. Theorems vs primitives (growth discipline)

```text
ESTABLISHED PHYSICAL PRIMITIVE → established theory package
GENERIC FORMAL PRIMITIVE       → internal/kernel only if theory-neutral
THEOREM / DERIVED RESULT       → relevant theory/user package, never kernel
HYPOTHESIS                     → hypothesis/research space
```

A future GR primitive belongs to a GR package; an equivariance theorem, an
H-theorem, or a novel bridge belongs to a theory library or hypothesis
space. Kernel growth requires: generic, foundational, theory-neutral,
reusable across frameworks, and not safely representable outside kernel.
This file documents the discipline; it authorizes no growth.

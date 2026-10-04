# Plan 7 — Final Kernel Freeze + Theory/Epistemic Package Isolation

**Base:** Plan 6 gates A–G adopted unchanged (no-panic battery; 12-op provenance
matrix; corpus lock + pinned reverse allowlist; kernel physics allowlist;
candidate expression validity; manifest corruption battery; serialization pins
incl. strict `pow` rejection; replay safety; M1–M18 + new-gate mutants; §17 A–E
report). This plan adds **Gate H only**, plus integration notes.

**Locked decisions:** Gate H tests live in `core/object_test.go` (import boundary
+ physics allowlist, reusing its AST walk infra); assumption non-leakage proven
at **both** levels (ops subset invariant + package forbidden-key checks, never
exact-key pinning); this plan recorded as a **new plan7.md file**.

**Verified status quo (no repair needed):** production imports are clean —
`internal/kernel` imports nothing local; `core→kernel`; `ops→core+kernel`;
`session→core+ops+kernel`; `mechanics/relativity/hypothesis→core+kernel` only;
no assumption registries (assumptions flow solely via `mergeInputs` over
explicit inputs). Gate H is mechanical proof, not repair.

## H0. Epistemic rule (wording locked)

Package isolation does **not** demote special relativity to HYPOTHESIS.
`CorpusStatus = ESTABLISHED` stays for both frameworks. The rule is:

> An established physical framework remains explicitly bounded by its
> assumptions, conventions, scope and empirical status; it does not become an
> ambient axiom of the generic kernel.

`internal/kernel` = theory-neutral formal machinery (plus the two
spec-mandated identifiers). Theory packages own their objects, relations,
assumptions, postulates, conventions, scope, limitations, anomalies,
derivations, and corpus metadata. `hypothesis/` stays provisional and never
auto-promotes.

## H1. Import boundary (`core/object_test.go`)

Extend the existing AST walk: assert `internal/kernel` non-test sources import
zero local packages (`mechanics`/`relativity`/`hypothesis`/`core`/`ops`/
`session`); assert dependency direction kernel←core←ops←session with domains
depending on the substrate only. No reverse dependencies.

## H2. Kernel physics allowlist (same file)

Non-test `internal/kernel` may contain only `KindMinkowski` and
`LorentzFactorFunctionID`/`"lorentz_factor"` plus their closed-world
validation. `ops` may contain only the fixed Lorentz body — reject
`symbol == "c"/"m"`-style branches. Reject any other theory-specific symbol,
law, constant, or framework assumption.

## H3. Ops-level assumption-subset invariant (ops tests)

For every pure operation: output assumption keys ⊆ (input assumption keys ∪
that operation's explicitly required/generated keys — e.g. denominator
preconditions, `selected_branch/*`). An operation cannot invent an unrelated
framework assumption. Enumerate per-op required keys in-implementation.

## H4. Package forbidden-key checks (domain + hypothesis tests, both directions)

`mechanics` output ∩ relativity-only keys = ∅, and reverse. Relativity-only
keys include `rest_frame`, `rest_mass_nonnegative`, `speed_of_light_positive`,
`minkowski_spacetime`, `lorentz_symmetry`, `no_gravitational_dynamics`,
`special_relativistic_regime`; mechanics-only keys enumerated from its manifest
in-implementation. Never assert exact-key sets (legitimate accumulation from
multiple sources must keep working). Hypothesis assumptions stay
hypothesis-associated; no ESTABLISHED material leaks into candidates.

## H5. E=mc² assumption ownership (`relativity/derivation_test.go`)

Every assumption on the derived `m*c²` result (`rest_mass_nonnegative`,
`speed_of_light_positive`, `selected_branch/*`, rest-frame-derived) must trace
to an explicit input object or operation step. Harden the existing
no-`MassEnergyRelation` convention into a mechanical source assertion. Keep
`Differentiate(KineticEnergy,v) → m*v` and Lorentz-`1` traces unhardcoded.

## Integration

- **E=mc² chain preserved:** EnergyMomentumRelation → ZeroThreeMomentum →
  Substitute → Simplify → Solve(Energy) → Compare(gte) → SelectBranch → m*c².
- **Reverse allowlist (locked):** pinned explicit set `{NewKineticEnergy,
  Velocity, ZeroThreeMomentum, ZeroEnergy, ZeroVelocity, hypothesis concept}`;
  every other ESTABLISHED-yielding constructor must map to a manifest item.
- **Serialization:** strict `pow` rejection (`operator`/`justification` must be
  empty); rational pins `1/2`, `-3/4`, `2/1`, `0/1` + JSON-number rejection;
  determinism double-run (extend R only on gap). Spec text untouched.
- **Replay safety:** canonical-byte retention, byte-based replay, and
  `Session.Validate` stage order/semantics untouched.
- **File discipline:** no new files (H1/H2 → `core/object_test.go`; H3 → ops
  tests; H4 → domain/hypothesis tests; H5 → `relativity/derivation_test.go`;
  Gate A battery → `ops/negative_test.go` + `session/session_test.go`;
  Gate F → both `manifest_test.go`); no renames; no cross-package moves without
  an architecture violation; production diffs only for genuine findings
  (Gate A real panics, Gate G `validateShape` tightening with existing-usage
  green-check first).
- **Mutation additions:** non-empty `pow.operator` acceptance, extra
  ESTABLISHED constructor, unapproved kernel physics identifier, invalid
  candidate expression, malformed Limit/SelectBranch acceptance, ambient
  cross-package assumption injection. Rerun full M1–M18 + new mutants in a
  disposable copy (`-count=1`, restore per mutant); claim only rerun results.
  Protect: hardcoded Limit/E=mc², containment/compatibility bypasses,
  non-canonical acceptance, second mint path, promotion, physics leakage.
- **Non-blockers:** justification prose quality; absolute zero-physics rule;
  `go:embed` location; authenticity crypto; AI planner; EBP runtime; 1905
  corpus.

## Verification & report

`go test ./...`, `go vet ./...`, plus: 39-file tree, package graph, import
boundaries, kernel allowlist, theory isolation, manifest integrity, reverse
allowlist, containment, candidate validity, 12-op matrix, no-panic battery,
serialization pins, byte replay, E=mc² trace, mutation results, clean tree,
no new deps. Report in §17 A–E format with verdict `FREEZE-READY` /
`NOT FREEZE-READY` (concrete file/test-evidenced blockers only).

*Status: PLAN ONLY — no implementation performed. Awaiting explicit instruction.*

## Plan 8 implementation record (freeze pass executed)

- Gates A–H implemented tests-only except: symmetric `validateShape`
  strictness (production) and H7 source-pin realignment (production data:
  mechanics 10 items → "Classical Mechanics corpus", relativity 8 items →
  "Special Relativity corpus", code + manifest.json, cross-checks green).
- New tests: `TestMalformedNoPanicBattery`, `TestSessionStepNoPartialCommit`,
  `TestProvenanceMatrix12`, `TestAssumptionSubsetInvariant`,
  `TestNoRelativityLeakage`, `TestNoCrossFrameworkLeakage`,
  `TestFinalAssumptionOwnership`, `TestDerivationUsesNoStoredResult`,
  `TestInvalidCandidateExpressions`, `TestManifestCorruptionBattery` (×2),
  `TestReverseConstructorAllowlist` (×2 + hypothesis status pin),
  `TestKernelImportIndependence`, `TestKernelPhysicsBoundary`,
  `TestLorentzBodyClusterOnly`, `TestValidateCandidateContainmentStage`,
  `TestValidateOperationParamsStage`, `TestLimitDirectSubstitutionNonUnity`,
  `TestObjectCanonicalRoundTripRejectsNonCanonical`, strengthened
  `TestNoGenericFactory`; extended cross-checks, `TestParamsCanonicalRoundTrip`,
  `TestNewtonSecondLawManifestMatch` pins; H0 README section.
- Mutation campaign: M1–M18 + N1–N7 (pow-operator, ESTABLISHED ctor, kernel
  identifier, candidate expr, Limit/SelectBranch shapes, ambient assumption),
  25/25 KILLED in disposable copies (`-count=1`, per-mutant restore).
- `specs_v2_3.md` / `plan10_v2_3.md` untouched. 39-file tree intact.

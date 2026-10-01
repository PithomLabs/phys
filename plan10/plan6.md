# Plan 6 — Final Kernel Freeze / Hardening Pass

**Scope:** narrow hardening + mechanical proof only. Not a redesign, not a feature phase.
No new physics domains, no GR/QM/EM, no Einstein-1905 corpus reconstruction (deferred
past freeze), no parser/CAS/numerical engine/prover/solver/promotion/truth scoring.

**Frozen unless a gate below strictly requires it:** `specs_v2_3.md`, Plan 10
architecture, 39-file tree, package boundaries, dependency direction,
`internal/kernel` authority, `core` aliases, `ops` purity, `session` replay model,
existing corpora, existing E=mc² path, stdlib-only.

**Locked decisions:** reverse allowlist = pinned explicit set
`{NewKineticEnergy, Velocity, ZeroThreeMomentum, ZeroEnergy, ZeroVelocity,
hypothesis concept}` (no inference rule); `pow` params with non-empty `operator`
are **strictly rejected** (canonical `{"kind":"pow","exponent":"2/1","operator":"",`
`"justification":""}`), spec text untouched (implementation pin).

## 0. Baseline
Inspect tree; `go test ./...` + `go vet ./...` green before changes. Planning-time
verification found guards, provenance law, candidate `.Valid()` checks, manifest
cross-checks, and rational pins largely present — this pass is proof + narrow
hardening, not repair.

## 1. Primary acceptance target (do not weaken)
`EnergyMomentumRelation → ZeroThreeMomentum → Substitute → Simplify →
Solve(Energy) → Compare(Energy, ZeroEnergy, gte) → SelectBranch → m*c²`,
via real operations, never calling `MassEnergyRelation()`. Harden the existing
no-call convention into a mechanical source assertion in
`relativity/derivation_test.go`. Keep `Differentiate(KineticEnergy,v) → m*v`
and `Limit(LorentzFactor,v,0) → 1` unhardcoded.

## 2. Gate A — No-panic boundary
Battery (ops `negative_test.go`, `session/session_test.go`) asserting typed
errors, never panics, and no partial session commit: `Apply` wrong-arity /
unknown-op / invalid-object inputs; non-symbol variables for
Substitute/Differentiate/Limit; non-relation/non-quadratic Solve shapes;
malformed BranchSets (count, negation, operator, target, zero); `Session.Step`
failure leaves the draft unextended. Production fix only if the battery exposes
a real panic path. No general expression validator.

## 3. Gate B — 12-op provenance matrix
Table test over all 12 pure ops (+ `Session.Identify` separately):
HYPOTHESIS input → HYPOTHESIS; clean inputs → DERIVED; pure ops never
IDENTIFIED/POSTULATED/DEFINED; Identify → IDENTIFIED when permitted,
HYPOTHESIS when contaminated. Built on central `resultStatus`; no justification
quality semantics beyond non-empty/trimmed.

## 4. Gate C — Corpus lock + reverse allowlist
Extend manifest cross-checks to per-item `CorpusStatus` and `Assumptions`;
assert both frameworks ESTABLISHED and derived/candidate artifacts never
ESTABLISHED. Reverse allowlist test per domain package: every
ESTABLISHED-yielding exported constructor must be in the manifest map; every
non-manifest constructor must be in the pinned set above.

## 5. Gate D — Kernel physics boundary
Allowlist source test (`core/object_test.go` AST infra): `internal/kernel`
non-test sources may contain only `{KindMinkowski,
LorentzFactorFunctionID/"lorentz_factor"}` + their closed-world validation;
`ops` only the fixed Lorentz body (reject `symbol == "c"/"m"` branches).

## 6. Gate E — Candidate expression validity
Negative test (`hypothesis/candidate_test.go`): invalid `core.Expr` in
`Prediction.Relation` / `FalsificationCondition.ContradictingCondition` /
`RecoveryClaim.Condition` fails validation/seal. Production checks exist —
verify, add proof only. Structural validity only, never truth.

## 7. Gate F — Manifest corruption battery
In-memory mutations in both `manifest_test.go` files, each independently
rejected via `ValidateManifestBytes`/cross-check: canonical-expr edit, hash,
dimension, kind, prov-status, corpus-status, assumptions, source, unknown
constructor, extra item. Never mutate production corpus state or disk files.

## 8. Gate G — Serialization pins
Strict `pow` operator/justification rejection in `validateShape` (check existing
usages stay green first) + regression test for the canonical form and the
`"eq"` rejection. Extend rational tests: `1/2`, `-3/4`, `2→2/1`, `0→0/1`,
JSON-number rejection, no direct `*big.Rat` marshaling. Determinism double-run:
verify R covers it, extend only on gap.

## 9. Replay safety
Keep canonical-byte retention (`InputCanonicals`/`OutputCanonical`/
`ParamsCanonical` + hashes), byte-based replay, and `Session.Validate` stage
order/semantics untouched.

## 10. Mutation regression
Rerun full M1–M18 in a disposable copy (`-count=1`, restore per mutant); add
cheap mutants for the new gates (pow-operator accept, extra ESTABLISHED
constructor, kernel physics identifier, invalid candidate expr, Limit/SelectBranch
malformed-shape acceptance). Protect especially: hardcoded Limit, hardcoded
E=mc², containment bypass, compatibility bypass, non-canonical acceptance,
second mint path. Claim only rerun results.

## 11. Non-blockers (do not expand)
Justification prose quality; zero-physics-absolute kernel rule; `go:embed`
location; cryptographic authenticity; AI planner.

## 12. File discipline
No new files, no renames, no cross-package moves without an architecture
violation. Small production fix + focused test over refactoring. Out-of-gate
changes are out of scope.

## 13. Verification & report
`go test ./...`, `go vet ./...`, import boundary, 39-file tree, package graph,
kernel boundary, manifest integrity, containment, provenance, E=mc² trace,
mutation results, no generated files/dependency changes. Report in §17 A–E
format with verdict `FREEZE-READY` / `NOT FREEZE-READY` (concrete
file/test-evidenced blockers only; never declare on inspection alone).

*Status: PLAN ONLY — no implementation performed. Awaiting explicit instruction.*

# Plan: Add Derivation Package Contract to phys-gr/solution

## 0. Objective

Bring `phys-gr/solution` into the same derivation-package discipline used by
the `phys/relativity` package, and establish a minimal repository-wide norm
so future implementation agents do not drift away from the required
derivation checklist.

This is metadata/test-contract work.

Do NOT redesign phys-gr.
Do NOT modify the physics derivation architecture.
Do NOT modify phys, internal/kernel, frozen specs, manifests outside
phys-gr/solution, or the GR implementation contract.

Existing (retain byte-identical unless noted):

- `solution/derivation_test.go` (RETAIN, unchanged in purpose)
- `solution/schwarzschild.go` (RETAIN, untouched)
- `solution/schwarzschild_test.go` (RETAIN, untouched)

Add:

- `solution/README.md`
- `solution/manifest.json`
- `solution/manifest_test.go`
- `phys-gr/AGENTS.md` (new file; see section 6 decision below)

## 1. Read and use the existing precedent

Inspect the existing `phys/relativity` package:

- `phys/AGENTS.md`
- `relativity/README.md`
- `relativity/manifest.json`
- `relativity/manifest_test.go`
- `relativity/derivation_test.go`

Use them as structural precedent only.

Do NOT copy Special Relativity semantics into GR.

The GR derivation remains:

```text
Schwarzschild vacuum derivation
from static, spherically symmetric A(r), B(r)
under MTW, -+++, coordinates [t,r,theta,phi], Lambda=0.
```

Additional implementation sources already surveyed:

- `phys/PROTOCOL.md` — the 4-artifact package shape
  (README / manifest.json / manifest_test.go / derivation_test.go).
- `phys-gr/solution/schwarzschild.go` — `Primary()` (14 stages from
  `UnknownFunction A(r), B(r)`), `Secondary()` (independent mu-form check).
- `phys-gr/solution/derivation_test.go` — firewall (no `mu` in PRIMARY,
  `k2` first at stage 12), golden states, M1-M8 mutations,
  PRIMARY/SECONDARY independence.
- `phys-gr/symbolic/expr.go` — 9-node GRExpr enum; `Normalize` +
  `CanonicalBytes()` produce the manifest bytes.
- `phys-gr/metric/metric.go` — ansatz `diag(-A(r), B(r), r^2, r^2 sin^2θ)`,
  source of premises wording.

Environment facts pinned during planning:

- `phys-gr` is NOT a git repo; `phys` is (pre-existing untracked
  `plan10/` / `PROTOCOL.md` files form the baseline — record before work).
- `phys-gr/solution/` currently holds exactly `derivation_test.go`,
  `schwarzschild.go`, `schwarzschild_test.go`.
- `phys-gr` has NO `AGENTS.md`; `phys/AGENTS.md` section 8 holds the
  derivation checklist.

## 2. Add solution/README.md

Write a concise package-level README (human orientation only).

It must explicitly contain:

```text
Purpose
Physical question (vacuum field outside a static spherical mass)
Framework = General Relativity
Regime = static, spherically symmetric vacuum spacetime
Coordinates = [t,r,theta,phi]
Metric signature = -+++
Curvature convention = MTW
Lambda = 0
```

Premises:

- static/spherical ansatz
- vacuum Einstein equations
- asymptotic-flatness boundary condition

PRIMARY:

- starts from unknown A(r), B(r)
- performs the prescribed 14-stage derivation
- supplies candidates and machine-certifies them

SECONDARY:

- starts independently from closed-form Schwarzschild
- recomputes curvature/vacuum verification
- never rescues PRIMARY

Explain explicitly:

```text
Primary = derivation
Secondary = independent verification
```

Include the derivation outline:

```text
metric
inverse
Christoffel
Riemann
Ricci
AB identity
vacuum implication
AB=k1
asymptotic boundary
B=1/A
R_theta-theta
rA equation
rA=r+k2
closed forms
final vacuum verification
```

Include limitations:

- bounded symbolic algebra
- no general ODE solver
- no general theorem prover
- no general series engine
- machine certifies supplied candidates
- formal consistency is not empirical truth

Include the distinction:

- k2=-2mu correspondence is recorded later in GR-8/session;
  it is not a PRIMARY derivation premise.

## 3. Add solution/manifest.json

Create a DERIVATION manifest, not a fake corpus manifest.

Top-level schema (exact keys):

```text
schema_version
artifact_type
derivation_id
framework_id
name
entry_points
regime
coordinates
conventions
premises
assumptions
expected_result
verification
limitations
anomalies
falsification_conditions
```

Pinned values:

```json
{
  "schema_version": "1",
  "artifact_type": "derivation",
  "derivation_id": "schwarzschild_vacuum",
  "framework_id": "general_relativity",
  "name": "Schwarzschild Vacuum Derivation",
  "entry_points": {"primary": "solution.Primary", "secondary": "solution.Secondary"},
  "regime": "static, spherically symmetric vacuum spacetime",
  "coordinates": ["t", "r", "theta", "phi"],
  "conventions": {"metric_signature": "-+++", "curvature_convention": "MTW", "cosmological_constant": "0"},
  "verification": {"primary_derivation_trace": "GRStep", "primary_stages": 14,
                   "secondary_verification": "independent closed-form vacuum check"},
  "anomalies": []
}
```

Details:

- `premises` must state the actual starting assumptions: static/spherical
  ansatz with unknown A(r), B(r); vacuum `R_uv = 0`; asymptotic flatness
  (A,B -> 1).
- `assumptions` grounded in code only (e.g. `r_positive` as used in
  `schwarzschild_test.go`).
- `expected_result`: `{ "A": <GRExpr JSON>, "B": <GRExpr JSON>,
  "description": "<prose only>" }` where A = `1+k2/r` normalized and
  B = `(1+k2/r)^-1` normalized, using the exact canonical GRExpr JSON
  representation. Do not use a prose equation as the machine equation.
- `limitations`: the six bounded-algebra items as `{id, description}`
  entries.
- `anomalies`: `[]` — DECIDED (Q&A): do not invent a scientific anomaly
  merely to fill the field; the GR solution package is a derivation
  artifact, not a claim of a discovered anomaly.
- `falsification_conditions`: structured `{id, description}` entries, e.g.
  any vacuum component certifying NONZERO; trace replay failure.

GRExpr byte-generation procedure (implementation step 0) — DECIDED (Q&A):

- Throwaway `go run` program in `/tmp` (never committed) built against the
  actual `phys-gr/symbolic` implementation: construct `1+k2/r` and
  `Pow(1+k2/r,-1)`, run `Normalize`, print `CanonicalBytes()`.
- Paste the exact output bytes into `manifest.json`.
- Do NOT hand-write canonical JSON.

## 4. Add solution/manifest_test.go

Use strict typed Go structs for every manifest level.

Do NOT use:

```text
map[string]any
reflection
dynamic constructor discovery
runtime registries
```

Use strict JSON decoding with unknown-field rejection
(`json.Decoder.DisallowUnknownFields()`, applies at all nesting levels).
Load via `go:embed manifest.json` for hermetic builds (relativity
precedent); resolve sibling artifact paths caller-relative, never CWD.

The test must enforce:

A. `manifest.json` exists.
B. `README.md` exists.
C. `derivation_test.go` exists.
D. required manifest fields exist.
E. exact values for: `artifact_type`, `derivation_id`, `framework_id`,
   `coordinates`, metric signature, curvature convention, Lambda, primary
   entry point, secondary entry point.
F. `expected_result` parses structurally as GRExpr
   (via `symbolic.ParseCanonicalJSON`).
G. expected A is exactly `1 + k2/r` (constructor-built with `symbolic`,
   compared with `symbolic.Equal` — mirrors the stage-13 assertion in
   `derivation_test.go`).
H. expected B is exactly `(1 + k2/r)^-1` (same method).
I. `entry_points.primary` MUST equal the exact literal
   `"solution.Primary"`.

   `entry_points.secondary` MUST equal the exact literal
   `"solution.Secondary"`.

   This is metadata validation only.
   Do NOT dynamically resolve or invoke entry points from manifest strings.
   No reflection or registry is permitted.
J. no unknown manifest fields are accepted (decoder rejects; test pins
   with a mutated in-memory copy).
K. no prose field is treated as executable mathematics (the `description`
   string is never passed to any GRExpr parser).
L. the derivation package contains all four required artifacts.

Anti-shortcut firewall (satisfies section 5; follows the
`TestDerivationUsesNoStoredResult` precedent): AST/substring scan
asserting `schwarzschild.go` and `derivation_test.go` contain zero
references to `manifest.json` / `expected_result` / manifest-loading
calls.

Do not require `manifest_test.go` to re-prove the derivation; that is the
job of `derivation_test.go`.

## 5. Manifest must not become a derivation shortcut

This is critical.

Neither `Primary()` nor `derivation_test.go` may use the
`expected_result` section of `manifest.json` as a derivation premise.

The manifest documents the expected derivation/result metadata.
It does not supply the derivation input.

The PRIMARY path must still begin from A(r), B(r) UnknownFunctions and
derive the result through the existing trace. Enforced by the firewall
scan in section 4 plus the unchanged `derivation_test.go` firewall
(no `mu` in PRIMARY, `k2` first at stage 12).

## 6. Derivation-package norm — phys-gr/AGENTS.md (DECIDED)

DECIDED (Q&A, option 1): create a NEW `phys-gr/AGENTS.md`. Do NOT edit
`phys/AGENTS.md` — the `phys` tree is frozen and its `AGENTS.md` is part
of that repository's agent-facing contract.

Hierarchy:

```text
phys/AGENTS.md
    -> governs frozen phys repository

phys-gr/AGENTS.md
    -> governs GR implementation and future GR derivation agents
    -> adds the derivation-package norm, scoped to phys-gr
```

Content: a pointer to frozen `phys/AGENTS.md` as upstream workflow
authority, plus the Derivation Package Artifact Contract (the exact
policy from the spec):

Every durable derivation package MUST contain:

```text
README.md
manifest.json
manifest_test.go
derivation_test.go
```

Roles are fixed:

```text
README.md
    human orientation only

manifest.json
    machine-readable derivation metadata

manifest_test.go
    package-contract and manifest validation

derivation_test.go
    executable derivation and anti-shortcut verification
```

The package MUST fail its tests if any required artifact is missing.

The manifest MUST NOT be used as a derivation premise or shortcut.

`manifest_test.go` MUST use strict typed decoding and MUST NOT use
reflection/dynamic registries.

A derivation is complete only when:

- README exists
- manifest exists and validates
- manifest test passes
- derivation test executes the real derivation
- replay/validation passes
- assumptions/conventions/provenance are checked
- no stored target result is used as a derivation premise

A derivation package MUST distinguish:

```text
formal derivation
independent verification
interpretation
empirical truth
```

The machine validates formal structure; human researchers and empirical
reality remain final scientific authority.

This subsection changes workflow only. It introduces no code semantics.

## 7. Do not generalize beyond this minimal protocol

Do NOT create:

```text
derivation framework
generic derivation interface
derivation registry
reflection-based discovery
repository-wide plugin mechanism
new kernel APIs
new GR AST nodes
new session abstractions
```

This is a repository convention enforced locally by each derivation
package's `manifest_test.go`.

## 8. Validation

Run:

```text
go test -count=1 ./...
go vet ./...
gofmt -l .
```

Also run specifically:

```text
go test ./solution -count=1 -v     (workdir: phys-gr)
```

Verify:

- `solution/README.md` exists
- `solution/manifest.json` exists
- `solution/manifest_test.go` exists
- `solution/derivation_test.go` exists (unchanged in purpose)

Verify the existing 14-stage derivation test remains unchanged in purpose
and still passes.

Verify `phys` tracked files remain untouched:

```text
git -C /home/chaschel/Documents/go/phys status --short   # == pre-work baseline
```

## 9. Final report

Report:

- files added
- AGENTS.md section added (new `phys-gr/AGENTS.md`)
- manifest fields validated
- derivation test status
- full test status
- phys integrity status

Final assertions:

```text
ARCHITECTURAL_CHANGES = 0
KERNEL_CHANGES = 0
DERIVATION_PROTOCOL_ADDED = 1
GR_SOLUTION_PACKAGE_CONFORMS = 1
```

## Appendix: clarification decisions (locked)

1. AGENTS.md target -> create `phys-gr/AGENTS.md`; do not touch `phys`.
2. `anomalies` field -> `[]`; real limits live in `limitations`.
3. GRExpr bytes -> throwaway `go run` against `phys-gr/symbolic`;
   `manifest_test.go` parses them back, `derivation_test.go`
   independently verifies how `Primary()` derived them.

# Physics Compiler MVP — Implementation Plan (Plan 10 — `plan10_v2_3.md`)

**Normative source:** `specs_v2_3.md` = *Physics Compiler MVP — Normative Implementation Specification v2.3* (§0–§41), together with the implementation-plan prompt v2.3 (`prompt.md`). References use `specs_v2_3.md §N` and `REQ-*` IDs. Module path everywhere: `github.com/PithomLabs/phys`.

`plan9.md` (= Plan 9.2 against specs v2.2 + prompt v2.2) is **historical/background only** — never normative. Where this plan and `plan9.md` disagree, `specs_v2_3.md` wins.

This document is a precision/conformance revision of `plan10.md`; the architecture is fixed, and the coding agent has no remaining architectural choices (specs_v2_3.md §41 handoff invariant).

---

## 1. Executive architecture

```text
internal/kernel  ← actual immutable Object/Expr/Kind/Dimension/metadata authority
       ↑
core             ← public API/type aliases + package-level façade functions
       ↑
ops              ← pure transformations (12 operations, no ledger, no ambient state)
       ↑
session          ← ledger/replay/seal authority (imports core, ops, internal/kernel)

mechanics   → core + internal/kernel
relativity  → core + internal/kernel
hypothesis  → core + internal/kernel
```

- **No `core → ops` dependency exists.** `core` declares no operation functions, no ledger, no session types; `core` imports only `internal/kernel` and the standard library. `core` MUST NOT import `ops` or `session` (REQ-004-02).
- `ops` MUST NOT import `session` (REQ-004-01, no cycles). `session` MAY import `core`, `ops`, `internal/kernel` (specs_v2_3.md §4).
- **`session` MUST NOT import `mechanics`, `relativity`, or `hypothesis`.** Domain packages supply objects to the session purely as `core.Object` values; this is asserted by `session_test.go::TestSessionImportIndependence`.
- **Authority boundaries:**
  - *Constructor authority (MRC-001):* Go `internal/` visibility + `kernel.MintObject` invariant validation + fixed domain constructors. External/public API boundary only — never claimed as hostile-source protection (REQ-000-04).
  - *Operation authority (MRC-002..005, 008):* pure `ops` functions, evaluated per call.
  - *Session authority (MRC-006, 007):* only `session.Session` may commit a ledger, produce `IDENTIFIED`, or mint a `ResearchCandidate`.
  - *Candidate containment (MRC-008):* enforced in `ops` (status law) and `session` (artifact validation).
  - *Review artifact (v2.3 reconciliation A):* `Challenge` and `Review` are defined **exclusively** in `core/corpus.go` (specs_v2_3.md §26.0, §27); `session` consumes `core.Challenge`/`core.Review` **without redefining them and without aliases**; no `review/` package exists.
  - *`physvet`:* deferred entirely (REQ-002-15; specs_v2_3.md §36 records the v0.5 contract only).
- Concrete backing types for `Object`, `Expr`, `Kind`, `Dimension`, `Assumption`, `AssumptionSet`, `Convention`, `ConventionSet`, `Provenance`, `CorpusStatus` live in `internal/kernel`; `core` exposes only `type X = kernel.X` aliases plus façade constructors/helpers — `core` declares no methods on aliased types (specs_v2_3.md §4).

## 2. Exact repository tree

```text
go.mod
README.md

internal/kernel/
    types.go
    mint.go

core/
    object.go        expr.go       dimension.go  assumption.go  convention.go
    provenance.go    corpus.go     canonical.go  errors.go
    object_test.go   expr_test.go

ops/
    arithmetic.go    simplify.go   transform.go  relation.go     dispatch.go
    operations_test.go   negative_test.go

session/
    session.go       ledger.go     research_candidate.go
    session_test.go

mechanics/
    primitives.go    relations.go  manifest.json
    manifest_test.go relations_test.go

relativity/
    primitives.go    relations.go  manifest.json
    manifest_test.go derivation_test.go

hypothesis/
    candidate.go     candidate_test.go

docs/
    paper-translation.md
```

Exactly **39 files** including `go.mod` (24 Go source, 10 Go test, 2 JSON manifests, 3 config/doc: `go.mod`, `README.md`, `docs/paper-translation.md`). No convenience, scaffold, generated, or placeholder files (REQ-003-01). Tests are adjacent `*_test.go` only (REQ-003-02). Package assignments (pinned):

| Test file | Go package | Why |
|---|---|---|
| `core/object_test.go` | `core_test` (external) | REQ-003-03: visibility/authority assertions (may import `internal/kernel` for fixtures — in-module) |
| `core/expr_test.go` | `core_test` | public-API round-trips |
| `ops/operations_test.go`, `ops/negative_test.go` | `ops` | white-box canonicalization detail |
| `session/session_test.go` | `session_test` (external) | session public surface + replay/tamper/seal |
| `mechanics/manifest_test.go` | `mechanics` (internal) | reads `go:embed` var without new export surface |
| `mechanics/relations_test.go` | `mechanics_test` | uses `ops` (test-only edge) |
| `relativity/manifest_test.go` | `relativity` (internal) | `go:embed` access |
| `relativity/derivation_test.go` | `relativity_test` | uses `ops` |
| `hypothesis/candidate_test.go` | `hypothesis_test` | uses `ops`, `session`, `relativity` |

`go.mod` pinned content: `module github.com/PithomLabs/phys` + `go 1.24`, with no third-party requirements (REQ-000-01/02). An empty `require` block is the authored form; Go tooling may normalize its textual presentation, and conformance is judged by the absence of third-party dependencies — not by block syntax.

## 3. Dependency-ordered implementation sequence

Each step: purpose / dependencies / files / result / gate. Gates accumulate; final gate is §13.

**Step 1 — Immutable kernel representation, canonicalization, construction primitives**
- Purpose: closed `Expr` node set, `Kind`, `Dimension`, metadata value types, structural canonicalization, canonical JSON, SHA-256, defensive-copy accessors, **and the kernel-side construction primitives required by every public metadata constructor (review F3)**.
- Deps: none (creates `go.mod`).
- Files: `go.mod`, `internal/kernel/types.go`.
- Result:
  - `kernel.Expr/Kind/Dimension/Assumption/Convention/Provenance/CorpusStatus/Object` (unexported `Object` fields, accessors, no mint yet); expr/dimension/metadata canonical encoders+decoders; child-sorting (three-key rule); `math/big.Rat` helpers.
  - **Kernel-side constructors (F3) — one authoritative validation location:** `NewRational`, `NewAdd/NewMul/Neg/NewPow/NewSqrt/NewCall/NewRelation/NewBranchSet` (expr), kernel-side `Dimension` construction for the nine `Dimension*` façade constructors, `NewTextAssumption/NewExprAssumption/NewAssumptionSet`, `NewConvention(s)`, and the single deterministic provenance constructor. All invariant validation for these types lives **here, once**; `core` façade constructors delegate and never duplicate validation logic.
  - **Deterministic set ordering (F4):** each `AssumptionSet`/`ConventionSet` element is canonicalized → canonical element bytes → sorted lexicographically by canonical bytes → stored in a deterministic slice; exact duplicates deduplicated. Go map iteration NEVER determines canonical order. This ordering participates in: canonical JSON → set hash → object hash → ledger hash chain → candidate determinism. Assumption conflicts remain keyed on `(Kind, Key)` with differing canonical values → `AssumptionConflictError`; convention conflicts on same `Key` + different `Value` → `ConventionConflictError`.
  - **Kernel helper ownership (review 🟡-6):** `EntailsNonNegative`, `EntailsPositive`, `EntailsNonZero`, `AssumptionSet.Merge`, `ConventionSet.Merge`, metadata equality, metadata canonicalization, metadata hashing — all implemented in `internal/kernel/types.go`. No helper packages/files.
  - **Defensive copies at construction:** `NewRational` copies its `*big.Rat` input; `NewPow` copies its exponent; all slice/`*big.Rat`-returning accessors return copies. No caller-owned mutable pointer ever becomes part of immutable kernel state.
- Gate: `go build ./internal/kernel`.

**Step 2 — Public core façade: expressions, dimensions, metadata, errors**
- Purpose: aliases, façade constructors (each delegating to the Step 1 kernel primitive), typed errors, expr/metadata equality+hash.
- Deps: 1.
- Files: `core/expr.go`, `core/dimension.go`, `core/assumption.go`, `core/convention.go`, `core/provenance.go`, `core/errors.go`, `core/canonical.go` (expr/metadata parts), `core/expr_test.go`.
- Result: `NewSymbol/NewRational/NewAdd/NewMul/NewNeg/NewPow/NewSqrt/NewCall/NewRelation/NewBranchSet`; nine `Dimension*` constructors; assumption/convention/provenance constructors (v2.3 §11.1 seven-kind set, §11.3, §13.0.1); 11 typed errors; `EqualExpr/HashExpr/HashAssumptionSet/HashConventionSet/CanonicalExprJSON/ParseExprJSON`. Kernel exports controlled canonical helper functions to in-module packages; `core` exposes the public canonical façade. (No package accesses another package's unexported Go identifiers.)
- Gate: `go vet ./core && go test ./core` (round-trip, sorting, defensive copies, fixed hash vectors).

**Step 3 — Mint authority & object façade**
- Purpose: single production mint entry point; object equality/hash; authority negative tests.
- Deps: 2.
- Files: `internal/kernel/mint.go`, `core/object.go`, `core/canonical.go` (object part), `core/object_test.go`.
- Result: `kernel.MintObject(ObjectSpec) (Object, error)` with the 7-point validation contract (specs_v2_3.md §5.2.1); `core.Object = kernel.Object`; `EqualObject/HashObject`; no `core` forwarder of `MintObject`. Object canonical decoding is **internal to `internal/kernel`** (no exported `DecodeObjectJSON`/`ParseObject` anywhere — review 🟡-8/prompt §5); the internal decoder validates the same authoritative object invariants as `MintObject` before any decoded object can enter replay, and guarantees `Encode(Decode(b)) == b`.
- Gate: `go test ./core` — zero-object invalid, exact accessor surface, unexported fields (reflect), mint rejects bad specs, AST export-surface scan, import-independence scan.

**Step 4 — Pure arithmetic operations (MRC-002..005, 008 at operation time)**
- Deps: 3. Files: `ops/arithmetic.go` (+ shared private helpers: validity check, MRC precondition pipeline, assumption/convention merge, provenance result law).
- Result: `Add, Subtract, Multiply, Divide, Pow` with exact signatures of specs_v2_3.md §15.1–15.5.
- Gate: `go build ./ops`.

**Step 5 — Simplify engine**
- Deps: 4. Files: `ops/simplify.go`.
- Result: bottom-up recursive simplification: §8.7 `Pow(x,0)→1` (nonzero-safe base gate; `0^0` unsupported), §9.6 identity/rational-`Pow` rules, §9.7 repeated powers, §9.8 `Sqrt` rules gated by bounded entailment, §9.9 recursive `Relation`/bounded `BranchSet` simplification, sign normal form, re-sort; `Simplify(x core.Object)`; finiteness predicate (§6 of this plan).
- Gate: `go build ./ops`.

**Step 6 — Transform operations**
- Deps: 5. Files: `ops/transform.go`.
- Result: `Substitute`, `Differentiate` (bounded rules), `Limit` (expansion of fixed `lorentz_factor` body → substitution → simplify).
- Gate: `go build ./ops`.

**Step 7 — Relation operations**
- Deps: 5. Files: `ops/relation.go`.
- Result: `Compare`, `Solve` (pattern-only), `SelectBranch` (exact §15.12 six-step contract).
- Gate: `go build ./ops`.

**Step 8 — Dispatch + full ops test suite**
- Deps: 6, 7. Files: `ops/dispatch.go`, `ops/operations_test.go`, `ops/negative_test.go`.
- Result: `OperationID`, `OperationParams` (exact v2.3 §15.13 schema — see plan §9), `Apply` closed switch over exactly 12 IDs (`identify` rejected); **positional enforcement of §15.13.1** (inputs bound by index, never reordered; unused positions omitted); complete positive+negative tests (C, D, E, F, `TestApplyPositionalInputs` covering all 12 IDs + reorder rejection, REQ-032-01..07, 18..20).
- Gate: `go vet ./ops && go test ./ops`.

**Step 9 — Corpus layer**
- Deps: 3. Files: `core/corpus.go`.
- Result: `Manifest` typed structs, `ManifestDomain/Limit/Anomaly/Reduction/Item`, `Challenge`, `Review` (**defined in `core/corpus.go` per specs_v2_3.md §27 — the sole declaration site**), `ParseManifest/ValidateManifestBytes/CanonicalManifestJSON/Hash` with `encoding/json` strict decoding (`DisallowUnknownFields`), enum validation, canonical-expr decode + byte-canonical re-encode check, **no constructor calls, no I/O**.
- Gate: `go vet ./core`.

**Step 10 — Mechanics corpus package**
- Deps: 8, 9. Files: `mechanics/primitives.go` (with `//go:embed manifest.json`), `mechanics/relations.go`, `mechanics/manifest.json`, `mechanics/manifest_test.go`, `mechanics/relations_test.go`.
- Result: 9 thin wrappers + 3 relation constructors; 11-item manifest in canonical bytes; static `map[string]func() core.Object` cross-check; acceptance **A, B, G**.
- Gate: `go test ./mechanics`.

**Step 11 — Relativity corpus package**
- Deps: 8, 9. Files: `relativity/primitives.go` (embed), `relativity/relations.go`, `relativity/manifest.json`, `relativity/manifest_test.go`, `relativity/derivation_test.go`.
- Result: **8 nominal wrappers** (`Spacetime, MinkowskiMetric, RestMass, Energy, ThreeMomentum, FourMomentum, SpeedOfLight, Velocity`), **3 fixed relation/function constructors** (`LorentzFactor, EnergyMomentumRelation, MassEnergyRelation`), **3 zero constructors** (`ZeroThreeMomentum, ZeroEnergy, ZeroVelocity`), `RestFrameAssumption`; **exactly 10 manifest items** (F1 — `Velocity` and the zeros are NOT manifest items); cross-check; acceptance **H, I, J** (full §20 sequence, fixed-body limit).
- Gate: `go test ./relativity`.

**Step 12 — Session core: state machine, draft buffer, steps, Identify**
- Deps: 8. Files: `session/session.go`, `session/ledger.go`, `session/session_test.go` (first half).
- Result: nine actions (§16.1: `Postulate, Declare, Define, Step, Identify, Conclude, Draft, Commit, Seal`), state machine, `Session.Step` → `ops.Apply` **using the §15.13.1 positional indices and validating `OperationParams` against the operation ID**, session-only `Identify`, `Commit` (indices, `step-000001`, genesis hash, `StepEnvelope` chain), `Conclude`, `Validate` (17-step replay reconstituting calls by the same positions), `ParseLedgerJSON`, ledger tamper detection, determinism.
- Gate: `go vet ./session && go test ./session` (K, Q, R partial).

**Step 13 — ResearchCandidate seal & unverified loader**
- Deps: 12. Files: `session/research_candidate.go`, remainder of `session/session_test.go`.
- Result: support types (`FrameworkDependency, Prediction, FalsificationCondition, RecoveryClaim, AnomalyReference` defined in `session`; `Challenge/Review` **consumed from `core` without any session declaration or alias**), `DraftMetadata`, `Session.Seal`, private seal/validation constructor shared by `Seal` and `UnverifiedResearchCandidate.Validate`, `ParseResearchCandidateJSON` returning only `UnverifiedResearchCandidate`.
- Gate: `go test ./session` (S, REQ-032-10a, seal-requires-hypothesis, post-seal mutation).

**Step 14 — Hypothesis package**
- Deps: 13. Files: `hypothesis/candidate.go`, `hypothesis/candidate_test.go`.
- Result: `NewCandidateConcept` forcing `HYPOTHESIS` + `NONE`, explicit `Kind`/`Dimension`; containment tests; falsifiability + anomaly tests through `Seal`; AST no-promotion/no-forbidden-export scan.
- Gate: `go test ./hypothesis` (L, M, N, REQ-024-01/02, REQ-032-09/10).

**Step 15 — Documentation**
- Deps: none (content final after 10–14). Files: `README.md`, `docs/paper-translation.md`.
- Result: README with module identity, dependency diagram, **constructor authority as external/public API boundary** (REQ-000-04), MRC list + `mrc-v0.4` fallibility (no runtime override), integrity-vs-authenticity statement, non-goals/deferrals; translation doc with exactly `Common notation`, `Framework mapping`, `Ambiguity resolution` sections (specs_v2_3.md §35).
- Gate: `go test ./...` (docs-section assertions in `core/object_test.go`).

**Step 16 — End-to-end acceptance**
- Deps: all. Files: none (verification only).
- Result: full §13 gate: build, vet, all tests incl. A–S + `TestApplyPositionalInputs`, file-tree audit, import-graph audits, export-surface audits, positional-mapping verification, determinism double-run.
- Gate: §13.

## 4. Core model implementation

**`internal/kernel/types.go`**
- `type Expr struct { /* unexported: kind ExprKind; name string; rat *big.Rat; children []Expr; op RelationOperator; functionID string */ }` — immutable handle; zero value invalid (`Valid()==false`).
- `ExprKind` values pinned by specs_v2_3.md §8.2.1 (`ExprSymbol`..`ExprBranchSet`) with ordinals 0–9; `RelationOperator` (`RelationEq`..`RelationGte`) with canonical strings `eq,neq,lt,lte,gt,gte` (no aliases).
- `type Kind uint8` with exactly the 18 values of specs_v2_3.md §6 in listed order (ordinals stable for `mrc-v0.4`), `String()` returning the §6 names (`Mass` … `BranchSet`) — these strings are the manifest/JSON `kind` values.
- `type Dimension struct { /* [7]*big.Rat exponents M,L,T,I,Θ,N,J, immutable */ }` — every arithmetic returns reduced rationals with positive denominator; `Valid/Equal/Multiply/Divide/Pow(*big.Rat)/CanonicalJSON/Hash` implemented on the kernel type; `core` provides the nine fixed constructors (`Dimensionless()` … `DimensionEnergy()`) as façades over the Step 1 kernel-side constructor.
- Metadata types: `Assumption` (kind, key, closed union `TextValue(string)|ExprValue(Expr)`), `AssumptionSet` (immutable, F4 canonical-bytes ordering, dedup), `Convention{key,value}`, `ConventionSet` (F4 ordering), `Provenance{status,source,framework,parentHashes,assumptionHash,conventionHash,mrcVersion,justification}`, `CorpusStatus`.
- **Assumption kinds — exact v2.3 §11.1 set (replaces all v2.2 kinds):**

  | Go value | Canonical JSON |
  |---|---|
  | `Domain` | `domain` |
  | `Regime` | `regime` |
  | `Constraint` | `constraint` |
  | `Convention` | `convention` |
  | `Approximation` | `approximation` |
  | `MathPrecondition` | `math_precondition` |
  | `PhysicalAssumption` | `physical_assumption` |

  Exactly these seven; `Approximation` is reserved (no MVP operation mints `APPROXIMATED` provenance). `AssumptionKind.Convention` is a reserved category (conventions themselves live in `ConventionSet`).
- Canonical JSON encoders/decoders for every type (expr forms exactly specs_v2_3.md §10.2; dimension field order §10.3; rational strings `"1/2"`, integers `"2/1"`; uppercase provenance status; lowercase corpus status and assumption kind; `branch_set` expr kind).
- Child sorting: (1) node-kind ordinal, (2) lowercase hex SHA-256 child hash via `bytes.Compare`, (3) canonical child bytes via `bytes.Compare`; never maps/pointers/locale.
- Hashing: SHA-256 over canonical bytes; hex helper lowercase.
- Kernel helper ownership (step 1): `EntailsNonNegative/EntailsPositive/EntailsNonZero`, `AssumptionSet.Merge`, `ConventionSet.Merge`, metadata equality/canonicalization/hashing.

**`internal/kernel/mint.go`**
```go
type ObjectSpec struct { Name string; Kind Kind; Dimension Dimension; Expr Expr
    Assumptions AssumptionSet; Conventions ConventionSet; Provenance Provenance; CorpusStatus CorpusStatus }
func MintObject(spec ObjectSpec) (Object, error)
```
- Exactly one production mint entry point (§5.2.1): rejects invalid expr handle; rejects empty name *where the constructor contract requires one* (fixed domain/corpus constructors and hypothesis `id` require names; operation results mint with `Name == ""` — pinned rule); rejects invalid dimension; validates provenance (known status, valid hex parent hashes, `mrcVersion` string, justification only for `IDENTIFIED` or `HYPOTHESIS`-with-justification, `DERIVED` with empty parents permitted for corpus artifacts per §19.6); validates assumption/convention canonicality; normalizes to canonical internal form; returns typed error otherwise.
- Object canonical JSON field order (§10.4): `schema_version, valid, name, kind, dimension, expr, assumptions, conventions, provenance, corpus_status`.
- **`schema_version`:** required by v2.3 as the first canonical serialization field, but it is **serialization DTO state, not `kernel.Object` state** — it is never stored in the kernel struct and never computed by the kernel. Pinned value for MVP: `"1"` (string) for every artifact schema (object, ledger, session, candidate, manifest). **This value is a Plan 10 implementation pin** — specs_v2_3.md mandates the field and its position, not the value `"1"`.
- Replay decoding is **internal to `internal/kernel`** (no exported `DecodeObjectJSON` or `ParseObject` in any package). The internal decoder validates the same authoritative invariants `MintObject` validates before a decoded object may enter replay — a canonical byte sequence alone is insufficient if the decoded object violates kernel invariants — and guarantees `Encode(Decode(b)) == b`.

**`core` façade**
- `core/object.go`: `type Object = kernel.Object`; `EqualObject`, `HashObject`; `Kind`, `CorpusStatus` aliases + constants. Accessors are kernel methods: `Valid/Name/Kind/Dimension/Expr/Assumptions/Conventions/Provenance/CorpusStatus` (REQ-005-09); all return immutable copies (slices/`*big.Rat` defensive copies).
- `core/expr.go`: `type Expr = kernel.Expr` + the ten `New*` constructors (façade over kernel structural canonicalization) + inspection methods `Valid/Kind/SymbolName/RationalValue/Children/Base/Exponent/FunctionID/Arguments/RelationOperator/Left/Right/BranchTarget/Branches` (kind-guarded; empty/invalid where inapplicable; `RationalValue()`, `Exponent()`, and slice accessors return defensive copies per §8.0).
- **Authority boundary mechanics:** `kernel.Object` fields are unexported; outside-module callers cannot import `internal/kernel` (REQ-004-03, Go language rule); `core` re-exports no mint function (REQ-005-04/10); within-module packages (`ops`, `session`, domains, tests) call `kernel.MintObject` and are trusted repository code per REQ-000-03/04. The zero `core.Object{}` is invalid; there are no mutators; validation runs inside `MintObject` before any valid object becomes visible (REQ-005-11).
- **Visibility contract:** `internal/kernel` exports controlled canonical helper functions to in-module packages; `core` exposes the public canonical façade; object decoding stays internal-only. No package reads another package's unexported Go identifiers.
- Fixed data pins (all deterministic, no open choices):
  - **Manifest item names** (11 mechanics: `mass, time, position, velocity, acceleration, force, momentum, energy, Newton's second law, momentum relation, kinetic energy relation`; 10 relativity: `spacetime, Minkowski metric, rest mass, energy, three-momentum, four-momentum, speed of light, Lorentz factor, energy-momentum relation, mass-energy relation`).
  - **Non-manifest constructor object names (F1):** `velocity` (wrapper), `zero three-momentum`, `zero energy`, `zero velocity` — these objects exist but have NO manifest entry.
  - Provenance per constructor: status `DEFINED` for fixed corpus constructors except `MassEnergyRelation()` = `DERIVED` (§19.6); `ParentHashes` empty; `Framework` = `classical_mechanics` / `special_relativity`; `Source` pinned per manifest `source` field (see §7); `MRCVersion` = `mrc-v0.4`; `Justification` empty.
  - `CorpusStatus`: fixed corpus constructors → `ESTABLISHED`; all operation/`Identify` results → `NONE`; candidate concepts → forced `NONE` (§13.4, §24.2).

## 5. MRC implementation map

| Rule | Exact semantic check | Enforcement location | Failure | Test |
|---|---|---|---|---|
| MRC-001 Constructor/carrier integrity | Valid object only via allowed §5.4 paths; unexported `kernel.Object` fields; no public generic factory; zero object invalid; `MintObject` invariant validation before visibility | `internal/kernel/mint.go` (validation); `core/object.go` (no forwarder); fixed domain/hypothesis constructors; `ops` outputs; `session` authority; Go `internal/` rule | `InvalidObjectError`, or `ProvenanceError` when a public API rejects a construction path | `core/object_test.go::TestNoGenericFactory`, `TestZeroObjectInvalid`, `TestMintObjectRejectsInvalidSpec`, `TestKernelObjectFieldsUnexported`, `TestInternalKernelBoundary` (external pkg, REQ-003-03) |
| MRC-002 Dimensional compatibility | Add/Subtract/Compare equal dimensions; Multiply/Divide compose; `Pow` raises by `*big.Rat`; Substitute replacement dim = variable dim; Differentiate dim = target/wrt | `ops/arithmetic.go`, `ops/transform.go`, `ops/relation.go` (precondition pipeline) | `DimensionMismatchError` | `ops/negative_test.go::TestDimensionMismatch` (C), `TestSubstituteDimensionMismatch` |
| MRC-003 Physical-kind compatibility | Exact §6.2 table: Add/Subtract same named Kind+Dim / Expression+Expression only; Compare per eq-table; inequalities restricted to ordered kinds {Mass,RestMass,Time,Energy,KineticEnergy,SpeedOfLight,Expression}; Multiply/Divide/Pow → `Expression`; Solve input `Relation`/symbol target → `BranchSet`; SelectBranch `BranchSet`+`Relation` | same as MRC-002 | `CategoryMismatchError` | `ops/negative_test.go::TestKindMismatchEqualDimensions` (D — **fixture: `Mass` vs `RestMass`, equal dimension `M`, distinct named kinds; isolates MRC-003 without triggering MRC-002**), `ops/operations_test.go::TestKindCompatibilityTable` (incl. G's Expression-vs-named path) |
| MRC-004 Assumption compatibility | Merge = union of inputs + operation + introduced assumptions (F4 ordering, dedup); same `(Kind,Key)` different canonical value → conflict | `internal/kernel` (`AssumptionSet.Merge`) via `core/assumption.go` façade, used by all `ops` | `AssumptionConflictError` | `ops/negative_test.go::TestAssumptionConflict` (E) |
| MRC-005 Convention compatibility | Merge conventions (F4 ordering); same key different value → conflict | `internal/kernel` (`ConventionSet.Merge`) via `core/convention.go`, used by all `ops` | `ConventionConflictError` | `ops/negative_test.go::TestConventionConflict` (F) |
| MRC-006 Explicit physical identification | Only `Session.Identify`: two valid operands, equal dims, Compare-kind rules, non-empty trimmed justification; logs `Identification` step; `HYPOTHESIS` operand forces `HYPOTHESIS` output | `session/session.go::Identify` (there is no `ops.Identify`; `ops.Apply` rejects `identify`) | `IdentifyError` | `session/session_test.go::TestSessionIdentifyRecords`, `TestIdentifyRequiresJustification` (K, REQ-032-08) |
| MRC-007 Session/provenance authority | Only `session.Session` commits ledger / seals `ResearchCandidate`; no exported arbitrary-field `Step`/`Ledger` constructors (storage unexported; loader is validation-only) | `session/session.go`, `session/ledger.go`, `session/research_candidate.go` | `ProvenanceError` (also `LedgerValidationError` for structural failures) | `session/session_test.go::TestLedgerConstructionViaSessionOnly`, Q, S (REQ-032-11..14) |
| MRC-008 Candidate containment | Any op/session result depending on a `HYPOTHESIS` input is `HYPOTHESIS` (transitive); artifact presenting hypothesis-dependent output as trusted fails validation | `ops` provenance law (`arithmetic.go` shared helper, all 12 ops + `Identify`); `session` candidate validation pipeline | `CandidateContainmentError` | `hypothesis/candidate_test.go::TestHypothesisContamination` (L), `TestCandidateArtifactContainmentRejected` (REQ-032-09), S |

Enforcement-layer distinctions: *Go package/internal visibility* = MRC-001 (+REQ-004-03); *trusted mint authority* = `kernel.MintObject` only (MRC-001); *operation-time MRC* = MRC-002/003/004/005/008 inside `ops`; *session authority* = MRC-006/007 inside `session`; *candidate containment at artifact level* = MRC-008 in `session/research_candidate.go`; *`physvet`* = absent from MVP (REQ-002-15), its v0.5 contract recorded only in specs_v2_3.md §36/§14 of this plan.

## 6. Symbolic engine

**Node set (REQ-008-01):** `Symbol, Rational, Add, Mul, Neg, Pow, Sqrt, Call, Relation, BranchSet` only; ordinals 0–9 stable; forbidden nodes `Derivative, Limit, Function, Arbitrary, Eval, Callback, RawString` do not exist; `Expr` stores no callbacks (REQ-008-02) — assert by reflect type scan.

**Exact arithmetic:** all coefficients/exponents `math/big.Rat`, canonical `num/den` reduced, positive denominator, zero `0/1`; no `float64` anywhere in non-test sources (scan test).

**Constructor-time vs Simplify-time boundary (explicit v2.3 §8.5/§9.6 resolution, review 🟡-2):**

```text
Constructor-time (structural only):
    flatten same-kind Add/Mul
    combine exact rationals (§9.4)
    sign normal form: Mul(-1,x)→Neg(x), Neg(Rational(q))→Rational(-q), Neg(Neg(x))→x
    deterministic child sort

Simplify-time (semantic rewrites):
    Add/Mul identities, zero/one rules
    rational Pow rewrites (§9.6), Pow(x,0)→1 (§8.7)
    repeated powers (§9.7)
    Sqrt rewrites (§9.8)
    relation recursive simplification (§9.9)
```

This determines the canonical bytes of the unsimplified `Limit` body, which retains `Mul(1, …)` until `Simplify`.

**Canonical equality/hashing:** `EqualExpr` structural (never display strings); `EqualObject` = byte equality of object canonical JSON (all metadata); SHA-256 over canonical bytes, lowercase hex; deterministic across runs (fixed-vector tests).

**Canonical JSON:** expr forms exactly §10.2 (fixtures assert the literal example bytes); object order §10.4; dimension order §10.3; typed structs only, never `map[string]any`.

**`Simplify` rule set (bottom-up recursive, then re-normalize + re-sort):**
- §9.6 list: `Add(x,0)→x`, `Add()→0`, `Mul(x,1)→x`, `Mul()→1`, `Mul(x,0)→0` (finiteness gate), `Pow(x,1)→x`, `Neg(Neg(x))→x`, `Sqrt(1)→1`, `Sqrt(0)→0`.
- **§8.7 (review F2): `Pow(x,0)→1` only when `x` is a valid nonzero-safe base under the operation's assumptions.** Never unconditional; `0^0` remains `UnsupportedOperationError`.
- Exact rational `Pow` (§9.6): `Pow(Rational(q),n)→Rational(q^n)` for integer `n≥0`; for integer `n<0` when `q≠0`; `Pow(Rational(0),n<0)` → `UnsupportedOperationError`; `Pow(Rational(1),n)→Rational(1)` any exact `n`; `Pow(Rational(0),e>0)→Rational(0)`; `0^0` unsupported (error, never silently valued).
- §9.7: `x·x→Pow(x,2)`, `Pow(x,2)·Pow(x,2)→Pow(x,4)`, `Pow(Pow(x,a),b)→Pow(x,a·b)` only for nonnegative-integer rational `a,b` — no other power algebra (§9.7-MUST-02).
- §9.8: `Sqrt(Rational(q))` → exact root only for non-negative perfect-square rationals; `Sqrt(Pow(x,2))→x` only when entailment holds.
- **§9.9 — `Simplify` on `Relation`/`BranchSet` (review 🟡-3):** `Relation(eq,a,b)` and other relations are recursively simplified on **both sides while preserving the relation operator**; `BranchSet` is recursively simplified within its bounded structure; `Simplify` **never** creates `IDENTIFIED` provenance and no relation simplification is interpreted as physical identification. No general equation reasoning.

**Finiteness predicate (pins §9.6 "finite under current assumptions" — required by the golden limit trace; v2.3 reconciliation B):** `finite(e, A)` = `Rational`→true; `Symbol`→true (primitive physical quantity); `Add/Mul/Neg`→all children finite; `Pow(b,n)` with nonneg-integer `n`→finite(b); `Pow(b,n)` negative-integer `n`→finite(b) ∧ `A` entails `b≠0`. The entailment distinguishes:

```text
base != 0   → non-singularity/admissibility   (admits c^-1 under c > 0)
base > 0    → sign-sensitive rewrite          (Sqrt and sign rules only)
```

`0^0` → `UnsupportedOperationError`; `0^-1` is rejected by the `b≠0` gate. No Lorentz-specific simplification hack exists — the general finite-factor rule admits `c^-1` because `c>0` entails `c≠0`. `Sqrt(x)`→`A` entails `x≥0` ∧ finite(x), or `x`∈{Rational 0,1}; `Call`→expand fixed body, then finite; `Relation/BranchSet`→false. `Mul` with a `Rational 0` factor → `0` iff every other factor is finite.

**Bounded differentiation (§15.8):** `Rational→0`; `Symbol→1` if same symbol else `0`; `Add→` termwise; `Mul→` n-ary product rule; `Neg→Neg(d)`; `Pow(base,n)` nonneg-integer rational `n → n·base^(n-1)·d(base)`; result `Simplify`d; `Kind=Expression`; dimension `target/wrt`. Unsupported (`Sqrt, Call, Relation, BranchSet`, non-integer exponent) → `UnsupportedOperationError`. No `Derivative` node is created.

**Limit (§15.9):** if `Call(lorentz_factor, args)` present → replace by fixed body `1 / Sqrt(1 - Pow(v/c, 2))`, substituting the call's argument expressions for parameter symbol `v` (body built with `NewMul(NewRational(1), NewPow(NewSqrt(NewAdd(NewRational(1), NewNeg(NewPow(NewMul(v, NewPow(c,-1)), 2)))), -1))`); then direct substitution of `value` for `variable`'s single symbol; then `Simplify`; singular/unsupported → `UnsupportedOperationError`; result `Kind=Expression`, dimension = target dimension after substitution. **Division representation pinned: `a/b ≡ Mul(a, Pow(b,-1))`** (no division node). The implementation must traverse the body — returning `1` on function-ID match alone is non-conforming (§15.9.1-MUST-02; secondary verification = source audit of `ops/transform.go`).

**Pattern-only `Solve` (§15.11):** accept exactly `Relation(eq, Pow(Symbol(t),2), Expr)` with target object expr a single `Symbol` equal to `t`; return `BranchSet(targetExpr, [Sqrt(rhs), Neg(Sqrt(rhs))])`, kind `BranchSet`, dimension = target dimension; everything else → `UnsupportedOperationError`.

**`SelectBranch` — exact v2.3 §15.12 contract (review 🟡-4):**

```text
Shape:  BranchSet(target, [positiveBranch, Neg(positiveBranch)])
        + Relation(gte, target, zero)      (zero: matching kind and dimension)
```

The operation performs exactly these six steps:

1. validate branch-set shape (operands valid; `branches.Kind()==BranchSet`; `constraint.Kind()==Relation`; branches structurally `[b, Neg(b)]`);
2. validate constraint compatibility (operator `gte`; `constraint.Left() == branches.BranchTarget()`; `constraint.Right()` is `Rational 0`; dimensions match);
3. merge the constraint into assumptions as `Assumption{Kind: Constraint, Key: "selected_branch/" + hex(HashExpr(constraint.Expr())), Value: ExprValue(constraint.Expr())}`;
4. select the first (structurally nonnegative `Sqrt`) branch and simplify it under the merged assumptions;
5. preserve provenance subject to contamination;
6. return `Kind=Expression`, dimension = branch-set dimension.

Mass-energy case required: `Sqrt(Pow(Mul(m,Pow(c,2)),2)) → m·c²` because `m≥0` (from `RestMass`) and `Pow(c,2)` is structurally nonnegative (even exponent) per §9.8. **No sign-entailment-failure error is invented:** where v2.3 specifies none, the plan states only the narrow behavior required by the mass-energy path; `UnsupportedOperationError` is used only where the broader operation contract already requires it. No branch-selection engine exists.

**Defensive copies (review 🟡-5):** `NewRational` copies its input `*big.Rat`; `NewPow` copies its input exponent; `Pow` clones/reads the caller-supplied exponent (§15.5 — never retains the pointer); `Expr.RationalValue()` and `Expr.Exponent()` return copies; all slice-returning accessors return copies. No caller-owned mutable pointer becomes part of immutable kernel state.

**No CAS/root-finder/prover/parser/function registry:** `Solve` is one pattern; `lorentz_factor` is the sole `function_id`, hard-coded in `ops/transform.go`; entailment is the closed §9.8 rule list ("no other logical inference is permitted") implemented by the kernel `Entails*` helpers; no source text is ever parsed as mathematics.

## 7. Physics corpus

- **Files:** `mechanics/manifest.json`, `relativity/manifest.json`, embedded via `//go:embed manifest.json` in each package's `primitives.go` (blank `import "embed"`), exposed to tests via the internal test package — no new export surface, no new files.
- **Canonical file bytes:** each manifest file is byte-identical to `CanonicalManifestJSON(ParseManifest(file))` (compact, struct field order). `Manifest.Hash()` = SHA-256 of those canonical bytes. Tests assert file-bytes == canonical re-encoding (guards drift) and this hash is what callers put in `FrameworkDependency.ManifestHash`.
- **Strict typed structs:** `core/corpus.go` per specs_v2_3.md §17.0 (`Manifest`, `ManifestDomain`, `ManifestLimit`, `ManifestAnomaly`, `ManifestReduction`, `ManifestItem`) + `Challenge`, `Review` (§27, sole declaration site). Decoding uses `encoding/json` with `DisallowUnknownFields`; missing required fields, unknown fields, unknown enum values, non-canonical `canonical_expr`/`dimension` (must equal their canonical re-encoding), duplicate item IDs, empty `constructor` → `ManifestValidationError`.
- **`constructor` mapping:** every item carries `constructor` (e.g. `mechanics.NewMass`, `mechanics.NewtonSecondLaw`, `relativity.LorentzFactor`); each manifest test defines the static test-only `map[string]func() core.Object` resolving **every** constructor id; missing entry or unresolvable id fails the test (§17.8). No reflection registry.
- **Canonical expression decoding:** `ManifestItem.CanonicalExpr` decoded through the kernel expr decoder into the closed node set; `statement` is stored as an opaque string and **never parsed** (§17.6; verified by `TestStatementNotParsed` — a `statement` containing arbitrary math text validates fine, and no non-test file imports `go/parser`/`go/scanner`/`go/token`).
- **Constructor cross-check (§17.8):** nine-step test — parse manifest bytes → resolve `constructor` via static map → call → compare (1) `canonical_expr` structurally + by hash, (2) `dimension`, (3) `kind`, (4) `provenance_status`, (5) `assumptions`, (6) `source`/framework metadata.
- **Required items (F1 reconciliation):** mechanics **exactly 11** (§28.1: Mass, Time, Position, Velocity, Acceleration, Force, Momentum, Energy, NewtonSecondLaw, MomentumRelation, KineticEnergyRelation); relativity **exactly 10** (§28.2: Spacetime, MinkowskiMetric, RestMass, Energy, ThreeMomentum, FourMomentum, SpeedOfLight, LorentzFactor, EnergyMomentumRelation, MassEnergyRelation). Spec says "at least"; this plan pins the minimum set — **`Velocity`, `ZeroThreeMomentum`, `ZeroEnergy`, `ZeroVelocity`, `KineticEnergy` (parameterized), and `RestFrame` (assumption) are NOT manifest items.**
- **Pinned manifest data:** both frameworks `corpus_status: "established"` (loaded, never derived — REQ-013-04); mechanics `framework_id/classical_mechanics`, name `Classical Mechanics`, `domain[1]` (`classical_nonrelativistic`), `limits[1]` (`nonrelativistic_scope`), `anomalies[1]` (`nonrelativistic_regime`, framework `classical_mechanics`, status `scope_limit`); relativity `framework_id/special_relativity`, name `Special Relativity`, `domain[1]`, `limits[1]`, `anomalies[1]` = the exact §19.9 record (`id: no_gravity`, status `scope_limit`); relativity framework `assumptions[4]` — text values exactly `Minkowski spacetime, Lorentz symmetry, No gravitational dynamics in package, Special-relativistic regime` (§19.3), with plan-pinned kinds/keys: kind `domain`/key `minkowski_spacetime`, kind `physical_assumption`/key `lorentz_symmetry`, kind `domain`/key `no_gravitational_dynamics`, kind `domain`/key `special_relativistic_regime`; relativity convention recorded on `NewMinkowskiMetric()` and in manifest: `metric.signature = -+++` (§19.3).
- **`source` pins:** `NewtonSecondLaw → "Newton, Principia"` (spec §17.5 example); remaining pins are Plan 10 data pins (spec-silent): all other mechanics items → `"Classical Mechanics corpus"`; `EnergyMomentumRelation`, `MassEnergyRelation` → `"Einstein, 1905"`; other relativity items → `"Special Relativity corpus"` (must equal object `Source()` for cross-check).
- **Item provenance pins:** all items `DEFINED` except `MassEnergyRelation` = `DERIVED` with `derivable_from: ["EnergyMomentumRelation","RestFrame"]` (entries validated as non-empty strings only — `RestFrame` intentionally resolves to no item, so cross-reference existence is NOT checked; §19.6).
- **Anomaly/limitation records** are descriptive metadata; no runtime inference (§34.3).

## 8. Assumptions, conventions, provenance, and candidate containment

- **Assumption kinds:** the exact seven of §11.1 (`Domain, Regime, Constraint, Convention, Approximation, MathPrecondition, PhysicalAssumption`; JSON `domain, regime, constraint, convention, approximation, math_precondition, physical_assumption`) — no other kind value may exist (see §4 table).
- **Deterministic set ordering (F4):** both sets store elements canonicalized → sorted lexicographically by canonical bytes → deduplicated; never Go map order; ordering participates in canonical JSON → set hash → object hash → ledger hash chain → candidate determinism (slice R).
- **Assumption merge (§11.4):** result = union(inputs) + required operation assumptions + explicitly introduced assumptions; exact duplicates deduplicated; same `(Kind,Key)` with different canonical value → `AssumptionConflictError`. No subsumption/theorem proving.
- **Denominator precondition (§11.5):** `Divide(a,b)` adds `Assumption{MathPrecondition, "denominator/"+hex(HashExpr(b.Expr())), ExprValue(Compare(b,0,neq))}` — distinct keys per denominator.
- **Sign assumptions (§11.6):** `RestMass()` carries `Constraint/rest_mass_nonnegative = (m >= 0)`; `SpeedOfLight()` carries `Constraint/speed_of_light_positive = (c > 0)`; merged transitively through operations (so `EnergyMomentumRelation` carries both — §19.5-MUST-02).
- **Bounded entailment (§9.8, kernel-owned `Entails*` helpers, used only for sign/nonzero preconditions):** `Rational q≥0`→nonneg; `Symbol` with explicit `x≥0`→nonneg; explicit `x>0`→positive; `Mul` all factors nonneg→nonneg; `Pow(x, even nonneg integer)`→nonneg; positive-coeff × nonneg factor → nonneg; positive entails nonzero; `EntailsNonZero` additionally accepts a direct structural match of an assumption expr `Relation(neq, e, 0)`. Nothing else — no derived inference beyond this closed list.
- **Conventions (§12):** `Convention{Key,Value}` non-empty strings; merge conflicts (same key, different value) → `ConventionConflictError`; conventions serialize/hashes participate in object canonical form; no equation may live in a convention string (REQ-012-02 — type system permits only strings).
- **Provenance propagation law (§13.2) — exact, no invented ordering:** pure `ops` (all 12): any input `HYPOTHESIS` → output `HYPOTHESIS`, else `DERIVED`. `Session.Identify`: any input `HYPOTHESIS` → `HYPOTHESIS`, else `IDENTIFIED`. Assertion actions: `Postulate` requires and records `POSTULATED`; `Define` requires and records `DEFINED`; `Declare` preserves existing status. `APPROXIMATED` exists in the enum but no MVP operation mints it (REQ-013-01). Only `Session.Identify` creates `IDENTIFIED`. `HYPOTHESIS` is never raised to a trusted status by any operation. Source/Framework inheritance (§13.5): preserve iff all inputs identical and non-empty, else empty — purely mechanical.
- **Corpus status separation:** a distinct human-curated axis (`NONE/ESTABLISHED/CONTESTED/SUPERSEDED/FALSIFIED`, lowercase JSON); loaded from manifests, never computed, never revised at runtime (REQ-013-04); no API produces `PHYSICALLY_TRUE/TRUTH_SCORE/PROBABILITY_OF_TRUTH/BEST_THEORY` (REQ-013-05).
- **Hypothesis contamination law (MRC-008):** input-level (ops + Identify status law) and artifact-level: candidate validation recomputes provenance per step; any hypothesis-dependent step output presented with a trusted status, any `Hypothesis` field not `HYPOTHESIS`, or any hypothesis-dependent output carrying corpus status ≠ `NONE` → `CandidateContainmentError`. There is no promotion API of any kind.

## 9. Derivation ledger

- **`session.Session`** (state machine, REQ-016-01): `New → Drafting → Committed → Concluded → Sealed`.
  - Exactly nine actions (§16.1): `Postulate, Declare, Define, Step, Identify, Conclude, Draft, Commit, Seal` (+ read-only `Validate`, `CanonicalJSON`).
  - `Draft(derivationID, label, DraftMetadata)` once, `New→Drafting`; `DerivationID` non-empty caller-supplied (deterministic, never generated); records `DerivationID, Label, MRCVersion, GenesisHash`; all slices copied at Draft time; **no separate candidate-metadata mutation methods exist** (§16.2.1).
  - `Postulate/Declare/Define/Step/Identify` legal only in `Drafting`; `Commit` only in `Drafting` (empty draft → `LedgerValidationError`, state unchanged); `Conclude` only in `Committed`; `Seal` only in `Concluded`; `Validate` legal in every state, never mutates; after `Sealed` every mutation → `ProvenanceError` (REQ-032-22).
  - Assertion actions record without `ops.Apply`; `Postulate` requires status `POSTULATED` (exercised via `kernel.MintObject` fixture from `session_test` — the single mint path, no second helper), `Define` requires `DEFINED`, `Declare` requires valid and preserves status.

- **Operation invocation contract — exact v2.3 §15.13/§15.13.1 (v2.3 reconciliation C):**

  `ops.Apply`, every `Session.Step` record, and every replay step in `Session.Validate` MUST use **exactly** the following positional schema. Inputs are never reordered or packed differently; unused input positions are omitted; replay decoders reconstitute calls using exactly these indices.

  | `OperationID` | `inputs[0]` | `inputs[1]` | `inputs[2]` | `params` fields used | Result |
  |---|---|---|---|---|---|
  | `add` | `a` | `b` | — | — | Sum `core.Object` |
  | `subtract` | `a` | `b` | — | — | Difference `core.Object` |
  | `multiply` | `a` | `b` | — | — | Product `core.Object` |
  | `divide` | `a` | `b` | — | — | Quotient `core.Object` |
  | `pow` | `base` | — | — | `params.Exponent` (exact rational `"num/den"`) | Power `core.Object` |
  | `simplify` | `target` | — | — | — | Simplified `core.Object` |
  | `substitute` | `target` | `variable` | `replacement` | — | Substituted `core.Object` |
  | `differentiate` | `target` | `wrt` | — | — | Derivative `core.Object` |
  | `limit` | `target` | `variable` | `value` | — | Limit `core.Object` |
  | `compare` | `a` | `b` | — | `params.Operator` (`RelationOperator`) | Relation `core.Object` |
  | `solve` | `relation` | `target` | — | — | BranchSet `core.Object` |
  | `select_branch` | `branches` | `constraint` | — | — | Selected branch `core.Object` |

  The exact `OperationID` values: `add, subtract, multiply, divide, pow, simplify, substitute, differentiate, limit, compare, solve, select_branch` — the dispatch switch is closed over exactly these IDs; `identify` is deliberately excluded from `ops.Apply` (session-owned); no plugin registry or runtime registration exists.

  `OperationParams` — exact field order and canonical JSON (v2.3 §15.13):

  ```go
  type OperationParams struct {
      Kind          string
      Exponent      string
      Operator      core.RelationOperator
      Justification string
  }
  ```

  ```json
  {"kind":"pow","exponent":"2/1","operator":"eq","justification":""}
  ```

  Permitted `Kind` values: `empty | pow | compare | identify`. Binding: `pow` → `Exponent` required; `compare` → `Operator` required; `identify` → `Justification` required (stored in a session identification step, **never accepted by `ops.Apply`**); `empty` → all other operations and assertion actions. Unused fields encode their empty/default values. `Session.Step` validates params against the operation ID before calling `ops.Apply`; assertion steps record `Kind:"empty"`; `Session.Identify` records `Kind:"identify"` with the justification.

- **`Session.Step`:** validate state → validate inputs → reject `operation == "identify"` → validate `OperationParams` against op id (per the binding above) → `ops.Apply` (§15.13.1 indices) → build canonical step record (inputs, params, output all retained) → append to draft buffer → return result. Callers can never supply their own output for a mechanical operation.
- **`Session.Identify`:** the only identification API (no `ops.Identify`, no ambient/global session state — REQ-015-01/MRC-006): state check; both objects valid; equal dims; Compare-kind rules; non-empty trimmed justification → else `IdentifyError`; builds `Relation(eq, a.Expr(), b.Expr())`; status `IDENTIFIED` unless contamination; merged assumptions/conventions; justification stored in provenance; appends `Identification` step (params `Kind:"identify"`); returns object.
- **Draft buffer:** `Step`/`Identify` append fully specified entry material (label, kind, op, canonical inputs, params, canonical output, metadata hashes, status); entries unexported and immutable externally (§16.8).
- **`Commit`:** freeze buffer in order; indices from 1; `StepID = step-000001`-style (derived solely from index, no UUIDs; §16.16); first `PreviousStepHash` = 64 lowercase zeros (genesis; §16.15); compute each `CurrentStepHash`; retain canonical inputs/outputs/params; draft emptied; ledger immutable through public API.
- **Step fields, exact order (§16.12):** `StepID, Index, Label, StepKind, Operation, InputHashes, InputCanonicals, ParamsCanonical, OutputHash, OutputCanonical, AssumptionHash, ConventionHash, ProvenanceStatus, MRCVersion, PreviousStepHash, CurrentStepHash`. `StepKind ∈ {Assertion, Transformation, Identification}`; assertion operations `postulate|declare|define`; transformation = the 12 `OperationID`s; identification = `identify`. Assertion steps pin: `InputHashes=[h]`, `OutputHash=h`, `OutputCanonical` = same canonical object (retained-object consistency on replay).
- **`StepEnvelope`:** `CurrentStepHash = SHA256(CanonicalJSON(StepEnvelope{PreviousHash, Step: body}))` where body = all Step fields except `CurrentStepHash` (no self-reference; §16.17). Derivation hash = last step's `CurrentStepHash` (§16.18).
- **Replay substrate:** hashes verify integrity; `InputCanonicals[]` + `OutputCanonical` + `ParamsCanonical` make replay possible (hashes are not invertible — §16.13). Invariants `HashObject(decoded InputCanonical)==InputHash` and same for output are asserted per step. Decoding uses the internal kernel decoder with full invariant validation (plan §4).
- **`Session.Validate` — the exact 17-step ordered pipeline of §16.19** (structure/state → genesis → StepID/index → decode inputs → verify InputHash → decode outputs → verify OutputHash → assumption/convention hashes → recompute chain → verify params vs operation ID → replay assertions → replay transformations via `ops.Apply` at §15.13.1 indices → replay identification → compare replayed output hash → provenance/MRC → conclusion hash → candidate containment), any mismatch → `LedgerValidationError` (wrapping the specific typed error).
- **Tamper detection (§16.20):** (1) change only `OutputCanonical` → hash mismatch; (2) also recompute `OutputHash` → current-step hash mismatch; (3) recompute the full chain → replay divergence (re-executed `ops.Apply` output ≠ retained output, unless identical). Tests use `ParseLedgerJSON` (load-as-is, never repairs/recomputes) + `Ledger.Validate`.
- **Integrity ≠ authenticity:** the SHA-256 chain detects mutation only; anyone rewriting the whole ledger can recompute hashes; signing/authenticity is out of MVP — stated in README and in candidate validation messages; `ResearchCandidate.Validate` explicitly does not treat self-recomputed hashes as authenticity proof (§26.10).
- **Seal rules:** non-candidate derivations stop at `Concluded` and are checked with `Session.Validate()`; `Seal` requires `DraftMetadata.Hypothesis` valid with status `HYPOTHESIS` (invalid/absent/non-hypothesis → `ProvenanceError`, state stays `Concluded`); before sealing, full validation runs (§16.11); on success ledger becomes immutable, `DerivationHash` = last step's `CurrentStepHash`, state `Sealed`, `ResearchCandidate` returned.
- **How replay avoids an import cycle:** ledger/step/replay types live in `session`, which imports `core`, `ops`, and `internal/kernel`. `core` never imports `ops` or `session`, so `ops.Apply` can be re-executed during `session.Validate` without `ops` knowing anything about ledgers. `session` never imports `mechanics/relativity/hypothesis`; replay decodes retained canonical objects through the internal kernel decoder (mint/decode authority) and feeds them back into `ops.Apply`. Cycle-freedom is verified by `go build ./...` (REQ-004-01) plus `TestSessionImportIndependence`.
- **External parsing is unverified:** `ParseLedgerJSON` returns a `Ledger` for validation only — it creates no session authority and mints no objects outside the loaded records; `ParseResearchCandidateJSON` returns only `UnverifiedResearchCandidate` (no `core.Object` accessors, never `ResearchCandidate`).

## 10. Hypothesis and ResearchCandidate

- **Candidate concept (`hypothesis.NewCandidateConcept(id, kind, dimension, expr, assumptions, conventions) (core.Object, error)`):** `id` non-empty; `Kind` and `Dimension` explicit and caller-supplied (REQ-024-02), immutable after construction; `expr` valid; provenance **forced** to `HYPOTHESIS` with `Source="hypothesis"`, `Framework=""`, `MRCVersion="mrc-v0.4"`; corpus status **forced** `NONE`; no caller-supplied provenance; no promotion API (REQ-024-01). Mints via `kernel.MintObject`.
- **Trusted `ResearchCandidate` ownership:** constructed exclusively inside `session`. `Session.Seal()` is the canonical production minting action. `UnverifiedResearchCandidate.Validate()` **delegates to the same unexported sealing/validation constructor** (one private path — no independent trusted-construction logic); `ParseResearchCandidateJSON` performs outer-schema validation only, stores canonical bytes, exposes `CanonicalJSON()` and `Validate()` and **no `core.Object`-returning method** (REQ-032-10a).
- **Public surface exactly specs_v2_3.md §26.10** (14 accessors + `Validate` + `CanonicalJSON`; collection accessors return copies; unexported storage; no public struct literal/constructor) + unverified surface exactly as specified.
- **Validation pipeline (§26.8), in order:** schema → canonicalization → ledger-hash validation → derivation replay (17-step) → provenance/candidate containment → framework-reference integrity (structural only: non-empty framework IDs, 64-hex manifest-hash format, non-empty assumption-hash strings, well-formed anomaly IDs — `session` never imports `mechanics`/`relativity`; manifest/anomaly existence is cross-checked externally in tests) → falsifiability/anomaly field validation (IDs non-empty when entries present; `Review.Challenge.Category` ∈ the eight §27.3 categories). Success = **artifact internal consistency only, never physical truth**.
- **Containment failures → `CandidateContainmentError`:** `Hypothesis` not `HYPOTHESIS`; hypothesis-dependent derivation output marked trusted (recomputed status ≠ retained, or corpus status ≠ `NONE`); promotion-like internal state claimed (§26.10).
- **Falsifiability structures:** `Prediction{ID, Observable, Relation(core.Expr), Assumptions}`, `FalsificationCondition{ID, TargetClaim, ContradictingCondition(core.Expr), Regime}`, `RecoveryClaim{ID, Description, FromFramework, Condition(core.Expr)}`, `AnomalyReference{ID, Framework, Description}` — all typed, deterministic canonical encoders, no dynamic maps; they record research-test structure without any truth/false decision (§25-MUST-01).
- **Review / Challenge (v2.3 reconciliation A):** data-only types defined **exclusively** in `core/corpus.go` (§27.1 `Challenge{StepID, Category, Severity, Description}`, §27.2 `Review{DerivationID, Challenges, ReviewerNotes}`); `session` **consumes `core.Challenge`/`core.Review` without redeclaring them and without any alias implying a second canonical declaration**; no `review/` package exists. Categories exactly `CategoryError, DimensionError, AssumptionConflict, ConventionConflict, UnsupportedIdentification, InvalidReduction, ProvenanceProblem, CandidateOverreach` (§27.3). `DraftMetadata.ReviewHistory []core.Review`.
- **External reviewer boundary:** the library defines schemas only — no reviewer selection, orchestration, queues, auto-challenge, or auto-correction (§27.4; deferral row §14). Intended loop: candidate JSON → external reviewer → Review JSON → human.
- **Candidate fields (canonical order §26.2):** `ID (= Draft DerivationID), Hypothesis, Premises, Assumptions, Derivation (ledger; final hash == LedgerHash), FrameworkDependencies, Predictions, FalsificationConditions, RecoveryClaims, AnomalyReferences, ReviewHistory, MRCVersion, LedgerHash`. `Hypothesis` can never be replaced by a trusted status.

## 11. Canonical vertical slices

Test-function inventory (names are normative for this plan; all A–S map to them):

| Slice (acceptance) | Test function | File |
|---|---|---|
| A typed mechanics | `TestTypedMechanicsConstructors` | `mechanics/relations_test.go` |
| B `F = ma` + manifest agreement | `TestNewtonSecondLawManifestMatch` | `mechanics/manifest_test.go` (+ `TestNewtonSecondLawConstruct` in `relations_test.go`) |
| C dimension rejection | `TestDimensionMismatch` | `ops/negative_test.go` |
| D kind rejection | `TestKindMismatchEqualDimensions` | `ops/negative_test.go` |
| E assumption conflict | `TestAssumptionConflict` | `ops/negative_test.go` |
| F convention conflict | `TestConventionConflict` | `ops/negative_test.go` |
| G differentiate `½mv²` | `TestDifferentiateKineticEnergy` | `mechanics/relations_test.go` |
| H energy-momentum relation | `TestEnergyMomentumRelation` | `relativity/derivation_test.go` |
| I `E=mc²` derivation | `TestMassEnergyDerivation` | `relativity/derivation_test.go` |
| J Lorentz limit = 1 | `TestLorentzFactorLimit` | `relativity/derivation_test.go` |
| K Identify firewall | `TestSimplifyNeverIdentifies` + `TestSessionIdentifyRecords` | `ops/negative_test.go`, `session/session_test.go` |
| L candidate containment | `TestHypothesisContamination` + `TestCandidateArtifactContainmentRejected` | `hypothesis/candidate_test.go` |
| M falsifiability | `TestSealedCandidatePreservesFalsifiability` | `hypothesis/candidate_test.go` |
| N anomaly | `TestCandidateReferencesManifestAnomaly` | `hypothesis/candidate_test.go` |
| O manifest | `TestMechanicsManifestCrossCheck`, `TestRelativityManifestCrossCheck` | both `manifest_test.go` |
| P exact rational round-trip | `TestExactRationalRoundTrip` | `core/expr_test.go` |
| Q ledger tamper/replay | `TestLedgerTamperDetection` | `session/session_test.go` |
| R determinism | `TestDerivationDeterminism` | `session/session_test.go` |
| S sealed handoff | `TestSealResearchCandidate` | `session/session_test.go` |
| (positional contract, §15.13.1) | `TestApplyPositionalInputs` | `ops/operations_test.go` |

Sequence and required content:

1. **A:** construct all nine mechanics wrappers only via fixed constructors; assert `CoreObject()` exact accessor, kind, dimension, symbol, `DEFINED` provenance, `ESTABLISHED` corpus status.
2. **B:** `NewtonSecondLaw()` expr canonical bytes equal manifest `canonical_expr` fixture; `F = m*a` relation shape; cross-check (slice O machinery) validates item vs constructor.
3. **C:** `Add` and `Compare` with mismatched dimensions → `DimensionMismatchError` via `errors.As`; no panic.
4. **D:** **equal-dimension, distinct named kinds — `Mass (M)` vs `RestMass (M)`** → attempt an MRC-003-incompatible operation (e.g. `Add`) → `CategoryMismatchError` via `errors.As`. The pair shares dimension `M`, so MRC-002 is never triggered; the test isolates MRC-003. Do not substitute a dimensionally different pair (e.g. `KineticEnergy` vs `Momentum`).
5. **E:** two valid objects, same `(Constraint,key)`, different structured values → merge → `AssumptionConflictError`.
6. **F:** conflicting `ConventionSet`s merged by an operation → `ConventionConflictError`.
7. **G:** `K := mechanics.NewKineticEnergy(m,v)`; `Differentiate(K.CoreObject(), v.CoreObject())` → canonical `m*v`, `Kind=Expression`, dimension Momentum; `Compare(result, NewMomentum().CoreObject(), eq)` accepted (Expression-vs-named, MRC-003 §6.2); exercises product rule + power rule + exact rational simplification.
8. **H:** `EnergyMomentumRelation()` → canonical expr `E^2 = (p*c)^2 + (m*c^2)^2` (byte fixture), dimension `Energy^2`, assumptions contain `rest_mass_nonnegative` + `speed_of_light_positive`, manifest match.
9. **I:** the exact §20 sequence with real operations and a session-less ops chain *and* the session-driven variant in `session_test` for S:
   `construct EnergyMomentumRelation()` → `construct ZeroThreeMomentum()` (carries `RestFrameAssumption`) → `Substitute(rel, ThreeMomentum(), ZeroThreeMomentum())` → `Simplify` → `Solve(…, Energy())` → `Compare(Energy(), ZeroEnergy(), gte)` → `SelectBranch(…, energy_nonnegative)`; assert each golden-trace state (step 2 `E²=(0*c)²+(m c²)²`, step 3 `E²=(m c²)²`, step 4 `BranchSet(E,[Sqrt(…),Neg(Sqrt(…))])`, final `m*c^2` canonical bytes). No hardcoded final result; the test source contains no `MassEnergyRelation()` call (secondary: source audit).
10. **J:** `Limit(LorentzFactor(), Velocity(), ZeroVelocity())` → exact `1/1`, via body expansion (source audit verifies no ID-match shortcut).
11. **K:** both sides — `Simplify(x)` never yields `IDENTIFIED` and appends no identification event; `Session.Identify(a,b,justification)` yields `IDENTIFIED`, records justification, records an `Identification` step visible after sealing; `HYPOTHESIS` operand → result `HYPOTHESIS` with attempt preserved.
12. **L:** create candidate concept → derive with it (Add/Simplify/Solve/Limit/Compare/SelectBranch/Identify) → every downstream status `HYPOTHESIS`; `TestNoForbiddenExportedAPI` (std `go/ast` export scan over the module) proves absence of `Promote, Trust, ApproveHypothesis, PromoteToEstablished, SetCorpusStatus, ApproveException, OverrideMRC, BypassMRC, Integrate, Series, Taylor, Simulate, RankTheories, ScoreTruth`; crafted artifact with hypothesis-dependent trusted output → `CandidateContainmentError`.
13. **M:** candidate with ≥1 `Prediction`, ≥1 `FalsificationCondition`; `Seal()`; sealed accessors return the same canonical values (round-trip through `CanonicalJSON`).
14. **N:** `AnomalyReference{ID:"no_gravity", …}`; external cross-check against `relativity/manifest.json` anomalies (test-side, not in `session`).
15. **O:** both manifests: `go:embed` bytes → `ValidateManifestBytes` → static constructor map resolves every item → all nine cross-check steps → corpus status preserved (`"established"`) → file bytes == `CanonicalManifestJSON`; item sets exactly 11 mechanics / 10 relativity (F1); negative: `TestManifestRejectsInconsistent` (unknown field, bad enum, non-canonical expr, duplicate ID) → `ManifestValidationError` (REQ-032-15).
16. **P:** expr → canonical JSON → expr byte-identical + equal hash; representative object round-trip through the internal kernel decoder (`TestObjectCanonicalRoundTrip`, `session_test.go`); `Simplify(Pow(Rational(0),2))→0`, `Simplify(Pow(Rational(1),-1))→1`, `Simplify(Pow(x,0))→1` for nonzero-safe `x` (F2), `0^0` → `UnsupportedOperationError`.
17. **Q:** all three §16.20 tamper forms + corrupted chain (REQ-032-11..14) detected; each returns `LedgerValidationError` (or wrapped typed error) from `Ledger.Validate` / candidate validation.
18. **R:** the full §20 derivation executed twice in-process + once through a fresh `Session` → byte-identical canonical artifacts, hashes, ledger JSON, candidate JSON (no hidden ambient state, REQ-015-01). Positional contract coverage for §15.13.1 lives in `TestApplyPositionalInputs` (all 12 IDs, wrong-order inputs rejected, `ops.Apply`/`Session.Step`/replay share indices).
19. **S:** full candidate flow: `New → Draft(id,label,metadata{Hypothesis, Premises, Predictions, FalsificationConditions, AnomalyReferences, FrameworkDependencies(manifest hashes computed from canonical manifest bytes), …}) → Define/Declare premises → Step×N (substitute, simplify, solve, compare(gte), select_branch) → Commit → Conclude(final) → Seal` → candidate validates, `LedgerHash == DerivationHash`, `ParseResearchCandidateJSON(candidateJSON).Validate()` round-trips to an equal trusted candidate; negative: `Seal` without hypothesis → `ProvenanceError` state `Concluded`; `Commit` empty → `LedgerValidationError`; post-seal `Step` → `ProvenanceError`; `Conclude` with wrong hash → error.

## 12. Coverage matrix

Conventions: `REQ-§N-MUST-nn` = deterministic coverage ID for a normative clause at specs_v2_3.md §N without an explicit `REQ-*` label (assigned per specs_v2_3.md §39; subsection-precise where the spec numbers subsections). "Secondary" = source/AST/import audit or §13 gate unless noted. Every labeled `REQ-*` of v2.3, every MRC ID, every unnumbered MUST/MUST NOT clause, and acceptance A–S appear. Rows were re-derived from the complete `specs_v2_3.md`; v2.2-only rows are gone.

### 12.1 Module, product boundary, scope (§0–§3)

| Requirement ID | Implementation location | Primary test | Secondary verification |
|---|---|---|---|
| REQ-000-01 | `go.mod` | build gate (`go build ./...`) | §13 G-Module: `head go.mod` |
| REQ-000-02 (MVP MUST NOT require third-party Go modules) | `go.mod` | `go test ./...` + `go list -deps ./...` | §13 G-Module/G-Imports: zero third-party modules regardless of `require` block formatting |
| REQ-000-03 | README §Authority; §1 architecture | S (`TestSealResearchCandidate`) | README human-authority section audit |
| REQ-000-04 | README "Constructor authority" (external/public API boundary, not hostile-source) | `core/object_test.go::TestInternalKernelBoundary` | README text audit (no adversarial-protection claim) |
| REQ-§0.1-MUST-01 (module = `github.com/PithomLabs/phys`) | `go.mod` | build gate | §13 G-Module |
| REQ-§0.1-MUST-02 (all repo-local imports use `github.com/PithomLabs/phys/...`) | every file import path | §13 G-Imports (AST import audit) | `go list ./...` |
| REQ-§0.1-MUST-03 (Go standard library only; no third-party deps) | `go.mod` (no third-party requirements) | §13 G-Imports (zero third-party in `go list -deps`) | — |
| REQ-§0.3-MUST-01 (MUST NOT claim adversarial write-access protection) | README threat-model wording | `TestInternalKernelBoundary` | README text audit |
| REQ-§1-MUST-01 (implementation demonstrates all six product requirements) | plan §1, §12.1 | §13 review + rows REQ-001-01..06 | — |
| REQ-001-01 | `core/`+`ops/` | A, C, D | MRC table §5 |
| REQ-001-02 | both manifests + domain packages | H, I, J, O | §34 gate (three mechanisms) |
| REQ-001-03 | `core/corpus.go` + manifests | O | — |
| REQ-001-04 | metadata types + `session` | A–S (matrix §12.9) | — |
| REQ-001-05 | no truth APIs anywhere | L, S + `hypothesis_test::TestNoForbiddenExportedAPI` (source scan for `TRUTH_SCORE` etc.) | §13 G-Audit |
| REQ-001-06 | `hypothesis/`, `session/` | L, M, N, S | — |
| REQ-§2-MUST-01 (MVP MUST NOT implement any listed non-goal) | tree + tests | `TestNoDeferredPackages` + `TestForbiddenSourceSurface` | §13 G-Scope |
| REQ-002-01 (no `.phys`/physics-specific language) | tree: no language pkg | `core/object_test::TestRepositoryTreeExact` | §13 G-Scope |
| REQ-002-02 (no lexer/parser/frontend) | no parser imports (non-test) | `core/object_test::TestForbiddenSourceSurface` (bans `go/parser`,`go/scanner`,`go/token`) | — |
| REQ-002-03 (no custom compiler language) | tree exactness; no language artifacts | `TestRepositoryTreeExact` | G-Scope + README non-goals |
| REQ-002-04 (no general mathematical ontology) | fixed 18 kinds + closed 10-node set | `core/expr_test::TestExprKindEnumExact`, `core/object_test::TestKindEnumExact` | G-Audit (no open-ended kind registry) |
| REQ-002-05 (no full CAS) | bounded ops | `ops/negative_test::TestSolveUnsupportedForm`, `TestDifferentiateUnsupportedForms` | — |
| REQ-002-06 (no theorem prover) | bounded entailment | `ops/operations_test::TestSignEntailmentBounded` (denied-inference case) | — |
| REQ-002-07 (no numerical execution) | no floats | `TestForbiddenSourceSurface` (bans `float64`/`float32`) | — |
| REQ-002-08 (no integration/optimization/simulation/Monte Carlo/GPU/ODE) | no numerics | `TestForbiddenSourceSurface` (bans `math/rand`, `os/exec`, `net/http`) + `TestNoDeferredPackages` | — |
| REQ-002-09 (no empirical-data ingestion) | no adapter pkgs | `core/object_test::TestNoDeferredPackages` | G-Scope |
| REQ-002-10 (no theory ranking/truth scoring) | no truth API | `hypothesis_test::TestNoForbiddenExportedAPI` (banned identifiers incl. `Score`,`Rank`) | — |
| REQ-002-11 (no corpus-status inference) | corpus status loaded only | `manifest_test::TestCorpusStatusPreserved` + `TestNoForbiddenExportedAPI` (`SetCorpusStatus`) | source audit |
| REQ-002-12 (no hypothesis promotion) | no promotion API | `TestNoForbiddenExportedAPI` | — |
| REQ-002-13 (no MRC bypass/override API) | no bypass API | `TestNoForbiddenExportedAPI` (`Override`,`Bypass`) + `session_test::TestSessionSurfaceExact` | — |
| REQ-002-14 (no MRC exception workflow) | no exception API | same as 002-13 | README MRC-fallibility note |
| REQ-002-15 (no `physvet`) | no pkg | `TestNoDeferredPackages` | G-Scope |
| REQ-002-16 (no electromagnetism pkg) | no pkg | `TestNoDeferredPackages` | G-Scope |
| REQ-002-17 (no qm pkg) | no pkg | `TestNoDeferredPackages` | G-Scope |
| REQ-002-18 (no qft pkg) | no pkg | `TestNoDeferredPackages` | G-Scope |
| REQ-002-19 (no statmech pkg) | no pkg | `TestNoDeferredPackages` | G-Scope |
| REQ-002-20 (no general relativity content) | no pkg | `TestNoDeferredPackages` + relativity manifest scope records | G-Scope |
| REQ-002-21 (no general diff-geometry/tensor algebra) | no tensor nodes/pkgs | `TestExprKindEnumExact` + `TestNoDeferredPackages` | — |
| REQ-002-22 (no symbolic integration) | no API | `TestNoForbiddenExportedAPI` (`Integrate`) | — |
| REQ-002-23 (no series expansion) | no API | `TestNoForbiddenExportedAPI` (`Series`,`Taylor`) | §14 deferral |
| REQ-002-24 (no multi-agent orchestration runtime) | no `review` pkg | `TestNoDeferredPackages` | — |
| REQ-002-25 (no EBP 2.1 integration) | no API/pkg | `TestNoForbiddenExportedAPI`/G-Scope | README non-goals |
| REQ-002-26 (no general theory-management platform) | tree | G-Scope (`TestRepositoryTreeExact`) | README |
| REQ-§3-MUST-01 (create/modify only the §3 files) | §2 tree | `TestRepositoryTreeExact` | G-Scope |
| REQ-003-01 | §2 tree | `TestRepositoryTreeExact` (exact 39-file set) | G-Scope |
| REQ-003-02 | test placement | `TestRepositoryTreeExact` (all tests adjacent `*_test.go`) | — |
| REQ-003-03 | `package core_test` authority tests | `core/object_test.go` (runs as external package) | compiler enforces package clause |

### 12.2 Dependency graph (§4)

| Requirement ID | Implementation location | Primary test | Secondary |
|---|---|---|---|
| REQ-004-01 (no dependency cycle) | package graph | `go test ./...` (compile) | `go list ./...` edge audit |
| REQ-004-02 | `core` imports | `core/object_test::TestCoreImportIndependence` (parses `core/*.go` imports; forbids `ops`,`session`) | G-Imports |
| REQ-004-03 | `internal/kernel` | `core/object_test::TestInternalKernelBoundary` (path contains `/internal/`; `core` export surface lacks `MintObject`) | Go toolchain enforces `internal/` rule |
| REQ-§4-MUST-01 (`core` MUST NOT import `ops`) | `core` imports | `TestCoreImportIndependence` | G-Imports |
| REQ-§4-MUST-02 (`core` MUST NOT import `session`) | `core` imports | `TestCoreImportIndependence` | G-Imports |
| REQ-§4-MUST-03 (`ops` MUST NOT import `session`) | `ops` imports | `session/session_test::TestOpsImportsNoSession` (parses `ops/*.go`) | G-Imports |

### 12.3 Object model, kinds, dimensions (§5–§7)

| Requirement ID | Implementation location | Primary test | Secondary |
|---|---|---|---|
| REQ-005-01 (`CoreObject()` accessor on every domain type) | domain wrappers | A (`TestTypedMechanicsConstructors`) + relativity ctor test H | compile-time signature check |
| REQ-005-02 (ops accept `core.Object`) | `ops` signatures | `ops/operations_test` (all 12 via `core.Object`) | G-Imports |
| REQ-005-03 (wrappers don't duplicate symbolic logic) | thin wrappers | A | G-Audit: domain export surface = constructors+accessors only |
| REQ-005-04 / REQ-005-10 | `core` surface | `core/object_test::TestNoGenericFactory` (AST: no `NewObject`, no `MintObject` in `core`) | — |
| REQ-005-05 (no domain ctor letting callers pick Kind/Dim/Provenance/corpus status) | domain ctor surfaces | `TestNoGenericFactory` extended to domain packages (AST: no exported fn taking `Kind`/`CorpusStatus` params) | — |
| REQ-005-06 | unexported fields | `core/object_test::TestKernelObjectFieldsUnexported` (reflect) | — |
| REQ-005-07 | zero invalid | `TestZeroObjectInvalid` | — |
| REQ-005-08 | no mutators | `TestObjectAccessorSurfaceExact` (AST method scan: only listed accessors, value receivers) | — |
| REQ-005-09 | accessor names/types | `TestObjectAccessorSurfaceExact` | A |
| REQ-§5.3-MUST-01 (returned metadata immutable from caller) | accessors return copies | `TestObjectDefensiveCopies` (mutate returned slice/`Rat`; object unchanged) | — |
| REQ-005-11 | mint validation | `TestMintObjectRejectsInvalidSpec` (bad expr, empty required name, bad dimension, bad provenance) | — |
| REQ-005-12 + REQ-§5.5-MUST-01 (public ops reject `!Valid()` first) | ops validity check first | `ops/negative_test::TestInvalidObjectRejected` (loops all 12 `Apply` IDs + `Identify`) | — |
| REQ-§5.2.1-MUST-01 (exactly one production minting entry point) | `internal/kernel/mint.go` | `TestNoGenericFactory` (AST: single `func MintObject` in module) | source audit |
| REQ-§5.2.1-MUST-02 (`MintObject` 7-point contract: validity, name, kind, dimension, assumptions, conventions, provenance) | `internal/kernel/mint.go` | `TestMintObjectRejectsInvalidSpec` (one failing case per contract point) | — |
| REQ-§5.2.1-MUST-03 (`MintObject` exists only under `internal/kernel`, not re-exported) | `core` export surface | `TestNoGenericFactory` (AST export list) | G-Imports |
| REQ-§5.4-MUST-01 (`core` must not forward mint) | `core/object.go` | `TestNoGenericFactory` | AST export list |
| REQ-006-01 | 18 Kind values | `core/object_test::TestKindEnumExact` (ordinal + string table) | — |
| REQ-006-02 | ctor-assigned immutable | `TestObjectAccessorSurfaceExact` + A | — |
| REQ-006-03 (compatibility table exactly) | `ops/arithmetic.go`, `ops/relation.go` | `ops/operations_test::TestKindCompatibilityTable` (table-driven over full §6.2 grid) | D |
| REQ-§6-MUST-01 (ordinals stable in `mrc-v0.4`) | `kernel/types.go` | `TestKindEnumExact` (hardcoded ordinal table) | — |
| REQ-§6.2-MUST-01 (inequality operands ordered scalars: Mass, RestMass, Time, Energy, KineticEnergy, SpeedOfLight, Expression) | `ops/relation.go` | `TestInequalityOrderedKinds` (Velocity rejected; each listed kind accepted) | G |
| REQ-§6.2-MUST-02 (SelectBranch input `BranchSet`, constraint `Relation`) | `ops/relation.go` | `TestKindCompatibilityTable` (SelectBranch rows) | I |
| REQ-§7.0-MUST-01 (fixed constructors for all MVP dimensions; kernel-side constructor + core delegation — F3) | `internal/kernel/types.go`, `core/dimension.go` | `core/expr_test::TestDimensionConstructorsAndAPI` (each §7.2 exponent vector) | compile |
| REQ-§7.0-MUST-02 (`Dimension` exposes the exact read-only API) | `core/dimension.go` | `TestDimensionConstructorsAndAPI` (AST method list) | — |
| REQ-007-01 + REQ-§7.1-MUST-01 (big.Rat exponents) | `kernel/types.go` | `TestRationalExactSerialization` (`"1/2"`, `-3/4`, `"0/1"`) | `TestForbiddenSourceSurface` (no float) |
| REQ-007-02 | structural equality | `TestDimensionEqual` | — |
| REQ-007-03 + REQ-§7.3-MUST-01/02 (deterministic ops; canonical reduced) | `kernel` dimension ops | `TestDimensionArithmetic` (Multiply/Divide/Pow incl. fractional `Pow`), `TestDimensionDeterminism` (double-run bytes) | — |
| REQ-§7.2-MUST-01 (required dimensions table of §7.2) | `core/dimension.go` constructors | `TestDimensionConstructorsAndAPI` (each constructor vs §7.2 exponents) | — |

### 12.4 Expressions, canonicalization, hashing (§8–§10)

| Requirement ID | Implementation location | Primary test | Secondary |
|---|---|---|---|
| REQ-008-01 | closed node set | `core/expr_test::TestExprKindEnumExact` (exactly 10 constants/ordinals) | — |
| REQ-008-02 | no callbacks in expr | `TestExprStoresNoCallbacks` (reflect field-type scan of `kernel.Expr`) | — |
| REQ-008-03 | symbols are data | `TestSymbolIsDataOnly` (symbol named `";x=1+1"` round-trips unchanged; no evaluation) | `TestForbiddenSourceSurface` |
| REQ-008-04 + REQ-§8.4-MUST-01 (exact rational serialization; no float64) | `TestExactRationalRoundTrip` (P) | — | `TestForbiddenSourceSurface` |
| REQ-008-05 + REQ-§8.5-MUST-01 (no map order; constructor-time canonicalization list) | structural canonicalization | `TestCanonicalOrderingStable` (REQ-032-17) + `TestConstructorCanonicalization` (flatten/combine/sign/sort) | — |
| REQ-008-06 (exact `Pow` exponent) | `kernel` | `TestPowExponentExact` (type is `*big.Rat`; `"2/1"` bytes) | `TestForbiddenSourceSurface` |
| REQ-§8.0-MUST-01 (public `core` exposes the fixed expr constructors/helpers) | `core/expr.go` | `core/expr_test::TestExprConstructorSurface` | compile |
| REQ-§8.0-MUST-02 (read-only inspection methods exist on `Expr`) | `core/expr.go` | `TestExprAccessors` | — |
| REQ-§8.0-MUST-03 (kind-check before accessors) | `ops` discipline | shared precondition + `negative_test` unsupported-form tests (REQ-032-18..20) | source audit |
| REQ-§8.0-MUST-04 (`RationalValue()` defensive copy) | `kernel` accessors | `TestExprDefensiveCopies` (mutate returned `Rat`; tree unchanged) | — |
| REQ-§8.0-MUST-05 (`Exponent()` defensive copy) | `kernel` accessors | `TestExprDefensiveCopies` | — |
| REQ-§8.0-MUST-06 (slice accessors return copies; incl. construction-side copies of `NewRational`/`NewPow` inputs — prompt §5) | `kernel` constructors/accessors | `TestExprDefensiveCopies` (mutate returned slice + caller-owned inputs; tree unchanged) | — |
| REQ-§8.1-MUST-01 (closed expression tree of §8.1) | `kernel` | `TestExprKindEnumExact` | — |
| REQ-§8.2-MUST-01 (stable node ordinals) | `kernel` | `TestExprKindEnumExact` (ordinal table) | — |
| REQ-§8.2.1-MUST-01 (`ExprKind` stable values exactly as listed) | `kernel` | `TestExprKindEnumExact` (name+ordinal table) | — |
| REQ-§8.2.1-MUST-02 (`RelationOperator` exactly `eq,neq,lt,lte,gt,gte`) | `kernel` | `TestRelationOperatorEnumExact` | — |
| REQ-§8.3-MUST-01 (symbol content never executable) | `kernel` symbol | `TestSymbolIsDataOnly` | `TestForbiddenSourceSurface` |
| REQ-§8.6-MUST-01 (`Neg(Neg(x))→x` canonicalizer rules) | constructor canonicalization | `TestConstructorCanonicalization` | — |
| REQ-§8.7-MUST-01 (`Pow(x,1)→x`; **`Pow(x,0)→1` only for nonzero-safe base — F2**; `0^0` unsupported) | `ops/simplify.go` | `ops/operations_test::TestSimplifyIdentityRules` (incl. `Pow(x,0)` gated case + `0^0` error) | — |
| REQ-§8.8-MUST-01 (no sign rewrite without entailment) | `ops/simplify.go` | `ops/operations_test::TestSqrtSignRewriteRequiresAssumption` | — |
| REQ-009-01 + REQ-§9.1-MUST-01/02 (`EqualExpr` structural, never display strings; `EqualObject` canonical incl. metadata) | `core/canonical.go` | `core/expr_test::TestEqualExprStructural`, `core/object_test::TestEqualObjectMetadata` | — |
| REQ-009-02 | `core/canonical.go` (`EqualObject` byte-equality of object canonical JSON) | `TestEqualObjectMetadata` (metadata-only difference ⇒ unequal) | — |
| REQ-009-03 + REQ-§9.2-MUST-01 (deterministic hashes; lowercase hex in canonical JSON hash fields) | `kernel` hash helper | `TestHashDeterministic` (fixed vectors, double-run bytes) | §31 determinism test (R) |
| REQ-009-04 + REQ-§9.3-MUST-01 (three-key child order; `bytes.Compare` byte ordering) | `kernel` child sorter | `TestCanonicalOrderingStable` (REQ-032-17) | — |
| REQ-009-05 (exact rational combination, §9.4) | `kernel` rational arithmetic | `TestExactRationalRoundTrip` (P) + `TestRationalCombinationExact` | — |
| REQ-§9.6-MUST-01 (exact mechanical identity/zero/rational-`Pow` rule list) | `ops/simplify.go` | `ops/operations_test::TestSimplifyIdentityRules` (each §9.6 rule incl. `0^0` error) | — |
| REQ-§9.7-MUST-01 (combinations required by MVP derivations: `x·x`, `Pow(x,2)·Pow(x,2)`, `Pow(Pow(x,a),b)` nonneg-int only) | `ops/simplify.go` | `TestSimplifyRepeatedPowers` | — |
| REQ-§9.7-MUST-02 (no general branch-sensitive power algebra) | `ops/simplify.go` | `ops/negative_test::TestUnsupportedPowerAlgebra` | source audit (no CAS) |
| REQ-§9.8-MUST-01 (§9.8 `Sqrt` rewrite list) | `ops/simplify.go` | `TestSqrtRewriteRules` (`Sqrt(0)`, `Sqrt(1)`, `Sqrt(Pow(x,2))` gated case) | — |
| REQ-009-06 (bounded structural entailment only) | `internal/kernel` `Entails*` helpers via `ops/simplify.go` | `ops/operations_test::TestSignEntailmentBounded` (denied-inference case) | — |
| REQ-§9.9-MUST-01 (`Simplify` MUST NOT mint `IDENTIFIED` from `Relation(eq,…)`; relation simplification preserves operator — 🟡-3) | `ops/simplify.go` | K: `ops/negative_test::TestSimplifyNeverIdentifies`; relation-both-sides case in I (`Substitute`→`Simplify` on relation) | — |
| REQ-010-01 (typed structs) | all canonical encoders | `TestCanonicalStructsTyped` (reflect: no `map[string]any` fields) | G-Audit |
| REQ-010-02 (no `map[string]any`) | all canonical encoders | `TestCanonicalStructsTyped` | G-Audit (source scan) |
| REQ-010-03 (fixed field order = struct declaration order) | encoders in `kernel`/`core` | `TestExactRationalRoundTrip`, `TestObjectCanonicalRoundTrip`, O | — |
| REQ-010-04 (no JSON numbers for rationals/exponents) | rational/dimension encoders | P, `TestRationalStringNotNumber` (JSON token type `string`) | — |
| REQ-010-05 (canonical rational strings `"1/2"`) | rational encoder | P (`"1/2"`, `-3/4`, `"0/1"`) | — |
| REQ-010-06 (no timestamps/random IDs/pointers/env/map order) | all encoders | `TestNoVolatileFields` (source scan: `time.Now`, `rand`, `%p`; F4 set-order scan) | G-Audit |
| REQ-010-07 (canonical round-trip tests for expr + representative objects) | `core/expr_test.go`, `session/session_test.go` | P + `TestObjectCanonicalRoundTrip` | §13 gate |
| REQ-§10.3-MUST-01 (dimension JSON field order) | `kernel` dimension encoder | `TestDimensionJSONFieldOrder` (literal field sequence) | — |
| REQ-§10.4-MUST-01 (object JSON field order, incl. `schema_version` first — serialization DTO only) | `kernel` object encoder | `TestObjectCanonicalRoundTrip` (literal field sequence) | — |
| REQ-§10.5-MUST-01 (byte-identical canonical form + identical hash on round-trip) | expr/object encode-decode | P + `TestObjectCanonicalRoundTrip` + `TestHashDeterministic` | — |

### 12.5 Assumptions, conventions, provenance (§11–§13)

| Requirement ID | Implementation location | Primary test | Secondary |
|---|---|---|---|
| REQ-§11.0-MUST-01 (`AssumptionSet` immutable, exact exposed surface) | `core/assumption.go` | `TestAssumptionSetSurface` | — |
| REQ-§11.0-MUST-02 (`Assumption` read-only `Kind`/`Key`/typed-Value accessors) | `kernel` metadata | `TestAssumptionAccessors` | — |
| REQ-§11.0-MUST-03 (`Values()` returns a copy) | `core/assumption.go` | `TestAssumptionValuesCopy` (mutate result; set unchanged) | — |
| REQ-§11.1-MUST-01 (exact seven assumption kinds + lowercase JSON strings of §11.1 — replaces all v2.2 kinds) | `kernel/types.go` kind enum, `core/assumption.go` | `TestAssumptionKindEnumExact` (7 Go values + 7 JSON strings) | G-Audit |
| REQ-§11.2-MUST-01 (no free-form string relationship encoding) | `Assumption` value union | `TestAssumptionStructuredValue` (E fixtures use `ExprValue`) | G-Audit |
| REQ-011-01 (structured expression values, never equation strings) | `core/assumption.go` | E + `TestAssumptionStructuredValue` | — |
| REQ-§11.3-MUST-01 (deterministic constructors per §11.3; kernel-side primitive + core delegation — F3) | `internal/kernel/types.go`, `core/assumption.go` | `TestAssumptionConstructorsDeterministic` (double-run bytes) | — |
| REQ-§11.3-MUST-02 (keys non-empty) | constructors | `TestAssumptionKeyNonEmpty` (empty key ⇒ error) | — |
| REQ-§11.4-MUST-01 (same `(Kind,Key)` different canonical value ⇒ conflict) | `AssumptionSet.Merge` in `internal/kernel` | E: `ops/negative_test::TestAssumptionConflict` | — |
| REQ-011-02 (deterministic set union + exact conflict detection; F4 canonical-bytes element order, no map iteration) | merge used by all `ops` | E + `TestAssumptionMergeUnion` (order-independence: differently-built sets → identical canonical bytes) | §31 determinism (R) |
| REQ-§11.6-MUST-01 (`RestMass()` carries `m >= 0`) | `mechanics`/`relativity` primitives | `TestSignAssumptionsCarried` (I/J assumptions check) | — |
| REQ-§11.6-MUST-02 (`SpeedOfLight()` carries `c > 0`) | `relativity/primitives.go` | `TestSignAssumptionsCarried` | — |
| REQ-§12.0-MUST-01 (`ConventionSet` exact surface) | `core/convention.go` | `TestConventionSetSurface` | — |
| REQ-§12.0-MUST-02 (`Values()` returns a copy) | `core/convention.go` | `TestConventionValuesCopy` | — |
| REQ-012-01 (conflicts rejected deterministically; F4 ordering) | `ConventionSet.Merge` in `internal/kernel` | F: `ops/negative_test::TestConventionConflict` | — |
| REQ-012-02 (no equation solely as convention string) | `Convention{Key,Value string}` type | `TestConventionStringsOnly` + G-Audit | — |
| REQ-§13.0-MUST-01 (`Provenance` read-only accessors) | `kernel` provenance | `TestProvenanceAccessors` | — |
| REQ-§13.0-MUST-02 (`ParentHashes()` returns a copy) | `kernel` provenance | `TestProvenanceParentHashesCopy` | — |
| REQ-§13.0.1-MUST-01 (read-only status values + one deterministic provenance constructor; kernel-side primitive + `core.NewProvenance` delegation — F3) | `internal/kernel/types.go`, `core/provenance.go` | `TestProvenanceConstructorSingle` (AST: exactly one constructor) | G-Audit |
| REQ-013-01 (`APPROXIMATED` in enum, never minted by MVP op) | `kernel` status enum | `TestProvenanceStatusEnumExact` + ops result-status sweep never yields `APPROXIMATED` | §39 coverage note |
| REQ-013-02 (exact deterministic status law) | `ops` provenance helper + `session` actions | `TestProvenanceStatusLaw` (table over all 12 ops × statuses; Identify; Postulate/Define/Declare) | K, L, S |
| REQ-013-03 (parents = canonical object hashes) | `ops`/`session` result build | `TestProvenanceParentHashes` (parents equal input `HashObject`s) | — |
| REQ-§13.4-MUST-01 (implementing agent MUST NOT infer/revise corpus status) | no corpus-status API | `TestCorpusStatusPreserved` + `TestNoForbiddenExportedAPI` (`SetCorpusStatus`) | G-Audit |
| REQ-013-04 (library never computes corpus status) | manifest load path | O: `TestCorpusStatusPreserved` | — |
| REQ-013-05 (no `PHYSICALLY_TRUE`/`TRUTH_SCORE`/`PROBABILITY_OF_TRUTH`/`BEST_THEORY` API) | whole tree | `hypothesis_test::TestNoForbiddenExportedAPI` (banned identifiers incl. `ScoreTruth`) | G-Audit |
| REQ-§13.5-MUST-01 (source/framework inheritance: preserve iff all inputs identical and non-empty, else empty) | `ops` provenance helper | `TestProvenanceInheritance` (mixed inputs ⇒ empty) | — |

### 12.6 MRC rules and operation contracts (§14–§15)

| Requirement ID | Implementation location | Primary test | Secondary |
|---|---|---|---|
| MRC-001 (+ REQ-§14.1-MUST-01: valid object originates only via §5.4 paths) | `internal/kernel/mint.go`, `core/object.go`, domain/hypothesis constructors | `TestNoGenericFactory`, `TestMintObjectRejectsInvalidSpec`, `TestZeroObjectInvalid`, `TestInternalKernelBoundary` | G-Audit (single `MintObject`) |
| MRC-002 | `ops` precondition pipeline | C: `ops/negative_test::TestDimensionMismatch` | — |
| MRC-003 | `ops` kind table (§6.2) | D: `TestKindMismatchEqualDimensions` (**fixture `Mass` vs `RestMass`**) + `TestKindCompatibilityTable` | G |
| MRC-004 | `AssumptionSet.Merge` (`internal/kernel`) | E: `TestAssumptionConflict` | — |
| MRC-005 | `ConventionSet.Merge` (`internal/kernel`) | F: `TestConventionConflict` | — |
| MRC-006 | `session/session.go::Identify` only | K: `TestSessionIdentifyRecords`, `TestIdentifyRequiresJustification` (REQ-032-08) | G-Audit (no `ops.Identify`) |
| MRC-007 (+ REQ-§14.7-MUST-01: all committed steps pass through Session) | `session` ledger authority | `TestLedgerConstructionViaSessionOnly`, Q, S | G-Audit |
| MRC-008 (+ REQ-§14.8-MUST-01/02: hypothesis-dependent result is `HYPOTHESIS`; artifact presenting it trusted ⇒ validation failure) | `ops` status law + `session` candidate validation | L: `TestHypothesisContamination`, `TestCandidateArtifactContainmentRejected` | S |
| REQ-015-01 (no global/ambient session state) | all `ops` functions | R: `TestDerivationDeterminism` (fresh-Session variant) | G-Audit (no package-level mutable state in `ops`) |
| REQ-§15-MUST-01 (common semantics: no mutation, validity first, MRC, merge, fresh result, provenance, canonical) | shared `ops` precondition helper | `TestOperationCommonSemantics` (inputs byte-unchanged after every op) | — |
| REQ-§15.5-MUST-01 (`Pow` clones/reads exponent; no retained caller pointer) | `ops/arithmetic.go` | `TestPowClonesExponent` (mutate caller `*big.Rat`; result unchanged) | — |
| REQ-§15.6-MUST-01 (`Simplify` MUST NOT: create `IDENTIFIED`, create identification relation, change corpus status, invoke `Session.Identify`, consult session state) | `ops/simplify.go` | K: `TestSimplifyNeverIdentifies` + `TestCorpusStatusPreserved` | G-Audit |
| REQ-§15.7-MUST-01 (`variable.Expr()` is a single `Symbol` node) | `ops/transform.go` | `TestSubstituteContract` (non-symbol variable ⇒ error) | — |
| REQ-§15.7-MUST-02 (replacement dimension must equal variable dimension) | `ops/transform.go` | `TestSubstituteDimensionMismatch` | — |
| REQ-§15.7-MUST-03 (variable and replacement kind must be identical) | `ops/transform.go` | `TestSubstituteContract` (kind mismatch ⇒ error) | — |
| REQ-§15.7-MUST-04 (all occurrences of the exact variable symbol replaced) | `ops/transform.go` | `TestSubstituteContract` (repeated symbol fully replaced) | G |
| REQ-§15.8-MUST-01 (`wrt.Expr()` single `Symbol`) | `ops/transform.go` | `TestDifferentiateRequiresSingleSymbol` | — |
| REQ-§15.8-MUST-02 (result simplified after differentiation) | `ops/transform.go` | G: `TestDifferentiateKineticEnergy` | — |
| REQ-§15.9-MUST-01 (`Call(lorentz_factor)` expands §15.9.1 body first) | `ops/transform.go` | J: `TestLorentzFactorLimit` | source audit |
| REQ-§15.9.1-MUST-01 (exact variable/value forms accepted by `Limit`) | `ops/transform.go` | J + `TestLimitForms` | — |
| REQ-§15.9.1-MUST-02 (MUST NOT return `1` on function-ID match alone) | `ops/transform.go` | `TestLorentzFactorLimit` (traverses body) | G-Audit: source audit of `ops/transform.go` (no ID shortcut) |
| REQ-§15.11-MUST-01 (Solve target = single `Symbol` matching squared symbol) | `ops/relation.go` | I step 4 + `ops/negative_test::TestSolveUnsupportedForm` (REQ-032-18) | — |
| REQ-§15.12-MUST-01 (SelectBranch exact six-step §15.12 sequence; assumption key `selected_branch/<hash>`; no invented sign-entailment error — 🟡-4) | `ops/relation.go` | I step 6 + `TestSelectBranchContract` | — |
| REQ-§15.12-MUST-02 (mass-energy positive branch simplifies `Sqrt(Pow(m·c²,2)) → m·c²`) | `ops/relation.go` | I final assertion | — |
| REQ-§15.13-MUST-01 (one fixed dispatch mechanism for replay) | `ops/dispatch.go` | `TestApplyDispatchMatchesSessionReplay` (R) | — |
| REQ-§15.13-MUST-02 (`OperationParams` exact field order, canonical JSON, permitted `Kind` ∈ `empty\|pow\|compare\|identify`, per-kind field binding) | `ops/dispatch.go`, `session/session.go` | `TestParamsCanonicalRoundTrip` + `TestSessionSurfaceExact` | R |
| REQ-§15.13.1-MUST-01 (`ops.Apply`, `Session.Step`, and replay MUST adhere to the §15.13.1 positional schema — plan §9 table) | `ops/dispatch.go`, `session/session.go`, `session/ledger.go` | **`TestApplyPositionalInputs`** (all 12 IDs at exact indices) | R (`TestApplyDispatchMatchesSessionReplay`) |
| REQ-§15.13.1-MUST-02 (inputs MUST NOT be reordered/packed; unused positions omitted; replay reconstitutes by index) | `ops/dispatch.go`, replay decoder | `TestApplyPositionalInputs` (wrong-order inputs rejected; replay index reconstitution) | G-Audit |
| REQ-§15.13.1-MUST-03 (dispatch switch closed over exactly the 12 IDs; `identify` excluded; no registry) | `ops/dispatch.go` | `TestDispatchClosed` (unknown ID ⇒ error; `identify` ⇒ rejected) | G-Audit |

### 12.7 Session and derivation ledger (§16)

| Requirement ID | Implementation location | Primary test | Secondary |
|---|---|---|---|
| REQ-016-01 (state machine matches lifecycle `New→Drafting→Committed→Concluded→Sealed`) | `session/session.go` | `TestSessionStateMachine` (every legal + illegal transition) | REQ-032-21 |
| REQ-§16.1-MUST-01 (exact nine actions: Postulate, Declare, Define, Step, Identify, Conclude, Draft, Commit, Seal) | `session/session.go` | `TestSessionSurfaceExact` (AST method list) | — |
| REQ-§16.2-MUST-01 (empty `Commit` ⇒ `LedgerValidationError`, state stays `Drafting`) | `session/session.go` | `TestCommitEmptyDraftFails` | S |
| REQ-§16.2-MUST-02 (`Seal` without valid `HYPOTHESIS` ⇒ `ProvenanceError`, state stays `Concluded`) | `session/session.go` | `TestSealRequiresHypothesis` | S |
| REQ-§16.2.1-MUST-01 (`Hypothesis` valid + `HYPOTHESIS`; `Seal` rejects otherwise; all slices copied at Draft) | `DraftMetadata` validation | `TestSealRequiresHypothesis` | — |
| REQ-§16.2.1-MUST-02 (no separate candidate-metadata mutation methods) | `session` public surface | `TestSessionSurfaceExact` (AST method list) | G-Audit |
| REQ-§16.3-MUST-01 (exact method signatures of §16.3 incl. read-only `Validate`/`CanonicalJSON`) | `session/session.go` | `TestSessionSurfaceExact` | — |
| REQ-§16.4-MUST-01 (`DerivationID` non-empty, caller-supplied, never auto-generated) | `Session.Draft` | `TestDraftDerivationID` (empty ⇒ error; value preserved verbatim) | — |
| REQ-§16.5-MUST-01 (assertion actions: Postulate/Define status requirements, Declare preserves; recorded without `ops.Apply`; replay checks retained object) | `session/session.go` | `TestAssertionStepReplay` (assertion-step replay consistency + wrong-status rejections) | K |
| REQ-§16.6-MUST-01 (Step: state→inputs→reject `identify`→`ops.Apply` at §15.13.1 indices→canonical record→buffer→result; no caller-supplied output) | `session/session.go` | `TestStepRecordsApplyResult` + `TestStepRejectsIdentifyOperation` + `TestApplyPositionalInputs` | MRC-007 |
| REQ-§16.7-MUST-01 (Identify contract incl. trimmed justification, `IDENTIFIED`/`HYPOTHESIS`, justification in provenance, identification step) | `session.Identify` | K: `TestSessionIdentifyRecords`, `TestIdentifyRequiresJustification` | — |
| REQ-§16.7-MUST-02 (no global active-session state) | `session` | `TestDerivationDeterminism` (fresh instance) | G-Audit |
| REQ-§16.8-MUST-01 (draft entries fully specified, externally immutable) | draft buffer (unexported) | `TestSessionSurfaceExact` (no draft-accessor mutation path) | — |
| REQ-§16.9-MUST-01 (Commit: index from 1, `step-000001` IDs, genesis, chain hashes, retained canonicals, immutable ledger) | `session/ledger.go` | `TestCommitLedgerShape` (StepID/genesis/chain assertions) | — |
| REQ-§16.9-MUST-02 (draft buffer empty after commit) | `session` | `TestCommitClearsDraft` | — |
| REQ-§16.10-MUST-01 (Conclude object: valid; hash == last step output hash; becomes conclusion) | `Session.Conclude` | `TestConcludeMatchesFinalOutput` (wrong hash ⇒ error) | S |
| REQ-§16.11-MUST-01 (full validation before sealing) | `Session.Seal` | `TestSealValidatesFirst` (corrupt-then-seal ⇒ error) | Q |
| REQ-§16.12-MUST-01 (exact 16-field Step order of §16.12) | `session/ledger.go` | `TestStepFieldOrder` (literal JSON key sequence) | — |
| REQ-§16.13-MUST-01 (each step retains canonical inputs + params + output) | `session/ledger.go` | `TestReplayUsesRetainedCanonicals` | — |
| REQ-§16.13-MUST-02 (`HashObject(decoded input) == InputHash`) | `Session.Validate` | Q + `TestReplayUsesRetainedCanonicals` | — |
| REQ-§16.13-MUST-03 (`HashObject(decoded output) == OutputHash`) | `Session.Validate` | Q | — |
| REQ-§16.14-MUST-01 (`ParamsCanonical` encodes params for deterministic replay) | `session/ledger.go` | `TestParamsCanonicalRoundTrip` | R |
| REQ-§16.15-MUST-01 (genesis `PreviousStepHash` = 64 lowercase zeros) | `session/ledger.go` | `TestCommitLedgerShape` | — |
| REQ-§16.16-MUST-01 (StepID deterministic from Index; no UUIDs) | `session/ledger.go` | `TestCommitLedgerShape` | `TestNoVolatileFields` |
| REQ-§16.17-MUST-01 (`StepEnvelope` hash formula; no self-reference) | `session/ledger.go` | `TestStepEnvelopeHashFormula` | — |
| REQ-§16.18-MUST-01 (derivation hash = last step's `CurrentStepHash`) | `session/ledger.go` | `TestDerivationHashEqualsLastStep` | S |
| REQ-§16.19-MUST-01 (exact 17-step `Validate` order, incl. replay at §15.13.1 indices) | `Session.Validate` | `TestValidatePipelineOrder` (fault-injection per step) | Q |
| REQ-§16.19-MUST-02 (mismatch ⇒ `LedgerValidationError` or wrapped typed error) | `Session.Validate` | Q: `TestLedgerTamperDetection` | REQ-032-11..14 |
| REQ-§16.20-MUST-01 (changing only `OutputCanonical` ⇒ hash-mismatch detection) | `Ledger.Validate` | Q: `TestLedgerTamperDetection` mutation 1 | REQ-032-12 |
| REQ-§16.20-MUST-02 (changing `OutputCanonical` + recomputing only `OutputHash` ⇒ current-step hash mismatch) | `Ledger.Validate` | Q: `TestLedgerTamperDetection` mutation 2 | REQ-032-13 |
| REQ-§16.20-MUST-03 (changing retained output + recomputing all hashes ⇒ replay divergence detection) | `Ledger.Validate` | Q: `TestLedgerTamperDetection` mutation 3 | REQ-032-14 |
| REQ-§16.21-MUST-01 (`Ledger` exact value surface) | `session/ledger.go` | `TestLedgerSurfaceExact` (AST) | — |
| REQ-§16.21-MUST-02 (`Step` unexported storage, read-only accessors for all §16.12 fields) | `session/ledger.go` | `TestStepAccessorsExact` (AST + reflect) | — |

### 12.8 Corpus manifests and domain packages (§17–§19)

| Requirement ID | Implementation location | Primary test | Secondary |
|---|---|---|---|
| REQ-§17.0-MUST-01 (typed manifest structures of §17.0) | `core/corpus.go` | `TestManifestStructSurface` (AST field list/order) | — |
| REQ-§17.0-MUST-02 (`ParseManifest`, `ValidateManifestBytes`, `CanonicalManifestJSON`, `Manifest.Hash`) | `core/corpus.go` | O: `TestMechanicsManifestCrossCheck`, `TestRelativityManifestCrossCheck` | — |
| REQ-§17.6-MUST-01 (`statement` never parsed as mathematics) | `core/corpus.go` | `TestStatementNotParsed` | G-Audit (no parser imports) |
| REQ-§17.7-MUST-01 (loading: typed decode, closed node set, field/enum validation, canonical structure) | `core/corpus.go` | O + `TestManifestRejectsInconsistent` (REQ-032-15) | — |
| REQ-§17.7-MUST-02 (`ValidateManifestBytes` exists in `core/corpus.go`; no constructors, no I/O) | `core/corpus.go` | `TestValidateManifestBytesPure` (source scan: no ctor calls/`os.ReadFile`) | G-Audit |
| REQ-§17.8-MUST-01 (static test-only `map[string]func() core.Object`) | both `manifest_test.go` | O (map present; no reflection) | G-Audit |
| REQ-§17.8-MUST-02 (map resolves every `constructor` id) | both `manifest_test.go` | O (unresolved id ⇒ fail) | — |
| REQ-§17.8-MUST-03 (nine-step cross-check incl. six comparisons) | both `manifest_test.go` | O (nine steps asserted per item) | — |
| REQ-§17.9-MUST-01 (manifest files contain explicit corpus status values) | both `manifest.json` | O: `TestCorpusStatusPreserved` (field present per item) | — |
| REQ-§17.9-MUST-02 (implementation loads values, never derives them) | `core/corpus.go` loader | O (mismatch fixture ⇒ preserved verbatim) | G-Audit |
| REQ-§17.9-MUST-03 (tests verify loading preserves declared status) | both `manifest_test.go` | O | — |
| REQ-§18.0-MUST-01 (exact mechanics constructor signatures: 8 primitives + `NewKineticEnergy` + 3 relations) | `mechanics/primitives.go`, `mechanics/relations.go` | A: `TestTypedMechanicsConstructors` (compile-time signatures) | — |
| REQ-§18.1-MUST-01 (nine thin nominal wrappers with `CoreObject()`) | `mechanics/primitives.go` | A | G-Audit (wrapper-only exports) |
| REQ-§18.1-MUST-02 (exact constructor naming `NewMass`…`NewEnergy`) | `mechanics/primitives.go` | A + O | — |
| REQ-§18.2-MUST-01 (three zero-argument relation constructors) | `mechanics/relations.go` | B: `TestNewtonSecondLawConstruct` | — |
| REQ-§18.2-MUST-02 (`NewKineticEnergy(m,v)` exact expression) | `mechanics/relations.go` | G (input to derivative) | — |
| REQ-§18.3-MUST-01 (manifest records classical scope + limitation/anomaly record) | `mechanics/manifest.json` | O + N machinery | — |
| REQ-§19.0-MUST-01 (exact relativity constructor signatures: 8 wrappers incl. `Velocity`, 3 fixed, 3 zeros, `RestFrameAssumption`) | `relativity/primitives.go`, `relativity/relations.go` | H: `TestEnergyMomentumRelation` (signatures) | — |
| REQ-§19.1-MUST-01 (thin nominal wrappers) | `relativity/primitives.go` | H | G-Audit |
| REQ-§19.1-MUST-02 (fixed relation/function constructors: `LorentzFactor`, `EnergyMomentumRelation`, `MassEnergyRelation`) | `relativity/relations.go` | H, I, J | — |
| REQ-§19.3-MUST-01 (manifest identifies the four framework assumptions + `metric.signature = -+++` convention) | `relativity/manifest.json` | H (assumptions check) + O | — |
| REQ-§19.5-MUST-01 (`EnergyMomentumRelation()` exact form) | `relativity/relations.go` | H (byte fixture) | — |
| REQ-§19.5-MUST-02 (carries `m >= 0` and `c > 0` through constituents) | `relativity/relations.go` | H + `TestSignAssumptionsCarried` | — |
| REQ-§19.6-MUST-01 (`MassEnergyRelation` manifest provenance `DERIVED` + metadata) | `relativity/manifest.json` | O | — |
| REQ-§19.6-MUST-02 (vertical derivation test MUST NOT call `MassEnergyRelation()`) | `relativity/derivation_test.go` | I | G-Audit (source scan of test) |
| REQ-§19.7-MUST-01 (zero constructors exposed with exact values) | `relativity/relations.go` | I (steps 2, 5) | — |
| REQ-§19.9-MUST-01 (manifest contains structured scope limitation/anomaly record `no_gravity`) | `relativity/manifest.json` | N: `TestCandidateReferencesManifestAnomaly` | — |

### 12.9 Canonical vertical slices (§20–§25)

| Requirement ID | Implementation location | Primary test | Secondary |
|---|---|---|---|
| REQ-§20.1-MUST-01 (six-step §20 sequence executed with real operations) | slice I test | I: `TestMassEnergyDerivation` (each golden-trace state asserted) | — |
| REQ-§20.1-MUST-02 (test MUST NOT call `MassEnergyRelation()`) | slice I test | I | G-Audit (test source scan) |
| REQ-§21-MUST-01 (derivative nontrivial; exercises Expression-vs-named comparison) | slice G test | G: `TestDifferentiateKineticEnergy` | — |
| REQ-§21-MUST-02 (comparison accepted under MRC-003) | `ops/relation.go` + G test | G (`Compare(result, Momentum(), eq)` accepted) | D table |
| REQ-§22-MUST-01 (limit result from fixed-body expansion + substitution + simplify) | slice J test | J: `TestLorentzFactorLimit` | G-Audit |
| REQ-§23-MUST-01 (firewall test establishes both sides) | slice K tests | K: `TestSimplifyNeverIdentifies` + `TestSessionIdentifyRecords` | — |
| REQ-§23-MUST-02 (Simplify path: no `IDENTIFIED`, no identification event) | `ops/simplify.go` | K | — |
| REQ-§23-MUST-03 (Identify path: full recording contract) | `session.Identify` | K | — |
| REQ-024-01 (candidate construction forces `HYPOTHESIS`) | `hypothesis/candidate.go` | L: `TestHypothesisContamination` | — |
| REQ-024-02 (explicit `Kind` and `Dimension` required) | `hypothesis.NewCandidateConcept` | `TestCandidateConceptRequiresKindDimension` | — |
| REQ-§24.3-MUST-01 (containment test: downstream statuses all `HYPOTHESIS`) | slice L test | L (op sweep + Identify) | — |
| REQ-§24.4-MUST-01 (no promotion-equivalent exported names anywhere) | whole tree | L: `TestNoForbiddenExportedAPI` (banned list) | G-Audit |
| REQ-§25-MUST-01 (structure without true/false decision) | `session` falsifiability types | M: `TestSealedCandidatePreservesFalsifiability` | G-Audit (no truth fields) |
| REQ-§25.5-MUST-01 (provisional candidate constructed with ≥1 prediction + ≥1 falsification condition) | slice M test | M | — |
| REQ-§25.5-MUST-02 (sealed artifact preserves them byte-identically) | `Session.Seal` | M (accessor + `CanonicalJSON` round-trip) | — |

### 12.10 ResearchCandidate and review artifacts (§26–§27)

| Requirement ID | Implementation location | Primary test | Secondary |
|---|---|---|---|
| REQ-§26.0-MUST-01 (five support types defined in `session`; **`Challenge`/`Review` defined exclusively in `core/corpus.go`, consumed by `session` without redefinition and without aliases — v2.3 reconciliation A**) | `session/research_candidate.go`, `core/corpus.go` | `TestSupportTypeSurface` (AST: `core` declares Challenge/Review; `session` declares only the five support types; no session alias) | G-Audit (no `review/` package) |
| REQ-§26.0-MUST-02 (deterministic canonical encoders; no dynamic maps) | support types | `TestSupportTypeCanonicalDeterministic` (double-run) | REQ-010-02 |
| REQ-§26.2-MUST-01 (candidate canonical field order of §26.2) | `session/research_candidate.go` | `TestCandidateFieldOrder` (literal JSON key sequence) | — |
| REQ-§26.3-MUST-01 (`Hypothesis` is candidate concept with `HYPOTHESIS`) | `DraftMetadata` + `Seal` | L, S | — |
| REQ-§26.3-MUST-02 (artifact never replaces it with trusted status) | candidate validation | `TestCandidateArtifactContainmentRejected` | REQ-032-09 |
| REQ-§26.6-MUST-01 (derivation final hash == `LedgerHash`) | `Seal` | S: `TestSealResearchCandidate` | — |
| REQ-§26.8-MUST-01 (validation pipeline, exact order) | `ResearchCandidate.Validate` | `TestValidatePipelineOrder` (fault-injection per stage) | S |
| REQ-§26.8-MUST-02 (framework-reference integrity = structural validity only inside `session`) | `ResearchCandidate.Validate` | `TestFrameworkReferenceStructuralOnly` (malformed hash rejected; unknown-but-well-formed accepted in `session`, cross-checked externally) | G-Imports |
| REQ-§26.8-MUST-03 (inside `session`: framework-reference integrity is deterministic structural validity only — non-empty framework IDs, valid manifest-hash format, required non-empty assumption-hash strings, well-formed anomaly IDs; `session` MUST NOT import `mechanics` or `relativity`; existence cross-checked externally in tests) | `session` imports + `ResearchCandidate.Validate` | `TestSessionImportIndependence` (AST import scan) + `TestFrameworkReferenceStructuralOnly` | G-Imports |
| REQ-§26.9-MUST-01 (external JSON never decoded into trusted value; unverified wrapper) | `ParseResearchCandidateJSON` | REQ-032-10a: `TestParseResearchCandidateJSONUnverifiedOnly` | G-Audit |
| REQ-§26.10-MUST-01 (unexported storage; no arbitrary public constructor) | `session` | `TestCandidateSurfaceExact` (AST: no exported struct literal/ctor) | — |
| REQ-§26.10-MUST-02 (`UnverifiedResearchCandidate.Validate` delegates to same private seal path; no independent path) | `session/research_candidate.go` | `TestUnverifiedValidateDelegatesToSeal` | G-Audit (one private ctor) |
| REQ-§26.10-MUST-03 (exact 14 accessors + `Validate` + `CanonicalJSON`) | `session` | `TestCandidateSurfaceExact` | — |
| REQ-§26.10-MUST-04 (collection accessors return copies) | `session` | `TestCandidateAccessorsCopy` | — |
| REQ-§26.10-MUST-05 (exact unverified surface: type + 3 funcs) | `session` | `TestUnverifiedSurfaceExact` (AST) | — |
| REQ-§26.10-MUST-06 (`ParseResearchCandidateJSON` MUST NOT return `ResearchCandidate` and MUST NOT expose `core.Object` accessors; outer-schema validation + stored canonical bytes only) | `session/research_candidate.go` | REQ-032-10a: `TestParseResearchCandidateJSONUnverifiedOnly` | G-Audit |
| REQ-§26.10-MUST-07 (`UnverifiedResearchCandidate.Validate()` is the only public path from external JSON to a trusted candidate; runs the full §26.8 pipeline before success) | `session/research_candidate.go` | REQ-032-10a + `TestValidatePipelineOrder` | G-Audit (single public entry) |
| REQ-§26.10-MUST-08 (artifact containment: `Hypothesis` not `HYPOTHESIS` / hypothesis-dependent derivation output marked trusted / promotion-like state claimed ⇒ `CandidateContainmentError`) | `session` candidate validation | L: `TestCandidateArtifactContainmentRejected` | REQ-032-09 |
| REQ-§26.10-MUST-09 (candidate validation MUST NOT treat self-recomputed hashes as authenticity proof; internal consistency only) | `ResearchCandidate.Validate` | `TestValidateDoesNotTrustRecomputedHashes` (rewrite all hashes ⇒ still consistency-only) | README integrity-vs-authenticity audit |
| REQ-§27-MUST-01 (review artifact is data only; `Challenge`/`Review` defined in `core/corpus.go`; no `review/` package in MVP — v2.3 reconciliation A) | `core/corpus.go` | `TestSupportTypeSurface` (AST: types declared only in `core`; no `review/` package in tree) | G-Audit, G-Scope |
| REQ-§27.1-MUST-01 (`Challenge` exact fields: `StepID, Category, Severity, Description`) | `core/corpus.go` | `TestChallengeReviewFieldOrder` (AST field order) | — |
| REQ-§27.2-MUST-01 (`Review` exact fields: `DerivationID, Challenges, ReviewerNotes`; `session` consumes `core.Review` without redefinition or alias) | `core/corpus.go`, `session/research_candidate.go` | `TestChallengeReviewFieldOrder` + `TestSupportTypeSurface` | G-Audit (no `session` alias) |
| REQ-§27.3-MUST-01 (exact eight categories: `CategoryError, DimensionError, AssumptionConflict, ConventionConflict, UnsupportedIdentification, InvalidReduction, ProvenanceProblem, CandidateOverreach`) | `core/corpus.go` | `TestReviewCategoryEnumExact` (ordinal + string table) | — |
| REQ-§27.4-MUST-01 (orchestration boundary: library implements no reviewer selection, agent orchestration, review queues, auto-challenge, auto-correction) | no orchestration code exists | `TestNoForbiddenExportedAPI` + `TestNoDeferredPackages` | G-Scope; deferral row §14 |

### 12.11 Manifest data, MRC fallibility, errors, determinism, negative tests, acceptance (§28–§33)

| Requirement ID | Implementation location | Primary test | Secondary |
|---|---|---|---|
| REQ-§28.1-MUST-01 (mechanics manifest entries for at least the 11 items: `Mass, Time, Position, Velocity, Acceleration, Force, Momentum, Energy, NewtonSecondLaw, MomentumRelation, KineticEnergyRelation`; `KineticEnergy` parameterized — represented by `KineticEnergyRelation`, not a zero-arg item) | `mechanics/manifest.json` | O: `TestMechanicsManifestCrossCheck` (exact item-ID set assertion) | — |
| REQ-§28.2-MUST-01 (relativity manifest entries for at least the 10 items: `Spacetime, MinkowskiMetric, RestMass, Energy, ThreeMomentum, FourMomentum, SpeedOfLight, LorentzFactor, EnergyMomentumRelation, MassEnergyRelation`; `Velocity` is a wrapper, **not** a manifest item; `RestFrame` is an assumption — F1) | `relativity/manifest.json` | O: `TestRelativityManifestCrossCheck` (exact item-ID set assertion; `Velocity` absent) | §12.8 rows, G-Audit |
| REQ-§28.3-MUST-01 (each populated framework manifest includes ≥1 structured limitation/anomaly record) | both `manifest.json` files | O + §34.3 row | — |
| REQ-§28.3-MUST-02 (candidate test demonstrates referencing ≥1 such record) | slice N test | N: `TestCandidateReferencesManifestAnomaly` | — |
| REQ-§29-MUST-01 (`MRCVersion = mrc-v0.4` recorded; every rule ID stable) | `core` constants + `session` records (steps, candidate) | `TestMRCVersionRecorded` (constant `"mrc-v0.4"` in every Step/candidate; rule-ID table) | — |
| REQ-§29-MUST-02 (no runtime MRC override; no runtime exception registry) | no such API exists | `TestNoForbiddenExportedAPI` (`Override`, `Bypass`, `Exception`) | G-Audit; §14 deferral 20 |
| REQ-§29-MUST-03 (MVP implements only version recording, stable rule identifiers, deterministic rule enforcement — MRC fallibility stated in README) | `core`/`ops`/`session` | `TestMRCVersionRecorded` + G-Audit | README audit |
| REQ-§30-MUST-01 (exact 11 typed errors: `DimensionMismatchError, CategoryMismatchError, AssumptionConflictError, ConventionConflictError, IdentifyError, ProvenanceError, CandidateContainmentError, InvalidObjectError, UnsupportedOperationError, ManifestValidationError, LedgerValidationError`) | `core/errors.go` | `TestErrorTypesExact` (AST: exactly these 11 exported error types, no more) | — |
| REQ-§30-MUST-02 (all implement Go's `error` interface) | `core/errors.go` | `TestErrorTypesImplementError` (compile-time assertions + loop) | — |
| REQ-§30-MUST-03 (all inspectable with `errors.As`) | `core/errors.go` | `TestErrorsAsInspectable` (every negative test uses `errors.As`) | C–F, Q |
| REQ-§30-MUST-04 (normal MRC violations MUST NOT panic) | `ops`, `session` | `TestNoPanicOnMRCViolation` (table; any `recover` ⇒ test failure) | — |
| REQ-§30-MUST-05 (errors MUST NOT contain nondeterministic values; SHOULD carry deterministic diagnostics: operation ID, offending kinds/dimensions, stable rule ID) | error construction sites | `TestErrorDiagnosticsDeterministic` (double-run byte equality of `Error()` strings) | — |
| REQ-§31-MUST-01 (two identical derivation runs ⇒ byte-identical canonical artifacts) | whole pipeline | R: `TestDerivationDeterminism` | §13 gate 7 |
| REQ-§31-MUST-02 (the 12 listed artifact classes deterministic: expression/object/assumption/convention/manifest/step/step-chain/ledger/candidate canonical forms + hashes) | all encoders | `TestNoVolatileFields` + R | G-Audit |
| REQ-§31-MUST-03 (no artifact depends on timestamps, random UUIDs, pointer addresses, map iteration order, environment-dependent path ordering, locale-sensitive ordering, process IDs, hostnames) | all encoders | `TestNoVolatileFields` (source audit bans `time.Now`, `uuid`, `os.Hostname`, `os.Getpid`; runtime byte-compare) | G-Audit |
| REQ-§31.1-MUST-01 (same canonical derivation executed twice in one test process; all canonical outputs compared) | slice R test | R: `TestDerivationDeterminism` (run 1 vs run 2) | — |
| REQ-§31.1-MUST-02 (SHOULD additionally run through a fresh `Session` instance — implemented: run 3 uses `session.New` to prove no hidden mutable global state; REQ-015-01) | slice R test | R: `TestDerivationDeterminism` (run 3) | §12.7 REQ-§16.7-MUST-02 |
| REQ-§32-MUST-01 (tests exist for all listed failure classes) | all `*_test.go` | §13 gate 6 (inventory check) | — |
| REQ-032-01 dimension mismatch | `ops` | `TestDimensionMismatch` (C) | — |
| REQ-032-02 kind mismatch with equal dimensions (**fixture `Mass` vs `RestMass`** — F5) | `ops` | `TestKindMismatchEqualDimensions` (D) | — |
| REQ-032-03 assumption conflict | `AssumptionSet.Merge` via `ops` | `TestAssumptionConflict` (E) | — |
| REQ-032-04 convention conflict | `ConventionSet.Merge` via `ops` | `TestConventionConflict` (F) | — |
| REQ-032-05 invalid zero `core.Object` | `ops`, `session` | `TestInvalidObjectRejected` (all 12 `Apply` IDs + `Identify`) | — |
| REQ-032-06 absence of public generic object factory | `core` | `TestNoGenericFactory` (AST) | G-Audit |
| REQ-032-07 `Simplify` cannot create `IDENTIFIED` provenance | `ops/simplify.go` | `TestSimplifyNeverIdentifies` (K) | — |
| REQ-032-08 missing Identify justification | `session` | `TestIdentifyRequiresJustification` | — |
| REQ-032-09 candidate contamination bypass | `session` validation | `TestCandidateArtifactContainmentRejected` (L) | — |
| REQ-032-10 absence of promotion API | whole tree | `TestNoForbiddenExportedAPI` | G-Audit |
| REQ-032-10a external candidate JSON loads only as `UnverifiedResearchCandidate`; no public parse path returns trusted `ResearchCandidate` or `core.Object` values | `session` | `TestParseResearchCandidateJSONUnverifiedOnly` | G-Audit |
| REQ-032-11 corrupted ledger hash chain | `session/ledger.go` | `TestLedgerTamperDetection` case 1 (Q) | — |
| REQ-032-12 modified intermediate canonical output | `session/ledger.go` | Q case 2 | — |
| REQ-032-13 modified output, recomputed local hash, stale step hash | `session/ledger.go` | Q case 3 | — |
| REQ-032-14 recomputed full chain, replay divergence | `session/ledger.go` | Q case 4 | — |
| REQ-032-15 inconsistent manifest | `core/corpus.go` | `TestManifestRejectsInconsistent` (O; unknown field, bad enum, non-canonical expr, duplicate ID) | — |
| REQ-032-16 exact-rational canonical JSON round-trip | `kernel` encoder | P: `TestExactRationalRoundTrip` | — |
| REQ-032-17 nondeterministic canonicalization attempt | `kernel` sorter | `TestCanonicalOrderingStable` (different build orders ⇒ same bytes) | — |
| REQ-032-18 unsupported general solver form | `ops/relation.go` | `TestSolveUnsupportedForm` | — |
| REQ-032-19 unsupported derivative node/function form | `ops/transform.go` | `TestDifferentiateUnsupportedForms` (`Sqrt`, `Call`, `Relation`, `BranchSet`, non-integer exponent) | — |
| REQ-032-20 unsupported general limit form | `ops/transform.go` | `TestLimitUnsupportedForm` | — |
| REQ-032-21 invalid session state transition | `session` | `TestSessionStateMachine` (every illegal transition ⇒ typed error) | — |
| REQ-032-22 post-seal mutation attempt | `session` | `TestPostSealMutationRejected` (all nine actions after `Sealed`) | — |
| REQ-§33-MUST-01 (MVP complete only when every acceptance test A–S passes) | §11 slices | §13 gate 5 | — |

**Acceptance A–S (specs_v2_3.md §33)** — every row maps to a §11 slice test:

| Acceptance | Requirement | Test function | File |
|---|---|---|---|
| A typed mechanics | REQ-001-01, REQ-005-01 | `TestTypedMechanicsConstructors` | `mechanics/relations_test.go` |
| B classical relation | REQ-001-03, REQ-§28.1 | `TestNewtonSecondLawConstruct` + `TestNewtonSecondLawManifestMatch` | `mechanics/*_test.go` |
| C dimension rejection | MRC-002, REQ-032-01 | `TestDimensionMismatch` | `ops/negative_test.go` |
| D category rejection | MRC-003, REQ-032-02 | `TestKindMismatchEqualDimensions` | `ops/negative_test.go` |
| E assumption conflict | MRC-004, REQ-032-03 | `TestAssumptionConflict` | `ops/negative_test.go` |
| F convention conflict | MRC-005, REQ-032-04 | `TestConventionConflict` | `ops/negative_test.go` |
| G differentiation | §21 | `TestDifferentiateKineticEnergy` | `mechanics/relations_test.go` |
| H relativistic relation | REQ-001-02 | `TestEnergyMomentumRelation` | `relativity/derivation_test.go` |
| I mass-energy derivation | §20 | `TestMassEnergyDerivation` | `relativity/derivation_test.go` |
| J relativistic limit | §22 | `TestLorentzFactorLimit` | `relativity/derivation_test.go` |
| K Identify firewall | MRC-006, REQ-013-02 | `TestSimplifyNeverIdentifies` + `TestSessionIdentifyRecords` | `ops/negative_test.go`, `session/session_test.go` |
| L candidate containment | MRC-008, REQ-024 | `TestHypothesisContamination` + `TestCandidateArtifactContainmentRejected` | `hypothesis/candidate_test.go` |
| M falsifiability | §25 | `TestSealedCandidatePreservesFalsifiability` | `hypothesis/candidate_test.go` |
| N anomaly | §28.3 | `TestCandidateReferencesManifestAnomaly` | `hypothesis/candidate_test.go` |
| O manifest | REQ-001-03 | `TestMechanicsManifestCrossCheck` + `TestRelativityManifestCrossCheck` | both `manifest_test.go` |
| P exact rational round-trip | REQ-032-16 | `TestExactRationalRoundTrip` | `core/expr_test.go` |
| Q ledger tamper/replay | REQ-032-11..14 | `TestLedgerTamperDetection` | `session/session_test.go` |
| R determinism | REQ-015-01, §31 | `TestDerivationDeterminism` | `session/session_test.go` |
| S candidate handoff | REQ-000-03, §26 | `TestSealResearchCandidate` | `session/session_test.go` |

Plus the non-acceptance positional-contract row of §11: `TestApplyPositionalInputs` (`ops/operations_test.go`) covering §15.13.1 for all 12 IDs.

### 12.12 Evaluation, documentation, deferrals, plan-level requirements (§34–§41 + prompt v2.3)

| Requirement ID | Implementation location | Primary test | Secondary |
|---|---|---|---|
| REQ-§34-MUST-01 (each populated framework evaluated through three complementary mechanisms: deterministic derivation `F=m*a` / `E=m*c^2`, executable reduction `lim(v→0) LorentzFactor(v)=1`, limitation/anomaly record) | slice tests B/G/I/J/N + both manifests | §13 gate 5 (three-mechanism inventory per framework) | — |
| REQ-§34.2-MUST-01 (MVP executable reduction is the Lorentz-factor limit; Newtonian kinetic-energy series recovery deferred — no series code) | `ops/transform.go`, slice J | J: `TestLorentzFactorLimit` | `TestNoForbiddenExportedAPI` (`Series`) |
| REQ-§34.3-MUST-01 (framework manifest contains ≥1 limitation/anomaly record; machine-readable context only — library never infers what to do about an anomaly) | both `manifest.json` files | O + N (cross-ref REQ-§28.3-MUST-01) | — |
| REQ-§35-MUST-01 (`docs/paper-translation.md` contains the three sections: `Common notation`, `Framework mapping`, `Ambiguity resolution`) | `docs/paper-translation.md` | `TestDocsSectionsPresent` (docs audit in `core/object_test.go`) | — |
| REQ-§35-MUST-02 (translation document MUST NOT be required for runtime execution; manifests are the canonical machine-readable metadata) | docs only | `TestNoDocsImport` (source scan: no package imports `docs`) | G-Imports |
| REQ-§36-MUST-01 (no `physvet` in MVP; future `physvet` MUST consume the same public/core semantic contracts and metadata — recorded as v0.5+ constraint) | plan §14 deferral 1 | `TestNoDeferredPackages` (no `physvet` package) | G-Scope |
| REQ-§36-MUST-02 (future `physvet` MUST NOT independently redefine dimensions, physical kinds, assumptions, conventions, provenance, operation contracts, corpus status — recorded, not implemented) | plan §14 deferral 1 note | `TestNoDeferredPackages` | G-Scope |
| REQ-§37-MUST-01 (all 22 v0.5+ deferrals outside MVP; no deferred feature leaks in via placeholder packages or speculative abstractions) | plan §14 (below) | `TestNoDeferredPackages` + `TestRepositoryTreeExact` | G-Scope |
| REQ-§38-MUST-01 (implementation plan MUST NOT add files outside the §3 tree) | plan §2 tree | `TestRepositoryTreeExact` (39 files) | §13 gate 3 |
| REQ-§38-MUST-02 (prohibited: convenience-only files, placeholder abstractions, compatibility shims, future scaffolding, generated registries, empty packages) | tree | `TestRepositoryTreeExact` + `TestNoDeferredPackages` | G-Scope |
| REQ-§39-MUST-01 (this matrix maps every normative `MUST`/`MUST NOT` clause to an implementation location and ≥1 test; no normative clause omitted) | plan §12 | §13 gate 8 (completeness check against `must_by_sec` inventory) | — |
| REQ-§39-MUST-02 (deterministic coverage IDs such as `REQ-§15.8-MUST-01` assigned where no explicit `REQ-*` label exists) | plan §12 ID convention | §13 gate 8 | — |
| REQ-§39-MUST-03 (grouped rows such as `REQ-002-01..26`, `REQ-005-01..12`, `REQ-032-01..22` expanded into individual requirement rows) | plan §12.1/12.3/12.11 | §13 gate 8 | — |
| REQ-§40-MUST-01 (final deterministic loop invariant: manifest → typed construction → immutable object → explicit assumptions → pure ops → MRC → session draft/commit → Identify → replayable ledger → provisional hypothesis → prediction + falsification → anomaly context → sealed candidate → human review; the system stops there) | plan §1–§11 | §13 gates 2/4/5 (architecture + A–S loop) | — |
| REQ-§40-MUST-02 (governing design principle quoted verbatim: open representation/closed authority; explicit assumptions/explicit insight; deterministic mechanics/fallible MRC; machine-readable physics/human scientific judgment) | README + plan §1 | README audit | — |
| REQ-§41-MUST-01 (handoff invariant: a coding agent MUST NOT need to choose among alternative architectures; any implementation requiring reopening one of the 21 resolved decisions is non-conforming) | plan §1–§12 | §13 gate 10 (plan conformance review) | — |
| REQ-§Prompt-MUST-01 (prompt "Governing source of truth": `specs_v2_3.md` is normative and authoritative; v2.2/`plan9.md` never normative) | plan header, plan §1 | §13 gate 10 | — |
| REQ-§Prompt-MUST-02 (prompt §1: invariant architecture preserved — `core`/`ops`/`session` import restrictions, kernel authority, aliases, sole mint, pure `ops`, `Session.Identify` only, session-owned ledger/candidate, contamination, no generic factory, no truth API, no MRC bypass; no `core.SealedDerivation`, no `review/` package, no domain registry, no alternative layering) | plan §1, §4, §5 | §13 gates 2/4 | — |
| REQ-§Prompt-MUST-03 (prompt §2: all v2.2 authority language replaced — no `Normative source: plan9.md = specs v2.2` anywhere; `plan9`/`plan10.md` marked historical only) | plan header | §13 gate 10 (text scan of this plan) | — |
| REQ-§Prompt-MUST-04 (prompt §3A: `Challenge`/`Review` defined exclusively in `core/corpus.go`; `session` consumes without redefinition; no alias implying a second declaration; no `review/` package) | `core/corpus.go`, `session` | `TestSupportTypeSurface` (AST) | §12.10 rows |
| REQ-§Prompt-MUST-05 (prompt §3B: §9.6 finite-factor rule used exactly — `base != 0` ⇒ non-singularity/admissibility; `base > 0` ⇒ sign-sensitive rewrite; `c^-1` allowed under `c > 0`; `0^0` ⇒ `UnsupportedOperationError`; `0^-1` rejected; no Lorentz-specific simplification hack) | `ops/simplify.go` | `TestFiniteFactorRules` (`Pow(x,0)`, `Pow(c,-1)` under `c>0`, `0^0`, `0^-1` cases) | §12.4 rows |
| REQ-§Prompt-MUST-06 (prompt §3C: exact §15.13.1 positional mapping reproduced in plan §9; `ops.Apply`, `Session.Step`, and replay use the same mapping; inputs never reordered; `identify` excluded; `OperationParams` schema + canonical JSON + permitted `Kind`s pinned — never opaque) | `ops/dispatch.go`, `session` | `TestApplyPositionalInputs`, `TestParamsCanonicalRoundTrip` | §12.6 rows |
| REQ-§Prompt-MUST-07 (prompt §4: all five net-valid findings F1–F5 incorporated — 10-item relativity manifest with `Velocity` wrapper distinction; gated `Pow(x,0)`; explicit kernel-side metadata construction in Step 1 with single authoritative validation site; deterministic `AssumptionSet`/`ConventionSet` ordering; Test D fixture `Mass` vs `RestMass` used consistently) | plan §3, §6, §7, §11, §12 | O (manifest set), P/F2 cases, D, F4 ordering tests | §13 gate 8 |
| REQ-§Prompt-MUST-08 (prompt §5 precision items: constructor-time vs Simplify-time boundary resolved; `Simplify` on `Relation`/`BranchSet` recursive, never `IDENTIFIED`; exact §15.12 `SelectBranch` six-step contract with `selected_branch/<hash>` key and no invented sign-entailment error; defensive copies enumerated; kernel helper ownership in `internal/kernel/types.go`; no public `DecodeObjectJSON`/`ParseObject`; `schema_version:"1"` labeled Plan 10 implementation pin, serialization DTO only; kernel exports controlled canonical helpers) | plan §4, §6 | `TestConstructorCanonicalization`, REQ-§9.9-MUST-01 row (relation/BranchSet simplification exercised in slice I), `TestSelectBranchContract`, `TestExprDefensiveCopies`/`TestObjectDefensiveCopies` | — |
| REQ-§Prompt-MUST-09 (prompt §6: stale v2.2 material removed — no v2.2 assumption-kind lists (`fixed_constant`, `state_relation`, `framework_relation`); actual v2.3 kinds `Domain, Regime, Constraint, Convention, Approximation, MathPrecondition, PhysicalAssumption` with JSON `domain…physical_assumption`; no old Challenge/Review placement; no `OperationParams` ambiguity; no old manifest counts; no stale REQ semantics) | plan §8, §12.5, §12.10 | `TestAssumptionKindEnumExact` (7 Go constants, 7 lowercase JSON values) | §13 gate 10 (text scan) |
| REQ-§Prompt-MUST-10 (prompt §7: matrix rebuilt against the complete `specs_v2_3.md` — every explicit `REQ-*`, expanded `REQ-002-01..26`, every unnumbered `MUST`/`MUST NOT` with deterministic IDs, every MRC-001..008, A–S, implementation location, ≥1 named test, §15.13.1 positional mapping, §§18–41 covered, no obsolete v2.2-only references; mechanically verified) | plan §12 | §13 gate 8 | — |
| REQ-§Prompt-MUST-11 (prompt §8+§9: Plan 10 execution strengths preserved — purpose/dependencies/files/result/gate steps, dependency-ordered sequence, explicit test-package plan, A–S inventory, file-budget audit, source/AST/import audits, determinism run (twice + fresh `Session`), Lorentz-body and E=mc² anti-hardcoding checks; plan stays independent of Plan 9.2) | plan §2, §3, §11, §13 | §13 gates 3/4/5/7 | — |
| REQ-§Prompt-MUST-12 (prompt §10+§11: final checks satisfied — 39 files; exactly 12 `ops` operations with `identify` excluded; exactly nine session actions per §16.1; unverified-only candidate loading; exact 10 relativity manifest items; gated `Pow(x,0)`; deterministic canonicalization for all 10 listed artifacts; canonical object field order with serialization-only metadata kept out of kernel storage; no premature `Open Spec Items: NONE` claim; deliverable = `plan10_v2_3.md` only, no code) | plan, workflow | §13 gates 1/3/5/8/10 | — |

---

## 13. Acceptance gate

The MVP is accepted only when every condition below passes from the repository root, with zero warnings and zero skipped tests.

**Exact commands (run in order):**

```text
go build ./...
go vet ./...
go test ./...
```

**Gate checklist:**

1. **G-Module** — `go.mod` contains `module github.com/PithomLabs/phys` and `go 1.24`; it declares no third-party requirements (an empty or tooling-normalized `require` block is equally conforming; no `go.sum` is needed); every internal import resolves under `github.com/PithomLabs/phys/...` (REQ-000-01, REQ-000-02, REQ-§0.1-MUST-01/02).
2. **G-Imports** — `go list -deps ./...` contains only standard-library packages (zero third-party modules); AST import audit asserts: `core` imports neither `ops` nor `session` (REQ-§4-MUST-01/02); `ops` imports neither `session` nor domain packages (REQ-004-01); `session` imports neither `mechanics`/`relativity`/`hypothesis` (REQ-§26.8-MUST-03, prompt §1); no non-test file imports `go/parser`, `go/scanner`, `go/token`, `math/rand`, `net/http`, or `os/exec`; `kernel` is `internal/` (REQ-004-03).
3. **G-Scope** — exact file-tree audit: exactly the 39 files of plan §2 (24 Go source, 10 `*_test.go`, 2 `manifest.json`, 3 config/doc), no file outside the §2 list (REQ-003-01, REQ-§38-MUST-01/02), no deferred package (`physvet`, `electromagnetism`, `qm`, `qft`, `statmech`, `review`), no scaffold/generated/convenience files; tests are adjacent `*_test.go` only (REQ-003-02), with external test packages wherever package visibility is the subject (REQ-003-03).
4. **G-Audit** — source/AST audits: no `float32`/`float64`/`math.` in non-test sources (REQ-007-01 secondary); no `map[string]any` in canonical encoders (REQ-010-02); exactly one `func MintObject` in the module, under `internal/kernel`, not re-exported (REQ-005-04, REQ-§5.2.1-MUST-01/03, REQ-§5.4-MUST-01); no exported `NewObject` in public packages; no promotion/bypass/truth API names (`Promote`, `Trust`, `ApproveHypothesis`, `PromoteToEstablished`, `SetCorpusStatus`, `ApproveException`, `OverrideMRC`, `BypassMRC`, `Integrate`, `Series`, `Taylor`, `Simulate`, `RankTheories`, `ScoreTruth`, `TruthScore`, `ProbabilityOfTruth`, `BestTheory`) (REQ-013-05, REQ-§24.4-MUST-01, REQ-§29-MUST-02); `ops/transform.go` expands the fixed Lorentz body rather than matching a function ID (REQ-§15.9.1-MUST-02); `Simplify` never creates `IDENTIFIED` (REQ-§15.6-MUST-01); no package-level mutable state in `ops`/`session` (REQ-015-01, REQ-§16.7-MUST-02); `statement` never parsed (REQ-§17.6-MUST-01); no `ops.Identify` export (MRC-006); no `review/` package and `Challenge`/`Review` declared only in `core` (REQ-§27-MUST-01).
5. **Acceptance A–S** — every function named in §12.11's acceptance table exists with the exact name and package, and passes under `go test ./...` (19 rows, 21 test functions; specs_v2_3.md §33), **plus `TestApplyPositionalInputs` proving the §15.13.1 positional contract for all 12 `OperationID`s (wrong-order inputs rejected; `ops.Apply`, `Session.Step`, and replay share the same indices)** (REQ-§15.13.1-MUST-01/02/03).
6. **Negative suite** — all 23 REQ-032 rows (01–22 plus 10a) present and passing (specs_v2_3.md §32), including all four ledger tamper cases (REQ-032-11..14).
7. **Determinism** — `TestDerivationDeterminism` runs the canonical derivation twice in-process and once with a fresh `Session`; canonical artifacts, hashes, ledger JSON, and candidate JSON are byte-identical (REQ-§31-MUST-01/02/03, REQ-§31.1-MUST-01/02).
8. **Coverage completeness** — every MUST/MUST NOT clause in `specs_v2_3.md` §0–§41 and every MRC-001..008 ID appears in plan §12 with an implementation location and ≥1 named test; grouped REQ ranges are expanded to individual rows; deterministic coverage IDs follow the `REQ-§N-MUST-nn` convention (REQ-§39-MUST-01/02/03). Verification is mechanical: extract every `MUST`/`MUST NOT` line and every labeled `REQ-*` from `specs_v2_3.md`, then confirm each appears in §12 before claiming `Open Spec Items: NONE`.
9. **Documentation** — README contains module identity, dependency diagram, constructor-authority statement (external/public API boundary, not cryptographic — REQ-000-04), MRC list with `mrc-v0.4` fallibility note (REQ-000-03, REQ-§29-MUST-01/03), integrity-vs-authenticity statement, and non-goals/deferrals; `docs/paper-translation.md` contains the three §35 sections (REQ-§35-MUST-01/02).
10. **Plan conformance** — implementation invents no architecture beyond plan §1–§12, adds no file outside plan §2, uses stdlib only, and leaves no normative requirement without a test mapping; `schema_version:"1"` appears only as the labeled Plan 10 serialization pin, never as kernel state (REQ-§Prompt-MUST-12, REQ-§Prompt-MUST-08). Any implementation that requires reopening one of the 21 resolved decisions of specs_v2_3.md §41 is non-conforming (REQ-§41-MUST-01).

A build that requires reopening any specs_v2_3.md §41 decision is non-conforming.

---

## 14. Explicit deferrals

All of the following are v0.5+ and excluded from MVP (specs_v2_3.md §37, repeated exactly):

1. `physvet` implementation.
2. General symbolic series expansion.
3. General truncation-order propagation.
4. Full 1905 Einstein low-velocity series derivation.
5. Symbolic integration.
6. Broader symbolic calculus.
7. General equation solving.
8. General power/branch algebra.
9. Richer tensor/index machinery.
10. General relativity.
11. Quantum mechanics.
12. Quantum field theory.
13. Statistical mechanics.
14. Electromagnetism.
15. Numerical execution.
16. Empirical-data adapters.
17. Reviewer-AI orchestration.
18. Automatic corpus governance.
19. Automatic hypothesis promotion.
20. MRC exception/override workflow.
21. EBP 2.1 integration.
22. General theory-management platform.

No deferred feature may leak into MVP through placeholder packages or speculative abstractions (REQ-§37-MUST-01; enforced by `TestNoDeferredPackages` + `TestRepositoryTreeExact`). The prompt's deferral list (general series expansion, full 1905 Einstein derivation, broader symbolic analysis, integration, richer tensor/index machinery, GR/QM/QFT/stat-mech/EM, numerical execution, empirical adapters, reviewer orchestration, MRC exception workflow, automatic corpus governance, automatic hypothesis promotion, EBP 2.1, physvet) is a subset of this 22-item list — no additions to scope. Reviewer-AI orchestration also covers specs_v2_3.md §27.4 (no orchestration code exists).

---

## 15. Scope accounting

**File count:** exactly **39** files including `go.mod` (38 excluding it), against the ≤40 budget of specs_v2_3.md §38 (REQ-§38-MUST-01).

**Exact split:**

| Class | Count | Files |
|---|---|---|
| Go source | 24 | `internal/kernel/{types,mint}.go` (2); `core/{object,expr,dimension,assumption,convention,provenance,corpus,canonical,errors}.go` (9); `ops/{arithmetic,simplify,transform,relation,dispatch}.go` (5); `session/{session,ledger,research_candidate}.go` (3); `mechanics/{primitives,relations}.go` (2); `relativity/{primitives,relations}.go` (2); `hypothesis/candidate.go` (1) |
| Go test | 10 | `core/{object,expr}_test.go` (2); `ops/{operations,negative}_test.go` (2); `session/session_test.go` (1); `mechanics/{manifest,relations}_test.go` (2); `relativity/{manifest,derivation}_test.go` (2); `hypothesis/candidate_test.go` (1) |
| JSON manifest | 2 | `mechanics/manifest.json`, `relativity/manifest.json` |
| Config/documentation | 3 | `go.mod`, `README.md`, `docs/paper-translation.md` |

**Major exported types:** `core.Object` (= `kernel.Object`), `core.Expr`, `core.Kind` (18 constants), `core.Dimension` (+9 constructors), `core.Assumption`/`AssumptionSet`, `core.Convention`/`ConventionSet`, `core.Provenance`, `core.CorpusStatus`, `core.Manifest` (+ `ManifestDomain/Limit/Anomaly/Reduction/Item`), `core.Challenge`, `core.Review`, the 11 typed errors (`DimensionMismatchError` … `LedgerValidationError`), `ops.OperationID`, `ops.OperationParams`, `session.Session`, `session.Ledger`, `session.Step` (read-only accessors), `session.DraftMetadata`, `session.ResearchCandidate`, `session.UnverifiedResearchCandidate`, `session.FrameworkDependency`, `session.Prediction`, `session.FalsificationCondition`, `session.RecoveryClaim`, `session.AnomalyReference`, `hypothesis.NewCandidateConcept` (returns a `core.Object`).

**Major operations:** 12 pure operations (`Add`, `Subtract`, `Multiply`, `Divide`, `Pow`, `Simplify`, `Substitute`, `Differentiate`, `Limit`, `Compare`, `Solve`, `SelectBranch`) behind the closed `ops.Apply` dispatch with the exact §15.13.1 positional schema (plan §9); `identify` excluded from `ops.Apply`; 9 session actions (`Postulate`, `Declare`, `Define`, `Step`, `Identify`, `Conclude`, `Draft`, `Commit`, `Seal`) plus read-only `Validate`/`CanonicalJSON`; manifest API (`ParseManifest`, `ValidateManifestBytes`, `CanonicalManifestJSON`, `Manifest.Hash`); candidate API (`Seal`, `Validate`, `CanonicalJSON`, `ParseResearchCandidateJSON`); 9 mechanics + 8 relativity value constructors and the relation constructors of §18.2/§19.1; `kernel.MintObject` (the single internal mint point).

**Expected dependency count:** **0** third-party Go modules — `go.mod` declares no third-party requirements (REQ-000-02); the empty `require` block is a presentation detail, not a conformance criterion. Intra-repo compile-time package edges (12 total, verified by G-Imports): `internal/kernel` → ∅; `core` → `internal/kernel`; `ops` → `core`, `internal/kernel`; `session` → `core`, `ops`, `internal/kernel`; `mechanics` → `core`, `internal/kernel`; `relativity` → `core`, `internal/kernel`; `hypothesis` → `core`, `internal/kernel`. Test-only edges from `*_test.go` into `ops`/`session`/domain packages are allowed and additional (REQ-003-02/03).

**Canonicalization inventory (prompt §10):** deterministic ordering with no map-order dependence is pinned for `Expr`, `AssumptionSet`, `ConventionSet`, `Object`, `Manifest`, `OperationParams`, `Step`, `StepEnvelope`, `Ledger`, `ResearchCandidate`; canonical object field order is `schema_version, valid, name, kind, dimension, expr, assumptions, conventions, provenance, corpus_status`, where `schema_version:"1"` is a **Plan 10 implementation pin** on the serialization DTO only — authoritative `kernel.Object` storage carries no serialization-only state (prompt §5).

---

## Revision Audit

### Preserved Architecture

Plan 10's architecture is unchanged from `plan10.md` wherever v2.3 permits it, and `specs_v2_3.md` §41 confirms every element of it as a resolved, normative choice:

- Layering `internal/kernel → core → ops → session`, domains (`mechanics`, `relativity`, `hypothesis`) → `core` + `internal/kernel`; no cycles; `core` never imports `ops`/`session`; `ops` never imports `session`; `session` never imports domain packages (§1, §12.2).
- Single production mint (`kernel.MintObject`), aliases in `core`, no generic public factory, no public object decoder (internal kernel decode only).
- Exactly 12 pure operations behind the closed `ops.Apply` dispatch; `Session.Identify` as the only identification API; no ambient session state.
- `session` owns ledger, replay, and sealing; nine session actions (§16.1); genesis/StepID/StepEnvelope/17-step Validate contracts (§16.15–16.19).
- Candidate containment (MRC-008) enforced in `ops` and `session`; `ResearchCandidate` sealed only by `Session.Seal()`; external JSON loads only as `UnverifiedResearchCandidate`.
- `Challenge`/`Review` in `core/corpus.go` only; no `review/` package; no truth/adjudication, promotion, bypass, or orchestration APIs anywhere.
- Exact 39-file tree (24 source / 10 test / 2 manifest / 3 config-doc), stdlib only, `github.com/PithomLabs/phys`, `go 1.24`.
- Execution strengths kept: 16 dependency-ordered steps with purpose/deps/files/result/gate, explicit test-package plan, A–S acceptance inventory, file-budget audit, source/AST/import audits, determinism run (twice + fresh `Session`), Lorentz-body and `E=mc²` anti-hardcoding checks.

### Net-valid Review Findings Adopted

All findings re-checked against the complete `specs_v2_3.md`; only those still valid under v2.3 were adopted:

- **F1** — relativity manifest pinned to exactly the §28.2 10-item set; `Velocity` is a wrapper, not a manifest item; `RestFrame` is an assumption; mechanics pinned to the §28.1 11-item set (§7, §12.8, §12.11).
- **F2** — `Pow(x,0) → 1` gated on a nonzero-safe base added to Simplify; `0^0` remains `UnsupportedOperationError`; `0^-1` rejected; §9.6 rules intact (§6, Step 5, slice P).
- **F3** — Step 1 explicitly establishes kernel-side construction primitives for `Dimension`, `Assumption`, `AssumptionSet`, `Convention`, `ConventionSet`, `Provenance`, with one authoritative validation site; Step 2 façade constructors delegate (Steps 1–2).
- **F4** — deterministic `AssumptionSet`/`ConventionSet` ordering pinned: canonicalize → canonical bytes → lexicographic sort → dedupe; never map order; participates in JSON → set hash → object hash → ledger chain → candidate determinism (Step 1, §8, slice R).
- **F5** — Test D fixture is exactly `Mass` vs `RestMass` (equal dimension `M`, distinct kinds ⇒ `CategoryMismatchError`, isolates MRC-003) in the slice table, test body, matrix (REQ-032-02, MRC-003 row), and acceptance D.
- **Reconciliation A** — `Challenge`/`Review` exclusively in `core/corpus.go`; `session` consumes without redefinition or alias (§1, §10, §12.10 rows REQ-§26.0/§27.x).
- **Reconciliation B** — §9.6 finite-factor rule stated exactly (`base != 0` admissibility vs `base > 0` sign rewrite), `c^-1` allowed under `c > 0`, no Lorentz-specific hack (§6, REQ-§Prompt-MUST-05).
- **Reconciliation C** — exact §15.13.1 positional table and `OperationParams` schema reproduced once in §9 and shared by `ops.Apply`, `Session.Step`, and replay; inputs never reordered; `identify` excluded (§9, §12.6 rows REQ-§15.13.1-MUST-01/02/03, new test `TestApplyPositionalInputs`).
- **Precision items (prompt §5)** — constructor-time vs Simplify-time boundary resolved; recursive `Simplify` on `Relation`/`BranchSet` that never creates `IDENTIFIED`; exact §15.12 `SelectBranch` six-step contract with `selected_branch/<hash>` key and no invented sign-entailment error; enumerated defensive copies; kernel helper ownership in `internal/kernel/types.go`; no public `DecodeObjectJSON`/`ParseObject`; `schema_version:"1"` labeled a Plan 10 serialization pin, kept out of `kernel.Object`; controlled kernel canonical-helper exports.
- **Matrix/structure** — coverage matrix rebuilt row-by-row from `specs_v2_3.md` (subsection-precise deterministic IDs replacing v2.2's flat `REQ-§N` rows); new rows for §15.13/§15.13.1, §11.1, §26.8/26.10 restructure, §27, §29, §39; nine session actions per §16.1.

### Stale v2.2 Material Removed

- `plan9.md`/specs v2.2 as normative source → replaced by the header declaration that `specs_v2_3.md` is the sole normative source; `plan9.md`/`plan10.md` appear only as historical/background (prompt §2).
- v2.2 assumption-kind lists (`fixed_constant`, `state_relation`, `framework_relation`) → exact seven v2.3 §11.1 kinds with lowercase JSON values (§4, §8, row REQ-§11.1-MUST-01, test `TestAssumptionKindEnumExact`).
- "Challenge/Review re-declared as aliases of `core` types" → v2.3 consumption model with no alias (§10, §12.10).
- Opaque/ambiguous `OperationParams` → exact §15.13 field order, canonical JSON, permitted `Kind` set, and per-kind binding (§9).
- Old/ambiguous manifest item counts and the `Velocity`-as-manifest ambiguity → exact §28.1/§28.2 sets (§7, §12.11).
- Stale REQ semantics: flat `REQ-§18/§19/§20…` IDs → subsection-precise `REQ-§18.0/§19.0/§20.1…`; `REQ-011-02b` → `REQ-§11.1-MUST-01`; v2.2 §26.8/§26.9 rows restructured to v2.3's §26.9/§26.10 split; "seven-stage pipeline" wording normalized to §26.8's seven-stage list; `ReviewNotes` → `ReviewerNotes`.
- Dimensionally invalid Test D example (`KineticEnergy` vs `Momentum`) → `Mass` vs `RestMass` (F5).
- No `Normative source: plan9.md = specs v2.2` statement, no v2.2-only matrix row, and no v2.2 manifest count remains anywhere in this document.

### v2.3 Coverage Verification

Mechanical checks performed against the complete `specs_v2_3.md` (4633 lines, §0–§41) before writing this audit:

- **115 of 115** explicit `REQ-*` labels extracted from v2.3 appear in this plan (0 missing).
- **126 of 126** spec sections containing at least one `MUST`/`MUST NOT` clause map to at least one §12 matrix row (via a labeled `REQ-*` or a deterministic `REQ-§N-MUST-nn` ID) — 0 uncovered.
- Every `MRC-001`…`MRC-008` ID appears in §12.6 with implementation location and named test.
- Acceptance `A`–`S` (19 rows, 21 test functions) mapped in §12.11 and §11; positional contract covered by `TestApplyPositionalInputs` (all 12 `OperationID`s, §15.13.1).
- Grouped ranges expanded individually: `REQ-002-01..26` (26 rows), `REQ-032-01..22` + `10a` (23 rows), `REQ-005-*`, `REQ-008-*`, `REQ-010-*`, etc.
- The exact §15.13.1 table and `OperationParams` schema appear once in §9 and are referenced by rows in §12.6.
- Stale-term scan: no non-historical `plan9`/v2.2 normative reference, no v2.2 assumption kinds, no `ReviewNotes`, no old manifest counts.
- §13 gate 8 encodes the same mechanical procedure as an implementation-time re-check.

### Open Spec Items

NONE

---

The architecture is fixed by `specs_v2_3.md` §41 and this plan adds nothing beyond it: a coding agent reading only `specs_v2_3.md` + `plan10_v2_3.md` has no remaining architectural choices to make — no packages or files to invent, no positional mappings to infer, no manifest membership to guess, and no contradictory instructions to resolve. Any implementation that must reopen one of the 21 resolved §41 decisions, or add a file outside plan §2, is non-conforming.

# Plan 8 — Freeze Delta (G1–G11 consolidation over Plan 6 + Plan 7)

**Status:** PLAN ONLY — no implementation performed. Awaiting explicit instruction.

**Precedence (G8):** `specs_v2_3.md` is authoritative; Plan 7 is additive gates;
conflicts resolve to spec. **Location (G9):** this file lives in `plan10/`,
outside the module tree — never counted toward the 39-file implementation tree
(REQ-003-01).

**Scope:** narrow hardening + mechanical proof. No redesign, no features, no new
physics domains (GR/QM/EM deferred), no 1905 corpus, no parser/CAS/numerics/
prover/solver/promotion/truth/EBP/AI-planner. 39-file tree, package graph,
stdlib-only preserved. Tests live in existing files only.

## Delta 1 — H2 corrected (G1/G2/G3)

`internal/kernel` allowlist = spec-mandated closed machinery, not two strings:

- closed Kind enum (§6, 18 values + ordinals), ExprKind (§8.2.1),
  AssumptionKind (§11.1), ProvenanceStatus (§13.1), CorpusStatus (§13.4);
- `KindMinkowski`, `LorentzFactorFunctionID` / `"lorentz_factor"`;
- only the closed-world validation those mandates require.

Forbidden: laws, formulas, constants, framework assumptions, symbol-dispatch
masquerading as algebra. `ops` keeps the full mandated Lorentz hook cluster
(body constructor + `expandCalls` + Call-finiteness), nothing more.

Scanner mechanics (pinned): `go/parser` **without** comment mode (comment-blind
by construction); `"v"`/`"c"` literals allowed only inside the Lorentz
body-construction function in `ops/transform.go` (function-scoped, resolved
in-implementation). Test lives in `core/object_test.go` beside existing AST
infra. Explicitly NOT required: crypto authenticity, `unsafe`/`reflect`/
`linkname` resistance, general adversarial-AST fuzzing (Gate A battery is the
panic boundary).

## Delta 2 — H3/H4 tightened (G4/G5/G6)

- **G4:** new `TestProvenanceAssumptionPreservation`
  (`ops/operations_test.go`): per pure op, inherited assumptions preserved
  **verbatim** (canonical values, not keys); only enumerated generated keys
  (`denominator/*`, `selected_branch/*`, others ∅) may appear; same
  key + different value → `AssumptionConflictError` (existing merge law).
- **G5:** reverse allowlist status-independent — **every** exported constructor
  in `mechanics`/`relativity`/`hypothesis` ⊆ manifest constructor IDs ∪ pinned
  set `{NewKineticEnergy, Velocity, ZeroThreeMomentum, ZeroEnergy,
  ZeroVelocity, hypothesis concept}`. No inference from status/naming.
- **G6:** leakage checks scoped to **constructor outputs** (mixed-framework op
  derivations legitimately merge keys; never exact-pin full sets; forbidden
  cross-framework keys). MRC-008 wording: no hypothesis-derived step output
  carries `ESTABLISHED` or non-`HYPOTHESIS` provenance (candidates may contain
  established premises by design).

## Delta 3 — H6 symmetric params (supersedes pow-only)

Per-kind unused-field rejection in `validateShape`, tested explicitly for
all four `OperationParams.Kind` values:

- `empty` → Exponent="", Operator="", Justification="";
- `pow` → Exponent required, Operator="", Justification="";
- `compare` → Operator required, Exponent="", Justification="";
- `identify` → Justification required, Exponent="", Operator="".

Align `operations_test.go` fixtures to canonical empty form; the old
`"operator":"eq"` representation lives only in the negative test. Spec text
untouched (implementation pin).

## Delta 4 — H7 baseline pins (G11, before freeze)

- **Sources, both corpora (code + manifest.json, cross-checks stay green):**
  mechanics keeps NewtonSecondLaw → `"Newton, Principia"`, other 10 items →
  `"Classical Mechanics corpus"`; relativity keeps EMR/MER → `"Einstein,
  1905"`, other 8 → `"Special Relativity corpus"`. Provenance-metadata
  alignment only.
- **Corpus status:** keep framework-level ESTABLISHED asserts; add per-item
  `corpus_status` to cross-checks; `established`→`contested` flip fails the
  **freeze baseline** test while remaining schema-legal (no enum prohibition;
  status is loaded, never derived).
- **H0 → README (G10):** kernel/theory/hypothesis wording + E=mc² demo status
  (formal derivation from encoded SR premises ≠ 1905 reconstruction ≠
  empirical proof); extend docs-section assertion minimally if needed.

## Unchanged (Plan 6 + Plan 7 as written)

Gate A battery (`ops/negative_test.go`, `session/session_test.go`, typed
errors, no partial commit); Gate B 12-op matrix; Gate C cross-check extensions;
Gate E invalid-expr candidate negatives; Gate F in-memory corruption battery;
rational/determinism pins; replay safety (retention + stage order untouched);
strict no-`MassEnergyRelation` source test + full E=mc²/`m*v`/Lorentz traces;
M1–M18 + new-gate mutants in disposable copy (`-count=1`, restore per mutant,
claim only rerun results); §17 A–E report with `FREEZE-READY` /
`NOT FREEZE-READY`.

## Non-blockers (do not expand)

Justification prose quality; absolute zero-physics rule; `go:embed` location;
authenticity crypto; AI planner; EBP runtime; 1905 corpus.

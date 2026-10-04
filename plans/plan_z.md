# Physics Compiler in Go — Design (Phase B architecture, Phase A compatible)

**Design stance, one sentence:** one canonical IR as the only stable contract, all frontends (Go SDK, YAML artifacts, future physics language) desugar into it, a deterministic pass suite produces verdicts, and closed constructs with typed signatures shrink — but never eliminate — the label-injection surface.

The pasted document's core principles are taken as constraints, not suggestions: the MVP is not discarded (it becomes the semantic backend), the language is **earned from corpus evidence** rather than invented, and the language encodes physics semantics, not just math syntax. Everything below is that program made concrete in Go.

---

## 1. The pipeline

```text
Frontend 1: Go SDK (struct literals)      ┐
Frontend 2: YAML/JSON hypothesis artifact ├─→ desugar ─→ canonical IR
Frontend 3: physics language source (v0)  ┘                │
                                              │
                                    SEMANTIC PASS SUITE
                                    (deterministic, ordered)
                                              │
                     ┌────────────────────────┼──────────────────────┐
                     ▼                        ▼                      ▼
              Verdict artifact          Lowered execution IR    Lean obligations
              REJECTED / WF+flags /     → Program B             → proof backend
              DEFERRED (+ spec stamp)   (B-test specs, nulls)   (derivation claims)
```

Three decisions frozen here, carrying forward the earlier review findings:

1. **Edge→graph classification is mechanical** (relation registry decides DERIVATIONAL vs PHYSICAL; the author never classifies edges directly).
2. **One canonical serialization** — the IR's ground truth is canonical JSON (sorted keys, rational numbers only). Go structs, YAML, and the future language are authoring surfaces, never storage formats.
3. **No floats in the claim IR.** Numeric approximation is Program B's business; the compiler's world is exact and symbolic.

---

## 2. Package layout

```text
physcomp/
  cmd/physc/                 # CLI: physc check|lower|emit
  internal/
    sdk/                     # Frontend 1: Go SDK → IR (the MVP's path, preserved)
    artifact/                # Frontend 2: YAML/JSON → IR
    lex/ parse/ ast/         # Frontend 3: physics language (Phase B, built last)
    ir/                      # canonical claim IR + content addressing
    expr/                    # expression AST: types, dimensions, operator signatures
    spec/                    # spec loader, fact typing, versioning
    sema/                    # the pass suite (Program A's checker, restructured)
    construct/               # closed constructs with enforced signatures
    lower/                   # lowering to execution IR / obligations
    diag/                    # diagnostics + repair directives
    hash/                    # canonical serialization, content addressing
  spec/semantic-spec.yaml    # single authoritative spec, versioned
  testdata/golden/           # corpus with .want verdict files
```

The MVP's `program-a/` checker does not die — its rules become `sema` passes over `ir`, its `semantic-spec.yaml` stays exactly where it is, and the Go SDK stays as Frontend 1. Nothing is rewritten; it is re-homed.

---

## 3. Canonical IR

```go
package ir

type ID string // H(canonical JSON) — content-addressed

type Unit struct {
    SpecVersion    string  // "semantic-spec@v0.3.1" — stamped on every verdict
    GrammarVersion string  // "physlang@v0"; empty for SDK/artifact frontends
    Source         FrontendKind
    Nodes          []Node
    Edges          []Edge
    FactsUsed      []FactUse // audit trail: which spec facts were consulted
}

type Node struct {
    ID       ID
    Kind     Kind        // CONCEPT | CLAIM | BRIDGE | CONSTRUCT_APP
    Name     string      // "physics.causal_order"
    MathType string      // "PARTIAL_ORDER"
    Dim      expr.Dim    // rational-exponent dimension vector
    Role     Component   // DANCER | DANCE | STAGE | JOINT
    Regime   Regime      // declared regimes + escalation rules
    Prov     Provenance  // BASELINE | POSTULATE | DERIVED | THEOREM_REF | UNVERIFIED
    MRC      MRC         // math-reality-check metadata (physical_job, work_level)
    Repr     []ID        // representation edges: physical node → math structures
}

type Edge struct {
    From, To ID
    Rel      string      // resolved through the registry below
    Mode     Mode        // for CLAIM nodes: DERIVE | REPRESENT | IDENTIFY | ...
    Bridge   *ID         // mandatory for cross-qualifier DERIVES/IDENTIFIES
}

// Frozen decision: the registry classifies edges, never the author.
var Registry = map[string]GraphClass{
    "ASSUMES": DERIVATIONAL, "REPRESENTS": DERIVATIONAL,
    "DERIVES": DERIVATIONAL, "IDENTIFIES": DERIVATIONAL,
    "INTERPRETS": DERIVATIONAL,
    "CONSTRAINS": PHYSICAL, "COUPLES": PHYSICAL,
}
```

`FactsUsed` is the audit trail that makes `UNKNOWN_TO_A` honest: every verdict records which spec facts it leaned on, so a later spec correction can identify which prior verdicts are invalidated.

---

## 4. The type system — five layers

This is the heart of the compiler, and the answer to "physics semantics, not math syntax":

| Layer | Example | Checked how |
|---|---|---|
| 1. MathType | `TENSOR`, `PARTIAL_ORDER` | structural compatibility |
| 2. Dimension | `L T⁻¹` | rational-exponent algebra |
| 3. Component/Role | `STAGE` + causal-structure role | role compatibility table (advisory, per the F4 fix) |
| 4. Nominal physics qualifiers | `PhysicalTime` ≠ `RGScale` ≠ `CoordinateTime` | **crossing requires a bridge constructor** |
| 5. Provenance | `BASELINE` vs `DERIVED` vs `UNVERIFIED` | K10 discipline: provenance can never silently upgrade |

Layer 4 is where the earlier `CausalOrder ↛ PhysicalTime` problem gets its correct form: these are **distinct nominal types**, and there is no implicit conversion — but there are explicit **conversion constructors**, which are typed bridges with mandatory obligations:

```go
var BridgeRGToTime = construct.Bridge{
    From: "physics.RGScale", To: "physics.PhysicalTime",
    Requires: []construct.Obligation{
        {Kind: "B_TEST", Spec: "observables in T_RG match H4 scale-factor time in overlap regime"},
        {Kind: "NULL",   Spec: "arbitrary monotone reparametrization of k must NOT match equally well"},
    },
}
```

This is the F4 resolution carried into the compiler: the identification is not banned by grammar and not permitted by labeling — it is **a typed operation whose obligations are structurally mandatory**. H0/H0′ from the stress test fail here before any semantic pass runs: there is no bridge constructor from bare `PARTIAL_ORDER` to `PhysicalTime` with zero obligations.

---

## 5. Closed constructs vs hypothesis blocks — the injection-surface answer

The compiler's key structural move: **frequent derivation chains stop being free-floating `DERIVES` edges and become typed constructs whose signatures encode their semantic prerequisites.**

```go
package construct

type Port struct {
    Name       string
    MathType   string
    Dim        expr.Dim
    Constraint func(*ir.Unit, ir.ID) *diag.Diagnostic // semantic prerequisite
}

type Signature struct {
    Requires    []Port
    Produces    []Port
    Obligations []Obligation // B-tests and nulls — structurally mandatory
}

type Construct struct {
    Name string
    Sig  Signature
    App  func(inputs []ir.ID) []ir.Node // builds IR nodes ONLY — never evaluates
}
```

The Noether example from the pasted document, fully typed:

```go
var Noether = construct.Construct{
    Name: "noether",
    Sig: construct.Signature{
        Requires: []construct.Port{
            {Name: "action",   MathType: "ACTION"},
            {Name: "symmetry", MathType: "CONTINUOUS_SYMMETRY",
             Constraint: symmetryOfAction}, // the symmetry must be OF that action
        },
        Produces: []construct.Port{
            {Name: "current", MathType: "VECTOR_FIELD"},
            {Name: "charge",  MathType: "SCALAR"},
        },
        Obligations: []construct.Obligation{
            {Kind: "B_TEST", Spec: "∂_μ j^μ = 0 on-shell"},
            {Kind: "NULL",   Spec: "broken-symmetry control: conservation must fail"},
        },
    },
}
```

Now `Group DERIVES ConservationLaw` is not merely flagged — it is **ill-typed**: no construct and no bridge has that signature, and free `DERIVES` edges targeting layer-4 qualifier crossings are rejected for lack of a bridge. Compare with the MVP, where that claim was author-labeled and only *flagged*. This is the concrete mechanism by which the compiler improves on Program A: **the injection surface shrinks from "everything is a label" to "everything outside the construct library is a label."**

The escape hatch is first-class and permanent — novel physics must never be blocked by a missing construct:

```text
hypothesis B2_clock_universality {
    mode:     IDENTIFY
    from:     physics.causal_order
    to:       physics.physical_time
    bridge:   clock_universality
    regime:   semiclassical
    null:     path-independent-clock control
}
```

Hypothesis blocks route through the full mode/masquerade machinery and receive flags. Missing constructs are a **DEFERRED**, never a REJECTED.

**First construct candidates** (promoted from the stress-test corpus, not invented): `Clock(C, Φ)`, `Noether(A, G)`, `ConditionalDynamics(H, Clock)` (Page–Wootters), `Guidance(j, ρ)`, `PathSum(Templates)`, `RGFlow(Γ_k)`.

---

## 6. The pass suite

```go
package sema

type Pass interface {
    Name() string
    After() []string
    Run(*ir.Unit) []diag.Diagnostic
}

var Suite = []Pass{
    SchemaPass{},     // E1xxx: shape, required fields
    ResolvePass{},    // E2xxx: names exist, packs resolved
    TypePass{},       // E3xxx: MathType compatibility, construct signatures
    DimensionPass{},  // E4xxx: dimensional algebra
    RegimePass{},     // E5xxx: cross-regime DERIVES → escalation obligation
    ModePass{},       // E6xxx: declared mode vs actual edge structure (masquerades)
    FactPass{},       // E7xxx: spec-fact lookup, fact-typed (MATH/PROVENANCE/PHYSICS)
    CyclePass{},      // E8xxx: two-graph rules; physical cycles legal
    BridgePass{},     // E9xxx: bridge presence AND obligation shape (anti-ceremony)
    ARIPass{},        // E10xx: presence + role-edge coverage, not linear chain
    VerdictPass{},    // aggregation — always last
}
```

Verdict aggregation is one frozen rule (the P3 fix):

```text
any E-level diagnostic            → REJECTED
else unresolved DEFERRED obligation → DEFERRED
else                               → WELL_FORMED + flags
```

Every flag carries a machine-readable obligation marked `UNOWNED_UNTIL_CLAIMED` — the compiler emits obligations, it does not service them. `App` in constructs builds IR nodes only; evaluation is structurally impossible inside the compiler. That is the A/B boundary enforced in the type system, not in prose (the F5 fix).

---

## 7. Diagnostics and repair

Modeled on the Go compiler: never stop at the first error, deterministic order, every error carries a repair directive.

```go
type Diagnostic struct {
    Code    string     // "E6103_MODE_MASQUERADE"
    Span    Span       // source location or IR node path
    Node    ir.ID
    Msg     string
    Repairs []Repair   // allowed patterns; a repair NEVER auto-accepts
}
```

Example output for the stress-test control case:

```text
E7210 TARGET_ALREADY_PRESENT_IN_ASSUMPTIONS
  node:  physics.growth_step_n (LABEL_TIME)
  edge:  IDENTIFIES → physics.physical_time
  fact:  LABEL_TIME HAS_PROPERTY ORDERED_SEQUENCE [PROVENANCE]
  repair: none — restructure: physical content lives in ensemble measure μ, not in n
```

---

## 8. Lowering and backends

The compiler emits three artifact kinds, and owns none of the downstream:

```go
type Backend interface {
    Name() string
    Accept(lower.Unit) bool
    Emit(lower.Unit) (Artifact, error)
}
// "program-b"        → execution spec: operations, required observables, nulls, kill conditions
// "lean-obligation"  → proof obligations for DERIVATION_CLAIM nodes (never a homegrown prover)
// "verdict"          → verdict artifact for Solvent/EBP (emitted, not maintained)
```

The Lean backend is the structural answer to the "don't build a CAS/theorem prover" rule: derivation claims with proof-shaped obligations are **lowered out** of the compiler entirely. The compiler is a dispatcher, not a mathematician.

---

## 9. Language discovery protocol — how the grammar is earned

Language v0 is not designed; it is **promoted from corpus evidence**:

```text
run 10–20 genuinely different hypotheses through Frontends 1/2
        ↓
log: construct frequency, recurring chain shapes, MRC flag sites,
     UNRESOLVED_BY_STATIC_ANALYSIS sites, what Program B requests
        ↓
recurring chain + formalizable signature + adversarial golden cases exist
        ↓
CONSTRUCT CANDIDATE → corpus-gated promotion → stdlib v(n+1)
```

Construct admission test (this is the anti-"LaTeX with a parser" gate):

1. The signature must encode semantic prerequisites (constraints on inputs, not just types).
2. It must be corpus-justified (seen in ≥N real hypotheses), not anticipated.
3. It must ship with adversarial golden cases (misuse attempts, not just correct uses).
4. It must contain no hardcoded physics conclusions — physics bets live in obligations, never in signatures.

The grammar itself is versioned (`GrammarVersion` in every Unit) and its changes are spec-grade events: full corpus regression before merge. A grammar edit is a hypothesis-grade event, governed like one.

---

## 10. Determinism engineering

- **Content addressing everywhere:** `ID = H(canonical JSON)`; Unit hash = H(unit ∥ spec-version ∥ suite-version). Same snapshot + same versions → same verdict, byte-for-byte.
- **Sorted maps, no global state, no goroutines in the pass suite.** Parallelism, if ever, happens at the corpus level, not inside a check.
- **Rationals only** in `expr.Dim` and coefficients.
- **Go-native fuzzing** on the parser and the artifact desugarer (`go test -fuzz`) — adversarial input is expected, so fuzzing is part of the definition of done.
- **Golden corpus with `.want` files** in `testdata/golden/`, including the adversarial categories: relabeling attacks, import-as-derivation, ceremonial bridges, label-time smuggling (H0/H0′), legitimate physical cycles.
- **Double implementation** of the two most dangerous passes (`CyclePass`, `FactPass`): a second naive implementation in test code, compared on the full corpus.
- **Self-hosting from day one:** `spec/semantic-spec.yaml` is itself fed through the checker as corpus case #1. The spec is checked by the thing it specifies, every time it changes.

---

## 11. What this design does not claim (carried forward, restated for the compiler)

1. **Typed constructs narrow the injection surface; they do not close it.** Hypothesis blocks remain author-declared, and undeclared semantic escalation inside them remains the adversarial layer's responsibility. The compiler's guarantee is: *everything that passes through a construct or bridge signature is structurally honest; everything else is flagged, not guaranteed.*
2. **Verdicts are unchanged:** REJECTED / WELL_FORMED+flags / DEFERRED. WELL_FORMED means formulable, never true.
3. **The compiler proves nothing, computes nothing, adjudicates nothing.** It validates formulation under a versioned spec, lowers obligations to parties who can discharge them, and records which facts it leaned on so its own verdicts are correctable.

---

## 12. Build order

| Step | Deliverable | Gate |
|---|---|---|
| 1 | Extract canonical IR + content addressing from MVP | MVP corpus byte-identical through new path |
| 2 | Pass framework + diagnostics + verdict aggregation | all 30 golden cases pass with `.want` files |
| 3 | SDK + artifact frontends desugar to IR | 10 fresh hypotheses round-trip |
| 4 | `expr` layer: dimensions, operator signatures, rationals | dimension errors caught on corpus |
| 5 | First 3 constructs (`Clock`, `Noether`, `ConditionalDynamics`) + bridge constructors | H0/H0′ rejected at type level, not flag level |
| 6 | Lowering stubs: Program-B spec, Lean obligations, verdict artifact | one hypothesis end-to-end |
| 7 | Corpus expansion + construct promotion pipeline | first stdlib v0.2 promotion from evidence |
| 8 | Language grammar v0 **from the corpus**, Frontend 3 | grammar versioned, self-hosted, fuzzed |

Steps 1–3 are restructuring of working code, not new science. The genuinely new machinery begins at step 4, and it is exactly the machinery the MVP deferred: the expression layer and typed constructs — built only now, because the corpus now exists to justify their signatures.

**The one-line summary:** a Go physics compiler is a *hub-and-spoke*: canonical content-addressed IR at the hub; disposable frontends; a deterministic pass suite as Program A's semantic backend; typed constructs as the mechanism that converts recurring derivation chains from author claims into compiler-checked signatures; and lowering — never proving, never computing — for everything downstream.

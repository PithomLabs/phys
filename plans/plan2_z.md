# The Compiler, Re-Founded on Go Principles

**The governing correction first.** The previous design made a category error of its own: it built *conceptual* ideas (jurisdiction law, precedent citation, adversarial routing, obligation adjudication) into compiler machinery. Those are mental models for humans — methodology, not grammar. Under the new principles, the partition is ruthless:

> **The compiler enforces syntax, types, and interfaces — nothing else. It renders intent into artifacts and diagnostics; analysis tools examine artifacts; humans interpret meaning. The compiler has no opinions.**

This is not a loss of rigor. It is where rigor actually lives: Go's compiler is powerful precisely because it enforces only what is mechanically decidable, and expresses everything else through *library API design, vet heuristics, and documentation*. We steal that entire philosophy.

---

## 1. The partition: conceptual vs compiler-native

Every construct from the previous designs, dispositioned:

| Concept | Disposition | Compiler-native form (if any) |
|---|---|---|
| Dancer / Dance / Stage | **Human mental model** | *Package membership.* `std/geom` types are "Stage"; `std/qm`, `std/mat` are "Dancer"; `std/dyn` is "Dance". The compiler checks types, never metaphysics. |
| Reality-First, no-sacred-cow | **Methodology** | One mechanical residue: the physics-to-math authority direction becomes *declaration-before-use* (§3). |
| Jurisdiction map / silence zones | **Documentation** | *API surface.* `alg.LieAlgebra` has no `.Probability()` and no `.Distance()` — silence is an absent method (§4). |
| Masquerade taxonomy | **Half compiler, half vet** | The grammar-expressible half is a compile error (missing premise edge, interface unsatisfied). The semantic half is `phys vet` warnings. |
| K09/K10 (time/Born import) | **Type system** | Provenance types: `Derived[T]`, `Assumed[T]`, `Proved[T]` — no implicit conversions (§3). |
| Two graphs (derivation vs physical) | **Type system** | Go's import-cycle ban + runtime values (§3.3). |
| Obligations, nulls, HIGH_OBLIGATION | **Removed from compiler** | Required *arguments* of constructors (providing the object is grammar; judging it is human). Unproven `Derived` nodes are reported, not adjudicated. |
| Precedent cases, KB-as-law | **Documentation** | Doc comments on std packages, rendered by `phys doc` (pkg.go.dev, not case law). |
| Adversarial routing, EBP, verdicts beyond pass/fail | **External** | Not the compiler's job, explicitly. |
| MRC work levels | **Vet hints** | "Unproven derivation relied on by downstream claims" — a warning, not a verdict. |
| Regimes | **Declarations** | `regime` annotations on derivations; *cross-regime use* is vet; regime presence is compile-checked. |

The deep point of the partition: **every previous "semantic check" reclassifies as one of four things — a type rule, a vet heuristic, a document, or an external tool.** Nothing is lost except false authority.

---

## 2. The Go-compiler mapping (the core design theorem)

Each Go mechanism is load-bearing for exactly one physics-grammar requirement:

| Go mechanism | Physics requirement it enforces |
|---|---|
| **Import-cycle ban** | Derivation acyclicity — *for free, structurally* |
| **Declaration-before-use** | Premises precede conclusions; anti-circularity as "undefined: τ" |
| **"Imported and not used"** | Assumption minimality — uncited premises are *errors*, not bookkeeping |
| **"Declared and not used"** | Ornamental derivations flagged |
| **Interfaces (structural)** | Capability grammar: what a structure can *do* |
| **Nominal types** | Physical qualifiers: `PhysicalTime` ≠ `RGScale` as distinct types |
| **No implicit conversions** | Provenance and qualifier discipline |
| **Generics** | Parametric structure: `Tensor[D, F]`, `Hilbert[S]` |
| **Embedding** | Building algebraic structures (Clifford embeds a vector space + product) |
| **Export data** | Hypothesis artifacts: downstream never re-parses source |
| **`go vet`** | Masquerade heuristics, dimensional suspicion, regime crossings |
| **pkg.go.dev** | The KB as *documentation of std*, not law |
| **`go fix` + compatibility promise** | Spec-version migrations, artifact stability |
| **MVS** | Framework/vocabulary version selection |
| **gopls** | Live typechecking for the AI authoring loop |

Two of these deserve expansion because they *are* the design.

### 2.1 Two graphs = imports + values

This is the cleanest formalization yet of the derivation/physical split, and it is pure Go:

- **`derive` blocks are packages.** A derivation *imports* its premises. Go bans import cycles: if derivation A needs B's result and B needs A's, that is `E-cycle`, exactly like Go's `import cycle not allowed`. No cycle checker to write — the module system *is* the derivation DAG.
- **Couplings are values.** Matter ↔ Geometry backreaction is not a dependency of one on the other; it is a `dyn.System` *value* referencing both. In Go, types may reference each other at runtime while packages never import cyclically — the identical move makes backreaction, self-consistency, and fixed points legal while derivation circularity is impossible.

The two-graph machinery of previous designs collapses into: **compile-time = acyclic; runtime = anything.**

### 2.2 Provenance is a type parameter

```go
type Assumed[T] struct{ decl Ref }          // from postulate/assume
type Derived[T] struct{ from []Ref }        // only a derive block may produce this
type Proved[T]  struct{ cert CertID }       // only an external proof ref may produce this
```

`Duration()` returns `Derived[Interval]`. There is no function `(Poset) → Derived[PhysicalTime]` in std — so H0 is not flagged, it is **a type error**. And because a `derive` block cannot reference values declared after it (declaration-before-use, already Go behavior), a "derived" time that silently presupposes itself fails as `undefined`, not as a deep semantic audit.

---

## 3. The language and the spec

Following Go's own method: **the spec is a normative prose document** (lexical structure, declarations, expressions, derivation statements, regime clauses, package system, compatibility), the compiler implements it, and conformance tests enforce implementation-of-spec. We do not invent a formal semantics for the spec in v0 — Go didn't either; the spec + conformance corpus *is* the semantics contract.

### 3.1 The physics grammar = a small interface set

The words in the spec that carry physics meaning are interfaces, each with spec-defined method semantics:

```go
type Spatial   interface{ Distance(Spatial) Interval }
type Temporal  interface{ Duration() Interval }
type Observable interface{ Spectrum() Set[Real] }
type Dynamical interface{ Flow(param Parameter) Transform }
type Stateful  interface{ Expect(O Observable) Real }
```

To call something *time* in this language is to claim it implements `Temporal` — i.e., to supply its duration semantics. That is grammar, not adjudication: the language defines what its own vocabulary words mean, mechanically. `cst.CausalSet` does not implement `Temporal`; `ms.Clock` does. H0′ ("growth step n is time") fails as `growthStep does not implement ms.Temporal (missing method Duration)` — the classic Go error message, now doing philosophy's job.

### 3.2 Source sketch — the canonical stress test, compiled

```go
package pothour

import (
	cst "frameworks/cst"     // vocabulary only; carries no physics authority
	ms  "std/measure"
	qm  "std/qm"
)

postulate C    cst.CausalSet          // declared, not derived
postulate Φ    qm.Field[over = C]

define  clock  = ms.Clock{Order: C, Matter: Φ}   // signature *demands* both arguments
derive  τ      ms.Interval = clock.Duration()    // Derived[Interval], premises recorded
```

Rejects the compiler *must* produce:

```text
./pothour.ph:9:14: E0212: cannot derive ms.Interval from cst.CausalSet:
    cst.CausalSet does not implement ms.Temporal (missing method Duration)
    suggestion: ms.Clock{Order: C, Matter: Φ} implements ms.Temporal
                given postulates C, Φ — did you mean to derive via clock?

./pothour.ph:12:3: E0104: assumed and not used: Φ
    (a derivation citing only C cannot produce a Temporal result)
```

Note what happened: **the role of math is enforced by API surface.** Silence zones became absent methods; the "bridge" became a constructor whose signature *requires* the matter sector — you cannot even express "order alone yields duration" without the compiler pointing at the missing argument. Import-as-derivation dies as a missing-premise type error; no law library needed.

### 3.3 The backreaction control

```go
package backreact

import (
	geom "std/geom"
	mat  "std/mat"
	dyn  "std/dyn"
)

postulate g   geom.Metric
postulate Φ   mat.ScalarField[on = geom.Manifold.Of(g)]

couple   S    = dyn.System{Stage: g, Matter: Φ}   // a value, not a dependency
derive   δg   geom.Perturbation = dyn.Backreaction(S)
```

`couple` introduces a runtime edge between values; `derive` imports only from prior declarations. The cycle the earlier designs needed a special checker to permit is now just *two values in one scope*, and the cycle the checker had to forbid is structurally unrepresentable.

---

## 4. Core standard libraries — jurisdiction as API

The std is where "the role of math relative to grammar" actually lives. Small, curated, and **math-only**: physics frameworks (BM/CDT/CST/AS) are *user-level* modules above std, versioned separately, carrying no normative weight — the no-sacred-cow policy, expressed as a module boundary.

| Package | Encodes | Notably *absent* (the silence, as API) |
|---|---|---|
| `std/dim` | Rational-exponent dimensions; typed constants require units | — |
| `std/alg` | Groups, Lie brackets, algebras, representations, fusion | `.Probability`, `.Distance`, `.Flow` |
| `std/topo` | Spaces, invariants, classification | `.Dynamics`, `.Scale` |
| `std/geom` | Manifolds, metrics, distance, curvature | `.Probability` |
| `std/an` (analysis) | Limits, integration, differential equations, flows | `.Genesis` (no initial-condition selection) |
| `std/measure` | Probability spaces, clocks, intervals, expectation | — |
| `std/qm` | Hilbert spaces, states, operators, observables | `.Born` — probabilities come from `measure`, never from algebra |
| `std/dyn` | Actions, generators, constraints, coupled systems | — |
| `std/inf` | Information, entropy (coarse-graining is a declared argument) | — |

Doc comments on these packages carry the KB content — what each structure is *for*, its known limits, its boundary cases — rendered by `phys doc`. Humans consult it; the compiler never reads it. **The knowledge base is demoted from law to pkg.go.dev, and that is the correct rank for prose.**

---

## 5. Compiler pipeline and diagnostics

Standard Go-compiler phases, reused nearly unchanged:

```text
lex → parse → AST
  → pass 1: declarations (postulates, defines, imports, regimes)   [two-pass, types2-style]
  → pass 2: bodies (derive blocks, expressions)
  → phys IR (unified-IR analogue; nodes carry provenance + spans)
  → export data (content-addressed artifact: decls, premises, results, regimes, spans)
```

Diagnostics follow Go discipline exactly — never stop at the first error, span-precise, suggestion-bearing:

```text
E0212  interface not implemented (capability grammar)
E0104  assumed and not used          ← premise minimality
E0106  declared and not used         ← ornamental result
E0308  dimension mismatch: L·T⁻² vs L (mismatched exponent T)
E0415  undefined: τ (premise declared after use — derivation cycle)
E0503  implicit conversion Derived[PartialOrder] → Derived[PhysicalTime] undefined
E0601  derivation import cycle: A → B → A
E0702  regime annotation missing on cross-framework import

V2031  vet: type named in declaration but no method of it appears in any derivation
       (notation-as-structure)
V2044  vet: result consumed outside its declared regime
V2051  vet: derivation returns T where caller treats as U — nominal qualifiers differ
```

Vet never blocks. The compiler *fails* only on E-codes — that is the entire enforcement budget, and it is spent only on what is mechanically decidable.

---

## 6. Toolchain — representing intent, analyzing output, never interpreting

Principle 5, concretized as binaries (the `go` command pattern: orchestrator + tools; the compiler proper never runs anything):

| Tool | Analogue | Function |
|---|---|---|
| `physc` | `compile` | diagnostics + export artifact; exit 0/1 |
| `phys vet` | `vet` | suspicious-pattern warnings over export data |
| `phys doc` | `pkg.go.dev` | renders structure: dependency graph, assumption ledger, unproven-`Derived` report, regime map — **the human's analysis surface** |
| `phys types` | `go/types` | embeddable type-checker library — the AI authoring loop calls it live (gopls-style), so hypotheses arrive pre-typechecked |
| `phys fix` | `go fix` | mechanical migration across spec versions |
| `phys prove` | (external) | ships `Proved[]`-shaped obligations to Lean/proof backends — like `go test` invoking external binaries |
| Program B link | cgo/external linking | B consumes export data only; never source |

The intent-representation claim is now precise: **the source file is the intent; the type checker is the conversation about that intent; `phys doc` is the intent made inspectable; nothing in the toolchain asserts what the hypothesis means physically.** Analysis is a report; interpretation is a human act.

---

## 7. What this design explicitly no longer does

Removed relative to the previous iteration — listed so the removal is deliberate, not silent:

1. No jurisdiction verdicts, precedent law, or obligation adjudication inside the compiler.
2. No adversarial routing, no WELL_FORMED-as-endorsement — exit status is *compiles / does not compile*, plus vet warnings; "DEFERRED" becomes "compiles, contains unproven `Derived` nodes" (a doc report, not a verdict).
3. No KB decision table — the KB is doc comments.
4. No component metaphysics (Dancer/Dance/Stage) in the type system — package membership conveys it to humans.
5. No claimed semantic guarantees beyond the type system's reach — the honest scope line: *the compiler guarantees the artifact is grammatically well-formed, premise-complete, cycle-free, and provenance-honest; every deeper question is a tool report or a human judgment.*

---

## 8. Build order (delta only)

1. Spec v0 prose + conformance corpus (the spec *is* the first artifact of the project).
2. `phys/types` two-pass checker + diagnostics — the reusable core, built first so the AI authoring loop exists from day one.
3. `std/dim`, `std/alg`, `std/measure`, `std/qm`, `std/dyn` with interface grammar — the five packages the stress tests need.
4. `physc` + export data; port the two canonical stress tests and the H0/H0′ rejects as conformance cases (they are now *compile errors with exact messages*, the strongest golden cases the project has had).
5. `phys vet`, `phys doc`; frameworks module; `phys fix` deferred until the first spec revision actually happens.

**One-sentence summary:** a Go-philosophy physics compiler — spec-first, two-pass typed, cycle-banning imports for derivations, values for couplings, interfaces as the physics vocabulary, provenance as types, std packages as the executable form of "where math applies," vet for suspicion, docs for knowledge, tools for analysis — and no opinion anywhere about what the physics means.

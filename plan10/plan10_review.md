I reviewed Plan 10 against the **current v2.3 specification and finalized Plan 9.2**. My conclusion is:

**Plan 10 is not materially better architecturally. It is mostly a more compact, execution-oriented restatement of the same architecture. However, as written, it is actually less safe than Plan 9.2 because it is still anchored to v2.2 and contains several regressions against v2.3.**

So I would **not refactor Plan 9.2 toward Plan 10**. Your instinct to avoid that risk is sound.

### The big picture

The two plans share essentially the same architecture:

```text
internal/kernel
      ↓
core
      ↓
ops
      ↓
session

mechanics / relativity / hypothesis
      ↓
core + internal/kernel
```

Plan 10 explicitly uses that same graph and the same 39-file tree. 

Plan 9.2, however, explicitly declares **v2.3 as the sole normative authority**, whereas Plan 10 opens by saying its normative source is **`plan9.md` / v2.2**. That alone makes Plan 10 unsuitable as the canonical handoff in its current form.  Plan 9.2 has the correct v2.3 authority chain. 

### What Plan 10 genuinely does better

There are a few useful things in Plan 10, but they are mostly **execution/readability improvements**, not architectural improvements.

| Plan 10 feature                                    | Adopt into 9.2?      | Why                                                             |
| -------------------------------------------------- | -------------------- | --------------------------------------------------------------- |
| 16-step dependency-ordered implementation sequence | **Yes, optionally**  | Very agent-friendly: purpose → deps → files → result → gate     |
| Explicit test-package placement                    | **Yes**              | Good clarification of internal vs external package tests        |
| Very explicit acceptance-test inventory            | **Yes, selectively** | Useful operational detail, though 9.2 already covers most of it |
| Strongly explicit final build/vet/test gates       | **Yes, selectively** | Good handoff discipline                                         |
| Expected dependency/import inventory               | **Yes, optionally**  | Useful verification checklist                                   |

But importantly, **Plan 9.2 already contains most of these ideas**. Its own implementation phases already identify dependencies, files, results and verification gates, so importing Plan 10 wholesale buys relatively little. 

### The material problems in Plan 10

These are the reasons I would **not use Plan 10 as-is**.

#### 1. It is still based on v2.2

This is the biggest problem.

Plan 10 literally starts with:

> “Normative source: `plan9.md` = specs v2.2 ...”

while the governing specification is now v2.3. 

That matters because v2.3 did not merely rename things. It introduced three explicit reconciliations, including the normative `ops.Apply` positional schema. Plan 9.2 incorporates those reconciliations explicitly. 

So Plan 10 is not a clean competing implementation of the same contract.

#### 2. It misses the v2.3 `ops.Apply` positional contract

This is a particularly important omission.

Plan 9.2 explicitly has the complete mapping:

```text
add          inputs[0]=a inputs[1]=b
subtract     inputs[0]=a inputs[1]=b
...
substitute   target, variable, replacement
differentiate target, wrt
limit        target, variable, value
compare      a, b + Operator
solve        relation, target
select_branch branches, constraint
```

and explicitly says `Session.Step` and replay must use those positions. 

Plan 10 describes the 12-operation switch, but its matrix does **not incorporate the v2.3 §15.13.1 positional mapping**. That means it is missing one of the reasons v2.3 exists.

That's a strong reason to keep 9.2 as the canonical plan.

#### 3. Challenge/Review handling regresses

Plan 10 says:

> “`Challenge/Review` re-declared as aliases of `core` types”

in its session implementation sequence. 

The current v2.3 contract is more precise: `Challenge` and `Review` are defined **exclusively in `core/corpus.go`**, and `session` consumes them **without redefining them**. 

An alias is technically different from defining a second struct, so this isn't the catastrophic kind of import-cycle error we previously rejected. But it is nevertheless an unnecessary public surface deviation from the settled v2.3 contract.

Plan 9.2 has already made the cleaner choice. 

**Do not adopt this from Plan 10.**

#### 4. The `schema_version = "1"` pin is an unnecessary new contract

Plan 10 says:

> `schema_version` pinned value: `"1"` for every artifact schema in MVP

and extends that to object, ledger, session, candidate and manifest. 

That is not a v2.3 requirement.

More importantly, Plan 9.2 deliberately resolved the related architectural ambiguity differently: `schema_version` belongs to the **canonical serialization representation**, not the authoritative kernel `Object` state. Its canonical codec rules explicitly avoid adding a serialization-only field to `kernel.Object`. 

So this is exactly the sort of extra contract I would avoid introducing during implementation.

**Do not adopt.**

#### 5. Plan 10 exports an internal object decoder

Plan 10 specifies:

```go
func DecodeObjectJSON(data []byte) (Object, error)
```

inside `internal/kernel`. 

This isn't automatically architecturally disastrous because the package is under `internal/`, but Plan 9.2 deliberately hardened this boundary so object canonical decoding is **internal/unexported**, with no public `ParseObject`/`DecodeObject` surface. 

Again, Plan 10 creates API surface where 9.2 intentionally doesn't.

**Do not adopt.**

#### 6. Plan 10's Test D contains a real dimensional mistake

This one is concrete.

Plan 10 says its equal-dimension kind-mismatch test can use:

> "`KineticEnergy` vs `Momentum`-family fixture"

But kinetic energy has dimensions of energy, while momentum does not. Therefore that is not a clean test of **kind mismatch with equal dimensions**. 

Plan 9.2 already fixed this by using an Energy/KineticEnergy-style same-dimension pair. 

**Definitely do not import Plan 10's Test D.**

### One important thing Plan 10 confirms

The fact that Plan 10 independently converged on:

* the same 39-file tree,
* the same kernel/core/ops/session layering,
* the same pure `ops`,
* the same session-owned ledger,
* the same bounded symbolic engine,
* the same candidate containment,
* the same 16-ish implementation stages,

is useful evidence that **the core architecture is stable rather than an artifact of one agent's plan**.

In other words, the independent plan is validating the architecture more than challenging it.

### My recommendation for your two-agent strategy

I would **not merge Plan 10 into Plan 9.2**.

I would use:

```text
                 v2.3 normative spec
                         │
             ┌───────────┴───────────┐
             │                       │
        Plan 9.2                  Plan 10
        canonical              independent branch
             │                       │
          Gemini                 second agent
             │                       │
        Implementation A       Implementation B
             │                       │
             └────────── compare ────┘
```

But there is one critical condition:

**Do not implement Plan 10 exactly as attached if the purpose is to compare coding-agent quality.**

It should first be converted into a **Plan 10-v2.3** with only these surgical corrections:

1. replace every normative v2.2/`plan9.md` reference with v2.3;
2. remove the session `Challenge`/`Review` aliases;
3. remove the `"schema_version": "1"` pin unless separately desired;
4. restore internal/unexported object decoding;
5. fix Test D to use genuinely equal-dimension, different-kind objects;
6. add the complete §15.13.1 positional `ops.Apply` table;
7. ensure its coverage matrix is against v2.3, not v2.2.

Those are **contract corrections**, not architectural refactoring.

### Final assessment

I'd characterize them this way:

**Plan 9.2:**
Architecture = settled
v2.3 compliance = strong
Handoff safety = strong
Detail = intentionally exhaustive

**Plan 10:**
Architecture = essentially the same
Execution readability = somewhat better/cleaner
v2.3 compliance = weaker as currently written
Risk of accidental implementation drift = higher

So **there is no compelling reason to disturb Plan 9.2**.

The most sensible use of Plan 10 is actually as an **independent second implementation/checklist**, not as a source for refactoring the canonical plan. And because you specifically want to see whether Gemini and another coding agent arrive at comparable working implementations, maintaining the same v2.3 contract for both is much more valuable than making the plans look identical.

  
## Z

give prompt for opencode to revise plan 10 per attached v2.3 specs, consolidate review below where net valid. Adversarial Review: `plan10.md` vs. specs v2.2

## Verdict

**No architectural blockers found. Conditional GO for hand-off** — contingent on fixing one real internal inconsistency (relativity manifest item list), ~5 precision gaps, and spot-verifying plan citations against spec sections not included in the review packet (both documents are truncated; see Scope).

---

## Scope caveat (material)

- `plan10.md` cuts off mid-slice-19 (S) — plan §12–§13 (final gate) and the tail of slice S were **not reviewable**.
- `specs_v2_2.md` cuts off mid-§15.11 — spec §15.11 (rest), §15.12, §16–§41 are **not independently verifiable** here. All plan claims citing those sections are marked ⚪ below.

---

## Blocker scan — what I attacked and why it holds

| Attack vector | Result |
|---|---|
| Can `core` satisfy "MUST expose methods" (§7.0, §8.0, §11.0, §12.0, §13.0) while "declaring no methods on aliased types"? | ✅ Holds. Kernel methods surface through `type X = kernel.X` aliases. Spec's own §4 note anticipates this. |
| Does `core`→`ops` leakage exist anywhere? | ✅ No. Façade constructors/canonical functions are pure kernel delegations; corpus.go has no ops need; Step ordering enforces it. |
| Replay cycle: `session` re-executes `ops.Apply` during `Validate` | ✅ Clean. `ops` knows nothing of ledgers; decoded objects re-enter via `kernel.DecodeObjectJSON`. `go build ./...` + import scan closes it. |
| REQ-003-01 file budget: 39 files, mixed test packages, embed access | ✅ Arithmetic checks out (24 src / 10 test / 2 json / 3 doc); `go:embed` var accessed via internal test package; mixed `pkg`/`pkg_test` per dir is legal Go. |
| Mint-authority circumvention via public API | ✅ Holds. Unexported kernel fields + no core forwarder + REQ-003-03 external tests. In-module minting is correctly framed per REQ-000-03/04. |
| Golden traces (G, H, I, J) executable under the bounded rules? | ✅ I hand-traced all four: product-rule zero-terms collapse under finiteness; `c⁻¹` admitted via `c>0`→positive→nonzero; `Sqrt(Pow(m·c²,2))` unwraps via `m≥0` + even-exponent rule; Lorentz limit reaches `1/1` through the body path, not ID-matching. |
| Constructor-vs-Simplify boundary vs. spec §8.5 ("remove safe identities") | ✅ Defensible. §9.6 preamble is explicit and more specific; plan's reading is the right one. (See 🟡-2 to document it.) |
| Non-goals compliance | ✅ No CAS/prover/parser/numeric/series; `Solve` is one pattern; entailment closed; `statement` never parsed; no `go/parser` in non-test files. |

---

## 🟠 Must-fix before hand-off

**F1 — Relativity manifest item list is internally inconsistent (strongest finding).**
Three places disagree:
- Plan §7: 10 items — `Spacetime, MinkowskiMetric, RestMass, Energy, ThreeMomentum, FourMomentum, SpeedOfLight, LorentzFactor, EnergyMomentumRelation, MassEnergyRelation` (no Velocity).
- Plan §4 pinned names: includes `velocity` in the manifest name list (⇒ 11 items).
- Step 11: "8 wrappers + Velocity wrapper … 10-item manifest" — ambiguous whether Velocity is a manifest item.

Since the plan's selling point is "fixed data pins (all deterministic, no open choices)," an item-list contradiction in the corpus is exactly what causes implementation drift and manifest cross-check failures (slice O hashes and item-count assertions will break one way or the other). Mechanics is consistent (11 = 8 + 3 everywhere); reconcile relativity to one list and one count. ⚪ Confirm against unseen §28.2.

**F2 — `Pow(x,0) → 1` is missing from the plan's Simplify rule list.**
Spec §8.7 states it normatively ("Pow(x, 0) simplifies to 1 only when x is a valid nonzero-safe base under the operation's assumptions"); plan §6's rule list omits it. No golden trace needs it, so this is cheap to fix — add it (Simplify-time, nonzero-safe gate, `0^0` still unsupported) or record an explicit deviation note.

**F3 — Kernel-side metadata constructors are implied but never declared.**
`kernel.Provenance/Assumption/ConventionSet/Dimension` have unexported fields (plan §4), so `core.NewProvenance`, `NewTextAssumption/NewExprAssumption/NewAssumptionSet`, `NewConvention(s)`, and the nine `Dimension*` constructors **cannot be written without kernel-side constructors**. Add them explicitly to Step 1/2 deliverables, and pin that provenance invariant logic lives **once** (kernel), with `core.NewProvenance` delegating — otherwise `MintObject`'s provenance validation and `core.NewProvenance`'s can silently diverge.

**F4 — Canonical element ordering of `AssumptionSet`/`ConventionSet` is unpinned.**
§9.3 sorts only `Add`/`Mul` children. Set serialization order feeds `HashAssumptionSet`/`HashConventionSet` → object hashes → ledger chain → candidate `LedgerHash`. Without a pinned ordering (e.g., kind-ordinal then key, or canonical-bytes), determinism tests (slice R) are luck. Pin it in Step 1.

**F5 — Slice D's example fixture is dimensionally wrong as written.**
"Energy vs Torque-style — KineticEnergy vs Momentum-family" — Energy and Momentum do **not** share dimensions. Valid equal-dim/different-kind pairs from the 18 kinds: `Mass`/`RestMass`, `Velocity`/`SpeedOfLight`, `Position`/`Spacetime`, `Momentum`/`ThreeMomentum`/`FourMomentum`. Fix the example so the test isn't built on a false premise.

---

## 🟡 Clarify / pin (low risk, do before implementation starts)

1. **`OperationParams` schema is opaque.** Step 12's "otherwise `Kind=="empty"`" and §9's "params `Kind:"identify"`" reference a schema the plan never defines. Pin fields, zero-value convention, and canonical JSON (⚪ verify vs. plan9 §15/§16).
2. **Document the §8.5 vs §9.6 resolution** in the plan (one sentence): constructor-time = flatten/combine/sign/sort only; identity removal and all §9.6–§9.8 rewrites are Simplify-time. This also determines the canonical bytes of the unsimplified Limit body (which retains `Mul(1, …)`).
3. **`Simplify` on `Kind=Relation`/`BranchSet` objects**: slice I requires it (Substitute → Simplify on a relation). Pin: kind preserved, sides recursed, `Relation(eq,a,b)` never becomes `IDENTIFIED` (§9.9).
4. **`SelectBranch` failure mode when sign entailment fails**: error vs. unsimplified result — currently unpinned. Also confirm the `Constraint`/`selected_branch/<hash>` assumption keying and the (branchset, constraint) signature against unseen §15.12.
5. **Restate defensive-copy requirements** where ops touch `*big.Rat`: `Pow`/`NewPow` exponent clone (§15.5, §8.0), `NewRational` copy.
6. **State where `Entails*`/`Merge`/`Equal`/`CanonicalJSON`/`Hash` are implemented** (kernel `types.go`) and which step builds them — Step 5 depends on them but Steps 1–2 don't name them.
7. **`schema_version: "1"`** — plan pins it; the shown spec text lists the field but not the value. Either cite the pin or mark it as a plan-level decision.
8. **`DecodeObjectJSON` should route through the same invariant validation as `MintObject`** (not just canonical round-trip) so a structurally-canonical-but-invalid object can't enter replay.

---

## ⚪ Verify against unseen sections before hand-off

Plan citations I could not check because the spec excerpt ends at §15.11:

- §15.11 Solve pattern orientation (LHS-only `Pow(t,2)`? RHS rejection?)
- §15.12 SelectBranch (see 🟡-4)
- §16.2.1 / §16.12 / §16.19 / §16.20 — nine-action surface, exact 16-field Step order, the 17-step `Validate` pipeline, tamper triple
- §19.5-MUST-02 / §19.6 — `EnergyMomentumRelation` assumption carry, `DERIVED`-with-empty-parents, `derivable_from` non-existence tolerance for `RestFrame`
- §26.2 / §26.8 / §26.10 — candidate field order, pipeline order, 14-accessor surface, REQ-032-10a
- §27.3 — eight Review categories
- §28.1/§28.2 — manifest item lists, `source` pins, `EnergyMomentumRelation` = `DEFINED`
- §32 — REQ-032-01..20 existence and wording
- Plan's own §12–§13 (truncated) — final gate content

---

## Recommendation

**Proceed to hand-off after F1–F5 are applied and the ⚪ list is spot-checked.** The architecture itself — alias-mediated kernel methods, single mint path, `internal/` boundary, session-only identification/sealing, test-only import edges, anti-hardcoding source audits — is sound, spec-conformant, and implementable in Go 1.24 stdlib-only within the 39-file budget. Nothing found requires re-planning packages, files, APIs, or the dependency graph; all defects are content-precision issues fixable in a single plan revision.

One process note for the hand-off packet: attach the **complete** `plan10.md` and the **complete** spec, since the truncations above cover precisely the session/candidate/manifest sections where the remaining risk is concentrated.
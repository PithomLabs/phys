## Z 

# Adversarial Review: `plan10.md` vs. specs v2.2

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

## chatgpt


Yes. After comparing the adversarial review against the **complete v2.3 specification and the current Plan 10**, there are several **net-valid findings worth carrying forward**.

The important distinction is that **none of them require architectural replanning**. They are implementation-precision or plan-consistency fixes.

### Net-valid findings

| Finding                                                       | Verdict | Why                                                                                                                                                                                                                                                                                                                                                                                                   |
| ------------------------------------------------------------- | ------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **F1 — Relativity manifest ambiguity**                        | ✅ Valid | v2.3 §28.2 requires exactly the ten named corpus items: Spacetime, MinkowskiMetric, RestMass, Energy, ThreeMomentum, FourMomentum, SpeedOfLight, LorentzFactor, EnergyMomentumRelation, MassEnergyRelation. Plan 10 also lists `velocity` among the “manifest name strings,” while separately describing a Velocity wrapper. The Step 11 wording “8 wrappers + Velocity wrapper” is also ambiguous.   |
| **F2 — Missing `Pow(x,0)`**                                   | ✅ Valid | v2.3 §8.7 explicitly requires `Pow(x,0) → 1` only when the base is nonzero-safe under active assumptions, while `0^0` remains unsupported. Plan 10's Simplify rule inventory does not explicitly include this rule.                                                                                                                                                                                   |
| **F3 — Kernel-side metadata construction**                    | ✅ Valid | The spec requires public constructors for provenance, assumptions, conventions, and dimensions while the underlying concrete types live in `internal/kernel` with authoritative fields. Therefore Plan 10 needs an explicit kernel-side construction mechanism rather than leaving this implicit. v2.3 also makes `NewProvenance` a required public constructor.                                      |
| **F4 — Deterministic ordering of assumption/convention sets** | ✅ Valid | v2.3 requires deterministic set union and canonical artifact field ordering, and these sets participate in canonical serialization and hashing. Plan 10 does not pin how elements within the sets are ordered. That is a genuine determinism gap because those hashes feed object/ledger integrity.                                                                                                   |
| **F5 — Slice D fixture**                                      | ✅ Valid | `KineticEnergy` and `Momentum` do **not** have equal dimensions, so they cannot serve as the intended “equal dimension, different kind” MRC-003 fixture. The plan itself says the test is specifically for equal-dimension kind rejection. Use `Mass`/`RestMass`, `Velocity`/`SpeedOfLight`, or another actually equal-dimensional pair.                                                              |

### The clarifications are also mostly net-valid

**1. `OperationParams` needs to be made explicit — definitely valid.**
This is no longer something that needs verification against unseen material. v2.3 explicitly specifies:

```go
type OperationParams struct {
    Kind          string
    Exponent      string
    Operator      core.RelationOperator
    Justification string
}
```

and §15.13.1 gives the exact positional mapping for all 12 operations. Plan 10 should reproduce that contract explicitly. 

**2. Constructor-time vs. Simplify-time boundary — valid.**
v2.3 is explicit: constructor-time normalization is structural only; the identity, zero, power, and sign-sensitive rewrites are Simplify-time. Plan 10 should state this clearly because it affects canonical bytes and replay behavior. 

**3. Relation/BranchSet simplification — valid as implementation precision.**
v2.3 explicitly requires recursive relation normalization and preservation of the relation operator, and it explicitly forbids converting equality into `IDENTIFIED`. Plan 10 should pin this behavior rather than leaving the generic “Simplify” description underspecified. 

**4. `SelectBranch` failure behavior — valid to clarify, but do not invent a new semantic rule.**
The review correctly identified this as a place where the plan should follow the exact v2.3 contract rather than have the coding agent improvise.

**5. Defensive copies — valid.**
The spec explicitly requires `Values()` and `ParentHashes()` to return copies, and the immutable model makes this important. The same principle should be explicit for `*big.Rat`-backed values.

**6. Ownership of `Entails*`, `Merge`, `Equal`, canonicalization, hashing — valid.**
The spec gives these APIs but Plan 10's implementation sequence could leave the exact ownership implicit. Pinning them to the kernel-backed types removes another implementation choice.

**7. `schema_version: "1"` — valid finding, but lower importance.**
Plan 10 explicitly pins `"1"`, while the v2.3 material we inspected specifies the field and canonical ordering but does not itself establish that exact value as a normative requirement. So Plan 10 should either mark `"1"` as a **Plan 10 implementation pin** or avoid presenting it as something mandated by v2.3.

**8. Decode path should enforce object validity — valid and worthwhile.**
The review's concern is sound: replay input should not merely be syntactically canonical; it must reconstruct a valid immutable object consistent with the same object invariants. The exact phrase “must use the same function” is stronger than the spec requires, but the underlying implementation requirement is legitimate.

### One thing from the adversarial review that is no longer valid as stated

The review's **⚪ “unseen sections” caveat is obsolete**.

That review was performed against a truncated v2.2 packet. We now have the complete v2.3 specification, and the allegedly unseen areas have been confirmed. For example, v2.3 explicitly confirms the **nine session actions** and the **10-item relativity manifest**.  

Also, the review's F1 should be characterized more precisely: it is **not an architectural inconsistency**. It is a **terminology/count ambiguity in Plan 10 that should be eliminated before implementation**.

### Bottom line

The adversarial review was substantively useful. I would carry forward:

**F1–F5 = all net-valid**, with F1 treated as a precision/consistency issue rather than a blocker.

Of the eight yellow clarifications, **all eight are reasonable implementation hardening**, with the caveat that `schema_version: "1"` should be labeled a Plan 10 pin rather than attributed to v2.3, and `DecodeObjectJSON` should be described as sharing the object invariant-validation semantics rather than claiming the spec mandates a particular helper call.

So the correct disposition is:

**Architecture: unchanged.
Plan: revise.
Coding handoff: after these precision fixes, no architectural reopening is warranted.**

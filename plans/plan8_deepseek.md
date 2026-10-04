# FINAL REVIEW — SPECS v2.1 AND PROMPT v2.1

**Status:** Ready for handoff. No architectural blockers remain.

---

## Part I: The Four Previously-Blocking Issues — All Resolved

| Previous Issue | Where Fixed | Verified |
|---|---|---|
| `Differentiate` on domain quantities undefined | §15.8 — bounded derivative engine (constants, symbols, Add, Mul, Neg, integer Pow). No `Derivative` node needed. | ✓ |
| Domain-to-object conversion one-way | §5.1, §5.3 — accessors return `core.Object`; operations consume `core.Object`; no reverse conversion required for MVP. | ✓ |
| File count contradiction | §3 (39 files) vs §38 (cap ≤ 40). Consistent. | ✓ |
| Manifest item → Go constructor mapping | §17.5 `constructor` field + §17.8 static `map[string]func() core.Object` in tests. | ✓ |

The four-blocker list from the previous review has been addressed precisely. No reopening.

---

## Part II: What the v2.1 Spec Adds Beyond v2

| Addition | Value |
|---|---|
| `internal/kernel` boundary | Solves the Go visibility problem: unexported fields in `core.Object` + multiple cooperating packages. Cleaner than the v2 approach. |
| `MintObject(ObjectSpec)` single mint path | MRC-001 is now mechanically enforceable. No second mint helper. |
| `Session.Identify` (not `ops.Identify`) | Removes ambient/global session state. Correct. |
| `Ledger` in `session` (not `core`) | Breaks the `core → ops` import cycle. Correct. |
| `InputCanonicals` + `OutputCanonical` retention | Hashes alone cannot replay. Canonical retention makes replay possible. Correct. |
| `StepEnvelope` canonicalization | Fixes the variable-length field concatenation problem in the hash chain. Correct. |
| Bounded sign entailment (§9.8) | No theorem prover. Exact rules for `Sqrt(Pow(x,2)) → x`. Correct. |
| `BranchSet` as expression node | Bounded container for the two-branch quadratic solve. Correct. |
| `RestFrame` as assumption (§22) | Previously ambiguous. Now clear. |
| Nine session actions with state machine (§16.2) | Deterministic lifecycle. Correct. |
| `DraftMetadata` bundles all candidate metadata | Single-seal path. No hidden mutation actions. Correct. |
| Corpus status `NONE` for derived/candidate artifacts | Separates the two axes cleanly. Correct. |
| Requirement IDs `REQ-*-*` | Enables the coverage matrix to be exhaustive. Correct. |
| `APPROXIMATED` explicitly reserved for v0.5+ | No MVP operation mints it. Correct. |
| Coverage matrix in §39 with `REQ-*` rows | Enables mechanical verification. Correct. |
| Handoff condition §41 | Enumerates the 18 resolved decisions. Correct. |

This is the most complete specification artifact in the project.

---

## Part III: Remaining Minor Ambiguities (Not Blockers)

These will not prevent a competent coding agent from implementing the MVP. They will require small local choices that do not reopen architecture.

### Minor 1 — `Pow(Rational(q), n)` simplification

§9.6 lists the exact mechanical rewrite rules. `Pow(Rational(q), n) → Rational(q^n)` for integer n is not listed, but it is needed for the mass-energy derivation (`Pow(Rational(0), 2) → Rational(0)`).

**Impact:** The coding agent will add this rule when tracing the derivation. It is a natural mechanical extension, not an architectural choice.

**Suggested clarification:** Add to §9.6:
```
Pow(Rational(q), Rational(n)) → Rational(q^n) for integer n ≥ 0
```

### Minor 2 — `Draft` metadata for non-hypothesis derivations

`Draft` requires a `DraftMetadata` with `Hypothesis core.Object`. For the mass-energy derivation (no hypothesis), what is passed?

**Impact:** The coding agent will either pass an invalid `core.Object` for the Hypothesis field (and rely on `Seal` checking it), or add an optional variant.

**Suggested clarification:** Add to §16.2.1:
```
For mechanical derivations without a candidate, DraftMetadata.Hypothesis
MAY be an invalid zero core.Object. Seal MUST fail with ProvenanceError
if Hypothesis is invalid or lacks HYPOTHESIS provenance.
```

### Minor 3 — Minkowski metric expression

§19.1 lists `MinkowskiMetric` as a populated primitive but does not specify its canonical expression. The manifest must contain a canonical_expr for it.

**Impact:** The coding agent will choose (e.g., `Symbol("η")` or a structural representation). This is not used in any derivation.

**Suggested clarification:** Add to §19.1:
```
MinkowskiMetric canonical symbol: η
```

### Minor 4 — Convention attachment scope

§19.3 says the manifest records `metric.signature = -+++`. It is unclear whether the convention attaches to all relativity objects, only the Minkowski metric, or only the manifest.

**Impact:** The coding agent will choose. The convention conflict test (Test F) does not depend on this.

**Suggested clarification:** Add to §19.3:
```
The metric.signature convention attaches to NewMinkowskiMetric() and to
any object derived from it via Substitute or Add.
```

### Minor 5 — Symbol `v` shared between mechanics and relativity

Both `mechanics.NewVelocity()` and `relativity.NewVelocity()` produce `Symbol("v")`. Within each package's derivations, this is consistent. Across packages, symbols are not mixed in MVP.

**Impact:** No impact. The coding agent will keep the packages isolated, as the canonical derivations do not mix them.

---

## Part IV: What I Verified

**Internal consistency:**
- §5.4 and §5.2.1 agree on the mint path.
- §13.2 and §15.1–15.12 agree on provenance propagation.
- §16.3 and §16.4–16.20 agree on session semantics.
- §17.0 and §17.1–17.9 agree on the manifest API.
- §18.0 and §18.1–18.3 agree on mechanics constructors.
- §19.0 and §19.1–19.8 agree on relativity constructors.
- §38 counts 39 files in §3. ✓
- §39 coverage matrix references the `REQ-*` IDs used throughout. ✓

**No contradictions remain.** The previous v2 review identified 7 contradictions; all are resolved in v2.1.

**No missing architectural decisions.** The 18 decisions listed in §41 cover: module identity, constructor authority, ledger placement, replay substrate, Identify ownership, structured assumptions, branch representation, Pow result kind, Lorentz factor body, Sqrt sign handling, zero constructors, differentiation test, solver scope, manifest mapping, file budget, provenance law, terminology, and MRC fallibility.

**MRC rule semantics are concrete.** MRC-001 through MRC-008 each have:
- an exact semantic check
- an enforcement location
- a failure error type
- at least one mapped test

**Coverage matrix is exhaustive.** §39 maps every `REQ-*` group to implementation location and tests. The prompt requires the plan to expand `REQ-002-01..26` into individual rows.

**The five previously-missing pieces are present:**
1. Paper translation protocol → §35 + `docs/paper-translation.md`
2. Adversarial review protocol → §27 (Challenge/Review structures)
3. Equality and canonicalization → §8, §9, §10
4. Symbolic expression type → §8, §10.2
5. Multi-agent workflow → §27.4 (external boundary), §24.4 (no orchestration)

---

## Part V: The Handoff Readiness

The prompt's final constraint is:

> "A coding agent MUST NOT need to choose among alternative architectures to implement this MVP."

**This is now satisfied.** The five minor ambiguities do not constitute alternative architectures; they are local choices the coding agent makes during implementation and documents in code. None reopens a settled decision.

The prompt's planning rules prohibit:
- multiple architectures ✓
- reopening settled decisions ✓
- substituting another object abstraction ✓
- moving ledger/session authority ✓
- public generic factory ✓
- public `ops.Identify` ✓
- ambient session state ✓
- `.phys` parser ✓
- `physvet` in MVP ✓
- deferred domain packages ✓
- general CAS ✓
- general root solver ✓
- theorem prover ✓
- numerical execution ✓
- mathematical source strings ✓
- `map[string]any` ✓
- floats for exact values ✓
- silent omission of requirements ✓
- "the coding agent can decide" ✓

The prompt is disciplined.

---

## Part VI: The Verdict

**Ship it. Hand it to the coding agent.**

The architecture is complete. The specification is precise. The scope is minimal (39 files). The deferrals are explicit. The coverage matrix is mandatory and exhaustive. The handoff condition §41 enumerates every resolved decision. The prompt is disciplined and prohibits the failure modes identified in the previous reviews.

The five minor ambiguities are not blockers. They are the kind of small decisions every specification leaves to the implementer, and they do not reopen architecture. A competent coding agent will resolve them locally and document the choices in code comments.

If you want to eliminate even these minor ambiguities before shipping, spend 15 minutes adding the five suggested clarifications to §9.6, §16.2.1, §19.1, §19.3, and one sentence about symbol isolation. That is the maximum additional value you can extract from another review cycle.

Otherwise, hand off now.

---

## The One-Paragraph Summary

> The v2.1 specs and prompt are ready for handoff. All four previously-blocking issues are resolved: `Differentiate` is bounded to a concrete rule set, domain objects expose `CoreObject()`, the file count is consistent (39 files in §3, cap ≤ 40 in §38), and the manifest-to-constructor mapping is specified via the `constructor` field plus a static test map. The spec adds significant value beyond v2: the `internal/kernel` boundary solves the Go visibility problem, the `MintObject` single-path constructor enforces MRC-001, `Session.Identify` eliminates ambient session state, the `Ledger` lives in `session` to break the `core → ops` cycle, canonical inputs and outputs are retained for replay, the `StepEnvelope` fixes the hash chain, bounded sign entailment avoids a theorem prover, `BranchSet` is a bounded expression node, `RestFrame` is an assumption, the nine session actions have a state machine, `DraftMetadata` bundles candidate metadata, corpus status `NONE` is separate from provenance, `APPROXIMATED` is reserved for v0.5+, and the coverage matrix in §39 enables mechanical verification. Five minor ambiguities remain: `Pow(Rational, n)` simplification, `Draft` metadata for non-hypothesis derivations, the Minkowski metric expression, the convention attachment scope, and the shared `v` symbol. None are architectural blockers; the coding agent will resolve them locally and document the choices. The prompt prohibits all the failure modes identified in previous reviews. The handoff condition §41 enumerates the 18 resolved decisions. The MVP is now specifiable as a single vertical slice from corpus to sealed ResearchCandidate, and the loop is deterministic. Ship it.
# FINAL REVIEW — SPECS v2 AND PROMPT v2

**Status:** Ready to ship after four fixes. Not five, not ten. Four.

I have reviewed v2 against the earlier v1, my previous review, and the prompt. v2 has addressed every one of the ten blocking issues from the previous review. The architecture is sound. The specification is close. Four remaining issues would force the coding agent to make architectural decisions.

---

## Part I: What v2 Fixed (Confirmed)

| Previous Gap | Where Fixed |
|---|---|
| Common physics object abstraction | §4 — two-layer model with `CoreObject()` accessor |
| MRC-001/003/007/008 semantics | §13 — concrete checks for each rule |
| Operation signatures | §14 — exact signatures for all 13 operations |
| Candidate containment mechanism | §12.2 — provenance contamination law |
| Expression node ambiguity | §7.1 — no `Derivative`/`Limit`/`Function` nodes; `Call` for functions |
| Manifest-constructor cross-check | §19.4 — `canonical_expr` field |
| `RestFrame` type | §22 — it is a constraint assumption |
| `Step` vs `Identify` ledger relationship | §15.2 — `Identify` is a specialized committed step |
| External JSON validation API | §27.3 — `Validate()` |
| Paper-translation docs status | §29 — auxiliary guidance, not canonical |
| `cmd/physvet/` in package layout | Removed |
| Session action count | §15.1 — nine actions, correctly listed |
| `physvet` deferred | §36 — clear |

That is a comprehensive response to the previous review. v2 is the strongest specification artifact in this project.

---

## Part II: The Four Blocking Issues

### Issue 1 — `Differentiate` on domain quantities is undefined

**Where the conflict lives:**

- §7.1 forbids `Derivative` as an expression node.
- §14.8 says `Differentiate` uses "the existing expression nodes; there is no `Derivative` node."
- §23.2 requires `Differentiate(Velocity(v), Time(t))` to produce "the expected symbolic derivative structure."

**The problem:** If `Velocity(v)` is a domain object whose `Expr()` is `Symbol("v")`, then `Differentiate(Symbol("v"), Time)` cannot be evaluated. Symbolic differentiation of `v` with respect to `t` requires either a `Derivative` node or knowledge that `v` is `dx/dt`. Neither exists.

**The coding agent will invent a resolution.** It will either (a) add a `Derivative` node, violating §7.1; (b) reject the operation, violating §23.2; or (c) produce an arbitrary symbol, violating the spirit of symbolic differentiation.

**Fix — insert after §14.8:**

> When `Differentiate` is applied to a domain quantity whose expression is a bare `Symbol`, the result is a new `Symbol` representing the derivative. The new symbol's provenance records the differentiation. The new symbol's `Kind` is determined by the domain package via a registered derivative map. If no derivative map entry exists for the input symbol, the operation returns `UnsupportedOperationError`.
>
> For MVP, the only registered derivative is:
> ```text
> Symbol("v") → Symbol("a")   with Kind = Acceleration
> ```
> All other domain-symbol differentiation is unsupported in MVP.

This is the minimal resolution that satisfies §23.2 without adding a `Derivative` node.

### Issue 2 — Domain-to-object conversion is one-way

**Where the gap lives:**

- §4.1 says domain types wrap `core.Object` and expose `CoreObject() core.Object`.
- §14 says all operations return `(core.Object, error)`.

**The problem:** There is no way to go from a `core.Object` back to a domain type. If a test does `Differentiate(Velocity(v).CoreObject(), Time(t).CoreObject())` and gets a `core.Object` with `Kind == Acceleration`, the test cannot convert it to `mechanics.Acceleration`.

**The coding agent will invent a conversion API.** The prompt forbids this.

**Fix — insert after §4.1:**

> For MVP, generic operations return `core.Object`. Domain-typed values are only required at construction time. Callers work with `core.Object` directly after construction. Domain packages MAY expose `FromObject(core.Object) (DomainType, error)` for convenience, but this is not required for the MVP and the coding agent should not add it unless a specific test requires it.

This is the minimum fix. It tells the coding agent "don't invent this."

### Issue 3 — File count contradiction

**Where the conflict lives:**

- §3 lists 47 files.
- §37.1 says "Target no more than 40 repository files total."

**The problem:** The spec contradicts itself. The coding agent will either trim the file list (changing the spec) or ignore the cap.

**Fix — change §37.1 to:**

> Target no more than 50 repository files total, excluding `go.mod`, `go.sum`, generated artifacts, and VCS metadata. Any plan or implementation exceeding the cap MUST justify every excess file against a normative MVP requirement.
>
> Convenience-only files, deferred-domain scaffolding, compatibility shims, and placeholder abstractions are prohibited.

50 is the realistic cap given the current file list. Alternatively, consolidate the 8 `core/*_test.go` files into 3 test files. I prefer raising the cap — the file-per-concern layout is cleaner.

### Issue 4 — Manifest item to Go constructor mapping

**Where the gap lives:**

- §19.4 says manifest tests must "call the specified domain constructor/relation."
- The manifest item has an `id` field like `"NewtonSecondLaw"`.
- No mapping from `id` to a Go function is specified.

**The problem:** The manifest test cannot know which Go function to call for `NewtonSecondLaw` without a mapping.

**Fix — insert into §19.2:**

> Each manifest item MUST include a `constructor` field naming the Go identifier of the function that produces the item. Example:
> ```json
> {
>   "id": "NewtonSecondLaw",
>   "constructor": "mechanics.NewtonSecondLaw",
>   "canonical_expr": {...}
> }
> ```
>
> The manifest test decodes `constructor` and uses a static `map[string]func() core.Object` defined in the test file to invoke the function. The mapping is part of the test, not the production code.

This tells the coding agent exactly how the manifest test resolves `id` to a Go function without requiring runtime reflection.

---

## Part III: The Two Minor Issues (Fix If Time Permits)

### Minor 1 — `HYPOTHESIS` object vs `ResearchCandidate` terminology

§27.1 says "The hypothesis package may create provisional candidate concepts marked `HYPOTHESIS`; only the session may seal the final research artifact." This is correct but the terminology "candidate concept" vs "research candidate" is easy to confuse.

**Fix — insert after §16.1:**

> Two distinct things are called "candidate" in the spec:
> 1. A **candidate concept** is a `core.Object` with `Provenance.Status == HYPOTHESIS`. It is created by the hypothesis package.
> 2. A **`ResearchCandidate`** is a sealed artifact that carries a complete derivation, predictions, and falsification conditions. It is created only by a `Session`.

This is a clarity fix, not a functional one.

### Minor 2 — Coverage matrix requirement duplication

The coverage matrix requirement appears in §34, §38.1, and the prompt twice. The duplication is harmless but noisy. No fix required.

---

## Part IV: What to Do

1. Apply Fix 1 (Differentiate on domain quantities) — **blocking**
2. Apply Fix 2 (Domain-to-object conversion) — **blocking**
3. Apply Fix 3 (File count) — **blocking**
4. Apply Fix 4 (Manifest constructor mapping) — **blocking**
5. Optionally apply Minor 1 (terminology)
6. Hand to the coding agent

The prompt itself needs no changes. The prompt correctly references the specs. The prompt's "coverage matrix" duplication is harmless.

**Time estimate to apply all four fixes:** 20 minutes of editing. Do not do another full review pass after applying these — the specs will be ready.

---

## Part V: The Go/No-Go

**Go.** Apply the four fixes. Hand to the coding agent.

The architecture is sound. The specification is specific. The scope is minimal. The deferrals are explicit. The MRC semantics are concrete. The ledger is hash-chained and replayable. The candidate containment is enforceable. The human boundary is clear.

The four fixes are surgical. None require architectural change. None reopen settled decisions. All four are the kind of detail that is easy to miss when writing a spec and easy to fix when pointed out.

After these four fixes, v2 will be a specification that a coding agent can execute without architectural decision-making. That has been the goal throughout. It is achievable now.

---

## The One-Paragraph Summary

> The v2 spec has addressed every blocking issue from the previous review: the common object abstraction, the MRC rule semantics, the operation signatures, the candidate containment mechanism, the expression node ambiguity, the manifest-constructor cross-check, the RestFrame type, the Step vs Identify ledger relationship, the external JSON validation API, the paper-translation docs status, the `cmd/physvet/` removal, the session action count, and the `physvet` deferral. Four blocking issues remain: (1) `Differentiate` on domain quantities is undefined because §7.1 forbids `Derivative` nodes while §23.2 requires differentiating `Velocity`; (2) there is no way to convert a `core.Object` back to a domain type; (3) the file count in §3 (47 files) contradicts the cap in §37.1 (40 files); (4) the manifest item to Go constructor mapping is unspecified, so the manifest test cannot resolve `id` to a function. Two minor issues are the terminology overlap between "candidate concept" and "ResearchCandidate" and the duplication of the coverage matrix requirement across three sections. The four blocking fixes are surgical: a paragraph for each. Time to apply is 20 minutes. After the fixes, the spec is ready for the coding agent. Do not do another full review pass. The architecture is sound. Ship it.
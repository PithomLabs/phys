# Final Review: v2 Prompt + v2 Specs — Dispatch Verdict

**Bottom line first: dispatch-ready after eleven surgical patches, roughly 90 minutes of editing.** Every blocking finding from the prior round is resolved — this is the first artifact pair in the corpus where the resolution is complete rather than partial:

- **B1 (two-layer model): resolved** — §4.1's wrapper pattern + `CoreObject()` + §5's `Kind` enum with the Expression exception in MRC-003 is exactly the right design, and the Energy/Torque rationale is now in-spec.
- **B2 (containment law): resolved** — §12.2 states the propagation law with the laundering list; §16.2/§27.1 split candidate-minting (AI) from sealing (session).
- **B3 (reduction test): resolved honestly** — §23.4's `LorentzFactor(v→0)=1` with the Newtonian-KE deferral stated, and §20 exempts mechanics.
- **H1/H3/M3/M4: all resolved** — rational-as-string JSON with round-trip test P and negative 14; manifest validator as pure function over bytes; external test packages mandated; StepEnvelope hashing specified.
- **P1/P2: resolved** — Open Spec Items protocol is in the prompt; coverage matrix is mandatory in three places.

The remaining defects are no longer *missing mechanisms* — they are **spec-internal contradictions and derivation-breaking gaps** that the coverage matrix itself would eventually catch, but only after the coding agent has burned a cycle. Patch them now. Patch text is paste-ready.

---

## Blocking patches

### 1. The ledger cannot replay — canonical outputs aren't retained (§15.5 vs §15.7)

§15.5 stores only hashes; §15.7 step 3 says "reconstruct inputs from canonical hashes." SHA-256 is not invertible. Worse: Test Q ("modified intermediate ledger output") has nothing to modify if only hashes exist.

> **Patch (§15.5):** "Each committed step MUST additionally retain `OutputCanonical` (the canonical JSON of the step's output object) alongside `OutputHash`. Hashes establish integrity; retained canonical outputs are the replay substrate. Replay reconstructs inputs from retained canonical outputs and verifies every recorded hash against its retained form. The tamper tests modify a retained canonical output and verify detection by hash mismatch and/or replay divergence."

### 2. The canonical derivation cannot be built — no zero-quantity constructors exist (§4, §23.3)

Test I requires `Substitute(relation, p, 0)` and the branch constraint `E ≥ 0`. `0` must be a `core.Object` with dimension M L T⁻¹ / M L² T⁻² — but §4.3 forbids generic factories and **no constructor produces it**. Test I is unimplementable as specified.

> **Patch (new §4.5):** "Each populated domain package MUST expose fixed constant constructors required by its canonical derivations, at minimum `relativity.ZeroThreeMomentum()` (M L T⁻¹, Kind ThreeMomentum, expr `Rational(0)`) and `relativity.ZeroEnergy()` (M L² T⁻², Kind Energy, expr `Rational(0)`). Constants are trusted objects with provenance `DEFINED`. The canonical derivation substitutes `ZeroThreeMomentum()` and constructs the branch constraint as `Compare(E, ZeroEnergy(), gte)`."

### 3. `Identify` has two contradictory identities (§14.13 vs §15.1)

§14.13 defines a package-level `Identify(a, b, justification)`; §15.1 lists `Identify` as a session action; §14.13 says it "records the identification in the active session" — with no session parameter. Ambient session state would violate the determinism rules.

> **Patch (§14 preamble + §14.13):** "Package operations are pure functions: they validate MRC, propagate assumptions/conventions/provenance, and return new immutable objects. They never read or write a ledger. Recording occurs only through session methods. Accordingly `Identify` is a session method — `Session.Identify(a, b, justification)` — and no package-level `ops.Identify` exists. `Session.Step(label, object)` records a result produced by a §14 function; the session refuses to record the same output hash twice within one derivation."

### 4. `LorentzFactor` has no body, and the sqrt rules can't complete either canonical test (§7.7–7.8, Tests I and J)

`Limit` on a `Call` node requires the function's definition expression; none is specified. And the final step of Test I — `Sqrt((mc²)²) → mc²` — is exactly the transform §7.7 forbids "without a supplied assumption," while the supplied assumptions (`m ≥ 0`, `c > 0`) have no specified entailment mechanism. Test J also needs `Sqrt(1) = 1`.

> **Patch (§7.8):** "Each MVP `function_id` has a fixed internal definition body with declared preconditions. Registry: `lorentz_factor(v) := 1 / Sqrt(1 − Pow(v/c, 2))`, precondition `v < c`. `Substitute` and `Limit` operate through the definition body; `Limit` on a `Call` substitutes the limit point into the body and simplifies."
>
> **Patch (§7.7):** "`Simplify` MUST evaluate `Sqrt` of a non-negative perfect-square rational to its exact root. `Simplify` MAY transform `Sqrt(Pow(x,2)) → x` only when the current assumption set entails `x ≥ 0` — in MVP, every multiplicative factor of the base is a positive rational constant or a symbol carrying an explicit non-negativity precondition (e.g., `m ≥ 0`). Absent entailment, the expression is left unchanged and the derivation remains incomplete. No other sign-dependent transform is permitted."

### 5. `Pow` has no result kind (§14.5)

The derivation builds `(pc)²` and `E²` via `Pow`. MRC-003 assigns Expression kind only to `Multiply`/`Divide`. `Add((pc)², (mc²)²)` must be Expression+Expression to pass — so `Pow` must produce Expression too.

> **Patch (§14.5):** "`Pow` returns Kind `Expression`, dimensions transformed by the exact exponent — matching `Multiply`/`Divide`."

---

## High/medium patches

### 6. The differentiation micro-test is trivial as written (§23.2)

`Differentiate(Velocity(v), Time(t))` where `v` is an independent symbol yields `0` — it exercises nothing, and no arbitrary time-dependence is expressible (no function nodes). Replace with a test that is nontrivial *and* covers the otherwise-untested MRC-003 clause (Expression vs named comparison):

> **Patch (§23.2):** "Construct `K := KineticEnergy(m, v)` (canonical `1/2 · m · v²`). Execute `Differentiate(K, v)`. Verify: (1) the canonical result is `m · v`; (2) `Compare(result, MomentumRelationLHS, eq)` is accepted under MRC-003 (derived `Expression` vs named `Momentum`, dimensions M L T⁻¹) and yields a true relation. This exercises product rule, power rule, canonical equality, and the Expression-vs-named comparison clause."

### 7. Draft/Commit state machine is ambiguous (§15.1)

"Draft" appears as both an action and implicitly a stage. Prescribe the flow:

> **Patch (§15.1):** "`Step` and `Identify` append fully hashed entries to the session's draft buffer (uncommitted). `Commit` finalizes the draft buffer into immutable, hash-chained ledger steps in recording order. `Conclude(object)` records the final object, which MUST equal the last committed step's output. `Seal` freezes the ledger, fixes the derivation hash, and enables `ResearchCandidate` minting."

### 8. The spec contradicts its own file cap (§3 vs §37.1)

The normative layout sums to **47 files** (core 19, ops 11, mechanics 5, relativity 7, hypothesis 2, docs 3) against a 40-file cap. Four reviewers will each find this.

> **Patch (§37.1):** "The normative §3 layout is the MVP file baseline (~47 files including adjacent tests). The cap in this section applies to files beyond §3; any file outside §3 requires justification against a normative requirement."

### 9. `APPROXIMATED` is unreachable, but the coverage matrix demands a test for every status (§12.1 vs §34)

No MVP operation produces it (series algebra deferred). The matrix requirement and the status enum are now in tension.

> **Patch (§12.1):** "`APPROXIMATED` is reserved in MVP: no MVP operation produces it. It is exempt from the §34 coverage matrix for this release and becomes mandatory at v0.5."

### 10. Canonical-JSON rules don't explicitly cover the artifacts they must cover (§9, §15.6, §32)

`StepEnvelope` and `ResearchCandidate` canonical JSON are mandated but their schema discipline is only defined for expressions.

> **Patch (§9):** "§9's canonical-JSON rules — explicit structs, fixed field order, rationals as strings, no floats, no `map[string]any` — apply to every canonical artifact encoding: expressions, assumption sets, convention sets, provenance records, `StepEnvelope`, and `ResearchCandidate`."

---

## Trivial (one minute each)

- **Prompt:** the coverage-matrix paragraph appears twice verbatim at the end of Planning rules — delete one.
- **§16.2:** add "The candidate constructor requires explicit declaration of the concept's `Kind` and `Dimension`; both become immutable authoritative metadata; provenance is fixed to `HYPOTHESIS`."
- **§3:** `core/corpus.go` has no direct test file — either accept package-level coverage via the manifest tests (state it in the matrix) or add `corpus_test.go` under the amended baseline.

---

## Dispatch checklist

Hand the coding agent:

1. **Both documents with the eleven patches applied** (patches 1–5 are blocking; 6–10 prevent reviewer findings; 11 is hygiene).
2. **The standing instruction** that the Open Spec Items protocol is active: any residual under-determination gets flagged there, never improvised.
3. **The two acceptance yardsticks**, named explicitly: §35's tests A–S and the §34 coverage matrix. Tell the agent the matrix is the *first* deliverable of the plan, not the last — writing it before the sequence forces every spec item to find its test early, which is how you avoid a second review round.
4. **The scope tripwire:** §37.1 as amended, plus the prompt's own rule — a plan that proposes a file outside §3 must justify it inline.
5. **One sentence of intent for the agent:** "The MVP's purpose is to prove the *architecture* (containment, provenance, identify-gate, replay) on two small corpora — not to do physics. Tests I and J validate the substrate; the specs say so; do not embellish them."

With those patches in, the pair has no known internal contradictions, no unimplementable tests, and no decision points left for the agent to improvise. That is the condition you were waiting for — ship it.
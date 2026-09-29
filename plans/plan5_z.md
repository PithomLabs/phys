# Adversarial Review: plan4_review.md — The Four-Plan Verdict, Audited

**Method and scope limit, stated first.** The four plans themselves are not in front of me — only the review that adjudicates them. So I cannot independently audit its scorecard (whether DeepSeek's packaging is truly "Strongest," whether Gemini overreaches as claimed). What I *can* do, and will do, is: audit the review's **verdict and proposed synthesis** against your six criteria; test its internal consistency against the corpus's own accumulated findings (the carried-forward mechanism ledger from the v0.1/v0.3 reviews); and discharge your explicit request — the **collective gaps** no plan (and the review itself) caught. Where I concur with the review, my concurrence is conditional on its characterizations being accurate; I flag which conclusions survive regardless.

**Headline verdict.** The review's recommendation — *Z's derivation/provenance spine + DeepSeek's package architecture + physvet as a layer + a first-class epistemic system* — is the right synthesis direction, and its five additions are genuinely the right five *epistemic* additions. But the synthesis has a structural blind spot: **every one of its five additions is epistemic; none is mechanical.** The v0.3 review established four blocking mechanical fixes (carriers-not-strings, canonical identity + exact rationals, conventions, dimensions/immutability). The plan4 synthesis mentions *none of them*. Adopted as written, v0.4 would gain the epistemic spine and silently regress the mechanical one — the exact failure mode (silent amputation of load-bearing mechanisms, absent from the removal list) that killed v0.3's freeze. The review that caught everyone else's gaps has one of its own.

---

## 1. The six criteria, audited

| Criterion | Review's handling | Audit |
|---|---|---|
| 1. Primitives + MRC at compiler level | Endorses Z's five levels + physvet gate | Right mechanism, one framing error to fix (§2 below) |
| 2. Relativity as Popperian eval target | Knowledge-status metadata | Half right — an eval target needs **executable anatomy**, not just metadata (§4) |
| 3. Packages as AI corpus | Machine-readable manifest contract | Right, but manifests must be **validated, not just written** (§5) |
| 4. Collective gaps | Five additions | Good list, incomplete — see §6 for the missing seven |
| 5. No adjudication | Rejects truth scores, rejects Qwen's overreach | Correct; one residual to sharpen (§3) |
| 6. Hypothesis formation, human delegation | Hypothesis struct + ResearchCandidate handoff | Correct boundary; the **candidate-object tension** is unaddressed (§6) |

---

## 2. Criterion 1 sharpening: bridges are typed, and the review's own example proves it

The review claims MRC "is fundamentally about permitted bridges, not just types," citing `Energy + Temperature` becoming meaningful via multiplication by $k_B$. Decompose that example: $k_B$ has dimension $\text{Energy}\cdot\text{Temperature}^{-1}$; $k_B \cdot T$ is an ordinary typed multiplication yielding Energy. **No bridge mechanism is needed — the type system already did the bridge.** The objects that genuinely exceed types are exactly two, and the corpus has names for both: **identity claims** (the `Identify` gate — the Einstein moment) and **structure-required operations** (a `Poset` has no `Duration` — the problem-of-time case).

Why this matters: if you believe bridges are extra-typal, you build a bridge registry — a lookup table, which is the jurisdiction matrix reborn. If you recognize that dimensional bridges are typed constructors and identity bridges are the `Identify` gate, you need **no new mechanism at all** — the five-level enforcement table from the accepted Einstein design already covers criterion 1 completely. The review's framing risks re-growing the machinery the project spent three rounds killing.

## 3. Criterion 5 sharpening: status is data, not computation

The review's status taxonomy (`Definition, Postulate, Law, Identity, Established Relation, ... Empirical Result`) conflates two axes that must be separate:

- **Provenance status** — computable by the library from the ledger: `DEFINED / POSTULATED / DERIVED / IDENTIFIED / APPROXIMATED(order n) / HYPOTHESIS`. The library mints these.
- **Corpus status** — curated by humans, versioned, carried as data: `ESTABLISHED / CONTESTED / SUPERSERDED / **FALSIFIED**`.

"Law" is neither. "Empirical Result" is evidence, not status. And the library must **never compute** corpus status — it reads it, propagates it through provenance closures, and renders it. The review's ban on truth scores is correct but incomplete: the sharper rule is that corpus status is *input* the library consumes, never *output* it produces.

## 4. Criteria 2+3: what an eval target actually is

The review says relativity should be "a target for derivation/evaluation" and gives it a rich manifest. But **a manifest is not an eval target — an executable anatomy is.** For criterion 2 to bind, each framework package needs four things beyond API + metadata:

1. **A derivation suite** — golden derivations with expected canonical outputs and ledger shapes. The Einstein 1905 ledger is the canonical example: the eval is "from this package's postulates, reproduce this ledger, hash-identically." (Note: this *finally* discharges the differentiator-test finding that has recurred in the v0.1 and v0.3 reviews — the thing that would certify "SymPy plus units" as conformant. Derivation suites are the differentiator. Third time raised; this time it has a natural home.)
2. **An executable reduction suite** — Popperian refinement means frameworks stand in *reduction relations*: GR → SR in the flat limit, SR → Newtonian at $v \ll c$, QM → classical as $\hbar \to 0$ (with declared caveats). `physics.Reduce(relation, limitCondition)` should run the symbolic limit and return the predecessor's result, checkably. "Reduces to: classical limit" as a manifest *string* is documentation; as an *executable test* it is physics.
3. **Distinguishing predictions as structured data** — where the framework departs from its predecessor (SR vs Newtonian at $v \sim c$). This is what makes it a *falsifiability* target: the departures are the handles.
4. **An anomaly registry** — see §7.

## 5. Criterion 3 sharpening: manifests rot unless validated

The package contract (Concept, Assumptions, Regime, Derivable-from, Reduces-to, Status...) is right. But YAML-adjacent prose next to code drifts. The contract must be: (a) represented as data the library can **query at runtime** (`relativity.Manifest()`); (b) **cross-checked by tests** wherever mechanically possible — if the manifest says "Dimensions: M L² T⁻²" for `EnergyMomentumRelation`, a test asserts the actual object's `Dim()` equals it; if it says "Reduces to: X under C," the reduction suite runs it. Corpus edits are spec-grade events with corpus-wide regression — the standing governance rule, applied to a new artifact class.

## 6. The collective gaps (your explicit request #4)

The review found five. Here are the ones it — and the four plans — missed.

**G1 — The mechanical-spine regression (highest severity).** The five additions are all epistemic. The accepted blocking fixes from the v0.3 review are absent from the synthesis: carriers-not-strings constructors, canonical form + structural equality + exact rational serialization, conventions as value metadata with conflict errors, dimensions in core, immutability/closed object model. A v0.4 frozen without these repeats v0.3's failure: the epistemology without the receipts. **Merge both lists; they are not in tension — the ledger's premise closures need canonical hashing, the eval suites need exact serialization, the reductions need convention discipline.**

**G2 — Assumption algebra.** The review requires results to "carry" assumptions but specifies no **propagation semantics** — and without propagation, the assumption system is a comment field. Concrete requirement: values carry assumption sets; every operation defines its merge rule (`ExpandSeries` *adds* a truncation-order assumption; `Solve` may add nondegeneracy; `Reduce` requires the limit condition and refuses without it); contradictory merge = typed error, not silent intersection. This is not a nicety — it is load-bearing for criterion 6, because composing frameworks (ToE work) means composing assumption sets, and a silent intersection of SR's and QM's assumptions is precisely how a fake "bridge" gets laundered. **Assumption-conflict surfacing is the mechanical core of the bridge question.**

**G3 — The candidate-object tension.** The review's §7 demands `hypothesis.NewObject` — an open factory. The accepted v0.3 review established the **closed object model** (unexported fields, constructor monopoly) as the anti-forger mechanism. These collide: an open factory is the impostor hole reopened. Resolution, explicitly: candidate objects are a distinct provenance class (`CANDIDATE`) — they participate in all operations (that is the point; the AI must be able to compute with new concepts), but operations on them tag results as candidate-contaminated; they can never mint canonical `DERIVED` status, enter the corpus, or forge core provenance; promotion to corpus object is a human-gated, spec-grade event. Open representation, closed authority.

**G4 — Falsification cases in the corpus.** A Popperian corpus that contains only survivors cannot teach falsification. The corpus must encode *deaths*: Mercury perihelion vs Newtonian gravity, ultraviolet catastrophe vs classical statmech, the aether. As **negative derivations** — structured cases where the predecessor framework's prediction and the observed divergence are both represented, with the anomaly as the pivot. These are eval targets of a second kind: "identify which assumption of framework F is contradicted by anomaly A." That is a testable AI capability and the most direct operationalization of criterion 2.

**G5 — physvet dual-source-of-truth.** Adopting physvet is fine (it's `go vet` for physics — implementable with `go/ast`, no new language, consistent with the v0.3 §47 removals). But a static analyzer that re-implements dimension or contract logic will drift from the library's semantics. Rule: **physvet must consume the library's contract metadata** (it can — manifests, constructor signatures, error taxonomy) and never re-derive physics rules independently. One source of truth; two enforcement points (construction-time in the library, whole-program patterns in vet).

**G6 — The handoff artifact needs format and ownership.** `ResearchCandidate` is the correct terminus (criterion 6) and correctly leaves EBP outside — but it must be *shaped* as the EBP/Solvent input without the library knowing about EBP (the d6 boundary, preserved via artifact schema, not integration). And only the session/ledger mints it — the AI cannot hand-craft one, or the chain-launderer attack returns at the handoff boundary. Exact serialization applies here more than anywhere: this artifact is what humans adjudicate.

**G7 — Bootstrapping order.** Nobody specified who writes the first trusted derivations. Order: core + mechanics (simplest suite) → SR (the Einstein suite) → GR, QM, QFT, statmech. The first framework package is both the eval target and the proof that the library can host one — build it before the epistemic layer is finalized, because the epistemic layer's fields should be *derived from* what the first suite actually needed, not specified in the abstract (the earned-language principle, applied to metadata).

## 7. The Popper inversion the review missed

The review treats Popper as a *labeling* requirement (don't call relativity true). But Popper's philosophy, taken seriously as a design constraint, inverts the corpus's center of gravity: **hypothesis generation should be anomaly-driven.** If the corpus packages each carry a structured registry of open problems and known failures (GR: singularities, the cosmological-constant problem; QM: measurement; statmech: the arrow), then the AI's criterion-6 workflow becomes mechanically guided: *read the anomaly registry → locate the assumptions it strains → formulate a candidate that modifies exactly those assumptions → the ledger shows which anomalies the candidate addresses and which it leaves*. Falsifiability stops being a philosophical posture and becomes a **data structure** the hypothesis layer must cite. That is a stronger answer to criterion 6 than the Hypothesis struct alone — and it is absent from all four plans and the review.

## 8. v0.4 non-negotiables — merged and prioritized

| # | Item | Source |
|---|---|---|
| 1 | Mechanical spine restored (G1: carriers, canonical/exact, conventions, dims, immutability) | v0.3 review, carried |
| 2 | Assumption algebra with propagation + conflict errors (G2) | New |
| 3 | Provenance/corpus status split as two axes; library never computes corpus status (§3) | Sharpened from review |
| 4 | Candidate-object class with closed authority (G3) | New — resolves review's §7 vs closed model |
| 5 | Eval-target anatomy: derivation suite, executable reductions, distinguishing predictions, falsification cases (§4, G4) | New — subsumes differentiator finding |
| 6 | Anomaly registries per framework; hypothesis layer must cite them (§7) | New |
| 7 | Corpus manifests validated by tests; spec-grade edit governance (§5) | Sharpened from review |
| 8 | physvet consumes library metadata; no second truth source (G5) | Sharpened from review |
| 9 | Session-minted ResearchCandidate, EBP-ignorant, exactly serialized (G6) | Sharpened from review |
| 10 | Review's five (assumptions, status, hypothesis vocab, package contract, human handoff) | Adopted as epistemic layer |

---

**Closing verdict.** Concurrence with the review's choice of Z as spine and its synthesis formula — conditionally. Its five additions complete the *epistemic* layer this project has needed since design0, and its rejection of Qwen's epistemic overreach and Gemini's truth-adjacent framing is correct per criterion 5. But the synthesis as written would ship the epistemology and quietly drop the receipts — canonical identity, exactness, conventions, dimensions, the closed object model — because none of them survived into its five-point list, exactly as none survived into v0.3's §47. And it stops one step short of the Popper implication: the corpus should not merely *refuse* to call relativity final — it should carry the anomalies that make the next theory findable. The library doesn't judge hypotheses; it builds the hunting ground where hypotheses are aimed, derivable, reducible, and falsifiable — and then it stops, and hands a sealed artifact to humans. That is the v0.4 worth freezing.
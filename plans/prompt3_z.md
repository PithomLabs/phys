# Adversarial Review: prompt3.md — As an Instrument for Eliciting an Implementation Plan

**Method.** The artifact under review is a *prompt*, so the review asks: what will a competent coding agent actually build from this text, where will it be forced to improvise, and can the plan it produces be accepted by four independent reviewers against what the prompt demands? The corpus's accumulated findings are the audit baseline — a prompt that silently drops a previously-accepted mechanism gets the same verdict a spec would.

**Headline verdict.** This is the most disciplined artifact in the entire corpus — the decisions-vs-questions separation (§2), the discipline rule (§21), and the deliverable schema (§20) are exactly what a plan-elicitation prompt needs. But it has one outright bug, one silently dropped top-priority finding, and a structural flaw in its acceptance design: **the vertical slice exercises only seven of the ten required components**, so a conforming plan can ship candidate/anomaly/hypothesis machinery as untested stubs and still pass §17 and §20M. Patch nine items and this prompt will elicit a reviewable plan; as written, it will elicit a plan whose tests certify the parts that were already going to work.

---

## What genuinely survives

- **§2 decisions-not-questions** — prevents the relitigation loop that consumed four design rounds.
- **§2E (fallible MRC, scoped exceptions)** and **§2D (one semantic source of truth, physvet deferrable)** — both carried findings, correctly resolved.
- **§5's propagation equation and conflict-surfacing** — the assumption algebra (G2), stated compactly.
- **§6's two-axis status split with "may carry but must not infer"** — the exact sharpened rule.
- **§9 "open representation does not imply open authority"** — the G3 tension resolved in one sentence.
- **§13's mechanical spine list** — G1 substantially present (canonical form, exact rationals, conventions, conflict detection, immutability).
- **§14 ledger authority ("AI must not be able to hand-author a trusted derivation artifact")** — F1/G6.
- **§19 including "integration with EBP 2.1 itself"** — the d6 boundary held.
- **§20 A–N schema + four-reviewer framing** — reviewable deliverable shape.

---

## Findings

### F1 (Blocking — outright bug): dangling reference

§21: "Whenever a proposed component does not directly support the MVP vertical slice **or one of the six project criteria**, defer it." The six criteria appear nowhere in the prompt. The coding agent cannot enforce a discipline rule whose reference is undefined — and the four reviewers cannot audit its enforcement. Inline them (primitives + compiler-level MRC; frameworks as Popperian eval targets; packages as machine-readable corpus; no adjudication; hypothesis formation with human delegation) or rewrite §21 to reference §1's mission bullets.

### F2 (Blocking): carriers-not-strings is gone — the single most important carried finding

§13 says "immutable/controlled physics objects" and "constructor-controlled authoritative metadata." That is a statement about *metadata*, not about the v0.3 review's #1 fix: **strings are names, never carriers; relational objects require typed carriers** (`Metric(on Spacetime)`, `StressEnergy(on Spacetime)`, `State(of System)`), and every constructor is audited against the rule. Without it, the agent builds `Mass("m")`-style name-keyed constructors throughout, `physics.Time("t")` reconstructs, and the entire MRC story sits on a foundation the Go compiler cannot check. Add as an explicit decision line in §2 (a new §2G) — this is precisely the kind of thing that must be a *decision, not a question to reopen*.

### F3 (Blocking): the acceptance target does not cover the requirements

The vertical slice (1905 + one classical relation) exercises: corpus read, construction, assumptions, per-step MRC, Identify, ledger, canonicalization, review artifact. It does **not** exercise: **candidate objects (§9), anomaly inspection (§11), falsifiability structure (§10), promotion blocking**. §11 even says "demonstrate that an AI can inspect an anomaly" — but nothing binds that demonstration to acceptance, and §20M is agent-authored, so the agent's own acceptance criteria will mirror its own implementation. Result: a plan can satisfy the prompt while §9–§11 are spec-only code. Fix: require three micro-demos in §20M as mandatory — (a) **candidate contamination**: propose a candidate object, derive with it, verify results are marked provisional and promotion is blocked; (b) **anomaly inspection**: one structured anomaly (UV catastrophe or Mercury perihelion) with a candidate response citing it; (c) **negative tests**: dimension error, category error, convention conflict, assumption conflict.

### F4 (Blocking): the 1905 slice is unimplementable from the §12 operation list

The derivation mechanically requires **series expansion with a tagged truncation order** (γ−1 → β²/2, with the order-2 tag entering provenance per §14's own "approximations" requirement) and **Solve/Compare** for the conclusion. §12's list — construct, arithmetic, power, substitute, differentiate, integrate, simplify, limit/reduction, identify — contains none of these. "Limit/reduction" is not a substitute: a limit operation cannot emit a second-order series carrying an order tag, and the tagged truncation is the mechanism that prevents the classic silent-dropping lie. The agent will improvise an expansion, and the truncation-provenance chain the corpus spent two reviews establishing will break in the improvisation. Add `ExpandSeries`/`Truncate(order)` and `Solve` to §12 explicitly. Related trap: if `Integrate` is included, its domain parameter must be a **physics-typed object** (the ban line cuts through math-typed domains — the v0.3 F4 finding); if that's awkward for MVP, defer Integrate and say so.

### F5 (High): "Prefer a machine-readable representation" — modal verbs leak

§6: "Do not make comments the sole canonical metadata source. **Prefer** a machine-readable representation that can be validated by tests." An agent satisfies "prefer" with doc comments plus one JSON file. The accepted finding was normative: manifests are data (JSON), accessible at runtime (`relativity.Manifest()`), **cross-validated by tests against the actual objects** (manifest says `M L² T⁻²` → test asserts `Dim()` agrees). Change "Prefer" to "Must."

### F6 (High): ledger integrity covers bypass but not tampering

§14 blocks hand-authored artifacts that *bypass* the ledger. It does not block **post-hoc tampering** with a legitimate ledger — edit one recorded step's result and the trusted provenance survives. The cheap, deterministic mechanism is **replay validation** (each recorded step re-executes from prior canonical results and must reproduce the recorded output) plus **hash-chaining** of the ledger. This is also the actual mechanism behind §15's "provenance problem" reviewer finding — without it, that finding is a label with no detector. Extend the minting rule while here: **only the session mints the `ResearchCandidate`**; a hand-authored candidate artifact fails validation (otherwise the chain-launderer returns at the handoff boundary).

### F7 (Should): acceptance is self-graded — seed the must-pass list

§20M asks the agent to state acceptance criteria. Seed the non-negotiables so reviewers have a fixed yardstick: dimension error, category error, assumption conflict, convention conflict, candidate contamination + blocked promotion, ledger bypass rejection, ledger tamper rejection, byte-identical determinism across runs, exact rational round-trip through serialization, truncation order present in the final result's provenance closure. Also ask the plan for **per-package size estimates** — the four reviewers need a scope yardstick, and "conservative" is unauditable without one.

### F8 (Should): package sprawl invited

§2C lists seven packages; the slice needs two. Say explicitly: **mechanics and relativity (SR subset) are populated; EM/QM/QFT/statmech are deferred** — absent directories, not stubs. Stubs rot and reviewers will waste their effort auditing scaffolding.

### F9 (Should): §2E is a scope trap

"The architecture must permit future human-validated exceptions" — read literally, an agent will build an exception-registry workflow (approval states, evidence requirements) in the MVP. What the MVP needs is only **rule versioning** (MRC rules carry stable IDs so a future exception can be scoped against them). Reword §2E: the *permission* is retained via rule versioning; the *mechanism* is a §20N deferral.

### F10 (Should): Identify is under-specified for the thing the whole corpus is about

Require, in one line each: `Identify` takes **two explicit operands** plus a **recorded physical justification** (the AI's stated reason — Einstein's "the coefficient of v² is the inertial mass"); it mints `IDENTIFIED` provenance; it is **never reachable through `Simplify`**; identifications are surfaced in the review artifact as a first-class list (the reviewer's primary attack surface).

### F11 (Nice): export equality

§13 has deterministic structural equality as a safeguard; the ledger, review, and any route comparison need it as **public API** (`ops.Eq`, canonical `Hash()`). One sentence. Similarly: make "AI reads corpus" mechanical (the Manifest accessor), pin the Go toolchain and module path, request golden-table tests, and add to §20N the explicit deferrals — EM/QM/QFT/statmech, the exception workflow, review-by-recomputation (`deriv.Diff`), and the CST/problem-of-time framework — so no reviewer counts their absence as a defect.

---

## Agent-misreading catalog (how the prompt fails if unpatched)

| Misreading | Mechanism | Finding |
|---|---|---|
| **The String-Constructor agent** | "controlled objects" read as name-keyed constructors | F2 |
| **The CAS agent** | missing expansion op → improvises a general simplifier | F4 |
| **The Doc-Comment agent** | "prefer machine-readable" satisfied with prose | F5 |
| **The Stub-Sprawl agent** | seven listed packages → seven scaffolds | F8 |
| **The Exception-Registry agent** | "must permit exceptions" read as "build the workflow" | F9 |
| **The Fixture agent** | hand-authored "trusted" ledgers in test data | F6/F7 |
| **The Recompute-Early agent** | builds the review engine instead of the substrate | F11 deferral |

---

## The patch (redline, in order)

1. **§2 add decision G:** strings are names, never carriers; relational objects require typed carriers; every constructor audited. (F2)
2. **§21:** inline the six criteria. (F1)
3. **§12:** add `ExpandSeries`/`Truncate(order)`, `Solve`, `Compare`; add the Integrate domain caveat or defer Integrate. (F4)
4. **§6:** "Prefer" → "Must"; manifest as JSON + runtime accessor + cross-validation tests. (F5)
5. **§14:** add replay validation + hash-chaining; extend minting to the ResearchCandidate. (F6)
6. **§17/§20M:** add the three micro-demos (candidate contamination, anomaly inspection, negative tests) as mandatory acceptance; seed the must-pass test list. (F3, F7)
7. **§2C:** explicitly scope populated packages to mechanics + SR. (F8)
8. **§2E:** reword to rule-versioning now, exception mechanism deferred. (F9)
9. **§2B/§20:** Identify signature requirements (operands + recorded justification + never via Simplify); exported `Eq`/`Hash`; size estimates requested; §20N deferral list extended. (F10, F11)

One substantive suggestion beyond the patches: make the **classical-relation demo do double duty** — derive the classical kinetic energy as the **second-order reduction of the relativistic KE** at v ≪ c. A single derivation then exercises truncation tagging, reduction-to-prior-framework (§10's recovery clause), and assumption propagation (nonrelativistic regime assumption minted by the limit) — which directly answers §17's demand that the slice "exercise the architecture rather than merely proving that arithmetic works."

---

## Closing verdict

As a prompt, this is a B+ instrument with A-grade intent: its decisions section will save the project another four rounds of architecture relitigation, and its discipline rule is the right instinct. Its failures are the corpus's familiar two, in prompt form — **it forgot to carry the highest-value enforcement rule (F2), and it demands components its own acceptance test never exercises (F3)**. Patch the nine items — an afternoon of editing — and the four reviewers will be auditing a plan rather than an agent's improvisations. The prompt's own §21 standard applies to it: a discipline rule is only as strong as its ability to catch the deferral it didn't anticipate, and right now the un-patched prompt defers exactly the mechanisms this project exists to enforce.
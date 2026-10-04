# Adversarial Review: The MVP Prompt + Normative Specs Pair

**Method.** The pair is audited against the corpus's accumulated findings ledger (the review trail this project has generated since design0), against each document's own instruments (the specs' §32 definition-of-done, §16's evaluation-anatomy requirement, §27's test suite; the prompt's "mechanical, no architectural decisions" demand), and against Go-specific execution realities. The key structural fact up front: **the prompt's discipline has finally outpaced the specs' determinacy.** With "specs win" and "do not say the coding agent can decide" both in force, every remaining spec ambiguity corners the plan agent — it cannot choose, cannot ask, so it will either stall or silently improvise. The residual defects have changed species: no longer *missing mechanisms* but *unresolved decision points*.

---

## 1. Resolution audit — the carried ledger

For the first time in this corpus, essentially every prior blocking finding is resolved, and traceably:

| Carried finding | Status | Where |
|---|---|---|
| Carriers-not-strings (v0.3 F2) | **Resolved** | §4.1 — "Strings MUST NEVER be used as carriers," with the exact valid/invalid pattern distinction |
| Session-minted provenance (v0.3 F1) | **Resolved** | §12.1, §22.1 (extended to ResearchCandidate) |
| Canonical form / exact rationals / structural Eq+Hash (v0.3 F3) | **Resolved in memory** | §5.1, §6.2–6.4 — wire encoding gap remains (H1 below) |
| Conventions first-class (v0.3 F5) | **Resolved** | §8 |
| Dimensions in core (v0.3 F7) | **Resolved** | §5 |
| Closed object model / impostor interfaces (v0.3 F6, F8) | **Resolved** | §2.2 (unexported-method interfaces), §4.2–4.3 |
| No stub packages (prompt3 F8) | **Resolved** | §2.1 |
| Seeded must-pass tests + scope accounting (prompt3 F7) | **Resolved** | prompt validation focus; §27 A–K; prompt item 15 |
| Six criteria dangling reference (prompt3 F1) | **Resolved** | §0 |
| Exception-registry trap (prompt3 F9) | **Resolved** | §10.2, §29 — versioning only, no bypass API |
| Identify contract (prompt3 F10) | **Resolved** | §11.2 |
| Replay + hash chain (prompt3 F6) | **Resolved** | §12.4–12.5 |
| Ledger-tamper & coverage micro-demos (prompt3 F3) | **Resolved** | Tests G, H, I, J, K; negatives 5–12 |

The 1905 → energy-momentum swap is a *declared* scope decision (§17.2, §31), not a silent amputation — the corpus's delta-ledger demand is finally satisfied by construction. This is the strongest artifact pair the project has produced. What follows is the residue.

---

## 2. Blocking findings — the plan cannot be mechanical until these are resolved

### B1 — The two-layer question is unanswered: nominal types vs `core.Expr`, and MRC-003 has no mechanism

§4.1 demands distinct nominal Go types (`Mass`, `Force`, `FourMomentum`…); §6.1 demands one exported immutable `core.Expr` handle. **Nothing specifies their composition.** Is `Mass` an `Expr`? Does it wrap one? What does `ops.Add` accept — nominal types, Exprs, or an interface both satisfy? If `Add` takes a common interface, `Add(Mass, Time)` *compiles* and fails at runtime — meaning the nominal-type layer does nothing at the expression level, and MRC-003 (category compatibility) silently collapses into MRC-002 (dimensions). If `Add` is nominal-typed, the expression engine can't combine quantities at all.

The distinction is real, not pedantic: torque and energy share dimensions (M L² T⁻²) — adding them is dimensionally valid and categorically meaningless. Test B ("reject invalid addition") and negative test 2 ("incompatible physical/category operation") are *different tests* only if MRC-003 has its own mechanism (e.g., a quantity-kind tag checked in ops, distinct from the dimension vector). The specs assert the rule ID (§10.1) and never specify the enforcement. Prompt item 3 (data model) and item 4 (MRC map) both hang on this. **Specify:** the composition pattern (nominal types as thin wrappers carrying an unexported `core.Expr` + kind tag + dimension), and that MRC-003 = kind-tag compatibility checked at operation entry, disjoint from MRC-002.

### B2 — Provenance contamination semantics are unspecified, and Test G is load-bearing on them

§19.1: "Candidate-derived results MUST remain candidate/provisional." By what rule? Test G ("derive with it, verify candidate contamination") cannot be written without a **provenance propagation law**, and the specs define propagation only for assumptions (§7.3). The needed rule is the dual of §7.3: *result provenance = f(input provenances, operation)* — specifically, `HYPOTHESIS` is contagious (any candidate input contaminates the output), `DERIVED` requires all inputs in trusted statuses, and no operation can *raise* provenance rank. Without this stated, the plan agent must invent the entire containment machinery — exactly the architectural decision the prompt forbids. One paragraph in §19 fixes it.

### B3 — §16's executable-reduction requirement is unsatisfiable as written, and it contradicts §0's MUST

§0 requirement 2: frameworks MUST be "corpus/evaluation targets." §16 makes evaluation-target status conditional on "at least one executable reduction/limit test." But the physically meaningful reductions (SR → Newtonian: K → ½mv², p → mv at v ≪ c) are **first-order-in-v/c statements requiring series expansion — deferred to v0.5** (§31.2–3). The `Limit` op cannot honestly produce them (0/0 forms; no L'Hôpital, no series). The only MVP-executable reductions are point evaluations — `γ → 1` as v → 0 — which are close to vacuous. As written, relativity *cannot* satisfy §16, so by the spec's own definition it is not an evaluation target, contradicting §0. **Resolve one of three ways, explicitly:** (a) name the honest-but-weak reduction (`LorentzFactor(v→0) = 1`, rest-frame consistency) and accept its weakness; (b) scope §16.2 to v0.5 alongside series algebra and amend §0.2's wording; (c) note that mechanics, having no predecessor, is exempt (also unstated — see L2). Silence here guarantees the four reviewers disagree about whether Test E's suite satisfies §16.

---

## 3. High findings

### H1 — Exact rationals have no wire encoding, and the prompt mandates a test the specs never include

The prompt's mandatory validation focus requires "exact rational round-trip." The specs put `big.Rat` in memory (§5.1, §6.2) and JSON at the boundary (§22.3) — but **JSON has no rational type**, and neither §25 nor §26 nor §27 contains a serialization round-trip test. This is the v0.3 review's float-smuggling finding, surviving at exactly the boundary the reviewer AI consumes. Specify the encoding (string-pair `["num","den"]`, or decimal-string rationals; floats banned from canonical JSON) and add the round-trip test to §25.

### H2 — Three operations are specified but exercised by no test

`Differentiate`, `Limit`, `Compare` (§11.1) appear in no canonical test (§27) and no negative test (§25). §32's definition-of-done checklist doesn't mention them either — so an implementation can ship them untested and be "done." This is the dead-spec-surface pattern the corpus has flagged in every round. The fix is one beautiful test that also *disambiguates B1*: **differentiate the momentum relation with respect to time and recover Newton's second law** — `d(mv)/dt = F` (Test A2). It forces the nominal-type/Expr composition to be concrete, gives `Differentiate` real coverage, and ties Test A to the op layer. `Limit` is covered iff B3 is resolved with option (a); `Compare` needs one use (e.g., branch-set comparison in Solve) or should be trimmed.

### H3 — Test K is unimplementable against `go:embed` as specified

§13.3 mandates loading manifests via Go's embedded-file mechanism — **compile-time immutable**. Test K requires validating a *mutated* manifest at test time. Contradiction, unless the validation logic is a **pure function over manifest bytes** (embedded bytes in production; corrupted bytes in Test K). One sentence in §13.4. Without it, the plan agent either hardcodes embed-only (Test K impossible) or abandons embed (violating §13.3).

---

## 4. Medium findings

- **M1 — Dead epistemic surface.** The `Approximation` assumption kind (§7.2's own example!) and the `APPROXIMATED` provenance status (§9.1) are specified and exercised by *nothing* — the residue of the 1905 deferral. Either add a micro-truncation test (truncate `1/(1−x)` to `1+x` under a declared `|x|≪1` assumption — no physics needed, purely to exercise the kind/status/tagging path) or move both to §31. Do not freeze them unexercised.
- **M2 — Solve semantics.** Test E's step 5 ("take the positive-energy root under an explicit precondition") is load-bearing and undefined. Specify MVP scope: Solve performs explicit isolation of one variable in polynomial-style relations, returns a branch set, and branch selection is an assumption-driven act (`Constraint, E, E > 0`), not a silent simplification.
- **M3 — In-package tests defeat the object model.** Go test files in `package mechanics` can write unexported fields — negative test 5 ("direct construction using forbidden fields") proves nothing if it lives inside the package. Require external test packages (`package mechanics_test`) for all construction-authority and immutability negatives.
- **M4 — Hash-chain concatenation is ambiguous.** `SHA256(canonical(step_without_currentHash) + previousHash)` — concatenating variable-length fields is injectivity-unsafe. Hash the *canonical serialization of a struct containing both* (the §6 canonical machinery already exists). Minor cryptographic hygiene, but this chain is the anti-tamper centerpiece.

## 5. Low findings

- **L1** — Assumption `Value` equality (§7.3's "identical collapse") needs canonical value forms or structurally identical assumptions won't collapse.
- **L2** — §16's four requirements are stated for "a framework" without exempting mechanics (no predecessor → no reduction target).
- **L3** — `Commit` (ledger) vs `Seal` (candidate) semantics never distinguished; one sentence each.
- **L4** — Who mints `HYPOTHESIS`? AI-callable candidate constructors (consistent with open-representation/closed-authority) — say so, since §22.1's session-only rule for `ResearchCandidate` otherwise looks contradictory.
- **L5** — Rendering/inspection from v0.3 is gone but absent from §31's deferral list.
- **L6** — Physics pedantry, for the record: deriving E = mc² by setting p = 0 in E² = (pc)² + (mc²)² is special-case extraction, not derivation of the relation from dynamics. This is acceptable *for substrate validation* — and the manifest's `derivable_from` field makes it honest — but reviewers should not mistake Test E for a physics demonstration; it is an architecture demonstration. The 1905 derivation remains the real one, correctly deferred.

---

## 6. Prompt-side findings

The prompt is in strong shape — seeded validation list, scope accounting, single-choice discipline, specs-win precedence, "no placeholder packages." Two additions:

- **P1 — An open-spec-item protocol.** B1–B3 are corners: the agent can't choose (prompt) and can't ask (one-shot). Add: *"If the specs under-determine a required plan element, enumerate it in a dedicated `Open Spec Items` section rather than improvising."* This converts silent divergence into visible flagging — the `UNKNOWN_TO_A` discipline, applied to the planning layer. (Better: patch the specs first, per the table below.)
- **P2 — Require a coverage matrix.** Mandate a table mapping **every normative MUST and every §10.1 rule ID → the test that exercises it** (prompt item 11 extension). This single artifact would have mechanically caught H1, H2, and M1 — the entire class of "specified but never exercised" defects — instead of requiring a reviewer to hunt them. Highest-leverage prompt edit available.

---

## 7. Patch list

| # | Fix | Doc | Blocking | Cost |
|---|---|---|---|---|
| 1 | Two-layer composition + MRC-003 kind-tag mechanism (B1) | specs §4/§6/§10 | **Yes** | ½ day |
| 2 | Provenance propagation law; HYPOTHESIS contagious; no rank-raising (B2) | specs §19 | **Yes** | ½ day |
| 3 | Name the MVP reduction or scope §16.2; exempt mechanics (B3) | specs §0/§16 | **Yes** | 1 hr |
| 4 | Rational wire encoding + round-trip test (H1) | specs §6/§22/§25 | **Yes** | 2 hrs |
| 5 | Test A2: d(mv)/dt = F; resolve Limit/Compare coverage (H2) | specs §27 | Yes | 2 hrs |
| 6 | Manifest validator as pure function over bytes (H3) | specs §13.4 | Yes | 15 min |
| 7 | APPROXIMATED micro-test or deferral (M1) | specs §7/§9/§31 | pre-freeze | 1 hr |
| 8 | Solve branch/selection semantics (M2) | specs §11 | pre-freeze | 1 hr |
| 9 | External test packages for authority negatives (M3) | specs §25 | pre-freeze | 15 min |
| 10 | Hash the canonical struct, not concatenation (M4) | specs §12.4 | pre-freeze | 15 min |
| 11 | L1–L5 clarifications | specs | pre-freeze | 1 hr |
| 12 | Open-spec-item protocol + coverage matrix requirement (P1, P2) | prompt | Yes | 30 min |

---

## 8. Closing verdict

Credit where the corpus has never given it before: **this pair resolves the entire carried findings ledger with no silent amputations** — carriers-not-strings, session authority, replay chains, exact rationals, conventions, dimensions, closed authority, no-stub discipline, the seeded test suite. The prompt finally enforces what four rounds of reviews were building toward: a plan agent with no room to improvise architecture.

Which is precisely why the remaining spec ambiguities are now the critical path. B1–B3 are not style issues; they are decision points the prompt forbids the agent from making and the specs never made themselves — the two-layer data model, the containment law, the reduction test. Left unresolved, the four reviewers will receive four plans that differ exactly there, and the review process will re-litigate in plan-review what should have been settled in spec. Patch the twelve items — a day of editing, most of it single paragraphs — add the coverage matrix, and this pair is ready to dispatch. The specs' own §32 standard says it best: *anything not required to make the loop work belongs outside the MVP* — but everything the loop's tests depend on must be inside the spec. Right now, three such things are not.
## Gemini

**Architectural Contradiction in Kernel Isolation (H0 vs. H2)**
Plan 7 mandates that `internal/kernel` must act as "theory-neutral semantic machinery" containing no theory-specific assumptions or physical laws. However, Gate H2 explicitly whitelists the inclusion of `KindMinkowski` and `LorentzFactorFunctionID` within this generic kernel. This creates a structural loophole: an attacker or developer could smuggle relativistic axioms—such as a specific metric signature or coordinate constraints—into the universal base layer under the guise of closed-world validation for these specific identifiers. This violates the core epistemic rule that the kernel must not contain the special-relativity corpus.

**Assumption Forgery via Weak Ops-Level Invariants (H3 & H4)**
Gate H3 allows operations to output an assumption set comprising input keys combined with the "operation's explicitly required/generated keys". Furthermore, Gate H4 specifically instructs the system to "never assert exact-key sets," relying only on preventing the intersection of mutually exclusive forbidden keys (e.g., mechanics vs. relativity). This presents a severe injection vector: a malicious operation inside the `hypothesis/` space could dynamically generate non-forbidden, tautological, or conflicting assumption keys. Because the system cannot assert exact-key sets, the `mergeInputs` function will blindly accumulate these fabricated assumptions, allowing a hypothesis to poison the assumption chain without triggering a package-boundary alert.

**Spoofing Epistemic Authority (Cryptographic Omission)**
The specifications strictly forbid a provisional hypothesis from automatically acquiring `ESTABLISHED` status. However, Plan 7 explicitly marks "authenticity crypto" as a non-blocker for this integration. Because the architecture relies on "byte-based replay" for safety without cryptographic signatures, an attacker can intercept the serialized ledger, manually execute a hypothesis derivation, and edit the raw JSON to label the output as `ESTABLISHED`. Without cryptographic provenance sealing the payload, the replay engine will deserialize and trust the altered bytes, bypassing the epistemic isolation rules entirely.

**Static AST Walk Evasion (H1)**
Gate H1 relies on an Abstract Syntax Tree (AST) walk located in `core/object_test.go` to enforce that the kernel imports zero local theory packages. This static check is easily bypassed in Go. An adversary could use the `reflect` package, `unsafe` pointers, or `//go:linkname` compiler directives to dynamically link or invoke functions from `mechanics` or `relativity` at runtime. The static AST parser will register zero illegal imports, yet the runtime isolation will be fundamentally breached.

**Unaddressed Runtime Panics (Replay Loop Vulnerability)**
The integration relies on a "no-panic battery" and strict serialization pins, such as `pow` rejection and JSON-number rejection, to secure the replay loop. However, restricting structural fixes to a "pinned explicit set" fails to address generalized state machine fragility. If an attacker feeds mathematically invalid arrays into operations like `Solve` or `SelectBranch`—bypassing the specific historical payloads tested in the Gate A battery—the system will still panic at runtime during the `Session.Validate` stage. The engine protects against known corrupted states but lacks comprehensive bounds-checking for novel, maliciously crafted AST structures.

## Z

# Adversarial Review: `plan7.md` (Gate H — Kernel Freeze) against specs v2.3 + Gate H origin

## Verdict

**Conditional GO.** Gate H is the right freeze at the right time: it locks a verified status quo (audit-corroborated), adopts my prior R3/R5 recommendations, and faithfully encodes the falsifiability-isolation principle — including the critical terminology correction that SR stays `ESTABLISHED`. But **one gate as written will fail on conforming code** (H2's allowlist mis-enumerates the kernel), the scanner mechanics have two false-positive traps, and four governance one-liners are missing. All fixable in an editing pass before implementation.

---

## 1. Fidelity to the Gate H origin (`specs_falsifiabilty.md`)

| Origin section | plan7 gate | Fidelity |
|---|---|---|
| A — SR outside kernel, identifier allowlist | H2 | ✅ (but see G1) |
| B — OS folder as unit of theory separation | H0 | ✅ recorded; ⚠ no documentation home pinned |
| C — ESTABLISHED ≠ unrestricted truth; no EBP in kernel | H0 + non-blockers | ✅ terminology correction preserved ("does not demote SR to HYPOTHESIS") |
| D — Hypothesis isolation (four "must not" bullets) | H4 sentence | ⚠ compressed; wording imprecise (G6) |
| E — No ambient assumption leakage, regression tests | H3 + H4 | ✅ both levels, non-leakage not exact-pinning — correct call |
| F — Import-boundary regression test | H1 | ✅ |
| G — No new directories during freeze | File discipline | ✅ zero new files, all placements legal |

The origin's most important nuance — *architectural isolation ≠ epistemic demotion* — survives intact in H0. That was the trap; plan7 avoided it.

---

## 2. 🔴 G1 — H2's kernel allowlist is mis-enumerated and fails on conforming code

H2: *"Non-test `internal/kernel` may contain only `KindMinkowski` and `LorentzFactorFunctionID`/`"lorentz_factor"` plus their closed-world validation."*

This is false about the current, audit-PASSED kernel. `kernel/types.go` must contain the **entire §6 Kind enum** — `KindMass, KindTime, KindForce, KindMomentum, KindVelocity, KindEnergy, KindSpacetime, …` (18 values) plus the `String()` mapping table emitting `"Force"`, `"Momentum"`, `"Energy"` as literal strings — all spec-mandated (REQ-006-01, stable ordinals). Add the §11.1 AssumptionKind, §13.1 ProvenanceStatus, §13.4 CorpusStatus, and §8.2.1 enums. Every one of these is physics-bearing text in the kernel.

The error's provenance is traceable: the corpus audit's FIND-002 said "the only physics-specific content is the spec-mandated `lorentz_factor` ID and `KindMinkowski`" — meaning the only *distinctive theory* content *beyond mandated machinery*. plan7 over-literalized that into an exhaustive two-item allowlist. As written, Gate H fails the exact code the audit passed, and the implementing agent "fixes" it by weakening the scan (whitelisting failures ad hoc) — the worst outcome.

**Fix:** H2 allowlist = *the spec-mandated closed enumerations (§6, §8.2.1, §11.1, §13.1, §13.4) + `LorentzFactorFunctionID` + its closed-world Call validation — and nothing else*: no laws, no formulas, no physics symbol constants, no framework assumptions. Same completeness nit for ops: allow the **body cluster** (`lorentzFactorBody` + `expandCall` + the Call-finiteness hook), not just "the fixed Lorentz body" — the dispatch and finiteness expansion are part of the same mandated machinery.

---

## 3. 🟠 Scanner mechanics — two false-positive traps

**G2 — Comment blindness is mandatory.** `ops/relation.go:180` contains the comment *"mass-energy path: Sqrt(Pow(m·c²,2)) → m·c²"*; `simplify.go` comments say *"primitive physical quantity."* An AST scan that includes comments trips H2 on current code. Pin: comment-ignoring scan (standard `go/ast` behavior — state it, since someone may use `go/parser` with comment mode).

**G3 — The body builder legitimately contains `"v"` and `"c"` literals.** "Reject `symbol == "c"/"m"`-style branches" needs a function-scoped allowance: `Symbol("v")`/`Symbol("c")` construction is legal *inside* `lorentzFactorBody` only. Pin the allowance per named function, not per file — otherwise either the gate fails or the exception swallows the gate.

---

## 4. 🟡 Invariant precision

**G4 — H3 is one-directional.** `output keys ⊆ input keys ∪ required` catches *invented* assumptions but not *mutated values* under an inherited key (same key, changed value passes a keys-only check). Add the second half of the spec's own merge law: output assumptions ⊇ merge(inputs) verbatim (value included). This may already exist in Plan 6's 12-op matrix — ⚪ verify; if not, add it.

**G5 — Reverse-allowlist invariant is mis-stated.** The pinned set `{NewKineticEnergy, Velocity, zeros…, hypothesis concept}` includes two constructors that do **not** yield ESTABLISHED (`NewKineticEnergy` is DERIVED/NONE; the hypothesis concept is HYPOTHESIS/NONE). So "every *other* ESTABLISHED-yielding constructor must map to a manifest item" exempts any future DERIVED-yielding smuggled constructor from mapping. Tighten to the status-independent bijection: *every exported constructor in mechanics/relativity/hypothesis ⊆ manifest constructor IDs ∪ pinned allowlist — period.* This is the drift-proof form of my R3.

**G6 — H4 scoping and wording.** (a) The checks must be scoped to **constructor outputs** — a mixed-framework ops derivation legitimately merges both frameworks' keys, and must remain legal. (b) "No ESTABLISHED material leaks into candidates" is ambiguous — candidates *contain* established premises by design. Restate in MRC-008 terms: *no hypothesis-derived step output carries `ESTABLISHED` corpus status or non-`HYPOTHESIS` provenance.*

**G7 — pow-only strict rejection.** The serialization pin locks strict empty-field rejection for `pow` only. My N6 was all-ops unused-field rejection. Either confirm Plan 6 covers the other eleven IDs or record pow-only as the accepted scope — silently asymmetric strictness is a replay-tamper surface.

---

## 5. 🟡 Governance gaps

- **G8 — No precedence header.** plan10_v2_3.md opens with "specs win"; plan7 has no normative-source statement. Add one line: *specs_v2_3.md is authoritative; plan7 is additive gates; conflicts resolve to spec.*
- **G9 — plan7.md's own location.** "Recorded as a new plan7.md file" must live **outside the module tree** (like plan10/specs), or it breaks the 39-file count and REQ-003-01. One sentence.
- **G10 — H0 has no home.** The locked epistemic wording should land in the README (with the existing docs-section assertion) — this simultaneously closes my prior recommendation to document the E=mc² demo's epistemic status ("derivation from trusted corpus, not a historical/physical proof"). Or explicitly defer; don't leave it homeless.
- **G11 — Two inherited dispositions are invisible.** FIND-001 (mechanics source strings — plan pinned `"Classical Mechanics corpus"`, code says `"Newton, Principia"` everywhere) and my R2 (corpus-status valid-value flip — flipping `established`→`contested` passes every detector the audit cited) are not addressed anywhere in plan7. They may live inside Plan 6's "corpus lock" — ⚪ verify; if not, they must enter the freeze, since a freeze that locks an unresolved pin divergence isn't a freeze.

---

## 6. What I verified holds (hand-checked)

- **H3 is true of current op semantics.** Per-op required-key table: `divide` → `denominator/*`, `select_branch` → `selected_branch/*`, all ten others → ∅. No op invents framework keys.
- **H5's ownership list matches the actual final object.** Re-derived the trace: final `m·c²` carries exactly `rest_mass_nonnegative` (← RestMass via EMR), `speed_of_light_positive` (← SpeedOfLight via EMR), `rest_frame` (← ZeroThreeMomentum via Substitute, surviving Simplify/Solve), `selected_branch/*` (op-generated). The constraint from `Compare(Energy, ZeroEnergy, gte)` adds nothing — and `Energy` is a legal ordered kind for `gte`. H5's list is exact.
- **Mutation methodology fixes the audit's R5 weakness** — disposable copy, `-count=1`, restore per mutant, "claim only rerun results." That converts the corruption battery from analytical to executed.
- **Non-blockers list is correct** — especially "1905 corpus" (consistent with the Einstein-derivation analysis: series expansion is banned, `APPROXIMATED` reserved) and "absolute zero-physics rule" (the two mandated identifiers make a literal zero-rule impossible).
- **File discipline holds**: every gate maps to an existing test file; H1/H2's AST infra reuse in `core/object_test.go` is legal (in-module test, `go/parser` in tests only). Minor option: the ops-side H2 scan could live in `ops`'s own test file — the plan's choice is compliant, just note it.

---

## 7. ⚪ Unverifiable in this packet

Plan 6 itself (gates A–G, the M1–M18 mutant inventory, the §17 A–E report format) was not provided — this review covers Gate H's deltas only. Specifically unverifiable: whether the corpus lock includes the corpus-status valid-flip detector and per-item `corpus_status` cross-check (R2), the FIND-001 disposition (G11), and whether value-preservation is asserted per-op (G4).

---

## Bottom line

**Proceed to implementation after applying G1** (the allowlist must describe the real kernel: mandated enums + the two distinctive identifiers) **and pinning G2/G3** into the gate text — those three are the difference between a gate that passes conforming code and one that trains the agent to ignore gates. Fold G4–G7 into the same pass; add the four governance one-liners (G8–G10, plus G11 verification). The gate design itself — epistemic rule without demotion, two-level non-leakage proof, ownership-traced E=mc², executable mutants, zero new files — is exactly the right freeze shape, and it demonstrably encodes the falsifiability-isolation principle you asked for without turning special relativity into a hypothesis.
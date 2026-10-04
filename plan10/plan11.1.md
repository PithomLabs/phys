# Plan 11.1 — GR Implementation-Plan Phase (agreed plan, no implementation yet)

Status: planning only. No code, corpus, spec, or tree changes made under
this plan. Awaiting explicit instruction before PASS 0 or PASS 1 executes.

Locked decisions (Q&A with maintainer):
- D1: Pins first (PASS 0), plan after (PASS 1) — two separate agent tasks.
- D2: Static-spherical workload anchor (GR-0…GR-8), index algebra embedded.
- D3: Level-2 libraries are separate Go modules; 42-file tree never touched.

## 0. Baseline (verified read-only before planning)

- 42-file re-freeze AUTHORIZED (adv_review11, DOCUMENTATION PASS).
- `specs_v2_3.md` byte-for-byte frozen; both manifests frozen.
- Pin A grounding: `Session.Identify` already mints `Kind: core.KindRelation`
  (`session/session.go:278`) — test pin only, zero production diff expected.
- Pin B grounding: `differentiateExpr` already fail-closes negative-integer
  exponents with `UnsupportedOperationError` — test pin only.
- Pin C grounding: `AGENTS.md:174` kind-checklist line overclaims (plan11_z
  A3) — one-line doc caveat, no file-count change.

## 1. Task decomposition

```text
PASS 0 — Conformance Pins (implement + verify + freeze baseline)
    ↓ hard gate: any pin failure blocks, never worked around
PASS 1 — GR Implementation Plan (document only; assumes Pass 0 done)
    ↓
GR workload begins under the plan
```

The GR plan must describe a verified substrate, not one being changed
while the plan is written. A planning agent with two jobs could smuggle an
implementation decision into the architecture — hence the split.

## 2. PASS 0 scope (separate task, executes first)

1. **Pin A test** (`session/session_test.go`):
   `Identify(named physical object, expression object, justification)`
   → `KindRelation`, `Expr = Relation(eq, …)`, `IDENTIFIED` provenance;
   plus `Identify(Energy, E/c², …)` must not yield `RestMass` (no covert
   semantic ascription through session authority).
2. **Pin B test** (`ops` tests):
   `Differentiate(Pow(x, negative-integer))` → `UnsupportedOperationError`.
   Do NOT expand differentiation to make GR easier.
3. **Pin C doc line** (`AGENTS.md` §8 checklist only):
   named-kind compatibility protects additive operations only while a named
   kind is retained; equal-dimension `Expression + Expression` mixing is
   permitted by design; semantic preservation is higher-layer work.
   Do NOT change MRC-003.
4. Run `go build ./...`, `go vet ./...`, `go test ./... -count=1` green;
   verify baseline frozen. Any non-conformance found → halt for human
   adjudication, never a silent fix or workaround.

## 3. PASS 1 scope (plan document only)

Produce `plan10/GR_IMPLEMENTATION_PLAN.md`. No GR code, no kernel changes,
no shared libraries, no spec edits, no corpus edits in this task.
Required sections (per task §14):

1. Executive Summary
2. Frozen Baseline and Non-Negotiable Invariants
3. Three-Level Architecture
4. GR Workload Objectives
5. Pre-GR Conformance Pins (records PASS 0 as satisfied prerequisite)
6. GR Package/Repository Structure
7. Representation Strategy
8. Staged GR Implementation Phases
9. Canonical GR Workloads / Worked Examples
10. Failure Taxonomy and Evidence Logging
11. Existing Six-Concern Ledger (C1–C8, routing frozen)
12. Level-1 → Level-2 Promotion Criteria
13. Level-2 → Level-3 Growth Gate
14. Silent-Wrongness Test Strategy
15. Mutation/Negative Testing Strategy
16. Canonicalization / Replay / Artifact Impact Analysis
17. Human Review / Independent Veto Points
18. Definition of Done for the GR Stress Test
19. Explicit Non-Goals
20. Expected Outcomes (NO-GROWTH / LEVEL-2-GROWTH /
    KERNEL-GROWTH CANDIDATE / ESCALATE-TO-SPEC)

## 4. Workload anchor (locked)

Primary thread — static, spherically symmetric spacetime:

```text
coordinates / chart → metric → inverse metric → Christoffel symbols →
covariant derivative → Riemann curvature → Ricci tensor/scalar →
Einstein tensor → vacuum field equation → Schwarzschild exterior →
weak-field / Newtonian limit
```

This forces the hard encounters honestly (`r⁻¹`/`r⁻²` negative powers,
`gμν`/`Γ`/`R` index identity, `∂g/∂x` differentiation demands,
`gμν·gνρ=δμρ` contraction, `Rμν=0` structured relations, Newtonian-limit
assumption/regime semantics). Index algebra is developed INSIDE this
thread (GR-2), never as a standalone generic system first. Linearized
gravity is reserved for the later controlled progression, not the entry.

Passes GR-0…GR-8 (representation baseline; coordinates+metric; inverse+
index ops; connection; curvature; Einstein tensor + vacuum equation;
Schwarzschild; weak-field limit; failure classification + Gate evidence).
Each pass records: objective, files/packages, dependencies, capabilities,
worked examples, tests, expected failures, evidence, level implicated,
exit criteria — and where the substrate bends, breaks, or silently
misrepresents.

## 5. Failure taxonomy (exact, one-of-five per obstacle)

SPEC-INTENDED-BOUND / PACKAGE-SOLVABLE / REPRESENTABLE-BUT-UNFAITHFUL /
UNREPRESENTABLE / SILENTLY-WRONG, each logged with requirement, minimal
example, current vs expected result, category rationale, Level-1/2
workarounds, Level-3 implication. No wishlist entries. Silent-wrongness
(symbol ambiguity, contraction/raising mistakes, dimension-valid but
semantically invalid ops, type decay, inconsistent-assumption
coexistence) is the primary adversarial thread.

## 6. Routing and gates (locked)

- C1/C2/C4 → Level-2 investigation (never kernel, never weaken Substitute).
- C3 → Level-2 constraint/entailment layer. C5/C6 → future/closed.
- L1→L2: named second theory + worked example + shared semantics.
- L2→L3: proven L2 failure to preserve a generic invariant + full 20-field
  Growth Evidence Record. Burden: generic, theory-neutral, non-semantic,
  foundational, reusable, invariant-bearing, unlivable above kernel.
- Default outcome: NO KERNEL GROWTH.

## 7. Locked module rule (D3)

Level-2 libraries are separate Go modules (e.g. `physmath`) depending
one-way on frozen `phys`; `phys` never depends back. The 42-file tree is
never modified by promotion. Only a proven Level-3 primitive touches
`phys`, via the full Growth Gate. Kernel invariants (§12 list:
immutability, mint authority, no generic factory, exact rationals,
canonicalization, ordinal stability, MRC, provenance, containment,
session/ledger/replay/hash integrity) remain untouched by all GR work.

## 8. Review inputs reconciled

plan11 (triage: limitations not defects) + plan11_z (A1 pin Identify kind;
A2 pin negative-pow; A3 kind-doc caveat; A4 ladder governance; A5 citation
hygiene) + plan11_review (KindRelation over Expression; C1–C8 ledger;
workload-not-redesign) + deepseek/qwen/gemini convergence on the ladder.
Net-valid conclusions only; plan11_review's pre-GR lock is the starting
position: test expressiveness of the kernel, never pre-expand it.

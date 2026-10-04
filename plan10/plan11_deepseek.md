## Integration: three-level promotion + six-concern triage

The review's conclusion is correct and worth locking: **the six criticisms do not, by themselves, produce a kernel-growth candidate.** But two of them are genuine boundary findings that the gate must now know how to route. The routing is exactly the three-level promotion model you proposed. Below is the amendment to `KERNEL_GROWTH_GATE.md` that encodes it, plus a triage that binds each of the six concerns to a specific gate destination.

---

### 1. Amendment to the Gate — promotion is three-level, not two

Replace §8's "candidate / no-growth" binary with an explicit ladder. Every admissible proposal must first be routed to a level *before* G1–G7 are applied.

```text
LEVEL 1  Theory-local userland
         general_relativity/, quantum/, hypothesis/, etc.
         Imports core + kernel. Cannot modify either.

             ↓ requires demonstrated reuse across ≥2 theories

LEVEL 2  Shared mathematical / physics library
         tensor/, manifold/, hilbert/, index/, etc.
         Imports core + kernel. Cannot modify either.
         Human-curated promotion only.

             ↓ requires a generic invariant the trusted substrate itself
               must enforce, or it cannot be enforced safely at all

LEVEL 3  Kernel primitive (internal/kernel)
         Subject to G1–G7 in full.
```

Rules:

```text
R9   A proposal is first routed to a level. G1–G7 apply only at Level 3.
R10  Level 1→2 promotion is human-curated. "Two theories both use it"
     is a candidate signal, not a promotion.
R11  Level 2→3 promotion requires a Level-2 attempt that demonstrably
     FAILS to enforce the invariant without kernel support (G5).
R12  GR, QM, QFT never promote directly to Level 3. They supply failures
     that may later justify a Level-2 artifact, which may later justify
     a Level-3 candidate.
R13  Absent a Level-3 candidate, the gate's default outcome remains
     NO KERNEL GROWTH.
```

The Go analogy survives cleanly under this ladder: `sync.Mutex` is Level 3; `encoding/json` is Level 2; your application's `OrderBook` is Level 1. Kernel growth means admitting a new `sync.Mutex`, not admitting JSON.

---

### 2. Triage of the six concerns against the gate

Each concern is routed to exactly one destination. None reaches Level 3 on the evidence in this review.

| # | Concern | Destination | Gate status |
|---|---------|-------------|-------------|
| 1 | Type decay (`Expression + Expression` loses nominal kind) | **Level 2 candidate** | Real gap; not yet kernel. Needs GR *and* QM workloads to demonstrate a shared invariant, not merely shared discomfort. |
| 2 | Substitution catch-22 | **Level 2 candidate** — *semantic ascription* | Kernel is correct to refuse. Do not weaken `Substitute`. The missing capability is a controlled ascription layer above the kernel. |
| 3 | Assumption semantic conflicts | **Level 2 candidate** — assumption logic | Kernel already provides the integrity boundary. A solver/entailment layer belongs above it. |
| 4 | Symbol identity collision | **Level 2 candidate** — namespaced symbol layer | Kernel AST stays syntactic. Entity identity is a higher-layer abstraction. |
| 5 | Transcendental / symbolic-exponent hazard | **Not applicable to MVP** | Not expressible through current APIs. Record as a *pre-commitment* for whenever `Pow` or `Call` are extended. |
| 6 | Closed `Kind` ontology / no subtyping | **Level 3 discipline, correctly closed** | No action. Adding physics nominal kinds to the kernel is exactly the inversion R7/G2 forbids. |

The pattern is: **five of six concerns point upward to Level 2, one is already correctly closed, and none points to Level 3.** That is precisely the outcome the gate was designed to produce — and it is a *positive* result of the gate, not a null one.

---

### 3. The one concern that is genuinely kernel-adjacent

Concern 1 (type decay) deserves separate treatment because it is the only one whose failure mode is *already observable* in the frozen MVP:

```text
E_expr = m × c²        → KindExpression, dim M L² T⁻²
τ_expr = r × F         → KindExpression, dim M L² T⁻²
Add(E_expr, τ_expr)    → passes MRC-002 and MRC-003
```

That is a real semantic hole at the kernel boundary — two dimensionally identical, physically distinct derived quantities can be added. But note carefully:

```text
The failure is not "we lack a Kind for energy-and-work."
The failure is "Expression has no ascription mechanism, and the kernel
              correctly refuses to invent one."
```

So concern 1 does **not** become a Level-3 candidate by passing G1–G7. It becomes a Level-2 candidate: a *semantic ascription layer* that lets a theory package say "treat this `Expression` as `Energy` under justification J, recorded in provenance" — without the kernel manufacturing physical meaning. That is the correct long-term home, and it preserves the MVP's most important discipline (already encoded in `AGENTS.md` §2 and specs §…):

> **Algebra does not manufacture physical meaning.**

The gate's correct response to concern 1 is therefore:

```text
Outcome B-adjacent: record as a Level-2 roadmap item,
not a kernel candidate. GR and QM workloads must each produce the
failure in situ before the ascription layer is specified.
```

---

### 4. What this changes about the gate's inputs

Two new evidence requirements follow from the triage:

```text
E1  GR workload log must record, per failure, which LEVEL the failure
    points at (1, 2, or 3) — with justification. Level-3 claims without
    a Level-2 counter-attempt are inadmissible per R11.

E2  The gate maintains a running "concern ledger" (this table) so that
    future proposals cannot silently re-litigate a concern already
    routed to Level 2.
```

`E1` prevents the failure log from becoming a wish list. `E2` prevents the same six concerns from reappearing dressed as kernel proposals after GR begins.

---

### 5. Where this leaves the architecture

The review's bottom line survives the gate unchanged and is now formally encoded:

```text
kernel robustness                 : strong        (Level 3, no growth)
MVP symbolic capability           : intentionally narrow (Level 3 frozen)
multi-theory math expressiveness  : insufficient, extensible at Level 2
evidence for kernel growth        : none from this review
```

Two forward items, both Level 2, both human-curated:

```text
L2-1  Semantic ascription layer (addresses concerns 1 and 2 together)
L2-2  Namespaced symbol / entity identity layer (addresses concern 4)
```

Concern 3 folds into L2-1's justification machinery. Concern 5 is a pre-commitment note. Concern 6 is closed.

The next concrete action is not to design any of these. It is to run the GR workload with `E1` in force, let it produce failures tagged by level, and only then decide whether L2-1 or L2-2 is the first shared library worth extracting. That is the gate working as intended: **GR supplies the stress; the ladder decides the response; the default remains no kernel growth.**
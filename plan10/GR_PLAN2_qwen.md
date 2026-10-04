# Adversarial Review of `GR_IMPLEMENTATION_PLAN_20261002.md`

## Executive Verdict: CLEARED FOR HANDOFF
The revised plan is **architecturally bulletproof**. It has successfully internalized every finding from the previous adversarial review cycle. The explicit resolution of the `core.Object` bridge (Section 6), the precise handling of the differentiation boundary (Section 7), the correction of the silent-wrongness semantics (Section 12), and the enumeration of the 20-field evidence record (Section 15) completely eliminate the ambiguities that previously existed. 

The 18 consolidated decisions (D1–D18) are strictly enforced throughout the document. There are no remaining architectural blockers, and the coding agent has zero latitude to invent alternative architectures.

However, to make this document **100% immune to coding-agent hallucinations** regarding determinism and data translation, there are **three micro-ambiguities** that must be explicitly pinned before handoff. These do not change the architecture; they merely close the final gaps in the execution contract.

---

## 1. Verification of Previous Review Points (All Resolved)

| Previous Finding | Resolution in v2 Plan | Status |
| :--- | :--- | :--- |
| **GR ↔ `core.Object` bridge was implicit** | Section 6 now explicitly defines the `GRExpr → ToCoreExpr() → hypothesis.NewCandidateConcept() → core.Object` pipeline. | ✅ **LOCKED** |
| **Trig/Unknown functions undefined** | Section 6 & 9 explicitly mandate real `Sin`, `Cos`, and `UnknownFunction` nodes in the GR-local algebra. No placeholders. | ✅ **LOCKED** |
| **GR-3a/3b wording contradiction** | Section 7 explicitly states GR-3b "necessarily takes over" after GR-3a fails; independence is the correct relation. | ✅ **LOCKED** |
| **Silent-wrongness semantics flawed** | Section 12 correctly classifies bare `Symbol("f")` differentiation to `0` as `ACCEPTED-CORRECT`, reserving `SILENTLY-WRONG` for actual semantic failures. | ✅ **LOCKED** |
| **Level-2 outcome contradiction** | Section 13 & 18 replace `LEVEL-2-GROWTH` with `LEVEL-2-CANDIDATE-PENDING`, requiring a named second consumer. | ✅ **LOCKED** |
| **20-field evidence record undefined** | Section 15 explicitly enumerates all 20 fields, including `MinimalKernelChange` and `HumanApproval`. | ✅ **LOCKED** |
| **Mutation testing protocol vague** | Section 17 explicitly mandates using a "disposable clone/worktree" and forbids mutating the authoritative `phys` checkout. | ✅ **LOCKED** |

---

## 2. The Final Three Micro-Ambiguities (Inject Before Handoff)

While the architecture is perfect, a coding agent might trip over the **determinism invariants** and **structured data translation** at the exact boundary where GR-local math meets the frozen kernel. 

Inject the following three clarifications into the plan to close these gaps.

### Micro-Ambiguity 1: Deterministic IDs for Bridged Objects
**The Risk:** When the GR package bridges a scalar expression to a `core.Object` via `hypothesis.NewCandidateConcept(id, ...)`, the agent might use a random UUID, a timestamp, or a non-deterministic string for the `id`. This violates the strict determinism requirement of the system (spec §31).
**The Fix:** Add to **Section 6 (Representation Strategy)**:
> *"The `id` parameter for bridged candidate concepts MUST be generated deterministically (e.g., the lowercase hex SHA-256 hash of the canonical `core.Expr` bytes) to preserve the system's global determinism invariant. No random UUIDs, timestamps, or non-deterministic counters are permitted."*

### Micro-Ambiguity 2: Structured Assumption Translation
**The Risk:** The GR-local algebra will naturally track assumptions like `r > 0` or `M != 0` as strings or raw expressions. When bridging to `core.Object`, the agent must pass a `core.AssumptionSet`. If the agent passes raw strings, the kernel's bounded entailment and MRC rules will fail to recognize them, breaking the derivation.
**The Fix:** Add to **Section 6 (Representation Strategy)**:
> *"When bridging to `core.Object`, GR-local assumptions (e.g., `r > 0`, `M != 0`) MUST be translated into structured `core.Assumption` values (e.g., `Constraint` or `MathPrecondition` with `ExprValue` containing a canonical `core.Expr` relation). The kernel's bounded entailment and MRC rules operate strictly on these structured values, not raw strings."*

### Micro-Ambiguity 3: GR-Local Canonicalization Rules
**The Risk:** The plan mandates a GR-local replay trace (`GRStep` with SHA-256 hashes). However, it does not explicitly state that the GR-local structures (`GRExpr`, `Tensor`, `IndexSlot`) must follow the same strict canonicalization rules as the kernel. The agent might use Go `map` iteration to sort tensor components, breaking the replay hash.
**The Fix:** Add to **Section 8 (Level-1 GR Package Design)** under the Tensor Engine or Section 9 (GR-Local Mathematics):
> *"The GR-local symbolic algebra (`GRExpr`, `Tensor`, `IndexSlot`, etc.) MUST implement deterministic canonicalization identical in spirit to the kernel: no Go map iteration for ordering, fixed struct field ordering for JSON serialization, and exact rational string formatting. This ensures the GR-local replay trace and artifact hashes are strictly reproducible across runs."*

---

## 3. Final Handoff Status

With the injection of those three defensive clarifications, the `GR_IMPLEMENTATION_PLAN_20261002.md` transitions from "architecturally sound" to **flawlessly executable**. 

The coding agent now has:
1. A strictly bounded frozen kernel (`phys`).
2. An explicit, type-safe bridge protocol for GR-local math.
3. A clear failure taxonomy and evidence-gathering protocol.
4. Zero ambiguity regarding what belongs in Level 1, Level 2, or Level 3.

**You are cleared to hand off to the coding agent.** The GR stress test will now definitively prove the boundary of the kernel without accidentally compromising its integrity.
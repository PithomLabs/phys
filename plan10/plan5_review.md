Plan 5 is **very close, but I would make two technical corrections before execution**. The overall strategy is sound: tests-only, exact M1–M18 roster, reproduce the survivors first, then add minimum proof tests. :chatgpt-content-reference{index="0"}

## 1. Fix M13 — use `add` with a `pow` params object

Current Plan 5 proposes:

> valid `pow` step → tamper `ParamsCanonical` to `Kind:"empty"` → recompute chain. :chatgpt-content-reference{index="1"}

I would **change this**.

The problem is that if Stage 10 is removed, the later replay stages may still notice that the altered Pow parameters produce a different result. That makes the test potentially prove **Stage 14 replay-output comparison**, not specifically Stage 10.

Use an operation whose implementation does not consume those parameters.

### Better M13 fixture

Create a valid:

```text
Operation = "add"
inputs = [a, b]
```

Then replace `ParamsCanonical` with a **structurally valid but operation-incompatible** parameter object:

```json
{
  "kind": "pow",
  "exponent": "2/1",
  "operator": "",
  "justification": ""
}
```

`OperationParams` is structurally valid, but `pow` parameters are illegal for `add`. The v2.3 binding explicitly requires `empty` for operations other than `pow`, `compare`, and `identify`. :chatgpt-content-reference{index="2"}

Then:

```text
recompute current-step/chain hashes
→ stages 1–9 pass
→ Stage 10 must reject
```

If Stage 10 is removed:

```text
ops.Apply(add, ...)
→ Add ignores irrelevant params
→ same output
→ replay-output comparison also passes
```

Therefore the mutant **can only be killed by the Stage 10 binding check**.

That is a much stronger mutation test.

Change M13 in Plan 5 to:

> `TestValidateOperationParamsStage`: commit an `add` step, replace `ParamsCanonical` with a canonically valid `Kind:"pow"` parameter object, recompute the step/chain hashes, verify `ParseOperationParams` succeeds, and assert `Session.Validate()` rejects the operation/parameter mismatch. Verify `Ledger.Validate()` accepts the structurally valid ledger. Removing Stage 10 must make the test pass.

---

## 2. Tighten M12 — isolate Stage 17 from Stage 15

Plan 5 currently allows:

> “output carries trusted status **or non-NONE corpus**” :chatgpt-content-reference{index="3"}

The **trusted-provenance** variant is potentially caught earlier by Stage 15, because Stage 15 is explicitly responsible for provenance/MRC validation. The v2.3 distinction is:

```text
15 → provenance propagation / MRC version
17 → candidate-containment invariants
```

:chatgpt-content-reference{index="4"}

For the mutation test, you want something that Stage 15 legitimately allows but Stage 17 must reject.

I would therefore make the primary M12 fixture:

```text
HYPOTHESIS dependency remains HYPOTHESIS
but retained output has CorpusStatus != NONE
```

while all provenance propagation remains correct.

That targets the artifact-level containment rule:

> a hypothesis-dependent output with corpus status other than `NONE` must fail candidate containment. :chatgpt-content-reference{index="5"}

Then prove:

```text
stages 1–16 → pass
stage 17 → CandidateContainmentError
```

The existing Plan 5 requirement to assert a containment-specific message and verify `Ledger.Validate()` does not produce the containment verdict is good. :chatgpt-content-reference{index="6"}

I would explicitly say:

> **Do not use the trusted-provenance mutation as the primary Stage-17 isolation fixture unless Phase 1 demonstrates that Stage 15 does not already reject it.**

That keeps the test causally precise.

---

## 3. M18 is correct

This part is good:

```text
Limit(Add(x,2), x, 3) → 5
```

It directly proves that `Limit` performs actual direct substitution/simplification rather than returning `1`. :chatgpt-content-reference{index="7"}

That is consistent with v2.3: `Limit` is direct substitution followed by simplification, while the Lorentz-specific test must separately produce exactly `1`. :chatgpt-content-reference{index="8"}

So **do not change M18**.

---

## One more small improvement

Phase 1 is exactly the right idea:

> reproduce M12/M13/M18 **before** adding tests. :chatgpt-content-reference{index="9"}

I would strengthen it slightly:

> When a mutant survives, determine **the first validation stage reached before the test fails**. Record that stage. This is necessary to prove that the new test isolates the intended stage rather than merely causing some later failure.

That makes the mutation campaign much more rigorous.

---

## Verdict

I would classify Plan 5 as:

```text
Architecture:        ✅
Strategy:            ✅
M18:                 ✅
M12:                 ⚠️ tighten stage isolation
M13:                 ❌ change fixture
Mutation roster:     ✅
18/18 target:        ✅
```

So **don't implement Plan 5 verbatim yet**.

Make these two edits:

```text
M12 → primarily corpus-status containment, demonstrated to survive stages 1–16
M13 → add operation + incompatible but structurally valid pow params
```

Then Plan 5 is ready.

The underlying principle is important: **each surviving mutant needs a test whose failure is attributable to the exact semantic invariant that mutant removed**, not merely a test that happens to fail somewhere later in the pipeline. That's the right standard for calling the 18-mutant campaign meaningful.
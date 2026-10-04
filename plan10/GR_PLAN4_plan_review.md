I read the attached **`GR_PLAN4_plan.md`** carefully. It is materially stronger than the earlier 20261004 hardening-only file: it correctly turns the review findings into a build specification for the **full 20261005 contract**, rather than trying to execute a 129-line patch list as though it were the contract. GR_PLAN4_plan

The major corrections are correctly captured:

- the **`Subst` contradiction is fixed** by allowing Symbol or zero-order `UnknownFunction` targets with arbitrary `GRExpr` replacements and derivative closure through order 2; GR_PLAN4_plan
- GR-8 `Identify` is correctly changed to **`Relation / HYPOTHESIS / NONE`**, consistent with the frozen contamination law; GR_PLAN4_plan
- the undefined `Reduce` operation is removed; GR_PLAN4_plan
- canonical inventory, chart propagation, operation namespaces, AST invariants, adversarial-test placement, and the flat-space check are all pinned; GR_PLAN4_plan
- DeepSeek's replay/certification details are incorporated, especially **candidate → residual → substitution → Normalize → ZeroTest**, with `UNDECIDED` failing closed; GR_PLAN4_plan
- the useful Gemini implementation discipline has been incorporated without importing its rejected architecture; the plan explicitly rejects the `int64` Rational pin, generic sorted-key serialization, the incorrect `HYPOTHESIS` BridgePayload example, and tree-reopening proposals. GR_PLAN4_plan

### Two things I would still correct in the plan before execution

**1. Clarify the authority sentence in §1.**

This line:

> “Follow procedures in `GR_PLAN4_review.md` (this file's instruction)…”

is ambiguous because the uploaded file is `GR_PLAN4_plan.md`, not `GR_PLAN4_review.md`. GR_PLAN4_plan

Use:

```text
Use this file as the execution-building instruction.
The attached review files are evidence/guidance only.
The final 20261005 contract itself is the sole normative execution contract.
```

That eliminates any possibility that the coding agent treats one of the review documents as authoritative.

**2. Make the embedded Chart/ConventionSet rule explicit.**

The plan says:

> “R6/R7 reconcile Index vs IndexSlot; Chart/ConventionSet embedded.” GR_PLAN4_plan

That is directionally correct, but “embedded” should specify **where**. Otherwise the Appendix D schema inventory can still leave the agent wondering whether standalone encoders are required.

Pin:

```text
Chart and ConventionSet have no standalone canonical artifact schema.
When present in a GR artifact, they are serialized using their fixed
embedded representation as specified by that artifact's Appendix D schema.
They are never independently hashed as GR artifact kinds.
```

That closes the residual R7 ambiguity cleanly.

### One bookkeeping issue

The plan says the new build scope is based on **28 hardening items**, while also adding the Z/DeepSeek/Gemini/Qwen-derived procedures and the 15-rule discipline section. That's fine conceptually, but don't let the final audit count those as “28 items” unless the numbering actually does so. The earlier review already caught this kind of counting mismatch. The 20261005 contract should have one authoritative item count and one authoritative checklist count. The plan already says those counts must be verified only after the audit passes. GR_PLAN4_plan

With those two wording pins and the count audit, **I would proceed to generate `GR_IMPLEMENTATION_PLAN_20261005.md` from this plan**. The architecture itself does not need another redesign.
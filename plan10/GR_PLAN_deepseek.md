# Adversarial Review — GR_IMPLEMENTATION_PLAN.md

The revision is a substantial improvement over `plan11.1.md`. The three-level ladder, the C1–C8 ledger, the GR-3a/GR-3b split, the enumerated 20-field evidence record, and the corrected Pin A fixture all address the earlier review directly. What follows is only what still breaks, ordered by severity. Every finding is new — resolved items are not re-litigated.

---

## A. Critical — the plan contradicts itself, or is not implementable as written

### A1. §13 and §18 are in direct conflict about whether LEVEL-2 can occur

§13 defines Level-1 → Level-2 promotion as requiring:

> "shared abstraction + **named second established consumer (not hypothetical)** + worked second-theory example + demonstrated common semantics + assumption/provenance non-leakage."

During the GR workload, **no second established theory exists**. GR is the first serious workload. QM/QFT do not yet exist as established consumers. Therefore no Level-2 promotion is possible during this plan's execution.

But §18 lists as an admissible terminal outcome:

> "LEVEL-2-GROWTH — Demonstrated cross-theory abstraction (named 2nd consumer + worked example); new shared module above kernel."

These two sections cannot both be true. Either:

- §18 must remove LEVEL-2-GROWTH from the admissible outcomes for this plan (and defer it to a later plan whose execution spans GR *and* QM/QFT), or
- §13 must relax the "named second established consumer" requirement to allow a *placeholder* or *fixture* second consumer (with an explicit note that promotion at this stage is provisional).

As written, an implementing agent reaching §18 will either fabricate a second consumer or silently expand §13. Both are violations.

This is the most consequential gap in the revision.

### A2. The interop contract between GR-local math and `core.Expr` is undefined

§6.2 says:

> "Scalar rational algebra uses `core.Expr` … Everything with index/variance/chart/frame semantics is a GR-local Go struct holding `core.Expr` components plus explicit index metadata."

§9 says GR-local math must provide:

1. a bounded differentiator (returns… what?),
2. trigonometric functions — `sin(θ)`, `cos(θ)` (returns… what?),
3. index/tensor algebra,
4. inverse metric, contraction, etc.

But `core.Expr` is **closed** (§3.5) and has no `sin` node, and `Call` admits only `lorentz_factor`. So a GR-local `sin(θ)` cannot be a `core.Expr`. It must be a GR-local type.

Now consider a metric component like `g_φφ = r² sin²θ`. Per §6.2 this is "a GR-local Go struct holding `core.Expr` components." But `sin²θ` is not a `core.Expr`. So the component must hold either:

- a mixed representation (some fields `core.Expr`, some GR-local), or
- a GR-local math type that wraps `core.Expr` for the parts it can, or
- a GR-local math type that is disjoint from `core.Expr` entirely.

§9 does not say which. Without this, §9 is not implementable: the implementing agent cannot define `symbolic/differentiate.go` or `symbolic/trig.go` without knowing what type they return and how those returns compose with `core.Expr` in the metric, Christoffel, and curvature structs.

This is a design gap, not a doc gap. §9 must state the unifying type (or declare that GR-local math is a *separate* algebra with an explicit conversion layer to `core.Expr`).

---

## B. Significant — should be fixed before PASS 0 is authorized

### B1. §7's allowed-operations table overstates the surface actually used by the workload

§7 permits `Compare`, `Solve`, `SelectBranch`, `Session.Identify`, and `Limit`. §5's workload (GR-0…GR-8) exercises none of these in an obvious way:

- `Compare` builds relation artifacts — used where?
- `Solve` accepts only `Relation(eq, t², rhs)` — used where in GR?
- `SelectBranch` needs a `gte`-against-zero constraint — used where?
- `Session.Identify` — identifies what against what? GR doesn't mint `core.Object`s (§6.4), so what are the operands?
- `Limit` — one Call body (`lorentz_factor`); GR-8's Newtonian limit is not a `Limit` on the Lorentz factor.

Either the workload must be extended to actually exercise these operations, or §7 must be trimmed to "allowed in principle, not exercised by this workload." Otherwise the table reads as planned usage that will silently not occur, and the Definition of Done (§20) will not catch the gap.

### B2. §7's GR-3b clause is self-contradictory

> "GR-3b must not modify `phys.Differentiate`, add negative-power support to the kernel, create a hidden kernel extension, or **claim kernel-level GR differentiation support**."
>
> "`symbolic/differentiate.go` … **does not bypass `phys.Differentiate`**."

If GR has its own differentiator, then for GR purposes the local one *is* used instead of `phys.Differentiate`. That is bypass by ordinary meaning. The intended rule is: GR-3b must not *modify*, *replace*, *shadow*, or *re-export* `phys.Differentiate`, and must not be presented as kernel capability. "Does not bypass" invites the implementing agent to invent a usage where GR-3b somehow calls into `phys.Differentiate` — which is impossible for negative exponents (Pin B fail-closes). Reword.

### B3. §3.3 asserts PASS 0 edits stay within the 42-file tree without any verification step

> "PASS0 adds content to existing authorized files only; exact 42-file membership is unchanged."

Three files are edited (§16): `session/session_test.go`, an `ops` test file, `AGENTS.md`. The plan asserts — but does not verify — that these three are (a) already in the 42-file set, and (b) the *only* files PASS 0 edits. If the `ops` test file is new, or if `AGENTS.md` is not counted in the 42, the membership claim is false.

Fix: add to §16 a one-line membership check:

```text
PASS0 verifies, before any edit, that the three target files are
members of the frozen 42-file set; a non-member target halts for
human adjudication (never a silent add).
```

### B4. §15 field 18 is undefined

> "minimal kernel change (**exact diff-shape**, no implementation)"

"Diff-shape" is not a defined artifact format. Does it mean:

- a `git diff` against `internal/kernel` (which would be implementation),
- an interface-signature delta,
- a text description of the change surface,
- a structural diff against the existing kernel `Object`/`Expr` types?

Without a format, this field will be filled differently by every implementing agent, and the evidence records will be incomparable. Define the format (e.g. "an enumerated list of affected exported identifiers, with the *kind* of change — added signature, changed invariant, new constructor — and no code").

### B5. §15 field 20's "human approval" is undefined

> "human approval (signature/decision)"

What is a signature here? A commit message line, a GPG-signed commit, a filename convention, an out-of-band record? The plan requires human approval as a terminal field of a Growth Evidence Record and offers no artifact format. Define it, or state that records without approval are provisional and cannot trigger a Level-3 change.

### B6. §18's "three false readings" is itself confusing

> "Prevent three false readings: 'kernel passed GR' (false — userland did work), 'kernel failed, therefore grow it' (premature), 'Level 1 worked around it, therefore kernel had no limitation' (also false)."

The first is described as *false*, yet §1 states:

> "`NO-GROWTH` is an acceptable successful outcome."

If NO-GROWTH means the kernel was sufficient for GR, then "kernel passed GR" *is* the outcome, not a false reading. The intended distinction is subtler: the kernel did not *do* the substantive work; userland did, and the kernel merely provided just enough primitives. But as written, the reader is told to avoid a statement that is literally the target outcome.

Reword along the lines of:

```text
Prevent three misreadings:
1. "Kernel passed GR" — meaning "kernel sufficed". Correct as a
   boundary statement. NOT "kernel did the work".
2. "Kernel failed, therefore grow it" — premature; failure is not
   evidence until the L1/L2 attempts fail.
3. "Level 1 worked around it, therefore no limitation" — false;
   the workaround is evidence of a boundary, not its absence.
```

### B7. §4's "Level-1 theory packages live outside frozen phys" has unaddressed grandfathering

> "Level-1 theory packages follow the same one-way dependency discipline as Level-2 libraries: they live outside frozen `phys`, depend only on the public `phys` API, and never modify the kernel tree."

But `mechanics/` and `relativity/` already live *inside* `phys`. They are Level-1 theory packages that predate this rule. The plan does not say:

- whether they are grandfathered,
- whether they must eventually be extracted,
- or whether the rule applies only to *new* Level-1 packages.

Without this, an implementing agent will either treat the rule as retroactive (which requires modifying the frozen tree) or ignore it for new packages that happen to be convenient to place inside `phys`. State the scope: "This rule applies to new Level-1 packages added after this plan; existing `mechanics/` and `relativity/` are grandfathered."

### B8. §12 probe 2 misclassifies user error as substrate silent-wrongness

> "field `f(t)` represented as bare `Symbol('f')`; `Differentiate(f-expr, t)` yields `0` … → `SILENTLY-WRONG` candidate."

The substrate is behaving correctly: a bare symbol is a constant with respect to any other symbol. If the user chose to represent a function as a bare symbol, that is a *representation choice* — arguably a documentation gap in the GR-local math, not a substrate defect. Labelling this SILENTLY-WRONG will pollute the taxonomy: it turns every misuse of a syntactic AST into a kernel finding.

The probe is worth running, but its classification should be "userland representation error" unless there is a *plausible innocent representation* under which the substrate silently produces a wrong result. Reword §12 probe 2 to distinguish "misuse available" from "misuse unavoidable given the substrate's vocabulary."

---

## C. Minor — clarifications, inconsistencies, and word-level fixes

### C1. §2.6's "no level modifies the level below it" is ambiguous about direction

> "No level modifies the level below it except through the Growth Gate, and Level 3 changes touch phys only via full gate evidence plus human approval."

Level 3 modifying Level 2, and Level 2 modifying Level 1, are strange directions. The intended rule is almost certainly: *no level modifies any other level except through the Growth Gate*. Reword.

### C2. §5 still does not bound linearized gravity

> "Linearized gravity is reserved for later, not the entry."

"Later" remains unbounded. Is it a GR-9, a second workload after the primary thread, a separate plan? If it is out of scope for this plan, say so explicitly. If it is in scope for a checkpoint, say which.

### C3. §5 names `sin²θ` as an encounter but does not tie it to a specific pass

> "and `sin²θ` trigonometric demand (see §9)."

Which pass first encounters it? GR-1 (metric) obviously, but the plan should say, so §20's "Trig ownership demonstrated" has a concrete check. Otherwise the DoD item is unfalsifiable.

### C4. §7's `Differentiate` row does not say what GR uses for *positive*-integer differentiation

The table lists `Differentiate` as "GR-3a probe only; expect fail-closed on negative-integer powers." But GR needs to differentiate positive-integer-power expressions too (e.g. a polynomial approximation in the Newtonian limit). Does GR use `phys.Differentiate` for those, or its own local differentiator for all cases? §9 implies the latter. State it.

### C5. §16's mutation test framing is confused

> "Flip Pin A fixture back to `Identify(Energy, E/c²)` → must fail dimension check (proves corrected fixture is non-vacuous)."

This mutates the *test fixture*, not the code under test. That is a test-sanity check, not a mutation test. The plan calls the whole section "Mutation Testing" but includes both:
- mutations of the test (this one), and
- mutations of production code (the others).

Split them: "test-sanity checks" vs "code mutation probes." Otherwise "mutation testing" will be read as a single discipline and both categories will be executed inconsistently.

### C6. §19's "no pre-building" is ambiguous about the new GR-local directories

> "No `semantic/`, `symbols/`, `constraints/`, `tensor/`, `manifold/`, `calculus/`, `hilbert/` pre-building before GR generates evidence."

But §8 defines `phys-gr/index/`, `phys-gr/metric/`, `phys-gr/curvature/`, etc. Are those forbidden pre-building too, or is the ban on Level-2/3 only? The list mixes plausible Level-2 names (`tensor/`, `manifold/`) with plausible Level-1 names (`semantic/`, `symbols/`). Clarify that the ban is on Level-2/3 pre-building, and that Level-1 GR directories may be created when the workload reaches the corresponding pass.

### C7. §3.3 and §16 disagree slightly on baseline wording

§3.3 says the tree is "frozen" (present tense). §16's PASS 0 gate says "baseline re-verified and frozen." The earlier review's B7 fix landed as "re-verified and frozen," which is fine — but §3.3's "is frozen" should be "is frozen *after* PASS 0 completes." The two sections describe the same state differently.

---

## D. What the revision gets right (so the criticism is calibrated)

- The three-level ladder is consistently applied across §§2, 13, 14, 18.
- The GR-3a/GR-3b split (§7) is exactly the right pattern: probe the boundary, classify it, then continue in userland, with an explicit prohibition on pretending the kernel now supports the capability.
- The C1–C8 ledger (§11) is self-contained, and the routing is frozen — this closes the earlier review's A2.
- The 20-field evidence record (§15) is enumerated — this closes the earlier review's B3.
- Pin A's corrected fixture (`Identify(RestMass, E/c²)`) is dimensionally coherent and non-vacuous — this closes the earlier review's C1.
- The failure taxonomy (§10) now explicitly states that UNREPRESENTABLE is not growth evidence until L1+L2 fail — this closes the earlier review's B5-adjacent concern.
- §3.3's "exact 42-file membership is unchanged" reconciles the earlier B6, subject to B3 above.
- The final decision principle (§20) is exactly the right framing and matches the architectural contract.

---

## E. Recommended order of fixes before PASS 0 is authorized

```text
1.  A1  Reconcile §13 and §18 (LEVEL-2 admissible only if §13 is relaxed,
        or remove it from this plan's outcomes).
2.  A2  Define the interop type between GR-local math and core.Expr.
3.  B1  Trim or motivate §7's unused allowed-operations.
4.  B2  Reword GR-3b's "does not bypass phys.Differentiate".
5.  B3  Add PASS 0 membership verification for its three file targets.
6.  B4  Define §15 field 18's "diff-shape" format.
7.  B5  Define §15 field 20's human-approval artifact.
8.  B6  Reword §18's three false readings.
9.  B7  Grandfather or extract existing in-phys Level-1 packages.
10. B8  Reclassify §12 probe 2 (user error vs substrate silent-wrongness).
11. C1-C7  Fold into the editorial pass (they do not block PASS 0).
```

A1 and A2 are blocking: PASS 0 cannot be authorized while they stand, because §9 and §18 cannot be executed as written. The B items should be resolved in the same revision that authorizes PASS 0. The C items are editorial and can land with §16's test-plan section when PASS 1 is drafted.

The plan is close to executable. A1 is the difference between "defensible plan" and "plan that silently contradicts its own outcomes list."
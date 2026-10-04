I reviewed **Plan 9** against the adversarial report, our locked decisions, and the earlier freeze architecture. It is close, but I would make **one substantive correction and one explicit decision** before implementation.

### 1. File-tree correction is sound

The decision to **preserve the three `reality_*.md` files but relocate them into `plan10/`** is compatible with the frozen implementation tree because `plan10/` is outside the module-count boundary. The same applies to moving the disposable `.gut` metadata and review helper out of the repository. plan9

The remaining issue is `.gitignore`: Plan 9 correctly observes that keeping it leaves 40 counted files, while removing it yields 39. plan9

Given the frozen requirement of exactly 39 files, the correct decision is:

> **Remove `.gitignore` from the counted repository tree.**

Do not weaken the specification just to retain it.

### 2. I would change B2 slightly

Plan 9 proposes:

```text
internal/kernel/types.go
    DimensionEnergySquared()
```

as an exported kernel helper. plan9

Functionally this fixes the public API violation, but architecturally I am slightly uncomfortable with it because we've just established:

```text
internal/kernel = theory-neutral machinery
relativity      = owns special-relativity physics
```

A helper literally named `DimensionEnergySquared` is **physics-domain-specific naming inside the generic kernel**.

The cleaner freeze interpretation is:

```text
public core:
    no DimensionEnergySquared()

relativity:
    owns whatever local construction is needed for E² dimension

internal/kernel:
    retains only generic dimension machinery
```

However, **do not redesign the dimension API** to accomplish this. Use the smallest existing generic dimension operation available. If the existing implementation already provides a generic way to square/multiply dimensions, use it locally in `relativity`. Only put a kernel helper there if no existing generic mechanism can express it without introducing broader API changes.

This is consistent with the purpose of the H2 boundary: the kernel should contain mandated generic machinery, while theory-specific meaning stays in the theory package.

### 3. `CoreObject()` decision remains correct

Keep the user-locked decision:

```go
CoreObject() core.Object
```

and leave free constructors unchanged.

Plan 9 correctly captures that scope. plan9

### 4. Everything else is appropriate

The Plan correctly preserves:

- E=mc² chain
- no `MassEnergyRelation()` executable reference
- no new physics
- no Einstein-1905 corpus
- no architecture redesign
- existing invariant suite
- fresh public API audit
- full M1–M18/N1–N7 rerun
- final `FREEZE-READY` / `NOT FREEZE-READY` discipline. plan9

The insistence on reviewing `attack_review_main.go` before moving it is also appropriate because that file was only partially reviewed in the adversarial pass. plan9

## Final decision

I would amend Plan 9 to:

```text
A. Preserve the three reality*.md files by moving them to plan10/
B. Move .gut out
C. Review and then move/remove attack_review_main.go as appropriate
D. Remove .gitignore so the implementation tree is exactly 39 files
E. Remove public core.DimensionEnergySquared()
F. Prefer a relativity-local implementation using existing generic Dimension
   machinery; do not add physics-specific naming to internal/kernel unless
   absolutely necessary
G. Change CoreObject() methods only
H. Rerun full verification + mutations
```

Everything else stays frozen.

So **Plan 9 is implementation-ready after that B2 refinement and the explicit `.gitignore = remove` decision**.
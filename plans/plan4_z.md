# The Einstein Test — Making Pen-and-Paper Derivation the Design Spec

**The concession lands where it should.** Enforcing MRC at the compiler level (rather than as a documentation principle, as v0.3 §23 had it) is exactly right, and the Einstein framing is the best thing that has happened to this spec — because E=mc² is not a metaphor here. It is a **decomposable worked example**, and decomposing it tells us precisely what the library must contain, what it must never pretend to do, and where the compiler's value actually lives.

## 1. What the 1905 derivation mechanically required

Strip the paper to its operations. Einstein did:

```text
1. POSTULATE   two inertial frames; energy conservation holds in both
2. DECLARE     "let L denote the energy emitted as light"   ← new symbol, mid-derivation
3. DERIVE      energy difference between frames: KE(S) − KE(S′) = L(γ − 1)
4. EXPAND      γ to second order in β = v/c          ← controlled, tagged truncation
5. IDENTIFY    L·v²/2c²  with  ½mv²                  ← THE Einstein moment
6. CONCLUDE    m = L/c², i.e. E = mc²
```

Now the crucial observation: **five of the six steps are mechanical. One is not.** Step 5 is the physics — the claim that the coefficient of v² is the inertial mass, i.e., that emitting energy reduces mass. Everything else is pen-and-paper algebra: subtracting equations, series expansion to a stated order, comparing coefficients, solving.

And that ratio is the design spec:

> **The compiler makes the five mechanical steps impeccable and the one insight step un-hideable. It never performs the insight.**

This is also exactly why the adversarial reviewer needs the ledger: a human reviewer of the 1905 paper doesn't re-check the algebra line by line so much as they *attack step 5* — "why is that coefficient the mass?" The compiler's entire value is surrounding that step with receipts: dimensions checked, structures compared, provenance sealed, nothing silent.

## 2. The derivation, in library form

```go
d := deriv.New("mass-energy-1905")

// ── Ground: postulates and declarations (provenance ASSUMED — never derivable)
c     := d.Postulate("c", physics.SpeedOfLight())
gamma := physics.LorentzFactor(c)                  // function of β := v/c
L     := d.Declare("L", physics.Energy())          // Einstein's introduced symbol
v     := d.Declare("v", physics.Velocity())

// ── Mechanical: frame-wise conservation (postulated) and the energy difference
bal := d.Postulate("conservation", relativity.EnergyConserved())
diff := d.Step("frame-difference",
    physics.Subtract(relativity.KE(bal, "S"), relativity.KE(bal, "S'")))
// → diff = Scale(L, gamma − 1)

// ── Mechanical: controlled expansion — truncation is TAGGED, never silent
approx := d.Step("expand",
    physics.ExpandSeries(diff, beta, 2))           // keeps order β², tags the truncation

// ── THE insight: explicit, dimension-checked, structure-checked, provenance-minted
massEq, err := physics.Identify(approx, physics.KineticEnergyForm(m, v))
// internally: Dim(L·v²/c²) = M ✓   quadratic-in-v energy ✓   nondegenerate match ✓
// the physical claim — "inertia depends on energy content" — is NOT verified;
// it is minted as an IDENTIFIED node, vet-flagged, ledger-recorded

// ── Mechanical: solve and seal
result := d.Conclude("m = L/c²", physics.Solve(massEq, physics.Mass()))
// Prov: DERIVED — premises {c, conservation, LorentzFactor, identification, order-2 truncation}
```

Every element the Einstein test forces, the corpus already designed and the library shape supports — the derivation proves the design was right all along:

- **`d.Declare`** — introducing new symbols mid-derivation (Einstein's "let L denote…"). Pen and paper is fluid; the ledger must allow cheap introduction and distinguish it from committed steps.
- **`physics.ExpandSeries(expr, β, 2)`** — controlled approximation whose truncation is **a tag on the returned value**, not a silently dropped tail. The order-2 truncation is a premise of the final result and appears in its provenance closure. Silent dropping of terms is one of the two classic ways derivations lie; this makes it impossible.
- **`physics.Identify(a, b)`** — the MRC-critical operation, and the reconciliation of the whole conversation: insight stays with the AI, but the *form* of the insight goes through a typed gate. Dimensions must match. Structures must be compatible. The identification cannot happen implicitly in algebra — `physics.Simplify` will never quietly decide two objects are "the same." Every `Identify` is vet-flagged and ledger-minted.

## 3. MRC at the compiler level — the enforcement table

Reinstated per your correction, mapped to the accepted review findings:

| Level | Mechanism | Enforced where | Finding it implements |
|---|---|---|---|
| 1. Construction | Typed constructors only; **carriers, never strings**; unexported fields; closed object model | Go compiler | F2, F6 |
| 2. Dimension | Rational-exponent dims checked in every operation | Go compiler + op entry | F7 |
| 3. Operation contract | Conditions as typed errors (`Contract` requires inverse metric; `Trace` requires finite/trace-class) | library, at every call | F4 |
| 4. **Identification gate** | `Identify` explicit; dimension + structure pre-check; mints `IDENTIFIED` provenance; always vet-flagged | library | the Einstein step |
| 5. Ledger closure | Only `deriv.Session` mints `DERIVED`; every derived node carries its premise closure; reviewer diffs routes by canonical form | `deriv` package | F1, F3 |

Note what this table resolves about the ban line: to support steps 3–5 the library *must* internally run series expansion, coefficient comparison, and algebraic solving — heavy mathematics. The ban holds because it bans a **public object ontology**, not internal machinery: the AI never constructs a `Series` or a `Polynomial`; it calls physics-named operations (`ExpandSeries`, `Compare`, `Solve`) whose signatures are typed in physics objects. The math is the ink, not the paper.

## 4. What stays non-mechanical — stated as a feature

The compiler verifies the **form** of step 5 (dimensions close, structures match); it cannot verify the **truth** of it (that inertial mass is what changes — that took experiments and decades of acceptance). This boundary is not a limitation to apologize for; it is the division of labor that makes the reviewer AI effective:

```text
mechanical steps  → reviewer recomputes them (deriv.Diff, route comparison, convention flips)
insight steps     → reviewer finds every Identify in the ledger and attacks its premises
final equation    → canonical hash; two routes agreeing byte-for-byte is a checkable event
```

An AI presenting a "derivation" in this system has exactly two honest options: route the insight through `Identify` (exposed, flagged, auditable) or fake a lemma (and the lemma registry is versioned and vetted — forged lemmas are visible). The chain-laundering attack from the review dies precisely because the *only* minter of derivation provenance is the session, and sessions record everything.

## 5. The one refinement the Einstein test adds that no prior round had

**Drafts.** Real pen-and-paper work is iterative — you write a line, cross it out, introduce a symbol, commit to a cleaner form later. The session should support this without polluting provenance:

```go
draft := d.Draft("try: expand in γ directly")   // working memory; no provenance weight
// ... abandoned ...
d.Commit("frame-difference", finalForm)          // the ledger takes only what's named
```

Cheap entry for exploration, expensive commitment for the ledger — the same discipline this project has converged on at every scale, now at the scale of one derivation.

## 6. v0.4 delta (short)

1. **MRC moves from §23–§26 (documentation principle) to the enforcement table above** — construction, dimension, contract, identification gate, ledger.
2. **Provenance and session restored** as the derivation substrate (F1): `Postulate / Declare / Define / Step / Lemma / Identify / Conclude / Draft / Commit`.
3. **Carriers-not-strings rule** in the object model (F2); `dim` restored to core (F7); canonical `Eq`/`Hash` + exact rational encoding (F3); conventions as value metadata (F5).
4. **Operation signatures table** — the real spec content — expanded to include the algebraic spine the Einstein test demands: `ExpandSeries`, `Compare`, `Substitute`, `Solve`, `KeepUpTo` (tagged truncation).
5. **New canonical stress test, ranked with the problem-of-time pair:** the 1905 derivation end-to-end, asserting (a) the ledger's shape, (b) that `Identify` flags, (c) that the truncation tag survives into the result's provenance, (d) that removing the conservation postulate makes `Conclude` fail.

**One sentence:** the compiler is the pencil and the notebook — it makes Einstein's five mechanical steps exact, tagged, and reproducible, forces his one moment of insight through a typed, flagged, recorded gate instead of letting it hide in algebra, and thereby gives the reviewer AI something it can actually attack; the intelligence does the physics, and the library makes the physics legible.
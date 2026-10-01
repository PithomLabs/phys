# IMPLEMENTATION PROMPT — AI-Agent Documentation Layer (v2.4 docs amendment)

You are the documentation implementation agent for the Physics Compiler
(`github.com/PithomLabs/phys`), currently frozen at **39 implementation files, FREEZE-READY**.

## 0. Locked contract (do not renegotiate)

```text
39 frozen implementation files
+ AGENTS.md
+ mechanics/README.md
+ relativity/README.md
= exactly 42 repository files
```

- `plan10/specs_v2_3.md` remains **byte-for-byte frozen**. Do not touch it.
- The v2.4 amendment is a separate document stored **outside the scanned
  repository tree** (see §9). It is NOT a 43rd repository file.
- `TestRepositoryTreeExact` is extended for **exactly the three paths above**
  and remains closed-world otherwise.
- No hypothesis/README, no core/ops/session/kernel READMEs, no root README
  modification, no implementation/corpus/manifest semantic changes.

## 1. READ BEFORE WRITING

Read the complete repository, not just selected files. At minimum:

```text
internal/kernel/
core/
ops/
session/
mechanics/
relativity/
hypothesis/
README.md
plan10/specs_v2_3.md
plan10/plan10_v2_3.md
plan10/plan7.md
plan10/plan8.md
plan10/plan9.md
plan10/plan9_imp.md
mechanics/manifest.json
relativity/manifest.json
```

Read the Go implementations behind the manifest constructors and the
derivation tests. Authority precedence:

```text
specs_v2_3.md → actual implementation → manifest.json → tests → AGENTS.md → package README.md
```

Do not invent behavior or capabilities. Where the implementation is
deliberately bounded (bounded `Solve`, no CAS/prover, special-relativity
only), document the boundary.

## 2. AGENTS.md (repo root, new)

Global operating manual for AI agents. Required content:

1. **What this system is**: agent = theorist/formulator; library = formal
   symbolic substrate + corpus; humans/empirical reality = scientific
   authority. The library never adjudicates physical truth.
2. **Three epistemic layers**: ESTABLISHED framework vs DERIVED artifact vs
   HYPOTHESIS. `ESTABLISHED ≠ universally true`; formally derivable ≠
   empirically established; HYPOTHESIS is never auto-promoted.
3. **Canonical reasoning workflow checklist** (§5 of draft: 22 steps from
   physical question → framework → README → manifest → assumptions/regime →
   premises → derivation with supported ops only → provenance/contamination
   checks → session recording → hypothesis only if needed).
4. **Assumption-first reasoning**: deterministic union, exact
   `(Kind, Key)` conflict detection, structured `ExprValue` assumptions,
   operation-generated preconditions (e.g. §11.5 denominator, sign
   assumptions). Never silently introduce/discard assumptions; never prose
   over structured values.
5. **Math discipline**: construction vs transformation vs identification vs
   hypothesis formation; `Simplify`/`Compare`/`Solve` (one exact quadratic
   pattern only)/`SelectBranch`/`Identify` (session-owned) semantics.
6. **Derivation checklist** (20 items, draft §8): premises, framework,
   assumptions, conventions, regime, validity, dimensions, kinds, supported
   op, canonical params, provenance, contamination, branches, preconditions,
   no corpus-shortcut, canonical form, session use, interpretation split.
7. **Hypothesis workflow**: observation → framework → premises →
   assumptions → transformation → gap/anomaly → candidate → derivation →
   predictions → falsification conditions → recovery/limit claims → research
   candidate → human evaluation. `Prediction` / `FalsificationCondition` /
   `RecoveryClaim` / `AnomalyReference` are recorded, never decided.
8. **Manifest-first discovery**: per-item fields (`id`, `constructor`,
   `kind`, `statement`, `canonical_expr`, `dimension`, `provenance_status`,
   `source`, `assumptions`, `derivable_from`, `reduces_to`, `known_limits`,
   `anomalies`, `falsification_conditions`). `statement` = orientation only;
   `canonical_expr` = machine-checkable object. Never parse prose as math.
9. **Package boundaries**: kernel = theory-neutral machinery;
   `mechanics/` = classical only; `relativity/` = special relativity only
   (Minkowski, Lorentz symmetry, special-relativistic regime, no
   gravitational dynamics, metric `-+++`); `hypothesis/` = provisional
   space. No implicit assumption merging across frameworks.
10. **E=mc² rule**: stored `MassEnergyRelation` is a corpus artifact, never
    the derivation source. The demonstration chain is
    `EnergyMomentumRelation → ZeroThreeMomentum → Substitute → Simplify →
    Solve → Compare(gte) → SelectBranch → m*c²` — formal derivation from
    encoded premises, not the 1905 historical reconstruction, not empirical
    discovery. Hypothesis/candidate workflow (when to form one + checklist).
11. **Anti-patterns** (draft §15, 12 items).
12. **Theorem-vs-primitive rule + kernel-growth boundary** (draft §§21–22):
    theorems/results belong to theory/user packages; kernel grows only for
    generic, foundational, theory-neutral, cross-framework capabilities.
    Document the discipline; implement nothing.

## 3. Package READMEs (new, ~1 page each, orientation not manifest dump)

Template: Framework Identity (ID + corpus status) / Scope / Not in Scope /
Assumptions / Conventions / Core Primitives / Core Relations / Derivation
Targets / Limits and Known Boundaries / AI Reasoning Guidance / Canonical
Metadata (`manifest.json`).

- `relativity/README.md`: `special_relativity`, ESTABLISHED, SR-only scope,
  Minkowski/Lorentz/regime/no-gravity assumptions, `-+++` convention,
  Spacetime/MinkowskiMetric/RestMass/Energy/ThreeMomentum/FourMomentum/
  SpeedOfLight/Velocity primitives, LorentzFactor/EnergyMomentumRelation/
  MassEnergyRelation relations, E=mc² target, no-gravity limit. Verify every
  value against manifest/code.
- `mechanics/README.md`: framework ID/status/scope from manifest,
  Mass/Time/Position/Velocity/Acceleration/Force/Momentum/Energy/
  KineticEnergy primitives, F=ma / p=mv / K=½mv² relations, limits/anomalies
  from manifest. Verify against manifest/code.

## 4. Forbidden in this pass

Modify production `.go` (except §6 allowlist edit), tests (except §6),
manifests, kernel/core/ops/session semantics, public APIs, corpus content,
package boundaries, root README, specs_v2_3.md bytes. No GR, no 1905
reconstruction, no new theorems/physics/primitives. No fourth (or fifth)
documentation file anywhere in the scanned tree.

## 5. v2.4 amendment document (outside the scanned tree)

Write `plan10/specs_v2_4_amendment.md` (plan10/ is excluded from
`TestRepositoryTreeExact`; verify exclusion holds). Contents only: §3 tree
gains exactly the 3 enumerated paths (39 + 3 = 42 exact); §38 guardrail set
to that exact configuration; all else v2.3 unchanged. **Do not place this
file (or any other new file) inside a scanned directory.**

## 6. The one permitted test edit

Extend `TestRepositoryTreeExact`'s allowlist with exactly:

```text
AGENTS.md
mechanics/README.md
relativity/README.md
```

No generic documentation exception. Nothing else in the test changes.

## 7. Verification gate (all mandatory)

```text
go build ./... → go vet ./... → go test ./... -count=1 (all green)
filesystem tree recount → exactly 42/42 against §3 + amendment
public API audit (9 dimension ctors; no new exports)
doc↔code↔manifest cross-checks (every README claim traced)
```

## 8. Final report

```text
## Documentation Implementation Report
### Files created: AGENTS.md, mechanics/README.md, relativity/README.md
### Amendment: plan10/specs_v2_4_amendment.md (outside tree)
### Test change: 3-path allowlist extension (exact diff)
### Code/corpus/kernel changes: NONE
### Verification: build/vet/tests, 42/42 tree, cross-checks
### Final status: DOCUMENTATION-READY (re-freeze by independent reviewer only)
```

Goal: an AI agent reading `AGENTS.md` + package README + manifest can reason
about the physics without reverse-engineering the Go implementation. The
implementer does NOT declare the re-freeze; the independent adversarial
reviewer does.

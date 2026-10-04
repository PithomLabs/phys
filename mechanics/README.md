# Classical Mechanics

## Framework Identity
Framework ID: `classical_mechanics`
Corpus Status: `ESTABLISHED` (human-curated manifest metadata, never computed)

## Scope
Point-particle Newtonian mechanics in the nonrelativistic regime
(`v << c`): inertial mass, absolute time, position/velocity/acceleration
kinematics, force, momentum, and (kinetic) energy.

## Not in Scope
Relativistic, quantum, gravitational, or continuum physics. Anything at
velocities comparable to `c` is outside this package (see anomaly below).

## Assumptions
The manifest records no framework-level assumption entries; the regime is
carried by the limit below. Per-object content is fixed by the
constructors (e.g. `NewtonSecondLaw` carries no extra assumptions).

## Conventions
None recorded for this package (contrast `relativity/`'s metric signature).

## Core Primitives
Fixed constructors only; all `DEFINED` / `ESTABLISHED` unless noted:

```text
NewMass / NewTime / NewPosition / NewVelocity /
NewAcceleration / NewForce / NewMomentum / NewEnergy
NewKineticEnergy(m, v)   parameterized K = ½·m·v²; DERIVED / NONE (not a manifest item)
```

## Core Relations

```text
NewtonSecondLaw       F = m·a   (source: Newton, Principia)
MomentumRelation      p = m·v
KineticEnergyRelation K = ½·m·v²
```

`NewtonSecondLaw` is the `F = ma` acceptance target: its expression must
match the manifest `canonical_expr` byte-for-byte.

## Derivation Targets
- `Differentiate(KineticEnergy, Velocity)` → canonical `m·v`
  (`Kind=Expression`, Momentum dimension) via product + power rules —
  exercises the bounded differentiation engine, not a stored answer.
- Session-ledged derivations over these premises with full replay.

## Limits and Known Boundaries
- Limit `nonrelativistic`: `v << c`, classical limit of special relativity.
- Anomaly `galilean_noninvariance` (`known_limitation`, on
  `NewtonSecondLaw` and `KineticEnergyRelation`): fails at velocities
  comparable to `c`; motivates relativity.

## AI Reasoning Guidance
Start from the fixed constructors; check dimensions (M/L/T compositions)
and kinds (`KineticEnergy` ≢ `Energy` for add/compare) at every step;
merge assumptions explicitly; never paste `F = m·a` prose as a premise —
use the constructor or derive it. Crossing into `v ~ c` means leaving this
framework: switch packages, do not stretch this one.

## Canonical Metadata
`mechanics/manifest.json` (11 items, 9 domains) — machine-readable source
of truth. This README is orientation only.

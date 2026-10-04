# Special Relativity

## Framework Identity
Framework ID: `special_relativity`
Corpus Status: `ESTABLISHED` (human-curated manifest metadata, never computed)

## Scope
**Special relativity only**: flat-spacetime kinematics and dynamics —
spacetime, Minkowski metric, rest mass, energy, three/four-momentum,
invariant speed of light, velocity, Lorentz factor, and the
energy–momentum / mass–energy relations.

## Not in Scope
General relativity, gravitational dynamics, curved spacetime, black-hole
physics, quantum mechanics, quantum gravity. See `no_gravity` below.

## Assumptions
Framework entries from the manifest (structured, not prose):

```text
constraint  rest_mass_nonnegative    m >= 0   (carried by RestMass)
constraint  speed_of_light_positive  c > 0    (carried by SpeedOfLight)
domain      minkowski_spacetime      Minkowski spacetime
physical    lorentz_symmetry         Lorentz symmetry
domain      no_gravitational_dynamics
domain      special_relativistic_regime
```

`EnergyMomentumRelation` carries both sign assumptions; they transitively
enable the E=mc² simplification through bounded sign entailment.

## Conventions
`metric.signature = -+++`, recorded on `NewMinkowskiMetric()` and in the
manifest. Convention conflicts (same key, different value) fail operations.

## Core Primitives
Fixed constructors; `DEFINED` / `ESTABLISHED`:

```text
NewSpacetime / NewMinkowskiMetric / NewRestMass / NewEnergy /
NewThreeMomentum / NewFourMomentum / NewSpeedOfLight / NewVelocity
```

`NewVelocity` is a constructor helper, not a manifest item. Rest-frame and
zero helpers (non-manifest): `ZeroThreeMomentum` (carries the `rest_frame`
assumption `p = 0`), `ZeroEnergy`, `ZeroVelocity`, `RestFrameAssumption()`.

## Core Relations

```text
LorentzFactor            γ = 1/sqrt(1−v²/c²) as Call(lorentz_factor, v)
EnergyMomentumRelation   E² = (p·c)² + (m·c²)²  (source: Einstein, 1905)
MassEnergyRelation       E = m·c²  DERIVED corpus item,
                         derivable_from EnergyMomentumRelation + RestFrame
```

The E² dimension is `Energy×Energy` composed locally in this package via
generic dimension multiplication — no special kernel machinery exists.

## Derivation Targets
Primary: **E = mc²**, derived never copied —
`EnergyMomentumRelation → ZeroThreeMomentum → Substitute → Simplify →
Solve(Energy) → Compare(Energy, ZeroEnergy, gte) → SelectBranch → m·c²`,
each golden state asserted in `derivation_test.go` (which scans itself for
any use of the stored relation as a premise). Secondary:
`Limit(LorentzFactor, Velocity, ZeroVelocity) → 1` via fixed-body Call
expansion, never an ID shortcut.

## Limits and Known Boundaries
- `flat_spacetime`: no gravitational fields; spacetime is Minkowski.
- `classical_limit`: `v << c` recovers classical mechanics.
- Anomaly `no_gravity` (`scope_limit`, on `EnergyMomentumRelation`):
  cannot describe event horizons or black holes; motivates GR as future
  work, not present content.

## AI Reasoning Guidance
Read the 6 framework assumptions before transforming anything; carry
`m >= 0` / `c > 0` or sign-sensitive rewrites (`Sqrt(Pow(x,2)) → x`) fail
closed. `Solve` accepts exactly `Relation(eq, t², rhs)`; `SelectBranch`
needs the `gte`-against-zero constraint. This README orients; the proof
obligations live in the manifest + derivation tests.

## Canonical Metadata
`relativity/manifest.json` (10 items, 8 domains) — machine-readable source
of truth. This README is orientation only.

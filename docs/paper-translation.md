# Paper Translation — Physics Compiler MVP (Plan 10 v2.3 / specs_v2_3.md)

This document is auxiliary documentation. It is not canonical machine-readable physics metadata, is not parsed at runtime, and introduces no new package, abstraction, dependency, or architecture change.

---

## Common notation

- `core.Object` — the single immutable generic carrier (§5).
- `core.Expr` — closed expression node set (§6; 9 node kinds).
- `core.Kind` — physical kind enum (18 values, stable ordinals §6).
- `core.Dimension` — dimension vector (§7; `Divide`, `Pow` methods).
- `ops.OperationID` — 12 closed operations (`add`, `subtract`, `multiply`, `divide`, `pow`, `simplify`, `substitute`, `differentiate`, `limit`, `compare`, `solve`, `select_branch`). No plugin registry.
- `session.Session` — the only authority that commits a ledger, produces `IDENTIFIED`, or seals a `ResearchCandidate` (§16 / MRC-006 / MRC-007).

## Framework mapping

| Plan 10 term | Source reference | Normative section |
|---|---|---|
| `core.Object` | `specs_v2_3.md` §5 | Object carrier / Mint authority |
| `core.Kind` | `specs_v2_3.md` §6 | 18-value enum |
| `core.Expr` | `specs_v2_3.md` §7 / §9 | Expression / rewrite rules |
| `internal/kernel.ObjectSpec` | `specs_v2_3.md` §5.2.1 | Mint contract (7-point validation) |
| `ops.Apply` | `specs_v2_3.md` §15 / §15.13.1 | Positional operation dispatch |
| `session.Session` | `specs_v2_3.md` §16 | State machine (New → Drafting → Committed → Concluded → Sealed) |
| `ResearchCandidate` | `specs_v2_3.md` §26 | Sealed artifact; `UnverifiedResearchCandidate.Validate()` delegates to same private constructor |
| `MRC-001` | `specs_v2_3.md` §5 / Plan 10 §1 | Constructor authority (public/package boundary, not cryptographic/hostile-source protection) |
| `MRC-006` | `specs_v2_3.md` §16.4 | Identify firewall (`Session.Identify` only) |
| `MRC-008` | `specs_v2_3.md` §26.8 / §26.9 | Candidate containment (`CandidateContainmentError`) |

## Ambiguity resolution

1. **No truth decision:** The library defines no truth API (`TRUTH_SCORE`, `TRUST`, `PROMOTE`, etc.). `ResearchCandidate.Validate()` asserts internal consistency only (§26.8); it never claims physical truth (§26.10).
2. **No cryptographic protection claim:** `internal/kernel.Object` fields are unexported; the `internal/` visibility rule is a Go-language package boundary (REQ-004-03), not adversarial-source protection (REQ-000-03 / §0.3). The README states this explicitly.
3. **Replay ≠ authenticity:** The SHA-256 chain detects mutation only. Anyone rewriting the full ledger can recompute hashes; signing/authenticity is deferred (§16.10 / REQ-002-15).
4. **Position contract (§15.13.1):** Inputs are never reordered. Replay in `Validate` uses exactly the same indices as `Apply`.
5. **Assumption mode (§11 / F2):** `EntailsNonNegative` and `EntailsPositive` match only expression-valued assumptions (`expr` mode). Text-valued assumptions do not trigger bounded sign entailment. This is pinned by the mechanics/relativity manifest construction and the derivation golden sequence (`E² = (pc)² + (mc²)²` → `m*c²`).

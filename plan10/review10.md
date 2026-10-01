This adversarial review is doing exactly what we wanted: it found **three concrete documentation/freeze issues without finding a production-code regression**. The correct status is therefore **not ready to re-freeze**. adv_review10

### What the review established

**F-001 is the most important finding.**  
`TestRepositoryTreeExact` is too permissive: its allowlist contains five extra paths, so the test would fail to enforce the intended exact 42-file boundary. The reviewer identified the precise five entries and the exact correction: reduce the allowlist to the 39 frozen files plus the three authorized documentation files. adv_review10 adv_review10

**F-002 catches a real conceptual problem in `AGENTS.md`.**  
The sentence saying every object belongs to exactly one epistemic layer is misleading because **provenance status and corpus status are independent axes**. The review correctly uses `MassEnergyRelation` as the counterexample: it is `DERIVED` by provenance while `ESTABLISHED` by corpus status. adv_review10

The proposed correction is exactly the right direction:

> Objects carry two independent status axes: provenance and corpus status.

That is actually important for the AI-agent layer because conflating those two can cause an agent to mistake “derived” for “unestablished,” or “established corpus item” for “primitive/fundamental.”

**F-003 is a metadata precision issue.**  
The mechanics README says the anomaly applies to “all three relations,” while the manifest associates it with only two relations. That needs to be made manifest-exact. adv_review10

### One correction to my previous review

The adversarial reviewer also resolved something I had flagged earlier: **`galilean_noninvariance` is in fact the current mechanics anomaly ID**. So my earlier suspicion that the canonical ID was `nonrelativistic_regime` was wrong; the independent review directly checked the current manifest and confirms `galilean_noninvariance`. adv_review10

That is exactly why the independent pass was valuable.

### What is already clean

The reviewer independently verified:

- exact 42 non-plan files currently present
- v2.3 spec unchanged
- amendment document outside the repository tree
- mechanics and relativity manifest metadata
- `MassEnergyRelation`'s `DERIVED` + `ESTABLISHED` combination
- package boundaries
- closed 12-operation surface
- bounded `Limit`, `Differentiate`, and `Solve`
- identification authority
- no promotion API
- E=mc² anti-hardcoding
- Lorentz-limit behavior
- full build, vet, and test suite passing. adv_review10 adv_review10

So this is a **clean documentation/conformance remediation**, not a kernel or physics-architecture reopening.

### Recommended next step

Do **only F-001 through F-003**, then rerun the independent review rather than accepting the implementation agent's own remediation report as sufficient. The current reviewer explicitly withholds re-freeze authorization until those defects are fixed. adv_review10

The important thing is: **do not weaken the tree test, do not modify the manifests, and do not reinterpret the epistemic model to make the README fit.** Fix the documentation and allowlist to match the already-frozen semantics.
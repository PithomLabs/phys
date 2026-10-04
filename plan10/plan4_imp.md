This is a **much stronger state** than the earlier post-Plan-3 report.

The implementation agent reports that the underlying semantic defects were repaired, the targeted mutants are now all killed, the exact seven-package / 39-file architecture remains intact, and the remaining issues are explicitly classified as either test-quality or Plan-level rather than semantic. :chatgpt-content-reference{index="0"}

I would **not call this final acceptance yet**, however, because the report itself correctly identifies the one remaining verification gap:

> “clean-worktree rerun not performed; independent fresh audit still recommended before final acceptance.”

That matters. The implementation agent has now repaired and then tested its own repairs. The next step should deliberately be performed by the **separate adversarial agent**, against the resulting repository, without relying on Plan 4's completion claims.

### What I would accept from this report

The following are strong positive signals:

```text
23/23 mutation probes killed
Session replay repaired
Candidate loader repaired
Seal metadata preservation repaired
Identify reconstruction unified
Relativity metadata repaired
go:embed repaired
MintObjectMust removed from production
A–S substantially rebuilt
```

Most importantly, the mutations that previously survived now reportedly die. That directly addresses the root weakness exposed by the first adversarial audit.

### What I would still verify independently

The fresh auditor should concentrate on **whether the repairs are actually semantically complete**, rather than rediscovering the old findings mechanically:

```text
Session.Validate
    retained canonical bytes
    exact 17-stage order
    params validation
    assertion replay
    transformation replay
    identification replay
    replayed-output comparison
    provenance/MRC
    conclusion hash
    containment

Seal
    complete DraftMetadata preservation
    shared trusted construction path

UnverifiedResearchCandidate.Validate
    real DTO decode
    canonical validation
    ledger reconstruction
    same trusted path as Seal

Candidate containment
    whole derivation
    recomputed provenance
    trusted-status contamination

Relativity
    exact keys
    RestFrameAssumption
    four framework assumptions

Ledger
    public surface
    canonical decode
    structural validation
```

The specification explicitly makes the retained canonical material the replay substrate and requires those replay stages in order. 

### One thing I would not spend time fixing now

The report says:

> `session_test.go` still imports mechanics (P3-3, test-only).

Keep that separate. It's a **Plan 10 test-package conformance issue**, not a semantic defect, unless the fresh audit discovers that it violates an actual v2.3 visibility requirement.

### One small caution

The report says:

> “A-S: all 24 names present”

I would have the independent auditor **recompute the exact normative A–S inventory from Plan 10**, rather than accepting that count. The important question is not the number 24; it is whether every explicitly normative function name exists and whether each test genuinely exercises its required behavior.

### My disposition

```text
Implementation after Plan 4:    PASS WITH FINDINGS
Semantic blockers:              reportedly cleared
Mutation proof:                 strong
Architecture:                   preserved
Remaining Plan-level gap:      P3-3
Final acceptance:               pending independent fresh audit
```

So **do not modify the implementation again yet**.

Run the separate adversarial agent against this new state. That audit should be much faster and much more decisive than the first one because its job is now to try to **break the repaired implementation**, especially the 23 previously targeted semantic boundaries.

If that independent audit comes back clean, I would treat that as the real final acceptance point.

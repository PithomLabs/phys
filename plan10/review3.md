The third adversarial audit is **the right final finding set**: the implementation itself now passes the substantive audits, but the proof is still incomplete in two places. The report is explicit that M9 and M15 are the only surviving mutants, with everything else passing, including replay, candidate handling, relativity, canonicalization, A–S, REQ-032, and clean-copy verification. :chatgpt-content-reference{index="0"}

I would **not change production semantics** for either one yet. These are proof gaps.

### M9 — valid and important

The specification requires the internal object replay decoder to support canonical object round-trip:

```text
Object
→ canonical JSON
→ Object
```

with byte-identical canonical form/hash, and explicitly requires representative object round-trip tests. 

The audit found that `LoadObjectJSON` can accept a non-canonical encoding because the test suite never proves rejection of such input. :chatgpt-content-reference{index="2"}

The right test is therefore **not merely another successful round-trip test**. Add a negative test that:

```text
1. creates a valid canonical object JSON;
2. changes only canonical representation
   (for example field ordering / permitted whitespace / another
   actually non-canonical representation supported by the parser);
3. calls internal LoadObjectJSON;
4. asserts rejection;
5. verifies errors.As gives the required typed error;
6. verifies canonical input still succeeds.
```

The exact mutation must preserve semantic content while violating canonical-byte form.

That proves the decoder is not merely a permissive JSON parser.

### M15 — valid and even more clearly a test-design gap

The spec requires **exactly one production `MintObject` entry point** and explicitly prohibits a second minting helper with a different semantic contract. :chatgpt-content-reference{index="3"}

The audit showed the current `TestNoGenericFactory` only recognizes exact identifiers such as `MintObject` and `NewObject`, so a mutant named `MintObjectAlt` survives. :chatgpt-content-reference{index="4"}

The correct fix is **not** to forbid every function whose name happens to contain “Mint.” Instead, make the test enforce the actual architectural invariant.

A good test should mechanically inspect production AST and establish:

```text
Allowed object-producing authority:
    internal/kernel.MintObject

Allowed non-mint object decoder:
    internal/kernel.LoadObjectJSON

No other production function in the kernel may constitute
a second Object-minting entry point.
```

For example, the source test can identify production functions that return `Object` / `(Object,error)` and then distinguish the explicitly allowed decoder from the mint authority. Better still, combine that with an AST/source inspection of calls and construction behavior so a renamed `MintObjectAlt` cannot evade the test merely by changing its identifier.

The important property is **semantic detection**, not name matching.

### One subtle issue in the audit

The report says M9 first fails as “TestExactRationalRoundTrip.” That test name is misleading for the mutation. The actual Plan 10 requirement is that P includes both exact rational/expression round-trip and **representative object** canonical round-trip; Plan 10 separately names `TestObjectCanonicalRoundTrip` as coverage for the object codec. :chatgpt-content-reference{index="5"}

So I would add a dedicated:

```text
TestObjectCanonicalRoundTripRejectsNonCanonical
```

or similar, rather than making `TestExactRationalRoundTrip` responsible for object-decoder rejection.

And for M15:

```text
TestNoGenericFactory
```

can remain the normative A–S/REQ-032 test, but it should be strengthened so it detects the architectural invariant rather than exact identifier strings.

## What I would tell the implementation agent

> **Address M9 and M15 as TEST-PROOF corrections only unless direct inspection reveals a production defect. Do not alter the architecture or add files.**
>
> **M9:** add a substantive internal object-codec negative test proving `LoadObjectJSON` rejects a semantically equivalent but non-canonical object encoding, while canonical encoding succeeds and round-trips byte-identically. Keep `TestExactRationalRoundTrip` focused on its intended P coverage and use a dedicated object-codec test where appropriate.
>
> **M15:** strengthen the existing `TestNoGenericFactory`/related source audit so a renamed second minting entry point such as `MintObjectAlt` is detected. Do not rely on exact identifier matching. Detect the architectural invariant: exactly one production object-minting authority (`internal/kernel.MintObject`), while explicitly allowing the internal object decoder and other non-mint object-returning functions that the specification permits. Use AST/source semantics, not a name blacklist.
>
> Before changes, reproduce M9 and M15 individually in a disposable copy and record the surviving behavior. After changes, rerun exactly M1–M18 with `-count=1`, requiring every mutant to fail for the intended reason. Restore the pristine repository and run `go build ./...`, `go vet ./...`, and `go test ./...` plus a fresh clean-copy check.
>
> Do not update documentation claims until the 18/18 result is genuinely reproduced.

The target remains **18/18 meaningful kills**, not merely a higher percentage. The independent auditor has already verified that M12, M13, and M18 now have proper stage-isolated proof, while M9 and M15 remain the only survivors. :chatgpt-content-reference{index="6"}

At this point, **do not touch the core implementation unless M9 reproduction proves the decoder itself is wrong**. The current evidence points to the tests failing to prove two already-specified invariants.

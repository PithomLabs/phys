# Physics Compiler — Specification Amendment v2.4 (documentation only)

**Amends:** `plan10/specs_v2_3.md` (which remains byte-for-byte frozen).
**Scope:** repository tree + guardrail only. No semantic change.

## §3 — Repository tree (amended)

The frozen **implementation tree** is unchanged: exactly the 39 files listed
in specs_v2_3.md §3 (24 Go source, 10 Go test, 2 JSON manifests,
`go.mod` + `README.md` + `docs/paper-translation.md`).

Additionally, the repository contains exactly three explicitly authorized
AI-agent documentation files:

```text
AGENTS.md
mechanics/README.md
relativity/README.md
```

No other documentation, helper, scaffold, or generated file is permitted in
scanned directories. The frozen repository tree is therefore exactly:

```text
39 frozen implementation files + 3 authorized documentation files = 42 files
```

## §38 — Repository file-count guardrail (amended)

The guardrail is the exact approved configuration above: 42 files maximum,
consisting of the 39 frozen implementation files plus the 3 explicitly
enumerated documentation files. This is not a loose allowance; any further
addition requires its own amendment.

## Non-reopened (explicit)

This amendment does not touch: `internal/kernel`, public APIs, corpus
semantics, operation semantics, provenance rules, session/replay rules,
manifests, physics content, or package boundaries. It authorizes no
general documentation exception — the three paths are enumerated.

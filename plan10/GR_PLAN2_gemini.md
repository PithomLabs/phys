```python
with open('GR_IMPLEMENTATION_PLAN.md', 'r') as f:
    plan_text = f.read()

print("Length of GR_IMPLEMENTATION_PLAN.md:", len(plan_text))


```

```text
Length of GR_IMPLEMENTATION_PLAN.md: 24289


```

```python
import os, glob

go_files = glob.glob('*.go')
print("Go files found:", go_files)

for gf in go_files:
    print(f"--- {gf} ---")
    with open(gf, 'r') as f:
        content = f.read()
        print(content[:500])
        print("...\n")


```

```text
Go files found: ['relations.go', 'primitives.go', 'manifest_test.go', 'derivation_test.go']
--- relations.go ---
package relativity

import (
	"math/big"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

// LorentzFactor returns the Lorentz factor γ = 1/sqrt(1 - v^2/c^2) as a Call.
func LorentzFactor() kernel.Object {
	v, _ := kernel.NewSymbol("v")
	call, _ := kernel.NewCall(kernel.LorentzFactorFunctionID, v)
	cPosExpr := kernel.NewRelation(kernel.RelationGt, kernel.NewSymbolMust("c"), kernel.NewRational(big.NewRat(0, 1)))
	cPos := kernel.NewExprAssumptionMust(kernel.Assu
...

--- primitives.go ---
package relativity

import (
	"math/big"

	_ "embed"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

//go:embed manifest.json
var ManifestJSON []byte

// Spacetime is the canonical spacetime object (symbol s).
type Spacetime struct {
	Object kernel.Object
}

func NewSpacetime() Spacetime {
	sym, _ := kernel.NewSymbol("s")
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Special Relativity corpus", "special_relativity",
		nil, [32]byte{}, [32]byte{}, ke
...

--- manifest_test.go ---
package relativity

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

func TestRelativityManifestCrossCheck(t *testing.T) {
	m, err := core.ParseManifest(ManifestJSON)
	if err != nil {
		t.Fatalf("ParseManifest: %v", err)
	}
	if m.FrameworkID != "special_relativity" {
		t.Errorf("framework_id = %q", m.FrameworkID)
	}
	if m.CorpusStatus != cor
...

--- derivation_test.go ---
package relativity

import (
	"go/ast"
	"go/parser"
	"go/token"
	"math/big"
	"os"
	"runtime"
	"testing"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
	"github.com/PithomLabs/phys/ops"
)

// TestFullDerivationSequence tests the complete §20 derivation:
// EnergyMomentumRelation -> ZeroThreeMomentum -> Substitute -> Simplify -> Solve -> Compare -> SelectBranch
// Final result: m*c^2

// mustMint is a test-only convenience; production code uses MintObject directl
...



```

```python
with open('specs_v2_3.md', 'r') as f:
    specs_text = f.read()

# Let's search for key items in specs_v2_3.md
keywords = ['NewCandidateConcept', 'Differentiate', 'Identify', 'KindRelation', 'KindExpression', 'SelectBranch', 'Solve', 'Substitute', 'Simplify', 'HYPOTHESIS']
for kw in keywords:
    count = specs_text.count(kw)
    print(f"Keyword '{kw}': {count} occurrences")


```

```text
Keyword 'NewCandidateConcept': 1 occurrences
Keyword 'Differentiate': 8 occurrences
Keyword 'Identify': 25 occurrences
Keyword 'KindRelation': 0 occurrences
Keyword 'KindExpression': 0 occurrences
Keyword 'SelectBranch': 6 occurrences
Keyword 'Solve': 7 occurrences
Keyword 'Substitute': 8 occurrences
Keyword 'Simplify': 11 occurrences
Keyword 'HYPOTHESIS': 26 occurrences


```

```python
import re

# Find section titles or tables in specs_v2_3.md related to Kind, hypothesis, ops, session
for line in specs_text.split('\n'):
    if line.startswith('#') or 'Kind' in line:
        if any(h in line for h in ['#', 'Kind', 'kind', 'Enum', 'enum']):
            print(line[:120])


```

```text
# Physics Compiler MVP — Normative Implementation Specification v2.3
## 0. Authority, module identity, and governing boundary
### 0.1 Exact Go module identity
### 0.2 Product boundary
### 0.3 Threat model for constructor authority
# 1. Product requirements
# 2. Explicit non-goals
# 3. Exact MVP repository layout
# 4. Package dependency graph and authority boundary
internal/kernel  ← concrete immutable Object/Expr/Kind/Dimension authority
The concrete backing definitions for `Object`, `Expr`, `Kind`, and `Dimension` live in `internal/kernel`.
# 5. Core physical-object abstraction
## 5.1 Two-layer representation
**REQ-005-05**: No domain package may expose a public constructor that permits callers to choose `Kind`, `Dimension`, `P
## 5.2.1 Exact internal object mint API
    Kind          Kind
`Object`, `Expr`, `Kind`, `Dimension`, and the metadata value types used by `ObjectSpec` are concrete kernel types. `cor
## 5.2 Internal implementation of `core.Object`
Kind
## 5.3 Exact read-only accessors
Kind() Kind
## 5.4 Trusted construction
## 5.5 Invalid-object rule
# 6. Physical kind
`Dimension` and `Kind` are intentionally independent.
The MVP `Kind` enumeration is exactly:
**REQ-006-01**: The implementation MUST include every Kind listed above.
**REQ-006-02**: `Kind` MUST be constructor-assigned and immutable.
## 6.1 Why kind exists
MRC-003 therefore checks `Kind` separately from `Dimension`.
## 6.2 Kind compatibility table
same named Kind + same Dimension → allowed
named Kind + different named Kind → rejected
named Kind + Expression           → rejected
named Kind + same named Kind       → allowed when dimensions match
named Kind + Expression             → allowed when dimensions match
Expression + named Kind             → allowed when dimensions match
# 7. Dimensions
## 7.0 Exact public Dimension API
## 7.1 Base dimensions
## 7.2 Required dimensions
## 7.3 Dimension operations
# 8. Symbolic expression representation
## 8.0 Exact public expression construction and inspection API
Kind() ExprKind
A method MAY return an empty/invalid value when the method does not apply to the node kind; operation code MUST check `E
## 8.1 Closed expression tree
## 8.2 Node ordinals
## 8.2.1 Exact expression/operator enum surface
`ExprKind` MUST expose these stable values:
## 8.3 Symbol
## 8.4 Rational
## 8.5 Add and Mul
## 8.6 Neg
## 8.7 Pow
## 8.8 Sqrt
## 8.9 Call
# 9. Canonicalization and mechanical rewrite rules
## 9.1 Canonical equality
## 9.2 Canonical hashing
## 9.3 Child sorting
## 9.4 Rational combination
## 9.5 Canonical scalar/sign normal form
## 9.6 Identity, zero, and exact rational-power rules
## 9.7 Repeated powers
## 9.8 Sqrt rewrite rules and bounded sign entailment
## 9.9 Relation normalization
# 10. Canonical JSON schema
## 10.1 Canonical JSON requirements
## 10.2 Expression JSON
## 10.3 Dimension JSON
## 10.4 Object JSON
## 10.5 Round-trip requirement
# 11. Assumption system
## 11.0 Exact assumption metadata API
func NewTextAssumption(kind AssumptionKind, key, value string) (Assumption, error)
func NewExprAssumption(kind AssumptionKind, key string, value Expr) (Assumption, error)
`Assumption` MUST expose read-only accessors for `Kind`, `Key`, and typed Value mode/content.
## 11.1 Assumption kinds
## 11.2 Structured assumption value
Kind
## 11.3 Construction
## 11.4 Merge
Two assumptions with the same `(Kind, Key)` and different canonical values MUST yield:
## 11.5 Denominator preconditions
Kind = MathPrecondition
## 11.6 Sign assumptions used in MVP
# 12. Conventions
## 12.0 Exact convention metadata API
Conventions are separate immutable metadata from assumptions, even though `AssumptionKind.Convention` remains a reserved
# 13. Provenance and corpus status
## 13.0 Exact provenance API
## 13.0.1 Exact provenance construction helpers
## 13.1 Provenance statuses
## 13.2 Deterministic provenance propagation
## 13.3 Provenance record
## 13.4 Corpus status
## 13.5 Source and framework inheritance
# 14. MRC version and rules
## 14.1 MRC-001 — Constructor/carrier integrity
## 14.2 MRC-002 — Dimensional compatibility
## 14.3 MRC-003 — Physical-kind compatibility
## 14.4 MRC-004 — Assumption compatibility
## 14.5 MRC-005 — Convention compatibility
## 14.6 MRC-006 — Explicit physical identification
## 14.7 MRC-007 — Session/provenance authority
## 14.8 MRC-008 — Candidate containment
# 15. Operation API and common operation semantics
## 15.1 Add
same Kind as operands
## 15.2 Subtract
## 15.3 Multiply
## 15.4 Divide
## 15.5 Pow
## 15.6 Simplify
## 15.7 Substitute
## 15.8 Differentiate
Every successful `Differentiate` call returns a `core.Object` with `Kind=Expression` and the exact derivative dimension 
## 15.9 Limit
### 15.9.1 Fixed Lorentz-factor function body
Every successful `Limit` call returns a `core.Object` with `Kind=Expression` and the exact resulting limit dimension.
## 15.10 Compare
Kind = Relation
## 15.11 Solve
## 15.12 SelectBranch
## 15.13 Pure-operation dispatch
Unused fields are encoded with their empty/default values. The only permitted `OperationParams.Kind` values are:
    Kind          string
### 15.13.1 Normative positional argument mapping
# 16. Session and derivation ledger
## 16.1 Exact nine session actions
## 16.2 Session lifecycle
## 16.2.1 Exact candidate-draft metadata
## 16.3 Exact session method signatures
## 16.4 Draft
## 16.5 Assertion actions
## 16.6 Step execution
## 16.7 Identify ownership
## 16.8 Draft buffer
## 16.9 Commit
## 16.10 Conclude
## 16.11 Seal
## 16.12 Step structure
StepKind
`StepKind` values:
## 16.13 Input retention
## 16.14 Operation parameters
## 16.15 Genesis hash
## 16.16 StepID
## 16.17 StepEnvelope hash
## 16.18 Derivation hash
## 16.19 Replay validation
## 16.20 Tamper behavior
## 16.21 Exact Ledger value surface
# 17. Corpus manifests
## 17.0 Exact public manifest API
    Kind                    string
## 17.1 Manifest top-level schema
## 17.2 Domain entries
## 17.3 Limit entries
## 17.4 Anomaly entries
## 17.5 Manifest item schema
## 17.6 Statement field
## 17.7 Canonical manifest expression validation
## 17.8 Constructor cross-check
6. compare manifest `kind` with `object.Kind()`;
## 17.9 Corpus status ownership
# 18. Mechanics corpus
## 18.0 Exact mechanics constructor surface
## 18.1 Populated primitives
## 18.2 Mechanics relations
with `Kind = KineticEnergy` and `Dimension = Energy`.
## 18.3 Mechanics framework metadata
# 19. Special-relativity corpus
## 19.0 Exact relativity constructor surface
## 19.1 Populated primitives
## 19.2 RestFrame
Kind: Constraint
## 19.3 Fixed special-relativity assumptions
## 19.4 Lorentz factor
Kind = Expression
## 19.5 Energy-momentum relation
Kind = Relation
## 19.6 Mass-energy relation
## 19.7 Zero constructors
    Kind = ThreeMomentum
    Kind = Energy
    Kind = Velocity
## 19.8 Symbol isolation
## 19.9 Relativity limitation/anomaly
# 20. Canonical E = mc² vertical derivation
## 20.1 Required sequence
## 20.2 Golden trace
### Step 1 — input relation
### Step 2 — substitute rest-frame zero momentum
### Step 3 — mechanical simplify
### Step 4 — solve exact accepted quadratic form
### Step 5 — explicit positive-energy constraint
### Step 6 — select positive branch and simplify
### Final selected solution
## 20.3 Why the sign step is valid in the MVP
# 21. Calculus micro-test
# 22. Relativistic limit micro-test
# 23. Identify firewall micro-test
# 24. Hypothesis candidate concept
## 24.1 Candidate concept definition
## 24.2 Candidate concept constructor
    kind core.Kind,
**REQ-024-02**: Candidate concept construction MUST require explicit `Kind` and `Dimension`.
## 24.3 Candidate containment test
## 24.4 No promotion API
# 25. Falsifiability structures
## 25.1 Prediction
## 25.2 FalsificationCondition
## 25.3 RecoveryClaim
## 25.4 AnomalyReference
## 25.5 Required falsifiability test
# 26. ResearchCandidate artifact
## 26.0 Exact ResearchCandidate support types
## 26.1 Authority
## 26.2 Required fields
## 26.3 Hypothesis field
## 26.4 Premises
## 26.5 Framework dependencies
## 26.6 Derivation
## 26.7 Review history
## 26.8 Validation pipeline
## 26.9 Candidate containment at artifact level
## 26.10 ResearchCandidate public surface
# 27. Review artifact
## 27.1 Challenge
## 27.2 Review
## 27.3 Categories
## 27.4 Orchestration boundary
# 28. Manifest data requirements
## 28.1 Mechanics manifest required items
## 28.2 Relativity manifest required items
## 28.3 Limitation/anomaly requirements
# 29. MRC fallibility and revision boundary
# 30. Mandatory error taxonomy
# 31. Determinism requirements
## 31.1 Determinism test
# 32. Mandatory negative tests
# 33. Explicit acceptance tests A–S
## A — Typed mechanics
## B — Classical relation
## C — Dimension rejection
## D — Category rejection
## E — Assumption conflict
## F — Convention conflict
## G — Differentiation
## H — Relativistic relation
## I — Mass-energy derivation
## J — Relativistic limit
## K — Identify firewall
## L — Candidate containment
## M — Falsifiability
## N — Anomaly
## O — Manifest
## P — Exact rational round-trip
## Q — Ledger tamper/replay
## R — Determinism
## S — Candidate handoff
# 34. Framework evaluation anatomy
## 34.1 Deterministic derivation test
## 34.2 Executable self-consistency/reduction
## 34.3 Limitation/anomaly record
# 35. Paper translation documentation
# 36. `physvet` boundary
# 37. Explicit v0.5+ deferrals
# 38. Scope guardrail
# 39. Normative requirement and test coverage matrix
# 40. Final architectural invariants
# 41. Final handoff condition — no architectural blockers
8. **Pow result:** `Kind=Expression`.


```

```python
# Let's extract section 24.2 and other specific signatures from specs_v2_3.md
def get_section(title, text):
    pos = text.find(title)
    if pos == -1:
        return "NOT FOUND"
    return text[pos:pos+1500]

print("--- Section 24.2 ---")
print(get_section("## 24.2 Candidate concept constructor", specs_text))

print("--- Section 15.8 Differentiate ---")
print(get_section("## 15.8 Differentiate", specs_text))

print("--- Section 16.3 Exact session method signatures ---")
print(get_section("## 16.3 Exact session method signatures", specs_text))

print("--- Section 16.7 Identify ownership ---")
print(get_section("## 16.7 Identify ownership", specs_text))


```

```text
--- Section 24.2 ---
## 24.2 Candidate concept constructor

The exact constructor is conceptually:

```go
func NewCandidateConcept(
    id string,
    kind core.Kind,
    dimension core.Dimension,
    expr core.Expr,
    assumptions core.AssumptionSet,
    conventions core.ConventionSet,
) (core.Object, error)

```

Requirements:

* `id` non-empty;
* `kind` explicit and immutable after construction;
* `dimension` explicit and immutable after construction;
* `expr` valid;
* assumptions/conventions deterministic;
* provenance forced to `HYPOTHESIS`;
* corpus status forced to `NONE`;
* no caller-supplied provenance status;
* no promotion API.

**REQ-024-01**: Candidate concept construction MUST force `HYPOTHESIS` provenance.

**REQ-024-02**: Candidate concept construction MUST require explicit `Kind` and `Dimension`.

## 24.3 Candidate containment test

The test MUST:

1. create a candidate concept;
2. derive with it;
3. verify downstream result is `HYPOTHESIS`;
4. verify `Simplify`, `Solve`, `Limit`, `Compare`, `SelectBranch`, and `Session.Identify` do not upgrade it;
5. verify no promotion API exists;
6. verify a hypothesis-dependent ResearchCandidate cannot validate as trusted corpus authority.

The last assertion is a validation/integrity assertion, not a scientific-truth judgment.

## 24.4 No promotion API

The MVP source tree MUST NOT expose exported functions/methods named or semantically equivalent to:

```text
Promote
Trust
ApproveHypothesis
PromoteToEstablished
SetCorpusStatus

```

The cand
--- Section 15.8 Differentiate ---

## 15.8 Differentiate

Exact signature:

```go
func Differentiate(target, wrt core.Object) (core.Object, error)

```

`wrt.Expr()` MUST be a single `Symbol` node.

The bounded MVP derivative engine supports only:

```text
Rational        → 0
Symbol          → 1 if same symbol, otherwise 0
Add             → derivative of each term
Mul             → n-ary product rule
Neg             → Neg(derivative)
Pow             → integer-constant power rule

```

For `Pow(base,n)` where `n` is a non-negative integer rational:

```text
d(base^n)/dx = n * base^(n-1) * d(base)/dx

```

The result MUST be simplified after differentiation.

Unsupported:

```text
Sqrt
Call
Relation
BranchSet
non-integer symbolic exponent

```

These return `UnsupportedOperationError`.

Every successful `Differentiate` call returns a `core.Object` with `Kind=Expression` and the exact derivative dimension computed from the target dimension divided by the `wrt` dimension.

The MVP does not create a `Derivative` node.

## 15.9 Limit

Exact signature:

```go
func Limit(target, variable, value core.Object) (core.Object, error)

```

The MVP performs direct substitution followed by simplification.

For `Call(function="lorentz_factor")`, it MUST expand the fixed definition in §15.9.1 first.

No L'Hôpital algorithm exists.

No series expansion exists.

No general asymptotic engine exists.

Unsupported or singular results return `UnsupportedOperationError`.

### 15.9.1 Fixed Lorentz-factor function body

The only MVP function
--- Section 16.3 Exact session method signatures ---

## 16.3 Exact session method signatures

The session MUST expose methods equivalent to:

```go
func New() *Session

func (s *Session) Draft(derivationID, label string, metadata DraftMetadata) error

func (s *Session) Postulate(label string, object core.Object) error

func (s *Session) Declare(label string, object core.Object) error

func (s *Session) Define(label string, object core.Object) error

func (s *Session) Step(label string, operation ops.OperationID, inputs []core.Object, params ops.OperationParams) (core.Object, error)

func (s *Session) Identify(a, b core.Object, justification string) (core.Object, error)

func (s *Session) Conclude(object core.Object) error

func (s *Session) Commit() error

func (s *Session) Seal() (ResearchCandidate, error)

func (s *Session) Validate() error

func (s *Session) CanonicalJSON() ([]byte, error)

```

Exact pointer/value receiver details MAY follow ordinary Go conventions, but the semantic signatures above are fixed.

## 16.4 Draft

`Draft` records:

```text
DerivationID
Label
MRCVersion
GenesisHash

```

It clears no existing state because a session starts empty.

`DerivationID` MUST be non-empty and caller-supplied. It is therefore deterministic and not auto-generated.

## 16.5 Assertion actions

`Postulate(label,obj)` requires:

```text
obj.Valid() == true
obj.Provenance().Status == POSTULATED

```

MVP does not require a public domain constructor that emits `POSTULATED`; the action is exercised with an `internal/kernel` test fixt
--- Section 16.7 Identify ownership ---

## 16.7 Identify ownership

There is **no** package-level:

```text
ops.Identify

```

`Session.Identify` is the only identification API.

It MUST:

1. validate session state;
2. validate both objects;
3. validate dimensions and kinds under MRC-006;
4. require non-empty justification after trimming whitespace;
5. construct a Relation(eq, a.Expr(), b.Expr());
6. set provenance to `IDENTIFIED` unless hypothesis contamination forces `HYPOTHESIS`;
7. preserve the merged assumptions/conventions;
8. record the justification in provenance;
9. append an identification step to the draft buffer;
10. return the immutable relation object.

It MUST NOT consult or modify any global active-session state.

## 16.8 Draft buffer

`Step` and `Identify` append fully specified step material to the session's draft buffer.

Each draft entry contains enough material to later form a committed hash-chain record.

No draft entry is externally mutable through the public API.

## 16.9 Commit

`Commit` freezes the draft buffer into the committed ledger in recording order.

Commit MUST:

1. assign deterministic indices beginning at one;
2. assign `StepID` as:

```text
step-000001
step-000002
...

```

3. set the genesis/previous hash chain;
4. compute each current step hash;
5. preserve canonical inputs/outputs/params;
6. make the committed ledger immutable through the public API.

After commit, the draft buffer MUST be empty.

## 16.10 Conclude

`Conclude(object)` is legal only after `Commit`.

The supplied

```



```python?code_reference&code_event_index=6
import glob

all_files = glob.glob('**/*', recursive=True)
for f in sorted(all_files):
    if f.endswith('.go') or f.endswith('.json') or f.endswith('.md'):
        print(f)


```

```text
AGENTS.md
GR_IMPLEMENTATION_PLAN.md
GR_IMPLEMENTATION_PLAN_20261002.md
GR_PLAN_review.md
README.md
adv_review5.md
derivation_test.go
manifest.json
manifest_test.go
physics_compiler_mvp_implementation_plan_prompt_v2_2.md
physics_compiler_mvp_specs_v2_2.md
plan10_v2_3.md
plan11.1.md
plan11.1_overview.md
plan7.md
plan9.1_claude.md
plan9.2.md
plan9_claude.md
primitives.go
prompt2.md
relations.go
specs_falsifiabilty.md
specs_v2_3.md


```

An adversarial line-by-line review of `GR_IMPLEMENTATION_PLAN_20261002.md` confirms that the primary core architecture (frozen `phys` kernel substrate, external `phys-gr` Level-1 module, HYPOTHESIS scalar bridge, derivation-first Schwarzschild, and non-ratio Growth Gate) is **sound and fully intact** . The locked decisions (D1–D18) eliminate major architectural loopholes .

However, reading the execution contract line-by-line reveals **7 specific implementation ambiguities and missing concrete signatures** that could block or cause non-deterministic choices for an automated coding agent.

---

## Architectural Soundness Verification

The fundamental boundary choices in v2 are completely solid:

1. **Kernel Freeze & Zero Leakage:** `phys` remains 100% frozen . No GR concepts enter `internal/kernel` .
2. **Bridge Isolation:** The bridge between `phys-gr` and `phys` is restricted to `GRExpr → core.Expr → hypothesis.NewCandidateConcept` for kernel-compatible subexpressions only . Level-1 constructs (`Sin`, `Cos`, `UnknownFunction`) never pass to `phys` .
3. **No Kernel Differentiation Contamination:** GR-3a probing `phys.Differentiate` fail-closed behavior on negative integer exponents before GR-3b engages `phys-gr` local differentiation cleanly distinguishes substrate testing from theory execution .
4. **Non-Ratio Promotion:** Decoupling kernel growth from effective code volume or operation counts prevents synthetic test inflation .

---

## Detailed Line-by-Line Findings & Concrete Resolutions

To ensure **zero ambiguity** for a coding agent, the following 7 items must be pinned directly in the execution contract:

### 1. PASS0 Pin B Exact File Target Is Ambiguous

* **Plan Location:** §16 (`PASS0 — Conformance Pins`) & §3.3 .
* **Issue:** §16 states Pin B target is "`ops` tests" . In a strict 42-file frozen repository where file creation is forbidden, specifying "`ops` tests" leaves the agent guessing between potential test files (e.g., `ops/ops_test.go` vs `ops/differentiate_test.go`) .
* **Resolution:** Pin B target file path must be explicitly pinned as `ops/differentiate_test.go` (or the exact existing file name in the 42-file inventory).

### 2. Scalar Bridge Is Missing the Inverse Unwrapping Contract (`FromCoreExpr`)

* **Plan Location:** §6 (`Representation Strategy`) & §7 (`Kernel Interaction Points`) .
* **Issue:** §6 specifies `ToCoreExpr()` and `hypothesis.NewCandidateConcept(...)` to bridge `GRExpr → core.Object` for kernel operations like `Simplify` . §7 mentions taking the result of `phys.Simplify` and "reintegrating into `GRExpr`" . However, no inverse parsing contract `FromCoreExpr(core.Expr) (GRExpr, error)` is defined.
* **Resolution:** Define `phys-gr/symbolic/bridge.go` API explicitly:
```go
func ToCoreExpr(e GRExpr) (core.Expr, error)
func FromCoreExpr(e core.Expr) (GRExpr, error)
func ToCoreObject(e GRExpr, dim core.Dimension) (*core.Object, error)

```



### 3. Dimension Passing Contract for Bridge Wrapping

* **Plan Location:** §6 (`Representation Strategy`) .
* **Issue:** §6 notes `hypothesis.NewCandidateConcept` requires an explicit `Dimension` . Metric components ($g_{tt}, g_{rr}$) are dimensionless, but coordinate components ($r^2 d\theta^2$) carry dimensions ($L^2$). `GRExpr` does not natively store unit dimensions for arbitrary symbolic expressions .
* **Resolution:** `ToCoreObject` must accept an explicit `core.Dimension` parameter from the caller (e.g. `ToCoreObject(expr, core.Dimensionless)`), or fall back to deriving dimension if all constituent symbols have registered dimensions in `phys-gr/coordinates`.

### 4. Replay Trace Hash Chaining Mechanics (`GRStep`)

* **Plan Location:** §7 (`Kernel Interaction Points`) & Appendix C .
* **Issue:** §7 defines `GRStep` fields: `StepID / OperationID / InputCanonical / ParamsCanonical / OutputCanonical / CurrentHash` . It does not state whether `CurrentHash` is calculated over the step contents alone or chained with `PreviousHash`.
* **Resolution:** Explicitly define `CurrentHash` as a Merkle-chained hash to match the kernel session invariant:

$$\text{CurrentHash} = \text{SHA256}(\text{PreviousHash} \mathbin{\Vert} \text{CanonicalJSON}(\text{Step}_{n}))$$



For `StepID = 0`, $\text{PreviousHash} = \text{"0000000000000000000000000000000000000000000000000000000000000000"}$.

### 5. Tensor Engine Index Storage Standard

* **Plan Location:** §8 (`Level-1 GR Package Design`) .
* **Issue:** §8 outlines `Tensor` with `Rank`, `IndexSlots[]`, `Components`, `Symmetries` . If `Components` storage format is left open, a coding agent might choose sparse component maps with symmetry keys, leading to complex index permutation bugs during tensor contractions .
* **Resolution:** Require full dense tuple indexing for `Components`:
```go
type ComponentKey [4]int // (x0, r, theta, phi) = (0, 1, 2, 3)
type Tensor struct {
    Rank       int
    Slots      []IndexSlot
    Components map[ComponentKey]GRExpr
    ChartID    string
}

```


Symmetries are enforced on write/construction by validating $T_{i j} = T_{j i}$ or $R_{i j k l} = -R_{j i k l}$.

### 6. Precise Bounded Solver Steps for GR-8 Derivation

* **Plan Location:** §5 (`GR Workload`) & §9 (`GR-Local Mathematics Required`) .
* **Issue:** §5 mandates deriving Schwarzschild from unknown $A(r), B(r)$ via "bounded differential/algebraic reduction (problem-specific bounded local solver; NO general ODE solver)" . Without a rigid sequence, an agent might attempt to implement generic term rewrite engines that fail to terminate.
* **Resolution:** Hardcode the exact 4-step canonical transformation path in `phys-gr/solution/schwarzschild.go`:
1. **Linear Combination:** Compute $B(r) R_{tt} + A(r) R_{rr} = 0 \implies A'(r)B(r) + A(r)B'(r) = 0 \implies (A(r)B(r))' = 0$.
2. **Boundary Substitution:** Apply $A(\infty)B(\infty) = 1 \implies A(r)B(r) = 1 \implies B(r) = \frac{1}{A(r)}$.
3. **Angular Reduction:** Substitute $B(r) = \frac{1}{A(r)}$ into $R_{\theta\theta} = 0 \implies 1 - A(r) - r A'(r) = 0 \implies (r A(r))' = 1$.
4. **Integration & Mass Identification:** Integrate $(r A(r))' = 1 \implies r A(r) = r + C_1 \implies A(r) = 1 + \frac{C_1}{r}$. Match weak field limit $C_1 = -2GM/c^2$.



### 7. Explicit Local Differentiator Derivative Rules

* **Plan Location:** §9 (`GR-Local Mathematics Required`) .
* **Issue:** §9 specifies a "bounded differentiator" in `phys-gr/symbolic/differentiate.go` .
* **Resolution:** Pin the complete table of supported operations in `differentiate.go`:
* **Sum/Difference:** $\frac{d}{dx}(u \pm v) = u' \pm v'$
* **Product Rule:** $\frac{d}{dx}(u \cdot v) = u' v + u v'$
* **Power Rule:** $\frac{d}{dx}(u^n) = n u^{n-1} u'$ (where $n \in \mathbb{Q}$)
* **Trig Derivatives:** $\frac{d}{dx}\sin(u) = \cos(u) u'$, $\frac{d}{dx}\cos(u) = -\sin(u) u'$
* **Unknown Function:** $\frac{d}{dr} A(r) = A'(r)$, $\frac{d}{dr} A'(r) = A''(r)$
* **Constants / Other Variables:** $\frac{d}{dx}(c) = 0$, $\frac{d}{dx}(y) = 0$ ($y \neq x$)



---

## Action Plan for Coding Agent

| Task | Target Package | Action |
| --- | --- | --- |
| **PASS0 Gate** | `phys` repository | Implement Pin A (`session/session_test.go`), Pin B (`ops/differentiate_test.go`), and Pin C (`AGENTS.md`) . Verify `go test ./...` passes without altering frozen file list . |
| **Module Setup** | `phys-gr` repository | Initialize sibling module `[github.com/PithomLabs/phys-gr](https://github.com/PithomLabs/phys-gr)` with `go.mod` requiring `phys` . |
| **Symbolic Algebra** | `phys-gr/symbolic` | Implement `GRExpr` AST, `differentiate.go` (fixed 7-rule set), `normalize.go` ($\sin^2 + \cos^2 \to 1$), and `bridge.go` (`ToCoreExpr`, `FromCoreExpr`, `ToCoreObject`) . |
| **Tensor Engine** | `phys-gr/metric`, `index` | Implement dense 4-tuple key mapping, Christoffel assembly, Riemann, Ricci, and Einstein tensor calculators . |
| **Derivation & Limit** | `phys-gr/solution`, `limit` | Implement the 4-step canonical ODE solver for Schwarzschild and SI dimensional audit for weak-field Newtonian limit . |
| **Adversarial Probes** | `phys-gr/adversarial` | Run silent wrongness probes (Probe A, B, C) and log classification outputs . |

With these 7 clarifying resolutions incorporated, the plan has **zero architectural blockers** and provides an unequivocal specification for an automated coding agent .
Yes. Use the following as the **single planning prompt for the Physics Compiler**. It is deliberately written to force the coding agent to consolidate the four source plans rather than average them together or reintroduce the older research-methodology architecture.

The strongest foundation is the Z plan's separation of conceptual scaffolding from compiler-native structure and its Go-style interfaces/tooling, combined with Qwen's conventional compiler pipeline, DeepSeek's broader math/physics grammar inventory, and Gemini's tensor/index/dimension concerns.    

```text id="w8o0cm"
You are the senior compiler architect responsible for drafting the
implementation plan for the next major project:

                    THE PHYSICS COMPILER

DO NOT IMPLEMENT CODE YET.

Your task is to inspect the repository and the attached/reference design
documents, reconcile them against the locked design below, and produce ONE
complete implementation plan.

This is a language/compiler project.

It is NOT an extension of the old research-methodology system.

Do not produce a piecemeal plan that leaves fundamental language, type,
interface, composition, AST, or IR decisions for later.

The goal is to leave us with a coherent v0 compiler blueprint that can be
implemented directly.

============================================================
1. AUTHORITATIVE DESIGN DIRECTION
============================================================

The Physics Compiler is a programming-language implementation for
expressing mathematical and physical structures.

Its central purpose is:

    REPRESENT HUMAN INTENT
    +
    ENFORCE THE FORMAL LANGUAGE
    +
    PROVIDE ANALYSIS TOOLS
    +
    EMIT A PRECISE REPRESENTATION

The compiler does NOT:

    adjudicate physics
    litigate papers
    decide physical truth
    interpret experimental results
    select theories
    promote claims
    maintain evidence
    manage research methodology
    act as an adversarial reviewer
    act as EBP
    act as Solvent
    decide what a physical theory "really means"

The compiler does not interpret the physical significance of the program.

It parses what the human wrote, represents it formally, checks the formal
semantics of the language, and exposes the resulting representation and
diagnostics for humans and downstream tools to analyze.

Use the principle:

    REPRESENT, DON'T INTERPRET.

More precisely:

    The compiler enforces the formal semantics of its language.
    It does not infer or adjudicate the truth or physical meaning of
    the resulting theory.

============================================================
2. MOST IMPORTANT SEPARATION
============================================================

Separate:

    HUMAN CONCEPTUAL SCAFFOLDING

from:

    COMPILER-NATIVE STRUCTURE

Human conceptual scaffolding includes, unless there is a concrete
language-level reason to encode a piece mechanically:

    Dancer / Dance / Stage
    Reality-First methodology
    No Sacred Cows
    philosophical "jurisdiction" narratives
    research-methodology workflows
    paper-review methodology
    epistemic promotion concepts
    adversarial review methodology
    broad research maps

These may appear in:

    documentation
    examples
    doc comments
    analysis output
    annotations

but MUST NOT become the core ontology of the compiler.

Compiler-native concerns are:

    grammar
    tokens
    syntax
    AST
    declarations
    expressions
    names/scopes
    types
    type parameters
    interfaces/capabilities
    composition
    dimensions
    tensor/index structure
    domains/codomains
    mathematical operators
    physical operators
    relations
    dependencies
    imports/packages
    lowering
    IR
    export data
    diagnostics
    standard library APIs

The compiler must care about what is mechanically represented and checked,
not about preserving a philosophical taxonomy.

============================================================
3. USE GO'S COMPILER DESIGN AS THE PRIMARY ENGINEERING MODEL
============================================================

Follow the broad architectural discipline of the Go toolchain.

The v0 pipeline should be conceptually:

    source
      ↓
    scanner / lexer
      ↓
    parser
      ↓
    AST
      ↓
    name resolution / scope
      ↓
    type checking
      ↓
    interface / capability checking
      ↓
    composition checking
      ↓
    semantic analysis
      ↓
    typed Physics IR
      ↓
    export / analysis representation
      ↓
    downstream tooling / Program B

Do not invent exotic compiler architecture.

Use standard compiler engineering practices:

    deterministic compilation
    source positions/spans
    clear diagnostics
    separate frontend and IR
    reusable type-checking library
    package/module boundaries
    conformance tests
    standard library
    tooling around compiler output

Qwen's scanner → parser → AST → types → typecheck → semantic analysis →
IR organization is the preferred structural starting point.

============================================================
4. IMPORTANT CORRECTION TO THE GO ANALOGY
============================================================

Do NOT literally equate:

    physics derivation graph = Go package graph

and do not force every physics derivation into a separate Go-style package.

Use Go package/import architecture as inspiration for SOFTWARE MODULARITY.

The Physics Compiler must still have its own:

    physics dependency graph
    derivation dependencies
    physical relationships

inside the language/IR.

Similarly:

    software dependency cycle
        ≠
    physical mutual dependence

The compiler should preserve the distinction.

Use Go's structural ideas, not superficial one-to-one metaphors.

============================================================
5. THE LANGUAGE ITSELF
============================================================

We are no longer merely building the existing Program A SDK.

We are designing an actual source language, provisionally:

    .phys

The implementation plan must define a coherent v0 language.

Do NOT just describe a set of Go structs.

Define the language at the level of:

    lexical grammar
    declarations
    expressions
    operations
    relations
    blocks
    package/import syntax
    types
    interfaces
    generic parameters
    composition
    assumptions/postulates
    derivation statements
    representation statements
    identification statements
    constraints
    couplings
    regimes
    dimensions
    tensor/index expressions where supported
    annotations/doc comments

The plan must include concrete syntax examples.

============================================================
6. MATHEMATICS MUST HAVE A REAL GRAMMAR
============================================================

Represent mathematics as an actual language rather than a list of names.

The compiler should be able to represent mathematical structures such as:

    Scalar
    Vector
    Covector
    Tensor
    Function
    Functional
    Set
    Relation
    Map
    Operator
    VectorSpace
    InnerProductSpace
    HilbertSpace
    Group
    Algebra
    Manifold
    Metric
    Connection
    Curvature
    Topology
    Measure
    ProbabilityMeasure
    Spectrum
    Invariant
    Limit
    Variation
    Expectation
    Correlation
    Distribution
    etc.

This list is a candidate inventory, NOT a giant closed compiler enum.

The implementation plan must instead determine:

    primitive types
    parameterized types
    constructors
    interfaces
    operators
    relations
    type classes / capability-style contracts where appropriate
    standard library packages

DeepSeek's broader math inventory may be used as a source for standard
library candidates, but do not turn the entire list into hardcoded
compiler ontology. :contentReference[oaicite:4]{index=4}

============================================================
7. PHYSICS MUST ALSO HAVE A REAL GRAMMAR
============================================================

Represent physical structures explicitly.

Potential v0 standard-library concepts include:

    State
    System
    Subsystem
    Field
    Particle
    Observable
    Measurement
    Configuration
    Action
    Lagrangian
    Hamiltonian
    Constraint
    Generator
    Symmetry
    Conservation
    Charge
    Current
    Energy
    Momentum
    Mass
    Spin
    GaugeStructure
    Spacetime
    Geometry
    Coupling
    InitialCondition
    BoundaryCondition
    Regime
    Scale
    Vacuum
    GroundState
    ExcitedState
    etc.

Again:

    these are library concepts,
    not necessarily compiler keywords,
    not necessarily hardcoded ontology,
    and not all must be implemented in v0.

Determine the smallest useful core vocabulary and a principled extension
mechanism.

============================================================
8. INTERFACES ARE CENTRAL
============================================================

This is one of the most important requirements.

Use Go-style structural interfaces as a model for expressing CAPABILITY
GRAMMAR.

A concept should be usable in an operation because it satisfies the
required formal contract, not because a philosophy table says that the
concept "belongs" to that operation.

Illustrative examples:

    Spatial
    Temporal
    Observable
    Dynamical
    Stateful
    Measurable
    Differentiable
    Integrable
    Composable
    TensorLike
    etc.

For example, conceptually:

    type Temporal interface {
        Duration() Interval
    }

means that a value is usable where Temporal is required if it satisfies
the language's Duration contract.

This does NOT mean:

    "Nature has certified this as physical time."

It means:

    "This program has supplied the formal interface required by the
     language's Temporal type."

The implementation plan must determine:

    how interfaces are declared
    how interfaces are satisfied
    whether satisfaction is structural
    how method/operator signatures are represented
    whether interfaces may embed other interfaces
    whether generic interfaces are supported
    how interface satisfaction is checked
    how interfaces interact with physical and mathematical types
    how users extend capabilities without changing compiler code

Avoid a giant `PhysRole` or `MathKind` taxonomy as the primary semantic
mechanism.

============================================================
9. COMPOSITION IS THE CENTER OF THE TYPE SYSTEM
============================================================

This must be stronger than merely "constructors exist."

The language should treat composition as a first-class operation.

At minimum distinguish:

A. Mathematical structural composition

    VectorSpace + InnerProduct
        → InnerProductSpace

    InnerProductSpace + Completeness
        → HilbertSpace

    Manifold + Metric
        → Metric structure

    Metric + Connection
        → Curvature structure

B. Physical composition

    Field + Spacetime
        → FieldOnSpacetime

    Matter + Geometry
        → CoupledSystem

C. Capability/interface composition

    HilbertSpace
        satisfies InnerProductSpace
        satisfies CompleteSpace

D. Parametric composition

    Field[T]
    State[S]
    Operator[A,B]
    Tensor[Rank, Space]

Use Go-style generics where appropriate, but do not mechanically copy Go's
generic syntax into `.phys` unless the language design supports it cleanly.

The implementation plan must define the composition model BEFORE defining a
large standard-library inventory.

Do not manually enumerate every complex concept.

Use:

    primitives
    constructors
    interfaces
    parameters
    properties
    composition rules

to build complex structures.

============================================================
10. TYPES MUST BE MORE PRECISE THAN NAMES
============================================================

Avoid semantic systems where:

    "Group"

    "HilbertSpace"

    "Time"

    "Field"

are merely string labels.

Types must carry enough structure to support checking.

The plan must address:

    nominal vs structural typing
    type identity
    parameterized types
    domain/codomain
    dimensions
    index structure
    field domains
    scalar/value types
    physical qualifiers
    regime qualifiers
    interface satisfaction
    aliases
    type construction

Do not decide everything through a single giant enum.

The type system is the core of the compiler.

============================================================
11. MATHEMATICAL DOMAIN / CODOMAIN
============================================================

Operators and maps must have explicit signatures.

Examples:

    Operator[A,B]
    Map[A,B]
    Function[A,B]

A compiler error should be possible because:

    operation requires A → B
    supplied object is C

not because:

    "C is philosophically not allowed here."

The plan should explain how signatures are represented and checked.

============================================================
12. TENSOR / INDEX / DIMENSION SUPPORT
============================================================

Gemini's tensor/index machinery should be incorporated selectively.

The compiler should eventually support formal structures such as:

    tensor rank
    covariant / contravariant indices
    index binding
    contraction
    metric-mediated raising/lowering
    dimension vectors
    dimensional consistency
    derivative-induced type/rank changes

Examples that the language should be capable of representing:

    A^μ B_μ

versus:

    A_μ B_μ

when the corresponding index rules make the second expression invalid.

Dimensional analysis should support symbolic dimension vectors such as:

    M L T^-2

and algebraic composition of dimensions.

Do not build a unit-conversion empire in v0.

Separate:

    dimension
    unit
    numerical value

where appropriate.

Gemini's tensor/index/dimension ideas are valuable inputs here. :contentReference[oaicite:5]{index=5}

============================================================
13. NO GENERAL-PURPOSE PHYSICS "JURISDICTION MATRIX"
============================================================

Do NOT revive the previous design as:

    MathType × PhysicsRole × Regime
        → allowed / forbidden

as the primary compiler architecture.

Instead prefer:

    typed operation signatures
    interfaces/capabilities
    constructors
    domain/codomain
    composition
    explicitly declared relations
    formal language semantics

The compiler should reject a program when a formal language contract is
violated.

Example:

    differentiate(causal_order)

should fail because the supplied object does not satisfy the formal input
contract of the differentiation operator.

It should NOT fail because a philosophical table says:

    "causal order is not calculus."

Likewise, a novel construct must remain expressible if the language grammar
and type/interface contracts permit it.

============================================================
14. OPEN PHYSICS MUST REMAIN EXPRESSIBLE
============================================================

The compiler must not become a gatekeeper of currently accepted physics.

A user must be able to define a new concept, interface, constructor, or
composition within the rules of the language.

For example:

    RelationalTime

may be introduced by the user.

If it is intended to satisfy:

    Temporal

the compiler should check whether its required formal contract exists.

It should not say:

    "RelationalTime is not in the approved physics ontology."

The language defines formal contracts.

Humans determine the physical interpretation and scientific significance.

============================================================
15. PROVENANCE
============================================================

The language may distinguish structurally between:

    assumed
    derived
    externally proved / certified

but do NOT turn these into an epistemic research ledger.

A useful direction from the previous plans is:

    Assumed[T]
    Derived[T]
    Proved[T]

or an equivalent compiler-native representation.

The implementation plan must ensure these constructors cannot be forged
merely by writing a label.

However:

    Proved[T]

must mean only that the value is associated with a recognized proof
certificate or proof artifact.

It MUST NOT mean:

    physically true.

============================================================
16. DERIVATION VS PHYSICAL COUPLING
============================================================

Preserve the important distinction:

    derivational dependency
        vs
    physical mutual relation

Derivation dependencies should have a well-defined formal structure.

Physical systems may contain mutual coupling:

    Matter ↔ Geometry

without implying deductive circularity.

The compiler must represent both without conflating them.

Do not simply copy the old two-graph implementation into the new compiler.

Redesign it naturally for the language/IR.

============================================================
17. SOURCE-LEVEL SEMANTICS
============================================================

The language should distinguish constructs such as:

    postulate
    define
    derive
    represent
    identify
    constrain
    couple
    measure
    construct
    recover
    emerge
    limit

where these are actually necessary.

Do NOT blindly preserve the five-mode AI-facing API from Program A.

This is a programming language now.

The language must be expressive enough to distinguish:

    definition
    assumption
    construction
    derivation
    representation
    physical relationship

without forcing all of these into one "relation + mode" object.

Determine the smallest coherent source-language vocabulary.

============================================================
18. SEMANTIC RULES
============================================================

The compiler should enforce formal language semantics such as:

    name resolution
    scope
    declaration-before-use where appropriate
    duplicate declarations
    undefined symbols
    type compatibility
    interface satisfaction
    argument count
    domain/codomain compatibility
    generic parameter compatibility
    composition validity
    dimensional consistency
    tensor/index validity
    invalid implicit conversions
    invalid provenance conversions
    invalid operator application
    invalid derivation dependencies
    import/module constraints
    regime/type qualifiers where those are actual language semantics

Do not convert philosophical judgments into hard compiler errors.

============================================================
19. ERROR MODEL
============================================================

Follow Go's diagnostic philosophy.

Diagnostics should be:

    precise
    source-position aware
    deterministic
    actionable
    multi-error where practical
    suggestion-capable

Define:

    error codes
    source spans
    related notes
    suggestions
    severity model

Do not use the old Program A artifact verdict model as the compiler's
primary user-facing model.

A compiler should naturally have:

    compile success
    compile failure

plus separate warnings/tool reports.

Do not create:

    WELL_FORMED
    DEFERRED
    UNKNOWN_TO_A

as the main compiler user experience.

Where useful, downstream analysis tools may report unresolved or
unproven constructs.

============================================================
20. GO-LIKE TOOLING
============================================================

Design a toolchain modeled on Go's ecosystem.

At minimum evaluate:

    physc        compiler
    phys vet     heuristic/static analysis
    phys doc     structural documentation/inspection
    phys fmt     source formatting
    phys types   reusable type-checking library
    phys test    language/conformance tests where justified

Do not blindly implement every command.

The implementation plan must determine the minimal useful v0 toolchain.

`phys vet` is where heuristic/suspicion analysis belongs.

The compiler itself should not pretend that every useful warning is a
hard language error.

`phys doc` should expose the structure of the program:

    declarations
    interfaces
    composition
    dependencies
    regimes
    dimensions
    exported objects
    IR structure

It must not interpret the physics for the user.

============================================================
21. STANDARD LIBRARY
============================================================

The standard library is a major part of the language design.

Follow the Go principle:

    SMALL COMPILER CORE
    +
    RICH STANDARD LIBRARY
    +
    EXTENSIBLE USER PACKAGES

Design a standard library hierarchy such as:

    std/math/...
    std/physics/...
    std/tensor/...
    std/dim/...
    std/analysis/...

but do not lock these exact names without repository justification.

Potential math areas:

    algebra
    linear
    calculus / analysis
    geometry
    topology
    measure
    probability
    tensor
    spectral

Potential physics areas:

    mechanics
    fields
    dynamics
    quantum
    relativity
    statistical
    measurement

The exact v0 set must be determined from concrete language requirements.

Standard library packages must be normal language/library constructs,
not hidden compiler law.

Frameworks such as:

    BM
    CDT
    CST
    AS

must live above the core standard library as user/framework packages.

They are libraries, not compiler truth.

============================================================
22. CONCEPT LIBRARY / EXTENSIBILITY
============================================================

The compiler core must not contain a giant list of every mathematical and
physical concept.

A user should be able to add something new through language/library
mechanisms without modifying compiler internals.

The plan must clearly distinguish:

    compiler primitives
    standard-library types
    standard-library interfaces
    user-defined types
    user-defined interfaces
    framework packages

and specify exactly what requires compiler changes versus ordinary library
changes.

============================================================
23. AST DESIGN
============================================================

The AST is now a real source-language AST.

Design it carefully.

It should represent:

    declarations
    identifiers
    literals
    type expressions
    generic parameters
    function/operator applications
    composition
    fields/indexing
    tensor/index syntax if v0
    relations
    derive blocks
    assumptions
    bridges/representation constructs
    annotations

Use Go-style AST interface patterns where useful:

    Node
    Expr
    Decl
    Stmt
    TypeExpr

with source positions/spans.

Do not create an enormous AST covering every imaginable piece of physics.

Define the minimum grammar needed by the v0 language.

============================================================
24. TYPED AST / TYPE INFORMATION
============================================================

Follow the Go compiler philosophy of keeping syntax representation
distinct from semantic type information.

Define how the compiler attaches:

    resolved symbols
    types
    interfaces
    dimensions
    provenance
    source spans
    constant/type information

to AST or companion structures.

Determine whether v0 uses:

    decorated AST
    side tables
    typed AST

and why.

============================================================
25. IR DESIGN
============================================================

The Physics IR is critical.

Do not make it simply:

    copy of AST as JSON.

It should be a canonical semantic representation suitable for:

    analysis
    export
    downstream computation
    Program B
    future backends

Define IR invariants.

At minimum consider representation of:

    declarations
    typed objects
    operations
    compositions
    interfaces/capabilities
    dimensions
    tensor/index structures
    dependencies
    assumptions
    provenance
    regimes
    physical relations
    expressions
    spans/source maps
    package/import identity

Determine which constructs disappear during lowering and which remain.

============================================================
26. EXPRESSION LANGUAGE
============================================================

Unlike Program A, the Physics Compiler DOES require a mathematical
expression language because the goal is now a real programming language
for physics.

However:

    do not invent the full universe of mathematics.

Define a minimal expression core supporting the v0 needs, such as:

    literals
    names
    application
    addition
    multiplication
    powers
    equality
    approximate equality if justified
    indexing
    derivatives
    integrals
    tensor contractions
    function/map application
    constructors

Then define how new mathematical constructs are introduced by standard
libraries.

Do not implement a CAS merely because the language can express equations.

Compiler representation ≠ symbolic consequence engine.

============================================================
27. PROGRAM B BOUNDARY
============================================================

The compiler should compile/represent a physics program.

Program B performs actual computational work.

Therefore define an explicit boundary:

    .phys source
       ↓
    Physics Compiler
       ↓
    Physics IR / export data
       ↓
    Program B / numerical backend / symbolic backend / proof backend

The compiler may lower mathematical operations into a structured IR.

It must not itself become:

    numerical simulator
    CAS
    theorem prover
    GPU execution engine

in v0.

Do not prematurely choose C/LLVM/OpenCL backends.

First stabilize language and IR.

============================================================
28. LEAN / PROOF SYSTEM BOUNDARY
============================================================

Lean 4 is NOT a core compiler dependency.

If the language can produce formal mathematical propositions requiring proof,
define an external proof boundary.

Conceptually:

    Physics Compiler
        ↓
    proof obligation / proposition
        ↓
    external verifier such as Lean

But the compiler must never turn:

    Lean proof
        →
    physical truth

A formal proof establishes only the formally stated mathematical proposition
under its formal assumptions.

Do not implement Lean integration in the v0 compiler unless the plan can
justify a minimal interface without complicating the compiler core.

============================================================
29. ANALYSIS WITHOUT INTERPRETATION
============================================================

The compiler toolchain should allow users to inspect the compiled program.

Useful reports may include:

    dependency graph
    type graph
    interface satisfaction
    composition tree
    assumptions
    provenance
    dimensions
    tensor/index structure
    regime usage
    exported declarations
    IR dump
    source-to-IR mapping

But reports must describe the representation.

They must not say:

    "this theory is correct"
    "this physics is wrong"
    "this hypothesis is likely true"
    "this theory is better"

Those belong to humans and separate research systems.

============================================================
30. NO PAPER-RESEARCH METHODOLOGY
============================================================

This requirement is absolute.

Do NOT include in compiler architecture:

    EBP
    literature evidence tracking
    adversarial paper litigation
    claim promotion
    scientific scorecards
    research workflow states
    paper review logic
    null-result adjudication
    theory ranking
    citation governance
    research context graph

The compiler may support source references or documentation as language
metadata, but it must not manage the research process.

============================================================
31. USE THE FOUR ATTACHED PLANS AS INPUTS, NOT AUTHORITIES
============================================================

You have four candidate designs:

    plan2_z.md
    plan2_qwen.md
    plan2_deepseek.md
    plan2_gemini.md

Do NOT average them.

Disposition them explicitly.

The intended consolidation is approximately:

    Z
      → primary architecture/philosophy/tooling/interface direction

    Qwen
      → concrete compiler pipeline and AST/typechecker organization

    DeepSeek
      → broader math/physics grammar inventory and vocabulary candidates

    Gemini
      → tensor/index/dimension representation and checking ideas

But independently verify every inherited idea against the locked principles.

Explicitly reject designs that introduce:

    giant jurisdiction matrices
    philosophical compiler laws
    Dancer/Dance/Stage as compiler ontology
    premature code generation
    premature execution backend
    theorem proving inside the compiler
    scientific truth judgments
    research-methodology machinery

============================================================
32. REQUIRED REPOSITORY INSPECTION
============================================================

Before producing the plan, inspect:

    existing repository structure
    existing Program A implementation
    Program A artifacts/IR
    existing concepts/specs
    examples
    tests
    module/package conventions
    Go version
    existing build/test tooling

Determine:

    what should be reused
    what should be treated as prototype/reference
    what should not be carried forward

Do not assume Program A's architecture should become the Physics
Compiler's architecture.

Treat Program A as a valuable experimental semantic prototype.

============================================================
33. MIGRATION / REUSE STRATEGY
============================================================

The plan must explicitly answer:

    What can Program A contribute?

Potential reusable pieces might include:

    concept registry ideas
    dimension model
    canonicalization ideas
    deterministic hashing
    provenance representation
    relation semantics
    test corpus knowledge
    examples

Potentially replaceable pieces might include:

    Go SDK-only authoring
    artifact-first frontend
    verdict system
    old MRC implementation
    old mode/relation surface
    old jurisdiction model

Do not preserve something solely because it already exists.

============================================================
34. VERSIONING / SPEC
============================================================

Follow Go's preference for a clearly defined language specification.

Define:

    language version
    package/module versioning
    standard library versioning
    compatibility policy
    source compatibility
    export-data compatibility if applicable

The plan must decide whether the language spec is:

    prose + conformance tests
    formal grammar + prose semantics
    another justified form

Do not confuse:

    language specification
    standard library
    research methodology

They are different artifacts.

============================================================
35. CONFORMANCE CORPUS
============================================================

The Physics Compiler needs a compiler/language conformance suite.

Design tests for:

    lexical errors
    parse errors
    scope/name errors
    type errors
    interface errors
    composition errors
    dimension errors
    tensor/index errors
    invalid conversions
    invalid generic instantiation
    invalid operator application
    derivation dependency errors
    package/import errors
    valid programs
    novel user-defined concepts
    physics-framework packages

Include canonical examples:

    Problem of Time
    Matter/Geometry backreaction

but transform them into LANGUAGE tests, not research adjudication tests.

For example:

    valid composition
    invalid interface satisfaction
    legal physical coupling
    invalid derivation dependency

not:

    "the compiler decides whether the Problem of Time is solved."

============================================================
36. GOLDEN EXAMPLES
============================================================

The plan must include a small number of representative `.phys` examples.

At minimum:

    hello-world / smallest valid program
    mathematical composition
    tensor/index example
    dimensional analysis
    field on spacetime
    symmetry/conservation structure
    causal-order/time example
    matter/geometry coupling example
    generic type example
    user-defined interface example
    framework import example

For each example explain:

    source
    AST structure
    type information
    semantic checks
    resulting IR

This is essential for validating that the proposed language actually
represents human intent rather than merely providing an impressive list
of compiler components.

============================================================
37. STANDARD LIBRARY TEST OFTEN, NOT LAST
============================================================

The standard library is not a final add-on.

The plan must use representative library packages early to validate:

    type system
    interfaces
    composition
    generics
    dimensions
    tensor/index grammar
    import/package behavior

The compiler core should remain small because the standard library carries
the richer domain vocabulary.

============================================================
38. ANALYSIS TOOL DESIGN
============================================================

Analysis tools should consume:

    AST
    typed AST
    Physics IR
    export data

rather than reimplementing compiler semantics.

Define a reusable inspection API.

Potential reports:

    `phys doc graph`
    `phys doc types`
    `phys doc interfaces`
    `phys doc dimensions`
    `phys doc deps`
    `phys doc ir`

Do not make these reports interpretive.

============================================================
39. SECURITY / EXECUTION
============================================================

`.phys` source is declarative/program source, but the toolchain will
eventually interact with external programs.

The plan must identify:

    what the compiler itself executes
    what it never executes
    where external execution occurs
    how untrusted input is handled

Do not put sandboxing or execution orchestration into the compiler core.

============================================================
40. IMPLEMENTATION STRUCTURE
============================================================

Propose a clean Go module/package structure.

The plan should likely resemble the conceptual separation:

    cmd/
    internal/scanner/
    internal/token/
    internal/parser/
    internal/ast/
    internal/types/
    internal/typecheck/
    internal/semantics/
    internal/ir/
    internal/lower/
    internal/export/
    internal/diagnostics/
    std/
    tools/
    testdata/

but choose the final package boundaries based on the actual repository.

Important:

    distinguish compiler internals
    from public reusable type-checking libraries
    from standard library packages
    from framework packages
    from analysis tools

Do not create packages merely to make the directory tree impressive.

============================================================
41. IMPLEMENTATION ORDER
============================================================

Produce one dependency-aware implementation sequence.

It should NOT postpone core language decisions.

The sequence should likely proceed conceptually through:

    repository assessment
       ↓
    language specification v0
       ↓
    lexical grammar
       ↓
    AST
       ↓
    names/scopes
       ↓
    type representation
       ↓
    interfaces/capabilities
       ↓
    composition
       ↓
    mathematical expressions
       ↓
    dimensions/index system
       ↓
    semantic analysis
       ↓
    typed Physics IR
       ↓
    minimal std library
       ↓
    compiler CLI
       ↓
    diagnostics
       ↓
    analysis tools
       ↓
    conformance corpus
       ↓
    framework packages
       ↓
    export/program-B boundary

But determine the true dependency order rather than blindly using this
sequence.

============================================================
42. NO PREMATURE CODEGEN
============================================================

Do NOT make v0 dependent on:

    LLVM
    C
    OpenCL
    GPU kernels
    CPU code generation
    numerical solver generation

The first compiler backend is:

    typed Physics IR / export representation

Later backends may consume it.

Language and IR must stabilize before heavy code generation.

============================================================
43. WHAT "DONE" MEANS FOR V0
============================================================

Define MVP completion in compiler terms.

It should include, at minimum:

    source scanner
    parser
    AST
    name resolution
    type system
    interfaces
    composition
    mathematical expression core
    physics grammar core
    dimension support
    tensor/index support to the explicitly defined v0 boundary
    semantic checking
    Physics IR
    minimal standard library
    compiler CLI
    diagnostics
    formatting if included
    conformance test suite
    representative `.phys` examples
    analysis/inspection tooling
    Program B export boundary

It does NOT include:

    scientific truth
    theorem proving
    full symbolic algebra
    numerical simulation
    paper review
    research management

============================================================
44. "NO LATER" ARCHITECTURAL COMPLETENESS AUDIT
============================================================

Before declaring the plan finished, explicitly answer YES/NO for:

    Can the language express mathematical objects?
    Can it express physical objects?
    Can mathematical structures compose?
    Can physical structures compose?
    Can interfaces express capabilities?
    Can user-defined interfaces extend the language?
    Can generic/parameterized structures exist?
    Can domains/codomains be represented?
    Can tensor/index structure be represented?
    Can dimensions be represented?
    Can expressions be represented?
    Can assumptions be represented?
    Can derivations be represented?
    Can physical couplings be represented?
    Can provenance be represented structurally?
    Can package/import structure be represented?
    Can the compiler type-check all of the above?
    Can the compiler distinguish syntax errors from type errors?
    Can it distinguish type errors from semantic-language errors?
    Can it preserve open physics?
    Can new concepts be added without modifying compiler internals?
    Can standard-library packages carry richer physics?
    Can framework packages exist outside compiler core?
    Can the typed representation be lowered to Physics IR?
    Can the IR be consumed by Program B?
    Can users inspect the compiled representation?
    Can users receive useful source diagnostics?
    Can the language evolve without rewriting the compiler architecture?

Every NO must either be solved in the plan or explicitly removed from v0
scope.

Do not hide unresolved architecture under "future work."

============================================================
45. REQUIRED PLAN OUTPUT
============================================================

Produce:

    PHYSICS_COMPILER_PLAN.md

The document must contain:

1. Executive Summary
2. Repository Assessment
3. Final Architectural Decisions
4. Conceptual vs Compiler Boundary
5. Language Philosophy
6. `.phys` Language Overview
7. Lexical Grammar
8. Syntax / EBNF or equivalent
9. AST Design
10. Name Resolution and Scoping
11. Mathematical Type System
12. Physical Type System
13. Interface / Capability System
14. Composition Model
15. Generic / Parameterized Types
16. Mathematical Expression System
17. Tensor / Index System
18. Dimension System
19. Physics Semantic Constructs
20. Semantic Analysis
21. Dependency / Derivation Model
22. Provenance Model
23. Physics IR
24. Lowering
25. Standard Library Architecture
26. Core Math Standard Library
27. Core Physics Standard Library
28. Framework Package Architecture
29. Compiler Toolchain
30. Diagnostics
31. `phys vet`
32. `phys doc`
33. Formatting / Testing tools
34. Program B Boundary
35. External Proof Boundary
36. Conformance Test Suite
37. Golden Examples
38. Repository / Package Layout
39. Implementation Sequence
40. Acceptance Criteria
41. Residual Limitations
42. NO-LATER Audit
43. Decisions inherited from the four candidate plans
44. Decisions explicitly rejected from the four plans

============================================================
46. REQUIRED COMPARATIVE DISPOSITION
============================================================

Include a table:

    Idea
    Source plan
    Adopt?
    Why?
    Where implemented?

For every major idea from the four plans, classify:

    ADOPT
    ADAPT
    REJECT

Do not merely say one plan is "best."

Show exactly how the consolidated design was derived.

============================================================
47. REQUIRED DESIGN ARTIFACTS INSIDE THE PLAN
============================================================

The plan must contain concrete examples of:

    representative `.phys` source
    token classes
    AST node families
    type structures
    interface declarations
    generic type examples
    composition examples
    dimension representation
    tensor/index representation
    semantic errors
    typed IR
    standard-library package structure
    framework package structure
    compiler command examples

The purpose is to prove that the architecture can actually represent
physics and mathematics.

============================================================
48. FINAL DESIGN PRINCIPLE
============================================================

The plan should culminate in this architecture:

    HUMAN / AI
        ↓
    .phys SOURCE
        ↓
    SCANNER
        ↓
    PARSER
        ↓
    AST
        ↓
    NAME RESOLUTION
        ↓
    TYPE SYSTEM
        ↓
    INTERFACES / CAPABILITIES
        ↓
    COMPOSITION
        ↓
    SEMANTIC ANALYSIS
        ↓
    TYPED PHYSICS IR
        ↓
    EXPORT / ANALYSIS
        ↓
    PROGRAM B / EXTERNAL TOOLS

with:

    STANDARD LIBRARY
        ↑
    compiler core

and:

    FRAMEWORK PACKAGES
        ↑
    standard/user packages

The conceptual layer sits beside the compiler:

    conceptual models
        →
    documentation / examples / human analysis

not inside the compiler's ontology.

============================================================
49. FINAL INSTRUCTION
============================================================

Do not implement code.

Do not ask me to choose between multiple alternatives unless the repository
makes the choice genuinely impossible.

Where the architecture above is already decided, make the decision.

Where a detail is genuinely undecided, choose the smallest design that is
consistent with:

    Go compiler engineering
    interfaces
    composition
    explicit grammar
    strong typing
    extensibility
    mathematical precision
    physical expressiveness
    no interpretation
    no adjudication
    no research methodology

Then explain why.

The objective is to leave a single coherent implementation blueprint
for a REAL PROGRAMMING LANGUAGE AND COMPILER FOR PHYSICS.

The old Program A semantic checker was a prototype.

The Physics Compiler is the language-level successor.

Do not accidentally rebuild Program A with a parser attached.

Build the compiler architecture around:

    GRAMMAR
    TYPES
    INTERFACES
    COMPOSITION
    EXPRESSIONS
    PHYSICS STRUCTURES
    MATHEMATICAL STRUCTURES
    IR
    STANDARD LIBRARY
    TOOLING

and nothing that belongs to scientific adjudication.
```

### One important design correction embedded in this prompt

The biggest conceptual shift is this:

```text
OLD PROGRAM A

Concept → relation → MRC → verdict
```

becomes:

```text
PHYSICS COMPILER

Syntax
  ↓
Type
  ↓
Interface
  ↓
Composition
  ↓
Formal semantic validity
  ↓
IR
```

That is much closer to the Go model. The Z proposal's strongest idea is that interfaces can encode **capability grammar**, while types and operations carry the mechanically enforceable structure.  

And the prompt intentionally rejects the parts that would pull us backward: a hardcoded jurisdiction system, conceptual Dancer/Dance/Stage types, early LLVM/C/GPU backends, and compiler-level scientific adjudication. The candidate plans contain all of those temptations in different forms.  

**This is the prompt I would now use to draft the Physics Compiler implementation plan.**


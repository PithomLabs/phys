## 1. Compiler Architecture & Standard Library Layout

Following Go's compiler layout (`cmd/compile`), the toolchain cleanly splits into a deterministic parsing/type-checking frontend, a canonical Intermediate Representation (IR), code-generation backends, and a modular standard library (`std/`).

```text
Source Code (.phys)
      │
      ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ cmd/compile │ ──> │ cmd/compile │ ──> │ cmd/compile │ ──> │ cmd/compile │
│   /syntax   │     │   /types    │     │     /ir     │     │   /codegen  │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
  Lexing & AST       Units, Tensors,      Math/Physics        C/LLVM/OpenCL
  Constructs         & Index Rules        Canonical DAG       Executables
                                               │
                                               ▼
                                      ┌─────────────────┐
                                      │ std/            │
                                      │  ├─ math/       │
                                      │  ├─ physics/    │
                                      │  └─ analysis/   │
                                      └─────────────────┘

```

### Core Compiler Packages (`cmd/compile/`)

* **`syntax`**: Lexes and parses source declarations into explicit Abstract Syntax Tree (AST) nodes without performing physical or mathematical validation.
* **`types`**: Validates structural constraints: dimension consistency ($M^a L^b T^c$), index contraction matching, and operator type compatibility.
* **`ir`**: Transforms AST into a canonical directed acyclic graph (DAG) representing exact mathematical operations and state transformations.
* **`codegen`**: Emits high-performance evaluation targets (C99, LLVM IR, CPU/GPU kernels) and introspection metadata.

### Standard Library Packages (`std/`)

* **`std/math/algebra`**: Provides concrete types for groups, Clifford/Lie algebras, matrices, and non-commutative brackets.
* **`std/math/calculus`**: Defines differential operators ($\partial_\mu$, $\nabla$), integration domains, and boundary conditions.
* **`std/physics/units`**: Standard SI/Natural unit definitions and dimensional basis vectors.
* **`std/physics/fields`**: Generic field types (`ScalarField`, `VectorField`, `TensorField`) over specified manifolds.
* **`std/analysis/inspect`**: Reflection and trace primitives enabling runtime inspection of intermediate states, conservation laws, and tensor indices.

---

## 2. Separation: Conceptual Ideas vs. Compiler Invariants

To keep the compiler focused on mechanical evaluation rather than physical interpretation, human mental models are separated from compiler IR invariants.

| Layer | Conceptual Ideas (Human Mental Models) | Compiler IR Invariants (Math & Physics) |
| --- | --- | --- |
| **Purpose** | Documentation, interpretation, naming conventions, paper cross-references. | Executable state transformations, type checking, memory layout, operational verification. |
| **Compiler Treatment** | Stripped or attached purely as debug metadata (`DocString`, `Annotation`, `SymbolicTag`). | Hard type bounds (`TensorRank`, `DimensionalVector`, `OperatorSignature`). |
| **Examples** | "Bohmian trajectory", "Gauge fixing choice", "Unparticles", "Dark Energy field". | $d$-dimensional vector array, scalar differential equation, Lie group generator $T^a$, array memory layout. |
| **Enforcement** | Unenforced. The compiler does not judge whether a physical model is "correct" or "real". | Enforced statically. Strict index contractions, dimensional balance, valid operator domain/codomain. |

---

## 3. Math and Physics Grammar Representation

The AST and IR explicitly mirror mathematical and physical grammar. Conceptual terms are replaced by concrete mathematical types and operations.

### AST Node Definitions

```go
// TensorNode represents tensor definitions and index configurations.
type TensorNode struct {
    Name       string
    Rank       int
    Covariant  []string // e.g., ["mu", "nu"]
    Contravariant []string // e.g., ["alpha"]
    ValueType  Type     // e.g., Float64, Complex128
}

// DifferentialOpNode represents formal calculus operations.
type DifferentialOpNode struct {
    Operator   string // e.g., "Partial", "Covariant", "Laplacian"
    Index      string // e.g., "mu"
    Target     Node   // Field or Expression
}

// CommutatorNode represents algebraic non-commutativity.
type CommutatorNode struct {
    OpA, OpB   Node
    IsAntisym  bool // true = [A, B], false = {A, B}
}

```

### Type Checking & Grammar Rules

1. **Index Contraction Grammar**: An index appearing once in a contravariant position and once in a covariant position represents an implicit or explicit summation ($\mathbf{A}^\mu \mathbf{B}_\mu$). Duplicated indices with matching variance (e.g., $\mathbf{A}_\mu \mathbf{B}_\mu$ without a metric) trigger a compile-time error (`InvalidIndexVarianceError`).
2. **Dimensional Analysis Grammar**: Algebraic addition and subtraction require identical physical dimensions. Multiplication compounds dimensions.

$$\text{Type}(A \pm B) \implies \text{Dim}(A) \equiv \text{Dim}(B)$$


3. **Operator Domain Rule**: Differential operators acting on field types must return explicit derivative fields with updated tensor ranks or index structures (e.g., $\partial_\mu \phi \implies \text{Rank 1 Covariant Field}$).

---

## 4. Syntactic & Semantic Guidelines

### Syntactic Guidelines

* **Explicit over Implicit**: No implicit metric conversions. Lowering or raising indices requires an explicit contraction call with a metric tensor ($g_{\mu\nu} A^\nu$).
* **Unambiguous Index Bounding**: Tensor expressions must declare index scopes explicitly or utilize Standard Einstein Summation syntax validated during type analysis.
* **Separation of Definition and State**: Defining a physical field equation ($\mathcal{L} = \frac{1}{2} \partial_\mu \phi \partial^\mu \phi - V(\phi)$) is syntactically distinct from instantiating a numerical solver grid with initial/boundary conditions.

### Semantic Guidelines (Non-Arbitrating Scope)

* **No Physics Litigation**: The compiler never rejects a program on the grounds that a theoretical physics model is unphysical, untested, or mathematically speculative. If the syntax, units, tensor contractions, and linear operations are structurally valid, the program compiles.
* **Deterministic Transformation**: The compiler's single duty is to compile human-written operational intent into verifiable computations and output inspection streams.

### Analysis & Inspection Interface (`std/analysis/inspect`)

To allow users to evaluate program output and mathematical consistency without compiler bias, the standard library exposes direct inspection hooks:

```go
package main

import (
    "std/analysis/inspect"
    "std/math/calculus"
    "std/physics/fields"
)

func AnalyzeField(f fields.ScalarField) {
    // Inspect index balance and dimensional properties statically/dynamically
    inspect.DumpAST(f)
    inspect.AssertDimension(f, "Energy / Volume")
    
    // Output state traces for user evaluation
    trace := inspect.TraceDivergence(calculus.Grad(f))
    inspect.WriteVTK("output_grid.vtk", trace)
}

```

End-users evaluate physics paper hypotheses by running programs and analyzing emitted data grids, tensor traces, and numerical convergence metrics—keeping the compiler an objective tool for execution and representation.

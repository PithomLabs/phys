## 1. The "Pen-and-Paper" AI Paradigm

If the ultimate goal of the Physics Compiler is to enable AI agents to perform theoretical derivations like Einstein—manipulating concepts, conservation laws, and invariant principles on a symbolic "paper" substrate—then the compiler's role is **formal proof generation and structural verification**.

```text
       PHYSICS LITERATURE / PAPER
                   │
                   ▼
     AI PROPOSER AGENT (Translation)
     • Maps notation to physical concepts
     • Emits Go code with physical assertions
                   │
                   ▼
┌──────────────────────────────────────────────┐
│  GO COMPILER + STATIC MRC ANALYZER (physvet) │
│  • Enforces Math-to-Physics Bridges          │
│  • Validates Dimensional Homogeneity         │
│  • Verifies Index Contraction & Variance     │
│  • Fails build on physical category errors   │
└──────────────────┬───────────────────────────┘
                   │
                   ▼
    SYMBOLIC DERIVATION ENGINE (pkg/phys)
    • Executes formal transformations (d/dt, δS, Lim)
    • Builds deterministic Derivation AST Trace
                   │
                   ▼
    ADVERSARIAL REVIEWER AGENT (Verification)
    • Inspects Derivation Trace
    • Challenges assumptions & boundary limits
    • Emits formal verification / counter-proof

```

Instead of evaluating numbers, the engine operates on **symbolic physical equivalence classes**, keeping an immutable ledger of every algebraic step, functional variation, and limit evaluation.

---

## 2. Enforcing MRC at the Compiler Level (`physvet`)

Delegating Math Reality Check (MRC) to prompt guidelines allows hallucination leaks. To keep Go as the source language while guaranteeing structural rigor, MRC is enforced at compile time using a custom Go static analyzer tool (`physvet` built via `golang.org/x/tools/go/analysis`).

When an AI agent runs `go test` or `physvet ./...`, the analyzer intercepts the Go AST before symbolic execution and applies three static verification passes:

```text
Go Source File (.go)
        │
        ▼
┌────────────────────────────────────────────────────────┐
│ Pass 1: Category & Bridge Audit                        │
│ - Rejects unbridged raw constants / bare math scalars. │
│ - Requires explicit physical constructors:             │
│   relativity.AsStressEnergy(tensor)                    │
└───────────────────────┬────────────────────────────────┘
                        │
                        ▼
┌────────────────────────────────────────────────────────┐
│ Pass 2: Dimensional Vector Flow                        │
│ - Traces [M, L, T, Q, Θ] statically across operators. │
│ - Flags addition/equality of non-matching dimensions. │
└───────────────────────┬────────────────────────────────┘
                        │
                        ▼
┌────────────────────────────────────────────────────────┐
│ Pass 3: Tensor Variance & Non-Commutativity Audit     │
│ - Verifies covariant/contravariant index pairing.     │
│ - Flags improper reordering of non-commuting operators.│
└───────────────────────┬────────────────────────────────┘
                        │
                        ▼
      PASS: Proceed to Symbolic Engine
      FAIL: Abort Build with Detailed MRC Error Trace

```

### Static MRC Verification Example

If an AI agent attempts to construct an equation assigning a raw temperature scalar directly to energy without a dimensional bridge:

```go
// INVALID GO CODE (Caught statically by physvet)
T := statmech.Temperature("T")
E := relativity.Energy("E")

eq := operations.Equal(E, T) // Compiler Error!

```

`physvet` throws a build error during analysis:

```text
./derivation.go:14:18: MRC Category Error: Cannot equate Physical Type 'Energy' [M^1 L^2 T^-2] 
with Physical Type 'Temperature' [Θ^1] without explicit Boltzmann Bridge: thermal.ScaleBoltzmann(T)

```

---

## 3. Concrete Symbolic Derivation: $E = m_0 c^2$

Here is how an AI Proposer agent uses the Go library to derive mass-energy equivalence from Special Relativity invariants, generating a machine-readable derivation trace for the Adversarial Reviewer agent.

### A. The Proposer Agent Script (`derivation.go`)

```go
package main

import (
	"fmt"

	"github.com/physics-compiler/phys/pkg/phys/calculus"
	"github.com/physics-compiler/phys/pkg/phys/inspect"
	"github.com/physics-compiler/phys/pkg/phys/operations"
	"github.com/physics-compiler/phys/pkg/phys/relativity"
	"github.com/physics-compiler/phys/pkg/phys/units"
)

// DeriveMassEnergyEquivalence executes symbolic steps to derive E = m_0 * c^2
func DeriveMassEnergyEquivalence() (operations.Equation, *inspect.DerivationTrace) {
	trace := inspect.NewDerivationTrace("Relativistic Mass-Energy Equivalence")

	// 1. Postulate: Spacetime Four-Momentum Invariant
	// P^μ P_μ = - (m_0 * c)^2
	m0 := relativity.RestMass("m_0")
	c := units.SpeedOfLight
	
	invariantRHS := operations.Negate(operations.Square(operations.Multiply(m0, c)))
	P := relativity.FourMomentum("P")
	invariantLHS := relativity.MinkowskiDot(P, P)

	postulate := operations.Equal(invariantLHS, invariantRHS)
	trace.AddStep("Postulate Four-Momentum Norm Invariant", postulate)

	// 2. Expand Four-Momentum in terms of Energy E and 3-Momentum p:
	// P^μ = (E/c, p_x, p_y, p_z)  =>  P^μ P_μ = - (E/c)^2 + |p|^2
	E := relativity.Energy("E")
	p := relativity.ThreeMomentum("p")

	expandedLHS := operations.Subtract(
		operations.Square(p),
		operations.Square(operations.Divide(E, c)),
	)
	
	expandedEq := operations.Equal(expandedLHS, invariantRHS)
	trace.AddStep("Expand P^μ in Components (E/c, p)", expandedEq)

	// 3. Algebraic Transformation: Solve for E^2
	// E^2 = (m_0 * c^2)^2 + (p * c)^2
	eSquared := operations.Add(
		operations.Square(operations.Multiply(m0, operations.Square(c))),
		operations.Square(operations.Multiply(p, c)),
	)
	energyDispersion := operations.Equal(operations.Square(E), eSquared)
	trace.AddStep("Solve for Energy Dispersion relation E^2(p)", energyDispersion)

	// 4. Rest Frame Condition: Take limit as 3-momentum p -> 0
	restFrameEq := calculus.Limit(energyDispersion, p, operations.ZeroVector)
	trace.AddStep("Apply Rest-Frame Constraint (p -> 0)", restFrameEq)
	// Output at this step: E^2 = (m_0 * c^2)^2

	// 5. Derive Positive Energy Root: E = m_0 * c^2
	finalEquation := operations.SquareRoot(restFrameEq)
	trace.AddStep("Extract Positive Energy Eigenvalue", finalEquation)

	return finalEquation, trace
}

func main() {
	eq, trace := DeriveMassEnergyEquivalence()
	fmt.Println("Derived Equation:", inspect.Render(eq))
	fmt.Println("\n--- Derivation Trace for Adversarial Review ---")
	trace.DumpJSON()
}

```

---

### B. Machine-Readable Derivation Trace (Emitted AST)

This JSON trace is passed directly to the **Adversarial Reviewer AI** to audit every step of the proof:

```json
{
  "title": "Relativistic Mass-Energy Equivalence",
  "steps": [
    {
      "step": 1,
      "rule": "POSTULATE_MINKOWSKI_INVARIANT",
      "expression": "Equal(MinkowskiDot(P, P), Negate(Square(Multiply(RestMass(m_0), SpeedOfLight))))",
      "latex": "P^\\mu P_\\mu = -(m_0 c)^2"
    },
    {
      "step": 2,
      "rule": "COMPONENT_EXPANSION",
      "expression": "Equal(Subtract(Square(ThreeMomentum(p)), Square(Divide(Energy(E), SpeedOfLight))), Negate(Square(Multiply(RestMass(m_0), SpeedOfLight))))",
      "latex": "|\\mathbf{p}|^2 - \\frac{E^2}{c^2} = -(m_0 c)^2"
    },
    {
      "step": 3,
      "rule": "SYMBOLIC_SOLVE_VARIABLE",
      "target": "Energy(E)",
      "expression": "Equal(Square(Energy(E)), Add(Square(Multiply(RestMass(m_0), Square(SpeedOfLight))), Square(Multiply(ThreeMomentum(p), SpeedOfLight))))",
      "latex": "E^2 = (m_0 c^2)^2 + (|\\mathbf{p}|c)^2"
    },
    {
      "step": 4,
      "rule": "LIMIT_EVALUATION",
      "condition": "p -> 0",
      "expression": "Equal(Square(Energy(E)), Square(Multiply(RestMass(m_0), Square(SpeedOfLight))))",
      "latex": "E^2 = (m_0 c^2)^2"
    },
    {
      "step": 5,
      "rule": "POSITIVE_ROOT_EXTRACTION",
      "expression": "Equal(Energy(E), Multiply(RestMass(m_0), Square(SpeedOfLight)))",
      "latex": "E = m_0 c^2"
    }
  ]
}

```

---

## 4. How the Adversarial Reviewer Audits the Trace

When the Adversarial Reviewer receives this output, it does not re-read paper prose. It executes structural verification passes against the derivation trace:

1. **Conservation & Invariance Verification**: Validates that Step 1 uses a valid Lorentz invariant ($P^\mu P_\mu$).
2. **Dimension Flow Audit**:

$$\text{Step 3 RHS: } [M^1 L^2 T^{-2}]^2 + [M^1 L^2 T^{-2}]^2 = [M^2 L^4 T^{-4}]$$



Confirms dimensional homogeneity across all addition operations.
3. **Boundary Condition Check**: Tests the limit step ($p \to 0$) to ensure no division by zero or invalid dimension collapse occurred.
4. **Physical Sanity Verdict**: Confirms that for $m_0 > 0$ and $c > 0$, $E > 0$, verifying the derived relation $E = m_0 c^2$ as a formally valid result.
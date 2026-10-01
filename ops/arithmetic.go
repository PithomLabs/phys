// Package ops implements the twelve pure symbolic operations of
// specs_v2_3.md §15. Every operation is a closed function over core.Object:
// it never mutates inputs, never touches session or ambient state, validates
// object validity first, applies the applicable MRC rules, merges
// assumptions/conventions, builds a fresh immutable result through the
// internal kernel mint, and propagates provenance by the exact §13.2 law.
package ops

import (
	"math/big"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

// ---------------------------------------------------------------------------
// Shared private helpers (validity, MRC pipeline, merge, provenance law)
// ---------------------------------------------------------------------------

// validInputs enforces MRC-001 first for every operation.
func validInputs(op string, objs ...core.Object) error {
	for _, o := range objs {
		if !o.Valid() {
			return core.InvalidObjectError{Operation: op}
		}
	}
	return nil
}

// dimensionString renders a canonical dimension for diagnostics.
func dimensionString(o core.Object) string {
	b, err := o.Dimension().CanonicalJSON()
	if err != nil {
		return ""
	}
	return string(b)
}

// checkAdditiveCompat applies MRC-002 then MRC-003 for Add/Subtract
// (specs_v2_3.md §6.2 additive table).
func checkAdditiveCompat(op string, a, b core.Object) error {
	if !a.Dimension().Equal(b.Dimension()) {
		return core.DimensionMismatchError{
			Operation: op,
			Left:      dimensionString(a),
			Right:     dimensionString(b),
		}
	}
	ka, kb := a.Kind(), b.Kind()
	switch {
	case ka == core.KindExpression && kb == core.KindExpression:
		return nil
	case ka == core.KindExpression || kb == core.KindExpression:
		return core.CategoryMismatchError{Operation: op, Left: ka.String(), Right: kb.String()}
	case ka != kb:
		return core.CategoryMismatchError{Operation: op, Left: ka.String(), Right: kb.String()}
	default:
		return nil
	}
}

// mergeInputs applies MRC-004 (assumptions) then MRC-005 (conventions):
// deterministic set union over the inputs plus any operation-required
// assumptions, with exact keyed conflict detection (§11.4).
func mergeInputs(inputs []core.Object, extra []core.Assumption) (core.AssumptionSet, core.ConventionSet, error) {
	assumptions := inputs[0].Assumptions()
	conventions := inputs[0].Conventions()
	var err error
	for _, o := range inputs[1:] {
		if assumptions, err = assumptions.Merge(o.Assumptions()); err != nil {
			return core.AssumptionSet{}, core.ConventionSet{}, err
		}
	}
	if len(extra) > 0 {
		if assumptions, err = assumptions.Merge(core.NewAssumptionSet(extra...)); err != nil {
			return core.AssumptionSet{}, core.ConventionSet{}, err
		}
	}
	for _, o := range inputs[1:] {
		if conventions, err = conventions.Merge(o.Conventions()); err != nil {
			return core.AssumptionSet{}, core.ConventionSet{}, err
		}
	}
	return assumptions, conventions, nil
}

// resultStatus implements the exact §13.2 pure-operation status law:
// any HYPOTHESIS input contaminates the output; otherwise DERIVED.
func resultStatus(inputs []core.Object) core.ProvenanceStatus {
	for _, o := range inputs {
		if o.Provenance().Status() == core.StatusHypothesis {
			return core.StatusHypothesis
		}
	}
	return core.StatusDerived
}

// inheritString implements §13.5 Source/Framework inheritance: preserve iff
// all input values are identical and non-empty, otherwise empty.
func inheritString(values []string) string {
	if len(values) == 0 {
		return ""
	}
	first := values[0]
	if first == "" {
		return ""
	}
	for _, v := range values[1:] {
		if v != first {
			return ""
		}
	}
	return first
}

// buildResult completes the operation pipeline: merge metadata (MRC-004/005),
// propagate provenance (§13.2/§13.5, parents = input canonical object hashes
// in input order, REQ-013-03), and mint the fresh immutable result through
// the single kernel mint authority with corpus status NONE.
func buildResult(kind core.Kind, dim core.Dimension, expr core.Expr,
	inputs []core.Object, extra []core.Assumption) (core.Object, error) {
	assumptions, conventions, err := mergeInputs(inputs, extra)
	if err != nil {
		return core.Object{}, err
	}
	sources := make([]string, 0, len(inputs))
	frameworks := make([]string, 0, len(inputs))
	parents := make([]string, 0, len(inputs))
	for _, o := range inputs {
		p := o.Provenance()
		sources = append(sources, p.Source())
		frameworks = append(frameworks, p.Framework())
		parents = append(parents, kernel.Hex(core.HashObject(o)))
	}
	provenance, err := kernel.NewProvenance(
		resultStatus(inputs),
		inheritString(sources),
		inheritString(frameworks),
		parents,
		core.HashAssumptionSet(assumptions),
		core.HashConventionSet(conventions),
		kernel.MRCVersion,
		"",
	)
	if err != nil {
		return core.Object{}, err
	}
	return kernel.MintObject(kernel.ObjectSpec{
		Name:         "",
		Kind:         kind,
		Dimension:    dim,
		Expr:         expr,
		Assumptions:  assumptions,
		Conventions:  conventions,
		Provenance:   provenance,
		CorpusStatus: kernel.CorpusNone,
	})
}

// denominatorPrecondition builds the §11.5 required operation assumption for
// a symbolic division: MathPrecondition keyed by the denominator's canonical
// expression hash with the structured value b != 0.
func denominatorPrecondition(b core.Object) (core.Assumption, error) {
	key := "denominator/" + kernel.Hex(core.HashExpr(b.Expr()))
	return core.NewExprAssumption(
		core.AssumptionMathPrecondition,
		key,
		core.NewRelation(core.RelationNeq, b.Expr(), core.NewRational(big.NewRat(0, 1))),
	)
}

// ---------------------------------------------------------------------------
// §15.1–§15.5 pure arithmetic operations
// ---------------------------------------------------------------------------

// Add implements specs_v2_3.md §15.1: equal dimensions, compatible kinds,
// merged metadata, same Kind and Dimension as the operands.
func Add(a, b core.Object) (core.Object, error) {
	if err := validInputs("add", a, b); err != nil {
		return core.Object{}, err
	}
	if err := checkAdditiveCompat("add", a, b); err != nil {
		return core.Object{}, err
	}
	return buildResult(a.Kind(), a.Dimension(), core.NewAdd(a.Expr(), b.Expr()),
		[]core.Object{a, b}, nil)
}

// Subtract implements specs_v2_3.md §15.2 with the same compatibility
// contract as Add.
func Subtract(a, b core.Object) (core.Object, error) {
	if err := validInputs("subtract", a, b); err != nil {
		return core.Object{}, err
	}
	if err := checkAdditiveCompat("subtract", a, b); err != nil {
		return core.Object{}, err
	}
	return buildResult(a.Kind(), a.Dimension(),
		core.NewAdd(a.Expr(), core.NewNeg(b.Expr())),
		[]core.Object{a, b}, nil)
}

// Multiply implements specs_v2_3.md §15.3: result kind Expression, result
// dimension a.dimension * b.dimension, operand kinds unconstrained.
func Multiply(a, b core.Object) (core.Object, error) {
	if err := validInputs("multiply", a, b); err != nil {
		return core.Object{}, err
	}
	return buildResult(core.KindExpression, a.Dimension().Multiply(b.Dimension()),
		core.NewMul(a.Expr(), b.Expr()),
		[]core.Object{a, b}, nil)
}

// Divide implements specs_v2_3.md §15.4: result kind Expression, result
// dimension a.dimension / b.dimension, plus the exact §11.5 denominator
// precondition.
func Divide(a, b core.Object) (core.Object, error) {
	if err := validInputs("divide", a, b); err != nil {
		return core.Object{}, err
	}
	precondition, err := denominatorPrecondition(b)
	if err != nil {
		return core.Object{}, err
	}
	quotient := core.NewMul(a.Expr(), core.NewPow(b.Expr(), big.NewRat(-1, 1)))
	return buildResult(core.KindExpression, a.Dimension().Divide(b.Dimension()),
		quotient,
		[]core.Object{a, b}, []core.Assumption{precondition})
}

// Pow implements specs_v2_3.md §15.5: the exponent is cloned by the kernel
// constructor (never retained as a caller-owned pointer), result kind
// Expression, result dimension base.dimension ^ exponent.
func Pow(base core.Object, exponent *big.Rat) (core.Object, error) {
	if err := validInputs("pow", base); err != nil {
		return core.Object{}, err
	}
	if exponent == nil {
		return core.Object{}, core.UnsupportedOperationError{
			Operation: "pow",
			Reason:    "exact rational exponent required",
		}
	}
	return buildResult(core.KindExpression, base.Dimension().Pow(exponent),
		core.NewPow(base.Expr(), exponent),
		[]core.Object{base}, nil)
}

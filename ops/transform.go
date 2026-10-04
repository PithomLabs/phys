package ops

import (
	"errors"
	"math/big"

	"github.com/PithomLabs/phys/core"
)

// ---------------------------------------------------------------------------
// Transform operations (specs_v2_3.md §15.7–§15.9)
// ---------------------------------------------------------------------------

// Substitute implements specs_v2_3.md §15.7: it replaces every occurrence of
// the variable's single symbol with the replacement expression. The target
// kind is unchanged; assumptions/conventions/provenance merge under the
// common rules. No simplification is performed — §15.7 does not require it
// (the golden trace calls Simplify explicitly as its own step).
func Substitute(target, variable, replacement core.Object) (core.Object, error) {
	const op = "substitute"
	if err := validInputs(op, target, variable, replacement); err != nil {
		return core.Object{}, err
	}
	// MRC-002 before MRC-003.
	if !replacement.Dimension().Equal(variable.Dimension()) {
		return core.Object{}, core.DimensionMismatchError{
			Operation: op,
			Left:      dimensionString(replacement),
			Right:     dimensionString(variable),
		}
	}
	if variable.Kind() != replacement.Kind() {
		return core.Object{}, core.CategoryMismatchError{
			Operation: op,
			Left:      variable.Kind().String(),
			Right:     replacement.Kind().String(),
		}
	}
	sym, err := singleSymbol(op, variable)
	if err != nil {
		return core.Object{}, err
	}
	expr, err := substituteExpr(target.Expr(), sym, replacement.Expr())
	if err != nil {
		return core.Object{}, err
	}
	return buildResult(target.Kind(), target.Dimension(), expr,
		[]core.Object{target, variable, replacement}, nil)
}

// Differentiate implements specs_v2_3.md §15.8: the bounded derivative
// engine. The result is simplified, has Kind=Expression, and carries the
// exact dimension target/wrt. No Derivative node is created.
func Differentiate(target, wrt core.Object) (core.Object, error) {
	const op = "differentiate"
	if err := validInputs(op, target, wrt); err != nil {
		return core.Object{}, err
	}
	sym, err := singleSymbol(op, wrt)
	if err != nil {
		return core.Object{}, err
	}
	assumptions, _, err := mergeInputs([]core.Object{target, wrt}, nil)
	if err != nil {
		return core.Object{}, err
	}
	deriv, err := differentiateExpr(target.Expr(), sym)
	if err != nil {
		return core.Object{}, err
	}
	simplified, err := simplifyExpr(deriv, assumptions)
	if err != nil {
		return core.Object{}, err
	}
	return buildResult(core.KindExpression,
		target.Dimension().Divide(wrt.Dimension()), simplified,
		[]core.Object{target, wrt}, nil)
}

// Limit implements specs_v2_3.md §15.9: expand the fixed §15.9.1 function
// body, substitute the value for the variable's single symbol, then
// simplify. The result has Kind=Expression and the target dimension.
// Singular or unsupported results surface as UnsupportedOperationError.
func Limit(target, variable, value core.Object) (core.Object, error) {
	const op = "limit"
	if err := validInputs(op, target, variable, value); err != nil {
		return core.Object{}, err
	}
	if !value.Dimension().Equal(variable.Dimension()) {
		return core.Object{}, core.DimensionMismatchError{
			Operation: op,
			Left:      dimensionString(value),
			Right:     dimensionString(variable),
		}
	}
	sym, err := singleSymbol(op, variable)
	if err != nil {
		return core.Object{}, err
	}
	assumptions, _, err := mergeInputs([]core.Object{target, variable, value}, nil)
	if err != nil {
		return core.Object{}, err
	}
	// 1. expand the fixed body (§15.9.1 — traversal, never an ID shortcut);
	// 2. substitute the value for the variable symbol;
	// 3. simplify exact zeros, powers, roots, and the denominator.
	expanded := expandCalls(target.Expr())
	if containsCall(expanded) {
		// Every admissible Call has a fixed one-argument body; anything
		// left after expansion cannot be evaluated by the bounded engine.
		return core.Object{}, core.UnsupportedOperationError{
			Operation: op,
			Reason:    "limit cannot expand the contained Call",
		}
	}
	expr, err := substituteExpr(expanded, sym, value.Expr())
	if err != nil {
		return core.Object{}, err
	}
	simplified, err := simplifyExpr(expr, assumptions)
	if err != nil {
		return core.Object{}, rewrap(op, err)
	}
	return buildResult(core.KindExpression, target.Dimension(), simplified,
		[]core.Object{target, variable, value}, nil)
}

// rewrap attributes an inner unsupported/singular failure to the operation
// whose contract produced it, preserving the reason text.
func rewrap(op string, err error) error {
	var u core.UnsupportedOperationError
	if errors.As(err, &u) {
		return core.UnsupportedOperationError{Operation: op, Reason: u.Reason}
	}
	return err
}

// ---------------------------------------------------------------------------
// Shared contract helpers
// ---------------------------------------------------------------------------

// singleSymbol enforces "the operand's expression is a single Symbol node"
// (§15.7, §15.8, §15.9) and returns the symbol name.
func singleSymbol(op string, o core.Object) (string, error) {
	e := o.Expr()
	if e.Kind() != core.ExprSymbol {
		return "", core.UnsupportedOperationError{
			Operation: op,
			Reason:    "operand expression must be a single Symbol node",
		}
	}
	return e.SymbolName(), nil
}

// ---------------------------------------------------------------------------
// Expression-level substitution
// ---------------------------------------------------------------------------

// substituteExpr replaces every occurrence of sym with repl, rebuilding the
// tree through the kernel constructors (flatten/combine/sign/sort re-run).
func substituteExpr(e core.Expr, sym string, repl core.Expr) (core.Expr, error) {
	switch e.Kind() {
	case core.ExprSymbol:
		if e.SymbolName() == sym {
			return repl, nil
		}
		return e, nil
	case core.ExprRational:
		return e, nil
	case core.ExprAdd:
		terms := make([]core.Expr, 0, len(e.Children()))
		for _, c := range e.Children() {
			s, err := substituteExpr(c, sym, repl)
			if err != nil {
				return core.Expr{}, err
			}
			terms = append(terms, s)
		}
		return core.NewAdd(terms...), nil
	case core.ExprMul:
		factors := make([]core.Expr, 0, len(e.Children()))
		for _, c := range e.Children() {
			s, err := substituteExpr(c, sym, repl)
			if err != nil {
				return core.Expr{}, err
			}
			factors = append(factors, s)
		}
		return core.NewMul(factors...), nil
	case core.ExprNeg:
		s, err := substituteExpr(e.Children()[0], sym, repl)
		if err != nil {
			return core.Expr{}, err
		}
		return core.NewNeg(s), nil
	case core.ExprPow:
		base, err := substituteExpr(e.Base(), sym, repl)
		if err != nil {
			return core.Expr{}, err
		}
		return core.NewPow(base, e.Exponent()), nil
	case core.ExprSqrt:
		s, err := substituteExpr(e.Children()[0], sym, repl)
		if err != nil {
			return core.Expr{}, err
		}
		return core.NewSqrt(s), nil
	case core.ExprCall:
		args := make([]core.Expr, 0, len(e.Children()))
		for _, c := range e.Children() {
			s, err := substituteExpr(c, sym, repl)
			if err != nil {
				return core.Expr{}, err
			}
			args = append(args, s)
		}
		return core.NewCall(e.FunctionID(), args...)
	case core.ExprRelation:
		lhs, err := substituteExpr(e.Left(), sym, repl)
		if err != nil {
			return core.Expr{}, err
		}
		rhs, err := substituteExpr(e.Right(), sym, repl)
		if err != nil {
			return core.Expr{}, err
		}
		return core.NewRelation(e.RelationOperator(), lhs, rhs), nil
	case core.ExprBranchSet:
		target, err := substituteExpr(e.BranchTarget(), sym, repl)
		if err != nil {
			return core.Expr{}, err
		}
		branches := make([]core.Expr, 0, len(e.Branches()))
		for _, c := range e.Branches() {
			s, err := substituteExpr(c, sym, repl)
			if err != nil {
				return core.Expr{}, err
			}
			branches = append(branches, s)
		}
		return core.NewBranchSet(target, branches...), nil
	default:
		return core.Expr{}, core.UnsupportedOperationError{
			Operation: "substitute",
			Reason:    "expression outside the closed node set",
		}
	}
}

// ---------------------------------------------------------------------------
// Bounded differentiation (§15.8)
// ---------------------------------------------------------------------------

func differentiateExpr(e core.Expr, wrt string) (core.Expr, error) {
	switch e.Kind() {
	case core.ExprRational:
		return core.NewRational(big.NewRat(0, 1)), nil
	case core.ExprSymbol:
		if e.SymbolName() == wrt {
			return core.NewRational(big.NewRat(1, 1)), nil
		}
		return core.NewRational(big.NewRat(0, 1)), nil
	case core.ExprAdd:
		terms := make([]core.Expr, 0, len(e.Children()))
		for _, c := range e.Children() {
			d, err := differentiateExpr(c, wrt)
			if err != nil {
				return core.Expr{}, err
			}
			terms = append(terms, d)
		}
		return core.NewAdd(terms...), nil
	case core.ExprMul:
		factors := e.Children()
		if len(factors) == 0 {
			return core.NewRational(big.NewRat(0, 1)), nil
		}
		terms := make([]core.Expr, 0, len(factors))
		for i, f := range factors {
			df, err := differentiateExpr(f, wrt)
			if err != nil {
				return core.Expr{}, err
			}
			piece := make([]core.Expr, 0, len(factors))
			piece = append(piece, factors[:i]...)
			piece = append(piece, df)
			piece = append(piece, factors[i+1:]...)
			terms = append(terms, core.NewMul(piece...))
		}
		return core.NewAdd(terms...), nil
	case core.ExprNeg:
		d, err := differentiateExpr(e.Children()[0], wrt)
		if err != nil {
			return core.Expr{}, err
		}
		return core.NewNeg(d), nil
	case core.ExprPow:
		base := e.Base()
		n := e.Exponent()
		if !n.IsInt() || n.Sign() < 0 {
			return core.Expr{}, core.UnsupportedOperationError{
				Operation: "differentiate",
				Reason:    "power rule limited to non-negative integer exponents",
			}
		}
		dbase, err := differentiateExpr(base, wrt)
		if err != nil {
			return core.Expr{}, err
		}
		if n.Sign() == 0 {
			// The coefficient n vanishes; the term is identically zero.
			return core.NewRational(big.NewRat(0, 1)), nil
		}
		if n.Cmp(big.NewRat(1, 1)) == 0 {
			// n·base^(n-1)·d(base) with n-1 = 0: the base^0 factor is a
			// formal identity of the rule, not an evaluation — the bound
			// Pow(x,0) gate must not leak into d(x)/dx = 1.
			return dbase, nil
		}
		pow := core.NewPow(base, new(big.Rat).Sub(n, big.NewRat(1, 1)))
		return core.NewMul(
			core.NewRational(n), pow, dbase), nil
	case core.ExprSqrt, core.ExprCall, core.ExprRelation, core.ExprBranchSet:
		return core.Expr{}, core.UnsupportedOperationError{
			Operation: "differentiate",
			Reason:    "node outside the bounded derivative rules",
		}
	default:
		return core.Expr{}, core.UnsupportedOperationError{
			Operation: "differentiate",
			Reason:    "expression outside the closed node set",
		}
	}
}

// ---------------------------------------------------------------------------
// Fixed-body expansion traversal (§15.9)
// ---------------------------------------------------------------------------

// expandCalls walks the whole tree and replaces every lorentz_factor Call by
// the fixed §15.9.1 body. There is no ID-match shortcut returning 1: the
// body tree is always traversed (G-Audit pins this by source review).
func expandCalls(e core.Expr) core.Expr {
	switch e.Kind() {
	case core.ExprCall:
		expanded := expandCall(e)
		if expanded.Kind() == core.ExprCall {
			// Unknown function id (unreachable in MVP) — expand args only.
			args := make([]core.Expr, 0, len(e.Children()))
			for _, c := range e.Children() {
				args = append(args, expandCalls(c))
			}
			out, err := core.NewCall(e.FunctionID(), args...)
			if err != nil {
				return e
			}
			return out
		}
		return expandCalls(expanded)
	case core.ExprAdd:
		terms := make([]core.Expr, 0, len(e.Children()))
		for _, c := range e.Children() {
			terms = append(terms, expandCalls(c))
		}
		return core.NewAdd(terms...)
	case core.ExprMul:
		factors := make([]core.Expr, 0, len(e.Children()))
		for _, c := range e.Children() {
			factors = append(factors, expandCalls(c))
		}
		return core.NewMul(factors...)
	case core.ExprNeg:
		return core.NewNeg(expandCalls(e.Children()[0]))
	case core.ExprPow:
		return core.NewPow(expandCalls(e.Base()), e.Exponent())
	case core.ExprSqrt:
		return core.NewSqrt(expandCalls(e.Children()[0]))
	case core.ExprRelation:
		return core.NewRelation(e.RelationOperator(),
			expandCalls(e.Left()), expandCalls(e.Right()))
	case core.ExprBranchSet:
		branches := make([]core.Expr, 0, len(e.Branches()))
		for _, c := range e.Branches() {
			branches = append(branches, expandCalls(c))
		}
		return core.NewBranchSet(expandCalls(e.BranchTarget()), branches...)
	default:
		return e
	}
}

// containsCall reports whether any Call node remains in the tree.
func containsCall(e core.Expr) bool {
	if e.Kind() == core.ExprCall {
		return true
	}
	for _, c := range e.Children() {
		if containsCall(c) {
			return true
		}
	}
	return false
}

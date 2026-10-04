package ops

import (
	"math/big"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

// ---------------------------------------------------------------------------
// Simplify engine (specs_v2_3.md §8.7, §9.6–§9.9; plan §6)
// ---------------------------------------------------------------------------

// Simplify implements specs_v2_3.md §15.6: purely mechanical, bottom-up
// recursive rewriting of the expression tree, followed by constructor
// re-normalization (flatten, exact rational combination, sign normal form,
// three-key sort). It never creates IDENTIFIED provenance, never changes
// corpus status, and consults no session state.
func Simplify(x core.Object) (core.Object, error) {
	if err := validInputs("simplify", x); err != nil {
		return core.Object{}, err
	}
	simplified, err := simplifyExpr(x.Expr(), x.Assumptions())
	if err != nil {
		return core.Object{}, err
	}
	return buildResult(x.Kind(), x.Dimension(), simplified, []core.Object{x}, nil)
}

func simplifyExpr(e core.Expr, A core.AssumptionSet) (core.Expr, error) {
	switch e.Kind() {
	case core.ExprSymbol, core.ExprRational:
		return e, nil

	case core.ExprAdd:
		terms := make([]core.Expr, 0, len(e.Children()))
		for _, c := range e.Children() {
			s, err := simplifyExpr(c, A)
			if err != nil {
				return core.Expr{}, err
			}
			terms = append(terms, s)
		}
		// Add() → 0; Add(x, 0) → x (the constructor keeps a combined
		// zero rational term for exactly this Simplify-time rule).
		filtered := make([]core.Expr, 0, len(terms))
		for _, t := range terms {
			if t.Kind() == core.ExprRational && t.RationalValue().Sign() == 0 {
				continue
			}
			filtered = append(filtered, t)
		}
		if len(filtered) == 0 {
			return core.NewRational(big.NewRat(0, 1)), nil
		}
		m := core.NewAdd(filtered...)
		if len(m.Children()) == 1 {
			return m.Children()[0], nil
		}
		return m, nil

	case core.ExprMul:
		factors := make([]core.Expr, 0, len(e.Children()))
		for _, c := range e.Children() {
			s, err := simplifyExpr(c, A)
			if err != nil {
				return core.Expr{}, err
			}
			factors = append(factors, s)
		}
		return simplifyMul(factors, A)

	case core.ExprNeg:
		s, err := simplifyExpr(e.Children()[0], A)
		if err != nil {
			return core.Expr{}, err
		}
		return core.NewNeg(s), nil

	case core.ExprPow:
		base, err := simplifyExpr(e.Base(), A)
		if err != nil {
			return core.Expr{}, err
		}
		return simplifyPow(base, e.Exponent(), A)

	case core.ExprSqrt:
		s, err := simplifyExpr(e.Children()[0], A)
		if err != nil {
			return core.Expr{}, err
		}
		return simplifySqrt(s, A)
	case core.ExprCall:
		args := make([]core.Expr, 0, len(e.Children()))
		for _, c := range e.Children() {
			s, err := simplifyExpr(c, A)
			if err != nil {
				return core.Expr{}, err
			}
			args = append(args, s)
		}
		return core.NewCall(e.FunctionID(), args...)

	case core.ExprRelation:
		lhs, err := simplifyExpr(e.Left(), A)
		if err != nil {
			return core.Expr{}, err
		}
		rhs, err := simplifyExpr(e.Right(), A)
		if err != nil {
			return core.Expr{}, err
		}
		return core.NewRelation(e.RelationOperator(), lhs, rhs), nil

	case core.ExprBranchSet:
		target, err := simplifyExpr(e.BranchTarget(), A)
		if err != nil {
			return core.Expr{}, err
		}
		branches := make([]core.Expr, 0, len(e.Branches()))
		for _, c := range e.Branches() {
			s, err := simplifyExpr(c, A)
			if err != nil {
				return core.Expr{}, err
			}
			branches = append(branches, s)
		}
		return core.NewBranchSet(target, branches...), nil

	default:
		return core.Expr{}, core.UnsupportedOperationError{
			Operation: "simplify",
			Reason:    "expression outside the closed node set",
		}
	}
}

// simplifyMul applies the §9.6 zero/one rules, the §9.7 repeated-power
// combinations, and re-normalizes through the constructor.
func simplifyMul(factors []core.Expr, A core.AssumptionSet) (core.Expr, error) {
	if len(factors) == 0 {
		return core.NewRational(big.NewRat(1, 1)), nil // Mul() → 1
	}
	// Mul(x, 0) → 0 when every other factor is finite (§9.6, plan §6).
	zeroSeen := false
	others := make([]core.Expr, 0, len(factors))
	for _, f := range factors {
		if f.Kind() == core.ExprRational && f.RationalValue().Sign() == 0 {
			zeroSeen = true
			continue
		}
		others = append(others, f)
	}
	if zeroSeen {
		allFinite := true
		for _, f := range others {
			if !finite(f, A) {
				allFinite = false
				break
			}
		}
		if allFinite {
			return core.NewRational(big.NewRat(0, 1)), nil
		}
		// Not finite: the zero factor stays; fall through with the full list.
		others = factors
	}
	// Mul(x, 1) → x: drop unit rational factors (§9.6).
	kept := make([]core.Expr, 0, len(others))
	for _, f := range others {
		if f.Kind() == core.ExprRational && f.RationalValue().Cmp(big.NewRat(1, 1)) == 0 {
			continue
		}
		kept = append(kept, f)
	}
	if len(kept) == 0 {
		return core.NewRational(big.NewRat(1, 1)), nil
	}
	if len(kept) == 1 {
		return kept[0], nil
	}
	// §9.7 repeated powers over identical non-rational factors.
	kept = collapseRepeatedPowers(kept)
	if len(kept) == 1 {
		return kept[0], nil
	}
	return core.NewMul(kept...), nil
}

// collapseRepeatedPowers groups structurally identical factors:
// x · x → Pow(x, n) for a non-Pow factor, and Pow(x, a) · … (n times) →
// Pow(x, a*n) when a is a non-negative integer (§9.7). Rationale-coefficients
// combine through the constructor, never through grouping. Factors outside
// those exact combinations are left untouched — no other power algebra.
func collapseRepeatedPowers(factors []core.Expr) []core.Expr {
	type group struct {
		e core.Expr
		n int
	}
	var groups []group
	for _, f := range factors {
		if f.Kind() == core.ExprRational {
			groups = append(groups, group{e: f, n: 1})
			continue
		}
		merged := false
		for i := range groups {
			if groups[i].e.Kind() != core.ExprRational &&
				core.EqualExpr(groups[i].e, f) {
				groups[i].n++
				merged = true
				break
			}
		}
		if !merged {
			groups = append(groups, group{e: f, n: 1})
		}
	}
	changed := false
	out := make([]core.Expr, 0, len(groups))
	for _, g := range groups {
		if g.n == 1 || g.e.Kind() == core.ExprRational {
			for i := 0; i < g.n; i++ {
				out = append(out, g.e)
			}
			continue
		}
		if g.e.Kind() == core.ExprPow {
			a := g.e.Exponent()
			if isNonNegInteger(a) {
				pow := core.NewPow(g.e.Base(),
					new(big.Rat).Mul(a, big.NewRat(int64(g.n), 1)))
				out = append(out, pow)
				changed = true
				continue
			}
			for i := 0; i < g.n; i++ {
				out = append(out, g.e)
			}
			continue
		}
		out = append(out, core.NewPow(g.e, big.NewRat(int64(g.n), 1)))
		changed = true
	}
	if !changed {
		return factors
	}
	return out
}

// simplifyPow applies the §8.7/§9.6/§9.7 exact rational-power rules.
func simplifyPow(base core.Expr, exponent *big.Rat, A core.AssumptionSet) (core.Expr, error) {
	// §9.7 nested powers: only when both exponents are non-negative integers.
	if base.Kind() == core.ExprPow && isNonNegInteger(base.Exponent()) &&
		isNonNegInteger(exponent) {
		product := new(big.Rat).Mul(base.Exponent(), exponent)
		return simplifyPow(base.Base(), product, A)
	}
	// Exact rational base rules (§9.6) — checked before the identity rules
	// so that 0^0 and 0^-1 are decided by the rational branch.
	if base.Kind() == core.ExprRational {
		q := base.RationalValue()
		if q.Sign() == 0 {
			switch {
			case exponent.Sign() == 0:
				return core.Expr{}, core.UnsupportedOperationError{
					Operation: "simplify",
					Reason:    "0^0 is unsupported",
				}
			case exponent.Sign() < 0:
				return core.Expr{}, core.UnsupportedOperationError{
					Operation: "simplify",
					Reason:    "0 raised to a negative exponent is unsupported",
				}
			default:
				return core.NewRational(big.NewRat(0, 1)), nil
			}
		}
		if q.Cmp(big.NewRat(1, 1)) == 0 {
			return core.NewRational(big.NewRat(1, 1)), nil // any exact exponent
		}
		if exponent.IsInt() {
			v := ratIntPow(q, exponent)
			if v == nil {
				return core.Expr{}, core.UnsupportedOperationError{
					Operation: "simplify",
					Reason:    "rational exponent out of range",
				}
			}
			return core.NewRational(v), nil
		}
		// Non-integer exponent of a general rational: keep the structural Pow.
		return core.NewPow(base, exponent), nil
	}
	// General base identity rules (§8.7, §9.6).
	if exponent.Cmp(big.NewRat(1, 1)) == 0 {
		return base, nil // Pow(x, 1) → x
	}
	if exponent.Sign() == 0 {
		// Pow(x, 0) → 1 only when x is nonzero-safe under current
		// assumptions; an unentailed base keeps its structure (F2).
		if A.EntailsNonZero(base) {
			return core.NewRational(big.NewRat(1, 1)), nil
		}
		return core.NewPow(base, exponent), nil
	}
	return core.NewPow(base, exponent), nil
}

// simplifySqrt applies the §9.8 rewrite rules, gated by bounded entailment.
func simplifySqrt(child core.Expr, A core.AssumptionSet) (core.Expr, error) {
	if child.Kind() == core.ExprRational {
		if root, ok := ratSqrt(child.RationalValue()); ok {
			return core.NewRational(root), nil
		}
		return core.NewSqrt(child), nil
	}
	// Sqrt(Pow(x, 2)) → x only when the assumptions entail x >= 0.
	if child.Kind() == core.ExprPow &&
		child.Exponent().Cmp(big.NewRat(2, 1)) == 0 &&
		A.EntailsNonNegative(child.Base()) {
		return child.Base(), nil
	}
	return core.NewSqrt(child), nil
}

// ---------------------------------------------------------------------------
// Finiteness predicate (plan §6; pins §9.6 "finite under current assumptions")
// ---------------------------------------------------------------------------

// finite reports whether e is a finite supported symbolic factor under A.
func finite(e core.Expr, A core.AssumptionSet) bool {
	switch e.Kind() {
	case core.ExprRational, core.ExprSymbol:
		return true
	case core.ExprAdd, core.ExprMul, core.ExprNeg:
		for _, c := range e.Children() {
			if !finite(c, A) {
				return false
			}
		}
		return true
	case core.ExprPow:
		b := e.Base()
		if !finite(b, A) {
			return false
		}
		n := e.Exponent()
		switch {
		case n.IsInt() && n.Sign() >= 0:
			return true
		case n.IsInt() && n.Sign() < 0:
			return A.EntailsNonZero(b)
		default:
			return false // fractional exponents: no MVP admissibility rule
		}
	case core.ExprSqrt:
		c := e.Children()[0]
		if c.Kind() == core.ExprRational &&
			(c.RationalValue().Sign() == 0 || c.RationalValue().Cmp(big.NewRat(1, 1)) == 0) {
			return finite(c, A)
		}
		return A.EntailsNonNegative(c) && finite(c, A)
	case core.ExprCall:
		return finite(expandCall(e), A)
	default:
		// Relation / BranchSet are not factors.
		return false
	}
}

// ---------------------------------------------------------------------------
// Fixed lorentz-factor body (§15.9.1) — expanded only by finiteness and Limit
// ---------------------------------------------------------------------------

// lorentzFactorBody builds the sole MVP function definition:
// 1 / Sqrt(1 - Pow(v/c, 2)). The Mul(1, …) form is retained until Simplify.
func lorentzFactorBody(v core.Expr) core.Expr {
	c, err := core.NewSymbol("c")
	if err != nil {
		panic("ops: fixed symbol c") // literal identifier, cannot fail
	}
	one := core.NewRational(big.NewRat(1, 1))
	ratio := core.NewMul(v, core.NewPow(c, big.NewRat(-1, 1)))
	squared := core.NewPow(ratio, big.NewRat(2, 1))
	denominator := core.NewSqrt(core.NewAdd(one, core.NewNeg(squared)))
	return core.NewMul(one, core.NewPow(denominator, big.NewRat(-1, 1)))
}

// expandCall expands the fixed MVP function body; other expressions pass
// through unchanged. There is no runtime function-definition registry.
func expandCall(e core.Expr) core.Expr {
	if e.Kind() == core.ExprCall &&
		e.FunctionID() == kernel.LorentzFactorFunctionID &&
		len(e.Children()) == 1 {
		return lorentzFactorBody(e.Children()[0])
	}
	return e
}

// ---------------------------------------------------------------------------
// Exact rational helpers
// ---------------------------------------------------------------------------

func isNonNegInteger(r *big.Rat) bool { return r.IsInt() && r.Sign() >= 0 }

// ratIntPow computes q^n for integer n within int64 range; nil when out of
// range (the caller keeps the structural Pow).
func ratIntPow(q *big.Rat, n *big.Rat) *big.Rat {
	if !n.IsInt() || !n.Num().IsInt64() {
		return nil
	}
	e := n.Num().Int64()
	negative := e < 0
	if negative {
		e = -e
	}
	num := new(big.Int).Exp(q.Num(), big.NewInt(e), nil)
	den := new(big.Int).Exp(q.Denom(), big.NewInt(e), nil)
	if negative {
		num, den = den, num
	}
	return new(big.Rat).SetFrac(num, den)
}

// ratSqrt returns the exact square root of a non-negative perfect-square
// rational.
func ratSqrt(q *big.Rat) (*big.Rat, bool) {
	if q.Sign() < 0 {
		return nil, false
	}
	sn := new(big.Int).Sqrt(q.Num())
	if new(big.Int).Mul(sn, sn).Cmp(q.Num()) != 0 {
		return nil, false
	}
	sd := new(big.Int).Sqrt(q.Denom())
	if new(big.Int).Mul(sd, sd).Cmp(q.Denom()) != 0 {
		return nil, false
	}
	return new(big.Rat).SetFrac(sn, sd), true
}

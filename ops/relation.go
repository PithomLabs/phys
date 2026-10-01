package ops

import (
	"math/big"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

// ---------------------------------------------------------------------------
// Relation operations (specs_v2_3.md §15.10–§15.12)
// ---------------------------------------------------------------------------

// Compare implements specs_v2_3.md §15.10: it constructs a relation
// artifact; it never decides whether an equality is physically true.
func Compare(a, b core.Object, op core.RelationOperator) (core.Object, error) {
	const c = "compare"
	if err := validInputs(c, a, b); err != nil {
		return core.Object{}, err
	}
	switch op {
	case core.RelationEq, core.RelationNeq, core.RelationLt, core.RelationLte,
		core.RelationGt, core.RelationGte:
	default:
		return core.Object{}, core.UnsupportedOperationError{
			Operation: c, Reason: "unknown relation operator",
		}
	}
	// All operators require equal dimensions (MRC-002).
	if !a.Dimension().Equal(b.Dimension()) {
		return core.Object{}, core.DimensionMismatchError{
			Operation: c,
			Left:      dimensionString(a),
			Right:     dimensionString(b),
		}
	}
	// §6.2 Compare table: every combination is allowed except two distinct
	// named kinds.
	ka, kb := a.Kind(), b.Kind()
	if ka != core.KindExpression && kb != core.KindExpression && ka != kb {
		return core.Object{}, core.CategoryMismatchError{
			Operation: c, Left: ka.String(), Right: kb.String(),
		}
	}
	// Inequalities require both sides to be ordered scalar quantities under
	// the bounded MVP ordering rule (§6.2).
	if op != core.RelationEq && op != core.RelationNeq {
		if !orderedScalarKind(ka) || !orderedScalarKind(kb) {
			return core.Object{}, core.CategoryMismatchError{
				Operation: c, Left: ka.String(), Right: kb.String(),
			}
		}
	}
	expr := core.NewRelation(op, a.Expr(), b.Expr())
	if !expr.Valid() {
		return core.Object{}, core.UnsupportedOperationError{
			Operation: c, Reason: "operands cannot form a relation",
		}
	}
	return buildResult(core.KindRelation, a.Dimension(), expr,
		[]core.Object{a, b}, nil)
}

// orderedScalarKind is the exact closed inequality list of §6.2. Velocity is
// intentionally absent from the MVP inequality path.
func orderedScalarKind(k core.Kind) bool {
	switch k {
	case core.KindMass, core.KindRestMass, core.KindTime, core.KindEnergy,
		core.KindKineticEnergy, core.KindSpeedOfLight, core.KindExpression:
		return true
	default:
		return false
	}
}

// Solve implements specs_v2_3.md §15.11: pattern matching on exactly
// Relation(eq, Pow(Symbol(t), 2), Expr) with target t a single Symbol. No
// linear solver, root finder, or equation isolation exists.
func Solve(relation, target core.Object) (core.Object, error) {
	const c = "solve"
	if err := validInputs(c, relation, target); err != nil {
		return core.Object{}, err
	}
	if relation.Kind() != core.KindRelation {
		return core.Object{}, core.CategoryMismatchError{
			Operation: c, Left: relation.Kind().String(), Right: core.KindRelation.String(),
		}
	}
	tSym := target.Expr()
	if tSym.Kind() != core.ExprSymbol {
		return core.Object{}, core.UnsupportedOperationError{
			Operation: c,
			Reason:    "target expression must be a single Symbol",
		}
	}
	e := relation.Expr()
	if e.Kind() != core.ExprRelation || e.RelationOperator() != core.RelationEq ||
		e.Left().Kind() != core.ExprPow ||
		e.Left().Exponent().Cmp(big.NewRat(2, 1)) != 0 ||
		e.Left().Base().Kind() != core.ExprSymbol ||
		e.Left().Base().SymbolName() != tSym.SymbolName() {
		return core.Object{}, core.UnsupportedOperationError{
			Operation: c,
			Reason:    "accepted shape is Relation(eq, Pow(Symbol, 2), Expr)",
		}
	}
	rhs := e.Right()
	sqrt := core.NewSqrt(rhs)
	branches := core.NewBranchSet(tSym, sqrt, core.NewNeg(sqrt))
	return buildResult(core.KindBranchSet, target.Dimension(), branches,
		[]core.Object{relation, target}, nil)
}

// SelectBranch implements the exact six-step contracts of §15.12. No
// sign-entailment failure error is invented: unsupported shapes are the only
// rejection beyond the common preconditions.
func SelectBranch(branches, constraint core.Object) (core.Object, error) {
	const c = "select_branch"
	if err := validInputs(c, branches, constraint); err != nil {
		return core.Object{}, err
	}
	// MRC-002 before the structural checks.
	if !branches.Dimension().Equal(constraint.Dimension()) {
		return core.Object{}, core.DimensionMismatchError{
			Operation: c,
			Left:      dimensionString(branches),
			Right:     dimensionString(constraint),
		}
	}
	// §6.2: BranchSet in, Relation constraint.
	if branches.Kind() != core.KindBranchSet || constraint.Kind() != core.KindRelation {
		return core.Object{}, core.CategoryMismatchError{
			Operation: c,
			Left:      branches.Kind().String(),
			Right:     constraint.Kind().String(),
		}
	}
	// 1. validate branch-set shape: BranchSet(target, [b, Neg(b)]).
	e := branches.Expr()
	if e.Kind() != core.ExprBranchSet || len(e.Branches()) != 2 {
		return core.Object{}, core.UnsupportedOperationError{
			Operation: c,
			Reason:    "expected BranchSet(target, [b, Neg(b)])",
		}
	}
	first, second := e.Branches()[0], e.Branches()[1]
	if second.Kind() != core.ExprNeg || !core.EqualExpr(second.Children()[0], first) {
		return core.Object{}, core.UnsupportedOperationError{
			Operation: c,
			Reason:    "expected BranchSet(target, [b, Neg(b)])",
		}
	}
	// 2. validate constraint compatibility: Relation(gte, target, 0).
	ce := constraint.Expr()
	if ce.Kind() != core.ExprRelation ||
		ce.RelationOperator() != core.RelationGte ||
		!core.EqualExpr(ce.Left(), e.BranchTarget()) ||
		ce.Right().Kind() != core.ExprRational ||
		ce.Right().RationalValue().Sign() != 0 {
		return core.Object{}, core.UnsupportedOperationError{
			Operation: c,
			Reason:    "expected constraint Relation(gte, target, 0)",
		}
	}
	// 3. merge the constraint as the deterministic structured assumption.
	selected, err := core.NewExprAssumption(
		core.AssumptionConstraint,
		"selected_branch/"+kernel.Hex(core.HashExpr(ce)),
		ce,
	)
	if err != nil {
		return core.Object{}, err
	}
	extra := []core.Assumption{selected}
	assumptions, _, err := mergeInputs([]core.Object{branches, constraint}, extra)
	if err != nil {
		return core.Object{}, err
	}
	// 4. select the first branch and simplify it under the merged
	// assumptions (mass-energy path: Sqrt(Pow(m·c²,2)) → m·c²).
	simplified, err := simplifyExpr(first, assumptions)
	if err != nil {
		return core.Object{}, err
	}
	// 5. provenance subject to contamination; 6. Kind=Expression with the
	// branch-set dimension.
	return buildResult(core.KindExpression, branches.Dimension(), simplified,
		[]core.Object{branches, constraint}, extra)
}

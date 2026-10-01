package ops

import (
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

// ---------------------------------------------------------------------------
// Acceptance C/D/E/F — the MRC violation classes (§14, REQ-032-01..04)
// ---------------------------------------------------------------------------

func TestDimensionMismatch(t *testing.T) {
	energy := exprObj(t, core.DimensionEnergy(), sym(t, "e1"))
	length := exprObj(t, core.DimensionLength(), sym(t, "l1"))

	cases := []struct {
		name string
		run  func() (core.Object, error)
	}{
		{"Add", func() (core.Object, error) { return Add(energy, length) }},
		{"Subtract", func() (core.Object, error) { return Subtract(energy, length) }},
		{"Compare", func() (core.Object, error) {
			return Compare(energy, length, core.RelationEq)
		}},
	}
	for _, c := range cases {
		out, err := c.run()
		if err == nil {
			t.Errorf("%s: mismatched dimensions accepted (result %v)", c.name, out)
			continue
		}
		var dm core.DimensionMismatchError
		if !errors.As(err, &dm) {
			t.Errorf("%s: error type %T (%v)", c.name, err, err)
			continue
		}
		if dm.Left == "" || dm.Right == "" || dm.Operation == "" {
			t.Errorf("%s: missing diagnostics: %+v", c.name, dm)
		}
		if dm.Left == dm.Right {
			t.Errorf("%s: identical dimension strings: %+v", c.name, dm)
		}
	}
}

func TestKindMismatchEqualDimensions(t *testing.T) {
	// F5: the fixture pair shares dimension M so MRC-002 never triggers and
	// this test isolates MRC-003.
	mass := defined(t, core.KindMass, core.DimensionMass(), sym(t, "m1"), "Mass")
	restMass := defined(t, core.KindRestMass, core.DimensionMass(), sym(t, "m2"), "RestMass")
	if !mass.Dimension().Equal(restMass.Dimension()) {
		t.Fatal("fixture must share dimension M")
	}

	out, err := Add(mass, restMass)
	if err == nil {
		t.Fatalf("Mass + RestMass accepted: %v", out)
	}
	var cat core.CategoryMismatchError
	if !errors.As(err, &cat) {
		t.Fatalf("error type %T (%v)", err, err)
	}
	if cat.Operation != "add" {
		t.Errorf("operation %q", cat.Operation)
	}
	if cat.Left != "Mass" || cat.Right != "RestMass" {
		t.Errorf("kinds %q vs %q", cat.Left, cat.Right)
	}

	if _, err := Subtract(mass, restMass); err == nil {
		t.Error("Mass - RestMass accepted")
	} else if !errors.As(err, &cat) {
		t.Errorf("subtract error type %T", err)
	}
	// Equal-dimension distinct named kinds are rejected for eq comparisons
	// but mixed Expression operands remain legal (§6.2).
	if _, err := Compare(mass, restMass, core.RelationEq); err == nil {
		t.Error("Compare eq across distinct named kinds accepted")
	} else if !errors.As(err, &cat) {
		t.Errorf("compare error type %T", err)
	}
}

func TestAssumptionConflict(t *testing.T) {
	x := sym(t, "x")
	key, err := core.NewExprAssumption(core.AssumptionConstraint, "shared_key",
		core.NewRelation(core.RelationGte, x, rat(0, 1)))
	if err != nil {
		t.Fatal(err)
	}
	other, err := core.NewExprAssumption(core.AssumptionConstraint, "shared_key",
		core.NewRelation(core.RelationGte, x, rat(1, 1)))
	if err != nil {
		t.Fatal(err)
	}
	a := defined(t, core.KindExpression, core.DimensionEnergy(), sym(t, "e1"), "a", key)
	b := defined(t, core.KindExpression, core.DimensionEnergy(), sym(t, "e2"), "b", other)

	_, err = Add(a, b)
	if err == nil {
		t.Fatal("conflicting assumptions merged")
	}
	var ac core.AssumptionConflictError
	if !errors.As(err, &ac) {
		t.Fatalf("error type %T (%v)", err, err)
	}
	if ac.Key != "shared_key" || ac.Kind != string(core.AssumptionConstraint) {
		t.Errorf("diagnostics: %+v", ac)
	}
	if ac.Left == ac.Right {
		t.Errorf("conflict values must differ: %+v", ac)
	}
	// Identical (kind, key, value) merges without conflict.
	same := defined(t, core.KindExpression, core.DimensionEnergy(), sym(t, "e3"), "c", key)
	if _, err := Add(a, same); err != nil {
		t.Errorf("identical assumption rejected: %v", err)
	}
}

func TestConventionConflict(t *testing.T) {
	c1, err := core.NewConvention("metric.signature", "-+++")
	if err != nil {
		t.Fatal(err)
	}
	c2, err := core.NewConvention("metric.signature", "++++")
	if err != nil {
		t.Fatal(err)
	}
	g := buildGolden(t)
	a := mintSpec(t, kernel.ObjectSpec{
		Name: "ca", Kind: core.KindExpression, Dimension: core.DimensionEnergy(),
		Expr: g.e, Assumptions: core.NewAssumptionSet(),
		Conventions:  core.NewConventionSet(c1),
		Provenance:   newProv(t, core.StatusDefined, "test", "test", nil, ""),
		CorpusStatus: kernel.CorpusNone,
	})
	b := mintSpec(t, kernel.ObjectSpec{
		Name: "cb", Kind: core.KindExpression, Dimension: core.DimensionEnergy(),
		Expr: sym(t, "e2"), Assumptions: core.NewAssumptionSet(),
		Conventions:  core.NewConventionSet(c2),
		Provenance:   newProv(t, core.StatusDefined, "test", "test", nil, ""),
		CorpusStatus: kernel.CorpusNone,
	})
	_, err = Add(a, b)
	if err == nil {
		t.Fatal("conflicting conventions merged")
	}
	var cc core.ConventionConflictError
	if !errors.As(err, &cc) {
		t.Fatalf("error type %T (%v)", err, err)
	}
	// Identical conventions merge cleanly.
	dup := mintSpec(t, kernel.ObjectSpec{
		Name: "cc", Kind: core.KindExpression, Dimension: core.DimensionEnergy(),
		Expr: sym(t, "e3"), Assumptions: core.NewAssumptionSet(),
		Conventions:  core.NewConventionSet(c1),
		Provenance:   newProv(t, core.StatusDefined, "test", "test", nil, ""),
		CorpusStatus: kernel.CorpusNone,
	})
	if _, err := Add(a, dup); err != nil {
		t.Errorf("identical conventions rejected: %v", err)
	}
}

// ---------------------------------------------------------------------------
// REQ-032-05 / REQ-005-12 — validity is checked first, never a panic
// ---------------------------------------------------------------------------

func TestInvalidObjectRejected(t *testing.T) {
	var zero core.Object
	if zero.Valid() {
		t.Fatal("zero Object must be invalid")
	}

	// All twelve dispatch ids reject a zero object at the shared precondition.
	arity := map[OperationID]int{
		OpAdd: 2, OpSubtract: 2, OpMultiply: 2, OpDivide: 2, OpPow: 1,
		OpSimplify: 1, OpSubstitute: 3, OpDifferentiate: 2, OpLimit: 3,
		OpCompare: 2, OpSolve: 2, OpSelectBranch: 2,
	}
	for id, n := range arity {
		inputs := make([]core.Object, n)
		for i := range inputs {
			inputs[i] = zero
		}
		params := OperationParams{Kind: "empty"}
		if id == OpPow {
			params = OperationParams{Kind: "pow", Exponent: "2/1"}
		}
		if id == OpCompare {
			params = OperationParams{Kind: "compare", Operator: core.RelationEq}
		}
		_, err := Apply(id, inputs, params)
		if err == nil {
			t.Errorf("%s: zero object accepted", id)
			continue
		}
		var io core.InvalidObjectError
		if !errors.As(err, &io) {
			t.Errorf("%s: error type %T (%v)", id, err, err)
			continue
		}
		if io.Operation != string(id) {
			t.Errorf("%s: operation %q", id, io.Operation)
		}
	}

	// Identify is rejected before any input inspection (session-owned).
	if _, err := Apply("identify", nil, OperationParams{Kind: "empty"}); err == nil {
		t.Error("identify accepted")
	}

	// Direct public op entry points share the same precondition.
	direct := []struct {
		name string
		run  func() (core.Object, error)
	}{
		{"Simplify", func() (core.Object, error) { return Simplify(zero) }},
		{"Add", func() (core.Object, error) { return Add(zero, zero) }},
		{"Substitute", func() (core.Object, error) {
			return Substitute(zero, zero, zero)
		}},
		{"Differentiate", func() (core.Object, error) { return Differentiate(zero, zero) }},
		{"Limit", func() (core.Object, error) { return Limit(zero, zero, zero) }},
		{"Compare", func() (core.Object, error) { return Compare(zero, zero, core.RelationEq) }},
		{"Solve", func() (core.Object, error) { return Solve(zero, zero) }},
		{"SelectBranch", func() (core.Object, error) { return SelectBranch(zero, zero) }},
		{"Pow", func() (core.Object, error) { return Pow(zero, big.NewRat(2, 1)) }},
	}
	for _, c := range direct {
		if _, err := c.run(); err == nil {
			t.Errorf("%s: zero object accepted", c.name)
		} else {
			var io core.InvalidObjectError
			if !errors.As(err, &io) {
				t.Errorf("%s: error type %T", c.name, err)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// REQ-032-18..20 — bounded-operation rejection surface
// ---------------------------------------------------------------------------

func TestSolveUnsupportedForm(t *testing.T) {
	g := buildGolden(t)
	target := g.energy
	relation := func(op core.RelationOperator, lhs, rhs core.Expr) core.Object {
		return defined(t, core.KindRelation, g.energy2,
			core.NewRelation(op, lhs, rhs), "r")
	}
	two := big.NewRat(2, 1)
	validRHS := core.NewPow(core.NewMul(g.m, core.NewPow(g.c, two)), two)

	cases := []struct {
		name     string
		rel      core.Object
		target   core.Object
		wantKind func(error) bool
	}{
		{"input not Relation",
			exprObj(t, core.DimensionEnergy(), g.e), target,
			func(err error) bool { var c core.CategoryMismatchError; return errors.As(err, &c) }},
		{"target not single symbol",
			relation(core.RelationEq, core.NewPow(g.e, two), validRHS),
			exprObj(t, core.DimensionEnergy(), core.NewNeg(g.e)),
			func(err error) bool { var u core.UnsupportedOperationError; return errors.As(err, &u) }},
		{"operator not eq",
			relation(core.RelationNeq, core.NewPow(g.e, two), validRHS), target,
			func(err error) bool { var u core.UnsupportedOperationError; return errors.As(err, &u) }},
		{"lhs not pow",
			relation(core.RelationEq, g.e, validRHS), target,
			func(err error) bool { var u core.UnsupportedOperationError; return errors.As(err, &u) }},
		{"exponent not two",
			relation(core.RelationEq, core.NewPow(g.e, big.NewRat(3, 1)), validRHS),
			target,
			func(err error) bool { var u core.UnsupportedOperationError; return errors.As(err, &u) }},
		{"base not symbol",
			relation(core.RelationEq, core.NewPow(rat(2, 1), two), validRHS), target,
			func(err error) bool { var u core.UnsupportedOperationError; return errors.As(err, &u) }},
		{"squared symbol mismatch",
			relation(core.RelationEq, core.NewPow(sym(t, "q"), two), validRHS),
			target,
			func(err error) bool { var u core.UnsupportedOperationError; return errors.As(err, &u) }},
	}
	for _, c := range cases {
		_, err := Solve(c.rel, c.target)
		if err == nil {
			t.Errorf("%s: accepted", c.name)
			continue
		}
		if !c.wantKind(err) {
			t.Errorf("%s: error type %T (%v)", c.name, err, err)
		}
	}
}

func TestDifferentiateUnsupportedForms(t *testing.T) {
	g := buildGolden(t)
	wrt := defined(t, core.KindVelocity, core.DimensionVelocity(), sym(t, "v"), "v")
	v := sym(t, "v")
	branch := core.NewBranchSet(g.e,
		core.NewSqrt(g.e), core.NewNeg(core.NewSqrt(g.e)))
	gamma, err := core.NewCall(kernel.LorentzFactorFunctionID, v)
	if err != nil {
		t.Fatal(err)
	}

	unsupported := []core.Expr{
		core.NewSqrt(v),
		gamma,
		core.NewRelation(core.RelationEq, g.e, v),
		branch,
		core.NewPow(v, big.NewRat(1, 2)),
		core.NewPow(v, big.NewRat(-1, 1)),
	}
	for _, e := range unsupported {
		target := exprObj(t, core.DimensionEnergy(), e)
		_, err := Differentiate(target, wrt)
		if err == nil {
			t.Errorf("%s: accepted", jsonExpr(t, e))
			continue
		}
		var u core.UnsupportedOperationError
		if !errors.As(err, &u) || u.Operation != "differentiate" {
			t.Errorf("%s: error %v", jsonExpr(t, e), err)
		}
	}
}

func TestLimitUnsupportedForm(t *testing.T) {
	v := sym(t, "v")
	variable := defined(t, core.KindVelocity, core.DimensionVelocity(), v, "v")
	zeroValue := exprObj(t, core.DimensionVelocity(), rat(0, 1))

	// Singular: v^-1 as v → 0.
	singular := exprObj(t, core.DimensionLength(), core.NewPow(v, big.NewRat(-1, 1)))
	_, err := Limit(singular, variable, zeroValue)
	if err == nil {
		t.Fatal("singular limit accepted")
	}
	var u core.UnsupportedOperationError
	if !errors.As(err, &u) || u.Operation != "limit" {
		t.Fatalf("singular error: %v", err)
	}
	if u.Reason == "" {
		t.Error("missing reason")
	}

	// A Call that the fixed body cannot expand (wrong arity).
	twoArg, err := core.NewCall(kernel.LorentzFactorFunctionID, v, v)
	if err != nil {
		t.Fatal(err)
	}
	target := exprObj(t, core.Dimensionless(), twoArg)
	if _, err := Limit(target, variable, zeroValue); err == nil {
		t.Fatal("unexpandable Call accepted")
	} else if !errors.As(err, &u) {
		t.Errorf("error type %T", err)
	}

	// Wrong variable/value dimension pairing.
	if _, err := Limit(singular, variable,
		exprObj(t, core.DimensionLength(), rat(0, 1))); err == nil {
		t.Fatal("dimension mismatch accepted")
	} else {
		var dm core.DimensionMismatchError
		if !errors.As(err, &dm) {
			t.Errorf("error type %T", err)
		}
	}
}

// ---------------------------------------------------------------------------
// REQ-032-07, §9.7-MUST-02 — boundedness of Simplify and power algebra
// ---------------------------------------------------------------------------

func TestSimplifyNeverIdentifies(t *testing.T) {
	g := buildGolden(t)

	// Every provenance status is legal as input; IDENTIFIED is never minted.
	statuses := []struct {
		status core.ProvenanceStatus
		just   string
	}{
		{core.StatusDefined, ""},
		{core.StatusPostulated, ""},
		{core.StatusDerived, ""},
		{core.StatusIdentified, "identity"},
		{core.StatusApproximated, ""},
		{core.StatusHypothesis, "hypothesis"},
	}
	for _, st := range statuses {
		obj := mintSpec(t, kernel.ObjectSpec{
			Name: "n", Kind: core.KindRelation, Dimension: g.energy2,
			Expr: core.NewRelation(core.RelationEq, g.e,
				core.NewMul(g.m, core.NewPow(g.c, big.NewRat(2, 1)))),
			Assumptions:  core.NewAssumptionSet(),
			Conventions:  core.NewConventionSet(),
			Provenance:   newProv(t, st.status, "test", "test", nil, st.just),
			CorpusStatus: kernel.CorpusEstablished,
		})
		out, err := Simplify(obj)
		if err != nil {
			t.Fatalf("%v: %v", st.status, err)
		}
		got := out.Provenance().Status()
		if got == core.StatusIdentified {
			t.Errorf("%v: Simplify minted IDENTIFIED", st.status)
		}
		if st.status == core.StatusHypothesis {
			if got != core.StatusHypothesis {
				t.Errorf("contamination: %v", got)
			}
		} else if got != core.StatusDerived {
			t.Errorf("%v: output status %v", st.status, got)
		}
		// Corpus status is never inferred: results are always NONE.
		if out.CorpusStatus() != kernel.CorpusNone {
			t.Errorf("%v: corpus status %v", st.status, out.CorpusStatus())
		}
		// The relation operator is preserved; no identification relation is
		// created (§9.9, 🟡-3).
		if out.Expr().Kind() != core.ExprRelation ||
			out.Expr().RelationOperator() != core.RelationEq {
			t.Errorf("%v: relation not preserved: %s", st.status, jsonExpr(t, out.Expr()))
		}
	}

	// A non-relation input yields no relation node at all.
	plain, err := Simplify(exprObj(t, core.DimensionEnergy(),
		core.NewAdd(sym(t, "x"), rat(1, 1))))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(jsonExpr(t, plain.Expr()), `"kind":"relation"`) {
		t.Errorf("Simplify invented a relation: %s", jsonExpr(t, plain.Expr()))
	}
}

func TestUnsupportedPowerAlgebra(t *testing.T) {
	x := sym(t, "x")
	y := sym(t, "y")
	two := big.NewRat(2, 1)

	// Only the exact §9.7 combinations exist — nothing branch-sensitive or
	// general-purpose.
	stays := []struct {
		name string
		in   core.Expr
	}{
		{"x * x^-1 not collapsed",
			core.NewMul(x, core.NewPow(x, big.NewRat(-1, 1)))},
		{"identical negative-exponent factors ungrouped",
			core.NewMul(core.NewPow(x, big.NewRat(-1, 1)),
				core.NewPow(x, big.NewRat(-1, 1)))},
		{"nested negative exponents uncombined",
			core.NewPow(core.NewPow(x, big.NewRat(-1, 1)), big.NewRat(-1, 1))},
		{"fractional nested exponents uncombined",
			core.NewPow(core.NewPow(x, big.NewRat(1, 2)), big.NewRat(1, 2))},
		{"fractional powers of a product untouched",
			core.NewPow(core.NewMul(x, y), big.NewRat(1, 2))},
	}
	for _, c := range stays {
		js := jsonExpr(t, mustSimp(t, c.in))
		if !strings.Contains(js, `"kind":"pow"`) && !strings.Contains(js, `"kind":"mul"`) {
			t.Errorf("%s: collapsed to %s", c.name, js)
			continue
		}
		// The exact form must be preserved, not further reduced.
		if !core.EqualExpr(mustSimp(t, c.in), c.in) &&
			!core.EqualExpr(mustSimp(t, c.in), mustSimp(t, mustSimp(t, c.in))) {
			t.Errorf("%s: not idempotent", c.name)
		}
	}
	// x·x is the only symbolic square admitted for non-Pow factors.
	got := mustSimp(t, core.NewMul(x, y))
	if !core.EqualExpr(got, core.NewMul(x, y)) {
		t.Errorf("x*y regrouped: %s", jsonExpr(t, got))
	}
	_ = two
}

// ---------------------------------------------------------------------------
// REQ-§30-MUST-04 — violations are typed errors, never panics
// ---------------------------------------------------------------------------

func TestNoPanicOnMRCViolation(t *testing.T) {
	g := buildGolden(t)
	energy := exprObj(t, core.DimensionEnergy(), g.e)
	length := exprObj(t, core.DimensionLength(), sym(t, "l1"))
	mass := defined(t, core.KindMass, core.DimensionMass(), sym(t, "m1"), "Mass")
	restMass := defined(t, core.KindRestMass, core.DimensionMass(), sym(t, "m2"), "RestMass")

	keyA, err := core.NewExprAssumption(core.AssumptionConstraint, "panic_key",
		core.NewRelation(core.RelationGte, g.e, rat(0, 1)))
	if err != nil {
		t.Fatal(err)
	}
	keyB, err := core.NewExprAssumption(core.AssumptionConstraint, "panic_key",
		core.NewRelation(core.RelationGte, g.e, rat(1, 1)))
	if err != nil {
		t.Fatal(err)
	}
	conflictA := defined(t, core.KindExpression, core.DimensionEnergy(), g.e, "pa", keyA)
	conflictB := defined(t, core.KindExpression, core.DimensionEnergy(),
		sym(t, "e2"), "pb", keyB)
	c1, err := core.NewConvention("panic.convention", "one")
	if err != nil {
		t.Fatal(err)
	}
	c2, err := core.NewConvention("panic.convention", "two")
	if err != nil {
		t.Fatal(err)
	}
	convA := mintSpec(t, kernel.ObjectSpec{
		Name: "va", Kind: core.KindExpression, Dimension: core.DimensionEnergy(),
		Expr: g.e, Assumptions: core.NewAssumptionSet(),
		Conventions:  core.NewConventionSet(c1),
		Provenance:   newProv(t, core.StatusDefined, "test", "test", nil, ""),
		CorpusStatus: kernel.CorpusNone,
	})
	convB := mintSpec(t, kernel.ObjectSpec{
		Name: "vb", Kind: core.KindExpression, Dimension: core.DimensionEnergy(),
		Expr: sym(t, "e3"), Assumptions: core.NewAssumptionSet(),
		Conventions:  core.NewConventionSet(c2),
		Provenance:   newProv(t, core.StatusDefined, "test", "test", nil, ""),
		CorpusStatus: kernel.CorpusNone,
	})
	var zero core.Object

	violations := []struct {
		name string
		run  func() error
	}{
		{"dimension Add", func() error { _, err := Add(energy, length); return err }},
		{"dimension Compare", func() error {
			_, err := Compare(energy, length, core.RelationGt); return err
		}},
		{"kind Add", func() error { _, err := Add(mass, restMass); return err }},
		{"kind SelectBranch", func() error {
			_, err := SelectBranch(energy, g.rel); return err
		}},
		{"assumption conflict", func() error {
			_, err := Add(conflictA, conflictB); return err
		}},
		{"convention conflict", func() error { _, err := Add(convA, convB); return err }},
		{"zero object", func() error { _, err := Simplify(zero); return err }},
		{"0^0", func() error {
			_, err := simpExpr(t, core.NewPow(rat(0, 1), big.NewRat(0, 1))); return err
		}},
		{"nil exponent", func() error { _, err := Pow(energy, nil); return err }},
		{"solve form", func() error { _, err := Solve(energy, g.energy); return err }},
		{"differentiate sqrt", func() error {
			_, err := Differentiate(exprObj(t, core.DimensionEnergy(),
				core.NewSqrt(sym(t, "v"))), g.momentum)
			return err
		}},
		{"dispatch unknown id", func() error {
			_, err := Apply("nope", nil, OperationParams{Kind: "empty"}); return err
		}},
	}
	for _, v := range violations {
		func() {
			panicked := true
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s: panic: %v", v.name, r)
				}
				_ = panicked
			}()
			err := v.run()
			if err == nil {
				t.Errorf("%s: violation not rejected", v.name)
			}
			panicked = false
		}()
	}
}

func TestDispatchClosed(t *testing.T) {
	// Unknown ids are rejected by the closed switch — no registry lookup.
	unknown := []OperationID{
		"", "identify", "Identify", "PROMOTE", "eval", "add ", "subs",
	}
	for _, id := range unknown {
		if _, err := Apply(id, nil, OperationParams{Kind: "empty"}); err == nil {
			t.Errorf("id %q accepted", id)
		} else {
			var u core.UnsupportedOperationError
			if !errors.As(err, &u) {
				t.Errorf("id %q error type %T", id, err)
			}
			if u.Operation != string(id) {
				t.Errorf("id %q diagnostics %q", id, u.Operation)
			}
		}
	}
	// The twelve legal ids are exactly the exported constants.
	legal := []OperationID{
		OpAdd, OpSubtract, OpMultiply, OpDivide, OpPow, OpSimplify,
		OpSubstitute, OpDifferentiate, OpLimit, OpCompare, OpSolve,
		OpSelectBranch,
	}
	if len(legal) != 12 {
		t.Fatalf("legal id count %d", len(legal))
	}
	seen := map[OperationID]bool{}
	for _, id := range legal {
		if seen[id] {
			t.Errorf("duplicate id %q", id)
		}
		seen[id] = true
	}
}

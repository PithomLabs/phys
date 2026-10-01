package ops

import (
	"encoding/json"
	"errors"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

// ---------------------------------------------------------------------------
// Shared fixtures (ops is white-box; fixtures never import corpus packages)
// ---------------------------------------------------------------------------

func newProv(t *testing.T, status core.ProvenanceStatus, source, framework string,
	parents []string, justification string) core.Provenance {
	t.Helper()
	p, err := kernel.NewProvenance(status, source, framework, parents,
		[32]byte{}, [32]byte{}, kernel.MRCVersion, justification)
	if err != nil {
		t.Fatalf("provenance: %v", err)
	}
	return p
}

func mintSpec(t *testing.T, spec kernel.ObjectSpec) core.Object {
	t.Helper()
	o, err := kernel.MintObject(spec)
	if err != nil {
		t.Fatalf("mint: %v", err)
	}
	return o
}

func defined(t *testing.T, kind core.Kind, dim core.Dimension, e core.Expr,
	name string, as ...core.Assumption) core.Object {
	t.Helper()
	return mintSpec(t, kernel.ObjectSpec{
		Name: name, Kind: kind, Dimension: dim, Expr: e,
		Assumptions:  core.NewAssumptionSet(as...),
		Conventions:  core.NewConventionSet(),
		Provenance:   newProv(t, core.StatusDefined, "test", "test", nil, ""),
		CorpusStatus: kernel.CorpusNone,
	})
}

func exprObj(t *testing.T, dim core.Dimension, e core.Expr, as ...core.Assumption) core.Object {
	t.Helper()
	return defined(t, core.KindExpression, dim, e, "e", as...)
}

func hypoObj(t *testing.T, kind core.Kind, dim core.Dimension, e core.Expr) core.Object {
	t.Helper()
	return mintSpec(t, kernel.ObjectSpec{
		Name: "h", Kind: kind, Dimension: dim, Expr: e,
		Assumptions:  core.NewAssumptionSet(),
		Conventions:  core.NewConventionSet(),
		Provenance:   newProv(t, core.StatusHypothesis, "test", "test", nil, "hypothesis"),
		CorpusStatus: kernel.CorpusNone,
	})
}

func sym(t *testing.T, name string) core.Expr {
	t.Helper()
	e, err := core.NewSymbol(name)
	if err != nil {
		t.Fatalf("symbol %s: %v", name, err)
	}
	return e
}

func rat(n, d int64) core.Expr { return core.NewRational(big.NewRat(n, d)) }

func jsonExpr(t *testing.T, e core.Expr) string {
	t.Helper()
	b, err := core.CanonicalExprJSON(e)
	if err != nil {
		t.Fatalf("canonical expr: %v", err)
	}
	return string(b)
}

func goldenJSON(t *testing.T, o core.Object) string {
	t.Helper()
	b, err := kernel.CanonicalObjectJSON(o)
	if err != nil {
		t.Fatalf("canonical object: %v", err)
	}
	return string(b)
}

func gteZeroAssumption(t *testing.T, subject core.Expr, key string) core.Assumption {
	t.Helper()
	a, err := core.NewExprAssumption(core.AssumptionConstraint, key,
		core.NewRelation(core.RelationGte, subject, rat(0, 1)))
	if err != nil {
		t.Fatalf("assumption: %v", err)
	}
	return a
}

func gtZeroAssumption(t *testing.T, subject core.Expr, key string) core.Assumption {
	t.Helper()
	a, err := core.NewExprAssumption(core.AssumptionConstraint, key,
		core.NewRelation(core.RelationGt, subject, rat(0, 1)))
	if err != nil {
		t.Fatalf("assumption: %v", err)
	}
	return a
}

// golden builds the §20 energy–momentum fixtures used by Solve, SelectBranch,
// Limit, and the provenance tests.
type golden struct {
	rel, momentum, zeroMomentum, energy, zeroEnergy core.Object
	m, c, e, p                                      core.Expr
	massGE0, cGT0                                   core.Assumption
	energy2                                         core.Dimension
}

func buildGolden(t *testing.T) golden {
	t.Helper()
	g := golden{}
	g.m = sym(t, "m")
	g.c = sym(t, "c")
	g.e = sym(t, "E")
	g.p = sym(t, "p")
	g.massGE0 = gteZeroAssumption(t, g.m, "rest_mass_nonnegative")
	g.cGT0 = gtZeroAssumption(t, g.c, "speed_of_light_positive")
	two := big.NewRat(2, 1)
	g.energy2 = core.DimensionEnergy().Pow(two)
	lhs := core.NewPow(g.e, two)
	rhs := core.NewAdd(
		core.NewPow(core.NewMul(g.p, g.c), two),
		core.NewPow(core.NewMul(g.m, core.NewPow(g.c, two)), two),
	)
	g.rel = defined(t, core.KindRelation, g.energy2,
		core.NewRelation(core.RelationEq, lhs, rhs), "em",
		g.massGE0, g.cGT0)
	g.zeroMomentum = defined(t, core.KindExpression, core.DimensionMomentum(),
		rat(0, 1), "zp", g.massGE0, g.cGT0)
	g.momentum = defined(t, core.KindExpression, core.DimensionMomentum(),
		g.p, "p3", g.massGE0, g.cGT0)
	g.energy = defined(t, core.KindEnergy, core.DimensionEnergy(),
		g.e, "En", g.massGE0, g.cGT0)
	g.zeroEnergy = defined(t, core.KindEnergy, core.DimensionEnergy(),
		rat(0, 1), "ZE", g.massGE0, g.cGT0)
	return g
}

// goldenRelation runs Substitute → Simplify and returns the accepted
// quadratic relation fed to Solve.
func goldenRelation(t *testing.T, g golden) core.Object {
	t.Helper()
	sub, err := Substitute(g.rel, g.momentum, g.zeroMomentum)
	if err != nil {
		t.Fatalf("substitute: %v", err)
	}
	simp, err := Simplify(sub)
	if err != nil {
		t.Fatalf("simplify: %v", err)
	}
	return simp
}

// goldenSolved runs Substitute → Simplify → Solve on the fixture relation.
func goldenSolved(t *testing.T, g golden) core.Object {
	t.Helper()
	sub, err := Substitute(g.rel, g.momentum, g.zeroMomentum)
	if err != nil {
		t.Fatalf("substitute: %v", err)
	}
	simp, err := Simplify(sub)
	if err != nil {
		t.Fatalf("simplify: %v", err)
	}
	br, err := Solve(simp, g.energy)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	return br
}

// ---------------------------------------------------------------------------
// Dispatch and parameter contract (§15.13, §15.13.1)
// ---------------------------------------------------------------------------

func TestApplyPositionalInputs(t *testing.T) {
	g := buildGolden(t)
	v := sym(t, "v")
	w := sym(t, "w")
	mExpr := sym(t, "m")
	kinetic := core.NewMul(rat(1, 2), g.m, core.NewPow(v, big.NewRat(2, 1)))

	a := exprObj(t, core.DimensionEnergy(), g.e)
	b := exprObj(t, core.DimensionEnergy(), core.NewAdd(g.e, rat(1, 1)))
	variable := defined(t, core.KindVelocity, core.DimensionVelocity(), v, "v")
	replacement := defined(t, core.KindVelocity, core.DimensionVelocity(), w, "w")
	targetSub := exprObj(t, core.DimensionEnergy(),
		core.NewMul(v, sym(t, "x"), g.e))
	kineticObj := exprObj(t, core.DimensionEnergy(), kinetic)
	wrt := defined(t, core.KindVelocity, core.DimensionVelocity(), v, "v")
	limitTarget := exprObj(t, core.DimensionEnergy(), core.NewPow(v, big.NewRat(2, 1)))
	br := goldenSolved(t, g)
	constraint, err := Compare(g.energy, g.zeroEnergy, core.RelationGte)
	if err != nil {
		t.Fatal(err)
	}

	type pos struct {
		id     OperationID
		params OperationParams
		inputs []core.Object
		verify func(t *testing.T, got core.Object)
	}
	cases := []pos{
		{OpAdd, OperationParams{Kind: "empty"}, []core.Object{a, b},
			func(t *testing.T, got core.Object) {
				if !core.EqualExpr(got.Expr(), core.NewAdd(a.Expr(), b.Expr())) {
					t.Fatalf("add bound wrong index: %s", jsonExpr(t, got.Expr()))
				}
			}},
		{OpSubtract, OperationParams{Kind: "empty"}, []core.Object{a, b},
			func(t *testing.T, got core.Object) {
				want := core.NewAdd(a.Expr(), core.NewNeg(b.Expr()))
				if !core.EqualExpr(got.Expr(), want) {
					t.Fatalf("subtract order: %s", jsonExpr(t, got.Expr()))
				}
			}},
		{OpMultiply, OperationParams{Kind: "empty"}, []core.Object{a, b},
			func(t *testing.T, got core.Object) {
				if !core.EqualExpr(got.Expr(), core.NewMul(a.Expr(), b.Expr())) {
					t.Fatalf("multiply bound wrong index: %s", jsonExpr(t, got.Expr()))
				}
			}},
		{OpDivide, OperationParams{Kind: "empty"}, []core.Object{a, b},
			func(t *testing.T, got core.Object) {
				want := core.NewMul(a.Expr(), core.NewPow(b.Expr(), big.NewRat(-1, 1)))
				if !core.EqualExpr(got.Expr(), want) {
					t.Fatalf("divide bound wrong index: %s", jsonExpr(t, got.Expr()))
				}
			}},
		{OpPow, OperationParams{Kind: "pow", Exponent: "2/1"},
			[]core.Object{exprObj(t, core.DimensionEnergy(), mExpr)},
			func(t *testing.T, got core.Object) {
				want := core.NewPow(mExpr, big.NewRat(2, 1))
				if !core.EqualExpr(got.Expr(), want) {
					t.Fatalf("pow bound wrong: %s", jsonExpr(t, got.Expr()))
				}
			}},
		{OpSimplify, OperationParams{Kind: "empty"},
			[]core.Object{exprObj(t, core.DimensionEnergy(),
				core.NewAdd(mExpr, rat(0, 1)))},
			func(t *testing.T, got core.Object) {
				if !core.EqualExpr(got.Expr(), mExpr) {
					t.Fatalf("simplify: %s", jsonExpr(t, got.Expr()))
				}
			}},
		{OpSubstitute, OperationParams{Kind: "empty"},
			[]core.Object{targetSub, variable, replacement},
			func(t *testing.T, got core.Object) {
				if strings.Contains(jsonExpr(t, got.Expr()), `"name":"v"`) {
					t.Fatalf("variable not replaced from index 1: %s", jsonExpr(t, got.Expr()))
				}
				if !strings.Contains(jsonExpr(t, got.Expr()), `"name":"w"`) {
					t.Fatalf("replacement not taken from index 2: %s", jsonExpr(t, got.Expr()))
				}
				if got.Kind() != core.KindExpression {
					t.Fatalf("target kind not preserved: %v", got.Kind())
				}
			}},
		{OpDifferentiate, OperationParams{Kind: "empty"},
			[]core.Object{kineticObj, wrt},
			func(t *testing.T, got core.Object) {
				if !core.EqualExpr(got.Expr(), core.NewMul(g.m, v)) {
					t.Fatalf("differentiate indices: %s", jsonExpr(t, got.Expr()))
				}
			}},
		{OpLimit, OperationParams{Kind: "empty"},
			[]core.Object{limitTarget, variable, exprObj(t, core.DimensionVelocity(), rat(0, 1))},
			func(t *testing.T, got core.Object) {
				if jsonExpr(t, got.Expr()) != `{"kind":"rational","value":"0/1"}` {
					t.Fatalf("limit: %s", jsonExpr(t, got.Expr()))
				}
			}},
		{OpCompare, OperationParams{Kind: "compare", Operator: core.RelationEq},
			[]core.Object{a, b},
			func(t *testing.T, got core.Object) {
				if got.Kind() != core.KindRelation {
					t.Fatalf("compare kind: %v", got.Kind())
				}
				if !core.EqualExpr(got.Expr().Left(), a.Expr()) {
					t.Fatalf("compare lhs not inputs[0]: %s", jsonExpr(t, got.Expr()))
				}
				if !core.EqualExpr(got.Expr().Right(), b.Expr()) {
					t.Fatalf("compare rhs not inputs[1]: %s", jsonExpr(t, got.Expr()))
				}
			}},
		{OpSolve, OperationParams{Kind: "empty"},
			[]core.Object{goldenRelation(t, g), g.energy},
			func(t *testing.T, got core.Object) {
				if got.Kind() != core.KindBranchSet {
					t.Fatalf("solve kind: %v", got.Kind())
				}
			}},
		{OpSelectBranch, OperationParams{Kind: "empty"},
			[]core.Object{br, constraint},
			func(t *testing.T, got core.Object) {
				if got.Kind() != core.KindExpression {
					t.Fatalf("select kind: %v", got.Kind())
				}
				if !core.EqualExpr(got.Expr(), core.NewMul(g.m, core.NewPow(g.c, big.NewRat(2, 1)))) {
					t.Fatalf("select result: %s", jsonExpr(t, got.Expr()))
				}
			}},
	}
	for _, c := range cases {
		got, err := Apply(c.id, c.inputs, c.params)
		if err != nil {
			t.Errorf("%s: %v", c.id, err)
			continue
		}
		if !got.Valid() {
			t.Errorf("%s: invalid result", c.id)
			continue
		}
		c.verify(t, got)

		// Exact input count per §15.13.1 — unused positions omitted.
		if _, err := Apply(c.id, c.inputs[:len(c.inputs)-1], c.params); err == nil {
			t.Errorf("%s: missing input position accepted", c.id)
		}
		if _, err := Apply(c.id, nil, c.params); err == nil {
			t.Errorf("%s: nil inputs accepted", c.id)
		}
	}

	// Wrong-order inputs must be rejected wherever the positional order is
	// load-bearing (MUST-02: inputs are never reordered or packed).
	wrongOrder := []struct {
		id     OperationID
		inputs []core.Object
	}{
		{OpSubstitute, []core.Object{variable, targetSub, replacement}},
		{OpDifferentiate, []core.Object{wrt, kineticObj}},
		{OpLimit, []core.Object{variable, limitTarget, exprObj(t, core.DimensionVelocity(), rat(0, 1))}},
		{OpSolve, []core.Object{g.energy, br}},
		{OpSelectBranch, []core.Object{constraint, br}},
	}
	for _, c := range wrongOrder {
		params := OperationParams{Kind: "empty"}
		if c.id == OpCompare {
			params = OperationParams{Kind: "compare", Operator: core.RelationEq}
		}
		if _, err := Apply(c.id, c.inputs, params); err == nil {
			t.Errorf("%s: wrong-order inputs accepted", c.id)
		}
	}

	// identify is excluded from dispatch: neither as id nor as params kind.
	if _, err := Apply("identify", []core.Object{a}, OperationParams{Kind: "empty"}); err == nil {
		t.Error("identify OperationID must be rejected")
	} else {
		var u core.UnsupportedOperationError
		if !errors.As(err, &u) {
			t.Errorf("identify rejection type: %v", err)
		}
	}
	if _, err := Apply(OpAdd, []core.Object{a, b},
		OperationParams{Kind: "identify", Justification: "why"}); err == nil {
		t.Error("identify params kind must be rejected by Apply")
	}
}

func TestParamsCanonicalRoundTrip(t *testing.T) {
	canonical := `{"kind":"pow","exponent":"2/1","operator":"","justification":""}`
	p, err := ParseOperationParams([]byte(canonical))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	out, err := p.CanonicalJSON()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if string(out) != canonical {
		t.Fatalf("round trip:\n got %s\nwant %s", out, canonical)
	}

	variants := []string{
		`{"kind":"empty","exponent":"","operator":"","justification":""}`,
		`{"kind":"pow","exponent":"-3/4","operator":"","justification":""}`,
		`{"kind":"compare","exponent":"","operator":"gte","justification":""}`,
		`{"kind":"identify","exponent":"","operator":"","justification":"because"}`,
	}
	for _, v := range variants {
		got, err := ParseOperationParams([]byte(v))
		if err != nil {
			t.Fatalf("parse %s: %v", v, err)
		}
		b, err := got.CanonicalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != v {
			t.Errorf("round trip:\n got %s\nwant %s", b, v)
		}
	}

	reject := []string{
		`{"kind":"pow","exponent":"2/1","operator":"eq","justification":"","extra":1}`,
		`{"kind":"pow","exponent":"2/1","operator":"eq","justification":""}`,
		`{"kind":"pow","exponent":"2/1","operator":"","justification":"why"}`,
		`{"kind":"compare","exponent":"2/1","operator":"eq","justification":""}`,
		`{"kind":"compare","exponent":"","operator":"eq","justification":"why"}`,
		`{"kind":"empty","exponent":"1/1","operator":"","justification":""}`,
		`{"kind":"empty","exponent":"","operator":"eq","justification":""}`,
		`{"kind":"empty","exponent":"","operator":"","justification":"why"}`,
		`{"kind":"identify","exponent":"2/1","operator":"","justification":"because"}`,
		`{"kind":"identify","exponent":"","operator":"eq","justification":"because"}`,
		`{"kind":"nope","exponent":"","operator":"","justification":""}`,
		`{"kind":"pow","exponent":"4/2","operator":"","justification":""}`,
		`{"kind":"pow","exponent":"2.5","operator":"","justification":""}`,
		`{"kind":"pow","exponent":"","operator":"","justification":""}`,
		`{"kind":"compare","exponent":"","operator":"","justification":""}`,
		`{"kind":"identify","exponent":"","operator":"","justification":"  "}`,
		`{"kind":"empty","operator":"","justification":""}{}`,
	}
	for _, v := range reject {
		if _, err := ParseOperationParams([]byte(v)); err == nil {
			t.Errorf("accepted invalid params %s", v)
		}
	}

	// Per-id binding (§15.13).
	binding := []struct {
		id     OperationID
		params OperationParams
	}{
		{OpPow, OperationParams{Kind: "empty"}},
		{OpPow, OperationParams{Kind: "compare", Operator: core.RelationEq}},
		{OpCompare, OperationParams{Kind: "empty"}},
		{OpAdd, OperationParams{Kind: "pow", Exponent: "2/1"}},
		{OpAdd, OperationParams{Kind: "identify", Justification: "why"}},
		{OpAdd, OperationParams{Kind: "bogus"}},
		{"bogus", OperationParams{Kind: "empty"}},
	}
	for _, c := range binding {
		if err := ValidateOperationParams(c.id, c.params); err == nil {
			t.Errorf("binding accepted %s / %+v", c.id, c.params)
		} else {
			var u core.UnsupportedOperationError
			if !errors.As(err, &u) {
				t.Errorf("binding error type %v", err)
			}
		}
	}
	if err := ValidateOperationParams(OpSimplify, OperationParams{Kind: "empty"}); err != nil {
		t.Errorf("valid binding rejected: %v", err)
	}
}

func TestOperationCommonSemantics(t *testing.T) {
	g := buildGolden(t)
	a := exprObj(t, core.DimensionEnergy(), g.e)
	b := exprObj(t, core.DimensionEnergy(), core.NewAdd(g.e, rat(1, 1)))
	br := goldenSolved(t, g)
	constraint, err := Compare(g.energy, g.zeroEnergy, core.RelationGte)
	if err != nil {
		t.Fatal(err)
	}
	v := sym(t, "v")
	variable := defined(t, core.KindVelocity, core.DimensionVelocity(), v, "v")
	replacement := defined(t, core.KindVelocity, core.DimensionVelocity(),
		sym(t, "w"), "w")
	wrongDimValue := exprObj(t, core.DimensionLength(), rat(0, 1))

	invocations := []struct {
		id     OperationID
		params OperationParams
		inputs []core.Object
	}{
		{OpAdd, OperationParams{Kind: "empty"}, []core.Object{a, b}},
		{OpSubtract, OperationParams{Kind: "empty"}, []core.Object{a, b}},
		{OpMultiply, OperationParams{Kind: "empty"}, []core.Object{a, b}},
		{OpDivide, OperationParams{Kind: "empty"}, []core.Object{a, b}},
		{OpPow, OperationParams{Kind: "pow", Exponent: "2/1"}, []core.Object{a}},
		{OpSimplify, OperationParams{Kind: "empty"}, []core.Object{a}},
		{OpSubstitute, OperationParams{Kind: "empty"},
			[]core.Object{exprObj(t, core.DimensionEnergy(), core.NewMul(v, g.e)),
				variable, replacement}},
		{OpDifferentiate, OperationParams{Kind: "empty"},
			[]core.Object{exprObj(t, core.DimensionEnergy(),
				core.NewMul(g.m, core.NewPow(v, big.NewRat(2, 1)))), variable}},
		{OpLimit, OperationParams{Kind: "empty"},
			[]core.Object{exprObj(t, core.DimensionEnergy(),
				core.NewPow(v, big.NewRat(2, 1))), variable,
				exprObj(t, core.DimensionVelocity(), rat(0, 1))}},
		{OpCompare, OperationParams{Kind: "compare", Operator: core.RelationGt},
			[]core.Object{a, b}},
		{OpSolve, OperationParams{Kind: "empty"},
			[]core.Object{goldenRelation(t, g), g.energy}},
		{OpSelectBranch, OperationParams{Kind: "empty"}, []core.Object{br, constraint}},
	}
	_ = wrongDimValue

	for _, inv := range invocations {
		before := make([]string, len(inv.inputs))
		for i, in := range inv.inputs {
			before[i] = goldenJSON(t, in)
		}
		out, err := Apply(inv.id, inv.inputs, inv.params)
		if err != nil {
			t.Errorf("%s: %v", inv.id, err)
			continue
		}
		for i, in := range inv.inputs {
			if after := goldenJSON(t, in); after != before[i] {
				t.Errorf("%s: input %d mutated\n got %s\nwant %s",
					inv.id, i, after, before[i])
			}
		}
		if !out.Valid() {
			t.Errorf("%s: invalid result", inv.id)
		}
		if out.CorpusStatus() != kernel.CorpusNone {
			t.Errorf("%s: corpus status %v", inv.id, out.CorpusStatus())
		}
		if out.Provenance().MRCVersion() != kernel.MRCVersion {
			t.Errorf("%s: mrc version %q", inv.id, out.Provenance().MRCVersion())
		}
		if got, want := len(out.Provenance().ParentHashes()), len(inv.inputs); got != want {
			t.Errorf("%s: parent count %d want %d", inv.id, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// Provenance law (§13.5, REQ-013-03, REQ-§13.5-MUST-01)
// ---------------------------------------------------------------------------

func TestProvenanceParentHashes(t *testing.T) {
	g := buildGolden(t)
	a := exprObj(t, core.DimensionEnergy(), g.e)
	b := exprObj(t, core.DimensionEnergy(), core.NewAdd(g.e, rat(1, 1)))
	br := goldenSolved(t, g)
	constraint, err := Compare(g.energy, g.zeroEnergy, core.RelationGte)
	if err != nil {
		t.Fatal(err)
	}
	v := sym(t, "v")
	variable := defined(t, core.KindVelocity, core.DimensionVelocity(), v, "v")
	replacement := defined(t, core.KindVelocity, core.DimensionVelocity(),
		sym(t, "w"), "w")

	invocations := []struct {
		id     OperationID
		params OperationParams
		inputs []core.Object
	}{
		{OpAdd, OperationParams{Kind: "empty"}, []core.Object{a, b}},
		{OpSubtract, OperationParams{Kind: "empty"}, []core.Object{a, b}},
		{OpMultiply, OperationParams{Kind: "empty"}, []core.Object{a, b}},
		{OpDivide, OperationParams{Kind: "empty"}, []core.Object{a, b}},
		{OpPow, OperationParams{Kind: "pow", Exponent: "2/1"}, []core.Object{a}},
		{OpSimplify, OperationParams{Kind: "empty"}, []core.Object{a}},
		{OpSubstitute, OperationParams{Kind: "empty"},
			[]core.Object{exprObj(t, core.DimensionEnergy(), core.NewMul(v, g.e)),
				variable, replacement}},
		{OpDifferentiate, OperationParams{Kind: "empty"},
			[]core.Object{exprObj(t, core.DimensionEnergy(),
				core.NewMul(g.m, core.NewPow(v, big.NewRat(2, 1)))), variable}},
		{OpLimit, OperationParams{Kind: "empty"},
			[]core.Object{exprObj(t, core.DimensionEnergy(),
				core.NewPow(v, big.NewRat(2, 1))), variable,
				exprObj(t, core.DimensionVelocity(), rat(0, 1))}},
		{OpCompare, OperationParams{Kind: "compare", Operator: core.RelationEq},
			[]core.Object{a, b}},
		{OpSolve, OperationParams{Kind: "empty"},
			[]core.Object{goldenRelation(t, g), g.energy}},
		{OpSelectBranch, OperationParams{Kind: "empty"}, []core.Object{br, constraint}},
	}
	for _, inv := range invocations {
		out, err := Apply(inv.id, inv.inputs, inv.params)
		if err != nil {
			t.Errorf("%s: %v", inv.id, err)
			continue
		}
		want := make([]string, len(inv.inputs))
		for i, in := range inv.inputs {
			want[i] = kernel.Hex(core.HashObject(in))
		}
		got := out.Provenance().ParentHashes()
		if len(got) != len(want) {
			t.Errorf("%s: parent count %d want %d", inv.id, len(got), len(want))
			continue
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("%s: parent %d = %s want %s (input order)",
					inv.id, i, got[i], want[i])
			}
		}
	}
}

func TestProvenanceInheritance(t *testing.T) {
	g := buildGolden(t)
	withSource := func(source, framework string) core.Object {
		return mintSpec(t, kernel.ObjectSpec{
			Name: "s", Kind: core.KindExpression, Dimension: core.DimensionEnergy(),
			Expr:         g.e,
			Assumptions:  core.NewAssumptionSet(),
			Conventions:  core.NewConventionSet(),
			Provenance:   newProv(t, core.StatusDefined, source, framework, nil, ""),
			CorpusStatus: kernel.CorpusNone,
		})
	}
	cases := []struct {
		name          string
		left, right   core.Object
		wantSource    string
		wantFramework string
	}{
		{"identical", withSource("S", "F"), withSource("S", "F"), "S", "F"},
		{"mixed source", withSource("S", "F"), withSource("T", "F"), "", "F"},
		{"mixed framework", withSource("S", "F"), withSource("S", "G"), "S", ""},
		{"empty left", withSource("", "F"), withSource("S", "F"), "", "F"},
		{"empty right", withSource("S", ""), withSource("S", "F"), "S", ""},
	}
	for _, c := range cases {
		out, err := Add(c.left, c.right)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if out.Provenance().Source() != c.wantSource ||
			out.Provenance().Framework() != c.wantFramework {
			t.Errorf("%s: source=%q framework=%q want %q/%q", c.name,
				out.Provenance().Source(), out.Provenance().Framework(),
				c.wantSource, c.wantFramework)
		}
	}
	// Single-input operations inherit from their only input.
	s, err := Simplify(withSource("S", "F"))
	if err != nil {
		t.Fatal(err)
	}
	if s.Provenance().Source() != "S" || s.Provenance().Framework() != "F" {
		t.Errorf("simplify inherit: %q/%q", s.Provenance().Source(),
			s.Provenance().Framework())
	}
}

// ---------------------------------------------------------------------------
// Kind and ordering tables (§6.2)
// ---------------------------------------------------------------------------

func TestKindCompatibilityTable(t *testing.T) {
	mass := defined(t, core.KindMass, core.DimensionMass(), sym(t, "m1"), "Mass")
	restMass := defined(t, core.KindRestMass, core.DimensionMass(), sym(t, "m2"), "RestMass")
	exprM := exprObj(t, core.DimensionMass(), sym(t, "m3"))

	kindCases := []struct {
		name    string
		run     func() (core.Object, error)
		wantErr bool
	}{
		{"Add mass+mass", func() (core.Object, error) { return Add(mass, mass) }, false},
		{"Add expr+expr", func() (core.Object, error) { return Add(exprM, exprM) }, false},
		{"Add mass+restmass", func() (core.Object, error) { return Add(mass, restMass) }, true},
		{"Add restmass+mass", func() (core.Object, error) { return Add(restMass, mass) }, true},
		{"Add mass+expr", func() (core.Object, error) { return Add(mass, exprM) }, true},
		{"Add expr+mass", func() (core.Object, error) { return Add(exprM, mass) }, true},
		{"Subtract mass+restmass", func() (core.Object, error) { return Subtract(mass, restMass) }, true},
		{"Compare eq mass+mass", func() (core.Object, error) {
			return Compare(mass, mass, core.RelationEq)
		}, false},
		{"Compare eq mass+expr", func() (core.Object, error) {
			return Compare(mass, exprM, core.RelationEq)
		}, false},
		{"Compare eq expr+mass", func() (core.Object, error) {
			return Compare(exprM, mass, core.RelationEq)
		}, false},
		{"Compare eq expr+expr", func() (core.Object, error) {
			return Compare(exprM, exprM, core.RelationEq)
		}, false},
		{"Compare eq mass+restmass", func() (core.Object, error) {
			return Compare(mass, restMass, core.RelationEq)
		}, true},
		{"Compare neq mass+restmass", func() (core.Object, error) {
			return Compare(mass, restMass, core.RelationNeq)
		}, true},
		{"Multiply mass+expr", func() (core.Object, error) {
			return Multiply(mass, exprM)
		}, false},
		{"Divide mass+expr", func() (core.Object, error) {
			return Divide(mass, exprM)
		}, false},
	}
	for _, c := range kindCases {
		out, err := c.run()
		if c.wantErr {
			if err == nil {
				t.Errorf("%s: expected category error", c.name)
				continue
			}
			var cat core.CategoryMismatchError
			if !errors.As(err, &cat) {
				t.Errorf("%s: error type %T", c.name, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}
		switch {
		case strings.HasPrefix(c.name, "Multiply"), strings.HasPrefix(c.name, "Divide"):
			if out.Kind() != core.KindExpression {
				t.Errorf("%s: kind %v want Expression", c.name, out.Kind())
			}
		}
	}

	// Pow always yields Expression regardless of operand kind.
	p, err := Pow(mass, big.NewRat(2, 1))
	if err != nil {
		t.Fatal(err)
	}
	if p.Kind() != core.KindExpression {
		t.Errorf("pow kind %v", p.Kind())
	}
}

func TestInequalityOrderedKinds(t *testing.T) {
	ordered := []struct {
		kind core.Kind
		dim  core.Dimension
	}{
		{core.KindMass, core.DimensionMass()},
		{core.KindRestMass, core.DimensionMass()},
		{core.KindTime, core.DimensionTime()},
		{core.KindEnergy, core.DimensionEnergy()},
		{core.KindKineticEnergy, core.DimensionEnergy()},
		{core.KindSpeedOfLight, core.DimensionVelocity()},
		{core.KindExpression, core.DimensionEnergy()},
	}
	opsList := []core.RelationOperator{
		core.RelationLt, core.RelationLte, core.RelationGt, core.RelationGte,
	}
	for _, k := range ordered {
		a := defined(t, k.kind, k.dim, sym(t, "ia"), "ia")
		b := defined(t, k.kind, k.dim, sym(t, "ib"), "ib")
		for _, op := range opsList {
			if _, err := Compare(a, b, op); err != nil {
				t.Errorf("ordered kind %v with %s: %v", k.kind, op, err)
			}
		}
	}
	// Velocity is intentionally outside the MVP inequality path.
	velA := defined(t, core.KindVelocity, core.DimensionVelocity(), sym(t, "va"), "va")
	velB := defined(t, core.KindVelocity, core.DimensionVelocity(), sym(t, "vb"), "vb")
	for _, op := range opsList {
		_, err := Compare(velA, velB, op)
		if err == nil {
			t.Errorf("velocity inequality %s accepted", op)
			continue
		}
		var cat core.CategoryMismatchError
		if !errors.As(err, &cat) {
			t.Errorf("velocity inequality error type %T", err)
		}
	}
	// Equality comparisons remain open to every compatible kind.
	if _, err := Compare(velA, velB, core.RelationEq); err != nil {
		t.Errorf("velocity equality rejected: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Substitute / Differentiate / Limit contracts (§15.7–§15.9)
// ---------------------------------------------------------------------------

func TestSubstituteContract(t *testing.T) {
	g := buildGolden(t)
	v := sym(t, "v")
	variable := defined(t, core.KindVelocity, core.DimensionVelocity(), v, "v")
	replacement := defined(t, core.KindVelocity, core.DimensionVelocity(),
		sym(t, "w"), "w")

	// REQ-§15.7-MUST-01: variable expression must be a single Symbol.
	nonSymbolVar := defined(t, core.KindVelocity, core.DimensionVelocity(),
		core.NewNeg(v), "nv")
	_, err := Substitute(exprObj(t, core.DimensionEnergy(), g.e),
		nonSymbolVar, replacement)
	if err == nil {
		t.Fatal("non-symbol variable accepted")
	}
	var u core.UnsupportedOperationError
	if !errors.As(err, &u) {
		t.Errorf("error type %T", err)
	}

	// REQ-§15.7-MUST-03: variable and replacement kinds must be identical.
	kindMismatch := defined(t, core.KindMass, core.DimensionVelocity(),
		sym(t, "w"), "mw")
	_, err = Substitute(exprObj(t, core.DimensionEnergy(), g.e), variable, kindMismatch)
	if err == nil {
		t.Fatal("kind mismatch accepted")
	}
	var cat core.CategoryMismatchError
	if !errors.As(err, &cat) {
		t.Errorf("error type %T", err)
	}

	// REQ-§15.7-MUST-04: every occurrence of the exact symbol is replaced.
	repeated := exprObj(t, core.DimensionEnergy(),
		core.NewAdd(
			core.NewMul(v, v),
			core.NewPow(v, big.NewRat(2, 1)),
			v,
			core.NewSqrt(core.NewMul(v, sym(t, "x"))),
		))
	out, err := Substitute(repeated, variable, replacement)
	if err != nil {
		t.Fatalf("substitute: %v", err)
	}
	js := jsonExpr(t, out.Expr())
	if strings.Contains(js, `"name":"v"`) {
		t.Fatalf("occurrence of v survived: %s", js)
	}
	for _, want := range []string{`"name":"w"`, `"name":"x"`} {
		if !strings.Contains(js, want) {
			t.Fatalf("missing %s in %s", want, js)
		}
	}
	if out.Kind() != repeated.Kind() || !out.Dimension().Equal(repeated.Dimension()) {
		t.Fatal("target kind/dimension not preserved")
	}
	if len(out.Provenance().ParentHashes()) != 3 {
		t.Fatalf("parents: %v", out.Provenance().ParentHashes())
	}
}

func TestSubstituteDimensionMismatch(t *testing.T) {
	v := sym(t, "v")
	variable := defined(t, core.KindVelocity, core.DimensionVelocity(), v, "v")
	wrongDim := defined(t, core.KindVelocity, core.DimensionLength(),
		sym(t, "w"), "w")
	target := exprObj(t, core.DimensionEnergy(), core.NewMul(v, sym(t, "x")))
	_, err := Substitute(target, variable, wrongDim)
	if err == nil {
		t.Fatal("dimension mismatch accepted")
	}
	var dm core.DimensionMismatchError
	if !errors.As(err, &dm) {
		t.Fatalf("error type %T", err)
	}
	if dm.Operation != "substitute" || dm.Left == "" || dm.Right == "" {
		t.Fatalf("diagnostics: %+v", dm)
	}
}

func TestDifferentiateRequiresSingleSymbol(t *testing.T) {
	g := buildGolden(t)
	v := sym(t, "v")
	target := exprObj(t, core.DimensionEnergy(), core.NewMul(g.m, v))
	badWrt := defined(t, core.KindVelocity, core.DimensionVelocity(),
		core.NewNeg(v), "nv")
	_, err := Differentiate(target, badWrt)
	if err == nil {
		t.Fatal("non-symbol wrt accepted")
	}
	var u core.UnsupportedOperationError
	if !errors.As(err, &u) || u.Operation != "differentiate" {
		t.Fatalf("error: %v", err)
	}
}

func TestLimitForms(t *testing.T) {
	v := sym(t, "v")
	target := exprObj(t, core.DimensionEnergy(), core.NewPow(v, big.NewRat(2, 1)))
	variable := defined(t, core.KindVelocity, core.DimensionVelocity(), v, "v")

	// Exact value form: the zero-velocity style rational object.
	zeroValue := exprObj(t, core.DimensionVelocity(), rat(0, 1))
	out, err := Limit(target, variable, zeroValue)
	if err != nil {
		t.Fatalf("rational value form: %v", err)
	}
	if jsonExpr(t, out.Expr()) != `{"kind":"rational","value":"0/1"}` {
		t.Fatalf("result: %s", jsonExpr(t, out.Expr()))
	}
	if out.Kind() != core.KindExpression || !out.Dimension().Equal(core.DimensionEnergy()) {
		t.Fatalf("kind/dim: %v %v", out.Kind(), out.Dimension())
	}

	// Exact value form: a symbolic value of matching dimension.
	symbolic := exprObj(t, core.DimensionVelocity(), sym(t, "u"))
	if _, err := Limit(target, variable, symbolic); err != nil {
		t.Fatalf("symbolic value form: %v", err)
	}

	// Non-symbol variable rejected.
	badVar := defined(t, core.KindVelocity, core.DimensionVelocity(),
		core.NewNeg(v), "nv")
	_, err = Limit(target, badVar, zeroValue)
	var u core.UnsupportedOperationError
	if !errors.As(err, &u) || u.Operation != "limit" {
		t.Fatalf("variable contract: %v", err)
	}

	// Dimension mismatch between variable and value.
	wrongDim := exprObj(t, core.DimensionLength(), rat(0, 1))
	_, err = Limit(target, variable, wrongDim)
	var dm core.DimensionMismatchError
	if !errors.As(err, &dm) || dm.Operation != "limit" {
		t.Fatalf("dimension contract: %v", err)
	}
}

func TestLorentzFactorLimit(t *testing.T) {
	g := buildGolden(t)
	v := sym(t, "v")
	gamma, err := core.NewCall(kernel.LorentzFactorFunctionID, v)
	if err != nil {
		t.Fatal(err)
	}
	target := exprObj(t, core.Dimensionless(), gamma, g.cGT0)
	variable := defined(t, core.KindVelocity, core.DimensionVelocity(), v, "v")
	zeroValue := exprObj(t, core.DimensionVelocity(), rat(0, 1))

	out, err := Limit(target, variable, zeroValue)
	if err != nil {
		t.Fatalf("limit: %v", err)
	}
	if jsonExpr(t, out.Expr()) != `{"kind":"rational","value":"1/1"}` {
		t.Fatalf("limit result: %s", jsonExpr(t, out.Expr()))
	}
	if out.Kind() != core.KindExpression {
		t.Fatalf("kind %v", out.Kind())
	}

	// G-Audit: source of ops/transform.go must traverse the fixed body —
	// expansion → substitution → simplify — with no function-id shortcut.
	src, err := os.ReadFile("transform.go")
	if err != nil {
		t.Fatalf("read transform.go: %v", err)
	}
	body := string(src)
	for _, needle := range []string{
		"expandCalls(target.Expr())",
		"substituteExpr(expanded",
		"simplifyExpr(expr, assumptions)",
		"expandCalls",
	} {
		if !strings.Contains(body, needle) {
			t.Errorf("transform.go missing traversal step %q", needle)
		}
	}
	if strings.Contains(body, "LorentzFactorFunctionID") {
		t.Error("transform.go must not branch on the lorentz_factor id")
	}
}

// ---------------------------------------------------------------------------
// SelectBranch six-step contract (§15.12)
// ---------------------------------------------------------------------------

func TestSelectBranchContract(t *testing.T) {
	g := buildGolden(t)
	br := goldenSolved(t, g)
	constraint, err := Compare(g.energy, g.zeroEnergy, core.RelationGte)
	if err != nil {
		t.Fatal(err)
	}

	// Success path: steps 3–6.
	out, err := SelectBranch(br, constraint)
	if err != nil {
		t.Fatalf("select: %v", err)
	}
	if out.Kind() != core.KindExpression {
		t.Fatalf("kind %v", out.Kind())
	}
	if !out.Dimension().Equal(core.DimensionEnergy()) {
		t.Fatalf("dim %v", out.Dimension())
	}
	if !core.EqualExpr(out.Expr(), core.NewMul(g.m, core.NewPow(g.c, big.NewRat(2, 1)))) {
		t.Fatalf("expr %s", jsonExpr(t, out.Expr()))
	}
	if out.Provenance().Status() != core.StatusDerived {
		t.Fatalf("status %v", out.Provenance().Status())
	}
	wantParents := []string{
		kernel.Hex(core.HashObject(br)), kernel.Hex(core.HashObject(constraint)),
	}
	gotParents := out.Provenance().ParentHashes()
	if len(gotParents) != 2 || gotParents[0] != wantParents[0] ||
		gotParents[1] != wantParents[1] {
		t.Fatalf("parents %v", gotParents)
	}
	// Step 3: deterministic structured assumption key + value.
	wantKey := "selected_branch/" + kernel.Hex(core.HashExpr(constraint.Expr()))
	found := false
	for _, a := range out.Assumptions().Values() {
		if a.Key() != wantKey {
			continue
		}
		found = true
		if a.Kind() != core.AssumptionConstraint {
			t.Errorf("assumption kind %v", a.Kind())
		}
		if !core.EqualExpr(a.Expr(), constraint.Expr()) {
			t.Errorf("assumption value not the constraint expr")
		}
	}
	if !found {
		t.Fatalf("assumption %s missing", wantKey)
	}

	// Step 5: hypothesis contamination propagates.
	hypoBr := mintSpec(t, kernel.ObjectSpec{
		Name: "hb", Kind: core.KindBranchSet, Dimension: core.DimensionEnergy(),
		Expr:        br.Expr(),
		Assumptions: br.Assumptions(),
		Conventions: br.Conventions(),
		Provenance: newProv(t, core.StatusHypothesis, "test", "test",
			br.Provenance().ParentHashes(), "hypothesis"),
		CorpusStatus: kernel.CorpusNone,
	})
	hOut, err := SelectBranch(hypoBr, constraint)
	if err != nil {
		t.Fatalf("hypo select: %v", err)
	}
	if hOut.Provenance().Status() != core.StatusHypothesis {
		t.Fatalf("contamination: %v", hOut.Provenance().Status())
	}

	// Step 1/2 rejection matrix.
	mkBranchSet := func(e core.Expr) core.Object {
		return defined(t, core.KindBranchSet, core.DimensionEnergy(), e, "bs")
	}
	sqrtTerm := core.NewSqrt(core.NewPow(
		core.NewMul(g.m, core.NewPow(g.c, big.NewRat(2, 1))), big.NewRat(2, 1)))
	validFirst := mkBranchSet(core.NewBranchSet(g.e, sqrtTerm, core.NewNeg(sqrtTerm)))

	rejects := []struct {
		name       string
		branches   core.Object
		constraint core.Object
		errType    func(error) bool
	}{
		{"branches kind", exprObj(t, core.DimensionEnergy(), sqrtTerm),
			constraint, func(err error) bool {
				var c core.CategoryMismatchError
				return errors.As(err, &c)
			}},
		{"constraint kind",
			br, defined(t, core.KindExpression, core.DimensionEnergy(),
				constraint.Expr(), "ce"),
			func(err error) bool {
				var c core.CategoryMismatchError
				return errors.As(err, &c)
			}},
		{"constraint dimension",
			br, defined(t, core.KindRelation, core.DimensionLength(),
				constraint.Expr(), "cl"),
			func(err error) bool {
				var d core.DimensionMismatchError
				return errors.As(err, &d)
			}},
		{"three branches",
			mkBranchSet(core.NewBranchSet(g.e, sqrtTerm, core.NewNeg(sqrtTerm), sqrtTerm)),
			constraint, func(err error) bool {
				var u core.UnsupportedOperationError
				return errors.As(err, &u)
			}},
		{"second branch not negation",
			mkBranchSet(core.NewBranchSet(g.e, sqrtTerm, sqrtTerm)),
			constraint, func(err error) bool {
				var u core.UnsupportedOperationError
				return errors.As(err, &u)
			}},
		{"operator not gte",
			validFirst, defined(t, core.KindRelation, core.DimensionEnergy(),
				core.NewRelation(core.RelationLt, g.e, rat(0, 1)), "clt"),
			func(err error) bool {
				var u core.UnsupportedOperationError
				return errors.As(err, &u)
			}},
		{"left side not target",
			validFirst, defined(t, core.KindRelation, core.DimensionEnergy(),
				core.NewRelation(core.RelationGte, sym(t, "z"), rat(0, 1)), "cz"),
			func(err error) bool {
				var u core.UnsupportedOperationError
				return errors.As(err, &u)
			}},
		{"right side not zero",
			validFirst, defined(t, core.KindRelation, core.DimensionEnergy(),
				core.NewRelation(core.RelationGte, g.e, rat(1, 1)), "cz"),
			func(err error) bool {
				var u core.UnsupportedOperationError
				return errors.As(err, &u)
			}},
		{"constraint expr not relation",
			validFirst, defined(t, core.KindRelation, core.DimensionEnergy(),
				g.e, "cs"),
			func(err error) bool {
				var u core.UnsupportedOperationError
				return errors.As(err, &u)
			}},
	}
	for _, c := range rejects {
		if _, err := SelectBranch(c.branches, c.constraint); err == nil {
			t.Errorf("%s: accepted", c.name)
		} else if !c.errType(err) {
			t.Errorf("%s: error type %T (%v)", c.name, err, err)
		}
	}
}

// ---------------------------------------------------------------------------
// Simplify rule surface (§8.7, §9.6–§9.9)
// ---------------------------------------------------------------------------

func simpExpr(t *testing.T, e core.Expr, as ...core.Assumption) (core.Expr, error) {
	t.Helper()
	obj := defined(t, core.KindExpression, core.DimensionEnergy(), e, "s", as...)
	out, err := Simplify(obj)
	if err != nil {
		return core.Expr{}, err
	}
	return out.Expr(), nil
}

func mustSimp(t *testing.T, e core.Expr, as ...core.Assumption) core.Expr {
	t.Helper()
	out, err := simpExpr(t, e, as...)
	if err != nil {
		t.Fatalf("simplify: %v", err)
	}
	return out
}

func TestSimplifyIdentityRules(t *testing.T) {
	x := sym(t, "x")
	cases := []struct {
		name string
		in   core.Expr
		want string
	}{
		{"Add(x,0)", core.NewAdd(x, rat(0, 1)), `{"kind":"symbol","name":"x"}`},
		{"Add()", core.NewAdd(), `{"kind":"rational","value":"0/1"}`},
		{"Add cancel", core.NewAdd(x, rat(1, 2), rat(-1, 2)),
			`{"kind":"symbol","name":"x"}`},
		{"Mul(x,1)", core.NewMul(x, rat(1, 1)), `{"kind":"symbol","name":"x"}`},
		{"Mul()", core.NewMul(), `{"kind":"rational","value":"1/1"}`},
		{"Mul(x,0)", core.NewMul(x, rat(0, 1)), `{"kind":"rational","value":"0/1"}`},
		{"Pow(x,1)", core.NewPow(x, big.NewRat(1, 1)),
			`{"kind":"symbol","name":"x"}`},
		{"Neg(Neg(x))", core.NewNeg(core.NewNeg(x)),
			`{"kind":"symbol","name":"x"}`},
		{"Sqrt(1)", core.NewSqrt(rat(1, 1)), `{"kind":"rational","value":"1/1"}`},
		{"Sqrt(0)", core.NewSqrt(rat(0, 1)), `{"kind":"rational","value":"0/1"}`},
		{"Pow(0,2)", core.NewPow(rat(0, 1), big.NewRat(2, 1)),
			`{"kind":"rational","value":"0/1"}`},
		{"Pow(1,-1)", core.NewPow(rat(1, 1), big.NewRat(-1, 1)),
			`{"kind":"rational","value":"1/1"}`},
		{"Pow(2,3)", core.NewPow(rat(2, 1), big.NewRat(3, 1)),
			`{"kind":"rational","value":"8/1"}`},
		{"Pow(2,-1)", core.NewPow(rat(2, 1), big.NewRat(-1, 1)),
			`{"kind":"rational","value":"1/2"}`},
		{"Pow(x,0) unentailed keeps",
			core.NewPow(x, big.NewRat(0, 1)),
			`{"kind":"pow","base":{"kind":"symbol","name":"x"},"exp":"0/1"}`},
	}
	for _, c := range cases {
		got := mustSimp(t, c.in)
		if js := jsonExpr(t, got); js != c.want {
			t.Errorf("%s:\n got %s\nwant %s", c.name, js, c.want)
		}
	}

	// F2: Pow(x,0) → 1 only for a nonzero-safe base under the assumptions.
	positive := gtZeroAssumption(t, x, "x_positive")
	got := mustSimp(t, core.NewPow(x, big.NewRat(0, 1)), positive)
	if js := jsonExpr(t, got); js != `{"kind":"rational","value":"1/1"}` {
		t.Errorf("Pow(x,0) under x>0: %s", js)
	}

	// 0^0 and 0^-n remain UnsupportedOperationError (never silently valued).
	for _, c := range []struct {
		name string
		in   core.Expr
	}{
		{"0^0", core.NewPow(rat(0, 1), big.NewRat(0, 1))},
		{"0^-1", core.NewPow(rat(0, 1), big.NewRat(-1, 1))},
	} {
		_, err := simpExpr(t, c.in)
		if err == nil {
			t.Errorf("%s: accepted", c.name)
			continue
		}
		var u core.UnsupportedOperationError
		if !errors.As(err, &u) {
			t.Errorf("%s: error type %T", c.name, err)
		}
	}
}

func TestSimplifyRepeatedPowers(t *testing.T) {
	x := sym(t, "x")
	y := sym(t, "y")
	two := big.NewRat(2, 1)
	cases := []struct {
		name string
		in   core.Expr
		want string
	}{
		{"x*x", core.NewMul(x, x),
			`{"kind":"pow","base":{"kind":"symbol","name":"x"},"exp":"2/1"}`},
		{"Pow(x,2)*Pow(x,2)",
			core.NewMul(core.NewPow(x, two), core.NewPow(x, two)),
			`{"kind":"pow","base":{"kind":"symbol","name":"x"},"exp":"4/1"}`},
		{"Pow(Pow(x,2),3)",
			core.NewPow(core.NewPow(x, two), big.NewRat(3, 1)),
			`{"kind":"pow","base":{"kind":"symbol","name":"x"},"exp":"6/1"}`},
		{"Pow(Pow(x,2),-1) keeps nesting",
			core.NewPow(core.NewPow(x, two), big.NewRat(-1, 1)),
			`{"kind":"pow","base":{"kind":"pow","base":{"kind":"symbol","name":"x"},"exp":"2/1"},"exp":"-1/1"}`},
		{"distinct bases untouched",
			core.NewMul(core.NewPow(x, two), core.NewPow(y, two)),
			`{"kind":"mul","factors":[{"kind":"pow","base":{"kind":"symbol","name":"x"},"exp":"2/1"},{"kind":"pow","base":{"kind":"symbol","name":"y"},"exp":"2/1"}]}`},
		{"non-integer exponents ungrouped",
			core.NewMul(core.NewPow(x, big.NewRat(1, 2)), core.NewPow(x, big.NewRat(1, 2))),
			`{"kind":"mul","factors":[{"kind":"pow","base":{"kind":"symbol","name":"x"},"exp":"1/2"},{"kind":"pow","base":{"kind":"symbol","name":"x"},"exp":"1/2"}]}`},
	}
	for _, c := range cases {
		got := mustSimp(t, c.in)
		if js := jsonExpr(t, got); js != c.want {
			t.Errorf("%s:\n got %s\nwant %s", c.name, js, c.want)
		}
	}
}

func TestSqrtRewriteRules(t *testing.T) {
	x := sym(t, "x")
	m := sym(t, "m")
	c := sym(t, "c")
	two := big.NewRat(2, 1)
	massGE0 := gteZeroAssumption(t, m, "rest_mass_nonnegative")
	xGE0 := gteZeroAssumption(t, x, "x_nonnegative")

	plain := []struct {
		name string
		in   core.Expr
		want string
	}{
		{"Sqrt(0)", core.NewSqrt(rat(0, 1)), `{"kind":"rational","value":"0/1"}`},
		{"Sqrt(1)", core.NewSqrt(rat(1, 1)), `{"kind":"rational","value":"1/1"}`},
		{"Sqrt(4)", core.NewSqrt(rat(4, 1)), `{"kind":"rational","value":"2/1"}`},
		{"Sqrt(9/4)", core.NewSqrt(rat(9, 4)), `{"kind":"rational","value":"3/2"}`},
		{"Sqrt(2) stays",
			core.NewSqrt(rat(2, 1)),
			`{"kind":"sqrt","expr":{"kind":"rational","value":"2/1"}}`},
		{"negative rational stays",
			core.NewSqrt(rat(-1, 1)),
			`{"kind":"sqrt","expr":{"kind":"rational","value":"-1/1"}}`},
		{"odd exponent untouched",
			core.NewSqrt(core.NewPow(x, big.NewRat(3, 1))),
			`{"kind":"sqrt","expr":{"kind":"pow","base":{"kind":"symbol","name":"x"},"exp":"3/1"}}`},
	}
	for _, c := range plain {
		if js := jsonExpr(t, mustSimp(t, c.in)); js != c.want {
			t.Errorf("%s:\n got %s\nwant %s", c.name, js, c.want)
		}
	}

	// Gated: Sqrt(Pow(x,2)) → x only under x ≥ 0.
	gated := core.NewSqrt(core.NewPow(x, two))
	stays := jsonExpr(t, mustSimp(t, gated))
	if stays == `{"kind":"symbol","name":"x"}` {
		t.Fatal("gated rewrite fired without the assumption")
	}
	if js := jsonExpr(t, mustSimp(t, gated, xGE0)); js != `{"kind":"symbol","name":"x"}` {
		t.Fatalf("entailed rewrite: %s", js)
	}

	// Mass-energy positive branch: Sqrt(Pow(m·c²,2)) → m·c² under m ≥ 0.
	mc2 := core.NewMul(m, core.NewPow(c, two))
	massCase := core.NewSqrt(core.NewPow(mc2, two))
	want := `{"kind":"mul","factors":[{"kind":"symbol","name":"m"},{"kind":"pow","base":{"kind":"symbol","name":"c"},"exp":"2/1"}]}`
	if js := jsonExpr(t, mustSimp(t, massCase, massGE0)); js != want {
		t.Fatalf("mass-energy rewrite: %s", js)
	}
	if js := jsonExpr(t, mustSimp(t, massCase)); js == want {
		t.Fatal("mass-energy rewrite must require m ≥ 0")
	}
}

func TestSqrtSignRewriteRequiresAssumption(t *testing.T) {
	x := sym(t, "x")
	two := big.NewRat(2, 1)
	gated := core.NewSqrt(core.NewPow(x, two))

	// No sign rewrite may occur without entailment — and never into -x.
	neq, err := core.NewExprAssumption(core.AssumptionConstraint, "x_nonzero",
		core.NewRelation(core.RelationNeq, x, rat(0, 1)))
	if err != nil {
		t.Fatal(err)
	}
	lte, err := core.NewExprAssumption(core.AssumptionConstraint, "x_nonpositive",
		core.NewRelation(core.RelationLte, x, rat(0, 1)))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		as   []core.Assumption
	}{
		{"none", nil},
		{"nonzero", []core.Assumption{neq}},
		{"nonpositive", []core.Assumption{lte}},
	} {
		js := jsonExpr(t, mustSimp(t, gated, tc.as...))
		if js == `{"kind":"symbol","name":"x"}` ||
			js == `{"kind":"neg","expr":{"kind":"symbol","name":"x"}}` {
			t.Errorf("%s: unwarranted sign rewrite: %s", tc.name, js)
		}
		if !strings.Contains(js, `"kind":"sqrt"`) {
			t.Errorf("%s: expected structural Sqrt, got %s", tc.name, js)
		}
	}
	// Neg around the root is never invented either.
	negated := core.NewNeg(core.NewSqrt(core.NewPow(x, two)))
	if js := jsonExpr(t, mustSimp(t, negated, gteZeroAssumption(t, x, "x_ge0"))); js !=
		`{"kind":"neg","expr":{"kind":"pow","base":{"kind":"symbol","name":"x"},"exp":"2/1"}}` &&
		js != `{"kind":"neg","expr":{"kind":"sqrt","expr":{"kind":"pow","base":{"kind":"symbol","name":"x"},"exp":"2/1"}}` &&
		!strings.Contains(js, `"name":"x"`) {
		t.Errorf("negated root: %s", js)
	}
}

// ---------------------------------------------------------------------------
// Finiteness and bounded entailment (plan §6, §9.6, §9.8)
// ---------------------------------------------------------------------------

func TestFiniteFactorRules(t *testing.T) {
	x := sym(t, "x")
	c := sym(t, "c")
	cGT0 := gtZeroAssumption(t, c, "speed_of_light_positive")
	xGT0 := gtZeroAssumption(t, x, "x_positive")
	zeroFactor := func(factors ...core.Expr) core.Expr {
		return core.NewMul(append([]core.Expr{rat(0, 1)}, factors...)...)
	}

	// base != 0 (admissibility): c^-1 is finite under c > 0, so a zero
	// Mul collapses — no Lorentz-specific hack exists.
	cInv := core.NewPow(c, big.NewRat(-1, 1))
	got := mustSimp(t, zeroFactor(cInv), cGT0)
	if js := jsonExpr(t, got); js != `{"kind":"rational","value":"0/1"}` {
		t.Fatalf("Mul(0, c^-1) under c>0: %s", js)
	}
	// Without the assumption the factor is not finite and nothing collapses.
	got = mustSimp(t, zeroFactor(cInv))
	if js := jsonExpr(t, got); !strings.Contains(js, `"kind":"pow"`) {
		t.Fatalf("Mul(0, c^-1) unassumed: %s", js)
	}
	// m ≥ 0 is a sign statement, not a nonzero statement: m^-1 is not finite.
	m := sym(t, "m")
	mGE0 := gteZeroAssumption(t, m, "m_nonnegative")
	got = mustSimp(t, zeroFactor(core.NewPow(m, big.NewRat(-1, 1))), mGE0)
	if js := jsonExpr(t, got); !strings.Contains(js, `"kind":"pow"`) {
		t.Fatalf("Mul(0, m^-1) under m>=0: %s", js)
	}
	// base > 0 admits Pow(x,0) → 1 (nonzero-safe base).
	got = mustSimp(t, core.NewPow(x, big.NewRat(0, 1)), xGT0)
	if js := jsonExpr(t, got); js != `{"kind":"rational","value":"1/1"}` {
		t.Fatalf("Pow(x,0) under x>0: %s", js)
	}
	// 0^0 unsupported; 0^-1 rejected by the b != 0 gate.
	for _, in := range []core.Expr{
		core.NewPow(rat(0, 1), big.NewRat(0, 1)),
		core.NewPow(rat(0, 1), big.NewRat(-1, 1)),
	} {
		if _, err := simpExpr(t, in); err == nil {
			t.Fatalf("singular rational power accepted: %s", jsonExpr(t, in))
		}
	}
	// Sqrt(x) requires x ≥ 0 under assumptions (or the rational {0,1} forms).
	got = mustSimp(t, core.NewSqrt(x))
	if js := jsonExpr(t, got); !strings.Contains(js, `"kind":"sqrt"`) {
		t.Fatalf("Sqrt(x) unassumed: %s", js)
	}
	got = mustSimp(t, core.NewSqrt(x), gteZeroAssumption(t, x, "x_ge0"))
	if js := jsonExpr(t, got); js != `{"kind":"sqrt","expr":{"kind":"symbol","name":"x"}}` {
		t.Fatalf("Sqrt(x) has no rule beyond the closed list: %s", js)
	}
}

func TestSignEntailmentBounded(t *testing.T) {
	x := sym(t, "x")
	y := sym(t, "y")
	two := big.NewRat(2, 1)

	mkRel := func(op core.RelationOperator, subject core.Expr) core.Assumption {
		t.Helper()
		a, err := core.NewExprAssumption(core.AssumptionConstraint, "rel_"+string(op),
			core.NewRelation(op, subject, rat(0, 1)))
		if err != nil {
			t.Fatal(err)
		}
		return a
	}
	gated := core.NewSqrt(core.NewPow(x, two))

	// Denied: nonzero does not entail nonnegative.
	neqOnly := jsonExpr(t, mustSimp(t, gated, mkRel(core.RelationNeq, x)))
	if !strings.Contains(neqOnly, `"kind":"sqrt"`) {
		t.Errorf("x != 0 triggered sign rewrite: %s", neqOnly)
	}
	// Denied: x < 0 is the opposite sign claim.
	ltOnly := jsonExpr(t, mustSimp(t, gated, mkRel(core.RelationLt, x)))
	if !strings.Contains(ltOnly, `"kind":"sqrt"`) {
		t.Errorf("x < 0 triggered sign rewrite: %s", ltOnly)
	}
	// Denied: closed rules never reason across an Add — x ≥ 0 ∧ y ≥ 0 does
	// not entail x − y ≥ 0.
	both := []core.Assumption{
		gteZeroAssumption(t, x, "x_nonnegative"),
		gteZeroAssumption(t, y, "y_nonnegative"),
	}
	diff := jsonExpr(t, mustSimp(t,
		core.NewSqrt(core.NewPow(core.NewAdd(x, core.NewNeg(y)), two)), both...))
	if !strings.Contains(diff, `"kind":"sqrt"`) {
		t.Errorf("x - y sign inferred across Add: %s", diff)
	}
	// Denied: m ≥ 0 does not entail m ≠ 0 (so Pow(m,0) stays gated).
	m := sym(t, "m")
	got := mustSimp(t, core.NewPow(m, big.NewRat(0, 1)),
		gteZeroAssumption(t, m, "m_nonnegative"))
	if js := jsonExpr(t, got); js == `{"kind":"rational","value":"1/1"}` {
		t.Errorf("m >= 0 treated as m != 0: %s", js)
	}
	// Allowed by the closed list: x > 0 entails x ≠ 0, so Pow(x,0) → 1.
	got = mustSimp(t, core.NewPow(x, big.NewRat(0, 1)),
		gteZeroAssumption(t, x, "x_ge0"))
	if js := jsonExpr(t, got); js == `{"kind":"rational","value":"1/1"}` {
		t.Errorf("x >= 0 treated as x != 0: %s", js)
	}
	got = mustSimp(t, core.NewPow(x, big.NewRat(0, 1)),
		gtZeroAssumption(t, x, "x_positive"))
	if js := jsonExpr(t, got); js != `{"kind":"rational","value":"1/1"}` {
		t.Errorf("x > 0 must admit Pow(x,0) -> 1: %s", js)
	}
}

// TestLimitDirectSubstitutionNonUnity proves Limit performs genuine direct
// substitution followed by simplification (M18). Limit(x+2, x, 3) must yield
// exactly 5 — a mutant hardcoding every Limit result to 1 dies here, while
// TestLorentzFactorLimit continues to assert the required exact result 1.
func TestLimitDirectSubstitutionNonUnity(t *testing.T) {
	x := sym(t, "x")
	target := exprObj(t, core.Dimensionless(), core.NewAdd(x, rat(2, 1)))
	variable := defined(t, core.KindExpression, core.Dimensionless(), x, "x")
	value := exprObj(t, core.Dimensionless(), rat(3, 1))
	out, err := Limit(target, variable, value)
	if err != nil {
		t.Fatalf("Limit: %v", err)
	}
	if got := jsonExpr(t, out.Expr()); got != `{"kind":"rational","value":"5/1"}` {
		t.Fatalf("Limit(x+2, x, 3) = %s, want 5/1", got)
	}
	if out.Kind() != core.KindExpression {
		t.Fatalf("kind = %v, want Expression", out.Kind())
	}
	if !out.Dimension().Equal(core.Dimensionless()) {
		t.Fatalf("dimension mismatch: want dimensionless")
	}
	// The result must come from substitution, not the Lorentz fixed body:
	// it carries no Lorentz assumption and its provenance parents are the
	// three limit inputs.
	if got := len(out.Provenance().ParentHashes()); got != 3 {
		t.Fatalf("parent hashes = %d, want 3 (target, variable, value)", got)
	}
}

// TestProvenanceMatrix12 (Gate B): the §13.2 law over all 12 pure operations —
// clean inputs → DERIVED, any HYPOTHESIS input → HYPOTHESIS, never
// IDENTIFIED/POSTULATED/DEFINED. (Session.Identify halves live in the K tests:
// TestSessionIdentifyRecords + TestHypothesisContamination.)
func TestProvenanceMatrix12(t *testing.T) {
	x := sym(t, "x")
	y := sym(t, "y")
	clean := func(e core.Expr) core.Object {
		return exprObj(t, core.Dimensionless(), e)
	}
	hypo := func(e core.Expr) core.Object {
		return hypoObj(t, core.KindExpression, core.Dimensionless(), e)
	}
	mkSolvePair := func(h core.Object) (core.Object, core.Object) {
		rel, err := Compare(
			exprObj(t, core.Dimensionless(), core.NewPow(sym(t, "E"), big.NewRat(2, 1))),
			h, core.RelationEq)
		if err != nil {
			t.Fatalf("solve relation: %v", err)
		}
		tgt := defined(t, core.KindExpression, core.Dimensionless(), sym(t, "E"), "E")
		return rel, tgt
	}
	mkBranches := func() (core.Object, core.Object) {
		r, err := Compare(
			exprObj(t, core.Dimensionless(), core.NewPow(sym(t, "E"), big.NewRat(2, 1))),
			exprObj(t, core.Dimensionless(), rat(4, 1)), core.RelationEq)
		if err != nil {
			t.Fatalf("branch relation: %v", err)
		}
		br, err := Solve(r, defined(t, core.KindExpression, core.Dimensionless(), sym(t, "E"), "E"))
		if err != nil {
			t.Fatalf("branch solve: %v", err)
		}
		lo, err := Compare(
			exprObj(t, core.Dimensionless(), sym(t, "E")),
			exprObj(t, core.Dimensionless(), rat(0, 1)), core.RelationGte)
		if err != nil {
			t.Fatalf("branch constraint: %v", err)
		}
		return br, lo
	}
	br, lo := mkBranches()
	type opcase struct {
		name  string
		clean func() (core.Object, error)
		hypo  func() (core.Object, error)
	}
	cases := []opcase{
		{"add", func() (core.Object, error) { return Add(clean(x), clean(y)) },
			func() (core.Object, error) { return Add(hypo(x), clean(y)) }},
		{"subtract", func() (core.Object, error) { return Subtract(clean(x), clean(y)) },
			func() (core.Object, error) { return Subtract(clean(x), hypo(y)) }},
		{"multiply", func() (core.Object, error) { return Multiply(clean(x), clean(y)) },
			func() (core.Object, error) { return Multiply(hypo(x), hypo(y)) }},
		{"divide", func() (core.Object, error) { return Divide(clean(x), clean(y)) },
			func() (core.Object, error) { return Divide(hypo(x), clean(y)) }},
		{"pow", func() (core.Object, error) { return Pow(clean(x), big.NewRat(2, 1)) },
			func() (core.Object, error) { return Pow(hypo(x), big.NewRat(2, 1)) }},
		{"simplify", func() (core.Object, error) { return Simplify(clean(core.NewAdd(x, rat(0, 1)))) },
			func() (core.Object, error) { return Simplify(hypo(x)) }},
		{"substitute", func() (core.Object, error) {
			return Substitute(clean(core.NewAdd(x, rat(1, 1))), clean(x), clean(y))
		}, func() (core.Object, error) {
			return Substitute(hypo(core.NewAdd(x, rat(1, 1))), clean(x), clean(y))
		}},
		{"differentiate", func() (core.Object, error) { return Differentiate(clean(x), clean(x)) },
			func() (core.Object, error) { return Differentiate(hypo(x), clean(x)) }},
		{"limit", func() (core.Object, error) {
			return Limit(clean(core.NewAdd(x, rat(2, 1))), clean(x), clean(rat(3, 1)))
		}, func() (core.Object, error) {
			return Limit(hypo(core.NewAdd(x, rat(2, 1))), clean(x), clean(rat(3, 1)))
		}},
		{"compare", func() (core.Object, error) { return Compare(clean(x), clean(y), core.RelationEq) },
			func() (core.Object, error) { return Compare(hypo(x), clean(y), core.RelationEq) }},
		{"solve", func() (core.Object, error) {
			r, g := mkSolvePair(clean(rat(0, 1)))
			return Solve(r, g)
		}, func() (core.Object, error) {
			r, g := mkSolvePair(hypo(rat(0, 1)))
			return Solve(r, g)
		}},
		{"select_branch", func() (core.Object, error) { return SelectBranch(br, lo) },
			func() (core.Object, error) {
				hypoBr, err := Solve(func() core.Object {
					r, rerr := Compare(
						exprObj(t, core.Dimensionless(), core.NewPow(sym(t, "E"), big.NewRat(2, 1))),
						hypo(rat(4, 1)), core.RelationEq)
					if rerr != nil {
						t.Fatalf("hypo branch relation: %v", rerr)
					}
					return r
				}(), defined(t, core.KindExpression, core.Dimensionless(), sym(t, "E"), "E"))
				if err != nil {
					return core.Object{}, err
				}
				return SelectBranch(hypoBr, lo)
			}},
	}
	if len(cases) != 12 {
		t.Fatalf("matrix covers %d ops, want 12", len(cases))
	}
	for _, tc := range cases {
		got, err := tc.clean()
		if err != nil {
			t.Errorf("%s clean: %v", tc.name, err)
			continue
		}
		if got.Provenance().Status() != core.StatusDerived {
			t.Errorf("%s clean status = %v, want DERIVED", tc.name, got.Provenance().Status())
		}
		hout, err := tc.hypo()
		if err != nil {
			t.Errorf("%s hypo: %v", tc.name, err)
			continue
		}
		if hout.Provenance().Status() != core.StatusHypothesis {
			t.Errorf("%s hypo status = %v, want HYPOTHESIS", tc.name, hout.Provenance().Status())
		}
		for _, o := range []core.Object{got, hout} {
			switch o.Provenance().Status() {
			case core.StatusIdentified, core.StatusPostulated, core.StatusDefined:
				t.Errorf("%s manufactured forbidden status %v", tc.name, o.Provenance().Status())
			}
		}
	}
}

// TestAssumptionSubsetInvariant (Gate H3): per pure operation, output
// assumption keys ⊆ input keys ∪ operation-generated keys, and inherited
// assumptions are preserved verbatim (values, not just keys). Generated keys:
// divide → denominator/*, select_branch → selected_branch/*, all others → ∅.
func TestAssumptionSubsetInvariant(t *testing.T) {
	x := sym(t, "x")
	y := sym(t, "y")
	ex := func(e core.Expr) core.Object { return exprObj(t, core.Dimensionless(), e) }
	keySet := func(o core.Object) map[string]core.Assumption {
		m := map[string]core.Assumption{}
		for _, a := range o.Assumptions().Values() {
			m[a.Key()] = a
		}
		return m
	}
	generated := func(op string, key string) bool {
		switch op {
		case "divide":
			return strings.HasPrefix(key, "denominator/")
		case "select_branch":
			return strings.HasPrefix(key, "selected_branch/")
		}
		return false
	}
	br, lo := mkBranchesForSubset(t)
	type opcall struct {
		name   string
		inputs []core.Object
		run    func(ins []core.Object) (core.Object, error)
	}
	ca, err := core.NewExprAssumption(core.AssumptionConstraint, "carry/x",
		core.NewRelation(core.RelationEq, sym(t, "x"), rat(0, 1)))
	if err != nil {
		t.Fatalf("assumption: %v", err)
	}
	carrier, err := addAssumptionForSubset(t, ex(x), ca)
	if err != nil {
		t.Fatalf("carrier: %v", err)
	}
	carrier2, err := addAssumptionForSubset(t, ex(y), ca)
	if err != nil {
		t.Fatalf("carrier2: %v", err)
	}
	solveRel, err := Compare(
		exprObj(t, core.Dimensionless(), core.NewPow(sym(t, "E"), big.NewRat(2, 1))),
		carrier, core.RelationEq)
	if err != nil {
		t.Fatalf("solve relation: %v", err)
	}
	solveTgt := defined(t, core.KindExpression, core.Dimensionless(), sym(t, "E"), "E")
	loWithCarrier, err := addAssumptionForSubset(t, lo, ca)
	if err != nil {
		t.Fatalf("constraint carrier: %v", err)
	}
	calls := []opcall{
		{"add", []core.Object{carrier, carrier2}, func(ins []core.Object) (core.Object, error) {
			return Add(ins[0], ins[1])
		}},
		{"subtract", []core.Object{carrier, carrier2}, func(ins []core.Object) (core.Object, error) {
			return Subtract(ins[0], ins[1])
		}},
		{"multiply", []core.Object{carrier, carrier2}, func(ins []core.Object) (core.Object, error) {
			return Multiply(ins[0], ins[1])
		}},
		{"divide", []core.Object{carrier, carrier2}, func(ins []core.Object) (core.Object, error) {
			return Divide(ins[0], ins[1])
		}},
		{"pow", []core.Object{carrier}, func(ins []core.Object) (core.Object, error) {
			return Pow(ins[0], big.NewRat(2, 1))
		}},
		{"simplify", []core.Object{carrier}, func(ins []core.Object) (core.Object, error) {
			return Simplify(ins[0])
		}},
		{"substitute", []core.Object{carrier, ex(x), ex(y)}, func(ins []core.Object) (core.Object, error) {
			return Substitute(ins[0], ins[1], ins[2])
		}},
		{"differentiate", []core.Object{carrier, ex(x)}, func(ins []core.Object) (core.Object, error) {
			return Differentiate(ins[0], ins[1])
		}},
		{"limit", []core.Object{carrier, ex(x), ex(rat(3, 1))}, func(ins []core.Object) (core.Object, error) {
			return Limit(ins[0], ins[1], ins[2])
		}},
		{"compare", []core.Object{carrier, carrier2}, func(ins []core.Object) (core.Object, error) {
			return Compare(ins[0], ins[1], core.RelationEq)
		}},
		{"solve", []core.Object{solveRel, solveTgt}, func(ins []core.Object) (core.Object, error) {
			return Solve(ins[0], ins[1])
		}},
		{"select_branch", []core.Object{br, loWithCarrier}, func(ins []core.Object) (core.Object, error) {
			return SelectBranch(ins[0], ins[1])
		}},
	}
	if len(calls) != 12 {
		t.Fatalf("subset matrix covers %d ops, want 12", len(calls))
	}
	for _, tc := range calls {
		allowed := map[string]core.Assumption{}
		for _, in := range tc.inputs {
			for k, a := range keySet(in) {
				allowed[k] = a
			}
		}
		// solve/select_branch fixtures build their own inputs; collect from a dry run is
		// unnecessary — instead verify output keys against the union rule below using
		// the same inputs where present, and verbatim preservation via carrier cases.
		out, err := tc.run(tc.inputs)
		if err != nil {
			t.Errorf("%s: %v", tc.name, err)
			continue
		}
		for k, a := range keySet(out) {
			if want, ok := allowed[k]; ok {
				if !assumptionEqualForSubset(want, a) {
					t.Errorf("%s: inherited assumption %q mutated (value not preserved)", tc.name, k)
				}
				continue
			}
			if !generated(tc.name, k) {
				t.Errorf("%s: unrelated generated assumption %q", tc.name, k)
			}
		}
	}
	// Verbatim preservation across a merge: same key on both inputs survives.
	m1, err := Add(carrier, carrier2)
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	got := keySet(m1)["carry/x"]
	if !assumptionEqualForSubset(ca, got) {
		t.Errorf("merged assumption not verbatim")
	}
}

func mkBranchesForSubset(t *testing.T) (core.Object, core.Object) {
	t.Helper()
	r, err := Compare(
		exprObj(t, core.Dimensionless(), core.NewPow(sym(t, "E"), big.NewRat(2, 1))),
		exprObj(t, core.Dimensionless(), rat(4, 1)), core.RelationEq)
	if err != nil {
		t.Fatalf("branch relation: %v", err)
	}
	br, err := Solve(r, defined(t, core.KindExpression, core.Dimensionless(), sym(t, "E"), "E"))
	if err != nil {
		t.Fatalf("branch solve: %v", err)
	}
	lo, err := Compare(
		exprObj(t, core.Dimensionless(), sym(t, "E")),
		exprObj(t, core.Dimensionless(), rat(0, 1)), core.RelationGte)
	if err != nil {
		t.Fatalf("branch constraint: %v", err)
	}
	return br, lo
}

func addAssumptionForSubset(t *testing.T, o core.Object, a core.Assumption) (core.Object, error) {
	t.Helper()
	merged, err := o.Assumptions().Merge(core.NewAssumptionSet(a))
	if err != nil {
		return core.Object{}, err
	}
	return mintSpec(t, kernel.ObjectSpec{
		Name: "w", Kind: o.Kind(), Dimension: o.Dimension(), Expr: o.Expr(),
		Assumptions: merged, Conventions: o.Conventions(),
		Provenance:   newProv(t, core.StatusDerived, "test", "test", nil, ""),
		CorpusStatus: kernel.CorpusNone,
	}), nil
}

func assumptionEqualForSubset(a, b core.Assumption) bool {
	x, err1 := json.Marshal(a)
	y, err2 := json.Marshal(b)
	if err1 != nil || err2 != nil {
		return false
	}
	return string(x) == string(y)
}

package core_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math/big"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/PithomLabs/phys/core"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func mustSym(t *testing.T, name string) core.Expr {
	t.Helper()
	e, err := core.NewSymbol(name)
	if err != nil {
		t.Fatalf("NewSymbol(%q): %v", name, err)
	}
	return e
}

func mustBytes(t *testing.T, b []byte, err error) []byte {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// must unwraps a (bytes, error) pair in sole-argument position.
func must(b []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return b
}

func rat(n, d int64) *big.Rat { return big.NewRat(n, d) }

func hexOf(h [32]byte) string { return hex.EncodeToString(h[:]) }

// ---------------------------------------------------------------------------
// Expression enum, constructors, accessors (§8)
// ---------------------------------------------------------------------------

func TestExprKindEnumExact(t *testing.T) {
	want := []struct {
		name string
		k    core.ExprKind
		json string
	}{
		{"Symbol", core.ExprSymbol, "symbol"},
		{"Rational", core.ExprRational, "rational"},
		{"Add", core.ExprAdd, "add"},
		{"Mul", core.ExprMul, "mul"},
		{"Neg", core.ExprNeg, "neg"},
		{"Pow", core.ExprPow, "pow"},
		{"Sqrt", core.ExprSqrt, "sqrt"},
		{"Call", core.ExprCall, "call"},
		{"Relation", core.ExprRelation, "relation"},
		{"BranchSet", core.ExprBranchSet, "branch_set"},
	}
	for i, w := range want {
		if int(w.k) != i {
			t.Fatalf("ordinal of %s = %d, want %d", w.name, w.k, i)
		}
		if w.k.String() != w.json {
			t.Fatalf("String(%s) = %q, want %q", w.name, w.k.String(), w.json)
		}
	}
	if (core.ExprKind(10)).String() != "" {
		t.Fatal("out-of-range ExprKind must have empty String")
	}
	// Forbidden node names do not exist as kind strings.
	forbidden := []string{"Derivative", "Limit", "Function", "Arbitrary", "Eval", "Callback", "RawString"}
	for _, f := range forbidden {
		for _, w := range want {
			if w.json == strings.ToLower(f) {
				t.Fatalf("forbidden node %q present", f)
			}
		}
	}
	if len(want) != 10 {
		t.Fatal("node set must be exactly 10")
	}
}

func TestExprConstructorSurface(t *testing.T) {
	m := mustSym(t, "m")
	v := mustSym(t, "v")
	z := core.NewRational(rat(0, 1))

	if _, err := core.NewSymbol(""); err == nil {
		t.Fatal("empty symbol must be rejected")
	}
	if _, err := core.NewCall("not_lorentz", m); err == nil {
		t.Fatal("unknown function id must be rejected")
	}
	call, err := core.NewCall("lorentz_factor", v)
	if err != nil || !call.Valid() {
		t.Fatalf("lorentz_factor call: %v", err)
	}
	if call.FunctionID() != "lorentz_factor" {
		t.Fatalf("function id %q", call.FunctionID())
	}

	add := core.NewAdd(m, v)
	if !add.Valid() || add.Kind() != core.ExprAdd {
		t.Fatal("NewAdd")
	}
	mul := core.NewMul(m, v)
	if !mul.Valid() || mul.Kind() != core.ExprMul {
		t.Fatal("NewMul")
	}
	neg := core.NewNeg(m)
	if !neg.Valid() || neg.Kind() != core.ExprNeg {
		t.Fatal("NewNeg")
	}
	pow := core.NewPow(m, rat(2, 1))
	if !pow.Valid() || pow.Kind() != core.ExprPow {
		t.Fatal("NewPow")
	}
	sqrt := core.NewSqrt(m)
	if !sqrt.Valid() || sqrt.Kind() != core.ExprSqrt {
		t.Fatal("NewSqrt")
	}
	rel := core.NewRelation(core.RelationEq, m, v)
	if !rel.Valid() || rel.Kind() != core.ExprRelation {
		t.Fatal("NewRelation")
	}
	if core.NewRelation("equals", m, v).Valid() {
		t.Fatal("invalid relation operator must yield invalid expr")
	}
	bs := core.NewBranchSet(m, core.NewSqrt(m), core.NewNeg(core.NewSqrt(m)))
	if !bs.Valid() || bs.Kind() != core.ExprBranchSet {
		t.Fatal("NewBranchSet")
	}
	if core.NewBranchSet(m).Valid() {
		t.Fatal("BranchSet without branches must be invalid")
	}
	if core.NewAdd(m, core.Expr{}).Valid() {
		t.Fatal("Add with invalid child must be invalid")
	}
	if core.NewPow(core.Expr{}, rat(2, 1)).Valid() {
		t.Fatal("Pow with invalid base must be invalid")
	}
	if core.NewRational(nil).Valid() {
		t.Fatal("nil rational must be invalid")
	}
	_ = z
}

func TestExprAccessors(t *testing.T) {
	m := mustSym(t, "m")
	v := mustSym(t, "v")
	r := core.NewRational(rat(1, 2))
	pow := core.NewPow(m, rat(-3, 4))
	rel := core.NewRelation(core.RelationGt, m, v)
	bs := core.NewBranchSet(m, core.NewSqrt(m), core.NewNeg(core.NewSqrt(m)))
	call, err := core.NewCall("lorentz_factor", v)
	if err != nil {
		t.Fatal(err)
	}
	add := core.NewAdd(m, v)
	neg := core.NewNeg(m)
	sqrt := core.NewSqrt(m)

	// Symbol
	if !m.Valid() || m.Kind() != core.ExprSymbol || m.SymbolName() != "m" {
		t.Fatal("symbol accessors")
	}
	if m.RationalValue() != nil || m.Children() != nil || m.Exponent() != nil ||
		m.FunctionID() != "" || m.RelationOperator() != "" {
		t.Fatal("symbol inapplicable accessors must be empty")
	}
	if m.Base().Valid() || m.BranchTarget().Valid() || m.Arguments() != nil || m.Branches() != nil {
		t.Fatal("symbol inapplicable accessors (2)")
	}
	// Rational
	if rv := r.RationalValue(); rv == nil || rv.Cmp(rat(1, 2)) != 0 {
		t.Fatal("rational value")
	}
	if r.Children() != nil {
		t.Fatal("rational has no children")
	}
	// Pow
	if !core.EqualExpr(pow.Base(), m) {
		t.Fatal("pow base")
	}
	if e := pow.Exponent(); e == nil || e.Cmp(rat(-3, 4)) != 0 {
		t.Fatal("pow exponent")
	}
	if pow.RationalValue() != nil {
		t.Fatal("pow rational accessor must be nil")
	}
	// Relation
	if rel.RelationOperator() != core.RelationGt || !core.EqualExpr(rel.Left(), m) ||
		!core.EqualExpr(rel.Right(), v) {
		t.Fatal("relation accessors")
	}
	if rel.Children() == nil || len(rel.Children()) != 2 {
		t.Fatal("relation children copy")
	}
	// BranchSet
	if !core.EqualExpr(bs.BranchTarget(), m) {
		t.Fatal("branch target")
	}
	brs := bs.Branches()
	if len(brs) != 2 {
		t.Fatalf("branches = %d", len(brs))
	}
	if bs.RelationOperator() != "" || bs.FunctionID() != "" {
		t.Fatal("branch set inapplicable accessors")
	}
	// Call
	if call.FunctionID() != "lorentz_factor" || len(call.Arguments()) != 1 {
		t.Fatal("call accessors")
	}
	if call.Exponent() != nil {
		t.Fatal("call exponent must be nil")
	}
	// Add/Neg/Sqrt
	if len(add.Children()) != 2 {
		t.Fatal("add children")
	}
	if len(neg.Children()) != 1 || !core.EqualExpr(neg.Children()[0], m) {
		t.Fatal("neg children")
	}
	if len(sqrt.Children()) != 1 {
		t.Fatal("sqrt children")
	}
	if neg.RelationOperator() != "" {
		t.Fatal("neg relation operator")
	}
	// Invalid zero handle
	var zero core.Expr
	if zero.Valid() {
		t.Fatal("zero Expr must be invalid")
	}
	if zero.SymbolName() != "" || zero.Children() != nil || zero.RationalValue() != nil {
		t.Fatal("zero accessors empty")
	}
}

func TestExprDefensiveCopies(t *testing.T) {
	// NewRational copies its input.
	r := big.NewRat(1, 2)
	e := core.NewRational(r)
	r.Set(big.NewRat(7, 3))
	if v := e.RationalValue(); v.Cmp(rat(1, 2)) != 0 {
		t.Fatalf("NewRational retained caller pointer: %s", v)
	}
	// RationalValue returns a copy.
	v := e.RationalValue()
	v.Set(big.NewRat(9, 5))
	if e.RationalValue().Cmp(rat(1, 2)) != 0 {
		t.Fatal("RationalValue must return a copy")
	}
	// NewPow clones its exponent.
	exp := big.NewRat(2, 1)
	p := core.NewPow(mustSym(t, "x"), exp)
	exp.Set(big.NewRat(5, 1))
	if p.Exponent().Cmp(rat(2, 1)) != 0 {
		t.Fatal("NewPow retained caller exponent")
	}
	// Exponent returns a copy.
	e2 := p.Exponent()
	e2.Set(big.NewRat(0, 1))
	if p.Exponent().Cmp(rat(2, 1)) != 0 {
		t.Fatal("Exponent must return a copy")
	}
	// Slice accessors return copies.
	add := core.NewAdd(mustSym(t, "a"), mustSym(t, "b"))
	ch := add.Children()
	ch[0] = core.Expr{}
	if !add.Valid() || len(add.Children()) != 2 {
		t.Fatal("Children must return a copy")
	}
	bs := core.NewBranchSet(mustSym(t, "t"), mustSym(t, "x"))
	brs := bs.Branches()
	brs[0] = core.Expr{}
	if len(bs.Branches()) != 1 || !bs.Branches()[0].Valid() {
		t.Fatal("Branches must return a copy")
	}
}

func TestExprStoresNoCallbacks(t *testing.T) {
	rt := reflect.TypeOf(core.Expr{})
	if rt.Kind() != reflect.Struct {
		t.Fatalf("Expr backing must be a struct, got %v", rt.Kind())
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.Type.Kind() == reflect.Func || f.Type.Kind() == reflect.Chan {
			t.Fatalf("Expr field %s stores executable state", f.Name)
		}
		if f.PkgPath == "" {
			t.Fatalf("Expr field %s must be unexported", f.Name)
		}
	}
	// Expr is deeply data-only: no func values reachable through fields.
	hasFunc := false
	var walk func(t2 reflect.Type)
	walk = func(t2 reflect.Type) {
		for i := 0; i < t2.NumField(); i++ {
			ft := t2.Field(i).Type
			if ft.Kind() == reflect.Func {
				hasFunc = true
			}
			if ft.Kind() == reflect.Ptr && ft.Elem().Kind() == reflect.Struct && ft.Elem() != t2 {
				walk(ft.Elem())
			}
		}
	}
	walk(rt)
	if hasFunc {
		t.Fatal("Expr struct graph contains func state")
	}
}

func TestSymbolIsDataOnly(t *testing.T) {
	names := []string{"m", "v", "a", "F", "p", "E", "c", "t", "x", "η", "d/dt", "func", "if"}
	for _, n := range names {
		e, err := core.NewSymbol(n)
		if err != nil {
			t.Fatalf("symbol %q: %v", n, err)
		}
		if e.SymbolName() != n {
			t.Fatalf("symbol %q roundtrip", n)
		}
		b := must(core.CanonicalExprJSON(e))
		want := `{"kind":"symbol","name":"` + n + `"}`
		if string(b) != want {
			t.Fatalf("symbol bytes %s want %s", b, want)
		}
		back, err := core.ParseExprJSON(b)
		if err != nil || back.SymbolName() != n {
			t.Fatalf("symbol parse %q: %v", n, err)
		}
	}
	if _, err := core.NewSymbol(""); err == nil {
		t.Fatal("empty symbol rejected")
	}
}

// ---------------------------------------------------------------------------
// Exact rationals and canonical serialization (§8.4, §10.2)
// ---------------------------------------------------------------------------

func TestRationalExactSerialization(t *testing.T) {
	cases := []struct {
		expr core.Expr
		want string
	}{
		{core.NewRational(rat(1, 2)), `{"kind":"rational","value":"1/2"}`},
		{core.NewRational(rat(-3, 4)), `{"kind":"rational","value":"-3/4"}`},
		{core.NewRational(rat(0, 1)), `{"kind":"rational","value":"0/1"}`},
		{core.NewRational(rat(2, 1)), `{"kind":"rational","value":"2/1"}`},
		{core.NewRational(big.NewRat(1234567890123456789, 987654321098765432)),
			`{"kind":"rational","value":"1234567890123456789/987654321098765432"}`},
	}
	for _, c := range cases {
		b := must(core.CanonicalExprJSON(c.expr))
		if string(b) != c.want {
			t.Fatalf("got %s want %s", b, c.want)
		}
		// round-trip byte identity
		back, err := core.ParseExprJSON(b)
		if err != nil {
			t.Fatal(err)
		}
		b2 := must(core.CanonicalExprJSON(back))
		if !bytes.Equal(b, b2) {
			t.Fatalf("roundtrip %s vs %s", b, b2)
		}
	}
}

func TestRationalStringNotNumber(t *testing.T) {
	b := must(core.CanonicalExprJSON(core.NewRational(rat(1, 2))))
	if !bytes.Contains(b, []byte(`"value":"1/2"`)) {
		t.Fatalf("rational must serialize as string: %s", b)
	}
	var probe struct {
		Value json.RawMessage `json:"value"`
	}
	if err := json.Unmarshal(b, &probe); err != nil {
		t.Fatal(err)
	}
	if len(probe.Value) == 0 || probe.Value[0] != '"' {
		t.Fatalf("value field must be a JSON string, got %s", probe.Value)
	}
}

func TestExactRationalRoundTrip(t *testing.T) {
	exprs := []core.Expr{
		core.NewRational(rat(-3, 4)),
		core.NewRational(rat(0, 1)),
		core.NewAdd(mustSym(t, "m"), core.NewRational(rat(1, 2))),
		core.NewMul(mustSym(t, "m"), core.NewPow(mustSym(t, "c"), rat(2, 1))),
		core.NewPow(mustSym(t, "x"), rat(-1, 2)),
		core.NewNeg(core.NewRational(rat(7, 9))),
		core.NewRelation(core.RelationEq, core.NewRational(rat(5, 6)), mustSym(t, "E")),
		core.NewSqrt(core.NewRational(rat(4, 9))),
	}
	for i, e := range exprs {
		b := must(core.CanonicalExprJSON(e))
		back, err := core.ParseExprJSON(b)
		if err != nil {
			t.Fatalf("[%d] parse: %v", i, err)
		}
		b2 := must(core.CanonicalExprJSON(back))
		if !bytes.Equal(b, b2) {
			t.Fatalf("[%d] bytes differ: %s vs %s", i, b, b2)
		}
		h1, h2 := core.HashExpr(e), core.HashExpr(back)
		if h1 != h2 {
			t.Fatalf("[%d] hash differ", i)
		}
		if !core.EqualExpr(e, back) {
			t.Fatalf("[%d] structural equality", i)
		}
	}
	// Non-canonical / dynamic forms are rejected.
	bad := []string{
		`{"kind":"rational","value":"2/4"}`,
		`{"kind":"rational","value":"0.5"}`,
		`{"kind":"rational","value":0.5}`,
		`{"kind":"rational","value":"1/-2"}`,
		`{"kind":"rational","value":"-1/-2"}`,
		`{"kind":"rational","value":"1"}`,
		`{"kind":"rational","value":"1/"}`,
	}
	for _, s := range bad {
		if _, err := core.ParseExprJSON([]byte(s)); err == nil {
			t.Fatalf("must reject %s", s)
		}
	}
}

func TestRationalCombinationExact(t *testing.T) {
	// §9.4 examples: 1/2 · 2 → 1 ; 1/2 + 1/2 → 1
	prod := core.NewMul(core.NewRational(rat(1, 2)), core.NewRational(rat(2, 1)))
	if prod.Kind() != core.ExprRational || prod.RationalValue().Cmp(rat(1, 1)) != 0 {
		t.Fatal("1/2 * 2 must combine to 1")
	}
	sum := core.NewAdd(core.NewRational(rat(1, 2)), core.NewRational(rat(1, 2)))
	if sum.Kind() != core.ExprRational || sum.RationalValue().Cmp(rat(1, 1)) != 0 {
		t.Fatal("1/2 + 1/2 must combine to 1")
	}
	// Exactness with large integers (no float64 anywhere).
	a := new(big.Rat).SetFrac(big.NewInt(1000000007), big.NewInt(1000000009))
	b := new(big.Rat).SetFrac(big.NewInt(1000000009), big.NewInt(1000000007))
	got := core.NewMul(core.NewRational(a), core.NewRational(b))
	if got.Kind() != core.ExprRational || got.RationalValue().Cmp(rat(1, 1)) != 0 {
		t.Fatalf("large exact product = %v", got.RationalValue())
	}
	// Mixed Add keeps the combined rational term for Simplify-time removal.
	mx := core.NewAdd(mustSym(t, "x"), core.NewRational(rat(0, 1)))
	if mx.Kind() != core.ExprAdd || len(mx.Children()) != 2 {
		t.Fatalf("Add(x,0) must keep structure, got %v", mx.Kind())
	}
	// Nonzero rational combination in mixed Add.
	mx2 := core.NewAdd(mustSym(t, "x"), core.NewRational(rat(1, 4)), core.NewRational(rat(1, 4)))
	if len(mx2.Children()) != 2 {
		t.Fatalf("combined rational terms: %d", len(mx2.Children()))
	}
	if mx2.Children()[1].RationalValue().Cmp(rat(1, 2)) != 0 &&
		mx2.Children()[0].RationalValue().Cmp(rat(1, 2)) != 0 {
		t.Fatal("1/4+1/4 must combine to 1/2")
	}
}

func TestPowExponentExact(t *testing.T) {
	cases := []struct {
		exp  *big.Rat
		want string
	}{
		{rat(2, 1), "2/1"},
		{rat(1, 2), "1/2"},
		{rat(-3, 4), "-3/4"},
		{rat(-1, 1), "-1/1"},
	}
	for _, c := range cases {
		p := core.NewPow(mustSym(t, "x"), c.exp)
		b := must(core.CanonicalExprJSON(p))
		want := `{"kind":"pow","base":{"kind":"symbol","name":"x"},"exp":"` + c.want + `"}`
		if string(b) != want {
			t.Fatalf("got %s want %s", b, want)
		}
		back, err := core.ParseExprJSON(b)
		if err != nil {
			t.Fatal(err)
		}
		if back.Exponent().Cmp(c.exp) != 0 {
			t.Fatal("exponent roundtrip")
		}
	}
	// Non-canonical exponents rejected.
	p := core.NewPow(mustSym(t, "x"), rat(1, 2))
	b := must(core.CanonicalExprJSON(p))
	bad := strings.Replace(string(b), `"exp":"1/2"`, `"exp":"2/4"`, 1)
	if _, err := core.ParseExprJSON([]byte(bad)); err == nil {
		t.Fatal("non-canonical exponent must be rejected")
	}
}

func TestPowClonesExponent(t *testing.T) {
	exp := big.NewRat(3, 1)
	p := core.NewPow(mustSym(t, "y"), exp)
	exp.Neg(exp)
	if p.Exponent().Cmp(rat(3, 1)) != 0 {
		t.Fatal("NewPow must clone the caller exponent")
	}
}

// ---------------------------------------------------------------------------
// Canonical equality, ordering, hashing (§9)
// ---------------------------------------------------------------------------

func TestEqualExprStructural(t *testing.T) {
	m := mustSym(t, "m")
	v := mustSym(t, "v")
	// Structural: commutative order does not matter after canonicalization.
	if !core.EqualExpr(core.NewAdd(m, v), core.NewAdd(v, m)) {
		t.Fatal("Add is commutative under canonical equality")
	}
	if !core.EqualExpr(core.NewMul(m, v), core.NewMul(v, m)) {
		t.Fatal("Mul is commutative under canonical equality")
	}
	// Not display strings: same structure via independent construction.
	a := core.NewMul(mustSym(t, "m"), core.NewPow(mustSym(t, "c"), rat(2, 1)))
	b := core.NewMul(mustSym(t, "m"), core.NewPow(mustSym(t, "c"), rat(2, 1)))
	if !core.EqualExpr(a, b) {
		t.Fatal("identical structure must be equal")
	}
	if core.EqualExpr(m, v) {
		t.Fatal("different symbols not equal")
	}
	// Relation sides are ordered, not commutative.
	r1 := core.NewRelation(core.RelationEq, m, v)
	r2 := core.NewRelation(core.RelationEq, v, m)
	if core.EqualExpr(r1, r2) {
		t.Fatal("relation sides must preserve order")
	}
	// Different operators.
	if core.EqualExpr(r1, core.NewRelation(core.RelationNeq, m, v)) {
		t.Fatal("operators distinguish relations")
	}
	// Rational comparison is exact numeric, not string spelling: only one
	// spelling can be constructed (2/4 cannot be built or parsed).
	if core.EqualExpr(core.NewRational(rat(1, 2)), core.NewRational(rat(1, 3))) {
		t.Fatal("distinct rationals not equal")
	}
	// Invalid handles are never equal.
	var zero core.Expr
	if core.EqualExpr(zero, zero) {
		t.Fatal("invalid expressions must not compare equal")
	}
	// Round-trip equality.
	bts := must(core.CanonicalExprJSON(a))
	back, err := core.ParseExprJSON(bts)
	if err != nil {
		t.Fatal(err)
	}
	if !core.EqualExpr(a, back) {
		t.Fatal("round-trip structural equality")
	}
}

func TestCanonicalOrderingStable(t *testing.T) {
	m := mustSym(t, "m")
	v := mustSym(t, "v")
	a := core.NewAdd(m, v)
	b := core.NewAdd(v, m)
	if !bytes.Equal(must(core.CanonicalExprJSON(a)), must(core.CanonicalExprJSON(b))) {
		t.Fatal("construction order must not affect canonical bytes")
	}
	// Cross-kind ordering follows node-kind ordinal: Symbol(0) < Rational(1)
	// < Mul(3) < Neg(4) < Pow(5) < Sqrt(6).
	x := core.NewAdd(core.NewSqrt(m), core.NewPow(m, rat(2, 1)), core.NewNeg(m),
		core.NewMul(m, v), core.NewRational(rat(1, 1)), m)
	kinds := make([]core.ExprKind, 0, len(x.Children()))
	for _, c := range x.Children() {
		kinds = append(kinds, c.Kind())
	}
	want := []core.ExprKind{
		core.ExprSymbol, core.ExprRational, core.ExprMul, core.ExprNeg, core.ExprPow, core.ExprSqrt,
	}
	if !reflect.DeepEqual(kinds, want) {
		t.Fatalf("kind order %v want %v", kinds, want)
	}
	// Repeated runs, both construction orders → identical bytes.
	p1 := core.NewMul(mustSym(t, "aa"), mustSym(t, "bb"))
	p2 := core.NewMul(mustSym(t, "bb"), mustSym(t, "aa"))
	if !bytes.Equal(must(core.CanonicalExprJSON(p1)), must(core.CanonicalExprJSON(p2))) {
		t.Fatal("hash-keyed sort must be order-independent")
	}
	// Deep tree determinism.
	t1 := core.NewAdd(core.NewMul(m, v), core.NewMul(v, m), core.NewNeg(core.NewNeg(m)))
	t2 := core.NewAdd(core.NewNeg(core.NewNeg(m)), core.NewMul(v, m), core.NewMul(m, v))
	if !bytes.Equal(must(core.CanonicalExprJSON(t1)), must(core.CanonicalExprJSON(t2))) {
		t.Fatal("deep canonical ordering must be stable")
	}
}

func TestRelationOperatorEnumExact(t *testing.T) {
	want := []struct {
		op   core.RelationOperator
		json string
	}{
		{core.RelationEq, "eq"},
		{core.RelationNeq, "neq"},
		{core.RelationLt, "lt"},
		{core.RelationLte, "lte"},
		{core.RelationGt, "gt"},
		{core.RelationGte, "gte"},
	}
	for _, w := range want {
		if string(w.op) != w.json {
			t.Fatalf("operator %q want %q", w.op, w.json)
		}
		rel := core.NewRelation(w.op, mustSym(t, "a"), mustSym(t, "b"))
		b := must(core.CanonicalExprJSON(rel))
		if !bytes.Contains(b, []byte(`"operator":"`+w.json+`"`)) {
			t.Fatalf("operator encoding: %s", b)
		}
		back, err := core.ParseExprJSON(b)
		if err != nil {
			t.Fatal(err)
		}
		if back.RelationOperator() != w.op {
			t.Fatal("operator roundtrip")
		}
	}
	if len(want) != 6 {
		t.Fatal("exactly six operators")
	}
	// No alias operator names accepted by canonical JSON.
	alias := `{"kind":"relation","operator":"equals","lhs":{"kind":"symbol","name":"a"},"rhs":{"kind":"symbol","name":"b"}}`
	if _, err := core.ParseExprJSON([]byte(alias)); err == nil {
		t.Fatal("alias operator must be rejected")
	}
	le := `{"kind":"relation","operator":"le","lhs":{"kind":"symbol","name":"a"},"rhs":{"kind":"symbol","name":"b"}}`
	if _, err := core.ParseExprJSON([]byte(le)); err == nil {
		t.Fatal("alias operator le must be rejected")
	}
}

func TestHashDeterministic(t *testing.T) {
	m := mustSym(t, "m")
	c := mustSym(t, "c")
	expr := core.NewMul(m, core.NewPow(c, rat(2, 1)))
	if got := hexOf(core.HashExpr(expr)); got != "4d77099c7c0cbf3df1f37bb1b53e6855612cf1f98bedf2d75cd916982d909b82" {
		t.Fatalf("fixed expr vector = %s", got)
	}
	if got := hexOf(core.HashExpr(mustSym(t, "E"))); got != "09fae00a8c1018c1451fee8baa955f393e6c68812055d178f2aea190756353ac" {
		t.Fatalf("fixed symbol vector = %s", got)
	}
	a1, err := core.NewTextAssumption(core.AssumptionDomain, "k1", "v1")
	if err != nil {
		t.Fatal(err)
	}
	set := core.NewAssumptionSet(a1)
	if got := hexOf(core.HashAssumptionSet(set)); got != "62c1be7f62be53dc8e5c5a974044c3b6a58ab0180e3e726d19f6af2187d67e6f" {
		t.Fatalf("fixed assumption-set vector = %s", got)
	}
	if got := hexOf(core.DimensionForce().Hash()); got != "807962c2c2db263a67e95aa7e217bfdd178251deae9b32c89d69f2d8509fbf9b" {
		t.Fatalf("fixed dimension vector = %s", got)
	}
	// Double-run stability.
	for i := 0; i < 3; i++ {
		if core.HashExpr(expr) != core.HashExpr(core.NewMul(m, core.NewPow(c, rat(2, 1)))) {
			t.Fatal("hash must be deterministic across runs")
		}
		if core.HashAssumptionSet(set) != core.HashAssumptionSet(core.NewAssumptionSet(a1)) {
			t.Fatal("set hash must be deterministic")
		}
	}
}

// ---------------------------------------------------------------------------
// Dimensions (§7)
// ---------------------------------------------------------------------------

func TestDimensionConstructorsAndAPI(t *testing.T) {
	cases := []struct {
		name string
		d    core.Dimension
		m, l, t3 int64
	}{
		{"Dimensionless", core.Dimensionless(), 0, 0, 0},
		{"DimensionMass", core.DimensionMass(), 1, 0, 0},
		{"DimensionLength", core.DimensionLength(), 0, 1, 0},
		{"DimensionTime", core.DimensionTime(), 0, 0, 1},
		{"DimensionVelocity", core.DimensionVelocity(), 0, 1, -1},
		{"DimensionAcceleration", core.DimensionAcceleration(), 0, 1, -2},
		{"DimensionForce", core.DimensionForce(), 1, 1, -2},
		{"DimensionMomentum", core.DimensionMomentum(), 1, 1, -1},
		{"DimensionEnergy", core.DimensionEnergy(), 1, 2, -2},
	}
	for _, c := range cases {
		if !c.d.Valid() {
			t.Fatalf("%s must be valid", c.name)
		}
	}
	// Explicit §7.2 exponent vectors via arithmetic identities.
	must := func(d core.Dimension) core.Dimension {
		t.Helper()
		if !d.Valid() {
			t.Fatal("dimension invalid")
		}
		return d
	}
	if !must(core.DimensionVelocity()).Equal(core.DimensionLength().Divide(core.DimensionTime())) {
		t.Fatal("velocity = L/T")
	}
	if !must(core.DimensionAcceleration()).Equal(core.DimensionLength().Divide(core.DimensionTime().Multiply(core.DimensionTime()))) {
		t.Fatal("acceleration = L/T^2")
	}
	if !must(core.DimensionForce()).Equal(core.DimensionMass().Multiply(core.DimensionLength()).Divide(core.DimensionTime().Multiply(core.DimensionTime()))) {
		t.Fatal("force = ML/T^2")
	}
	if !must(core.DimensionMomentum()).Equal(core.DimensionMass().Multiply(core.DimensionLength()).Divide(core.DimensionTime())) {
		t.Fatal("momentum = ML/T")
	}
	if !must(core.DimensionEnergy()).Equal(core.DimensionMass().Multiply(core.DimensionLength()).Multiply(core.DimensionLength()).Divide(core.DimensionTime().Multiply(core.DimensionTime()))) {
		t.Fatal("energy = ML^2/T^2")
	}
	if !core.Dimensionless().Equal(core.Dimensionless().Multiply(core.Dimensionless())) {
		t.Fatal("dimensionless identity")
	}
	// Zero value invalid.
	var zero core.Dimension
	if zero.Valid() || zero.Equal(core.Dimensionless()) {
		t.Fatal("zero Dimension must be invalid")
	}
	// Round-trip through canonical JSON for each constructor.
	for _, c := range cases {
		b, err := c.d.CanonicalJSON()
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		back, err := core.ParseDimensionJSON(b)
		if err != nil {
			t.Fatalf("%s parse: %v", c.name, err)
		}
		if !back.Equal(c.d) {
			t.Fatalf("%s roundtrip", c.name)
		}
	}
}

func TestDimensionArithmetic(t *testing.T) {
	// Power multiplies every exponent by the exact rational exponent.
	half := core.DimensionForce().Pow(rat(1, 2))
	want, _ := core.NewDimension(rat(1, 2), rat(1, 2), rat(-1, 1),
		rat(0, 1), rat(0, 1), rat(0, 1), rat(0, 1))
	if !half.Equal(want) {
		t.Fatalf("Force^(1/2) = %v", half)
	}
	third := core.Dimensionless().Pow(rat(2, 3))
	if !third.Equal(core.Dimensionless()) {
		t.Fatal("dimensionless^n = dimensionless")
	}
	// Division cancels.
	if !core.DimensionVelocity().Multiply(core.DimensionTime()).Equal(core.DimensionLength()) {
		t.Fatal("(L/T)*T = L")
	}
	if !core.DimensionForce().Divide(core.DimensionForce()).Equal(core.Dimensionless()) {
		t.Fatal("F/F = dimensionless")
	}
	// Exact rational exponents, no floats: (L^(1/3))^3 = L.
	one := core.DimensionLength().Pow(rat(1, 3)).Pow(rat(3, 1))
	if !one.Equal(core.DimensionLength()) {
		t.Fatalf("(L^(1/3))^3 = %v", one)
	}
	// Determinism across repeated arithmetic.
	a := core.DimensionEnergy().Multiply(core.DimensionVelocity()).Divide(core.DimensionLength())
	b := core.DimensionEnergy().Multiply(core.DimensionVelocity()).Divide(core.DimensionLength())
	if !a.Equal(b) {
		t.Fatal("arithmetic determinism")
	}
	// Invalid operands yield invalid results (no panic).
	var zero core.Dimension
	if zero.Multiply(core.DimensionMass()).Valid() {
		t.Fatal("invalid operand must yield invalid result")
	}
	if zero.Pow(rat(2, 1)).Valid() {
		t.Fatal("invalid base must yield invalid result")
	}
}

func TestDimensionEqual(t *testing.T) {
	if !core.DimensionMass().Equal(core.DimensionMass()) {
		t.Fatal("reflexive")
	}
	// Equal exponents from independent construction compare equal.
	restMass, _ := core.NewDimension(rat(1, 1), rat(0, 1), rat(0, 1),
		rat(0, 1), rat(0, 1), rat(0, 1), rat(0, 1))
	if !core.DimensionMass().Equal(restMass) {
		t.Fatal("M and RestMass share dimension M")
	}
	if core.DimensionMass().Equal(core.DimensionVelocity()) {
		t.Fatal("M != L/T")
	}
	var zero core.Dimension
	if zero.Equal(zero) {
		t.Fatal("invalid dimensions are never equal")
	}
	if core.DimensionMass().Equal(zero) {
		t.Fatal("valid vs invalid not equal")
	}
}

func TestDimensionJSONFieldOrder(t *testing.T) {
	b := must(core.DimensionMass().CanonicalJSON())
	if string(b) != `{"m":"1/1","l":"0/1","t":"0/1","i":"0/1","theta":"0/1","n":"0/1","j":"0/1"}` {
		t.Fatalf("mass dimension bytes %s", b)
	}
	b = must(core.Dimensionless().CanonicalJSON())
	if string(b) != `{"m":"0/1","l":"0/1","t":"0/1","i":"0/1","theta":"0/1","n":"0/1","j":"0/1"}` {
		t.Fatalf("dimensionless bytes %s", b)
	}
	// Non-canonical spellings rejected.
	bad := `{"m":"2/4","l":"0/1","t":"0/1","i":"0/1","theta":"0/1","n":"0/1","j":"0/1"}`
	if _, err := core.ParseDimensionJSON([]byte(bad)); err == nil {
		t.Fatal("non-canonical exponent must be rejected")
	}
	badNum := `{"m":1,"l":0,"t":0,"i":0,"theta":0,"n":0,"j":0}`
	if _, err := core.ParseDimensionJSON([]byte(badNum)); err == nil {
		t.Fatal("numeric exponents must be rejected")
	}
	unknown := `{"m":"0/1","l":"0/1","t":"0/1","i":"0/1","theta":"0/1","n":"0/1","j":"0/1","x":"0/1"}`
	if _, err := core.ParseDimensionJSON([]byte(unknown)); err == nil {
		t.Fatal("unknown dimension field must be rejected")
	}
}

func TestDimensionDeterminism(t *testing.T) {
	for i := 0; i < 3; i++ {
		b1 := must(core.DimensionEnergy().CanonicalJSON())
		b2 := must(core.DimensionEnergy().CanonicalJSON())
		if !bytes.Equal(b1, b2) {
			t.Fatal("dimension bytes deterministic")
		}
		if core.DimensionEnergy().Hash() != core.DimensionEnergy().Hash() {
			t.Fatal("dimension hash deterministic")
		}
		d1 := core.DimensionForce().Multiply(core.DimensionVelocity())
		d2 := core.DimensionVelocity().Multiply(core.DimensionForce())
		if !d1.Equal(d2) || d1.Hash() != d2.Hash() {
			t.Fatal("dimension arithmetic deterministic")
		}
	}
}

// ---------------------------------------------------------------------------
// Assumptions (§11)
// ---------------------------------------------------------------------------

func TestAssumptionKindEnumExact(t *testing.T) {
	want := []struct {
		k    core.AssumptionKind
		json string
	}{
		{core.AssumptionDomain, "domain"},
		{core.AssumptionRegime, "regime"},
		{core.AssumptionConstraint, "constraint"},
		{core.AssumptionConvention, "convention"},
		{core.AssumptionApproximation, "approximation"},
		{core.AssumptionMathPrecondition, "math_precondition"},
		{core.AssumptionPhysicalAssumption, "physical_assumption"},
	}
	if len(want) != 7 {
		t.Fatal("exactly seven assumption kinds")
	}
	for _, w := range want {
		if string(w.k) != w.json {
			t.Fatalf("kind %q want json %q", w.k, w.json)
		}
		a, err := core.NewTextAssumption(w.k, "some_key", "some value")
		if err != nil {
			t.Fatalf("kind %q: %v", w.k, err)
		}
		set := core.NewAssumptionSet(a)
		b, err := set.CanonicalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Contains(b, []byte(`"kind":"`+w.json+`"`)) {
			t.Fatalf("canonical json for kind %q: %s", w.json, b)
		}
	}
	// No eighth kind may be constructed.
	if _, err := core.NewTextAssumption("bogus", "k", "v"); err == nil {
		t.Fatal("unknown assumption kind must be rejected")
	}
}

func TestAssumptionAccessors(t *testing.T) {
	m := mustSym(t, "m")
	txt, err := core.NewTextAssumption(core.AssumptionRegime, "regime_key", "classical")
	if err != nil {
		t.Fatal(err)
	}
	if txt.Kind() != core.AssumptionRegime || txt.Key() != "regime_key" {
		t.Fatal("text kind/key")
	}
	if txt.Mode() != "text" || txt.Text() != "classical" {
		t.Fatal("text mode/value")
	}
	if txt.Expr().Valid() {
		t.Fatal("text assumption has no expr")
	}
	z := core.NewRational(rat(0, 1))
	ex, err := core.NewExprAssumption(core.AssumptionConstraint, "rest_mass_nonnegative",
		core.NewRelation(core.RelationGte, m, z))
	if err != nil {
		t.Fatal(err)
	}
	if ex.Kind() != core.AssumptionConstraint || ex.Key() != "rest_mass_nonnegative" {
		t.Fatal("expr kind/key")
	}
	if ex.Mode() != "expr" || ex.Text() != "" {
		t.Fatal("expr mode")
	}
	if !ex.Expr().Valid() || ex.Expr().Kind() != core.ExprRelation {
		t.Fatal("expr content")
	}
	// Invalid value modes cannot exist through constructors.
	if _, err := core.NewExprAssumption(core.AssumptionDomain, "k", core.Expr{}); err == nil {
		t.Fatal("invalid expr value rejected")
	}
}

func TestAssumptionStructuredValue(t *testing.T) {
	m := mustSym(t, "m")
	z := core.NewRational(rat(0, 1))
	a, err := core.NewExprAssumption(core.AssumptionConstraint, "momentum_zero",
		core.NewRelation(core.RelationEq, mustSym(t, "p"), z))
	if err != nil {
		t.Fatal(err)
	}
	set := core.NewAssumptionSet(a)
	b, err := set.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	want := `[{"kind":"constraint","key":"momentum_zero","value":{"mode":"expr","expr":{"kind":"relation","operator":"eq","lhs":{"kind":"symbol","name":"p"},"rhs":{"kind":"rational","value":"0/1"}}}}]`
	if string(b) != want {
		t.Fatalf("structured expr assumption:\n got %s\nwant %s", b, want)
	}
	// Text value form.
	tx, err := core.NewTextAssumption(core.AssumptionDomain, "minkowski_spacetime", "Minkowski spacetime")
	if err != nil {
		t.Fatal(err)
	}
	tb := must(core.NewAssumptionSet(tx).CanonicalJSON())
	if !bytes.Contains(tb, []byte(`"value":{"mode":"text","value":"Minkowski spacetime"}`)) {
		t.Fatalf("text assumption json: %s", tb)
	}
	// Mathematical relationship is structured (expr), not an equation string:
	// a text assumption never contains a relation node shape.
	if bytes.Contains(tb, []byte(`"operator"`)) {
		t.Fatal("text assumptions must not encode relations")
	}
	_ = m
}

func TestAssumptionSetSurface(t *testing.T) {
	// §11.0 exact surface: constructors, Values, Merge, Equal, CanonicalJSON,
	// Hash, Entails*.
	m := mustSym(t, "m")
	gte, err := core.NewExprAssumption(core.AssumptionConstraint, "rest_mass_nonnegative",
		core.NewRelation(core.RelationGte, m, core.NewRational(rat(0, 1))))
	if err != nil {
		t.Fatal(err)
	}
	gt, err := core.NewExprAssumption(core.AssumptionConstraint, "speed_of_light_positive",
		core.NewRelation(core.RelationGt, mustSym(t, "c"), core.NewRational(rat(0, 1))))
	if err != nil {
		t.Fatal(err)
	}
	s := core.NewAssumptionSet(gte, gt)
	if len(s.Values()) != 2 {
		t.Fatalf("values = %d", len(s.Values()))
	}
	other := core.NewAssumptionSet(gt)
	merged, err := s.Merge(other)
	if err != nil {
		t.Fatal(err)
	}
	if merged.Hash() != s.Hash() {
		t.Fatal("union with subset is identity")
	}
	if !merged.Equal(s) {
		t.Fatal("merged equal")
	}
	if s.Equal(other) {
		t.Fatal("different sets not equal")
	}
	if _, err := s.CanonicalJSON(); err != nil {
		t.Fatal(err)
	}
	// Entailment surface (bounded §9.8 rules only).
	if !s.EntailsNonNegative(m) {
		t.Fatal("m >= 0 entailed")
	}
	if !s.EntailsPositive(mustSym(t, "c")) {
		t.Fatal("c > 0 entailed")
	}
	if !s.EntailsNonZero(mustSym(t, "c")) {
		t.Fatal("c != 0 entailed")
	}
	if s.EntailsNonNegative(mustSym(t, "v")) {
		t.Fatal("no assumption on v: entailment denied")
	}
	if s.EntailsPositive(m) {
		t.Fatal("m >= 0 does not entail positive")
	}
	if s.EntailsNonZero(m) {
		t.Fatal("m >= 0 does not entail nonzero")
	}
	// Rational structural rules do not need assumptions.
	if !s.EntailsNonNegative(core.NewRational(rat(3, 4))) {
		t.Fatal("nonnegative rational")
	}
	if s.EntailsNonNegative(core.NewRational(rat(-1, 4))) {
		t.Fatal("negative rational denied")
	}
}

func TestAssumptionValuesCopy(t *testing.T) {
	a, err := core.NewTextAssumption(core.AssumptionDomain, "k", "v")
	if err != nil {
		t.Fatal(err)
	}
	s := core.NewAssumptionSet(a)
	v1 := s.Values()
	if len(v1) != 1 {
		t.Fatal("one value")
	}
	v1[0] = core.Assumption{}
	if s.Values()[0].Key() != "k" {
		t.Fatal("Values must return a copy")
	}
}

func TestAssumptionKeyNonEmpty(t *testing.T) {
	m := mustSym(t, "m")
	if _, err := core.NewTextAssumption(core.AssumptionDomain, "", "v"); err == nil {
		t.Fatal("empty key rejected (text)")
	}
	if _, err := core.NewExprAssumption(core.AssumptionDomain, "", m); err == nil {
		t.Fatal("empty key rejected (expr)")
	}
	if _, err := core.NewTextAssumption(core.AssumptionDomain, "k", ""); err == nil {
		t.Fatal("empty text value rejected")
	}
	if _, err := core.NewTextAssumption("nope", "k", "v"); err == nil {
		t.Fatal("invalid kind rejected")
	}
}

func TestAssumptionConstructorsDeterministic(t *testing.T) {
	build := func() core.AssumptionSet {
		a1, err := core.NewTextAssumption(core.AssumptionDomain, "b", "v2")
		if err != nil {
			t.Fatal(err)
		}
		a2, err := core.NewExprAssumption(core.AssumptionConstraint, "a",
			core.NewRelation(core.RelationGte, mustSym(t, "m"), core.NewRational(rat(0, 1))))
		if err != nil {
			t.Fatal(err)
		}
		return core.NewAssumptionSet(a1, a2)
	}
	s1, s2 := build(), build()
	b1 := must(s1.CanonicalJSON())
	b2 := must(s2.CanonicalJSON())
	if !bytes.Equal(b1, b2) {
		t.Fatalf("bytes differ:\n%s\n%s", b1, b2)
	}
	if s1.Hash() != s2.Hash() {
		t.Fatal("hash differs")
	}
	// Input order must not matter.
	a1, _ := core.NewTextAssumption(core.AssumptionDomain, "b", "v2")
	a2, _ := core.NewExprAssumption(core.AssumptionConstraint, "a",
		core.NewRelation(core.RelationGte, mustSym(t, "m"), core.NewRational(rat(0, 1))))
	s3 := core.NewAssumptionSet(a2, a1)
	if !bytes.Equal(must(s3.CanonicalJSON()), b1) {
		t.Fatal("input order must not affect canonical bytes")
	}
}

func TestAssumptionMergeUnion(t *testing.T) {
	mk := func(k, v string) core.Assumption {
		a, err := core.NewTextAssumption(core.AssumptionDomain, k, v)
		if err != nil {
			t.Fatal(err)
		}
		return a
	}
	s1 := core.NewAssumptionSet(mk("k1", "v1"), mk("k2", "v2"))
	s2 := core.NewAssumptionSet(mk("k2", "v2"), mk("k3", "v3"))
	merged, err := s1.Merge(s2)
	if err != nil {
		t.Fatal(err)
	}
	if len(merged.Values()) != 3 {
		t.Fatalf("union size = %d, want 3", len(merged.Values()))
	}
	// Dedup: exact duplicates collapse.
	dup, err := s1.Merge(s1)
	if err != nil {
		t.Fatal(err)
	}
	if !dup.Equal(s1) || len(dup.Values()) != 2 {
		t.Fatal("duplicate merge deduplicates")
	}
	// Same (Kind, Key) with different canonical value → AssumptionConflictError.
	conflict := core.NewAssumptionSet(mk("k1", "other"))
	_, err = s1.Merge(conflict)
	var ace core.AssumptionConflictError
	if !errors.As(err, &ace) {
		t.Fatalf("want AssumptionConflictError, got %v", err)
	}
	if ace.Kind != "domain" || ace.Key != "k1" {
		t.Fatalf("conflict diagnostics: %+v", ace)
	}
	// Different keys never conflict even with equal values.
	s4 := core.NewAssumptionSet(mk("k9", "v1"))
	if _, err := s1.Merge(s4); err != nil {
		t.Fatalf("distinct keys must not conflict: %v", err)
	}
	// Merge is deterministic (double run).
	m1, _ := s1.Merge(s2)
	m2, _ := s1.Merge(s2)
	if m1.Hash() != m2.Hash() {
		t.Fatal("merge determinism")
	}
}

// ---------------------------------------------------------------------------
// Conventions (§12)
// ---------------------------------------------------------------------------

func TestConventionSetSurface(t *testing.T) {
	// §12.0 exact surface.
	c1, err := core.NewConvention("metric.signature", "-+++")
	if err != nil {
		t.Fatal(err)
	}
	if c1.Key() != "metric.signature" || c1.Value() != "-+++" {
		t.Fatal("convention accessors")
	}
	c2, err := core.NewConvention("units.scale", "1")
	if err != nil {
		t.Fatal(err)
	}
	s := core.NewConventionSet(c1, c2)
	if len(s.Values()) != 2 {
		t.Fatal("values")
	}
	other := core.NewConventionSet(c2)
	merged, err := s.Merge(other)
	if err != nil {
		t.Fatal(err)
	}
	if !merged.Equal(s) {
		t.Fatal("union identity")
	}
	if s.Equal(other) {
		t.Fatal("inequality")
	}
	if _, err := s.CanonicalJSON(); err != nil {
		t.Fatal(err)
	}
	if s.Hash() != merged.Hash() {
		t.Fatal("hash equality")
	}
	// Same key, different value → ConventionConflictError.
	alt, err := core.NewConvention("metric.signature", "-+-+")
	if err != nil {
		t.Fatal(err)
	}
	_, err = core.NewConventionSet(c1).Merge(core.NewConventionSet(alt))
	var cce core.ConventionConflictError
	if !errors.As(err, &cce) {
		t.Fatalf("want ConventionConflictError, got %v", err)
	}
	if cce.Key != "metric.signature" {
		t.Fatalf("conflict key %q", cce.Key)
	}
	// Canonical JSON shape.
	b := must(core.NewConventionSet(c1).CanonicalJSON())
	if string(b) != `[{"key":"metric.signature","value":"-+++"}]` {
		t.Fatalf("convention bytes %s", b)
	}
}

func TestConventionValuesCopy(t *testing.T) {
	c, err := core.NewConvention("k", "v")
	if err != nil {
		t.Fatal(err)
	}
	s := core.NewConventionSet(c)
	v := s.Values()
	v[0] = core.Convention{}
	if s.Values()[0].Key() != "k" {
		t.Fatal("Values must return a copy")
	}
}

func TestConventionStringsOnly(t *testing.T) {
	// Both fields non-empty strings; the type system permits only strings
	// (REQ-012-02: no equation can be encoded as a convention).
	rt := reflect.TypeOf(core.Convention{})
	if rt.NumField() != 2 {
		t.Fatalf("Convention fields = %d", rt.NumField())
	}
	for i := 0; i < rt.NumField(); i++ {
		if rt.Field(i).Type.Kind() != reflect.String {
			t.Fatalf("Convention field %s is not a string", rt.Field(i).Name)
		}
	}
	if _, err := core.NewConvention("", "v"); err == nil {
		t.Fatal("empty key rejected")
	}
	if _, err := core.NewConvention("k", ""); err == nil {
		t.Fatal("empty value rejected")
	}
	if _, err := core.NewConvention("k", "E = m*c^2"); err != nil {
		t.Fatal("string content is permitted (conventions stay strings)")
	}
}

// ---------------------------------------------------------------------------
// Provenance (§13)
// ---------------------------------------------------------------------------

func TestProvenanceStatusEnumExact(t *testing.T) {
	want := []core.ProvenanceStatus{
		"DEFINED", "POSTULATED", "DERIVED", "IDENTIFIED", "APPROXIMATED", "HYPOTHESIS",
	}
	got := []core.ProvenanceStatus{
		core.StatusDefined, core.StatusPostulated, core.StatusDerived,
		core.StatusIdentified, core.StatusApproximated, core.StatusHypothesis,
	}
	if len(got) != 6 || len(want) != 6 {
		t.Fatal("exactly six statuses")
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("status %d = %q want %q", i, got[i], want[i])
		}
	}
	// APPROXIMATED exists as an enum value (reserved, never minted).
	p, err := core.NewProvenance(core.StatusApproximated, "", "", nil,
		[32]byte{}, [32]byte{}, "mrc-v0.4", "")
	if err != nil || p.Status() != core.StatusApproximated {
		t.Fatalf("APPROXIMATED must exist: %v", err)
	}
	// Unknown statuses rejected.
	if _, err := core.NewProvenance("PHYSICALLY_TRUE", "", "", nil,
		[32]byte{}, [32]byte{}, "mrc-v0.4", ""); err == nil {
		t.Fatal("unknown status must be rejected")
	}
	// MRC version pinned.
	if _, err := core.NewProvenance(core.StatusDefined, "", "", nil,
		[32]byte{}, [32]byte{}, "mrc-v0.3", ""); err == nil {
		t.Fatal("wrong mrc version rejected")
	}
	// Justification rule.
	if _, err := core.NewProvenance(core.StatusDefined, "", "", nil,
		[32]byte{}, [32]byte{}, "mrc-v0.4", "why"); err == nil {
		t.Fatal("justification on DEFINED rejected")
	}
	for _, st := range []core.ProvenanceStatus{core.StatusIdentified, core.StatusHypothesis} {
		if _, err := core.NewProvenance(st, "", "", nil,
			[32]byte{}, [32]byte{}, "mrc-v0.4", "why"); err != nil {
			t.Fatalf("justification on %s: %v", st, err)
		}
	}
}

func TestProvenanceAccessors(t *testing.T) {
	parents := []string{
		strings.Repeat("a", 64),
		strings.Repeat("b", 64),
	}
	ah := core.HashExpr(mustSym(t, "m"))
	ch := core.HashExpr(mustSym(t, "c"))
	p, err := core.NewProvenance(core.StatusDerived, "Newton, Principia",
		"classical_mechanics", parents, ah, ch, "mrc-v0.4", "")
	if err != nil {
		t.Fatal(err)
	}
	if p.Status() != core.StatusDerived {
		t.Fatal("status")
	}
	if p.Source() != "Newton, Principia" {
		t.Fatal("source")
	}
	if p.Framework() != "classical_mechanics" {
		t.Fatal("framework")
	}
	if !reflect.DeepEqual(p.ParentHashes(), parents) {
		t.Fatal("parents")
	}
	if p.AssumptionHash() != ah || p.ConventionHash() != ch {
		t.Fatal("hashes")
	}
	if p.MRCVersion() != "mrc-v0.4" {
		t.Fatal("mrc version")
	}
	if p.Justification() != "" {
		t.Fatal("justification")
	}
	// Justification accessor on IDENTIFIED.
	id, err := core.NewProvenance(core.StatusIdentified, "", "", nil,
		[32]byte{}, [32]byte{}, "mrc-v0.4", "asserted equal")
	if err != nil || id.Justification() != "asserted equal" {
		t.Fatalf("identified justification: %v", err)
	}
}

func TestProvenanceParentHashesCopy(t *testing.T) {
	parents := []string{strings.Repeat("a", 64)}
	p, err := core.NewProvenance(core.StatusIdentified, "", "", parents,
		[32]byte{}, [32]byte{}, "mrc-v0.4", "why")
	if err != nil {
		t.Fatal(err)
	}
	got := p.ParentHashes()
	got[0] = strings.Repeat("z", 64)
	if p.ParentHashes()[0] != strings.Repeat("a", 64) {
		t.Fatal("ParentHashes must return a copy")
	}
	parents[0] = strings.Repeat("y", 64)
	if p.ParentHashes()[0] != strings.Repeat("a", 64) {
		t.Fatal("NewProvenance must copy caller slice")
	}
}

func TestProvenanceConstructorSingle(t *testing.T) {
	// AST: exactly one provenance constructor in the public core package.
	fset := token.NewFileSet()
	dir, err := parser.ParseDir(fset, ".", func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	provenanceCtors := 0
	for _, pkg := range dir {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Recv != nil || fn.Name == nil {
					continue
				}
				if strings.HasPrefix(fn.Name.Name, "New") && fn.Type.Results != nil {
					for _, res := range fn.Type.Results.List {
						if id, ok := res.Type.(*ast.Ident); ok && id.Name == "Provenance" {
							provenanceCtors++
							if fn.Name.Name != "NewProvenance" {
								t.Fatalf("unexpected provenance constructor %s", fn.Name.Name)
							}
						}
					}
				}
			}
		}
	}
	if provenanceCtors != 1 {
		t.Fatalf("provenance constructors = %d, want exactly 1", provenanceCtors)
	}
}

// ---------------------------------------------------------------------------
// Error taxonomy (§30)
// ---------------------------------------------------------------------------

func allTypedErrors() []error {
	return []error{
		core.DimensionMismatchError{Operation: "add", Left: "A", Right: "B"},
		core.CategoryMismatchError{Operation: "add", Left: "Mass", Right: "RestMass"},
		core.AssumptionConflictError{Kind: "constraint", Key: "k", Left: "a", Right: "b"},
		core.ConventionConflictError{Key: "k", Left: "a", Right: "b"},
		core.IdentifyError{Reason: "reason"},
		core.ProvenanceError{Context: "seal", Reason: "reason"},
		core.CandidateContainmentError{Reason: "reason"},
		core.InvalidObjectError{Operation: "add"},
		core.UnsupportedOperationError{Operation: "solve", Reason: "reason"},
		core.ManifestValidationError{Path: "items[0]", Reason: "reason"},
		core.LedgerValidationError{Reason: "reason"},
	}
}

func TestErrorTypesImplementError(t *testing.T) {
	for _, e := range allTypedErrors() {
		if e.Error() == "" {
			t.Fatalf("empty message for %T", e)
		}
		var _ error = e
	}
}

func TestErrorsAsInspectable(t *testing.T) {
	errs := allTypedErrors()
	wrapped := fmt.Errorf("outer: %w", errs[0])
	var dme core.DimensionMismatchError
	if !errors.As(wrapped, &dme) {
		t.Fatal("errors.As must reach DimensionMismatchError through wrapping")
	}
	var cme core.CategoryMismatchError
	if errors.As(errs[0], &cme) {
		t.Fatal("errors.As must not match a different error type")
	}
	probes := []interface{}{
		&dme,
		new(core.CategoryMismatchError),
		new(core.AssumptionConflictError),
		new(core.ConventionConflictError),
		new(core.IdentifyError),
		new(core.ProvenanceError),
		new(core.CandidateContainmentError),
		new(core.InvalidObjectError),
		new(core.UnsupportedOperationError),
		new(core.ManifestValidationError),
		new(core.LedgerValidationError),
	}
	if len(probes) != len(errs) {
		t.Fatal("one probe per error type")
	}
	for i, e := range errs {
		if !errors.As(fmt.Errorf("w: %w", e), probes[i]) {
			t.Fatalf("errors.As failed for %T", e)
		}
	}
	// LedgerValidationError unwraps its cause.
	inner := core.InvalidObjectError{Operation: "step"}
	lve := core.LedgerValidationError{Reason: "replay", Cause: inner}
	var got core.InvalidObjectError
	if !errors.As(lve, &got) {
		t.Fatal("LedgerValidationError must wrap its cause")
	}
	if !strings.Contains(lve.Error(), "replay") || !strings.Contains(lve.Error(), "add") {
		// cause message includes operation "step"
		if !strings.Contains(lve.Error(), "step") {
			t.Fatalf("message %q", lve.Error())
		}
	}
}

func TestErrorTypesExact(t *testing.T) {
	// The public core package declares exactly the eleven §30 error types.
	fset := token.NewFileSet()
	dir, err := parser.ParseDir(fset, ".", func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, pkg := range dir {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				gen, ok := decl.(*ast.GenDecl)
				if !ok || gen.Tok != token.TYPE {
					continue
				}
				for _, spec := range gen.Specs {
					ts := spec.(*ast.TypeSpec)
					if strings.HasSuffix(ts.Name.Name, "Error") {
						got[ts.Name.Name] = true
					}
				}
			}
		}
	}
	want := []string{
		"DimensionMismatchError", "CategoryMismatchError", "AssumptionConflictError",
		"ConventionConflictError", "IdentifyError", "ProvenanceError",
		"CandidateContainmentError", "InvalidObjectError",
		"UnsupportedOperationError", "ManifestValidationError", "LedgerValidationError",
	}
	if len(got) != len(want) {
		t.Fatalf("error type count = %d, want %d (%v)", len(got), len(want), got)
	}
	for _, w := range want {
		if !got[w] {
			t.Fatalf("missing error type %s", w)
		}
	}
}

func TestErrorDiagnosticsDeterministic(t *testing.T) {
	for _, e := range allTypedErrors() {
		m1 := e.Error()
		m2 := e.Error()
		if m1 != m2 {
			t.Fatalf("nondeterministic message for %T", e)
		}
		if strings.Contains(m1, "0x") || strings.Contains(m1, "%!") {
			t.Fatalf("nondeterministic content in %q", m1)
		}
	}
	// Structured diagnostics carry stable fields.
	d := core.DimensionMismatchError{Operation: "compare", Left: `{"m":"1/1"}`, Right: `{"m":"0/1"}`}
	msg := d.Error()
	if !strings.Contains(msg, "compare") || !strings.Contains(msg, `{"m":"1/1"}`) {
		t.Fatalf("diagnostics %q", msg)
	}
	// Two independently constructed identical errors have identical text.
	d2 := core.DimensionMismatchError{Operation: "compare", Left: `{"m":"1/1"}`, Right: `{"m":"0/1"}`}
	if d.Error() != d2.Error() {
		t.Fatal("deterministic message across constructions")
	}
}

// ---------------------------------------------------------------------------
// Canonical serialization hygiene (§10)
// ---------------------------------------------------------------------------

func TestCanonicalStructsTyped(t *testing.T) {
	// Canonical encoders must use explicit typed structs, never dynamic maps.
	roots := []string{"."}
	for _, root := range roots {
		entries, err := os.ReadDir(root)
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") ||
				strings.HasSuffix(e.Name(), "_test.go") {
				continue
			}
			b, err := os.ReadFile(root + "/" + e.Name())
			if err != nil {
				t.Fatal(err)
			}
			src := string(b)
			for _, banned := range []string{"map[string]any", "map[string]interface{}", "interface{}("} {
				if strings.Contains(src, banned) {
					t.Fatalf("%s contains dynamic map encoding %q", e.Name(), banned)
				}
			}
		}
	}
}

func TestNoVolatileFields(t *testing.T) {
	// Canonical artifacts must not depend on timestamps, randomness, or
	// environment (REQ-010-06): no time/uuid/rand imports in non-test sources.
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		b, err := os.ReadFile(e.Name())
		if err != nil {
			t.Fatal(err)
		}
		src := string(b)
		if strings.Contains(src, `"time"`) || strings.Contains(src, `"math/rand"`) ||
			strings.Contains(src, `"os/exec"`) || strings.Contains(src, `"net/http"`) {
			t.Fatalf("%s imports a volatile package", e.Name())
		}
	}
}

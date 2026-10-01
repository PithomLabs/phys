package relativity

import (
	"go/ast"
	"go/parser"
	"go/token"
	"math/big"
	"os"
	"runtime"
	"testing"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
	"github.com/PithomLabs/phys/ops"
)

// TestFullDerivationSequence tests the complete §20 derivation:
// EnergyMomentumRelation -> ZeroThreeMomentum -> Substitute -> Simplify -> Solve -> Compare -> SelectBranch
// Final result: m*c^2

// mustMint is a test-only convenience; production code uses MintObject directly.
func mustMint(t *testing.T, spec kernel.ObjectSpec) kernel.Object {
	t.Helper()
	obj, err := kernel.MintObject(spec)
	if err != nil {
		t.Fatalf("MintObject: %v", err)
	}
	return obj
}

func TestFullDerivationSequence(t *testing.T) {
	// Step 1: Get the energy-momentum relation
	emr := EnergyMomentumRelation() // E^2 = (p*c)^2 + (m*c^2)^2

	// Step 2: Zero three-momentum assumption
	zeroP := ZeroThreeMomentum()

	// Step 3: Substitute p = 0 into EMR
	pVar := mustMint(t, kernel.ObjectSpec{
		Name: "p", Kind: kernel.KindThreeMomentum,
		Dimension:   core.DimensionMomentum(),
		Expr:        kernel.NewSymbolMust("p"),
		Assumptions: kernel.NewAssumptionSet(),
		Conventions: kernel.NewConventionSet(),
		Provenance: kernel.NewProvenanceMust(kernel.StatusDefined, "", "", nil,
			[32]byte{}, [32]byte{}, kernel.MRCVersion, ""),
		CorpusStatus: kernel.CorpusNone,
	})
	subst, err := ops.Substitute(emr, pVar, zeroP)
	if err != nil {
		t.Fatalf("Substitute: %v", err)
	}
	if subst.Expr().Kind() != kernel.ExprRelation {
		t.Fatalf("result not a relation: %v", subst.Expr().Kind())
	}
	t.Logf("After Substitute assumptions: %v", subst.Assumptions().Values())

	// Step 4: Simplify the relation
	simp, err := ops.Simplify(subst)
	if err != nil {
		t.Fatalf("Simplify: %v", err)
	}
	if simp.Expr().Kind() != kernel.ExprRelation {
		t.Fatalf("simplified not a relation: %v", simp.Expr().Kind())
	}
	t.Logf("After Simplify assumptions: %v", simp.Assumptions().Values())

	// Step 5: Solve for E (E = sqrt(m^2*c^4))
	eTarget := mustMint(t, kernel.ObjectSpec{
		Name: "E", Kind: kernel.KindEnergy,
		Dimension:   core.DimensionEnergy(),
		Expr:        kernel.NewSymbolMust("E"),
		Assumptions: kernel.NewAssumptionSet(),
		Conventions: kernel.NewConventionSet(),
		Provenance: kernel.NewProvenanceMust(kernel.StatusDefined, "", "", nil,
			[32]byte{}, [32]byte{}, kernel.MRCVersion, ""),
		CorpusStatus: kernel.CorpusNone,
	})
	solved, err := ops.Solve(simp, eTarget)
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	if solved.Kind() != kernel.KindBranchSet {
		t.Fatalf("solved kind = %v, want BranchSet", solved.Kind())
	}
	t.Logf("After Solve assumptions: %v", solved.Assumptions().Values())
	t.Logf("After Solve branches: %v", solved.Expr())

	// Step 6: Compare Energy() with ZeroEnergy() using GTE
	energy := NewEnergy()
	zeroEnergy := ZeroEnergy()
	t.Logf("Energy assumptions: %v", energy.CoreObject().Assumptions().Values())
	t.Logf("ZeroEnergy assumptions: %v", zeroEnergy.Assumptions().Values())
	cmp, err := ops.Compare(energy.CoreObject(), zeroEnergy, kernel.RelationGte)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	if cmp.Expr().RelationOperator() != kernel.RelationGte {
		t.Fatalf("compare operator != gte")
	}
	t.Logf("Compare assumptions: %v", cmp.Assumptions().Values())

	// Step 7: SelectBranch with the constraint E >= 0
	selected, err := ops.SelectBranch(solved, cmp)
	if err != nil {
		t.Fatalf("SelectBranch: %v", err)
	}
	t.Logf("After SelectBranch: %v", selected.Expr())
	t.Logf("Selected assumptions: %v", selected.Assumptions().Values())

	// Final result should be m*c^2 (positive branch)
	if selected.Kind() != kernel.KindExpression {
		t.Fatalf("selected kind = %v, want Expression", selected.Kind())
	}

	e := selected.Expr()
	if e.Kind() != kernel.ExprMul {
		t.Fatalf("selected expr not mul: %v", e.Kind())
	}
	factors := e.Children()
	if len(factors) != 2 {
		t.Fatalf("expected 2 factors, got %d", len(factors))
	}
	foundM := false
	foundCSq := false
	for _, f := range factors {
		if f.Kind() == kernel.ExprSymbol && f.SymbolName() == "m" {
			foundM = true
		}
		if f.Kind() == kernel.ExprPow {
			base := f.Base()
			if base.Kind() == kernel.ExprSymbol && base.SymbolName() == "c" {
				if f.Exponent().Cmp(big.NewRat(2, 1)) == 0 {
					foundCSq = true
				}
			}
		}
	}
	if !foundM || !foundCSq {
		t.Errorf("selected expression not m*c^2: m=%v c^2=%v", foundM, foundCSq)
	}
}
func TestLorentzFactorLimit(t *testing.T) {
	// J: Limit(LorentzFactor(), Velocity(), ZeroVelocity()) expands the fixed
	// body, substitutes v=0, and simplifies to exactly 1/1.
	out, err := ops.Limit(LorentzFactor(), NewVelocity().CoreObject(), ZeroVelocity())
	if err != nil {
		t.Fatalf("Limit: %v", err)
	}
	if out.Kind() != kernel.KindExpression {
		t.Fatalf("kind = %v, want Expression", out.Kind())
	}
	raw, err := kernel.CanonicalExprJSON(out.Expr())
	if err != nil {
		t.Fatalf("canonical: %v", err)
	}
	if string(raw) != `{"kind":"rational","value":"1/1"}` {
		t.Fatalf("limit result = %s, want 1/1", raw)
	}
	// A hardcoded 1 minted without body traversal would lack the merged
	// speed-of-light assumption carried by the Lorentz factor.
	if !hasKey(out, "speed_of_light_positive") {
		t.Errorf("limit result missing merged assumption speed_of_light_positive (have %v)", assumptionKeys(out))
	}
}

// assumptionKeys returns the sorted key list of an object's assumptions.
func assumptionKeys(o kernel.Object) []string {
	keys := []string{}
	for _, a := range o.Assumptions().Values() {
		keys = append(keys, a.Key())
	}
	return keys
}

func hasKey(o kernel.Object, key string) bool {
	for _, k := range assumptionKeys(o) {
		if k == key {
			return true
		}
	}
	return false
}

// TestEnergyMomentumRelation asserts the canonical E^2=(p*c)^2+(m*c^2)^2
// shape, the Energy^2 dimension, the two required sign assumptions, and the
// manifest cross-check. Acceptance H.
func TestEnergyMomentumRelation(t *testing.T) {
	rel := EnergyMomentumRelation()
	if rel.Kind() != kernel.KindRelation {
		t.Fatalf("kind = %v, want Relation", rel.Kind())
	}
	if !rel.Dimension().Equal(core.DimensionEnergy().Multiply(core.DimensionEnergy())) {
		t.Fatalf("dimension mismatch: want Energy^2")
	}
	for _, key := range []string{"rest_mass_nonnegative", "speed_of_light_positive"} {
		if !hasKey(rel, key) {
			t.Errorf("missing assumption %q (have %v)", key, assumptionKeys(rel))
		}
	}
	e := rel.Expr()
	if e.Kind() != kernel.ExprRelation || e.RelationOperator() != kernel.RelationEq {
		t.Fatalf("not an eq relation")
	}
	if e.Left().Kind() != kernel.ExprPow {
		t.Fatalf("lhs not Pow(E,2)")
	}
	m, err := core.ParseManifest(ManifestJSON)
	if err != nil {
		t.Fatalf("ParseManifest: %v", err)
	}
	found := false
	for _, item := range m.Items {
		if item.ID == "EnergyMomentumRelation" {
			found = true
			if !core.EqualExpr(item.CanonicalExpr, rel.Expr()) {
				t.Errorf("manifest canonical_expr mismatch")
			}
		}
	}
	if !found {
		t.Fatalf("EnergyMomentumRelation missing from manifest")
	}
}

// TestMassEnergyDerivation executes the exact §20 sequence with real
// operations and asserts each golden-trace state. Acceptance I.
// (Deliberately never calls MassEnergyRelation: the result is derived.)
func TestMassEnergyDerivation(t *testing.T) {
	two := big.NewRat(2, 1)
	sym := func(name string) kernel.Expr {
		e, err := kernel.NewSymbol(name)
		if err != nil {
			t.Fatalf("symbol %s: %v", name, err)
		}
		return e
	}
	mustObj := func(name string, kind kernel.Kind, dim core.Dimension, expr kernel.Expr) kernel.Object {
		prov, err := kernel.NewProvenance(kernel.StatusDefined, "", "", nil,
			[32]byte{}, [32]byte{}, kernel.MRCVersion, "")
		if err != nil {
			t.Fatalf("provenance: %v", err)
		}
		obj, err := kernel.MintObject(kernel.ObjectSpec{
			Name: name, Kind: kind, Dimension: dim, Expr: expr,
			Assumptions: kernel.NewAssumptionSet(), Conventions: kernel.NewConventionSet(),
			Provenance: prov, CorpusStatus: kernel.CorpusNone,
		})
		if err != nil {
			t.Fatalf("mint %s: %v", name, err)
		}
		return obj
	}

	emr := EnergyMomentumRelation()
	zeroP := ZeroThreeMomentum()
	// ZeroThreeMomentum must carry the rest-frame constraint (F-005).
	if !hasKey(zeroP, "rest_frame") {
		t.Fatalf("ZeroThreeMomentum lacks RestFrameAssumption (have %v)", assumptionKeys(zeroP))
	}
	if zeroP.Kind() != kernel.KindThreeMomentum || !zeroP.Dimension().Equal(core.DimensionMomentum()) {
		t.Fatalf("ZeroThreeMomentum kind/dimension mismatch")
	}

	pVar := mustObj("p", kernel.KindThreeMomentum, core.DimensionMomentum(), sym("p"))
	subst, err := ops.Substitute(emr, pVar, zeroP)
	if err != nil {
		t.Fatalf("Substitute: %v", err)
	}
	// Golden state 2: E^2 = (0*c)^2 + (m*c^2)^2 — rhs mentions Rational 0.
	substRaw, _ := kernel.CanonicalExprJSON(subst.Expr())
	if !contains(string(substRaw), `"value":"0/1"`) {
		t.Fatalf("golden state 2 missing zero momentum: %s", substRaw)
	}

	simp, err := ops.Simplify(subst)
	if err != nil {
		t.Fatalf("Simplify: %v", err)
	}
	// Golden state 3: E^2 = (m*c^2)^2.
	mc2 := kernel.NewMul(sym("m"), kernel.NewPow(sym("c"), two))
	wantRHS := kernel.NewPow(mc2, two)
	wantRel := kernel.NewRelation(kernel.RelationEq, kernel.NewPow(sym("E"), two), wantRHS)
	if !core.EqualExpr(simp.Expr(), wantRel) {
		got, _ := kernel.CanonicalExprJSON(simp.Expr())
		want, _ := kernel.CanonicalExprJSON(wantRel)
		t.Fatalf("golden state 3 mismatch:\n got %s\nwant %s", got, want)
	}

	eTarget := mustObj("E", kernel.KindEnergy, core.DimensionEnergy(), sym("E"))
	solved, err := ops.Solve(simp, eTarget)
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	// Golden state 4: BranchSet(E, [Sqrt(...), Neg(Sqrt(...))]).
	if solved.Kind() != kernel.KindBranchSet {
		t.Fatalf("golden state 4 kind = %v, want BranchSet", solved.Kind())
	}
	if len(solved.Expr().Branches()) != 2 {
		t.Fatalf("golden state 4 branch count = %d, want 2", len(solved.Expr().Branches()))
	}

	cmp, err := ops.Compare(NewEnergy().CoreObject(), ZeroEnergy(), kernel.RelationGte)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	if cmp.Expr().RelationOperator() != kernel.RelationGte {
		t.Fatalf("constraint operator != gte")
	}

	selected, err := ops.SelectBranch(solved, cmp)
	if err != nil {
		t.Fatalf("SelectBranch: %v", err)
	}
	// Final golden state: m*c^2.
	wantFinal := kernel.NewMul(sym("m"), kernel.NewPow(sym("c"), two))
	if !core.EqualExpr(selected.Expr(), wantFinal) {
		got, _ := kernel.CanonicalExprJSON(selected.Expr())
		want, _ := kernel.CanonicalExprJSON(wantFinal)
		t.Fatalf("final mismatch:\n got %s\nwant %s", got, want)
	}
	if selected.Provenance().Status() != kernel.StatusDerived {
		t.Fatalf("final status = %v, want DERIVED", selected.Provenance().Status())
	}
	// A hardcoded m*c^2 minted without the derivation would lack the merged
	// sign assumptions; the real replayed result carries both.
	for _, key := range []string{"rest_mass_nonnegative", "speed_of_light_positive"} {
		if !hasKey(selected, key) {
			t.Errorf("final result missing merged assumption %q (have %v)", key, assumptionKeys(selected))
		}
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (func() bool {
		for i := 0; i+len(needle) <= len(haystack); i++ {
			if haystack[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	}())
}

// deriveMassEnergyForOwnership replays the §20 operation chain without golden
// assertions, returning the final selected object for ownership analysis.
func deriveMassEnergyForOwnership(t *testing.T) kernel.Object {
	t.Helper()
	sym := func(name string) kernel.Expr {
		e, err := kernel.NewSymbol(name)
		if err != nil {
			t.Fatalf("symbol %s: %v", name, err)
		}
		return e
	}
	mustObj := func(name string, kind kernel.Kind, dim core.Dimension, expr kernel.Expr) kernel.Object {
		prov, err := kernel.NewProvenance(kernel.StatusDefined, "", "", nil,
			[32]byte{}, [32]byte{}, kernel.MRCVersion, "")
		if err != nil {
			t.Fatalf("provenance: %v", err)
		}
		obj, err := kernel.MintObject(kernel.ObjectSpec{
			Name: name, Kind: kind, Dimension: dim, Expr: expr,
			Assumptions: kernel.NewAssumptionSet(), Conventions: kernel.NewConventionSet(),
			Provenance: prov, CorpusStatus: kernel.CorpusNone,
		})
		if err != nil {
			t.Fatalf("mint %s: %v", name, err)
		}
		return obj
	}
	emr := EnergyMomentumRelation()
	zeroP := ZeroThreeMomentum()
	pVar := mustObj("p", kernel.KindThreeMomentum, core.DimensionMomentum(), sym("p"))
	subst, err := ops.Substitute(emr, pVar, zeroP)
	if err != nil {
		t.Fatalf("Substitute: %v", err)
	}
	simp, err := ops.Simplify(subst)
	if err != nil {
		t.Fatalf("Simplify: %v", err)
	}
	eTarget := mustObj("E", kernel.KindEnergy, core.DimensionEnergy(), sym("E"))
	solved, err := ops.Solve(simp, eTarget)
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	cmp, err := ops.Compare(NewEnergy().CoreObject(), ZeroEnergy(), kernel.RelationGte)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	selected, err := ops.SelectBranch(solved, cmp)
	if err != nil {
		t.Fatalf("SelectBranch: %v", err)
	}
	return selected
}

// TestFinalAssumptionOwnership (Gate H5): every assumption on the derived
// m*c² traces to an explicit input object or operation step — none arrives
// ambiently. Trace: rest_mass_nonnegative + speed_of_light_positive ←
// EnergyMomentumRelation inputs; rest_frame ← ZeroThreeMomentum via
// Substitute; selected_branch/* ← SelectBranch operation contract.
func TestFinalAssumptionOwnership(t *testing.T) {
	if testing.Short() {
		t.Skip("ownership requires full derivation")
	}
	selected := deriveMassEnergyForOwnership(t)
	keys := map[string]bool{}
	for _, a := range selected.Assumptions().Values() {
		keys[a.Key()] = true
	}
	for _, want := range []string{"rest_mass_nonnegative", "speed_of_light_positive", "rest_frame"} {
		if !keys[want] {
			t.Errorf("final result lacks traced assumption %q (have %v)", want, keys)
		}
	}
	branch := 0
	for k := range keys {
		if len(k) >= 16 && k[:16] == "selected_branch/" {
			branch++
		}
	}
	if branch != 1 {
		t.Errorf("want exactly 1 selected_branch/* assumption, have %v", keys)
	}
	if len(keys) != 4 {
		t.Errorf("final assumption keys = %v, want exactly the 4 traced keys", keys)
	}
}

// TestDerivationUsesNoStoredResult (E=mc² firewall): derivation_test.go must
// contain zero non-comment references to MassEnergyRelation — the result must
// come from the operation pipeline, never corpus retrieval.
func TestDerivationUsesNoStoredResult(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller unavailable")
	}
	raw, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	pf, err := parser.ParseFile(fset, file, raw, 0)
	if err != nil {
		t.Fatal(err)
	}
	ast.Inspect(pf, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && id.Name == "MassEnergyRelation" {
			t.Errorf("forbidden reference to MassEnergyRelation at %v", fset.Position(id.Pos()))
		}
		return true
	})
}

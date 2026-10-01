package mechanics

import (
	"math/big"
	"testing"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
	"github.com/PithomLabs/phys/ops"
)

// TestNewtonSecondLawManifestMatch tests the surface of the three relation constructors.
// Acceptance A: mechanics relations have correct expressions and metadata.
func TestNewtonSecondLawManifestMatch(t *testing.T) {
	tests := []struct {
		name           string
		fn             func() kernel.Object
		wantKind       kernel.Kind
		wantDim        kernel.Dimension
		wantProvStatus kernel.ProvenanceStatus
		wantSrc        string
	}{
		{"NewtonSecondLaw", NewtonSecondLaw, kernel.KindRelation, core.DimensionForce(), kernel.StatusDefined, "Newton, Principia"},
		{"MomentumRelation", MomentumRelation, kernel.KindRelation, core.DimensionMomentum(), kernel.StatusDefined, "Classical Mechanics corpus"},
		{"KineticEnergyRelation", KineticEnergyRelation, kernel.KindRelation, core.DimensionEnergy(), kernel.StatusDefined, "Classical Mechanics corpus"},
	}
	for _, tc := range tests {
		obj := tc.fn()
		if obj.Kind() != tc.wantKind {
			t.Errorf("%s: kind = %v, want %v", tc.name, obj.Kind(), tc.wantKind)
		}
		if !obj.Dimension().Equal(tc.wantDim) {
			t.Errorf("%s: dimension = %v, want %v", tc.name, obj.Dimension(), tc.wantDim)
		}
		if obj.Provenance().Status() != tc.wantProvStatus {
			t.Errorf("%s: provenance status = %v, want %v", tc.name, obj.Provenance().Status(), tc.wantProvStatus)
		}
		if obj.Provenance().Source() != tc.wantSrc {
			t.Errorf("%s: source = %q, want %q", tc.name, obj.Provenance().Source(), tc.wantSrc)
		}
		if obj.CorpusStatus() != kernel.CorpusEstablished {
			t.Errorf("%s: corpus status = %v, want ESTABLISHED", tc.name, obj.CorpusStatus())
		}
		if obj.Expr().Kind() != kernel.ExprRelation {
			t.Errorf("%s: expr not a relation", tc.name)
		}
		if obj.Expr().RelationOperator() != kernel.RelationEq {
			t.Errorf("%s: operator != eq", tc.name)
		}
	}
}

// TestDifferentiateKineticEnergy tests the parameterized constructor.
// Acceptance B: NewKineticEnergy(m,v) produces correct expression.
func TestDifferentiateKineticEnergy(t *testing.T) {
	m := NewMass()
	v := NewVelocity()
	k := NewKineticEnergy(m, v)
	obj := k.CoreObject()
	if obj.Kind() != kernel.KindKineticEnergy {
		t.Errorf("kind = %v, want KineticEnergy", obj.Kind())
	}
	if !obj.Dimension().Equal(core.DimensionEnergy()) {
		t.Errorf("dimension = %v, want Energy", obj.Dimension())
	}
	e := obj.Expr()
	if e.Kind() != kernel.ExprMul {
		t.Fatalf("expr not mul: %v", e.Kind())
	}
	factors := e.Children()
	if len(factors) != 3 {
		t.Fatalf("expected 3 factors, got %d", len(factors))
	}
	foundM := false
	foundHalf := false
	foundVSq := false
	for _, f := range factors {
		switch f.Kind() {
		case kernel.ExprSymbol:
			if f.SymbolName() == "m" {
				foundM = true
			}
		case kernel.ExprRational:
			if f.RationalValue().Cmp(big.NewRat(1, 2)) == 0 {
				foundHalf = true
			}
		case kernel.ExprPow:
			base := f.Base()
			if base.Kind() == kernel.ExprSymbol && base.SymbolName() == "v" {
				if f.Exponent().Cmp(big.NewRat(2, 1)) == 0 {
					foundVSq = true
				}
			}
		}
	}
	if !foundM || !foundHalf || !foundVSq {
		t.Errorf("expression factors mismatch: m=%v 1/2=%v v^2=%v", foundM, foundHalf, foundVSq)
	}
	if obj.Provenance().Status() != kernel.StatusDerived {
		t.Errorf("provenance status = %v, want DERIVED", obj.Provenance().Status())
	}
	if obj.CorpusStatus() != kernel.CorpusNone {
		t.Errorf("corpus status = %v, want NONE", obj.CorpusStatus())
	}
}

// TestTypedMechanicsConstructors constructs all eight fixed mechanics wrappers
// via fixed constructors only and asserts the exact accessor, kind,
// dimension, symbol, DEFINED provenance, and ESTABLISHED corpus status.
// Acceptance A.
func TestTypedMechanicsConstructors(t *testing.T) {
	type wrapper struct {
		name string
		obj  kernel.Object
		kind kernel.Kind
		dim  kernel.Dimension
		sym  string
	}
	cases := []wrapper{
		{"Mass", NewMass().CoreObject(), kernel.KindMass, core.DimensionMass(), "m"},
		{"Time", NewTime().CoreObject(), kernel.KindTime, core.DimensionTime(), "t"},
		{"Position", NewPosition().CoreObject(), kernel.KindPosition, core.DimensionLength(), "x"},
		{"Velocity", NewVelocity().CoreObject(), kernel.KindVelocity, core.DimensionVelocity(), "v"},
		{"Acceleration", NewAcceleration().CoreObject(), kernel.KindAcceleration, core.DimensionAcceleration(), "a"},
		{"Force", NewForce().CoreObject(), kernel.KindForce, core.DimensionForce(), "F"},
		{"Momentum", NewMomentum().CoreObject(), kernel.KindMomentum, core.DimensionMomentum(), "p"},
		{"Energy", NewEnergy().CoreObject(), kernel.KindEnergy, core.DimensionEnergy(), "E"},
	}
	for _, tc := range cases {
		if !tc.obj.Valid() {
			t.Errorf("%s: invalid object", tc.name)
			continue
		}
		if tc.obj.Kind() != tc.kind {
			t.Errorf("%s: kind = %v, want %v", tc.name, tc.obj.Kind(), tc.kind)
		}
		if !tc.obj.Dimension().Equal(tc.dim) {
			t.Errorf("%s: dimension mismatch", tc.name)
		}
		if tc.obj.Expr().Kind() != kernel.ExprSymbol || tc.obj.Expr().SymbolName() != tc.sym {
			t.Errorf("%s: symbol = %v, want %s", tc.name, tc.obj.Expr(), tc.sym)
		}
		if tc.obj.Provenance().Status() != kernel.StatusDefined {
			t.Errorf("%s: provenance = %v, want DEFINED", tc.name, tc.obj.Provenance().Status())
		}
		if tc.obj.CorpusStatus() != kernel.CorpusEstablished {
			t.Errorf("%s: corpus status = %v, want ESTABLISHED", tc.name, tc.obj.CorpusStatus())
		}
		if tc.obj.Name() == "" {
			t.Errorf("%s: empty name", tc.name)
		}
	}
}

// TestDifferentiateKineticEnergyOps differentiates 1/2 m v^2 with respect to
// v through ops.Differentiate and asserts the canonical m*v result, the
// Expression kind, the Momentum dimension, and Compare acceptance.
// Acceptance G (operational half; constructor half is TestDifferentiateKineticEnergy).
func TestDifferentiateKineticEnergyOps(t *testing.T) {
	m := NewMass()
	v := NewVelocity()
	k := NewKineticEnergy(m, v)
	deriv, err := ops.Differentiate(k.CoreObject(), v.CoreObject())
	if err != nil {
		t.Fatalf("Differentiate: %v", err)
	}
	if deriv.Kind() != kernel.KindExpression {
		t.Errorf("kind = %v, want Expression", deriv.Kind())
	}
	if !deriv.Dimension().Equal(core.DimensionMomentum()) {
		t.Errorf("dimension mismatch: want Momentum")
	}
	want := core.NewMul(
		mustSymbol(t, "m"),
		mustSymbol(t, "v"),
	)
	if !core.EqualExpr(deriv.Expr(), want) {
		t.Errorf("derivative mismatch:\n got %s\nwant %s", jsonOf(t, deriv.Expr()), jsonOf(t, want))
	}
	// Expression-vs-named Compare must be accepted (MRC-003 §6.2).
	cmp, err := ops.Compare(deriv, NewMomentum().CoreObject(), kernel.RelationEq)
	if err != nil {
		t.Errorf("Compare(deriv, Momentum, eq): %v", err)
	} else if cmp.Expr().RelationOperator() != kernel.RelationEq {
		t.Errorf("compare operator != eq")
	}
}

func mustSymbol(t *testing.T, name string) kernel.Expr {
	t.Helper()
	e, err := kernel.NewSymbol(name)
	if err != nil {
		t.Fatalf("symbol %s: %v", name, err)
	}
	return e
}

func jsonOf(t *testing.T, e kernel.Expr) string {
	t.Helper()
	raw, err := kernel.CanonicalExprJSON(e)
	if err != nil {
		t.Fatalf("canonical expr: %v", err)
	}
	return string(raw)
}

// TestNoRelativityLeakage (Gate H4, mechanics side): ops derivations from
// mechanics objects must carry none of the relativity-only framework keys.
// Forbidden-intersection form: legitimate accumulation is untouched.
func TestNoRelativityLeakage(t *testing.T) {
	relOnly := []string{
		"rest_frame", "rest_mass_nonnegative", "speed_of_light_positive",
		"minkowski_spacetime", "lorentz_symmetry", "no_gravitational_dynamics",
		"special_relativistic_regime",
	}
	m := NewMass()
	v := NewVelocity()
	k := NewKineticEnergy(m, v)
	deriv, err := ops.Differentiate(k.CoreObject(), v.CoreObject())
	if err != nil {
		t.Fatalf("differentiate: %v", err)
	}
	sum, err := ops.Add(m.CoreObject(), m.CoreObject())
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	for _, o := range []kernel.Object{deriv, sum} {
		for _, a := range o.Assumptions().Values() {
			for _, rk := range relOnly {
				if a.Key() == rk {
					t.Errorf("mechanics result acquired relativity key %q", rk)
				}
			}
		}
	}
}

package mechanics

import (
	"math/big"

	_ "embed"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

//go:embed manifest.json
var ManifestJSON []byte

// Mass is the canonical mass object (symbol m, dimension M).
type Mass struct {
	Object kernel.Object
}

func NewMass() Mass {
	sym, _ := kernel.NewSymbol("m")
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Classical Mechanics corpus", "classical_mechanics",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "m",
		Kind:         kernel.KindMass,
		Dimension:    core.DimensionMass(),
		Expr:         sym,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return Mass{obj}
}

// Time is the canonical time object (symbol t, dimension T).
type Time struct {
	Object kernel.Object
}

func NewTime() Time {
	sym, _ := kernel.NewSymbol("t")
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Classical Mechanics corpus", "classical_mechanics",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "t",
		Kind:         kernel.KindTime,
		Dimension:    core.DimensionTime(),
		Expr:         sym,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return Time{obj}
}

// Position is the canonical position object (symbol x, dimension L).
type Position struct {
	Object kernel.Object
}

func NewPosition() Position {
	sym, _ := kernel.NewSymbol("x")
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Classical Mechanics corpus", "classical_mechanics",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "x",
		Kind:         kernel.KindPosition,
		Dimension:    core.DimensionLength(),
		Expr:         sym,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return Position{obj}
}

// Velocity is the canonical velocity object (symbol v, dimension L/T).
type Velocity struct {
	Object kernel.Object
}

func NewVelocity() Velocity {
	sym, _ := kernel.NewSymbol("v")
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Classical Mechanics corpus", "classical_mechanics",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "v",
		Kind:         kernel.KindVelocity,
		Dimension:    core.DimensionVelocity(),
		Expr:         sym,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return Velocity{obj}
}

// Acceleration is the canonical acceleration object (symbol a, dimension L/T^2).
type Acceleration struct {
	Object kernel.Object
}

func NewAcceleration() Acceleration {
	sym, _ := kernel.NewSymbol("a")
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Classical Mechanics corpus", "classical_mechanics",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "a",
		Kind:         kernel.KindAcceleration,
		Dimension:    core.DimensionAcceleration(),
		Expr:         sym,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return Acceleration{obj}
}

// Force is the canonical force object (symbol F, dimension M*L/T^2).
type Force struct {
	Object kernel.Object
}

func NewForce() Force {
	sym, _ := kernel.NewSymbol("F")
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Classical Mechanics corpus", "classical_mechanics",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "F",
		Kind:         kernel.KindForce,
		Dimension:    core.DimensionForce(),
		Expr:         sym,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return Force{obj}
}

// Momentum is the canonical momentum object (symbol p, dimension M*L/T).
type Momentum struct {
	Object kernel.Object
}

func NewMomentum() Momentum {
	sym, _ := kernel.NewSymbol("p")
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Classical Mechanics corpus", "classical_mechanics",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "p",
		Kind:         kernel.KindMomentum,
		Dimension:    core.DimensionMomentum(),
		Expr:         sym,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return Momentum{obj}
}

// Energy is the canonical energy object (symbol E, dimension M*L^2/T^2).
type Energy struct {
	Object kernel.Object
}

func NewEnergy() Energy {
	sym, _ := kernel.NewSymbol("E")
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Classical Mechanics corpus", "classical_mechanics",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "E",
		Kind:         kernel.KindEnergy,
		Dimension:    core.DimensionEnergy(),
		Expr:         sym,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return Energy{obj}
}

// KineticEnergy is the canonical kinetic energy object (symbol K, dimension Energy).
type KineticEnergy struct {
	Object kernel.Object
}

// NewKineticEnergy constructs the kinetic energy relation K = 1/2 * m * v^2.
func NewKineticEnergy(m Mass, v Velocity) KineticEnergy {
	half := kernel.NewRational(big.NewRat(1, 2))
	mExpr := m.Object.Expr()
	vExpr := v.Object.Expr()
	vSq := kernel.NewPow(vExpr, big.NewRat(2, 1))
	expr := kernel.NewMul(half, mExpr, vSq)
	prov, _ := kernel.NewProvenance(kernel.StatusDerived, "", "", nil,
		[32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "K",
		Kind:         kernel.KindKineticEnergy,
		Dimension:    core.DimensionEnergy(),
		Expr:         expr,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusNone,
	})
	return KineticEnergy{obj}
}

// CoreObject returns the underlying core.Object.
func (m Mass) CoreObject() core.Object       { return m.Object }
func (t Time) CoreObject() core.Object       { return t.Object }
func (p Position) CoreObject() core.Object   { return p.Object }
func (v Velocity) CoreObject() core.Object   { return v.Object }
func (a Acceleration) CoreObject() core.Object { return a.Object }
func (f Force) CoreObject() core.Object      { return f.Object }
func (m Momentum) CoreObject() core.Object   { return m.Object }
func (e Energy) CoreObject() core.Object     { return e.Object }
func (k KineticEnergy) CoreObject() core.Object { return k.Object }

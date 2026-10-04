package mechanics

import (
	"math/big"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

// NewtonSecondLaw returns the relation F = m*a.
func NewtonSecondLaw() kernel.Object {
	m, _ := kernel.NewSymbol("m")
	a, _ := kernel.NewSymbol("a")
	lhs, _ := kernel.NewSymbol("F")
	rhs := kernel.NewMul(m, a)
	expr := kernel.NewRelation(kernel.RelationEq, lhs, rhs)
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Newton, Principia", "classical_mechanics",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "NewtonSecondLaw",
		Kind:         kernel.KindRelation,
		Dimension:    core.DimensionForce(),
		Expr:         expr,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return obj
}

// MomentumRelation returns the relation p = m*v.
func MomentumRelation() kernel.Object {
	m, _ := kernel.NewSymbol("m")
	v, _ := kernel.NewSymbol("v")
	lhs, _ := kernel.NewSymbol("p")
	rhs := kernel.NewMul(m, v)
	expr := kernel.NewRelation(kernel.RelationEq, lhs, rhs)
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Classical Mechanics corpus", "classical_mechanics",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "MomentumRelation",
		Kind:         kernel.KindRelation,
		Dimension:    core.DimensionMomentum(),
		Expr:         expr,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return obj
}

// KineticEnergyRelation returns the relation K = 1/2*m*v^2.
func KineticEnergyRelation() kernel.Object {
	half := kernel.NewRational(big.NewRat(1, 2))
	m, _ := kernel.NewSymbol("m")
	v, _ := kernel.NewSymbol("v")
	vSq := kernel.NewPow(v, big.NewRat(2, 1))
	lhs, _ := kernel.NewSymbol("K")
	rhs := kernel.NewMul(half, m, vSq)
	expr := kernel.NewRelation(kernel.RelationEq, lhs, rhs)
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Classical Mechanics corpus", "classical_mechanics",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "KineticEnergyRelation",
		Kind:         kernel.KindRelation,
		Dimension:    core.DimensionEnergy(),
		Expr:         expr,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return obj
}

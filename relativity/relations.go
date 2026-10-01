package relativity

import (
	"math/big"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

// LorentzFactor returns the Lorentz factor γ = 1/sqrt(1 - v^2/c^2) as a Call.
func LorentzFactor() kernel.Object {
	v, _ := kernel.NewSymbol("v")
	call, _ := kernel.NewCall(kernel.LorentzFactorFunctionID, v)
	cPosExpr := kernel.NewRelation(kernel.RelationGt, kernel.NewSymbolMust("c"), kernel.NewRational(big.NewRat(0, 1)))
	cPos := kernel.NewExprAssumptionMust(kernel.AssumptionConstraint, "speed_of_light_positive", cPosExpr)
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Special Relativity corpus", "special_relativity",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "LorentzFactor",
		Kind:         kernel.KindExpression,
		Dimension:    core.Dimensionless(),
		Expr:         call,
		Assumptions:  kernel.NewAssumptionSet(cPos),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return obj
}

// EnergyMomentumRelation returns E^2 = (p*c)^2 + (m*c^2)^2.
func EnergyMomentumRelation() kernel.Object {
	p, _ := kernel.NewSymbol("p")
	c, _ := kernel.NewSymbol("c")
	m, _ := kernel.NewSymbol("m")
	E_sym, _ := kernel.NewSymbol("E")
	two := big.NewRat(2, 1)
	pc := kernel.NewMul(p, c)
	pcSq := kernel.NewPow(pc, two)
	mcSq := kernel.NewMul(m, kernel.NewPow(c, two))
	mcSqSq := kernel.NewPow(mcSq, two)
	lhs := kernel.NewPow(E_sym, two)
	rhs := kernel.NewAdd(pcSq, mcSqSq)
	expr := kernel.NewRelation(kernel.RelationEq, lhs, rhs)
	mNonNegExpr := kernel.NewRelation(kernel.RelationGte, kernel.NewSymbolMust("m"), kernel.NewRational(big.NewRat(0, 1)))
	mNonNeg := kernel.NewExprAssumptionMust(kernel.AssumptionConstraint, "rest_mass_nonnegative", mNonNegExpr)
	cPosExpr := kernel.NewRelation(kernel.RelationGt, kernel.NewSymbolMust("c"), kernel.NewRational(big.NewRat(0, 1)))
	cPos := kernel.NewExprAssumptionMust(kernel.AssumptionConstraint, "speed_of_light_positive", cPosExpr)
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Einstein, 1905", "special_relativity",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "EnergyMomentumRelation",
		Kind:         kernel.KindRelation,
		Dimension:    core.DimensionEnergy().Multiply(core.DimensionEnergy()),
		Expr:         expr,
		Assumptions:  kernel.NewAssumptionSet(mNonNeg, cPos),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return obj
}

// MassEnergyRelation returns E = m*c^2.
func MassEnergyRelation() kernel.Object {
	m, _ := kernel.NewSymbol("m")
	c, _ := kernel.NewSymbol("c")
	lhs, _ := kernel.NewSymbol("E")
	rhs := kernel.NewMul(m, kernel.NewPow(c, big.NewRat(2, 1)))
	expr := kernel.NewRelation(kernel.RelationEq, lhs, rhs)
	mNonNegExpr := kernel.NewRelation(kernel.RelationGte, kernel.NewSymbolMust("m"), kernel.NewRational(big.NewRat(0, 1)))
	mNonNeg := kernel.NewExprAssumptionMust(kernel.AssumptionConstraint, "rest_mass_nonnegative", mNonNegExpr)
	cPosExpr := kernel.NewRelation(kernel.RelationGt, kernel.NewSymbolMust("c"), kernel.NewRational(big.NewRat(0, 1)))
	cPos := kernel.NewExprAssumptionMust(kernel.AssumptionConstraint, "speed_of_light_positive", cPosExpr)
	prov, _ := kernel.NewProvenance(kernel.StatusDerived, "Einstein, 1905", "special_relativity",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "MassEnergyRelation",
		Kind:         kernel.KindRelation,
		Dimension:    core.DimensionEnergy(),
		Expr:         expr,
		Assumptions:  kernel.NewAssumptionSet(mNonNeg, cPos),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return obj
}

// ZeroThreeMomentum returns the zero three-momentum object.
func ZeroThreeMomentum() kernel.Object {
	zero := kernel.NewRational(big.NewRat(0, 1))
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Special Relativity corpus", "special_relativity",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "ZeroThreeMomentum",
		Kind:         kernel.KindThreeMomentum,
		Dimension:    core.DimensionMomentum(),
		Expr:         zero,
		Assumptions:  kernel.NewAssumptionSet(RestFrameAssumption()),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return obj
}

// ZeroEnergy returns the zero energy object.
func ZeroEnergy() kernel.Object {
	zero := kernel.NewRational(big.NewRat(0, 1))
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Special Relativity corpus", "special_relativity",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "ZeroEnergy",
		Kind:         kernel.KindEnergy,
		Dimension:    core.DimensionEnergy(),
		Expr:         zero,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return obj
}

// ZeroVelocity returns the zero velocity object.
func ZeroVelocity() kernel.Object {
	zero := kernel.NewRational(big.NewRat(0, 1))
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Special Relativity corpus", "special_relativity",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "ZeroVelocity",
		Kind:         kernel.KindVelocity,
		Dimension:    core.DimensionVelocity(),
		Expr:         zero,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return obj
}

// RestFrameAssumption returns the rest frame assumption: p = 0.
func RestFrameAssumption() kernel.Assumption {
	p, _ := kernel.NewSymbol("p")
	zero := kernel.NewRational(big.NewRat(0, 1))
	expr := kernel.NewRelation(kernel.RelationEq, p, zero)
	a, _ := kernel.NewExprAssumption(kernel.AssumptionConstraint, "rest_frame", expr)
	return a
}

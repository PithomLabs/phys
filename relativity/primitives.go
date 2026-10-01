package relativity

import (
	"math/big"

	_ "embed"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

//go:embed manifest.json
var ManifestJSON []byte

// Spacetime is the canonical spacetime object (symbol s).
type Spacetime struct {
	Object kernel.Object
}

func NewSpacetime() Spacetime {
	sym, _ := kernel.NewSymbol("s")
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Einstein, 1905", "special_relativity",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "s",
		Kind:         kernel.KindSpacetime,
		Dimension:    core.DimensionLength(),
		Expr:         sym,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return Spacetime{obj}
}

// MinkowskiMetric is the canonical Minkowski metric object (symbol η).
type MinkowskiMetric struct {
	Object kernel.Object
}

func NewMinkowskiMetric() MinkowskiMetric {
	sym, _ := kernel.NewSymbol("eta")
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Einstein, 1905", "special_relativity",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "eta",
		Kind:         kernel.KindMinkowski,
		Dimension:    core.Dimensionless(),
		Expr:         sym,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(kernel.NewConventionMust("metric.signature", "-+++")),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return MinkowskiMetric{obj}
}

// RestMass is the canonical rest mass object (symbol m).
type RestMass struct {
	Object kernel.Object
}

func NewRestMass() RestMass {
	sym, _ := kernel.NewSymbol("m")
	mNonNegExpr := kernel.NewRelation(kernel.RelationGte, kernel.NewSymbolMust("m"), kernel.NewRational(big.NewRat(0, 1)))
	mNonNeg := kernel.NewExprAssumptionMust(kernel.AssumptionConstraint, "rest_mass_nonnegative", mNonNegExpr)
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Einstein, 1905", "special_relativity",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "m",
		Kind:         kernel.KindRestMass,
		Dimension:    core.DimensionMass(),
		Expr:         sym,
		Assumptions:  kernel.NewAssumptionSet(mNonNeg),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return RestMass{obj}
}

// Energy is the canonical energy object (symbol E).
type Energy struct {
	Object kernel.Object
}

func NewEnergy() Energy {
	sym, _ := kernel.NewSymbol("E")
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Einstein, 1905", "special_relativity",
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

// ThreeMomentum is the canonical three-momentum object (symbol p).
type ThreeMomentum struct {
	Object kernel.Object
}

func NewThreeMomentum() ThreeMomentum {
	sym, _ := kernel.NewSymbol("p")
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Einstein, 1905", "special_relativity",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "p",
		Kind:         kernel.KindThreeMomentum,
		Dimension:    core.DimensionMomentum(),
		Expr:         sym,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return ThreeMomentum{obj}
}

// FourMomentum is the canonical four-momentum object (symbol P).
type FourMomentum struct {
	Object kernel.Object
}

func NewFourMomentum() FourMomentum {
	sym, _ := kernel.NewSymbol("P")
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Einstein, 1905", "special_relativity",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "P",
		Kind:         kernel.KindFourMomentum,
		Dimension:    core.DimensionMomentum(),
		Expr:         sym,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return FourMomentum{obj}
}

// SpeedOfLight is the canonical speed of light object (symbol c).
type SpeedOfLight struct {
	Object kernel.Object
}

func NewSpeedOfLight() SpeedOfLight {
	sym, _ := kernel.NewSymbol("c")
	cPosExpr := kernel.NewRelation(kernel.RelationGt, kernel.NewSymbolMust("c"), kernel.NewRational(big.NewRat(0, 1)))
	cPos := kernel.NewExprAssumptionMust(kernel.AssumptionConstraint, "speed_of_light_positive", cPosExpr)
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Einstein, 1905", "special_relativity",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	obj, _ := kernel.MintObject(kernel.ObjectSpec{
		Name:         "c",
		Kind:         kernel.KindSpeedOfLight,
		Dimension:    core.DimensionVelocity(),
		Expr:         sym,
		Assumptions:  kernel.NewAssumptionSet(cPos),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	})
	return SpeedOfLight{obj}
}

// Velocity is the canonical velocity object (symbol v).
type Velocity struct {
	Object kernel.Object
}

func NewVelocity() Velocity {
	sym, _ := kernel.NewSymbol("v")
	prov, _ := kernel.NewProvenance(kernel.StatusDefined, "Einstein, 1905", "special_relativity",
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

// CoreObject returns the underlying kernel.Object.
func (s Spacetime) CoreObject() kernel.Object       { return s.Object }
func (m MinkowskiMetric) CoreObject() kernel.Object { return m.Object }
func (m RestMass) CoreObject() kernel.Object        { return m.Object }
func (e Energy) CoreObject() kernel.Object          { return e.Object }
func (p ThreeMomentum) CoreObject() kernel.Object   { return p.Object }
func (p FourMomentum) CoreObject() kernel.Object    { return p.Object }
func (c SpeedOfLight) CoreObject() kernel.Object    { return c.Object }
func (v Velocity) CoreObject() kernel.Object        { return v.Object }

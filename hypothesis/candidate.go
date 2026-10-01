// Hypothesis package — candidate concept construction (§24 / MRC-008 / Plan 10 v2.3).
// Only the session authority (via Session.Seal) produces trusted
// ResearchCandidate artifacts. This package provides the unverified
// concept constructor (forced HYPOTHESIS / NONE, no promotion APIs).
package hypothesis

import (
	"fmt"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

// ---------------------------------------------------------------------------
// CandidateConcept — the unverified hypothesis concept (§24.2)
// ---------------------------------------------------------------------------

type CandidateConcept struct {
	id         string
	dimension  core.Dimension
	expr       core.Expr
	assumptions core.AssumptionSet
}

// NewCandidateConcept creates a forced HYPOTHESIS concept (REQ-024-01 / §24.1).
// Status forced to HYPOTHESIS; Source forced to "hypothesis"; Framework empty;
// CorpusStatus forced to NONE. No caller-supplied provenance allowed.
func NewCandidateConcept(id string, kind core.Kind, dimension core.Dimension, expr core.Expr, assumptions ...core.Assumption) (core.Object, error) {
	if id == "" {
		return core.Object{}, fmt.Errorf("candidate concept: non-empty id required")
	}
	prov, err := kernel.NewProvenance(kernel.StatusHypothesis, "hypothesis", "", nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	if err != nil {
		return core.Object{}, fmt.Errorf("candidate provenance: %w", err)
	}
	obj, err := kernel.MintObject(kernel.ObjectSpec{
		Name:         id,
		Kind:         kind,
		Dimension:    dimension,
		Expr:         expr,
		Assumptions:  core.NewAssumptionSet(assumptions...),
		Conventions:  core.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusNone,
	})
	if err != nil {
		return core.Object{}, fmt.Errorf("candidate mint: %w", err)
	}
	return obj, nil
}

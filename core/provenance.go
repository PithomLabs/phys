package core

import (
	"github.com/PithomLabs/phys/internal/kernel"
)

// Aliases over the kernel provenance metadata (specs_v2_3.md §13).
type (
	Provenance       = kernel.Provenance
	ProvenanceStatus = kernel.ProvenanceStatus
)

// The exact provenance statuses of §13.1; APPROXIMATED exists but no MVP
// operation mints it (REQ-013-01).
const (
	StatusDefined      = kernel.StatusDefined
	StatusPostulated   = kernel.StatusPostulated
	StatusDerived      = kernel.StatusDerived
	StatusIdentified   = kernel.StatusIdentified
	StatusApproximated = kernel.StatusApproximated
	StatusHypothesis   = kernel.StatusHypothesis
)

// NewProvenance is the single deterministic public provenance constructor
// (§13.0.1); it validates status-specific invariants but never mints an
// Object. It is a façade over the kernel primitive (plan F3).
func NewProvenance(
	status ProvenanceStatus,
	source string,
	framework string,
	parentHashes []string,
	assumptionHash [32]byte,
	conventionHash [32]byte,
	mrcVersion string,
	justification string,
) (Provenance, error) {
	return kernel.NewProvenance(
		status, source, framework, parentHashes,
		assumptionHash, conventionHash, mrcVersion, justification,
	)
}

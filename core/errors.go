package core

import (
	"github.com/PithomLabs/phys/internal/kernel"
)

// The exact mandatory error taxonomy of specs_v2_3.md §30. All eleven typed
// errors are declared once (kernel) and aliased here so that errors.As works
// through the public API. All messages are deterministic.
type (
	// DimensionMismatchError reports an MRC-002 failure.
	DimensionMismatchError = kernel.DimensionMismatchError
	// CategoryMismatchError reports an MRC-003 failure.
	CategoryMismatchError = kernel.CategoryMismatchError
	// AssumptionConflictError reports an MRC-004 failure.
	AssumptionConflictError = kernel.AssumptionConflictError
	// ConventionConflictError reports an MRC-005 failure.
	ConventionConflictError = kernel.ConventionConflictError
	// IdentifyError reports an MRC-006 failure.
	IdentifyError = kernel.IdentifyError
	// ProvenanceError reports a provenance/MRC-007 failure.
	ProvenanceError = kernel.ProvenanceError
	// CandidateContainmentError reports an MRC-008 failure.
	CandidateContainmentError = kernel.CandidateContainmentError
	// InvalidObjectError reports rejection of an invalid zero object.
	InvalidObjectError = kernel.InvalidObjectError
	// UnsupportedOperationError reports a bounded-operation rejection.
	UnsupportedOperationError = kernel.UnsupportedOperationError
	// ManifestValidationError reports a manifest validation failure.
	ManifestValidationError = kernel.ManifestValidationError
	// LedgerValidationError reports a ledger structural/replay failure.
	LedgerValidationError = kernel.LedgerValidationError
)

package core

import (
	"github.com/PithomLabs/phys/internal/kernel"
)

// Aliases over the kernel convention metadata (specs_v2_3.md §12).
type (
	Convention    = kernel.Convention
	ConventionSet = kernel.ConventionSet
)

// NewConvention builds a convention; both fields must be non-empty.
func NewConvention(key, value string) (Convention, error) {
	return kernel.NewConvention(key, value)
}

// NewConventionSet canonicalizes, sorts, and deduplicates conventions.
func NewConventionSet(values ...Convention) ConventionSet {
	return kernel.NewConventionSet(values...)
}

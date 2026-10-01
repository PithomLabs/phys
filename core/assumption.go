package core

import (
	"github.com/PithomLabs/phys/internal/kernel"
)

// Aliases over the kernel assumption metadata (specs_v2_3.md §11).
type (
	Assumption     = kernel.Assumption
	AssumptionKind = kernel.AssumptionKind
	AssumptionSet  = kernel.AssumptionSet
)

// The exact seven assumption kinds of §11.1 (Approximation reserved).
const (
	AssumptionDomain             = kernel.AssumptionDomain
	AssumptionRegime             = kernel.AssumptionRegime
	AssumptionConstraint         = kernel.AssumptionConstraint
	AssumptionConvention         = kernel.AssumptionConvention
	AssumptionApproximation      = kernel.AssumptionApproximation
	AssumptionMathPrecondition   = kernel.AssumptionMathPrecondition
	AssumptionPhysicalAssumption = kernel.AssumptionPhysicalAssumption
)

// NewTextAssumption builds a text-valued assumption (façade over the kernel
// primitive; validation lives in the kernel).
func NewTextAssumption(kind AssumptionKind, key, value string) (Assumption, error) {
	return kernel.NewTextAssumption(kind, key, value)
}

// NewExprAssumption builds an expression-valued assumption.
func NewExprAssumption(kind AssumptionKind, key string, value Expr) (Assumption, error) {
	return kernel.NewExprAssumption(kind, key, value)
}

// NewAssumptionSet canonicalizes, sorts, and deduplicates assumptions.
func NewAssumptionSet(values ...Assumption) AssumptionSet {
	return kernel.NewAssumptionSet(values...)
}

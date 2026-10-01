package core

import (
	"math/big"

	"github.com/PithomLabs/phys/internal/kernel"
)

// Dimension is the public alias over the kernel seven-base-dimension value
// (specs_v2_3.md §7).
type Dimension = kernel.Dimension

// dim builds a dimension from integer exponents M, L, T (the other four
// bases are zero in every fixed MVP constructor). Construction cannot fail
// for literal integer exponents.
func dim(m, l, t int64) Dimension {
	d, _ := kernel.NewDimension(
		big.NewRat(m, 1), big.NewRat(l, 1), big.NewRat(t, 1),
		big.NewRat(0, 1), big.NewRat(0, 1), big.NewRat(0, 1), big.NewRat(0, 1),
	)
	return d
}

// The exact fixed dimension constructors of §7.0. Each is a façade over the
// kernel-side dimension constructor (plan F3).

// Dimensionless returns the all-zero dimension.
func Dimensionless() Dimension { return dim(0, 0, 0) }

// DimensionMass returns M.
func DimensionMass() Dimension { return dim(1, 0, 0) }

// DimensionLength returns L.
func DimensionLength() Dimension { return dim(0, 1, 0) }

// DimensionTime returns T.
func DimensionTime() Dimension { return dim(0, 0, 1) }

// DimensionVelocity returns L T^-1.
func DimensionVelocity() Dimension { return dim(0, 1, -1) }

// DimensionAcceleration returns L T^-2.
func DimensionAcceleration() Dimension { return dim(0, 1, -2) }

// DimensionForce returns M L T^-2.
func DimensionForce() Dimension { return dim(1, 1, -2) }

// DimensionMomentum returns M L T^-1.
func DimensionMomentum() Dimension { return dim(1, 1, -1) }

// DimensionEnergy returns M L^2 T^-2.
func DimensionEnergy() Dimension { return dim(1, 2, -2) }

// DimensionEnergySquared returns M^2 L^4 T^-4.
func DimensionEnergySquared() Dimension { return dim(2, 4, -4) }

// NewDimension builds a dimension from exact rational exponents for the
// seven SI base dimensions (M, L, T, I, Θ, N, J).
func NewDimension(m, l, t, i, theta, n, j *big.Rat) (Dimension, error) {
	return kernel.NewDimension(m, l, t, i, theta, n, j)
}

// ParseDimensionJSON decodes a byte-canonical dimension document.
func ParseDimensionJSON(data []byte) (Dimension, error) {
	return kernel.ParseDimensionJSON(data)
}

// Package core is the public façade of the Physics Compiler MVP. It exposes
// aliases over the internal kernel representation plus package-level
// façade constructors and helpers; it declares no methods on aliased types
// and performs no validation of its own (validation lives once, in
// internal/kernel).
package core

import (
	"math/big"

	"github.com/PithomLabs/phys/internal/kernel"
)

// Aliases over the kernel expression representation (specs_v2_3.md §8).
type (
	Expr             = kernel.Expr
	ExprKind         = kernel.ExprKind
	RelationOperator = kernel.RelationOperator
)

// Closed expression node ordinals (§8.2.1).
const (
	ExprSymbol    = kernel.ExprSymbol
	ExprRational  = kernel.ExprRational
	ExprAdd       = kernel.ExprAdd
	ExprMul       = kernel.ExprMul
	ExprNeg       = kernel.ExprNeg
	ExprPow       = kernel.ExprPow
	ExprSqrt      = kernel.ExprSqrt
	ExprCall      = kernel.ExprCall
	ExprRelation  = kernel.ExprRelation
	ExprBranchSet = kernel.ExprBranchSet
)

// The closed relation operator set (§8.2.1); canonical JSON strings.
const (
	RelationEq  = kernel.RelationEq
	RelationNeq = kernel.RelationNeq
	RelationLt  = kernel.RelationLt
	RelationLte = kernel.RelationLte
	RelationGt  = kernel.RelationGt
	RelationGte = kernel.RelationGte
)

// The exact public expression constructors/helpers of §8.0, each a façade
// over the kernel-side structural canonicalization primitive (plan F3: one
// authoritative validation location).

// NewSymbol builds a Symbol node.
func NewSymbol(name string) (Expr, error) { return kernel.NewSymbol(name) }

// NewRational builds a Rational node, copying the caller's value.
func NewRational(value *big.Rat) Expr { return kernel.NewRational(value) }

// NewAdd builds an Add node (flatten, exact rational combination, sort).
func NewAdd(terms ...Expr) Expr { return kernel.NewAdd(terms...) }

// NewMul builds a Mul node (flatten, exact rational coefficients, sign
// normal form, sort).
func NewMul(factors ...Expr) Expr { return kernel.NewMul(factors...) }

// NewNeg builds a Neg node with sign normal form.
func NewNeg(expr Expr) Expr { return kernel.NewNeg(expr) }

// NewPow builds a Pow node, cloning the exact rational exponent.
func NewPow(base Expr, exponent *big.Rat) Expr { return kernel.NewPow(base, exponent) }

// NewSqrt builds a Sqrt node.
func NewSqrt(expr Expr) Expr { return kernel.NewSqrt(expr) }

// NewCall builds the sole MVP Call form (lorentz_factor).
func NewCall(functionID string, args ...Expr) (Expr, error) {
	return kernel.NewCall(functionID, args...)
}

// NewRelation builds a Relation node preserving operator and side order.
func NewRelation(op RelationOperator, lhs, rhs Expr) Expr {
	return kernel.NewRelation(op, lhs, rhs)
}

// NewBranchSet builds a bounded BranchSet (target first, then branches).
func NewBranchSet(target Expr, branches ...Expr) Expr {
	return kernel.NewBranchSet(target, branches...)
}

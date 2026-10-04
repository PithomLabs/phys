package core

import (
	"github.com/PithomLabs/phys/internal/kernel"
)

// EqualExpr compares canonical symbolic structure, never display strings
// (§9.1, REQ-009-01).
func EqualExpr(a, b Expr) bool { return kernel.EqualExpr(a, b) }

// HashExpr returns the SHA-256 of the expression's canonical JSON bytes
// (§9.2).
func HashExpr(e Expr) [32]byte { return kernel.HashExpr(e) }

// HashAssumptionSet returns the SHA-256 of the canonical assumption-set
// bytes (§9.2).
func HashAssumptionSet(s AssumptionSet) [32]byte { return s.Hash() }

// HashConventionSet returns the SHA-256 of the canonical convention-set
// bytes (§9.2).
func HashConventionSet(s ConventionSet) [32]byte { return s.Hash() }

// CanonicalExprJSON encodes an expression in its exact §10.2 canonical form.
func CanonicalExprJSON(e Expr) ([]byte, error) { return kernel.CanonicalExprJSON(e) }

// ParseExprJSON decodes canonical expression JSON into the closed node set,
// enforcing invariants and byte-level canonicality.
func ParseExprJSON(data []byte) (Expr, error) { return kernel.ParseExprJSON(data) }

// EqualObject compares canonical object representation, including all
// authoritative physical metadata — never display strings (§9.1,
// REQ-009-01/02).
func EqualObject(a, b Object) bool { return kernel.EqualObject(a, b) }

// HashObject returns the SHA-256 of the object's canonical JSON bytes
// (§9.2).
func HashObject(o Object) [32]byte { return kernel.HashObject(o) }

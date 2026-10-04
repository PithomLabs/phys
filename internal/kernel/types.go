// Package kernel is the internal immutable representation authority for the
// Physics Compiler MVP (mrc-v0.4). It owns the closed Expr node set, Kind,
// Dimension, metadata value types, structural canonicalization, canonical JSON,
// SHA-256 hashing, defensive-copy accessors, and the kernel-side construction
// primitives required by every public metadata constructor.
//
// This package sits below the module's public import boundary: outside-module
// callers cannot import internal/kernel (Go language rule). Nothing in this
// package performs I/O, uses floating point, or stores callbacks.
package kernel

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"
)

// ---------------------------------------------------------------------------
// Fixed MVP constants
// ---------------------------------------------------------------------------

// MRCVersion is the single provenance/ledger MRC version string for MVP
// artifacts.
const MRCVersion = "mrc-v0.4"

// SchemaVersion is the Plan 10 serialization pin for every canonical artifact
// schema (object, ledger, session, candidate, manifest).
const SchemaVersion = "1"

// LorentzFactorFunctionID is the only function_id permitted in MVP Call nodes.
const LorentzFactorFunctionID = "lorentz_factor"

// GenesisStepHash is the PreviousStepHash of the first committed step (64
// lowercase hex zeros).
const GenesisStepHash = "0000000000000000000000000000000000000000000000000000000000000000"

// ---------------------------------------------------------------------------
// Mandatory error taxonomy (specs_v2_3.md §30)
// ---------------------------------------------------------------------------
//
// All eleven typed errors are declared here (kernel emits AssumptionConflict,
// ConventionConflict, InvalidObject, and Provenance errors during
// construction/mint) and are re-exported by core/errors.go as type aliases so
// errors.As works through the public API. All messages are deterministic.

// DimensionMismatchError reports an MRC-002 dimensional compatibility failure.
type DimensionMismatchError struct {
	Operation string
	Left      string
	Right     string
}

func (e DimensionMismatchError) Error() string {
	return "dimension mismatch in " + e.Operation + ": " + e.Left + " != " + e.Right
}

// CategoryMismatchError reports an MRC-003 physical-kind compatibility failure.
type CategoryMismatchError struct {
	Operation string
	Left      string
	Right     string
}

func (e CategoryMismatchError) Error() string {
	return "category mismatch in " + e.Operation + ": " + e.Left + " vs " + e.Right
}

// AssumptionConflictError reports an MRC-004 assumption merge conflict.
type AssumptionConflictError struct {
	Kind  string
	Key   string
	Left  string
	Right string
}

func (e AssumptionConflictError) Error() string {
	return "assumption conflict on (" + e.Kind + ", " + e.Key + ")"
}

// ConventionConflictError reports an MRC-005 convention merge conflict.
type ConventionConflictError struct {
	Key   string
	Left  string
	Right string
}

func (e ConventionConflictError) Error() string {
	return "convention conflict on key " + e.Key
}

// IdentifyError reports an MRC-006 identification failure.
type IdentifyError struct {
	Reason string
}

func (e IdentifyError) Error() string {
	return "identify failed: " + e.Reason
}

// ProvenanceError reports an MRC-007 / provenance-contract failure.
type ProvenanceError struct {
	Context string
	Reason  string
}

func (e ProvenanceError) Error() string {
	if e.Context == "" {
		return "provenance error: " + e.Reason
	}
	return "provenance error in " + e.Context + ": " + e.Reason
}

// CandidateContainmentError reports an MRC-008 candidate containment failure.
type CandidateContainmentError struct {
	Reason string
}

func (e CandidateContainmentError) Error() string {
	return "candidate containment error: " + e.Reason
}

// InvalidObjectError reports rejection of an invalid (zero) core.Object.
type InvalidObjectError struct {
	Operation string
}

func (e InvalidObjectError) Error() string {
	if e.Operation == "" {
		return "invalid object"
	}
	return "invalid object in " + e.Operation
}

// UnsupportedOperationError reports a bounded-operation rejection.
type UnsupportedOperationError struct {
	Operation string
	Reason    string
}

func (e UnsupportedOperationError) Error() string {
	if e.Reason == "" {
		return "unsupported operation: " + e.Operation
	}
	return "unsupported operation " + e.Operation + ": " + e.Reason
}

// ManifestValidationError reports a corpus manifest validation failure.
type ManifestValidationError struct {
	Path   string
	Reason string
}

func (e ManifestValidationError) Error() string {
	if e.Path == "" {
		return "manifest validation error: " + e.Reason
	}
	return "manifest validation error at " + e.Path + ": " + e.Reason
}

// LedgerValidationError reports a derivation-ledger structural or replay
// failure. Cause carries the specific wrapped typed error when one exists.
type LedgerValidationError struct {
	Reason string
	Cause  error
}

func (e LedgerValidationError) Error() string {
	if e.Cause != nil {
		return "ledger validation error: " + e.Reason + ": " + e.Cause.Error()
	}
	return "ledger validation error: " + e.Reason
}

// Unwrap exposes the wrapped typed error for errors.As/Is.
func (e LedgerValidationError) Unwrap() error { return e.Cause }

// ---------------------------------------------------------------------------
// ExprKind — closed expression node set (specs_v2_3.md §8.2/§8.2.1)
// ---------------------------------------------------------------------------

// ExprKind is the closed MVP expression node-kind ordinal enum. Ordinals are
// part of canonicalization and stable for mrc-v0.4.
type ExprKind uint8

// The exact node ordinals 0..9.
const (
	ExprSymbol    ExprKind = 0 // Symbol
	ExprRational  ExprKind = 1 // Rational
	ExprAdd       ExprKind = 2 // Add
	ExprMul       ExprKind = 3 // Mul
	ExprNeg       ExprKind = 4 // Neg
	ExprPow       ExprKind = 5 // Pow
	ExprSqrt      ExprKind = 6 // Sqrt
	ExprCall      ExprKind = 7 // Call
	ExprRelation  ExprKind = 8 // Relation
	ExprBranchSet ExprKind = 9 // BranchSet
)

// String returns the canonical JSON node-kind name (specs_v2_3.md §10.2).
// Unknown ordinals yield the empty string; no alias names exist.
func (k ExprKind) String() string {
	switch k {
	case ExprSymbol:
		return "symbol"
	case ExprRational:
		return "rational"
	case ExprAdd:
		return "add"
	case ExprMul:
		return "mul"
	case ExprNeg:
		return "neg"
	case ExprPow:
		return "pow"
	case ExprSqrt:
		return "sqrt"
	case ExprCall:
		return "call"
	case ExprRelation:
		return "relation"
	case ExprBranchSet:
		return "branch_set"
	default:
		return ""
	}
}

func (k ExprKind) valid() bool { return k <= ExprBranchSet }

// ---------------------------------------------------------------------------
// RelationOperator (specs_v2_3.md §8.2.1, §9.9)
// ---------------------------------------------------------------------------

// RelationOperator is the closed relation operator set with canonical JSON
// strings; no alias operator names are accepted.
type RelationOperator string

// The exact operator values.
const (
	RelationEq  RelationOperator = "eq"
	RelationNeq RelationOperator = "neq"
	RelationLt  RelationOperator = "lt"
	RelationLte RelationOperator = "lte"
	RelationGt  RelationOperator = "gt"
	RelationGte RelationOperator = "gte"
)

func validRelationOperator(op RelationOperator) bool {
	switch op {
	case RelationEq, RelationNeq, RelationLt, RelationLte, RelationGt, RelationGte:
		return true
	default:
		return false
	}
}

// ---------------------------------------------------------------------------
// Kind — physical kind (specs_v2_3.md §6)
// ---------------------------------------------------------------------------

// Kind is the physical kind of an object. Ordinals are stable for mrc-v0.4.
// Kind and Dimension are intentionally independent (MRC-003).
type Kind uint8

// The exact 18 MVP kinds in listed order (§6).
const (
	KindMass          Kind = 0
	KindTime          Kind = 1
	KindPosition      Kind = 2
	KindVelocity      Kind = 3
	KindAcceleration  Kind = 4
	KindForce         Kind = 5
	KindMomentum      Kind = 6
	KindEnergy        Kind = 7
	KindKineticEnergy Kind = 8
	KindRestMass      Kind = 9
	KindThreeMomentum Kind = 10
	KindFourMomentum  Kind = 11
	KindSpeedOfLight  Kind = 12
	KindSpacetime     Kind = 13
	KindMinkowski     Kind = 14
	KindExpression    Kind = 15
	KindRelation      Kind = 16
	KindBranchSet     Kind = 17
)

// String returns the manifest/JSON kind name (§6 listing).
func (k Kind) String() string {
	switch k {
	case KindMass:
		return "Mass"
	case KindTime:
		return "Time"
	case KindPosition:
		return "Position"
	case KindVelocity:
		return "Velocity"
	case KindAcceleration:
		return "Acceleration"
	case KindForce:
		return "Force"
	case KindMomentum:
		return "Momentum"
	case KindEnergy:
		return "Energy"
	case KindKineticEnergy:
		return "KineticEnergy"
	case KindRestMass:
		return "RestMass"
	case KindThreeMomentum:
		return "ThreeMomentum"
	case KindFourMomentum:
		return "FourMomentum"
	case KindSpeedOfLight:
		return "SpeedOfLight"
	case KindSpacetime:
		return "Spacetime"
	case KindMinkowski:
		return "MinkowskiMetric"
	case KindExpression:
		return "Expression"
	case KindRelation:
		return "Relation"
	case KindBranchSet:
		return "BranchSet"
	default:
		return ""
	}
}

func (k Kind) valid() bool { return k <= KindBranchSet }

// NamedKind reports whether k is a fixed named physical quantity kind (as
// opposed to the operation-derived kinds Expression, Relation, BranchSet).
func (k Kind) NamedKind() bool {
	return k.valid() && k <= KindSpeedOfLight
}

// orderedKind reports whether k may appear on either side of an inequality
// under the bounded MVP ordering rule (§6.2).
func orderedKind(k Kind) bool {
	switch k {
	case KindMass, KindRestMass, KindTime, KindEnergy, KindKineticEnergy,
		KindSpeedOfLight, KindExpression:
		return true
	default:
		return false
	}
}

// ---------------------------------------------------------------------------
// AssumptionKind — exact §11.1 seven-kind set
// ---------------------------------------------------------------------------

// AssumptionKind is the exact MVP assumption kind enum; the string values are
// the canonical JSON strings of §11.1.
type AssumptionKind string

// The exact seven assumption kinds. Approximation is reserved in MVP.
const (
	AssumptionDomain             AssumptionKind = "domain"
	AssumptionRegime             AssumptionKind = "regime"
	AssumptionConstraint         AssumptionKind = "constraint"
	AssumptionConvention         AssumptionKind = "convention"
	AssumptionApproximation      AssumptionKind = "approximation"
	AssumptionMathPrecondition   AssumptionKind = "math_precondition"
	AssumptionPhysicalAssumption AssumptionKind = "physical_assumption"
)

func validAssumptionKind(k AssumptionKind) bool {
	switch k {
	case AssumptionDomain, AssumptionRegime, AssumptionConstraint,
		AssumptionConvention, AssumptionApproximation,
		AssumptionMathPrecondition, AssumptionPhysicalAssumption:
		return true
	default:
		return false
	}
}

// ---------------------------------------------------------------------------
// ProvenanceStatus — exact §13.1 enum
// ---------------------------------------------------------------------------

// ProvenanceStatus is the provenance status enum. APPROXIMATED exists but no
// MVP operation mints it (REQ-013-01).
type ProvenanceStatus string

// The exact provenance statuses (canonical JSON = uppercase).
const (
	StatusDefined      ProvenanceStatus = "DEFINED"
	StatusPostulated   ProvenanceStatus = "POSTULATED"
	StatusDerived      ProvenanceStatus = "DERIVED"
	StatusIdentified   ProvenanceStatus = "IDENTIFIED"
	StatusApproximated ProvenanceStatus = "APPROXIMATED"
	StatusHypothesis   ProvenanceStatus = "HYPOTHESIS"
)

func validProvenanceStatus(s ProvenanceStatus) bool {
	switch s {
	case StatusDefined, StatusPostulated, StatusDerived,
		StatusIdentified, StatusApproximated, StatusHypothesis:
		return true
	default:
		return false
	}
}

// ---------------------------------------------------------------------------
// CorpusStatus — exact §13.4 enum
// ---------------------------------------------------------------------------

// CorpusStatus is the human-curated corpus status axis. The library never
// computes or revises it (REQ-013-04).
type CorpusStatus string

// The exact corpus statuses (canonical JSON = lowercase snake case).
const (
	CorpusNone        CorpusStatus = "NONE"
	CorpusEstablished CorpusStatus = "ESTABLISHED"
	CorpusContested   CorpusStatus = "CONTESTED"
	CorpusSuperseded  CorpusStatus = "SUPERSEDED"
	CorpusFalsified   CorpusStatus = "FALSIFIED"
)

// UnmarshalJSON implements custom JSON unmarshaling for CorpusStatus
// to accept lowercase snake_case values from manifest files.
func (c *CorpusStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "none":
		*c = CorpusNone
	case "established":
		*c = CorpusEstablished
	case "contested":
		*c = CorpusContested
	case "superseded":
		*c = CorpusSuperseded
	case "falsified":
		*c = CorpusFalsified
	default:
		return fmt.Errorf("invalid corpus status %q", s)
	}
	return nil
}

func validCorpusStatus(s CorpusStatus) bool {
	switch s {
	case CorpusNone, CorpusEstablished, CorpusContested, CorpusSuperseded, CorpusFalsified:
		return true
	default:
		return false
	}
}

// CorpusStatusJSON returns the canonical lowercase JSON spelling of s.
func CorpusStatusJSON(s CorpusStatus) string {
	switch s {
	case CorpusNone:
		return "none"
	case CorpusEstablished:
		return "established"
	case CorpusContested:
		return "contested"
	case CorpusSuperseded:
		return "superseded"
	case CorpusFalsified:
		return "falsified"
	default:
		return ""
	}
}

// ParseCorpusStatusJSON parses a canonical lowercase corpus status string.
func ParseCorpusStatusJSON(s string) (CorpusStatus, bool) {
	switch s {
	case "none":
		return CorpusNone, true
	case "established":
		return CorpusEstablished, true
	case "contested":
		return CorpusContested, true
	case "superseded":
		return CorpusSuperseded, true
	case "falsified":
		return CorpusFalsified, true
	default:
		return "", false
	}
}

// ---------------------------------------------------------------------------
// Hashing and canonical JSON helpers
// ---------------------------------------------------------------------------

// HashBytes returns the SHA-256 digest of b.
func HashBytes(b []byte) [32]byte { return sha256.Sum256(b) }

// Hex returns the lowercase hexadecimal form of h.
func Hex(h [32]byte) string { return hex.EncodeToString(h[:]) }

// isLowerHex64 reports whether s is exactly 64 lowercase hexadecimal
// characters.
func isLowerHex64(s string) bool {
	if len(s) != 64 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}

// marshalCanonical encodes v as compact canonical JSON without HTML escaping.
// Field order is fixed by struct declaration order (REQ-010-03); no map type
// is ever passed to it (REQ-010-02).
func marshalCanonical(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
}

// strictDecode decodes data into v with unknown-field rejection.
func strictDecode(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}
	return nil
}

// ratString returns the canonical "num/den" form: reduced, positive
// denominator, zero as "0/1" (§8.4, REQ-010-05).
func ratString(r *big.Rat) string {
	if r == nil {
		return ""
	}
	num := r.Num()
	den := r.Denom()
	// big.Rat keeps denominators positive; guard anyway for exactness.
	if den.Sign() < 0 {
		num = new(big.Int).Neg(num)
		den = new(big.Int).Neg(den)
	}
	return num.String() + "/" + den.String()
}

// parseRatExact parses a canonical "num/den" rational string and rejects any
// non-canonical spelling (floats, unreduced fractions, negative or zero
// denominators). This keeps exact rationals out of JSON numbers
// (REQ-010-04/05).
func parseRatExact(s string) (*big.Rat, error) {
	if s == "" {
		return nil, fmt.Errorf("empty rational")
	}
	parts := strings.Split(s, "/")
	if len(parts) > 2 {
		return nil, fmt.Errorf("malformed rational %q", s)
	}
	parseInt := func(p string) (*big.Int, error) {
		if p == "" {
			return nil, fmt.Errorf("malformed rational %q", s)
		}
		body := p
		if body[0] == '-' {
			body = body[1:]
		}
		if body == "" {
			return nil, fmt.Errorf("malformed rational %q", s)
		}
		for i := 0; i < len(body); i++ {
			if body[i] < '0' || body[i] > '9' {
				return nil, fmt.Errorf("malformed rational %q", s)
			}
		}
		n, ok := new(big.Int).SetString(p, 10)
		if !ok {
			return nil, fmt.Errorf("malformed rational %q", s)
		}
		return n, nil
	}
	num, err := parseInt(parts[0])
	if err != nil {
		return nil, err
	}
	den := big.NewInt(1)
	if len(parts) == 2 {
		den, err = parseInt(parts[1])
		if err != nil {
			return nil, err
		}
		if den.Sign() == 0 {
			return nil, fmt.Errorf("zero denominator in rational %q", s)
		}
	}
	r := new(big.Rat).SetFrac(num, den)
	// Canonicality: the input spelling must equal the canonical spelling.
	if ratString(r) != s {
		return nil, fmt.Errorf("non-canonical rational %q", s)
	}
	return r, nil
}

// ---------------------------------------------------------------------------
// Expr — immutable closed symbolic expression handle (§8)
// ---------------------------------------------------------------------------
//
// The zero Expr value is invalid (Valid() == false). An Expr stores only
// data: no callbacks, no functions, no raw strings (REQ-008-02).

// Expr is the immutable symbolic expression handle backed by a closed internal
// representation. Its fields are unexported; all accessors are kind-guarded
// and return defensive copies for mutable representations (*big.Rat, slices).
type Expr struct {
	kind       ExprKind
	name       string
	rat        *big.Rat // Rational value; Pow exponent
	children   []Expr   // Add/Mul terms/factors; Neg/Pow/Sqrt body; Call args; Relation lhs/rhs; BranchSet target+branches
	op         RelationOperator
	functionID string
}

// Valid reports structural validity of the expression handle. The zero value
// is invalid.
func (e Expr) Valid() bool {
	if !e.kind.valid() {
		return false
	}
	allValid := func(xs []Expr) bool {
		for _, x := range xs {
			if !x.Valid() {
				return false
			}
		}
		return true
	}
	switch e.kind {
	case ExprSymbol:
		return e.name != ""
	case ExprRational:
		return e.rat != nil
	case ExprAdd:
		for _, c := range e.children {
			if !c.Valid() || c.kind == ExprAdd {
				return false // nested same-kind nodes must be flattened
			}
		}
		return true
	case ExprMul:
		for _, c := range e.children {
			if !c.Valid() || c.kind == ExprMul {
				return false
			}
		}
		return true
	case ExprNeg:
		return len(e.children) == 1 && e.children[0].Valid()
	case ExprPow:
		return len(e.children) == 1 && e.children[0].Valid() && e.rat != nil
	case ExprSqrt:
		return len(e.children) == 1 && e.children[0].Valid()
	case ExprCall:
		return e.functionID == LorentzFactorFunctionID && allValid(e.children)
	case ExprRelation:
		return validRelationOperator(e.op) && len(e.children) == 2 && allValid(e.children)
	case ExprBranchSet:
		return len(e.children) >= 2 && allValid(e.children)
	default:
		return false
	}
}

// Kind returns the node kind (never a valid kind for invalid handles: the
// zero handle reports ExprSymbol, so callers must check Valid first).
func (e Expr) Kind() ExprKind { return e.kind }

// SymbolName returns the identifier for Symbol nodes; "" otherwise.
func (e Expr) SymbolName() string {
	if e.kind == ExprSymbol {
		return e.name
	}
	return ""
}

// RationalValue returns a defensive copy of the Rational value, or nil when
// the node is not a Rational.
func (e Expr) RationalValue() *big.Rat {
	if e.kind != ExprRational || e.rat == nil {
		return nil
	}
	return new(big.Rat).Set(e.rat)
}

// Children returns a copy of the node's child list, or nil when inapplicable.
func (e Expr) Children() []Expr {
	switch e.kind {
	case ExprAdd, ExprMul, ExprNeg, ExprPow, ExprSqrt, ExprCall, ExprRelation, ExprBranchSet:
		out := make([]Expr, len(e.children))
		copy(out, e.children)
		return out
	default:
		return nil
	}
}

// Base returns the Pow base, or the invalid Expr when inapplicable.
func (e Expr) Base() Expr {
	if e.kind != ExprPow || len(e.children) != 1 {
		return Expr{}
	}
	return e.children[0]
}

// Exponent returns a defensive copy of the Pow exponent, or nil.
func (e Expr) Exponent() *big.Rat {
	if e.kind != ExprPow || e.rat == nil {
		return nil
	}
	return new(big.Rat).Set(e.rat)
}

// FunctionID returns the Call function id, or "".
func (e Expr) FunctionID() string {
	if e.kind == ExprCall {
		return e.functionID
	}
	return ""
}

// Arguments returns a copy of the Call argument list, or nil.
func (e Expr) Arguments() []Expr {
	if e.kind != ExprCall {
		return nil
	}
	out := make([]Expr, len(e.children))
	copy(out, e.children)
	return out
}

// RelationOperator returns the Relation operator, or "" when inapplicable.
func (e Expr) RelationOperator() RelationOperator {
	if e.kind == ExprRelation {
		return e.op
	}
	return ""
}

// Left returns the Relation left side, or the invalid Expr.
func (e Expr) Left() Expr {
	if e.kind != ExprRelation || len(e.children) != 2 {
		return Expr{}
	}
	return e.children[0]
}

// Right returns the Relation right side, or the invalid Expr.
func (e Expr) Right() Expr {
	if e.kind != ExprRelation || len(e.children) != 2 {
		return Expr{}
	}
	return e.children[1]
}

// BranchTarget returns the BranchSet target, or the invalid Expr.
func (e Expr) BranchTarget() Expr {
	if e.kind != ExprBranchSet || len(e.children) < 2 {
		return Expr{}
	}
	return e.children[0]
}

// Branches returns a copy of the BranchSet branch list (excluding the
// target), or nil.
func (e Expr) Branches() []Expr {
	if e.kind != ExprBranchSet || len(e.children) < 2 {
		return nil
	}
	out := make([]Expr, len(e.children)-1)
	copy(out, e.children[1:])
	return out
}

// ---------------------------------------------------------------------------
// Expr construction — kernel-side primitives (§8.0, constructor-time
// structural canonicalization only: flattening, exact rational combination,
// sign normal form, deterministic child ordering — §9.6)
// ---------------------------------------------------------------------------

// NewSymbol builds a Symbol node. Symbol names are data, never code.
func NewSymbol(name string) (Expr, error) {
	if name == "" {
		return Expr{}, fmt.Errorf("symbol name must be non-empty")
	}
	return Expr{kind: ExprSymbol, name: name}, nil
}

// NewSymbolMust is like NewSymbol but panics on error.
func NewSymbolMust(name string) Expr {
	e, err := NewSymbol(name)
	if err != nil {
		panic(err)
	}
	return e
}

// NewRational builds a Rational node, copying the caller's *big.Rat so no
// caller-owned mutable pointer becomes part of immutable kernel state.
// A nil value yields the invalid Expr (callers must not pass nil).
func NewRational(value *big.Rat) Expr {
	if value == nil {
		return Expr{}
	}
	return Expr{kind: ExprRational, rat: new(big.Rat).Set(value)}
}

// NewRationalString parses an exact canonical "num/den" string.
func NewRationalString(s string) (Expr, error) {
	r, err := parseRatExact(s)
	if err != nil {
		return Expr{}, err
	}
	return Expr{kind: ExprRational, rat: r}, nil
}

// NewAdd builds an Add node: flattens nested Add children, combines exact
// rational terms (§9.4), and sorts children by the three-key rule (§9.3).
// Invalid children yield the invalid Expr.
func NewAdd(terms ...Expr) Expr {
	flat := make([]Expr, 0, len(terms))
	for _, t := range terms {
		if !t.Valid() {
			return Expr{}
		}
		if t.kind == ExprAdd {
			flat = append(flat, t.children...)
		} else {
			flat = append(flat, t)
		}
	}
	var sum *big.Rat
	rest := make([]Expr, 0, len(flat))
	for _, t := range flat {
		if t.kind == ExprRational {
			if sum == nil {
				sum = new(big.Rat)
			}
			sum.Add(sum, t.rat)
		} else {
			rest = append(rest, t)
		}
	}
	if sum == nil {
		if len(rest) == 0 {
			return Expr{kind: ExprAdd, children: []Expr{}}
		}
		sortExprChildren(rest)
		return Expr{kind: ExprAdd, children: rest}
	}
	if len(rest) == 0 {
		// Entire expression is a rational scalar.
		return Expr{kind: ExprRational, rat: sum}
	}
	children := append([]Expr{{kind: ExprRational, rat: sum}}, rest...)
	sortExprChildren(children)
	return Expr{kind: ExprAdd, children: children}
}

// NewMul builds a Mul node: flattens nested Mul children, combines exact
// rational coefficients (§9.4), applies the sign normal form (§9.5), and
// sorts children by the three-key rule. Invalid children yield the invalid
// Expr.
func NewMul(factors ...Expr) Expr {
	flat := make([]Expr, 0, len(factors))
	for _, f := range factors {
		if !f.Valid() {
			return Expr{}
		}
		if f.kind == ExprMul {
			flat = append(flat, f.children...)
		} else {
			flat = append(flat, f)
		}
	}
	if len(flat) == 0 {
		return Expr{kind: ExprMul, children: []Expr{}}
	}
	var prod *big.Rat
	rest := make([]Expr, 0, len(flat))
	for _, f := range flat {
		if f.kind == ExprRational {
			if prod == nil {
				prod = new(big.Rat).Set(f.rat)
			} else {
				prod.Mul(prod, f.rat)
			}
		} else {
			rest = append(rest, f)
		}
	}
	if prod == nil {
		sortExprChildren(rest)
		return Expr{kind: ExprMul, children: rest}
	}
	if len(rest) == 0 {
		// Entire expression is a rational scalar.
		return Expr{kind: ExprRational, rat: prod}
	}
	if prod.Sign() < 0 {
		// Sign normal form (§9.5): a negative scalar coefficient on a
		// symbolic product becomes Neg around the positive-coefficient
		// product; Mul(-1, x...) → Neg(x...) absorbs the unit coefficient.
		absProd := new(big.Rat).Abs(prod)
		if absProd.Cmp(big.NewRat(1, 1)) == 0 {
			if len(rest) == 1 {
				return Expr{kind: ExprNeg, children: []Expr{rest[0]}}
			}
			sortExprChildren(rest)
			return Expr{kind: ExprNeg, children: []Expr{{kind: ExprMul, children: rest}}}
		}
		children := append([]Expr{{kind: ExprRational, rat: absProd}}, rest...)
		sortExprChildren(children)
		mul := Expr{kind: ExprMul, children: children}
		return Expr{kind: ExprNeg, children: []Expr{mul}}
	}
	children := append([]Expr{{kind: ExprRational, rat: prod}}, rest...)
	sortExprChildren(children)
	return Expr{kind: ExprMul, children: children}
}

// NewNeg builds a Neg node with sign normal form: Neg(Neg(x)) → x and
// Neg(Rational(q)) → Rational(-q) (§9.5).
func NewNeg(e Expr) Expr {
	if !e.Valid() {
		return Expr{}
	}
	switch e.kind {
	case ExprNeg:
		return e.children[0]
	case ExprRational:
		return Expr{kind: ExprRational, rat: new(big.Rat).Neg(e.rat)}
	default:
		return Expr{kind: ExprNeg, children: []Expr{e}}
	}
}

// NewPow builds a Pow node, copying the exact rational exponent. Algebraic
// rewrites (Pow(x,1), Pow(x,0), …) are Simplify-time rules (§9.6).
func NewPow(base Expr, exponent *big.Rat) Expr {
	if !base.Valid() || exponent == nil {
		return Expr{}
	}
	return Expr{kind: ExprPow, children: []Expr{base}, rat: new(big.Rat).Set(exponent)}
}

// NewSqrt builds a Sqrt node.
func NewSqrt(e Expr) Expr {
	if !e.Valid() {
		return Expr{}
	}
	return Expr{kind: ExprSqrt, children: []Expr{e}}
}

// NewCall builds the sole MVP Call form: lorentz_factor. No arbitrary
// function registry exists (§8.9).
func NewCall(functionID string, args ...Expr) (Expr, error) {
	if functionID != LorentzFactorFunctionID {
		return Expr{}, fmt.Errorf("unsupported function id %q", functionID)
	}
	children := make([]Expr, len(args))
	copy(children, args)
	for _, a := range children {
		if !a.Valid() {
			return Expr{}, fmt.Errorf("invalid call argument")
		}
	}
	return Expr{kind: ExprCall, functionID: functionID, children: children}, nil
}

// NewRelation builds a Relation node, preserving operator and side order
// (§9.9). Invalid operator or operands yield the invalid Expr.
func NewRelation(op RelationOperator, lhs, rhs Expr) Expr {
	if !validRelationOperator(op) || !lhs.Valid() || !rhs.Valid() {
		return Expr{}
	}
	return Expr{kind: ExprRelation, op: op, children: []Expr{lhs, rhs}}
}

// NewBranchSet builds a bounded BranchSet: children[0] is the target,
// children[1:] are the branches in declared order (the first branch is the
// structurally nonnegative root selected by SelectBranch).
func NewBranchSet(target Expr, branches ...Expr) Expr {
	if !target.Valid() || len(branches) == 0 {
		return Expr{}
	}
	children := make([]Expr, 0, len(branches)+1)
	children = append(children, target)
	for _, b := range branches {
		if !b.Valid() {
			return Expr{}
		}
		children = append(children, b)
	}
	return Expr{kind: ExprBranchSet, children: children}
}

// ---------------------------------------------------------------------------
// Child sorting — three-key rule (§9.3)
// ---------------------------------------------------------------------------

// sortExprChildren orders Add/Mul children by (1) node-kind ordinal, (2)
// lowercase hex SHA-256 child hash via bytes.Compare, (3) canonical child
// bytes via bytes.Compare. Never maps, pointers, or locale.
func sortExprChildren(children []Expr) {
	type sortEntry struct {
		ord   uint8
		hash  [32]byte
		canon []byte
		e     Expr
	}
	entries := make([]sortEntry, len(children))
	for i, c := range children {
		canon, err := CanonicalExprJSON(c)
		if err != nil {
			canon = nil
		}
		entries[i] = sortEntry{ord: uint8(c.kind), hash: HashBytes(canon), canon: canon, e: c}
	}
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].ord != entries[j].ord {
			return entries[i].ord < entries[j].ord
		}
		if c := bytes.Compare(entries[i].hash[:], entries[j].hash[:]); c != 0 {
			return c < 0
		}
		return bytes.Compare(entries[i].canon, entries[j].canon) < 0
	})
	for i, en := range entries {
		children[i] = en.e
	}
}

// ---------------------------------------------------------------------------
// Canonical expression JSON (§10.2)
// ---------------------------------------------------------------------------

type exprSymbolDTO struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type exprRationalDTO struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type exprAddDTO struct {
	Kind  string            `json:"kind"`
	Terms []json.RawMessage `json:"terms"`
}

type exprMulDTO struct {
	Kind    string            `json:"kind"`
	Factors []json.RawMessage `json:"factors"`
}

type exprNegDTO struct {
	Kind string          `json:"kind"`
	Expr json.RawMessage `json:"expr"`
}

type exprPowDTO struct {
	Kind string          `json:"kind"`
	Base json.RawMessage `json:"base"`
	Exp  string          `json:"exp"`
}

type exprSqrtDTO struct {
	Kind string          `json:"kind"`
	Expr json.RawMessage `json:"expr"`
}

type exprCallDTO struct {
	Kind     string            `json:"kind"`
	Function string            `json:"function"`
	Args     []json.RawMessage `json:"args"`
}

type exprRelationDTO struct {
	Kind     string           `json:"kind"`
	Operator RelationOperator `json:"operator"`
	LHS      json.RawMessage  `json:"lhs"`
	RHS      json.RawMessage  `json:"rhs"`
}

type exprBranchSetDTO struct {
	Kind     string            `json:"kind"`
	Target   json.RawMessage   `json:"target"`
	Branches []json.RawMessage `json:"branches"`
}

// CanonicalExprJSON encodes e in its exact §10.2 canonical form. Invalid
// expressions yield an error; encoding never sorts or rewrites (internal
// state is already canonical by construction).
func CanonicalExprJSON(e Expr) ([]byte, error) {
	if !e.Valid() {
		return nil, fmt.Errorf("cannot encode invalid expression")
	}
	switch e.kind {
	case ExprSymbol:
		return marshalCanonical(exprSymbolDTO{Kind: "symbol", Name: e.name})
	case ExprRational:
		return marshalCanonical(exprRationalDTO{Kind: "rational", Value: ratString(e.rat)})
	case ExprAdd:
		terms := rawChildren(e.children)
		return marshalCanonical(exprAddDTO{Kind: "add", Terms: terms})
	case ExprMul:
		factors := rawChildren(e.children)
		return marshalCanonical(exprMulDTO{Kind: "mul", Factors: factors})
	case ExprNeg:
		body, err := CanonicalExprJSON(e.children[0])
		if err != nil {
			return nil, err
		}
		return marshalCanonical(exprNegDTO{Kind: "neg", Expr: body})
	case ExprPow:
		base, err := CanonicalExprJSON(e.children[0])
		if err != nil {
			return nil, err
		}
		return marshalCanonical(exprPowDTO{Kind: "pow", Base: base, Exp: ratString(e.rat)})
	case ExprSqrt:
		body, err := CanonicalExprJSON(e.children[0])
		if err != nil {
			return nil, err
		}
		return marshalCanonical(exprSqrtDTO{Kind: "sqrt", Expr: body})
	case ExprCall:
		args := rawChildren(e.children)
		return marshalCanonical(exprCallDTO{Kind: "call", Function: e.functionID, Args: args})
	case ExprRelation:
		lhs, err := CanonicalExprJSON(e.children[0])
		if err != nil {
			return nil, err
		}
		rhs, err := CanonicalExprJSON(e.children[1])
		if err != nil {
			return nil, err
		}
		return marshalCanonical(exprRelationDTO{Kind: "relation", Operator: e.op, LHS: lhs, RHS: rhs})
	case ExprBranchSet:
		target, err := CanonicalExprJSON(e.children[0])
		if err != nil {
			return nil, err
		}
		branches := rawChildren(e.children[1:])
		return marshalCanonical(exprBranchSetDTO{Kind: "branch_set", Target: target, Branches: branches})
	default:
		return nil, fmt.Errorf("unknown expression kind")
	}
}

func rawChildren(children []Expr) []json.RawMessage {
	out := make([]json.RawMessage, len(children))
	for i, c := range children {
		b, err := CanonicalExprJSON(c)
		if err != nil {
			// Children of a Valid parent are always encodable.
			b = nil
		}
		out[i] = json.RawMessage(b)
	}
	return out
}

// ParseExprJSON decodes canonical expression JSON into the closed node set,
// validating node invariants and byte-level canonicality:
// Encode(Decode(data)) == data. It is exported for in-module packages only
// (internal/kernel is unreachable from outside the module).
func ParseExprJSON(data []byte) (Expr, error) {
	e, err := decodeExpr(data)
	if err != nil {
		return Expr{}, err
	}
	if !e.Valid() {
		return Expr{}, fmt.Errorf("decoded expression violates invariants")
	}
	canonical, err := CanonicalExprJSON(e)
	if err != nil {
		return Expr{}, err
	}
	if !bytes.Equal(canonical, data) {
		return Expr{}, fmt.Errorf("non-canonical expression JSON")
	}
	return e, nil
}

// MarshalJSON implements json.Marshaler using CanonicalExprJSON.
func (e Expr) MarshalJSON() ([]byte, error) {
	return CanonicalExprJSON(e)
}

func decodeExpr(data []byte) (Expr, error) {
	var head struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(data, &head); err != nil {
		return Expr{}, fmt.Errorf("expression kind: %w", err)
	}
	switch head.Kind {
	case "symbol":
		var dto exprSymbolDTO
		if err := strictDecode(data, &dto); err != nil {
			return Expr{}, err
		}
		return NewSymbol(dto.Name)
	case "rational":
		var dto exprRationalDTO
		if err := strictDecode(data, &dto); err != nil {
			return Expr{}, err
		}
		return NewRationalString(dto.Value)
	case "add":
		var dto exprAddDTO
		if err := strictDecode(data, &dto); err != nil {
			return Expr{}, err
		}
		terms, err := decodeRawChildren(dto.Terms)
		if err != nil {
			return Expr{}, err
		}
		return NewAdd(terms...), nil
	case "mul":
		var dto exprMulDTO
		if err := strictDecode(data, &dto); err != nil {
			return Expr{}, err
		}
		factors, err := decodeRawChildren(dto.Factors)
		if err != nil {
			return Expr{}, err
		}
		return NewMul(factors...), nil
	case "neg":
		var dto exprNegDTO
		if err := strictDecode(data, &dto); err != nil {
			return Expr{}, err
		}
		body, err := decodeExpr(dto.Expr)
		if err != nil {
			return Expr{}, err
		}
		return NewNeg(body), nil
	case "pow":
		var dto exprPowDTO
		if err := strictDecode(data, &dto); err != nil {
			return Expr{}, err
		}
		base, err := decodeExpr(dto.Base)
		if err != nil {
			return Expr{}, err
		}
		exp, err := parseRatExact(dto.Exp)
		if err != nil {
			return Expr{}, err
		}
		return NewPow(base, exp), nil
	case "sqrt":
		var dto exprSqrtDTO
		if err := strictDecode(data, &dto); err != nil {
			return Expr{}, err
		}
		body, err := decodeExpr(dto.Expr)
		if err != nil {
			return Expr{}, err
		}
		return NewSqrt(body), nil
	case "call":
		var dto exprCallDTO
		if err := strictDecode(data, &dto); err != nil {
			return Expr{}, err
		}
		args, err := decodeRawChildren(dto.Args)
		if err != nil {
			return Expr{}, err
		}
		return NewCall(dto.Function, args...)
	case "relation":
		var dto exprRelationDTO
		if err := strictDecode(data, &dto); err != nil {
			return Expr{}, err
		}
		if !validRelationOperator(dto.Operator) {
			return Expr{}, fmt.Errorf("invalid relation operator %q", string(dto.Operator))
		}
		lhs, err := decodeExpr(dto.LHS)
		if err != nil {
			return Expr{}, err
		}
		rhs, err := decodeExpr(dto.RHS)
		if err != nil {
			return Expr{}, err
		}
		return NewRelation(dto.Operator, lhs, rhs), nil
	case "branch_set":
		var dto exprBranchSetDTO
		if err := strictDecode(data, &dto); err != nil {
			return Expr{}, err
		}
		target, err := decodeExpr(dto.Target)
		if err != nil {
			return Expr{}, err
		}
		branches, err := decodeRawChildren(dto.Branches)
		if err != nil {
			return Expr{}, err
		}
		return NewBranchSet(target, branches...), nil
	default:
		return Expr{}, fmt.Errorf("unknown expression kind %q", head.Kind)
	}
}

func decodeRawChildren(raws []json.RawMessage) ([]Expr, error) {
	out := make([]Expr, len(raws))
	for i, r := range raws {
		e, err := decodeExpr(r)
		if err != nil {
			return nil, err
		}
		out[i] = e
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Expression equality and hashing (§9.1, §9.2)
// ---------------------------------------------------------------------------

// EqualExpr compares canonical symbolic structure (never display strings).
// Invalid expressions are never equal to anything, including each other.
func EqualExpr(a, b Expr) bool {
	if !a.Valid() || !b.Valid() {
		return false
	}
	if a.kind != b.kind {
		return false
	}
	switch a.kind {
	case ExprSymbol:
		return a.name == b.name
	case ExprRational:
		return a.rat.Cmp(b.rat) == 0
	case ExprAdd, ExprMul:
		if len(a.children) != len(b.children) {
			return false
		}
		for i := range a.children {
			if !EqualExpr(a.children[i], b.children[i]) {
				return false
			}
		}
		return true
	case ExprNeg, ExprSqrt:
		return EqualExpr(a.children[0], b.children[0])
	case ExprPow:
		return EqualExpr(a.children[0], b.children[0]) && a.rat.Cmp(b.rat) == 0
	case ExprCall:
		if a.functionID != b.functionID || len(a.children) != len(b.children) {
			return false
		}
		for i := range a.children {
			if !EqualExpr(a.children[i], b.children[i]) {
				return false
			}
		}
		return true
	case ExprRelation:
		return a.op == b.op &&
			EqualExpr(a.children[0], b.children[0]) &&
			EqualExpr(a.children[1], b.children[1])
	case ExprBranchSet:
		if len(a.children) != len(b.children) {
			return false
		}
		for i := range a.children {
			if !EqualExpr(a.children[i], b.children[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// HashExpr returns the SHA-256 of the expression's canonical JSON bytes;
// invalid expressions hash to the zero digest.
func HashExpr(e Expr) [32]byte {
	b, err := CanonicalExprJSON(e)
	if err != nil {
		return [32]byte{}
	}
	return HashBytes(b)
}

// ---------------------------------------------------------------------------
// Dimension — seven exact rational SI exponents (§7)
// ---------------------------------------------------------------------------

// Dimension is the immutable seven-base-dimension exponent vector
// (M, L, T, I, Θ, N, J) backed by math/big.Rat. The zero value is invalid.
type Dimension struct {
	exp [7]*big.Rat
}

// NewDimension builds a Dimension from exact rational exponents in the order
// M, L, T, I, Θ, N, J. Inputs are copied and canonicalized (reduced, positive
// denominator). Nil exponents or nil vectors are rejected.
func NewDimension(m, l, t, i, theta, n, j *big.Rat) (Dimension, error) {
	exps := [7]*big.Rat{m, l, t, i, theta, n, j}
	return newDimension(exps)
}

func newDimension(exps [7]*big.Rat) (Dimension, error) {
	out := [7]*big.Rat{}
	for idx, e := range exps {
		if e == nil {
			return Dimension{}, fmt.Errorf("nil dimension exponent at index %d", idx)
		}
		r := new(big.Rat).Set(e)
		if r.Denom().Sign() < 0 {
			r.SetFrac(new(big.Int).Neg(r.Num()), new(big.Int).Neg(r.Denom()))
		}
		out[idx] = r
	}
	return Dimension{exp: out}, nil
}

// Valid reports whether all seven exponents are present.
func (d Dimension) Valid() bool {
	for _, e := range d.exp {
		if e == nil {
			return false
		}
	}
	return true
}

// Equal reports exact structural equality (REQ-007-02).
func (d Dimension) Equal(other Dimension) bool {
	if !d.Valid() || !other.Valid() {
		return false
	}
	for i := 0; i < 7; i++ {
		if d.exp[i].Cmp(other.exp[i]) != 0 {
			return false
		}
	}
	return true
}

// Multiply returns the componentwise product (exponent sum).
func (d Dimension) Multiply(other Dimension) Dimension {
	if !d.Valid() || !other.Valid() {
		return Dimension{}
	}
	var out [7]*big.Rat
	for i := 0; i < 7; i++ {
		out[i] = new(big.Rat).Add(d.exp[i], other.exp[i])
	}
	dim, _ := newDimension(out)
	return dim
}

// Divide returns the componentwise quotient (exponent difference).
func (d Dimension) Divide(other Dimension) Dimension {
	if !d.Valid() || !other.Valid() {
		return Dimension{}
	}
	var out [7]*big.Rat
	for i := 0; i < 7; i++ {
		out[i] = new(big.Rat).Sub(d.exp[i], other.exp[i])
	}
	dim, _ := newDimension(out)
	return dim
}

// Pow multiplies every base exponent by the exact rational exponent
// (§7.3).
func (d Dimension) Pow(exponent *big.Rat) Dimension {
	if !d.Valid() || exponent == nil {
		return Dimension{}
	}
	var out [7]*big.Rat
	for i := 0; i < 7; i++ {
		out[i] = new(big.Rat).Mul(d.exp[i], exponent)
	}
	dim, _ := newDimension(out)
	return dim
}

type dimensionDTO struct {
	M     string `json:"m"`
	L     string `json:"l"`
	T     string `json:"t"`
	I     string `json:"i"`
	Theta string `json:"theta"`
	N     string `json:"n"`
	J     string `json:"j"`
}

// CanonicalJSON encodes the dimension in the exact §10.3 field order using
// canonical rational strings (never JSON numbers).
func (d Dimension) CanonicalJSON() ([]byte, error) {
	if !d.Valid() {
		return nil, fmt.Errorf("cannot encode invalid dimension")
	}
	return marshalCanonical(dimensionDTO{
		M: ratString(d.exp[0]), L: ratString(d.exp[1]), T: ratString(d.exp[2]),
		I: ratString(d.exp[3]), Theta: ratString(d.exp[4]), N: ratString(d.exp[5]),
		J: ratString(d.exp[6]),
	})
}

// MarshalJSON implements json.Marshaler using CanonicalJSON.
func (d Dimension) MarshalJSON() ([]byte, error) {
	return d.CanonicalJSON()
}

// Hash returns the SHA-256 of the canonical dimension JSON; invalid
// dimensions hash to the zero digest.
func (d Dimension) Hash() [32]byte {
	b, err := d.CanonicalJSON()
	if err != nil {
		return [32]byte{}
	}
	return HashBytes(b)
}

// ParseDimensionJSON decodes canonical dimension JSON, enforcing canonical
// rational spellings (round-trip safe).
func ParseDimensionJSON(data []byte) (Dimension, error) {
	var dto dimensionDTO
	if err := strictDecode(data, &dto); err != nil {
		return Dimension{}, err
	}
	fields := [7]string{dto.M, dto.L, dto.T, dto.I, dto.Theta, dto.N, dto.J}
	var exps [7]*big.Rat
	for idx, f := range fields {
		r, err := parseRatExact(f)
		if err != nil {
			return Dimension{}, fmt.Errorf("dimension field %d: %w", idx, err)
		}
		exps[idx] = r
	}
	d, err := newDimension(exps)
	if err != nil {
		return Dimension{}, err
	}
	canonical, err := d.CanonicalJSON()
	if err != nil {
		return Dimension{}, err
	}
	if !bytes.Equal(canonical, data) {
		return Dimension{}, fmt.Errorf("non-canonical dimension JSON")
	}
	return d, nil
}

// UnmarshalJSON implements custom JSON unmarshaling for Dimension
// to accept the canonical rational string format (without strict canonicality check).
func (d *Dimension) UnmarshalJSON(data []byte) error {
	var dto dimensionDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return err
	}
	fields := [7]string{dto.M, dto.L, dto.T, dto.I, dto.Theta, dto.N, dto.J}
	var exps [7]*big.Rat
	for idx, f := range fields {
		r, err := parseRatExact(f)
		if err != nil {
			return fmt.Errorf("dimension field %d: %w", idx, err)
		}
		exps[idx] = r
	}
	*d = Dimension{exp: exps}
	return nil
}

// ---------------------------------------------------------------------------
// Assumption — structured metadata (§11)
// ---------------------------------------------------------------------------

// Assumption is immutable structured metadata: Kind, Key, and a closed
// TextValue(string) | ExprValue(Expr) union. Assumptions are metadata, not
// trusted physical objects.
type Assumption struct {
	kind     AssumptionKind
	key      string
	exprMode bool
	text     string
	expr     Expr
}

// NewTextAssumption builds a text-valued assumption. Kind, key, and value
// must be non-empty.
func NewTextAssumption(kind AssumptionKind, key, value string) (Assumption, error) {
	if !validAssumptionKind(kind) {
		return Assumption{}, fmt.Errorf("invalid assumption kind %q", string(kind))
	}
	if key == "" {
		return Assumption{}, fmt.Errorf("assumption key must be non-empty")
	}
	if value == "" {
		return Assumption{}, fmt.Errorf("text assumption value must be non-empty")
	}
	return Assumption{kind: kind, key: key, exprMode: false, text: value}, nil
}

// NewTextAssumptionMust is like NewTextAssumption but panics on error.
// For use in manifest/constructor code where inputs are known valid.
func NewTextAssumptionMust(kind AssumptionKind, key, value string) Assumption {
	a, err := NewTextAssumption(kind, key, value)
	if err != nil {
		panic(err)
	}
	return a
}

// NewExprAssumption builds an expression-valued assumption. Mathematical
// assumptions MUST use structured expression values, never equation strings
// (REQ-011-01).
func NewExprAssumption(kind AssumptionKind, key string, value Expr) (Assumption, error) {
	if !validAssumptionKind(kind) {
		return Assumption{}, fmt.Errorf("invalid assumption kind %q", string(kind))
	}
	if key == "" {
		return Assumption{}, fmt.Errorf("assumption key must be non-empty")
	}
	if !value.Valid() {
		return Assumption{}, fmt.Errorf("assumption expression must be valid")
	}
	return Assumption{kind: kind, key: key, exprMode: true, expr: value}, nil
}

// NewExprAssumptionMust is like NewExprAssumption but panics on error.
func NewExprAssumptionMust(kind AssumptionKind, key string, value Expr) Assumption {
	a, err := NewExprAssumption(kind, key, value)
	if err != nil {
		panic(err)
	}
	return a
}

// Kind returns the assumption kind.
func (a Assumption) Kind() AssumptionKind { return a.kind }

// Key returns the non-empty assumption key.
func (a Assumption) Key() string { return a.key }

// Mode returns the structured value mode: "text" or "expr".
func (a Assumption) Mode() string {
	if a.exprMode {
		return "expr"
	}
	return "text"
}

// Text returns the text-mode content ("" in expr mode).
func (a Assumption) Text() string {
	if a.exprMode {
		return ""
	}
	return a.text
}

// Expr returns the expr-mode content (the invalid Expr in text mode).
func (a Assumption) Expr() Expr {
	if !a.exprMode {
		return Expr{}
	}
	return a.expr
}

func (a Assumption) valid() bool {
	if !validAssumptionKind(a.kind) || a.key == "" {
		return false
	}
	if a.exprMode {
		return a.expr.Valid()
	}
	return a.text != ""
}

type assumptionValueDTO struct {
	Mode  string          `json:"mode"`
	Value string          `json:"value,omitempty"`
	Expr  json.RawMessage `json:"expr,omitempty"`
}

type assumptionDTO struct {
	Kind  string             `json:"kind"`
	Key   string             `json:"key"`
	Value assumptionValueDTO `json:"value"`
}

func canonicalAssumptionJSON(a Assumption) ([]byte, error) {
	if !a.valid() {
		return nil, fmt.Errorf("invalid assumption")
	}
	var val assumptionValueDTO
	if a.exprMode {
		body, err := CanonicalExprJSON(a.expr)
		if err != nil {
			return nil, err
		}
		val = assumptionValueDTO{Mode: "expr", Expr: body}
	} else {
		val = assumptionValueDTO{Mode: "text", Value: a.text}
	}
	return marshalCanonical(assumptionDTO{Kind: string(a.kind), Key: a.key, Value: val})
}

func decodeAssumption(data []byte) (Assumption, error) {
	var dto assumptionDTO
	if err := strictDecode(data, &dto); err != nil {
		return Assumption{}, err
	}
	switch dto.Value.Mode {
	case "text":
		return NewTextAssumption(AssumptionKind(dto.Kind), dto.Key, dto.Value.Value)
	case "expr":
		e, err := ParseExprJSON(dto.Value.Expr)
		if err != nil {
			return Assumption{}, err
		}
		return NewExprAssumption(AssumptionKind(dto.Kind), dto.Key, e)
	default:
		return Assumption{}, fmt.Errorf("invalid assumption value mode %q", dto.Value.Mode)
	}
}

// ParseAssumptionJSON parses a canonical assumption JSON into an Assumption.
// Exported for use by core/corpus.go manifest loading.
func ParseAssumptionJSON(data []byte) (Assumption, error) {
	return decodeAssumption(data)
}

// Valid reports whether the assumption is well-formed.
func (a Assumption) Valid() bool {
	return a.valid()
}

// UnmarshalJSON implements json.Unmarshaler for Assumption.
func (a *Assumption) UnmarshalJSON(data []byte) error {
	parsed, err := decodeAssumption(data)
	if err != nil {
		return err
	}
	*a = parsed
	return nil
}

// MarshalJSON implements json.Marshaler using canonical assumption JSON.
func (a Assumption) MarshalJSON() ([]byte, error) {
	return canonicalAssumptionJSON(a)
}

// ---------------------------------------------------------------------------
// AssumptionSet — immutable deterministic set (§11.0, §11.4; plan F4)
// ---------------------------------------------------------------------------

// AssumptionSet is an immutable deterministic set of assumptions. Elements
// are canonicalized, sorted lexicographically by canonical bytes, and
// deduplicated at construction; map iteration never determines order.
type AssumptionSet struct {
	values []Assumption // sorted by canonical bytes, deduplicated
}

// NewAssumptionSet canonicalizes, sorts, and deduplicates values (plan F4).
// Invalid entries are dropped; the constructor cannot fail (§11.0).
func NewAssumptionSet(values ...Assumption) AssumptionSet {
	type entry struct {
		a Assumption
		b []byte
	}
	entries := make([]entry, 0, len(values))
	for _, v := range values {
		if !v.valid() {
			continue
		}
		b, err := canonicalAssumptionJSON(v)
		if err != nil {
			continue
		}
		entries = append(entries, entry{a: v, b: b})
	}
	sort.Slice(entries, func(i, j int) bool {
		return bytes.Compare(entries[i].b, entries[j].b) < 0
	})
	out := make([]Assumption, 0, len(entries))
	var prev []byte
	for i, e := range entries {
		if i > 0 && bytes.Equal(prev, e.b) {
			continue // exact duplicate
		}
		out = append(out, e.a)
		prev = e.b
	}
	return AssumptionSet{values: out}
}

// Values returns a copy of the canonical element slice.
func (s AssumptionSet) Values() []Assumption {
	out := make([]Assumption, len(s.values))
	copy(out, s.values)
	return out
}

// Merge returns the deterministic union of both sets with exact conflict
// detection: same (Kind, Key) with different canonical values yields
// AssumptionConflictError (§11.4, MRC-004).
func (s AssumptionSet) Merge(other AssumptionSet) (AssumptionSet, error) {
	type assocKey struct {
		kind AssumptionKind
		key  string
	}
	all := make([]Assumption, 0, len(s.values)+len(other.values))
	all = append(all, s.values...)
	all = append(all, other.values...)
	seen := make(map[assocKey][]byte, len(all))
	for _, a := range all {
		b, err := canonicalAssumptionJSON(a)
		if err != nil {
			return AssumptionSet{}, err
		}
		k := assocKey{kind: a.kind, key: a.key}
		if prev, ok := seen[k]; ok {
			if !bytes.Equal(prev, b) {
				return AssumptionSet{}, AssumptionConflictError{
					Kind: string(a.kind), Key: a.key,
					Left: string(prev), Right: string(b),
				}
			}
			continue
		}
		seen[k] = b
	}
	return NewAssumptionSet(all...), nil
}

// Equal reports canonical set equality.
func (s AssumptionSet) Equal(other AssumptionSet) bool {
	a, errA := s.CanonicalJSON()
	b, errB := other.CanonicalJSON()
	if errA != nil || errB != nil {
		return false
	}
	return bytes.Equal(a, b)
}

// CanonicalJSON encodes the set as its canonical element array.
func (s AssumptionSet) CanonicalJSON() ([]byte, error) {
	// Marshal element bytes directly for a canonical array encoding.
	parts := make([]json.RawMessage, 0, len(s.values))
	for _, a := range s.values {
		body, err := canonicalAssumptionJSON(a)
		if err != nil {
			return nil, err
		}
		parts = append(parts, json.RawMessage(body))
	}
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, p := range parts {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.Write(p)
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

// Hash returns the SHA-256 of the canonical set JSON.
func (s AssumptionSet) Hash() [32]byte {
	b, err := s.CanonicalJSON()
	if err != nil {
		return [32]byte{}
	}
	return HashBytes(b)
}

// ParseAssumptionSetJSON decodes a canonical assumption-set array.
func ParseAssumptionSetJSON(data []byte) (AssumptionSet, error) {
	var raws []json.RawMessage
	if err := json.Unmarshal(data, &raws); err != nil {
		return AssumptionSet{}, fmt.Errorf("assumption set: %w", err)
	}
	values := make([]Assumption, 0, len(raws))
	for _, r := range raws {
		a, err := decodeAssumption(r)
		if err != nil {
			return AssumptionSet{}, err
		}
		values = append(values, a)
	}
	set := NewAssumptionSet(values...)
	canonical, err := set.CanonicalJSON()
	if err != nil {
		return AssumptionSet{}, err
	}
	if !bytes.Equal(canonical, data) {
		return AssumptionSet{}, fmt.Errorf("non-canonical assumption set JSON")
	}
	return set, nil
}

// ---------------------------------------------------------------------------
// Bounded sign entailment (§9.8) — kernel-owned helpers; "no other logical
// inference is permitted"
// ---------------------------------------------------------------------------

// EntailsNonNegative implements the closed §9.8 nonnegativity rules.
func (s AssumptionSet) EntailsNonNegative(e Expr) bool {
	if !e.Valid() {
		return false
	}
	switch e.Kind() {
	case ExprRational:
		return e.rat.Sign() >= 0
	case ExprMul:
		if len(e.children) == 0 {
			return false
		}
		for _, f := range e.children {
			if !s.EntailsNonNegative(f) {
				return false
			}
		}
		return true
	case ExprPow:
		n := e.rat
		if n == nil || n.Sign() < 0 || !n.IsInt() {
			return false
		}
		// Even nonnegative integer exponent → nonnegative.
		return n.Num().Bit(0) == 0
	default:
		return s.directInequality(e, RelationGte, false) ||
			s.directInequality(e, RelationLte, true)
	}
}

// EntailsPositive implements the closed §9.8 positivity rules; positive
// entails nonzero.
func (s AssumptionSet) EntailsPositive(e Expr) bool {
	if !e.Valid() {
		return false
	}
	switch e.Kind() {
	case ExprRational:
		return e.rat.Sign() > 0
	case ExprMul:
		if len(e.children) == 0 {
			return false
		}
		anyPositive := false
		for _, f := range e.children {
			if s.EntailsPositive(f) {
				anyPositive = true
				continue
			}
			if !s.EntailsNonNegative(f) {
				return false
			}
		}
		return anyPositive
	default:
		return s.directInequality(e, RelationGt, false) ||
			s.directInequality(e, RelationLt, true)
	}
}

// EntailsNonZero implements the closed §9.8 nonzero rules: positive entails
// nonzero, plus a direct structural match of an assumption expr
// Relation(neq, e, 0).
func (s AssumptionSet) EntailsNonZero(e Expr) bool {
	if s.EntailsPositive(e) {
		return true
	}
	return s.directInequality(e, RelationNeq, false)
}

// directInequality reports whether some expression-valued assumption is a
// direct structural inequality of the requested shape against e:
//
//	flipped=false (op gte/gt/neq): Relation(op, e, 0)
//	flipped=true  (op lte/lt):     Relation(op, 0, e)
//
// This is the only structural match permitted; no derived inference exists.
func (s AssumptionSet) directInequality(e Expr, op RelationOperator, flipped bool) bool {
	for _, a := range s.values {
		if !a.exprMode {
			continue
		}
		v := a.expr
		if !v.Valid() || v.kind != ExprRelation || v.op != op || len(v.children) != 2 {
			continue
		}
		lhs, rhs := v.children[0], v.children[1]
		if flipped {
			if lhs.kind == ExprRational && lhs.rat.Sign() == 0 && EqualExpr(rhs, e) {
				return true
			}
			continue
		}
		if rhs.kind == ExprRational && rhs.rat.Sign() == 0 && EqualExpr(lhs, e) {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Convention — immutable key/value metadata (§12)
// ---------------------------------------------------------------------------

// Convention is immutable metadata: both Key and Value are non-empty
// strings. No physical equation may be encoded solely as a convention string
// (REQ-012-02 — the type system permits only strings).
type Convention struct {
	key   string
	value string
}

// NewConvention builds a convention; both fields must be non-empty.
func NewConvention(key, value string) (Convention, error) {
	if key == "" {
		return Convention{}, fmt.Errorf("convention key must be non-empty")
	}
	if value == "" {
		return Convention{}, fmt.Errorf("convention value must be non-empty")
	}
	return Convention{key: key, value: value}, nil
}

// NewConventionMust is like NewConvention but panics on error.
// For use in manifest/constructor code where inputs are known valid.
func NewConventionMust(key, value string) Convention {
	c, err := NewConvention(key, value)
	if err != nil {
		panic(err)
	}
	return c
}

// Key returns the convention key.
func (c Convention) Key() string { return c.key }

// Value returns the convention value.
func (c Convention) Value() string { return c.value }

func (c Convention) valid() bool { return c.key != "" && c.value != "" }

type conventionDTO struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func canonicalConventionJSON(c Convention) ([]byte, error) {
	if !c.valid() {
		return nil, fmt.Errorf("invalid convention")
	}
	return marshalCanonical(conventionDTO{Key: c.key, Value: c.value})
}

// ---------------------------------------------------------------------------
// ConventionSet — immutable deterministic set (§12.0; plan F4)
// ---------------------------------------------------------------------------

// ConventionSet is an immutable deterministic set of conventions with the
// same canonical ordering rules as AssumptionSet.
type ConventionSet struct {
	values []Convention
}

// NewConventionSet canonicalizes, sorts, and deduplicates values.
func NewConventionSet(values ...Convention) ConventionSet {
	type entry struct {
		c Convention
		b []byte
	}
	entries := make([]entry, 0, len(values))
	for _, v := range values {
		if !v.valid() {
			continue
		}
		b, err := canonicalConventionJSON(v)
		if err != nil {
			continue
		}
		entries = append(entries, entry{c: v, b: b})
	}
	sort.Slice(entries, func(i, j int) bool {
		return bytes.Compare(entries[i].b, entries[j].b) < 0
	})
	out := make([]Convention, 0, len(entries))
	var prev []byte
	for i, e := range entries {
		if i > 0 && bytes.Equal(prev, e.b) {
			continue
		}
		out = append(out, e.c)
		prev = e.b
	}
	return ConventionSet{values: out}
}

// Values returns a copy of the canonical element slice.
func (s ConventionSet) Values() []Convention {
	out := make([]Convention, len(s.values))
	copy(out, s.values)
	return out
}

// Merge returns the deterministic union of both sets. Same key with
// different value yields ConventionConflictError (MRC-005).
func (s ConventionSet) Merge(other ConventionSet) (ConventionSet, error) {
	all := make([]Convention, 0, len(s.values)+len(other.values))
	all = append(all, s.values...)
	all = append(all, other.values...)
	seen := make(map[string]string, len(all))
	for _, c := range all {
		if prev, ok := seen[c.key]; ok {
			if prev != c.value {
				return ConventionSet{}, ConventionConflictError{
					Key: c.key, Left: prev, Right: c.value,
				}
			}
			continue
		}
		seen[c.key] = c.value
	}
	return NewConventionSet(all...), nil
}

// Equal reports canonical set equality.
func (s ConventionSet) Equal(other ConventionSet) bool {
	a, errA := s.CanonicalJSON()
	b, errB := other.CanonicalJSON()
	if errA != nil || errB != nil {
		return false
	}
	return bytes.Equal(a, b)
}

// CanonicalJSON encodes the set as its canonical element array.
func (s ConventionSet) CanonicalJSON() ([]byte, error) {
	parts := make([]json.RawMessage, 0, len(s.values))
	for _, c := range s.values {
		body, err := canonicalConventionJSON(c)
		if err != nil {
			return nil, err
		}
		parts = append(parts, json.RawMessage(body))
	}
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, p := range parts {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.Write(p)
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

// Hash returns the SHA-256 of the canonical set JSON.
func (s ConventionSet) Hash() [32]byte {
	b, err := s.CanonicalJSON()
	if err != nil {
		return [32]byte{}
	}
	return HashBytes(b)
}

// ParseConventionSetJSON decodes a canonical convention-set array.
func ParseConventionSetJSON(data []byte) (ConventionSet, error) {
	var raws []json.RawMessage
	if err := json.Unmarshal(data, &raws); err != nil {
		return ConventionSet{}, fmt.Errorf("convention set: %w", err)
	}
	values := make([]Convention, 0, len(raws))
	for _, r := range raws {
		var dto conventionDTO
		if err := strictDecode(r, &dto); err != nil {
			return ConventionSet{}, err
		}
		c, err := NewConvention(dto.Key, dto.Value)
		if err != nil {
			return ConventionSet{}, err
		}
		values = append(values, c)
	}
	set := NewConventionSet(values...)
	canonical, err := set.CanonicalJSON()
	if err != nil {
		return ConventionSet{}, err
	}
	if !bytes.Equal(canonical, data) {
		return ConventionSet{}, fmt.Errorf("non-canonical convention set JSON")
	}
	return set, nil
}

// ---------------------------------------------------------------------------
// Provenance — deterministic record (§13)
// ---------------------------------------------------------------------------

// Provenance is the immutable provenance record (§13.3).
type Provenance struct {
	status         ProvenanceStatus
	source         string
	framework      string
	parentHashes   []string
	assumptionHash [32]byte
	conventionHash [32]byte
	mrcVersion     string
	justification  string
}

// NewProvenance is the single deterministic provenance constructor (§13.0.1).
// It validates status-specific invariants but does not mint an Object.
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
	if !validProvenanceStatus(status) {
		return Provenance{}, fmt.Errorf("invalid provenance status %q", string(status))
	}
	for i, h := range parentHashes {
		if !isLowerHex64(h) {
			return Provenance{}, fmt.Errorf("invalid parent hash at index %d", i)
		}
	}
	if mrcVersion != MRCVersion {
		return Provenance{}, fmt.Errorf("invalid mrc version %q", mrcVersion)
	}
	if justification != "" && status != StatusIdentified && status != StatusHypothesis {
		return Provenance{}, fmt.Errorf("justification permitted only for IDENTIFIED or HYPOTHESIS")
	}
	parents := make([]string, len(parentHashes))
	copy(parents, parentHashes)
	return Provenance{
		status:         status,
		source:         source,
		framework:      framework,
		parentHashes:   parents,
		assumptionHash: assumptionHash,
		conventionHash: conventionHash,
		mrcVersion:     mrcVersion,
		justification:  justification,
	}, nil
}

// NewProvenanceMust is like NewProvenance but panics on error.
// For use in tests and fixed constructor code where inputs are known valid.
func NewProvenanceMust(
	status ProvenanceStatus,
	source string,
	framework string,
	parentHashes []string,
	assumptionHash [32]byte,
	conventionHash [32]byte,
	mrcVersion string,
	justification string,
) Provenance {
	p, err := NewProvenance(status, source, framework, parentHashes,
		assumptionHash, conventionHash, mrcVersion, justification)
	if err != nil {
		panic(err)
	}
	return p
}

// Status returns the provenance status.
func (p Provenance) Status() ProvenanceStatus { return p.status }

// Source returns the source attribution string.
func (p Provenance) Source() string { return p.source }

// Framework returns the framework identifier.
func (p Provenance) Framework() string { return p.framework }

// ParentHashes returns a copy of the parent object hashes in input order.
func (p Provenance) ParentHashes() []string {
	out := make([]string, len(p.parentHashes))
	copy(out, p.parentHashes)
	return out
}

// AssumptionHash returns the assumption-set hash recorded at mint time.
func (p Provenance) AssumptionHash() [32]byte { return p.assumptionHash }

// ConventionHash returns the convention-set hash recorded at mint time.
func (p Provenance) ConventionHash() [32]byte { return p.conventionHash }

// MRCVersion returns the MRC version string ("mrc-v0.4" for MVP).
func (p Provenance) MRCVersion() string { return p.mrcVersion }

// Justification returns the identification justification ("" unless
// IDENTIFIED or a recorded HYPOTHESIS identification attempt).
func (p Provenance) Justification() string { return p.justification }

func (p Provenance) valid() bool {
	if !validProvenanceStatus(p.status) || p.mrcVersion != MRCVersion {
		return false
	}
	if p.justification != "" && p.status != StatusIdentified && p.status != StatusHypothesis {
		return false
	}
	for _, h := range p.parentHashes {
		if !isLowerHex64(h) {
			return false
		}
	}
	return true
}

type provenanceDTO struct {
	Status         string   `json:"status"`
	Source         string   `json:"source"`
	Framework      string   `json:"framework"`
	ParentHashes   []string `json:"parent_hashes"`
	AssumptionHash string   `json:"assumption_hash"`
	ConventionHash string   `json:"convention_hash"`
	MRCVersion     string   `json:"mrc_version"`
	Justification  string   `json:"justification"`
}

func canonicalProvenanceJSON(p Provenance) ([]byte, error) {
	if !p.valid() {
		return nil, fmt.Errorf("invalid provenance")
	}
	parents := p.parentHashes
	if parents == nil {
		parents = []string{}
	}
	return marshalCanonical(provenanceDTO{
		Status:         string(p.status),
		Source:         p.source,
		Framework:      p.framework,
		ParentHashes:   parents,
		AssumptionHash: Hex(p.assumptionHash),
		ConventionHash: Hex(p.conventionHash),
		MRCVersion:     p.mrcVersion,
		Justification:  p.justification,
	})
}

func decodeProvenance(data []byte) (Provenance, error) {
	var dto provenanceDTO
	if err := strictDecode(data, &dto); err != nil {
		return Provenance{}, err
	}
	var assumptionHash, conventionHash [32]byte
	if len(dto.AssumptionHash) == 64 {
		copy(assumptionHash[:], mustDecodeHex(dto.AssumptionHash))
	}
	if len(dto.ConventionHash) == 64 {
		copy(conventionHash[:], mustDecodeHex(dto.ConventionHash))
	}
	return NewProvenance(
		ProvenanceStatus(dto.Status),
		dto.Source,
		dto.Framework,
		dto.ParentHashes,
		assumptionHash,
		conventionHash,
		dto.MRCVersion,
		dto.Justification,
	)
}

func mustDecodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil
	}
	return b
}

// EqualProvenance reports canonical equality of two provenance records.
func EqualProvenance(a, b Provenance) bool {
	x, errA := canonicalProvenanceJSON(a)
	y, errB := canonicalProvenanceJSON(b)
	if errA != nil || errB != nil {
		return false
	}
	return bytes.Equal(x, y)
}

// ---------------------------------------------------------------------------
// Object — immutable physical-object carrier (§5)
// ---------------------------------------------------------------------------

// Object is the one immutable generic carrier used by ops. All authoritative
// fields are unexported (REQ-005-06); the zero value is invalid
// (REQ-005-07); no mutator exists (REQ-005-08).
type Object struct {
	valid        bool
	name         string
	kind         Kind
	dimension    Dimension
	expr         Expr
	assumptions  AssumptionSet
	conventions  ConventionSet
	provenance   Provenance
	corpusStatus CorpusStatus
}

// Valid reports mint-time validity.
func (o Object) Valid() bool { return o.valid }

// Name returns the object name ("" for operation results).
func (o Object) Name() string { return o.name }

// Kind returns the physical kind.
func (o Object) Kind() Kind { return o.kind }

// Dimension returns the dimension value.
func (o Object) Dimension() Dimension { return o.dimension }

// Expr returns the expression handle (immutable).
func (o Object) Expr() Expr { return o.expr }

// Assumptions returns the immutable assumption set.
func (o Object) Assumptions() AssumptionSet { return o.assumptions }

// Conventions returns the immutable convention set.
func (o Object) Conventions() ConventionSet { return o.conventions }

// Provenance returns the provenance record.
func (o Object) Provenance() Provenance { return o.provenance }

// CorpusStatus returns the human-curated corpus status.
func (o Object) CorpusStatus() CorpusStatus { return o.corpusStatus }

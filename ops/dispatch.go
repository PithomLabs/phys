package ops

import (
	"encoding/json"
	"io"
	"math/big"
	"strings"

	"github.com/PithomLabs/phys/core"
)

// ---------------------------------------------------------------------------
// Pure-operation dispatch (specs_v2_3.md §15.13, §15.13.1)
// ---------------------------------------------------------------------------

// OperationID is the closed set of dispatch identifiers. identify is
// deliberately excluded: identification is session-owned (MRC-006).
type OperationID string

// The exact OperationID values — the dispatch switch is closed over exactly
// these twelve; no registry or runtime registration exists.
const (
	OpAdd           OperationID = "add"
	OpSubtract      OperationID = "subtract"
	OpMultiply      OperationID = "multiply"
	OpDivide        OperationID = "divide"
	OpPow           OperationID = "pow"
	OpSimplify      OperationID = "simplify"
	OpSubstitute    OperationID = "substitute"
	OpDifferentiate OperationID = "differentiate"
	OpLimit         OperationID = "limit"
	OpCompare       OperationID = "compare"
	OpSolve         OperationID = "solve"
	OpSelectBranch  OperationID = "select_branch"
)

// OperationParams has the exact field order and canonical JSON shape of
// §15.13. Unused fields encode their empty/default values.
type OperationParams struct {
	Kind          string                `json:"kind"`
	Exponent      string                `json:"exponent"`
	Operator      core.RelationOperator `json:"operator"`
	Justification string                `json:"justification"`
}

// CanonicalJSON encodes the params in canonical form (struct field order).
func (p OperationParams) CanonicalJSON() ([]byte, error) {
	return json.Marshal(p)
}

// ParseOperationParams strictly decodes canonical params bytes
// (DisallowUnknownFields) and validates the field binding generically.
func ParseOperationParams(data []byte) (OperationParams, error) {
	var p OperationParams
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&p); err != nil {
		return OperationParams{}, err
	}
	if err := decodeExact(data); err != nil {
		return OperationParams{}, err
	}
	if err := p.validateShape(); err != nil {
		return OperationParams{}, err
	}
	return p, nil
}

// decodeExact rejects trailing content after the single JSON document.
func decodeExact(data []byte) error {
	dec := json.NewDecoder(strings.NewReader(string(data)))
	var raw json.RawMessage
	if err := dec.Decode(&raw); err != nil {
		return err
	}
	var extra json.RawMessage
	if err := dec.Decode(&extra); err != io.EOF {
		return core.UnsupportedOperationError{
			Operation: "params", Reason: "trailing content after params document",
		}
	}
	return nil
}

// validateShape enforces the permitted Kind values and the per-kind field
// binding: pow → Exponent, compare → Operator, identify → Justification,
// empty → all other fields.
// Freeze invariant (Plan 8 H6): every OperationParams kind explicitly defines
// its meaningful fields; every unused field must hold its canonical
// empty/default value. This keeps canonical parameter encodings unique per
// operation semantics (spec §15.13; the v2.3 "operator":"eq" pow example is
// superseded by this operational pin, spec text unchanged).
func (p OperationParams) validateShape() error {
	switch p.Kind {
	case "empty":
		if p.Exponent != "" || p.Operator != "" || p.Justification != "" {
			return core.UnsupportedOperationError{
				Operation: "params", Reason: "empty params require all empty fields",
			}
		}
	case "pow":
		if _, err := parseExactRational(p.Exponent); err != nil {
			return core.UnsupportedOperationError{
				Operation: "params", Reason: "pow params require an exact rational exponent",
			}
		}
		if p.Operator != "" || p.Justification != "" {
			return core.UnsupportedOperationError{
				Operation: "params", Reason: "pow params require empty operator and justification",
			}
		}
	case "compare":
		switch p.Operator {
		case core.RelationEq, core.RelationNeq, core.RelationLt, core.RelationLte,
			core.RelationGt, core.RelationGte:
		default:
			return core.UnsupportedOperationError{
				Operation: "params", Reason: "compare params require a relation operator",
			}
		}
		if p.Exponent != "" || p.Justification != "" {
			return core.UnsupportedOperationError{
				Operation: "params", Reason: "compare params require empty exponent and justification",
			}
		}
	case "identify":
		if strings.TrimSpace(p.Justification) == "" {
			return core.UnsupportedOperationError{
				Operation: "params", Reason: "identify params require a justification",
			}
		}
		if p.Exponent != "" || p.Operator != "" {
			return core.UnsupportedOperationError{
				Operation: "params", Reason: "identify params require empty exponent and operator",
			}
		}
	default:
		return core.UnsupportedOperationError{
			Operation: "params", Reason: "unknown params kind",
		}
	}
	return nil
}

// ValidateOperationParams checks params against the operation ID exactly as
// §15.13 requires (Session.Step calls this before dispatch). identify is
// never accepted: neither as an OperationID nor as a params kind.
func ValidateOperationParams(id OperationID, params OperationParams) error {
	if err := params.validateShape(); err != nil {
		return err
	}
	switch id {
	case OpPow:
		if params.Kind != "pow" {
			return core.UnsupportedOperationError{
				Operation: string(id), Reason: "pow requires params kind pow",
			}
		}
	case OpCompare:
		if params.Kind != "compare" {
			return core.UnsupportedOperationError{
				Operation: string(id), Reason: "compare requires params kind compare",
			}
		}
	case OpAdd, OpSubtract, OpMultiply, OpDivide, OpSimplify, OpSubstitute,
		OpDifferentiate, OpLimit, OpSolve, OpSelectBranch:
		if params.Kind != "empty" {
			return core.UnsupportedOperationError{
				Operation: string(id), Reason: "operation requires params kind empty",
			}
		}
	default:
		return core.UnsupportedOperationError{
			Operation: string(id), Reason: "unknown operation id",
		}
	}
	return nil
}

// Apply is the single fixed dispatch mechanism for session replay. Inputs
// are bound by index exactly per §15.13.1 — never reordered, never packed;
// unused positions are omitted (the input count is exact).
func Apply(id OperationID, inputs []core.Object, params OperationParams) (core.Object, error) {
	if err := ValidateOperationParams(id, params); err != nil {
		return core.Object{}, err
	}
	count := func(n int) error {
		if len(inputs) != n {
			return core.UnsupportedOperationError{
				Operation: string(id),
				Reason:    "input count does not match the positional schema",
			}
		}
		return nil
	}
	switch id {
	case OpAdd:
		if err := count(2); err != nil {
			return core.Object{}, err
		}
		return Add(inputs[0], inputs[1])
	case OpSubtract:
		if err := count(2); err != nil {
			return core.Object{}, err
		}
		return Subtract(inputs[0], inputs[1])
	case OpMultiply:
		if err := count(2); err != nil {
			return core.Object{}, err
		}
		return Multiply(inputs[0], inputs[1])
	case OpDivide:
		if err := count(2); err != nil {
			return core.Object{}, err
		}
		return Divide(inputs[0], inputs[1])
	case OpPow:
		if err := count(1); err != nil {
			return core.Object{}, err
		}
		exponent, err := parseExactRational(params.Exponent)
		if err != nil {
			return core.Object{}, err
		}
		return Pow(inputs[0], exponent)
	case OpSimplify:
		if err := count(1); err != nil {
			return core.Object{}, err
		}
		return Simplify(inputs[0])
	case OpSubstitute:
		if err := count(3); err != nil {
			return core.Object{}, err
		}
		return Substitute(inputs[0], inputs[1], inputs[2])
	case OpDifferentiate:
		if err := count(2); err != nil {
			return core.Object{}, err
		}
		return Differentiate(inputs[0], inputs[1])
	case OpLimit:
		if err := count(3); err != nil {
			return core.Object{}, err
		}
		return Limit(inputs[0], inputs[1], inputs[2])
	case OpCompare:
		if err := count(2); err != nil {
			return core.Object{}, err
		}
		return Compare(inputs[0], inputs[1], params.Operator)
	case OpSolve:
		if err := count(2); err != nil {
			return core.Object{}, err
		}
		return Solve(inputs[0], inputs[1])
	case OpSelectBranch:
		if err := count(2); err != nil {
			return core.Object{}, err
		}
		return SelectBranch(inputs[0], inputs[1])
	default:
		// Unknown ids — including "identify" — are rejected: the switch is
		// closed over exactly the twelve OperationIDs.
		return core.Object{}, core.UnsupportedOperationError{
			Operation: string(id), Reason: "unknown operation id",
		}
	}
}

// parseExactRational accepts only the canonical "num/den" form: parsed,
// reduced, positive-denominator encoding must equal the input byte-for-byte.
func parseExactRational(s string) (*big.Rat, error) {
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		return nil, core.UnsupportedOperationError{
			Operation: "params", Reason: "exponent is not an exact rational",
		}
	}
	if canon := r.Num().String() + "/" + r.Denom().String(); canon != s {
		return nil, core.UnsupportedOperationError{
			Operation: "params", Reason: "exponent is not in canonical num/den form",
		}
	}
	return r, nil
}

package kernel

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ---------------------------------------------------------------------------
// Mint authority (specs_v2_3.md §5.2.1)
// ---------------------------------------------------------------------------

// ObjectSpec is the exact mint input of specs_v2_3.md §5.2.1. It is a
// value-type carrier: Object carries no serialization-only state.
type ObjectSpec struct {
	Name         string
	Kind         Kind
	Dimension    Dimension
	Expr         Expr
	Assumptions  AssumptionSet
	Conventions  ConventionSet
	Provenance   Provenance
	CorpusStatus CorpusStatus
}

// nameRequired reports whether the pinned constructor contract requires a
// non-empty Name for this provenance status (plan §5): fixed domain/corpus
// constructors mint DEFINED and session fixtures mint POSTULATED with names;
// operation results mint with Name == "" (DERIVED, IDENTIFIED, and HYPOTHESIS
// propagation results). The hypothesis package enforces its non-empty `id`
// itself, since HYPOTHESIS operation results are unnamed.
func nameRequired(s ProvenanceStatus) bool {
	return s == StatusDefined || s == StatusPostulated
}

// MintObject is the single production minting entry point for Object values
// (specs_v2_3.md §5.2.1). It performs the seven-point contract: reject an
// invalid expression handle; reject empty names where the constructor
// contract requires one; reject invalid dimensions; validate provenance/hash
// metadata; validate assumption/convention canonicality; normalize all
// immutable metadata into canonical internal form; return a fully valid
// immutable Object or a typed error.
func MintObject(spec ObjectSpec) (Object, error) {
	// 1. expression handle
	if !spec.Expr.Valid() {
		return Object{}, InvalidObjectError{Operation: "mint"}
	}
	// 2. name where the constructor contract requires one
	if nameRequired(spec.Provenance.status) && spec.Name == "" {
		return Object{}, InvalidObjectError{Operation: "mint"}
	}
	// 3. kind and dimension
	if !spec.Kind.valid() {
		return Object{}, InvalidObjectError{Operation: "mint"}
	}
	if !spec.Dimension.Valid() {
		return Object{}, InvalidObjectError{Operation: "mint"}
	}
	// 4. provenance/hash metadata (status, hex parent hashes, mrc version,
	// justification rule — already enforced by NewProvenance and rechecked
	// here so that every minted object passes the authoritative invariant
	// validation regardless of construction path)
	if !spec.Provenance.valid() {
		return Object{}, ProvenanceError{Context: "mint", Reason: "invalid provenance metadata"}
	}
	// 5. assumption/convention canonicality
	if _, err := spec.Assumptions.CanonicalJSON(); err != nil {
		return Object{}, InvalidObjectError{Operation: "mint"}
	}
	if _, err := spec.Conventions.CanonicalJSON(); err != nil {
		return Object{}, InvalidObjectError{Operation: "mint"}
	}
	if !validCorpusStatus(spec.CorpusStatus) {
		return Object{}, InvalidObjectError{Operation: "mint"}
	}
	// 6. normalize immutable metadata into canonical internal form (set
	// normalization is idempotent: canonical input stays byte-identical)
	return Object{
		valid:        true,
		name:         spec.Name,
		kind:         spec.Kind,
		dimension:    spec.Dimension,
		expr:         spec.Expr,
		assumptions:  NewAssumptionSet(spec.Assumptions.values...),
		conventions:  NewConventionSet(spec.Conventions.values...),
		provenance:   spec.Provenance,
		corpusStatus: spec.CorpusStatus,
	}, nil
}

// ---------------------------------------------------------------------------
// Canonical object JSON (§10.4) and round-trip decoder (§10.5)
// ---------------------------------------------------------------------------

// objectDTO fixes the §10.4 canonical field order. schema_version "1" is a
// Plan 10 implementation pin on the serialization DTO only — it is never
// stored in kernel.Object and never computed by the kernel.
type objectDTO struct {
	SchemaVersion string          `json:"schema_version"`
	Valid         bool            `json:"valid"`
	Name          string          `json:"name"`
	Kind          string          `json:"kind"`
	Dimension     json.RawMessage `json:"dimension"`
	Expr          json.RawMessage `json:"expr"`
	Assumptions   json.RawMessage `json:"assumptions"`
	Conventions   json.RawMessage `json:"conventions"`
	Provenance    json.RawMessage `json:"provenance"`
	CorpusStatus  string          `json:"corpus_status"`
}

// schemaVersionPin is the Plan 10 serialization pin for every MVP artifact
// schema (object, ledger, session, candidate, manifest).
const schemaVersionPin = "1"

// CanonicalObjectJSON encodes a valid object in the exact §10.4 field order.
func CanonicalObjectJSON(o Object) ([]byte, error) {
	if !o.valid {
		return nil, fmt.Errorf("cannot encode invalid object")
	}
	dim, err := o.dimension.CanonicalJSON()
	if err != nil {
		return nil, err
	}
	body, err := CanonicalExprJSON(o.expr)
	if err != nil {
		return nil, err
	}
	assumptions, err := o.assumptions.CanonicalJSON()
	if err != nil {
		return nil, err
	}
	conventions, err := o.conventions.CanonicalJSON()
	if err != nil {
		return nil, err
	}
	provenance, err := canonicalProvenanceJSON(o.provenance)
	if err != nil {
		return nil, err
	}
	return marshalCanonical(objectDTO{
		SchemaVersion: schemaVersionPin,
		Valid:         true,
		Name:          o.name,
		Kind:          o.kind.String(),
		Dimension:     dim,
		Expr:          body,
		Assumptions:   assumptions,
		Conventions:   conventions,
		Provenance:    provenance,
		CorpusStatus:  CorpusStatusJSON(o.corpusStatus),
	})
}

// LoadObjectJSON decodes byte-canonical object JSON. It validates the same
// authoritative object invariants as MintObject before any decoded object can
// enter replay, and guarantees Encode(Decode(b)) == b (§10.5). The name is
// deliberately not DecodeObjectJSON/ParseObject (review 🟡-8): object decoding
// stays an internal-kernel concern.
func LoadObjectJSON(data []byte) (Object, error) {
	var dto objectDTO
	if err := strictDecode(data, &dto); err != nil {
		return Object{}, err
	}
	if dto.SchemaVersion != schemaVersionPin {
		return Object{}, fmt.Errorf("unsupported schema_version %q", dto.SchemaVersion)
	}
	if !dto.Valid {
		return Object{}, fmt.Errorf("serialized object must carry valid=true")
	}
	kind, ok := ParseKind(dto.Kind)
	if !ok {
		return Object{}, fmt.Errorf("unknown object kind %q", dto.Kind)
	}
	dimension, err := ParseDimensionJSON(dto.Dimension)
	if err != nil {
		return Object{}, fmt.Errorf("object dimension: %w", err)
	}
	expr, err := ParseExprJSON(dto.Expr)
	if err != nil {
		return Object{}, fmt.Errorf("object expr: %w", err)
	}
	assumptions, err := ParseAssumptionSetJSON(dto.Assumptions)
	if err != nil {
		return Object{}, fmt.Errorf("object assumptions: %w", err)
	}
	conventions, err := ParseConventionSetJSON(dto.Conventions)
	if err != nil {
		return Object{}, fmt.Errorf("object conventions: %w", err)
	}
	provenance, err := decodeProvenance(dto.Provenance)
	if err != nil {
		return Object{}, fmt.Errorf("object provenance: %w", err)
	}
	corpusStatus, ok := ParseCorpusStatusJSON(dto.CorpusStatus)
	if !ok {
		return Object{}, fmt.Errorf("unknown corpus status %q", dto.CorpusStatus)
	}
	obj, err := MintObject(ObjectSpec{
		Name:         dto.Name,
		Kind:         kind,
		Dimension:    dimension,
		Expr:         expr,
		Assumptions:  assumptions,
		Conventions:  conventions,
		Provenance:   provenance,
		CorpusStatus: corpusStatus,
	})
	if err != nil {
		return Object{}, err
	}
	reencoded, err := CanonicalObjectJSON(obj)
	if err != nil {
		return Object{}, err
	}
	if !bytes.Equal(reencoded, data) {
		return Object{}, fmt.Errorf("non-canonical object JSON")
	}
	return obj, nil
}

// EqualObject reports canonical object equality including all authoritative
// physical metadata (§9.1). Objects without a canonical form (invalid) are
// never equal.
func EqualObject(a, b Object) bool {
	x, errA := CanonicalObjectJSON(a)
	y, errB := CanonicalObjectJSON(b)
	if errA != nil || errB != nil {
		return false
	}
	return bytes.Equal(x, y)
}

// HashObject returns the SHA-256 of the object's canonical JSON bytes
// (§9.2). An object without a canonical form hashes to the zero value.
func HashObject(o Object) [32]byte {
	b, err := CanonicalObjectJSON(o)
	if err != nil {
		return [32]byte{}
	}
	return HashBytes(b)
}

// ParseKind resolves a canonical §6 kind name (the manifest/JSON `kind`
// value) to its stable ordinal.
func ParseKind(name string) (Kind, bool) {
	for k := KindMass; k <= KindBranchSet; k++ {
		if k.String() == name {
			return k, true
		}
	}
	return KindMass, false
}

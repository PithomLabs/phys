// Ledger replay substrate (§16.12–§16.20, C1/C2).
//
// Committed steps retain IMMUTABLE CANONICAL BYTES as the authoritative replay
// substrate: InputCanonicals/OutputCanonical are canonical object JSON bytes
// produced once at Commit; ParamsCanonical are canonical params bytes.
// Hashes are computed FROM those retained bytes. Validation decodes the
// retained bytes through the internal kernel decoder (invariant validation +
// Encode(Decode(b))==b) and replays only from decoded values.
//
// Public Step accessors expose decoded value copies.
package session

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
	"github.com/PithomLabs/phys/ops"
)

// StepKind represents the kind of a derivation step (§16.12).
type StepKind string

const (
	StepKindAssertion      StepKind = "assertion"
	StepKindTransformation StepKind = "transformation"
	StepKindIdentification StepKind = "identification"
)

// stepEnvelope is the committed step record. Canonical object material is
// retained as canonical JSON bytes (C1); live core.Object values never cross
// the commit boundary.
type stepEnvelope struct {
	StepID           string            `json:"step_id"`
	Index            int               `json:"index"`
	Label            string            `json:"label"`
	StepKind         string            `json:"step_kind"`
	Operation        string            `json:"operation"`
	InputHashes      []string          `json:"input_hashes"`
	InputCanonicals  []json.RawMessage `json:"input_canonicals"`
	ParamsCanonical  json.RawMessage   `json:"params_canonical"`
	OutputHash       string            `json:"output_hash"`
	OutputCanonical  json.RawMessage   `json:"output_canonical"`
	AssumptionHash   string            `json:"assumption_hash"`
	ConventionHash   string            `json:"convention_hash"`
	ProvenanceStatus string            `json:"provenance_status"`
	MRCVersion       string            `json:"mrc_version"`
	PreviousStepHash string            `json:"previous_step_hash"`
	CurrentStepHash  string            `json:"current_step_hash"`
}

// Step represents a single step in the derivation ledger (§16.12).
// All fields are read-only; construction is internal.
type Step struct {
	stepEnvelope
}

// Accessors for Step fields (§16.12). Canonical material is decoded from the
// retained bytes on each call; callers receive value copies.
func (s Step) StepID() string     { return s.stepEnvelope.StepID }
func (s Step) Index() int         { return s.stepEnvelope.Index }
func (s Step) Label() string      { return s.stepEnvelope.Label }
func (s Step) StepKind() StepKind { return StepKind(s.stepEnvelope.StepKind) }
func (s Step) Operation() string  { return s.stepEnvelope.Operation }
func (s Step) InputHashes() []string {
	return append([]string(nil), s.stepEnvelope.InputHashes...)
}
func (s Step) InputCanonicals() []core.Object {
	objs, err := decodeStepInputs(s.stepEnvelope)
	if err != nil {
		return nil
	}
	return objs
}
func (s Step) InputCanonicalBytes() [][]byte {
	out := make([][]byte, len(s.stepEnvelope.InputCanonicals))
	for i, raw := range s.stepEnvelope.InputCanonicals {
		out[i] = append([]byte(nil), raw...)
	}
	return out
}
func (s Step) ParamsCanonical() []byte {
	return append([]byte(nil), s.stepEnvelope.ParamsCanonical...)
}
func (s Step) OutputHash() string { return s.stepEnvelope.OutputHash }
func (s Step) OutputCanonical() core.Object {
	obj, err := decodeStepOutput(s.stepEnvelope)
	if err != nil {
		return core.Object{}
	}
	return obj
}
func (s Step) OutputCanonicalBytes() []byte {
	return append([]byte(nil), s.stepEnvelope.OutputCanonical...)
}
func (s Step) AssumptionHash() string   { return s.stepEnvelope.AssumptionHash }
func (s Step) ConventionHash() string   { return s.stepEnvelope.ConventionHash }
func (s Step) ProvenanceStatus() string { return s.stepEnvelope.ProvenanceStatus }
func (s Step) MRCVersion() string       { return s.stepEnvelope.MRCVersion }
func (s Step) PreviousStepHash() string { return s.stepEnvelope.PreviousStepHash }
func (s Step) CurrentStepHash() string  { return s.stepEnvelope.CurrentStepHash }

// Ledger represents a committed derivation ledger (§16.12, §16.21).
// It is immutable through the public API.
type Ledger struct {
	steps     []stepEnvelope
	finalHash string
	sealed    bool
	valid     bool
}

// Steps returns a copy of the ledger steps (§16.21).
func (l Ledger) Steps() []Step {
	steps := make([]Step, len(l.steps))
	for i, se := range l.steps {
		steps[i] = Step{se}
	}
	return steps
}

// DerivationHash returns the final step's current hash (§16.18, §26.6).
func (l Ledger) DerivationHash() string {
	return l.finalHash
}

// ---------------------------------------------------------------------------
// Ledger wire format (typed DTO, strict decoding)
// ---------------------------------------------------------------------------

type ledgerDTO struct {
	SchemaVersion string         `json:"schema_version"`
	Steps         []stepEnvelope `json:"steps"`
	FinalHash     string         `json:"final_hash"`
	Sealed        bool           `json:"sealed"`
}

// ParseLedgerJSON loads a ledger for validation (§16.21).
// It performs strict typed decoding and retains canonical bytes as-is; it
// creates no session authority and mints no objects outside validation.
func ParseLedgerJSON(data []byte) (Ledger, error) {
	var dto ledgerDTO
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&dto); err != nil {
		return Ledger{}, fmt.Errorf("parse ledger: %w", err)
	}
	if dto.SchemaVersion != schemaVersionPin {
		return Ledger{}, fmt.Errorf("parse ledger: unsupported schema_version %q", dto.SchemaVersion)
	}
	l := Ledger{steps: dto.Steps, finalHash: dto.FinalHash, sealed: dto.Sealed}
	l.valid = len(l.steps) > 0
	return l, nil
}

// Validate performs structural/integrity validation of the ledger artifact.
// The full ordered replay pipeline (§16.19) belongs to Session.Validate;
// this entry point never performs operation replay, provenance judgment,
// conclusion checks, or containment decisions.
func (l Ledger) Validate() error {
	if l.steps == nil {
		return core.LedgerValidationError{Reason: "no ledger"}
	}
	if len(l.steps) == 0 {
		return core.LedgerValidationError{Reason: "ledger empty"}
	}
	for i, step := range l.steps {
		if err := checkStepIDFormat(step); err != nil {
			return core.LedgerValidationError{Reason: fmt.Sprintf("step %d: %v", i, err)}
		}
		inputs, err := decodeStepInputs(step)
		if err != nil {
			return core.LedgerValidationError{Reason: fmt.Sprintf("step %d: %v", i, err)}
		}
		if err := verifyStepInputHashes(step, inputs); err != nil {
			return core.LedgerValidationError{Reason: fmt.Sprintf("step %d: %v", i, err)}
		}
		output, err := decodeStepOutput(step)
		if err != nil {
			return core.LedgerValidationError{Reason: fmt.Sprintf("step %d: %v", i, err)}
		}
		if err := verifyStepOutputHash(step, output); err != nil {
			return core.LedgerValidationError{Reason: fmt.Sprintf("step %d: %v", i, err)}
		}
		if err := verifyStepMetadataHashes(step, output); err != nil {
			return core.LedgerValidationError{Reason: fmt.Sprintf("step %d: %v", i, err)}
		}
		if _, err := decodeStepParams(step); err != nil {
			return core.LedgerValidationError{Reason: fmt.Sprintf("step %d: %v", i, err)}
		}
		if err := verifyStepStatusAndVersion(step, output); err != nil {
			return core.LedgerValidationError{Reason: fmt.Sprintf("step %d: %v", i, err)}
		}
		if i == 0 {
			if !isLowerHex64(step.PreviousStepHash) {
				return core.LedgerValidationError{Reason: fmt.Sprintf("step %d: genesis hash malformed", i)}
			}
		} else if err := checkChainLink(l.steps[i-1], step); err != nil {
			return core.LedgerValidationError{Reason: fmt.Sprintf("step %d: %v", i, err)}
		}
		if recomputeStepHash(step) != step.CurrentStepHash {
			return core.LedgerValidationError{Reason: fmt.Sprintf("step %d: current-step hash mismatch", i)}
		}
	}
	return nil
}

// CanonicalJSON returns the canonical JSON representation of the ledger.
func (l Ledger) CanonicalJSON() ([]byte, error) {
	if l.steps == nil {
		return nil, fmt.Errorf("no ledger")
	}
	return marshalCanonical(ledgerDTO{
		SchemaVersion: schemaVersionPin,
		Steps:         l.steps,
		FinalHash:     l.finalHash,
		Sealed:        l.sealed,
	})
}

// ---------------------------------------------------------------------------
// Commit-time construction (canonicalize once, retain bytes)
// ---------------------------------------------------------------------------

// newLedgerInternal freezes draft entries into an immutable ledger. Each
// input/output object is canonicalized exactly once; hashes are computed from
// the retained canonical bytes (C1).
func newLedgerInternal(entries []stepEntry, derivationID, genesisHash string) (Ledger, error) {
	l := Ledger{steps: []stepEnvelope{}, valid: true}
	prev := genesisHash
	for i, entry := range entries {
		index := i + 1
		inputBytes := make([]json.RawMessage, len(entry.inputs))
		inputHashes := make([]string, len(entry.inputs))
		for j, obj := range entry.inputs {
			raw, err := kernel.CanonicalObjectJSON(obj)
			if err != nil {
				return Ledger{}, fmt.Errorf("commit step %d input %d: %w", index, j, err)
			}
			inputBytes[j] = json.RawMessage(raw)
			inputHashes[j] = fmt.Sprintf("%x", kernel.HashObject(obj))
		}
		outputRaw, err := kernel.CanonicalObjectJSON(entry.output)
		if err != nil {
			return Ledger{}, fmt.Errorf("commit step %d output: %w", index, err)
		}
		paramsRaw, err := entry.params.CanonicalJSON()
		if err != nil {
			return Ledger{}, fmt.Errorf("commit step %d params: %w", index, err)
		}
		se := stepEnvelope{
			StepID:           fmt.Sprintf("step-%06d", index),
			Index:            index,
			Label:            entry.label,
			StepKind:         string(entry.kind),
			Operation:        entry.operation,
			InputHashes:      inputHashes,
			InputCanonicals:  inputBytes,
			ParamsCanonical:  json.RawMessage(paramsRaw),
			OutputHash:       fmt.Sprintf("%x", kernel.HashObject(entry.output)),
			OutputCanonical:  json.RawMessage(outputRaw),
			AssumptionHash:   fmt.Sprintf("%x", core.HashAssumptionSet(entry.output.Assumptions())),
			ConventionHash:   fmt.Sprintf("%x", core.HashConventionSet(entry.output.Conventions())),
			ProvenanceStatus: string(entry.output.Provenance().Status()),
			MRCVersion:       kernel.MRCVersion,
			PreviousStepHash: prev,
		}
		se.CurrentStepHash = recomputeStepHash(se)
		l.steps = append(l.steps, se)
		prev = se.CurrentStepHash
	}
	l.finalHash = prev
	return l, nil
}

// canonicalStepBody computes the canonical body for step hashing: all Step
// fields except CurrentStepHash (no self-reference; §16.17).
func canonicalStepBody(se stepEnvelope) []byte {
	b := stepEnvelope{
		StepID: se.StepID, Index: se.Index, Label: se.Label,
		StepKind: se.StepKind, Operation: se.Operation,
		InputHashes:      se.InputHashes,
		InputCanonicals:  se.InputCanonicals,
		ParamsCanonical:  se.ParamsCanonical,
		OutputHash:       se.OutputHash,
		OutputCanonical:  se.OutputCanonical,
		AssumptionHash:   se.AssumptionHash,
		ConventionHash:   se.ConventionHash,
		ProvenanceStatus: se.ProvenanceStatus,
		MRCVersion:       se.MRCVersion,
		PreviousStepHash: se.PreviousStepHash,
		CurrentStepHash:  "",
	}
	data, _ := json.Marshal(b)
	return data
}

// recomputeStepHash returns the lowercase hex SHA-256 of the canonical body.
func recomputeStepHash(se stepEnvelope) string {
	return fmt.Sprintf("%064x", sha256.Sum256(canonicalStepBody(se)))
}

// ---------------------------------------------------------------------------
// Leaf validation primitives (C2: each performs exactly one check)
// ---------------------------------------------------------------------------

func checkStepIDFormat(step stepEnvelope) error {
	expected := fmt.Sprintf("step-%06d", step.Index)
	if step.Index < 1 {
		return fmt.Errorf("index %d out of range", step.Index)
	}
	if step.StepID != expected {
		return fmt.Errorf("step ID format error: got %s want %s", step.StepID, expected)
	}
	return nil
}

func decodeStepInputs(step stepEnvelope) ([]core.Object, error) {
	objs := make([]core.Object, len(step.InputCanonicals))
	for j, raw := range step.InputCanonicals {
		obj, err := kernel.LoadObjectJSON(raw)
		if err != nil {
			return nil, fmt.Errorf("input %d decode: %w", j, err)
		}
		objs[j] = obj
	}
	return objs, nil
}

func decodeStepOutput(step stepEnvelope) (core.Object, error) {
	obj, err := kernel.LoadObjectJSON(step.OutputCanonical)
	if err != nil {
		return core.Object{}, fmt.Errorf("output decode: %w", err)
	}
	return obj, nil
}

func verifyStepInputHashes(step stepEnvelope, inputs []core.Object) error {
	if len(step.InputHashes) != len(inputs) {
		return fmt.Errorf("input hash count mismatch: %d hashes vs %d canonicals",
			len(step.InputHashes), len(inputs))
	}
	for j, obj := range inputs {
		if want := fmt.Sprintf("%x", kernel.HashObject(obj)); step.InputHashes[j] != want {
			return fmt.Errorf("input hash mismatch at input %d", j)
		}
	}
	return nil
}

func verifyStepOutputHash(step stepEnvelope, output core.Object) error {
	if want := fmt.Sprintf("%x", kernel.HashObject(output)); step.OutputHash != want {
		return fmt.Errorf("output hash mismatch")
	}
	return nil
}

func verifyStepMetadataHashes(step stepEnvelope, output core.Object) error {
	if want := fmt.Sprintf("%x", core.HashAssumptionSet(output.Assumptions())); step.AssumptionHash != want {
		return fmt.Errorf("assumption hash mismatch")
	}
	if want := fmt.Sprintf("%x", core.HashConventionSet(output.Conventions())); step.ConventionHash != want {
		return fmt.Errorf("convention hash mismatch")
	}
	return nil
}

func decodeStepParams(step stepEnvelope) (ops.OperationParams, error) {
	params, err := ops.ParseOperationParams(step.ParamsCanonical)
	if err != nil {
		return ops.OperationParams{}, fmt.Errorf("params decode: %w", err)
	}
	return params, nil
}

func verifyStepStatusAndVersion(step stepEnvelope, output core.Object) error {
	if step.ProvenanceStatus != string(output.Provenance().Status()) {
		return fmt.Errorf("provenance status mismatch: step %q vs output %q",
			step.ProvenanceStatus, string(output.Provenance().Status()))
	}
	if step.MRCVersion != kernel.MRCVersion {
		return fmt.Errorf("mrc version mismatch: got %q", step.MRCVersion)
	}
	return nil
}

func checkChainLink(prev, step stepEnvelope) error {
	if step.PreviousStepHash != prev.CurrentStepHash {
		return fmt.Errorf("chain linkage mismatch")
	}
	return nil
}

func isLowerHex64(s string) bool {
	if len(s) != 64 {
		return false
	}
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}

func mustCanonicalParams(p ops.OperationParams) []byte {
	b, err := p.CanonicalJSON()
	if err != nil {
		panic(err)
	}
	return b
}

// marshalCanonical marshals with HTML escaping disabled for stable bytes.
func marshalCanonical(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	out := buf.Bytes()
	return bytes.TrimSuffix(out, []byte("\n")), nil
}

var _ = strings.TrimSpace

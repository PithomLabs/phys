// Package session: session authority, ledger/replay, and identification
// firewall (specs_v2_3.md §16, §26; Plan 10 §12). Imports only core,
// ops, and internal/kernel — never mechanics, relativity, or hypothesis.
package session

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
	"github.com/PithomLabs/phys/ops"
)

// schemaVersionPin is the Plan 10 serialization pin for every MVP artifact
// schema (object, ledger, session, candidate, manifest). DTO state only.
const schemaVersionPin = "1"

// ---------------------------------------------------------------------------
// SessionState
// ---------------------------------------------------------------------------
type SessionState string

const (
	StateNew       SessionState = "new"
	StateDrafting  SessionState = "drafting"
	StateCommitted SessionState = "committed"
	StateConcluded SessionState = "concluded"
	StateSealed    SessionState = "sealed"
)

// ---------------------------------------------------------------------------
// DraftMetadata (§16.2.1) — exact v2.3 typed model
// ---------------------------------------------------------------------------
type DraftMetadata struct {
	DerivationID            string                   `json:"derivation_id"`
	Label                   string                   `json:"label"`
	Hypothesis              core.Object              `json:"hypothesis"`
	Premises                []core.Object            `json:"premises"`
	Assumptions             core.AssumptionSet       `json:"assumptions"`
	FrameworkDependencies   []FrameworkDependency    `json:"framework_dependencies"`
	Predictions             []Prediction             `json:"predictions"`
	FalsificationConditions []FalsificationCondition `json:"falsification_conditions"`
	RecoveryClaims          []RecoveryClaim          `json:"recovery_claims"`
	AnomalyReferences       []AnomalyReference       `json:"anomaly_references"`
	ReviewHistory           []core.Review            `json:"review_history"`
}

// ---------------------------------------------------------------------------
// Session
// ---------------------------------------------------------------------------
type Session struct {
	state               SessionState
	derivationID        string
	label               string
	mrcVersion          string
	genesisHash         string
	draftEntries        []stepEntry
	ledger              Ledger
	metadata            DraftMetadata
	conclusionSet       bool
	conclusionHash      string
	conclusionCanonical json.RawMessage
}

// stepEntry — unexported draft material (§16.8). Live objects are acceptable
// in the draft buffer; Commit canonicalizes them into retained bytes (C1).
type stepEntry struct {
	label     string
	kind      StepKind
	operation string
	inputs    []core.Object
	params    ops.OperationParams
	output    core.Object
}

// ---------------------------------------------------------------------------
// New (§16.2) — creates a new session in StateNew
// ---------------------------------------------------------------------------
func New() *Session {
	return &Session{
		state:        StateNew,
		mrcVersion:   kernel.MRCVersion,
		genesisHash:  "0000000000000000000000000000000000000000000000000000000000000000",
		draftEntries: []stepEntry{},
	}
}

// ---------------------------------------------------------------------------
// Draft (§16.2) — only entry from New to Drafting
// ---------------------------------------------------------------------------
func (s *Session) Draft(id, label string, meta DraftMetadata) error {
	if s.state != StateNew {
		return core.ProvenanceError{Context: "draft", Reason: fmt.Sprintf("requires new state: %s", s.state)}
	}
	if id == "" {
		return core.ProvenanceError{Context: "draft", Reason: "derivation ID required"}
	}
	s.state = StateDrafting
	s.derivationID = id
	s.label = label
	meta.DerivationID = id
	meta.Label = label
	meta.Premises = append([]core.Object(nil), meta.Premises...)
	meta.FrameworkDependencies = append([]FrameworkDependency(nil), meta.FrameworkDependencies...)
	meta.Predictions = append([]Prediction(nil), meta.Predictions...)
	meta.FalsificationConditions = append([]FalsificationCondition(nil), meta.FalsificationConditions...)
	meta.RecoveryClaims = append([]RecoveryClaim(nil), meta.RecoveryClaims...)
	meta.AnomalyReferences = append([]AnomalyReference(nil), meta.AnomalyReferences...)
	meta.ReviewHistory = append([]core.Review(nil), meta.ReviewHistory...)
	s.metadata = meta
	s.draftEntries = []stepEntry{}
	return nil
}

// ---------------------------------------------------------------------------
// Nine actions (§16.1)
// ---------------------------------------------------------------------------

func (s *Session) Postulate(label string, obj core.Object) error {
	if s.state != StateDrafting {
		return core.ProvenanceError{Context: "postulate", Reason: fmt.Sprintf("requires drafting: %s", s.state)}
	}
	if !obj.Valid() {
		return core.InvalidObjectError{Operation: "postulate"}
	}
	if obj.Provenance().Status() != core.StatusPostulated {
		return core.ProvenanceError{Context: "postulate", Reason: fmt.Sprintf("requires POSTULATED: got %s", obj.Provenance().Status())}
	}
	s.draftEntries = append(s.draftEntries, stepEntry{
		label: label, kind: StepKindAssertion,
		operation: "postulate",
		inputs:    []core.Object{obj},
		params:    ops.OperationParams{Kind: "empty"},
		output:    obj,
	})
	return nil
}

func (s *Session) Declare(label string, obj core.Object) error {
	if s.state != StateDrafting {
		return core.ProvenanceError{Context: "declare", Reason: fmt.Sprintf("requires drafting: %s", s.state)}
	}
	if !obj.Valid() {
		return core.InvalidObjectError{Operation: "declare"}
	}
	s.draftEntries = append(s.draftEntries, stepEntry{
		label: label, kind: StepKindAssertion,
		operation: "declare",
		inputs:    []core.Object{obj},
		params:    ops.OperationParams{Kind: "empty"},
		output:    obj,
	})
	return nil
}

func (s *Session) Define(label string, obj core.Object) error {
	if s.state != StateDrafting {
		return core.ProvenanceError{Context: "define", Reason: fmt.Sprintf("requires drafting: %s", s.state)}
	}
	if !obj.Valid() {
		return core.InvalidObjectError{Operation: "define"}
	}
	if obj.Provenance().Status() != core.StatusDefined {
		return core.ProvenanceError{Context: "define", Reason: fmt.Sprintf("requires DEFINED: got %s", obj.Provenance().Status())}
	}
	s.draftEntries = append(s.draftEntries, stepEntry{
		label: label, kind: StepKindAssertion,
		operation: "define",
		inputs:    []core.Object{obj},
		params:    ops.OperationParams{Kind: "empty"},
		output:    obj,
	})
	return nil
}

func (s *Session) Step(label string, op ops.OperationID, inputs []core.Object, params ops.OperationParams) (core.Object, error) {
	if s.state != StateDrafting {
		return core.Object{}, core.ProvenanceError{Context: "step", Reason: fmt.Sprintf("requires drafting: %s", s.state)}
	}
	if op == "identify" {
		return core.Object{}, core.ProvenanceError{Context: "step", Reason: "identify rejected: use Session.Identify"}
	}
	if err := ops.ValidateOperationParams(op, params); err != nil {
		return core.Object{}, fmt.Errorf("invalid params for %s: %w", op, err)
	}
	cp := append([]core.Object(nil), inputs...)
	result, err := ops.Apply(op, cp, params)
	if err != nil {
		return core.Object{}, fmt.Errorf("step %s failed: %w", op, err)
	}
	s.draftEntries = append(s.draftEntries, stepEntry{
		label: label, kind: StepKindTransformation,
		operation: string(op), inputs: cp, params: params, output: result,
	})
	return result, nil
}

// Identify — the only identification API (MRC-006, §16.7). Delegates to the
// session-owned pure reconstruction helper shared with replay.
func (s *Session) Identify(a, b core.Object, justification string) (core.Object, error) {
	if s.state != StateDrafting {
		return core.Object{}, core.ProvenanceError{Context: "identify", Reason: fmt.Sprintf("requires drafting: %s", s.state)}
	}
	result, err := reconstructIdentification(a, b, justification)
	if err != nil {
		return core.Object{}, err
	}
	s.draftEntries = append(s.draftEntries, stepEntry{
		label: "identify", kind: StepKindIdentification,
		operation: "identify", inputs: []core.Object{a, b},
		params: ops.OperationParams{Kind: "identify", Justification: strings.TrimSpace(justification)},
		output: result,
	})
	return result, nil
}

// reconstructIdentification is the private pure session-owned identification
// path (MRC-006). It performs the full Identify semantics without touching
// session state; both Identify and replay use it.
func reconstructIdentification(a, b core.Object, justification string) (core.Object, error) {
	trimmed := strings.TrimSpace(justification)
	if trimmed == "" {
		return core.Object{}, core.IdentifyError{Reason: "justification required"}
	}
	if !a.Valid() || !b.Valid() {
		return core.Object{}, core.InvalidObjectError{Operation: "identify"}
	}
	if !a.Dimension().Equal(b.Dimension()) {
		return core.Object{}, core.DimensionMismatchError{
			Operation: "identify",
			Left:      dimensionStringOf(a),
			Right:     dimensionStringOf(b),
		}
	}
	// Compare-kind rules (§6.2): every combination allowed except two
	// distinct named kinds (e.g. Mass vs RestMass at equal dimension M).
	ka, kb := a.Kind(), b.Kind()
	if ka != core.KindExpression && kb != core.KindExpression && ka != kb {
		return core.Object{}, core.CategoryMismatchError{
			Operation: "identify", Left: ka.String(), Right: kb.String(),
		}
	}
	mergedAssumptions, err := a.Assumptions().Merge(b.Assumptions())
	if err != nil {
		return core.Object{}, err
	}
	mergedConventions, err := a.Conventions().Merge(b.Conventions())
	if err != nil {
		return core.Object{}, err
	}
	status := core.StatusIdentified
	if a.Provenance().Status() == core.StatusHypothesis || b.Provenance().Status() == core.StatusHypothesis {
		status = core.StatusHypothesis
	}
	relation := core.NewRelation(core.RelationEq, a.Expr(), b.Expr())
	if !relation.Valid() {
		return core.Object{}, core.IdentifyError{Reason: "operands cannot form a relation"}
	}
	prov, err := kernel.NewProvenance(
		status,
		inheritProvenanceString(a.Provenance().Source(), b.Provenance().Source()),
		inheritProvenanceString(a.Provenance().Framework(), b.Provenance().Framework()),
		[]string{
			fmt.Sprintf("%x", kernel.HashObject(a)),
			fmt.Sprintf("%x", kernel.HashObject(b)),
		},
		core.HashAssumptionSet(mergedAssumptions),
		core.HashConventionSet(mergedConventions),
		kernel.MRCVersion,
		trimmed,
	)
	if err != nil {
		return core.Object{}, fmt.Errorf("identify provenance: %w", err)
	}
	result, err := kernel.MintObject(kernel.ObjectSpec{
		Name: "", Kind: core.KindRelation, Dimension: a.Dimension(),
		Expr:         relation,
		Assumptions:  mergedAssumptions,
		Conventions:  mergedConventions,
		Provenance:   prov,
		CorpusStatus: kernel.CorpusNone,
	})
	if err != nil {
		return core.Object{}, fmt.Errorf("identify mint: %w", err)
	}
	return result, nil
}

func inheritProvenanceString(a, b string) string {
	if a != "" && a == b {
		return a
	}
	return ""
}

func dimensionStringOf(o core.Object) string {
	b, err := o.Dimension().CanonicalJSON()
	if err != nil {
		return ""
	}
	return string(b)
}

func (s *Session) Conclude(obj core.Object) error {
	if s.state != StateCommitted {
		return core.ProvenanceError{Context: "conclude", Reason: fmt.Sprintf("requires committed: %s", s.state)}
	}
	if !obj.Valid() {
		return core.InvalidObjectError{Operation: "conclude"}
	}
	if len(s.ledger.steps) == 0 {
		return core.LedgerValidationError{Reason: "conclude requires a non-empty ledger"}
	}
	last := s.ledger.steps[len(s.ledger.steps)-1]
	if want := fmt.Sprintf("%x", kernel.HashObject(obj)); want != last.OutputHash {
		return core.LedgerValidationError{Reason: fmt.Sprintf("conclude object hash mismatch: got %s, expected %s", want, last.OutputHash)}
	}
	raw, err := kernel.CanonicalObjectJSON(obj)
	if err != nil {
		return core.LedgerValidationError{Reason: fmt.Sprintf("conclude canonicalization: %v", err)}
	}
	s.conclusionSet = true
	s.conclusionHash = fmt.Sprintf("%x", kernel.HashObject(obj))
	s.conclusionCanonical = json.RawMessage(raw)
	s.state = StateConcluded
	return nil
}

func (s *Session) Commit() error {
	if s.state != StateDrafting {
		return core.ProvenanceError{Context: "commit", Reason: fmt.Sprintf("requires drafting: %s", s.state)}
	}
	if len(s.draftEntries) == 0 {
		return core.LedgerValidationError{Reason: "commit requires non-empty draft"}
	}
	ledger, err := newLedgerInternal(s.draftEntries, s.derivationID, s.genesisHash)
	if err != nil {
		return err
	}
	s.ledger = ledger
	s.state = StateCommitted
	s.draftEntries = []stepEntry{}
	return nil
}

func (s *Session) Seal() (ResearchCandidate, error) {
	if s.state != StateConcluded {
		return ResearchCandidate{}, core.ProvenanceError{Context: "seal", Reason: fmt.Sprintf("requires concluded: %s", s.state)}
	}
	if !s.metadata.Hypothesis.Valid() || s.metadata.Hypothesis.Provenance().Status() != core.StatusHypothesis {
		return ResearchCandidate{}, core.ProvenanceError{Context: "seal", Reason: "requires HYPOTHESIS hypothesis"}
	}
	if err := s.Validate(); err != nil {
		return ResearchCandidate{}, fmt.Errorf("seal validation failed: %w", err)
	}
	candidate, err := buildTrustedCandidate(s.derivationID, s.ledger, s.metadata)
	if err != nil {
		return ResearchCandidate{}, fmt.Errorf("seal construction failed: %w", err)
	}
	if err := candidate.Validate(); err != nil {
		return ResearchCandidate{}, fmt.Errorf("seal candidate validation failed: %w", err)
	}
	s.state = StateSealed
	s.ledger.sealed = true
	return candidate, nil
}

// Validate — the exact 17-stage ordered replay pipeline of §16.19.
// Session.Validate is authoritative for session semantics; it orchestrates
// atomic leaf primitives (C2) in the mandated order. Read-only.
func (s *Session) Validate() error {
	// Stage 1: session structure/state.
	if s.ledger.steps == nil {
		return core.LedgerValidationError{Reason: "no ledger"}
	}
	if !validStateForValidate(s.state) {
		return core.LedgerValidationError{Reason: fmt.Sprintf("invalid state for validate: %s", s.state)}
	}
	if len(s.ledger.steps) == 0 {
		return core.LedgerValidationError{Reason: "ledger empty"}
	}
	// Stage 2: genesis.
	if s.ledger.steps[0].PreviousStepHash != s.genesisHash {
		return core.LedgerValidationError{Reason: "genesis mismatch"}
	}
	for i, step := range s.ledger.steps {
		tag := fmt.Sprintf("step %d", i)
		// Stage 3: StepID/index.
		if err := checkStepIDFormat(step); err != nil {
			return core.LedgerValidationError{Reason: tag + ": " + err.Error()}
		}
		// Stage 4: reconstruct/decode retained input canonicals.
		inputs, err := decodeStepInputs(step)
		if err != nil {
			return core.LedgerValidationError{Reason: tag + ": " + err.Error()}
		}
		// Stage 5: verify InputHash against canonical input.
		if err := verifyStepInputHashes(step, inputs); err != nil {
			return core.LedgerValidationError{Reason: tag + ": " + err.Error()}
		}
		// Stage 6: reconstruct/decode retained output canonical.
		output, err := decodeStepOutput(step)
		if err != nil {
			return core.LedgerValidationError{Reason: tag + ": " + err.Error()}
		}
		// Stage 7: verify OutputHash.
		if err := verifyStepOutputHash(step, output); err != nil {
			return core.LedgerValidationError{Reason: tag + ": " + err.Error()}
		}
		// Stage 8: verify AssumptionHash and ConventionHash.
		if err := verifyStepMetadataHashes(step, output); err != nil {
			return core.LedgerValidationError{Reason: tag + ": " + err.Error()}
		}
		// Stage 9: recompute current step hash/chain.
		if recomputeStepHash(step) != step.CurrentStepHash {
			return core.LedgerValidationError{Reason: tag + ": current-step hash mismatch"}
		}
		if i > 0 {
			if err := checkChainLink(s.ledger.steps[i-1], step); err != nil {
				return core.LedgerValidationError{Reason: tag + ": " + err.Error()}
			}
		}
		// Stage 10: validate OperationParams against OperationID.
		params, err := decodeStepParams(step)
		if err != nil {
			return core.LedgerValidationError{Reason: tag + ": " + err.Error()}
		}
		if err := checkParamsBinding(step, params); err != nil {
			return core.LedgerValidationError{Reason: tag + ": " + err.Error()}
		}
		// Stage 11: replay assertion actions.
		// Stage 12: replay transformations via ops.Apply.
		// Stage 13: replay identifications via the session-owned path.
		replayed, err := replayStep(step, inputs, params)
		if err != nil {
			return core.LedgerValidationError{Reason: tag + ": replay failed: " + err.Error()}
		}
		// Stage 14: compare replayed result against retained output/hash.
		replayedRaw, err := kernel.CanonicalObjectJSON(replayed)
		if err != nil {
			return core.LedgerValidationError{Reason: tag + ": replayed canonicalization: " + err.Error()}
		}
		if string(replayedRaw) != string(step.OutputCanonical) {
			return core.LedgerValidationError{Reason: tag + ": replayed output differs from retained output"}
		}
		if want := fmt.Sprintf("%x", kernel.HashObject(replayed)); want != step.OutputHash {
			return core.LedgerValidationError{Reason: tag + ": replayed output hash mismatch"}
		}
		// Stage 15: validate provenance and MRC rules.
		if err := checkStepProvenance(step, inputs, output); err != nil {
			return core.LedgerValidationError{Reason: tag + ": " + err.Error()}
		}
	}
	// Stage 16: validate conclusion hash.
	if err := s.checkConclusion(); err != nil {
		return err
	}
	// Stage 17: validate candidate containment.
	if err := checkContainment(s.ledger.steps); err != nil {
		return err
	}
	return nil
}

func validStateForValidate(st SessionState) bool {
	switch st {
	case StateDrafting, StateCommitted, StateConcluded, StateSealed:
		return true
	default:
		return false
	}
}

// checkParamsBinding enforces the §15.13 per-kind binding between the step's
// operation/kind and its params (stage 10).
func checkParamsBinding(step stepEnvelope, params ops.OperationParams) error {
	switch StepKind(step.StepKind) {
	case StepKindAssertion:
		switch step.Operation {
		case "postulate", "declare", "define":
		default:
			return fmt.Errorf("unknown assertion operation %q", step.Operation)
		}
		if params.Kind != "empty" {
			return fmt.Errorf("assertion step requires params kind empty")
		}
	case StepKindTransformation:
		if err := ops.ValidateOperationParams(ops.OperationID(step.Operation), params); err != nil {
			return fmt.Errorf("operation params: %w", err)
		}
	case StepKindIdentification:
		if step.Operation != "identify" {
			return fmt.Errorf("identification step requires operation identify")
		}
		if params.Kind != "identify" || strings.TrimSpace(params.Justification) == "" {
			return fmt.Errorf("identification step requires params kind identify with justification")
		}
	default:
		return fmt.Errorf("unknown step kind %q", step.StepKind)
	}
	return nil
}

// replayStep re-executes one step from retained canonical inputs and params:
// assertions return the retained object after status checks (stage 11),
// transformations go through ops.Apply at §15.13.1 indices (stage 12),
// identifications go through the session-owned helper (stage 13, never ops).
func replayStep(step stepEnvelope, inputs []core.Object, params ops.OperationParams) (core.Object, error) {
	switch StepKind(step.StepKind) {
	case StepKindAssertion:
		if len(inputs) != 1 {
			return core.Object{}, fmt.Errorf("assertion requires exactly 1 input")
		}
		output, err := decodeStepOutput(step)
		if err != nil {
			return core.Object{}, err
		}
		inputRaw, err := kernel.CanonicalObjectJSON(inputs[0])
		if err != nil {
			return core.Object{}, err
		}
		if string(inputRaw) != string(step.OutputCanonical) {
			return core.Object{}, fmt.Errorf("assertion input differs from retained output")
		}
		switch step.Operation {
		case "postulate":
			if output.Provenance().Status() != core.StatusPostulated {
				return core.Object{}, fmt.Errorf("postulate requires POSTULATED output")
			}
		case "define":
			if output.Provenance().Status() != core.StatusDefined {
				return core.Object{}, fmt.Errorf("define requires DEFINED output")
			}
		case "declare":
		default:
			return core.Object{}, fmt.Errorf("unknown assertion operation %q", step.Operation)
		}
		return output, nil
	case StepKindTransformation:
		return ops.Apply(ops.OperationID(step.Operation), inputs, params)
	case StepKindIdentification:
		if len(inputs) != 2 {
			return core.Object{}, fmt.Errorf("identification requires 2 inputs")
		}
		return reconstructIdentification(inputs[0], inputs[1], params.Justification)
	default:
		return core.Object{}, fmt.Errorf("unknown step kind %q", step.StepKind)
	}
}

// checkStepProvenance enforces the deterministic provenance law per step
// (stage 15): operation results follow the ops contamination law, Identify
// results follow the session law, assertions follow their status classes.
func checkStepProvenance(step stepEnvelope, inputs []core.Object, output core.Object) error {
	if err := verifyStepStatusAndVersion(step, output); err != nil {
		return err
	}
	anyHypothesis := false
	for _, in := range inputs {
		if in.Provenance().Status() == core.StatusHypothesis {
			anyHypothesis = true
		}
	}
	switch StepKind(step.StepKind) {
	case StepKindTransformation:
		want := core.StatusDerived
		if anyHypothesis {
			want = core.StatusHypothesis
		}
		if output.Provenance().Status() != want {
			return fmt.Errorf("transformation provenance: got %s want %s",
				output.Provenance().Status(), want)
		}
	case StepKindIdentification:
		want := core.StatusIdentified
		if anyHypothesis {
			want = core.StatusHypothesis
		}
		if output.Provenance().Status() != want {
			return fmt.Errorf("identification provenance: got %s want %s",
				output.Provenance().Status(), want)
		}
		if output.Provenance().Justification() == "" {
			return fmt.Errorf("identification provenance lacks justification")
		}
	}
	return nil
}

// checkConclusion validates the conclusion hash (stage 16).
func (s *Session) checkConclusion() error {
	switch s.state {
	case StateConcluded, StateSealed:
		if !s.conclusionSet {
			return core.LedgerValidationError{Reason: "concluded session lacks conclusion record"}
		}
		last := s.ledger.steps[len(s.ledger.steps)-1]
		if s.conclusionHash != last.OutputHash {
			return core.LedgerValidationError{Reason: "conclusion hash mismatch"}
		}
		if string(s.conclusionCanonical) == "" {
			return core.LedgerValidationError{Reason: "conclusion canonical missing"}
		}
		concluded, err := kernel.LoadObjectJSON(s.conclusionCanonical)
		if err != nil {
			return core.LedgerValidationError{Reason: fmt.Sprintf("conclusion decode: %v", err)}
		}
		if want := fmt.Sprintf("%x", kernel.HashObject(concluded)); want != last.OutputHash {
			return core.LedgerValidationError{Reason: "conclusion object hash mismatch"}
		}
	}
	return nil
}

// checkContainment enforces MRC-008 over retained derivation material (C4,
// stage 17 and candidate validation). Expected provenance is recomputed by
// replay; any hypothesis-dependent output presented with a trusted status or
// a non-NONE corpus status fails with CandidateContainmentError.
func checkContainment(steps []stepEnvelope) error {
	for i, step := range steps {
		inputs, err := decodeStepInputs(step)
		if err != nil {
			return core.LedgerValidationError{Reason: fmt.Sprintf("step %d: %v", i, err)}
		}
		output, err := decodeStepOutput(step)
		if err != nil {
			return core.LedgerValidationError{Reason: fmt.Sprintf("step %d: %v", i, err)}
		}
		anyHypothesis := false
		for _, in := range inputs {
			if in.Provenance().Status() == core.StatusHypothesis {
				anyHypothesis = true
			}
		}
		if anyHypothesis {
			if output.Provenance().Status() != core.StatusHypothesis {
				return core.CandidateContainmentError{
					Reason: fmt.Sprintf("step %d: hypothesis-dependent output presented as %s",
						i, output.Provenance().Status()),
				}
			}
			if output.CorpusStatus() != kernel.CorpusNone {
				return core.CandidateContainmentError{
					Reason: fmt.Sprintf("step %d: hypothesis-dependent output carries corpus status %v", i, output.CorpusStatus()),
				}
			}
		}
	}
	return nil
}

// CanonicalJSON — read-only serialization with schema_version (DTO only).
func (s *Session) CanonicalJSON() ([]byte, error) {
	if s.state == StateNew {
		return nil, fmt.Errorf("session not drafted")
	}
	var ledgerRaw json.RawMessage
	if s.ledger.valid {
		raw, err := s.ledger.CanonicalJSON()
		if err != nil {
			return nil, err
		}
		ledgerRaw = json.RawMessage(raw)
	}
	return marshalCanonical(sessionDTO{
		SchemaVersion: schemaVersionPin,
		State:         string(s.state),
		DerivationID:  s.derivationID,
		Label:         s.label,
		MRCVersion:    s.mrcVersion,
		Ledger:        ledgerRaw,
	})
}

type sessionDTO struct {
	SchemaVersion string          `json:"schema_version"`
	State         string          `json:"state"`
	DerivationID  string          `json:"derivation_id"`
	Label         string          `json:"label"`
	MRCVersion    string          `json:"mrc_version"`
	Ledger        json.RawMessage `json:"ledger,omitempty"`
}

// Ledger returns the session's ledger (read-only).
func (s *Session) Ledger() Ledger {
	return s.ledger
}

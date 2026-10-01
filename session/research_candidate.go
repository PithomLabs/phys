// ResearchCandidate — trusted candidate artifact (§26, Plan 10 v2.3).
// Only the single private constructor buildTrustedCandidate may create a
// trusted ResearchCandidate; Session.Seal and UnverifiedResearchCandidate
// .Validate share that path. ParseResearchCandidateJSON performs only
// outer-schema validation and exposes no core.Object accessors (REQ-032-10a).
package session

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

// ---------------------------------------------------------------------------
// Candidate support types (exact v2.3: Expr conditions, no Object)
// ---------------------------------------------------------------------------

type Prediction struct {
	ID          string             `json:"id"`
	Observable  string             `json:"observable"`
	Relation    core.Expr          `json:"relation"`
	Assumptions core.AssumptionSet `json:"assumptions"`
}

type FalsificationCondition struct {
	ID                     string    `json:"id"`
	TargetClaim            string    `json:"target_claim"`
	ContradictingCondition core.Expr `json:"contradicting_condition"`
	Regime                 string    `json:"regime"`
}

type RecoveryClaim struct {
	ID            string    `json:"id"`
	Description   string    `json:"description"`
	FromFramework string    `json:"from_framework"`
	Condition     core.Expr `json:"condition"`
}

type AnomalyReference struct {
	ID          string `json:"id"`
	Framework   string `json:"framework"`
	Description string `json:"description"`
}

type FrameworkDependency struct {
	FrameworkID      string   `json:"framework_id"`
	ManifestHash     string   `json:"manifest_hash"`
	AssumptionHashes []string `json:"assumption_hashes"`
}

// ---------------------------------------------------------------------------
// Canonical Expr / AssumptionSet wire helpers
// ---------------------------------------------------------------------------

func marshalExpr(e core.Expr) (json.RawMessage, error) {
	raw, err := kernel.CanonicalExprJSON(e)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(raw), nil
}

func parseExpr(raw json.RawMessage) (core.Expr, error) {
	if len(bytes.TrimSpace(raw)) == 0 {
		return core.Expr{}, fmt.Errorf("empty expr canonical")
	}
	e, err := kernel.ParseExprJSON(raw)
	if err != nil {
		return core.Expr{}, err
	}
	reencoded, err := kernel.CanonicalExprJSON(e)
	if err != nil {
		return core.Expr{}, err
	}
	if !bytes.Equal(bytes.TrimSpace(reencoded), bytes.TrimSpace(raw)) {
		return core.Expr{}, fmt.Errorf("non-canonical expr JSON")
	}
	return e, nil
}

func marshalAssumptionSet(s core.AssumptionSet) (json.RawMessage, error) {
	raw, err := s.CanonicalJSON()
	if err != nil {
		return nil, err
	}
	return json.RawMessage(raw), nil
}

func parseAssumptionSet(raw json.RawMessage) (core.AssumptionSet, error) {
	if len(bytes.TrimSpace(raw)) == 0 {
		return core.AssumptionSet{}, fmt.Errorf("empty assumption set canonical")
	}
	s, err := kernel.ParseAssumptionSetJSON(raw)
	if err != nil {
		return core.AssumptionSet{}, err
	}
	reencoded, err := s.CanonicalJSON()
	if err != nil {
		return core.AssumptionSet{}, err
	}
	if !bytes.Equal(bytes.TrimSpace(reencoded), bytes.TrimSpace(raw)) {
		return core.AssumptionSet{}, fmt.Errorf("non-canonical assumption set JSON")
	}
	return s, nil
}

// predictionWire etc. are the JSON shapes of the support types (Expr and
// AssumptionSet travel as canonical bytes since neither has a JSON decoder
// on the value type itself).
type predictionWire struct {
	ID          string          `json:"id"`
	Observable  string          `json:"observable"`
	Relation    json.RawMessage `json:"relation"`
	Assumptions json.RawMessage `json:"assumptions"`
}

type falsificationWire struct {
	ID                     string          `json:"id"`
	TargetClaim            string          `json:"target_claim"`
	ContradictingCondition json.RawMessage `json:"contradicting_condition"`
	Regime                 string          `json:"regime"`
}

type recoveryWire struct {
	ID            string          `json:"id"`
	Description   string          `json:"description"`
	FromFramework string          `json:"from_framework"`
	Condition     json.RawMessage `json:"condition"`
}

func marshalPrediction(p Prediction) (predictionWire, error) {
	rel, err := marshalExpr(p.Relation)
	if err != nil {
		return predictionWire{}, err
	}
	asm, err := marshalAssumptionSet(p.Assumptions)
	if err != nil {
		return predictionWire{}, err
	}
	return predictionWire{ID: p.ID, Observable: p.Observable, Relation: rel, Assumptions: asm}, nil
}

func parsePrediction(w predictionWire) (Prediction, error) {
	rel, err := parseExpr(w.Relation)
	if err != nil {
		return Prediction{}, fmt.Errorf("prediction relation: %w", err)
	}
	asm, err := parseAssumptionSet(w.Assumptions)
	if err != nil {
		return Prediction{}, fmt.Errorf("prediction assumptions: %w", err)
	}
	return Prediction{ID: w.ID, Observable: w.Observable, Relation: rel, Assumptions: asm}, nil
}

func marshalFalsification(f FalsificationCondition) (falsificationWire, error) {
	cond, err := marshalExpr(f.ContradictingCondition)
	if err != nil {
		return falsificationWire{}, err
	}
	return falsificationWire{ID: f.ID, TargetClaim: f.TargetClaim, ContradictingCondition: cond, Regime: f.Regime}, nil
}

func parseFalsification(w falsificationWire) (FalsificationCondition, error) {
	cond, err := parseExpr(w.ContradictingCondition)
	if err != nil {
		return FalsificationCondition{}, fmt.Errorf("falsification condition: %w", err)
	}
	return FalsificationCondition{ID: w.ID, TargetClaim: w.TargetClaim, ContradictingCondition: cond, Regime: w.Regime}, nil
}

func marshalRecovery(r RecoveryClaim) (recoveryWire, error) {
	cond, err := marshalExpr(r.Condition)
	if err != nil {
		return recoveryWire{}, err
	}
	return recoveryWire{ID: r.ID, Description: r.Description, FromFramework: r.FromFramework, Condition: cond}, nil
}

func parseRecovery(w recoveryWire) (RecoveryClaim, error) {
	cond, err := parseExpr(w.Condition)
	if err != nil {
		return RecoveryClaim{}, fmt.Errorf("recovery condition: %w", err)
	}
	return RecoveryClaim{ID: w.ID, Description: w.Description, FromFramework: w.FromFramework, Condition: cond}, nil
}

// MarshalJSON preserves deterministic canonical encoding for the support types.
func (p Prediction) MarshalJSON() ([]byte, error) {
	w, err := marshalPrediction(p)
	if err != nil {
		return nil, err
	}
	return marshalCanonical(w)
}

func (p *Prediction) UnmarshalJSON(data []byte) error {
	var w predictionWire
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&w); err != nil {
		return err
	}
	parsed, err := parsePrediction(w)
	if err != nil {
		return err
	}
	*p = parsed
	return nil
}

func (f FalsificationCondition) MarshalJSON() ([]byte, error) {
	w, err := marshalFalsification(f)
	if err != nil {
		return nil, err
	}
	return marshalCanonical(w)
}

func (f *FalsificationCondition) UnmarshalJSON(data []byte) error {
	var w falsificationWire
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&w); err != nil {
		return err
	}
	parsed, err := parseFalsification(w)
	if err != nil {
		return err
	}
	*f = parsed
	return nil
}

func (r RecoveryClaim) MarshalJSON() ([]byte, error) {
	w, err := marshalRecovery(r)
	if err != nil {
		return nil, err
	}
	return marshalCanonical(w)
}

func (r *RecoveryClaim) UnmarshalJSON(data []byte) error {
	var w recoveryWire
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&w); err != nil {
		return err
	}
	parsed, err := parseRecovery(w)
	if err != nil {
		return err
	}
	*r = parsed
	return nil
}

// ---------------------------------------------------------------------------
// ResearchCandidate (§26.10) — trusted sealed artifact
// ---------------------------------------------------------------------------

type ResearchCandidate struct {
	id                      string
	hypothesis              core.Object
	premises                []core.Object
	assumptions             core.AssumptionSet
	derivation              Ledger
	frameworkDependencies   []FrameworkDependency
	predictions             []Prediction
	falsificationConditions []FalsificationCondition
	recoveryClaims          []RecoveryClaim
	anomalyReferences       []AnomalyReference
	reviewHistory           []core.Review
	mrcVersion              string
	ledgerHash              string
}

// ---------------------------------------------------------------------------
// UnverifiedResearchCandidate (§26.10) — unverified loader output
// ---------------------------------------------------------------------------

type UnverifiedResearchCandidate struct {
	bytes []byte
}

// ---------------------------------------------------------------------------
// Access methods (§26.10)
// ---------------------------------------------------------------------------

func (c ResearchCandidate) ID() string              { return c.id }
func (c ResearchCandidate) Hypothesis() core.Object { return c.hypothesis }
func (c ResearchCandidate) Premises() []core.Object {
	p := make([]core.Object, len(c.premises))
	copy(p, c.premises)
	return p
}
func (c ResearchCandidate) Assumptions() core.AssumptionSet { return c.assumptions }
func (c ResearchCandidate) Derivation() Ledger              { return c.derivation }
func (c ResearchCandidate) FrameworkDependencies() []FrameworkDependency {
	d := make([]FrameworkDependency, len(c.frameworkDependencies))
	copy(d, c.frameworkDependencies)
	return d
}
func (c ResearchCandidate) Predictions() []Prediction {
	p := make([]Prediction, len(c.predictions))
	copy(p, c.predictions)
	return p
}
func (c ResearchCandidate) FalsificationConditions() []FalsificationCondition {
	f := make([]FalsificationCondition, len(c.falsificationConditions))
	copy(f, c.falsificationConditions)
	return f
}
func (c ResearchCandidate) RecoveryClaims() []RecoveryClaim {
	r := make([]RecoveryClaim, len(c.recoveryClaims))
	copy(r, c.recoveryClaims)
	return r
}
func (c ResearchCandidate) AnomalyReferences() []AnomalyReference {
	a := make([]AnomalyReference, len(c.anomalyReferences))
	copy(a, c.anomalyReferences)
	return a
}
func (c ResearchCandidate) ReviewHistory() []core.Review {
	r := make([]core.Review, len(c.reviewHistory))
	copy(r, c.reviewHistory)
	return r
}
func (c ResearchCandidate) MRCVersion() string { return c.mrcVersion }
func (c ResearchCandidate) LedgerHash() string { return c.ledgerHash }

// ---------------------------------------------------------------------------
// Single trusted construction path (shared by Seal and unverified Validate)
// ---------------------------------------------------------------------------

func buildTrustedCandidate(derivationID string, ledger Ledger, meta DraftMetadata) (ResearchCandidate, error) {
	if derivationID == "" {
		return ResearchCandidate{}, core.ProvenanceError{Context: "candidate", Reason: "derivation ID required"}
	}
	if len(ledger.steps) == 0 {
		return ResearchCandidate{}, core.LedgerValidationError{Reason: "candidate requires a non-empty derivation"}
	}
	if !meta.Hypothesis.Valid() || meta.Hypothesis.Provenance().Status() != core.StatusHypothesis {
		return ResearchCandidate{}, core.CandidateContainmentError{Reason: "candidate hypothesis must carry HYPOTHESIS status"}
	}
	for i, p := range meta.Premises {
		if !p.Valid() {
			return ResearchCandidate{}, core.InvalidObjectError{Operation: fmt.Sprintf("candidate premise %d", i)}
		}
	}
	premises := append([]core.Object(nil), meta.Premises...)
	return ResearchCandidate{
		id:                      derivationID,
		hypothesis:              meta.Hypothesis,
		premises:                premises,
		assumptions:             meta.Assumptions,
		derivation:              ledger,
		frameworkDependencies:   append([]FrameworkDependency(nil), meta.FrameworkDependencies...),
		predictions:             append([]Prediction(nil), meta.Predictions...),
		falsificationConditions: append([]FalsificationCondition(nil), meta.FalsificationConditions...),
		recoveryClaims:          append([]RecoveryClaim(nil), meta.RecoveryClaims...),
		anomalyReferences:       append([]AnomalyReference(nil), meta.AnomalyReferences...),
		reviewHistory:           append([]core.Review(nil), meta.ReviewHistory...),
		mrcVersion:              kernel.MRCVersion,
		ledgerHash:              ledger.finalHash,
	}, nil
}

// ---------------------------------------------------------------------------
// Validation (§26.8) — in exact order
// ---------------------------------------------------------------------------

func (c ResearchCandidate) Validate() error {
	// 1. Schema / canonicalization (structural only).
	if c.id == "" {
		return fmt.Errorf("candidate derivation ID empty")
	}
	if !isLowerHex64(c.ledgerHash) {
		return fmt.Errorf("candidate ledger hash malformed")
	}
	if c.mrcVersion != kernel.MRCVersion {
		return fmt.Errorf("candidate mrc version mismatch: got %q", c.mrcVersion)
	}
	if _, err := kernel.CanonicalObjectJSON(c.hypothesis); err != nil {
		return fmt.Errorf("candidate hypothesis canonicalization: %w", err)
	}
	for i, p := range c.premises {
		if _, err := kernel.CanonicalObjectJSON(p); err != nil {
			return fmt.Errorf("candidate premise %d canonicalization: %w", i, err)
		}
	}
	// 2. Ledger-hash validation: derivation hash == ledger hash.
	if c.derivation.DerivationHash() != c.ledgerHash {
		return fmt.Errorf("ledger hash mismatch: derivation %s vs ledger %s", c.ledgerHash, c.derivation.DerivationHash())
	}
	// 3. Derivation replay: structural ledger validation + replay-derived
	// containment over retained canonical material.
	if err := c.derivation.Validate(); err != nil {
		return fmt.Errorf("derivation replay failed: %w", err)
	}
	if err := checkContainment(c.derivation.steps); err != nil {
		return fmt.Errorf("derivation containment failed: %w", err)
	}
	// 4. Provenance / candidate containment (§26.9): Hypothesis must be
	// HYPOTHESIS; any hypothesis-dependent step output must already have
	// been rejected above if presented as trusted.
	if c.hypothesis.Provenance().Status() != core.StatusHypothesis {
		return core.CandidateContainmentError{Reason: "hypothesis status must be HYPOTHESIS"}
	}
	// 5. Framework-reference integrity.
	for _, fd := range c.frameworkDependencies {
		if fd.FrameworkID == "" {
			return fmt.Errorf("framework dependency: empty framework ID")
		}
		if !isLowerHex64(fd.ManifestHash) {
			return fmt.Errorf("framework dependency manifest hash malformed: %q", fd.ManifestHash)
		}
		for _, ah := range fd.AssumptionHashes {
			if ah == "" {
				return fmt.Errorf("framework dependency: empty assumption hash")
			}
		}
	}
	// 6. Falsifiability / anomaly field validation.
	for _, p := range c.predictions {
		if p.ID == "" || p.Observable == "" {
			return fmt.Errorf("prediction requires non-empty ID and observable")
		}
		if !p.Relation.Valid() {
			return fmt.Errorf("prediction %q has invalid relation", p.ID)
		}
	}
	for _, f := range c.falsificationConditions {
		if f.ID == "" || f.TargetClaim == "" || f.Regime == "" {
			return fmt.Errorf("falsification condition requires non-empty ID, target claim, and regime")
		}
		if !f.ContradictingCondition.Valid() {
			return fmt.Errorf("falsification condition %q has invalid condition", f.ID)
		}
	}
	for _, r := range c.recoveryClaims {
		if r.ID == "" || r.Description == "" || r.FromFramework == "" {
			return fmt.Errorf("recovery claim requires non-empty ID, description, and framework")
		}
		if !r.Condition.Valid() {
			return fmt.Errorf("recovery claim %q has invalid condition", r.ID)
		}
	}
	if len(c.predictions) > 0 && len(c.falsificationConditions) == 0 {
		return fmt.Errorf("candidate falsifiability: predictions require falsification conditions")
	}
	for _, ref := range c.anomalyReferences {
		if ref.ID == "" || ref.Framework == "" || ref.Description == "" {
			return fmt.Errorf("anomaly reference requires non-empty ID, framework, and description")
		}
	}
	// 7. Review challenge category validation (§27.3).
	for _, review := range c.reviewHistory {
		for _, ch := range review.Challenges {
			if !core.ValidReviewCategory(core.ReviewCategory(ch.Category)) {
				return fmt.Errorf("invalid review challenge category: %s", ch.Category)
			}
		}
	}
	// 8. Success: artifact internal consistency only (never physical truth).
	return nil
}

// ---------------------------------------------------------------------------
// Canonical JSON (§26.2 field order) via the typed wire representation
// ---------------------------------------------------------------------------

type candidateWire struct {
	SchemaVersion           string                `json:"schema_version"`
	DerivationID            string                `json:"derivation_id"`
	Hypothesis              json.RawMessage       `json:"hypothesis"`
	Premises                []json.RawMessage     `json:"premises"`
	Assumptions             json.RawMessage       `json:"assumptions"`
	Derivation              json.RawMessage       `json:"derivation"`
	FrameworkDependencies   []FrameworkDependency `json:"framework_dependencies"`
	Predictions             []predictionWire      `json:"predictions"`
	FalsificationConditions []falsificationWire   `json:"falsification_conditions"`
	RecoveryClaims          []recoveryWire        `json:"recovery_claims"`
	AnomalyReferences       []AnomalyReference    `json:"anomaly_references"`
	ReviewHistory           []core.Review         `json:"review_history"`
	MRCVersion              string                `json:"mrc_version"`
	LedgerHash              string                `json:"ledger_hash"`
}

func (c ResearchCandidate) wire() (candidateWire, error) {
	hyp, err := kernel.CanonicalObjectJSON(c.hypothesis)
	if err != nil {
		return candidateWire{}, err
	}
	premises := make([]json.RawMessage, len(c.premises))
	for i, p := range c.premises {
		raw, err := kernel.CanonicalObjectJSON(p)
		if err != nil {
			return candidateWire{}, err
		}
		premises[i] = json.RawMessage(raw)
	}
	if premises == nil {
		premises = []json.RawMessage{}
	}
	asm, err := marshalAssumptionSet(c.assumptions)
	if err != nil {
		return candidateWire{}, err
	}
	derivation, err := c.derivation.CanonicalJSON()
	if err != nil {
		return candidateWire{}, err
	}
	preds := make([]predictionWire, len(c.predictions))
	for i, p := range c.predictions {
		w, err := marshalPrediction(p)
		if err != nil {
			return candidateWire{}, err
		}
		preds[i] = w
	}
	if preds == nil {
		preds = []predictionWire{}
	}
	fals := make([]falsificationWire, len(c.falsificationConditions))
	for i, f := range c.falsificationConditions {
		w, err := marshalFalsification(f)
		if err != nil {
			return candidateWire{}, err
		}
		fals[i] = w
	}
	if fals == nil {
		fals = []falsificationWire{}
	}
	recs := make([]recoveryWire, len(c.recoveryClaims))
	for i, r := range c.recoveryClaims {
		w, err := marshalRecovery(r)
		if err != nil {
			return candidateWire{}, err
		}
		recs[i] = w
	}
	if recs == nil {
		recs = []recoveryWire{}
	}
	fds := c.frameworkDependencies
	if fds == nil {
		fds = []FrameworkDependency{}
	}
	anoms := c.anomalyReferences
	if anoms == nil {
		anoms = []AnomalyReference{}
	}
	revs := c.reviewHistory
	if revs == nil {
		revs = []core.Review{}
	}
	return candidateWire{
		SchemaVersion: schemaVersionPin, DerivationID: c.id,
		Hypothesis: json.RawMessage(hyp), Premises: premises, Assumptions: asm,
		Derivation:            json.RawMessage(derivation),
		FrameworkDependencies: fds, Predictions: preds,
		FalsificationConditions: fals, RecoveryClaims: recs,
		AnomalyReferences: anoms, ReviewHistory: revs,
		MRCVersion: c.mrcVersion, LedgerHash: c.ledgerHash,
	}, nil
}

func (c ResearchCandidate) CanonicalJSON() ([]byte, error) {
	w, err := c.wire()
	if err != nil {
		return nil, err
	}
	return marshalCanonical(w)
}

// ---------------------------------------------------------------------------
// Unverified loader (§26.10): typed intermediate DTO, shared trusted path
// ---------------------------------------------------------------------------

func (uc UnverifiedResearchCandidate) CanonicalJSON() ([]byte, error) {
	return append([]byte(nil), uc.bytes...), nil
}

func (uc UnverifiedResearchCandidate) Validate() (ResearchCandidate, error) {
	var w candidateWire
	dec := json.NewDecoder(bytes.NewReader(uc.bytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&w); err != nil {
		return ResearchCandidate{}, fmt.Errorf("unverified candidate decode: %w", err)
	}
	if w.SchemaVersion != schemaVersionPin {
		return ResearchCandidate{}, fmt.Errorf("unverified candidate: unsupported schema_version %q", w.SchemaVersion)
	}
	hypothesis, err := kernel.LoadObjectJSON(w.Hypothesis)
	if err != nil {
		return ResearchCandidate{}, fmt.Errorf("unverified candidate hypothesis: %w", err)
	}
	premises := make([]core.Object, len(w.Premises))
	for i, raw := range w.Premises {
		obj, err := kernel.LoadObjectJSON(raw)
		if err != nil {
			return ResearchCandidate{}, fmt.Errorf("unverified candidate premise %d: %w", i, err)
		}
		premises[i] = obj
	}
	assumptions, err := parseAssumptionSet(w.Assumptions)
	if err != nil {
		return ResearchCandidate{}, fmt.Errorf("unverified candidate assumptions: %w", err)
	}
	ledger, err := ParseLedgerJSON(w.Derivation)
	if err != nil {
		return ResearchCandidate{}, fmt.Errorf("unverified candidate derivation: %w", err)
	}
	preds := make([]Prediction, len(w.Predictions))
	for i, pw := range w.Predictions {
		p, err := parsePrediction(pw)
		if err != nil {
			return ResearchCandidate{}, fmt.Errorf("unverified candidate prediction %d: %w", i, err)
		}
		preds[i] = p
	}
	fals := make([]FalsificationCondition, len(w.FalsificationConditions))
	for i, fw := range w.FalsificationConditions {
		f, err := parseFalsification(fw)
		if err != nil {
			return ResearchCandidate{}, fmt.Errorf("unverified candidate falsification %d: %w", i, err)
		}
		fals[i] = f
	}
	recs := make([]RecoveryClaim, len(w.RecoveryClaims))
	for i, rw := range w.RecoveryClaims {
		r, err := parseRecovery(rw)
		if err != nil {
			return ResearchCandidate{}, fmt.Errorf("unverified candidate recovery %d: %w", i, err)
		}
		recs[i] = r
	}
	meta := DraftMetadata{
		DerivationID:            w.DerivationID,
		Hypothesis:              hypothesis,
		Premises:                premises,
		Assumptions:             assumptions,
		FrameworkDependencies:   w.FrameworkDependencies,
		Predictions:             preds,
		FalsificationConditions: fals,
		RecoveryClaims:          recs,
		AnomalyReferences:       w.AnomalyReferences,
		ReviewHistory:           w.ReviewHistory,
	}
	candidate, err := buildTrustedCandidate(w.DerivationID, ledger, meta)
	if err != nil {
		return ResearchCandidate{}, fmt.Errorf("unverified candidate construction: %w", err)
	}
	if candidate.ledgerHash != w.LedgerHash {
		return ResearchCandidate{}, fmt.Errorf("unverified candidate ledger hash mismatch")
	}
	if err := candidate.Validate(); err != nil {
		return ResearchCandidate{}, fmt.Errorf("unverified candidate failed validation: %w", err)
	}
	return candidate, nil
}

func ParseResearchCandidateJSON(data []byte) (UnverifiedResearchCandidate, error) {
	var raw map[string]json.RawMessage
	dec := json.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&raw); err != nil {
		return UnverifiedResearchCandidate{}, fmt.Errorf("outer schema validation: %w", err)
	}
	for _, field := range []string{"schema_version", "derivation_id", "hypothesis", "derivation", "ledger_hash"} {
		if _, ok := raw[field]; !ok {
			return UnverifiedResearchCandidate{}, fmt.Errorf("outer schema validation: missing field %q", field)
		}
	}
	cp := append([]byte(nil), data...)
	return UnverifiedResearchCandidate{bytes: cp}, nil
}

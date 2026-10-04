package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
)

// ---------------------------------------------------------------------------
// Manifest types (§17.0)
// ---------------------------------------------------------------------------

type Manifest struct {
	SchemaVersion string           `json:"schema_version"`
	FrameworkID   string           `json:"framework_id"`
	FrameworkName string           `json:"framework_name"`
	CorpusStatus  CorpusStatus     `json:"corpus_status"`
	Assumptions   []Assumption     `json:"assumptions"`
	Domain        []ManifestDomain `json:"domain"`
	Limits        []ManifestLimit  `json:"limits"`
	Anomalies     []ManifestAnomaly `json:"anomalies"`
	Items         []ManifestItem   `json:"items"`
}

type ManifestDomain struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ManifestLimit struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type ManifestAnomaly struct {
	ID                 string   `json:"id"`
	Framework          string   `json:"framework"`
	Description        string   `json:"description"`
	Status             string   `json:"status"`
	RelatedAssumptions []string `json:"related_assumptions"`
	RelatedItems       []string `json:"related_items"`
	ResearchRelevance  string   `json:"research_relevance"`
}

type ManifestReduction struct {
	ID        string `json:"id"`
	Condition string `json:"condition"`
	Result    string `json:"result"`
}

type ManifestItem struct {
	ID                      string             `json:"id"`
	Constructor             string             `json:"constructor"`
	Kind                    string             `json:"kind"`
	Name                    string             `json:"name"`
	Statement               string             `json:"statement"`
	CanonicalExpr           Expr               `json:"canonical_expr"`
	Dimension               Dimension          `json:"dimension"`
	ProvenanceStatus        ProvenanceStatus   `json:"provenance_status"`
	Source                  string             `json:"source"`
	Assumptions             []Assumption       `json:"assumptions"`
	DerivableFrom           []string           `json:"derivable_from"`
	ReducesTo               []ManifestReduction `json:"reduces_to"`
	KnownLimits             []string           `json:"known_limits"`
	Anomalies               []string           `json:"anomalies"`
	FalsificationConditions []string           `json:"falsification_conditions"`
}

// UnmarshalJSON implements custom JSON unmarshaling for ManifestItem
// to decode canonical_expr into the closed Expr node set.
func (item *ManifestItem) UnmarshalJSON(data []byte) error {
	type Alias ManifestItem
	aux := &struct {
		CanonicalExpr json.RawMessage `json:"canonical_expr"`
		*Alias
	}{
		Alias: (*Alias)(item),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if len(aux.CanonicalExpr) > 0 {
		expr, err := ParseExprJSON(aux.CanonicalExpr)
		if err != nil {
			return fmt.Errorf("canonical_expr: %w", err)
		}
		item.CanonicalExpr = expr
	}
	return nil
}

// ---------------------------------------------------------------------------
// Challenge and Review (§27.1, §27.2)
// ---------------------------------------------------------------------------

type Challenge struct {
	StepID      string `json:"step_id"`
	Category    string `json:"category"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

type Review struct {
	DerivationID  string       `json:"derivation_id"`
	Challenges    []Challenge  `json:"challenges"`
	ReviewerNotes string       `json:"reviewer_notes"`
}

// ReviewCategory holds the exact eight MVP categories (§27.3).
type ReviewCategory string

const (
	ReviewCategoryError                ReviewCategory = "CategoryError"
	ReviewCategoryDimensionError       ReviewCategory = "DimensionError"
	ReviewCategoryAssumptionConflict   ReviewCategory = "AssumptionConflict"
	ReviewCategoryConventionConflict   ReviewCategory = "ConventionConflict"
	ReviewCategoryUnsupportedIdentification ReviewCategory = "UnsupportedIdentification"
	ReviewCategoryInvalidReduction     ReviewCategory = "InvalidReduction"
	ReviewCategoryProvenanceProblem    ReviewCategory = "ProvenanceProblem"
	ReviewCategoryCandidateOverreach   ReviewCategory = "CandidateOverreach"
)

func ValidReviewCategory(c ReviewCategory) bool {
	switch c {
	case ReviewCategoryError,
		ReviewCategoryDimensionError,
		ReviewCategoryAssumptionConflict,
		ReviewCategoryConventionConflict,
		ReviewCategoryUnsupportedIdentification,
		ReviewCategoryInvalidReduction,
		ReviewCategoryProvenanceProblem,
		ReviewCategoryCandidateOverreach:
		return true
	default:
		return false
	}
}

// ---------------------------------------------------------------------------
// Parsing and validation (§17.7, §17.8)
// ---------------------------------------------------------------------------

func ParseManifest(data []byte) (Manifest, error) {
	var m Manifest
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&m); err != nil {
		return Manifest{}, ManifestValidationError{Path: "root", Reason: err.Error()}
	}
	if err := validateManifest(m); err != nil {
		return Manifest{}, err
	}
	return m, nil
}

func ValidateManifestBytes(data []byte) (Manifest, error) {
	return ParseManifest(data)
}

func validateManifest(m Manifest) error {
	if m.SchemaVersion == "" {
		return ManifestValidationError{Path: "schema_version", Reason: "missing"}
	}
	if m.FrameworkID == "" {
		return ManifestValidationError{Path: "framework_id", Reason: "missing"}
	}
	if m.FrameworkName == "" {
		return ManifestValidationError{Path: "framework_name", Reason: "missing"}
	}
	if !validCorpusStatus(m.CorpusStatus) {
		return ManifestValidationError{Path: "corpus_status", Reason: "invalid value: " + string(m.CorpusStatus)}
	}
	for i, a := range m.Assumptions {
		if !a.Valid() {
			return ManifestValidationError{Path: "assumptions[" + itoa(i) + "]", Reason: "invalid"}
		}
	}
	seenDomains := make(map[string]bool)
	for i, d := range m.Domain {
		if d.ID == "" {
			return ManifestValidationError{Path: "domain[" + itoa(i) + "].id", Reason: "missing"}
		}
		if seenDomains[d.ID] {
			return ManifestValidationError{Path: "domain[" + itoa(i) + "].id", Reason: "duplicate: " + d.ID}
		}
		seenDomains[d.ID] = true
	}
	seenLimits := make(map[string]bool)
	for i, l := range m.Limits {
		if l.ID == "" {
			return ManifestValidationError{Path: "limits[" + itoa(i) + "].id", Reason: "missing"}
		}
		if seenLimits[l.ID] {
			return ManifestValidationError{Path: "limits[" + itoa(i) + "].id", Reason: "duplicate: " + l.ID}
		}
		seenLimits[l.ID] = true
	}
	seenAnomalies := make(map[string]bool)
	for i, a := range m.Anomalies {
		if a.ID == "" {
			return ManifestValidationError{Path: "anomalies[" + itoa(i) + "].id", Reason: "missing"}
		}
		if seenAnomalies[a.ID] {
			return ManifestValidationError{Path: "anomalies[" + itoa(i) + "].id", Reason: "duplicate: " + a.ID}
		}
		seenAnomalies[a.ID] = true
	}
	seenItems := make(map[string]bool)
	for i, item := range m.Items {
		if item.ID == "" {
			return ManifestValidationError{Path: "items[" + itoa(i) + "].id", Reason: "missing"}
		}
		if seenItems[item.ID] {
			return ManifestValidationError{Path: "items[" + itoa(i) + "].id", Reason: "duplicate: " + item.ID}
		}
		seenItems[item.ID] = true
		if item.Constructor == "" {
			return ManifestValidationError{Path: "items[" + itoa(i) + "].constructor", Reason: "missing"}
		}
		if item.Kind == "" {
			return ManifestValidationError{Path: "items[" + itoa(i) + "].kind", Reason: "missing"}
		}
		if !validKind(item.Kind) {
			return ManifestValidationError{Path: "items[" + itoa(i) + "].kind", Reason: "invalid: " + item.Kind}
		}
		if item.Name == "" {
			return ManifestValidationError{Path: "items[" + itoa(i) + "].name", Reason: "missing"}
		}
		if !item.CanonicalExpr.Valid() {
			return ManifestValidationError{Path: "items[" + itoa(i) + "].canonical_expr", Reason: "invalid"}
		}
		if !item.Dimension.Valid() {
			return ManifestValidationError{Path: "items[" + itoa(i) + "].dimension", Reason: "invalid"}
		}
		if !validProvenanceStatus(item.ProvenanceStatus) {
			return ManifestValidationError{Path: "items[" + itoa(i) + "].provenance_status", Reason: "invalid: " + string(item.ProvenanceStatus)}
		}
		for j, a := range item.Assumptions {
			if !a.Valid() {
				return ManifestValidationError{Path: "items[" + itoa(i) + "].assumptions[" + itoa(j) + "]", Reason: "invalid"}
			}
		}
		for j, r := range item.ReducesTo {
			if r.ID == "" {
				return ManifestValidationError{Path: "items[" + itoa(i) + "].reduces_to[" + itoa(j) + "].id", Reason: "missing"}
			}
		}
	}
	return nil
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

func validCorpusStatus(s CorpusStatus) bool {
	switch s {
	case CorpusNone, CorpusEstablished, CorpusContested, CorpusSuperseded, CorpusFalsified:
		return true
	default:
		return false
	}
}

func validKind(k string) bool {
	switch k {
	case "Mass", "RestMass", "Time", "Energy", "KineticEnergy",
		"SpeedOfLight", "Velocity", "Position", "Acceleration",
		"Force", "Momentum", "ThreeMomentum", "FourMomentum",
		"Spacetime", "MinkowskiMetric", "Expression", "Relation", "BranchSet":
		return true
	default:
		return false
	}
}

func validProvenanceStatus(s ProvenanceStatus) bool {
	switch s {
	case StatusDefined, StatusPostulated, StatusDerived, StatusIdentified, StatusApproximated, StatusHypothesis:
		return true
	default:
		return false
	}
}

// ---------------------------------------------------------------------------
// Canonical JSON and Hash (§17.7, §17.8)
// ---------------------------------------------------------------------------

func CanonicalManifestJSON(m Manifest) ([]byte, error) {
	// Use a deterministic encoder that sorts map keys (though our structs
	// don't use maps). json.Marshal respects struct field order.
	return json.Marshal(m)
}

func (m Manifest) Hash() [32]byte {
	b, err := CanonicalManifestJSON(m)
	if err != nil {
		return [32]byte{}
	}
	h := sha256.Sum256(b)
	return h
}

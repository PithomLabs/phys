// Hypothesis candidate tests (§24 / MRC-008 / L / M / N / Plan 10 v2.3).
package hypothesis_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/hypothesis"
	"github.com/PithomLabs/phys/internal/kernel"
	"github.com/PithomLabs/phys/mechanics"
	"github.com/PithomLabs/phys/ops"
	"github.com/PithomLabs/phys/relativity"
	"github.com/PithomLabs/phys/session"
)

func mustSymbol(t *testing.T, name string) core.Expr {
	t.Helper()
	e, err := core.NewSymbol(name)
	if err != nil {
		t.Fatalf("symbol: %v", err)
	}
	return e
}

// TestHypothesisContamination proves the MRC-008 input-level law: every
// downstream status of a HYPOTHESIS input is HYPOTHESIS, across Add,
// Simplify, and session Identify. Acceptance L (first half).
func TestHypothesisContamination(t *testing.T) {
	h, err := hypothesis.NewCandidateConcept("h_c", core.KindExpression, core.Dimensionless(), mustSymbol(t, "q"))
	if err != nil {
		t.Fatalf("candidate concept: %v", err)
	}
	if h.Provenance().Status() != core.StatusHypothesis {
		t.Fatalf("concept status = %v, want HYPOTHESIS", h.Provenance().Status())
	}
	if h.CorpusStatus() != core.CorpusNone {
		t.Fatalf("concept corpus status = %v, want NONE", h.CorpusStatus())
	}
	plain, err := hypothesis.NewCandidateConcept("plain", core.KindExpression, core.Dimensionless(), mustSymbol(t, "r"))
	if err != nil {
		t.Fatalf("second concept: %v", err)
	}
	added, err := ops.Add(h, plain)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if added.Provenance().Status() != core.StatusHypothesis {
		t.Errorf("Add(h, x) status = %v, want HYPOTHESIS", added.Provenance().Status())
	}
	simp, err := ops.Simplify(h)
	if err != nil {
		t.Fatalf("Simplify: %v", err)
	}
	if simp.Provenance().Status() != core.StatusHypothesis {
		t.Errorf("Simplify(h) status = %v, want HYPOTHESIS", simp.Provenance().Status())
	}
	// Session Identify with a HYPOTHESIS operand stays HYPOTHESIS.
	s := session.New()
	if err := s.Draft("hyp_ident", "l", session.DraftMetadata{Hypothesis: h}); err != nil {
		t.Fatalf("draft: %v", err)
	}
	out, err := s.Identify(h, plain, "hypothesis identification attempt")
	if err != nil {
		t.Fatalf("Identify: %v", err)
	}
	if out.Provenance().Status() != core.StatusHypothesis {
		t.Errorf("Identify(h, x) status = %v, want HYPOTHESIS", out.Provenance().Status())
	}
}

// sealFlow builds a real sealed candidate whose derivation consumes the
// hypothesis concept, returning the session and candidate.
func sealFlow(t *testing.T, id string) (*session.Session, session.ResearchCandidate, core.Object) {
	t.Helper()
	h, err := hypothesis.NewCandidateConcept(id+"_h", core.KindEnergy, core.DimensionEnergy(), mustSymbol(t, "E"))
	if err != nil {
		t.Fatalf("candidate concept: %v", err)
	}
	eSym := mustSymbol(t, "E")
	rel := core.NewRelation(core.RelationEq, eSym, eSym)
	s := session.New()
	meta := session.DraftMetadata{
		DerivationID: id,
		Hypothesis:   h,
		Premises:     []core.Object{relativity.NewEnergy().CoreObject()},
		Assumptions:  core.NewAssumptionSet(),
		Predictions: []session.Prediction{{
			ID: "hp1", Observable: "E", Relation: rel, Assumptions: core.NewAssumptionSet(),
		}},
		FalsificationConditions: []session.FalsificationCondition{{
			ID: "hf1", TargetClaim: "hp1",
			ContradictingCondition: core.NewRelation(core.RelationNeq, eSym, eSym),
			Regime:                 "lab",
		}},
		AnomalyReferences: []session.AnomalyReference{{
			ID: "no_gravity", Framework: "special_relativity", Description: "flat only",
		}},
	}
	if err := s.Draft(id, "hypothesis flow", meta); err != nil {
		t.Fatalf("draft: %v", err)
	}
	energy := relativity.NewEnergy().CoreObject()
	if err := s.Declare("energy", energy); err != nil {
		t.Fatalf("declare: %v", err)
	}
	// Derive WITH the hypothesis concept: downstream stays HYPOTHESIS.
	mixed, err := s.Step("mix", ops.OpAdd, []core.Object{h, h}, ops.OperationParams{Kind: "empty"})
	if err != nil {
		t.Fatalf("step: %v", err)
	}
	if mixed.Provenance().Status() != core.StatusHypothesis {
		t.Fatalf("mixed status = %v, want HYPOTHESIS", mixed.Provenance().Status())
	}
	if err := s.Commit(); err != nil {
		t.Fatalf("commit: %v", err)
	}
	if err := s.Conclude(mixed); err != nil {
		t.Fatalf("conclude: %v", err)
	}
	cand, err := s.Seal()
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	return s, cand, h
}

// TestCandidateArtifactContainmentRejected crafts an artifact whose
// hypothesis-dependent derivation output is presented with a trusted
// hypothesis field and proves validation rejects it with
// CandidateContainmentError. Acceptance L (second half, REQ-032-09).
func TestCandidateArtifactContainmentRejected(t *testing.T) {
	_, cand, _ := sealFlow(t, "contain_l")
	raw, err := cand.CanonicalJSON()
	if err != nil {
		t.Fatalf("canonical: %v", err)
	}
	// Replace the HYPOTHESIS hypothesis with a trusted DEFINED object.
	trusted := mechanics.NewMass().CoreObject()
	trustedRaw, err := json.Marshal(mustCanonicalObject(t, trusted))
	if err != nil {
		t.Fatalf("trusted canonical: %v", err)
	}
	var wire map[string]json.RawMessage
	if err := json.Unmarshal(raw, &wire); err != nil {
		t.Fatalf("wire: %v", err)
	}
	wire["hypothesis"] = trustedRaw
	swapped, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("re-encode: %v", err)
	}
	uc, err := session.ParseResearchCandidateJSON(swapped)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	_, err = uc.Validate()
	if err == nil {
		t.Fatal("trusted-hypothesis artifact must fail validation")
	}
	var target core.CandidateContainmentError
	if !errors.As(err, &target) {
		t.Fatalf("error = %T (%v), want CandidateContainmentError", err, err)
	}
}

func mustCanonicalObject(t *testing.T, o core.Object) json.RawMessage {
	t.Helper()
	raw, err := kernel.CanonicalObjectJSON(o)
	if err != nil {
		t.Fatalf("canonical object: %v", err)
	}
	return json.RawMessage(raw)
}

// TestSealedCandidatePreservesFalsifiability seals predictions and
// falsification conditions and proves they survive with identical canonical
// values through CanonicalJSON round-trips. Acceptance M.
func TestSealedCandidatePreservesFalsifiability(t *testing.T) {
	_, cand, _ := sealFlow(t, "falsif_m")
	if len(cand.Predictions()) != 1 {
		t.Fatalf("predictions = %d, want 1", len(cand.Predictions()))
	}
	if len(cand.FalsificationConditions()) != 1 {
		t.Fatalf("falsification conditions = %d, want 1", len(cand.FalsificationConditions()))
	}
	p := cand.Predictions()[0]
	if p.ID != "hp1" || p.Observable != "E" || !p.Relation.Valid() {
		t.Errorf("prediction corrupted: %+v", p.ID)
	}
	raw, err := cand.CanonicalJSON()
	if err != nil {
		t.Fatalf("canonical: %v", err)
	}
	uc, err := session.ParseResearchCandidateJSON(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	rt, err := uc.Validate()
	if err != nil {
		t.Fatalf("revalidate: %v", err)
	}
	if len(rt.Predictions()) != 1 || !core.EqualExpr(rt.Predictions()[0].Relation, p.Relation) {
		t.Errorf("prediction relation not preserved through round-trip")
	}
	if len(rt.FalsificationConditions()) != 1 {
		t.Errorf("falsification conditions lost through round-trip")
	}
}

// TestCandidateReferencesManifestAnomaly seals an anomaly reference and
// cross-checks it externally against relativity/manifest.json anomalies
// (session never imports relativity). Acceptance N.
func TestCandidateReferencesManifestAnomaly(t *testing.T) {
	_, cand, _ := sealFlow(t, "anomaly_n")
	refs := cand.AnomalyReferences()
	if len(refs) != 1 {
		t.Fatalf("anomaly references = %d, want 1", len(refs))
	}
	if refs[0].ID != "no_gravity" || refs[0].Framework != "special_relativity" {
		t.Fatalf("anomaly reference corrupted: %+v", refs[0])
	}
	m, err := core.ParseManifest(relativity.ManifestJSON)
	if err != nil {
		t.Fatalf("parse relativity manifest: %v", err)
	}
	found := false
	for _, a := range m.Anomalies {
		if a.ID == refs[0].ID && a.Framework == refs[0].Framework {
			found = true
		}
	}
	if !found {
		t.Errorf("anomaly %q not present in relativity manifest", refs[0].ID)
	}
	sum := sha256.Sum256(relativity.ManifestJSON)
	_ = sum
}

// TestNoForbiddenExportedAPI scans the module AST for promotion, truth,
// bypass, and deferred-runtime APIs. None may exist.
func TestNoForbiddenExportedAPI(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	root := filepath.Dir(filepath.Dir(file))
	forbidden := []string{
		"Promote", "Trust", "ApproveHypothesis", "PromoteToEstablished",
		"SetCorpusStatus", "ApproveException", "OverrideMRC", "BypassMRC",
		"Integrate", "Series", "Taylor", "Simulate", "RankTheories",
		"ScoreTruth", "TRUTH_SCORE", "PROBABILITY_OF_TRUTH", "BEST_THEORY",
		"PHYSICALLY_TRUE",
	}
	fset := token.NewFileSet()
	var hits []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "plan10" || name == "plans" || name == ".opencode" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		ast.Inspect(f, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !ok || !id.IsExported() {
				return true
			}
			for _, ban := range forbidden {
				if id.Name == ban || strings.Contains(id.Name, ban) {
					hits = append(hits, fmt.Sprintf("%s: %s", path, id.Name))
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// Filter: references in comments/strings are not AST idents, so any hit
	// is a real identifier. Allow the test's own ban list literals: they are
	// BasicLit strings, not idents — no filtering needed.
	if len(hits) > 0 {
		t.Errorf("forbidden exported API hits:\n%s", strings.Join(hits, "\n"))
	}
}

// TestCandidateDerivationHashMatchesLedger proves the sealed ledger hash is
// the derivation hash end to end.
func TestCandidateDerivationHashMatchesLedger(t *testing.T) {
	s, cand, _ := sealFlow(t, "hash_n")
	if cand.LedgerHash() != s.Ledger().DerivationHash() {
		t.Errorf("ledger hash %s != derivation hash %s", cand.LedgerHash(), s.Ledger().DerivationHash())
	}
	raw, err := cand.CanonicalJSON()
	if err != nil {
		t.Fatalf("canonical: %v", err)
	}
	var wire map[string]json.RawMessage
	if err := json.Unmarshal(raw, &wire); err != nil {
		t.Fatalf("wire: %v", err)
	}
	var ledgerWire map[string]json.RawMessage
	if err := json.Unmarshal(wire["derivation"], &ledgerWire); err != nil {
		t.Fatalf("ledger wire: %v", err)
	}
	var finalHash string
	if err := json.Unmarshal(ledgerWire["final_hash"], &finalHash); err != nil {
		t.Fatalf("final hash: %v", err)
	}
	if finalHash != cand.LedgerHash() {
		t.Errorf("embedded final_hash %s != ledger hash %s", finalHash, cand.LedgerHash())
	}
	if !bytes.Contains(raw, []byte(`"schema_version":"1"`)) {
		t.Errorf("candidate canonical JSON lacks schema_version")
	}
}

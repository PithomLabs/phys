package session_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unsafe"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
	"github.com/PithomLabs/phys/mechanics"
	"github.com/PithomLabs/phys/ops"
	"github.com/PithomLabs/phys/session"
)

// ---------------------------------------------------------------------------
// Fixtures (in-module kernel minting; no production hooks)
// ---------------------------------------------------------------------------

func hypothesisFixture(t *testing.T, id string) core.Object {
	t.Helper()
	sym, err := core.NewSymbol("h_" + id)
	if err != nil {
		t.Fatalf("symbol: %v", err)
	}
	prov, err := kernel.NewProvenance(kernel.StatusHypothesis, "hypothesis", "", nil,
		[32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	if err != nil {
		t.Fatalf("provenance: %v", err)
	}
	obj, err := kernel.MintObject(kernel.ObjectSpec{
		Name: id, Kind: core.KindExpression, Dimension: core.Dimensionless(),
		Expr: sym, Assumptions: core.NewAssumptionSet(), Conventions: core.NewConventionSet(),
		Provenance: prov, CorpusStatus: kernel.CorpusNone,
	})
	if err != nil {
		t.Fatalf("mint hypothesis: %v", err)
	}
	return obj
}

func postulatedFixture(t *testing.T, name string) core.Object {
	t.Helper()
	sym, err := core.NewSymbol(name)
	if err != nil {
		t.Fatalf("symbol: %v", err)
	}
	prov, err := kernel.NewProvenance(kernel.StatusPostulated, "test", "test_framework", nil,
		[32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	if err != nil {
		t.Fatalf("provenance: %v", err)
	}
	obj, err := kernel.MintObject(kernel.ObjectSpec{
		Name: name, Kind: core.KindExpression, Dimension: core.Dimensionless(),
		Expr: sym, Assumptions: core.NewAssumptionSet(), Conventions: core.NewConventionSet(),
		Provenance: prov, CorpusStatus: kernel.CorpusNone,
	})
	if err != nil {
		t.Fatalf("mint postulated: %v", err)
	}
	return obj
}

func restMassKindFixture(t *testing.T) core.Object {
	t.Helper()
	sym, err := core.NewSymbol("m")
	if err != nil {
		t.Fatalf("symbol: %v", err)
	}
	prov, err := kernel.NewProvenance(kernel.StatusDefined, "test", "test_framework", nil,
		[32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	if err != nil {
		t.Fatalf("provenance: %v", err)
	}
	obj, err := kernel.MintObject(kernel.ObjectSpec{
		Name: "restmass", Kind: core.KindRestMass, Dimension: core.DimensionMass(),
		Expr: sym, Assumptions: core.NewAssumptionSet(), Conventions: core.NewConventionSet(),
		Provenance: prov, CorpusStatus: kernel.CorpusNone,
	})
	if err != nil {
		t.Fatalf("mint restmass-kind: %v", err)
	}
	return obj
}

func mustDraft(t *testing.T, id string) *session.Session {
	t.Helper()
	s := session.New()
	if err := s.Draft(id, "label_"+id, session.DraftMetadata{
		DerivationID: id,
		Hypothesis:   hypothesisFixture(t, id),
		Assumptions:  core.NewAssumptionSet(),
	}); err != nil {
		t.Fatalf("draft: %v", err)
	}
	return s
}

// simpleCommittedSession builds Draft → Define(mass) → add(mass,mass) → Commit.
func simpleCommittedSession(t *testing.T, id string) (*session.Session, core.Object) {
	t.Helper()
	s := mustDraft(t, id)
	mass := mechanics.NewMass().CoreObject()
	if err := s.Define("mass", mass); err != nil {
		t.Fatalf("define: %v", err)
	}
	sum, err := s.Step("add", ops.OpAdd, []core.Object{mass, mass}, ops.OperationParams{Kind: "empty"})
	if err != nil {
		t.Fatalf("step: %v", err)
	}
	if err := s.Commit(); err != nil {
		t.Fatalf("commit: %v", err)
	}
	return s, sum
}

// ---------------------------------------------------------------------------
// White-box fault injection on the session's OWN committed ledger.
// Uses reflect+unsafe in test only; no production hook exists.
// ---------------------------------------------------------------------------

func sessionLedgerSteps(t *testing.T, s *session.Session) reflect.Value {
	t.Helper()
	rv := reflect.ValueOf(s).Elem().FieldByName("ledger")
	ledger := reflect.NewAt(rv.Type(), unsafe.Pointer(rv.UnsafeAddr())).Elem()
	return ledger.FieldByName("steps")
}

func setSessionStepString(t *testing.T, s *session.Session, idx int, field, value string) {
	t.Helper()
	steps := sessionLedgerSteps(t, s)
	if idx < 0 || idx >= steps.Len() {
		t.Fatalf("step index %d out of range", idx)
	}
	st := reflect.NewAt(steps.Index(idx).Type(), unsafe.Pointer(steps.Index(idx).UnsafeAddr())).Elem()
	fh := st.FieldByName(field)
	reflect.NewAt(fh.Type(), unsafe.Pointer(fh.UnsafeAddr())).Elem().SetString(value)
}

func setSessionString(t *testing.T, s *session.Session, field, value string) {
	t.Helper()
	rv := reflect.ValueOf(s).Elem().FieldByName(field)
	reflect.NewAt(rv.Type(), unsafe.Pointer(rv.UnsafeAddr())).Elem().SetString(value)
}

func setSessionStepBytes(t *testing.T, s *session.Session, idx int, field string, value []byte) {
	t.Helper()
	steps := sessionLedgerSteps(t, s)
	if idx < 0 || idx >= steps.Len() {
		t.Fatalf("step index %d out of range", idx)
	}
	st := reflect.NewAt(steps.Index(idx).Type(), unsafe.Pointer(steps.Index(idx).UnsafeAddr())).Elem()
	fh := st.FieldByName(field)
	reflect.NewAt(fh.Type(), unsafe.Pointer(fh.UnsafeAddr())).Elem().SetBytes(value)
}

// mirrorStep mirrors the ledger JSON wire format for hash recomputation.
type mirrorStep struct {
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

func mirrorBodyHash(m mirrorStep) string {
	m.CurrentStepHash = ""
	raw, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return strings.ToLower(strings.TrimSpace(strings.ToUpper(""))) + lowerHex(sha256.Sum256(raw))
}

func lowerHex(h [32]byte) string {
	const digits = "0123456789abcdef"
	out := make([]byte, 64)
	for i, b := range h {
		out[2*i] = digits[b>>4]
		out[2*i+1] = digits[b&0xf]
	}
	return string(out)
}

func ledgerMirror(t *testing.T, s *session.Session) []mirrorStep {
	t.Helper()
	raw, err := s.Ledger().CanonicalJSON()
	if err != nil {
		t.Fatalf("ledger canonical: %v", err)
	}
	var dto struct {
		Steps []mirrorStep `json:"steps"`
	}
	if err := json.Unmarshal(raw, &dto); err != nil {
		t.Fatalf("mirror decode: %v", err)
	}
	return dto.Steps
}

// ---------------------------------------------------------------------------
// Draft / state machine
// ---------------------------------------------------------------------------

func TestSessionDraftAndState(t *testing.T) {
	s := session.New()
	if err := s.Draft("test_01", "test_derivation", session.DraftMetadata{
		DerivationID: "test_01",
		Hypothesis:   hypothesisFixture(t, "test_01"),
		Assumptions:  core.NewAssumptionSet(),
	}); err != nil {
		t.Fatalf("draft failed: %v", err)
	}
	if err := s.Draft("again", "x", session.DraftMetadata{}); err == nil {
		t.Error("second Draft must fail")
	} else {
		var target core.ProvenanceError
		if !errors.As(err, &target) {
			t.Errorf("second Draft error = %T, want ProvenanceError", err)
		}
	}
	if err := session.New().Draft("", "x", session.DraftMetadata{}); err == nil {
		t.Error("empty derivation ID must fail")
	}
}

func TestSessionStateMachineTransitions(t *testing.T) {
	// Validate with no ledger and illegal state must fail structurally
	// (stage-1 gate), never panic or pass.
	var stage1 core.LedgerValidationError
	if err := session.New().Validate(); !errors.As(err, &stage1) {
		t.Errorf("New().Validate() = %v, want LedgerValidationError", err)
	}
	// Step before Draft.
	if _, err := session.New().Step("x", ops.OpAdd, nil, ops.OperationParams{Kind: "empty"}); err == nil {
		t.Error("Step in New state must fail")
	}
	// Commit in New.
	if err := session.New().Commit(); err == nil {
		t.Error("Commit in New state must fail")
	}
	// Empty commit stays Drafting: a later valid commit must still succeed.
	s := mustDraft(t, "stm_empty")
	var lerr core.LedgerValidationError
	if err := s.Commit(); !errors.As(err, &lerr) {
		t.Fatalf("empty Commit error = %v, want LedgerValidationError", err)
	}
	mass := mechanics.NewMass().CoreObject()
	if err := s.Define("mass", mass); err != nil {
		t.Fatalf("define after empty commit: %v", err)
	}
	if _, err := s.Step("add", ops.OpAdd, []core.Object{mass, mass}, ops.OperationParams{Kind: "empty"}); err != nil {
		t.Fatalf("step after empty commit: %v", err)
	}
	if err := s.Commit(); err != nil {
		t.Fatalf("commit after empty commit: %v", err)
	}
	// Conclude before Commit.
	s2 := mustDraft(t, "stm_conc")
	if err := s2.Conclude(mass); err == nil {
		t.Error("Conclude in Drafting must fail")
	}
	// Seal before Concluded.
	s3, _ := simpleCommittedSession(t, "stm_seal")
	if err := func() error { _, err := s3.Seal(); return err }(); err == nil {
		t.Error("Seal in Committed must fail")
	}
	// Post-seal mutations rejected.
	s4, sum := simpleCommittedSession(t, "stm_post")
	if err := s4.Conclude(sum); err != nil {
		t.Fatalf("conclude: %v", err)
	}
	if _, err := s4.Seal(); err != nil {
		t.Fatalf("seal: %v", err)
	}
	var perr core.ProvenanceError
	if _, err := s4.Step("x", ops.OpAdd, []core.Object{mass, mass}, ops.OperationParams{Kind: "empty"}); !errors.As(err, &perr) {
		t.Errorf("post-seal Step error = %v, want ProvenanceError", err)
	}
	if _, err := s4.Identify(mass, mass, "j"); !errors.As(err, &perr) {
		t.Errorf("post-seal Identify error = %v, want ProvenanceError", err)
	}
	if err := s4.Commit(); !errors.As(err, &perr) {
		t.Errorf("post-seal Commit error = %v, want ProvenanceError", err)
	}
}

// ---------------------------------------------------------------------------
// K — Identify firewall (MRC-006)
// ---------------------------------------------------------------------------

func TestIdentifyRequiresJustification(t *testing.T) {
	s := mustDraft(t, "identify_k")
	mass := mechanics.NewMass().CoreObject()
	for _, j := range []string{"", "   ", "\t\n "} {
		if _, err := s.Identify(mass, mass, j); err == nil {
			t.Errorf("justification %q must fail", j)
		} else {
			var target core.IdentifyError
			if !errors.As(err, &target) {
				t.Errorf("justification %q error = %T, want IdentifyError", j, err)
			}
		}
	}
}

func TestSessionIdentifyRecords(t *testing.T) {
	s := mustDraft(t, "identify_rec")
	massA := mechanics.NewMass().CoreObject()
	massB := mechanics.NewMass().CoreObject()
	got, err := s.Identify(massA, massB, "  linking mass observations  ")
	if err != nil {
		t.Fatalf("Identify: %v", err)
	}
	if got.Provenance().Status() != core.StatusIdentified {
		t.Fatalf("status = %v, want IDENTIFIED", got.Provenance().Status())
	}
	if got.Provenance().Justification() != "linking mass observations" {
		t.Errorf("justification not trimmed/recorded: %q", got.Provenance().Justification())
	}
	if got.Kind() != core.KindRelation || got.Expr().RelationOperator() != core.RelationEq {
		t.Errorf("identify result must be an eq Relation")
	}
	if !core.EqualExpr(got.Expr().Left(), massA.Expr()) || !core.EqualExpr(got.Expr().Right(), massB.Expr()) {
		t.Errorf("identify relation sides mismatch")
	}
	// Wrong-kind identification at equal dimensions must fail (MRC-003).
	rest := restMassKindFixture(t)
	if _, err := s.Identify(massA, rest, "x"); err == nil {
		t.Error("Mass vs RestMass Identify must fail")
	} else {
		var target core.CategoryMismatchError
		if !errors.As(err, &target) {
			t.Errorf("kind-mismatch error = %T, want CategoryMismatchError", err)
		}
	}
	// HYPOTHESIS operand forces HYPOTHESIS output with attempt preserved.
	h := hypothesisFixture(t, "identify_rec_h")
	h2 := hypothesisFixture(t, "identify_rec_h2")
	hgot, err := s.Identify(h, h2, "hypothesis attempt")
	if err != nil {
		t.Fatalf("hypothesis Identify: %v", err)
	}
	if hgot.Provenance().Status() != core.StatusHypothesis {
		t.Errorf("hypothesis identify status = %v, want HYPOTHESIS", hgot.Provenance().Status())
	}
	if hgot.Provenance().Justification() != "hypothesis attempt" {
		t.Errorf("hypothesis attempt justification lost: %q", hgot.Provenance().Justification())
	}
	// Commit and prove the Identification step is recorded and validates.
	if err := s.Commit(); err != nil {
		t.Fatalf("commit: %v", err)
	}
	found := false
	for _, st := range s.Ledger().Steps() {
		if st.StepKind() == session.StepKindIdentification {
			found = true
		}
	}
	if !found {
		t.Error("no Identification step recorded in ledger")
	}
	if err := s.Validate(); err != nil {
		t.Errorf("Validate after Identify: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Surface: all nine actions + read-only views in one valid flow
// ---------------------------------------------------------------------------

func TestSessionSurfaceNineActions(t *testing.T) {
	s := mustDraft(t, "surface_01")
	mass := mechanics.NewMass().CoreObject()
	if err := s.Postulate("p", postulatedFixture(t, "g")); err != nil {
		t.Fatalf("Postulate: %v", err)
	}
	if err := s.Declare("d", mass); err != nil {
		t.Fatalf("Declare: %v", err)
	}
	if err := s.Define("m", mass); err != nil {
		t.Fatalf("Define: %v", err)
	}
	sum, err := s.Step("add", ops.OpAdd, []core.Object{mass, mass}, ops.OperationParams{Kind: "empty"})
	if err != nil {
		t.Fatalf("Step: %v", err)
	}
	if _, err := s.Identify(mass, mass, "surface"); err != nil {
		t.Fatalf("Identify: %v", err)
	}
	if err := s.Commit(); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	// Conclude must target the final output (the identification result).
	last := s.Ledger().Steps()[len(s.Ledger().Steps())-1].OutputCanonical()
	_ = sum
	if err := s.Conclude(last); err != nil {
		t.Fatalf("Conclude: %v", err)
	}
	if err := s.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if _, err := s.CanonicalJSON(); err != nil {
		t.Fatalf("CanonicalJSON: %v", err)
	}
	if _, err := s.Seal(); err != nil {
		t.Fatalf("Seal: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Import independence (production graph audits, test-side)
// ---------------------------------------------------------------------------

func repoDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	return filepath.Dir(file)
}

func nonTestImports(t *testing.T, dir string) map[string][]string {
	t.Helper()
	fset := token.NewFileSet()
	found := map[string][]string{}
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		name := fi.Name()
		return !fi.IsDir() && !(len(name) > 8 && name[len(name)-8:] == "_test.go")
	}, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse %s: %v", dir, err)
	}
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			for _, imp := range f.Imports {
				path := strings.Trim(imp.Path.Value, `"`)
				found[path] = append(found[path], fset.Position(imp.Pos()).Filename)
			}
		}
	}
	return found
}

func TestSessionImportIndependence(t *testing.T) {
	for imp := range nonTestImports(t, repoDir(t)) {
		if strings.Contains(imp, "PithomLabs/phys/mechanics") ||
			strings.Contains(imp, "PithomLabs/phys/relativity") ||
			strings.Contains(imp, "PithomLabs/phys/hypothesis") {
			t.Errorf("session production code imports domain package %s", imp)
		}
	}
}

func TestOpsImportsNoSession(t *testing.T) {
	opsDir := filepath.Join(filepath.Dir(repoDir(t)), "ops")
	for imp := range nonTestImports(t, opsDir) {
		if strings.Contains(imp, "PithomLabs/phys/session") {
			t.Errorf("ops production code imports session: %s", imp)
		}
	}
}

// ---------------------------------------------------------------------------
// Ledger exterior: strict parsing + artifact-level tamper
// ---------------------------------------------------------------------------

func TestLedgerParseAndValidateExists(t *testing.T) {
	if _, err := session.ParseLedgerJSON([]byte(`{"schema_version":"1","steps":[],"final_hash":"","sealed":false}`)); err != nil {
		t.Fatalf("parse well-formed empty ledger: %v", err)
	} else {
		// Empty steps parse but never validate.
		l, _ := session.ParseLedgerJSON([]byte(`{"schema_version":"1","steps":[],"final_hash":"","sealed":false}`))
		var target core.LedgerValidationError
		if err := l.Validate(); !errors.As(err, &target) {
			t.Errorf("empty ledger Validate = %v, want LedgerValidationError", err)
		}
	}
	if _, err := session.ParseLedgerJSON([]byte(`{"steps":[]}`)); err == nil {
		t.Error("missing schema_version must fail strict parse")
	}
	if _, err := session.ParseLedgerJSON([]byte(`{"schema_version":"1","steps":[],"final_hash":"","sealed":false,"extra":1}`)); err == nil {
		t.Error("unknown field must fail strict parse")
	}
}

func TestLedgerTamperViaParse(t *testing.T) {
	s, _ := simpleCommittedSession(t, "parse_tamper")
	raw, err := s.Ledger().CanonicalJSON()
	if err != nil {
		t.Fatalf("canonical: %v", err)
	}
	var dto map[string]json.RawMessage
	if err := json.Unmarshal(raw, &dto); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	// Decode steps into the byte-preserving mirror struct (a generic map
	// round-trip would re-sort keys and break canonicality, masking the
	// hash check under test).
	var steps []mirrorStep
	if err := json.Unmarshal(dto["steps"], &steps); err != nil {
		t.Fatalf("steps: %v", err)
	}
	// Tamper the LAST step only: linkage is intact, so only the
	// current-step hash recomputation can catch it.
	steps[len(steps)-1].CurrentStepHash = "0000000000000000000000000000000000000000000000000000000000000000"
	rawSteps, err := json.Marshal(steps)
	if err != nil {
		t.Fatalf("steps encode: %v", err)
	}
	dto["steps"] = rawSteps
	tampered, _ := json.Marshal(dto)
	l2, err := session.ParseLedgerJSON(tampered)
	if err != nil {
		t.Fatalf("parse tampered: %v", err)
	}
	var target core.LedgerValidationError
	err = l2.Validate()
	if !errors.As(err, &target) {
		t.Errorf("tampered ledger Validate = %v, want LedgerValidationError", err)
	} else if !strings.Contains(err.Error(), "current-step") {
		t.Errorf("tampered ledger failed at wrong stage: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Q — ledger tamper/replay on the SESSION'S OWN committed ledger (§16.20)
// ---------------------------------------------------------------------------

func TestLedgerTamperDetection(t *testing.T) {
	t.Run("output_only", func(t *testing.T) {
		s, _ := simpleCommittedSession(t, "q_out")
		steps := s.Ledger().Steps()
		raw := steps[0].OutputCanonicalBytes()
		// Flip a symbol byte inside the canonical object, preserving valid JSON.
		idx := bytes.Index(raw, []byte(`"name":"m"`))
		if idx < 0 {
			t.Fatalf("anchor not found in %s", raw)
		}
		raw[idx+len(`"name":"`)] = 'n'
		setSessionStepBytes(t, s, 0, "OutputCanonical", raw)
		var target core.LedgerValidationError
		err := s.Validate()
		if !errors.As(err, &target) {
			t.Fatalf("output-only tamper Validate = %v, want LedgerValidationError", err)
		}
		if !strings.Contains(err.Error(), "output hash") {
			t.Errorf("output-only tamper failed at wrong stage: %v", err)
		}
	})
	t.Run("output_and_hash", func(t *testing.T) {
		s, _ := simpleCommittedSession(t, "q_outhash")
		steps := s.Ledger().Steps()
		raw := append([]byte(nil), steps[0].OutputCanonicalBytes()...)
		idx := bytes.Index(raw, []byte(`"name":"m"`))
		if idx < 0 {
			t.Fatalf("anchor not found")
		}
		raw[idx+len(`"name":"`)] = 'n'
		obj, err := kernel.LoadObjectJSON(raw)
		if err != nil {
			t.Fatalf("tampered output must still decode: %v", err)
		}
		setSessionStepBytes(t, s, 0, "OutputCanonical", raw)
		setSessionStepString(t, s, 0, "OutputHash", lowerHex(kernel.HashObject(obj)))
		var target core.LedgerValidationError
		err = s.Validate()
		if !errors.As(err, &target) {
			t.Fatalf("output+hash tamper Validate = %v, want LedgerValidationError", err)
		}
		if !strings.Contains(err.Error(), "current-step") {
			t.Errorf("output+hash tamper failed at wrong stage: %v", err)
		}
	})
	t.Run("full_chain_recompute_diverges_on_replay", func(t *testing.T) {
		s, _ := simpleCommittedSession(t, "q_replay")
		mirror := ledgerMirror(t, s)
		if len(mirror) == 0 {
			t.Fatal("no steps")
		}
		// Fidelity check: mirror recomputation must reproduce the stored hash.
		if got := mirrorBodyHash(mirror[1]); got != mirror[1].CurrentStepHash {
			t.Fatalf("mirror unfaithful: recomputed %s vs stored %s", got, mirror[1].CurrentStepHash)
		}
		// Substitute a different valid object as retained output and
		// recompute every dependent hash consistently. The alt object is
		// DERIVED (like the original add result) so that only transformation
		// replay can distinguish it — structural checks all pass.
		timeExpr := mechanics.NewTime().CoreObject().Expr()
		altProv, err := kernel.NewProvenance(kernel.StatusDerived, "", "", nil,
			[32]byte{}, [32]byte{}, kernel.MRCVersion, "")
		if err != nil {
			t.Fatalf("provenance: %v", err)
		}
		timeObj, err := kernel.MintObject(kernel.ObjectSpec{
			Name: "", Kind: core.KindExpression, Dimension: core.DimensionTime(),
			Expr: timeExpr, Assumptions: core.NewAssumptionSet(),
			Conventions: core.NewConventionSet(), Provenance: altProv,
			CorpusStatus: kernel.CorpusNone,
		})
		if err != nil {
			t.Fatalf("mint alt: %v", err)
		}
		timeRaw, err := kernel.CanonicalObjectJSON(timeObj)
		if err != nil {
			t.Fatalf("canonical time: %v", err)
		}
		m := mirror[1]
		m.OutputCanonical = json.RawMessage(timeRaw)
		m.OutputHash = lowerHex(kernel.HashObject(timeObj))
		m.AssumptionHash = lowerHex(core.HashAssumptionSet(timeObj.Assumptions()))
		m.ConventionHash = lowerHex(core.HashConventionSet(timeObj.Conventions()))
		m.CurrentStepHash = mirrorBodyHash(m)
		setSessionStepBytes(t, s, 1, "OutputCanonical", m.OutputCanonical)
		setSessionStepString(t, s, 1, "OutputHash", m.OutputHash)
		setSessionStepString(t, s, 1, "AssumptionHash", m.AssumptionHash)
		setSessionStepString(t, s, 1, "ConventionHash", m.ConventionHash)
		setSessionStepString(t, s, 1, "CurrentStepHash", m.CurrentStepHash)
		var target core.LedgerValidationError
		err = s.Validate()
		if !errors.As(err, &target) {
			t.Fatalf("recomputed-chain tamper Validate = %v, want LedgerValidationError (replay divergence)", err)
		}
		if !strings.Contains(err.Error(), "replayed output") {
			t.Errorf("recomputed-chain tamper failed at wrong stage: %v", err)
		}
	})
	t.Run("step_id_format", func(t *testing.T) {
		s, _ := simpleCommittedSession(t, "q_stepid")
		// Corrupt the StepID string but recompute the chain hash, so only
		// the StepID/index format check can fire.
		mirror := ledgerMirror(t, s)
		m := mirror[1]
		m.StepID = "step-000099"
		m.CurrentStepHash = mirrorBodyHash(m)
		setSessionStepString(t, s, 1, "StepID", m.StepID)
		setSessionStepString(t, s, 1, "CurrentStepHash", m.CurrentStepHash)
		var target core.LedgerValidationError
		err := s.Validate()
		if !errors.As(err, &target) {
			t.Fatalf("step-id tamper Validate = %v, want LedgerValidationError", err)
		}
		if !strings.Contains(err.Error(), "step ID format") {
			t.Errorf("step-id tamper failed at wrong stage: %v", err)
		}
	})
	t.Run("broken_chain", func(t *testing.T) {
		s, _ := simpleCommittedSession(t, "q_chain")
		// Break linkage but recompute the step hash so only the linkage
		// check itself can fire.
		mirror := ledgerMirror(t, s)
		m := mirror[1]
		m.PreviousStepHash = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
		m.CurrentStepHash = mirrorBodyHash(m)
		setSessionStepString(t, s, 1, "PreviousStepHash", m.PreviousStepHash)
		setSessionStepString(t, s, 1, "CurrentStepHash", m.CurrentStepHash)
		var target core.LedgerValidationError
		err := s.Validate()
		if !errors.As(err, &target) {
			t.Fatalf("broken chain Validate = %v, want LedgerValidationError", err)
		}
		if !strings.Contains(err.Error(), "chain linkage") {
			t.Errorf("broken chain failed at wrong stage: %v", err)
		}
	})
	t.Run("params_binding", func(t *testing.T) {
		s, _ := simpleCommittedSession(t, "q_params")
		// An add step carrying compare params must fail stage 10. The chain
		// hash is recomputed so stages 1-9 pass and binding is isolated.
		badParams := []byte(`{"kind":"compare","exponent":"","operator":"eq","justification":""}`)
		mirror := ledgerMirror(t, s)
		m := mirror[1]
		m.ParamsCanonical = json.RawMessage(badParams)
		m.CurrentStepHash = mirrorBodyHash(m)
		setSessionStepBytes(t, s, 1, "ParamsCanonical", badParams)
		setSessionStepString(t, s, 1, "CurrentStepHash", m.CurrentStepHash)
		var target core.LedgerValidationError
		err := s.Validate()
		if !errors.As(err, &target) {
			t.Fatalf("params tamper Validate = %v, want LedgerValidationError", err)
		}
		if !strings.Contains(err.Error(), "operation params:") {
			t.Errorf("params tamper failed at wrong stage: %v", err)
		}
	})
	t.Run("provenance_justification", func(t *testing.T) {
		// An IDENTIFIED output whose provenance lacks justification is
		// structurally consistent (hashes recomputed) but violates the
		// stage-15 provenance law: replay divergence alone cannot catch it.
		s := mustDraft(t, "q_prov")
		mass := mechanics.NewMass().CoreObject()
		if _, err := s.Identify(mass, mass, "recorded reason"); err != nil {
			t.Fatalf("identify: %v", err)
		}
		if err := s.Commit(); err != nil {
			t.Fatalf("commit: %v", err)
		}
		mirror := ledgerMirror(t, s)
		m := mirror[0]
		raw := append([]byte(nil), m.OutputCanonical...)
		// Strip the justification inside the retained canonical bytes.
		oldJ := `"justification":"recorded reason"`
		if !bytes.Contains(raw, []byte(oldJ)) {
			t.Fatalf("justification anchor missing in %s", raw)
		}
		raw = bytes.Replace(raw, []byte(oldJ), []byte(`"justification":""`), 1)
		obj, err := kernel.LoadObjectJSON(raw)
		if err != nil {
			t.Fatalf("stripped output must still decode: %v", err)
		}
		m.OutputCanonical = json.RawMessage(raw)
		m.OutputHash = lowerHex(kernel.HashObject(obj))
		m.CurrentStepHash = mirrorBodyHash(m)
		setSessionStepBytes(t, s, 0, "OutputCanonical", m.OutputCanonical)
		setSessionStepString(t, s, 0, "OutputHash", m.OutputHash)
		setSessionStepString(t, s, 0, "CurrentStepHash", m.CurrentStepHash)
		var target core.LedgerValidationError
		if err := s.Validate(); !errors.As(err, &target) {
			t.Fatalf("justification-stripped Validate = %v, want LedgerValidationError", err)
		}
	})
	t.Run("input_only", func(t *testing.T) {
		s, _ := simpleCommittedSession(t, "q_input")
		steps := s.Ledger().Steps()
		raw := append([]byte(nil), steps[1].InputCanonicalBytes()[0]...)
		idx := bytes.Index(raw, []byte(`"name":"m"`))
		if idx < 0 {
			t.Fatalf("anchor missing")
		}
		raw[idx+len(`"name":"`)] = 'q'
		// Replace first input bytes via unsafe slice elem.
		stepsVal := sessionLedgerSteps(t, s)
		st := reflect.NewAt(stepsVal.Index(1).Type(), unsafe.Pointer(stepsVal.Index(1).UnsafeAddr())).Elem()
		fh := st.FieldByName("InputCanonicals")
		elem := reflect.NewAt(fh.Type().Elem(), unsafe.Pointer(fh.Index(0).UnsafeAddr())).Elem()
		elem.SetBytes(raw)
		var target core.LedgerValidationError
		err := s.Validate()
		if !errors.As(err, &target) {
			t.Fatalf("input-only tamper Validate = %v, want LedgerValidationError", err)
		}
		if !strings.Contains(err.Error(), "input hash") {
			t.Errorf("input-only tamper failed at wrong stage: %v", err)
		}
	})
	t.Run("output_hash_only", func(t *testing.T) {
		s, _ := simpleCommittedSession(t, "q_outhashonly")
		// Only the hash string changes; bytes stay canonical and consistent.
		setSessionStepString(t, s, 1, "OutputHash",
			"0000000000000000000000000000000000000000000000000000000000000000")
		var target core.LedgerValidationError
		err := s.Validate()
		if !errors.As(err, &target) {
			t.Fatalf("output-hash-only tamper Validate = %v, want LedgerValidationError", err)
		}
		if !strings.Contains(err.Error(), "output hash") {
			t.Errorf("output-hash-only tamper failed at wrong stage: %v", err)
		}
	})
	t.Run("identification_replay", func(t *testing.T) {
		s := mustDraft(t, "q_ident_replay")
		mass := mechanics.NewMass().CoreObject()
		if _, err := s.Identify(mass, mass, "replay me"); err != nil {
			t.Fatalf("identify: %v", err)
		}
		if err := s.Commit(); err != nil {
			t.Fatalf("commit: %v", err)
		}
		mirror := ledgerMirror(t, s)
		m := mirror[0]
		if got := mirrorBodyHash(m); got != m.CurrentStepHash {
			t.Fatalf("mirror unfaithful")
		}
		// Swap retained output for a different valid relation, recompute all
		// dependent hashes: only identification replay can catch this.
		alt, err := kernel.LoadObjectJSON(m.OutputCanonical)
		if err != nil {
			t.Fatalf("decode: %v", err)
		}
		_ = alt
		timeObj := mechanics.NewTime().CoreObject()
		relExpr := core.NewRelation(core.RelationEq, timeObj.Expr(), timeObj.Expr())
		prov, err := kernel.NewProvenance(kernel.StatusIdentified, "", "", nil,
			[32]byte{}, [32]byte{}, kernel.MRCVersion, "replay me")
		if err != nil {
			t.Fatalf("provenance: %v", err)
		}
		altObj, err := kernel.MintObject(kernel.ObjectSpec{
			Name: "", Kind: core.KindRelation, Dimension: core.DimensionTime(),
			Expr: relExpr, Assumptions: core.NewAssumptionSet(),
			Conventions: core.NewConventionSet(), Provenance: prov,
			CorpusStatus: kernel.CorpusNone,
		})
		if err != nil {
			t.Fatalf("mint alt: %v", err)
		}
		altRaw, err := kernel.CanonicalObjectJSON(altObj)
		if err != nil {
			t.Fatalf("canonical alt: %v", err)
		}
		m.OutputCanonical = json.RawMessage(altRaw)
		m.OutputHash = lowerHex(kernel.HashObject(altObj))
		m.CurrentStepHash = mirrorBodyHash(m)
		setSessionStepBytes(t, s, 0, "OutputCanonical", m.OutputCanonical)
		setSessionStepString(t, s, 0, "OutputHash", m.OutputHash)
		setSessionStepString(t, s, 0, "CurrentStepHash", m.CurrentStepHash)
		var target core.LedgerValidationError
		err = s.Validate()
		if !errors.As(err, &target) {
			t.Fatalf("identification-replay tamper Validate = %v, want LedgerValidationError", err)
		}
		if !strings.Contains(err.Error(), "replayed output") {
			t.Errorf("identification-replay tamper failed at wrong stage: %v", err)
		}
	})
	t.Run("containment_corpus_status", func(t *testing.T) {
		// Assertion replay is vacuous (retained object is authoritative), so
		// a Declare step whose hypothesis-derived material carries a trusted
		// corpus status passes stages 1-16 and must fail stage 17 with
		// CandidateContainmentError. Transformation steps cannot reach stage
		// 17 inconsistently because replay equality (stage 14) subsumes them.
		s := mustDraft(t, "q_contain")
		h := hypothesisFixture(t, "q_contain_h")
		if err := s.Declare("h", h); err != nil {
			t.Fatalf("declare: %v", err)
		}
		if err := s.Commit(); err != nil {
			t.Fatalf("commit: %v", err)
		}
		mirror := ledgerMirror(t, s)
		m := mirror[0]
		upgrade := func(raw json.RawMessage) json.RawMessage {
			if !bytes.Contains(raw, []byte(`"corpus_status":"none"`)) {
				t.Fatalf("corpus anchor missing in %s", raw)
			}
			out := bytes.Replace(raw, []byte(`"corpus_status":"none"`), []byte(`"corpus_status":"established"`), 1)
			if _, err := kernel.LoadObjectJSON(out); err != nil {
				t.Fatalf("upgraded material must still decode: %v", err)
			}
			return out
		}
		inUp := upgrade(m.InputCanonicals[0])
		outUp := upgrade(m.OutputCanonical)
		inObj, _ := kernel.LoadObjectJSON(inUp)
		outObj, _ := kernel.LoadObjectJSON(outUp)
		m.InputCanonicals[0] = inUp
		m.InputHashes[0] = lowerHex(kernel.HashObject(inObj))
		m.OutputCanonical = outUp
		m.OutputHash = lowerHex(kernel.HashObject(outObj))
		m.CurrentStepHash = mirrorBodyHash(m)
		stepsVal := sessionLedgerSteps(t, s)
		st := reflect.NewAt(stepsVal.Index(0).Type(), unsafe.Pointer(stepsVal.Index(0).UnsafeAddr())).Elem()
		fh := st.FieldByName("InputCanonicals")
		elem := reflect.NewAt(fh.Type().Elem(), unsafe.Pointer(fh.Index(0).UnsafeAddr())).Elem()
		elem.SetBytes(inUp)
		// InputHashes is []string: rewrite the element directly.
		ih := st.FieldByName("InputHashes")
		ihElem := reflect.NewAt(ih.Type().Elem(), unsafe.Pointer(ih.Index(0).UnsafeAddr())).Elem()
		ihElem.SetString(lowerHex(kernel.HashObject(inObj)))
		setSessionStepBytes(t, s, 0, "OutputCanonical", outUp)
		setSessionStepString(t, s, 0, "OutputHash", lowerHex(kernel.HashObject(outObj)))
		setSessionStepString(t, s, 0, "CurrentStepHash", mirrorBodyHash(m))
		err := s.Validate()
		if err == nil {
			t.Fatal("trusted-status hypothesis output must fail validation")
		}
		var target core.CandidateContainmentError
		if !errors.As(err, &target) {
			t.Fatalf("containment error = %T (%v), want CandidateContainmentError", err, err)
		}
	})
	t.Run("conclusion_tamper", func(t *testing.T) {
		s, sum := simpleCommittedSession(t, "q_conc")
		if err := s.Conclude(sum); err != nil {
			t.Fatalf("conclude: %v", err)
		}
		setSessionString(t, s, "conclusionHash",
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
		var target core.LedgerValidationError
		err := s.Validate()
		if !errors.As(err, &target) {
			t.Fatalf("conclusion tamper Validate = %v, want LedgerValidationError", err)
		}
		if !strings.Contains(err.Error(), "conclusion") {
			t.Errorf("conclusion tamper failed at wrong stage: %v", err)
		}
	})
	t.Run("assertion_replay", func(t *testing.T) {
		s := mustDraft(t, "q_assert")
		mass := mechanics.NewMass().CoreObject()
		if err := s.Declare("mass", mass); err != nil {
			t.Fatalf("declare: %v", err)
		}
		if err := s.Commit(); err != nil {
			t.Fatalf("commit: %v", err)
		}
		mirror := ledgerMirror(t, s)
		m := mirror[0]
		timeObj := mechanics.NewTime().CoreObject()
		timeRaw, err := kernel.CanonicalObjectJSON(timeObj)
		if err != nil {
			t.Fatalf("canonical time: %v", err)
		}
		m.OutputCanonical = json.RawMessage(timeRaw)
		m.OutputHash = lowerHex(kernel.HashObject(timeObj))
		m.CurrentStepHash = mirrorBodyHash(m)
		setSessionStepBytes(t, s, 0, "OutputCanonical", m.OutputCanonical)
		setSessionStepString(t, s, 0, "OutputHash", m.OutputHash)
		setSessionStepString(t, s, 0, "CurrentStepHash", m.CurrentStepHash)
		var target core.LedgerValidationError
		err = s.Validate()
		if !errors.As(err, &target) {
			t.Fatalf("assertion-replay tamper Validate = %v, want LedgerValidationError", err)
		}
		if !strings.Contains(err.Error(), "assertion") {
			t.Errorf("assertion-replay tamper failed at wrong stage: %v", err)
		}
	})
	t.Run("status_string", func(t *testing.T) {
		// Only the provenance-status string changes (output bytes intact);
		// chain recomputed. Structural hash stages pass; the stage-15
		// status/version reconciliation must fire.
		s, _ := simpleCommittedSession(t, "q_status")
		mirror := ledgerMirror(t, s)
		m := mirror[1]
		m.ProvenanceStatus = "POSTULATED"
		m.CurrentStepHash = mirrorBodyHash(m)
		setSessionStepString(t, s, 1, "ProvenanceStatus", m.ProvenanceStatus)
		setSessionStepString(t, s, 1, "CurrentStepHash", m.CurrentStepHash)
		var target core.LedgerValidationError
		err := s.Validate()
		if !errors.As(err, &target) {
			t.Fatalf("status-string tamper Validate = %v, want LedgerValidationError", err)
		}
		if !strings.Contains(err.Error(), "provenance status mismatch") {
			t.Errorf("status-string tamper failed at wrong stage: %v", err)
		}
	})
	t.Run("identification_justification_strip", func(t *testing.T) {
		s := mustDraft(t, "q_ident")
		mass := mechanics.NewMass().CoreObject()
		if _, err := s.Identify(mass, mass, "keep me"); err != nil {
			t.Fatalf("identify: %v", err)
		}
		if err := s.Commit(); err != nil {
			t.Fatalf("commit: %v", err)
		}
		mirror := ledgerMirror(t, s)
		mj := mirror[0]
		mj.ParamsCanonical = json.RawMessage([]byte(`{"kind":"identify","exponent":"","operator":"","justification":""}`))
		mj.CurrentStepHash = mirrorBodyHash(mj)
		setSessionStepBytes(t, s, 0, "ParamsCanonical", mj.ParamsCanonical)
		setSessionStepString(t, s, 0, "CurrentStepHash", mj.CurrentStepHash)
		var target core.LedgerValidationError
		if err := s.Validate(); !errors.As(err, &target) {
			t.Fatalf("stripped justification Validate = %v, want LedgerValidationError", err)
		}
	})
}

// ---------------------------------------------------------------------------
// R — determinism
// ---------------------------------------------------------------------------

func TestDerivationDeterminism(t *testing.T) {
	s1, _ := simpleCommittedSession(t, "det_same")
	s2, _ := simpleCommittedSession(t, "det_same")
	raw1, err := s1.Ledger().CanonicalJSON()
	if err != nil {
		t.Fatalf("canonical 1: %v", err)
	}
	raw2, err := s2.Ledger().CanonicalJSON()
	if err != nil {
		t.Fatalf("canonical 2: %v", err)
	}
	if !bytes.Equal(raw1, raw2) {
		t.Fatalf("same derivation diverged:\n%s\n%s", raw1, raw2)
	}
	if s1.Ledger().DerivationHash() != s2.Ledger().DerivationHash() {
		t.Fatalf("derivation hash diverged")
	}
	// Fresh session, same operations, same ID → identical bytes.
	s3, _ := simpleCommittedSession(t, "det_same")
	raw3, _ := s3.Ledger().CanonicalJSON()
	if !bytes.Equal(raw1, raw3) {
		t.Fatalf("fresh-session derivation diverged")
	}
	// Session-level canonical JSON is deterministic too.
	sj1, _ := s1.CanonicalJSON()
	sj2, _ := s2.CanonicalJSON()
	if !bytes.Equal(sj1, sj2) {
		t.Fatalf("session canonical JSON diverged")
	}
}

// ---------------------------------------------------------------------------
// S — sealed handoff
// ---------------------------------------------------------------------------

func sealFixtureMetadata(t *testing.T, id string, h core.Object) session.DraftMetadata {
	t.Helper()
	eSym, err := core.NewSymbol("E")
	if err != nil {
		t.Fatalf("symbol: %v", err)
	}
	rel := core.NewRelation(core.RelationEq, eSym, eSym)
	contra := core.NewRelation(core.RelationNeq, eSym, eSym)
	sum := sha256.Sum256([]byte("classical_mechanics"))
	return session.DraftMetadata{
		DerivationID: id,
		Hypothesis:   h,
		Premises:     []core.Object{mechanics.NewMass().CoreObject()},
		Assumptions:  core.NewAssumptionSet(),
		FrameworkDependencies: []session.FrameworkDependency{{
			FrameworkID: "classical_mechanics", ManifestHash: lowerHex(sum),
			AssumptionHashes: []string{"abc"},
		}},
		Predictions: []session.Prediction{{
			ID: "p1", Observable: "E", Relation: rel, Assumptions: core.NewAssumptionSet(),
		}},
		FalsificationConditions: []session.FalsificationCondition{{
			ID: "f1", TargetClaim: "p1", ContradictingCondition: contra, Regime: "lab",
		}},
		RecoveryClaims: []session.RecoveryClaim{{
			ID: "r1", Description: "rec", FromFramework: "classical_mechanics", Condition: rel,
		}},
		AnomalyReferences: []session.AnomalyReference{{
			ID: "no_gravity", Framework: "special_relativity", Description: "flat only",
		}},
		ReviewHistory: []core.Review{{
			DerivationID: id,
			Challenges: []core.Challenge{{
				StepID: "step-000001", Category: string(core.ReviewCategoryDimensionError),
				Severity: "low", Description: "checked",
			}},
			ReviewerNotes: "notes",
		}},
	}
}

func TestSealResearchCandidate(t *testing.T) {
	h := hypothesisFixture(t, "seal_s")
	s := session.New()
	if err := s.Draft("seal_s", "sealed flow", sealFixtureMetadata(t, "seal_s", h)); err != nil {
		t.Fatalf("draft: %v", err)
	}
	mass := mechanics.NewMass().CoreObject()
	if err := s.Define("mass", mass); err != nil {
		t.Fatalf("define: %v", err)
	}
	sum, err := s.Step("add", ops.OpAdd, []core.Object{mass, mass}, ops.OperationParams{Kind: "empty"})
	if err != nil {
		t.Fatalf("step: %v", err)
	}
	if err := s.Commit(); err != nil {
		t.Fatalf("commit: %v", err)
	}
	if err := s.Conclude(sum); err != nil {
		t.Fatalf("conclude: %v", err)
	}
	cand, err := s.Seal()
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	// Metadata preservation across the seal.
	if cand.ID() != "seal_s" {
		t.Errorf("candidate ID = %q", cand.ID())
	}
	if !core.EqualObject(cand.Hypothesis(), h) {
		t.Errorf("hypothesis not preserved")
	}
	if len(cand.Premises()) != 1 || !core.EqualObject(cand.Premises()[0], mass) {
		t.Errorf("premises not preserved")
	}
	if len(cand.Predictions()) != 1 || cand.Predictions()[0].ID != "p1" {
		t.Fatalf("predictions not preserved")
	}
	if !core.EqualExpr(cand.Predictions()[0].Relation, core.NewRelation(core.RelationEq,
		mustCoreSymbol(t, "E"), mustCoreSymbol(t, "E"))) {
		t.Errorf("prediction relation not preserved")
	}
	if len(cand.FalsificationConditions()) != 1 || cand.FalsificationConditions()[0].ID != "f1" {
		t.Errorf("falsification conditions not preserved")
	}
	if len(cand.RecoveryClaims()) != 1 || cand.RecoveryClaims()[0].ID != "r1" {
		t.Errorf("recovery claims not preserved")
	}
	if len(cand.AnomalyReferences()) != 1 || cand.AnomalyReferences()[0].ID != "no_gravity" {
		t.Errorf("anomaly references not preserved")
	}
	if len(cand.ReviewHistory()) != 1 || len(cand.ReviewHistory()[0].Challenges) != 1 {
		t.Fatalf("review history not preserved")
	}
	if len(cand.FrameworkDependencies()) != 1 || cand.FrameworkDependencies()[0].FrameworkID != "classical_mechanics" {
		t.Errorf("framework dependencies not preserved")
	}
	if cand.LedgerHash() != s.Ledger().DerivationHash() {
		t.Errorf("ledger hash %s != derivation hash %s", cand.LedgerHash(), s.Ledger().DerivationHash())
	}
	if err := cand.Validate(); err != nil {
		t.Errorf("sealed candidate Validate: %v", err)
	}
	// External round-trip through the unverified loader.
	raw, err := cand.CanonicalJSON()
	if err != nil {
		t.Fatalf("candidate canonical: %v", err)
	}
	raw2, err := cand.CanonicalJSON()
	if err != nil || !bytes.Equal(raw, raw2) {
		t.Fatalf("candidate canonical JSON not deterministic")
	}
	uc, err := session.ParseResearchCandidateJSON(raw)
	if err != nil {
		t.Fatalf("parse candidate: %v", err)
	}
	rt, err := uc.Validate()
	if err != nil {
		t.Fatalf("unverified Validate: %v", err)
	}
	if rt.ID() != cand.ID() || rt.LedgerHash() != cand.LedgerHash() {
		t.Errorf("round-trip candidate mismatch")
	}
	if len(rt.Predictions()) != 1 || len(rt.ReviewHistory()) != 1 {
		t.Errorf("round-trip metadata loss")
	}
	// Post-seal mutation rejection.
	var perr core.ProvenanceError
	if _, err := s.Step("x", ops.OpAdd, []core.Object{mass, mass}, ops.OperationParams{Kind: "empty"}); !errors.As(err, &perr) {
		t.Errorf("post-seal Step = %v, want ProvenanceError", err)
	}
	if _, err := s.Seal(); !errors.As(err, &perr) {
		t.Errorf("second Seal = %v, want ProvenanceError", err)
	}
}

func mustCoreSymbol(t *testing.T, name string) core.Expr {
	t.Helper()
	e, err := core.NewSymbol(name)
	if err != nil {
		t.Fatalf("symbol: %v", err)
	}
	return e
}

func TestSealValidatesFirst(t *testing.T) {
	// Conclude verifies the conclusion object against the retained hash.
	s, sum := simpleCommittedSession(t, "seal_val")
	setSessionStepString(t, s, 1, "OutputHash",
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	var lerr core.LedgerValidationError
	if err := s.Conclude(sum); !errors.As(err, &lerr) {
		t.Fatalf("conclude over hash-tampered ledger = %v, want LedgerValidationError", err)
	}
	// Seal runs full validation: conclude cleanly, then corrupt the retained
	// output bytes, then Seal must fail (kills a Seal-skips-Validate mutant).
	s2, sum2 := simpleCommittedSession(t, "seal_val2")
	if err := s2.Conclude(sum2); err != nil {
		t.Fatalf("conclude: %v", err)
	}
	steps := s2.Ledger().Steps()
	raw := append([]byte(nil), steps[1].OutputCanonicalBytes()...)
	raw[len(raw)-2] ^= 0x01
	setSessionStepBytes(t, s2, 1, "OutputCanonical", raw)
	if _, err := s2.Seal(); !errors.As(err, &lerr) {
		t.Fatalf("seal over tampered ledger = %v, want LedgerValidationError", err)
	}
	// Params tamper passes structural + containment validation (params shape
	// is well-formed, provenance law holds) so only Session.Validate stage
	// 10 rejects it: Seal must invoke the full session pipeline, not merely
	// the candidate checks. Kills a Seal-skips-Validate mutant.
	s3, sum3 := simpleCommittedSession(t, "seal_val3")
	if err := s3.Conclude(sum3); err != nil {
		t.Fatalf("conclude: %v", err)
	}
	badParams := []byte(`{"kind":"compare","exponent":"","operator":"eq","justification":""}`)
	mirror := ledgerMirror(t, s3)
	m := mirror[1]
	m.ParamsCanonical = json.RawMessage(badParams)
	m.CurrentStepHash = mirrorBodyHash(m)
	setSessionStepBytes(t, s3, 1, "ParamsCanonical", badParams)
	setSessionStepString(t, s3, 1, "CurrentStepHash", m.CurrentStepHash)
	if _, err := s3.Seal(); !errors.As(err, &lerr) {
		t.Fatalf("seal over params-tampered ledger = %v, want LedgerValidationError", err)
	}
}

func TestSealRequiresConcludedAndHypothesis(t *testing.T) {
	// Seal without Conclude.
	s, _ := simpleCommittedSession(t, "seal_neg1")
	var perr core.ProvenanceError
	if _, err := s.Seal(); !errors.As(err, &perr) {
		t.Errorf("pre-conclude Seal = %v, want ProvenanceError", err)
	}
	// Seal with non-HYPOTHESIS hypothesis stays Concluded.
	s2 := session.New()
	meta := sealFixtureMetadata(t, "seal_neg2", mechanics.NewMass().CoreObject())
	if err := s2.Draft("seal_neg2", "neg", meta); err != nil {
		t.Fatalf("draft: %v", err)
	}
	mass := mechanics.NewMass().CoreObject()
	if err := s2.Define("mass", mass); err != nil {
		t.Fatalf("define: %v", err)
	}
	sum, err := s2.Step("add", ops.OpAdd, []core.Object{mass, mass}, ops.OperationParams{Kind: "empty"})
	if err != nil {
		t.Fatalf("step: %v", err)
	}
	if err := s2.Commit(); err != nil {
		t.Fatalf("commit: %v", err)
	}
	if err := s2.Conclude(sum); err != nil {
		t.Fatalf("conclude: %v", err)
	}
	if _, err := s2.Seal(); !errors.As(err, &perr) {
		t.Errorf("non-hypothesis Seal = %v, want ProvenanceError", err)
	}
	// State stayed Concluded: Validate still passes and Seal still reaches
	// the hypothesis check (not the state check).
	if err := s2.Validate(); err != nil {
		t.Errorf("Validate after failed Seal: %v", err)
	}
	if _, err := s2.Seal(); !errors.As(err, &perr) {
		t.Errorf("repeat Seal = %v, want ProvenanceError", err)
	}
	// Conclude with the wrong hash.
	s3, _ := simpleCommittedSession(t, "seal_neg3")
	var lerr core.LedgerValidationError
	if err := s3.Conclude(mechanics.NewTime().CoreObject()); !errors.As(err, &lerr) {
		t.Errorf("wrong-hash Conclude = %v, want LedgerValidationError", err)
	}
}

// ---------------------------------------------------------------------------
// Containment: trusted presentation of hypothesis material is rejected;
// errors.As sees the concrete typed error (F-009 probe).
// ---------------------------------------------------------------------------

func TestCandidateContainmentBasic(t *testing.T) {
	h := hypothesisFixture(t, "contain")
	s := session.New()
	if err := s.Draft("contain", "containment", sealFixtureMetadata(t, "contain", h)); err != nil {
		t.Fatalf("draft: %v", err)
	}
	mass := mechanics.NewMass().CoreObject()
	if err := s.Define("mass", mass); err != nil {
		t.Fatalf("define: %v", err)
	}
	sum, err := s.Step("add", ops.OpAdd, []core.Object{mass, mass}, ops.OperationParams{Kind: "empty"})
	if err != nil {
		t.Fatalf("step: %v", err)
	}
	if err := s.Commit(); err != nil {
		t.Fatalf("commit: %v", err)
	}
	if err := s.Conclude(sum); err != nil {
		t.Fatalf("conclude: %v", err)
	}
	cand, err := s.Seal()
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	raw, err := cand.CanonicalJSON()
	if err != nil {
		t.Fatalf("canonical: %v", err)
	}
	// Swap the HYPOTHESIS hypothesis bytes for a DEFINED object: derivation
	// replay still passes, but containment must fire.
	massRaw, err := kernel.CanonicalObjectJSON(mass)
	if err != nil {
		t.Fatalf("mass canonical: %v", err)
	}
	var wire map[string]json.RawMessage
	if err := json.Unmarshal(raw, &wire); err != nil {
		t.Fatalf("wire decode: %v", err)
	}
	wire["hypothesis"] = json.RawMessage(massRaw)
	swapped, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("wire encode: %v", err)
	}
	uc, err := session.ParseResearchCandidateJSON(swapped)
	if err != nil {
		t.Fatalf("parse swapped: %v", err)
	}
	_, err = uc.Validate()
	if err == nil {
		t.Fatal("swapped-hypothesis candidate must fail validation")
	}
	var target core.CandidateContainmentError
	if !errors.As(err, &target) {
		t.Fatalf("swapped-hypothesis error = %T (%v), want CandidateContainmentError via errors.As", err, err)
	}
	// A Validate that always succeeds would miss this; the probe above kills it.
}

func TestValidateMutation_EmptyLedger(t *testing.T) {
	l, err := session.ParseLedgerJSON([]byte(`{"schema_version":"1","steps":[],"final_hash":"","sealed":false}`))
	if err != nil {
		t.Fatalf("ParseLedgerJSON: %v", err)
	}
	var target core.LedgerValidationError
	if err := l.Validate(); !errors.As(err, &target) {
		t.Errorf("empty ledger Validate = %v, want LedgerValidationError", err)
	}
}

// TestValidateCandidateContainmentStage proves Session.Validate stage 17 is
// indispensable (M12). Primary fixture is corpus-status containment only: the
// HYPOTHESIS dependency stays HYPOTHESIS (stage 15 has nothing to reject)
// while the retained output carries CorpusStatus != NONE. Stages 1–16 pass by
// construction; only stage 17 may fire, with CandidateContainmentError.
func TestValidateCandidateContainmentStage(t *testing.T) {
	s := mustDraft(t, "stage17")
	h := hypothesisFixture(t, "stage17_h")
	if err := s.Declare("h", h); err != nil {
		t.Fatalf("declare: %v", err)
	}
	if err := s.Commit(); err != nil {
		t.Fatalf("commit: %v", err)
	}
	mirror := ledgerMirror(t, s)
	if len(mirror) != 1 {
		t.Fatalf("steps = %d, want 1", len(mirror))
	}
	m := mirror[0]
	upgrade := func(raw json.RawMessage) json.RawMessage {
		if !bytes.Contains(raw, []byte(`"corpus_status":"none"`)) {
			t.Fatalf("corpus anchor missing in %s", raw)
		}
		out := bytes.Replace(raw, []byte(`"corpus_status":"none"`), []byte(`"corpus_status":"established"`), 1)
		if _, err := kernel.LoadObjectJSON(out); err != nil {
			t.Fatalf("upgraded material must still decode: %v", err)
		}
		return out
	}
	inUp := upgrade(m.InputCanonicals[0])
	outUp := upgrade(m.OutputCanonical)
	inObj, err := kernel.LoadObjectJSON(inUp)
	if err != nil {
		t.Fatalf("decode upgraded input: %v", err)
	}
	outObj, err := kernel.LoadObjectJSON(outUp)
	if err != nil {
		t.Fatalf("decode upgraded output: %v", err)
	}
	if inObj.Provenance().Status() != core.StatusHypothesis || outObj.Provenance().Status() != core.StatusHypothesis {
		t.Fatalf("fixture must preserve HYPOTHESIS propagation (in=%v out=%v)",
			inObj.Provenance().Status(), outObj.Provenance().Status())
	}
	m.InputCanonicals[0] = inUp
	m.InputHashes[0] = lowerHex(kernel.HashObject(inObj))
	m.OutputCanonical = outUp
	m.OutputHash = lowerHex(kernel.HashObject(outObj))
	m.CurrentStepHash = mirrorBodyHash(m)
	stepsVal := sessionLedgerSteps(t, s)
	st := reflect.NewAt(stepsVal.Index(0).Type(), unsafe.Pointer(stepsVal.Index(0).UnsafeAddr())).Elem()
	fh := st.FieldByName("InputCanonicals")
	elem := reflect.NewAt(fh.Type().Elem(), unsafe.Pointer(fh.Index(0).UnsafeAddr())).Elem()
	elem.SetBytes(inUp)
	ih := st.FieldByName("InputHashes")
	ihElem := reflect.NewAt(ih.Type().Elem(), unsafe.Pointer(ih.Index(0).UnsafeAddr())).Elem()
	ihElem.SetString(lowerHex(kernel.HashObject(inObj)))
	setSessionStepBytes(t, s, 0, "OutputCanonical", outUp)
	setSessionStepString(t, s, 0, "OutputHash", lowerHex(kernel.HashObject(outObj)))
	setSessionStepString(t, s, 0, "CurrentStepHash", m.CurrentStepHash)
	// Ledger-level structural validation must NOT produce the containment
	// verdict: this isolates stage 17 of Session.Validate.
	if err := s.Ledger().Validate(); err != nil {
		t.Fatalf("Ledger.Validate must pass structurally, got: %v", err)
	}
	err = s.Validate()
	if err == nil {
		t.Fatal("corpus-status containment must fail Session.Validate")
	}
	var target core.CandidateContainmentError
	if !errors.As(err, &target) {
		t.Fatalf("stage-17 error = %T (%v), want CandidateContainmentError via errors.As", err, err)
	}
	if !strings.Contains(err.Error(), "corpus status") {
		t.Errorf("stage-17 failure attributed to wrong invariant: %v", err)
	}
}

// TestValidateOperationParamsStage proves Session.Validate stage 10 is
// indispensable (M13). Fixture: a valid `add` step whose ParamsCanonical is
// replaced with a canonically valid but operation-incompatible `pow` params
// object; chain hashes recomputed so stages 1–9 pass. `ops.Apply`
// re-validates params at replay, so without stage 10 the failure would
// surface at stage 12 with a different message — the stage-10 message pin is
// what isolates stage 10.
func TestValidateOperationParamsStage(t *testing.T) {
	s, _ := simpleCommittedSession(t, "stage10")
	powParams := []byte(`{"kind":"pow","exponent":"2/1","operator":"","justification":""}`)
	if _, err := ops.ParseOperationParams(powParams); err != nil {
		t.Fatalf("fixture params must be shape-valid: %v", err)
	}
	mirror := ledgerMirror(t, s)
	m := mirror[1]
	m.ParamsCanonical = json.RawMessage(powParams)
	m.CurrentStepHash = mirrorBodyHash(m)
	setSessionStepBytes(t, s, 1, "ParamsCanonical", powParams)
	setSessionStepString(t, s, 1, "CurrentStepHash", m.CurrentStepHash)
	// Ledger-level structural validation passes shape-valid params: only
	// session stage 10 judges operation binding.
	if err := s.Ledger().Validate(); err != nil {
		t.Fatalf("Ledger.Validate must pass shape-valid params, got: %v", err)
	}
	var target core.LedgerValidationError
	err := s.Validate()
	if !errors.As(err, &target) {
		t.Fatalf("params-binding tamper Validate = %v, want LedgerValidationError", err)
	}
	if !strings.Contains(err.Error(), "operation params:") {
		t.Errorf("params-binding failure attributed to wrong stage (want stage 10): %v", err)
	}
}

// TestSessionStepNoPartialCommit (Gate A): a failed Session.Step must leave
// the draft and ledger unchanged — Commit afterwards must behave as if the
// failed call never happened.
func TestSessionStepNoPartialCommit(t *testing.T) {
	s := mustDraft(t, "nopartial")
	mass := mechanics.NewMass().CoreObject()
	if _, err := s.Step("bad", ops.OpAdd, []core.Object{mass}, ops.OperationParams{Kind: "empty"}); err == nil {
		t.Fatal("wrong-arity Step must fail")
	}
	if _, err := s.Identify(mass, mass, "   "); err == nil {
		t.Fatal("empty-justification Identify must fail")
	}
	// Draft must be empty: Commit must report empty-draft, not commit one step.
	var target core.LedgerValidationError
	if err := s.Commit(); !errors.As(err, &target) {
		t.Fatalf("Commit after failures = %v, want LedgerValidationError (empty draft)", err)
	}
	// The session remains usable.
	if err := s.Define("mass", mass); err != nil {
		t.Fatalf("define: %v", err)
	}
	if _, err := s.Step("add", ops.OpAdd, []core.Object{mass, mass}, ops.OperationParams{Kind: "empty"}); err != nil {
		t.Fatalf("step: %v", err)
	}
	if err := s.Commit(); err != nil {
		t.Fatalf("commit: %v", err)
	}
	if got := len(s.Ledger().Steps()); got != 2 {
		t.Fatalf("ledger steps = %d, want 2 (no partial append)", got)
	}
}

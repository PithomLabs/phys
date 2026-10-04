package relativity

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

func TestRelativityManifestCrossCheck(t *testing.T) {
	m, err := core.ParseManifest(ManifestJSON)
	if err != nil {
		t.Fatalf("ParseManifest: %v", err)
	}
	if m.FrameworkID != "special_relativity" {
		t.Errorf("framework_id = %q", m.FrameworkID)
	}
	if m.CorpusStatus != core.CorpusEstablished {
		t.Errorf("corpus_status = %q", m.CorpusStatus)
	}
	// F1: exactly 10 items (Velocity and zeros NOT manifest items)
	if len(m.Items) != 10 {
		t.Errorf("expected 10 items, got %d", len(m.Items))
	}
	if len(m.Domain) != 8 {
		t.Errorf("expected 8 domains, got %d", len(m.Domain))
	}
}

func TestManifestCanonicalRoundTrip(t *testing.T) {
	m, err := core.ParseManifest(ManifestJSON)
	if err != nil {
		t.Fatalf("ParseManifest: %v", err)
	}
	canon, err := core.CanonicalManifestJSON(m)
	if err != nil {
		t.Fatalf("CanonicalManifestJSON: %v", err)
	}
	m2, err := core.ParseManifest(canon)
	if err != nil {
		t.Fatalf("re-parse canonical: %v", err)
	}
	if m.Hash() != m2.Hash() {
		t.Error("hash mismatch after round-trip")
	}
}

func TestManifestConstructorCrossCheck(t *testing.T) {
	constructors := manifestConstructors()

	m, err := core.ParseManifest(ManifestJSON)
	if err != nil {
		t.Fatalf("ParseManifest: %v", err)
	}

	for _, item := range m.Items {
		ctor, ok := constructors[item.Constructor]
		if !ok {
			t.Errorf("item %s: unknown constructor %q", item.ID, item.Constructor)
			continue
		}
		obj := ctor()
		if !core.EqualExpr(item.CanonicalExpr, obj.Expr()) {
			t.Errorf("item %s: canonical_expr mismatch", item.ID)
		}
		if !item.Dimension.Equal(obj.Dimension()) {
			t.Errorf("item %s: dimension mismatch: manifest %v vs obj %v", item.ID, item.Dimension, obj.Dimension())
		}
		if item.Kind != obj.Kind().String() {
			t.Errorf("item %s: kind mismatch: manifest %q vs obj %q", item.ID, item.Kind, obj.Kind().String())
		}
		if string(item.ProvenanceStatus) != string(obj.Provenance().Status()) {
			t.Errorf("item %s: provenance_status mismatch", item.ID)
		}
		if item.Source != "" && obj.Provenance().Source() != item.Source {
			t.Errorf("item %s: source mismatch", item.ID)
		}
		if obj.CorpusStatus() != m.CorpusStatus {
			t.Errorf("item %s: corpus status %v != framework %v", item.ID, obj.CorpusStatus(), m.CorpusStatus)
		}
		if want := core.NewAssumptionSet(item.Assumptions...); !want.Equal(obj.Assumptions()) {
			t.Errorf("item %s: assumptions mismatch", item.ID)
		}
	}
}

// TestManifestEmbedPinned pins the manifest loading mechanism: the manifest
// must be embedded at compile time via go:embed, never read from disk at
// runtime (replacing embed with a file read breaks hermetic builds).
func TestManifestEmbedPinned(t *testing.T) {
	src, err := readSourceFile("primitives.go")
	if err != nil {
		t.Fatalf("read primitives.go: %v", err)
	}
	if !containsSubstr(src, "//go:embed manifest.json") {
		t.Errorf("primitives.go must embed the manifest via //go:embed")
	}
	if containsSubstr(src, "os.ReadFile") || containsSubstr(src, "os.Open") {
		t.Errorf("primitives.go must not read the manifest from disk at runtime")
	}
}

func readSourceFile(name string) (string, error) {
	return readFileAtCaller(name)
}

func containsSubstr(haystack, needle string) bool {
	if len(needle) == 0 {
		return true
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

func readFileAtCaller(name string) (string, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("caller unavailable")
	}
	raw, err := os.ReadFile(filepath.Join(filepath.Dir(file), name))
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// TestReverseConstructorAllowlist (Gate C): closed exported constructor set.
func TestReverseConstructorAllowlist(t *testing.T) {
	allowed := map[string]bool{
		"NewSpacetime": true, "NewMinkowskiMetric": true, "NewRestMass": true,
		"NewEnergy": true, "NewThreeMomentum": true, "NewFourMomentum": true,
		"NewSpeedOfLight": true, "LorentzFactor": true, "EnergyMomentumRelation": true,
		"MassEnergyRelation": true, "NewVelocity": true,
		"ZeroThreeMomentum": true, "ZeroEnergy": true, "ZeroVelocity": true,
	}
	// RestFrameAssumption is a metadata helper, not an object constructor
	// (returns Assumption); pinned separately per G5 refinement.
	allowedHelpers := map[string]bool{"RestFrameAssumption": true}
	names := exportedConstructors(t, []string{"primitives.go", "relations.go"})
	if len(names) != len(allowed)+len(allowedHelpers) {
		t.Fatalf("exported constructors = %v, want exactly %d pinned names", names, len(allowed)+len(allowedHelpers))
	}
	for _, n := range names {
		if !allowed[n] && !allowedHelpers[n] {
			t.Errorf("unapproved constructor %q (not manifest-backed, not allowlisted)", n)
		}
	}
}

// exportedConstructors parses the given same-package source files and returns
// exported receiverless function names (constructors), sorted.
func exportedConstructors(t *testing.T, files []string) []string {
	t.Helper()
	fset := token.NewFileSet()
	var names []string
	for _, f := range files {
		src, err := readSourceFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		pf, err := parser.ParseFile(fset, f, src, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", f, err)
		}
		for _, decl := range pf.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Name.IsExported() {
				names = append(names, fn.Name.Name)
			}
		}
	}
	sort.Strings(names)
	return names
}

func manifestConstructors() map[string]func() kernel.Object {
	return map[string]func() kernel.Object{
		"relativity.NewSpacetime":           func() kernel.Object { return NewSpacetime().CoreObject() },
		"relativity.NewMinkowskiMetric":     func() kernel.Object { return NewMinkowskiMetric().CoreObject() },
		"relativity.NewRestMass":            func() kernel.Object { return NewRestMass().CoreObject() },
		"relativity.NewEnergy":              func() kernel.Object { return NewEnergy().CoreObject() },
		"relativity.NewThreeMomentum":       func() kernel.Object { return NewThreeMomentum().CoreObject() },
		"relativity.NewFourMomentum":        func() kernel.Object { return NewFourMomentum().CoreObject() },
		"relativity.NewSpeedOfLight":        func() kernel.Object { return NewSpeedOfLight().CoreObject() },
		"relativity.LorentzFactor":          func() kernel.Object { return LorentzFactor() },
		"relativity.EnergyMomentumRelation": func() kernel.Object { return EnergyMomentumRelation() },
		"relativity.MassEnergyRelation":     func() kernel.Object { return MassEnergyRelation() },
	}
}

// TestManifestCorruptionBattery (Gate F): each independent in-memory
// corruption must be rejected by schema validation or constructor
// cross-check. Never touches production files or disk state.
func TestManifestCorruptionBattery(t *testing.T) {
	ctors := manifestConstructors()
	withItems := func(mut func(items []any)) []byte {
		var m map[string]any
		if err := json.Unmarshal(ManifestJSON, &m); err != nil {
			t.Fatalf("decode: %v", err)
		}
		mut(m["items"].([]any))
		raw, err := json.Marshal(m)
		if err != nil {
			t.Fatalf("encode: %v", err)
		}
		return raw
	}
	withRoot := func(mut func(m map[string]any)) []byte {
		var m map[string]any
		if err := json.Unmarshal(ManifestJSON, &m); err != nil {
			t.Fatalf("decode: %v", err)
		}
		mut(m)
		raw, err := json.Marshal(m)
		if err != nil {
			t.Fatalf("encode: %v", err)
		}
		return raw
	}
	// crossMismatch replays the constructor cross-check; true = rejected.
	crossMismatch := func(raw []byte) bool {
		t.Helper()
		m, err := core.ParseManifest(raw)
		if err != nil {
			return true
		}
		if m.CorpusStatus != core.CorpusEstablished {
			return true
		}
		for _, item := range m.Items {
			ctor, ok := ctors[item.Constructor]
			if !ok {
				return true
			}
			obj := ctor()
			if !core.EqualExpr(item.CanonicalExpr, obj.Expr()) {
				return true
			}
			if !item.Dimension.Equal(obj.Dimension()) {
				return true
			}
			if item.Kind != obj.Kind().String() {
				return true
			}
			if string(item.ProvenanceStatus) != string(obj.Provenance().Status()) {
				return true
			}
			if obj.CorpusStatus() != m.CorpusStatus {
				return true
			}
			if want := core.NewAssumptionSet(item.Assumptions...); !want.Equal(obj.Assumptions()) {
				return true
			}
			if item.Source != "" && obj.Provenance().Source() != item.Source {
				return true
			}
		}
		return false
	}
	rejected := func(name string, raw []byte) {
		t.Helper()
		if _, err := core.ValidateManifestBytes(raw); err != nil {
			return
		}
		if crossMismatch(raw) {
			return
		}
		t.Errorf("%s: corruption undetected by schema and cross-check", name)
	}
	first := func(items []any) map[string]any { return items[0].(map[string]any) }
	rejected("canonical_expr", withItems(func(items []any) {
		first(items)["canonical_expr"] = map[string]any{"kind": "symbol", "name": "zzz"}
	}))
	rejected("dimension", withItems(func(items []any) {
		first(items)["dimension"] = map[string]any{"m": "2/1", "l": "0/1", "t": "0/1", "i": "0/1", "theta": "0/1", "n": "0/1", "j": "0/1"}
	}))
	rejected("kind", withItems(func(items []any) {
		first(items)["kind"] = "Energy"
	}))
	rejected("provenance_status", withItems(func(items []any) {
		first(items)["provenance_status"] = "HYPOTHESIS"
	}))
	rejected("corpus_status_flip", withRoot(func(m map[string]any) {
		m["corpus_status"] = "contested"
	}))
	rejected("assumptions", withItems(func(items []any) {
		first(items)["assumptions"] = []any{
			map[string]any{"kind": "constraint", "key": "injected", "value": map[string]any{"mode": "text", "value": "x"}},
		}
	}))
	rejected("source", withItems(func(items []any) {
		first(items)["source"] = "Bogus Source"
	}))
	rejected("unknown_constructor", withItems(func(items []any) {
		first(items)["constructor"] = "pkg.NoSuchConstructor"
	}))
	rejected("extra_item", withRoot(func(m map[string]any) {
		items := m["items"].([]any)
		dup := map[string]any{}
		for k, v := range first(items) {
			dup[k] = v
		}
		dup["id"] = "injected_item"
		dup["constructor"] = "pkg.NoSuchConstructor"
		m["items"] = append(items, dup)
	}))
}

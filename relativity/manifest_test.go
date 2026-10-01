package relativity

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
	constructors := map[string]func() kernel.Object{
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

package mechanics

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

func TestMechanicsManifestCrossCheck(t *testing.T) {
	m, err := core.ParseManifest(ManifestJSON)
	if err != nil {
		t.Fatalf("ParseManifest: %v", err)
	}
	if m.FrameworkID != "classical_mechanics" {
		t.Errorf("framework_id = %q", m.FrameworkID)
	}
	if m.CorpusStatus != core.CorpusEstablished {
		t.Errorf("corpus_status = %q", m.CorpusStatus)
	}
	if len(m.Items) != 11 {
		t.Errorf("expected 11 items, got %d", len(m.Items))
	}
	if len(m.Domain) != 9 {
		t.Errorf("expected 10 domains, got %d", len(m.Domain))
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
		"mechanics.NewMass":               func() kernel.Object { return NewMass().CoreObject() },
		"mechanics.NewTime":               func() kernel.Object { return NewTime().CoreObject() },
		"mechanics.NewPosition":           func() kernel.Object { return NewPosition().CoreObject() },
		"mechanics.NewVelocity":           func() kernel.Object { return NewVelocity().CoreObject() },
		"mechanics.NewAcceleration":       func() kernel.Object { return NewAcceleration().CoreObject() },
		"mechanics.NewForce":              func() kernel.Object { return NewForce().CoreObject() },
		"mechanics.NewMomentum":           func() kernel.Object { return NewMomentum().CoreObject() },
		"mechanics.NewEnergy":             func() kernel.Object { return NewEnergy().CoreObject() },
		"mechanics.NewtonSecondLaw":       func() kernel.Object { return NewtonSecondLaw() },
		"mechanics.MomentumRelation":      func() kernel.Object { return MomentumRelation() },
		"mechanics.KineticEnergyRelation": func() kernel.Object { return KineticEnergyRelation() },
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

package core_test

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/PithomLabs/phys/core"
	"github.com/PithomLabs/phys/internal/kernel"
)

// ---------------------------------------------------------------------------
// Object authority (§5)
// ---------------------------------------------------------------------------

func TestZeroObjectInvalid(t *testing.T) {
	var z core.Object
	if z.Valid() {
		t.Fatal("zero core.Object must be invalid")
	}
	if z.Name() != "" || z.Kind() != core.KindMass || z.CorpusStatus() != "" {
		t.Fatal("zero accessors")
	}
	if z.Dimension().Valid() || z.Expr().Valid() {
		t.Fatal("zero dimension/expr must be invalid")
	}
	if len(z.Assumptions().Values()) != 0 || len(z.Conventions().Values()) != 0 {
		t.Fatal("zero metadata sets empty")
	}
	// An invalid object has no canonical form and never equals a minted one.
	if core.EqualObject(z, z) {
		t.Fatal("invalid objects have no canonical representation")
	}
}

func TestObjectAccessorSurfaceExact(t *testing.T) {
	// REQ-005-08/09: the Object method set is exactly the nine §5.3
	// read-only accessors, all with value receivers (no mutators).
	fset := token.NewFileSet()
	dir, err := parser.ParseDir(fset, filepath.Join("..", "internal", "kernel"),
		func(fi os.FileInfo) bool { return !strings.HasSuffix(fi.Name(), "_test.go") }, 0)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"Valid":        "bool",
		"Name":         "string",
		"Kind":         "Kind",
		"Dimension":    "Dimension",
		"Expr":         "Expr",
		"Assumptions":  "AssumptionSet",
		"Conventions":  "ConventionSet",
		"Provenance":   "Provenance",
		"CorpusStatus": "CorpusStatus",
	}
	got := map[string]string{}
	for _, pkg := range dir {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Recv == nil || len(fn.Recv.List) != 1 {
					continue
				}
				recv := fn.Recv.List[0].Type
				star, isStar := recv.(*ast.StarExpr)
				if isStar {
					recv = star.X
				}
				id, ok := recv.(*ast.Ident)
				if !ok || id.Name != "Object" {
					continue
				}
				if isStar {
					t.Fatalf("Object method %s uses a pointer receiver", fn.Name.Name)
				}
				if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
					t.Fatalf("Object method %s must have exactly one result", fn.Name.Name)
				}
				res, ok := fn.Type.Results.List[0].Type.(*ast.Ident)
				if !ok {
					t.Fatalf("Object method %s result type", fn.Name.Name)
				}
				if _, exists := got[fn.Name.Name]; exists {
					t.Fatalf("duplicate Object method %s", fn.Name.Name)
				}
				got[fn.Name.Name] = res.Name
			}
		}
	}
	if len(got) != len(want) {
		t.Fatalf("Object method set = %v, want exactly %v", got, want)
	}
	for name, res := range want {
		if got[name] != res {
			t.Fatalf("Object.%s result = %q, want %q", name, got[name], res)
		}
	}
}

func TestKernelObjectFieldsUnexported(t *testing.T) {
	rt := reflect.TypeOf(core.Object{})
	want := []string{
		"valid", "name", "kind", "dimension", "expr",
		"assumptions", "conventions", "provenance", "corpusStatus",
	}
	if rt.NumField() != len(want) {
		t.Fatalf("Object fields = %d, want %d", rt.NumField(), len(want))
	}
	for i, name := range want {
		f := rt.Field(i)
		if f.Name != name {
			t.Fatalf("field %d = %q, want %q", i, f.Name, name)
		}
		if f.PkgPath == "" {
			t.Fatalf("field %s is exported", f.Name)
		}
	}
}

// validMintSpec builds a fully valid mint spec for negative-case derivation.
func validMintSpec(t *testing.T) kernel.ObjectSpec {
	t.Helper()
	prov, err := kernel.NewProvenance(kernel.StatusDefined, "fixture", "classical_mechanics",
		nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	if err != nil {
		t.Fatal(err)
	}
	m, err := kernel.NewSymbol("m")
	if err != nil {
		t.Fatal(err)
	}
	return kernel.ObjectSpec{
		Name:         "mass",
		Kind:         kernel.KindMass,
		Dimension:    core.DimensionMass(),
		Expr:         m,
		Assumptions:  kernel.NewAssumptionSet(),
		Conventions:  kernel.NewConventionSet(),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusEstablished,
	}
}

func TestMintObjectRejectsInvalidSpec(t *testing.T) {
	// One failing case per §5.2.1 contract point.
	badExpr := validMintSpec(t)
	badExpr.Expr = core.Expr{}
	if _, err := kernel.MintObject(badExpr); !errors.As(err, new(core.InvalidObjectError)) {
		t.Fatalf("invalid expr: %v", err)
	}

	emptyName := validMintSpec(t)
	emptyName.Name = ""
	if _, err := kernel.MintObject(emptyName); !errors.As(err, new(core.InvalidObjectError)) {
		t.Fatalf("empty required name: %v", err)
	}

	badDim := validMintSpec(t)
	badDim.Dimension = core.Dimension{}
	if _, err := kernel.MintObject(badDim); !errors.As(err, new(core.InvalidObjectError)) {
		t.Fatalf("invalid dimension: %v", err)
	}

	badProv := validMintSpec(t)
	badProv.Provenance = core.Provenance{}
	if _, err := kernel.MintObject(badProv); !errors.As(err, new(core.ProvenanceError)) {
		t.Fatalf("invalid provenance: %v", err)
	}

	badKind := validMintSpec(t)
	badKind.Kind = core.Kind(200)
	if _, err := kernel.MintObject(badKind); !errors.As(err, new(core.InvalidObjectError)) {
		t.Fatalf("invalid kind: %v", err)
	}

	badCorpus := validMintSpec(t)
	badCorpus.CorpusStatus = ""
	if _, err := kernel.MintObject(badCorpus); !errors.As(err, new(core.InvalidObjectError)) {
		t.Fatalf("invalid corpus status: %v", err)
	}

	// Positive path: the valid spec mints a fully valid immutable object.
	obj, err := kernel.MintObject(validMintSpec(t))
	if err != nil {
		t.Fatal(err)
	}
	if !obj.Valid() || obj.Name() != "mass" || obj.Kind() != kernel.KindMass ||
		obj.CorpusStatus() != kernel.CorpusEstablished || !obj.Dimension().Valid() {
		t.Fatal("minted object accessors")
	}

	// Pinned name rule: operation results mint empty names.
	derived := validMintSpec(t)
	derived.Name = ""
	derived.Provenance, err = kernel.NewProvenance(kernel.StatusDerived, "", "", nil,
		[32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := kernel.MintObject(derived); err != nil {
		t.Fatalf("DERIVED operation results must mint with empty names: %v", err)
	}
	// POSTULATED also requires a name (session fixture contract).
	posted := validMintSpec(t)
	posted.Name = ""
	posted.Provenance, err = kernel.NewProvenance(kernel.StatusPostulated, "", "", nil,
		[32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := kernel.MintObject(posted); !errors.As(err, new(core.InvalidObjectError)) {
		t.Fatalf("POSTULATED requires a name: %v", err)
	}
	// HYPOTHESIS operation results may be unnamed (status propagation).
	hyp := derived
	hyp.Provenance, err = kernel.NewProvenance(kernel.StatusHypothesis, "", "", nil,
		[32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := kernel.MintObject(hyp); err != nil {
		t.Fatalf("HYPOTHESIS operation results may be unnamed: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Boundary & authority scans (§0, §2, §4, §5.2.1)
// ---------------------------------------------------------------------------

// moduleRoot returns the repository root from the core package directory.
func moduleRoot() string { return ".." }

// eachGoFile walks the module's Go sources, skipping .git and the plan
// document directories (task docs, never packages).
func eachGoFile(t *testing.T, includeTests bool, fn func(path string, src []byte)) {
	t.Helper()
	err := filepath.Walk(moduleRoot(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			base := info.Name()
			if base == ".git" || base == "plan10" || base == "plans" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if !includeTests && strings.HasSuffix(path, "_test.go") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fn(path, b)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestObjectCanonicalRoundTripRejectsNonCanonical(t *testing.T) {
	// §10.5: canonical Object JSON round-trips byte-identically through the
	// internal replay decoder, while a semantically equivalent but
	// non-canonical encoding must be rejected (decoder is not a permissive
	// JSON parser). Kills removal of the LoadObjectJSON re-encode check.
	obj, err := kernel.MintObject(validMintSpec(t))
	if err != nil {
		t.Fatalf("mint fixture: %v", err)
	}
	canonical, err := kernel.CanonicalObjectJSON(obj)
	if err != nil {
		t.Fatalf("canonical encode: %v", err)
	}
	back, err := kernel.LoadObjectJSON(canonical)
	if err != nil {
		t.Fatalf("canonical decode: %v", err)
	}
	reencoded, err := kernel.CanonicalObjectJSON(back)
	if err != nil {
		t.Fatalf("re-encode: %v", err)
	}
	if !bytes.Equal(reencoded, canonical) {
		t.Fatalf("round-trip not byte-identical:\n got %s\nwant %s", reencoded, canonical)
	}
	if kernel.HashObject(back) != kernel.HashObject(obj) {
		t.Fatal("round-trip hash mismatch")
	}
	// Same semantic content, non-canonical bytes: a single insignificant
	// space after the first top-level comma. Nested canonical values stay
	// byte-identical so every sub-component check passes — only the
	// top-level re-encode equality check may reject this.
	prefix := []byte(`{"schema_version":"1",`)
	if !bytes.HasPrefix(canonical, prefix) {
		t.Fatalf("unexpected canonical prefix in %s", canonical)
	}
	spaced := append(append([]byte(nil), canonical[:len(prefix)]...), append([]byte(" "), canonical[len(prefix):]...)...)
	if bytes.Equal(spaced, canonical) {
		t.Fatal("fixture is not actually non-canonical")
	}
	if _, err := kernel.LoadObjectJSON(spaced); err == nil {
		t.Fatal("non-canonical encoding must be rejected")
	} else if !strings.Contains(err.Error(), "non-canonical object JSON") {
		t.Fatalf("rejection attributed to wrong check (want top-level): %v", err)
	}
}

func TestNoGenericFactory(t *testing.T) {
	// REQ-005-04/10, REQ-§5.2.1-MUST-01/03, REQ-§5.4-MUST-01, REQ-032-06:
	// exactly one func MintObject in the module, under internal/kernel;
	// no NewObject anywhere; core exports neither.
	fset := token.NewFileSet()
	mintCount := 0
	newObjectCount := 0
	coreExports := map[string]bool{}
	// M15 semantic authority audit: every production function in
	// internal/kernel returning Object, and where Object values are built.
	kernelObjectEntries := map[string]bool{}
	kernelLiteralHolders := map[string]bool{}
	kernelMintCallers := map[string]bool{}
	kernelCanonCheckers := map[string]bool{}
	domainParamBanned := map[string]bool{
		"Kind": true, "CorpusStatus": true, "Provenance": true, "Dimension": true,
	}
	err := filepath.Walk(moduleRoot(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == "plan10" || info.Name() == "plans" {
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
		rel, err := filepath.Rel(moduleRoot(), path)
		if err != nil {
			return err
		}
		pkgDir := filepath.Dir(rel)
		inKernel := strings.Contains(filepath.ToSlash(pkgDir), "internal/kernel")
		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if inKernel && returnsObject(d) {
					kernelObjectEntries[d.Name.Name] = true
					ast.Inspect(d.Body, func(n ast.Node) bool {
						lit, ok := n.(*ast.CompositeLit)
						if ok {
							if id, ok := lit.Type.(*ast.Ident); ok && id.Name == "Object" && len(lit.Elts) > 0 {
								kernelLiteralHolders[d.Name.Name] = true
							}
						}
						call, ok := n.(*ast.CallExpr)
						if ok {
							switch fn := call.Fun.(type) {
							case *ast.Ident:
								if fn.Name == "MintObject" {
									kernelMintCallers[d.Name.Name] = true
								}
							case *ast.SelectorExpr:
								if x, ok := fn.X.(*ast.Ident); ok && x.Name == "bytes" && fn.Sel.Name == "Equal" {
									kernelCanonCheckers[d.Name.Name] = true
								}
							}
						}
						return true
					})
				}
				if d.Recv != nil {
					break
				}
				switch d.Name.Name {
				case "MintObject":
					mintCount++
					if f.Name.Name != "kernel" ||
						!strings.Contains(filepath.ToSlash(pkgDir), "internal/kernel") {
						t.Errorf("MintObject must live in internal/kernel, found in %s", rel)
					}
				case "NewObject":
					newObjectCount++
				}
				// REQ-005-05: domain packages (mechanics, relativity) must
				// not expose constructors taking these parameter types.
				// The hypothesis package is exempt (REQ-024-02).
				if pkgDir == "mechanics" || pkgDir == "relativity" {
					if d.Name.IsExported() && d.Type.Params != nil {
						for _, field := range d.Type.Params.List {
							var id *ast.Ident
							switch tt := field.Type.(type) {
							case *ast.Ident:
								id = tt
							case *ast.SelectorExpr:
								if x, ok := tt.X.(*ast.Ident); ok && x.Name == "core" {
									id = tt.Sel
								}
							}
							if id != nil && domainParamBanned[id.Name] {
								t.Errorf("%s exports constructor %s taking %s (REQ-005-05)",
									pkgDir, d.Name.Name, id.Name)
							}
						}
					}
				}
			case *ast.GenDecl:
				if d.Tok == token.TYPE && f.Name.Name == "core" && d.Lparen == 0 {
					// single-type declaration
				}
			}
		}
		if f.Name.Name == "core" {
			for _, decl := range f.Decls {
				switch d := decl.(type) {
				case *ast.FuncDecl:
					if d.Recv == nil && d.Name.IsExported() {
						coreExports[d.Name.Name] = true
					}
				case *ast.GenDecl:
					for _, spec := range d.Specs {
						switch s := spec.(type) {
						case *ast.TypeSpec:
							if s.Name.IsExported() {
								coreExports[s.Name.Name] = true
							}
						case *ast.ValueSpec:
							for _, n := range s.Names {
								if n.IsExported() {
									coreExports[n.Name] = true
								}
							}
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if mintCount != 1 {
		t.Fatalf("func MintObject count = %d, want exactly 1", mintCount)
	}
	if newObjectCount != 0 {
		t.Fatalf("NewObject found %d time(s)", newObjectCount)
	}
	for _, banned := range []string{"MintObject", "NewObject", "ParseObject", "DecodeObjectJSON"} {
		if coreExports[banned] {
			t.Fatalf("core exports banned symbol %s", banned)
		}
	}
	// M15 primary invariant: the exact allowed production Object entry-point
	// set is {MintObject (sole mint authority), LoadObjectJSON (allowed
	// internal decoder)}. Any additional name — e.g. MintObjectAlt — fails
	// regardless of identifier, so renaming cannot evade detection.
	for name := range kernelObjectEntries {
		if name != "MintObject" && name != "LoadObjectJSON" {
			t.Errorf("second Object entry point %q in internal/kernel", name)
		}
	}
	if !kernelObjectEntries["MintObject"] || !kernelObjectEntries["LoadObjectJSON"] {
		t.Errorf("kernel Object entry points = %v, want exactly {MintObject LoadObjectJSON}", kernelObjectEntries)
	}
	// Secondary safeguards: single construction site; decoder-only shape.
	if len(kernelLiteralHolders) != 1 || !kernelLiteralHolders["MintObject"] {
		t.Errorf("Object construction sites = %v, want exactly {MintObject}", kernelLiteralHolders)
	}
	if !kernelMintCallers["LoadObjectJSON"] || !kernelCanonCheckers["LoadObjectJSON"] {
		t.Errorf("LoadObjectJSON must decode via MintObject with a canonical-equality check")
	}
}

func returnsObject(d *ast.FuncDecl) bool {
	if d.Type.Results == nil {
		return false
	}
	for _, field := range d.Type.Results.List {
		switch tt := field.Type.(type) {
		case *ast.Ident:
			if tt.Name == "Object" {
				return true
			}
		}
	}
	return false
}

func TestCoreImportIndependence(t *testing.T) {
	// REQ-004-02, REQ-§4-MUST-01/02: core imports neither ops nor session;
	// every core import is stdlib or internal/kernel (G-Imports shape).
	fset := token.NewFileSet()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	seenKernel := false
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		f, err := parser.ParseFile(fset, e.Name(), nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, im := range f.Imports {
			path := strings.Trim(im.Path.Value, `"`)
			isTest := strings.HasSuffix(e.Name(), "_test.go")
			switch {
			case path == "github.com/PithomLabs/phys/internal/kernel":
				seenKernel = true
			case path == "github.com/PithomLabs/phys/core" && isTest:
				// external test package core_test references the package under test
			case strings.HasPrefix(path, "github.com/PithomLabs/phys/"):
				t.Errorf("core/%s imports %s", e.Name(), path)
			case strings.Contains(path, "/ops") || strings.Contains(path, "/session"):
				t.Errorf("core/%s imports %s", e.Name(), path)
			case strings.Contains(path, "."):
				// dotted first path element = third-party module
				if !strings.HasPrefix(path, "github.com/") && !isStdlib(path) {
					t.Errorf("core/%s imports non-stdlib %s", e.Name(), path)
				}
			}
		}
	}
	if !seenKernel {
		t.Fatal("core must alias the kernel types (imports internal/kernel)")
	}
}

// isStdlib reports whether an import path has no dotted first element.
func isStdlib(path string) bool {
	first := path
	if i := strings.IndexByte(path, '/'); i >= 0 {
		first = path[:i]
	}
	return !strings.Contains(first, ".")
}

func TestInternalKernelBoundary(t *testing.T) {
	// REQ-004-03: the authority lives below Go's internal/ visibility rule.
	kernelDir := filepath.Join("..", "internal", "kernel")
	if !strings.Contains(filepath.ToSlash(filepath.Clean(kernelDir)), "/internal/") {
		t.Fatal("kernel package must live under internal/")
	}
	// Public external callers cannot import internal/kernel: enforce the
	// directory shape the toolchain then applies.
	if _, err := os.Stat(filepath.Join(kernelDir, "types.go")); err != nil {
		t.Fatalf("kernel sources: %v", err)
	}
	// REQ-000-04 / REQ-§0.3-MUST-01: constructor authority is an
	// external/public API boundary statement, never an adversarial-protection
	// claim. README audit runs once the README exists (step 15).
	readme, err := os.ReadFile(filepath.Join(moduleRoot(), "README.md"))
	if err == nil {
		text := string(readme)
		for _, want := range []string{"constructor authority", "not cryptographic"} {
			if !strings.Contains(text, want) {
				t.Errorf("README must state %q", want)
			}
		}
		for _, banned := range []string{
			"adversarial protection", "tamper-proof", "tamperproof",
			"cryptographically secure", "unhackable", "secure against",
		} {
			if strings.Contains(text, banned) {
				t.Errorf("README must not claim %q", banned)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Kind enum (§6)
// ---------------------------------------------------------------------------

func TestKindEnumExact(t *testing.T) {
	want := []string{
		"Mass", "Time", "Position", "Velocity", "Acceleration", "Force",
		"Momentum", "Energy", "KineticEnergy", "RestMass", "ThreeMomentum",
		"FourMomentum", "SpeedOfLight", "Spacetime", "MinkowskiMetric",
		"Expression", "Relation", "BranchSet",
	}
	if len(want) != 18 {
		t.Fatal("exactly 18 kinds")
	}
	for i, name := range want {
		k := core.Kind(i)
		if k.String() != name {
			t.Fatalf("Kind(%d).String() = %q, want %q", i, k.String(), name)
		}
		back, ok := kernel.ParseKind(name)
		if !ok || back != k {
			t.Fatalf("ParseKind(%q) round-trip", name)
		}
	}
	// Stable ordinal anchors for mrc-v0.4 (§8.2-style pinning).
	if core.KindMass != 0 || core.KindExpression != 15 || core.KindBranchSet != 17 {
		t.Fatal("kind ordinals are stable for mrc-v0.4")
	}
	if _, ok := kernel.ParseKind("Bogus"); ok {
		t.Fatal("unknown kind name must be rejected")
	}
	if core.Kind(18).String() != "" {
		t.Fatal("out-of-range kind has empty name")
	}
}

// ---------------------------------------------------------------------------
// Object equality, copies, tree & source hygiene
// ---------------------------------------------------------------------------

func TestEqualObjectMetadata(t *testing.T) {
	// REQ-009-01/02: canonical object equality covers all metadata.
	base, err := kernel.MintObject(validMintSpec(t))
	if err != nil {
		t.Fatal(err)
	}
	same, err := kernel.MintObject(validMintSpec(t))
	if err != nil {
		t.Fatal(err)
	}
	if !core.EqualObject(base, same) || core.HashObject(base) != core.HashObject(same) {
		t.Fatal("identical specs must produce equal objects with equal hashes")
	}

	changes := func(mutate func(*kernel.ObjectSpec)) bool {
		spec := validMintSpec(t)
		mutate(&spec)
		o, err := kernel.MintObject(spec)
		if err != nil {
			t.Fatalf("mutated spec: %v", err)
		}
		return core.EqualObject(base, o)
	}
	if changes(func(s *kernel.ObjectSpec) { s.Name = "other" }) {
		t.Fatal("name difference must be visible to EqualObject")
	}
	if changes(func(s *kernel.ObjectSpec) { s.CorpusStatus = kernel.CorpusContested }) {
		t.Fatal("corpus-status difference must be visible")
	}
	if changes(func(s *kernel.ObjectSpec) { s.Kind = kernel.KindRestMass }) {
		t.Fatal("kind difference must be visible")
	}
	if changes(func(s *kernel.ObjectSpec) {
		s.Provenance, _ = kernel.NewProvenance(kernel.StatusDefined, "other",
			"classical_mechanics", nil, [32]byte{}, [32]byte{}, kernel.MRCVersion, "")
	}) {
		t.Fatal("provenance difference must be visible")
	}
	if changes(func(s *kernel.ObjectSpec) {
		a, _ := kernel.NewTextAssumption(kernel.AssumptionDomain, "regime", "classical")
		s.Assumptions = kernel.NewAssumptionSet(a)
	}) {
		t.Fatal("assumption difference must be visible")
	}
	if changes(func(s *kernel.ObjectSpec) {
		c, _ := kernel.NewConvention("units.scale", "2")
		s.Conventions = kernel.NewConventionSet(c)
	}) {
		t.Fatal("convention difference must be visible")
	}
	if changes(func(s *kernel.ObjectSpec) {
		s.Expr = kernel.NewRational(big.NewRat(1, 1))
	}) {
		t.Fatal("expr difference must be visible")
	}
	if core.EqualObject(base, core.Object{}) {
		t.Fatal("zero object is never equal to a minted object")
	}
}

func TestObjectDefensiveCopies(t *testing.T) {
	prov, err := kernel.NewProvenance(kernel.StatusIdentified, "src", "framework",
		[]string{strings.Repeat("a", 64)}, [32]byte{}, [32]byte{},
		kernel.MRCVersion, "because")
	if err != nil {
		t.Fatal(err)
	}
	m := kernel.NewRational(big.NewRat(1, 1))
	powExpr := kernel.NewPow(m, big.NewRat(2, 1))
	a, err := kernel.NewTextAssumption(kernel.AssumptionDomain, "regime", "classical")
	if err != nil {
		t.Fatal(err)
	}
	c, err := kernel.NewConvention("metric.signature", "-+++")
	if err != nil {
		t.Fatal(err)
	}
	obj, err := kernel.MintObject(kernel.ObjectSpec{
		Name:         "fixture",
		Kind:         kernel.KindExpression,
		Dimension:    core.Dimensionless(),
		Expr:         powExpr,
		Assumptions:  kernel.NewAssumptionSet(a),
		Conventions:  kernel.NewConventionSet(c),
		Provenance:   prov,
		CorpusStatus: kernel.CorpusNone,
	})
	if err != nil {
		t.Fatal(err)
	}
	before := core.HashObject(obj)

	vs := obj.Assumptions().Values()
	vs[0] = core.Assumption{}
	if obj.Assumptions().Values()[0].Key() != "regime" {
		t.Fatal("Assumptions() must return a defensive copy")
	}
	cs := obj.Conventions().Values()
	cs[0] = core.Convention{}
	if obj.Conventions().Values()[0].Key() != "metric.signature" {
		t.Fatal("Conventions() must return a defensive copy")
	}
	ph := obj.Provenance().ParentHashes()
	ph[0] = strings.Repeat("z", 64)
	if obj.Provenance().ParentHashes()[0] != strings.Repeat("a", 64) {
		t.Fatal("ParentHashes() must return a defensive copy")
	}
	exp := obj.Expr().Exponent()
	exp.Add(exp, big.NewRat(1, 1))
	if obj.Expr().Exponent().Cmp(big.NewRat(2, 1)) != 0 {
		t.Fatal("Exponent() must return a defensive copy")
	}
	if core.HashObject(obj) != before {
		t.Fatal("mutating returned handles must not change the object hash")
	}
}

func TestRepositoryTreeExact(t *testing.T) {
	// REQ-§2-MUST-01, REQ-003-01/02: no file outside the plan §2 tree (plus
	// the pre-existing task documents), no scaffold/generated files.
	allowed := map[string]bool{
		"go.mod":                         true,
		"README.md":                      true,
		".gitignore":                     true,
		".gut":                           true,
		"reality-first-physics-v8.0.md":  true,
		"reality_check.md":               true,
		"reality_first_plain_english.md": true,
		// implementation files (plan §2)
		"internal/kernel/types.go":      true,
		"internal/kernel/mint.go":       true,
		"core/object.go":                true,
		"core/expr.go":                  true,
		"core/dimension.go":             true,
		"core/assumption.go":            true,
		"core/convention.go":            true,
		"core/provenance.go":            true,
		"core/corpus.go":                true,
		"core/canonical.go":             true,
		"core/errors.go":                true,
		"core/object_test.go":           true,
		"core/expr_test.go":             true,
		"ops/arithmetic.go":             true,
		"ops/simplify.go":               true,
		"ops/transform.go":              true,
		"ops/relation.go":               true,
		"ops/dispatch.go":               true,
		"ops/operations_test.go":        true,
		"ops/negative_test.go":          true,
		"session/session.go":            true,
		"session/ledger.go":             true,
		"session/research_candidate.go": true,
		"session/session_test.go":       true,
		"mechanics/primitives.go":       true,
		"mechanics/relations.go":        true,
		"mechanics/manifest.json":       true,
		"mechanics/manifest_test.go":    true,
		"mechanics/relations_test.go":   true,
		"relativity/primitives.go":      true,
		"relativity/relations.go":       true,
		"relativity/manifest.json":      true,
		"relativity/manifest_test.go":   true,
		"relativity/derivation_test.go": true,
		"hypothesis/candidate.go":       true,
		"hypothesis/candidate_test.go":  true,
		"docs/paper-translation.md":     true,
	}
	err := filepath.Walk(moduleRoot(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, rerr := filepath.Rel(moduleRoot(), path)
		if rerr != nil {
			return rerr
		}
		rel = filepath.ToSlash(rel)
		if info.IsDir() {
			if rel == ".git" || rel == "plan10" || rel == "plans" || rel == ".opencode" {
				return filepath.SkipDir
			}
			return nil
		}
		if !allowed[rel] {
			t.Errorf("unexpected file in repository tree: %s", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestForbiddenSourceSurface(t *testing.T) {
	// REQ-§2-MUST-01, REQ-002-02 (no lexer/parser/frontend), REQ-002-07
	// (no numerical execution), REQ-002-08: banned imports and tokens in
	// non-test sources.
	bannedImports := []string{
		"go/parser", "go/scanner", "go/token",
		"math", "math/rand", "os/exec", "net/http", "unsafe",
	}
	bannedTokens := []string{"float32", "float64"}
	eachGoFile(t, false, func(path string, src []byte) {
		rel, _ := filepath.Rel(moduleRoot(), path)
		for _, imp := range bannedImports {
			if strings.Contains(string(src), `"`+imp+`"`) {
				t.Errorf("%s imports banned package %q", rel, imp)
			}
		}
		for _, tok := range bannedTokens {
			if strings.Contains(string(src), tok) {
				t.Errorf("%s contains banned token %q", rel, tok)
			}
		}
	})
}

func TestNoDeferredPackages(t *testing.T) {
	// REQ-002-09, REQ-002-15, REQ-§2-MUST-01: no deferred package directory
	// and no scaffold dependency files anywhere in the implementation tree.
	bannedDirs := []string{
		"physvet", "electromagnetism", "qm", "qft", "statmech", "review",
		"adapter", "adapters", "vendor",
	}
	bannedFiles := []string{"go.sum", "go.work"}
	seen := map[string]bool{}
	err := filepath.Walk(moduleRoot(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, rerr := filepath.Rel(moduleRoot(), path)
		if rerr != nil {
			return rerr
		}
		rel = filepath.ToSlash(rel)
		if info.IsDir() {
			if rel == ".git" || rel == "plan10" || rel == "plans" || rel == ".opencode" {
				return filepath.SkipDir
			}
			seen[info.Name()] = true
			return nil
		}
		base := info.Name()
		for _, banned := range bannedFiles {
			if base == banned {
				t.Errorf("banned file present: %s", rel)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, banned := range bannedDirs {
		if seen[banned] {
			t.Errorf("deferred package directory present: %s", banned)
		}
	}
}

func TestNoDocsImport(t *testing.T) {
	// REQ-§35-MUST-02: no package imports docs (the translation document is
	// never a runtime dependency).
	eachGoFile(t, true, func(path string, src []byte) {
		rel, _ := filepath.Rel(moduleRoot(), path)
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Errorf("%s: %v", rel, err)
			return
		}
		for _, im := range f.Imports {
			p := strings.Trim(im.Path.Value, `"`)
			if strings.HasPrefix(p, "github.com/PithomLabs/phys/docs") ||
				strings.HasSuffix(p, "/docs") {
				t.Errorf("%s imports docs", rel)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// Corpus / Manifest tests (§17, §27)
// ---------------------------------------------------------------------------

func TestManifestStructSurface(t *testing.T) {
	// REQ-§17.0-MUST-01: exact typed manifest structures with correct field order.
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "corpus.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	want := []struct {
		name   string
		fields []string
	}{
		{"Manifest", []string{
			"SchemaVersion", "FrameworkID", "FrameworkName", "CorpusStatus",
			"Assumptions", "Domain", "Limits", "Anomalies", "Items"}},
		{"ManifestDomain", []string{"ID", "Name", "Description"}},
		{"ManifestLimit", []string{"ID", "Description"}},
		{"ManifestAnomaly", []string{
			"ID", "Framework", "Description", "Status",
			"RelatedAssumptions", "RelatedItems", "ResearchRelevance"}},
		{"ManifestReduction", []string{"ID", "Condition", "Result"}},
		{"ManifestItem", []string{
			"ID", "Constructor", "Kind", "Name", "Statement",
			"CanonicalExpr", "Dimension", "ProvenanceStatus", "Source",
			"Assumptions", "DerivableFrom", "ReducesTo",
			"KnownLimits", "Anomalies", "FalsificationConditions"}},
		{"Challenge", []string{"StepID", "Category", "Severity", "Description"}},
		{"Review", []string{"DerivationID", "Challenges", "ReviewerNotes"}},
	}
	for _, w := range want {
		var found bool
		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok || ts.Name.Name != w.name {
					continue
				}
				found = true
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					t.Errorf("%s is not a struct", w.name)
					continue
				}
				var got []string
				if st.Fields != nil {
					for _, field := range st.Fields.List {
						for _, name := range field.Names {
							got = append(got, name.Name)
						}
					}
				}
				if len(got) != len(w.fields) {
					t.Errorf("%s: field count %d != %d", w.name, len(got), len(w.fields))
				}
				for i, wantField := range w.fields {
					if i >= len(got) || got[i] != wantField {
						t.Errorf("%s field %d = %q, want %q", w.name, i, got[i], wantField)
					}
				}
				break
			}
		}
		if !found {
			t.Errorf("type %s not found in corpus.go", w.name)
		}
	}
}

func TestManifestRejectsInconsistent(t *testing.T) {
	// REQ-032-15: unknown field, bad enum, non-canonical expr, duplicate ID → core.ManifestValidationError
	base := `{
		"schema_version": "1",
		"framework_id": "test",
		"framework_name": "Test",
		"corpus_status": "established",
		"assumptions": [],
		"domain": [],
		"limits": [],
		"anomalies": [],
		"items": [],
		"unknown_field": 1
	}`
	// Unknown field
	_, err := core.ParseManifest([]byte(base))
	if err == nil {
		t.Fatal("unknown field accepted")
	}
	var mve core.ManifestValidationError
	if !errors.As(err, &mve) {
		t.Errorf("error type %T", err)
	}
	// Bad corpus_status enum
	badEnum := strings.Replace(base, `"established"`, `"bogus"`, 1)
	_, err = core.ParseManifest([]byte(badEnum))
	if err == nil {
		t.Fatal("bad corpus_status accepted")
	}
	if !errors.As(err, &mve) {
		t.Errorf("bad enum error type %T", err)
	}
	// Non-canonical expr (extra field in expression)
	badExpr := strings.Replace(base, `"items": []`, `"items": [{"id":"i1","constructor":"c","kind":"Expression","name":"n","statement":"","canonical_expr":{"kind":"symbol","name":"x","extra":1},"dimension":{},"provenance_status":"DEFINED","source":"","assumptions":[],"derivable_from":[],"reduces_to":[],"known_limits":[],"anomalies":[],"falsification_conditions":[]}]`, 1)
	_, err = core.ParseManifest([]byte(badExpr))
	if err == nil {
		t.Fatal("non-canonical expr accepted")
	}
	if !errors.As(err, &mve) {
		t.Errorf("bad expr error type %T", err)
	}
	// Duplicate item ID
	dupID := strings.Replace(base, `"items": []`, `"items": [{"id":"i1","constructor":"c","kind":"Expression","name":"n","statement":"","canonical_expr":{"kind":"symbol","name":"x"},"dimension":{},"provenance_status":"DEFINED","source":"","assumptions":[],"derivable_from":[],"reduces_to":[],"known_limits":[],"anomalies":[],"falsification_conditions":[]},{"id":"i1","constructor":"c2","kind":"Expression","name":"n2","statement":"","canonical_expr":{"kind":"symbol","name":"y"},"dimension":{},"provenance_status":"DEFINED","source":"","assumptions":[],"derivable_from":[],"reduces_to":[],"known_limits":[],"anomalies":[],"falsification_conditions":[]}]`, 1)
	_, err = core.ParseManifest([]byte(dupID))
	if err == nil {
		t.Fatal("duplicate item ID accepted")
	}
	if !errors.As(err, &mve) {
		t.Errorf("duplicate ID error type %T", err)
	}
}

func TestChallengeReviewFieldOrder(t *testing.T) {
	// REQ-§27.1/27.2-MUST-01: exact field order in AST for Challenge and Review.
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "corpus.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []struct {
		name   string
		fields []string
	}{
		{"Challenge", []string{"StepID", "Category", "Severity", "Description"}},
		{"Review", []string{"DerivationID", "Challenges", "ReviewerNotes"}},
	} {
		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok || ts.Name.Name != want.name {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				var got []string
				if st.Fields != nil {
					for _, field := range st.Fields.List {
						for _, name := range field.Names {
							got = append(got, name.Name)
						}
					}
				}
				for i, f := range want.fields {
					if i >= len(got) || got[i] != f {
						t.Errorf("%s field order %d = %q, want %q", want.name, i, got[i], f)
					}
				}
			}
		}
	}
}

func TestReviewCategoryEnumExact(t *testing.T) {
	// REQ-§27.3-MUST-01: exact eight categories with correct string values.
	want := []struct {
		constant core.ReviewCategory
		value    string
	}{
		{core.ReviewCategoryError, "CategoryError"},
		{core.ReviewCategoryDimensionError, "DimensionError"},
		{core.ReviewCategoryAssumptionConflict, "AssumptionConflict"},
		{core.ReviewCategoryConventionConflict, "ConventionConflict"},
		{core.ReviewCategoryUnsupportedIdentification, "UnsupportedIdentification"},
		{core.ReviewCategoryInvalidReduction, "InvalidReduction"},
		{core.ReviewCategoryProvenanceProblem, "ProvenanceProblem"},
		{core.ReviewCategoryCandidateOverreach, "CandidateOverreach"},
	}
	if len(want) != 8 {
		t.Fatal("exactly 8 review categories")
	}
	for _, w := range want {
		if string(w.constant) != w.value {
			t.Errorf("core.ReviewCategory %q = %q, want %q", w.constant, string(w.constant), w.value)
		}
		if !core.ValidReviewCategory(w.constant) {
			t.Errorf("core.ValidReviewCategory(%q) = false", w.value)
		}
	}
	if core.ValidReviewCategory("BogusCategory") {
		t.Fatal("invalid category accepted")
	}
	// Ordinal stability: constants are typed string aliases; order in source
	// is the ordinal. The eight categories must appear in the order above.
	// We verify by checking the const block order in corpus.go source.
	src, err := os.ReadFile("corpus.go")
	if err != nil {
		t.Fatal(err)
	}
	order := []string{
		"ReviewCategoryError",
		"ReviewCategoryDimensionError",
		"ReviewCategoryAssumptionConflict",
		"ReviewCategoryConventionConflict",
		"ReviewCategoryUnsupportedIdentification",
		"ReviewCategoryInvalidReduction",
		"ReviewCategoryProvenanceProblem",
		"ReviewCategoryCandidateOverreach",
	}
	last := -1
	for _, name := range order {
		idx := strings.Index(string(src), name)
		if idx == -1 {
			t.Errorf("category %s not found in corpus.go", name)
			continue
		}
		if idx < last {
			t.Errorf("category %s appears before previous in source", name)
		}
		last = idx
	}
}

func TestStatementNotParsed(t *testing.T) {
	// REQ-§17.6-MUST-01: statement is stored as opaque string, never parsed.
	// Verify by embedding arbitrary non-mathematical text that would break a parser.
	malformed := `{
		"schema_version": "1",
		"framework_id": "test",
		"framework_name": "Test",
		"corpus_status": "established",
		"assumptions": [],
		"domain": [],
		"limits": [],
		"anomalies": [],
		"items": [{
			"id": "i1",
			"constructor": "test.NewObj",
			"kind": "Expression",
			"name": "test",
			"statement": "∀x ∃y: x² + y² = z² 🧮 ∫∂∇",
			"canonical_expr": {"kind":"symbol","name":"x"},
			"dimension": {"m":"0/1","l":"0/1","t":"0/1","i":"0/1","theta":"0/1","n":"0/1","j":"0/1"},
			"provenance_status": "DEFINED",
			"source": "",
			"assumptions": [],
			"derivable_from": [],
			"reduces_to": [],
			"known_limits": [],
			"anomalies": [],
			"falsification_conditions": []
		}]
	}`
	m, err := core.ParseManifest([]byte(malformed))
	if err != nil {
		t.Fatalf("valid manifest with arbitrary statement rejected: %v", err)
	}
	if len(m.Items) != 1 {
		t.Fatal("item not parsed")
	}
	if m.Items[0].Statement != "∀x ∃y: x² + y² = z² 🧮 ∫∂∇" {
		t.Errorf("statement not preserved: %q", m.Items[0].Statement)
	}
	// G-Audit: ensure no non-test file imports go/parser, go/scanner, go/token
	// (this test itself imports them for AST inspection — that's allowed)
}

func TestValidateManifestBytesPure(t *testing.T) {
	// REQ-§17.7-MUST-02: core.ValidateManifestBytes has no constructor calls and no I/O.
	// Source scan: no calls to MintObject, New*, or os.ReadFile/io.ReadFile.
	src, err := os.ReadFile("corpus.go")
	if err != nil {
		t.Fatal(err)
	}
	disallowed := []string{
		"MintObject(", "NewMass(", "NewTime(", "NewPosition(", "NewVelocity(",
		"NewAcceleration(", "NewForce(", "NewMomentum(", "NewEnergy(",
		"NewKineticEnergy(", "NewRestMass(", "NewSpacetime(", "NewMinkowskiMetric(",
		"NewThreeMomentum(", "NewFourMomentum(", "NewSpeedOfLight(",
		"LorentzFactor(", "EnergyMomentumRelation(", "MassEnergyRelation(",
		"ZeroThreeMomentum(", "ZeroEnergy(", "ZeroVelocity(",
		"os.ReadFile", "io.ReadFile",
	}
	for _, d := range disallowed {
		if strings.Contains(string(src), d) {
			t.Errorf("corpus.go contains disallowed call %q", d)
		}
	}
	// core.ValidateManifestBytes just calls core.ParseManifest (which we verify above).
	valid := `{"schema_version":"1","framework_id":"test","framework_name":"Test","corpus_status":"established","assumptions":[],"domain":[],"limits":[],"anomalies":[],"items":[]}`
	if _, err := core.ValidateManifestBytes([]byte(valid)); err != nil {
		t.Fatalf("core.ValidateManifestBytes failed: %v", err)
	}
}

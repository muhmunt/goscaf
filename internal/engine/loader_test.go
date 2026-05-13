package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/muhmunt/goscaf/internal/wizard"
)

func TestLoadPreset_REST(t *testing.T) {
	r := &wizard.Result{Transport: "rest"}
	s, fsys, err := loadPreset(r)
	if err != nil {
		t.Fatalf("loadPreset(rest): %v", err)
	}
	if s.Name != "rest-service" {
		t.Errorf("name = %q, want rest-service", s.Name)
	}
	if fsys == nil {
		t.Fatal("fs is nil")
	}
	if len(s.Files) == 0 {
		t.Error("no files in scaffold")
	}
}

func TestLoadPreset_GRPC(t *testing.T) {
	r := &wizard.Result{Transport: "grpc"}
	s, _, err := loadPreset(r)
	if err != nil {
		t.Fatalf("loadPreset(grpc): %v", err)
	}
	if s.Name != "grpc-service" {
		t.Errorf("name = %q, want grpc-service", s.Name)
	}
}

func TestLoadPreset_Both_RoutesToGRPC(t *testing.T) {
	r := &wizard.Result{Transport: "both"}
	s, _, err := loadPreset(r)
	if err != nil {
		t.Fatalf("loadPreset(both): %v", err)
	}
	if s.Name != "grpc-service" {
		t.Errorf("transport=both should resolve to grpc preset, got %q", s.Name)
	}
}

func TestLoadLayer(t *testing.T) {
	tests := []struct {
		layer    string
		wantName string
	}{
		{"handler", "add-handler"},
		{"usecase", "add-usecase"},
		{"repository", "add-repository"},
	}

	for _, tt := range tests {
		t.Run(tt.layer, func(t *testing.T) {
			s, fsys, err := loadLayer(tt.layer)
			if err != nil {
				t.Fatalf("loadLayer(%q): %v", tt.layer, err)
			}
			if s.Name != tt.wantName {
				t.Errorf("name = %q, want %q", s.Name, tt.wantName)
			}
			if fsys == nil {
				t.Error("fs is nil")
			}
			if len(s.Files) == 0 {
				t.Error("no files in scaffold")
			}
		})
	}
}

func TestLoadLayer_Unknown(t *testing.T) {
	_, _, err := loadLayer("nonexistent")
	if err == nil {
		t.Error("expected error for unknown layer, got nil")
	}
}

func TestLoadCustom(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "scaffold.yaml")

	err := os.WriteFile(path, []byte(`
name: my-custom
version: "1.0.0"
files:
  - src: files/main.go.tmpl
    dst: main.go
hooks:
  pre: []
  post: []
`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	s, fsys, err := loadCustom(path)
	if err != nil {
		t.Fatalf("loadCustom: %v", err)
	}
	if s.Name != "my-custom" {
		t.Errorf("name = %q, want my-custom", s.Name)
	}
	if fsys == nil {
		t.Error("fs is nil")
	}
}

func TestLoadCustom_MissingFile(t *testing.T) {
	_, _, err := loadCustom("/does/not/exist/scaffold.yaml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadCustom_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "scaffold.yaml")
	os.WriteFile(path, []byte(":\tnot: valid: yaml:"), 0o644)

	_, _, err := loadCustom(path)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

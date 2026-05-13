package engine

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/muhmunt/goscaf/internal/schema"
	"github.com/muhmunt/goscaf/internal/templates"
	"github.com/muhmunt/goscaf/internal/wizard"
	"gopkg.in/yaml.v3"
)

func load(r *wizard.Result) (*schema.Scaffold, fs.FS, error) {
	if r.TemplatePath != "" {
		return loadCustom(r.TemplatePath)
	}
	if r.Layer != "" {
		return loadLayer(r.Layer)
	}
	return loadPreset(r)
}

func loadCustom(path string) (*schema.Scaffold, fs.FS, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("read custom scaffold.yaml: %w", err)
	}
	return parseScaffold(data, os.DirFS(filepath.Dir(path)))
}

func loadLayer(layer string) (*schema.Scaffold, fs.FS, error) {
	path := fmt.Sprintf("layers/%s/scaffold.yaml", layer)
	data, err := templates.FS.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("unknown layer %q — valid values: handler, usecase, repository", layer)
	}
	return parseScaffold(data, templates.FS)
}

func loadPreset(r *wizard.Result) (*schema.Scaffold, fs.FS, error) {
	preset := resolvePreset(r)
	path := fmt.Sprintf("presets/%s/scaffold.yaml", preset)
	data, err := templates.FS.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("read built-in preset %q: %w", preset, err)
	}
	return parseScaffold(data, templates.FS)
}

func parseScaffold(data []byte, fsys fs.FS) (*schema.Scaffold, fs.FS, error) {
	var s schema.Scaffold
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, nil, fmt.Errorf("parse scaffold.yaml: %w", err)
	}
	return &s, fsys, nil
}

func resolvePreset(r *wizard.Result) string {
	switch r.Transport {
	case "grpc", "both":
		return "grpc"
	default:
		return "rest"
	}
}

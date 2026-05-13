package engine

import (
	"bytes"
	"fmt"
	"io/fs"
	"strings"
	"text/template"
)

// funcMap is shared across all template executions.
var funcMap = template.FuncMap{
	// title converts "user-profile" or "order_item" into "UserProfile" / "OrderItem".
	"title": func(s string) string {
		parts := strings.FieldsFunc(s, func(r rune) bool { return r == '-' || r == '_' })
		for i, p := range parts {
			if len(p) > 0 {
				parts[i] = strings.ToUpper(p[:1]) + p[1:]
			}
		}
		return strings.Join(parts, "")
	},
}

// evalCondition executes a Go template expression (e.g. "{{eq .transport "rest"}}") and
// returns true if the result is non-empty/truthy.
func evalCondition(condition string, vars map[string]any) (bool, error) {
	if condition == "" {
		return true, nil
	}

	// Strip surrounding {{ }} so we can wrap in {{if ...}}
	expr := strings.TrimSpace(condition)
	expr = strings.TrimPrefix(expr, "{{")
	expr = strings.TrimSuffix(expr, "}}")
	expr = strings.TrimSpace(expr)

	tmplStr := fmt.Sprintf(`{{if %s}}1{{end}}`, expr)
	tmpl, err := template.New("cond").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return false, fmt.Errorf("parse condition %q: %w", condition, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return false, fmt.Errorf("execute condition %q: %w", condition, err)
	}

	return buf.String() == "1", nil
}

// renderString renders a Go template string (e.g. "cmd/{{.service_name}}/main.go").
func renderString(s string, vars map[string]any) (string, error) {
	tmpl, err := template.New("str").Funcs(funcMap).Parse(s)
	if err != nil {
		return "", fmt.Errorf("parse string template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("execute string template: %w", err)
	}

	return buf.String(), nil
}

// renderTemplate reads a .tmpl file from fsys and renders it with vars.
func renderTemplate(fsys fs.FS, src string, vars map[string]any) ([]byte, error) {
	data, err := fs.ReadFile(fsys, src)
	if err != nil {
		return nil, fmt.Errorf("read template file %q: %w", src, err)
	}

	tmpl, err := template.New(src).Funcs(funcMap).Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse template %q: %w", src, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return nil, fmt.Errorf("execute template %q: %w", src, err)
	}

	return buf.Bytes(), nil
}

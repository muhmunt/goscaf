package engine

import (
	"testing"
	"testing/fstest"
)

func TestEvalCondition(t *testing.T) {
	vars := map[string]any{
		"transport":      "rest",
		"database":       "none",
		"include_docker": true,
		"include_ci":     false,
	}

	tests := []struct {
		name      string
		condition string
		want      bool
	}{
		{"empty is always true", "", true},
		{"eq match", `{{eq .transport "rest"}}`, true},
		{"eq no match", `{{eq .transport "grpc"}}`, false},
		{"ne false", `{{ne .database "none"}}`, false},
		{"ne true", `{{ne .transport "grpc"}}`, true},
		{"or first matches", `{{or (eq .transport "rest") (eq .transport "both")}}`, true},
		{"or second matches", `{{or (eq .transport "grpc") (eq .transport "rest")}}`, true},
		{"or neither matches", `{{or (eq .transport "grpc") (eq .transport "both")}}`, false},
		{"and both true", `{{and (eq .transport "rest") (eq .database "none")}}`, true},
		{"and one false", `{{and (eq .transport "rest") (ne .database "none")}}`, false},
		{"bool true", `{{.include_docker}}`, true},
		{"bool false", `{{.include_ci}}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evalCondition(tt.condition, vars)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("evalCondition(%q) = %v, want %v", tt.condition, got, tt.want)
			}
		})
	}
}

func TestEvalCondition_InvalidTemplate(t *testing.T) {
	_, err := evalCondition("{{not valid", map[string]any{})
	if err == nil {
		t.Error("expected error for malformed template, got nil")
	}
}

func TestRenderString(t *testing.T) {
	vars := map[string]any{
		"service_name":  "order-service",
		"resource_name": "user-profile",
	}

	tests := []struct {
		input string
		want  string
	}{
		{"no vars", "no vars"},
		{"cmd/{{.service_name}}/main.go", "cmd/order-service/main.go"},
		{"internal/handler/{{.resource_name}}_handler.go", "internal/handler/user-profile_handler.go"},
		{"{{.resource_name | title}}Handler", "UserProfileHandler"},
		{"{{.service_name | title}}Service", "OrderServiceService"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := renderString(tt.input, vars)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("renderString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRenderTemplate(t *testing.T) {
	fsys := fstest.MapFS{
		"service.go.tmpl": &fstest.MapFile{
			Data: []byte("package main\n\ntype {{.resource_name | title}}Service struct{}\n"),
		},
	}

	vars := map[string]any{"resource_name": "order-item"}
	got, err := renderTemplate(fsys, "service.go.tmpl", vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "package main\n\ntype OrderItemService struct{}\n"
	if string(got) != want {
		t.Errorf("renderTemplate =\n%s\nwant:\n%s", got, want)
	}
}

func TestRenderTemplate_MissingFile(t *testing.T) {
	_, err := renderTemplate(fstest.MapFS{}, "nonexistent.tmpl", map[string]any{})
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestRenderTemplate_InvalidTemplate(t *testing.T) {
	fsys := fstest.MapFS{
		"bad.tmpl": &fstest.MapFile{Data: []byte("{{not valid")},
	}
	_, err := renderTemplate(fsys, "bad.tmpl", map[string]any{})
	if err == nil {
		t.Error("expected error for malformed template, got nil")
	}
}

func TestTitleFunc(t *testing.T) {
	fn := funcMap["title"].(func(string) string)

	tests := []struct {
		input string
		want  string
	}{
		{"user", "User"},
		{"order", "Order"},
		{"user-profile", "UserProfile"},
		{"order_item", "OrderItem"},
		{"my-order-service", "MyOrderService"},
		{"already", "Already"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := fn(tt.input)
			if got != tt.want {
				t.Errorf("title(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

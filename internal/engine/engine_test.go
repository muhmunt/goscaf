package engine_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/muhmunt/goscaf/internal/engine"
	"github.com/muhmunt/goscaf/internal/wizard"
)

// TestRun_CustomTemplate verifies that a custom scaffold.yaml renders files
// with the correct variable substitution and the title function.
func TestRun_CustomTemplate(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()

	writeFile(t, filepath.Join(srcDir, "scaffold.yaml"), `
name: test-scaffold
version: "1.0.0"
files:
  - src: files/main.go.tmpl
    dst: "cmd/{{.service_name}}/main.go"
  - src: files/handler.go.tmpl
    dst: "internal/handler/handler.go"
    condition: '{{eq .transport "rest"}}'
hooks:
  pre: []
  post: []
`)
	writeFile(t, filepath.Join(srcDir, "files/main.go.tmpl"),
		"package main\n// {{.service_name | title}} entry point\n")
	writeFile(t, filepath.Join(srcDir, "files/handler.go.tmpl"),
		"package handler\n// {{.service_name | title}}Handler\n")

	r := &wizard.Result{
		ServiceName:  "order-service",
		ModulePath:   "github.com/test/order-service",
		Transport:    "rest",
		OutputDir:    outDir,
		TemplatePath: filepath.Join(srcDir, "scaffold.yaml"),
	}

	if err := engine.Run(r); err != nil {
		t.Fatalf("Run: %v", err)
	}

	assertContains(t, filepath.Join(outDir, "cmd/order-service/main.go"),
		"// OrderService entry point")
	assertContains(t, filepath.Join(outDir, "internal/handler/handler.go"),
		"// OrderServiceHandler")
}

// TestRun_ConditionExcludesFile verifies that files with a false condition are skipped.
func TestRun_ConditionExcludesFile(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()

	writeFile(t, filepath.Join(srcDir, "scaffold.yaml"), `
name: test-scaffold
version: "1.0.0"
files:
  - src: files/base.go.tmpl
    dst: "base.go"
  - src: files/grpc.go.tmpl
    dst: "grpc.go"
    condition: '{{eq .transport "grpc"}}'
hooks:
  pre: []
  post: []
`)
	writeFile(t, filepath.Join(srcDir, "files/base.go.tmpl"), "package main\n")
	writeFile(t, filepath.Join(srcDir, "files/grpc.go.tmpl"), "package main // grpc\n")

	r := &wizard.Result{
		ServiceName:  "test-svc",
		Transport:    "rest",
		OutputDir:    outDir,
		TemplatePath: filepath.Join(srcDir, "scaffold.yaml"),
	}

	if err := engine.Run(r); err != nil {
		t.Fatalf("Run: %v", err)
	}

	assertExists(t, filepath.Join(outDir, "base.go"))
	assertNotExists(t, filepath.Join(outDir, "grpc.go"))
}

// TestRun_DstPathRendered verifies that dst paths containing template vars are rendered.
func TestRun_DstPathRendered(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()

	writeFile(t, filepath.Join(srcDir, "scaffold.yaml"), `
name: test-scaffold
version: "1.0.0"
files:
  - src: files/handler.go.tmpl
    dst: "internal/handler/{{.resource_name}}_handler.go"
hooks:
  pre: []
  post: []
`)
	writeFile(t, filepath.Join(srcDir, "files/handler.go.tmpl"), "package handler\n")

	r := &wizard.Result{
		ServiceName:  "user",
		OutputDir:    outDir,
		TemplatePath: filepath.Join(srcDir, "scaffold.yaml"),
	}

	if err := engine.Run(r); err != nil {
		t.Fatalf("Run: %v", err)
	}

	assertExists(t, filepath.Join(outDir, "internal/handler/user_handler.go"))
}

// TestRun_AddUsecaseLayer verifies the built-in usecase layer generates the correct file.
func TestRun_AddUsecaseLayer(t *testing.T) {
	outDir := t.TempDir()

	r := &wizard.Result{
		ServiceName: "user",
		ModulePath:  "github.com/test/my-service",
		Layer:       "usecase",
		OutputDir:   outDir,
	}

	if err := engine.Run(r); err != nil {
		t.Fatalf("Run add usecase: %v", err)
	}

	assertContains(t, filepath.Join(outDir, "internal/usecase/user_usecase.go"),
		"type UserUseCase struct")
	assertContains(t, filepath.Join(outDir, "internal/usecase/user_usecase.go"),
		"func NewUserUseCase(")
	assertContains(t, filepath.Join(outDir, "internal/usecase/user_usecase.go"),
		"github.com/test/my-service/internal/repository")
}

// TestRun_AddRepositoryLayer verifies the built-in repository layer generates the correct file.
func TestRun_AddRepositoryLayer(t *testing.T) {
	outDir := t.TempDir()

	r := &wizard.Result{
		ServiceName: "order",
		ModulePath:  "github.com/test/my-service",
		Layer:       "repository",
		OutputDir:   outDir,
	}

	if err := engine.Run(r); err != nil {
		t.Fatalf("Run add repository: %v", err)
	}

	assertContains(t, filepath.Join(outDir, "internal/repository/order_repository.go"),
		"type OrderRepository struct")
	assertContains(t, filepath.Join(outDir, "internal/repository/order_repository.go"),
		"func NewOrderRepository(")
}

// TestRun_AddHandlerLayer_HTTP verifies that only the HTTP handler is generated for transport=rest.
func TestRun_AddHandlerLayer_HTTP(t *testing.T) {
	outDir := t.TempDir()

	r := &wizard.Result{
		ServiceName: "payment",
		ModulePath:  "github.com/test/my-service",
		Transport:   "rest",
		Layer:       "handler",
		OutputDir:   outDir,
	}

	if err := engine.Run(r); err != nil {
		t.Fatalf("Run add handler (rest): %v", err)
	}

	assertContains(t, filepath.Join(outDir, "internal/handler/payment_handler.go"),
		"type PaymentHandler struct")
	assertContains(t, filepath.Join(outDir, "internal/handler/payment_handler.go"),
		"func NewPaymentHandler(")
	assertNotExists(t, filepath.Join(outDir, "internal/handler/payment_grpc_handler.go"))
}

// TestRun_AddHandlerLayer_GRPC verifies that only the gRPC handler is generated for transport=grpc.
func TestRun_AddHandlerLayer_GRPC(t *testing.T) {
	outDir := t.TempDir()

	r := &wizard.Result{
		ServiceName: "payment",
		ModulePath:  "github.com/test/my-service",
		Transport:   "grpc",
		Layer:       "handler",
		OutputDir:   outDir,
	}

	if err := engine.Run(r); err != nil {
		t.Fatalf("Run add handler (grpc): %v", err)
	}

	assertContains(t, filepath.Join(outDir, "internal/handler/payment_grpc_handler.go"),
		"type PaymentGRPCHandler struct")
	assertNotExists(t, filepath.Join(outDir, "internal/handler/payment_handler.go"))
}

// TestRun_AddHandlerLayer_Both verifies that both HTTP and gRPC handlers are generated for transport=both.
func TestRun_AddHandlerLayer_Both(t *testing.T) {
	outDir := t.TempDir()

	r := &wizard.Result{
		ServiceName: "payment",
		ModulePath:  "github.com/test/my-service",
		Transport:   "both",
		Layer:       "handler",
		OutputDir:   outDir,
	}

	if err := engine.Run(r); err != nil {
		t.Fatalf("Run add handler (both): %v", err)
	}

	assertExists(t, filepath.Join(outDir, "internal/handler/payment_handler.go"))
	assertExists(t, filepath.Join(outDir, "internal/handler/payment_grpc_handler.go"))
}

// TestRun_HyphenatedResourceName verifies title converts "order-item" to "OrderItem".
func TestRun_HyphenatedResourceName(t *testing.T) {
	outDir := t.TempDir()

	r := &wizard.Result{
		ServiceName: "order-item",
		ModulePath:  "github.com/test/my-service",
		Layer:       "repository",
		OutputDir:   outDir,
	}

	if err := engine.Run(r); err != nil {
		t.Fatalf("Run: %v", err)
	}

	assertContains(t, filepath.Join(outDir, "internal/repository/order-item_repository.go"),
		"type OrderItemRepository struct")
}

// ── helpers ──────────────────────────────────────────────────────────────────

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func assertContains(t *testing.T, path, substr string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !strings.Contains(string(data), substr) {
		t.Errorf("%s: missing %q\ngot:\n%s", path, substr, data)
	}
}

func assertExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %s", path)
	}
}

func assertNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Errorf("expected file NOT to exist: %s", path)
	}
}

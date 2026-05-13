# goscaf

A CLI that bootstraps production-ready Go microservices from templates — interactively or fully non-interactive.

Pick a transport, database, and architecture pattern. `goscaf` generates the directory structure, wires the layers together, and drops in a Dockerfile, docker-compose, CI workflow, and Makefile. Single binary, no runtime dependencies.

---

## Install

**From source**

```bash
git clone https://github.com/muhmunt/goscaf
cd goscaf
make install        # go install → $GOPATH/bin/goscaf
```

**Go install**

```bash
go install github.com/muhmunt/goscaf/cmd/goscaf@latest
```

**Pre-built binary** — download from [Releases](https://github.com/muhmunt/goscaf/releases), unzip, and move to your `$PATH`.

---

## Quick start

```bash
goscaf new
```

Answer a few questions:

```
? Service name?          order-service
? Go module path?        github.com/muhmunt/order-service
? Transport layer?       REST (gin)
? Database?              PostgreSQL
? Architecture?          Clean Architecture
? Go version?            1.22
? Include Docker?        Yes
? Include GitHub CI?     Yes
```

Generated output:

```
order-service/
├── cmd/order-service/main.go        # HTTP server with graceful shutdown
├── internal/
│   ├── handler/handler.go           # HTTP handler (List, GetByID, Create, Update, Delete)
│   ├── usecase/usecase.go           # Business logic layer
│   └── repository/repository.go    # Data access layer
├── config/config.yaml
├── Dockerfile                       # Multi-stage build
├── docker-compose.yml               # With postgres healthcheck
├── .github/workflows/ci.yml         # Vet + test + GoReleaser on tag
├── Makefile
└── go.mod
```

Then:

```bash
cd order-service
go mod tidy
go run ./cmd/order-service
```

---

## Commands

### `goscaf new`

Bootstraps a full microservice.

```bash
goscaf new                        # interactive wizard
goscaf new --yes --name pay-svc   # non-interactive, uses flags + defaults
goscaf new --dry-run --yes --name pay-svc --transport grpc  # preview only
```

| Flag | Short | Description |
|---|---|---|
| `--name` | | Service name (e.g. `order-service`) |
| `--module` | | Go module path (e.g. `github.com/org/order-service`) |
| `--transport` | | `rest` \| `grpc` \| `both` |
| `--database` | | `postgres` \| `mysql` \| `mongodb` \| `none` |
| `--arch` | | `clean` \| `hexagonal` \| `standard` |
| `--go-version` | | Go version for `go.mod` and Dockerfile (default: `1.22`) |
| `--no-docker` | | Skip Dockerfile + docker-compose |
| `--no-ci` | | Skip GitHub Actions workflow |
| `--template` | `-t` | Path to a custom `scaffold.yaml` |
| `--output` | `-o` | Output directory (default: service name) |
| `--dry-run` | | Print files that would be generated — write nothing |
| `--yes` | `-y` | Skip prompts, use flags + defaults |

### `goscaf add`

Adds a single layer to an existing service. Run from inside your service directory.

```bash
goscaf add handler      # HTTP handler, gRPC handler, or both
goscaf add usecase      # use case wired to repository
goscaf add repository   # repository stub ready for a DB client
```

Example — adding a `user` handler interactively:

```
? Resource name?    user
? Module path?      github.com/muhmunt/order-service
? Transport?        HTTP (REST)

  ✓ internal/handler/user_handler.go
```

Non-interactive:

```bash
goscaf add handler --yes --name user --module github.com/org/order-service --transport rest
goscaf add handler --dry-run --yes --name invoice --module github.com/org/billing --transport both
```

| Flag | Short | Description |
|---|---|---|
| `--name` | | Resource name (e.g. `user`, `order`) |
| `--module` | | Go module path of the existing service |
| `--transport` | | `rest` \| `grpc` \| `both` (handler layer only) |
| `--dry-run` | | Print files that would be generated — write nothing |
| `--yes` | `-y` | Skip prompts, use flags + defaults |

### `goscaf list`

Shows all available presets and layers.

```bash
goscaf list
```

```
Presets  (goscaf new)
  rest    REST microservice — net/http + gin router, Clean Architecture
  grpc    gRPC microservice — buf for protobuf codegen

Layers   (goscaf add <layer>)
  handler       HTTP handler, gRPC handler, or both (--transport rest|grpc|both)
  usecase       Business logic wired to repository
  repository    Data access stub ready for DB injection
```

---

## Presets

### REST (`--transport rest`)

Generates a service using `net/http` + `gin` router.

| File | Description |
|---|---|
| `cmd/<name>/main.go` | HTTP server, graceful shutdown |
| `internal/handler/handler.go` | Gin handler with health endpoint |
| `internal/usecase/usecase.go` | Business logic (Clean/Hexagonal only) |
| `internal/repository/repository.go` | Data access (if database selected) |
| `config/config.yaml` | App + server + DB config |
| `Dockerfile` | Multi-stage, scratch-based binary |
| `docker-compose.yml` | App + DB service with healthcheck |
| `.github/workflows/ci.yml` | Vet, test, GoReleaser on tag |
| `Makefile` | `run`, `build`, `test`, `lint`, `docker-up/down` |

### gRPC (`--transport grpc` or `both`)

Generates a gRPC service with [buf](https://buf.build) for protobuf code generation.

Additional files over REST:

| File | Description |
|---|---|
| `proto/<name>/v1/service.proto` | Starter proto with Health RPC |
| `internal/handler/grpc_handler.go` | Handler stub with `buf generate` instructions |
| `buf.yaml` | Buf module config with lint + breaking change rules |
| `buf.gen.yaml` | Code gen config (go + grpc-go plugins) |

After scaffolding:

```bash
cd my-service
go mod tidy
make proto          # buf generate → writes gen/my-service/v1/*.pb.go
go run ./cmd/my-service
```

---

## Custom templates

You can replace the built-in preset entirely with your own `scaffold.yaml`.

### scaffold.yaml schema

```yaml
name: my-template
description: "Description shown during scaffolding"
version: "1.0.0"

# Variables — map to wizard prompts or CLI flags
vars:
  - name: service_name
    type: string          # string | bool | select | multiselect
    required: true
    prompt: "Service name?"
    validate:
      pattern: "^[a-z][a-z0-9-]*$"
      message: "Lowercase letters, numbers, and hyphens only"

  - name: transport
    type: select
    prompt: "Transport?"
    options:
      - label: "REST"
        value: rest
      - label: "gRPC"
        value: grpc
    default: rest

# Files — each entry maps a template to an output path
# src: path to .tmpl file, relative to this scaffold.yaml
# dst: output path, supports Go template vars
# condition: Go template expression — file is skipped if false
files:
  - src: files/main.go.tmpl
    dst: "cmd/{{.service_name}}/main.go"

  - src: files/handler.go.tmpl
    dst: "internal/handler/handler.go"
    condition: '{{eq .transport "rest"}}'

# Hooks — shell commands run before/after file generation
hooks:
  pre: []
  post:
    - cmd: "go mod tidy"
```

### Template variables

All variables defined in `vars` are available in `.tmpl` files and in `dst` paths using Go's `text/template` syntax:

```go
// files/handler.go.tmpl
package handler

type {{.service_name | title}}Handler struct{}
```

### Template functions

| Function | Input | Output | Example |
|---|---|---|---|
| `title` | `string` | PascalCase string | `"order-item" → "OrderItem"` |

The `title` function splits on `-` and `_`, then capitalises each word — safe for use as Go type names.

### Running a custom template

```bash
goscaf new --template ./path/to/scaffold.yaml
```

`src` paths in your scaffold.yaml are relative to the directory containing it. See [`examples/custom-template/`](examples/custom-template/) for a working example.

---

## Architecture

```
goscaf/
├── cmd/goscaf/           # Entry point
├── internal/
│   ├── cli/              # cobra commands (new, add, list)
│   ├── wizard/           # huh interactive prompts + --yes mode
│   ├── engine/           # load → render → write → hooks (+ --dry-run)
│   ├── schema/           # scaffold.yaml Go types
│   └── templates/        # Embedded presets, layers, shared templates
│       ├── presets/
│       │   ├── rest/
│       │   └── grpc/
│       ├── layers/
│       │   ├── handler/
│       │   ├── usecase/
│       │   └── repository/
│       └── shared/       # Dockerfile, CI, Makefile, go.mod
└── examples/
    └── custom-template/
```

The engine pipeline for every `goscaf` invocation:

```
wizard (collect vars)
  → loader (parse scaffold.yaml from embedded FS or custom path)
    → renderer (evalCondition + renderTemplate per file)
      → writer (MkdirAll + WriteFile)
        → hooks (pre/post shell commands)
```

---

## Development

```bash
make build          # build ./bin/goscaf
make test           # go test -race -cover ./...
make lint           # golangci-lint run ./...
make install        # go install → $GOPATH/bin
make release-dry    # goreleaser --snapshot (no publish)
```

**Releasing** — push a `v*` tag. GitHub Actions runs GoReleaser and publishes binaries for Linux, macOS, and Windows (amd64 + arm64).

```bash
git tag v0.1.0
git push origin v0.1.0
```

---

## Contributing

1. Fork the repo
2. Add your preset or layer under `internal/templates/`
3. Write a `scaffold.yaml` and `.tmpl` files
4. Add tests in `internal/engine/`
5. Open a pull request

**Adding a new preset** — create `internal/templates/presets/<name>/scaffold.yaml` and a `files/` directory. Register the preset name in `internal/engine/loader.go:resolvePreset`.

**Adding a new layer** — create `internal/templates/layers/<name>/scaffold.yaml` and `files/`. The `goscaf add <name>` command will pick it up automatically.

---

## License

MIT

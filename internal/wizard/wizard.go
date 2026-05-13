package wizard

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

// Result holds user answers collected during the interactive session.
type Result struct {
	ServiceName   string
	ModulePath    string
	Transport     string // rest | grpc | both
	Database      string // postgres | mysql | mongodb | none
	Architecture  string // standard | clean | hexagonal
	GoVersion     string
	IncludeDocker bool
	IncludeCI     bool
	OutputDir     string
	TemplatePath  string
	Layer         string // set by RunAdd
	DryRun        bool
	Yes           bool // skip interactive prompts, use flag values + defaults
}

// RunNew runs the interactive wizard for `goscaf new`.
// r is pre-populated from CLI flags; wizard fills in anything that's still zero.
func RunNew(r *Result) (*Result, error) {
	if r == nil {
		r = &Result{}
	}

	// Apply defaults for fields not set by flags
	if r.GoVersion == "" {
		r.GoVersion = "1.22"
	}
	if r.Transport == "" {
		r.Transport = "rest"
	}
	if r.Database == "" {
		r.Database = "none"
	}
	if r.Architecture == "" {
		r.Architecture = "clean"
	}

	if r.Yes {
		return applyNewDefaults(r)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Service name?").
				Description("Lowercase letters, numbers, and hyphens only.").
				Placeholder("order-service").
				Value(&r.ServiceName).
				Validate(validateServiceName),

			huh.NewInput().
				Title("Go module path?").
				Placeholder("github.com/org/order-service").
				Value(&r.ModulePath),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Transport layer?").
				Options(
					huh.NewOption("REST  (gin)", "rest"),
					huh.NewOption("gRPC", "grpc"),
					huh.NewOption("Both", "both"),
				).
				Value(&r.Transport),

			huh.NewSelect[string]().
				Title("Database?").
				Options(
					huh.NewOption("PostgreSQL", "postgres"),
					huh.NewOption("MySQL", "mysql"),
					huh.NewOption("MongoDB", "mongodb"),
					huh.NewOption("None", "none"),
				).
				Value(&r.Database),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Architecture?").
				Options(
					huh.NewOption("Clean Architecture", "clean"),
					huh.NewOption("Hexagonal", "hexagonal"),
					huh.NewOption("Standard (flat internal/)", "standard"),
				).
				Value(&r.Architecture),

			huh.NewInput().
				Title("Go version?").
				Placeholder("1.22").
				Value(&r.GoVersion),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Include Docker + docker-compose?").
				Value(&r.IncludeDocker),

			huh.NewConfirm().
				Title("Include GitHub Actions CI?").
				Value(&r.IncludeCI),
		),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}

	return applyNewDefaults(r)
}

// RunAdd runs the interactive wizard for `goscaf add`.
// r is pre-populated from CLI flags; wizard fills in anything that's still zero.
func RunAdd(layer string, r *Result) (*Result, error) {
	if r == nil {
		r = &Result{}
	}
	r.Layer = layer
	if r.OutputDir == "" {
		r.OutputDir = "."
	}
	if r.Transport == "" {
		r.Transport = "rest"
	}

	if r.Yes {
		if r.ServiceName == "" {
			return nil, fmt.Errorf("--name is required when using --yes")
		}
		if r.ModulePath == "" {
			return nil, fmt.Errorf("--module is required when using --yes")
		}
		return r, nil
	}

	groups := []*huh.Group{
		huh.NewGroup(
			huh.NewInput().
				Title("Resource name?").
				Description("e.g. user, order, payment").
				Placeholder("user").
				Value(&r.ServiceName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("resource name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Module path?").
				Description("Go module path of the existing service").
				Placeholder("github.com/org/my-service").
				Value(&r.ModulePath).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("module path is required")
					}
					return nil
				}),
		),
	}

	if layer == "handler" {
		groups = append(groups, huh.NewGroup(
			huh.NewSelect[string]().
				Title("Transport?").
				Options(
					huh.NewOption("HTTP (REST)", "rest"),
					huh.NewOption("gRPC", "grpc"),
					huh.NewOption("Both", "both"),
				).
				Value(&r.Transport),
		))
	}

	if err := huh.NewForm(groups...).Run(); err != nil {
		return nil, err
	}

	return r, nil
}

func applyNewDefaults(r *Result) (*Result, error) {
	if r.ServiceName == "" {
		return nil, fmt.Errorf("service name is required (use --name or answer the prompt)")
	}
	if r.ModulePath == "" {
		r.ModulePath = fmt.Sprintf("github.com/muhmunt/%s", r.ServiceName)
	}
	if r.OutputDir == "" {
		r.OutputDir = r.ServiceName
	}
	return r, nil
}

func validateServiceName(s string) error {
	if s == "" {
		return fmt.Errorf("service name is required")
	}
	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
			return fmt.Errorf("only lowercase letters, numbers, and hyphens allowed")
		}
	}
	return nil
}

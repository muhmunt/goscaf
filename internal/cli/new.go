package cli

import (
	"fmt"

	"github.com/muhmunt/goscaf/internal/engine"
	"github.com/muhmunt/goscaf/internal/wizard"
	"github.com/spf13/cobra"
)

var (
	nameFlag         string
	moduleFlag       string
	transportFlag    string
	databaseFlag     string
	architectureFlag string
	goVersionFlag    string
	noDockerFlag     bool
	noCIFlag         bool
	templateFlag     string
	outputFlag       string
	dryRunFlag       bool
	yesFlag          bool
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Bootstrap a new Go microservice",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := &wizard.Result{
			ServiceName:   nameFlag,
			ModulePath:    moduleFlag,
			Transport:     transportFlag,
			Database:      databaseFlag,
			Architecture:  architectureFlag,
			GoVersion:     goVersionFlag,
			IncludeDocker: !noDockerFlag,
			IncludeCI:     !noCIFlag,
			OutputDir:     outputFlag,
			TemplatePath:  templateFlag,
			DryRun:        dryRunFlag,
			Yes:           yesFlag,
		}

		result, err := wizard.RunNew(r)
		if err != nil {
			return fmt.Errorf("wizard: %w", err)
		}
		return engine.Run(result)
	},
}

func init() {
	newCmd.Flags().StringVar(&nameFlag, "name", "", "service name (e.g. order-service)")
	newCmd.Flags().StringVar(&moduleFlag, "module", "", "Go module path (e.g. github.com/org/order-service)")
	newCmd.Flags().StringVar(&transportFlag, "transport", "", "transport layer: rest | grpc | both")
	newCmd.Flags().StringVar(&databaseFlag, "database", "", "database: postgres | mysql | mongodb | none")
	newCmd.Flags().StringVar(&architectureFlag, "arch", "", "architecture: clean | hexagonal | standard")
	newCmd.Flags().StringVar(&goVersionFlag, "go-version", "", "Go version for go.mod and Dockerfile (default: 1.22)")
	newCmd.Flags().BoolVar(&noDockerFlag, "no-docker", false, "skip Docker + docker-compose files")
	newCmd.Flags().BoolVar(&noCIFlag, "no-ci", false, "skip GitHub Actions CI workflow")
	newCmd.Flags().StringVarP(&templateFlag, "template", "t", "", "path to a custom scaffold.yaml")
	newCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "output directory (default: service name)")
	newCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "print files that would be generated without writing them")
	newCmd.Flags().BoolVarP(&yesFlag, "yes", "y", false, "skip interactive prompts and use flag values + defaults")
}

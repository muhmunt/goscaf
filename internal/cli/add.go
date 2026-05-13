package cli

import (
	"fmt"

	"github.com/muhmunt/goscaf/internal/engine"
	"github.com/muhmunt/goscaf/internal/wizard"
	"github.com/spf13/cobra"
)

var (
	addNameFlag      string
	addModuleFlag    string
	addTransportFlag string
	addDryRunFlag    bool
	addYesFlag       bool
)

var addCmd = &cobra.Command{
	Use:   "add <layer>",
	Short: "Add a layer to an existing service (handler|usecase|repository)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r := &wizard.Result{
			ServiceName: addNameFlag,
			ModulePath:  addModuleFlag,
			Transport:   addTransportFlag,
			DryRun:      addDryRunFlag,
			Yes:         addYesFlag,
		}

		result, err := wizard.RunAdd(args[0], r)
		if err != nil {
			return fmt.Errorf("wizard: %w", err)
		}
		return engine.Run(result)
	},
}

func init() {
	addCmd.Flags().StringVar(&addNameFlag, "name", "", "resource name (e.g. user, order)")
	addCmd.Flags().StringVar(&addModuleFlag, "module", "", "Go module path of the existing service")
	addCmd.Flags().StringVar(&addTransportFlag, "transport", "", "transport: rest | grpc | both (handler layer only)")
	addCmd.Flags().BoolVar(&addDryRunFlag, "dry-run", false, "print files that would be generated without writing them")
	addCmd.Flags().BoolVarP(&addYesFlag, "yes", "y", false, "skip interactive prompts and use flag values + defaults")
}

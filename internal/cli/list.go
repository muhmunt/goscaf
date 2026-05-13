package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available presets and layers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Presets  (goscaf new)")
		fmt.Println("  rest    REST microservice — Gin router, Clean Architecture")
		fmt.Println("  grpc    gRPC microservice — buf for protobuf codegen")
		fmt.Println()
		fmt.Println("Layers   (goscaf add <layer>)")
		fmt.Println("  handler       HTTP handler, gRPC handler, or both (--transport rest|grpc|both)")
		fmt.Println("  usecase       Business logic wired to repository")
		fmt.Println("  repository    Data access stub ready for DB injection")
		fmt.Println()
		fmt.Println("Flags available on all commands:")
		fmt.Println("  --dry-run    Print files that would be generated without writing them")
		fmt.Println("  --yes, -y    Skip prompts — use flag values and defaults (non-interactive)")
	},
}

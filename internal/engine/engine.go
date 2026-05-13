package engine

import (
	"fmt"

	"github.com/muhmunt/goscaf/internal/wizard"
)

// Run orchestrates: load scaffold → render templates → write files → run hooks.
func Run(r *wizard.Result) error {
	scaffold, fsys, err := load(r)
	if err != nil {
		return fmt.Errorf("load scaffold: %w", err)
	}

	vars := buildVars(r)

	if r.DryRun {
		fmt.Printf("\n[dry-run] %q — files that would be written to ./%s\n\n", scaffold.Name, r.OutputDir)
	} else {
		fmt.Printf("\nScaffolding %q into ./%s\n\n", scaffold.Name, r.OutputDir)
	}

	if !r.DryRun {
		if err := runHooks(scaffold.Hooks.Pre, vars); err != nil {
			return fmt.Errorf("pre-hooks: %w", err)
		}
	}

	for _, f := range scaffold.Files {
		ok, err := evalCondition(f.Condition, vars)
		if err != nil {
			return fmt.Errorf("eval condition for %q: %w", f.Src, err)
		}
		if !ok {
			continue
		}

		dstPath, err := renderString(f.Dst, vars)
		if err != nil {
			return fmt.Errorf("render dst path %q: %w", f.Dst, err)
		}

		outPath := r.OutputDir + "/" + dstPath

		if r.DryRun {
			fmt.Printf("  ~ %s\n", outPath)
			continue
		}

		content, err := renderTemplate(fsys, f.Src, vars)
		if err != nil {
			return fmt.Errorf("render template %q: %w", f.Src, err)
		}

		if err := writeFile(outPath, content); err != nil {
			return fmt.Errorf("write %q: %w", outPath, err)
		}

		fmt.Printf("  ✓ %s\n", outPath)
	}

	if r.DryRun {
		fmt.Println("\nDry run complete — no files written.")
		return nil
	}

	if err := runHooks(scaffold.Hooks.Post, vars); err != nil {
		return fmt.Errorf("post-hooks: %w", err)
	}

	if r.Layer != "" {
		fmt.Printf("\nDone! %s layer added for %q.\n", r.Layer, r.ServiceName)
	} else {
		fmt.Printf("\nDone! Next steps:\n  cd %s\n  go mod tidy\n  go run ./cmd/%s\n",
			r.OutputDir, r.ServiceName)
	}
	return nil
}

func buildVars(r *wizard.Result) map[string]any {
	return map[string]any{
		"service_name":   r.ServiceName,
		"resource_name":  r.ServiceName,
		"module_path":    r.ModulePath,
		"transport":      r.Transport,
		"database":       r.Database,
		"architecture":   r.Architecture,
		"go_version":     r.GoVersion,
		"include_docker": r.IncludeDocker,
		"include_ci":     r.IncludeCI,
	}
}

package engine

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/muhmunt/goscaf/internal/schema"
)

func writeFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create directories: %w", err)
	}
	return os.WriteFile(path, content, 0o644)
}

func runHooks(hooks []schema.Hook, vars map[string]any) error {
	for _, h := range hooks {
		ok, err := evalCondition(h.Condition, vars)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}

		cmd, err := renderString(h.Cmd, vars)
		if err != nil {
			return err
		}

		fmt.Printf("  $ %s\n", cmd)

		parts := strings.Fields(cmd)
		if len(parts) == 0 {
			continue
		}

		c := exec.Command(parts[0], parts[1:]...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("hook %q failed: %w", cmd, err)
		}
	}
	return nil
}

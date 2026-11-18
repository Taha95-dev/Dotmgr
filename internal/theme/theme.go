package theme

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"dotmgr/internal/config"
)

func Apply(themeDir string) error {
	yamlPath := filepath.Join(themeDir, "theme.yaml")
	t, err := config.LoadTheme(yamlPath)
	if err != nil {
		return err
	}

	// 1. Native Go File Operations
	for _, file := range t.Files {
		src := filepath.Join(themeDir, file.Source)
		dst := config.ExpandHome(file.Target)

		os.MkdirAll(filepath.Dir(dst), 0755)

		if file.Action == "copy" {
			if err := copyFile(src, dst); err != nil {
				return fmt.Errorf("copy failed for %s: %w", file.Source, err)
			}
			fmt.Printf("  -> Copied %s\n", file.Source)
		} else {
			os.Remove(dst)
			if err := os.Symlink(src, dst); err != nil {
				return fmt.Errorf("symlink failed for %s: %w", file.Source, err)
			}
			fmt.Printf("  -> Linked %s\n", file.Source)
		}
	}

	// 2. Native Commands
	if len(t.Commands) > 0 {
		script := strings.Join(t.Commands, " && ")
		cmd := exec.Command("bash", "-c", script)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("command sequence failed: %w", err)
		}
	}

	// 3. Custom Script Escape Hatch
	if t.Script != "" {
		shell := t.Shell
		if shell == "" { shell = "bash" }
		scriptPath := config.ExpandHome(t.Script)

		cmd := exec.Command(shell, scriptPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("custom script failed: %w", err)
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil { return err }
	return os.WriteFile(dst, data, 0644)
}

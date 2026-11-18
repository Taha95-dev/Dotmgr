package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"dotmgr/internal/config"
)

func Execute(configPath string) error {
	cfg, err := config.LoadCommands(configPath)
	if err != nil {
		return err
	}

	// 1. External Script Mode
	if cfg.Script != "" {
		shell := cfg.Shell
		if shell == "" { shell = "bash" }

		scriptPath := config.ExpandHome(cfg.Script)
		fmt.Printf("🏃 Running custom script: %s (using %s)\n", scriptPath, shell)

		cmd := exec.Command(shell, scriptPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("script execution failed: %w", err)
		}
		return nil
	}

	// 2. Inline Commands Mode
	if len(cfg.Commands) > 0 {
		fmt.Println("🏃 Running inline commands...")
		fullScript := strings.Join(cfg.Commands, " && ")

		cmd := exec.Command("bash", "-c", fullScript)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("command sequence failed: %w", err)
		}
		return nil
	}

	fmt.Println("  [~] No commands or scripts defined.")
	return nil
}

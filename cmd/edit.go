package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"dotmgr/internal/config"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [symlinks|deps|config]",
	Short: "Edit dotmgr configuration files",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "nano" // fallback
		}

		var targetFile string
		switch args[0] {
		case "symlinks":
			targetFile = config.ExpandHome("~/.config/dotmgr/symlinks/symlink.yaml")
		case "deps":
			targetFile = config.ExpandHome("~/.config/dotmgr/deps/deps.yaml")
		case "config":
			targetFile = config.ExpandHome("~/.config/dotmgr/config.yaml")
		default:
			fmt.Printf("❌ Unknown config file: %s\n", args[0])
			return
		}

		cmdExec := exec.Command(editor, targetFile)
		cmdExec.Stdin = os.Stdin
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr

		if err := cmdExec.Run(); err != nil {
			fmt.Printf("❌ Failed to open %s: %v\n", editor, err)
		}
	},
}

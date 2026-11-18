package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"dotmgr/internal/symlink"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if your symlinks are healthy",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()
		configPath := filepath.Join(home, ".config", "dotmgr", "symlinks", "symlink.yaml")

		fmt.Println("🩺 Checking symlink health...")
		symlink.Status(configPath)
	},
}

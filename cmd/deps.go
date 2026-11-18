package cmd

import (
	"fmt"

	"dotmgr/internal/config"
	"dotmgr/internal/deps"

	"github.com/spf13/cobra"
)

var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Install system dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.ExpandHome("~/.config/dotmgr/deps/deps.yaml")

		fmt.Println("📦 Checking dependencies...")
		if err := deps.Install(configPath); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println("✅ Dependencies are up to date!")
	},
}

package cmd

import (
	"fmt"
	"dotmgr/internal/config"
	"dotmgr/internal/runner"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run custom system setup commands",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.ExpandHome("~/.config/dotmgr/commands/commands.yaml")

		if err := runner.Execute(configPath); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println("✅ Done!")
	},
}

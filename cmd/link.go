package cmd

import (
	"fmt"

	"dotmgr/internal/config"
	"dotmgr/internal/symlink"

	"github.com/spf13/cobra"
)

var dryRunLink bool

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Symlink your core dotfiles",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.ExpandHome("~/.config/dotmgr/symlinks/symlink.yaml")
		fmt.Println("🔗 Linking core dotfiles...")
		if err := symlink.Apply(configPath, dryRunLink); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println("✅ Done!")
	},
}

func init() {
	linkCmd.Flags().BoolVar(&dryRunLink, "dry-run", false, "Show what would be linked without making changes")
}

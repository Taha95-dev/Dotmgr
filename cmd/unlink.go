package cmd

import (
	"fmt"
	"dotmgr/internal/config"
	"dotmgr/internal/symlink"
	"github.com/spf13/cobra"
)

var dryRunUnlink bool

var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Remove symlinks and restore backups",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.ExpandHome("~/.config/dotmgr/symlinks/symlink.yaml")
		fmt.Println("🔓 Unlinking dotfiles and restoring backups...")

		if err := symlink.Unlink(configPath, dryRunUnlink); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println("✅ Done!")
	},
}

func init() {
	unlinkCmd.Flags().BoolVar(&dryRunUnlink, "dry-run", false, "Show what would be unlinked without making changes")
}

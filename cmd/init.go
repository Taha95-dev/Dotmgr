package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"dotmgr/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the dotmgr directory structure",
	Run: func(cmd *cobra.Command, args []string) {
		baseDir := config.ExpandHome("~/.config/dotmgr")
		symlinksDir := filepath.Join(baseDir, "symlinks")
		themesDir := filepath.Join(baseDir, "themes")
		depsDir := filepath.Join(baseDir, "deps")

		// 1. Create directories
		for _, dir := range []string{symlinksDir, themesDir, depsDir} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("❌ Failed to create %s: %v\n", dir, err)
				return
			}
		}
		fmt.Println("  [+] Created directory structure")

		// 2. Create default symlink.yaml
		yamlPath := filepath.Join(symlinksDir, "symlink.yaml")
		if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
			defaultYaml := `base_dir: "~/.config/dotmgr/symlinks"
links: []
`
			os.WriteFile(yamlPath, []byte(defaultYaml), 0644)
			fmt.Println("  [+] Created default symlink.yaml")
		}

		// 3. Create default deps.yaml
		depsPath := filepath.Join(depsDir, "deps.yaml")
		if _, err := os.Stat(depsPath); os.IsNotExist(err) {
			defaultDeps := `deps:
  - name: "git"
  - name: "fzf"
`
			os.WriteFile(depsPath, []byte(defaultDeps), 0644)
			fmt.Println("  [+] Created default deps.yaml")
		}

		fmt.Println("✅ dotmgr initialized! Next: run 'dotmgr link' or 'dotmgr theme'.")
	},
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dotmgr/internal/config"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [path-to-file]",
	Short: "Add a file to the dotfiles repo and symlink it",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetPath := args[0]
		home, _ := os.UserHomeDir()

		// 1. Expand ~ to the actual home path
		targetPath = config.ExpandHome(targetPath)

		// 2. Verify the file exists
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			fmt.Printf("❌ File not found: %s\n", targetPath)
			return
		}

		// 3. Determine the relative path for the repo
		// e.g., /home/taha/.zshrc becomes .zshrc
		relPath, err := filepath.Rel(home, targetPath)
		if err != nil {
			fmt.Printf("❌ File must be inside your home directory.\n")
			return
		}

		// 4. Copy the file to the symlinks directory
		symlinksBase := config.ExpandHome("~/.config/dotmgr/symlinks")
		repoSource := filepath.Join(symlinksBase, relPath)

		os.MkdirAll(filepath.Dir(repoSource), 0755)
		if err := copyFile(targetPath, repoSource); err != nil {
			fmt.Printf("❌ Failed to copy file: %v\n", err)
			return
		}
		fmt.Printf("  [+] Copied %s to repo\n", relPath)

		// 5. Delete the original file and create the symlink
		os.Remove(targetPath)
		if err := os.Symlink(repoSource, targetPath); err != nil {
			fmt.Printf("❌ Failed to create symlink: %v\n", err)
			return
		}
		fmt.Printf("  [+] Symlinked %s\n", targetPath)

		// 6. Append to symlink.yaml
		// In a real app, you'd parse and re-marshal the YAML, but for QoL,
		// we can just append text safely.
		yamlPath := filepath.Join(symlinksBase, "symlink.yaml")
		appendYaml(yamlPath, relPath)

		fmt.Println("✅ Successfully added to dotfiles!")
	},
}

func appendYaml(path, relPath string) {
	// Read existing content
	data, _ := os.ReadFile(path)
	content := string(data)

	// If it ends with "links: []\n", remove that so we can add an item
	if strings.Contains(content, "links: []") {
		content = strings.Replace(content, "links: []\n", "links:\n", 1)
	}

	// Build the new YAML entry
	newEntry := fmt.Sprintf("  - source: \"%s\"\n    target: \"~/%s\"\n    type: \"file\"\n", relPath, relPath)

	os.WriteFile(path, []byte(content+newEntry), 0644)
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil { return err }
	return os.WriteFile(dst, data, 0644)
}

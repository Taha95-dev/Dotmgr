package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dotmgr/internal/config"
	"dotmgr/internal/fzf"
	"dotmgr/internal/theme"
	"dotmgr/internal/wm"

	"github.com/spf13/cobra"
)

var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Switch themes using FZF",
	Run: func(cmd *cobra.Command, args []string) {
		currentWM := wm.Detect()
		themesDir := config.ExpandHome("~/.config/dotmgr/themes")

		entries, err := os.ReadDir(themesDir)
		if err != nil {
			fmt.Println("❌ Could not read themes directory:", err)
			return
		}

		var validThemes []string
		for _, e := range entries {
			if !e.IsDir() { continue }

			yamlPath := filepath.Join(themesDir, e.Name(), "theme.yaml")
			t, err := config.LoadTheme(yamlPath)
			if err != nil { continue }

			// Only show themes that match your current Window Manager
			if strings.ToLower(t.WM) == currentWM {
				validThemes = append(validThemes, e.Name())
			}
		}

		if len(validThemes) == 0 {
			fmt.Printf("No themes found for WM: %s\n", currentWM)
			return
		}

		selected, _ := fzf.Select(validThemes)
		if selected == "" { return }

		fmt.Printf("🎨 Applying theme: %s\n", selected)
		themePath := filepath.Join(themesDir, selected)
		if err := theme.Apply(themePath); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		}
	},
}

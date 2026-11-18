package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configDir string

var rootCmd = &cobra.Command{
	Use:   "dotmgr",
	Short: "A fast, YAML-driven dotfile and theme manager.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Add the global flag here
	rootCmd.PersistentFlags().StringVar(&configDir, "config", "", "Path to dotmgr config directory (default is ~/.config/dotmgr)")

	rootCmd.AddCommand(linkCmd)
	rootCmd.AddCommand(themeCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(depsCmd)
	rootCmd.AddCommand(unlinkCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(doctorCmd)
  rootCmd.AddCommand(runCmd)
}

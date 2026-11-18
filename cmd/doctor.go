package cmd

import (
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check for required system tools",
	Run: func(cmd *cobra.Command, args []string) {
		// Tools that dotmgr itself requires
		requiredTools := []string{"git", "fzf"}
		// Tools that your themes likely require
		optionalTools := []string{"starship", "zoxide", "eza", "bat"}

		fmt.Println("🩺 Running environment diagnostics...")

		fmt.Println("\n  Required Tools:")
		for _, tool := range requiredTools {
			if _, err := exec.LookPath(tool); err == nil {
				fmt.Printf("    [✅] %s is installed\n", tool)
			} else {
				fmt.Printf("    [❌] %s is MISSING (dotmgr needs this!)\n", tool)
			}
		}

		fmt.Println("\n  Optional Tools (for themes):")
		for _, tool := range optionalTools {
			if _, err := exec.LookPath(tool); err == nil {
				fmt.Printf("    [✅] %s is installed\n", tool)
			} else {
				fmt.Printf("    [⚠️ ] %s is missing\n", tool)
			}
		}
	},
}

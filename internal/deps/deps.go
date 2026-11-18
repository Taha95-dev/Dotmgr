package deps

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"dotmgr/internal/config"
)

// detectManager checks which package manager is available
func detectManager(preferred string) (string, error) {
	if preferred != "" {
		if _, err := exec.LookPath(preferred); err == nil {
			return preferred, nil
		}
		fmt.Printf("  [!] Specified manager '%s' not found, falling back...\n", preferred)
	}

	if _, err := exec.LookPath("pacman"); err == nil { return "pacman", nil }
	if _, err := exec.LookPath("apt"); err == nil { return "apt", nil }
	if _, err := exec.LookPath("brew"); err == nil { return "brew", nil }

	return "", fmt.Errorf("could not detect a supported package manager")
}

// Install reads the deps.yaml and installs missing packages in a batch
func Install(configPath string) error {
	cfg, err := config.LoadDeps(configPath)
	if err != nil {
		return err
	}

	// Map of manager -> list of missing packages
	missingPackages := make(map[string][]string)

	for _, dep := range cfg.Deps {
		manager, err := detectManager(dep.Manager)
		if err != nil {
			return err
		}

		// 1. Check if the package is already installed
		isInstalled := false
		switch manager {
		case "pacman":
			if err := exec.Command("pacman", "-Q", dep.Name).Run(); err == nil {
				isInstalled = true
			}
		case "apt":
			if err := exec.Command("dpkg", "-s", dep.Name).Run(); err == nil {
				isInstalled = true
			}
		case "brew":
			if err := exec.Command("brew", "list", dep.Name).Run(); err == nil {
				isInstalled = true
			}
		}

		if isInstalled {
			fmt.Printf("  [~] Already installed: %s\n", dep.Name)
			continue
		}

		// 2. Add to the missing list for this manager
		missingPackages[manager] = append(missingPackages[manager], dep.Name)
	}

	// 3. Install all missing packages in batches
	for manager, packages := range missingPackages {
		if len(packages) == 0 {
			continue
		}

		fmt.Printf("  [+] Installing %d packages with %s...\n", len(packages), manager)

		var cmd *exec.Cmd
		switch manager {
		case "pacman":
			args := append([]string{"pacman", "-S", "--needed", "--noconfirm"}, packages...)
			cmd = exec.Command("sudo", args...)
		case "apt":
			args := append([]string{"apt", "install", "-y"}, packages...)
			cmd = exec.Command("sudo", args...)
		case "brew":
			args := append([]string{"install"}, packages...)
			cmd = exec.Command("brew", args...)
		}

		// Print the exact command being run (for transparency)
		fmt.Printf("      Running: %s\n", strings.Join(cmd.Args, " "))

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			fmt.Printf("  ❌ Failed to install batch: %v\n", err)
		}
	}

	return nil
}

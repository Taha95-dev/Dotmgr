package symlink

import (
	"fmt"
	"os"
	"path/filepath"
	"dotmgr/internal/config"
)

func Apply(configPath string, dryRun bool) error {
	cfg, err := config.LoadSymlinks(configPath)
	if err != nil {
		return err
	}
	baseDir := config.ExpandHome(cfg.BaseDir)

	for _, link := range cfg.Links {
		src := filepath.Join(baseDir, link.Source)
		dst := config.ExpandHome(link.Target)

		if _, err := os.Stat(src); os.IsNotExist(err) {
			fmt.Printf("  ❌ Source does not exist: %s\n", link.Source)
			continue
		}

		// If it's a directory, we ALWAYS use the ignore logic
		if link.Type == "dir" {
			// Combine the user's custom ignores with our default ignores
			allIgnores := append(config.DefaultIgnores(), link.Ignore...)
			if err := linkDirectoryWithIgnores(src, dst, allIgnores, dryRun); err != nil {
				return err
			}
			continue
		}

		// Normal file symlink
		if dryRun {
			fmt.Printf("  [Dry Run] Would link %s -> %s\n", src, dst)
			continue
		}

		if err := createSymlink(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func linkDirectoryWithIgnores(srcDir, dstDir string, ignores []string, dryRun bool) error {
	// 1. Ensure the destination directory exists (DO NOT back up an existing directory)
	if !dryRun {
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			return fmt.Errorf("failed to create target dir %s: %w", dstDir, err)
		}
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read source dir %s: %w", srcDir, err)
	}

	for _, entry := range entries {
		name := entry.Name()

		// 2. Check if this entry is in the ignore list
		ignored := false
		for _, ig := range ignores {
			// Use filepath.Match to support wildcards like *.so
			matched, _ := filepath.Match(ig, name)
			if matched {
				ignored = true
				break
			}
		}

		if ignored {
			fmt.Printf("  [~] Ignored (leaving local): %s\n", name)
			continue
		}

		// 3. Link the non-ignored files
		srcPath := filepath.Join(srcDir, name)
		dstPath := filepath.Join(dstDir, name)

		if dryRun {
			fmt.Printf("  [Dry Run] Would link %s -> %s\n", srcPath, dstPath)
		} else {
			if err := createSymlink(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func createSymlink(src, dst string) error {
	// Backup existing if it's not our symlink
	if _, err := os.Lstat(dst); err == nil {
		existingTarget, _ := os.Readlink(dst)
		if existingTarget == src {
			fmt.Printf("  [~] Already linked: %s\n", dst)
			return nil
		}

		backupPath := dst + ".bak"
		fmt.Printf("  [!] Backing up %s -> %s\n", dst, backupPath)
		os.RemoveAll(backupPath)
		if err := os.Rename(dst, backupPath); err != nil {
			return fmt.Errorf("failed to backup %s: %w", dst, err)
		}
	}

	if err := os.Symlink(src, dst); err != nil {
		return fmt.Errorf("failed to link %s -> %s: %w", src, dst, err)
	}
	fmt.Printf("  [+] Linked: %s\n", dst)
	return nil
}

func Status(configPath string) error {
	cfg, err := config.LoadSymlinks(configPath)
	if err != nil {
		return err
	}
	baseDir := config.ExpandHome(cfg.BaseDir)

	broken := 0
	for _, link := range cfg.Links {
		src := filepath.Join(baseDir, link.Source)
		dst := config.ExpandHome(link.Target)

		info, err := os.Lstat(dst)
		if err != nil {
			fmt.Printf("  [X] Missing: %s\n", link.Target)
			broken++
			continue
		}

		if info.Mode()&os.ModeSymlink == 0 {
			fmt.Printf("  [!] Not a symlink: %s\n", link.Target)
			broken++
			continue
		}

		existingTarget, _ := os.Readlink(dst)
		if existingTarget != src {
			fmt.Printf("  [!] Wrong target: %s\n", link.Target)
			broken++
		}
	}

	if broken == 0 {
		fmt.Println("  ✅ All links are healthy.")
	}
	return nil
}

func Unlink(configPath string, dryRun bool) error {
	cfg, err := config.LoadSymlinks(configPath)
	if err != nil {
		return err
	}

	baseDir := config.ExpandHome(cfg.BaseDir)

	for _, link := range cfg.Links {
		dst := config.ExpandHome(link.Target)
		backupPath := dst + ".bak"
		expectedSrc := filepath.Join(baseDir, link.Source)

		// 1. Check if the target is a symlink
		info, err := os.Lstat(dst)
		if err != nil {
			fmt.Printf("  [~] Not found: %s\n", link.Target)
			continue
		}

		if info.Mode()&os.ModeSymlink == 0 {
			fmt.Printf("  [!] Not a symlink, skipping: %s\n", link.Target)
			continue
		}

		// 2. SAFETY CHECK: Does it point to our dotmgr directory?
		actualSrc, err := os.Readlink(dst)
		if err != nil || actualSrc != expectedSrc {
			fmt.Printf("  [!] Warning: %s does not point to dotmgr. Skipping to be safe.\n", link.Target)
			continue
		}

		// 3. Dry Run Check
		if dryRun {
			fmt.Printf("  [Dry Run] Would remove: %s\n", link.Target)
			if _, err := os.Stat(backupPath); err == nil {
				fmt.Printf("  [Dry Run] Would restore backup: %s\n", backupPath)
			}
			continue
		}

		// 4. Remove the symlink
		if err := os.Remove(dst); err != nil {
			fmt.Printf("  ❌ Failed to remove symlink: %s\n", link.Target)
			continue
		}
		fmt.Printf("  [-] Unlinked: %s\n", link.Target)

		// 5. Restore backup if it exists
		if _, err := os.Stat(backupPath); err == nil {
			if err := os.Rename(backupPath, dst); err != nil {
				fmt.Printf("  ❌ Failed to restore backup: %s\n", backupPath)
			} else {
				fmt.Printf("  [+] Restored backup: %s\n", link.Target)
			}
		}
	}
	return nil
}

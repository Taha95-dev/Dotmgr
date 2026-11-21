package symlink

import (
	"fmt"
	"io"
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

		// Handle Directories
		if link.Type == "dir" {
			allIgnores := append(config.DefaultIgnores(), link.Ignore...)
			if err := linkDirectory(src, dst, link.Action, allIgnores, dryRun); err != nil {
				return err
			}
			continue
		}

		// Handle Single Files
		if dryRun {
			fmt.Printf("  [Dry Run] Would %s: %s -> %s\n", link.Action, src, dst)
			continue
		}

		if link.Action == "copy" {
			if err := copyFile(src, dst); err != nil {
				return err
			}
			fmt.Printf("  [+] Copied: %s\n", dst)
		} else {
			if err := createSymlink(src, dst); err != nil {
				return err
			}
		}
	}
	return nil
}

func linkDirectory(srcDir, dstDir, action string, ignores []string, dryRun bool) error {
	if !dryRun {
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			return fmt.Errorf("failed to create target dir %s: %w", dstDir, err)
		}
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()

		// Check ignores
		ignored := false
		for _, ig := range ignores {
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

		srcPath := filepath.Join(srcDir, name)
		dstPath := filepath.Join(dstDir, name)

		if dryRun {
			fmt.Printf("  [Dry Run] Would %s: %s -> %s\n", action, srcPath, dstPath)
		} else {
			if entry.IsDir() {
				// Recursive directory copy (if action is copy)
				if action == "copy" {
					if err := copyDir(srcPath, dstPath); err != nil { return err }
					fmt.Printf("  [+] Copied dir: %s\n", dstPath)
				} else {
					if err := createSymlink(srcPath, dstPath); err != nil { return err }
				}
			} else {
				// File
				if action == "copy" {
					if err := copyFile(srcPath, dstPath); err != nil { return err }
					fmt.Printf("  [+] Copied: %s\n", dstPath)
				} else {
					if err := createSymlink(srcPath, dstPath); err != nil { return err }
				}
			}
		}
	}
	return nil
}

func createSymlink(src, dst string) error {
	// Backup existing
	if _, err := os.Lstat(dst); err == nil {
		existingTarget, _ := os.Readlink(dst)
		if existingTarget == src {
			fmt.Printf("  [~] Already linked: %s\n", dst)
			return nil
		}
		os.RemoveAll(dst + ".bak")
		os.Rename(dst, dst+".bak")
	}

	if err := os.Symlink(src, dst); err != nil {
		return err
	}
	fmt.Printf("  [+] Linked: %s\n", dst)
	return nil
}

func copyFile(src, dst string) error {
	// Backup existing
	if _, err := os.Stat(dst); err == nil {
		os.RemoveAll(dst + ".bak")
		os.Rename(dst, dst+".bak")
	}

	in, err := os.Open(src)
	if err != nil { return err }
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil { return err }
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }

		relPath, err := filepath.Rel(src, path)
		if err != nil { return err }
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return copyFile(path, targetPath)
	})
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

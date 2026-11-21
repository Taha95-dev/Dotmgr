package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func ExpandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func DefaultIgnores() []string {
	return []string{
		".DS_Store",
		"*.so", // Wildcards are handled in the symlink engine
		"*.o",
	}
}

// ==========================================
// SYMLINK MODELS
// ==========================================

type Symlink struct {
	Source string   `yaml:"source"`
	Target string   `yaml:"target"`
	Type   string   `yaml:"type,omitempty"`
  Action string   `yaml:"action,omitempty"` // "symlink" (default) or "copy"
	Ignore []string `yaml:"ignore,omitempty"`
}

type SymlinkConfig struct {
	BaseDir string    `yaml:"base_dir"`
	Links   []Symlink `yaml:"links"`
}

func LoadSymlinks(path string) (*SymlinkConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read symlink config: %w", err)
	}
	var cfg SymlinkConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid symlink yaml: %w", err)
	}
	return &cfg, nil
}

// ==========================================
// THEME MODELS
// ==========================================

type ThemeFile struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
	Action string `yaml:"action,omitempty"`
}

type Theme struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	WM          string      `yaml:"wm"`
	Files       []ThemeFile `yaml:"files"`
	Commands    []string    `yaml:"commands,omitempty"`
	Script      string      `yaml:"script,omitempty"`
	Shell       string      `yaml:"shell,omitempty"`
}

func LoadTheme(path string) (*Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read theme file: %w", err)
	}
	var t Theme
	if err := yaml.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("invalid theme yaml: %w", err)
	}
	return &t, nil
}

// ==========================================
// DEPENDENCY MODELS
// ==========================================

type Dependency struct {
	Name    string `yaml:"name"`
	Manager string `yaml:"manager,omitempty"`
}

type DepsConfig struct {
	Deps []Dependency `yaml:"deps"`
}

func LoadDeps(path string) (*DepsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read deps config: %w", err)
	}
	var cfg DepsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid deps yaml: %w", err)
	}
	return &cfg, nil
}

// ==========================================
// COMMANDS MODELS (NEW)
// ==========================================

type CommandsConfig struct {
	Commands []string `yaml:"commands,omitempty"`
	Script   string   `yaml:"script,omitempty"`
	Shell    string   `yaml:"shell,omitempty"`
}

func LoadCommands(path string) (*CommandsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read commands config: %w", err)
	}
	var cfg CommandsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid commands yaml: %w", err)
	}
	return &cfg, nil
}

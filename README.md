# dotmgr

A fast, YAML-driven dotfile, theme, and environment manager written in Go.

`dotmgr` replaces messy `bootstrap.sh` scripts with clean, declarative YAML. It handles symlinking, theme switching, package installation, and custom setup scripts, all from a single binary.

## Why `dotmgr`?
Most dotfile managers just create symlinks. `dotmgr` manages your entire environment lifecycle:
1. **Declarative:** Define your environment in YAML, not Bash.
2. **Safe:** Automatically backs up existing files before overwriting.
3. **Modular:** Separate your core configs (`symlinks`), aesthetics (`themes`), packages (`deps`), and setup logic (`commands`).

## Installation

```bash
git clone https://github.com/Taha95-dev/dotmgr.git
cd dotmgr
go build -o dotmgr .
sudo install -m 755 dotmgr /usr/local/bin/dotmgr
```
# Quick Start

Initialize the directory structure:
```bash
dotmgr init
```
This creates ~/.config/dotmgr/ with the following structure
```text

~/.config/dotmgr/
├── commands/
│   └── commands.yaml    # Custom setup scripts (e.g., git config, neovim sync)
├── deps/
│   └── deps.yaml        # System packages (pacman, apt, brew)
├── symlinks/
│   └── symlink.yaml     # Core dotfiles to link
└── themes/
    └── my_theme/
        └── theme.yaml   # Window manager themes (dwl, hyprland, etc.)
```
YAML Syntax Guide
1. Symlinks (symlink.yaml)

Map files from your dotmgr repo to your home directory.
```yaml

base_dir: "~/.config/dotmgr/symlinks"

links:
  # Link a single file
  - source: "zsh/.zshrc"
    target: "~/.zshrc"
    type: "file"

  # Link a directory (with automatic ignores)
  - source: "nvim"
    target: "~/.config/nvim"
    type: "dir"
    ignore:
      - "lazy-lock.json"
      - "parser"
      - "*.so"
```
Note: dotmgr automatically ignores .git, .gitignore, .DS_Store, *.so, and *.o by default.

2. Dependencies (deps.yaml)

List your required packages. dotmgr will auto-detect your package manager (pacman, apt, or brew) and batch install them in one command.
```yaml

deps:
  - name: "neovim"
    manager: "pacman"
  - name: "fzf"
    manager: "pacman"
  - name: "starship"
```
3. Themes (theme.yaml)

Themes are tied to a specific Window Manager (wm). dotmgr theme will only show you themes that match your current WM.
```yaml

name: "Ember"
description: "Warm dark theme"
wm: "dwl"

files:
  - source: "alacritty.toml"
    target: "~/.config/alacritty/alacritty.toml"
    action: "symlink" # or "copy"

# Inline commands to run after copying files
commands:
  - "cd ~/dwl && make clean && make && sudo make install"
  - "pkill waybar; waybar &"

# OR: Point to a custom shell script (Power Mode)
script: "~/.config/dotmgr/commands/reload_dwl.sh"
shell: "bash"
```
4. Custom Commands (commands.yaml)

For system setup that doesn't fit into themes (like authenticating GitHub or changing your shell). You can use inline commands or point to a script.

# Inline:
```yaml

commands:
  - "git config --global user.name 'Taha'"
  - "gsettings set org.gnome.desktop.interface color-scheme 'prefer-dark'"
```
# External Script (Recommended for complex logic):
```yaml

script: "~/.config/dotmgr/commands/setup.sh"
shell: "bash"
```
| Commands Reference | Command Description |
| ------------------ | -------------------- |
| dotmgr init |	Create the default ~/.config/dotmgr structure. |
| dotmgr link |	Symlink core dotfiles. Use --dry-run to preview. |
| dotmgr unlink |	Remove symlinks and restore .bak files safely. |
| dotmgr status |	Check if your symlinks are healthy and pointing to the right place. |
| dotmgr add <file> |	Move a file into the dotmgr repo and symlink it automatically. |
| dotmgr theme |	Launch an FZF menu to switch themes (WM filtered). |
| dotmgr deps |	Install missing system packages in a single batch. |
| dotmgr run |	Run custom system setup scripts. |
| dotmgr doctor |	Check for required tools (git, fzf) and optional theme tools. |

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

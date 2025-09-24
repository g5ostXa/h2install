package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var packages = []string{
	"hyprland",
	"uwsm",
	"alacritty",
	"aquamarine",
	"waybar",
	"rofi",
	"libnotify",
	"dunst",
	"cliphist",
	"wlogout",
	"xdg-desktop-portal-hyprland",
	"xdg-desktop-portal-gtk",
	"qt5-wayland",
	"qt6-wayland",
	"waypaper",
	"hyprpicker",
	"hyprlock",
	"hyprcursor",
	"hypridle",
	"hyprgraphics",
	"hyprland-qt-support",
	"hyprland-qtutils",
	"hyprlang",
	"hyprls-git",
	"hyprwayland-scanner",
	"otf-font-awesome",
	"woff2-font-awesome",
	"ttf-fira-sans",
	"ttf-fira-code",
	"ttf-firacode-nerd",
	"brightnessctl",
	"neovim",
	"nautilus",
	"fastfetch",
	"pavucontrol",
	"pulseaudio",
	"bibata-cursor-theme",
	"dracula-icons-theme",
	"tokyonight-gtk-theme-git",
	"python-pywal16",
	"gtk2",
	"gtk3",
	"gtk4",
	"nwg-look",
	"swww",
	"fish",
	"starship",
	"python-pip",
	"eza",
	"swappy",
	"vscodium-bin",
	"firefox",
}

var symlinks = [][2]string{
	{"dotfiles/gtk/.Xresources", ".Xresources"},
	{"dotfiles/alacritty", ".config/alacritty"},
	{"dotfiles/dunst", ".config/dunst"},
	{"dotfiles/gtk", ".config/gtk"},
	{"dotfiles/hypr", ".config/hypr"},
	{"dotfiles/nvim", ".config/nvim"},
	{"dotfiles/rofi", ".config/rofi"},
	{"dotfiles/starship/starship.toml", ".config/starship.toml"},
	{"dotfiles/swappy", ".config/swappy"},
	{"dotfiles/vim", ".config/vim"},
	{"dotfiles/wal", ".config/wal"},
	{"dotfiles/waybar", ".config/waybar"},
	{"dotfiles/wlogout", ".config/wlogout"},
	{"dotfiles/fastfetch", ".config/fastfetch"},
	{"dotfiles/fish", ".config/fish"},
	{"dotfiles/pacseek", ".config/pacseek"},
	{"dotfiles/waypaper", ".config/waypaper"},
	{"dotfiles/uwsm", ".config/uwsm"},
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Chmod(0644)
}

func installBashrc(home string) error {
	src := filepath.Join(home, "Downloads", "hyprarch2", ".bashrc")
	if _, err := os.Stat(src); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return copyFile(src, filepath.Join(home, ".bashrc"))
}

func installWallpaper(dest string) error {
	if err := os.RemoveAll(dest); err != nil {
		return err
	}
	const repo = "https://github.com/g5ostXa/wallpaper.git"
	return runCommand("git", "clone", "--depth=1", repo, dest)
}

func installPackages(helper string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	helperDir := filepath.Join(home, ".cache", helper+"-bin")
	if err := os.RemoveAll(helperDir); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(helperDir), 0o755); err != nil {
		return err
	}

	repo := "https://aur.archlinux.org/" + helper + "-bin.git"
	if err := runCommand("git", "clone", "--depth=1", repo, helperDir); err != nil {
		return err
	}
	defer os.RemoveAll(helperDir)

	prev, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(prev)

	if err := os.Chdir(helperDir); err != nil {
		return err
	}

	if err := runCommand("makepkg", "-si", "--noconfirm"); err != nil {
		return err
	}

	args := append([]string{"-S", "--needed", "--noconfirm"}, packages...)
	return runCommand(helper, args...)
}

func createSymlinks(home string) error {
	for _, pair := range symlinks {
		src := filepath.Join(home, filepath.FromSlash(pair[0]))
		dst := filepath.Join(home, filepath.FromSlash(pair[1]))

		if err := os.RemoveAll(dst); err != nil {
			return err
		}
		if err := os.Symlink(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func cleanup(home string) error {
	os.RemoveAll(filepath.Join(home, "Downloads", "hyprarch2"))

	script := filepath.Join(home, "src", "Scripts", "cleanup.sh")
	if _, err := os.Stat(script); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return runCommand("bash", script)
}

func runPacman(home string) error {
	script := filepath.Join(home, "Downloads", "hyprarch2", "src", "Scripts", "pacman.sh")
	if _, err := os.Stat(script); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return runCommand("bash", script)
}

func main() {
	helper := flag.String("aur-helper", "paru", "AUR helper to use")
	noWallpaper := flag.Bool("no-wallpaper", false, "skip wallpaper installation")
	flag.Parse()

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := runPacman(home); err != nil {
		fmt.Fprintf(os.Stderr, "pacman: %v\n", err)
		os.Exit(1)
	}

	if err := installPackages(*helper); err != nil {
		fmt.Fprintf(os.Stderr, "packages: %v\n", err)
		os.Exit(1)
	}

	if !*noWallpaper {
		if err := installWallpaper(filepath.Join(home, "wallpaper")); err != nil {
			fmt.Fprintf(os.Stderr, "wallpaper: %v\n", err)
			os.Exit(1)
		}
	}

	if err := installBashrc(home); err != nil {
		fmt.Fprintf(os.Stderr, "bashrc: %v\n", err)
		os.Exit(1)
	}

	if err := createSymlinks(home); err != nil {
		fmt.Fprintf(os.Stderr, "symlinks: %v\n", err)
		os.Exit(1)
	}

	if err := cleanup(home); err != nil {
		fmt.Fprintf(os.Stderr, "cleanup: %v\n", err)
		os.Exit(1)
	}
}

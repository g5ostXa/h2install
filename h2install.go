// ===== h2install.go =====
package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var packages = []string{
	"hyprland",
	"hyprpolkitagent",
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
	"waypaper",
	"hyprpicker",
	"hyprlock",
	"hyprcursor",
	"hypridle",
	"hyprgraphics",
	"hyprlang",
	"hyprls-git",
	"hyprwayland-scanner",
	"otf-font-awesome",
	"woff2-font-awesome",
	"ttf-fira-sans",
	"ttf-fira-code",
	"ttf-firacode-nerd",
	"gnu-free-fonts",
	"brightnessctl",
	"neovim",
	"nautilus",
	"fastfetch",
	"pavucontrol",
	"pipewire",
	"pipewire-pulse",
	"pipewire-alsa",
	"pipewire-jack",
	"wireplumber",
	"bibata-cursor-theme",
	"dracula-icons-theme",
	"tokyonight-gtk-theme-git",
	"python-pywal16",
	"gtk3",
	"gtk4",
	"swww",
	"fish",
	"starship",
	"python-pip",
	"eza",
	"swappy",
	"firefox-nightly-bin",
	"vscodium-bin",
	"ccache",
	"jq",
	"pacman-contrib",
	"fzf",
	"ttf-0xproto-nerd",
}

var symlinks = [][2]string{
	{"dotfiles/gtk/.Xresources", ".Xresources"},
	{"dotfiles/alacritty", ".config/alacritty"},
	{"dotfiles/dunst", ".config/dunst"},
	{"dotfiles/gtk/", ".config/gtk"},
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

func run(name string, args ...string) error {
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

	return out.Chmod(0o644)
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

func installWallpaper(home string) error {
	dest := filepath.Join(home, "wallpaper")
	if err := os.RemoveAll(dest); err != nil {
		return err
	}
	return run("git", "clone", "--depth=1", "https://github.com/g5ostXa/wallpaper.git", dest)
}

func installPackages(home, helper string) error {
	helperDir := filepath.Join(home, ".cache", helper+"-bin")
	if err := os.RemoveAll(helperDir); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(helperDir), 0o755); err != nil {
		return err
	}

	if err := run("git", "clone", "--depth=1", "https://aur.archlinux.org/"+helper+"-bin.git", helperDir); err != nil {
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

	if err := run("makepkg", "-si", "--noconfirm"); err != nil {
		return err
	}

	args := append([]string{"-S", "--needed", "--noconfirm"}, packages...)
	return run(helper, args...)
}

func createSymlinks(home string) error {
	for _, pair := range symlinks {
		src := filepath.Join(home, pair[0])
		dst := filepath.Join(home, pair[1])

		if err := os.RemoveAll(dst); err != nil {
			return err
		}
		if err := os.Symlink(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	script := filepath.Join(home, "Downloads", "hyprarch2", "src", "Scripts", "pacman.sh")
	if _, err := os.Stat(script); err == nil {
		if err := run("bash", script); err != nil {
			log.Fatal(err)
		}
	}

	if err := installPackages(home, "paru"); err != nil {
		log.Fatal(err)
	}
	if err := installWallpaper(home); err != nil {
		log.Fatal(err)
	}
	if err := installBashrc(home); err != nil {
		log.Fatal(err)
	}
	if err := createSymlinks(home); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ===== Package List =====
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
}

func runCommand(dryRun bool, name string, args ...string) error {
	fmt.Printf(">> %s %s\n", name, strings.Join(args, " "))
	if dryRun {
		return nil
	}
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func copyFile(src, dst string, dryRun bool) error {
	if dryRun {
		fmt.Printf("Would copy %s → %s\n", src, dst)
		return nil
	}
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
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Chmod(0644)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func printBanner(msg string) {
	if _, err := exec.LookPath("figlet"); err == nil {
		runCommand(false, "figlet", "-f", "smslant", msg)
	} else {
		fmt.Printf("=== %s ===\n", msg)
	}
}

func promptYesNo(msg string) bool {
	fmt.Printf("%s (y/n): ", msg)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}

func promptChoice(msg string, choices []string) string {
	fmt.Println(msg)
	for i, c := range choices {
		fmt.Printf("[%d] %s\n", i+1, c)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Choose option: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		for i, c := range choices {
			if input == fmt.Sprintf("%d", i+1) || strings.EqualFold(input, c) {
				return c
			}
		}
	}
}

// ===== Bashrc =====
func installBashrc(home string, dryRun bool) error {
	src := filepath.Join(home, "Downloads", "hyprarch2", ".bashrc")
	dst := filepath.Join(home, ".bashrc")

	if !fileExists(src) {
		fmt.Println(";; No .bashrc found in hyprarch2, skipping.")
		return nil
	}

	if dryRun {
		fmt.Printf("Would copy %s → %s\n", src, dst)
		return nil
	}

	if err := copyFile(src, dst, dryRun); err != nil {
		return fmt.Errorf("failed to install .bashrc: %w", err)
	}

	fmt.Println(";; Installed new .bashrc successfully.")
	return nil
}

// ===== Wallpaper =====
func installWallpaper(repoURL, dest string, force, dryRun bool) error {
	fmt.Println(";; Installing wallpapers...")
	if fileExists(dest) {
		if !force {
			fmt.Printf(";; Wallpaper dir exists (%s), skipping.\n", dest)
			return nil
		}
		if dryRun {
			fmt.Printf("Would remove existing wallpaper dir %s\n", dest)
		} else {
			os.RemoveAll(dest)
		}
	}
	return runCommand(dryRun, "git", "clone", "--depth=1", repoURL, dest)
}

// ===== AUR Helper + Packages =====
func installPackages(helper string, dryRun bool) error {
	home, _ := os.UserHomeDir()
	helperRepo := fmt.Sprintf("https://aur.archlinux.org/%s-bin.git", helper)
	helperDir := filepath.Join(home, ".cache", helper+"-bin")

	if dryRun {
		fmt.Printf("Would install AUR helper %s and packages: %v\n", helper, packages)
		return nil
	}

	os.MkdirAll(filepath.Dir(helperDir), 0755)
	if err := runCommand(false, "git", "clone", "--depth=1", helperRepo, helperDir); err != nil {
		return err
	}
	defer os.RemoveAll(helperDir)

	if err := os.Chdir(helperDir); err != nil {
		return err
	}
	if err := runCommand(false, "makepkg", "-si", "--noconfirm"); err != nil {
		return err
	}

	return runCommand(false, helper, append([]string{"-S", "--needed", "--noconfirm"}, packages...)...)
}

// ===== Symlinks =====
func createSymlinks(home string, force, dryRun bool) {
	links := map[string]string{
		"~/dotfiles/gtk/.Xresources":        "~/.Xresources",
		"~/dotfiles/alacritty":              "~/.config/alacritty",
		"~/dotfiles/dunst":                  "~/.config/dunst",
		"~/dotfiles/gtk":                    "~/.config/gtk",
		"~/dotfiles/hypr":                   "~/.config/hypr",
		"~/dotfiles/nvim":                   "~/.config/nvim",
		"~/dotfiles/rofi":                   "~/.config/rofi",
		"~/dotfiles/starship/starship.toml": "~/.config/starship.toml",
		"~/dotfiles/swappy":                 "~/.config/swappy",
		"~/dotfiles/vim":                    "~/.config/vim",
		"~/dotfiles/wal":                    "~/.config/wal",
		"~/dotfiles/waybar":                 "~/.config/waybar",
		"~/dotfiles/wlogout":                "~/.config/wlogout",
		"~/dotfiles/fastfetch":              "~/.config/fastfetch",
		"~/dotfiles/fish":                   "~/.config/fish",
		"~/dotfiles/pacseek":                "~/.config/pacseek",
		"~/dotfiles/waypaper":               "~/.config/waypaper",
		"~/dotfiles/uwsm":                   "~/.config/uwsm",
	}

	for src, dst := range links {
		srcPath := strings.Replace(src, "~", home, 1)
		dstPath := strings.Replace(dst, "~", home, 1)
		if dryRun {
			fmt.Printf("Would link %s → %s\n", srcPath, dstPath)
			continue
		}
		if fileExists(dstPath) {
			if !force {
				fmt.Printf(";; Skipping %s (exists)\n", dstPath)
				continue
			}
			os.RemoveAll(dstPath)
		}
		if err := os.Symlink(srcPath, dstPath); err != nil {
			fmt.Printf(";; Failed %s → %s: %v\n", srcPath, dstPath, err)
		} else {
			fmt.Printf(";; Linked %s → %s\n", srcPath, dstPath)
		}
	}
}

// ===== Cleanup =====
func cleanup(home string, skip, dryRun bool) {
	hyprarchDir := filepath.Join(home, "Downloads", "hyprarch2")
	if dryRun {
		fmt.Printf("Would remove %s\n", hyprarchDir)
		return
	}
	if fileExists(hyprarchDir) {
		os.RemoveAll(hyprarchDir)
		fmt.Printf(";; Removed %s\n", hyprarchDir)
	}
	if skip {
		return
	}
	cleanupScript := filepath.Join(home, "src", "Scripts", "cleanup.sh")
	if fileExists(cleanupScript) {
		runCommand(false, "bash", cleanupScript)
	}
}

// ===== Main =====
func main() {
	dryRun := flag.Bool("dry-run", false, "simulate actions without executing")
	force := flag.Bool("force", false, "force overwrite of files/symlinks")
	skipWallpaper := flag.Bool("skip-wallpaper", false, "skip wallpaper installation")
	skipCleanup := flag.Bool("skip-cleanup", false, "skip cleanup.sh")
	flag.Parse()

	home, _ := os.UserHomeDir()

	printBanner("Installer")
	fmt.Println("Welcome to h2install")

	if !promptYesNo("Do you want to start the installation now?") {
		fmt.Println("Installation cancelled.")
		return
	}

	// Pacman bootstrap
	pacmanConfig := filepath.Join(home, "Downloads", "hyprarch2", "src", "Scripts", "pacman.sh")
	if fileExists(pacmanConfig) && !*dryRun {
		printBanner("pacman.sh")
		runCommand(false, "bash", pacmanConfig)
	}

	// Choose AUR helper
	choice := promptChoice("Which AUR helper do you want to install?", []string{"paru", "yay", "cancel"})
	if choice == "cancel" {
		fmt.Println(";; Operation canceled.")
		return
	}
	if err := installPackages(choice, *dryRun); err != nil {
		fmt.Printf("Error installing packages: %v\n", err)
	}

	// Wallpapers
	if !*skipWallpaper {
		if err := installWallpaper("https://github.com/g5ostXa/wallpaper.git", filepath.Join(home, "wallpaper"), *force, *dryRun); err != nil {
			fmt.Printf("Error installing wallpapers: %v\n", err)
		}
	}

	// Bashrc
	if err := installBashrc(home, *dryRun); err != nil {
		fmt.Printf("Error handling .bashrc: %v\n", err)
	}

	// Symlinks
	createSymlinks(home, *force, *dryRun)

	// Cleanup
	cleanup(home, *skipCleanup, *dryRun)

	printBanner("h2install")
	fmt.Println(";; Installation status: COMPLETE")
}

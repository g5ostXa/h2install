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
	"time"
)

// ===== Global variables =====
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
	"pipewire",
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

// ===== Utility functions =====
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

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Chmod(0644)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ===== Banner =====
func printBanner(msg string) {
	if _, err := exec.LookPath("figlet"); err == nil {
		runCommand("figlet", "-f", "smslant", msg)
	} else {
		fmt.Printf("=== %s ===\n", msg)
	}
}

// ===== Prompt =====
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

// ===== .bashrc handling =====
func backupAndReplaceBashrc(homeDir, newBashrcPath string, force bool) error {
	bashrc := filepath.Join(homeDir, ".bashrc")

	if fileExists(bashrc) {
		backupName := fmt.Sprintf(".bashrc.bak.%s", time.Now().Format("20060102-150405"))
		backupPath := filepath.Join(homeDir, backupName)

		if err := copyFile(bashrc, backupPath); err != nil {
			return fmt.Errorf("failed to backup .bashrc: %w", err)
		}
		fmt.Printf(";; Backed up existing .bashrc to %s\n", backupPath)

		if !force {
			fmt.Println(";; Skipping replacement (use --force to overwrite).")
			return nil
		}
	}

	if err := copyFile(newBashrcPath, bashrc); err != nil {
		return fmt.Errorf("failed to copy new .bashrc: %w", err)
	}
	fmt.Println(";; Installed new .bashrc successfully.")
	return nil
}

// ===== Wallpapers =====
func installWallpaper(repoURL, dest string, dryRun bool) error {
	fmt.Println(";; Installing wallpapers...")
	if dryRun {
		fmt.Printf("Would clone %s into %s\n", repoURL, dest)
		return nil
	}
	if fileExists(dest) {
		if err := os.RemoveAll(dest); err != nil {
			return fmt.Errorf("failed to remove old wallpaper dir: %w", err)
		}
	}
	return runCommand("git", "clone", repoURL, dest)
}

// ===== AUR Helper + Packages =====
func installPackages(helper string, dryRun bool) error {
	home, _ := os.UserHomeDir()
	helperRepo := fmt.Sprintf("https://aur.archlinux.org/%s-bin.git", helper)
	helperDir := filepath.Join(home, helper+"-bin")

	if dryRun {
		fmt.Printf("Would install AUR helper %s and packages: %v\n", helper, packages)
		return nil
	}

	if err := runCommand("git", "clone", helperRepo); err != nil {
		return fmt.Errorf("failed to clone %s: %w", helper, err)
	}
	if err := os.Chdir(helperDir); err != nil {
		return fmt.Errorf("failed to cd into %s: %w", helperDir, err)
	}
	if err := runCommand("makepkg", "-si", "--noconfirm"); err != nil {
		return fmt.Errorf("failed to build %s: %w", helper, err)
	}
	if err := runCommand(helper, append([]string{"-S", "--needed", "--noconfirm"}, packages...)...); err != nil {
		return fmt.Errorf("failed to install packages: %w", err)
	}
	if err := os.RemoveAll(helperDir); err == nil {
		fmt.Printf(";; Removed %s\n", helperDir)
	}
	return nil
}

// ===== Symlinks =====
func createSymlinks(home string, force bool, dryRun bool) {
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

		if fileExists(dstPath) && force {
			os.Remove(dstPath)
		}
		if err := os.Symlink(srcPath, dstPath); err != nil {
			fmt.Printf(";; Failed to link %s → %s: %v\n", srcPath, dstPath, err)
		} else {
			fmt.Printf(";; Linked %s → %s\n", srcPath, dstPath)
		}
	}
}

// ===== Cleanup =====
func cleanup(home string, skip bool, dryRun bool) {
	hyprarchDir := filepath.Join(home, "Downloads", "hyprarch2")
	cleanupScript := filepath.Join(home, "src", "Scripts", "cleanup.sh")

	if dryRun {
		fmt.Printf("Would cleanup %s and maybe run cleanup.sh\n", hyprarchDir)
		return
	}
	if fileExists(hyprarchDir) {
		os.RemoveAll(hyprarchDir)
		fmt.Printf(";; Removed %s\n", hyprarchDir)
	}
	if fileExists(cleanupScript) && !skip {
		runCommand("bash", cleanupScript)
	}
	if _, err := exec.LookPath("trash-empty"); err == nil {
		runCommand("trash-empty")
	}
}

// ===== Main =====
func main() {
	dryRun := flag.Bool("dry-run", false, "simulate actions without executing")
	force := flag.Bool("force", false, "force overwrite of files/symlinks")
	skipWallpaper := flag.Bool("skip-wallpaper", false, "skip wallpaper installation")
	skipCleanup := flag.Bool("skip-cleanup", false, "skip cleanup.sh and trash-empty")
	flag.Parse()

	home, _ := os.UserHomeDir()

	// Banner
	printBanner("Installer")
	fmt.Println("Welcome to hyprarch2")

	if !promptYesNo("Do you want to start the installation now?") {
		fmt.Println("Installation cancelled.")
		return
	}

	// Pacman config
	pacmanConfig := filepath.Join(home, "Downloads", "hyprarch2", "src", "Scripts", "pacman.sh")
	if fileExists(pacmanConfig) && !*dryRun {
		printBanner("pacman.sh")
		runCommand("bash", pacmanConfig)
	}

	// AUR helper choice
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
		if err := installWallpaper("https://github.com/g5ostXa/wallpaper.git", filepath.Join(home, "wallpaper"), *dryRun); err != nil {
			fmt.Printf("Error installing wallpapers: %v\n", err)
		}
	}

	// .bashrc handling
	newBashrc := filepath.Join(home, "Downloads", "hyprarch2", ".bashrc")
	if err := backupAndReplaceBashrc(home, newBashrc, *force); err != nil {
		fmt.Printf("Error handling .bashrc: %v\n", err)
	}

	// Symlinks
	createSymlinks(home, *force, *dryRun)

	// Cleanup
	cleanup(home, *skipCleanup, *dryRun)

	printBanner("hyprarch2")
	fmt.Println(";; Installation status: COMPLETE")
}

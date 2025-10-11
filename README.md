## h2install
This is the new installer for [`hyprarch2`](https://github.com/g5ostXa/hyprarch2).

#### Why is this in a different repository?
- For now, it's only a helper for the main [`install.sh`](https://github.com/g5ostXa/hyprarch2/blob/master/src/install.sh), but in the future this will become the main installer.
- Having this in a separate repo helps for maintenance while allowing us to visualize weather we want to keep working on this, or keep the old installer.

#### How can I use this?
- Simply visit https://github.com/g5ostXa/hyprarch2 and follow the instructions.

> [!CAUTION]
> - It is recommended installing via the [`installer`](https://github.com/g5ostXa/h2install), which is managed by [`install.sh`](https://github.com/g5ostXa/hyprarch2/blob/master/src/install.sh).
> - The installer installs [`dotfiles/`](https://github.com/g5ostXa/hyprarch2/tree/master/dotfiles) in your home directroy and create symlinks that point to `~/.config/`
> - We are currently working on renaming `~/dotfiles/` to `~/.config/`, but for now we still use symlinks.
> - This is NOT compatible with a different distro than upstream [`Archlinux`](https://archlinux.org).

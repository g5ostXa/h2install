## h2install


> [!CAUTION]
> - This is the new installer for [`hyprarch2`](https://github.com/g5ostXa/hyprarch2) which is written in golang.
> - It is still very unstable, even though it works fine for me at the momment.
> - I recommend running in `--dry-run` mode first.

To do that, edit [`install.sh`](https://github.com/g5ostXa/hyprarch2/src/install.sh) and add `--dry-run` when running the installer like shown below:
```bash
func_main() {
	src_check && src_copy && target_check

	if [ -f "/etc/issue" ]; then
		sudo chown root:root /etc/issue
	else
		echo -e "${YELLOW}Failed to copy issue to /etc...${RC}"
		sleep 1
	fi

	cd "$HOME/Downloads" && git clone --depth=1 https://github.com/g5ostXa/h2install.git
	cd h2install && rm -rf .git/ && go mod tidy && go build -o h2installer && ./h2installer --dry-run

}
```

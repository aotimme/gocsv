#!/bin/sh
# Compatibility entry point for existing installation commands.
set -eu

main() {
    local_script=$(dirname "$0")/install.sh
    if [ -f "$local_script" ]; then
        exec sh "$local_script" "$@"
    fi

    tmp_dir=$(mktemp -d)
    trap 'rm -rf "$tmp_dir"' EXIT
    trap 'exit 1' HUP INT TERM
    curl -fLsS https://raw.githubusercontent.com/aotimme/gocsv/main/scripts/install.sh > "$tmp_dir/install.sh"
    sh "$tmp_dir/install.sh" "$@"
}

main "$@"

#!/bin/sh
# Install a GoCSV release without requiring root. Safe to pipe into sh.
set -eu

fail() {
    printf 'Error: %s\n' "$*" >&2
    exit 1
}

main() {
    case $(uname -s) in
        Darwin) os=darwin ;;
        Linux) os=linux ;;
        *) fail 'Supported operating systems are Linux and macOS.' ;;
    esac
    case $(uname -m) in
        x86_64|amd64) arch=amd64 ;;
        arm64|aarch64) arch=arm64 ;;
        i386|i486|i586|i686) arch=386 ;;
        armv7l|armv8l) arch=arm ;;
        ppc64le|riscv64) arch=$(uname -m) ;;
        *) fail "Unsupported architecture: $(uname -m)" ;;
    esac
    case "$os:$arch" in
        darwin:amd64|darwin:arm64|linux:*) ;;
        *) fail "No release binary is available for $os/$arch." ;;
    esac

    command -v unzip >/dev/null 2>&1 || fail 'unzip is required to install GoCSV.'
    if command -v curl >/dev/null 2>&1; then
        downloader=curl
    elif command -v wget >/dev/null 2>&1; then
        downloader=wget
    else
        fail 'curl or wget is required to install GoCSV.'
    fi

    install_dir=${GOCSV_INSTALL_DIR:-"$HOME/.local/bin"}
    archive="gocsv-$os-$arch"
    version=${GOCSV_VERSION:-latest}
    case "$version" in
        latest) release_path=latest/download ;;
        .|..|*[!a-zA-Z0-9._+-]*) fail 'Set GOCSV_VERSION to a release tag such as v1.0.7.' ;;
        *) release_path="download/$version" ;;
    esac
    url="https://github.com/aotimme/gocsv/releases/$release_path/$archive.zip"
    tmp_dir=$(mktemp -d)
    staged_binary=
    trap 'rm -rf "$tmp_dir"; if [ -n "$staged_binary" ]; then rm -f "$staged_binary"; fi' EXIT
    trap 'exit 1' HUP INT TERM

    printf 'Downloading %s (%s)...\n' "$archive" "$version"
    if [ "$downloader" = curl ]; then
        curl -fLsS "$url" -o "$tmp_dir/release.zip" || fail "Download failed: $url"
    else
        wget -q "$url" -O "$tmp_dir/release.zip" || fail "Download failed: $url"
    fi
    unzip -q "$tmp_dir/release.zip" -d "$tmp_dir" || fail 'Could not extract the release archive.'
    binary="$tmp_dir/$archive/gocsv"
    [ -f "$binary" ] || fail 'The release archive does not contain the expected gocsv binary.'
    mkdir -p "$install_dir" || fail "Cannot create $install_dir. Set GOCSV_INSTALL_DIR to a writable directory."
    [ ! -d "$install_dir/gocsv" ] || fail "The destination $install_dir/gocsv is a directory."
    # Stage on the same filesystem so the final rename replaces the old binary atomically.
    staged_binary=$(mktemp "$install_dir/.gocsv.XXXXXX") || fail "Cannot create a temporary file in $install_dir."
    install -m 755 "$binary" "$staged_binary" || fail "Cannot copy the binary to $install_dir."
    mv -f "$staged_binary" "$install_dir/gocsv" || fail "Cannot install to $install_dir. Set GOCSV_INSTALL_DIR to a writable directory."
    staged_binary=

    printf 'GoCSV installed to %s/gocsv\n' "$install_dir"
    case ":${PATH:-}:" in
        *:"$install_dir":*) printf 'Run: gocsv help\n' ;;
        *) printf 'Add %s to your PATH, then run: gocsv help\n' "$install_dir" ;;
    esac
}

main "$@"

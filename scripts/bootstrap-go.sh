#!/usr/bin/env bash
set -euo pipefail

root_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)
target_dir="$root_dir/.tools/go"
version=1.24.1

if [[ -x "$target_dir/bin/go" ]]; then
    "$target_dir/bin/go" version
    exit 0
fi
if [[ -e "$target_dir" ]]; then
    printf 'Installation directory already exists: %s\n' "$target_dir" >&2
    exit 1
fi

case "$(uname -s)/$(uname -m)" in
    Linux/x86_64)
        platform=linux-amd64
        checksum=cb2396bae64183cdccf81a9a6df0aea3bce9511fc21469fb89a0c00470088073 ;;
    Linux/aarch64|Linux/arm64)
        platform=linux-arm64
        checksum=8df5750ffc0281017fb6070fba450f5d22b600a02081dceef47966ffaf36a3af ;;
    Darwin/x86_64)
        platform=darwin-amd64
        checksum=addbfce2056744962e2d7436313ab93486660cf7a2e066d171b9d6f2da7c7abe ;;
    Darwin/arm64)
        platform=darwin-arm64
        checksum=295581b5619acc92f5106e5bcb05c51869337eb19742fdfa6c8346c18e78ff88 ;;
    *)
        printf 'Unsupported platform. Install Go %s or newer from https://go.dev/dl/.\n' "$version" >&2
        exit 1 ;;
esac

# Checksums are from https://go.dev/dl/?mode=json&include=all.
mkdir -p "$root_dir/.tools"
temp_dir=$(mktemp -d "$root_dir/.tools/bootstrap.XXXXXX")
trap 'rm -rf -- "$temp_dir"' EXIT
archive="go${version}.${platform}.tar.gz"
curl --fail --location --retry 3 --connect-timeout 15 --max-time 300 \
    "https://go.dev/dl/$archive" --output "$temp_dir/$archive"
(
    cd -- "$temp_dir"
    if command -v sha256sum >/dev/null 2>&1; then
        printf '%s  %s\n' "$checksum" "$archive" | sha256sum --check -
    else
        printf '%s  %s\n' "$checksum" "$archive" | shasum -a 256 --check -
    fi
)
tar -xzf "$temp_dir/$archive" -C "$temp_dir"
mv -- "$temp_dir/go" "$target_dir"
"$target_dir/bin/go" version

#!/usr/bin/env sh
# Portable installer for volume-speeder (Linux: CLI, macOS: GUI)
# Defaults assume GitHub repo owner and name, but can be overridden:
#   OWNER=<owner> REPO=<repo> TAG=<tag> sh install.sh
# Or flags: --owner <o> --repo <r> --tag <t> --name <install_name>
# Linux installs binary to /usr/local/bin/<install_name> (default: volume-speeder)
# macOS installs the .app to /Applications
# Note: Without an Apple Developer ID, the macOS app is unsigned and Gatekeeper may block it.
# Pass --trust (or TRUST=1) to remove the quarantine attribute automatically after install.

set -eu

OWNER_DEFAULT="axlroden"
REPO_DEFAULT="volume-speeder"
INSTALL_NAME_DEFAULT="volume-speeder"

OWNER="${OWNER:-$OWNER_DEFAULT}"
REPO="${REPO:-$REPO_DEFAULT}"
TAG="${TAG:-}"
INSTALL_NAME="${INSTALL_NAME:-$INSTALL_NAME_DEFAULT}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
TRUST="${TRUST:-}"

usage() {
  echo "Usage: OWNER=<owner> REPO=<repo> [TAG=<tag>] sh install.sh"
  echo "Flags: --owner <owner> --repo <repo> [--tag <tag>] [--name <install_name>] [--dir <install_dir>] [--trust]"
  echo "Env: TRUST=1   # remove macOS quarantine after install"
}

# Basic flag parsing (POSIX)
while [ $# -gt 0 ]; do
  case "$1" in
    --owner) shift; OWNER=${1:-$OWNER};;
    --repo) shift; REPO=${1:-$REPO};;
    --tag) shift; TAG=${1:-$TAG};;
    --name) shift; INSTALL_NAME=${1:-$INSTALL_NAME};;
    --dir) shift; INSTALL_DIR=${1:-$INSTALL_DIR};;
  --trust) TRUST=1;;
    -h|--help) usage; exit 0;;
    *) echo "Unknown arg: $1" >&2; usage; exit 2;;
  esac
  shift || true
done

need_cmd() { command -v "$1" >/dev/null 2>&1 || { echo "Missing required command: $1" >&2; exit 1; }; }

detect_os() {
  u=$(uname -s 2>/dev/null || echo unknown)
  case "$u" in
    Linux) echo linux;;
    Darwin) echo darwin;;
    *) echo "Unsupported OS: $u" >&2; exit 1;;
  esac
}

detect_arch() {
  a=$(uname -m 2>/dev/null || echo unknown)
  case "$a" in
    x86_64|amd64) echo amd64;;
    aarch64|arm64) echo arm64;;
    *) echo "Unsupported arch: $a" >&2; exit 1;;
  esac
}

http_get() {
  # $1=url $2=out
  url="$1" out="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$out"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "$out" "$url"
  else
    echo "Need curl or wget" >&2; exit 1
  fi
}

http_get_stdout() {
  url="$1"
  if command -v curl >/dev/null 2>&1; then
    if [ -n "${GH_TOKEN:-}" ]; then
      curl -fsSL -H "Authorization: Bearer $GH_TOKEN" "$url"
    else
      curl -fsSL "$url"
    fi
  else
    if [ -n "${GH_TOKEN:-}" ]; then
      wget -qO- --header "Authorization: Bearer $GH_TOKEN" "$url"
    else
      wget -qO- "$url"
    fi
  fi
}

sudo_cp() {
  src="$1" dst="$2"
  if [ -w "$(dirname "$dst")" ]; then
    cp "$src" "$dst"
  else
    need_cmd sh
    if command -v sudo >/dev/null 2>&1; then
      sudo cp "$src" "$dst"
    else
      echo "Need sudo to write to $(dirname "$dst")" >&2; exit 1
    fi
  fi
}

make_executable() {
  f="$1"
  chmod +x "$f" 2>/dev/null || { command -v sudo >/dev/null 2>&1 && sudo chmod +x "$f"; }
}

LATEST_TAG() {
  api="https://api.github.com/repos/$OWNER/$REPO/releases/latest"
  http_get_stdout "$api" | awk -F '"' '/tag_name/ {print $4; exit}'
}

find_asset_url() {
  tag="$1" pattern="$2"
  api="https://api.github.com/repos/$OWNER/$REPO/releases/tags/$tag"
  http_get_stdout "$api" | awk -v p="$pattern" -F '"' '/browser_download_url/ { if ($4 ~ p) print $4 }' | head -n1
}

install_linux() {
  need_cmd tar
  os=linux
  arch=$(detect_arch)
  tag="$TAG"
  if [ -z "$tag" ]; then tag=$(LATEST_TAG); fi
  [ -n "$tag" ] || { echo "Could not determine latest tag. Set TAG explicitly." >&2; exit 1; }
  pat="-${os}-${arch}\\.tar\\.gz$"
  url=$(find_asset_url "$tag" "$pat")
  tmpdir=$(mktemp -d)
  trap 'rm -rf "$tmpdir"' EXIT INT HUP TERM

  if [ -n "$url" ]; then
    echo "Downloading $url"
    tarball="$tmpdir/app.tar.gz"
    http_get "$url" "$tarball"
    (cd "$tmpdir" && tar xzf "$tarball")
    # Pick first executable file
    bin=$(find "$tmpdir" -maxdepth 2 -type f -perm -111 ! -name "*.sh" | head -n1 || true)
    [ -n "$bin" ] || { echo "No binary found in tarball" >&2; exit 1; }
    dst="$INSTALL_DIR/$INSTALL_NAME"
    echo "Installing to $dst"
    sudo_cp "$bin" "$dst"
    make_executable "$dst"
  else
    echo "No release asset found for $os/$arch @ $tag. Falling back to build from source."
    need_cmd go
    gobin=$(go env GOBIN || true)
    if [ -z "$gobin" ]; then gobin="$(go env GOPATH)/bin"; fi
  pkg="github.com/$OWNER/$REPO/app@$tag"
    echo "go install $pkg"
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 go install "$pkg"
  src="$gobin/$REPO"
    if [ ! -f "$src" ]; then
      # some modules may produce binary named differently; search
      src=$(ls -1 "$gobin" | grep -E "^$REPO$|volume-boost|volume.*spee" | head -n1 || true)
      src="$gobin/$src"
    fi
    [ -f "$src" ] || { echo "Built binary not found at $src" >&2; exit 1; }
    dst="$INSTALL_DIR/$INSTALL_NAME"
    echo "Installing to $dst"
    sudo_cp "$src" "$dst"
    make_executable "$dst"
  fi

  # Ensure config exists
  conf="$HOME/.config/volume-speeder/config"
  if [ ! -f "$conf" ]; then
    echo "Creating default config at $conf"
    mkdir -p "$(dirname "$conf")"
    echo "multiplier=3" > "$conf"
  fi

  echo "\nInstalled $INSTALL_NAME. Try:\n  $INSTALL_NAME get\n  $INSTALL_NAME set 5\n"
}

install_macos() {
  os=darwin
  arch=universal
  tag="$TAG"
  if [ -z "$tag" ]; then tag=$(LATEST_TAG); fi
  [ -n "$tag" ] || { echo "Could not determine latest tag. Set TAG explicitly." >&2; exit 1; }
  # Prefer dmg, fallback to zip
  url=$(find_asset_url "$tag" "-macos-${arch}\\.dmg$")
  kind=dmg
  if [ -z "$url" ]; then
    url=$(find_asset_url "$tag" "-macos-${arch}\\.zip$")
    kind=zip
  fi
  [ -n "$url" ] || { echo "No macOS artifact found for tag $tag" >&2; exit 1; }

  tmpdir=$(mktemp -d)
  trap 'rm -rf "$tmpdir"' EXIT INT HUP TERM

  echo "Downloading $url"
  asset="$tmpdir/app.$kind"
  http_get "$url" "$asset"

  if [ "$kind" = dmg ]; then
    need_cmd hdiutil
    vol=$(hdiutil attach "$asset" -nobrowse -quiet | awk 'END{print $3}')
    [ -n "$vol" ] || { echo "Failed to mount DMG" >&2; exit 1; }
    app=$(find "$vol" -maxdepth 1 -type d -name "*.app" | head -n1 || true)
    [ -n "$app" ] || { echo "No .app found in DMG" >&2; hdiutil detach "$vol" -quiet || true; exit 1; }
    dest="/Applications/$(basename "$app")"
    echo "Installing to $dest"
    if [ -w "/Applications" ]; then
      cp -R "$app" "$dest"
    else
      need_cmd sudo
      sudo cp -R "$app" "$dest"
    fi
    hdiutil detach "$vol" -quiet || true
    echo "Installed $(basename "$app") to /Applications"
  else
    need_cmd unzip
    unzip -q "$asset" -d "$tmpdir"
    app=$(find "$tmpdir" -type d -name "*.app" | head -n1 || true)
    [ -n "$app" ] || { echo "No .app found in zip" >&2; exit 1; }
    dest="/Applications/$(basename "$app")"
    echo "Installing to $dest"
    if [ -w "/Applications" ]; then
      cp -R "$app" "$dest"
    else
      need_cmd sudo
      sudo cp -R "$app" "$dest"
    fi
    echo "Installed $(basename "$app") to /Applications"
  fi

  # Optionally remove quarantine so the app opens without Gatekeeper prompts
  if [ -n "$TRUST" ]; then
    if command -v xattr >/dev/null 2>&1; then
      echo "Removing quarantine attribute from $dest"
      xattr -dr com.apple.quarantine "$dest" 2>/dev/null || {
        if command -v sudo >/dev/null 2>&1; then sudo xattr -dr com.apple.quarantine "$dest"; fi
      }
    fi
  else
    echo "Note: On first launch, macOS may warn about an unsigned app."
    echo "- Right-click the app in /Applications and choose Open, then Open again"
    echo "- Or remove quarantine manually: xattr -dr com.apple.quarantine \"$dest\""
    echo "- Or re-run this installer with --trust (or TRUST=1) to do this automatically"
  fi
}

main() {
  os=$(detect_os)
  case "$os" in
    linux) install_linux ;;
    darwin) install_macos ;;
  esac
}

main "$@"

#!/bin/env sh
# Cairn CLI installer
# Usage: curl -sSL https://raw.githubusercontent.com/ndy40/cairn/main/install.sh | sh
#
# MIT License - https://github.com/ndy40/cairn

set -eu

# ── Defaults ──────────────────────────────────────────────────────────────────

INSTALL_DIR="${CAIRN_INSTALL_DIR:-"/usr/local/bin"}"
VERSION=""
NON_INTERACTIVE=0
WITH_EXTENSION=0
WITH_DESKTOP=0
GITHUB_REPO="ndy40/cairn"
TMPDIR_CAIRN=""
DESKTOP_ARCHIVE=""

# ── Utility functions ─────────────────────────────────────────────────────────

log_info() {
    printf '  %s\n' "$*"
}

log_error() {
    printf '  ERROR: %s\n' "$*" >&2
}

log_success() {
    printf '  %s\n' "$*"
}

has_command() {
    command -v "$1" >/dev/null 2>&1
}

# ── Cleanup ───────────────────────────────────────────────────────────────────

cleanup() {
    if [ -n "$TMPDIR_CAIRN" ] && [ -d "$TMPDIR_CAIRN" ]; then
        rm -rf "$TMPDIR_CAIRN"
    fi
}

trap cleanup EXIT INT TERM

# ── Help ──────────────────────────────────────────────────────────────────────

show_help() {
    cat <<'HELP'
Cairn CLI Installer

USAGE:
    install.sh [OPTIONS]

OPTIONS:
    -d, --install-dir DIR    Installation directory (default: /usr/local/bin)
    -v, --version VERSION    Install a specific version (e.g., v0.0.1)
    -y, --non-interactive    Skip all prompts, install CLI only
        --with-extension     Also install Vicinae extension (non-interactive mode)
        --with-desktop       Also install the Cairn desktop app (non-interactive mode)
        --no-autostart       Deprecated, does nothing: start at login is now set in
                             Cairn's Settings
    -h, --help               Show this help message

ENVIRONMENT VARIABLES:
    CAIRN_INSTALL_DIR        Override install directory (flag takes precedence)

EXIT CODES:
    0    Installation completed successfully
    1    General error (network failure, unexpected error)
    2    Unsupported platform (OS or architecture)
    3    Checksum verification failed
    4    Permission denied (cannot write to install directory)

EXAMPLES:
    # Install latest version
    curl -sSL https://raw.githubusercontent.com/ndy40/cairn/main/install.sh | sh

    # Install to a custom directory
    sh install.sh --install-dir ~/bin

    # Install a specific version
    sh install.sh --version v0.0.1

    # Non-interactive install (CI/CD)
    sh install.sh -y

    # Non-interactive with Vicinae extension
    sh install.sh -y --with-extension

    # Non-interactive with the desktop app
    sh install.sh -y --with-desktop
HELP
}

# ── Argument parsing ──────────────────────────────────────────────────────────

parse_args() {
    while [ $# -gt 0 ]; do
        case "$1" in
            -d|--install-dir)
                [ $# -ge 2 ] || { log_error "Missing argument for $1"; exit 1; }
                INSTALL_DIR="$2"
                shift 2
                ;;
            -v|--version)
                [ $# -ge 2 ] || { log_error "Missing argument for $1"; exit 1; }
                VERSION="$2"
                shift 2
                ;;
            -y|--non-interactive)
                NON_INTERACTIVE=1
                shift
                ;;
            --with-extension)
                WITH_EXTENSION=1
                shift
                ;;
            --with-desktop)
                WITH_DESKTOP=1
                shift
                ;;
            --no-autostart)
                # Kept so existing scripts don't fail; the app now owns the login item.
                log_info "Start at login is now set in Cairn's Settings"
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                log_info "Run 'install.sh --help' for usage information."
                exit 1
                ;;
        esac
    done
}

# ── Platform detection ────────────────────────────────────────────────────────

detect_os() {
    case "$(uname -s)" in
        Linux)  echo "linux" ;;
        Darwin) echo "darwin" ;;
        *)
            log_error "Unsupported operating system: $(uname -s)"
            log_info "Cairn supports Linux and macOS only."
            log_info "For other platforms, download binaries from:"
            log_info "  https://github.com/$GITHUB_REPO/releases"
            exit 2
            ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   echo "amd64" ;;
        aarch64|arm64)  echo "arm64" ;;
        *)
            log_error "Unsupported architecture: $(uname -m)"
            log_info "Cairn supports amd64 and arm64 only."
            log_info "For other architectures, build from source:"
            log_info "  https://github.com/$GITHUB_REPO"
            exit 2
            ;;
    esac
}

# ── Download helpers ──────────────────────────────────────────────────────────

download() {
    url="$1"
    dest="$2"

    if has_command curl; then
        http_code=$(curl -sL -w '%{http_code}' -o "$dest" "$url")
        if [ "$http_code" -ge 400 ] 2>/dev/null || [ "$http_code" = "000" ]; then
            rm -f "$dest"
            log_error "Download failed: HTTP $http_code for $url"
            return 1
        fi
    elif has_command wget; then
        if ! wget -qO "$dest" "$url" 2>/dev/null; then
            rm -f "$dest"
            log_error "Download failed: $url"
            return 1
        fi
    else
        log_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
}

resolve_version() {
    if [ -n "$VERSION" ]; then
        echo "$VERSION"
        return
    fi

    # Fetch latest release tag from GitHub API
    latest_url="https://api.github.com/repos/$GITHUB_REPO/releases/latest"
    tmpfile=$(mktemp)

    if has_command curl; then
        http_code=$(curl -sL -w '%{http_code}' -o "$tmpfile" "$latest_url")
        if [ "$http_code" = "403" ]; then
            rm -f "$tmpfile"
            log_error "GitHub API rate limit exceeded."
            log_info "Specify a version explicitly: install.sh --version v0.0.1"
            log_info "Or wait and try again later."
            exit 1
        fi
        if [ "$http_code" = "404" ]; then
            rm -f "$tmpfile"
            log_error "No releases found for $GITHUB_REPO."
            log_info "Check https://github.com/$GITHUB_REPO/releases"
            exit 1
        fi
        if [ "$http_code" -ge 400 ] 2>/dev/null || [ "$http_code" = "000" ]; then
            rm -f "$tmpfile"
            log_error "Failed to query GitHub API (HTTP $http_code)."
            log_info "Check your internet connection or specify a version with --version."
            exit 1
        fi
    elif has_command wget; then
        if ! wget -qO "$tmpfile" "$latest_url" 2>/dev/null; then
            rm -f "$tmpfile"
            log_error "Failed to query GitHub API."
            log_info "Check your internet connection or specify a version with --version."
            exit 1
        fi
    else
        rm -f "$tmpfile"
        log_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi

    tag=$(grep '"tag_name"' "$tmpfile" | head -1 | sed 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/')
    rm -f "$tmpfile"

    if [ -z "$tag" ]; then
        log_error "Could not determine latest release version."
        log_info "The API response may be unexpected. Specify a version explicitly:"
        log_info "  install.sh --version v0.0.1"
        exit 1
    fi

    echo "$tag"
}

# ── Checksum verification ────────────────────────────────────────────────────

verify_checksum() {
    file="$1"
    expected="$2"

    if has_command sha256sum; then
        actual=$(sha256sum "$file" | awk '{print $1}')
    elif has_command shasum; then
        actual=$(shasum -a 256 "$file" | awk '{print $1}')
    else
        log_error "Neither sha256sum nor shasum found. Cannot verify download integrity."
        exit 1
    fi

    if [ "$actual" != "$expected" ]; then
        log_error "Checksum verification failed!"
        log_info "Expected: $expected"
        log_info "Actual:   $actual"
        log_info "The downloaded file may be corrupted. Please try again."
        exit 3
    fi
}

# ── Binary installation ──────────────────────────────────────────────────────

install_binary() {
    src="$1"
    dest_dir="$2"
    dest="$dest_dir/${3:-cairn}"

    # Create install directory if needed
    if ! mkdir -p "$dest_dir" 2>/dev/null; then
        log_error "Permission denied: cannot create directory '$dest_dir'."
        log_info "Try one of the following:"
        log_info "  - Use a different directory: install.sh --install-dir ~/bin"
        log_info "  - Fix permissions: sudo mkdir -p '$dest_dir' && sudo chown \$(whoami) '$dest_dir'"
        exit 4
    fi

    if [ ! -w "$dest_dir" ]; then
        log_error "Permission denied: cannot write to '$dest_dir'."
        log_info "Try one of the following:"
        log_info "  - Use a different directory: install.sh --install-dir ~/bin"
        log_info "  - Fix permissions: sudo chown \$(whoami) '$dest_dir'"
        exit 4
    fi

    # Copy to a temporary file in the same directory, then rename it over the
    # old binary. Unlike copying in place, a rename works while the old binary
    # is running (copying fails with "Text file busy"), and the old binary is
    # left untouched if anything goes wrong.
    tmp="$dest_dir/.$(basename "$dest").new.$$"
    if ! err=$(cp "$src" "$tmp" 2>&1 && chmod +x "$tmp" 2>&1 && mv -f "$tmp" "$dest" 2>&1); then
        rm -f "$tmp"
        log_error "Could not install '$dest': $err"
        exit 1
    fi
}

check_path() {
    dir="$1"

    case ":$PATH:" in
        *":$dir:"*) return 0 ;;
    esac

    log_info ""
    log_info "NOTE: '$dir' is not in your PATH."
    log_info "Add it by running one of the following:"
    log_info ""

    if [ -f "$HOME/.zshrc" ]; then
        log_info "  echo 'export PATH=\"$dir:\$PATH\"' >> ~/.zshrc && source ~/.zshrc"
    fi
    if [ -f "$HOME/.bashrc" ]; then
        log_info "  echo 'export PATH=\"$dir:\$PATH\"' >> ~/.bashrc && source ~/.bashrc"
    fi
    if [ -f "$HOME/.profile" ]; then
        log_info "  echo 'export PATH=\"$dir:\$PATH\"' >> ~/.profile"
    fi
    if [ ! -f "$HOME/.zshrc" ] && [ ! -f "$HOME/.bashrc" ] && [ ! -f "$HOME/.profile" ]; then
        log_info "  export PATH=\"$dir:\$PATH\""
        log_info ""
        log_info "  Add this to your shell profile to make it permanent."
    fi
}

# ── Vicinae extension ────────────────────────────────────────────────────────

detect_vicinae() {
    has_command vici
}

prompt_extension() {
    if [ "$NON_INTERACTIVE" = 1 ]; then
        if [ "$WITH_EXTENSION" = 1 ]; then
            return 0  # install extension
        fi
        return 1  # skip extension
    fi

    # If stdin is not a terminal, treat as non-interactive
    if [ ! -t 0 ]; then
        return 1
    fi

    printf '  Vicinae detected. Install Cairn extension? [y/N] '
    read -r answer
    case "$answer" in
        [yY]|[yY][eE][sS]) return 0 ;;
        *) return 1 ;;
    esac
}

install_extension() {
    version="$1"
    ext_archive="cairn-vicinae-extension.tar.gz"
    ext_url="https://github.com/$GITHUB_REPO/releases/download/$version/$ext_archive"
    checksums_file="$TMPDIR_CAIRN/checksums.txt"

    log_info "Downloading Vicinae extension..."

    ext_dest="$TMPDIR_CAIRN/$ext_archive"
    if ! download "$ext_url" "$ext_dest"; then
        log_error "Failed to download Vicinae extension."
        log_info "The extension may not be available for this release."
        log_info "You can install it manually from: https://github.com/$GITHUB_REPO"
        return 1
    fi

    # Verify checksum if checksums file exists
    if [ -f "$checksums_file" ]; then
        expected=$(grep "$ext_archive" "$checksums_file" | awk '{print $1}')
        if [ -n "$expected" ]; then
            verify_checksum "$ext_dest" "$expected"
        fi
    fi

    # Extract and install to Vicinae extensions directory
    ext_tmp="$TMPDIR_CAIRN/extension"
    mkdir -p "$ext_tmp"
    tar -xzf "$ext_dest" -C "$ext_tmp"

    # Determine Vicinae extensions directory
    os="$2"
    case "$os" in
        darwin)
            ext_dir="$HOME/Library/Application Support/vicinae/extensions/cairn"
            ;;
        linux)
            ext_dir="${XDG_DATA_HOME:-$HOME/.local/share}/vicinae/extensions/cairn"
            ;;
    esac

    mkdir -p "$ext_dir"
    cp -r "$ext_tmp"/* "$ext_dir"/

    log_success "Vicinae extension installed to $ext_dir"
}

# ── Desktop app ──────────────────────────────────────────────────────────────

# desktop_archive_name prints the release archive name for the desktop app on
# the given OS and architecture.
desktop_archive_name() {
    case "$1" in
        linux)  echo "cairn-desktop-linux-$2.tar.gz" ;;
        darwin) echo "cairn-desktop-darwin-universal.zip" ;;
    esac
}

# desktop_available reports whether the release ships a desktop archive for
# this platform. Every release asset is listed in checksums.txt, so this needs
# no extra request; releases made before the desktop app existed have no entry.
desktop_available() {
    archive=$(desktop_archive_name "$1" "$2")
    [ -n "$archive" ] && grep -q " $archive\$" "$TMPDIR_CAIRN/checksums.txt"
}

prompt_desktop() {
    if [ "$WITH_DESKTOP" = 1 ]; then
        return 0
    fi
    if [ "$NON_INTERACTIVE" = 1 ] || [ ! -t 0 ]; then
        return 1
    fi

    printf '  Install the Cairn desktop app? [y/N] '
    read -r answer
    case "$answer" in
        [yY]|[yY][eE][sS]) return 0 ;;
        *) return 1 ;;
    esac
}

# download_desktop_archive fetches a desktop release archive and verifies it
# against checksums.txt, setting DESKTOP_ARCHIVE to the downloaded path. It is
# not run in a subshell so its messages are shown and a checksum mismatch
# aborts the installer, as it does for the CLI binary.
download_desktop_archive() {
    version="$1"
    archive="$2"
    dest="$TMPDIR_CAIRN/$archive"

    if ! download "https://github.com/$GITHUB_REPO/releases/download/$version/$archive" "$dest"; then
        log_error "Failed to download the desktop app."
        log_info "The desktop app may not be available for $version."
        return 1
    fi

    expected=$(grep " $archive\$" "$TMPDIR_CAIRN/checksums.txt" | awk '{print $1}')
    if [ -z "$expected" ]; then
        log_error "Could not find checksum for $archive in checksums.txt"
        return 1
    fi
    verify_checksum "$dest" "$expected"
    DESKTOP_ARCHIVE="$dest"
}

# desktop_running reports whether this user is running cairn-desktop.
desktop_running() {
    has_command pgrep && pgrep -x -u "$(id -u)" cairn-desktop >/dev/null 2>&1
}

# stop_desktop quits the running app and waits up to 5 seconds for it to exit.
stop_desktop() {
    pkill -x -u "$(id -u)" cairn-desktop 2>/dev/null || true
    i=0
    while desktop_running && [ "$i" -lt 5 ]; do
        sleep 1
        i=$((i + 1))
    done
}

install_desktop_linux() {
    version="$1"
    arch="$2"
    name="cairn-desktop-linux-$arch"

    log_info "Downloading desktop app..."
    download_desktop_archive "$version" "$name.tar.gz" || return 1
    tar -xzf "$DESKTOP_ARCHIVE" -C "$TMPDIR_CAIRN"

    was_running=0
    if desktop_running; then
        was_running=1
    fi

    install_binary "$TMPDIR_CAIRN/$name/cairn-desktop" "$INSTALL_DIR" cairn-desktop

    data_dir="${XDG_DATA_HOME:-$HOME/.local/share}"
    icon="$data_dir/cairn/cairn-desktop.png"
    mkdir -p "$data_dir/cairn" "$data_dir/applications"
    cp "$TMPDIR_CAIRN/$name/cairn-desktop.png" "$icon"
    # `cairn update` compares this with the latest release.
    echo "$version" > "$data_dir/cairn/desktop-version.txt"

    write_desktop_entry "$data_dir/applications/cairn-desktop.desktop" "$INSTALL_DIR/cairn-desktop" "$icon"

    log_success "Desktop app installed to $INSTALL_DIR/cairn-desktop"

    # The old version keeps running from the replaced file until it exits, so
    # restart it to pick up the update.
    if [ "$was_running" = 1 ]; then
        log_info "Restarting the desktop app..."
        stop_desktop
        if [ -n "${WAYLAND_DISPLAY:-}${DISPLAY:-}" ]; then
            nohup "$INSTALL_DIR/cairn-desktop" --background >/dev/null 2>&1 &
        else
            log_info "No graphical session found; start it with: cairn-desktop --background"
        fi
    fi

    # GTK4 and WebKitGTK 6.0 are runtime dependencies that are not bundled.
    if has_command ldconfig && ! ldconfig -p 2>/dev/null | grep -q 'libwebkitgtk-6.0'; then
        log_info ""
        log_info "NOTE: WebKitGTK 6.0 was not found. Install it before running cairn-desktop:"
        log_info "  Fedora:        sudo dnf install webkitgtk6.0"
        log_info "  Debian/Ubuntu: sudo apt install libwebkitgtk-6.0-4"
        log_info "  Arch:          sudo pacman -S webkitgtk-6.0"
    fi

    log_info ""
    log_info "Open Cairn from the app menu or with Ctrl+Shift+Space, or bind a"
    log_info "system keyboard shortcut to: $INSTALL_DIR/cairn-desktop --toggle"
    log_info "To start Cairn at login, turn it on in Cairn's Settings."
}

write_desktop_entry() {
    cat > "$1" <<EOF_DESKTOP
[Desktop Entry]
Type=Application
Name=Cairn
Comment=Search, add and open your Cairn bookmarks
Exec=$2
Icon=$3
Terminal=false
Categories=Utility;
StartupWMClass=cairn-desktop
EOF_DESKTOP
}

install_desktop_darwin() {
    version="$1"
    app_dir="$HOME/Applications"
    app="$app_dir/cairn-desktop.app"

    log_info "Downloading desktop app..."
    download_desktop_archive "$version" "cairn-desktop-darwin-universal.zip" || return 1

    # Unpack first so a failed extraction never leaves the user without an app.
    unpacked="$TMPDIR_CAIRN/desktop-app"
    mkdir -p "$unpacked"
    if ! ditto -x -k "$DESKTOP_ARCHIVE" "$unpacked"; then
        log_error "Could not unpack the desktop app."
        return 1
    fi

    was_running=0
    if desktop_running; then
        was_running=1
        stop_desktop
    fi

    mkdir -p "$app_dir"
    rm -rf "$app"
    mv "$unpacked/cairn-desktop.app" "$app"
    # `cairn update` compares this with the latest release.
    mkdir -p "$HOME/Library/Application Support/cairn"
    echo "$version" > "$HOME/Library/Application Support/cairn/desktop-version.txt"

    log_success "Desktop app installed to $app"

    if [ "$was_running" = 1 ]; then
        log_info "Restarting the desktop app..."
        open -g "$app" --args --background
        return
    fi

    log_info ""
    log_info "Start it now with: open '$app'"
    log_info "Then bring it forward from anywhere with Cmd+Shift+Space."
    log_info "To start Cairn at login, turn it on in Cairn's Settings."
}

install_desktop() {
    case "$2" in
        linux)  install_desktop_linux "$1" "$3" ;;
        darwin) install_desktop_darwin "$1" ;;
    esac
}

# ── Main ──────────────────────────────────────────────────────────────────────

main() {
    parse_args "$@"

    log_info "Cairn CLI Installer"
    log_info ""

    # Detect platform
    os=$(detect_os)
    arch=$(detect_arch)
    log_info "Detected platform: ${os}/${arch}"

    # Resolve version
    log_info "Resolving version..."
    version=$(resolve_version)
    log_info "Installing cairn $version"

    # Create temp directory
    TMPDIR_CAIRN=$(mktemp -d)

    # Construct download URLs
    binary_name="cairn-${os}-${arch}"
    binary_url="https://github.com/$GITHUB_REPO/releases/download/$version/$binary_name"
    checksums_url="https://github.com/$GITHUB_REPO/releases/download/$version/checksums.txt"

    # Download checksums
    log_info "Downloading checksums..."
    checksums_file="$TMPDIR_CAIRN/checksums.txt"
    if ! download "$checksums_url" "$checksums_file"; then
        log_info "Verify the release exists: https://github.com/$GITHUB_REPO/releases/tag/$version"
        exit 1
    fi

    # Validate checksums file is not an HTML error page
    if head -1 "$checksums_file" | grep -qi '<html'; then
        log_error "Release $version not found (got HTML instead of checksums)."
        log_info "Available releases: https://github.com/$GITHUB_REPO/releases"
        exit 1
    fi

    # Download binary
    log_info "Downloading cairn ${os}/${arch}..."
    binary_file="$TMPDIR_CAIRN/$binary_name"
    if ! download "$binary_url" "$binary_file"; then
        log_info "Verify the release exists: https://github.com/$GITHUB_REPO/releases/tag/$version"
        exit 1
    fi

    # Verify checksum
    log_info "Verifying checksum..."
    expected_checksum=$(grep "$binary_name" "$checksums_file" | awk '{print $1}')
    if [ -z "$expected_checksum" ]; then
        log_error "Could not find checksum for $binary_name in checksums.txt"
        exit 3
    fi
    verify_checksum "$binary_file" "$expected_checksum"
    log_success "Checksum verified."

    # Install binary
    log_info "Installing to $INSTALL_DIR..."
    install_binary "$binary_file" "$INSTALL_DIR"
    log_success "cairn $version installed to $INSTALL_DIR/cairn"

    # Check PATH
    check_path "$INSTALL_DIR"

    # Vicinae extension
    if detect_vicinae; then
        if prompt_extension; then
            install_extension "$version" "$os"
        else
            log_info "Skipping Vicinae extension installation."
        fi
    fi

    # Desktop app
    if ! desktop_available "$os" "$arch"; then
        if [ "$WITH_DESKTOP" = 1 ]; then
            log_info "The desktop app is not available for $version on ${os}/${arch}; skipping it."
        fi
    elif prompt_desktop; then
        install_desktop "$version" "$os" "$arch" || log_info "Skipping desktop app; the CLI is installed."
    fi

    log_info ""
    log_success "Installation complete! cairn $version installed to $INSTALL_DIR/cairn"
    log_info ""

    # Check if using standard location
    if [ "$INSTALL_DIR" = "/usr/local/bin" ] || [ "$INSTALL_DIR" = "/usr/bin" ]; then
        log_info "✓ $INSTALL_DIR is in your system PATH"
        log_info "  Run 'cairn --help' to get started"
    else
        log_info "Next steps:"
        log_info "  1. Ensure $INSTALL_DIR is in your PATH"
        log_info "  2. Run 'cairn --help' to get started"
        log_info ""
        log_info "To add $INSTALL_DIR to PATH, add this line to your shell profile:"
        log_info "  export PATH=\"$INSTALL_DIR:\$PATH\""
    fi
}

main "$@"

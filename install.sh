#!/bin/sh
# Cairn CLI installer
# Usage: curl -sSL https://raw.githubusercontent.com/ndy40/cairn/main/install.sh | sh
#
# MIT License - https://github.com/ndy40/cairn

set -eu

# ── Defaults ──────────────────────────────────────────────────────────────────

INSTALL_DIR="${CAIRN_INSTALL_DIR:-"$HOME/.local/bin"}"
VERSION=""
NON_INTERACTIVE=0
WITH_EXTENSION=0
GITHUB_REPO="ndy40/cairn"
TMPDIR_CAIRN=""

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
    -d, --install-dir DIR    Installation directory (default: ~/.local/bin)
    -v, --version VERSION    Install a specific version (e.g., v0.0.1)
    -y, --non-interactive    Skip all prompts, install CLI only
        --with-extension     Also install Vicinae extension (non-interactive mode)
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
        curl -fsSL -o "$dest" "$url" 2>/dev/null
    elif has_command wget; then
        wget -qO "$dest" "$url" 2>/dev/null
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
    if has_command curl; then
        tag=$(curl -fsSL "$latest_url" 2>/dev/null | grep '"tag_name"' | head -1 | sed 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/')
    elif has_command wget; then
        tag=$(wget -qO- "$latest_url" 2>/dev/null | grep '"tag_name"' | head -1 | sed 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/')
    else
        log_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi

    if [ -z "$tag" ]; then
        log_error "Could not determine latest release version."
        log_info "Check your internet connection or specify a version with --version."
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
    dest="$dest_dir/cairn"

    # Create install directory if needed
    if ! mkdir -p "$dest_dir" 2>/dev/null; then
        log_error "Permission denied: cannot create directory '$dest_dir'."
        log_info "Try one of the following:"
        log_info "  - Use a different directory: install.sh --install-dir ~/bin"
        log_info "  - Fix permissions: sudo mkdir -p '$dest_dir' && sudo chown \$(whoami) '$dest_dir'"
        exit 4
    fi

    # Backup existing binary
    if [ -f "$dest" ]; then
        cp "$dest" "$dest.bak"
    fi

    # Copy new binary
    if ! cp "$src" "$dest" 2>/dev/null; then
        # Restore backup if copy failed
        if [ -f "$dest.bak" ]; then
            mv "$dest.bak" "$dest"
        fi
        log_error "Permission denied: cannot write to '$dest'."
        log_info "Try one of the following:"
        log_info "  - Use a different directory: install.sh --install-dir ~/bin"
        log_info "  - Fix permissions: sudo chown \$(whoami) '$dest_dir'"
        exit 4
    fi

    chmod +x "$dest"

    # Remove backup on success
    rm -f "$dest.bak"
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
        log_error "Failed to download checksums file."
        log_info "Check your internet connection or verify the version exists:"
        log_info "  https://github.com/$GITHUB_REPO/releases"
        exit 1
    fi

    # Download binary
    log_info "Downloading cairn ${os}/${arch}..."
    binary_file="$TMPDIR_CAIRN/$binary_name"
    if ! download "$binary_url" "$binary_file"; then
        log_error "Failed to download cairn binary."
        log_info "Check your internet connection or verify the version exists:"
        log_info "  https://github.com/$GITHUB_REPO/releases/tag/$version"
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

    log_info ""
    log_success "Installation complete! Run 'cairn --help' to get started."
}

main "$@"

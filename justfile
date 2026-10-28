# https://just.systems

default:
    echo 'Hello, world!'

build:
    go build -o cairn  ./cmd/cairn/

_docs:
    hugo build

update-docs: _docs
    git add docs/
    git commit -m 'docs: Updated documentation'

publish-docs: update-docs
    git push

serve-docs: update-docs
    hugo serve

[working-directory: "vicinae-extension"]
build-ext:
    npm run build

# Desktop app (desktop/, a separate Go module built with Wails v3).
# Requires the wails3 CLI and, on Linux, GTK4 + WebKitGTK 6.0 dev packages.
[working-directory: "desktop"]
build-desktop:
    wails3 build

[working-directory: "desktop"]
dev-desktop:
    wails3 dev

[working-directory: "desktop"]
test-desktop:
    go test ./...

[working-directory: "desktop"]
package-desktop:
    wails3 package

# Throwaway folder for try-desktop, and the bookmarks it copies on first run.
desktop_sandbox := "/tmp/cairn-test"
default_db := if os() == "macos" { env_var("HOME") / "Library/Application Support/cairn/bookmarks.db" } else { env_var_or_default("XDG_DATA_HOME", env_var("HOME") / ".local/share") / "cairn/bookmarks.db" }

# The copy is made once; try-desktop-reset starts again from your current
# bookmarks. On Linux the app's cairn.json, window state, settings and login
# item also stay in the sandbox. Sync is off there, because the sandbox has no
# cairn.json with Dropbox credentials.
#
# Build and run cairn-desktop with debug logging on a sandboxed copy of your bookmarks
try-desktop src_db=default_db: build-desktop
    #!/usr/bin/env sh
    set -eu
    if pgrep -x cairn-desktop >/dev/null; then
        echo "cairn-desktop is already running; quit it first (tray menu, or: pkill cairn-desktop)." >&2
        exit 1
    fi
    sandbox="{{desktop_sandbox}}"
    mkdir -p "$sandbox/config"
    if [ ! -f "$sandbox/bookmarks.db" ] && [ -f "{{src_db}}" ]; then
        # Copy the WAL files too so the copy includes changes not yet checkpointed.
        for suffix in "" -wal -shm; do
            if [ -f "{{src_db}}$suffix" ]; then
                cp "{{src_db}}$suffix" "$sandbox/bookmarks.db$suffix"
            fi
        done
        echo "Copied {{src_db}} to $sandbox"
    fi
    XDG_CONFIG_HOME="$sandbox/config" \
    CAIRN_DB_PATH="$sandbox/bookmarks.db" \
    CAIRN_DESKTOP_DEBUG=1 \
        ./desktop/bin/cairn-desktop

# Delete the try-desktop sandbox (bookmark copy, window state and settings).
try-desktop-reset:
    rm -rf "{{desktop_sandbox}}"

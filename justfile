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

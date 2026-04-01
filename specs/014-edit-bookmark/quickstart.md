# Quickstart: Edit Bookmark

**Feature**: 014-edit-bookmark

## Prerequisites

- Go 1.25.0+
- Node.js (for vicinae-extension only)
- cairn binary built and on PATH

## Build & Test

```bash
# Build cairn
cd /home/ndy40/Development/Go/cairn
go build -o cairn ./cmd/cairn

# Run Go tests
go test ./internal/store/... ./cmd/cairn/...

# Test CLI edit with URL
./cairn add https://example.com
./cairn edit 1 --url=https://new-example.com
./cairn list --json  # verify URL changed

# Test duplicate rejection
./cairn add https://other.com
./cairn edit 1 --url=https://other.com  # should fail with exit 1
```

## Vicinae Extension

```bash
cd vicinae-extension
npm install
npm run build
# Test by loading extension in Raycast dev mode
```

## Key Files to Modify

### Go (Store layer)
- `internal/store/bookmark.go` — Add `URL *string` to `BookmarkPatch`, update `UpdateFields` for URL/domain handling and duplicate check

### Go (CLI)
- `cmd/cairn/main.go` — Add `--url` flag to `cmdEdit`, pass to `runEdit`

### Go (TUI)
- `internal/model/edit.go` — Add URL text input field, Tab navigation between fields
- `internal/model/app.go` — Update `updateEdit` to call `UpdateFields` with URL patch

### TypeScript (Extension)
- `vicinae-extension/src/bm.ts` — Add `bmEdit()` function
- `vicinae-extension/src/bm-list.tsx` — Add Edit action button
- `vicinae-extension/src/bm-search.tsx` — Add Edit action button
- New file: `vicinae-extension/src/bm-edit.tsx` — Edit form component

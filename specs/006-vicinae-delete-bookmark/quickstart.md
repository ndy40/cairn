# Quickstart: Delete Bookmarks from Vicinae Extension

**Feature**: 006-vicinae-delete-bookmark
**Date**: 2026-03-11

## Prerequisites

- Node.js (for vicinae-extension development)
- Vicinae launcher installed with `vici` CLI
- `cairn` CLI built and in PATH

## Build & Test

```bash
# Navigate to extension
cd vicinae-extension

# Install dependencies (if needed)
npm install

# Development mode with hot-reload
vici develop

# Lint
vici lint

# Format
npx biome format --write src
```

## New Files

None — all changes are to existing files.

## Modified Files

| File           | Change                                               |
|----------------|------------------------------------------------------|
| `src/bm.ts`    | Add `bmDelete(id)` function                          |
| `src/bm-list.tsx` | Add Delete action to action panel, add list refresh |
| `src/bm-search.tsx` | Add Delete action to action panel, add list refresh |

## Testing Strategy

1. **Manual test in Vicinae launcher**:
   - Open List Bookmarks command
   - Select a bookmark → action panel should show "Delete Bookmark"
   - Choose "Delete Bookmark" → confirmation dialog appears
   - Cancel → bookmark remains
   - Confirm → bookmark removed, success toast shown, list refreshes
2. **Test search view**: Same flow in Search Bookmarks command
3. **Test error handling**: Delete a bookmark, then try deleting same ID again (should show "not found" error)
4. **Test with no bookmarks**: Empty list should show empty view, no delete actions available

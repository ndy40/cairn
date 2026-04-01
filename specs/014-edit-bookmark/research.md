# Research: Edit Bookmark

**Feature**: 014-edit-bookmark
**Date**: 2026-04-01

## Technical Context — No Unknowns

All technologies are established in the project. No NEEDS CLARIFICATION items from the spec.

## Decisions

### D1: URL duplicate detection during edit

- **Decision**: Check for duplicate URLs before updating, excluding the bookmark being edited (self-match is allowed).
- **Rationale**: The store already has `ExistsByURL(rawURL)` and a UNIQUE constraint on the `url` column. However, `ExistsByURL` doesn't exclude the current bookmark. For edit, we need a query like `SELECT COUNT(1) FROM bookmarks WHERE url = ? AND id != ?`. The UNIQUE constraint on UPDATE will also catch this, but an explicit check provides a clearer error message.
- **Alternatives considered**: (a) Rely solely on UNIQUE constraint — rejected because the SQLite error message is generic and harder to map to a user-friendly message. (b) Remove UNIQUE constraint — rejected, it's a valuable safety net.

### D2: Domain recalculation on URL change

- **Decision**: Reuse the existing `extractDomain(rawURL)` function (internal to store package) when URL is updated.
- **Rationale**: Domain is derived from URL. When URL changes, domain must be recalculated to stay consistent with list/search display.
- **Alternatives considered**: Make domain user-editable — rejected, domain is always derived from URL.

### D3: TUI edit panel — field layout

- **Decision**: Add a URL text input above the existing tags input in the TUI EditModel. Use Tab key to move focus between fields.
- **Rationale**: Matches existing patterns in the codebase (charmbracelet/bubbles textinput). URL is the more prominent field and should appear first.
- **Alternatives considered**: (a) Separate edit modes for URL vs tags — rejected, unnecessary complexity. (b) Single-field form with field selector — rejected, two fields are manageable.

### D4: Vicinae extension edit pattern

- **Decision**: Add a `bmEdit(id, url?, tags?)` function to `bm.ts` that calls `cairn edit <id> [--url=<url>] [--tags=<tags>]`. Add an Edit action button to the list/search views that opens a pre-filled form.
- **Rationale**: Follows the established pattern of bmDelete/bmPin/bmAdd — thin wrapper calling cairn CLI. The add form (`bm-add.tsx`) provides the UX template for the edit form.
- **Alternatives considered**: Inline editing in the list view — rejected, Raycast doesn't support inline editing natively; a form is the standard pattern.

### D5: No schema migration needed

- **Decision**: No DDL changes. The `url` and `domain` columns already exist and are writable. `BookmarkPatch` gains a `URL` field but this is a Go struct change, not a schema change.
- **Rationale**: The existing schema already supports updating these columns. The `UpdateFields` method dynamically builds the SET clause, so adding URL/domain to the clause is sufficient.

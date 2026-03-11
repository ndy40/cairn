# Research: Delete Bookmarks from Vicinae Extension

**Feature**: 006-vicinae-delete-bookmark
**Date**: 2026-03-11

## Decision 1: Confirmation Dialog API

**Decision**: Use `confirmAlert()` from `@vicinae/api` with `Alert.ActionStyle.Destructive` for the confirm button.

**Rationale**: The Vicinae API provides a built-in `confirmAlert()` function that returns `Promise<boolean>`. It supports a destructive action style which visually signals danger (red button). This is the idiomatic approach for the platform.

**Alternatives considered**:
- Custom modal component: Rejected — unnecessary complexity when a built-in API exists.
- No confirmation (direct delete): Rejected — spec requires confirmation to prevent accidental deletions (FR-003).

## Decision 2: Success/Error Feedback

**Decision**: Use `showToast()` from `@vicinae/api` with `Toast.Style.Success` for successful deletions and `Toast.Style.Failure` for errors.

**Rationale**: The existing extension already uses `showHUD()` for clipboard copy feedback. `showToast()` is more appropriate for delete operations because it supports styled feedback (success/failure) and doesn't require the same fire-and-forget semantics as HUD.

**Alternatives considered**:
- `showHUD()`: Too simple — no error styling, meant for brief confirmations.
- Inline error state in the list: Rejected — overly complex for a transient notification.

## Decision 3: CLI Wrapper Pattern

**Decision**: Add a `bmDelete(id: number)` function to `bm.ts` following the existing `bmAdd()` pattern, returning `{ exitCode: number; stderr: string }`.

**Rationale**: Consistent with the established pattern in the codebase (`bmList`, `bmSearch`, `bmAdd`). Returns exit code and stderr for error handling.

**Alternatives considered**:
- Inline spawnSync in the component: Rejected — violates the existing separation of concerns where all CLI calls go through `bm.ts`.

## Decision 4: List Refresh After Deletion

**Decision**: After a successful delete, re-call `bmList()` or `bmSearch(query)` and update the React state with the new results.

**Rationale**: The simplest approach — reload the full list from the CLI. The bookmark count is small enough that a full reload is instantaneous. This avoids maintaining a separate optimistic-delete state that could drift from the actual database.

**Alternatives considered**:
- Client-side filter (remove item from state without reloading): Faster UI update but could show stale data if another process modified bookmarks.
- Optimistic update + rollback: Over-engineered for this use case.

## Decision 5: Shared BookmarkListItem Component

**Decision**: The `BookmarkListItem` component is currently duplicated in `bm-list.tsx` and `bm-search.tsx`. The delete action will need an `onDelete` callback prop to trigger list refresh from the parent. Rather than adding a third copy, extract the shared component or accept the duplication and add the delete action to both copies.

**Rationale**: The component is small (~30 lines). Adding the callback makes the two copies slightly different (list vs search refresh behavior), so keeping them separate is acceptable. Alternatively, a shared component with an `onDelete` callback could be extracted — but this is a minor refactoring decision.

**Decision**: Keep separate copies in each file, add `onDelete` callback to both. This maintains the current code style and avoids adding new files.

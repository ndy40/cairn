# Data Model: Delete Bookmarks from Vicinae Extension

**Feature**: 006-vicinae-delete-bookmark
**Date**: 2026-03-11

## No Schema Changes

This feature requires no database schema changes. It uses the existing `cairn delete <id>` CLI command which operates on the existing `bookmarks` table.

## Existing Entity: Bookmark (TypeScript interface)

The existing `Bookmark` interface in `bm.ts` already includes the `id` field needed for deletion:

| Field         | Type             | Used for Delete |
|---------------|------------------|-----------------|
| `id`          | number           | Yes — passed to `cairn delete <id>` |
| `URL`         | string           | No              |
| `Domain`      | string           | No              |
| `Title`       | string           | Yes — shown in confirmation dialog |
| Other fields  | various          | No              |

## CLI Delete Contract

| Aspect     | Detail                       |
|------------|------------------------------|
| Command    | `cairn delete <id>`          |
| Exit code 0| Successfully deleted         |
| Exit code 1| Bookmark not found           |
| Exit code 3| Error (database, arguments)  |
| Stderr     | Error message on failure     |

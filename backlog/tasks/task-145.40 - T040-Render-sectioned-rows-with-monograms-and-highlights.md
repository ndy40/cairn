---
id: TASK-145.40
title: T040 Render sectioned rows with monograms and highlights
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:58:00Z'
updated_date: '2026-10-27T04:18:00Z'
labels:
  - desktop
  - us6
dependencies: []
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 148000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T040 in specs/018-native-desktop-window/tasks.md.

Rewrite `render()` in `desktop/frontend/src/main.ts`: empty query → "Pinned" and "Recent" section headings (non-selectable `role="presentation"` items) with selection skipping them; rows with a 28px monogram tile (first letter of the domain, tint from a stable hash of the domain over a fixed accessible palette), title with case-insensitive query-term matches wrapped in `<b>` built with `textContent` only, domain, first 3 tags then "+N", an SVG pin glyph with `aria-label="Pinned"`, and an "↩ Open" hint on the selected row; empty states "No bookmarks yet" + add hint, and `No bookmarks match "<query>"` + "Type + before a URL to add it."; keep the selection by id across refreshes
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 An empty query shows Pinned then Recent sections
- [ ] #2 Rows show a monogram, the title with matches in bold, the domain, up to 3 tags then +N, a pin glyph and an Open hint on the selected row
- [ ] #3 Empty states give one line and the next step
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
ndw-rows, merged into 018-native-desktop-window (0bedb24). Code complete; awaiting visual/contrast QA (rows 22-24). Deviations: pure white text and outlined chips on macOS selection, --text-2-strong on Linux tints, for 4.5:1.
<!-- SECTION:NOTES:END -->

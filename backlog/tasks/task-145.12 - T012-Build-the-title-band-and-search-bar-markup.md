---
id: TASK-145.12
title: T012 Build the title band and search bar markup
status: In Progress
assignee:
  - '@claude'
created_date: '2026-10-27T03:57:00Z'
updated_date: '2026-10-27T04:08:00Z'
labels:
  - desktop
  - us1
dependencies:
  - TASK-145.8
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 120000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T012 in specs/018-native-desktop-window/tasks.md.

Restructure the top of `desktop/frontend/index.html`: a `<header class="titleband">` containing a `<label>`-wrapped `<input type="search" id="query" aria-label="Search bookmarks">`, an Add `<button id="add-btn" aria-label="Add bookmark">` with an inline SVG plus icon, and a main-menu `<button id="menu-btn" aria-label="Main menu" aria-haspopup="menu">` (shown on Linux only); keep `#results` as `role="listbox"`; remove the old placeholder text "(type + to add)"
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The search field has an accessible label and sits in the title band (macOS) or the bar under the title bar (Linux)
- [ ] #2 Add and main-menu buttons are real buttons with accessible names; the main menu button shows on Linux only
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implemented in ndw-core, merged in 176f4f7. Code complete; awaiting manual QA (quickstart.md rows 1-10) before ticking criteria.
<!-- SECTION:NOTES:END -->

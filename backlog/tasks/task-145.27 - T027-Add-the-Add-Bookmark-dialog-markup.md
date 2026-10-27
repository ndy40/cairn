---
id: TASK-145.27
title: T027 Add the Add Bookmark dialog markup
status: To Do
assignee: []
created_date: '2026-10-27T03:58:00Z'
labels:
  - desktop
  - us4
dependencies:
  - TASK-145.12
references:
  - specs/018-native-desktop-window/tasks.md
parent_task_id: TASK-145
ordinal: 135000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
T027 in specs/018-native-desktop-window/tasks.md.

In `desktop/frontend/index.html`, replace the `#add-view` form with `<dialog id="add-dialog">` containing a heading "Add Bookmark", a labelled URL input, a labelled Tags input ("comma separated"), helper text "Title and description are fetched automatically.", and Cancel/Add buttons (Add is `type="submit"`, default)
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 The add form is a modal dialog with labelled URL and Tags fields, helper text, and Cancel then Add buttons
<!-- AC:END -->

---
id: TASK-12
title: >-
  Implement root bubbletea Model and application state machine in
  internal/model/app.go
status: Done
assignee: []
created_date: '2026-03-06 04:05'
updated_date: '2026-03-06 17:10'
labels:
  - 'phase:foundation'
dependencies: []
priority: high
ordinal: 89000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create the top-level bubbletea model that owns the application state machine (browse, search, add, confirm-delete modes), routes messages to the active sub-model, and composes the final view with a persistent footer.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 AppState type is defined with constants: StateBrowse, StateSearch, StateAdd, StateConfirmDelete
- [ ] #2 Root Model struct holds: current AppState, store reference, and sub-model fields (BrowseModel, SearchModel, AddModel)
- [ ] #3 Init() returns a tea.Cmd that loads bookmarks from the store on startup
- [ ] #4 Update() handles tea.KeyMsg{ctrl+c} globally to quit in any state
- [ ] #5 View() delegates to the active sub-model view and appends the mode-appropriate footer
- [ ] #6 Running ./bm launches, shows 'No bookmarks' or bookmark list, and quits cleanly on Ctrl+C
<!-- AC:END -->

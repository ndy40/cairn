---
id: TASK-13
title: Implement add bookmark modal sub-model in internal/model/add.go
status: Done
assignee: []
created_date: '2026-03-06 04:05'
updated_date: '2026-03-06 17:10'
labels:
  - 'story:us1'
dependencies: []
priority: high
ordinal: 53000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create the AddModel sub-model that handles the add-bookmark modal: a text input field pre-filled with the clipboard URL, status messages for errors, and keyboard handling for save/cancel.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 AddModel contains a bubbles/textinput.Model that is focused when the modal is shown
- [ ] #2 Update() handles Enter key to emit a save-confirmation message with the current URL value
- [ ] #3 Update() handles Escape key to emit a cancel message
- [ ] #4 Update() passes all other keypresses to the textinput for editing
- [ ] #5 URL() method returns the current value of the text input
- [ ] #6 SetStatus(msg) updates a status string that is rendered inside the modal for error/info feedback
<!-- AC:END -->

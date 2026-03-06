# Contract: Keyboard Shortcuts

**Feature**: 001-tui-bookmark-manager
**Date**: 2026-03-06

This document defines the complete keyboard interface contract for the TUI application. All shortcuts are active in the specified modes and must be consistently displayed in the status bar footer.

---

## Application Modes

| Mode | Description |
|------|-------------|
| `browse` | Default mode: scrollable bookmark list, no active text input |
| `search` | Search input is focused; list updates in real time as user types |
| `add` | Modal overlay is active for adding a new bookmark |
| `confirm-delete` | Delete confirmation dialog is active |

---

## Global Shortcuts (active in all modes)

| Key | Action | Notes |
|-----|--------|-------|
| `Ctrl+C` | Quit application | Always available; terminates without confirmation |
| `Ctrl+P` | Add bookmark from clipboard | Opens `add` modal with clipboard URL pre-filled; active in `browse` and `search` modes |
| `?` | Toggle help overlay | Shows full shortcut reference; closes on any key |

---

## Browse Mode Shortcuts

| Key | Action | Notes |
|-----|--------|-------|
| `↑` / `k` | Move selection up | Vim-style `k` also supported |
| `↓` / `j` | Move selection down | Vim-style `j` also supported |
| `Enter` | Open selected bookmark | Opens URL in system default browser |
| `d` or `Delete` | Delete selected bookmark | Transitions to `confirm-delete` mode |
| `/` | Focus search input | Transitions to `search` mode |
| `g` | Jump to top of list | |
| `G` | Jump to bottom of list | |

---

## Search Mode Shortcuts

| Key | Action | Notes |
|-----|--------|-------|
| Any printable char | Append to search term | List filters in real time |
| `Backspace` | Delete last character of search term | |
| `Ctrl+A` | Clear entire search term | Returns to full list |
| `↑` / `↓` | Move selection within filtered list | |
| `Enter` | Open selected bookmark | Opens URL in system default browser |
| `Escape` | Clear search and return to browse mode | Full list restored |

---

## Add Mode Shortcuts (modal)

| Key | Action | Notes |
|-----|--------|-------|
| `Enter` | Confirm and save bookmark | Triggers page fetch, then saves |
| `Escape` | Cancel and close modal | No bookmark saved |
| `Backspace` | Edit URL in input field | User can modify the pasted URL before saving |

---

## Confirm-Delete Mode Shortcuts (dialog)

| Key | Action | Notes |
|-----|--------|-------|
| `y` / `Enter` | Confirm deletion | Bookmark permanently removed |
| `n` / `Escape` | Cancel deletion | Returns to browse mode, bookmark unchanged |

---

## Status Bar Content

The persistent footer always shows the shortcuts relevant to the current mode:

| Mode | Footer Content |
|------|----------------|
| `browse` | `[/] Search  [Ctrl+P] Add  [Enter] Open  [d] Delete  [?] Help  [Ctrl+C] Quit` |
| `search` | `[Esc] Clear  [Enter] Open  [Ctrl+P] Add  [Ctrl+C] Quit` |
| `add` | `[Enter] Save  [Esc] Cancel` |
| `confirm-delete` | `[y/Enter] Confirm Delete  [n/Esc] Cancel` |

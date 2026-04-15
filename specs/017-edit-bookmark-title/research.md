# Research: Edit Bookmark Title (017)

No external research was needed. All decisions were resolved directly from the existing codebase.

---

## Decision 1: Where to add the title field in the TUI

**Decision**: Add a `titleInput textinput.Model` alongside the existing `urlInput` and `tagsInput` in `internal/model/edit.go`, using the same charmbracelet/bubbles `textinput` component already in use.

**Rationale**: The `EditModel` struct already manages two `textinput.Model` fields (`urlInput`, `tagsInput`) with a simple integer `activeField` index driving Tab/Shift+Tab navigation. Adding a third field (`editFieldTitle = 2`) follows the exact same pattern with zero new dependencies or abstractions. The title currently renders as a read-only bold header — repurposing that position as an editable input is the most natural UX placement.

**Alternatives considered**:
- A separate "rename" keybinding outside the edit panel: rejected — breaks discoverability and makes the flow inconsistent with URL and tags editing.
- A modal prompt overlay: rejected — adds complexity not present elsewhere in the TUI.

---

## Decision 2: Tab order for the new title field

**Decision**: Title field is first in tab order — `Title → URL → Tags`.

**Rationale**: Title is the field users most commonly need to correct when fetching fails. Placing it first means a single Enter (or Tab) sequence gets them there immediately. URL and tags are secondary corrections. This matches the visual layout: title appears above the URL in the current view.

**Alternatives considered**:
- Title last (`URL → Tags → Title`): inconsistent with visual layout order.
- Title second (`URL → Title → Tags`): no clear advantage; URL and tags are already adjacent and users are used to that pairing.

---

## Decision 3: Validation approach in the TUI

**Decision**: Validate on save (Enter key) inside the `EditModel`. If `titleInput.Value()` is empty after trimming, set an error string on the model and return without emitting a save message. Display the error inline below the title field.

**Rationale**: The existing TUI has no form-level validation infrastructure — errors are handled at the `App` level via `err string` state. For a field-level constraint (non-empty title), inline validation at the model level is the simplest approach and avoids threading error state through the parent `App` model.

**Alternatives considered**:
- Disabling the save action when title is empty: requires observable state in the key handler — more invasive.
- Validating at the `App` level: leaks field-specific knowledge into the parent.

---

## Decision 4: `bmEdit()` signature change in the extension

**Decision**: Add an optional `title?: string` parameter to `bmEdit()` in `bm.ts`. When provided and non-empty, append `--title <value>` to the CLI args.

**Rationale**: `bmEdit()` already has an identical optional pattern for `url` and `tags`. The change is a one-line conditional — no refactor required, no interface breaks. The call site in `bm-edit.tsx` passes title only when it differs from the original, mirroring the existing URL and tags change-detection pattern.

**Alternatives considered**:
- A new `bmEditTitle()` function: unnecessary duplication; the single function approach keeps the call site simple.

---

## Decision 5: Maximum title length (500 characters)

**Decision**: 500 characters, enforced via `CharLimit` on the `textinput.Model` (TUI) and `maxLength` on the `Form.TextField` (extension).

**Rationale**: No existing bookmarks exceed this — the column has no explicit SQLite constraint so the application layer is the right enforcement point. 500 chars accommodates even unusually long academic paper titles. The URL field already uses a 2048-char limit; titles are consistently shorter.

**Alternatives considered**:
- 255 chars: too short for some real-world page titles.
- No limit: could produce awkward display in list view.

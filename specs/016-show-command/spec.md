# Feature 016: `cairn show` — Metadata Summary Command

## Overview

Add a `show` subcommand to the `cairn` CLI that prints a concise summary of the current installation state. The output helps users quickly verify their setup without having to run multiple commands.

---

## User Stories

### US1 — P1: Display database path and bookmark count

**As a** cairn user,  
**I want** to run `cairn show` and see the path to the active database and how many bookmarks are stored,  
**So that** I can confirm I'm using the right database and have a sense of its contents.

**Acceptance Criteria**:
- `cairn show` prints `Database:   <resolved db path>`.
- `cairn show` prints `Bookmarks:  <n>` where `<n>` is the total row count in the bookmarks table.
- The count reflects the live database state (not a cached value).

---

### US2 — P2: Display sync status and last sync date

**As a** cairn user,  
**I want** `cairn show` to also report whether sync is configured and when it last ran,  
**So that** I can tell at a glance whether my data is being backed up and is up to date.

**Acceptance Criteria**:
- When sync is **not** configured, `cairn show` prints `Sync:       not configured`.
- When sync **is** configured, `cairn show` prints `Sync:       configured (dropbox)` (or the relevant backend name).
- When sync is configured and has run at least once, `cairn show` prints `Last sync:  <date-time in UTC>`.
- When sync is configured but has never run, `cairn show` prints `Last sync:  never`.
- The sync fields are omitted (or clearly labelled "not configured") when `LoadConfig` returns nil.

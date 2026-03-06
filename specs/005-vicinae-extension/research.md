# Research: Vicinae Extension for Bookmark Manager

**Feature**: 005-vicinae-extension
**Date**: 2026-03-06

---

## Decision 1: Extension Technology Stack

**Decision**: Use TypeScript with the `@vicinae/api` SDK. Build with `vici build` / `vici develop`. The extension is a new package in `vicinae-extension/` at the repository root.

**Rationale**: Vicinae extensions are exclusively TypeScript/React. The SDK (`@vicinae/api`) is the only supported way to build Vicinae extensions; it provides the `List`, `Form`, `Detail`, and `Action` components used by all platform extensions. The Vicinae API is documented as "mostly compatible" with Raycast, so established Raycast extension patterns apply directly.

**Alternatives considered**:
- Script Commands (shell scripts registered in Vicinae): Simpler but provide no rich UI — no list navigation, no form inputs, no inline search. Rejected: the spec requires real-time search results and a form for adding bookmarks.
- dmenu mode: Text-only stdin/stdout piping. No interactive form. Rejected.

---

## Decision 2: CLI Invocation Pattern

**Decision**: Execute `bm` commands via Node.js `child_process.spawnSync` (synchronous) for simple one-shot calls (list, add) and via streaming/async exec for real-time search updates as the user types.

**Rationale**: Vicinae extensions run in a Node.js-compatible environment. `spawnSync` avoids shell injection by passing arguments as an array (no shell interpolation). For the search-as-you-type flow, the query is debounced and a fresh `bm search <query> --json` call is made on each change.

**Alternatives considered**:
- `execSync` with shell string interpolation: Vulnerable to shell injection if the query contains special characters. Rejected.
- Direct database access from the extension: Violates FR-011 (all access must go through the CLI). Rejected.

---

## Decision 3: CLI Output Format

**Decision**: Use `--json` flag for machine-readable output from both `bm list --json` and `bm search <query> --json`. The JSON array includes all Bookmark fields (id, title, url, domain, tags, created_at).

**Rationale**: The `--json` flag is already implemented in the `bm` CLI and outputs a properly structured JSON array via `json.Encode`. This gives the extension typed data without parsing tab-separated text.

**Alternatives considered**:
- Parse tab-separated output from `bm list`: Fragile if titles contain tabs. Rejected.

---

## Decision 4: CLI Prerequisite — bm add --tags Flag

**Decision**: The `bm add` CLI command must be extended with an optional `--tags` flag before the extension's Add Bookmark form can pass tags. This is an intra-feature prerequisite.

**Rationale**: The current `bm add <url>` command hard-codes `nil` for the tags parameter (see `cmd/bm/main.go:runAdd`). The spec requires the Add form to support tags. The fix is minimal: add `--tags` string flag to the `add` FlagSet in `main.go`, split on comma, and pass to `s.Insert()`.

**Impact**: The Go CLI (`cmd/bm/main.go`) requires a small change as part of this feature. This is one additional source file modified versus the extension-only scope originally assumed.

**Alternatives considered**:
- Ship the extension without tags support initially: Possible but contradicts FR-007. Rejected.
- Pass tags as a second positional argument: Non-standard, harder to use. Rejected.

---

## Decision 5: Real-Time Search Architecture

**Decision**: In the Search command, call `bm search <query> --json` each time the search input changes (debounced). If the query is empty, fall back to `bm list --json` to show all bookmarks.

**Rationale**: This matches the spec requirement that results filter in real time. Since `bm search` already does FTS5 + fuzzy ranking server-side, the extension simply delegates to the CLI. When the query is empty, showing the full list (via `bm list --json`) is the most natural behaviour.

**Alternatives considered**:
- Load all bookmarks once and filter client-side in the extension: Would diverge from the existing fuzzy ranking logic. Rejected.

---

## Decision 6: Clipboard Pre-fill for Add Form

**Decision**: On load of the Add Bookmark command, read the system clipboard. If the clipboard content starts with `http://` or `https://`, pre-fill the URL field with it.

**Rationale**: The Vicinae API (Raycast-compatible) provides a clipboard utility. Pre-filling removes the copy-paste step for the most common add-bookmark flow. Validation before pre-fill prevents garbage data appearing in the URL field.

**Alternatives considered**:
- Always pre-fill regardless of URL validity: Would confuse users who have non-URL text on clipboard. Rejected.

---

## Files Changed

| Location | File | Change |
|----------|------|--------|
| `cmd/bm/main.go` | existing | Add `--tags` flag to `bm add` subcommand |
| `vicinae-extension/` | new package | Three-command TypeScript extension |
| `vicinae-extension/package.json` | new | Extension manifest with commands + `@vicinae/api` dependency |
| `vicinae-extension/src/search.tsx` | new | Search Bookmarks command |
| `vicinae-extension/src/list.tsx` | new | List Bookmarks command |
| `vicinae-extension/src/add.tsx` | new | Add Bookmark command |
| `vicinae-extension/tsconfig.json` | new | TypeScript configuration |

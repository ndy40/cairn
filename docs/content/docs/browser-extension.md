---
title: "Vicinae Extension"
weight: 60
---

# Vicinae Extension

Cairn includes a [Vicinae](https://vicinae.com) extension that lets you save, search, and browse bookmarks directly from the Vicinae launcher — without leaving your current context.

## Features

- **Search** — find bookmarks from the Vicinae command palette
- **Browse** — list all bookmarks in a panel
- **Save** — add the current browser tab as a bookmark (with optional tags)

## Installation

### Via Installer (Recommended)

If Vicinae is installed, run:

```sh
sh install.sh --with-extension
```

Or on a fresh install:

```sh
curl -sSL https://raw.githubusercontent.com/ndy40/cairn/main/install.sh | sh -s -- -y --with-extension
```

### Manual

The extension source is in `vicinae-extension/` in the repository. Build with:

```sh
cd vicinae-extension
npm install
npm run build
```

Then register the extension with the `vici` toolchain.

## Requirements

- [Vicinae](https://vicinae.com) installed and running
- `cairn` binary in `PATH`

The extension delegates all operations to the `cairn` CLI — no separate daemon or server is needed.

## Usage

Open the Vicinae launcher and use the `cairn` prefix to interact:

| Action | Example |
|--------|---------|
| Search bookmarks | Type `cairn golang` |
| Browse all bookmarks | Open the Cairn panel |
| Save current page | Use the "Save to Cairn" action |

## Commands Invoked

The extension calls these CLI commands internally:

| Extension action | CLI command |
|-----------------|-------------|
| Search | `cairn search <query> --json` |
| List | `cairn list --json` |
| Save | `cairn add <url> --tags <tags>` |
| Delete | `cairn delete <id>` |
| Pin | `cairn pin <id>` |

# Cairn Vicinae Extension

A [Vicinae](https://vicinae.com) extension for managing bookmarks through the Cairn CLI.

## Features

- **Search Bookmarks** — search your saved bookmarks and open them in the browser
- **List Bookmarks** — browse all bookmarks with details
- **Add Bookmark** — save the current page as a bookmark

All commands work by calling the `cairn` CLI under the hood, so the Cairn CLI must be installed and available on your PATH.

## Prerequisites

- [Cairn CLI](../README.md#installation) installed and on your PATH
- [Vicinae](https://vicinae.com) installed

## Installation

### Via the install script

If you already have Vicinae installed, the Cairn install script can set up the extension automatically:

```sh
curl -sSL https://raw.githubusercontent.com/ndy40/cairn/main/install.sh | sh
```

The installer detects Vicinae and prompts you to install the extension.

### Manual installation

```sh
cd vicinae-extension
npm install
npm run dev
```

This starts the extension in development mode using `vici develop`.

## Development

```sh
npm install       # Install dependencies
npm run dev       # Start in development mode
npm run build     # Build for production
npm run lint      # Run linter
npm run format    # Format source code
```

## How it works

The extension uses Node.js `child_process` to call the `cairn` CLI with the `--json` flag for structured output. Results are rendered using React components via the Vicinae API.

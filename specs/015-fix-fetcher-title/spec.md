# Feature 015: Fix Fetcher Title Extraction

## Overview

The `Fetch` function in `internal/fetcher/fetcher.go` uses `extractTitle` to obtain a bookmark title from a web page. Currently, the function prioritises `og:title` over the standard HTML `<title>` tag. For sites like YouTube that serve limited or bot-blocked HTML, neither tag resolves and the function falls back to the hostname — so a YouTube video bookmark is saved as `youtube.com` instead of the video title.

This feature addresses two issues:
1. **Priority change**: The HTML `<title>` tag should take priority over `og:title`, since it is present in all well-formed HTML pages and is the canonical title visible to users.
2. **YouTube / bot-detection resilience**: The fetcher should attempt to work around pages that serve minimal HTML to plain HTTP clients by improving the request headers (Accept, Accept-Language) and handling redirect flows (e.g., YouTube consent pages).

---

## User Stories

### US1 — P1: HTML `<title>` tag is preferred over `og:title` (priority swap)

**As a** bookmark user,  
**I want** the bookmark title to use the HTML `<title>` tag when it is available,  
**So that** the title shown matches what I see in the browser tab, not a marketing/social meta tag.

**Acceptance Criteria**:
- When a page has both a `<title>` tag and an `og:title` meta tag, `extractTitle` returns the `<title>` value.
- When a page has only `og:title` (no `<title>`), `extractTitle` returns the `og:title` value as before.
- When neither is present, `extractTitle` returns the fallback hostname.

---

### US2 — P2: YouTube and bot-gated pages return a meaningful title

**As a** bookmark user,  
**I want** saving a YouTube (or similar JS-heavy / bot-gated) URL to capture the video title,  
**So that** I do not end up with `youtube.com` as the bookmark name.

**Acceptance Criteria**:
- Fetching a YouTube watch URL returns the video title (or at worst the channel name), not `youtube.com`.
- The fetcher passes realistic browser-like headers (`User-Agent`, `Accept`, `Accept-Language`) to reduce bot blocking.
- If the server returns a non-2xx status code, the error path still returns the fallback hostname alongside the error.
- The change does not regress title extraction for regular web pages.

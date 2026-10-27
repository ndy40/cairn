import type {Bookmark} from "../bindings/github.com/ndy40/cairn/desktop";

// Rows show up to this many tags; the rest collapse into a "+N" chip.
const maxTags = 3;
// Number of monogram tints defined in style.css (.tint-0 … .tint-5).
const tintCount = 6;

const svgNS = "http://www.w3.org/2000/svg";

// Letters and digits in any script. Built with RegExp so tsc's default (ES5)
// target accepts the "u" flag.
const edgePunct = new RegExp("^[^\\p{L}\\p{N}]+|[^\\p{L}\\p{N}]+$", "gu");
const letterOrDigit = new RegExp("[\\p{L}\\p{N}]", "u");

export interface RowHandlers {
    // Pointer moved over result i.
    hover(i: number): void;
    // Result i was clicked.
    open(i: number): void;
}

// renderResults rebuilds the result list. With an empty query the results
// (pinned first, then recent, as Search returns them) are split under
// "Pinned" and "Recent" headings. Headings are role="presentation" items, so
// indices passed to the handlers are always indices into results.
//
// All bookmark text is set via textContent: titles come from fetched web pages
// and must never be interpreted as HTML.
export function renderResults(list: HTMLUListElement, results: Bookmark[], selected: number, query: string, on: RowHandlers) {
    list.replaceChildren();
    setCount(results.length);

    if (results.length === 0) {
        list.append(emptyState(query));
        return;
    }

    const match = matcher(query);
    const sectioned = query === "";
    let section = "";

    results.forEach((b, i) => {
        if (sectioned) {
            const next = b.pinned ? "Pinned" : "Recent";
            if (next !== section) {
                section = next;
                list.append(heading(section));
            }
        }
        const li = row(b, i === selected, match);
        li.addEventListener("mousemove", () => on.hover(i));
        li.addEventListener("click", () => {
            on.hover(i);
            on.open(i);
        });
        list.append(li);
    });
}

function heading(text: string): HTMLLIElement {
    const li = document.createElement("li");
    li.className = "section";
    li.setAttribute("role", "presentation");
    li.textContent = text;
    return li;
}

function row(b: Bookmark, selected: boolean, match: RegExp | null): HTMLLIElement {
    const li = document.createElement("li");
    li.className = "result" + (selected ? " is-selected" : "");
    li.id = `bookmark-${b.id}`;
    li.setAttribute("role", "option");
    li.setAttribute("aria-selected", String(selected));

    const domain = b.domain || hostOf(b.url);

    const tile = document.createElement("span");
    tile.className = `monogram tint-${hash(domain) % tintCount}`;
    tile.setAttribute("aria-hidden", "true");
    tile.textContent = monogram(domain);

    const text = document.createElement("span");
    text.className = "text";

    const title = document.createElement("span");
    title.className = "title";
    appendHighlighted(title, b.title || b.url, match);
    title.title = b.title || b.url;

    const meta = document.createElement("span");
    meta.className = "meta";
    const dom = document.createElement("span");
    dom.className = "domain";
    dom.textContent = domain;
    meta.append(dom);

    const tags = b.tags ?? [];
    for (const tag of tags.slice(0, maxTags)) {
        meta.append(chip(tag));
    }
    if (tags.length > maxTags) {
        const more = chip(`+${tags.length - maxTags}`);
        more.title = tags.slice(maxTags).join(", ");
        more.setAttribute("aria-label", `${tags.length - maxTags} more tags: ${more.title}`);
        meta.append(more);
    }

    text.append(title, meta);
    li.append(tile, text);

    if (b.pinned) {
        li.append(pinGlyph());
    }
    if (selected) {
        li.append(openHint());
    }
    return li;
}

function chip(text: string): HTMLSpanElement {
    const span = document.createElement("span");
    span.className = "tag";
    span.textContent = text;
    return span;
}

function pinGlyph(): SVGSVGElement {
    const svg = document.createElementNS(svgNS, "svg");
    svg.setAttribute("class", "pin");
    svg.setAttribute("viewBox", "0 0 24 24");
    svg.setAttribute("width", "14");
    svg.setAttribute("height", "14");
    svg.setAttribute("role", "img");
    svg.setAttribute("aria-label", "Pinned");
    const path = document.createElementNS(svgNS, "path");
    path.setAttribute("d", "M15 3l6 6-3 1-4 4 1 5-2 2-4-4-5 5-1-1 5-5-4-4 2-2 5 1 4-4z");
    path.setAttribute("fill", "currentColor");
    svg.append(path);
    return svg;
}

// openHint marks the selected row; Enter is announced by the listbox already,
// so the hint is hidden from assistive tech.
function openHint(): HTMLSpanElement {
    const hint = document.createElement("span");
    hint.className = "open-hint";
    hint.setAttribute("aria-hidden", "true");
    const key = document.createElement("kbd");
    key.textContent = "↩";
    hint.append(key, "Open");
    return hint;
}

function emptyState(query: string): HTMLLIElement {
    const li = document.createElement("li");
    li.className = "empty";
    li.setAttribute("role", "presentation");

    const head = document.createElement("strong");
    const next = document.createElement("span");
    if (query) {
        head.textContent = `No bookmarks match "${query}"`;
        next.append("Type ", code("+"), " before a URL to add it.");
    } else {
        head.textContent = "No bookmarks yet";
        const mac = document.documentElement.dataset.os === "darwin";
        const key = document.createElement("kbd");
        key.textContent = mac ? "⌘N" : "Ctrl+N";
        next.append("Press ", key, " to add one, or run ", code("cairn add <url>"), ".");
    }
    li.append(head, next);
    return li;
}

function code(text: string): HTMLElement {
    const el = document.createElement("code");
    el.textContent = text;
    return el;
}

// setCount updates the result count in the footer (macOS and KDE).
function setCount(n: number) {
    const el = document.getElementById("result-count");
    if (el) {
        el.textContent = n === 0 ? "" : n === 1 ? "1 result" : `${n} results`;
    }
}

// matcher builds one case-insensitive pattern for the query's terms, longest
// first so "fts5" wins over "fts". Search syntax characters are dropped.
function matcher(query: string): RegExp | null {
    const terms = Array.from(new Set(query.toLowerCase().split(/\s+/)
        .map((t) => t.replace(edgePunct, ""))
        .filter((t) => t !== "")));
    if (terms.length === 0) {
        return null;
    }
    terms.sort((a, b) => b.length - a.length);
    return new RegExp(terms.map((t) => t.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")).join("|"), "giu");
}

// appendHighlighted appends text to el with matches wrapped in <b>, using text
// nodes only.
function appendHighlighted(el: HTMLElement, text: string, match: RegExp | null) {
    if (!match) {
        el.textContent = text;
        return;
    }
    let last = 0;
    match.lastIndex = 0;
    for (let m = match.exec(text); m !== null; m = match.exec(text)) {
        const start = m.index;
        if (start > last) {
            el.append(text.slice(last, start));
        }
        const b = document.createElement("b");
        b.textContent = m[0];
        el.append(b);
        last = start + m[0].length;
    }
    if (last < text.length) {
        el.append(text.slice(last));
    }
}

function hostOf(url: string): string {
    try {
        return new URL(url).hostname;
    } catch {
        return "";
    }
}

function monogram(domain: string): string {
    const letter = domain.replace(/^www\./, "").match(letterOrDigit)?.[0];
    return letter ? letter.toUpperCase() : "?";
}

// hash is FNV-1a over the domain, so a domain always gets the same tint.
function hash(s: string): number {
    let h = 0x811c9dc5;
    for (let i = 0; i < s.length; i++) {
        h ^= s.charCodeAt(i);
        h = Math.imul(h, 0x01000193);
    }
    return h >>> 0;
}

// The selection is drawn in a neutral colour while the window is inactive.
window.addEventListener("blur", () => {
    document.documentElement.dataset.inactive = "";
});
window.addEventListener("focus", () => {
    delete document.documentElement.dataset.inactive;
});

import {Events, Window} from "@wailsio/runtime";
import {BookmarkService, WindowService} from "../bindings/github.com/ndy40/cairn/desktop";
import type {Bookmark} from "../bindings/github.com/ndy40/cairn/desktop";

const queryEl = document.getElementById("query") as HTMLInputElement;
const resultsEl = document.getElementById("results") as HTMLUListElement;
const statusEl = document.getElementById("status") as HTMLSpanElement;
const searchView = document.getElementById("search-view") as HTMLElement;
const addView = document.getElementById("add-view") as HTMLFormElement;
const addURLEl = document.getElementById("add-url") as HTMLInputElement;
const addTagsEl = document.getElementById("add-tags") as HTMLInputElement;
const searchHints = document.getElementById("hints-search") as HTMLElement;
const addHints = document.getElementById("hints-add") as HTMLElement;

const searchDebounceMs = 80;

type Mode = "search" | "add";

let mode: Mode = "search";
let results: Bookmark[] = [];
let selected = 0;
let searchTimer: ReturnType<typeof setTimeout> | undefined;
// Guards against a slow, older search overwriting a newer one.
let searchSeq = 0;
let saving = false;

function setStatus(message: string) {
    statusEl.textContent = message;
    // Long messages are truncated in the footer; show the full text on hover.
    statusEl.title = message;
}

// errorMessage extracts the Go error text from a rejected binding call,
// capitalised for display since Go error strings are lowercase by convention.
function errorMessage(err: unknown): string {
    const text = err instanceof Error ? err.message : String(err);
    return text.charAt(0).toUpperCase() + text.slice(1);
}

// render rebuilds the result list. Text is set via textContent because titles
// come from fetched web pages and must never be interpreted as HTML.
function render() {
    resultsEl.replaceChildren();

    if (results.length === 0) {
        const empty = document.createElement("li");
        empty.className = "empty";
        empty.textContent = queryEl.value ? "No matching bookmarks" : "No bookmarks yet";
        resultsEl.append(empty);
        return;
    }

    results.forEach((b, i) => {
        const li = document.createElement("li");
        li.className = "result" + (i === selected ? " is-selected" : "");
        li.setAttribute("role", "option");
        li.setAttribute("aria-selected", String(i === selected));

        const title = document.createElement("div");
        title.className = "title";
        title.textContent = (b.pinned ? "📌 " : "") + (b.title || b.url);

        const meta = document.createElement("div");
        meta.className = "meta";
        const domain = document.createElement("span");
        domain.className = "domain";
        domain.textContent = b.domain;
        meta.append(domain);
        for (const tag of b.tags ?? []) {
            const t = document.createElement("span");
            t.className = "tag";
            t.textContent = tag;
            meta.append(t);
        }

        li.append(title, meta);
        li.addEventListener("mousemove", () => select(i));
        li.addEventListener("click", () => openSelected());
        resultsEl.append(li);
    });
}

function select(i: number) {
    if (i === selected || i < 0 || i >= results.length) {
        return;
    }
    selected = i;
    render();
    resultsEl.children[selected]?.scrollIntoView({block: "nearest"});
}

async function runSearch() {
    const seq = ++searchSeq;
    try {
        const found = await BookmarkService.Search(queryEl.value.trim());
        if (seq !== searchSeq) {
            return;
        }
        results = found ?? [];
        selected = 0;
        render();
    } catch (err) {
        setStatus(`Search failed: ${errorMessage(err)}`);
    }
}

async function openSelected() {
    const b = results[selected];
    if (!b) {
        return;
    }
    try {
        await BookmarkService.Open(b.id);
        await Window.Hide();
    } catch (err) {
        setStatus(`Could not open: ${errorMessage(err)}`);
    }
}

function setMode(next: Mode) {
    mode = next;
    searchView.hidden = next !== "search";
    addView.hidden = next !== "add";
    searchHints.hidden = next !== "search";
    addHints.hidden = next !== "add";
}

// showAdd switches to the add form. A URL typed after "+" in the search box
// wins; otherwise the field is pre-filled from the clipboard when it holds a URL.
async function showAdd(typedURL = "") {
    setMode("add");
    setStatus("");
    addTagsEl.value = "";
    addURLEl.value = typedURL;
    addURLEl.focus();

    if (!typedURL) {
        try {
            const clip = await BookmarkService.ClipboardURL();
            // Don't clobber anything typed while the clipboard was being read.
            if (mode === "add" && !addURLEl.value) {
                addURLEl.value = clip;
                addURLEl.select();
            }
        } catch {
            // Clipboard access is best-effort.
        }
    }
}

function showSearch() {
    setMode("search");
    queryEl.value = "";
    queryEl.focus();
    runSearch();
}

async function saveBookmark() {
    if (saving) {
        return;
    }
    saving = true;
    addURLEl.disabled = addTagsEl.disabled = true;
    setStatus("Saving…");
    try {
        const res = await BookmarkService.Add(addURLEl.value, addTagsEl.value);
        showSearch();
        const title = res?.bookmark.title ?? "";
        setStatus(res?.titleFetched ? `Saved: ${title}` : `Saved (title unavailable): ${title}`);
    } catch (err) {
        setStatus(errorMessage(err));
        addURLEl.focus();
    } finally {
        saving = false;
        addURLEl.disabled = addTagsEl.disabled = false;
    }
}

queryEl.addEventListener("input", () => {
    if (queryEl.value.startsWith("+")) {
        showAdd(queryEl.value.slice(1).trim());
        return;
    }
    clearTimeout(searchTimer);
    searchTimer = setTimeout(runSearch, searchDebounceMs);
});

addView.addEventListener("submit", (e) => {
    e.preventDefault();
    saveBookmark();
});

document.addEventListener("keydown", (e) => {
    if (e.key === "n" && (e.ctrlKey || e.metaKey) && mode === "search") {
        e.preventDefault();
        showAdd();
        return;
    }

    if (mode === "add") {
        if (e.key === "Escape" && !saving) {
            e.preventDefault();
            showSearch();
            setStatus("");
        } else if (e.key === "Enter") {
            e.preventDefault();
            saveBookmark();
        }
        return;
    }

    switch (e.key) {
        case "ArrowDown":
            e.preventDefault();
            select(selected + 1);
            break;
        case "ArrowUp":
            e.preventDefault();
            select(selected - 1);
            break;
        case "Enter":
            e.preventDefault();
            openSelected();
            break;
        case "Escape":
            e.preventDefault();
            // First Esc clears the search; Esc on an empty box dismisses the popup.
            if (queryEl.value) {
                queryEl.value = "";
                setStatus("");
                runSearch();
            } else {
                Window.Hide();
            }
            break;
    }
});

// The popup is reused rather than reloaded, so each time Go shows it, reset to
// the requested view and re-query so bookmarks changed elsewhere appear.
Events.On("launcher:show", (event) => {
    setStatus("");
    if (event.data === "add") {
        showAdd();
    } else {
        showSearch();
    }
});

// A background sync finished: refresh results without disturbing the query.
Events.On("bookmarks:changed", () => {
    if (mode === "search") {
        runSearch();
    }
});

Events.On("sync:warning", (event) => {
    setStatus(event.data);
});

// Tell Go about every press inside the popup. On GNOME/Wayland a press (click,
// drag or edge resize) can make the window briefly lose focus, and Go uses this
// to avoid hiding the popup in that case. Capture phase so drags still report.
window.addEventListener("pointerdown", () => {
    WindowService.PointerDown();
}, {capture: true});

queryEl.focus();
runSearch();

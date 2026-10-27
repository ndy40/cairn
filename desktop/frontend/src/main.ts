import {Events} from "@wailsio/runtime";
import {BookmarkService, WindowService} from "../bindings/github.com/ndy40/cairn/desktop";
import type {Bookmark} from "../bindings/github.com/ndy40/cairn/desktop";
import {renderResults} from "./render";
import {clearBanner, setStatus, showBanner} from "./messages";

const queryEl = document.getElementById("query") as HTMLInputElement;
const searchField = document.getElementById("search-field") as HTMLLabelElement;
const addButton = document.getElementById("add-btn") as HTMLButtonElement;
const resultsEl = document.getElementById("results") as HTMLUListElement;
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
// The query the current results are for; a refresh of the same query keeps the
// selected bookmark, a new query selects the top result.
let resultsQuery = "";
let searchTimer: ReturnType<typeof setTimeout> | undefined;
// Guards against a slow, older search overwriting a newer one.
let searchSeq = 0;
let saving = false;

// errorMessage extracts the Go error text from a rejected binding call,
// capitalised for display since Go error strings are lowercase by convention.
function errorMessage(err: unknown): string {
    const text = err instanceof Error ? err.message : String(err);
    return text.charAt(0).toUpperCase() + text.slice(1);
}

function render() {
    renderResults(resultsEl, results, selected, resultsQuery, {hover: select, open: () => openSelected()});
}

function select(i: number) {
    if (i === selected || i < 0 || i >= results.length) {
        return;
    }
    selected = i;
    render();
    resultsEl.querySelector(".is-selected")?.scrollIntoView({block: "nearest"});
}

async function runSearch() {
    const seq = ++searchSeq;
    const query = queryEl.value.trim();
    try {
        const found = await BookmarkService.Search(query);
        if (seq !== searchSeq) {
            return;
        }
        const keep = query === resultsQuery ? results[selected]?.id : undefined;
        results = found ?? [];
        resultsQuery = query;
        selected = Math.max(0, results.findIndex((b) => b.id === keep));
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
        // Go decides whether to put Cairn away ("Hide Cairn after opening a link").
        await WindowService.AfterOpen();
    } catch (err) {
        setStatus(`Could not open: ${errorMessage(err)}`);
    }
}

function setMode(next: Mode) {
    mode = next;
    searchView.hidden = next !== "search";
    searchField.hidden = addButton.hidden = next !== "search";
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

// focusSearch is used each time the window comes forward: the query is kept
// and selected so typing replaces it, and results are refreshed so bookmarks
// changed elsewhere appear.
function focusSearch() {
    setMode("search");
    queryEl.focus();
    queryEl.select();
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

addButton.addEventListener("click", () => showAdd());

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
            // Esc clears the search and never closes or hides the window.
            if (queryEl.value) {
                queryEl.value = "";
                setStatus("");
                runSearch();
            }
            break;
    }
});

// Every way of bringing the window forward (launch, shortcut, tray, --toggle)
// sends this with the view to open.
Events.On("launcher:show", (event) => {
    setStatus("");
    if (event.data === "add") {
        showAdd();
    } else {
        focusSearch();
    }
});

// A background sync pull finished: refresh results without disturbing the query.
Events.On("bookmarks:changed", () => {
    if (mode === "search") {
        runSearch();
    }
});

Events.On("sync:warning", (event) => {
    showBanner(event.data);
});

Events.On("sync:ok", () => {
    clearBanner();
});

// Background progress ("Finishing sync…"); "" clears it.
Events.On("status", (event) => {
    setStatus(event.data);
});

// applyPlatform tags <html> with data-os and data-desktop so the CSS can draw
// the macOS title band or the Linux search bar and the desktop variants.
async function applyPlatform() {
    try {
        const p = await WindowService.Platform();
        document.documentElement.dataset.os = p.os;
        if (p.desktop) {
            document.documentElement.dataset.desktop = p.desktop;
        }
    } catch {
        // Without platform info the neutral styles still work.
    }
}

applyPlatform();
queryEl.focus();
runSearch();

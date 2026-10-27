import {Events} from "@wailsio/runtime";
import {BookmarkService, WindowService} from "../bindings/github.com/ndy40/cairn/desktop";
import type {Bookmark} from "../bindings/github.com/ndy40/cairn/desktop";
import {initAdd, isAddOpen, openAdd} from "./add";
import {renderResults} from "./render";
import {clearBanner, setStatus, showBanner} from "./messages";

const queryEl = document.getElementById("query") as HTMLInputElement;
const addButton = document.getElementById("add-btn") as HTMLButtonElement;
const resultsEl = document.getElementById("results") as HTMLUListElement;

const searchDebounceMs = 80;

let results: Bookmark[] = [];
let selected = 0;
// The query the current results are for; a refresh of the same query keeps the
// selected bookmark, a new query selects the top result.
let resultsQuery = "";
let searchTimer: ReturnType<typeof setTimeout> | undefined;
// Guards against a slow, older search overwriting a newer one.
let searchSeq = 0;

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

// reveal shows b in an unfiltered list and selects it. A bookmark outside the
// first page of results (e.g. behind many pinned ones) is put at the top.
async function reveal(b: Bookmark) {
    queryEl.value = "";
    await runSearch();
    let i = results.findIndex((r) => r.id === b.id);
    if (i < 0) {
        results.unshift(b);
        i = 0;
    }
    selected = i;
    render();
    resultsEl.children[selected]?.scrollIntoView({block: "nearest"});
}

// focusSearch is used each time the window comes forward: the query is kept
// and selected so typing replaces it, and results are refreshed so bookmarks
// changed elsewhere appear.
function focusSearch() {
    if (isAddOpen()) {
        return;
    }
    queryEl.focus();
    queryEl.select();
    runSearch();
}

queryEl.addEventListener("input", () => {
    if (queryEl.value.startsWith("+")) {
        const typed = queryEl.value.slice(1).trim();
        queryEl.value = "";
        runSearch();
        openAdd(typed);
        return;
    }
    clearTimeout(searchTimer);
    searchTimer = setTimeout(runSearch, searchDebounceMs);
});

addButton.addEventListener("click", () => openAdd());

document.addEventListener("keydown", (e) => {
    // The open <dialog> handles its own keys: Enter submits, Esc cancels.
    if (isAddOpen()) {
        return;
    }

    if (e.key === "n" && (e.ctrlKey || e.metaKey)) {
        e.preventDefault();
        openAdd();
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
        openAdd();
    } else {
        focusSearch();
    }
});

// A background sync pull finished: refresh results without disturbing the query.
Events.On("bookmarks:changed", () => {
    runSearch();
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

// Background progress from Go, e.g. "Finishing sync…" while quitting; "" clears it.
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

initAdd({
    reveal,
    setStatus,
    errorMessage,
    closed: () => queryEl.focus(),
});
applyPlatform();
queryEl.focus();
runSearch();

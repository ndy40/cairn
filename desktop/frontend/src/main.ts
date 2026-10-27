import {Events} from "@wailsio/runtime";
import {BookmarkService, WindowService} from "../bindings/github.com/ndy40/cairn/desktop";
import type {Bookmark, PlatformInfo} from "../bindings/github.com/ndy40/cairn/desktop";
import {initAdd, openAdd} from "./add";
import {renderResults} from "./render";
import {clearBanner, setStatus, showBanner, showToast} from "./messages";
import {initMenus, isMenuOpen, openRowMenu} from "./menus";
import type {MenuAction} from "./menus";
import {keyAction} from "./keys";
import type {KeyAction} from "./keys";
import {initSettings, openSettings, platformChanged} from "./settings";
import {initPanels, openAbout, openShortcuts} from "./panels";
import {applyAccent} from "./accent";

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

// pageSize is how many rows fit in the visible list, for PageUp/PageDown.
function pageSize(): number {
    const row = resultsEl.querySelector(".result") as HTMLElement | null;
    return row ? Math.max(1, Math.floor(resultsEl.clientHeight / row.offsetHeight) - 1) : 1;
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

async function copySelected() {
    const b = results[selected];
    if (!b) {
        return;
    }
    try {
        await BookmarkService.CopyLink(b.id);
        showToast("Link copied");
    } catch (err) {
        setStatus(`Could not copy: ${errorMessage(err)}`);
    }
}

// togglePin pins or unpins the selected bookmark; Go then sends
// bookmarks:changed "pin", which refreshes the list keeping the selection.
async function togglePin() {
    const b = results[selected];
    if (!b) {
        return;
    }
    try {
        await BookmarkService.SetPinned(b.id, !b.pinned);
    } catch (err) {
        setStatus(`Could not pin: ${errorMessage(err)}`);
    }
}

// syncNow pulls and pushes now and reports the result in the status area.
async function syncNow() {
    setStatus("Syncing…");
    try {
        await BookmarkService.SyncNow();
        setStatus("Synced");
    } catch (err) {
        setStatus(errorMessage(err));
    }
}

type Action = KeyAction | MenuAction;

// run performs an action from the keyboard map, a menu or the macOS menu bar.
function run(action: Action) {
    switch (action) {
        case "find":
            queryEl.focus();
            queryEl.select();
            break;
        case "add":
            openAdd();
            break;
        case "open":
            openSelected();
            break;
        case "copy":
            copySelected();
            break;
        case "pin":
            togglePin();
            break;
        case "sync":
            syncNow();
            break;
        case "settings":
            openSettings();
            break;
        case "shortcuts":
            openShortcuts();
            break;
        case "about":
            openAbout();
            break;
        case "close":
            WindowService.CloseWindow();
            break;
        case "quit":
            WindowService.Quit();
            break;
    }
}

// modalOpen reports whether a dialog or popup menu has the keyboard.
function modalOpen(): boolean {
    return document.querySelector("dialog[open]") !== null || isMenuOpen();
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
    resultsEl.querySelector(".is-selected")?.scrollIntoView({block: "nearest"});
}

// focusSearch is used each time the window comes forward: the query is kept
// and selected so typing replaces it, and results are refreshed so bookmarks
// changed elsewhere appear.
function focusSearch() {
    if (document.querySelector("dialog[open]")) {
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
    const modal = modalOpen();
    const act = keyAction(e, modal);
    if (act) {
        e.preventDefault();
        run(act);
        return;
    }
    // An open <dialog> or menu handles its own keys: Enter submits, Esc cancels.
    if (modal) {
        return;
    }

    if ((e.key === "F10" && e.shiftKey) || e.key === "ContextMenu") {
        e.preventDefault();
        openRowMenu();
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
        case "PageDown":
            e.preventDefault();
            select(Math.min(selected + pageSize(), results.length - 1));
            break;
        case "PageUp":
            e.preventDefault();
            select(Math.max(selected - pageSize(), 0));
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

// Panels opened from a menu, the tray or the Dock, after the window came forward.
Events.On("ui:open", (event) => {
    run(event.data as Action);
});

// macOS menu bar items that act on the selected row.
Events.On("menu:action", (event) => {
    run(event.data as Action);
});

// A sync pull or a pin changed bookmarks: refresh results without disturbing
// the query, keeping the selected bookmark.
Events.On("bookmarks:changed", () => {
    runSearch();
});

Events.On("sync:warning", (event) => {
    showBanner(event.data);
});

Events.On("sync:ok", () => {
    clearBanner();
});

// Background progress from Go, e.g. "Finishing sync…" while quitting; "" clears it.
Events.On("status", (event) => {
    setStatus(event.data);
});

// applyPlatform tags <html> with data-os and data-desktop so the CSS can draw
// the macOS title band or the Linux search bar and the desktop variants.
async function applyPlatform() {
    try {
        showPlatform(await WindowService.Platform());
    } catch {
        // Without platform info the neutral styles still work.
    }
}

function showPlatform(p: PlatformInfo) {
    document.documentElement.dataset.os = p.os;
    if (p.desktop) {
        document.documentElement.dataset.desktop = p.desktop;
    }
    applyAccent(p.accentColor);
}

// A tray host came or went, or the accent colour changed.
Events.On("platform:changed", (event) => {
    showPlatform(event.data);
    platformChanged(event.data);
});

initAdd({
    reveal,
    setStatus,
    errorMessage,
    closed: () => queryEl.focus(),
});
initMenus({
    run,
    rowAt: (el) => {
        const id = el.closest("li.result")?.id.replace("bookmark-", "");
        return id ? results.findIndex((b) => String(b.id) === id) : -1;
    },
    select,
    selected: () => {
        const bookmark = results[selected];
        const row = resultsEl.querySelector<HTMLElement>(".is-selected");
        return bookmark && row ? {bookmark, row} : undefined;
    },
});
initSettings({errorMessage, closed: () => queryEl.focus()});
initPanels(() => queryEl.focus());
applyPlatform();
queryEl.focus();
runSearch();

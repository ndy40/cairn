import {BookmarkService} from "../bindings/github.com/ndy40/cairn/desktop";
import type {Bookmark} from "../bindings/github.com/ndy40/cairn/desktop";

// The Add Bookmark modal. main.ts owns the list and the footer status; this
// module owns the <dialog> and reports back through the hooks given to initAdd.

const dialog = document.getElementById("add-dialog") as HTMLDialogElement;
const form = document.getElementById("add-form") as HTMLFormElement;
const urlEl = document.getElementById("add-url") as HTMLInputElement;
const tagsEl = document.getElementById("add-tags") as HTMLInputElement;
const messageEl = document.getElementById("add-message") as HTMLParagraphElement;
const cancelButton = document.getElementById("add-cancel") as HTMLButtonElement;
const submitButton = document.getElementById("add-submit") as HTMLButtonElement;

// Text of ErrAlreadyBookmarked in service.go, which reaches JS as the error message.
const alreadyBookmarked = "already bookmarked";

export interface AddHooks {
    // reveal clears the search, refreshes the list and selects b.
    reveal(b: Bookmark): Promise<void>;
    setStatus(message: string): void;
    errorMessage(err: unknown): string;
    // closed runs after the modal closes without a save, to refocus search.
    closed(): void;
}

let hooks: AddHooks;
let saving = false;
// Bumped on every open so a slow clipboard read can't fill a later session.
let session = 0;

export function isAddOpen(): boolean {
    return dialog.open;
}

// openAdd shows the modal. Text typed after "+" in search is used as-is with
// the cursor at the end; otherwise the URL is pre-filled from the clipboard
// and selected so typing replaces it. Opening it again keeps what was typed.
export async function openAdd(prefill = "") {
    if (dialog.open) {
        if (prefill && !saving) {
            setURL(prefill);
        }
        urlEl.focus();
        return;
    }

    const current = ++session;
    tagsEl.value = "";
    setMessage("");
    setURL(prefill);
    dialog.showModal();
    urlEl.focus();

    if (prefill) {
        return;
    }
    try {
        const clip = await BookmarkService.ClipboardURL();
        // Don't clobber anything typed while the clipboard was being read.
        if (current === session && dialog.open && !urlEl.value && clip) {
            urlEl.value = clip;
            urlEl.select();
        }
    } catch {
        // Clipboard access is best-effort.
    }
}

function setURL(text: string) {
    urlEl.value = text;
    urlEl.setSelectionRange(text.length, text.length);
}

function setMessage(text: string, isError = false) {
    messageEl.textContent = text;
    messageEl.classList.toggle("is-error", isError);
}

function close() {
    if (saving) {
        return;
    }
    dialog.close();
}

function setSaving(on: boolean) {
    saving = on;
    urlEl.disabled = tagsEl.disabled = submitButton.disabled = cancelButton.disabled = on;
}

async function save() {
    if (saving) {
        return;
    }
    setSaving(true);
    // Add fetches the page before it returns.
    setMessage("Fetching title…");
    try {
        const res = await BookmarkService.Add(urlEl.value, tagsEl.value);
        setSaving(false);
        dialog.close();
        if (res) {
            await hooks.reveal(res.bookmark);
            hooks.setStatus(res.titleFetched ? `Added: ${res.bookmark.title}` : "Added (title unavailable)");
        }
    } catch (err) {
        setSaving(false);
        if (err instanceof Error && err.message === alreadyBookmarked) {
            await showExisting();
            return;
        }
        setMessage(hooks.errorMessage(err), true);
        urlEl.focus();
    }
}

async function showExisting() {
    const url = urlEl.value;
    dialog.close();
    try {
        const existing = await BookmarkService.FindByURL(url);
        if (existing) {
            await hooks.reveal(existing);
        }
    } catch {
        // The duplicate message below is still accurate.
    }
    // TODO(toast): show as toast
    hooks.setStatus("Already bookmarked");
}

export function initAdd(h: AddHooks) {
    hooks = h;

    // method="dialog" would close on submit; close only after a successful save.
    form.addEventListener("submit", (e) => {
        e.preventDefault();
        save();
    });

    cancelButton.addEventListener("click", close);

    // Esc fires "cancel"; keep the modal while a save is in flight.
    dialog.addEventListener("cancel", (e) => {
        if (saving) {
            e.preventDefault();
        }
    });

    dialog.addEventListener("close", () => {
        setMessage("");
        hooks.closed();
    });
}

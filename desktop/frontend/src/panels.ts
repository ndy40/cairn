// The Keyboard Shortcuts and About panels.

import {Browser} from "@wailsio/runtime";
import {WindowService} from "../bindings/github.com/ndy40/cairn/desktop";
import {shortcutList} from "./keys";

// The Desktop App user guide (same as Help › Cairn Help on macOS).
const helpURL = "https://ndy40.github.io/cairn/docs/desktop-app/";

const shortcutsDialog = document.getElementById("shortcuts-dialog") as HTMLDialogElement;
const shortcutListEl = document.getElementById("shortcut-list") as HTMLDListElement;
const aboutDialog = document.getElementById("about-dialog") as HTMLDialogElement;
const aboutVersion = document.getElementById("about-version") as HTMLParagraphElement;
const aboutHelp = document.getElementById("about-help") as HTMLButtonElement;

// openShortcuts lists the keyboard map with the platform's key names.
export function openShortcuts() {
    shortcutListEl.replaceChildren();
    for (const s of shortcutList()) {
        const dt = document.createElement("dt");
        const key = document.createElement("kbd");
        key.textContent = s.keys;
        dt.append(key);
        const dd = document.createElement("dd");
        dd.textContent = s.does;
        shortcutListEl.append(dt, dd);
    }
    show(shortcutsDialog);
}

export async function openAbout() {
    show(aboutDialog);
    try {
        aboutVersion.textContent = `Version ${await WindowService.Version()}`;
    } catch {
        aboutVersion.textContent = "";
    }
}

function show(dialog: HTMLDialogElement) {
    if (!dialog.open) {
        dialog.showModal();
    }
    dialog.querySelector<HTMLButtonElement>("[type=submit]")?.focus();
}

export function initPanels(closed: () => void) {
    aboutHelp.addEventListener("click", () => {
        Browser.OpenURL(helpURL).catch(() => undefined);
    });
    shortcutsDialog.addEventListener("close", closed);
    aboutDialog.addEventListener("close", closed);
}

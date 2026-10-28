// The Settings panel (Preferences on Linux). Each switch applies at once
// through WindowService; nothing needs a restart.

import {Clipboard} from "@wailsio/runtime";
import {WindowService} from "../bindings/github.com/ndy40/cairn/desktop";
import type {PlatformInfo} from "../bindings/github.com/ndy40/cairn/desktop";
import {hotkeyLabel, isMac} from "./keys";

const dialog = document.getElementById("settings-dialog") as HTMLDialogElement;
const title = document.getElementById("settings-title") as HTMLHeadingElement;
const loginSwitch = document.getElementById("login-switch") as HTMLButtonElement;
const menubarRow = document.getElementById("menubar-row") as HTMLElement;
const menubarSwitch = document.getElementById("menubar-switch") as HTMLButtonElement;
const backgroundRow = document.getElementById("background-row") as HTMLElement;
const backgroundSwitch = document.getElementById("background-switch") as HTMLButtonElement;
const backgroundHelp = document.getElementById("background-help") as HTMLSpanElement;
const hideSwitch = document.getElementById("hide-switch") as HTMLButtonElement;
const shortcutKey = document.getElementById("shortcut-key") as HTMLElement;
const shortcutHelp = document.getElementById("shortcut-help") as HTMLSpanElement;
const fallback = document.getElementById("shortcut-fallback") as HTMLElement;
const shortcutError = document.getElementById("shortcut-error") as HTMLParagraphElement;
const toggleCommand = document.getElementById("toggle-command") as HTMLElement;
const copyCommand = document.getElementById("copy-command") as HTMLButtonElement;
const message = document.getElementById("settings-message") as HTMLParagraphElement;
const note = document.getElementById("settings-note") as HTMLParagraphElement;

// Setting keys accepted by WindowService.SetSetting.
type SettingKey = "showInMenuBar" | "runInBackground" | "hideAfterOpen";

// How long the Copy Command button says "Copied".
const copiedMs = 2000;

let errorMessage: (err: unknown) => string = String;

function setSwitch(el: HTMLButtonElement, on: boolean) {
    el.setAttribute("aria-checked", String(on));
}

function isOn(el: HTMLButtonElement): boolean {
    return el.getAttribute("aria-checked") === "true";
}

// openSettings shows the panel with the current values read from Go and the
// OS (the login item is read live, so one made by an older installer shows).
export async function openSettings() {
    title.textContent = isMac() ? "Settings" : "Preferences";
    message.textContent = "";
    if (!dialog.open) {
        dialog.showModal();
    }
    loginSwitch.focus();
    await refresh();
}

export function isSettingsOpen(): boolean {
    return dialog.open;
}

async function refresh() {
    const [settings, shortcut, platform] = await Promise.all([
        WindowService.Settings(),
        WindowService.ShortcutStatus(),
        WindowService.Platform(),
    ]);
    setSwitch(hideSwitch, settings.hideAfterOpen);
    setSwitch(menubarSwitch, settings.showInMenuBar);
    setSwitch(backgroundSwitch, settings.runInBackground ?? false);
    showPlatform(platform);

    try {
        setSwitch(loginSwitch, await WindowService.OpenAtLogin());
        loginSwitch.disabled = false;
    } catch (err) {
        loginSwitch.disabled = true;
        message.textContent = `Could not read the login item: ${errorMessage(err)}`;
    }

    if (shortcut.hotkey === "off") {
        shortcutKey.textContent = "Off";
        shortcutHelp.textContent = "Turned off with desktop_hotkey in cairn.json.";
    } else {
        shortcutKey.textContent = hotkeyLabel(shortcut.hotkey);
        shortcutHelp.textContent = shortcut.registered
            ? "Brings Cairn forward from any app. Set desktop_hotkey in cairn.json to change it."
            : "Not available.";
    }
    const failed = shortcut.hotkey !== "off" && !shortcut.registered;
    fallback.hidden = !failed;
    shortcutError.textContent = failed ? `The shortcut could not be registered: ${shortcut.error || "unknown error"}.` : "";
    toggleCommand.textContent = shortcut.toggleCommand;
}

// showPlatform shows the rows for this platform. Run in background needs a
// tray to return to, so it is disabled with an explanation without one.
function showPlatform(p: PlatformInfo) {
    const mac = p.os === "darwin";
    menubarRow.hidden = !mac;
    backgroundRow.hidden = mac;
    backgroundSwitch.disabled = !p.trayAvailable;
    if (p.trayAvailable) {
        backgroundHelp.textContent = "Keep Cairn in the system tray when its window is closed.";
    } else {
        backgroundHelp.textContent = "Requires a system tray. On GNOME, install the AppIndicator extension.";
        setSwitch(backgroundSwitch, false);
    }
    note.textContent = mac
        ? "Closing the window keeps Cairn running. Quit with ⌘Q."
        : "Quit Cairn with Ctrl+Q or from the main menu.";
}

// platformChanged updates an open panel when a tray host comes or goes.
export function platformChanged(p: PlatformInfo) {
    if (dialog.open) {
        refresh().catch(() => showPlatform(p));
    }
}

// toggle flips a switch at once and applies it, reverting on failure.
async function toggle(el: HTMLButtonElement, apply: (on: boolean) => Promise<unknown>) {
    if (el.disabled) {
        return;
    }
    const on = !isOn(el);
    setSwitch(el, on);
    message.textContent = "";
    try {
        await apply(on);
    } catch (err) {
        setSwitch(el, !on);
        message.textContent = errorMessage(err);
    }
}

function setting(key: SettingKey) {
    return (on: boolean) => WindowService.SetSetting(key, on);
}

export function initSettings(opts: { errorMessage(err: unknown): string; closed(): void }) {
    errorMessage = opts.errorMessage;
    loginSwitch.addEventListener("click", () => toggle(loginSwitch, (on) => WindowService.SetOpenAtLogin(on)));
    menubarSwitch.addEventListener("click", () => toggle(menubarSwitch, setting("showInMenuBar")));
    backgroundSwitch.addEventListener("click", () => toggle(backgroundSwitch, setting("runInBackground")));
    hideSwitch.addEventListener("click", () => toggle(hideSwitch, setting("hideAfterOpen")));
    copyCommand.addEventListener("click", async () => {
        try {
            await Clipboard.SetText(toggleCommand.textContent ?? "");
            // A toast would sit under the modal, so confirm on the button.
            copyCommand.textContent = "Copied";
            setTimeout(() => copyCommand.textContent = "Copy Command", copiedMs);
        } catch (err) {
            message.textContent = `Could not copy: ${errorMessage(err)}`;
        }
    });
    dialog.addEventListener("close", opts.closed);
}

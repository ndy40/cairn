// The app-wide keyboard map (FR-028). On macOS the menu bar owns these
// accelerators, so the web view leaves them alone and each action fires once;
// on Linux, which has no menu bar, they are handled here.

export type KeyAction = "find" | "copy" | "add" | "close" | "quit" | "settings" | "shortcuts";

export function isMac(): boolean {
    return document.documentElement.dataset.os === "darwin";
}

// keyAction maps a keydown to an app action, or undefined. modalOpen limits
// it to Close and Quit while a dialog or menu has the keyboard.
export function keyAction(e: KeyboardEvent, modalOpen: boolean): KeyAction | undefined {
    if (isMac() || !e.ctrlKey || e.altKey || e.metaKey) {
        return undefined;
    }
    const key = e.key.toLowerCase();
    let act: KeyAction | undefined;
    if (e.shiftKey) {
        // Ctrl+? is Ctrl+Shift+/ on most layouts.
        if (key === "c") {
            act = "copy";
        } else if (key === "?" || key === "/") {
            act = "shortcuts";
        }
    } else {
        act = ({f: "find", n: "add", w: "close", q: "quit", ",": "settings", "/": "shortcuts"} as Record<string, KeyAction>)[key];
    }
    if (modalOpen && act !== "close" && act !== "quit") {
        return undefined;
    }
    return act;
}

export interface Shortcut {
    keys: string;
    does: string;
}

// shortcutList is shown in the Keyboard Shortcuts panel, with the platform's
// key names.
export function shortcutList(): Shortcut[] {
    const mac = isMac();
    const mod = mac ? "⌘" : "Ctrl+";
    const shiftMod = mac ? "⇧⌘" : "Shift+Ctrl+";
    return [
        {keys: `${mod}F`, does: "Search"},
        {keys: "↑ ↓", does: "Move the selection"},
        {keys: mac ? "↩" : "Enter", does: "Open the link"},
        {keys: `${shiftMod}C`, does: "Copy the link"},
        {keys: `${mod}N`, does: "Add a bookmark"},
        {keys: mac ? "esc" : "Esc", does: "Clear the search, or cancel"},
        {keys: mac ? "⇧F10" : "Shift+F10", does: "Bookmark actions"},
        {keys: `${mod}W`, does: "Close the window"},
        {keys: `${mod}Q`, does: "Quit Cairn"},
        {keys: `${mod},`, does: mac ? "Settings" : "Preferences"},
        {keys: mac ? "⌘/" : "Ctrl+?", does: "Keyboard shortcuts"},
    ];
}

// hotkeyLabel formats a Wails accelerator such as "CmdOrCtrl+Shift+Space"
// with the platform's key names.
export function hotkeyLabel(accelerator: string): string {
    const mac = isMac();
    const names: Record<string, string> = mac
        ? {cmdorctrl: "⌘", cmd: "⌘", command: "⌘", super: "⌘", ctrl: "⌃", control: "⌃", shift: "⇧", alt: "⌥", option: "⌥", optionoralt: "⌥", space: "Space"}
        : {cmdorctrl: "Ctrl", ctrl: "Ctrl", control: "Ctrl", cmd: "Super", super: "Super", shift: "Shift", alt: "Alt", optionoralt: "Alt", space: "Space"};
    const parts = accelerator.split("+").map((p) => names[p.trim().toLowerCase()] ?? p.trim().toUpperCase());
    return parts.join(mac ? "" : "+");
}

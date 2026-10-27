// HTML popup menus: the Linux main menu (☰ in the search bar) and the row
// context menu. Both are role="menu" with arrow-key navigation; Esc or a
// click elsewhere closes them and returns focus to where it was.

import type {Bookmark} from "../bindings/github.com/ndy40/cairn/desktop";

export type MenuAction = "add" | "sync" | "settings" | "shortcuts" | "about" | "quit" | "open" | "copy" | "pin";

export interface MenuHooks {
    run(action: MenuAction): void;
    // rowAt returns the index of the result row containing el, or -1.
    rowAt(el: Element): number;
    // select makes result i the selected row.
    select(i: number): void;
    // selected returns the selected bookmark and its row element.
    selected(): { bookmark: Bookmark; row: HTMLElement } | undefined;
}

// Gap between a menu and the window edge.
const edge = 6;

class PopupMenu {
    private returnFocus: HTMLElement | null = null;

    // opener is the button that toggles the menu, whose own clicks are not
    // "clicks elsewhere".
    constructor(private readonly el: HTMLElement, private readonly run: (action: MenuAction) => void,
                private readonly opener?: HTMLElement) {
        el.addEventListener("click", (e) => {
            const item = (e.target as Element).closest<HTMLButtonElement>("[role=menuitem]");
            if (item && !item.disabled) {
                this.close();
                this.run(item.dataset.action as MenuAction);
            }
        });
        el.addEventListener("keydown", (e) => this.onKey(e));
        // Clicks and focus moving elsewhere close the menu.
        document.addEventListener("pointerdown", (e) => {
            const target = e.target as Node;
            if (this.isOpen() && !el.contains(target) && !this.opener?.contains(target)) {
                this.close(false);
            }
        }, true);
        el.addEventListener("focusout", (e) => {
            if (this.isOpen() && !el.contains(e.relatedTarget as Node | null)) {
                this.close(false);
            }
        });
        window.addEventListener("blur", () => this.close(false));
        window.addEventListener("resize", () => this.close(false));
    }

    isOpen(): boolean {
        return !this.el.hidden;
    }

    // openAt shows the menu with its top-left corner at x, y, kept inside the
    // window, and focuses the first item.
    openAt(x: number, y: number) {
        this.returnFocus = document.activeElement as HTMLElement | null;
        this.el.hidden = false;
        const rect = this.el.getBoundingClientRect();
        const left = Math.max(edge, Math.min(x, window.innerWidth - rect.width - edge));
        const top = y + rect.height + edge > window.innerHeight ? Math.max(edge, y - rect.height) : y;
        this.el.style.left = `${left}px`;
        this.el.style.top = `${top}px`;
        this.items()[0]?.focus();
    }

    // openBelow shows the menu under button, right-aligned with it.
    openBelow(button: HTMLElement) {
        const r = button.getBoundingClientRect();
        this.el.hidden = false;
        const width = this.el.getBoundingClientRect().width;
        this.openAt(r.right - width, r.bottom + 4);
        button.setAttribute("aria-expanded", "true");
    }

    close(restoreFocus = true) {
        if (!this.isOpen()) {
            return;
        }
        this.el.hidden = true;
        this.opener?.setAttribute("aria-expanded", "false");
        if (restoreFocus) {
            this.returnFocus?.focus();
        }
        this.returnFocus = null;
    }

    private items(): HTMLButtonElement[] {
        return Array.from(this.el.querySelectorAll<HTMLButtonElement>("[role=menuitem]:not(:disabled):not([hidden])"));
    }

    private onKey(e: KeyboardEvent) {
        const items = this.items();
        const i = items.indexOf(document.activeElement as HTMLButtonElement);
        let next = -1;
        switch (e.key) {
            case "ArrowDown":
                next = (i + 1) % items.length;
                break;
            case "ArrowUp":
                next = (i - 1 + items.length) % items.length;
                break;
            case "Home":
                next = 0;
                break;
            case "End":
                next = items.length - 1;
                break;
            case "Escape":
                this.close();
                break;
            case "Tab":
                this.close();
                break;
            case "Enter":
            case " ":
                // The button's click handler runs the item.
                e.stopPropagation();
                return;
            default:
                return;
        }
        // The menu owns these keys; the window's handlers must not see them.
        e.preventDefault();
        e.stopPropagation();
        items[next]?.focus();
    }
}

let mainMenu: PopupMenu;
let rowMenu: PopupMenu;
let hooks: MenuHooks;

const menuButton = document.getElementById("menu-btn") as HTMLButtonElement;
const pinItem = document.getElementById("row-menu-pin") as HTMLButtonElement;

// isMenuOpen reports whether a popup menu has the keyboard.
export function isMenuOpen(): boolean {
    return mainMenu.isOpen() || rowMenu.isOpen();
}

// openRowMenu shows the context menu for the selected row, under it (for
// Shift+F10 and the Menu key) or at x, y (for a right-click).
export function openRowMenu(x?: number, y?: number) {
    const sel = hooks.selected();
    if (!sel) {
        return;
    }
    pinItem.textContent = sel.bookmark.pinned ? "Unpin" : "Pin";
    if (x === undefined || y === undefined) {
        const r = sel.row.getBoundingClientRect();
        x = r.left + 24;
        y = r.bottom;
    }
    mainMenu.close(false);
    rowMenu.openAt(x, y);
}

export function initMenus(h: MenuHooks) {
    hooks = h;
    mainMenu = new PopupMenu(document.getElementById("main-menu") as HTMLElement, h.run, menuButton);
    rowMenu = new PopupMenu(document.getElementById("row-menu") as HTMLElement, h.run);

    menuButton.setAttribute("aria-expanded", "false");
    menuButton.addEventListener("click", () => {
        if (mainMenu.isOpen()) {
            mainMenu.close();
            return;
        }
        rowMenu.close(false);
        mainMenu.openBelow(menuButton);
    });

    const results = document.getElementById("results") as HTMLElement;
    results.addEventListener("contextmenu", (e) => {
        const i = hooks.rowAt(e.target as Element);
        e.preventDefault();
        if (i < 0) {
            return;
        }
        hooks.select(i);
        openRowMenu(e.clientX, e.clientY);
    });
}

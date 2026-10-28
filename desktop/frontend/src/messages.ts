// Messages by urgency: toasts for passing confirmations, a banner for problems
// that persist, and the status area for background progress. All text is set
// via textContent.

const toastMs = 4000;

const toastEl = document.getElementById("toast") as HTMLElement;
const toastText = document.getElementById("toast-text") as HTMLSpanElement;
const toastAction = document.getElementById("toast-action") as HTMLButtonElement;
const bannerEl = document.getElementById("banner") as HTMLElement;
const bannerText = document.getElementById("banner-text") as HTMLSpanElement;
const bannerDismiss = document.getElementById("banner-dismiss") as HTMLButtonElement;
const statusEl = document.getElementById("status") as HTMLSpanElement;

export interface ToastAction {
    label: string;
    run(): void;
}

let toastTimer: ReturnType<typeof setTimeout> | undefined;
let toastRun: (() => void) | undefined;
// Whether the visible toast came from setStatus, so setStatus("") clears it.
let toastIsStatus = false;

// showToast shows a passing confirmation for about 4 seconds, replacing any
// toast already shown. The optional action adds a button (e.g. "Show").
export function showToast(text: string, action?: ToastAction) {
    clearTimeout(toastTimer);
    toastIsStatus = false;
    toastText.textContent = text;
    toastRun = action?.run;
    toastAction.textContent = action?.label ?? "";
    toastAction.hidden = !action;
    toastEl.hidden = false;
    toastTimer = setTimeout(hideToast, toastMs);
}

function hideToast() {
    clearTimeout(toastTimer);
    toastEl.hidden = true;
    toastRun = undefined;
    toastIsStatus = false;
}

toastAction.addEventListener("click", () => {
    const run = toastRun;
    hideToast();
    run?.();
});

// showBanner shows a persistent problem under the title area. It stays until
// dismissed or clearBanner (sync:ok), and comes back on the next call.
export function showBanner(text: string) {
    bannerText.textContent = text;
    bannerEl.hidden = false;
}

export function clearBanner() {
    bannerEl.hidden = true;
    bannerText.textContent = "";
}

bannerDismiss.addEventListener("click", clearBanner);

// setStatus reports background progress: in the footer status area on macOS
// and KDE, and as a toast on GNOME, which has no footer. "" clears it.
export function setStatus(message: string) {
    if (document.documentElement.dataset.desktop === "gnome") {
        if (message) {
            showToast(message);
            toastIsStatus = true;
        } else if (toastIsStatus) {
            hideToast();
        }
        return;
    }
    statusEl.textContent = message;
    // Long messages are truncated in the footer; show the full text on hover.
    statusEl.title = message;
}

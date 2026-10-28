// The selection and default buttons follow the system accent colour: CSS
// AccentColor on macOS, the Settings portal's colour on Linux (sent in
// PlatformInfo.accentColor). The colour is darkened until white text on it
// reaches 4.5:1, and a text variant is adjusted to 4.5:1 on the list.

const minContrast = 4.5;
// Each step mixes this much black (or white) into the colour.
const step = 0.05;

type RGB = [number, number, number];

let systemAccent = "";

// applyAccent sets --accent and --accent-text from the platform's accent, or
// removes them so the stylesheet's fallback tokens apply.
export function applyAccent(linuxAccent: string | undefined) {
    systemAccent = linuxAccent ?? "";
    update();
}

function update() {
    const root = document.documentElement.style;
    // Start from the stylesheet's fallback tokens, not colours set here earlier.
    root.removeProperty("--accent");
    root.removeProperty("--accent-text");
    const base = parse(systemAccent) ?? macAccent();
    if (!base) {
        return;
    }
    const list = resolve("background-color", "var(--list)") ?? [255, 255, 255];
    const towards: RGB = luminance(list) > 0.5 ? [0, 0, 0] : [255, 255, 255];
    root.setProperty("--accent", hex(adjust(base, [255, 255, 255], [0, 0, 0])));
    root.setProperty("--accent-text", hex(adjust(base, list, towards)));
}

// macAccent resolves CSS AccentColor where the web view supports it.
function macAccent(): RGB | undefined {
    if (document.documentElement.dataset.os !== "darwin" || !CSS.supports("color", "AccentColor")) {
        return undefined;
    }
    return resolve("color", "AccentColor");
}

// resolve computes a CSS colour value (a keyword or var()) to RGB.
function resolve(property: "color" | "background-color", value: string): RGB | undefined {
    const probe = document.createElement("span");
    probe.style.setProperty(property, value);
    probe.style.display = "none";
    document.body.append(probe);
    const color = parse(getComputedStyle(probe).getPropertyValue(property));
    probe.remove();
    return color;
}

// adjust mixes c towards `towards` until text of colour `on` reaches 4.5:1.
function adjust(c: RGB, on: RGB, towards: RGB): RGB {
    let out = c;
    for (let t = step; contrast(out, on) < minContrast && t <= 1; t += step) {
        out = c.map((v, i) => Math.round(v + (towards[i] - v) * t)) as RGB;
    }
    return out;
}

function parse(css: string): RGB | undefined {
    let m = /^#([0-9a-f]{2})([0-9a-f]{2})([0-9a-f]{2})$/i.exec(css.trim());
    if (m) {
        return [parseInt(m[1], 16), parseInt(m[2], 16), parseInt(m[3], 16)];
    }
    m = /^rgba?\(\s*(\d+)[,\s]+(\d+)[,\s]+(\d+)/.exec(css.trim());
    return m ? [Number(m[1]), Number(m[2]), Number(m[3])] : undefined;
}

function hex(c: RGB): string {
    return "#" + c.map((v) => v.toString(16).padStart(2, "0")).join("");
}

// luminance is the WCAG relative luminance.
function luminance(c: RGB): number {
    const [r, g, b] = c.map((v) => {
        const s = v / 255;
        return s <= 0.03928 ? s / 12.92 : Math.pow((s + 0.055) / 1.055, 2.4);
    });
    return 0.2126 * r + 0.7152 * g + 0.0722 * b;
}

function contrast(a: RGB, b: RGB): number {
    const [hi, lo] = [luminance(a), luminance(b)].sort((x, y) => y - x);
    return (hi + 0.05) / (lo + 0.05);
}

// The text variant depends on the list colour, which follows light/dark mode.
window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", update);

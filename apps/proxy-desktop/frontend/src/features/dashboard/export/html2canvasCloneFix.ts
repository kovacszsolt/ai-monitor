/** html2canvas cannot parse Tailwind v4 / modern CSS color functions (oklab, oklch, color-mix, …). */

const MODERN_COLOR_RE = /oklab|oklch|color-mix\(|color\(|(?:^|[\s,/])(?:lab|lch)\(/i;
const COLOR_FUNC_START_RE = /(?:color-mix|oklab|oklch|lab|lch|color)\(/gi;

const EXPORT_CANVAS_BG = "#131315";

const COLOR_PROPERTIES = [
  "color",
  "background-color",
  "border-top-color",
  "border-right-color",
  "border-bottom-color",
  "border-left-color",
  "outline-color",
  "text-decoration-color",
  "caret-color",
  "column-rule-color",
  "-webkit-text-stroke-color",
] as const;

const COLOR_IN_SHORTHANDS = ["box-shadow", "text-shadow", "background", "border-color"] as const;

const SVG_COLOR_PROPS = ["fill", "stroke", "color", "stop-color", "flood-color", "lighting-color"] as const;

let rgbCanvas: HTMLCanvasElement | null = null;

function needsColorSanitize(value: string): boolean {
  return MODERN_COLOR_RE.test(value);
}

function colorToRgb(color: string): string {
  const trimmed = color.trim();
  if (!trimmed || trimmed === "none" || trimmed === "transparent") return trimmed;
  if (!needsColorSanitize(trimmed)) return trimmed;
  if (!rgbCanvas) rgbCanvas = document.createElement("canvas");
  const ctx = rgbCanvas.getContext("2d");
  if (!ctx) return "rgb(0, 0, 0)";
  try {
    ctx.fillStyle = "#000000";
    ctx.fillStyle = trimmed;
    ctx.fillRect(0, 0, 1, 1);
    const [r, g, b, a] = ctx.getImageData(0, 0, 1, 1).data;
    return a < 255 ? `rgba(${r}, ${g}, ${b}, ${(a / 255).toFixed(3)})` : `rgb(${r}, ${g}, ${b})`;
  } catch {
    return "rgb(0, 0, 0)";
  }
}

function findModernColorFunctionSpan(text: string, fromIndex = 0): { start: number; end: number } | null {
  COLOR_FUNC_START_RE.lastIndex = fromIndex;
  const m = COLOR_FUNC_START_RE.exec(text);
  if (!m) return null;
  const openParen = m.index + m[0].length - 1;
  let depth = 0;
  for (let i = openParen; i < text.length; i++) {
    const ch = text[i];
    if (ch === "(") depth++;
    else if (ch === ")") {
      depth--;
      if (depth === 0) return { start: m.index, end: i + 1 };
    }
  }
  return null;
}

export function sanitizeCSSValue(value: string): string {
  if (!needsColorSanitize(value)) return value;
  let result = value;
  for (let guard = 0; guard < 64; guard++) {
    const span = findModernColorFunctionSpan(result);
    if (!span) break;
    const fn = result.slice(span.start, span.end);
    const rgb = colorToRgb(fn);
    result = result.slice(0, span.start) + rgb + result.slice(span.end);
    if (!needsColorSanitize(result)) break;
  }
  if (needsColorSanitize(result)) {
    return colorToRgb(result);
  }
  return result;
}

function stripAllStylesheets(clonedDoc: Document): void {
  clonedDoc.querySelectorAll("style, link[rel='stylesheet']").forEach((node) => node.remove());
}

function disableAllStylesheets(clonedDoc: Document): void {
  try {
    for (let i = 0; i < clonedDoc.styleSheets.length; i++) {
      const sheet = clonedDoc.styleSheets[i];
      if (sheet) sheet.disabled = true;
    }
  } catch {
    // ignore SecurityError on cross-origin sheets
  }
}

function sanitizeDocumentShell(clonedDoc: Document): void {
  const html = clonedDoc.documentElement;
  html.style.setProperty("background-color", EXPORT_CANVAS_BG, "important");
  html.style.setProperty("color", "rgb(229, 225, 228)", "important");
  const body = clonedDoc.body;
  if (body) {
    body.style.setProperty("background-color", EXPORT_CANVAS_BG, "important");
    body.style.setProperty("color", "rgb(229, 225, 228)", "important");
  }
}

function stripClasses(root: HTMLElement): void {
  root.removeAttribute("class");
  root.querySelectorAll("[class]").forEach((el) => el.removeAttribute("class"));
}

function sanitizeInlineStyleAttribute(el: HTMLElement | SVGElement): void {
  const inline = el.getAttribute("style");
  if (!inline || !needsColorSanitize(inline)) return;
  el.setAttribute("style", sanitizeCSSValue(inline));
}

const SKIP_INLINE_PROPS = new Set(["content", "webkitappregion"]);

function sanitizePropertyValue(prop: string, value: string): string {
  if (!needsColorSanitize(value)) return value;
  let sanitized = sanitizeCSSValue(value);
  if (needsColorSanitize(sanitized)) {
    if (prop === "background-image" || prop.includes("image")) {
      sanitized = "none";
    } else if (
      prop.endsWith("-color") ||
      prop === "color" ||
      prop === "fill" ||
      prop === "stroke" ||
      prop === "caret-color"
    ) {
      sanitized = colorToRgb(value);
    } else {
      sanitized = sanitizeCSSValue(value);
    }
  }
  return sanitized;
}

function inlineAllComputedStyles(source: HTMLElement, clone: HTMLElement): void {
  const view = source.ownerDocument.defaultView;
  const computed = view?.getComputedStyle(source);
  if (!computed) return;

  for (let i = 0; i < computed.length; i++) {
    const prop = computed.item(i);
    if (!prop || SKIP_INLINE_PROPS.has(prop)) continue;
    const value = computed.getPropertyValue(prop);
    if (!value) continue;

    const sanitized = sanitizePropertyValue(prop, value);

    try {
      clone.style.setProperty(prop, sanitized, computed.getPropertyPriority(prop));
    } catch {
      // ignore unsupported properties on clone
    }
  }
}

function inlineSvgComputedColors(source: SVGElement, clone: SVGElement): void {
  const computed = source.ownerDocument.defaultView?.getComputedStyle(source);
  if (!computed) return;

  for (const prop of SVG_COLOR_PROPS) {
    const value = computed.getPropertyValue(prop);
    if (!value || value === "none") continue;
    const sanitized = sanitizePropertyValue(prop, value);
    try {
      clone.style.setProperty(prop, sanitized, "important");
    } catch {
      // ignore
    }
    if (prop === "fill" || prop === "stroke") {
      clone.setAttribute(prop, sanitized);
    }
  }
}

function walkInlineAllComputed(source: Element, clone: Element): void {
  if (source instanceof HTMLElement && clone instanceof HTMLElement) {
    inlineAllComputedStyles(source, clone);
    sanitizeInlineStyleAttribute(clone);
  } else if (source instanceof SVGElement && clone instanceof SVGElement) {
    inlineSvgComputedColors(source, clone);
    sanitizeInlineStyleAttribute(clone);
  }
  const sourceChildren = source.children;
  const cloneChildren = clone.children;
  for (let i = 0; i < sourceChildren.length; i++) {
    const cloneChild = cloneChildren.item(i);
    if (cloneChild) walkInlineAllComputed(sourceChildren[i], cloneChild);
  }
}

function forceSafeColorsOnClone(clone: HTMLElement): void {
  const win = clone.ownerDocument.defaultView;
  const computed = win?.getComputedStyle(clone);
  if (!computed) return;

  for (const prop of COLOR_PROPERTIES) {
    const value = computed.getPropertyValue(prop);
    if (!value || value === "none") continue;
    if (needsColorSanitize(value)) {
      clone.style.setProperty(prop, colorToRgb(value), "important");
    }
  }

  for (const prop of COLOR_IN_SHORTHANDS) {
    const value = computed.getPropertyValue(prop);
    if (!value) continue;
    if (needsColorSanitize(value)) {
      const sanitized = sanitizePropertyValue(prop, value);
      clone.style.setProperty(prop, sanitized, "important");
    }
  }

  clone.style.setProperty("backdrop-filter", "none", "important");
  clone.style.setProperty("filter", "none", "important");
}

function forceSafeColorsOnSvgClone(clone: SVGElement): void {
  const computed = clone.ownerDocument.defaultView?.getComputedStyle(clone);
  if (!computed) return;
  for (const prop of SVG_COLOR_PROPS) {
    const value = computed.getPropertyValue(prop);
    if (!value || value === "none") continue;
    if (needsColorSanitize(value)) {
      const rgb = colorToRgb(value);
      clone.style.setProperty(prop, rgb, "important");
      if (prop === "fill" || prop === "stroke") {
        clone.setAttribute(prop, rgb);
      }
    }
  }
}

function walkForceSafeColors(cloneRoot: HTMLElement): void {
  forceSafeColorsOnClone(cloneRoot);
  cloneRoot.querySelectorAll("*").forEach((node) => {
    if (node instanceof HTMLElement) {
      forceSafeColorsOnClone(node);
    } else if (node instanceof SVGElement) {
      forceSafeColorsOnSvgClone(node);
    }
  });
}

function isExportableElement(el: Element): el is HTMLElement | SVGElement {
  return el instanceof HTMLElement || el instanceof SVGElement;
}

function applyFlattenToSourceElement(el: HTMLElement | SVGElement): void {
  const computed = el.ownerDocument.defaultView?.getComputedStyle(el);
  if (!computed) return;

  for (let i = 0; i < computed.length; i++) {
    const prop = computed.item(i);
    if (!prop || SKIP_INLINE_PROPS.has(prop)) continue;
    const value = computed.getPropertyValue(prop);
    if (!value || !needsColorSanitize(value)) continue;
    try {
      el.style.setProperty(prop, sanitizePropertyValue(prop, value), "important");
    } catch {
      // ignore
    }
  }

  if (el instanceof SVGElement) {
    for (const prop of SVG_COLOR_PROPS) {
      const value = computed.getPropertyValue(prop);
      if (!value || value === "none" || !needsColorSanitize(value)) continue;
      const sanitized = sanitizePropertyValue(prop, value);
      el.setAttribute(prop, sanitized);
      el.style.setProperty(prop, sanitized, "important");
    }
  }
}

function sanitizeElementInlineStyles(el: HTMLElement | SVGElement): void {
  for (let i = 0; i < el.style.length; i++) {
    const prop = el.style.item(i);
    if (!prop) continue;
    const value = el.style.getPropertyValue(prop);
    if (needsColorSanitize(value)) {
      el.style.setProperty(prop, sanitizePropertyValue(prop, value), "important");
    }
  }
  sanitizeInlineStyleAttribute(el);
  if (el instanceof SVGElement) {
    for (const attr of ["fill", "stroke"] as const) {
      const v = el.getAttribute(attr);
      if (v && needsColorSanitize(v)) {
        el.setAttribute(attr, colorToRgb(v));
      }
    }
  }
}

function sanitizeEntireCloneSubtree(root: Element): void {
  const visit = (el: Element) => {
    if (isExportableElement(el)) {
      sanitizeElementInlineStyles(el);
      const win = el.ownerDocument.defaultView;
      const computed = win?.getComputedStyle(el);
      if (computed) {
        for (const prop of COLOR_PROPERTIES) {
          const value = computed.getPropertyValue(prop);
          if (value && needsColorSanitize(value)) {
            el.style.setProperty(prop, colorToRgb(value), "important");
          }
        }
        for (const prop of COLOR_IN_SHORTHANDS) {
          const value = computed.getPropertyValue(prop);
          if (value && needsColorSanitize(value)) {
            el.style.setProperty(prop, sanitizePropertyValue(prop, value), "important");
          }
        }
      }
    }
    for (const child of el.children) visit(child);
  };
  visit(root);
}

const EXPORT_PSEUDO_STYLE_ID = "dashboard-export-suppress-pseudos";

type ElementStyleBackup = {
  style: string | null;
  fill: string | null;
  stroke: string | null;
};

function installExportPseudoSuppress(): void {
  if (document.getElementById(EXPORT_PSEUDO_STYLE_ID)) return;
  const style = document.createElement("style");
  style.id = EXPORT_PSEUDO_STYLE_ID;
  style.textContent = `
    *::before,
    *::after {
      content: none !important;
      display: none !important;
    }
  `;
  document.head.appendChild(style);
}

function removeExportPseudoSuppress(): void {
  document.getElementById(EXPORT_PSEUDO_STYLE_ID)?.remove();
}

/**
 * html2canvas copies getComputedStyle (including oklab) during clone — before onclone.
 * Flatten colors on the live export subtree first, then restore after capture.
 */
export function beginExportColorFlatten(exportRoot: HTMLElement): () => void {
  installExportPseudoSuppress();
  const backups = new WeakMap<Element, ElementStyleBackup>();
  const scope = exportRoot.ownerDocument?.documentElement ?? exportRoot;
  const nodes: Element[] = [scope, ...scope.querySelectorAll("*")];

  for (const el of nodes) {
    if (!isExportableElement(el)) continue;
    backups.set(el, {
      style: el.getAttribute("style"),
      fill: el instanceof SVGElement ? el.getAttribute("fill") : null,
      stroke: el instanceof SVGElement ? el.getAttribute("stroke") : null,
    });
    applyFlattenToSourceElement(el);
  }

  return () => {
    for (const el of nodes) {
      if (!isExportableElement(el)) continue;
      const backup = backups.get(el);
      if (!backup) continue;
      if (backup.style === null) el.removeAttribute("style");
      else el.setAttribute("style", backup.style);
      if (el instanceof SVGElement) {
        if (backup.fill === null) el.removeAttribute("fill");
        else el.setAttribute("fill", backup.fill);
        if (backup.stroke === null) el.removeAttribute("stroke");
        else el.setAttribute("stroke", backup.stroke);
      }
    }
    removeExportPseudoSuppress();
  };
}

/** Make cloned DOM safe for html2canvas (no oklab in stylesheets — inline resolved RGB only). */
export function prepareHtml2CanvasClone(
  sourceRoot: HTMLElement,
  clonedDoc: Document,
  cloneRoot: HTMLElement,
): void {
  stripAllStylesheets(clonedDoc);
  disableAllStylesheets(clonedDoc);
  sanitizeDocumentShell(clonedDoc);
  walkInlineAllComputed(sourceRoot, cloneRoot);
  walkForceSafeColors(cloneRoot);
  sanitizeEntireCloneSubtree(cloneRoot);
  stripClasses(cloneRoot);
  clonedDoc.querySelectorAll("html2canvaspseudoelement").forEach((node) => {
    if (node instanceof HTMLElement) sanitizeEntireCloneSubtree(node);
  });
  cloneRoot.style.setProperty("background-color", colorToRgb(getComputedBackground(sourceRoot)), "important");
}

function getComputedBackground(el: HTMLElement): string {
  const bg = el.ownerDocument.defaultView?.getComputedStyle(el).backgroundColor;
  return bg && bg !== "transparent" ? bg : EXPORT_CANVAS_BG;
}

import { beginExportColorFlatten, prepareHtml2CanvasClone } from "@/features/dashboard/export/html2canvasCloneFix";
import { SaveExportFile } from "../../../../wailsjs/go/main/App";

export type DashboardExportFormat = "png" | "html" | "pdf";

export type DashboardExportMeta = {
  layoutName: string;
  exportedAtIso: string;
  filterSummary: string;
};

function slugFilenamePart(value: string): string {
  const slug = value
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
  return slug || "dashboard";
}

export function defaultExportFilename(layoutName: string, format: DashboardExportFormat): string {
  const date = new Date().toISOString().slice(0, 10);
  const ext = format === "html" ? "html" : format;
  return `${slugFilenamePart(layoutName)}-${date}.${ext}`;
}

function isWailsApp(): boolean {
  return typeof window !== "undefined" && typeof (window as Window & { go?: unknown }).go !== "undefined";
}

function bytesToBase64(bytes: Uint8Array): string {
  let binary = "";
  const chunk = 0x8000;
  for (let i = 0; i < bytes.length; i += chunk) {
    binary += String.fromCharCode(...bytes.subarray(i, i + chunk));
  }
  return btoa(binary);
}

async function saveFile(filename: string, bytes: Uint8Array, mime: string): Promise<string | null> {
  const base64 = bytesToBase64(bytes);
  if (isWailsApp()) {
    const savedPath = await SaveExportFile(filename, base64);
    if (savedPath === "") return null;
    if (!savedPath) {
      throw new Error("Export failed — save dialog did not return a path");
    }
    return savedPath;
  }

  try {
    const savedPath = await SaveExportFile(filename, base64);
    if (savedPath) return savedPath;
    if (savedPath === "") return null;
  } catch {
    // vite dev without Wails — fall through to browser download
  }

  const blob = new Blob([bytes as BlobPart], { type: mime });
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = filename;
  anchor.rel = "noopener";
  anchor.click();
  URL.revokeObjectURL(url);
  return filename;
}

async function waitForCharts(): Promise<void> {
  await new Promise((resolve) => requestAnimationFrame(() => resolve(undefined)));
  await new Promise((resolve) => setTimeout(resolve, 400));
}

export async function captureDashboardElement(element: HTMLElement): Promise<HTMLCanvasElement> {
  const scrollParent = element.closest("main");
  scrollParent?.scrollTo({ top: 0, behavior: "instant" });
  await waitForCharts();

  const endFlatten = beginExportColorFlatten(element);
  await new Promise<void>((resolve) => {
    requestAnimationFrame(() => requestAnimationFrame(() => resolve()));
  });

  const { default: html2canvas } = await import("html2canvas");
  try {
    return await html2canvas(element, {
      backgroundColor: "#131315",
      scale: 2,
      logging: false,
      useCORS: true,
      scrollX: 0,
      scrollY: 0,
      onclone: (clonedDoc, cloneRoot) => {
        prepareHtml2CanvasClone(element, clonedDoc, cloneRoot);
      },
    });
  } catch (err) {
    const message = err instanceof Error ? err.message : String(err);
    throw new Error(`Dashboard capture failed: ${message}`);
  } finally {
    endFlatten();
  }
}

function buildHtmlDocument(meta: DashboardExportMeta, pngDataUrl: string): string {
  const title = `Dashboard — ${meta.layoutName}`;
  const escapedLayout = escapeHtml(meta.layoutName);
  const escapedFilters = escapeHtml(meta.filterSummary);
  const escapedAt = escapeHtml(meta.exportedAtIso);
  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>${escapeHtml(title)}</title>
  <style>
    body { margin: 0; padding: 1.5rem; font-family: ui-sans-serif, system-ui, sans-serif; background: #09090b; color: #fafafa; }
    h1 { font-size: 1.25rem; margin: 0 0 0.5rem; }
    .meta { color: #a1a1aa; font-size: 0.875rem; margin: 0 0 1.25rem; line-height: 1.5; }
    img { display: block; max-width: 100%; height: auto; border-radius: 0.5rem; border: 1px solid #27272a; }
  </style>
</head>
<body>
  <h1>Dashboard — ${escapedLayout}</h1>
  <p class="meta">Exported ${escapedAt}<br />${escapedFilters}</p>
  <img src="${pngDataUrl}" alt="Dashboard export" />
</body>
</html>`;
}

function escapeHtml(value: string): string {
  return value
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}

function canvasToPngBytes(canvas: HTMLCanvasElement): Promise<Uint8Array> {
  return new Promise((resolve, reject) => {
    canvas.toBlob((blob) => {
      if (!blob) {
        reject(new Error("Failed to encode PNG"));
        return;
      }
      blob.arrayBuffer().then((buf) => resolve(new Uint8Array(buf)), reject);
    }, "image/png");
  });
}

async function canvasToPdfBytes(canvas: HTMLCanvasElement, meta: DashboardExportMeta): Promise<Uint8Array> {
  const { jsPDF } = await import("jspdf");
  const imgData = canvas.toDataURL("image/png");
  const pdf = new jsPDF({ orientation: "portrait", unit: "mm", format: "a4" });
  const pageWidth = pdf.internal.pageSize.getWidth();
  const pageHeight = pdf.internal.pageSize.getHeight();
  const margin = 10;
  const contentWidth = pageWidth - margin * 2;

  pdf.setFontSize(14);
  pdf.text(`Dashboard — ${meta.layoutName}`, margin, margin);
  pdf.setFontSize(9);
  pdf.setTextColor(120);
  const metaLines = pdf.splitTextToSize(
    `Exported ${meta.exportedAtIso}\n${meta.filterSummary}`,
    contentWidth,
  );
  pdf.text(metaLines, margin, margin + 7);

  const metaBlockHeight = 7 + metaLines.length * 4;
  const imgTop = margin + metaBlockHeight + 4;
  const imgHeight = (canvas.height * contentWidth) / canvas.width;
  const availableFirstPage = pageHeight - imgTop - margin;

  let heightLeft = imgHeight;
  let position = 0;

  pdf.addImage(imgData, "PNG", margin, imgTop, contentWidth, imgHeight);
  heightLeft -= availableFirstPage;

  while (heightLeft > 0) {
    pdf.addPage();
    position = heightLeft - imgHeight;
    pdf.addImage(imgData, "PNG", margin, position, contentWidth, imgHeight);
    heightLeft -= pageHeight;
  }

  pdf.setTextColor(0);
  const out = pdf.output("arraybuffer");
  return new Uint8Array(out);
}

export async function exportDashboard(
  element: HTMLElement,
  format: DashboardExportFormat,
  meta: DashboardExportMeta,
): Promise<string | null> {
  const canvas = await captureDashboardElement(element);
  const filename = defaultExportFilename(meta.layoutName, format);

  if (format === "png") {
    const bytes = await canvasToPngBytes(canvas);
    return saveFile(filename, bytes, "image/png");
  }

  if (format === "html") {
    const pngDataUrl = canvas.toDataURL("image/png");
    const html = buildHtmlDocument(meta, pngDataUrl);
    const bytes = new TextEncoder().encode(html);
    return saveFile(filename, bytes, "text/html;charset=utf-8");
  }

  const pdfBytes = await canvasToPdfBytes(canvas, meta);
  return saveFile(filename, pdfBytes, "application/pdf");
}

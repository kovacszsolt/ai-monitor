import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { service } from "../../../../wailsjs/go/models";

const PROTOCOLS = ["openai_chat", "openai_responses", "openai_images", "anthropic_messages"] as const;
export const MODEL_MODALITIES = ["text-text", "text-image", "text_image-text"] as const;
export type ModelModality = (typeof MODEL_MODALITIES)[number];

export type ModelRouteFormProps = {
  publicName?: string;
  setPublicName?: (v: string) => void;
  showPublic?: boolean;
  providerID?: string;
  setProviderID?: (v: string) => void;
  providers?: service.ProviderDTO[];
  upstream: string;
  setUpstream: (v: string) => void;
  protocol: string;
  setProtocol: (v: string) => void;
  modality: string;
  setModality: (v: string) => void;
  inputUSD: string;
  setInputUSD: (v: string) => void;
  cachedUSD: string;
  setCachedUSD: (v: string) => void;
  outputUSD: string;
  setOutputUSD: (v: string) => void;
};

export function ModelRouteForm(props: ModelRouteFormProps) {
  return (
    <div className="grid gap-3 py-2">
      {props.showPublic && (
        <div className="grid gap-2">
          <Label>Public model name</Label>
          <Input value={props.publicName ?? ""} onChange={(e) => props.setPublicName?.(e.target.value)} />
        </div>
      )}
      {props.providers && props.setProviderID && (
        <div className="grid gap-2">
          <Label>Provider</Label>
          <select
            className="h-9 rounded-md border border-border bg-card px-3 text-sm"
            value={props.providerID}
            onChange={(e) => props.setProviderID?.(e.target.value)}
          >
            {props.providers.map((p) => (
              <option key={p.id} value={p.id}>{p.name}</option>
            ))}
          </select>
        </div>
      )}
      <div className="grid gap-2">
        <Label>Upstream model</Label>
        <Input value={props.upstream} onChange={(e) => props.setUpstream(e.target.value)} />
      </div>
      <div className="grid gap-2">
        <Label>Protocol</Label>
        <select
          className="h-9 rounded-md border border-border bg-card px-3 text-sm"
          value={props.protocol}
          onChange={(e) => props.setProtocol(e.target.value)}
        >
          {PROTOCOLS.map((p) => <option key={p} value={p}>{p}</option>)}
        </select>
      </div>
      <div className="grid gap-2">
        <Label>Modality</Label>
        <select
          className="h-9 rounded-md border border-border bg-card px-3 text-sm"
          value={props.modality}
          onChange={(e) => props.setModality(e.target.value)}
        >
          {MODEL_MODALITIES.map((m) => (
            <option key={m} value={m}>{m}</option>
          ))}
        </select>
        <p className="text-xs text-muted-foreground">
          Input→output: text-text, text-image (generation), text_image-text (vision).
        </p>
      </div>
      <p className="text-xs text-muted-foreground">Token pricing (optional, USD per 1M tokens). Leave blank if unknown.</p>
      <div className="grid grid-cols-3 gap-2">
        <div className="grid gap-2">
          <Label>Input</Label>
          <Input inputMode="decimal" placeholder="e.g. 3" value={props.inputUSD} onChange={(e) => props.setInputUSD(e.target.value)} />
        </div>
        <div className="grid gap-2">
          <Label>Cached</Label>
          <Input inputMode="decimal" placeholder="e.g. 0.75" value={props.cachedUSD} onChange={(e) => props.setCachedUSD(e.target.value)} />
        </div>
        <div className="grid gap-2">
          <Label>Output</Label>
          <Input inputMode="decimal" placeholder="e.g. 15" value={props.outputUSD} onChange={(e) => props.setOutputUSD(e.target.value)} />
        </div>
      </div>
    </div>
  );
}

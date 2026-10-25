export type ChatRole = "user" | "assistant" | "system";

export type ChatMessage = {
  id: string;
  role: ChatRole;
  content: string;
};

export const CHAT_REQUEST_TAG = "desktop-chat";

export const CHAT_PREFS_STORAGE_KEY = "proxy-desktop.chat.preferences";

export type ChatPreferencesV1 = {
  projectId: string;
  providerId: string;
  modelRouteId: string;
  apiKeyId: string;
  stream: boolean;
};

export const defaultChatPreferences = (): ChatPreferencesV1 => ({
  projectId: "",
  providerId: "",
  modelRouteId: "",
  apiKeyId: "",
  stream: false,
});

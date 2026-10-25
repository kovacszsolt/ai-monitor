import { useCallback, useMemo, useState } from "react";
import {
  CHAT_PREFS_STORAGE_KEY,
  defaultChatPreferences,
  type ChatPreferencesV1,
} from "@/features/chat/types";

function readStoredPreferences(): ChatPreferencesV1 {
  try {
    const raw = localStorage.getItem(CHAT_PREFS_STORAGE_KEY);
    if (!raw) return defaultChatPreferences();
    const parsed = JSON.parse(raw) as Partial<ChatPreferencesV1>;
    return {
      ...defaultChatPreferences(),
      ...parsed,
      stream: Boolean(parsed.stream),
    };
  } catch {
    return defaultChatPreferences();
  }
}

function writeStoredPreferences(prefs: ChatPreferencesV1) {
  try {
    localStorage.setItem(CHAT_PREFS_STORAGE_KEY, JSON.stringify(prefs));
  } catch {
    // ignore quota / private mode
  }
}

export function useChatPreferences() {
  const [preferences, setPreferencesState] = useState<ChatPreferencesV1>(readStoredPreferences);

  const setPreferences = useCallback((patch: Partial<ChatPreferencesV1>) => {
    setPreferencesState((prev) => {
      const next = { ...prev, ...patch };
      writeStoredPreferences(next);
      return next;
    });
  }, []);

  const setPreference = useCallback(
    <K extends keyof ChatPreferencesV1>(key: K, value: ChatPreferencesV1[K]) => {
      setPreferences({ [key]: value });
    },
    [setPreferences],
  );

  return useMemo(
    () => ({
      preferences,
      setPreferences,
      setPreference,
    }),
    [preferences, setPreferences, setPreference],
  );
}

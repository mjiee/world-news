import { create } from "zustand";
import { saveSystemConfig } from "@/services";
import { isWeb } from "@/utils/platform";

interface LanguageState {
  language: string;
  setLanguage: (newLanguage: string) => void;
}

export const useLanguageStore = create<LanguageState>((set) => {
  return {
    language: "en",
    setLanguage: (newLanguage: string) => {
      if (!isWeb()) return saveSystemConfig({ key: "language", value: newLanguage });

      set((state) => ({ ...state, language: newLanguage }));
    },
  };
});

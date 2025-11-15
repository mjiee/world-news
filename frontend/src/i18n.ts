// src/i18n.ts
import i18n from "i18next";
import LanguageDetector from "i18next-browser-languagedetector";
import { initReactI18next } from "react-i18next";

import commonEn from "@/assets/locales/en/common.json";
import newsEn from "@/assets/locales/en/news.json";
import settingsEn from "@/assets/locales/en/settings.json";
import task from "@/assets/locales/en/task.json";
import commonZh from "@/assets/locales/zh/common.json";
import newsZh from "@/assets/locales/zh/news.json";
import settingsZh from "@/assets/locales/zh/settings.json";
import taskZh from "@/assets/locales/zh/task.json";

const resources = {
  en: { common: commonEn, settings: settingsEn, news: newsEn, task: task },
  zh: { common: commonZh, settings: settingsZh, news: newsZh, task: taskZh },
};

i18n
  .use(LanguageDetector) // detect user language
  .use(initReactI18next) // pass the i18n instance to react-i18next.
  .init({
    debug: true,
    fallbackLng: "en", // default language
    interpolation: {
      escapeValue: false, // not needed for react as it escapes by default
    },
    defaultNS: "common",
    resources: resources,
  });

export default i18n;

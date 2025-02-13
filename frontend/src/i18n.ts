// src/i18n.ts
import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";

import commonEn from "@/assets/locales/en/common.json";
import commonZh from "@/assets/locales/zh/common.json";
import homeEn from "@/assets/locales/en/home.json";
import homeZh from "@/assets/locales/zh/home.json";
import settingsEn from "@/assets/locales/en/settings.json";
import settingsZh from "@/assets/locales/zh/settings.json";
import newsEn from "@/assets/locales/en/news.json";
import newsZh from "@/assets/locales/zh/news.json";

const resources = {
  en: { common: commonEn, home: homeEn, settings: settingsEn, news: newsEn },
  zh: { common: commonZh, home: homeZh, settings: settingsZh, news: newsZh },
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

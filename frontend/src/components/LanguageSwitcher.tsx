import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { ActionIcon } from "@mantine/core";
import { useLanguageStore } from "@/stores";

const [en, zh] = ["en", "zh"];

export function LanguageSwitcher() {
  const { language, setLanguage } = useLanguageStore((state) => state);
  const { i18n } = useTranslation();

  const handleLanguageChange = () => {
    const lang = i18n.language === en ? zh : en;

    i18n.changeLanguage(lang);
    setLanguage(lang);
  };

  useEffect(() => {
    if (language === i18n.language) return;

    if (i18n.language !== en && i18n.language !== zh) {
      i18n.changeLanguage(language);
    } else if (language !== i18n.language) {
      setLanguage(i18n.language);
    }
  }, []);

  return (
    <ActionIcon variant="default" radius="xl" aria-label="Settings" onClick={handleLanguageChange}>
      {i18n.language === en ? "中" : "EN"}
    </ActionIcon>
  );
}

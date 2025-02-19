import { ActionIcon } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { saveSystemConfig } from "@/services";

const [en, zh] = ["en", "zh"];

export function LanguageSwitcher() {
  const { i18n } = useTranslation();

  const saveLanguage = (language: string) => {
    saveSystemConfig({ key: "language", value: language });
  };

  const handleLanguageChange = () => {
    const language = i18n.language === en ? zh : en;

    saveLanguage(language);
    i18n.changeLanguage(language);
  };

  return (
    <ActionIcon variant="default" radius="xl" aria-label="Settings" onClick={handleLanguageChange}>
      {i18n.language === en ? "中" : "EN"}
    </ActionIcon>
  );
}

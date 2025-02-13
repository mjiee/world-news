import { ActionIcon } from "@mantine/core";
import { useTranslation } from "react-i18next";

export function LanguageSwitcher() {
  const { i18n } = useTranslation();

  const handleLanguageChange = () => {
    if (i18n.language === "en") {
      i18n.changeLanguage("zh");
    } else {
      i18n.changeLanguage("en");
    }
  };

  return (
    <ActionIcon variant="default" radius="xl" aria-label="Settings" onClick={handleLanguageChange}>
      {i18n.language === "en" ? "中" : "EN"}
    </ActionIcon>
  );
}

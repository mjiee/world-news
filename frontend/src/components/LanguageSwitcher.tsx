import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { ActionIcon } from "@mantine/core";
import { saveSystemConfig, SystemConfigKey } from "@/services";
import { GolbalLanguage } from "@/stores/languageStore";
import { isWeb } from "@/utils/platform";

const [en, zh] = ["en", "zh"];

export function LanguageSwitcher() {
  const { i18n } = useTranslation();

  const handleLanguageChange = () => {
    const lang = i18n.language === en ? zh : en;

    if (!isWeb()) saveSystemConfig({ key: SystemConfigKey.Language, value: lang }, true);

    i18n.changeLanguage(lang);
    GolbalLanguage.setLanguage(lang);
  };

  useEffect(() => {
    if (GolbalLanguage.getLanguage() === i18n.language) return;

    if (i18n.language !== en && i18n.language !== zh) {
      i18n.changeLanguage(GolbalLanguage.getLanguage());
    } else if (GolbalLanguage.getLanguage() !== i18n.language) {
      GolbalLanguage.setLanguage(i18n.language);
    }
  }, []);

  return (
    <ActionIcon variant="default" radius="xl" aria-label="Settings" onClick={handleLanguageChange}>
      {i18n.language === en ? "中" : "EN"}
    </ActionIcon>
  );
}

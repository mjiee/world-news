import { useEffect } from "react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";
import { Button, Stack, TextInput, Autocomplete, PasswordInput } from "@mantine/core";
import { useForm } from "@mantine/form";
import { getSystemConfig, saveSystemConfig, SystemConfigKey, TranslaterConfig } from "@/services";
import { isWeb } from "@/utils/platform";

export const translateKey = "translate";
const translationPlatforms = ["google", "microsoft", "aliyun", "baidu"];

// news translation
export function NewsTranslation({ item }: { item: string | null }) {
  const { t } = useTranslation();

  const form = useForm<TranslaterConfig>({
    initialValues: {
      platform: "",
      appId: "",
      appSecret: "",
    },
    validate: {
      platform: (value) => validatRequiredField(value, "translate.label.platform"),
      appId: (value) => validatRequiredField(value, "translate.label.app_id"),
    },
  });

  // validate required field
  function validatRequiredField(value: string, label: string): null | string {
    if (value) return null;
    return t("translate.validate.required", { ns: "settings", label: t(label, { ns: "settings" }) });
  }

  // fetch translate config
  const fetchTranslateConfig = async () => {
    const resp = await getSystemConfig<TranslaterConfig>({ key: SystemConfigKey.Translater }, !isWeb());
    if (!resp || !resp.value) return;

    form.setValues(resp.value);
  };

  useEffect(() => {
    if (item === translateKey) fetchTranslateConfig();
    else form.reset();
  }, [item]);

  return (
    <form
      onSubmit={form.onSubmit((values) => {
        saveSystemConfig({ key: SystemConfigKey.Translater, value: values }, !isWeb());
        toast.success("ok");
      })}
    >
      <Stack align="stretch">
        <Autocomplete
          withAsterisk
          data={translationPlatforms}
          label={t("translate.label.platform", { ns: "settings" })}
          key={form.key("platform")}
          {...form.getInputProps("platform")}
        />
        <TextInput
          withAsterisk
          label={t("translate.label.app_id", { ns: "settings" })}
          key={form.key("appId")}
          {...form.getInputProps("appId")}
        />

        <PasswordInput
          label={t("translate.label.app_secret", { ns: "settings" })}
          key={form.key("appSecret")}
          {...form.getInputProps("appSecret")}
        />
        <Button type="submit">{t("button.save")}</Button>
      </Stack>
    </form>
  );
}

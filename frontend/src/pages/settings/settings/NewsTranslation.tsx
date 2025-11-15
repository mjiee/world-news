import { useEffect, useState } from "react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";
import { Button, Stack, TextInput, Autocomplete, PasswordInput } from "@mantine/core";
import { useForm } from "@mantine/form";
import { getSystemConfig, saveSystemConfig, SystemConfigKey, TranslaterConfig } from "@/services";
import { isWeb } from "@/utils/platform";

// News Translation Component
export default function NewsTranslation() {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);

  const form = useForm<TranslaterConfig>({
    initialValues: { platform: "", appId: "", appSecret: "" },
    validate: {
      platform: (value) =>
        !value
          ? t("translate.validate.required", {
              ns: "settings",
              label: t("translate.label.platform", { ns: "settings" }),
            })
          : null,
      appId: (value) =>
        !value
          ? t("translate.validate.required", { ns: "settings", label: t("translate.label.app_id", { ns: "settings" }) })
          : null,
    },
  });

  useEffect(() => {
    fetchConfig();
  }, []);

  const fetchConfig = async () => {
    setLoading(true);
    const resp = await getSystemConfig<TranslaterConfig>({ key: SystemConfigKey.Translater }, !isWeb());
    if (resp?.value) form.setValues(resp.value);
    setLoading(false);
  };

  return (
    <form
      onSubmit={form.onSubmit((values) => {
        saveSystemConfig({ key: SystemConfigKey.Translater, value: values }, !isWeb());
        toast.success("ok");
      })}
    >
      <Stack gap="md">
        <Autocomplete
          withAsterisk
          data={["google", "microsoft", "aliyun", "baidu"]}
          label={t("translate.label.platform", { ns: "settings" })}
          {...form.getInputProps("platform")}
        />
        <TextInput
          withAsterisk
          label={t("translate.label.app_id", { ns: "settings" })}
          {...form.getInputProps("appId")}
        />
        <PasswordInput
          label={t("translate.label.app_secret", { ns: "settings" })}
          {...form.getInputProps("appSecret")}
        />
        <Button type="submit" loading={loading}>
          {t("button.save")}
        </Button>
      </Stack>
    </form>
  );
}

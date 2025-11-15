import { getSystemConfig, saveSystemConfig, SystemConfigKey, TextAIConfig } from "@/services";
import { isWeb } from "@/utils/platform";
import { validateUrl } from "@/utils/url";
import { Autocomplete, Button, PasswordInput, Stack, TextInput } from "@mantine/core";
import { useForm } from "@mantine/form";
import { useEffect } from "react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";

// Text AI Component
export default function TextAi() {
  const { t } = useTranslation();

  const defaultConfigs: TextAIConfig[] = [
    { platform: "ChatGPT-GPT-4o", model: "gpt-4o", apiKey: "", apiUrl: "https://api.openai.com/v1/chat/completions" },
    {
      platform: "DeepSeek-chat",
      model: "deepseek-chat",
      apiKey: "",
      apiUrl: "https://api.deepseek.com/chat/completions",
    },
    {
      platform: "Aliyun-DeepSeek-R1",
      model: "deepseek-r1",
      apiKey: "",
      apiUrl: "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
    },
    {
      platform: "Aliyun-Qwen-Plus",
      model: "qwen-plus",
      apiKey: "",
      apiUrl: "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
    },
  ];

  const form = useForm<TextAIConfig>({
    initialValues: { platform: "", model: "", apiUrl: "", apiKey: "" },
    onValuesChange: (values) => {
      const config = defaultConfigs.find((item) => item.platform === values.platform);
      if (config && (config.model !== values.model || config.apiUrl !== values.apiUrl)) {
        form.setValues(config);
      }
    },
    validate: {
      platform: (value) =>
        !value
          ? t("text_ai.validate.required", {
              ns: "settings",
              label: t("text_ai.label.ai_platform", { ns: "settings" }),
            })
          : null,
      model: (value) =>
        !value
          ? t("text_ai.validate.required", { ns: "settings", label: t("text_ai.label.model", { ns: "settings" }) })
          : null,
      apiUrl: (value) => (!validateUrl(value) ? t("text_ai.validate.invalid_url", { ns: "settings" }) : null),
      apiKey: (value) =>
        !value
          ? t("text_ai.validate.required", { ns: "settings", label: t("text_ai.label.api_key", { ns: "settings" }) })
          : null,
    },
  });

  useEffect(() => {
    fetchConfig();
  }, []);

  const fetchConfig = async () => {
    const resp = await getSystemConfig<TextAIConfig>({ key: SystemConfigKey.TextAi }, !isWeb());
    if (resp?.value) form.setValues(resp.value);
  };

  return (
    <form
      onSubmit={form.onSubmit((values) => {
        saveSystemConfig({ key: SystemConfigKey.TextAi, value: values }, !isWeb());
        toast.success("ok");
      })}
    >
      <Stack gap="md">
        <Autocomplete
          withAsterisk
          data={defaultConfigs.map((item) => item.platform)}
          label={t("text_ai.label.ai_platform", { ns: "settings" })}
          {...form.getInputProps("platform")}
        />
        <TextInput withAsterisk label={t("text_ai.label.model", { ns: "settings" })} {...form.getInputProps("model")} />
        <TextInput
          withAsterisk
          label={t("text_ai.label.api_url", { ns: "settings" })}
          {...form.getInputProps("apiUrl")}
        />
        <PasswordInput
          withAsterisk
          label={t("text_ai.label.api_key", { ns: "settings" })}
          {...form.getInputProps("apiKey")}
        />
        <Button type="submit">{t("button.save")}</Button>
      </Stack>
    </form>
  );
}

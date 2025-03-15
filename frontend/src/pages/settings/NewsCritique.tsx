import { useEffect } from "react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";
import { Button, Stack, Textarea, TextInput, Autocomplete, PasswordInput, NumberInput } from "@mantine/core";
import { useForm } from "@mantine/form";
import { getSystemConfig, saveSystemConfig, SystemConfigKey, OpenAIConfig } from "@/services";
import { validateUrl } from "@/utils/url";

const defaultAiConfigs: OpenAIConfig[] = [
  { description: "ChatGPT-GPT-4o", model: "gpt-4o", apiKey: "", apiUrl: "https://api.openai.com/v1/chat/completions" },
  {
    description: "DeepSeek-chat",
    model: "deepseek-chat",
    apiKey: "",
    apiUrl: "https://api.deepseek.com/chat/completions",
  },
  {
    description: "Aliyun-DeepSeek-R1",
    model: "deepseek-r1",
    apiKey: "",
    apiUrl: "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
  },
  {
    description: "Aliyun-Qwen-Plus",
    model: "qwen-plus",
    apiKey: "",
    apiUrl: "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
  },
];

// news critique
export function NewsCritique() {
  const { t } = useTranslation();

  const aiform = useForm<OpenAIConfig>({
    mode: "uncontrolled",
    initialValues: {
      description: "",
      model: "",
      apiUrl: "",
      apiKey: "",
      maxTokens: 4096,
      systemPrompt: t("news_critique.default_system_prompt", { ns: "settings" }),
      assistantPrompt: "",
    },
    onValuesChange: (values) => {
      const config = defaultAiConfigs.find((item) => item.description === values.description);

      if (config && (config.model !== values.model || config.apiUrl !== values.apiUrl)) {
        aiform.setValues(config);
      }
    },
    validate: {
      description: (value) => validatRequiredField(value, "news_critique.label.ai_platform"),
      model: (value) => validatRequiredField(value, "news_critique.label.model"),
      apiUrl: (value) => (validateUrl(value) ? null : t("news_critique.validate.invalid_url", { ns: "settings" })),
      apiKey: (value) => validatRequiredField(value, "news_critique.label.api_key"),
    },
  });

  // validate required field
  function validatRequiredField(value: string, label: string): null | string {
    if (value) return null;
    return t("news_critique.validate.required", { ns: "settings", label: t(label, { ns: "settings" }) });
  }

  // fetch ai config
  const fetchAiConfig = async () => {
    const resp = await getSystemConfig<OpenAIConfig>({ key: SystemConfigKey.OpenAI });
    if (!resp || !resp.value) return;

    aiform.setValues(resp.value);
  };

  useEffect(() => {
    fetchAiConfig();
  }, []);

  return (
    <form
      onSubmit={aiform.onSubmit((values) => {
        saveSystemConfig({ key: SystemConfigKey.OpenAI, value: values });

        toast.success("ok");
      })}
    >
      <Stack align="stretch">
        <Autocomplete
          withAsterisk
          data={defaultAiConfigs.map((item) => item.description)}
          label={t("news_critique.label.ai_platform", { ns: "settings" })}
          key={aiform.key("description")}
          {...aiform.getInputProps("description")}
        />
        <TextInput
          withAsterisk
          label={t("news_critique.label.model", { ns: "settings" })}
          key={aiform.key("model")}
          {...aiform.getInputProps("model")}
        />
        <TextInput
          withAsterisk
          label={t("news_critique.label.api_url", { ns: "settings" })}
          key={aiform.key("apiUrl")}
          {...aiform.getInputProps("apiUrl")}
        />
        <PasswordInput
          withAsterisk
          label={t("news_critique.label.api_key", { ns: "settings" })}
          key={aiform.key("apiKey")}
          {...aiform.getInputProps("apiKey")}
        />
        <NumberInput
          label={t("news_critique.label.max_toknes", { ns: "settings" })}
          key={aiform.key("maxTokens")}
          {...aiform.getInputProps("maxTokens")}
        />
        <Textarea
          resize="vertical"
          label={t("news_critique.label.system_prompt", { ns: "settings" })}
          key={aiform.key("systemPrompt")}
          {...aiform.getInputProps("systemPrompt")}
        />
        <Textarea
          resize="vertical"
          label={t("news_critique.label.assistant_prompt", { ns: "settings" })}
          key={aiform.key("assistantPrompt")}
          {...aiform.getInputProps("assistantPrompt")}
        />
        <Button type="submit">{t("button.save")}</Button>
      </Stack>
    </form>
  );
}

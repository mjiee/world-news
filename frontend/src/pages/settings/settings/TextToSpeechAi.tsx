import { getSystemConfig, saveSystemConfig, SystemConfigKey, TextToSpeechAIConfig } from "@/services";
import { isWeb } from "@/utils/platform";
import {
  ActionIcon,
  Autocomplete,
  Box,
  Button,
  Group,
  PasswordInput,
  Stack,
  Switch,
  Text,
  TextInput,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { IconPlus, IconTrash } from "@tabler/icons-react";
import { useEffect } from "react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";
import styles from "../styles/settings.module.css";

// Text to Speech AI Component
export default function TextToSpeechAi() {
  const { t } = useTranslation();

  const form = useForm<TextToSpeechAIConfig>({
    initialValues: { platform: "", model: "", apiKey: "", voices: [] },
    validate: {
      platform: (value) =>
        !value
          ? t("tts_ai.validate.required", { ns: "settings", label: t("tts_ai.label.ai_platform", { ns: "settings" }) })
          : null,
      model: (value) =>
        !value
          ? t("tts_ai.validate.required", { ns: "settings", label: t("tts_ai.label.model", { ns: "settings" }) })
          : null,
      apiKey: (value) =>
        !value
          ? t("tts_ai.validate.required", { ns: "settings", label: t("tts_ai.label.api_key", { ns: "settings" }) })
          : null,
    },
  });

  useEffect(() => {
    fetchConfig();
  }, []);

  const fetchConfig = async () => {
    const resp = await getSystemConfig<TextToSpeechAIConfig>({ key: SystemConfigKey.TextToSpeechAi }, !isWeb());
    if (resp?.value) form.setValues({ ...resp.value, voices: resp.value.voices || [] });
  };

  const addVoice = () => {
    const current = form.getValues().voices || [];
    form.setFieldValue("voices", [...current, { id: "", name: "", description: "" }]);
  };

  const removeVoice = (idx: number) => {
    const current = form.getValues().voices || [];
    form.setFieldValue(
      "voices",
      current.filter((_, i) => i !== idx),
    );
  };

  return (
    <form
      onSubmit={form.onSubmit((values) => {
        saveSystemConfig({ key: SystemConfigKey.TextToSpeechAi, value: values }, !isWeb());
        toast.success("ok");
      })}
    >
      <Stack gap="md">
        <Autocomplete
          withAsterisk
          data={["doubao"]}
          label={t("tts_ai.label.ai_platform", { ns: "settings" })}
          {...form.getInputProps("platform")}
        />
        <TextInput withAsterisk label={t("tts_ai.label.model", { ns: "settings" })} {...form.getInputProps("model")} />
        <PasswordInput
          withAsterisk
          label={t("tts_ai.label.api_key", { ns: "settings" })}
          {...form.getInputProps("apiKey")}
        />
        <Switch label={t("tts_ai.label.auto_task", { ns: "settings" })} {...form.getInputProps("autoTask")} />

        <Box>
          <Group justify="space-between" mb="xs">
            <Text size="sm" fw={500}>
              {t("tts_ai.label.voices", { ns: "settings" })}
            </Text>
            <ActionIcon variant="light" onClick={addVoice}>
              <IconPlus size={16} />
            </ActionIcon>
          </Group>

          <Stack gap="xs">
            {(form.getValues().voices || []).map((_, idx) => (
              <Group key={idx} align="center" wrap="nowrap" className={styles.voiceRow}>
                <TextInput
                  placeholder={t("tts_ai.label.voice_id", { ns: "settings" })}
                  style={{ flex: 1 }}
                  {...form.getInputProps(`voices.${idx}.id`)}
                />
                <TextInput
                  placeholder={t("tts_ai.label.voice_name", { ns: "settings" })}
                  style={{ flex: 1 }}
                  {...form.getInputProps(`voices.${idx}.name`)}
                />
                <TextInput
                  placeholder={t("tts_ai.label.voice_model", { ns: "settings" })}
                  style={{ flex: 1 }}
                  {...form.getInputProps(`voices.${idx}.model`)}
                />
                <TextInput
                  placeholder={t("tts_ai.label.voice_description", { ns: "settings" })}
                  style={{ flex: 2 }}
                  {...form.getInputProps(`voices.${idx}.description`)}
                />
                <ActionIcon color="red" variant="light" onClick={() => removeVoice(idx)}>
                  <IconTrash size={16} />
                </ActionIcon>
              </Group>
            ))}
          </Stack>
        </Box>

        <Button type="submit">{t("button.save")}</Button>
      </Stack>
    </form>
  );
}

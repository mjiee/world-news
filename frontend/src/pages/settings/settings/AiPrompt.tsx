import { useEffect } from "react";
import { Text, Stack, Textarea, Button, TextInput, ActionIcon, Paper, Group } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { IconPlus, IconTrash } from "@tabler/icons-react";
import {
  getSystemConfig,
  NewsCritiquePrompt,
  PodcastScriptPrompt,
  saveSystemConfig,
  SystemConfigKey,
} from "@/services/systemConfigApi";
import { useForm } from "@mantine/form";
import toast from "react-hot-toast";
import { isWeb } from "@/utils/platform";
import styles from "../styles/settings.module.css";

// Critique Prompt Component
export function CritiquePrompt() {
  const { t } = useTranslation();

  const form = useForm<NewsCritiquePrompt>({
    initialValues: {
      systemPrompt: t("critique_prompt.default", { ns: "settings" }),
    },
    validate: {
      systemPrompt: (value) => (!value ? t("critique_prompt.validate", { ns: "settings" }) : null),
    },
  });

  useEffect(() => {
    fetchConfig();
  }, []);

  const fetchConfig = async () => {
    const resp = await getSystemConfig<NewsCritiquePrompt>({ key: SystemConfigKey.NewsCritiquePromptKey }, !isWeb());
    if (resp?.value) form.setValues(resp.value);
  };

  return (
    <form
      onSubmit={form.onSubmit((values) => {
        saveSystemConfig({ key: SystemConfigKey.NewsCritiquePromptKey, value: values }, !isWeb());
        toast.success("ok");
      })}
    >
      <Stack gap="md">
        <Textarea
          resize="vertical"
          label={t("critique_prompt.label", { ns: "settings" })}
          minRows={4}
          {...form.getInputProps("systemPrompt")}
        />
        <Button type="submit">{t("button.save")}</Button>
      </Stack>
    </form>
  );
}

// Podcast Prompt Component
export function PodcastPrompt() {
  const { t } = useTranslation();

  const form = useForm<PodcastScriptPrompt>({
    initialValues: {
      systemPrompt: t("podcast_prompt.default.system", { ns: "settings" }),
      stylizePrompts: [],
    },
    validate: {
      systemPrompt: (value) =>
        !value
          ? t("podcast_prompt.validate", {
              ns: "settings",
              label: t("podcast_prompt.label.system", { ns: "settings" }),
            })
          : null,
    },
  });

  useEffect(() => {
    fetchConfig();
  }, []);

  const fetchConfig = async () => {
    const resp = await getSystemConfig<PodcastScriptPrompt>({ key: SystemConfigKey.PodcastScriptPromptKey }, !isWeb());
    if (resp?.value) form.setValues(resp.value);
  };

  const addStyle = () => {
    const current = form.getValues().stylizePrompts || [];
    form.setFieldValue("stylizePrompts", [...current, { style: "", prompt: "" }]);
  };

  const removeStyle = (idx: number) => {
    const current = form.getValues().stylizePrompts || [];
    form.setFieldValue(
      "stylizePrompts",
      current.filter((_, i) => i !== idx),
    );
  };

  return (
    <form
      onSubmit={form.onSubmit((values) => {
        saveSystemConfig({ key: SystemConfigKey.PodcastScriptPromptKey, value: values }, !isWeb());
        toast.success("ok");
      })}
    >
      <Stack gap="md">
        <Textarea
          resize="vertical"
          label={t("podcast_prompt.label.system", { ns: "settings" })}
          minRows={3}
          {...form.getInputProps("systemPrompt")}
        />
        <Textarea
          resize="vertical"
          label={t("podcast_prompt.label.approval", { ns: "settings" })}
          placeholder={t("podcast_prompt.default.approval", { ns: "settings" })}
          minRows={2}
          {...form.getInputProps("approvalPrompt")}
        />
        <Textarea
          resize="vertical"
          label={t("podcast_prompt.label.rewrite", { ns: "settings" })}
          placeholder={t("podcast_prompt.default.rewrite", { ns: "settings" })}
          minRows={2}
          {...form.getInputProps("rewritePrompt")}
        />
        <Textarea
          resize="vertical"
          label={t("podcast_prompt.label.merge", { ns: "settings" })}
          placeholder={t("podcast_prompt.default.merge", { ns: "settings" })}
          minRows={2}
          {...form.getInputProps("mergePrompt")}
        />
        <Textarea
          resize="vertical"
          label={t("podcast_prompt.label.classify", { ns: "settings" })}
          placeholder={t("podcast_prompt.default.classify", { ns: "settings" })}
          minRows={2}
          {...form.getInputProps("classifyPrompt")}
        />

        <Group gap="sm">
          <ActionIcon variant="filled" onClick={addStyle} size="lg">
            <IconPlus />
          </ActionIcon>
          <Text size="sm" c="dimmed">
            {t("podcast_prompt.label.add_style", { ns: "settings" })}
          </Text>
        </Group>

        {(form.getValues().stylizePrompts || []).map((item, idx) => (
          <Paper key={idx} p="md" withBorder className={styles.styleItem}>
            <Group align="flex-start" gap="md" wrap="nowrap">
              <Stack style={{ flex: 1 }} gap="sm">
                <TextInput
                  label={t("podcast_prompt.label.style", { ns: "settings" })}
                  withAsterisk
                  {...form.getInputProps(`stylizePrompts.${idx}.style`)}
                />
                <Textarea
                  label={t("podcast_prompt.label.podcast", { ns: "settings" })}
                  withAsterisk
                  resize="vertical"
                  minRows={2}
                  {...form.getInputProps(`stylizePrompts.${idx}.prompt`)}
                />
              </Stack>
              <ActionIcon variant="light" color="red" onClick={() => removeStyle(idx)} size="lg" mt={28}>
                <IconTrash size={18} />
              </ActionIcon>
            </Group>
          </Paper>
        ))}

        <Button type="submit">{t("button.save")}</Button>
      </Stack>
    </form>
  );
}

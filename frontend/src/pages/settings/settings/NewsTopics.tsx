import { useState, useEffect } from "react";
import toast from "react-hot-toast";
import { Button, PillsInput, Pill, Stack } from "@mantine/core";
import { useField } from "@mantine/form";
import { useTranslation } from "react-i18next";
import { getSystemConfig, saveSystemConfig, SystemConfigKey } from "@/services";
import styles from "../styles/settings.module.css";

// News Topics Component
export default function NewsTopics() {
  const { t } = useTranslation("settings");
  const [topics, setTopics] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);

  const field = useField({
    initialValue: "",
    validateOnChange: true,
    validate: (value) => {
      if (value.length > 0 && value.trim().length === 0) return t("news_topic.topic_validate.required");
      if (topics.includes(value.trim())) return t("news_topic.topic_validate.duplicate");
      return null;
    },
  });

  useEffect(() => {
    fetchNewsTopics();
  }, []);

  const fetchNewsTopics = async () => {
    setLoading(true);
    const resp = await getSystemConfig<string[]>({ key: SystemConfigKey.NewsTopics });
    if (resp?.value && resp?.value?.length > 0) setTopics(resp.value);
    setLoading(false);
  };

  const saveTopics = async (newTopics: string[]) => {
    await saveSystemConfig({ key: SystemConfigKey.NewsTopics, value: newTopics });
    toast.success("ok");
  };

  const addTopic = () => {
    const value = field.getValue().trim();
    if (!value || topics.includes(value)) return;

    const updated = [...topics, value];
    setTopics(updated);
    saveTopics(updated);
    field.setValue("");
  };

  const removeTopic = (value: string) => {
    const updated = topics.filter((t) => t !== value);
    setTopics(updated);
    saveTopics(updated);
  };

  return (
    <Stack gap="md">
      <PillsInput {...field.getInputProps()}>
        <Pill.Group>
          {topics.map((value) => (
            <Pill key={value} withRemoveButton size="lg" onRemove={() => removeTopic(value)} className={styles.pill}>
              {value}
            </Pill>
          ))}
          <PillsInput.Field
            value={field.getValue()}
            onChange={(e) => field.setValue(e.currentTarget.value)}
            placeholder={t("news_topic.input_placeholder")}
          />
        </Pill.Group>
      </PillsInput>
      <Button onClick={addTopic} loading={loading}>
        {t("news_topic.add_button")}
      </Button>
    </Stack>
  );
}

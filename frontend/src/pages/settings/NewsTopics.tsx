import { useState, useEffect } from "react";
import { Button, PillsInput, Pill, Stack, TextInput } from "@mantine/core";
import { useField } from "@mantine/form";
import { useTranslation } from "react-i18next";
import { getSystemConfig, saveSystemConfig } from "@/services";

// news topic settings
const newsTopicKey = "newsTopic";

export function NewsTopics() {
  const { t } = useTranslation("settings");
  const [topics, setTopics] = useState<string[]>([]);

  const field = useField({
    initialValue: "",
    validateOnChange: true,
    validate: (value) => {
      if (value.length > 0 && value.trim().length == 0) return t("news_topic.topic_validate.required");

      if (topics.includes(value.trim())) return t("news_topic.topic_validate.duplicate");

      return null;
    },
  });

  // save news topics
  const saveNewsTopics = async () => {
    await saveSystemConfig({ key: newsTopicKey, value: topics });
  };

  // fetch news topics
  const fetchNewsTopics = async () => {
    const resp = await getSystemConfig<string[]>({ key: newsTopicKey });

    if (!resp || !resp.value) return;
    if (resp.value.length === 0) return;

    setTopics(resp.value);
  };

  useEffect(() => {
    fetchNewsTopics();
  }, []);

  // add news topic
  const addNewsTopicHandle = () => {
    const value = field.getValue().trim();

    if (value.length === 0) return;

    if (topics.includes(value)) {
      field.setValue("");

      return;
    }

    setTopics([...topics, value]);
    saveNewsTopics();
    field.setValue("");
  };

  // remove news topic
  const removeNewsTopicHandle = (value: string) => {
    setTopics(topics.filter((topic) => topic !== value));
    saveNewsTopics();
  };

  return (
    <Stack w={"100%"} align="stretch" justify="flex-start" gap="md">
      <PillsInput {...field.getInputProps()} variant="unstyled">
        <Pill.Group>
          {topics.map((value) => (
            <Pill key={value} withRemoveButton size="lg" onRemove={() => removeNewsTopicHandle(value)}>
              {value}
            </Pill>
          ))}
          <PillsInput.Field
            value={field.getValue()}
            onChange={(event) => field.setValue(event.currentTarget.value)}
            placeholder={t("news_topic.input_placeholder")}
          />
        </Pill.Group>
      </PillsInput>
      <Button variant="default" onClick={addNewsTopicHandle}>
        {t("news_topic.add_button")}
      </Button>
    </Stack>
  );
}

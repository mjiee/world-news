import {
  createAudio,
  editScript,
  getSystemConfig,
  PodcastScript,
  SystemConfigKey,
  TaskStage,
  TaskStageStatus,
  TextToSpeechAIConfig,
} from "@/services";
import {
  ActionIcon,
  Alert,
  Badge,
  Box,
  Button,
  Card,
  Group,
  Modal,
  NumberInput,
  Select,
  Stack,
  Text,
  Textarea,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { IconAlertCircle, IconEdit, IconPlayerPlay, IconTrash } from "@tabler/icons-react";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import classes from "../styles/taskDetail.module.css";

const emotions = [
  "happy",
  "sad",
  "angry",
  "surprised",
  "fear",
  "hate",
  "excited",
  "coldness",
  "neutral",
  "depressed",
  "lovey-dovey",
  "shy",
  "comfort",
  "tension",
  "tender",
  "storytelling",
  "radio",
  "magnetic",
  "advertising",
  "vocal - fry",
  "ASMR",
  "news",
  "entertainment",
];

export default function StageScriptedCard({ stage, onRefresh }: { stage: TaskStage; onRefresh: () => void }) {
  const { t } = useTranslation();

  if (stage.status === TaskStageStatus.Failed) {
    return (
      <Alert icon={<IconAlertCircle />} color="red" radius="md" mb={20}>
        <Text size="sm" fw={500} mb={4}>
          {t("error_reason", { ns: "task" })}:
        </Text>
        <Text size="sm">{stage.reason || "Unknown Error"}</Text>
      </Alert>
    );
  }

  const [loading, setLoading] = useState(false);
  const [editingIndex, setEditingIndex] = useState<number | null>(null);
  const [scripts, setScripts] = useState<PodcastScript[]>(stage?.audio?.scripts || []);
  const [voices, setVoices] = useState<{ value: string; label: string }[]>([]);
  const editForm = useForm({
    mode: "uncontrolled",
    initialValues: { content: "", speaker: "", emotion: "", speechRate: 0, volume: 50 },
  });

  const createPodcastAudio = async () => {
    setLoading(true);
    await createAudio(stage.id);
    setLoading(false);
    onRefresh();
  };

  const openEditModal = (index: number, script: PodcastScript) => {
    setEditingIndex(index);
    editForm.setValues(script);
  };

  const saveScripts = async (scripts: PodcastScript[]) => {
    if (scripts.length === 0) return;
    setLoading(true);
    await editScript(stage.id, scripts);
    setLoading(false);
    setScripts(scripts);
  };

  const deleteScript = async (index: number) => {
    await saveScripts(scripts.filter((_, i) => i !== index));
  };

  const saveEdit = async () => {
    if (editingIndex === null || !editForm) return;
    await saveScripts(
      scripts.map((script, i) => (i === editingIndex ? { ...script, ...editForm.getValues() } : script)),
    );
    setEditingIndex(null);
    editForm.setValues({});
  };

  const loadVoices = async () => {
    const resp = await getSystemConfig<TextToSpeechAIConfig>({ key: SystemConfigKey.TextToSpeechAi }, true);

    if (!resp || !resp.value || !resp.value.voices) return [];

    setVoices(resp.value.voices.map((v) => ({ value: v.id, label: v.name })));
  };

  useEffect(() => {
    loadVoices();
  }, [stage]);

  const scriptLabel = (label: string, color: string, value: string | number) => {
    return (
      <Box>
        <Text size="xs" c="dimmed">
          {label}
        </Text>
        <Badge variant="light" color={color}>
          {value}
        </Badge>
      </Box>
    );
  };

  return (
    <>
      <Stack gap="md" mb="xl">
        {scripts.map((script, index) => (
          <Card key={index} shadow="sm" padding="lg" radius="md" withBorder className={classes.scriptcard}>
            <Stack gap="sm">
              <Group justify="space-between" align="flex-start">
                <Badge size="lg" variant="gradient" gradient={{ from: "violet", to: "grape" }}>
                  # {index + 1}
                </Badge>
                <Group gap="xs">
                  <ActionIcon variant="light" color="blue" onClick={() => openEditModal(index, script)}>
                    <IconEdit size={16} />
                  </ActionIcon>
                  <ActionIcon variant="light" color="red" onClick={() => deleteScript(index)}>
                    <IconTrash size={16} />
                  </ActionIcon>
                </Group>
              </Group>

              <Box>
                <Text size="sm" fw={500} c="dimmed" mb={4}>
                  {t("podcast.scripted.content", { ns: "task" })}:
                </Text>
                <Text size="sm">{script.content}</Text>
              </Box>

              <Group gap="md">
                {scriptLabel(
                  t("podcast.scripted.voice", { ns: "task" }),
                  "cyan",
                  voices.find((v) => v.value === script.speaker)?.label ?? script.speaker,
                )}
                {scriptLabel(t("podcast.scripted.emotion", { ns: "task" }), "pink", script.emotion)}
                {scriptLabel(t("podcast.scripted.speed", { ns: "task" }), "grape", script.speechRate)}
                {scriptLabel(t("podcast.scripted.volume", { ns: "task" }), "indigo", script.volume)}
              </Group>
            </Stack>
          </Card>
        ))}
      </Stack>

      {stage.status === TaskStageStatus.Completed && (
        <Button
          size="lg"
          variant="gradient"
          gradient={{ from: "violet", to: "grape" }}
          leftSection={<IconPlayerPlay size={18} />}
          onClick={createPodcastAudio}
          loading={loading}
          fullWidth
        >
          {t("podcast.scripted.generate", { ns: "task" })}
        </Button>
      )}

      <Modal opened={editingIndex !== null} onClose={() => setEditingIndex(null)} title="编辑脚本" size="lg">
        <Stack gap="md">
          <Textarea label={t("podcast.scripted.content", { ns: "task" })} {...editForm.getInputProps("content")} />
          <Select
            label={t("podcast.scripted.voice", { ns: "task" })}
            data={voices}
            {...editForm.getInputProps("speaker")}
          />
          <Select
            label={t("podcast.scripted.emotion", { ns: "task" })}
            data={emotions}
            {...editForm.getInputProps("emotion")}
          />
          <NumberInput
            label={t("podcast.scripted.speed", { ns: "task" })}
            min={0}
            max={2}
            step={0.1}
            {...editForm.getInputProps("speechRate")}
          />
          <NumberInput
            label={t("podcast.scripted.volume", { ns: "task" })}
            min={0}
            max={100}
            {...editForm.getInputProps("volume")}
          />

          <Group justify="flex-end" mt="md">
            <Button variant="light" onClick={() => setEditingIndex(null)}>
              {t("button.cancel")}
            </Button>
            <Button variant="gradient" gradient={{ from: "violet", to: "grape" }} loading={loading} onClick={saveEdit}>
              {t("button.save")}
            </Button>
          </Group>
        </Stack>
      </Modal>
    </>
  );
}

import {
  createAudio,
  editScript,
  getSystemConfig,
  PodcastScript,
  SystemConfigKey,
  TaskStage,
  TaskStageStatus,
  textToSpeech,
  TextToSpeechAIConfig,
} from "@/services";
import { buildAudioSrc } from "@/stores";
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
import { IconAlertCircle, IconEdit, IconPlayerPlay, IconPlus, IconRefresh, IconTrash } from "@tabler/icons-react";
import { useEffect, useRef, useState } from "react";
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
  "news",
  "entertainment",
];

const DEFAULT_FORM_VALUES = {
  text: "",
  speaker: "",
  emotion: "",
  speed: 0,
  volume: 50,
  silence: 0.2,
};

type ModalState = { mode: "edit" | "add"; index: number } | null;

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
  const [modalState, setModalState] = useState<ModalState>(null);
  const [playingIndex, setPlayingIndex] = useState<number | null>(null);
  const audioRef = useRef<HTMLAudioElement | null>(null);

  const [scripts, setScripts] = useState<PodcastScript[]>(stage?.audio?.scripts || []);
  const [voices, setVoices] = useState<{ value: string; label: string }[]>([]);

  const form = useForm({ mode: "uncontrolled", initialValues: DEFAULT_FORM_VALUES });

  const openModal = (state: ModalState) => {
    if (!state) return;
    form.setValues(scripts[state.index]);
    setModalState(state);
  };

  const closeModal = () => {
    setModalState(null);
    form.reset();
  };

  const createPodcastAudio = async () => {
    setLoading(true);
    await createAudio(stage.id);
    setLoading(false);
    onRefresh();
  };

  const reGenerateScriptAudio = async (script: PodcastScript, index: number) => {
    setLoading(true);
    const resp = await textToSpeech(script);
    if (resp) {
      script.audio = resp;
      await editScript(
        stage.id,
        scripts.map((s, i) => (i === index ? { ...script } : s)),
      );
    }
    setLoading(false);
  };

  const togglePlayAudio = (index: number, format: string, audioData: string) => {
    if (playingIndex === index) {
      audioRef.current?.pause();
      setPlayingIndex(null);
      return;
    }
    audioRef.current?.pause();
    const audio = new Audio(buildAudioSrc(format, audioData));
    audioRef.current = audio;
    audio.play();
    setPlayingIndex(index);
    audio.onended = () => setPlayingIndex(null);
  };

  const saveScripts = async (updated: PodcastScript[]) => {
    if (updated.length === 0) return;
    setLoading(true);
    await editScript(stage.id, updated);
    setLoading(false);
    setScripts(updated);
  };

  const deleteScript = (index: number) => saveScripts(scripts.filter((_, i) => i !== index));

  const handleSave = async () => {
    if (!modalState) return;
    const values = form.getValues();

    if (modalState.mode === "edit") {
      await saveScripts(scripts.map((s, i) => (i === modalState.index ? { ...s, ...values } : s)));
    } else {
      const newScript: PodcastScript = {
        format: scripts[modalState.index]?.format ?? "",
        ...values,
      };
      await saveScripts([...scripts.slice(0, modalState.index + 1), newScript, ...scripts.slice(modalState.index + 1)]);
    }

    closeModal();
  };

  const loadVoices = async () => {
    const resp = await getSystemConfig<TextToSpeechAIConfig>({ key: SystemConfigKey.TextToSpeechAi }, true);
    if (!resp?.value?.voices) return;
    setVoices(resp.value.voices.map((v) => ({ value: v.id, label: v.name })));
  };

  useEffect(() => {
    loadVoices();
  }, [stage]);

  const scriptLabel = (label: string, color: string, value: string | number) => (
    <Box>
      <Text size="xs" c="dimmed">
        {label}
      </Text>
      <Badge variant="light" color={color}>
        {value}
      </Badge>
    </Box>
  );

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
                  {script.audio && (
                    <ActionIcon
                      variant="light"
                      color={loading ? "gray" : playingIndex === index ? "orange" : "green"}
                      disabled={loading}
                      onClick={() => togglePlayAudio(index, script.format, script.audio!)}
                    >
                      <IconPlayerPlay size={16} />
                    </ActionIcon>
                  )}
                  <ActionIcon
                    variant="light"
                    color={loading ? "gray" : "yellow"}
                    disabled={loading}
                    onClick={() => reGenerateScriptAudio(script, index)}
                  >
                    <IconRefresh size={16} />
                  </ActionIcon>
                  <ActionIcon variant="light" color="blue" onClick={() => openModal({ mode: "edit", index })}>
                    <IconEdit size={16} />
                  </ActionIcon>
                  <ActionIcon variant="light" color="teal" onClick={() => openModal({ mode: "add", index })}>
                    <IconPlus size={16} />
                  </ActionIcon>
                  <ActionIcon variant="light" color="red" onClick={() => deleteScript(index)}>
                    <IconTrash size={16} />
                  </ActionIcon>
                </Group>
              </Group>

              <Box>
                <Text size="sm" fw={500} c="dimmed" mb={4}>
                  {t("podcast.scripted.text", { ns: "task" })}:
                </Text>
                <Text size="sm">{script.text}</Text>
              </Box>

              <Group gap="md">
                {scriptLabel(
                  t("podcast.scripted.voice", { ns: "task" }),
                  "cyan",
                  voices.find((v) => v.value === script.speaker)?.label ?? script.speaker,
                )}
                {scriptLabel(t("podcast.scripted.emotion", { ns: "task" }), "pink", script.emotion)}
                {scriptLabel(t("podcast.scripted.speed", { ns: "task" }), "grape", script.speed)}
                {scriptLabel(t("podcast.scripted.volume", { ns: "task" }), "indigo", script.volume)}
                {scriptLabel(t("podcast.scripted.silence", { ns: "task" }), "cyan", script.silence)}
              </Group>
            </Stack>
          </Card>
        ))}
      </Stack>

      {stage.status === TaskStageStatus.Completed && (
        <Group grow>
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
        </Group>
      )}

      <Modal
        opened={modalState !== null}
        onClose={closeModal}
        title={t(modalState?.mode === "edit" ? "podcast.scripted.edit" : "podcast.scripted.add", { ns: "task" })}
        size="lg"
      >
        <Stack gap="md">
          <Textarea label={t("podcast.scripted.text", { ns: "task" })} {...form.getInputProps("text")} />
          <Select
            label={t("podcast.scripted.voice", { ns: "task" })}
            data={voices}
            {...form.getInputProps("speaker")}
          />
          <Select
            label={t("podcast.scripted.emotion", { ns: "task" })}
            data={emotions}
            {...form.getInputProps("emotion")}
          />
          <NumberInput
            label={t("podcast.scripted.speed", { ns: "task" })}
            min={0}
            max={2}
            step={0.1}
            {...form.getInputProps("speed")}
          />
          <NumberInput
            label={t("podcast.scripted.volume", { ns: "task" })}
            min={0}
            max={100}
            {...form.getInputProps("volume")}
          />
          <NumberInput
            label={t("podcast.scripted.silence", { ns: "task" })}
            min={0}
            max={4}
            step={0.1}
            {...form.getInputProps("silence")}
          />

          <Group justify="flex-end" mt="md">
            <Button variant="light" onClick={closeModal}>
              {t("button.cancel")}
            </Button>
            <Button
              variant="gradient"
              gradient={{ from: "violet", to: "grape" }}
              loading={loading}
              onClick={handleSave}
            >
              {t("button.save")}
            </Button>
          </Group>
        </Stack>
      </Modal>
    </>
  );
}

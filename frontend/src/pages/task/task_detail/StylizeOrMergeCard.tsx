import {
  createScript,
  getSystemConfig,
  restyleArticle,
  SystemConfigKey,
  TaskStage,
  TaskStageName,
  TaskStageStatus,
  TextToSpeechAIConfig,
} from "@/services";
import { useMergeArticleStore } from "@/stores";
import { md } from "@/utils/md";
import { Alert, Box, Button, Card, Group, Modal, MultiSelect, Stack, Text, Textarea } from "@mantine/core";
import { useField } from "@mantine/form";
import { IconAlertCircle } from "@tabler/icons-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import classes from "../styles/taskDetail.module.css";

enum ActionType {
  Default = 0,
  CreateScript = 1,
  RestyleArticle = 2,
}

export default function StylizeOrMergeCard({ stage, onRefresh }: { stage: TaskStage; onRefresh: () => void }) {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [edit, setEdit] = useState(ActionType.Default);
  const [voices, setVoices] = useState<{ value: string; label: string }[]>([]);
  const { addStage } = useMergeArticleStore();

  const promptField = useField({ initialValue: stage.prompt });
  const voiceField = useField({ initialValue: [] });

  const createPodcastScript = async () => {
    setLoading(true);
    await createScript(stage.id, voiceField.getValue());
    setLoading(false);
    setEdit(ActionType.Default);
    voiceField.setValue([]);
    onRefresh();
  };

  const restylePodcastArticle = async () => {
    setLoading(true);
    await restyleArticle(stage.id, promptField.getValue());
    setLoading(false);
    setEdit(ActionType.Default);
    onRefresh();
  };

  const loadVoices = async () => {
    const resp = await getSystemConfig<TextToSpeechAIConfig>({ key: SystemConfigKey.TextToSpeechAi }, true);

    if (!resp || !resp.value || !resp.value.voices) return [];

    setVoices(resp.value.voices.map((v) => ({ value: v.id, label: v.name })));
  };

  const actionButton = (label: string, onClick: () => void) => {
    return (
      <Button variant="gradient" gradient={{ from: "violet", to: "grape" }} onClick={onClick} loading={loading}>
        {label}
      </Button>
    );
  };

  return (
    <>
      <Card shadow="sm" padding="lg" radius="md" withBorder mb={20}>
        {stage.prompt && stage.prompt.length > 0 && (
          <Box mb="md">
            <Text size="sm" fw={500} c="dimmed" mb="xs">
              {t("podcast.prompt", { ns: "task" })}:
            </Text>
            <Text size="sm" c="gray.7">
              {stage.prompt}
            </Text>
          </Box>
        )}

        {stage.status === TaskStageStatus.Failed && (
          <Alert icon={<IconAlertCircle />} color="red" radius="md" mb="md">
            <Text size="sm" fw={500} mb={4}>
              {t("error_reason", { ns: "task" })}:
            </Text>
            <Text size="sm">{stage.reason || "Unknown Error"}</Text>
          </Alert>
        )}

        {stage.output && stage.output.length > 0 && (
          <>
            <Box mb="md" className={classes.contentbox}>
              <div
                dangerouslySetInnerHTML={{
                  __html: md.render(stage.output),
                }}
                style={{ fontSize: "0.95rem", lineHeight: "1.6" }}
              />
            </Box>
            <Group justify="center" grow>
              {actionButton(t("podcast.stylize.scripted", { ns: "task" }), () => setEdit(ActionType.CreateScript))}
              {actionButton(t("podcast.stylize.rewrite", { ns: "task" }), () => setEdit(ActionType.RestyleArticle))}
              {TaskStageName.Merge != stage.stage &&
                actionButton(t("podcast.stylize.merge", { ns: "task" }), () => addStage(stage))}
            </Group>
          </>
        )}
      </Card>
      <Modal
        opened={edit != ActionType.Default}
        onClose={() => setEdit(ActionType.Default)}
        title={
          edit == ActionType.RestyleArticle
            ? t("podcast.stylize.rewrite", { ns: "task" })
            : t("podcast.stylize.scripted", { ns: "task" })
        }
        size="lg"
      >
        <Stack gap="md">
          {edit == ActionType.RestyleArticle && (
            <Textarea {...promptField.getInputProps()} label={t("podcast.prompt", { ns: "task" })} />
          )}
          {edit == ActionType.CreateScript && (
            <MultiSelect
              data={voices}
              {...voiceField.getInputProps()}
              label={t("podcast.voice", { ns: "task" })}
              onDropdownOpen={loadVoices}
            />
          )}
          <Group justify="flex-end" mt="md">
            <Button variant="light" onClick={() => setEdit(ActionType.Default)}>
              {t("button.cancel")}
            </Button>
            <Button
              variant="gradient"
              gradient={{ from: "violet", to: "grape" }}
              loading={loading}
              onClick={edit == ActionType.RestyleArticle ? restylePodcastArticle : createPodcastScript}
            >
              {t("button.save")}
            </Button>
          </Group>
        </Stack>
      </Modal>
    </>
  );
}

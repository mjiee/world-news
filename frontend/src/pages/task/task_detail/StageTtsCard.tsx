import { downloadAudio, TaskStage, TaskStageStatus } from "@/services";
import { useAudioPlayStore } from "@/stores";
import { ActionIcon, Alert, Badge, Box, Button, Group, Modal, Stack, Text, TextInput } from "@mantine/core";
import { useField } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { IconAlertCircle, IconCheck, IconDownload, IconPlaylistAdd } from "@tabler/icons-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

export default function StageTtsCard({ stage }: { stage: TaskStage }) {
  const { t } = useTranslation();
  const [hasDownload, sethasDownload] = useState<boolean>(false);
  const { addAudio, inPlayList } = useAudioPlayStore();
  const fileNameField = useField({ initialValue: "" });
  const [opened, { open, close }] = useDisclosure(false);

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

  if (stage.status === TaskStageStatus.Processing || !stage.audio) return null;

  const isInPlaylist = inPlayList(stage.id);

  const handleAddToPlaylist = () => {
    if (isInPlaylist) return;
    addAudio(stage.id, stage.audio!);
  };

  const onDownload = async () => {
    await downloadAudio(stage.id, fileNameField.getValue());
    sethasDownload(true);
    close();
  };

  return (
    <>
      <Stack gap="md">
        <Group justify="space-between" align="center">
          <Box>
            <Text size="lg" fw={600} mb={4}>
              {t("podcast.tts.sub_title", { ns: "task" })}
            </Text>
            {stage.audio.duration && (
              <Text size="sm" c="dimmed">
                {stage.audio.duration}
              </Text>
            )}
          </Box>
          <Group gap="xs">
            <ActionIcon
              variant="light"
              size="lg"
              radius="md"
              onClick={handleAddToPlaylist}
              color={isInPlaylist ? "green" : "blue"}
            >
              {isInPlaylist ? <IconCheck size={20} /> : <IconPlaylistAdd size={20} />}
            </ActionIcon>
            <ActionIcon variant="light" size="lg" radius="md" onClick={open} disabled={hasDownload}>
              {hasDownload ? <IconCheck size={20} /> : <IconDownload size={20} />}
            </ActionIcon>
          </Group>
        </Group>

        {stage.audio.voices && stage.audio.voices.length > 0 && (
          <Group gap="xs">
            {stage.audio.voices.map((voice, idx) => (
              <Badge key={idx} variant="light" size="md">
                {voice.name || voice.id}
              </Badge>
            ))}
          </Group>
        )}
      </Stack>

      <Modal opened={opened} onClose={close} title={t("podcast.tts.download", { ns: "task" })} size="lg">
        <Stack gap="md">
          <TextInput
            label={t("podcast.tts.file_name", { ns: "task" })}
            description={t("podcast.tts.download_desc", { ns: "task" })}
            {...fileNameField.getInputProps()}
          />
          <Button variant="gradient" gradient={{ from: "violet", to: "grape" }} onClick={onDownload}>
            {t("button.save")}
          </Button>
        </Stack>
      </Modal>
    </>
  );
}

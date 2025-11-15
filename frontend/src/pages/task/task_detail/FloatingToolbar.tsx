import { getSystemConfig, mergeArticle, MergeArticleRequest, SystemConfigKey, TextToSpeechAIConfig } from "@/services";
import { GolbalLanguage, useMergeArticleStore } from "@/stores";
import {
  ActionIcon,
  Affix,
  Badge,
  Box,
  Button,
  Divider,
  Group,
  Modal,
  MultiSelect,
  Paper,
  Stack,
  Text,
  TextInput,
  Transition,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { IconLayersDifference, IconRefresh, IconTrash } from "@tabler/icons-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router";
import styles from "../styles/floatingToolbar.module.css";

interface FloatingToolbarProps {
  onRefresh: () => void;
}

export default function FloatingToolbar({ onRefresh }: FloatingToolbarProps) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [opened, { open, close }] = useDisclosure(false);
  const { stages, removeStage, resetStage } = useMergeArticleStore();
  const [voices, setVoices] = useState<{ value: string; label: string }[]>([]);
  const [loading, setLoading] = useState(false);

  const mergeform = useForm<MergeArticleRequest>({
    mode: "uncontrolled",
    initialValues: { voiceIds: [], title: "" },
    validate: {
      title: (value) => {
        if (value) return null;
        return t("validate.required", { label: t("podcast.merge.title", { ns: "task" }) });
      },
    },
  });

  const loadVoices = async () => {
    const resp = await getSystemConfig<TextToSpeechAIConfig>({ key: SystemConfigKey.TextToSpeechAi }, true);

    if (!resp || !resp.value || !resp.value.voices) return [];

    setVoices(resp.value.voices.map((v) => ({ value: v.id, label: v.name })));
  };

  const onMergeArticle = async () => {
    setLoading(true);
    const resp = await mergeArticle(GolbalLanguage.getLanguage(), {
      ...mergeform.getValues(),
      stageIds: stages.map((s) => s.id),
    });
    resetStage();
    mergeform.reset();
    setLoading(false);
    close();
    if (resp && resp.batchNo) navigate("/task/" + resp.batchNo);
  };

  return (
    <>
      <Affix position={{ bottom: 30, right: 50 }}>
        <Transition transition="slide-up" mounted={true}>
          {(transitionStyles) => (
            <Stack style={transitionStyles}>
              <ActionIcon onClick={onRefresh} size="xl" variant="filled" color="blue" className={styles.floatingButton}>
                <IconRefresh size={20} />
              </ActionIcon>

              <Box className={styles.mergeButtonWrapper}>
                <ActionIcon onClick={open} size="xl" variant="filled" color="violet" className={styles.floatingButton}>
                  <IconLayersDifference size={20} />
                </ActionIcon>
                {stages.length > 0 && <Box className={styles.badge}>{stages.length}</Box>}
              </Box>
            </Stack>
          )}
        </Transition>
      </Affix>

      <Modal
        opened={opened}
        onClose={close}
        title={
          <Text size="lg" fw={600} c="violet">
            {t("podcast.stage.merge", { ns: "task" })}
          </Text>
        }
        size="lg"
        radius="md"
        overlayProps={{
          backgroundOpacity: 0.55,
          blur: 3,
        }}
      >
        <Stack gap="lg">
          {stages.length > 0 && (
            <>
              <Box>
                <Text size="sm" fw={500} mb="xs" c="dimmed">
                  {t("podcast.merge.content", { ns: "task" })} ({stages.length})
                </Text>
                <Stack gap="xs">
                  {stages.map((stage, index) => (
                    <Paper key={index} p="md" radius="md" withBorder className={styles.stageItem}>
                      <Group justify="space-between" wrap="nowrap">
                        <Stack gap={6} style={{ flex: 1, minWidth: 0 }}>
                          <Text size="sm" lineClamp={2} fw={500}>
                            {stage.output}
                          </Text>
                          <Badge variant="light" color="violet" size="sm" radius="sm">
                            {stage.batchNo}
                          </Badge>
                        </Stack>
                        <ActionIcon variant="subtle" color="red" onClick={() => removeStage(stage.id)} radius="md">
                          <IconTrash size={18} />
                        </ActionIcon>
                      </Group>
                    </Paper>
                  ))}
                </Stack>
              </Box>
              <Divider />
            </>
          )}

          <MultiSelect
            data={voices}
            label={t("podcast.merge.voice", { ns: "task" })}
            searchable
            onDropdownOpen={loadVoices}
            key={mergeform.key("voiceIds")}
            {...mergeform.getInputProps("voiceIds")}
            classNames={{ label: styles.formLabel }}
          />

          <TextInput
            required
            label={t("podcast.merge.title", { ns: "task" })}
            key={mergeform.key("title")}
            {...mergeform.getInputProps("title")}
            classNames={{ label: styles.formLabel }}
          />

          <Group justify="flex-end" mt="md">
            <Button variant="light" onClick={close} color="gray" radius="md">
              {t("button.cancel")}
            </Button>
            <Button
              disabled={loading}
              loading={loading}
              onClick={onMergeArticle}
              variant="gradient"
              gradient={{ from: "violet", to: "grape", deg: 135 }}
              radius="md"
            >
              {t("button.save")}
            </Button>
          </Group>
        </Stack>
      </Modal>
    </>
  );
}

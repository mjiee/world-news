import { SourceLabel } from "@/components";
import {
  autoTask,
  createTask,
  deleteNews,
  getNewsDetail,
  getSystemConfig,
  NewsDetail,
  saveFavorite,
  SystemConfigKey,
  TextToSpeechAIConfig,
  translateNews,
} from "@/services";
import { GolbalLanguage, useTaskPollingStore } from "@/stores";
import { isWeb } from "@/utils/platform";
import { ActionIcon, Badge, Button, Group, Modal, MultiSelect, Stack, Text } from "@mantine/core";
import { useField } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { IconBroadcast, IconLanguage, IconStar, IconStarFilled, IconTrash } from "@tabler/icons-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

interface NewsCardFooterProps {
  news: NewsDetail;
  updatePage: (page: number) => void;
  updateTitle: (title: string) => void;
  showTask?: boolean;
}

// news card footer
export default function NewsCardFooter({ news, updatePage, updateTitle, showTask }: NewsCardFooterProps) {
  const [favorited, setFavorited] = useState<boolean>(news?.favorited ?? false);
  const { addToQueue } = useTaskPollingStore();
  const translateTitle = async () => {
    const resp = await translateNews({ contents: [news.title], toLang: GolbalLanguage.getLanguage() });

    if (resp && resp.length > 0) updateTitle(resp[0]);
  };

  // save news favorite
  const saveNewsFavorite = async () => {
    await saveFavorite({ id: news.id, favorited: !favorited });
    setFavorited(!favorited);

    if (favorited) return;

    const resp = await getSystemConfig<TextToSpeechAIConfig>({ key: SystemConfigKey.TextToSpeechAi }, !isWeb());
    if (!resp || !resp.value || !resp.value.autoTask) return;

    const newsDetail = await getNewsDetail({ id: news.id });
    if (!newsDetail) return;

    const autoResp = await autoTask(GolbalLanguage.getLanguage(), newsDetail);
    if (autoResp) addToQueue(autoResp.batchNo);
  };

  const newsCardfooter = (txt: string, color: string = "dimmed") => (
    <Badge variant="light" color={color} size="sm">
      {txt}
    </Badge>
  );

  return (
    <Group justify="space-between" mt={5} mb="xs" onClick={(event) => event.stopPropagation()}>
      <Group gap="xs">
        {news.source && <SourceLabel source={news.source} size="sm" />}
        {news.topic && newsCardfooter(news.topic, "green")}
        {news.publishedAt && (
          <Text size="xs" c="dimmed">
            {news.publishedAt}
          </Text>
        )}
      </Group>
      <Group gap="xs">
        <ActionIcon variant="subtle" color="gray" size="sm" onClick={translateTitle}>
          <IconLanguage />
        </ActionIcon>
        <ActionIcon variant="subtle" color={favorited ? "yellow" : "gray"} size="sm" onClick={saveNewsFavorite}>
          {favorited ? <IconStarFilled /> : <IconStar />}
        </ActionIcon>
        {showTask && <CreateTaskButton newsId={news.id} />}
        <DeleteNewsButton newsId={news.id} updatePage={updatePage} />
      </Group>
    </Group>
  );
}

// delete news button
interface DeleteNewsButtonProps {
  newsId: number;
  updatePage: (page: number) => void;
}

function DeleteNewsButton({ newsId, updatePage }: DeleteNewsButtonProps) {
  const [opened, { open, close }] = useDisclosure(false);
  const { t } = useTranslation();

  // click ok handler
  const clickOkHandler = async () => {
    await deleteNews({ id: newsId });
    close();
    updatePage(0);
  };

  return (
    <>
      <Modal opened={opened} onClose={close} withCloseButton={false}>
        <p>{t("news_list.delete_label", { ns: "news" })}</p>
        <Group justify="flex-end">
          <Button onClick={clickOkHandler}>{t("button.ok")}</Button>
          <Button onClick={close} variant="default">
            {t("button.cancel")}
          </Button>
        </Group>
      </Modal>
      <ActionIcon variant="subtle" color="gray" size="sm" onClick={open}>
        <IconTrash />
      </ActionIcon>
    </>
  );
}

// create task button
interface CreateTaskButtonProps {
  newsId: number;
}

function CreateTaskButton({ newsId }: CreateTaskButtonProps) {
  const { t } = useTranslation();
  const [voices, setVoices] = useState<{ value: string; label: string }[]>([]);
  const [opened, { open, close }] = useDisclosure(false);
  const [loading, setLoading] = useState(false);
  const voiceField = useField({ initialValue: [] });

  const createPodcastcaTask = async () => {
    setLoading(true);
    const news = await getNewsDetail({ id: newsId });

    if (!news) return;

    await createTask(GolbalLanguage.getLanguage(), news, voiceField.getValue());
    setLoading(false);
    voiceField.setValue([]);
    close();
  };

  const loadVoices = async () => {
    const resp = await getSystemConfig<TextToSpeechAIConfig>({ key: SystemConfigKey.TextToSpeechAi }, true);

    if (!resp || !resp.value || !resp.value.voices) return [];

    setVoices(resp.value.voices.map((v) => ({ value: v.id, label: v.name })));
  };

  return (
    <>
      <ActionIcon variant="subtle" color="gray" size="sm" onClick={open}>
        <IconBroadcast />
      </ActionIcon>

      <Modal opened={opened} onClose={close} withCloseButton={false} title={"创建任务"} size="lg">
        <Stack gap="md">
          <MultiSelect data={voices} {...voiceField.getInputProps()} label="播音" onDropdownOpen={loadVoices} />
          <Button
            variant="gradient"
            gradient={{ from: "violet", to: "grape" }}
            loading={loading}
            onClick={createPodcastcaTask}
          >
            {t("button.save")}
          </Button>
        </Stack>
      </Modal>
    </>
  );
}

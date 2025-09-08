import { deleteNews, NewsDetail, saveFavorite, translateNews } from "@/services";
import { GolbalLanguage } from "@/stores";
import { ActionIcon, Badge, Button, Group, Modal, Text } from "@mantine/core";
import { useState } from "react";
import { useDisclosure } from "@mantine/hooks";
import { useTranslation } from "react-i18next";
import { SourceLabel } from "@/components";
import IconTrash from "@/assets/icons/IconTrash.svg?react";
import IconLanguage from "@/assets/icons/IconLanguage.svg?react";
import IconStar from "@/assets/icons/IconStar.svg?react";
import IconStarFilled from "@/assets/icons/IconStarFilled.svg?react";

interface NewsCardFooterProps {
  news: NewsDetail;
  updatePage: (page: number) => void;
  updateTitle: (title: string) => void;
}

// news card footer
export default function NewsCardFooter({ news, updatePage, updateTitle }: NewsCardFooterProps) {
  const [favorited, setFavorited] = useState<boolean>(news?.favorited ?? false);

  const translateTitle = async () => {
    const resp = await translateNews({ contents: [news.title], toLang: GolbalLanguage.getLanguage() });

    if (resp && resp.length > 0) updateTitle(resp[0]);
  };

  // save news favorite
  const saveNewsFavorite = async () => {
    await saveFavorite({ id: news.id, favorited: !favorited });
    setFavorited(!favorited);
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

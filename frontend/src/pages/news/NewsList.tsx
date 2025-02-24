import { useState, useEffect } from "react";
import { useNavigate } from "react-router";
import { useTranslation } from "react-i18next";
import {
  AspectRatio,
  Modal,
  Button,
  Card,
  Group,
  Image,
  Select,
  SimpleGrid,
  Text,
  ActionIcon,
  Pagination,
} from "@mantine/core";
import { useForm, UseFormReturnType } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import {
  getSystemConfig,
  queryNews,
  getCrawlingRecord,
  deleteNews,
  SystemConfigKey,
  NewsDetail,
  NewsWebsiteValue,
} from "@/services";
import { getPageNumber } from "@/utils/pagination";
import { getHost } from "@/utils/url";
import { httpx } from "wailsjs/go/models";
import classes from "./styles/newsList.module.css";
import IconSearch from "@/assets/icons/IconSearch.svg?react";
import IconTrash from "@/assets/icons/IconTrash.svg?react";

// news list component
interface NewsListProps {
  recordId: number;
}

export function NewsList({ recordId }: NewsListProps) {
  const searchFrom = useForm({
    mode: "uncontrolled",
    initialValues: { source: "", topic: "" },
  });

  const [newsList, setNewsList] = useState<NewsDetail[]>([]);
  const [pagination, setPagination] = useState<httpx.Pagination>({ page: 1, limit: 25, total: 0 });
  const [loading, setLoading] = useState<boolean>(true);

  // fetch news
  const fetchNews = async () => {
    if (!loading) return;

    const resp = await queryNews({ ...searchFrom.getValues(), recordId: recordId, pagination: pagination });

    setLoading(false);

    if (!resp || !resp.data) return;

    setNewsList(resp.data);
    setPagination({ ...pagination, total: resp.total });
  };

  useEffect(() => {
    fetchNews();
  }, [loading]);

  // update page
  const updatePageHandler = (page: number) => {
    setPagination({ ...pagination, page: page });
    setLoading(true);
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  // searchHandler
  const searchHandler = async () => {
    setLoading(true);
  };

  return (
    <>
      <SearchNews recordId={recordId} searchFrom={searchFrom} searchHandler={searchHandler} />
      <SimpleGrid cols={{ base: 1, sm: 2 }}>
        {newsList.map((item) => (
          <NewsCard key={item.id} news={item} updatePage={updatePageHandler} />
        ))}
      </SimpleGrid>
      <Pagination p="md" value={pagination.page} total={getPageNumber(pagination)} onChange={updatePageHandler} />
    </>
  );
}

// search news component
interface SearchNewsProps {
  recordId: number;
  searchFrom: UseFormReturnType<{ source: string; topic: string }>;
  searchHandler: () => void;
}

function SearchNews({ recordId, searchFrom, searchHandler }: SearchNewsProps) {
  const { t } = useTranslation("news");
  const [sources, setSources] = useState<string[]>([]);
  const [topics, setTopics] = useState<string[]>([]);

  // fetch data
  const fetchData = async () => {
    if (recordId > 0) {
      const resp = await getCrawlingRecord({ id: recordId });

      if (!resp || !resp.config) return;

      setSources(resp?.config?.sources);
      setTopics(resp.config.topics);
    } else {
      const sourceConfig = await getSystemConfig<NewsWebsiteValue[]>({ key: SystemConfigKey.NewsWebsites });
      const topicsConfig = await getSystemConfig<string[]>({ key: SystemConfigKey.NewsTopics });

      if (topicsConfig && topicsConfig.value) setTopics(topicsConfig.value);

      if (sourceConfig && sourceConfig.value) {
        setSources(sourceConfig?.value?.map((item) => getHost(item.url)));
      }
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const select = (label: string, key: string, data: string[]) => (
    <Select
      label={label}
      limit={100}
      data={data}
      searchable
      clearable
      key={searchFrom.key(key)}
      {...searchFrom.getInputProps(key)}
    />
  );

  return (
    <Group gap="sm" p="md" mb="md" align="flex-end">
      {select(t("news_list.search.source"), "source", sources)}
      {select(t("news_list.search.topic"), "topic", topics)}
      <ActionIcon onClick={searchHandler} variant="filled" aria-label="Settings" mb={3}>
        <IconSearch />
      </ActionIcon>
    </Group>
  );
}

// news card component
interface NewsCardProps {
  news: NewsDetail;
  updatePage: (page: number) => void;
}

function NewsCard({ news, updatePage }: NewsCardProps) {
  const navigate = useNavigate();

  return (
    <Card key={news.id} p="md" radius="md" className={classes.card}>
      <div onClick={() => navigate("/news/detail/" + news.id)}>
        {news.images.length === 0 ? (
          <p>{news.contents[0]}</p>
        ) : (
          <AspectRatio ratio={1920 / 1080}>
            <Image src={news.images[0]} fallbackSrc="https://placehold.co/200x100?text=Placeholder" />
          </AspectRatio>
        )}

        <Text mt="md">{news.title}</Text>
      </div>
      <Group justify="space-between" mt={5} mb="xs">
        <Group gap="xs">
          {newsCardfooter(news.source)}
          {newsCardfooter(news.topic)}
          {newsCardfooter(news.publishedAt)}
        </Group>

        <DeleteNewsButton newsId={news.id} updatePage={updatePage} />
      </Group>
    </Card>
  );
}

// news card footer
const newsCardfooter = (txt: string) => {
  return (
    <Text c="dimmed" size="xs" tt="uppercase" fw={600}>
      {txt}
    </Text>
  );
};

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
    updatePage(1);
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

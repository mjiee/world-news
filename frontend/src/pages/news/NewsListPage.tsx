import { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router";
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
  Title,
  ActionIcon,
  Pagination,
} from "@mantine/core";
import { useForm, UseFormReturnType } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { FetchNewsButton } from "@/components";
import {
  getSystemConfig,
  queryNews,
  getCrawlingRecord,
  deleteNews,
  SystemConfigKey,
  NewsDetail,
  NewsWebsiteValue,
  translateNews,
} from "@/services";
import dayjs from "dayjs";
import { DateInput } from "@/components";
import { GolbalLanguage, useRemoteServiceStore } from "@/stores";
import { getPageNumber } from "@/utils/pagination";
import { getHost } from "@/utils/url";
import { httpx } from "wailsjs/go/models";
import classes from "./styles/newsList.module.css";
import IconTrash from "@/assets/icons/IconTrash.svg?react";
import IconLanguage from "@/assets/icons/IconLanguage.svg?react";

// news list page
export function NewsListPage() {
  const { recordId } = useParams();
  const [newsList, setNewsList] = useState<NewsDetail[]>([]);
  const [pagination, setPagination] = useState<httpx.Pagination>({ page: 1, limit: 20, total: 0 });
  const [loading, setLoading] = useState<boolean>(true);
  const enableService = useRemoteServiceStore((state) => state.enable);

  const recordID = Number(recordId) || 0;

  const searchFrom = useForm({
    mode: "uncontrolled",
    initialValues: { source: "", topic: "", publishDate: "" },
  });

  // fetch news
  const fetchNews = async () => {
    if (!loading) return;

    const resp = await queryNews({ ...searchFrom.getValues(), recordId: recordID, pagination: pagination });

    setLoading(false);

    if (!resp || !resp.data) return;

    setNewsList(resp.data);
    setPagination({ ...pagination, total: resp.total });
  };

  useEffect(() => {
    fetchNews();
  }, [loading, enableService]);

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
      <SearchNews recordId={recordID} searchFrom={searchFrom} searchHandler={searchHandler} />
      <SimpleGrid cols={{ base: 1, sm: 2, lg: 3, xl: 4 }}>
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
  searchFrom: UseFormReturnType<{ source: string; topic: string; publishDate: string }>;
  searchHandler: () => void;
}

function SearchNews({ recordId, searchFrom, searchHandler }: SearchNewsProps) {
  const { t } = useTranslation();
  const [sources, setSources] = useState<string[]>([]);
  const [topics, setTopics] = useState<string[]>([]);
  const [publishDate, setPublishDate] = useState<Date | null>(null);

  const setSearchPublishDate = (date: Date | null) => {
    setPublishDate(date);
    searchFrom.setFieldValue("publishDate", date ? dayjs(date).format("YYYY-MM-DD HH:mm:ss") : "");
  };

  // fetch data
  const fetchData = async () => {
    let sourceData: string[] = [];
    let topicsData: string[] = [];

    if (recordId > 0) {
      const resp = await getCrawlingRecord({ id: recordId });

      if (!resp || !resp.config) return;

      sourceData = resp.config.sources;
      topicsData = resp.config.topics;
    } else {
      const sourceConfig = await getSystemConfig<NewsWebsiteValue[]>({ key: SystemConfigKey.NewsWebsites });
      const topicsConfig = await getSystemConfig<string[]>({ key: SystemConfigKey.NewsTopics });

      if (topicsConfig && topicsConfig.value) topicsData = topicsConfig?.value;

      if (sourceConfig && sourceConfig.value) sourceData = sourceConfig?.value?.map((item) => getHost(item.url));
    }

    setSources(sourceData.filter((item, index) => sourceData.indexOf(item) === index));
    setTopics(topicsData.filter((item, index) => topicsData.indexOf(item) === index));
  };

  useEffect(() => {
    fetchData();
  }, []);

  const select = (key: string, data: string[]) => (
    <Select
      placeholder={t("news_list.search." + key, { ns: "news" })}
      limit={100}
      data={data}
      searchable
      clearable
      key={searchFrom.key(key)}
      {...searchFrom.getInputProps(key)}
    />
  );

  return (
    <Group gap="sm" p="md" mb="md" align="flex-end" justify="center">
      {select("source", sources)}
      {select("topic", topics)}
      <DateInput
        placeholder={t("news_list.search.publish_date", { ns: "news" })}
        value={publishDate}
        onChange={setSearchPublishDate}
      />
      <Button onClick={searchHandler} variant="filled" aria-label="Settings">
        {t("button.search")}
      </Button>
      <FetchNewsButton />
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
  const [title, setTitle] = useState(news.title);

  // translate title
  const translateTitle = async () => {
    const resp = await translateNews({ id: 0, texts: [title], toLang: GolbalLanguage.getLanguage() });

    if (resp && resp.length > 0) setTitle(resp[0]);
  };

  return (
    <Card key={news.id} p="md" radius="md" className={classes.card} onClick={() => navigate("/news/detail/" + news.id)}>
      <NewsCardContent news={news} />

      <Title order={4} c="blue" mt="md">
        {title}
      </Title>

      <Group justify="space-between" mt={5} mb="xs" onClick={(event) => event.stopPropagation()}>
        <NewsCardFooter news={news} />

        <Group gap="xs">
          <ActionIcon variant="subtle" color="gray" size="sm" onClick={translateTitle}>
            <IconLanguage />
          </ActionIcon>
          <DeleteNewsButton newsId={news.id} updatePage={updatePage} />
        </Group>
      </Group>
    </Card>
  );
}

// news card content
function NewsCardContent({ news }: { news: NewsDetail }) {
  return (
    <div>
      {news.images.length === 0 ? (
        <p>{news.contents[0]}</p>
      ) : (
        <AspectRatio ratio={1920 / 1080}>
          <Image src={news.images[0]} fallbackSrc="https://placehold.co/200x100?text=Placeholder" />
        </AspectRatio>
      )}
    </div>
  );
}

// news card footer
function NewsCardFooter({ news }: { news: NewsDetail }) {
  const newsCardfooter = (txt: string) => (
    <Text c="dimmed" size="xs" fw={600}>
      {txt}
    </Text>
  );

  return (
    <Group gap="xs">
      {newsCardfooter(news.source)}
      {newsCardfooter(news.topic)}
      {newsCardfooter(news.publishedAt)}
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

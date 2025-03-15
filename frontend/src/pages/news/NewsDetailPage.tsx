import { useState, useEffect } from "react";
import { useParams } from "react-router";
import { useTranslation } from "react-i18next";
import {
  Center,
  Container,
  Title,
  ActionIcon,
  CopyButton,
  Text,
  Flex,
  Box,
  Image,
  Divider,
  Loader,
  Paper,
  Group,
  Button,
} from "@mantine/core";
import { getNewsDetail, NewsDetail, critiqueNews, translateNews } from "@/services";
import IconCopy from "@/assets/icons/IconCopy.svg?react";
import IconCheck from "@/assets/icons/IconCheck.svg?react";

// News detail page
export function NewsDetailPage() {
  const { newsId } = useParams();
  const [newsDetail, setNewsDetail] = useState<NewsDetail>();
  const [loading, setLoading] = useState<boolean>(true);

  // fetch news
  const fetchNews = async () => {
    if (!newsId) return;

    const resp = await getNewsDetail({ id: Number(newsId) });

    if (!resp) return;

    setNewsDetail(resp);
    setLoading(false);
  };

  useEffect(() => {
    fetchNews();
  }, []);

  return (
    <>
      {loading || newsDetail === undefined ? (
        <Center h={200}>
          <Loader color="blue" size="xl" type="bars" />
        </Center>
      ) : (
        <Container size="md">
          <Title mb="md" c="blue">
            {newsDetail?.title}
          </Title>
          <Text c="dimmed">{newsDetail?.publishedAt}</Text>
          <NewsLink link={newsDetail?.link} />
          <Divider my="md" />
          <Paper shadow="md" radius="md" withBorder p="lg">
            {newsBody(newsDetail?.contents, newsDetail?.images)} <NewsExtension newsId={newsDetail.id} />
          </Paper>
        </Container>
      )}
    </>
  );
}

// news body
const newsBody = (contents: string[], images: string[]) => {
  const maxLength = Math.max(contents.length, images.length);

  return [...Array(maxLength)].map((_, idx) => {
    return (
      <Box key={idx}>
        {idx < images.length && <Image src={images[idx]} fallbackSrc="https://placehold.co/400x50?text=Placeholder" />}
        {idx < contents.length && <p>{contents[idx]}</p>}
      </Box>
    );
  });
};

// news link
interface NewsLinkProps {
  link: string;
}

function NewsLink({ link }: NewsLinkProps) {
  const { t } = useTranslation("news");

  return (
    <Flex>
      <Text c="dimmed">{link}</Text>
      <CopyButton value={link} timeout={2000}>
        {({ copied, copy }) => (
          <ActionIcon p={3} color={copied ? "teal" : "gray"} variant="subtle" onClick={copy}>
            {copied ? <IconCheck /> : <IconCopy />}
          </ActionIcon>
        )}
      </CopyButton>
    </Flex>
  );
}

// news extension
interface NewsExtensionProps {
  newsId: number;
}

const critique = "critique";
const translate = "translate";

function NewsExtension({ newsId }: NewsExtensionProps) {
  const { t, i18n } = useTranslation("news");
  const [extension, setExtension] = useState<string>("");
  const [contents, setContents] = useState<string[]>([]);
  const [loading, setLoading] = useState<boolean>(false);

  const onClickHandle = async (obj: string) => {
    setLoading(true);

    let resp: string[] | undefined;

    if (obj === critique) {
      resp = await critiqueNews({ id: newsId });
    } else if (obj === translate) {
      resp = await translateNews({ id: newsId, toLang: i18n.language, texts: [] });
    }

    if (!resp) {
      setLoading(false);
      return;
    }

    setExtension(obj);
    setContents(resp);
    setLoading(false);
  };

  return (
    <>
      <Group>
        <Button loading={loading} onClick={() => onClickHandle(translate)}>
          {t("news_detail." + translate)}
        </Button>
        <Button loading={loading} onClick={() => onClickHandle(critique)}>
          {t("news_detail." + critique)}
        </Button>
      </Group>

      {extension !== "" && (
        <Container mt="md" p="md" bg="#f4f6f9">
          <Title order={4} c="blue">
            {t("news_detail." + extension)}
          </Title>
          {contents.map((item, idx) => (
            <p key={idx}>{item}</p>
          ))}
        </Container>
      )}
    </>
  );
}

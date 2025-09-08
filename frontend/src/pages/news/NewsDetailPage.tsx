import { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router";
import { useTranslation } from "react-i18next";
import {
  Container,
  Title,
  ActionIcon,
  CopyButton,
  Text,
  Badge,
  Flex,
  Box,
  Image,
  Divider,
  Paper,
  Group,
  Anchor,
  Stack,
  Transition,
  Affix,
  Popover,
} from "@mantine/core";
import MarkdownIt from "markdown-it";
import { Loading, SourceLabel } from "@/components";
import { getNewsDetail, NewsDetail, critiqueNews, translateNews, saveFavorite } from "@/services";
import classes from "./styles/newsDetail.module.css";
import IconCopy from "@/assets/icons/IconCopy.svg?react";
import IconCheck from "@/assets/icons/IconCheck.svg?react";
import IconLanguage from "@/assets/icons/IconLanguage.svg?react";
import IconAi from "@/assets/icons/IconAi.svg?react";
import IconArrowLeft from "@/assets/icons/IconArrowLeft.svg?react";
import IconX from "@/assets/icons/IconX.svg?react";
import IconStar from "@/assets/icons/IconStar.svg?react";
import IconStarFilled from "@/assets/icons/IconStarFilled.svg?react";

// News detail page
export function NewsDetailPage() {
  const { t } = useTranslation();
  const { newsId } = useParams();
  const navigate = useNavigate();
  const [newsDetail, setNewsDetail] = useState<NewsDetail>();
  const [translations, setTranslations] = useState<string[]>([]);
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

  if (loading || newsDetail === undefined) return <Loading />;

  const newsLabel = (txt: string, color: string = "dimmed") => (
    <Badge variant="light" color={color} size="md">
      {txt}
    </Badge>
  );

  return (
    <Container size="md">
      <Group mb="xl">
        <ActionIcon variant="subtle" color="gray" size="lg" onClick={() => navigate(-1)} aria-label={t("button.back")}>
          <IconArrowLeft />
        </ActionIcon>
        <Text c="dimmed" size="sm">
          {t("button.back")}
        </Text>
      </Group>
      <Title order={1} size="h2" c="dark" pb="xs">
        {newsDetail?.title}
      </Title>

      <Group mb="xs">
        <SourceLabel source={newsDetail.source} />
        {newsDetail.topic && newsLabel(newsDetail.topic, "green")}
        {newsDetail.publishedAt && (
          <Text c="dimmed" size="sm">
            {newsDetail.publishedAt}
          </Text>
        )}
      </Group>

      <NewsLink link={newsDetail?.link} />
      <Divider my="lg" />
      <Paper shadow="md" radius="md" withBorder p="lg">
        <NewsBody contents={newsDetail?.contents} images={newsDetail?.images} translations={translations} />
      </Paper>
      <FloatingToolbar newsDetail={newsDetail} setTranslations={setTranslations} />
    </Container>
  );
}

// news body
interface NewsBodyProps {
  contents: string[] | undefined;
  images: string[] | undefined;
  translations: string[] | undefined;
}

const imgFallbackSrc = "https://placehold.co/400x50?text=Placeholder";

function NewsBody({ contents, images, translations }: NewsBodyProps) {
  const safeContents = contents || [];
  const safeImages = images || [];
  const safeTranslations = translations || [];
  const maxLength = Math.max(safeContents.length, safeImages.length, safeTranslations.length);

  return (
    <Stack gap="lg">
      {[...Array(maxLength)].map((_, idx) => (
        <Box key={idx}>
          {idx < safeImages.length && (
            <Box mb="lg">
              <Image
                src={safeImages[idx]}
                fallbackSrc={imgFallbackSrc}
                radius="md"
                fit="contain"
                style={{ maxHeight: "500px" }}
              />
            </Box>
          )}
          {idx < safeContents.length && (
            <Text size="md" lh={1.7} c="dark.7" className={classes.content}>
              {safeContents[idx]}
            </Text>
          )}
          {idx < safeTranslations.length && (
            <Text size="sm" lh={1.8} c="gray.6" className={classes.translation}>
              {safeTranslations[idx]}
            </Text>
          )}
        </Box>
      ))}
    </Stack>
  );
}

// news link
interface NewsLinkProps {
  link: string;
}

function NewsLink({ link }: NewsLinkProps) {
  const { t } = useTranslation("news");

  return (
    <Flex>
      <Text c="dimmed" size="sm">
        {t("news_detail.link")}
      </Text>
      <Anchor href={link} target="_blank" rel="noopener noreferrer" size="sm" truncate style={{ maxWidth: "500px" }}>
        {link}
      </Anchor>
      <CopyButton value={link} timeout={2000}>
        {({ copied, copy }) => (
          <ActionIcon p={3} color={copied ? "teal" : "gray"} variant="subtle" size="sm" onClick={copy}>
            {copied ? <IconCheck /> : <IconCopy />}
          </ActionIcon>
        )}
      </CopyButton>
    </Flex>
  );
}

// floating toolbar
interface FloatingToolbarProps {
  newsDetail: NewsDetail | undefined;
  setTranslations: (translations: string[]) => void;
}

const critique = "critique";
const translate = "translate";
const favorite = "favorite";
const md = new MarkdownIt();

function FloatingToolbar({ newsDetail, setTranslations }: FloatingToolbarProps) {
  const { t, i18n } = useTranslation();
  const [extension, setExtension] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [critiqueContent, setCritiqueContent] = useState<string>("");
  const [favorited, setFavorited] = useState<boolean>(newsDetail?.favorited ?? false);

  const onClickHandle = async (obj: string) => {
    setExtension(obj);
    setLoading(true);

    let resp: string[] | undefined;

    if (obj === critique && newsDetail) {
      resp = await critiqueNews({ title: newsDetail.title, contents: newsDetail.contents ?? [] });
    } else if (obj === translate && newsDetail) {
      resp = await translateNews({ toLang: i18n.language, contents: newsDetail.contents ?? [] });
    } else if (obj == favorite && newsDetail) {
      await saveFavorite({ id: newsDetail.id, favorited: !favorited });
      setFavorited(!favorited);
    }

    if (!resp) {
      setLoading(false);
      return;
    }

    if (obj === translate) {
      setTranslations(resp);
    } else if (obj === critique) {
      setCritiqueContent(md.render(resp.join("\n")));
    }

    setLoading(false);
  };

  const clearExtension = () => {
    setExtension("");
    setCritiqueContent("");
  };

  const actionButton = (obj: string, label: string, color: string) => {
    const button = (
      <ActionIcon
        variant={obj === favorite ? "light" : "filled"}
        color={color}
        size="xl"
        loading={loading && extension === obj}
        onClick={() => onClickHandle(obj)}
        disabled={loading}
        aria-label={label}
      >
        {obj === critique && <IconAi />}
        {obj === translate && <IconLanguage />}
        {obj === favorite && favorited && <IconStarFilled />}
        {obj === favorite && !favorited && <IconStar />}
      </ActionIcon>
    );

    if (obj === critique) {
      return (
        <Popover
          opened={extension === critique}
          onClose={clearExtension}
          position="left-end"
          offset={{ mainAxis: 10, crossAxis: 0 }}
        >
          <Popover.Target>{button}</Popover.Target>
          <Popover.Dropdown>
            <Stack gap="md" style={{ width: "500px" }}>
              <Group justify="space-between" align="center">
                <Text fw={600} c="blue.7">
                  {t("news_detail.critique", { ns: "news" })}
                </Text>
                <ActionIcon variant="subtle" color="gray" onClick={clearExtension} aria-label="Close">
                  <IconX />
                </ActionIcon>
              </Group>
              <Box style={{ maxHeight: "550px", overflow: "auto" }}>
                {loading ? (
                  <Loading />
                ) : (
                  <div dangerouslySetInnerHTML={{ __html: critiqueContent }} className={classes.critique} />
                )}
              </Box>
            </Stack>
          </Popover.Dropdown>
        </Popover>
      );
    }

    return button;
  };

  return (
    <Affix position={{ bottom: 30, right: 50 }}>
      <Transition transition="slide-up" mounted={true}>
        {(transitionStyles) => (
          <Stack gap="sm" style={{ ...transitionStyles }}>
            {actionButton(favorite, t("news_detail.favorite", { ns: "news" }), "yellow")}
            {actionButton(translate, t("news_detail.translate", { ns: "news" }), "blue")}
            {actionButton(critique, t("news_detail.critique", { ns: "news" }), "green")}
          </Stack>
        )}
      </Transition>
    </Affix>
  );
}

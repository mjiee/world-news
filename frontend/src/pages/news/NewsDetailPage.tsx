import { useState, useEffect } from "react";
import { useParams } from "react-router";
import { useTranslation } from "react-i18next";
import {
  Container,
  Title,
  ActionIcon,
  CopyButton,
  Text,
  Flex,
  Box,
  LoadingOverlay,
  Image,
  Divider,
} from "@mantine/core";
import { BackHeader } from "@/components/BackHeader";
import { getNewsDetail, NewsDetail } from "@/services";
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
      <BackHeader />
      <>
        {loading || newsDetail === undefined ? (
          <Box pos="relative">
            <LoadingOverlay visible={true} zIndex={1000} overlayProps={{ radius: "sm", blur: 2 }} />
          </Box>
        ) : (
          <Container size="md">
            <Title mb="md">{newsDetail?.title}</Title>
            <Text c="dimmed">{newsDetail?.publishedAt}</Text>
            <NewsLink link={newsDetail?.link} />
            <Divider my="md" />
            {newsBody(newsDetail?.contents, newsDetail?.images)}
          </Container>
        )}
      </>
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

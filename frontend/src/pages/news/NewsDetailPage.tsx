import { useState, useEffect } from "react";
import { useParams } from "react-router";
import { Container, Title, ActionIcon, CopyButton, Text, Flex, Loader } from "@mantine/core";
import { IconCopy, IconCheck } from "@tabler/icons-react";
import { BackHeader } from "@/components/BackHeader";
import { getNewsDetail, NewsDetail } from "@/services";

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

  return loading || newsDetail === undefined ? (
    <Loader color="blue" />
  ) : (
    <>
      <BackHeader />
      <Container size="md">
        <Title>{newsDetail?.title}</Title>
        <p style={{ color: "var(--mantine-color-gray-5)" }}>{newsDetail?.publishedAt}</p>
        <NewsLink link={newsDetail?.link} />
        {newsDetail?.contents.map((item, idx) => <p key={idx}>{item}</p>)}
        {newsDetail?.images.map((item, idx) => <img key={idx} src={item} alt="news" />)}
      </Container>
    </>
  );
}

interface NewsLinkProps {
  link: string;
}

function NewsLink({ link }: NewsLinkProps) {
  return (
    <Flex>
      <Text c="blue">{link}</Text>
      <CopyButton value={link} timeout={2000}>
        {({ copied, copy }) => (
          <ActionIcon color={copied ? "teal" : "gray"} variant="subtle" onClick={copy}>
            {copied ? <IconCheck size={16} /> : <IconCopy size={16} />}
          </ActionIcon>
        )}
      </CopyButton>
    </Flex>
  );
}

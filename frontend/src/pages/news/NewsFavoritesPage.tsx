import { Loading, Pagination } from "@/components";
import { NewsDetail, queryNews } from "@/services";
import { useRemoteServiceStore } from "@/stores";
import { getPageNumber } from "@/utils/pagination";
import { Card, Stack, Image, Title, Text, AspectRatio, Flex } from "@mantine/core";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { httpx } from "wailsjs/go/models";
import NewsCardFooter from "./components/NewsCardFooter";

// news favorites page
export function NewsFavoritesPage() {
  const [newsList, setNewsList] = useState<NewsDetail[]>([]);
  const [pagination, setPagination] = useState<httpx.Pagination>({ page: 1, limit: 20, total: 0 });
  const [loading, setLoading] = useState<boolean>(true);
  const enableService = useRemoteServiceStore((state) => state.enable);

  // fetch news
  const fetchNews = async () => {
    if (!loading) return;

    const resp = await queryNews({ favorited: true, pagination: pagination });

    setLoading(false);

    if (!resp || !resp.data) return;

    setNewsList(resp.data);
    setPagination({ ...pagination, total: resp.total });
  };

  // update page
  const updatePageHandler = (page: number) => {
    if (page) setPagination({ ...pagination, page: page });
    setLoading(true);
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  useEffect(() => {
    fetchNews();
  }, [loading, enableService]);

  return (
    <>
      {loading ? (
        <Loading />
      ) : (
        <Stack gap="lg" p="md">
          <Stack gap="md">
            {newsList.map((item) => (
              <NewsCard key={item.id} news={item} updatePage={updatePageHandler} />
            ))}
          </Stack>
          <Pagination value={pagination.page} total={getPageNumber(pagination)} onChange={updatePageHandler} />
        </Stack>
      )}
    </>
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

  return (
    <Card shadow="sm" radius="md" onClick={() => navigate("/news/detail/" + news.id)}>
      <Flex justify="flex-start" align="center" direction={{ base: "column", sm: "row" }} gap="md">
        {news.images && news.images.length > 0 && (
          <AspectRatio ratio={16 / 9} style={{ flexShrink: 0 }} w={{ base: "100%", sm: 200 }}>
            <Image src={news.images[0]} alt={news.title} radius="sm" fit="cover" style={{ flexShrink: 0 }} />
          </AspectRatio>
        )}

        <Stack gap="xs" style={{ flex: 1 }}>
          <Title order={4} c="blue.7" mt="md" lineClamp={2}>
            {news.title}
          </Title>

          {news.contents && news.contents.length > 0 && (
            <Text size="sm" lineClamp={2}>
              {news.contents[0]}
            </Text>
          )}

          <NewsCardFooter news={news} updatePage={updatePage} updateTitle={setTitle} />
        </Stack>
      </Flex>
    </Card>
  );
}

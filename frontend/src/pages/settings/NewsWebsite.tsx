import { useState, useEffect } from "react";
import { Button, Table, Pill, Stack, Pagination } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { getSystemConfig, crawlingWebsite, hasCrawlingTask } from "@/services";
import { getPageData, getPageNumber } from "@/utils/pagination";

// news website settings
const newsWebsiteCollectionKey = "newsWebsiteCollection";
const newsWebsiteKey = "newsWebsite";

interface NewsWebsiteValue {
  url: string;
  selectors?: string[];
}

export function NewsWebsiteCollection() {
  const [data, setData] = useState<NewsWebsiteValue[]>([]);

  const fetchData = async () => {
    const resp = await getSystemConfig<NewsWebsiteValue[]>({ key: newsWebsiteCollectionKey });

    if (!resp || !resp.value) return;
    if (resp.value.length === 0) return;

    setData(resp.value);
  };

  useEffect(() => {
    fetchData();
  }, []);

  return <WebsiteTable websites={data} />;
}

export function NewsWebsite() {
  const [data, setData] = useState<NewsWebsiteValue[]>([]);
  const { t } = useTranslation("settings");
  const [loading, setLoading] = useState(true);

  // fetch news website
  const fetchNewsWebsite = async () => {
    const processingTask = await hasCrawlingTask();

    if (!processingTask) setLoading(false);

    const resp = await getSystemConfig<NewsWebsiteValue[]>({ key: newsWebsiteKey });

    if (!resp || !resp.value) return;
    if (resp.value?.length === 0) return;

    setData(resp.value);
  };

  useEffect(() => {
    fetchNewsWebsite();
  }, []);

  // crawling website handle
  const crawlingWebsiteHandle = async () => {
    setLoading(true);
    await crawlingWebsite();
  };

  return (
    <Stack w={"100%"} align="stretch" justify="flex-start" gap="md">
      <Button variant="default" disabled={loading} onClick={crawlingWebsiteHandle}>
        {t("news_website.button.update_news_website")}
      </Button>
      <WebsiteTable websites={data} />
    </Stack>
  );
}

interface WebsiteTableProps {
  websites: NewsWebsiteValue[];
}

function WebsiteTable({ websites }: WebsiteTableProps) {
  const { t } = useTranslation("settings");
  const [page, setPage] = useState<number>(1);
  const [data, setData] = useState<NewsWebsiteValue[]>([]);

  // update data
  useEffect(() => {
    setData(getPageData<NewsWebsiteValue>(websites, page, 25));
  }, [websites, page]);

  // update page
  const updatePageHandle = (page: number) => {
    setPage(page);
    setData(getPageData<NewsWebsiteValue>(websites, page, 25));
  };

  // table header
  const tableHeader = (
    <Table.Tr>
      <Table.Th>{t("news_website.table.head.website")}</Table.Th>
      <Table.Th>{t("news_website.table.head.selector")}</Table.Th>
    </Table.Tr>
  );

  // table body
  const tableBody = data.map((item) => (
    <Table.Tr key={item.url}>
      <Table.Td>{item.url}</Table.Td>
      <Table.Td>
        <Pill.Group>{item.selectors?.map((value) => <Pill key={value}>{value}</Pill>)}</Pill.Group>
      </Table.Td>
    </Table.Tr>
  ));

  return (
    <Stack w={"100%"}>
      <Table>
        <Table.Thead>{tableHeader}</Table.Thead>
        <Table.Tbody>{tableBody}</Table.Tbody>
      </Table>
      <Pagination value={page} onChange={updatePageHandle} total={getPageNumber({ limit: 25, total: websites.length })} />
    </Stack>
  );
}

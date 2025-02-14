import { useState, useEffect } from "react";
import { Button, Table, Pill, Stack } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { getSystemConfig, crawlingWebsite } from "@/services";

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

  return <WebsiteTable data={data} />;
}

export function NewsWebsite() {
  const [data, setData] = useState<NewsWebsiteValue[]>([]);
  const { t } = useTranslation("settings");

  // fetch news website
  const fetchNewsWebsite = async () => {
    const resp = await getSystemConfig<NewsWebsiteValue[]>({ key: newsWebsiteKey });

    if (!resp || !resp.value) return;
    if (resp.value.length === 0) return;

    setData(resp.value);
  };

  useEffect(() => {
    fetchNewsWebsite();
  }, []);

  return (
    <Stack w={"100%"} align="stretch" justify="flex-start" gap="md">
      <Button variant="default" onClick={() => crawlingWebsite()}>
        {t("news_website.button.update_news_website")}
      </Button>
      <WebsiteTable data={data} />
    </Stack>
  );
}

interface WebsiteTableProps {
  data: NewsWebsiteValue[];
}

function WebsiteTable({ data }: WebsiteTableProps) {
  const { t } = useTranslation("settings");

  const tableHeader = (
    <Table.Tr>
      <Table.Th>{t("news_website.table.head.website")}</Table.Th>
      <Table.Th>{t("news_website.table.head.selector")}</Table.Th>
    </Table.Tr>
  );

  const tableBody = data.map((item) => (
    <Table.Tr key={item.url}>
      <Table.Td>{item.url}</Table.Td>
      <Table.Td>
        <Pill.Group>{item.selectors?.map((value) => <Pill key={value}>{value}</Pill>)}</Pill.Group>
      </Table.Td>
    </Table.Tr>
  ));

  return (
    <Table>
      <Table.Thead>{tableHeader}</Table.Thead>
      <Table.Tbody>{tableBody}</Table.Tbody>
    </Table>
  );
}

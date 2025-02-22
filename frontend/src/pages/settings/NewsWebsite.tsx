import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Button, Table, Pill, Stack, Pagination, Modal, Group, JsonInput, Flex } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { useField } from "@mantine/form";
import {
  getSystemConfig,
  saveSystemConfig,
  crawlingWebsite,
  hasCrawlingTask,
  SystemConfigKey,
  NewsWebsiteValue,
} from "@/services";
import { getPageData, getPageNumber } from "@/utils/pagination";
import { validateUrl } from "@/utils/url";

export function NewsWebsiteCollection() {
  const [data, setData] = useState<NewsWebsiteValue[]>([]);

  const fetchData = async () => {
    const resp = await getSystemConfig<NewsWebsiteValue[]>({ key: SystemConfigKey.NewsWebsiteCollections });

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

    const resp = await getSystemConfig<NewsWebsiteValue[]>({ key: SystemConfigKey.NewsWebsites });

    if (!resp || !resp.value) return;
    if (resp.value?.length === 0) return;

    setData(resp.value);
  };

  useEffect(() => {
    fetchNewsWebsite();
  }, []);

  // refresh data
  const refreshData = async () => {
    await fetchNewsWebsite();
  };

  // crawling website handle
  const crawlingWebsiteHandle = async () => {
    setLoading(true);
    await crawlingWebsite();
  };

  return (
    <Stack w={"100%"} align="stretch" justify="flex-start" gap="md">
      <Flex gap="md">
        <Button variant="default" disabled={loading} onClick={crawlingWebsiteHandle}>
          {t("news_website.button.fetch_news_website")}
        </Button>
        <UploadNewsWebsite refershData={refreshData} />
      </Flex>
      <WebsiteTable websites={data} />
    </Stack>
  );
}

// website table
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
      <Table.Td>{JSON.stringify(item.selector)}</Table.Td>
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

// UploadNewsWebsite upload news website
interface UploadNewsWebsiteProps {
  refershData: () => void;
}

function UploadNewsWebsite({ refershData }: UploadNewsWebsiteProps) {
  const { t } = useTranslation();
  const [opened, { open, close }] = useDisclosure(false);

  const website = useField({
    initialValue: "",
    validateOnChange: true,
    validate: (value) => {
      if (validateWebsiteValue(value)) return null;

      return t("news_website.upload_website.invalid_value", { ns: "settings" });
    },
  });

  // click ok handler
  const clickOkHandler = async () => {
    if (!website.validate()) return;

    const websites: NewsWebsiteValue[] = JSON.parse(website.getValue());

    await saveSystemConfig({ key: SystemConfigKey.NewsWebsites, value: websites });

    website.reset();
    close();
    refershData();
  };

  return (
    <>
      <Modal
        title={t("news_website.upload_website.title", { ns: "settings" })}
        opened={opened}
        onClose={close}
        withCloseButton={false}
      >
        <JsonInput {...website.getInputProps()} placeholder='[{"url": "https://news.com"}]' formatOnBlur rows={4} />
        <Group justify="flex-end">
          <Button onClick={clickOkHandler}>{t("button.ok")}</Button>
          <Button onClick={close} variant="default">
            {t("button.cancel")}
          </Button>
        </Group>
      </Modal>
      <Button variant="default" onClick={open}>
        {t("button.upload")}
      </Button>
    </>
  );
}

// validateWebsiteValue validate website value
function validateWebsiteValue(value: string): boolean {
  if (!value) return false;

  try {
    const websites: NewsWebsiteValue[] = JSON.parse(value);

    websites.forEach((item) => {
      if (!validateUrl(item.url)) throw new Error("invalid url");
    });

    return true;
  } catch (_) {
    return false;
  }
}

import {
  crawlingWebsite,
  getSystemConfig,
  hasCrawlingTask,
  NewsWebsiteValue,
  saveSystemConfig,
  SystemConfigKey,
} from "@/services";
import { validateUrl } from "@/utils/url";
import { Button, Group, JsonInput, Modal, Paper, Stack, Text } from "@mantine/core";
import { useField } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { useEffect, useState } from "react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";
import styles from "../styles/settings.module.css";

// News Website Collection Component
export function NewsWebsiteCollection() {
  const [data, setData] = useState<NewsWebsiteValue[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    const resp = await getSystemConfig<NewsWebsiteValue[]>({ key: SystemConfigKey.NewsWebsiteCollections });
    if (resp?.value) setData(resp.value);
    setLoading(false);
  };

  return <WebsiteList websites={data} loading={loading} />;
}

// News Website Component
export function NewsWebsite() {
  const { t } = useTranslation("settings");
  const [data, setData] = useState<NewsWebsiteValue[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    const processingTask = await hasCrawlingTask();
    if (!processingTask) setLoading(false);

    const resp = await getSystemConfig<NewsWebsiteValue[]>({ key: SystemConfigKey.NewsWebsites });
    if (resp?.value) setData(resp.value);
    setLoading(false);
  };

  const crawl = async () => {
    setLoading(true);
    await crawlingWebsite();
    await fetchData();
  };

  return (
    <Stack gap="md">
      <Group>
        <Button variant="default" disabled={loading} onClick={crawl}>
          {t("news_website.button.fetch_news_website")}
        </Button>
        <UploadNewsWebsite onSuccess={fetchData} />
      </Group>
      <WebsiteList websites={data} loading={loading} />
    </Stack>
  );
}

// Website List Component
interface WebsiteListProps {
  websites: NewsWebsiteValue[];
  loading: boolean;
}

function WebsiteList({ websites, loading }: WebsiteListProps) {
  if (loading) return <Text c="dimmed">Loading...</Text>;
  if (!websites.length) return <Text c="dimmed">No data available</Text>;

  return (
    <Stack gap="xs">
      {websites.map((item, idx) => (
        <Paper key={idx} p="md" withBorder className={styles.websiteItem}>
          <Stack gap="xs">
            <Text size="sm" fw={500}>
              {item.url}
            </Text>
            <Text size="xs" c="dimmed" className={styles.selectorText}>
              {JSON.stringify(item.selector)}
            </Text>
          </Stack>
        </Paper>
      ))}
    </Stack>
  );
}

// Upload News Website Component
function UploadNewsWebsite({ onSuccess }: { onSuccess: () => void }) {
  const { t } = useTranslation();
  const [opened, { open, close }] = useDisclosure(false);

  const website = useField({
    initialValue: "",
    validate: (value) => {
      if (!value) return false;
      try {
        const websites: NewsWebsiteValue[] = JSON.parse(value);
        return websites.every((item) => validateUrl(item.url));
      } catch {
        return false;
      }
    },
  });

  const handleSave = async () => {
    if (!website.validate()) return;
    const websites: NewsWebsiteValue[] = JSON.parse(website.getValue());
    await saveSystemConfig({ key: SystemConfigKey.NewsWebsites, value: websites });
    website.reset();
    close();
    onSuccess();
    toast.success("ok");
  };

  return (
    <>
      <Modal opened={opened} onClose={close} title={t("news_website.upload_website.title", { ns: "settings" })}>
        <Stack gap="md">
          <JsonInput {...website.getInputProps()} placeholder='[{"url": "https://news.com"}]' formatOnBlur rows={6} />
          <Group justify="flex-end">
            <Button variant="default" onClick={close}>
              {t("button.cancel")}
            </Button>
            <Button onClick={handleSave}>{t("button.ok")}</Button>
          </Group>
        </Stack>
      </Modal>
      <Button variant="default" onClick={open}>
        {t("button.upload")}
      </Button>
    </>
  );
}

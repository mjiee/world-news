import { useState, useEffect, useCallback, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { Button, Group, Modal, MultiSelect, Stack } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { useForm } from "@mantine/form";
import { hasCrawlingTask, crawlingNews, getSystemConfig, NewsWebsiteValue, SystemConfigKey } from "@/services";
import { DateInput } from "./DateInput";
import { getSecondLevelDomain } from "@/utils/url";

export function FetchNewsButton() {
  const { t } = useTranslation();
  const [opened, { open, close }] = useDisclosure(false);
  const [disabled, setDisabled] = useState(false);
  const [sources, setSources] = useState<string[]>([]);
  const [topics, setTopics] = useState<string[]>([]);

  const form = useForm({
    mode: "uncontrolled",
    initialValues: { sources: [], topics: [], startTime: "" },
  });

  const setStartTime = useCallback(
    (date: string | null) => {
      form.setFieldValue("startTime", date ?? "");
    },
    [form],
  );

  // Memoize processed data to avoid recalculation
  const processedSources = useMemo(() => [...new Set(sources.map(getSecondLevelDomain).filter(Boolean))], [sources]);

  const processedTopics = useMemo(() => [...new Set(topics.filter(Boolean))], [topics]);

  // Fetch config data
  const fetchData = useCallback(async () => {
    const [sourceConfig, topicsConfig] = await Promise.all([
      getSystemConfig<NewsWebsiteValue[]>({ key: SystemConfigKey.NewsWebsites }),
      getSystemConfig<string[]>({ key: SystemConfigKey.NewsTopics }),
    ]);

    if (sourceConfig?.value) {
      setSources(sourceConfig.value.map((item) => item.url));
    }
    if (topicsConfig?.value) {
      setTopics(topicsConfig.value);
    }
  }, []);

  // Check if crawling task is running
  const checkCanFetchNews = useCallback(async () => {
    const hasTask = await hasCrawlingTask();
    setDisabled(hasTask ?? false);
  }, []);

  const handleSubmit = useCallback(async () => {
    await crawlingNews(form.getValues());
    handleClose();
  }, [form]);

  const handleClose = useCallback(() => {
    form.reset();
    close();
  }, [form, close]);

  // Load data when modal opens
  useEffect(() => {
    if (opened) {
      fetchData();
      checkCanFetchNews();
    }
  }, [opened, fetchData, checkCanFetchNews]);

  // Reusable MultiSelect component
  const renderMultiSelect = (key: string, data: string[]) => (
    <MultiSelect
      label={t(`news_list.fetch_news.${key}`, { ns: "news" })}
      data={data}
      searchable
      clearable
      maxDropdownHeight={200}
      disabled={disabled}
      key={form.key(key)}
      {...form.getInputProps(key)}
    />
  );

  return (
    <>
      <Modal opened={opened} onClose={handleClose} withCloseButton={false}>
        <Stack>
          <DateInput label={t("news_list.fetch_news.start_time", { ns: "news" })} onChange={setStartTime} />
          {renderMultiSelect("sources", processedSources)}
          {renderMultiSelect("topics", processedTopics)}
        </Stack>
        <Group justify="flex-end" mt="md">
          <Button onClick={handleClose} variant="default">
            {t("button.cancel")}
          </Button>
          <Button onClick={handleSubmit} disabled={disabled}>
            {t("button.ok")}
          </Button>
        </Group>
      </Modal>

      <Button onClick={open}>{t("news_list.fetch_news.button", { ns: "news" })}</Button>
    </>
  );
}

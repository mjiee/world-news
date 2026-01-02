import { Loading, Pagination } from "@/components";
import {
  CrawlingRecord,
  CrawlingRecordStatus,
  CrawlingRecordType,
  deleteCrawlingRecord,
  queryCrawlingRecords,
  updateCrawlingRecordStatus,
} from "@/services";
import { getPageNumber } from "@/utils/pagination";
import { ActionIcon, Badge, Box, Button, Card, Flex, Group, Modal, Space, Stack, Text, Title } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { IconCalendar, IconClock, IconEye, IconPlayerPause, IconTrash } from "@tabler/icons-react";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router";
import { httpx } from "wailsjs/go/models";
import styles from "./styles/crawlingRecord.module.css";

const STATUS_COLORS: Record<CrawlingRecordStatus, string> = {
  [CrawlingRecordStatus.ProcessingCrawlingRecord]: "blue",
  [CrawlingRecordStatus.CompletedCrawlingRecord]: "green",
  [CrawlingRecordStatus.FailedCrawlingRecord]: "red",
  [CrawlingRecordStatus.PausedCrawlingRecord]: "gray",
};

export function CrawlingRecordPage() {
  const { t } = useTranslation("task");
  const [records, setRecords] = useState<CrawlingRecord[]>([]);
  const [pagination, setPagination] = useState<httpx.Pagination>({
    page: 1,
    limit: 25,
    total: 0,
  });
  const [loading, setLoading] = useState<boolean>(true);

  const fetchCrawlingRecords = async () => {
    if (!loading) return;

    const resp = await queryCrawlingRecords({ pagination: pagination });

    setLoading(false);

    if (!resp || !resp.data) return;

    setRecords(resp.data);
    setPagination({ ...pagination, total: resp.total });
  };

  useEffect(() => {
    fetchCrawlingRecords();
  }, [loading]);

  const updatePageHandler = (page: number) => {
    setPagination({ ...pagination, page: page });
    setLoading(true);
  };

  return (
    <>
      <Box mb="xs">
        <Title order={2} c="violet" mb={4}>
          {t("news.page_title")}
        </Title>
        <Text c="dimmed" size="sm">
          {t("news.page_description", { total: pagination.total })}
        </Text>
      </Box>

      {loading ? (
        <Loading />
      ) : (
        <Stack gap="md">
          {records.map((record) => (
            <RecordCard key={record.id} record={record} updatePage={updatePageHandler} />
          ))}
        </Stack>
      )}
      <Space h="xl" />
      <Pagination page={pagination.page} total={getPageNumber(pagination)} onChange={updatePageHandler} />
    </>
  );
}

interface RecordCardProps {
  record: CrawlingRecord;
  updatePage: (page: number) => void;
}

function RecordCard({ record, updatePage }: RecordCardProps) {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const handlePause = async () => {
    await updateCrawlingRecordStatus({
      id: record.id,
      status: CrawlingRecordStatus.PausedCrawlingRecord,
    });
    updatePage(1);
  };

  const handleView = () => {
    navigate("/news/list/" + record.id);
  };

  const isProcessing = record.status === CrawlingRecordStatus.ProcessingCrawlingRecord;
  const showViewButton = record.recordType !== CrawlingRecordType.CrawlingWebsite;

  return (
    <Card className={styles.recordCard} shadow="sm" padding="lg" radius="md">
      <Flex gap="md" align="flex-start">
        <Box className={styles.infoSection}>
          <Flex gap="xs" align="center" mb="xs">
            <Text size="sm" c="dimmed" fw={500}>
              ID:
            </Text>
            <Text size="sm" fw={600}>
              #{record.id}
            </Text>
            <Badge color={STATUS_COLORS[record.status]} variant="light" size="sm">
              {t("news.record.status." + record.status, { ns: "task" })}
            </Badge>
          </Flex>

          <Text size="lg" fw={500} mb="md">
            {t("news.record.record_type." + record.recordType, { ns: "task" })}
          </Text>

          <Group gap="xl">
            <Flex gap="xs" align="center">
              <Text size="sm" c="dimmed">
                {t("news.record.quantity", { ns: "task" })}:
              </Text>
              <Text size="sm" fw={500}>
                {record.quantity}
              </Text>
            </Flex>

            <Flex gap="xs" align="center">
              <IconCalendar size={16} className={styles.icon} />
              <Text size="xs" c="dimmed">
                {record.startTime}
              </Text>
            </Flex>

            {!isProcessing && record.endTime && (
              <Flex gap="xs" align="center">
                <IconClock size={16} className={styles.icon} />
                <Text size="xs" c="dimmed">
                  {record.endTime}
                </Text>
              </Flex>
            )}
          </Group>
        </Box>

        <Group gap="xs" className={styles.actionSection}>
          {showViewButton && (
            <ActionIcon variant="light" size="lg" onClick={handleView} title={t("button.view")}>
              <IconEye size={18} />
            </ActionIcon>
          )}

          {isProcessing && (
            <ActionIcon
              variant="light"
              color="orange"
              size="lg"
              onClick={handlePause}
              title={t("news.button.pause", { ns: "task" })}
            >
              <IconPlayerPause size={18} />
            </ActionIcon>
          )}

          {!isProcessing && <DeleteRecordButton record={record} updatePage={updatePage} />}
        </Group>
      </Flex>
    </Card>
  );
}

interface DeleteRecordButtonProps {
  record: CrawlingRecord;
  updatePage: (page: number) => void;
}

function DeleteRecordButton({ record, updatePage }: DeleteRecordButtonProps) {
  const [opened, { open, close }] = useDisclosure(false);
  const { t } = useTranslation();

  const clickOkHandler = async () => {
    await deleteCrawlingRecord({ id: record.id });
    close();
    updatePage(1);
  };

  return (
    <>
      <Modal opened={opened} onClose={close} title={t("button.delete", { ns: "record" })} centered>
        <Stack gap="md">
          <Text>{t("news.button.delete_label", { date: record.startTime, ns: "task" })}</Text>
          <Group justify="flex-end">
            <Button onClick={close} variant="default">
              {t("button.cancel")}
            </Button>
            <Button onClick={clickOkHandler} color="red">
              {t("button.ok")}
            </Button>
          </Group>
        </Stack>
      </Modal>

      <ActionIcon variant="light" color="red" size="lg" onClick={open} title={t("button.delete")}>
        <IconTrash size={18} />
      </ActionIcon>
    </>
  );
}

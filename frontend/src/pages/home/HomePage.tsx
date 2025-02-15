import { useState, useEffect } from "react";
import { useNavigate } from "react-router";
import { Button, Container, Avatar, Group, Table, Modal, LoadingOverlay, Pagination } from "@mantine/core";
import { DateInput } from "@mantine/dates";
import { useDisclosure } from "@mantine/hooks";
import { useTranslation } from "react-i18next";
import dayjs from "dayjs";
import { LanguageSwitcher } from "@/components";
import { hasCrawlingTask, crawlingNews, queryCrawlingRecords, deleteCrawlingRecord, CrawlingRecord } from "@/services";
import { httpx } from "wailsjs/go/models";
import { getPageNumber } from "@/utils/pagination";

import styles from "@/assets/styles/header.module.css";
import "dayjs/locale/en";
import "dayjs/locale/zh";

// Application homepage
export function HomePage() {
  return (
    <>
      <HomeHeader />
      <CrawlingRecords />
    </>
  );
}

function HomeHeader() {
  let navigate = useNavigate();
  const { t } = useTranslation("home");

  return (
    <header className={styles.header}>
      <Container size="md" className={styles.inner}>
        <Avatar size={28} name="World News" color="initials" />
        <Group>
          <FetchNewsButton />
          <Button onClick={() => navigate("/settings")}>{t("header.button.settings")}</Button>
          <LanguageSwitcher />
        </Group>
      </Container>
    </header>
  );
}

// fetch news button
function FetchNewsButton() {
  const { t, i18n } = useTranslation();
  const [opened, { open, close }] = useDisclosure(false);
  const [startTime, setStartTime] = useState<Date | null>(null);
  const [disabled, setDisabled] = useState<boolean>(false);

  // click ok handler
  const clickOkHandler = () => {
    crawlingNews({ startTime: startTime ? dayjs(startTime).format("YYYY-MM-DD HH:mm:ss") : "" });
    setStartTime(null);
    close();
  };

  // click cancel handler
  const clickCancelHandler = () => {
    setStartTime(null);
    close();
  };

  // check can fetch news
  const checkCanFetchNews = async () => {
    const processingTask = await hasCrawlingTask();

    if (processingTask) setDisabled(true);
    else setDisabled(false);
  };

  useEffect(() => {
    checkCanFetchNews();
  }, [opened]);

  return (
    <>
      <Modal opened={opened} onClose={close} withCloseButton={false}>
        <DateInput
          maxDate={new Date()}
          locale={i18n.language}
          label={t("header.label.start_time", { ns: "home" })}
          valueFormat="YYYY-MM-DD"
          value={startTime}
          onChange={setStartTime}
        />
        <Group justify="flex-end" mt="md">
          <Button disabled={disabled} type="submit" onClick={clickOkHandler}>
            {t("button.ok")}
          </Button>
          <Button onClick={clickCancelHandler} variant="default">
            {t("button.cancel")}
          </Button>
        </Group>
      </Modal>
      <Button onClick={open}>{t("header.button.fetch_news", { ns: "home" })}</Button>
    </>
  );
}

// crawling records
function CrawlingRecords() {
  const { t } = useTranslation("home");

  const tableHeader = (
    <Table.Tr>
      <Table.Th>ID</Table.Th>
      <Table.Th>{t("crawling_records.table.head.date")}</Table.Th>
      <Table.Th>{t("crawling_records.table.head.quantity")}</Table.Th>
      <Table.Th>{t("crawling_records.table.head.status")}</Table.Th>
      <Table.Th />
    </Table.Tr>
  );

  return (
    <Container size="md">
      <Table>
        <Table.Thead>{tableHeader}</Table.Thead>
        <RecordTableBody />
      </Table>
    </Container>
  );
}

function RecordTableBody() {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [records, setRecords] = useState<CrawlingRecord[]>([]);
  const [pagination, setPagination] = useState<httpx.Pagination>({ page: 1, limit: 25, total: 0 });
  const [loading, setLoading] = useState<boolean>(true);

  // fetch crawling records
  const fetchCrawlingRecords = async () => {
    if (!loading) return;

    const resp = await queryCrawlingRecords({ pagination: pagination });

    setLoading(false);

    if (!resp || !resp.data) return;

    setRecords(resp.data);
  };

  useEffect(() => {
    fetchCrawlingRecords();
  }, [loading]);

  // update page
  const updatePageHandler = (page: number) => {
    setPagination({ ...pagination, page: page });
    setLoading(true);
  };

  // record table body
  const recordTableBody = records.map((item) => (
    <Table.Tr key={item.id}>
      <Table.Td>{item.id}</Table.Td>
      <Table.Td>{item.date}</Table.Td>
      <Table.Td>{item.quantity}</Table.Td>
      <Table.Td>{t("crawling_records.table.body.status." + item.status, { ns: "home" })}</Table.Td>
      <Table.Td>
        <Button.Group>
          <Button variant="default" size="xs" onClick={() => navigate("/news/list/" + item.id)}>
            {t("button.view")}
          </Button>
          <DeleteRecordButton recordId={item.id} date={item.date} updatePage={updatePageHandler} />
        </Button.Group>
      </Table.Td>
    </Table.Tr>
  ));

  return (
    <>
      {loading ? <LoadingOverlay visible={loading} /> : <Table.Tbody>{recordTableBody}</Table.Tbody>}
      <Pagination value={pagination.page} total={getPageNumber(pagination)} onChange={updatePageHandler} />
    </>
  );
}

interface DeleteRecordButtonProps {
  recordId: number;
  date: String;
  updatePage: (page: number) => void;
}

function DeleteRecordButton({ recordId, date, updatePage }: DeleteRecordButtonProps) {
  const [opened, { open, close }] = useDisclosure(false);
  const { t } = useTranslation();

  // click ok handler
  const clickOkHandler = async () => {
    await deleteCrawlingRecord({ id: recordId });
    close();
    updatePage(1);
  };

  return (
    <>
      <Modal opened={opened} onClose={close} withCloseButton={false}>
        <p>{t("crawling_records.button.delete_label", { date, ns: "home" })}</p>
        <Group justify="flex-end">
          <Button onClick={clickOkHandler}>{t("button.ok")}</Button>
          <Button onClick={close} variant="default">
            {t("button.cancel")}
          </Button>
        </Group>
      </Modal>
      <Button variant="default" size="xs" onClick={open}>
        {t("button.delete")}
      </Button>
    </>
  );
}

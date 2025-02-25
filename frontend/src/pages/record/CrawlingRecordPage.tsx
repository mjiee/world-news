import { useState, useEffect } from "react";
import { useNavigate } from "react-router";
import { useTranslation } from "react-i18next";
import { Button, Container, Group, Table, Modal, LoadingOverlay, Pagination, Box } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { BackHeader } from "@/components/BackHeader";
import {
  queryCrawlingRecords,
  deleteCrawlingRecord,
  updateCrawlingRecordStatus,
  CrawlingRecord,
  CrawlingRecordType,
  CrawlingRecordStatus,
} from "@/services";
import { getPageNumber } from "@/utils/pagination";
import { httpx } from "wailsjs/go/models";

// crawling record page
export function CrawlingRecordPage() {
  return (
    <>
      <BackHeader />
      <Container size="md">
        <CrawlingRecords />
      </Container>
    </>
  );
}

// crawling records
function CrawlingRecords() {
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
  const recordTableBody = records.map((item, idx) => (
    <RecordTableBody key={idx} record={item} updatePage={updatePageHandler} />
  ));

  return (
    <>
      <Button variant="default" onClick={() => setLoading(true)}>
        {t("button.refresh")}
      </Button>
      {loading ? (
        <Box pos="relative">
          <LoadingOverlay visible={loading} zIndex={1000} overlayProps={{ radius: "sm", blur: 2 }} />
        </Box>
      ) : (
        <Table>
          <RecordTableHeader />
          <Table.Tbody>{recordTableBody}</Table.Tbody>
        </Table>
      )}
      <Pagination value={pagination.page} total={getPageNumber(pagination)} onChange={updatePageHandler} />
    </>
  );
}

function RecordTableHeader() {
  const { t } = useTranslation("record");

  return (
    <Table.Thead>
      <Table.Tr>
        <Table.Th>ID</Table.Th>
        <Table.Th>{t("table.head.date")}</Table.Th>
        <Table.Th>{t("table.head.record_type")}</Table.Th>
        <Table.Th>{t("table.head.quantity")}</Table.Th>
        <Table.Th>{t("table.head.status")}</Table.Th>
        <Table.Th />
      </Table.Tr>
    </Table.Thead>
  );
}

interface RecordTableBodyProps {
  record: CrawlingRecord;
  updatePage: (page: number) => void;
}

function RecordTableBody({ record, updatePage }: RecordTableBodyProps) {
  const { t } = useTranslation();

  return (
    <Table.Tr key={record.id}>
      <Table.Td>{record.id}</Table.Td>
      <Table.Td>{record.date}</Table.Td>
      <Table.Td>{t("table.body.record_type." + record.recordType, { ns: "record" })}</Table.Td>
      <Table.Td>{record.quantity}</Table.Td>
      <Table.Td>{t("table.body.status." + record.status, { ns: "record" })}</Table.Td>
      <Table.Td>
        <Button.Group>
          {viewRecordButton(record)}
          {pauseRecordButton(record)}
          <DeleteRecordButton record={record} updatePage={updatePage} />
        </Button.Group>
      </Table.Td>
    </Table.Tr>
  );
}

// view record button
const viewRecordButton = (record: CrawlingRecord) => {
  const { t } = useTranslation();
  const navigate = useNavigate();

  return record.recordType === CrawlingRecordType.CrawlingWebsite ? (
    <></>
  ) : (
    <Button variant="default" size="xs" onClick={() => navigate("/news/list/" + record.id)}>
      {t("button.view")}
    </Button>
  );
};

// pause record button
const pauseRecordButton = (record: CrawlingRecord) => {
  const { t } = useTranslation("record");

  return record.status !== CrawlingRecordStatus.ProcessingCrawlingRecord ? (
    <></>
  ) : (
    <Button
      variant="default"
      size="xs"
      onClick={() => updateCrawlingRecordStatus({ id: record.id, status: CrawlingRecordStatus.PausedCrawlingRecord })}
    >
      {t("button.pause")}
    </Button>
  );
};

// delete record button
interface DeleteRecordButtonProps {
  record: CrawlingRecord;
  updatePage: (page: number) => void;
}

function DeleteRecordButton({ record, updatePage }: DeleteRecordButtonProps) {
  const [opened, { open, close }] = useDisclosure(false);
  const { t } = useTranslation();

  // click ok handler
  const clickOkHandler = async () => {
    await deleteCrawlingRecord({ id: record.id });
    close();
    updatePage(1);
  };

  const deleteLabel = (date: String) => <p>{t("button.delete_label", { date, ns: "record" })}</p>;

  return record.status === CrawlingRecordStatus.ProcessingCrawlingRecord ? (
    <></>
  ) : (
    <>
      <Modal opened={opened} onClose={close} withCloseButton={false}>
        {deleteLabel(record.date)}
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

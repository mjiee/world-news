import { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router";
import { Container, Table, Button, Modal, Group, Pagination, LoadingOverlay } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { useTranslation } from "react-i18next";
import { BackHeader } from "@/components/BackHeader";
import { queryNews, deleteNews, NewsDetail } from "@/services";
import { httpx } from "wailsjs/go/models";
import { getPageNumber } from "@/utils/pagination";

// News list page
export function NewsListPage() {
  return (
    <>
      <BackHeader />
      <NewsTable />
    </>
  );
}

function NewsTable() {
  const { recordId } = useParams();
  const [newsList, setNewsList] = useState<NewsDetail[]>([]);
  const [pagination, setPagination] = useState<httpx.Pagination>({ page: 1, limit: 25, total: 0 });
  const [loading, setLoading] = useState<boolean>(true);

  // fetch news
  const fetchNews = async () => {
    if (!recordId || !loading) return;

    const resp = await queryNews({ recordId: Number(recordId), pagination: pagination });

    setLoading(false);

    if (!resp || !resp.data) return;

    setNewsList(resp.data);
  };

  useEffect(() => {
    fetchNews();
  }, [loading]);

  // update page
  const updatePageHandler = (page: number) => {
    setPagination({ ...pagination, page: page });
    setLoading(true);
  };

  // news table body
  const newsTableBody = newsList.map((news) => <NewsTableBody key={news.id} newsDetail={news} updatePage={updatePageHandler} />);

  return (
    <Container size="md">
      {loading ? (
        <LoadingOverlay visible={loading} />
      ) : (
        <Table>
          <NewsTableHeader />
          <Table.Tbody>{newsTableBody}</Table.Tbody>
        </Table>
      )}
      <Pagination value={pagination.page} total={getPageNumber(pagination)} onChange={updatePageHandler} />
    </Container>
  );
}

function NewsTableHeader() {
  const { t } = useTranslation("news");

  return (
    <Table.Thead>
      <Table.Tr>
        <Table.Th>ID</Table.Th>
        <Table.Th>{t("news_list.news_table.head.date")}</Table.Th>
        <Table.Th>{t("news_list.news_table.head.title")}</Table.Th>
        <Table.Th />
      </Table.Tr>
    </Table.Thead>
  );
}

interface NewsTableBodyProps {
  newsDetail: NewsDetail;
  updatePage: (page: number) => void;
}

function NewsTableBody({ newsDetail, updatePage }: NewsTableBodyProps) {
  const navigate = useNavigate();
  const { t } = useTranslation();

  return (
    <Table.Tr key={newsDetail.id}>
      <Table.Td>{newsDetail.id}</Table.Td>
      <Table.Td>{newsDetail.publishedAt}</Table.Td>
      <Table.Td>{newsDetail.title}</Table.Td>
      <Table.Td>
        <Button.Group>
          <Button variant="default" size="xs" onClick={() => navigate("/news/detail/" + newsDetail.id)}>
            {t("button.view")}
          </Button>
          <DeleteNewsButton newsId={newsDetail.id} updatePage={updatePage} />
        </Button.Group>
      </Table.Td>
    </Table.Tr>
  );
}

interface DeleteNewsButtonProps {
  newsId: number;
  updatePage: (page: number) => void;
}

function DeleteNewsButton({ newsId, updatePage }: DeleteNewsButtonProps) {
  const [opened, { open, close }] = useDisclosure(false);
  const { t } = useTranslation();

  // click ok handler
  const clickOkHandler = async () => {
    await deleteNews({ id: newsId });
    close();
    updatePage(1);
  };

  return (
    <>
      <Modal opened={opened} onClose={close} withCloseButton={false}>
        <p>{t("news_list.delete_label", { ns: "news" })}</p>
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

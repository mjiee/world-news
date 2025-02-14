import { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router";
import { Container, Table, Button, Modal, Group, Loader } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { useTranslation } from "react-i18next";
import { BackHeader } from "@/components/BackHeader";
import { queryNews, deleteNews, NewsDetail } from "@/services";

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
  const { t } = useTranslation("news");

  const tableHeader = (
    <Table.Tr>
      <Table.Th>ID</Table.Th>
      <Table.Th>{t("news_list.news_table.head.date")}</Table.Th>
      <Table.Th>{t("news_list.news_table.head.title")}</Table.Th>
      <Table.Th />
    </Table.Tr>
  );

  return (
    <>
      <Container size="md">
        <Table>
          <Table.Thead>{tableHeader}</Table.Thead>
          <Table.Tbody>
            <NewsTableBody />
          </Table.Tbody>
        </Table>
      </Container>
    </>
  );
}

function NewsTableBody() {
  const { recordId } = useParams();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [newsList, setNewsList] = useState<NewsDetail[]>([]);
  const [page, setPage] = useState<number>(1);
  const [loading, setLoading] = useState<boolean>(true);

  // fetch news
  const fetchNews = async () => {
    if (!recordId || !loading) return;

    const resp = await queryNews({ recordId: Number(recordId), pagination: { page: page, limit: 25 } });

    if (!resp || !resp.data) return;

    setLoading(false);
    setNewsList(resp.data);
  };

  useEffect(() => {
    fetchNews();
  }, [loading]);

  // update page
  const updatePageHandler = (page: number) => {
    setPage(page);
    setLoading(true);
  };

  return loading ? (
    <Loader color="blue" />
  ) : (
    <>
      {newsList.map((item) => (
        <Table.Tr key={item.id}>
          <Table.Td>{item.id}</Table.Td>
          <Table.Td>{item.publishedAt}</Table.Td>
          <Table.Td>{item.title}</Table.Td>
          <Table.Td>
            <Button.Group>
              <Button variant="default" size="xs" onClick={() => navigate("/news/detail/" + item.id)}>
                {t("button.view")}
              </Button>
              <DeleteNewsButton newsId={item.id} updatePage={updatePageHandler} />
            </Button.Group>
          </Table.Td>
        </Table.Tr>
      ))}
    </>
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

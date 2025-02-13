import { useNavigate } from "react-router";
import { Container, Table, Button, Modal, Group } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { useTranslation } from "react-i18next";
import { BackHeader } from "@/components/BackHeader";

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

const data = [
  {
    id: 1,
    date: "2024.01.02",
    title: "eftSection and rightSection allow adding icons or any other element to the left and right side of the button",
  },
  { id: 2, date: "2024.01.02", title: "Completed" },
];

function NewsTableBody() {
  let navigate = useNavigate();
  const { t } = useTranslation();

  const rows = data.map((item) => (
    <Table.Tr key={item.id}>
      <Table.Td>{item.id}</Table.Td>
      <Table.Td>{item.date}</Table.Td>
      <Table.Td>{item.title}</Table.Td>
      <Table.Td>
        <Button.Group>
          <Button variant="default" size="xs" onClick={() => navigate("/news/detail/" + item.id)}>
            {t("button.view")}
          </Button>
          <DeleteNewsButton newsId={item.id} />
        </Button.Group>
      </Table.Td>
    </Table.Tr>
  ));

  return <>{rows}</>;
}

interface DeleteNewsButtonProps {
  newsId: number;
}

function DeleteNewsButton({ newsId }: DeleteNewsButtonProps) {
  const [opened, { open, close }] = useDisclosure(false);
  const { t } = useTranslation();

  return (
    <>
      <Modal opened={opened} onClose={close} withCloseButton={false}>
        <p>{t("news_list.delete_label", { ns: "news" })}</p>
        <Group justify="flex-end">
          <Button onClick={close}>{t("button.ok")}</Button>
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

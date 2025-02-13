import { useNavigate } from "react-router";
import { Button, ActionIcon, Container, Avatar, Group, Table, Modal, TextInput } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { useForm } from "@mantine/form";
import { useTranslation } from "react-i18next";
import { LanguageSwitcher } from "@/components";
import styles from "@/assets/styles/header.module.css";

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
          {import.meta.env.VITE_PLATFORM}
          <FetchNewsButton />
          <Button onClick={() => navigate("/settings")}>{t("header.button.settings")}</Button>
          <LanguageSwitcher />
        </Group>
      </Container>
    </header>
  );
}

function FetchNewsButton() {
  const [opened, { open, close }] = useDisclosure(false);
  const { t } = useTranslation();

  const form = useForm({
    mode: "uncontrolled",
    initialValues: {
      startTime: "",
    },
  });

  return (
    <>
      <Modal opened={opened} onClose={close} withCloseButton={false}>
        <form onSubmit={form.onSubmit((values) => console.log(values))}>
          <TextInput
            label={t("header.label.start_time", { ns: "home" })}
            key={form.key("startTime")}
            {...form.getInputProps("startTime")}
          />
          <Group justify="flex-end" mt="md">
            <Button type="submit" onClick={close}>
              {t("button.ok")}
            </Button>
            <Button onClick={close} variant="default">
              {t("button.cancel")}
            </Button>
          </Group>
        </form>
      </Modal>
      <Button onClick={open}>{t("header.button.fetch_news", { ns: "home" })}</Button>
    </>
  );
}

const records = [
  { id: 1, date: "2024.01.04", quantity: 123, status: "processing" },
  { id: 2, date: "2024.01.10", quantity: 333, status: "completed" },
];

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
        <Table.Tbody>
          <RecordTableBody />
        </Table.Tbody>
      </Table>
    </Container>
  );
}

function RecordTableBody() {
  let navigate = useNavigate();
  const { t } = useTranslation();

  const rows = records.map((item) => (
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
          <DeleteRecordButton recordId={item.id} date={item.date} />
        </Button.Group>
      </Table.Td>
    </Table.Tr>
  ));

  return <>{rows}</>;
}

interface DeleteRecordButtonProps {
  recordId: number;
  date: String;
}

function DeleteRecordButton({ recordId, date }: DeleteRecordButtonProps) {
  const [opened, { open, close }] = useDisclosure(false);
  const { t } = useTranslation();

  return (
    <>
      <Modal opened={opened} onClose={close} withCloseButton={false}>
        <p>{t("crawling_records.button.delete_label", { date, ns: "home" })}</p>
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

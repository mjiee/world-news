import { useNavigate } from "react-router";
import { Button, Container, Avatar, Group, Table, Modal, TextInput } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { useForm } from "@mantine/form";
import styles from "@/styles/header.module.css";

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

  return (
    <header className={styles.header}>
      <Container size="md" className={styles.inner}>
        <Avatar size={28} name="World News" color="initials" />
        <Group gap={5}>
          <FetchNewsButton />
          <Button onClick={() => navigate("/settings")}>Settings</Button>
        </Group>
      </Container>
    </header>
  );
}

function FetchNewsButton() {
  const [opened, { open, close }] = useDisclosure(false);

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
          <TextInput label="Start Time" key={form.key("startTime")} {...form.getInputProps("startTime")} />
          <Group justify="flex-end" mt="md">
            <Button type="submit" onClick={close}>
              OK
            </Button>
            <Button onClick={close} variant="default">
              Cancel
            </Button>
          </Group>
        </form>
      </Modal>
      <Button onClick={open}>Fetch News</Button>
    </>
  );
}

const records = [
  { id: 1, date: "2024.01.04", quantity: 123, status: "Pending" },
  { id: 2, date: "2024.01.10", quantity: 333, status: "Completed" },
];

function CrawlingRecords() {
  const tableHeader = (
    <Table.Tr>
      <Table.Th>ID</Table.Th>
      <Table.Th>Date</Table.Th>
      <Table.Th>Quantity</Table.Th>
      <Table.Th>Status</Table.Th>
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

  const rows = records.map((item) => (
    <Table.Tr key={item.id}>
      <Table.Td>{item.id}</Table.Td>
      <Table.Td>{item.date}</Table.Td>
      <Table.Td>{item.quantity}</Table.Td>
      <Table.Td>{item.status}</Table.Td>
      <Table.Td>
        <Button.Group>
          <Button variant="default" size="xs" onClick={() => navigate("/news/list/" + item.id)}>
            view
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

  return (
    <>
      <Modal opened={opened} onClose={close} withCloseButton={false}>
        <p>Do you want to delete this record ({date})?</p>
        <Group justify="flex-end">
          <Button onClick={close}>OK</Button>
          <Button onClick={close} variant="default">
            Cancel
          </Button>
        </Group>
      </Modal>
      <Button variant="default" size="xs" onClick={open}>
        delete
      </Button>
    </>
  );
}

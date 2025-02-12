import { useNavigate } from "react-router";
import { Button, Container, Avatar, Group, Table } from "@mantine/core";
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
          <Button>Fetch News</Button>
          <Button onClick={() => navigate("/settings")}>Settings</Button>
        </Group>
      </Container>
    </header>
  );
}

const records = [
  { id: 1, date: "2024.01.02", quantity: 123, status: "Pending" },
  { id: 2, date: "2024.01.02", quantity: 333, status: "Completed" },
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
        <Table.Thead>{tableHeader} </Table.Thead>
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
          <Button variant="default" size="xs">
            delete
          </Button>
        </Button.Group>
      </Table.Td>
    </Table.Tr>
  ));

  return <>{rows}</>;
}

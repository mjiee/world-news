import { useNavigate } from "react-router";
import { Container, Table, Button } from "@mantine/core";
import { HeaderMenu } from "@/components/HeaderMenu";

// News list page
export function NewsListPage() {
  return (
    <>
      <HeaderMenu />
      <NewsTable />
    </>
  );
}

function NewsTable() {
  const tableHeader = (
    <Table.Tr>
      <Table.Th>ID</Table.Th>
      <Table.Th>Date</Table.Th>
      <Table.Th>Title</Table.Th>
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

  const rows = data.map((item) => (
    <Table.Tr key={item.id}>
      <Table.Td>{item.id}</Table.Td>
      <Table.Td>{item.date}</Table.Td>
      <Table.Td>{item.title}</Table.Td>
      <Table.Td>
        <Button.Group>
          <Button variant="default" size="xs" onClick={() => navigate("/news/detail/" + item.id)}>
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

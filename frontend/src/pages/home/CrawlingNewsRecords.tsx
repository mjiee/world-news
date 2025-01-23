import { Checkbox, Table } from "@mantine/core";
import { useState } from "react";

const records = [
  { id: 1, date: "2024.01.02", quantity: 123, status: "Pending" },
  { id: 2, date: "2024.01.02", quantity: 333, status: "Completed" },
];

export function CrawlingNewsRecords() {
  return (
    <Table>
      <RecordTableHeader />
      <RecordTableBody />
    </Table>
  );
}

function RecordTableHeader() {
  return (
    <Table.Thead>
      <Table.Tr>
        <Table.Th>ID</Table.Th>
        <Table.Th>Date</Table.Th>
        <Table.Th>Quantity</Table.Th>
        <Table.Th>Status</Table.Th>
        <Table.Th />
      </Table.Tr>
    </Table.Thead>
  );
}

function RecordTableBody() {
  const rows = records.map((item) => (
    <Table.Tr key={item.id}>
      <Table.Td>{item.id}</Table.Td>
      <Table.Td>{item.date}</Table.Td>
      <Table.Td>{item.quantity}</Table.Td>
      <Table.Td>{item.status}</Table.Td>
      <Table.Td>view delete</Table.Td>
    </Table.Tr>
  ));

  return <Table.Tbody>{rows}</Table.Tbody>;
}

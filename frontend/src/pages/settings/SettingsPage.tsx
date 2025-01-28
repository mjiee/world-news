import { useNavigate } from "react-router";
import { Container, Text, Flex, Space, Stack, Button, Title, Table, ActionIcon, Group, Pill } from "@mantine/core";
import { IconPencil, IconTrash } from "@tabler/icons-react";

// Application settings page
export function SettingsPage() {
  let navigate = useNavigate();

  return (
    <Stack h={300} bg="var(--mantine-color-body)" align="flex-start" justify="flex-start" gap="md">
      <Button onClick={() => navigate(-1)}>Return</Button>
      <AggregationWebsites />
      <NewsKeywords />
      <GlobalNewsWebsites />
    </Stack>
  );
}

const data = [
  { id: 1, url: "www.baidu.com" },
  { id: 2, url: "www.typescript.com" },
];

function AggregationWebsites() {
  return <WebsiteData title="News Website Aggregation" websiteType="AggregationWebsit" data={data} />;
}

function GlobalNewsWebsites() {
  return <WebsiteData title="Global News Websites" websiteType="NewsWebsite" data={data} />;
}

interface WebsiteDataProps {
  title: string;
  websiteType: string;
  data: {
    id: number;
    url: string;
  }[];
}

function WebsiteData(props: WebsiteDataProps) {
  const rows = props.data.map((item) => (
    <Flex key={item.id} gap="xl" justify="flex-start" align="center" direction="row" wrap="wrap">
      <Text>{item.url}</Text>
      <Space w="md" />
      <Group gap={0} justify="flex-end">
        <ActionIcon variant="subtle" color="gray">
          <IconPencil size={16} stroke={1.5} />
        </ActionIcon>
        <ActionIcon variant="subtle" color="red">
          <IconTrash size={16} stroke={1.5} />
        </ActionIcon>
      </Group>
    </Flex>
  ));

  return (
    <>
      <Title order={2}>{props.title}</Title>
      <Stack>{rows}</Stack>
    </>
  );
}

const keys = [
  { id: 1, key: "key1" },
  { id: 2, key: "key2" },
];

function NewsKeywords() {
  const pills = keys.map((keyData) => (
    <Pill key={keyData.id} withRemoveButton>
      {keyData.key}
    </Pill>
  ));

  return (
    <>
      <Title order={2}>News Keywords</Title>
      <Pill.Group>{pills}</Pill.Group>
    </>
  );
}

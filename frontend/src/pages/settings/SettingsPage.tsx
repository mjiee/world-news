import {
  Button,
  Container,
  Accordion,
  Text,
  Table,
  Space,
  Switch,
  Avatar,
  Group,
  Pill,
  ActionIcon,
  Flex,
  Stack,
  Modal,
  TextInput,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { IconPencil, IconPlus } from "@tabler/icons-react";
import { HeaderMenu } from "@/components/HeaderMenu";

const settingsItems = [
  {
    id: "topic",
    image: "https://img.icons8.com/?size=100&id=46893&format=png&color=000000",
    label: "News Topic",
    description: "Configure the news topics you are interested in to filter news.",
    content: <NewsTopics />,
  },
  {
    id: "collection",
    image: "https://img.icons8.com/?size=100&id=IwtVX5J92E9k&format=png&color=000000",
    label: "News Website Collection",
    description: "The news website collection is used to access global news websites.",
    content: <NewsWebsiteCollection />,
  },
  {
    id: "website",
    image: "https://img.icons8.com/?size=100&id=42835&format=png&color=000000",
    label: "News Website",
    description: "The URLs of global news websites obtained from the news website collection.",
    content: <NewsWebsite />,
  },
  {
    id: "service",
    image: "https://img.icons8.com/?size=100&id=104308&format=png&color=000000",
    label: "Remote service",
    description: "Enable cloud-deployed services for global news crawling, and use this feature when the local network is poor.",
    content: <RemoteService />,
  },
];

// Application settings page
export function SettingsPage() {
  const items = settingsItems.map((item) => (
    <Accordion.Item value={item.id} key={item.label}>
      <Accordion.Control>
        <SettingsLabel {...item} />
      </Accordion.Control>
      <Accordion.Panel>
        <Group wrap="nowrap">
          <Space w="54px" />
          {item.content}
        </Group>
      </Accordion.Panel>
    </Accordion.Item>
  ));

  return (
    <>
      <HeaderMenu />
      <Container size="md">
        <Accordion chevronPosition="right" variant="contained">
          {items}
        </Accordion>
      </Container>
    </>
  );
}

interface SettingsLabelProps {
  label: string;
  image: string;
  description: string;
}

function SettingsLabel({ label, image, description }: SettingsLabelProps) {
  return (
    <Group wrap="nowrap">
      <Avatar src={image} radius="xl" size="lg" />
      <div>
        <Text>{label}</Text>
        <Text size="sm" c="dimmed" fw={400}>
          {description}
        </Text>
      </div>
    </Group>
  );
}

const keys = [{ key: "key1" }, { key: "key2" }];

function NewsTopics() {
  const pills = keys.map((keyData) => (
    <Pill key={keyData.key} withRemoveButton size="lg">
      {keyData.key}
    </Pill>
  ));

  return (
    <>
      <Pill.Group>{pills}</Pill.Group>
      <AddNewsTopicButton />
    </>
  );
}

function AddNewsTopicButton() {
  const [opened, { open, close }] = useDisclosure(false);

  return (
    <>
      <ActionIcon variant="default" size="sm" onClick={open}>
        <IconPlus />
      </ActionIcon>
      <Modal title="News Topic" opened={opened} onClose={close} withCloseButton={false}>
        <TextInput />
        <Group justify="flex-end" mt="md">
          <Button type="submit" onClick={close}>
            OK
          </Button>
          <Button onClick={close} variant="default">
            Cancel
          </Button>
        </Group>
      </Modal>
    </>
  );
}

const data = [
  { url: "www.baidu.com", selectors: ["aa, cc", "bb"] },
  { url: "www.typescript.com", selectors: ["cc", "dd"] },
];

function NewsWebsiteCollection() {
  return <WebsiteTable />;
}

function NewsWebsite() {
  return (
    <Stack w={"100%"} align="stretch" justify="flex-start" gap="md">
      <Button variant="default">Update news website</Button>
      <WebsiteTable />
    </Stack>
  );
}

function WebsiteTable() {
  const tableHeader = (
    <Table.Tr>
      <Table.Th>Website</Table.Th>
      <Table.Th>Selector</Table.Th>
    </Table.Tr>
  );

  const tableBody = data.map((item) => (
    <Table.Tr key={item.url}>
      <Table.Td>{item.url}</Table.Td>
      <Table.Td>
        <Pill.Group>
          {item.selectors.map((value) => (
            <Pill key={value}>{value}</Pill>
          ))}
        </Pill.Group>
      </Table.Td>
    </Table.Tr>
  ));

  return (
    <Table>
      <Table.Thead>{tableHeader}</Table.Thead>
      <Table.Tbody>{tableBody}</Table.Tbody>
    </Table>
  );
}

function RemoteService() {
  return (
    <div>
      <Switch defaultChecked label="Enable remote services." />
      <Flex gap="lg" justify="flex-start" align="center" direction="row" wrap="wrap">
        <ActionIcon variant="default">
          <IconPencil />
        </ActionIcon>
        <p>
          https://127.0.0.1:8080 <span style={{ color: "var(--mantine-color-gray-5)" }}>(service host)</span>
        </p>
      </Flex>
    </div>
  );
}

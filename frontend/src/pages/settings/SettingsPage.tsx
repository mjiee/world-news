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
import { useTranslation } from "react-i18next";
import { BackHeader } from "@/components/BackHeader";

const settingsItems = [
  {
    id: "topic",
    image: "https://img.icons8.com/?size=100&id=46893&format=png&color=000000",
    content: <NewsTopics />,
  },
  {
    id: "collection",
    image: "https://img.icons8.com/?size=100&id=IwtVX5J92E9k&format=png&color=000000",
    content: <NewsWebsiteCollection />,
  },
  {
    id: "website",
    image: "https://img.icons8.com/?size=100&id=42835&format=png&color=000000",
    content: <NewsWebsite />,
  },
  {
    id: "service",
    image: "https://img.icons8.com/?size=100&id=104308&format=png&color=000000",
    content: <RemoteService />,
  },
];

// Application settings page
export function SettingsPage() {
  const items = settingsItems.map((item) => (
    <Accordion.Item key={item.id} value={item.id}>
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
      <BackHeader />
      <Container size="md">
        <Accordion chevronPosition="right" variant="contained">
          {items}
        </Accordion>
      </Container>
    </>
  );
}

interface SettingsLabelProps {
  id: string;
  image: string;
}

function SettingsLabel({ id, image }: SettingsLabelProps) {
  const { t } = useTranslation("settings");

  return (
    <Group wrap="nowrap">
      <Avatar src={image} radius="xl" size="lg" />
      <div>
        <Text>{t("settings_label." + id + ".label")}</Text>
        <Text size="sm" c="dimmed" fw={400}>
          {t("settings_label." + id + ".description")}
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
  const { t } = useTranslation();

  return (
    <>
      <ActionIcon variant="default" size="sm" onClick={open}>
        <IconPlus />
      </ActionIcon>
      <Modal title={t("news_topic.title", { ns: "settings" })} opened={opened} onClose={close} withCloseButton={false}>
        <TextInput />
        <Group justify="flex-end" mt="md">
          <Button type="submit" onClick={close}>
            {t("button.ok")}
          </Button>
          <Button onClick={close} variant="default">
            {t("button.cancel")}
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
  const { t } = useTranslation("settings");

  return (
    <Stack w={"100%"} align="stretch" justify="flex-start" gap="md">
      <Button variant="default">{t("news_website.button.update_news_website")}</Button>
      <WebsiteTable />
    </Stack>
  );
}

function WebsiteTable() {
  const { t } = useTranslation("settings");

  const tableHeader = (
    <Table.Tr>
      <Table.Th>{t("news_website.table.head.website")}</Table.Th>
      <Table.Th>{t("news_website.table.head.selector")}</Table.Th>
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
  const { t } = useTranslation("settings");

  return (
    <div>
      <Switch defaultChecked label={t("remote_service.lable.enable_remote_service")} />
      <Flex gap="lg" justify="flex-start" align="center" direction="row" wrap="wrap">
        <ActionIcon variant="default">
          <IconPencil />
        </ActionIcon>
        <p>
          https://127.0.0.1:8080
          <span style={{ color: "var(--mantine-color-gray-5)" }}>({t("remote_service.lable.service_host")})</span>
        </p>
      </Flex>
    </div>
  );
}

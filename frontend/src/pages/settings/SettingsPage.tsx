import { Container, Accordion, Text, Space, Switch, Avatar, Group, ActionIcon, Flex } from "@mantine/core";
import { IconPencil } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";
import { BackHeader } from "@/components/BackHeader";
import { NewsTopics } from "./NewsTopics";
import { NewsWebsite, NewsWebsiteCollection } from "./NewsWebsite";
import { isWeb } from "@/utils/platform";

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

// settings item label
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

// remote service settings

function RemoteService() {
  const { t } = useTranslation("settings");

  return isWeb() ? (
    <p>{t("remote_service.web_not_support")}</p>
  ) : (
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

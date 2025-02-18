import { Container, Accordion, Text, Space, Switch, Avatar, Group, ActionIcon, Flex } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { BackHeader } from "@/components/BackHeader";
import { NewsTopics } from "./NewsTopics";
import { NewsWebsite, NewsWebsiteCollection } from "./NewsWebsite";
import { isWeb } from "@/utils/platform";
import IconPencil from "@/assets/icons/IconPencil.svg";
import settingsTopic from "@/assets/images/settings_topic.png";
import settingsService from "@/assets/images/settings_service.png";
import settingsCollection from "@/assets/images/settings_collection.png";
import settingsWebsite from "@/assets/images/settings_website.png";

const settingsItems = [
  { id: "topic", image: settingsTopic, content: <NewsTopics /> },
  { id: "collection", image: settingsCollection, content: <NewsWebsiteCollection /> },
  { id: "website", image: settingsWebsite, content: <NewsWebsite /> },
  { id: "service", image: settingsService, content: <RemoteService /> },
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
          <Avatar src={IconPencil} alt="eidt" variant="default" radius="sm" size="sm" />
        </ActionIcon>
        <p>
          https://127.0.0.1:8080
          <span style={{ color: "var(--mantine-color-gray-5)" }}>({t("remote_service.lable.service_host")})</span>
        </p>
      </Flex>
    </div>
  );
}

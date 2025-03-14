import { Accordion, Text, Avatar, Group } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { NewsTopics } from "./NewsTopics";
import { NewsWebsite, NewsWebsiteCollection } from "./NewsWebsite";
import { RemoteService } from "./RemoteService";
import { NewsCritique } from "./NewsCritique";
import settingsTopic from "@/assets/images/settings_topic.png";
import settingsService from "@/assets/images/settings_service.png";
import settingsCollection from "@/assets/images/settings_collection.png";
import settingsWebsite from "@/assets/images/settings_website.png";
import settingsCritique from "@/assets/images/settings_critique.png";

const settingsItems = [
  { id: "topic", image: settingsTopic, content: <NewsTopics /> },
  { id: "collection", image: settingsCollection, content: <NewsWebsiteCollection /> },
  { id: "website", image: settingsWebsite, content: <NewsWebsite /> },
  { id: "service", image: settingsService, content: <RemoteService /> },
  { id: "critique", image: settingsCritique, content: <NewsCritique /> },
];

// Application settings page
export function SettingsPage() {
  return (
    <Accordion chevronPosition="right" variant="contained">
      {settingsItems.map((item) => (
        <Accordion.Item key={item.id} value={item.id}>
          <Accordion.Control>
            <SettingsLabel {...item} />
          </Accordion.Control>
          <Accordion.Panel mx="md">{item.content}</Accordion.Panel>
        </Accordion.Item>
      ))}
    </Accordion>
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

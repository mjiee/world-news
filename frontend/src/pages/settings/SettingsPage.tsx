import { useState } from "react";
import { Accordion, Text, Avatar, Group } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { NewsTopics, topicKey } from "./NewsTopics";
import { NewsWebsite, NewsWebsiteCollection, collectionKey, websiteKey } from "./NewsWebsite";
import { RemoteService, serviceKey } from "./RemoteService";
import { NewsCritique, critiqueKey } from "./NewsCritique";
import { NewsTranslation, translateKey } from "./NewsTranslation";
import settingsTopic from "@/assets/images/settings_topic.png";
import settingsService from "@/assets/images/settings_service.png";
import settingsCollection from "@/assets/images/settings_collection.png";
import settingsWebsite from "@/assets/images/settings_website.png";
import settingsCritique from "@/assets/images/settings_critique.png";
import settingsTranslate from "@/assets/images/settings_translate.png";

const settingsItems = [
  { id: topicKey, image: settingsTopic, content: (value: string | null) => <NewsTopics item={value} /> },
  {
    id: collectionKey,
    image: settingsCollection,
    content: (value: string | null) => <NewsWebsiteCollection item={value} />,
  },
  { id: websiteKey, image: settingsWebsite, content: (value: string | null) => <NewsWebsite item={value} /> },
  { id: serviceKey, image: settingsService, content: (_: string | null) => <RemoteService /> },
  { id: critiqueKey, image: settingsCritique, content: (value: string | null) => <NewsCritique item={value} /> },
  { id: translateKey, image: settingsTranslate, content: (value: string | null) => <NewsTranslation item={value} /> },
];

// Application settings page
export function SettingsPage() {
  const [value, setValue] = useState<string | null>(null);

  return (
    <Accordion chevronPosition="right" variant="contained" value={value} onChange={setValue}>
      {settingsItems.map((item) => (
        <Accordion.Item key={item.id} value={item.id}>
          <Accordion.Control>
            <SettingsLabel {...item} />
          </Accordion.Control>
          <Accordion.Panel mx="md">{item.content(value)}</Accordion.Panel>
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

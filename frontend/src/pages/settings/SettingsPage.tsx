import { Box, Card, Group, Modal, SimpleGrid, Stack, Text } from "@mantine/core";
import {
  IconChevronRight,
  IconFileText,
  IconGlobe,
  IconLanguage,
  IconMessage,
  IconMoodEdit,
  IconRobot,
  IconScript,
  IconServer,
  IconVolume,
} from "@tabler/icons-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { CritiquePrompt, PodcastPrompt } from "./settings/AiPrompt";
import NewsTopics from "./settings/NewsTopics";
import NewsTranslation from "./settings/NewsTranslation";
import { NewsWebsite, NewsWebsiteCollection } from "./settings/NewsWebsite";
import RemoteService from "./settings/RemoteService";
import TextAi from "./settings/TextAi";
import TextToSpeechAi from "./settings/TextToSpeechAi";
import styles from "./styles/settings.module.css";

interface SettingItemConfig {
  key: string;
  icon: React.ComponentType<any>;
  color: string;
  component: React.ComponentType;
  modalSize?: "sm" | "md" | "lg" | "xl";
}

const settingsItems: SettingItemConfig[] = [
  { key: "topic", icon: IconMessage, color: "#FF6B6B", component: NewsTopics },
  { key: "collection", icon: IconGlobe, color: "#4ECDC4", component: NewsWebsiteCollection },
  { key: "website", icon: IconFileText, color: "#45B7D1", component: NewsWebsite, modalSize: "xl" },
  { key: "service", icon: IconServer, color: "#96CEB4", component: RemoteService },
  { key: "text_ai", icon: IconRobot, color: "#FFEAA7", component: TextAi },
  { key: "text_to_speech_ai", icon: IconVolume, color: "#DFE6E9", component: TextToSpeechAi, modalSize: "xl" },
  { key: "critique_prompt", icon: IconMoodEdit, color: "#A29BFE", component: CritiquePrompt },
  { key: "podcast_prompt", icon: IconScript, color: "#FD79A8", component: PodcastPrompt, modalSize: "xl" },
  { key: "translate", icon: IconLanguage, color: "#74B9FF", component: NewsTranslation },
];

export function SettingsPage() {
  const { t } = useTranslation("settings");
  const [activeModal, setActiveModal] = useState<string | null>(null);

  const openModal = (key: string) => setActiveModal(key);
  const closeModal = () => setActiveModal(null);

  const activeItem = settingsItems.find((item) => item.key === activeModal);

  return (
    <Stack gap="xl">
      <SimpleGrid cols={{ base: 1, sm: 2, lg: 3 }} spacing="lg">
        {settingsItems.map((item) => (
          <SettingCard key={item.key} item={item} onClick={() => openModal(item.key)} />
        ))}
      </SimpleGrid>

      {activeItem && (
        <SettingModal
          opened={!!activeModal}
          onClose={closeModal}
          title={t(`settings_label.${activeItem.key}.label`)}
          size={activeItem.modalSize || "lg"}
        >
          <activeItem.component />
        </SettingModal>
      )}
    </Stack>
  );
}

// Setting Card Component
function SettingCard({ item, onClick }: { item: SettingItemConfig; onClick: () => void }) {
  const { t } = useTranslation("settings");
  const Icon = item.icon;

  return (
    <Card className={styles.settingCard} onClick={onClick} shadow="sm" padding="lg" radius="md" withBorder>
      <Group justify="space-between" wrap="nowrap">
        <Group gap="md" wrap="nowrap">
          <Box className={styles.iconWrapper} style={{ backgroundColor: item.color }}>
            <Icon size={24} stroke={2} color="white" />
          </Box>
          <Box style={{ flex: 1 }}>
            <Text fw={600} size="md">
              {t("settings_label." + item.key + ".label")}
            </Text>
            <Text size="xs" c="dimmed" lineClamp={2}>
              {t("settings_label." + item.key + ".description")}
            </Text>
          </Box>
        </Group>
        <IconChevronRight size={20} className={styles.chevron} />
      </Group>
    </Card>
  );
}

// Setting Modal Component

interface SettingModalProps {
  opened: boolean;
  onClose: () => void;
  title: string;
  children: React.ReactNode;
  size?: string;
}

function SettingModal({ opened, onClose, title, children, size = "lg" }: SettingModalProps) {
  return (
    <Modal
      opened={opened}
      onClose={onClose}
      title={
        <Text fw={600} size="lg">
          {title}
        </Text>
      }
      size={size}
      centered
      classNames={{
        content: styles.modalContent,
        header: styles.modalHeader,
      }}
    >
      {children}
    </Modal>
  );
}

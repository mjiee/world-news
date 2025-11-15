import { ActionIcon, Box, Divider, Drawer, Group, Stack, Title, Tooltip } from "@mantine/core";
import { useDisclosure, useMediaQuery } from "@mantine/hooks";
import {
  IconChevronLeft,
  IconChevronRight,
  IconHistory,
  IconHome,
  IconList,
  IconMenu2,
  IconSettings,
  IconStar,
} from "@tabler/icons-react";
import { useTranslation } from "react-i18next";
import { Outlet, useLocation, useNavigate } from "react-router";

import appicon from "@/assets/images/appicon.png";
import styles from "@/assets/styles/appLayout.module.css";
import { LanguageSwitcher } from "@/components";
import { buildAudioSrc } from "@/stores";
import { HeaderAudioPlayer, SidebarAudioPlayer, useAudioPlayer } from "./AudioPlayer";

export function AppLayout() {
  const [mobileOpened, { toggle: toggleMobile, close: closeMobile }] = useDisclosure();
  const [collapsed, { toggle: toggleCollapsed }] = useDisclosure(false);
  const isMobile = useMediaQuery("(max-width: 768px)");
  const audioData = useAudioPlayer();
  const { audioRef, currentAudio } = audioData;

  return (
    <div className={styles.layout}>
      {isMobile ? (
        <>
          <Header onMenuClick={toggleMobile} audioData={audioData} />
          <Drawer opened={mobileOpened} onClose={closeMobile} size="280px" padding={0} withCloseButton={false}>
            <Sidebar collapsed={false} onToggle={toggleCollapsed} isMobile={true} audioData={audioData} />
          </Drawer>
        </>
      ) : (
        <aside className={`${styles.aside} ${collapsed ? styles.collapsed : ""}`}>
          <Sidebar collapsed={collapsed} onToggle={toggleCollapsed} isMobile={false} audioData={audioData} />
        </aside>
      )}

      {currentAudio && (
        <audio ref={audioRef} src={buildAudioSrc(currentAudio.audio)} preload="auto" style={{ display: "none" }} />
      )}

      <main className={`${styles.main} ${collapsed ? styles.expanded : ""} ${isMobile ? styles.mobile : ""}`}>
        <Outlet />
      </main>
    </div>
  );
}

interface NavItemProps {
  label: string;
  icon: React.ComponentType<any>;
  active: boolean;
  collapsed: boolean;
  onClick: () => void;
}

function NavItem({ label, icon: Icon, active, collapsed, onClick }: NavItemProps) {
  const content = (
    <button className={`${styles.navItem} ${active ? styles.active : ""}`} onClick={onClick}>
      <Icon size={22} stroke={1.5} />
      {!collapsed && <span>{label}</span>}
    </button>
  );

  return collapsed ? (
    <Tooltip label={label} position="right" withArrow>
      {content}
    </Tooltip>
  ) : (
    content
  );
}

interface SidebarProps {
  collapsed: boolean;
  onToggle: () => void;
  isMobile: boolean;
  audioData: ReturnType<typeof useAudioPlayer>;
}

function Sidebar({ collapsed, onToggle, isMobile, audioData }: SidebarProps) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();

  const navItems = [
    { label: t("navbar.news"), icon: IconHome, path: "/" },
    { label: t("navbar.favorites"), icon: IconStar, path: "/news/favorites" },
    { label: t("navbar.tasks"), icon: IconList, path: "/tasks" },
    { label: t("navbar.records"), icon: IconHistory, path: "/records" },
    { label: t("navbar.settings"), icon: IconSettings, path: "/settings" },
  ];

  return (
    <Stack className={styles.sidebar} gap={0}>
      <Box className={styles.logo}>
        <img src={appicon} alt="Logo" />
        {!collapsed && <Title order={4}>World News</Title>}
      </Box>

      <Divider my="md" />

      <Stack className={styles.nav} gap={4}>
        {navItems.map((item) => (
          <NavItem
            key={item.path}
            label={item.label}
            icon={item.icon}
            active={location.pathname === item.path}
            collapsed={collapsed && !isMobile}
            onClick={() => navigate(item.path)}
          />
        ))}
      </Stack>

      <Divider my="md" />

      {<SidebarAudioPlayer audioData={audioData} collapsed={collapsed} />}

      <Divider my="md" />

      <Group justify={collapsed ? "center" : "space-between"}>
        {!collapsed && <LanguageSwitcher />}
        {!isMobile && (
          <ActionIcon onClick={onToggle} variant="subtle" size="lg">
            {collapsed ? <IconChevronRight size={20} /> : <IconChevronLeft size={20} />}
          </ActionIcon>
        )}
      </Group>
    </Stack>
  );
}

function Header({ onMenuClick, audioData }: { onMenuClick: () => void; audioData: ReturnType<typeof useAudioPlayer> }) {
  return (
    <header className={styles.header}>
      <Group gap="sm" wrap="nowrap">
        <ActionIcon onClick={onMenuClick} size="lg" variant="subtle">
          <IconMenu2 size={24} />
        </ActionIcon>
        <img src={appicon} alt="Logo" className={styles.logoSmall} />
        <Title order={4}>World News</Title>
        <HeaderAudioPlayer audioData={audioData} />
      </Group>
    </header>
  );
}

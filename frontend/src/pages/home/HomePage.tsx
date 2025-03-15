import { Outlet, useNavigate } from "react-router";
import { useTranslation } from "react-i18next";
import { Flex, Group, Avatar, Title, Text, Button, Menu, Burger } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { LanguageSwitcher } from "@/components";
import styles from "@/assets/styles/header.module.css";
import appicon from "@/assets/images/appicon.png";

// page layout
export function HomePage() {
  return (
    <>
      <header className={styles.header}>
        <Group justify="space-between">
          <Flex align="center">
            <Avatar pr="10" src={appicon} variant="default" radius="sm" />
            <Title c="white" fw={700} order={2} visibleFrom="xs">
              World News
            </Title>
          </Flex>

          <Group>
            <DesktopNav />
            <MobileNav />
            <LanguageSwitcher />
          </Group>
        </Group>
      </header>
      <div className={styles.container}>
        <Outlet />
      </div>
    </>
  );
}

// desktop nav
function DesktopNav() {
  const { t } = useTranslation("home");
  const navigate = useNavigate();

  const navText = (title: string) => (
    <Text c="white" size="xl" fw={700}>
      {t(title)}
    </Text>
  );

  return (
    <Button.Group visibleFrom="sm">
      <Button onClick={() => navigate("/")} variant="subtle">
        {navText("header.button.home")}
      </Button>
      <Button onClick={() => navigate("/records")} variant="subtle">
        {navText("header.button.records")}
      </Button>
      <Button onClick={() => navigate("/settings")} variant="subtle">
        {navText("header.button.settings")}
      </Button>
    </Button.Group>
  );
}

// mobile nav
function MobileNav() {
  const [opened, { toggle }] = useDisclosure();
  const { t } = useTranslation("home");
  const navigate = useNavigate();

  const onClickHandle = (route: string) => {
    navigate(route);
    toggle();
  };

  const navText = (title: string) => (
    <Text size="md" fw={500}>
      {t(title)}
    </Text>
  );

  return (
    <Menu position="bottom-end">
      <Menu.Target>
        <Burger opened={opened} onClick={toggle} lineSize={3} size="md" hiddenFrom="sm" color="white" />
      </Menu.Target>

      <Menu.Dropdown>
        <Menu.Item onClick={() => onClickHandle("/")}>{navText("header.button.home")}</Menu.Item>
        <Menu.Item onClick={() => onClickHandle("/records")}>{navText("header.button.records")}</Menu.Item>
        <Menu.Item onClick={() => onClickHandle("/settings")}>{navText("header.button.settings")}</Menu.Item>
      </Menu.Dropdown>
    </Menu>
  );
}

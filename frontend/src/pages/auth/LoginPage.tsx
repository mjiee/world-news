import { Link } from "react-router";
import { useTranslation } from "react-i18next";
import { Avatar, Button, Container, Paper, TextInput, Title } from "@mantine/core";
import { LanguageSwitcher } from "@/components";
import { useRemoteServiceStore } from "@/stores";
import styles from "@/assets/styles/header.module.css";
import appicon from "@/assets/images/appicon.png";
import classes from "./styles/auth.module.css";

// login page
export function LoginPage() {
  const { t } = useTranslation();
  const token = useRemoteServiceStore((state) => state.token);
  const setToken = useRemoteServiceStore((state) => state.setToken);

  return (
    <>
      <header className={styles.header}>
        <Container size="md" className={styles.inner}>
          <Avatar src={appicon} variant="default" radius="sm" />
          <LanguageSwitcher />
        </Container>
      </header>
      <Container size={460} my={30}>
        <Title className={classes.title} ta="center">
          {t("auth.title")}
        </Title>

        <Paper withBorder shadow="md" p={30} radius="md" mt="xl">
          <TextInput
            label={t("auth.label")}
            placeholder="0123456"
            required
            value={token}
            onChange={(event) => setToken(event.currentTarget.value)}
          />
          <Link to="/">
            <Button fullWidth mt="xl">
              {t("auth.button")}
            </Button>
          </Link>
        </Paper>
      </Container>
    </>
  );
}

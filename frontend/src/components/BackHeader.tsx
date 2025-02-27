import { useNavigate } from "react-router";
import { useTranslation } from "react-i18next";
import { Container, Button, Avatar } from "@mantine/core";
import styles from "@/assets/styles/header.module.css";
import appicon from "@/assets/images/appicon.png";

export function BackHeader() {
  let navigate = useNavigate();
  const { t } = useTranslation("common");

  return (
    <header className={styles.header}>
      <Container size="md" className={styles.inner}>
        <Avatar onClick={() => navigate("/")} src={appicon} variant="default" radius="sm" />
        <Button variant="default" onClick={() => navigate(-1)}>
          {t("button.back")}
        </Button>
      </Container>
    </header>
  );
}

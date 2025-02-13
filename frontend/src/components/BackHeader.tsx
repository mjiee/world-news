import { useNavigate } from "react-router";
import { useTranslation } from "react-i18next";
import { Container, Button, Avatar } from "@mantine/core";
import styles from "@/assets/styles/header.module.css";

export function BackHeader() {
  let navigate = useNavigate();
  const { t } = useTranslation("common");

  return (
    <header className={styles.header}>
      <Container size="md" className={styles.inner}>
        <Avatar size={28} name="World News" color="initials" />
        <Button variant="default" onClick={() => navigate(-1)}>
          {t("button.back")}
        </Button>
      </Container>
    </header>
  );
}

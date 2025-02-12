import { useNavigate } from "react-router";
import { Container, Button, Avatar } from "@mantine/core";
import styles from "@/styles/header.module.css";

export function HeaderMenu() {
  let navigate = useNavigate();

  return (
    <header className={styles.header}>
      <Container size="md" className={styles.inner}>
        <Avatar size={28} name="World News" color="initials" />
        <Button variant="default" onClick={() => navigate(-1)}>
          Return
        </Button>
      </Container>
    </header>
  );
}

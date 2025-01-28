import { useNavigate } from "react-router";
import { Container, Avatar, Group, Button } from "@mantine/core";
import classes from "./styles/header.module.css";

export function HomeHeader() {
  let navigate = useNavigate();

  return (
    <header className={classes.header}>
      <Container size="md" className={classes.inner}>
        <Avatar size={28} name="World News" color="initials" />
        <Group gap={5}>
          <Button onClick={() => navigate("/settings")}>Settings</Button>
          <Button>Fetch News</Button>
        </Group>
      </Container>
    </header>
  );
}

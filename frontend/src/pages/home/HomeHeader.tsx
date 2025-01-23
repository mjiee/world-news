import { Container, Avatar, Group, Button } from "@mantine/core";
import classes from "./styles/header.module.css";

export function HomeHeader() {
  return (
    <header className={classes.header}>
      <Container size="md" className={classes.inner}>
        <Avatar size={28} name="World News" color="initials" />
        <Group gap={5}>
          <Button>Settings</Button>
          <Button>Fetch News</Button>
        </Group>
      </Container>
    </header>
  );
}

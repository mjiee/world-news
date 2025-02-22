import { Container } from "@mantine/core";
import { NewsList } from "@/pages/news/NewsList";
import { HomeHeader } from "./HomeHeader";

// Application homepage
export function HomePage() {
  return (
    <>
      <HomeHeader />
      <Container size="md">
        <NewsList recordId={0} />
      </Container>
    </>
  );
}

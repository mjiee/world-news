import { useParams } from "react-router";
import { Container } from "@mantine/core";
import { BackHeader } from "@/components/BackHeader";
import { NewsList } from "./NewsList";

// News list page
export function NewsListPage() {
  const { recordId } = useParams();

  return (
    <>
      <BackHeader />
      <Container size="md">
        <NewsList recordId={parseInt(recordId ?? "0")} />
      </Container>
    </>
  );
}

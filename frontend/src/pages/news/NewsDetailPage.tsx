import { Container, Title, ActionIcon, CopyButton, Text, Flex } from "@mantine/core";
import { IconCopy, IconCheck } from "@tabler/icons-react";
import { BackHeader } from "@/components/BackHeader";

const data = {
  title: "Example article.",
  link: "https://mantine.dev/",
  contents: [
    "Let’s talk for a moment about how we talk about our teams.",
    "There can be a perception that as a manager of an organization you are in control at all times. Part of that control can invariably be perceived as how you appear to be in charge, are competent, or how you personally perform.",
    "function that should be called to copy given value to clipboard",
  ],
  publishedAt: "2024.01.02",
  images: [
    "https://images.unsplash.com/photo-1508193638397-1c4234db14d8?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=400&q=80",
    "https://images.unsplash.com/photo-1559494007-9f5847c49d94?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=400&q=80",
    "https://images.unsplash.com/photo-1510798831971-661eb04b3739?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=400&q=80",
  ],
};

// News detail page
export function NewsDetailPage() {
  const content = data.contents.map((item, idx) => <p key={idx}>{item}</p>);
  const image = data.images.map((item, idx) => <img key={idx} src={item} alt="news" />);

  return (
    <>
      <BackHeader />
      <Container size="md">
        <Title>{data.title}</Title>
        <p style={{ color: "var(--mantine-color-gray-5)" }}>{data.publishedAt}</p>
        <NewsLink link={data.link} />
        {content}
        {image}
      </Container>
    </>
  );
}

interface NewsLinkProps {
  link: string;
}

function NewsLink({ link }: NewsLinkProps) {
  return (
    <Flex>
      <Text c="blue">{link}</Text>
      <CopyButton value={link} timeout={2000}>
        {({ copied, copy }) => (
          <ActionIcon color={copied ? "teal" : "gray"} variant="subtle" onClick={copy}>
            {copied ? <IconCheck size={16} /> : <IconCopy size={16} />}
          </ActionIcon>
        )}
      </CopyButton>
    </Flex>
  );
}

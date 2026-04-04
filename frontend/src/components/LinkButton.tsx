import { ActionIcon, Anchor, CopyButton, Flex, Text } from "@mantine/core";
import { IconCheck, IconCopy } from "@tabler/icons-react";

// link button
interface LinkButtonProps {
  link: string;
  label?: string;
}

export function LinkButton({ link, label }: LinkButtonProps) {
  return (
    <Flex>
      {label && (
        <Text c="dimmed" size="sm">
          {label}
        </Text>
      )}
      <Anchor href={link} target="_blank" rel="noopener noreferrer" size="sm" truncate style={{ maxWidth: "500px" }}>
        {link}
      </Anchor>
      <CopyButton value={link} timeout={2000}>
        {({ copied, copy }) => (
          <ActionIcon p={3} color={copied ? "teal" : "gray"} variant="subtle" size="sm" onClick={copy}>
            {copied ? <IconCheck /> : <IconCopy />}
          </ActionIcon>
        )}
      </CopyButton>
    </Flex>
  );
}

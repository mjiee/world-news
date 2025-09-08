import { ActionIcon, Badge, Popover, Stack, Tooltip } from "@mantine/core";
import IconChevronUp from "@/assets/icons/IconChevronUp.svg?react";
import IconChevronDown from "@/assets/icons/IconChevronDown.svg?react";
import { saveWebsiteWeight } from "@/services";
import { useState } from "react";
import { useTranslation } from "react-i18next";

interface SourceLabelProps {
  source: string;
  size?: "sm" | "md";
}

export function SourceLabel({ source, size = "md" }: SourceLabelProps) {
  const { t } = useTranslation("settings");
  const [loading, setLoading] = useState<boolean>(false);

  const saveSourceWeight = async (step: number) => {
    setLoading(true);
    await saveWebsiteWeight({ website: source, step: step });
    setLoading(false);
  };

  const sourceLable = (opt: string) => (
    <Tooltip label={t("source_label." + opt)} position="right" withArrow>
      <ActionIcon
        variant="light"
        size="md"
        disabled={loading}
        onClick={() => saveSourceWeight(opt === "incr" ? 1 : -1)}
      >
        {opt === "incr" ? <IconChevronUp /> : <IconChevronDown />}
      </ActionIcon>
    </Tooltip>
  );

  return (
    <Popover position="right" withArrow>
      <Popover.Target>
        <Badge variant="light" size={size} color="blue">
          {source}
        </Badge>
      </Popover.Target>
      <Popover.Dropdown>
        <Stack align="center" justify="flex-start" gap="sm">
          {sourceLable("incr")}
          {sourceLable("decr")}
        </Stack>
      </Popover.Dropdown>
    </Popover>
  );
}

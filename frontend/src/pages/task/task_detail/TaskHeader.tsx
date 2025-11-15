import { PodcastTask, PodcastTaskResult, TaskStageStatus } from "@/services";
import { ActionIcon, Badge, Box, Flex, Group, Progress, Stack, Text, Title, Tooltip } from "@mantine/core";
import { IconRefresh } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";
import classes from "../styles/taskDetail.module.css";

interface TaskHeaderProps {
  task: PodcastTask;
  loading: boolean;
  onRefresh: () => void;
}

export default function TaskHeader({ task, loading, onRefresh }: TaskHeaderProps) {
  const { t } = useTranslation("task");
  const processing = task?.stages?.some((s) => s.status === TaskStageStatus.Processing);

  const getStatusColor = () => {
    if (task?.result === PodcastTaskResult.Failed) return "red";
    if (processing) return "blue";
    return "green";
  };

  const getProgressValue = () => {
    if (!task?.stages) return 0;
    const completed = task.stages.filter((s) => s.status === TaskStageStatus.Completed).length;
    return (completed / task.stages.length) * 100;
  };

  return (
    <Stack gap="md" mb="xl">
      <Flex justify="space-between" align="flex-start" wrap="wrap" gap="md">
        <Box style={{ flex: 1, minWidth: "250px" }}>
          <Title order={1} size="h2" mb="xs" className={classes.headertitle}>
            {task?.title || "Podcast Task"}
          </Title>
          <Group gap="xs">
            <Badge variant="gradient" gradient={{ from: "violet", to: "grape" }} size="lg">
              {task.batchNo}
            </Badge>
            <Badge color={getStatusColor()} size="lg" variant="light">
              {processing ? t("podcast.status.processing") : t("podcast.result." + task.result)}
            </Badge>
          </Group>
        </Box>
        <Tooltip label="Refresh">
          <ActionIcon
            variant="gradient"
            gradient={{ from: "violet", to: "grape" }}
            size="lg"
            onClick={onRefresh}
            loading={loading}
          >
            <IconRefresh size={18} />
          </ActionIcon>
        </Tooltip>
      </Flex>

      <Box>
        <Group justify="space-between" mb="xs">
          <Text size="sm" fw={500} c="dimmed">
            {t("podcast.progress")}
          </Text>
          <Text size="sm" fw={600} c="violet">
            {Math.round(getProgressValue())}%
          </Text>
        </Group>
        <Progress value={getProgressValue()} size="lg" radius="xl" animated={processing} className={classes.progress} />
      </Box>
    </Stack>
  );
}

import { TaskStage, TaskStageName, TaskStageStatus } from "@/services";
import { md } from "@/utils/md";
import { Accordion, Alert, Badge, Box, Card, Group, Text } from "@mantine/core";
import { IconAlertCircle, IconCheck, IconClock, IconSparkles } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";
import classes from "../styles/taskDetail.module.css";
import StageScriptedCard from "./StageScriptedCard";
import StageTtsCard from "./StageTtsCard";
import StylizeOrMergeCard from "./StylizeOrMergeCard";

interface TaskStageListProps {
  stages: TaskStage[];
  onRefresh: () => void;
}

export default function TaskStageList({ stages, onRefresh }: TaskStageListProps) {
  if (!stages) return <></>;

  const { t } = useTranslation();

  return (
    <Accordion variant="separated" radius="md" className={classes.accordion}>
      {stages.map((stage) => (
        <Accordion.Item key={stage.id} value={`${stage.id}`}>
          <Accordion.Control icon={<StageIcon status={stage.status} />}>
            <Group justify="space-between" wrap="nowrap" style={{ width: "100%" }}>
              <Group gap="sm">
                <IconSparkles size={18} color="#667eea" />
                <Text fw={600} size="md">
                  {t("podcast.stage." + stage.stage, { ns: "task", defaultValue: stage.stage })}
                </Text>
              </Group>
              <Badge size="sm" variant="dot" color={getStatusColor(stage.status)} style={{ marginRight: "1rem" }}>
                {t("podcast.status." + stage.status, { ns: "task", defaultValue: stage.status })}
              </Badge>
            </Group>
          </Accordion.Control>
          <Accordion.Panel>
            <Box p="md">
              <StageCard stage={stage} onRefresh={onRefresh} />
            </Box>
          </Accordion.Panel>
        </Accordion.Item>
      ))}
    </Accordion>
  );
}

// get status color
const getStatusColor = (status: TaskStageStatus) => {
  switch (status) {
    case TaskStageStatus.Completed:
      return "violet";
    case TaskStageStatus.Processing:
      return "blue";
    case TaskStageStatus.Failed:
      return "red";
    default:
      return "gray";
  }
};

// stage icon
function StageIcon({ status }: { status: string }) {
  const iconProps = { size: 20, strokeWidth: 2 };

  switch (status) {
    case TaskStageStatus.Completed:
      return <IconCheck {...iconProps} color="#667eea" />;
    case TaskStageStatus.Processing:
      return <IconClock {...iconProps} color="#3b82f6" />;
    case TaskStageStatus.Failed:
      return <IconAlertCircle {...iconProps} color="#ef4444" />;
    default:
      return <IconClock {...iconProps} color="#9ca3af" />;
  }
}

function StageCard({ stage, onRefresh }: { stage: TaskStage; onRefresh: () => void }) {
  const { t } = useTranslation("task");

  switch (stage.stage) {
    case TaskStageName.TTS:
      return <StageTtsCard stage={stage} />;
    case TaskStageName.Scripted:
      return <StageScriptedCard stage={stage} onRefresh={onRefresh} />;
    case TaskStageName.Stylize:
    case TaskStageName.Merge:
      return <StylizeOrMergeCard stage={stage} onRefresh={onRefresh} />;
  }

  return (
    <Card
      shadow="sm"
      padding="lg"
      radius="md"
      withBorder
      mb={20}
      style={{
        borderColor: stage.status === TaskStageStatus.Failed ? "#f87171" : "#e5e7eb",
        background: stage.status === TaskStageStatus.Failed ? "#fef2f2" : "white",
      }}
    >
      {stage.prompt && stage.prompt.length > 0 && (
        <Box mb="md">
          <Text size="sm" fw={500} c="dimmed" mb="xs">
            {t("podcast.prompt")}:
          </Text>
          {stage.prompt.split(/\r?\n/).map((line, index) => (
            <Text key={index} size="sm" c="gray.7">
              {line}
            </Text>
          ))}
        </Box>
      )}

      {stage.status === TaskStageStatus.Failed && (
        <Alert icon={<IconAlertCircle />} color="red" radius="md" mb="md">
          <Text size="sm" fw={500} mb={4}>
            {t("error_reason")}:
          </Text>
          <Text size="sm">{stage.reason || "Unknown Error"}</Text>
        </Alert>
      )}

      {stage.output && stage.output.length > 0 && (
        <Box className={classes.contentbox}>
          <div
            dangerouslySetInnerHTML={{
              __html: md.render(stage.output),
            }}
            style={{ fontSize: "0.95rem", lineHeight: "1.6" }}
          />
        </Box>
      )}

      {stage.taskAi && (
        <Group gap="xs" mt="md">
          <Badge variant="light" color="cyan">
            {stage.taskAi.platform}
          </Badge>
          <Badge variant="light" color="grape">
            {stage.taskAi.model}
          </Badge>
        </Group>
      )}
    </Card>
  );
}

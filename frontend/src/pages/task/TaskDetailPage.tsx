import { Loading } from "@/components";
import { getTask, PodcastTask } from "@/services";
import { ActionIcon, Divider, Group, Paper, Text } from "@mantine/core";
import { IconArrowLeft } from "@tabler/icons-react";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useParams } from "react-router";
import FloatingToolbar from "./task_detail/FloatingToolbar";
import TaskHeader from "./task_detail/TaskHeader";
import TaskStageList from "./task_detail/TaskStageList";

// task detail
export function TaskDetailPage() {
  const { t } = useTranslation();
  const { batchNo } = useParams();
  const [task, setTask] = useState<PodcastTask>();
  const [loading, setLoading] = useState<boolean>(true);
  const navigate = useNavigate();

  const fetchTaskDetail = async () => {
    if (!batchNo) return;

    const resp = await getTask(batchNo);
    if (resp) setTask(resp);
    setLoading(false);
  };

  useEffect(() => {
    fetchTaskDetail();
  }, [batchNo]);

  const handleRefresh = () => {
    setLoading(true);
    fetchTaskDetail();
  };

  if (loading || task === undefined) return <Loading />;

  return (
    <>
      <Group mb="sm" gap="xs">
        <ActionIcon variant="subtle" color="gray" size="lg" onClick={() => navigate(-1)} aria-label={t("button.back")}>
          <IconArrowLeft />
        </ActionIcon>
        <Text c="dimmed" size="sm">
          {t("button.back")}
        </Text>
      </Group>
      <Paper shadow="xl" p="xl" radius="lg">
        <TaskHeader task={task} loading={loading} onRefresh={handleRefresh} />

        <Divider mb="xl" />

        {task.stages && <TaskStageList stages={task.stages} onRefresh={handleRefresh} />}
      </Paper>
      <FloatingToolbar onRefresh={handleRefresh} />
    </>
  );
}

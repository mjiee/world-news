import { Loading } from "@/components";
import { getTask, PodcastTask } from "@/services";
import { Divider, Paper } from "@mantine/core";
import { useEffect, useState } from "react";
import { useParams } from "react-router";
import FloatingToolbar from "./task_detail/FloatingToolbar";
import TaskHeader from "./task_detail/TaskHeader";
import TaskStageList from "./task_detail/TaskStageList";

// task detail
export function TaskDetailPage() {
  const { batchNo } = useParams();
  const [task, setTask] = useState<PodcastTask>();
  const [loading, setLoading] = useState<boolean>(true);

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
      <Paper shadow="xl" p="xl" radius="lg">
        <TaskHeader task={task} loading={loading} onRefresh={handleRefresh} />

        <Divider mb="xl" />

        {task.stages && <TaskStageList stages={task.stages} onRefresh={handleRefresh} />}
      </Paper>
      <FloatingToolbar onRefresh={handleRefresh} />
    </>
  );
}

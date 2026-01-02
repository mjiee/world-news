import { Pagination } from "@/components";
import { deleteTask, PodcastTask, PodcastTaskResult, queryTask } from "@/services";
import { getPageNumber } from "@/utils/pagination";
import { ActionIcon, Badge, Box, Button, Card, Group, Space, Stack, Text, Title } from "@mantine/core";
import { IconEye, IconTrash } from "@tabler/icons-react";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router";
import { httpx } from "wailsjs/go/models";
import styles from "./styles/taskListPage.module.css";
import FloatingToolbar from "./task_detail/FloatingToolbar";

export function TaskListPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [tasks, setTasks] = useState<PodcastTask[]>([]);
  const [pagination, setPagination] = useState<httpx.Pagination>({ page: 1, limit: 25, total: 0 });
  const [loading, setLoading] = useState<boolean>(true);

  const updatePageHandler = (page: number) => {
    if (page) setPagination({ ...pagination, page: page });
    setLoading(true);
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  function getResultColor(result: PodcastTaskResult) {
    switch (result) {
      case PodcastTaskResult.Completed:
        return "green";
      case PodcastTaskResult.Failed:
        return "red";
      default:
        return "dimmed";
    }
  }

  const handleDelete = async (batchNo: string) => {
    await deleteTask(batchNo);
    setLoading(true);
  };

  const fetchTasks = async () => {
    if (!loading) return;

    const resp = await queryTask({ pagination: pagination });

    setLoading(false);
    if (!resp || !resp.data) return;
    setTasks(resp.data);
    setPagination({ ...pagination, total: resp.total });
  };

  useEffect(() => {
    fetchTasks();
  }, [loading]);

  return (
    <>
      <Box mb="xs">
        <Title order={2} c="violet" mb={4}>
          {t("podcast.page_title", { ns: "task" })}
        </Title>
        <Text c="dimmed" size="sm">
          {t("podcast.page_description", { total: pagination.total, ns: "task" })}
        </Text>
      </Box>

      <Stack gap="md">
        {tasks.map((task) => (
          <Card key={task.batchNo} padding="lg" radius="md" withBorder className={styles.card}>
            <Group justify="space-between" align="flex-start" wrap="nowrap">
              <Box style={{ flex: 1, minWidth: 0 }}>
                <Text size="lg" fw={500} mb="sm">
                  {task.title}
                </Text>

                <Group gap="xs">
                  <Badge variant="light" color="violet">
                    {task.batchNo}
                  </Badge>

                  {task.result && (
                    <Badge color={getResultColor(task.result)} variant="dot">
                      {t("podcast.result." + task.result, { ns: "task" })}
                    </Badge>
                  )}

                  <Text size="sm" c="dimmed">
                    {task.createdAt}
                  </Text>
                </Group>
              </Box>

              <Group gap="xs" style={{ flexShrink: 0 }}>
                <Button
                  variant="light"
                  color="violet"
                  leftSection={<IconEye size={16} />}
                  onClick={() => navigate("/task/" + task.batchNo)}
                  radius="md"
                >
                  {t("button.view")}
                </Button>

                <ActionIcon
                  variant="light"
                  color="red"
                  size="lg"
                  radius="md"
                  onClick={() => handleDelete(task.batchNo)}
                >
                  <IconTrash size={18} />
                </ActionIcon>
              </Group>
            </Group>
          </Card>
        ))}
      </Stack>
      <Space h="xl" />
      <Pagination page={pagination.page} total={getPageNumber(pagination)} onChange={updatePageHandler} />
      <FloatingToolbar onRefresh={() => setLoading(true)} />
    </>
  );
}

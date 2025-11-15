import { AppLayout } from "@/components";
import {
  CrawlingRecordPage,
  NewsDetailPage,
  NewsFavoritesPage,
  NewsListPage,
  SettingsPage,
  TaskDetailPage,
  TaskListPage,
} from "@/pages";
import { useRemoteServiceStore } from "@/stores";
import { Notification, Space } from "@mantine/core";
import { useEffect } from "react";
import toast, { Toaster, resolveValue } from "react-hot-toast";
import { Route, Routes } from "react-router";

function App() {
  const getService = useRemoteServiceStore((state) => state.getService);

  useEffect(() => {
    getService();
  }, []);

  return (
    <div>
      <Routes>
        <Route path="/" element={<AppLayout />}>
          <Route index element={<NewsListPage />} />
          <Route path="records" element={<CrawlingRecordPage />} />
          <Route path="settings" element={<SettingsPage />} />
          <Route path="news/list/:recordId" element={<NewsListPage />} />
          <Route path="news/detail/:newsId" element={<NewsDetailPage />} />
          <Route path="news/favorites" element={<NewsFavoritesPage />} />
          <Route path="tasks" element={<TaskListPage />} />
          <Route path="task/:batchNo" element={<TaskDetailPage />} />
        </Route>
      </Routes>
      <Toaster
        toastOptions={{
          duration: 3000,
          removeDelay: 500,
        }}
        children={(props) => (
          <Notification color={props?.type == "error" ? "red" : "blue"} onClose={() => toast.dismiss(props.id)}>
            {resolveValue(props.message, props)}
          </Notification>
        )}
      />
      <Space h="xl" />
    </div>
  );
}

export default App;

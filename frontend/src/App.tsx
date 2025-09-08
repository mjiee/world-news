import { useEffect } from "react";
import { Routes, Route } from "react-router";
import toast, { Toaster, resolveValue } from "react-hot-toast";
import { Notification, Space } from "@mantine/core";
import { HomePage } from "@/pages/home";
import { SettingsPage } from "@/pages/settings";
import { NewsDetailPage, NewsListPage, NewsFavoritesPage } from "@/pages/news";
import { LoginPage } from "@/pages/auth";
import { CrawlingRecordPage } from "./pages/record";
import { useRemoteServiceStore } from "@/stores";

function App() {
  const getService = useRemoteServiceStore((state) => state.getService);

  useEffect(() => {
    getService();
  }, []);

  return (
    <div>
      <Routes>
        <Route path="/" element={<HomePage />}>
          <Route index element={<NewsListPage />} />
          <Route path="records" element={<CrawlingRecordPage />} />
          <Route path="settings" element={<SettingsPage />} />
          <Route path="news/list/:recordId" element={<NewsListPage />} />
          <Route path="news/detail/:newsId" element={<NewsDetailPage />} />
          <Route path="news/favorites" element={<NewsFavoritesPage />} />
        </Route>

        <Route path="login" element={<LoginPage />} />
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

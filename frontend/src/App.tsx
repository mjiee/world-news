import { Routes, Route } from "react-router";
import toast, { Toaster, resolveValue } from "react-hot-toast";
import { Notification } from "@mantine/core";
import { HomePage } from "@/pages/home";
import { SettingsPage } from "@/pages/settings";
import { NewsDetailPage, NewsListPage } from "@/pages/news";
import { initRemoteService } from "@/stores";

function App() {
  initRemoteService();

  return (
    <div>
      <Routes>
        <Route index element={<HomePage />} />
        <Route path="settings" element={<SettingsPage />} />
        <Route path="news/list/:recordId" element={<NewsListPage />} />
        <Route path="news/detail/:newsId" element={<NewsDetailPage />} />
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
    </div>
  );
}

export default App;

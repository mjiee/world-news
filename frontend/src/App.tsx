import { useEffect } from "react";
import { Routes, Route } from "react-router";
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
    </div>
  );
}

export default App;

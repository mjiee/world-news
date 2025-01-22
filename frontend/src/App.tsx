import { Routes, Route } from "react-router";
import { HomePage } from "@/pages/home";
import { SettingsPage } from "@/pages/settings";
import { NewsDetailPage, NewsListPage } from "@/pages/news";

function App() {
  return (
    <div>
      <Routes>
        <Route index element={<HomePage />} />
        <Route path="settings" element={<SettingsPage />} />
        <Route path="news/list/:id" element={<NewsListPage />} />
        <Route path="news/detail/:id" element={<NewsDetailPage />} />
      </Routes>
    </div>
  );
}

export default App;

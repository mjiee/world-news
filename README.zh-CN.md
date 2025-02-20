# World-News 🌍📰

欢迎使用 **World-News** —— 你的全球新闻查看一站式桌面应用，来自多个新闻来源。无论你是想了解最新头条，探索不同地区的新闻，还是关注特定话题，World-News 都能让你轻松访问和阅读来自世界各地的最相关报道。

## 🚀 主要特性

- 实时获取来自多个全球新闻来源的最新新闻。
- 下载后支持离线阅读 — 随时随地查看新闻。
- 简单直观的用户界面，提供无缝便捷的新闻浏览体验。
- 灵活的数据源配置，允许用户自定义新闻来源。
- 支持桌面安装和网页部署，满足不同平台的需求。

## ⚙️ 技术栈

- **Wails** - 一个轻量级框架，用于使用 Go 和 Web 技术构建跨平台桌面应用。
- **React** - 用于构建用户界面的 JavaScript 库，应用于前端开发。
- **Mantine** - 现代化的 React 组件库，提供 UI 元素和 hooks。
- **Zustand** - 一个极简的 React 状态管理库，用于管理应用状态。
- **Gin** - 快速的 Go Web 框架，处理后端逻辑和 API 请求。
- **Gorm** - Go 的 ORM，用于与数据库进行交互。
- **Zap** - 一个结构化、分级的 Go 日志库，用于应用日志记录。
- **SQLite** - 轻量级、无服务器的 SQL 数据库，用于本地数据存储。

## 🛠️ 安装与设置

确保你已安装以下依赖项：

- **Go (>=1.23)**
- **Node.js (>=22)**

### 1. 克隆仓库

```bash
git clone https://github.com/mjiee/world-news.git
cd world-news
```

### 2. 构建桌面应用

按照以下步骤构建并运行桌面版本的应用：

```bash
# 构建应用
make build
```

构建完成后，你可以本地运行应用：

```bash
# 在 Linux/macOS 上运行
./bin/world-news

# 在 Windows 上运行
world-news.exe
```

### 2. 后端部署

后端使用 Gin 构建，提供获取新闻的 API。你可以通过以下方式部署后端：

#### 使用 Docker Compose 部署

运行以下命令：

```bash
# 使用 Docker Compose 启动后端
docker compose up -d
```

然后，在浏览器中访问 http://localhost:9010。

#### 本地部署

1. 确保本地有运行中的 PostgreSQL 数据库。
2. 设置数据库连接字符串：

```bash
WORLD_NEWS_DB_ADDR="host=localhost port=5432 user=world_news password=world_news dbname=world_news sslmode=disable TimeZone=Asia/Shanghai"
```

3. 构建并运行后端：

```bash
# 构建项目
make build-web

# 运行后端服务
./world-news
```

然后，在浏览器中访问 http://localhost:9010。

## ⚠️ 重要提示

- **学习项目**: World-News 是一个学习项目，旨在探索 Go、React、Wails 和其他 Web 技术。它不适用于生产或商业用途。
- **使用限制**: 请避免将本应用用于不道德活动或违反新闻提供商服务条款的数据抓取。使用本应用即表示你同意不将获取的数据用于非法或不道德的用途。

## 🙋‍♂️ 贡献

如果你希望贡献代码，欢迎 fork 该仓库并提交 pull request！这是一个开源项目，欢迎任何改进或建议。

## 📄 许可证

World-News 使用 [MIT](LICENSE) 许可证。

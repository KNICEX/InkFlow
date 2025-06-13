



# InkFlow ✒️ - 基于 Go 的内容分享社区平台

InkFlow 是一个使用 **Go + Gin + GORM** 构建的模块化 IT 内容分享社区平台，支持内容创作、AI 审核、评论互动、标签搜索、通知提醒等完整的内容分发闭环。项目采用模块化架构，具备良好的可扩展性和高并发支持能力，适合作为中型后端项目学习参考或二次开发基础。

## ✨ 项目特色

- 🧱 **模块化设计**：各业务独立封装，领域清晰，便于维护和扩展
- 🖊️ **草稿发布流程**：草稿编辑、审核、发布分离，支持 AI 内容审核
- 🔎 **全文搜索系统**：接入 Meilisearch，支持 Kafka 异步索引同步
- 📊 **交互与热度排名**：点赞、收藏、浏览行为实时记录，热度动态排行
- 🔔 **事件驱动通知**：基于事件总线的系统通知提醒机制
- 🧠 **审核接入 AI**：对用户上传内容自动进行智能审核过滤

---
### 🔗 前端仓库地址

InkFlow 项目前后端分离，前端代码请访问以下仓库：

👉 **前端仓库地址**： [https://github.com/KNICEX/ink-flow-web](https://github.com/KNICEX/ink-flow-web)

请按照前端项目中的 `README.md` 启动指引进行部署。

---

### 🔗 数据库测试数据生成仓库地址

InkFlow 项目前后端测试数据分离，测试数据生成代码请访问以下仓库：

👉 **测试数据生成仓库地址**： [https://github.com/KNICEX/InkFlow-TestData](https://github.com/KNICEX/InkFlow-TestData)

---

### 🖼️ Web 端 UI 展示

InkFlow 提供简洁直观的前端界面，支持内容创作、发布、审核、互动、搜索等完整社区功能。

以下为部分 UI 截图：

#### 🏠 首页 Feed 流

![首页 Feed 流](./docs/images/home_feed.png)

#### ✍️ 发布内容界面

![发布内容](./docs/images/post_editor.png)

#### 🔍 搜索与标签浏览

![内容关系](./docs/images/content_relation.png)

#### 📬 通知中心与评论互动

![系统通知](./docs/images/notification_system.png)

> 💡 如果你在本地运行前端项目，可通过 `http://localhost:5173` 或指定端口访问 UI 页面。

---


## 🧩 功能模块

| 模块          | 描述                                                         |
|---------------|--------------------------------------------------------------|
| `user`        | 用户注册、登录、资料维护、JWT 鉴权                          |
| `ink`         | 内容创作、草稿管理、发布审核、Feed 推送                    |
| `review`      | AI 审核模块，异步处理内容合规性审查                         |
| `comment`     | 评论与回复、分页查询、敏感词过滤                             |
| `relation`    | 用户关注/粉丝、关注事件发送                                  |
| `interactive` | 点赞、收藏、浏览行为统计                                     |
| `ranking`     | 内容热度评分与热门榜单                                       |
| `search`      | Meilisearch 搜索接入，Kafka 消息同步                         |
| `notification`| 系统通知模块，处理审核/关注/互动等事件消息                 |
| `recommend`   | 推荐模块，基于 Gorse 实现个性化内容推荐                      |
| `bff`         | Backend for Frontend，整合所有模块 API，基于 Gin 实现       |


---

## 🛠️ 技术栈

- 语言：Go 1.21+
- Web 框架：[Gin](https://gin-gonic.com/)
- ORM：GORM + PostgreSQL
- 缓存：Redis
- 消息队列：Kafka
- 搜索引擎：[Meilisearch](https://www.meilisearch.com/)
- 推荐服务：Gorse
- 审核服务：AI 审核平台（可对接 google/Gemini）
- 日志与配置：zap + viper

---

## 🧾 项目结构

```bash
inkflow/
├── config/             # 配置文件（数据库、Redis、Kafka 等）
├── pkg/                # 通用工具库（JWT、中间件、日志封装等）
├── internal/
│   ├── bff/            # API 层，Gin 路由入口，整合所有模块服务
│   ├── user/           # 用户模块（含 service/domain/events/repo）
│   ├── ink/            # 内容模块（草稿/发布/审核/Feed流）
│   ├── review/         # 审核模块
│   ├── comment/        # 评论模块
│   ├── search/         # 搜索模块（Kafka 同步 Meilisearch）
│   ├── relation/       # 关注模块
│   ├── ranking/        # 热度排名模块
│   ├── notification/   # 通知模块
│   └── interactive/    # 点赞浏览收藏模块
├── main.go             # 应用程序入口
```

## 🚀 快速启动

### ✅ 环境依赖

- Go 1.21+
- Docker（用于启动 PostgreSQL、Redis、Kafka、Meilisearch）
- Git

---

### 1. 克隆项目

```bash
git clone https://github.com/yourname/inkflow.git
cd inkflow
```
---

### 2. 启动依赖服务（推荐）

使用内置 `docker-compose.yml` 启动 PostgreSQL、Redis、Kafka、Meilisearch：

```bash
docker-compose up -d
```
---

### 3. 配置服务

复制配置模板并根据本地环境修改参数：

```bash
cp config/config.yaml.example config/config.yaml
```
#### 🛠️ 示例配置文件内容（`config.yaml`）

```yaml
postgres:
  dsn: "host=localhost user=root password=root dbname=ink_flow port=15432"

redis:
  addr: 127.0.0.1:16379
  password: 123456

# 可选：如使用 ElasticSearch（默认未启用）
# es:
#   addr: 127.0.0.1:9200
#   sniff: false

meilisearch:
  addr: http://127.0.0.1:7700
  master_key: your_master_key

email:
  smtp:
    username: example@gmail.com
    password: password
    port: 587
    host: smtp.gmail.com
    from_name: InkFlow

llm:
  gemini:
    key:
      - key

otel:
  grpc:
    endpoint: localhost:4317
    insecure: true

kafka:
  addrs:
    - localhost:9094

file:
  cloudinary:
    key: your_key
    secret: your_secret
    cloud_name: your_cloud_name

temporal:
  addr: localhost:7233
  namespace: inkflow
  domain: default

gorse:
  addr: http://localhost:8088
  api_key: your_api_key

```




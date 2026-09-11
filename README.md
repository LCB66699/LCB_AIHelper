# GopherAI

GopherAI 是一个基于 Go 和 Vue 3 的 AI 应用服务平台，提供多模型对话、流式响应、RAG 知识库问答、MCP 工具调用和可选的图像识别能力。

## 功能

- OpenAI、Ollama 和阿里云 RAG 模型接入，支持同步与 SSE 流式对话
- 基于 Redis 向量检索的 RAG 文件知识库，按用户隔离文档
- 基于 `mcp-go` 的 MCP HTTP 服务端与客户端，支持天气工具调用
- JWT 认证、用户会话和历史消息管理
- MySQL 持久化、Redis 缓存与向量索引、RabbitMQ 异步消息处理
- 基于 ONNX Runtime 和 MobileNetV2 的可选图像识别
- Docker Compose 本地部署，以及 GitHub Actions CI/CD 工作流

## 技术栈

- 后端：Go 1.24、Gin、GORM、Eino
- 前端：Vue 3、Vue Router、Element Plus、Axios
- 基础设施：MySQL 8.4、Redis Stack、RabbitMQ
- AI 与集成：OpenAI、Ollama、ONNX Runtime、MCP

## 项目结构

```text
GopherAI-v2/
├── common/          AI、RAG、MCP、消息队列和图像识别组件
├── config/          配置加载与环境变量覆盖
├── controller/      HTTP 请求处理
├── dao/             数据访问层
├── middleware/      JWT 等中间件
├── router/          API 路由
├── service/         用户、会话、文件和图像业务
├── vue-frontend/    Vue 3 前端
├── docker-compose.yml
└── scripts/         本地 CI/CD 脚本
```

## 本地运行

### Docker Compose（推荐）

要求：Docker Desktop 和 Docker Compose。

```powershell
cd GopherAI-v2
Copy-Item .env.example .env
```

编辑 `.env`，至少配置数据库、RabbitMQ、`JWT_KEY`；使用聊天或 RAG 时还需配置 AI 模型参数，例如 `OPENAI_API_KEY`、`OPENAI_MODEL_NAME` 和 `OPENAI_BASE_URL`。`.env` 只保存在本地，不应提交到 Git。

```powershell
.\scripts\deploy-local.ps1
```

启动后：

- 前端：<http://localhost:8080>
- API：<http://localhost:9090>
- 健康检查：<http://localhost:9090/healthz>

停止服务但保留数据：

```powershell
.\scripts\stop-local.ps1
```

如需启用图像识别，准备 ONNX Runtime、MobileNetV2 模型和 ImageNet 标签文件，并执行：

```powershell
.\scripts\deploy-local.ps1 -WithImageRecognition
```

### 本地开发

后端入口为 `GopherAI-v2/main.go`，前端位于 `GopherAI-v2/vue-frontend`。依赖 MySQL、Redis 和 RabbitMQ，具体配置见 `GopherAI-v2/config/config.toml` 与 `.env.example`。

## 测试与 CI/CD

本地质量检查：

```powershell
cd GopherAI-v2
.\scripts\ci-local.ps1
```

GitHub Actions 配置位于仓库根目录的 `.github/workflows/`：

- `ci.yml`：在推送和 Pull Request 中执行 Go 测试、MCP 测试、前端检查和 Docker 构建
- `cd-local.yml`：通过手动触发和 self-hosted runner 部署到本机

## 配置说明

推荐通过环境变量覆盖配置，变量前缀为 `GOPHERAI_`，例如：

```text
GOPHERAI_HTTP_PORT=9090
GOPHERAI_MYSQL_HOST=localhost
GOPHERAI_REDIS_HOST=localhost
GOPHERAI_RABBITMQ_HOST=localhost
GOPHERAI_JWT_KEY=<long-random-key>
GOPHERAI_MCP_BASE_URL=http://localhost:8081/mcp
```

完整变量和本地部署流程见 [DEPLOYMENT.md](GopherAI-v2/DEPLOYMENT.md)。

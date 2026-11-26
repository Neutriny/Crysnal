# Crysnal 开发规划

> 最后更新：2026-11-26
>
> 已确认决策：**单用户**、**Windows 优先（控制台统一 UTF-8）**、规划阶段不动手写代码。

## 一、项目现状

| 部分 | 状态 |
| --- | --- |
| 后端 | `main.go` 为模板代码，`go.mod` 无依赖，`internal/` 未创建，从零开始 |
| 前端 | 侧边栏布局已完成，路由骨架齐全：仪表盘、服务端列表/创建、隧道列表、账户设置、关于；页面为空壳 |
| 技术栈 | Go 1.27 + Vue 3 + Naive UI + Pinia + UnoCSS，vite-plus 工具链，pnpm 包管理 |

整体判断：前端壳子已就绪，主要工作在后端，然后前后端对接。

## 二、产品定位

单机/单节点 Minecraft 服务端管理面板（类似 MCSManager 单节点版、MSLX）。

核心场景：**本机开服 → 多实例管理 → Web 控制台 → 配置/文件管理 → frp 内网穿透**。

## 三、功能分期

### MVP（v0.1）—— 能开服、能管服

- 实例 CRUD：创建（核心 jar、内存、端口、JVM 参数）、启动/停止/重启/删除
- 进程管理：Go 托管 Java 子进程，stdin 命令注入，stdout/stderr 实时采集
- Web 控制台：WebSocket 实时日志 + 发送命令
- `server.properties` 可视化编辑、EULA 确认流程
- 单用户登录鉴权

### v0.2 —— 降低开服门槛

- 核心下载：Vanilla / Paper / Purpur / Fabric 等版本清单与下载（官方源 + BMCLAPI 镜像）
- Java 管理：自动检测本机 JDK、按 MC 版本推荐/下载运行时
- 文件管理器：浏览/编辑/上传/下载（限定实例目录内）
- 备份与计划任务：定时重启、定时备份、定时执行命令

### v0.3 —— 完善生态

- frp 隧道管理（内嵌 frpc 进程，对接现有「隧道列表」页面）
- 玩家管理：白名单/封禁/OP（控制台命令或 RCON）
- 历史性能监控图表

## 四、后端设计

### 技术选型

| 用途 | 选择 | 理由 |
| --- | --- | --- |
| HTTP 路由 | `chi` | 轻量、兼容标准库，适合嵌入 SPA |
| WebSocket | `coder/websocket` | 维护活跃（gorilla 已归档） |
| 存储 | SQLite（`modernc.org/sqlite`，纯 Go 免 CGO） | 实例数据量小，一步到位免迁移 |
| 日志 | `log/slog` | 标准库够用 |
| 配置 | YAML | 面板自身配置 |
| 计划任务 | `robfig/cron/v3` | v0.2 引入 |

### 模块划分

遵循 AGENTS.md：业务代码放 `internal/<domain>/`，不往 `main` 包堆逻辑。

```text
internal/
  config/    # 面板配置：端口、数据目录、密钥
  server/    # HTTP 装配、路由、中间件、SPA embed
  instance/  # 实例领域：模型、CRUD、生命周期状态机、handler
  process/   # 进程托管：启停、stdin/stdout、日志环形缓冲、事件总线
  core/      # 服务端核心：版本清单、下载器（v0.2）
  java/      # Java 检测与管理（v0.2）
  user/      # 单用户与鉴权
  storage/   # SQLite 持久化
  task/      # 计划任务（v0.2）
  frp/       # frpc 托管（v0.3）
```

### 关键设计点

1. **实例状态机**：`Stopped → Starting → Running → Stopping → Stopped`，异常退出 → `Crashed`。停止优先发 `stop` 命令优雅关服，超时（如 30s）再强杀。
2. **日志分发**：每实例「日志文件 + 内存环形缓冲（最近 N 行）+ 发布订阅」。WS 新连接先回放缓冲再实时推送，进程输出与 WS 推送解耦。
3. **跨平台进程管理（Windows 优先）**：Windows 无 SIGTERM，强杀用 `taskkill /PID <pid> /T /F` 杀进程树（Java 可能派生子进程）。
4. **控制台编码统一 UTF-8**：
   - JVM 参数：`-Dfile.encoding=UTF-8 -Dstdout.encoding=UTF-8 -Dstderr.encoding=UTF-8`（JDK 18+ 默认 UTF-8，旧版靠参数强制）。
   - Go 侧读子进程 stdout/stderr 按 UTF-8 处理，不做 GBK 转码。
   - 若调用 PowerShell（下载核心、查询进程等），脚本首行固定：

     ```powershell
     $OutputEncoding = [console]::InputEncoding = [console]::OutputEncoding = [Text.UTF8Encoding]::UTF8
     ```

5. **面板重启后的状态恢复**：MVP 简单标记为 Stopped，后续可支持按 PID 重新 attach。
6. **安全**：登录态用 session token（WS 鉴权用 query token）；文件管理严格限制在实例目录内防目录穿越；端口冲突检测。
7. **路径**：使用 `filepath`，实例目录默认放面板数据目录下（如 `data/instances/<id>/`）。

### 数据模型草案（单用户）

- `settings`：单用户密码哈希、面板偏好（键值或单行表）
- `instances(id, name, dir, core_type, core_version, jar_path, java_path, min_mem, max_mem, port, jvm_args, auto_start, created_at)`
- 实例运行时状态不持久化（内存态）
- `tasks` / `backups` / `frp_tunnels`：v0.2 / v0.3 追加

### API 草案

```text
POST   /api/auth/login
POST   /api/auth/logout
PUT    /api/auth/password            # 账户设置：修改密码

GET    /api/instances                POST /api/instances
GET    /api/instances/:id            PUT/DELETE /api/instances/:id
POST   /api/instances/:id/start | stop | restart | kill
GET    /api/instances/:id/properties PUT /api/instances/:id/properties

GET    /api/cores/versions?type=paper        # v0.2
WS     /api/instances/:id/console?token=...  # 日志回放+推送，接收命令
```

## 五、前端计划

- **API 层**：fetch 封装（统一 baseURL、错误处理、鉴权），Pinia store 管理实例列表与当前实例状态。
- **控制台组件**：引入 `xterm.js` 做终端渲染（Naive UI 无终端组件），接 WebSocket。
- 页面开发顺序：实例列表 → 创建实例（分步表单：选核心 → 基础配置 → EULA）→ 实例详情（控制台 + 属性编辑）→ 仪表盘 → 登录页。
- 路由懒加载 + 登录守卫。
- 「账户设置」= 修改密码 + 基础偏好。
- 遵循 AGENTS.md：路由组件在 `web/src/views`，大驼峰命名、以 `View` 结尾；可复用组件在 `web/src/components`。

## 六、里程碑

| 里程碑 | 内容 | 产出 |
| --- | --- | --- |
| M1 | 后端骨架：`internal/config` + `internal/server` + `internal/instance` + `internal/process` + 控制台 WS | 能用 API 创建实例并启停、看日志 |
| M2 | 前端对接：列表/创建/控制台页面 | 端到端可用的最小开服流程 |
| M3 | 登录鉴权 + 仪表盘 + properties 编辑 | MVP 完成，发布 v0.1 |
| M4 | 核心下载 + Java 管理 | v0.2 核心 |
| M5 | 文件管理 + 备份 + 计划任务 | v0.2 完成 |
| M6 | frp + 玩家管理 | v0.3 |

## 七、风险与对策

| 风险 | 对策 |
| --- | --- |
| Forge/NeoForge 安装流程复杂（需运行 installer） | MVP 只支持已有 jar 启动；核心下载先做 Vanilla/Paper 等直链核心 |
| 国内下载源可用性 | 版本清单与下载均支持镜像切换（官方源 + BMCLAPI） |
| 长时间运行的日志量与 WS 背压 | 环形缓冲 + 慢客户端断连策略 |
| Windows 编码问题 | 全链路 UTF-8（见「关键设计点」第 4 条） |

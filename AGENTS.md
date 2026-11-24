# AGENTS.md

Crysnal —— Minecraft 服务器管理面板。Go 后端 + Vue 3 SPA，单仓库。

## 目录

- `main.go`：入口。业务代码放 `internal/<domain>/`，不要往 `main` 包堆逻辑。
- `web/`：前端（Vue 3 + TS + Vite + Pinia + Vue Router），产物由 Go `embed` 提供。

## 命令

- 后端：`go build ./...`、`go test ./...`、`go vet ./...`、`gofmt -l .`
- 前端（`web/` 下）：`pnpm install`、`pnpm dev`、`pnpm build`、`pnpm format`
- 包管理器只用 pnpm。

## 通用

- 遵循 `.editorconfig`：Go 用 tab，其余用 2 空格，LF 换行。

# Crysnal

> [!WARNING]
> 项目目前处于前期开发阶段

开源 Minecraft 服务器管理面板

![Go](https://shieldcn.dev/badge/Go-v12.7.3-FFFFFF.svg?variant=branded&brand=go&size=default&mode=light)
![pnpm](https://shieldcn.dev/badge/pnpm-v12-FFFFFF.svg?variant=branded&brand=pnpm&size=default&gradient=ff6b6b%2Cfeca57&logoColor=fff&mode=light)
![Vue.js](https://shieldcn.dev/badge/Vue.js-v3-FFFFFF.svg?variant=branded&brand=vuedotjs&size=default&mode=light)
![TypeScript](https://shieldcn.dev/badge/TypeScript-v7-FFFFFF.svg?variant=branded&brand=typescript&size=default&gradient=667eea%2C764ba2&logoColor=fff&mode=light)

![Stars + License](https://shieldcn.dev/group/github/stars/Neutriny/Crysnal+github/license/Neutriny/Crysnal+github/realese/Neutriny/Crysnal+github/dt/Neutriny/Crysnal.svg&size=default&mode=light)

# 概述

**Crysnal**，发音为“ _Creesnal_ ”，由 _**Crystal**_ 和 _**Panel**_ 两个单词组成，有“水晶面板”之美好含义

这是一款开源免费的 Minecraft 服务器管理面板，使用 [**Go**](https://golang.google.cn/) 和 [**Vue 3**](https://cn.vuejs.org/)开发

网页路由和界面布局设计参考了 [**MSLX**](https://github.com/MSLTeam/MSLX)

# 字体

界面使用 **HarmonyOS Sans SC**（Copyright 2021 Huawei Device Co., Ltd.，HarmonyOS Sans Fonts License Agreement），随前端资源以 woff2 形式分发，许可全文见 [`web/public/fonts/LICENSE.txt`](./web/public/fonts/LICENSE.txt)。其他第三方组件见 [THIRD_PARTY_NOTICES.md](./THIRD_PARTY_NOTICES.md)。

# 本地开发

> [!WARNING]
> 本项目目前只能在 Windows 系统上开发，MacOS 和 Linux 请等待开发者后续支持

### 下载源码

```powershell
git clone https://github.com/Neutriny/Crysnal.git crysnal
cd crysnal
```

### 前端开发

安装 [pnpm](https://pnpm.io/zh/installation) 和 [Vite+](https://viteplus.nodejs.cn/guide/#install-vp)

> [!NOTE]
> 
> 运行前端命令前先 `cd web`

详细开发命令见 [package.json](web/package.json)

### 后端构建

安装 Go 和 [Taskfile](https://taskfile.dev/docs/installation#official-package-managers)

↪ 对于 Windows 系统，可以使用 Winget 安装：

```powershell
winget install GoLang.Go
winget install Task.Task
```


```bash
task # 默认参数dev
task dev #构建前端并嵌入二进制
task clean #清理构建产物
```
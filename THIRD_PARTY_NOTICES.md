# 第三方组件与许可

本项目源代码以 [Mozilla Public License 2.0](./LICENSE) 发布。下列第三方资源**不受** MPL-2.0 约束，各自保留原有许可证。

## 字体：HarmonyOS Sans SC

`web/public/fonts/HarmonyOS_Sans_SC_Regular.woff2`、`HarmonyOS_Sans_SC_Medium.woff2` 随前端资源一起分发（构建后位于 `dist/fonts/`，并由 Go `embed` 打包进可执行文件）。

- 版权声明：
  ```
  Copyright 2021 Huawei Device Co., Ltd.
  ```
- 许可证：**HarmonyOS Sans Fonts License Agreement**，全文见 [`web/public/fonts/LICENSE.txt`](./web/public/fonts/LICENSE.txt)（构建产物中为 `dist/fonts/LICENSE.txt`，部署后可通过 `/fonts/LICENSE.txt` 访问）。该文件必须随字体一并分发，不得删改。
- 上游来源：华为官方 HarmonyOS Sans 字体包中的 `HarmonyOS Sans/HarmonyOS_Sans_SC/`，字重 `Regular` 与 `Medium`。

### 未做修改

许可第 2 条第 2 项禁止对字体及其任何组件做修改，因此这里**只做无损容器转换**：

```bash
fonttools ttLib.woff2 compress HarmonyOS_Sans_SC_Regular.ttf -o HarmonyOS_Sans_SC_Regular.woff2
```

- 字形轮廓、hinting、`cmap`、度量、name 表均与上游 TTF 逐表一致；
- **未子集化**，完整字符集保留；
- 转换仅改变 sfnt 的封装方式（ttf → woff2），解码后与原始 TTF 等价。

### 合规对应

| 许可条款 | 要求 | 落实位置 |
|---|---|---|
| 第 2 条第 1 项 | 在软件中**显著声明**使用了 HarmonyOS Sans Fonts | [`README.md`](./README.md) 的“字体”一节、界面“关于”页（`web/src/views/AboutView.vue`） |
| 第 2 条第 2 项 | 不得修改字体或其任何组件 | 仅做 ttf → woff2 无损转换，未子集化、未改轮廓（见上） |
| 第 2 条第 3 项 | 不得单独分发或出售字体；随软件分发不受限 | 字体仅作为本面板的前端资源随包分发，不提供单独下载 |
| 第 2 条第 4 项 | 任何副本都需保留版权声明与本协议 | `web/public/fonts/LICENSE.txt` 随字体进入 `dist/fonts/` 与可执行文件 |

### 与项目许可证的关系

字体按上述许可分发，不适用本项目的 MPL-2.0；该限制不影响使用该字体呈现的任何界面、文档或截图。

## 字体：JetBrains Maple Mono

`web/public/fonts/JetBrainsMapleMono-Regular.woff2` 随前端资源一起分发（构建后位于 `dist/fonts/`，并由 Go `embed` 打包进可执行文件），用于控制台日志等等宽字体场景。

- 版权声明：

  ```text
  Copyright 2020 The JetBrains Mono Project Authors (https://github.com/JetBrains/JetBrainsMono)
  Copyright 2022 The Maple Mono Project Authors (https://github.com/subframe7536/maple-font)
  Copyright 2025 Space Time (https://github.com/SpaceTimee/Fusion-JetBrainsMapleMono)
  ```

- 许可证：**SIL Open Font License 1.1**，全文见 [`web/public/fonts/JetBrainsMapleMono-OFL.txt`](./web/public/fonts/JetBrainsMapleMono-OFL.txt)（构建产物中为 `dist/fonts/JetBrainsMapleMono-OFL.txt`，部署后可通过 `/fonts/JetBrainsMapleMono-OFL.txt` 访问）。该文件必须随字体一并分发，不得删改。
- 上游来源：[SpaceTimee/Fusion-JetBrainsMapleMono](https://github.com/SpaceTimee/Fusion-JetBrainsMapleMono)，融合 JetBrains Mono 与 Maple Mono；使用 NL（No Ligatures）变体，字重 `Regular`。

### 未做修改

与 HarmonyOS Sans SC 相同，只做无损容器转换（ttf → woff2），未子集化、未改轮廓，字体名称表未改动：

```bash
fonttools ttLib.woff2 compress JetBrainsMapleMono-Regular.ttf -o JetBrainsMapleMono-Regular.woff2
```

### 合规对应

| OFL 条款 | 要求 | 落实位置 |
| --- | --- | --- |
| 条件 2 | 任何副本都需包含版权声明与许可证 | `web/public/fonts/JetBrainsMapleMono-OFL.txt` 随字体进入 `dist/fonts/` 与可执行文件 |
| 条件 3 | 修改版不得使用保留字体名 | 仅无损容器转换，字体名称表未改动 |
| 条件 5 | 字体须整体按 OFL 分发 | 字体按 OFL-1.1 分发，不适用本项目的 MPL-2.0 |

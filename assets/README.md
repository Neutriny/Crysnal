# assets

Windows 可执行文件的资源素材，由链接器合并进 `crysnal.exe`（本项目没有系统托盘，因此不含托盘图标）。

| 文件 | 说明 |
| --- | --- |
| `app.ico` | 16 / 24 / 32 / 48 / 64 / 128（32bpp DIB）+ 256（PNG 帧）：exe 图标、任务栏、开始菜单、Alt+Tab |
| `versioninfo.json` | 资源版本信息（文件说明、产品名、内部名、原始文件名） |

构建时 `task gen-res` 会调用 `goversioninfo`（版本固定在 `Taskfile.yml`）把这两者编译成
`rsrc_windows_amd64.syso`，放在仓库根目录，由 `go build` 自动链接；该文件是构建产物，已加入 `.gitignore`。

注意：`.syso` 文件名带 `_windows_amd64` 后缀，因此只在 windows/amd64 构建时生效，其他平台构建会自动忽略它
（要支持 windows/arm64 需另生成一份，或在 `Taskfile.yml` 中加 `-platform-specific`）。

## 版本信息为何要写版本字符串

`versioninfo.json` 的 `StringFileInfo` 里 `FileVersion` / `ProductVersion` 不能省略：goversioninfo
缺少这两个字段时写出的版本资源无法被 `VerQueryValue` 读取，资源管理器和 `FileVersionInfo` 会显示空白。

## 重新生成各档 PNG

`16` / `32` 直接在设计稿里按该尺寸单独导出（小尺寸需要手动微调，下采样会糊）；其余从 512 母版取整数倍：

```powershell
magick logo-512.png -filter Lanczos -resize 48x48   logo-48.png
magick logo-512.png -filter Lanczos -resize 64x64   logo-64.png
magick logo-512.png -filter Lanczos -resize 128x128 logo-128.png
```

`app.ico` 用 DIB 帧（≤128）+ PNG 帧（256）混合打包：`magick -define icon:auto-resize` 会把 256 存成
未压缩 BMP（约 372KB），而纯 PNG 帧则可能让走 GDI `LoadImage` 的代码读不到图标。

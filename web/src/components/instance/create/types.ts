// 创建实例向导的共享表单数据，各步骤面板通过 defineModel 读写
export interface CreateInstanceForm {
  name: string
  path: string
  core: CoreFile
}

// 核心文件：在线下载或本地已有
export interface CoreFile {
  /** online = 在线下载，local = 浏览本地文件 */
  source: 'online' | 'local'
  /** 在线下载：核心类型，如 vanilla / paper / fabric，v0.2 接下载 API */
  coreType: string
  /** 在线下载：MC 版本，如 1.21.1 */
  mcVersion: string
  /** 在线下载：下载镜像，v0.2 支持镜像切换 */
  mirror: string
  /** 本地文件：服务端 jar 的完整路径 */
  localPath: string
}

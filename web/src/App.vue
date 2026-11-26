<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { NConfigProvider, NCollapse, NCollapseItem } from 'naive-ui'
import type { GlobalThemeOverrides } from 'naive-ui'

// 折叠面板标题默认左右内边距为 0，这里改成和菜单链接一致（12px），使标题与「仪表盘」左对齐
const themeOverrides: GlobalThemeOverrides = {
  common: {
    fontFamily:
      "'HarmonyOS Sans SC', -apple-system, 'PingFang SC', 'Hiragino Sans GB', 'Segoe UI', 'Microsoft YaHei', 'Noto Sans CJK SC', 'Source Han Sans SC', 'WenQuanYi Micro Hei', system-ui, sans-serif",
  },
  Collapse: {
    titlePadding: '8px 0 8px 12px',
    titleFontSize: '15px',
  },
}
</script>

<template>
  <NConfigProvider :theme-overrides="themeOverrides">
    <div class="layout">
      <aside class="sidebar">
        <div class="brand">Crysnal 管理中心</div>

        <nav class="menu">
          <RouterLink class="menu-link" to="/dashboard/base">仪表盘</RouterLink>

          <NCollapse :default-expanded-names="['instance', 'frp', 'setting']" arrow-placement="right">
            <NCollapseItem title="服务端管理" name="instance">
              <div class="submenu">
                <RouterLink class="menu-link" to="/instance/list">服务端列表</RouterLink>
                <RouterLink class="menu-link" to="/instance/create">创建服务端</RouterLink>
              </div>
            </NCollapseItem>

            <NCollapseItem title="隧道管理" name="frp">
              <div class="submenu">
                <RouterLink class="menu-link" to="/frp/list">隧道列表</RouterLink>
              </div>
            </NCollapseItem>

            <NCollapseItem title="设置" name="setting">
              <div class="submenu">
                <RouterLink class="menu-link" to="/setting/profile">账户设置</RouterLink>
              </div>
            </NCollapseItem>
          </NCollapse>

          <RouterLink class="menu-link" to="/about">关于</RouterLink>
        </nav>
      </aside>

      <main class="content">
        <RouterView />
      </main>
    </div>
  </NConfigProvider>
</template>

<style scoped>
.layout {
  display: flex;
  height: 100vh;
}

.sidebar {
  width: 200px;
  flex: 0 0 200px;
  overflow-y: auto;
  padding: 16px 12px;
  border-right: 1px solid #efeff5;
  background: #fff;
}

.brand {
  padding: 4px 12px 16px;
  font-size: 20px;
  font-weight: 600;
}

.menu {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.menu-link {
  display: block;
  padding: 8px 12px;
  font-size: 15px;
  line-height: 22px;
  border-radius: 4px;
  color: #333;
  text-decoration: none;
  transition: background-color 0.2s ease;
}

.menu-link:hover {
  background: #f2f2f5;
}

.menu-link.router-link-active {
  color: #18a058;
  background: #f0faf4;
}

/* 分组间距与 .menu 的 gap 保持一致，并去掉 item 之间的分割线 */
.sidebar .menu :deep(.n-collapse .n-collapse-item) {
  margin: 4px 0 0 0;
  border-top: none;
}

.sidebar .menu :deep(.n-collapse .n-collapse-item:first-child) {
  margin-top: 0;
}

.sidebar .menu :deep(.n-collapse .n-collapse-item__header) {
  line-height: 22px;
}

/* naive-ui 会强制首个 item 的 header padding-top: 0，这里补回来才能等高 */
.sidebar .menu :deep(.n-collapse .n-collapse-item:first-child > .n-collapse-item__header) {
  padding-top: 8px;
}

/* 子菜单不要离分组标题太远 */
.sidebar .menu :deep(.n-collapse .n-collapse-item__content-inner) {
  padding-top: 4px;
}

.submenu {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.submenu .menu-link {
  padding-left: 24px;
}

.content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}
</style>

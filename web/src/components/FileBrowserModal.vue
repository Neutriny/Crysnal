<script setup lang="ts">
import { MdArrowBack, MdDocument, MdFolder } from '@vicons/ionicons4'
import {
  NAlert,
  NButton,
  NCheckbox,
  NEmpty,
  NIcon,
  NInput,
  NModal,
  NScrollbar,
  NSpin,
} from 'naive-ui'
import { computed, ref, watch } from 'vue'

defineOptions({ name: 'FileBrowserModal' })

// 与后端 GET /api/system/fs/list 的 dirEntry 对应
interface DirEntry {
  name: string
  path: string
  isDir: boolean
  size: number
  modTime: string
}

// 与后端 GET /api/system/fs/list 的响应对应
interface DirListResponse {
  path: string
  parent: string
  entries: DirEntry[]
  truncated: boolean
  error?: string
}

const show = defineModel<boolean>('show', { required: true })

const props = withDefaults(
  defineProps<{
    title?: string
    /** 只显示这些扩展名的文件（不含点），如 ['jar']；为空表示不过滤 */
    extensions?: string[]
    /** 打开时定位到的初始目录，为空则从盘符列表（此电脑）开始 */
    initialPath?: string
  }>(),
  { title: '选择文件', extensions: () => [], initialPath: '' },
)

const emit = defineEmits<{ confirm: [path: string] }>()

// currentPath 为空表示正在展示「此电脑」的盘符列表
const currentPath = ref('')
const parentPath = ref('')
const entries = ref<DirEntry[]>([])
const truncated = ref(false)
const keyword = ref('')
const onlyExtensions = ref(true)
const loading = ref(false)
const error = ref('')
const selected = ref<DirEntry>()
const pathInput = ref('')

const extFilter = computed(() => props.extensions.join(','))
const filterActive = computed(() => onlyExtensions.value && extFilter.value !== '')

// 关键字过滤只在前端做，切换目录时清空
const visibleEntries = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return entries.value
  return entries.value.filter((entry) => entry.name.toLowerCase().includes(kw))
})

async function load(path: string) {
  loading.value = true
  error.value = ''
  try {
    const query = new URLSearchParams()
    if (path) query.set('path', path)
    if (filterActive.value) query.set('ext', extFilter.value)

    const res = await fetch(`/api/system/fs/list?${query}`)
    const data = (await res.json().catch(() => null)) as DirListResponse | null
    if (!res.ok || !data) {
      error.value = data?.error ?? `读取目录失败（HTTP ${res.status}）`
      return
    }
    currentPath.value = data.path ?? ''
    parentPath.value = data.parent ?? ''
    entries.value = data.entries ?? []
    truncated.value = data.truncated ?? false
    selected.value = undefined
  } catch {
    error.value = '无法连接面板后端'
  } finally {
    loading.value = false
  }
}

function submitPath() {
  const target = pathInput.value.trim()
  if (!target) return
  keyword.value = ''
  void load(target)
}

function goUp() {
  keyword.value = ''
  void load(parentPath.value)
}

function pick(entry: DirEntry) {
  selected.value = entry
}

function open(entry: DirEntry) {
  if (!entry.isDir) return
  keyword.value = ''
  void load(entry.path)
}

function openOrConfirm(entry: DirEntry) {
  if (entry.isDir) {
    open(entry)
    return
  }
  selected.value = entry
  confirm()
}

function confirm() {
  if (!selected.value) return
  emit('confirm', selected.value.path)
  show.value = false
}

watch(currentPath, (path) => {
  pathInput.value = path
})

watch(filterActive, () => {
  if (show.value) void load(currentPath.value)
})

watch(show, (visible) => {
  if (!visible) return
  keyword.value = ''
  selected.value = undefined
  error.value = ''
  void load(props.initialPath || currentPath.value)
})

function formatSize(size: number) {
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let value = size
  let unit = 0
  while (value >= 1024 && unit < units.length - 1) {
    value /= 1024
    unit++
  }
  return unit === 0 ? `${value} B` : `${value.toFixed(value >= 100 ? 0 : 1)} ${units[unit]}`
}

function formatTime(raw: string) {
  if (!raw) return ''
  const date = new Date(raw)
  if (Number.isNaN(date.getTime())) return ''
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`
}
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    :title="title"
    class="file-browser"
    :bordered="false"
    :mask-closable="false"
    :style="{ width: '600px' }"
  >
    <div class="toolbar">
      <n-button :disabled="!parentPath" @click="goUp">
        <template #icon>
          <n-icon>
            <MdArrowBack />
          </n-icon>
        </template>
        上一级
      </n-button>
      <n-input
        v-model:value="pathInput"
        class="path-input"
        :placeholder="currentPath ? '目录路径' : '此电脑'"
        @keyup.enter="submitPath"
        @blur="submitPath"
      />
    </div>

    <div class="filters">
      <n-input v-model:value="keyword" class="keyword-input" clearable placeholder="在当前目录中筛选" />
      <n-checkbox v-if="extensions.length" v-model:checked="onlyExtensions">
        仅显示 {{ extensions.map((ext) => `*.${ext}`).join(' / ') }}
      </n-checkbox>
    </div>

    <n-spin :show="loading">
      <div class="list">
        <n-scrollbar style="height: 100%">
          <n-alert v-if="error" type="error" :bordered="false" class="alert">{{ error }}</n-alert>

        <div
          v-for="entry in visibleEntries"
          :key="entry.path"
          class="row"
          :class="{ selected: selected?.path === entry.path }"
          :title="entry.path"
          @click="pick(entry)"
          @dblclick="openOrConfirm(entry)"
        >
          <n-icon class="row-icon" :class="{ folder: entry.isDir }">
            <MdFolder v-if="entry.isDir" />
            <MdDocument v-else />
          </n-icon>
          <span class="row-name" :title="entry.name">{{ entry.name }}</span>
          <span class="row-meta">
            <span v-if="!entry.isDir" class="row-size">{{ formatSize(entry.size) }}</span>
            <span class="row-time">{{ formatTime(entry.modTime) }}</span>
          </span>
        </div>

        <n-empty
          v-if="!loading && !visibleEntries.length"
          :description="error ? '无法显示该目录' : '这里没有可选择的项目'"
          class="list-empty"
        />

          <div v-if="truncated" class="hint">
            项目过多，仅显示前 {{ visibleEntries.length }} 项，可用筛选缩小范围
          </div>
        </n-scrollbar>
      </div>
    </n-spin>

    <template #footer>
      <div class="footer">
        <span class="selection" :class="{ placeholder: !selected }" :title="selected?.path ?? ''">
          {{ selected ? selected.path : '尚未选择文件（双击文件可直接确定）' }}
        </span>
        <div class="actions">
          <n-button @click="show = false">取消</n-button>
          <n-button type="primary" :disabled="!selected" @click="confirm">确定</n-button>
        </div>
      </div>
    </template>
  </n-modal>
</template>

<style scoped>
.file-browser.n-card {
  width: 600px;
  max-width: calc(100vw - 24px);
}

.toolbar {
  display: flex;
  gap: 8px;
}

.path-input {
  flex: 1;
  min-width: 0;
}

.filters {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: 10px;
}

.keyword-input {
  max-width: 240px;
}

/* 固定高度，超出部分在列表内部滚动 */
.list {
  height: 270px;
  margin-top: 10px;
  border: 1px solid #efeff5;
  border-radius: 4px;
  overflow: hidden;
}

.alert {
  margin: 8px;
}

.row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 5px 10px;
  font-size: 13px;
  cursor: default;
  user-select: none;
}

.row:hover {
  background: #f7f7fa;
}

.row.selected {
  background: #f0faf4;
}

.row-icon {
  flex: 0 0 auto;
  font-size: 18px;
  color: #909399;
}

.row-icon.folder {
  color: #e6a23c;
}

.row-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.row-meta {
  flex: 0 0 auto;
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #909399;
}

.row-size {
  width: 64px;
  text-align: right;
}

.row-time {
  width: 108px;
  text-align: right;
}

.list-empty {
  padding: 32px 0;
}

.hint {
  padding: 0 10px 8px;
  font-size: 12px;
  color: #909399;
}

.footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  /* 固定高度：选中项切换时不让弹窗尺寸跳动 */
  height: 34px;
}

.selection {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 20px;
  font-size: 13px;
  color: #333;
}

.selection.placeholder {
  color: #909399;
}

.actions {
  flex: 0 0 auto;
  display: flex;
  gap: 8px;
}
</style>

<style>
/* NScrollbar 的滚动容器由 naive-ui 内部渲染，拿不到 scoped 属性，这里用组件独有类名隔离 */
.file-browser .n-scrollbar-container {
  overscroll-behavior: contain;
}
</style>

<script setup lang="ts">
import { MdFolderOpen } from '@vicons/ionicons4'
import {
  NButton,
  NForm,
  NFormItem,
  NIcon,
  NInput,
  NInputGroup,
  NSelect,
  NTabs,
  NTabPane,
} from 'naive-ui'
import type { FormRules } from 'naive-ui'
import { ref } from 'vue'
import FileBrowserModal from '@/components/FileBrowserModal.vue'
import type { CreateInstanceForm } from './types'

defineOptions({ name: 'CoreFilePanel' })
const formValue = defineModel<CreateInstanceForm>({ required: true })

// v0.2 接后端版本清单 API 后改为动态加载
const coreTypeOptions = [
  { label: 'Vanilla（官方原版）', value: 'vanilla' },
  { label: 'Paper', value: 'paper' },
  { label: 'Purpur', value: 'purpur' },
  { label: 'Fabric', value: 'fabric' },
]

const mirrorOptions = [
  { label: '官方源', value: 'official' },
  { label: 'BMCLAPI（国内镜像）', value: 'bmclapi' },
]

const rules: FormRules = {
  'core.coreType': [{ required: true, message: '请选择核心类型', trigger: 'change' }],
  'core.mcVersion': [{ required: true, message: '请选择 MC 版本', trigger: 'change' }],
  'core.localPath': [{ required: true, message: '请指定服务端核心文件路径', trigger: 'blur' }],
}

// 本地文件由面板内置的文件浏览器选择（浏览器拿不到真实路径，也不弹原生对话框）
const browserVisible = ref(false)

function browseLocalFile() {
  browserVisible.value = true
}

function onFilePicked(path: string) {
  formValue.value.core.localPath = path
}
</script>

<template>
  <n-form :model="formValue" :rules="rules" label-placement="top">
    <n-tabs v-model:value="formValue.core.source" type="segment" animated class="core-tabs">
      <n-tab-pane name="online" tab="在线下载">
        <n-form-item label="核心类型" path="core.coreType">
          <n-select
            v-model:value="formValue.core.coreType"
            :options="coreTypeOptions"
            placeholder="选择服务端核心"
          />
        </n-form-item>
        <n-form-item label="MC 版本" path="core.mcVersion">
          <n-select
            v-model:value="formValue.core.mcVersion"
            :options="[]"
            placeholder="选择版本，接版本清单后可选"
          />
        </n-form-item>
        <n-form-item label="下载源" path="core.mirror">
          <n-select v-model:value="formValue.core.mirror" :options="mirrorOptions" />
        </n-form-item>
      </n-tab-pane>
      <n-tab-pane name="local" tab="选择本地文件">
        <n-form-item label="核心文件路径" path="core.localPath">
          <n-input-group>
            <n-input
              v-model:value="formValue.core.localPath"
              placeholder="服务端 jar 完整路径，如 D:\servers\paper-1.21.1.jar"
            />
            <n-button @click="browseLocalFile">
              <template #icon>
                <n-icon>
                  <MdFolderOpen />
                </n-icon>
              </template>
            </n-button>
          </n-input-group>
        </n-form-item>
      </n-tab-pane>
    </n-tabs>
  </n-form>

  <FileBrowserModal
    v-model:show="browserVisible"
    title="选择服务端核心文件"
    :extensions="['jar']"
    @confirm="onFilePicked"
  />
</template>

<style scoped>
/* segment 模式默认撑满整行，这里限制轨道宽度让标签紧凑；pane 不受影响 */
.core-tabs :deep(.n-tabs-rail) {
  max-width: 320px;
}

/* 标签过长时省略，避免撑破轨道 */
.core-tabs :deep(.n-tabs-tab__label) {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>

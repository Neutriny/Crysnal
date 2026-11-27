<script setup lang="ts">
defineOptions({ name: 'InstanceCreateView' })
import { MdArrowRoundBack, MdArrowRoundForward } from '@vicons/ionicons4'
import {
  NButton,
  NButtonGroup,
  NIcon,
  NLayout,
  NLayoutContent,
  NLayoutSider,
  NStep,
  NSteps,
} from 'naive-ui'
import { ref } from 'vue'
import BasicInfoPanel from '@/components/instance/create/BasicInfoPanel.vue'
import CoreFilePanel from '@/components/instance/create/CoreFilePanel.vue'
import JavaEnvPanel from '@/components/instance/create/JavaEnvPanel.vue'
import ResourcePanel from '@/components/instance/create/ResourcePanel.vue'
import ConfirmPanel from '@/components/instance/create/ConfirmPanel.vue'
import CreatingPanel from '@/components/instance/create/CreatingPanel.vue'
import DonePanel from '@/components/instance/create/DonePanel.vue'
import type { CreateInstanceForm } from '@/components/instance/create/types'

// 各步骤面板，与下方 n-step 一一对应
const panels = [
  BasicInfoPanel,
  CoreFilePanel,
  JavaEnvPanel,
  ResourcePanel,
  ConfirmPanel,
  CreatingPanel,
  DonePanel,
]
const stepCount = panels.length

const currentRef = ref<number | null>(1)
const current = currentRef

const formValue = ref<CreateInstanceForm>({
  name: '',
  path: '',
  core: {
    source: 'online',
    coreType: '',
    mcVersion: '',
    mirror: 'official',
    localPath: '',
  },
})

function next() {
  if (currentRef.value === null) currentRef.value = 1
  else if (currentRef.value >= stepCount) currentRef.value = null
  else currentRef.value++
}

function prev() {
  if (currentRef.value === 0) currentRef.value = null
  else if (currentRef.value === null) currentRef.value = stepCount
  else currentRef.value--
}
</script>

<template>
  <n-layout has-sider style="height: 100%">
    <n-layout-sider width="270" content-style="padding: 16px;">
      <n-steps vertical :current="current as number">
        <n-step title="基本信息" description="填写实例名称和路径" />
        <n-step title="核心文件" description="指定核心文件及下载" />
        <n-step title="Java 环境" description="配置 Java 运行时" />
        <n-step title="资源配置" description="设置内存与 JVM 参数" />
        <n-step title="确认信息" description="核对并提交" />
        <n-step title="创建实例" description="提交并等待创建" />
        <n-step title="完成" description="查看创建结果" />
      </n-steps>
    </n-layout-sider>
    <n-layout-content content-style="padding: 24px; display: flex; flex-direction: column;">
      <div class="step-body">
        <!-- 每一步对应一个独立组件，current 为 null（未开始）时显示提示 -->
        <component :is="panels[current - 1]" v-if="current" v-model="formValue" />
        <section v-else>点击右箭头开始创建</section>
      </div>
      <div class="step-actions">
        <n-button-group>
          <n-button @click="prev">
            <template #icon>
              <n-icon>
                <MdArrowRoundBack />
              </n-icon>
            </template>
          </n-button>
          <n-button @click="next">
            <template #icon>
              <n-icon>
                <MdArrowRoundForward />
              </n-icon>
            </template>
          </n-button>
        </n-button-group>
      </div>
    </n-layout-content>
  </n-layout>
</template>

<style scoped>
.step-body {
  flex: 1;
  min-height: 0;
}

.step-actions {
  display: flex;
  justify-content: flex-end;
  padding-top: 16px;
}
</style>

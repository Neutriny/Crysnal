<script setup lang="ts">
import type { CreateInstanceForm } from './types'

defineOptions({ name: 'CreatingPanel' })
defineModel<CreateInstanceForm>({ required: true })
</script>

<template>
  <section>创建实例：提交并等待创建（待实现）</section>
</template>

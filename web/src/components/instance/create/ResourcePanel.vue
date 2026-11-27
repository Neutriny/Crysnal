<script setup lang="ts">
import type { CreateInstanceForm } from './types'

defineOptions({ name: 'ResourcePanel' })
defineModel<CreateInstanceForm>({ required: true })
</script>

<template>
  <section>资源配置：设置内存与 JVM 参数（待实现）</section>
</template>

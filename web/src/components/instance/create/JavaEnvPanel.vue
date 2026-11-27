<script setup lang="ts">
import type { CreateInstanceForm } from './types'

defineOptions({ name: 'JavaEnvPanel' })
defineModel<CreateInstanceForm>({ required: true })
</script>

<template>
  <section>Java 环境：配置 Java 运行时（待实现）</section>
</template>

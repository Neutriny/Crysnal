<script setup lang="ts">
import type { CreateInstanceForm } from './types'

defineOptions({ name: 'DonePanel' })
defineModel<CreateInstanceForm>({ required: true })
</script>

<template>
  <section>完成：查看创建结果（待实现）</section>
</template>

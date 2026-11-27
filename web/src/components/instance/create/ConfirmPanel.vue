<script setup lang="ts">
import type { CreateInstanceForm } from './types'

defineOptions({ name: 'ConfirmPanel' })
defineModel<CreateInstanceForm>({ required: true })
</script>

<template>
  <section>确认信息：核对并提交（待实现）</section>
</template>

<script setup lang="ts">
import { useCounterStore } from '@/stores/counter'

const counter = useCounterStore()
</script>

<template>
  <h1>仪表盘</h1>

  <p>下面这组数字来自 Pinia store，多个页面共享同一份状态。</p>
  <p>count: {{ counter.count }}</p>
  <p>doubleCount: {{ counter.doubleCount }}</p>
  <button @click="counter.increment()">+1</button>
</template>

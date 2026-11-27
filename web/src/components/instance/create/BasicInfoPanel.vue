<script setup lang="ts">
import { NForm, NFormItem, NInput } from 'naive-ui'
import type { FormRules } from 'naive-ui'
import type { CreateInstanceForm } from './types'

defineOptions({ name: 'BasicInfoPanel' })

const formValue = defineModel<CreateInstanceForm>({ required: true })

const rules: FormRules = {
  name: [{ required: true, message: '请输入实例名称', trigger: 'blur' }],
}
</script>

<template>
  <n-form :model="formValue" :rules="rules" label-placement="top">
    <n-form-item label="实例名称" path="name">
      <n-input
        v-model:value="formValue.name"
        placeholder="新建服务器（可以填中文）"
      />
    </n-form-item>
    <n-form-item label="实例路径" path="path">
      <n-input v-model:value="formValue.path" placeholder="留空使用默认路径" />
      <template #feedback>默认：%USERPROFILE%\Crysnal\instances\实例名，可改为其他目录</template>
    </n-form-item>
  </n-form>
</template>

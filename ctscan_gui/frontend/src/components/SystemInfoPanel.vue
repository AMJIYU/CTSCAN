<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetSystemInfo } from '../../wailsjs/go/pkg/App'

const sysInfo = ref({
  hostname: '',
  os: '',
  arch: '',
  cpu_cores: 0,
  go_version: ''
})

const refresh = () => {
  GetSystemInfo().then(info => {
    sysInfo.value = info
  })
}

onMounted(() => {
  refresh()
})

defineExpose({ refresh })
</script>
<template>
  <el-descriptions title="系统基本信息" :column="2" border>
    <el-descriptions-item label="主机名">{{ sysInfo.hostname }}</el-descriptions-item>
    <el-descriptions-item label="操作系统">{{ sysInfo.os }}</el-descriptions-item>
    <el-descriptions-item label="架构">{{ sysInfo.arch }}</el-descriptions-item>
    <el-descriptions-item label="CPU核心数">{{ sysInfo.cpu_cores }}</el-descriptions-item>
    <el-descriptions-item label="Go版本">{{ sysInfo.go_version }}</el-descriptions-item>
  </el-descriptions>
</template> 
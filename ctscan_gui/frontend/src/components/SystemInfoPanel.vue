<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetSystemInfo } from '../../wailsjs/go/pkg/App'

interface DiskInfo {
  mount_point: string
  total_size: number
  used_size: number
  free_size: number
  usage: number
}

interface SystemInfo {
  hostname: string
  os: string
  arch: string
  cpu_cores: number
  kernel_version: string
  cpu_usage: number
  total_memory: number
  memory_usage: number
  disks: DiskInfo[]
}

const sysInfo = ref<SystemInfo>({
  hostname: '',
  os: '',
  arch: '',
  cpu_cores: 0,
  kernel_version: '',
  cpu_usage: 0,
  total_memory: 0,
  memory_usage: 0,
  disks: []
})

const formatMemory = (bytes: number) => {
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let size = bytes
  let unitIndex = 0
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  return `${size.toFixed(2)} ${units[unitIndex]}`
}

const refresh = () => {
  GetSystemInfo().then(info => {
    sysInfo.value = info
  })
}

onMounted(() => {
  refresh()
  // 每5秒刷新一次系统信息
  setInterval(refresh, 5000)
})

defineExpose({ refresh })
</script>
<template>
  <el-descriptions title="系统基本信息" :column="2" border>
    <el-descriptions-item label="主机名">{{ sysInfo.hostname }}</el-descriptions-item>
    <el-descriptions-item label="操作系统">{{ sysInfo.os }}</el-descriptions-item>
    <el-descriptions-item label="架构">{{ sysInfo.arch }}</el-descriptions-item>
    <el-descriptions-item label="CPU核心数">{{ sysInfo.cpu_cores }}</el-descriptions-item>
    <el-descriptions-item label="内核版本">{{ sysInfo.kernel_version }}</el-descriptions-item>
    <el-descriptions-item label="CPU使用率">
      <el-progress :percentage="sysInfo.cpu_usage" :format="(val) => val.toFixed(1) + '%'" />
    </el-descriptions-item>
    <el-descriptions-item label="内存总量">{{ formatMemory(sysInfo.total_memory) }}</el-descriptions-item>
    <el-descriptions-item label="内存使用率">
      <el-progress :percentage="sysInfo.memory_usage" :format="(val) => val.toFixed(1) + '%'" />
    </el-descriptions-item>
  </el-descriptions>

  <el-divider />

  <el-descriptions title="磁盘信息" :column="1" border>
    <el-descriptions-item v-for="disk in sysInfo.disks" :key="disk.mount_point">
      <div class="disk-info">
        <h4>{{ disk.mount_point }}</h4>
        <div class="disk-details">
          <div>总容量: {{ formatMemory(disk.total_size) }}</div>
          <div>已用空间: {{ formatMemory(disk.used_size) }}</div>
          <div>可用空间: {{ formatMemory(disk.free_size) }}</div>
          <div>使用率:
            <el-progress :percentage="disk.usage" :format="(val) => val.toFixed(1) + '%'" />
          </div>
        </div>
      </div>
    </el-descriptions-item>
  </el-descriptions>
</template>

<style scoped>
.disk-info {
  padding: 10px;
}

.disk-info h4 {
  margin: 0 0 10px 0;
  color: #409EFF;
}

.disk-details {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.disk-details > div {
  display: flex;
  align-items: center;
  gap: 10px;
}

.disk-details > div:last-child {
  grid-column: 1 / -1;
}
</style> 
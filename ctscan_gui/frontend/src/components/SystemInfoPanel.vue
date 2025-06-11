<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetSystemInfo } from '../../wailsjs/go/pkg/App'
import { Folder } from '@element-plus/icons-vue'

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

  <div class="info-card">
    <div class="card-header">
      <el-icon :size="18" color="#409EFF"><Folder /></el-icon>
      <h3>磁盘信息</h3>
    </div>
    <div class="disk-list">
      <div v-for="disk in sysInfo.disks" :key="disk.mount_point" class="disk-item">
        <div class="disk-name">
          <el-icon><Folder /></el-icon>
          <span>{{ disk.mount_point }}</span>
        </div>
        <div class="disk-stats">
          <span class="stat-item">
            <span class="label">总容量</span>
            <span class="value">{{ formatMemory(disk.total_size) }}</span>
          </span>
          <span class="stat-item">
            <span class="label">已用</span>
            <span class="value">{{ formatMemory(disk.used_size) }}</span>
          </span>
          <span class="stat-item">
            <span class="label">可用</span>
            <span class="value">{{ formatMemory(disk.free_size) }}</span>
          </span>
          <span class="stat-item">
            <span class="label">使用率</span>
            <span class="value">{{ (disk.used_size / disk.total_size * 100).toFixed(1) }}%</span>
          </span>
        </div>
        <el-progress 
          :percentage="(disk.used_size / disk.total_size * 100)" 
          :status="(disk.used_size / disk.total_size * 100) > 90 ? 'exception' : ''"
          :stroke-width="8"
          :show-text="false"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.info-card {
  margin-top: 24px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  color: #1a202c;
  font-weight: 600;
  font-size: 16px;
}

.disk-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.disk-item {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.disk-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.08);
}

.disk-name {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  color: #1a202c;
  font-weight: 600;
  font-size: 15px;
}

.disk-name .el-icon {
  color: #409EFF;
}

.disk-stats {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.stat-item .label {
  font-size: 12px;
  color: #718096;
}

.stat-item .value {
  font-family: monospace;
  font-size: 13px;
  color: #2d3748;
  font-weight: 500;
}

:deep(.el-descriptions) {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
}

:deep(.el-descriptions__title) {
  color: #1a202c;
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 24px;
}

:deep(.el-descriptions__label) {
  color: #4a5568;
  font-weight: 500;
}

:deep(.el-descriptions__content) {
  color: #1a202c;
}

:deep(.el-progress-bar__outer) {
  background-color: rgba(0, 0, 0, 0.05);
  border-radius: 4px;
}

:deep(.el-progress-bar__inner) {
  border-radius: 4px;
}

:deep(.el-divider) {
  margin: 32px 0;
  border-color: rgba(0, 0, 0, 0.05);
}
</style> 
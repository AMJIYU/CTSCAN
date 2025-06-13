<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetSystemInfo, SaveSystemInfo } from '../../wailsjs/go/pkg/App'
import { Monitor, Cpu, Memo, Folder } from '@element-plus/icons-vue'

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

const systemInfo = ref<SystemInfo>({
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

const loading = ref(false)

// 添加 refresh 方法，用于重新获取系统信息
const refresh = async () => {
  loading.value = true
  try {
    const info = await GetSystemInfo()
    systemInfo.value = info
    // 保存到数据库
    await SaveSystemInfo(info)
  } catch (error) {
    console.error('获取系统信息失败:', error)
  } finally {
    
    loading.value = false
  }
}

onMounted(() => {
  refresh()
})

// 暴露 refresh 方法，供父组件调用
defineExpose({ refresh })
</script>

<template>
  <div class="system-info-panel">
    <div class="info-section">
      <div class="section-header">
        <el-icon :size="20" color="#409EFF"><Monitor /></el-icon>
        <h3>系统基本信息</h3>
      </div>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="主机名">{{ systemInfo.hostname }}</el-descriptions-item>
        <el-descriptions-item label="操作系统">{{ systemInfo.os }}</el-descriptions-item>
        <el-descriptions-item label="架构">{{ systemInfo.arch }}</el-descriptions-item>
        <el-descriptions-item label="CPU核心数">{{ systemInfo.cpu_cores }}</el-descriptions-item>
        <el-descriptions-item label="内核版本">{{ systemInfo.kernel_version }}</el-descriptions-item>
        <el-descriptions-item label="CPU使用率">
          <el-progress :percentage="systemInfo.cpu_usage" :format="(val) => val.toFixed(1) + '%'" />
        </el-descriptions-item>
        <el-descriptions-item label="总内存">
          {{ (systemInfo.total_memory / (1024 * 1024 * 1024)).toFixed(2) }} GB
        </el-descriptions-item>
        <el-descriptions-item label="内存使用率">
          <el-progress :percentage="systemInfo.memory_usage" :format="(val) => val.toFixed(1) + '%'" />
        </el-descriptions-item>
      </el-descriptions>
    </div>

    <div class="info-section">
      <div class="section-header">
        <el-icon :size="20" color="#409EFF"><Folder /></el-icon>
        <h3>磁盘信息</h3>
      </div>
      <el-table :data="systemInfo.disks" style="width: 100%" border>
        <el-table-column prop="mount_point" label="挂载点" />
        <el-table-column prop="total_size" label="总大小">
          <template #default="{ row }">
            {{ (row.total_size / (1024 * 1024 * 1024)).toFixed(2) }} GB
          </template>
        </el-table-column>
        <el-table-column prop="used_size" label="已用大小">
          <template #default="{ row }">
            {{ (row.used_size / (1024 * 1024 * 1024)).toFixed(2) }} GB
          </template>
        </el-table-column>
        <el-table-column prop="free_size" label="可用大小">
          <template #default="{ row }">
            {{ (row.free_size / (1024 * 1024 * 1024)).toFixed(2) }} GB
          </template>
        </el-table-column>
        <el-table-column prop="usage" label="使用率">
          <template #default="{ row }">
            <el-progress :percentage="row.usage" :format="(val) => val.toFixed(1) + '%'" />
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<style scoped>
.system-info-panel {
  padding: 0;
}

.info-section {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(0, 0, 0, 0.05);
  margin-bottom: 24px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 20px;
}

.section-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1a202c;
}

:deep(.el-descriptions) {
  background: transparent;
  padding: 0;
  box-shadow: none;
  border: none;
}

:deep(.el-descriptions__label) {
  font-weight: 500;
  color: #4a5568;
  width: 100px;
}

:deep(.el-descriptions__content) {
  color: #2d3748;
}

:deep(.el-progress-bar__outer) {
  background-color: rgba(0, 0, 0, 0.05);
}

:deep(.el-progress-bar__inner) {
  background: linear-gradient(90deg, #409EFF 0%, #66b1ff 100%);
}

:deep(.el-table) {
  background: transparent;
  border-radius: 8px;
  overflow: hidden;
}

:deep(.el-table th) {
  background-color: rgba(0, 0, 0, 0.02);
  font-weight: 600;
  color: #1a202c;
  padding: 12px 16px;
}

:deep(.el-table td) {
  color: #2d3748;
  padding: 12px 16px;
}

:deep(.el-table tr:hover > td) {
  background-color: rgba(64, 158, 255, 0.05);
}
</style> 
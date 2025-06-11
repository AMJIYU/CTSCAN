<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetStartupItems } from '../../wailsjs/go/pkg/App'
import { Document, Folder, Timer, InfoFilled, CopyDocument } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

interface StartupItem {
  name: string
  path: string
  type: string
  enabled: boolean
  lastModTime: string
  size: number
  description: string
}

const startupItems = ref<StartupItem[]>([])

const refresh = async () => {
  try {
    const items = await GetStartupItems()
    startupItems.value = items
  } catch (error) {
    console.error('获取启动项失败:', error)
  }
}

// 格式化文件大小
const formatSize = (size: number) => {
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB'
  return (size / 1024 / 1024).toFixed(1) + ' MB'
}

// 格式化时间
const formatTime = (time: string) => {
  return new Date(time).toLocaleString()
}

// 复制路径到剪贴板
const copyPath = async (path: string) => {
  try {
    await navigator.clipboard.writeText(path)
    ElMessage({
      type: 'success',
      message: '路径已复制到剪贴板',
      duration: 2000
    })
  } catch (err) {
    ElMessage({
      type: 'error',
      message: '复制失败',
      duration: 2000
    })
  }
}

onMounted(() => {
  refresh()
})

defineExpose({ refresh })
</script>

<template>
  <div class="startup-panel">
    <div class="panel-header">
      <el-icon :size="20" color="#409EFF"><Document /></el-icon>
      <h2>开机启动项</h2>
    </div>
    
    <el-table 
      :data="startupItems" 
      style="width: 100%" 
      size="small"
      border
      :resizable="true"
    >
      <el-table-column 
        prop="name" 
        label="名称" 
        min-width="180"
        resizable
        show-overflow-tooltip
      >
        <template #default="{ row }">
          <div class="name-cell">
            <el-icon><Document /></el-icon>
            <span class="text-ellipsis">{{ row.name }}</span>
          </div>
        </template>
      </el-table-column>
      
      <el-table-column 
        prop="type" 
        label="类型" 
        width="120" 
        align="center"
        resizable
      >
        <template #default="{ row }">
          <el-tag size="small" :type="row.type.includes('User') ? 'success' : 'warning'">
            {{ row.type }}
          </el-tag>
        </template>
      </el-table-column>
      
      <el-table-column 
        prop="path" 
        label="路径" 
        min-width="250"
        resizable
        show-overflow-tooltip
      >
        <template #default="{ row }">
          <div class="path-cell">
            <el-icon><Folder /></el-icon>
            <span class="text-ellipsis">{{ row.path }}</span>
            <el-button
              class="copy-button"
              type="primary"
              link
              @click="copyPath(row.path)"
            >
              <el-icon><CopyDocument /></el-icon>
            </el-button>
          </div>
        </template>
      </el-table-column>
      
      <el-table-column 
        label="信息" 
        width="200" 
        align="right"
        resizable
      >
        <template #default="{ row }">
          <div class="info-cell">
            <div class="info-item">
              <el-icon><Timer /></el-icon>
              <span>{{ formatTime(row.lastModTime) }}</span>
            </div>
            <div class="info-item">
              <span class="size">{{ formatSize(row.size) }}</span>
            </div>
          </div>
        </template>
      </el-table-column>

      <el-table-column 
        prop="description" 
        label="描述" 
        min-width="200"
        resizable
        show-overflow-tooltip
      >
        <template #default="{ row }">
          <div class="description-cell" v-if="row.description">
            <el-icon><InfoFilled /></el-icon>
            <span class="text-ellipsis">{{ row.description }}</span>
          </div>
          <span v-else class="no-description">无描述</span>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<style scoped>
.startup-panel {
  padding: 0;
}

.panel-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 20px;
}

.panel-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: #1a202c;
  margin: 0;
}

:deep(.el-table) {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

:deep(.el-table th) {
  background-color: rgba(0, 0, 0, 0.02);
  font-weight: 600;
  color: #1a202c;
  padding: 12px 0;
}

:deep(.el-table td) {
  padding: 12px 0;
  color: #4a5568;
}

:deep(.el-table tr:hover > td) {
  background-color: rgba(64, 158, 255, 0.05);
}

:deep(.el-table__border-left-patch) {
  background-color: rgba(0, 0, 0, 0.02);
}

:deep(.el-table__column-resize-proxy) {
  background-color: #409EFF;
}

:deep(.el-table__column-resize-handle) {
  background-color: #409EFF;
}

:deep(.el-table__column-resize-handle:hover) {
  background-color: #66b1ff;
}

.name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.name-cell .el-icon {
  color: #409EFF;
  font-size: 16px;
  flex-shrink: 0;
}

.text-ellipsis {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 0;
}

.description-cell {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #4a5568;
  min-width: 0;
}

.description-cell .el-icon {
  color: #909399;
  font-size: 14px;
  flex-shrink: 0;
}

.no-description {
  color: #909399;
  font-style: italic;
  font-size: 13px;
}

.path-cell {
  display: flex;
  align-items: center;
  gap: 6px;
  font-family: monospace;
  font-size: 13px;
  color: #2d3748;
  min-width: 0;
  padding-right: 8px;
}

.path-cell .el-icon {
  color: #909399;
  font-size: 14px;
  flex-shrink: 0;
}

.copy-button {
  padding: 2px;
  margin-left: 4px;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.path-cell:hover .copy-button {
  opacity: 1;
}

.copy-button:hover {
  background-color: rgba(64, 158, 255, 0.1);
  border-radius: 4px;
}

.info-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: flex-end;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #718096;
}

.info-item .el-icon {
  font-size: 14px;
}

.size {
  font-family: monospace;
  color: #409EFF;
}

:deep(.el-tag) {
  border-radius: 4px;
  font-size: 12px;
  height: 20px;
  line-height: 18px;
  padding: 0 6px;
}
</style> 
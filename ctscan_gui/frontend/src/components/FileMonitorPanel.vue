<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Document, Timer, User, Folder, Link } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { GetSensitiveFileInfo } from '../../wailsjs/go/pkg/App'

interface FileInfo {
  path: string
  size: number
  mode: string
  mod_time: string
  create_time: string
  access_time: string
  change_time: string
  is_dir: boolean
  is_symlink: boolean
  owner: string
  group: string
  permissions: string
  description: string
}

const files = ref<FileInfo[]>([])
const loading = ref(false)

// 分页相关
const currentPage = ref(1)
const pageSize = ref(10)
const total = computed(() => filteredFiles.value.length)

// 筛选条件
const filters = ref({
  path: ''
})

// 格式化文件大小
const formatFileSize = (size: number) => {
  if (size === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(size) / Math.log(1024))
  return `${(size / Math.pow(1024, i)).toFixed(2)} ${units[i]}`
}

// 格式化时间
const formatTime = (timeStr: string) => {
  if (!timeStr) return { display: '', full: '' }
  
  const date = new Date(timeStr)
  
  // 如果时间无效，返回空对象
  if (isNaN(date.getTime())) return { display: timeStr, full: timeStr }
  
  // 格式化具体时间
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hour = String(date.getHours()).padStart(2, '0')
  const minute = String(date.getMinutes()).padStart(2, '0')
  const second = String(date.getSeconds()).padStart(2, '0')
  
  const formattedTime = `${year}-${month}-${day} ${hour}:${minute}:${second}`
  
  return {
    display: formattedTime,
    full: formattedTime
  }
}

// 获取文件类型图标
const getFileTypeIcon = (file: FileInfo) => {
  if (file.is_dir) return Folder
  if (file.is_symlink) return Link
  return Document
}

// 获取文件类型标签
const getFileTypeTag = (file: FileInfo) => {
  if (file.is_dir) return { type: 'warning', text: '目录' }
  if (file.is_symlink) return { type: 'success', text: '链接' }
  return { type: '', text: '文件' }
}

// 筛选后的数据
const filteredFiles = computed(() => {
  return files.value.filter(file => {
    return !filters.value.path || file.path.toLowerCase().includes(filters.value.path.toLowerCase())
  })
})

// 当前页数据
const currentPageData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredFiles.value.slice(start, end)
})

// 处理页码变化
const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

// 处理每页条数变化
const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
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

// 刷新文件信息
const refreshFiles = async () => {
  loading.value = true
  try {
    files.value = await GetSensitiveFileInfo()
  } catch (error) {
    console.error('获取文件信息失败:', error)
    ElMessage({
      type: 'error',
      message: '获取文件信息失败',
      duration: 2000
    })
  } finally {
    loading.value = false
  }
}

// 时间排序方法
const sortByTime = (a: string, b: string) => {
  const dateA = new Date(a).getTime()
  const dateB = new Date(b).getTime()
  return dateA - dateB
}

onMounted(async () => {
  await refreshFiles()
})
</script>

<template>
  <div class="file-monitor-panel">
    <div class="panel-header">
      <div class="header-left">
        <h2>敏感文件监控</h2>
        <el-tag 
          size="small" 
          type="info" 
          class="record-type-tag"
        >
          系统文件
        </el-tag>
      </div>
      <div class="header-actions">
        <span class="total-count">共 {{ total }} 个文件</span>
        <el-button 
          type="primary" 
          link 
          @click="refreshFiles"
          :loading="loading"
        >
          刷新
        </el-button>
      </div>
    </div>

    <div class="table-container" v-loading="loading">
      <el-table 
        :data="currentPageData" 
        style="width: 100%"
        border
        size="small"
        v-if="files.length > 0"
      >
        <el-table-column prop="path" label="文件路径" min-width="130">
          <template #header>
            <div class="table-header">
              <span>文件路径</span>
              <el-input
                v-model="filters.path"
                placeholder="筛选路径"
                size="small"
                clearable
              />
            </div>
          </template>
          <template #default="{ row }">
            <div class="path-cell">
              <el-icon><component :is="getFileTypeIcon(row)" /></el-icon>
              <div class="path-info">
                <span class="path-text">{{ row.path }}</span>
                <span v-if="row.description" class="file-description">{{ row.description }}</span>
              </div>
              <el-button
                type="primary"
                link
                size="small"
                class="copy-button"
                @click="copyPath(row.path)"
              >
                复制
              </el-button>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="修改时间" min-width="120" sortable :sort-method="(a, b) => sortByTime(a.mod_time, b.mod_time)">
          <template #default="{ row }">
            <div class="time-cell">
              <el-icon><Timer /></el-icon>
              <span class="time-text">{{ formatTime(row.mod_time).display }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="创建时间" min-width="120" sortable :sort-method="(a, b) => sortByTime(a.create_time, b.create_time)">
          <template #default="{ row }">
            <div class="time-cell">
              <el-icon><Timer /></el-icon>
              <span class="time-text">{{ formatTime(row.create_time).display }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="访问时间" min-width="120" sortable :sort-method="(a, b) => sortByTime(a.access_time, b.access_time)">
          <template #default="{ row }">
            <div class="time-cell">
              <el-icon><Timer /></el-icon>
              <span class="time-text">{{ formatTime(row.access_time).display }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="属性修改时间" min-width="120" sortable :sort-method="(a, b) => sortByTime(a.change_time, b.change_time)">
          <template #default="{ row }">
            <div class="time-cell">
              <el-icon><Timer /></el-icon>
              <span class="time-text">{{ formatTime(row.change_time).display }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="类型" width="60">
          <template #default="{ row }">
            <el-tag 
              :type="getFileTypeTag(row).type"
              size="small"
            >
              {{ getFileTypeTag(row).text }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="大小" width="90">
          <template #default="{ row }">
            <span v-if="!row.is_dir">{{ formatFileSize(row.size) }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column label="权限" width="100">
          <template #default="{ row }">
            <el-tooltip
              :content="row.mode"
              placement="top"
              :show-after="500"
            >
              <span>{{ row.permissions }}</span>
            </el-tooltip>
          </template>
        </el-table-column>

        <el-table-column label="所有者" width="150">
          <template #default="{ row }">
            <div class="owner-cell">
              <el-icon><User /></el-icon>
              <div class="owner-info">
                <span class="owner-name">{{ row.owner }}</span>
                <span class="group-name">{{ row.group }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <el-empty 
        v-else 
        description="暂无敏感文件信息"
      />
    </div>

    <div class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<style scoped>
.file-monitor-panel {
  padding: 0;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.panel-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: #1a202c;
  margin: 0;
}

.record-type-tag {
  font-size: 12px;
  height: 20px;
  line-height: 18px;
  padding: 0 6px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.total-count {
  color: #909399;
  font-size: 14px;
}

.table-container {
  border-radius: 8px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

:deep(.el-table) {
  --el-table-border-color: rgba(0, 0, 0, 0.05);
  --el-table-header-bg-color: rgba(0, 0, 0, 0.02);
}

:deep(.el-table th) {
  font-weight: 600;
  color: #1a202c;
}

:deep(.el-table td) {
  color: #4a5568;
}

:deep(.el-table tr:hover > td) {
  background-color: rgba(64, 158, 255, 0.05);
}

.path-cell,
.time-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.path-cell .el-icon,
.time-cell .el-icon {
  color: #909399;
  font-size: 16px;
}

.path-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
}

:deep(.el-tag) {
  border-radius: 4px;
  font-size: 12px;
  height: 20px;
  line-height: 18px;
  padding: 0 6px;
}

:deep(.el-loading-mask) {
  backdrop-filter: blur(2px);
}

:deep(.el-loading-spinner) {
  .el-loading-text {
    color: #409EFF;
    font-size: 14px;
    margin-top: 8px;
  }
  
  .circular {
    width: 30px;
    height: 30px;
  }
}

.table-header {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.table-header span {
  font-weight: bold;
  color: #606266;
}

.path-cell {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  position: relative;
  padding-right: 60px;
}

.copy-button {
  opacity: 0;
  transition: opacity 0.2s ease;
  flex-shrink: 0;
  position: absolute;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
}

.path-cell:hover .copy-button {
  opacity: 1;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  padding: 10px 0;
}

:deep(.el-pagination) {
  --el-pagination-button-color: #409EFF;
  --el-pagination-hover-color: #66b1ff;
}

:deep(.el-pagination .el-select .el-input) {
  width: 110px;
}

:deep(.el-pagination .el-pagination__total) {
  margin-right: 16px;
}

:deep(.el-pagination .el-pagination__sizes) {
  margin-right: 16px;
}

:deep(.el-pagination .el-pagination__jump) {
  margin-left: 16px;
}

:deep(.el-pagination .el-pagination__jump .el-input__inner) {
  text-align: center;
}

.time-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.time-cell .el-icon {
  color: #909399;
  font-size: 16px;
}

.time-text {
  color: #606266;
  font-size: 13px;
  font-family: monospace;
}

.owner-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.owner-cell .el-icon {
  color: #909399;
  font-size: 16px;
}

.owner-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.owner-name {
  color: #303133;
  font-size: 13px;
  font-weight: 500;
}

.group-name {
  color: #909399;
  font-size: 12px;
}

.path-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.file-description {
  color: #909399;
  font-size: 12px;
  font-style: italic;
}
</style> 
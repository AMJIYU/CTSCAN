<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { GetAllProcesses } from '../../wailsjs/go/pkg/App'
import { Monitor, Document, Connection, Timer, Key } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

interface ProcessInfo {
  pid: number
  name: string
  ppid: number
  parent_name: string
  create_time: number
  exe: string
  file_ctime: number
  file_mtime: number
  md5: string
  signature: string
  cpu_percent: number
  mem_percent: number
}

const processes = ref<ProcessInfo[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const sortProp = ref('cpu_percent')
const sortOrder = ref('descending')
const loading = ref(false)
const isFirstLoad = ref(true)

// 筛选条件
const filters = ref({
  pid: '',
  name: '',
  ppid: '',
  parent_name: '',
  exe: '',
  md5: '',
  signature: ''
})

// 重置筛选条件
const resetFilters = () => {
  filters.value = {
    pid: '',
    name: '',
    ppid: '',
    parent_name: '',
    exe: '',
    md5: '',
    signature: ''
  }
  currentPage.value = 1
}

// 筛选后的数据
const filteredProcesses = computed(() => {
  return processes.value.filter(process => {
    return (
      (!filters.value.pid || process.pid.toString().includes(filters.value.pid)) &&
      (!filters.value.name || process.name.toLowerCase().includes(filters.value.name.toLowerCase())) &&
      (!filters.value.ppid || process.ppid.toString().includes(filters.value.ppid)) &&
      (!filters.value.parent_name || (process.parent_name && process.parent_name.toLowerCase().includes(filters.value.parent_name.toLowerCase()))) &&
      (!filters.value.exe || process.exe.toLowerCase().includes(filters.value.exe.toLowerCase())) &&
      (!filters.value.md5 || (process.md5 && process.md5.toLowerCase().includes(filters.value.md5.toLowerCase()))) &&
      (!filters.value.signature || (process.signature && process.signature.toLowerCase().includes(filters.value.signature.toLowerCase())))
    )
  })
})

// 计算排序后的数据
const sortedProcesses = computed(() => {
  const sorted = [...filteredProcesses.value]
  if (sortProp.value && sortOrder.value) {
    sorted.sort((a, b) => {
      const aValue = a[sortProp.value as keyof ProcessInfo]
      const bValue = b[sortProp.value as keyof ProcessInfo]
      if (sortOrder.value === 'ascending') {
        return aValue > bValue ? 1 : -1
      } else {
        return aValue < bValue ? 1 : -1
      }
    })
  }
  return sorted
})

// 计算当前页的数据
const currentPageData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return sortedProcesses.value.slice(start, end)
})

// 监听筛选条件变化
watch(filters, () => {
  total.value = filteredProcesses.value.length
  if (currentPage.value > Math.ceil(total.value / pageSize.value)) {
    currentPage.value = 1
  }
}, { deep: true })

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

const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
}

const handleSortChange = ({ prop, order }: { prop: string, order: string }) => {
  sortProp.value = prop
  sortOrder.value = order
}

const refresh = () => {
  // 如果不是第一次加载，直接返回
  if (!isFirstLoad.value) {
    return
  }
  
  loading.value = true
  GetAllProcesses().then(list => {
    processes.value = list
    total.value = list.length
    isFirstLoad.value = false
  }).finally(() => {
    loading.value = false
  })
}

// 强制刷新方法
const forceRefresh = () => {
  loading.value = true
  GetAllProcesses().then(list => {
    processes.value = list
    total.value = list.length
  }).finally(() => {
    loading.value = false
  })
}

// 格式化时间戳
const formatTimestamp = (timestamp: number): string => {
  if (!timestamp) return '-'
  // 将毫秒级时间戳转换为日期对象
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  })
}

onMounted(() => {
  refresh()
})

defineExpose({ refresh, forceRefresh })
</script>

<template>
  <div class="process-panel">
    <div class="panel-header">
      <el-icon :size="20" color="#409EFF"><Monitor /></el-icon>
      <h2>进程信息</h2>
      <span class="total-count">共 {{ total }} 个进程</span>
      <div class="header-actions">
        <el-button 
          type="primary" 
          link 
          @click="resetFilters"
          class="reset-button"
        >
          重置筛选
        </el-button>
        <el-button 
          type="primary" 
          link 
          @click="forceRefresh"
          :loading="loading"
          class="refresh-button"
        >
          刷新
        </el-button>
      </div>
    </div>

    <el-table 
      v-loading="loading"
      element-loading-text="正在加载进程信息..."
      element-loading-background="rgba(255, 255, 255, 0.9)"
      :data="currentPageData" 
      style="width: 100%" 
      border
      :resizable="true"
      size="small"
      @sort-change="handleSortChange"
    >
      <el-table-column 
        prop="pid" 
        label="PID" 
        width="70"
        align="center"
        resizable
      >
        <template #header>
          <div class="table-header">
            <span>PID</span>
            <el-input
              v-model="filters.pid"
              placeholder="筛选PID"
              size="small"
              clearable
            />
          </div>
        </template>
        <template #default="{ row }">
          <span class="pid-value">{{ row.pid }}</span>
        </template>
      </el-table-column>

      <el-table-column 
        prop="name" 
        label="进程名" 
        min-width="100"
        resizable
        show-overflow-tooltip
      >
        <template #header>
          <div class="table-header">
            <span>进程名</span>
            <el-input
              v-model="filters.name"
              placeholder="筛选进程名"
              size="small"
              clearable
            />
          </div>
        </template>
        <template #default="{ row }">
          <div class="process-name">
            <el-icon><Document /></el-icon>
            <span>{{ row.name }}</span>
          </div>
        </template>
      </el-table-column>

      <el-table-column 
        prop="exe" 
        label="可执行文件路径" 
        min-width="180"
        resizable
        show-overflow-tooltip
      >
        <template #header>
          <div class="table-header">
            <span>可执行文件路径</span>
            <el-input
              v-model="filters.exe"
              placeholder="筛选路径"
              size="small"
              clearable
            />
          </div>
        </template>
        <template #default="{ row }">
          <div class="path-cell">
            <span class="path-text">{{ row.exe }}</span>
            <el-button
              type="primary"
              link
              size="small"
              class="copy-button"
              @click="copyPath(row.exe)"
            >
              复制
            </el-button>
          </div>
        </template>
      </el-table-column>

      <el-table-column 
        prop="signature" 
        label="签名信息" 
        min-width="150"
        resizable
        show-overflow-tooltip
      >
        <template #header>
          <div class="table-header">
            <span>签名信息</span>
            <el-input
              v-model="filters.signature"
              placeholder="筛选签名"
              size="small"
              clearable
            />
          </div>
        </template>
        <template #default="{ row }">
          <div v-if="row.signature" class="signature-cell">
            <el-icon><Key /></el-icon>
            <span>{{ row.signature }}</span>
          </div>
          <span v-else class="no-value">-</span>
        </template>
      </el-table-column>

      <el-table-column 
        prop="ppid" 
        label="父PID" 
        width="70"
        align="center"
        resizable
      >
        <template #header>
          <div class="table-header">
            <span>父PID</span>
            <el-input
              v-model="filters.ppid"
              placeholder="筛选父PID"
              size="small"
              clearable
            />
          </div>
        </template>
        <template #default="{ row }">
          <span class="pid-value">{{ row.ppid }}</span>
        </template>
      </el-table-column>

      <el-table-column 
        prop="parent_name" 
        label="父进程名" 
        min-width="100"
        resizable
        show-overflow-tooltip
      >
        <template #header>
          <div class="table-header">
            <span>父进程名</span>
            <el-input
              v-model="filters.parent_name"
              placeholder="筛选父进程名"
              size="small"
              clearable
            />
          </div>
        </template>
        <template #default="{ row }">
          <div class="process-name" v-if="row.parent_name">
            <el-icon><Document /></el-icon>
            <span>{{ row.parent_name }}</span>
          </div>
          <span v-else class="no-parent">-</span>
        </template>
      </el-table-column>

      <el-table-column 
        prop="md5" 
        label="MD5" 
        min-width="120"
        resizable
        show-overflow-tooltip
      >
        <template #header>
          <div class="table-header">
            <span>MD5</span>
            <el-input
              v-model="filters.md5"
              placeholder="筛选MD5"
              size="small"
              clearable
            />
          </div>
        </template>
        <template #default="{ row }">
          <span v-if="row.md5" class="md5-value">{{ row.md5 }}</span>
          <span v-else class="no-value">-</span>
        </template>
      </el-table-column>

      <el-table-column 
        prop="cpu_percent" 
        label="CPU使用率" 
        width="90"
        align="center"
        sortable="custom"
        resizable
      >
        <template #default="{ row }">
          <el-tag 
            size="small" 
            :type="row.cpu_percent > 80 ? 'danger' : row.cpu_percent > 50 ? 'warning' : 'info'"
          >
            {{ row.cpu_percent ? row.cpu_percent.toFixed(1) + '%' : '0%' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column 
        prop="mem_percent" 
        label="内存使用率" 
        width="90"
        align="center"
        sortable="custom"
        resizable
      >
        <template #default="{ row }">
          <el-tag 
            size="small" 
            :type="row.mem_percent > 80 ? 'danger' : row.mem_percent > 50 ? 'warning' : 'info'"
          >
            {{ row.mem_percent ? row.mem_percent.toFixed(1) + '%' : '0%' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column 
        prop="create_time" 
        label="创建时间" 
        min-width="140"
        sortable
        resizable
        show-overflow-tooltip
      >
        <template #default="{ row }">
          <div class="time-value">
            <el-icon><Timer /></el-icon>
            <span>{{ formatTimestamp(row.create_time) }}</span>
          </div>
        </template>
      </el-table-column>
  </el-table>

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
.process-panel {
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

.total-count {
  margin-left: auto;
  color: #909399;
  font-size: 14px;
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

.process-name {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.process-name .el-icon {
  color: #409EFF;
  font-size: 16px;
  flex-shrink: 0;
}

.process-name span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.pid-value {
  font-family: monospace;
  font-size: 13px;
  color: #409EFF;
  font-weight: 500;
}

.no-parent {
  color: #909399;
  font-style: italic;
}

.time-value {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #2d3748;
}

.time-value .el-icon {
  color: #909399;
  font-size: 14px;
}

.path-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  position: relative;
  padding-right: 60px;
}

.path-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
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

.md5-value {
  font-family: monospace;
  font-size: 13px;
  color: #2d3748;
}

.signature-cell {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #2d3748;
}

.signature-cell .el-icon {
  color: #909399;
  font-size: 14px;
  flex-shrink: 0;
}

.no-value {
  color: #909399;
  font-style: italic;
}

:deep(.el-tag) {
  border-radius: 4px;
  font-size: 12px;
  height: 20px;
  line-height: 18px;
  padding: 0 6px;
  min-width: 60px;
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

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-left: auto;
}

.reset-button,
.refresh-button {
  margin-left: 0;
}
</style> 
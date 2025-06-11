<template>
  <div class="shell-history-panel">
    <div class="panel-header">
      <div class="header-left">
        <h2>命令执行记录</h2>
        <el-tag 
          size="small" 
          type="info" 
          class="record-type-tag"
        >
          历史记录
        </el-tag>
      </div>
      <div class="header-actions">
        <span class="total-count">共 {{ total }} 条记录</span>
        <el-button 
          type="primary" 
          link 
          @click="resetFilters"
        >
          重置筛选
        </el-button>
        <el-button 
          type="primary" 
          link 
          @click="refresh"
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
        v-if="records.length > 0"
      >
        <el-table-column prop="time" label="时间" min-width="180">
          <template #header>
            <div class="table-header">
              <span>时间</span>
              <el-input
                v-model="filters.time"
                placeholder="筛选时间"
                size="small"
                clearable
              />
            </div>
          </template>
          <template #default="{ row }">
            <div class="time-cell">
              <el-icon><Timer /></el-icon>
              <span>{{ row.time }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="command" label="命令" min-width="300">
          <template #header>
            <div class="table-header">
              <span>命令</span>
              <el-input
                v-model="filters.command"
                placeholder="筛选命令"
                size="small"
                clearable
              />
            </div>
          </template>
          <template #default="{ row }">
            <div class="command-cell">
              <el-icon><Operation /></el-icon>
              <span class="command-text">{{ row.command }}</span>
              <el-button
                type="primary"
                link
                size="small"
                class="copy-button"
                @click="copyCommand(row.command)"
              >
                复制
              </el-button>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="user" label="用户" min-width="120">
          <template #header>
            <div class="table-header">
              <span>用户</span>
              <el-input
                v-model="filters.user"
                placeholder="筛选用户"
                size="small"
                clearable
              />
            </div>
          </template>
          <template #default="{ row }">
            <div class="user-cell">
              <el-icon><User /></el-icon>
              <span>{{ row.user }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="shell" label="Shell" min-width="120">
          <template #header>
            <div class="table-header">
              <span>Shell</span>
              <el-input
                v-model="filters.shell"
                placeholder="筛选Shell"
                size="small"
                clearable
              />
            </div>
          </template>
          <template #default="{ row }">
            <el-tag 
              size="small" 
              :type="getShellTagType(row.shell)"
            >
              {{ row.shell }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      
      <el-empty 
        v-else 
        description="暂无命令执行记录"
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

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Timer, Operation, User } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { GetShellHistory } from '../../wailsjs/go/pkg/App'

interface ShellHistory {
  time: string;
  command: string;
  user: string;
  shell: string;
}

const records = ref<ShellHistory[]>([])
const loading = ref(false)

// 分页相关
const currentPage = ref(1)
const pageSize = ref(10)
const total = computed(() => filteredRecords.value.length)

// 筛选条件
const filters = ref({
  time: '',
  command: '',
  user: '',
  shell: ''
})

// 重置筛选条件
const resetFilters = () => {
  filters.value = {
    time: '',
    command: '',
    user: '',
    shell: ''
  }
  currentPage.value = 1
}

// 筛选后的数据
const filteredRecords = computed(() => {
  return records.value.filter(record => {
    return (
      (!filters.value.time || record.time.includes(filters.value.time)) &&
      (!filters.value.command || record.command.toLowerCase().includes(filters.value.command.toLowerCase())) &&
      (!filters.value.user || record.user.toLowerCase().includes(filters.value.user.toLowerCase())) &&
      (!filters.value.shell || record.shell.toLowerCase().includes(filters.value.shell.toLowerCase()))
    )
  })
})

// 当前页数据
const currentPageData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredRecords.value.slice(start, end)
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

// 获取 Shell 标签类型
const getShellTagType = (shell: string) => {
  switch (shell.toLowerCase()) {
    case 'bash':
      return 'success'
    case 'zsh':
      return 'warning'
    case 'fish':
      return 'info'
    case 'powershell':
      return 'danger'
    case 'cmd':
      return ''
    default:
      return 'info'
  }
}

// 复制命令到剪贴板
const copyCommand = async (command: string) => {
  try {
    await navigator.clipboard.writeText(command)
    ElMessage({
      type: 'success',
      message: '命令已复制到剪贴板',
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

// 获取记录
const refresh = async () => {
  loading.value = true
  try {
    const response = await GetShellHistory()
    records.value = response
  } catch (error) {
    console.error('获取命令执行记录失败:', error)
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

<style scoped>
.shell-history-panel {
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

.time-cell,
.command-cell,
.user-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.time-cell .el-icon,
.command-cell .el-icon,
.user-cell .el-icon {
  color: #909399;
  font-size: 16px;
}

.command-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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

.command-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  position: relative;
  padding-right: 60px;
}

.command-text {
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

.command-cell:hover .copy-button {
  opacity: 1;
}
</style>

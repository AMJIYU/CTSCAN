<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { GetLoginSuccessRecords, SaveLoginSuccess } from '../../wailsjs/go/pkg/App'
import { ElMessage } from 'element-plus'
import { Timer, User, Location, Document, Key } from '@element-plus/icons-vue'

interface LoginSuccess {
  time: string
  event_id: string
  event_type: string
  source: string
  username: string
  ip_address: string
}

const loginSuccessRecords = ref<LoginSuccess[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 筛选条件
const filters = ref({
  time: '',
  event_type: '',
  source: '',
  username: '',
  ip_address: ''
})

// 重置筛选条件
const resetFilters = () => {
  filters.value = {
    time: '',
    event_type: '',
    source: '',
    username: '',
    ip_address: ''
  }
  currentPage.value = 1
}

// 筛选后的数据
const filteredRecords = computed(() => {
  return loginSuccessRecords.value.filter(record => {
    return (
      (!filters.value.time || formatTime(record.time).includes(filters.value.time)) &&
      (!filters.value.event_type || record.event_type.includes(filters.value.event_type)) &&
      (!filters.value.source || record.source.toLowerCase().includes(filters.value.source.toLowerCase())) &&
      (!filters.value.username || record.username.toLowerCase().includes(filters.value.username.toLowerCase())) &&
      (!filters.value.ip_address || (record.ip_address && record.ip_address.toLowerCase().includes(filters.value.ip_address.toLowerCase())))
    )
  })
})

// 计算当前页的数据
const currentPageData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredRecords.value.slice(start, end)
})

// 监听筛选条件变化
watch(filters, () => {
  total.value = filteredRecords.value.length
  if (currentPage.value > Math.ceil(total.value / pageSize.value)) {
    currentPage.value = 1
  }
}, { deep: true })

// 格式化时间
const formatTime = (time: string) => {
  if (!time) return '-'
  const date = new Date(time)
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

// 添加 refresh 方法，用于重新获取登录成功记录
const refresh = async () => {
  loading.value = true
  try {
    const records = await GetLoginSuccessRecords()
    loginSuccessRecords.value = records
    // 保存到数据库
    await SaveLoginSuccess(records)
    total.value = records.length
  } catch (error) {
    console.error('获取登录成功记录失败:', error)
  } finally {
    loading.value = false
  }
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
}

onMounted(() => {
  refresh()
})

// 暴露 refresh 方法，供父组件调用
defineExpose({ refresh })
</script>

<template>
  <div class="login-success-panel">
    <div class="info-section">
      <div class="section-header">
        <el-icon :size="20" color="#409EFF"><Key /></el-icon>
        <h3>登录成功记录</h3>
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
      <el-table 
        :data="currentPageData" 
        style="width: 100%" 
        border
        v-loading="loading"
        element-loading-text="正在加载登录成功记录..."
        element-loading-background="rgba(255, 255, 255, 0.8)"
      >
        <el-table-column prop="time" label="时间" width="180" show-overflow-tooltip>
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
              <span>{{ formatTime(row.time) }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="username" label="用户名" width="120" show-overflow-tooltip>
          <template #header>
            <div class="table-header">
              <span>用户名</span>
              <el-input
                v-model="filters.username"
                placeholder="筛选用户名"
                size="small"
                clearable
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="ip_address" label="IP地址" width="140" show-overflow-tooltip>
          <template #header>
            <div class="table-header">
              <span>IP地址</span>
              <el-input
                v-model="filters.ip_address"
                placeholder="筛选IP"
                size="small"
                clearable
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="source" label="来源" width="120" show-overflow-tooltip>
          <template #header>
            <div class="table-header">
              <span>来源</span>
              <el-input
                v-model="filters.source"
                placeholder="筛选来源"
                size="small"
                clearable
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="event_type" label="事件类型" width="120" show-overflow-tooltip>
          <template #header>
            <div class="table-header">
              <span>事件类型</span>
              <el-input
                v-model="filters.event_type"
                placeholder="筛选类型"
                size="small"
                clearable
              />
            </div>
          </template>
        </el-table-column>
      </el-table>
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
.login-success-panel {
  padding: 0;
}

.info-section {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(0, 0, 0, 0.05);
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

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.total-count {
  color: #909399;
  font-size: 14px;
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

.time-cell,
.source-cell,
.username-cell,
.ip-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.time-cell .el-icon,
.source-cell .el-icon,
.username-cell .el-icon,
.ip-cell .el-icon {
  color: #909399;
  font-size: 16px;
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

:deep(.el-loading-spinner .el-loading-text) {
  color: #409EFF;
  font-size: 14px;
  margin-top: 8px;
}

:deep(.el-loading-spinner .circular) {
  width: 30px;
  height: 30px;
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
</style> 
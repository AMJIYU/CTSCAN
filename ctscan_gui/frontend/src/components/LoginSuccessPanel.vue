<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { GetLoginSuccessRecords } from '../../wailsjs/go/pkg/App'
import { ElMessage } from 'element-plus'
import { Timer, User, Location, Document } from '@element-plus/icons-vue'

const records = ref<any[]>([])
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
  return records.value.filter(record => {
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

// 添加 refresh 方法，用于重新获取登录成功日志信息
const refresh = async () => {
  loading.value = true
  try {
    const logs = await GetLoginSuccessRecords()
    records.value = logs
    total.value = logs.length
  } catch (error) {
    ElMessage.error('获取登录记录失败')
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
    <div class="panel-header">
      <h2>登录成功记录</h2>
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
              <span>{{ formatTime(row.time) }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="event_type" label="事件类型" width="120">
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
          <template #default="{ row }">
            <el-tag 
              size="small" 
              :type="row.event_type === '当前登录' ? 'success' : 'info'"
            >
              {{ row.event_type }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="source" label="来源" min-width="200">
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
          <template #default="{ row }">
            <div class="source-cell">
              <el-icon><Document /></el-icon>
              <span>{{ row.source }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="username" label="用户名" min-width="120">
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
          <template #default="{ row }">
            <div class="username-cell">
              <el-icon><User /></el-icon>
              <span>{{ row.username || '-' }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="ip_address" label="IP地址" min-width="140">
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
          <template #default="{ row }">
            <div class="ip-cell">
              <el-icon><Location /></el-icon>
              <span>{{ row.ip_address || '-' }}</span>
            </div>
          </template>
        </el-table-column>
      </el-table>
      
      <el-empty 
        v-else 
        description="暂无登录记录"
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
.login-success-panel {
  padding: 0;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.panel-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: #1a202c;
  margin: 0;
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
</style> 
<template>
  <div class="login-failed-panel">
    <div class="panel-header">
      <div class="header-left">
        <h2>登录失败记录</h2>
        <el-tag 
          size="small" 
          type="danger" 
          class="record-type-tag"
        >
          失败
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
          @click="fetchRecords"
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
              type="danger"
            >
              {{ row.event_type }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="source" label="来源" min-width="120">
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
        
        <el-table-column prop="reason" label="失败原因" min-width="200">
          <template #header>
            <div class="table-header">
              <span>失败原因</span>
              <el-input
                v-model="filters.reason"
                placeholder="筛选原因"
                size="small"
                clearable
              />
            </div>
          </template>
          <template #default="{ row }">
            <div v-if="row.reason" class="reason-cell">
              <el-tag 
                size="small" 
                type="danger"
                effect="plain"
              >
                {{ row.reason }}
              </el-tag>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
      </el-table>
      
      <el-empty 
        v-else 
        description="暂无登录失败记录"
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
import { Timer, Document, User, Location } from '@element-plus/icons-vue'

const records = ref([])
const loading = ref(false)

// 分页相关
const currentPage = ref(1)
const pageSize = ref(10)
const total = computed(() => filteredRecords.value.length)

// 筛选条件
const filters = ref({
  time: '',
  event_type: '',
  source: '',
  username: '',
  ip_address: '',
  reason: ''
})

// 重置筛选条件
const resetFilters = () => {
  filters.value = {
    time: '',
    event_type: '',
    source: '',
    username: '',
    ip_address: '',
    reason: ''
  }
  currentPage.value = 1
}

// 筛选后的数据
const filteredRecords = computed(() => {
  return records.value.filter(record => {
    return (
      (!filters.value.time || record.time.includes(filters.value.time)) &&
      (!filters.value.event_type || record.event_type.includes(filters.value.event_type)) &&
      (!filters.value.source || record.source.toLowerCase().includes(filters.value.source.toLowerCase())) &&
      (!filters.value.username || record.username.toLowerCase().includes(filters.value.username.toLowerCase())) &&
      (!filters.value.ip_address || record.ip_address.toLowerCase().includes(filters.value.ip_address.toLowerCase())) &&
      (!filters.value.reason || (record.reason && record.reason.toLowerCase().includes(filters.value.reason.toLowerCase())))
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

// 获取记录
const fetchRecords = async () => {
  loading.value = true
  try {
    const response = await window.go.main.App.GetLoginFailedRecords()
    records.value = response
  } catch (error) {
    console.error('获取登录失败记录失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchRecords()
})
</script>

<style scoped>
.login-failed-panel {
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
  background-color: rgba(245, 108, 108, 0.05);
}

.time-cell,
.source-cell,
.username-cell,
.ip-cell,
.reason-cell {
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

:deep(.el-tag--danger.el-tag--plain) {
  background-color: rgba(245, 108, 108, 0.1);
  border-color: rgba(245, 108, 108, 0.2);
  color: #f56c6c;
}

:deep(.el-loading-mask) {
  backdrop-filter: blur(2px);
}

:deep(.el-loading-spinner) {
  .el-loading-text {
    color: #F56C6C;
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
  --el-pagination-button-color: #F56C6C;
  --el-pagination-hover-color: #f78989;
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

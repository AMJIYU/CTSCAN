<script setup lang="ts">
import { ref, computed, onMounted, defineExpose } from 'vue'
import { Timer, User, Location, InfoFilled, Monitor } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { GetRDPLoginLogs, SaveRDPLogin } from '../../wailsjs/go/pkg/App'

interface RDPLoginInfo {
  time: string
  username: string
  ip: string
  status: string
  description: string
}

const logs = ref<RDPLoginInfo[]>([])
const loading = ref(false)

// 分页相关
const currentPage = ref(1)
const pageSize = ref(10)
const total = computed(() => filteredLogs.value.length)

// 筛选条件
const filters = ref({
  username: '',
  ip: '',
  status: ''
})

// 筛选后的数据
const filteredLogs = computed(() => {
  return logs.value.filter(log => {
    return (
      (!filters.value.username || log.username.toLowerCase().includes(filters.value.username.toLowerCase())) &&
      (!filters.value.ip || log.ip.includes(filters.value.ip)) &&
      (!filters.value.status || log.status === filters.value.status)
    )
  })
})

// 当前页数据
const currentPageData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredLogs.value.slice(start, end)
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

// 刷新日志
const refresh = async () => {
  loading.value = true
  try {
    const response = await GetRDPLoginLogs()
    logs.value = response
    await SaveRDPLogin(response)
  } catch (error) {
    console.error('获取RDP登录记录失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await refresh()
})

defineExpose({ refresh })
</script>

<template>
  <div class="rdp-login-panel">
    <div class="panel-header">
      <div class="header-left">
        <h2>RDP登录日志</h2>
        <el-tag 
          size="small" 
          type="info" 
          class="record-type-tag"
        >
          远程登录
        </el-tag>
      </div>
      <div class="header-actions">
        <span class="total-count">共 {{ total }} 条记录</span>
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
        v-if="logs.length > 0"
      >
        <el-table-column label="时间" min-width="180" sortable>
          <template #default="{ row }">
            <div class="time-cell">
              <el-icon><Timer /></el-icon>
              <span class="time-text">{{ row.time }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="用户名" min-width="120">
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
            <div class="user-cell">
              <el-icon><User /></el-icon>
              <span>{{ row.username }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="IP地址" min-width="150">
          <template #header>
            <div class="table-header">
              <span>IP地址</span>
              <el-input
                v-model="filters.ip"
                placeholder="筛选IP"
                size="small"
                clearable
              />
            </div>
          </template>
          <template #default="{ row }">
            <div class="ip-cell">
              <el-icon><Location /></el-icon>
              <span>{{ row.ip }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100">
          <template #header>
            <div class="table-header">
              <span>状态</span>
              <el-select
                v-model="filters.status"
                placeholder="全部"
                size="small"
                clearable
              >
                <el-option label="成功" value="成功" />
                <el-option label="失败" value="失败" />
              </el-select>
            </div>
          </template>
          <template #default="{ row }">
            <el-tag 
              :type="row.status === '成功' ? 'success' : 'danger'"
              size="small"
            >
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="描述" min-width="300">
          <template #default="{ row }">
            <div class="description-cell">
              <el-icon><InfoFilled /></el-icon>
              <span>{{ row.description }}</span>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <el-empty 
        v-else 
        :description="loading ? '正在加载日志...' : '暂无RDP登录日志'"
      >
        <template #image>
          <el-icon :size="60" color="#909399"><Monitor /></el-icon>
        </template>
        <template #description>
          <p>{{ loading ? '正在加载日志...' : '当前系统未检测到RDP登录记录' }}</p>
        </template>
      </el-empty>
    </div>

    <div class="pagination-container" v-if="logs.length > 0">
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
.rdp-login-panel {
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
.user-cell,
.ip-cell,
.description-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.time-cell .el-icon,
.user-cell .el-icon,
.ip-cell .el-icon,
.description-cell .el-icon {
  color: #909399;
  font-size: 16px;
}

.time-text {
  color: #606266;
  font-size: 13px;
  font-family: monospace;
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

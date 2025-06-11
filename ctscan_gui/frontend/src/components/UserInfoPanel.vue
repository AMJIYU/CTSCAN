<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { GetUserInfo, GetAllUsers } from '../../wailsjs/go/pkg/App'
import { User, UserFilled, CopyDocument } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const userInfo = ref({
  username: '',
  uid: '',
  gid: '',
  home_dir: '',
  name: ''
})

const allUsers = ref<{ username: string; uid: string; gid: string; home_dir: string; name: string }[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)

// 筛选条件
const filters = ref({
  username: '',
  uid: '',
  gid: '',
  home_dir: '',
  name: ''
})

// 重置筛选条件
const resetFilters = () => {
  filters.value = {
    username: '',
    uid: '',
    gid: '',
    home_dir: '',
    name: ''
  }
  currentPage.value = 1
}

// 筛选后的数据
const filteredUsers = computed(() => {
  return allUsers.value.filter(user => {
    return (
      (!filters.value.username || user.username.toLowerCase().includes(filters.value.username.toLowerCase())) &&
      (!filters.value.uid || user.uid.toLowerCase().includes(filters.value.uid.toLowerCase())) &&
      (!filters.value.gid || user.gid.toLowerCase().includes(filters.value.gid.toLowerCase())) &&
      (!filters.value.home_dir || user.home_dir.toLowerCase().includes(filters.value.home_dir.toLowerCase())) &&
      (!filters.value.name || user.name.toLowerCase().includes(filters.value.name.toLowerCase()))
    )
  })
})

// 计算当前页的数据
const currentPageData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredUsers.value.slice(start, end)
})

const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
}

// 添加 refresh 方法，用于重新获取用户信息
const refresh = async () => {
  loading.value = true
  try {
    const [info, users] = await Promise.all([
      GetUserInfo(),
      GetAllUsers()
    ])
    userInfo.value = info
    allUsers.value = users
    total.value = users.length
  } catch (error) {
    console.error('获取用户信息失败:', error)
  } finally {
    loading.value = false
  }
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

// 监听筛选条件变化
watch(filters, () => {
  total.value = filteredUsers.value.length
  if (currentPage.value > Math.ceil(total.value / pageSize.value)) {
    currentPage.value = 1
  }
}, { deep: true })

onMounted(() => {
  refresh()
})

// 暴露 refresh 方法，供父组件调用
defineExpose({ refresh })
</script>

<template>
  <div class="user-info-panel">
    <!-- 当前用户信息 -->
    <div class="info-section">
      <div class="section-header">
        <el-icon :size="20" color="#409EFF"><User /></el-icon>
        <h3>当前用户信息</h3>
      </div>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="用户名">{{ userInfo.username }}</el-descriptions-item>
        <el-descriptions-item label="UID">{{ userInfo.uid }}</el-descriptions-item>
        <el-descriptions-item label="GID">{{ userInfo.gid }}</el-descriptions-item>
        <el-descriptions-item label="家目录">
          <span class="path-value">{{ userInfo.home_dir }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="显示名">{{ userInfo.name }}</el-descriptions-item>
      </el-descriptions>
    </div>

    <!-- 系统用户列表 -->
    <div class="info-section">
      <div class="section-header">
        <el-icon :size="20" color="#409EFF"><UserFilled /></el-icon>
        <h3>系统用户列表</h3>
        <span class="total-count">共 {{ total }} 个用户</span>
        <el-button 
          type="primary" 
          link 
          @click="resetFilters"
          class="reset-button"
        >
          重置筛选
        </el-button>
      </div>
      <el-table 
        v-loading="loading"
        element-loading-text="正在加载用户信息..."
        element-loading-background="rgba(255, 255, 255, 0.9)"
        :data="currentPageData" 
        style="width: 100%"
        border
        :resizable="true"
        size="small"
      >
        <el-table-column 
          prop="username" 
          label="用户名" 
          min-width="120"
          resizable
        >
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
        <el-table-column 
          prop="uid" 
          label="UID" 
          min-width="80"
          align="center"
          resizable
        >
          <template #header>
            <div class="table-header">
              <span>UID</span>
              <el-input
                v-model="filters.uid"
                placeholder="筛选UID"
                size="small"
                clearable
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column 
          prop="gid" 
          label="GID" 
          min-width="80"
          align="center"
          resizable
        >
          <template #header>
            <div class="table-header">
              <span>GID</span>
              <el-input
                v-model="filters.gid"
                placeholder="筛选GID"
                size="small"
                clearable
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column 
          prop="home_dir" 
          label="家目录" 
          min-width="200"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="table-header">
              <span>家目录</span>
              <el-input
                v-model="filters.home_dir"
                placeholder="筛选家目录"
                size="small"
                clearable
              />
            </div>
          </template>
          <template #default="{ row }">
            <div class="path-cell">
              <span class="path-text">{{ row.home_dir }}</span>
              <el-button
                type="primary"
                link
                size="small"
                class="copy-button"
                @click="copyPath(row.home_dir)"
              >
                复制
              </el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column 
          prop="name" 
          label="显示名" 
          min-width="150"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="table-header">
              <span>显示名</span>
              <el-input
                v-model="filters.name"
                placeholder="筛选显示名"
                size="small"
                clearable
              />
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
  </div>
</template>

<style scoped>
.user-info-panel {
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

.path-value {
  font-family: monospace;
  background: rgba(0, 0, 0, 0.02);
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 13px;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.copy-button:hover {
  background-color: rgba(64, 158, 255, 0.1);
  border-radius: 4px;
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

.user-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-cell .el-icon {
  color: #409EFF;
  font-size: 16px;
}

:deep(.el-table__empty-block) {
  background: transparent;
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

.reset-button {
  margin-left: 16px;
}
</style> 
<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { GetStartupItems } from '../../wailsjs/go/pkg/App'
import { Timer, Document, Filter } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { pkg } from '../../wailsjs/go/models'

const startupItems = ref<pkg.StartupItem[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)

// 筛选条件
const filters = ref({
  name: '',
  type: '',
  path: '',
  enabled: ''
})

// 重置筛选条件
const resetFilters = () => {
  filters.value = {
    name: '',
    type: '',
    path: '',
    enabled: ''
  }
  currentPage.value = 1
}

// 筛选后的数据
const filteredItems = computed(() => {
  return startupItems.value.filter(item => {
    return (
      (!filters.value.name || item.name.toLowerCase().includes(filters.value.name.toLowerCase())) &&
      (!filters.value.type || item.type.toLowerCase().includes(filters.value.type.toLowerCase())) &&
      (!filters.value.path || item.path.toLowerCase().includes(filters.value.path.toLowerCase())) &&
      (!filters.value.enabled || item.enabled.toString().toLowerCase().includes(filters.value.enabled.toLowerCase()))
    )
  })
})

// 计算当前页的数据
const currentPageData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredItems.value.slice(start, end)
})

// 更新总数
const updateTotal = () => {
  total.value = filteredItems.value.length
  if (currentPage.value > Math.ceil(total.value / pageSize.value)) {
    currentPage.value = 1
  }
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
}

// 监听筛选条件变化
watch(filters, () => {
  updateTotal()
}, { deep: true })

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

// 添加 refresh 方法，用于重新获取启动项信息
const refresh = async () => {
  loading.value = true
  try {
    const items = await GetStartupItems()
    startupItems.value = items
    updateTotal()
  } catch (error) {
    console.error('获取启动项信息失败:', error)
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
  <div class="startup-panel">
    <div class="info-card">
      <div class="card-header">
        <el-icon :size="18" color="#409EFF"><Timer /></el-icon>
        <h3>开机启动项</h3>
        <span class="total-count">共 {{ total }} 个启动项</span>
      </div>

      <el-table 
        :data="currentPageData"
        style="width: 100%"
        v-loading="loading"
        element-loading-text="正在加载启动项信息..."
        element-loading-background="rgba(255, 255, 255, 0.8)"
        border
        :resizable="true"
      >
        <el-table-column 
          prop="name" 
          label="名称" 
          min-width="120"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="table-header">
              <span>名称</span>
              <el-input
                v-model="filters.name"
                placeholder="筛选名称"
                size="small"
                clearable
              />
            </div>
          </template>
        </el-table-column>

        <el-table-column 
          prop="type" 
          label="类型" 
          min-width="120"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="table-header">
              <span>类型</span>
              <el-input
                v-model="filters.type"
                placeholder="筛选类型"
                size="small"
                clearable
              />
            </div>
          </template>
        </el-table-column>

        <el-table-column 
          prop="path" 
          label="路径" 
          min-width="200"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="table-header">
              <span>路径</span>
              <el-input
                v-model="filters.path"
                placeholder="筛选路径"
                size="small"
                clearable
              />
            </div>
          </template>
        </el-table-column>

        <el-table-column 
          prop="enabled" 
          label="状态" 
          width="100"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="table-header">
              <span>状态</span>
              <el-input
                v-model="filters.enabled"
                placeholder="筛选状态"
                size="small"
                clearable
              />
            </div>
          </template>
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column 
          prop="description" 
          label="描述" 
          min-width="200"
          resizable
          show-overflow-tooltip
        />
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
.startup-panel {
  padding: 16px;
}

.info-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 16px;
}

.card-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  gap: 8px;
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  color: #303133;
}

.total-count {
  margin-left: auto;
  color: #909399;
  font-size: 14px;
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
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
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

:deep(.el-table) {
  --el-table-border-color: #EBEEF5;
  --el-table-header-bg-color: #F5F7FA;
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

:deep(.el-table .cell) {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style> 
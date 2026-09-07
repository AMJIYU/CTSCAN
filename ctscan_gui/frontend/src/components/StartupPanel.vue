<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { GetStartupItems, SaveStartupItems } from '../../wailsjs/go/pkg/App'
import { Timer, Filter } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { pkg } from '../../wailsjs/go/models'
import { formatBytes } from '../utils/logExport'
import type { LogSnapshot } from '../utils/logExport'

type StartupRow = pkg.StartupItem & {
  category?: string
  location?: string
  image_path?: string
  publisher?: string
}

const startupItems = ref<StartupRow[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)
const startupLoadTimeoutMs = 25000

// 筛选条件
const filters = ref({
  keyword: '',
  category: '',
  name: '',
  type: '',
  path: '',
  location: '',
  publisher: '',
  lastModTime: '',
  size: '',
  description: '',
  enabled: ''
})

// 重置筛选条件
const resetFilters = () => {
  filters.value = {
    keyword: '',
    category: '',
    name: '',
    type: '',
    path: '',
    location: '',
    publisher: '',
    lastModTime: '',
    size: '',
    description: '',
    enabled: ''
  }
  currentPage.value = 1
}

const textValue = (value: unknown) => String(value ?? '').toLowerCase()
const displayCategory = (item: StartupRow) => item.category || 'Other'
const displayImagePath = (item: StartupRow) => item.image_path || item.path || ''
const displayLocation = (item: StartupRow) => item.location || item.path || ''
const displayPublisher = (item: StartupRow) => item.publisher || ''
const displayStatus = (enabled: boolean) => enabled ? '启用' : '禁用'
const formatStartupTime = (value: unknown) => {
  if (!value) return '-'
  const date = new Date(String(value))
  if (Number.isNaN(date.getTime()) || date.getFullYear() <= 1970) return '-'
  return date.toLocaleString('zh-CN', { hour12: false })
}

const categoryOptions = computed(() => {
  const categories = startupItems.value.map(displayCategory).filter(Boolean)
  return Array.from(new Set(categories)).sort((a, b) => a.localeCompare(b))
})

const categoryStats = computed(() => {
  const counts = new Map<string, number>()
  startupItems.value.forEach(item => {
    const category = displayCategory(item)
    counts.set(category, (counts.get(category) || 0) + 1)
  })
  return Array.from(counts.entries())
    .map(([category, count]) => ({ category, count }))
    .sort((a, b) => b.count - a.count || a.category.localeCompare(b.category))
})

// 筛选后的数据
const filteredItems = computed(() => {
  return startupItems.value.filter(item => {
    const keyword = filters.value.keyword.toLowerCase()
    const searchable = [
      item.name,
      item.type,
      item.path,
      item.description,
      displayCategory(item),
      displayImagePath(item),
      displayLocation(item),
      displayPublisher(item),
      formatStartupTime(item.lastModTime),
      formatBytes(item.size)
    ].map(textValue).join('\n')

    return (
      (!keyword || searchable.includes(keyword)) &&
      (!filters.value.category || displayCategory(item) === filters.value.category) &&
      (!filters.value.name || textValue(item.name).includes(filters.value.name.toLowerCase())) &&
      (!filters.value.type || textValue(item.type).includes(filters.value.type.toLowerCase())) &&
      (!filters.value.path || textValue(displayImagePath(item)).includes(filters.value.path.toLowerCase())) &&
      (!filters.value.location || textValue(displayLocation(item)).includes(filters.value.location.toLowerCase())) &&
      (!filters.value.publisher || textValue(displayPublisher(item)).includes(filters.value.publisher.toLowerCase())) &&
      (!filters.value.lastModTime || textValue(formatStartupTime(item.lastModTime)).includes(filters.value.lastModTime.toLowerCase())) &&
      (!filters.value.size || textValue(`${item.size} ${formatBytes(item.size)}`).includes(filters.value.size.toLowerCase())) &&
      (!filters.value.description || textValue(item.description).includes(filters.value.description.toLowerCase())) &&
      (!filters.value.enabled || String(item.enabled) === filters.value.enabled)
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

const getErrorMessage = (error: unknown, fallback: string) => {
  if (error instanceof Error && error.message) return error.message
  if (typeof error === 'string' && error.trim()) return error
  if (error && typeof error === 'object') {
    const value = error as Record<string, unknown>
    if (typeof value.message === 'string' && value.message.trim()) return value.message
    if (typeof value.error === 'string' && value.error.trim()) return value.error
  }
  return fallback
}

const withTimeout = async <T,>(task: Promise<T>, timeoutMs: number, message: string): Promise<T> => {
  let timer: number | undefined
  try {
    return await Promise.race([
      task,
      new Promise<T>((_, reject) => {
        timer = window.setTimeout(() => reject(new Error(message)), timeoutMs)
      })
    ])
  } finally {
    if (timer !== undefined) {
      window.clearTimeout(timer)
    }
  }
}

const saveStartupItemsInBackground = async (items: StartupRow[]) => {
  try {
    await SaveStartupItems(items)
  } catch (error) {
    console.warn('保存启动项到数据库失败:', error)
  }
}

// 监听筛选条件变化
watch(filters, () => {
  updateTotal()
}, { deep: true })

// 复制命令到剪贴板
const copyCommand = async (command: string) => {
  if (!command) {
    ElMessage({
      type: 'warning',
      message: '没有可复制的内容',
      duration: 1600
    })
    return
  }
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
    const items = await withTimeout(GetStartupItems(), startupLoadTimeoutMs, '获取开机启动项超时，已停止等待')
    const safeItems = Array.isArray(items) ? items as StartupRow[] : []
    startupItems.value = safeItems
    updateTotal()
    void saveStartupItemsInBackground(safeItems)
  } catch (error) {
    console.error('获取启动项信息失败:', error)
    ElMessage.error(getErrorMessage(error, '获取启动项信息失败'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  refresh()
})

const getLogSnapshot = (): LogSnapshot => ({
  title: '开机启动项',
  description: '开机启动项与持久化位置查询结果。',
  filters: { ...filters.value },
  sections: [
    {
      title: '开机启动项',
      columns: [
        { key: 'category', label: '分类' },
        { key: 'name', label: '名称' },
        { key: 'type', label: '类型' },
        { key: 'enabled', label: '状态' },
        { key: 'image_path', label: '镜像/命令路径' },
        { key: 'location', label: '注册/触发位置' },
        { key: 'publisher', label: '发布者' },
        { key: 'last_mod_time', label: '文件修改时间' },
        { key: 'size', label: '大小' },
        { key: 'description', label: '描述' }
      ],
      rows: filteredItems.value.map(item => ({
        category: displayCategory(item),
        name: item.name,
        type: item.type,
        enabled: displayStatus(item.enabled),
        image_path: displayImagePath(item),
        location: displayLocation(item),
        publisher: displayPublisher(item),
        last_mod_time: formatStartupTime(item.lastModTime),
        size: formatBytes(item.size),
        description: item.description
      }))
    }
  ]
})

// 暴露 refresh 方法，供父组件调用
defineExpose({ refresh, getLogSnapshot })
</script>

<template>
  <div class="startup-panel">
    <div class="info-card">
      <div class="card-header">
        <el-icon :size="18" color="#409EFF"><Timer /></el-icon>
        <h3>Autoruns 风格启动项</h3>
        <span class="total-count">共 {{ total }} 个启动项</span>
        <el-button
          type="primary"
          link
          :icon="Filter"
          @click="resetFilters"
        >
          重置筛选
        </el-button>
      </div>

      <div class="startup-filter-bar">
        <el-input
          v-model="filters.keyword"
          placeholder="全局搜索：名称、路径、注册表位置、描述..."
          size="small"
          clearable
        />
        <el-select
          v-model="filters.category"
          placeholder="全部分类"
          size="small"
          clearable
          filterable
        >
          <el-option
            v-for="category in categoryOptions"
            :key="category"
            :label="category"
            :value="category"
          />
        </el-select>
        <el-select
          v-model="filters.enabled"
          placeholder="全部状态"
          size="small"
          clearable
        >
          <el-option label="启用" value="true" />
          <el-option label="禁用" value="false" />
        </el-select>
      </div>

      <div class="category-strip" v-if="categoryStats.length">
        <button
          class="category-chip"
          :class="{ active: !filters.category }"
          @click="filters.category = ''"
        >
          Everything
          <span>{{ startupItems.length }}</span>
        </button>
        <button
          v-for="item in categoryStats"
          :key="item.category"
          class="category-chip"
          :class="{ active: filters.category === item.category }"
          @click="filters.category = item.category"
        >
          {{ item.category }}
          <span>{{ item.count }}</span>
        </button>
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
          label="分类" 
          width="150"
          fixed="left"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="table-header">
              <span>分类</span>
              <el-select
                v-model="filters.category"
                placeholder="全部"
                size="small"
                clearable
                filterable
              >
                <el-option
                  v-for="category in categoryOptions"
                  :key="category"
                  :label="category"
                  :value="category"
                />
              </el-select>
            </div>
          </template>
          <template #default="{ row }">
            <el-tag size="small" effect="plain">{{ displayCategory(row) }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column 
          prop="name" 
          label="名称" 
          min-width="180"
          fixed="left"
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
          label="启动类型" 
          min-width="150"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="table-header">
              <span>启动类型</span>
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
          prop="enabled" 
          label="状态" 
          width="100"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="table-header">
              <span>状态</span>
              <el-select
                v-model="filters.enabled"
                placeholder="全部"
                size="small"
                clearable
              >
                <el-option label="启用" value="true" />
                <el-option label="禁用" value="false" />
              </el-select>
            </div>
          </template>
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ displayStatus(row.enabled) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column 
          label="镜像/命令路径" 
          min-width="320"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="table-header">
              <span>镜像/命令路径</span>
              <el-input
                v-model="filters.path"
                placeholder="筛选镜像路径"
                size="small"
                clearable
              />
            </div>
          </template>
          <template #default="{ row }">
            <div class="path-cell">
              <span class="path-text">{{ displayImagePath(row) || '-' }}</span>
              <el-button
                type="primary"
                link
                size="small"
                class="copy-button"
                @click="copyCommand(displayImagePath(row))"
              >
                复制
              </el-button>
            </div>
          </template>
        </el-table-column>

        <el-table-column 
          label="注册/触发位置" 
          min-width="360"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="table-header">
              <span>注册/触发位置</span>
              <el-input 
                v-model="filters.location"
                placeholder="筛选注册表键/任务路径"
                size="small"
                clearable
              />
            </div>
          </template>
          <template #default="{ row }">
            <div class="path-cell">
              <span class="path-text">{{ displayLocation(row) || '-' }}</span>
              <el-button
                type="primary"
                link
                size="small"
                class="copy-button"
                @click="copyCommand(displayLocation(row))"
              >
                复制
              </el-button>
            </div>
          </template>
        </el-table-column>

        <el-table-column 
          label="发布者" 
          min-width="180"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="table-header">
              <span>发布者</span>
              <el-input
                v-model="filters.publisher"
                placeholder="筛选发布者"
                size="small"
                clearable
              />
            </div>
          </template>
          <template #default="{ row }">{{ displayPublisher(row) || '-' }}</template>
        </el-table-column>

        <el-table-column 
          label="文件修改时间" 
          min-width="170"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="table-header">
              <span>文件修改时间</span>
              <el-input
                v-model="filters.lastModTime"
                placeholder="筛选时间"
                size="small"
                clearable
              />
            </div>
          </template>
          <template #default="{ row }">{{ formatStartupTime(row.lastModTime) }}</template>
        </el-table-column>

        <el-table-column 
          label="大小" 
          min-width="110"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="table-header">
              <span>大小</span>
              <el-input
                v-model="filters.size"
                placeholder="筛选大小"
                size="small"
                clearable
              />
            </div>
          </template>
          <template #default="{ row }">{{ formatBytes(row.size) }}</template>
        </el-table-column>

        <el-table-column 
          prop="description" 
          label="描述" 
          min-width="260"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="table-header">
              <span>描述</span>
              <el-input
                v-model="filters.description"
                placeholder="筛选描述"
                size="small"
                clearable
              />
            </div>
          </template>
          <template #default="{ row }">{{ row.description || '-' }}</template>
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

.startup-filter-bar {
  display: grid;
  grid-template-columns: minmax(260px, 1fr) 220px 140px;
  gap: 12px;
  margin-bottom: 14px;
  padding: 12px;
  background: #f8fbff;
  border: 1px solid #e4edf8;
  border-radius: 10px;
}

.category-strip {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding: 0 0 14px;
  margin-bottom: 8px;
}

.category-chip {
  flex: 0 0 auto;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  height: 30px;
  padding: 0 10px;
  color: #475569;
  background: #ffffff;
  border: 1px solid #dbe5f0;
  border-radius: 999px;
  cursor: pointer;
  font-size: 12px;
  line-height: 1;
}

.category-chip span {
  min-width: 20px;
  padding: 3px 6px;
  color: #64748b;
  background: #eef4fb;
  border-radius: 999px;
  font-size: 11px;
}

.category-chip.active {
  color: #ffffff;
  background: #409EFF;
  border-color: #409EFF;
}

.category-chip.active span {
  color: #1d4f8f;
  background: #ffffff;
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

@media (max-width: 1100px) {
  .startup-filter-bar {
    grid-template-columns: 1fr;
  }
}
</style> 

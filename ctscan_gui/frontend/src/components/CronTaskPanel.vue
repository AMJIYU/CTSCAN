<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { GetCronTasks } from '../../wailsjs/go/pkg/App'
import { Timer, Document } from '@element-plus/icons-vue'

interface CronTask {
  line: string
  type?: string
  schedule?: string
  command?: string
  user?: string
  lastRun?: string
  nextRun?: string
  status?: string
}

const cronTasks = ref<CronTask[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 筛选条件
const filters = ref({
  type: '',
  command: '',
  schedule: '',
  status: '',
  lastRun: '',
  nextRun: ''
})

// 重置筛选条件
const resetFilters = () => {
  filters.value = {
    type: '',
    command: '',
    schedule: '',
    status: '',
    lastRun: '',
    nextRun: ''
  }
  currentPage.value = 1
}

// 筛选后的数据
const filteredTasks = computed(() => {
  return cronTasks.value.filter(task => {
    return (
      (!filters.value.type || (task.type && task.type.toLowerCase().includes(filters.value.type.toLowerCase()))) &&
      (!filters.value.command || (task.command && task.command.toLowerCase().includes(filters.value.command.toLowerCase()))) &&
      (!filters.value.schedule || (task.schedule && task.schedule.toLowerCase().includes(filters.value.schedule.toLowerCase()))) &&
      (!filters.value.status || (task.status && task.status.toLowerCase().includes(filters.value.status.toLowerCase()))) &&
      (!filters.value.lastRun || (task.lastRun && task.lastRun.toLowerCase().includes(filters.value.lastRun.toLowerCase()))) &&
      (!filters.value.nextRun || (task.nextRun && task.nextRun.toLowerCase().includes(filters.value.nextRun.toLowerCase())))
    )
  })
})

// 计算当前页的数据
const currentPageData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredTasks.value.slice(start, end)
})

// 更新总数
const updateTotal = () => {
  total.value = filteredTasks.value.length
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

// 解析计划任务内容
const parseTask = (line: string): CronTask => {
  const task: CronTask = { line }
  
  // 解析 Windows 任务
  if (line.includes('TaskName:')) {
    const parts = line.split(';').map(p => p.trim())
    parts.forEach(part => {
      if (part.startsWith('TaskName:')) {
        task.type = 'Windows'
        task.command = part.replace('TaskName:', '').trim()
      } else if (part.startsWith('State:')) {
        task.status = part.replace('State:', '').trim()
      } else if (part.startsWith('LastRunTime:')) {
        task.lastRun = part.replace('LastRunTime:', '').trim()
      } else if (part.startsWith('NextRunTime:')) {
        task.nextRun = part.replace('NextRunTime:', '').trim()
      }
    })
  } 
  // 解析 Unix 任务
  else {
    task.type = 'Unix'
    const parts = line.split(/\s+/)
    if (parts.length >= 6) {
      task.schedule = parts.slice(0, 5).join(' ')
      task.command = parts.slice(5).join(' ')
      task.user = '当前用户'
    }
  }
  
  return task
}

// 解析crontab表达式
const parseCronExpression = (expression: string): string => {
  const parts = expression.split(/\s+/)
  if (parts.length !== 5) return expression

  const [minute, hour, dayOfMonth, month, dayOfWeek] = parts
  let description = ''

  // 解析分钟
  if (minute === '*') {
    description += '每分钟'
  } else if (minute.includes('*/')) {
    const interval = minute.split('*/')[1]
    description += `每${interval}分钟`
  } else {
    description += `在${minute}分`
  }

  // 解析小时
  if (hour === '*') {
    description += '每小时'
  } else if (hour.includes('*/')) {
    const interval = hour.split('*/')[1]
    description += `每${interval}小时`
  } else {
    description += `在${hour}时`
  }

  // 解析日期
  if (dayOfMonth === '*') {
    description += '每天'
  } else if (dayOfMonth.includes('*/')) {
    const interval = dayOfMonth.split('*/')[1]
    description += `每${interval}天`
  } else {
    description += `在${dayOfMonth}日`
  }

  // 解析月份
  if (month === '*') {
    description += '每月'
  } else if (month.includes('*/')) {
    const interval = month.split('*/')[1]
    description += `每${interval}个月`
  } else {
    const months = month.split(',').map(m => {
      const monthNames = ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月']
      return monthNames[parseInt(m) - 1] || m
    })
    description += `在${months.join('、')}`
  }

  // 解析星期
  if (dayOfWeek === '*') {
    description += '每周'
  } else if (dayOfWeek.includes('*/')) {
    const interval = dayOfWeek.split('*/')[1]
    description += `每${interval}周`
  } else {
    const weekdays = dayOfWeek.split(',').map(w => {
      const weekdayNames = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
      return weekdayNames[parseInt(w)] || w
    })
    description += `在${weekdays.join('、')}`
  }

  return description
}

const refresh = async () => {
  loading.value = true
  try {
    const tasks = await GetCronTasks()
    cronTasks.value = tasks.map(task => parseTask(task.line))
  } catch (error) {
    console.error('获取计划任务失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  refresh()
})

defineExpose({ refresh })
</script>

<template>
  <div class="cron-task-panel">
    <div class="panel-header">
      <el-icon :size="20" color="#409EFF"><Timer /></el-icon>
      <h2>计划任务</h2>
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
      :data="currentPageData" 
      style="width: 100%" 
      border
      :resizable="true"
      size="small"
      v-loading="loading"
      element-loading-text="正在加载计划任务信息..."
      element-loading-background="rgba(255, 255, 255, 0.8)"
    >
      <el-table-column 
        prop="type" 
        label="类型" 
        width="100"
        align="center"
        resizable
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
        <template #default="{ row }">
          <el-tag 
            size="small" 
            :type="row.type === 'Windows' ? 'primary' : 'success'"
          >
            {{ row.type }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column 
        prop="command" 
        label="任务内容" 
        min-width="200"
        resizable
        show-overflow-tooltip
      >
        <template #header>
          <div class="table-header">
            <span>任务内容</span>
            <el-input
              v-model="filters.command"
              placeholder="筛选任务内容"
              size="small"
              clearable
            />
          </div>
        </template>
        <template #default="{ row }">
          <div class="task-content">
            <el-icon><Document /></el-icon>
            <span>{{ row.command }}</span>
          </div>
        </template>
      </el-table-column>

      <el-table-column 
        prop="schedule" 
        label="执行计划" 
        min-width="150"
        resizable
        show-overflow-tooltip
      >
        <template #header>
          <div class="table-header">
            <span>执行计划</span>
            <el-input
              v-model="filters.schedule"
              placeholder="筛选执行计划"
              size="small"
              clearable
            />
          </div>
        </template>
        <template #default="{ row }">
          <div v-if="row.schedule" class="schedule-value">
            <div class="schedule-expression">{{ row.schedule }}</div>
            <div class="schedule-description">{{ parseCronExpression(row.schedule) }}</div>
          </div>
          <span v-else class="no-schedule">未设置</span>
        </template>
      </el-table-column>

      <el-table-column 
        prop="status" 
        label="状态" 
        width="100"
        align="center"
        resizable
      >
        <template #header>
          <div class="table-header">
            <span>状态</span>
            <el-input
              v-model="filters.status"
              placeholder="筛选状态"
              size="small"
              clearable
            />
          </div>
        </template>
        <template #default="{ row }">
          <el-tag 
            v-if="row.status"
            size="small" 
            :type="row.status === 'Running' ? 'success' : 'info'"
          >
            {{ row.status }}
          </el-tag>
          <span v-else class="no-status">-</span>
        </template>
      </el-table-column>

      <el-table-column 
        prop="lastRun" 
        label="上次运行" 
        min-width="150"
        resizable
        show-overflow-tooltip
      >
        <template #header>
          <div class="table-header">
            <span>上次运行</span>
            <el-input
              v-model="filters.lastRun"
              placeholder="筛选上次运行"
              size="small"
              clearable
            />
          </div>
        </template>
        <template #default="{ row }">
          <span v-if="row.lastRun" class="time-value">{{ row.lastRun }}</span>
          <span v-else class="no-time">-</span>
        </template>
      </el-table-column>

      <el-table-column 
        prop="nextRun" 
        label="下次运行" 
        min-width="150"
        resizable
        show-overflow-tooltip
      >
        <template #header>
          <div class="table-header">
            <span>下次运行</span>
            <el-input
              v-model="filters.nextRun"
              placeholder="筛选下次运行"
              size="small"
              clearable
            />
          </div>
        </template>
        <template #default="{ row }">
          <span v-if="row.nextRun" class="time-value">{{ row.nextRun }}</span>
          <span v-else class="no-time">-</span>
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
.cron-task-panel {
  padding: 0;
}

.panel-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 20px;
  padding: 0 16px;
}

.panel-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: #1a202c;
  margin: 0;
}

.reset-button {
  margin-left: auto;
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

.task-content {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.task-content .el-icon {
  color: #409EFF;
  font-size: 16px;
  flex-shrink: 0;
}

.task-content span {
  font-family: monospace;
  font-size: 13px;
  color: #2d3748;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.schedule-value {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.schedule-expression {
  font-family: monospace;
  font-size: 13px;
  color: #409EFF;
  background: rgba(64, 158, 255, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
}

.schedule-description {
  font-size: 12px;
  color: #606266;
  line-height: 1.4;
}

.no-schedule {
  color: #909399;
  font-style: italic;
}

.time-value {
  font-family: monospace;
  font-size: 13px;
  color: #2d3748;
}

.no-time {
  color: #909399;
  font-style: italic;
}

.no-status {
  color: #909399;
  font-style: italic;
}

:deep(.el-tag) {
  border-radius: 4px;
  font-size: 12px;
  height: 20px;
  line-height: 18px;
  padding: 0 6px;
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

:deep(.el-loading-spinner .el-loading-text) {
  color: #409EFF;
  font-size: 14px;
  margin-top: 8px;
}

:deep(.el-loading-spinner .circular) {
  width: 30px;
  height: 30px;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
  padding: 0 16px;
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
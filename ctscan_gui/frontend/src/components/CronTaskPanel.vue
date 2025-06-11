<script setup lang="ts">
import { ref, onMounted } from 'vue'
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

const refresh = async () => {
  try {
    const tasks = await GetCronTasks()
    cronTasks.value = tasks.map(task => parseTask(task.line))
  } catch (error) {
    console.error('获取计划任务失败:', error)
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
    </div>
    
    <el-table 
      :data="cronTasks" 
      style="width: 100%" 
      border
      :resizable="true"
      size="small"
    >
      <el-table-column 
        prop="type" 
        label="类型" 
        width="100"
        align="center"
        resizable
      >
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
        <template #default="{ row }">
          <span v-if="row.schedule" class="schedule-value">{{ row.schedule }}</span>
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
        <template #default="{ row }">
          <span v-if="row.nextRun" class="time-value">{{ row.nextRun }}</span>
          <span v-else class="no-time">-</span>
        </template>
      </el-table-column>
    </el-table>
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
}

.panel-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: #1a202c;
  margin: 0;
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
  font-family: monospace;
  font-size: 13px;
  color: #409EFF;
  background: rgba(64, 158, 255, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
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
</style> 
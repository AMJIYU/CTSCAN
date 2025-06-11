<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetCronTasks } from '../../wailsjs/go/pkg/App'

const cronTasks = ref<{ line: string }[]>([])

const refresh = async () => {
  try {
    const tasks = await GetCronTasks()
    cronTasks.value = tasks
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
    <h2>计划任务</h2>
    <el-table :data="cronTasks" style="width: 100%">
      <el-table-column prop="line" label="计划任务内容">
        <template #default="{ row }">
          <pre class="task-content">{{ row.line }}</pre>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<style scoped>
.cron-task-panel {
  padding: 0;
}

h2 {
  font-size: 20px;
  font-weight: 600;
  color: #1a202c;
  margin-bottom: 24px;
}

:deep(.el-table) {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
}

:deep(.el-table th) {
  background-color: rgba(0, 0, 0, 0.02);
  font-weight: 600;
  color: #1a202c;
}

:deep(.el-table td) {
  color: #4a5568;
  padding: 0;
}

:deep(.el-table tr:hover > td) {
  background-color: rgba(64, 158, 255, 0.05);
}

.task-content {
  white-space: pre-wrap;
  word-break: break-all;
  font-family: monospace;
  margin: 0;
  padding: 16px;
  background: rgba(0, 0, 0, 0.02);
  border-radius: 4px;
  font-size: 14px;
  line-height: 1.5;
}
</style> 
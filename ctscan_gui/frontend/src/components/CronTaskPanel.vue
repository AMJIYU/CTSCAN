<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetCronTasks } from '../../wailsjs/go/pkg/App'

const cronTasks = ref<{ line: string }[]>([])

const refresh = () => {
  GetCronTasks().then(list => {
    cronTasks.value = list
  })
}

onMounted(() => {
  refresh()
})

defineExpose({ refresh })
</script>
<template>
  <div class="cron-task-panel">
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
  padding: 16px;
}

.task-content {
  white-space: pre-wrap;
  word-break: break-all;
  font-family: monospace;
  margin: 0;
  padding: 8px;
  background-color: #f5f7fa;
  border-radius: 4px;
}
</style> 
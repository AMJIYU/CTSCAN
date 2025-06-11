<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetStartupItems } from '../../wailsjs/go/pkg/App'

const startupItems = ref<{ name: string; path: string }[]>([])

const refresh = async () => {
  try {
    const items = await GetStartupItems()
    startupItems.value = items
  } catch (error) {
    console.error('获取启动项失败:', error)
  }
}

onMounted(() => {
  refresh()
})

defineExpose({ refresh })
</script>

<template>
  <div class="startup-panel">
    <h2>开机启动项</h2>
    <el-table :data="startupItems" style="width: 100%">
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="path" label="路径">
        <template #default="scope">
          <div class="path-cell">{{ scope.row.path }}</div>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<style scoped>
.startup-panel {
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
}

:deep(.el-table tr:hover > td) {
  background-color: rgba(64, 158, 255, 0.05);
}

.path-cell {
  font-family: monospace;
  padding: 4px 8px;
  background: rgba(0, 0, 0, 0.02);
  border-radius: 4px;
  word-break: break-all;
}
</style> 
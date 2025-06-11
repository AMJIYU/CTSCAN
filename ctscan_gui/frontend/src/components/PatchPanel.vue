<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetPatchRecords } from '../../wailsjs/go/pkg/App'

const patchInfo = ref<{ time: string; title: string; description: string; status: string; kb: string }[]>([])

const refresh = () => {
  GetPatchRecords().then(info => {
    patchInfo.value = info
  })
}

onMounted(() => {
  refresh()
})

defineExpose({ refresh })
</script>
<template>
  <div class="patch-panel">
    <h2>系统更新记录</h2>
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>时间</th>
            <th>标题</th>
            <th>描述</th>
            <th>状态</th>
            <th>KB编号</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="record in patchInfo" :key="record.kb">
            <td>{{ record.time }}</td>
            <td>{{ record.title }}</td>
            <td>{{ record.description }}</td>
            <td>{{ record.status }}</td>
            <td>{{ record.kb }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.patch-panel {
  padding: 20px;
}

.table-container {
  margin-top: 20px;
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  background-color: #fff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

th, td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

th {
  background-color: #f5f5f5;
  font-weight: 600;
}

tr:hover {
  background-color: #f9f9f9;
}
</style> 
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetLoginSuccessRecords } from '../../wailsjs/go/pkg/App'

const records = ref<any[]>([])
const loginSuccessLogs = ref<any[]>([])

// 添加 refresh 方法，用于重新获取登录成功日志信息
const refresh = () => {
  GetLoginSuccessRecords().then(logs => {
    loginSuccessLogs.value = logs
  })
}

onMounted(() => {
  GetLoginSuccessRecords().then(list => {
    records.value = list
  })
  refresh()
})

// 暴露 refresh 方法，供父组件调用
defineExpose({ refresh })
</script>
<template>
  <div class="login-success-panel">
    <h2>登录成功记录</h2>
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>时间</th>
            <th>事件ID</th>
            <th>事件类型</th>
            <th>来源</th>
            <th>用户名</th>
            <th>IP地址</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="record in records" :key="record.time">
            <td>{{ record.time }}</td>
            <td>{{ record.event_id }}</td>
            <td>{{ record.event_type }}</td>
            <td>{{ record.source }}</td>
            <td>{{ record.username }}</td>
            <td>{{ record.ip_address }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.login-success-panel {
  padding: 0;
}

.table-container {
  margin-top: 24px;
  overflow-x: auto;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
}

table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  background-color: transparent;
}

th, td {
  padding: 16px;
  text-align: left;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  font-size: 14px;
}

th {
  background-color: rgba(0, 0, 0, 0.02);
  font-weight: 600;
  color: #1a202c;
  white-space: nowrap;
}

td {
  color: #4a5568;
}

tr:hover {
  background-color: rgba(64, 158, 255, 0.05);
}

tr:last-child td {
  border-bottom: none;
}
</style> 
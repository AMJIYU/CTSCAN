<template>
  <div class="login-failed-panel">
    <h2>登录失败记录</h2>
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

<script>
export default {
  name: 'LoginFailedPanel',
  data() {
    return {
      records: []
    }
  },
  async mounted() {
    await this.fetchRecords()
  },
  methods: {
    async fetchRecords() {
      try {
        const response = await window.go.main.App.GetLoginFailedRecords()
        this.records = response
      } catch (error) {
        console.error('获取登录失败记录失败:', error)
      }
    }
  }
}
</script>

<style scoped>
.login-failed-panel {
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

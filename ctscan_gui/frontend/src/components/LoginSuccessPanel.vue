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
  <el-table :data="records" style="width: 100%">
    <el-table-column prop="time" label="时间" />
    <el-table-column prop="event_id" label="事件ID" />
    <el-table-column prop="login_type" label="登入类型" />
    <el-table-column prop="source_ip" label="源IP" />
    <el-table-column prop="username" label="用户名" />
    <el-table-column prop="workstation" label="工作站" />
    <el-table-column prop="subject_user" label="主体用户名" />
    <el-table-column prop="subject_domain" label="主体域" />
    <el-table-column prop="process" label="进程" />
  </el-table>
</template> 
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetAllProcesses } from '../../wailsjs/go/main/App'

const processes = ref<{ pid: number; name: string; cmdline: string; user: string; exe: string }[]>([])

onMounted(() => {
  GetAllProcesses().then(list => {
    processes.value = list
  })
})
</script>
<template>
  <el-table :data="processes" style="width: 100%">
    <el-table-column prop="pid" label="PID" width="80" />
    <el-table-column prop="name" label="进程名" />
    <el-table-column prop="user" label="用户" />
    <el-table-column prop="cmdline" label="命令行" />
    <el-table-column prop="exe" label="启动路径" />
  </el-table>
</template> 
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetAllProcesses } from '../../wailsjs/go/pkg/App'

const processes = ref<any[]>([])

const processList = ref([])

const refresh = () => {
  GetAllProcesses().then(list => {
    processes.value = list
  })
}

onMounted(() => {
  refresh()
})

defineExpose({ refresh })
</script>
<template>
  <el-table :data="processes" style="width: 100%" :default-sort="{ prop: 'cpu_percent', order: 'descending' }">
    <el-table-column prop="pid" label="PID" width="80" />
    <el-table-column prop="name" label="进程名" />
    <el-table-column prop="ppid" label="父PID" />
    <el-table-column prop="parent_name" label="父进程名" />
    <el-table-column prop="cpu_percent" label="CPU使用率" sortable
      :formatter="row => row.cpu_percent ? row.cpu_percent.toFixed(2) + '%' : '0%'" />
    <el-table-column prop="mem_percent" label="内存使用率" sortable
      :formatter="row => row.mem_percent ? row.mem_percent.toFixed(2) + '%' : '0%'" />
    <el-table-column prop="create_time" label="创建时间" sortable
      :formatter="row => row.create_time ? new Date(row.create_time * 1000).toLocaleString() : ''" />
    <el-table-column prop="exe" label="可执行文件路径" />
    <el-table-column prop="file_ctime" label="文件创建时间" sortable
      :formatter="row => row.file_ctime ? new Date(row.file_ctime * 1000).toLocaleString() : ''" />
    <el-table-column prop="file_mtime" label="文件修改时间" sortable
      :formatter="row => row.file_mtime ? new Date(row.file_mtime * 1000).toLocaleString() : ''" />
    <el-table-column prop="md5" label="MD5" />
    <el-table-column prop="signature" label="签名信息" />
  </el-table>
</template> 
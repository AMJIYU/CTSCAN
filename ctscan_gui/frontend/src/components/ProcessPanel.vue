<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetAllProcesses } from '../../wailsjs/go/main/App'

const processes = ref<{
  pid: number;
  name: string;
  ppid: number;
  parent_name: string;
  create_time: number;
  exe: string;
  file_ctime: number;
  file_mtime: number;
  md5: string;
  signature: string;
}[]>([])

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
    <el-table-column prop="ppid" label="父PID" />
    <el-table-column prop="parent_name" label="父进程名" />
    <el-table-column prop="create_time" label="创建时间"
      :formatter="row => row.create_time ? new Date(row.create_time * 1000).toLocaleString() : ''" />
    <el-table-column prop="exe" label="可执行文件路径" />
    <el-table-column prop="file_ctime" label="文件创建时间"
      :formatter="row => row.file_ctime ? new Date(row.file_ctime * 1000).toLocaleString() : ''" />
    <el-table-column prop="file_mtime" label="文件修改时间"
      :formatter="row => row.file_mtime ? new Date(row.file_mtime * 1000).toLocaleString() : ''" />
    <el-table-column prop="md5" label="MD5" />
    <el-table-column prop="signature" label="签名信息" />
  </el-table>
</template> 
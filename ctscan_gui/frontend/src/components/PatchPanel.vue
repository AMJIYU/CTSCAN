<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetPatchInfo } from '../../wailsjs/go/main/App'

const patches = ref<{ name: string; date: string; status: string }[]>([])

onMounted(() => {
  GetPatchInfo().then(list => {
    patches.value = list
  })
})
</script>
<template>
  <el-table :data="patches" style="width: 100%">
    <el-table-column prop="name" label="补丁名称" />
    <el-table-column prop="date" label="安装日期" />
    <el-table-column prop="status" label="状态" />
  </el-table>
</template> 
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetPatchInfo } from '../../wailsjs/go/pkg/App'

const patchInfo = ref<{ name: string; date: string; status: string }[]>([])

const refresh = () => {
  GetPatchInfo().then(info => {
    patchInfo.value = info
  })
}

onMounted(() => {
  refresh()
})

defineExpose({ refresh })
</script>
<template>
  <el-table :data="patchInfo" style="width: 100%">
    <el-table-column prop="name" label="补丁名称" />
    <el-table-column prop="date" label="安装日期" />
    <el-table-column prop="status" label="状态" />
  </el-table>
</template> 
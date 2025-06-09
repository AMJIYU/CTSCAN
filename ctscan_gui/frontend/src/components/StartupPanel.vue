<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetStartupItems } from '../../wailsjs/go/main/App'

const startupItems = ref<{ name: string; path: string }[]>([])

onMounted(() => {
  GetStartupItems().then(list => {
    startupItems.value = list
  })
})
</script>
<template>
  <el-table :data="startupItems" style="width: 100%">
    <el-table-column prop="name" label="名称" />
    <el-table-column prop="path" label="路径" />
  </el-table>
</template> 
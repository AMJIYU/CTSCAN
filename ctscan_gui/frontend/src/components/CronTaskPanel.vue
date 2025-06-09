<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetCronTasks } from '../../wailsjs/go/main/App'

const cronTasks = ref<{ line: string }[]>([])

onMounted(() => {
  GetCronTasks().then(list => {
    cronTasks.value = list
  })
})
</script>
<template>
  <el-table :data="cronTasks" style="width: 100%">
    <el-table-column prop="line" label="计划任务内容" />
  </el-table>
</template> 
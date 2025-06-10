<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetCronTasks } from '../../wailsjs/go/pkg/App'

const cronTasks = ref<{ line: string }[]>([])

const refresh = () => {
  GetCronTasks().then(list => {
    cronTasks.value = list
  })
}

onMounted(() => {
  refresh()
})

defineExpose({ refresh })
</script>
<template>
  <el-table :data="cronTasks" style="width: 100%">
    <el-table-column prop="line" label="计划任务内容" />
  </el-table>
</template> 
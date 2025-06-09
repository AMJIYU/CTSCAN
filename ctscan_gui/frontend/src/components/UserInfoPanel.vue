<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetUserInfo } from '../../wailsjs/go/main/App'

const userInfo = ref({
  username: '',
  uid: '',
  gid: '',
  home_dir: '',
  name: ''
})

onMounted(() => {
  GetUserInfo().then(info => {
    userInfo.value = info
  })
})
</script>
<template>
  <el-descriptions title="当前用户信息" :column="2" border>
    <el-descriptions-item label="用户名">{{ userInfo.username }}</el-descriptions-item>
    <el-descriptions-item label="UID">{{ userInfo.uid }}</el-descriptions-item>
    <el-descriptions-item label="GID">{{ userInfo.gid }}</el-descriptions-item>
    <el-descriptions-item label="家目录">{{ userInfo.home_dir }}</el-descriptions-item>
    <el-descriptions-item label="显示名">{{ userInfo.name }}</el-descriptions-item>
  </el-descriptions>
</template> 
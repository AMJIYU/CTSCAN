<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetUserInfo, GetAllUsers } from '../../wailsjs/go/pkg/App'

const userInfo = ref({
  username: '',
  uid: '',
  gid: '',
  home_dir: '',
  name: ''
})

const allUsers = ref<{ username: string; uid: string; gid: string; home_dir: string; name: string }[]>([])

// 添加 refresh 方法，用于重新获取用户信息
const refresh = async () => {
  try {
    const [info, users] = await Promise.all([
      GetUserInfo(),
      GetAllUsers()
    ])
    userInfo.value = info
    allUsers.value = users
  } catch (error) {
    console.error('获取用户信息失败:', error)
  }
}

onMounted(() => {
  refresh()
})

// 暴露 refresh 方法，供父组件调用
defineExpose({ refresh })
</script>

<template>
  <div class="user-info-panel">
    <el-descriptions title="当前用户信息" :column="2" border>
      <el-descriptions-item label="用户名">{{ userInfo.username }}</el-descriptions-item>
      <el-descriptions-item label="UID">{{ userInfo.uid }}</el-descriptions-item>
      <el-descriptions-item label="GID">{{ userInfo.gid }}</el-descriptions-item>
      <el-descriptions-item label="家目录">{{ userInfo.home_dir }}</el-descriptions-item>
      <el-descriptions-item label="显示名">{{ userInfo.name }}</el-descriptions-item>
    </el-descriptions>
    
    <h4 style="margin-top:24px;">所有系统用户</h4>
    <el-table :data="allUsers" style="width: 100%">
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="uid" label="UID" />
      <el-table-column prop="gid" label="GID" />
      <el-table-column prop="home_dir" label="家目录" />
      <el-table-column prop="name" label="显示名" />
    </el-table>
  </div>
</template>

<style scoped>
.user-info-panel {
  padding: 0;
}

:deep(.el-descriptions) {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
}

:deep(.el-descriptions__title) {
  margin-bottom: 24px;
  font-size: 20px;
  font-weight: 600;
  color: #1a202c;
}

:deep(.el-descriptions__label) {
  font-weight: 500;
  color: #4a5568;
}

:deep(.el-descriptions__content) {
  color: #2d3748;
}

:deep(.el-table) {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
}

:deep(.el-table th) {
  background-color: rgba(0, 0, 0, 0.02);
  font-weight: 600;
  color: #1a202c;
}

:deep(.el-table td) {
  color: #4a5568;
}

:deep(.el-table tr:hover > td) {
  background-color: rgba(64, 158, 255, 0.05);
}
</style> 
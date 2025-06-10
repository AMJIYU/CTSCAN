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
const refresh = () => {
  GetUserInfo().then(info => {
    userInfo.value = info
  })
}

onMounted(() => {
  refresh()
  GetAllUsers().then(list => {
    allUsers.value = list
  })
})

// 暴露 refresh 方法，供父组件调用
defineExpose({ refresh })
</script>
<template>
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
</template> 
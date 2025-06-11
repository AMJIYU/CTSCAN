<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetUserInfo, GetAllUsers } from '../../wailsjs/go/pkg/App'
import { User, UserFilled, CopyDocument } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

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

// 复制路径到剪贴板
const copyPath = async (path: string) => {
  try {
    await navigator.clipboard.writeText(path)
    ElMessage({
      type: 'success',
      message: '路径已复制到剪贴板',
      duration: 2000
    })
  } catch (err) {
    ElMessage({
      type: 'error',
      message: '复制失败',
      duration: 2000
    })
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
    <!-- 当前用户信息 -->
    <div class="info-section">
      <div class="section-header">
        <el-icon :size="20" color="#409EFF"><User /></el-icon>
        <h3>当前用户信息</h3>
      </div>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="用户名">{{ userInfo.username }}</el-descriptions-item>
        <el-descriptions-item label="UID">{{ userInfo.uid }}</el-descriptions-item>
        <el-descriptions-item label="GID">{{ userInfo.gid }}</el-descriptions-item>
        <el-descriptions-item label="家目录">
          <span class="path-value">{{ userInfo.home_dir }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="显示名">{{ userInfo.name }}</el-descriptions-item>
      </el-descriptions>
    </div>

    <!-- 系统用户列表 -->
    <div class="info-section">
      <div class="section-header">
        <el-icon :size="20" color="#409EFF"><UserFilled /></el-icon>
        <h3>系统用户列表</h3>
      </div>
      <el-table 
        :data="allUsers" 
        style="width: 100%"
        border
        :resizable="true"
        size="small"
      >
        <el-table-column 
          prop="username" 
          label="用户名" 
          min-width="120"
          resizable
        >
          <template #default="{ row }">
            <div class="user-cell">
              <el-icon><User /></el-icon>
              <span>{{ row.username }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column 
          prop="uid" 
          label="UID" 
          min-width="80"
          align="center"
          resizable
        />
        <el-table-column 
          prop="gid" 
          label="GID" 
          min-width="80"
          align="center"
          resizable
        />
        <el-table-column 
          prop="home_dir" 
          label="家目录" 
          min-width="200"
          resizable
          show-overflow-tooltip
        >
          <template #default="{ row }">
            <div class="path-cell">
              <span class="path-value">{{ row.home_dir }}</span>
              <el-button
                class="copy-button"
                type="primary"
                link
                @click="copyPath(row.home_dir)"
              >
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column 
          prop="name" 
          label="显示名" 
          min-width="150"
          resizable
          show-overflow-tooltip
        />
      </el-table>
    </div>
  </div>
</template>

<style scoped>
.user-info-panel {
  padding: 0;
}

.info-section {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(0, 0, 0, 0.05);
  margin-bottom: 24px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 20px;
}

.section-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1a202c;
}

:deep(.el-descriptions) {
  background: transparent;
  padding: 0;
  box-shadow: none;
  border: none;
}

:deep(.el-descriptions__label) {
  font-weight: 500;
  color: #4a5568;
  width: 100px;
}

:deep(.el-descriptions__content) {
  color: #2d3748;
}

.path-cell {
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 0;
}

.path-value {
  font-family: monospace;
  background: rgba(0, 0, 0, 0.02);
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 13px;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.copy-button {
  padding: 2px;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.path-cell:hover .copy-button {
  opacity: 1;
}

.copy-button:hover {
  background-color: rgba(64, 158, 255, 0.1);
  border-radius: 4px;
}

:deep(.el-table__column-resize-proxy) {
  background-color: #409EFF;
}

:deep(.el-table__column-resize-handle) {
  background-color: #409EFF;
}

:deep(.el-table__column-resize-handle:hover) {
  background-color: #66b1ff;
}

:deep(.el-table) {
  background: transparent;
  border-radius: 8px;
  overflow: hidden;
}

:deep(.el-table th) {
  background-color: rgba(0, 0, 0, 0.02);
  font-weight: 600;
  color: #1a202c;
  padding: 12px 16px;
}

:deep(.el-table td) {
  color: #2d3748;
  padding: 12px 16px;
}

:deep(.el-table tr:hover > td) {
  background-color: rgba(64, 158, 255, 0.05);
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-cell .el-icon {
  color: #409EFF;
  font-size: 16px;
}

:deep(.el-table__empty-block) {
  background: transparent;
}
</style> 
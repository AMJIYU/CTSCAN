<script setup lang="ts">
import { ref, onMounted } from 'vue';
import SystemInfoPanel from './SystemInfoPanel.vue'
import UserInfoPanel from './UserInfoPanel.vue'
import NetworkInfoPanel from './NetworkInfoPanel.vue'
import StartupPanel from './StartupPanel.vue'
import CronTaskPanel from './CronTaskPanel.vue'
import ProcessPanel from './ProcessPanel.vue'
import LoginSuccessPanel from './LoginSuccessPanel.vue'
import LoginFailedPanel from './LoginFailedPanel.vue'
import ShellHistoryPanel from './ShellHistoryPanel.vue'
import FileMonitorPanel from './FileMonitorPanel.vue'
import RdploginPanel from './RdploginPanel.vue'
import {
  Monitor,
  User,
  Connection,
  Timer,
  Calendar,
  Operation,
  Key,
  Warning,
  Monitor as RdpIcon,
  Cpu,
  UploadFilled,
  Document
} from '@element-plus/icons-vue'
import { ElMessage, ElLoading } from 'element-plus'

// 使用ref引用每个选项卡组件
const systemInfoRef = ref<InstanceType<typeof SystemInfoPanel> | null>(null);
const userInfoRef = ref<InstanceType<typeof UserInfoPanel> | null>(null);
const networkInfoRef = ref<InstanceType<typeof NetworkInfoPanel> | null>(null);
const startupRef = ref<InstanceType<typeof StartupPanel> | null>(null);
const cronTaskRef = ref<InstanceType<typeof CronTaskPanel> | null>(null);
const processRef = ref();
const loginSuccessRef = ref();
const loginFailedRef = ref();
const shellHistoryRef = ref();
const fileMonitorRef = ref<InstanceType<typeof FileMonitorPanel> | null>(null);
const rdploginRef = ref<InstanceType<typeof RdploginPanel> | null>(null);

// 当前激活的面板
const activePanel = ref('system');

// 面板配置
const panels = [
  { id: 'system', name: '系统基本信息', icon: Monitor, component: SystemInfoPanel },
  { id: 'user', name: '用户信息', icon: User, component: UserInfoPanel },
  { id: 'network', name: '网络信息', icon: Connection, component: NetworkInfoPanel },
  { id: 'startup', name: '开机启动项', icon: Timer, component: StartupPanel },
  { id: 'cron', name: '任务计划', icon: Calendar, component: CronTaskPanel },
  { id: 'process', name: '进程排查', icon: Operation, component: ProcessPanel },
  { id: 'login-success', name: '登入成功', icon: Key, component: LoginSuccessPanel },
  { id: 'login-failed', name: '登入失败', icon: Warning, component: LoginFailedPanel },
  { id: 'shell-history', name: '命令记录', icon: Operation, component: ShellHistoryPanel },
  { id: 'rdp', name: 'RDP登入', icon: RdpIcon, component: RdploginPanel },
  { id: 'file-monitor', name: '文件监控', icon: Document, component: FileMonitorPanel }
];

// 添加重新获取信息的方法
const refreshInfo = async () => {
  const loading = ElLoading.service({
    lock: true,
    text: '正在重新获取所有信息...',
    background: 'rgba(255, 255, 255, 0.7)',
  })
  
  try {
    // 调用所有选项卡组件的刷新方法，重新获取当前所有选项卡的信息
    await Promise.all([
      systemInfoRef.value?.refresh(),
      userInfoRef.value?.refresh(),
      networkInfoRef.value?.refresh(),
      startupRef.value?.refresh(),
      cronTaskRef.value?.refresh(),
      processRef.value?.refresh(),
      loginSuccessRef.value?.refresh(),
      loginFailedRef.value?.refresh(),
      shellHistoryRef.value?.refresh(),
      rdploginRef.value?.refresh(),
      fileMonitorRef.value?.refresh()
    ])
    
    ElMessage({
      type: 'success',
      message: '所有信息已更新完成',
      duration: 2000
    })
  } catch (error) {
    console.error('刷新信息失败:', error)
    ElMessage({
      type: 'error',
      message: '刷新信息时出现错误',
      duration: 2000
    })
  } finally {
    loading.close()
  }
}

// 切换面板的方法
const handlePanelChange = (panel: string) => {
  activePanel.value = panel
  // 根据当前面板刷新数据
  switch (panel) {
    case 'system':
      systemInfoRef.value?.refresh()
      break
    case 'user':
      userInfoRef.value?.refresh()
      break
    case 'network':
      networkInfoRef.value?.refresh()
      break
    case 'startup':
      startupRef.value?.refresh()
      break
    case 'cron':
      cronTaskRef.value?.refresh()
      break
    case 'rdp':
      rdploginRef.value?.refresh()
      break
    case 'file-monitor':
      fileMonitorRef.value?.refresh()
      break
  }
}

const activeTab = ref('system')

// 刷新当前标签页的数据
const refreshCurrentTab = () => {
  if (activeTab.value === 'system' && systemInfoRef.value) {
    systemInfoRef.value.refresh()
  } else if (activeTab.value === 'network' && networkInfoRef.value) {
    networkInfoRef.value.refresh()
  }
}

onMounted(() => {
  refreshCurrentTab()
})
</script>

<template>
  <div class="overview-container">
    <!-- 顶部操作卡片 -->
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card shadow="never" class="action-card">
          <div class="card-content">
            <el-icon :size="48" color="#409EFF"><Cpu /></el-icon>
            <div class="text-content">
              <h3>本机扫描分析</h3>
              <p>准备扫描</p>
            </div>
            <el-button @click="refreshInfo">点击重新开始</el-button>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never" class="action-card">
           <div class="card-content">
            <el-icon :size="48" color="#67C23A"><UploadFilled /></el-icon>
            <div class="text-content">
              <h3>上传文件分析</h3>
              <p>选择文件进行分析</p>
            </div>
            <el-button>点击上传文件</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 详情区域 -->
    <div class="details-section">
      <!-- 左侧导航栏 -->
      <div class="sidebar">
        <div
          v-for="panel in panels"
          :key="panel.id"
          class="nav-item"
          :class="{ active: activePanel === panel.id }"
          @click="handlePanelChange(panel.id)"
        >
          <el-icon><component :is="panel.icon" /></el-icon>
          <span>{{ panel.name }}</span>
        </div>
      </div>

      <!-- 右侧内容区域 -->
      <div class="content-area">
        <SystemInfoPanel v-if="activePanel === 'system'" ref="systemInfoRef" />
        <UserInfoPanel v-if="activePanel === 'user'" ref="userInfoRef" />
        <NetworkInfoPanel v-if="activePanel === 'network'" ref="networkInfoRef" />
        <StartupPanel v-if="activePanel === 'startup'" ref="startupRef" />
        <CronTaskPanel v-if="activePanel === 'cron'" ref="cronTaskRef" />
        <ProcessPanel v-if="activePanel === 'process'" ref="processRef" />
        <LoginSuccessPanel v-if="activePanel === 'login-success'" ref="loginSuccessRef" />
        <LoginFailedPanel v-if="activePanel === 'login-failed'" ref="loginFailedRef" />
        <ShellHistoryPanel v-if="activePanel === 'shell-history'" ref="shellHistoryRef" />
        <RdploginPanel v-if="activePanel === 'rdp'" ref="rdploginRef" />
        <FileMonitorPanel v-if="activePanel === 'file-monitor'" ref="fileMonitorRef" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.overview-container {
  display: flex;
  flex-direction: column;
  gap: 32px;
  background: transparent;
  padding: 0;
}

.stats-section {
  flex: 0 0 120px;
  margin-bottom: 0;
}

.details-section {
  flex: 1 1 auto;
  display: flex;
  gap: 24px;
}

.sidebar {
  width: 150px;
  flex-shrink: 0;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(0, 0, 0, 0.05);
  height: fit-content;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  margin-bottom: 8px;
  border-radius: 8px;
  color: #4a5568;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.nav-item:hover {
  background: rgba(64, 158, 255, 0.05);
  color: #409EFF;
}

.nav-item.active {
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  color: #fff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
}

.nav-item .el-icon {
  margin-right: 12px;
  font-size: 18px;
}

.content-area {
  flex: 1;

  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(0, 0, 0, 0.05);
  min-height: 600px;

}

.action-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(0, 0, 0, 0.05);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
}

.action-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(64, 158, 255, 0.12);
  border-color: rgba(64, 158, 255, 0.2);
}

.card-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 32px;
}

.card-content .text-content {
  flex-grow: 1;
  margin-left: 32px;
}

.card-content h3 {
  margin: 0 0 12px 0;
  font-size: 28px;
  font-weight: 700;
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: -0.5px;
  position: relative;
  display: inline-block;
}

.card-content h3::after {
  content: '';
  position: absolute;
  bottom: -4px;
  left: 0;
  width: 40px;
  height: 3px;
  background: linear-gradient(90deg, #409EFF 0%, #66b1ff 100%);
  border-radius: 2px;
}

.card-content p {
  margin: 0;
  font-size: 16px;
  color: #4a5568;
  font-weight: 500;
  margin-top: 16px;
  opacity: 0.85;
}

/* 主选项卡和子选项卡统一美化 */
.detail-tabs,
.log-tabs {
  border: none;
  background: transparent;
  border-radius: 24px;
  padding: 16px 0 0 0;
  overflow: hidden;
}

:deep(.detail-tabs .el-tabs__header),
:deep(.log-tabs .el-tabs__header) {
  background: transparent;
  border-bottom: none;
  margin-bottom: 24px;
  padding-left: 16px;
}

:deep(.detail-tabs .el-tabs__item),
:deep(.log-tabs .el-tabs__item) {
  border-radius: 16px 16px 0 0;
  margin-right: 16px;
  padding: 16px 36px;
  font-size: 16px;
  color: #4a5568;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  font-weight: 600;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(0, 0, 0, 0.05);
  letter-spacing: 0.3px;
  position: relative;
  overflow: hidden;
}

:deep(.detail-tabs .el-tabs__item:hover),
:deep(.log-tabs .el-tabs__item:hover) {
  background: rgba(255, 255, 255, 0.98);
  color: #409EFF !important;
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(64, 158, 255, 0.08);
  border-color: rgba(64, 158, 255, 0.2);
}

:deep(.detail-tabs .el-tabs__item.is-active),
:deep(.log-tabs .el-tabs__item.is-active) {
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  color: #fff !important;
  font-weight: 600;
  box-shadow: 0 8px 20px rgba(64, 158, 255, 0.15);
  transform: translateY(-2px);
  border-color: transparent;
}

:deep(.detail-tabs .el-tabs__content),
:deep(.log-tabs .el-tabs__content) {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 24px;
  padding: 40px;
  min-height: 200px;
  margin: 0 16px 16px 16px;
  border: 1px solid rgba(0, 0, 0, 0.05);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
}

:deep(.detail-tabs .el-tabs__nav-wrap::after) {
  display: none;
}

:deep(.detail-tabs .el-tabs__active-bar) {
  display: none;
}

/* 自定义滚动条 */
:deep(.el-scrollbar__bar) {
  opacity: 0.6;
}

:deep(.el-scrollbar__bar:hover) {
  opacity: 0.8;
}

.placeholder-content {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 400px;
}
</style>

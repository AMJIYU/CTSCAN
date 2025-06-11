<script setup lang="ts">
import { ref } from 'vue';
import SystemInfoPanel from './SystemInfoPanel.vue'
import UserInfoPanel from './UserInfoPanel.vue'
import NetworkInfoPanel from './NetworkInfoPanel.vue'
import StartupPanel from './StartupPanel.vue'
import CronTaskPanel from './CronTaskPanel.vue'
import ProcessPanel from './ProcessPanel.vue'
import LoginSuccessPanel from './LoginSuccessPanel.vue'

// 使用ref引用每个选项卡组件
const systemInfoRef = ref();
const userInfoRef = ref();
const networkInfoRef = ref();
const startupRef = ref();
const cronTaskRef = ref();
const processRef = ref();
const patchRef = ref();
const loginSuccessRef = ref();

// 添加重新获取信息的方法
const refreshInfo = () => {
  // 调用所有选项卡组件的刷新方法，重新获取当前所有选项卡的信息
  systemInfoRef.value?.refresh();
  userInfoRef.value?.refresh();
  networkInfoRef.value?.refresh();
  startupRef.value?.refresh();
  cronTaskRef.value?.refresh();
  processRef.value?.refresh();
  patchRef.value?.refresh();
  loginSuccessRef.value?.refresh();
  console.log('重新获取所有选项卡信息');
}

</script>

<template>
  <div class="overview-container">
    <!-- 顶部操作卡片 -->
    <el-row :gutter="24">
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

 

    <!-- 详情标签页 -->
    <div class="details-section">
       <el-tabs type="border-card" shadow="never" class="detail-tabs">
        <el-tab-pane label="系统基本信息">
          <SystemInfoPanel ref="systemInfoRef" />
        </el-tab-pane>
        
        <el-tab-pane label="用户">
          <UserInfoPanel ref="userInfoRef" />
        </el-tab-pane>
        <el-tab-pane label="网络信息">
          <NetworkInfoPanel ref="networkInfoRef" />
        </el-tab-pane>
        <el-tab-pane label="开机启动项">
          <StartupPanel ref="startupRef" />
        </el-tab-pane>
        <el-tab-pane label="任务计划">
          <CronTaskPanel ref="cronTaskRef" />
        </el-tab-pane>
        <el-tab-pane label="进程排查">
          <ProcessPanel ref="processRef" />
        </el-tab-pane>
       
        <el-tab-pane label="登入成功">
          <LoginSuccessPanel ref="loginSuccessRef" />
        </el-tab-pane>
        <el-tab-pane label="登入失败">
          <div>这里显示登入失败日志</div>
        </el-tab-pane>
        <el-tab-pane label="RDP登入">
          <div>这里显示RDP登入日志</div>
        </el-tab-pane>
        <!-- <div class="tab-placeholder">
           请执行扫描或上传文件, 然后选择标签页查看数据.
        </div> -->
      </el-tabs>
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
</style>

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
  background: #fff;
  padding: 32px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.stats-section {
  flex: 0 0 120px;
  margin-bottom: 0;
}

.details-section {
  flex: 1 1 auto;
}

.action-card {
  border: 1px solid #ebeef5;
  border-radius: 8px;
  transition: all 0.3s;
}

.action-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.card-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px;
}

.card-content .text-content {
  flex-grow: 1;
  margin-left: 24px;
}

.card-content h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.card-content p {
  margin: 0;
  font-size: 14px;
  color: #909399;
}

.stats-section h3 {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 24px;
}

.stat-card {
  text-align: center;
  border: none;
  border-radius: 8px;
  padding: 24px 16px;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  line-height: 1.2;
  margin-bottom: 8px;
}

.stat-label {
  color: #606266;
  font-size: 14px;
  font-weight: 500;
}

/* 主选项卡和子选项卡统一美化 */
.detail-tabs,
.log-tabs {
  border: none;
  background: #f6f8fa;
  border-radius: 16px;
  box-shadow: 0 4px 18px 0 rgba(64,158,255,0.10);
  padding: 8px 0 0 0;
}
:deep(.detail-tabs .el-tabs__header),
:deep(.log-tabs .el-tabs__header) {
  background: transparent;
  border-bottom: none;
  margin-bottom: 8px;
  padding-left: 8px;
}
:deep(.detail-tabs .el-tabs__item),
:deep(.log-tabs .el-tabs__item) {
  border-radius: 10px 10px 0 0;
  margin-right: 10px;
  padding: 10px 28px;
  font-size: 16px;
  color: #3a4a5a;
  background: #eaf2ff;
  font-weight: 500;
  transition: background 0.2s, color 0.2s, box-shadow 0.2s;
  box-shadow: 0 1px 2px rgba(64,158,255,0.04);
  border: 1.5px solid transparent;
  letter-spacing: 1px;
}
:deep(.detail-tabs .el-tabs__item:hover),
:deep(.log-tabs .el-tabs__item:hover) {
  background: #d0e7ff;
  color: #1769aa !important;
  font-weight: bold;
  text-shadow: 0 1px 0 #fff, 0 0 2px #eaf2ff;
  box-shadow: 0 2px 8px rgba(64,158,255,0.10);
}
:deep(.detail-tabs .el-tabs__item.is-active),
:deep(.log-tabs .el-tabs__item.is-active) {
  background: linear-gradient(90deg, #409eff 0%, #66b1ff 100%);
  color: #fff !important;
  font-weight: bold;
  box-shadow: 0 4px 16px rgba(64,158,255,0.13);
  border-bottom: 2.5px solid #fff;
  letter-spacing: 1.5px;
}
:deep(.detail-tabs .el-tabs__content),
:deep(.log-tabs .el-tabs__content) {
  background: #fff;
  border-radius: 0 0 16px 16px;
  padding: 32px 24px;
  min-height: 160px;
  box-shadow: 0 2px 8px rgba(64,158,255,0.07);
}
</style>

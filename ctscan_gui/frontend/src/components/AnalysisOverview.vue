<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { GetStats, GetSystemInfo, GetUserInfo, GetNetworkInfo, GetStartupItems, GetCronTasks, GetAllProcesses } from '../../wailsjs/go/main/App';
import SystemInfoPanel from './SystemInfoPanel.vue'
import UserInfoPanel from './UserInfoPanel.vue'
import NetworkInfoPanel from './NetworkInfoPanel.vue'
import StartupPanel from './StartupPanel.vue'
import CronTaskPanel from './CronTaskPanel.vue'
import ProcessPanel from './ProcessPanel.vue'
import PatchPanel from './PatchPanel.vue'

// 定义与后端匹配的统计信息结构
interface Stats {
  total: number;
  network: number;
  startup: number;
  tasks: number;
  process: number;
}

// 统计数据的响应式引用
const stats = ref<Stats>({
  total: 0,
  network: 0,
  startup: 0,
  tasks: 0,
  process: 0,
});

// 系统信息
const sysInfo = ref({
  hostname: '',
  os: '',
  arch: '',
  cpu_cores: 0,
  go_version: ''
})

const userInfo = ref({
  username: '',
  uid: '',
  gid: '',
  home_dir: '',
  name: ''
})

interface NetworkInfo {
  hostname: string;
  ips: string[];
  macs: string[];
  interfaces: string[];
}

const networkInfo = ref<NetworkInfo>({
  hostname: '',
  ips: [],
  macs: [],
  interfaces: []
})

// 统计卡片的配置
const statItems = ref([
    { key: 'total', label: '总异常数', color: '#eaf2ff', iconColor: '#409eff' },
    { key: 'network', label: '网络', color: '#eef8e9', iconColor: '#67c23a' },
    { key: 'startup', label: '启动项', color: '#fef8e7', iconColor: '#e6a23c' },
    { key: 'tasks', label: '任务', color: '#f9eefe', iconColor: '#a262d5' },
    { key: 'process', label: '进程', color: '#ffeae9', iconColor: '#f56c6c' },
]);

const startupItems = ref<{ name: string; path: string }[]>([])
const cronTasks = ref<{ line: string }[]>([])
const processes = ref<{ pid: number; name: string; cmdline: string; user: string; exe: string }[]>([])

// 组件挂载后从 Go 后端获取数据
onMounted(() => {
  GetStats().then(result => {
    stats.value = result;
  });
  GetSystemInfo().then(info => {
    sysInfo.value = info;
  });
  GetUserInfo().then(info => {
    userInfo.value = info;
  });
  GetNetworkInfo().then(info => {
    networkInfo.value = info;
  });
  GetStartupItems().then(list => {
    startupItems.value = list;
  });
  GetCronTasks().then(list => {
    cronTasks.value = list;
  });
  GetAllProcesses().then(list => {
    processes.value = list;
  });
});

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
            <el-button>点击重新开始</el-button>
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

    <!-- 统计部分 -->
    <div class="stats-section">
      <h3>概览统计</h3>
      <el-row :gutter="24">
        <el-col :span="24 / statItems.length" v-for="item in statItems" :key="item.key">
           <el-card shadow="never" class="stat-card" :style="{ backgroundColor: item.color }">
              <div class="stat-value" :style="{ color: item.iconColor }">{{ stats[item.key as keyof Stats] }}</div>
              <div class="stat-label">{{ item.label }}</div>
           </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 详情标签页 -->
    <div class="details-section">
       <el-tabs type="border-card" shadow="never" class="detail-tabs">
        <el-tab-pane label="系统基本信息">
          <SystemInfoPanel />
        </el-tab-pane>
        <el-tab-pane label="系统补丁信息">
          <PatchPanel />
        </el-tab-pane>
        <el-tab-pane label="用户">
          <UserInfoPanel />
        </el-tab-pane>
        <el-tab-pane label="网络信息">
          <NetworkInfoPanel />
        </el-tab-pane>
        <el-tab-pane label="开机启动项">
          <StartupPanel />
        </el-tab-pane>
        <el-tab-pane label="任务计划">
          <CronTaskPanel />
        </el-tab-pane>
        <el-tab-pane label="进程排查">
          <ProcessPanel />
        </el-tab-pane>
        <el-tab-pane label="痕迹&日志"></el-tab-pane>
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

.detail-tabs {
  border: 1px solid #ebeef5;
  box-shadow: none;
  border-radius: 8px;
  overflow: hidden;
}

.tab-placeholder {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 240px;
  color: #909399;
  font-size: 14px;
  background-color: #fafafa;
  border-radius: 0 0 8px 8px;
}

:deep(.el-tabs__header) {
  margin: 0;
  background-color: #fafafa;
  border-bottom: 1px solid #ebeef5;
}

:deep(.el-tabs__item) {
  height: 48px;
  line-height: 48px;
  font-size: 14px;
  font-weight: 500;
}

:deep(.el-tabs__item.is-active) {
  font-weight: 600;
}
</style>

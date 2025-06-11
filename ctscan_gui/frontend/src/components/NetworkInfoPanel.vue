<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetNetworkInfo, GetNetworkConnections } from '../../wailsjs/go/pkg/App'

interface InterfaceStats {
  name: string;
  bytes_sent: number;
  bytes_recv: number;
  packets_sent: number;
  packets_recv: number;
}

interface NetworkInfo {
  hostname: string;
  ips: string[];
  macs: string[];
  interfaces: string[];
  interface_stats: InterfaceStats[];
}

const networkInfo = ref<NetworkInfo>({
  hostname: '',
  ips: [],
  macs: [],
  interfaces: [],
  interface_stats: []
})

const connections = ref<{ proto: string; local_addr: string; remote_addr: string; status: string; pid: number }[]>([])

// 添加 refresh 方法，用于重新获取网络信息
const refresh = async () => {
  try {
    const [info, conns] = await Promise.all([
      GetNetworkInfo(),
      GetNetworkConnections()
    ])
    networkInfo.value = info
    connections.value = conns
  } catch (error) {
    console.error('获取网络信息失败:', error)
  }
}

onMounted(() => {
  refresh()
})

// 暴露 refresh 方法，供父组件调用
defineExpose({ refresh })
</script>

<template>
  <div class="network-info-panel">
    <el-descriptions title="网络基本信息" :column="1" border>
      <el-descriptions-item label="主机名">{{ networkInfo.hostname }}</el-descriptions-item>
      <el-descriptions-item label="网卡">{{ networkInfo.interfaces.join(', ') }}</el-descriptions-item>
      <el-descriptions-item label="IP地址">
        <div v-for="ip in networkInfo.ips" :key="ip" class="ip-item">{{ ip }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="MAC地址">
        <div v-for="mac in networkInfo.macs" :key="mac" class="mac-item">{{ mac }}</div>
      </el-descriptions-item>
    </el-descriptions>

    <h4 style="margin-top:24px;">网卡流量统计</h4>
    <el-table :data="networkInfo.interface_stats" style="width: 100%">
      <el-table-column prop="name" label="网卡名称" />
      <el-table-column label="发送流量">
        <template #default="scope">
          {{ (scope.row.bytes_sent / 1024 / 1024).toFixed(2) }} MB
          ({{ scope.row.packets_sent }} 个数据包)
        </template>
      </el-table-column>
      <el-table-column label="接收流量">
        <template #default="scope">
          {{ (scope.row.bytes_recv / 1024 / 1024).toFixed(2) }} MB
          ({{ scope.row.packets_recv }} 个数据包)
        </template>
      </el-table-column>
    </el-table>

    <h4 style="margin-top:24px;">网络连接详情</h4>
    <el-table :data="connections" style="width: 100%">
      <el-table-column prop="proto" label="协议" />
      <el-table-column prop="local_addr" label="本地地址" />
      <el-table-column prop="remote_addr" label="远程地址" />
      <el-table-column prop="status" label="状态" />
      <el-table-column prop="pid" label="PID" />
    </el-table>
  </div>
</template>

<style scoped>
.network-info-panel {
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

.ip-item, .mac-item {
  margin-bottom: 4px;
  padding: 4px 8px;
  background: rgba(0, 0, 0, 0.02);
  border-radius: 4px;
  font-family: monospace;
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

h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1a202c;
  margin-bottom: 16px;
}
</style> 
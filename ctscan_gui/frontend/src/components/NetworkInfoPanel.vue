<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { GetNetworkInfo, GetNetworkConnections } from '../../wailsjs/go/pkg/App'
import { Monitor, Connection, DataLine } from '@element-plus/icons-vue'

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

// 计算是否有流量统计
const hasTrafficStats = computed(() => {
  return networkInfo.value.interface_stats.some(stat => 
    stat.bytes_sent > 0 || stat.bytes_recv > 0
  )
})

// 过滤出有流量的网卡统计
const activeInterfaceStats = computed(() => {
  return networkInfo.value.interface_stats.filter(stat => 
    stat.bytes_sent > 0 || stat.bytes_recv > 0
  )
})

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
    <!-- 基本信息卡片 -->
    <div class="info-card">
      <div class="card-header">
        <el-icon :size="18" color="#409EFF"><Monitor /></el-icon>
        <h3>网络基本信息</h3>
      </div>
      <div class="interface-list">
        <div v-for="(iface, index) in networkInfo.interfaces" :key="iface" class="interface-item">
          <div class="interface-info">
            <div class="interface-name">
              <el-icon><Connection /></el-icon>
              <span>{{ iface }}</span>
            </div>
            <div class="interface-details">
              <template v-if="networkInfo.ips[index]">
                <span class="label">IP</span>
                <span class="value ip-tag">{{ networkInfo.ips[index] }}</span>
              </template>
              <template v-if="networkInfo.macs[index]">
                <span class="label">MAC</span>
                <span class="value mac-tag">{{ networkInfo.macs[index] }}</span>
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 网卡流量统计 -->
    <div class="info-card" v-if="hasTrafficStats">
      <div class="card-header">
        <el-icon :size="18" color="#409EFF"><DataLine /></el-icon>
        <h3>网卡流量统计</h3>
      </div>
      <el-table :data="activeInterfaceStats" style="width: 100%" size="small">
        <el-table-column prop="name" label="网卡名称" min-width="120">
          <template #default="{ row }">
            <div class="interface-name">
              <el-icon><Connection /></el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="发送流量" min-width="160">
          <template #default="scope">
            <div class="traffic-info">
              <span class="traffic-value">{{ (scope.row.bytes_sent / 1024 / 1024).toFixed(2) }} MB</span>
              <span class="traffic-packets">{{ scope.row.packets_sent }} 包</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="接收流量" min-width="160">
          <template #default="scope">
            <div class="traffic-info">
              <span class="traffic-value">{{ (scope.row.bytes_recv / 1024 / 1024).toFixed(2) }} MB</span>
              <span class="traffic-packets">{{ scope.row.packets_recv }} 包</span>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 网络连接详情 -->
    <div class="info-card">
      <div class="card-header">
        <el-icon :size="18" color="#409EFF"><Connection /></el-icon>
        <h3>网络连接详情</h3>
      </div>
      <el-table :data="connections" style="width: 100%" size="small">
        <el-table-column prop="proto" label="协议" width="65" align="center">
          <template #default="{ row }">
            <el-tag size="small" :type="row.proto === 'tcp' ? 'primary' : 'success'" class="proto-tag">
              {{ row.proto.toUpperCase() }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="local_addr" label="本地地址" min-width="100">
          <template #default="{ row }">
            <span class="address-value">{{ row.local_addr }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="remote_addr" label="远程地址" min-width="100">
          <template #default="{ row }">
            <span class="address-value">{{ row.remote_addr }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="150" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'ESTABLISHED' ? 'success' : 'info'" size="small" class="status-tag">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="pid" label="PID" width="85" align="center">
          <template #default="{ row }">
            <span class="pid-value">{{ row.pid }}</span>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<style scoped>
.network-info-panel {
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.info-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.info-card:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.08);
  transform: translateY(-1px);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.card-header h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: #1a202c;
}

.interface-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.interface-item {
  background: rgba(0, 0, 0, 0.02);
  border-radius: 4px;
  padding: 8px 12px;
  transition: all 0.2s ease;
}

.interface-item:hover {
  background: rgba(64, 158, 255, 0.05);
}

.interface-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.interface-name {
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 120px;
}

.interface-name .el-icon {
  color: #409EFF;
  font-size: 14px;
}

.interface-name span {
  font-weight: 600;
  color: #1a202c;
  font-size: 13px;
}

.interface-details {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  overflow-x: auto;
}

.interface-details .label {
  font-size: 12px;
  color: #718096;
  min-width: 30px;
}

.interface-details .value {
  font-family: monospace;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 3px;
  background: rgba(255, 255, 255, 0.8);
  white-space: nowrap;
}

.interface-details .ip-tag {
  color: #409EFF;
}

.interface-details .mac-tag {
  color: #67C23A;
}

:deep(.el-table) {
  background: transparent;
  border-radius: 6px;
  overflow: hidden;
}

:deep(.el-table th) {
  background-color: rgba(0, 0, 0, 0.02);
  font-weight: 600;
  color: #1a202c;
  padding: 8px 0;
  font-size: 13px;
}

:deep(.el-table td) {
  color: #2d3748;
  padding: 8px 0;
  font-size: 13px;
}

:deep(.el-table .cell) {
  padding-left: 12px;
  padding-right: 12px;
}

:deep(.el-table tr:hover > td) {
  background-color: rgba(64, 158, 255, 0.05);
}

.traffic-info {
  display: flex;
  align-items: center;
  gap: 6px;
}

.traffic-value {
  font-weight: 500;
  color: #2d3748;
  font-size: 13px;
}

.traffic-packets {
  font-size: 12px;
  color: #718096;
  background: rgba(0, 0, 0, 0.02);
  padding: 2px 4px;
  border-radius: 3px;
}

.address-value {
  font-family: monospace;
  font-size: 12px;
  color: #2d3748;
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.pid-value {
  font-family: monospace;
  font-size: 12px;
  color: #409EFF;
  font-weight: 500;
  display: inline-block;
  min-width: 40px;
}

:deep(.el-tag) {
  border-radius: 3px;
  font-size: 12px;
  height: 20px;
  line-height: 18px;
  padding: 0 6px;
  min-width: 45px;
  text-align: center;
}

:deep(.proto-tag) {
  min-width: 40px;
}

:deep(.status-tag) {
  min-width: 70px;
}

:deep(.el-table__header) {
  th {
    background-color: rgba(0, 0, 0, 0.02);
    font-weight: 600;
    color: #1a202c;
  }
}

:deep(.el-table__body) {
  td {
    color: #2d3748;
  }
}

:deep(.el-table__empty-block) {
  background: transparent;
}

:deep(.el-table--small) {
  font-size: 13px;
}
</style> 
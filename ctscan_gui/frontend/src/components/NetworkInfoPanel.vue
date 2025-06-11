<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { GetNetworkInfo, GetNetworkConnections } from '../../wailsjs/go/pkg/App'
import { Monitor, Connection, DataLine, CopyDocument, Filter } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

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
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)

// 筛选条件
const filters = ref({
  proto: '',
  local_addr: '',
  remote_addr: '',
  status: '',
  pid: ''
})

// 重置筛选条件
const resetFilters = () => {
  filters.value = {
    proto: '',
    local_addr: '',
    remote_addr: '',
    status: '',
    pid: ''
  }
  currentPage.value = 1
}

// 筛选后的数据
const filteredConnections = computed(() => {
  return connections.value.filter(conn => {
    return (
      (!filters.value.proto || conn.proto.toLowerCase().includes(filters.value.proto.toLowerCase())) &&
      (!filters.value.local_addr || conn.local_addr.toLowerCase().includes(filters.value.local_addr.toLowerCase())) &&
      (!filters.value.remote_addr || conn.remote_addr.toLowerCase().includes(filters.value.remote_addr.toLowerCase())) &&
      (!filters.value.status || conn.status.toLowerCase().includes(filters.value.status.toLowerCase())) &&
      (!filters.value.pid || conn.pid.toString().includes(filters.value.pid))
    )
  })
})

// 计算当前页的数据
const currentPageData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredConnections.value.slice(start, end)
})

// 更新总数
const updateTotal = () => {
  total.value = filteredConnections.value.length
  if (currentPage.value > Math.ceil(total.value / pageSize.value)) {
    currentPage.value = 1
  }
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
}

// 监听筛选条件变化
watch(filters, () => {
  updateTotal()
}, { deep: true })

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
  loading.value = true
  try {
    const [info, conns] = await Promise.all([
      GetNetworkInfo(),
      GetNetworkConnections()
    ])
    networkInfo.value = info
    connections.value = conns
    updateTotal()
  } catch (error) {
    console.error('获取网络信息失败:', error)
  } finally {
    loading.value = false
  }
}

// 复制地址到剪贴板
const copyAddress = async (address: string) => {
  try {
    await navigator.clipboard.writeText(address)
    ElMessage({
      type: 'success',
      message: '地址已复制到剪贴板',
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
              <div class="basic-info">
                <template v-if="networkInfo.ips[index]">
                  <span class="label">IP</span>
                  <span class="value ip-tag">{{ networkInfo.ips[index] }}</span>
                </template>
                <template v-if="networkInfo.macs[index]">
                  <span class="label">MAC</span>
                  <span class="value mac-tag">{{ networkInfo.macs[index] }}</span>
                </template>
              </div>
              <template v-if="networkInfo.interface_stats[index]">
                <div class="traffic-stats">
                  <div class="traffic-item">
                    <span class="label">发送</span>
                    <span class="value">{{ (networkInfo.interface_stats[index].bytes_sent / 1024 / 1024).toFixed(2) }} MB</span>
                    <span class="packets">{{ networkInfo.interface_stats[index].packets_sent }} 包</span>
                  </div>
                  <div class="traffic-item">
                    <span class="label">接收</span>
                    <span class="value">{{ (networkInfo.interface_stats[index].bytes_recv / 1024 / 1024).toFixed(2) }} MB</span>
                    <span class="packets">{{ networkInfo.interface_stats[index].packets_recv }} 包</span>
                  </div>
                </div>
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 网络连接详情 -->
    <div class="info-card">
      <div class="card-header">
        <el-icon :size="18" color="#409EFF"><Connection /></el-icon>
        <h3>网络连接详情</h3>
        <span class="total-count">共 {{ total }} 个连接</span>
      </div>

      <el-table 
        v-loading="loading"
        element-loading-text="正在加载网络连接信息..."
        element-loading-background="rgba(255, 255, 255, 0.9)"
        :data="currentPageData" 
        style="width: 100%" 
        size="small"
        border
        :resizable="true"
      >
        <el-table-column 
          prop="proto" 
          label="协议" 
          min-width="65" 
          align="center"
          resizable
        >
          <template #header>
            <div class="column-header">
              <span>协议</span>
              <el-input
                v-model="filters.proto"
                placeholder="筛选"
                clearable
                size="small"
                class="header-filter"
              />
            </div>
          </template>
          <template #default="{ row }">
            <el-tag size="small" :type="row.proto === 'tcp' ? 'primary' : 'success'" class="proto-tag">
              {{ row.proto.toUpperCase() }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column 
          prop="local_addr" 
          label="本地地址" 
          min-width="150"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="column-header">
              <span>本地地址</span>
              <el-input
                v-model="filters.local_addr"
                placeholder="筛选"
                clearable
                size="small"
                class="header-filter"
              />
            </div>
          </template>
          <template #default="{ row }">
            <div class="address-cell">
              <span class="address-text">{{ row.local_addr }}</span>
              <el-button
                type="primary"
                link
                size="small"
                class="copy-button"
                @click="copyAddress(row.local_addr)"
              >
                复制
              </el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column 
          prop="remote_addr" 
          label="远程地址" 
          min-width="150"
          resizable
          show-overflow-tooltip
        >
          <template #header>
            <div class="column-header">
              <span>远程地址</span>
              <el-input
                v-model="filters.remote_addr"
                placeholder="筛选"
                clearable
                size="small"
                class="header-filter"
              />
            </div>
          </template>
          <template #default="{ row }">
            <div class="address-cell">
              <span class="address-text">{{ row.remote_addr }}</span>
              <el-button
                type="primary"
                link
                size="small"
                class="copy-button"
                @click="copyAddress(row.remote_addr)"
              >
                复制
              </el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column 
          prop="status" 
          label="状态" 
          min-width="85"
          align="center"
          resizable
        >
          <template #header>
            <div class="column-header">
              <span>状态</span>
              <el-input
                v-model="filters.status"
                placeholder="筛选"
                clearable
                size="small"
                class="header-filter"
              />
            </div>
          </template>
          <template #default="{ row }">
            <el-tag :type="row.status === 'ESTABLISHED' ? 'success' : 'info'" size="small" class="status-tag">
              {{ row.status }}
            </el-tag>
      </template>
    </el-table-column>
        <el-table-column 
          prop="pid" 
          label="PID" 
          min-width="65"
          align="center"
          resizable
        >
          <template #header>
            <div class="column-header">
              <span>PID</span>
              <el-input
                v-model="filters.pid"
                placeholder="筛选"
                clearable
                size="small"
                class="header-filter"
              />
            </div>
          </template>
          <template #default="{ row }">
            <span class="pid-value">{{ row.pid }}</span>
      </template>
    </el-table-column>
  </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
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
  align-items: flex-start;
  gap: 12px;
}

.interface-name {
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 100px;
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
  flex-direction: column;
  gap: 4px;
  flex: 1;
}

.basic-info {
  display: flex;
  align-items: center;
  gap: 12px;
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

.traffic-stats {
  display: flex;
  gap: 12px;
}

.traffic-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.traffic-item .value {
  color: #2d3748;
  font-weight: 500;
}

.traffic-item .packets {
  font-size: 12px;
  color: #718096;
  background: rgba(0, 0, 0, 0.02);
  padding: 2px 4px;
  border-radius: 3px;
}

.address-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  position: relative;
  padding-right: 60px;
}

.address-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
}

.copy-button {
  opacity: 0;
  transition: opacity 0.2s ease;
  flex-shrink: 0;
  position: absolute;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
}

.address-cell:hover .copy-button {
  opacity: 1;
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

:deep(.el-table__column-resize-proxy) {
  background-color: #409EFF;
}

:deep(.el-table__column-resize-handle) {
  background-color: #409EFF;
}

:deep(.el-table__column-resize-handle:hover) {
  background-color: #66b1ff;
}

.total-count {
  margin-left: auto;
  color: #909399;
  font-size: 14px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  padding: 10px 0;
}

:deep(.el-pagination) {
  --el-pagination-button-color: #409EFF;
  --el-pagination-hover-color: #66b1ff;
}

:deep(.el-pagination .el-select .el-input) {
  width: 110px;
}

:deep(.el-pagination .el-pagination__total) {
  margin-right: 16px;
}

:deep(.el-pagination .el-pagination__sizes) {
  margin-right: 16px;
}

:deep(.el-pagination .el-pagination__jump) {
  margin-left: 16px;
}

:deep(.el-pagination .el-pagination__jump .el-input__inner) {
  text-align: center;
}

:deep(.el-loading-mask) {
  backdrop-filter: blur(2px);
}

:deep(.el-loading-spinner) {
  .el-loading-text {
    color: #409EFF;
    font-size: 14px;
    margin-top: 8px;
  }
  
  .circular {
    width: 30px;
    height: 30px;
  }
}

.column-header {
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: center;
}

.column-header span {
  font-weight: 600;
  color: #1a202c;
  font-size: 13px;
}

.header-filter {
  width: 100%;
}

.header-filter :deep(.el-input__wrapper) {
  padding: 0 8px;
  background-color: rgba(0, 0, 0, 0.02);
}

.header-filter :deep(.el-input__inner) {
  height: 24px;
  font-size: 12px;
}

.header-filter :deep(.el-input__inner::placeholder) {
  color: #909399;
}

.filter-container,
.filter-input {
  display: none;
}
</style> 
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
const refresh = () => {
  GetNetworkInfo().then(info => {
    networkInfo.value = info
  })
}

onMounted(() => {
  refresh()
  GetNetworkConnections().then(list => {
    connections.value = list
  })
})

// 暴露 refresh 方法，供父组件调用
defineExpose({ refresh })
</script>
<template>
  <el-descriptions title="网络信息" :column="1" border>
    <el-descriptions-item label="主机名">{{ networkInfo.hostname }}</el-descriptions-item>
    <el-descriptions-item label="网卡">{{ networkInfo.interfaces.join(', ') }}</el-descriptions-item>
    <el-descriptions-item label="IP地址">
      <div v-for="ip in networkInfo.ips" :key="ip">{{ ip }}</div>
    </el-descriptions-item>
    <el-descriptions-item label="MAC地址">
      <div v-for="mac in networkInfo.macs" :key="mac">{{ mac }}</div>
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
</template> 
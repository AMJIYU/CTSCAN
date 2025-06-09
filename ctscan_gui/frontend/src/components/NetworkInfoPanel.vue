<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetNetworkInfo, GetNetworkConnections } from '../../wailsjs/go/main/App'

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

const connections = ref<{ proto: string; local_addr: string; remote_addr: string; status: string; pid: number }[]>([])

onMounted(() => {
  GetNetworkInfo().then(info => {
    networkInfo.value = info
  })
  GetNetworkConnections().then(list => {
    connections.value = list
  })
})
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
  <h4 style="margin-top:24px;">网络连接详情</h4>
  <el-table :data="connections" style="width: 100%">
    <el-table-column prop="proto" label="协议" />
    <el-table-column prop="local_addr" label="本地地址" />
    <el-table-column prop="remote_addr" label="远程地址" />
    <el-table-column prop="status" label="状态" />
    <el-table-column prop="pid" label="PID" />
  </el-table>
</template> 
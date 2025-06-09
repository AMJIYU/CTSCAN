<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetNetworkInfo } from '../../wailsjs/go/main/App'

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

onMounted(() => {
  GetNetworkInfo().then(info => {
    networkInfo.value = info
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
</template> 
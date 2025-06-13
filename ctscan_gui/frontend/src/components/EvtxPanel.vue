<script setup lang="ts">
import { ref, onMounted, computed, PropType } from 'vue'
import { ElMessage, ElLoading } from 'element-plus'
import { Document, Search, Filter, Download, Key, Warning } from '@element-plus/icons-vue'
import { ParseEVTXFile } from '../../wailsjs/go/pkg/App'
import { pkg } from '../../wailsjs/go/models'

// 定义组件事件
const emit = defineEmits(['update:events'])

const events = ref<pkg.EVTXEvent[]>([])
const loading = ref(false)
const searchQuery = ref('')
const selectedLevel = ref('')
const selectedChannel = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const dialogVisible = ref(false)
const selectedEvent = ref<pkg.EVTXEvent | null>(null)
const quickFilter = ref('')

// 分页相关
const total = computed(() => filteredEvents.value.length)

// 获取所有可用的级别和通道
const levels = ref<string[]>([])
const channels = ref<string[]>([])

// 过滤后的事件列表
const filteredEvents = computed(() => {
  if (!events.value) {
    return []
  }
  return events.value.filter(event => {
    const matchesSearch = searchQuery.value === '' || 
      event.description.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      event.provider.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    const matchesLevel = selectedLevel.value === '' || event.level === selectedLevel.value
    const matchesChannel = selectedChannel.value === '' || event.channel === selectedChannel.value
    
    // 快速筛选
    let matchesQuickFilter = true
    if (quickFilter.value === 'login-success') {
      matchesQuickFilter = event.event_id === 4624 || event.event_id === 4648
    } else if (quickFilter.value === 'login-failed') {
      matchesQuickFilter = event.event_id === 4625 || event.event_id === 4647
    }
    
    return matchesSearch && matchesLevel && matchesChannel && matchesQuickFilter
  })
})

// 分页后的事件列表
const paginatedEvents = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredEvents.value.slice(start, end)
})

// 处理页码变化
const handlePageChange = (page: number) => {
  currentPage.value = page
}

// 处理每页条数变化
const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
}

// 解析 EVTX 文件
const parseEvtxFile = async (filePath: string) => {
  loading.value = true
  try {
    const result = await ParseEVTXFile(filePath)
    events.value = result
    
    // 提取所有唯一的级别和通道
    const uniqueLevels = new Set<string>()
    const uniqueChannels = new Set<string>()
    
    result.forEach(event => {
      uniqueLevels.add(event.level)
      uniqueChannels.add(event.channel)
    })
    
    levels.value = Array.from(uniqueLevels)
    channels.value = Array.from(uniqueChannels)
    
    ElMessage({
      type: 'success',
      message: `成功解析 ${result.length} 条事件记录`,
      duration: 2000
    })
  } catch (error) {
    console.error('解析EVTX文件失败:', error)
    ElMessage({
      type: 'error',
      message: error instanceof Error ? error.message : '解析EVTX文件失败',
      duration: 2000
    })
  } finally {
    loading.value = false
  }
}

// 导出事件数据为JSON
const exportEvents = () => {
  const dataStr = JSON.stringify(filteredEvents.value, null, 2)
  const dataUri = 'data:application/json;charset=utf-8,'+ encodeURIComponent(dataStr)
  
  const exportFileDefaultName = `evtx_events_${new Date().toISOString()}.json`
  
  const linkElement = document.createElement('a')
  linkElement.setAttribute('href', dataUri)
  linkElement.setAttribute('download', exportFileDefaultName)
  linkElement.click()
}

// 获取事件级别的样式
const getLevelType = (level: string | undefined) => {
  if (!level) return ''
  const styles = {
    '错误': 'danger',
    '警告': 'warning',
    '信息': 'info',
    '成功': 'success'
  }
  return styles[level as keyof typeof styles] || ''
}

// 获取登入类型描述
const getLogonTypeDescription = (event: pkg.EVTXEvent | null) => {
  if (!event) return ''
  
  // 检查是否是登入相关事件
  if (![4624, 4648, 4625, 4647].includes(event.event_id)) return ''
  
  // 调试输出
  console.log('事件数据:', event.event_data)
  
  // 尝试不同的数据访问路径
  let logonType: string | number | null = null
  
  // 尝试直接从 event_data 获取
  if (event.event_data && typeof event.event_data === 'object') {
    // 遍历所有键
    for (const key in event.event_data) {
      if (key.includes('LogonType')) {
        const value = event.event_data[key]
        if (typeof value === 'string' || typeof value === 'number') {
          logonType = value
          break
        }
      }
    }
  }
  
  if (!logonType) return ''
  
  const logonTypes: { [key: string]: string } = {
    '2': '本地交互式登入',
    '3': '网络登入',
    '4': '批处理登入',
    '5': '服务登入',
    '7': '工作站解锁',
    '8': '网络明文登入',
    '9': '新凭证登入',
    '10': '远程交互式登入 (RDP)',
    '11': '缓存交互式登入'
  }
  
  return logonTypes[String(logonType)] || `未知登入类型 (${logonType})`
}

// 处理行点击
const handleRowClick = (row: pkg.EVTXEvent) => {
  selectedEvent.value = row
  dialogVisible.value = true
}

// 暴露方法给父组件
defineExpose({
  parseEvtxFile,
  setEvents: (ev: pkg.EVTXEvent[]) => {
    console.log('EvtxPanel 收到事件数据:', ev)
    events.value = ev
    // 提取所有唯一的级别和通道
    const uniqueLevels = new Set<string>()
    const uniqueChannels = new Set<string>()
    ev.forEach(event => {
      uniqueLevels.add(event.level)
      uniqueChannels.add(event.channel)
    })
    levels.value = Array.from(uniqueLevels)
    channels.value = Array.from(uniqueChannels)
    console.log('更新后的级别列表:', levels.value)
    console.log('更新后的通道列表:', channels.value)
  },
  refresh: () => {
    return Promise.resolve()
  }
})
</script>

<template>
  <div class="evtx-panel">
    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="search-section">
        <el-input
          v-model="searchQuery"
          placeholder="搜索事件..."
          :prefix-icon="Search"
          clearable
        />
      </div>
      
      <div class="filters">
        <!-- 快速筛选按钮组 -->
        <el-radio-group v-model="quickFilter" size="large">
          <el-radio-button label="">全部</el-radio-button>
          <el-radio-button label="login-success">
            <el-icon><Key /></el-icon>
            登入成功
          </el-radio-button>
          <el-radio-button label="login-failed">
            <el-icon><Warning /></el-icon>
            登入失败
          </el-radio-button>
        </el-radio-group>

        <el-select
          v-model="selectedLevel"
          placeholder="事件级别"
          clearable
          style="width: 120px"
        >
          <el-option
            v-for="level in levels"
            :key="level"
            :label="level"
            :value="level"
          />
        </el-select>
        
        <el-select
          v-model="selectedChannel"
          placeholder="事件通道"
          clearable
          style="width: 150px"
        >
          <el-option
            v-for="channel in channels"
            :key="channel"
            :label="channel"
            :value="channel"
          />
        </el-select>
        
        <el-button
          type="primary"
          :icon="Download"
          @click="exportEvents"
        >
          导出数据
        </el-button>
      </div>
    </div>

    <!-- 事件列表 -->
    <el-table
      v-loading="loading"
      :data="paginatedEvents"
      style="width: 100%"
      height="calc(100vh - 250px)"
      border
      @row-click="handleRowClick"
    >
      <el-table-column
        prop="time"
        label="时间"
        width="180"
        sortable
      />
      
      <el-table-column
        prop="event_id"
        label="事件ID"
        width="100"
        sortable
      />
      
      <el-table-column
        label="登入类型"
        width="150"
      >
        <template #default="{ row }">
          {{ getLogonTypeDescription(row) }}
        </template>
      </el-table-column>
      
      <el-table-column
        prop="level"
        label="级别"
        width="100"
      >
        <template #default="{ row }">
          <el-tag
            :type="getLevelType(row.level)"
            effect="dark"
          >
            {{ row.level }}
          </el-tag>
        </template>
      </el-table-column>
      
      <el-table-column
        prop="channel"
        label="通道"
        width="150"
      />
      
      <el-table-column
        prop="provider"
        label="提供者"
        width="150"
      />
      
      <el-table-column
        prop="description"
        label="描述"
        min-width="200"
        show-overflow-tooltip
      />
      
      <el-table-column
        label="详细信息"
        width="100"
        fixed="right"
      >
        <template #default="{ row }">
          <el-button
            type="primary"
            link
            @click="handleRowClick(row)"
          >
            查看详情
          </el-button>
        </template>
      </el-table-column>
      
      <el-table-column
        type="expand"
        width="1"
      >
        <template #default="{ row }">
          <div class="event-details">
            <el-descriptions
              :column="2"
              border
            >
              <el-descriptions-item label="计算机">
                {{ row.computer }}
              </el-descriptions-item>
              <el-descriptions-item label="用户ID">
                {{ row.user_id }}
              </el-descriptions-item>
            </el-descriptions>
            
            <div class="event-data">
              <h4>事件数据</h4>
              <el-descriptions
                :column="1"
                border
              >
                <el-descriptions-item
                  v-for="(value, key) in row.event_data"
                  :key="key"
                  :label="key"
                >
                  {{ value }}
                </el-descriptions-item>
              </el-descriptions>
            </div>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页器 -->
    <el-pagination
      v-if="total > pageSize"
      background
      layout="prev, pager, next"
      :total="total"
      :current-page="currentPage"
      :page-size="pageSize"
      @current-change="handlePageChange"
      @size-change="handleSizeChange"
    />

    <!-- 事件详情对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="事件详情"
      width="80%"
      :destroy-on-close="true"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="时间 Time">{{ selectedEvent?.time }}</el-descriptions-item>
        <el-descriptions-item label="事件ID EventID">{{ selectedEvent?.event_id }}</el-descriptions-item>
        <el-descriptions-item label="提供者 Provider">{{ selectedEvent?.provider }}</el-descriptions-item>
        <el-descriptions-item label="级别 Level">
          <el-tag :type="getLevelType(selectedEvent?.level)">{{ selectedEvent?.level }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="通道 Channel">{{ selectedEvent?.channel }}</el-descriptions-item>
        <el-descriptions-item label="计算机 Computer">{{ selectedEvent?.computer }}</el-descriptions-item>
        <el-descriptions-item label="用户ID UserID">{{ selectedEvent?.user_id }}</el-descriptions-item>
        <el-descriptions-item label="事件记录ID EventRecordID">{{ selectedEvent?.event_record_id }}</el-descriptions-item>
        <el-descriptions-item label="版本 Version">{{ selectedEvent?.version }}</el-descriptions-item>
        <el-descriptions-item label="限定符 Qualifiers">{{ selectedEvent?.qualifiers }}</el-descriptions-item>
        <el-descriptions-item label="任务 Task">{{ selectedEvent?.task }}</el-descriptions-item>
        <el-descriptions-item label="操作码 Opcode">{{ selectedEvent?.opcode }}</el-descriptions-item>
        <el-descriptions-item label="关键词 Keywords">{{ selectedEvent?.keywords }}</el-descriptions-item>
        <el-descriptions-item label="进程ID ProcessID">{{ selectedEvent?.process_id }}</el-descriptions-item>
        <el-descriptions-item label="线程ID ThreadID">{{ selectedEvent?.thread_id }}</el-descriptions-item>
      </el-descriptions>

      <el-divider>消息 Message</el-divider>
      <div class="message-content">{{ selectedEvent?.message }}</div>

      <el-divider>系统信息 SystemInfo</el-divider>
      <el-descriptions :column="2" border>
        <el-descriptions-item
          v-for="(value, key) in selectedEvent?.system_info"
          :key="key"
          :label="`${key} ${key}`"
        >
          {{ value }}
        </el-descriptions-item>
      </el-descriptions>

      <el-divider>事件数据 EventData</el-divider>
      <el-descriptions :column="2" border>
        <el-descriptions-item
          v-for="(value, key) in selectedEvent?.event_data"
          :key="key"
          :label="`${key} ${key}`"
        >
          {{ value }}
        </el-descriptions-item>
      </el-descriptions>

      <el-divider>用户数据 UserData</el-divider>
      <el-descriptions :column="2" border>
        <el-descriptions-item
          v-for="(value, key) in selectedEvent?.user_data"
          :key="key"
          :label="`${key} ${key}`"
        >
          {{ value }}
        </el-descriptions-item>
      </el-descriptions>

      <el-divider>描述 Description</el-divider>
      <div class="description-content">{{ selectedEvent?.description }}</div>
    </el-dialog>
  </div>
</template>

<style scoped>
.evtx-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}

.search-section {
  flex: 1;
  max-width: 400px;
}

.filters {
  display: flex;
  gap: 12px;
  align-items: center;
}

/* 快速筛选按钮组样式 */
:deep(.el-radio-group) {
  margin-right: 16px;
}

:deep(.el-radio-button__inner) {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 16px;
}

:deep(.el-radio-button__inner .el-icon) {
  margin-right: 4px;
}

:deep(.el-radio-button:first-child .el-radio-button__inner) {
  border-radius: 4px 0 0 4px;
}

:deep(.el-radio-button:last-child .el-radio-button__inner) {
  border-radius: 0 4px 4px 0;
}

:deep(.el-radio-button__original-radio:checked + .el-radio-button__inner) {
  background-color: var(--el-color-primary);
  border-color: var(--el-color-primary);
  box-shadow: -1px 0 0 0 var(--el-color-primary);
}

.event-details {
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
}

.event-data {
  margin-top: 16px;
}

.event-data h4 {
  margin: 0 0 12px 0;
  color: #606266;
  font-size: 14px;
}

:deep(.el-table) {
  --el-table-border-color: #e4e7ed;
  --el-table-header-bg-color: #f5f7fa;
}

:deep(.el-table__expanded-cell) {
  background: #f8f9fa;
}

:deep(.el-descriptions) {
  margin-bottom: 16px;
}

:deep(.el-descriptions__label) {
  width: 120px;
  font-weight: 500;
}

:deep(.el-pagination) {
  margin-top: 16px;
  justify-content: center;
  padding: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}

.message-content,
.description-content {
  white-space: pre-wrap;
  word-break: break-all;
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
  margin: 10px 0;
  font-family: monospace;
}

.el-divider {
  margin: 20px 0;
}
</style>

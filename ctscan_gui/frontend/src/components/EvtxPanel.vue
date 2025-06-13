<script setup lang="ts">
import { ref,  computed,  } from 'vue'
import { ElMessage,  } from 'element-plus'
import { Search,  Key, Warning } from '@element-plus/icons-vue'
import { ParseEVTXFile } from '../../wailsjs/go/pkg/App'
import { pkg } from '../../wailsjs/go/models'

// 定义组件事件
const emit = defineEmits(['update:events'])

const events = ref<pkg.EVTXEvent[]>([])
const loading = ref(false)
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const dialogVisible = ref(false)
const selectedEvent = ref<pkg.EVTXEvent | null>(null)
const quickFilter = ref('')

// 分页相关
const total = computed(() => filteredEvents.value.length)

// 过滤后的事件列表
const filteredEvents = computed(() => {
  if (!events.value) {
    return []
  }
  return events.value.filter(event => {
    const matchesSearch = searchQuery.value === '' || 
      event.description.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      event.provider.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    // 快速筛选
    let matchesQuickFilter = true
    if (quickFilter.value === 'login-success') {
      matchesQuickFilter = event.event_id === 4624 || event.event_id === 4648
    } else if (quickFilter.value === 'login-failed') {
      matchesQuickFilter = event.event_id === 4625 || event.event_id === 4647
    }
    
    return matchesSearch && matchesQuickFilter
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
    console.log('更新后的事件列表:', events.value)
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
      :cell-style="{ padding: '4px 0' }"
      :header-cell-style="{ padding: '8px 0' }"
      class="custom-table"
    >
      <el-table-column
        prop="time"
        label="时间"
        width="150"
        sortable
      />
      
      <el-table-column
        prop="event_id"
        label="事件ID"
        width="80"
        sortable
      />
      
      <el-table-column
        prop="channel"
        label="事件通道"
        width="120"
      />
      
      <el-table-column
        label="登入类型"
        width="120"
      >
        <template #default="{ row }">
          {{ getLogonTypeDescription(row) }}
        </template>
      </el-table-column>

      <el-table-column
        label="用户名"
        width="120"
      >
        <template #default="{ row }">
          <div v-if="[4624, 4648, 4625, 4647].includes(row.event_id)">
            {{ row.event_data?.TargetUserName || '-' }}
          </div>
        </template>
      </el-table-column>

      <el-table-column
        label="域名"
        width="120"
      >
        <template #default="{ row }">
          <div v-if="[4624, 4648, 4625, 4647].includes(row.event_id)">
            {{ row.event_data?.TargetDomainName || '-' }}
          </div>
        </template>
      </el-table-column>

      <el-table-column
        label="来源IP"
        width="120"
      >
        <template #default="{ row }">
          <div v-if="[4624, 4648, 4625, 4647].includes(row.event_id)">
            {{ row.event_data?.IpAddress || '-' }}
          </div>
        </template>
      </el-table-column>

      <el-table-column
        label="工作站"
        width="120"
      >
        <template #default="{ row }">
          <div v-if="[4624, 4648, 4625, 4647].includes(row.event_id)">
            {{ row.event_data?.WorkstationName || '-' }}
          </div>
        </template>
      </el-table-column>
      
      <el-table-column
        label="详细信息"
        width="80"
        fixed="right"
      >
        <template #default="{ row }">
          <el-button
            type="primary"
            link
            @click="handleRowClick(row)"
          >
            查看
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
    <div class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </div>

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

      <el-divider>描述 Description</el-divider>
      <div class="description-content">{{ selectedEvent?.description }}</div>

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

:deep(.custom-table) {
  --el-table-border-color: #ebeef5;
  --el-table-header-bg-color: #f8fafc;
  --el-table-row-height: 40px;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

:deep(.custom-table .el-table__header) {
  background: #f8fafc;
}

:deep(.custom-table .el-table__header th) {
  background: #f8fafc;
  color: #1f2937;
  font-weight: 600;
  font-size: 13px;
  border-bottom: 1px solid #e5e7eb;
}

:deep(.custom-table .el-table__row) {
  height: 40px;
  transition: all 0.2s ease;
}

:deep(.custom-table .el-table__row:hover) {
  background-color: #f1f5f9 !important;
}

:deep(.custom-table .el-table__row td) {
  border-bottom: 1px solid #f1f5f9;
  color: #4b5563;
  font-size: 13px;
}

:deep(.custom-table .cell) {
  padding: 0 8px;
}

:deep(.custom-table .el-button--primary.is-link) {
  font-size: 13px;
  padding: 2px 4px;
}

:deep(.custom-table .el-table__expand-icon) {
  color: #6b7280;
}

:deep(.custom-table .el-table__expand-icon--expanded) {
  transform: rotate(90deg);
}

:deep(.custom-table .el-table__expanded-cell) {
  background: #f8fafc;
  padding: 16px;
}

:deep(.custom-table .el-table__expanded-cell .el-descriptions) {
  margin-bottom: 12px;
}

:deep(.custom-table .el-table__expanded-cell .el-descriptions__label) {
  color: #6b7280;
  font-weight: 500;
}

:deep(.custom-table .el-table__expanded-cell .el-descriptions__content) {
  color: #1f2937;
}

:deep(.custom-table .el-tag) {
  border-radius: 4px;
  padding: 0 6px;
  height: 20px;
  line-height: 18px;
  font-size: 12px;
}

:deep(.custom-table .el-tag--dark) {
  border: none;
}

:deep(.custom-table .el-tag--success) {
  background-color: #10b981;
}

:deep(.custom-table .el-tag--warning) {
  background-color: #f59e0b;
}

:deep(.custom-table .el-tag--danger) {
  background-color: #ef4444;
}

:deep(.custom-table .el-tag--info) {
  background-color: #6b7280;
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

.login-info {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  line-height: 1.4;
}

.login-info .label {
  color: #909399;
  font-weight: 500;
}

.login-info .value {
  color: #303133;
}
</style>

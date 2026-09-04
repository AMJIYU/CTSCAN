<script setup lang="ts">
import { ref, computed, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Key, Warning, Connection, RefreshLeft } from '@element-plus/icons-vue'
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

const loginEventIds = [4624, 4625, 4648]
const humanLogonTypes = new Set(['2', '3', '7', '8', '9', '10', '11'])
const remoteLogonTypes = new Set(['3', '10'])
type TagType = '' | 'success' | 'warning' | 'info' | 'danger'
type ColumnFilterKey = 'time' | 'eventId' | 'eventType' | 'logonType' | 'sourceIp' | 'userName' | 'workstation' | 'subjectUserName' | 'subjectDomain' | 'process'

const columnFilterDefinitions: Array<{ key: ColumnFilterKey; label: string; placeholder: string }> = [
  { key: 'time', label: '时间', placeholder: '2026-09-03' },
  { key: 'eventId', label: '事件ID', placeholder: '4624' },
  { key: 'eventType', label: '事件类型', placeholder: 'RDP 登录成功' },
  { key: 'logonType', label: '登录类型', placeholder: '10 / RDP' },
  { key: 'sourceIp', label: '源IP', placeholder: '192.168.17.57' },
  { key: 'userName', label: '用户名', placeholder: 'jiyu' },
  { key: 'workstation', label: '工作站', placeholder: 'YU' },
  { key: 'subjectUserName', label: '主体用户名', placeholder: 'YU$' },
  { key: 'subjectDomain', label: '主体域', placeholder: 'WORKGROUP' },
  { key: 'process', label: '进程', placeholder: 'User32' }
]

const columnFilters = reactive<Record<ColumnFilterKey, string>>({
  time: '',
  eventId: '',
  eventType: '',
  logonType: '',
  sourceIp: '',
  userName: '',
  workstation: '',
  subjectUserName: '',
  subjectDomain: '',
  process: ''
})

// 分页相关
const total = computed(() => filteredEvents.value.length)
const hasColumnFilters = computed(() => Object.values(columnFilters).some(value => value.trim() !== ''))

// 过滤后的事件列表
const filteredEvents = computed(() => {
  if (!events.value) {
    return []
  }
  const query = searchQuery.value.trim().toLowerCase()
  const matched = events.value.filter(event => {
    const matchesSearch = query === '' || getSearchableEventText(event).includes(query)
    
    // 快速筛选
    let matchesQuickFilter = true
    if (quickFilter.value === 'login-success') {
      matchesQuickFilter = isHumanLogonSuccess(event)
    } else if (quickFilter.value === 'login-failed') {
      matchesQuickFilter = event.event_id === 4625
    } else if (quickFilter.value === 'rdp-login') {
      matchesQuickFilter = isRemoteLogon(event)
    }
    
    return matchesSearch && matchesQuickFilter && matchesColumnFilters(event)
  })

  return sortEventsByNewest(matched)
})

// 分页后的事件列表
const paginatedEvents = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredEvents.value.slice(start, end)
})

watch([searchQuery, quickFilter], () => {
  currentPage.value = 1
})

watch(columnFilters, () => {
  currentPage.value = 1
}, { deep: true })

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

const normalizeDisplayValue = (value: unknown) => {
  if (value === undefined || value === null) return '-'
  const text = String(value).trim()
  return text === '' ? '-' : text
}

const getEventDataValue = (event: pkg.EVTXEvent | null, key: string) => {
  if (!event?.event_data || typeof event.event_data !== 'object') return '-'
  return normalizeDisplayValue(event.event_data[key])
}

const getLogonTypeValue = (event: pkg.EVTXEvent | null) => {
  return getEventDataValue(event, 'LogonType')
}

const isLoginEvent = (event: pkg.EVTXEvent | null) => {
  return !!event && loginEventIds.includes(event.event_id)
}

const isNoiseAccount = (event: pkg.EVTXEvent) => {
  const user = getEventDataValue(event, 'TargetUserName').toUpperCase()
  return user === '-' || user === 'SYSTEM' || user.endsWith('$')
}

const isHumanLogonSuccess = (event: pkg.EVTXEvent) => {
  if (event.event_id !== 4624 || isNoiseAccount(event)) return false
  return humanLogonTypes.has(getLogonTypeValue(event))
}

const isRemoteLogon = (event: pkg.EVTXEvent) => {
  if (event.event_id !== 4624 || isNoiseAccount(event)) return false
  return remoteLogonTypes.has(getLogonTypeValue(event)) && getEventDataValue(event, 'IpAddress') !== '-'
}

const getSearchableEventText = (event: pkg.EVTXEvent) => {
  const values = [
    event.time,
    event.time_utc,
    event.time_local,
    event.event_id,
    getEventTypeLabel(event),
    getLogonTypeDisplay(event),
    event.provider,
    event.level,
    event.channel,
    event.computer,
    event.user_id,
    event.description,
    event.message,
    ...Object.values(event.event_data || {}),
    ...Object.values(event.system_info || {}),
    ...Object.values(event.user_data || {})
  ]
  return values.map(normalizeDisplayValue).join(' ').toLowerCase()
}

const sortEventsByNewest = (list: pkg.EVTXEvent[]) => {
  return [...list].sort((a, b) => {
    const timeCompare = normalizeDisplayValue(b.time).localeCompare(normalizeDisplayValue(a.time))
    if (timeCompare !== 0) return timeCompare
    return (b.event_record_id || 0) - (a.event_record_id || 0)
  })
}

const getEventTypeLabel = (event: pkg.EVTXEvent | null) => {
  if (!event) return '-'
  if (event.event_type) return event.event_type

  if (event.event_id === 4624) {
    const logonType = getLogonTypeValue(event)
    if (logonType === '10') return 'RDP 登录成功'
    if (logonType === '3') return '网络登录成功'
    if (logonType === '2') return '本地登录成功'
    if (logonType === '5') return '服务登录成功'
    return '登录成功'
  }
  if (event.event_id === 4625) return '登录失败'
  if (event.event_id === 4648) return '显式凭据登录'
  if (event.event_id === 4688) return '进程创建'
  if (event.event_id === 4720) return '用户创建'
  if (event.event_id === 1102) return '安全日志清除'
  return '普通安全事件'
}

const getEventTypeTag = (event: pkg.EVTXEvent | null): TagType => {
  if (!event) return ''
  if (event.event_id === 4624) return 'success'
  if (event.event_id === 4625 || event.event_id === 1102) return 'danger'
  if ([4648, 4672, 4697, 4698, 4719, 4720, 4724, 4732].includes(event.event_id)) return 'warning'
  return 'info'
}

const getEventTypeClass = (event: pkg.EVTXEvent | null) => {
  return `event-type-tag event-type-tag--${getEventTypeTag(event) || 'info'}`
}

const getLogonTypeDisplay = (event: pkg.EVTXEvent | null) => {
  const logonType = getLogonTypeValue(event)
  if (logonType === '-') return '-'
  const description = getLogonTypeDescription(event)
  return description ? `${logonType} / ${description}` : logonType
}

const getColumnFilterValue = (event: pkg.EVTXEvent, key: ColumnFilterKey) => {
  switch (key) {
    case 'time':
      return normalizeDisplayValue(event.time)
    case 'eventId':
      return normalizeDisplayValue(event.event_id)
    case 'eventType':
      return getEventTypeLabel(event)
    case 'logonType':
      return getLogonTypeDisplay(event)
    case 'sourceIp':
      return getEventDataValue(event, 'IpAddress')
    case 'userName':
      return getEventDataValue(event, 'TargetUserName')
    case 'workstation':
      return getEventDataValue(event, 'WorkstationName')
    case 'subjectUserName':
      return getEventDataValue(event, 'SubjectUserName')
    case 'subjectDomain':
      return getEventDataValue(event, 'SubjectDomainName')
    case 'process':
      return getEventDataValue(event, 'LogonProcessName')
    default:
      return '-'
  }
}

const matchesColumnFilters = (event: pkg.EVTXEvent) => {
  return columnFilterDefinitions.every(({ key }) => {
    const filter = columnFilters[key].trim().toLowerCase()
    if (filter === '') return true
    return getColumnFilterValue(event, key).toLowerCase().includes(filter)
  })
}

const clearColumnFilters = () => {
  columnFilterDefinitions.forEach(({ key }) => {
    columnFilters[key] = ''
  })
}

// 获取登入类型描述
const getLogonTypeDescription = (event: pkg.EVTXEvent | null) => {
  if (!isLoginEvent(event)) return ''
  
  const logonType = getLogonTypeValue(event)
  if (logonType === '-') return ''
  
  const logonTypes: { [key: string]: string } = {
    '2': '本地交互式登录',
    '3': '网络登录',
    '4': '批处理登录',
    '5': '服务登录',
    '7': '工作站解锁',
    '8': '网络明文登录',
    '9': '新凭据登录',
    '10': '远程交互式登录 (RDP)',
    '11': '缓存交互式登录'
  }
  
  return logonTypes[logonType] || `未知登入类型 (${logonType})`
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
    events.value = ev
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
            登录成功
          </el-radio-button>
          <el-radio-button label="rdp-login">
            <el-icon><Connection /></el-icon>
            RDP登录
          </el-radio-button>
          <el-radio-button label="login-failed">
            <el-icon><Warning /></el-icon>
            登录失败
          </el-radio-button>
        </el-radio-group>
      </div>
    </div>

    <div class="column-filter-bar">
      <div class="column-filter-heading">
        <span>列筛选</span>
        <el-button
          :icon="RefreshLeft"
          size="small"
          text
          :disabled="!hasColumnFilters"
          @click="clearColumnFilters"
        >
          清空
        </el-button>
      </div>
      <div class="column-filter-grid">
        <label
          v-for="filter in columnFilterDefinitions"
          :key="filter.key"
          class="column-filter-item"
        >
          <span>{{ filter.label }}</span>
          <el-input
            v-model="columnFilters[filter.key]"
            :placeholder="filter.placeholder"
            size="small"
            clearable
          />
        </label>
      </div>
    </div>

    <!-- 事件列表 -->
    <el-table
      v-loading="loading"
      :data="paginatedEvents"
      style="width: 100%"
      height="calc(100vh - 360px)"
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
        show-overflow-tooltip
      />
      
      <el-table-column
        prop="event_id"
        label="事件ID"
        width="80"
        sortable
        show-overflow-tooltip
      />

      <el-table-column
        label="事件类型"
        width="165"
        show-overflow-tooltip
      >
        <template #default="{ row }">
          <el-tag :class="getEventTypeClass(row)" disable-transitions>
            {{ getEventTypeLabel(row) }}
          </el-tag>
        </template>
      </el-table-column>
      
      <el-table-column
        label="登录类型"
        width="170"
        show-overflow-tooltip
      >
        <template #default="{ row }">
          {{ getLogonTypeDisplay(row) }}
        </template>
      </el-table-column>

      <el-table-column
        label="源IP"
        width="140"
        show-overflow-tooltip
      >
        <template #default="{ row }">
          {{ getEventDataValue(row, 'IpAddress') }}
        </template>
      </el-table-column>

      <el-table-column
        label="用户名"
        width="140"
        show-overflow-tooltip
      >
        <template #default="{ row }">
          {{ getEventDataValue(row, 'TargetUserName') }}
        </template>
      </el-table-column>

      <el-table-column
        label="工作站"
        width="120"
        show-overflow-tooltip
      >
        <template #default="{ row }">
          {{ getEventDataValue(row, 'WorkstationName') }}
        </template>
      </el-table-column>

      <el-table-column
        label="主体用户名"
        width="120"
        show-overflow-tooltip
      >
        <template #default="{ row }">
          {{ getEventDataValue(row, 'SubjectUserName') }}
        </template>
      </el-table-column>

      <el-table-column
        label="主体域"
        width="140"
        show-overflow-tooltip
      >
        <template #default="{ row }">
          {{ getEventDataValue(row, 'SubjectDomainName') }}
        </template>
      </el-table-column>

      <el-table-column
        label="进程"
        width="120"
        show-overflow-tooltip
      >
        <template #default="{ row }">
          {{ getEventDataValue(row, 'LogonProcessName') }}
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
        <el-descriptions-item label="UTC 时间 UTC Time">{{ selectedEvent?.time_utc || '-' }}</el-descriptions-item>
        <el-descriptions-item label="本机换算 Local Time">{{ selectedEvent?.time_local || '-' }}</el-descriptions-item>
        <el-descriptions-item label="事件ID EventID">{{ selectedEvent?.event_id }}</el-descriptions-item>
        <el-descriptions-item label="事件类型 EventType">
          <el-tag :class="getEventTypeClass(selectedEvent)" disable-transitions>
            {{ getEventTypeLabel(selectedEvent) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="登录类型 LogonType">{{ getLogonTypeDisplay(selectedEvent) }}</el-descriptions-item>
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

.column-filter-bar {
  padding: 12px 14px;
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.06);
}

.column-filter-heading {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  color: #334155;
  font-size: 13px;
  font-weight: 700;
}

.column-filter-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(145px, 1fr));
  gap: 10px;
}

.column-filter-item {
  display: flex;
  flex-direction: column;
  gap: 5px;
  min-width: 0;
}

.column-filter-item span {
  color: #64748b;
  font-size: 12px;
  line-height: 1;
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
  min-height: 22px;
  line-height: 20px;
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

:deep(.event-type-tag) {
  max-width: 100%;
  height: 22px;
  padding: 0 8px;
  border: 0 !important;
  color: #ffffff !important;
  font-weight: 700;
  letter-spacing: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

:deep(.event-type-tag--success) {
  background: #0f766e !important;
}

:deep(.event-type-tag--warning) {
  background: #9a3412 !important;
}

:deep(.event-type-tag--danger) {
  background: #b91c1c !important;
}

:deep(.event-type-tag--info) {
  background: #334155 !important;
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

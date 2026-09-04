<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Box, FolderOpened, Refresh, VideoPlay } from '@element-plus/icons-vue'
import {
  GetSysinternalsTools,
  InstallSysinternalsTool,
  LaunchSysinternalsTool,
  OpenSysinternalsToolsFolder
} from '../../wailsjs/go/pkg/App'
import { pkg } from '../../wailsjs/go/models'

const tools = ref<pkg.SysinternalsTool[]>([])
const loading = ref(false)
const activeToolId = ref('')

const refresh = async () => {
  loading.value = true
  try {
    tools.value = await GetSysinternalsTools()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '获取微软工具状态失败')
  } finally {
    loading.value = false
  }
}

const runToolTask = async (tool: pkg.SysinternalsTool, task: () => Promise<void>, successMessage: string) => {
  activeToolId.value = tool.id
  try {
    await task()
    ElMessage.success(successMessage)
    await refresh()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '操作失败')
  } finally {
    activeToolId.value = ''
  }
}

const launchTool = (tool: pkg.SysinternalsTool) => {
  runToolTask(tool, () => LaunchSysinternalsTool(tool.id), `${tool.name} 已启动`)
}

const installTool = (tool: pkg.SysinternalsTool) => {
  runToolTask(tool, async () => {
    await InstallSysinternalsTool(tool.id)
  }, `${tool.name} 已释放`)
}

const openToolsFolder = async () => {
  try {
    await OpenSysinternalsToolsFolder()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '打开工具目录失败')
  }
}

onMounted(refresh)

defineExpose({ refresh })
</script>

<template>
  <div class="sysinternals-panel" v-loading="loading">
    <div class="panel-toolbar">
      <div class="toolbar-copy">
        <h2>微软排查工具</h2>
        <p>使用随 CTScan 封装的原版 EXE，启动后保留工具自身全部功能。</p>
      </div>
      <div class="toolbar-actions">
        <el-button :icon="FolderOpened" @click="openToolsFolder">工具目录</el-button>
        <el-button :icon="Refresh" @click="refresh">刷新</el-button>
      </div>
    </div>

    <div class="tool-grid">
      <section
        v-for="tool in tools"
        :key="tool.id"
        class="tool-card"
      >
        <div class="tool-card-header">
          <div>
            <h3>{{ tool.name }}</h3>
            <span>{{ tool.vendor }}</span>
          </div>
          <el-tag :type="tool.available ? 'success' : 'info'" effect="dark">
            {{ tool.available ? '已释放' : '已封装' }}
          </el-tag>
        </div>

        <p class="tool-description">{{ tool.description }}</p>

        <dl class="tool-meta">
          <div>
            <dt>用途</dt>
            <dd>{{ tool.best_for }}</dd>
          </div>
          <div>
            <dt>文件</dt>
            <dd>{{ tool.file_name }}</dd>
          </div>
          <div>
            <dt>封装</dt>
            <dd>{{ tool.packaged ? '已内置到 CTScan' : '未找到内置资源' }}</dd>
          </div>
          <div>
            <dt>权限</dt>
            <dd>{{ tool.requires_admin ? '管理员权限启动' : '普通权限启动' }}</dd>
          </div>
          <div>
            <dt>路径</dt>
            <dd>{{ tool.local_path || '未找到本地文件' }}</dd>
          </div>
        </dl>

        <div class="tool-actions">
          <el-button
            type="primary"
            :icon="VideoPlay"
            :disabled="!tool.supported"
            :loading="activeToolId === tool.id"
            @click="launchTool(tool)"
          >
            启动
          </el-button>
          <el-button
            :icon="Box"
            :disabled="!tool.supported"
            :loading="activeToolId === tool.id"
            @click="installTool(tool)"
          >
            释放/更新
          </el-button>
        </div>

        <div v-if="!tool.supported" class="platform-note">
          当前系统不能运行 Windows Sysinternals EXE；打包到 Windows 后可直接使用。
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.sysinternals-panel {
  display: flex;
  flex-direction: column;
  gap: 18px;
  min-height: 520px;
}

.panel-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.06);
}

.toolbar-copy {
  text-align: left;
}

.toolbar-copy h2 {
  margin: 0;
  color: #0f172a;
  font-size: 18px;
  line-height: 1.3;
}

.toolbar-copy p {
  margin: 6px 0 0;
  color: #64748b;
  font-size: 13px;
}

.toolbar-actions {
  display: flex;
  gap: 10px;
  flex-shrink: 0;
}

.tool-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 16px;
}

.tool-card {
  display: flex;
  flex-direction: column;
  gap: 14px;
  min-width: 0;
  padding: 18px;
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.06);
  text-align: left;
}

.tool-card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.tool-card-header h3 {
  margin: 0;
  color: #111827;
  font-size: 17px;
  line-height: 1.25;
}

.tool-card-header span {
  display: inline-block;
  margin-top: 4px;
  color: #64748b;
  font-size: 12px;
}

.tool-description {
  margin: 0;
  color: #334155;
  font-size: 13px;
  line-height: 1.6;
}

.tool-meta {
  display: grid;
  gap: 10px;
  margin: 0;
}

.tool-meta div {
  display: grid;
  grid-template-columns: 52px minmax(0, 1fr);
  gap: 8px;
}

.tool-meta dt {
  color: #64748b;
  font-size: 12px;
}

.tool-meta dd {
  margin: 0;
  color: #1f2937;
  font-size: 12px;
  line-height: 1.45;
  overflow-wrap: anywhere;
}

.tool-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: auto;
}

.platform-note {
  padding: 8px 10px;
  color: #92400e;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 6px;
  font-size: 12px;
  line-height: 1.5;
}
</style>

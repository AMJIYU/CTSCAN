<script setup lang="ts">
import { ref } from 'vue';
import AnalysisOverview from './components/AnalysisOverview.vue';
import {
  Setting,
  VideoCamera,
  Document,
  Help,
  DataAnalysis,
} from '@element-plus/icons-vue';

// 默认激活的菜单项
const activeMenu = ref('overview');
</script>

<template>
  <div class="app-root">
    <!-- 顶部横向菜单栏 -->
    <div class="app-menubar">
      <div class="menubar-left">
        <span class="logo">CT<span>Scan</span></span>
      </div>
      <div class="menubar-right">
        <div class="header-right">
          <span class="welcome-text">你好, 欢迎使用 CTScan!</span>
          <el-button size="small" class="lang-btn">En</el-button>
        </div>
      </div>
    </div>
    <el-container class="main-container">
      <!-- 主内容区 -->
      <el-main class="main-content-wrapper">
        <div class="content-body">
          <!-- 渲染主视图 -->
          <AnalysisOverview />
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<style>
.app-root {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #f6f8fc 0%, #f1f5f9 100%);
}

.app-menubar {
  height: 48px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;
  font-size: 14px;
  user-select: none;
  position: relative;
  z-index: 10;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.menubar-left {
  display: flex;
  align-items: center;
}

.logo {
  font-weight: 700;
  font-size: 20px;
  color: #1a202c;
  letter-spacing: -0.5px;
  position: relative;
}

.logo span {
  color: #409EFF;
  font-weight: 800;
}

.main-container {
  flex: 1;
  height: calc(100vh - 48px);
  background: transparent;
  position: relative;
  padding: 24px;
}

.main-content-wrapper {
  padding: 0;
  display: flex;
  flex-direction: column;
  background: transparent;
  min-height: 0;
  height: 100%;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
  font-size: 14px;
  color: #4a5568;
  font-weight: 500;
}

.welcome-text {
  color: #4a5568;
  font-weight: 500;
}

.lang-btn {
  margin-left: 8px;
  padding: 4px 12px;
  font-weight: 500;
  transition: all 0.3s ease;
  background: transparent;
  color: #4a5568;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 13px;
}

.lang-btn:hover {
  color: #409EFF;
  border-color: #409EFF;
  background: rgba(64, 158, 255, 0.05);
}

.content-body {
  padding: 0;
  background: transparent;
  flex-grow: 1;
  overflow-y: auto;
  max-width: 1600px;
  margin: 0 auto;
  width: 100%;
}

/* 自定义滚动条样式 */
.content-body::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.content-body::-webkit-scrollbar-track {
  background: transparent;
}

.content-body::-webkit-scrollbar-thumb {
  background: #e2e8f0;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.content-body::-webkit-scrollbar-thumb:hover {
  background: #cbd5e1;
}
</style>

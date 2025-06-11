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
        <span class="logo">CTScan</span>
        <span class="menu-item">编辑</span>
        <span class="menu-item">视图</span>
        <span class="menu-item">窗口</span>
        <span class="menu-item">帮助</span>
      </div>
      <div class="menubar-right">
        <div class="header-right">
          你好, 欢迎使用 CTScan!
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
  background-color: #f8fafc;
}

.app-menubar {
  height: 56px;
  background: #ffffff;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 40px;
  font-size: 15px;
  user-select: none;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
  position: relative;
  z-index: 10;
}

.menubar-left {
  display: flex;
  align-items: center;
  gap: 32px;
}

.logo {
  font-weight: 700;
  font-size: 20px;
  color: #409EFF;
  margin-right: 24px;
  letter-spacing: 0.5px;
  background: linear-gradient(135deg, #409EFF 0%, #66b1ff 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.menu-item {
  color: #4a5568;
  cursor: pointer;
  padding: 8px 16px;
  border-radius: 8px;
  transition: all 0.3s ease;
  font-weight: 500;
}

.menu-item:hover {
  background: #f0f7ff;
  color: #409EFF;
  transform: translateY(-1px);
}

.main-container {
  flex: 1;
  height: calc(100vh - 56px);
  background-color: #f8fafc;
  position: relative;
}

.main-content-wrapper {
  padding: 0;
  display: flex;
  flex-direction: column;
  background-color: #f8fafc;
  min-height: 0;
  height: 100%;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
  font-size: 15px;
  color: #4a5568;
  font-weight: 500;
}

.lang-btn {
  margin-left: 8px;
  border-radius: 8px;
  padding: 4px 16px;
  font-weight: 500;
  transition: all 0.3s ease;
  border: 1px solid #e2e8f0;
  background: #ffffff;
  color: #4a5568;
}

.lang-btn:hover {
  background: #f0f7ff;
  color: #409EFF;
  border-color: #409EFF;
  transform: translateY(-1px);
}

.content-body {
  padding: 24px;
  background-color: #f8fafc;
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
  background: #f1f5f9;
  border-radius: 4px;
}

.content-body::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.content-body::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}
</style>

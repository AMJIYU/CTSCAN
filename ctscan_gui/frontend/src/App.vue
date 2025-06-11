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
  height: 68px;
  background: #ffffff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 40px;
  font-size: 14px;
  user-select: none;
  position: relative;
  z-index: 10;
}

.menubar-left {
  display: flex;
  align-items: center;
  gap: 40px;
}

.logo {
  font-weight: 700;
  font-size: 22px;
  color: #1a202c;
  margin-right: 40px;
  letter-spacing: -0.3px;
  position: relative;
}

.logo span {
  color: #409EFF;
}

.menu-item {
  color: #4a5568;
  cursor: pointer;
  padding: 6px 0;
  font-weight: 500;
  position: relative;
  transition: color 0.2s ease;
}

.menu-item::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: #409EFF;
  transform: scaleX(0);
  transition: transform 0.2s ease;
  transform-origin: right;
}

.menu-item:hover {
  color: #409EFF;
}

.menu-item:hover::after {
  transform: scaleX(1);
  transform-origin: left;
}

.main-container {
  flex: 1;
  height: calc(100vh - 68px);
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
  gap: 24px;
  font-size: 14px;
  color: #4a5568;
  font-weight: 500;
}

.lang-btn {
  margin-left: 8px;
  padding: 4px 12px;
  font-weight: 500;
  transition: all 0.2s ease;
  background: transparent;
  color: #4a5568;
  border: none;
  position: relative;
}

.lang-btn::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: #409EFF;
  transform: scaleX(0);
  transition: transform 0.2s ease;
  transform-origin: right;
}

.lang-btn:hover {
  color: #409EFF;
}

.lang-btn:hover::after {
  transform: scaleX(1);
  transform-origin: left;
}

.content-body {
  padding: 32px;
  background-color: #f8fafc;
  flex-grow: 1;
  overflow-y: auto;
  max-width: 1600px;
  margin: 0 auto;
  width: 100%;
}

/* 自定义滚动条样式 */
.content-body::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.content-body::-webkit-scrollbar-track {
  background: transparent;
}

.content-body::-webkit-scrollbar-thumb {
  background: #e2e8f0;
  border-radius: 3px;
  transition: all 0.3s ease;
}

.content-body::-webkit-scrollbar-thumb:hover {
  background: #cbd5e1;
}
</style>

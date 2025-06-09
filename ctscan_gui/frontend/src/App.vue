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
      <!-- 侧边栏 -->
      <el-aside width="260px" class="sidebar">
        <div class="start-button-container">
          <el-button type="primary" size="large" class="start-button" :icon="DataAnalysis">开始分析</el-button>
        </div>
        <el-menu :default-active="activeMenu" class="sidebar-menu" background-color="transparent">
          <el-menu-item index="api" class="sidebar-menu-item">
            <el-icon><Setting /></el-icon>
            <span>设置 API</span>
          </el-menu-item>
          <div class="sidebar-group-title">分析模块</div>
          <el-menu-item index="overview" class="sidebar-menu-item analysis-overview-menu">
            <el-icon><DataAnalysis /></el-icon>
            <span>分析概览</span>
          </el-menu-item>
        </el-menu>
        <div class="os-info">OS: N/A<br>来源: 本地</div>
      </el-aside>
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
}
.app-menubar {
  height: 44px;
  background: #fff;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;
  font-size: 15px;
  user-select: none;
}
.menubar-left {
  display: flex;
  align-items: center;
  gap: 24px;
}
.logo {
  font-weight: bold;
  font-size: 18px;
  color: #409EFF;
  margin-right: 16px;
}
.menu-item {
  color: #303133;
  cursor: pointer;
  padding: 0 8px;
  border-radius: 4px;
  transition: background 0.2s;
}
.menu-item:hover {
  background: #f2f6fc;
}

.main-container {
  flex: 1;
  height: calc(100vh - 44px);
}

.sidebar {
  background-color: #fff;
  border-right: 1px solid #ebeef5;
  display: flex;
  flex-direction: column;
  padding: 0;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.03);
  position: relative;
}
.start-button-container {
  padding: 32px 24px 16px 24px;
}
.start-button {
  width: 100%;
  height: 48px;
  font-size: 18px;
  border-radius: 12px;
}
.sidebar-menu {
  border-right: none !important;
  margin-top: 8px;
  background: transparent;
}
.sidebar-menu-item {
  height: 44px;
  line-height: 44px;
  margin: 4px 0;
  border-radius: 8px;
  font-size: 16px;
}
.analysis-overview-menu.is-active {
  background-color: #ecf5ff !important;
  color: #409EFF !important;
  border-radius: 8px;
}
.sidebar-group-title {
  font-size: 13px;
  color: #909399;
  margin: 18px 0 4px 24px;
  font-weight: 500;
}
.os-info {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  padding: 18px 24px;
  font-size: 13px;
  color: #909399;
  border-top: 1px solid #ebeef5;
  background-color: #fafafa;
  line-height: 1.7;
}

.main-content-wrapper {
  padding: 0;
  display: flex;
  flex-direction: column;
  background-color: #f5f7fa;
  min-height: 0;
  padding-left: 0;
}
.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 15px;
  color: #303133;
}
.lang-btn {
  margin-left: 8px;
  border-radius: 8px;
  padding: 2px 14px;
}
.content-body {
  padding: 16px;
  background-color: #f5f7fa;
  flex-grow: 1;
  overflow-y: auto;
}
</style>

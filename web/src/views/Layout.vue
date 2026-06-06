<template>
  <el-container class="layout-container">
    <el-aside width="220px" class="sidebar">
      <div class="logo">
        <h2>HU MQ</h2>
        <p>企业级消息队列平台</p>
      </div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#001529"
        text-color="#ffffff99"
        active-text-color="#fff"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataAnalysis /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/clusters">
          <el-icon><Monitor /></el-icon>
          <span>集群管理</span>
        </el-menu-item>
        <el-menu-item index="/topics">
          <el-icon><Collection /></el-icon>
          <span>Topic管理</span>
        </el-menu-item>
        <el-menu-item index="/consumers">
          <el-icon><UserFilled /></el-icon>
          <span>消费组监控</span>
        </el-menu-item>
        <el-menu-item index="/messages">
          <el-icon><Search /></el-icon>
          <span>消息追踪</span>
        </el-menu-item>
        <el-menu-item index="/alerts">
          <el-icon><Bell /></el-icon>
          <span>告警中心</span>
        </el-menu-item>
        <el-menu-item index="/acl">
          <el-icon><Lock /></el-icon>
          <span>权限管理</span>
        </el-menu-item>
        <el-menu-item index="/ops">
          <el-icon><Setting /></el-icon>
          <span>运维工具</span>
        </el-menu-item>
        <el-menu-item index="/guide">
          <el-icon><Reading /></el-icon>
          <span>使用指引</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <div class="header-right">
          <span class="user-info">{{ username }}</span>
          <el-button type="danger" text @click="logout">退出登录</el-button>
        </div>
      </el-header>
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()
const activeMenu = computed(() => route.path)
const username = localStorage.getItem('username') || '管理员'

function logout() {
  localStorage.clear()
  router.push('/login')
}
</script>

<style scoped>
.layout-container { height: 100vh; }
.sidebar { background: #001529; overflow-y: auto; }
.logo { padding: 20px; text-align: center; color: #fff; border-bottom: 1px solid #ffffff1a; }
.logo h2 { margin: 0; font-size: 20px; }
.logo p { margin: 4px 0 0; font-size: 12px; color: #ffffff66; }
.header { background: #fff; display: flex; align-items: center; justify-content: flex-end; border-bottom: 1px solid #eee; padding: 0 24px; }
.header-right { display: flex; align-items: center; gap: 16px; }
.user-info { font-size: 14px; color: #333; }
.main-content { background: #f0f2f5; padding: 24px; }
</style>
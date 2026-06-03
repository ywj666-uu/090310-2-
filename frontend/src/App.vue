<template>
  <div id="app">
    <el-container v-if="isLoggedIn">
      <el-header class="app-header">
        <div class="header-left">
          <h1>🌍 碳足迹追踪</h1>
        </div>
        <div class="header-right">
          <span class="username">{{ username }}</span>
          <el-button type="danger" size="small" @click="logout">退出</el-button>
        </div>
      </el-header>
      <el-container>
        <el-aside width="200px">
          <el-menu :default-active="activeMenu" router>
            <el-menu-item index="/dashboard">
              <span>📊 数据概览</span>
            </el-menu-item>
            <el-menu-item index="/record">
              <span>📝 记录碳排</span>
            </el-menu-item>
            <el-menu-item index="/trend">
              <span>📈 趋势分析</span>
            </el-menu-item>
            <el-menu-item index="/goal">
              <span>🎯 减排目标</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        <el-main>
          <router-view />
        </el-main>
      </el-container>
    </el-container>
    <router-view v-else />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const isLoggedIn = computed(() => !!localStorage.getItem('token'))
const username = computed(() => localStorage.getItem('username') || '')
const activeMenu = computed(() => route.path)

function logout() {
  localStorage.removeItem('token')
  localStorage.removeItem('username')
  router.push('/login')
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Microsoft YaHei', sans-serif;
  background: #f0f2f5;
}

.app-header {
  background: linear-gradient(135deg, #2d8c4e, #1a6b3a);
  color: white;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.app-header h1 {
  font-size: 20px;
  font-weight: 600;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.username {
  font-size: 14px;
}

.el-aside {
  background: white;
  min-height: calc(100vh - 60px);
  box-shadow: 2px 0 6px rgba(0, 0, 0, 0.05);
}

.el-main {
  padding: 20px;
  min-height: calc(100vh - 60px);
}
</style>

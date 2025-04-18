<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { RouterView } from 'vue-router'
import Header from './components/Header.vue'
import Footer from './components/Footer.vue'
import TokenExpiryAlert from './components/TokenExpiryAlert.vue'
import { setupTokenExpiryMonitor } from './utils/tokenExpiry'

// 控制Token过期提示框的显示
const showTokenAlert = ref(false)
const minutesLeft = ref(3)
let cleanupTokenMonitor = null

// 处理Token续期成功
const handleRefreshSuccess = () => {
  showTokenAlert.value = false
}

// 处理关闭提示框
const handleCloseAlert = () => {
  showTokenAlert.value = false
}

onMounted(() => {
  // 设置Token过期监控，3分钟前提醒
  cleanupTokenMonitor = setupTokenExpiryMonitor({
    warningThreshold: 3, // 3分钟前提醒
    autoRefresh: false, // 不自动刷新，由用户决定是否刷新
    onExpiringSoon: () => {
      showTokenAlert.value = true
    }
  })
})

onUnmounted(() => {
  // 清理Token监控
  if (cleanupTokenMonitor) {
    cleanupTokenMonitor()
  }
})
</script>

<template>
  <div class="app">
    <Header />

    <main style="min-height: calc(80vh)">
      <div class="container-fluid px-0">
        <RouterView />
      </div>
    </main>

    <Footer />
    
    <!-- Token过期提示框 -->
    <TokenExpiryAlert 
      :visible="showTokenAlert" 
      :minutes-left="minutesLeft"
      @close="handleCloseAlert"
      @refresh-success="handleRefreshSuccess"
    />
  </div>
</template>

<style>
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

main {
  flex: 1;
}
</style>

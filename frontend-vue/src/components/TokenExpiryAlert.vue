<script setup>
import { ref } from 'vue'
import { refreshToken } from '../services/auth'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  minutesLeft: {
    type: Number,
    default: 3
  }
})

const emit = defineEmits(['close', 'refresh-success'])
const loading = ref(false)
const error = ref('')

const handleRefresh = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const newToken = await refreshToken()
    if (newToken) {
      emit('refresh-success')
    } else {
      error.value = '续期失败，请重新登录'
    }
  } catch (err) {
    error.value = '续期失败，请重新登录'
    console.error('Token refresh error:', err)
  } finally {
    loading.value = false
  }
}

const handleClose = () => {
  emit('close')
}
</script>

<template>
  <div v-if="visible" class="token-expiry-alert">
    <div class="alert-container">
      <div class="alert-header">
        <h4>登录状态即将过期</h4>
        <button class="close-btn" @click="handleClose">&times;</button>
      </div>
      <div class="alert-body">
        <p>您的登录状态将在 {{ minutesLeft }} 分钟后过期，请点击续期按钮继续使用。</p>
        <div v-if="error" class="error-message">{{ error }}</div>
      </div>
      <div class="alert-footer">
        <button 
          class="refresh-btn" 
          @click="handleRefresh" 
          :disabled="loading"
        >
          {{ loading ? '续期中...' : '立即续期' }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.token-expiry-alert {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 9999;
}

.alert-container {
  background-color: white;
  border-radius: 8px;
  width: 400px;
  max-width: 90%;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  overflow: hidden;
}

.alert-header {
  padding: 16px 20px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.alert-header h4 {
  margin: 0;
  font-size: 18px;
  color: #343a40;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #6c757d;
}

.alert-body {
  padding: 20px;
}

.error-message {
  color: #dc3545;
  margin-top: 10px;
  font-size: 14px;
}

.alert-footer {
  padding: 16px 20px;
  border-top: 1px solid #e9ecef;
  display: flex;
  justify-content: flex-end;
}

.refresh-btn {
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 8px 16px;
  cursor: pointer;
  font-size: 14px;
}

.refresh-btn:hover {
  background-color: #0069d9;
}

.refresh-btn:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}
</style>
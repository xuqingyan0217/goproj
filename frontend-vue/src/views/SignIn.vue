<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { login } from '../services/auth'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const email = ref('')
const password = ref('')
const error = ref('')

const handleSubmit = async () => {
  try {
    const response = await login(email.value, password.value, route.query.next?.toString() || '/')
    if (response.code === -1) {
      error.value = response.error || '登录失败'
      return
    }
    userStore.setUserId(response.user_id)
    if (typeof response.cart_num === 'number') {
      userStore.setCartNum(response.cart_num)
    }
    router.push(route.query.next?.toString() || '/')
  } catch (err) {
    error.value = '登录失败，请稍后重试'+ err
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
      <h2 class="login-title">欢迎登录</h2>
      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label for="email" class="form-label">邮箱 <span class="text-danger">*</span></label>
          <input
            type="email"
            class="form-control"
            id="email"
            v-model="email"
            required
            placeholder="请输入邮箱"
          />
        </div>
        <div class="form-group">
          <label for="password" class="form-label">密码 <span class="text-danger">*</span></label>
          <input
            type="password"
            class="form-control"
            id="password"
            v-model="password"
            required
            placeholder="请输入密码"
          />
        </div>
        <div class="form-group text-center">
          还没有账号？ <router-link to="/sign-up" class="signup-link">立即注册</router-link>
        </div>
        <div v-if="error" class="alert alert-danger" role="alert">
          {{ error }}
        </div>
        <button type="submit" class="login-btn">登录</button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 20px;
  margin-top: -120px;
}

.login-card {
  background: white;
  padding: 2rem;
  border-radius: 15px;
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
  transition: transform 0.3s ease;
}

.login-card:hover {
  transform: translateY(-5px);
}

.login-title {
  text-align: center;
  color: #2c3e50;
  margin-bottom: 2rem;
  font-weight: 600;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-control {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 0.8rem;
  width: 100%;
  transition: border-color 0.3s ease;
}

.form-control:focus {
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.2);
  outline: none;
}

.login-btn {
  width: 100%;
  padding: 0.8rem;
  background: #3498db;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.login-btn:hover {
  background: #2980b9;
}

.signup-link {
  color: #3498db;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.3s ease;
}

.signup-link:hover {
  color: #2980b9;
}

.alert {
  margin-bottom: 1rem;
  padding: 0.8rem;
  border-radius: 8px;
  font-size: 0.9rem;
}

.form-label {
  color: #2c3e50;
  font-weight: 500;
  margin-bottom: 0.5rem;
  display: block;
}
</style>
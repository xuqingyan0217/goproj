<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { register } from '../services/auth'

const router = useRouter()
const email = ref('')
const password = ref('')
const passwordConfirm = ref('')
const error = ref('')

// 验证规则
const emailError = computed(() => {
  if (!email.value) return ''
  if (email.value.length < 5) return '邮箱长度不能少于5个字符'
  if (email.value.length > 100) return '邮箱长度不能超过100个字符'
  const emailRegex = /^[\w-]+(\.[\w-]+)*@[\w-]+(\.[\w-]+)+$/
  if (!emailRegex.test(email.value)) return '请输入有效的邮箱地址'
  return ''
})

const passwordError = computed(() => {
  if (!password.value) return ''
  if (password.value.length < 8) return '密码长度至少为8个字符'
  if (password.value.length > 20) return '密码长度不能超过20个字符'
  const passwordRegex = /^[a-zA-Z\d]*[a-z]+[a-zA-Z\d]*[A-Z]+[a-zA-Z\d]*\d+[a-zA-Z\d]*$|^[a-zA-Z\d]*[a-z]+[a-zA-Z\d]*\d+[a-zA-Z\d]*[A-Z]+[a-zA-Z\d]*$|^[a-zA-Z\d]*[A-Z]+[a-zA-Z\d]*[a-z]+[a-zA-Z\d]*\d+[a-zA-Z\d]*$|^[a-zA-Z\d]*[A-Z]+[a-zA-Z\d]*\d+[a-zA-Z\d]*[a-z]+[a-zA-Z\d]*$|^[a-zA-Z\d]*\d+[a-zA-Z\d]*[a-z]+[a-zA-Z\d]*[A-Z]+[a-zA-Z\d]*$|^[a-zA-Z\d]*\d+[a-zA-Z\d]*[A-Z]+[a-zA-Z\d]*[a-z]+[a-zA-Z\d]*$/
  if (!passwordRegex.test(password.value)) return '密码必须包含大小写字母和数字'
  return ''
})

const passwordConfirmError = computed(() => {
  if (!passwordConfirm.value) return ''
  if (password.value !== passwordConfirm.value) return '两次输入的密码不一致'
  return ''
})

const isFormValid = computed(() => {
  return !emailError.value && !passwordError.value && !passwordConfirmError.value &&
         email.value && password.value && passwordConfirm.value
})

const handleSubmit = async () => {
  if (!isFormValid.value) {
    error.value = '请修正表单中的错误'
    return
  }

  try {
    const response = await register(email.value, password.value, passwordConfirm.value)
    if (response.code === -1) {
      error.value = response.error || '注册失败，请稍后重试'
      return
    }
    router.push('/sign-in')
  } catch (err) {
    error.value = '注册失败，请稍后重试'
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
      <h2 class="login-title">注册账号</h2>
      <p class="login-subtitle">创建您的账户，开始购物之旅</p>
      <form @submit.prevent="handleSubmit" class="login-form">
        <div class="form-group">
          <label for="email" class="form-label">邮箱 <span class="text-danger">*</span></label>
          <input
            type="email"
            id="email"
            v-model="email"
            class="form-control"
            :class="{ 'error': emailError }"
            required
            placeholder="请输入邮箱地址"
          />
          <p v-if="emailError" class="error-message">{{ emailError }}</p>
        </div>
        <div class="form-group">
          <label for="password" class="form-label">密码 <span class="text-danger">*</span></label>
          <input
            type="password"
            id="password"
            v-model="password"
            class="form-control"
            :class="{ 'error': passwordError }"
            required
            placeholder="请输入密码"
          />
          <p v-if="passwordError" class="error-message">{{ passwordError }}</p>
        </div>
        <div class="form-group">
          <label for="password-confirm" class="form-label">确认密码 <span class="text-danger">*</span></label>
          <input
            type="password"
            id="password-confirm"
            v-model="passwordConfirm"
            class="form-control"
            :class="{ 'error': passwordConfirmError }"
            required
            placeholder="请再次输入密码"
          />
          <p v-if="passwordConfirmError" class="error-message">{{ passwordConfirmError }}</p>
        </div>

        <div v-if="error" class="alert-error">
          <svg class="alert-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
          <span>{{ error }}</span>
        </div>

        <div class="form-footer">
          <router-link to="/sign-in" class="login-link">已有账号？立即登录</router-link>
        </div>

        <button
          type="submit"
          :disabled="!isFormValid"
          class="submit-button"
        >
          注册
        </button>
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
  margin-top: -80px;
}

.login-card {
  background: white;
  padding: 2.5rem;
  border-radius: 20px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 450px;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.login-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.15);
}

.login-title {
  font-size: 2rem;
  font-weight: 700;
  color: #2c3e50;
  text-align: center;
  margin-bottom: 0.5rem;
}

.login-subtitle {
  text-align: center;
  color: #6b7280;
  margin-bottom: 2rem;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-label {
  font-size: 0.9rem;
  font-weight: 500;
  color: #4b5563;
}

.form-control {
  padding: 0.75rem 1rem;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  font-size: 1rem;
  transition: all 0.3s ease;
  width: 100%;
  outline: none;
}

.form-control:focus {
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.form-control.error {
  border-color: #ef4444;
}

.error-message {
  font-size: 0.85rem;
  color: #ef4444;
  margin-top: 0.25rem;
}

.alert-error {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background-color: #fee2e2;
  border-radius: 10px;
  color: #dc2626;
}

.alert-icon {
  width: 1.25rem;
  height: 1.25rem;
  flex-shrink: 0;
}

.form-footer {
  text-align: center;
  margin-top: 1rem;
}

.login-link {
  color: #6366f1;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.3s ease;
}

.login-link:hover {
  color: #4f46e5;
  text-decoration: underline;
}

.submit-button {
  background-color: #6366f1;
  color: white;
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 10px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  width: 100%;
  margin-top: 1rem;
}

.submit-button:hover:not(:disabled) {
  background-color: #4f46e5;
  transform: translateY(-1px);
}

.submit-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

@media (max-width: 640px) {
  .login-card {
    padding: 1.5rem;
  }

  .login-title {
    font-size: 1.75rem;
  }
}
</style>
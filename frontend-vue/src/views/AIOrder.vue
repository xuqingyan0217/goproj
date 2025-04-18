<template>
  <div class="container">
    <h1>AI 预下单</h1>
    <form @submit.prevent="submitOrder" class="order-form">
      <textarea
        v-model="orderInput"
        placeholder="请输入您的订单需求,eg:选择商品名称带有t-shirt的商品，每个购买1件"
        required
      ></textarea>
      <button type="submit" :disabled="isSubmitting">
        <span class="spinner" v-if="isSubmitting"></span>
        {{ isSubmitting ? '提交中...' : '提交订单' }}
      </button>
    </form>
    <div v-if="orderResult.show" :class="['result', orderResult.success ? 'success' : 'error']">
      <p>{{ orderResult.message }}</p>
    </div>

    <h1>AI 获取订单信息</h1>
    <form @submit.prevent="getOrderInfo" class="order-form">
      <input
        v-model="orderListInput"
        type="text"
        placeholder="请输入您的需求,eg:查询出我的的订单信息"
        required
      >
      <button type="submit" :disabled="isGettingInfo">
        <span class="spinner" v-if="isGettingInfo"></span>
        {{ isGettingInfo ? '获取中...' : '获取订单信息' }}
      </button>
    </form>
    <div v-if="orderListResult.show" :class="['result', orderListResult.success ? 'success' : 'error']">
      <p>{{ orderListResult.message }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { submitAIOrder, getAIOrderInfo } from '../services/aieino'
import { useUserStore } from '../stores/user'

const store = useUserStore()
const orderInput = ref('')
const orderListInput = ref('')
const isSubmitting = ref(false)
const isGettingInfo = ref(false)
const orderResult = ref({
  show: false,
  success: false,
  message: ''
})
const orderListResult = ref({
  show: false,
  success: false,
  message: ''
})

const submitOrder = async () => {
  isSubmitting.value = true
  try {
    const response = await submitAIOrder(orderInput.value)
    
    orderResult.value = {
      show: true,
      success: true,
      message: `预下单成功！订单信息：${response.data.orderId}`
    }
    
    // 重置表单
    orderInput.value = ''
  } catch (error) {
    orderResult.value = {
      show: true,
      success: false,
      message: `预下单失败：${error.response?.data?.error || '未知错误'}`
    }
  } finally {
    isSubmitting.value = false
  }
}

const getOrderInfo = async () => {
  isGettingInfo.value = true
  try {
    const response = await getAIOrderInfo(orderListInput.value, store.userId)
    
    orderListResult.value = {
      show: true,
      success: true,
      message: `查询该用户订单信息：${JSON.stringify(response.data.orderInfo)}`
    }
    
    // 重置表单
    orderListInput.value = ''
  } catch (error) {
    console.error('获取订单信息错误:', error)
    orderListResult.value = {
      show: true,
      success: false,
      message: `获取订单信息失败：${error.response?.data?.message || error.response?.data?.error || error.message || '未知错误'}`
    }
  } finally {
    isGettingInfo.value = false
  }
}
</script>

<style scoped>
.container {
  max-width: 500px;
  margin: 0 auto;
  padding: 20px;
  font-family: Arial, sans-serif;
}

.order-form {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 20px;
}

textarea,
input {
  padding: 10px;
  font-size: 16px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

textarea {
  height: 100px;
  resize: vertical;
}

button {
  padding: 10px 20px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.2s;
}

button:hover {
  background-color: #0056b3;
}

.result {
  margin-top: 10px;
  padding: 15px;
  border-radius: 4px;
  font-size: 14px;
}

.success {
  background-color: #e8f5e9;
  color: #2e7d32;
  border: 1px solid #c8e6c9;
}

.error {
  background-color: #ffebee;
  color: #c62828;
  border: 1px solid #ffcdd2;
}

h1 {
  font-size: 24px;
  margin-bottom: 20px;
  color: #333;
}

.spinner {
  display: inline-block;
  width: 16px;
  height: 16px;
  margin-right: 8px;
  border: 2px solid #ffffff;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}
</style>
<template>
  <div class="container row p-5 d-flex justify-content-center">
    <div class="text-danger h3 text-center mb-5">
      {{ status === 'cancel' ? '正在取消支付，请稍候...' : '正在处理支付，请稍候...' }}
    </div>
    <div class="spinner-grow text-primary" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
    <div class="spinner-grow text-secondary" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
    <div class="spinner-grow text-success" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
    <div class="spinner-grow text-danger" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
    <div class="spinner-grow text-warning" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
    <div class="spinner-grow text-info" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
    <div class="spinner-grow text-light" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
    <div class="spinner-grow text-dark" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getCheckoutResult } from '../services/checkout'

const router = useRouter()
const route = useRoute()
const status = ref(route.query.status || 'success')

const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms))

onMounted(async () => {
  try {
    await sleep(3000) // 等待3秒
    await getCheckoutResult()
    router.push(`/checkout/result?status=${status.value}`)
  } catch (error) {
    console.error('Failed to check status:', error)
  }
})
</script>
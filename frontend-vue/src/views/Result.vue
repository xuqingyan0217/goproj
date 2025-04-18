<template>
  <div class="container row p-5 d-flex justify-content-center">
    <template v-if="status === 'success'">
      <i class="fa-regular fa-circle-check fs-1 text-success"></i>
      <div class="text-center fs-3">
        恭喜您，订单已成功支付！
      </div>
    </template>
    <template v-else>
      <i class="fa-regular fa-circle-xmark fs-1 text-warning"></i>
      <div class="text-center fs-3">
        您已取消支付，可以稍后在订单中心继续处理订单。
      </div>
    </template>
  </div>
  <div class="d-flex justify-content-center">
    <router-link to="/order" class="btn btn-info">查看订单</router-link>
    <router-link to="/" class="btn btn-success ms-5">返回首页</router-link>
  </div>
</template>

<script setup>
import { useUserStore } from '../stores/user'
import { useRoute } from 'vue-router'
import { onMounted, ref } from 'vue'

const store = useUserStore()
const route = useRoute()
const status = ref(route.query.status || 'success')

onMounted(() => {
  // 如果是支付成功，清空购物车数量
  if (status.value === 'success') {
    store.setCartNum(0)
  }
})
</script>
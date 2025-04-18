<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCartItems, removeFromCart, updateCartItemQuantity } from '../services/cart'

const router = useRouter()
const items = ref([])
const total = ref(0)
const loading = ref(true)
const error = ref('')

const fetchCartItems = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await getCartItems()
    items.value = Array.isArray(response.items) ? response.items : []
    total.value = response.total
  } catch (err) {
    console.error('Failed to fetch cart items:', err)
    error.value = '获取购物车数据失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

const handleRemoveItem = async (productId) => {
  try {
    await removeFromCart(productId)
    await fetchCartItems()
  } catch (err) {
    console.error('Failed to remove item:', err)
  }
}

const handleUpdateQuantity = async (productId, quantity) => {
  try {
    await updateCartItemQuantity(productId, quantity)
    await fetchCartItems()
  } catch (err) {
    console.error('Failed to update quantity:', err)
  }
}

const goToCheckout = () => {
  router.push('/checkout')
}

onMounted(() => {
  fetchCartItems()
})
</script>

<template>
  <div class="container py-5">
    <div class="row justify-content-center">
      <div class="col-12 col-lg-10">
        <h2 class="mb-4 text-center fw-bold">我的购物车</h2>
        <div v-if="loading" class="text-center py-5">
          <div class="spinner-border text-primary" role="status">
            <span class="visually-hidden">加载中...</span>
          </div>
        </div>
        <div v-else-if="error" class="alert alert-danger shadow-sm" role="alert">
          <i class="fas fa-exclamation-circle me-2"></i>{{ error }}
        </div>
        <div v-else-if="items.length > 0" class="card shadow-sm border-0">
          <ul class="list-group list-group-flush">
            <li v-for="item in items" :key="item.id" class="list-group-item py-4 border-bottom">
              <div class="row align-items-center g-4">
                <div class="col-12 col-md-2">
                  <div class="rounded overflow-hidden shadow-sm">
                    <img :src="item.Picture" class="img-fluid w-100" style="object-fit: cover; height: 120px" :alt="item.Name">
                  </div>
                </div>
                <div class="col-12 col-md-4">
                  <h5 class="mb-2 fw-bold">{{ item.Name }}</h5>
                  <p class="mb-1 text-muted fs-5">¥{{ item.Price }}</p>
                </div>
                <div class="col-12 col-md-3">
                  <div class="input-group input-group-lg shadow-sm">
                    <button class="btn btn-outline-primary" 
                            @click="handleUpdateQuantity(item.id, item.Qty - 1)" 
                            :disabled="item.Qty <= 1">
                      <i class="fas fa-minus"></i>
                    </button>
                    <input type="number" class="form-control text-center border-primary" 
                           v-model="item.Qty" readonly>
                    <button class="btn btn-outline-primary" 
                            @click="handleUpdateQuantity(item.id, item.Qty + 1)">
                      <i class="fas fa-plus"></i>
                    </button>
                  </div>
                </div>
                <div class="col-8 col-md-2">
                  <p class="mb-0 fs-5 fw-bold text-primary text-end">¥{{ (item.Price * item.Qty).toFixed(2) }}</p>
                </div>
                <div class="col-4 col-md-1 text-end">
                  <button class="btn btn-outline-danger rounded-circle" 
                          @click="handleRemoveItem(item.id)">
                    <i class="fas fa-trash"></i>
                  </button>
                </div>
              </div>
            </li>
          </ul>
          <div class="card-footer bg-white py-4 border-0">
            <div class="d-flex flex-column flex-md-row justify-content-between align-items-center gap-3">
              <div class="d-flex align-items-baseline">
                <span class="text-muted me-2 fs-5">总计:</span>
                <h3 class="mb-0 text-primary">¥{{ Number(total).toFixed(2) }}</h3>
              </div>
              <button class="btn btn-primary btn-lg px-5" @click="goToCheckout">
                <i class="fas fa-shopping-cart me-2"></i>去结算
              </button>
            </div>
          </div>
        </div>
        <div v-else class="text-center py-5">
          <div class="mb-4">
            <i class="fas fa-shopping-cart fa-4x text-muted"></i>
          </div>
          <h3 class="text-muted mb-4">购物车是空的</h3>
          <router-link to="/" class="btn btn-primary btn-lg">
            <i class="fas fa-store me-2"></i>去购物
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>
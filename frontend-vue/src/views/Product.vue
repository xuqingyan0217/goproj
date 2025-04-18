<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { addToCart as addToCartService } from '../services/cart'
import { getProduct } from '../services/product'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const item = ref(null)
const quantity = ref(1)
const addingToCart = ref(false)
const message = ref('')

const fetchProduct = async () => {
  try {
    item.value = await getProduct(route.query.id)
  } catch (err) {
    console.error('Failed to fetch product:', err)
  }
}

const addToCart = async () => {
  if (!userStore.userId) {
    router.push(`/sign-in?next=/product?id=${route.query.id}`)
    return
  }

  addingToCart.value = true
  message.value = ''

  try {
    const result = await addToCartService(item.value.id, quantity.value)
    message.value = '成功添加到购物车'
    // 直接使用返回的购物车数量更新store
    userStore.setCartNum(result.cartNum)
  } catch (err) {
    console.error('Error:', err)
    message.value = '添加失败，请重试'
  } finally {
    addingToCart.value = false
  }
}

onMounted(() => {
  fetchProduct()
})
</script>

<template>
  <div class="container py-5" v-if="item">
    <div class="row">
      <div class="col-md-6">
        <div id="productPicture" class="carousel slide">
          <div class="carousel-indicators">
            <button type="button" data-bs-target="#productPicture" data-bs-slide-to="0" class="active" aria-current="true" aria-label="Slide 1"></button>
            <button type="button" data-bs-target="#productPicture" data-bs-slide-to="1" aria-label="Slide 2"></button>
            <button type="button" data-bs-target="#productPicture" data-bs-slide-to="2" aria-label="Slide 3"></button>
          </div>
          <div class="carousel-inner">
            <div class="carousel-item active">
              <img :src="item.picture" class="d-block w-100 rounded" :alt="item.name">
            </div>
            <div class="carousel-item">
              <img :src="item.picture" class="d-block w-100 rounded" :alt="item.name">
            </div>
            <div class="carousel-item">
              <img :src="item.picture" class="d-block w-100 rounded" :alt="item.name">
            </div>
          </div>
          <button class="carousel-control-prev" type="button" data-bs-target="#productPicture" data-bs-slide="prev">
            <span class="carousel-control-prev-icon" aria-hidden="true"></span>
            <span class="visually-hidden">Previous</span>
          </button>
          <button class="carousel-control-next" type="button" data-bs-target="#productPicture" data-bs-slide="next">
            <span class="carousel-control-next-icon" aria-hidden="true"></span>
            <span class="visually-hidden">Next</span>
          </button>
        </div>
      </div>
      <div class="col-md-6">
        <h1 class="mb-4">{{ item.name }}</h1>
        <p class="text-muted mb-4">{{ item.description }}</p>
        <div class="d-flex align-items-center mb-4">
          <h2 class="text-primary mb-0">￥{{ item.price }}</h2>
        </div>
        <div class="mb-4">
          <label class="form-label">数量</label>
          <div class="input-group" style="width: 150px;">
            <button class="btn btn-outline-secondary" type="button" @click="quantity = Math.max(1, quantity - 1)">-</button>
            <input type="number" class="form-control text-center" v-model="quantity" min="1">
            <button class="btn btn-outline-secondary" type="button" @click="quantity++">+</button>
          </div>
        </div>
        <button class="btn btn-primary" @click="addToCart" :disabled="addingToCart">
          <span class="spinner-border spinner-border-sm me-2" v-if="addingToCart"></span>
          {{ addingToCart ? '添加中...' : '加入购物车' }}
        </button>
        <router-link to="/cart" class="btn btn-outline-primary ms-2">
          <i class="fas fa-shopping-cart me-1"></i>查看购物车
        </router-link>
        <div class="alert" :class="message.includes('成功') ? 'alert-success' : 'alert-danger'" role="alert" v-if="message">
          {{ message }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.carousel img {
  width: 100%;
  max-width: 600px;
  height: auto;
  max-height: 500px;
  object-fit: contain;
  margin: 0 auto;
}

.carousel-control-prev,
.carousel-control-next {
  width: 5%;
  background-color: rgba(0, 0, 0, 0.3);
  border-radius: 0 3px 3px 0;
  z-index: 10;
}

.carousel-control-next {
  border-radius: 3px 0 0 3px;
}

.carousel-control-prev:hover,
.carousel-control-next:hover {
  background-color: rgba(0, 0, 0, 0.5);
}

.carousel-control-prev-icon,
.carousel-control-next-icon {
  width: 2rem;
  height: 2rem;
}
</style>
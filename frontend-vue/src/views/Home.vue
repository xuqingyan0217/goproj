<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '../stores/user'
import request from '../utils/request'

const store = useUserStore()
const items = ref([])

const fetchHomeData = async () => {
  try {
    const response = await request.get('/api')
    const data = response.data
    items.value = data.items || []
    
    // 更新用户状态到store
    if (data.user_id) {
      store.setUserId(data.user_id)
    }
    if (typeof data.cart_num === 'number') {
      store.setCartNum(data.cart_num)
    }
  } catch (err) {
    console.error('Failed to fetch home data:', err)
    store.setError('Failed to load page data')
  }
}

onMounted(() => {
  fetchHomeData()
})
</script>

<template>
  <div class="row g-3 mx-0">
    <div
      v-for="item in items"
      :key="item.id"
      class="col-xl-3 col-lg-4 col-md-6 col-sm-12 px-2"
    >
      <router-link :to="`/product?id=${item.id}`" class="text-decoration-none">
        <div class="card h-100 border-0 shadow-sm product-card">
          <div class="card-img-container">
            <img :src="item.picture" class="card-img-top" :alt="item.name" />
          </div>
          <div class="card-body d-flex flex-column">
            <h5 class="card-title text-dark mb-2 product-name">{{ item.name }}</h5>
            <p class="card-text text-muted mb-3 flex-grow-1">{{ item.description || '暂无描述' }}</p>
            <div class="d-flex justify-content-between align-items-center">
              <h4 class="card-title text-primary mb-0">￥{{ item.price }}</h4>
              <span class="badge bg-primary">{{ item.categories || '未分类' }}</span>
            </div>
          </div>
        </div>
      </router-link>
    </div>
  </div>
</template>

<style scoped>
.product-card {
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  border-radius: 12px;
  overflow: hidden;
}

.product-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
}

.card-img-container {
  height: 280px;
  overflow: hidden;
}

.card img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.product-card:hover img {
  transform: scale(1.05);
}

.product-name {
  font-size: 1.1rem;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.card-text {
  font-size: 0.9rem;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}
</style>
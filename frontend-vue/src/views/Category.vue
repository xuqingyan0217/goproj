<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const products = ref([])
const loading = ref(true)
const error = ref(null)

const fetchProducts = async () => {
  const category = route.params.category
  loading.value = true
  error.value = null
  try {
    const response = await fetch(`/category/${category}`)
    if (response.ok) {
      const data = await response.json()
      products.value = data.items || []
    } else {
      error.value = '获取商品列表失败'
    }
  } catch (err) {
    error.value = '服务器错误'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchProducts()
})

watch(() => route.params.category, () => {
  fetchProducts()
})
</script>

<template>
  <div class="container mt-4">
    <h2 class="text-center mb-4">{{ route.params.category }} 分类</h2>

    <div v-if="loading" class="text-center">
      <div class="spinner-border" role="status">
        <span class="visually-hidden">加载中...</span>
      </div>
    </div>

    <div v-else-if="error" class="alert alert-danger" role="alert">
      {{ error }}
    </div>

    <div v-else class="row row-cols-1 row-cols-md-3 g-4">
      <div v-for="product in products" :key="product.id" class="col">
        <div class="card h-100">
          <img :src="product.picture" class="card-img-top" :alt="product.name">
          <div class="card-body">
            <h5 class="card-title">{{ product.name }}</h5>
            <p class="card-text">{{ product.description }}</p>
            <p class="card-text"><strong>价格: ¥{{ product.price }}</strong></p>
            <router-link :to="`/product?id=${product.id}`" class="btn btn-primary">查看详情</router-link>
          </div>
        </div>
      </div>

      <div v-if="products.length === 0" class="col-12 text-center">
        <p>该分类下暂无商品</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card-img-top {
  height: 200px;
  object-fit: cover;
}

.card {
  transition: transform 0.2s;
}

.card:hover {
  transform: translateY(-5px);
}
</style>
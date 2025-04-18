<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { ref, onMounted, watch } from 'vue'
import { getAdminStats } from '../services/admin'

const router = useRouter()
const store = useUserStore()
const stats = ref({
  totalProducts: 0,
  totalCategories: 0,
  totalOrders: 0,
  recentProducts: []
})

// 图片预览相关的响应式数据
const showPreview = ref(false)
const previewImage = ref('')
const previewScale = ref(1)

// 打开图片预览
const openPreview = (imageUrl) => {
  previewImage.value = imageUrl
  showPreview.value = true
  previewScale.value = 1
}

// 关闭图片预览
const closePreview = () => {
  showPreview.value = false
  previewScale.value = 1
}

// 缩放图片
const zoomImage = (factor) => {
  previewScale.value *= factor
}

const fetchStats = async () => {
  try {
    const data = await getAdminStats()
    stats.value = data
  } catch (err) {
    console.error('Failed to fetch stats:', err)
  }
}

onMounted(() => {
  fetchStats()
})

watch(() => router.currentRoute.value.path, () => {
  if (router.currentRoute.value.path === '/admin') {
    fetchStats()
  }
})
</script>

<template>
  <div class="container-fluid">
    <nav class="navbar navbar-expand-lg bg-body-tertiary mb-3">
      <div class="container-fluid">
        <a class="navbar-brand" href="/admin">
          <i class="fas fa-cogs me-2"></i>商品管理
        </a>
        <button
          class="navbar-toggler"
          type="button"
          data-bs-toggle="collapse"
          data-bs-target="#navbarNav"
          aria-controls="navbarNav"
          aria-expanded="false"
          aria-label="Toggle navigation"
        >
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarNav">
          <ul class="navbar-nav">
            <li class="nav-item">
              <router-link class="nav-link" to="/admin">
                <i class="fas fa-home me-1"></i>首页
              </router-link>
            </li>
            <li class="nav-item">
              <router-link class="nav-link" to="/admin/product/create">
                <i class="fas fa-plus-circle me-1"></i>创建商品
              </router-link>
            </li>
            <li class="nav-item">
              <router-link class="nav-link" to="/admin/product/delete">
                <i class="fas fa-trash me-1"></i>删除商品
              </router-link>
            </li>
            <li class="nav-item">
              <router-link class="nav-link" to="/admin/product/update">
                <i class="fas fa-edit me-1"></i>更新商品
              </router-link>
            </li>
          </ul>
        </div>
      </div>
    </nav>

    <!-- 数据统计卡片 -->
    <div class="container-fluid mb-4" v-if="$route.path === '/admin'">
      <div class="row g-4">
        <div class="col-xl-3 col-lg-6">
          <div class="card border-0 shadow-sm h-100 bg-primary bg-gradient text-white">
            <div class="card-body">
              <div class="d-flex justify-content-between align-items-center">
                <div>
                  <h6 class="card-title mb-0">总商品数</h6>
                  <h2 class="mt-2 mb-0">{{ stats.totalProducts }}</h2>
                </div>
                <i class="fas fa-box fa-2x opacity-50"></i>
              </div>
            </div>
          </div>
        </div>

        <div class="col-xl-3 col-lg-6">
          <div class="card border-0 shadow-sm h-100 bg-success bg-gradient text-white">
            <div class="card-body">
              <div class="d-flex justify-content-between align-items-center">
                <div>
                  <h6 class="card-title mb-0">商品分类数</h6>
                  <h2 class="mt-2 mb-0">{{ stats.totalCategories }}</h2>
                </div>
                <i class="fas fa-tags fa-2x opacity-50"></i>
              </div>
            </div>
          </div>
        </div>

        <div class="col-xl-3 col-lg-6">
          <div class="card border-0 shadow-sm h-100 bg-info bg-gradient text-white">
            <div class="card-body">
              <div class="d-flex justify-content-between align-items-center">
                <div>
                  <h6 class="card-title mb-0">订单总数</h6>
                  <h2 class="mt-2 mb-0">{{ stats.totalOrders }}</h2>
                </div>
                <i class="fas fa-shopping-cart fa-2x opacity-50"></i>
              </div>
            </div>
          </div>
        </div>

        <div class="col-xl-3 col-lg-6">
          <div class="card border-0 shadow-sm h-100">
            <div class="card-body">
              <h6 class="card-title">最新商品</h6>
              <ul class="list-unstyled mb-0">
                <li v-for="product in stats.recentProducts" :key="product.id" class="mb-2 d-flex align-items-center">
                  <img 
                    v-if="product.picture" 
                    :src="product.picture" 
                    class="me-2 product-thumbnail" 
                    :alt="product.name"
                  />
                  <small class="text-muted">{{ product.name }}</small>
                </li>
              </ul>
            </div>
          </div>
        </div>

        <!-- 图片预览模态框 -->
        <div v-if="showPreview" class="image-preview-modal" @click.self="closePreview">
          <div class="preview-content">
            <img 
              :src="previewImage" 
              :style="{ transform: `scale(${previewScale})` }" 
              class="preview-image"
            />
            <div class="preview-controls">
              <button class="btn btn-light me-2" @click="zoomImage(1.2)">
                <i class="fas fa-search-plus"></i>
              </button>
              <button class="btn btn-light me-2" @click="zoomImage(0.8)">
                <i class="fas fa-search-minus"></i>
              </button>
              <button class="btn btn-light" @click="closePreview">
                <i class="fas fa-times"></i>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="container-fluid py-3">
      <router-view></router-view>
    </div>
  </div>
</template>

<style scoped>
.navbar {
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.card {
  transition: transform 0.2s;
}

.card:hover {
  transform: translateY(-5px);
}

.bg-gradient {
  background-image: linear-gradient(45deg, rgba(255,255,255,0.15) 0%, rgba(255,255,255,0) 100%);
}

.product-thumbnail {
  width: 40px;
  height: 40px;
  object-fit: cover;
  border-radius: 4px;
}

.image-preview-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1050;
}

.preview-content {
  position: relative;
  max-width: 90vw;
  max-height: 90vh;
}

.preview-image {
  max-width: 100%;
  max-height: 80vh;
  object-fit: contain;
  transition: transform 0.3s ease;
}

.preview-controls {
  position: absolute;
  bottom: -50px;
  left: 50%;
  transform: translateX(-50%);
  background-color: rgba(255, 255, 255, 0.9);
  padding: 8px;
  border-radius: 20px;
  display: flex;
  gap: 8px;
}

.preview-controls .btn {
  width: 40px;
  height: 40px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}
</style>
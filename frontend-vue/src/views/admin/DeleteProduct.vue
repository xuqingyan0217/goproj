<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { deleteProduct } from '../../services/product'

const router = useRouter()
const items = ref([])

const fetchProducts = async () => {
  try {
    const response = await fetch('/api')
    const data = await response.json()
    items.value = data.items || []
  } catch (err) {
    console.error('Failed to fetch products:', err)
  }
}

const handleDelete = async (id) => {
  if (confirm('确定要删除这个产品吗？')) {
    try {
      await deleteProduct(id)
      router.push('/admin')
    } catch (error) {
      console.error('Error:', error)
      alert('删除失败，请重试')
    }
  }
}

onMounted(() => {
  fetchProducts()
})
</script>

<template>
  <div class="container-fluid">
    <div class="row g-4">
      <div v-for="item in items" :key="item.id" class="col-xl-3 col-lg-4 col-md-6 col-sm-12">
        <div class="card h-100 border-0 shadow-sm">
          <img :src="item.picture" class="card-img-top" :alt="item.name" />
          <div class="card-body">
            <p class="card-text">{{ item.name }}</p>
            <h5 class="card-title">￥{{ item.price }}</h5>
            <p class="text-muted">ID: {{ item.id }}</p>
            <button @click="handleDelete(item.id)" class="btn btn-danger">Delete</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  transition: transform 0.2s;
}

.card:hover {
  transform: translateY(-5px);
}

.card img {
  height: 200px;
  object-fit: cover;
}
</style>
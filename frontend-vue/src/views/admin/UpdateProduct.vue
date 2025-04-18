<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getProduct, updateProduct } from '../../services/product'

const router = useRouter()
const route = useRoute()
const productData = ref({
  id: '',
  name: '',
  description: '',
  price: '',
  picture: '',
  categories: []
})

const selectedCategories = ref([])
const categories = ref([])
const newCategory = ref('')

const fetchCategories = async () => {
  try {
    const response = await fetch('/category')
    if (response.ok) {
      const data = await response.json()
      categories.value = data.categories || []
    }
  } catch (err) {
    console.error('Failed to fetch categories:', err)
  }
}

const fetchProduct = async () => {
  try {
    const item = await getProduct(productData.value.id)
    productData.value = {
      id: item.id,
      name: item.name,
      description: item.description,
      price: item.price,
      picture: item.picture,
      categories: Array.isArray(item.categories) ? item.categories : []
    }
    selectedCategories.value = [...productData.value.categories]
  } catch (err) {
    console.error('Error:', err)
  }
}

const addNewCategory = () => {
  if (newCategory.value && !categories.value.includes(newCategory.value)) {
    categories.value.push(newCategory.value)
    selectedCategories.value.push(newCategory.value)
    productData.value.categories = [...selectedCategories.value]
    newCategory.value = ''
  }
}

const toggleCategory = (category) => {
  const index = selectedCategories.value.indexOf(category)
  if (index === -1) {
    selectedCategories.value.push(category)
  } else {
    selectedCategories.value.splice(index, 1)
  }
  productData.value.categories = [...selectedCategories.value]
}

const validateForm = () => {
  const errors = {}
  
  if (!productData.value.name) {
    errors.name = '商品名称不能为空'
  } else if (productData.value.name.length > 100) {
    errors.name = '商品名称不能超过100个字符'
  }

  if (!productData.value.description) {
    errors.description = '商品描述不能为空'
  }

  if (!productData.value.price) {
    errors.price = '商品价格不能为空'
  } else if (isNaN(productData.value.price) || Number(productData.value.price) <= 0) {
    errors.price = '请输入有效的价格'
  }

  if (!productData.value.picture) {
    errors.picture = '商品图片链接不能为空'
  } else {
    // 检查是否为本地静态资源路径
    const isStaticPath = productData.value.picture.startsWith('/static/image/')
    if (!isStaticPath) {
      // 如果不是静态资源路径，则验证是否为有效的URL
      try {
        new URL(productData.value.picture)
      } catch (e) {
        errors.picture = '请输入有效的图片URL或静态资源路径'
      }
    }
  }

  if (selectedCategories.value.length === 0) {
    errors.categories = '商品分类不能为空'
  }

  return Object.keys(errors).length === 0
}

const handleSubmit = async () => {
  if (!validateForm()) return

  try {
    await updateProduct({
      id: productData.value.id,
      name: productData.value.name,
      description: productData.value.description,
      price: parseFloat(productData.value.price),
      picture: productData.value.picture,
      categories: productData.value.categories
    })
    router.push('/admin')
  } catch (err) {
    console.error('Error:', err)
    alert('更新商品失败，请重试')
  }
}

onMounted(() => {
  fetchCategories()
})
</script>

<template>
  <div class="container">
    <div class="row justify-content-center">
      <div class="col-md-8 col-lg-6">
        <div class="card border-0 shadow-sm">
          <div class="card-header bg-primary bg-gradient text-white py-3">
            <h5 class="card-title mb-0">
              <i class="fas fa-edit me-2"></i>更新商品
            </h5>
          </div>
          <div class="card-body p-4">
            <form @submit.prevent="handleSubmit">
              <div class="mb-4">
                <label class="form-label">商品ID</label>
                <div class="input-group">
                  <span class="input-group-text"><i class="fas fa-hashtag"></i></span>
                  <input type="text" class="form-control" v-model="productData.id" placeholder="请输入商品ID" />
                </div>
              </div>

              <div class="mb-4">
                <label class="form-label">商品名称</label>
                <div class="input-group">
                  <span class="input-group-text"><i class="fas fa-tag"></i></span>
                  <input type="text" class="form-control" v-model="productData.name" placeholder="请输入商品名称" />
                </div>
              </div>

              <div class="mb-4">
                <label class="form-label">商品描述</label>
                <div class="input-group">
                  <span class="input-group-text"><i class="fas fa-align-left"></i></span>
                  <textarea class="form-control" v-model="productData.description" rows="3" placeholder="请输入商品描述"></textarea>
                </div>
              </div>

              <div class="mb-4">
                <label class="form-label">商品价格</label>
                <div class="input-group">
                  <span class="input-group-text"><i class="fas fa-yen-sign"></i></span>
                  <input type="number" step="0.1" class="form-control" v-model="productData.price" placeholder="请输入商品价格" />
                </div>
              </div>

              <div class="mb-4">
                <label class="form-label">商品图片</label>
                <div class="input-group mb-3">
                  <span class="input-group-text"><i class="fas fa-image"></i></span>
                  <input type="text" class="form-control" v-model="productData.picture" placeholder="请输入图片URL" />
                </div>
                <div class="image-preview mb-3" v-if="productData.picture">
                  <img :src="productData.picture" class="img-thumbnail" alt="商品图片预览" />
                </div>
              </div>

              <div class="mb-4">
                <label class="form-label">商品分类</label>
                <div class="input-group mb-3">
                  <span class="input-group-text"><i class="fas fa-folder"></i></span>
                  <input type="text" class="form-control" v-model="newCategory" placeholder="输入新分类" />
                  <button type="button" class="btn btn-outline-primary" @click="addNewCategory">添加分类</button>
                </div>
                <div class="category-list">
                  <div class="form-check" v-for="category in categories" :key="category">
                    <input class="form-check-input" type="checkbox" :id="category" 
                           :checked="selectedCategories.includes(category)"
                           @change="toggleCategory(category)">
                    <label class="form-check-label" :for="category">{{ category }}</label>
                  </div>
                </div>
                <div class="selected-categories mt-2" v-if="selectedCategories.length > 0">
                  <span class="badge bg-primary me-2 mb-2" v-for="category in selectedCategories" :key="category">
                    {{ category }}
                    <i class="fas fa-times ms-1" @click="toggleCategory(category)" style="cursor: pointer;"></i>
                  </span>
                </div>
              </div>

              <div class="d-grid gap-2">
                <button type="submit" class="btn btn-primary btn-lg">
                  <i class="fas fa-save me-2"></i>保存更新
                </button>
              </div>
            </form>
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

.image-preview img {
  max-height: 200px;
  object-fit: contain;
}

.bg-gradient {
  background-image: linear-gradient(45deg, rgba(255,255,255,0.15) 0%, rgba(255,255,255,0) 100%);
}

.category-list {
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid #dee2e6;
  border-radius: 0.25rem;
  padding: 0.5rem;
}
</style>
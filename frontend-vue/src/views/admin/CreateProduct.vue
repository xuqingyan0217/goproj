<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { createProduct } from '../../services/product'

const router = useRouter()
const productData = ref({
  name: '',
  description: '',
  price: '',
  picture: '',
  categories: []
})
const errors = ref({})
const categories = ref([])
const newCategory = ref('')
const selectedCategories = ref([])

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

const removeCategory = (category) => {
  const index = selectedCategories.value.indexOf(category)
  if (index !== -1) {
    selectedCategories.value.splice(index, 1)
    productData.value.categories = [...selectedCategories.value]
  }
}

const validateForm = () => {
  errors.value = {}
  
  if (!productData.value.name) {
    errors.value.name = '商品名称不能为空'
  } else if (productData.value.name.length > 100) {
    errors.value.name = '商品名称不能超过100个字符'
  }

  if (!productData.value.description) {
    errors.value.description = '商品描述不能为空'
  }

  if (!productData.value.price) {
    errors.value.price = '商品价格不能为空'
  } else if (isNaN(productData.value.price) || Number(productData.value.price) <= 0) {
    errors.value.price = '请输入有效的价格'
  }
  if (!productData.value.picture) {
    errors.value.picture = '商品图片链接不能为空'
  } else {
    // 检查是否为本地静态资源路径
    const isStaticPath = productData.value.picture.startsWith('/static/image/')
    if (!isStaticPath) {
      // 如果不是静态资源路径，则验证是否为有效的URL
      try {
        new URL(productData.value.picture)
      } catch (e) {
        errors.value.picture = '请输入有效的图片URL或静态资源路径'
      }
    }
  }
  if (selectedCategories.value.length === 0) {
    errors.value.categories = '商品分类不能为空'
  }
  return Object.keys(errors.value).length === 0
}

const handleSubmit = async () => {
  if (!validateForm()) return
  
  try {
    const productDataToSubmit = {
      name: productData.value.name,
      description: productData.value.description,
      price: parseFloat(productData.value.price),
      picture: productData.value.picture,
      categories: selectedCategories.value
    }

    await createProduct(productDataToSubmit)
    router.push('/admin')
  } catch (err) {
    console.error('Error:', err)
    alert('创建商品失败：网络错误，请稍后重试')
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
              <i class="fas fa-plus-circle me-2"></i>创建新商品
            </h5>
          </div>
          <div class="card-body p-4">
            <form @submit.prevent="handleSubmit">
              <div class="mb-4">
                <label class="form-label">商品名称</label>
                <div class="input-group">
                  <span class="input-group-text"><i class="fas fa-tag"></i></span>
                  <input type="text" class="form-control" :class="{'is-invalid': errors.name}" v-model="productData.name" placeholder="请输入商品名称" />
                </div>
                <div class="invalid-feedback" v-if="errors.name">{{ errors.name }}</div>
              </div>

              <div class="mb-4">
                <label class="form-label">商品描述</label>
                <div class="input-group">
                  <span class="input-group-text"><i class="fas fa-align-left"></i></span>
                  <textarea class="form-control" :class="{'is-invalid': errors.description}" v-model="productData.description" rows="3" placeholder="请输入商品描述"></textarea>
                </div>
                <div class="invalid-feedback" v-if="errors.description">{{ errors.description }}</div>
              </div>

              <div class="mb-4">
                <label class="form-label">商品价格</label>
                <div class="input-group">
                  <span class="input-group-text"><i class="fas fa-yen-sign"></i></span>
                  <input type="number" step="0.1" class="form-control" :class="{'is-invalid': errors.price}" v-model="productData.price" placeholder="请输入商品价格" />
                </div>
                <div class="invalid-feedback" v-if="errors.price">{{ errors.price }}</div>
              </div>

              <div class="mb-4">
                <label class="form-label">商品图片</label>
                <div class="input-group mb-3">
                  <span class="input-group-text"><i class="fas fa-image"></i></span>
                  <input type="text" class="form-control" :class="{'is-invalid': errors.picture}" v-model="productData.picture" placeholder="请输入图片URL" />
                </div>
                <div class="image-preview mb-3" v-if="productData.picture">
                  <img :src="productData.picture" class="img-thumbnail" alt="商品图片预览" />
                </div>
                <div class="invalid-feedback" v-if="errors.picture">{{ errors.picture }}</div>
              </div>

              <div class="mb-4">
                <label class="form-label">商品分类</label>
                <div class="category-tags mb-2">
                  <span v-for="category in selectedCategories" :key="category" class="badge bg-primary me-2 mb-2">
                    {{ category }}
                    <i class="fas fa-times ms-1 cursor-pointer" @click="removeCategory(category)"></i>
                  </span>
                </div>
                <div class="input-group mb-2">
                  <span class="input-group-text"><i class="fas fa-folder"></i></span>
                  <input type="text" class="form-control" v-model="newCategory" @keyup.enter="addNewCategory" placeholder="输入新分类并按回车添加" />
                  <button class="btn btn-outline-primary" type="button" @click="addNewCategory">
                    <i class="fas fa-plus"></i>
                  </button>
                </div>
                <div class="category-suggestions">
                  <span v-for="category in categories" 
                        :key="category" 
                        class="badge" 
                        :class="selectedCategories.includes(category) ? 'bg-secondary' : 'bg-light text-dark'" 
                        style="cursor: pointer"
                        @click="toggleCategory(category)">
                    {{ category }}
                  </span>
                </div>
                <div class="invalid-feedback" v-if="errors.categories">{{ errors.categories }}</div>
              </div>

              <div class="d-grid gap-2">
                <button type="submit" class="btn btn-primary btn-lg">
                  <i class="fas fa-save me-2"></i>保存商品
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

.category-tags .badge {
  font-size: 0.9rem;
  padding: 0.5rem 0.75rem;
}

.category-suggestions .badge {
  font-size: 0.9rem;
  padding: 0.5rem 0.75rem;
  margin-right: 0.5rem;
  margin-bottom: 0.5rem;
  transition: all 0.2s;
}

.category-suggestions .badge:hover {
  transform: translateY(-2px);
}

.cursor-pointer {
  cursor: pointer;
}
</style>
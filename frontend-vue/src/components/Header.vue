<script setup>
import { useUserStore } from '../stores/user'
import { useRouter } from 'vue-router'
import { ref, onMounted } from 'vue'
import { getCategories } from '../services/header'

const store = useUserStore()
const router = useRouter()
const categories = ref([])

const onSearch = () => {
  if (store.searchQuery) {
    router.push(`/search?q=${encodeURIComponent(store.searchQuery)}`)
  }
}

const fetchCategories = async () => {
  try {
    categories.value = await getCategories()
  } catch (err) {
    console.error('Failed to fetch categories:', err)
  }
}

onMounted(() => {
  fetchCategories()
})

// 添加页面焦点事件监听，确保类别列表及时更新
if (typeof window !== 'undefined') {
  window.addEventListener('focus', () => {
    fetchCategories()
  })
}
</script>

<template>
  <header>
    <nav class="navbar navbar-expand-lg bg-body-tertiary">
      <div class="container-fluid">
        <img class="navbar-brand" src="/static/image/logo.jpg" alt="CloudWeGo" style="height: 3em" />
        <button
          class="navbar-toggler"
          type="button"
          data-bs-toggle="collapse"
          data-bs-target="#navbarSupportedContent"
          aria-controls="navbarSupportedContent"
          aria-expanded="false"
          aria-label="Toggle navigation"
        >
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarSupportedContent">
          <ul class="navbar-nav me-auto mb-2 mb-lg-0">
            <li class="nav-item">
              <router-link class="nav-link active" aria-current="page" to="/">Home</router-link>
            </li>
            <li class="nav-item dropdown">
              <a
                class="nav-link dropdown-toggle"
                href="#"
                role="button"
                data-bs-toggle="dropdown"
                aria-expanded="false"
              >
                Categories
              </a>
              <ul class="dropdown-menu">
                <li v-for="category in categories" :key="category">
                  <router-link class="dropdown-item" :to="`/category/${category}`">{{ category }}</router-link>
                </li>
              </ul>
            </li>
            <li class="nav-item">
              <router-link class="nav-link" to="/about">About</router-link>
            </li>
          </ul>

          <form class="d-flex ms-auto" role="search" @submit.prevent="onSearch">
            <input
              class="form-control me-2"
              type="search"
              name="q"
              placeholder="Search"
              aria-label="Search"
              v-model="store.searchQuery"
            />
            <button class="btn btn-outline-success" type="submit">Search</button>
          </form>

          <div class="ms-3">
            <router-link to="/cart" class="text-decoration-none cart-icon-wrapper">
              <i class="fa-solid fa-cart-shopping fa-lg"></i>
              <span v-if="store.cartNum > 0" class="cart-badge">{{ store.cartNum }}</span>
            </router-link>
          </div>

          <template v-if="store.userId">
            <div class="dropdown">
              <div class="ms-3 dropdown-toggle" data-bs-toggle="dropdown">
                <i class="fa-solid fa-user fa-lg"></i>
                <span>User</span>
              </div>
              <ul class="dropdown-menu dropdown-menu-end" style="margin-top: 15px">
                <li><router-link class="dropdown-item" to="/order">Order Center</router-link></li>
                <li><router-link class="dropdown-item" to="/admin">Admin Center</router-link></li>
                <li><router-link class="dropdown-item" to="/ai">AI Order Center</router-link></li>
                <li>
                  <button class="dropdown-item" @click="store.logout">Logout</button>
                </li>
              </ul>
            </div>
          </template>
          <template v-else>
            <div class="ms-3">
              <router-link class="btn btn-primary" to="/sign-in">Sign In</router-link>
            </div>
          </template>
        </div>
      </div>
    </nav>
    <div class="bg-primary text-center text-white pt-1 pb-1">
      This website is hosted for demo purposes only. It is not an actual shop.
    </div>
    <div v-if="store.error" class="alert alert-danger text-center" role="alert">{{ store.error }}</div>
    <div v-if="store.warning" class="alert alert-warning text-center" role="alert">{{ store.warning }}</div>
  </header>
</template>

<style scoped>
.cart-icon-wrapper {
  position: relative;
  display: inline-block;
  color: #333;
}

.cart-badge {
  position: absolute;
  top: -8px;
  right: -8px;
  min-width: 16px;
  height: 16px;
  line-height: 16px;
  text-align: center;
  background-color: #dc3545;
  color: white;
  border-radius: 50%;
  font-size: 0.75rem;
  padding: 0 4px;
}
</style>
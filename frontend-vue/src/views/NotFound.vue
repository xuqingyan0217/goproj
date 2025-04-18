<script setup>
</script>

<template>
  <div class="not-found">
    <div class="container text-center py-5">
      <h1 class="display-1">404</h1>
      <h2 class="mb-4">页面未找到</h2>
      <p class="lead mb-4">抱歉，您访问的页面不存在。</p>
      <router-link to="/" class="btn btn-primary">返回首页</router-link>
    </div>
  </div>
</template>

<style scoped>
.not-found {
  min-height: 60vh;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
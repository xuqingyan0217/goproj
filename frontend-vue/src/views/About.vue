<script setup>
</script>

<template>
  <div class="container mt-4">
    <div class="row">
      <div class="col-md-8 offset-md-2">
        <h1 class="text-center mb-4">关于我们</h1>
        <div class="card">
          <div class="card-body">
            <h5 class="card-title">CloudWeGo商城</h5>
            <p class="card-text">
              CloudWeGo商城是一个基于现代微服务架构构建的电子商务平台演示项目。我们致力于为用户提供最佳的在线购物体验。
            </p>
            <h5 class="mt-4">我们的特点</h5>
            <ul>
              <li>丰富的商品种类</li>
              <li>便捷的购物流程</li>
              <li>安全的支付系统</li>
              <li>快速的配送服务</li>
              <li>专业的客户支持</li>
            </ul>
            <h5 class="mt-4">联系我们</h5>
            <p>
              如果您有任何问题或建议，请随时与我们联系：
              <br>
              邮箱：support@cloudwego.com
              <br>
              电话：xxx-xxx-xxxx
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}

.card-title {
  color: #2c3e50;
}

ul {
  list-style-type: none;
  padding-left: 0;
}

ul li {
  margin-bottom: 0.5rem;
  padding-left: 1.5rem;
  position: relative;
}

ul li:before {
  content: "✓";
  position: absolute;
  left: 0;
  color: #28a745;
}
</style>
<template>
  <div class="container mt-4">
    <div class="row row-cols-1 row-cols-md-2 g-4">
      <div v-for="order in orders" :key="order.OrderId" class="col">
        <div class="card h-100">
          <div class="card-header d-flex justify-content-between align-items-center bg-light">
            <div>
              <span class="text-muted">订单号: {{ order.OrderId }}</span>
              <span v-if="order.UserCurrency === 'WAIT'" class="ms-2 badge bg-warning text-dark">订单尚未完成 ({{ getCountdown(order.CreatedDate) }})</span>
              <span v-else class="ms-2 badge bg-success">订单已完成</span>
            </div>
            <span class="badge bg-secondary">{{ order.CreatedDate }}</span>
          </div>
          <div class="card-body p-3">
            <div class="order-items">
              <div v-for="item in order.Items" :key="item.ProductName" class="order-item mb-2">
                <div class="d-flex align-items-center">
                  <img :src="item.Picture" class="product-image me-3" :alt="item.ProductName">
                  <div class="flex-grow-1">
                    <h6 class="mb-1">{{ item.ProductName }}</h6>
                    <div class="d-flex justify-content-between align-items-center">
                      <span class="text-muted">数量: {{ item.Qty }}</span>
                      <span class="text-primary">¥{{ item.Cost }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="card-footer bg-white">
            <div class="d-flex justify-content-between align-items-center">
              <button v-if="order.UserCurrency === 'WAIT'" @click="cancelOrderHandler(order.OrderId)" class="btn btn-sm btn-outline-danger me-2">
                取消订单
              </button>
              <button v-if="order.UserCurrency === 'WAIT'" @click="handlePayment(order.OrderId)" class="btn btn-sm btn-primary">
                立即支付
              </button>
              <span class="fw-bold ms-auto">总计: ¥{{ order.Cost }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <!-- 支付模态框 -->
  <div v-if="showPaymentModal" class="modal fade show d-block" tabindex="-1" style="background-color: rgba(0,0,0,0.5);">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">输入支付信息</h5>
          <button type="button" class="btn-close" @click="closePaymentModal"></button>
        </div>
        <div class="modal-body">
          <div class="mb-3">
            <label class="form-label">信用卡号</label>
            <input type="text" class="form-control" value="424242424242424242" readonly>
          </div>
          <div class="mb-3">
            <label class="form-label">到期月份</label>
            <input type="text" class="form-control" value="12" readonly>
          </div>
          <div class="mb-3">
            <label class="form-label">到期年份</label>
            <input type="text" class="form-control" value="2030" readonly>
          </div>
          <div class="mb-3">
            <label class="form-label">支付方式</label>
            <input type="text" class="form-control" value="card" readonly>
          </div>
          <div class="mb-3">
            <label for="email" class="form-label">邮箱</label>
            <input type="email" class="form-control" id="email" v-model="email" placeholder="请输入您的邮箱地址">
          </div>
          <div class="mb-3">
            <label for="cvv" class="form-label">CVV码</label>
            <input type="password" class="form-control" id="cvv" v-model="cvv" maxlength="3" placeholder="请输入信用卡背面的3位数字">
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="closePaymentModal">取消</button>
          <button type="button" class="btn btn-primary" @click="processPayment" :disabled="isProcessingPayment || !cvv || cvv.length !== 3">
            <span v-if="isProcessingPayment" class="spinner-border spinner-border-sm me-2"></span>
            {{ isProcessingPayment ? '处理中...' : '确认支付' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { getOrders, cancelOrder } from '../services/order'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import { repayOrder } from '../services/checkout'

// 支付相关的状态
const showPaymentModal = ref(false)
const currentOrderId = ref('')
const cvv = ref('')
const email = ref('')
const isProcessingPayment = ref(false)
const autoCheckInterval = ref(null)

const router = useRouter()
const orders = ref([])

// 计算倒计时时间
const getCountdown = (createdDate) => {
  const now = new Date()
  const orderDate = new Date(createdDate)
  const timeDiff = now - orderDate
  const totalSeconds = 600 - Math.floor(timeDiff / 1000) // 10分钟 = 600秒
  
  if (totalSeconds <= 0) {
    return '即将取消'
  }
  
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60
  
  return `${minutes}分${seconds}秒后自动取消`
}

// 更新倒计时的定时器
const countdownTimer = ref(null)

// 更新所有订单的倒计时显示
const updateCountdowns = () => {
  orders.value = [...orders.value]
}

onMounted(() => {
  fetchOrders()
  // 启动定时检查，每分钟检查一次
  autoCheckInterval.value = setInterval(checkAndCancelTimeoutOrders, 60000)
  // 启动倒计时更新，每秒更新一次
  countdownTimer.value = setInterval(updateCountdowns, 1000)
})

onUnmounted(() => {
  // 清理定时器
  if (autoCheckInterval.value) {
    clearInterval(autoCheckInterval.value)
  }
  if (countdownTimer.value) {
    clearInterval(countdownTimer.value)
  }
})

// 检查并自动取消超时订单
const checkAndCancelTimeoutOrders = () => {
  const now = new Date()
  orders.value.forEach(async order => {
    if (order.UserCurrency === 'WAIT') {
      const orderDate = new Date(order.CreatedDate)
      const timeDiff = now - orderDate
      const minutesDiff = Math.floor(timeDiff / (1000 * 60))
      
      if (minutesDiff >= 10) {
        try {
          await cancelOrder(order.OrderId)
          await fetchOrders() // 刷新订单列表
        } catch (error) {
          console.error('Failed to auto cancel order:', error)
        }
      }
    }
  })
}

// 处理支付按钮点击
const handlePayment = (orderId) => {
  currentOrderId.value = orderId
  showPaymentModal.value = true
  cvv.value = ''
}

// 关闭支付模态框
const closePaymentModal = () => {
  showPaymentModal.value = false
  currentOrderId.value = ''
  cvv.value = ''
  isProcessingPayment.value = false
}

// 处理支付请求
const processPayment = async () => {
  if (!cvv.value || cvv.value.length !== 3) {
    ElMessage.warning('请输入正确的CVV码')
    return
  }

  if (!email.value || !email.value.includes('@')) {
    ElMessage.warning('请输入有效的邮箱地址')
    return
  }

  const currentOrder = orders.value.find(order => order.OrderId === currentOrderId.value)
  if (!currentOrder) {
    ElMessage.error('订单信息不存在')
    return
  }

  isProcessingPayment.value = true
  try {
    await repayOrder(currentOrderId.value, cvv.value, email.value, currentOrder.Cost)
    ElMessage.success('支付成功')
    closePaymentModal()
    await fetchOrders() // 刷新订单列表
  } catch (error) {
    console.error('Payment failed:', error)
    ElMessage.error('支付失败，请重试')
  } finally {
    isProcessingPayment.value = false
  }
}

// 获取订单列表
const fetchOrders = async () => {
  try {
    const response = await getOrders()
    // 对订单按创建时间进行降序排序
    orders.value = response.data.orders.sort((a, b) => {
      return new Date(b.CreatedDate) - new Date(a.CreatedDate)
    })
  } catch (error) {
    console.error('Failed to fetch orders:', error)
    ElMessage.error('获取订单列表失败')
  }
}

// 取消订单处理函数
const cancelOrderHandler = async (orderId) => {
  try {
    // 显示确认对话框
    await ElMessageBox.confirm('确定要取消此订单吗？此操作不可逆。', '取消订单', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    // 调用取消订单API
    await cancelOrder(orderId)
    ElMessage.success('订单已成功取消')
    
    // 重新获取订单列表以更新UI
    await fetchOrders()
  } catch (error) {
    if (error === 'cancel') {
      // 用户取消了确认对话框
      return
    }
    console.error('Failed to cancel order:', error)
    ElMessage.error('取消订单失败')
  }
}

onMounted(() => {
  fetchOrders()
  // 启动定时检查，每分钟检查一次
  autoCheckInterval.value = setInterval(checkAndCancelTimeoutOrders, 60000)
})

onUnmounted(() => {
  // 清理定时器
  if (autoCheckInterval.value) {
    clearInterval(autoCheckInterval.value)
  }
})
</script>

<style scoped>
.card {
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.product-image {
  width: 60px;
  height: 60px;
  object-fit: contain;
  border-radius: 4px;
}

.order-items {
  max-height: 300px;
  overflow-y: auto;
}

.order-item {
  padding: 8px;
  border-radius: 6px;
  background-color: #f8f9fa;
}

.order-item:hover {
  background-color: #f0f1f2;
}

.badge {
  font-weight: normal;
}
</style>
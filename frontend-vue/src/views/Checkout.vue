<template>
  <div class="row mb-5">
    <div class="col-lg-8 col-sm-12">
      <div class="card shadow-sm rounded p-4 mb-4">
        <form @submit.prevent="handleSubmit">
          <h4 class="mb-4 border-bottom pb-2">Contact</h4>
          <label for="email" class="form-label col-12">
            <input class="form-control form-control-lg" id="email" type="email" placeholder="Email" v-model="formData.email"
                   aria-label="email">
          </label>
          <h4 class="mb-4 mt-5 border-bottom pb-2">Delivery</h4>
          <div class="mb-4 mt-3 col-12 row">
            <label for="firstname" class="col-md-6 col-sm-12 mb-3 mb-md-0">
              <input type="text" id="firstname" class="form-control" placeholder="First name"
                     v-model="formData.firstname">
            </label>
            <label for="lastname" class="col-md-6 col-sm-12">
              <input type="text" id="lastname" class="form-control" placeholder="Last name" 
                     v-model="formData.lastname">
            </label>
          </div>
          <label for="street" class="mb-4 col-12 form-label">
            <input type="text" class="form-control" placeholder="Street" v-model="formData.street"
                   id="street">
          </label>
          <label for="zipcode" class="mb-4 form-label col-12">
            <input type="text" class="form-control" id="zipcode" v-model="formData.zipcode" placeholder="zipcode">
          </label>
          <div class="mb-4 col-12 row">
            <label for="city" class="col-md-6 col-sm-12 mb-3 mb-md-0">
              <input type="text" id="city" class="form-control" placeholder="City" v-model="formData.city">
            </label>
            <label for="province" class="col-md-6 col-sm-12">
              <input type="text" id="province" class="form-control" v-model="formData.province" placeholder="Province">
            </label>
          </div>
          <label for="country" class="mb-4 form-label col-12">
            <input type="text" class="form-control" id="country" v-model="formData.country" placeholder="Country">
          </label>
          <h4 class="mb-4 mt-5 border-bottom pb-2">Payment</h4>
          <label for="card-num" class="form-label col-12 mb-4">
            <input type="text" id="card-num" class="form-control" v-model="formData.cardNum" placeholder="Card number">
          </label>
          <div class="mb-4 col-12 row">
            <label for="expiration-month" class="col-md-4 col-sm-12 mb-3 mb-md-0">
              <input type="text" id="expiration-month" v-model="formData.expirationMonth" class="form-control"
                     placeholder="Expiration Month">
            </label>
            <label for="expiration-year" class="col-md-4 col-sm-12 mb-3 mb-md-0">
              <input type="text" id="expiration-year" v-model="formData.expirationYear" class="form-control"
                     placeholder="Expiration Year">
            </label>
            <label for="cvv" class="col-md-4 col-sm-12">
              <input type="text" id="cvv" class="form-control" v-model="formData.cvv" placeholder="cvv" required>
            </label>
          </div>
          <div class="payment-methods mb-4">
            <div class="form-check mb-2 payment-option">
              <input class="form-check-input" type="radio" v-model="formData.payment" id="card" value="card" checked>
              <label class="form-check-label" for="card">Card</label>
            </div>
            <div class="form-check mb-2 payment-option disabled">
              <input class="form-check-input" type="radio" v-model="formData.payment" id="stripe" value="stripe" disabled>
              <label class="form-check-label" for="stripe">Stripe</label>
            </div>
            <div class="form-check mb-2 payment-option disabled">
              <input class="form-check-input" type="radio" v-model="formData.payment" id="paypal" value="paypal" disabled>
              <label class="form-check-label" for="paypal">Paypal</label>
            </div>
            <div class="form-check mb-2 payment-option disabled">
              <input class="form-check-input" type="radio" v-model="formData.payment" id="wechat" value="wechat" disabled>
              <label class="form-check-label" for="wechat">Wechat</label>
            </div>
            <div class="form-check payment-option disabled">
              <input class="form-check-input" type="radio" v-model="formData.payment" id="alipay" value="alipay" disabled>
              <label class="form-check-label" for="alipay">Alipay</label>
            </div>
          </div>
          <div class="mt-4 mb-3 text-end">
            <div class="h4 mb-4 text-danger">Total: ¥{{ total }}</div>
            <div class="d-flex justify-content-between">
              <button type="button" class="btn btn-outline-secondary btn-lg px-5" @click="handleCancel">取消</button>
              <button type="submit" class="btn btn-success btn-lg px-5">Pay</button>
            </div>
          </div>
        </form>
      </div>
    </div>
    <div class="col-lg-4 col-sm-12">
      <div class="card shadow-sm rounded">
        <div class="card-header bg-light">
          <h5 class="mb-0">Order Summary</h5>
        </div>
        <ul class="list-group list-group-flush">
          <li v-for="item in items" :key="item.id" class="list-group-item hover-effect">
            <div class="card border-0">
              <div class="card-body row align-items-center p-2">
                <div class="col-4">
                  <img :src="item.Picture" class="img-fluid" style="max-height: 80px; object-fit: contain" alt="">
                </div>
                <div class="col-8">
                  <h6 class="mb-2">{{ item.Name }}</h6>
                  <div class="text-muted mb-1">Single Price: ¥{{ item.Price }}</div>
                  <div class="text-muted">Qty: {{ item.Qty }}</div>
                </div>
              </div>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { getCheckout, submitCheckout } from '../services/checkout'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const items = ref([])
const total = computed(() => {
  return items.value.reduce((sum, item) => sum + item.Price * item.Qty, 0).toFixed(2)
})

const formData = ref({
  email: 'abc@example.com',
  firstname: '三',
  lastname: '张',
  street: '枣林街 长江路80号',
  zipcode: '310000',
  city: '南阳市',
  province: '河南省',
  country: '中国',
  cardNum: '424242424242424242',
  expirationMonth: '12',
  expirationYear: '2030',
  cvv: '',
  payment: 'card',
  flag: 1
})

onMounted(async () => {
  try {
    const response = await getCheckout()
    items.value = response.items
  } catch (error) {
    console.error('Failed to fetch cart items:', error)
  }
})

const handleSubmit = async () => {
  try {
    formData.value.flag = 1
    await submitCheckout(formData.value)
    router.push('/checkout/waiting?status=success')
  } catch (error) {
    console.error('Checkout failed:', error)
  }
}
const handleCancel = () => {
  ElMessageBox.confirm('确定要取消支付吗？您可以稍后在订单列表中查看并处理此订单。', '取消支付', {
    confirmButtonText: '确定',
    cancelButtonText: '返回',
    type: 'warning'
  }).then(async () => {
    formData.value.flag = 0
    try {
      await submitCheckout(formData.value)
      router.push('/checkout/waiting?status=cancel')
      ElMessage.info('您可以在订单列表中查看并管理订单')
    } catch (error) {
      console.error('Cancel payment failed:', error)
      ElMessage.error('取消支付失败')
    }
  }).catch(() => {})
}
</script>

<style scoped>
.hover-effect:hover {
  background-color: #f8f9fa;
  transition: background-color 0.3s ease;
}

.payment-option {
  padding: 0.75rem;
  border-radius: 0.375rem;
  transition: background-color 0.3s ease;
}

.payment-option:hover:not(.disabled) {
  background-color: #f8f9fa;
}

.payment-option.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-control:focus {
  box-shadow: 0 0 0 0.2rem rgba(40, 167, 69, 0.25);
  border-color: #28a745;
}

.btn-success {
  transition: all 0.3s ease;
}

.btn-success:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 5px rgba(0,0,0,0.2);
}
</style>

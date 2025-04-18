import service from '../utils/request'

// 订单相关的API服务

// 获取订单列表
export const getOrders = async () => {
  try {
    const response = await service.get('/order')
    console.log('获取订单列表响应:', response.data)
    return { data: { orders: response.data.orders || [] } }
  } catch (error) {
    console.error('获取订单列表失败:', error)
    throw error
  }
}

// 取消订单
export const cancelOrder = async (orderId) => {
  try {
    const response = await service.post('/order/cancel', { order_id: orderId })
    console.log('取消订单响应:', response.data)
    return response.data
  } catch (error) {
    console.error('取消订单失败:', error)
    throw error
  }
}

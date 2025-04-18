import service from '../utils/request'

// AI预下单
export const submitAIOrder = (orderText) => {
  return service.post('/ai/aiorder', {
    order: orderText
  })
}

// 获取AI订单信息
export const getAIOrderInfo = (orderList) => {
  return service.post('/ai/ailists', {
    orderList: orderList
  })
}
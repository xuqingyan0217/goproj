import service from '../utils/request'

// 获取结账信息
export const getCheckout = async () => {
  try {
    const response = await service.get('/checkout')
    return response.data
  } catch (error) {
    console.error('Error getting checkout info:', error)
    throw error
  }
}

// 提交结账请求
export const submitCheckout = async (formData) => {
  try {
    // 验证月份值
    const month = parseInt(formData.expirationMonth, 10)
    if (isNaN(month) || month < 1 || month > 12) {
      throw new Error('无效的月份值，请输入1-12之间的数字')
    }

    const requestData = {
      email: formData.email,
      firstname: formData.firstname,
      lastname: formData.lastname,
      street: formData.street,
      city: formData.city,
      province: formData.province,
      country: formData.country,
      zipcode: formData.zipcode,
      cardNum: formData.cardNum,
      expirationMonth: parseInt(formData.expirationMonth, 10),
      expirationYear: parseInt(formData.expirationYear, 10),
      cvv: parseInt(formData.cvv, 10),
      payment: 'card',
      flag: formData.flag
    }

    const formDataParams = new URLSearchParams()
    for (const [key, value] of Object.entries(requestData)) {
      formDataParams.append(key, value)
    }
    const response = await service.post('/checkout/waiting', formDataParams, {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      }
    })
    return response.data
  } catch (error) {
    console.error('Error submitting checkout:', error)
    throw error
  }
}

// 获取结账结果
export const getCheckoutResult = async () => {
  try {
    const response = await service.get('/checkout/result')
    return response.data
  } catch (error) {
    console.error('Error getting checkout result:', error)
    throw error
  }
}

// 订单二次支付
export const repayOrder = async (orderId, cvv, email, total, expirationMonth = 12, expirationYear = 2030) => {
  try {
    const requestData = {
      order_id: orderId,
      cardNum: '424242424242424242',
      cvv: parseInt(cvv, 10),
      expirationMonth: parseInt(expirationMonth, 10),
      expirationYear: parseInt(expirationYear, 10),
      email: email,
      total: parseFloat(total)
    }

    const formDataParams = new URLSearchParams()
    for (const [key, value] of Object.entries(requestData)) {
      formDataParams.append(key, value)
    }

    const response = await service.post('/checkout/repay', formDataParams, {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      }
    })
    return response.data
  } catch (error) {
    console.error('Error repaying order:', error)
    throw error
  }
}
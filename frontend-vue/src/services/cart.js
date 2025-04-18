import service from '../utils/request'

// 购物车相关的API服务

// 获取购物车商品列表
export const getCartItems = async () => {
  try {
    const response = await service.get('/cart')
    return {
      items: response.data.items || [],
      total: response.data.total || 0
    }
  } catch (error) {
    throw new Error('获取购物车失败')
  }
}

// 添加商品到购物车
export const addToCart = async (productId, quantity) => {
  try {
    const response = await service.post('/cart', {
      product_id: productId,
      product_num: quantity
    })
    return {
      success: true,
      cartNum: response.data.cart_num || 0
    }
  } catch (error) {
    throw new Error('添加到购物车失败')
  }
}

// 从购物车中删除商品
export const removeFromCart = async (productId) => {
  try {
    const response = await service.delete(`/cart?product_id=${productId}`)
    return response.data
  } catch (error) {
    throw new Error('从购物车删除失败')
  }
}

// 更新购物车商品数量
export const updateCartItemQuantity = async (productId, quantity) => {
  try {
    const response = await service.post('/cart/update', {
      product_id: productId,
      quantity: quantity
    })
    return response.data
  } catch (error) {
    throw new Error('更新购物车数量失败')
  }
}
import service from '../utils/request'

// 商品相关的API服务

// 获取商品详情
export const getProduct = async (id) => {
  try {
    const response = await service.get(`/product?id=${id}`)
    return response.data.item
  } catch (error) {
    console.error('获取商品详情失败:', error)
    throw error
  }
}

// 搜索商品
export const searchProducts = async (query) => {
  try {
    const response = await service.get(`/search?q=${encodeURIComponent(query)}`)
    return response.data.items
  } catch (error) {
    console.error('搜索商品失败:', error)
    throw error
  }
}

// 创建商品
export const createProduct = async (productData) => {
  try {
    const response = await service.post('/product/api/create', productData)
    return response.data
  } catch (error) {
    console.error('创建商品失败:', error)
  }
}

// 更新商品
export const updateProduct = async (productData) => {
  try {
    const formattedData = {
      ...productData,
      id: parseInt(productData.id, 10)
    }
    const response = await service.post('/product/api/update', formattedData)
    console.log('更新商品响应:', response.data)
    return response.data
  } catch (error) {
    console.error('更新商品失败:', error)
  }
}

// 删除商品
export const deleteProduct = async (id) => {
  try {
    const response = await service.delete(`/product/api/delete?id=${id}`)
    return response.data
  } catch (error) {
    console.error('删除商品失败:', error)
  }
}
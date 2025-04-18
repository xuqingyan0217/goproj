import service from '../utils/request'

// 获取管理后台统计数据
export const getAdminStats = async () => {
  try {
    // 获取商品数据
    const productResponse = await service.get('/api')
    const productData = productResponse.data
    
    // 获取分类数据
    const categoryResponse = await service.get('/category')
    const categoryData = categoryResponse.data
    
    // 获取订单数据
    const orderResponse = await service.get('/order')
    const orderData = orderResponse.data
    
    return {
      totalProducts: productData.items?.length || 0,
      recentProducts: productData.items?.slice(0, 5) || [],
      totalCategories: categoryData.categories?.length || 0,
      totalOrders: orderData.orders?.length || 0
    }
  } catch (error) {
    throw new Error('获取管理后台统计数据失败')
  }
}
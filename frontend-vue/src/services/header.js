import service from '../utils/request'

// 获取所有商品分类
export const getCategories = async () => {
  try {
    const response = await service.get('/category')
    return response.data.categories || []
  } catch (error) {
    throw new Error('获取分类列表失败')
  }
}
import service from '../utils/request'

// 获取首页数据
export const getHomeData = async () => {
  try {
    const response = await service.get('/api')
    return response.data
  } catch (error) {
    throw new Error('获取首页数据失败')
  }
}
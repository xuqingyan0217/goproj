import { setToken, removeToken } from '../utils/token'
import service from '../utils/request'

export const login = async (email, password, next = '/') => {
  try {
    const response = await service.post(`/auth/login?next=${next}`, {
      email,
      password
    })
    
    // 从响应头中获取token并保存
    const authHeader = response.headers['authorization']
    if (authHeader && authHeader.startsWith('Bearer ')) {
      const token = authHeader.substring(7) // 去掉'Bearer '前缀
      setToken(token)
    }
    
    return response.data
  } catch (error) {
    if (error.response && error.response.data && error.response.data.message) {
      throw new Error(error.response.data.message)
    }
    throw new Error('登录失败，请检查网络连接')
  }
}

export const register = async (email, password, passwordConfirm) => {
  try {
    const response = await service.post('/auth/register', {
      email,
      password,
      password_confirm: passwordConfirm
    })
    return response.data
  } catch (error) {
    if (error.response && error.response.data && error.response.data.message) {
      throw new Error(error.response.data.message)
    }
    throw new Error('注册失败，请稍后重试')
  }
}

export const logout = async () => {
  try {
    const response = await service.post('/auth/logout')
    // 登出时清除token
    removeToken()
    return response.data
  } catch (error) {
    throw new Error('Logout failed')
  }
}

// 刷新token
export const refreshToken = async () => {
  try {
    const response = await service.post('/auth/refresh-token')
    
    // 从响应头中获取新的token
    const authHeader = response.headers['authorization']
    if (authHeader && authHeader.startsWith('Bearer ')) {
      const token = authHeader.substring(7) // 去掉'Bearer '前缀
      setToken(token)
      return token
    }
    return null
  } catch (error) {
    console.error('Token refresh failed:', error)
    return null
  }
}
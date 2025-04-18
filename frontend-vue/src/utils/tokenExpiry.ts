// JWT Token过期检测工具
import { getToken, removeToken } from './token'

// 解析JWT token
const parseJwt = (token: string): any => {
  try {
    // 获取payload部分（JWT的第二部分）
    const base64Url = token.split('.')[1]
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    )
    return JSON.parse(jsonPayload)
  } catch (error) {
    console.error('Invalid token format:', error)
    return null
  }
}

// 检查token是否过期
export const isTokenExpired = (): boolean => {
  const token = getToken()
  if (!token) return true

  try {
    const decoded = parseJwt(token)
    if (!decoded || !decoded.exp) return true

    // exp是Unix时间戳（秒），需要转换为毫秒
    const expirationTime = decoded.exp * 1000
    return Date.now() >= expirationTime
  } catch (error) {
    console.error('Error checking token expiration:', error)
    return true
  }
}

// 检查token是否即将过期（默认5分钟内）
export const isTokenExpiringSoon = (minutesThreshold: number = 5): boolean => {
  const token = getToken()
  if (!token) return true

  try {
    const decoded = parseJwt(token)
    if (!decoded || !decoded.exp) return true

    const expirationTime = decoded.exp * 1000
    const warningTime = expirationTime - minutesThreshold * 60 * 1000
    return Date.now() >= warningTime
  } catch (error) {
    console.error('Error checking token expiration:', error)
    return true
  }
}

// 获取token过期时间
export const getTokenExpiryTime = (): Date | null => {
  const token = getToken()
  if (!token) return null

  try {
    const decoded = parseJwt(token)
    if (!decoded || !decoded.exp) return null

    return new Date(decoded.exp * 1000)
  } catch (error) {
    console.error('Error getting token expiry time:', error)
    return null
  }
}

// @ts-ignore
// 导入刷新token的函数
import { refreshToken } from '../services/auth'

// 设置token过期监控
export const setupTokenExpiryMonitor = (options: {
  checkInterval?: number; // 检查间隔（毫秒）
  warningThreshold?: number; // 警告阈值（分钟）
  onExpiringSoon?: () => void; // 即将过期回调
  onExpired?: () => void; // 已过期回调
  autoRefresh?: boolean; // 是否自动刷新token
} = {}): (() => void) => {
  const {
    checkInterval = 60000, // 默认每分钟检查一次
    warningThreshold = 5, // 默认5分钟前警告
    autoRefresh = true, // 默认自动刷新token
    onExpiringSoon = async () => {
      console.warn(`Token will expire in less than ${warningThreshold} minutes`)
      // 自动刷新token
      if (autoRefresh) {
        try {
          const newToken = await refreshToken()
          if (newToken) {
            console.log('Token refreshed successfully')
            return // 刷新成功，不需要进一步处理
          }
        } catch (error) {
          console.error('Failed to refresh token:', error)
        }
      }
    },
    onExpired = () => {
      console.warn('Token has expired')
      removeToken()
      
      // 清除用户状态和购物车数据
      localStorage.removeItem('userId')
      localStorage.removeItem('cartNum')
      localStorage.clear() // 清除所有localStorage数据
      
      // 重定向到登录页
      const currentPath = window.location.pathname
      window.location.href = `/sign-in?next=${currentPath}`
    }
  } = options

  // 启动定时器
  const intervalId = setInterval(() => {
    if (isTokenExpired()) {
      onExpired()
      clearInterval(intervalId)
    } else if (isTokenExpiringSoon(warningThreshold)) {
      onExpiringSoon()
    }
  }, checkInterval)

  // 返回清理函数
  return () => clearInterval(intervalId)
}
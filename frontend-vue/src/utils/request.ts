import axios from 'axios'
import { getToken } from './token'
import { useUserStore } from '../stores/user'

// 创建axios实例
const service = axios.create({
  baseURL: '', // API的base_url
  timeout: 15000 // 请求超时时间
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    // 在请求发送之前添加token
    const token = getToken()
    if (token) {
      // 将token添加到请求头中
      config.headers['Authorization'] = `Bearer ${token}`
    }
    
    return config
  },
  error => {
    // 请求错误处理
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    // 如果响应中包含新的token（例如刷新token的情况），可以在这里处理
    const authHeader = response.headers['authorization']
    if (authHeader && authHeader.startsWith('Bearer ')) {
      const token = authHeader.substring(7) // 去掉'Bearer '前缀
      // 更新token
      import('./token').then(tokenModule => {
        tokenModule.setToken(token)
      })
    }
    return response
  },
  error => {
    // 响应错误处理
    if (error.response) {
      // 服务器返回错误状态码
      if (error.response.status === 401) {
        // 未授权，可能是token过期，显示友好提示
        alert('登录已过期，请重新登录')
        
        // 延迟1秒后再执行重定向，让用户有时间看到提示
        setTimeout(() => {
          import('./token').then(tokenModule => {
            tokenModule.removeToken()
            // 清除用户状态
            const userStore = useUserStore()
            userStore.setUserId(null)
            userStore.setCartNum(0)
            // 获取当前路径，用于登录后重定向回来
            const currentPath = window.location.pathname
            window.location.href = `/sign-in?next=${currentPath}`
          })
        }, 1000)
      } else if (error.response.status === 403) {
        // 权限不足，显示弹窗提示
        alert('权限不足')
        return Promise.reject(new Error('权限不足'))
      }
    }
    return Promise.reject(error)
  }
)

export default service
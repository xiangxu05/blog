import axios, { type AxiosInstance, type AxiosResponse, type AxiosError } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

// 创建 axios 实例
const api: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
  timeout: 10000,
  withCredentials: true, // 支持跨域cookie
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // 依赖cookie自动发送，无需手动添加认证信息
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response: AxiosResponse) => {
    return response
  },
  (error: AxiosError) => {
    // 处理错误响应
    if (error.response) {
      const { status, data } = error.response
      
      switch (status) {
        case 401:
          // 未授权，跳转到登录页（cookie会自动过期）
          ElMessage.error('登录已过期，请重新登录')
          router.push('/login')
          break
        case 403:
          ElMessage.error('权限不足')
          router.push('/403')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          // 显示后端返回的错误信息
          const message = (data as any)?.message || '请求失败'
          ElMessage.error(message)
      }
    } else if (error.request) {
      // 网络错误
      ElMessage.error('网络连接失败，请检查网络设置')
    } else {
      // 其他错误
      ElMessage.error('请求配置错误')
    }
    
    return Promise.reject(error)
  }
)

export default api
import axios, { type AxiosResponse, type AxiosError } from 'axios'
import { ElMessage, ElLoading } from 'element-plus'
import type { LoadingInstance } from 'element-plus/es/components/loading/src/loading'

// 创建axios实例
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000,
  withCredentials: true // 支持跨域cookie
})

// 全局loading实例
let loadingInstance: LoadingInstance | null = null
let requestCount = 0

// 显示loading
const showLoading = () => {
  if (requestCount === 0) {
    loadingInstance = ElLoading.service({
      text: '加载中...',
      background: 'rgba(0, 0, 0, 0.7)'
    })
  }
  requestCount++
}

// 隐藏loading
const hideLoading = () => {
  requestCount--
  if (requestCount === 0 && loadingInstance) {
    loadingInstance.close()
    loadingInstance = null
  }
}

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    // 显示loading（可通过config.hideLoading = true来禁用）
    if (!(config as any).hideLoading) {
      showLoading()
    }

    // 可以在这里添加token等认证信息
    // const token = localStorage.getItem('token')
    // if (token) {
    //   config.headers.Authorization = `Bearer ${token}`
    // }

    return config
  },
  (error) => {
    hideLoading()
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    hideLoading()

    const { code, message, data } = response.data

    // 根据后端API规范处理响应
    if (code === 200 || code === 201) {
      return data
    } else {
      ElMessage.error(message || '请求失败')
      return Promise.reject(new Error(message || '请求失败'))
    }
  },
  (error: AxiosError) => {
    hideLoading()

    // 处理HTTP错误状态码
    if (error.response) {
      const { status, data } = error.response
      let message = '请求失败'

      switch (status) {
        case 401:
          message = '未授权，请重新登录'
          // 可以在这里处理登录过期逻辑
          // router.push('/login')
          break
        case 403:
          message = '拒绝访问'
          break
        case 404:
          message = '请求地址出错'
          break
        case 408:
          message = '请求超时'
          break
        case 500:
          message = '服务器内部错误'
          break
        case 501:
          message = '服务未实现'
          break
        case 502:
          message = '网关错误'
          break
        case 503:
          message = '服务不可用'
          break
        case 504:
          message = '网关超时'
          break
        case 505:
          message = 'HTTP版本不受支持'
          break
        default:
          message = (data as any)?.message || `连接错误${status}`
      }

      ElMessage.error(message)
    } else if (error.request) {
      ElMessage.error('网络连接异常，请检查网络')
    } else {
      ElMessage.error('请求配置错误')
    }

    return Promise.reject(error)
  }
)

export default request

// 导出常用的请求方法
export const get = (url: string, params?: any, config?: any) => {
  return request.get(url, { params, ...config })
}

export const post = (url: string, data?: any, config?: any) => {
  return request.post(url, data, config)
}

export const put = (url: string, data?: any, config?: any) => {
  return request.put(url, data, config)
}

export const del = (url: string, config?: any) => {
  return request.delete(url, config)
}
import { ElMessage, ElNotification } from 'element-plus'
import type { App } from 'vue'

// 错误类型枚举
export enum ErrorType {
  NETWORK = 'NETWORK',
  API = 'API',
  VALIDATION = 'VALIDATION',
  PERMISSION = 'PERMISSION',
  UNKNOWN = 'UNKNOWN'
}

// 错误信息接口
export interface ErrorInfo {
  type: ErrorType
  message: string
  code?: string | number
  details?: any
  timestamp: number
}

// 错误处理器类
class ErrorHandler {
  private errorLog: ErrorInfo[] = []
  private maxLogSize = 100

  // 处理错误
  handleError(error: any, type: ErrorType = ErrorType.UNKNOWN, showMessage = true): ErrorInfo {
    const errorInfo: ErrorInfo = {
      type,
      message: this.extractErrorMessage(error),
      code: error.code || error.status,
      details: error,
      timestamp: Date.now()
    }

    // 记录错误日志
    this.logError(errorInfo)

    // 显示错误消息
    if (showMessage) {
      this.showErrorMessage(errorInfo)
    }

    // 上报错误（可选）
    this.reportError(errorInfo)

    return errorInfo
  }

  // 提取错误消息
  private extractErrorMessage(error: any): string {
    if (typeof error === 'string') {
      return error
    }

    if (error?.message) {
      return error.message
    }

    if (error?.response?.data?.message) {
      return error.response.data.message
    }

    if (error?.response?.statusText) {
      return error.response.statusText
    }

    return '未知错误'
  }

  // 显示错误消息
  private showErrorMessage(errorInfo: ErrorInfo) {
    const { type, message } = errorInfo

    switch (type) {
      case ErrorType.NETWORK:
        ElMessage.error({
          message: `网络错误: ${message}`,
          duration: 5000
        })
        break

      case ErrorType.API:
        ElMessage.error({
          message: `接口错误: ${message}`,
          duration: 4000
        })
        break

      case ErrorType.VALIDATION:
        ElMessage.warning({
          message: `验证错误: ${message}`,
          duration: 3000
        })
        break

      case ErrorType.PERMISSION:
        ElNotification.error({
          title: '权限错误',
          message,
          duration: 5000
        })
        break

      default:
        ElMessage.error({
          message,
          duration: 4000
        })
    }
  }

  // 记录错误日志
  private logError(errorInfo: ErrorInfo) {
    console.error('[ErrorHandler]', errorInfo)

    // 添加到错误日志
    this.errorLog.unshift(errorInfo)

    // 限制日志大小
    if (this.errorLog.length > this.maxLogSize) {
      this.errorLog = this.errorLog.slice(0, this.maxLogSize)
    }

    // 存储到本地存储（可选）
    try {
      localStorage.setItem('errorLog', JSON.stringify(this.errorLog.slice(0, 10)))
    } catch (e) {
      console.warn('无法保存错误日志到本地存储')
    }
  }

  // 上报错误（可选实现）
  private reportError(errorInfo: ErrorInfo) {
    // 这里可以实现错误上报逻辑
    // 比如发送到错误监控服务
    if (import.meta.env.PROD) {
      // 生产环境才上报
      // sendErrorToService(errorInfo)
    }
  }

  // 获取错误日志
  getErrorLog(): ErrorInfo[] {
    return [...this.errorLog]
  }

  // 清除错误日志
  clearErrorLog() {
    this.errorLog = []
    localStorage.removeItem('errorLog')
  }

  // 处理Promise拒绝
  handlePromiseRejection(event: PromiseRejectionEvent) {
    this.handleError(event.reason, ErrorType.UNKNOWN)
    event.preventDefault()
  }

  // 处理全局错误
  handleGlobalError(event: ErrorEvent) {
    this.handleError({
      message: event.message,
      filename: event.filename,
      lineno: event.lineno,
      colno: event.colno,
      error: event.error
    }, ErrorType.UNKNOWN)
  }
}

// 创建全局错误处理器实例
export const errorHandler = new ErrorHandler()

// Vue插件安装函数
export function installErrorHandler(app: App) {
  // 设置全局错误处理
  app.config.errorHandler = (err, instance, info) => {
    console.error('Vue Error:', err, info)
    errorHandler.handleError(err, ErrorType.UNKNOWN)
  }

  // 监听未捕获的Promise拒绝
  window.addEventListener('unhandledrejection', (event) => {
    errorHandler.handlePromiseRejection(event)
  })

  // 监听全局错误
  window.addEventListener('error', (event) => {
    errorHandler.handleGlobalError(event)
  })

  // 将错误处理器添加到全局属性
  app.config.globalProperties.$errorHandler = errorHandler
}

// 导出便捷方法
export const handleNetworkError = (error: any) => errorHandler.handleError(error, ErrorType.NETWORK)
export const handleApiError = (error: any) => errorHandler.handleError(error, ErrorType.API)
export const handleValidationError = (error: any) => errorHandler.handleError(error, ErrorType.VALIDATION)
export const handlePermissionError = (error: any) => errorHandler.handleError(error, ErrorType.PERMISSION)
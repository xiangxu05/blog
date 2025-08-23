import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export type Theme = 'light' | 'dark'

export interface BreadcrumbItem {
  title: string
  path?: string
}

export const useAppStore = defineStore('app', () => {
  // 状态
  const theme = ref<Theme>('light')
  const sidebarCollapsed = ref(false)
  const globalLoading = ref(false)
  const pageTitle = ref('个人博客')
  const isMobile = ref(false)
  
  // 面包屑导航
  const breadcrumbs = ref<BreadcrumbItem[]>([{ title: '首页', path: '/' }])

  // 切换主题
  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
    localStorage.setItem('theme', theme.value)
    updateThemeClass()
  }

  // 设置主题
  const setTheme = (newTheme: Theme) => {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
    updateThemeClass()
  }

  // 更新主题类名
  const updateThemeClass = () => {
    const html = document.documentElement
    if (theme.value === 'dark') {
      html.classList.add('dark')
    } else {
      html.classList.remove('dark')
    }
  }

  // 初始化主题
  const initTheme = () => {
    const savedTheme = localStorage.getItem('theme') as Theme
    if (savedTheme && ['light', 'dark'].includes(savedTheme)) {
      theme.value = savedTheme
    } else {
      // 检测系统主题偏好
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      theme.value = prefersDark ? 'dark' : 'light'
    }
    updateThemeClass()
  }

  // 切换侧边栏
  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value
    localStorage.setItem('sidebarCollapsed', String(sidebarCollapsed.value))
  }

  // 设置侧边栏状态
  const setSidebarCollapsed = (collapsed: boolean) => {
    sidebarCollapsed.value = collapsed
    localStorage.setItem('sidebarCollapsed', String(collapsed))
  }

  // 初始化侧边栏状态
  const initSidebar = () => {
    const saved = localStorage.getItem('sidebarCollapsed')
    if (saved !== null) {
      sidebarCollapsed.value = saved === 'true'
    }
  }

  // 设置页面标题
  const setPageTitle = (title: string) => {
    pageTitle.value = title
    document.title = `${title} - ${import.meta.env.VITE_APP_TITLE}`
  }

  // 检测移动端
  const checkMobile = () => {
    isMobile.value = window.innerWidth < 1024
    // 移动端默认收起侧边栏
    if (isMobile.value) {
      sidebarCollapsed.value = true
    }
  }

  // 初始化移动端检测
  const initMobileDetection = () => {
    checkMobile()
    window.addEventListener('resize', checkMobile)
  }

  // 设置面包屑
  const setBreadcrumbs = (crumbs: BreadcrumbItem[]) => {
    breadcrumbs.value = crumbs
  }

  // 添加面包屑
  const addBreadcrumb = (crumb: BreadcrumbItem) => {
    breadcrumbs.value.push(crumb)
  }

  // 重置面包屑
  const resetBreadcrumbs = () => {
    breadcrumbs.value = [{ title: '首页', path: '/' }]
  }

  // 设置全局加载状态
  const setGlobalLoading = (isLoading: boolean) => {
    globalLoading.value = isLoading
  }

  // 初始化应用状态
  const initApp = () => {
    initTheme()
    initSidebar()
    initMobileDetection()
  }

  // 销毁方法（清理事件监听器）
  const destroy = () => {
    window.removeEventListener('resize', checkMobile)
  }

  return {
    // 状态
    theme,
    sidebarCollapsed,
    globalLoading,
    pageTitle,
    breadcrumbs,
    isMobile,
    
    // 方法
    toggleTheme,
    setTheme,
    toggleSidebar,
    setSidebarCollapsed,
    setPageTitle,
    setBreadcrumbs,
    addBreadcrumb,
    resetBreadcrumbs,
    setGlobalLoading,
    initApp,
    destroy
  }
})
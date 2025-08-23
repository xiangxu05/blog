import './assets/main.css'
import './styles/global.scss'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'

import App from './App.vue'
import router from './router'
import { installErrorHandler } from './utils/errorHandler'
import { useAppStore } from './stores/app'
import { useUserStore } from './stores/user'

const app = createApp(App)
const pinia = createPinia()

// 安装插件
app.use(pinia)
app.use(router)
app.use(ElementPlus)

// 安装错误处理器
installErrorHandler(app)

// 初始化应用状态
const appStore = useAppStore()
const userStore = useUserStore()

// 初始化应用
appStore.initApp()

// 尝试获取用户信息（如果已登录）
userStore.fetchUserInfo().catch(() => {
  // 忽略错误，用户未登录是正常情况
})

app.mount('#app')

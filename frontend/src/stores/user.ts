import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/api'
import { ElMessage } from 'element-plus'

// 用户信息接口
export interface UserInfo {
  id: number
  name: string
  email: string
  avatar?: string
  role: string
  bio?: string
  website?: string
  github?: string
  created_at?: string
}

// 登录响应接口
interface LoginResponse {
  token: string
  user: UserInfo
}

export const useUserStore = defineStore('user', () => {
  // 状态
  const userInfo = ref<UserInfo | null>(null)
  const token = ref<string | null>(null)
  const isLoggedIn = ref(false)
  const loading = ref(false)

  // 计算属性
  const isAdmin = computed(() => {
    return userInfo.value?.role === 'admin'
  })

  const avatar = computed(() => {
    return userInfo.value?.avatar || '/default-avatar.png'
  })

  const displayName = computed(() => {
    return userInfo.value?.name || '未知用户'
  })

  // 设置用户信息
  const setUser = (user: UserInfo) => {
    userInfo.value = user
    isLoggedIn.value = true
  }

  // 设置 token
  const setToken = (newToken: string) => {
    token.value = newToken
    // 可以在这里设置到 localStorage 或其他持久化存储
    localStorage.setItem('auth_token', newToken)
  }

  // 登录
  const login = async (email: string, password: string, remember = false) => {
    try {
      loading.value = true
      const response = await api.post('/users/login', {
        email,
        password,
        remember
      })

      const { token: authToken, user } = response.data
      setUser(user)
      setToken(authToken)

      ElMessage.success('登录成功')
      return response.data
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // 注册
  const register = async (userData: {
    name: string
    email: string
    password: string
    confirmPassword: string
  }) => {
    try {
      loading.value = true
      await api.post('/users/register', userData)
      ElMessage.success('注册成功，请登录')
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取用户信息
  const fetchUserInfo = async () => {
    try {
      const response = await api.get('/users/profile')
      const user = response.data
      setUser(user)
      return user
    } catch (error) {
      // 如果获取用户信息失败，说明未登录
      userInfo.value = null
      isLoggedIn.value = false
      token.value = null
      localStorage.removeItem('auth_token')
      throw error
    }
  }

  // 更新用户信息
  const updateProfile = async (data: {
    name?: string
    email?: string
    avatar?: string
    bio?: string
    website?: string
    github?: string
  }) => {
    try {
      loading.value = true
      await api.put('/users/profile', data)

      // 更新本地用户信息
      if (userInfo.value) {
        Object.assign(userInfo.value, data)
      }

      ElMessage.success('更新成功')
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // 修改密码
  const changePassword = async (data: {
    oldPassword: string
    newPassword: string
    confirmPassword: string
  }) => {
    try {
      loading.value = true
      await api.put('/users/password', data)
      ElMessage.success('密码修改成功')
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // 注销
  const logout = async () => {
    try {
      await api.get('/users/logout')
    } catch (error) {
      // 即使注销接口失败，也要清除本地状态
      console.error('注销接口调用失败:', error)
    } finally {
      // 清除本地状态
      userInfo.value = null
      isLoggedIn.value = false
      token.value = null
      localStorage.removeItem('auth_token')
      ElMessage.success('已退出登录')
    }
  }

  // 刷新会话
  const refreshSession = async () => {
    try {
      await fetchUserInfo()
    } catch (error) {
      // 会话已过期，清除本地状态
      userInfo.value = null
      isLoggedIn.value = false
      token.value = null
      localStorage.removeItem('auth_token')
    }
  }

  return {
    // 状态
    userInfo,
    isLoggedIn,
    loading,
    token,

    // 计算属性
    isAdmin,
    avatar,
    displayName,

    // 方法
    setUser,
    setToken,
    login,
    register,
    fetchUserInfo,
    updateProfile,
    changePassword,
    logout,
    refreshSession
  }
})
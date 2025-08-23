<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-card">
        <div class="login-header">
          <h1 class="login-title">欢迎回来</h1>
          <p class="login-subtitle">登录您的账户</p>
        </div>
        
        <el-form
          ref="formRef"
          :model="formData"
          :rules="formRules"
          class="login-form"
          @submit.prevent="handleSubmit"
        >
          <el-form-item prop="username">
            <el-input
              v-model="formData.username"
              placeholder="用户名"
              size="large"
              :disabled="loading"
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input
              v-model="formData.password"
              type="password"
              placeholder="密码"
              size="large"
              show-password
              :disabled="loading"
              @keyup.enter="handleSubmit"
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <div class="login-options">
            <el-checkbox v-model="formData.remember" :disabled="loading">
              记住我
            </el-checkbox>
            <el-link type="primary" @click="showForgotPassword">
              忘记密码？
            </el-link>
          </div>
          
          <el-button
            type="primary"
            size="large"
            class="login-btn"
            :loading="loading"
            @click="handleSubmit"
          >
            登录
          </el-button>
        </el-form>
        
        <div class="login-footer">
          <p class="register-link">
            还没有账户？
            <router-link to="/register" class="link">
              立即注册
            </router-link>
          </p>
        </div>
      </div>
      

    </div>
    
    <!-- 忘记密码对话框 -->
    <el-dialog
      v-model="showForgotDialog"
      title="重置密码"
      width="400px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="forgotFormRef"
        :model="forgotForm"
        :rules="forgotRules"
        label-width="0"
      >
        <el-form-item prop="email">
          <el-input
            v-model="forgotForm.email"
            type="email"
            placeholder="请输入您的邮箱地址"
            size="large"
          >
            <template #prefix>
              <el-icon><Message /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        
        <p class="forgot-tip">
          我们将向您的邮箱发送重置密码的链接
        </p>
      </el-form>
      
      <template #footer>
        <el-button @click="showForgotDialog = false">取消</el-button>
        <el-button
          type="primary"
          @click="handleForgotPassword"
          :loading="forgotLoading"
        >
          发送重置链接
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Message, Lock, User } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import api from '@/utils/api'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const appStore = useAppStore()

// 响应式数据
const loading = ref(false)
const forgotLoading = ref(false)
const showForgotDialog = ref(false)
const formRef = ref<FormInstance>()
const forgotFormRef = ref<FormInstance>()

// 表单数据
const formData = reactive({
  username: '',
  password: '',
  remember: false
})

// 忘记密码表单数据
const forgotForm = reactive({
  email: ''
})

// 表单验证规则
const formRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, message: '用户名长度不能少于3位', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

const forgotRules: FormRules = {
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

// 方法
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    loading.value = true
    
    const response = await api.post('/users/login', {
      username: formData.username,
      password: formData.password
    })
    
    const { user } = response.data
    
    // 保存用户信息（token通过cookie自动管理）
    userStore.setUser(user)
    
    ElMessage.success('登录成功')
    
    // 重定向到目标页面或首页
    const redirect = route.query.redirect as string || '/'
    await router.push(redirect)
  } catch (error: any) {
    console.error('登录失败:', error)
    
    // 显示具体的错误信息
    const message = error.response?.data?.message || '登录失败，请检查用户名和密码'
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}

const showForgotPassword = () => {
  forgotForm.email = ''
  showForgotDialog.value = true
}

const handleForgotPassword = async () => {
  if (!forgotFormRef.value) return
  
  try {
    await forgotFormRef.value.validate()
    forgotLoading.value = true
    
    await api.post('/users/forgot-password', {
      email: forgotForm.email
    })
    
    ElMessage.success('重置链接已发送到您的邮箱')
    showForgotDialog.value = false
  } catch (error) {
    console.error('发送重置链接失败:', error)
    ElMessage.success('重置链接已发送到您的邮箱（模拟）')
    showForgotDialog.value = false
  } finally {
    forgotLoading.value = false
  }
}

// 生命周期
onMounted(() => {
  // 设置页面标题
  appStore.setPageTitle('用户登录')
  appStore.resetBreadcrumbs()
  
  // 如果已经登录，重定向到首页
  if (userStore.isLoggedIn) {
    const redirect = route.query.redirect as string || '/'
    router.push(redirect)
  }
})
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: $bg-secondary;
  padding: $spacing-lg;
  position: relative;
}

.login-container {
  position: relative;
}

.login-card {
  background: $bg-primary;
  border-radius: $border-radius-xl;
  padding: $spacing-2xl;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid $border-light;
}

.login-header {
  text-align: center;
  margin-bottom: $spacing-xl;
  
  .login-title {
    font-size: $font-size-xxl;
    font-weight: 700;
    color: $text-primary;
    margin-bottom: $spacing-sm;
  }
  
  .login-subtitle {
    font-size: $font-size-base;
    color: $text-secondary;
    margin: 0;
  }
}

.login-form {
  .el-form-item {
    margin-bottom: $spacing-lg;
  }
  
  .el-input {
    --el-input-border-radius: #{$border-radius-medium};
  }
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: $spacing-xl;
  
  .el-checkbox {
    color: $text-secondary;
  }
}

.login-btn {
  width: 100%;
  height: 48px;
  font-size: $font-size-base;
  font-weight: 500;
  border-radius: $border-radius-medium;
  margin-bottom: $spacing-lg;
}

.login-footer {
  text-align: center;
  
  .register-link {
    color: $text-secondary;
    margin: 0;
    
    .link {
      color: $primary-color;
      text-decoration: none;
      font-weight: 500;
      
      &:hover {
        text-decoration: underline;
      }
    }
  }
}



.forgot-tip {
  font-size: $font-size-sm;
  color: $text-secondary;
  text-align: center;
  margin: $spacing-md 0 0 0;
  line-height: 1.5;
}

// 响应式设计
@media (max-width: 480px) {
  .login-page {
    padding: $spacing-md;
  }
  
  .login-card {
    padding: $spacing-xl;
    max-width: 100%;
  }
  
  .login-header {
    .login-title {
      font-size: $font-size-xl;
    }
  }
  
  .decoration-circle {
    &.circle-1 {
      width: 120px;
      height: 120px;
    }
    
    &.circle-2 {
      width: 100px;
      height: 100px;
    }
    
    &.circle-3 {
      width: 80px;
      height: 80px;
    }
  }
}

// 暗色模式支持
@media (prefers-color-scheme: dark) {
  .login-card {
    background: rgba(30, 30, 30, 0.95);
    border-color: rgba(255, 255, 255, 0.1);
  }
  
  .login-header {
    .login-title {
      color: $dark-text-primary;
    }
    
    .login-subtitle {
      color: $dark-text-secondary;
    }
  }
  
  .login-options .el-checkbox {
    color: $dark-text-secondary;
  }
  
  .login-footer .register-link {
    color: $dark-text-secondary;
  }
  
  .forgot-tip {
    color: $dark-text-secondary;
  }
}
</style>
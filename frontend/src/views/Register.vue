<template>
  <div class="register-page">
    <div class="register-container">
      <div class="register-card">
        <div class="register-header">
          <h1 class="register-title">创建账户</h1>
          <p class="register-subtitle">加入我们的社区</p>
        </div>
        
        <el-form
          ref="formRef"
          :model="formData"
          :rules="formRules"
          class="register-form"
          @submit.prevent="handleSubmit"
        >
          <el-form-item prop="name">
            <el-input
              v-model="formData.name"
              placeholder="用户名"
              size="large"
              :disabled="loading"
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item prop="email">
            <el-input
              v-model="formData.email"
              type="email"
              placeholder="邮箱地址"
              size="large"
              :disabled="loading"
            >
              <template #prefix>
                <el-icon><Message /></el-icon>
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
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item prop="confirmPassword">
            <el-input
              v-model="formData.confirmPassword"
              type="password"
              placeholder="确认密码"
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
          
          <div class="register-options">
            <el-checkbox v-model="formData.agreeTerms" :disabled="loading">
              我已阅读并同意
              <el-link type="primary" @click="showTermsDialog = true">
                服务条款
              </el-link>
              和
              <el-link type="primary" @click="showPrivacyDialog = true">
                隐私政策
              </el-link>
            </el-checkbox>
          </div>
          
          <el-button
            type="primary"
            size="large"
            class="register-btn"
            :loading="loading"
            :disabled="!formData.agreeTerms"
            @click="handleSubmit"
          >
            注册账户
          </el-button>
        </el-form>
        
        <div class="register-footer">
          <p class="login-link">
            已有账户？
            <router-link to="/login" class="link">
              立即登录
            </router-link>
          </p>
        </div>
      </div>
      
      <!-- 装饰性背景 -->
      <div class="background-decoration">
        <div class="decoration-circle circle-1"></div>
        <div class="decoration-circle circle-2"></div>
        <div class="decoration-circle circle-3"></div>
      </div>
    </div>
    
    <!-- 服务条款对话框 -->
    <el-dialog
      v-model="showTermsDialog"
      title="服务条款"
      width="600px"
      :close-on-click-modal="false"
    >
      <div class="terms-content">
        <h3>1. 服务说明</h3>
        <p>本网站是一个个人博客平台，为用户提供文章阅读、评论等服务。</p>
        
        <h3>2. 用户责任</h3>
        <p>用户应当遵守相关法律法规，不得发布违法、有害信息。</p>
        
        <h3>3. 隐私保护</h3>
        <p>我们承诺保护用户隐私，不会泄露用户个人信息。</p>
        
        <h3>4. 免责声明</h3>
        <p>本网站对用户发布的内容不承担责任，如有侵权请及时联系我们。</p>
        
        <h3>5. 条款变更</h3>
        <p>我们保留随时修改本条款的权利，修改后的条款将在网站上公布。</p>
      </div>
      
      <template #footer>
        <el-button type="primary" @click="showTermsDialog = false">
          我已了解
        </el-button>
      </template>
    </el-dialog>
    
    <!-- 隐私政策对话框 -->
    <el-dialog
      v-model="showPrivacyDialog"
      title="隐私政策"
      width="600px"
      :close-on-click-modal="false"
    >
      <div class="privacy-content">
        <h3>1. 信息收集</h3>
        <p>我们只收集必要的用户信息，包括用户名、邮箱等基本信息。</p>
        
        <h3>2. 信息使用</h3>
        <p>收集的信息仅用于提供服务、改善用户体验，不会用于其他商业目的。</p>
        
        <h3>3. 信息保护</h3>
        <p>我们采用安全措施保护用户信息，防止未经授权的访问、使用或泄露。</p>
        
        <h3>4. 信息共享</h3>
        <p>除法律要求外，我们不会与第三方共享用户个人信息。</p>
        
        <h3>5. Cookie使用</h3>
        <p>我们使用Cookie来改善用户体验，用户可以通过浏览器设置管理Cookie。</p>
      </div>
      
      <template #footer>
        <el-button type="primary" @click="showPrivacyDialog = false">
          我已了解
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Message, Lock } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import api from '@/utils/api'

const router = useRouter()
const userStore = useUserStore()
const appStore = useAppStore()

// 响应式数据
const loading = ref(false)
const showTermsDialog = ref(false)
const showPrivacyDialog = ref(false)
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: '',
  agreeTerms: false
})

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '用户名长度在 2 到 20 个字符', trigger: 'blur' },
    {
      pattern: /^[a-zA-Z0-9_\u4e00-\u9fa5]+$/,
      message: '用户名只能包含字母、数字、下划线和中文',
      trigger: 'blur'
    }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' },
    {
      pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d@$!%*?&]{6,}$/,
      message: '密码必须包含大小写字母和数字',
      trigger: 'blur'
    }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== formData.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 方法
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    if (!formData.agreeTerms) {
      ElMessage.warning('请先同意服务条款和隐私政策')
      return
    }
    
    loading.value = true
    
    const response = await api.post('/users/register', {
      username: formData.name,
      email: formData.email,
      password: formData.password
    })
    
    const { user } = response.data
    
    // 保存用户信息（token通过cookie自动管理）
    userStore.setUser(user)
    
    ElMessage.success('注册成功，欢迎加入！')
    
    // 重定向到首页
    await router.push('/')
  } catch (error: any) {
    console.error('注册失败:', error)
    
    // 检查是否是邮箱已存在的错误
    if (error.response?.status === 409) {
      ElMessage.error('该邮箱已被注册，请使用其他邮箱')
      return
    }
    
    // 显示具体的错误信息
    const message = error.response?.data?.message || '注册失败，请稍后重试'
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}

// 生命周期
onMounted(() => {
  // 设置页面标题
  appStore.setPageTitle('用户注册')
  appStore.resetBreadcrumbs()
  
  // 如果已经登录，重定向到首页
  if (userStore.isLoggedIn) {
    router.push('/')
  }
})
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: $spacing-lg;
  position: relative;
  overflow: hidden;
}

.register-container {
  position: relative;
  z-index: 2;
}

.register-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: $border-radius-xl;
  padding: $spacing-2xl;
  width: 100%;
  max-width: 450px;
  box-shadow: $shadow-lg;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.register-header {
  text-align: center;
  margin-bottom: $spacing-xl;
  
  .register-title {
    font-size: $font-size-xxl;
    font-weight: 700;
    color: $text-primary;
    margin-bottom: $spacing-sm;
  }
  
  .register-subtitle {
    font-size: $font-size-base;
    color: $text-secondary;
    margin: 0;
  }
}

.register-form {
  .el-form-item {
    margin-bottom: $spacing-lg;
  }
  
  .el-input {
    --el-input-border-radius: #{$border-radius-medium};
  }
}

.register-options {
  margin-bottom: $spacing-xl;
  
  .el-checkbox {
    color: $text-secondary;
    line-height: 1.5;
    
    :deep(.el-checkbox__label) {
      font-size: $font-size-sm;
    }
  }
}

.register-btn {
  width: 100%;
  height: 48px;
  font-size: $font-size-base;
  font-weight: 500;
  border-radius: $border-radius-medium;
  margin-bottom: $spacing-lg;
  
  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}

.register-footer {
  text-align: center;
  
  .login-link {
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

.background-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1;
  overflow: hidden;
}

.decoration-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 6s ease-in-out infinite;
  
  &.circle-1 {
    width: 200px;
    height: 200px;
    top: 10%;
    left: 10%;
    animation-delay: 0s;
  }
  
  &.circle-2 {
    width: 150px;
    height: 150px;
    top: 60%;
    right: 10%;
    animation-delay: 2s;
  }
  
  &.circle-3 {
    width: 100px;
    height: 100px;
    bottom: 20%;
    left: 20%;
    animation-delay: 4s;
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
  }
}

.terms-content,
.privacy-content {
  max-height: 400px;
  overflow-y: auto;
  padding: $spacing-md;
  
  h3 {
    font-size: $font-size-lg;
    font-weight: 600;
    color: $text-primary;
    margin: $spacing-lg 0 $spacing-sm 0;
    
    &:first-child {
      margin-top: 0;
    }
  }
  
  p {
    font-size: $font-size-base;
    color: $text-secondary;
    line-height: 1.6;
    margin-bottom: $spacing-md;
  }
}

// 响应式设计
@media (max-width: 480px) {
  .register-page {
    padding: $spacing-md;
  }
  
  .register-card {
    padding: $spacing-xl;
    max-width: 100%;
  }
  
  .register-header {
    .register-title {
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
  .register-card {
    background: rgba(30, 30, 30, 0.95);
    border-color: rgba(255, 255, 255, 0.1);
  }
  
  .register-header {
    .register-title {
      color: $dark-text-primary;
    }
    
    .register-subtitle {
      color: $dark-text-secondary;
    }
  }
  
  .register-options .el-checkbox {
    color: $dark-text-secondary;
  }
  
  .register-footer .login-link {
    color: $dark-text-secondary;
  }
  
  .terms-content,
  .privacy-content {
    h3 {
      color: $dark-text-primary;
    }
    
    p {
      color: $dark-text-secondary;
    }
  }
}
</style>
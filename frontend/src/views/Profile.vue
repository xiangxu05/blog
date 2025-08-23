<template>
  <div class="profile-page">
    <div class="page-header">
      <h1 class="page-title">个人资料</h1>
      <p class="page-description">管理您的个人信息</p>
    </div>
    
    <div class="profile-content">
      <div class="profile-card">
        <div class="avatar-section">
          <div class="avatar-container">
            <el-avatar
              :size="120"
              :src="userInfo.avatar"
              :icon="UserFilled"
              class="user-avatar"
            />
            <el-button
              type="primary"
              size="small"
              circle
              class="avatar-edit-btn"
              @click="handleAvatarEdit"
            >
              <el-icon><Edit /></el-icon>
            </el-button>
          </div>
          
          <div class="user-basic-info">
            <h2 class="username">{{ userInfo.name || '未设置' }}</h2>
            <p class="user-email">{{ userInfo.email }}</p>
            <el-tag :type="userInfo.role === 'admin' ? 'danger' : 'info'" size="small">
              {{ userInfo.role === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
          </div>
        </div>
        
        <el-divider />
        
        <div class="profile-form">
          <el-form
            ref="formRef"
            :model="formData"
            :rules="formRules"
            label-width="100px"
            label-position="left"
          >
            <el-form-item label="用户名" prop="name">
              <el-input
                v-model="formData.name"
                placeholder="请输入用户名"
                :disabled="loading"
              />
            </el-form-item>
            
            <el-form-item label="邮箱" prop="email">
              <el-input
                v-model="formData.email"
                type="email"
                placeholder="请输入邮箱"
                :disabled="loading"
              />
            </el-form-item>
            
            <el-form-item label="个人简介" prop="bio">
              <el-input
                v-model="formData.bio"
                type="textarea"
                :rows="4"
                placeholder="介绍一下自己吧..."
                :disabled="loading"
                maxlength="200"
                show-word-limit
              />
            </el-form-item>
            
            <el-form-item label="个人网站" prop="website">
              <el-input
                v-model="formData.website"
                placeholder="https://example.com"
                :disabled="loading"
              >
                <template #prefix>
                  <el-icon><Link /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item label="GitHub" prop="github">
              <el-input
                v-model="formData.github"
                placeholder="GitHub 用户名"
                :disabled="loading"
              >
                <template #prefix>
                  <el-icon><Platform /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item>
              <el-button
                type="primary"
                @click="handleSubmit"
                :loading="loading"
              >
                保存修改
              </el-button>
              <el-button @click="handleReset" :disabled="loading">
                重置
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
      
      <div class="security-card">
        <h3 class="card-title">
          <el-icon><Lock /></el-icon>
          安全设置
        </h3>
        
        <div class="security-items">
          <div class="security-item">
            <div class="item-info">
              <h4>修改密码</h4>
              <p>定期修改密码以保护账户安全</p>
            </div>
            <el-button type="primary" plain @click="showPasswordDialog = true">
              修改密码
            </el-button>
          </div>
          
          <el-divider />
          
          <div class="security-item">
            <div class="item-info">
              <h4>登录记录</h4>
              <p>查看最近的登录活动</p>
            </div>
            <el-button type="info" plain @click="showLoginHistory">
              查看记录
            </el-button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 修改密码对话框 -->
    <el-dialog
      v-model="showPasswordDialog"
      title="修改密码"
      width="400px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="100px"
      >
        <el-form-item label="当前密码" prop="currentPassword">
          <el-input
            v-model="passwordForm.currentPassword"
            type="password"
            placeholder="请输入当前密码"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="passwordForm.newPassword"
            type="password"
            placeholder="请输入新密码"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            show-password
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="showPasswordDialog = false">取消</el-button>
        <el-button
          type="primary"
          @click="handlePasswordChange"
          :loading="passwordLoading"
        >
          确认修改
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import {
  UserFilled,
  Edit,
  Link,
  Platform,
  Lock
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import api from '@/utils/api'

const userStore = useUserStore()

// 响应式数据
const loading = ref(false)
const passwordLoading = ref(false)
const showPasswordDialog = ref(false)
const formRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()

// 用户信息
const userInfo = ref({
  id: 1,
  name: '用户名',
  email: 'user@example.com',
  role: 'user',
  avatar: '',
  bio: '',
  website: '',
  github: ''
})

// 表单数据
const formData = reactive({
  name: '',
  email: '',
  bio: '',
  website: '',
  github: ''
})

// 密码表单数据
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '用户名长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  website: [
    { type: 'url', message: '请输入正确的网址格式', trigger: 'blur' }
  ]
}

const passwordRules: FormRules = {
  currentPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
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
const fetchUserInfo = async () => {
  try {
    loading.value = true
    const response = await api.get('/users/profile')
    userInfo.value = response.data
    
    // 更新表单数据
    Object.assign(formData, {
      name: userInfo.value.name || '',
      email: userInfo.value.email || '',
      bio: userInfo.value.bio || '',
      website: userInfo.value.website || '',
      github: userInfo.value.github || ''
    })
  } catch (error) {
    console.error('获取用户信息失败:', error)
    ElMessage.error('获取用户信息失败')
    
    // 使用模拟数据
    userInfo.value = {
      id: 1,
      name: '博主',
      email: 'admin@example.com',
      role: 'admin',
      avatar: '',
      bio: '这是一个热爱技术的开发者',
      website: 'https://example.com',
      github: 'username'
    }
    
    Object.assign(formData, {
      name: userInfo.value.name,
      email: userInfo.value.email,
      bio: userInfo.value.bio,
      website: userInfo.value.website,
      github: userInfo.value.github
    })
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    loading.value = true
    
    await api.put('/users/profile', formData)
    
    // 更新用户信息
    Object.assign(userInfo.value, formData)
    
    // 更新 store 中的用户信息
    await userStore.updateProfile(formData)
    
    ElMessage.success('个人信息更新成功')
  } catch (error) {
    console.error('更新个人信息失败:', error)
    ElMessage.error('更新个人信息失败')
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  Object.assign(formData, {
    name: userInfo.value.name || '',
    email: userInfo.value.email || '',
    bio: userInfo.value.bio || '',
    website: userInfo.value.website || '',
    github: userInfo.value.github || ''
  })
  
  formRef.value?.clearValidate()
}

const handleAvatarEdit = () => {
  ElMessageBox.prompt('请输入头像链接', '修改头像', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputValue: userInfo.value.avatar,
    inputPlaceholder: 'https://example.com/avatar.jpg'
  }).then(async ({ value }) => {
    try {
      await api.put('/users/avatar', { avatar: value })
      userInfo.value.avatar = value
      ElMessage.success('头像更新成功')
    } catch (error) {
      console.error('更新头像失败:', error)
      ElMessage.error('更新头像失败')
    }
  }).catch(() => {
    // 用户取消
  })
}

const handlePasswordChange = async () => {
  if (!passwordFormRef.value) return
  
  try {
    await passwordFormRef.value.validate()
    passwordLoading.value = true
    
    await api.put('/users/password', {
      old_password: passwordForm.currentPassword,
      new_password: passwordForm.newPassword,
      confirm_password: passwordForm.confirmPassword
    })
    
    ElMessage.success('密码修改成功')
    showPasswordDialog.value = false
    
    // 重置表单
    Object.assign(passwordForm, {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    })
    passwordFormRef.value?.resetFields()
  } catch (error) {
    console.error('修改密码失败:', error)
    ElMessage.error('修改密码失败')
  } finally {
    passwordLoading.value = false
  }
}

const showLoginHistory = () => {
  ElMessage.info('登录记录功能开发中...')
}

// 生命周期
onMounted(() => {
  fetchUserInfo()
})
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.profile-page {
  max-width: 800px;
  margin: 0 auto;
  padding: $spacing-lg;
}

.page-header {
  text-align: center;
  margin-bottom: $spacing-xl;
  
  .page-title {
    font-size: $font-size-xxl;
    font-weight: 700;
    color: $text-primary;
    margin-bottom: $spacing-sm;
  }
  
  .page-description {
    font-size: $font-size-lg;
    color: $text-secondary;
  }
}

.profile-content {
  display: flex;
  flex-direction: column;
  gap: $spacing-xl;
}

.profile-card,
.security-card {
  background: $bg-secondary;
  border-radius: $border-radius-large;
  padding: $spacing-xl;
  box-shadow: $shadow-sm;
  border: 1px solid $border-light;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: $spacing-xl;
  margin-bottom: $spacing-lg;
  
  .avatar-container {
    position: relative;
    
    .user-avatar {
      border: 4px solid $border-light;
    }
    
    .avatar-edit-btn {
      position: absolute;
      bottom: 0;
      right: 0;
      width: 32px;
      height: 32px;
    }
  }
  
  .user-basic-info {
    flex: 1;
    
    .username {
      font-size: $font-size-xl;
      font-weight: 600;
      color: $text-primary;
      margin-bottom: $spacing-xs;
    }
    
    .user-email {
      font-size: $font-size-base;
      color: $text-secondary;
      margin-bottom: $spacing-sm;
    }
  }
}

.profile-form {
  .el-form-item {
    margin-bottom: $spacing-lg;
  }
}

.card-title {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  font-size: $font-size-lg;
  font-weight: 600;
  color: $text-primary;
  margin-bottom: $spacing-lg;
}

.security-items {
  .security-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: $spacing-md 0;
    
    .item-info {
      flex: 1;
      
      h4 {
        font-size: $font-size-base;
        font-weight: 500;
        color: $text-primary;
        margin-bottom: $spacing-xs;
      }
      
      p {
        font-size: $font-size-sm;
        color: $text-secondary;
        margin: 0;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .profile-page {
    padding: $spacing-md;
  }
  
  .profile-card,
  .security-card {
    padding: $spacing-lg;
  }
  
  .avatar-section {
    flex-direction: column;
    text-align: center;
    gap: $spacing-lg;
  }
  
  .security-item {
    flex-direction: column;
    align-items: flex-start;
    gap: $spacing-md;
    
    .item-info {
      text-align: left;
    }
  }
}

// 暗色模式支持
@media (prefers-color-scheme: dark) {
  .page-title {
    color: $dark-text-primary;
  }
  
  .page-description {
    color: $dark-text-secondary;
  }
  
  .profile-card,
  .security-card {
    background: $dark-bg-secondary;
    border-color: $dark-border-light;
  }
  
  .avatar-section {
    .user-avatar {
      border-color: $dark-border-light;
    }
    
    .user-basic-info {
      .username {
        color: $dark-text-primary;
      }
      
      .user-email {
        color: $dark-text-secondary;
      }
    }
  }
  
  .card-title {
    color: $dark-text-primary;
  }
  
  .security-items .security-item .item-info {
    h4 {
      color: $dark-text-primary;
    }
    
    p {
      color: $dark-text-secondary;
    }
  }
}
</style>
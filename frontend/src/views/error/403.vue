<template>
  <div class="error-page">
    <div class="error-container">
      <div class="error-icon">
        <el-icon :size="120" color="#f56c6c">
          <Lock />
        </el-icon>
      </div>
      <h1 class="error-title">403</h1>
      <p class="error-message">抱歉，您没有访问权限</p>
      <p class="error-description">
        您当前的账户权限不足以访问此页面，请联系管理员或使用具有相应权限的账户登录。
      </p>
      <div class="error-actions">
        <el-button type="primary" @click="goHome">
          <el-icon><House /></el-icon>
          返回首页
        </el-button>
        <el-button @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回上页
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { House, ArrowLeft, Lock } from '@element-plus/icons-vue'

const router = useRouter()

const goHome = () => {
  router.push('/')
}

const goBack = () => {
  router.go(-1)
}
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.error-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, $bg-primary 0%, #f8f8f8 100%);
  padding: $spacing-lg;
}

.error-container {
  text-align: center;
  max-width: 500px;
  padding: $spacing-xl;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: $border-radius-large;
  box-shadow: $shadow-lg;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.error-icon {
  margin-bottom: $spacing-lg;
  animation: pulse 2s infinite;
}

.error-title {
  font-size: 4rem;
  font-weight: 700;
  color: #f56c6c;
  margin: $spacing-md 0;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
}

.error-message {
  font-size: $font-size-xl;
  font-weight: 600;
  color: $text-primary;
  margin-bottom: $spacing-md;
}

.error-description {
  font-size: $font-size-base;
  color: $text-secondary;
  line-height: $line-height-relaxed;
  margin-bottom: $spacing-xl;
}

.error-actions {
  display: flex;
  gap: $spacing-md;
  justify-content: center;
  flex-wrap: wrap;
  
  .el-button {
    min-width: 120px;
    height: 40px;
    border-radius: $border-radius-medium;
    font-weight: 500;
    transition: all $transition-normal;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: $shadow-md;
    }
  }
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .error-container {
    padding: $spacing-lg;
    margin: $spacing-md;
  }
  
  .error-title {
    font-size: 3rem;
  }
  
  .error-actions {
    flex-direction: column;
    align-items: center;
    
    .el-button {
      width: 100%;
      max-width: 200px;
    }
  }
}

// 暗色模式支持
@media (prefers-color-scheme: dark) {
  .error-page {
    background: linear-gradient(135deg, $dark-bg-primary 0%, #1a1a1a 100%);
  }
  
  .error-container {
    background: rgba(0, 0, 0, 0.8);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }
  
  .error-message {
    color: $dark-text-primary;
  }
  
  .error-description {
    color: $dark-text-secondary;
  }
}
</style>
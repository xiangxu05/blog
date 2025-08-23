<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import {
  Menu as IconMenu,
  Search as IconSearch,
  Moon as IconMoon,
  Sunny as IconSunny,
  User as IconUser,
  Setting as IconSetting,
  SwitchButton as IconLogout,
  Edit as IconEdit
} from '@element-plus/icons-vue'

const router = useRouter()
const appStore = useAppStore()
const userStore = useUserStore()

const searchQuery = ref('')
const showUserMenu = ref(false)

const isDark = computed(() => appStore.theme === 'dark')
const isLoggedIn = computed(() => userStore.isLoggedIn)
const userInfo = computed(() => userStore.userInfo)
const isAdmin = computed(() => userStore.isAdmin)

// 切换主题
const toggleTheme = () => {
  appStore.toggleTheme()
}

// 切换侧边栏
const toggleSidebar = () => {
  appStore.toggleSidebar()
}

// 搜索功能
const handleSearch = () => {
  if (!searchQuery.value.trim()) {
    ElMessage.warning('请输入搜索关键词')
    return
  }
  router.push({
    name: 'Search',
    query: { q: searchQuery.value.trim() }
  })
  searchQuery.value = ''
}

// 用户登出
const handleLogout = async () => {
  try {
    await userStore.logout()
    ElMessage.success('退出登录成功')
    router.push('/')
  } catch (error) {
    console.error('退出登录失败:', error)
  }
}

// 导航到用户资料
const goToProfile = () => {
  router.push('/profile')
  showUserMenu.value = false
}

// 导航到管理后台
const goToAdmin = () => {
  router.push('/admin')
  showUserMenu.value = false
}

// 导航到写文章
const goToWrite = () => {
  router.push('/write')
  showUserMenu.value = false
}
</script>

<template>
  <header class="app-header glass">
    <div class="header-container">
      <!-- 左侧：Logo和菜单按钮 -->
      <div class="header-left">
        <!-- 移动端菜单按钮 -->
        <el-button
          v-if="appStore.isMobile"
          class="menu-btn"
          type="text"
          @click="toggleSidebar"
        >
          <el-icon><IconMenu /></el-icon>
        </el-button>

        <!-- Logo -->
        <router-link to="/" class="logo">
          <h1>个人博客</h1>
        </router-link>
      </div>

      <!-- 中间：搜索框 -->
      <div class="header-center">
        <div class="search-box">
          <el-input
            v-model="searchQuery"
            placeholder="搜索文章..."
            class="search-input"
            @keyup.enter="handleSearch"
          >
            <template #suffix>
              <el-button
                type="text"
                @click="handleSearch"
              >
                <el-icon><IconSearch /></el-icon>
              </el-button>
            </template>
          </el-input>
        </div>
      </div>

      <!-- 右侧：用户操作 -->
      <div class="header-right">
        <!-- 主题切换 -->
        <el-button
          class="theme-btn"
          type="text"
          @click="toggleTheme"
        >
          <el-icon>
            <IconSunny v-if="isDark" />
            <IconMoon v-else />
          </el-icon>
        </el-button>

        <!-- 未登录状态 -->
        <template v-if="!isLoggedIn">
          <el-button
            type="text"
            @click="router.push('/login')"
          >
            登录
          </el-button>
          <el-button
            type="primary"
            size="small"
            @click="router.push('/register')"
          >
            注册
          </el-button>
        </template>

        <!-- 已登录状态 -->
        <template v-else>
          <!-- 写文章按钮（管理员和作者可见） -->
          <el-button
            v-if="isAdmin"
            type="primary"
            size="small"
            @click="goToWrite"
          >
            <el-icon><IconEdit /></el-icon>
            写文章
          </el-button>

          <!-- 用户菜单 -->
          <el-dropdown
            @visible-change="(visible: boolean) => showUserMenu = visible"
          >
            <div class="user-avatar">
              <el-avatar
                :size="32"
                :src="userInfo?.avatar"
              >
                <el-icon><IconUser /></el-icon>
              </el-avatar>
              <span class="username">{{ userInfo?.name }}</span>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="goToProfile">
                  <el-icon><IconUser /></el-icon>
                  个人资料
                </el-dropdown-item>
                <el-dropdown-item
                  v-if="isAdmin"
                  @click="goToAdmin"
                >
                  <el-icon><IconSetting /></el-icon>
                  管理后台
                </el-dropdown-item>
                <el-dropdown-item
                  divided
                  @click="handleLogout"
                >
                  <el-icon><IconLogout /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </div>
    </div>
  </header>
</template>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.app-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  z-index: $z-fixed;
  border-bottom: 1px solid $border-light;
  backdrop-filter: $backdrop-blur;
  -webkit-backdrop-filter: $backdrop-blur;

  .dark & {
    border-bottom-color: $dark-border-light;
  }
}

.header-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 $spacing-lg;

  @media (max-width: $breakpoint-md) {
    padding: 0 $spacing-md;
  }
}

.header-left {
  display: flex;
  align-items: center;
  gap: $spacing-md;

  .menu-btn {
    padding: $spacing-sm;
    
    .el-icon {
      font-size: 20px;
    }
  }

  .logo {
    text-decoration: none;
    color: inherit;
    
    h1 {
      font-size: $font-size-xl;
      font-weight: $font-weight-bold;
      margin: 0;
      background: linear-gradient(135deg, $primary-color, $accent-color);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;

      .dark & {
        background: linear-gradient(135deg, $dark-text-primary, $accent-color);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
      }
    }
  }
}

.header-center {
  flex: 1;
  max-width: 400px;
  margin: 0 $spacing-lg;

  @media (max-width: $breakpoint-md) {
    display: none;
  }

  .search-box {
    .search-input {
      :deep(.el-input__wrapper) {
        border-radius: $border-radius-large;
        background-color: $bg-secondary;
        border: 1px solid transparent;
        transition: all $transition-fast;

        &:hover {
          border-color: $border-medium;
        }

        &.is-focus {
          border-color: $accent-color;
          box-shadow: 0 0 0 2px rgba($accent-color, 0.2);
        }

        .dark & {
          background-color: $dark-bg-tertiary;
          
          &:hover {
            border-color: $dark-border-medium;
          }
        }
      }
    }
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: $spacing-sm;

  .theme-btn {
    padding: $spacing-sm;
    border-radius: $border-radius-medium;
    
    .el-icon {
      font-size: 18px;
    }

    &:hover {
      background-color: $bg-secondary;

      .dark & {
        background-color: $dark-bg-tertiary;
      }
    }
  }

  .user-avatar {
    display: flex;
    align-items: center;
    gap: $spacing-sm;
    padding: $spacing-xs $spacing-sm;
    border-radius: $border-radius-medium;
    cursor: pointer;
    transition: background-color $transition-fast;

    &:hover {
      background-color: $bg-secondary;

      .dark & {
        background-color: $dark-bg-tertiary;
      }
    }

    .username {
      font-size: $font-size-sm;
      font-weight: $font-weight-medium;
      color: $text-primary;

      .dark & {
        color: $dark-text-primary;
      }

      @media (max-width: $breakpoint-sm) {
        display: none;
      }
    }
  }
}

// 移动端适配
@media (max-width: $breakpoint-md) {
  .header-container {
    .header-center {
      display: none;
    }
  }
}
</style>
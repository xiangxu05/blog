<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import {
  House as IconHome,
  Document as IconArticle,
  Collection as IconCategory,
  User as IconUser,
  Setting as IconSetting,
  DataAnalysis as IconAnalytics,
  Edit as IconEdit,
  FolderOpened as IconFolder,
  Star as IconStar,
  Clock as IconClock
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const userStore = useUserStore()

const activeMenu = ref('')

const isCollapsed = computed(() => appStore.sidebarCollapsed)
const isLoggedIn = computed(() => userStore.isLoggedIn)
const isAdmin = computed(() => userStore.isAdmin)

// 菜单项接口
interface MenuItem {
  index: string
  title: string
  icon: any | null
  show: boolean
  isDivider?: boolean
}

// 菜单项配置
const menuItems = computed((): MenuItem[] => {
  const items = [
    {
      index: '/',
      title: '首页',
      icon: IconHome,
      show: true
    },
    {
      index: '/articles',
      title: '文章列表',
      icon: IconArticle,
      show: true
    },
    {
      index: '/categories',
      title: '分类目录',
      icon: IconCategory,
      show: true
    },
    {
      index: '/archives',
      title: '文章归档',
      icon: IconFolder,
      show: true
    },
    {
      index: '/favorites',
      title: '我的收藏',
      icon: IconStar,
      show: isLoggedIn.value
    },
    {
      index: '/recent',
      title: '最近阅读',
      icon: IconClock,
      show: isLoggedIn.value
    }
  ]

  // 管理员菜单
  if (isAdmin.value) {
    items.push(
      {
        index: 'admin-divider',
        title: '管理功能',
        icon: null as any,
        show: true,
        isDivider: true
      } as MenuItem,
      {
        index: '/admin',
        title: '管理后台',
        icon: IconSetting,
        show: true
      },
      {
        index: '/admin/articles',
        title: '文章管理',
        icon: IconEdit,
        show: true
      },
      {
        index: '/admin/analytics',
        title: '数据统计',
        icon: IconAnalytics,
        show: true
      }
    )
  }

  return items.filter(item => item.show)
})

// 监听路由变化，更新活跃菜单
watch(
  () => route.path,
  (newPath) => {
    activeMenu.value = newPath
  },
  { immediate: true }
)

// 菜单点击处理
const handleMenuClick = (index: string) => {
  if (index.includes('divider')) return
  
  router.push(index)
  
  // 移动端点击菜单后收起侧边栏
  if (appStore.isMobile) {
    appStore.toggleSidebar()
  }
}

// 切换侧边栏折叠状态
const toggleCollapse = () => {
  appStore.toggleSidebar()
}
</script>

<template>
  <aside 
    class="app-sidebar"
    :class="{ 
      collapsed: isCollapsed,
      mobile: appStore.isMobile
    }"
  >
    <div class="sidebar-content">
      <!-- 折叠按钮（桌面端） -->
      <div 
        v-if="!appStore.isMobile" 
        class="collapse-btn"
        @click="toggleCollapse"
      >
        <el-icon>
          <component :is="isCollapsed ? 'ArrowRight' : 'ArrowLeft'" />
        </el-icon>
      </div>

      <!-- 导航菜单 -->
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapsed && !appStore.isMobile"
        :unique-opened="true"
        class="sidebar-menu"
        @select="handleMenuClick"
      >
        <template v-for="item in menuItems" :key="item.index">
          <!-- 分割线 -->
          <div 
            v-if="item.isDivider"
            class="menu-divider"
          >
            <span v-if="!isCollapsed || appStore.isMobile" class="divider-text">
              {{ item.title }}
            </span>
          </div>
          
          <!-- 菜单项 -->
          <el-menu-item 
            v-else
            :index="item.index"
            class="menu-item"
          >
            <el-icon v-if="item.icon">
              <component :is="item.icon" />
            </el-icon>
            <template #title>
              <span>{{ item.title }}</span>
            </template>
          </el-menu-item>
        </template>
      </el-menu>

      <!-- 底部信息（展开状态显示） -->
      <div 
        v-if="!isCollapsed || appStore.isMobile" 
        class="sidebar-footer"
      >
        <div class="footer-info">
          <p class="site-info">个人博客系统</p>
          <p class="version">v1.0.0</p>
        </div>
      </div>
    </div>
  </aside>
</template>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.app-sidebar {
  position: fixed;
  top: 64px;
  left: 0;
  bottom: 0;
  width: 260px;
  background-color: $bg-primary;
  border-right: 1px solid $border-light;
  transition: width $transition-normal;
  z-index: $z-sticky;
  overflow: hidden;

  &.collapsed {
    width: 64px;
  }

  &.mobile {
    width: 260px;
    box-shadow: $shadow-lg;
  }

  .dark & {
    background-color: $dark-bg-secondary;
    border-right-color: $dark-border-light;
  }
}

.sidebar-content {
  display: flex;
  flex-direction: column;
  height: 100%;
  position: relative;
}

.collapse-btn {
  position: absolute;
  top: $spacing-md;
  right: -12px;
  width: 24px;
  height: 24px;
  background-color: $bg-primary;
  border: 1px solid $border-light;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 1;
  transition: all $transition-fast;

  &:hover {
    background-color: $bg-secondary;
    transform: scale(1.1);
  }

  .dark & {
    background-color: $dark-bg-secondary;
    border-color: $dark-border-light;

    &:hover {
      background-color: $dark-bg-tertiary;
    }
  }

  .el-icon {
    font-size: 12px;
    color: $text-secondary;

    .dark & {
      color: $dark-text-secondary;
    }
  }
}

.sidebar-menu {
  flex: 1;
  border: none;
  padding: $spacing-md 0;

  :deep(.el-menu-item) {
    height: 48px;
    line-height: 48px;
    margin: 0 $spacing-sm;
    border-radius: $border-radius-medium;
    transition: all $transition-fast;

    &:hover {
      background-color: $bg-secondary !important;

      .dark & {
        background-color: $dark-bg-tertiary !important;
      }
    }

    &.is-active {
      background-color: rgba($accent-color, 0.1) !important;
      color: $accent-color !important;

      .el-icon {
        color: $accent-color !important;
      }
    }

    .el-icon {
      margin-right: $spacing-sm;
      font-size: 18px;
      color: $text-secondary;
      transition: color $transition-fast;

      .dark & {
        color: $dark-text-secondary;
      }
    }

    span {
      font-size: $font-size-sm;
      font-weight: $font-weight-medium;
    }
  }

  // 折叠状态样式
  &.el-menu--collapse {
    :deep(.el-menu-item) {
      padding: 0 20px;
      
      .el-icon {
        margin-right: 0;
      }
    }
  }
}

.menu-divider {
  margin: $spacing-md 0;
  padding: 0 $spacing-md;
  position: relative;

  &::before {
    content: '';
    position: absolute;
    top: 50%;
    left: $spacing-md;
    right: $spacing-md;
    height: 1px;
    background-color: $border-light;
    transform: translateY(-50%);

    .dark & {
      background-color: $dark-border-light;
    }
  }

  .divider-text {
    display: inline-block;
    padding: 0 $spacing-sm;
    background-color: $bg-primary;
    color: $text-tertiary;
    font-size: $font-size-xs;
    font-weight: $font-weight-medium;
    text-transform: uppercase;
    letter-spacing: 0.5px;

    .dark & {
      background-color: $dark-bg-secondary;
      color: $dark-text-tertiary;
    }
  }
}

.sidebar-footer {
  padding: $spacing-md;
  border-top: 1px solid $border-light;
  margin-top: auto;

  .dark & {
    border-top-color: $dark-border-light;
  }

  .footer-info {
    text-align: center;

    .site-info {
      font-size: $font-size-sm;
      font-weight: $font-weight-medium;
      color: $text-primary;
      margin: 0 0 $spacing-xs 0;

      .dark & {
        color: $dark-text-primary;
      }
    }

    .version {
      font-size: $font-size-xs;
      color: $text-tertiary;
      margin: 0;

      .dark & {
        color: $dark-text-tertiary;
      }
    }
  }
}

// 移动端适配
@media (max-width: $breakpoint-lg) {
  .app-sidebar {
    transform: translateX(-100%);
    transition: transform $transition-normal;

    &:not(.collapsed) {
      transform: translateX(0);
    }
  }
}
</style>
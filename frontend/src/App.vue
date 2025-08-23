<script setup lang="ts">
import { computed } from 'vue'
import { RouterView } from 'vue-router'
import { useAppStore } from './stores/app'
import { useUserStore } from './stores/user'
import AppHeader from './components/layout/AppHeader.vue'
import AppSidebar from './components/layout/AppSidebar.vue'
import AppFooter from './components/layout/AppFooter.vue'

const appStore = useAppStore()
const userStore = useUserStore()

const isDark = computed(() => appStore.theme === 'dark')
const sidebarCollapsed = computed(() => appStore.sidebarCollapsed)
</script>

<template>
  <div id="app" :class="{ dark: isDark }">
    <!-- 全局加载遮罩 -->
    <div v-if="appStore.globalLoading" class="global-loading">
      <div class="loading-spinner">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>
    </div>

    <!-- 主布局 -->
    <div class="app-layout">
      <!-- 顶部导航 -->
      <AppHeader />

      <!-- 主内容区域 -->
      <div class="app-main">
        <!-- 侧边栏 -->
        <AppSidebar v-if="!appStore.isMobile" />

        <!-- 内容区域 -->
        <main 
          class="app-content" 
          :class="{ 
            'sidebar-collapsed': sidebarCollapsed && !appStore.isMobile,
            'no-sidebar': appStore.isMobile
          }"
        >
          <div class="content-wrapper">
            <!-- 面包屑导航 -->
            <el-breadcrumb 
              v-if="appStore.breadcrumbs.length > 1" 
              class="breadcrumb"
              separator="/"
            >
              <el-breadcrumb-item 
                v-for="item in appStore.breadcrumbs" 
                :key="item.path"
                :to="item.path"
              >
                {{ item.title }}
              </el-breadcrumb-item>
            </el-breadcrumb>

            <!-- 页面内容 -->
            <div class="page-content">
              <RouterView v-slot="{ Component }">
                <transition name="fade" mode="out-in">
                  <component :is="Component" />
                </transition>
              </RouterView>
            </div>
          </div>
        </main>
      </div>

      <!-- 底部 -->
      <AppFooter />
    </div>

    <!-- 移动端侧边栏遮罩 -->
    <div 
      v-if="appStore.isMobile && !sidebarCollapsed" 
      class="sidebar-overlay"
      @click="appStore.toggleSidebar()"
    />

    <!-- 移动端侧边栏 -->
    <transition name="slide-left">
      <AppSidebar 
        v-if="appStore.isMobile && !sidebarCollapsed" 
        class="mobile-sidebar"
      />
    </transition>
  </div>
</template>

<style lang="scss">
@import './styles/variables.scss';

#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.global-loading {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(4px);
  z-index: $z-modal;
  display: flex;
  align-items: center;
  justify-content: center;

  .dark & {
    background: rgba(0, 0, 0, 0.8);
  }

  .loading-spinner {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;

    .spinner {
      width: 40px;
      height: 40px;
      border: 3px solid rgba($primary-color, 0.3);
      border-top: 3px solid $primary-color;
      border-radius: 50%;
      animation: spin 1s linear infinite;
    }

    p {
       color: $text-primary;
       font-size: $font-size-sm;
       margin: 0;
     }
  }
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.app-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.app-main {
  display: flex;
  flex: 1;
  position: relative;
}

.app-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  transition: margin-left $transition-normal;
  margin-left: 260px;

  &.sidebar-collapsed {
    margin-left: 64px;
  }

  &.no-sidebar {
    margin-left: 0;
  }

  @media (max-width: $breakpoint-lg) {
    margin-left: 0;
  }
}

.content-wrapper {
  flex: 1;
  padding: $spacing-lg;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;

  @media (max-width: $breakpoint-md) {
    padding: $spacing-md;
  }
}

.breadcrumb {
  margin-bottom: $spacing-lg;
  padding: $spacing-sm 0;
  border-bottom: 1px solid $border-light;

  .dark & {
    border-bottom-color: $dark-border-light;
  }
}

.page-content {
  flex: 1;
}

.sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: $z-modal-backdrop;
}

.mobile-sidebar {
  position: fixed;
  top: 64px;
  left: 0;
  bottom: 0;
  z-index: $z-modal;
}

// 过渡动画
.slide-left-enter-active,
.slide-left-leave-active {
  transition: transform $transition-normal;
}

.slide-left-enter-from {
  transform: translateX(-100%);
}

.slide-left-leave-to {
  transform: translateX(-100%);
}
</style>

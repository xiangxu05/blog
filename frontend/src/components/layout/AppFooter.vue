<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

const currentYear = computed(() => new Date().getFullYear())
const isDark = computed(() => appStore.theme === 'dark')

// 社交链接配置
const socialLinks = [
  {
    name: 'GitHub',
    url: 'https://github.com',
    icon: 'github'
  },
  {
    name: 'Twitter',
    url: 'https://twitter.com',
    icon: 'twitter'
  },
  {
    name: 'Email',
    url: 'mailto:contact@example.com',
    icon: 'email'
  }
]

// 友情链接
const friendLinks = [
  {
    name: 'Vue.js',
    url: 'https://vuejs.org'
  },
  {
    name: 'Element Plus',
    url: 'https://element-plus.org'
  },
  {
    name: 'TypeScript',
    url: 'https://www.typescriptlang.org'
  }
]
</script>

<template>
  <footer class="app-footer">
    <div class="footer-container">
      <!-- 主要内容区域 -->
      <div class="footer-main">
        <!-- 网站信息 -->
        <div class="footer-section">
          <h3 class="section-title">关于本站</h3>
          <p class="site-description">
            这是一个基于 Vue 3 + TypeScript 构建的现代化个人博客系统，
            采用苹果风格设计，提供优雅的阅读体验。
          </p>
          <div class="social-links">
            <a 
              v-for="link in socialLinks" 
              :key="link.name"
              :href="link.url"
              :title="link.name"
              class="social-link"
              target="_blank"
              rel="noopener noreferrer"
            >
              <el-icon>
                <component :is="link.icon" />
              </el-icon>
            </a>
          </div>
        </div>

        <!-- 快速链接 -->
        <div class="footer-section">
          <h3 class="section-title">快速链接</h3>
          <ul class="link-list">
            <li><router-link to="/">首页</router-link></li>
            <li><router-link to="/articles">文章列表</router-link></li>
            <li><router-link to="/categories">分类目录</router-link></li>
            <li><router-link to="/archives">文章归档</router-link></li>
          </ul>
        </div>

        <!-- 友情链接 -->
        <div class="footer-section">
          <h3 class="section-title">友情链接</h3>
          <ul class="link-list">
            <li 
              v-for="link in friendLinks" 
              :key="link.name"
            >
              <a 
                :href="link.url" 
                target="_blank" 
                rel="noopener noreferrer"
              >
                {{ link.name }}
              </a>
            </li>
          </ul>
        </div>

        <!-- 技术栈 -->
        <div class="footer-section">
          <h3 class="section-title">技术栈</h3>
          <div class="tech-stack">
            <span class="tech-item">Vue 3</span>
            <span class="tech-item">TypeScript</span>
            <span class="tech-item">Element Plus</span>
            <span class="tech-item">Vite</span>
            <span class="tech-item">Pinia</span>
          </div>
        </div>
      </div>

      <!-- 底部版权信息 -->
      <div class="footer-bottom">
        <div class="copyright">
          <p>
            © {{ currentYear }} 个人博客系统. All rights reserved.
          </p>
          <p class="powered-by">
            Powered by 
            <a href="https://vuejs.org" target="_blank" rel="noopener noreferrer">Vue.js</a>
            & 
            <a href="https://element-plus.org" target="_blank" rel="noopener noreferrer">Element Plus</a>
          </p>
        </div>
        
        <div class="footer-links">
          <router-link to="/privacy">隐私政策</router-link>
          <router-link to="/terms">使用条款</router-link>
          <router-link to="/sitemap">网站地图</router-link>
        </div>
      </div>
    </div>
  </footer>
</template>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.app-footer {
  background-color: $bg-secondary;
  border-top: 1px solid $border-light;
  margin-top: auto;
  padding: $spacing-2xl 0 $spacing-lg 0;

  .dark & {
    background-color: $dark-bg-secondary;
    border-top-color: $dark-border-light;
  }
}

.footer-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 $spacing-lg;

  @media (max-width: $breakpoint-md) {
    padding: 0 $spacing-md;
  }
}

.footer-main {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: $spacing-xl;
  margin-bottom: $spacing-xl;

  @media (max-width: $breakpoint-md) {
    grid-template-columns: 1fr;
    gap: $spacing-lg;
  }
}

.footer-section {
  .section-title {
    font-size: $font-size-lg;
    font-weight: $font-weight-semibold;
    color: $text-primary;
    margin: 0 0 $spacing-md 0;

    .dark & {
      color: $dark-text-primary;
    }
  }

  .site-description {
    font-size: $font-size-sm;
    line-height: $line-height-relaxed;
    color: $text-secondary;
    margin: 0 0 $spacing-md 0;

    .dark & {
      color: $dark-text-secondary;
    }
  }

  .link-list {
    list-style: none;
    padding: 0;
    margin: 0;

    li {
      margin-bottom: $spacing-sm;

      a {
        color: $text-secondary;
        text-decoration: none;
        font-size: $font-size-sm;
        transition: color $transition-fast;

        &:hover {
          color: $accent-color;
        }

        .dark & {
          color: $dark-text-secondary;

          &:hover {
            color: $accent-color;
          }
        }
      }
    }
  }

  .social-links {
    display: flex;
    gap: $spacing-sm;

    .social-link {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 36px;
      height: 36px;
      border-radius: $border-radius-medium;
      background-color: $bg-primary;
      color: $text-secondary;
      text-decoration: none;
      transition: all $transition-fast;
      border: 1px solid $border-light;

      &:hover {
        background-color: $accent-color;
        color: white;
        transform: translateY(-2px);
        box-shadow: $shadow-md;
      }

      .dark & {
        background-color: $dark-bg-tertiary;
        border-color: $dark-border-light;
        color: $dark-text-secondary;

        &:hover {
          background-color: $accent-color;
          color: white;
        }
      }

      .el-icon {
        font-size: 16px;
      }
    }
  }

  .tech-stack {
    display: flex;
    flex-wrap: wrap;
    gap: $spacing-sm;

    .tech-item {
      display: inline-block;
      padding: $spacing-xs $spacing-sm;
      background-color: $bg-primary;
      color: $text-secondary;
      font-size: $font-size-xs;
      font-weight: $font-weight-medium;
      border-radius: $border-radius-small;
      border: 1px solid $border-light;
      transition: all $transition-fast;

      &:hover {
        background-color: $accent-color;
        color: white;
        transform: translateY(-1px);
      }

      .dark & {
        background-color: $dark-bg-tertiary;
        border-color: $dark-border-light;
        color: $dark-text-secondary;

        &:hover {
          background-color: $accent-color;
          color: white;
        }
      }
    }
  }
}

.footer-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: $spacing-lg;
  border-top: 1px solid $border-light;

  .dark & {
    border-top-color: $dark-border-light;
  }

  @media (max-width: $breakpoint-md) {
    flex-direction: column;
    gap: $spacing-md;
    text-align: center;
  }

  .copyright {
    p {
      margin: 0;
      font-size: $font-size-xs;
      color: $text-tertiary;
      line-height: 1.4;

      .dark & {
        color: $dark-text-tertiary;
      }

      &.powered-by {
        margin-top: $spacing-xs;

        a {
          color: $accent-color;
          text-decoration: none;

          &:hover {
            text-decoration: underline;
          }
        }
      }
    }
  }

  .footer-links {
    display: flex;
    gap: $spacing-md;

    @media (max-width: $breakpoint-md) {
      justify-content: center;
    }

    a {
      font-size: $font-size-xs;
      color: $text-tertiary;
      text-decoration: none;
      transition: color $transition-fast;

      &:hover {
        color: $accent-color;
      }

      .dark & {
        color: $dark-text-tertiary;

        &:hover {
          color: $accent-color;
        }
      }
    }
  }
}
</style>
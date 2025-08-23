<template>
  <div class="admin-page">
    <div class="admin-header">
      <h1 class="page-title">管理后台</h1>
      <p class="page-description">欢迎回来，{{ userStore.userInfo?.name }}</p>
    </div>
    
    <div class="admin-stats">
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon size="24"><Document /></el-icon>
          </div>
          <div class="stat-content">
            <h3 class="stat-number">{{ stats.totalArticles }}</h3>
            <p class="stat-label">文章总数</p>
          </div>
          <div class="stat-trend positive">
            <el-icon><TrendCharts /></el-icon>
            <span>+{{ stats.newArticlesThisMonth }}</span>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon size="24"><View /></el-icon>
          </div>
          <div class="stat-content">
            <h3 class="stat-number">{{ formatNumber(stats.totalViews) }}</h3>
            <p class="stat-label">总浏览量</p>
          </div>
          <div class="stat-trend positive">
            <el-icon><TrendCharts /></el-icon>
            <span>+{{ formatNumber(stats.viewsThisMonth) }}</span>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon size="24"><ChatDotRound /></el-icon>
          </div>
          <div class="stat-content">
            <h3 class="stat-number">{{ stats.totalComments }}</h3>
            <p class="stat-label">评论总数</p>
          </div>
          <div class="stat-trend positive">
            <el-icon><TrendCharts /></el-icon>
            <span>+{{ stats.newCommentsThisMonth }}</span>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon size="24"><User /></el-icon>
          </div>
          <div class="stat-content">
            <h3 class="stat-number">{{ stats.totalUsers }}</h3>
            <p class="stat-label">用户总数</p>
          </div>
          <div class="stat-trend positive">
            <el-icon><TrendCharts /></el-icon>
            <span>+{{ stats.newUsersThisMonth }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <div class="admin-content">
      <div class="content-grid">
        <!-- 快速操作 -->
        <div class="quick-actions">
          <h2 class="section-title">快速操作</h2>
          <div class="action-buttons">
            <el-button
              type="primary"
              size="large"
              @click="$router.push('/admin/articles/new')"
            >
              <el-icon><EditPen /></el-icon>
              写新文章
            </el-button>
            
            <el-button
              type="success"
              size="large"
              @click="$router.push('/admin/articles')"
            >
              <el-icon><Document /></el-icon>
              管理文章
            </el-button>
            
            <el-button
              type="info"
              size="large"
              @click="$router.push('/admin/categories')"
            >
              <el-icon><Collection /></el-icon>
              管理分类
            </el-button>
            
            <el-button
              type="warning"
              size="large"
              @click="$router.push('/admin/settings')"
            >
              <el-icon><Setting /></el-icon>
              系统设置
            </el-button>
          </div>
        </div>
        
        <!-- 最近文章 -->
        <div class="recent-articles">
          <div class="section-header">
            <h2 class="section-title">最近文章</h2>
            <el-link type="primary" @click="$router.push('/admin/articles')">
              查看全部
            </el-link>
          </div>
          
          <div class="article-list" v-loading="articlesLoading">
            <div
              v-for="article in recentArticles"
              :key="article.id"
              class="article-item"
              @click="$router.push(`/admin/articles/${article.id}/edit`)"
            >
              <div class="article-info">
                <h3 class="article-title">{{ article.title }}</h3>
                <p class="article-meta">
                  <span class="article-date">{{ formatDate(article.created_at) }}</span>
                  <el-tag :type="getStatusType(article.status)" size="small">
                    {{ getStatusText(article.status) }}
                  </el-tag>
                </p>
              </div>
              
              <div class="article-stats">
                <span class="stat-item">
                  <el-icon><View /></el-icon>
                  {{ article.views }}
                </span>
                <span class="stat-item">
                  <el-icon><ChatDotRound /></el-icon>
                  {{ article.comments_count }}
                </span>
              </div>
            </div>
            
            <div v-if="recentArticles.length === 0" class="empty-state">
              <el-empty description="暂无文章" />
            </div>
          </div>
        </div>
        
        <!-- 系统信息 -->
        <div class="system-info">
          <h2 class="section-title">系统信息</h2>
          <div class="info-list">
            <div class="info-item">
              <span class="info-label">系统版本</span>
              <span class="info-value">v1.0.0</span>
            </div>
            <div class="info-item">
              <span class="info-label">运行时间</span>
              <span class="info-value">{{ systemUptime }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">存储空间</span>
              <span class="info-value">{{ storageUsage }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">最后备份</span>
              <span class="info-value">{{ lastBackup }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import {
  Document,
  View,
  ChatDotRound,
  User,
  TrendCharts,
  EditPen,
  Collection,
  Setting
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import api from '@/utils/api'

const userStore = useUserStore()
const appStore = useAppStore()

// 响应式数据
const articlesLoading = ref(false)

// 统计数据
const stats = reactive({
  totalArticles: 0,
  totalViews: 0,
  totalComments: 0,
  totalUsers: 0,
  newArticlesThisMonth: 0,
  viewsThisMonth: 0,
  newCommentsThisMonth: 0,
  newUsersThisMonth: 0
})

// 最近文章
const recentArticles = ref<any[]>([])

// 系统信息
const systemUptime = ref('7天 12小时')
const storageUsage = ref('2.3GB / 10GB')
const lastBackup = ref('2024-01-15 02:00')

// 计算属性
const formatNumber = (num: number) => {
  if (num >= 1000000) {
    return (num / 1000000).toFixed(1) + 'M'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num.toString()
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    published: 'success',
    draft: 'info',
    archived: 'warning'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    published: '已发布',
    draft: '草稿',
    archived: '已归档'
  }
  return statusMap[status] || '未知'
}

// 方法
const fetchStats = async () => {
  try {
    const response = await api.get('/admin/stats')
    Object.assign(stats, response.data)
  } catch (error) {
    console.error('获取统计数据失败:', error)
    
    // 使用模拟数据
    Object.assign(stats, {
      totalArticles: 42,
      totalViews: 15680,
      totalComments: 234,
      totalUsers: 89,
      newArticlesThisMonth: 5,
      viewsThisMonth: 2340,
      newCommentsThisMonth: 18,
      newUsersThisMonth: 12
    })
  }
}

const fetchRecentArticles = async () => {
  try {
    articlesLoading.value = true
    const response = await api.get('/api/admin/articles/recent')
    recentArticles.value = response.data
  } catch (error) {
    console.error('获取最近文章失败:', error)
    
    // 使用模拟数据
    recentArticles.value = [
      {
        id: 1,
        title: 'Vue 3 Composition API 深度解析',
        status: 'published',
        created_at: '2024-01-15T10:30:00Z',
        views: 1234,
        comments_count: 15
      },
      {
        id: 2,
        title: 'TypeScript 高级类型系统',
        status: 'draft',
        created_at: '2024-01-14T15:20:00Z',
        views: 856,
        comments_count: 8
      },
      {
        id: 3,
        title: 'Vite 构建工具最佳实践',
        status: 'published',
        created_at: '2024-01-13T09:15:00Z',
        views: 2341,
        comments_count: 23
      },
      {
        id: 4,
        title: 'CSS Grid 布局完全指南',
        status: 'published',
        created_at: '2024-01-12T14:45:00Z',
        views: 1876,
        comments_count: 19
      }
    ]
  } finally {
    articlesLoading.value = false
  }
}

// 生命周期
onMounted(async () => {
  appStore.setPageTitle('管理后台')
  appStore.setBreadcrumbs([
    { title: '首页', path: '/' },
    { title: '管理后台', path: '/admin' }
  ])
  
  await Promise.all([
    fetchStats(),
    fetchRecentArticles()
  ])
})
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.admin-page {
  padding: $spacing-lg;
  max-width: 1200px;
  margin: 0 auto;
}

.admin-header {
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

.admin-stats {
  margin-bottom: $spacing-xl;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: $spacing-lg;
}

.stat-card {
  background: $bg-secondary;
  border-radius: $border-radius-large;
  padding: $spacing-lg;
  box-shadow: $shadow-sm;
  border: 1px solid $border-light;
  display: flex;
  align-items: center;
  gap: $spacing-md;
  transition: $transition-fast;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: $shadow-md;
  }
  
  .stat-icon {
    width: 48px;
    height: 48px;
    border-radius: $border-radius-medium;
    background: linear-gradient(135deg, $primary-color, #333333);
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
  }
  
  .stat-content {
    flex: 1;
    
    .stat-number {
      font-size: $font-size-xl;
      font-weight: 700;
      color: $text-primary;
      margin-bottom: $spacing-xs;
    }
    
    .stat-label {
      font-size: $font-size-sm;
      color: $text-secondary;
      margin: 0;
    }
  }
  
  .stat-trend {
    display: flex;
    align-items: center;
    gap: $spacing-xs;
    font-size: $font-size-sm;
    font-weight: 500;
    
    &.positive {
      color: $success-color;
    }
    
    &.negative {
      color: $error-color;
    }
  }
}

.admin-content {
  .content-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: $spacing-xl;
  }
}

.quick-actions,
.recent-articles,
.system-info {
  background: $bg-secondary;
  border-radius: $border-radius-large;
  padding: $spacing-xl;
  box-shadow: $shadow-sm;
  border: 1px solid $border-light;
}

.section-title {
  font-size: $font-size-lg;
  font-weight: 600;
  color: $text-primary;
  margin-bottom: $spacing-lg;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: $spacing-lg;
  
  .section-title {
    margin-bottom: 0;
  }
}

.action-buttons {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: $spacing-md;
  
  .el-button {
    height: 48px;
    display: flex;
    flex-direction: column;
    gap: $spacing-xs;
    
    .el-icon {
      font-size: 18px;
    }
  }
}

.article-list {
  .article-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: $spacing-md;
    border-radius: $border-radius-medium;
    border: 1px solid $border-light;
    margin-bottom: $spacing-md;
    cursor: pointer;
    transition: $transition-fast;
    
    &:hover {
      background: $bg-secondary;
      border-color: $primary-color;
    }
    
    &:last-child {
      margin-bottom: 0;
    }
  }
  
  .article-info {
    flex: 1;
    
    .article-title {
      font-size: $font-size-base;
      font-weight: 500;
      color: $text-primary;
      margin-bottom: $spacing-xs;
      
      &:hover {
        color: $primary-color;
      }
    }
    
    .article-meta {
      display: flex;
      align-items: center;
      gap: $spacing-sm;
      margin: 0;
      
      .article-date {
        font-size: $font-size-sm;
        color: $text-secondary;
      }
    }
  }
  
  .article-stats {
    display: flex;
    gap: $spacing-md;
    
    .stat-item {
      display: flex;
      align-items: center;
      gap: $spacing-xs;
      font-size: $font-size-sm;
      color: $text-secondary;
    }
  }
}

.info-list {
  .info-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: $spacing-md 0;
    border-bottom: 1px solid $border-light;
    
    &:last-child {
      border-bottom: none;
    }
    
    .info-label {
      font-size: $font-size-base;
      color: $text-secondary;
    }
    
    .info-value {
      font-size: $font-size-base;
      font-weight: 500;
      color: $text-primary;
    }
  }
}

.empty-state {
  padding: $spacing-xl 0;
}

// 响应式设计
@media (max-width: 1024px) {
  .admin-content .content-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .admin-page {
    padding: $spacing-md;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .stat-card {
    padding: $spacing-md;
  }
  
  .quick-actions,
  .recent-articles,
  .system-info {
    padding: $spacing-lg;
  }
  
  .action-buttons {
    grid-template-columns: 1fr;
  }
  
  .article-item {
    flex-direction: column;
    align-items: flex-start;
    gap: $spacing-sm;
  }
  
  .article-stats {
    align-self: flex-end;
  }
}

// 暗色模式支持
@media (prefers-color-scheme: dark) {
  .admin-header {
    .page-title {
      color: $dark-text-primary;
    }
    
    .page-description {
      color: $dark-text-secondary;
    }
  }
  
  .stat-card,
  .quick-actions,
  .recent-articles,
  .system-info {
    background: $dark-bg-secondary;
    border-color: $dark-border-light;
  }
  
  .section-title {
    color: $dark-text-primary;
  }
  
  .stat-card {
    .stat-content .stat-number {
      color: $dark-text-primary;
    }
    
    .stat-content .stat-label {
      color: $dark-text-secondary;
    }
  }
  
  .article-list .article-item {
    border-color: $dark-border-light;
    
    &:hover {
      background: $dark-bg-secondary;
    }
    
    .article-info .article-title {
      color: $dark-text-primary;
    }
    
    .article-meta .article-date {
      color: $dark-text-secondary;
    }
    
    .article-stats .stat-item {
      color: $dark-text-secondary;
    }
  }
  
  .info-list .info-item {
    border-color: $dark-border-light;
    
    .info-label {
      color: $dark-text-secondary;
    }
    
    .info-value {
      color: $dark-text-primary;
    }
  }
}
</style>
<template>
  <div class="archives-page">
    <div class="page-header">
      <h1 class="page-title">文章归档</h1>
      <p class="page-description">按时间浏览所有文章</p>
    </div>
    
    <div class="archives-stats" v-if="!loading">
      <div class="stat-item">
        <span class="stat-number">{{ totalArticles }}</span>
        <span class="stat-label">篇文章</span>
      </div>
      <div class="stat-item">
        <span class="stat-number">{{ totalYears }}</span>
        <span class="stat-label">年</span>
      </div>
      <div class="stat-item">
        <span class="stat-number">{{ totalMonths }}</span>
        <span class="stat-label">个月</span>
      </div>
    </div>
    
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="5" animated />
    </div>
    
    <div v-else-if="archives.length === 0" class="empty-state">
      <el-empty description="暂无文章" />
    </div>
    
    <div v-else class="archives-timeline">
      <div
        v-for="yearGroup in archives"
        :key="yearGroup.year"
        class="year-group"
      >
        <div class="year-header">
          <h2 class="year-title">{{ yearGroup.year }}</h2>
          <span class="year-count">{{ yearGroup.totalCount }} 篇</span>
        </div>
        
        <div class="months-container">
          <div
            v-for="monthGroup in yearGroup.months"
            :key="monthGroup.month"
            class="month-group"
          >
            <div class="month-header">
              <h3 class="month-title">{{ monthGroup.monthName }}</h3>
              <span class="month-count">{{ monthGroup.articles.length }} 篇</span>
            </div>
            
            <div class="articles-list">
              <div
                v-for="article in monthGroup.articles"
                :key="article.id"
                class="article-item"
                @click="goToArticle(article.id)"
              >
                <div class="article-date">
                  {{ formatDay(article.created_at) }}
                </div>
                
                <div class="article-content">
                  <h4 class="article-title">{{ article.title }}</h4>
                  <div class="article-meta">
                    <span class="article-category" v-if="article.category">
                      <el-icon><PriceTag /></el-icon>
                      {{ article.category.name }}
                    </span>
                    <span class="article-views">
                      <el-icon><View /></el-icon>
                      {{ article.views || 0 }}
                    </span>
                  </div>
                </div>
                
                <div class="article-arrow">
                  <el-icon><ArrowRight /></el-icon>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  PriceTag,
  View,
  ArrowRight
} from '@element-plus/icons-vue'
import api from '@/utils/api'

interface Article {
  id: number
  title: string
  created_at: string
  views?: number
  category?: {
    id: number
    name: string
  }
}

interface MonthGroup {
  month: number
  monthName: string
  articles: Article[]
}

interface YearGroup {
  year: number
  totalCount: number
  months: MonthGroup[]
}

const router = useRouter()

// 响应式数据
const loading = ref(true)
const archives = ref<YearGroup[]>([])
const allArticles = ref<Article[]>([])

// 计算属性
const totalArticles = computed(() => allArticles.value.length)
const totalYears = computed(() => archives.value.length)
const totalMonths = computed(() => {
  return archives.value.reduce((total, year) => total + year.months.length, 0)
})

// 月份名称映射
const monthNames = [
  '一月', '二月', '三月', '四月', '五月', '六月',
  '七月', '八月', '九月', '十月', '十一月', '十二月'
]

// 方法
const fetchArchives = async () => {
  try {
    loading.value = true
    const response = await api.get('/articles/archives')
    allArticles.value = Array.isArray(response.data) ? response.data : []
    processArchives()
  } catch (error) {
    console.error('获取归档数据失败:', error)
    ElMessage.error('获取归档数据失败')
    
    // 模拟数据
    allArticles.value = [
      {
        id: 1,
        title: 'Vue 3 Composition API 深入理解',
        created_at: '2024-01-15T10:30:00Z',
        views: 256,
        category: { id: 1, name: '技术分享' }
      },
      {
        id: 2,
        title: 'TypeScript 最佳实践',
        created_at: '2024-01-10T14:20:00Z',
        views: 189,
        category: { id: 1, name: '技术分享' }
      },
      {
        id: 3,
        title: '2024年的第一篇随笔',
        created_at: '2024-01-01T09:00:00Z',
        views: 98,
        category: { id: 2, name: '生活随笔' }
      },
      {
        id: 4,
        title: 'Vite 构建优化技巧',
        created_at: '2023-12-28T16:45:00Z',
        views: 312,
        category: { id: 1, name: '技术分享' }
      },
      {
        id: 5,
        title: '年终总结：我的2023',
        created_at: '2023-12-31T20:00:00Z',
        views: 445,
        category: { id: 2, name: '生活随笔' }
      },
      {
        id: 6,
        title: 'React vs Vue：选择指南',
        created_at: '2023-11-15T11:30:00Z',
        views: 567,
        category: { id: 1, name: '技术分享' }
      },
      {
        id: 7,
        title: '秋天的思考',
        created_at: '2023-10-20T15:20:00Z',
        views: 123,
        category: { id: 2, name: '生活随笔' }
      },
      {
        id: 8,
        title: 'Node.js 性能优化实战',
        created_at: '2023-09-05T13:15:00Z',
        views: 389,
        category: { id: 1, name: '技术分享' }
      }
    ]
    processArchives()
  } finally {
    loading.value = false
  }
}

const processArchives = () => {
  const grouped: { [year: number]: { [month: number]: Article[] } } = {}
  
  // 确保allArticles.value是数组
  if (!Array.isArray(allArticles.value)) {
    allArticles.value = []
    return
  }
  
  // 按年月分组
  allArticles.value.forEach(article => {
    const date = new Date(article.created_at)
    const year = date.getFullYear()
    const month = date.getMonth() + 1
    
    if (!grouped[year]) {
      grouped[year] = {}
    }
    
    if (!grouped[year][month]) {
      grouped[year][month] = []
    }
    
    grouped[year][month].push(article)
  })
  
  // 转换为所需格式并排序
  archives.value = Object.keys(grouped)
    .map(year => parseInt(year))
    .sort((a, b) => b - a) // 年份倒序
    .map(year => {
      const months = Object.keys(grouped[year])
        .map(month => parseInt(month))
        .sort((a, b) => b - a) // 月份倒序
        .map(month => ({
          month,
          monthName: monthNames[month - 1],
          articles: grouped[year][month].sort((a, b) => 
            new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
          )
        }))
      
      return {
        year,
        totalCount: months.reduce((total, month) => total + month.articles.length, 0),
        months
      }
    })
}

const goToArticle = (id: number) => {
  router.push(`/article/${id}`)
}

const formatDay = (dateString: string) => {
  return new Date(dateString).getDate().toString().padStart(2, '0')
}

// 生命周期
onMounted(() => {
  fetchArchives()
})
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.archives-page {
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

.archives-stats {
  display: flex;
  justify-content: center;
  gap: $spacing-xl;
  margin-bottom: $spacing-xl;
  padding: $spacing-lg;
  background: $bg-secondary;
  border-radius: $border-radius-large;
  box-shadow: $shadow-sm;
  
  .stat-item {
    text-align: center;
    
    .stat-number {
      display: block;
      font-size: $font-size-xl;
      font-weight: 700;
      color: $primary-color;
      margin-bottom: $spacing-xs;
    }
    
    .stat-label {
      font-size: $font-size-sm;
      color: $text-secondary;
    }
  }
}

.loading-container {
  padding: $spacing-xl 0;
}

.empty-state {
  text-align: center;
  padding: $spacing-2xl 0;
}

.archives-timeline {
  position: relative;
  
  &::before {
    content: '';
    position: absolute;
    left: 30px;
    top: 0;
    bottom: 0;
    width: 2px;
    background: linear-gradient(to bottom, $primary-color, transparent);
  }
}

.year-group {
  margin-bottom: $spacing-xl;
  
  .year-header {
    display: flex;
    align-items: center;
    gap: $spacing-md;
    margin-bottom: $spacing-lg;
    position: relative;
    
    &::before {
      content: '';
      position: absolute;
      left: 24px;
      width: 14px;
      height: 14px;
      background: $primary-color;
      border-radius: 50%;
      border: 3px solid $bg-primary;
      z-index: 1;
    }
    
    .year-title {
      font-size: $font-size-xl;
      font-weight: 700;
      color: $text-primary;
      margin-left: 60px;
    }
    
    .year-count {
      font-size: $font-size-sm;
      color: $text-secondary;
      background: $bg-secondary;
      padding: $spacing-xs $spacing-sm;
      border-radius: $border-radius-small;
    }
  }
}

.months-container {
  margin-left: 60px;
}

.month-group {
  margin-bottom: $spacing-lg;
  
  .month-header {
    display: flex;
    align-items: center;
    gap: $spacing-md;
    margin-bottom: $spacing-md;
    
    .month-title {
      font-size: $font-size-lg;
      font-weight: 600;
      color: $text-primary;
    }
    
    .month-count {
      font-size: $font-size-xs;
      color: $text-secondary;
      background: $bg-primary;
      padding: 2px $spacing-xs;
      border-radius: $border-radius-small;
    }
  }
}

.articles-list {
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;
}

.article-item {
  display: flex;
  align-items: center;
  padding: $spacing-md;
  background: $bg-secondary;
  border-radius: $border-radius-medium;
  border: 1px solid $border-light;
  cursor: pointer;
  transition: all $transition-normal;
  
  &:hover {
    transform: translateX(4px);
    box-shadow: $shadow-md;
    border-color: $primary-color;
    
    .article-arrow {
      transform: translateX(4px);
    }
  }
  
  .article-date {
    flex-shrink: 0;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: $primary-color;
    color: white;
    border-radius: $border-radius-medium;
    font-weight: 600;
    font-size: $font-size-sm;
    margin-right: $spacing-md;
  }
  
  .article-content {
    flex: 1;
    min-width: 0;
    
    .article-title {
      font-size: $font-size-base;
      font-weight: 500;
      color: $text-primary;
      margin-bottom: $spacing-xs;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    
    .article-meta {
      display: flex;
      gap: $spacing-md;
      font-size: $font-size-xs;
      color: $text-secondary;
      
      .article-category,
      .article-views {
        display: flex;
        align-items: center;
        gap: $spacing-xs;
      }
    }
  }
  
  .article-arrow {
    flex-shrink: 0;
    color: $text-secondary;
    transition: transform $transition-normal;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .archives-page {
    padding: $spacing-md;
  }
  
  .archives-stats {
    gap: $spacing-md;
    padding: $spacing-md;
    
    .stat-item .stat-number {
      font-size: $font-size-lg;
    }
  }
  
  .archives-timeline::before {
    left: 20px;
  }
  
  .year-group .year-header {
    &::before {
      left: 14px;
    }
    
    .year-title {
      margin-left: 50px;
      font-size: $font-size-lg;
    }
  }
  
  .months-container {
    margin-left: 50px;
  }
  
  .article-item {
    padding: $spacing-sm;
    
    .article-date {
      width: 32px;
      height: 32px;
      font-size: $font-size-xs;
    }
    
    .article-content .article-meta {
      flex-direction: column;
      gap: $spacing-xs;
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
  
  .archives-stats {
    background: $dark-bg-secondary;
    
    .stat-item .stat-label {
      color: $dark-text-secondary;
    }
  }
  
  .year-group .year-header {
    &::before {
      border-color: $dark-bg-primary;
    }
    
    .year-title {
      color: $dark-text-primary;
    }
    
    .year-count {
      background: $dark-bg-secondary;
      color: $dark-text-secondary;
    }
  }
  
  .month-group .month-header {
    .month-title {
      color: $dark-text-primary;
    }
    
    .month-count {
      background: $dark-bg-primary;
      color: $dark-text-secondary;
    }
  }
  
  .article-item {
    background: $dark-bg-secondary;
    border-color: $dark-border-light;
    
    &:hover {
      border-color: $primary-color;
    }
    
    .article-content {
      .article-title {
        color: $dark-text-primary;
      }
      
      .article-meta {
        color: $dark-text-secondary;
      }
    }
    
    .article-arrow {
      color: $dark-text-secondary;
    }
  }
}
</style>
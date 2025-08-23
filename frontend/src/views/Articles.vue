<template>
  <div class="articles-page">
    <div class="page-header">
      <h1 class="page-title">文章列表</h1>
      <p class="page-description">浏览所有发布的文章</p>
    </div>
    
    <div class="articles-filters">
      <div class="filter-group">
        <el-input
          v-model="searchQuery"
          placeholder="搜索文章..."
          :prefix-icon="Search"
          clearable
          @input="handleSearch"
          class="search-input"
        />
      </div>
      
      <div class="filter-group">
        <el-select
          v-model="selectedCategory"
          placeholder="选择分类"
          clearable
          @change="handleCategoryChange"
          class="category-select"
        >
          <el-option
            v-for="category in categories"
            :key="category.id"
            :label="category.name"
            :value="category.id"
          />
        </el-select>
      </div>
    </div>
    
    <div class="articles-content">
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="3" animated />
        <el-skeleton :rows="3" animated />
        <el-skeleton :rows="3" animated />
      </div>
      
      <div v-else-if="articles.length === 0" class="empty-state">
        <el-empty description="暂无文章" />
      </div>
      
      <div v-else class="articles-grid">
        <article
          v-for="article in articles"
          :key="article.id"
          class="article-card"
          @click="goToArticle(article.id)"
        >
          <div class="article-header">
            <h2 class="article-title">{{ article.title }}</h2>
            <div class="article-meta">
              <span class="article-date">
                <el-icon><Calendar /></el-icon>
                {{ formatDate(article.created_at) }}
              </span>
              <span class="article-category">
                <el-icon><PriceTag /></el-icon>
                {{ article.category?.name || '未分类' }}
              </span>
            </div>
          </div>
          
          <div class="article-summary">
            {{ article.summary || '暂无摘要' }}
          </div>
          
          <div class="article-footer">
            <div class="article-stats">
              <span class="stat-item">
                <el-icon><View /></el-icon>
                {{ article.views || 0 }}
              </span>
              <span class="stat-item">
                <el-icon><ChatDotRound /></el-icon>
                {{ article.comments_count || 0 }}
              </span>
            </div>
            <el-button type="primary" link>
              阅读更多
              <el-icon><ArrowRight /></el-icon>
            </el-button>
          </div>
        </article>
      </div>
      
      <div v-if="articles.length > 0" class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Search,
  Calendar,
  PriceTag,
  View,
  ChatDotRound,
  ArrowRight
} from '@element-plus/icons-vue'
import { useAppStore } from '@/stores/app'
import api from '@/utils/api'

interface Article {
  id: number
  title: string
  summary?: string
  created_at: string
  views?: number
  comments_count?: number
  category?: {
    id: number
    name: string
  }
}

interface Category {
  id: number
  name: string
}

const router = useRouter()
const appStore = useAppStore()

// 响应式数据
const loading = ref(false)
const articles = ref<Article[]>([])
const categories = ref<Category[]>([])
const searchQuery = ref('')
const selectedCategory = ref<number | null>(null)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 计算属性
const filteredArticles = computed(() => {
  let result = articles.value
  
  if (searchQuery.value) {
    result = result.filter(article => 
      article.title.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }
  
  if (selectedCategory.value) {
    result = result.filter(article => 
      article.category?.id === selectedCategory.value
    )
  }
  
  return result
})

// 方法
const fetchArticles = async () => {
  try {
    loading.value = true
    const response = await api.get('/articles', {
      params: {
        page: currentPage.value,
        limit: pageSize.value,
        // category: selectedCategory.value,
        // search: searchQuery.value
      }
    })
    
    articles.value = response.data.articles || []
    total.value = response.data.total || 0
  } catch (error) {
    console.error('获取文章列表失败:', error)
    ElMessage.error('获取文章列表失败')
    // 模拟数据
    articles.value = [
      {
        id: 1,
        title: '欢迎来到我的博客',
        summary: '这是我的第一篇博客文章，欢迎大家阅读和交流。',
        created_at: '2024-01-15T10:30:00Z',
        views: 128,
        comments_count: 5,
        category: { id: 1, name: '随笔' }
      },
      {
        id: 2,
        title: 'Vue 3 开发实践',
        summary: '分享一些 Vue 3 开发中的实用技巧和最佳实践。',
        created_at: '2024-01-10T14:20:00Z',
        views: 256,
        comments_count: 12,
        category: { id: 2, name: '技术' }
      }
    ]
    total.value = 2
  } finally {
    loading.value = false
  }
}

const fetchCategories = async () => {
  try {
    const response = await api.get('/categories')
    categories.value = response.data || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
    // 模拟数据
    categories.value = [
      { id: 1, name: '随笔' },
      { id: 2, name: '技术' },
      { id: 3, name: '生活' }
    ]
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchArticles()
}

const handleCategoryChange = () => {
  currentPage.value = 1
  fetchArticles()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchArticles()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchArticles()
}

const goToArticle = (id: number) => {
  router.push(`/article/${id}`)
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// 生命周期
onMounted(() => {
  fetchCategories()
  fetchArticles()
})
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.articles-page {
  max-width: 1200px;
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

.articles-filters {
  display: flex;
  gap: $spacing-md;
  margin-bottom: $spacing-xl;
  flex-wrap: wrap;
  
  .filter-group {
    flex: 1;
    min-width: 200px;
  }
  
  .search-input {
    max-width: 400px;
  }
  
  .category-select {
    min-width: 150px;
  }
}

.articles-content {
  .loading-container {
    .el-skeleton {
      margin-bottom: $spacing-lg;
    }
  }
  
  .empty-state {
    text-align: center;
    padding: $spacing-2xl 0;
  }
}

.articles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: $spacing-lg;
  margin-bottom: $spacing-xl;
}

.article-card {
  background: $bg-secondary;
  border-radius: $border-radius-large;
  padding: $spacing-lg;
  box-shadow: $shadow-sm;
  transition: all $transition-normal;
  cursor: pointer;
  border: 1px solid $border-light;
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: $shadow-lg;
    border-color: $primary-color;
  }
  
  .article-header {
    margin-bottom: $spacing-md;
    
    .article-title {
      font-size: $font-size-lg;
      font-weight: 600;
      color: $text-primary;
      margin-bottom: $spacing-sm;
      line-height: $line-height-tight;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
    
    .article-meta {
      display: flex;
      gap: $spacing-md;
      font-size: $font-size-sm;
      color: $text-secondary;
      
      .article-date,
      .article-category {
        display: flex;
        align-items: center;
        gap: $spacing-xs;
      }
    }
  }
  
  .article-summary {
    font-size: $font-size-base;
    color: $text-secondary;
    line-height: $line-height-relaxed;
    margin-bottom: $spacing-md;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  
  .article-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    
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
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: $spacing-xl;
}

// 响应式设计
@media (max-width: 768px) {
  .articles-page {
    padding: $spacing-md;
  }
  
  .articles-filters {
    flex-direction: column;
    
    .filter-group {
      min-width: auto;
    }
  }
  
  .articles-grid {
    grid-template-columns: 1fr;
    gap: $spacing-md;
  }
  
  .article-card {
    padding: $spacing-md;
  }
}

// 暗色模式支持
@media (prefers-color-scheme: dark) {
  .article-card {
    background: $dark-bg-secondary;
    border-color: $dark-border-light;
    
    .article-title {
      color: $dark-text-primary;
    }
    
    .article-summary {
      color: $dark-text-secondary;
    }
    
    .article-meta {
      color: $dark-text-secondary;
    }
    
    .article-stats .stat-item {
      color: $dark-text-secondary;
    }
    
    &:hover {
      border-color: $primary-color;
    }
  }
  
  .page-title {
    color: $dark-text-primary;
  }
  
  .page-description {
    color: $dark-text-secondary;
  }
}
</style>
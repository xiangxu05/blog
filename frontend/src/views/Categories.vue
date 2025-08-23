<template>
  <div class="categories-page">
    <div class="page-header">
      <h1 class="page-title">文章分类</h1>
      <p class="page-description">按分类浏览文章</p>
    </div>
    
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="3" animated />
    </div>
    
    <div v-else-if="categories.length === 0" class="empty-state">
      <el-empty description="暂无分类" />
    </div>
    
    <div v-else class="categories-grid">
      <div
        v-for="category in categories"
        :key="category.id"
        class="category-card"
        @click="goToCategory(category.id)"
      >
        <div class="category-icon">
          <el-icon :size="32" :color="getRandomColor()">
            <Folder />
          </el-icon>
        </div>
        
        <div class="category-info">
          <h3 class="category-name">{{ category.name }}</h3>
          <p class="category-description">{{ category.description || '暂无描述' }}</p>
          
          <div class="category-stats">
            <span class="article-count">
              <el-icon><Document /></el-icon>
              {{ category.article_count || 0 }} 篇文章
            </span>
            
            <span v-if="category.last_updated" class="last-updated">
              <el-icon><Clock /></el-icon>
              {{ formatDate(category.last_updated) }}
            </span>
          </div>
        </div>
        
        <div class="category-arrow">
          <el-icon><ArrowRight /></el-icon>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Folder,
  Document,
  Clock,
  ArrowRight
} from '@element-plus/icons-vue'
import api from '@/utils/api'

interface Category {
  id: number
  name: string
  description?: string
  article_count?: number
  last_updated?: string
}

const router = useRouter()

// 响应式数据
const loading = ref(true)
const categories = ref<Category[]>([])

// 颜色数组用于随机分配图标颜色
const colors = [
  '#409EFF', // 主蓝色
  '#67C23A', // 成功绿
  '#E6A23C', // 警告橙
  '#F56C6C', // 危险红
  '#909399', // 信息灰
  '#9C27B0', // 紫色
  '#FF9800', // 橙色
  '#4CAF50', // 绿色
  '#2196F3', // 蓝色
  '#FF5722'  // 深橙色
]

// 方法
const fetchCategories = async () => {
  try {
    loading.value = true
    const response = await api.get('/categories')
    categories.value = response.data || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
    ElMessage.error('获取分类列表失败')
    
    // 模拟数据
    categories.value = [
      {
        id: 1,
        name: '技术分享',
        description: '分享各种技术心得和经验',
        article_count: 15,
        last_updated: '2024-01-15T10:30:00Z'
      },
      {
        id: 2,
        name: '生活随笔',
        description: '记录生活中的点点滴滴',
        article_count: 8,
        last_updated: '2024-01-12T14:20:00Z'
      },
      {
        id: 3,
        name: '学习笔记',
        description: '学习过程中的总结和思考',
        article_count: 12,
        last_updated: '2024-01-10T09:15:00Z'
      },
      {
        id: 4,
        name: '项目实战',
        description: '实际项目开发经验分享',
        article_count: 6,
        last_updated: '2024-01-08T16:45:00Z'
      },
      {
        id: 5,
        name: '工具推荐',
        description: '好用的开发工具和资源推荐',
        article_count: 4,
        last_updated: '2024-01-05T11:30:00Z'
      },
      {
        id: 6,
        name: '读书笔记',
        description: '读书心得和知识总结',
        article_count: 3,
        last_updated: '2024-01-03T13:20:00Z'
      }
    ]
  } finally {
    loading.value = false
  }
}

const goToCategory = (categoryId: number) => {
  router.push({
    name: 'Articles',
    query: { category: categoryId }
  })
}

const getRandomColor = () => {
  return colors[Math.floor(Math.random() * colors.length)]
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

// 生命周期
onMounted(() => {
  fetchCategories()
})
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.categories-page {
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

.loading-container {
  padding: $spacing-xl 0;
}

.empty-state {
  text-align: center;
  padding: $spacing-2xl 0;
}

.categories-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: $spacing-lg;
}

.category-card {
  display: flex;
  align-items: center;
  padding: $spacing-lg;
  background: $bg-secondary;
  border-radius: $border-radius-large;
  box-shadow: $shadow-sm;
  border: 1px solid $border-light;
  cursor: pointer;
  transition: all $transition-normal;
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: $shadow-lg;
    border-color: $primary-color;
    
    .category-arrow {
      transform: translateX(4px);
    }
  }
  
  .category-icon {
    flex-shrink: 0;
    margin-right: $spacing-md;
    padding: $spacing-sm;
    background: rgba(64, 158, 255, 0.1);
    border-radius: $border-radius-medium;
  }
  
  .category-info {
    flex: 1;
    min-width: 0;
    
    .category-name {
      font-size: $font-size-lg;
      font-weight: 600;
      color: $text-primary;
      margin-bottom: $spacing-xs;
    }
    
    .category-description {
      font-size: $font-size-sm;
      color: $text-secondary;
      margin-bottom: $spacing-sm;
      line-height: $line-height-normal;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
    
    .category-stats {
      display: flex;
      gap: $spacing-md;
      font-size: $font-size-xs;
      color: $text-secondary;
      
      .article-count,
      .last-updated {
        display: flex;
        align-items: center;
        gap: $spacing-xs;
      }
    }
  }
  
  .category-arrow {
    flex-shrink: 0;
    color: $text-secondary;
    transition: transform $transition-normal;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .categories-page {
    padding: $spacing-md;
  }
  
  .categories-grid {
    grid-template-columns: 1fr;
    gap: $spacing-md;
  }
  
  .category-card {
    padding: $spacing-md;
    
    .category-info {
      .category-stats {
        flex-direction: column;
        gap: $spacing-xs;
      }
    }
  }
}

@media (max-width: 480px) {
  .category-card {
    flex-direction: column;
    text-align: center;
    
    .category-icon {
      margin-right: 0;
      margin-bottom: $spacing-sm;
    }
    
    .category-arrow {
      margin-top: $spacing-sm;
      transform: rotate(90deg);
    }
    
    &:hover .category-arrow {
      transform: rotate(90deg) translateX(4px);
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
  
  .category-card {
    background: $dark-bg-secondary;
    border-color: $dark-border-light;
    
    &:hover {
      border-color: $primary-color;
    }
    
    .category-info {
      .category-name {
        color: $dark-text-primary;
      }
      
      .category-description {
        color: $dark-text-secondary;
      }
      
      .category-stats {
        color: $dark-text-secondary;
      }
    }
    
    .category-arrow {
      color: $dark-text-secondary;
    }
  }
}
</style>
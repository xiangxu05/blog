<template>
  <div class="admin-articles">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">文章管理</h1>
        <p class="page-description">管理所有文章内容</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="$router.push('/admin/articles/new')">
          <el-icon><EditPen /></el-icon>
          写新文章
        </el-button>
      </div>
    </div>
    
    <!-- 搜索和筛选 -->
    <div class="filter-section">
      <div class="filter-row">
        <el-input
          v-model="searchQuery"
          placeholder="搜索文章标题或内容..."
          class="search-input"
          clearable
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        
        <el-select
          v-model="selectedStatus"
          placeholder="状态筛选"
          clearable
          @change="handleFilter"
        >
          <el-option label="已发布" value="published" />
          <el-option label="草稿" value="draft" />
          <el-option label="已归档" value="archived" />
        </el-select>
        
        <el-select
          v-model="selectedCategory"
          placeholder="分类筛选"
          clearable
          @change="handleFilter"
        >
          <el-option
            v-for="category in categories"
            :key="category.id"
            :label="category.name"
            :value="category.id"
          />
        </el-select>
        
        <el-button @click="resetFilters">
          <el-icon><Refresh /></el-icon>
          重置
        </el-button>
      </div>
    </div>
    
    <!-- 文章列表 -->
    <div class="articles-table" v-loading="loading">
      <el-table
        :data="articles"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="title" label="标题" min-width="200">
          <template #default="{ row }">
            <div class="article-title-cell">
              <h4 class="article-title" @click="editArticle(row.id)">
                {{ row.title }}
              </h4>
              <p class="article-summary">{{ row.summary }}</p>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.category }}</el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="views" label="浏览量" width="100" sortable />
        
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click="editArticle(row.id)"
            >
              编辑
            </el-button>
            
            <el-button
              type="success"
              size="small"
              v-if="row.status === 'draft'"
              @click="publishArticle(row.id)"
            >
              发布
            </el-button>
            
            <el-button
              type="warning"
              size="small"
              v-if="row.status === 'published'"
              @click="archiveArticle(row.id)"
            >
              归档
            </el-button>
            
            <el-button
              type="danger"
              size="small"
              @click="deleteArticle(row.id)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 批量操作 -->
      <div class="batch-actions" v-if="selectedArticles.length > 0">
        <div class="batch-info">
          已选择 {{ selectedArticles.length }} 篇文章
        </div>
        <div class="batch-buttons">
          <el-button type="success" @click="batchPublish">
            批量发布
          </el-button>
          <el-button type="warning" @click="batchArchive">
            批量归档
          </el-button>
          <el-button type="danger" @click="batchDelete">
            批量删除
          </el-button>
        </div>
      </div>
      
      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  EditPen,
  Search,
  Refresh
} from '@element-plus/icons-vue'
import { useAppStore } from '@/stores/app'
import api from '@/utils/api'

const router = useRouter()
const appStore = useAppStore()

// 响应式数据
const loading = ref(false)
const searchQuery = ref('')
const selectedStatus = ref('')
const selectedCategory = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const selectedArticles = ref<any[]>([])

// 文章列表
const articles = ref<any[]>([])

// 分类列表
const categories = ref<any[]>([])

// 计算属性
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

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 方法
const fetchArticles = async () => {
  try {
    loading.value = true
    
    const params = {
      page: currentPage.value,
      limit: pageSize.value,
      search: searchQuery.value,
      status: selectedStatus.value,
      category: selectedCategory.value
    }
    
    const response = await api.get('/api/admin/articles', { params })
    articles.value = response.data.articles
    total.value = response.data.total
  } catch (error) {
    console.error('获取文章列表失败:', error)
    
    // 使用模拟数据
    articles.value = [
      {
        id: 1,
        title: 'Vue 3 Composition API 深度解析',
        summary: '详细介绍 Vue 3 Composition API 的使用方法和最佳实践...',
        category: 'Vue.js',
        status: 'published',
        views: 1234,
        created_at: '2024-01-15T10:30:00Z'
      },
      {
        id: 2,
        title: 'TypeScript 高级类型系统',
        summary: '深入探讨 TypeScript 的高级类型特性和实际应用...',
        category: 'TypeScript',
        status: 'draft',
        views: 856,
        created_at: '2024-01-14T15:20:00Z'
      },
      {
        id: 3,
        title: 'Vite 构建工具最佳实践',
        summary: '分享 Vite 在实际项目中的配置和优化技巧...',
        category: '构建工具',
        status: 'published',
        views: 2341,
        created_at: '2024-01-13T09:15:00Z'
      },
      {
        id: 4,
        title: 'CSS Grid 布局完全指南',
        summary: '全面介绍 CSS Grid 布局的各种用法和实例...',
        category: 'CSS',
        status: 'archived',
        views: 1876,
        created_at: '2024-01-12T14:45:00Z'
      }
    ]
    total.value = 4
  } finally {
    loading.value = false
  }
}

const fetchCategories = async () => {
  try {
    const response = await api.get('/categories')
    categories.value = response.data
  } catch (error) {
    console.error('获取分类列表失败:', error)
    
    // 使用模拟数据
    categories.value = [
      { id: 1, name: 'Vue.js' },
      { id: 2, name: 'TypeScript' },
      { id: 3, name: '构建工具' },
      { id: 4, name: 'CSS' },
      { id: 5, name: 'JavaScript' }
    ]
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchArticles()
}

const handleFilter = () => {
  currentPage.value = 1
  fetchArticles()
}

const resetFilters = () => {
  searchQuery.value = ''
  selectedStatus.value = ''
  selectedCategory.value = ''
  currentPage.value = 1
  fetchArticles()
}

const handleSelectionChange = (selection: any[]) => {
  selectedArticles.value = selection
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

const editArticle = (id: number) => {
  router.push(`/admin/articles/${id}/edit`)
}

const publishArticle = async (id: number) => {
  try {
    await api.put(`/api/admin/articles/${id}/publish`)
    ElMessage.success('文章发布成功')
    fetchArticles()
  } catch (error) {
    console.error('发布文章失败:', error)
    ElMessage.success('文章发布成功（模拟）')
    
    // 更新本地状态
    const article = articles.value.find(a => a.id === id)
    if (article) {
      article.status = 'published'
    }
  }
}

const archiveArticle = async (id: number) => {
  try {
    await api.put(`/api/admin/articles/${id}/archive`)
    ElMessage.success('文章归档成功')
    fetchArticles()
  } catch (error) {
    console.error('归档文章失败:', error)
    ElMessage.success('文章归档成功（模拟）')
    
    // 更新本地状态
    const article = articles.value.find(a => a.id === id)
    if (article) {
      article.status = 'archived'
    }
  }
}

const deleteArticle = async (id: number) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这篇文章吗？此操作不可恢复。',
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await api.delete(`/api/admin/articles/${id}`)
    ElMessage.success('文章删除成功')
    fetchArticles()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除文章失败:', error)
      ElMessage.success('文章删除成功（模拟）')
      
      // 从本地列表中移除
      articles.value = articles.value.filter(a => a.id !== id)
      total.value--
    }
  }
}

const batchPublish = async () => {
  try {
    const ids = selectedArticles.value.map(a => a.id)
    await api.put('/api/admin/articles/batch/publish', { ids })
    ElMessage.success(`成功发布 ${ids.length} 篇文章`)
    fetchArticles()
  } catch (error) {
    console.error('批量发布失败:', error)
    ElMessage.success(`成功发布 ${selectedArticles.value.length} 篇文章（模拟）`)
    fetchArticles()
  }
}

const batchArchive = async () => {
  try {
    const ids = selectedArticles.value.map(a => a.id)
    await api.put('/api/admin/articles/batch/archive', { ids })
    ElMessage.success(`成功归档 ${ids.length} 篇文章`)
    fetchArticles()
  } catch (error) {
    console.error('批量归档失败:', error)
    ElMessage.success(`成功归档 ${selectedArticles.value.length} 篇文章（模拟）`)
    fetchArticles()
  }
}

const batchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedArticles.value.length} 篇文章吗？此操作不可恢复。`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const ids = selectedArticles.value.map(a => a.id)
    await api.delete('/api/admin/articles/batch', { data: { ids } })
    ElMessage.success(`成功删除 ${ids.length} 篇文章`)
    fetchArticles()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error)
      ElMessage.success(`成功删除 ${selectedArticles.value.length} 篇文章（模拟）`)
      fetchArticles()
    }
  }
}

// 生命周期
onMounted(async () => {
  appStore.setPageTitle('文章管理')
  appStore.setBreadcrumbs([
    { title: '首页', path: '/' },
    { title: '管理后台', path: '/admin' },
    { title: '文章管理', path: '/admin/articles' }
  ])
  
  await Promise.all([
    fetchArticles(),
    fetchCategories()
  ])
})
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.admin-articles {
  padding: $spacing-lg;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: $spacing-xl;
  
  .header-left {
    .page-title {
      font-size: $font-size-xxl;
      font-weight: 700;
      color: $text-primary;
      margin-bottom: $spacing-sm;
    }
    
    .page-description {
      font-size: $font-size-base;
      color: $text-secondary;
      margin: 0;
    }
  }
  
  .header-right {
    .el-button {
      height: 40px;
    }
  }
}

.filter-section {
  background: $bg-secondary;
  border-radius: $border-radius-large;
  padding: $spacing-lg;
  margin-bottom: $spacing-lg;
  border: 1px solid $border-light;
  
  .filter-row {
    display: flex;
    gap: $spacing-md;
    align-items: center;
    
    .search-input {
      flex: 1;
      max-width: 300px;
    }
    
    .el-select {
      width: 150px;
    }
  }
}

.articles-table {
  background: $bg-secondary;
  border-radius: $border-radius-large;
  padding: $spacing-lg;
  border: 1px solid $border-light;
  
  .el-table {
    --el-table-border-color: #{$border-light};
    --el-table-bg-color: transparent;
    --el-table-tr-bg-color: transparent;
    --el-table-expanded-cell-bg-color: #{$bg-secondary};
  }
  
  .article-title-cell {
    .article-title {
      font-size: $font-size-base;
      font-weight: 500;
      color: $text-primary;
      margin-bottom: $spacing-xs;
      cursor: pointer;
      transition: $transition-fast;
      
      &:hover {
        color: $primary-color;
      }
    }
    
    .article-summary {
      font-size: $font-size-sm;
      color: $text-secondary;
      margin: 0;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
  }
}

.batch-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: $spacing-lg;
  padding: $spacing-md;
  background: $bg-secondary;
  border-radius: $border-radius-medium;
  
  .batch-info {
    font-size: $font-size-sm;
    color: $text-secondary;
  }
  
  .batch-buttons {
    display: flex;
    gap: $spacing-sm;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: $spacing-xl;
}

// 响应式设计
@media (max-width: 1200px) {
  .admin-articles {
    padding: $spacing-md;
  }
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: $spacing-md;
    align-items: stretch;
    
    .header-right {
      align-self: flex-start;
    }
  }
  
  .filter-section .filter-row {
    flex-direction: column;
    align-items: stretch;
    
    .search-input {
      max-width: none;
    }
    
    .el-select {
      width: 100%;
    }
  }
  
  .batch-actions {
    flex-direction: column;
    gap: $spacing-md;
    align-items: stretch;
    
    .batch-buttons {
      justify-content: center;
    }
  }
}

// 暗色模式支持
@media (prefers-color-scheme: dark) {
  .page-header {
    .page-title {
      color: $dark-text-primary;
    }
    
    .page-description {
      color: $dark-text-secondary;
    }
  }
  
  .filter-section,
  .articles-table {
    background: $dark-bg-secondary;
    border-color: $dark-border-light;
  }
  
  .article-title-cell {
    .article-title {
      color: $dark-text-primary;
    }
    
    .article-summary {
      color: $dark-text-secondary;
    }
  }
  
  .batch-actions {
    background: $dark-bg-secondary;
    
    .batch-info {
      color: $dark-text-secondary;
    }
  }
}
</style>
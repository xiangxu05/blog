<template>
  <div class="article-detail-page">
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="8" animated />
    </div>
    
    <div v-else-if="!article" class="error-container">
      <el-result
        icon="warning"
        title="文章未找到"
        sub-title="抱歉，您访问的文章不存在或已被删除"
      >
        <template #extra>
          <el-button type="primary" @click="goBack">
            返回文章列表
          </el-button>
        </template>
      </el-result>
    </div>
    
    <article v-else class="article-content">
      <header class="article-header">
        <h1 class="article-title">{{ article.title }}</h1>
        
        <div class="article-meta">
          <div class="meta-item">
            <el-icon><User /></el-icon>
            <span>{{ article.author?.name || '匿名' }}</span>
          </div>
          
          <div class="meta-item">
            <el-icon><Calendar /></el-icon>
            <span>{{ formatDate(article.created_at) }}</span>
          </div>
          
          <div class="meta-item">
            <el-icon><PriceTag /></el-icon>
            <span>{{ article.category?.name || '未分类' }}</span>
          </div>
          
          <div class="meta-item">
            <el-icon><View /></el-icon>
            <span>{{ article.views || 0 }} 次阅读</span>
          </div>
        </div>
        
        <div v-if="article.summary" class="article-summary">
          {{ article.summary }}
        </div>
      </header>
      
      <div class="article-body">
        <div v-if="contentLoading" class="content-loading">
          <el-skeleton :rows="10" animated />
        </div>
        
        <div v-else-if="contentError" class="content-error">
          <el-alert
            title="内容加载失败"
            type="error"
            :description="contentError"
            show-icon
            :closable="false"
          />
          <el-button type="primary" @click="loadContent" style="margin-top: 16px;">
            重新加载
          </el-button>
        </div>
        
        <div v-else class="markdown-content" v-html="renderedContent"></div>
      </div>
      
      <footer class="article-footer">
        <div class="article-actions">
          <el-button @click="goBack">
            <el-icon><ArrowLeft /></el-icon>
            返回列表
          </el-button>
          
          <div class="share-buttons">
            <el-button type="primary" @click="shareArticle">
              <el-icon><Share /></el-icon>
              分享
            </el-button>
          </div>
        </div>
        
        <div class="article-tags" v-if="article.tags && article.tags.length > 0">
          <span class="tags-label">标签：</span>
          <el-tag
            v-for="tag in article.tags"
            :key="tag.id"
            type="info"
            size="small"
            style="margin-right: 8px;"
          >
            {{ tag.name }}
          </el-tag>
        </div>
      </footer>
    </article>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  User,
  Calendar,
  PriceTag,
  View,
  ArrowLeft,
  Share
} from '@element-plus/icons-vue'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import api from '@/utils/api'

interface Article {
  id: number
  title: string
  summary?: string
  content?: string
  store_id?: string
  created_at: string
  updated_at?: string
  views?: number
  author?: {
    id: number
    name: string
  }
  category?: {
    id: number
    name: string
  }
  tags?: Array<{
    id: number
    name: string
  }>
}

const route = useRoute()
const router = useRouter()

// 响应式数据
const loading = ref(true)
const contentLoading = ref(false)
const article = ref<Article | null>(null)
const contentError = ref('')

// 计算属性
const renderedContent = computed(() => {
  if (!article.value?.content) return ''
  
  // 配置 marked
  marked.setOptions({
    breaks: true,
    gfm: true
  })
  
  return marked(article.value.content)
})

// 方法
const fetchArticle = async () => {
  try {
    loading.value = true
    const articleId = route.params.id as string
    const version = route.params.version as string || 'latest'
    
    // 使用版本支持的接口
    const response = await api.get(`/articles/${articleId}/${version}`)
    article.value = response.data
    
    // 如果文章有 store_id，需要单独加载内容
    if (article.value?.store_id && !article.value.content) {
      await loadContent()
    }
    
    // 增加阅读量
    incrementViews()
    
  } catch (error) {
    console.error('获取文章失败:', error)
    ElMessage.error('文章加载失败')
    
    // 模拟数据
    const articleId = route.params.id as string
    article.value = {
      id: parseInt(articleId),
      title: '示例文章标题',
      summary: '这是一篇示例文章的摘要内容。',
      content: `# 欢迎阅读

这是一篇示例文章的内容。

## 代码示例

\`\`\`javascript
const message = 'Hello, World!';
console.log(message);
\`\`\`

## 列表示例

- 项目 1
- 项目 2
- 项目 3

感谢您的阅读！`,
      created_at: '2024-01-15T10:30:00Z',
      views: 128,
      author: { id: 1, name: '博主' },
      category: { id: 1, name: '技术' },
      tags: [
        { id: 1, name: 'Vue' },
        { id: 2, name: 'TypeScript' }
      ]
    }
  } finally {
    loading.value = false
  }
}

const loadContent = async () => {
  if (!article.value?.store_id) return
  
  try {
    contentLoading.value = true
    contentError.value = ''
    
    const response = await api.get(`/articles/${article.value.id}/content`)
    if (article.value) {
      article.value.content = response.data.content
    }
  } catch (error) {
    console.error('加载文章内容失败:', error)
    contentError.value = '内容加载失败，请稍后重试'
  } finally {
    contentLoading.value = false
  }
}

const incrementViews = async () => {
  try {
    await api.post(`/articles/${article.value?.id}/view`)
    if (article.value) {
      article.value.views = (article.value.views || 0) + 1
    }
  } catch (error) {
    console.error('更新阅读量失败:', error)
  }
}

const goBack = () => {
  router.push('/articles')
}

const shareArticle = () => {
  if (navigator.share && article.value) {
    navigator.share({
      title: article.value.title,
      text: article.value.summary || '',
      url: window.location.href
    }).catch(err => {
      console.error('分享失败:', err)
      copyToClipboard()
    })
  } else {
    copyToClipboard()
  }
}

const copyToClipboard = () => {
  navigator.clipboard.writeText(window.location.href).then(() => {
    ElMessage.success('链接已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 生命周期
onMounted(() => {
  fetchArticle()
})
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.article-detail-page {
  max-width: 800px;
  margin: 0 auto;
  padding: $spacing-lg;
}

.loading-container,
.error-container {
  padding: $spacing-xl 0;
}

.article-content {
  background: $bg-secondary;
  border-radius: $border-radius-large;
  overflow: hidden;
  box-shadow: $shadow-sm;
  border: 1px solid $border-light;
}

.article-header {
  padding: $spacing-xl;
  border-bottom: 1px solid $border-light;
  
  .article-title {
    font-size: $font-size-xxl;
    font-weight: 700;
    color: $text-primary;
    line-height: $line-height-tight;
    margin-bottom: $spacing-lg;
  }
  
  .article-meta {
    display: flex;
    flex-wrap: wrap;
    gap: $spacing-lg;
    margin-bottom: $spacing-lg;
    
    .meta-item {
      display: flex;
      align-items: center;
      gap: $spacing-xs;
      font-size: $font-size-sm;
      color: $text-secondary;
    }
  }
  
  .article-summary {
    font-size: $font-size-lg;
    color: $text-secondary;
    line-height: $line-height-relaxed;
    padding: $spacing-md;
    background: $bg-primary;
    border-radius: $border-radius-medium;
    border-left: 4px solid $primary-color;
  }
}

.article-body {
  padding: $spacing-xl;
  
  .content-loading,
  .content-error {
    text-align: center;
    padding: $spacing-xl 0;
  }
}

.article-footer {
  padding: $spacing-xl;
  border-top: 1px solid $border-light;
  
  .article-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: $spacing-lg;
  }
  
  .article-tags {
    .tags-label {
      font-size: $font-size-sm;
      color: $text-secondary;
      margin-right: $spacing-sm;
    }
  }
}

// Markdown 内容样式
:deep(.markdown-content) {
  line-height: $line-height-relaxed;
  color: $text-primary;
  
  h1, h2, h3, h4, h5, h6 {
    margin: $spacing-lg 0 $spacing-md 0;
    font-weight: 600;
    line-height: $line-height-tight;
    
    &:first-child {
      margin-top: 0;
    }
  }
  
  h1 { font-size: $font-size-xxl; }
  h2 { font-size: $font-size-xl; }
  h3 { font-size: $font-size-lg; }
  h4 { font-size: $font-size-base; }
  
  p {
    margin: $spacing-md 0;
    
    &:first-child {
      margin-top: 0;
    }
    
    &:last-child {
      margin-bottom: 0;
    }
  }
  
  ul, ol {
    margin: $spacing-md 0;
    padding-left: $spacing-lg;
    
    li {
      margin: $spacing-xs 0;
    }
  }
  
  blockquote {
    margin: $spacing-lg 0;
    padding: $spacing-md $spacing-lg;
    background: $bg-primary;
    border-left: 4px solid $primary-color;
    border-radius: $border-radius-small;
    
    p {
      margin: 0;
      color: $text-secondary;
    }
  }
  
  code {
    background: $bg-primary;
    padding: 2px 6px;
    border-radius: $border-radius-small;
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    font-size: 0.9em;
    color: $primary-color;
  }
  
  pre {
    margin: $spacing-lg 0;
    padding: $spacing-lg;
    background: #f6f8fa;
    border-radius: $border-radius-medium;
    overflow-x: auto;
    
    code {
      background: none;
      padding: 0;
      color: inherit;
    }
  }
  
  table {
    width: 100%;
    margin: $spacing-lg 0;
    border-collapse: collapse;
    
    th, td {
      padding: $spacing-sm $spacing-md;
      border: 1px solid $border-light;
      text-align: left;
    }
    
    th {
      background: $bg-primary;
      font-weight: 600;
    }
  }
  
  img {
    max-width: 100%;
    height: auto;
    border-radius: $border-radius-medium;
    margin: $spacing-md 0;
  }
  
  hr {
    margin: $spacing-xl 0;
    border: none;
    border-top: 1px solid $border-light;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .article-detail-page {
    padding: $spacing-md;
  }
  
  .article-header {
    padding: $spacing-lg;
    
    .article-title {
      font-size: $font-size-xl;
    }
    
    .article-meta {
      flex-direction: column;
      gap: $spacing-sm;
    }
  }
  
  .article-body {
    padding: $spacing-lg;
  }
  
  .article-footer {
    padding: $spacing-lg;
    
    .article-actions {
      flex-direction: column;
      gap: $spacing-md;
      align-items: stretch;
    }
  }
}

// 暗色模式支持
@media (prefers-color-scheme: dark) {
  .article-content {
    background: $dark-bg-secondary;
    border-color: $dark-border-light;
  }
  
  .article-header {
    border-bottom-color: $dark-border-light;
    
    .article-title {
      color: $dark-text-primary;
    }
    
    .article-meta .meta-item {
      color: $dark-text-secondary;
    }
    
    .article-summary {
      background: $dark-bg-primary;
      color: $dark-text-secondary;
    }
  }
  
  .article-footer {
    border-top-color: $dark-border-light;
  }
  
  :deep(.markdown-content) {
    color: $dark-text-primary;
    
    blockquote {
      background: $dark-bg-primary;
      
      p {
        color: $dark-text-secondary;
      }
    }
    
    code {
      background: $dark-bg-primary;
    }
    
    pre {
      background: #1e1e1e;
    }
    
    table {
      th, td {
        border-color: $dark-border-light;
      }
      
      th {
        background: $dark-bg-primary;
      }
    }
    
    hr {
      border-top-color: $dark-border-light;
    }
  }
}
</style>
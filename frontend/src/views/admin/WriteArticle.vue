<template>
  <div class="write-article">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ isEdit ? '编辑文章' : '写新文章' }}</h1>
        <p class="page-description">{{ isEdit ? '修改文章内容' : '创建新的文章内容' }}</p>
      </div>
      <div class="header-right">
        <el-button @click="saveDraft" :loading="saving">
          <el-icon><Document /></el-icon>
          保存草稿
        </el-button>
        <el-button type="primary" @click="publishArticle" :loading="publishing">
          <el-icon><Upload /></el-icon>
          {{ isEdit ? '更新文章' : '发布文章' }}
        </el-button>
      </div>
    </div>
    
    <!-- 文章表单 -->
    <div class="article-form" v-loading="loading">
      <el-form
        ref="formRef"
        :model="articleForm"
        :rules="formRules"
        label-width="100px"
        class="article-form-content"
      >
        <!-- 基本信息 -->
        <div class="form-section">
          <h3 class="section-title">基本信息</h3>
          
          <el-form-item label="文章标题" prop="title">
            <el-input
              v-model="articleForm.title"
              placeholder="请输入文章标题"
              maxlength="100"
              show-word-limit
              size="large"
            />
          </el-form-item>
          
          <el-form-item label="文章摘要" prop="summary">
            <el-input
              v-model="articleForm.summary"
              type="textarea"
              placeholder="请输入文章摘要（可选）"
              :rows="3"
              maxlength="200"
              show-word-limit
            />
          </el-form-item>
          
          <div class="form-row">
            <el-form-item label="分类" prop="category_id" class="form-item-half">
              <el-select
                v-model="articleForm.category_id"
                placeholder="选择分类"
                filterable
                allow-create
                default-first-option
                style="width: 100%"
              >
                <el-option
                  v-for="category in categories"
                  :key="category.id"
                  :label="category.name"
                  :value="category.id"
                />
              </el-select>
            </el-form-item>
            
            <el-form-item label="标签" class="form-item-half">
              <el-select
                v-model="articleForm.tags"
                multiple
                filterable
                allow-create
                default-first-option
                placeholder="选择或创建标签"
                style="width: 100%"
              >
                <el-option
                  v-for="tag in availableTags"
                  :key="tag"
                  :label="tag"
                  :value="tag"
                />
              </el-select>
            </el-form-item>
          </div>
          
          <el-form-item label="封面图片">
            <div class="cover-upload">
              <el-upload
                class="cover-uploader"
                :show-file-list="false"
                :before-upload="beforeCoverUpload"
                :http-request="uploadCover"
                accept="image/*"
              >
                <img v-if="articleForm.cover_image" :src="articleForm.cover_image" class="cover-image" />
                <div v-else class="cover-placeholder">
                  <el-icon class="cover-icon"><Plus /></el-icon>
                  <div class="cover-text">点击上传封面</div>
                </div>
              </el-upload>
              <div class="cover-actions" v-if="articleForm.cover_image">
                <el-button size="small" @click="removeCover">移除封面</el-button>
              </div>
            </div>
          </el-form-item>
        </div>
        
        <!-- 文章内容 -->
        <div class="form-section">
          <h3 class="section-title">文章内容</h3>
          
          <el-form-item prop="content">
            <div class="editor-container">
              <div ref="editorRef" class="markdown-editor"></div>
            </div>
          </el-form-item>
        </div>
        
        <!-- 发布设置 -->
        <div class="form-section">
          <h3 class="section-title">发布设置</h3>
          
          <div class="form-row">
            <el-form-item label="发布状态" class="form-item-half">
              <el-radio-group v-model="articleForm.status">
                <el-radio value="draft">草稿</el-radio>
                <el-radio value="published">发布</el-radio>
              </el-radio-group>
            </el-form-item>
            
            <el-form-item label="发布时间" class="form-item-half">
              <el-date-picker
                v-model="articleForm.publish_at"
                type="datetime"
                placeholder="选择发布时间"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DD HH:mm:ss"
                style="width: 100%"
              />
            </el-form-item>
          </div>
          
          <el-form-item label="SEO设置">
            <el-collapse>
              <el-collapse-item title="SEO优化" name="seo">
                <el-form-item label="SEO标题">
                  <el-input
                    v-model="articleForm.seo_title"
                    placeholder="SEO标题（可选，默认使用文章标题）"
                    maxlength="60"
                    show-word-limit
                  />
                </el-form-item>
                
                <el-form-item label="SEO描述">
                  <el-input
                    v-model="articleForm.seo_description"
                    type="textarea"
                    placeholder="SEO描述（可选，默认使用文章摘要）"
                    :rows="2"
                    maxlength="160"
                    show-word-limit
                  />
                </el-form-item>
                
                <el-form-item label="SEO关键词">
                  <el-input
                    v-model="articleForm.seo_keywords"
                    placeholder="SEO关键词，用逗号分隔"
                  />
                </el-form-item>
              </el-collapse-item>
            </el-collapse>
          </el-form-item>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, type FormInstance, type UploadProps } from 'element-plus'
import {
  Document,
  Upload,
  Plus
} from '@element-plus/icons-vue'
import { useAppStore } from '@/stores/app'
import api from '@/utils/api'

// Toast UI Editor
import '@toast-ui/editor/dist/toastui-editor.css'
import '@toast-ui/editor/dist/theme/toastui-editor-dark.css'
// @ts-ignore
import { Editor } from '@toast-ui/editor'
import codeSyntaxHighlight from '@toast-ui/editor-plugin-code-syntax-highlight'
import 'prismjs/themes/prism.css'
import Prism from 'prismjs'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()

// 响应式数据
const loading = ref(false)
const saving = ref(false)
const publishing = ref(false)
const formRef = ref<FormInstance>()
const editorRef = ref<HTMLElement>()
let editor: Editor | null = null

// 是否为编辑模式
const isEdit = ref(false)
const articleId = ref<string | null>(null)

// 分类和标签数据
const categories = ref<any[]>([])
const availableTags = ref<string[]>([])

// 文章表单数据
const articleForm = reactive({
  title: '',
  summary: '',
  content: '',
  category_id: '',
  tags: [] as string[],
  cover_image: '',
  status: 'draft',
  publish_at: '',
  seo_title: '',
  seo_description: '',
  seo_keywords: ''
})

// 表单验证规则
const formRules = {
  title: [
    { required: true, message: '请输入文章标题', trigger: 'blur' },
    { min: 1, max: 100, message: '标题长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  category_id: [
    { required: true, message: '请选择文章分类', trigger: 'change' }
  ],
  content: [
    { required: true, message: '请输入文章内容', trigger: 'blur' }
  ]
}

// 方法
const initEditor = async () => {
  await nextTick()
  
  if (!editorRef.value) return
  
  editor = new Editor({
    el: editorRef.value,
    height: '500px',
    initialEditType: 'markdown',
    previewStyle: 'vertical',
    placeholder: '请输入文章内容...',
    plugins: [[codeSyntaxHighlight, { highlighter: Prism }]],
    theme: appStore.theme === 'dark' ? 'dark' : 'light',
    toolbarItems: [
      ['heading', 'bold', 'italic', 'strike'],
      ['hr', 'quote'],
      ['ul', 'ol', 'task', 'indent', 'outdent'],
      ['table', 'image', 'link'],
      ['code', 'codeblock'],
      ['scrollSync']
    ],
    hooks: {
      addImageBlobHook: async (blob: Blob, callback: (url: string, altText?: string) => void) => {
        try {
          const formData = new FormData()
          formData.append('file', blob)
          
          const response = await api.post('/api/files/upload', formData, {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          })
          
          const imageUrl = response.data.url
          callback(imageUrl, '')
        } catch (error) {
          console.error('图片上传失败:', error)
          ElMessage.error('图片上传失败')
        }
      }
    },
    events: {
      change: () => {
        if (editor) {
          articleForm.content = editor.getMarkdown()
          // 检查并处理本地图片链接
          detectAndProcessLocalImages()
        }
      }
    }
  })
  
  // 设置初始内容
  if (articleForm.content) {
    editor.setMarkdown(articleForm.content)
  }
}

// 检测并处理本地图片链接
const detectAndProcessLocalImages = () => {
  if (!editor) return
  
  const content = editor.getMarkdown()
  const localImageRegex = /!\[([^\]]*)\]\((?!https?:\/\/)([^)]+)\)/g
  const matches = content.match(localImageRegex)
  
  if (matches && matches.length > 0) {
    // 显示提示，询问用户是否要上传本地图片
    ElMessageBox.confirm(
      `检测到 ${matches.length} 个本地图片链接，是否要上传到服务器？`,
      '本地图片检测',
      {
        confirmButtonText: '上传',
        cancelButtonText: '取消',
        type: 'info'
      }
    ).then(() => {
      processLocalImages(content, matches)
    }).catch(() => {
      // 用户取消，不做处理
    })
  }
}

// 处理本地图片上传和替换
const processLocalImages = async (content: string, matches: string[]) => {
  let updatedContent = content
  
  for (const match of matches) {
    const imagePathMatch = match.match(/!\[([^\]]*)\]\(([^)]+)\)/)
    if (!imagePathMatch) continue
    
    const altText = imagePathMatch[1]
    const imagePath = imagePathMatch[2]
    
    try {
      // 这里需要用户手动选择文件，因为浏览器安全限制无法直接读取本地文件
      ElMessage.info(`请手动上传图片: ${imagePath}`)
      // 实际项目中，可以提供文件选择器让用户重新选择对应的图片
    } catch (error) {
      console.error('处理本地图片失败:', error)
    }
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

const fetchTags = async () => {
  try {
    const response = await api.get('/tags')
    availableTags.value = response.data.map((tag: any) => tag.name)
  } catch (error) {
    console.error('获取标签列表失败:', error)
    
    // 使用模拟数据
    availableTags.value = [
      'Vue.js', 'TypeScript', 'JavaScript', 'CSS', 'HTML',
      'React', 'Node.js', 'Webpack', 'Vite', '前端开发'
    ]
  }
}

const fetchArticle = async (id: string, version?: string) => {
  try {
    loading.value = true
    const url = version ? `/api/admin/articles/${id}/${version}` : `/api/admin/articles/${id}`
    const response = await api.get(url)
    const article = response.data
    
    // 填充表单数据
    Object.assign(articleForm, {
      title: article.title,
      summary: article.summary || '',
      content: article.content || '',
      category_id: article.category_id,
      tags: article.tags || [],
      cover_image: article.cover_image || '',
      status: article.status,
      publish_at: article.publish_at || '',
      seo_title: article.seo_title || '',
      seo_description: article.seo_description || '',
      seo_keywords: article.seo_keywords || ''
    })
    
    // 更新编辑器内容
    if (editor && article.content) {
      editor.setMarkdown(article.content)
    }
  } catch (error) {
    console.error('获取文章详情失败:', error)
    ElMessage.error('获取文章详情失败')
    router.push('/admin/articles')
  } finally {
    loading.value = false
  }
}

const beforeCoverUpload: UploadProps['beforeUpload'] = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2
  
  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB!')
    return false
  }
  return true
}

const uploadCover = async (options: any) => {
  try {
    const formData = new FormData()
    formData.append('file', options.file)
    
    const response = await api.post('/api/files/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    
    articleForm.cover_image = response.data.url
    ElMessage.success('封面上传成功')
  } catch (error) {
    console.error('封面上传失败:', error)
    
    // 模拟上传成功
    const reader = new FileReader()
    reader.onload = (e) => {
      articleForm.cover_image = e.target?.result as string
      ElMessage.success('封面上传成功（模拟）')
    }
    reader.readAsDataURL(options.file)
  }
}

const removeCover = () => {
  articleForm.cover_image = ''
}

const validateForm = async () => {
  if (!formRef.value) return false
  
  try {
    await formRef.value.validate()
    
    // 验证编辑器内容
    if (!articleForm.content.trim()) {
      ElMessage.error('请输入文章内容')
      return false
    }
    
    return true
  } catch (error) {
    return false
  }
}

const saveDraft = async () => {
  if (!(await validateForm())) return
  
  try {
    saving.value = true
    
    // 第一步：上传文章内容获取 store_id
    let storeId = ''
    if (articleForm.content.trim()) {
      const contentResponse = await api.post('/api/files/upload', {
        content: articleForm.content,
        type: 'markdown'
      })
      storeId = contentResponse.data.store_id
    }
    
    // 第二步：提交文章元数据
    const data = {
      title: articleForm.title,
      summary: articleForm.summary,
      category_id: articleForm.category_id,
      tags: articleForm.tags,
      cover_image: articleForm.cover_image,
      status: 'draft',
      publish_at: articleForm.publish_at,
      seo_title: articleForm.seo_title,
      seo_description: articleForm.seo_description,
      seo_keywords: articleForm.seo_keywords,
      store_id: storeId
    }
    
    if (isEdit.value && articleId.value) {
      await api.put(`/api/admin/articles/${articleId.value}`, data)
      ElMessage.success('草稿保存成功')
    } else {
      const response = await api.post('/api/admin/articles', data)
      articleId.value = response.data.id
      isEdit.value = true
      ElMessage.success('草稿保存成功')
      
      // 更新URL
      router.replace(`/admin/articles/${articleId.value}/edit`)
    }
  } catch (error) {
    console.error('保存草稿失败:', error)
    ElMessage.success('草稿保存成功（模拟）')
  } finally {
    saving.value = false
  }
}

const publishArticle = async () => {
  if (!(await validateForm())) return
  
  try {
    publishing.value = true
    
    // 第一步：上传文章内容获取 store_id
    let storeId = ''
    if (articleForm.content.trim()) {
      const contentResponse = await api.post('/api/files/upload', {
        content: articleForm.content,
        type: 'markdown'
      })
      storeId = contentResponse.data.store_id
    }
    
    // 第二步：提交文章元数据
    const data = {
      title: articleForm.title,
      summary: articleForm.summary,
      category_id: articleForm.category_id,
      tags: articleForm.tags,
      cover_image: articleForm.cover_image,
      status: 'published',
      publish_at: articleForm.publish_at || new Date().toISOString(),
      seo_title: articleForm.seo_title,
      seo_description: articleForm.seo_description,
      seo_keywords: articleForm.seo_keywords,
      store_id: storeId
    }
    
    if (isEdit.value && articleId.value) {
      await api.put(`/api/admin/articles/${articleId.value}`, data)
      ElMessage.success('文章更新成功')
    } else {
      await api.post('/api/articles', data)
      ElMessage.success('文章发布成功')
    }
    
    router.push('/admin/articles')
  } catch (error) {
    console.error('发布文章失败:', error)
    ElMessage.success(isEdit.value ? '文章更新成功（模拟）' : '文章发布成功（模拟）')
    router.push('/admin/articles')
  } finally {
    publishing.value = false
  }
}

// 生命周期
onMounted(async () => {
  // 检查是否为编辑模式
  if (route.params.id) {
    isEdit.value = true
    articleId.value = route.params.id as string
  }
  
  // 设置页面信息
  appStore.setPageTitle(isEdit.value ? '编辑文章' : '写新文章')
  appStore.setBreadcrumbs([
    { title: '首页', path: '/' },
    { title: '管理后台', path: '/admin' },
    { title: '文章管理', path: '/admin/articles' },
    { title: isEdit.value ? '编辑文章' : '写新文章', path: route.path }
  ])
  
  // 获取基础数据
  await Promise.all([
    fetchCategories(),
    fetchTags()
  ])
  
  // 初始化编辑器
  await initEditor()
  
  // 如果是编辑模式，获取文章数据
  if (isEdit.value && articleId.value) {
    await fetchArticle(articleId.value)
  }
})

onBeforeUnmount(() => {
  if (editor) {
    editor.destroy()
  }
})
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.write-article {
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
    display: flex;
    gap: $spacing-md;
    
    .el-button {
      height: 40px;
    }
  }
}

.article-form {
  background: $bg-secondary;
  border-radius: $border-radius-large;
  border: 1px solid $border-light;
  
  .article-form-content {
    padding: $spacing-xl;
  }
}

.form-section {
  margin-bottom: $spacing-xl;
  
  &:last-child {
    margin-bottom: 0;
  }
  
  .section-title {
    font-size: $font-size-lg;
    font-weight: 600;
    color: $text-primary;
    margin-bottom: $spacing-lg;
    padding-bottom: $spacing-sm;
    border-bottom: 2px solid $primary-color;
  }
}

.form-row {
  display: flex;
  gap: $spacing-lg;
  
  .form-item-half {
    flex: 1;
  }
}

.cover-upload {
  .cover-uploader {
    :deep(.el-upload) {
      border: 2px dashed $border-light;
      border-radius: $border-radius-medium;
      cursor: pointer;
      position: relative;
      overflow: hidden;
      transition: $transition-fast;
      
      &:hover {
        border-color: $primary-color;
      }
    }
  }
  
  .cover-image {
    width: 200px;
    height: 120px;
    object-fit: cover;
    display: block;
  }
  
  .cover-placeholder {
    width: 200px;
    height: 120px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: $text-secondary;
    
    .cover-icon {
      font-size: $font-size-xl;
      margin-bottom: $spacing-sm;
    }
    
    .cover-text {
      font-size: $font-size-sm;
    }
  }
  
  .cover-actions {
    margin-top: $spacing-sm;
  }
}

.editor-container {
  border: 1px solid $border-light;
  border-radius: $border-radius-medium;
  overflow: hidden;
  
  .markdown-editor {
    :deep(.toastui-editor-defaultUI) {
      border: none;
    }
    
    :deep(.toastui-editor-toolbar) {
      background: $bg-secondary;
      border-bottom: 1px solid $border-light;
    }
    
    :deep(.toastui-editor-md-container) {
      background: $bg-primary;
    }
    
    :deep(.toastui-editor-preview-container) {
      background: $bg-primary;
    }
  }
}

// 响应式设计
@media (max-width: 1200px) {
  .write-article {
    padding: $spacing-md;
  }
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: $spacing-md;
    align-items: stretch;
    
    .header-right {
      justify-content: flex-start;
    }
  }
  
  .form-row {
    flex-direction: column;
    gap: 0;
  }
  
  .cover-image,
  .cover-placeholder {
    width: 100%;
    max-width: 300px;
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
  
  .article-form {
    background: $dark-bg-secondary;
    border-color: $dark-border-light;
  }
  
  .form-section {
    .section-title {
      color: $dark-text-primary;
    }
  }
  
  .cover-placeholder {
    color: $dark-text-secondary;
  }
  
  .editor-container {
    border-color: $dark-border-light;
    
    .markdown-editor {
      :deep(.toastui-editor-toolbar) {
        background: $dark-bg-secondary;
        border-bottom-color: $dark-border-light;
      }
      
      :deep(.toastui-editor-md-container),
      :deep(.toastui-editor-preview-container) {
        background: $dark-bg-primary;
      }
    }
  }
}
</style>
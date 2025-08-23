<template>
  <div class="home">
    <div class="hero-section">
      <h1>欢迎来到我的博客</h1>
      <p>分享技术、记录生活</p>
    </div>
    
    <div class="content-section">
      <h2>最新文章</h2>
      <div class="article-list">
        <div class="article-card" v-for="i in 3" :key="i">
          <h3>示例文章标题 {{ i }}</h3>
          <p>这是一篇示例文章的摘要内容...</p>
          <div class="article-meta">
            <span>2024-01-{{ 10 + i }}</span>
            <span>技术分享</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

onMounted(() => {
  appStore.setPageTitle('首页')
  appStore.setBreadcrumbs([{ title: '首页', path: '/' }])
})
</script>

<style lang="scss" scoped>
.home {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.hero-section {
  text-align: center;
  padding: 4rem 0;
  
  h1 {
    font-size: 3rem;
    font-weight: 700;
    margin-bottom: 1rem;
    color: #1D1D1F;
  }
  
  p {
    font-size: 1.2rem;
    color: #8E8E93;
  }
}

.content-section {
  margin-top: 4rem;
  
  h2 {
    font-size: 2rem;
    margin-bottom: 2rem;
    color: #1D1D1F;
  }
}

.article-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
}

.article-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.07);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 15px rgba(0, 0, 0, 0.1);
  }
  
  h3 {
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
    color: #1D1D1F;
  }
  
  p {
    color: #8E8E93;
    line-height: 1.6;
    margin-bottom: 1rem;
  }
}

.article-meta {
  display: flex;
  gap: 1rem;
  font-size: 0.875rem;
  color: #C7C7CC;
}

.dark {
  .hero-section h1,
  .content-section h2,
  .article-card h3 {
    color: #FFFFFF;
  }
  
  .article-card {
    background: #2C2C2E;
  }
}
</style>
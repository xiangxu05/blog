import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import { ElMessage } from 'element-plus'

// 路由元信息类型定义
declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    requiresAuth?: boolean
    requiresAdmin?: boolean
    breadcrumbs?: Array<{ title: string; path?: string }>
  }
}

// 路由配置
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: {
      title: '首页',
      breadcrumbs: [{ title: '首页', path: '/' }]
    }
  },
  {
    path: '/articles',
    name: 'Articles',
    component: () => import('@/views/Articles.vue'),
    meta: {
      title: '文章列表',
      breadcrumbs: [
        { title: '首页', path: '/' },
        { title: '文章列表' }
      ]
    }
  },
  {
    path: '/article/:id/:version?',
    name: 'ArticleDetail',
    component: () => import('@/views/Detail.vue'),
    meta: {
      title: '文章详情',
      breadcrumbs: [
        { title: '首页', path: '/' },
        { title: '文章列表', path: '/articles' },
        { title: '文章详情' }
      ]
    }
  },
  {
    path: '/categories',
    name: 'Categories',
    component: () => import('@/views/Categories.vue'),
    meta: {
      title: '分类',
      breadcrumbs: [
        { title: '首页', path: '/' },
        { title: '分类' }
      ]
    }
  },
  {
    path: '/archives',
    name: 'Archives',
    component: () => import('@/views/Archives.vue'),
    meta: {
      title: '归档',
      breadcrumbs: [
        { title: '首页', path: '/' },
        { title: '归档' }
      ]
    }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: {
      title: '登录',
      requiresGuest: true
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: {
      title: '注册',
      requiresGuest: true
    }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: {
      title: '个人资料',
      requiresAuth: true,
      breadcrumbs: [
        { title: '首页', path: '/' },
        { title: '个人资料', path: '/profile' }
      ]
    }
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('@/views/admin/Admin.vue'),
    meta: {
      title: '管理后台',
      requiresAuth: true,
      requiresAdmin: true,
      breadcrumbs: [
        { title: '首页', path: '/' },
        { title: '管理后台' }
      ]
    }
  },
  {
    path: '/admin/articles',
    name: 'AdminArticles',
    component: () => import('@/views/admin/Articles.vue'),
    meta: {
      title: '文章管理',
      requiresAuth: true,
      requiresAdmin: true,
      breadcrumbs: [
        { title: '首页', path: '/' },
        { title: '管理后台', path: '/admin' },
        { title: '文章管理' }
      ]
    }
  },
  {
    path: '/admin/write',
    name: 'AdminWrite',
    component: () => import('@/views/admin/WriteArticle.vue'),
    meta: {
      title: '写文章',
      requiresAuth: true,
      requiresAdmin: true,
      breadcrumbs: [
        { title: '首页', path: '/' },
        { title: '管理后台', path: '/admin' },
        { title: '写文章' }
      ]
    }
  },
  {
    path: '/403',
    name: 'Forbidden',
    component: () => import('@/views/error/403.vue'),
    meta: {
      title: '访问被拒绝',
      breadcrumbs: [
        { title: '首页', path: '/' },
        { title: '访问被拒绝' }
      ]
    }
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: {
      title: '页面未找到',
      breadcrumbs: [
        { title: '首页', path: '/' },
        { title: '页面未找到' }
      ]
    }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404'
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  const appStore = useAppStore()

  // 设置页面标题
  if (to.meta.title) {
    appStore.setPageTitle(to.meta.title)
    document.title = `${to.meta.title} - 个人博客`
  }

  // 设置面包屑
  if (to.meta.breadcrumbs) {
    appStore.setBreadcrumbs(to.meta.breadcrumbs)
  }

  // 检查是否需要认证
  if (to.meta.requiresAuth) {
    // 如果用户未登录，重定向到登录页
    if (!userStore.isLoggedIn) {
      ElMessage.warning('请先登录')
      next({
        name: 'Login',
        query: { redirect: to.fullPath }
      })
      return
    }

    // 检查是否需要管理员权限
    if (to.meta.requiresAdmin) {
      if (!userStore.isAdmin) {
        ElMessage.error('您没有访问权限')
        next({ name: 'Forbidden' })
        return
      }
    }
  }

  // 如果已登录用户访问登录/注册页，重定向到首页
  if ((to.name === 'Login' || to.name === 'Register') && userStore.isLoggedIn) {
    next({ name: 'Home' })
    return
  }

  next()
})

// 路由后置守卫
router.afterEach(() => {
  // 页面切换后滚动到顶部
  window.scrollTo(0, 0)
})

export default router

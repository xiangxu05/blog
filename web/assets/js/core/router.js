// SPA 路由系统 - router.js

class Router {
  constructor() {
    this.routes = new Map();
    this.currentRoute = null;
    this.beforeHooks = [];
    this.afterHooks = [];
    this.initialized = false;
  }

  // 初始化路由系统
  init() {
    if (this.initialized) return;

    // 注册路由配置
    this.registerRoutes();

    // 设置全局导航守卫
    page('*', (ctx, next) => {
      this.globalGuard(ctx, next);
    });

    // 启动路由
    page.start();
    this.initialized = true;

    console.log('SPA 路由系统已初始化');
  }

  // 注册所有路由
  registerRoutes() {
    // 定义路由配置
    const routeConfigs = {
      // 公开页面
      '/': {
        template: 'home.html',
        title: '首页',
        requireAuth: false
      },
      '/login': {
        template: 'login.html',
        title: '用户登录',
        requireAuth: false,
        redirectIfAuth: '/' // 已登录用户访问登录页时重定向
      },
      '/register': {
        template: 'register.html',
        title: '用户注册',
        requireAuth: false,
        redirectIfAuth: '/'
      },
      '/articles': {
        template: 'articles.html',
        title: '文章列表',
        requireAuth: false
      },
      '/articles/:id/:version': {
        template: 'article-detail.html',
        title: '文章详情',
        requireAuth: false
      },

      // 需要登录的页面
      '/profile': {
        template: 'profile.html',
        title: '个人资料',
        requireAuth: true
      },

      // 管理员页面
      '/admin': {
        template: 'admin/dashboard.html',
        title: '管理后台',
        requireAuth: true,
        requireRole: 'admin'
      },
      '/admin/articles': {
        template: 'admin/article-list.html',
        title: '文章管理',
        requireAuth: true,
        requireRole: 'admin'
      },
      '/admin/articles/new': {
        template: 'admin/article-editor.html',
        title: '新建文章',
        requireAuth: true,
        requireRole: 'admin'
      },
      '/admin/articles/:id/edit': {
        template: 'admin/article-editor.html',
        title: '编辑文章',
        requireAuth: true,
        requireRole: 'admin'
      },
      '/admin/files': {
        template: 'admin/file-manager.html',
        title: '文件管理',
        requireAuth: true,
        requireRole: 'admin'
      },

      // 错误页面
      '/404': {
        template: '404.html',
        title: '页面不存在',
        requireAuth: false
      }
    };

    // 注册路由
    Object.entries(routeConfigs).forEach(([path, config]) => {
      this.routes.set(path, config);
      page(path, (ctx) => this.handleRoute(ctx, config));
    });

    // 404处理
    page('*', (ctx) => {
      this.handleRoute(ctx, {
        template: '404.html',
        title: '页面不存在',
        requireAuth: false
      });
    });
  }

  // 全局路由守卫
  async globalGuard(ctx, next) {
    const route = this.routes.get(ctx.pathname) || this.routes.get(ctx.routePath);

    // 执行前置钩子
    for (const hook of this.beforeHooks) {
      const result = await hook(ctx, route);
      if (result === false) return;
    }

    // 检查用户认证状态
    const isAuthenticated = window.appState.isAuthenticated();
    const userRole = window.appState.getState('user.role');

    // 如果需要登录但用户未登录
    if (route && route.requireAuth && !isAuthenticated) {
      this.showNotification('请先登录', 'warning');
      page.redirect('/login?redirect=' + encodeURIComponent(ctx.pathname));
      return;
    }

    // 如果需要管理员权限但用户不是管理员
    if (route && route.requireRole === 'admin' && userRole !== 'admin') {
      this.showNotification('权限不足', 'error');
      page.redirect('/404');
      return;
    }

    // 如果已登录用户访问登录/注册页面
    if (route && route.redirectIfAuth && isAuthenticated) {
      page.redirect(route.redirectIfAuth);
      return;
    }

    next();
  }

  // 处理路由
  async handleRoute(ctx, config) {
    try {
      // 显示加载状态
      window.appState.setLoading(true);

      // 更新当前路由状态
      this.currentRoute = {
        path: ctx.pathname,
        params: ctx ? ctx.params : {},
        query: this.parseQuery(ctx.querystring),
        config
      };

      window.appState.setState('ui.currentRoute', ctx.pathname);

      // 设置页面标题
      if (config.title) {
        document.title = `${config.title} - 个人博客`;
      }

      // 加载页面模板
      await this.loadTemplate(config.template, this.currentRoute);

      // 更新导航栏状态
      this.updateNavigation();

      // 执行后置钩子
      for (const hook of this.afterHooks) {
        await hook(ctx, config);
      }

    } catch (error) {
      console.error('路由处理失败:', error);
      this.showNotification('页面加载失败', 'error');

      // 加载错误页面
      await this.loadTemplate('error.html', { error });
    } finally {
      // 隐藏加载状态
      window.appState.setLoading(false);
    }
  }

  // 加载页面模板
  async loadTemplate(templateName, routeData = {}) {
    const contentContainer = document.getElementById('page-content');
    if (!contentContainer) {
      throw new Error('页面内容容器不存在');
    }

    try {
      // 获取模板内容
      const templateUrl = `/templates/pages/${templateName}`;
      const response = await fetch(templateUrl);

      if (!response.ok) {
        throw new Error(`模板加载失败: ${response.status}`);
      }

      const templateContent = await response.text();

      // 渲染模板
      contentContainer.innerHTML = templateContent;

      // 执行模板内的脚本，确保 initXxxPage 已定义
      executeTemplateScripts(contentContainer);

      // 添加页面动画
      contentContainer.classList.add('fade-in');

      // 执行页面特定的初始化逻辑
      await this.initializePage(templateName, routeData);

      // 滚动到顶部
      window.scrollTo(0, 0);

    } catch (error) {
      console.error('模板加载失败:', error);

      // 显示错误页面
      contentContainer.innerHTML = `
        <div class="container py-5">
          <div class="text-center">
            <h2>页面加载失败</h2>
            <p class="text-muted">抱歉，页面暂时无法访问</p>
            <button class="btn btn-primary-custom" onclick="window.location.reload()">
              <i class="fas fa-refresh me-2"></i>重新加载
            </button>
          </div>
        </div>
      `;
    }
  }

  // 初始化页面特定逻辑
  async initializePage(templateName, routeData) {
    const pageName = templateName.replace('.html', '').replace('/', '-');
    const initFunction = `init${this.toPascalCase(pageName)}Page`;

    // 检查是否存在页面初始化函数
    if (typeof window[initFunction] === 'function') {
      try {
        await window[initFunction](routeData);
      } catch (error) {
        console.error(`页面初始化失败 (${initFunction}):`, error);
      }
    }
  }

  // 更新导航栏状态
  updateNavigation() {
    const navLinks = document.querySelectorAll('.navbar-nav .nav-link');
    const currentPath = this.currentRoute.path;

    navLinks.forEach(link => {
      link.classList.remove('active');

      const href = link.getAttribute('href');
      if (href && (href === currentPath ||
        (href !== '/' && currentPath.startsWith(href)))) {
        link.classList.add('active');
      }
    });

    // 更新用户状态显示
    this.updateUserNavigation();
  }

  // 更新用户导航状态
  updateUserNavigation() {
    const isAuthenticated = window.appState.isAuthenticated();
    const userProfile = window.appState.getState('user.profile');
    const isAdmin = window.appState.isAdmin();

    // 控制导航元素显示/隐藏
    const guestElements = document.querySelectorAll('#nav-guest, #nav-guest-register');
    const userElements = document.querySelectorAll('#nav-user');
    const adminElements = document.querySelectorAll('#nav-admin');

    if (isAuthenticated) {
      guestElements.forEach(el => el.classList.add('d-none'));
      userElements.forEach(el => el.classList.remove('d-none'));

      // 更新用户信息显示
      const avatarEl = document.getElementById('nav-avatar');
      const usernameEl = document.getElementById('nav-username');

      if (avatarEl && userProfile) {
        avatarEl.src = userProfile.avatar || '/assets/images/default-avatar.png';
        avatarEl.alt = userProfile.nickname || userProfile.username;
      }

      if (usernameEl && userProfile) {
        usernameEl.textContent = userProfile.nickname || userProfile.username;
      }

      // 控制管理员菜单
      if (isAdmin) {
        adminElements.forEach(el => el.classList.remove('d-none'));
      } else {
        adminElements.forEach(el => el.classList.add('d-none'));
      }

    } else {
      guestElements.forEach(el => el.classList.remove('d-none'));
      userElements.forEach(el => el.classList.add('d-none'));
      adminElements.forEach(el => el.classList.add('d-none'));
    }
  }

  // 解析查询参数
  parseQuery(queryString) {
    const params = new URLSearchParams(queryString);
    const result = {};
    for (const [key, value] of params) {
      result[key] = value;
    }
    return result;
  }

  // 字符串转PascalCase
  toPascalCase(str) {
    return str.replace(/-([a-z])/g, (g) => g[1].toUpperCase())
      .replace(/^[a-z]/, (g) => g.toUpperCase());
  }

  // 添加前置路由钩子
  beforeEach(hook) {
    this.beforeHooks.push(hook);
  }

  // 添加后置路由钩子
  afterEach(hook) {
    this.afterHooks.push(hook);
  }

  // 编程式导航
  push(path) {
    page(path);
  }

  // 替换当前路由
  replace(path) {
    page.replace(path);
  }

  // 返回上一页
  back() {
    window.history.back();
  }

  // 前进到下一页
  forward() {
    window.history.forward();
  }

  // 获取当前路由信息
  getCurrentRoute() {
    return this.currentRoute;
  }

  // 显示通知
  showNotification(message, type = 'info') {
    window.appState.addNotification({
      message,
      type,
      duration: 3000
    });
  }

  // 销毁路由系统
  destroy() {
    if (this.initialized) {
      page.stop();
      this.routes.clear();
      this.beforeHooks = [];
      this.afterHooks = [];
      this.currentRoute = null;
      this.initialized = false;
    }
  }
}

// 创建全局路由实例
window.router = new Router();

// 导出路由实例（如果使用模块系统）
if (typeof module !== 'undefined' && module.exports) {
  module.exports = Router;
}


// 显式执行模板中的 <script>，确保 window.initXxxPage 已定义
function executeTemplateScripts(container) {
  const scripts = Array.from(container.querySelectorAll('script'));
  scripts.forEach((oldScript) => {
    const newScript = document.createElement('script');
    if (oldScript.type) {
      newScript.type = oldScript.type;
    }
    if (oldScript.src) {
      // 外链脚本：复制 src 并按顺序同步执行
      newScript.src = oldScript.src;
      newScript.async = false;
    } else {
      // 内联脚本：复制文本内容
      newScript.textContent = oldScript.textContent;
    }
    document.body.appendChild(newScript);
    // 去除原脚本节点，避免重复解析
    if (oldScript.parentNode) {
      oldScript.parentNode.removeChild(oldScript);
    }
  });
}
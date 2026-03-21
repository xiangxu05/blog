// 应用主入口 - app.js

class BlogApp {
  constructor() {
    this.initialized = false;
    this.components = new Map();
    this.ux = null; // UX管理器实例
  }
  
  // 初始化应用
  async init() {
    if (this.initialized) return;
    
    try {
      console.log('初始化个人博客应用...');
      
      // 1. 初始化UX管理器
      this.ux = new UXManager();
      // 暴露通知系统以兼容内联调用
      this.notification = this.ux.notification;
      
      // 2. 初始化主题系统
      this.initTheme();
      
      // 3. 初始化事件监听器
      this.initEventListeners();
      
      // 4. 初始化通知系统
      this.initNotificationSystem();
      
      // 5. 检查用户登录状态
      await this.checkAuthStatus();
      
      // 6. 初始化路由系统
      window.router.init();
      
      // 7. 初始化全局组件
      this.initGlobalComponents();

      // 8. 全局搜索
      this.setupGlobalSearch();
      
      this.initialized = true;
      console.log('个人博客应用初始化完成');
      
      // 隐藏初始加载指示器
      this.hideInitialLoading();
      
    } catch (error) {
      console.error('应用初始化失败:', error);
      this.showError('应用初始化失败，请刷新页面重试');
    }
  }
  
  // 初始化主题系统
  initTheme() {
    const savedTheme = window.appState.getState('ui.theme');
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    
    // 设置初始主题
    const theme = savedTheme || (prefersDark ? 'dark' : 'light');
    window.appState.setTheme(theme);
    
    // 监听系统主题变化
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
      if (!window.appState.getState('ui.theme')) {
        window.appState.setTheme(e.matches ? 'dark' : 'light');
      }
    });
    
    console.log(`主题系统初始化完成: ${theme}`);
  }
  
  // 初始化事件监听器
  initEventListeners() {
    // 主题切换按钮
    document.addEventListener('click', (e) => {
      if (e.target.matches('#theme-toggle, #theme-toggle *')) {
        e.preventDefault();
        this.toggleTheme();
      }
    });
    
    // 用户注销
    document.addEventListener('click', (e) => {
      if (e.target.matches('#logout-btn')) {
        e.preventDefault();
        this.handleLogout();
      }
    });
    
    // 全局键盘快捷键
    document.addEventListener('keydown', (e) => {
      this.handleGlobalKeydown(e);
    });
    
    // 监听网络状态变化
    window.addEventListener('online', () => {
      this.showNotification('网络连接已恢复', 'success');
    });
    
    window.addEventListener('offline', () => {
      this.showNotification('网络连接已断开', 'warning');
    });
    
    console.log('全局事件监听器初始化完成');
  }
  
  // 初始化通知系统
  initNotificationSystem() {
    // 监听通知状态变化
    window.appState.subscribe('ui.notifications', (notifications) => {
      this.updateNotifications(notifications);
    });
    
    console.log('通知系统初始化完成');
  }
  
  // 检查用户认证状态
  async checkAuthStatus() {
    try {
      // 尝试获取用户信息来验证登录状态
      const userProfile = await window.apiClient.user.getProfile();
      
      // 更新用户状态
      window.appState.setUser(userProfile);
      
      console.log('用户认证状态: 已登录', userProfile);
      
    } catch (error) {
      // 如果获取用户信息失败，清除本地状态
      window.appState.clearUser();
      console.log('用户认证状态: 未登录');
    }
  }
  
  // 初始化全局组件
  initGlobalComponents() {
    // 注册全局组件
    this.registerComponent('navbar', new NavbarComponent());
    this.registerComponent('toast', new ToastComponent());
    this.registerComponent('modal', new ModalComponent());
    
    // 初始化所有组件
    this.components.forEach((component, name) => {
      try {
        if (typeof component.init === 'function') {
          component.init();
        }
      } catch (error) {
        console.error(`组件初始化失败 (${name}):`, error);
      }
    });
    
    console.log('全局组件初始化完成');
  }
  
  // 注册组件
  registerComponent(name, component) {
    this.components.set(name, component);
  }
  
  // 获取组件
  getComponent(name) {
    return this.components.get(name);
  }
  
  // 切换主题
  toggleTheme() {
    const currentTheme = window.appState.getState('ui.theme');
    const newTheme = currentTheme === 'light' ? 'dark' : 'light';
    
    // 添加切换动画
    document.body.setAttribute('data-theme-switching', 'true');
    
    setTimeout(() => {
      window.appState.setTheme(newTheme);
      
      // 更新主题切换按钮图标
      this.updateThemeToggleIcon(newTheme);
      
      setTimeout(() => {
        document.body.removeAttribute('data-theme-switching');
      }, 300);
    }, 50);
  }
  
  // 更新主题切换按钮图标
  updateThemeToggleIcon(theme) {
    const themeToggle = document.getElementById('theme-toggle');
    if (themeToggle) {
      const icon = themeToggle.querySelector('i');
      if (icon) {
        icon.className = theme === 'light' ? 'fas fa-moon' : 'fas fa-sun';
      }
    }
  }
  
  // 处理用户注销
  async handleLogout() {
    console.log('handleLogout 方法被调用');
    try {
      // 显示确认对话框
      console.log('显示确认对话框');
      if (!confirm('确定要退出登录吗？')) {
        console.log('用户取消了退出登录');
        return;
      }
      
      console.log('用户确认退出登录，开始调用API');
      // 显示加载状态
      window.appState.setLoading(true);
      
      // 调用注销API
      console.log('正在调用 logout API...');
      await window.apiClient.user.logout();
      console.log('logout API 调用成功');
      
      // 清除用户状态
      window.appState.clearUser();
      
      // 显示成功消息
      this.showNotification('已成功退出登录', 'success');
      
      // 重定向到首页
      window.router.push('/');
      
    } catch (error) {
      console.error('注销失败:', error);
      this.showNotification('注销失败，请重试', 'error');
    } finally {
      window.appState.setLoading(false);
    }
  }
  
  // 处理全局键盘快捷键
  handleGlobalKeydown(e) {
    // Ctrl/Cmd + K: 快速搜索
    if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
      e.preventDefault();
      this.openSearch();
    }
    
    // Escape: 关闭模态框/通知
    if (e.key === 'Escape') {
      this.closeModals();
    }
    
    // Alt + T: 切换主题
    if (e.altKey && e.key === 't') {
      e.preventDefault();
      this.toggleTheme();
    }
  }
  
  // 打开搜索（跳转文章列表并带 search 参数）
  openSearch() {
    const modalEl = document.getElementById('global-search-modal');
    if (!modalEl || typeof bootstrap === 'undefined') {
      if (typeof page === 'function') page('/articles');
      return;
    }
    const modal = bootstrap.Modal.getOrCreateInstance(modalEl);
    modal.show();
    setTimeout(() => {
      const input = document.getElementById('global-search-input');
      if (input) {
        input.focus();
        input.select();
      }
    }, 200);
  }

  setupGlobalSearch() {
    const form = document.getElementById('global-search-form');
    const input = document.getElementById('global-search-input');
    document.getElementById('navbar-search-btn')?.addEventListener('click', () => this.openSearch());
    form?.addEventListener('submit', (e) => {
      e.preventDefault();
      const q = (input?.value || '').trim();
      const modalEl = document.getElementById('global-search-modal');
      if (modalEl && typeof bootstrap !== 'undefined') {
        bootstrap.Modal.getInstance(modalEl)?.hide();
      }
      if (typeof page === 'function') {
        const qs = q ? `?search=${encodeURIComponent(q)}` : '';
        page('/articles' + qs);
      } else {
        window.location.href = '/articles' + (q ? `?search=${encodeURIComponent(q)}` : '');
      }
    });
  }
  
  // 关闭模态框
  closeModals() {
    // 关闭Bootstrap模态框
    const modals = document.querySelectorAll('.modal.show');
    modals.forEach(modal => {
      const modalInstance = bootstrap.Modal.getInstance(modal);
      if (modalInstance) {
        modalInstance.hide();
      }
    });
    
    // 关闭通知
    const toasts = document.querySelectorAll('.toast.show');
    toasts.forEach(toast => {
      const toastInstance = bootstrap.Toast.getInstance(toast);
      if (toastInstance) {
        toastInstance.hide();
      }
    });
  }
  
  // 更新通知显示
  updateNotifications(notifications) {
    const latestNotification = notifications[notifications.length - 1];
    if (latestNotification) {
      this.showToast(latestNotification);
    }
  }
  
  // 显示Toast通知
  showToast(notification) {
    const toastEl = document.getElementById('toast');
    const toastIcon = document.getElementById('toast-icon');
    const toastTitle = document.getElementById('toast-title');
    const toastBody = document.getElementById('toast-body');
    
    if (!toastEl) return;
    
    // 设置通知内容
    if (toastIcon) {
      const iconClass = this.getNotificationIcon(notification.type);
      toastIcon.className = iconClass;
    }
    
    if (toastTitle) {
      toastTitle.textContent = this.getNotificationTitle(notification.type);
    }
    
    if (toastBody) {
      toastBody.textContent = notification.message;
    }
    
    // 设置通知样式
    toastEl.className = `toast ${this.getNotificationClass(notification.type)}`;
    
    // 显示通知
    const toast = new bootstrap.Toast(toastEl, {
      delay: notification.duration || 3000
    });
    toast.show();
    
    // 自动移除通知
    setTimeout(() => {
      window.appState.removeNotification(notification.id);
    }, notification.duration || 3000);
  }
  
  // 获取通知图标
  getNotificationIcon(type) {
    const icons = {
      success: 'fas fa-check-circle text-success',
      error: 'fas fa-exclamation-circle text-danger',
      warning: 'fas fa-exclamation-triangle text-warning',
      info: 'fas fa-info-circle text-info'
    };
    return icons[type] || icons.info;
  }
  
  // 获取通知标题
  getNotificationTitle(type) {
    const titles = {
      success: '成功',
      error: '错误',
      warning: '警告',
      info: '提示'
    };
    return titles[type] || titles.info;
  }
  
  // 获取通知样式类
  getNotificationClass(type) {
    const classes = {
      success: 'border-success',
      error: 'border-danger',
      warning: 'border-warning',
      info: 'border-info'
    };
    return classes[type] || classes.info;
  }
  
  // 显示通知（便捷方法）
  showNotification(message, type = 'info', duration = 3000) {
    window.appState.addNotification({
      message,
      type,
      duration
    });
  }
  
  // 隐藏初始加载指示器
  hideInitialLoading() {
    const loadingEl = document.getElementById('loading');
    if (loadingEl) {
      loadingEl.style.display = 'none';
    }
  }
  
  // 显示错误信息
  showError(message) {
    const contentContainer = document.getElementById('page-content');
    if (contentContainer) {
      contentContainer.innerHTML = `
        <div class="container py-5">
          <div class="alert alert-danger-custom">
            <h4><i class="fas fa-exclamation-circle me-2"></i>应用错误</h4>
            <p>${message}</p>
            <button class="btn btn-primary-custom" onclick="window.location.reload()">
              重新加载
            </button>
          </div>
        </div>
      `;
    }
  }
}

// 简单的导航栏组件
class NavbarComponent {
  init() {
    // 监听用户状态变化以更新导航栏
    window.appState.subscribe('user.*', () => {
      this.updateNavbar();
    });
  }
  
  updateNavbar() {
    // 路由系统会处理导航栏更新
    if (window.router && window.router.updateUserNavigation) {
      window.router.updateUserNavigation();
    }
  }
}

// 简单的Toast组件
class ToastComponent {
  init() {
    // Toast 通知已在应用主类中处理
  }
}

// 简单的模态框组件
class ModalComponent {
  init() {
    // 模态框功能由Bootstrap提供
  }
}

// 创建应用实例
const app = new BlogApp();

// DOM加载完成后初始化应用
document.addEventListener('DOMContentLoaded', () => {
  app.init().catch(error => {
    console.error('应用启动失败:', error);
  });
});

// 全局暴露应用实例
window.blogApp = app;

// 导出应用类（如果使用模块系统）
if (typeof module !== 'undefined' && module.exports) {
  module.exports = BlogApp;
}
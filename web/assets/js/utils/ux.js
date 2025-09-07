// 用户体验优化组件 - ux.js

// 通知系统
class NotificationSystem {
  constructor() {
    this.container = null;
    this.notifications = new Map();
    this.counter = 0;
    this.init();
  }

  init() {
    // 创建通知容器
    this.container = document.createElement('div');
    this.container.className = 'notification-container';
    this.container.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      z-index: 9999;
      pointer-events: none;
    `;
    document.body.appendChild(this.container);
  }

  show(message, type = 'info', duration = 4000, options = {}) {
    const id = ++this.counter;
    const notification = this.createNotification(id, message, type, options);

    this.container.appendChild(notification);
    this.notifications.set(id, notification);

    // 触发显示动画
    requestAnimationFrame(() => {
      notification.classList.add('show');
    });

    // 自动隐藏
    if (duration > 0) {
      setTimeout(() => {
        this.hide(id);
      }, duration);
    }

    return id;
  }

  createNotification(id, message, type, options) {
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.style.pointerEvents = 'auto';

    const iconMap = {
      success: 'fas fa-check-circle',
      error: 'fas fa-exclamation-circle',
      warning: 'fas fa-exclamation-triangle',
      info: 'fas fa-info-circle'
    };

    notification.innerHTML = `
      <div class="notification-content">
        <div class="notification-icon">
          <i class="${iconMap[type] || iconMap.info}"></i>
        </div>
        <div class="notification-message">${message}</div>
        <button class="notification-close" onclick="window.blogApp.notification.hide(${id})">
          <i class="fas fa-times"></i>
        </button>
      </div>
    `;

    // 添加点击关闭事件
    if (options.clickToClose !== false) {
      notification.addEventListener('click', (e) => {
        if (!e.target.closest('.notification-close')) {
          this.hide(id);
        }
      });
    }

    return notification;
  }

  hide(id) {
    const notification = this.notifications.get(id);
    if (!notification) return;

    notification.classList.remove('show');
    notification.addEventListener('transitionend', () => {
      if (notification.parentNode) {
        notification.parentNode.removeChild(notification);
      }
      this.notifications.delete(id);
    }, { once: true });
  }

  clear() {
    this.notifications.forEach((notification, id) => {
      this.hide(id);
    });
  }
}

// 加载状态管理器
class LoadingManager {
  constructor() {
    this.loadingCount = 0;
    this.overlay = null;
    this.init();
  }

  init() {
    // 创建全局加载遮罩
    this.overlay = document.createElement('div');
    this.overlay.className = 'loading-overlay';
    this.overlay.style.display = 'none';
    this.overlay.innerHTML = `
      <div class="loading-content text-center">
        <div class="loading-spinner"></div>
        <div class="loading-text text-white mt-3">加载中...</div>
      </div>
    `;
    document.body.appendChild(this.overlay);
  }

  show(text = '加载中...') {
    this.loadingCount++;
    this.overlay.querySelector('.loading-text').textContent = text;
    this.overlay.style.display = 'flex';
    document.body.style.overflow = 'hidden';
  }

  hide() {
    this.loadingCount = Math.max(0, this.loadingCount - 1);
    if (this.loadingCount === 0) {
      this.overlay.style.display = 'none';
      document.body.style.overflow = '';
    }
  }

  isLoading() {
    return this.loadingCount > 0;
  }
}

// 错误处理器
class ErrorHandler {
  constructor() {
    this.errorLog = [];
    this.init();
  }

  init() {
    // 监听全局错误
    window.addEventListener('error', (e) => {
      this.handleError(e.error, '页面错误');
    });

    window.addEventListener('unhandledrejection', (e) => {
      this.handleError(e.reason, 'Promise 错误');
    });
  }

  handleError(error, context = '未知错误') {
    // 安全检查：确保error不为null或undefined
    const safeError = error || new Error('未知错误');

    const errorInfo = {
      message: safeError.message || String(safeError),
      stack: safeError.stack,
      context,
      timestamp: Date.now(),
      url: window.location.href,
      userAgent: navigator.userAgent
    };

    // 记录错误
    this.errorLog.push(errorInfo);

    // 控制台输出
    console.error(`[${context}]`, safeError);

    // 用户友好的错误提示
    this.showUserFriendlyError(safeError, context);

    // 可选：发送错误到服务器
    this.reportError(errorInfo);
  }

  showUserFriendlyError(error, context) {
    let message = '发生了一个错误，请稍后重试';

    // 根据错误类型显示不同消息
    if (error.message && error.message.includes('网络')) {
      message = '网络连接失败，请检查网络后重试';
    } else if (error.message && error.message.includes('权限')) {
      message = '权限不足，请登录后重试';
    } else if (error.message && error.message.includes('404')) {
      message = '请求的资源不存在';
    }

    // 安全检查：确保notification对象存在
    if (window.blogApp && window.blogApp.notification && typeof window.blogApp.notification.show === 'function') {
      window.blogApp.notification.show(message, 'error', 5000);
    } else {
      // 降级处理：使用原生alert
      console.error('通知系统未初始化，使用console.error输出:', message);
      console.error('原始错误:', error);
    }
  }

  reportError(errorInfo) {
    // 这里可以发送错误到服务器
    // 避免在错误报告时再次出错
    try {
      // fetch('/api/errors', {
      //   method: 'POST',
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify(errorInfo)
      // }).catch(() => {}); // 静默处理报告失败
    } catch (e) {
      // 静默处理
    }
  }

  getErrorLog() {
    return [...this.errorLog];
  }

  clearErrorLog() {
    this.errorLog = [];
  }
}

// 性能监控器
class PerformanceMonitor {
  constructor() {
    this.metrics = {
      pageLoadTime: 0,
      apiCalls: [],
      memoryUsage: 0,
      renderTime: 0
    };
    this.init();
  }

  init() {
    // 监控页面加载时间
    window.addEventListener('load', () => {
      const loadTime = performance.timing.loadEventEnd - performance.timing.navigationStart;
      this.metrics.pageLoadTime = loadTime;
      // console.log(`页面加载时间: ${loadTime}ms`);
    });

    // 定期监控内存使用情况
    if (performance.memory) {
      setInterval(() => {
        this.metrics.memoryUsage = performance.memory.usedJSHeapSize;
      }, 30000);
    }
  }

  recordApiCall(url, startTime, endTime, success) {
    const duration = endTime - startTime;
    this.metrics.apiCalls.push({
      url,
      duration,
      success,
      timestamp: Date.now()
    });

    // 只保留最近100次调用记录
    if (this.metrics.apiCalls.length > 100) {
      this.metrics.apiCalls.shift();
    }

    if (duration > 3000) {
      console.warn(`API调用缓慢: ${url} (${duration}ms)`);
    }
  }

  getMetrics() {
    return { ...this.metrics };
  }
}

// 用户体验辅助功能
class UXHelpers {
  constructor() {
    this.init();
  }

  init() {
    this.initScrollToTop();
    this.initBackButton();
    this.initKeyboardNavigation();
    this.initTooltips();
  }

  // 回到顶部按钮
  initScrollToTop() {
    const scrollToTopBtn = document.createElement('button');
    scrollToTopBtn.className = 'scroll-to-top-btn';
    scrollToTopBtn.innerHTML = '<i class="fas fa-arrow-up"></i>';
    scrollToTopBtn.title = '返回顶部';
    scrollToTopBtn.style.cssText = `
      position: fixed;
      bottom: 30px;
      right: 30px;
      width: 50px;
      height: 50px;
      border-radius: 50%;
      background-color: var(--color-primary);
      color: white;
      border: none;
      box-shadow: var(--shadow-lg);
      z-index: 1000;
      opacity: 0;
      visibility: hidden;
      transition: all var(--transition-normal) var(--ease-out);
      cursor: pointer;
    `;

    scrollToTopBtn.addEventListener('click', () => {
      window.scrollTo({
        top: 0,
        behavior: 'smooth'
      });
    });

    document.body.appendChild(scrollToTopBtn);

    // 监听滚动事件
    let ticking = false;
    window.addEventListener('scroll', () => {
      if (!ticking) {
        requestAnimationFrame(() => {
          const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
          if (scrollTop > 300) {
            scrollToTopBtn.style.opacity = '1';
            scrollToTopBtn.style.visibility = 'visible';
          } else {
            scrollToTopBtn.style.opacity = '0';
            scrollToTopBtn.style.visibility = 'hidden';
          }
          ticking = false;
        });
        ticking = true;
      }
    });
  }

  // 浏览器后退按钮处理
  initBackButton() {
    let isNavigating = false;

    window.addEventListener('beforeunload', (e) => {
      if (window.appState && window.appState.isDirty && !isNavigating) {
        e.preventDefault();
        e.returnValue = '您有未保存的更改，确定要离开吗？';
        return e.returnValue;
      }
    });

    // 监听路由变化
    window.addEventListener('popstate', () => {
      isNavigating = true;
      setTimeout(() => {
        isNavigating = false;
      }, 100);
    });
  }

  // 键盘导航
  initKeyboardNavigation() {
    document.addEventListener('keydown', (e) => {
      // ESC 键关闭模态框等
      if (e.key === 'Escape') {
        // 关闭打开的模态框
        const openModals = document.querySelectorAll('.modal.show');
        openModals.forEach(modal => {
          const bsModal = bootstrap.Modal.getInstance(modal);
          if (bsModal) {
            bsModal.hide();
          }
        });

        // 清除通知
        if (e.ctrlKey || e.metaKey) {
          window.blogApp.notification.clear();
        }
      }

      // Ctrl/Cmd + K 聚焦搜索框
      if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        const searchInput = document.querySelector('input[type="search"], input[placeholder*="搜索"]');
        if (searchInput) {
          searchInput.focus();
        }
      }
    });
  }

  // 工具提示初始化
  initTooltips() {
    // 为所有带有 title 属性的元素添加工具提示
    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (mutation.type === 'childList') {
          mutation.addedNodes.forEach((node) => {
            if (node.nodeType === Node.ELEMENT_NODE) {
              this.initElementTooltips(node);
            }
          });
        }
      });
    });

    observer.observe(document.body, {
      childList: true,
      subtree: true
    });

    // 初始化现有元素
    this.initElementTooltips(document.body);
  }

  initElementTooltips(element) {
    const elementsWithTooltips = element.querySelectorAll('[title]:not([data-tooltip-initialized])');
    elementsWithTooltips.forEach(el => {
      el.setAttribute('data-tooltip-initialized', 'true');
      // 这里可以使用 Bootstrap 的 tooltip 或自定义实现
      if (window.bootstrap && bootstrap.Tooltip) {
        new bootstrap.Tooltip(el);
      }
    });
  }
}

// 导出UX管理器
class UXManager {
  constructor() {
    this.notification = new NotificationSystem();
    this.loading = new LoadingManager();
    this.error = new ErrorHandler();
    this.performance = new PerformanceMonitor();
    this.helpers = new UXHelpers();
  }

  // 便捷方法
  showNotification(message, type, duration) {
    return this.notification.show(message, type, duration);
  }

  showLoading(text) {
    this.loading.show(text);
  }

  hideLoading() {
    this.loading.hide();
  }

  handleError(error, context) {
    this.error.handleError(error, context);
  }

  // 显示确认对话框
  confirm(message, title = '确认') {
    return new Promise((resolve) => {
      // 创建自定义确认对话框
      const modal = document.createElement('div');
      modal.className = 'modal fade';
      modal.innerHTML = `
        <div class="modal-dialog modal-sm">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">${title}</h5>
            </div>
            <div class="modal-body">
              <p>${message}</p>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary-custom" data-bs-dismiss="modal">取消</button>
              <button type="button" class="btn btn-primary-custom" id="confirm-ok">确认</button>
            </div>
          </div>
        </div>
      `;

      document.body.appendChild(modal);
      const bsModal = new bootstrap.Modal(modal);

      modal.querySelector('#confirm-ok').addEventListener('click', () => {
        bsModal.hide();
        resolve(true);
      });

      modal.addEventListener('hidden.bs.modal', () => {
        document.body.removeChild(modal);
        resolve(false);
      });

      bsModal.show();
    });
  }

  // 显示提示对话框
  alert(message, title = '提示', type = 'info') {
    return new Promise((resolve) => {
      const iconMap = {
        info: 'fas fa-info-circle text-primary',
        success: 'fas fa-check-circle text-success',
        warning: 'fas fa-exclamation-triangle text-warning',
        error: 'fas fa-exclamation-circle text-danger'
      };

      const modal = document.createElement('div');
      modal.className = 'modal fade';
      modal.innerHTML = `
        <div class="modal-dialog modal-sm">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">
                <i class="${iconMap[type]} me-2"></i>
                ${title}
              </h5>
            </div>
            <div class="modal-body">
              <p>${message}</p>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-primary-custom" data-bs-dismiss="modal">确定</button>
            </div>
          </div>
        </div>
      `;

      document.body.appendChild(modal);
      const bsModal = new bootstrap.Modal(modal);

      modal.addEventListener('hidden.bs.modal', () => {
        document.body.removeChild(modal);
        resolve();
      });

      bsModal.show();
    });
  }

  // 防抖函数
  debounce(func, wait, immediate) {
    let timeout;
    return function executedFunction(...args) {
      const later = () => {
        timeout = null;
        if (!immediate) func.apply(this, args);
      };
      const callNow = immediate && !timeout;
      clearTimeout(timeout);
      timeout = setTimeout(later, wait);
      if (callNow) func.apply(this, args);
    };
  }

  // 节流函数
  throttle(func, limit) {
    let inThrottle;
    return function executedFunction(...args) {
      if (!inThrottle) {
        func.apply(this, args);
        inThrottle = true;
        setTimeout(() => inThrottle = false, limit);
      }
    };
  }

  // 复制到剪贴板
  async copyToClipboard(text) {
    try {
      await navigator.clipboard.writeText(text);
      this.showNotification('已复制到剪贴板', 'success', 2000);
      return true;
    } catch (err) {
      console.error('复制失败:', err);
      this.showNotification('复制失败', 'error', 2000);
      return false;
    }
  }

  // 格式化文件大小
  formatFileSize(bytes) {
    const sizes = ['B', 'KB', 'MB', 'GB'];
    if (bytes === 0) return '0 B';
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i];
  }

  // 格式化时间
  formatTime(timestamp) {
    const date = new Date(timestamp * 1000);
    const now = new Date();
    const diff = now - date;

    if (diff < 60000) {
      return '刚刚';
    } else if (diff < 3600000) {
      return `${Math.floor(diff / 60000)} 分钟前`;
    } else if (diff < 86400000) {
      return `${Math.floor(diff / 3600000)} 小时前`;
    } else if (diff < 604800000) {
      return `${Math.floor(diff / 86400000)} 天前`;
    } else {
      return date.toLocaleDateString('zh-CN');
    }
  }

  // 检测设备类型
  getDeviceType() {
    const width = window.innerWidth;
    if (width < 576) return 'mobile';
    if (width < 768) return 'mobile-large';
    if (width < 992) return 'tablet';
    if (width < 1200) return 'desktop';
    return 'desktop-large';
  }

  // 平滑滚动到元素
  scrollToElement(element, offset = 0) {
    if (typeof element === 'string') {
      element = document.querySelector(element);
    }

    if (element) {
      const top = element.offsetTop - offset;
      window.scrollTo({
        top,
        behavior: 'smooth'
      });
    }
  }
}

// 如果在浏览器环境中，添加到全局
if (typeof window !== 'undefined') {
  window.UXManager = UXManager;
}
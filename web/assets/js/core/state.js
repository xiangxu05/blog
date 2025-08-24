// 全局状态管理 - state.js

class AppState {
  constructor() {
    // 初始化状态
    this.state = {
      // 用户状态
      user: {
        isAuthenticated: false,
        profile: null,
        role: null // 'user' | 'admin'
      },
      
      // 文章状态
      articles: {
        list: [],
        current: null,
        currentContent: '',
        pagination: {
          page: 1,
          pageSize: 10,
          total: 0,
          totalPages: 0
        },
        loading: false
      },
      
      // 文件状态
      files: {
        uploadQueue: [],
        list: [],
        loading: false
      },
      
      // UI状态
      ui: {
        theme: 'light',
        loading: false,
        notifications: [],
        currentRoute: '',
        sidebarOpen: false
      }
    };
    
    // 状态变更监听器
    this.listeners = {};
    
    // 从localStorage恢复状态
    this.restoreState();
  }
  
  // 获取状态
  getState(path) {
    if (!path) return this.state;
    
    const keys = path.split('.');
    let current = this.state;
    
    for (const key of keys) {
      if (current[key] === undefined) return undefined;
      current = current[key];
    }
    
    return current;
  }
  
  // 设置状态
  setState(path, value) {
    const keys = path.split('.');
    let current = this.state;
    
    // 导航到父级对象
    for (let i = 0; i < keys.length - 1; i++) {
      if (current[keys[i]] === undefined) {
        current[keys[i]] = {};
      }
      current = current[keys[i]];
    }
    
    // 设置值
    const lastKey = keys[keys.length - 1];
    const oldValue = current[lastKey];
    current[lastKey] = value;
    
    // 触发监听器
    this.notifyListeners(path, value, oldValue);
    
    // 持久化特定状态
    this.persistState();
  }
  
  // 批量更新状态
  batchUpdate(updates) {
    Object.entries(updates).forEach(([path, value]) => {
      this.setState(path, value);
    });
  }
  
  // 添加状态监听器
  subscribe(path, callback) {
    if (!this.listeners[path]) {
      this.listeners[path] = [];
    }
    this.listeners[path].push(callback);
    
    // 返回取消订阅函数
    return () => {
      this.listeners[path] = this.listeners[path].filter(cb => cb !== callback);
    };
  }
  
  // 通知监听器
  notifyListeners(path, newValue, oldValue) {
    // 精确路径匹配
    if (this.listeners[path]) {
      this.listeners[path].forEach(callback => {
        callback(newValue, oldValue, path);
      });
    }
    
    // 通配符匹配 (例如 'user.*' 匹配 'user.profile')
    Object.keys(this.listeners).forEach(listenerPath => {
      if (listenerPath.endsWith('*')) {
        const prefix = listenerPath.slice(0, -1);
        if (path.startsWith(prefix)) {
          this.listeners[listenerPath].forEach(callback => {
            callback(newValue, oldValue, path);
          });
        }
      }
    });
  }
  
  // 持久化状态到localStorage
  persistState() {
    const persistentData = {
      user: this.state.user,
      ui: {
        theme: this.state.ui.theme
      }
    };
    
    try {
      localStorage.setItem('blogAppState', JSON.stringify(persistentData));
    } catch (error) {
      console.warn('无法保存状态到localStorage:', error);
    }
  }
  
  // 从localStorage恢复状态
  restoreState() {
    try {
      const savedState = localStorage.getItem('blogAppState');
      if (savedState) {
        const parsed = JSON.parse(savedState);
        
        // 恢复用户状态
        if (parsed.user) {
          this.state.user = { ...this.state.user, ...parsed.user };
        }
        
        // 恢复UI状态
        if (parsed.ui) {
          this.state.ui = { ...this.state.ui, ...parsed.ui };
        }
      }
    } catch (error) {
      console.warn('无法从localStorage恢复状态:', error);
    }
  }
  
  // 清除状态
  clearState() {
    this.state = {
      user: {
        isAuthenticated: false,
        profile: null,
        role: null
      },
      articles: {
        list: [],
        current: null,
        currentContent: '',
        pagination: {
          page: 1,
          pageSize: 10,
          total: 0,
          totalPages: 0
        },
        loading: false
      },
      files: {
        uploadQueue: [],
        list: [],
        loading: false
      },
      ui: {
        theme: this.state.ui.theme, // 保持主题设置
        loading: false,
        notifications: [],
        currentRoute: '',
        sidebarOpen: false
      }
    };
    
    // 清除localStorage
    try {
      localStorage.removeItem('blogAppState');
    } catch (error) {
      console.warn('无法清除localStorage:', error);
    }
  }
  
  // 用户相关的便捷方法
  setUser(userProfile) {
    this.batchUpdate({
      'user.isAuthenticated': true,
      'user.profile': userProfile,
      'user.role': userProfile.role || 'user'
    });
  }
  
  clearUser() {
    this.batchUpdate({
      'user.isAuthenticated': false,
      'user.profile': null,
      'user.role': null
    });
  }
  
  isAuthenticated() {
    return this.getState('user.isAuthenticated');
  }
  
  isAdmin() {
    return this.getState('user.role') === 'admin';
  }
  
  // 文章相关的便捷方法
  setArticles(articles, pagination = null) {
    const updates = {
      'articles.list': articles,
      'articles.loading': false
    };
    
    if (pagination) {
      updates['articles.pagination'] = {
        ...this.getState('articles.pagination'),
        ...pagination
      };
    }
    
    this.batchUpdate(updates);
  }
  
  setCurrentArticle(article, content = '') {
    this.batchUpdate({
      'articles.current': article,
      'articles.currentContent': content
    });
  }
  
  // 通知相关的便捷方法
  addNotification(notification) {
    const notifications = this.getState('ui.notifications');
    const newNotifications = [...notifications, {
      id: Date.now(),
      timestamp: new Date(),
      ...notification
    }];
    this.setState('ui.notifications', newNotifications);
  }
  
  removeNotification(id) {
    const notifications = this.getState('ui.notifications');
    const newNotifications = notifications.filter(n => n.id !== id);
    this.setState('ui.notifications', newNotifications);
  }
  
  // 主题相关的便捷方法
  setTheme(theme) {
    this.setState('ui.theme', theme);
    document.body.setAttribute('data-theme', theme);
  }
  
  toggleTheme() {
    const currentTheme = this.getState('ui.theme');
    const newTheme = currentTheme === 'light' ? 'dark' : 'light';
    this.setTheme(newTheme);
  }
  
  // 加载状态管理
  setLoading(loading) {
    this.setState('ui.loading', loading);
  }
  
  setArticlesLoading(loading) {
    this.setState('articles.loading', loading);
  }
  
  setFilesLoading(loading) {
    this.setState('files.loading', loading);
  }
}

// 创建全局状态实例
window.appState = new AppState();

// 调试用：在控制台暴露状态管理
if (typeof window !== 'undefined') {
  window.debugState = {
    getState: (path) => window.appState.getState(path),
    setState: (path, value) => window.appState.setState(path, value),
    clearState: () => window.appState.clearState(),
    subscribe: (path, callback) => window.appState.subscribe(path, callback)
  };
}
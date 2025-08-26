// API 客户端 - api.js

class ApiClient {
  constructor() {
    this.baseURL = `${window.location.origin}/api`;
    this.defaultHeaders = {
      'Content-Type': 'application/json'
    };
  }

  // 通用请求方法
  async request(endpoint, options = {}) {
    const config = {
      credentials: 'include', // 自动携带Session Cookie
      headers: {
        ...this.defaultHeaders,
        ...options.headers
      },
      ...options
    };

    // 如果是FormData，移除Content-Type让浏览器自动设置
    if (config.body instanceof FormData) {
      delete config.headers['Content-Type'];
    }

    try {
      const response = await fetch(`${this.baseURL}${endpoint}`, config);

      // 检查响应状态
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ msg: '请求失败' }));
        throw new ApiError(errorData.msg || '请求失败', response.status, errorData);
      }

      // 解析响应数据
      const data = await response.json();
      return data.data; // 后端统一返回格式中的data字段

    } catch (error) {
      console.error('API请求错误:', error);

      // 处理网络错误
      if (error instanceof TypeError && error.message.includes('fetch')) {
        throw new ApiError('网络连接失败，请检查网络连接', 0);
      }

      // 处理会话过期
      if (error.status === 401) {
        // 清除用户状态
        window.appState.clearUser();
        // 显示登录提示
        this.showNotification('登录已过期，请重新登录', 'warning');
        // 重定向到登录页面（如果不在登录页面）
        if (!window.location.pathname.includes('/login')) {
          page('/login');
        }
      }

      throw error;
    }
  }

  // 显示通知的辅助方法
  showNotification(message, type = 'info') {
    window.appState.addNotification({
      message,
      type,
      duration: 3000
    });
  }

  // 用户相关API
  user = {
    // 用户注册
    register: (userData) => {
      return this.request('/users/register', {
        method: 'POST',
        body: JSON.stringify({
          username: userData.username,
          password: userData.password,
          confirm_password: userData.confirmPassword,
          email: userData.email
        })
      });
    },

    // 用户登录
    login: (credentials) => {
      return this.request('/users/login', {
        method: 'POST',
        body: JSON.stringify({
          username: credentials.username,
          password: credentials.password
        })
      });
    },

    // 用户注销
    logout: () => {
      return this.request('/users/logout');
    },

    // 获取用户资料
    getProfile: () => {
      return this.request('/users/profile');
    },

    // 更新用户资料
    updateProfile: (profileData) => {
      return this.request('/users/profile', {
        method: 'PUT',
        body: JSON.stringify(profileData)
      });
    },

    // 修改密码
    updatePassword: (passwordData) => {
      return this.request('/users/password', {
        method: 'PUT',
        body: JSON.stringify({
          old_password: passwordData.oldPassword,
          new_password: passwordData.newPassword,
          confirm_password: passwordData.confirmPassword
        })
      });
    },

    // 获取其他用户信息
    getUserInfo: (userId) => {
      return this.request(`/users/${userId}`);
    },

    // 删除用户账号
    deleteAccount: () => {
      return this.request('/users/delete', {
        method: 'DELETE'
      });
    },

    // 刷新会话
    refreshSession: () => {
      return this.request('/users/refresh');
    }
  };

  // 文章相关API
  articles = {
    // 获取文章列表
    getList: (params = {}) => {
      const queryParams = new URLSearchParams();

      if (params.page) queryParams.append('page', params.page);
      if (params.pageSize) queryParams.append('pageSize', params.pageSize);
      if (params.category) queryParams.append('category', params.category);
      if (params.search) queryParams.append('search', params.search);
      if (params.search_type) queryParams.append('search_type', params.search_type);
      if (params.search_fields) queryParams.append('search_fields', params.search_fields);
      if (params.sort) queryParams.append('sort', params.sort);

      const queryString = queryParams.toString();
      const endpoint = queryString ? `/articles?${queryString}` : '/articles';

      return this.request(endpoint);
    },

    // 获取文章详情
    getDetail: (id, version) => {
      return this.request(`/articles/${id}/${version}`);
    },

    // 创建文章
    create: (articleData) => {
      return this.request('/articles', {
        method: 'POST',
        body: JSON.stringify({
          title: articleData.title,
          description: articleData.description,
          category: articleData.category,
          tags: articleData.tags,
          store_id: articleData.store_id
        })
      });
    },

    // 更新文章
    update: (id, articleData) => {
      return this.request(`/articles/${id}`, {
        method: 'PUT',
        body: JSON.stringify({
          title: articleData.title,
          description: articleData.description,
          category: articleData.category,
          tags: articleData.tags,
          store_id: articleData.store_id
        })
      });
    },

    // 删除文章
    delete: (id) => {
      return this.request(`/articles/${id}`, {
        method: 'DELETE'
      });
    },

    // 获取热门文章
    getHot: () => {
      return this.request('/articles/hot');
    }
  };

  // 文件相关API
  files = {
    // 上传文件
    upload: (file, onProgress = null) => {
      const formData = new FormData();
      formData.append('file', file);

      return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();

        // 监听上传进度
        if (onProgress && xhr.upload) {
          xhr.upload.addEventListener('progress', (e) => {
            if (e.lengthComputable) {
              const percentComplete = (e.loaded / e.total) * 100;
              onProgress(percentComplete);
            }
          });
        }

        xhr.onload = () => {
          if (xhr.status >= 200 && xhr.status < 300) {
            try {
              const response = JSON.parse(xhr.responseText);
              resolve(response.data);
            } catch (error) {
              reject(new ApiError('响应解析失败', xhr.status));
            }
          } else {
            try {
              const errorData = JSON.parse(xhr.responseText);
              reject(new ApiError(errorData.msg || '上传失败', xhr.status, errorData));
            } catch (error) {
              reject(new ApiError('上传失败', xhr.status));
            }
          }
        };

        xhr.onerror = () => {
          reject(new ApiError('网络错误', 0));
        };

        xhr.open('POST', `${this.baseURL}/files/upload`);
        xhr.withCredentials = true; // 发送cookies
        xhr.send(formData);
      });
    },

    // 下载文件内容
    download: async (fileId) => {
      const config = {
        credentials: 'include',
        headers: this.defaultHeaders
      };

      try {
        const response = await fetch(`${this.baseURL}/files/download/${fileId}`, config);

        if (!response.ok) {
          const errorData = await response.json().catch(() => ({ msg: '文件下载失败' }));
          throw new ApiError(errorData.msg || '文件下载失败', response.status, errorData);
        }

        // 检查响应的Content-Type
        const contentType = response.headers.get('content-type');
        if (contentType && contentType.includes('application/json')) {
          // 如果是JSON响应，按原来的方式处理
          const data = await response.json();
          return data.data;
        } else {
          // 如果是文件下载，直接返回文本内容
          return await response.text();
        }

      } catch (error) {
        console.error('文件下载错误:', error);
        if (error instanceof TypeError && error.message.includes('fetch')) {
          throw new ApiError('网络连接失败，请检查网络连接', 0);
        }
        throw error;
      }
    },

    // 获取文件信息
    getInfo: (fileId) => {
      return this.request(`/files/${fileId}`);
    },

    // 获取文件列表（管理员）
    getList: () => {
      return this.request('/files');
    },

    // 删除文件（管理员）
    delete: (fileId) => {
      return this.request(`/files/${fileId}`, {
        method: 'DELETE'
      });
    },

    // 获取文件下载URL
    getDownloadUrl: (fileId) => {
      return `${this.baseURL}/files/download/${fileId}`;
    }
  };

  // 管理员相关API
  admin = {
    // 获取管理员文章列表
    getArticles: (params = {}) => {
      const queryParams = new URLSearchParams();

      if (params.page) queryParams.append('page', params.page);
      if (params.pageSize) queryParams.append('pageSize', params.pageSize);
      if (params.status) queryParams.append('status', params.status);
      if (params.category) queryParams.append('category', params.category);
      if (params.search) queryParams.append('search', params.search);
      if (params.search_type) queryParams.append('search_type', params.search_type);
      if (params.search_fields) queryParams.append('search_fields', params.search_fields);
      if (params.sort) queryParams.append('sort', params.sort);

      const queryString = queryParams.toString();
      const endpoint = queryString ? `/admin/articles?${queryString}` : '/admin/articles';

      return this.request(endpoint);
    }
  };

  // 分类相关API
  categories = {
    // 获取分类列表
    getList: () => {
      return this.request('/category');
    }
  };

  // 健康检查
  health = {
    check: () => {
      return this.request('/health');
    }
  };
}

// API错误类
class ApiError extends Error {
  constructor(message, status, data = null) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.data = data;
  }
}

// 文章服务类 - 处理文章的复杂业务逻辑
class ArticleService {
  constructor(apiClient) {
    this.api = apiClient;
  }

  // 获取文章详情（包含内容）- 两步加载机制
  async getArticleDetail(id, version) {
    try {
      // 第一步：获取文章元数据
      const articleMeta = await this.api.articles.getDetail(id, version);

      // 第二步：根据store_id下载文章内容
      if (articleMeta.store_id) {
        const content = await this.api.files.download(articleMeta.store_id);
        return {
          ...articleMeta,
          content: content
        };
      }

      return articleMeta;
    } catch (error) {
      console.error('获取文章详情失败:', error);
      throw error;
    }
  }

  // 发布文章 - 完整流程
  async publishArticle(title, description, markdownContent, category = '', tags = '') {
    try {
      // 第一步：上传Markdown内容作为文件
      const contentBlob = new Blob([markdownContent], { type: 'text/markdown' });
      const contentFile = new File([contentBlob], 'article.md', {
        type: 'text/markdown'
      });

      const uploadResult = await this.api.files.upload(contentFile);

      // 第二步：创建文章记录
      const articleData = {
        title,
        description,
        store_id: uploadResult.id,
        category,
        tags
      };

      return await this.api.articles.create(articleData);
    } catch (error) {
      console.error('发布文章失败:', error);
      throw error;
    }
  }

  // 更新文章内容
  async updateArticle(articleId, title, description, markdownContent, category = '', tags = '') {
    try {
      // 上传新的内容文件
      const contentBlob = new Blob([markdownContent], { type: 'text/markdown' });
      const contentFile = new File([contentBlob], 'article.md', {
        type: 'text/markdown'
      });

      const uploadResult = await this.api.files.upload(contentFile);

      // 更新文章记录
      const articleData = {
        title,
        description,
        store_id: uploadResult.id,
        category,
        tags
      };

      return await this.api.articles.update(articleId, articleData);
    } catch (error) {
      console.error('更新文章失败:', error);
      throw error;
    }
  }
}

// 创建全局API客户端实例
window.apiClient = new ApiClient();
window.articleService = new ArticleService(window.apiClient);

// 导出（如果使用模块系统）
if (typeof module !== 'undefined' && module.exports) {
  module.exports = { ApiClient, ApiError, ArticleService };
}
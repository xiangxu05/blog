// 全局变量
let currentUser = null;
let currentPage = 1;
let pageSize = 10;

// API 基础URL
const API_BASE = '/api';

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', async function() {
    await checkAuthStatus();
    loadHome();
    loadPopularArticles();
    loadCategories();
    
    // 绑定事件
    bindEvents();
});

// 绑定事件
function bindEvents() {
    // 登录表单
    document.getElementById('loginForm').addEventListener('submit', handleLogin);
    
    // 注册表单
    document.getElementById('registerForm').addEventListener('submit', handleRegister);
    
    // 创建文章表单
    document.getElementById('createArticleForm').addEventListener('submit', handleCreateArticle);
    
    // 上传文章表单
    document.getElementById('uploadArticleForm').addEventListener('submit', handleUploadArticle);
    
    // 修改文章表单
    document.getElementById('editArticleForm').addEventListener('submit', handleEditArticle);
    
    // 文件选择事件
    document.getElementById('markdownFile').addEventListener('change', handleFileSelect);
    document.getElementById('editMarkdownFile').addEventListener('change', handleEditFileSelect);
    
    // 搜索框回车事件
    document.getElementById('searchInput').addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            searchArticles();
        }
    });
}

// 检查认证状态
async function checkAuthStatus() {
    const username = localStorage.getItem('username');
    
    if (username) {
        try {
            // 获取完整的用户信息
            const response = await apiRequest('/users/get');
            if (response.code === 200 && response.data) {
                currentUser = response.data;
                updateAuthUI(true);
            } else {
                // 如果获取用户信息失败，清除本地存储
                localStorage.removeItem('username');
                currentUser = null;
                updateAuthUI(false);
            }
        } catch (error) {
            console.error('获取用户信息失败:', error);
            // 网络错误时，暂时使用本地存储的用户名
            currentUser = { username };
            updateAuthUI(true);
        }
    } else {
        currentUser = null;
        updateAuthUI(false);
    }
}

// 更新认证UI
function updateAuthUI(isLoggedIn) {
    const authNav = document.getElementById('authNav');
    const userNav = document.getElementById('userNav');
    const usernameSpan = document.getElementById('username');
    
    if (isLoggedIn) {
        authNav.classList.add('d-none');
        userNav.classList.remove('d-none');
        usernameSpan.textContent = currentUser.username;
    } else {
        authNav.classList.remove('d-none');
        userNav.classList.add('d-none');
    }
}

// API 请求封装
async function apiRequest(url, options = {}) {
    const defaultOptions = {
        headers: {
            'Content-Type': 'application/json',
        },
        credentials: 'include', // 包含cookies
    };
    
    const finalOptions = { ...defaultOptions, ...options };
    
    try {
        const response = await fetch(API_BASE + url, finalOptions);
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.message || '请求失败');
        }
        
        return data;
    } catch (error) {
        console.error('API请求错误:', error);
        showAlert('错误: ' + error.message, 'danger');
        throw error;
    }
}

// 显示提示信息
function showAlert(message, type = 'info') {
    const alertDiv = document.createElement('div');
    alertDiv.className = `alert alert-${type} alert-dismissible fade show position-fixed`;
    alertDiv.style.cssText = 'top: 80px; right: 20px; z-index: 9999; min-width: 300px;';
    alertDiv.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
    `;
    
    document.body.appendChild(alertDiv);
    
    // 3秒后自动消失
    setTimeout(() => {
        if (alertDiv.parentNode) {
            alertDiv.parentNode.removeChild(alertDiv);
        }
    }, 3000);
}

// 显示加载动画
function showLoading() {
    document.getElementById('loading').style.display = 'block';
    document.getElementById('mainContent').style.display = 'none';
}

// 隐藏加载动画
function hideLoading() {
    document.getElementById('loading').style.display = 'none';
    document.getElementById('mainContent').style.display = 'block';
}

// 登录处理
async function handleLogin(e) {
    e.preventDefault();
    
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;
    
    try {
        const response = await apiRequest('/users/login', {
            method: 'POST',
            body: JSON.stringify({ username, password })
        });
        
        if (response.code === 200) {
            localStorage.setItem('username', username);
            
            // 登录成功后获取完整的用户信息
            try {
                const userResponse = await apiRequest('/users/get');
                if (userResponse.code === 200 && userResponse.data) {
                    currentUser = userResponse.data;
                } else {
                    currentUser = { username };
                }
            } catch (error) {
                console.error('获取用户信息失败:', error);
                currentUser = { username };
            }
            
            updateAuthUI(true);
            bootstrap.Modal.getInstance(document.getElementById('loginModal')).hide();
            showAlert('登录成功！', 'success');
            
            // 重置表单
            document.getElementById('loginForm').reset();
        }
    } catch (error) {
        // 错误已在apiRequest中处理
    }
}

// 注册处理
async function handleRegister(e) {
    e.preventDefault();
    
    const username = document.getElementById('registerUsername').value;
    const email = document.getElementById('registerEmail').value;
    const password = document.getElementById('registerPassword').value;
    const confirmPassword = document.getElementById('confirmPassword').value;
    
    if (password !== confirmPassword) {
        showAlert('两次输入的密码不一致', 'danger');
        return;
    }
    
    try {
        const response = await apiRequest('/users/register', {
            method: 'POST',
            body: JSON.stringify({ 
                username, 
                email, 
                password,
                confirm_password: confirmPassword 
            })
        });
        
        if (response.code === 200) {
            bootstrap.Modal.getInstance(document.getElementById('registerModal')).hide();
            showAlert('注册成功！请登录', 'success');
            
            // 重置表单
            document.getElementById('registerForm').reset();
            
            // 显示登录模态框
            setTimeout(() => {
                showLogin();
            }, 1000);
        }
    } catch (error) {
        // 错误已在apiRequest中处理
    }
}

// 退出登录
function logout() {
    if (currentUser) {
        // 调用退出API
        apiRequest('/users/logout', { method: 'GET' }).catch(() => {});
    }
    
    currentUser = null;
    localStorage.removeItem('username');
    updateAuthUI(false);
    showAlert('已退出登录', 'info');
    loadHome();
}

// 显示登录模态框
function showLogin() {
    new bootstrap.Modal(document.getElementById('loginModal')).show();
}

// 显示注册模态框
function showRegister() {
    new bootstrap.Modal(document.getElementById('registerModal')).show();
}

// 显示创建文章模态框
function showCreateArticle() {
    if (!currentUser) {
        showAlert('请先登录', 'warning');
        showLogin();
        return;
    }
    new bootstrap.Modal(document.getElementById('createArticleModal')).show();
}

// 显示上传文章模态框
function showUploadArticle() {
    if (!currentUser) {
        showAlert('请先登录', 'warning');
        showLogin();
        return;
    }
    // 重置表单
    document.getElementById('uploadArticleForm').reset();
    document.getElementById('filePreview').classList.add('d-none');
    document.getElementById('uploadSubmitBtn').disabled = true;
    
    new bootstrap.Modal(document.getElementById('uploadArticleModal')).show();
}

// 文件选择处理
function handleFileSelect(e) {
    const file = e.target.files[0];
    if (!file) {
        document.getElementById('filePreview').classList.add('d-none');
        document.getElementById('uploadSubmitBtn').disabled = true;
        return;
    }
    
    // 检查文件类型
    if (!file.name.toLowerCase().endsWith('.md') && !file.name.toLowerCase().endsWith('.markdown')) {
        showAlert('请选择 .md 或 .markdown 格式的文件', 'warning');
        e.target.value = '';
        return;
    }
    
    // 显示文件信息
    const fileInfo = document.getElementById('fileInfo');
    fileInfo.innerHTML = `
        <strong>文件名：</strong>${file.name}<br>
        <strong>大小：</strong>${(file.size / 1024).toFixed(2)} KB<br>
        <strong>类型：</strong>${file.type || 'text/markdown'}
    `;
    
    // 读取文件内容
    const reader = new FileReader();
    reader.onload = function(e) {
        const content = e.target.result;
        parseMarkdownFile(content);
    };
    reader.readAsText(file, 'UTF-8');
}

// 解析Markdown文件
function parseMarkdownFile(content) {
    try {
        // 解析YAML前置元数据
        const yamlMatch = content.match(/^---\s*\n([\s\S]*?)\n---\s*\n([\s\S]*)$/);
        
        let metadata = {};
        let articleContent = content;
        
        if (yamlMatch) {
            const yamlContent = yamlMatch[1];
            articleContent = yamlMatch[2];
            
            // 简单的YAML解析（支持基本的key: value格式）
            const yamlLines = yamlContent.split('\n');
            yamlLines.forEach(line => {
                const match = line.match(/^([^:]+):\s*(.*)$/);
                if (match) {
                    const key = match[1].trim();
                    let value = match[2].trim();
                    
                    // 处理数组格式 [item1, item2]
                    if (value.startsWith('[') && value.endsWith(']')) {
                        value = value.slice(1, -1).split(',').map(item => item.trim());
                    }
                    
                    metadata[key] = value;
                }
            });
        }
        
        // 填充表单字段
        document.getElementById('parsedTitle').value = metadata.title || '';
        document.getElementById('parsedDescription').value = metadata.description || '';
        
        // 处理分类
        let category = '';
        if (metadata.categories) {
            if (Array.isArray(metadata.categories)) {
                category = metadata.categories[0] || '';
            } else {
                category = metadata.categories;
            }
        }
        document.getElementById('parsedCategory').value = category;
        
        // 处理关键字/标签
        let keywords = '';
        if (metadata.keywords) {
            if (Array.isArray(metadata.keywords)) {
                keywords = metadata.keywords.join(', ');
            } else {
                keywords = metadata.keywords.replace(/，/g, ','); // 替换中文逗号
            }
        }
        document.getElementById('parsedKeywords').value = keywords;
        
        // 显示内容预览（限制长度）
        const contentPreview = articleContent.length > 500 
            ? articleContent.substring(0, 500) + '...'
            : articleContent;
        document.getElementById('parsedContent').value = contentPreview;
        
        // 存储完整内容到隐藏属性
        document.getElementById('parsedContent').dataset.fullContent = articleContent;
        
        // 显示预览区域
        document.getElementById('filePreview').classList.remove('d-none');
        document.getElementById('uploadSubmitBtn').disabled = false;
        
    } catch (error) {
        console.error('解析Markdown文件失败:', error);
        showAlert('解析文件失败，请检查文件格式', 'danger');
    }
}

// 上传文章处理
async function handleUploadArticle(e) {
    e.preventDefault();
    
    if (!currentUser) {
        showAlert('请先登录', 'warning');
        return;
    }
    
    const title = document.getElementById('parsedTitle').value;
    const description = document.getElementById('parsedDescription').value;
    const category = document.getElementById('parsedCategory').value;
    const tagsInput = document.getElementById('parsedKeywords').value;
    const content = document.getElementById('parsedContent').dataset.fullContent;
    
    if (!title || !content) {
        showAlert('标题和内容不能为空', 'warning');
        return;
    }
    
    // 智能处理图片路径
    const imageCheck = await processMarkdownImages(content);
    if (imageCheck.invalidImages > 0) {
        const proceed = confirm(`检测到 ${imageCheck.invalidImages} 个无效图片路径，是否继续发布？`);
        if (!proceed) {
            return;
        }
    }
    
    const tags = tagsInput ? tagsInput.split(',').map(tag => tag.trim()).filter(tag => tag) : [];
    
    try {
        const response = await apiRequest('/articles', {
            method: 'POST',
            body: JSON.stringify({
                title,
                description,
                category: category || '未分类',
                tags,
                content
            })
        });
        
        if (response.code === 200) {
            bootstrap.Modal.getInstance(document.getElementById('uploadArticleModal')).hide();
            showAlert('文章上传成功！', 'success');
            
            // 重置表单
            document.getElementById('uploadArticleForm').reset();
            document.getElementById('filePreview').classList.add('d-none');
            
            // 刷新文章列表
            loadArticles();
        }
    } catch (error) {
        // 错误已在apiRequest中处理
    }
}

// 创建文章处理
async function handleCreateArticle(e) {
    e.preventDefault();
    
    if (!currentUser) {
        showAlert('请先登录', 'warning');
        return;
    }
    
    const title = document.getElementById('articleTitle').value;
    const description = document.getElementById('articleDescription').value;
    const category = document.getElementById('articleCategory').value;
    const tagsInput = document.getElementById('articleTags').value;
    const content = document.getElementById('articleContent').value;
    
    // 智能处理图片路径
    const imageCheck = await processMarkdownImages(content);
    if (imageCheck.invalidImages > 0) {
        const proceed = confirm(`检测到 ${imageCheck.invalidImages} 个无效图片路径，是否继续发布？`);
        if (!proceed) {
            return;
        }
    }
    
    const tags = tagsInput ? tagsInput.split(',').map(tag => tag.trim()).filter(tag => tag) : [];
    
    try {
        const response = await apiRequest('/articles', {
            method: 'POST',
            body: JSON.stringify({
                title,
                description,
                category,
                tags,
                content
            })
        });
        
        if (response.code === 200) {
            bootstrap.Modal.getInstance(document.getElementById('createArticleModal')).hide();
            showAlert('文章创建成功！', 'success');
            
            // 重置表单
            document.getElementById('createArticleForm').reset();
            
            // 刷新文章列表
            loadArticles();
        }
    } catch (error) {
        // 错误已在apiRequest中处理
    }
}

// 加载首页
function loadHome() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
        <div class="text-center py-5">
            <h2 class="mb-4">欢迎来到我的博客</h2>
            <p class="lead mb-4">这里记录着我的技术学习历程、生活感悟和思考</p>
            <div class="row">
                <div class="col-md-4 mb-4">
                    <div class="card h-100 text-center">
                        <div class="card-body">
                            <i class="fas fa-code fa-3x text-primary mb-3"></i>
                            <h5 class="card-title">技术分享</h5>
                            <p class="card-text">分享编程技术、开发经验和最佳实践</p>
                        </div>
                    </div>
                </div>
                <div class="col-md-4 mb-4">
                    <div class="card h-100 text-center">
                        <div class="card-body">
                            <i class="fas fa-lightbulb fa-3x text-warning mb-3"></i>
                            <h5 class="card-title">思考感悟</h5>
                            <p class="card-text">记录生活中的思考和感悟</p>
                        </div>
                    </div>
                </div>
                <div class="col-md-4 mb-4">
                    <div class="card h-100 text-center">
                        <div class="card-body">
                            <i class="fas fa-heart fa-3x text-danger mb-3"></i>
                            <h5 class="card-title">生活记录</h5>
                            <p class="card-text">分享生活中的美好时刻</p>
                        </div>
                    </div>
                </div>
            </div>
            <button class="btn btn-primary btn-lg mt-4" onclick="loadArticles()">
                <i class="fas fa-arrow-right me-2"></i>查看所有文章
            </button>
        </div>
    `;
}

// 加载文章列表
async function loadArticles(page = 1) {
    showLoading();
    
    try {
        const response = await apiRequest(`/articles?limit=${pageSize}&offset=${(page - 1) * pageSize}`);
        
        if (response.code === 200) {
            renderArticleList(response.data.articles, response.data.total, page);
        }
    } catch (error) {
        document.getElementById('contentArea').innerHTML = '<div class="alert alert-danger">加载文章失败</div>';
    } finally {
        hideLoading();
    }
}

// 渲染文章列表
function renderArticleList(articles, total, page) {
    const contentArea = document.getElementById('contentArea');
    
    if (!articles || articles.length === 0) {
        contentArea.innerHTML = '<div class="alert alert-info">暂无文章</div>';
        return;
    }
    
    let html = '<h2 class="mb-4">文章列表</h2>';
    
    articles.forEach(article => {
        const createdAt = new Date(article.created_at).toLocaleDateString('zh-CN');
        const tags = article.tags ? article.tags.map(tag => 
            `<span class="tag-badge">${tag.name}</span>`
        ).join('') : '';
        
        html += `
            <div class="article-card card mb-4">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-start mb-2">
                        <h5 class="card-title mb-0">
                            <a href="#" onclick="loadArticle(${article.id})" class="text-decoration-none">
                                ${article.title}
                            </a>
                        </h5>
                        <span class="category-badge">${article.category ? article.category.name : '未分类'}</span>
                    </div>
                    <p class="text-muted small mb-2">
                        <i class="fas fa-user me-1"></i>${article.user ? article.user.username : '匿名'}
                        <i class="fas fa-calendar ms-3 me-1"></i>${createdAt}
                        <i class="fas fa-eye ms-3 me-1"></i>${article.views || 0} 次浏览
                    </p>
                    <p class="card-text">${article.content && article.content.description ? article.content.description : '暂无描述'}</p>
                    <div class="d-flex justify-content-between align-items-center">
                        <div>${tags}</div>
                        <a href="#" onclick="loadArticle(${article.id})" class="btn btn-outline-primary btn-sm">
                            阅读更多 <i class="fas fa-arrow-right ms-1"></i>
                        </a>
                    </div>
                </div>
            </div>
        `;
    });
    
    // 添加分页
    const totalPages = Math.ceil(total / pageSize);
    if (totalPages > 1) {
        html += renderPagination(page, totalPages, 'loadArticles');
    }
    
    contentArea.innerHTML = html;
}

// 渲染分页
function renderPagination(currentPage, totalPages, functionName) {
    let html = '<nav aria-label="文章分页"><ul class="pagination justify-content-center">';
    
    // 上一页
    if (currentPage > 1) {
        html += `<li class="page-item"><a class="page-link" href="#" onclick="${functionName}(${currentPage - 1})">上一页</a></li>`;
    }
    
    // 页码
    for (let i = Math.max(1, currentPage - 2); i <= Math.min(totalPages, currentPage + 2); i++) {
        const active = i === currentPage ? 'active' : '';
        html += `<li class="page-item ${active}"><a class="page-link" href="#" onclick="${functionName}(${i})">${i}</a></li>`;
    }
    
    // 下一页
    if (currentPage < totalPages) {
        html += `<li class="page-item"><a class="page-link" href="#" onclick="${functionName}(${currentPage + 1})">下一页</a></li>`;
    }
    
    html += '</ul></nav>';
    return html;
}

// 加载单篇文章
async function loadArticle(id) {
    showLoading();
    
    try {
        const response = await apiRequest(`/articles/${id}`);
        
        if (response.code === 200) {
            renderArticle(response.data);
        }
    } catch (error) {
        document.getElementById('contentArea').innerHTML = '<div class="alert alert-danger">加载文章失败</div>';
    } finally {
        hideLoading();
    }
}

// 渲染单篇文章
function renderArticle(article) {
    const contentArea = document.getElementById('contentArea');
    const createdAt = new Date(article.created_at).toLocaleDateString('zh-CN');
    const tags = article.tags ? article.tags.map(tag => 
        `<span class="tag-badge">${tag.name}</span>`
    ).join('') : '';
    
    const html = `
        <div class="article-detail">
            <nav aria-label="breadcrumb">
                <ol class="breadcrumb">
                    <li class="breadcrumb-item"><a href="#" onclick="loadHome()">首页</a></li>
                    <li class="breadcrumb-item"><a href="#" onclick="loadArticles()">文章</a></li>
                    <li class="breadcrumb-item active">${article.title}</li>
                </ol>
            </nav>
            
            <div class="card">
                <div class="card-body">
                    <h1 class="card-title mb-3">${article.title}</h1>
                    
                    <div class="article-meta mb-4 pb-3 border-bottom">
                        <div class="row">
                            <div class="col-md-6">
                                <p class="text-muted mb-1">
                                    <i class="fas fa-user me-2"></i>作者：${article.user ? article.user.username : '匿名'}
                                </p>
                                <p class="text-muted mb-1">
                                    <i class="fas fa-calendar me-2"></i>发布时间：${createdAt}
                                </p>
                            </div>
                            <div class="col-md-6">
                                <p class="text-muted mb-1">
                                    <i class="fas fa-folder me-2"></i>分类：
                                    <span class="category-badge">${article.category ? article.category.name : '未分类'}</span>
                                </p>
                                <p class="text-muted mb-1">
                                    <i class="fas fa-eye me-2"></i>浏览量：${article.views || 0}
                                </p>
                            </div>
                        </div>
                        ${tags ? `<div class="mt-2"><i class="fas fa-tags me-2"></i>${tags}</div>` : ''}
                    </div>
                    
                    <div class="article-content">
                        ${formatContent(article.content ? article.content.content : '内容加载失败')}
                    </div>
                </div>
            </div>
            
            <div class="text-center mt-4">
                <button class="btn btn-outline-primary" onclick="loadArticles()">
                    <i class="fas fa-arrow-left me-2"></i>返回文章列表
                </button>
            </div>
        </div>
    `;
    
    contentArea.innerHTML = html;
}

// 格式化文章内容
function formatContent(content) {
    if (!content) return '';
    
    try {
        // 配置marked选项
        marked.setOptions({
            highlight: function(code, lang) {
                if (lang && hljs.getLanguage(lang)) {
                    try {
                        return hljs.highlight(code, { language: lang }).value;
                    } catch (err) {
                        console.warn('代码高亮失败:', err);
                    }
                }
                return hljs.highlightAuto(code).value;
            },
            breaks: true,
            gfm: true
        });
        
        // 使用marked渲染Markdown
        return marked.parse(content);
    } catch (error) {
        console.error('Markdown渲染失败:', error);
        // 降级处理：简单的文本格式化
        return content.split('\n').map(paragraph => 
            paragraph.trim() ? `<p>${escapeHtml(paragraph)}</p>` : ''
        ).join('');
    }
}

// HTML转义函数
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// 搜索文章
async function searchArticles() {
    const keyword = document.getElementById('searchInput').value.trim();
    
    if (!keyword) {
        showAlert('请输入搜索关键词', 'warning');
        return;
    }
    
    showLoading();
    
    try {
        const response = await apiRequest(`/articles/search?keyword=${encodeURIComponent(keyword)}&limit=${pageSize}&offset=0`);
        
        if (response.code === 200) {
            renderSearchResults(response.data.articles, keyword, response.data.total);
        }
    } catch (error) {
        document.getElementById('contentArea').innerHTML = '<div class="alert alert-danger">搜索失败</div>';
    } finally {
        hideLoading();
    }
}

// 渲染搜索结果
function renderSearchResults(articles, keyword, total) {
    const contentArea = document.getElementById('contentArea');
    
    let html = `<h2 class="mb-4">搜索结果 "${keyword}" (${total} 篇文章)</h2>`;
    
    if (!articles || articles.length === 0) {
        html += '<div class="alert alert-info">没有找到相关文章</div>';
    } else {
        articles.forEach(article => {
            const createdAt = new Date(article.created_at).toLocaleDateString('zh-CN');
            const tags = article.tags ? article.tags.map(tag => 
                `<span class="tag-badge">${tag.name}</span>`
            ).join('') : '';
            
            html += `
                <div class="article-card card mb-4">
                    <div class="card-body">
                        <div class="d-flex justify-content-between align-items-start mb-2">
                            <h5 class="card-title mb-0">
                                <a href="#" onclick="loadArticle(${article.id})" class="text-decoration-none">
                                    ${article.title}
                                </a>
                            </h5>
                            <span class="category-badge">${article.category ? article.category.name : '未分类'}</span>
                        </div>
                        <p class="text-muted small mb-2">
                            <i class="fas fa-user me-1"></i>${article.user ? article.user.username : '匿名'}
                            <i class="fas fa-calendar ms-3 me-1"></i>${createdAt}
                            <i class="fas fa-eye ms-3 me-1"></i>${article.views || 0} 次浏览
                        </p>
                        <p class="card-text">${article.content && article.content.description ? article.content.description : '暂无描述'}</p>
                        <div class="d-flex justify-content-between align-items-center">
                            <div>${tags}</div>
                            <a href="#" onclick="loadArticle(${article.id})" class="btn btn-outline-primary btn-sm">
                                阅读更多 <i class="fas fa-arrow-right ms-1"></i>
                            </a>
                        </div>
                    </div>
                </div>
            `;
        });
    }
    
    contentArea.innerHTML = html;
}

// 加载热门文章
async function loadPopular() {
    showLoading();
    
    try {
        const response = await apiRequest('/articles/popular?limit=10');
        
        if (response.code === 200) {
            renderPopularArticles(response.data);
        }
    } catch (error) {
        document.getElementById('contentArea').innerHTML = '<div class="alert alert-danger">加载热门文章失败</div>';
    } finally {
        hideLoading();
    }
}

// 渲染热门文章页面
function renderPopularArticles(articles) {
    const contentArea = document.getElementById('contentArea');
    
    let html = '<h2 class="mb-4"><i class="fas fa-fire me-2"></i>热门文章</h2>';
    
    if (!articles || articles.length === 0) {
        html += '<div class="alert alert-info">暂无热门文章</div>';
    } else {
        articles.forEach((article, index) => {
            const createdAt = new Date(article.created_at).toLocaleDateString('zh-CN');
            const tags = article.tags ? article.tags.map(tag => 
                `<span class="tag-badge">${tag.name}</span>`
            ).join('') : '';
            
            html += `
                <div class="article-card card mb-4">
                    <div class="card-body">
                        <div class="d-flex align-items-start">
                            <div class="me-3">
                                <span class="badge bg-primary fs-6">${index + 1}</span>
                            </div>
                            <div class="flex-grow-1">
                                <div class="d-flex justify-content-between align-items-start mb-2">
                                    <h5 class="card-title mb-0">
                                        <a href="#" onclick="loadArticle(${article.id})" class="text-decoration-none">
                                            ${article.title}
                                        </a>
                                    </h5>
                                    <span class="category-badge">${article.category ? article.category.name : '未分类'}</span>
                                </div>
                                <p class="text-muted small mb-2">
                                    <i class="fas fa-user me-1"></i>${article.user ? article.user.username : '匿名'}
                                    <i class="fas fa-calendar ms-3 me-1"></i>${createdAt}
                                    <i class="fas fa-eye ms-3 me-1"></i>${article.views || 0} 次浏览
                                </p>
                                <p class="card-text">${article.content && article.content.description ? article.content.description : '暂无描述'}</p>
                                <div class="d-flex justify-content-between align-items-center">
                                    <div>${tags}</div>
                                    <a href="#" onclick="loadArticle(${article.id})" class="btn btn-outline-primary btn-sm">
                                        阅读更多 <i class="fas fa-arrow-right ms-1"></i>
                                    </a>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            `;
        });
    }
    
    contentArea.innerHTML = html;
}

// 加载侧边栏热门文章
async function loadPopularArticles() {
    try {
        const response = await apiRequest('/articles/popular?limit=5');
        
        if (response.code === 200) {
            renderSidebarPopular(response.data);
        }
    } catch (error) {
        console.error('加载侧边栏热门文章失败:', error);
    }
}

// 渲染侧边栏热门文章
function renderSidebarPopular(articles) {
    const container = document.getElementById('popularArticles');
    
    if (!articles || articles.length === 0) {
        container.innerHTML = '<p class="text-muted">暂无热门文章</p>';
        return;
    }
    
    let html = '';
    articles.forEach(article => {
        html += `
            <div class="d-flex mb-3">
                <div class="flex-grow-1">
                    <h6 class="mb-1">
                        <a href="#" onclick="loadArticle(${article.id})" class="text-decoration-none text-dark">
                            ${article.title}
                        </a>
                    </h6>
                    <small class="text-muted">
                        <i class="fas fa-eye me-1"></i>${article.views || 0} 次浏览
                    </small>
                </div>
            </div>
        `;
    });
    
    container.innerHTML = html;
}

// 加载分类
async function loadCategories() {
    try {
        const response = await apiRequest('/articles/categories');
        
        if (response.code === 200) {
            renderSidebarCategories(response.data);
        }
    } catch (error) {
        console.error('加载分类失败:', error);
    }
}

// 渲染侧边栏分类
function renderSidebarCategories(categories) {
    const container = document.getElementById('categoriesList');
    
    if (!categories || categories.length === 0) {
        container.innerHTML = '<p class="text-muted">暂无分类</p>';
        return;
    }
    
    let html = '';
    categories.forEach(category => {
        html += `
            <a href="#" onclick="loadCategoryArticles(${category.id}, '${category.name}')" 
               class="btn btn-outline-secondary btn-sm me-2 mb-2">
                ${category.name}
            </a>
        `;
    });
    
    container.innerHTML = html;
}

// 加载分类文章
async function loadCategoryArticles(categoryId, categoryName, page = 1) {
    showLoading();
    
    try {
        const response = await apiRequest(`/articles/category/${categoryId}?limit=${pageSize}&offset=${(page - 1) * pageSize}`);
        
        if (response.code === 200) {
            renderCategoryArticles(response.data.articles, categoryName, response.data.total, page, categoryId);
        }
    } catch (error) {
        document.getElementById('contentArea').innerHTML = '<div class="alert alert-danger">加载分类文章失败</div>';
    } finally {
        hideLoading();
    }
}

// 渲染分类文章
function renderCategoryArticles(articles, categoryName, total, page, categoryId) {
    const contentArea = document.getElementById('contentArea');
    
    let html = `<h2 class="mb-4">分类：${categoryName} (${total} 篇文章)</h2>`;
    
    if (!articles || articles.length === 0) {
        html += '<div class="alert alert-info">该分类下暂无文章</div>';
    } else {
        articles.forEach(article => {
            const createdAt = new Date(article.created_at).toLocaleDateString('zh-CN');
            const tags = article.tags ? article.tags.map(tag => 
                `<span class="tag-badge">${tag.name}</span>`
            ).join('') : '';
            
            html += `
                <div class="article-card card mb-4">
                    <div class="card-body">
                        <div class="d-flex justify-content-between align-items-start mb-2">
                            <h5 class="card-title mb-0">
                                <a href="#" onclick="loadArticle(${article.id})" class="text-decoration-none">
                                    ${article.title}
                                </a>
                            </h5>
                            <span class="category-badge">${article.category ? article.category.name : '未分类'}</span>
                        </div>
                        <p class="text-muted small mb-2">
                            <i class="fas fa-user me-1"></i>${article.user ? article.user.username : '匿名'}
                            <i class="fas fa-calendar ms-3 me-1"></i>${createdAt}
                            <i class="fas fa-eye ms-3 me-1"></i>${article.views || 0} 次浏览
                        </p>
                        <p class="card-text">${article.content && article.content.description ? article.content.description : '暂无描述'}</p>
                        <div class="d-flex justify-content-between align-items-center">
                            <div>${tags}</div>
                            <a href="#" onclick="loadArticle(${article.id})" class="btn btn-outline-primary btn-sm">
                                阅读更多 <i class="fas fa-arrow-right ms-1"></i>
                            </a>
                        </div>
                    </div>
                </div>
            `;
        });
        
        // 添加分页
        const totalPages = Math.ceil(total / pageSize);
        if (totalPages > 1) {
            html += renderCategoryPagination(page, totalPages, categoryId, categoryName);
        }
    }
    
    contentArea.innerHTML = html;
}

// 渲染分类分页
function renderCategoryPagination(currentPage, totalPages, categoryId, categoryName) {
    let html = '<nav aria-label="分类文章分页"><ul class="pagination justify-content-center">';
    
    // 上一页
    if (currentPage > 1) {
        html += `<li class="page-item"><a class="page-link" href="#" onclick="loadCategoryArticles(${categoryId}, '${categoryName}', ${currentPage - 1})">上一页</a></li>`;
    }
    
    // 页码
    for (let i = Math.max(1, currentPage - 2); i <= Math.min(totalPages, currentPage + 2); i++) {
        const active = i === currentPage ? 'active' : '';
        html += `<li class="page-item ${active}"><a class="page-link" href="#" onclick="loadCategoryArticles(${categoryId}, '${categoryName}', ${i})">${i}</a></li>`;
    }
    
    // 下一页
    if (currentPage < totalPages) {
        html += `<li class="page-item"><a class="page-link" href="#" onclick="loadCategoryArticles(${categoryId}, '${categoryName}', ${currentPage + 1})">下一页</a></li>`;
    }
    
    html += '</ul></nav>';
    return html;
}

// 加载我的文章
async function loadMyArticles() {
    if (!currentUser || !currentUser.id) {
        showAlert('请先登录', 'warning');
        showLogin();
        return;
    }
    
    showLoading();
    
    try {
        const articlesResponse = await apiRequest(`/articles/user/${currentUser.id}?limit=${pageSize}&offset=0`);
        
        if (articlesResponse.code === 200) {
            renderMyArticles(articlesResponse.data.articles, articlesResponse.data.total);
        } else {
            document.getElementById('contentArea').innerHTML = '<div class="alert alert-warning">暂无文章数据</div>';
        }
    } catch (error) {
        console.error('加载我的文章失败:', error);
        document.getElementById('contentArea').innerHTML = '<div class="alert alert-danger">加载我的文章失败，请稍后重试</div>';
    } finally {
        hideLoading();
    }
}

// 渲染我的文章
function renderMyArticles(articles, total) {
    const contentArea = document.getElementById('contentArea');
    
    let html = `<h2 class="mb-4">我的文章 (${total} 篇)</h2>`;
    
    if (!articles || articles.length === 0) {
        html += `
            <div class="alert alert-info">
                您还没有发布任何文章
                <div class="mt-3">
                    <button class="btn btn-primary me-2" onclick="showCreateArticle()">
                        <i class="fas fa-plus me-1"></i>写文章
                    </button>
                    <button class="btn btn-outline-primary" onclick="showUploadArticle()">
                        <i class="fas fa-upload me-1"></i>上传文章
                    </button>
                </div>
            </div>
        `;
    } else {
        articles.forEach(article => {
            const createdAt = new Date(article.created_at).toLocaleDateString('zh-CN');
            const tags = article.tags ? article.tags.map(tag => 
                `<span class="tag-badge">${tag.name}</span>`
            ).join('') : '';
            
            html += `
                <div class="article-card card mb-4">
                    <div class="card-body">
                        <div class="d-flex justify-content-between align-items-start mb-2">
                            <h5 class="card-title mb-0">
                                <a href="#" onclick="loadArticle(${article.id})" class="text-decoration-none">
                                    ${article.title}
                                </a>
                            </h5>
                            <span class="category-badge">${article.category ? article.category.name : '未分类'}</span>
                        </div>
                        <p class="text-muted small mb-2">
                            <i class="fas fa-calendar me-1"></i>${createdAt}
                            <i class="fas fa-eye ms-3 me-1"></i>${article.views || 0} 次浏览
                        </p>
                        <p class="card-text">${article.content && article.content.description ? article.content.description : '暂无描述'}</p>
                        <div class="d-flex justify-content-between align-items-center">
                            <div>${tags}</div>
                            <div>
                                <a href="#" onclick="loadArticle(${article.id})" class="btn btn-outline-primary btn-sm me-2">
                                    <i class="fas fa-eye me-1"></i>查看
                                </a>
                                <button class="btn btn-outline-secondary btn-sm me-2" onclick="showEditArticle(${article.id})">
                                    <i class="fas fa-edit me-1"></i>修改
                                </button>
                                <button class="btn btn-outline-danger btn-sm" onclick="deleteArticle(${article.id})">
                                    <i class="fas fa-trash me-1"></i>删除
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            `;
        });
    }
    
    contentArea.innerHTML = html;
}

// 删除文章
async function deleteArticle(id) {
    if (!confirm('确定要删除这篇文章吗？此操作不可恢复。')) {
        return;
    }
    
    try {
        const response = await apiRequest(`/articles/${id}`, {
            method: 'DELETE'
        });
        
        if (response.code === 200) {
            showAlert('文章删除成功', 'success');
            loadMyArticles(); // 重新加载我的文章列表
        }
    } catch (error) {
        // 错误已在apiRequest中处理
    }
}

// 显示修改文章模态框
async function showEditArticle(id) {
    if (!currentUser) {
        showAlert('请先登录', 'warning');
        return;
    }
    
    try {
        // 获取文章详情
        const response = await apiRequest(`/articles/${id}`);
        if (response.code === 200) {
            const article = response.data;
            
            // 填充表单
            document.getElementById('editArticleId').value = article.id;
            document.getElementById('editParsedTitle').value = article.title;
            document.getElementById('editParsedCategory').value = article.category ? article.category.name : '';
            document.getElementById('editParsedDescription').value = article.content ? article.content.description : '';
            
            // 处理标签
            const tags = article.tags ? article.tags.map(tag => tag.name).join(', ') : '';
            document.getElementById('editParsedKeywords').value = tags;
            
            // 填充内容
            document.getElementById('editParsedContent').value = article.content ? article.content.content : '';
            
            // 显示模态框
            new bootstrap.Modal(document.getElementById('editArticleModal')).show();
        }
    } catch (error) {
        showAlert('获取文章信息失败', 'danger');
    }
}

// 处理修改文章的文件选择
function handleEditFileSelect(e) {
    const file = e.target.files[0];
    if (!file) return;
    
    // 检查文件类型
    if (!file.name.match(/\.(md|markdown)$/i)) {
        showAlert('请选择 .md 或 .markdown 格式的文件', 'warning');
        e.target.value = '';
        return;
    }
    
    const reader = new FileReader();
    reader.onload = function(event) {
        const content = event.target.result;
        parseMarkdownFileForEdit(content);
    };
    reader.readAsText(file);
}

// 为修改文章解析Markdown文件
function parseMarkdownFileForEdit(content) {
    try {
        // 解析YAML前置元数据
        const yamlMatch = content.match(/^---\s*\n([\s\S]*?)\n---\s*\n([\s\S]*)$/);
        
        let metadata = {};
        let articleContent = content;
        
        if (yamlMatch) {
            const yamlContent = yamlMatch[1];
            articleContent = yamlMatch[2];
            
            // 简单的YAML解析
            const yamlLines = yamlContent.split('\n');
            yamlLines.forEach(line => {
                const match = line.match(/^([^:]+):\s*(.*)$/);
                if (match) {
                    const key = match[1].trim();
                    let value = match[2].trim();
                    
                    // 处理数组格式
                    if (value.startsWith('[') && value.endsWith(']')) {
                        value = value.slice(1, -1).split(',').map(item => item.trim());
                    }
                    
                    metadata[key] = value;
                }
            });
        }
        
        // 更新表单字段
        if (metadata.title) {
            document.getElementById('editParsedTitle').value = metadata.title;
        }
        if (metadata.description) {
            document.getElementById('editParsedDescription').value = metadata.description;
        }
        
        // 处理分类
        if (metadata.categories) {
            let category = '';
            if (Array.isArray(metadata.categories)) {
                category = metadata.categories[0] || '';
            } else {
                category = metadata.categories;
            }
            document.getElementById('editParsedCategory').value = category;
        }
        
        // 处理关键字/标签
        if (metadata.keywords) {
            let keywords = '';
            if (Array.isArray(metadata.keywords)) {
                keywords = metadata.keywords.join(', ');
            } else {
                keywords = metadata.keywords.replace(/，/g, ',');
            }
            document.getElementById('editParsedKeywords').value = keywords;
        }
        
        // 更新内容
        document.getElementById('editParsedContent').value = articleContent;
        
    } catch (error) {
        console.error('解析Markdown文件失败:', error);
        showAlert('解析文件失败，请检查文件格式', 'danger');
    }
}

// 处理修改文章提交
async function handleEditArticle(e) {
    e.preventDefault();
    
    if (!currentUser) {
        showAlert('请先登录', 'warning');
        return;
    }
    
    const articleId = document.getElementById('editArticleId').value;
    const title = document.getElementById('editParsedTitle').value;
    const description = document.getElementById('editParsedDescription').value;
    const category = document.getElementById('editParsedCategory').value;
    const tagsInput = document.getElementById('editParsedKeywords').value;
    const content = document.getElementById('editParsedContent').value;
    
    if (!title || !content) {
        showAlert('标题和内容不能为空', 'warning');
        return;
    }
    
    const tags = tagsInput ? tagsInput.split(',').map(tag => tag.trim()).filter(tag => tag) : [];
    
    try {
        const response = await apiRequest(`/articles/${articleId}`, {
            method: 'PUT',
            body: JSON.stringify({
                title,
                description,
                category: category || '未分类',
                tags,
                content
            })
        });
        
        if (response.code === 200) {
            bootstrap.Modal.getInstance(document.getElementById('editArticleModal')).hide();
            showAlert('文章修改成功！', 'success');
            
            // 重置表单
            document.getElementById('editArticleForm').reset();
            
            // 刷新我的文章列表
            loadMyArticles();
        }
    } catch (error) {
        // 错误已在apiRequest中处理
    }
}

// 工具函数：格式化日期
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
}

// 工具函数：截取文本
function truncateText(text, maxLength = 100) {
    if (text.length <= maxLength) {
        return text;
    }
    return text.substring(0, maxLength) + '...';
}

// 图片上传相关函数
let currentUploadedImageUrl = null;

// 显示图片上传区域
function showImageUpload(type) {
    const uploadArea = document.getElementById(`${type}ImageUploadArea`);
    uploadArea.style.display = 'block';
    
    // 重置上传状态
    const fileInput = document.getElementById(`${type}ImageFile`);
    const preview = document.getElementById(`${type}ImagePreview`);
    const insertBtn = document.getElementById(`${type}InsertImageBtn`);
    
    fileInput.value = '';
    preview.style.display = 'none';
    insertBtn.disabled = true;
    currentUploadedImageUrl = null;
}

// 隐藏图片上传区域
function hideImageUpload(type) {
    const uploadArea = document.getElementById(`${type}ImageUploadArea`);
    uploadArea.style.display = 'none';
}

// 处理图片上传
async function handleImageUpload(input, type) {
    const file = input.files[0];
    if (!file) return;
    
    // 检查文件类型
    const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
    if (!allowedTypes.includes(file.type)) {
        showAlert('请选择有效的图片文件（JPG、PNG、GIF、WebP）', 'warning');
        input.value = '';
        return;
    }
    
    // 检查文件大小（5MB）
    if (file.size > 5 * 1024 * 1024) {
        showAlert('图片大小不能超过5MB', 'warning');
        input.value = '';
        return;
    }
    
    // 显示预览
    const preview = document.getElementById(`${type}ImagePreview`);
    const previewImg = document.getElementById(`${type}PreviewImg`);
    const reader = new FileReader();
    
    reader.onload = function(e) {
        previewImg.src = e.target.result;
        preview.style.display = 'block';
    };
    reader.readAsDataURL(file);
    
    // 上传图片
    try {
        showLoading();
        
        const formData = new FormData();
        formData.append('image', file);
        
        const response = await fetch('/api/articles/upload-image', {
            method: 'POST',
            body: formData,
            credentials: 'include'
        });
        
        const result = await response.json();
        
        if (result.code === 200) {
            currentUploadedImageUrl = result.data.url;
            const insertBtn = document.getElementById(`${type}InsertImageBtn`);
            insertBtn.disabled = false;
            showAlert('图片上传成功！', 'success');
        } else {
            throw new Error(result.msg || '上传失败');
        }
    } catch (error) {
        showAlert('图片上传失败：' + error.message, 'danger');
        input.value = '';
        preview.style.display = 'none';
    } finally {
        hideLoading();
    }
}

// 插入图片到内容
function insertImageToContent(type) {
    if (!currentUploadedImageUrl) {
        showAlert('请先上传图片', 'warning');
        return;
    }
    
    const contentTextarea = type === 'create' ? 
        document.getElementById('articleContent') : 
        document.getElementById('editParsedContent');
    
    const imageMarkdown = `\n![图片描述](${currentUploadedImageUrl})\n`;
    
    // 在光标位置插入图片
    const cursorPos = contentTextarea.selectionStart;
    const textBefore = contentTextarea.value.substring(0, cursorPos);
    const textAfter = contentTextarea.value.substring(contentTextarea.selectionEnd);
    
    contentTextarea.value = textBefore + imageMarkdown + textAfter;
    
    // 设置光标位置到插入内容之后
    const newCursorPos = cursorPos + imageMarkdown.length;
    contentTextarea.setSelectionRange(newCursorPos, newCursorPos);
    contentTextarea.focus();
    
    // 隐藏上传区域
    hideImageUpload(type);
    
    showAlert('图片已插入到内容中', 'success');
}

// 智能路径处理：扫描Markdown内容中的图片引用
function scanMarkdownImages(content) {
    const imageRegex = /!\[([^\]]*)\]\(([^)]+)\)/g;
    const images = [];
    let match;
    
    while ((match = imageRegex.exec(content)) !== null) {
        images.push({
            alt: match[1],
            src: match[2],
            fullMatch: match[0]
        });
    }
    
    return images;
}

// 检查图片路径是否有效
async function validateImagePath(imagePath) {
    try {
        const response = await fetch(imagePath, { method: 'HEAD' });
        return response.ok;
    } catch {
        return false;
    }
}

// 智能处理Markdown中的图片路径
async function processMarkdownImages(content) {
    const images = scanMarkdownImages(content);
    const invalidImages = [];
    
    for (const image of images) {
        // 跳过已经是完整URL的图片
        if (image.src.startsWith('http://') || image.src.startsWith('https://') || image.src.startsWith('/static/')) {
            continue;
        }
        
        // 检查相对路径图片是否存在
        const isValid = await validateImagePath(image.src);
        if (!isValid) {
            invalidImages.push(image);
        }
    }
    
    if (invalidImages.length > 0) {
        const imageList = invalidImages.map(img => `- ${img.alt || '未命名'}: ${img.src}`).join('\n');
        showAlert(`检测到以下图片路径可能无效，请检查或重新上传：\n${imageList}`, 'warning');
    }
    
    return {
        totalImages: images.length,
        invalidImages: invalidImages.length,
        validImages: images.length - invalidImages.length
    };
}
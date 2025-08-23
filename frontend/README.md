# 个人博客系统 - 前端

一个基于 Vue 3 + TypeScript + Vite 构建的现代化个人博客前端应用。

## ✨ 特性

- 🚀 **现代化技术栈**: Vue 3 + TypeScript + Vite
- 🎨 **优雅设计**: 苹果风格的响应式布局
- 📝 **Markdown 编辑器**: 集成 Toast UI Editor，支持代码高亮
- 🔐 **用户认证**: 完整的注册/登录/权限管理系统
- 📱 **响应式设计**: 完美适配桌面端和移动端
- 🌙 **主题切换**: 支持明暗主题切换
- 🔍 **搜索功能**: 文章搜索和分类筛选
- 📊 **管理后台**: 功能完整的文章管理系统
- ⚡ **性能优化**: 路由懒加载、代码分割

## 🛠️ 技术栈

- **框架**: Vue 3 (Composition API)
- **语言**: TypeScript
- **构建工具**: Vite
- **路由**: Vue Router 4
- **状态管理**: Pinia
- **UI 组件库**: Element Plus
- **HTTP 客户端**: Axios
- **Markdown 编辑器**: Toast UI Editor
- **代码高亮**: Prism.js
- **样式**: SCSS

## 📦 项目结构

```
src/
├── assets/          # 静态资源
├── components/      # 公共组件
│   ├── common/      # 通用组件
│   └── layout/      # 布局组件
├── router/          # 路由配置
├── stores/          # Pinia 状态管理
├── styles/          # 全局样式
├── utils/           # 工具函数
├── views/           # 页面组件
│   ├── admin/       # 管理后台页面
│   ├── article/     # 文章相关页面
│   └── error/       # 错误页面
├── App.vue          # 根组件
└── main.ts          # 入口文件
```

## 🚀 快速开始

### 环境要求

- Node.js >= 16.0.0
- npm >= 7.0.0

### 安装依赖

```bash
npm install
```

### 开发环境

```bash
npm run dev
```

访问 http://localhost:8080

### 生产构建

```bash
npm run build
```

### 预览构建结果

```bash
npm run preview
```

## 🔧 配置

### 环境变量

创建 `.env.local` 文件配置本地环境变量：

```env
# API 基础地址
VITE_API_BASE_URL=http://localhost:3000/api

# 应用标题
VITE_APP_TITLE=我的个人博客
```

### API 接口

项目使用 RESTful API，主要接口包括：

- `GET /api/articles` - 获取文章列表
- `GET /api/articles/:id` - 获取文章详情
- `GET /api/categories` - 获取分类列表
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/register` - 用户注册
- `GET /api/auth/profile` - 获取用户信息

## 📱 功能模块

### 前台功能

- **首页**: 文章列表展示，支持分页和搜索
- **文章详情**: Markdown 渲染，代码高亮
- **分类页面**: 按分类浏览文章
- **归档页面**: 按时间归档浏览
- **用户认证**: 登录、注册、个人资料管理

### 后台功能

- **仪表盘**: 数据统计和快速操作
- **文章管理**: 文章的增删改查
- **文章编辑**: Markdown 编辑器，支持图片上传
- **分类管理**: 分类的管理
- **用户管理**: 用户信息管理

## 🎨 设计特色

- **苹果风格**: 简洁优雅的设计语言
- **响应式布局**: 完美适配各种屏幕尺寸
- **流畅动画**: 丰富的过渡动画效果
- **主题切换**: 明暗主题无缝切换
- **现代化 UI**: 使用 Element Plus 组件库

## 🔐 权限管理

- **游客**: 可浏览文章，无法评论
- **普通用户**: 可浏览、评论文章
- **管理员**: 拥有所有权限，可管理文章和用户

## 📈 性能优化

- **路由懒加载**: 按需加载页面组件
- **代码分割**: 第三方库单独打包
- **图片优化**: 支持 WebP 格式
- **缓存策略**: 合理的缓存配置
- **Tree Shaking**: 移除未使用的代码

## 🤝 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [Vue.js](https://vuejs.org/) - 渐进式 JavaScript 框架
- [Element Plus](https://element-plus.org/) - Vue 3 组件库
- [Toast UI Editor](https://ui.toast.com/tui-editor) - Markdown 编辑器
- [Vite](https://vitejs.dev/) - 下一代前端构建工具

---

如有问题或建议，欢迎提交 Issue 或 Pull Request！

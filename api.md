### 1. 用户模块（User）

#### 认证相关
- **POST /api/auth/register** - 用户注册
```json
请求参数:
{
    "username": "string",
    "password": "string",
    "confirm_password": "string",
    "email": "string",
    "nickname": "string"
}

响应数据:
{
    "code": 200,
    "message": "注册成功",
    "data": {
        "user_id": "number",
        "username": "string",
        "nickname": "string",
        "email": "string",
        "avatar": "string",
        "created_at": "string"
    }
}
```

- **POST /api/auth/login** - 用户登录
```json
请求参数:
{
    "username": "string",
    "password": "string",
    "remember_me": "boolean"
}

响应数据:
{
    "code": 200,
    "message": "登录成功",
    "data": {
        "token": "string",
        "expires_in": "number",
        "user": {
            "user_id": "number",
            "username": "string",
            "nickname": "string",
            "email": "string",
            "avatar": "string",
            "role": "string"
        }
    }
}
```

- **POST /api/auth/logout** - 用户注销
```json
响应数据:
{
    "code": 200,
    "message": "注销成功"
}
```

- **POST /api/auth/refresh** - 刷新Token
```json
响应数据:
{
    "code": 200,
    "message": "刷新成功",
    "data": {
        "token": "string",
        "expires_in": "number"
    }
}
```

#### 用户信息管理
- **GET /api/users/profile** - 获取当前用户信息（需要认证）
```json
响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "user_id": "number",
        "username": "string",
        "nickname": "string",
        "email": "string",
        "avatar": "string",
        "bio": "string",
        "website": "string",
        "location": "string",
        "role": "string",
        "created_at": "string",
        "updated_at": "string",
        "article_count": "number",
        "comment_count": "number"
    }
}
```

- **PUT /api/users/profile** - 更新用户信息（需要认证）
```json
请求参数:
{
    "nickname": "string",
    "email": "string",
    "bio": "string",
    "website": "string",
    "location": "string"
}
```

- **PUT /api/users/password** - 修改密码（需要认证）
```json
请求参数:
{
    "old_password": "string",
    "new_password": "string",
    "confirm_password": "string"
}
```

- **PUT /api/users/avatar** - 更新头像（需要认证）
```json
请求参数:
{
    "avatar": "string (base64 或 文件ID)"
}
```

- **GET /api/users/{user_id}** - 获取其他用户公开信息
```json
响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "user_id": "number",
        "username": "string",
        "nickname": "string",
        "avatar": "string",
        "bio": "string",
        "website": "string",
        "location": "string",
        "created_at": "string",
        "article_count": "number",
        "recent_articles": [
            {
                "article_id": "number",
                "title": "string",
                "created_at": "string"
            }
        ]
    }
}
```
### 2. 文章模块（Article）

- **GET /api/articles** - 获取文章列表
```json
查询参数:
{
    "page": "number (默认1)",
    "page_size": "number (默认10, 最大50)",
    "keyword": "string (搜索关键词)",
    "category_id": "number (分类ID)",
    "tag_ids": "string (标签ID，逗号分隔)",
    "author_id": "number (作者ID)",
    "status": "string (published/draft)",
    "sort": "string (created_at/updated_at/views/likes)",
    "order": "string (desc/asc)"
}

响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "articles": [
            {
                "article_id": "number",
                "title": "string",
                "summary": "string",
                "cover_image": "string",
                "author": {
                    "user_id": "number",
                    "username": "string",
                    "nickname": "string",
                    "avatar": "string"
                },
                "category": {
                    "category_id": "number",
                    "name": "string"
                },
                "tags": [
                    {
                        "tag_id": "number",
                        "name": "string"
                    }
                ],
                "created_at": "string",
                "updated_at": "string",
                "views": "number",
                "likes": "number",
                "comments_count": "number",
                "status": "string"
            }
        ],
        "pagination": {
            "current_page": "number",
            "page_size": "number",
            "total_pages": "number",
            "total_count": "number",
            "has_next": "boolean",
            "has_prev": "boolean"
        }
    }
}
```

- **GET /api/articles/{article_id}** - 获取文章详情
```json
响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "article_id": "number",
        "title": "string",
        "content": "string",
        "summary": "string",
        "cover_image": "string",
        "author": {
            "user_id": "number",
            "username": "string",
            "nickname": "string",
            "avatar": "string",
            "bio": "string"
        },
        "category": {
            "category_id": "number",
            "name": "string",
            "description": "string"
        },
        "tags": [
            {
                "tag_id": "number",
                "name": "string"
            }
        ],
        "created_at": "string",
        "updated_at": "string",
        "views": "number",
        "likes": "number",
        "comments_count": "number",
        "status": "string",
        "is_liked": "boolean (当前用户是否点赞)",
        "prev_article": {
            "article_id": "number",
            "title": "string"
        },
        "next_article": {
            "article_id": "number",
            "title": "string"
        }
    }
}
```

- **POST /api/articles** - 创建文章（需要认证）
```json
请求参数:
{
    "title": "string",
    "content": "string",
    "summary": "string",
    "cover_image": "string",
    "category_id": "number",
    "tag_ids": "array[number]",
    "status": "string (draft/published)"
}

响应数据:
{
    "code": 201,
    "message": "创建成功",
    "data": {
        "article_id": "number",
        "title": "string",
        "status": "string",
        "created_at": "string"
    }
}
```

- **PUT /api/articles/{article_id}** - 更新文章（需要认证）
```json
请求参数:
{
    "title": "string",
    "content": "string",
    "summary": "string",
    "cover_image": "string",
    "category_id": "number",
    "tag_ids": "array[number]",
    "status": "string"
}
```

- **DELETE /api/articles/{article_id}** - 删除文章（需要认证）
```json
响应数据:
{
    "code": 200,
    "message": "删除成功"
}
```

- **POST /api/articles/{article_id}/like** - 点赞文章（需要认证）
```json
响应数据:
{
    "code": 200,
    "message": "点赞成功",
    "data": {
        "likes": "number"
    }
}
```

- **DELETE /api/articles/{article_id}/like** - 取消点赞（需要认证）
```json
响应数据:
{
    "code": 200,
    "message": "取消点赞成功",
    "data": {
        "likes": "number"
    }
}
```

- **GET /api/articles/{article_id}/related** - 获取相关文章
```json
查询参数:
{
    "limit": "number (默认5)"
}

响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": [
        {
            "article_id": "number",
            "title": "string",
            "summary": "string",
            "cover_image": "string",
            "created_at": "string",
            "views": "number"
        }
    ]
}
```
### 3. 分类/标签模块（Category/Tag）

#### 分类管理
- **GET /api/categories** - 获取分类列表
```json
查询参数:
{
    "include_count": "boolean (是否包含文章数量)"
}

响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": [
        {
            "category_id": "number",
            "name": "string",
            "description": "string",
            "slug": "string",
            "article_count": "number",
            "created_at": "string"
        }
    ]
}
```

- **GET /api/categories/{category_id}** - 获取分类详情
```json
响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "category_id": "number",
        "name": "string",
        "description": "string",
        "slug": "string",
        "article_count": "number",
        "created_at": "string"
    }
}
```

- **POST /api/categories** - 创建分类（需要管理员权限）
```json
请求参数:
{
    "name": "string",
    "description": "string",
    "slug": "string"
}
```

- **PUT /api/categories/{category_id}** - 更新分类（需要管理员权限）
```json
请求参数:
{
    "name": "string",
    "description": "string",
    "slug": "string"
}

响应数据:
{
    "code": 200,
    "message": "更新成功",
    "data": {
        "category_id": "number",
        "name": "string",
        "description": "string",
        "slug": "string",
        "updated_at": "string"
    }
}
```

- **DELETE /api/categories/{category_id}** - 删除分类（需要管理员权限）
```json
响应数据:
{
    "code": 200,
    "message": "删除成功"
}
```

#### 标签管理
- **GET /api/tags** - 获取标签列表
```json
查询参数:
{
    "keyword": "string (搜索关键词)",
    "limit": "number (返回数量限制)",
    "include_count": "boolean (是否包含文章数量)"
}

响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": [
        {
            "tag_id": "number",
            "name": "string",
            "color": "string",
            "article_count": "number",
            "created_at": "string"
        }
    ]
}
```

- **GET /api/tags/{tag_id}** - 获取标签详情
```json
响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "tag_id": "number",
        "name": "string",
        "color": "string",
        "article_count": "number",
        "created_at": "string"
    }
}
```

- **POST /api/tags** - 创建标签（需要认证）
```json
请求参数:
{
    "name": "string",
    "color": "string"
}

响应数据:
{
    "code": 201,
    "message": "创建成功",
    "data": {
        "tag_id": "number",
        "name": "string",
        "color": "string",
        "created_at": "string"
    }
}
```

- **PUT /api/tags/{tag_id}** - 更新标签（需要认证）
```json
请求参数:
{
    "name": "string",
    "color": "string"
}

响应数据:
{
    "code": 200,
    "message": "更新成功",
    "data": {
        "tag_id": "number",
        "name": "string",
        "color": "string",
        "updated_at": "string"
    }
}
```

- **DELETE /api/tags/{tag_id}** - 删除标签（需要认证）
```json
响应数据:
{
    "code": 200,
    "message": "删除成功"
}
```

### 4. 评论模块（Comment）

- **GET /api/articles/{article_id}/comments** - 获取文章评论
```json
查询参数:
{
    "page": "number (默认1)",
    "page_size": "number (默认10)",
    "sort": "string (created_at/likes)",
    "order": "string (desc/asc)"
}

响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "comments": [
            {
                "comment_id": "number",
                "content": "string",
                "author": {
                    "user_id": "number",
                    "username": "string",
                    "nickname": "string",
                    "avatar": "string"
                },
                "parent_id": "number (父评论ID，null表示顶级评论)",
                "reply_to": {
                    "user_id": "number",
                    "nickname": "string"
                },
                "likes": "number",
                "is_liked": "boolean",
                "created_at": "string",
                "updated_at": "string",
                "replies": [
                    {
                        "comment_id": "number",
                        "content": "string",
                        "author": {...},
                        "reply_to": {...},
                        "likes": "number",
                        "is_liked": "boolean",
                        "created_at": "string"
                    }
                ]
            }
        ],
        "pagination": {
            "current_page": "number",
            "page_size": "number",
            "total_pages": "number",
            "total_count": "number"
        }
    }
}
```

- **POST /api/articles/{article_id}/comments** - 添加评论（需要认证）
```json
请求参数:
{
    "content": "string",
    "parent_id": "number (可选，回复某条评论)",
    "reply_to_user_id": "number (可选，@某个用户)"
}

响应数据:
{
    "code": 201,
    "message": "评论成功",
    "data": {
        "comment_id": "number",
        "content": "string",
        "created_at": "string"
    }
}
```

- **PUT /api/comments/{comment_id}** - 更新评论（需要认证，仅限作者）
```json
请求参数:
{
    "content": "string"
}

响应数据:
{
    "code": 200,
    "message": "更新成功",
    "data": {
        "comment_id": "number",
        "content": "string",
        "updated_at": "string"
    }
}
```

- **DELETE /api/comments/{comment_id}** - 删除评论（需要认证，仅限作者或管理员）
```json
响应数据:
{
    "code": 200,
    "message": "删除成功"
}
```

- **POST /api/comments/{comment_id}/like** - 点赞评论（需要认证）
```json
响应数据:
{
    "code": 200,
    "message": "点赞成功",
    "data": {
        "likes": "number"
    }
}
```

- **DELETE /api/comments/{comment_id}/like** - 取消点赞评论（需要认证）
```json
响应数据:
{
    "code": 200,
    "message": "取消点赞成功",
    "data": {
        "likes": "number"
    }
}
```

### 5. 搜索模块（Search）

- **GET /api/search** - 全局搜索
```json
查询参数:
{
    "q": "string (搜索关键词)",
    "type": "string (搜索类型: all/articles/users/tags)",
    "page": "number (默认1)",
    "page_size": "number (默认10)"
}

响应数据:
{
    "code": 200,
    "message": "搜索成功",
    "data": {
        "articles": [
            {
                "article_id": "number",
                "title": "string",
                "summary": "string",
                "highlight": "string (高亮摘要)",
                "author": {...},
                "created_at": "string"
            }
        ],
        "users": [
            {
                "user_id": "number",
                "username": "string",
                "nickname": "string",
                "avatar": "string",
                "bio": "string"
            }
        ],
        "tags": [
            {
                "tag_id": "number",
                "name": "string",
                "article_count": "number"
            }
        ],
        "total_count": {
            "articles": "number",
            "users": "number",
            "tags": "number"
        }
    }
}
```

- **GET /api/search/suggestions** - 搜索建议
```json
查询参数:
{
    "q": "string (输入的关键词)"
}

响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": [
        {
            "text": "string (建议文本)",
            "type": "string (类型: article/tag/user)"
        }
    ]
}
```

### 6. 文件上传模块（Upload）

- **POST /api/upload/image** - 上传图片（需要认证）
```json
请求参数: (multipart/form-data)
{
    "file": "File (图片文件)",
    "type": "string (upload_type: avatar/cover/content)"
}

响应数据:
{
    "code": 200,
    "message": "上传成功",
    "data": {
        "file_id": "string",
        "url": "string",
        "filename": "string",
        "size": "number",
        "mime_type": "string",
        "width": "number",
        "height": "number",
        "created_at": "string"
    }
}
```

- **POST /api/upload/file** - 上传文件（需要认证）
```json
请求参数: (multipart/form-data)
{
    "file": "File",
    "type": "string (upload_type: attachment)"
}

响应数据:
{
    "code": 200,
    "message": "上传成功",
    "data": {
        "file_id": "string",
        "url": "string",
        "filename": "string",
        "size": "number",
        "mime_type": "string",
        "created_at": "string"
    }
}
```

- **GET /api/files/{file_id}** - 获取文件信息
```json
响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "file_id": "string",
        "url": "string",
        "filename": "string",
        "size": "number",
        "mime_type": "string",
        "created_at": "string"
    }
}
```

- **DELETE /api/files/{file_id}** - 删除文件（需要认证，仅限上传者）
```json
响应数据:
{
    "code": 200,
    "message": "删除成功"
}
```

### 7. 统计模块（Statistics）

- **GET /api/stats/overview** - 获取网站概览统计
```json
响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total_articles": "number",
        "total_users": "number",
        "total_comments": "number",
        "total_views": "number",
        "total_likes": "number"
    }
}
```

- **GET /api/stats/articles/popular** - 获取热门文章
```json
查询参数:
{
    "period": "string (week/month/year)",
    "limit": "number (默认10)"
}

响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": [
        {
            "article_id": "number",
            "title": "string",
            "summary": "string",
            "cover_image": "string",
            "author": {
                "user_id": "number",
                "username": "string",
                "nickname": "string",
                "avatar": "string"
            },
            "views": "number",
            "likes": "number",
            "comments_count": "number",
            "created_at": "string"
        }
    ]
}
```

- **GET /api/stats/articles/recent** - 获取最新文章
```json
查询参数:
{
    "limit": "number (默认10)"
}

响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": [
        {
            "article_id": "number",
            "title": "string",
            "summary": "string",
            "cover_image": "string",
            "author": {
                "user_id": "number",
                "username": "string",
                "nickname": "string",
                "avatar": "string"
            },
            "created_at": "string"
        }
    ]
}
```

### 8. 通知模块（Notification）

- **GET /api/notifications** - 获取通知列表（需要认证）
```json
查询参数:
{
    "page": "number",
    "page_size": "number",
    "type": "string (comment/like/follow/system)",
    "is_read": "boolean"
}

响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "notifications": [
            {
                "notification_id": "number",
                "type": "string",
                "title": "string",
                "content": "string",
                "data": "object (额外数据)",
                "is_read": "boolean",
                "created_at": "string"
            }
        ],
        "unread_count": "number",
        "pagination": {...}
    }
}
```

- **PUT /api/notifications/{notification_id}/read** - 标记通知为已读（需要认证）
```json
响应数据:
{
    "code": 200,
    "message": "标记成功",
    "data": {
        "notification_id": "number",
        "is_read": true
    }
}
```

- **PUT /api/notifications/read-all** - 标记所有通知为已读（需要认证）
```json
响应数据:
{
    "code": 200,
    "message": "标记成功",
    "data": {
        "count": "number (标记为已读的通知数量)"
    }
}
```

### 9. 管理员模块（Admin）

- **GET /api/admin/dashboard** - 管理员仪表板（需要管理员权限）
```json
响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "total_users": "number",
        "total_articles": "number",
        "total_comments": "number",
        "total_views": "number",
        "recent_users": [
            {
                "user_id": "number",
                "username": "string",
                "nickname": "string",
                "avatar": "string",
                "created_at": "string"
            }
        ],
        "recent_articles": [
            {
                "article_id": "number",
                "title": "string",
                "author": {
                    "user_id": "number",
                    "username": "string",
                    "nickname": "string"
                },
                "status": "string",
                "created_at": "string"
            }
        ],
        "stats_by_date": [
            {
                "date": "string",
                "users_count": "number",
                "articles_count": "number",
                "comments_count": "number",
                "views_count": "number"
            }
        ]
    }
}
```

- **GET /api/admin/users** - 用户管理列表（需要管理员权限）
```json
查询参数:
{
    "page": "number",
    "page_size": "number",
    "keyword": "string",
    "status": "string (active/banned)"
}

响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "users": [
            {
                "user_id": "number",
                "username": "string",
                "nickname": "string",
                "email": "string",
                "avatar": "string",
                "role": "string",
                "status": "string",
                "created_at": "string",
                "last_login_at": "string",
                "article_count": "number",
                "comment_count": "number"
            }
        ],
        "pagination": {
            "current_page": "number",
            "page_size": "number",
            "total_pages": "number",
            "total_count": "number"
        }
    }
}
```

- **PUT /api/admin/users/{user_id}/status** - 更改用户状态（需要管理员权限）
```json
请求参数:
{
    "status": "string (active/banned)"
}

响应数据:
{
    "code": 200,
    "message": "更新成功",
    "data": {
        "user_id": "number",
        "status": "string"
    }
}
```

- **GET /api/admin/articles** - 文章管理列表（需要管理员权限）
```json
查询参数:
{
    "page": "number",
    "page_size": "number",
    "keyword": "string",
    "status": "string (published/draft/pending)",
    "author_id": "number",
    "category_id": "number"
}

响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "articles": [
            {
                "article_id": "number",
                "title": "string",
                "summary": "string",
                "author": {
                    "user_id": "number",
                    "username": "string",
                    "nickname": "string"
                },
                "category": {
                    "category_id": "number",
                    "name": "string"
                },
                "status": "string",
                "views": "number",
                "likes": "number",
                "comments_count": "number",
                "created_at": "string",
                "updated_at": "string"
            }
        ],
        "pagination": {
            "current_page": "number",
            "page_size": "number",
            "total_pages": "number",
            "total_count": "number"
        }
    }
}
```

- **PUT /api/admin/articles/{article_id}/status** - 更改文章状态（需要管理员权限）
```json
请求参数:
{
    "status": "string (published/draft/pending)"
}

响应数据:
{
    "code": 200,
    "message": "更新成功",
    "data": {
        "article_id": "number",
        "status": "string"
    }
}
```
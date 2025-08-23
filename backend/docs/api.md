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

### 8. 管理员模块（Admin）

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

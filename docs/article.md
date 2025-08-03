### 3. 文章模块（Article）

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
    "sort": "string (created_at/updated_at/views)",
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

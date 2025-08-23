### 3. 文章模块（Article）

- **GET /api/articles/:id/:version** - 获取文章列表

```json
响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "id": "number",
        "created_at": "string",
        "updated_at": "string",
        "article_id": "number",
        "description": "string",
        "version": "number",
        "store_id": "number",
    }
}
```

- **GET /api/articles** - 获取文章列表

```json
查询参数:
{
    "page": "number (默认1)",
    "page_size": "number (默认10, 最大50)",
}

响应数据:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "articles": [
            {
                "id": "number",
                "created_at": "string",
                "updated_at": "string",
                "user_id : "number",
                "title": "string",
                "article_id": "number",
                "version" : "number",
                "category" : "string",
                "tags" : "string",
                "views" : "number",
            },
            ...
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

- **POST /api/articles** - 创建文章（需要认证）

```json
请求参数:
{
    "title": "string",
    "description": "string",
    "store_id": "number",
    "category": "string",
    "tags": "string",
}

响应数据:
{
    "code": 201,
    "message": "创建成功"
}
```

- **PUT /api/articles/:id** - 更新文章（需要认证）

```json
请求参数:
{
    "title": "string",
    "description": "string",
    "store_id": "number",
    "category": "string",
    "tags": "string",
}
```

- **DELETE /api/articles/:id** - 删除文章（需要认证）

```json
响应数据:
{
    "code": 200,
    "message": "删除成功"
}
```

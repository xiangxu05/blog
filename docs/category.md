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

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
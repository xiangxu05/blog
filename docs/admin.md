### 管理员模块（Admin）

说明：
- 所有接口默认返回结构为：
```json
{
  "code": 200,
  "message": "string",
  "data": {}
}
```
- 除特别说明外，以下接口均需要认证（登录）且需要管理员角色（`user_role=admin`）。

#### 1. 后台统计

- GET /api/admin/backend_statistics — 获取后台统计（需要认证，管理员）

```json
响应数据:
{
  "code": 200,
  "message": "获取后台统计信息成功",
  "data": {
    "articles_num": "number",
    "views": "number",
    "files_num": "number",
    "users_num": "number",
    "store_usage": "string", 
    "capacity": "string",
    "last_backup": "string"
  }
}
```

对照的公开统计接口（无需管理员）：

- GET /api/website_statistics — 获取网站概览统计（公开）

```json
响应数据:
{
  "code": 200,
  "message": "获取网站统计信息成功",
  "data": {
    "article_num": "number",
    "views": "number",
    "categories": "number",
    "last_update": "string",
    "last_login": "string"
  }
}
```

#### 2. 管理员文章管理

- GET /api/admin/articles — 获取文章列表（管理员视图，需认证+管理员）

```json
查询参数:
{
  "page": "number (默认1)",
  "pageSize": "number (默认10, 最大50)", 
  "category": "string",
  "sort": "string (默认 created_at_desc)",
  "search": "string",
  "search_type": "string (默认 fuzzy)",
  "search_fields": "string (默认 all)",
  "status": "string (按状态筛选)"
}

响应数据:
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "articles": [ { ... } ],
    "pagination": {
      "current_page": "number",
      "page_size": "number",
      "total_pages": "number",
      "total_count": "number"
    }
  }
}
```

- GET /api/admin/articles/{id} — 获取文章详情（需认证+管理员）

注意：当前实现路由为 `{id}`，但处理函数期望 `{id}/{version}` 参数对（否则会报参数绑定错误）。建议按业务修复为以下其一：
- 路由改为 `GET /api/admin/articles/{id}/{version}` 使用“按版本获取”，或
- 保持当前路由并改用“获取最新版本”的处理函数

```json
响应数据:
{
  "code": 200,
  "message": "获取成功",
  "data": { ... }
}
```

- POST /api/admin/articles — 创建文章（需认证+管理员）

```json
请求参数:
{
  "title": "string",
  "description": "string",
  "category": "string",
  "tags": "string",
  "store_id": "number"
}

响应数据:
{
  "code": 200,
  "message": "创建成功"
}
```

- PUT /api/admin/articles/{id} — 更新文章（需认证+管理员）

```json
请求参数:
{
  "title": "string",
  "description": "string",
  "category": "string",
  "tags": "string",
  "store_id": "number"
}

响应数据:
{
  "code": 200,
  "message": "更新成功"
}
```

- DELETE /api/admin/articles/{id} — 删除文章（需认证+管理员）

```json
响应数据:
{
  "code": 200,
  "message": "删除成功"
}
```

#### 3. 文件管理（Admin）

- GET /api/files — 获取文件列表（需认证+管理员）

```json
查询参数:
{
  "search": "string",
  "type": "string"
}

响应数据:
{
  "code": 200,
  "message": "获取文件列表成功",
  "data": [
    {
      "file_id": "number",
      "filename": "string",
      "size": "number",
      "mime_type": "string",
      "created_at": "number"
    }
  ]
}
```

- POST /api/files/upload — 上传文件（需认证+管理员，multipart/form-data，key: file）

```json
响应数据:
{
  "code": 200,
  "message": "文件上传成功",
  "data": {
    "file_id": "number",
    "filename": "string",
    "size": "number",
    "mime_type": "string",
    "created_at": "number"
  }
}
```

- DELETE /api/files/{file_id} — 删除文件（需认证+管理员）

```json
响应数据:
{
  "code": 200,
  "message": "文件删除成功"
}
```

- GET /api/files/backup — 下载备份（需认证+管理员）

```
返回文件内容（附件下载）
```

#### 4. 其他说明
- 所有管理员接口通过 `UserAuth()` + `RoleAuth()` 中间件校验，`RoleAuth` 要求 `user_role = admin`。
- 公开接口中存在文章访问计数（`IncViews()`）与登录、备份等时间指标更新（`UpdateLastLogin()`、`UpdateLastBackup()` 等），已体现在统计接口数据中。



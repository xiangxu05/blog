### 2. 文件上传模块（Upload）

- **POST /api/files/upload** - 上传文件（需要认证）

```json
请求参数: (multipart/form-data),"key": "file", "value": file

响应数据:
{
    "code": 200,
    "message": "上传成功",
    "data": {
        "file_id": "string",
        "filename": "string",
        "size": "number",
        "mime_type": "string",
        "created_at": "string"
    }
}
```

- **GET /api/files/** - 获取所有文件信息

```json

响应数据:
{
    "code": 200,
    "message": "查询成功",
    "data": {
        {
            "file_id": "string",
            "filename": "string",
            "size": "number",
            "mime_type": "string",
            "created_at": "string"
        },
        ...
    }
}
```

- **GET /api/files/{file_id}** - 获取文件信息

```json

响应数据:
{
    "code": 200,
    "message": "查询成功",
    "data": {
        "file_id": "string",
        "filename": "string",
        "size": "number",
        "mime_type": "string",
        "created_at": "string"
    }
}
```

- **GET /api/files/download/{file_id}** - 下载文件

```
  返回文件内容
```

- **DELETE /api/files/{file_id}** - 删除文件（需要认证，仅限上传者）

```json
响应数据:
{
    "code": 200,
    "message": "删除成功"
}
```

### 1. 用户模块（User）

#### 认证相关

- **POST /api/users/register** - 用户注册

```json
请求参数:
{
    "username": "string",
    "password": "string",
    "confirm_password": "string",
    "email": "string",
}

响应数据:
{
    "code": 201,
    "message": "注册成功"
}
```

- **POST /api/users/login** - 用户登录

```json
请求参数:
{
    "username": "string",
    "password": "string"
}

响应数据:
{
    "code": 200,
    "message": "登录成功",
    "data": {
        "expires_in": "number",
        "user": {
            "user_id": "number",
            "username": "string",
            "nickname": "string",
            "email": "string",
            "avatar": "string",
            "role": "string",
            "created_at": "string"
        }
    }
}
```

- **GET /api/users/logout** - 用户注销

```json
响应数据:
{
    "code": 200,
    "message": "注销成功"
}
```

- **DELETE /api/users/delete** - 删除用户

```json
响应数据:
{
    "code": 200,
    "message": "删除成功"
}
```

- **GET /api/users/refresh** - 刷新登出时间

```json
响应数据:
{
    "code": 200,
    "message": "刷新成功",
    "data": {
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
        "role": "string",
        "created_at": "string"
    }
}
```

- **PUT /api/users/profile** - 更新用户信息（需要认证）

```json
请求参数:
{
    "nickname": "string",
    "email": "string",
    "avatar": "string"
}
响应数据:
{
    "code": 200,
    "message": "更新成功"
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
响应数据:
{
    "code": 200,
    "message": "修改成功"
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
        "role": "string",
        "created_at": "string"
    }
}
```

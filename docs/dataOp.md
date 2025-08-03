以下是你提供内容的 **Markdown 格式整理版**，包括语法高亮、结构层次优化和可读性提升：

---

# 🚀 GORM Gen 自动生成 DAO 层使用指南

## 📂 文件查看记录

* 查看自动生成的 DAO 文件结构和接口：

  * `dao/gen.go`：第 1 - 50 行
  * `dao/tb_user.gen.go`：第 1 - 100 行（用户表）
  * `dao/tb_article.gen.go`：第 1 - 80 行（文章表）

---

## 📋 1. 初始化 DAO

首先需要初始化数据库连接并设置默认 DAO：

```go
package main

import (
    "blog/dao"
    "github.com/glebarez/sqlite"
    "gorm.io/gorm"
)

func main() {
    // 连接数据库
    db, err := gorm.Open(sqlite.Open("data/blogData.db"), &gorm.Config{})
    if err != nil {
        panic(err)
    }

    // 设置默认 DAO
    dao.SetDefault(db)

    // 现在可以使用 dao.User, dao.Article, dao.Category 等全局变量查询
}
```

---

## 🔍 2. 基础查询操作

### 查询单条记录

```go
// 查询第一个用户
user, err := dao.User.First()

// 根据 ID 查询
user, err := dao.User.Where(dao.User.ID.Eq(1)).First()

// 根据用户名查询
user, err := dao.User.Where(dao.User.Username.Eq("xiangxu")).First()
```

### 查询多条记录

```go
// 查询所有用户
users, err := dao.User.Find()

// 条件查询
adminUsers, err := dao.User.Where(dao.User.Role.Eq("admin")).Find()

// 分页查询
users, count, err := dao.User.FindByPage(0, 10) // offset=0, limit=10
```

---

## 🎯 3. 高级查询条件

### 多条件查询

```go
// AND 条件
users, err := dao.User.Where(
    dao.User.Role.Eq("admin"),
    dao.User.Email.Like("%@foxmail.com"),
).Find()

// OR 条件
users, err := dao.User.Where(
    dao.User.Role.Eq("admin"),
).Or(
    dao.User.Role.Eq("editor"),
).Find()
```

### 排序与限制

```go
// 排序
users, err := dao.User.Order(dao.User.CreatedAt.Desc()).Find()

// 限制数量
users, err := dao.User.Limit(5).Find()

// 跳过记录
users, err := dao.User.Offset(10).Limit(5).Find()
```

---

## 📰 4. 文章查询示例

### 基础查询

```go
// 查询所有文章
articles, err := dao.Article.Find()

// 标题模糊查询
articles, err := dao.Article.Where(
    dao.Article.Title.Like("%Go语言%"),
).Find()

// 分类查询
articles, err := dao.Article.Where(
    dao.Article.CategoryID.Eq(1),
).Find()
```

### 关联查询（预加载）

```go
// 查询文章并预加载用户信息
articles, err := dao.Article.Preload(dao.Article.User).Find()

// 查询文章并预加载分类和标签
articles, err := dao.Article.
    Preload(dao.Article.Category).
    Preload(dao.Article.Tags).
    Find()

// 查询单篇文章的完整信息
article, err := dao.Article.
    Preload(dao.Article.User).
    Preload(dao.Article.Category).
    Preload(dao.Article.Tags).
    Where(dao.Article.ID.Eq(1)).
    First()
```

---

## 🏷️ 5. 多对多关系操作（文章标签）

```go
article := &model_def.Article{ID: 1}
tags := []*model_def.Tag{{ID: 1}, {ID: 2}}

// 添加标签
err := dao.Article.Tags.Model(article).Append(tags...)

// 替换标签
err := dao.Article.Tags.Model(article).Replace(tags...)

// 删除标签
err := dao.Article.Tags.Model(article).Delete(tags...)

// 查询所有标签
tags, err := dao.Article.Tags.Model(article).Find()
```

---

## 💾 6. 增删改操作

### 创建记录

```go
// 创建新用户
newUser := &model_def.User{
    Username: "testuser",
    Password: "hashedpassword",
    Email:    "test@example.com",
    Role:     "user",
}
err := dao.User.Create(newUser)

// 批量创建
users := []*model_def.User{user1, user2, user3}
err := dao.User.CreateInBatches(users, 100)
```

### 更新记录

```go
// 更新单个字段
info, err := dao.User.Where(dao.User.ID.Eq(1)).Update(dao.User.Email, "new@example.com")

// 更新多个字段
info, err := dao.User.Where(dao.User.ID.Eq(1)).Updates(map[string]interface{}{
    "email": "new@example.com",
    "role":  "admin",
})
```

### 删除记录

```go
// 软删除
info, err := dao.User.Where(dao.User.ID.Eq(1)).Delete()

// 硬删除
info, err := dao.User.Unscoped().Where(dao.User.ID.Eq(1)).Delete()
```

---

## 📊 7. 统计查询

```go
// 总记录数
count, err := dao.User.Count()

// 条件统计
count, err := dao.User.Where(dao.User.Role.Eq("admin")).Count()

// 分组统计
type Result struct {
    Role  string
    Count int64
}

var results []Result
err := dao.User.Select(dao.User.Role, dao.User.ID.Count().As("count")).
    Group(dao.User.Role).
    Scan(&results)
```

---

## 🔄 8. 事务操作

```go
// 使用事务
err := dao.Q.Transaction(func(tx *dao.Query) error {
    if err := tx.User.Create(newUser); err != nil {
        return err
    }

    if err := tx.Article.Create(newArticle); err != nil {
        return err
    }

    return nil
})
```

---

## 🛠️ 9. 实用技巧

### 调试查询

```go
// 打印 SQL
users, err := dao.User.Debug().Where(dao.User.Role.Eq("admin")).Find()
```

### 原生 SQL 查询

```go
var users []model_def.User
err := dao.User.UnderlyingDB().Raw("SELECT * FROM tb_user WHERE role = ?", "admin").Scan(&users)
```

---

## ✅ 总结

GORM Gen 自动生成的 DAO 层提供：

* 🔍 类型安全：编译时字段检查
* ⚡ 高性能：预编译查询
* 🔗 关联查询：自动处理外键与预加载
* 📄 分页支持：内置分页函数
* 🛡️ 事务支持：复杂操作保证原子性
* 🐛 调试友好：支持 Debug 模式打印 SQL

> 利用这些功能，你可以高效、安全地进行业务数据库开发。

---

如需我继续为你整理代码片段、生成文档或添加搜索索引等功能，也可以继续提问！

package endpoint

import (
	"blog/dao"
	"blog/internal/message"
	"blog/internal/service"
	"blog/model_def"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserMsg struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

func GetUserHandler(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		message.SendMsg(c, http.StatusUnauthorized, "用户不存在", nil)
		return
	}
	user, err := dao.User.WithContext(c.Request.Context()).
		Where(dao.User.ID.Eq(userId.(int32))).First()
	if err != nil {
		message.SendMsg(c, http.StatusUnauthorized, "用户不存在", nil)
		return
	}
	resp := &UserMsg{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		Email:    user.Email,
		Avatar:   user.Avatar,
	}
	message.SendMsg(c, http.StatusOK, "获取用户成功", resp)
}

func UpdateUserHandler(c *gin.Context) {
	var updateUser *model_def.User
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}
	userId, ok := c.Get("user_id")
	if !ok {
		message.SendMsg(c, http.StatusUnauthorized, "用户不存在", nil)
		return
	}
	if err := service.UpdateUser(c, userId.(int32), updateUser); err != nil {
		message.SendMsg(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	message.SendMsg(c, http.StatusOK, "更新用户成功", nil)
}

func DeleteUserHandler(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		message.SendMsg(c, http.StatusUnauthorized, "用户不存在", nil)
		return
	}
	role, _ := c.Get("user_role")
	if role == "admin" {
		message.SendMsg(c, http.StatusForbidden, "不能删除管理员账号", nil)
		return
	}
	if err := service.DeleteUser(c, userId.(int32)); err != nil {
		message.SendMsg(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "删除用户成功", nil)
}

func RegisterHandler(c *gin.Context) {
	var req service.RegisterRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 调用注册服务
	resp, err := service.RegisterUser(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "服务器内部错误",
		})
		return
	}

	if resp.Success {
		c.JSON(http.StatusCreated, resp)
	} else {
		c.JSON(http.StatusBadRequest, resp)
	}
}

func LoginHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		message.SendMsg(c, http.StatusBadRequest, "request parameters are incorrect", nil)
		return
	}

	// 调用登录服务
	uid, err := service.LoginService(c, req.Username, req.Password)
	if err != nil {
		message.SendMsg(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	// 存储会话
	err = service.StoreSession(c, uid)
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "登录成功", nil)
}

func LogoutHandler(c *gin.Context) {
	if err := service.LogoutService(c); err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "注销失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "注销成功", nil)
}

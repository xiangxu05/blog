package endpoint

import (
	"blog/internal/message"
	"blog/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 注册
func RegisterHandler(c *gin.Context) {
	var req message.RegisterRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败"+err.Error(), nil)
		return
	}

	// 调用注册服务
	err := service.RegisterUser(c, &req)
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "注册失败: "+err.Error(), nil)
		return
	}

	message.SendMsg(c, http.StatusCreated, "注册成功", nil)
}

// 登录
func LoginHandler(c *gin.Context) {
	var req message.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message.SendMsg(c, http.StatusBadRequest, "请求参数错误！", nil)
		return
	}

	// 调用登录服务
	uid, err := service.LoginService(c, &req)
	if err != nil {
		message.SendMsg(c, http.StatusUnauthorized, "登录失败: "+err.Error(), nil)
		return
	}

	// 存储会话
	err = service.StoreSession(c, uid)
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "会话存储失败: "+err.Error(), nil)
		return
	}

	// 获取用户信息
	user, err := service.GetUser(c, uid)
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "获取用户信息失败: "+err.Error(), nil)
		return
	}

	// 生成响应
	resp := &message.LoginResponse{
		ExpiresIn: c.GetInt64("expires_in"),
		UserInfo: message.UserInfo{
			ID:          user.ID,
			Username:    user.Username,
			Nickname:    user.Nickname,
			Email:       user.Email,
			Avatar:      user.Avatar,
			Role:        user.Role,
			SelfIntro:   user.SelfIntro,
			PersonalWeb: user.PersonalWeb,
			Location:    user.Location,
			CreateAt:    user.CreatedAt.Unix(),
		},
	}
	message.SendMsg(c, http.StatusOK, "登录成功", resp)
}

// 注销登录
func LogoutHandler(c *gin.Context) {
	if err := service.LogoutService(c); err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "注销失败: "+err.Error(), nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "注销成功", nil)
}

// 删除用户
func DeleteUserHandler(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		message.SendMsg(c, http.StatusUnauthorized, "获取用户ID失败", nil)
		return
	}
	role, _ := c.Get("user_role")
	if role == "admin" {
		message.SendMsg(c, http.StatusForbidden, "不能删除管理员账号", nil)
		return
	}
	if err := service.DeleteUser(c, userId.(int32)); err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "删除用户失败: "+err.Error(), nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "删除成功", nil)
}

// 刷新会话
func RefreshHandler(c *gin.Context) {
	err := service.RefreshSession(c)
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "刷新会话失败: "+err.Error(), nil)
		return
	}
	m := map[string]int64{
		"expires_in": c.GetInt64("expires_in"),
	}
	message.SendMsg(c, http.StatusOK, "刷新成功", m)
}

func GetUserHandler(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		message.SendMsg(c, http.StatusUnauthorized, "用户不存在", nil)
		return
	}

	user, err := service.GetUser(c, userId.(int32))
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := &message.UserInfo{
		ID:          user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Email:       user.Email,
		Avatar:      user.Avatar,
		Role:        user.Role,
		SelfIntro:   user.SelfIntro,
		PersonalWeb: user.PersonalWeb,
		Location:    user.Location,
		CreateAt:    user.CreatedAt.Unix(),
	}
	message.SendMsg(c, http.StatusOK, "获取成功", resp)
}

func UpdateUserHandler(c *gin.Context) {
	var req message.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败: "+err.Error(), nil)
		return
	}
	userId, ok := c.Get("user_id")
	if !ok {
		message.SendMsg(c, http.StatusUnauthorized, "用户不存在", nil)
		return
	}
	if err := service.UpdateUser(c, userId.(int32), &req); err != nil {
		message.SendMsg(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	message.SendMsg(c, http.StatusOK, "更新成功", nil)
}

func UpdatePasswordHandler(c *gin.Context) {
	var req message.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}
	userId, ok := c.Get("user_id")
	if !ok {
		message.SendMsg(c, http.StatusUnauthorized, "用户不存在", nil)
		return
	}
	if err := service.UpdatePassword(c, userId.(int32), &req); err != nil {
		message.SendMsg(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "修改成功", nil)
}

func GetOtherUserHandler(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}
	user, err := service.GetUser(c, int32(user_id))
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := &message.GetUserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Role:        user.Role,
		SelfIntro:   user.SelfIntro,
		PersonalWeb: user.PersonalWeb,
		Location:    user.Location,
		CreateAt:    user.CreatedAt.Unix(),
	}
	message.SendMsg(c, http.StatusOK, "获取成功", resp)
}

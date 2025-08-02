package service

import (
	"blog/dao"
	"blog/model_def"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
)

func LoginService(c *gin.Context, username, password string) (int32, error) {
	ctx := c.Request.Context()
	user, err := dao.User.WithContext(ctx).Where(dao.User.Username.Eq(username)).First()
	if err != nil || user == nil {
		return 0, errors.New("用户未注册")
	}

	// 计算 SHA-256 哈希
	hash := hashPassword(password)
	if hash != user.Password {
		return 0, errors.New("password is incorrect")
	}

	return user.ID, nil
}

func LogoutService(c *gin.Context) error {
	sid, err := c.Cookie("sid")
	if err != nil {
		return err
	}
	if _, err := dao.Session.WithContext(c.Request.Context()).
		Where(dao.Session.SessionId.Eq(sid)).Delete(); err != nil {
		return err
	}
	return nil
}

type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=20"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  int32  `json:"user_id,omitempty"`
}

func RegisterUser(c *gin.Context, req *RegisterRequest) (*RegisterResponse, error) {
	ctx := c.Request.Context()
	// 基础验证
	if err := validateRegisterRequest(req); err != nil {
		return &RegisterResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// 检查用户名是否已存在
	existingUser, err := dao.User.WithContext(ctx).
		Where(dao.User.Username.Eq(req.Username)).
		First()
	if err == nil && existingUser != nil {
		return &RegisterResponse{
			Success: false,
			Message: "用户名已存在",
		}, nil
	}

	// 检查邮箱是否已存在
	existingEmail, err := dao.User.WithContext(ctx).
		Where(dao.User.Email.Eq(req.Email)).
		First()
	if err == nil && existingEmail != nil {
		return &RegisterResponse{
			Success: false,
			Message: "邮箱已被注册",
		}, nil
	}

	// 密码加密
	hashedPassword := hashPassword(req.Password)

	// 创建用户
	newUser := &model_def.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user", // 默认角色
	}

	err = dao.User.WithContext(ctx).Create(newUser)
	if err != nil {
		return &RegisterResponse{
			Success: false,
			Message: "注册失败，请稍后重试",
		}, err
	}

	return &RegisterResponse{
		Success: true,
		Message: "注册成功",
		UserID:  newUser.ID,
	}, nil
}

func UpdateUser(c *gin.Context, userID int32, updateUser *model_def.User) error {
	ctx := c.Request.Context()
	// 检查用户是否存在
	existingUser, err := dao.User.WithContext(ctx).
		Where(dao.User.ID.Eq(userID)).
		First()
	if err != nil || existingUser == nil {
		return errors.New("user does not exist")
	}
	// 更新用户信息
	if updateUser.Email != "" {
		existingUser.Email = updateUser.Email
	}
	if updateUser.Avatar != "" {
		existingUser.Avatar = updateUser.Avatar
	}
	if updateUser.Password != "" {
		existingUser.Password = hashPassword(updateUser.Password)
	}
	if updateUser.Nickname != "" {
		existingUser.Nickname = updateUser.Nickname
	}
	// 保存更新后的用户信息
	if _, err := dao.User.WithContext(ctx).
		Where(dao.User.ID.Eq(userID)).
		Updates(existingUser); err != nil {
		return err
	}
	LogoutService(c)
	return nil
}

func DeleteUser(c *gin.Context, userID int32) error {
	ctx := c.Request.Context()
	// 检查用户是否存在
	existingUser, err := dao.User.WithContext(ctx).
		Where(dao.User.ID.Eq(userID)).
		First()
	if err != nil || existingUser == nil {
		return errors.New("user does not exist")
	}
	err = dao.Q.Transaction(func(tx *dao.Query) error {
		// 先删除用户的会话
		if _, err = tx.Session.WithContext(ctx).
			Where(tx.Session.UserID.Eq(userID)).
			Unscoped().Delete(); err != nil {
			return err
		}

		// 硬删除用户（使用 Unscoped()）
		if _, err = tx.User.WithContext(ctx).
			Where(tx.User.ID.Eq(userID)).
			Unscoped().Delete(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// 验证注册请求
func validateRegisterRequest(req *RegisterRequest) error {
	// 密码确认
	if req.Password != req.ConfirmPassword {
		return errors.New("两次密码输入不一致")
	}

	// 用户名格式验证
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	if !usernameRegex.MatchString(req.Username) {
		return errors.New("用户名只能包含字母、数字和下划线，长度3-20个字符")
	}

	return nil
}

// 密码加密
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

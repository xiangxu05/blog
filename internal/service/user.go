package service

import (
	"blog/dao"
	"blog/internal/message"
	"blog/model_def"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUser(c *gin.Context, req *message.RegisterRequest) error {
	ctx := c.Request.Context()
	// 基础验证
	if err := validateRegisterRequest(req); err != nil {
		return err
	}

	// 检查用户名是否已存在
	existingUser, err := dao.User.WithContext(ctx).
		Where(dao.User.Username.Eq(req.Username)).
		First()
	if err == nil && existingUser != nil {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	existingEmail, err := dao.User.WithContext(ctx).
		Where(dao.User.Email.Eq(req.Email)).
		First()
	if err == nil && existingEmail != nil {
		return errors.New("邮箱已被注册")
	}

	// 密码加密
	hashedPassword := hashPassword(req.Password)

	// 创建用户
	newUser := &model_def.User{
		Nickname: req.Username,
		Username: req.Username,
		Password: hashedPassword,
		Role:     "user",
		Email:    req.Email,
	}

	err = dao.User.WithContext(ctx).Create(newUser)
	if err != nil {
		return errors.New("注册失败，请稍后重试" + err.Error())
	}

	return nil
}

// 验证注册请求
func validateRegisterRequest(req *message.RegisterRequest) error {
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

func LoginService(c *gin.Context, req *message.LoginRequest) (int32, error) {
	ctx := c.Request.Context()
	user, err := dao.User.WithContext(ctx).Where(dao.User.Username.Eq(req.Username)).First()
	if err != nil || user == nil {
		return -1, errors.New("用户未注册")
	}

	// 计算 SHA-256 哈希
	hash := hashPassword(req.Password)
	if hash != user.Password {
		return -1, errors.New("密码错误")
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

func DeleteUser(c *gin.Context, userID int32) error {
	ctx := c.Request.Context()

	err := dao.Q.Transaction(func(tx *dao.Query) error {
		// 先删除用户的会话
		if _, err := tx.Session.WithContext(ctx).
			Where(tx.Session.UserID.Eq(userID)).
			Unscoped().Delete(); err != nil {
			return err
		}

		// 硬删除用户（使用 Unscoped()）
		if _, err := tx.User.WithContext(ctx).
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

// GetAdminBlogger 取 role=admin 中 id 最小者作为博主；后续若有多名管理员，展示仍稳定指向首位管理员。
func GetAdminBlogger(c *gin.Context) (*model_def.User, error) {
	ctx := c.Request.Context()
	u, err := dao.User.WithContext(ctx).
		Where(dao.User.Role.Eq("admin")).
		Order(dao.User.ID.Asc()).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

func GetUser(c *gin.Context, userID int32) (*model_def.User, error) {
	ctx := c.Request.Context()
	existingUser, err := dao.User.WithContext(ctx).
		Where(dao.User.ID.Eq(userID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return existingUser, nil
}

func UpdateUser(c *gin.Context, userID int32, updateUser *message.UpdateProfileRequest) error {
	ctx := c.Request.Context()
	// 检查用户是否存在
	existingUser, err := dao.User.WithContext(ctx).
		Where(dao.User.ID.Eq(userID)).
		First()
	if err != nil || existingUser == nil {
		return errors.New("用户不存在")
	}
	// 更新用户信息
	if updateUser.Email != "" {
		existingUser.Email = updateUser.Email
	}
	if updateUser.Avatar != "" {
		existingUser.Avatar = updateUser.Avatar
	}
	if updateUser.Nickname != "" {
		existingUser.Nickname = updateUser.Nickname
	}
	if updateUser.SelfIntro != "" {
		existingUser.SelfIntro = updateUser.SelfIntro
	}
	if updateUser.PersonalWeb != "" {
		existingUser.PersonalWeb = updateUser.PersonalWeb
	}
	if updateUser.Location != "" {
		existingUser.Location = updateUser.Location
	}

	// 保存更新后的用户信息
	if _, err := dao.User.WithContext(ctx).
		Where(dao.User.ID.Eq(userID)).
		Updates(existingUser); err != nil {
		return err
	}
	return nil
}

func UpdatePassword(c *gin.Context, userID int32, updatePassword *message.UpdatePasswordRequest) error {
	ctx := c.Request.Context()
	// 检查用户是否存在
	existingUser, err := dao.User.WithContext(ctx).
		Where(dao.User.ID.Eq(userID)).
		First()
	if err != nil || existingUser == nil {
		return errors.New("用户不存在")
	}
	// 验证旧密码
	if hashPassword(updatePassword.OldPassword) != existingUser.Password {
		return errors.New("旧密码错误")
	}
	// 验证新密码
	if updatePassword.NewPassword != updatePassword.ConfirmPassword {
		return errors.New("两次新密码输入不一致")
	}

	// 加密新密码
	hashedPassword := hashPassword(updatePassword.NewPassword)
	existingUser.Password = hashedPassword
	// 保存更新后的用户信息
	if _, err := dao.User.WithContext(ctx).
		Where(dao.User.ID.Eq(userID)).
		Updates(existingUser); err != nil {
		return err
	}
	LogoutService(c)
	return nil
}

package service

import (
	"blog/dao"
	"blog/model_def"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func StoreSession(c *gin.Context, uid int32) error {
	// 先删除旧会话
	err := DeleteSessionByUID(c, uid)
	if err != nil {
		return err
	}

	// 设置会话 cookie, 生成uuid作为会话ID
	sessionId := uuid.New().String()
	// 将会话ID存储到数据库中
	expiresIn := time.Now().Add(24 * time.Hour).Unix()
	err = dao.Session.WithContext(c.Request.Context()).Create(&model_def.Session{
		UserID:    uid,
		SessionId: sessionId,
		Ctime:     time.Now().Unix(), // 创建时间
		Etime:     expiresIn,         // 过期时间，1小时后
	})
	if err != nil {
		return err
	}
	c.Set("expires_in", expiresIn)
	// 设置会话 cookie, 过期时间为24小时
	c.SetCookie("sid", sessionId, 24*3600, "/", "", false, true)
	return nil
}

func GetSession(c *gin.Context) (*model_def.Session, error) {
	sid, err := c.Cookie("sid")
	if err != nil {
		return nil, err
	}
	session, err := dao.Session.WithContext(c.Request.Context()).
		Where(dao.Session.SessionId.Eq(sid)).First()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func DeleteSession(c *gin.Context) error {
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

func DeleteSessionByUID(c *gin.Context, uid int32) error {
	if _, err := dao.Session.WithContext(c.Request.Context()).
		Where(dao.Session.UserID.Eq(uid)).Delete(); err != nil {
		return err
	}
	return nil
}

func RefreshSession(c *gin.Context) error {
	// 获取用户ID
	uid, exists := c.Get("user_id")
	if !exists {
		return errors.New("未找到用户ID")
	}
	// 存储新会话
	return StoreSession(c, uid.(int32))
}

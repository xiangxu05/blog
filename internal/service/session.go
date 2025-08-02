package service

import (
	"blog/dao"
	"blog/model_def"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func StoreSession(c *gin.Context, uid int32) error {
	// 设置会话 cookie, 生成uuid作为会话ID
	sessionId := uuid.New().String()
	// 将会话ID存储到数据库中
	dao.Session.WithContext(c.Request.Context()).Create(&model_def.Session{
		UserID:    uid,
		SessionId: sessionId,
		Ctime:     time.Now().Unix(),                    // 创建时间
		Etime:     time.Now().Add(1 * time.Hour).Unix(), // 过期时间，1小时后
	})
	// 设置会话 cookie
	c.SetCookie("sid", sessionId, 3600, "/", "", false, true)
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

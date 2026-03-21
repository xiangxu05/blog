package middleware

import (
	"blog/dao"
	"blog/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

// OptionalUserAuth 若存在有效会话则写入 user_id 等，否则继续（不 401）
func OptionalUserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := service.GetSession(c)
		if err != nil {
			c.Next()
			return
		}
		if session.Etime < time.Now().Unix() {
			c.Next()
			return
		}
		user, err := dao.User.WithContext(c.Request.Context()).
			Where(dao.User.ID.Eq(session.UserID)).First()
		if err != nil {
			c.Next()
			return
		}
		c.Set("user_id", int32(user.ID))
		c.Set("nickname", user.Nickname)
		c.Set("user_name", user.Username)
		c.Set("user_role", user.Role)
		c.Set("email", user.Email)
		c.Set("avatar", user.Avatar)
		c.Next()
	}
}

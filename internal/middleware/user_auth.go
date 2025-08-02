package middleware

import (
	"blog/dao"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie("sid")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}
		// 验证 sid 是否存在
		session, err := dao.Session.WithContext(c.Request.Context()).
			Where(dao.Session.SessionId.Eq(sid)).First()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效sid"})
			c.Abort()
			return
		}
		user, err := dao.User.WithContext(c.Request.Context()).
			Where(dao.User.ID.Eq(session.UserID)).First()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			c.Abort()
			return
		}

		c.Set("user_id", int32(user.ID))
		c.Set("user_name", user.Username)
		c.Set("user_role", user.Role)
		c.Set("nickname", user.Nickname)
		c.Next()
	}
}

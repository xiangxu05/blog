package middleware

import (
	"blog/internal/message"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("user_role")
		if !ok {
			message.SendMsg(c, http.StatusUnauthorized, "用户不存在", nil)
			c.Abort()
			return
		}
		if role != "admin" {
			message.SendMsg(c, http.StatusForbidden, "权限不足", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

package middleware

import (
	"blog/dao"
	"blog/internal/message"
	"blog/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 验证 sid 是否存在
		session, err := service.GetSession(c)
		if err != nil {
			message.SendMsg(c, http.StatusUnauthorized, "未登录", nil)
			c.Abort()
			return
		}

		// 查询用户是否存在
		user, err := dao.User.WithContext(c.Request.Context()).
			Where(dao.User.ID.Eq(session.UserID)).First()
		if err != nil {
			message.SendMsg(c, http.StatusUnauthorized, "用户不存在", nil)
			c.Abort()
			return
		}

		if session.Etime < time.Now().Unix() {
			message.SendMsg(c, http.StatusUnauthorized, "会话已过期", nil)
			err := service.DeleteSession(c)
			if err != nil {
				message.SendMsg(c, http.StatusInternalServerError, "会话删除失败", nil)
			}
			c.Abort()
			return
		}

		c.Set("user_id", int32(user.ID))
		c.Set("nickname", user.Nickname)
		c.Set("user_name", user.Username)
		c.Set("user_role", user.Role)
		c.Set("email", user.Email)
		c.Set("create_at", user.CreatedAt)
		c.Set("avatar", user.Avatar)
		c.Next()
	}
}

package endpoint

import (
	"blog/internal/message"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	message.SendMsg(c, http.StatusOK, "ok", nil)
}

package endpoint

import (
	"blog/internal/message"
	"blog/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ListCommentsHandler GET /api/articles/:id/comments
func ListCommentsHandler(c *gin.Context) {
	articleStr := c.Param("id")
	articleID, err := strconv.ParseInt(articleStr, 10, 32)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "无效的文章 ID", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var uid *int32
	if v, ok := c.Get("user_id"); ok {
		u := v.(int32)
		uid = &u
	}

	data, err := service.ListComments(c, int32(articleID), page, pageSize, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message.SendMsg(c, http.StatusNotFound, "文章不存在", nil)
			return
		}
		log.Errorf("获取评论失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "获取评论失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "获取成功", data)
}

// CreateCommentHandler POST /api/articles/:id/comments（需登录）
func CreateCommentHandler(c *gin.Context) {
	articleStr := c.Param("id")
	articleID, err := strconv.ParseInt(articleStr, 10, 32)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "无效的文章 ID", nil)
		return
	}

	uid, ok := c.Get("user_id")
	if !ok {
		message.SendMsg(c, http.StatusUnauthorized, "未登录", nil)
		return
	}
	userID := uid.(int32)

	req := &message.CommentCreateRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}

	data, err := service.CreateComment(c, int32(articleID), userID, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message.SendMsg(c, http.StatusNotFound, "文章不存在", nil)
			return
		}
		switch err.Error() {
		case "父评论不存在", "评论内容过短", "评论内容过长":
			message.SendMsg(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		log.Errorf("发表评论失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "发表评论失败", nil)
		return
	}
	message.SendMsg(c, http.StatusCreated, "评论成功", data)
}

// DeleteCommentHandler DELETE /api/comments/:comment_id（需登录）
func DeleteCommentHandler(c *gin.Context) {
	commentStr := c.Param("comment_id")
	cid, err := strconv.ParseUint(commentStr, 10, 32)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "无效的评论 ID", nil)
		return
	}

	uid, ok := c.Get("user_id")
	if !ok {
		message.SendMsg(c, http.StatusUnauthorized, "未登录", nil)
		return
	}
	userID := uid.(int32)

	role, _ := c.Get("user_role")
	roleStr, _ := role.(string)

	err = service.DeleteComment(c, uint(cid), userID, roleStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message.SendMsg(c, http.StatusNotFound, "评论不存在", nil)
			return
		}
		if err.Error() == "forbidden" {
			message.SendMsg(c, http.StatusForbidden, "无权删除该评论", nil)
			return
		}
		log.Errorf("删除评论失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "删除评论失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "删除成功", nil)
}

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

// GetArticleSocialHandler GET /api/articles/:id/social（可选登录：返回是否已赞/已收藏）
func GetArticleSocialHandler(c *gin.Context) {
	articleStr := c.Param("id")
	articleID, err := strconv.ParseInt(articleStr, 10, 32)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "无效的文章 ID", nil)
		return
	}
	var uid *int32
	if v, ok := c.Get("user_id"); ok {
		u := v.(int32)
		uid = &u
	}
	data, err := service.GetArticleSocial(c, int32(articleID), uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message.SendMsg(c, http.StatusNotFound, "文章不存在", nil)
			return
		}
		log.Errorf("获取文章社交状态失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "获取失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "获取成功", data)
}

// PostArticleLikeHandler POST /api/articles/:id/like
func PostArticleLikeHandler(c *gin.Context) {
	articleStr := c.Param("id")
	articleID, err := strconv.ParseInt(articleStr, 10, 32)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "无效的文章 ID", nil)
		return
	}
	uid := c.MustGet("user_id").(int32)
	if err := service.LikeArticle(c, int32(articleID), uid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message.SendMsg(c, http.StatusNotFound, "文章不存在", nil)
			return
		}
		log.Errorf("点赞失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "操作失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "已点赞", nil)
}

// DeleteArticleLikeHandler DELETE /api/articles/:id/like
func DeleteArticleLikeHandler(c *gin.Context) {
	articleStr := c.Param("id")
	articleID, err := strconv.ParseInt(articleStr, 10, 32)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "无效的文章 ID", nil)
		return
	}
	uid := c.MustGet("user_id").(int32)
	if err := service.UnlikeArticle(c, int32(articleID), uid); err != nil {
		log.Errorf("取消赞失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "操作失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "已取消赞", nil)
}

// ListMyBookmarksHandler GET /api/articles/bookmarks（须登录）
func ListMyBookmarksHandler(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}
	pageSizeStr := c.DefaultQuery("page_size", "10")
	pageSizeNum, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}
	if pageSizeNum > 50 {
		pageSizeNum = 50
	}
	uid := c.MustGet("user_id").(int32)
	data, err := service.ListMyBookmarks(c, uid, pageNum, pageSizeNum)
	if err != nil {
		log.Errorf("获取收藏列表失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "获取收藏列表失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "获取成功", data)
}

// PostArticleBookmarkHandler POST /api/articles/:id/bookmark
func PostArticleBookmarkHandler(c *gin.Context) {
	articleStr := c.Param("id")
	articleID, err := strconv.ParseInt(articleStr, 10, 32)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "无效的文章 ID", nil)
		return
	}
	uid := c.MustGet("user_id").(int32)
	if err := service.BookmarkArticle(c, int32(articleID), uid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message.SendMsg(c, http.StatusNotFound, "文章不存在", nil)
			return
		}
		log.Errorf("收藏失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "操作失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "已收藏", nil)
}

// DeleteArticleBookmarkHandler DELETE /api/articles/:id/bookmark
func DeleteArticleBookmarkHandler(c *gin.Context) {
	articleStr := c.Param("id")
	articleID, err := strconv.ParseInt(articleStr, 10, 32)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "无效的文章 ID", nil)
		return
	}
	uid := c.MustGet("user_id").(int32)
	if err := service.UnbookmarkArticle(c, int32(articleID), uid); err != nil {
		log.Errorf("取消收藏失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "操作失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "已取消收藏", nil)
}

// UpdateCommentHandler PUT /api/comments/:comment_id
func UpdateCommentHandler(c *gin.Context) {
	cidStr := c.Param("comment_id")
	cid, err := strconv.ParseUint(cidStr, 10, 32)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "无效的评论 ID", nil)
		return
	}
	uid := c.MustGet("user_id").(int32)
	req := &message.CommentUpdateRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		message.SendMsg(c, http.StatusBadRequest, "参数绑定失败", nil)
		return
	}
	err = service.UpdateComment(c, uint(cid), uid, req.Content)
	if err != nil {
		if err.Error() == "forbidden" {
			message.SendMsg(c, http.StatusForbidden, "无权修改", nil)
			return
		}
		if err.Error() == "评论内容长度无效" {
			message.SendMsg(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message.SendMsg(c, http.StatusNotFound, "评论不存在", nil)
			return
		}
		log.Errorf("更新评论失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "更新失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "更新成功", nil)
}

// PostCommentLikeHandler POST /api/comments/:comment_id/like
func PostCommentLikeHandler(c *gin.Context) {
	cidStr := c.Param("comment_id")
	cid, err := strconv.ParseUint(cidStr, 10, 32)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "无效的评论 ID", nil)
		return
	}
	uid := c.MustGet("user_id").(int32)
	if err := service.ToggleCommentLike(c, uint(cid), uid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message.SendMsg(c, http.StatusNotFound, "评论不存在", nil)
			return
		}
		log.Errorf("评论点赞失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "操作失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "已点赞", nil)
}

// DeleteCommentLikeHandler DELETE /api/comments/:comment_id/like
func DeleteCommentLikeHandler(c *gin.Context) {
	cidStr := c.Param("comment_id")
	cid, err := strconv.ParseUint(cidStr, 10, 32)
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "无效的评论 ID", nil)
		return
	}
	uid := c.MustGet("user_id").(int32)
	if err := service.RemoveCommentLike(c, uint(cid), uid); err != nil {
		log.Errorf("取消评论赞失败: %v", err)
		message.SendMsg(c, http.StatusInternalServerError, "操作失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "已取消", nil)
}

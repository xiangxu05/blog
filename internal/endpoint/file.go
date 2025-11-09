package endpoint

import (
	"blog/internal/message"
	"blog/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetFileHandler(c *gin.Context) {
	fileID := c.Param("file_id")
	fileInfo, err := service.GetFileInfo(c, fileID)
	if err != nil {
		message.SendMsg(c, http.StatusNotFound, "文件不存在", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "文件获取成功", fileInfo)
}

func GetFileListHandler(c *gin.Context) {
	// 获取分页参数
	page := c.DefaultQuery("page", "1")
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	// 支持 pageSize 和 page_size 两种格式
	pageSizeStr := c.Query("page_size")
	if pageSizeStr == "" {
		pageSizeStr = c.Query("pageSize")
	}
	if pageSizeStr == "" {
		pageSizeStr = "24"
	}
	pageSizeNum, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 24
	}
	if pageSizeNum > 100 {
		pageSizeNum = 100 // 限制最大页面大小
	}

	// 获取其他查询参数
	search := c.Query("search")
	fileType := c.Query("type")
	sort := c.DefaultQuery("sort", "created_at_desc")

	fileList, err := service.GetFileList(c, pageNum, pageSizeNum, search, fileType, sort)
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "获取文件列表失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "获取文件列表成功", fileList)
}

func UploadFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		message.SendMsg(c, http.StatusBadRequest, "文件上传失败", nil)
		return
	}
	fileInfo, err := service.UploadFile(c, file)
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "文件上传失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "文件上传成功", fileInfo)
}

func DownloadFileHandler(c *gin.Context) {
	fileID := c.Param("file_id")
	err := service.DownloadFile(c, fileID)
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "文件下载失败", nil)
		return
	}
	// 文件下载成功，不需要发送额外的JSON响应
	// service.DownloadFile 已经通过 c.FileAttachment 发送了文件内容
}

func DeleteFileHandler(c *gin.Context) {
	fileID := c.Param("file_id")
	err := service.DeleteFile(c, fileID)
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, "文件删除失败", nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "文件删除成功", nil)
}

func BackupHandler(c *gin.Context) {
	err := service.Backup(c)
	if err != nil {
		message.SendMsg(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	// 如果没有错误，service.Backup已经通过c.FileAttachment发送了文件
	// 不需要额外的响应
}

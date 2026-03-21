package endpoint

import (
	"blog/internal/message"
	"blog/internal/monitor"
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

// BackupStartHandler 异步创建备份任务，立即返回 job_id，避免长时间占用 HTTP 连接。
func BackupStartHandler(m *monitor.Monitor) gin.HandlerFunc {
	return func(c *gin.Context) {
		jobID, filename, err := service.StartBackupAsync(m)
		if err != nil {
			message.SendMsg(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		message.SendMsg(c, http.StatusOK, "备份任务已创建", gin.H{"job_id": jobID, "filename": filename})
	}
}

// BackupStatusHandler 查询备份任务状态：pending / running / done / failed
func BackupStatusHandler(c *gin.Context) {
	job := c.Query("job")
	if job == "" {
		message.SendMsg(c, http.StatusBadRequest, "缺少 job 参数", nil)
		return
	}
	state, filename, errMsg := service.GetBackupJobStatus(job)
	if state == "" {
		message.SendMsg(c, http.StatusNotFound, errMsg, nil)
		return
	}
	message.SendMsg(c, http.StatusOK, "ok", gin.H{
		"state":    state,
		"filename": filename,
		"error":    errMsg,
	})
}

// BackupDownloadHandler 下载已生成好的 zip（短响应，仅读盘发送文件）
func BackupDownloadHandler(c *gin.Context) {
	fn := c.Query("file")
	if fn == "" {
		message.SendMsg(c, http.StatusBadRequest, "缺少 file 参数", nil)
		return
	}
	if err := service.ServeBackupDownload(c, fn); err != nil {
		message.SendMsg(c, http.StatusNotFound, err.Error(), nil)
		return
	}
}

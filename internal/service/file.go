package service

import (
	"blog/dao"
	"blog/internal/message"
	"blog/model_def"
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const ArticleStorePath = "data/article/"
const ImageStorePath = "data/image/"

// type FileInfo struct {
// 	FileID    int32 `json:"file_id"`    // 文件ID
// 	Filename  string `json:"filename"`   // 文件名
// 	Size      int64  `json:"size"`       // 文件大小
// 	MimeType  string `json:"mime_type"`  // 文件类型
// 	CreatedAt string `json:"created_at"` // 创建时间
// }

func GetFileInfo(c *gin.Context, fileID string) (*message.FileInfo, error) {
	fileIDInt, err := strconv.Atoi(fileID)
	if err != nil {
		return nil, err
	}
	file, err := dao.File.WithContext(c.Request.Context()).Where(dao.File.ID.Eq(int32(fileIDInt))).First()
	if err != nil {
		return nil, err
	}
	fileInfo := &message.FileInfo{
		FileID:    file.ID,
		Filename:  file.Filename,
		Size:      file.Size,
		MimeType:  file.MimeType,
		CreatedAt: file.CreatedAt,
	}
	return fileInfo, nil
}

func GetFileList(c *gin.Context) ([]*message.FileInfo, error) {
	files, err := dao.File.WithContext(c.Request.Context()).Find()
	if err != nil {
		return nil, err
	}
	fileInfos := make([]*message.FileInfo, 0, len(files))
	for _, file := range files {
		fileInfos = append(fileInfos, &message.FileInfo{
			FileID:    file.ID,
			Filename:  file.Filename,
			Size:      file.Size,
			MimeType:  file.MimeType,
			CreatedAt: file.CreatedAt,
		})
	}
	return fileInfos, nil
}

func UploadFile(c *gin.Context, file *multipart.FileHeader) (*message.FileInfo, error) {
	// 保存到本地
	storePath, err := SaveFileToLocal(c, file)
	if err != nil {
		return nil, err
	}

	// 从上下文中获取用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		return nil, errors.New("user_id not found in context")
	}
	userIdInt, ok := userId.(int32)
	if !ok {
		return nil, errors.New("user_id type error")
	}

	// 保存文件信息到数据库
	fileInfo := &model_def.File{
		Filename:  file.Filename,
		StorePath: storePath,
		OwnerID:   userIdInt,
		Size:      file.Size,
		MimeType:  filepath.Ext(file.Filename),
		CreatedAt: time.Now().Unix(),
	}
	err = dao.File.WithContext(c.Request.Context()).Create(fileInfo)
	if err != nil {
		return nil, err
	}

	return &message.FileInfo{
		FileID:    fileInfo.ID,
		Filename:  fileInfo.Filename,
		Size:      fileInfo.Size,
		MimeType:  fileInfo.MimeType,
		CreatedAt: fileInfo.CreatedAt,
	}, nil
}

func SaveFileToLocal(c *gin.Context, file *multipart.FileHeader) (string, error) {
	//确保路径存在
	err := os.MkdirAll(ArticleStorePath, os.ModePerm)
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(ImageStorePath, os.ModePerm)
	if err != nil {
		return "", err
	}

	// 生成文件名
	uid := uuid.New().String()
	filename := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename)) + "-" + uid + filepath.Ext(file.Filename)
	fileType := filepath.Ext(filename)
	var filePath string
	if fileType == ".md" {
		// 生成文件路径
		filePath = filepath.Join(ArticleStorePath, filename)
	} else {
		// 生成文件路径
		filePath = filepath.Join(ImageStorePath, filename)
	}

	// 保存文件
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func DownloadFile(c *gin.Context, fileID string) error {
	// 从数据库获取文件信息
	fileIDInt, err := strconv.Atoi(fileID)
	if err != nil {
		return err
	}
	file, err := dao.File.WithContext(c.Request.Context()).Where(dao.File.ID.Eq(int32(fileIDInt))).First()
	if err != nil {
		return err
	}
	// 检查文件是否存在
	if _, err := os.Stat(file.StorePath); os.IsNotExist(err) {
		return errors.New("file not found")
	}
	// 发送文件
	c.FileAttachment(file.StorePath, file.Filename)
	return nil
}

func DeleteFile(c *gin.Context, fileID string) error {
	// 从数据库获取文件信息
	fileIDInt, err := strconv.Atoi(fileID)
	if err != nil {
		return err
	}
	file, err := dao.File.WithContext(c.Request.Context()).Where(dao.File.ID.Eq(int32(fileIDInt))).First()
	if err != nil {
		return err
	}
	// 检查文件是否存在
	if _, err = os.Stat(file.StorePath); os.IsNotExist(err) {
		return errors.New("file not found")
	}
	// 删除文件
	err = os.Remove(file.StorePath)
	if err != nil {
		return err
	}
	// 删除数据库记录
	_, err = dao.File.WithContext(c.Request.Context()).Delete(file)
	if err != nil {
		return err
	}
	return nil
}

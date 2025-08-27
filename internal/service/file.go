package service

import (
	"archive/zip"
	"blog/dao"
	"blog/internal/message"
	"blog/model_def"
	"errors"
	"fmt"
	"io"
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
const FilesStorePath = "data/files/"

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
		FileID:   file.ID,
		Filename: file.Filename,
		Size:     file.Size,
		MimeType: file.MimeType,
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
			FileID:   file.ID,
			Filename: file.Filename,
			Size:     file.Size,
			MimeType: file.MimeType,
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
	}
	err = dao.File.WithContext(c.Request.Context()).Create(fileInfo)
	if err != nil {
		return nil, err
	}

	return &message.FileInfo{
		FileID:   fileInfo.ID,
		Filename: fileInfo.Filename,
		Size:     fileInfo.Size,
		MimeType: fileInfo.MimeType,
	}, nil
}

func SaveFileToLocal(c *gin.Context, file *multipart.FileHeader) (string, error) {
	//确保路径存在
	err := os.MkdirAll(ArticleStorePath, os.ModePerm)
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(FilesStorePath, os.ModePerm)
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
		filePath = filepath.Join(FilesStorePath, filename)
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

// Backup 实现备份功能
func Backup(c *gin.Context) error {
	// 生成备份文件名（包含时间戳）
	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("backup_%s.zip", timestamp)
	backupPath := filepath.Join("backup", backupName)

	// 确保备份目录存在
	err := os.MkdirAll("backup", os.ModePerm)
	if err != nil {
		return fmt.Errorf("创建备份目录失败: %w", err)
	}

	// 打包数据目录
	err = ZipFolder("data", backupPath)
	if err != nil {
		return fmt.Errorf("备份失败: %w", err)
	}
	c.FileAttachment(backupPath, backupName)
	return nil
}

// ZipFolder 将指定文件夹打包成ZIP文件
func ZipFolder(sourceDir, zipPath string) error {
	// 创建ZIP文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("创建ZIP文件失败: %w", err)
	}
	defer zipFile.Close()

	// 创建ZIP写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历源目录
	return filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录本身
		if info.IsDir() {
			return nil
		}

		// 计算ZIP内的相对路径
		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			return err
		}

		// 统一使用正斜杠（ZIP标准）
		relPath = strings.ReplaceAll(relPath, "\\", "/")

		// 在ZIP中创建文件
		zipFileWriter, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// 打开源文件
		srcFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// 复制文件内容
		_, err = io.Copy(zipFileWriter, srcFile)
		return err
	})
}

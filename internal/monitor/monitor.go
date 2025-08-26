package monitor

import (
	"blog/dao"
	"blog/internal/message"
	"blog/model_def"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/disk"
	"github.com/xiangxu05/logger/v2"
	"gorm.io/gorm"
)

type Monitor struct {
	WebInfo *model_def.WebInfo
	mtx     sync.Mutex
	ctx     context.Context
}

var log = logger.GetLogger()

func NewMonitor(ctx context.Context) *Monitor {
	m := &Monitor{
		ctx: ctx,
		mtx: sync.Mutex{},
	}
	webInfo, err := GetWebInfo()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 记录不存在，创建默认记录
			now := time.Now()
			webInfo = &model_def.WebInfo{
				LastUpdate: now,
				LastLogin:  now,
				LastClear:  now,
				LastBackup: now,
			}
			m.WebInfo = webInfo
			if err := m.RefreshWebInfo(); err != nil {
				log.Errorf("Failed to create initial web info: %v", err)
			}
		} else {
			log.Errorf("Failed to query web info: %v", err)
		}
	} else {
		m.WebInfo = webInfo
	}
	return m
}

func (m *Monitor) Start() {
	// 启动监控
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-m.ctx.Done():
			saveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			m.ctx = saveCtx
			if err := m.RefreshWebInfo(); err != nil {
				log.Errorf("Failed to refresh web info: %v", err)
			}
			log.Infof("Monitor 保存监控数据并退出")
			return
		case <-ticker.C:
			// 刷新监控数据
			if err := m.RefreshWebInfo(); err != nil {
				log.Errorf("Failed to refresh web info: %v", err)
			}
		}
	}
}

func GetWebInfo() (*model_def.WebInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	webInfo, err := dao.WebInfo.WithContext(ctx).First()
	if err != nil {
		return nil, err
	}
	return webInfo, nil
}

func (m *Monitor) RefreshWebInfo() error {
	if m.WebInfo == nil {
		return fmt.Errorf("WebInfo is nil")
	}

	// 先查询是否存在记录
	existingInfo, err := dao.WebInfo.WithContext(m.ctx).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 记录不存在，创建新记录
			err = dao.WebInfo.WithContext(m.ctx).Create(m.WebInfo)
			return err
		}
		// 其他错误
		return err
	}

	// 记录存在，更新记录（保持原有的ID）
	m.WebInfo.ID = existingInfo.ID
	m.UpdateInfos()
	_, err = dao.WebInfo.WithContext(m.ctx).Where(dao.WebInfo.ID.Eq(m.WebInfo.ID)).Updates(m.WebInfo)
	return err
}

func (m *Monitor) UpdateInfos() error {
	articleNum, err := dao.Article.WithContext(m.ctx).Count()
	if err != nil {
		return err
	}
	fileNum, err := dao.File.WithContext(m.ctx).Count()
	if err != nil {
		return err
	}
	userNum, err := dao.User.WithContext(m.ctx).Count()
	if err != nil {
		return err
	}
	categoriesNum, err := dao.Category.WithContext(m.ctx).Count()
	if err != nil {
		return err
	}
	commentsNum, err := dao.Comment.WithContext(m.ctx).Count()
	if err != nil {
		return err
	}
	storeUsage, capacity, err := m.DiskCheck()
	if err != nil {
		return err
	}
	m.WebInfo.ArticleNum = int(articleNum)
	m.WebInfo.FileNum = int(fileNum)
	m.WebInfo.UserNum = int(userNum)
	m.WebInfo.CategoriesNum = int(categoriesNum)
	m.WebInfo.CommentsNum = int(commentsNum)
	m.WebInfo.StoreUsage = storeUsage
	m.WebInfo.Capacity = capacity
	return nil
}

func (m *Monitor) DiskCheck() (string, string, error) {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)

	// 获取当前项目所在分区的信息
	usage, err := disk.Usage(currentDir)
	if err != nil {
		fmt.Printf("Error getting disk usage: %v\n", err)
		return "", "", err
	}
	return formatBytes(uint64(usage.Used)), formatBytes(uint64(usage.Total)), nil
}

func formatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func (m *Monitor) IncViews() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.mtx.Lock()
		defer m.mtx.Unlock()
		m.WebInfo.Views++
		c.Next()
	}
}

func (m *Monitor) UpdateLastUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.mtx.Lock()
		defer m.mtx.Unlock()
		m.WebInfo.LastUpdate = time.Now()
		c.Next()
	}
}

func (m *Monitor) UpdateLastLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.mtx.Lock()
		defer m.mtx.Unlock()
		m.WebInfo.LastLogin = time.Now()
		c.Next()
	}
}

func (m *Monitor) UpdateLastClear() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.mtx.Lock()
		defer m.mtx.Unlock()
		m.WebInfo.LastClear = time.Now()
		c.Next()
	}
}

func (m *Monitor) UpdateLastBackup() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.mtx.Lock()
		defer m.mtx.Unlock()
		m.WebInfo.LastBackup = time.Now()
		c.Next()
	}
}

func (m *Monitor) GetWebsiteStatisticsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		websiteStatisticsResponse := &message.WebsiteStatisticsResponse{
			ArticlesNum:   m.WebInfo.ArticleNum,
			Views:         m.WebInfo.Views,
			CategoriesNum: m.WebInfo.CategoriesNum,
			LastUpdate:    m.WebInfo.LastUpdate,
			LastLogin:     m.WebInfo.LastLogin,
		}
		message.SendMsg(c, http.StatusOK, "获取网站统计信息成功", websiteStatisticsResponse)
	}
}

func (m *Monitor) GetBackendStatisticsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		backendStatisticsResponse := &message.BackendStatisticsResponse{
			ArticlesNum: m.WebInfo.ArticleNum,
			Views:       m.WebInfo.Views,
			FilesNum:    m.WebInfo.FileNum,
			UsersNum:    m.WebInfo.UserNum,
			StoreUsage:  m.WebInfo.StoreUsage,
			Capacity:    m.WebInfo.Capacity,
			LastBackup:  m.WebInfo.LastBackup,
		}
		message.SendMsg(c, http.StatusOK, "获取后台统计信息成功", backendStatisticsResponse)
	}
}

// type WebInfo struct {
// 	ID            uint      `gorm:"primaryKey"`
// 	ArticleNum    int       `gorm:"column:article_num"`
// 	FileNum       int       `gorm:"column:file_num"`
// 	UserNum       int       `gorm:"column:user_num"`
// 	CategoriesNum int       `gorm:"column:categories"`
// 	CommentsNum   int       `gorm:"column:comments"`
// 	StoreUsage    int       `gorm:"column:store_usage"`
// 	Capacity      int       `gorm:"column:capacity"`
// 	Views         int       `gorm:"column:views"`
// 	LastUpdate    time.Time `gorm:"column:last_update"`
// 	LastLogin     time.Time `gorm:"column:last_login"`
// 	LastClear     time.Time `gorm:"column:last_clear"`
// 	LastBackup    time.Time `gorm:"column:last_backup"`
// }

package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"blog/internal/monitor"

	"github.com/gin-gonic/gin"
)

var backupNameRe = regexp.MustCompile(`^backup_\d{8}_\d{6}\.zip$`)

type backupJob struct {
	State     string // pending, running, done, failed
	Filename  string
	ErrMsg    string
	CreatedAt time.Time
}

var (
	backupJobsMu sync.RWMutex
	backupJobs   = make(map[string]*backupJob) // key: job id (timestamp 20060102_150405)
)

// StartBackupAsync 立即返回任务 ID，在后台 goroutine 中打包；连接不必长时间占用。
func StartBackupAsync(m *monitor.Monitor) (jobID string, filename string, err error) {
	ts := time.Now().Format("20060102_150405")
	jobID = ts
	filename = fmt.Sprintf("backup_%s.zip", ts)

	if err := os.MkdirAll("backup", 0o755); err != nil {
		return "", "", fmt.Errorf("创建备份目录失败: %w", err)
	}

	backupJobsMu.Lock()
	if _, exists := backupJobs[jobID]; exists {
		backupJobsMu.Unlock()
		return "", "", errors.New("该秒已有备份任务，请稍后重试")
	}
	backupJobs[jobID] = &backupJob{
		State:     "pending",
		Filename:  filename,
		CreatedAt: time.Now(),
	}
	backupJobsMu.Unlock()

	go runBackupJob(m, jobID, filename)
	return jobID, filename, nil
}

func runBackupJob(m *monitor.Monitor, jobID, filename string) {
	setJob := func(state, errMsg string) {
		backupJobsMu.Lock()
		defer backupJobsMu.Unlock()
		j, ok := backupJobs[jobID]
		if !ok {
			return
		}
		j.State = state
		j.ErrMsg = errMsg
	}

	setJob("running", "")

	tmpPath := filepath.Join("backup", filename+".tmp")
	finalPath := filepath.Join("backup", filename)

	// 清理可能残留的 tmp
	_ = os.Remove(tmpPath)

	if err := ZipFolder("data", tmpPath); err != nil {
		_ = os.Remove(tmpPath)
		setJob("failed", err.Error())
		return
	}

	if err := os.Rename(tmpPath, finalPath); err != nil {
		_ = os.Remove(tmpPath)
		setJob("failed", err.Error())
		return
	}

	setJob("done", "")
	if m != nil {
		m.MarkBackupCompleted()
	}
}

// GetBackupJobStatus 查询任务状态：pending / running / done / failed
func GetBackupJobStatus(jobID string) (state, filename, errMsg string) {
	backupJobsMu.RLock()
	defer backupJobsMu.RUnlock()
	j, ok := backupJobs[jobID]
	if !ok {
		return "", "", "任务不存在或已过期"
	}
	return j.State, j.Filename, j.ErrMsg
}

// ServeBackupDownload 仅发送已生成好的 zip（短连接，不再在请求里慢慢打 zip）。
func ServeBackupDownload(c *gin.Context, filename string) error {
	if !backupNameRe.MatchString(filename) {
		return errors.New("非法文件名")
	}
	absPath, err := filepath.Abs(filepath.Join("backup", filename))
	if err != nil {
		return err
	}
	root, err := filepath.Abs("backup")
	if err != nil {
		return err
	}
	if !strings.HasPrefix(absPath, root+string(os.PathSeparator)) && absPath != root {
		return errors.New("非法路径")
	}
	st, err := os.Stat(absPath)
	if err != nil || st.IsDir() {
		return errors.New("文件不存在或尚未就绪")
	}
	c.FileAttachment(absPath, filename)
	return nil
}

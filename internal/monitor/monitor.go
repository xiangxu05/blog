package monitor

import (
	"blog/dao"
	"blog/model_def"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

type Monitor struct {
	WebInfo *model_def.WebInfo
	Mtx     sync.Mutex
	ctx     context.Context
}

func NewMonitor(ctx context.Context) *Monitor {
	m := &Monitor{
		ctx: ctx,
	}
	webInfo, err := dao.WebInfo.WithContext(ctx).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 记录不存在，创建默认记录
			now := time.Now()
			webInfo = &model_def.WebInfo{
				ArticleNum: 0,
				FileNum:    0,
				UserNum:    0,
				Views:      0,
				Categories: 0,
				Users:      0,
				Comments:   0,
				StoreUsage: 0,
				Capacity:   0,
				LastUpdate: now,
				LastLogin:  now,
				LastClear:  now,
				LastBackup: now,
			}
			m.WebInfo = webInfo
			if err := m.RefreshWebInfo(); err != nil {
				log.Printf("Failed to create initial web info: %v", err)
			}
		} else {
			log.Printf("Failed to query web info: %v", err)
		}
	} else {
		m.WebInfo = webInfo
	}
	return m
}

func (m *Monitor) Start() {
	// 启动监控
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-m.ctx.Done():
				log.Println("Monitor stopped")
				return
			case <-ticker.C:
				// 刷新监控数据
				if err := m.RefreshWebInfo(); err != nil {
					log.Printf("Failed to refresh web info: %v", err)
				}
			}
		}
	}()
}

func (m *Monitor) GetWebInfo() (*model_def.WebInfo, error) {
	webInfo, err := dao.WebInfo.WithContext(m.ctx).First()
	if err != nil {
		return nil, err
	}
	m.WebInfo = webInfo
	if err != nil {
		return nil, err
	}
	return m.WebInfo, nil
}

func (m *Monitor) RefreshWebInfo() error {
	if m.WebInfo == nil {
		return fmt.Errorf("WebInfo is nil")
	}

	// 更新时间戳
	m.WebInfo.LastUpdate = time.Now()

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
	_, err = dao.WebInfo.WithContext(m.ctx).Updates(m.WebInfo)
	return err
}

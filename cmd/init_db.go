package cmd

import (
	"blog/model_def"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var initDbCmd = &cobra.Command{
	Use:   "init_db",
	Short: "初始化数据库",
	Run: func(cmd *cobra.Command, args []string) {
		InitDataBase()
	},
}

func init() {
	rootCmd.AddCommand(initDbCmd)
}

// 清空数据库所有数据（保留表结构）
func clearAllData(db *gorm.DB) error {
	// 按依赖关系逆序删除数据，避免外键约束问题
	// 使用 Unscoped() 进行硬删除，彻底清空数据
	if err := db.Unscoped().Where("1 = 1").Delete(&model_def.Article{}).Error; err != nil {
		return fmt.Errorf("清空 articles 表失败: %v", err)
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&model_def.ArticleTag{}).Error; err != nil {
		return fmt.Errorf("清空 article_tags 表失败: %v", err)
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&model_def.ArticleContent{}).Error; err != nil {
		return fmt.Errorf("清空 article_contents 表失败: %v", err)
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&model_def.Comment{}).Error; err != nil {
		return fmt.Errorf("清空 comments 表失败: %v", err)
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&model_def.Category{}).Error; err != nil {
		return fmt.Errorf("清空 categories 表失败: %v", err)
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&model_def.Tag{}).Error; err != nil {
		return fmt.Errorf("清空 tags 表失败: %v", err)
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&model_def.User{}).Error; err != nil {
		return fmt.Errorf("清空 users 表失败: %v", err)
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&model_def.Session{}).Error; err != nil {
		return fmt.Errorf("清空 sessions 表失败: %v", err)
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&model_def.File{}).Error; err != nil {
		return fmt.Errorf("清空 file 表失败: %v", err)
	}
	return nil
}

func InitDataBase() {
	// 使用相对路径连接 SQLite 数据库（相对于 blog 目录）
	dbPath := "data/blogData.db"

	// 确保数据库目录存在
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("无法创建数据库目录 %s: %v", dbDir, err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Auto migrate tables
	if err := db.AutoMigrate(&model_def.User{}); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	if err := db.AutoMigrate(&model_def.Article{}); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	if err := db.AutoMigrate(&model_def.ArticleContent{}); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	if err := db.AutoMigrate(&model_def.Category{}); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	if err := db.AutoMigrate(&model_def.ArticleTag{}); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	if err := db.AutoMigrate(&model_def.Tag{}); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	if err := db.AutoMigrate(&model_def.Comment{}); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	if err := db.AutoMigrate(&model_def.Session{}); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	if err := db.AutoMigrate(&model_def.File{}); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// Clear all table data
	log.Println("Starting to clear all table data...")
	if err := clearAllData(db); err != nil {
		log.Printf("Failed to clear data: %v", err)
	}
	// Create default admin user with encrypted password (SHA256)
	adminUser := &model_def.User{Username: "admin", Password: "123456789"}
	// Calculate SHA-256 hash
	hash := sha256.Sum256([]byte(adminUser.Password))

	// Convert to hexadecimal string
	hashStr := hex.EncodeToString(hash[:])
	adminUser.Password = hashStr
	adminUser.Role = "admin"
	adminUser.Email = "admin@qq.com"
	if err := db.Create(adminUser).Error; err != nil {
		log.Printf("Failed to create default user: %v", err)
	}

	log.Println("All table data cleared")

	// Close database connection
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

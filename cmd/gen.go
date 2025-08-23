package cmd

import (
	"blog/model_def"
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "生成数据库迁移文件",
	Run: func(cmd *cobra.Command, args []string) {
		GenDataBase()
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}

func GenDataBase() {
	// 使用绝对路径连接 SQLite 数据库
	dbPath := "data/blogData.db"
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 创建生成器实例
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./dao",   // 输出目录
		ModelPkgPath: "./model", // model 包路径
		Mode:         gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(db) // Bind database

	// Use existing model definitions to avoid compatibility issues with reverse generation
	fmt.Println("Starting DAO layer code generation...")

	// Apply basic models
	g.ApplyBasic(
		&model_def.User{},
		&model_def.Article{},
		&model_def.ArticleVersion{},
		&model_def.Comment{},
		&model_def.Session{},
		&model_def.File{},
		&model_def.Category{},
	)

	fmt.Println("Starting code generation execution...")
	// Execute code generation
	g.Execute()

	fmt.Println("DAO layer code generation completed!")

	// Close database connection
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

package cmd

import (
	"blog/dao"
	"blog/internal/router"

	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var blogsysCmd = &cobra.Command{
	Use:   "blogsys",
	Short: "启动博客系统",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := gorm.Open(sqlite.Open("data/blogData.db"), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		dao.SetDefault(db)

		r := router.InitRouter()
		r.Run(":8080")
	},
}

func init() {
	rootCmd.AddCommand(blogsysCmd)
}

package cmd

import (
	"blog/dao"
	"blog/internal/server"

	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var port string

var blogsysCmd = &cobra.Command{
	Use:   "blogsys",
	Short: "启动博客系统",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := gorm.Open(sqlite.Open("data/blogData.db"), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		dao.SetDefault(db)

		server := server.NewServer(db)
		if err := server.Run(":" + port); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(blogsysCmd)
	blogsysCmd.Flags().StringVarP(&port, "port", "p", "8080", "端口")
}

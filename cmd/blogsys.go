package cmd

import (
	"blog/dao"
	appmigrate "blog/internal/migrate"
	"blog/internal/server"
	"fmt"
	"strconv"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var port string

// normalizePort 修正误用 `-port=8081` 时被 Cobra 解析成 `-p` 值为 `ort=8081` 的情况。
// 正确写法请使用：`--port=8081`、`-p 8081` or `-p8081`。
func normalizePort(s string) (string, error) {
	s = strings.TrimSpace(s)
	if strings.Contains(s, "=") {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) == 2 {
			s = strings.TrimSpace(parts[1])
		}
	}
	n, err := strconv.ParseUint(s, 10, 16)
	if err != nil || n == 0 {
		return "", fmt.Errorf("无效端口 %q，请使用 --port=8081 或 -p 8081", s)
	}
	return strconv.FormatUint(n, 10), nil
}

var blogsysCmd = &cobra.Command{
	Use:   "blogsys",
	Short: "启动博客系统",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := gorm.Open(sqlite.Open("data/blogData.db"), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		if err := appmigrate.AutoMigrate(db); err != nil {
			panic(err)
		}

		dao.SetDefault(db)

		listenPort, err := normalizePort(port)
		if err != nil {
			panic(err)
		}

		server := server.NewServer(db)
		if err := server.Run(":" + listenPort); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(blogsysCmd)
	blogsysCmd.Flags().StringVarP(&port, "port", "p", "80", "监听端口（可用 --port=8081 或 -p 8081，勿写 -port=8081）")
}

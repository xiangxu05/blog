package cmd

import (
	"github.com/spf13/cobra"
)

var regenCmd = &cobra.Command{
	Use:   "regen",
	Short: "重新生成数据库",
	Run: func(cmd *cobra.Command, args []string) {
		// 先清空数据库
		InitDataBase()
		// 再重新生成数据库
		GenDataBase()
	},
}

func init() {
	rootCmd.AddCommand(regenCmd)
}

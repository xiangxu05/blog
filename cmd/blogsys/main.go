package main

import (
	"blog/dao"
	"blog/internal/router"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("data/blogData.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dao.SetDefault(db)

	r := router.InitRouter()
	r.Run(":80")
}

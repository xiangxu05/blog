package message

import "time"

type WebsiteStatisticsResponse struct {
	ArticlesNum   int       `json:"article_num"`
	Views         int       `json:"views"`
	CategoriesNum int       `json:"categories"`
	LastUpdate    time.Time `json:"last_update"`
	LastLogin     time.Time `json:"last_login"`
}

type BackendStatisticsResponse struct {
	ArticlesNum int       `json:"articles_num"`
	Views       int       `json:"views"`
	FilesNum    int       `json:"files_num"`
	UsersNum    int       `json:"users_num"`
	StoreUsage  string    `json:"store_usage"`
	Capacity    string    `json:"capacity"`
	LastBackup  time.Time `json:"last_backup"`
}

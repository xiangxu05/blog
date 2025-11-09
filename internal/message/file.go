package message

type FileInfo struct {
	FileID    int32  `json:"file_id"`    // 文件ID
	Filename  string `json:"filename"`   // 文件名
	Size      int64  `json:"size"`       // 文件大小
	MimeType  string `json:"mime_type"`  // 文件类型
	CreatedAt int64  `json:"created_at"` // 创建时间
}

type FileListResponse struct {
	Files      []FileInfo     `json:"files"`
	Pagination PaginationInfo `json:"pagination"`
}

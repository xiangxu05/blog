package service

import (
	"blog/dao"
	"blog/internal/message"
	"blog/model_def"
	"context"
	"errors"
	"math"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	commentMaxLen = 2000
	commentMinLen = 1
)

// ListComments 分页获取文章评论（时间正序，便于阅读）；uid 非空时填充是否已赞
func ListComments(c *gin.Context, articleID int32, page, pageSize int, uid *int32) (*message.CommentListResponse, error) {
	ctx := c.Request.Context()

	_, err := dao.Article.WithContext(ctx).Where(dao.Article.ID.Eq(articleID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	total, err := dao.Comment.WithContext(ctx).Where(dao.Comment.ArticleID.Eq(articleID)).Count()
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	rows, err := dao.Comment.WithContext(ctx).
		Where(dao.Comment.ArticleID.Eq(articleID)).
		Order(dao.Comment.CreatedAt).
		Offset(offset).
		Limit(pageSize).
		Find()
	if err != nil {
		return nil, err
	}

	userIDs := make(map[int32]struct{})
	var parentIDList []uint
	seenParent := make(map[uint]struct{})
	for _, row := range rows {
		userIDs[row.UserID] = struct{}{}
		if row.ParentID != nil {
			pid := uint(*row.ParentID)
			if _, ok := seenParent[pid]; !ok {
				seenParent[pid] = struct{}{}
				parentIDList = append(parentIDList, pid)
			}
		}
	}

	parentMap := make(map[uint]*model_def.Comment)
	if len(parentIDList) > 0 {
		var parents []*model_def.Comment
		db := dao.Comment.WithContext(ctx).UnderlyingDB()
		if err := db.Where("id IN ?", parentIDList).Find(&parents).Error; err != nil {
			return nil, err
		}
		for _, p := range parents {
			parentMap[p.ID] = p
			userIDs[p.UserID] = struct{}{}
		}
	}

	users, err := loadUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	items := make([]message.CommentItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, buildCommentItem(row, users, parentMap))
	}
	for i := range items {
		n, err := CommentLikeCount(c, items[i].ID)
		if err != nil {
			return nil, err
		}
		items[i].LikeCount = int(n)
		if uid != nil {
			liked, err := UserLikedComment(c, items[i].ID, *uid)
			if err != nil {
				return nil, err
			}
			items[i].Liked = liked
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	return &message.CommentListResponse{
		Comments: items,
		Pagination: message.PaginationInfo{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  totalPages,
			TotalCount:  int(total),
		},
	}, nil
}

func loadUsersByIDs(ctx context.Context, idSet map[int32]struct{}) (map[int32]*model_def.User, error) {
	if len(idSet) == 0 {
		return map[int32]*model_def.User{}, nil
	}
	ids := make([]int32, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	var users []*model_def.User
	db := dao.User.WithContext(ctx).UnderlyingDB()
	if err := db.Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	m := make(map[int32]*model_def.User, len(users))
	for _, u := range users {
		m[u.ID] = u
	}
	return m, nil
}

func buildCommentItem(row *model_def.Comment, users map[int32]*model_def.User, parentMap map[uint]*model_def.Comment) message.CommentItem {
	au := authorFromUser(users[row.UserID])
	edited := !row.UpdatedAt.Equal(row.CreatedAt)
	item := message.CommentItem{
		ID:        row.ID,
		Content:   row.Content,
		CreatedAt: row.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: row.UpdatedAt.Format("2006-01-02 15:04:05"),
		Edited:    edited,
		ParentID:  row.ParentID,
		Author:    au,
	}
	if row.ParentID != nil {
		pid := uint(*row.ParentID)
		if p, ok := parentMap[pid]; ok {
			pu := users[p.UserID]
			item.ReplyTo = &message.CommentReplyTo{
				UserID:   p.UserID,
				Nickname: nicknameOf(pu),
			}
		}
	}
	return item
}

func authorFromUser(u *model_def.User) message.CommentAuthor {
	if u == nil {
		return message.CommentAuthor{}
	}
	return message.CommentAuthor{
		UserID:   u.ID,
		Username: u.Username,
		Nickname: nicknameOf(u),
		Avatar:   u.Avatar,
	}
}

func nicknameOf(u *model_def.User) string {
	if u == nil {
		return ""
	}
	if strings.TrimSpace(u.Nickname) != "" {
		return u.Nickname
	}
	return u.Username
}

// CreateComment 发表评论
func CreateComment(c *gin.Context, articleID int32, userID int32, req *message.CommentCreateRequest) (*message.CommentCreated, error) {
	ctx := c.Request.Context()

	content := strings.TrimSpace(req.Content)
	if len([]rune(content)) < commentMinLen {
		return nil, errors.New("评论内容过短")
	}
	if len([]rune(content)) > commentMaxLen {
		return nil, errors.New("评论内容过长")
	}

	_, err := dao.Article.WithContext(ctx).Where(dao.Article.ID.Eq(articleID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	var parentID *int32
	if req.ParentID != nil {
		pid := uint(*req.ParentID)
		_, err := dao.Comment.WithContext(ctx).Where(dao.Comment.ID.Eq(pid), dao.Comment.ArticleID.Eq(articleID)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("父评论不存在")
			}
			return nil, err
		}
		parentID = req.ParentID
	}

	m := &model_def.Comment{
		ArticleID: articleID,
		UserID:    userID,
		Content:   content,
		ParentID:  parentID,
	}
	if err := dao.Comment.WithContext(ctx).Create(m); err != nil {
		return nil, err
	}

	return &message.CommentCreated{
		ID:        m.ID,
		Content:   m.Content,
		CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// DeleteComment 删除评论（作者或管理员）
func DeleteComment(c *gin.Context, commentID uint, userID int32, role string) error {
	ctx := c.Request.Context()
	row, err := dao.Comment.WithContext(ctx).Where(dao.Comment.ID.Eq(commentID)).First()
	if err != nil {
		return err
	}
	if row.UserID != userID && role != "admin" {
		return errors.New("forbidden")
	}
	_, err = dao.Comment.WithContext(ctx).Where(dao.Comment.ID.Eq(commentID)).Delete()
	return err
}
